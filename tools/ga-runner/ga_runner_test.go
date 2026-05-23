package main

import (
	"bytes"
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
	if !hasIssue(issues, expectedResultSchemaMessage()) {
		t.Fatalf("expected schema issue, got %v", issues)
	}
}

func TestValidateResultFileRequiresV2Fields(t *testing.T) {
	repo := newTestRepo(t)
	path := repo.Path("results", "SIM-alpha", "scenario-one", "model-a", "20260519-101500.json")
	raw := `{
  "schema": "promisegrid.ga.result.v2",
  "result_id": "SIM-alpha-scenario-one-model-a-20260519-101500",
  "run_group_id": "ga-test",
  "cell_id": "ga-test-000001",
  "sim_id": "SIM-alpha",
  "scenario_id": "scenario-one",
  "model_id": "model-a",
  "timestamp_utc": "20260519-101500",
  "result_path": "results/SIM-alpha/scenario-one/model-a/20260519-101500.json",
  "runner": {"tool": "ga-runner"},
  "source": {
    "repo_commit": "abc123",
    "sim_path": "simulations/SIM-alpha/",
    "scenario_path": "scenarios/scenario-one/scenario-one.md",
    "root_contract_paths": ["results/RUN-PROTOCOL.md", "scenarios/README.md"],
    "files": [{"path": "simulations/SIM-alpha/README.md", "sha256": "` + strings.Repeat("a", 64) + `"}],
    "simulation_tree_hash": "` + strings.Repeat("b", 64) + `"
  },
  "rubric": {
    "rubric_version": "ga-rubric-20260522-v2",
    "score_scale": "0..5",
    "score_meanings": {"0": "no fit", "5": "strong fit"},
    "axes": ["scenario_fit", "promisegrid_alignment", "auditability", "evolution_safety", "layer_boundary_clarity", "failure_handling", "implementation_plausibility", "promise_vocabulary", "simplicity_durability", "risk_penalty"]
  },
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
  "fitness": {"raw": 0, "normalized_0_100": 0, "confidence_0_1": 0.7},
  "assessment": {"rationale": "fixture", "strengths": [], "weaknesses": [], "risks": [], "open_questions": [], "authority_boundary": "Evidence only; does not settle PromiseGrid design."}
}`
	if err := ensureParent(path); err != nil {
		t.Fatalf("make parent: %v", err)
	}
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatalf("write result: %v", err)
	}
	issues := validateResultFile(repo, path)
	if !hasIssue(issues, "scores.promise_vocabulary is required for "+resultSchemaV2) {
		t.Fatalf("expected promise_vocabulary issue, got %v", issues)
	}
	if !hasIssue(issues, "scores.simplicity_durability is required for "+resultSchemaV2) {
		t.Fatalf("expected simplicity_durability issue, got %v", issues)
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

func TestDiscoverTrackedPopulationReadsSimulationMetadata(t *testing.T) {
	repo := newGitTestRepo(t)
	writeTestFile(t, repo.Path("simulations", "SIM-parent", "README.md"), "# Parent\n")
	writeTestFile(t, repo.Path("simulations", "SIM-parent", "SIM-META.json"), "{\n  \"schema\": \"promisegrid.sim.meta.v1\",\n  \"role\": \"negative-control\"\n}\n")
	gitAdd(t, repo, "simulations/SIM-parent/README.md", "simulations/SIM-parent/SIM-META.json")

	population, err := discoverTrackedPopulation(repo)
	if err != nil {
		t.Fatalf("discover tracked population with metadata: %v", err)
	}
	if len(population) != 1 || population[0].Role != simRoleNegativeCtl {
		t.Fatalf("expected negative-control role in tracked population, got %#v", population)
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

func TestBuildGenerationPlanIncludesRequestedSimsAndScenarios(t *testing.T) {
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
	plan, err := buildGenerationPlan(population, scenarios, PlanOptions{
		RunGroupID:         "ga-include",
		ModelID:            "model-a",
		ShuffleSeed:        "5",
		ParentCount:        2,
		ScenarioCount:      1,
		ChildCount:         1,
		MaxPromotions:      1,
		IncludeSimIDs:      []string{"SIM-c", "SIM-c"},
		IncludeScenarioIDs: []string{"scenario-c", "scenario-b"},
	})
	if err != nil {
		t.Fatalf("build include plan: %v", err)
	}
	if !stringSliceContains(parentIDs(plan.Parents), "SIM-c") {
		t.Fatalf("included sim was not selected: %#v", plan.Parents)
	}
	if got := strings.Join(scenarioIDs(plan.Scenarios), ","); got != "scenario-c,scenario-b" {
		t.Fatalf("included scenarios should be preserved and expand sample when needed, got %s", got)
	}

	_, err = buildGenerationPlan(population, scenarios, PlanOptions{
		RunGroupID:    "ga-include-missing",
		ModelID:       "model-a",
		ParentCount:   2,
		ScenarioCount: 1,
		ChildCount:    1,
		MaxPromotions: 1,
		IncludeSimIDs: []string{"SIM-missing"},
	})
	if err == nil || !strings.Contains(err.Error(), `included simulation "SIM-missing" was not discovered`) {
		t.Fatalf("expected missing included sim error, got %v", err)
	}
}

func TestBuildGenerationPlanExcludesNegativeControlParentsUnlessIncluded(t *testing.T) {
	population := []PopulationSim{
		{SimID: "SIM-a", Path: "simulations/SIM-a/", TreeHash: "a"},
		{SimID: "SIM-b", Path: "simulations/SIM-b/", TreeHash: "b", Role: simRoleNegativeCtl},
		{SimID: "SIM-c", Path: "simulations/SIM-c/", TreeHash: "c"},
	}
	scenarios := []Scenario{
		{ScenarioID: "scenario-a", Path: "scenarios/scenario-a/scenario-a.md"},
		{ScenarioID: "scenario-b", Path: "scenarios/scenario-b/scenario-b.md"},
	}

	plan, err := buildGenerationPlan(population, scenarios, PlanOptions{
		RunGroupID:    "ga-negative-default",
		ModelID:       "model-a",
		ShuffleSeed:   "7",
		ParentCount:   2,
		ScenarioCount: 1,
		ChildCount:    1,
		MaxPromotions: 1,
	})
	if err != nil {
		t.Fatalf("build default plan with negative control: %v", err)
	}
	if stringSliceContains(parentIDs(plan.Parents), "SIM-b") {
		t.Fatalf("negative-control sim should be excluded by default: %#v", plan.Parents)
	}

	plan, err = buildGenerationPlan(population, scenarios, PlanOptions{
		RunGroupID:    "ga-negative-include",
		ModelID:       "model-a",
		ShuffleSeed:   "7",
		ParentCount:   2,
		ScenarioCount: 1,
		ChildCount:    1,
		MaxPromotions: 1,
		IncludeSimIDs: []string{"SIM-b"},
	})
	if err != nil {
		t.Fatalf("build include plan with negative control: %v", err)
	}
	if !stringSliceContains(parentIDs(plan.Parents), "SIM-b") || !plan.Parents[0].ExplicitInclude {
		t.Fatalf("explicit include should preserve negative-control parent: %#v", plan.Parents)
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
	if result.Schema != resultSchemaV2 {
		t.Fatalf("score wrote schema %q, want %q", result.Schema, resultSchemaV2)
	}
	if result.Rubric.RubricVersion != rubricVersionV2 {
		t.Fatalf("score wrote rubric version %q, want %q", result.Rubric.RubricVersion, rubricVersionV2)
	}
	if result.Scores.PromiseVocabulary != 4 || result.Scores.SimplicityDurability != 4 {
		t.Fatalf("score did not capture new v2 axes: %#v", result.Scores)
	}
	if result.Fitness.Raw != 38 || result.Fitness.Normalized0To100 != 76 {
		t.Fatalf("score did not recompute deterministic fitness: %#v", result.Fitness)
	}
}

func TestRunScoreWritesZeroValuedV2Axes(t *testing.T) {
	repo := newGAFixtureRepo(t)
	initGAStateForTest(t, repo, "ga-score-zero-v2-axes")
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			return ProviderResponse{
				Text:        scorePayloadWithZeroV2AxesJSON(),
				RequestID:   "req-score-zero",
				ResponseID:  "resp-score-zero",
				ServiceTier: defaultServiceTier,
				UsageJSON:   `{"input_tokens":1000,"input_tokens_details":{"cached_tokens":100},"output_tokens":500}`,
			}, nil
		},
	}
	var out strings.Builder
	err := runScoreWithProvider(context.Background(), repo, provider, scoreOptions{
		RunGroupID:       "ga-score-zero-v2-axes",
		Target:           "parents",
		ProviderName:     "fake",
		APIModel:         "model-a",
		ReasoningEffort:  "medium",
		OutputContract:   outputContractJSONSchemaStrict,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
	}, &out)
	if err != nil {
		t.Fatalf("score with zero-valued v2 axes: %v\n%s", err, out.String())
	}
	state := mustReadGAState(t, repo, "ga-score-zero-v2-axes")
	resultPath := repo.Abs(state.Cells[0].ResultPath)
	bytes, err := os.ReadFile(resultPath)
	if err != nil {
		t.Fatalf("read result: %v", err)
	}
	for _, want := range []string{`"promise_vocabulary": 0`, `"simplicity_durability": 0`} {
		if !strings.Contains(string(bytes), want) {
			t.Fatalf("result omitted required zero-valued v2 axis %s:\n%s", want, string(bytes))
		}
	}
	if issues := validateResultFile(repo, resultPath); len(issues) != 0 {
		t.Fatalf("zero-valued v2 axes should validate as present scores, got %v", issues)
	}
}

func TestRunScoreUsesPerCellAPIModelWhenOptionIsEmpty(t *testing.T) {
	repo := newGAFixtureRepo(t)
	initGAStateForTest(t, repo, "ga-score-state-api-model")
	state := mustReadGAState(t, repo, "ga-score-state-api-model")
	parentIDs := map[string]bool{}
	for _, parent := range state.Parents {
		parentIDs[parent.SimID] = true
	}
	for index := range state.Cells {
		if parentIDs[state.Cells[index].SimID] {
			state.Cells[index].APIModel = "gpt-5.3-codex"
		}
	}
	writeTestState(t, repo, "ga-score-state-api-model", state)
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			if request.APIModel != "gpt-5.3-codex" {
				t.Fatalf("score request api model = %q, want %q", request.APIModel, "gpt-5.3-codex")
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
		RunGroupID:       "ga-score-state-api-model",
		Target:           "parents",
		ProviderName:     "fake",
		ReasoningEffort:  "xhigh",
		MaxOutputTokens:  4000,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
	}, &out)
	if err != nil {
		t.Fatalf("score: %v\n%s", err, out.String())
	}
}

func TestAuditSimulationVocabularyClassifiesAllowedAndHardHitSims(t *testing.T) {
	repo := newGAFixtureRepo(t)
	writeTestFile(t, repo.Path("simulations", "SIM-hard-hit", "README.md"), "# Hard Hit\n\nThis sim uses a claim card artifact.\n")
	writeTestFile(t, repo.Path("simulations", "SIM-udp-feed-v0-conformance", "README.md"), "# UDP Fixture\n\nA UDP-feed v0 conformance fixture.\n")
	gitAdd(t, repo,
		"simulations/SIM-hard-hit/README.md",
		"simulations/SIM-udp-feed-v0-conformance/README.md",
	)
	status, _, err := auditSimulationVocabulary(repo, auditSourceState{
		Mode:          auditSourceResolutionHistorical,
		ActiveSimPath: "simulations/SIM-hard-hit",
	}, "SIM-hard-hit")
	if err != nil {
		t.Fatalf("audit hard-hit sim: %v", err)
	}
	if status != "hard_hit" {
		t.Fatalf("hard-hit sim status = %q, want %q", status, "hard_hit")
	}
	status, _, err = auditSimulationVocabulary(repo, auditSourceState{
		Mode:          auditSourceResolutionHistorical,
		ActiveSimPath: "simulations/SIM-udp-feed-v0-conformance",
	}, "SIM-udp-feed-v0-conformance")
	if err != nil {
		t.Fatalf("audit udp sim: %v", err)
	}
	if status != "clean" {
		t.Fatalf("udp fixture status = %q, want %q", status, "clean")
	}
}

func TestAuditCanonicalV1ResultsUsesCanonicalFallbackForPromotedResults(t *testing.T) {
	repo := newGAFixtureRepo(t)
	writeTestFile(t, repo.Path("simulations", "SIM-promoted-clean", "README.md"), "# Promoted Clean\n\nAlice promises this payload meets the protocol specification referred to by this pCID.\n")
	writeTestFile(t, repo.Path("simulations", "SIM-promoted-clean", "QUESTION.md"), "# Questions\n\nCan Bob keep the outer grid minimal?\n")
	gitAdd(t, repo,
		"simulations/SIM-promoted-clean/README.md",
		"simulations/SIM-promoted-clean/QUESTION.md",
	)
	writePromotedCanonicalV1Result(t, repo, "SIM-promoted-clean", "SIM-promoted-child", "scenario-one", "openai-gpt-5.4-xhigh", "20260519-101520", 88)
	records, err := auditCanonicalV1Results(repo, "", "")
	if err != nil {
		t.Fatalf("audit canonical promoted result: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("record count = %d, want 1", len(records))
	}
	record := records[0]
	if record.SourceResolution != auditSourceResolutionCanonicalFallback {
		t.Fatalf("source resolution = %q, want %q", record.SourceResolution, auditSourceResolutionCanonicalFallback)
	}
	if !record.ExactMatch {
		t.Fatalf("expected promoted canonical fallback to remain exact-match")
	}
	if record.VocabularyStatus != "clean" {
		t.Fatalf("vocabulary status = %q, want clean", record.VocabularyStatus)
	}
}

func TestAuditCanonicalV1ResultsDetectsCanonicalFallbackDrift(t *testing.T) {
	repo := newGAFixtureRepo(t)
	writeTestFile(t, repo.Path("simulations", "SIM-promoted-drift", "README.md"), "# Promoted Drift\n\nAlice promises this payload meets the protocol specification referred to by this pCID.\n")
	writeTestFile(t, repo.Path("simulations", "SIM-promoted-drift", "QUESTION.md"), "# Questions\n\nCan Bob detect drift?\n")
	gitAdd(t, repo,
		"simulations/SIM-promoted-drift/README.md",
		"simulations/SIM-promoted-drift/QUESTION.md",
	)
	writePromotedCanonicalV1Result(t, repo, "SIM-promoted-drift", "SIM-promoted-drift-child", "scenario-one", "openai-gpt-5.4-xhigh", "20260519-101521", 84)
	writeTestFile(t, repo.Path("simulations", "SIM-promoted-drift", "README.md"), "# Promoted Drift\n\nThis drift changes the canonical bytes after scoring.\n")
	records, err := auditCanonicalV1Results(repo, "", "")
	if err != nil {
		t.Fatalf("audit drifted promoted result: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("record count = %d, want 1", len(records))
	}
	record := records[0]
	if record.SourceResolution != auditSourceResolutionCanonicalFallback {
		t.Fatalf("source resolution = %q, want %q", record.SourceResolution, auditSourceResolutionCanonicalFallback)
	}
	if record.ExactMatch {
		t.Fatalf("drifted canonical fallback remained exact-match")
	}
}

func TestRunAuditReportsSourceResolutionCounts(t *testing.T) {
	repo := newGAFixtureRepo(t)
	writeTestFile(t, repo.Path("simulations", "SIM-promoted-audit-grid-envelope", "README.md"), "# Promoted Audit\n\nAlice promises this payload meets the protocol specification referred to by this pCID.\n")
	gitAdd(t, repo, "simulations/SIM-promoted-audit-grid-envelope/README.md")
	writePromotedCanonicalV1Result(t, repo, "SIM-promoted-audit-grid-envelope", "SIM-promoted-audit-child", "scenario-one", "openai-gpt-5.4-xhigh", "20260519-101522", 87)
	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "audit",
		"-repo-root", repo.Root,
	}, &out, &out)
	if err != nil {
		t.Fatalf("audit command: %v\n%s", err, out.String())
	}
	text := out.String()
	if !strings.Contains(text, "source_canonical_fallback=1") {
		t.Fatalf("audit output missing canonical fallback count:\n%s", text)
	}
	if !strings.Contains(text, "source_resolution=canonical_fallback") {
		t.Fatalf("audit output missing source resolution summary:\n%s", text)
	}
}

func TestRunBackfillInitUsesPromotedCanonicalFallbackResults(t *testing.T) {
	repo := newGAFixtureRepo(t)
	writeTestFile(t, repo.Path("simulations", "SIM-promoted-hard-hit", "README.md"), "# Promoted Hard Hit\n\nThis sim still uses a claim card artifact.\n")
	writeTestFile(t, repo.Path("simulations", "SIM-promoted-hard-hit", "QUESTION.md"), "# Questions\n\nCan Alice simplify this later?\n")
	gitAdd(t, repo,
		"simulations/SIM-promoted-hard-hit/README.md",
		"simulations/SIM-promoted-hard-hit/QUESTION.md",
	)
	writePromotedCanonicalV1Result(t, repo, "SIM-promoted-hard-hit", "SIM-promoted-hard-hit-child", "scenario-one", "openai-gpt-5.4-xhigh", "20260519-101523", 83)
	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "backfill-init",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-backfill-promoted",
		"-timestamp", "20260522-130500",
	}, &out, &out)
	if err != nil {
		t.Fatalf("backfill-init promoted result: %v\n%s", err, out.String())
	}
	state := mustReadGAState(t, repo, "ga-backfill-promoted")
	if len(state.Parents) != 1 || len(state.Cells) != 1 {
		t.Fatalf("unexpected promoted backfill state counts: parents=%d cells=%d", len(state.Parents), len(state.Cells))
	}
	if state.Parents[0].SimID != "SIM-promoted-hard-hit" {
		t.Fatalf("backfill parent sim = %q, want %q", state.Parents[0].SimID, "SIM-promoted-hard-hit")
	}
}

