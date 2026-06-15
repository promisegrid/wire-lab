package main

import (
	"os"
	"path/filepath"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/decision"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/pcid"
)

func TestAnalyzeRunSummarizesEventsAndMonitorReport(t *testing.T) {
	runDir := t.TempDir()
	writeFile(t, filepath.Join(runDir, "alice.jsonl"), ""+
		`{"observer":"alice","event":"promise_sent","outcome":"kept","peer":"bob","detail":"sent"}`+"\n"+
		`{"observer":"alice","event":"send_failed","outcome":"broken","peer":"bob","detail":"connection refused"}`+"\n"+
		`{"observer":"alice","event":"shipping_label_received","outcome":"kept","peer":"ups_label_printer","detail":"tracking"}`+"\n"+
		`{"observer":"alice","event":"printer_port_print_confirmed","outcome":"kept","peer":"printer_port","detail":"spool"}`+"\n"+
		`{"observer":"alice","event":"local_resource_exhausted","outcome":"non_commitment","peer":"bob","detail":"capacity exhausted"}`+"\n"+
		`{"observer":"alice","event":"direct_peer_unchanged","outcome":"kept","peer":"bob","detail":"outcome=non_commitment trust=0"}`+"\n"+
		`{"observer":"alice","event":"wasm_module_instantiated","outcome":"kept","peer":"victor","detail":"runtime=wazero"}`+"\n"+
		`{"observer":"alice","event":"wasm_adapter_ack_received","outcome":"kept","peer":"victor","detail":"pcid=relationship_v1"}`+"\n"+
		`{"observer":"alice","event":"stdio_cbor_ack_event","outcome":"kept","peer":"victor","detail":"exact_sha256=test"}`+"\n"+
		`{"observer":"alice","event":"bearer_token_exchange_rate_observed","outcome":"kept","peer":"grace","detail":"local market signal"}`+"\n"+
		`{"observer":"alice","event":"mixed_version_successor_pcid_selected","outcome":"kept","peer":"bob","detail":"current pCID"}`+"\n"+
		`{"observer":"alice","event":"run_internal_restart_recovery_observed","outcome":"kept","peer":"victor","detail":"same-run recovery"}`+"\n")
	writeFile(t, filepath.Join(runDir, "bob.jsonl"), ""+
		`{"observer":"bob","event":"decision_rejected","outcome":"malformed","peer":"alice","detail":"bad target"}`+"\n"+
		`{"observer":"bob","event":"direct_peer_added","outcome":"kept","peer":"alice","detail":"trust=2"}`+"\n")
	writeFile(t, filepath.Join(runDir, "monitor-report.json"), `{
  "promise_theory_fit": 4,
  "autonomy": 5,
  "protocol_validity": 3,
  "local_trust_correctness": 4,
  "imposition_avoidance": 5,
  "summary": "test",
  "concerns": ["synthetic"]
}`)

	summary, err := analyzeRun(runDir)
	if err != nil {
		t.Fatalf("analyze run: %v", err)
	}
	if summary.TotalEvents != 14 {
		t.Fatalf("total events = %d, want 14", summary.TotalEvents)
	}
	if summary.EventCounts["promise_sent"] != 1 || summary.FailureCounts["send_failed"] != 1 || summary.FailureCounts["decision_rejected"] != 1 {
		t.Fatalf("unexpected counts: %#v failures %#v", summary.EventCounts, summary.FailureCounts)
	}
	if summary.MonitorReport == nil || summary.MonitorReport.ProtocolValidity != 3 {
		t.Fatalf("monitor report not summarized: %#v", summary.MonitorReport)
	}
	if summary.ShippingCounts["shipping_label_received"] != 1 {
		t.Fatalf("shipping counts not summarized: %#v", summary.ShippingCounts)
	}
	if summary.ShippingCounts["printer_port_print_confirmed"] != 1 {
		t.Fatalf("printer-port counts not summarized: %#v", summary.ShippingCounts)
	}
	if summary.RelationshipTransitionCounts["direct_peer_added"] != 1 {
		t.Fatalf("relationship transitions not summarized: %#v", summary.RelationshipTransitionCounts)
	}
	if summary.LocalResourceCounts["local_resource_exhausted"] != 1 {
		t.Fatalf("local resource counts not summarized: %#v", summary.LocalResourceCounts)
	}
	if summary.ResourceTrustCouplingCounts["alice"] != 1 {
		t.Fatalf("resource/trust coupling not summarized: %#v", summary.ResourceTrustCouplingCounts)
	}
	if summary.RuntimeAdapterEventCounts["wasm_adapter_ack_received"] != 1 {
		t.Fatalf("runtime adapter event counts not summarized: %#v", summary.RuntimeAdapterEventCounts)
	}
	if summary.DecentralizedMonitorCounts["bearer_token_exchange_rate_observed"] != 1 {
		t.Fatalf("decentralized monitor counts not summarized: %#v", summary.DecentralizedMonitorCounts)
	}
	if summary.MigrationCounts["mixed_version_successor_pcid_selected"] != 1 {
		t.Fatalf("migration counts not summarized: %#v", summary.MigrationCounts)
	}
	if summary.RestartCounts["run_internal_restart_recovery_observed"] != 1 {
		t.Fatalf("restart counts not summarized: %#v", summary.RestartCounts)
	}
}

