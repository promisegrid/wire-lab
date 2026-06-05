package runtime

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/config"
	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/decision"
	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/pcid"
	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/production"
	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/relationship"
)

func TestNodeWithNoDirectPeersRecordsLocalNonCommitment(t *testing.T) {
	cfg := singleNodeTestConfig(t)
	node := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := node.Run(ctx); err != nil {
		t.Fatalf("run node: %v", err)
	}
	if len(node.events) == 0 {
		t.Fatalf("node should record local evidence")
	}
	if len(node.events) < 2 {
		t.Fatalf("node should record kernel registration and turn events: %#v", node.events)
	}
	if node.events[0].Event != "app_kernel_registration_skipped" || node.events[1].Event != "local_non_commitment" {
		t.Fatalf("node events did not record local non-commitment: %#v", node.events)
	}
}

func TestRelationshipStatePersistsAcrossRuns(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	alice.observeOutcome("bob", "kept")
	if err := alice.saveRelationshipState(); err != nil {
		t.Fatalf("save relationship state: %v", err)
	}
	reloadedAlice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	if err := reloadedAlice.loadRelationshipState(); err != nil {
		t.Fatalf("load relationship state: %v", err)
	}
	if reloadedAlice.ledger.Trust("bob") != 1 {
		t.Fatalf("reloaded trust = %d, want 1", reloadedAlice.ledger.Trust("bob"))
	}
}

func TestResourcePromiseChecksLocalCapacity(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	err := alice.checkLocalResourcePromise(map[string]string{
		"field_resource": "storage",
		"field_units":    "3",
	})
	if err == nil {
		t.Fatalf("storage promise beyond local capacity should fail")
	}
}

func TestBrokenPromiseStakeCostReducesBudget(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	alice.applyBrokenPromiseCost("bob", map[string]string{"field_stake": "2"}, "test broken promise")
	if alice.budget != 3 {
		t.Fatalf("budget after stake cost = %d, want 3", alice.budget)
	}
}

func TestRunTurnRepairsMissingActAndTarget(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], missingShapeDecider{}, decision.FakeMonitor{})
	if err := alice.runTurn(context.Background(), 0); err != nil {
		t.Fatalf("run turn: %v", err)
	}
	if !hasEvent(alice.events, "decision_repaired") {
		t.Fatalf("expected decision_repaired event, got %#v", alice.events)
	}
}

func TestCandidateDiscoveryCanFormDirectPeer(t *testing.T) {
	cfg := candidateOnlyTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	discoveryFields := map[string]string{"field_promise_about": decision.PromiseAboutLinkDiscovery}
	if !alice.canDialTarget("bob", discoveryFields) {
		t.Fatalf("candidate discovery should be dialable")
	}
	alice.observeOutcome("bob", relationship.OutcomeDiscoveryKept)
	if !alice.canDial("bob") {
		t.Fatalf("kept candidate discovery should form a direct peer")
	}
	if !hasEvent(alice.events, string(relationship.TransitionAdded)) {
		t.Fatalf("expected direct_peer_added event, got %#v", alice.events)
	}
}

func TestProtocolForFieldsRoutesShippingPromises(t *testing.T) {
	cfg := shippingTestConfig(t)
	fulfillment := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	protocolName, _ := fulfillment.protocolForFields(map[string]string{"field_promise_about": production.PromiseWeighPackage})
	if protocolName != pcid.PostalScaleV1 {
		t.Fatalf("weigh package protocol = %s, want %s", protocolName, pcid.PostalScaleV1)
	}
	protocolName, _ = fulfillment.protocolForFields(map[string]string{"field_promise_about": production.PromisePrintLabel})
	if protocolName != pcid.UPSLabelV1 {
		t.Fatalf("print label protocol = %s, want %s", protocolName, pcid.UPSLabelV1)
	}
	protocolName, _ = fulfillment.protocolForFields(map[string]string{"field_promise_about": production.PromiseShipmentUpdate})
	if protocolName != pcid.AccountingV1 {
		t.Fatalf("shipment update protocol = %s, want %s", protocolName, pcid.AccountingV1)
	}
}