func TestRunBackfillInitWritesTargetedState(t *testing.T) {
	repo := newGAFixtureRepo(t)
	writeTestFile(t, repo.Path("simulations", "SIM-claim-card-target", "README.md"), "# Claim Card\n\nThis sim still talks about a claim card artifact.\n")
	writeTestFile(t, repo.Path("simulations", "SIM-claim-card-target", "QUESTION.md"), "# Questions\n\nCan Alice and Bob simplify this?\n")
	writeTestFile(t, repo.Path("simulations", "SIM-grid-envelope-clean", "README.md"), "# Grid Envelope Clean\n\nAlice promises this payload meets the protocol specification referred to by this pCID.\n")
	writeTestFile(t, repo.Path("simulations", "SIM-grid-envelope-clean", "QUESTION.md"), "# Questions\n\nCan Bob keep this envelope minimal and durable?\n")
	gitAdd(t, repo,
		"simulations/SIM-claim-card-target/README.md",
		"simulations/SIM-claim-card-target/QUESTION.md",
		"simulations/SIM-grid-envelope-clean/README.md",
		"simulations/SIM-grid-envelope-clean/QUESTION.md",
	)
	writeExactMatchV1Result(t, repo, "SIM-claim-card-target", "scenario-one", "openai-gpt-5.4-xhigh", "20260519-101500", 82)
	writeExactMatchV1Result(t, repo, "SIM-grid-envelope-clean", "scenario-one", "openai-gpt-5.4-xhigh", "20260519-101501", 91)
	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "backfill-init",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-backfill",
		"-timestamp", "20260522-130000",
		"-clean-envelope-count", "1",
	}, &out, &out)
	if err != nil {
		t.Fatalf("backfill-init: %v\n%s", err, out.String())
	}
	state := mustReadGAState(t, repo, "ga-backfill")
	if len(state.Parents) != 2 || len(state.Cells) != 2 || len(state.Children) != 0 {
		t.Fatalf("unexpected backfill state counts: parents=%d cells=%d children=%d", len(state.Parents), len(state.Cells), len(state.Children))
	}
	if state.Cells[0].APIModel != "gpt-5.4" {
		t.Fatalf("backfill cell api_model = %q, want %q", state.Cells[0].APIModel, "gpt-5.4")
	}
	if !strings.HasSuffix(state.Cells[0].ResultPath, "20260522-130000.json") {
		t.Fatalf("backfill result path = %q, want new timestamp", state.Cells[0].ResultPath)
	}
}

