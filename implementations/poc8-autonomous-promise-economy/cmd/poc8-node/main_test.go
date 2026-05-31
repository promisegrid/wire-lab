package main

import (
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc8-autonomous-promise-economy/protocol"
	"promisegrid.dev/wire-lab/implementations/poc8-autonomous-promise-economy/token"
)

func TestPayloadKindsShareOneProtocolCID(t *testing.T) {
	// Intent: POC8 has one promise-economy pCID; message kind selects a payload
	// variant under that one protocol instead of minting a pCID per message type.
	// Source: DI-sirus
	node := newNode("alice")
	payloadKinds := []string{
		kindNeedAdvertisement,
		kindOfferPromise,
		kindCounterPromise,
		kindAcceptancePromise,
		kindTokenIssuePromise,
		kindTokenRedemptionPromise,
		kindOutcomeObservation,
		kindExchangeRateQuote,
	}
	for _, payloadKind := range payloadKinds {
		message, messageErr := node.newWireMessage("bob", payloadKind, []string{"bob"}, map[string]string{"resource_kind": "storage"}, token.Token{})
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

func TestRouteBetweenNeighborsAndRelays(t *testing.T) {
	// Intent: Routes remain signed payload promises; each relay only observes and
	// promises the next hop in its local route position. Source: DI-sirus
	cases := []struct {
		fromNode string
		toNode   string
		want     string
	}{
		{fromNode: "alice", toNode: "bob", want: "bob"},
		{fromNode: "alice", toNode: "carol", want: "bob,carol"},
		{fromNode: "dave", toNode: "bob", want: "carol,bob"},
		{fromNode: "mallory", toNode: "alice", want: "alice"},
	}
	for _, testCase := range cases {
		route, routeErr := routeBetween(testCase.fromNode, testCase.toNode)
		if routeErr != nil {
			t.Fatalf("route %s -> %s: %v", testCase.fromNode, testCase.toNode, routeErr)
		}
		if got := encodeRoute(route); got != testCase.want {
			t.Fatalf("route %s -> %s = %s, want %s", testCase.fromNode, testCase.toNode, got, testCase.want)
		}
	}
}
