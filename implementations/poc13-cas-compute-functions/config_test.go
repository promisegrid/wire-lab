package poc13

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigRejectsStoredAPIKey(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(configPath, []byte(`{"api_key":"not allowed"}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	if _, err := LoadConfig(configPath); err == nil {
		t.Fatalf("config with api_key should be rejected")
	}
}

func TestExampleConfigLoads(t *testing.T) {
	cfg, err := LoadConfig("config.example.json")
	if err != nil {
		t.Fatalf("load example config: %v", err)
	}
	if len(cfg.Agents) != 8 {
		t.Fatalf("agent count = %d, want 8", len(cfg.Agents))
	}
	if _, ok := cfg.Container("alice-bob"); !ok {
		t.Fatalf("alice-bob container missing")
	}
}