func TestRunBackfillInitAppliesStagedModelAndEffortOverrides(t *testing.T) {
	repo := newGAFixtureRepo(t)
	writeTestFile(t, repo.Path("simulations", "SIM-claim-card-target", "README.md"), "# Claim Card\n\nThis sim still talks about a claim card artifact.\n")
	writeTestFile(t, repo.Path("simulations", "SIM-claim-card-target", "QUESTION.md"), "# Questions\n\nCan Alice and Bob simplify this?\n")
	gitAdd(t, repo,
		"simulations/SIM-claim-card-target/README.md",
		"simulations/SIM-claim-card-target/QUESTION.md",
	)
	writeExactMatchV1Result(t, repo, "SIM-claim-card-target", "scenario-one", "openai-gpt-5.4-xhigh", "20260519-101500", 82)
	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "backfill-init",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-backfill-high",
		"-timestamp", "20260522-131000",
		"-staged-model-id", "openai-gpt-5.4-high",
		"-staged-reasoning-effort", "high",
	}, &out, &out)
	if err != nil {
		t.Fatalf("backfill-init staged override: %v\n%s", err, out.String())
	}
	state := mustReadGAState(t, repo, "ga-backfill-high")
	if state.ModelID != "openai-gpt-5.4-high" {
		t.Fatalf("backfill state model_id = %q, want %q", state.ModelID, "openai-gpt-5.4-high")
	}
	if len(state.Cells) != 1 {
		t.Fatalf("unexpected staged backfill cell count: %d", len(state.Cells))
	}
	cell := state.Cells[0]
	if cell.ModelID != "openai-gpt-5.4-high" {
		t.Fatalf("backfill cell model_id = %q, want %q", cell.ModelID, "openai-gpt-5.4-high")
	}
	if cell.APIModel != "gpt-5.4" {
		t.Fatalf("backfill cell api_model = %q, want %q", cell.APIModel, "gpt-5.4")
	}
	if cell.ReasoningEffort != "high" {
		t.Fatalf("backfill cell reasoning_effort = %q, want %q", cell.ReasoningEffort, "high")
	}
	if !strings.HasSuffix(cell.ResultPath, "/openai-gpt-5.4-high/20260522-131000.json") {
		t.Fatalf("backfill result path = %q, want staged model/timestamp suffix", cell.ResultPath)
	}
}

func TestRunBackfillInitDedupesRepeatedHistoricalPairs(t *testing.T) {
	repo := newGAFixtureRepo(t)
	writeTestFile(t, repo.Path("simulations", "SIM-claim-card-target", "README.md"), "# Claim Card\n\nThis sim still talks about a claim card artifact.\n")
	writeTestFile(t, repo.Path("simulations", "SIM-claim-card-target", "QUESTION.md"), "# Questions\n\nCan Alice and Bob simplify this?\n")
	gitAdd(t, repo,
		"simulations/SIM-claim-card-target/README.md",
		"simulations/SIM-claim-card-target/QUESTION.md",
	)
	writeExactMatchV1Result(t, repo, "SIM-claim-card-target", "scenario-one", "openai-gpt-5.4-xhigh", "20260519-101500", 82)
	writeExactMatchV1Result(t, repo, "SIM-claim-card-target", "scenario-one", "openai-gpt-5.4-low", "20260519-101700", 80)
	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "backfill-init",
		"-repo-root", repo.Root,
		"-run-group-id", "ga-backfill-dedupe",
		"-timestamp", "20260522-132000",
	}, &out, &out)
	if err != nil {
		t.Fatalf("backfill-init dedupe: %v\n%s", err, out.String())
	}
	state := mustReadGAState(t, repo, "ga-backfill-dedupe")
	if len(state.Parents) != 1 || len(state.Cells) != 1 {
		t.Fatalf("unexpected deduped backfill counts: parents=%d cells=%d", len(state.Parents), len(state.Cells))
	}
	cell := state.Cells[0]
	if cell.SimID != "SIM-claim-card-target" || cell.ScenarioID != "scenario-one" {
		t.Fatalf("unexpected deduped cell identity: %+v", cell)
	}
	if cell.APIModel != "gpt-5.4" {
		t.Fatalf("deduped backfill chose api_model %q, want %q", cell.APIModel, "gpt-5.4")
	}
}

