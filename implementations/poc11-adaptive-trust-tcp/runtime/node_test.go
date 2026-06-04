package runtime

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc11-adaptive-trust-tcp/config"
	"promisegrid.dev/wire-lab/implementations/poc11-adaptive-trust-tcp/decision"
)

func TestNodeWithNoDirectPeersRecordsLocalNonCommitment(t *testing.T) {
	cfg := singleNodeTestConfig(t)
	node := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := node.Run(ctx); err != nil {
		t.Fatalf("run node: %v", err)
	}
	if len(node.events) == 0 {
		t.Fatalf("node should record local evidence")
	}
	if len(node.events) < 2 {
		t.Fatalf("node should record server and turn events: %#v", node.events)
	}
	if node.events[0].Event != "server_skipped" || node.events[1].Event != "local_non_commitment" {
		t.Fatalf("node events did not record local non-commitment: %#v", node.events)
	}
}

func TestRelationshipStatePersistsAcrossRuns(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	alice.observeOutcome("bob", "kept")
	if err := alice.saveRelationshipState(); err != nil {
		t.Fatalf("save relationship state: %v", err)
	}
	reloadedAlice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	if err := reloadedAlice.loadRelationshipState(); err != nil {
		t.Fatalf("load relationship state: %v", err)
	}
	if reloadedAlice.ledger.Trust("bob") != 1 {
		t.Fatalf("reloaded trust = %d, want 1", reloadedAlice.ledger.Trust("bob"))
	}
}

func TestResourcePromiseChecksLocalCapacity(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	err := alice.checkLocalResourcePromise(map[string]string{
		"field_resource": "storage",
		"field_units":    "3",
	})
	if err == nil {
		t.Fatalf("storage promise beyond local capacity should fail")
	}
}

func TestBrokenPromiseStakeCostReducesBudget(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	alice.applyBrokenPromiseCost("bob", map[string]string{"field_stake": "2"}, "test broken promise")
	if alice.budget != 3 {
		t.Fatalf("budget after stake cost = %d, want 3", alice.budget)
	}
}

func TestRunTurnRepairsMissingActAndTarget(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], missingShapeDecider{}, decision.FakeMonitor{})
	if err := alice.runTurn(context.Background(), 0); err != nil {
		t.Fatalf("run turn: %v", err)
	}
	if !hasEvent(alice.events, "decision_repaired") {
		t.Fatalf("expected decision_repaired event, got %#v", alice.events)
	}
}

func singleNodeTestConfig(t *testing.T) config.Config {
	t.Helper()
	return config.Config{
		RunID:                 "test",
		RunRoot:               filepath.Join(t.TempDir(), "run"),
		ListenPort:            0,
		ProviderBaseURL:       "https://example.invalid/v1/responses",
		APIKeyEnv:             "OPENAI_API_KEY",
		AgentModel:            "model",
		MonitorModel:          "model",
		ReasoningEffort:       "medium",
		ServiceTier:           "default",
		RequestTimeoutSeconds: 1,
		StartupDelayMillis:    1,
		TurnDelayMillis:       1,
		MaxTurns:              1,
		MaxAgentCalls:         1,
		MaxMonitorCalls:       1,
		MonitorNode:           "alice",
		StrongTrustThreshold:  2,
		WeakTrustThreshold:    -2,
		TrustDecayPerRound:    0,
		Agents: []config.AgentConfig{{
			Name:       "alice",
			Persona:    "tester",
			Motivation: "test",
			Budget:     1,
			Capacity:   1,
		}},
		Containers: []config.ContainerConfig{{Name: "alice", Agents: []string{"alice"}}},
	}
}

func twoNodeTestConfig(t *testing.T) config.Config {
	t.Helper()
	return config.Config{
		RunID:                 "test",
		RunRoot:               filepath.Join(t.TempDir(), "run"),
		ListenPort:            0,
		ProviderBaseURL:       "https://example.invalid/v1/responses",
		APIKeyEnv:             "OPENAI_API_KEY",
		AgentModel:            "model",
		MonitorModel:          "model",
		ReasoningEffort:       "medium",
		ServiceTier:           "default",
		RequestTimeoutSeconds: 1,
		StartupDelayMillis:    1,
		TurnDelayMillis:       1,
		MaxTurns:              1,
		MaxAgentCalls:         1,
		MaxMonitorCalls:       1,
		MonitorNode:           "alice",
		StrongTrustThreshold:  2,
		WeakTrustThreshold:    -2,
		TrustDecayPerRound:    0,
		Agents: []config.AgentConfig{{
			Name:         "alice",
			Persona:      "tester",
			Motivation:   "test",
			InitialPeers: []string{"bob"},
			Budget:       5,
			Capacity:     1,
		}, {
			Name:         "bob",
			Persona:      "tester",
			Motivation:   "test",
			InitialPeers: []string{"alice"},
			Budget:       5,
			Capacity:     1,
		}},
		Containers: []config.ContainerConfig{
			{Name: "alice", Agents: []string{"alice"}},
			{Name: "bob", Agents: []string{"bob"}},
		},
	}
}

type missingShapeDecider struct{}

func (missingShapeDecider) Decide(_ context.Context, _ decision.Observation) (decision.PromiseDecision, error) {
	return decision.PromiseDecision{
		Promise: "Alice promises one bounded local exchange.",
		Fields:  map[string]any{"purpose": "repair-test"},
	}, nil
}

func hasEvent(events []decision.Event, eventName string) bool {
	for _, event := range events {
		if event.Event == eventName {
			return true
		}
	}
	return false
}
