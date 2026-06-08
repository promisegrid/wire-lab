package poc13

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// AnalysisSummary is the deterministic POC13 analyzer report.
type AnalysisSummary struct {
	RunDir                        string         `json:"run_dir"`
	TotalEvents                   int            `json:"total_events"`
	EventCounts                   map[string]int `json:"event_counts"`
	OutcomeCounts                 map[string]int `json:"outcome_counts"`
	AgentCounts                   map[string]int `json:"agent_counts"`
	ProtocolCounts                map[string]int `json:"protocol_counts"`
	TrustDrivenChoiceCounts       map[string]int `json:"trust_driven_choice_counts"`
	EconomicsCounts               map[string]int `json:"economics_counts"`
	VerificationCounts            map[string]int `json:"verification_counts"`
	ReplicaRecoveryCounts         map[string]int `json:"replica_recovery_counts"`
	RPCDriftCounts                map[string]int `json:"rpc_drift_counts"`
	PlaceholderLiveDecisionCounts map[string]int `json:"placeholder_live_decision_counts"`
	ScoreReport                   ScoreReport    `json:"score_report"`
	MissingRequiredEventNames     []string       `json:"missing_required_event_names,omitempty"`
}

// ScoreReport is a compact POC13 fitness report for humans and scripts.
// Intent: The analyzer should report evidence quality directly, not force
// readers to infer POC health from raw event counts only. Source: DI-lupag
type ScoreReport struct {
	Overall      int      `json:"overall"`
	Transport    int      `json:"transport"`
	Storage      int      `json:"storage"`
	Compute      int      `json:"compute"`
	Economics    int      `json:"economics"`
	Trust        int      `json:"trust"`
	Verification int      `json:"verification"`
	Replica      int      `json:"replica"`
	Concerns     []string `json:"concerns,omitempty"`
}

// AnalyzeRun summarizes one POC13 JSONL directory or its parent run directory.
func AnalyzeRun(inputPath string) (AnalysisSummary, error) {
	logDir, resolveErr := ResolveRunLogDir(inputPath)
	if resolveErr != nil {
		return AnalysisSummary{}, resolveErr
	}
	summary := AnalysisSummary{
		RunDir:                        logDir,
		EventCounts:                   make(map[string]int),
		OutcomeCounts:                 make(map[string]int),
		AgentCounts:                   make(map[string]int),
		ProtocolCounts:                make(map[string]int),
		TrustDrivenChoiceCounts:       make(map[string]int),
		EconomicsCounts:               make(map[string]int),
		VerificationCounts:            make(map[string]int),
		ReplicaRecoveryCounts:         make(map[string]int),
		RPCDriftCounts:                make(map[string]int),
		PlaceholderLiveDecisionCounts: make(map[string]int),
	}
	logPaths := jsonlLogPaths(logDir)
	sort.Strings(logPaths)
	for _, logPath := range logPaths {
		if err := summarizePOC13Log(logPath, &summary); err != nil {
			return AnalysisSummary{}, err
		}
	}
	summary.MissingRequiredEventNames = missingRequiredEvents(summary)
	summary.ScoreReport = computeScores(summary)
	return summary, nil
}

// ResolveRunLogDir accepts either `/run/poc13/<run_id>` or the nested
// `/run/poc13/<run_id>/run` JSONL path.
func ResolveRunLogDir(inputPath string) (string, error) {
	if _, statErr := os.Stat(inputPath); statErr != nil {
		return "", statErr
	}
	if len(jsonlLogPaths(inputPath)) > 0 {
		return inputPath, nil
	}
	nestedRunPath := filepath.Join(inputPath, "run")
	if _, statErr := os.Stat(nestedRunPath); statErr == nil && len(jsonlLogPaths(nestedRunPath)) > 0 {
		return nestedRunPath, nil
	}
	return "", fmt.Errorf("%s contains no POC13 JSONL logs", inputPath)
}

