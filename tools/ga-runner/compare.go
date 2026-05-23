package main

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
)

type backfillComparisonRecord struct {
	SimID            string
	ScenarioID       string
	V1               FitnessResult
	V2               FitnessResult
	V1Path           string
	V2Path           string
	VocabularyStatus string
	SourceResolution string
	DeltaNormalized  float64
}

type backfillComparisonSimSummary struct {
	SimID            string
	CellCount        int
	V1Average        float64
	V2Average        float64
	DeltaAverage     float64
	V1Rank           int
	V2Rank           int
	RankDelta        int
	VocabularyStatus string
	SourceResolution string
	Family           string
}

// Intent: Turn targeted rubric-v2 backfill evidence into a durable review
// document so operators can inspect rank drift before expanding rescoring scope.
// The report is additive derived evidence under results/reports/ and does not
// rewrite any scored artifact bytes. Source: DI-zuzup
func runCompareBackfill(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("compare-backfill", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	runGroupID := fs.String("run-group-id", "", "run group ID whose completed v2 backfill cells should be compared against canonical v1 evidence")
	outputPath := fs.String("output", "", "optional report path; defaults to results/reports/<run-group-id>-comparison.md")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*runGroupID) == "" {
		return errUsage("compare-backfill: -run-group-id is required")
	}
	repo, err := openRepo(*repoRoot)
	if err != nil {
		return err
	}
	stateFile, err := statePath(repo, *runGroupID)
	if err != nil {
		return err
	}
	state, err := readGAState(stateFile)
	if err != nil {
		return err
	}
	if len(state.Cells) == 0 {
		return fmt.Errorf("compare-backfill: state has no cells")
	}
	for _, cell := range state.Cells {
		if cell.Status != "done" {
			return fmt.Errorf("compare-backfill: cell %s status=%s; all cells must be done", cell.CellID, cell.Status)
		}
	}
	records, err := auditCanonicalV1Results(repo, "", "")
	if err != nil {
		return err
	}
	index := indexAuditRecordsByPair(records)
	report, comparisons, simSummaries, unmatched, ambiguous, err := buildBackfillComparisonReport(repo, state, index)
	if err != nil {
		return err
	}
	target := strings.TrimSpace(*outputPath)
	if target == "" {
		target = filepath.ToSlash(filepath.Join("results", "reports", *runGroupID+"-comparison.md"))
	}
	if err := writeFile(repo.Abs(target), report); err != nil {
		return err
	}
	return writeFormat(stdout, "report=%s cells=%d sims=%d unmatched=%d ambiguous=%d\n",
		repo.Rel(repo.Abs(target)),
		len(comparisons),
		len(simSummaries),
		unmatched,
		ambiguous,
	)
}

type auditRecordIndex map[string][]auditRecord

// indexAuditRecordsByPair prepares exact-match historical v1 evidence for cheap
// lookup by the v2 backfill cell's sim/scenario pair.
//
// Intent: `compare-backfill` should only compare against the stable historical
// evidence set that qualified for targeted backfill, while keeping any
// remaining ambiguity visible in the report. Source: DI-zuzup
func indexAuditRecordsByPair(records []auditRecord) auditRecordIndex {
	index := auditRecordIndex{}
	for _, record := range records {
		if !record.ExactMatch {
			continue
		}
		key := backfillPairKey(record.Result.SimID, record.Result.ScenarioID)
		index[key] = append(index[key], record)
	}
	for key := range index {
		sort.Slice(index[key], func(i, j int) bool {
			left := index[key][i]
			right := index[key][j]
			if preferComparisonAuditRecord(left, right) {
				return true
			}
			if preferComparisonAuditRecord(right, left) {
				return false
			}
			return left.Path < right.Path
		})
	}
	return index
}

func backfillPairKey(simID string, scenarioID string) string {
	return simID + "\x00" + scenarioID
}

func preferComparisonAuditRecord(left auditRecord, right auditRecord) bool {
	if left.Result.Runner.APIModel != "" && right.Result.Runner.APIModel == "" {
		return true
	}
	if left.Result.Runner.APIModel == "" && right.Result.Runner.APIModel != "" {
		return false
	}
	if left.Result.TimestampUTC != right.Result.TimestampUTC {
		return left.Result.TimestampUTC > right.Result.TimestampUTC
	}
	if left.Result.ModelID != right.Result.ModelID {
		return left.Result.ModelID < right.Result.ModelID
	}
	return left.Path < right.Path
}