func TestAnalyzeRunAcceptsParentRunDirectory(t *testing.T) {
	parentDir := t.TempDir()
	logDir := filepath.Join(parentDir, "run")
	if err := os.Mkdir(logDir, 0o755); err != nil {
		t.Fatalf("make log dir: %v", err)
	}
	writeFile(t, filepath.Join(logDir, "alice.jsonl"), `{"observer":"alice","event":"promise_sent","outcome":"kept","peer":"bob","detail":"sent"}`+"\n")

	summary, err := analyzeRun(parentDir)
	if err != nil {
		t.Fatalf("analyze parent run dir: %v", err)
	}
	if summary.RunDir != logDir {
		t.Fatalf("summary run dir = %q, want %q", summary.RunDir, logDir)
	}
	if summary.TotalEvents != 1 {
		t.Fatalf("total events = %d, want 1", summary.TotalEvents)
	}
}

func TestAnalyzeRunRejectsDirectoryWithoutLogs(t *testing.T) {
	_, err := analyzeRun(t.TempDir())
	if err == nil {
		t.Fatalf("analyze without logs should fail")
	}
}

func TestValidateSummaryAcceptsCleanRegressionEvent(t *testing.T) {
	summary := cleanRegressionSummary()
	if err := validateSummary(summary, cleanRegressionCriteria()); err != nil {
		t.Fatalf("clean regression summary should pass: %v", err)
	}
}

func TestValidateSummaryRejectsMissingSuppression(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["promise_not_promised_suppressed"] = 0
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing suppression should fail")
	}
}

func TestValidateSummaryRejectsResourceTrustCoupling(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.ResourceTrustCouplingCounts["alice"] = 1
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("resource/trust coupling should fail")
	}
}

func TestValidateSummaryRejectsMissingAdapterEvent(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["stdio_cbor_ack_event"] = 0
	summary.ScoreReport = computeScores(summary)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing stdio adapter event should fail")
	}
}

func TestValidateSummaryRejectsForbiddenVocabulary(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.ForbiddenVocabularyCounts["run_events"] = 1
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("forbidden vocabulary should fail")
	}
}

func TestAnalyzeRunCountsForbiddenVocabulary(t *testing.T) {
	retiredWord := "evi" + "dence"
	retiredInterfaceWord := "boun" + "dary"
	runDir := t.TempDir()
	writeFile(t, filepath.Join(runDir, "alice.jsonl"), `{"observer":"alice","event":"promise_sent","outcome":"kept","peer":"bob","detail":"fresh `+retiredWord+` and `+retiredInterfaceWord+`"}`+"\n")
	writeFile(t, filepath.Join(runDir, "monitor-report.json"), `{"promise_theory_fit":5,"autonomy":5,"protocol_validity":5,"local_trust_correctness":5,"imposition_avoidance":5,"summary":"clean","concerns":["no issues"]}`)
	summary, err := analyzeRun(runDir)
	if err != nil {
		t.Fatalf("analyze run: %v", err)
	}
	if summary.ForbiddenVocabularyCounts["run_events"] < 2 {
		t.Fatalf("forbidden vocabulary not counted: %#v", summary.ForbiddenVocabularyCounts)
	}
}

