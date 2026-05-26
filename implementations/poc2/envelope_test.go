package main

import "testing"

func TestEnvelopeRoundTrip(t *testing.T) {
	protocolCID, cidErr := HelloProtocolCID()
	if cidErr != nil {
		t.Fatalf("HelloProtocolCID: %v", cidErr)
	}
	envelope, envelopeErr := NewEnvelope(protocolCID, map[string]string{
		"kind": "hello_v1",
		"from": "alice-hello-app",
		"to":   "bob",
		"text": "hello from Alice",
	})
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
		t.Fatalf("parsed pCID did not match hello protocol pCID")
	}
	fields, fieldsErr := parsedEnvelope.PayloadFields()
	if fieldsErr != nil {
		t.Fatalf("PayloadFields: %v", fieldsErr)
	}
	if fields["kind"] != "hello_v1" || fields["text"] != "hello from Alice" {
		t.Fatalf("unexpected payload fields: %#v", fields)
	}
}

func TestEnvelopeRequiresTag42(t *testing.T) {
	_, parseErr := ParseEnvelope([]byte{0x82, 0x58, 0x01, 0x00, 0x40})
	if parseErr == nil {
		t.Fatalf("expected tag-42 parse failure")
	}
}
