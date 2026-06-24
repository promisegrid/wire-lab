package protocol

import (
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
	ParentCIDs  []string
}

// NewEnvelope builds an older generic map payload and signs the envelope's two-slot
// signable view. New or reworked pCIDs should prefer NewEnvelopeFromPayload with
// pCID-owned CBOR arrays rather than extending the older generic map scaffold.
// Intent: Preserve existing POC16 flows while `DI-vipih` migrates new protocol
// work toward pCID-owned array payloads. Source: DI-vipih
func NewEnvelope(protocolCID ProtocolCID, fields map[string]string, signer string) (Envelope, error) {
	return NewEnvelopeWithParents(protocolCID, fields, nil, signer)
}

// NewEnvelopeWithParents builds an older generic map payload and signs it with
// parent links over the same signable view used by pCID-owned array payloads.
// Intent: POC16 ACKs must be able to parent-link the request message CID even
// when a not-yet-migrated pCID still uses the generic map scaffold. Source:
// DI-vopab
func NewEnvelopeWithParents(protocolCID ProtocolCID, fields map[string]string, parentCIDs []string, signer string) (Envelope, error) {
	payloadBytes, marshalErr := MarshalStringMap(fields)
	if marshalErr != nil {
		return Envelope{}, marshalErr
	}
	return NewEnvelopeFromPayloadWithParents(protocolCID, payloadBytes, parentCIDs, signer)
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
func NewEnvelopeFromPayloadWithParents(protocolCID ProtocolCID, payloadBytes []byte, parentCIDs []string, signer string) (Envelope, error) {
	copiedPayload := make([]byte, len(payloadBytes))
	copy(copiedPayload, payloadBytes)
	copiedParents := make([]string, len(parentCIDs))
	copy(copiedParents, parentCIDs)
	envelope := Envelope{ProtocolCID: protocolCID, Payload: copiedPayload, ParentCIDs: copiedParents}
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
	if len(envelope.ParentCIDs) == 0 {
		return EncodeGridMessage(envelope.ProtocolCID, payloadSlot)
	}
	parentSlot, parentSlotErr := parentCIDGridSlot(envelope.ParentCIDs)
	if parentSlotErr != nil {
		return nil, parentSlotErr
	}
	return EncodeGridMessage(envelope.ProtocolCID, parentSlot, payloadSlot)
}

func parentCIDGridSlot(parentCIDTexts []string) (GridSlot, error) {
	parentCIDs := make([][]byte, 0, len(parentCIDTexts))
	for _, parentCIDText := range parentCIDTexts {
		parsedCID, parseErr := ParseCIDText(parentCIDText)
		if parseErr != nil {
			return GridSlot{}, parseErr
		}
		cidBytes := parsedCID.Bytes()
		parentCIDs = append(parentCIDs, cidBytes)
	}
	parentArrayBytes, parentArrayErr := EncodeTag42LinkArray(parentCIDs...)
	if parentArrayErr != nil {
		return GridSlot{}, parentArrayErr
	}
	return RawCBORGridSlot(parentArrayBytes), nil
}

func decodeParentCIDs(slot GridSlot) ([]string, error) {
	parentCIDs, parentErr := DecodeTag42LinkArray(slot.RawCBOR)
	if parentErr != nil {
		return nil, parentErr
	}
	parentCIDTexts := make([]string, 0, len(parentCIDs))
	for _, parentCID := range parentCIDs {
		parentCIDText, textErr := CIDTextFromBytes(parentCID)
		if textErr != nil {
			return nil, fmt.Errorf("parent link must be CIDv1 raw sha2-256: %w", textErr)
		}
		parentCIDTexts = append(parentCIDTexts, parentCIDText)
	}
	return parentCIDTexts, nil
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
		parentCIDs, parentErr := decodeParentCIDs(gridMessage.Slots[0])
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
		return Envelope{ProtocolCID: gridMessage.ProtocolCID, Payload: payloadBytes, Proof: proof, ParentCIDs: parentCIDs}, nil
	default:
		return Envelope{}, fmt.Errorf("poc16 envelope must have payload/proof or parents/payload/proof slots, got %d protocol-defined slots", len(gridMessage.Slots))
	}
}

// Bytes serializes the complete signed envelope using the slot shape selected by
// this envelope's pCID-owned parent-link state.
// Intent: Preserve the common grid([42(pCID), payload, proof]) shape for most
// traffic while allowing selected pCIDs to exercise
// grid([42(pCID), parents, payload, proof]) in normal POC16 traffic. Source:
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
	if len(envelope.ParentCIDs) == 0 {
		return EncodeGridMessage(envelope.ProtocolCID, payloadSlot, proofSlot)
	}
	parentSlot, parentSlotErr := parentCIDGridSlot(envelope.ParentCIDs)
	if parentSlotErr != nil {
		return nil, parentSlotErr
	}
	return EncodeGridMessage(envelope.ProtocolCID, parentSlot, payloadSlot, proofSlot)
}

// ParseEnvelope parses the POC16 signed envelope variants currently exercised by
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

// PayloadFields decodes older map payloads and pCID-owned array payloads into
// local compatibility fields whose routing fields the POC16 kernel still needs.
// Intent: Kernel routing still needs promiser/promisee during the migration away
// from generic map payloads, but the envelope layer should not make generic maps
// the target protocol pattern. Source: DI-vipih; DI-dirat
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
	seed := sha256.Sum256([]byte("poc16 protocol signer: " + seedText))
	return ed25519.NewKeyFromSeed(seed[:])
}

// DeterministicPublicKey returns the POC public key for a named local agent.
func DeterministicPublicKey(seedText string) ed25519.PublicKey {
	return DeterministicPrivateKey(seedText).Public().(ed25519.PublicKey)
}
