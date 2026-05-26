package lib

import (
	"crypto/ed25519"
	"testing"
)

func TestSignatureSlotRoundTrip(t *testing.T) {
	publicKey, privateKey, keyErr := ed25519.GenerateKey(nil)
	if keyErr != nil {
		t.Fatalf("GenerateKey: %v", keyErr)
	}
	message := []byte("signed bytes")
	signature := ed25519.Sign(privateKey, message)
	decoded, decodeErr := DecodeSignatureSlot(EncodeSignatureSlot(SignatureSlot{PublicKey: publicKey, Signature: signature}))
	if decodeErr != nil {
		t.Fatalf("DecodeSignatureSlot: %v", decodeErr)
	}
	if !ed25519.Verify(ed25519.PublicKey(decoded.PublicKey), message, decoded.Signature) {
		t.Fatalf("decoded signature did not verify")
	}
}