// buildBackfillComparisonReport pairs each completed v2 backfill cell with the
// best historical exact-match v1 candidate, then renders a Markdown report from
// the matched set.
//
// Intent: Make the operator's first drift review deterministic and auditable,
// even when multiple historical v1 records exist for one sim/scenario pair.
// Source: DI-zuzup
func buildBackfillComparisonReport(repo Repo, state GAState, index auditRecordIndex) (string, []backfillComparisonRecord, []backfillComparisonSimSummary, int, int, error) {
	var comparisons []backfillComparisonRecord
	unmatched := 0
	ambiguous := 0
	seenResultPaths := map[string]bool{}
	for _, cell := range state.Cells {
		v2Path := repo.Abs(cell.SelectedResultPath())
		v2RelPath := repo.Rel(v2Path)
		if seenResultPaths[v2RelPath] {
			continue
		}
		seenResultPaths[v2RelPath] = true
		v2, err := readFitnessResult(v2Path)
		if err != nil {
			return "", nil, nil, 0, 0, err
		}
		key := backfillPairKey(cell.SimID, cell.ScenarioID)
		candidates := append([]auditRecord(nil), index[key]...)
		if len(candidates) == 0 {
			unmatched++
			continue
		}
		candidates = narrowAuditCandidatesByAPIModel(candidates, v2.Runner.APIModel)
		if len(candidates) > 1 {
			ambiguous++
		}
		selected := candidates[0]
		comparisons = append(comparisons, backfillComparisonRecord{
			SimID:            cell.SimID,
			ScenarioID:       cell.ScenarioID,
			V1:               selected.Result,
			V2:               v2,
			V1Path:           selected.Path,
			V2Path:           v2RelPath,
			VocabularyStatus: selected.VocabularyStatus,
			SourceResolution: selected.SourceResolution,
			DeltaNormalized:  v2.Fitness.Normalized0To100 - selected.Result.Fitness.Normalized0To100,
		})
	}
	sort.Slice(comparisons, func(i, j int) bool {
		if comparisons[i].SimID != comparisons[j].SimID {
			return comparisons[i].SimID < comparisons[j].SimID
		}
		return comparisons[i].ScenarioID < comparisons[j].ScenarioID
	})
	simSummaries := summarizeBackfillComparison(comparisons)
	report := renderBackfillComparisonMarkdown(state, comparisons, simSummaries, unmatched, ambiguous)
	return report, comparisons, simSummaries, unmatched, ambiguous, nil
}

// narrowAuditCandidatesByAPIModel keeps same-provider-model comparisons
// together when the historical corpus contains more than one exact-match v1
// record for the same sim/scenario pair.
//
// Intent: Prefer comparing rubric drift before model drift whenever the source
// evidence allows it. Source: DI-zuzup
func narrowAuditCandidatesByAPIModel(candidates []auditRecord, apiModel string) []auditRecord {
	if strings.TrimSpace(apiModel) == "" {
		return candidates
	}
	var filtered []auditRecord
	for _, candidate := range candidates {
		if strings.TrimSpace(candidate.Result.Runner.APIModel) == strings.TrimSpace(apiModel) {
			filtered = append(filtered, candidate)
		}
	}
	if len(filtered) > 0 {
		return filtered
	}
	return candidates
}

