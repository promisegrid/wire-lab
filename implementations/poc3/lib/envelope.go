package lib

import "fmt"

// Envelope is the local POC representation of grid([42(pCID), payload, ...]).
// Extra slots are protocol-owned proof or extension bytes, not kernel commands.
//
// Intent: Keep app/kernel and kernel/kernel boundaries wire-shaped while
// allowing the signed app to pressure-test a proof slot. Source: DI-horak.
type Envelope struct {
	ProtocolCID ProtocolCID
	Payload     []byte
	ExtraSlots  [][]byte
}

// NewEnvelope builds an envelope around deterministic protocol-owned fields.
func NewEnvelope(protocolCID ProtocolCID, fields map[string]string, extraSlots ...[]byte) (Envelope, error) {
	payloadBytes, marshalErr := MarshalStringMap(fields)
	if marshalErr != nil {
		return Envelope{}, marshalErr
	}
	copiedSlots := make([][]byte, 0, len(extraSlots))
	for _, slot := range extraSlots {
		copiedSlot := make([]byte, len(slot))
		copy(copiedSlot, slot)
		copiedSlots = append(copiedSlots, copiedSlot)
	}
	return Envelope{ProtocolCID: protocolCID, Payload: payloadBytes, ExtraSlots: copiedSlots}, nil
}

// Bytes serializes the envelope as CBOR grid([42(pCID), payload, ...]).
func (envelope Envelope) Bytes() ([]byte, error) {
	writer := &cborWriter{}
	if err := writer.writeArrayHeader(2 + len(envelope.ExtraSlots)); err != nil {
		return nil, err
	}
	if err := writer.writeTag(42); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(envelope.ProtocolCID.Tag42Bytes()); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(envelope.Payload); err != nil {
		return nil, err
	}
	for _, slot := range envelope.ExtraSlots {
		if err := writer.writeBytes(slot); err != nil {
			return nil, err
		}
	}
	return writer.buffer.Bytes(), nil
}

// ParseEnvelope reads the POC envelope and keeps every later slot as raw bytes.
func ParseEnvelope(envelopeBytes []byte) (Envelope, error) {
	reader := &cborReader{data: envelopeBytes}
	arrayLength, arrayErr := reader.readTypeAndLength(4)
	if arrayErr != nil {
		return Envelope{}, arrayErr
	}
	if arrayLength < 2 {
		return Envelope{}, fmt.Errorf("grid envelope must have at least two slots, got %d", arrayLength)
	}
	tagNumber, tagErr := reader.readTypeAndLength(6)
	if tagErr != nil {
		return Envelope{}, tagErr
	}
	if tagNumber != 42 {
		return Envelope{}, fmt.Errorf("slot 0 must be tag 42, got tag %d", tagNumber)
	}
	tagBytes, tagBytesErr := reader.readBytes()
	if tagBytesErr != nil {
		return Envelope{}, tagBytesErr
	}
	if len(tagBytes) < 2 || tagBytes[0] != 0x00 {
		return Envelope{}, fmt.Errorf("tag 42 payload must start with DAG-CBOR CID sentinel")
	}
	payloadBytes, payloadErr := reader.readBytes()
	if payloadErr != nil {
		return Envelope{}, payloadErr
	}
	extraSlots := make([][]byte, 0, int(arrayLength)-2)
	for index := uint64(2); index < arrayLength; index++ {
		slotBytes, slotErr := reader.readBytes()
		if slotErr != nil {
			return Envelope{}, slotErr
		}
		extraSlots = append(extraSlots, slotBytes)
	}
	if reader.offset != len(reader.data) {
		return Envelope{}, fmt.Errorf("trailing cbor bytes in envelope: %d", len(reader.data)-reader.offset)
	}
	return Envelope{ProtocolCID: NewProtocolCIDFromBytes(tagBytes[1:]), Payload: payloadBytes, ExtraSlots: extraSlots}, nil
}

// PayloadFields decodes protocol-owned payload bytes as the POC string map.
func (envelope Envelope) PayloadFields() (map[string]string, error) {
	return UnmarshalStringMap(envelope.Payload)
}

// EnvelopeKind returns the payload kind and decoded fields.
func EnvelopeKind(envelope Envelope) (string, map[string]string, error) {
	fields, payloadErr := envelope.PayloadFields()
	if payloadErr != nil {
		return "", nil, payloadErr
	}
	kind := fields["kind"]
	if kind == "" {
		return "", nil, fmt.Errorf("payload missing kind")
	}
	return kind, fields, nil
}
