package poc13

import (
	"context"
	"path/filepath"
	"testing"
)

func TestAllPOC13AgentsProduceRequiredEvidence(t *testing.T) {
	cfg := testConfig(t)
	for _, agent := range cfg.Agents {
		runner := NewRunner(cfg, agent, scriptedOnlyDecider{})
		if err := runner.Run(context.Background()); err != nil {
			t.Fatalf("run %s: %v", agent.Name, err)
		}
	}
	summary, err := AnalyzeRun(filepath.Join(cfg.RunRoot, cfg.RunID))
	if err != nil {
		t.Fatalf("analyze run: %v", err)
	}
	if err := ValidateAnalysis(summary); err != nil {
		t.Fatalf("validate analysis: %v", err)
	}
}

type scriptedOnlyDecider struct{}

func (scriptedOnlyDecider) Decide(context.Context, Config, AgentConfig, string) (DecisionResult, error) {
	return DecisionResult{Mode: "scripted", Text: "test promise judgment"}, nil
}

func testConfig(t *testing.T) Config {
	t.Helper()
	return Config{
		RunID:                 "test",
		RunRoot:               filepath.Join(t.TempDir(), "run"),
		ProviderBaseURL:       "https://example.invalid/v1/responses",
		APIKeyEnv:             "OPENAI_API_KEY",
		AgentModel:            "gpt-5.4",
		ReasoningEffort:       "medium",
		RequestTimeoutSeconds: 1,
		Containers: []ContainerConfig{
			{Name: "alice-bob", Agents: []string{"alice", "bob"}},
			{Name: "carol-dave", Agents: []string{"carol", "dave"}},
			{Name: "ellen-frank", Agents: []string{"ellen", "frank"}},
			{Name: "grace-mallory", Agents: []string{"grace", "mallory"}},
		},
		Agents: []AgentConfig{
			{Name: "alice", Role: "data_holder", Persona: "Alice", Peers: []string{"bob", "carol"}},
			{Name: "bob", Role: "storage_peer", Persona: "Bob", Peers: []string{"alice"}},
			{Name: "carol", Role: "compute_peer", Persona: "Carol", Peers: []string{"alice"}},
			{Name: "dave", Role: "cache_peer", Persona: "Dave", Peers: []string{"carol"}},
			{Name: "ellen", Role: "context_peer", Persona: "Ellen", Peers: []string{"carol"}},
			{Name: "frank", Role: "replication_peer", Persona: "Frank", Peers: []string{"bob"}},
			{Name: "grace", Role: "verifier_peer", Persona: "Grace", Peers: []string{"mallory"}},
			{Name: "mallory", Role: "adversary_peer", Persona: "Mallory", Peers: []string{"grace"}},
		},
	}
}
