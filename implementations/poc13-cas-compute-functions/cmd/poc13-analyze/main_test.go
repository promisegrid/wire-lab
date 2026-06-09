package main

import (
	"os"
	"path/filepath"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/decision"
	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/pcid"
)

func TestAnalyzeRunSummarizesEventsAndMonitorReport(t *testing.T) {
	runDir := t.TempDir()
	writeFile(t, filepath.Join(runDir, "alice.jsonl"), ""+
		`{"observer":"alice","event":"promise_sent","outcome":"kept","peer":"bob","detail":"sent"}`+"\n"+
		`{"observer":"alice","event":"send_failed","outcome":"broken","peer":"bob","detail":"connection refused"}`+"\n"+
		`{"observer":"alice","event":"shipping_label_received","outcome":"kept","peer":"ups_label_printer","detail":"tracking"}`+"\n"+
		`{"observer":"alice","event":"printer_port_print_confirmed","outcome":"kept","peer":"printer_port","detail":"spool"}`+"\n"+
		`{"observer":"alice","event":"local_resource_exhausted","outcome":"non_commitment","peer":"bob","detail":"capacity exhausted"}`+"\n"+
		`{"observer":"alice","event":"direct_peer_unchanged","outcome":"kept","peer":"bob","detail":"outcome=non_commitment trust=0"}`+"\n")
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
	if summary.TotalEvents != 8 {
		t.Fatalf("total events = %d, want 8", summary.TotalEvents)
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

func TestValidateSummaryAcceptsCleanRegressionEvidence(t *testing.T) {
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
			pcid.CASStorageV1:     1,
			pcid.CIDComputeV1:     1,
			pcid.EvidenceReportV1: 1,
		},
		ShippingCounts: map[string]int{
			"accounting_updated":                    1,
			"accounting_update_duplicate":           1,
			"accounting_update_duplicate_confirmed": 1,
		},
		RelationshipTransitionCounts: map[string]int{},
		LocalResourceCounts:          map[string]int{},
		ResourceTrustCouplingCounts:  map[string]int{},
		MonitorReport: &decision.MonitorReport{
			PromiseTheoryFit:      5,
			Autonomy:              5,
			ProtocolValidity:      4,
			LocalTrustCorrectness: 5,
			ImpositionAvoidance:   5,
		},
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
