package relationship

import "testing"

func TestBrokenPromiseRemovesDirectPeer(t *testing.T) {
	ledger := NewLedger([]string{"bob"}, []string{"bob"}, 2, -2, 0)
	transition := ledger.ObserveOutcome("bob", OutcomeBroken)
	if ledger.CanDial("bob") {
		t.Fatalf("broken promise should remove direct dial promise")
	}
	if transition != TransitionRemoved {
		t.Fatalf("transition = %s, want %s", transition, TransitionRemoved)
	}
}

func TestRepairCanRestoreDirectPeer(t *testing.T) {
	ledger := NewLedger([]string{"bob"}, []string{"bob"}, 2, -2, 0)
	ledger.ObserveOutcome("bob", OutcomeBroken)
	ledger.ObserveOutcome("bob", OutcomeRepairKept)
	ledger.ObserveOutcome("bob", OutcomeRepairKept)
	ledger.ObserveOutcome("bob", OutcomeRepairKept)
	if !ledger.CanDial("bob") {
		t.Fatalf("kept repair promises should restore direct adjacency")
	}
}

func TestDiscoveryCanCreateDirectPeer(t *testing.T) {
	ledger := NewLedger([]string{"bob"}, nil, 2, -2, 0)
	transition := ledger.ObserveOutcome("bob", OutcomeDiscoveryKept)
	if !ledger.CanDial("bob") {
		t.Fatalf("kept discovery promise should create direct adjacency")
	}
	if transition != TransitionAdded {
		t.Fatalf("transition = %s, want %s", transition, TransitionAdded)
	}
}

func TestNonCommitmentDoesNotReduceTrust(t *testing.T) {
	ledger := NewLedger([]string{"bob"}, []string{"bob"}, 2, -2, 0)
	transition := ledger.ObserveOutcome("bob", OutcomeNonCommitment)
	if ledger.Trust("bob") != 0 {
		t.Fatalf("non-commitment trust = %d, want 0", ledger.Trust("bob"))
	}
	if !ledger.CanDial("bob") {
		t.Fatalf("ordinary non-commitment should not remove an existing direct peer promise")
	}
	if transition != TransitionUnchanged {
		t.Fatalf("transition = %s, want %s", transition, TransitionUnchanged)
	}
}

func TestDecayReducesPositiveTrust(t *testing.T) {
	ledger := NewLedger([]string{"bob"}, []string{"bob"}, 2, -2, 1)
	ledger.ObserveOutcome("bob", OutcomeKept)
	ledger.ObserveOutcome("bob", OutcomeKept)
	ledger.DecayRound()
	if ledger.Trust("bob") != 1 {
		t.Fatalf("trust after decay = %d, want 1", ledger.Trust("bob"))
	}
}

func TestTenRoundDecayAndRepairScenario(t *testing.T) {
	ledger := NewLedger([]string{"bob"}, []string{"bob"}, 2, -2, 1)
	ledger.ObserveOutcome("bob", OutcomeBroken)
	for roundIndex := 0; roundIndex < 10; roundIndex++ {
		ledger.DecayRound()
	}
	if ledger.CanDial("bob") {
		t.Fatalf("broken relationship should stay non-direct after ten decay rounds")
	}
	ledger.ObserveOutcome("bob", OutcomeRepairKept)
	ledger.ObserveOutcome("bob", OutcomeRepairKept)
	ledger.ObserveOutcome("bob", OutcomeRepairKept)
	if !ledger.CanDial("bob") {
		t.Fatalf("kept repair promises should restore direct relationship after long decay")
	}
}

func TestStateRoundTripKeepsLocalTrust(t *testing.T) {
	ledger := NewLedger([]string{"bob"}, []string{"bob"}, 2, -2, 0)
	ledger.ObserveOutcome("bob", OutcomeKept)
	ledger.ObserveOutcome("bob", OutcomeKept)
	state := ledger.Export()
	restored := NewLedger([]string{"bob"}, nil, 2, -2, 0)
	restored.ApplyState(state)
	if restored.Trust("bob") != 2 || !restored.CanDial("bob") {
		t.Fatalf("restored state = trust %d direct %v, want trust 2 direct true", restored.Trust("bob"), restored.CanDial("bob"))
	}
}

func TestStateRoundTripClearsRemovedDirectPeer(t *testing.T) {
	ledger := NewLedger([]string{"bob"}, []string{"bob"}, 2, -2, 0)
	ledger.ObserveOutcome("bob", OutcomeBroken)
	state := ledger.Export()
	restored := NewLedger([]string{"bob"}, []string{"bob"}, 2, -2, 0)
	restored.ApplyState(state)
	if restored.CanDial("bob") {
		t.Fatalf("restored state should preserve removed direct peer")
	}
}
