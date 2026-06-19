package runtime

import (
	"fmt"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/decision"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/pcid"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/production"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
)

// MessageShapeSpecimen is one exact CBOR grid message emitted only for operator
// review and analyzer coverage.
// Intent: These artifacts test pCID-owned slot-vector alternatives without
// making normal app/kernel traffic depend on every experimental shape. Source:
// DI-mosat
type MessageShapeSpecimen struct {
	Name               string
	ProtocolName       string
	EnvelopeBytes      []byte
	ParentExactSHA256  string
	ParentLinkLocation string
}

// runMessageShapeSpecimenWorkflow emits exact raw messages for the POC16
// multiarity, parent-link, and COSE goals.
// Intent: Alice promises only to publish local specimens into the observer-only
// artifact stream; no peer is commanded to accept these shapes as app traffic.
// Source: DI-mosat
func (node *Node) runMessageShapeSpecimenWorkflow() error {
	specimens, err := node.messageShapeSpecimens()
	if err != nil {
		return err
	}
	for _, specimen := range specimens {
		fields := map[string]string{
			"promise_about": "message_shape_specimen",
		}
		if specimen.ParentExactSHA256 != "" {
			fields["parent_exact_sha256"] = specimen.ParentExactSHA256
			fields["parent_link_location"] = specimen.ParentLinkLocation
		}
		node.emitMessageArtifact("shape_specimen", specimen.Name, specimen.ProtocolName, specimen.EnvelopeBytes, fields)
		node.record("message_shape_specimen_emitted", "kept", specimen.Name, "pcid="+specimen.ProtocolName+" shape="+specimen.Name+" exact_sha256="+protocol.HashExactBytes(specimen.EnvelopeBytes))
	}
	return nil
}

