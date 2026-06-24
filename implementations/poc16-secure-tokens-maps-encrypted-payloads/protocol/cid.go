package protocol

import (
	"bytes"
	"fmt"

	cidlib "github.com/ipfs/go-cid"
	mbase "github.com/multiformats/go-multibase"
	mh "github.com/multiformats/go-multihash"
)

// ProtocolCID identifies the protocol spec document that defines every slot
// after slot 0. It is not a payload CID.
// Intent: Keep POC16 aligned with grid([42(pCID), ...protocol-defined-slots])
// while making the CID library, not POC-local digest parsing, authoritative for
// pCID text and byte handling. Source: DI-timah; DI-sazip
type ProtocolCID struct {
	value cidlib.Cid
}

// NewProtocolCID content-addresses spec bytes into the supported CID profile.
func NewProtocolCID(specBytes []byte) ProtocolCID {
	return ProtocolCID{value: CIDForBytes(specBytes)}
}

// MustProtocolCIDText returns the authoritative pCID from a base32 CID string.
func MustProtocolCIDText(cidText string) ProtocolCID {
	protocolCID, err := ParseProtocolCIDText(cidText)
	if err != nil {
		panic(err)
	}
	return protocolCID
}

// ParseProtocolCIDText validates a printable pCID and returns its binary form.
func ParseProtocolCIDText(cidText string) (ProtocolCID, error) {
	parsedCID, parseErr := ParseCIDText(cidText)
	if parseErr != nil {
		return ProtocolCID{}, parseErr
	}
	return ProtocolCID{value: parsedCID}, nil
}

// ParseProtocolCIDFromBytes rebuilds a pCID from binary CIDv1 bytes.
func ParseProtocolCIDFromBytes(cidBytes []byte) (ProtocolCID, error) {
	parsedCID, parseErr := ParseCIDBytes(cidBytes)
	if parseErr != nil {
		return ProtocolCID{}, parseErr
	}
	return ProtocolCID{value: parsedCID}, nil
}

// Bytes returns binary CIDv1 bytes without the tag-42 sentinel.
func (protocolCID ProtocolCID) Bytes() []byte {
	return append([]byte(nil), protocolCID.value.Bytes()...)
}

// Tag42Bytes returns the DAG-CBOR tag-42 byte string payload.
func (protocolCID ProtocolCID) Tag42Bytes() []byte {
	cidBytes := protocolCID.Bytes()
	tagBytes := make([]byte, 0, len(cidBytes)+1)
	tagBytes = append(tagBytes, 0x00)
	tagBytes = append(tagBytes, cidBytes...)
	return tagBytes
}

// String renders the canonical CIDv1 base32 text form.
func (protocolCID ProtocolCID) String() string {
	return CIDText(protocolCID.value)
}

// Equal reports whether two pCIDs name the same protocol bytes.
func (protocolCID ProtocolCID) Equal(other ProtocolCID) bool {
	return bytes.Equal(protocolCID.Bytes(), other.Bytes())
}

// CIDForBytes returns the CIDv1 raw sha2-256 identifier for exact object bytes.
func CIDForBytes(content []byte) cidlib.Cid {
	multihash, multihashErr := mh.Sum(content, mh.SHA2_256, -1)
	if multihashErr != nil {
		panic("sha2-256 multihash should not fail: " + multihashErr.Error())
	}
	return cidlib.NewCidV1(cidlib.Raw, multihash)
}

// CIDForExactBytes returns the printable CID for exact retained bytes.
func CIDForExactBytes(exactBytes []byte) string {
	if len(exactBytes) == 0 {
		return ""
	}
	return CIDText(CIDForBytes(exactBytes))
}

// ParseCIDText validates a printable CIDv1 raw sha2-256 identifier.
func ParseCIDText(cidText string) (cidlib.Cid, error) {
	parsedCID, decodeErr := cidlib.Decode(cidText)
	if decodeErr != nil {
		return cidlib.Undef, decodeErr
	}
	if validateErr := validateSupportedCIDProfile(parsedCID); validateErr != nil {
		return cidlib.Undef, validateErr
	}
	if CIDText(parsedCID) != cidText {
		return cidlib.Undef, fmt.Errorf("cid text must be canonical base32")
	}
	return parsedCID, nil
}

// CIDBytesForExactBytes returns binary CID bytes for exact object bytes using
// the POC16-supported CID profile.
func CIDBytesForExactBytes(content []byte) []byte {
	return CIDForBytes(content).Bytes()
}

// ParseCIDBytes validates binary CIDv1 raw sha2-256 bytes.
func ParseCIDBytes(cidBytes []byte) (cidlib.Cid, error) {
	parsedCID, castErr := cidlib.Cast(cidBytes)
	if castErr != nil {
		return cidlib.Undef, castErr
	}
	if validateErr := validateSupportedCIDProfile(parsedCID); validateErr != nil {
		return cidlib.Undef, validateErr
	}
	if !bytes.Equal(parsedCID.Bytes(), cidBytes) {
		return cidlib.Undef, fmt.Errorf("cid bytes are not canonical")
	}
	return parsedCID, nil
}

// CIDTextFromBytes returns the printable CIDv1 base32 form for binary CID bytes.
func CIDTextFromBytes(cidBytes []byte) (string, error) {
	parsedCID, parseErr := ParseCIDBytes(cidBytes)
	if parseErr != nil {
		return "", parseErr
	}
	return CIDText(parsedCID), nil
}

// CIDText returns canonical CIDv1 base32 text with the multibase b prefix.
func CIDText(value cidlib.Cid) string {
	text, textErr := value.StringOfBase(mbase.Base32)
	if textErr != nil {
		panic("cid base32 rendering should not fail: " + textErr.Error())
	}
	return text
}

func validateSupportedCIDProfile(value cidlib.Cid) error {
	if !value.Defined() {
		return fmt.Errorf("cid is undefined")
	}
	if value.Version() != 1 {
		return fmt.Errorf("cid must be CIDv1, got v%d", value.Version())
	}
	if value.Type() != cidlib.Raw {
		return fmt.Errorf("cid codec must be raw, got 0x%x", value.Type())
	}
	decodedMultihash, decodeErr := mh.Decode(value.Hash())
	if decodeErr != nil {
		return decodeErr
	}
	if decodedMultihash.Code != mh.SHA2_256 {
		return fmt.Errorf("cid multihash must be sha2-256, got 0x%x", decodedMultihash.Code)
	}
	if decodedMultihash.Length != 32 {
		return fmt.Errorf("cid sha2-256 length must be 32 bytes, got %d", decodedMultihash.Length)
	}
	return nil
}
