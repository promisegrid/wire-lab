package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type scoreOptions struct {
	RepoRoot                 string
	RunGroupID               string
	Target                   string
	ProviderName             string
	APIModel                 string
	ReasoningEffort          string
	OutputContract           string
	ReasoningSummary         string
	ServiceTier              string
	APIKeyEnv                string
	OpenAIBaseURL            string
	Workers                  int
	RequestTimeout           time.Duration
	ProviderAttempts         int
	ProviderElapsed          time.Duration
	Stream                   bool
	StreamIdleTimeout        time.Duration
	StreamContentStdout      bool
	TextVerbosity            string
	MaxOutputTokens          int
	CostEstimateOutputTokens int
	MaxRunCostUSD            float64
	MaxCellCostUSD           float64
	InputPrice               float64
	CachedInputPrice         float64
	OutputPrice              float64
	RetryFailed              bool
	RerunDone                bool
	SkipFailedCells          bool
	DryRun                   bool
	JobDir                   string
}

type scorePayload struct {
	Scores     FitnessScores  `json:"scores"`
	Fitness    FitnessSummary `json:"fitness"`
	Assessment Assessment     `json:"assessment"`
}

type scorePayloadEnvelope struct {
	Scores     map[string]json.RawMessage `json:"scores"`
	Assessment map[string]json.RawMessage `json:"assessment"`
}

type scorePayloadParse struct {
	Payload scorePayload
	Raw     scorePayloadEnvelope
}

