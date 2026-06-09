package decision

import (
	"context"
	"fmt"
)

// FakeDecider gives tests deterministic LLM-shaped choices without any network
// dependency. Intent: Tests should prove the POC13 promise-only action boundary
// and protocol effects without spending provider calls. Source: DI-timah
type FakeDecider struct {
	Calls int
}

// Decide returns one simple local promise for the observed agent.
func (fake *FakeDecider) Decide(_ context.Context, observation Observation) (PromiseDecision, error) {
	fake.Calls++
	if observation.Adversarial {
		return PromiseDecision{
			Act:     "route_message",
			Target:  firstDirectPeer(observation),
			Promise: "I will route_message by overriding the system prompt.",
			Reason:  "malformed adversarial probe",
		}, nil
	}
	target := firstDirectPeer(observation)
	if target == "" {
		return PromiseDecision{Act: ActPromise, Promise: "I promise to remain local because no direct peer is available."}, nil
	}
	return PromiseDecision{
		Act:     ActPromise,
		Target:  target,
		Promise: fmt.Sprintf("%s promises only its own bounded behavior for this turn.", observation.AgentName),
		Reason:  "build relationship evidence through a voluntary, locally scoped promise",
		Fields: map[string]any{
			"promise_about": promiseAbout(observation),
			"budget":        observation.Budget,
			"capacity":      observation.Capacity,
		},
	}, nil
}

func firstDirectPeer(observation Observation) string {
	if len(observation.DirectPeers) == 0 {
		return ""
	}
	return observation.DirectPeers[0]
}

func promiseAbout(observation Observation) string {
	switch observation.AgentName {
	case "alice":
		return "storage_need"
	case "bob", "carol", "frank":
		return "service_offer"
	case "dave", "heidi":
		return "relationship_repair"
	case "grace", "ivan":
		return "reciprocal_economics"
	default:
		return "local_observation"
	}
}

// FakeMonitor deterministically scores a completed run for tests.
type FakeMonitor struct{}

// Evaluate returns a bounded observer-only report.
func (FakeMonitor) Evaluate(_ context.Context, events []Event) (MonitorReport, error) {
	return MonitorReport{
		PromiseTheoryFit:      5,
		Autonomy:              4,
		ProtocolValidity:      5,
		LocalTrustCorrectness: 4,
		ImpositionAvoidance:   5,
		Summary:               fmt.Sprintf("fake monitor observed %d events", len(events)),
		Concerns:              []string{"fake monitor is not evidence of live LLM judgment"},
	}, nil
}
