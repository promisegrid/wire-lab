package main

import (
	"os"
	"path/filepath"
	"testing"
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

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