func TestCompareBackfillWritesReport(t *testing.T) {
	repo := newGitTestRepo(t)
	writeTestFile(t, repo.Path("results", "RUN-PROTOCOL.md"), "# Run Protocol\n\nScore cells as evidence.\n")
	writeTestFile(t, repo.Path("scenarios", "README.md"), "# Scenarios\n\nUse 100-year PromiseGrid pressure.\n")
	writeTestFile(t, repo.Path("scenarios", "scenario-one", "scenario-one.md"), "# Scenario One\n\nAlice asks Bob to ship labels with auditable promises.\n")
	writeTestFile(t, repo.Path("simulations", "SIM-grid-envelope-clean", "README.md"), "# Grid Envelope\n\nAlice promises this payload meets the protocol specification referred to by this pCID.\n")
	writeTestFile(t, repo.Path("simulations", "SIM-grid-envelope-clean", "QUESTION.md"), "# Questions\n\nCan Bob audit the payload promise?\n")
	writeTestFile(t, repo.Path("simulations", "SIM-robot-app-semantics-conformance", "README.md"), "# App Semantics\n\nAlice promises this payload meets the protocol specification referred to by this pCID.\n")
	writeTestFile(t, repo.Path("simulations", "SIM-robot-app-semantics-conformance", "QUESTION.md"), "# Questions\n\nCan Carol audit partial implementation promises?\n")
	gitAdd(t, repo,
		"results/RUN-PROTOCOL.md",
		"scenarios/README.md",
		"scenarios/scenario-one/scenario-one.md",
		"simulations/SIM-grid-envelope-clean/README.md",
		"simulations/SIM-grid-envelope-clean/QUESTION.md",
		"simulations/SIM-robot-app-semantics-conformance/README.md",
		"simulations/SIM-robot-app-semantics-conformance/QUESTION.md",
	)
	writeExactMatchV1Result(t, repo, "SIM-grid-envelope-clean", "scenario-one", "openai-gpt-5.4-xhigh", "20260519-101500", 91)
	writeExactMatchV1Result(t, repo, "SIM-robot-app-semantics-conformance", "scenario-one", "openai-gpt-5.4-xhigh", "20260519-101501", 62)
	writeExactMatchV1Result(t, repo, "SIM-robot-app-semantics-conformance", "scenario-one", "openai-gpt-5.4-xhigh", "20260519-101502", 63)

	v2APath := repo.Path("results", "SIM-grid-envelope-clean", "scenario-one", "openai-gpt-5.4-medium", "20260523-220000.json")
	v2A := validResult(repo, v2APath)
	v2A.Schema = resultSchemaV2
	v2A.RunGroupID = "ga-backfill-medium"
	v2A.CellID = "ga-backfill-medium-000001"
	v2A.SimID = "SIM-grid-envelope-clean"
	v2A.ScenarioID = "scenario-one"
	v2A.ModelID = "openai-gpt-5.4-medium"
	v2A.TimestampUTC = "20260523-220000"
	v2A.ResultPath = repo.Rel(v2APath)
	v2A.Runner.APIModel = "gpt-5.4"
	v2A.Runner.ReasoningEffort = "medium"
	v2A.Source.RepoCommit = repo.GitCommit()
	v2A.Source.SimPath = "simulations/SIM-grid-envelope-clean/"
	v2A.Source.ScenarioPath = "scenarios/scenario-one/scenario-one.md"
	v2A.Source.RootContractPaths = []string{"results/RUN-PROTOCOL.md", "scenarios/README.md"}
	v2A.Source.Files = buildExactMatchSourceFiles(t, repo, "SIM-grid-envelope-clean", "scenario-one")
	v2A.Source.SimulationTreeHash = mustSimulationTreeHash(t, repo, "SIM-grid-envelope-clean")
	v2A.Rubric.RubricVersion = rubricVersionV2
	v2A.Rubric.Axes = rubricAxesV2
	v2A.Scores.PromiseVocabulary = 5
	v2A.Scores.SimplicityDurability = 4
	v2A.Fitness.Normalized0To100 = 88
	if err := writeFitnessResultAtomic(v2APath, v2A); err != nil {
		t.Fatalf("write v2 result A: %v", err)
	}

	v2BPath := repo.Path("results", "SIM-robot-app-semantics-conformance", "scenario-one", "openai-gpt-5.4-medium", "20260523-220000.json")
	v2B := v2A
	v2B.ResultID = "SIM-robot-app-semantics-conformance-scenario-one-openai-gpt-5.4-medium-20260523-220000"
	v2B.CellID = "ga-backfill-medium-000002"
	v2B.SimID = "SIM-robot-app-semantics-conformance"
	v2B.ResultPath = repo.Rel(v2BPath)
	v2B.Source.SimPath = "simulations/SIM-robot-app-semantics-conformance/"
	v2B.Source.Files = buildExactMatchSourceFiles(t, repo, "SIM-robot-app-semantics-conformance", "scenario-one")
	v2B.Source.SimulationTreeHash = mustSimulationTreeHash(t, repo, "SIM-robot-app-semantics-conformance")
	v2B.Fitness.Normalized0To100 = 70
	if err := writeFitnessResultAtomic(v2BPath, v2B); err != nil {
		t.Fatalf("write v2 result B: %v", err)
	}

	state := GAState{
		Schema:     stateSchemaV1,
		RunGroupID: "ga-backfill-medium",
		Cells: []GACell{
			{CellID: "ga-backfill-medium-000001", SimID: "SIM-grid-envelope-clean", ScenarioID: "scenario-one", ModelID: "openai-gpt-5.4-medium", ResultPath: repo.Rel(v2APath), Status: "done"},
			{CellID: "ga-backfill-medium-000002", SimID: "SIM-robot-app-semantics-conformance", ScenarioID: "scenario-one", ModelID: "openai-gpt-5.4-medium", ResultPath: repo.Rel(v2BPath), Status: "done"},
		},
	}
	writeTestState(t, repo, "ga-backfill-medium", state)

	var out bytes.Buffer
	err := runMain([]string{"ga-runner", "compare-backfill", "-repo-root", repo.Root, "-run-group-id", "ga-backfill-medium"}, &out, io.Discard)
	if err != nil {
		t.Fatalf("compare-backfill: %v\n%s", err, out.String())
	}
	if !strings.Contains(out.String(), "report=results/reports/ga-backfill-medium-comparison.md") {
		t.Fatalf("compare output missing report path:\n%s", out.String())
	}
	if !strings.Contains(out.String(), "rerun_groups=1") || !strings.Contains(out.String(), "ambiguous=0") {
		t.Fatalf("compare output missing rerun/ambiguity counts:\n%s", out.String())
	}
	reportPath := repo.Path("results", "reports", "ga-backfill-medium-comparison.md")
	reportBytes, err := os.ReadFile(reportPath)
	if err != nil {
		t.Fatalf("read report: %v", err)
	}
	report := string(reportBytes)
	if !strings.Contains(report, "# GA Backfill Comparison Report") {
		t.Fatalf("missing report title:\n%s", report)
	}
	if !strings.Contains(report, "`SIM-grid-envelope-clean`") || !strings.Contains(report, "`SIM-robot-app-semantics-conformance`") {
		t.Fatalf("missing sim detail in report:\n%s", report)
	}
	if !strings.Contains(report, "## Family Highlights") || !strings.Contains(report, "grid-envelope") || !strings.Contains(report, "conformance-family") {
		t.Fatalf("missing family highlights:\n%s", report)
	}
	if !strings.Contains(report, "- Historical rerun groups collapsed: `1`") || !strings.Contains(report, "- Ambiguous matched pairs: `0`") {
		t.Fatalf("missing rerun/ambiguity summary:\n%s", report)
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

func TestBuildScorePromptIncludesRequiredAxisChecklist(t *testing.T) {
	repo := newGAFixtureRepo(t)
	initGAStateForTest(t, repo, "ga-score-prompt-checklist")
	state := mustReadGAState(t, repo, "ga-score-prompt-checklist")
	scenario, err := scenarioFromState(state, state.Cells[0].ScenarioID)
	if err != nil {
		t.Fatalf("scenario from state: %v", err)
	}
	prompt, err := buildScorePrompt(repo, state, state.Cells[0], scenario)
	if err != nil {
		t.Fatalf("build prompt: %v", err)
	}
	for _, want := range []string{
		"Required score-axis checklist:",
		"`promise_vocabulary`",
		"`simplicity_durability`",
		"A response missing any required `scores` axis is invalid.",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("prompt missing %q:\n%s", want, prompt)
		}
	}
}

func TestRunScoreRetriesMissingV2Axes(t *testing.T) {
	repo := newGAFixtureRepo(t)
	initGAStateForTest(t, repo, "ga-score-schema-retry")
	calls := 0
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			calls++
			switch calls {
			case 1:
				return ProviderResponse{
					Text:        invalidScorePayloadMissingV2AxesJSON(),
					RequestID:   "req-score-1",
					ResponseID:  "resp-score-1",
					ServiceTier: defaultServiceTier,
					UsageJSON:   `{"input_tokens":1000,"input_tokens_details":{"cached_tokens":100},"output_tokens":500}`,
				}, nil
			case 2:
				if !strings.Contains(request.Prompt, "## Schema Correction") || !strings.Contains(request.Prompt, "`promise_vocabulary`") || !strings.Contains(request.Prompt, "`simplicity_durability`") {
					return ProviderResponse{}, fmt.Errorf("schema correction prompt missing missing-axis guidance")
				}
				return ProviderResponse{
					Text:        validScorePayloadJSON(),
					RequestID:   "req-score-2",
					ResponseID:  "resp-score-2",
					ServiceTier: defaultServiceTier,
					UsageJSON:   `{"input_tokens":1200,"input_tokens_details":{"cached_tokens":200},"output_tokens":700}`,
				}, nil
			case 3:
				return ProviderResponse{
					Text:        validScorePayloadJSON(),
					RequestID:   "req-score-3",
					ResponseID:  "resp-score-3",
					ServiceTier: defaultServiceTier,
					UsageJSON:   `{"input_tokens":900,"input_tokens_details":{"cached_tokens":90},"output_tokens":450}`,
				}, nil
			default:
				return ProviderResponse{}, fmt.Errorf("unexpected extra score call %d", calls)
			}
		},
	}
	var out strings.Builder
	err := runScoreWithProvider(context.Background(), repo, provider, scoreOptions{
		RunGroupID:       "ga-score-schema-retry",
		Target:           "parents",
		ProviderName:     "fake",
		APIModel:         "model-a",
		ReasoningEffort:  "medium",
		OutputContract:   outputContractPromptJSON,
		Workers:          1,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
	}, &out)
	if err != nil {
		t.Fatalf("score with schema retry: %v\n%s", err, out.String())
	}
	if calls != 3 {
		t.Fatalf("expected one schema-correction retry across two parent cells, got %d calls", calls)
	}
	state := mustReadGAState(t, repo, "ga-score-schema-retry")
	cell := state.Cells[0]
	if cell.Status != "done" || cell.Attempts != 2 {
		t.Fatalf("score state did not record retry success: %#v", cell)
	}
	if cell.InputTokens != 2200 || cell.CachedTokens != 300 || cell.OutputTokens != 1200 {
		t.Fatalf("score state did not accumulate retry usage: %#v", cell)
	}
	if cell.CostUSD <= 0 {
		t.Fatalf("score state missing accumulated retry cost: %#v", cell)
	}
	if cell.RequestID != "req-score-2" || cell.ResponseID != "resp-score-2" {
		t.Fatalf("score state did not keep latest provider ids: %#v", cell)
	}
	if state.Cells[1].Status != "done" || state.Cells[1].Attempts != 1 {
		t.Fatalf("second parent cell should succeed without retry: %#v", state.Cells[1])
	}
}

