package protocol

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GridSlot is one already-encoded CBOR item after slot 0 in a PromiseGrid grid
// message.
// Intent: POC16 needs to test pCID-owned arity and slot meaning without making
// the legacy `payload, proof` pair the only parseable envelope shape. Source:
// DI-mosat
type GridSlot struct {
	RawCBOR []byte
}

// GridMessage is a parsed grid message whose slot-0 pCID names the protocol
// spec that defines every following slot.
// Intent: The generic parser only identifies the pCID and preserves raw slot
// CBOR; protocol-specific code still owns the meaning of those slots. Source:
// DI-mosat
type GridMessage struct {
	ProtocolCID ProtocolCID
	Slots       []GridSlot
}

// ByteStringGridSlot wraps bytes as one CBOR byte-string slot.
func ByteStringGridSlot(value []byte) (GridSlot, error) {
	writer := &cborWriter{}
	if err := writer.writeBytes(value); err != nil {
		return GridSlot{}, err
	}
	return RawCBORGridSlot(writer.buffer.Bytes()), nil
}

// RawCBORGridSlot preserves a complete CBOR item as one pCID-owned slot.
func RawCBORGridSlot(value []byte) GridSlot {
	copiedValue := make([]byte, len(value))
	copy(copiedValue, value)
	return GridSlot{RawCBOR: copiedValue}
}

// ByteString returns the decoded value when this slot is exactly one CBOR byte
// string.
func (slot GridSlot) ByteString() ([]byte, bool) {
	reader := &cborReader{data: slot.RawCBOR}
	value, err := reader.readBytes()
	if err != nil || reader.offset != len(reader.data) {
		return nil, false
	}
	return value, true
}

// EncodeGridMessage writes grid([42(pCID), ...protocol-defined-slots]).
func EncodeGridMessage(protocolCID ProtocolCID, slots ...GridSlot) ([]byte, error) {
	writer := &cborWriter{}
	if err := writer.writeTag(gridTag); err != nil {
		return nil, err
	}
	if err := writer.writeArrayHeader(1 + len(slots)); err != nil {
		return nil, err
	}
	if err := writer.writeTag(dagCBORLinkTag); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(protocolCID.Tag42Bytes()); err != nil {
		return nil, err
	}
	for slotIndex, slot := range slots {
		if len(slot.RawCBOR) == 0 {
			return nil, fmt.Errorf("grid slot %d is empty", slotIndex+1)
		}
		if err := writer.writeRawCBOR(slot.RawCBOR); err != nil {
			return nil, err
		}
	}
	return writer.buffer.Bytes(), nil
}

// ParseGridMessage parses any definite-length grid array with slot-0 tag 42 and
// preserves every later slot as raw CBOR for pCID-specific dispatch.
func ParseGridMessage(messageBytes []byte) (GridMessage, error) {
	reader := &cborReader{data: messageBytes}
	outerTag, outerTagErr := reader.readTypeAndLength(6)
	if outerTagErr != nil {
		return GridMessage{}, outerTagErr
	}
	if outerTag != gridTag {
		return GridMessage{}, fmt.Errorf("outer envelope must be grid tag %d, got tag %d", gridTag, outerTag)
	}
	arrayLength, arrayErr := reader.readTypeAndLength(4)
	if arrayErr != nil {
		return GridMessage{}, arrayErr
	}
	if arrayLength == 0 {
		return GridMessage{}, fmt.Errorf("grid message must include slot 0 pCID")
	}
	tagNumber, tagErr := reader.readTypeAndLength(6)
	if tagErr != nil {
		return GridMessage{}, tagErr
	}
	if tagNumber != dagCBORLinkTag {
		return GridMessage{}, fmt.Errorf("slot 0 must be tag 42, got tag %d", tagNumber)
	}
	tagBytes, tagBytesErr := reader.readBytes()
	if tagBytesErr != nil {
		return GridMessage{}, tagBytesErr
	}
	if len(tagBytes) < 2 || tagBytes[0] != 0x00 {
		return GridMessage{}, fmt.Errorf("tag 42 payload must start with DAG-CBOR CID sentinel")
	}
	slots := make([]GridSlot, 0, int(arrayLength-1))
	for slotIndex := uint64(1); slotIndex < arrayLength; slotIndex++ {
		rawSlot, rawSlotErr := reader.readRawItem()
		if rawSlotErr != nil {
			return GridMessage{}, fmt.Errorf("slot %d: %w", slotIndex, rawSlotErr)
		}
		slots = append(slots, RawCBORGridSlot(rawSlot))
	}
	if reader.offset != len(reader.data) {
		return GridMessage{}, fmt.Errorf("trailing cbor bytes in grid message: %d", len(reader.data)-reader.offset)
	}
	return GridMessage{ProtocolCID: NewProtocolCIDFromBytes(tagBytes[1:]), Slots: slots}, nil
}

