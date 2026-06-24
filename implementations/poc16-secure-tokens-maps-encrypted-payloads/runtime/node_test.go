package runtime

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/config"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/decision"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/pcid"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/production"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/relationship"
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
		t.Fatalf("node should record local events")
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
	t.Fatalf("provider decision error should be local non-commitment event: %#v", node.events)
}

func TestResourcePromiseChecksLocalCapacity(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	err := alice.checkLocalResourcePromise(map[string]string{
		"resource": "storage",
		"units":    "3",
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
		t.Fatalf("local exhaustion should not create peer trust transition event: %#v", alice.events)
	}
	if !hasEvent(alice.events, "local_resource_exhausted") {
		t.Fatalf("local exhaustion should be recorded locally: %#v", alice.events)
	}
}

func TestBrokenPromiseStakeCostReducesBudget(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	alice.applyBrokenPromiseCost("bob", map[string]string{"stake": "2"}, "test broken promise")
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

func TestRepeatedPromiseSuppressedAfterJournalEvent(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	fields := map[string]string{
		"promise":       "I promise to send one repeated relationship promise.",
		"promise_about": "local_observation",
	}
	recordKey := alice.rememberOutstandingPromise("bob", pcid.RelationshipV1, "exact", fields)
	alice.resolveOutstandingPromise(recordKey, promiseStatusKept, "test")
	if !alice.suppressRepeatedPromise("bob", fields) {
		t.Fatalf("expected repeated promise to be suppressed")
	}
	if !hasEvent(alice.events, "promise_repeated_suppressed") {
		t.Fatalf("suppressed repeat should be visible in local events: %#v", alice.events)
	}
}

func TestNotPromisedSuppressesSemanticRetryWithoutTrustChange(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	fields := map[string]string{
		"promise":       "I promise to make one storage offer.",
		"promise_about": "storage_offer",
	}
	alice.rememberNonCommitment("bob", pcid.RelationshipV1, fields, "ack outcome not_promised")
	retryFields := map[string]string{
		"promise":       "I promise to make a revised storage offer.",
		"promise_about": "storage_offer",
	}
	if !alice.shouldSuppressNonCommittedPromise("bob", retryFields) {
		t.Fatalf("semantic retry after not_promised should be suppressed")
	}
	if alice.ledger.Trust("bob") != 0 {
		t.Fatalf("not_promised suppression changed peer trust to %d, want 0", alice.ledger.Trust("bob"))
	}
	if !hasEvent(alice.events, "promise_not_promised_suppressed") {
		t.Fatalf("suppressed not_promised retry should be visible in local events: %#v", alice.events)
	}
}

func TestCheckpointJournalIsReusableBeyondShipmentMap(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	record := checkpointRecord{
		Key:          checkpointKey(pcid.ProductionShippingV1, production.PromiseRedeemPrintCapability, "spool-1"),
		ProtocolName: pcid.ProductionShippingV1,
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
	discoveryFields := map[string]string{"promise_about": decision.PromiseAboutLinkDiscovery}
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

func TestFutureRepairCandidatePromiseCanBeHeardAfterMalformedEvent(t *testing.T) {
	// Intent: A peer may choose to hear a narrow future-repair promise after
	// malformed event records without treating that as permission for arbitrary
	// traffic or immediate trust repair. Source: DI-fijov
	cfg := computeRoutingTestConfig(t)
	grace := NewNode(cfg, cfg.Agents[2], &decision.FakeDecider{}, decision.FakeMonitor{})
	mallory := NewNode(cfg, cfg.Agents[4], &decision.FakeDecider{}, decision.FakeMonitor{})
	grace.observeOutcome("mallory", relationship.OutcomeMalformed)
	mallory.observeOutcome("grace", relationship.OutcomeMalformed)
	ordinaryFields := map[string]string{"promise_about": "ordinary_followup"}
	if grace.canAcceptFrom("mallory", ordinaryFields) {
		t.Fatalf("Grace should not accept arbitrary candidate traffic after malformed event")
	}
	if mallory.canDialTarget("grace", ordinaryFields) {
		t.Fatalf("Mallory should not dial arbitrary candidate traffic after malformed event")
	}
	repairFields := map[string]string{"promise_about": production.PromiseLabelFutureMalformedReport}
	if !grace.canAcceptFrom("mallory", repairFields) {
		t.Fatalf("Grace should be able to hear a narrow future-repair promise")
	}
	if !mallory.canDialTarget("grace", repairFields) {
		t.Fatalf("Mallory should be able to offer a narrow future-repair promise")
	}
}

func TestProtocolForFieldsRoutesShippingPromises(t *testing.T) {
	cfg := shippingTestConfig(t)
	fulfillment := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	protocolName, _ := fulfillment.protocolForFields(map[string]string{"promise_about": production.PromiseWeighPackage})
	if protocolName != pcid.ProductionShippingV1 {
		t.Fatalf("weigh package protocol = %s, want %s", protocolName, pcid.ProductionShippingV1)
	}
	protocolName, _ = fulfillment.protocolForFields(map[string]string{"promise_about": production.PromisePrintLabel})
	if protocolName != pcid.ProductionShippingV1 {
		t.Fatalf("print label protocol = %s, want %s", protocolName, pcid.ProductionShippingV1)
	}
	protocolName, _ = fulfillment.protocolForFields(map[string]string{"promise_about": production.PromiseShipmentUpdate})
	if protocolName != pcid.ProductionShippingV1 {
		t.Fatalf("shipment update protocol = %s, want %s", protocolName, pcid.ProductionShippingV1)
	}
	protocolName, _ = fulfillment.protocolForFields(map[string]string{"promise_about": production.PromiseIssuePrintCapability})
	if protocolName != pcid.ProductionShippingV1 {
		t.Fatalf("print capability protocol = %s, want %s", protocolName, pcid.ProductionShippingV1)
	}
}

func TestAutonomousProtocolMismatchReframesToRelationship(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	fields := map[string]string{
		"protocol":      pcid.CASStorageV1,
		"promise_about": "storage_capacity",
		"promise":       "Alice promises only local storage-capacity event.",
	}
	alice.normalizeAutonomousPromiseFields("bob", fields)
	if fields["protocol"] != pcid.RelationshipV1 {
		t.Fatalf("reframed protocol = %s, want %s", fields["protocol"], pcid.RelationshipV1)
	}
	if fields["promise_about"] != "local_observation" {
		t.Fatalf("reframed promise_about = %s, want local_observation", fields["promise_about"])
	}
	if fields["original_protocol"] != pcid.CASStorageV1 || fields["original_promise_about"] != "storage_capacity" {
		t.Fatalf("original protocol fields missing after reframe: %#v", fields)
	}
	if !hasEvent(alice.events, "promise_reframed_for_pcid_fit") {
		t.Fatalf("pCID reframe should be visible in local events: %#v", alice.events)
	}
}

func TestAutonomousConcreteCASPayloadKeepsProtocol(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	contentBytes := production.SampleContentBytes()
	fields := map[string]string{
		"protocol":      pcid.CASStorageV1,
		"promise_about": production.PromiseStoreContent,
		"content_cid":   production.ContentCID(contentBytes),
		"content_b64":   "nonempty-test-bytes",
		"promise":       "Alice promises to receive concrete CAS storage event.",
	}
	alice.normalizeAutonomousPromiseFields("bob", fields)
	if fields["protocol"] != pcid.CASStorageV1 {
		t.Fatalf("concrete CAS payload protocol = %s, want %s", fields["protocol"], pcid.CASStorageV1)
	}
	if hasEvent(alice.events, "promise_reframed_for_pcid_fit") {
		t.Fatalf("concrete CAS payload should not be reframed: %#v", alice.events)
	}
}

func TestNegativeAckVerdictsDoNotUpdateTrust(t *testing.T) {
	negativeVerdicts := []map[string]string{
		{"verdict": "broken"},
		{"verdict": "disagree"},
		{"variant_status": "not_promised"},
		{"storage_status": "price_refused"},
		{"compute_status": "capacity_refused"},
		{"cache_status": "miss"},
	}
	for _, ackFields := range negativeVerdicts {
		if eventUpdatesTrust(ackFields) {
			t.Fatalf("negative ACK fields should not update trust: %#v", ackFields)
		}
	}
	if !eventUpdatesTrust(map[string]string{"verdict": "kept"}) {
		t.Fatalf("kept ACK verdict should update trust")
	}
}

func TestNonTrustingAckEventUsesPrecisePromiseStatus(t *testing.T) {
	// Intent: Non-mutating ACK events should distinguish true duplicate
	// checkpoints from peer non-commitments and malformed/broken verdicts.
	// Source: DI-sihuz
	cases := []struct {
		name       string
		fields     map[string]string
		wantStatus promiseStatus
	}{{
		name:       "duplicate shipment checkpoint",
		fields:     map[string]string{duplicateShipmentEventField: "true"},
		wantStatus: promiseStatusDuplicate,
	}, {
		name:       "storage price refused",
		fields:     map[string]string{"storage_status": "price_refused"},
		wantStatus: promiseStatusNonCommitment,
	}, {
		name:       "compute capacity refused",
		fields:     map[string]string{"compute_status": "capacity_refused"},
		wantStatus: promiseStatusNonCommitment,
	}, {
		name:       "cache miss",
		fields:     map[string]string{"cache_status": "miss"},
		wantStatus: promiseStatusNonCommitment,
	}, {
		name:       "future repair only",
		fields:     map[string]string{"repair_status": "future_only"},
		wantStatus: promiseStatusNonCommitment,
	}, {
		name:       "unsupported variant",
		fields:     map[string]string{"variant_status": "not_promised"},
		wantStatus: promiseStatusNonCommitment,
	}, {
		name:       "replay refused",
		fields:     map[string]string{"replay_status": "not_promised"},
		wantStatus: promiseStatusNonCommitment,
	}, {
		name:       "broken verdict",
		fields:     map[string]string{"verdict": "broken"},
		wantStatus: promiseStatusBroken,
	}, {
		name:       "malformed verdict",
		fields:     map[string]string{"verdict": "malformed"},
		wantStatus: promiseStatusMalformed,
	}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if gotStatus := promiseStatusForNonTrustingEvent(tc.fields); gotStatus != tc.wantStatus {
				t.Fatalf("status = %s, want %s for %#v", gotStatus, tc.wantStatus, tc.fields)
			}
		})
	}
}

func TestRunScopedEventSummaryCountsAllNonCommitments(t *testing.T) {
	// Intent: Saved event summaries should reflect all local
	// non-commitment outcomes, not only receiver-side not-promised restraint
	// journal entries. Source: DI-sihuz
	cfg := singleNodeTestConfig(t)
	node := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	node.record("local_non_commitment", "non_commitment", "", "no direct peer")
	node.record("promise_resolved", "non_commitment", "bob", "status=non_commitment")
	if err := node.saveRunScopedState(); err != nil {
		t.Fatalf("save run scoped state: %v", err)
	}
	if !hasEventDetail(node.events, "event_run_store_saved", "non_commitments=2") {
		t.Fatalf("saved event summary did not include all non-commitments: %#v", node.events)
	}
}

func TestDeterministicShippingHandlersReturnEvent(t *testing.T) {
	cfg := shippingTestConfig(t)
	scale := NewNode(cfg, cfg.Agents[1], &decision.FakeDecider{}, decision.FakeMonitor{})
	message := parsedMessage{
		ProtocolName: pcid.ProductionShippingV1,
		Fields: map[string]string{
			"from":          "fulfillment",
			"promise_about": production.PromiseWeighPackage,
			"package_id":    "PKG-1001",
		},
	}
	ackResult, err := scale.handleProtocolPromise(message)
	if err != nil {
		t.Fatalf("scale handler: %v", err)
	}
	if ackResult.Fields["weight_ounces"] == "" || !hasEvent(scale.events, "package_weighed") {
		t.Fatalf("scale did not return weight event: %#v events %#v", ackResult.Fields, scale.events)
	}
}

// TestFulfillmentStartupWorkflowStepsUseDeterministicHandlers checks the
// deterministic handler sequence behind the live Docker workflow.
// Intent: Unit tests cannot open local TCP sockets in the Codex sandbox, so the
// handler-level test preserves the event chain while Docker Compose remains
// the live TCP validation path. Source: DI-parok
func TestFulfillmentStartupWorkflowStepsUseDeterministicHandlers(t *testing.T) {
	cfg := shippingTestConfig(t)
	fulfillment := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	accounting := NewNode(cfg, cfg.Agents[4], &decision.FakeDecider{}, decision.FakeMonitor{})
	addressAck, err := accounting.handleAccountingPromise(map[string]string{
		"from":          "fulfillment",
		"promise_about": production.PromiseAddressLookup,
		"order_id":      fulfillmentOrderID,
	})
	if err != nil {
		t.Fatalf("address lookup: %v", err)
	}
	scale := NewNode(cfg, cfg.Agents[1], &decision.FakeDecider{}, decision.FakeMonitor{})
	weightAck, err := scale.handlePostalScalePromise(map[string]string{
		"from":          "fulfillment",
		"promise_about": production.PromiseWeighPackage,
		"package_id":    fulfillmentPackageID,
	})
	if err != nil {
		t.Fatalf("package weighing: %v", err)
	}
	trackingNumber, costCents, err := production.LabelForShipment(fulfillmentPackageID, addressAck["shipping_address"], intField(weightAck, "weight_ounces"))
	if err != nil {
		t.Fatalf("label facts: %v", err)
	}
	labelBytes, err := production.LabelBytesForShipment(map[string]string{
		"package_id":      fulfillmentPackageID,
		"tracking_number": trackingNumber,
		"cost_cents":      strconv.Itoa(costCents),
	})
	if err != nil {
		t.Fatalf("label bytes: %v", err)
	}
	printerPort := NewNode(cfg, cfg.Agents[3], &decision.FakeDecider{}, decision.FakeMonitor{})
	capabilityFields := map[string]string{
		"from":                       "ups_label_printer",
		"to":                         "printer_port",
		"promise_about":              production.PromiseIssuePrintCapability,
		"print_capability_issuee":    "ups_label_printer",
		"print_capability_token_id":  "printcap-ups_label_printer",
		"print_capability_scope":     production.PrintCapabilityScope,
		"print_capability_max_bytes": strconv.Itoa(production.PrintCapabilityMaxBytes),
	}
	capabilityAck, err := printerPort.handlePrinterPortPromise(capabilityFields)
	if err != nil {
		t.Fatalf("capability issue: %v", err)
	}
	redemptionFields := map[string]string{
		"from":                       "ups_label_printer",
		"promise_about":              production.PromiseRedeemPrintCapability,
		"print_capability_issuee":    "ups_label_printer",
		"print_capability_token":     capabilityAck["print_capability_token"],
		"print_capability_token_id":  capabilityAck["print_capability_token_id"],
		"print_capability_scope":     capabilityAck["print_capability_scope"],
		"print_capability_max_bytes": capabilityAck["print_capability_max_bytes"],
		"label_bytes_hex":            hex.EncodeToString(labelBytes),
	}
	printAck, err := printerPort.handlePrinterPortPromise(redemptionFields)
	if err != nil {
		t.Fatalf("capability redemption: %v", err)
	}
	_, err = accounting.handleAccountingPromise(map[string]string{
		"from":            "fulfillment",
		"promise_about":   production.PromiseShipmentUpdate,
		"order_id":        fulfillmentOrderID,
		"tracking_number": trackingNumber,
		"cost_cents":      strconv.Itoa(costCents),
	})
	if err != nil {
		t.Fatalf("accounting update: %v", err)
	}
	if printAck["printer_spool_id"] == "" {
		t.Fatalf("printer port did not return spool event: %#v", printAck)
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
	for attempt := 0; attempt < 2; attempt++ {
		frameBytes := signedAccountingUpdateFrame(t, accounting, "semantic-duplicate-"+strconv.Itoa(attempt+1))
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
		if attempt == 0 && ackFields[duplicateShipmentEventField] == "true" {
			t.Fatalf("first accounting update should not be duplicate: %#v", ackFields)
		}
		if attempt == 1 && ackFields[duplicateShipmentEventField] != "true" {
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
	if len(accounting.checkpointJournal) == 0 {
		t.Fatalf("accounting checkpoint journal should contain semantic duplicate event")
	}
}

func TestExactEnvelopeReplayRejectedWithoutTrustGain(t *testing.T) {
	cfg := shippingTestConfig(t)
	accounting := NewNode(cfg, cfg.Agents[4], &decision.FakeDecider{}, decision.FakeMonitor{})
	frameBytes := signedAccountingUpdateFrame(t, accounting, "exact-replay")
	if _, err := accounting.handleFrame(frameBytes); err != nil {
		t.Fatalf("first accounting frame should be accepted: %v", err)
	}
	ackBytes, err := accounting.handleFrame(frameBytes)
	if err != nil {
		t.Fatalf("exact replay should return non-commitment ACK, not a transport error: %v", err)
	}
	ackEnvelope, parseErr := protocol.ParseEnvelope(ackBytes)
	if parseErr != nil {
		t.Fatalf("parse replay ack: %v", parseErr)
	}
	ackFields, fieldsErr := ackEnvelope.PayloadFields()
	if fieldsErr != nil {
		t.Fatalf("parse replay ack fields: %v", fieldsErr)
	}
	if ackFields["replay_status"] != "not_promised" {
		t.Fatalf("replay ack fields = %#v, want replay_status=not_promised", ackFields)
	}
	if accounting.ledger.Trust("fulfillment") != 1 {
		t.Fatalf("exact replay trust = %d, want first-message-only trust 1", accounting.ledger.Trust("fulfillment"))
	}
	if !hasEventOutcome(accounting.events, "replay_envelope_rejected", "non_commitment") {
		t.Fatalf("exact replay rejection should be local non-commitment event: %#v", accounting.events)
	}
}

func TestRunScopedStatePersistsWithinRun(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	bob := NewNode(cfg, cfg.Agents[1], &decision.FakeDecider{}, decision.FakeMonitor{})
	contentBytes := production.SampleContentBytes()
	contentCID := production.ContentCID(contentBytes)
	if _, err := bob.storeLocalCASObject(contentBytes, agentCASStoreOptions{
		Kind:         agentCASKindPeer,
		ProtocolName: pcid.CASStorageV1,
		Retention:    "test-run-local",
		Paid:         true,
	}); err != nil {
		t.Fatalf("store filesystem CAS object: %v", err)
	}
	bob.capabilityTokens["alice|"+contentCID] = "token-1"
	bob.computeCache["compute-1"] = map[string]string{"result_cid": "result-1"}
	bob.nonCommitmentJournal["nc-1"] = nonCommitmentRecord{Key: "nc-1", Peer: "alice", ProtocolName: pcid.CASStorageV1, PromiseAbout: production.PromiseStoreContent}
	bob.checkpointJournal["cp-1"] = checkpointRecord{Key: "cp-1", ProtocolName: pcid.CASStorageV1, PromiseAbout: production.PromiseStoreContent}
	bob.promiseJournal["pr-1"] = promiseRecord{Key: "pr-1", Peer: "alice", ProtocolName: pcid.CASStorageV1, Status: promiseStatusOutstanding}
	bob.replayJournal["cid-1"] = "alice|" + pcid.CASStorageV1
	if err := bob.saveRunScopedState(); err != nil {
		t.Fatalf("save run-scoped state: %v", err)
	}

	reloadedBob := NewNode(cfg, cfg.Agents[1], &decision.FakeDecider{}, decision.FakeMonitor{})
	if err := reloadedBob.loadRunScopedState(); err != nil {
		t.Fatalf("load run-scoped state: %v", err)
	}
	reloadedBytes, stored, readErr := reloadedBob.readLocalCASObject(contentCID)
	if readErr != nil {
		t.Fatalf("read reloaded filesystem CAS object: %v", readErr)
	}
	if !stored || string(reloadedBytes) != string(contentBytes) {
		t.Fatalf("reloaded CAS object mismatch")
	}
	savedState := bob.exportRunScopedState()
	if len(savedState.CASObjects) != 0 {
		t.Fatalf("new run-scoped state should omit base64 CAS bytes: %#v", savedState.CASObjects)
	}
	if reloadedBob.capabilityTokens["alice|"+contentCID] != "token-1" {
		t.Fatalf("reloaded capability token mismatch: %#v", reloadedBob.capabilityTokens)
	}
	if reloadedBob.computeCache["compute-1"]["result_cid"] != "result-1" {
		t.Fatalf("reloaded compute cache mismatch: %#v", reloadedBob.computeCache)
	}
	if reloadedBob.promiseJournal["pr-1"].Status != promiseStatusOutstanding || reloadedBob.replayJournal["cid-1"] == "" {
		t.Fatalf("reloaded journals mismatch promises=%#v replay=%#v", reloadedBob.promiseJournal, reloadedBob.replayJournal)
	}
	if !hasEvent(reloadedBob.events, "run_scoped_store_loaded") {
		t.Fatalf("load should record run-scoped store event: %#v", reloadedBob.events)
	}
}

func TestAgentCASStorageProfilesWriteExpectedFormats(t *testing.T) {
	cfg := casProfileCoverageTestConfig(t)
	cfg.RunID = "filesystem-profile-test"
	agentsByProfile := map[string]*Node{}
	for _, agentConfig := range cfg.Agents {
		node := NewNode(cfg, agentConfig, &decision.FakeDecider{}, decision.FakeMonitor{})
		agentsByProfile[node.agentCASStorageProfileFor()] = node
	}
	for _, profile := range []string{agentCASStorageProfileGenericBinary, agentCASStorageProfileTypedExtension, agentCASStorageProfileCBORWrapper} {
		if agentsByProfile[profile] == nil {
			t.Fatalf("profile %s was not assigned in test config", profile)
		}
	}
	genericNode := agentsByProfile[agentCASStorageProfileGenericBinary]
	genericCID, genericErr := genericNode.storeLocalCASObject([]byte("plain bytes"), agentCASStoreOptions{Kind: agentCASKindInternal})
	if genericErr != nil {
		t.Fatalf("store generic profile object: %v", genericErr)
	}
	if genericNode.agentCASStore[genericCID].ByteFormat != agentCASByteFormatBinary || !strings.HasSuffix(genericNode.agentCASStore[genericCID].RelativePath, ".bin") {
		t.Fatalf("generic profile metadata = %#v", genericNode.agentCASStore[genericCID])
	}
	typedNode := agentsByProfile[agentCASStorageProfileTypedExtension]
	cborBytes, cborErr := protocol.MarshalStringMap(map[string]string{"kind": "profile-test"})
	if cborErr != nil {
		t.Fatalf("marshal test cbor: %v", cborErr)
	}
	typedCID, typedErr := typedNode.storeLocalCASObject(cborBytes, agentCASStoreOptions{Kind: agentCASKindMessage})
	if typedErr != nil {
		t.Fatalf("store typed profile object: %v", typedErr)
	}
	if typedNode.agentCASStore[typedCID].ByteFormat != agentCASByteFormatCBOR || !strings.HasSuffix(typedNode.agentCASStore[typedCID].RelativePath, ".cbor") {
		t.Fatalf("typed profile metadata = %#v", typedNode.agentCASStore[typedCID])
	}
	wrapperNode := agentsByProfile[agentCASStorageProfileCBORWrapper]
	wrapperCID, wrapperErr := wrapperNode.storeLocalCASObject([]byte("wrapped exact bytes"), agentCASStoreOptions{Kind: agentCASKindPeer})
	if wrapperErr != nil {
		t.Fatalf("store wrapper profile object: %v", wrapperErr)
	}
	wrapperRecord := wrapperNode.agentCASStore[wrapperCID]
	if wrapperRecord.ByteFormat != agentCASByteFormatCBORWrapper || wrapperRecord.StoredCID == wrapperCID || !strings.HasSuffix(wrapperRecord.RelativePath, ".cbor") {
		t.Fatalf("wrapper profile metadata = %#v", wrapperRecord)
	}
	originalBytes, stored, readErr := wrapperNode.readLocalCASObject(wrapperCID)
	if readErr != nil {
		t.Fatalf("read wrapper profile object: %v", readErr)
	}
	if !stored || string(originalBytes) != "wrapped exact bytes" {
		t.Fatalf("wrapper read = %q stored=%v", string(originalBytes), stored)
	}
}

func TestAgentCASWrapperModesCoverCurrentAgents(t *testing.T) {
	cfg := casProfileCoverageTestConfig(t)
	cfg.RunID = "wrapper-mode-coverage-test"
	seenModes := map[string]bool{}
	for _, agentConfig := range cfg.Agents {
		node := NewNode(cfg, agentConfig, &decision.FakeDecider{}, decision.FakeMonitor{})
		if node.agentCASStorageProfileFor() == agentCASStorageProfileCBORWrapper {
			seenModes[node.agentCASWrapperModeFor()] = true
		}
	}
	for _, wrapperMode := range []string{agentCASWrapperModeOriginalKey, agentCASWrapperModeWrapperKey, agentCASWrapperModeDualKey} {
		if !seenModes[wrapperMode] {
			t.Fatalf("wrapper mode %s was not assigned across current agents: %#v", wrapperMode, seenModes)
		}
	}
}

func TestLegacyCASObjectsMigrateToFilesystemCAS(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	bob := NewNode(cfg, cfg.Agents[1], &decision.FakeDecider{}, decision.FakeMonitor{})
	contentBytes := production.SampleContentBytes()
	contentCID := production.ContentCID(contentBytes)
	legacyState := runScopedState{
		Version: 1,
		CASObjects: map[string]string{
			contentCID: base64.StdEncoding.EncodeToString(contentBytes),
		},
		AgentCASObjects: map[string]agentCASObject{
			contentCID: {
				CID:          contentCID,
				Kind:         agentCASKindPeer,
				Owner:        "bob",
				ProtocolName: pcid.CASStorageV1,
				Retention:    "legacy-test",
				Paid:         true,
			},
		},
	}
	stateBytes, marshalErr := json.MarshalIndent(legacyState, "", "  ")
	if marshalErr != nil {
		t.Fatalf("marshal legacy state: %v", marshalErr)
	}
	statePath := bob.runScopedStatePath()
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		t.Fatalf("create state dir: %v", err)
	}
	if err := os.WriteFile(statePath, append(stateBytes, '\n'), 0o644); err != nil {
		t.Fatalf("write legacy state: %v", err)
	}
	if err := bob.loadRunScopedState(); err != nil {
		t.Fatalf("load legacy state: %v", err)
	}
	migratedBytes, stored, readErr := bob.readLocalCASObject(contentCID)
	if readErr != nil {
		t.Fatalf("read migrated CAS object: %v", readErr)
	}
	if !stored || string(migratedBytes) != string(contentBytes) {
		t.Fatalf("migrated CAS mismatch stored=%v bytes=%q", stored, string(migratedBytes))
	}
	if err := bob.saveRunScopedState(); err != nil {
		t.Fatalf("save migrated state: %v", err)
	}
	savedBytes, readStateErr := os.ReadFile(statePath)
	if readStateErr != nil {
		t.Fatalf("read migrated state: %v", readStateErr)
	}
	if strings.Contains(string(savedBytes), "cas_objects_b64") {
		t.Fatalf("new state should omit legacy base64 CAS bytes: %s", string(savedBytes))
	}
}

func TestCapabilityTokenReplayReturnsNonCommitmentEvent(t *testing.T) {
	cfg := twoNodeTestConfig(t)
	bob := NewNode(cfg, cfg.Agents[1], &decision.FakeDecider{}, decision.FakeMonitor{})
	contentBytes := production.SampleContentBytes()
	contentCID := production.ContentCID(contentBytes)
	storeAck, storeErr := bob.handleCASStoragePromise(map[string]string{
		"from":          "alice",
		"promise_about": production.PromiseStoreContent,
		"content_cid":   contentCID,
		"content_b64":   base64.StdEncoding.EncodeToString(contentBytes),
		"credit_offer":  "4",
		"units":         "1",
	})
	if storeErr != nil {
		t.Fatalf("store content: %v", storeErr)
	}
	token := storeAck["capability_token"]
	if _, serveErr := bob.handleCASStoragePromise(map[string]string{
		"from":          "alice",
		"promise_about": production.PromiseServeContent,
		"content_cid":   contentCID,
		"token":         token,
	}); serveErr != nil {
		t.Fatalf("serve content: %v", serveErr)
	}
	replayAck, replayErr := bob.handleCASStoragePromise(map[string]string{
		"from":          "alice",
		"promise_about": production.PromiseServeContent,
		"content_cid":   contentCID,
		"token":         token,
	})
	if replayErr != nil {
		t.Fatalf("token replay should be non-commitment event, not handler error: %v", replayErr)
	}
	if replayAck["token_status"] != "not_promised" {
		t.Fatalf("replay ack = %#v, want token_status=not_promised", replayAck)
	}
	if !hasEventOutcome(bob.events, "capability_token_replay_rejected", "non_commitment") {
		t.Fatalf("token replay should be recorded as local non-commitment: %#v", bob.events)
	}
	if !hasEvent(bob.events, "gc_object_removed") || !hasEvent(bob.events, "gc_promise_ended") {
		t.Fatalf("token redemption should record GC event: %#v", bob.events)
	}
}

func TestRetentionPromiseBrokenRecordsLocalEvent(t *testing.T) {
	cfg := singleNodeTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	alice.recordRetentionPromiseBroken("alice", "test retention miss after simulated crash")
	if !hasEventOutcome(alice.events, "retention_promise_broken", "broken") {
		t.Fatalf("retention break case should be explicit local events: %#v", alice.events)
	}
}

func TestBadResultProbeReducesComputePromiserTrust(t *testing.T) {
	// Intent: Alice's own recomputation event record should reduce trust in the
	// compute promiser and leave a still-trusted alternate compute peer usable for
	// follow-up work. Source: DI-vahan
	cfg := computeRoutingTestConfig(t)
	alice := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, decision.FakeMonitor{})
	message := parsedMessage{Fields: computeAckFields(t)}
	if err := alice.verifyComputeAckLocally(message, "carol"); err != nil {
		t.Fatalf("verify compute ack locally: %v", err)
	}
	if alice.ledger.Trust("carol") != -3 {
		t.Fatalf("carol trust = %d, want -3 after malformed bad-result event", alice.ledger.Trust("carol"))
	}
	if alice.canDial("carol") {
		t.Fatalf("Alice should stop promising direct sends to Carol after malformed compute event")
	}
	if !alice.canDial("dave") {
		t.Fatalf("Alice should still be able to use Dave as the alternate compute peer")
	}
	if !hasEventOutcome(alice.events, "compute_result_locally_rejected", "malformed") {
		t.Fatalf("local bad-result rejection should be recorded: %#v", alice.events)
	}
}

func TestMalformedProofEventReducesIdentifiedPromiserTrust(t *testing.T) {
	// Intent: A parseable envelope with a stale proof should be attributed to the
	// claimed promiser as a local malformed event, not counted as random
	// transport noise. Source: DI-sunuf
	cfg := computeRoutingTestConfig(t)
	grace := NewNode(cfg, cfg.Agents[2], &decision.FakeDecider{}, decision.FakeMonitor{})
	fields := map[string]string{
		"act":           decision.ActPromise,
		"from":          "mallory",
		"to":            "grace",
		"turn":          "test",
		"promise":       "Mallory promises this parseable but mutated proof is valid.",
		"reason":        "test bad proof attribution",
		"promise_about": production.PromisePresentStorageReport,
	}
	payloadBytes, _, payloadErr := protocol.MarshalKnownArrayPayload(pcid.CASStorageV1, fields)
	if payloadErr != nil {
		t.Fatalf("marshal payload: %v", payloadErr)
	}
	envelope, envelopeErr := protocol.NewEnvelopeFromPayload(grace.Protocols.MustCID(pcid.CASStorageV1), payloadBytes, "mallory")
	if envelopeErr != nil {
		t.Fatalf("new envelope: %v", envelopeErr)
	}
	mutatedFields := copyStringMap(fields)
	mutatedFields["reason"] = "payload changed after signing"
	mutatedPayload, _, mutatedPayloadErr := protocol.MarshalKnownArrayPayload(pcid.CASStorageV1, mutatedFields)
	if mutatedPayloadErr != nil {
		t.Fatalf("marshal mutated payload: %v", mutatedPayloadErr)
	}
	envelope.Payload = mutatedPayload
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		t.Fatalf("envelope bytes: %v", bytesErr)
	}
	if _, err := grace.handleFrame(envelopeBytes); err == nil {
		t.Fatalf("mutated payload under stale proof should be rejected")
	}
	if grace.ledger.Trust("mallory") != -3 {
		t.Fatalf("mallory trust = %d, want -3 after bad proof", grace.ledger.Trust("mallory"))
	}
	if !hasEventOutcome(grace.events, "malformed_proof_observed", "malformed") {
		t.Fatalf("bad proof should be attributed as malformed event: %#v", grace.events)
	}
}

func TestCorruptCASEventReducesIdentifiedPromiserTrust(t *testing.T) {
	// Intent: A corrupt content-addressed storage event records a malformed
	// promise by the presenting peer, so it must reduce local trust just like a
	// bad proof or bad compute result. Source: DI-fijov
	cfg := computeRoutingTestConfig(t)
	grace := NewNode(cfg, cfg.Agents[2], &decision.FakeDecider{}, decision.FakeMonitor{})
	goodBytes := []byte("expected content")
	badBytes := []byte("different content")
	ackFields, err := grace.handleCASStoragePromise(map[string]string{
		"from":          "mallory",
		"promise_about": production.PromisePresentStorageReport,
		"content_cid":   production.ContentCID(goodBytes),
		"content_b64":   base64.StdEncoding.EncodeToString(badBytes),
	})
	if err != nil {
		t.Fatalf("corrupt CAS event should return a broken verdict ACK, not handler error: %v", err)
	}
	if ackFields["verdict"] != "broken" {
		t.Fatalf("corrupt CAS verdict = %#v, want broken", ackFields)
	}
	if grace.ledger.Trust("mallory") != -3 {
		t.Fatalf("mallory trust = %d, want -3 after corrupt CAS event", grace.ledger.Trust("mallory"))
	}
	if grace.canDial("mallory") {
		t.Fatalf("Grace should stop promising direct sends to Mallory after corrupt CAS event")
	}
	if !hasEventOutcome(grace.events, "cas_corrupt_bytes_rejected", "malformed") {
		t.Fatalf("corrupt CAS event should be recorded as malformed: %#v", grace.events)
	}
}

func TestFutureRepairPromiseDoesNotImmediatelyIncreaseTrust(t *testing.T) {
	// Intent: A promise to repair future behavior is a useful local event to
	// remember, but it is not proof that the future repair has already been kept.
	// Source: DI-fijov
	ackFields := map[string]string{
		"promise_about": production.PromiseLabelFutureMalformedReport,
		"repair_status": "future_only",
	}
	if eventUpdatesTrust(ackFields) {
		t.Fatalf("future-only repair promise should not immediately update trust: %#v", ackFields)
	}
}

func TestTrustCautionAllowsOnlyFutureRepairCandidateTraffic(t *testing.T) {
	// Intent: After malformed events, a peer may be heard only for the narrow
	// future-repair candidate promise; unsupported ordinary variants remain
	// non-promised until local trust is rebuilt. Source: DI-sihuz
	cfg := computeRoutingTestConfig(t)
	grace := NewNode(cfg, cfg.Agents[2], &decision.FakeDecider{}, decision.FakeMonitor{})
	grace.observeOutcome("mallory", relationship.OutcomeMalformed)
	unsupportedFields := map[string]string{"promise_about": production.PromiseUnsupportedVariantProbe}
	if grace.canAcceptFrom("mallory", unsupportedFields) {
		t.Fatalf("unsupported ordinary promise should not be accepted during trust caution")
	}
	repairFields := map[string]string{"promise_about": production.PromiseLabelFutureMalformedReport}
	if !grace.canAcceptFrom("mallory", repairFields) {
		t.Fatalf("future repair promise should remain hearable as candidate traffic")
	}
	grace.observeOutcome("mallory", relationship.OutcomeKept)
	if !hasEvent(grace.events, "trust_caution_recorded") {
		t.Fatalf("malformed event should record trust caution: %#v", grace.events)
	}
	if !hasEvent(grace.events, "trust_recovery_delayed") {
		t.Fatalf("kept event during caution should record delayed recovery: %#v", grace.events)
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
	if !printerPort.supportsProtocol(pcid.ProductionShippingV1) {
		t.Fatalf("printer_port should support production shipping pCID")
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

func TestNodeRunDoesNotUseMonitorMarkerFiles(t *testing.T) {
	cfg := singleNodeTestConfig(t)
	node := NewNode(cfg, cfg.Agents[0], &decision.FakeDecider{}, failingMonitor{})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := node.Run(ctx); err != nil {
		t.Fatalf("run node with unused monitor: %v", err)
	}
	for _, markerName := range []string{"monitor.done", "alice.done", "alice.turns_done"} {
		if _, err := os.Stat(filepath.Join(cfg.RunRoot, cfg.RunID, markerName)); err == nil {
			t.Fatalf("marker file %s should not exist after natural-exit run", markerName)
		} else if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("stat marker file %s: %v", markerName, err)
		}
	}
	if hasEventOutcome(node.events, "monitor_error", "non_commitment") || hasEventOutcome(node.events, "monitor_done", "non_commitment") {
		t.Fatalf("node runtime should not run the observer monitor: %#v", node.events)
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
		EventCollectorAddress: "event-collector:9200",
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
		EventCollectorAddress: "event-collector:9200",
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

func casProfileCoverageTestConfig(t *testing.T) config.Config {
	t.Helper()
	cfg := twoNodeTestConfig(t)
	cfg.RunRoot = filepath.Join(t.TempDir(), "run")
	agentNames := []string{"alice", "bob", "carol", "dave", "ellen", "frank", "grace", "mallory", "victor"}
	cfg.Agents = make([]config.AgentConfig, 0, len(agentNames))
	cfg.Containers = make([]config.ContainerConfig, 0, len(agentNames))
	for _, agentName := range agentNames {
		cfg.Agents = append(cfg.Agents, config.AgentConfig{
			Name:           agentName,
			Persona:        "cas profile tester",
			Motivation:     "test local CAS storage profile assignment",
			InitialPeers:   []string{"alice"},
			CandidatePeers: []string{"alice"},
			Budget:         5,
			Capacity:       5,
		})
		cfg.Containers = append(cfg.Containers, config.ContainerConfig{Name: agentName, Agents: []string{agentName}})
	}
	return cfg
}

func computeRoutingTestConfig(t *testing.T) config.Config {
	t.Helper()
	// Intent: Keep compute-route tests deterministic without opening sockets:
	// Alice begins with Carol and Dave as direct peers, while Grace/Mallory let
	// malformed-proof attribution run through the same local ledger code.
	// Source: DI-vahan; DI-sunuf
	return config.Config{
		RunID:                 "compute-routing-test",
		RunRoot:               filepath.Join(t.TempDir(), "run"),
		ListenPort:            0,
		EventCollectorAddress: "event-collector:9200",
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
		StrongTrustThreshold:  2,
		WeakTrustThreshold:    -2,
		TrustDecayPerRound:    0,
		Agents: []config.AgentConfig{{
			Name:           "alice",
			Persona:        "tester",
			Motivation:     "test",
			InitialPeers:   []string{"carol", "dave"},
			CandidatePeers: []string{"carol", "dave"},
			SupportedPCIDs: []string{pcid.RelationshipV1, pcid.CIDComputeV1, pcid.CASStorageV1},
			Budget:         5,
			Capacity:       5,
		}, {
			Name:           "carol",
			Persona:        "compute",
			Motivation:     "test",
			InitialPeers:   []string{"alice"},
			CandidatePeers: []string{"alice"},
			SupportedPCIDs: []string{pcid.RelationshipV1, pcid.CIDComputeV1},
			Budget:         5,
			Capacity:       5,
		}, {
			Name:           "grace",
			Persona:        "receiver",
			Motivation:     "test",
			InitialPeers:   []string{"mallory"},
			CandidatePeers: []string{"mallory"},
			SupportedPCIDs: []string{pcid.RelationshipV1, pcid.CASStorageV1},
			Budget:         5,
			Capacity:       5,
		}, {
			Name:           "dave",
			Persona:        "compute verifier",
			Motivation:     "test",
			InitialPeers:   []string{"alice"},
			CandidatePeers: []string{"alice"},
			SupportedPCIDs: []string{pcid.RelationshipV1, pcid.CIDComputeV1},
			Budget:         5,
			Capacity:       5,
		}, {
			Name:           "mallory",
			Persona:        "adversary",
			Motivation:     "test",
			InitialPeers:   []string{"grace"},
			CandidatePeers: []string{"grace"},
			SupportedPCIDs: []string{pcid.RelationshipV1, pcid.CASStorageV1},
			Budget:         5,
			Capacity:       5,
		}},
		Containers: []config.ContainerConfig{
			{Name: "alice", Agents: []string{"alice"}},
			{Name: "carol", Agents: []string{"carol"}},
			{Name: "grace", Agents: []string{"grace"}},
			{Name: "dave", Agents: []string{"dave"}},
			{Name: "mallory", Agents: []string{"mallory"}},
		},
	}
}

func shippingTestConfig(t *testing.T) config.Config {
	t.Helper()
	return config.Config{
		RunID:                 "shipping-test",
		RunRoot:               filepath.Join(t.TempDir(), "run"),
		ListenPort:            0,
		EventCollectorAddress: "event-collector:9200",
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
			SupportedPCIDs: []string{pcid.RelationshipV1, pcid.ProductionShippingV1},
			Budget:         5,
			Capacity:       5,
		}, {
			Name:           "ups_label_printer",
			Kind:           "ups_label_printer",
			Persona:        "printer",
			Motivation:     "print label",
			InitialPeers:   []string{"fulfillment", "printer_port"},
			CandidatePeers: []string{"fulfillment", "printer_port"},
			SupportedPCIDs: []string{pcid.RelationshipV1, pcid.ProductionShippingV1},
			Budget:         5,
			Capacity:       5,
		}, {
			Name:           "printer_port",
			Kind:           "printer_port",
			Persona:        "printer port",
			Motivation:     "local hardware access",
			InitialPeers:   []string{"ups_label_printer"},
			CandidatePeers: []string{"ups_label_printer"},
			SupportedPCIDs: []string{pcid.RelationshipV1, pcid.ProductionShippingV1},
			Budget:         5,
			Capacity:       5,
		}, {
			Name:           "accounting",
			Kind:           "accounting",
			Persona:        "accounting",
			Motivation:     "records",
			InitialPeers:   []string{"fulfillment"},
			CandidatePeers: []string{"fulfillment"},
			SupportedPCIDs: []string{pcid.RelationshipV1, pcid.ProductionShippingV1},
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

func signedAccountingUpdateFrame(t *testing.T, node *Node, exchangeID string) []byte {
	t.Helper()
	fields := map[string]string{
		"act":             decision.ActPromise,
		"from":            "fulfillment",
		"to":              "accounting",
		"turn":            "test",
		"promise":         "I promise to receive accounting's shipment checkpoint event for this order and tracking number.",
		"reason":          "duplicate checkpoint regression test",
		"exchange_id":     exchangeID,
		"promise_about":   production.PromiseShipmentUpdate,
		"order_id":        fulfillmentOrderID,
		"tracking_number": production.ContentCID([]byte("test-accounting-update-tracking")),
		"cost_cents":      "1776",
	}
	payloadBytes, _, payloadErr := protocol.MarshalKnownArrayPayload(pcid.ProductionShippingV1, fields)
	if payloadErr != nil {
		t.Fatalf("marshal accounting update payload: %v", payloadErr)
	}
	envelope, err := protocol.NewEnvelopeFromPayload(node.Protocols.MustCID(pcid.ProductionShippingV1), payloadBytes, "fulfillment")
	if err != nil {
		t.Fatalf("new accounting update envelope: %v", err)
	}
	frameBytes, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("accounting update envelope bytes: %v", err)
	}
	return frameBytes
}

func computeAckFields(t *testing.T) map[string]string {
	t.Helper()
	functionBytes := production.SampleFunctionBytes()
	inputBytes := production.SampleInputBytes()
	contextBytes := production.SampleContextBytes()
	resultBytes, executeErr := production.ExecuteFunction(functionBytes, inputBytes, contextBytes)
	if executeErr != nil {
		t.Fatalf("execute sample function: %v", executeErr)
	}
	badResultBytes := production.BadComputeResultBytes(resultBytes)
	return map[string]string{
		"promise_about": production.PromiseExecuteFunction,
		"function_cid":  production.ContentCID(functionBytes),
		"function_b64":  base64.StdEncoding.EncodeToString(functionBytes),
		"input_cid":     production.ContentCID(inputBytes),
		"input_b64":     base64.StdEncoding.EncodeToString(inputBytes),
		"context_cid":   production.ContentCID(contextBytes),
		"context_b64":   base64.StdEncoding.EncodeToString(contextBytes),
		"result_cid":    production.ContentCID(resultBytes),
		"result_b64":    base64.StdEncoding.EncodeToString(resultBytes),
		"bad_result_cid": production.ContentCID(
			badResultBytes,
		),
		"bad_result_b64": base64.StdEncoding.EncodeToString(badResultBytes),
	}
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

func hasEventDetail(events []decision.Event, eventName, detailPart string) bool {
	for _, event := range events {
		if event.Event == eventName && strings.Contains(event.Detail, detailPart) {
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
