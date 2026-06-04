package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// Config is the repo-local, non-secret runtime configuration for POC11.
// Intent: Keep topology, relationship thresholds, model selection, and bounded
// live-LLM execution choices visible while keeping API keys outside repo files.
// Source: DI-hotos
type Config struct {
	RunID                 string            `json:"run_id"`
	RunRoot               string            `json:"run_root"`
	ListenPort            int               `json:"listen_port"`
	ProviderBaseURL       string            `json:"provider_base_url"`
	APIKeyEnv             string            `json:"api_key_env"`
	AgentModel            string            `json:"agent_model"`
	MonitorModel          string            `json:"monitor_model"`
	ReasoningEffort       string            `json:"reasoning_effort"`
	ServiceTier           string            `json:"service_tier"`
	RequestTimeoutSeconds int               `json:"request_timeout_seconds"`
	StartupDelayMillis    int               `json:"startup_delay_millis"`
	TurnDelayMillis       int               `json:"turn_delay_millis"`
	ShutdownGraceMillis   int               `json:"shutdown_grace_millis"`
	MaxTurns              int               `json:"max_turns"`
	MaxAgentCalls         int               `json:"max_agent_calls"`
	MaxMonitorCalls       int               `json:"max_monitor_calls"`
	MonitorNode           string            `json:"monitor_node"`
	StrongTrustThreshold  int               `json:"strong_trust_threshold"`
	WeakTrustThreshold    int               `json:"weak_trust_threshold"`
	TrustDecayPerRound    int               `json:"trust_decay_per_round"`
	Agents                []AgentConfig     `json:"agents"`
	Containers            []ContainerConfig `json:"containers"`
}

// AgentConfig is the non-secret local description of one autonomous agent.
// Intent: Personas and motivations are deliberately separate from topology so
// a container can run several agent processes without becoming a shared mind or
// authority. Source: DI-hotos
type AgentConfig struct {
	Name           string   `json:"name"`
	Persona        string   `json:"persona"`
	Motivation     string   `json:"motivation"`
	InitialPeers   []string `json:"initial_peers"`
	CandidatePeers []string `json:"candidate_peers"`
	Budget         int      `json:"budget"`
	Capacity       int      `json:"capacity"`
	Adversarial    bool     `json:"adversarial"`
}

// ContainerConfig names the agents that share one Docker container and kernel
// supervisor. Each agent still runs as a separate process with its own TCP
// listener and local relationship ledger.
type ContainerConfig struct {
	Name   string   `json:"name"`
	Agents []string `json:"agents"`
}

// Load reads and validates a POC11 config file.
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

// Validate checks that the config is complete enough for a bounded POC11 run.
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
	if cfg.StrongTrustThreshold <= cfg.WeakTrustThreshold {
		return fmt.Errorf("strong_trust_threshold must be greater than weak_trust_threshold")
	}
	if cfg.TrustDecayPerRound < 0 {
		return fmt.Errorf("trust_decay_per_round must not be negative")
	}
	if cfg.ShutdownGraceMillis < 0 {
		return fmt.Errorf("shutdown_grace_millis must not be negative")
	}
	if len(cfg.Agents) < 11 {
		return fmt.Errorf("at least 11 agents are required")
	}
	seenAgents, agentErr := cfg.validateAgents()
	if agentErr != nil {
		return agentErr
	}
	if !seenAgents[cfg.MonitorNode] {
		return fmt.Errorf("monitor_node %q is not an agent", cfg.MonitorNode)
	}
	if containerErr := cfg.validateContainers(seenAgents); containerErr != nil {
		return containerErr
	}
	return cfg.validatePeerNames(seenAgents)
}

func (cfg Config) validateAgents() (map[string]bool, error) {
	seenAgents := make(map[string]bool, len(cfg.Agents))
	for _, agent := range cfg.Agents {
		if err := agent.Validate(); err != nil {
			return nil, err
		}
		if seenAgents[agent.Name] {
			return nil, fmt.Errorf("duplicate agent %q", agent.Name)
		}
		seenAgents[agent.Name] = true
	}
	return seenAgents, nil
}

