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

var scoreRules = []struct {
	Phrase string
	Score  int
}{
	{"strongest fit", 6}, {"best migration fit", 6}, {"best fit", 6},
	{"strong fit", 5}, {"good partial fit", 4}, {"good strict baseline", 4},
	{"partial fit", 3}, {"partial guardrail", 3}, {"partial with privacy risk", 2},
	{"partial but brittle", 2}, {"weak-to-partial fit", 2}, {"weak fit", 1},
	{"poor standalone fit", 0}, {"poor fit", 0}, {"negative-control", 1},
}

type CellKey struct {
	Sim      string
	Scenario string
}

type ResultCell struct {
	Key       CellKey
	Model     string
	Timestamp string
	Path      string
	Verdict   string
	Score     int
}

func runCompare(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("compare", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	oldModel := fs.String("old-model", "", "baseline model slug")
	newModel := fs.String("new-model", "", "comparison model slug")
	oldTS := fs.String("old-ts", "latest", "baseline timestamp or latest")
	newTS := fs.String("new-ts", "latest", "comparison timestamp or latest")
	output := fs.String("output", "", "output markdown path")
	includePrototype := fs.Bool("include-prototype", false, "include scripted prototype outputs")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *oldModel == "" || *newModel == "" {
		return errUsage("compare: --old-model and --new-model are required")
	}
	repo, err := openRepo(*repoRoot)
	if err != nil {
		return err
	}
	oldIndex, err := indexModelCells(repo, *oldModel, *includePrototype)
	if err != nil {
		return err
	}
	newIndex, err := indexModelCells(repo, *newModel, *includePrototype)
	if err != nil {
		return err
	}
	var keys []CellKey
	for key := range oldIndex {
		if _, ok := newIndex[key]; ok {
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Sim == keys[j].Sim {
			return keys[i].Scenario < keys[j].Scenario
		}
		return keys[i].Sim < keys[j].Sim
	})
	if len(keys) == 0 {
		return fmt.Errorf("no shared (sim, scenario) cells between the two models")
	}
	oldCells, err := loadCompareCells(repo, *oldModel, *oldTS, keys, oldIndex)
	if err != nil {
		return err
	}
	newCells, err := loadCompareCells(repo, *newModel, *newTS, keys, newIndex)
	if err != nil {
		return err
	}
	report := renderCompareReport(*oldModel, *newModel, oldCells, newCells)
	outPath := *output
	if outPath == "" {
		inferredTS := latestTimestamp(newCells)
		outPath = repo.Path("results", "comparisons", *oldModel+"_vs_"+*newModel+"_"+inferredTS+".md")
	} else {
		outPath = repo.Abs(outPath)
	}
	if err := writeFile(outPath, report); err != nil {
		return err
	}
	fmt.Fprintln(stdout, repo.Rel(outPath))
	return nil
}

func scoreVerdict(verdict string) int {
	lower := strings.ToLower(verdict)
	for _, rule := range scoreRules {
		if strings.Contains(lower, rule.Phrase) {
			return rule.Score
		}
	}
	switch {
	case strings.Contains(lower, "strong"):
		return 5
	case strings.Contains(lower, "good"):
		return 4
	case strings.Contains(lower, "partial"):
		return 3
	case strings.Contains(lower, "weak"):
		return 1
	case strings.Contains(lower, "poor"):
		return 0
	default:
		return 2
	}
}

func indexModelCells(repo Repo, model string, includePrototype bool) (map[CellKey]map[string]string, error) {
	index := map[CellKey]map[string]string{}
	root := repo.Path("results")
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" || path == repo.Path("results", "README.md") {
			return nil
		}
		parts := strings.Split(repo.Rel(path), "/")
		if len(parts) != 5 || parts[0] != "results" || parts[3] != model {
			return nil
		}
		if isPrototypeResultPath(path) && !includePrototype {
			return nil
		}
		key := CellKey{Sim: parts[1], Scenario: parts[2]}
		if index[key] == nil {
			index[key] = map[string]string{}
		}
		index[key][strings.TrimSuffix(parts[4], ".md")] = path
		return nil
	})
	return index, err
}

func chooseCellPath(byTimestamp map[string]string, requestedTS string, model string, key CellKey) (string, string, error) {
	if requestedTS == "latest" {
		var timestamps []string
		for timestamp := range byTimestamp {
			timestamps = append(timestamps, timestamp)
		}
		sort.Strings(timestamps)
		ts := timestamps[len(timestamps)-1]
		return ts, byTimestamp[ts], nil
	}
	path, ok := byTimestamp[requestedTS]
	if !ok {
		return "", "", fmt.Errorf("missing %s timestamp %s for cell (%s, %s)", model, requestedTS, key.Sim, key.Scenario)
	}
	return requestedTS, path, nil
}

func loadCompareCells(repo Repo, model string, requestedTS string, keys []CellKey, index map[CellKey]map[string]string) (map[CellKey]ResultCell, error) {
	out := map[CellKey]ResultCell{}
	for _, key := range keys {
		ts, path, err := chooseCellPath(index[key], requestedTS, model, key)
		if err != nil {
			return nil, err
		}
		verdict, err := extractVerdict(path)
		if err != nil {
			return nil, err
		}
		out[key] = ResultCell{
			Key:       key,
			Model:     model,
			Timestamp: ts,
			Path:      repo.Rel(path),
			Verdict:   verdict,
			Score:     scoreVerdict(verdict),
		}
	}
	return out, nil
}

