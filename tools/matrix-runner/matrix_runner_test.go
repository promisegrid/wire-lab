package main

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestManifestJobsViewAndValidate(t *testing.T) {
	repo := makeTestRepo(t)
	var stdout bytes.Buffer
	manifestPath := filepath.Join(repo.Root, "results", "manifests", "test.csv")
	err := runMain(context.Background(), []string{
		"matrix-runner", "manifest",
		"-repo-root", repo.Root,
		"-models", "openai-test-xhigh",
		"-run-group-id", "rg",
		"-timestamp", "20260519-000000",
		"-output", manifestPath,
	}, &stdout, &stdout)
	if err != nil {
		t.Fatalf("manifest: %v\n%s", err, stdout.String())
	}
	cells, err := readManifest(repo, manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(cells) != 1 {
		t.Fatalf("cells=%d want 1", len(cells))
	}
	if cells[0].ResultPath != "results/SIM-alpha/scenario-one/openai-test-xhigh/20260519-000000.md" {
		t.Fatalf("unexpected result path: %s", cells[0].ResultPath)
	}
	stdout.Reset()
	err = runMain(context.Background(), []string{
		"matrix-runner", "jobs",
		"-repo-root", repo.Root,
		"-manifest", manifestPath,
	}, &stdout, &stdout)
	if err != nil {
		t.Fatalf("jobs: %v\n%s", err, stdout.String())
	}
	resultPath := repo.Abs(cells[0].ResultPath)
	if err := writeFile(resultPath, validResult(cells[0])); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	err = runMain(context.Background(), []string{
		"matrix-runner", "view",
		"-repo-root", repo.Root,
		"-scenario", "scenario-one",
	}, &stdout, &stdout)
	if err != nil {
		t.Fatalf("view: %v\n%s", err, stdout.String())
	}
	viewText := stdout.String()
	if !strings.Contains(viewText, repo.Rel(resultPath)) || !strings.Contains(viewText, "good partial fit") {
		t.Fatalf("view missing result evidence:\n%s", viewText)
	}
	stdout.Reset()
	err = runMain(context.Background(), []string{
		"matrix-runner", "validate",
		"-repo-root", repo.Root,
		"-manifest", manifestPath,
	}, &stdout, &stdout)
	if err != nil {
		t.Fatalf("validate: %v\n%s", err, stdout.String())
	}
	if !strings.Contains(stdout.String(), "validated=1 failed=0") {
		t.Fatalf("unexpected validate output: %s", stdout.String())
	}
}

func TestRunCellPersistsRunningBeforeProviderCall(t *testing.T) {
	repo := makeTestRepo(t)
	cell := MatrixCell{
		RunGroupID:   "rg",
		Ordinal:      1,
		CellID:       "rg-000001-SIM-alpha--scenario-one--openai-test-xhigh",
		SimID:        "SIM-alpha",
		ScenarioID:   "scenario-one",
		ModelID:      "openai-test-xhigh",
		SimPath:      "simulations/SIM-alpha/",
		ScenarioPath: "scenarios/scenario-one/scenario-one.md",
		Timestamp:    "20260519-000000",
		ResultPath:   "results/SIM-alpha/scenario-one/openai-test-xhigh/20260519-000000.md",
		Status:       "queued",
	}
	state := &QueueState{
		Manifest:   "manifest.csv",
		RunGroupID: "rg",
		CreatedAt:  utcISO(),
		Cells: map[string]*CellState{
			cell.CellID: {
				CellID:     cell.CellID,
				Status:     "queued",
				ResultPath: cell.ResultPath,
			},
		},
	}
	statePath := filepath.Join(repo.Root, "results", "state", "rg.json")
	provider := fakeProvider{
		text: validResult(cell),
		check: func() {
			bytes, err := os.ReadFile(statePath)
			if err != nil {
				t.Fatalf("state not persisted before provider call: %v", err)
			}
			var saved QueueState
			if err := json.Unmarshal(bytes, &saved); err != nil {
				t.Fatal(err)
			}
			if got := saved.Cells[cell.CellID].Status; got != "running" {
				t.Fatalf("saved status=%s want running", got)
			}
		},
	}
	var stdout bytes.Buffer
	status := runOneCell(context.Background(), repo, provider, state, statePath, filepath.Join(repo.Root, "results", "jobs", "rg"), cell, state.Cells[cell.CellID], RunOptions{
		ProviderName:    "openai",
		APIModel:        "gpt-test",
		ReasoningEffort: "xhigh",
		MaxOutputTokens: 1000,
	}, &stdout)
	if status != "done" {
		t.Fatalf("status=%s state=%+v output=%s", status, state.Cells[cell.CellID], stdout.String())
	}
}

type fakeProvider struct {
	text  string
	check func()
}

func (f fakeProvider) Generate(ctx context.Context, request ProviderRequest) (ProviderResponse, error) {
	if f.check != nil {
		f.check()
	}
	return ProviderResponse{Text: f.text, RequestID: "req_test", ResponseID: "resp_test"}, nil
}

func makeTestRepo(t *testing.T) Repo {
	t.Helper()
	root := t.TempDir()
	repo := Repo{Root: root}
	mustWrite(t, repo.Path(".git", "HEAD"), "ref: refs/heads/test\n")
	mustWrite(t, repo.Path("results", "RUN-PROTOCOL.md"), "# Results Run Protocol\n\n## Result File Contract\n")
	mustWrite(t, repo.Path("simulations", "SIM-alpha", "README.md"), "# SIM-alpha\n")
	mustWrite(t, repo.Path("simulations", "SIM-alpha", "QUESTION.md"), "# Question\n")
	mustWrite(t, repo.Path("simulations", "SIM-alpha", "protocols", "draft.md"), "# Draft\n")
	mustWrite(t, repo.Path("scenarios", "README.md"), "# Scenarios\n")
	mustWrite(t, repo.Path("scenarios", "scenario-one", "README.md"), "# Scenario One\n")
	mustWrite(t, repo.Path("scenarios", "scenario-one", "scenario-one.md"), "# Scenario One Body\n")
	return repo
}

func mustWrite(t *testing.T, path string, text string) {
	t.Helper()
	if err := writeFile(path, text); err != nil {
		t.Fatal(err)
	}
}

func validResult(cell MatrixCell) string {
	return "# Result: " + cell.SimID + " / " + cell.ScenarioID + " / " + cell.ModelID + " / " + cell.Timestamp + "\n\n" +
		"## Result ID\n\n" + cell.SimID + "-" + cell.ScenarioID + "-" + cell.ModelID + "-" + cell.Timestamp + "\n\n" +
		"## Scenario\n\n- Scenario ID: `" + cell.ScenarioID + "`\n- Scenario path: `" + cell.ScenarioPath + "`\n\n" +
		"## Simulation\n\n- Simulation ID: `" + cell.SimID + "`\n- Simulation path: `" + cell.SimPath + "`\n- Simulation commit: `test`\n\n" +
		"## Runner\n\n- Runner/interface: `matrix-runner openai`\n- Model ID: `" + cell.ModelID + "`\n- Run timestamp UTC: `" + cell.Timestamp + "`\n- Operator: `matrix-runner`\n\n" +
		"## Prompt / Procedure\n\n- Run mode: `llm-doc-eval-blind`\n\n" +
		"## Observed Behavior\n\nThe simulation is evaluated against the scenario.\n\n" +
		"## Verdict\n\nEvidence verdict: good partial fit for test coverage\n\n" +
		"## Evidence Links\n\n- Scenario: `" + cell.ScenarioPath + "`\n- Simulation docs: `" + cell.SimPath + "README.md`\n\n" +
		"## Open Questions\n\n- None for test.\n\n" +
		"## Handoff Target\n\nNone.\n\n" +
		"## Authority Boundary\n\nThis result is evidence only.\n"
}