type scoreResponseAttempt struct {
	Response ProviderResponse
	Payload  scorePayload
	Raw      scorePayloadEnvelope
	Usage    UsageCost
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
	provider, err := buildProvider(providerBuildOptions{
		ProviderName:        options.ProviderName,
		APIKeyEnv:           options.APIKeyEnv,
		OpenAIBaseURL:       options.OpenAIBaseURL,
		DryRun:              options.DryRun,
		RequestTimeout:      options.RequestTimeout,
		ProviderMaxAttempts: options.ProviderAttempts,
		ProviderMaxElapsed:  options.ProviderElapsed,
		Stream:              options.Stream,
		StreamIdleTimeout:   options.StreamIdleTimeout,
		StreamContentStdout: options.StreamContentStdout,
	})
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
	reasoningEffort := fs.String("reasoning-effort", defaultScoreReasoningEffort, "provider reasoning effort")
	outputContract := fs.String("output-contract", outputContractJSONSchemaStrict, "score output contract: prompt_json or json_schema_strict")
	// Intent: Raw score commands default to no reasoning summary, while the
	// terminal canary opts in so stdout can show supported reasoning-summary
	// stream content. Source: DI-vadub
	reasoningSummary := fs.String("reasoning-summary", "", "provider reasoning summary mode such as auto; empty omits")
	// Intent: Default scoring calls to Flex and reject Priority so unattended
	// scoring cannot inherit expensive processing. Source: DI-mopob
	serviceTier := fs.String("service-tier", defaultServiceTier, "provider service tier: flex or default; priority is rejected")
	apiKeyEnv := fs.String("api-key-env", "OPENAI_API_KEY", "environment variable holding provider API key")
	openAIBaseURL := fs.String("openai-base-url", "", "optional OpenAI Responses API URL override")
	// Intent: Keep raw scoring serial by default while letting canaries opt into
	// bounded concurrent provider calls. Source: DI-juzus
	workers := fs.Int("workers", defaultScoreWorkers, "number of concurrent provider scoring workers")
	requestTimeout := fs.String("request-timeout", defaultRequestTimeout.String(), "per-request provider timeout as a Go duration")
	providerAttempts := fs.Int("provider-max-attempts", defaultProviderMaxAttempts, "maximum provider attempts per cell")
	providerElapsed := fs.String("provider-max-elapsed", defaultProviderMaxElapsed.String(), "maximum elapsed provider retry time per cell")
	// Intent: Stream by default so long score calls emit provider liveness
	// evidence before a timeout or retry. Source: DI-tufud
	stream := fs.Bool("stream", defaultProviderStream, "stream OpenAI Responses API events when supported")
	streamIdleTimeout := fs.String("stream-idle-timeout", defaultStreamIdleTimeout.String(), "maximum silence between streaming events as a Go duration")
	streamContentStdout := fs.Bool("stream-content-stdout", false, "print streamed reasoning-summary and output deltas to stdout")
	// Intent: Default score requests away from hard output-token caps; budget
	// estimates stay separate from provider request shape. Source: DI-pulap
	textVerbosity := fs.String("text-verbosity", defaultTextVerbosity, "provider text verbosity: low, medium, or high")
	maxOutputTokens := fs.Int("max-output-tokens", 0, "optional hard provider output-token cap; 0 omits")
	costEstimateOutputTokens := fs.Int("cost-estimate-output-tokens", defaultScoreCostEstimateOutputTokens, "output tokens used only for preflight cost estimates")
	maxRunCost := fs.Float64("max-run-cost-usd", 0, "stop before starting a cell that would exceed this run budget; 0 disables")
	maxCellCost := fs.Float64("max-cell-estimate-usd", 0, "skip cells whose preflight cost estimate exceeds this amount; 0 disables")
	inputPrice := fs.Float64("cost-input-usd-per-mtok", defaultInputUSDPerMTok, "uncached input price in USD per million tokens")
	cachedInputPrice := fs.Float64("cost-cached-input-usd-per-mtok", defaultCachedInputUSDPerMTok, "cached input price in USD per million tokens")
	outputPrice := fs.Float64("cost-output-usd-per-mtok", defaultOutputUSDPerMTok, "output price in USD per million tokens")
	retryFailed := fs.Bool("retry-failed", false, "retry failed cells")
	rerunDone := fs.Bool("rerun-done", false, "rerun done cells")
	// Intent: Let unattended canaries preserve individual cell failures as
	// skipped evidence instead of aborting the whole GA cycle. Source: DI-zikag
	skipFailedCells := fs.Bool("skip-failed-cells", false, "mark failed cells skipped and keep the score command successful")
	dryRun := fs.Bool("dry-run", false, "show selected work without writing state/results")
	jobDir := fs.String("job-dir", "", "prompt audit directory; defaults to results/jobs/<run-group-id>")
	if err := fs.Parse(args); err != nil {
		return scoreOptions{}, err
	}
	if *runGroupID == "" {
		return scoreOptions{}, errUsage("score: -run-group-id is required")
	}
	if *target != "parents" && *target != "children" && *target != "all" {
		return scoreOptions{}, errUsage("score: -target must be parents, children, or all")
	}
	normalizedServiceTier, err := normalizeServiceTier(*serviceTier)
	if err != nil {
		return scoreOptions{}, errUsage("score: " + err.Error())
	}
	normalizedOutputContract, err := normalizeOutputContract(*outputContract)
	if err != nil {
		return scoreOptions{}, errUsage("score: " + err.Error())
	}
	normalizedWorkers, err := normalizeWorkers(*workers)
	if err != nil {
		return scoreOptions{}, errUsage("score: " + err.Error())
	}
	parsedRequestTimeout, err := parsePositiveDurationFlag("request-timeout", *requestTimeout)
	if err != nil {
		return scoreOptions{}, errUsage("score: " + err.Error())
	}
	if *providerAttempts < 1 {
		return scoreOptions{}, errUsage("score: provider-max-attempts must be at least 1")
	}
	if *maxOutputTokens < 0 {
		return scoreOptions{}, errUsage("score: max-output-tokens must be zero or greater")
	}
	if *costEstimateOutputTokens < 0 {
		return scoreOptions{}, errUsage("score: cost-estimate-output-tokens must be zero or greater")
	}
	parsedProviderElapsed, err := parsePositiveDurationFlag("provider-max-elapsed", *providerElapsed)
	if err != nil {
		return scoreOptions{}, errUsage("score: " + err.Error())
	}
	parsedStreamIdleTimeout, err := parsePositiveDurationFlag("stream-idle-timeout", *streamIdleTimeout)
	if err != nil {
		return scoreOptions{}, errUsage("score: " + err.Error())
	}
	normalizedTextVerbosity, err := normalizeTextVerbosity(*textVerbosity)
	if err != nil {
		return scoreOptions{}, errUsage("score: " + err.Error())
	}
	return scoreOptions{
		RepoRoot:                 *repoRoot,
		RunGroupID:               *runGroupID,
		Target:                   *target,
		ProviderName:             *providerName,
		APIModel:                 *apiModel,
		ReasoningEffort:          *reasoningEffort,
		OutputContract:           normalizedOutputContract,
		ReasoningSummary:         *reasoningSummary,
		ServiceTier:              normalizedServiceTier,
		APIKeyEnv:                *apiKeyEnv,
		OpenAIBaseURL:            *openAIBaseURL,
		Workers:                  normalizedWorkers,
		RequestTimeout:           parsedRequestTimeout,
		ProviderAttempts:         *providerAttempts,
		ProviderElapsed:          parsedProviderElapsed,
		Stream:                   *stream,
		StreamIdleTimeout:        parsedStreamIdleTimeout,
		StreamContentStdout:      *streamContentStdout,
		TextVerbosity:            normalizedTextVerbosity,
		MaxOutputTokens:          *maxOutputTokens,
		CostEstimateOutputTokens: *costEstimateOutputTokens,
		MaxRunCostUSD:            *maxRunCost,
		MaxCellCostUSD:           *maxCellCost,
		InputPrice:               *inputPrice,
		CachedInputPrice:         *cachedInputPrice,
		OutputPrice:              *outputPrice,
		RetryFailed:              *retryFailed,
		RerunDone:                *rerunDone,
		SkipFailedCells:          *skipFailedCells,
		DryRun:                   *dryRun,
		JobDir:                   *jobDir,
	}, nil
}

