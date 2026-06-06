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
	RPCDriftCounts                map[string]int `json:"rpc_drift_counts"`
	PlaceholderLiveDecisionCounts map[string]int `json:"placeholder_live_decision_counts"`
	MissingRequiredEventNames     []string       `json:"missing_required_event_names,omitempty"`
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
		if isRPCDrift(event) {
			summary.RPCDriftCounts[event.Observer]++
		}
		if HasPlaceholderLiveDecision(event) {
			summary.PlaceholderLiveDecisionCounts[event.Observer]++
		}
	}
	return scanner.Err()
}

// ValidateAnalysis enforces the first POC13 acceptance gates.
// Intent: CAS and compute evidence should fail fast if it drifts toward RPC,
// loses corrupt-byte handling, or stops recording exact cache evidence. Source:
// DI-notig
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
	if len(failures) > 0 {
		return errors.New(strings.Join(failures, "; "))
	}
	return nil
}

func missingRequiredEvents(summary AnalysisSummary) []string {
	requiredEvents := []string{
		"cas_storage_promised",
		"cas_retention_promised",
		"cas_replication_promised",
		"cas_corrupt_bytes_rejected",
		"cas_corrupt_evidence_recorded",
		"compute_context_promised",
		"cid_compute_promised",
		"compute_result_promised",
		"compute_cache_checkpointed",
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
