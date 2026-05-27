package policy

import (
	"fmt"

	"promisegrid.dev/wire-lab/implementations/poc7-capability-token-exchange/token"
)

const (
	ActionIssue    = "issue_resource_promise"
	ActionAccept   = "accept_received_promise"
	ActionRedeem   = "redeem_held_promise"
	ActionTransfer = "transfer_held_promise"
	ActionTrade    = "accept_reciprocal_exchange"
	ActionQuote    = "quote_exchange_offer"
)

// ActionContext is the local information one agent uses when scoring an
// opportunity. Intent: Keep POC7 autonomy local and deterministic; no agent sees
// a global price, global reputation, or global policy table. Source: DI-rodog
type ActionContext struct {
	Agent              string
	Action             string
	Peer               string
	SourcePeer         string
	Recipient          string
	Issuer             string
	ResourceKind       string
	ResourceID         string
	DesiredKind        string
	DesiredID          string
	TransferRule       string
	IssuerTrust        int
	PeerTrust          int
	SourcePeerTrust    int
	RecipientTrust     int
	ExpectedRisk       int
	ExpectedCost       int
	ExpectedReciprocal int
	Token              token.Token
}

// Decision records one agent's private economic judgment. It is evidence about
// that agent's local reasoning, not a command for any other agent. Source:
// DI-rodog
type Decision struct {
	Agent     string `json:"agent"`
	Action    string `json:"action"`
	Accepted  bool   `json:"accepted"`
	Score     int    `json:"score"`
	Threshold int    `json:"threshold"`
	Reason    string `json:"reason"`
}

// Detail renders a compact evidence string for POC logs.
func (decision Decision) Detail() string {
	return fmt.Sprintf("%s decision action=%s accepted=%t score=%d threshold=%d reason=%s", decision.Agent, decision.Action, decision.Accepted, decision.Score, decision.Threshold, decision.Reason)
}

// AgentPolicy is one agent's deterministic local utility model.
// Intent: Make POC7 agents autonomous enough to accept or refuse opportunities
// while staying reproducible in tests and Docker runs. Source: DI-rodog
type AgentPolicy struct {
	Agent             string
	Thresholds        map[string]int
	ResourceValues    map[string]int
	EvidenceCuriosity int
	DisposalBias      int
	DefaultRisk       int
	DefaultCost       int
}

// ForNode returns the fixed local policy for the named POC agent.
func ForNode(nodeName string) AgentPolicy {
	policy := AgentPolicy{
		Agent: nodeName,
		Thresholds: map[string]int{
			ActionIssue:    1,
			ActionAccept:   2,
			ActionRedeem:   1,
			ActionTransfer: 1,
			ActionTrade:    2,
			ActionQuote:    0,
		},
		ResourceValues: map[string]int{
			"data":               4,
			"storage":            4,
			"compute":            4,
			"dataset-public":     2,
			"dataset-private":    6,
			"dataset-stale":      1,
			"storage-slot-store": 4,
			"storage-slot-read":  4,
			"storage-slot-trade": 3,
			"fib-55":             5,
		},
		DefaultRisk: 1,
		DefaultCost: 1,
	}
	switch nodeName {
	case "alice":
		policy.ResourceValues["data"] = 7
		policy.Thresholds[ActionIssue] = 0
	case "bob":
		policy.ResourceValues["storage"] = 7
	case "carol":
		policy.ResourceValues["compute"] = 7
		policy.ResourceValues["dataset-private"] = 7
	case "dave":
		policy.EvidenceCuriosity = 2
		policy.Thresholds[ActionRedeem] = 0
	case "mallory":
		policy.DisposalBias = 5
		policy.Thresholds[ActionTransfer] = 0
		policy.Thresholds[ActionAccept] = 0
	}
	return policy
}

// Decide scores one opportunity using only the caller's local context.
func (policy AgentPolicy) Decide(context ActionContext) Decision {
	score := policy.score(context)
	threshold := policy.threshold(context.Action)
	accepted := score >= threshold
	reason := fmt.Sprintf("local utility from resource=%d desired=%d issuer_trust=%d peer_trust=%d source_trust=%d recipient_trust=%d", policy.value(context.ResourceKind, context.ResourceID), policy.value(context.DesiredKind, context.DesiredID), context.IssuerTrust, context.PeerTrust, context.SourcePeerTrust, context.RecipientTrust)
	return Decision{Agent: policy.Agent, Action: context.Action, Accepted: accepted, Score: score, Threshold: threshold, Reason: reason}
}

// score applies the same small deterministic formula for every run. Intent:
// Keep the autonomy evidence explainable; the values are POC weights, not a
// final PromiseGrid economics model. Source: DI-rodog
func (policy AgentPolicy) score(context ActionContext) int {
	resourceValue := policy.value(context.ResourceKind, context.ResourceID)
	desiredValue := policy.value(context.DesiredKind, context.DesiredID)
	risk := context.ExpectedRisk
	if risk == 0 {
		risk = policy.DefaultRisk
	}
	cost := context.ExpectedCost
	if cost == 0 {
		cost = policy.DefaultCost
	}
	switch context.Action {
	case ActionIssue:
		return resourceValue + context.ExpectedReciprocal + context.PeerTrust - risk - cost
	case ActionAccept:
		return resourceValue + transferValue(context.TransferRule) + policy.EvidenceCuriosity + context.IssuerTrust + 2*context.SourcePeerTrust - risk
	case ActionRedeem:
		return resourceValue + context.IssuerTrust + context.SourcePeerTrust - cost
	case ActionTransfer:
		return resourceValue + policy.DisposalBias + context.RecipientTrust - cost
	case ActionTrade:
		return resourceValue + desiredValue + transferValue(context.TransferRule) + context.PeerTrust + context.IssuerTrust - risk - cost
	case ActionQuote:
		return context.PeerTrust + context.IssuerTrust
	default:
		return -100
	}
}

// threshold returns the minimum local utility this agent requires for an action.
func (policy AgentPolicy) threshold(action string) int {
	if threshold, ok := policy.Thresholds[action]; ok {
		return threshold
	}
	return 1
}

// value prefers resource-specific local values, then falls back to kind-level
// values so new demo resources still have a bounded default.
func (policy AgentPolicy) value(kind string, id string) int {
	if id != "" {
		if value, ok := policy.ResourceValues[id]; ok {
			return value
		}
	}
	if kind != "" {
		if value, ok := policy.ResourceValues[kind]; ok {
			return value
		}
	}
	return 1
}

// transferValue gives bearer promises a small liquidity premium in this POC.
func transferValue(rule string) int {
	if rule == token.TransferBearer {
		return 1
	}
	return 0
}
