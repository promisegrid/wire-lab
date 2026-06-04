package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"promisegrid.dev/wire-lab/implementations/poc11-adaptive-trust-tcp/decision"
)

// RunSummary is the deterministic JSON report emitted by poc11-analyze.
// Intent: Keep POC11 live runs comparable without committing provider logs or
// Docker volume state. Source: DI-nanud
type RunSummary struct {
	RunDir        string                  `json:"run_dir"`
	TotalEvents   int                     `json:"total_events"`
	EventCounts   map[string]int          `json:"event_counts"`
	OutcomeCounts map[string]int          `json:"outcome_counts"`
	AgentCounts   map[string]int          `json:"agent_counts"`
	FailureCounts map[string]int          `json:"failure_counts"`
	MonitorReport *decision.MonitorReport `json:"monitor_report,omitempty"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: poc11-analyze RUN_DIR\n")
		os.Exit(2)
	}
	summary, err := analyzeRun(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "poc11-analyze: %v\n", err)
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		fmt.Fprintf(os.Stderr, "poc11-analyze: encode summary: %v\n", err)
		os.Exit(1)
	}
}

// analyzeRun summarizes one POC11 run directory containing agent JSONL files and
// an optional monitor-report.json.
// Intent: The analyzer treats logs as evidence to count, not as authority over
// agents; it does not mutate run state. Source: DI-nanud
func analyzeRun(runDir string) (RunSummary, error) {
	summary := RunSummary{
		RunDir:        runDir,
		EventCounts:   make(map[string]int),
		OutcomeCounts: make(map[string]int),
		AgentCounts:   make(map[string]int),
		FailureCounts: make(map[string]int),
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
			fmt.Fprintf(os.Stderr, "poc11-analyze: close %s: %v\n", logPath, closeErr)
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
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
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