func jsonlLogPaths(logDir string) []string {
	logPaths, globErr := filepath.Glob(filepath.Join(logDir, "*.jsonl"))
	if globErr != nil {
		return nil
	}
	return logPaths
}

func summarizePOC13Log(logPath string, summary *AnalysisSummary) error {
	logFile, openErr := os.Open(logPath)
	if openErr != nil {
		return openErr
	}
	defer func() {
		closeErr := logFile.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "poc13-analyze: close %s: %v\n", logPath, closeErr)
		}
	}()
	scanner := bufio.NewScanner(logFile)
	for scanner.Scan() {
		var event Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return fmt.Errorf("%s: %w", logPath, err)
		}
		summary.TotalEvents++
		summary.EventCounts[event.Event]++
		summary.OutcomeCounts[event.Outcome]++
		summary.AgentCounts[event.Observer]++
		if event.PCID != "" {
			summary.ProtocolCounts[event.PCID]++
		}
		if event.Event == "trust_driven_peer_choice" {
			summary.TrustDrivenChoiceCounts[event.Observer]++
		}
		if strings.HasPrefix(event.Event, "economics_") {
			summary.EconomicsCounts[event.Observer]++
		}
		if strings.Contains(event.Event, "verification") || event.Event == "compute_result_locally_verified" || event.Event == "compute_result_peer_verified" {
			summary.VerificationCounts[event.Observer]++
		}
		if strings.HasPrefix(event.Event, "replica_recovery") || event.Event == "cas_replica_serve_promised" {
			summary.ReplicaRecoveryCounts[event.Observer]++
		}
		if isRPCDrift(event) {
			summary.RPCDriftCounts[event.Observer]++
		}
		if HasPlaceholderLiveDecision(event) {
			summary.PlaceholderLiveDecisionCounts[event.Observer]++
		}
	}
	return scanner.Err()
}

// ValidateAnalysis enforces the POC13 acceptance gates.
// Intent: CAS and compute evidence should fail fast if it drifts toward RPC,
// loses corrupt-byte handling, stops recording exact cache evidence, or fails
// to prove TCP delivery, concrete storage/retrieval, dynamic compute, trust,
// economics, repair, capability-token, replica-recovery, verification, and
// score/report evidence. Source: DI-notig; DI-fumol; DI-lupag
func ValidateAnalysis(summary AnalysisSummary) error {
	var failures []string
	if summary.TotalEvents == 0 {
		failures = append(failures, "total_events=0")
	}
	if len(summary.MissingRequiredEventNames) > 0 {
		failures = append(failures, "missing required events: "+strings.Join(summary.MissingRequiredEventNames, ", "))
	}
	if len(summary.RPCDriftCounts) > 0 {
		failures = append(failures, fmt.Sprintf("rpc_drift_counts=%v want empty", summary.RPCDriftCounts))
	}
	if len(summary.PlaceholderLiveDecisionCounts) > 0 {
		failures = append(failures, fmt.Sprintf("placeholder_live_decision_counts=%v want empty", summary.PlaceholderLiveDecisionCounts))
	}
	if summary.ProtocolCounts[CASStorageV1] == 0 {
		failures = append(failures, "cas_storage_v1 evidence missing")
	}
	if summary.ProtocolCounts[CIDComputeV1] == 0 {
		failures = append(failures, "cid_compute_v1 evidence missing")
	}
	if summary.ScoreReport.Overall < 4 {
		failures = append(failures, fmt.Sprintf("score_report.overall=%d want >=4", summary.ScoreReport.Overall))
	}
	if len(failures) > 0 {
		return errors.New(strings.Join(failures, "; "))
	}
	return nil
}