func renderCompareReport(oldModel string, newModel string, oldCells map[CellKey]ResultCell, newCells map[CellKey]ResultCell) string {
	keys := make([]CellKey, 0, len(oldCells))
	for key := range oldCells {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Sim == keys[j].Sim {
			return keys[i].Scenario < keys[j].Scenario
		}
		return keys[i].Sim < keys[j].Sim
	})
	verdictTextChanges := 0
	scoreChanges := 0
	for _, key := range keys {
		if oldCells[key].Verdict != newCells[key].Verdict {
			verdictTextChanges++
		}
		if oldCells[key].Score != newCells[key].Score {
			scoreChanges++
		}
	}
	oldRanked := simulationRank(oldCells)
	newRanked := simulationRank(newCells)
	oldRank, oldAvg := rankMaps(oldRanked)
	newRank, newAvg := rankMaps(newRanked)
	oldScenarios := scenarioAverages(oldCells)
	newScenarios := scenarioAverages(newCells)

	var lines []string
	lines = append(lines, "# Cross-Model Drift Report: "+oldModel+" vs "+newModel, "")
	lines = append(lines, "- Baseline model: `"+oldModel+"`")
	lines = append(lines, "- Comparison model: `"+newModel+"`")
	lines = append(lines, fmt.Sprintf("- Cells compared: `%d`", len(keys)))
	lines = append(lines, fmt.Sprintf("- Verdict text changes: `%d`", verdictTextChanges))
	lines = append(lines, fmt.Sprintf("- Score/rank changes: `%d`", scoreChanges), "")
	lines = append(lines, "## Simulation Ranking Shift", "")
	lines = append(lines, "| Simulation | Old avg score | Old rank | New avg score | New rank | Rank delta |")
	lines = append(lines, "|---|---:|---:|---:|---:|---:|")
	for _, item := range oldRanked {
		sim := item.Sim
		lines = append(lines, fmt.Sprintf("| `%s` | %.2f | %d | %.2f | %d | %+d |", sim, oldAvg[sim], oldRank[sim], newAvg[sim], newRank[sim], oldRank[sim]-newRank[sim]))
	}
	lines = append(lines, "", "## Per-Cell Drift", "")
	lines = append(lines, "| Simulation | Scenario | Old verdict | New verdict | Score delta |")
	lines = append(lines, "|---|---|---|---|---:|")
	for _, key := range keys {
		lines = append(lines, fmt.Sprintf("| `%s` | `%s` | %s | %s | %+d |", key.Sim, key.Scenario, markdownEscapeCell(oldCells[key].Verdict), markdownEscapeCell(newCells[key].Verdict), newCells[key].Score-oldCells[key].Score))
	}
	lines = append(lines, "", "## Scenario Aggregates", "")
	lines = append(lines, "| Scenario | Old avg score | New avg score | Delta |")
	lines = append(lines, "|---|---:|---:|---:|")
	var scenarios []string
	for scenario := range oldScenarios {
		scenarios = append(scenarios, scenario)
	}
	sort.Strings(scenarios)
	for _, scenario := range scenarios {
		lines = append(lines, fmt.Sprintf("| `%s` | %.2f | %.2f | %+.2f |", scenario, oldScenarios[scenario], newScenarios[scenario], newScenarios[scenario]-oldScenarios[scenario]))
	}
	lines = append(lines, "", "## Notes", "")
	lines = append(lines, "- This report compares verdict lines from result artifacts; it does not execute protocol harness code.")
	lines = append(lines, "- Scripted prototype plumbing outputs are excluded unless the comparison tool is run with `--include-prototype`.")
	lines = append(lines, "- Paths and line numbers for each verdict are in the source result files under `results/<sim>/<scenario>/<model>/<timestamp>.md`.")
	return joinLines(lines)
}

type rankItem struct {
	Sim string
	Avg float64
}

func simulationRank(cells map[CellKey]ResultCell) []rankItem {
	grouped := map[string][]int{}
	for _, cell := range cells {
		grouped[cell.Key.Sim] = append(grouped[cell.Key.Sim], cell.Score)
	}
	var ranked []rankItem
	for sim, scores := range grouped {
		ranked = append(ranked, rankItem{Sim: sim, Avg: average(scores)})
	}
	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Avg == ranked[j].Avg {
			return ranked[i].Sim < ranked[j].Sim
		}
		return ranked[i].Avg > ranked[j].Avg
	})
	return ranked
}

func rankMaps(items []rankItem) (map[string]int, map[string]float64) {
	ranks := map[string]int{}
	avgs := map[string]float64{}
	for index, item := range items {
		ranks[item.Sim] = index + 1
		avgs[item.Sim] = item.Avg
	}
	return ranks, avgs
}

func scenarioAverages(cells map[CellKey]ResultCell) map[string]float64 {
	grouped := map[string][]int{}
	for _, cell := range cells {
		grouped[cell.Key.Scenario] = append(grouped[cell.Key.Scenario], cell.Score)
	}
	out := map[string]float64{}
	for scenario, scores := range grouped {
		out[scenario] = average(scores)
	}
	return out
}

func average(values []int) float64 {
	total := 0
	for _, value := range values {
		total += value
	}
	return float64(total) / float64(len(values))
}

func latestTimestamp(cells map[CellKey]ResultCell) string {
	var timestamps []string
	for _, cell := range cells {
		timestamps = append(timestamps, cell.Timestamp)
	}
	sort.Strings(timestamps)
	return timestamps[len(timestamps)-1]
}
