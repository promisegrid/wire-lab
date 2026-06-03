package runtime

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc10-llm-autonomous-agents/config"
	"promisegrid.dev/wire-lab/implementations/poc10-llm-autonomous-agents/decision"
)

func TestNodeWithNoNeighborsUsesFakeObservationOnly(t *testing.T) {
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
		Agents: []config.AgentConfig{{
			Name:       "alice",
			Profile:    config.ProfileStructuredAction,
			Persona:    "tester",
			Motivation: "test",
			Neighbors:  nil,
		}},
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
}
