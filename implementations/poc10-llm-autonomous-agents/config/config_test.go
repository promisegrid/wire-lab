package config

import (
	"os"
	"path/filepath"
	"testing"
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
	configBytes := []byte(`{
		"run_id":"test",
		"run_root":"/run/poc10",
		"listen_port":9000,
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
		"max_agent_calls":1,
		"max_monitor_calls":1,
		"monitor_node":"alice",
		"agents":[{"name":"alice","profile":"structured_action","persona":"tester","motivation":"test","neighbors":[]}]
	}`)
	if err := os.WriteFile(configPath, configBytes, 0o600); err != nil {
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
	if len(cfg.Agents) != 7 {
		t.Fatalf("example config agents = %d, want 7", len(cfg.Agents))
	}
}