func TestRunScoreFailsAfterSchemaRetryStillMissingAxes(t *testing.T) {
	repo := newGAFixtureRepo(t)
	initGAStateForTest(t, repo, "ga-score-schema-retry-fail")
	calls := 0
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			calls++
			return ProviderResponse{
				Text:        invalidScorePayloadMissingV2AxesJSON(),
				RequestID:   fmt.Sprintf("req-score-%d", calls),
				ResponseID:  fmt.Sprintf("resp-score-%d", calls),
				ServiceTier: defaultServiceTier,
				UsageJSON:   `{"input_tokens":1000,"input_tokens_details":{"cached_tokens":100},"output_tokens":500}`,
			}, nil
		},
	}
	var out strings.Builder
	err := runScoreWithProvider(context.Background(), repo, provider, scoreOptions{
		RunGroupID:       "ga-score-schema-retry-fail",
		Target:           "parents",
		ProviderName:     "fake",
		APIModel:         "model-a",
		ReasoningEffort:  "medium",
		OutputContract:   outputContractPromptJSON,
		Workers:          1,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
	}, &out)
	if err == nil {
		t.Fatalf("expected schema retry failure")
	}
	if calls != 4 {
		t.Fatalf("expected both parent cells to retry before failure, got %d calls", calls)
	}
	state := mustReadGAState(t, repo, "ga-score-schema-retry-fail")
	for _, cell := range state.Cells {
		if strings.HasPrefix(cell.ResultPath, "proposals/") {
			continue
		}
		if cell.Status != "failed" || !strings.Contains(cell.ValidationMessage, "schema-correction retry still missing required score axes") {
			t.Fatalf("score state did not preserve schema retry failure: %#v", cell)
		}
	}
}

func TestRunScoreStrictStructuredOutputDoesNotSchemaRetry(t *testing.T) {
	repo := newGAFixtureRepo(t)
	initGAStateForTest(t, repo, "ga-score-structured-no-retry")
	calls := 0
	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			calls++
			if request.OutputContract != outputContractJSONSchemaStrict {
				return ProviderResponse{}, fmt.Errorf("output contract = %q, want %q", request.OutputContract, outputContractJSONSchemaStrict)
			}
			if request.OutputSchemaName != "ga_score_payload_v2" || len(request.OutputSchema) == 0 {
				return ProviderResponse{}, fmt.Errorf("structured output schema missing from request")
			}
			return ProviderResponse{
				Text:        invalidScorePayloadMissingV2AxesJSON(),
				RequestID:   fmt.Sprintf("req-score-%d", calls),
				ResponseID:  fmt.Sprintf("resp-score-%d", calls),
				ServiceTier: defaultServiceTier,
				UsageJSON:   `{"input_tokens":1000,"input_tokens_details":{"cached_tokens":100},"output_tokens":500}`,
			}, nil
		},
	}
	var out strings.Builder
	err := runScoreWithProvider(context.Background(), repo, provider, scoreOptions{
		RunGroupID:       "ga-score-structured-no-retry",
		Target:           "parents",
		ProviderName:     "fake",
		APIModel:         "model-a",
		ReasoningEffort:  "medium",
		OutputContract:   outputContractJSONSchemaStrict,
		Workers:          1,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
	}, &out)
	if err == nil {
		t.Fatalf("expected strict structured-output failure")
	}
	if calls != 2 {
		t.Fatalf("expected two total calls for two parent cells without schema retry, got %d", calls)
	}
	state := mustReadGAState(t, repo, "ga-score-structured-no-retry")
	for _, cell := range state.Cells {
		if strings.HasPrefix(cell.ResultPath, "proposals/") {
			continue
		}
		if cell.Attempts != 1 {
			t.Fatalf("strict structured-output cell should not schema-retry: %#v", cell)
		}
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
	generatedChildID := state.Children[0].ID()
	if generatedChildID == childID || !strings.Contains(generatedChildID, "-child-audit-path") {
		t.Fatalf("generate did not adopt descriptive child ID: planned=%s generated=%s", childID, generatedChildID)
	}
	if state.Children[0].ServiceTier != defaultServiceTier || state.Children[0].ServedServiceTier != defaultServiceTier {
		t.Fatalf("generate service tier not recorded: %#v", state.Children[0])
	}
	assertExists(t, repo.Abs(filepath.Join("proposals", "ga-generate", "simulations", generatedChildID, "README.md")))
	assertExists(t, repo.Abs(filepath.Join("proposals", "ga-generate", "simulations", generatedChildID, "QUESTION.md")))
	foundChildScoreCell := false
	for _, cell := range state.Cells {
		if cell.SimID != generatedChildID {
			continue
		}
		foundChildScoreCell = true
		if !strings.Contains(cell.ResultPath, "proposals/ga-generate/results/"+generatedChildID+"/") {
			t.Fatalf("child score cell was not repointed to descriptive result path: %#v", cell)
		}
	}
	if !foundChildScoreCell {
		t.Fatalf("no child score cell was repointed to %s", generatedChildID)
	}
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
			if !strings.Contains(request.Prompt, "## Compact Fitness Evidence From This Run") || !strings.Contains(request.Prompt, "normalized_0_100=") {
				return ProviderResponse{}, fmt.Errorf("generate prompt missing selected parent fitness evidence")
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
		if child.Operation != childOperationBreed || len(child.ParentIDs) != 2 {
			t.Fatalf("child did not breed two scored parents: %#v", child)
		}
		if child.ParentIDs[0] == child.ParentIDs[1] {
			t.Fatalf("breed reused the same parent twice: %#v", child)
		}
		for _, parentID := range child.ParentIDs {
			if parentID != "SIM-high" && parentID != "SIM-parent" && parentID != "SIM-low" {
				t.Fatalf("parent was not selected from scored parents: %#v", child)
			}
		}
	}
}

