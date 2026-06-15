package protocol

import (
	"crypto/ed25519"
	"testing"
)

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
		t.Fatalf("ParseEnvelope should still reject transport-only traffic without a proof")
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

func TestCOSESign1RejectsWrongAlgorithm(t *testing.T) {
	protectedHeader, err := testCOSEProtectedHeader(-7)
	if err != nil {
		t.Fatalf("protected header: %v", err)
	}
	coseBytes, err := testCOSESign1WithHeader([]byte("payload"), protectedHeader, nil)
	if err != nil {
		t.Fatalf("encode COSE: %v", err)
	}
	if err := VerifyCOSESign1(coseBytes, nil, "alice"); err == nil {
		t.Fatalf("COSE with wrong protected alg should fail")
	}
}

func TestCOSESign1RejectsUnprotectedHeader(t *testing.T) {
	protectedHeader, err := testCOSEProtectedHeader(coseAlgorithmEdDSA)
	if err != nil {
		t.Fatalf("protected header: %v", err)
	}
	unprotectedHeader, err := testCOSEUnprotectedHeader()
	if err != nil {
		t.Fatalf("unprotected header: %v", err)
	}
	coseBytes, err := testCOSESign1WithHeader([]byte("payload"), protectedHeader, unprotectedHeader)
	if err != nil {
		t.Fatalf("encode COSE: %v", err)
	}
	if err := VerifyCOSESign1(coseBytes, nil, "alice"); err == nil {
		t.Fatalf("COSE with unprotected header should fail")
	}
}

func TestCOSESign1RejectsMismatchedDetachedPayload(t *testing.T) {
	detached, err := EncodeCOSESign1([]byte("payload-a"), "alice", true)
	if err != nil {
		t.Fatalf("encode detached COSE: %v", err)
	}
	if err := VerifyCOSESign1(detached, []byte("payload-b"), "alice"); err == nil {
		t.Fatalf("detached COSE with mismatched payload should fail")
	}
}

func testCOSEProtectedHeader(algorithm int64) ([]byte, error) {
	writer := &cborWriter{}
	if err := writer.writeMapHeader(1); err != nil {
		return nil, err
	}
	if err := writer.writeSignedInt(coseHeaderAlg); err != nil {
		return nil, err
	}
	if err := writer.writeSignedInt(algorithm); err != nil {
		return nil, err
	}
	return writer.buffer.Bytes(), nil
}

func testCOSEUnprotectedHeader() ([]byte, error) {
	writer := &cborWriter{}
	if err := writer.writeMapHeader(1); err != nil {
		return nil, err
	}
	if err := writer.writeSignedInt(4); err != nil {
		return nil, err
	}
	if err := writer.writeString("kid"); err != nil {
		return nil, err
	}
	return writer.buffer.Bytes(), nil
}

func testCOSESign1WithHeader(payload, protectedHeader, unprotectedHeader []byte) ([]byte, error) {
	toSign, signErr := encodeCOSESignatureStructure(protectedHeader, payload)
	if signErr != nil {
		return nil, signErr
	}
	writer := &cborWriter{}
	if err := writer.writeTag(coseSign1Tag); err != nil {
		return nil, err
	}
	if err := writer.writeArrayHeader(4); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(protectedHeader); err != nil {
		return nil, err
	}
	if len(unprotectedHeader) == 0 {
		if err := writer.writeMapHeader(0); err != nil {
			return nil, err
		}
	} else if err := writer.writeRawCBOR(unprotectedHeader); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(payload); err != nil {
		return nil, err
	}
	signature := ed25519.Sign(DeterministicPrivateKey("alice"), toSign)
	if err := writer.writeBytes(signature); err != nil {
		return nil, err
	}
	return writer.buffer.Bytes(), nil
}
