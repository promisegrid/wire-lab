package signed

import (
	"crypto/ed25519"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc3/lib"
)

func TestSignedEnvelopeVerifiesExactSignableBytes(t *testing.T) {
	signableEnvelope, envelopeErr := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":      "signed_note_v1",
		"from":      "alice-signed-app",
		"from_node": "alice",
		"to":        "bob",
		"text":      "signed hello",
	})
	if envelopeErr != nil {
		t.Fatalf("NewEnvelope: %v", envelopeErr)
	}
	signableBytes, bytesErr := signableEnvelope.Bytes()
	if bytesErr != nil {
		t.Fatalf("Envelope.Bytes: %v", bytesErr)
	}
	publicKey, privateKey := deterministicKeyPair("alice:alice-signed-app")
	signature := ed25519.Sign(privateKey, signableBytes)
	if !ed25519.Verify(publicKey, signableBytes, signature) {
		t.Fatalf("signature did not verify")
	}
}
