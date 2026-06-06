package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadRejectsSecretBearingFields(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(configPath, []byte(`{"api_key":"not allowed"}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	if _, err := Load(configPath); err == nil {
		t.Fatalf("config with api_key should be rejected")
	}
}

func TestLoadAllowsAPIKeyEnvironmentName(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(configPath, validConfigBytes(), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.APIKeyEnv != "OPENAI_API_KEY" {
		t.Fatalf("api key env = %q, want OPENAI_API_KEY", cfg.APIKeyEnv)
	}
}

func TestExampleConfigLoads(t *testing.T) {
	cfg, err := Load(filepath.Join("..", "config.example.json"))
	if err != nil {
		t.Fatalf("load example config: %v", err)
	}
	if len(cfg.Agents) != 17 {
		t.Fatalf("example config agents = %d, want 17", len(cfg.Agents))
	}
	if _, _, endpointFound := cfg.EndpointFor("mallory"); !endpointFound {
		t.Fatalf("expected endpoint for mallory")
	}
	if appAddress, appAddressFound := cfg.KernelAppAddressForAgent("mallory"); !appAddressFound || appAddress != "127.0.0.1:9114" {
		t.Fatalf("mallory kernel app address = %q, %v; want 127.0.0.1:9114, true", appAddress, appAddressFound)
	}
	if host, port, endpointFound := cfg.KernelPeerEndpointForAgent("mallory"); !endpointFound || host != "mallory-oscar" || port != 9115 {
		t.Fatalf("mallory peer kernel endpoint = %q %d %v, want mallory-oscar 9115 true", host, port, endpointFound)
	}
	fulfillment, ok := cfg.Agent("fulfillment")
	if !ok {
		t.Fatalf("expected fulfillment agent")
	}
	if len(fulfillment.Protocols()) != 1 {
		t.Fatalf("fulfillment protocols = %#v, want relationship-only receive promise", fulfillment.Protocols())
	}
}

func TestAgentRejectsUnknownProtocolName(t *testing.T) {
	agent := AgentConfig{
		Name:           "alice",
		Persona:        "tester",
		Motivation:     "test",
		SupportedPCIDs: []string{"relationship_v1", "bogus_v1"},
	}
	if err := agent.Validate(); err == nil {
		t.Fatalf("unknown supported pCID should be rejected")
	}
}

func TestMonitorWaitTimeoutCoversAgentAndMonitorBudgets(t *testing.T) {
	cfg := Config{
		RequestTimeoutSeconds: 120,
		StartupDelayMillis:    500,
		TurnDelayMillis:       500,
		ShutdownGraceMillis:   45000,
		MaxTurns:              4,
		MaxAgentCalls:         4,
		MaxMonitorCalls:       1,
	}
	want := 648 * time.Second
	if got := cfg.MonitorWaitTimeout(); got != want {
		t.Fatalf("monitor wait timeout = %s, want %s", got, want)
	}
	if cfg.MonitorWaitTimeout() <= 90*time.Second {
		t.Fatalf("monitor wait timeout should exceed the old hard-coded 90s limit")
	}
}

func validConfigBytes() []byte {
	return []byte(`{
		"run_id":"test",
		"run_root":"/run/poc12",
		"listen_port":9100,
		"provider_base_url":"https://example.invalid/v1/responses",
		"api_key_env":"OPENAI_API_KEY",
		"agent_model":"model",
		"monitor_model":"model",
		"reasoning_effort":"medium",
		"service_tier":"default",
		"request_timeout_seconds":1,
		"startup_delay_millis":1,
		"turn_delay_millis":1,
		"max_turns":1,
		"max_agent_calls":12,
		"max_monitor_calls":1,
		"monitor_node":"alice",
		"strong_trust_threshold":2,
		"weak_trust_threshold":-2,
		"trust_decay_per_round":0,
		"agents":[
			{"name":"alice","persona":"tester","motivation":"test","initial_peers":["bob"],"candidate_peers":["bob"],"budget":2,"capacity":2},
			{"name":"bob","persona":"tester","motivation":"test","initial_peers":["alice"],"candidate_peers":["alice"],"budget":2,"capacity":2},
			{"name":"carol","persona":"tester","motivation":"test","initial_peers":["dave"],"candidate_peers":["dave"],"budget":2,"capacity":2},
			{"name":"dave","persona":"tester","motivation":"test","initial_peers":["carol"],"candidate_peers":["carol"],"budget":2,"capacity":2},
			{"name":"ellen","persona":"tester","motivation":"test","initial_peers":["frank"],"candidate_peers":["frank"],"budget":2,"capacity":2},
			{"name":"frank","persona":"tester","motivation":"test","initial_peers":["ellen"],"candidate_peers":["ellen"],"budget":2,"capacity":2},
			{"name":"grace","persona":"tester","motivation":"test","initial_peers":["heidi"],"candidate_peers":["heidi"],"budget":2,"capacity":2},
			{"name":"heidi","persona":"tester","motivation":"test","initial_peers":["grace"],"candidate_peers":["grace"],"budget":2,"capacity":2},
			{"name":"ivan","persona":"tester","motivation":"test","initial_peers":["judy"],"candidate_peers":["judy"],"budget":2,"capacity":2},
			{"name":"judy","persona":"tester","motivation":"test","initial_peers":["ivan"],"candidate_peers":["ivan"],"budget":2,"capacity":2},
			{"name":"mallory","persona":"tester","motivation":"test","initial_peers":["oscar"],"candidate_peers":["oscar"],"budget":2,"capacity":2,"adversarial":true},
			{"name":"oscar","persona":"tester","motivation":"test","initial_peers":["mallory"],"candidate_peers":["mallory"],"budget":2,"capacity":2,"adversarial":true}
		],
		"containers":[
			{"name":"alice","agents":["alice"]},
			{"name":"bob","agents":["bob"]},
			{"name":"carol","agents":["carol"]},
			{"name":"dave","agents":["dave"]},
			{"name":"ellen-frank","agents":["ellen","frank"]},
			{"name":"grace-heidi","agents":["grace","heidi"]},
			{"name":"ivan-judy","agents":["ivan","judy"]},
			{"name":"mallory-oscar","agents":["mallory","oscar"]}
		]
	}`)
}
