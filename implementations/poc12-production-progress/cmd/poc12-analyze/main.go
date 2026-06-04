package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

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
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func isShippingEvent(eventName string) bool {
	switch eventName {
	case "package_weighed", "package_weight_received", "shipping_address_promised", "shipping_address_received", "shipping_label_printed", "shipping_label_received", "accounting_updated", "accounting_update_confirmed":
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
