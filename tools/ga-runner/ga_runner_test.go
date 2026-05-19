package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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
	if first.Children[0].Operation != "mutation" || first.Children[1].Operation != "crossover" {
		t.Fatalf("unexpected child operations: %#v", first.Children)
	}
}

func TestBuildGenerationPlanValidatesCounts(t *testing.T) {
	population := []PopulationSim{{SimID: "SIM-a", TreeHash: "a"}}
	scenarios := []Scenario{{ScenarioID: "scenario-a"}}
	options := PlanOptions{
		ModelID:       "model-a",
		ParentCount:   1,
		ScenarioCount: 1,
		ChildCount:    1,
		MaxPromotions: 2,
	}
	_, err := buildGenerationPlan(population, scenarios, options)
	if err == nil || !strings.Contains(err.Error(), "max-promotions cannot exceed child-count") {
		t.Fatalf("expected max-promotions validation error, got %v", err)
	}
	options.MaxPromotions = 0
	options.ParentCount = 0
	_, err = buildGenerationPlan(population, scenarios, options)
	if err == nil || !strings.Contains(err.Error(), "parent-count must be positive") {
		t.Fatalf("expected parent-count validation error, got %v", err)
	}
}

func TestRunInitDryRunPrintsGenerationPlan(t *testing.T) {
	repo := newGitTestRepo(t)
	writeTestFile(t, repo.Path("simulations", "SIM-parent", "README.md"), "# Parent\n")
	writeTestFile(t, repo.Path("simulations", "SIM-child-untracked", "README.md"), "# Child\n")
	writeTestFile(t, repo.Path("scenarios", "scenario-one", "scenario-one.md"), "# One\n")
	writeTestFile(t, repo.Path("scenarios", "scenario-two", "scenario-two.md"), "# Two\n")
	gitAdd(t, repo, "simulations/SIM-parent/README.md", "scenarios/scenario-one/scenario-one.md", "scenarios/scenario-two/scenario-two.md")

	var out strings.Builder
	err := runMain([]string{
		"ga-runner", "init",
		"-repo-root", repo.Root,
		"-dry-run",
		"-model", "model-a",
		"-run-group-id", "ga-test",
		"-shuffle-seed", "7",
		"-parent-count", "1",
		"-scenario-count", "2",
		"-child-count", "2",
		"-max-promotions", "1",
	}, &out, &out)
	if err != nil {
		t.Fatalf("init dry-run plan: %v", err)
	}
	text := out.String()
	for _, want := range []string{"population=1", "plan run_group_id=ga-test model=model-a", "parent_score_cells=2", "child_score_cells=4", "planned-child-0001"} {
		if !strings.Contains(text, want) {
			t.Fatalf("expected %q in output:\n%s", want, text)
		}
	}
	if strings.Contains(text, "SIM-child-untracked") {
		t.Fatalf("untracked child should not appear in plan output:\n%s", text)
	}
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

func hasIssue(issues []string, want string) bool {
	for _, issue := range issues {
		if issue == want {
			return true
		}
	}
	return false
}
