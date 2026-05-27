package protocol

import (
	"bytes"
	"testing"
)

func TestSignedEnvelopeRoundTrip(t *testing.T) {
	protocolCID := NewProtocolCID([]byte("poc7 test protocol"))
	envelope, envelopeErr := NewEnvelope(protocolCID, map[string]string{
		"kind": "test_v1",
		"body": "hello",
	}, "alice")
	if envelopeErr != nil {
		t.Fatalf("new envelope: %v", envelopeErr)
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		t.Fatalf("bytes: %v", bytesErr)
	}
	if !bytes.Contains(envelopeBytes, []byte{0xd8, 0x2a}) {
		t.Fatalf("envelope does not contain tag 42: %x", envelopeBytes)
	}
	parsed, parseErr := ParseEnvelope(envelopeBytes)
	if parseErr != nil {
		t.Fatalf("parse: %v", parseErr)
	}
	if !parsed.ProtocolCID.Equal(protocolCID) {
		t.Fatalf("pCID changed: %s != %s", parsed.ProtocolCID, protocolCID)
	}
	if verifyErr := VerifyEnvelope(parsed); verifyErr != nil {
		t.Fatalf("verify: %v", verifyErr)
	}
	fields, fieldsErr := parsed.PayloadFields()
	if fieldsErr != nil {
		t.Fatalf("fields: %v", fieldsErr)
	}
	if fields["kind"] != "test_v1" || fields["body"] != "hello" {
		t.Fatalf("unexpected fields: %#v", fields)
	}
}

func TestEnvelopeSignatureDetectsPayloadMutation(t *testing.T) {
	protocolCID := NewProtocolCID([]byte("poc7 mutation protocol"))
	envelope, envelopeErr := NewEnvelope(protocolCID, map[string]string{"kind": "test_v1"}, "alice")
	if envelopeErr != nil {
		t.Fatalf("new envelope: %v", envelopeErr)
	}
	envelope.Payload, envelopeErr = MarshalStringMap(map[string]string{"kind": "rewritten_v1"})
	if envelopeErr != nil {
		t.Fatalf("rewrite payload: %v", envelopeErr)
	}
	if verifyErr := VerifyEnvelope(envelope); verifyErr == nil {
		t.Fatalf("mutated payload unexpectedly verified")
	}
}