func TestDeterministicShippingHandlersReturnEvidence(t *testing.T) {
	cfg := shippingTestConfig(t)
	scale := NewNode(cfg, cfg.Agents[1], &decision.FakeDecider{}, decision.FakeMonitor{})
	message := parsedMessage{
		ProtocolName: pcid.PostalScaleV1,
		Fields: map[string]string{
			"from":                "fulfillment",
			"field_promise_about": production.PromiseWeighPackage,
			"field_package_id":    "PKG-1001",
		},
	}
	ackFields, err := scale.handleProtocolPromise(message)
	if err != nil {
		t.Fatalf("scale handler: %v", err)
	}
	if ackFields["field_weight_ounces"] == "" || !hasEvent(scale.events, "package_weighed") {
		t.Fatalf("scale did not return weight evidence: %#v events %#v", ackFields, scale.events)
	}
}

// TestFulfillmentStartupWorkflowStepsUseDeterministicHandlers checks the
// deterministic handler sequence behind the live Docker workflow.
// Intent: Unit tests cannot open local TCP sockets in the Codex sandbox, so the
// handler-level test preserves the evidence chain while Docker Compose remains
// the live TCP validation path. Source: DI-parok
func TestFulfillmentStartupWorkflowStepsUseDeterministicHandlers(t *testing.T) {
	cfg := shippingTestConfig(t)
	fulfillment := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	accounting := NewNode(cfg, cfg.Agents[3], &decision.FakeDecider{}, decision.FakeMonitor{})
	addressAck, err := accounting.handleAccountingPromise(map[string]string{
		"from":                "fulfillment",
		"field_promise_about": production.PromiseAddressLookup,
		"field_order_id":      fulfillmentOrderID,
	})
	if err != nil {
		t.Fatalf("address lookup: %v", err)
	}
	scale := NewNode(cfg, cfg.Agents[1], &decision.FakeDecider{}, decision.FakeMonitor{})
	weightAck, err := scale.handlePostalScalePromise(map[string]string{
		"from":                "fulfillment",
		"field_promise_about": production.PromiseWeighPackage,
		"field_package_id":    fulfillmentPackageID,
	})
	if err != nil {
		t.Fatalf("package weighing: %v", err)
	}
	printer := NewNode(cfg, cfg.Agents[2], &decision.FakeDecider{}, decision.FakeMonitor{})
	labelAck, err := printer.handleUPSLabelPromise(map[string]string{
		"from":                   "fulfillment",
		"field_promise_about":    production.PromisePrintLabel,
		"field_package_id":       fulfillmentPackageID,
		"field_shipping_address": addressAck["field_shipping_address"],
		"field_weight_ounces":    weightAck["field_weight_ounces"],
	})
	if err != nil {
		t.Fatalf("label printing: %v", err)
	}
	_, err = accounting.handleAccountingPromise(map[string]string{
		"from":                  "fulfillment",
		"field_promise_about":   production.PromiseShipmentUpdate,
		"field_order_id":        fulfillmentOrderID,
		"field_tracking_number": labelAck["field_tracking_number"],
		"field_cost_cents":      labelAck["field_cost_cents"],
	})
	if err != nil {
		t.Fatalf("accounting update: %v", err)
	}
	fulfillment.record("fulfillment_workflow_completed", "kept", "accounting", "test")
	allReceiverEvents := append(append(accounting.events, scale.events...), printer.events...)
	for _, eventName := range []string{"shipping_address_promised", "package_weighed", "shipping_label_printed", "accounting_updated"} {
		if !hasEvent(allReceiverEvents, eventName) {
			t.Fatalf("receiver logs missing %s event: %#v", eventName, allReceiverEvents)
		}
	}
	if !hasEvent(fulfillment.events, "fulfillment_workflow_completed") {
		t.Fatalf("fulfillment did not record workflow completion: %#v", fulfillment.events)
	}
}

