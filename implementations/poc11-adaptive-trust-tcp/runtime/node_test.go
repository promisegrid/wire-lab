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
	cfg := config.Config{
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
