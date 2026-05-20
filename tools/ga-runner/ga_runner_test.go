package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestValidateResultFileAcceptsValidJSON(t *testing.T) {
	repo := newTestRepo(t)
	path := repo.Path("results", "SIM-alpha", "scenario-one", "openai-gpt-5.5-xhigh", "20260519-101500.json")
	result := validResult(repo, path)
	if err := writeFitnessResultAtomic(path, result); err != nil {
		t.Fatalf("write valid result: %v", err)
	}
	issues := validateResultFile(repo, path)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestFindResultFilesIgnoresMarkdownCanaries(t *testing.T) {
	repo := newTestRepo(t)
	jsonPath := repo.Path("results", "SIM-alpha", "scenario-one", "openai-gpt-5.5-xhigh", "20260519-101500.json")
	result := validResult(repo, jsonPath)
	if err := writeFitnessResultAtomic(jsonPath, result); err != nil {
		t.Fatalf("write JSON result: %v", err)
	}
	mdPath := repo.Path("results", "SIM-alpha", "scenario-one", "openai-gpt-5.5-xhigh", "20260519-101500.md")
	if err := ensureParent(mdPath); err != nil {
		t.Fatalf("make md parent: %v", err)
	}
	if err := os.WriteFile(mdPath, []byte("# old canary\n"), 0o644); err != nil {
		t.Fatalf("write md canary: %v", err)
	}
	paths, err := findResultFiles(repo, "", "")
	if err != nil {
		t.Fatalf("find result files: %v", err)
	}
	if len(paths) != 1 || paths[0] != jsonPath {
		t.Fatalf("expected only JSON path %q, got %v", jsonPath, paths)
	}
}

func TestFindResultFilesHonorsModelAndTimestampFilters(t *testing.T) {
	repo := newTestRepo(t)
	first := repo.Path("results", "SIM-alpha", "scenario-one", "model-a", "20260519-101500.json")
	second := repo.Path("results", "SIM-alpha", "scenario-one", "model-b", "20260519-101501.json")
	for _, path := range []string{first, second} {
		result := validResult(repo, path)
		if err := writeFitnessResultAtomic(path, result); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}
	paths, err := findResultFiles(repo, "model-b", "20260519-101501")
	if err != nil {
		t.Fatalf("find filtered result files: %v", err)
	}
	if len(paths) != 1 || paths[0] != second {
		t.Fatalf("expected only %q, got %v", second, paths)
	}
}

func TestValidateResultFileRejectsPathShape(t *testing.T) {
	repo := newTestRepo(t)
	path := repo.Path("results", "SIM-alpha", "scenario-one", "model-a", "20260519-101500.md")
	result := validResult(repo, strings.TrimSuffix(path, ".md")+".json")
	if err := ensureParent(path); err != nil {
		t.Fatalf("make parent: %v", err)
	}
	if err := os.WriteFile(path, mustJSON(t, result), 0o644); err != nil {
		t.Fatalf("write wrong extension result: %v", err)
	}
	issues := validateResultFile(repo, path)
	if !hasIssue(issues, "path extension must be .json") {
		t.Fatalf("expected extension issue, got %v", issues)
	}
}

func TestValidateResultFileRejectsSchemaMismatch(t *testing.T) {
	repo := newTestRepo(t)
	path := repo.Path("results", "SIM-alpha", "scenario-one", "model-a", "20260519-101500.json")
	result := validResult(repo, path)
	result.Schema = "wrong.schema"
	if err := writeFitnessResultAtomic(path, result); err != nil {
		t.Fatalf("write bad schema result: %v", err)
	}
	issues := validateResultFile(repo, path)
	if !hasIssue(issues, "schema must be "+resultSchemaV1) {
		t.Fatalf("expected schema issue, got %v", issues)
	}
}

func TestValidateResultFileRejectsScoreRange(t *testing.T) {
	repo := newTestRepo(t)
	path := repo.Path("results", "SIM-alpha", "scenario-one", "model-a", "20260519-101500.json")
	result := validResult(repo, path)
	result.Scores.ScenarioFit = 6
	if err := writeFitnessResultAtomic(path, result); err != nil {
		t.Fatalf("write bad score result: %v", err)
	}
	issues := validateResultFile(repo, path)
	if !hasIssue(issues, "scores.scenario_fit must be between 0 and 5") {
		t.Fatalf("expected score issue, got %v", issues)
	}
}

func TestWriteFitnessResultAtomicRoundTrip(t *testing.T) {
	repo := newTestRepo(t)
	path := repo.Path("results", "SIM-alpha", "scenario-one", "model-a", "20260519-101500.json")
	result := validResult(repo, path)
	if err := writeFitnessResultAtomic(path, result); err != nil {
		t.Fatalf("atomic write: %v", err)
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read written result: %v", err)
	}
	if !strings.Contains(string(bytes), resultSchemaV1) {
		t.Fatalf("written JSON did not contain schema: %s", string(bytes))
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Fatalf("temporary file should not remain; stat err=%v", err)
	}
}

func TestDiscoverTrackedPopulationExcludesUntrackedChildren(t *testing.T) {
	repo := newGitTestRepo(t)
	writeTestFile(t, repo.Path("simulations", "SIM-parent", "README.md"), "# Parent\n")
	writeTestFile(t, repo.Path("simulations", "SIM-parent", "QUESTION.md"), "# Question\n")
	writeTestFile(t, repo.Path("simulations", "SIM-child-untracked", "README.md"), "# Child\n")
	writeTestFile(t, repo.Path("simulations", "not-a-sim", "README.md"), "# Not a sim\n")
	gitAdd(t, repo, "simulations/SIM-parent/README.md", "simulations/SIM-parent/QUESTION.md", "simulations/not-a-sim/README.md")

	population, err := discoverTrackedPopulation(repo)
	if err != nil {
		t.Fatalf("discover tracked population: %v", err)
	}
	if len(population) != 1 {
		t.Fatalf("expected one tracked SIM population member, got %#v", population)
	}
	if population[0].SimID != "SIM-parent" {
		t.Fatalf("expected SIM-parent, got %s", population[0].SimID)
	}
	if len(population[0].Files) != 2 {
		t.Fatalf("expected two tracked files, got %#v", population[0].Files)
	}
	if population[0].TreeHash == "" {
		t.Fatalf("expected tree hash")
	}
}

func TestRunInitDryRunPrintsTrackedPopulation(t *testing.T) {
	repo := newGitTestRepo(t)
	writeTestFile(t, repo.Path("simulations", "SIM-parent", "README.md"), "# Parent\n")
	gitAdd(t, repo, "simulations/SIM-parent/README.md")

	var out strings.Builder
	err := runMain([]string{"ga-runner", "init", "-repo-root", repo.Root, "-dry-run"}, &out, &out)
	if err != nil {
		t.Fatalf("init dry-run: %v", err)
	}
	text := out.String()
	if !strings.Contains(text, "population=1") || !strings.Contains(text, "SIM-parent") {
		t.Fatalf("unexpected dry-run output: %s", text)
	}
}

func TestDiscoverScenariosFindsRootScenarioFiles(t *testing.T) {
	repo := newTestRepo(t)
	writeTestFile(t, repo.Path("scenarios", "scenario-b", "scenario-b.md"), "# B\n")
	writeTestFile(t, repo.Path("scenarios", "scenario-a", "scenario-a.md"), "# A\n")
	writeTestFile(t, repo.Path("scenarios", "scenario-c", "README.md"), "# ignored\n")

	scenarios, err := discoverScenarios(repo)
	if err != nil {
		t.Fatalf("discover scenarios: %v", err)
	}
	if got := scenarioIDs(scenarios); strings.Join(got, ",") != "scenario-a,scenario-b" {
		t.Fatalf("unexpected scenarios: %#v", got)
	}
}

func TestBuildGenerationPlanIsDeterministic(t *testing.T) {
	population := []PopulationSim{
		{SimID: "SIM-a", Path: "simulations/SIM-a/", TreeHash: "a"},
		{SimID: "SIM-b", Path: "simulations/SIM-b/", TreeHash: "b"},
		{SimID: "SIM-c", Path: "simulations/SIM-c/", TreeHash: "c"},
	}
	scenarios := []Scenario{
		{ScenarioID: "scenario-a", Path: "scenarios/scenario-a/scenario-a.md"},
		{ScenarioID: "scenario-b", Path: "scenarios/scenario-b/scenario-b.md"},
		{ScenarioID: "scenario-c", Path: "scenarios/scenario-c/scenario-c.md"},
	}
	options := PlanOptions{
		RunGroupID:    "ga-test",
		ModelID:       "model-a",
		ShuffleSeed:   "42",
		ParentCount:   2,
		ScenarioCount: 2,
		ChildCount:    3,
		MaxPromotions: 1,
	}
	first, err := buildGenerationPlan(population, scenarios, options)
	if err != nil {
		t.Fatalf("build first plan: %v", err)
	}
	second, err := buildGenerationPlan(population, scenarios, options)
	if err != nil {
		t.Fatalf("build second plan: %v", err)
	}
	if strings.Join(parentIDs(first.Parents), ",") != strings.Join(parentIDs(second.Parents), ",") {
		t.Fatalf("parent selection should be deterministic: %#v vs %#v", first.Parents, second.Parents)
	}
	if strings.Join(scenarioIDs(first.Scenarios), ",") != strings.Join(scenarioIDs(second.Scenarios), ",") {
		t.Fatalf("scenario selection should be deterministic: %#v vs %#v", first.Scenarios, second.Scenarios)
	}
	if len(first.ParentScoreCells) != 4 || len(first.ChildScoreCells) != 6 {
		t.Fatalf("unexpected cell counts: parent=%d child=%d", len(first.ParentScoreCells), len(first.ChildScoreCells))
	}
	for _, child := range first.Children {
		if child.Operation != childOperationBreed || len(child.ParentIDs) != 2 || len(distinctStrings(child.ParentIDs)) != 2 {
			t.Fatalf("unexpected breed child plan: %#v", child)
		}
	}
}

func TestBuildGenerationPlanValidatesCounts(t *testing.T) {
	population := []PopulationSim{{SimID: "SIM-a", TreeHash: "a"}}
	scenarios := []Scenario{{ScenarioID: "scenario-a"}}
	options := PlanOptions{
		ModelID:       "model-a",
		ParentCount:   2,
		ScenarioCount: 1,
		ChildCount:    1,
		MaxPromotions: 2,
	}
	_, err := buildGenerationPlan(population, scenarios, options)
	if err == nil || !strings.Contains(err.Error(), "max-promotions cannot exceed child-count") {
		t.Fatalf("expected max-promotions validation error, got %v", err)
	}
	options.MaxPromotions = 0
	options.ParentCount = 1
	_, err = buildGenerationPlan(population, scenarios, options)
	if err == nil || !strings.Contains(err.Error(), "parent-count must be at least 2") {
		t.Fatalf("expected parent-count validation error, got %v", err)
	}
}

func TestRunInitDryRunPrintsGenerationPlan(t *testing.T) {
	repo := newGitTestRepo(t)
	writeTestFile(t, repo.Path("simulations", "SIM-parent", "README.md"), "# Parent\n")
	writeTestFile(t, repo.Path("simulations", "SIM-second", "README.md"), "# Second\n")
	writeTestFile(t, repo.Path("simulations", "SIM-child-untracked", "README.md"), "# Child\n")
	writeTestFile(t, repo.Path("scenarios", "scenario-one", "scenario-one.md"), "# One\n")
	writeTestFile(t, repo.Path("scenarios", "scenario-two", "scenario-two.md"), "# Two\n")
	gitAdd(t, repo, "simulations/SIM-parent/README.md", "simulations/SIM-second/README.md", "scenarios/scenario-one/scenario-one.md", "scenarios/scenario-two/scenario-two.md")

	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "init",
		"-repo-root", repo.Root,
		"-dry-run",
		"-model", "model-a",
		"-run-group-id", "ga-test",
		"-shuffle-seed", "7",
		"-parent-count", "2",
		"-scenario-count", "2",
		"-child-count", "2",
		"-max-promotions", "1",
	}, &out, &out)
	if err != nil {
		t.Fatalf("init dry-run plan: %v", err)
	}
	text := out.String()
	for _, want := range []string{"population=2", "plan run_group_id=ga-test model=model-a", "parent_score_cells=4", "child_score_cells=4", "planned-child-0001", "operation=breed"} {
		if !strings.Contains(text, want) {
			t.Fatalf("expected %q in output:\n%s", want, text)
		}
	}
	if strings.Contains(text, "SIM-child-untracked") {
		t.Fatalf("untracked child should not appear in plan output:\n%s", text)
	}
}

