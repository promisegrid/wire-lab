package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/decision"
)

// RunSummary is the deterministic JSON report emitted by poc12-analyze.
// Intent: Keep POC12 live runs comparable without committing provider logs or
// Docker volume state. Source: DI-timah
type RunSummary struct {
	RunDir                       string                  `json:"run_dir"`
	TotalEvents                  int                     `json:"total_events"`
	EventCounts                  map[string]int          `json:"event_counts"`
	OutcomeCounts                map[string]int          `json:"outcome_counts"`
	AgentCounts                  map[string]int          `json:"agent_counts"`
	FailureCounts                map[string]int          `json:"failure_counts"`
	ShippingCounts               map[string]int          `json:"shipping_counts"`
	RelationshipTransitionCounts map[string]int          `json:"relationship_transition_counts"`
	LocalResourceCounts          map[string]int          `json:"local_resource_counts"`
	ResourceTrustCouplingCounts  map[string]int          `json:"resource_trust_coupling_counts"`
	MonitorReport                *decision.MonitorReport `json:"monitor_report,omitempty"`
}

// AcceptanceCriteria describes the evidence gates for a clean POC12 regression
// run.
// Intent: Keep the current clean-run expectations executable so a wrong log
// path, missing shipping workflow, lost non-commitment suppression, or renewed
// resource/trust coupling fails loudly instead of producing a plausible but
// empty report. Source: DI-jidah
type AcceptanceCriteria struct {
	MinTotalEvents                    int
	RequireMonitorReport              bool
	MinMonitorScore                   int
	RequireFulfillmentWorkflow        bool
	RequireAccountingUpdated          bool
	RequireAccountingDuplicate        bool
	RequireNotPromisedSuppression     bool
	RequireNoResourceTrustCoupling    bool
	RequireAccountingDuplicateConfirm bool
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: poc12-analyze RUN_DIR\n")
		os.Exit(2)
	}
	summary, err := analyzeRun(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "poc12-analyze: %v\n", err)
		os.Exit(1)
	}
	if err := validateSummary(summary, cleanRegressionCriteria()); err != nil {
		fmt.Fprintf(os.Stderr, "poc12-analyze: acceptance criteria failed: %v\n", err)
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		fmt.Fprintf(os.Stderr, "poc12-analyze: encode summary: %v\n", err)
		os.Exit(1)
	}
}

// analyzeRun summarizes one POC12 run directory containing agent JSONL files and
// an optional monitor-report.json.
// Intent: The analyzer treats logs as evidence to count, not as authority over
// agents; it does not mutate run state. Source: DI-timah
func analyzeRun(runDir string) (RunSummary, error) {
	logDir, resolveErr := resolveRunLogDir(runDir)
	if resolveErr != nil {
		return RunSummary{}, resolveErr
	}
	summary := RunSummary{
		RunDir:                       logDir,
		EventCounts:                  make(map[string]int),
		OutcomeCounts:                make(map[string]int),
		AgentCounts:                  make(map[string]int),
		FailureCounts:                make(map[string]int),
		ShippingCounts:               make(map[string]int),
		RelationshipTransitionCounts: make(map[string]int),
		LocalResourceCounts:          make(map[string]int),
		ResourceTrustCouplingCounts:  make(map[string]int),
	}
	logPaths := jsonlLogPaths(logDir)
	sort.Strings(logPaths)
	for _, logPath := range logPaths {
		if err := summarizeLog(logPath, &summary); err != nil {
			return RunSummary{}, err
		}
	}
	report, reportErr := readMonitorReport(filepath.Join(logDir, "monitor-report.json"))
	if reportErr != nil {
		return RunSummary{}, reportErr
	}
	summary.MonitorReport = report
	return summary, nil
}

// resolveRunLogDir accepts either the concrete JSONL log directory or the
// parent run directory that contains the `run/` evidence subdirectory.
// Intent: Operators should not be able to accidentally analyze the parent path
// and receive a misleading zero-event success report. Source: DI-jidah
func resolveRunLogDir(inputPath string) (string, error) {
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
	return "", fmt.Errorf("%s contains no JSONL logs; pass the log directory or its parent run directory", inputPath)
}