func TestSupportedProtocolIsLocalToAgent(t *testing.T) {
	cfg := shippingTestConfig(t)
	fulfillment := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	scale := NewNode(cfg, cfg.Agents[1], &decision.FakeDecider{}, decision.FakeMonitor{})
	if fulfillment.supportsProtocol(pcid.AccountingV1) {
		t.Fatalf("fulfillment should send accounting pCID but should not receive it")
	}
	if scale.supportsProtocol(pcid.AccountingV1) {
		t.Fatalf("postal scale should not support accounting pCID")
	}
}

func TestShutdownGraceRecordsBeforeDone(t *testing.T) {
	cfg := singleNodeTestConfig(t)
	cfg.ShutdownGraceMillis = 1
	node := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := node.Run(ctx); err != nil {
		t.Fatalf("run node: %v", err)
	}
	graceIndex := eventIndex(node.events, "shutdown_grace_elapsed")
	doneIndex := eventIndex(node.events, "node_done")
	if graceIndex < 0 || doneIndex < 0 || graceIndex > doneIndex {
		t.Fatalf("shutdown grace should be recorded before node_done: %#v", node.events)
	}
}

func TestDrainInflightRecordsOnce(t *testing.T) {
	cfg := singleNodeTestConfig(t)
	node := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	node.drainInflight(context.Background())
	node.drainInflight(context.Background())
	if countEvents(node.events, "inflight_drained") != 1 {
		t.Fatalf("drain should be recorded once: %#v", node.events)
	}
}

func singleNodeTestConfig(t *testing.T) config.Config {
	t.Helper()
	return config.Config{
		RunID:                 "test",
		RunRoot:               filepath.Join(t.TempDir(), "run"),
		ListenPort:            0,
		ProviderBaseURL:       "https://example.invalid/v1/responses",
		APIKeyEnv:             "OPENAI_API_KEY",
		AgentModel:            "model",
		MonitorModel:          "model",
		ReasoningEffort:       "medium",
		ServiceTier:           "default",
		RequestTimeoutSeconds: 1,
		StartupDelayMillis:    1,
		TurnDelayMillis:       1,
		MaxTurns:              1,
		MaxAgentCalls:         1,
		MaxMonitorCalls:       1,
		MonitorNode:           "alice",
		StrongTrustThreshold:  2,
		WeakTrustThreshold:    -2,
		TrustDecayPerRound:    0,
		Agents: []config.AgentConfig{{
			Name:       "alice",
			Persona:    "tester",
			Motivation: "test",
			Budget:     1,
			Capacity:   1,
		}},
		Containers: []config.ContainerConfig{{Name: "alice", Agents: []string{"alice"}}},
	}
}

func candidateOnlyTestConfig(t *testing.T) config.Config {
	cfg := twoNodeTestConfig(t)
	cfg.Agents[0].InitialPeers = nil
	cfg.Agents[0].CandidatePeers = []string{"bob"}
	cfg.Agents[1].InitialPeers = nil
	cfg.Agents[1].CandidatePeers = []string{"alice"}
	return cfg
}

func twoNodeTestConfig(t *testing.T) config.Config {
	t.Helper()
	return config.Config{
		RunID:                 "test",
		RunRoot:               filepath.Join(t.TempDir(), "run"),
		ListenPort:            0,
		ProviderBaseURL:       "https://example.invalid/v1/responses",
		APIKeyEnv:             "OPENAI_API_KEY",
		AgentModel:            "model",
		MonitorModel:          "model",
		ReasoningEffort:       "medium",
		ServiceTier:           "default",
		RequestTimeoutSeconds: 1,
		StartupDelayMillis:    1,
		TurnDelayMillis:       1,
		MaxTurns:              1,
		MaxAgentCalls:         1,
		MaxMonitorCalls:       1,
		MonitorNode:           "alice",
		StrongTrustThreshold:  2,
		WeakTrustThreshold:    -2,
		TrustDecayPerRound:    0,
		Agents: []config.AgentConfig{{
			Name:           "alice",
			Persona:        "tester",
			Motivation:     "test",
			InitialPeers:   []string{"bob"},
			CandidatePeers: []string{"bob"},
			Budget:         5,
			Capacity:       1,
		}, {
			Name:           "bob",
			Persona:        "tester",
			Motivation:     "test",
			InitialPeers:   []string{"alice"},
			CandidatePeers: []string{"alice"},
			Budget:         5,
			Capacity:       1,
		}},
		Containers: []config.ContainerConfig{
			{Name: "alice", Agents: []string{"alice"}},
			{Name: "bob", Agents: []string{"bob"}},
		},
	}
}