func TestRunInitWritesGAState(t *testing.T) {
	repo := newGAFixtureRepo(t)

	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "init",
		"-repo-root", repo.Root,
		"-model", "model-a",
		"-run-group-id", "ga-state",
		"-timestamp", "20260519-111500",
		"-parent-count", "2",
		"-scenario-count", "1",
		"-child-count", "1",
		"-max-promotions", "1",
	}, &out, &out)
	if err != nil {
		t.Fatalf("init state: %v\n%s", err, out.String())
	}

	state := mustReadGAState(t, repo, "ga-state")
	if state.Schema != stateSchemaV1 || state.RunGroupID != "ga-state" {
		t.Fatalf("unexpected state identity: %#v", state)
	}
	if len(state.Parents) != 2 || len(state.ScenarioSample) != 1 || len(state.Children) != 1 || len(state.Cells) != 3 {
		t.Fatalf("unexpected state counts: parents=%d scenarios=%d children=%d cells=%d", len(state.Parents), len(state.ScenarioSample), len(state.Children), len(state.Cells))
	}
	if !strings.HasPrefix(state.Children[0].ID(), "SIM-") || state.Children[0].Status != "queued" || state.Children[0].Operation != childOperationBreed || len(state.Children[0].ParentIDs) != 2 {
		t.Fatalf("unexpected planned child: %#v", state.Children[0])
	}
}

func TestRunScoreWritesValidatedFitnessResult(t *testing.T) {
	repo := newGAFixtureRepo(t)
	initGAStateForTest(t, repo, "ga-score")
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			if request.ServiceTier != defaultServiceTier {
				t.Fatalf("score request service tier = %q, want %q", request.ServiceTier, defaultServiceTier)
			}
			return ProviderResponse{
				Text:        validScorePayloadJSON(),
				RequestID:   "req-score",
				ResponseID:  "resp-score",
				ServiceTier: defaultServiceTier,
				UsageJSON:   `{"input_tokens":1000,"input_tokens_details":{"cached_tokens":100},"output_tokens":500}`,
			}, nil
		},
	}
	var out strings.Builder
	err := runScoreWithProvider(context.Background(), repo, provider, scoreOptions{
		RunGroupID:       "ga-score",
		Target:           "parents",
		ProviderName:     "fake",
		APIModel:         "model-a",
		ReasoningEffort:  "xhigh",
		MaxOutputTokens:  4000,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
	}, &out)
	if err != nil {
		t.Fatalf("score: %v\n%s", err, out.String())
	}
	state := mustReadGAState(t, repo, "ga-score")
	if state.Cells[0].Status != "done" || state.Cells[0].RequestID != "req-score" || state.Cells[0].CostUSD <= 0 {
		t.Fatalf("score state not updated: %#v", state.Cells[0])
	}
	if state.Cells[0].ServiceTier != defaultServiceTier || state.Cells[0].ServedServiceTier != defaultServiceTier {
		t.Fatalf("score service tier not recorded: %#v", state.Cells[0])
	}
	resultPath := repo.Abs(state.Cells[0].ResultPath)
	if issues := validateResultFile(repo, resultPath); len(issues) != 0 {
		t.Fatalf("result did not validate: %v", issues)
	}
	result := mustReadFitnessResult(t, resultPath)
	if result.Runner.ServiceTier != defaultServiceTier || result.Runner.ServedServiceTier != defaultServiceTier {
		t.Fatalf("result service tier not recorded: %#v", result.Runner)
	}
}

func TestRunScoreSkipsFailedCellWhenRequested(t *testing.T) {
	repo := newGAFixtureRepo(t)
	initGAStateForTest(t, repo, "ga-score-skip")
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			return ProviderResponse{}, fmt.Errorf("provider returned empty usable output")
		},
	}
	var out strings.Builder
	err := runScoreWithProvider(context.Background(), repo, provider, scoreOptions{
		RunGroupID:       "ga-score-skip",
		Target:           "parents",
		ProviderName:     "fake",
		APIModel:         "model-a",
		ReasoningEffort:  "xhigh",
		MaxOutputTokens:  4000,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
		SkipFailedCells:  true,
	}, &out)
	if err != nil {
		t.Fatalf("score with skip failed cells: %v\n%s", err, out.String())
	}
	state := mustReadGAState(t, repo, "ga-score-skip")
	if state.Cells[0].Status != "skipped" || !strings.Contains(state.Cells[0].ValidationMessage, "provider returned empty usable output") {
		t.Fatalf("score failure was not preserved as skipped: %#v", state.Cells[0])
	}
}

func TestRunScoreUsesConfiguredWorkers(t *testing.T) {
	repo := newGAFixtureRepo(t)
	addFixtureScenarios(t, repo, "scenario-two", "scenario-three")
	initGAStateForTestWithScenarioCount(t, repo, "ga-score-workers", 3)
	var mu sync.Mutex
	inFlight := 0
	maxInFlight := 0
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			mu.Lock()
			inFlight++
			if inFlight > maxInFlight {
				maxInFlight = inFlight
			}
			mu.Unlock()
			time.Sleep(25 * time.Millisecond)
			mu.Lock()
			inFlight--
			mu.Unlock()
			return ProviderResponse{
				Text:        validScorePayloadJSON(),
				RequestID:   "req-score",
				ResponseID:  "resp-score",
				ServiceTier: defaultServiceTier,
				UsageJSON:   `{"input_tokens":1000,"input_tokens_details":{"cached_tokens":100},"output_tokens":500}`,
			}, nil
		},
	}
	var out strings.Builder
	err := runScoreWithProvider(context.Background(), repo, provider, scoreOptions{
		RunGroupID:       "ga-score-workers",
		Target:           "parents",
		ProviderName:     "fake",
		APIModel:         "model-a",
		ReasoningEffort:  "xhigh",
		Workers:          3,
		MaxOutputTokens:  4000,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
	}, &out)
	if err != nil {
		t.Fatalf("score with workers: %v\n%s", err, out.String())
	}
	if maxInFlight < 2 {
		t.Fatalf("expected concurrent provider calls, max in flight was %d", maxInFlight)
	}
	state := mustReadGAState(t, repo, "ga-score-workers")
	doneParents := 0
	for _, cell := range state.Cells {
		if cell.SimID == "SIM-parent" && cell.Status == "done" {
			doneParents++
		}
	}
	if doneParents != 3 {
		t.Fatalf("expected 3 done parent cells, got %d", doneParents)
	}
}

