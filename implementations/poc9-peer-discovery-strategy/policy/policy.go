package policy

import (
	"fmt"

	"promisegrid.dev/wire-lab/implementations/poc9-peer-discovery-strategy/token"
)

const (
	ActionAdvertise = "advertise_need"
	ActionOffer     = "offer_promise"
	ActionCounter   = "counter_promise"
	ActionAccept    = "accept_promise_terms"
	ActionRedeem    = "redeem_held_promise"
	ActionTransfer  = "transfer_bearer_promise"
)

// ActionContext is the private information one POC9 agent uses when judging an
// opportunity. Intent: Keep economics local to the judging agent; there is no
// global price, global reputation, universal trust table, or policy authority.
// Source: DI-sipuz
type ActionContext struct {
	Agent          string
	Action         string
	Peer           string
	SourcePeer     string
	Issuer         string
	Recipient      string
	ResourceKind   string
	ResourceID     string
	TransferRule   string
	IssuerTrust    int
	PeerTrust      int
	SourceTrust    int
	RecipientTrust int
	ExpectedRisk   int
	ExpectedCost   int
	ExpectedReturn int
	Price          int
	MaxPrice       int
	Stake          int
	Scarcity       int
	ResourceValue  int
	Token          token.Token
}

// Decision records one agent's local economic judgment. It is evidence about
// that agent's reasoning, not a command or conformance decision for anyone else.
// Source: DI-sipuz
type Decision struct {
	Agent     string `json:"agent"`
	Action    string `json:"action"`
	Accepted  bool   `json:"accepted"`
	Score     int    `json:"score"`
	Threshold int    `json:"threshold"`
	Reason    string `json:"reason"`
}

func (decision Decision) Detail() string {
	return fmt.Sprintf("%s decision action=%s accepted=%t score=%d threshold=%d reason=%s", decision.Agent, decision.Action, decision.Accepted, decision.Score, decision.Threshold, decision.Reason)
}

// AgentPolicy is a deterministic local utility model for one POC9 agent.
// Intent: Make the demo reproducible while still requiring every accept,
// counteroffer, redemption, and transfer to clear a local threshold.
// Source: DI-sipuz
type AgentPolicy struct {
	Agent             string
	Thresholds        map[string]int
	ResourceValues    map[string]int
	EvidenceCuriosity int
	DisposalBias      int
	DefaultRisk       int
	DefaultCost       int
}

func ForNode(nodeName string) AgentPolicy {
	agentPolicy := AgentPolicy{
		Agent: nodeName,
		Thresholds: map[string]int{
			ActionAdvertise: 0,
			ActionOffer:     1,
			ActionCounter:   0,
			ActionAccept:    1,
			ActionRedeem:    1,
			ActionTransfer:  1,
		},
		ResourceValues: map[string]int{
			"data":                  4,
			"storage":               4,
			"compute":               4,
			"dataset-public":        2,
			"dataset-private":       7,
			"dataset-expired":       1,
			"alice-report":          7,
			"alice-report-store":    7,
			"alice-report-read":     7,
			"bob-storage-stake":     5,
			"fib-10":                5,
			"carol-compute-alice-1": 5,
		},
		DefaultRisk: 1,
		DefaultCost: 1,
	}
	switch nodeName {
	case "alice":
		agentPolicy.ResourceValues["storage"] = 8
		agentPolicy.ResourceValues["compute"] = 7
		agentPolicy.Thresholds[ActionCounter] = 0
	case "bob":
		agentPolicy.ResourceValues["storage"] = 8
		agentPolicy.Thresholds[ActionAccept] = 0
	case "carol":
		agentPolicy.ResourceValues["compute"] = 8
		agentPolicy.ResourceValues["bob-storage-stake"] = 7
	case "dave":
		agentPolicy.EvidenceCuriosity = 3
		agentPolicy.Thresholds[ActionRedeem] = 0
	case "mallory":
		agentPolicy.DisposalBias = 5
		agentPolicy.Thresholds[ActionTransfer] = 0
		agentPolicy.Thresholds[ActionOffer] = 0
	}
	return agentPolicy
}

func (agentPolicy AgentPolicy) Decide(context ActionContext) Decision {
	score := agentPolicy.score(context)
	threshold := agentPolicy.threshold(context.Action)
	accepted := score >= threshold
	reason := fmt.Sprintf("resource_value=%d price=%d max_price=%d stake=%d scarcity=%d issuer_trust=%d peer_trust=%d source_trust=%d", agentPolicy.value(context.ResourceKind, context.ResourceID, context.ResourceValue), context.Price, context.MaxPrice, context.Stake, context.Scarcity, context.IssuerTrust, context.PeerTrust, context.SourceTrust)
	return Decision{Agent: agentPolicy.Agent, Action: context.Action, Accepted: accepted, Score: score, Threshold: threshold, Reason: reason}
}

// score applies a deliberately small local formula. Intent: POC9 is testing
// promise-economy shape and incentives, not claiming this is the final
// PromiseGrid economics model. Source: DI-sipuz
func (agentPolicy AgentPolicy) score(context ActionContext) int {
	resourceValue := agentPolicy.value(context.ResourceKind, context.ResourceID, context.ResourceValue)
	risk := context.ExpectedRisk
	if risk == 0 {
		risk = agentPolicy.DefaultRisk
	}
	cost := context.ExpectedCost
	if cost == 0 {
		cost = agentPolicy.DefaultCost
	}
	pricePenalty := context.Price
	if context.MaxPrice > 0 && context.Price <= context.MaxPrice {
		pricePenalty = context.Price / 2
	}
	brokenEvidencePenalty := 0
	if context.IssuerTrust < 0 {
		brokenEvidencePenalty += 2
	}
	if context.SourceTrust < 0 {
		brokenEvidencePenalty += 2
	}
	switch context.Action {
	case ActionAdvertise:
		return resourceValue - risk
	case ActionOffer:
		return resourceValue + context.PeerTrust + context.ExpectedReturn - context.Scarcity - risk - cost
	case ActionCounter:
		return resourceValue + context.Stake + context.PeerTrust - pricePenalty - risk
	case ActionAccept:
		return resourceValue + transferValue(context.TransferRule) + context.Stake + context.ExpectedReturn + context.PeerTrust + context.IssuerTrust + context.SourceTrust + agentPolicy.EvidenceCuriosity - pricePenalty - risk - context.Scarcity - brokenEvidencePenalty
	case ActionRedeem:
		return resourceValue + context.IssuerTrust + context.SourceTrust + agentPolicy.EvidenceCuriosity - cost
	case ActionTransfer:
		return resourceValue + agentPolicy.DisposalBias + context.RecipientTrust - cost
	default:
		return -100
	}
}

func (agentPolicy AgentPolicy) threshold(action string) int {
	if threshold, ok := agentPolicy.Thresholds[action]; ok {
		return threshold
	}
	return 1
}

func (agentPolicy AgentPolicy) value(kind string, id string, override int) int {
	if override != 0 {
		return override
	}
	if id != "" {
		if value, ok := agentPolicy.ResourceValues[id]; ok {
			return value
		}
	}
	if kind != "" {
		if value, ok := agentPolicy.ResourceValues[kind]; ok {
			return value
		}
	}
	return 1
}

func transferValue(rule string) int {
	if rule == token.TransferBearer {
		return 2
	}
	if rule == token.TransferNonTransferable {
		return 1
	}
	return 0
}