// RawCIDV1SHA256Bytes returns the CIDv1 raw sha2-256 bytes for any exact object
// bytes. This is intentionally not a ProtocolCID; parent links name message
// objects, while pCID names protocol specs.
func RawCIDV1SHA256Bytes(content []byte) []byte {
	digest := sha256.Sum256(content)
	cidBytes := make([]byte, 0, 36)
	cidBytes = append(cidBytes, 0x01, 0x55, 0x12, 0x20)
	cidBytes = append(cidBytes, digest[:]...)
	return cidBytes
}

// EncodeTag42LinkArray writes a CBOR array of DAG-CBOR tag-42 links.
func EncodeTag42LinkArray(cidValues ...[]byte) ([]byte, error) {
	writer := &cborWriter{}
	if err := writer.writeArrayHeader(len(cidValues)); err != nil {
		return nil, err
	}
	for linkIndex, cidBytes := range cidValues {
		if len(cidBytes) == 0 {
			return nil, fmt.Errorf("tag42 link %d is empty", linkIndex)
		}
		if err := writer.writeTag(dagCBORLinkTag); err != nil {
			return nil, err
		}
		tagBytes := make([]byte, 0, len(cidBytes)+1)
		tagBytes = append(tagBytes, 0x00)
		tagBytes = append(tagBytes, cidBytes...)
		if err := writer.writeBytes(tagBytes); err != nil {
			return nil, err
		}
	}
	return writer.buffer.Bytes(), nil
}

// DecodeTag42LinkArray reads the narrow parent-link shape used by POC16 route
// traffic: an array of DAG-CBOR tag-42 CID links.
// Intent: Normal route traffic can carry pCID-owned parent slots without making
// parent links a universal envelope requirement for every pCID. Source:
// DI-kohuj
func DecodeTag42LinkArray(arrayBytes []byte) ([][]byte, error) {
	reader := &cborReader{data: arrayBytes}
	arrayLength, arrayErr := reader.readTypeAndLength(4)
	if arrayErr != nil {
		return nil, arrayErr
	}
	cidValues := make([][]byte, 0, int(arrayLength))
	for linkIndex := uint64(0); linkIndex < arrayLength; linkIndex++ {
		tagNumber, tagErr := reader.readTypeAndLength(6)
		if tagErr != nil {
			return nil, fmt.Errorf("parent link %d: %w", linkIndex, tagErr)
		}
		if tagNumber != dagCBORLinkTag {
			return nil, fmt.Errorf("parent link %d must be tag 42, got tag %d", linkIndex, tagNumber)
		}
		tagBytes, bytesErr := reader.readBytes()
		if bytesErr != nil {
			return nil, fmt.Errorf("parent link %d: %w", linkIndex, bytesErr)
		}
		if len(tagBytes) < 2 || tagBytes[0] != 0x00 {
			return nil, fmt.Errorf("parent link %d tag 42 payload must start with DAG-CBOR CID sentinel", linkIndex)
		}
		cidBytes := make([]byte, len(tagBytes)-1)
		copy(cidBytes, tagBytes[1:])
		cidValues = append(cidValues, cidBytes)
	}
	if reader.offset != len(reader.data) {
		return nil, fmt.Errorf("trailing cbor bytes in tag42 link array: %d", len(reader.data)-reader.offset)
	}
	return cidValues, nil
}

// ExactSHA256FromRawCIDV1 returns the hex digest when a POC raw CIDv1 sha2-256
// link names exact message bytes retained in the run CAS.
// Intent: Parent links should name exact retained message bytes, not payload
// bytes or protocol CIDs, so the analyzer can traverse byte-faithful history.
// Source: DI-kohuj
func ExactSHA256FromRawCIDV1(cidBytes []byte) (string, bool) {
	if len(cidBytes) != 36 || cidBytes[0] != 0x01 || cidBytes[1] != 0x55 || cidBytes[2] != 0x12 || cidBytes[3] != 0x20 {
		return "", false
	}
	return hex.EncodeToString(cidBytes[4:]), true
}