func TestRunScoreMarksOnlyWorkerOwnedCellsRunning(t *testing.T) {
	repo := newGAFixtureRepo(t)
	addFixtureScenarios(t, repo, "scenario-two", "scenario-three")
	initGAStateForTestWithScenarioCount(t, repo, "ga-score-worker-owned-running", 3)
	calls := 0
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			running, err := countRunningParentCells(repo, "ga-score-worker-owned-running")
			if err != nil {
				return ProviderResponse{}, err
			}
			// Intent: Catch the canary failure mode where all prepared jobs were
			// pre-marked running even though only one serialized worker owned a
			// provider call. Source: DI-juzus
			if running != 1 {
				return ProviderResponse{}, fmt.Errorf("expected exactly one running parent score cell, got %d", running)
			}
			calls++
			return ProviderResponse{
				Text:        validScorePayloadJSON(),
				RequestID:   "req-score",
				ResponseID:  "resp-score",
				ServiceTier: defaultServiceTier,
				UsageJSON:   `{"input_tokens":1000,"input_tokens_details":{"cached_tokens":100},"output_tokens":500}`,
			}, nil
		},
	}
	var out strings.Builder
	err := runScoreWithProvider(context.Background(), repo, provider, scoreOptions{
		RunGroupID:       "ga-score-worker-owned-running",
		Target:           "parents",
		ProviderName:     "fake",
		APIModel:         "model-a",
		ReasoningEffort:  "xhigh",
		Workers:          1,
		MaxOutputTokens:  4000,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
	}, &out)
	if err != nil {
		t.Fatalf("score with worker-owned running state: %v\n%s", err, out.String())
	}
	if calls != 6 {
		t.Fatalf("expected 6 serialized score calls, got %d", calls)
	}
}

func TestRunScoreResetsStaleRunningCellsBeforeDispatch(t *testing.T) {
	repo := newGAFixtureRepo(t)
	addFixtureScenarios(t, repo, "scenario-two", "scenario-three")
	initGAStateForTestWithScenarioCount(t, repo, "ga-score-stale-running", 3)
	state := mustReadGAState(t, repo, "ga-score-stale-running")
	for index := range state.Cells {
		if state.Cells[index].SimID == "SIM-parent" {
			markGACell(&state.Cells[index], "running", "stale interrupted score attempt")
		}
	}
	writeTestState(t, repo, "ga-score-stale-running", state)
	calls := 0
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			running, err := countRunningParentCells(repo, "ga-score-stale-running")
			if err != nil {
				return ProviderResponse{}, err
			}
			// Intent: A restarted score command must reclaim stale running cells
			// before dispatch; otherwise the progress file keeps reporting dead
			// workers from the interrupted process. Source: DI-juzus
			if running != 1 {
				return ProviderResponse{}, fmt.Errorf("expected one reclaimed running parent score cell, got %d", running)
			}
			calls++
			return ProviderResponse{
				Text:        validScorePayloadJSON(),
				RequestID:   "req-score",
				ResponseID:  "resp-score",
				ServiceTier: defaultServiceTier,
				UsageJSON:   `{"input_tokens":1000,"input_tokens_details":{"cached_tokens":100},"output_tokens":500}`,
			}, nil
		},
	}
	var out strings.Builder
	err := runScoreWithProvider(context.Background(), repo, provider, scoreOptions{
		RunGroupID:       "ga-score-stale-running",
		Target:           "parents",
		ProviderName:     "fake",
		APIModel:         "model-a",
		ReasoningEffort:  "xhigh",
		Workers:          1,
		MaxOutputTokens:  4000,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
	}, &out)
	if err != nil {
		t.Fatalf("score after stale running reset: %v\n%s", err, out.String())
	}
	if calls != 6 {
		t.Fatalf("expected 6 reclaimed score calls, got %d", calls)
	}
}

func TestRunScoreReservesBudgetBeforeConcurrentDispatch(t *testing.T) {
	repo := newGAFixtureRepo(t)
	addFixtureScenarios(t, repo, "scenario-two", "scenario-three")
	initGAStateForTestWithScenarioCount(t, repo, "ga-score-budget", 3)
	state := mustReadGAState(t, repo, "ga-score-budget")
	scenario, err := scenarioFromState(state, state.Cells[0].ScenarioID)
	if err != nil {
		t.Fatalf("scenario from state: %v", err)
	}
	prompt, err := buildScorePrompt(repo, state, state.Cells[0], scenario)
	if err != nil {
		t.Fatalf("build prompt: %v", err)
	}
	cost := CostConfig{
		InputUSDPerMTok:       defaultInputUSDPerMTok,
		CachedInputUSDPerMTok: defaultCachedInputUSDPerMTok,
		OutputUSDPerMTok:      defaultOutputUSDPerMTok,
	}
	estimate := cost.EstimatePromptCost(prompt, 4000)
	calls := 0
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			calls++
			return ProviderResponse{
				Text:        validScorePayloadJSON(),
				RequestID:   "req-score",
				ResponseID:  "resp-score",
				ServiceTier: defaultServiceTier,
				UsageJSON:   `{"input_tokens":1000,"input_tokens_details":{"cached_tokens":100},"output_tokens":500}`,
			}, nil
		},
	}
	var out strings.Builder
	err = runScoreWithProvider(context.Background(), repo, provider, scoreOptions{
		RunGroupID:       "ga-score-budget",
		Target:           "parents",
		ProviderName:     "fake",
		APIModel:         "model-a",
		ReasoningEffort:  "xhigh",
		Workers:          3,
		MaxOutputTokens:  4000,
		MaxRunCostUSD:    estimate.CostUSD * 1.5,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
	}, &out)
	if err != nil {
		t.Fatalf("score with budget reservation: %v\n%s", err, out.String())
	}
	if calls != 1 {
		t.Fatalf("expected budget reservation to dispatch 1 provider call, got %d", calls)
	}
	if !strings.Contains(out.String(), "budget-stop") {
		t.Fatalf("expected budget-stop output, got %q", out.String())
	}
}

func TestRunGenerateWritesChildTree(t *testing.T) {
	repo := newGAFixtureRepo(t)
	initGAStateForTest(t, repo, "ga-generate")
	state := mustReadGAState(t, repo, "ga-generate")
	childID := state.Children[0].ID()
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			if request.ServiceTier != defaultServiceTier {
				t.Fatalf("generate request service tier = %q, want %q", request.ServiceTier, defaultServiceTier)
			}
			return ProviderResponse{
				Text:        validChildBundleJSON(childID),
				RequestID:   "req-generate",
				ResponseID:  "resp-generate",
				ServiceTier: defaultServiceTier,
				UsageJSON:   `{"input_tokens":1200,"input_tokens_details":{"cached_tokens":200},"output_tokens":700}`,
			}, nil
		},
	}
	var out strings.Builder
	err := runGenerateWithProvider(context.Background(), repo, provider, generateOptions{
		RunGroupID:       "ga-generate",
		ProviderName:     "fake",
		APIModel:         "model-a",
		ReasoningEffort:  "xhigh",
		MaxOutputTokens:  6000,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
	}, &out)
	if err != nil {
		t.Fatalf("generate: %v\n%s", err, out.String())
	}
	state = mustReadGAState(t, repo, "ga-generate")
	if state.Children[0].Status != "generated" || state.Children[0].RequestID != "req-generate" || state.Children[0].TreeHash == "" {
		t.Fatalf("generate state not updated: %#v", state.Children[0])
	}
	if state.Children[0].ServiceTier != defaultServiceTier || state.Children[0].ServedServiceTier != defaultServiceTier {
		t.Fatalf("generate service tier not recorded: %#v", state.Children[0])
	}
	assertExists(t, repo.Path("simulations", childID, "README.md"))
	assertExists(t, repo.Path("simulations", childID, "QUESTION.md"))
}

func TestRunGenerateSelectsParentsByFitnessEvidence(t *testing.T) {
	repo := newGAFixtureRepo(t)
	addTrackedFixtureSim(t, repo, "SIM-high")
	addTrackedFixtureSim(t, repo, "SIM-low")
	writeTestFile(t, repo.Path("scenarios", "scenario-one", "extra-pressure.md"), "# Extra Pressure\n\nCarol needs delayed shipment audits.\n")
	gitAdd(t, repo, "scenarios/scenario-one/extra-pressure.md")
	initGAStateForTestWithParentsAndChildren(t, repo, "ga-generate-ranked", 3, 2)
	state := mustReadGAState(t, repo, "ga-generate-ranked")
	writeParentFitnessResults(t, repo, "ga-generate-ranked", map[string]float64{
		"SIM-high":   95,
		"SIM-parent": 60,
		"SIM-low":    20,
	})
	state = mustReadGAState(t, repo, "ga-generate-ranked")
	calls := 0
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			if !strings.Contains(request.Prompt, "expected to score higher than its parent set") {
				return ProviderResponse{}, fmt.Errorf("generate prompt does not state score-improvement goal")
			}
			if !strings.Contains(request.Prompt, "breed a child simulation from exactly two parent simulations") || !strings.Contains(request.Prompt, "- Operation: `breed`") {
				return ProviderResponse{}, fmt.Errorf("generate prompt missing breed operator language")
			}
			if !strings.Contains(request.Prompt, "## Compact Fitness Evidence From This Run") || !strings.Contains(request.Prompt, "normalized_0_100=95.00") {
				return ProviderResponse{}, fmt.Errorf("generate prompt missing high parent fitness evidence")
			}
			if strings.Contains(request.Prompt, "\"schema\":") || strings.Contains(request.Prompt, "\"runner\":") {
				return ProviderResponse{}, fmt.Errorf("generate prompt should not embed complete fitness JSON")
			}
			if strings.Contains(request.Prompt, "# Run Protocol") || strings.Contains(request.Prompt, "# Scenarios") {
				return ProviderResponse{}, fmt.Errorf("generate prompt should not repeat root boilerplate")
			}
			if !strings.Contains(request.Prompt, "## Scenario Pressure") || !strings.Contains(request.Prompt, "Alice asks Bob to ship labels") || !strings.Contains(request.Prompt, "Carol needs delayed shipment audits") {
				return ProviderResponse{}, fmt.Errorf("generate prompt missing scenario pressure")
			}
			if !strings.Contains(request.Prompt, "## Parent Simulation Documents") || !strings.Contains(request.Prompt, "A fixture simulation for ranked parent selection") {
				return ProviderResponse{}, fmt.Errorf("generate prompt missing parent simulation documents")
			}
			childID := promptChildID(t, request.Prompt)
			calls++
			return ProviderResponse{
				Text:        validChildBundleJSON(childID),
				RequestID:   "req-ranked",
				ResponseID:  "resp-ranked",
				ServiceTier: defaultServiceTier,
				UsageJSON:   `{"input_tokens":1200,"input_tokens_details":{"cached_tokens":200},"output_tokens":700}`,
			}, nil
		},
	}
	var out strings.Builder
	err := runGenerateWithProvider(context.Background(), repo, provider, generateOptions{
		RunGroupID:       "ga-generate-ranked",
		ProviderName:     "fake",
		APIModel:         "model-a",
		ReasoningEffort:  "xhigh",
		Workers:          1,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
	}, &out)
	if err != nil {
		t.Fatalf("generate with ranked parents: %v\n%s", err, out.String())
	}
	if calls != 2 {
		t.Fatalf("expected two child-generation calls, got %d", calls)
	}
	state = mustReadGAState(t, repo, "ga-generate-ranked")
	if state.Parents[0].SimID != "SIM-high" || !strings.Contains(state.Parents[0].Rationale, "avg_normalized_0_100=95.00") {
		t.Fatalf("parents were not fitness-ranked: %#v", state.Parents)
	}
	for _, child := range state.Children {
		if child.Operation != childOperationBreed || len(child.ParentIDs) != 2 || child.ParentIDs[0] != "SIM-high" {
			t.Fatalf("child did not breed top parent with tournament diversity: %#v", child)
		}
		if child.ParentIDs[1] == "SIM-high" {
			t.Fatalf("breed reused top parent twice: %#v", child)
		}
	}
}

