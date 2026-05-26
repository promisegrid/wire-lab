package lib

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

// SignatureSlot contains the public key and signature for signed envelope bytes.
// Intent: Let the signed app witness exact bytes without claiming global trust.
// Source: DI-horak.
type SignatureSlot struct {
	PublicKey []byte
	Signature []byte
}

// EncodeSignatureSlot serializes the POC signature slot as public-key/signature
// bytes joined by a stable colon separator in hex form.
func EncodeSignatureSlot(slot SignatureSlot) []byte {
	return []byte(hex.EncodeToString(slot.PublicKey) + ":" + hex.EncodeToString(slot.Signature))
}

// DecodeSignatureSlot decodes the POC signature slot.
func DecodeSignatureSlot(slotBytes []byte) (SignatureSlot, error) {
	encoded := string(slotBytes)
	separator := -1
	for index, char := range encoded {
		if char == ':' {
			separator = index
			break
		}
	}
	if separator < 0 {
		return SignatureSlot{}, fmt.Errorf("signature slot missing separator")
	}
	publicKey, publicErr := hex.DecodeString(encoded[:separator])
	if publicErr != nil {
		return SignatureSlot{}, publicErr
	}
	signature, signatureErr := hex.DecodeString(encoded[separator+1:])
	if signatureErr != nil {
		return SignatureSlot{}, signatureErr
	}
	if len(publicKey) != ed25519.PublicKeySize {
		return SignatureSlot{}, fmt.Errorf("unexpected public key length %d", len(publicKey))
	}
	if len(signature) != ed25519.SignatureSize {
		return SignatureSlot{}, fmt.Errorf("unexpected signature length %d", len(signature))
	}
	return SignatureSlot{PublicKey: publicKey, Signature: signature}, nil
}
