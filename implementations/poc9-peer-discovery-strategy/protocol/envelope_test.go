package protocol

import (
	"bytes"
	"testing"
)

// cborGridTagPrefix is the deterministic CBOR encoding of tag 0x67726964.
var cborGridTagPrefix = []byte{0xda, 0x67, 0x72, 0x69, 0x64}

func TestSignedEnvelopeRoundTrip(t *testing.T) {
	protocolCID := NewProtocolCID([]byte("poc9 test protocol"))
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
	// Intent: The complete envelope must be CBOR tag("grid") wrapping the slot
	// vector, not just the bare array. Source: DI-hanih
	if !bytes.HasPrefix(envelopeBytes, cborGridTagPrefix) {
		t.Fatalf("envelope is not wrapped in grid tag: %x", envelopeBytes)
	}
	if !bytes.Contains(envelopeBytes, []byte{0xd8, 0x2a}) {
		t.Fatalf("envelope does not contain tag 42: %x", envelopeBytes)
	}
	signableBytes, signableErr := envelope.SignableBytes()
	if signableErr != nil {
		t.Fatalf("signable bytes: %v", signableErr)
	}
	// Intent: Proofs cover the tagged grid signable view so signing cannot drift
	// from the byte-level envelope shape. Source: DI-hanih
	if !bytes.HasPrefix(signableBytes, cborGridTagPrefix) {
		t.Fatalf("signable view is not wrapped in grid tag: %x", signableBytes)
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

func TestParseEnvelopeRejectsUntaggedSlotVector(t *testing.T) {
	protocolCID := NewProtocolCID([]byte("poc9 untagged protocol"))
	envelope, envelopeErr := NewEnvelope(protocolCID, map[string]string{"kind": "test_v1"}, "alice")
	if envelopeErr != nil {
		t.Fatalf("new envelope: %v", envelopeErr)
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		t.Fatalf("bytes: %v", bytesErr)
	}
	// Intent: A bare [42(pCID), payload, proof] array is the historical bug. POC9
	// rejects it so tests fail if the outer grid tag disappears again. Source:
	// DI-hanih
	if _, parseErr := ParseEnvelope(envelopeBytes[len(cborGridTagPrefix):]); parseErr == nil {
		t.Fatalf("untagged slot vector unexpectedly parsed")
	}
}

func TestEnvelopeSignatureDetectsPayloadMutation(t *testing.T) {
	protocolCID := NewProtocolCID([]byte("poc9 mutation protocol"))
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
