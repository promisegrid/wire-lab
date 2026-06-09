package runtime

import (
	"context"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/config"
	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/decision"
	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/pcid"
	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/production"
	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/protocol"
	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/relationship"
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
	if !hasEvent(node.events, "app_kernel_registration_skipped") || !hasEvent(node.events, "local_non_commitment") {
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

func TestProviderDecisionErrorRecordsLocalNonCommitment(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	node := NewNode(cfg, cfg.Agents[0], errorDecider{}, decision.FakeMonitor{})
	if err := node.runTurn(context.Background(), 0); err != nil {
		node.recordDecisionError(err)
	}
	for _, event := range node.events {
		if event.Event == "decision_error" && event.Outcome == "non_commitment" {
			return
		}
	}
	t.Fatalf("provider decision error should be local non-commitment evidence: %#v", node.events)
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

func TestLocalResourceExhaustionDoesNotTouchPeerTrust(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	cfg.Agents[0].Capacity = 0
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	if err := alice.runTurn(context.Background(), 0); err != nil {
		t.Fatalf("run turn: %v", err)
	}
	if alice.ledger.Trust("bob") != 0 {
		t.Fatalf("local exhaustion changed peer trust to %d, want 0", alice.ledger.Trust("bob"))
	}
	if hasEvent(alice.events, string(relationship.TransitionUnchanged)) {
		t.Fatalf("local exhaustion should not create peer trust transition evidence: %#v", alice.events)
	}
	if !hasEvent(alice.events, "local_resource_exhausted") {
		t.Fatalf("local exhaustion should be recorded locally: %#v", alice.events)
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

func TestNotPromisedAckStaysNonCommitment(t *testing.T) {
	outcome, updatesPeerTrust := outcomeForSendError(ackOutcomeError{outcome: "not_promised"})
	if outcome != relationship.OutcomeNonCommitment {
		t.Fatalf("not_promised outcome = %s, want %s", outcome, relationship.OutcomeNonCommitment)
	}
	if updatesPeerTrust {
		t.Fatalf("not_promised should not update peer trust")
	}
	eventName, eventOutcome := sendEventForError(ackOutcomeError{outcome: "not_promised"})
	if eventName != "send_not_promised" || eventOutcome != "non_commitment" {
		t.Fatalf("not_promised send event = %s/%s, want send_not_promised/non_commitment", eventName, eventOutcome)
	}
}

func TestLocalSendFailureDoesNotUpdatePeerTrust(t *testing.T) {
	outcome, updatesPeerTrust := outcomeForSendError(context.DeadlineExceeded)
	if outcome != relationship.OutcomeNonCommitment {
		t.Fatalf("local send failure outcome = %s, want %s", outcome, relationship.OutcomeNonCommitment)
	}
	if updatesPeerTrust {
		t.Fatalf("local send failure should not update peer trust")
	}
	eventName, eventOutcome := sendEventForError(context.DeadlineExceeded)
	if eventName != "send_unavailable" || eventOutcome != "non_commitment" {
		t.Fatalf("local send failure event = %s/%s, want send_unavailable/non_commitment", eventName, eventOutcome)
	}
}

func TestRepeatedPromiseSuppressedAfterJournalEvidence(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	fields := map[string]string{
		"promise":             "I promise to send one repeated relationship promise.",
		"field_promise_about": "local_observation",
	}
	recordKey := alice.rememberOutstandingPromise("bob", pcid.RelationshipV1, "exact", fields)
	alice.resolveOutstandingPromise(recordKey, promiseStatusKept, "test")
	if !alice.suppressRepeatedPromise("bob", fields) {
		t.Fatalf("expected repeated promise to be suppressed")
	}
	if !hasEvent(alice.events, "promise_repeated_suppressed") {
		t.Fatalf("suppressed repeat should be visible in local evidence: %#v", alice.events)
	}
}

func TestNotPromisedSuppressesSemanticRetryWithoutTrustChange(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	fields := map[string]string{
		"promise":             "I promise to make one storage offer.",
		"field_promise_about": "storage_offer",
	}
	alice.rememberNonCommitment("bob", pcid.RelationshipV1, fields, "ack outcome not_promised")
	retryFields := map[string]string{
		"promise":             "I promise to make a revised storage offer.",
		"field_promise_about": "storage_offer",
	}
	if !alice.shouldSuppressNonCommittedPromise("bob", retryFields) {
		t.Fatalf("semantic retry after not_promised should be suppressed")
	}
	if alice.ledger.Trust("bob") != 0 {
		t.Fatalf("not_promised suppression changed peer trust to %d, want 0", alice.ledger.Trust("bob"))
	}
	if !hasEvent(alice.events, "promise_not_promised_suppressed") {
		t.Fatalf("suppressed not_promised retry should be visible in local evidence: %#v", alice.events)
	}
}

func TestCheckpointJournalIsReusableBeyondShipmentMap(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	record := checkpointRecord{
		Key:          checkpointKey(pcid.PrinterPortV1, production.PromiseRedeemPrintCapability, "spool-1"),
		ProtocolName: pcid.PrinterPortV1,
		PromiseAbout: production.PromiseRedeemPrintCapability,
		Subject:      "spool-1",
		Detail:       "printer port checkpoint regression test",
	}
	if alice.rememberCheckpoint(record) {
		t.Fatalf("first checkpoint record should be new")
	}
	if !alice.rememberCheckpoint(record) {
		t.Fatalf("second checkpoint record should be recognized as duplicate")
	}
	if len(alice.checkpointJournal) != 1 {
		t.Fatalf("checkpoint journal length = %d, want 1", len(alice.checkpointJournal))
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
	protocolName, _ = fulfillment.protocolForFields(map[string]string{"field_promise_about": production.PromiseIssuePrintCapability})
	if protocolName != pcid.PrinterPortV1 {
		t.Fatalf("print capability protocol = %s, want %s", protocolName, pcid.PrinterPortV1)
	}
}

func TestAutonomousProtocolMismatchReframesToRelationship(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	fields := map[string]string{
		"field_protocol":      pcid.CASStorageV1,
		"field_promise_about": "storage_capacity",
		"promise":             "Alice promises only local storage-capacity evidence.",
	}
	alice.normalizeAutonomousPromiseFields("bob", fields)
	if fields["field_protocol"] != pcid.RelationshipV1 {
		t.Fatalf("reframed protocol = %s, want %s", fields["field_protocol"], pcid.RelationshipV1)
	}
	if fields["field_promise_about"] != "local_observation" {
		t.Fatalf("reframed promise_about = %s, want local_observation", fields["field_promise_about"])
	}
	if fields["field_original_protocol"] != pcid.CASStorageV1 || fields["field_original_promise_about"] != "storage_capacity" {
		t.Fatalf("original protocol fields missing after reframe: %#v", fields)
	}
	if !hasEvent(alice.events, "promise_reframed_for_pcid_fit") {
		t.Fatalf("pCID reframe should be visible in local evidence: %#v", alice.events)
	}
}

func TestAutonomousConcreteCASPayloadKeepsProtocol(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	contentBytes := production.SampleContentBytes()
	fields := map[string]string{
		"field_protocol":      pcid.CASStorageV1,
		"field_promise_about": production.PromiseStoreContent,
		"field_content_cid":   production.ContentCID(contentBytes),
		"field_content_b64":   "nonempty-test-bytes",
		"promise":             "Alice promises to receive concrete CAS storage evidence.",
	}
	alice.normalizeAutonomousPromiseFields("bob", fields)
	if fields["field_protocol"] != pcid.CASStorageV1 {
		t.Fatalf("concrete CAS payload protocol = %s, want %s", fields["field_protocol"], pcid.CASStorageV1)
	}
	if hasEvent(alice.events, "promise_reframed_for_pcid_fit") {
		t.Fatalf("concrete CAS payload should not be reframed: %#v", alice.events)
	}
}

func TestNegativeAckVerdictsDoNotUpdateTrust(t *testing.T) {
	negativeVerdicts := []map[string]string{
		{"field_verdict": "broken"},
		{"field_verdict": "disagree"},
		{"field_variant_status": "not_promised"},
		{"field_storage_status": "price_refused"},
		{"field_compute_status": "capacity_refused"},
		{"field_cache_status": "miss"},
	}
	for _, ackFields := range negativeVerdicts {
		if evidenceUpdatesTrust(ackFields) {
			t.Fatalf("negative ACK fields should not update trust: %#v", ackFields)
		}
	}
	if !evidenceUpdatesTrust(map[string]string{"field_verdict": "kept"}) {
		t.Fatalf("kept ACK verdict should update trust")
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
	accounting := NewNode(cfg, cfg.Agents[4], &decision.FakeDecider{}, decision.FakeMonitor{})
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
	trackingNumber, costCents, err := production.LabelForShipment(fulfillmentPackageID, addressAck["field_shipping_address"], intField(weightAck, "field_weight_ounces"))
	if err != nil {
		t.Fatalf("label facts: %v", err)
	}
	labelBytes, err := production.LabelBytesForShipment(map[string]string{
		"field_package_id":      fulfillmentPackageID,
		"field_tracking_number": trackingNumber,
		"field_cost_cents":      strconv.Itoa(costCents),
	})
	if err != nil {
		t.Fatalf("label bytes: %v", err)
	}
	printerPort := NewNode(cfg, cfg.Agents[3], &decision.FakeDecider{}, decision.FakeMonitor{})
	capabilityFields := map[string]string{
		"from":                             "ups_label_printer",
		"to":                               "printer_port",
		"field_promise_about":              production.PromiseIssuePrintCapability,
		"field_print_capability_issuee":    "ups_label_printer",
		"field_print_capability_token_id":  "printcap-ups_label_printer",
		"field_print_capability_scope":     production.PrintCapabilityScope,
		"field_print_capability_max_bytes": strconv.Itoa(production.PrintCapabilityMaxBytes),
	}
	capabilityAck, err := printerPort.handlePrinterPortPromise(capabilityFields)
	if err != nil {
		t.Fatalf("capability issue: %v", err)
	}
	redemptionFields := map[string]string{
		"from":                             "ups_label_printer",
		"field_promise_about":              production.PromiseRedeemPrintCapability,
		"field_print_capability_issuee":    "ups_label_printer",
		"field_print_capability_token":     capabilityAck["field_print_capability_token"],
		"field_print_capability_token_id":  capabilityAck["field_print_capability_token_id"],
		"field_print_capability_scope":     capabilityAck["field_print_capability_scope"],
		"field_print_capability_max_bytes": capabilityAck["field_print_capability_max_bytes"],
		"field_label_bytes_hex":            hex.EncodeToString(labelBytes),
	}
	printAck, err := printerPort.handlePrinterPortPromise(redemptionFields)
	if err != nil {
		t.Fatalf("capability redemption: %v", err)
	}
	_, err = accounting.handleAccountingPromise(map[string]string{
		"from":                  "fulfillment",
		"field_promise_about":   production.PromiseShipmentUpdate,
		"field_order_id":        fulfillmentOrderID,
		"field_tracking_number": trackingNumber,
		"field_cost_cents":      strconv.Itoa(costCents),
	})
	if err != nil {
		t.Fatalf("accounting update: %v", err)
	}
	if printAck["field_printer_spool_id"] == "" {
		t.Fatalf("printer port did not return spool evidence: %#v", printAck)
	}
	fulfillment.record("fulfillment_workflow_completed", "kept", "accounting", "test")
	allReceiverEvents := append(append(accounting.events, scale.events...), printerPort.events...)
	for _, eventName := range []string{"shipping_address_promised", "package_weighed", "printer_capability_issued", "printer_port_printed", "accounting_updated"} {
		if !hasEvent(allReceiverEvents, eventName) {
			t.Fatalf("receiver logs missing %s event: %#v", eventName, allReceiverEvents)
		}
	}
	if !hasEvent(fulfillment.events, "fulfillment_workflow_completed") {
		t.Fatalf("fulfillment did not record workflow completion: %#v", fulfillment.events)
	}
}

func TestDuplicateAccountingUpdateDoesNotRepeatTrust(t *testing.T) {
	cfg := shippingTestConfig(t)
	accounting := NewNode(cfg, cfg.Agents[4], &decision.FakeDecider{}, decision.FakeMonitor{})
	frameBytes := signedAccountingUpdateFrame(t, accounting)
	for attempt := 0; attempt < 2; attempt++ {
		ackBytes, err := accounting.handleFrame(frameBytes)
		if err != nil {
			t.Fatalf("handle accounting update attempt %d: %v", attempt+1, err)
		}
		ackEnvelope, parseErr := protocol.ParseEnvelope(ackBytes)
		if parseErr != nil {
			t.Fatalf("parse accounting ack attempt %d: %v", attempt+1, parseErr)
		}
		ackFields, fieldsErr := ackEnvelope.PayloadFields()
		if fieldsErr != nil {
			t.Fatalf("parse accounting ack fields attempt %d: %v", attempt+1, fieldsErr)
		}
		if attempt == 0 && ackFields[duplicateShipmentEvidenceField] == "true" {
			t.Fatalf("first accounting update should not be duplicate: %#v", ackFields)
		}
		if attempt == 1 && ackFields[duplicateShipmentEvidenceField] != "true" {
			t.Fatalf("second accounting update should be duplicate: %#v", ackFields)
		}
	}
	if accounting.ledger.Trust("fulfillment") != 1 {
		t.Fatalf("duplicate accounting trust = %d, want 1", accounting.ledger.Trust("fulfillment"))
	}
	if countEvents(accounting.events, "accounting_updated") != 1 {
		t.Fatalf("accounting update should be recorded once: %#v", accounting.events)
	}
	if countEvents(accounting.events, "accounting_update_duplicate") != 1 {
		t.Fatalf("duplicate accounting update should be recorded once: %#v", accounting.events)
	}
	if len(accounting.checkpointJournal) != 1 {
		t.Fatalf("accounting checkpoint journal length = %d, want 1", len(accounting.checkpointJournal))
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
	printerPort := NewNode(cfg, cfg.Agents[3], &decision.FakeDecider{}, decision.FakeMonitor{})
	if !printerPort.supportsProtocol(pcid.PrinterPortV1) {
		t.Fatalf("printer_port should support printer_port pCID")
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

func TestMonitorFailureWritesNonAuthoritativeDoneMarker(t *testing.T) {
	cfg := singleNodeTestConfig(t)
	node := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, failingMonitor{})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := node.Run(ctx); err != nil {
		t.Fatalf("run node with failing monitor: %v", err)
	}
	if _, err := os.Stat(filepath.Join(cfg.RunRoot, cfg.RunID, "monitor.done")); err != nil {
		t.Fatalf("monitor done marker should exist after observer failure: %v", err)
	}
	if !hasEventOutcome(node.events, "monitor_error", "non_commitment") {
		t.Fatalf("monitor error should be local non-commitment evidence: %#v", node.events)
	}
	if !hasEventOutcome(node.events, "monitor_done", "non_commitment") {
		t.Fatalf("fallback monitor marker should be non-authoritative evidence: %#v", node.events)
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
			InitialPeers:   []string{"fulfillment", "printer_port"},
			CandidatePeers: []string{"fulfillment", "printer_port"},
			SupportedPCIDs: []string{pcid.RelationshipV1, pcid.UPSLabelV1},
			Budget:         5,
			Capacity:       5,
		}, {
			Name:           "printer_port",
			Kind:           "printer_port",
			Persona:        "printer port",
			Motivation:     "local hardware access",
			InitialPeers:   []string{"ups_label_printer"},
			CandidatePeers: []string{"ups_label_printer"},
			SupportedPCIDs: []string{pcid.RelationshipV1, pcid.PrinterPortV1},
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
			{Name: "shipping", Agents: []string{"fulfillment", "postal_scale", "ups_label_printer", "printer_port", "accounting"}},
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

type errorDecider struct{}

func (errorDecider) Decide(_ context.Context, _ decision.Observation) (decision.PromiseDecision, error) {
	return decision.PromiseDecision{}, context.DeadlineExceeded
}

type failingMonitor struct{}

func (failingMonitor) Evaluate(_ context.Context, _ []decision.Event) (decision.MonitorReport, error) {
	return decision.MonitorReport{}, errors.New("monitor provider unavailable")
}

func signedAccountingUpdateFrame(t *testing.T, node *Node) []byte {
	t.Helper()
	envelope, err := protocol.NewEnvelope(node.Protocols.MustCID(pcid.AccountingV1), map[string]string{
		"act":                   decision.ActPromise,
		"from":                  "fulfillment",
		"to":                    "accounting",
		"turn":                  "test",
		"promise":               "I promise to receive accounting's shipment checkpoint evidence for this order and tracking number.",
		"reason":                "duplicate checkpoint regression test",
		"field_promise_about":   production.PromiseShipmentUpdate,
		"field_order_id":        fulfillmentOrderID,
		"field_tracking_number": "1Z999AA10123456784",
		"field_cost_cents":      "1776",
	}, "fulfillment")
	if err != nil {
		t.Fatalf("new accounting update envelope: %v", err)
	}
	frameBytes, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("accounting update envelope bytes: %v", err)
	}
	return frameBytes
}

func hasEvent(events []decision.Event, eventName string) bool {
	for _, event := range events {
		if event.Event == eventName {
			return true
		}
	}
	return false
}

func hasEventOutcome(events []decision.Event, eventName, outcome string) bool {
	for _, event := range events {
		if event.Event == eventName && event.Outcome == outcome {
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