func TestRunGenerateExcludesNegativeControlParentsFromFitnessSelection(t *testing.T) {
	repo := newGAFixtureRepo(t)
	addTrackedFixtureSim(t, repo, "SIM-high")
	addTrackedFixtureSim(t, repo, "SIM-low")
	initGAStateForTestWithParentsAndChildren(t, repo, "ga-generate-negative-control", 3, 2)
	writeTestFile(t, repo.Path("simulations", "SIM-high", "SIM-META.json"), "{\n  \"schema\": \"promisegrid.sim.meta.v1\",\n  \"role\": \"negative-control\"\n}\n")
	gitAdd(t, repo, "simulations/SIM-high/SIM-META.json")
	writeParentFitnessResults(t, repo, "ga-generate-negative-control", map[string]float64{
		"SIM-high":   95,
		"SIM-parent": 60,
		"SIM-low":    20,
	})

	provider := fakeGAProvider{
		generate: func(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
			childID := promptChildID(t, request.Prompt)
			return ProviderResponse{
				Text:        validChildBundleJSON(childID),
				RequestID:   "req-negative-control",
				ResponseID:  "resp-negative-control",
				ServiceTier: defaultServiceTier,
				UsageJSON:   `{"input_tokens":1200,"input_tokens_details":{"cached_tokens":200},"output_tokens":700}`,
			}, nil
		},
	}
	var out strings.Builder
	err := runGenerateWithProvider(context.Background(), repo, provider, generateOptions{
		RunGroupID:       "ga-generate-negative-control",
		ProviderName:     "fake",
		APIModel:         "model-a",
		ReasoningEffort:  "medium",
		Workers:          1,
		InputPrice:       defaultInputUSDPerMTok,
		CachedInputPrice: defaultCachedInputUSDPerMTok,
		OutputPrice:      defaultOutputUSDPerMTok,
	}, &out)
	if err != nil {
		t.Fatalf("generate with negative-control fallback guard: %v\n%s", err, out.String())
	}

	state := mustReadGAState(t, repo, "ga-generate-negative-control")
	if state.Parents[0].SimID == "SIM-high" {
		t.Fatalf("negative-control parent should not lead the ranked parent pool: %#v", state.Parents)
	}
	for _, child := range state.Children {
		if stringSliceContains(child.ParentIDs, "SIM-high") {
			t.Fatalf("negative-control parent should not be selected for breed: %#v", child)
		}
	}
}

func TestEligibleRankedParentsAllowsExplicitNegativeControlOverride(t *testing.T) {
	repo := newGAFixtureRepo(t)
	addTrackedFixtureSim(t, repo, "SIM-high")
	addTrackedFixtureSim(t, repo, "SIM-low")
	ranked := []parentFitnessRank{
		{SimID: "SIM-high", AverageNormalized: 95, Samples: 3},
		{SimID: "SIM-parent", AverageNormalized: 60, Samples: 3},
		{SimID: "SIM-low", AverageNormalized: 20, Samples: 3},
	}
	state := GAState{
		Parents: []GAStateParent{
			{SimID: "SIM-high", Role: simRoleNegativeCtl, ExplicitInclude: true},
			{SimID: "SIM-parent"},
			{SimID: "SIM-low"},
		},
	}
	eligible, err := eligibleRankedParents(repo, state, ranked)
	if err != nil {
		t.Fatalf("eligible ranked parents: %v", err)
	}
	if len(eligible) != 3 || eligible[0].SimID != "SIM-high" {
		t.Fatalf("explicit include should preserve negative-control parent eligibility: %#v", eligible)
	}
}