func TestRunGenerateSkipsBreedChildWithTooFewParents(t *testing.T) {
	repo := newGAFixtureRepo(t)
	initGAStateForTest(t, repo, "ga-generate-one-parent")
	state := mustReadGAState(t, repo, "ga-generate-one-parent")
	state.Parents = state.Parents[:1]
	state.Children[0].Operation = "mutation"
	state.Children[0].ParentIDs = []string{state.Parents[0].SimID}
	writeTestState(t, repo, "ga-generate-one-parent", state)
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			return ProviderResponse{}, fmt.Errorf("provider should not be called for invalid breed parent count")
		},
	}
	var out strings.Builder
	err := runGenerateWithProvider(context.Background(), repo, provider, generateOptions{
		RunGroupID:         "ga-generate-one-parent",
		ProviderName:       "fake",
		APIModel:           "model-a",
		ReasoningEffort:    "xhigh",
		SkipFailedChildren: true,
		InputPrice:         defaultInputUSDPerMTok,
		CachedInputPrice:   defaultCachedInputUSDPerMTok,
		OutputPrice:        defaultOutputUSDPerMTok,
	}, &out)
	if err != nil {
		t.Fatalf("generate should skip invalid breed child: %v\n%s", err, out.String())
	}
	state = mustReadGAState(t, repo, "ga-generate-one-parent")
	if state.Children[0].Status != "skipped" || state.Children[0].Operation != childOperationBreed || !strings.Contains(state.Children[0].ValidationMessage, "exactly two parent IDs") {
		t.Fatalf("invalid breed child was not skipped with clear evidence: %#v", state.Children[0])
	}
}

func TestRunGenerateResetsStaleRunningChildrenBeforeDispatch(t *testing.T) {
	repo := newGAFixtureRepo(t)
	initGAStateForTestWithCounts(t, repo, "ga-generate-stale-running", 1, 3)
	state := mustReadGAState(t, repo, "ga-generate-stale-running")
	childIDs := make([]string, 0, len(state.Children))
	for index := range state.Children {
		childIDs = append(childIDs, state.Children[index].ID())
		markGAChild(&state.Children[index], "running", "stale interrupted child-generation attempt")
	}
	writeTestState(t, repo, "ga-generate-stale-running", state)
	calls := 0
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			running, err := countRunningChildren(repo, "ga-generate-stale-running")
			if err != nil {
				return ProviderResponse{}, err
			}
			// Intent: A restarted generation command must not leave every child
			// plan marked running before serialized workers actually own them.
			// Source: DI-juzus
			if running != 1 {
				return ProviderResponse{}, fmt.Errorf("expected one reclaimed running child, got %d", running)
			}
			childID := childIDs[calls]
			calls++
			return ProviderResponse{
				Text:        validChildBundleJSON(childID),
				RequestID:   "req-generate",
				ResponseID:  "resp-generate",
				ServiceTier: defaultServiceTier,
				UsageJSON:   `{"input_tokens":1200,"input_tokens_details":{"cached_tokens":200},"output_tokens":700}`,
			}, nil
		},
	}
	var out strings.Builder
	err := runGenerateWithProvider(context.Background(), repo, provider, generateOptions{
		RunGroupID:       "ga-generate-stale-running",
		ProviderName:     "fake",
		APIModel:         "model-a",
		ReasoningEffort:  "xhigh",
		Workers:          1,
		MaxOutputTokens:  6000,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
	}, &out)
	if err != nil {
		t.Fatalf("generate after stale running reset: %v\n%s", err, out.String())
	}
	if calls != 3 {
		t.Fatalf("expected 3 reclaimed child-generation calls, got %d", calls)
	}
}

func TestRunGenerateSkipsFailedChildAndChildScoreNoOps(t *testing.T) {
	repo := newGAFixtureRepo(t)
	initGAStateForTest(t, repo, "ga-generate-skip")
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			return ProviderResponse{}, fmt.Errorf("provider returned no child bundle")
		},
	}
	var out strings.Builder
	err := runGenerateWithProvider(context.Background(), repo, provider, generateOptions{
		RunGroupID:         "ga-generate-skip",
		ProviderName:       "fake",
		APIModel:           "model-a",
		ReasoningEffort:    "xhigh",
		MaxOutputTokens:    6000,
		InputPrice:         defaultInputUSDPerMTok,
		CachedInputPrice:   defaultCachedInputUSDPerMTok,
		OutputPrice:        defaultOutputUSDPerMTok,
		SkipFailedChildren: true,
	}, &out)
	if err != nil {
		t.Fatalf("generate with skip failed child: %v\n%s", err, out.String())
	}
	state := mustReadGAState(t, repo, "ga-generate-skip")
	if state.Children[0].Status != "skipped" || !strings.Contains(state.Children[0].ValidationMessage, "provider returned no child bundle") {
		t.Fatalf("child failure was not preserved as skipped: %#v", state.Children[0])
	}
	out.Reset()
	err = runScoreWithProvider(context.Background(), repo, provider, scoreOptions{
		RunGroupID:      "ga-generate-skip",
		Target:          "children",
		ProviderName:    "fake",
		APIModel:        "model-a",
		ReasoningEffort: "xhigh",
		SkipFailedCells: true,
	}, &out)
	if err != nil {
		t.Fatalf("child scoring should no-op when all child generation was skipped: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "processed=0") {
		t.Fatalf("expected child scoring no-op output, got %q", out.String())
	}
}

func TestServiceTierOptionsDefaultToFlexAndRejectPriority(t *testing.T) {
	scoreDefaults, err := parseScoreOptions([]string{"-run-group-id", "ga-score", "-api-model", "model-a"})
	if err != nil {
		t.Fatalf("parse score defaults: %v", err)
	}
	if scoreDefaults.ServiceTier != serviceTierFlex {
		t.Fatalf("score service tier = %q, want %q", scoreDefaults.ServiceTier, serviceTierFlex)
	}
	if scoreDefaults.ReasoningEffort != defaultScoreReasoningEffort || scoreDefaults.TextVerbosity != defaultTextVerbosity || scoreDefaults.MaxOutputTokens != 0 || scoreDefaults.CostEstimateOutputTokens != defaultScoreCostEstimateOutputTokens {
		t.Fatalf("score request-shaping defaults not applied: %#v", scoreDefaults)
	}
	if scoreDefaults.Workers != defaultScoreWorkers || scoreDefaults.RequestTimeout != defaultRequestTimeout || scoreDefaults.ProviderAttempts != defaultProviderMaxAttempts || scoreDefaults.ProviderElapsed != defaultProviderMaxElapsed {
		t.Fatalf("score throughput defaults not applied: %#v", scoreDefaults)
	}
	generateDefaults, err := parseGenerateOptions([]string{"-run-group-id", "ga-generate", "-api-model", "model-a"})
	if err != nil {
		t.Fatalf("parse generate defaults: %v", err)
	}
	if generateDefaults.ServiceTier != serviceTierFlex {
		t.Fatalf("generate service tier = %q, want %q", generateDefaults.ServiceTier, serviceTierFlex)
	}
	if generateDefaults.ReasoningEffort != defaultGenerateReasoningEffort || generateDefaults.TextVerbosity != defaultTextVerbosity || generateDefaults.MaxOutputTokens != 0 || generateDefaults.CostEstimateOutputTokens != defaultGenerateCostEstimateOutputTokens {
		t.Fatalf("generate request-shaping defaults not applied: %#v", generateDefaults)
	}
	if generateDefaults.Workers != defaultGenerateWorkers || generateDefaults.RequestTimeout != defaultRequestTimeout || generateDefaults.ProviderAttempts != defaultProviderMaxAttempts || generateDefaults.ProviderElapsed != defaultProviderMaxElapsed {
		t.Fatalf("generate throughput defaults not applied: %#v", generateDefaults)
	}
	scoreDefaultTier, err := parseScoreOptions([]string{"-run-group-id", "ga-score", "-api-model", "model-a", "-service-tier", "default"})
	if err != nil {
		t.Fatalf("parse explicit default tier: %v", err)
	}
	if scoreDefaultTier.ServiceTier != serviceTierDefault {
		t.Fatalf("score explicit service tier = %q, want %q", scoreDefaultTier.ServiceTier, serviceTierDefault)
	}
	scoreCustom, err := parseScoreOptions([]string{"-run-group-id", "ga-score", "-api-model", "model-a", "-workers", "3", "-request-timeout", "2m", "-provider-max-attempts", "4", "-provider-max-elapsed", "9m", "-text-verbosity", "high", "-max-output-tokens", "123", "-cost-estimate-output-tokens", "456"})
	if err != nil {
		t.Fatalf("parse custom throughput knobs: %v", err)
	}
	if scoreCustom.Workers != 3 || scoreCustom.RequestTimeout != 2*time.Minute || scoreCustom.ProviderAttempts != 4 || scoreCustom.ProviderElapsed != 9*time.Minute {
		t.Fatalf("custom throughput knobs not applied: %#v", scoreCustom)
	}
	if scoreCustom.TextVerbosity != textVerbosityHigh || scoreCustom.MaxOutputTokens != 123 || scoreCustom.CostEstimateOutputTokens != 456 {
		t.Fatalf("custom request-shaping knobs not applied: %#v", scoreCustom)
	}
	if _, err := parseGenerateOptions([]string{"-run-group-id", "ga-generate", "-api-model", "model-a", "-service-tier", "priority"}); err == nil {
		t.Fatalf("expected priority service tier to be rejected")
	}
	if _, err := parseScoreOptions([]string{"-run-group-id", "ga-score", "-api-model", "model-a", "-workers", "0"}); err == nil {
		t.Fatalf("expected zero workers to be rejected")
	}
	if _, err := parseGenerateOptions([]string{"-run-group-id", "ga-generate", "-api-model", "model-a", "-text-verbosity", "verbose"}); err == nil {
		t.Fatalf("expected invalid text verbosity to be rejected")
	}
	if _, err := parseScoreOptions([]string{"-run-group-id", "ga-score", "-api-model", "model-a", "-max-output-tokens", "-1"}); err == nil {
		t.Fatalf("expected negative max output tokens to be rejected")
	}
}

