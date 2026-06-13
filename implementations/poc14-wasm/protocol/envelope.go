package protocol

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	// gridTag is the CBOR tag for the PromiseGrid grid wrapper: ASCII "grid"
	// interpreted as 0x67726964 / decimal 1735551332.
	gridTag        = uint64(0x67726964)
	dagCBORLinkTag = uint64(42)
)

// ProtocolCID identifies the protocol spec document that defines every slot
// after slot 0. It is not a payload CID.
// Intent: Keep POC14 aligned with grid([42(pCID), ...protocol-defined-slots])
// while testing LLM-directed behavior above the envelope layer. Source: DI-timah
type ProtocolCID struct {
	cidBytes []byte
	digest   [32]byte
}

// NewProtocolCID hashes spec bytes into the POC CIDv1 raw sha2-256 form.
func NewProtocolCID(specBytes []byte) ProtocolCID {
	digest := sha256.Sum256(specBytes)
	cidBytes := make([]byte, 0, 36)
	cidBytes = append(cidBytes, 0x01, 0x55, 0x12, 0x20)
	cidBytes = append(cidBytes, digest[:]...)
	return ProtocolCID{cidBytes: cidBytes, digest: digest}
}

// NewProtocolCIDFromBytes rebuilds a POC pCID from binary CIDv1 bytes.
func NewProtocolCIDFromBytes(cidBytes []byte) ProtocolCID {
	copiedBytes := make([]byte, len(cidBytes))
	copy(copiedBytes, cidBytes)
	var digest [32]byte
	if len(copiedBytes) >= 36 {
		copy(digest[:], copiedBytes[len(copiedBytes)-32:])
	}
	return ProtocolCID{cidBytes: copiedBytes, digest: digest}
}

// Bytes returns binary CIDv1 bytes without the tag-42 sentinel.
func (protocolCID ProtocolCID) Bytes() []byte {
	cidBytes := make([]byte, len(protocolCID.cidBytes))
	copy(cidBytes, protocolCID.cidBytes)
	return cidBytes
}

// Tag42Bytes returns the DAG-CBOR tag-42 byte string payload.
func (protocolCID ProtocolCID) Tag42Bytes() []byte {
	tagBytes := make([]byte, 0, len(protocolCID.cidBytes)+1)
	tagBytes = append(tagBytes, 0x00)
	tagBytes = append(tagBytes, protocolCID.cidBytes...)
	return tagBytes
}

// String renders a stable debug form for event records.
func (protocolCID ProtocolCID) String() string {
	return "cidv1-raw-sha2-256:" + hex.EncodeToString(protocolCID.digest[:])
}

// Equal reports whether two pCIDs name the same protocol bytes.
func (protocolCID ProtocolCID) Equal(other ProtocolCID) bool {
	return bytes.Equal(protocolCID.cidBytes, other.cidBytes)
}

// Proof is an Ed25519 proof over the envelope's pCID-defined signable view.
type Proof struct {
	Signer    string
	PublicKey []byte
	Signature []byte
}

// Envelope represents grid([42(pCID), payload, proof]).
// Intent: The LLM can influence payload semantics, but Go owns the CBOR grid
// envelope and exact-byte signature mechanics. Source: DI-timah
type Envelope struct {
	ProtocolCID ProtocolCID
	Payload     []byte
	Proof       Proof
}

// NewEnvelope builds a legacy map payload and signs the envelope's two-slot
// signable view. New or reworked pCIDs should prefer NewEnvelopeFromPayload with
// pCID-owned CBOR arrays rather than extending the legacy field-map scaffold.
// Intent: Preserve existing POC14 flows while `DI-vipih` migrates new protocol
// work toward pCID-owned array payloads. Source: DI-vipih
func NewEnvelope(protocolCID ProtocolCID, fields map[string]string, signer string) (Envelope, error) {
	payloadBytes, marshalErr := MarshalStringMap(fields)
	if marshalErr != nil {
		return Envelope{}, marshalErr
	}
	return NewEnvelopeFromPayload(protocolCID, payloadBytes, signer)
}