func runScoreWithProvider(ctx context.Context, repo Repo, provider Provider, options scoreOptions, stdout io.Writer) error {
	serviceTier, err := normalizeServiceTier(options.ServiceTier)
	if err != nil {
		return err
	}
	options.ServiceTier = serviceTier
	textVerbosity, err := normalizeTextVerbosity(options.TextVerbosity)
	if err != nil {
		return err
	}
	options.TextVerbosity = textVerbosity
	outputContract, err := normalizeOutputContract(options.OutputContract)
	if err != nil {
		return err
	}
	options.OutputContract = outputContract
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
		// Intent: A canary may have no child score cells when all child generation
		// plans were skipped; treat that as a successful no-op only when the
		// caller explicitly asked to skip failed cells. Source: DI-zikag
		if options.SkipFailedCells {
			return writeFormat(stdout, "processed=0 failed=0 skipped=0 state=%s\n", repo.Rel(stateFile))
		}
		return fmt.Errorf("no score cells matched target %s", options.Target)
	}
	if !options.DryRun && resetRunningScoreCellsForRetry(&state, cells) {
		state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		if err := writeGAStateAtomic(stateFile, state); err != nil {
			return err
		}
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
	skipped := 0
	reservedCostUSD := 0.0
	var jobs []scoreJob
	for _, index := range cells {
		job, status, err := prepareScoreJob(repo, &state, stateFile, jobDir, index, options, cost, &reservedCostUSD, stdout)
		if err != nil {
			return err
		}
		processed++
		if job.Ready {
			jobs = append(jobs, job)
			continue
		}
		if status == "failed" {
			// Intent: Preserve the failure message in state while allowing the
			// terminal canary to continue into child generation and child scoring.
			// Source: DI-zikag
			if options.SkipFailedCells {
				skipGACellAfterFailure(&state.Cells[index])
				status = "skipped"
			} else {
				failed++
			}
		}
		if status == "skipped" {
			skipped++
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
	for result := range runScoreJobs(ctx, repo, provider, jobs, state, options, cost) {
		status := result.Status
		if !result.Final {
			state.Cells[result.Index] = result.Cell
			state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
			err := writeGAStateAtomic(stateFile, state)
			if result.Ack != nil {
				result.Ack <- err
			}
			if err != nil {
				return err
			}
			continue
		}
		if status == "failed" {
			if options.SkipFailedCells {
				skipGACellAfterFailure(&result.Cell)
				status = "skipped"
			} else {
				failed++
			}
		}
		if status == "skipped" {
			skipped++
		}
		if !options.DryRun {
			state.Cells[result.Index] = result.Cell
			state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
			if err := writeGAStateAtomic(stateFile, state); err != nil {
				return err
			}
		}
	}
	if err := writeFormat(stdout, "processed=%d failed=%d skipped=%d state=%s\n", processed, failed, skipped, repo.Rel(stateFile)); err != nil {
		return err
	}
	if failed > 0 {
		return fmt.Errorf("one or more score cells failed")
	}
	return nil
}

// effectiveScoreCostEstimateOutputTokens keeps old direct-call tests and
// hand-built options conservative even when only the CLI default path populated
// the estimate-only field.
//
// Intent: Preserve budget reservation without sending a provider hard cap.
// Source: DI-pulap
func effectiveScoreCostEstimateOutputTokens(options scoreOptions) int {
	if options.CostEstimateOutputTokens > 0 {
		return options.CostEstimateOutputTokens
	}
	return defaultScoreCostEstimateOutputTokens
}

// Intent: Let audit-first rubric-v2 rescoring reuse each cell's original API
// model when the operator omits `-api-model`, so rubric changes stay isolated
// from model drift by default. Source: DI-roruj
func effectiveScoreAPIModel(options scoreOptions, cell GACell) (string, error) {
	if strings.TrimSpace(options.APIModel) != "" {
		return strings.TrimSpace(options.APIModel), nil
	}
	if strings.TrimSpace(cell.APIModel) != "" {
		return strings.TrimSpace(cell.APIModel), nil
	}
	derived := deriveAPIModelFromModelID(cell.ModelID)
	if derived != "" {
		return derived, nil
	}
	return "", fmt.Errorf("score: -api-model is required when selected cells do not already carry api_model")
}

func selectedScoreCellIndexes(state GAState, options scoreOptions) []int {
	parentIDs := map[string]bool{}
	for _, parent := range state.Parents {
		parentIDs[parent.SimID] = true
	}
	childIDs := map[string]bool{}
	for _, child := range state.Children {
		// Intent: Score only materialized child simulation trees so skipped or
		// failed child-generation plans do not become missing-source failures.
		// Source: DI-zikag
		if child.Status == "generated" || child.Status == "accepted" || child.Status == "scored" {
			childIDs[child.ID()] = true
		}
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

type scoreJob struct {
	Ready    bool
	Index    int
	Cell     GACell
	Scenario Scenario
	Prompt   string
}

type scoreJobResult struct {
	Index  int
	Cell   GACell
	Status string
	Final  bool
	Ack    chan error
}

func resetRunningScoreCellsForRetry(state *GAState, indexes []int) bool {
	changed := false
	for _, index := range indexes {
		if state.Cells[index].Status != "running" {
			continue
		}
		// Intent: Treat stale `running` cells as interrupted prior work when a
		// user explicitly reruns the score command for the same state. Fresh
		// worker-owned `running` status is written only after a worker starts.
		// Source: DI-juzus
		markGACell(&state.Cells[index], "queued", "queued after incomplete prior score attempt")
		changed = true
	}
	return changed
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

func prepareScoreJob(repo Repo, state *GAState, stateFile string, jobDir string, index int, options scoreOptions, cost CostConfig, reservedCostUSD *float64, stdout io.Writer) (scoreJob, string, error) {
	cell := &state.Cells[index]
	scenario, err := scenarioFromState(*state, cell.ScenarioID)
	if err != nil {
		markGACell(cell, "failed", err.Error())
		return scoreJob{}, "failed", nil
	}
	if options.DryRun {
		if err := writeFormat(stdout, "dry-run score: %s -> %s\n", cell.CellID, cell.ResultPath); err != nil {
			markGACell(cell, "failed", err.Error())
			return scoreJob{}, "failed", nil
		}
		return scoreJob{}, cell.Status, nil
	}
	prompt, err := buildScorePrompt(repo, *state, *cell, scenario)
	if err != nil {
		markGACell(cell, "failed", err.Error())
		return scoreJob{}, "failed", nil
	}
	// Intent: Estimate and reserve token cost before dispatching concurrent
	// provider calls, so worker pools cannot oversubscribe the user-approved run
	// budget without sending a hard provider cap. Source: DI-gijom; DI-juzus;
	// DI-pulap
	estimate := cost.EstimatePromptCost(prompt, effectiveScoreCostEstimateOutputTokens(options))
	if cost.MaxCellEstimateUSD > 0 && estimate.CostUSD > cost.MaxCellEstimateUSD {
		markGACell(cell, "skipped", fmt.Sprintf("estimated cell cost %.6f exceeds max %.6f", estimate.CostUSD, cost.MaxCellEstimateUSD))
		return scoreJob{}, "skipped", nil
	}
	if cost.MaxRunUSD > 0 && stateActualCostUSD(*state)+*reservedCostUSD+estimate.CostUSD > cost.MaxRunUSD {
		if err := writeFormat(stdout, "budget-stop: actual_cost_usd=%.6f reserved_estimate_usd=%.6f estimated_next_cell_usd=%.6f max_run_cost_usd=%.6f\n", stateActualCostUSD(*state), *reservedCostUSD, estimate.CostUSD, cost.MaxRunUSD); err != nil {
			markGACell(cell, "failed", err.Error())
			return scoreJob{}, "failed", nil
		}
		return scoreJob{}, "budget-stop", nil
	}
	if err := writeFile(filepath.Join(jobDir, cell.CellID+".score.md"), prompt); err != nil {
		markGACell(cell, "failed", err.Error())
		return scoreJob{}, "failed", nil
	}
	if issues := validateResultFile(repo, repo.Abs(cell.ResultPath)); len(issues) == 0 {
		markGACell(cell, "done", "existing valid result")
		return scoreJob{}, "done", nil
	}
	apiModel, err := effectiveScoreAPIModel(options, *cell)
	if err != nil {
		markGACell(cell, "failed", err.Error())
		return scoreJob{}, "failed", nil
	}
	cell.Attempts++
	cell.Provider = options.ProviderName
	cell.APIModel = apiModel
	cell.ReasoningEffort = options.ReasoningEffort
	cell.OutputContract = options.OutputContract
	cell.ServiceTier = options.ServiceTier
	*reservedCostUSD += estimate.CostUSD
	return scoreJob{
		Ready:    true,
		Index:    index,
		Cell:     *cell,
		Scenario: scenario,
		Prompt:   prompt,
	}, "running", nil
}

func runScoreJobs(ctx context.Context, repo Repo, provider Provider, jobs []scoreJob, state GAState, options scoreOptions, cost CostConfig) <-chan scoreJobResult {
	results := make(chan scoreJobResult)
	if len(jobs) == 0 {
		close(results)
		return results
	}
	workers := options.Workers
	if workers < 1 {
		workers = defaultScoreWorkers
	}
	if workers > len(jobs) {
		workers = len(jobs)
	}
	jobCh := make(chan scoreJob)
	var waitGroup sync.WaitGroup
	waitGroup.Add(workers)
	for workerID := 0; workerID < workers; workerID++ {
		go func() {
			defer waitGroup.Done()
			for job := range jobCh {
				executeScoreJob(ctx, repo, provider, job, state, options, cost, results)
			}
		}()
	}
	go func() {
		for _, job := range jobs {
			jobCh <- job
		}
		close(jobCh)
		waitGroup.Wait()
		close(results)
	}()
	return results
}

func executeScoreJob(ctx context.Context, repo Repo, provider Provider, job scoreJob, state GAState, options scoreOptions, cost CostConfig, results chan<- scoreJobResult) {
	cell := job.Cell
	markGACell(&cell, "running", "provider call started")
	ack := make(chan error)
	results <- scoreJobResult{Index: job.Index, Cell: cell, Status: "running", Ack: ack}
	if err := <-ack; err != nil {
		return
	}
	callCtx := ctx
	cancel := func() {}
	if options.ProviderElapsed > 0 {
		// Intent: Bound the whole provider retry window per cell so a worker
		// slot cannot sit indefinitely even if the provider ignores client-side
		// request timeouts. Source: DI-juzus
		callCtx, cancel = context.WithTimeout(ctx, options.ProviderElapsed)
	}
	defer cancel()
	// Intent: Keep the scorer contract strict, but spend one narrow follow-up
	// provider call when a JSON-shaped score response omits required score
	// fields such as `promise_vocabulary`, `simplicity_durability`, or the PT
	// gate. This prevents repeated paid failures like the `SIM-suzuf`
	// focused-slice regression without inventing missing judgments locally.
	// Source: DI-kibuf; DI-movur
	attempt, err := executeScoreAttempt(callCtx, provider, options, cell, job.Prompt, cost)
	if err != nil {
		markGACell(&cell, "failed", err.Error())
		results <- scoreJobResult{Index: job.Index, Cell: cell, Status: "failed", Final: true}
		return
	}
	if missingAxes := missingRequiredScoreAxes(resultSchemaV4, attempt.Raw); len(missingAxes) > 0 {
		// Intent: Fail strict structured-output responses that omit required axes
		// instead of converting absent provider scores into local zero values.
		// Source: DI-vonot
		if effectiveScoreOutputContract(options, cell) != outputContractPromptJSON {
			markGACell(&cell, "failed", fmt.Sprintf("structured output response missing required score fields: %s", strings.Join(missingAxes, ", ")))
			results <- scoreJobResult{Index: job.Index, Cell: cell, Status: "failed", Final: true}
			return
		}
		cell.Attempts++
		retryPrompt := buildScoreSchemaCorrectionPrompt(job.Prompt, missingAxes)
		retryAttempt, retryErr := executeScoreAttempt(callCtx, provider, options, cell, retryPrompt, cost)
		if retryErr != nil {
			markGACell(&cell, "failed", retryErr.Error())
			results <- scoreJobResult{Index: job.Index, Cell: cell, Status: "failed", Final: true}
			return
		}
		attempt = mergeScoreAttempts(attempt, retryAttempt)
		if missingAxes = missingRequiredScoreAxes(resultSchemaV4, attempt.Raw); len(missingAxes) > 0 {
			markGACell(&cell, "failed", fmt.Sprintf("schema-correction retry still missing required score fields: %s", strings.Join(missingAxes, ", ")))
			results <- scoreJobResult{Index: job.Index, Cell: cell, Status: "failed", Final: true}
			return
		}
	}
	applyScoreAttemptToCell(&cell, attempt)
	result, err := buildFitnessResult(repo, state, cell, job.Scenario, attempt.Payload)
	if err != nil {
		markGACell(&cell, "failed", err.Error())
		results <- scoreJobResult{Index: job.Index, Cell: cell, Status: "failed", Final: true}
		return
	}
	if err := writeFitnessResultAtomic(repo.Abs(cell.ResultPath), result); err != nil {
		markGACell(&cell, "failed", err.Error())
		results <- scoreJobResult{Index: job.Index, Cell: cell, Status: "failed", Final: true}
		return
	}
	if issues := validateResultFile(repo, repo.Abs(cell.ResultPath)); len(issues) > 0 {
		markGACell(&cell, "failed", strings.Join(issues, "; "))
		results <- scoreJobResult{Index: job.Index, Cell: cell, Status: "failed", Final: true}
		return
	}
	markGACell(&cell, "done", "validated result")
	results <- scoreJobResult{Index: job.Index, Cell: cell, Status: "done", Final: true}
}

func executeScoreAttempt(ctx context.Context, provider Provider, options scoreOptions, cell GACell, prompt string, cost CostConfig) (scoreResponseAttempt, error) {
	response, err := provider.Generate(ctx, ProviderRequest{
		Provider:         options.ProviderName,
		APIModel:         cell.APIModel,
		ReasoningEffort:  cell.ReasoningEffort,
		ReasoningSummary: options.ReasoningSummary,
		ServiceTier:      cell.ServiceTier,
		TextVerbosity:    options.TextVerbosity,
		MaxOutputTokens:  options.MaxOutputTokens,
		OutputContract:   effectiveScoreOutputContract(options, cell),
		OutputSchemaName: "ga_score_payload_v4",
		OutputSchema:     scorePayloadJSONSchema(),
		Instructions:     "Return only valid JSON for the requested GA score payload. Do not include code fences or commentary.",
		Prompt:           prompt,
	})
	if err != nil {
		return scoreResponseAttempt{}, err
	}
	usageCost, err := parseScoreUsageCost(cost, response.UsageJSON)
	if err != nil {
		return scoreResponseAttempt{}, err
	}
	parsed, err := parseScorePayload(response.Text)
	if err != nil {
		return scoreResponseAttempt{}, err
	}
	return scoreResponseAttempt{
		Response: response,
		Payload:  parsed.Payload,
		Raw:      parsed.Raw,
		Usage:    usageCost,
	}, nil
}

func parseScoreUsageCost(cost CostConfig, usageJSON string) (UsageCost, error) {
	// Intent: Record measured provider usage for every score attempt so
	// schema-correction retries count toward state cost instead of vanishing from
	// the run ledger. Source: DI-kibuf
	if usageJSON == "" {
		if cost.BudgetEnabled() {
			return UsageCost{}, fmt.Errorf("missing provider usage metadata for cost-controlled run")
		}
		return UsageCost{}, nil
	}
	return cost.ParseUsage(usageJSON)
}

func applyScoreAttemptToCell(cell *GACell, attempt scoreResponseAttempt) {
	cell.RequestID = attempt.Response.RequestID
	cell.ResponseID = attempt.Response.ResponseID
	cell.ServedServiceTier = attempt.Response.ServiceTier
	cell.UsageJSON = attempt.Response.UsageJSON
	cell.InputTokens = attempt.Usage.InputTokens
	cell.CachedTokens = attempt.Usage.CachedInputTokens
	cell.OutputTokens = attempt.Usage.OutputTokens
	cell.CostUSD = attempt.Usage.CostUSD
}

func effectiveScoreOutputContract(options scoreOptions, cell GACell) string {
	if normalized, err := normalizeOutputContract(options.OutputContract); err == nil && normalized != "" {
		return normalized
	}
	if normalized, err := normalizeOutputContract(cell.OutputContract); err == nil && normalized != "" {
		return normalized
	}
	return outputContractJSONSchemaStrict
}

func mergeScoreAttempts(first scoreResponseAttempt, second scoreResponseAttempt) scoreResponseAttempt {
	second.Usage.InputTokens += first.Usage.InputTokens
	second.Usage.CachedInputTokens += first.Usage.CachedInputTokens
	second.Usage.OutputTokens += first.Usage.OutputTokens
	second.Usage.CostUSD += first.Usage.CostUSD
	return second
}

func parseScorePayload(text string) (scorePayloadParse, error) {
	clean := strings.TrimSpace(text)
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimPrefix(clean, "```")
	clean = strings.TrimSuffix(clean, "```")
	clean = strings.TrimSpace(clean)
	var payload scorePayload
	if err := json.Unmarshal([]byte(clean), &payload); err != nil {
		return scorePayloadParse{}, err
	}
	var raw scorePayloadEnvelope
	if err := json.Unmarshal([]byte(clean), &raw); err != nil {
		return scorePayloadParse{}, err
	}
	return scorePayloadParse{Payload: payload, Raw: raw}, nil
}

func missingRequiredScoreAxes(schema string, raw scorePayloadEnvelope) []string {
	if schema != resultSchemaV4 {
		return nil
	}
	var missing []string
	for _, field := range []string{
		"promise_vocabulary",
		"simplicity_durability",
		"envelope_discipline",
		"kernel_implementation_promises",
		"app_protocol_promise_semantics",
	} {
		if _, ok := raw.Scores[field]; !ok {
			missing = append(missing, "scores."+field)
		}
	}
	if _, ok := raw.Assessment["pt_gate"]; !ok {
		missing = append(missing, "assessment.pt_gate")
	}
	sort.Strings(missing)
	return missing
}

func buildScoreSchemaCorrectionPrompt(prompt string, missingAxes []string) string {
	var out strings.Builder
	out.WriteString(prompt)
	out.WriteString("\n## Schema Correction\n\n")
	out.WriteString("Your previous JSON response was invalid because it omitted required fields for the rubric-v4 layer-aware Promise Theory contract.\n")
	fmt.Fprintf(&out, "Missing required fields: `%s`.\n", strings.Join(missingAxes, "`, `"))
	out.WriteString("Return the full JSON object again, with all 13 `scores` axes present exactly once and a complete `assessment.pt_gate` object.\n")
	out.WriteString("A response missing any required score field or `assessment.pt_gate` is invalid.\n")
	return out.String()
}

func scorePayloadJSONSchema() map[string]interface{} {
	stringArraySchema := map[string]interface{}{
		"type":  "array",
		"items": map[string]interface{}{"type": "string"},
	}
	ptRuleAssessmentSchema := map[string]interface{}{
		"type":                 "object",
		"additionalProperties": false,
		"required":             []string{"status", "note"},
		"properties": map[string]interface{}{
			"status": map[string]interface{}{
				"type": "string",
				"enum": []string{ptRuleStatusPass, ptRuleStatusWarning, ptRuleStatusFail},
			},
			"note": map[string]interface{}{"type": "string"},
		},
	}
	ptGateSchema := map[string]interface{}{
		"type":                 "object",
		"additionalProperties": false,
		"required": []string{
			"status",
			"autonomous_agents",
			"scoped_intent",
			"no_promises_for_others",
			"no_guaranteed_outcomes",
			"local_trust_assessment",
			"accept_use_not_obligation",
			"violations",
		},
		"properties": map[string]interface{}{
			"status": map[string]interface{}{
				"type": "string",
				"enum": []string{ptGateStatusClean, ptGateStatusReframeNeeded, ptGateStatusInvalid},
			},
			"autonomous_agents":         ptRuleAssessmentSchema,
			"scoped_intent":             ptRuleAssessmentSchema,
			"no_promises_for_others":    ptRuleAssessmentSchema,
			"no_guaranteed_outcomes":    ptRuleAssessmentSchema,
			"local_trust_assessment":    ptRuleAssessmentSchema,
			"accept_use_not_obligation": ptRuleAssessmentSchema,
			"violations":                stringArraySchema,
		},
	}
	return map[string]interface{}{
		"type":                 "object",
		"additionalProperties": false,
		"required":             []string{"scores", "fitness", "assessment"},
		"properties": map[string]interface{}{
			"scores": map[string]interface{}{
				"type":                 "object",
				"additionalProperties": false,
				"required": []string{
					"scenario_fit",
					"promisegrid_alignment",
					"auditability",
					"evolution_safety",
					"layer_boundary_clarity",
					"failure_handling",
					"implementation_plausibility",
					"promise_vocabulary",
					"simplicity_durability",
					"envelope_discipline",
					"kernel_implementation_promises",
					"app_protocol_promise_semantics",
					"risk_penalty",
				},
				"properties": map[string]interface{}{
					"scenario_fit":                   boundedIntegerSchema(),
					"promisegrid_alignment":          boundedIntegerSchema(),
					"auditability":                   boundedIntegerSchema(),
					"evolution_safety":               boundedIntegerSchema(),
					"layer_boundary_clarity":         boundedIntegerSchema(),
					"failure_handling":               boundedIntegerSchema(),
					"implementation_plausibility":    boundedIntegerSchema(),
					"promise_vocabulary":             boundedIntegerSchema(),
					"simplicity_durability":          boundedIntegerSchema(),
					"envelope_discipline":            boundedIntegerSchema(),
					"kernel_implementation_promises": boundedIntegerSchema(),
					"app_protocol_promise_semantics": boundedIntegerSchema(),
					"risk_penalty":                   boundedIntegerSchema(),
				},
			},
			"fitness": map[string]interface{}{
				"type":                 "object",
				"additionalProperties": false,
				"required":             []string{"raw", "normalized_0_100", "confidence_0_1"},
				"properties": map[string]interface{}{
					"raw":              map[string]interface{}{"type": "number"},
					"normalized_0_100": map[string]interface{}{"type": "number"},
					"confidence_0_1":   map[string]interface{}{"type": "number", "minimum": 0, "maximum": 1},
				},
			},
			"assessment": map[string]interface{}{
				"type":                 "object",
				"additionalProperties": false,
				"required": []string{
					"rationale",
					"strengths",
					"weaknesses",
					"risks",
					"open_questions",
					"authority_boundary",
					"pt_gate",
				},
				"properties": map[string]interface{}{
					"rationale":          map[string]interface{}{"type": "string"},
					"strengths":          stringArraySchema,
					"weaknesses":         stringArraySchema,
					"risks":              stringArraySchema,
					"open_questions":     stringArraySchema,
					"authority_boundary": map[string]interface{}{"type": "string"},
					"pt_gate":            ptGateSchema,
				},
			},
		},
	}
}

func boundedIntegerSchema() map[string]interface{} {
	return map[string]interface{}{
		"type":    "integer",
		"minimum": 0,
		"maximum": 5,
	}
}

func buildScorePrompt(repo Repo, state GAState, cell GACell, scenario Scenario) (string, error) {
	docs, err := sourceDocumentsForPrompt(repo, state, cell.SimID, scenario)
	if err != nil {
		return "", err
	}
	// Intent: Bundle the exact local sim/scenario source documents into each
	// score prompt so a provider-backed cell has enough context without parsing
	// cross-file links itself. Source: DI-gijom
	//
	// Intent: Ask the scorer for the promise-first axes directly in the prompt
	// so new evidence penalizes claim-card drift and rewards small, durable
	// pCID-specific promises. Source: DI-roruj
	//
	// Intent: Make Burgess-grounded Promise Theory fundamentals explicit in the
	// scorer contract instead of relying on ambient model intuition about local
	// trust and autonomous agents. Source: DI-movur
	//
	// Intent: Score each candidate at its claimed protocol layer so an envelope
	// design is not penalized for leaving higher-layer promise accounting inside
	// the protocol payload. A signed envelope can itself be a scoped promise that
	// the payload is shaped according to the protocol specification named by the
	// pCID. Source: DI-pozom
	//
	// Intent: Penalize rejected hadit/jogoh-style base-envelope selector stacks
	// without forbidding pCID-selected higher-layer payload protocols from owning
	// their own evidence, freeze, transfer, or capability-token records. Source:
	// DI-kafiz
	//
	// Intent: Make Rubric V4 layer-aware in the active scorer by separating
	// envelope discipline, kernel implementation promises, and app-protocol
	// promise semantics instead of forcing one generic PromiseGrid alignment
	// score to carry all layer-specific concerns. Source: DI-ripuz
	//
	// Intent: Keep the scorer from rewarding conventional RPC, service-registry,
	// capability-table, or kernel-authority designs merely because they use
	// promise-shaped vocabulary. App/kernel interaction should score well when
	// apps make local pCID handling promises and kernels record local delivery
	// observations, not when the kernel claims authority over global service
	// capability facts. Source: DI-sitim
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
	out.WriteString("Axes: scenario_fit, promisegrid_alignment, auditability, evolution_safety, layer_boundary_clarity, failure_handling, implementation_plausibility, promise_vocabulary, simplicity_durability, envelope_discipline, kernel_implementation_promises, app_protocol_promise_semantics, risk_penalty.\n")
	out.WriteString("- Score the candidate at its claimed layer. Do not penalize an envelope-layer design for leaving promise accounting, storage semantics, computation semantics, or application-specific trust updates to the payload protocol when the layer boundary is explicit.\n")
	out.WriteString("- `promise_vocabulary`: reward Promise-Theory-correct, promise-first wording at the candidate's layer. For envelope sims, reward signed pCID-specific promises such as \"Alice promises these payload bytes are shaped according to the protocol specification named by this pCID.\" Penalize normative claim cards, generic claim headers/cards, conformance bundles, generic profile support claims, port claims, universal statement capsules, central trust-ledger framing, RPC dispatcher framing, service-registry framing, capability-table framing, dispatch authorization, and claims that a kernel globally knows or certifies another agent's abilities.\n")
	out.WriteString("- `simplicity_durability`: reward small, explicit, deterministic, 100-year durable artifacts that fit small devices. Penalize generic maps, cards, ledgers, bundles, mode matrices, capability maps, service catalogs, feature-shopping wrappers, and base-envelope selector-shopping stacks such as `env_pCID`/`sig_pCID`/`payload_pCID`.\n")
	out.WriteString("- `envelope_discipline`: reward alignment with `DN-jotob` and the current envelope direction: CBOR `grid([42(pCID), payload, ...])`, `pCID` as Protocol CID, protocol-owned slot roles after slot 0, local unknown-pCID behavior, and no universal proof-slot overreach.\n")
	out.WriteString("- `kernel_implementation_promises`: reward explicit local kernel implementation promises: apps promise a local kernel which pCIDs they will receive or handle, kernels promise best-effort delivery and local observation records, host assumptions are separated from promises, unsupported pCIDs/features are explicit, pCID adapter mappings are local promises, namespace/reference behavior is voluntary, and the kernel is not treated as an RPC authority, service registry, capability registry, permission issuer, or conformance judge.\n")
	out.WriteString("- `app_protocol_promise_semantics`: reward higher-layer/app protocols that model storage, computation, send/receive, reciprocal promises, selective sending, promise-as-capability-token behavior, make/break evidence, and local trust updates without command, request/response service, permission, policy-enforcement, or conformance-authority framing.\n")
	out.WriteString("- Treat phrases such as \"the kernel knows this app supports X\", \"registered service capability\", \"authorized dispatch\", or \"capability table\" as weak or invalid unless the design clearly reframes them as local promises made by specific agents plus local observations of kept, broken, refused, or timed-out promises.\n")
	out.WriteString("- Do not treat higher-layer pCID-owned payload protocols as selector shopping merely because they define their own signed refusal records, exact-byte observation evidence, freeze successor records, transfer semantics, or capability-as-promise-token payload behavior.\n")
	out.WriteString("- `scenario_fit` may be lower when a scenario asks for higher-layer behavior the candidate intentionally delegates to the payload protocol, but that delegation is not by itself a PT-gate violation when the layer boundary and local trust boundary are clear.\n")
	out.WriteString("- The runner recomputes `fitness.raw` and `fitness.normalized_0_100` from your axis scores with normal weighting. Use `fitness.confidence_0_1` for your confidence.\n\n")
	out.WriteString("## Promise Theory Fundamentals\n\n")
	out.WriteString("Apply these Mark Burgess reference notes while scoring:\n")
	for _, rule := range promiseTheoryRulesV1 {
		fmt.Fprintf(&out, "- %s\n", rule)
	}
	out.WriteString("Reference notes: Mark Burgess, *In Search of Certainty*; *Promise Theory: Principles and Applications*; *Thinking in Promises*.\n\n")
	out.WriteString("## PT Gate\n\n")
	out.WriteString("Classify the design as exactly one of `pt_clean`, `pt_reframe_needed`, or `pt_invalid`.\n")
	out.WriteString("- `pt_clean`: promise-first and locally trustworthy enough to compete normally.\n")
	out.WriteString("- `pt_reframe_needed`: technically interesting but drifts into non-PT framing, including RPC-like service dispatch, service-registry vocabulary, capability-table vocabulary, or kernel-known-support claims that could be repaired as local promises and observations; it may survive only as a question-home or rework candidate.\n")
	out.WriteString("- `pt_invalid`: relies on authority, imposition, global trust, policy enforcement, conformance authority, permission authority, or RPC-style command semantics as a load-bearing design premise; it cannot be promotable.\n")
	out.WriteString("Complete every PT rule check in `assessment.pt_gate`, and explain violations in `assessment.pt_gate.violations`.\n\n")
	out.WriteString("Required score-axis checklist: `scenario_fit`, `promisegrid_alignment`, `auditability`, `evolution_safety`, `layer_boundary_clarity`, `failure_handling`, `implementation_plausibility`, `promise_vocabulary`, `simplicity_durability`, `envelope_discipline`, `kernel_implementation_promises`, `app_protocol_promise_semantics`, `risk_penalty`.\n")
	out.WriteString("A response missing any required `scores` axis or `assessment.pt_gate` is invalid.\n\n")
	out.WriteString("## Source Documents\n\n")
	for _, doc := range docs {
		fmt.Fprintf(&out, "### `%s`\n\n```markdown\n%s\n```\n\n", doc.Path, strings.TrimSpace(doc.Text))
	}
	out.WriteString("## Required JSON Shape\n\n")
	out.WriteString(`{"scores":{"scenario_fit":0,"promisegrid_alignment":0,"auditability":0,"evolution_safety":0,"layer_boundary_clarity":0,"failure_handling":0,"implementation_plausibility":0,"promise_vocabulary":0,"simplicity_durability":0,"envelope_discipline":0,"kernel_implementation_promises":0,"app_protocol_promise_semantics":0,"risk_penalty":0},"fitness":{"raw":0,"normalized_0_100":0,"confidence_0_1":0.0},"assessment":{"rationale":"","strengths":[],"weaknesses":[],"risks":[],"open_questions":[],"authority_boundary":"Evidence only; does not settle PromiseGrid design.","pt_gate":{"status":"pt_clean","autonomous_agents":{"status":"pass","note":""},"scoped_intent":{"status":"pass","note":""},"no_promises_for_others":{"status":"pass","note":""},"no_guaranteed_outcomes":{"status":"pass","note":""},"local_trust_assessment":{"status":"pass","note":""},"accept_use_not_obligation":{"status":"pass","note":""},"violations":[]}}}`)
	out.WriteString("\n")
	return out.String(), nil
}

// Intent: Persist only the current score/result envelope for new score runs
// while the runner, not the provider, recomputes deterministic fitness from the
// returned axis scores. Source: DI-roruj; DI-movur; DI-ripuz
func buildFitnessResult(repo Repo, state GAState, cell GACell, scenario Scenario, payload scorePayload) (FitnessResult, error) {
	parts, issues := parseResultPath(repo, repo.Abs(cell.ResultPath))
	if len(issues) > 0 {
		return FitnessResult{}, fmt.Errorf("%s", strings.Join(issues, "; "))
	}
	files, err := sourceFilesForResult(repo, state, cell.SimID, scenario)
	if err != nil {
		return FitnessResult{}, err
	}
	simPath := simulationPathForState(state, cell.SimID)
	simHash, err := currentSimulationTreeHash(repo, simPath)
	if err != nil {
		return FitnessResult{}, err
	}
	scores := applyPTGateScorePolicy(payload.Scores, payload.Assessment.PTGate)
	return FitnessResult{
		Schema:       resultSchemaV4,
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
			OutputContract:    cell.OutputContract,
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
			SimPath:            strings.TrimSuffix(normalizeRelPath(simPath), "/") + "/",
			ScenarioPath:       scenario.Path,
			RootContractPaths:  []string{"results/RUN-PROTOCOL.md", "scenarios/README.md"},
			Files:              files,
			SimulationTreeHash: simHash,
		},
		Rubric: RubricInfo{
			RubricVersion:      rubricVersionForSchema(resultSchemaV4),
			ScoreScale:         "0..5",
			ScoreMeanings:      rubricScoreMeaningsForSchema(resultSchemaV4),
			Axes:               rubricAxesForSchema(resultSchemaV4),
			PromiseTheoryRules: rubricPromiseTheoryRulesForSchema(resultSchemaV4),
			PromiseTheoryRefs:  rubricPromiseTheoryReferencesForSchema(resultSchemaV4),
		},
		Scores:     scores,
		Fitness:    deterministicFitnessSummary(resultSchemaV4, scores, payload.Fitness.Confidence0To1),
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

func skipGACellAfterFailure(cell *GACell) {
	message := strings.TrimSpace(cell.ValidationMessage)
	if message == "" {
		message = "cell failed"
	}
	markGACell(cell, "skipped", "skipped after cell failure: "+message)
}