func (cfg Config) validateContainers(seenAgents map[string]bool) error {
	if len(cfg.Containers) == 0 {
		return fmt.Errorf("containers are required")
	}
	containerNames := make(map[string]bool, len(cfg.Containers))
	assignedAgents := make(map[string]string, len(cfg.Agents))
	for _, container := range cfg.Containers {
		if strings.TrimSpace(container.Name) == "" {
			return fmt.Errorf("container name is required")
		}
		if containerNames[container.Name] {
			return fmt.Errorf("duplicate container %q", container.Name)
		}
		containerNames[container.Name] = true
		if len(container.Agents) == 0 {
			return fmt.Errorf("container %q must name at least one agent", container.Name)
		}
		for _, agentName := range container.Agents {
			if !seenAgents[agentName] {
				return fmt.Errorf("container %q names unknown agent %q", container.Name, agentName)
			}
			if priorContainer, exists := assignedAgents[agentName]; exists {
				return fmt.Errorf("agent %q is assigned to both %q and %q", agentName, priorContainer, container.Name)
			}
			assignedAgents[agentName] = container.Name
		}
	}
	for agentName := range seenAgents {
		if _, exists := assignedAgents[agentName]; !exists {
			return fmt.Errorf("agent %q is not assigned to any container", agentName)
		}
	}
	return nil
}

func (cfg Config) validatePeerNames(seenAgents map[string]bool) error {
	for _, agent := range cfg.Agents {
		for _, peerName := range append(append([]string{}, agent.InitialPeers...), agent.CandidatePeers...) {
			if peerName == agent.Name {
				return fmt.Errorf("agent %q cannot peer with itself", agent.Name)
			}
			if !seenAgents[peerName] {
				return fmt.Errorf("agent %q names unknown peer %q", agent.Name, peerName)
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

// ShutdownGrace returns the listener grace period after active turns finish.
// Intent: Give lagging peers a bounded chance to complete already-planned sends
// before this node closes its TCP listener and writes done evidence.
// Source: DI-nanud
func (cfg Config) ShutdownGrace() time.Duration {
	return time.Duration(cfg.ShutdownGraceMillis) * time.Millisecond
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

// Container returns a named container config.
func (cfg Config) Container(name string) (ContainerConfig, bool) {
	for _, container := range cfg.Containers {
		if container.Name == name {
			return container, true
		}
	}
	return ContainerConfig{}, false
}

// AgentNames returns the configured agent names in config order.
func (cfg Config) AgentNames() []string {
	names := make([]string, 0, len(cfg.Agents))
	for _, agent := range cfg.Agents {
		names = append(names, agent.Name)
	}
	return names
}

// ListenPortFor gives each agent process a deterministic TCP port inside the
// container. Intent: Multiple agents can share one container while each keeps a
// separate local kernel boundary and listener. Source: DI-hotos
func (cfg Config) ListenPortFor(agentName string) (int, bool) {
	for agentIndex, agent := range cfg.Agents {
		if agent.Name == agentName {
			return cfg.ListenPort + agentIndex, true
		}
	}
	return 0, false
}

// ContainerForAgent returns the Docker service name that hosts one agent.
func (cfg Config) ContainerForAgent(agentName string) (string, bool) {
	for _, container := range cfg.Containers {
		for _, containerAgent := range container.Agents {
			if containerAgent == agentName {
				return container.Name, true
			}
		}
	}
	return "", false
}

// EndpointFor resolves the TCP host and port for a target agent.
func (cfg Config) EndpointFor(agentName string) (string, int, bool) {
	hostName, hostFound := cfg.ContainerForAgent(agentName)
	if !hostFound {
		return "", 0, false
	}
	listenPort, portFound := cfg.ListenPortFor(agentName)
	if !portFound {
		return "", 0, false
	}
	return hostName, listenPort, true
}

// CandidatePeersFor returns an agent's candidate peer names in deterministic
// order. The runtime may use this to present link-repair options without
// pretending that discovery is a central authority.
func (cfg Config) CandidatePeersFor(agent AgentConfig) []string {
	peers := append([]string{}, agent.CandidatePeers...)
	sort.Strings(peers)
	return peers
}

// Validate checks one agent config.
func (agent AgentConfig) Validate() error {
	if strings.TrimSpace(agent.Name) == "" {
		return fmt.Errorf("agent name is required")
	}
	if strings.TrimSpace(agent.Persona) == "" {
		return fmt.Errorf("agent %q persona is required", agent.Name)
	}
	if strings.TrimSpace(agent.Motivation) == "" {
		return fmt.Errorf("agent %q motivation is required", agent.Name)
	}
	if agent.Budget < 0 {
		return fmt.Errorf("agent %q budget must not be negative", agent.Name)
	}
	if agent.Capacity < 0 {
		return fmt.Errorf("agent %q capacity must not be negative", agent.Name)
	}
	return nil
}

// rejectSecretFields rejects config files that accidentally store credentials.
// Intent: The config may name an environment variable, but it must not become a
// credential file. Source: DI-hotos
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
