package decision

import (
	"context"
	"encoding/json"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc10-llm-autonomous-agents/config"
)

func TestValidateStructuredActionDerivesKind(t *testing.T) {
	decision, err := Validate(config.ProfileStructuredAction, Decision{Action: ActionAdvertiseNeed, Target: "bob"}, []string{"bob"})
	if err != nil {
		t.Fatalf("validate structured action: %v", err)
	}
	if decision.Kind != "need_advertisement" {
		t.Fatalf("kind = %q, want need_advertisement", decision.Kind)
	}
}

func TestValidateFreeformWrapsIntent(t *testing.T) {
	decision, err := Validate(config.ProfileFreeformIntent, Decision{Action: ActionFreeform, Target: "frank", Freeform: "I promise a vague trade."}, []string{"frank"})
	if err != nil {
		t.Fatalf("validate freeform: %v", err)
	}
	if decision.Kind != "freeform_intent" || decision.Promise == "" {
		t.Fatalf("freeform was not wrapped into a protocol payload: %#v", decision)
	}
}

func TestValidateRejectsNonNeighborTarget(t *testing.T) {
	_, err := Validate(config.ProfileStructuredAction, Decision{Action: ActionOfferPromise, Target: "mallory"}, []string{"bob"})
	if err == nil {
		t.Fatalf("non-neighbor target should be rejected")
	}
}

func TestValidateNormalizesMultiTargetToDirectNeighbor(t *testing.T) {
	decision, err := Validate(config.ProfileStructuredAction, Decision{Action: ActionOfferPromise, Target: "bob,ellen"}, []string{"bob", "ellen"})
	if err != nil {
		t.Fatalf("validate multi-target: %v", err)
	}
	if decision.Target != "bob" {
		t.Fatalf("target = %q, want first direct neighbor bob", decision.Target)
	}
}

func TestLiveStyleFieldsAcceptNonStringValues(t *testing.T) {
	var decoded Decision
	raw := []byte(`{"action":"offer_promise","target":"bob","fields":{"capacity_mb":100,"best_effort":true,"neighbors":["bob","ellen"]}}`)
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("decode live-style decision: %v", err)
	}
	obs := Observation{AgentName: "alice", Turn: 1}
	fields := Fields(obs, decoded)
	if fields["field_capacity_mb"] != "100" {
		t.Fatalf("capacity field = %q, want 100", fields["field_capacity_mb"])
	}
	if fields["field_best_effort"] != "true" {
		t.Fatalf("best_effort field = %q, want true", fields["field_best_effort"])
	}
	if fields["field_neighbors"] != `["bob","ellen"]` {
		t.Fatalf("neighbors field = %q, want JSON array", fields["field_neighbors"])
	}
}

func TestFakeDeciderNeedsNoProvider(t *testing.T) {
	fake := &FakeDecider{}
	decision, err := fake.Decide(context.Background(), Observation{
		AgentName:        "alice",
		Profile:          string(config.ProfileStructuredAction),
		NeighborPeers:    []string{"bob"},
		AvailableActions: AvailableActions(config.ProfileStructuredAction),
	})
	if err != nil {
		t.Fatalf("fake decide: %v", err)
	}
	if decision.Action == "" || fake.Calls != 1 {
		t.Fatalf("fake decision not recorded: %#v calls=%d", decision, fake.Calls)
	}
}
