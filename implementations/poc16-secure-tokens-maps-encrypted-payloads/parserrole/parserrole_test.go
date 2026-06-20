package parserrole

import (
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/config"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/pcid"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
)

func TestFrameIsResponseRequiresParentLink(t *testing.T) {
	// Intent: Parser roles must not swallow fresh pCID-owned promises merely
	// because their payload discusses an outcome; only parent-linked ACK frames
	// are session responses. Source: DI-gazin
	role := New(config.Config{}, "test")
	fields := map[string]string{
		"act":           "promise",
		"from":          "alice",
		"to":            "bob",
		"promise_about": "relationship_repair",
		"outcome":       "non_commitment",
		"promise":       "I promise to describe my local repair outcome.",
		"reason":        "fresh promise with outcome vocabulary",
	}
	payloadBytes, arrayPayload, payloadErr := protocol.MarshalKnownArrayPayload(pcid.RelationshipV1, fields)
	if payloadErr != nil {
		t.Fatalf("marshal relationship payload: %v", payloadErr)
	}
	if !arrayPayload {
		t.Fatalf("relationship_v1 test payload was not encoded as an array")
	}
	freshEnvelope, freshErr := protocol.NewEnvelopeFromPayload(role.Protocols.MustCID(pcid.RelationshipV1), payloadBytes, "alice")
	if freshErr != nil {
		t.Fatalf("build fresh envelope: %v", freshErr)
	}
	freshBytes, freshBytesErr := freshEnvelope.Bytes()
	if freshBytesErr != nil {
		t.Fatalf("fresh envelope bytes: %v", freshBytesErr)
	}
	freshIsResponse, freshResponseErr := role.frameIsResponse(freshBytes)
	if freshResponseErr != nil {
		t.Fatalf("classify fresh envelope: %v", freshResponseErr)
	}
	if freshIsResponse {
		t.Fatalf("fresh outcome-bearing promise was misclassified as a response")
	}
	parentHash := protocol.HashExactBytes([]byte("parent promise"))
	responseEnvelope, responseErr := protocol.NewEnvelopeFromPayloadWithParents(role.Protocols.MustCID(pcid.RelationshipV1), payloadBytes, []string{parentHash}, "alice")
	if responseErr != nil {
		t.Fatalf("build response envelope: %v", responseErr)
	}
	responseBytes, responseBytesErr := responseEnvelope.Bytes()
	if responseBytesErr != nil {
		t.Fatalf("response envelope bytes: %v", responseBytesErr)
	}
	responseIsResponse, responseClassifyErr := role.frameIsResponse(responseBytes)
	if responseClassifyErr != nil {
		t.Fatalf("classify response envelope: %v", responseClassifyErr)
	}
	if !responseIsResponse {
		t.Fatalf("parent-linked outcome payload was not classified as a response")
	}
}
