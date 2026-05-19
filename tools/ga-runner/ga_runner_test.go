package main

import (
	"encoding/json"
	"os"
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

func newTestRepo(t *testing.T) Repo {
	t.Helper()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatalf("make .git: %v", err)
	}
	return Repo{Root: root}
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