func shippingTestConfig(t *testing.T) config.Config {
	t.Helper()
	return config.Config{
		RunID:                 "shipping-test",
		RunRoot:               filepath.Join(t.TempDir(), "run"),
		ListenPort:            0,
		ProviderBaseURL:       "https://example.invalid/v1/responses",
		APIKeyEnv:             "OPENAI_API_KEY",
		AgentModel:            "model",
		MonitorModel:          "model",
		ReasoningEffort:       "medium",
		ServiceTier:           "default",
		RequestTimeoutSeconds: 1,
		StartupDelayMillis:    1,
		TurnDelayMillis:       1,
		MaxTurns:              1,
		MaxAgentCalls:         1,
		MaxMonitorCalls:       1,
		MonitorNode:           "fulfillment",
		StrongTrustThreshold:  2,
		WeakTrustThreshold:    -2,
		TrustDecayPerRound:    0,
		Agents: []config.AgentConfig{{
			Name:           "fulfillment",
			Kind:           "fulfillment",
			Persona:        "workflow",
			Motivation:     "ship package",
			InitialPeers:   []string{"postal_scale", "ups_label_printer", "accounting"},
			CandidatePeers: []string{"postal_scale", "ups_label_printer", "accounting"},
			SupportedPCIDs: []string{pcid.RelationshipV1},
			Budget:         5,
			Capacity:       5,
		}, {
			Name:           "postal_scale",
			Kind:           "postal_scale",
			Persona:        "scale",
			Motivation:     "weigh",
			InitialPeers:   []string{"fulfillment"},
			CandidatePeers: []string{"fulfillment"},
			SupportedPCIDs: []string{pcid.RelationshipV1, pcid.PostalScaleV1},
			Budget:         5,
			Capacity:       5,
		}, {
			Name:           "ups_label_printer",
			Kind:           "ups_label_printer",
			Persona:        "printer",
			Motivation:     "print label",
			InitialPeers:   []string{"fulfillment"},
			CandidatePeers: []string{"fulfillment"},
			SupportedPCIDs: []string{pcid.RelationshipV1, pcid.UPSLabelV1},
			Budget:         5,
			Capacity:       5,
		}, {
			Name:           "accounting",
			Kind:           "accounting",
			Persona:        "accounting",
			Motivation:     "records",
			InitialPeers:   []string{"fulfillment"},
			CandidatePeers: []string{"fulfillment"},
			SupportedPCIDs: []string{pcid.RelationshipV1, pcid.AccountingV1},
			Budget:         5,
			Capacity:       5,
		}},
		Containers: []config.ContainerConfig{
			{Name: "shipping", Agents: []string{"fulfillment", "postal_scale", "ups_label_printer", "accounting"}},
		},
	}
}

type missingShapeDecider struct{}

func (missingShapeDecider) Decide(_ context.Context, _ decision.Observation) (decision.PromiseDecision, error) {
	return decision.PromiseDecision{
		Promise: "Alice promises one bounded local exchange.",
		Fields:  map[string]any{"purpose": "repair-test"},
	}, nil
}

func hasEvent(events []decision.Event, eventName string) bool {
	for _, event := range events {
		if event.Event == eventName {
			return true
		}
	}
	return false
}

func eventIndex(events []decision.Event, eventName string) int {
	for eventIndex, event := range events {
		if event.Event == eventName {
			return eventIndex
		}
	}
	return -1
}

func countEvents(events []decision.Event, eventName string) int {
	count := 0
	for _, event := range events {
		if event.Event == eventName {
			count++
		}
	}
	return count
}
