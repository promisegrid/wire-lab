package decision

import (
	"context"
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