// summarizeBackfillComparison collapses cell-level v1/v2 evidence into sim-level
// averages and ranks so the report can show rank drift among the rescored sims.
//
// Intent: `tapur.36` requires rank-delta reporting for envelope contenders and
// conformance-family sims, not just raw per-cell deltas. Source: DI-zuzup
func summarizeBackfillComparison(records []backfillComparisonRecord) []backfillComparisonSimSummary {
	type aggregate struct {
		simID            string
		count            int
		v1Sum            float64
		v2Sum            float64
		vocabularyStatus string
		sourceResolution string
	}
	aggregates := map[string]*aggregate{}
	for _, record := range records {
		item, ok := aggregates[record.SimID]
		if !ok {
			item = &aggregate{
				simID:            record.SimID,
				vocabularyStatus: record.VocabularyStatus,
				sourceResolution: record.SourceResolution,
			}
			aggregates[record.SimID] = item
		}
		item.count++
		item.v1Sum += record.V1.Fitness.Normalized0To100
		item.v2Sum += record.V2.Fitness.Normalized0To100
		item.vocabularyStatus = strongerVocabularyStatus(item.vocabularyStatus, record.VocabularyStatus)
		item.sourceResolution = strongerSourceResolution(item.sourceResolution, record.SourceResolution)
	}
	var summaries []backfillComparisonSimSummary
	for _, item := range aggregates {
		v1Average := item.v1Sum / float64(item.count)
		v2Average := item.v2Sum / float64(item.count)
		summaries = append(summaries, backfillComparisonSimSummary{
			SimID:            item.simID,
			CellCount:        item.count,
			V1Average:        v1Average,
			V2Average:        v2Average,
			DeltaAverage:     v2Average - v1Average,
			VocabularyStatus: item.vocabularyStatus,
			SourceResolution: item.sourceResolution,
			Family:           comparisonSimFamily(item.simID),
		})
	}
	assignComparisonRanks(summaries, false)
	assignComparisonRanks(summaries, true)
	sort.Slice(summaries, func(i, j int) bool {
		if summaries[i].V2Rank != summaries[j].V2Rank {
			return summaries[i].V2Rank < summaries[j].V2Rank
		}
		return summaries[i].SimID < summaries[j].SimID
	})
	return summaries
}

func assignComparisonRanks(summaries []backfillComparisonSimSummary, useV2 bool) {
	type ranked struct {
		index int
		score float64
	}
	var items []ranked
	for index, summary := range summaries {
		score := summary.V1Average
		if useV2 {
			score = summary.V2Average
		}
		items = append(items, ranked{index: index, score: score})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].score != items[j].score {
			return items[i].score > items[j].score
		}
		return summaries[items[i].index].SimID < summaries[items[j].index].SimID
	})
	for rank, item := range items {
		if useV2 {
			summaries[item.index].V2Rank = rank + 1
			summaries[item.index].RankDelta = summaries[item.index].V1Rank - summaries[item.index].V2Rank
			continue
		}
		summaries[item.index].V1Rank = rank + 1
	}
}

func comparisonSimFamily(simID string) string {
	lower := strings.ToLower(simID)
	switch {
	case strings.Contains(lower, "grid-envelope"):
		return "grid-envelope"
	case strings.Contains(lower, "conformance"):
		return "conformance-family"
	default:
		return "other"
	}
}

func strongerVocabularyStatus(current string, next string) string {
	order := map[string]int{"": -1, "clean": 0, "soft_hit": 1, "hard_hit": 2}
	if order[next] > order[current] {
		return next
	}
	return current
}

func strongerSourceResolution(current string, next string) string {
	if current == "" {
		return next
	}
	if current == next {
		return current
	}
	return auditSourceResolutionMixed
}

