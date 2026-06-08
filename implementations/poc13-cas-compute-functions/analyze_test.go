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
		{Observer: "alice", Event: "runtime_readiness_promised", Outcome: "kept", Detail: "ready"},
		{Observer: "alice", Event: "peer_readiness_observed", Outcome: "kept", Detail: "peer ready"},
		{Observer: "alice", Event: "runtime_done_promised", Outcome: "kept", Detail: "done"},
		{Observer: "alice", Event: "tcp_message_sent", Outcome: "kept", PCID: CASStorageV1, Detail: "sent"},
		{Observer: "bob", Event: "tcp_message_received", Outcome: "kept", PCID: CASStorageV1, Detail: "received"},
		{Observer: "alice", Event: "tcp_message_send_failed", Outcome: "non_commitment", PCID: CASStorageV1, Detail: "bob unreachable"},
		{Observer: "alice", Event: "network_outage_variant_selected", Outcome: "kept", PCID: CASStorageV1, Detail: "variant"},
		{Observer: "bob", Event: "cas_storage_promised", Outcome: "kept", PCID: CASStorageV1, Detail: "storage"},
		{Observer: "bob", Event: "cas_retention_promised", Outcome: "kept", PCID: CASStorageV1, Detail: "retention"},
		{Observer: "frank", Event: "cas_replication_promised", Outcome: "kept", PCID: CASStorageV1, Detail: "replica"},
		{Observer: "bob", Event: "cas_bytes_stored", Outcome: "kept", PCID: CASStorageV1, Detail: "stored"},
		{Observer: "alice", Event: "cas_bytes_retrieved", Outcome: "kept", PCID: CASStorageV1, Detail: "retrieved"},
		{Observer: "frank", Event: "cas_replica_stored", Outcome: "kept", PCID: CASStorageV1, Detail: "replica stored"},
		{Observer: "frank", Event: "cas_replica_serve_promised", Outcome: "kept", PCID: CASStorageV1, Detail: "replica served"},
		{Observer: "alice", Event: "cas_multi_object_pressure", Outcome: "kept", PCID: CASStorageV1, Detail: "multi"},
		{Observer: "alice", Event: "primary_storage_unavailable", Outcome: "non_commitment", PCID: CASStorageV1, Detail: "primary unavailable"},
		{Observer: "frank", Event: "replica_capability_token_issued", Outcome: "kept", PCID: CASStorageV1, Detail: "replica issued"},
		{Observer: "alice", Event: "replica_capability_token_received", Outcome: "kept", PCID: CASStorageV1, Detail: "replica received"},
		{Observer: "frank", Event: "replica_capability_token_redeemed", Outcome: "kept", PCID: CASStorageV1, Detail: "replica redeemed"},
		{Observer: "frank", Event: "capability_token_ttl_promised", Outcome: "kept", PCID: CASStorageV1, Detail: "ttl"},
		{Observer: "alice", Event: "capability_token_ttl_observed", Outcome: "kept", PCID: CASStorageV1, Detail: "ttl observed"},
		{Observer: "frank", Event: "capability_token_expired", Outcome: "non_commitment", PCID: CASStorageV1, Detail: "expired"},
		{Observer: "frank", Event: "capability_token_revoked", Outcome: "non_commitment", PCID: CASStorageV1, Detail: "revoked"},
		{Observer: "alice", Event: "capability_token_renewal_requested", Outcome: "kept", PCID: CASStorageV1, Detail: "renew request"},
		{Observer: "alice", Event: "capability_token_renewed", Outcome: "kept", PCID: CASStorageV1, Detail: "renewed"},
		{Observer: "mallory", Event: "cas_corrupt_bytes_rejected", Outcome: "malformed", PCID: CASStorageV1, Detail: "corrupt"},
		{Observer: "grace", Event: "cas_corrupt_evidence_recorded", Outcome: "kept", PCID: CASStorageV1, Detail: "evidence"},
		{Observer: "ellen", Event: "compute_context_promised", Outcome: "kept", PCID: CIDComputeV1, Detail: "context"},
		{Observer: "carol", Event: "compute_function_executed", Outcome: "kept", PCID: CIDComputeV1, Detail: "executed"},
		{Observer: "carol", Event: "compute_alternate_function_executed", Outcome: "kept", PCID: CIDComputeV1, Detail: "sum"},
		{Observer: "alice", Event: "compute_followup_function_requested", Outcome: "kept", PCID: CIDComputeV1, Detail: "followup"},
		{Observer: "carol", Event: "cid_compute_promised", Outcome: "kept", PCID: CIDComputeV1, Detail: "compute"},
		{Observer: "carol", Event: "compute_result_promised", Outcome: "kept", PCID: CIDComputeV1, Detail: "result"},
		{Observer: "alice", Event: "compute_result_received", Outcome: "kept", PCID: CIDComputeV1, Detail: "result received"},
		{Observer: "carol", Event: "compute_bad_result_promised", Outcome: "malformed", PCID: CIDComputeV1, Detail: "bad result"},
		{Observer: "dave", Event: "compute_cache_checkpointed", Outcome: "kept", PCID: CIDComputeV1, Detail: "cache"},
		{Observer: "dave", Event: "compute_cache_miss", Outcome: "non_commitment", PCID: CIDComputeV1, Detail: "miss"},
		{Observer: "alice", Event: "compute_cache_miss_observed", Outcome: "non_commitment", PCID: CIDComputeV1, Detail: "miss observed"},
		{Observer: "dave", Event: "compute_cache_hit", Outcome: "kept", PCID: CIDComputeV1, Detail: "hit"},
		{Observer: "alice", Event: "compute_cache_reused", Outcome: "kept", PCID: CIDComputeV1, Detail: "reuse"},
		{Observer: "alice", Event: "compute_result_locally_verified", Outcome: "kept", PCID: CIDComputeV1, Detail: "local verify"},
		{Observer: "alice", Event: "compute_result_locally_rejected", Outcome: "malformed", PCID: CIDComputeV1, Detail: "local rejected"},
		{Observer: "grace", Event: "compute_result_peer_verified", Outcome: "kept", PCID: CIDComputeV1, Detail: "peer verify"},
		{Observer: "dave", Event: "compute_result_peer_rejected", Outcome: "malformed", PCID: CIDComputeV1, Detail: "peer rejected"},
		{Observer: "grace", Event: "compute_verifier_disagreement", Outcome: "non_commitment", PCID: CIDComputeV1, Detail: "disagree"},
		{Observer: "alice", Event: "compute_disagreement_resolved_locally", Outcome: "kept", PCID: EvidenceReportV1, Detail: "resolved"},
		{Observer: "alice", Event: "compute_verification_received", Outcome: "kept", PCID: EvidenceReportV1, Detail: "verify received"},
		{Observer: "alice", Event: "evidence_report_received", Outcome: "kept", PCID: EvidenceReportV1, Detail: "evidence report"},
		{Observer: "bob", Event: "capability_token_issued", Outcome: "kept", PCID: CASStorageV1, Detail: "issued"},
		{Observer: "alice", Event: "capability_token_received", Outcome: "kept", PCID: CASStorageV1, Detail: "received"},
		{Observer: "bob", Event: "capability_token_redeemed", Outcome: "kept", PCID: CASStorageV1, Detail: "redeemed"},
		{Observer: "alice", Event: "trust_driven_peer_choice", Outcome: "kept", PCID: CASStorageV1, Detail: "choice"},
		{Observer: "alice", Event: "persisted_trust_history_loaded", Outcome: "kept", Detail: "history"},
		{Observer: "alice", Event: "dynamic_peer_choice_from_persisted_trust", Outcome: "kept", PCID: CASStorageV1, Detail: "dynamic"},
		{Observer: "bob", Event: "economics_credit_accepted", Outcome: "kept", PCID: CASStorageV1, Detail: "credit"},
		{Observer: "bob", Event: "economics_price_refused", Outcome: "non_commitment", PCID: CASStorageV1, Detail: "price"},
		{Observer: "bob", Event: "economics_capacity_reserved", Outcome: "kept", PCID: CASStorageV1, Detail: "capacity"},
		{Observer: "carol", Event: "economics_capacity_refused", Outcome: "non_commitment", PCID: CIDComputeV1, Detail: "capacity refused"},
		{Observer: "alice", Event: "economics_credits_spent", Outcome: "kept", PCID: CASStorageV1, Detail: "spent"},
		{Observer: "bob", Event: "economics_credits_earned", Outcome: "kept", PCID: CASStorageV1, Detail: "earned"},
		{Observer: "alice", Event: "economics_payment_withheld", Outcome: "non_commitment", PCID: CIDComputeV1, Detail: "withheld"},
		{Observer: "alice", Event: "replica_recovery_requested", Outcome: "kept", PCID: CASStorageV1, Detail: "requested"},
		{Observer: "alice", Event: "replica_recovery_succeeded", Outcome: "kept", PCID: CASStorageV1, Detail: "succeeded"},
		{Observer: "alice", Event: "trust_updated", Outcome: "kept", Detail: "trust"},
		{Observer: "grace", Event: "trust_repair_promise_recorded", Outcome: "kept", PCID: CASStorageV1, Detail: "repair"},
		{Observer: "grace", Event: "unknown_pcid_not_promised", Outcome: "non_commitment", PCID: "unknown", Detail: "unknown"},
		{Observer: "grace", Event: "promise_variant_not_promised", Outcome: "non_commitment", PCID: CASStorageV1, Detail: "unsupported"},
		{Observer: "mallory", Event: "bad_proof_sent", Outcome: "malformed", PCID: CASStorageV1, Detail: "bad proof"},
		{Observer: "grace", Event: "bad_proof_rejected", Outcome: "malformed", PCID: CASStorageV1, Detail: "bad proof"},
		{Observer: "grace", Event: "key_rotation_promise_recorded", Outcome: "kept", PCID: EvidenceReportV1, Detail: "rotate"},
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
	if summary.ScoreReport.Overall != 5 {
		t.Fatalf("overall score = %d, want 5", summary.ScoreReport.Overall)
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
