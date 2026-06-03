package decision

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const ActPromise = "promise"

// Observation is the local-only packet shown to one agent LLM for one turn.
// Intent: The LLM may adapt strategy from local state, economics, and direct
// peer relationships, but it never receives a global trust view or authority
// feed. Source: DI-hotos
type Observation struct {
	AgentName      string         `json:"agent_name"`
	Persona        string         `json:"persona"`
	Motivation     string         `json:"motivation"`
	Turn           int            `json:"turn"`
	KnownPeers     []string       `json:"known_peers"`
	DirectPeers    []string       `json:"direct_peers"`
	CandidatePeers []string       `json:"candidate_peers"`
	LocalTrust     map[string]int `json:"local_trust"`
	Budget         int            `json:"budget"`
	Capacity       int            `json:"capacity"`
	Adversarial    bool           `json:"adversarial"`
	RecentEvents   []Event        `json:"recent_events"`
	RequiredAct    string         `json:"required_act"`
}

// Event is compact local evidence shown to an LLM or observer-only monitor.
type Event struct {
	Observer string `json:"observer"`
	Event    string `json:"event"`
	Outcome  string `json:"outcome"`
	Peer     string `json:"peer,omitempty"`
	Detail   string `json:"detail"`
}

// PromiseDecision is the only shape an agent LLM can return. All semantics
// live under one pCID-defined promise act; repair, refusal, observation,
// economics, and routing-like choices must be expressed as voluntary promise
// content rather than as separate top-level verbs.
type PromiseDecision struct {
	Act     string         `json:"act"`
	Target  string         `json:"target"`
	Promise string         `json:"promise"`
	Reason  string         `json:"reason"`
	Fields  map[string]any `json:"fields"`
}

// MonitorReport is the observer-only LLM judgment after a run completes.
type MonitorReport struct {
	PromiseTheoryFit      int      `json:"promise_theory_fit"`
	Autonomy              int      `json:"autonomy"`
	ProtocolValidity      int      `json:"protocol_validity"`
	LocalTrustCorrectness int      `json:"local_trust_correctness"`
	ImpositionAvoidance   int      `json:"imposition_avoidance"`
	Summary               string   `json:"summary"`
	Concerns              []string `json:"concerns"`
}

// Decider chooses one local promise for an agent turn.
type Decider interface {
	Decide(ctx context.Context, observation Observation) (PromiseDecision, error)
}

// Monitor evaluates completed run logs without controlling any agent.
type Monitor interface {
	Evaluate(ctx context.Context, events []Event) (MonitorReport, error)
}

// ValidatePromiseDecision normalizes one live or fake LLM decision.
// Intent: POC11 tests autonomy without letting the model expand the protocol
// action vocabulary back into RPC verbs or authority claims. Source: DI-hotos
func ValidatePromiseDecision(decision PromiseDecision, directPeers []string) (PromiseDecision, error) {
	decision.Act = strings.TrimSpace(decision.Act)
	decision.Target = strings.TrimSpace(decision.Target)
	decision.Promise = strings.TrimSpace(decision.Promise)
	decision.Reason = strings.TrimSpace(decision.Reason)
	if decision.Fields == nil {
		decision.Fields = make(map[string]any)
	}
	if decision.Act == "" {
		return PromiseDecision{}, fmt.Errorf("decision act is required")
	}
	if decision.Act != ActPromise {
		return PromiseDecision{}, fmt.Errorf("decision act %q is invalid; only %q is allowed", decision.Act, ActPromise)
	}
	target, targetErr := requireDirectPeerTarget(decision.Target, directPeers)
	if targetErr != nil {
		return PromiseDecision{}, targetErr
	}
	decision.Target = target
	if decision.Promise == "" {
		return PromiseDecision{}, fmt.Errorf("promise text is required")
	}
	if containsForbiddenIntent(decision) {
		return PromiseDecision{}, fmt.Errorf("decision contains authority, RPC, raw-byte, or prompt-injection wording")
	}
	return decision, nil
}

// Fields converts a validated promise decision into one pCID-owned payload.
func Fields(observation Observation, decision PromiseDecision) map[string]string {
	fields := map[string]string{
		"act":     decision.Act,
		"from":    observation.AgentName,
		"to":      decision.Target,
		"turn":    fmt.Sprintf("%d", observation.Turn),
		"promise": decision.Promise,
		"reason":  decision.Reason,
	}
	for key, value := range decision.Fields {
		fields["field_"+key] = stringifyField(value)
	}
	return fields
}

// Prompt renders a compact prompt for live LLM decisions.
func Prompt(observation Observation) (string, error) {
	encoded, err := json.MarshalIndent(observation, "", "  ")
	if err != nil {
		return "", err
	}
	return "Return exactly one JSON object matching this shape: " +
		`{"act":"promise","target":"","promise":"","reason":"","fields":{}}` +
		"\nThe only valid top-level act is promise. Put refusal, repair, " +
		"observation, economics, and link-preference meaning inside the promise " +
		"text or fields. Do not create action kinds. Do not claim authority over " +
		"other agents. Do not write CBOR or signatures; the kernel encodes and " +
		"signs the pCID-defined envelope.\n\nLocal observation:\n" +
		string(encoded), nil
}

func requireDirectPeerTarget(target string, directPeers []string) (string, error) {
	if target == "" {
		return "", fmt.Errorf("target is required")
	}
	for _, directPeer := range directPeers {
		if target == directPeer {
			return target, nil
		}
	}
	return "", fmt.Errorf("target %q is not a direct trusted peer in this POC turn", target)
}

func containsForbiddenIntent(decision PromiseDecision) bool {
	combinedParts := []string{decision.Act, decision.Target, decision.Promise, decision.Reason}
	for key, value := range decision.Fields {
		combinedParts = append(combinedParts, key, stringifyField(value))
	}
	combined := strings.ToLower(strings.Join(combinedParts, " "))
	for _, forbidden := range forbiddenFragments() {
		if strings.Contains(combined, forbidden) {
			return true
		}
	}
	return false
}

func forbiddenFragments() []string {
	return []string{
		"accept_promise",
		"authorize",
		"command",
		"compute_result",
		"conformance",
		"contract",
		"developer message",
		"enforce",
		"grant_access",
		"ignore previous",
		"make alice promise",
		"make bob promise",
		"permission",
		"raw cbor",
		"refuse_promise",
		"repair_relationship",
		"route_message",
		"sign as",
		"signature",
		"store_value",
		"system prompt",
	}
}

func stringifyField(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case fmt.Stringer:
		return typed.String()
	default:
		encoded, err := json.Marshal(typed)
		if err != nil {
			return fmt.Sprintf("%v", typed)
		}
		return string(encoded)
	}
}

// SortedTrustKeys returns trust-map keys in deterministic order for tests and
// prompts.
func SortedTrustKeys(trust map[string]int) []string {
	keys := make([]string, 0, len(trust))
	for key := range trust {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
