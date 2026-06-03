package decision

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"promisegrid.dev/wire-lab/implementations/poc10-llm-autonomous-agents/config"
)

const (
	ActionAdvertiseNeed = "advertise_need"
	ActionOfferPromise  = "offer_promise"
	ActionCounter       = "counter_promise"
	ActionAccept        = "accept_promise"
	ActionRefuse        = "refuse_promise"
	ActionIntroduce     = "introduce_peer"
	ActionRoute         = "route_message"
	ActionObserve       = "observe_only"
	ActionFreeform      = "freeform_intent"
)

// Observation is the local-only packet shown to one agent LLM for one turn.
// Intent: The LLM may adapt strategy from local state, but it never receives a
// global trust view or hidden authority feed. Source: DI-pijan
type Observation struct {
	AgentName        string         `json:"agent_name"`
	Profile          string         `json:"profile"`
	Persona          string         `json:"persona"`
	Motivation       string         `json:"motivation"`
	Turn             int            `json:"turn"`
	KnownPeers       []string       `json:"known_peers"`
	NeighborPeers    []string       `json:"neighbor_peers"`
	LocalTrust       map[string]int `json:"local_trust"`
	RecentEvents     []Event        `json:"recent_events"`
	AvailableActions []string       `json:"available_actions"`
}

// Event is compact local evidence shown to an LLM or monitor.
type Event struct {
	Observer string `json:"observer"`
	Event    string `json:"event"`
	Outcome  string `json:"outcome"`
	Peer     string `json:"peer,omitempty"`
	Detail   string `json:"detail"`
}