func missingRequiredEvents(summary AnalysisSummary) []string {
	requiredEvents := []string{
		"runtime_readiness_promised",
		"peer_readiness_observed",
		"runtime_done_promised",
		"tcp_message_sent",
		"tcp_message_received",
		"cas_storage_promised",
		"cas_retention_promised",
		"cas_replication_promised",
		"cas_bytes_stored",
		"cas_bytes_retrieved",
		"cas_replica_stored",
		"cas_replica_serve_promised",
		"primary_storage_unavailable",
		"cas_corrupt_bytes_rejected",
		"cas_corrupt_evidence_recorded",
		"compute_context_promised",
		"compute_function_executed",
		"cid_compute_promised",
		"compute_result_promised",
		"compute_result_received",
		"compute_cache_checkpointed",
		"compute_result_locally_verified",
		"compute_result_peer_verified",
		"compute_verification_received",
		"capability_token_issued",
		"capability_token_received",
		"capability_token_redeemed",
		"trust_driven_peer_choice",
		"economics_credit_accepted",
		"economics_price_refused",
		"economics_capacity_reserved",
		"economics_credits_spent",
		"economics_credits_earned",
		"replica_recovery_requested",
		"replica_recovery_succeeded",
		"trust_updated",
		"trust_repair_promise_recorded",
		"promise_envelope_validated",
	}
	var missing []string
	for _, eventName := range requiredEvents {
		if summary.EventCounts[eventName] == 0 {
			missing = append(missing, eventName)
		}
	}
	return missing
}

func computeScores(summary AnalysisSummary) ScoreReport {
	scores := ScoreReport{}
	addScore(&scores.Transport, summary.EventCounts["tcp_message_sent"] > 0 && summary.EventCounts["tcp_message_received"] > 0)
	addScore(&scores.Storage, summary.EventCounts["cas_bytes_stored"] > 0 && summary.EventCounts["cas_bytes_retrieved"] > 0)
	addScore(&scores.Compute, summary.EventCounts["compute_function_executed"] > 0 && summary.EventCounts["compute_result_received"] > 0)
	addScore(&scores.Economics, summary.EventCounts["economics_price_refused"] > 0 && summary.EventCounts["economics_credits_spent"] > 0 && summary.EventCounts["economics_credits_earned"] > 0)
	addScore(&scores.Trust, summary.EventCounts["trust_updated"] > 0 && summary.EventCounts["trust_driven_peer_choice"] > 0)
	addScore(&scores.Verification, summary.EventCounts["compute_result_locally_verified"] > 0 && summary.EventCounts["compute_result_peer_verified"] > 0)
	addScore(&scores.Replica, summary.EventCounts["replica_recovery_succeeded"] > 0)
	scores.Overall = (scores.Transport + scores.Storage + scores.Compute + scores.Economics + scores.Trust + scores.Verification + scores.Replica) / 7
	if len(summary.MissingRequiredEventNames) > 0 {
		scores.Concerns = append(scores.Concerns, "missing required events: "+strings.Join(summary.MissingRequiredEventNames, ", "))
	}
	if len(summary.RPCDriftCounts) > 0 {
		scores.Concerns = append(scores.Concerns, "RPC drift detected")
	}
	if len(summary.PlaceholderLiveDecisionCounts) > 0 {
		scores.Concerns = append(scores.Concerns, "placeholder live decisions detected")
	}
	return scores
}

func addScore(score *int, kept bool) {
	if kept {
		*score = 5
		return
	}
	*score = 1
}

func isRPCDrift(event Event) bool {
	text := strings.ToLower(event.Event + " " + event.Detail)
	for _, badTerm := range []string{"rpc", "command", "permission", "authorization", "conformance", "enforce"} {
		if strings.Contains(text, badTerm) {
			return true
		}
	}
	return false
}

// HasPlaceholderLiveDecision identifies live provider evidence that did not
// include a real provider judgment.
// Intent: A live POC13 run should not pass analyzer gates if provider calls
// succeeded but produced only placeholder text. Source: DI-lasuh
func HasPlaceholderLiveDecision(event Event) bool {
	if event.Event != "llm_decision_live" {
		return false
	}
	return strings.Contains(strings.ToLower(event.Detail), "provider returned no output_text")
}
