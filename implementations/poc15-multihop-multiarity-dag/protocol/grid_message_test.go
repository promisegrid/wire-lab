package protocol

import "testing"

func TestGridMessageParsesVariableArity(t *testing.T) {
	// Intent: POC15 must be able to inspect pCID-owned slot vectors without
	// assuming every grid message is the current payload/proof pair. Source:
	// DI-mosat
	protocolCID := NewProtocolCID([]byte("poc15 variable arity test protocol"))
	payloadSlot, err := ByteStringGridSlot([]byte("payload bytes"))
	if err != nil {
		t.Fatalf("payload slot: %v", err)
	}
	parentLinks, err := EncodeTag42LinkArray(RawCIDV1SHA256Bytes([]byte("parent envelope")))
	if err != nil {
		t.Fatalf("parent link array: %v", err)
	}
	proofSlot, err := ByteStringGridSlot([]byte("proof bytes"))
	if err != nil {
		t.Fatalf("proof slot: %v", err)
	}
	messageBytes, err := EncodeGridMessage(protocolCID, RawCBORGridSlot(parentLinks), payloadSlot, proofSlot)
	if err != nil {
		t.Fatalf("encode grid message: %v", err)
	}
	parsed, err := ParseGridMessage(messageBytes)
	if err != nil {
		t.Fatalf("parse grid message: %v", err)
	}
	if !parsed.ProtocolCID.Equal(protocolCID) {
		t.Fatalf("parsed pCID mismatch")
	}
	if len(parsed.Slots) != 3 {
		t.Fatalf("parsed slots=%d want 3", len(parsed.Slots))
	}
	if payloadBytes, ok := parsed.Slots[1].ByteString(); !ok || string(payloadBytes) != "payload bytes" {
		t.Fatalf("payload slot did not round trip as byte string")
	}
}

func TestGridMessageAcceptsTransportOnlyShape(t *testing.T) {
	protocolCID := NewProtocolCID([]byte("poc15 transport only test protocol"))
	payloadSlot, err := ByteStringGridSlot([]byte("session-authenticated payload"))
	if err != nil {
		t.Fatalf("payload slot: %v", err)
	}
	messageBytes, err := EncodeGridMessage(protocolCID, payloadSlot)
	if err != nil {
		t.Fatalf("encode grid message: %v", err)
	}
	parsed, err := ParseGridMessage(messageBytes)
	if err != nil {
		t.Fatalf("parse grid message: %v", err)
	}
	if len(parsed.Slots) != 1 {
		t.Fatalf("parsed slots=%d want 1", len(parsed.Slots))
	}
	if _, parseErr := ParseEnvelope(messageBytes); parseErr == nil {
		t.Fatalf("legacy ParseEnvelope should still reject non-three-slot traffic")
	}
}

func TestCOSESign1VerifiesAndRejectsTampering(t *testing.T) {
	payload := []byte("poc15 COSE payload")
	embedded, err := EncodeCOSESign1(payload, "alice", false)
	if err != nil {
		t.Fatalf("encode embedded COSE: %v", err)
	}
	if err := VerifyCOSESign1(embedded, nil, "alice"); err != nil {
		t.Fatalf("verify embedded COSE: %v", err)
	}
	detached, err := EncodeCOSESign1(payload, "alice", true)
	if err != nil {
		t.Fatalf("encode detached COSE: %v", err)
	}
	if err := VerifyCOSESign1(detached, payload, "alice"); err != nil {
		t.Fatalf("verify detached COSE: %v", err)
	}
	tampered := append([]byte(nil), detached...)
	tampered[len(tampered)-1] ^= 0x01
	if err := VerifyCOSESign1(tampered, payload, "alice"); err == nil {
		t.Fatalf("tampered detached COSE should fail verification")
	}
}
