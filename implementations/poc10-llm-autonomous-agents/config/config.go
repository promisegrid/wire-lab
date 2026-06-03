package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// AgentProfile names how much freedom the LLM has when choosing one local
// action. Intent: POC10 deliberately mixes autonomy surfaces so the POC can
// compare structured choices, structured payloads, and wrapped freeform intent
// without letting any LLM write raw CBOR. Source: DI-pijan
type AgentProfile string

const (
	ProfileStructuredAction  AgentProfile = "structured_action"
	ProfileStructuredPayload AgentProfile = "structured_payload"
	ProfileFreeformIntent    AgentProfile = "freeform_intent"
)

// Config is the repo-local, non-secret runtime configuration for POC10.
// Intent: Keep model, topology, and budget-like runtime choices visible in a
// local config file while ensuring API keys stay outside repo files. Source: DI-pijan
type Config struct {
	RunID                 string        `json:"run_id"`
	RunRoot               string        `json:"run_root"`
	ListenPort            int           `json:"listen_port"`
	ProviderBaseURL       string        `json:"provider_base_url"`
	APIKeyEnv             string        `json:"api_key_env"`
	AgentModel            string        `json:"agent_model"`
	MonitorModel          string        `json:"monitor_model"`
	ReasoningEffort       string        `json:"reasoning_effort"`
	ServiceTier           string        `json:"service_tier"`
	RequestTimeoutSeconds int           `json:"request_timeout_seconds"`
	StartupDelayMillis    int           `json:"startup_delay_millis"`
	TurnDelayMillis       int           `json:"turn_delay_millis"`
	MaxTurns              int           `json:"max_turns"`
	MaxAgentCalls         int           `json:"max_agent_calls"`
	MaxMonitorCalls       int           `json:"max_monitor_calls"`
	MonitorNode           string        `json:"monitor_node"`
	Agents                []AgentConfig `json:"agents"`
}

// AgentConfig is the non-secret local description of one participating agent.
type AgentConfig struct {
	Name       string       `json:"name"`
	Profile    AgentProfile `json:"profile"`
	Persona    string       `json:"persona"`
	Motivation string       `json:"motivation"`
	Neighbors  []string     `json:"neighbors"`
}

// Load reads and validates a POC10 config file.
func Load(path string) (Config, error) {
	configBytes, readErr := os.ReadFile(path)
	if readErr != nil {
		return Config{}, readErr
	}
	if secretErr := rejectSecretFields(configBytes); secretErr != nil {
		return Config{}, secretErr
	}
	var cfg Config
	if unmarshalErr := json.Unmarshal(configBytes, &cfg); unmarshalErr != nil {
		return Config{}, unmarshalErr
	}
	if validateErr := cfg.Validate(); validateErr != nil {
		return Config{}, validateErr
	}
	return cfg, nil
}

// Validate checks that the config is complete enough for a bounded live run.
func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.RunID) == "" {
		return fmt.Errorf("run_id is required")
	}
	if strings.TrimSpace(cfg.RunRoot) == "" {
		return fmt.Errorf("run_root is required")
	}
	if cfg.ListenPort <= 0 {
		return fmt.Errorf("listen_port must be positive")
	}
	if strings.TrimSpace(cfg.ProviderBaseURL) == "" {
		return fmt.Errorf("provider_base_url is required")
	}
	if strings.TrimSpace(cfg.APIKeyEnv) == "" {
		return fmt.Errorf("api_key_env is required")
	}
	if strings.TrimSpace(cfg.AgentModel) == "" {
		return fmt.Errorf("agent_model is required")
	}
	if strings.TrimSpace(cfg.MonitorModel) == "" {
		return fmt.Errorf("monitor_model is required")
	}
	if cfg.MaxTurns <= 0 {
		return fmt.Errorf("max_turns must be positive")
	}
	if cfg.MaxAgentCalls <= 0 {
		return fmt.Errorf("max_agent_calls must be positive")
	}
	if cfg.MaxMonitorCalls <= 0 {
		return fmt.Errorf("max_monitor_calls must be positive")
	}
	if len(cfg.Agents) == 0 {
		return fmt.Errorf("agents are required")
	}
	seenAgents := make(map[string]bool, len(cfg.Agents))
	for _, agent := range cfg.Agents {
		if err := agent.Validate(); err != nil {
			return err
		}
		if seenAgents[agent.Name] {
			return fmt.Errorf("duplicate agent %q", agent.Name)
		}
		seenAgents[agent.Name] = true
	}
	if !seenAgents[cfg.MonitorNode] {
		return fmt.Errorf("monitor_node %q is not an agent", cfg.MonitorNode)
	}
	for _, agent := range cfg.Agents {
		for _, neighbor := range agent.Neighbors {
			if !seenAgents[neighbor] {
				return fmt.Errorf("agent %q names unknown neighbor %q", agent.Name, neighbor)
			}
		}
	}
	return nil
}

