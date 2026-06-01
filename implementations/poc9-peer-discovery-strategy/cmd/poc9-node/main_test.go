package main

import (
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc9-peer-discovery-strategy/compute"
	"promisegrid.dev/wire-lab/implementations/poc9-peer-discovery-strategy/protocol"
	"promisegrid.dev/wire-lab/implementations/poc9-peer-discovery-strategy/storage"
	"promisegrid.dev/wire-lab/implementations/poc9-peer-discovery-strategy/token"
)

func TestPayloadKindsShareOneProtocolCID(t *testing.T) {
	// Intent: POC9 keeps discovery and economy messages under one pCID; Kind is a
	// protocol-owned payload variant, not a separate protocol selector. Source:
	// DI-sipuz
	node := newNode("alice")
	payloadKinds := []string{
		kindNeedAdvertisement,
		kindOfferPromise,
		kindCounterPromise,
		kindAcceptancePromise,
		kindTokenIssuePromise,
		kindTokenRedemptionPromise,
		kindOutcomeObservation,
		kindReferralPromise,
		kindIntroductionPromise,
		kindRoutePromise,
	}
	for _, payloadKind := range payloadKinds {
		message, messageErr := node.newWireMessage("bob", payloadKind, []string{"bob"}, map[string]string{"resource_kind": storage.Kind}, token.Token{})
		if messageErr != nil {
			t.Fatalf("new wire message %s: %v", payloadKind, messageErr)
		}
		envelopeBytes, envelopeErr := message.EnvelopeBytes()
		if envelopeErr != nil {
			t.Fatalf("envelope bytes %s: %v", payloadKind, envelopeErr)
		}
		envelope, parseErr := protocol.ParseEnvelope(envelopeBytes)
		if parseErr != nil {
			t.Fatalf("parse envelope %s: %v", payloadKind, parseErr)
		}
		if !envelope.ProtocolCID.Equal(appProtocolCID) {
			t.Fatalf("payload kind %s used unexpected pCID %s", payloadKind, envelope.ProtocolCID)
		}
		fields, fieldsErr := envelope.PayloadFields()
		if fieldsErr != nil {
			t.Fatalf("payload fields %s: %v", payloadKind, fieldsErr)
		}
		if fields["kind"] != payloadKind {
			t.Fatalf("payload kind field = %q, want %q", fields["kind"], payloadKind)
		}
	}
}

func TestSparseMeshRouteAndNeighbors(t *testing.T) {
	// Intent: The sparse mesh is deterministic and connected, while direct TCP is
	// still limited to mesh neighbors. Source: DI-sipuz
	mesh := NewMesh(meshEdges)
	if !mesh.HasEdge("alice", "bob") || !mesh.HasEdge("frank", "mallory") {
		t.Fatalf("expected approved sparse mesh edges to exist")
	}
	if mesh.HasEdge("alice", "carol") {
		t.Fatalf("alice and carol should not be direct TCP neighbors")
	}
	route, routeErr := mesh.Route("alice", "mallory")
	if routeErr != nil {
		t.Fatalf("mesh route alice->mallory: %v", routeErr)
	}
	if len(route) < 2 {
		t.Fatalf("alice->mallory should require a non-neighbor path, got %v", route)
	}
}

func TestRouteRequiresLocalDiscovery(t *testing.T) {
	// Intent: Non-neighbor routes must be locally learned from route promises;
	// POC9 no longer uses POC8's static all-pairs route table. Source: DI-sipuz
	alice := newNode("alice")
	if _, routeErr := alice.selectRoute("carol"); routeErr == nil {
		t.Fatalf("alice should not know a non-neighbor route before route evidence")
	}
	routePromise := WireMessage{
		FromNode: "bob",
		ToNode:   "alice",
		Payload: map[string]string{
			"target_node":    "carol",
			"promised_route": "bob,carol",
			"route":          "bob",
		},
	}
	if _, handleErr := alice.handleRoutePromise(routePromise); handleErr != nil {
		t.Fatalf("handle route promise: %v", handleErr)
	}
	if !alice.hasCompleted("route:carol") {
		t.Fatalf("route promise should create local route evidence")
	}
	alice = newNode("alice")
	alice.learnRoute("carol", []string{"bob", "carol"}, "bob")
	route, routeErr := alice.selectRoute("carol")
	if routeErr != nil {
		t.Fatalf("alice learned route to carol: %v", routeErr)
	}
	if got := encodeRoute(route); got != "bob,carol" {
		t.Fatalf("route = %s, want bob,carol", got)
	}
	if !alice.hasCompleted("route:carol") {
		t.Fatalf("learned non-neighbor route should create local route evidence")
	}
}

