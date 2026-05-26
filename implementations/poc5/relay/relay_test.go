package relay

import (
	"context"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc5/lib"
)

func TestParseRoutesRecordsLocalNextHopPromises(t *testing.T) {
	routes, err := ParseRoutes("Carol=carol:9200,Dave=ellen:9200")
	if err != nil {
		t.Fatalf("ParseRoutes returned error: %v", err)
	}
	if routes["Carol"] != "carol:9200" {
		t.Fatalf("Carol route = %q", routes["Carol"])
	}
	if routes["Dave"] != "ellen:9200" {
		t.Fatalf("Dave route = %q", routes["Dave"])
	}
}

func TestParseRoutesRejectsAmbiguousPromiseText(t *testing.T) {
	if _, err := ParseRoutes("Carol"); err == nil {
		t.Fatalf("ParseRoutes accepted malformed route promise")
	}
}

func TestWrapCarriesExactInnerBytes(t *testing.T) {
	inner, err := lib.NewEnvelope(lib.NewProtocolCID([]byte("inner spec")), map[string]string{
		"kind": "demo_v1",
	})
	if err != nil {
		t.Fatalf("NewEnvelope returned error: %v", err)
	}
	innerBytes, err := inner.Bytes()
	if err != nil {
		t.Fatalf("Bytes returned error: %v", err)
	}
	wrapper, err := Wrap("Alice", "alice-app", "Carol", "carol-app", innerBytes, "hash")
	if err != nil {
		t.Fatalf("Wrap returned error: %v", err)
	}
	kind, fields, err := lib.EnvelopeKind(wrapper)
	if err != nil {
		t.Fatalf("EnvelopeKind returned error: %v", err)
	}
	if kind != "relay_forward_v1" {
		t.Fatalf("kind = %q", kind)
	}
	roundTrip, err := lib.ParseHexBytes(fields["inner_hex"])
	if err != nil {
		t.Fatalf("ParseHexBytes returned error: %v", err)
	}
	if string(roundTrip) != string(innerBytes) {
		t.Fatalf("wrapper changed exact inner bytes")
	}
}

func TestProcessRelayFrameKeepsLocalRefusalNonFatal(t *testing.T) {
	inner, err := lib.NewEnvelope(lib.NewProtocolCID([]byte("inner spec")), map[string]string{
		"kind": "demo_v1",
	})
	if err != nil {
		t.Fatalf("NewEnvelope returned error: %v", err)
	}
	innerBytes, err := inner.Bytes()
	if err != nil {
		t.Fatalf("Bytes returned error: %v", err)
	}
	wrapper, err := Wrap("Alice", "alice-app", "Carol", "carol-app", innerBytes, "")
	if err != nil {
		t.Fatalf("Wrap returned error: %v", err)
	}
	wrapperBytes, err := wrapper.Bytes()
	if err != nil {
		t.Fatalf("Bytes returned error: %v", err)
	}
	relayApp := RelayApp{NodeName: "Alice", AppName: "alice-relay-app", RouteTable: map[string]string{}}
	if err := relayApp.processRelayFrame(context.Background(), wrapperBytes, "test"); err != nil {
		t.Fatalf("processRelayFrame returned error for local refusal: %v", err)
	}
}
