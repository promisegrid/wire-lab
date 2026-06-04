package decision

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const ActPromise = "promise"
const PromiseAboutLinkDiscovery = "link_discovery"

// Observation is the local-only packet shown to one agent LLM for one turn.
// Intent: The LLM may adapt strategy from local state, economics, and direct
// peer relationships, but it never receives a global trust view or authority
// feed. Source: DI-timah
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
	SupportedPCIDs []string       `json:"supported_pcids"`
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
// Intent: POC12 tests autonomy without letting the model expand the protocol
// action vocabulary back into RPC verbs or authority claims. Source: DI-timah
func ValidatePromiseDecision(decision PromiseDecision, directPeers []string) (PromiseDecision, error) {
	return ValidateObservedPromiseDecision(decision, Observation{DirectPeers: directPeers})
}

// ValidateObservedPromiseDecision normalizes a live or fake LLM decision against
// the full local observation. Intent: Candidate-peer link discovery is allowed
// only as a low-risk promise payload meaning, while ordinary promises still
// target exactly one current direct peer. Source: DI-timah
func ValidateObservedPromiseDecision(decision PromiseDecision, observation Observation) (PromiseDecision, error) {
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
	target, targetErr := requireObservedTarget(decision, observation)
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

// RepairPromiseDecision makes one bounded local repair attempt for common live
// LLM formatting mistakes. Intent: Improve POC12 protocol hygiene without
// broadening the single top-level promise action or accepting authority/RPC
// wording. Source: DI-timah
func RepairPromiseDecision(rawDecision PromiseDecision, observation Observation, validationErr error) (PromiseDecision, bool, error) {
	repairedDecision := rawDecision
	if repairedDecision.Fields == nil {
		repairedDecision.Fields = make(map[string]any)
	}
	repairedDecision.Fields["repair_source"] = validationErr.Error()
	repairedDecision.Fields["repair_policy"] = "single bounded local repair before rejection"
	if strings.TrimSpace(repairedDecision.Act) == "" {
		repairedDecision.Act = ActPromise
	}
	if strings.TrimSpace(repairedDecision.Target) == "" {
		repairedDecision.Target = singleRepairTarget(observation, IsLinkDiscoveryDecision(repairedDecision))
	} else if repairedTarget, ok := firstRepairTarget(repairedDecision.Target, observation, IsLinkDiscoveryDecision(repairedDecision)); ok {
		repairedDecision.Target = repairedTarget
		repairedDecision.Fields["repair_target_policy"] = "first allowed target from bundled target text"
	}
	if strings.TrimSpace(repairedDecision.Promise) == "" {
		repairedDecision.Promise = observation.AgentName + " promises only its own bounded non-commitment for this turn."
	}
	validDecision, repairErr := ValidateObservedPromiseDecision(repairedDecision, observation)
	if repairErr != nil {
		return PromiseDecision{}, false, repairErr
	}
	return validDecision, true, nil
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
// Intent: The prompt mirrors the strict provider schema while preserving the
// single top-level promise action and Go-owned CBOR/signature boundary.
// Source: DI-timah
func Prompt(observation Observation) (string, error) {
	encoded, err := json.MarshalIndent(observation, "", "  ")
	if err != nil {
		return "", err
	}
	return "Return exactly one JSON object matching this shape: " +
		`{"act":"promise","target":"","promise":"","reason":"","fields":[{"key":"","value":""}]}` +
		"\nThe target must be exactly one agent name, copied from direct_peers for ordinary promises. " +
		"Never return a comma-separated target list. A candidate_peers target is valid only for a low-risk link-discovery promise with fields including {\"key\":\"promise_about\",\"value\":\"link_discovery\"}. " +
		"Useful field keys are protocol, promise_about, package_id, order_id, weight_ounces, shipping_address, cost_cents, tracking_number, resource, units, stake, collateral, and discovery_reason. " +
		"Use protocol=relationship_v1 for trust/discovery/observation promises, protocol=postal_scale_v1 for package weighing, protocol=accounting_v1 for address lookup or shipment updates, and protocol=ups_label_v1 for label/cost/tracking promises. " +
		"Use resource=storage or resource=compute only when you personally promise fulfillment capacity, not when you advertise a need. " +
		"The only valid top-level act is promise. Put refusal, repair, observation, economics, and link-preference meaning inside the promise text or the fields key/value list. " +
		"Do not create action kinds. Do not claim authority over other agents. Do not write CBOR or signatures; the kernel encodes and signs the pCID-defined envelope.\n\nLocal observation:\n" +
		string(encoded), nil
}

func requireObservedTarget(decision PromiseDecision, observation Observation) (string, error) {
	target := decision.Target
	if target == "" {
		return "", fmt.Errorf("target is required")
	}
	for _, directPeer := range observation.DirectPeers {
		if target == directPeer {
			return target, nil
		}
	}
	if IsLinkDiscoveryDecision(decision) {
		for _, candidatePeer := range observation.CandidatePeers {
			if target == candidatePeer && target != observation.AgentName {
				return target, nil
			}
		}
	}
	return "", fmt.Errorf("target %q is not a direct trusted peer in this POC turn", target)
}

func singleRepairTarget(observation Observation, linkDiscovery bool) string {
	if len(observation.DirectPeers) == 1 {
		return observation.DirectPeers[0]
	}
	if linkDiscovery && len(observation.CandidatePeers) == 1 {
		return observation.CandidatePeers[0]
	}
	return ""
}

func firstRepairTarget(targetText string, observation Observation, linkDiscovery bool) (string, bool) {
	tokens := targetTokens(targetText)
	for _, token := range tokens {
		for _, directPeer := range observation.DirectPeers {
			if token == directPeer {
				return token, true
			}
		}
	}
	if linkDiscovery {
		for _, token := range tokens {
			for _, candidatePeer := range observation.CandidatePeers {
				if token == candidatePeer && token != observation.AgentName {
					return token, true
				}
			}
		}
	}
	return "", false
}

func targetTokens(targetText string) []string {
	split := strings.FieldsFunc(targetText, func(charValue rune) bool {
		return charValue == ',' || charValue == ';' || charValue == '|' || charValue == '/' || charValue == '&' || charValue == '\n' || charValue == '\t' || charValue == ' '
	})
	var tokens []string
	for _, token := range split {
		token = strings.TrimSpace(token)
		if token != "" && token != "and" {
			tokens = append(tokens, token)
		}
	}
	return tokens
}

// IsLinkDiscoveryDecision reports whether a promise asks for low-risk
// candidate-peer discovery rather than an ordinary direct-peer exchange.
// Intent: Link formation remains a promise payload meaning under the same pCID,
// not a new action kind or global routing command. Source: DI-timah
func IsLinkDiscoveryDecision(decision PromiseDecision) bool {
	for _, key := range []string{"promise_about", "meaning", "intent", "link_intent"} {
		if stringifyField(decision.Fields[key]) == PromiseAboutLinkDiscovery {
			return true
		}
	}
	return false
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
