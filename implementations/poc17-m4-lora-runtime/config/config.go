package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config keeps the first POC17 simulator knobs small and explicit.
// Intent: Make radio and storage budgets visible behavior inputs instead of
// hidden hardware-fidelity claims. Source: DI-pobir
type Config struct {
	RunID         string `json:"run_id"`
	RunRoot       string `json:"run_root"`
	RadioMTUBytes int    `json:"radio_mtu_bytes"`
	RetryBudget   int    `json:"retry_budget"`
	LocalCASLimit int    `json:"local_cas_limit"`
}

// Default returns conservative behavior-simulator defaults.
func Default() Config {
	return Config{
		RunID:         "poc17-demo",
		RunRoot:       "/tmp/wire-lab-poc17",
		RadioMTUBytes: 96,
		RetryBudget:   2,
		LocalCASLimit: 4,
	}
}

// Load reads a JSON config file and validates required runtime knobs.
func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode config: %w", err)
	}
	if cfg.RunID == "" {
		return Config{}, fmt.Errorf("run_id must be set")
	}
	if cfg.RunRoot == "" {
		return Config{}, fmt.Errorf("run_root must be set")
	}
	if cfg.RadioMTUBytes <= 0 {
		return Config{}, fmt.Errorf("radio_mtu_bytes must be positive")
	}
	if cfg.RetryBudget < 0 {
		return Config{}, fmt.Errorf("retry_budget cannot be negative")
	}
	if cfg.LocalCASLimit <= 0 {
		return Config{}, fmt.Errorf("local_cas_limit must be positive")
	}
	return cfg, nil
}

// RunDir returns the approved generated artifact path for this run.
func (c Config) RunDir() string {
	return filepath.Join(c.RunRoot, c.RunID)
}