func TestOpenAIProviderSendsExplicitServiceTier(t *testing.T) {
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requestBytes, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}
		var body openAIRequest
		if err := json.Unmarshal(requestBytes, &body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if body.ServiceTier != serviceTierFlex {
			t.Fatalf("request service tier = %q, want %q", body.ServiceTier, serviceTierFlex)
		}
		if body.Text == nil || body.Text.Verbosity != defaultTextVerbosity {
			t.Fatalf("request text verbosity = %#v, want %q", body.Text, defaultTextVerbosity)
		}
		if strings.Contains(string(requestBytes), "max_output_tokens") {
			t.Fatalf("default request unexpectedly sent hard max_output_tokens: %s", string(requestBytes))
		}
		return openAITestHTTPResponse(request, http.StatusOK, "req-tier", openAITestSuccessBody(serviceTierFlex)), nil
	})

	provider := OpenAIProvider{
		APIKey:  "test-key",
		BaseURL: "https://example.test/responses",
		Client:  &http.Client{Transport: transport},
		RetryPolicy: ProviderRetryPolicy{
			MaxAttempts: 1,
		},
	}
	response, err := provider.Generate(context.Background(), ProviderRequest{
		APIModel:    "model-a",
		ServiceTier: serviceTierFlex,
		Prompt:      "score this",
	})
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if response.ServiceTier != serviceTierFlex || response.RequestID != "req-tier" {
		t.Fatalf("unexpected response metadata: %#v", response)
	}
}

func TestOpenAIProviderLogsRequestAndResponseJSON(t *testing.T) {
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		return openAITestHTTPResponse(request, http.StatusOK, "req-debug", openAITestSuccessBody(serviceTierDefault)), nil
	})

	var debug strings.Builder
	provider := OpenAIProvider{
		APIKey:      "test-key",
		BaseURL:     "https://example.test/responses",
		Client:      &http.Client{Transport: transport},
		DebugWriter: &debug,
		RetryPolicy: ProviderRetryPolicy{
			MaxAttempts: 1,
		},
	}
	response, err := provider.Generate(context.Background(), ProviderRequest{
		APIModel:        "model-a",
		ServiceTier:     serviceTierDefault,
		ReasoningEffort: "xhigh",
		MaxOutputTokens: 12,
		Instructions:    "reply in json",
		Prompt:          "score this",
	})
	if err != nil {
		t.Fatalf("generate with debug logging: %v", err)
	}
	if response.ServiceTier != serviceTierDefault || response.RequestID != "req-debug" {
		t.Fatalf("unexpected response metadata: %#v", response)
	}
	logText := debug.String()
	for _, want := range []string{
		"event=request",
		"query_json=",
		`"model":"model-a"`,
		`"service_tier":"default"`,
		`"text":{"verbosity":"low"}`,
		`"reasoning":{"effort":"xhigh"}`,
		`"max_output_tokens":12`,
		"event=response",
		"status=200",
		`request_id="req-debug"`,
		"response_json=",
		`"service_tier":"default"`,
	} {
		if !strings.Contains(logText, want) {
			t.Fatalf("debug log missing %q:\n%s", want, logText)
		}
	}
	if strings.Contains(logText, "test-key") {
		t.Fatalf("debug log leaked API key:\n%s", logText)
	}
}

func TestOpenAIProviderRetriesFlex429WithoutChangingTier(t *testing.T) {
	attempts := 0
	var observedServiceTiers []string
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		attempts++
		var body openAIRequest
		if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
			t.Fatalf("decode retry body: %v", err)
		}
		observedServiceTiers = append(observedServiceTiers, body.ServiceTier)
		if attempts == 1 {
			return openAITestHTTPResponse(request, http.StatusTooManyRequests, "", `{"error":{"message":"Resource Unavailable"}}`), nil
		}
		return openAITestHTTPResponse(request, http.StatusOK, "", openAITestSuccessBody(serviceTierFlex)), nil
	})

	provider := OpenAIProvider{
		APIKey:  "test-key",
		BaseURL: "https://example.test/responses",
		Client:  &http.Client{Transport: transport},
		RetryPolicy: ProviderRetryPolicy{
			MaxAttempts:    2,
			MaxElapsed:     time.Second,
			InitialBackoff: time.Millisecond,
			MaxBackoff:     time.Millisecond,
		},
	}
	if _, err := provider.Generate(context.Background(), ProviderRequest{APIModel: "model-a", ServiceTier: serviceTierFlex, Prompt: "retry"}); err != nil {
		t.Fatalf("generate after retry: %v", err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
	for _, serviceTier := range observedServiceTiers {
		if serviceTier != serviceTierFlex {
			t.Fatalf("retry changed service tier: %v", observedServiceTiers)
		}
	}
}

func TestOpenAIProviderRetriesEmptyOutput(t *testing.T) {
	attempts := 0
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		attempts++
		if attempts == 1 {
			return openAITestHTTPResponse(request, http.StatusOK, "", `{"id":"resp-empty","status":"completed","usage":{"input_tokens":10,"input_tokens_details":{"cached_tokens":0},"output_tokens":0}}`), nil
		}
		return openAITestHTTPResponse(request, http.StatusOK, "", openAITestSuccessBody(serviceTierFlex)), nil
	})
	provider := OpenAIProvider{
		APIKey:  "test-key",
		BaseURL: "https://example.test/responses",
		Client:  &http.Client{Transport: transport},
		RetryPolicy: ProviderRetryPolicy{
			MaxAttempts:    2,
			MaxElapsed:     time.Second,
			InitialBackoff: time.Millisecond,
			MaxBackoff:     time.Millisecond,
		},
	}
	if _, err := provider.Generate(context.Background(), ProviderRequest{APIModel: "model-a", ServiceTier: serviceTierFlex, Prompt: "retry empty"}); err != nil {
		t.Fatalf("generate after empty-output retry: %v", err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}

func TestOpenAIProviderDoesNotRetryMaxOutputIncomplete(t *testing.T) {
	attempts := 0
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		attempts++
		body := `{"id":"resp-incomplete","status":"incomplete","incomplete_details":{"reason":"max_output_tokens"},"service_tier":"default","output":[{"type":"reasoning","summary":[]}],"usage":{"input_tokens":10,"input_tokens_details":{"cached_tokens":0},"output_tokens":4000,"output_tokens_details":{"reasoning_tokens":4000}}}`
		return openAITestHTTPResponse(request, http.StatusOK, "req-incomplete", body), nil
	})
	var debug strings.Builder
	provider := OpenAIProvider{
		APIKey:      "test-key",
		BaseURL:     "https://example.test/responses",
		Client:      &http.Client{Transport: transport},
		DebugWriter: &debug,
		RetryPolicy: ProviderRetryPolicy{
			MaxAttempts:    2,
			MaxElapsed:     time.Second,
			InitialBackoff: time.Millisecond,
			MaxBackoff:     time.Millisecond,
		},
	}
	_, err := provider.Generate(context.Background(), ProviderRequest{
		APIModel:        "model-a",
		ServiceTier:     serviceTierDefault,
		ReasoningEffort: "xhigh",
		Prompt:          "retry would waste budget",
	})
	if err == nil {
		t.Fatalf("expected max-output incomplete response to fail")
	}
	if attempts != 1 {
		t.Fatalf("max-output incomplete response retried %d attempts; want 1", attempts)
	}
	if !strings.Contains(err.Error(), `reason "max_output_tokens"`) || !strings.Contains(err.Error(), `"output_tokens":4000`) {
		t.Fatalf("error did not preserve cap/usage evidence: %v", err)
	}
	if strings.Contains(debug.String(), "event=retry_sleep") {
		t.Fatalf("max-output incomplete response should not retry:\n%s", debug.String())
	}
}