func TestValidateSummaryRejectsMissingDecentralizedMonitoring(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["bearer_token_exchange_rate_observed"] = 0
	summary.ScoreReport = computeScores(summary)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing decentralized monitoring event should fail")
	}
}

func TestValidateSummaryRejectsMissingMigrationEvent(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["mixed_version_successor_pcid_selected"] = 0
	summary.ScoreReport = computeScores(summary)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing migration event should fail")
	}
}

func TestValidateSummaryRejectsMissingRestartEvent(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["run_internal_restart_recovery_observed"] = 0
	summary.ScoreReport = computeScores(summary)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing restart event should fail")
	}
}

func TestValidateSummaryRejectsMissingPermanentDistrustEvent(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["permanent_distrust_decided"] = 0
	summary.MissingRequiredEventNames = missingRequiredEvents(summary)
	summary.ScoreReport = computeScores(summary)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing permanent distrust event should fail")
	}
}

func TestValidateSummaryRejectsMissingTransitExclusionEvent(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["transit_safe_route_selected"] = 0
	summary.MissingRequiredEventNames = missingRequiredEvents(summary)
	summary.ScoreReport = computeScores(summary)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing transit exclusion event should fail")
	}
}

func TestValidateSummaryRejectsMissingMigratedArrayPayloadEvent(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["pcid_owned_array_payload_received"] = 0
	summary.MissingRequiredEventNames = missingRequiredEvents(summary)
	summary.ScoreReport = computeScores(summary)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing migrated array payload event should fail")
	}
}

func TestValidateSummaryRejectsMissingArrayPayloadProtocol(t *testing.T) {
	summary := cleanRegressionSummary()
	delete(summary.ArrayPayloadProtocolCounts, pcid.RelationshipV1)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing relationship array payload coverage should fail")
	}
}

func cleanRegressionSummary() RunSummary {
	eventCounts := map[string]int{
		"fulfillment_workflow_completed":  1,
		"promise_not_promised_suppressed": 1,
	}
	for _, eventName := range requiredRegressionEvents() {
		eventCounts[eventName] = 1
	}
	summary := RunSummary{
		RunDir:        "test/run",
		TotalEvents:   len(eventCounts),
		EventCounts:   eventCounts,
		OutcomeCounts: map[string]int{},
		AgentCounts:   map[string]int{},
		FailureCounts: map[string]int{},
		ProtocolCounts: map[string]int{
			pcid.KernelReceiveV1: 1,
			pcid.CASStorageV1:    1,
			pcid.CIDComputeV1:    1,
			pcid.IdentityKeyV1:   1,
			pcid.RelationshipV1:  1,
			pcid.AccountingV1:    1,
			pcid.UPSLabelV1:      1,
			pcid.PostalScaleV1:   1,
			pcid.PrinterPortV1:   1,
		},
		ArrayPayloadProtocolCounts: map[string]int{},
		ShippingCounts: map[string]int{
			"accounting_updated":                    1,
			"accounting_update_duplicate":           1,
			"accounting_update_duplicate_confirmed": 1,
		},
		RelationshipTransitionCounts: map[string]int{},
		LocalResourceCounts:          map[string]int{},
		ResourceTrustCouplingCounts:  map[string]int{},
		ForbiddenVocabularyCounts:    map[string]int{},
		MonitorReport: &decision.MonitorReport{
			PromiseTheoryFit:      5,
			Autonomy:              5,
			ProtocolValidity:      4,
			LocalTrustCorrectness: 5,
			ImpositionAvoidance:   5,
		},
	}
	for _, protocolName := range requiredArrayPayloadProtocols() {
		summary.ArrayPayloadProtocolCounts[protocolName] = 1
	}
	summary.ScoreReport = computeScores(summary)
	return summary
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
