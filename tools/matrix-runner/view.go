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

// runView renders a read-only table from canonical result files.
//
// Intent: Replace committed scenario MATRIX.md state with an on-demand view so
// scenario trees remain input context and result trees remain output evidence.
// Source: DI-zamin
func runView(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("view", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	simFilter := fs.String("sim", "", "optional simulation ID filter")
	scenarioFilter := fs.String("scenario", "", "optional scenario ID filter")
	modelFilter := fs.String("model", "", "optional model ID filter")
	timestampFilter := fs.String("timestamp", "latest", "timestamp filter: latest, all, or exact YYYYMMDD-HHMMSS")
	output := fs.String("output", "", "optional markdown output path; default stdout")
	excludePrototype := fs.Bool("exclude-prototype", false, "hide scripted prototype results")
	if err := fs.Parse(args); err != nil {
		return err
	}
	repo, err := openRepo(*repoRoot)
	if err != nil {
		return err
	}
	rows, err := collectResultViewRows(repo, *simFilter, *scenarioFilter, *modelFilter, *timestampFilter, *excludePrototype)
	if err != nil {
		return err
	}
	text := renderResultView(*simFilter, *scenarioFilter, *modelFilter, *timestampFilter, rows)
	if *output == "" {
		return writeText(stdout, text)
	}
	if err := writeFile(repo.Abs(*output), text); err != nil {
		return err
	}
	return writeLine(stdout, repo.Rel(repo.Abs(*output)))
}

// ResultViewRow is one rendered result artifact in the generated inspection
// table. It deliberately mirrors the result path coordinates instead of keeping
// a second committed source of truth.
type ResultViewRow struct {
	Simulation string
	Scenario   string
	Model      string
	Timestamp  string
	ResultPath string
	Verdict    string
}

// collectResultViewRows scans result files and applies the requested filters.
// The default "latest" mode keeps one newest timestamp per
// simulation/scenario/model coordinate so the view stays compact.
func collectResultViewRows(repo Repo, simFilter string, scenarioFilter string, modelFilter string, timestampFilter string, excludePrototype bool) ([]ResultViewRow, error) {
	if timestampFilter == "" {
		timestampFilter = "latest"
	}
	if timestampFilter != "latest" && timestampFilter != "all" && strings.Contains(timestampFilter, timestampPlaceholder) {
		return nil, fmt.Errorf("view: timestamp filter must be latest, all, or a concrete timestamp")
	}
	var allRows []ResultViewRow
	err := filepath.WalkDir(repo.Path("results"), func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		simID, scenarioID, modelID, timestamp, abs, err := resultCoordinates(repo, path)
		if err != nil {
			return nil
		}
		if simFilter != "" && simID != simFilter {
			return nil
		}
		if scenarioFilter != "" && scenarioID != scenarioFilter {
			return nil
		}
		if modelFilter != "" && modelID != modelFilter {
			return nil
		}
		if timestampFilter != "latest" && timestampFilter != "all" && timestamp != timestampFilter {
			return nil
		}
		if excludePrototype && isPrototypeResultPath(abs) {
			return nil
		}
		verdict, err := extractVerdict(abs)
		if err != nil {
			return err
		}
		allRows = append(allRows, ResultViewRow{
			Simulation: simID,
			Scenario:   scenarioID,
			Model:      modelID,
			Timestamp:  timestamp,
			ResultPath: repo.Rel(abs),
			Verdict:    verdict,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if timestampFilter == "latest" {
		allRows = latestResultRows(allRows)
	}
	sort.Slice(allRows, func(i, j int) bool {
		left := allRows[i]
		right := allRows[j]
		return viewSortKey(left) < viewSortKey(right)
	})
	if len(allRows) == 0 {
		return nil, fmt.Errorf("view: no result files matched selection")
	}
	return allRows, nil
}

// resultCoordinates validates and returns the locked coordinates encoded in a
// result path. That path shape is the canonical linkage among result, scenario,
// simulation, model, and timestamp after scenario matrices were retired.
func resultCoordinates(repo Repo, result string) (simID, scenarioID, modelID, timestamp, abs string, err error) {
	abs = repo.Abs(result)
	rel := repo.Rel(abs)
	parts := strings.Split(rel, "/")
	if len(parts) != 5 || parts[0] != "results" || !strings.HasSuffix(parts[4], ".md") {
		err = fmt.Errorf("result path must have shape results/<sim-id>/<scenario-id>/<model-id>/<timestamp>.md")
		return
	}
	return parts[1], parts[2], parts[3], strings.TrimSuffix(parts[4], ".md"), abs, nil
}

// latestResultRows keeps the newest timestamp per result coordinate. The
// timestamp format is lexicographically sortable because it is YYYYMMDD-HHMMSS.
func latestResultRows(rows []ResultViewRow) []ResultViewRow {
	latest := map[string]ResultViewRow{}
	for _, row := range rows {
		key := row.Simulation + "\x00" + row.Scenario + "\x00" + row.Model
		if previous, ok := latest[key]; !ok || row.Timestamp > previous.Timestamp {
			latest[key] = row
		}
	}
	out := make([]ResultViewRow, 0, len(latest))
	for _, row := range latest {
		out = append(out, row)
	}
	return out
}

// viewSortKey gives generated tables stable ordering for review and tests.
func viewSortKey(row ResultViewRow) string {
	return row.Simulation + "\x00" + row.Scenario + "\x00" + row.Model + "\x00" + row.Timestamp
}

// renderResultView builds the Markdown table without writing it. Keeping render
// separate from I/O lets tests and future callers generate alternate views
// without reintroducing committed scenario-side result state.
func renderResultView(simFilter string, scenarioFilter string, modelFilter string, timestampFilter string, rows []ResultViewRow) string {
	var out strings.Builder
	fmt.Fprintln(&out, "# Result View")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "Generated from canonical result files under `results/`. Source: DI-zamin")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "## Filters")
	fmt.Fprintln(&out)
	fmt.Fprintf(&out, "- Simulation: `%s`\n", displayFilter(simFilter, "all"))
	fmt.Fprintf(&out, "- Scenario: `%s`\n", displayFilter(scenarioFilter, "all"))
	fmt.Fprintf(&out, "- Model: `%s`\n", displayFilter(modelFilter, "all"))
	fmt.Fprintf(&out, "- Timestamp: `%s`\n", displayFilter(timestampFilter, "latest"))
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "## Results")
	fmt.Fprintln(&out)
	fmt.Fprintln(&out, "| Simulation | Scenario | Model | Timestamp | Result | Verdict |")
	fmt.Fprintln(&out, "|---|---|---|---|---|---|")
	for _, row := range rows {
		fmt.Fprintf(
			&out,
			"| `%s` | `%s` | `%s` | `%s` | `%s` | %s |\n",
			markdownEscapeCell(row.Simulation),
			markdownEscapeCell(row.Scenario),
			markdownEscapeCell(row.Model),
			markdownEscapeCell(row.Timestamp),
			markdownEscapeCell(row.ResultPath),
			markdownEscapeCell(row.Verdict),
		)
	}
	return out.String()
}

// displayFilter normalizes empty flag values in the generated filter summary.
func displayFilter(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
