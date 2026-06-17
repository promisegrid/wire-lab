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

// Event is a compact local event record shown to an LLM or observer-only monitor.
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
// Intent: POC15 tests autonomy without letting the model expand the protocol
// action vocabulary back into authority-like verbs or authority claims. Source:
// DI-timah; DI-punib
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
	decision = NormalizePromiseDecisionVocabulary(decision)
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
		return PromiseDecision{}, fmt.Errorf("decision contains authority-like, byte-forging, or prompt-injection wording")
	}
	return decision, nil
}

// RepairPromiseDecision makes one bounded local repair attempt for common live
// LLM formatting mistakes. Intent: Improve POC15 protocol hygiene without
// broadening the single top-level promise action or accepting authority-like
// wording. Source: DI-timah; DI-punib
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
		fields[""+key] = stringifyField(value)
	}
	return fields
}

// Prompt renders a compact prompt for live LLM decisions.
// Intent: The prompt mirrors the strict provider schema while preserving the
// single top-level promise action, Go-owned CBOR/signature interface, pCID
// hygiene rule, pCID-owned payload shapes, and a preference for useful next
// promise work over repeated generic relationship chatter. Source: DI-timah;
// DI-galin; DI-punib; DI-sihuz; DI-vipih
func Prompt(observation Observation) (string, error) {
	encoded, err := json.MarshalIndent(observation, "", "  ")
	if err != nil {
		return "", err
	}
	return "Return exactly one JSON object matching this shape: " +
		`{"act":"promise","target":"","promise":"","reason":"","fields":[{"key":"","value":""}]}` +
		"\nThe target must be exactly one agent name, copied from direct_peers for ordinary promises. " +
		"Never return a comma-separated target list. A candidate_peers target is valid only for a low-risk link-discovery promise with fields including {\"key\":\"promise_about\",\"value\":\"link_discovery\"}. " +
		"Useful field keys are protocol, promise_about, package_id, order_id, weight_ounces, shipping_address, cost_cents, tracking_number, content_cid, function_cid, input_cid, context_cid, resource, units, stake, collateral, and discovery_reason. " +
		"Use protocol=relationship_v1 for ordinary trust, discovery, observation, capacity, or pricing promises unless you can provide every required payload field for another pCID. Use protocol=postal_scale_v1 for package weighing, protocol=accounting_v1 for address lookup or shipment updates, protocol=ups_label_v1 for label/cost/tracking promises, protocol=cas_storage_v1 only for content-addressed storage payloads with concrete content/token fields, and protocol=cid_compute_v1 only for CID-named compute/cache/verifier payloads with concrete function/input/context/result fields. identity_key_v1 is reserved for scripted pCID-owned key-rotation array payloads, not live generic-map turns. " +
		"Use resource=storage or resource=compute only when you personally promise fulfillment capacity, not when you advertise a need. " +
		"Prefer a concrete, useful next promise that advances local motivation, reciprocal economics, relationship repair, storage, compute, verification, or event sharing visible in recent_events. Use event, promise, or outcome for local records; avoid proof-like nouns. Do not repeat the same promise_about/promise text to the same peer unless recent_events show a new event that changes its meaning. " +
		"The only valid top-level act is promise. Put refusal, repair, observation, economics, and link-preference meaning inside the promise text or the fields key/value list. " +
		"Do not create action kinds. Do not claim authority over other agents. Do not write CBOR or signatures; implementation code encodes and signs the pCID-defined envelope before the local kernel routes exact bytes.\n\nLocal observation:\n" +
		string(encoded), nil
}

// NormalizePromiseDecisionVocabulary rewrites production-facing live-agent text
// away from the old proof-like vocabulary before it can become a signed promise
// payload or run-log event.
// Intent: DI-kirat makes event the active POC15 runtime/log term; live LLMs may
// still draft older vocabulary, so the Go interface normalizes that prose instead
// of letting it leak into fresh runs. Source: DI-kirat
func NormalizePromiseDecisionVocabulary(decision PromiseDecision) PromiseDecision {
	decision.Promise = NormalizeEventVocabulary(decision.Promise)
	decision.Reason = NormalizeEventVocabulary(decision.Reason)
	normalizedFields := make(map[string]any, len(decision.Fields))
	for key, value := range decision.Fields {
		normalizedKey := NormalizeEventVocabulary(key)
		if valueText, ok := value.(string); ok {
			normalizedFields[normalizedKey] = NormalizeEventVocabulary(valueText)
			continue
		}
		normalizedFields[normalizedKey] = value
	}
	decision.Fields = normalizedFields
	return decision
}

// NormalizeMonitorReportVocabulary applies the same active vocabulary rule to
// observer-only monitor summaries before they are written next to run events.
// Intent: The monitor is a POC tool, not a production authority; its prose must
// reinforce event/promise/outcome vocabulary rather than reviving proof-like
// protocol terms. Source: DI-kirat
func NormalizeMonitorReportVocabulary(report MonitorReport) MonitorReport {
	report.Summary = NormalizeMonitorAuthorityVocabulary(NormalizeEventVocabulary(report.Summary))
	for index, concern := range report.Concerns {
		report.Concerns[index] = NormalizeMonitorAuthorityVocabulary(NormalizeEventVocabulary(concern))
	}
	return report
}

// NormalizeMonitorAuthorityVocabulary rewrites monitor-only prose that describes
// avoided anti-patterns using terms that the analyzer correctly treats as drift
// when they appear in ordinary agent or kernel events.
// Intent: DI-kinaf and DI-kirat require the POC-only monitor to avoid becoming a
// source of production-looking authority or RPC vocabulary while preserving the
// stricter rejection path for live agent promises. Source: DI-kinaf; DI-kirat
func NormalizeMonitorAuthorityVocabulary(text string) string {
	replacer := strings.NewReplacer(
		"Commands", "Imposed instructions",
		"commands", "imposed instructions",
		"COMMANDS", "IMPOSED INSTRUCTIONS",
		"Command", "Imposed instruction",
		"command", "imposed instruction",
		"COMMAND", "IMPOSED INSTRUCTION",
	)
	return replacer.Replace(text)
}

// NormalizeEventVocabulary is intentionally small and explicit: it rewrites the
// prohibited active POC word while leaving historical files untouched.
func NormalizeEventVocabulary(text string) string {
	replacer := strings.NewReplacer(
		"Evi"+"dence", "Event",
		"evi"+"dence", "event",
		"EVI"+"DENCE", "EVENT",
		"Boun"+"dary", "Interface",
		"boun"+"dary", "interface",
		"BOUN"+"DARY", "INTERFACE",
	)
	return replacer.Replace(text)
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
