package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func runQueue(ctx context.Context, args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	manifest := fs.String("manifest", "", "matrix manifest CSV path")
	providerName := fs.String("provider", "openai", "provider name; v1 supports openai")
	apiModel := fs.String("api-model", "", "provider API model name")
	reasoningEffort := fs.String("reasoning-effort", "xhigh", "provider reasoning effort")
	apiKeyEnv := fs.String("api-key-env", "OPENAI_API_KEY", "environment variable holding provider API key")
	openAIBaseURL := fs.String("openai-base-url", "", "optional OpenAI Responses API URL override")
	maxOutputTokens := fs.Int("max-output-tokens", 12000, "maximum provider output tokens")
	stateFlag := fs.String("state", "", "queue state JSON path")
	jobDirFlag := fs.String("job-dir", "", "prompt audit directory")
	startIndex := fs.Int("start-index", 0, "zero-based queue start index")
	limit := fs.Int("limit", -1, "maximum selected cells this invocation")
	retryFailed := fs.Bool("retry-failed", false, "retry failed cells")
	rerunDone := fs.Bool("rerun-done", false, "rerun done cells")
	noMatrixUpdate := fs.Bool("no-matrix-update", false, "do not update scenario MATRIX.md")
	dryRun := fs.Bool("dry-run", false, "show selected work without writing state/results/matrices")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *manifest == "" {
		return errUsage("run: --manifest is required")
	}
	if *apiModel == "" && !*dryRun {
		return errUsage("run: --api-model is required for non-dry runs")
	}
	repo, err := openRepo(*repoRoot)
	if err != nil {
		return err
	}
	allCells, err := readManifest(repo, *manifest)
	if err != nil {
		return err
	}
	cells, err := selectedCells(allCells, *startIndex, *limit)
	if err != nil {
		return err
	}
	statePath := *stateFlag
	if statePath == "" {
		statePath = defaultStatePath(repo, cells)
	} else {
		statePath = repo.Abs(statePath)
	}
	jobDir := *jobDirFlag
	if jobDir == "" {
		jobDir = defaultJobDir(repo, cells)
	} else {
		jobDir = repo.Abs(jobDir)
	}
	state, err := loadOrCreateState(repo, statePath, repo.Abs(*manifest), cells)
	if err != nil {
		return err
	}
	provider, err := buildProvider(*providerName, *apiModel, *apiKeyEnv, *openAIBaseURL, *dryRun)
	if err != nil {
		return err
	}
	processed := 0
	failed := 0
	options := RunOptions{
		ProviderName:    *providerName,
		APIModel:        *apiModel,
		ReasoningEffort: *reasoningEffort,
		MaxOutputTokens: *maxOutputTokens,
		RetryFailed:     *retryFailed,
		RerunDone:       *rerunDone,
		UpdateMatrix:    !*noMatrixUpdate,
		DryRun:          *dryRun,
	}
	for _, cell := range cells {
		record := state.Cells[cell.CellID]
		if skip, reason := shouldSkip(record, options); skip {
			fmt.Fprintf(stdout, "skip: %s: %s\n", cell.CellID, reason)
			continue
		}
		status := runOneCell(ctx, repo, provider, state, statePath, jobDir, cell, record, options, stdout)
		processed++
		if status == "failed" {
			failed++
		}
		if !options.DryRun {
			if err := saveState(statePath, state); err != nil {
				return err
			}
		}
		fmt.Fprintln(stdout, formatCounts(state))
	}
	if !options.DryRun {
		if err := saveState(statePath, state); err != nil {
			return err
		}
	}
	fmt.Fprintf(stdout, "processed=%d failed=%d state=%s\n", processed, failed, statePath)
	if failed > 0 {
		return fmt.Errorf("one or more cells failed")
	}
	return nil
}

type RunOptions struct {
	ProviderName    string
	APIModel        string
	ReasoningEffort string
	MaxOutputTokens int
	RetryFailed     bool
	RerunDone       bool
	UpdateMatrix    bool
	DryRun          bool
}

