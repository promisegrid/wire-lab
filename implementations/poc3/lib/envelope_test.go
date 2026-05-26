package lib

import "testing"

func TestEnvelopeRoundTripWithExtraSlot(t *testing.T) {
	protocolCID := NewProtocolCID([]byte("test protocol"))
	envelope, envelopeErr := NewEnvelope(protocolCID, map[string]string{
		"kind": "test_v1",
		"text": "hello",
	}, []byte("proof-slot"))
	if envelopeErr != nil {
		t.Fatalf("NewEnvelope: %v", envelopeErr)
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		t.Fatalf("Envelope.Bytes: %v", bytesErr)
	}
	parsedEnvelope, parseErr := ParseEnvelope(envelopeBytes)
	if parseErr != nil {
		t.Fatalf("ParseEnvelope: %v", parseErr)
	}
	if !parsedEnvelope.ProtocolCID.Equal(protocolCID) {
		t.Fatalf("parsed pCID did not match")
	}
	if len(parsedEnvelope.ExtraSlots) != 1 || string(parsedEnvelope.ExtraSlots[0]) != "proof-slot" {
		t.Fatalf("unexpected extra slots: %#v", parsedEnvelope.ExtraSlots)
	}
	kind, fields, kindErr := EnvelopeKind(parsedEnvelope)
	if kindErr != nil {
		t.Fatalf("EnvelopeKind: %v", kindErr)
	}
	if kind != "test_v1" || fields["text"] != "hello" {
		t.Fatalf("unexpected payload: kind=%s fields=%#v", kind, fields)
	}
}

func TestEnvelopeRequiresTag42(t *testing.T) {
	_, parseErr := ParseEnvelope([]byte{0x82, 0x58, 0x01, 0x00, 0x40})
	if parseErr == nil {
		t.Fatalf("expected tag-42 parse failure")
	}
}
