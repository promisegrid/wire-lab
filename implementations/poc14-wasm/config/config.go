package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc14-wasm/pcid"
)

// Config is the repo-local, non-secret runtime configuration for POC14.
// Intent: Keep topology, relationship thresholds, model selection, and bounded
// live-LLM execution choices visible while keeping API keys outside repo files.
// Source: DI-timah
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
// authority. Source: DI-timah
type AgentConfig struct {
	Name           string   `json:"name"`
	Persona        string   `json:"persona"`
	Motivation     string   `json:"motivation"`
	InitialPeers   []string `json:"initial_peers"`
	CandidatePeers []string `json:"candidate_peers"`
	Budget         int      `json:"budget"`
	Capacity       int      `json:"capacity"`
	Adversarial    bool     `json:"adversarial"`
	Kind           string   `json:"kind"`
	SupportedPCIDs []string `json:"supported_pcids"`
}

// ContainerConfig names the app processes that share one Docker container and
// one local kernel process. Each app still owns its relationship ledger and
// promise judgment; the kernel only routes exact framed envelopes.
// Intent: Preserve the local-process app/kernel boundary instead of folding app
// trust or workflow policy into the container kernel. Source: DI-galin
type ContainerConfig struct {
	Name   string   `json:"name"`
	Agents []string `json:"agents"`
}

// Load reads and validates a POC14 config file.
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

// Validate checks that the config is complete enough for a bounded POC14 run.
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

// ShutdownGrace returns the receive-promise grace period after active turns finish.
// Intent: Give lagging peers a bounded chance to complete already-planned sends
// before this app closes its local kernel receive promises and writes done
// event records. Source: DI-galin
func (cfg Config) ShutdownGrace() time.Duration {
	return time.Duration(cfg.ShutdownGraceMillis) * time.Millisecond
}

// MonitorWaitTimeout returns the longest time a completed node should wait for
// the observer-only monitor marker before treating the run lifecycle as stuck.
// Intent: Early-finishing deterministic apps must not kill a valid run merely
// because live agents and the monitor can each spend their configured provider
// request budget before `monitor.done` appears. Source: DI-jupob
func (cfg Config) MonitorWaitTimeout() time.Duration {
	agentCallBudget := cfg.RequestTimeout() * time.Duration(cfg.MaxAgentCalls)
	monitorCallBudget := cfg.RequestTimeout() * time.Duration(cfg.MaxMonitorCalls)
	turnDelayBudget := cfg.TurnDelay() * time.Duration(cfg.MaxTurns+1)
	return cfg.StartupDelay() + cfg.ShutdownGrace() + turnDelayBudget + agentCallBudget + monitorCallBudget
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

// ListenPortFor gives legacy tests a deterministic per-agent TCP port. The live
// POC14 runtime uses KernelAppPortForContainer and KernelPeerPortForContainer
// instead so every container has one kernel boundary and local app processes
// register receive promises with that kernel.
// Intent: Keep old test helpers stable while moving production POC14 routing to
// the explicit local kernel/app split. Source: DI-galin
func (cfg Config) ListenPortFor(agentName string) (int, bool) {
	for agentIndex, agent := range cfg.Agents {
		if agent.Name == agentName {
			return cfg.ListenPort + agentIndex, true
		}
	}
	return 0, false
}

// ContainerIndex returns the stable index of one configured container.
func (cfg Config) ContainerIndex(containerName string) (int, bool) {
	for containerIndex, container := range cfg.Containers {
		if container.Name == containerName {
			return containerIndex, true
		}
	}
	return 0, false
}

// KernelAppPortForContainer returns the loopback port used by local app
// processes inside one container to promise receive capability and send
// outbound envelopes through their local kernel.
// Intent: Apps are local processes; they do not expose peer-listening sockets to
// other containers, and they do not bypass the local kernel for wire traffic.
// Source: DI-galin
func (cfg Config) KernelAppPortForContainer(containerName string) (int, bool) {
	containerIndex, ok := cfg.ContainerIndex(containerName)
	if !ok {
		return 0, false
	}
	return cfg.ListenPort + containerIndex*2, true
}

// KernelPeerPortForContainer returns the Docker-network port used by peer
// kernels to forward exact framed envelopes into one container.
// Intent: Kernel-to-kernel TCP carries bytes only; trust and promise judgment
// remain in app processes after local delivery. Source: DI-galin
func (cfg Config) KernelPeerPortForContainer(containerName string) (int, bool) {
	containerIndex, ok := cfg.ContainerIndex(containerName)
	if !ok {
		return 0, false
	}
	return cfg.ListenPort + containerIndex*2 + 1, true
}

// KernelAppAddressForAgent returns the loopback address of the local kernel for
// an app process. It deliberately returns 127.0.0.1 because apps are local
// processes in the same container as their kernel.
// Intent: Avoid modeling apps as remote services; app/kernel communication is a
// local process boundary. Source: DI-galin
func (cfg Config) KernelAppAddressForAgent(agentName string) (string, bool) {
	containerName, containerFound := cfg.ContainerForAgent(agentName)
	if !containerFound {
		return "", false
	}
	appPort, portFound := cfg.KernelAppPortForContainer(containerName)
	if !portFound {
		return "", false
	}
	return fmt.Sprintf("127.0.0.1:%d", appPort), true
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

// EndpointFor resolves the legacy per-agent endpoint used by older tests.
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

// KernelPeerEndpointForAgent resolves the peer-kernel endpoint for the
// container that hosts a target app. Kernels forward by target app name after
// parsing only enough payload to find the local receiver.
// Intent: The container kernel remains a transport router, not a service
// registry authority over app semantics. Source: DI-galin
func (cfg Config) KernelPeerEndpointForAgent(agentName string) (string, int, bool) {
	hostName, hostFound := cfg.ContainerForAgent(agentName)
	if !hostFound {
		return "", 0, false
	}
	peerPort, portFound := cfg.KernelPeerPortForContainer(hostName)
	if !portFound {
		return "", 0, false
	}
	return hostName, peerPort, true
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
	registry := pcid.NewRegistry()
	for _, protocolName := range agent.Protocols() {
		if !registry.Known(protocolName) {
			return fmt.Errorf("agent %q names unknown supported pCID %q", agent.Name, protocolName)
		}
	}
	return nil
}

// Protocols returns the protocol pCIDs this app locally promises to receive.
// Intent: POC14 tests kernel routing by slot-0 pCID to app receive promises;
// absence defaults to relationship_v1 so older generic agents remain
// relationship-only. Source: DI-galin
func (agent AgentConfig) Protocols() []string {
	if len(agent.SupportedPCIDs) == 0 {
		return []string{pcid.RelationshipV1}
	}
	protocols := append([]string{}, agent.SupportedPCIDs...)
	sort.Strings(protocols)
	return protocols
}

// Deterministic reports whether the agent is a local device/system handler
// rather than a live autonomous LLM actor.
// Intent: Keep local hardware, business-system, and heterogeneous runtime-adapter
// roles under Go-owned protocol behavior so LLM autonomy chooses relationship
// intent without inventing device, WASM, or stdio outcomes. Source: DI-pohaj;
// DI-linof
func (agent AgentConfig) Deterministic() bool {
	switch agent.Kind {
	case "postal_scale", "ups_label_printer", "printer_port", "accounting", "wasm_agent", "stdio_agent":
		return true
	default:
		return false
	}
}

// rejectSecretFields rejects config files that accidentally store credentials.
// Intent: The config may name an environment variable, but it must not become a
// credential file. Source: DI-timah
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
