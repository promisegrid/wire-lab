package poc13

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeRunAcceptsParentDirectoryAndValidates(t *testing.T) {
	parentDir := t.TempDir()
	logDir := filepath.Join(parentDir, "run")
	if err := os.Mkdir(logDir, 0o755); err != nil {
		t.Fatalf("make log dir: %v", err)
	}
	writePOC13Log(t, filepath.Join(logDir, "alice.jsonl"), []Event{
		{Observer: "alice", Event: "tcp_message_sent", Outcome: "kept", PCID: CASStorageV1, Detail: "sent"},
		{Observer: "bob", Event: "tcp_message_received", Outcome: "kept", PCID: CASStorageV1, Detail: "received"},
		{Observer: "bob", Event: "cas_storage_promised", Outcome: "kept", PCID: CASStorageV1, Detail: "storage"},
		{Observer: "bob", Event: "cas_retention_promised", Outcome: "kept", PCID: CASStorageV1, Detail: "retention"},
		{Observer: "frank", Event: "cas_replication_promised", Outcome: "kept", PCID: CASStorageV1, Detail: "replica"},
		{Observer: "bob", Event: "cas_bytes_stored", Outcome: "kept", PCID: CASStorageV1, Detail: "stored"},
		{Observer: "alice", Event: "cas_bytes_retrieved", Outcome: "kept", PCID: CASStorageV1, Detail: "retrieved"},
		{Observer: "frank", Event: "cas_replica_stored", Outcome: "kept", PCID: CASStorageV1, Detail: "replica stored"},
		{Observer: "mallory", Event: "cas_corrupt_bytes_rejected", Outcome: "malformed", PCID: CASStorageV1, Detail: "corrupt"},
		{Observer: "grace", Event: "cas_corrupt_evidence_recorded", Outcome: "kept", PCID: CASStorageV1, Detail: "evidence"},
		{Observer: "ellen", Event: "compute_context_promised", Outcome: "kept", PCID: CIDComputeV1, Detail: "context"},
		{Observer: "carol", Event: "compute_function_executed", Outcome: "kept", PCID: CIDComputeV1, Detail: "executed"},
		{Observer: "carol", Event: "cid_compute_promised", Outcome: "kept", PCID: CIDComputeV1, Detail: "compute"},
		{Observer: "carol", Event: "compute_result_promised", Outcome: "kept", PCID: CIDComputeV1, Detail: "result"},
		{Observer: "alice", Event: "compute_result_received", Outcome: "kept", PCID: CIDComputeV1, Detail: "result received"},
		{Observer: "dave", Event: "compute_cache_checkpointed", Outcome: "kept", PCID: CIDComputeV1, Detail: "cache"},
		{Observer: "bob", Event: "capability_token_issued", Outcome: "kept", PCID: CASStorageV1, Detail: "issued"},
		{Observer: "alice", Event: "capability_token_received", Outcome: "kept", PCID: CASStorageV1, Detail: "received"},
		{Observer: "bob", Event: "capability_token_redeemed", Outcome: "kept", PCID: CASStorageV1, Detail: "redeemed"},
		{Observer: "bob", Event: "economics_credit_accepted", Outcome: "kept", PCID: CASStorageV1, Detail: "credit"},
		{Observer: "alice", Event: "trust_updated", Outcome: "kept", Detail: "trust"},
		{Observer: "grace", Event: "trust_repair_promise_recorded", Outcome: "kept", PCID: CASStorageV1, Detail: "repair"},
		{Observer: "alice", Event: "promise_envelope_validated", Outcome: "kept", PCID: CASStorageV1, Detail: "exact"},
	})
	summary, err := AnalyzeRun(parentDir)
	if err != nil {
		t.Fatalf("analyze run: %v", err)
	}
	if summary.RunDir != logDir {
		t.Fatalf("summary run dir = %q, want %q", summary.RunDir, logDir)
	}
	if err := ValidateAnalysis(summary); err != nil {
		t.Fatalf("validate analysis: %v", err)
	}
}

func TestValidateAnalysisRejectsRPCDrift(t *testing.T) {
	summary := AnalysisSummary{
		TotalEvents:                   1,
		EventCounts:                   map[string]int{},
		ProtocolCounts:                map[string]int{CASStorageV1: 1, CIDComputeV1: 1},
		RPCDriftCounts:                map[string]int{"alice": 1},
		PlaceholderLiveDecisionCounts: map[string]int{},
	}
	if err := ValidateAnalysis(summary); err == nil {
		t.Fatalf("rpc drift should fail validation")
	}
}

func TestValidateAnalysisRejectsPlaceholderLiveDecision(t *testing.T) {
	summary := AnalysisSummary{
		TotalEvents:                   1,
		EventCounts:                   map[string]int{},
		ProtocolCounts:                map[string]int{CASStorageV1: 1, CIDComputeV1: 1},
		RPCDriftCounts:                map[string]int{},
		PlaceholderLiveDecisionCounts: map[string]int{"alice": 1},
	}
	if err := ValidateAnalysis(summary); err == nil {
		t.Fatalf("placeholder live decision should fail validation")
	}
}

func writePOC13Log(t *testing.T, path string, events []Event) {
	t.Helper()
	logFile, err := os.Create(path)
	if err != nil {
		t.Fatalf("create log: %v", err)
	}
	for _, event := range events {
		eventBytes, marshalErr := jsonMarshalEvent(event)
		if marshalErr != nil {
			t.Fatalf("marshal event: %v", marshalErr)
		}
		if _, writeErr := logFile.Write(append(eventBytes, '\n')); writeErr != nil {
			t.Fatalf("write event: %v", writeErr)
		}
	}
	if err := logFile.Close(); err != nil {
		t.Fatalf("close log: %v", err)
	}
}

func jsonMarshalEvent(event Event) ([]byte, error) {
	return json.Marshal(event)
}