func TestReferralDoesNotCreateTargetTrust(t *testing.T) {
	// Intent: Bob's referral about Carol is evidence about Bob's promise, not a
	// transfer of trust from Bob to Carol. Source: DI-sipuz
	alice := newNode("alice")
	message := WireMessage{FromNode: "bob", ToNode: "alice", Payload: map[string]string{"referred_peer": "carol", "resource_kind": compute.Kind, "resource_id": "fib-8"}}
	if _, handleErr := alice.handleReferralPromise(message); handleErr != nil {
		t.Fatalf("handle referral: %v", handleErr)
	}
	if alice.Wallet.Trust("carol") != 0 {
		t.Fatalf("referral should not create direct trust in Carol")
	}
	if !alice.hasCompleted("referral:carol:compute") {
		t.Fatalf("referral evidence should still be recorded for Alice's local strategy")
	}
}

func TestOrdinaryPublicPromiseGatesPrivateEscalation(t *testing.T) {
	// Intent: A probe is just an ordinary low-risk promise; private escalation is a
	// local strategy decision after public keep evidence, not a special wire kind.
	// Source: DI-sipuz
	alice := newNode("alice")
	privateOffer := WireMessage{FromNode: "bob", ToNode: "alice", Payload: map[string]string{"resource_kind": storage.Kind, "stage": "private"}}
	if alice.canAcceptOffer(privateOffer) {
		t.Fatalf("private offer before public evidence should be refused by local strategy")
	}
	alice.markCompleted("kept:bob:storage:public", "test public storage kept")
	if !alice.canAcceptOffer(privateOffer) {
		t.Fatalf("private offer after public evidence should be acceptable to local strategy")
	}
}

func TestMalformedTransportLowersRouteTrust(t *testing.T) {
	// Intent: Transport break evidence changes later route choices, but it remains
	// local evidence rather than a global ban or authority. Source: DI-sipuz
	carol := newNode("carol")
	for observationIndex := 0; observationIndex < 3; observationIndex++ {
		carol.recordTransportObservation("mallory", token.OutcomeKept, "earlier parseable discovery frame")
	}
	carol.recordTransportObservation("mallory", token.OutcomeBroken, "malformed bytes")
	carol.markCompleted("transport:malformed:mallory", "observed malformed bytes from mallory")
	message := WireMessage{FromNode: "mallory", ToNode: "carol", Payload: map[string]string{"target_node": "dave", "promised_route": "mallory,carol,dave"}}
	response, handleErr := carol.handleRoutePromise(message)
	if handleErr != nil {
		t.Fatalf("handle route promise: %v", handleErr)
	}
	if response.Outcome != token.OutcomeRefused {
		t.Fatalf("route from low-trust Mallory should be refused: %#v", response)
	}
	if !carol.hasCompleted("route_refused:mallory") {
		t.Fatalf("Carol should record local route refusal evidence for Mallory")
	}
}

func TestRouteRefusalChangesNextHopSelection(t *testing.T) {
	// Intent: A receiver's refusal is local evidence for future route selection, so
	// Mallory can route around Carol without treating Carol's refusal as global
	// authority. Source: DI-sipuz
	mallory := newNode("mallory")
	mallory.learnRoute("dave", []string{"carol", "dave"}, "carol")
	mallory.learnRoute("dave", []string{"frank", "ellen", "dave"}, "frank")
	mallory.markCompleted("route_refused_by:carol", "test route refusal")
	route, routeErr := mallory.selectRoute("dave")
	if routeErr != nil {
		t.Fatalf("select alternate route after Carol refusal: %v", routeErr)
	}
	if got := encodeRoute(route); got != "frank,ellen,dave" {
		t.Fatalf("route after Carol refusal = %s, want frank,ellen,dave", got)
	}
}
