package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var requiredHeadings = []string{
	"## Result ID",
	"## Scenario",
	"## Simulation",
	"## Runner",
	"## Prompt / Procedure",
	"## Observed Behavior",
	"## Verdict",
	"## Evidence Links",
	"## Open Questions",
	"## Handoff Target",
	"## Authority Boundary",
}

func runValidate(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	manifest := fs.String("manifest", "", "optional manifest path")
	model := fs.String("model", "", "optional model filter")
	timestamp := fs.String("timestamp", "", "optional timestamp filter")
	strictMatrix := fs.Bool("strict-matrix", false, "require each result path in scenario MATRIX.md")
	allowPrototype := fs.Bool("allow-prototype", false, "allow scripted prototype results")
	allowMissing := fs.Bool("allow-missing", false, "skip missing manifest result files")
	maxErrors := fs.Int("max-errors", 50, "maximum printed errors")
	if err := fs.Parse(args); err != nil {
		return err
	}
	repo, err := openRepo(*repoRoot)
	if err != nil {
		return err
	}
	var targets []string
	var missing []string
	if *manifest != "" {
		cells, err := readManifest(repo, *manifest)
		if err != nil {
			return err
		}
		for _, cell := range cells {
			if cell.ResultPath == "" || strings.Contains(cell.ResultPath, timestampPlaceholder) {
				missing = append(missing, fmt.Sprintf("%s: manifest row has no concrete result path", cell.CellID))
				continue
			}
			path := repo.Abs(cell.ResultPath)
			if _, err := os.Stat(path); err != nil {
				if !*allowMissing {
					missing = append(missing, fmt.Sprintf("%s: missing result file: %s", cell.CellID, repo.Rel(path)))
				}
				continue
			}
			targets = append(targets, path)
		}
	} else {
		targets, err = findResultFiles(repo, *model, *timestamp, *allowPrototype)
		if err != nil {
			return err
		}
	}
	sort.Strings(targets)
	bad := len(missing)
	printed := 0
	for _, msg := range missing {
		if printed < *maxErrors {
			fmt.Fprintln(stdout, msg)
			printed++
		}
	}
	for _, path := range targets {
		issues := validateResultFile(repo, path, *strictMatrix, *allowPrototype)
		if len(issues) == 0 {
			continue
		}
		bad++
		if printed < *maxErrors {
			fmt.Fprintf(stdout, "%s:\n", repo.Rel(path))
			for _, issue := range issues {
				fmt.Fprintf(stdout, "  - %s\n", issue)
			}
			printed++
		}
	}
	if len(targets) == 0 && len(missing) == 0 && !(*manifest != "" && *allowMissing) {
		return fmt.Errorf("no result files matched selection")
	}
	fmt.Fprintf(stdout, "validated=%d failed=%d\n", len(targets), bad)
	if bad > 0 {
		return fmt.Errorf("validation failed")
	}
	return nil
}

func findResultFiles(repo Repo, model string, timestamp string, allowPrototype bool) ([]string, error) {
	var paths []string
	root := repo.Path("results")
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" || path == repo.Path("results", "README.md") {
			return nil
		}
		relParts := strings.Split(repo.Rel(path), "/")
		if len(relParts) != 5 || relParts[0] != "results" {
			return nil
		}
		if !strings.HasPrefix(relParts[1], "SIM-") {
			return nil
		}
		if model != "" && relParts[3] != model {
			return nil
		}
		ts := strings.TrimSuffix(relParts[4], ".md")
		if timestamp != "" && ts != timestamp {
			return nil
		}
		if isPrototypeResultPath(path) && !allowPrototype {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	return paths, err
}

func validateResultFile(repo Repo, path string, strictMatrix bool, allowPrototype bool) []string {
	var issues []string
	textBytes, err := os.ReadFile(path)
	if err != nil {
		return []string{err.Error()}
	}
	text := string(textBytes)
	if isPrototypeText(text) && !allowPrototype {
		issues = append(issues, "prototype scripted result excluded by default; pass --allow-prototype for plumbing checks")
	}
	relParts := strings.Split(repo.Rel(path), "/")
	if len(relParts) != 5 || relParts[0] != "results" {
		return append(issues, "path shape must be results/<sim>/<scenario>/<model>/<timestamp>.md")
	}
	simID, scenarioID, modelID, filename := relParts[1], relParts[2], relParts[3], relParts[4]
	timestamp := strings.TrimSuffix(filename, ".md")
	lines := strings.Split(text, "\n")
	if len(lines) == 0 || !strings.HasPrefix(lines[0], "# Result: ") {
		issues = append(issues, "missing # Result header")
	} else if !strings.Contains(lines[0], "/ "+modelID+" / "+timestamp) {
		issues = append(issues, "header does not match path model/timestamp")
	}
	for _, heading := range requiredHeadings {
		if !strings.Contains(text, heading) {
			issues = append(issues, "missing heading: "+heading)
		}
	}
	if !strings.Contains(text, "Evidence verdict:") {
		issues = append(issues, "missing Evidence verdict line")
	}
	modelLine := firstLineWithPrefix(lines, "- Model ID: ")
	if !strings.Contains(modelLine, modelID) {
		issues = append(issues, "Model ID line does not match path model")
	}
	if strictMatrix && !matrixContainsResult(repo, scenarioID, repo.Rel(path)) {
		issues = append(issues, "scenario matrix does not reference this result path")
	}
	if simID == "" || scenarioID == "" {
		issues = append(issues, "empty sim or scenario path component")
	}
	return issues
}

func firstLineWithPrefix(lines []string, prefix string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			return line
		}
	}
	return ""
}

func isPrototypeResultPath(path string) bool {
	text, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return isPrototypeText(string(text))
}

func isPrototypeText(text string) bool {
	return strings.Contains(text, "Run mode: `scripted-doc-eval-blind`") ||
		strings.Contains(text, "Runner/interface: `results/tools/run_matrix_batch.py`")
}

func matrixContainsResult(repo Repo, scenarioID string, relResult string) bool {
	text, err := repo.ReadRel(filepath.ToSlash(filepath.Join("scenarios", scenarioID, "MATRIX.md")))
	if err != nil {
		return false
	}
	return strings.Contains(text, relResult)
}

func extractVerdict(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(bytes), "\n") {
		if strings.HasPrefix(line, "Evidence verdict:") {
			return strings.TrimSuffix(strings.TrimSpace(strings.TrimPrefix(line, "Evidence verdict:")), "."), nil
		}
	}
	return "", fmt.Errorf("missing Evidence verdict line in %s", path)
}
