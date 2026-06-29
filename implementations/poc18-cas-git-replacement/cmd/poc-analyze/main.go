// Command poc-analyze checks deterministic POC18 fixture results.
//
// Intent: Analyzer checks remain non-production review aids; they must not imply
// a global monitor or authority. Source: DI-jifuj
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/workspace"
)

// Report records local fixture checks.
type Report struct {
	RunRoot string            `json:"run_root"`
	Pass    bool              `json:"pass"`
	Checks  map[string]string `json:"checks"`
	Objects int               `json:"objects"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "poc-analyze: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	runRoot := flag.String("run-root", "/tmp/wire-lab-poc18-run", "approved POC18 runtime root")
	flag.Parse()
	result, resultErr := readResult(filepath.Join(*runRoot, "result.json"))
	if resultErr != nil {
		return resultErr
	}
	report, analyzeErr := analyze(*runRoot, result)
	if analyzeErr != nil {
		return analyzeErr
	}
	reportPath := filepath.Join(*runRoot, "analysis.json")
	if err := writeJSON(reportPath, report); err != nil {
		return err
	}
	fmt.Printf("pass=%t objects=%d report=%s\n", report.Pass, report.Objects, reportPath)
	if !report.Pass {
		return fmt.Errorf("analysis failed")
	}
	return nil
}

func readResult(path string) (workspace.IngestResult, error) {
	var result workspace.IngestResult
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		return workspace.IngestResult{}, readErr
	}
	if err := json.Unmarshal(content, &result); err != nil {
		return workspace.IngestResult{}, err
	}
	return result, nil
}

func analyze(runRoot string, result workspace.IngestResult) (Report, error) {
	checks := map[string]string{}
	requiredCounts := []string{
		"chunk",
		"chunk_manifest",
		"posix_node:regular",
		"posix_node:directory",
		"posix_node:symlink",
		"posix_node:fifo",
		"hardlink_label",
		"reference_set:directory",
		"reference_set:branch",
		"reference_set:workspace",
		"reference_set:logical_change",
		"reference_set:review_thread",
		"reference_set:release",
		"snapshot",
	}
	pass := true
	for _, countName := range requiredCounts {
		if result.Counts[countName] > 0 {
			checks[countName] = "kept"
			continue
		}
		checks[countName] = "missing"
		pass = false
	}
	if _, err := os.Stat(filepath.Join(result.CheckoutRoot, "README.md")); err == nil {
		checks["checkout_readme"] = "kept"
	} else {
		checks["checkout_readme"] = "missing"
		pass = false
	}
	if _, err := os.Lstat(filepath.Join(result.CheckoutRoot, "README-link.md")); err == nil {
		checks["checkout_symlink"] = "kept"
	} else {
		checks["checkout_symlink"] = "missing"
		pass = false
	}
	cas, openErr := store.Open(result.StoreRoot)
	if openErr != nil {
		return Report{}, openErr
	}
	entries, listErr := cas.List()
	if listErr != nil {
		return Report{}, listErr
	}
	return Report{RunRoot: runRoot, Pass: pass, Checks: checks, Objects: len(entries)}, nil
}

func writeJSON(path string, value any) error {
	content, marshalErr := json.MarshalIndent(value, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	content = append(content, '\n')
	return os.WriteFile(path, content, 0o644)
}