func jsonlLogPaths(runDir string) []string {
	logPaths, globErr := filepath.Glob(filepath.Join(runDir, "*.jsonl"))
	if globErr != nil {
		return nil
	}
	return logPaths
}

func cleanRegressionCriteria() AcceptanceCriteria {
	return AcceptanceCriteria{
		MinTotalEvents:                    1,
		RequireMonitorReport:              true,
		MinMonitorScore:                   4,
		RequireFulfillmentWorkflow:        true,
		RequireAccountingUpdated:          true,
		RequireAccountingDuplicate:        true,
		RequireNotPromisedSuppression:     true,
		RequireNoResourceTrustCoupling:    true,
		RequireAccountingDuplicateConfirm: true,
	}
}

// validateSummary checks the POC12 clean-run evidence gates without mutating
// logs or relationship state.
// Intent: The acceptance contract should catch regression evidence directly in
// analyzer output instead of requiring the operator to interpret counts by eye.
// Source: DI-jidah
func validateSummary(summary RunSummary, criteria AcceptanceCriteria) error {
	var failures []string
	if summary.TotalEvents < criteria.MinTotalEvents {
		failures = append(failures, fmt.Sprintf("total_events=%d below minimum %d", summary.TotalEvents, criteria.MinTotalEvents))
	}
	if criteria.RequireFulfillmentWorkflow && summary.EventCounts["fulfillment_workflow_completed"] != 1 {
		failures = append(failures, fmt.Sprintf("fulfillment_workflow_completed=%d want 1", summary.EventCounts["fulfillment_workflow_completed"]))
	}
	if criteria.RequireAccountingUpdated && summary.ShippingCounts["accounting_updated"] != 1 {
		failures = append(failures, fmt.Sprintf("accounting_updated=%d want 1", summary.ShippingCounts["accounting_updated"]))
	}
	if criteria.RequireAccountingDuplicate && summary.ShippingCounts["accounting_update_duplicate"] != 1 {
		failures = append(failures, fmt.Sprintf("accounting_update_duplicate=%d want 1", summary.ShippingCounts["accounting_update_duplicate"]))
	}
	if criteria.RequireAccountingDuplicateConfirm && summary.ShippingCounts["accounting_update_duplicate_confirmed"] != 1 {
		failures = append(failures, fmt.Sprintf("accounting_update_duplicate_confirmed=%d want 1", summary.ShippingCounts["accounting_update_duplicate_confirmed"]))
	}
	if criteria.RequireNotPromisedSuppression && summary.EventCounts["promise_not_promised_suppressed"] == 0 {
		failures = append(failures, "promise_not_promised_suppressed=0 want >0")
	}
	if criteria.RequireNoResourceTrustCoupling && len(summary.ResourceTrustCouplingCounts) != 0 {
		failures = append(failures, fmt.Sprintf("resource_trust_coupling_counts=%v want empty", summary.ResourceTrustCouplingCounts))
	}
	if criteria.RequireMonitorReport {
		if summary.MonitorReport == nil {
			failures = append(failures, "monitor_report missing")
		} else {
			failures = append(failures, monitorScoreFailures(*summary.MonitorReport, criteria.MinMonitorScore)...)
		}
	}
	if len(failures) > 0 {
		return errors.New(strings.Join(failures, "; "))
	}
	return nil
}

func monitorScoreFailures(report decision.MonitorReport, minimum int) []string {
	scores := map[string]int{
		"promise_theory_fit":      report.PromiseTheoryFit,
		"autonomy":                report.Autonomy,
		"protocol_validity":       report.ProtocolValidity,
		"local_trust_correctness": report.LocalTrustCorrectness,
		"imposition_avoidance":    report.ImpositionAvoidance,
	}
	names := make([]string, 0, len(scores))
	for name := range scores {
		names = append(names, name)
	}
	sort.Strings(names)
	var failures []string
	for _, name := range names {
		if scores[name] < minimum {
			failures = append(failures, fmt.Sprintf("%s=%d below minimum %d", name, scores[name], minimum))
		}
	}
	return failures
}