func TestParentSelectionUsesWeightedHighPlusUniformRandomScoredParent(t *testing.T) {
	state := GAState{RunGroupID: "ga-select"}
	ranked := []parentFitnessRank{
		{SimID: "SIM-high", AverageNormalized: 95, Samples: 3},
		{SimID: "SIM-mid", AverageNormalized: 60, Samples: 3},
		{SimID: "SIM-low", AverageNormalized: 20, Samples: 3},
	}
	seenFirstParents := map[string]int{}
	seenSecondParents := map[string]bool{}
	for index := 0; index < 600; index++ {
		child := GAChild{ChildID: fmt.Sprintf("SIM-child-%02d", index)}
		parentIDs, ok := parentSelectionForChild(state, ranked, child, index)
		if !ok {
			t.Fatalf("expected parent selection for child %d", index)
		}
		if len(parentIDs) != 2 {
			t.Fatalf("expected two parents: %#v", parentIDs)
		}
		if parentIDs[0] == parentIDs[1] {
			t.Fatalf("random parent reused high parent: %#v", parentIDs)
		}
		seenFirstParents[parentIDs[0]]++
		seenSecondParents[parentIDs[1]] = true
	}
	if seenFirstParents["SIM-high"] <= seenFirstParents["SIM-mid"] || seenFirstParents["SIM-mid"] <= seenFirstParents["SIM-low"] {
		t.Fatalf("weighted high parent did not favor higher ranks: %#v", seenFirstParents)
	}
	if !seenSecondParents["SIM-high"] || !seenSecondParents["SIM-mid"] || !seenSecondParents["SIM-low"] {
		t.Fatalf("uniform random scored parent did not preserve full diversity after exclusion: %#v", seenSecondParents)
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
	if scoreDefaults.OutputContract != outputContractJSONSchemaStrict {
		t.Fatalf("score output contract = %q, want %q", scoreDefaults.OutputContract, outputContractJSONSchemaStrict)
	}
	if scoreDefaults.ReasoningEffort != defaultScoreReasoningEffort || scoreDefaults.TextVerbosity != defaultTextVerbosity || scoreDefaults.MaxOutputTokens != 0 || scoreDefaults.CostEstimateOutputTokens != defaultScoreCostEstimateOutputTokens {
		t.Fatalf("score request-shaping defaults not applied: %#v", scoreDefaults)
	}
	if scoreDefaults.ReasoningSummary != "" || scoreDefaults.StreamContentStdout {
		t.Fatalf("score stream-content defaults not applied: %#v", scoreDefaults)
	}
	if scoreDefaults.Workers != defaultScoreWorkers || scoreDefaults.RequestTimeout != defaultRequestTimeout || scoreDefaults.ProviderAttempts != defaultProviderMaxAttempts || scoreDefaults.ProviderElapsed != defaultProviderMaxElapsed || scoreDefaults.Stream != defaultProviderStream || scoreDefaults.StreamIdleTimeout != defaultStreamIdleTimeout {
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
	if generateDefaults.ReasoningSummary != "" || generateDefaults.StreamContentStdout {
		t.Fatalf("generate stream-content defaults not applied: %#v", generateDefaults)
	}
	if generateDefaults.Workers != defaultGenerateWorkers || generateDefaults.RequestTimeout != defaultRequestTimeout || generateDefaults.ProviderAttempts != defaultProviderMaxAttempts || generateDefaults.ProviderElapsed != defaultProviderMaxElapsed || generateDefaults.Stream != defaultProviderStream || generateDefaults.StreamIdleTimeout != defaultStreamIdleTimeout {
		t.Fatalf("generate throughput defaults not applied: %#v", generateDefaults)
	}
	scoreDefaultTier, err := parseScoreOptions([]string{"-run-group-id", "ga-score", "-api-model", "model-a", "-service-tier", "default"})
	if err != nil {
		t.Fatalf("parse explicit default tier: %v", err)
	}
	if scoreDefaultTier.ServiceTier != serviceTierDefault {
		t.Fatalf("score explicit service tier = %q, want %q", scoreDefaultTier.ServiceTier, serviceTierDefault)
	}
	scoreCustom, err := parseScoreOptions([]string{"-run-group-id", "ga-score", "-api-model", "model-a", "-workers", "3", "-request-timeout", "2m", "-provider-max-attempts", "4", "-provider-max-elapsed", "9m", "-stream=false", "-stream-idle-timeout", "11s", "-stream-content-stdout=true", "-reasoning-summary", "auto", "-text-verbosity", "high", "-max-output-tokens", "123", "-cost-estimate-output-tokens", "456"})
	if err != nil {
		t.Fatalf("parse custom throughput knobs: %v", err)
	}
	if scoreCustom.Workers != 3 || scoreCustom.RequestTimeout != 2*time.Minute || scoreCustom.ProviderAttempts != 4 || scoreCustom.ProviderElapsed != 9*time.Minute || scoreCustom.Stream || scoreCustom.StreamIdleTimeout != 11*time.Second {
		t.Fatalf("custom throughput knobs not applied: %#v", scoreCustom)
	}
	if scoreCustom.TextVerbosity != textVerbosityHigh || scoreCustom.MaxOutputTokens != 123 || scoreCustom.CostEstimateOutputTokens != 456 {
		t.Fatalf("custom request-shaping knobs not applied: %#v", scoreCustom)
	}
	if scoreCustom.ReasoningSummary != "auto" || !scoreCustom.StreamContentStdout {
		t.Fatalf("custom stream-content knobs not applied: %#v", scoreCustom)
	}
	if _, err := parseGenerateOptions([]string{"-run-group-id", "ga-generate", "-api-model", "model-a", "-service-tier", "priority"}); err == nil {
		t.Fatalf("expected priority service tier to be rejected")
	}
	if _, err := parseScoreOptions([]string{"-run-group-id", "ga-score", "-api-model", "model-a", "-output-contract", "yaml"}); err == nil {
		t.Fatalf("expected invalid output contract to be rejected")
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

func TestOpenAIProviderSendsStructuredOutputFormat(t *testing.T) {
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requestBytes, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}
		var body openAIRequest
		if err := json.Unmarshal(requestBytes, &body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if body.Text == nil || body.Text.Format == nil {
			t.Fatalf("structured output request missing text.format: %s", string(requestBytes))
		}
		if body.Text.Format.Type != "json_schema" || !body.Text.Format.Strict {
			t.Fatalf("unexpected text.format: %#v", body.Text.Format)
		}
		if body.Text.Format.Name != "ga_score_payload_v2" {
			t.Fatalf("structured output schema name = %q, want ga_score_payload_v2", body.Text.Format.Name)
		}
		if _, ok := body.Text.Format.Schema["properties"]; !ok {
			t.Fatalf("structured output schema missing properties: %#v", body.Text.Format.Schema)
		}
		return openAITestHTTPResponse(request, http.StatusOK, "req-schema", openAITestSuccessBody(serviceTierFlex)), nil
	})

	provider := OpenAIProvider{
		APIKey:  "test-key",
		BaseURL: "https://example.test/responses",
		Client:  &http.Client{Transport: transport},
		RetryPolicy: ProviderRetryPolicy{
			MaxAttempts: 1,
		},
	}
	_, err := provider.Generate(context.Background(), ProviderRequest{
		APIModel:         "model-a",
		ServiceTier:      serviceTierFlex,
		Prompt:           "score this",
		OutputContract:   outputContractJSONSchemaStrict,
		OutputSchemaName: "ga_score_payload_v2",
		OutputSchema:     scorePayloadJSONSchema(),
	})
	if err != nil {
		t.Fatalf("generate: %v", err)
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

func TestOpenAIProviderStreamsResponsesEvents(t *testing.T) {
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		requestBytes, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("read stream request body: %v", err)
		}
		var body openAIRequest
		if err := json.Unmarshal(requestBytes, &body); err != nil {
			t.Fatalf("decode stream request body: %v", err)
		}
		if !body.Stream {
			t.Fatalf("stream request did not set stream=true: %s", string(requestBytes))
		}
		if body.Reasoning == nil || body.Reasoning.Summary != "auto" {
			t.Fatalf("stream request did not request reasoning summary: %#v", body.Reasoning)
		}
		streamBody := strings.Join([]string{
			`event: response.created`,
			`data: {"type":"response.created","response":{"id":"resp-stream","status":"in_progress"}}`,
			``,
			`event: response.reasoning_summary_text.delta`,
			`data: {"type":"response.reasoning_summary_text.delta","delta":"checking parent scores"}`,
			``,
			`event: response.reasoning_summary_part.added`,
			`data: {"type":"response.reasoning_summary_part.added","part":{"text":"checking score evidence"}}`,
			``,
			`event: response.reasoning_summary_part.done`,
			`data: {"type":"response.reasoning_summary_part.done","part":{"text":"checked parent scores"}}`,
			``,
			`event: response.output_text.delta`,
			`data: {"type":"response.output_text.delta","delta":"{\"scores\":{"}`,
			``,
			`event: response.output_text.delta`,
			`data: {"type":"response.output_text.delta","delta":"}}"}`,
			``,
			`event: response.completed`,
			`data: {"type":"response.completed","response":{"id":"resp-stream","status":"completed","service_tier":"default","usage":{"input_tokens":10,"input_tokens_details":{"cached_tokens":1},"output_tokens":2},"output":[]}}`,
			``,
		}, "\n")
		return openAITestHTTPResponse(request, http.StatusOK, "req-stream", streamBody), nil
	})
	var debug strings.Builder
	var streamContent strings.Builder
	provider := OpenAIProvider{
		APIKey:              "test-key",
		BaseURL:             "https://example.test/responses",
		Client:              &http.Client{Transport: transport},
		RequestTimeout:      time.Second,
		Stream:              true,
		StreamIdleTimeout:   time.Second,
		StreamContentWriter: &streamContent,
		DebugWriter:         &debug,
		RetryPolicy: ProviderRetryPolicy{
			MaxAttempts: 1,
		},
	}
	response, err := provider.Generate(context.Background(), ProviderRequest{
		APIModel:         "model-a",
		ServiceTier:      serviceTierDefault,
		ReasoningSummary: "auto",
		Prompt:           "stream this",
	})
	if err != nil {
		t.Fatalf("stream generate: %v\n%s", err, debug.String())
	}
	if response.Text != `{"scores":{}}`+"\n" || response.RequestID != "req-stream" || response.ResponseID != "resp-stream" || response.ServiceTier != serviceTierDefault {
		t.Fatalf("unexpected stream response: %#v", response)
	}
	if !strings.Contains(debug.String(), "event=stream_event") || !strings.Contains(debug.String(), `type="response.completed"`) {
		t.Fatalf("stream debug log missing liveness events:\n%s", debug.String())
	}
	for _, forbidden := range []string{
		`{"scores":{`,
		"checking parent scores",
		"checking score evidence",
		"response.output_text.delta",
		"response.reasoning_summary_text.delta",
		"response.reasoning_summary_part.added",
	} {
		if strings.Contains(streamContent.String(), forbidden) {
			t.Fatalf("stream content leaked suppressed stream event %q:\n%s", forbidden, streamContent.String())
		}
	}
	for _, want := range []string{
		`.`,
		`type=response.reasoning_summary_part.done delta="checked parent scores"`,
	} {
		if !strings.Contains(streamContent.String(), want) {
			t.Fatalf("stream content missing %q:\n%s", want, streamContent.String())
		}
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

func TestOpenAIProviderRetriesTransientHTTPStatuses(t *testing.T) {
	for _, statusCode := range []int{
		http.StatusRequestTimeout,
		http.StatusConflict,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
	} {
		t.Run(fmt.Sprintf("status_%d", statusCode), func(t *testing.T) {
			attempts := 0
			transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
				attempts++
				if attempts == 1 {
					return openAITestHTTPResponse(request, statusCode, "", `{"error":{"message":"temporary provider failure"}}`), nil
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
			if _, err := provider.Generate(context.Background(), ProviderRequest{APIModel: "model-a", ServiceTier: serviceTierFlex, Prompt: "retry transient"}); err != nil {
				t.Fatalf("generate after transient retry: %v", err)
			}
			if attempts != 2 {
				t.Fatalf("attempts = %d, want 2", attempts)
			}
		})
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
	assertMissing(t, repo.Abs(strings.TrimSuffix(proposalChildSimulationPath("ga-bad-generate", childID), "/")))
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
		"proposals/ga-test/simulations/SIM-child",
		repo.Rel(resultPath),
		"promotion required before git add",
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
	assertMissing(t, repo.Abs(strings.TrimSuffix(proposalChildSimulationPath("ga-test", "SIM-child"), "/")))
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
	for _, want := range []string{"cull children=SIM-child", "proposals/ga-test/simulations/SIM-child", "proposals/ga-test/results/SIM-child"} {
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
	assertExists(t, repo.Abs(strings.TrimSuffix(proposalChildSimulationPath("ga-test", "SIM-child"), "/")))
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
	if err == nil || !strings.Contains(err.Error(), "path must be proposals/ga-test/simulations/SIM-child/") {
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
	prefix := "- Temporary child ID: `"
	start := strings.Index(prompt, prefix)
	if start < 0 {
		t.Fatalf("prompt missing temporary child ID: %s", prompt)
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
	childPath := strings.TrimSuffix(proposalChildSimulationPath(runGroupID, childID), "/")
	writeTestFile(t, repo.Path(childPath, "README.md"), "# Child\n")
	writeTestFile(t, repo.Path(childPath, "QUESTION.md"), "# Question\n")
	treeHash, err := currentSimulationTreeHash(repo, childPath)
	if err != nil {
		t.Fatalf("hash child tree: %v", err)
	}
	resultPath := repo.Abs(proposalChildResultPath(runGroupID, childID, "scenario-one", "model-a", "20260519-101500"))
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

func stringSliceContains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
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

func writeExactMatchV1Result(t *testing.T, repo Repo, simID string, scenarioID string, modelID string, timestamp string, normalized float64) string {
	t.Helper()
	path := repo.Path("results", simID, scenarioID, modelID, timestamp+".json")
	result := validResult(repo, path)
	result.Schema = resultSchemaV1
	result.RunGroupID = "ga-v1"
	result.CellID = "ga-v1-000001"
	result.ModelID = modelID
	result.ResultPath = repo.Rel(path)
	result.Runner.APIModel = ""
	result.Source.RepoCommit = repo.GitCommit()
	result.Source.SimPath = filepath.ToSlash(filepath.Join("simulations", simID)) + "/"
	result.Source.ScenarioPath = filepath.ToSlash(filepath.Join("scenarios", scenarioID, scenarioID+".md"))
	sourcePaths := []string{
		"results/RUN-PROTOCOL.md",
		"scenarios/README.md",
		filepath.ToSlash(filepath.Join("scenarios", scenarioID, scenarioID+".md")),
		filepath.ToSlash(filepath.Join("simulations", simID, "README.md")),
	}
	questionPath := filepath.ToSlash(filepath.Join("simulations", simID, "QUESTION.md"))
	if _, err := os.Stat(repo.Abs(questionPath)); err == nil {
		sourcePaths = append(sourcePaths, questionPath)
	}
	result.Source.Files = nil
	for _, sourcePath := range sourcePaths {
		hash, err := sha256File(repo, sourcePath)
		if err != nil {
			t.Fatalf("hash source %s: %v", sourcePath, err)
		}
		result.Source.Files = append(result.Source.Files, SourceFile{Path: sourcePath, SHA256: hash})
	}
	treeHash, err := currentSimulationTreeHash(repo, filepath.ToSlash(filepath.Join("simulations", simID)))
	if err != nil {
		t.Fatalf("sim tree hash: %v", err)
	}
	result.Source.SimulationTreeHash = treeHash
	result.Fitness.Raw = normalized
	result.Fitness.Normalized0To100 = normalized
	if err := writeFitnessResultAtomic(path, result); err != nil {
		t.Fatalf("write exact-match result: %v", err)
	}
	return path
}

func writePromotedCanonicalV1Result(t *testing.T, repo Repo, finalSimID string, originalChildSimID string, scenarioID string, modelID string, timestamp string, normalized float64) string {
	t.Helper()
	path := repo.Path("results", finalSimID, scenarioID, modelID, timestamp+".json")
	result := validResult(repo, path)
	result.Schema = resultSchemaV1
	result.RunGroupID = "ga-v1-promoted"
	result.CellID = "ga-v1-promoted-000001"
	result.SimID = finalSimID
	result.ModelID = modelID
	result.ResultPath = repo.Rel(path)
	result.Runner.APIModel = ""
	result.Source.RepoCommit = repo.GitCommit()
	proposalRunGroupID := "ga-promoted-fixture"
	proposalSimPath := filepath.ToSlash(filepath.Join("proposals", proposalRunGroupID, "simulations", originalChildSimID))
	canonicalSimPath := filepath.ToSlash(filepath.Join("simulations", finalSimID))
	result.Source.SimPath = proposalSimPath + "/"
	result.Source.ScenarioPath = filepath.ToSlash(filepath.Join("scenarios", scenarioID, scenarioID+".md"))
	sourcePaths := []string{
		"results/RUN-PROTOCOL.md",
		"scenarios/README.md",
		filepath.ToSlash(filepath.Join("scenarios", scenarioID, scenarioID+".md")),
		filepath.ToSlash(filepath.Join(canonicalSimPath, "README.md")),
	}
	questionPath := filepath.ToSlash(filepath.Join(canonicalSimPath, "QUESTION.md"))
	if _, err := os.Stat(repo.Abs(questionPath)); err == nil {
		sourcePaths = append(sourcePaths, questionPath)
	}
	localSim, err := localMarkdownFiles(repo, canonicalSimPath, map[string]bool{
		"README.md":   true,
		"QUESTION.md": true,
	})
	if err != nil {
		t.Fatalf("local markdown files for promoted result: %v", err)
	}
	sourcePaths = append(sourcePaths, localSim...)
	result.Source.Files = nil
	for _, sourcePath := range uniqueStrings(sourcePaths) {
		hash, err := sha256File(repo, sourcePath)
		if err != nil {
			t.Fatalf("hash promoted source %s: %v", sourcePath, err)
		}
		storedPath := sourcePath
		if relativePath, ok := auditPathRelativeToRoot(sourcePath, canonicalSimPath); ok {
			storedPath = filepath.ToSlash(filepath.Join(proposalSimPath, relativePath))
		}
		result.Source.Files = append(result.Source.Files, SourceFile{Path: storedPath, SHA256: hash})
	}
	treeHash, err := currentSimulationTreeHash(repo, canonicalSimPath)
	if err != nil {
		t.Fatalf("canonical sim tree hash: %v", err)
	}
	result.Source.SimulationTreeHash = treeHash
	result.Promotion = &PromotionInfo{
		PromotionDI:                "DI-test-promotion",
		RunGroupID:                 proposalRunGroupID,
		OriginalChildSimID:         originalChildSimID,
		FinalSimID:                 finalSimID,
		OriginalProposalSimPath:    proposalSimPath + "/",
		OriginalProposalResultPath: filepath.ToSlash(filepath.Join("proposals", proposalRunGroupID, "results", originalChildSimID, scenarioID, modelID, timestamp+".json")),
		CanonicalResultPath:        repo.Rel(path),
		SourceProvenancePolicy:     "Test fixture: canonical storage identity differs from preserved historical source.* proposal paths.",
	}
	result.Fitness.Raw = normalized
	result.Fitness.Normalized0To100 = normalized
	if err := writeFitnessResultAtomic(path, result); err != nil {
		t.Fatalf("write promoted canonical result: %v", err)
	}
	return path
}

func buildExactMatchSourceFiles(t *testing.T, repo Repo, simID string, scenarioID string) []SourceFile {
	t.Helper()
	sourcePaths := []string{
		"results/RUN-PROTOCOL.md",
		"scenarios/README.md",
		filepath.ToSlash(filepath.Join("scenarios", scenarioID, scenarioID+".md")),
		filepath.ToSlash(filepath.Join("simulations", simID, "README.md")),
	}
	questionPath := filepath.ToSlash(filepath.Join("simulations", simID, "QUESTION.md"))
	if _, err := os.Stat(repo.Abs(questionPath)); err == nil {
		sourcePaths = append(sourcePaths, questionPath)
	}
	var files []SourceFile
	for _, sourcePath := range sourcePaths {
		hash, err := sha256File(repo, sourcePath)
		if err != nil {
			t.Fatalf("hash source %s: %v", sourcePath, err)
		}
		files = append(files, SourceFile{Path: sourcePath, SHA256: hash})
	}
	return files
}

func mustSimulationTreeHash(t *testing.T, repo Repo, simID string) string {
	t.Helper()
	treeHash, err := currentSimulationTreeHash(repo, filepath.ToSlash(filepath.Join("simulations", simID)))
	if err != nil {
		t.Fatalf("sim tree hash: %v", err)
	}
	return treeHash
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
    "promise_vocabulary": 4,
    "simplicity_durability": 4,
    "risk_penalty": 1
  },
  "fitness": {
    "raw": 0,
    "normalized_0_100": 0,
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

func invalidScorePayloadMissingV2AxesJSON() string {
	return `{
  "scores": {
    "scenario_fit": 2,
    "promisegrid_alignment": 1,
    "auditability": 3,
    "evolution_safety": 1,
    "layer_boundary_clarity": 1,
    "failure_handling": 1,
    "implementation_plausibility": 3,
    "risk_penalty": 4
  },
  "fitness": {
    "raw": 0,
    "normalized_0_100": 0,
    "confidence_0_1": 0.44
  },
  "assessment": {
    "rationale": "Old 8-axis fixture output.",
    "strengths": ["small fixture"],
    "weaknesses": ["missing rubric-v2 axes"],
    "risks": ["schema drift"],
    "open_questions": ["none for this test"],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}`
}

func scorePayloadWithZeroV2AxesJSON() string {
	return `{
  "scores": {
    "scenario_fit": 2,
    "promisegrid_alignment": 0,
    "auditability": 3,
    "evolution_safety": 1,
    "layer_boundary_clarity": 0,
    "failure_handling": 1,
    "implementation_plausibility": 3,
    "promise_vocabulary": 0,
    "simplicity_durability": 0,
    "risk_penalty": 5
  },
  "fitness": {
    "raw": 0,
    "normalized_0_100": 0,
    "confidence_0_1": 0.9
  },
  "assessment": {
    "rationale": "Zero-valued rubric-v2 axes are valid low scores, not absent fields.",
    "strengths": ["explicit low-score fixture"],
    "weaknesses": ["fixture-only evaluation"],
    "risks": ["serialization regression"],
    "open_questions": ["none for this test"],
    "authority_boundary": "Evidence only; does not settle PromiseGrid design."
  }
}`
}

func validChildBundleJSON(childID string) string {
	generatedID := generatedChildIDForTest(childID)
	return fmt.Sprintf(`{"child_id":%q,"design_delta_summary":"Tighten audit language while keeping the fixture small.","files":[{"path":"README.md","content":"# %s\n\nThis generated child keeps audit evidence local.\n"},{"path":"QUESTION.md","content":"# Questions\n\nDoes this child improve the audit path?\n"}]}`, generatedID, generatedID)
}

func generatedChildIDForTest(plannedID string) string {
	prefix, err := generatedChildIDPrefix(plannedID)
	if err != nil {
		return plannedID
	}
	return prefix + "audit-path"
}

func hasIssue(issues []string, want string) bool {
	for _, issue := range issues {
		if issue == want {
			return true
		}
	}
	return false
}