func buildProvider(providerName, apiModel, apiKeyEnv, openAIBaseURL string, dryRun bool) (Provider, error) {
	if dryRun {
		return nil, nil
	}
	switch providerName {
	case "openai":
		return OpenAIProvider{
			APIKey:  os.Getenv(apiKeyEnv),
			BaseURL: openAIBaseURL,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported provider %q", providerName)
	}
}

func shouldSkip(record *CellState, options RunOptions) (bool, string) {
	switch record.Status {
	case "done":
		if !options.RerunDone {
			return true, "already done"
		}
	case "failed":
		if !options.RetryFailed {
			return true, "failed; pass --retry-failed to retry"
		}
	case "skipped":
		if !options.RerunDone {
			return true, "terminal status skipped"
		}
	}
	return false, ""
}

func runOneCell(ctx context.Context, repo Repo, provider Provider, state *QueueState, statePath string, jobDir string, cell MatrixCell, record *CellState, options RunOptions, stdout io.Writer) string {
	resultPath := repo.Abs(cell.ResultPath)
	if options.DryRun {
		fmt.Fprintf(stdout, "dry-run: %s -> %s\n", cell.CellID, cell.ResultPath)
		return record.Status
	}
	prompt, err := PromptBuilder{Repo: repo}.BuildAPIPrompt(cell)
	if err != nil {
		markCell(record, "failed", err.Error())
		return "failed"
	}
	if err := writeFile(filepath.Join(jobDir, promptFilename(cell)), prompt); err != nil {
		markCell(record, "failed", err.Error())
		return "failed"
	}
	if issues := validateResultFile(repo, resultPath, false, false); len(issues) == 0 {
		if options.UpdateMatrix {
			if _, _, _, err := updateMatrixForResult(repo, resultPath, false); err != nil {
				markCell(record, "failed", err.Error())
				return "failed"
			}
		}
		markCell(record, "done", "existing valid result")
		return "done"
	}
	record.Attempts++
	record.Provider = options.ProviderName
	record.APIModel = options.APIModel
	record.ReasoningEffort = options.ReasoningEffort
	markCell(record, "running", "provider call started")
	// Intent: Persist the in-flight cell before the provider call so an
	// interrupted unattended run names the active cell on restart.
	// Source: DI-bujiv; DI-lulom
	if err := saveState(statePath, state); err != nil {
		markCell(record, "failed", err.Error())
		return "failed"
	}
	response, err := provider.Generate(ctx, ProviderRequest{
		Provider:        options.ProviderName,
		APIModel:        options.APIModel,
		ReasoningEffort: options.ReasoningEffort,
		MaxOutputTokens: options.MaxOutputTokens,
		Prompt:          prompt,
	})
	if err != nil {
		markCell(record, "failed", err.Error())
		return "failed"
	}
	record.RequestID = response.RequestID
	record.ResponseID = response.ResponseID
	record.UsageJSON = response.UsageJSON
	if err := writeFile(resultPath, response.Text); err != nil {
		markCell(record, "failed", err.Error())
		return "failed"
	}
	if issues := validateResultFile(repo, resultPath, false, false); len(issues) > 0 {
		markCell(record, "failed", strings.Join(issues, "; "))
		return "failed"
	}
	if options.UpdateMatrix {
		if _, _, _, err := updateMatrixForResult(repo, resultPath, false); err != nil {
			markCell(record, "failed", err.Error())
			return "failed"
		}
	}
	markCell(record, "done", "validated result")
	return "done"
}

func runProgress(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("progress", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	manifest := fs.String("manifest", "", "matrix manifest CSV path")
	stateFlag := fs.String("state", "", "queue state JSON path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *manifest == "" {
		return errUsage("progress: --manifest is required")
	}
	repo, err := openRepo(*repoRoot)
	if err != nil {
		return err
	}
	cells, err := readManifest(repo, *manifest)
	if err != nil {
		return err
	}
	statePath := *stateFlag
	if statePath == "" {
		statePath = defaultStatePath(repo, cells)
	} else {
		statePath = repo.Abs(statePath)
	}
	state, err := loadOrCreateState(repo, statePath, repo.Abs(*manifest), cells)
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout, formatCounts(state))
	fmt.Fprintf(stdout, "state=%s\n", statePath)
	return nil
}