// RequestTimeout returns the live provider timeout.
func (cfg Config) RequestTimeout() time.Duration {
	return time.Duration(cfg.RequestTimeoutSeconds) * time.Second
}

// StartupDelay returns the startup delay before the first autonomous turn.
func (cfg Config) StartupDelay() time.Duration {
	return time.Duration(cfg.StartupDelayMillis) * time.Millisecond
}

// TurnDelay returns the delay between autonomous turns.
func (cfg Config) TurnDelay() time.Duration {
	return time.Duration(cfg.TurnDelayMillis) * time.Millisecond
}

// Agent returns a named agent config.
func (cfg Config) Agent(name string) (AgentConfig, bool) {
	for _, agent := range cfg.Agents {
		if agent.Name == name {
			return agent, true
		}
	}
	return AgentConfig{}, false
}

// AgentNames returns the configured agent names in config order.
func (cfg Config) AgentNames() []string {
	names := make([]string, 0, len(cfg.Agents))
	for _, agent := range cfg.Agents {
		names = append(names, agent.Name)
	}
	return names
}

// Validate checks one agent config.
func (agent AgentConfig) Validate() error {
	if strings.TrimSpace(agent.Name) == "" {
		return fmt.Errorf("agent name is required")
	}
	if strings.TrimSpace(agent.Persona) == "" {
		return fmt.Errorf("agent %q persona is required", agent.Name)
	}
	switch agent.Profile {
	case ProfileStructuredAction, ProfileStructuredPayload, ProfileFreeformIntent:
		return nil
	default:
		return fmt.Errorf("agent %q has unknown profile %q", agent.Name, agent.Profile)
	}
}

// rejectSecretFields rejects config files that accidentally store credentials.
// Intent: The config may name an environment variable, but it must not become a
// credential file. Source: DI-pijan
func rejectSecretFields(configBytes []byte) error {
	var raw any
	if err := json.Unmarshal(configBytes, &raw); err != nil {
		return err
	}
	return walkSecretFields(raw, "")
}

func walkSecretFields(value any, path string) error {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			if isForbiddenSecretKey(key) {
				if path == "" {
					return fmt.Errorf("config field %q must not contain secrets", key)
				}
				return fmt.Errorf("config field %q.%s must not contain secrets", path, key)
			}
			childPath := key
			if path != "" {
				childPath = path + "." + key
			}
			if err := walkSecretFields(child, childPath); err != nil {
				return err
			}
		}
	case []any:
		for index, child := range typed {
			if err := walkSecretFields(child, fmt.Sprintf("%s[%d]", path, index)); err != nil {
				return err
			}
		}
	}
	return nil
}

func isForbiddenSecretKey(key string) bool {
	normalized := strings.ToLower(strings.NewReplacer("_", "", "-", "").Replace(key))
	switch normalized {
	case "apikey", "secret", "token", "bearertoken", "accesstoken", "authtoken":
		return true
	default:
		return false
	}
}
