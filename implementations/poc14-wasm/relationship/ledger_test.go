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

func TestMalformedEvidenceDelaysOrdinaryTrustRecovery(t *testing.T) {
	// Intent: Neat-looking kept promises after malformed evidence should first
	// work off local recovery caution instead of immediately rebuilding direct
	// trust in the peer. Source: DI-fijov
	ledger := NewLedger([]string{"mallory"}, []string{"mallory"}, 2, -2, 0)
	ledger.ObserveOutcome("mallory", OutcomeMalformed)
	for keptIndex := 0; keptIndex < recoveryCautionAfterNegativeEvidence; keptIndex++ {
		ledger.ObserveOutcome("mallory", OutcomeKept)
	}
	if ledger.Trust("mallory") != -3 {
		t.Fatalf("ordinary kept promises during caution changed trust to %d, want -3", ledger.Trust("mallory"))
	}
	if ledger.CanDial("mallory") {
		t.Fatalf("ordinary kept promises during caution should not restore direct adjacency")
	}
	ledger.ObserveOutcome("mallory", OutcomeKept)
	if ledger.Trust("mallory") != -2 {
		t.Fatalf("first post-caution kept promise changed trust to %d, want -2", ledger.Trust("mallory"))
	}
}

func TestTrustScoresSaturate(t *testing.T) {
	// Intent: Trust scores are local relationship evidence, not absolute
	// reputation points, so they stay in a small comparable range. Source:
	// DI-sihuz
	ledger := NewLedger([]string{"bob"}, []string{"bob"}, 2, -2, 0)
	for keptIndex := 0; keptIndex < 20; keptIndex++ {
		ledger.ObserveOutcome("bob", OutcomeKept)
	}
	if ledger.Trust("bob") != maxTrustScore {
		t.Fatalf("positive trust = %d, want max %d", ledger.Trust("bob"), maxTrustScore)
	}
	for brokenIndex := 0; brokenIndex < 20; brokenIndex++ {
		ledger.ObserveOutcome("bob", OutcomeBroken)
	}
	if ledger.Trust("bob") != minTrustScore {
		t.Fatalf("negative trust = %d, want min %d", ledger.Trust("bob"), minTrustScore)
	}
}

func TestCautionIsObservableWithoutMutableState(t *testing.T) {
	// Intent: Runtime tests and analyzer gates need to observe recovery caution
	// without reaching into ledger internals or creating shared authority.
	// Source: DI-sihuz
	ledger := NewLedger([]string{"mallory"}, []string{"mallory"}, 2, -2, 0)
	ledger.ObserveOutcome("mallory", OutcomeMalformed)
	if ledger.Caution("mallory") != recoveryCautionAfterNegativeEvidence {
		t.Fatalf("caution = %d, want %d", ledger.Caution("mallory"), recoveryCautionAfterNegativeEvidence)
	}
	ledger.ObserveOutcome("mallory", OutcomeKept)
	if ledger.Caution("mallory") != recoveryCautionAfterNegativeEvidence-1 {
		t.Fatalf("caution after kept = %d, want %d", ledger.Caution("mallory"), recoveryCautionAfterNegativeEvidence-1)
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

func TestPermanentDistrustBlocksFutureDialAndPersists(t *testing.T) {
	// Intent: Permanent distrust is a durable local restraint promise; kept or
	// repair-looking outcomes should not automatically restore direct contact.
	// Source: DI-dubih
	ledger := NewLedger([]string{"mallory"}, []string{"mallory"}, 2, -2, 0)
	ledger.PermanentlyDistrust("mallory")
	ledger.ObserveOutcome("mallory", OutcomeRepairKept)
	if ledger.CanDial("mallory") || ledger.CanAccept("mallory") {
		t.Fatalf("permanent distrust should block future dial and accept promises")
	}
	if !ledger.PermanentlyDistrusted("mallory") {
		t.Fatalf("permanent distrust state should be visible locally")
	}
	state := ledger.Export()
	restored := NewLedger([]string{"mallory"}, []string{"mallory"}, 2, -2, 0)
	restored.ApplyState(state)
	if !restored.PermanentlyDistrusted("mallory") || restored.CanDial("mallory") {
		t.Fatalf("restored permanent distrust = %v direct = %v, want true/false", restored.PermanentlyDistrusted("mallory"), restored.CanDial("mallory"))
	}
}

func TestTransitExclusionBlocksOnlyTransitHop(t *testing.T) {
	// Intent: Transit exclusion should reject Mallory as a hop through Alice's own
	// traffic while leaving non-Mallory paths and endpoint naming as separate local
	// decisions. Source: DI-dubih
	ledger := NewLedger([]string{"alice", "mallory", "carol", "frank"}, []string{"mallory", "frank"}, 2, -2, 0)
	ledger.ExcludeTransit("mallory")
	if ledger.RouteAllowed([]string{"alice", "mallory", "carol"}) {
		t.Fatalf("route with Mallory as transit should be rejected")
	}
	if !ledger.RouteAllowed([]string{"alice", "frank", "carol"}) {
		t.Fatalf("route through trusted Frank should be allowed")
	}
	if !ledger.RouteAllowed([]string{"alice", "frank", "mallory"}) {
		t.Fatalf("Mallory as endpoint should be a separate send/trust decision, not transit exclusion")
	}
	state := ledger.Export()
	restored := NewLedger([]string{"alice", "mallory", "carol", "frank"}, nil, 2, -2, 0)
	restored.ApplyState(state)
	if !restored.TransitExcluded("mallory") || restored.RouteAllowed([]string{"alice", "mallory", "carol"}) {
		t.Fatalf("restored transit exclusion should reject Mallory transit")
	}
}