// NewEnvelopeFromPayload signs an already-encoded pCID-owned payload.
// Intent: The envelope layer must not impose one universal payload shape; the
// pCID spec owns the bytes in slot 1. Source: DI-vipih
func NewEnvelopeFromPayload(protocolCID ProtocolCID, payloadBytes []byte, signer string) (Envelope, error) {
	copiedPayload := make([]byte, len(payloadBytes))
	copy(copiedPayload, payloadBytes)
	envelope := Envelope{ProtocolCID: protocolCID, Payload: copiedPayload}
	proof, proofErr := SignEnvelope(envelope, signer)
	if proofErr != nil {
		return Envelope{}, proofErr
	}
	envelope.Proof = proof
	return envelope, nil
}

// SignableBytes serializes grid([42(pCID), payload]) for proof generation.
func (envelope Envelope) SignableBytes() ([]byte, error) {
	writer := &cborWriter{}
	if err := writer.writeTag(gridTag); err != nil {
		return nil, err
	}
	if err := writer.writeArrayHeader(2); err != nil {
		return nil, err
	}
	if err := writer.writeTag(dagCBORLinkTag); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(envelope.ProtocolCID.Tag42Bytes()); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(envelope.Payload); err != nil {
		return nil, err
	}
	return writer.buffer.Bytes(), nil
}

// Bytes serializes the complete signed envelope as grid([42(pCID), payload, proof]).
func (envelope Envelope) Bytes() ([]byte, error) {
	proofBytes, proofErr := MarshalStringMap(map[string]string{
		"signer":     envelope.Proof.Signer,
		"public_key": hex.EncodeToString(envelope.Proof.PublicKey),
		"signature":  hex.EncodeToString(envelope.Proof.Signature),
	})
	if proofErr != nil {
		return nil, proofErr
	}
	writer := &cborWriter{}
	if err := writer.writeTag(gridTag); err != nil {
		return nil, err
	}
	if err := writer.writeArrayHeader(3); err != nil {
		return nil, err
	}
	if err := writer.writeTag(dagCBORLinkTag); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(envelope.ProtocolCID.Tag42Bytes()); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(envelope.Payload); err != nil {
		return nil, err
	}
	if err := writer.writeBytes(proofBytes); err != nil {
		return nil, err
	}
	return writer.buffer.Bytes(), nil
}

// ParseEnvelope parses grid([42(pCID), payload, proof]).
func ParseEnvelope(envelopeBytes []byte) (Envelope, error) {
	reader := &cborReader{data: envelopeBytes}
	outerTag, outerTagErr := reader.readTypeAndLength(6)
	if outerTagErr != nil {
		return Envelope{}, outerTagErr
	}
	if outerTag != gridTag {
		return Envelope{}, fmt.Errorf("outer envelope must be grid tag %d, got tag %d", gridTag, outerTag)
	}
	arrayLength, arrayErr := reader.readTypeAndLength(4)
	if arrayErr != nil {
		return Envelope{}, arrayErr
	}
	if arrayLength != 3 {
		return Envelope{}, fmt.Errorf("poc14 envelope must have three slots, got %d", arrayLength)
	}
	tagNumber, tagErr := reader.readTypeAndLength(6)
	if tagErr != nil {
		return Envelope{}, tagErr
	}
	if tagNumber != dagCBORLinkTag {
		return Envelope{}, fmt.Errorf("slot 0 must be tag 42, got tag %d", tagNumber)
	}
	tagBytes, tagBytesErr := reader.readBytes()
	if tagBytesErr != nil {
		return Envelope{}, tagBytesErr
	}
	if len(tagBytes) < 2 || tagBytes[0] != 0x00 {
		return Envelope{}, fmt.Errorf("tag 42 payload must start with DAG-CBOR CID sentinel")
	}
	payloadBytes, payloadErr := reader.readBytes()
	if payloadErr != nil {
		return Envelope{}, payloadErr
	}
	proofBytes, proofErr := reader.readBytes()
	if proofErr != nil {
		return Envelope{}, proofErr
	}
	if reader.offset != len(reader.data) {
		return Envelope{}, fmt.Errorf("trailing cbor bytes in envelope: %d", len(reader.data)-reader.offset)
	}
	proofFields, fieldsErr := UnmarshalStringMap(proofBytes)
	if fieldsErr != nil {
		return Envelope{}, fieldsErr
	}
	publicKey, publicErr := hex.DecodeString(proofFields["public_key"])
	if publicErr != nil {
		return Envelope{}, publicErr
	}
	signature, signatureErr := hex.DecodeString(proofFields["signature"])
	if signatureErr != nil {
		return Envelope{}, signatureErr
	}
	return Envelope{
		ProtocolCID: NewProtocolCIDFromBytes(tagBytes[1:]),
		Payload:     payloadBytes,
		Proof: Proof{
			Signer:    proofFields["signer"],
			PublicKey: publicKey,
			Signature: signature,
		},
	}, nil
}

