package lib

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
)

// SignatureSlot carries a POC public key and signature for exact signable bytes.
// Intent: Let the signed app witness bytes without claiming global trust.
// Source: DI-rarim.
type SignatureSlot struct {
	PublicKey []byte
	Signature []byte
}

// EncodeSignatureSlot serializes public key and signature bytes.
func EncodeSignatureSlot(slot SignatureSlot) []byte {
	return []byte(hex.EncodeToString(slot.PublicKey) + ":" + hex.EncodeToString(slot.Signature))
}

// DecodeSignatureSlot decodes public key and signature bytes.
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