// renderBackfillComparisonMarkdown writes the operator-facing review artifact in
// a stable Markdown shape so the report can be committed or cited from TODOs
// without depending on terminal formatting.
//
// Intent: Keep the comparison artifact readable as protocol evidence rather than
// as an ephemeral CLI-only summary. Source: DI-zuzup
func renderBackfillComparisonMarkdown(state GAState, comparisons []backfillComparisonRecord, simSummaries []backfillComparisonSimSummary, unmatched int, ambiguous int) string {
	var out strings.Builder
	out.WriteString("# GA Backfill Comparison Report\n\n")
	out.WriteString(fmt.Sprintf("- Run group: `%s`\n", state.RunGroupID))
	out.WriteString(fmt.Sprintf("- State: `results/state/%s.json`\n", state.RunGroupID))
	out.WriteString("- Comparison basis: latest exact-match canonical `promisegrid.ga.result.v1` record for the same `sim_id` + `scenario_id`, preferring the same `runner.api_model` when available.\n")
	out.WriteString(fmt.Sprintf("- Compared cells: `%d`\n", len(comparisons)))
	out.WriteString(fmt.Sprintf("- Unmatched v2 cells: `%d`\n", unmatched))
	out.WriteString(fmt.Sprintf("- Ambiguous matched pairs: `%d`\n\n", ambiguous))

	out.WriteString("## Sim Rank Drift\n\n")
	out.WriteString("| V2 rank | V1 rank | Δ rank | Sim | Family | Cells | V1 avg | V2 avg | Δ avg | Vocab | Source |\n")
	out.WriteString("|---:|---:|---:|---|---|---:|---:|---:|---:|---|---|\n")
	for _, summary := range simSummaries {
		out.WriteString(fmt.Sprintf("| %d | %d | %+d | `%s` | %s | %d | %.2f | %.2f | %+0.2f | %s | %s |\n",
			summary.V2Rank,
			summary.V1Rank,
			summary.RankDelta,
			summary.SimID,
			summary.Family,
			summary.CellCount,
			summary.V1Average,
			summary.V2Average,
			summary.DeltaAverage,
			summary.VocabularyStatus,
			summary.SourceResolution,
		))
	}
	out.WriteString("\n")

	out.WriteString("## Largest Cell Deltas\n\n")
	out.WriteString("| Δ score | Sim | Scenario | V1 | V2 | V1 result | V2 result |\n")
	out.WriteString("|---:|---|---|---:|---:|---|---|\n")
	for _, record := range topAbsoluteDeltaComparisons(comparisons, 12) {
		out.WriteString(fmt.Sprintf("| %+0.2f | `%s` | `%s` | %.2f | %.2f | `%s` | `%s` |\n",
			record.DeltaNormalized,
			record.SimID,
			record.ScenarioID,
			record.V1.Fitness.Normalized0To100,
			record.V2.Fitness.Normalized0To100,
			record.V1Path,
			record.V2Path,
		))
	}
	out.WriteString("\n")

	out.WriteString("## Family Highlights\n\n")
	for _, family := range []string{"grid-envelope", "conformance-family"} {
		section := filterSimSummariesByFamily(simSummaries, family)
		if len(section) == 0 {
			continue
		}
		out.WriteString(fmt.Sprintf("### %s\n\n", family))
		out.WriteString("| Sim | V1 avg | V2 avg | Δ avg | Vocab | Source |\n")
		out.WriteString("|---|---:|---:|---:|---|---|\n")
		for _, summary := range section {
			out.WriteString(fmt.Sprintf("| `%s` | %.2f | %.2f | %+0.2f | %s | %s |\n",
				summary.SimID,
				summary.V1Average,
				summary.V2Average,
				summary.DeltaAverage,
				summary.VocabularyStatus,
				summary.SourceResolution,
			))
		}
		out.WriteString("\n")
	}

	out.WriteString("## Cell Detail\n\n")
	out.WriteString("| Sim | Scenario | V1 model | V2 model | V1 | V2 | Δ | Vocab | Source |\n")
	out.WriteString("|---|---|---|---|---:|---:|---:|---|---|\n")
	for _, record := range comparisons {
		out.WriteString(fmt.Sprintf("| `%s` | `%s` | `%s` | `%s` | %.2f | %.2f | %+0.2f | %s | %s |\n",
			record.SimID,
			record.ScenarioID,
			record.V1.ModelID,
			record.V2.ModelID,
			record.V1.Fitness.Normalized0To100,
			record.V2.Fitness.Normalized0To100,
			record.DeltaNormalized,
			record.VocabularyStatus,
			record.SourceResolution,
		))
	}
	return out.String()
}

func topAbsoluteDeltaComparisons(records []backfillComparisonRecord, limit int) []backfillComparisonRecord {
	items := append([]backfillComparisonRecord(nil), records...)
	sort.Slice(items, func(i, j int) bool {
		left := items[i].DeltaNormalized
		if left < 0 {
			left = -left
		}
		right := items[j].DeltaNormalized
		if right < 0 {
			right = -right
		}
		if left != right {
			return left > right
		}
		if items[i].SimID != items[j].SimID {
			return items[i].SimID < items[j].SimID
		}
		return items[i].ScenarioID < items[j].ScenarioID
	})
	if len(items) > limit {
		items = items[:limit]
	}
	return items
}

func filterSimSummariesByFamily(summaries []backfillComparisonSimSummary, family string) []backfillComparisonSimSummary {
	var filtered []backfillComparisonSimSummary
	for _, summary := range summaries {
		if summary.Family == family {
			filtered = append(filtered, summary)
		}
	}
	return filtered
}