// PayloadFields decodes legacy map payloads and the small set of transitional
// pCID-owned array payloads whose routing fields the POC14 kernel still needs.
// Intent: Kernel routing still needs promiser/promisee during the migration away
// from field maps, but the envelope layer should not make field maps the target
// protocol pattern. Source: DI-vipih
func (envelope Envelope) PayloadFields() (map[string]string, error) {
	fields, fieldsErr := UnmarshalStringMap(envelope.Payload)
	if fieldsErr == nil {
		return fields, nil
	}
	identityKeyFields, identityKeyErr := IdentityKeyPayloadFields(envelope.Payload)
	if identityKeyErr == nil {
		return identityKeyFields, nil
	}
	casStorageFields, casStorageErr := CASStoragePayloadFields(envelope.Payload)
	if casStorageErr == nil {
		return casStorageFields, nil
	}
	cidComputeFields, cidComputeErr := CIDComputePayloadFields(envelope.Payload)
	if cidComputeErr == nil {
		return cidComputeFields, nil
	}
	return nil, fieldsErr
}

// SignEnvelope signs the envelope's pCID-defined signable view.
func SignEnvelope(envelope Envelope, signer string) (Proof, error) {
	signableBytes, signableErr := envelope.SignableBytes()
	if signableErr != nil {
		return Proof{}, signableErr
	}
	privateKey := DeterministicPrivateKey(signer)
	return Proof{
		Signer:    signer,
		PublicKey: DeterministicPublicKey(signer),
		Signature: ed25519.Sign(privateKey, signableBytes),
	}, nil
}

// VerifyEnvelope checks the Ed25519 proof over the pCID-defined signable view.
func VerifyEnvelope(envelope Envelope) error {
	if len(envelope.Proof.PublicKey) != ed25519.PublicKeySize {
		return fmt.Errorf("unexpected envelope public key length %d", len(envelope.Proof.PublicKey))
	}
	if len(envelope.Proof.Signature) != ed25519.SignatureSize {
		return fmt.Errorf("unexpected envelope signature length %d", len(envelope.Proof.Signature))
	}
	signableBytes, signableErr := envelope.SignableBytes()
	if signableErr != nil {
		return signableErr
	}
	if !ed25519.Verify(envelope.Proof.PublicKey, signableBytes, envelope.Proof.Signature) {
		return fmt.Errorf("envelope signature failed")
	}
	return nil
}

// DeterministicPrivateKey is a POC-only key source for reproducible signatures.
func DeterministicPrivateKey(seedText string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte("poc14 protocol signer: " + seedText))
	return ed25519.NewKeyFromSeed(seed[:])
}

// DeterministicPublicKey returns the POC public key for a named local agent.
func DeterministicPublicKey(seedText string) ed25519.PublicKey {
	return DeterministicPrivateKey(seedText).Public().(ed25519.PublicKey)
}

// HashExactBytes returns a sha256 hex digest for correlating local event records.
func HashExactBytes(exactBytes []byte) string {
	if len(exactBytes) == 0 {
		return ""
	}
	digest := sha256.Sum256(exactBytes)
	return hex.EncodeToString(digest[:])
}
