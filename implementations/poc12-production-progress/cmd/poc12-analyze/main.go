package main

import (
	"bufio"
	"encoding/json"
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
	summary := RunSummary{
		RunDir:                       runDir,
		EventCounts:                  make(map[string]int),
		OutcomeCounts:                make(map[string]int),
		AgentCounts:                  make(map[string]int),
		FailureCounts:                make(map[string]int),
		ShippingCounts:               make(map[string]int),
		RelationshipTransitionCounts: make(map[string]int),
		LocalResourceCounts:          make(map[string]int),
		ResourceTrustCouplingCounts:  make(map[string]int),
	}
	logPaths, globErr := filepath.Glob(filepath.Join(runDir, "*.jsonl"))
	if globErr != nil {
		return RunSummary{}, globErr
	}
	sort.Strings(logPaths)
	for _, logPath := range logPaths {
		if err := summarizeLog(logPath, &summary); err != nil {
			return RunSummary{}, err
		}
	}
	report, reportErr := readMonitorReport(filepath.Join(runDir, "monitor-report.json"))
	if reportErr != nil {
		return RunSummary{}, reportErr
	}
	summary.MonitorReport = report
	return summary, nil
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