func TestOpenAIProviderRetriesTimeout(t *testing.T) {
	transport := &timeoutThenSuccessTransport{}
	provider := OpenAIProvider{
		APIKey: "test-key",
		Client: &http.Client{Transport: transport},
		RetryPolicy: ProviderRetryPolicy{
			MaxAttempts:    2,
			MaxElapsed:     time.Second,
			InitialBackoff: time.Millisecond,
			MaxBackoff:     time.Millisecond,
		},
	}
	response, err := provider.Generate(context.Background(), ProviderRequest{APIModel: "model-a", ServiceTier: serviceTierFlex, Prompt: "retry timeout"})
	if err != nil {
		t.Fatalf("generate after timeout retry: %v", err)
	}
	if response.ServiceTier != serviceTierFlex || transport.Attempts != 2 {
		t.Fatalf("unexpected timeout retry response=%#v attempts=%d", response, transport.Attempts)
	}
}

func TestOpenAIProviderFailsAfterBoundedFlex429Retries(t *testing.T) {
	attempts := 0
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		attempts++
		return openAITestHTTPResponse(request, http.StatusTooManyRequests, "", `{"error":{"message":"Resource Unavailable"}}`), nil
	})

	provider := OpenAIProvider{
		APIKey:  "test-key",
		BaseURL: "https://example.test/responses",
		Client:  &http.Client{Transport: transport},
		RetryPolicy: ProviderRetryPolicy{
			MaxAttempts:    2,
			MaxElapsed:     time.Second,
			InitialBackoff: time.Millisecond,
			MaxBackoff:     time.Millisecond,
		},
	}
	_, err := provider.Generate(context.Background(), ProviderRequest{APIModel: "model-a", ServiceTier: serviceTierFlex, Prompt: "fail"})
	if err == nil {
		t.Fatalf("expected bounded retry failure")
	}
	if attempts != 2 || !strings.Contains(err.Error(), "after 2 attempts") {
		t.Fatalf("unexpected retry failure attempts=%d err=%v", attempts, err)
	}
}

func TestRunGenerateRejectsUnsafeBundlePath(t *testing.T) {
	repo := newGAFixtureRepo(t)
	initGAStateForTest(t, repo, "ga-bad-generate")
	state := mustReadGAState(t, repo, "ga-bad-generate")
	childID := state.Children[0].ID()
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			return ProviderResponse{
				Text: fmt.Sprintf(`{"child_id":%q,"design_delta_summary":"bad path","files":[{"path":"../README.md","content":"bad"},{"path":"QUESTION.md","content":"# Q"}]}`, childID),
			}, nil
		},
	}
	var out strings.Builder
	err := runGenerateWithProvider(context.Background(), repo, provider, generateOptions{
		RunGroupID:      "ga-bad-generate",
		ProviderName:    "fake",
		APIModel:        "model-a",
		ReasoningEffort: "xhigh",
		MaxOutputTokens: 6000,
	}, &out)
	if err == nil {
		t.Fatalf("expected unsafe bundle path error")
	}
	assertMissing(t, repo.Path("simulations", childID))
}

func TestRunAcceptRecordsAcceptanceAndPrintsStagePaths(t *testing.T) {
	repo := newGitTestRepo(t)
	resultPath := writeAcceptFixture(t, repo, "SIM-child", "ga-test")

	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "accept",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-test",
		"-child", "SIM-child",
		"-result", repo.Rel(resultPath),
		"-reviewer-note", "good enough to promote",
	}, &out, &out)
	if err != nil {
		t.Fatalf("accept child: %v\n%s", err, out.String())
	}

	state, err := readGAState(repo.Path("results", "state", "ga-test.json"))
	if err != nil {
		t.Fatalf("read updated state: %v", err)
	}
	if len(state.Acceptance) != 1 {
		t.Fatalf("expected one acceptance record, got %#v", state.Acceptance)
	}
	if state.Children[0].Status != "accepted" {
		t.Fatalf("expected child status accepted, got %s", state.Children[0].Status)
	}
	text := out.String()
	for _, want := range []string{
		"accepted children=SIM-child",
		"results/state/ga-test.json",
		"simulations/SIM-child",
		repo.Rel(resultPath),
		"git add --",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("expected %q in accept output:\n%s", want, text)
		}
	}
	if cached := gitOutput(t, repo, "diff", "--cached", "--name-only"); strings.TrimSpace(cached) != "" {
		t.Fatalf("accept must not stage changes; cached diff:\n%s", cached)
	}
}

func TestRunAcceptRejectsMissingState(t *testing.T) {
	repo := newTestRepo(t)
	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "accept",
		"-repo-root", repo.Root,
		"-run-group-id", "missing",
		"-child", "SIM-child",
		"-result", "results/SIM-child/scenario-one/model-a/20260519-101500.json",
		"-reviewer-note", "note",
	}, &out, &out)
	if err == nil {
		t.Fatalf("expected missing state error")
	}
}

func TestRunAcceptRejectsUnknownChild(t *testing.T) {
	repo := newTestRepo(t)
	resultPath := writeAcceptFixture(t, repo, "SIM-child", "ga-test")
	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "accept",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-test",
		"-child", "SIM-missing",
		"-result", repo.Rel(resultPath),
		"-reviewer-note", "note",
	}, &out, &out)
	if err == nil || !strings.Contains(err.Error(), "is not present in GA state") {
		t.Fatalf("expected unknown child error, got %v", err)
	}
}

func TestRunAcceptRejectsInvalidResult(t *testing.T) {
	repo := newTestRepo(t)
	resultPath := writeAcceptFixture(t, repo, "SIM-child", "ga-test")
	result := validResult(repo, resultPath)
	result.Scores.PromiseGridAlignment = 9
	if err := writeFitnessResultAtomic(resultPath, result); err != nil {
		t.Fatalf("write invalid result: %v", err)
	}

	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "accept",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-test",
		"-child", "SIM-child",
		"-result", repo.Rel(resultPath),
		"-reviewer-note", "note",
	}, &out, &out)
	if err == nil || !strings.Contains(err.Error(), "scores.promisegrid_alignment must be between 0 and 5") {
		t.Fatalf("expected invalid result error, got %v", err)
	}
}

func TestRunAcceptRejectsResultChildMismatch(t *testing.T) {
	repo := newTestRepo(t)
	writeAcceptFixture(t, repo, "SIM-child", "ga-test")
	otherResultPath := repo.Path("results", "SIM-other", "scenario-one", "model-a", "20260519-101500.json")
	otherResult := validResult(repo, otherResultPath)
	if err := writeFitnessResultAtomic(otherResultPath, otherResult); err != nil {
		t.Fatalf("write other result: %v", err)
	}

	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "accept",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-test",
		"-child", "SIM-child",
		"-result", repo.Rel(otherResultPath),
		"-reviewer-note", "note",
	}, &out, &out)
	if err == nil || !strings.Contains(err.Error(), "is not a selected child") {
		t.Fatalf("expected child mismatch error, got %v", err)
	}
}

func TestRunAcceptRejectsChildHashDrift(t *testing.T) {
	repo := newTestRepo(t)
	resultPath := writeAcceptFixture(t, repo, "SIM-child", "ga-test")
	state, err := readGAState(repo.Path("results", "state", "ga-test.json"))
	if err != nil {
		t.Fatalf("read state: %v", err)
	}
	state.Children[0].TreeHash = strings.Repeat("0", 64)
	if err := writeGAStateAtomic(repo.Path("results", "state", "ga-test.json"), state); err != nil {
		t.Fatalf("write drifted state: %v", err)
	}

	var out strings.Builder
	err = runMain([]string{
		"ga-runner", "accept",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-test",
		"-child", "SIM-child",
		"-result", repo.Rel(resultPath),
		"-reviewer-note", "note",
	}, &out, &out)
	if err == nil || !strings.Contains(err.Error(), "tree hash drift") {
		t.Fatalf("expected tree hash drift error, got %v", err)
	}
}

func TestRunCullDeletesRejectedChildAndRecordsState(t *testing.T) {
	repo := newTestRepo(t)
	resultPath := writeAcceptFixture(t, repo, "SIM-child", "ga-test")
	unrelatedPath := repo.Path("results", "SIM-other", "scenario-one", "model-a", "20260519-101500.json")
	writeTestFile(t, unrelatedPath, "{}\n")

	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "cull",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-test",
		"-child", "SIM-child",
		"-reason", "not good enough",
	}, &out, &out)
	if err != nil {
		t.Fatalf("cull child: %v\n%s", err, out.String())
	}
	assertMissing(t, repo.Path("simulations", "SIM-child"))
	assertMissing(t, filepath.Dir(filepath.Dir(filepath.Dir(resultPath))))
	assertExists(t, unrelatedPath)

	state, err := readGAState(repo.Path("results", "state", "ga-test.json"))
	if err != nil {
		t.Fatalf("read updated state: %v", err)
	}
	if state.Children[0].Status != "culled" {
		t.Fatalf("expected child status culled, got %s", state.Children[0].Status)
	}
	if len(state.Culling) != 1 || state.Culling[0].Reason != "not good enough" {
		t.Fatalf("expected culling record, got %#v", state.Culling)
	}
	text := out.String()
	for _, want := range []string{"cull children=SIM-child", "simulations/SIM-child", "results/SIM-child"} {
		if !strings.Contains(text, want) {
			t.Fatalf("expected %q in cull output:\n%s", want, text)
		}
	}
}