func (node *Node) messageShapeSpecimens() ([]MessageShapeSpecimen, error) {
	payloadBytes, payloadErr := protocol.MarshalRelationshipPayloadFields(map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "operator",
		"turn":          "startup",
		"promise":       "Alice promises this message-shape specimen is local POC16 coverage, not a command to any peer.",
		"reason":        "POC16 must compare pCID-owned slot vectors by exact bytes",
		"promise_about": "message_shape_specimen",
		"shape_scope":   "operator_review",
	})
	if payloadErr != nil {
		return nil, payloadErr
	}
	proofBytes, proofErr := protocol.MarshalStringMap(map[string]string{
		"signer":    node.Agent.Name,
		"proof":     "native-poc16-specimen-proof",
		"sign_view": "pCID-defined specimen view",
	})
	if proofErr != nil {
		return nil, proofErr
	}
	nativeEnvelope, nativeErr := protocol.NewEnvelopeFromPayload(node.Protocols.MustCID(pcid.MessageShapeNativeProofV1), payloadBytes, node.Agent.Name)
	if nativeErr != nil {
		return nil, nativeErr
	}
	nativeBytes, nativeBytesErr := nativeEnvelope.Bytes()
	if nativeBytesErr != nil {
		return nil, nativeBytesErr
	}
	parentCIDBytes := protocol.RawCIDV1SHA256Bytes(nativeBytes)
	parentLinks, parentErr := protocol.EncodeTag42LinkArray(parentCIDBytes)
	if parentErr != nil {
		return nil, parentErr
	}
	payloadSlot, payloadSlotErr := protocol.ByteStringGridSlot(payloadBytes)
	if payloadSlotErr != nil {
		return nil, payloadSlotErr
	}
	proofSlot, proofSlotErr := protocol.ByteStringGridSlot(proofBytes)
	if proofSlotErr != nil {
		return nil, proofSlotErr
	}
	transportBytes, transportErr := protocol.EncodeGridMessage(node.Protocols.MustCID(pcid.MessageShapeTransportV1), payloadSlot)
	if transportErr != nil {
		return nil, transportErr
	}
	envelopeParentBytes, envelopeParentErr := protocol.EncodeGridMessage(node.Protocols.MustCID(pcid.MessageShapeEnvelopeParentsV1), protocol.RawCBORGridSlot(parentLinks), payloadSlot, proofSlot)
	if envelopeParentErr != nil {
		return nil, envelopeParentErr
	}
	payloadParentBytes, payloadParentErr := protocol.EncodeGridMessage(node.Protocols.MustCID(pcid.MessageShapePayloadParentsV1), payloadSlot, protocol.RawCBORGridSlot(parentLinks), proofSlot)
	if payloadParentErr != nil {
		return nil, payloadParentErr
	}
	cosePayload, cosePayloadErr := protocol.EncodeCOSESign1(payloadBytes, node.Agent.Name, false)
	if cosePayloadErr != nil {
		return nil, cosePayloadErr
	}
	if verifyErr := protocol.VerifyCOSESign1(cosePayload, nil, node.Agent.Name); verifyErr != nil {
		return nil, verifyErr
	}
	node.record("message_shape_cose_payload_verified", "kept", "operator", "pcid="+pcid.MessageShapeCOSEPayloadV1+" alg=EdDSA")
	cosePayloadBytes, cosePayloadMessageErr := protocol.EncodeGridMessage(node.Protocols.MustCID(pcid.MessageShapeCOSEPayloadV1), protocol.RawCBORGridSlot(cosePayload))
	if cosePayloadMessageErr != nil {
		return nil, cosePayloadMessageErr
	}
	coseProof, coseProofErr := protocol.EncodeCOSESign1(payloadBytes, node.Agent.Name, true)
	if coseProofErr != nil {
		return nil, coseProofErr
	}
	if verifyErr := protocol.VerifyCOSESign1(coseProof, payloadBytes, node.Agent.Name); verifyErr != nil {
		return nil, verifyErr
	}
	node.record("message_shape_cose_proof_verified", "kept", "operator", "pcid="+pcid.MessageShapeCOSEProofV1+" alg=EdDSA detached_payload=true")
	tamperedCOSEProof := append([]byte(nil), coseProof...)
	tamperedCOSEProof[len(tamperedCOSEProof)-1] ^= 0x01
	if verifyErr := protocol.VerifyCOSESign1(tamperedCOSEProof, payloadBytes, node.Agent.Name); verifyErr == nil {
		return nil, fmt.Errorf("tampered COSE proof unexpectedly verified")
	}
	node.record("message_shape_cose_tamper_rejected", "malformed", "operator", "pcid="+pcid.MessageShapeCOSEProofV1+" tampered detached COSE proof rejected")
	coseProofBytes, coseProofMessageErr := protocol.EncodeGridMessage(node.Protocols.MustCID(pcid.MessageShapeCOSEProofV1), payloadSlot, protocol.RawCBORGridSlot(coseProof))
	if coseProofMessageErr != nil {
		return nil, coseProofMessageErr
	}
	parentExactHash := protocol.HashExactBytes(nativeBytes)
	node.record("message_shape_multiarity_specimens_promised", "kept", "operator", "pcid="+pcid.MessageShapeTransportV1+" pCID owns arity and slot meaning")
	node.record("message_shape_parent_link_specimens_promised", "kept", "operator", "pcid="+pcid.MessageShapeEnvelopeParentsV1+" parent_exact_sha256="+parentExactHash)
	node.record("message_shape_cose_specimens_promised", "kept", "operator", "pcid="+pcid.MessageShapeCOSEPayloadV1+" "+pcid.MessageShapeCOSEProofV1)
	node.record("message_shape_native_proof_specimen_emitted", "kept", "operator", "pcid="+pcid.MessageShapeNativeProofV1+" proof_style=native")
	node.record("message_shape_transport_specimen_emitted", "kept", "operator", "pcid="+pcid.MessageShapeTransportV1+" slot_count=2")
	node.record("message_shape_envelope_parent_specimen_emitted", "kept", "operator", "pcid="+pcid.MessageShapeEnvelopeParentsV1+" parent_slot=1")
	node.record("message_shape_payload_parent_specimen_emitted", "kept", "operator", "pcid="+pcid.MessageShapePayloadParentsV1+" parent_slot=2")
	node.record("message_shape_cose_payload_specimen_emitted", "kept", "operator", "pcid="+pcid.MessageShapeCOSEPayloadV1+" cose_slot=1")
	node.record("message_shape_cose_proof_specimen_emitted", "kept", "operator", "pcid="+pcid.MessageShapeCOSEProofV1+" cose_slot=2")
	return []MessageShapeSpecimen{
		{Name: "transport_payload_only", ProtocolName: pcid.MessageShapeTransportV1, EnvelopeBytes: transportBytes},
		{Name: "native_payload_proof", ProtocolName: pcid.MessageShapeNativeProofV1, EnvelopeBytes: nativeBytes},
		{Name: "envelope_parents_payload_proof", ProtocolName: pcid.MessageShapeEnvelopeParentsV1, EnvelopeBytes: envelopeParentBytes, ParentExactSHA256: parentExactHash, ParentLinkLocation: "envelope"},
		{Name: "payload_parents_proof", ProtocolName: pcid.MessageShapePayloadParentsV1, EnvelopeBytes: payloadParentBytes, ParentExactSHA256: parentExactHash, ParentLinkLocation: "payload"},
		{Name: "cose_as_payload", ProtocolName: pcid.MessageShapeCOSEPayloadV1, EnvelopeBytes: cosePayloadBytes},
		{Name: "cose_as_proof", ProtocolName: pcid.MessageShapeCOSEProofV1, EnvelopeBytes: coseProofBytes},
	}, nil
}

func (node *Node) recordMessageShapeSpecimenCoverage() {
	node.record("kernel_role_profile_recorded", "kept", "operator", "transport, app-interface, routing, local-resource, and event roles may be split or collapsed by runtime profile")
	node.record("message_shape_specimen_scope_recorded", "kept", "operator", "specimens are run-local raw artifacts, not receive promises or production APIs")
	node.record("message_shape_sample_workload_recorded", "kept", "operator", "sample_content_cid="+production.ContentCID(production.SampleContentBytes()))
	node.record("transport_proof_comparison_recorded", "kept", "operator", "transport/session signatures are hop-local promises; envelope proofs remain durable object-level promise records for retained message bytes")
}
