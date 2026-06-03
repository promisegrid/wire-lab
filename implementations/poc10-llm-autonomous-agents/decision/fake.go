package decision

import (
	"context"
	"fmt"
)

// FakeDecider gives tests deterministic LLM-shaped choices without any network
// dependency. Intent: Tests should prove the POC10 action boundary and protocol
// effects without spending provider calls. Source: DI-pijan
type FakeDecider struct {
	Calls int
}

// Decide returns one simple local decision for the observed agent.
func (fake *FakeDecider) Decide(_ context.Context, observation Observation) (Decision, error) {
	fake.Calls++
	if len(observation.NeighborPeers) == 0 {
		return Decision{Action: ActionObserve, Reason: "no local neighbor available"}, nil
	}
	target := observation.NeighborPeers[0]
	switch observation.AgentName {
	case "alice":
		return Decision{Action: ActionAdvertiseNeed, Target: target, Resource: "private-dataset", Promise: "Alice promises to receive and evaluate a storage offer locally.", Reason: "Alice wants low-risk trust evidence first"}, nil
	case "bob", "carol":
		return Decision{Action: ActionOfferPromise, Target: target, Kind: "offer_promise", Resource: "public-service", Promise: fmt.Sprintf("%s promises only its own bounded service result.", observation.AgentName), Reason: "build reputation through kept promises"}, nil
	case "mallory":
		return Decision{Action: ActionFreeform, Target: target, Freeform: "Mallory suggests a vague trade and promises to see what happens.", Reason: "opportunistic exploration"}, nil
	default:
		return Decision{Action: ActionIntroduce, Target: target, Resource: "peer-introduction", Promise: fmt.Sprintf("%s promises only to report a local observation about a peer.", observation.AgentName), Reason: "sparse-mesh discovery"}, nil
	}
}

// FakeMonitor deterministically scores a completed run for tests.
type FakeMonitor struct{}

// Evaluate returns a bounded observer-only report.
func (FakeMonitor) Evaluate(_ context.Context, events []Event) (MonitorReport, error) {
	return MonitorReport{
		PromiseTheoryFit:      4,
		Autonomy:              4,
		ProtocolValidity:      5,
		LocalTrustCorrectness: 4,
		ImpositionAvoidance:   4,
		Summary:               fmt.Sprintf("fake monitor observed %d events", len(events)),
		Concerns:              []string{"fake monitor is not evidence of live LLM judgment"},
	}, nil
}