func TestRunCullDryRunDeletesNothingAndWritesNoState(t *testing.T) {
	repo := newTestRepo(t)
	resultPath := writeAcceptFixture(t, repo, "SIM-child", "ga-test")

	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "cull",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-test",
		"-child", "SIM-child",
		"-reason", "not good enough",
		"-dry-run",
	}, &out, &out)
	if err != nil {
		t.Fatalf("dry-run cull child: %v\n%s", err, out.String())
	}
	assertExists(t, repo.Path("simulations", "SIM-child"))
	assertExists(t, resultPath)
	state, err := readGAState(repo.Path("results", "state", "ga-test.json"))
	if err != nil {
		t.Fatalf("read state: %v", err)
	}
	if state.Children[0].Status == "culled" || len(state.Culling) != 0 {
		t.Fatalf("dry-run must not write culling state: %#v", state)
	}
	if !strings.Contains(out.String(), "dry-run children=SIM-child") {
		t.Fatalf("expected dry-run output, got:\n%s", out.String())
	}
}

func TestRunCullRejectsUnknownChild(t *testing.T) {
	repo := newTestRepo(t)
	writeAcceptFixture(t, repo, "SIM-child", "ga-test")
	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "cull",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-test",
		"-child", "SIM-missing",
		"-reason", "reject",
	}, &out, &out)
	if err == nil || !strings.Contains(err.Error(), "is not present in GA state") {
		t.Fatalf("expected unknown child error, got %v", err)
	}
}

func TestRunCullRejectsAcceptedChild(t *testing.T) {
	repo := newTestRepo(t)
	writeAcceptFixture(t, repo, "SIM-child", "ga-test")
	state := mustReadGAState(t, repo, "ga-test")
	state.Children[0].Status = "accepted"
	writeTestState(t, repo, "ga-test", state)

	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "cull",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-test",
		"-child", "SIM-child",
		"-reason", "reject",
	}, &out, &out)
	if err == nil || !strings.Contains(err.Error(), "is accepted and cannot be culled") {
		t.Fatalf("expected accepted child error, got %v", err)
	}
}

func TestRunCullRejectsAlreadyCulledChild(t *testing.T) {
	repo := newTestRepo(t)
	writeAcceptFixture(t, repo, "SIM-child", "ga-test")
	state := mustReadGAState(t, repo, "ga-test")
	state.Children[0].Status = "culled"
	writeTestState(t, repo, "ga-test", state)

	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "cull",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-test",
		"-child", "SIM-child",
		"-reason", "reject",
	}, &out, &out)
	if err == nil || !strings.Contains(err.Error(), "is already culled") {
		t.Fatalf("expected already culled error, got %v", err)
	}
}

func TestRunCullRejectsUnsafeChildPath(t *testing.T) {
	repo := newTestRepo(t)
	writeAcceptFixture(t, repo, "SIM-child", "ga-test")
	state := mustReadGAState(t, repo, "ga-test")
	state.Children[0].Path = "results/SIM-child/"
	writeTestState(t, repo, "ga-test", state)

	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "cull",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-test",
		"-child", "SIM-child",
		"-reason", "reject",
	}, &out, &out)
	if err == nil || !strings.Contains(err.Error(), "path must be simulations/SIM-child/") {
		t.Fatalf("expected unsafe child path error, got %v", err)
	}
}

func newTestRepo(t *testing.T) Repo {
	t.Helper()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatalf("make .git: %v", err)
	}
	return Repo{Root: root}
}

func newGitTestRepo(t *testing.T) Repo {
	t.Helper()
	root := t.TempDir()
	runGit(t, root, "init")
	return Repo{Root: root}
}

func newGAFixtureRepo(t *testing.T) Repo {
	t.Helper()
	repo := newGitTestRepo(t)
	writeTestFile(t, repo.Path("results", "RUN-PROTOCOL.md"), "# Run Protocol\n\nScore cells as evidence.\n")
	writeTestFile(t, repo.Path("scenarios", "README.md"), "# Scenarios\n\nUse 100-year PromiseGrid pressure.\n")
	writeTestFile(t, repo.Path("scenarios", "scenario-one", "scenario-one.md"), "# Scenario One\n\nAlice asks Bob to ship labels with auditable promises.\n")
	writeTestFile(t, repo.Path("simulations", "SIM-parent", "README.md"), "# Parent Sim\n\nA small parent simulation for GA tests.\n")
	writeTestFile(t, repo.Path("simulations", "SIM-parent", "QUESTION.md"), "# Questions\n\nCan Alice and Bob audit the result?\n")
	writeTestFile(t, repo.Path("simulations", "SIM-second", "README.md"), "# Second Parent Sim\n\nAnother small parent simulation for GA breed tests.\n")
	writeTestFile(t, repo.Path("simulations", "SIM-second", "QUESTION.md"), "# Questions\n\nCan Carol compare the two parent designs?\n")
	gitAdd(t, repo,
		"results/RUN-PROTOCOL.md",
		"scenarios/README.md",
		"scenarios/scenario-one/scenario-one.md",
		"simulations/SIM-parent/README.md",
		"simulations/SIM-parent/QUESTION.md",
		"simulations/SIM-second/README.md",
		"simulations/SIM-second/QUESTION.md",
	)
	return repo
}

func initGAStateForTest(t *testing.T, repo Repo, runGroupID string) {
	t.Helper()
	initGAStateForTestWithScenarioCount(t, repo, runGroupID, 1)
}

func initGAStateForTestWithScenarioCount(t *testing.T, repo Repo, runGroupID string, scenarioCount int) {
	t.Helper()
	initGAStateForTestWithCounts(t, repo, runGroupID, scenarioCount, 1)
}

func initGAStateForTestWithCounts(t *testing.T, repo Repo, runGroupID string, scenarioCount int, childCount int) {
	t.Helper()
	initGAStateForTestWithParentsScenariosAndChildren(t, repo, runGroupID, 2, scenarioCount, childCount)
}

func initGAStateForTestWithParentsAndChildren(t *testing.T, repo Repo, runGroupID string, parentCount int, childCount int) {
	t.Helper()
	initGAStateForTestWithParentsScenariosAndChildren(t, repo, runGroupID, parentCount, 1, childCount)
}

func initGAStateForTestWithParentsScenariosAndChildren(t *testing.T, repo Repo, runGroupID string, parentCount int, scenarioCount int, childCount int) {
	t.Helper()
	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "init",
		"-repo-root", repo.Root,
		"-model", "model-a",
		"-run-group-id", runGroupID,
		"-timestamp", "20260519-111500",
		"-parent-count", fmt.Sprintf("%d", parentCount),
		"-scenario-count", fmt.Sprintf("%d", scenarioCount),
		"-child-count", fmt.Sprintf("%d", childCount),
		"-max-promotions", "1",
	}, &out, &out)
	if err != nil {
		t.Fatalf("init fixture state: %v\n%s", err, out.String())
	}
}

func addTrackedFixtureSim(t *testing.T, repo Repo, simID string) {
	t.Helper()
	writeTestFile(t, repo.Path("simulations", simID, "README.md"), "# "+simID+"\n\nA fixture simulation for ranked parent selection.\n")
	writeTestFile(t, repo.Path("simulations", simID, "QUESTION.md"), "# Questions\n\nCan Alice and Bob improve the score?\n")
	gitAdd(t, repo,
		filepath.ToSlash(filepath.Join("simulations", simID, "README.md")),
		filepath.ToSlash(filepath.Join("simulations", simID, "QUESTION.md")),
	)
}

func writeParentFitnessResults(t *testing.T, repo Repo, runGroupID string, normalizedBySimID map[string]float64) {
	t.Helper()
	state := mustReadGAState(t, repo, runGroupID)
	parentIDs := map[string]bool{}
	for _, parent := range state.Parents {
		parentIDs[parent.SimID] = true
	}
	for index := range state.Cells {
		cell := &state.Cells[index]
		if !parentIDs[cell.SimID] {
			continue
		}
		result := validResult(repo, repo.Abs(cell.ResultPath))
		result.RunGroupID = state.RunGroupID
		result.CellID = cell.CellID
		result.SimID = cell.SimID
		result.ScenarioID = cell.ScenarioID
		result.ModelID = cell.ModelID
		result.ResultPath = cell.ResultPath
		result.Source.SimPath = filepath.ToSlash(filepath.Join("simulations", cell.SimID)) + "/"
		result.Fitness.Normalized0To100 = normalizedBySimID[cell.SimID]
		if err := writeFitnessResultAtomic(repo.Abs(cell.ResultPath), result); err != nil {
			t.Fatalf("write parent result for %s: %v", cell.SimID, err)
		}
		cell.Status = "done"
	}
	writeTestState(t, repo, runGroupID, state)
}

func promptChildID(t *testing.T, prompt string) string {
	t.Helper()
	prefix := "- Child ID: `"
	start := strings.Index(prompt, prefix)
	if start < 0 {
		t.Fatalf("prompt missing child ID: %s", prompt)
	}
	start += len(prefix)
	end := strings.Index(prompt[start:], "`")
	if end < 0 {
		t.Fatalf("prompt has unterminated child ID: %s", prompt)
	}
	return prompt[start : start+end]
}

func countRunningParentCells(repo Repo, runGroupID string) (int, error) {
	state, err := readGAState(repo.Path("results", "state", runGroupID+".json"))
	if err != nil {
		return 0, err
	}
	parentIDs := map[string]bool{}
	for _, parent := range state.Parents {
		parentIDs[parent.SimID] = true
	}
	running := 0
	for _, cell := range state.Cells {
		if parentIDs[cell.SimID] && cell.Status == "running" {
			running++
		}
	}
	return running, nil
}

func countRunningChildren(repo Repo, runGroupID string) (int, error) {
	state, err := readGAState(repo.Path("results", "state", runGroupID+".json"))
	if err != nil {
		return 0, err
	}
	running := 0
	for _, child := range state.Children {
		if child.Status == "running" {
			running++
		}
	}
	return running, nil
}

