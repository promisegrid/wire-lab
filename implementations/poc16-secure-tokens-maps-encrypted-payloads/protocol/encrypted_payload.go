package protocol

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

const encryptedPayloadAlgorithm = "AES-256-GCM"
const encryptedPayloadKeyProfile = "poc16-local-derived-demo-key"

// EncryptedPayload is POC16's protocol-owned encrypted slot-1 payload profile.
//
// Intent: POC16 needs real ciphertext and authentication failure behavior without
// pretending to solve production key agreement. The key profile is explicitly
// POC-local; later pCIDs must replace it with relationship-specific key promises.
// Source: DI-vulit
type EncryptedPayload struct {
	Sender           string
	Recipient        string
	Context          string
	Algorithm        string
	KeyProfile       string
	NonceBase64      string
	CiphertextBase64 string
}

// EncryptPayloadForRecipient seals plaintext for one recipient using a POC-local
// derived key and AEAD additional data tied to the sender, recipient, and context.
func EncryptPayloadForRecipient(sender, recipient, context string, plaintext []byte) (EncryptedPayload, error) {
	if sender == "" || recipient == "" || context == "" {
		return EncryptedPayload{}, fmt.Errorf("encrypted payload needs sender, recipient, and context")
	}
	block, blockErr := aes.NewCipher(derivePOC16EncryptionKey(sender, recipient, context))
	if blockErr != nil {
		return EncryptedPayload{}, blockErr
	}
	aead, aeadErr := cipher.NewGCM(block)
	if aeadErr != nil {
		return EncryptedPayload{}, aeadErr
	}
	nonce := make([]byte, aead.NonceSize())
	if _, readErr := io.ReadFull(rand.Reader, nonce); readErr != nil {
		return EncryptedPayload{}, readErr
	}
	ciphertext := aead.Seal(nil, nonce, plaintext, encryptedPayloadAAD(sender, recipient, context))
	return EncryptedPayload{
		Sender:           sender,
		Recipient:        recipient,
		Context:          context,
		Algorithm:        encryptedPayloadAlgorithm,
		KeyProfile:       encryptedPayloadKeyProfile,
		NonceBase64:      base64.StdEncoding.EncodeToString(nonce),
		CiphertextBase64: base64.StdEncoding.EncodeToString(ciphertext),
	}, nil
}

// DecryptPayloadForRecipient opens a POC16 encrypted payload only for the named
// recipient and returns AEAD errors for wrong recipients or tampered bytes.
func DecryptPayloadForRecipient(payload EncryptedPayload, recipient string) ([]byte, error) {
	if payload.Algorithm != encryptedPayloadAlgorithm {
		return nil, fmt.Errorf("encrypted payload algorithm=%s, want %s", payload.Algorithm, encryptedPayloadAlgorithm)
	}
	if payload.KeyProfile != encryptedPayloadKeyProfile {
		return nil, fmt.Errorf("encrypted payload key profile=%s, want %s", payload.KeyProfile, encryptedPayloadKeyProfile)
	}
	if payload.Recipient != recipient {
		return nil, fmt.Errorf("encrypted payload recipient=%s, want %s", payload.Recipient, recipient)
	}
	nonce, nonceErr := base64.StdEncoding.DecodeString(payload.NonceBase64)
	if nonceErr != nil {
		return nil, nonceErr
	}
	ciphertext, ciphertextErr := base64.StdEncoding.DecodeString(payload.CiphertextBase64)
	if ciphertextErr != nil {
		return nil, ciphertextErr
	}
	block, blockErr := aes.NewCipher(derivePOC16EncryptionKey(payload.Sender, payload.Recipient, payload.Context))
	if blockErr != nil {
		return nil, blockErr
	}
	aead, aeadErr := cipher.NewGCM(block)
	if aeadErr != nil {
		return nil, aeadErr
	}
	return aead.Open(nil, nonce, ciphertext, encryptedPayloadAAD(payload.Sender, payload.Recipient, payload.Context))
}

// MarshalCBOR encodes the encrypted payload as a CBOR map because POC16 permits
// self-documenting map payloads when the pCID spec explicitly chooses them.
func (payload EncryptedPayload) MarshalCBOR() ([]byte, error) {
	return MarshalStringMap(payload.StringFields())
}

// EncryptedPayloadFromCBOR decodes the POC16 encrypted-payload map profile.
func EncryptedPayloadFromCBOR(payloadBytes []byte) (EncryptedPayload, error) {
	fields, fieldsErr := UnmarshalStringMap(payloadBytes)
	if fieldsErr != nil {
		return EncryptedPayload{}, fieldsErr
	}
	return EncryptedPayload{
		Sender:           fields["sender"],
		Recipient:        fields["recipient"],
		Context:          fields["context"],
		Algorithm:        fields["algorithm"],
		KeyProfile:       fields["key_profile"],
		NonceBase64:      fields["nonce_b64"],
		CiphertextBase64: fields["ciphertext_b64"],
	}, nil
}

// StringFields exposes non-secret encrypted-payload metadata for local event
// records and raw-message artifact indexing.
func (payload EncryptedPayload) StringFields() map[string]string {
	return map[string]string{
		"algorithm":      payload.Algorithm,
		"ciphertext_b64": payload.CiphertextBase64,
		"context":        payload.Context,
		"key_profile":    payload.KeyProfile,
		"nonce_b64":      payload.NonceBase64,
		"recipient":      payload.Recipient,
		"sender":         payload.Sender,
	}
}

func derivePOC16EncryptionKey(sender, recipient, context string) []byte {
	digest := sha256.Sum256([]byte("poc16-encrypted-payload-v1\x00" + sender + "\x00" + recipient + "\x00" + context))
	return digest[:]
}

func encryptedPayloadAAD(sender, recipient, context string) []byte {
	return []byte("poc16-encrypted-payload-v1\x00" + sender + "\x00" + recipient + "\x00" + context)
}
