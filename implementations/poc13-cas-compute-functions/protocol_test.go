package poc13

import "testing"

func TestEnvelopeRoundTripUsesGridTag42AndPromisePayload(t *testing.T) {
	registry := NewRegistry()
	envelope, err := NewEnvelope(registry.MustCID(CASStorageV1), map[string]string{
		"act":           "promise",
		"variant":       "store_acceptance",
		"promise_about": "store_content",
		"content_cid":   ContentCID([]byte("hello")),
	}, "alice")
	if err != nil {
		t.Fatalf("new envelope: %v", err)
	}
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("envelope bytes: %v", err)
	}
	if len(envelopeBytes) < 4 || envelopeBytes[0] != 0xda || envelopeBytes[1] != 0x67 || envelopeBytes[2] != 0x72 || envelopeBytes[3] != 0x69 {
		t.Fatalf("envelope does not start with CBOR grid tag bytes: %x", envelopeBytes[:4])
	}
	parsed, err := ParseEnvelope(envelopeBytes)
	if err != nil {
		t.Fatalf("parse envelope: %v", err)
	}
	if err := VerifyEnvelope(parsed); err != nil {
		t.Fatalf("verify envelope: %v", err)
	}
	fields, err := parsed.PayloadFields()
	if err != nil {
		t.Fatalf("payload fields: %v", err)
	}
	if fields["act"] != "promise" {
		t.Fatalf("payload act = %q, want promise", fields["act"])
	}
}

func TestContentCIDRejectsCorruptBytes(t *testing.T) {
	cid := ContentCID(sampleContentBytes())
	if !VerifyContentCID(sampleContentBytes(), cid) {
		t.Fatalf("sample content should verify")
	}
	if VerifyContentCID(corruptContentBytes(), cid) {
		t.Fatalf("corrupt content should not verify")
	}
}

func TestComputeCacheKeyIncludesContext(t *testing.T) {
	keyA := ComputeCacheKey(CIDComputeV1, "function", "input", "context-a", "result")
	keyB := ComputeCacheKey(CIDComputeV1, "function", "input", "context-b", "result")
	if keyA == keyB {
		t.Fatalf("cache key must include context identity")
	}
}
