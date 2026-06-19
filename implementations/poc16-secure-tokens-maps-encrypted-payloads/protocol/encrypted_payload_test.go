package protocol

import (
	"strings"
	"testing"
)

func TestEncryptedPayloadRoundTripAndRejectsWrongRecipient(t *testing.T) {
	payload, encryptErr := EncryptPayloadForRecipient("alice", "bob", "cas_storage_v1", []byte("private package bytes"))
	if encryptErr != nil {
		t.Fatalf("encrypt payload: %v", encryptErr)
	}
	encoded, marshalErr := payload.MarshalCBOR()
	if marshalErr != nil {
		t.Fatalf("marshal encrypted payload: %v", marshalErr)
	}
	decoded, decodeErr := EncryptedPayloadFromCBOR(encoded)
	if decodeErr != nil {
		t.Fatalf("decode encrypted payload: %v", decodeErr)
	}
	plaintext, decryptErr := DecryptPayloadForRecipient(decoded, "bob")
	if decryptErr != nil {
		t.Fatalf("decrypt payload: %v", decryptErr)
	}
	if string(plaintext) != "private package bytes" {
		t.Fatalf("plaintext = %q", plaintext)
	}
	if _, wrongRecipientErr := DecryptPayloadForRecipient(decoded, "carol"); wrongRecipientErr == nil || !strings.Contains(wrongRecipientErr.Error(), "recipient") {
		t.Fatalf("wrong recipient error = %v, want recipient failure", wrongRecipientErr)
	}
}

func TestEncryptedPayloadRejectsTamperedCiphertext(t *testing.T) {
	payload, encryptErr := EncryptPayloadForRecipient("alice", "bob", "cas_storage_v1", []byte("private package bytes"))
	if encryptErr != nil {
		t.Fatalf("encrypt payload: %v", encryptErr)
	}
	payload.CiphertextBase64 = payload.CiphertextBase64[:len(payload.CiphertextBase64)-2] + "AA"
	if _, decryptErr := DecryptPayloadForRecipient(payload, "bob"); decryptErr == nil {
		t.Fatalf("tampered ciphertext should fail")
	}
}
