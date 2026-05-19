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
	maxOutputTokens := fs.Int("max-output-tokens", 6000, "maximum provider output tokens")
	resultStyle := fs.String("result-style", "concise", "result style: concise or standard")
	maxRunCost := fs.Float64("max-run-cost-usd", 0, "stop before starting a cell that would exceed this run budget; 0 disables")
	maxCellEstimate := fs.Float64("max-cell-estimate-usd", 0, "skip cells whose preflight worst-case cost estimate exceeds this amount; 0 disables")
	inputPrice := fs.Float64("cost-input-usd-per-mtok", defaultInputUSDPerMTok, "uncached input price in USD per million tokens")
	cachedInputPrice := fs.Float64("cost-cached-input-usd-per-mtok", defaultCachedInputUSDPerMTok, "cached input price in USD per million tokens")
	outputPrice := fs.Float64("cost-output-usd-per-mtok", defaultOutputUSDPerMTok, "output price in USD per million tokens")
	stateFlag := fs.String("state", "", "queue state JSON path")
	jobDirFlag := fs.String("job-dir", "", "prompt audit directory")
	startIndex := fs.Int("start-index", 0, "zero-based queue start index")
	limit := fs.Int("limit", -1, "maximum selected cells this invocation")
	retryFailed := fs.Bool("retry-failed", false, "retry failed cells")
	rerunDone := fs.Bool("rerun-done", false, "rerun done cells")
	dryRun := fs.Bool("dry-run", false, "show selected work without writing state/results")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *manifest == "" {
		return errUsage("run: --manifest is required")
	}
	if *apiModel == "" && !*dryRun {
		return errUsage("run: --api-model is required for non-dry runs")
	}
	if !validResultStyle(*resultStyle) {
		return errUsage("run: --result-style must be concise or standard")
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
		ResultStyle:     *resultStyle,
		RetryFailed:     *retryFailed,
		RerunDone:       *rerunDone,
		DryRun:          *dryRun,
		Cost: CostConfig{
			InputUSDPerMTok:       *inputPrice,
			CachedInputUSDPerMTok: *cachedInputPrice,
			OutputUSDPerMTok:      *outputPrice,
			MaxRunUSD:             *maxRunCost,
			MaxCellEstimateUSD:    *maxCellEstimate,
		},
	}
	for _, cell := range cells {
		// Intent: Stop unattended batches before the next provider call when the
		// already-recorded usage reaches the user-approved run budget.
		// Source: DI-nugiv
		if options.Cost.MaxRunUSD > 0 && stateActualCostUSD(state) >= options.Cost.MaxRunUSD {
			if err := writeFormat(stdout, "budget-stop: actual_cost_usd=%.6f max_run_cost_usd=%.6f\n", stateActualCostUSD(state), options.Cost.MaxRunUSD); err != nil {
				return err
			}
			break
		}
		record := state.Cells[cell.CellID]
		if skip, reason := shouldSkip(record, options); skip {
			if err := writeFormat(stdout, "skip: %s: %s\n", cell.CellID, reason); err != nil {
				return err
			}
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
		if err := writeLine(stdout, formatCounts(state)); err != nil {
			return err
		}
		if err := writeFormat(stdout, "cost actual_usd=%.6f\n", stateActualCostUSD(state)); err != nil {
			return err
		}
		if status == "budget-stop" {
			break
		}
	}
	if !options.DryRun {
		if err := saveState(statePath, state); err != nil {
			return err
		}
	}
	if err := writeFormat(stdout, "processed=%d failed=%d state=%s\n", processed, failed, statePath); err != nil {
		return err
	}
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
	ResultStyle     string
	RetryFailed     bool
	RerunDone       bool
	DryRun          bool
	Cost            CostConfig
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
		if err := writeFormat(stdout, "dry-run: %s -> %s\n", cell.CellID, cell.ResultPath); err != nil {
			markCell(record, "failed", err.Error())
			return "failed"
		}
		return record.Status
	}
	prompt, err := PromptBuilder{Repo: repo, ResultStyle: options.ResultStyle}.BuildAPIPrompt(cell)
	if err != nil {
		markCell(record, "failed", err.Error())
		return "failed"
	}
	// Intent: Use a conservative per-cell estimate before spending API tokens;
	// actual accounting is recorded later from provider usage metadata.
	// Source: DI-nugiv
	estimate := options.Cost.EstimatePromptCost(prompt, options.MaxOutputTokens)
	if options.Cost.MaxCellEstimateUSD > 0 && estimate.CostUSD > options.Cost.MaxCellEstimateUSD {
		markCell(record, "skipped", fmt.Sprintf("estimated cell cost %.6f exceeds max %.6f", estimate.CostUSD, options.Cost.MaxCellEstimateUSD))
		return "skipped"
	}
	if options.Cost.MaxRunUSD > 0 && stateActualCostUSD(state)+estimate.CostUSD > options.Cost.MaxRunUSD {
		if err := writeFormat(stdout, "budget-stop: actual_cost_usd=%.6f estimated_next_cell_usd=%.6f max_run_cost_usd=%.6f\n", stateActualCostUSD(state), estimate.CostUSD, options.Cost.MaxRunUSD); err != nil {
			markCell(record, "failed", err.Error())
			return "failed"
		}
		return "budget-stop"
	}
	if err := writeFile(filepath.Join(jobDir, promptFilename(cell)), prompt); err != nil {
		markCell(record, "failed", err.Error())
		return "failed"
	}
	if issues := validateResultFile(repo, resultPath, false); len(issues) == 0 {
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
	// Intent: Persist actual provider usage so resumed runs and post-run reviews
	// can reason from measured cost instead of prompt-size estimates.
	// Source: DI-nugiv
	if response.UsageJSON != "" {
		usageCost, err := options.Cost.ParseUsage(response.UsageJSON)
		if err != nil {
			markCell(record, "failed", err.Error())
			return "failed"
		}
		record.InputTokens = usageCost.InputTokens
		record.CachedTokens = usageCost.CachedInputTokens
		record.OutputTokens = usageCost.OutputTokens
		record.CostUSD = usageCost.CostUSD
	} else if options.Cost.BudgetEnabled() {
		markCell(record, "failed", "missing provider usage metadata for cost-controlled run")
		return "failed"
	}
	if err := writeFile(resultPath, response.Text); err != nil {
		markCell(record, "failed", err.Error())
		return "failed"
	}
	if issues := validateResultFile(repo, resultPath, false); len(issues) > 0 {
		markCell(record, "failed", strings.Join(issues, "; "))
		return "failed"
	}
	if record.CostUSD > 0 {
		markCell(record, "done", fmt.Sprintf("validated result; cost_usd=%.6f", record.CostUSD))
	} else {
		markCell(record, "done", "validated result")
	}
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
	if err := writeLine(stdout, formatCounts(state)); err != nil {
		return err
	}
	if err := writeFormat(stdout, "cost actual_usd=%.6f\n", stateActualCostUSD(state)); err != nil {
		return err
	}
	return writeFormat(stdout, "state=%s\n", statePath)
}
