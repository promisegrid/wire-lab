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
// Intent: Keep POC15 aligned with grid([42(pCID), ...protocol-defined-slots])
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
	ProtocolCID        ProtocolCID
	Payload            []byte
	Proof              Proof
	ParentExactSHA256s []string
}

// NewEnvelope builds a legacy map payload and signs the envelope's two-slot
// signable view. New or reworked pCIDs should prefer NewEnvelopeFromPayload with
// pCID-owned CBOR arrays rather than extending the legacy field-map scaffold.
// Intent: Preserve existing POC15 flows while `DI-vipih` migrates new protocol
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
	return NewEnvelopeFromPayloadWithParents(protocolCID, payloadBytes, nil, signer)
}

// NewEnvelopeFromPayloadWithParents signs an already-encoded pCID-owned payload
// with optional envelope-parent links.
// Intent: DI-kohuj moves route_v1 parent-link pressure into ordinary traffic
// while keeping parent slots pCID-owned rather than universal. Source: DI-kohuj
func NewEnvelopeFromPayloadWithParents(protocolCID ProtocolCID, payloadBytes []byte, parentExactSHA256s []string, signer string) (Envelope, error) {
	copiedPayload := make([]byte, len(payloadBytes))
	copy(copiedPayload, payloadBytes)
	copiedParents := make([]string, len(parentExactSHA256s))
	copy(copiedParents, parentExactSHA256s)
	envelope := Envelope{ProtocolCID: protocolCID, Payload: copiedPayload, ParentExactSHA256s: copiedParents}
	proof, proofErr := SignEnvelope(envelope, signer)
	if proofErr != nil {
		return Envelope{}, proofErr
	}
	envelope.Proof = proof
	return envelope, nil
}

// SignableBytes serializes the pCID-defined signable view for proof generation.
// Intent: Parent slots must be covered by the durable envelope proof when a pCID
// chooses to place parent links outside the payload. Source: DI-kohuj
func (envelope Envelope) SignableBytes() ([]byte, error) {
	payloadSlot, payloadSlotErr := ByteStringGridSlot(envelope.Payload)
	if payloadSlotErr != nil {
		return nil, payloadSlotErr
	}
	if len(envelope.ParentExactSHA256s) == 0 {
		return EncodeGridMessage(envelope.ProtocolCID, payloadSlot)
	}
	parentSlot, parentSlotErr := parentExactSHA256GridSlot(envelope.ParentExactSHA256s)
	if parentSlotErr != nil {
		return nil, parentSlotErr
	}
	return EncodeGridMessage(envelope.ProtocolCID, parentSlot, payloadSlot)
}

func parentExactSHA256GridSlot(parentExactSHA256s []string) (GridSlot, error) {
	parentCIDs := make([][]byte, 0, len(parentExactSHA256s))
	for _, parentExactSHA256 := range parentExactSHA256s {
		parentDigest, decodeErr := hex.DecodeString(parentExactSHA256)
		if decodeErr != nil {
			return GridSlot{}, decodeErr
		}
		if len(parentDigest) != sha256.Size {
			return GridSlot{}, fmt.Errorf("parent exact sha256 must be %d bytes, got %d", sha256.Size, len(parentDigest))
		}
		cidBytes := make([]byte, 0, 36)
		cidBytes = append(cidBytes, 0x01, 0x55, 0x12, 0x20)
		cidBytes = append(cidBytes, parentDigest...)
		parentCIDs = append(parentCIDs, cidBytes)
	}
	parentArrayBytes, parentArrayErr := EncodeTag42LinkArray(parentCIDs...)
	if parentArrayErr != nil {
		return GridSlot{}, parentArrayErr
	}
	return RawCBORGridSlot(parentArrayBytes), nil
}

func decodeParentExactSHA256s(slot GridSlot) ([]string, error) {
	parentCIDs, parentErr := DecodeTag42LinkArray(slot.RawCBOR)
	if parentErr != nil {
		return nil, parentErr
	}
	parentExactSHA256s := make([]string, 0, len(parentCIDs))
	for _, parentCID := range parentCIDs {
		parentExactSHA256, ok := ExactSHA256FromRawCIDV1(parentCID)
		if !ok {
			return nil, fmt.Errorf("parent link must be CIDv1 raw sha2-256")
		}
		parentExactSHA256s = append(parentExactSHA256s, parentExactSHA256)
	}
	return parentExactSHA256s, nil
}

func proofBytesForEnvelope(envelope Envelope) ([]byte, error) {
	return MarshalStringMap(map[string]string{
		"signer":     envelope.Proof.Signer,
		"public_key": hex.EncodeToString(envelope.Proof.PublicKey),
		"signature":  hex.EncodeToString(envelope.Proof.Signature),
	})
}

func proofFromBytes(proofBytes []byte) (Proof, error) {
	proofFields, fieldsErr := UnmarshalStringMap(proofBytes)
	if fieldsErr != nil {
		return Proof{}, fieldsErr
	}
	publicKey, publicErr := hex.DecodeString(proofFields["public_key"])
	if publicErr != nil {
		return Proof{}, publicErr
	}
	signature, signatureErr := hex.DecodeString(proofFields["signature"])
	if signatureErr != nil {
		return Proof{}, signatureErr
	}
	return Proof{
		Signer:    proofFields["signer"],
		PublicKey: publicKey,
		Signature: signature,
	}, nil
}

func gridByteStringSlot(slot GridSlot, name string) ([]byte, error) {
	value, ok := slot.ByteString()
	if !ok {
		return nil, fmt.Errorf("%s must be a byte string", name)
	}
	return value, nil
}

