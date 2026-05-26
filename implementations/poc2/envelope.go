package main

import "fmt"

// Envelope is the local POC representation of grid([42(pCID), payload, ...]).
// Payload is already protocol-owned CBOR bytes; the envelope layer only carries
// and identifies it.
//
// Intent: Keep the kernel/app and kernel/kernel semantic boundary identical to
// the wire-shaped pCID envelope under test. Source: DI-ratij; DI-tijat.
type Envelope struct {
	ProtocolCID ProtocolCID
	Payload     []byte
}

func NewEnvelope(protocolCID ProtocolCID, fields map[string]string) (Envelope, error) {
	payloadBytes, marshalErr := marshalStringMap(fields)
	if marshalErr != nil {
		return Envelope{}, marshalErr
	}
	return Envelope{ProtocolCID: protocolCID, Payload: payloadBytes}, nil
}

func (envelope Envelope) Bytes() ([]byte, error) {
	writer := &cborWriter{}
	if err := writer.writeArrayHeader(2); err != nil {
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
	return writer.buffer.Bytes(), nil
}

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
	if arrayLength > 2 {
		return Envelope{}, fmt.Errorf("poc2 v0 does not parse later envelope slots")
	}
	if reader.offset != len(reader.data) {
		return Envelope{}, fmt.Errorf("trailing cbor bytes in envelope: %d", len(reader.data)-reader.offset)
	}
	cidBytes := make([]byte, len(tagBytes)-1)
	copy(cidBytes, tagBytes[1:])
	return Envelope{ProtocolCID: ProtocolCID{cidBytes: cidBytes}, Payload: payloadBytes}, nil
}

func (envelope Envelope) PayloadFields() (map[string]string, error) {
	return unmarshalStringMap(envelope.Payload)
}