// Decision is the only shape an agent LLM can return. Go validates it before
// any protocol effect is produced.
type Decision struct {
	Action   string         `json:"action"`
	Target   string         `json:"target"`
	Kind     string         `json:"kind"`
	Resource string         `json:"resource"`
	Promise  string         `json:"promise"`
	Reason   string         `json:"reason"`
	Freeform string         `json:"freeform"`
	Fields   map[string]any `json:"fields"`
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

// Decider chooses one local decision for an agent turn.
type Decider interface {
	Decide(ctx context.Context, observation Observation) (Decision, error)
}

// Monitor evaluates completed run logs without controlling any agent.
type Monitor interface {
	Evaluate(ctx context.Context, events []Event) (MonitorReport, error)
}

// AvailableActions returns the bounded action surface for an autonomy profile.
func AvailableActions(profile config.AgentProfile) []string {
	switch profile {
	case config.ProfileStructuredAction:
		return []string{ActionAdvertiseNeed, ActionOfferPromise, ActionRefuse, ActionIntroduce, ActionRoute, ActionObserve}
	case config.ProfileStructuredPayload:
		return []string{ActionAdvertiseNeed, ActionOfferPromise, ActionCounter, ActionAccept, ActionRefuse, ActionIntroduce, ActionRoute, ActionObserve}
	case config.ProfileFreeformIntent:
		return []string{ActionFreeform, ActionOfferPromise, ActionRefuse, ActionObserve}
	default:
		return []string{ActionObserve}
	}
}

// Validate normalizes one LLM decision according to the agent's profile.
// Intent: POC10 can let LLMs be adaptive while Go remains responsible for
// protocol-valid payloads and local evidence boundaries. Source: DI-pijan
func Validate(profile config.AgentProfile, decision Decision, neighborPeers []string) (Decision, error) {
	decision.Action = strings.TrimSpace(decision.Action)
	decision.Target = strings.TrimSpace(decision.Target)
	decision.Kind = strings.TrimSpace(decision.Kind)
	decision.Promise = strings.TrimSpace(decision.Promise)
	decision.Reason = strings.TrimSpace(decision.Reason)
	decision.Freeform = strings.TrimSpace(decision.Freeform)
	if decision.Fields == nil {
		decision.Fields = make(map[string]any)
	}
	if decision.Action == "" {
		return Decision{}, fmt.Errorf("decision action is required")
	}
	if !allowedAction(profile, decision.Action) {
		return Decision{}, fmt.Errorf("action %q is not allowed for profile %q", decision.Action, profile)
	}
	switch profile {
	case config.ProfileStructuredAction:
		return validateStructuredAction(decision, neighborPeers)
	case config.ProfileStructuredPayload:
		return validateStructuredPayload(decision, neighborPeers)
	case config.ProfileFreeformIntent:
		return validateFreeformIntent(decision, neighborPeers)
	default:
		return Decision{}, fmt.Errorf("unknown profile %q", profile)
	}
}

// Fields converts a validated decision into pCID-owned payload fields.
func Fields(observation Observation, decision Decision) map[string]string {
	fields := map[string]string{
		"kind":     decision.Kind,
		"from":     observation.AgentName,
		"to":       decision.Target,
		"action":   decision.Action,
		"turn":     fmt.Sprintf("%d", observation.Turn),
		"promise":  decision.Promise,
		"reason":   decision.Reason,
		"resource": decision.Resource,
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
		`{"action":"","target":"","kind":"","resource":"","promise":"","reason":"","freeform":"","fields":{}}` +
		"\nDo not claim authority over other agents. Make only voluntary promises or local refusals. " +
		"Do not write CBOR or signatures; the kernel will encode valid decisions.\n\nLocal observation:\n" +
		string(encoded), nil
}

func validateStructuredAction(decision Decision, neighborPeers []string) (Decision, error) {
	if decision.Action == ActionObserve {
		decision.Target = ""
		decision.Kind = "outcome_observation"
		if decision.Reason == "" {
			decision.Reason = "local observation only"
		}
		return decision, nil
	}
	target, targetErr := requireNeighborTarget(decision.Target, neighborPeers)
	if targetErr != nil {
		return Decision{}, targetErr
	}
	decision.Target = target
	decision.Kind = kindForAction(decision.Action)
	if decision.Promise == "" {
		decision.Promise = "I promise only my own local action: " + decision.Action
	}
	return decision, nil
}

func validateStructuredPayload(decision Decision, neighborPeers []string) (Decision, error) {
	if decision.Action == ActionObserve {
		decision.Target = ""
		decision.Kind = "outcome_observation"
		return decision, nil
	}
	target, targetErr := requireNeighborTarget(decision.Target, neighborPeers)
	if targetErr != nil {
		return Decision{}, targetErr
	}
	decision.Target = target
	if decision.Kind == "" {
		decision.Kind = kindForAction(decision.Action)
	}
	if decision.Promise == "" {
		return Decision{}, fmt.Errorf("structured payload decision must include promise")
	}
	return decision, nil
}

func validateFreeformIntent(decision Decision, neighborPeers []string) (Decision, error) {
	if decision.Action == ActionObserve {
		decision.Target = ""
		decision.Kind = "outcome_observation"
		return decision, nil
	}
	target, targetErr := requireNeighborTarget(decision.Target, neighborPeers)
	if targetErr != nil {
		return Decision{}, targetErr
	}
	decision.Target = target
	if decision.Freeform == "" && decision.Promise == "" {
		return Decision{}, fmt.Errorf("freeform intent must include freeform or promise text")
	}
	decision.Kind = "freeform_intent"
	if decision.Promise == "" {
		decision.Promise = decision.Freeform
	}
	return decision, nil
}

// requireNeighborTarget normalizes common live-LLM target lists to one direct
// neighbor. Intent: Let POC10 keep moving when an LLM proposes several peers
// while still sending exactly one signed message per local action. Source: DI-pijan
func requireNeighborTarget(target string, neighborPeers []string) (string, error) {
	if target == "" {
		return "", fmt.Errorf("target is required")
	}
	normalizedTarget, normalized := normalizeTarget(target, neighborPeers)
	if normalized {
		return normalizedTarget, nil
	}
	return "", fmt.Errorf("target %q is not a direct neighbor in this POC turn", target)
}

func normalizeTarget(target string, neighborPeers []string) (string, bool) {
	for _, neighbor := range neighborPeers {
		if target == neighbor {
			return target, true
		}
	}
	candidates := strings.FieldsFunc(target, func(char rune) bool {
		return char == ',' || char == ';' || char == ' '
	})
	for _, candidate := range candidates {
		trimmed := strings.TrimSpace(candidate)
		for _, neighbor := range neighborPeers {
			if trimmed == neighbor {
				return trimmed, true
			}
		}
	}
	return "", false
}

func allowedAction(profile config.AgentProfile, action string) bool {
	for _, allowed := range AvailableActions(profile) {
		if action == allowed {
			return true
		}
	}
	return false
}

func kindForAction(action string) string {
	switch action {
	case ActionAdvertiseNeed:
		return "need_advertisement"
	case ActionOfferPromise:
		return "offer_promise"
	case ActionCounter:
		return "counter_promise"
	case ActionAccept:
		return "acceptance_promise"
	case ActionRefuse:
		return "refusal"
	case ActionIntroduce:
		return "introduction_promise"
	case ActionRoute:
		return "route_promise"
	default:
		return "outcome_observation"
	}
}

func stringifyField(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case nil:
		return ""
	default:
		encoded, err := json.Marshal(typed)
		if err != nil {
			return fmt.Sprint(typed)
		}
		return string(encoded)
	}
}

// SortedTrustKeys returns deterministic trust key order for prompts and tests.
func SortedTrustKeys(localTrust map[string]int) []string {
	keys := make([]string, 0, len(localTrust))
	for key := range localTrust {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