func envelopeFromGridMessage(gridMessage GridMessage) (Envelope, error) {
	switch len(gridMessage.Slots) {
	case 2:
		payloadBytes, payloadErr := gridByteStringSlot(gridMessage.Slots[0], "payload slot")
		if payloadErr != nil {
			return Envelope{}, payloadErr
		}
		proofBytes, proofErr := gridByteStringSlot(gridMessage.Slots[1], "proof slot")
		if proofErr != nil {
			return Envelope{}, proofErr
		}
		proof, decodeErr := proofFromBytes(proofBytes)
		if decodeErr != nil {
			return Envelope{}, decodeErr
		}
		return Envelope{ProtocolCID: gridMessage.ProtocolCID, Payload: payloadBytes, Proof: proof}, nil
	case 3:
		parentExactSHA256s, parentErr := decodeParentExactSHA256s(gridMessage.Slots[0])
		if parentErr != nil {
			return Envelope{}, parentErr
		}
		payloadBytes, payloadErr := gridByteStringSlot(gridMessage.Slots[1], "payload slot")
		if payloadErr != nil {
			return Envelope{}, payloadErr
		}
		proofBytes, proofErr := gridByteStringSlot(gridMessage.Slots[2], "proof slot")
		if proofErr != nil {
			return Envelope{}, proofErr
		}
		proof, decodeErr := proofFromBytes(proofBytes)
		if decodeErr != nil {
			return Envelope{}, decodeErr
		}
		return Envelope{ProtocolCID: gridMessage.ProtocolCID, Payload: payloadBytes, Proof: proof, ParentExactSHA256s: parentExactSHA256s}, nil
	default:
		return Envelope{}, fmt.Errorf("poc15 envelope must have payload/proof or parents/payload/proof slots, got %d protocol-defined slots", len(gridMessage.Slots))
	}
}

// Bytes serializes the complete signed envelope using the slot shape selected by
// this envelope's pCID-owned parent-link state.
// Intent: Preserve the common grid([42(pCID), payload, proof]) shape for most
// traffic while allowing selected pCIDs to exercise
// grid([42(pCID), parents, payload, proof]) in normal POC15 traffic. Source:
// DI-kohuj
func (envelope Envelope) Bytes() ([]byte, error) {
	proofBytes, proofErr := proofBytesForEnvelope(envelope)
	if proofErr != nil {
		return nil, proofErr
	}
	payloadSlot, payloadSlotErr := ByteStringGridSlot(envelope.Payload)
	if payloadSlotErr != nil {
		return nil, payloadSlotErr
	}
	proofSlot, proofSlotErr := ByteStringGridSlot(proofBytes)
	if proofSlotErr != nil {
		return nil, proofSlotErr
	}
	if len(envelope.ParentExactSHA256s) == 0 {
		return EncodeGridMessage(envelope.ProtocolCID, payloadSlot, proofSlot)
	}
	parentSlot, parentSlotErr := parentExactSHA256GridSlot(envelope.ParentExactSHA256s)
	if parentSlotErr != nil {
		return nil, parentSlotErr
	}
	return EncodeGridMessage(envelope.ProtocolCID, parentSlot, payloadSlot, proofSlot)
}

// ParseEnvelope parses the POC15 signed envelope variants currently exercised by
// normal traffic: grid([42(pCID), payload, proof]) and route_v1's
// grid([42(pCID), parents, payload, proof]).
// Intent: Keep parsing strict for the current executable pCID-owned slot shapes
// so route multiarity is real without accepting arbitrary envelope layouts.
// Source: DI-kohuj
func ParseEnvelope(envelopeBytes []byte) (Envelope, error) {
	gridMessage, gridErr := ParseGridMessage(envelopeBytes)
	if gridErr != nil {
		return Envelope{}, gridErr
	}
	return envelopeFromGridMessage(gridMessage)
}

// PayloadFields decodes legacy map payloads and pCID-owned array payloads into
// local compatibility fields whose routing fields the POC15 kernel still needs.
// Intent: Kernel routing still needs promiser/promisee during the migration away
// from field maps, but the envelope layer should not make field maps the target
// protocol pattern. Source: DI-vipih; DI-dirat
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
	postalScaleFields, postalScaleErr := PostalScalePayloadFields(envelope.Payload)
	if postalScaleErr == nil {
		return postalScaleFields, nil
	}
	upsLabelFields, upsLabelErr := UPSLabelPayloadFields(envelope.Payload)
	if upsLabelErr == nil {
		return upsLabelFields, nil
	}
	accountingFields, accountingErr := AccountingPayloadFields(envelope.Payload)
	if accountingErr == nil {
		return accountingFields, nil
	}
	printerPortFields, printerPortErr := PrinterPortPayloadFields(envelope.Payload)
	if printerPortErr == nil {
		return printerPortFields, nil
	}
	kernelReceiveFields, kernelReceiveErr := KernelReceivePayloadFields(envelope.Payload)
	if kernelReceiveErr == nil {
		return kernelReceiveFields, nil
	}
	routeFields, routeErr := RoutePayloadFields(envelope.Payload)
	if routeErr == nil {
		return routeFields, nil
	}
	relationshipFields, relationshipErr := RelationshipPayloadFields(envelope.Payload)
	if relationshipErr == nil {
		return relationshipFields, nil
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
	seed := sha256.Sum256([]byte("poc15 protocol signer: " + seedText))
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
