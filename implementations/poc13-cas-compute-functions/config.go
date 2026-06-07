package poc13

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// Config is the non-secret runtime description for POC13.
// Intent: Keep provider selection, container layout, and run paths explicit
// while preventing API keys from entering committed config. Source: DI-notig
type Config struct {
	RunID                  string            `json:"run_id"`
	RunRoot                string            `json:"run_root"`
	ProviderBaseURL        string            `json:"provider_base_url"`
	APIKeyEnv              string            `json:"api_key_env"`
	AgentModel             string            `json:"agent_model"`
	ReasoningEffort        string            `json:"reasoning_effort"`
	ServiceTier            string            `json:"service_tier"`
	RequestTimeoutSeconds  int               `json:"request_timeout_seconds"`
	LiveDecisions          bool              `json:"live_decisions"`
	ListenPort             int               `json:"listen_port"`
	ReadinessTimeoutMillis int               `json:"readiness_timeout_millis"`
	CompletionIdleMillis   int               `json:"completion_idle_millis"`
	Containers             []ContainerConfig `json:"containers"`
	Agents                 []AgentConfig     `json:"agents"`
}

// ContainerConfig binds two local agents to one Docker container for the first
// POC13 runtime shape.
type ContainerConfig struct {
	Name   string   `json:"name"`
	Agents []string `json:"agents"`
}

// AgentConfig describes one autonomous local agent and the peers it is willing
// to consider during this bounded proof of concept.
type AgentConfig struct {
	Name    string   `json:"name"`
	Role    string   `json:"role"`
	Persona string   `json:"persona"`
	Peers   []string `json:"peers"`
}

// LoadConfig reads and validates a POC13 config file.
func LoadConfig(path string) (Config, error) {
	configBytes, readErr := os.ReadFile(path)
	if readErr != nil {
		return Config{}, readErr
	}
	if secretErr := rejectConfigSecrets(configBytes); secretErr != nil {
		return Config{}, secretErr
	}
	var cfg Config
	if err := json.Unmarshal(configBytes, &cfg); err != nil {
		return Config{}, err
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Validate checks only the structural fields needed by the local runtime.
func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.RunID) == "" {
		return fmt.Errorf("run_id is required")
	}
	if strings.TrimSpace(cfg.RunRoot) == "" {
		return fmt.Errorf("run_root is required")
	}
	if strings.TrimSpace(cfg.ProviderBaseURL) == "" {
		return fmt.Errorf("provider_base_url is required")
	}
	if strings.TrimSpace(cfg.APIKeyEnv) == "" {
		return fmt.Errorf("api_key_env is required")
	}
	if cfg.RequestTimeoutSeconds <= 0 {
		return fmt.Errorf("request_timeout_seconds must be positive")
	}
	if cfg.ListenPort < 0 {
		return fmt.Errorf("listen_port must be non-negative")
	}
	if cfg.ReadinessTimeoutMillis <= 0 {
		return fmt.Errorf("readiness_timeout_millis must be positive")
	}
	if cfg.CompletionIdleMillis <= 0 {
		return fmt.Errorf("completion_idle_millis must be positive")
	}
	if len(cfg.Containers) == 0 {
		return fmt.Errorf("at least one container is required")
	}
	if len(cfg.Agents) == 0 {
		return fmt.Errorf("at least one agent is required")
	}
	agentNames := make(map[string]bool)
	for _, agent := range cfg.Agents {
		if strings.TrimSpace(agent.Name) == "" {
			return fmt.Errorf("agent name is required")
		}
		if strings.TrimSpace(agent.Role) == "" {
			return fmt.Errorf("agent %s role is required", agent.Name)
		}
		agentNames[agent.Name] = true
	}
	for _, container := range cfg.Containers {
		if strings.TrimSpace(container.Name) == "" {
			return fmt.Errorf("container name is required")
		}
		if len(container.Agents) == 0 {
			return fmt.Errorf("container %s must name agents", container.Name)
		}
		for _, agentName := range container.Agents {
			if !agentNames[agentName] {
				return fmt.Errorf("container %s names unknown agent %s", container.Name, agentName)
			}
		}
	}
	return nil
}

// Timeout returns the provider decision timeout for live LLM calls.
func (cfg Config) Timeout() time.Duration {
	return time.Duration(cfg.RequestTimeoutSeconds) * time.Second
}

// ReadinessTimeout bounds how long one container waits for peer readiness
// evidence before failing the run.
// Intent: POC13 startup should be driven by explicit readiness markers rather
// than arbitrary sleeps while remaining bounded and diagnosable. Source:
// DI-mosil
func (cfg Config) ReadinessTimeout() time.Duration {
	return time.Duration(cfg.ReadinessTimeoutMillis) * time.Millisecond
}

// CompletionIdle is the quiet period required after local agents complete and
// no TCP handlers remain active before the runtime records done evidence.
// Source: DI-mosil
func (cfg Config) CompletionIdle() time.Duration {
	return time.Duration(cfg.CompletionIdleMillis) * time.Millisecond
}

// Container returns one named container config.
func (cfg Config) Container(name string) (ContainerConfig, bool) {
	for _, container := range cfg.Containers {
		if container.Name == name {
			return container, true
		}
	}
	return ContainerConfig{}, false
}

// Agent returns one named agent config.
func (cfg Config) Agent(name string) (AgentConfig, bool) {
	for _, agent := range cfg.Agents {
		if agent.Name == name {
			return agent, true
		}
	}
	return AgentConfig{}, false
}

func rejectConfigSecrets(configBytes []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(configBytes, &raw); err != nil {
		return err
	}
	for key := range raw {
		normalized := strings.ToLower(key)
		if strings.Contains(normalized, "api_key") && normalized != "api_key_env" {
			return fmt.Errorf("config field %q must not contain secrets", key)
		}
		if strings.Contains(normalized, "secret") || strings.Contains(normalized, "token") {
			return fmt.Errorf("config field %q must not contain secrets", key)
		}
	}
	return nil
}
