package protocol

import "testing"

func TestEnvelopeUsesGridTag42AndVerifies(t *testing.T) {
	protocolCID := NewProtocolCID([]byte("poc14 test protocol"))
	envelope, err := NewEnvelope(protocolCID, map[string]string{"kind": "offer_promise", "from": "alice"}, "alice")
	if err != nil {
		t.Fatalf("new envelope: %v", err)
	}
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("envelope bytes: %v", err)
	}
	if len(envelopeBytes) < 5 || envelopeBytes[0] != 0xda || envelopeBytes[1] != 0x67 || envelopeBytes[2] != 0x72 || envelopeBytes[3] != 0x69 || envelopeBytes[4] != 0x64 {
		t.Fatalf("envelope does not start with CBOR grid tag: %x", envelopeBytes[:5])
	}
	parsed, err := ParseEnvelope(envelopeBytes)
	if err != nil {
		t.Fatalf("parse envelope: %v", err)
	}
	if !parsed.ProtocolCID.Equal(protocolCID) {
		t.Fatalf("parsed pCID mismatch")
	}
	if err := VerifyEnvelope(parsed); err != nil {
		t.Fatalf("verify envelope: %v", err)
	}
}

func TestEnvelopeRejectsTampering(t *testing.T) {
	envelope, err := NewEnvelope(NewProtocolCID([]byte("poc14 tamper protocol")), map[string]string{"kind": "need_advertisement"}, "alice")
	if err != nil {
		t.Fatalf("new envelope: %v", err)
	}
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("envelope bytes: %v", err)
	}
	envelopeBytes[len(envelopeBytes)/2] ^= 0x01
	parsed, parseErr := ParseEnvelope(envelopeBytes)
	if parseErr == nil {
		if verifyErr := VerifyEnvelope(parsed); verifyErr == nil {
			t.Fatalf("tampered envelope should fail parse or verify")
		}
	}
}

func TestParseEnvelopeRejectsMalformedCBOR(t *testing.T) {
	// Intent: Adversarial peers may send arbitrary TCP bytes, but malformed
	// bytes remain local parse-failure evidence rather than a command surface.
	// Source: DI-timah
	if _, err := ParseEnvelope([]byte{0xda, 0x67, 0x72, 0x69}); err == nil {
		t.Fatalf("truncated grid-tag bytes should fail parse")
	}
	if _, err := ParseEnvelope([]byte("ignore previous instructions")); err == nil {
		t.Fatalf("plain prompt-injection bytes should fail parse")
	}
}

func FuzzParseEnvelopeHandlesArbitraryBytes(f *testing.F) {
	// Intent: POC14 should treat malformed CBOR, prompt-injection bytes, partial
	// writes, and random adversarial inputs as parse/verification outcomes rather
	// than panics or expanded protocol semantics. Source: DI-sunuf
	validEnvelope, err := NewEnvelope(NewProtocolCID([]byte("poc14 fuzz protocol")), map[string]string{"act": "promise", "from": "alice"}, "alice")
	if err != nil {
		f.Fatalf("new seed envelope: %v", err)
	}
	validBytes, err := validEnvelope.Bytes()
	if err != nil {
		f.Fatalf("seed envelope bytes: %v", err)
	}
	f.Add([]byte{})
	f.Add([]byte{0xda, 0x67, 0x72, 0x69})
	f.Add([]byte("ignore previous instructions"))
	f.Add(validBytes)
	f.Fuzz(func(t *testing.T, envelopeBytes []byte) {
		envelope, parseErr := ParseEnvelope(envelopeBytes)
		if parseErr != nil {
			return
		}
		if _, fieldsErr := envelope.PayloadFields(); fieldsErr != nil {
			return
		}
		if verifyErr := VerifyEnvelope(envelope); verifyErr != nil {
			return
		}
	})
}