func summarizeLog(logPath string, summary *RunSummary) error {
	logFile, openErr := os.Open(logPath)
	if openErr != nil {
		return openErr
	}
	defer func() {
		closeErr := logFile.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "poc12-analyze: close %s: %v\n", logPath, closeErr)
		}
	}()
	scanner := bufio.NewScanner(logFile)
	var previousEvent decision.Event
	hasPreviousEvent := false
	for scanner.Scan() {
		var event decision.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return fmt.Errorf("%s: %w", logPath, err)
		}
		summary.TotalEvents++
		summary.EventCounts[event.Event]++
		summary.OutcomeCounts[event.Outcome]++
		summary.AgentCounts[event.Observer]++
		if event.Outcome != "kept" {
			summary.FailureCounts[event.Event]++
		}
		if isShippingEvent(event.Event) {
			summary.ShippingCounts[event.Event]++
		}
		if isRelationshipTransition(event.Event) {
			summary.RelationshipTransitionCounts[event.Event]++
		}
		if isLocalResourceEvent(event) {
			summary.LocalResourceCounts[event.Event]++
		}
		if hasPreviousEvent && isResourceTrustCoupling(previousEvent, event) {
			summary.ResourceTrustCouplingCounts[event.Observer]++
		}
		previousEvent = event
		hasPreviousEvent = true
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func isShippingEvent(eventName string) bool {
	switch eventName {
	case "package_weighed", "package_weight_received", "shipping_address_promised", "shipping_address_received", "shipping_label_printed", "shipping_label_received", "accounting_updated", "accounting_update_confirmed", "accounting_update_duplicate", "accounting_update_duplicate_confirmed", "printer_capability_issued", "printer_capability_received", "printer_port_printed", "printer_port_print_confirmed":
		return true
	default:
		return false
	}
}

func isRelationshipTransition(eventName string) bool {
	switch eventName {
	case "direct_peer_added", "direct_peer_removed", "direct_peer_unchanged":
		return true
	default:
		return false
	}
}

// isLocalResourceEvent recognizes evidence about the observing app's own
// budget/capacity state.
// Intent: Analyzer output should make it visible if local resource exhaustion is
// accidentally coupled back into peer trust transitions. Source: DI-vujob
func isLocalResourceEvent(event decision.Event) bool {
	if event.Event == "local_resource_exhausted" {
		return true
	}
	if event.Event == "promise_withheld" && (strings.Contains(event.Detail, "budget exhausted") || strings.Contains(event.Detail, "capacity exhausted")) {
		return true
	}
	if event.Event == "resource_promise_broken" && strings.Contains(event.Detail, "local") {
		return true
	}
	return false
}

// isResourceTrustCoupling reports suspicious adjacency between local resource
// exhaustion and peer relationship transitions for one observer.
// Intent: This is a regression tripwire; local budget/capacity state should not
// immediately mutate peer trust without intervening promise evidence. Source:
// DI-vujob
func isResourceTrustCoupling(previousEvent, event decision.Event) bool {
	if previousEvent.Observer != event.Observer {
		return false
	}
	return (isRelationshipTransition(previousEvent.Event) && isLocalResourceEvent(event)) ||
		(isLocalResourceEvent(previousEvent) && isRelationshipTransition(event.Event))
}

func readMonitorReport(reportPath string) (*decision.MonitorReport, error) {
	reportBytes, readErr := os.ReadFile(reportPath)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return nil, nil
		}
		return nil, readErr
	}
	var report decision.MonitorReport
	if err := json.Unmarshal(reportBytes, &report); err != nil {
		return nil, err
	}
	return &report, nil
}
