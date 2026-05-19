package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"
)

type scoreOptions struct {
	RepoRoot         string
	RunGroupID       string
	Target           string
	ProviderName     string
	APIModel         string
	ReasoningEffort  string
	ServiceTier      string
	APIKeyEnv        string
	OpenAIBaseURL    string
	MaxOutputTokens  int
	MaxRunCostUSD    float64
	MaxCellCostUSD   float64
	InputPrice       float64
	CachedInputPrice float64
	OutputPrice      float64
	RetryFailed      bool
	RerunDone        bool
	DryRun           bool
	JobDir           string
}

type scorePayload struct {
	Scores     FitnessScores  `json:"scores"`
	Fitness    FitnessSummary `json:"fitness"`
	Assessment Assessment     `json:"assessment"`
}

func runScore(args []string, stdout io.Writer) error {
	options, err := parseScoreOptions(args)
	if err != nil {
		return err
	}
	repo, err := openRepo(options.RepoRoot)
	if err != nil {
		return err
	}
	provider, err := buildProvider(options.ProviderName, options.APIKeyEnv, options.OpenAIBaseURL, options.DryRun)
	if err != nil {
		return err
	}
	return runScoreWithProvider(context.Background(), repo, provider, options, stdout)
}

func parseScoreOptions(args []string) (scoreOptions, error) {
	fs := flag.NewFlagSet("score", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	runGroupID := fs.String("run-group-id", "", "GA run group to score")
	target := fs.String("target", "parents", "cells to score: parents, children, or all")
	providerName := fs.String("provider", "openai", "provider name; v1 supports openai")
	apiModel := fs.String("api-model", "", "provider API model name")
	reasoningEffort := fs.String("reasoning-effort", "xhigh", "provider reasoning effort")
	// Intent: Default scoring calls to Flex and reject Priority so unattended
	// scoring cannot inherit expensive processing. Source: DI-mopob
	serviceTier := fs.String("service-tier", defaultServiceTier, "provider service tier: flex or default; priority is rejected")
	apiKeyEnv := fs.String("api-key-env", "OPENAI_API_KEY", "environment variable holding provider API key")
	openAIBaseURL := fs.String("openai-base-url", "", "optional OpenAI Responses API URL override")
	maxOutputTokens := fs.Int("max-output-tokens", 4000, "maximum provider output tokens")
	maxRunCost := fs.Float64("max-run-cost-usd", 0, "stop before starting a cell that would exceed this run budget; 0 disables")
	maxCellCost := fs.Float64("max-cell-estimate-usd", 0, "skip cells whose preflight worst-case cost estimate exceeds this amount; 0 disables")
	inputPrice := fs.Float64("cost-input-usd-per-mtok", defaultInputUSDPerMTok, "uncached input price in USD per million tokens")
	cachedInputPrice := fs.Float64("cost-cached-input-usd-per-mtok", defaultCachedInputUSDPerMTok, "cached input price in USD per million tokens")
	outputPrice := fs.Float64("cost-output-usd-per-mtok", defaultOutputUSDPerMTok, "output price in USD per million tokens")
	retryFailed := fs.Bool("retry-failed", false, "retry failed cells")
	rerunDone := fs.Bool("rerun-done", false, "rerun done cells")
	dryRun := fs.Bool("dry-run", false, "show selected work without writing state/results")
	jobDir := fs.String("job-dir", "", "prompt audit directory; defaults to results/jobs/<run-group-id>")
	if err := fs.Parse(args); err != nil {
		return scoreOptions{}, err
	}
	if *runGroupID == "" {
		return scoreOptions{}, errUsage("score: -run-group-id is required")
	}
	if *apiModel == "" && !*dryRun {
		return scoreOptions{}, errUsage("score: -api-model is required for non-dry runs")
	}
	if *target != "parents" && *target != "children" && *target != "all" {
		return scoreOptions{}, errUsage("score: -target must be parents, children, or all")
	}
	normalizedServiceTier, err := normalizeServiceTier(*serviceTier)
	if err != nil {
		return scoreOptions{}, errUsage("score: " + err.Error())
	}
	return scoreOptions{
		RepoRoot:         *repoRoot,
		RunGroupID:       *runGroupID,
		Target:           *target,
		ProviderName:     *providerName,
		APIModel:         *apiModel,
		ReasoningEffort:  *reasoningEffort,
		ServiceTier:      normalizedServiceTier,
		APIKeyEnv:        *apiKeyEnv,
		OpenAIBaseURL:    *openAIBaseURL,
		MaxOutputTokens:  *maxOutputTokens,
		MaxRunCostUSD:    *maxRunCost,
		MaxCellCostUSD:   *maxCellCost,
		InputPrice:       *inputPrice,
		CachedInputPrice: *cachedInputPrice,
		OutputPrice:      *outputPrice,
		RetryFailed:      *retryFailed,
		RerunDone:        *rerunDone,
		DryRun:           *dryRun,
		JobDir:           *jobDir,
	}, nil
}

func runScoreWithProvider(ctx context.Context, repo Repo, provider Provider, options scoreOptions, stdout io.Writer) error {
	serviceTier, err := normalizeServiceTier(options.ServiceTier)
	if err != nil {
		return err
	}
	options.ServiceTier = serviceTier
	stateFile, err := statePath(repo, options.RunGroupID)
	if err != nil {
		return err
	}
	state, err := readGAState(stateFile)
	if err != nil {
		return err
	}
	if state.RunGroupID != options.RunGroupID {
		return fmt.Errorf("state run_group_id %q does not match requested %q", state.RunGroupID, options.RunGroupID)
	}
	cells := selectedScoreCellIndexes(state, options)
	if len(cells) == 0 {
		return fmt.Errorf("no score cells matched target %s", options.Target)
	}
	jobDir := options.JobDir
	if jobDir == "" {
		jobDir = repo.Path("results", "jobs", options.RunGroupID)
	} else {
		jobDir = repo.Abs(jobDir)
	}
	cost := CostConfig{
		InputUSDPerMTok:       options.InputPrice,
		CachedInputUSDPerMTok: options.CachedInputPrice,
		OutputUSDPerMTok:      options.OutputPrice,
		MaxRunUSD:             options.MaxRunCostUSD,
		MaxCellEstimateUSD:    options.MaxCellCostUSD,
	}
	processed := 0
	failed := 0
	for _, index := range cells {
		if cost.MaxRunUSD > 0 && stateActualCostUSD(state) >= cost.MaxRunUSD {
			if err := writeFormat(stdout, "budget-stop: actual_cost_usd=%.6f max_run_cost_usd=%.6f\n", stateActualCostUSD(state), cost.MaxRunUSD); err != nil {
				return err
			}
			break
		}
		status := scoreOneCell(ctx, repo, provider, &state, stateFile, jobDir, index, options, cost, stdout)
		processed++
		if status == "failed" {
			failed++
		}
		if !options.DryRun {
			state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
			if err := writeGAStateAtomic(stateFile, state); err != nil {
				return err
			}
		}
		if status == "budget-stop" {
			break
		}
	}
	if err := writeFormat(stdout, "processed=%d failed=%d state=%s\n", processed, failed, repo.Rel(stateFile)); err != nil {
		return err
	}
	if failed > 0 {
		return fmt.Errorf("one or more score cells failed")
	}
	return nil
}

func selectedScoreCellIndexes(state GAState, options scoreOptions) []int {
	parentIDs := map[string]bool{}
	for _, parent := range state.Parents {
		parentIDs[parent.SimID] = true
	}
	childIDs := map[string]bool{}
	for _, child := range state.Children {
		childIDs[child.ID()] = true
	}
	var indexes []int
	for index, cell := range state.Cells {
		if options.Target == "parents" && !parentIDs[cell.SimID] {
			continue
		}
		if options.Target == "children" && !childIDs[cell.SimID] {
			continue
		}
		if skipCellStatus(cell, options.RetryFailed, options.RerunDone) {
			continue
		}
		indexes = append(indexes, index)
	}
	return indexes
}

func skipCellStatus(cell GACell, retryFailed bool, rerunDone bool) bool {
	switch cell.Status {
	case "done":
		return !rerunDone
	case "failed":
		return !retryFailed
	case "skipped":
		return !rerunDone
	default:
		return false
	}
}

func scoreOneCell(ctx context.Context, repo Repo, provider Provider, state *GAState, stateFile string, jobDir string, index int, options scoreOptions, cost CostConfig, stdout io.Writer) string {
	cell := &state.Cells[index]
	scenario, err := scenarioFromState(*state, cell.ScenarioID)
	if err != nil {
		markGACell(cell, "failed", err.Error())
		return "failed"
	}
	if options.DryRun {
		if err := writeFormat(stdout, "dry-run score: %s -> %s\n", cell.CellID, cell.ResultPath); err != nil {
			markGACell(cell, "failed", err.Error())
			return "failed"
		}
		return cell.Status
	}
	prompt, err := buildScorePrompt(repo, *state, *cell, scenario)
	if err != nil {
		markGACell(cell, "failed", err.Error())
		return "failed"
	}
	// Intent: Estimate worst-case token cost before each provider call so a
	// long GA scoring run can stop before exceeding the user-approved budget.
	// Source: DI-gijom
	estimate := cost.EstimatePromptCost(prompt, options.MaxOutputTokens)
	if cost.MaxCellEstimateUSD > 0 && estimate.CostUSD > cost.MaxCellEstimateUSD {
		markGACell(cell, "skipped", fmt.Sprintf("estimated cell cost %.6f exceeds max %.6f", estimate.CostUSD, cost.MaxCellEstimateUSD))
		return "skipped"
	}
	if cost.MaxRunUSD > 0 && stateActualCostUSD(*state)+estimate.CostUSD > cost.MaxRunUSD {
		if err := writeFormat(stdout, "budget-stop: actual_cost_usd=%.6f estimated_next_cell_usd=%.6f max_run_cost_usd=%.6f\n", stateActualCostUSD(*state), estimate.CostUSD, cost.MaxRunUSD); err != nil {
			markGACell(cell, "failed", err.Error())
			return "failed"
		}
		return "budget-stop"
	}
	if err := writeFile(filepath.Join(jobDir, cell.CellID+".score.md"), prompt); err != nil {
		markGACell(cell, "failed", err.Error())
		return "failed"
	}
	if issues := validateResultFile(repo, repo.Abs(cell.ResultPath)); len(issues) == 0 {
		markGACell(cell, "done", "existing valid result")
		return "done"
	}
	cell.Attempts++
	cell.Provider = options.ProviderName
	cell.APIModel = options.APIModel
	cell.ReasoningEffort = options.ReasoningEffort
	cell.ServiceTier = options.ServiceTier
	markGACell(cell, "running", "provider call started")
	state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := writeGAStateAtomic(stateFile, *state); err != nil {
		markGACell(cell, "failed", err.Error())
		return "failed"
	}
	response, err := provider.Generate(ctx, ProviderRequest{
		Provider:        options.ProviderName,
		APIModel:        options.APIModel,
		ReasoningEffort: options.ReasoningEffort,
		ServiceTier:     options.ServiceTier,
		MaxOutputTokens: options.MaxOutputTokens,
		Instructions:    "Return only valid JSON for the requested GA score payload. Do not include code fences or commentary.",
		Prompt:          prompt,
	})
	if err != nil {
		markGACell(cell, "failed", err.Error())
		return "failed"
	}
	cell.RequestID = response.RequestID
	cell.ResponseID = response.ResponseID
	cell.ServedServiceTier = response.ServiceTier
	cell.UsageJSON = response.UsageJSON
	// Intent: Record measured provider usage in the state and result so resumed
	// runs and cost reviews use actual token counts when the provider returns
	// them. Source: DI-gijom
	if response.UsageJSON != "" {
		usageCost, err := cost.ParseUsage(response.UsageJSON)
		if err != nil {
			markGACell(cell, "failed", err.Error())
			return "failed"
		}
		cell.InputTokens = usageCost.InputTokens
		cell.CachedTokens = usageCost.CachedInputTokens
		cell.OutputTokens = usageCost.OutputTokens
		cell.CostUSD = usageCost.CostUSD
	} else if cost.BudgetEnabled() {
		markGACell(cell, "failed", "missing provider usage metadata for cost-controlled run")
		return "failed"
	}
	payload, err := parseScorePayload(response.Text)
	if err != nil {
		markGACell(cell, "failed", err.Error())
		return "failed"
	}
	result, err := buildFitnessResult(repo, *state, *cell, scenario, payload)
	if err != nil {
		markGACell(cell, "failed", err.Error())
		return "failed"
	}
	if err := writeFitnessResultAtomic(repo.Abs(cell.ResultPath), result); err != nil {
		markGACell(cell, "failed", err.Error())
		return "failed"
	}
	if issues := validateResultFile(repo, repo.Abs(cell.ResultPath)); len(issues) > 0 {
		markGACell(cell, "failed", strings.Join(issues, "; "))
		return "failed"
	}
	markGACell(cell, "done", "validated result")
	return "done"
}

func parseScorePayload(text string) (scorePayload, error) {
	clean := strings.TrimSpace(text)
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimPrefix(clean, "```")
	clean = strings.TrimSuffix(clean, "```")
	var payload scorePayload
	if err := json.Unmarshal([]byte(strings.TrimSpace(clean)), &payload); err != nil {
		return scorePayload{}, err
	}
	return payload, nil
}

func buildScorePrompt(repo Repo, state GAState, cell GACell, scenario Scenario) (string, error) {
	docs, err := sourceDocumentsForPrompt(repo, cell.SimID, scenario)
	if err != nil {
		return "", err
	}
	// Intent: Bundle the exact local sim/scenario source documents into each
	// score prompt so a provider-backed cell has enough context without parsing
	// cross-file links itself. Source: DI-gijom
	var out strings.Builder
	out.WriteString("# GA Score Cell\n\n")
	out.WriteString("Return only JSON with keys `scores`, `fitness`, and `assessment`.\n")
	out.WriteString("Do not include result identity, source metadata, code fences, or commentary.\n\n")
	out.WriteString("## Cell\n\n")
	fmt.Fprintf(&out, "- Run group ID: `%s`\n", state.RunGroupID)
	fmt.Fprintf(&out, "- Cell ID: `%s`\n", cell.CellID)
	fmt.Fprintf(&out, "- Simulation ID: `%s`\n", cell.SimID)
	fmt.Fprintf(&out, "- Scenario ID: `%s`\n", cell.ScenarioID)
	fmt.Fprintf(&out, "- Model ID: `%s`\n", cell.ModelID)
	fmt.Fprintf(&out, "- Result path: `%s`\n\n", cell.ResultPath)
	out.WriteString("## Rubric\n\n")
	out.WriteString("Score each axis from 0 to 5. Higher is better except `risk_penalty`, where 0 is low risk and 5 is severe risk.\n")
	out.WriteString("Axes: scenario_fit, promisegrid_alignment, auditability, evolution_safety, layer_boundary_clarity, failure_handling, implementation_plausibility, risk_penalty.\n\n")
	out.WriteString("## Source Documents\n\n")
	for _, doc := range docs {
		fmt.Fprintf(&out, "### `%s`\n\n```markdown\n%s\n```\n\n", doc.Path, strings.TrimSpace(doc.Text))
	}
	out.WriteString("## Required JSON Shape\n\n")
	out.WriteString(`{"scores":{"scenario_fit":0,"promisegrid_alignment":0,"auditability":0,"evolution_safety":0,"layer_boundary_clarity":0,"failure_handling":0,"implementation_plausibility":0,"risk_penalty":0},"fitness":{"raw":0,"normalized_0_100":0,"confidence_0_1":0.0},"assessment":{"rationale":"","strengths":[],"weaknesses":[],"risks":[],"open_questions":[],"authority_boundary":"Evidence only; does not settle PromiseGrid design."}}`)
	out.WriteString("\n")
	return out.String(), nil
}

func buildFitnessResult(repo Repo, state GAState, cell GACell, scenario Scenario, payload scorePayload) (FitnessResult, error) {
	parts, issues := parseResultPath(repo, repo.Abs(cell.ResultPath))
	if len(issues) > 0 {
		return FitnessResult{}, fmt.Errorf("%s", strings.Join(issues, "; "))
	}
	files, err := sourceFilesForResult(repo, cell.SimID, scenario)
	if err != nil {
		return FitnessResult{}, err
	}
	simHash, err := currentSimulationTreeHash(repo, filepath.ToSlash(filepath.Join("simulations", cell.SimID)))
	if err != nil {
		return FitnessResult{}, err
	}
	return FitnessResult{
		Schema:       resultSchemaV1,
		ResultID:     strings.Join([]string{parts.SimID, parts.ScenarioID, parts.ModelID, parts.Timestamp}, "-"),
		RunGroupID:   state.RunGroupID,
		CellID:       cell.CellID,
		SimID:        cell.SimID,
		ScenarioID:   cell.ScenarioID,
		ModelID:      cell.ModelID,
		TimestampUTC: parts.Timestamp,
		ResultPath:   cell.ResultPath,
		Runner: RunnerInfo{
			Tool:              "ga-runner",
			Provider:          cell.Provider,
			APIModel:          cell.APIModel,
			ReasoningEffort:   cell.ReasoningEffort,
			ServiceTier:       cell.ServiceTier,
			ServedServiceTier: cell.ServedServiceTier,
			RequestID:         cell.RequestID,
			ResponseID:        cell.ResponseID,
			InputTokens:       cell.InputTokens,
			CachedInputTokens: cell.CachedTokens,
			OutputTokens:      cell.OutputTokens,
			CostUSD:           cell.CostUSD,
		},
		Source: SourceInfo{
			RepoCommit:         state.RepoCommit,
			SimPath:            filepath.ToSlash(filepath.Join("simulations", cell.SimID)) + "/",
			ScenarioPath:       scenario.Path,
			RootContractPaths:  []string{"results/RUN-PROTOCOL.md", "scenarios/README.md"},
			Files:              files,
			SimulationTreeHash: simHash,
		},
		Rubric: RubricInfo{
			RubricVersion: "ga-rubric-20260519-v1",
			ScoreScale:    "0..5",
			ScoreMeanings: map[string]string{"0": "no fit or absent", "5": "strong fit", "risk_penalty": "0 low risk, 5 severe risk"},
			Axes:          []string{"scenario_fit", "promisegrid_alignment", "auditability", "evolution_safety", "layer_boundary_clarity", "failure_handling", "implementation_plausibility", "risk_penalty"},
		},
		Scores:     payload.Scores,
		Fitness:    payload.Fitness,
		Assessment: payload.Assessment,
	}, nil
}

func scenarioFromState(state GAState, scenarioID string) (Scenario, error) {
	for _, scenario := range state.ScenarioSample {
		if scenario.ScenarioID == scenarioID {
			return Scenario{ScenarioID: scenario.ScenarioID, Path: scenario.Path}, nil
		}
	}
	return Scenario{}, fmt.Errorf("scenario %s is not present in GA state", scenarioID)
}

func markGACell(cell *GACell, status string, message string) {
	cell.Status = status
	cell.ValidationMessage = message
	cell.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
}
