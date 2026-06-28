package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config keeps the first POC17 simulator knobs small and explicit.
// Intent: Make radio, retry, and storage limits visible behavior inputs instead
// of hidden hardware-fidelity claims. Source: DI-pobir; DI-gidul
type Config struct {
	RunID                 string `json:"run_id"`
	RunRoot               string `json:"run_root"`
	RadioMTUBytes         int    `json:"radio_mtu_bytes"`
	RetryLimit            int    `json:"retry_limit"`
	LocalCASLimit         int    `json:"local_cas_limit"`
	RAMByteLimit          uint64 `json:"ram_byte_limit"`
	FlashByteLimit        uint64 `json:"flash_byte_limit"`
	EnergyUnitLimit       uint64 `json:"energy_unit_limit"`
	RadioAirtimeByteLimit uint64 `json:"radio_airtime_byte_limit"`
}

// Default returns behavior-simulator defaults shaped by DN-zaraz's wideband
// LoRa frame-size limits without claiming hardware fidelity.
func Default() Config {
	return Config{
		RunID:                 "poc17-demo",
		RunRoot:               "/tmp/wire-lab-poc17",
		RadioMTUBytes:         200,
		RetryLimit:            2,
		LocalCASLimit:         4,
		RAMByteLimit:          4096,
		FlashByteLimit:        16384,
		EnergyUnitLimit:       100,
		RadioAirtimeByteLimit: 800,
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
	if cfg.RetryLimit < 0 {
		return Config{}, fmt.Errorf("retry_limit cannot be negative")
	}
	if cfg.LocalCASLimit <= 0 {
		return Config{}, fmt.Errorf("local_cas_limit must be positive")
	}
	if cfg.RAMByteLimit == 0 {
		return Config{}, fmt.Errorf("ram_byte_limit must be positive")
	}
	if cfg.FlashByteLimit == 0 {
		return Config{}, fmt.Errorf("flash_byte_limit must be positive")
	}
	if cfg.EnergyUnitLimit == 0 {
		return Config{}, fmt.Errorf("energy_unit_limit must be positive")
	}
	if cfg.RadioAirtimeByteLimit == 0 {
		return Config{}, fmt.Errorf("radio_airtime_byte_limit must be positive")
	}
	return cfg, nil
}

// RunDir returns the approved generated artifact path for this run.
func (c Config) RunDir() string {
	return filepath.Join(c.RunRoot, c.RunID)
}
