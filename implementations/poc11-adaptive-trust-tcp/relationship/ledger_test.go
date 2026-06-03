package relationship

import "testing"

func TestBrokenPromiseRemovesDirectPeer(t *testing.T) {
	ledger := NewLedger([]string{"bob"}, []string{"bob"}, 2, -2, 0)
	ledger.ObserveOutcome("bob", OutcomeBroken)
	if ledger.CanDial("bob") {
		t.Fatalf("broken promise should remove direct dial promise")
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

func TestDecayReducesPositiveTrust(t *testing.T) {
	ledger := NewLedger([]string{"bob"}, []string{"bob"}, 2, -2, 1)
	ledger.ObserveOutcome("bob", OutcomeKept)
	ledger.ObserveOutcome("bob", OutcomeKept)
	ledger.DecayRound()
	if ledger.Trust("bob") != 1 {
		t.Fatalf("trust after decay = %d, want 1", ledger.Trust("bob"))
	}
}