func addFixtureScenarios(t *testing.T, repo Repo, scenarioIDs ...string) {
	t.Helper()
	var paths []string
	for _, scenarioID := range scenarioIDs {
		rel := filepath.ToSlash(filepath.Join("scenarios", scenarioID, scenarioID+".md"))
		writeTestFile(t, repo.Path(rel), "# "+scenarioID+"\n\nAlice and Bob exercise "+scenarioID+".\n")
		paths = append(paths, rel)
	}
	gitAdd(t, repo, paths...)
}

type fakeGAProvider struct {
	generate func(context.Context, ProviderRequest) (ProviderResponse, error)
}

func (provider fakeGAProvider) Generate(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
	if provider.generate == nil {
		return ProviderResponse{}, fmt.Errorf("fake provider has no generate function")
	}
	return provider.generate(ctx, request)
}

func writeTestFile(t *testing.T, path string, text string) {
	t.Helper()
	if err := ensureParent(path); err != nil {
		t.Fatalf("make parent for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func gitAdd(t *testing.T, repo Repo, paths ...string) {
	t.Helper()
	args := append([]string{"add", "--"}, paths...)
	runGit(t, repo.Root, args...)
}

func runGit(t *testing.T, root string, args ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root}, args...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(output))
	}
}

func gitOutput(t *testing.T, repo Repo, args ...string) string {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", repo.Root}, args...)...)
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, string(output))
	}
	return string(output)
}

func writeAcceptFixture(t *testing.T, repo Repo, childID string, runGroupID string) string {
	t.Helper()
	childPath := filepath.ToSlash(filepath.Join("simulations", childID))
	writeTestFile(t, repo.Path(childPath, "README.md"), "# Child\n")
	writeTestFile(t, repo.Path(childPath, "QUESTION.md"), "# Question\n")
	treeHash, err := currentSimulationTreeHash(repo, childPath)
	if err != nil {
		t.Fatalf("hash child tree: %v", err)
	}
	resultPath := repo.Path("results", childID, "scenario-one", "model-a", "20260519-101500.json")
	result := validResult(repo, resultPath)
	result.Source.SimPath = childPath + "/"
	result.Source.SimulationTreeHash = treeHash
	if err := writeFitnessResultAtomic(resultPath, result); err != nil {
		t.Fatalf("write result: %v", err)
	}
	state := GAState{
		Schema:     stateSchemaV1,
		RunGroupID: runGroupID,
		CreatedAt:  "2026-05-19T10:43:01Z",
		UpdatedAt:  "2026-05-19T10:43:01Z",
		RepoCommit: "abc123",
		ModelID:    "model-a",
		Children: []GAChild{
			{
				ChildID:  childID,
				Path:     childPath + "/",
				TreeHash: treeHash,
				Status:   "scored",
			},
		},
		Cells: []GACell{
			{
				CellID:     "cell-1",
				SimID:      childID,
				ScenarioID: "scenario-one",
				ModelID:    "model-a",
				ResultPath: repo.Rel(resultPath),
				Status:     "done",
			},
		},
	}
	if err := writeGAStateAtomic(repo.Path("results", "state", runGroupID+".json"), state); err != nil {
		t.Fatalf("write state: %v", err)
	}
	return resultPath
}

func mustReadGAState(t *testing.T, repo Repo, runGroupID string) GAState {
	t.Helper()
	state, err := readGAState(repo.Path("results", "state", runGroupID+".json"))
	if err != nil {
		t.Fatalf("read state: %v", err)
	}
	return state
}

func mustReadFitnessResult(t *testing.T, path string) FitnessResult {
	t.Helper()
	bytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read result: %v", err)
	}
	var result FitnessResult
	if err := json.Unmarshal(bytes, &result); err != nil {
		t.Fatalf("parse result: %v", err)
	}
	return result
}

func writeTestState(t *testing.T, repo Repo, runGroupID string, state GAState) {
	t.Helper()
	if err := writeGAStateAtomic(repo.Path("results", "state", runGroupID+".json"), state); err != nil {
		t.Fatalf("write state: %v", err)
	}
}

func assertMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be missing, stat err=%v", path, err)
	}
}

func assertExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
}

func parentIDs(parents []PopulationSim) []string {
	var ids []string
	for _, parent := range parents {
		ids = append(ids, parent.SimID)
	}
	return ids
}

func scenarioIDs(scenarios []Scenario) []string {
	var ids []string
	for _, scenario := range scenarios {
		ids = append(ids, scenario.ScenarioID)
	}
	return ids
}

func openAITestSuccessBody(serviceTier string) string {
	return fmt.Sprintf(`{"id":"resp-test","status":"completed","output_text":"{}","service_tier":%q,"usage":{"input_tokens":10,"input_tokens_details":{"cached_tokens":2},"output_tokens":3}}`, serviceTier)
}

func openAITestHTTPResponse(request *http.Request, statusCode int, requestID string, body string) *http.Response {
	header := make(http.Header)
	if requestID != "" {
		header.Set("x-request-id", requestID)
	}
	return &http.Response{
		StatusCode: statusCode,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    request,
	}
}

type roundTripFunc func(request *http.Request) (*http.Response, error)

func (roundTrip roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return roundTrip(request)
}

type timeoutThenSuccessTransport struct {
	Attempts int
}

func (transport *timeoutThenSuccessTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	transport.Attempts++
	if transport.Attempts == 1 {
		return nil, timeoutTestError{}
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(openAITestSuccessBody(serviceTierFlex))),
		Request:    request,
	}, nil
}

type timeoutTestError struct{}

func (timeoutTestError) Error() string {
	return "test timeout"
}

func (timeoutTestError) Timeout() bool {
	return true
}

func (timeoutTestError) Temporary() bool {
	return true
}

func validResult(repo Repo, path string) FitnessResult {
	parts, _ := parseResultPath(repo, path)
	rel := repo.Rel(path)
	return FitnessResult{
		Schema:       resultSchemaV1,
		ResultID:     strings.Join([]string{parts.SimID, parts.ScenarioID, parts.ModelID, parts.Timestamp}, "-"),
		RunGroupID:   "ga-20260519-101500",
		CellID:       "ga-20260519-101500-000001",
		SimID:        parts.SimID,
		ScenarioID:   parts.ScenarioID,
		ModelID:      parts.ModelID,
		TimestampUTC: parts.Timestamp,
		ResultPath:   rel,
		Runner: RunnerInfo{
			Tool:            "ga-runner",
			Provider:        "openai",
			APIModel:        parts.ModelID,
			ReasoningEffort: "xhigh",
		},
		Source: SourceInfo{
			RepoCommit:        "abc123",
			SimPath:           "simulations/" + parts.SimID + "/",
			ScenarioPath:      "scenarios/" + parts.ScenarioID + "/" + parts.ScenarioID + ".md",
			RootContractPaths: []string{"scenarios/README.md", "results/RUN-PROTOCOL.md"},
			Files: []SourceFile{
				{Path: "simulations/" + parts.SimID + "/README.md", SHA256: strings.Repeat("a", 64)},
			},
			SimulationTreeHash: strings.Repeat("b", 64),
		},
		Rubric: RubricInfo{
			RubricVersion: "rubric-20260519-v1",
			ScoreScale:    "0..5",
			ScoreMeanings: map[string]string{"0": "no fit", "5": "strong fit"},
			Axes:          []string{"scenario_fit", "risk_penalty"},
		},
		Scores: FitnessScores{
			ScenarioFit:                4,
			PromiseGridAlignment:       4,
			Auditability:               4,
			EvolutionSafety:            3,
			LayerBoundaryClarity:       4,
			FailureHandling:            3,
			ImplementationPlausibility: 4,
			RiskPenalty:                1,
		},
		Fitness: FitnessSummary{
			Raw:              25,
			Normalized0To100: 78,
			Confidence0To1:   0.72,
		},
		Assessment: Assessment{
			Rationale:         "Strong enough for validation fixture.",
			Strengths:         []string{"clear fit"},
			Weaknesses:        []string{"limited fixture depth"},
			Risks:             []string{"test-only data"},
			OpenQuestions:     []string{"none for fixture"},
			AuthorityBoundary: "Evidence only.",
		},
	}
}

func mustJSON(t *testing.T, result FitnessResult) []byte {
	t.Helper()
	bytes, err := jsonMarshalIndent(result)
	if err != nil {
		t.Fatalf("marshal fixture: %v", err)
	}
	return bytes
}

func jsonMarshalIndent(result FitnessResult) ([]byte, error) {
	return json.MarshalIndent(result, "", "  ")
}

func validScorePayloadJSON() string {
	return `{
  "scores": {
    "scenario_fit": 4,
    "promisegrid_alignment": 4,
    "auditability": 4,
    "evolution_safety": 3,
    "layer_boundary_clarity": 4,
    "failure_handling": 3,
    "implementation_plausibility": 4,
    "risk_penalty": 1
  },
  "fitness": {
    "raw": 25,
    "normalized_0_100": 78,
    "confidence_0_1": 0.72
  },
  "assessment": {
    "rationale": "Strong enough for a provider-backed fixture.",
    "strengths": ["clear audit path"],
    "weaknesses": ["limited scenario depth"],
    "risks": ["fixture-only evaluation"],
    "open_questions": ["none for this test"],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}`
}

func validChildBundleJSON(childID string) string {
	return fmt.Sprintf(`{"child_id":%q,"design_delta_summary":"Tighten audit language while keeping the fixture small.","files":[{"path":"README.md","content":"# Generated Child\n\nThis generated child keeps audit evidence local.\n"},{"path":"QUESTION.md","content":"# Questions\n\nDoes this child improve the audit path?\n"}]}`, childID)
}

func hasIssue(issues []string, want string) bool {
	for _, issue := range issues {
		if issue == want {
			return true
		}
	}
	return false
}
