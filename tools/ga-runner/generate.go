package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type generateOptions struct {
	RepoRoot                 string
	RunGroupID               string
	ChildIDs                 []string
	ProviderName             string
	APIModel                 string
	ReasoningEffort          string
	ServiceTier              string
	APIKeyEnv                string
	OpenAIBaseURL            string
	Workers                  int
	RequestTimeout           time.Duration
	ProviderAttempts         int
	ProviderElapsed          time.Duration
	TextVerbosity            string
	MaxOutputTokens          int
	CostEstimateOutputTokens int
	MaxRunCostUSD            float64
	MaxChildCostUSD          float64
	InputPrice               float64
	CachedInputPrice         float64
	OutputPrice              float64
	RetryFailed              bool
	SkipFailedChildren       bool
	DryRun                   bool
	JobDir                   string
}

type childBundle struct {
	ChildID            string            `json:"child_id"`
	DesignDeltaSummary string            `json:"design_delta_summary"`
	Files              []childBundleFile `json:"files"`
}

type childBundleFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

func runGenerate(args []string, stdout io.Writer) error {
	options, err := parseGenerateOptions(args)
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
	})
	if err != nil {
		return err
	}
	return runGenerateWithProvider(context.Background(), repo, provider, options, stdout)
}

func parseGenerateOptions(args []string) (generateOptions, error) {
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	runGroupID := fs.String("run-group-id", "", "GA run group whose state contains queued child plans")
	providerName := fs.String("provider", "openai", "provider name; v1 supports openai")
	apiModel := fs.String("api-model", "", "provider API model name")
	// Intent: Child generation is more prone to reasoning-token exhaustion than
	// scoring, so default it to medium reasoning while preserving explicit
	// operator override. Source: DI-pulap
	reasoningEffort := fs.String("reasoning-effort", defaultGenerateReasoningEffort, "provider reasoning effort")
	// Intent: Default child generation to Flex and reject Priority so unattended
	// generation cannot inherit expensive processing. Source: DI-mopob
	serviceTier := fs.String("service-tier", defaultServiceTier, "provider service tier: flex or default; priority is rejected")
	apiKeyEnv := fs.String("api-key-env", "OPENAI_API_KEY", "environment variable holding provider API key")
	openAIBaseURL := fs.String("openai-base-url", "", "optional OpenAI Responses API URL override")
	// Intent: Keep raw generation serial by default while allowing bounded
	// worker pools for deliberately larger sync runs. Source: DI-juzus
	workers := fs.Int("workers", defaultGenerateWorkers, "number of concurrent provider child-generation workers")
	requestTimeout := fs.String("request-timeout", defaultRequestTimeout.String(), "per-request provider timeout as a Go duration")
	providerAttempts := fs.Int("provider-max-attempts", defaultProviderMaxAttempts, "maximum provider attempts per child")
	providerElapsed := fs.String("provider-max-elapsed", defaultProviderMaxElapsed.String(), "maximum elapsed provider retry time per child")
	// Intent: Default child-generation requests away from hard output-token caps;
	// budget estimates remain preflight-only. Source: DI-pulap
	textVerbosity := fs.String("text-verbosity", defaultTextVerbosity, "provider text verbosity: low, medium, or high")
	maxOutputTokens := fs.Int("max-output-tokens", 0, "optional hard provider output-token cap; 0 omits")
	costEstimateOutputTokens := fs.Int("cost-estimate-output-tokens", defaultGenerateCostEstimateOutputTokens, "output tokens used only for preflight cost estimates")
	maxRunCost := fs.Float64("max-run-cost-usd", 0, "stop before starting a child that would exceed this run budget; 0 disables")
	maxChildCost := fs.Float64("max-child-estimate-usd", 0, "skip children whose preflight cost estimate exceeds this amount; 0 disables")
	inputPrice := fs.Float64("cost-input-usd-per-mtok", defaultInputUSDPerMTok, "uncached input price in USD per million tokens")
	cachedInputPrice := fs.Float64("cost-cached-input-usd-per-mtok", defaultCachedInputUSDPerMTok, "cached input price in USD per million tokens")
	outputPrice := fs.Float64("cost-output-usd-per-mtok", defaultOutputUSDPerMTok, "output price in USD per million tokens")
	retryFailed := fs.Bool("retry-failed", false, "retry failed child generation")
	// Intent: Let unattended canaries preserve individual child-generation
	// failures as skipped evidence instead of aborting the whole GA cycle.
	// Source: DI-zikag
	skipFailedChildren := fs.Bool("skip-failed-children", false, "mark failed child generation skipped and keep the generate command successful")
	dryRun := fs.Bool("dry-run", false, "show selected children without writing state/trees")
	jobDir := fs.String("job-dir", "", "prompt audit directory; defaults to results/jobs/<run-group-id>")
	var childIDs stringListFlag
	fs.Var(&childIDs, "child", "child simulation ID to generate; repeatable")
	if err := fs.Parse(args); err != nil {
		return generateOptions{}, err
	}
	if *runGroupID == "" {
		return generateOptions{}, errUsage("generate: -run-group-id is required")
	}
	if *apiModel == "" && !*dryRun {
		return generateOptions{}, errUsage("generate: -api-model is required for non-dry runs")
	}
	normalizedServiceTier, err := normalizeServiceTier(*serviceTier)
	if err != nil {
		return generateOptions{}, errUsage("generate: " + err.Error())
	}
	normalizedWorkers, err := normalizeWorkers(*workers)
	if err != nil {
		return generateOptions{}, errUsage("generate: " + err.Error())
	}
	parsedRequestTimeout, err := parsePositiveDurationFlag("request-timeout", *requestTimeout)
	if err != nil {
		return generateOptions{}, errUsage("generate: " + err.Error())
	}
	if *providerAttempts < 1 {
		return generateOptions{}, errUsage("generate: provider-max-attempts must be at least 1")
	}
	if *maxOutputTokens < 0 {
		return generateOptions{}, errUsage("generate: max-output-tokens must be zero or greater")
	}
	if *costEstimateOutputTokens < 0 {
		return generateOptions{}, errUsage("generate: cost-estimate-output-tokens must be zero or greater")
	}
	parsedProviderElapsed, err := parsePositiveDurationFlag("provider-max-elapsed", *providerElapsed)
	if err != nil {
		return generateOptions{}, errUsage("generate: " + err.Error())
	}
	normalizedTextVerbosity, err := normalizeTextVerbosity(*textVerbosity)
	if err != nil {
		return generateOptions{}, errUsage("generate: " + err.Error())
	}
	return generateOptions{
		RepoRoot:                 *repoRoot,
		RunGroupID:               *runGroupID,
		ChildIDs:                 []string(childIDs),
		ProviderName:             *providerName,
		APIModel:                 *apiModel,
		ReasoningEffort:          *reasoningEffort,
		ServiceTier:              normalizedServiceTier,
		APIKeyEnv:                *apiKeyEnv,
		OpenAIBaseURL:            *openAIBaseURL,
		Workers:                  normalizedWorkers,
		RequestTimeout:           parsedRequestTimeout,
		ProviderAttempts:         *providerAttempts,
		ProviderElapsed:          parsedProviderElapsed,
		TextVerbosity:            normalizedTextVerbosity,
		MaxOutputTokens:          *maxOutputTokens,
		CostEstimateOutputTokens: *costEstimateOutputTokens,
		MaxRunCostUSD:            *maxRunCost,
		MaxChildCostUSD:          *maxChildCost,
		InputPrice:               *inputPrice,
		CachedInputPrice:         *cachedInputPrice,
		OutputPrice:              *outputPrice,
		RetryFailed:              *retryFailed,
		SkipFailedChildren:       *skipFailedChildren,
		DryRun:                   *dryRun,
		JobDir:                   *jobDir,
	}, nil
}

func runGenerateWithProvider(ctx context.Context, repo Repo, provider Provider, options generateOptions, stdout io.Writer) error {
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
	indexes := selectedGenerateChildIndexes(state, options)
	if len(indexes) == 0 {
		// Intent: Treat all-skipped child generation as a successful no-op only
		// when the caller explicitly chose skip-and-continue semantics.
		// Source: DI-zikag
		if options.SkipFailedChildren {
			return writeFormat(stdout, "processed=0 failed=0 skipped=0 state=%s\n", repo.Rel(stateFile))
		}
		return fmt.Errorf("no child plans matched generation selection")
	}
	if !options.DryRun && resetRunningGenerateChildrenForRetry(&state, indexes) {
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
		MaxCellEstimateUSD:    options.MaxChildCostUSD,
	}
	processed := 0
	failed := 0
	skipped := 0
	reservedCostUSD := 0.0
	var jobs []generateJob
	for _, index := range indexes {
		job, status, err := prepareGenerateJob(repo, &state, stateFile, jobDir, index, options, cost, &reservedCostUSD, stdout)
		if err != nil {
			return err
		}
		processed++
		if job.Ready {
			jobs = append(jobs, job)
			continue
		}
		if status == "failed" {
			// Intent: Preserve failed child-generation messages in state while
			// letting generated siblings continue through the GA cycle.
			// Source: DI-zikag
			if options.SkipFailedChildren {
				skipGAChildAfterFailure(&state.Children[index])
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
	for result := range runGenerateJobs(ctx, repo, provider, jobs, options, cost) {
		status := result.Status
		if !result.Final {
			state.Children[result.Index] = result.Child
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
			if options.SkipFailedChildren {
				skipGAChildAfterFailure(&result.Child)
				status = "skipped"
			} else {
				failed++
			}
		}
		if status == "skipped" {
			skipped++
		}
		if !options.DryRun {
			state.Children[result.Index] = result.Child
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
		return fmt.Errorf("one or more child generations failed")
	}
	return nil
}

// effectiveGenerateCostEstimateOutputTokens keeps budget estimates conservative
// for tests and direct callers that bypass CLI parsing.
//
// Intent: Preserve preflight cost controls while omitting default provider hard
// output caps that waste canary attempts. Source: DI-pulap
func effectiveGenerateCostEstimateOutputTokens(options generateOptions) int {
	if options.CostEstimateOutputTokens > 0 {
		return options.CostEstimateOutputTokens
	}
	return defaultGenerateCostEstimateOutputTokens
}

func selectedGenerateChildIndexes(state GAState, options generateOptions) []int {
	selected := map[string]bool{}
	for _, childID := range options.ChildIDs {
		selected[childID] = true
	}
	var indexes []int
	for index, child := range state.Children {
		childID := child.ID()
		if len(selected) > 0 && !selected[childID] {
			continue
		}
		if child.Status == "generated" || child.Status == "accepted" || child.Status == "culled" || child.Status == "skipped" {
			continue
		}
		if child.Status == "failed" && !options.RetryFailed {
			continue
		}
		indexes = append(indexes, index)
	}
	return indexes
}

type generateJob struct {
	Ready  bool
	Index  int
	Child  GAChild
	Prompt string
}

type generateJobResult struct {
	Index  int
	Child  GAChild
	Status string
	Final  bool
	Ack    chan error
}

func resetRunningGenerateChildrenForRetry(state *GAState, indexes []int) bool {
	changed := false
	for _, index := range indexes {
		if state.Children[index].Status != "running" {
			continue
		}
		// Intent: Treat stale `running` child-generation records as interrupted
		// prior work when a user reruns generation. Fresh `running` status is
		// written only after a worker owns the child. Source: DI-juzus
		markGAChild(&state.Children[index], "queued", "queued after incomplete prior child-generation attempt")
		changed = true
	}
	return changed
}

func prepareGenerateJob(repo Repo, state *GAState, stateFile string, jobDir string, index int, options generateOptions, cost CostConfig, reservedCostUSD *float64, stdout io.Writer) (generateJob, string, error) {
	child := &state.Children[index]
	if options.DryRun {
		if err := writeFormat(stdout, "dry-run generate: %s -> %s\n", child.ID(), child.Path); err != nil {
			child.Status = "failed"
			return generateJob{}, "failed", nil
		}
		return generateJob{}, child.Status, nil
	}
	prompt, err := buildGeneratePrompt(repo, *state, *child)
	if err != nil {
		markGAChild(child, "failed", err.Error())
		return generateJob{}, "failed", nil
	}
	// Intent: Stop before dispatching concurrent provider workers when a child
	// generation prompt would exceed per-child or reserved whole-run budget,
	// using estimate-only output tokens instead of a provider hard cap. Source:
	// DI-gijom; DI-juzus; DI-pulap
	estimate := cost.EstimatePromptCost(prompt, effectiveGenerateCostEstimateOutputTokens(options))
	if cost.MaxCellEstimateUSD > 0 && estimate.CostUSD > cost.MaxCellEstimateUSD {
		markGAChild(child, "failed", fmt.Sprintf("estimated child cost %.6f exceeds max %.6f", estimate.CostUSD, cost.MaxCellEstimateUSD))
		return generateJob{}, "failed", nil
	}
	if cost.MaxRunUSD > 0 && stateActualCostUSD(*state)+*reservedCostUSD+estimate.CostUSD > cost.MaxRunUSD {
		if err := writeFormat(stdout, "budget-stop: actual_cost_usd=%.6f reserved_estimate_usd=%.6f estimated_next_child_usd=%.6f max_run_cost_usd=%.6f\n", stateActualCostUSD(*state), *reservedCostUSD, estimate.CostUSD, cost.MaxRunUSD); err != nil {
			markGAChild(child, "failed", err.Error())
			return generateJob{}, "failed", nil
		}
		return generateJob{}, "budget-stop", nil
	}
	promptPath := filepath.Join(jobDir, child.ID()+".generate.md")
	if err := writeFile(promptPath, prompt); err != nil {
		markGAChild(child, "failed", err.Error())
		return generateJob{}, "failed", nil
	}
	child.PromptHash = hashText(prompt)
	child.Provider = options.ProviderName
	child.APIModel = options.APIModel
	child.ReasoningEffort = options.ReasoningEffort
	child.ServiceTier = options.ServiceTier
	*reservedCostUSD += estimate.CostUSD
	return generateJob{
		Ready:  true,
		Index:  index,
		Child:  *child,
		Prompt: prompt,
	}, "running", nil
}

func runGenerateJobs(ctx context.Context, repo Repo, provider Provider, jobs []generateJob, options generateOptions, cost CostConfig) <-chan generateJobResult {
	results := make(chan generateJobResult)
	if len(jobs) == 0 {
		close(results)
		return results
	}
	workers := options.Workers
	if workers < 1 {
		workers = defaultGenerateWorkers
	}
	if workers > len(jobs) {
		workers = len(jobs)
	}
	jobCh := make(chan generateJob)
	var waitGroup sync.WaitGroup
	waitGroup.Add(workers)
	for workerID := 0; workerID < workers; workerID++ {
		go func() {
			defer waitGroup.Done()
			for job := range jobCh {
				executeGenerateJob(ctx, repo, provider, job, options, cost, results)
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

func executeGenerateJob(ctx context.Context, repo Repo, provider Provider, job generateJob, options generateOptions, cost CostConfig, results chan<- generateJobResult) {
	child := job.Child
	markGAChild(&child, "running", "provider call started")
	ack := make(chan error)
	results <- generateJobResult{Index: job.Index, Child: child, Status: "running", Ack: ack}
	if err := <-ack; err != nil {
		return
	}
	callCtx := ctx
	cancel := func() {}
	if options.ProviderElapsed > 0 {
		// Intent: Bound each child-generation retry window so one slow child
		// cannot monopolize a worker for the old 30-minute timeout. Source:
		// DI-juzus
		callCtx, cancel = context.WithTimeout(ctx, options.ProviderElapsed)
	}
	defer cancel()
	response, err := provider.Generate(callCtx, ProviderRequest{
		Provider:        options.ProviderName,
		APIModel:        options.APIModel,
		ReasoningEffort: options.ReasoningEffort,
		ServiceTier:     options.ServiceTier,
		TextVerbosity:   options.TextVerbosity,
		MaxOutputTokens: options.MaxOutputTokens,
		Instructions:    "Return only valid JSON for the requested GA child file bundle. Do not include code fences or commentary.",
		Prompt:          job.Prompt,
	})
	if err != nil {
		markGAChild(&child, "failed", err.Error())
		results <- generateJobResult{Index: job.Index, Child: child, Status: "failed", Final: true}
		return
	}
	child.RequestID = response.RequestID
	child.ResponseID = response.ResponseID
	child.ServedServiceTier = response.ServiceTier
	child.UsageJSON = response.UsageJSON
	// Intent: Preserve actual usage from the generation provider so the state
	// file remains the restart and cost-review checkpoint for the GA run.
	// Source: DI-gijom
	if response.UsageJSON != "" {
		usageCost, err := cost.ParseUsage(response.UsageJSON)
		if err != nil {
			markGAChild(&child, "failed", err.Error())
			results <- generateJobResult{Index: job.Index, Child: child, Status: "failed", Final: true}
			return
		}
		child.InputTokens = usageCost.InputTokens
		child.CachedTokens = usageCost.CachedInputTokens
		child.OutputTokens = usageCost.OutputTokens
		child.CostUSD = usageCost.CostUSD
	} else if cost.BudgetEnabled() {
		markGAChild(&child, "failed", "missing provider usage metadata for cost-controlled run")
		results <- generateJobResult{Index: job.Index, Child: child, Status: "failed", Final: true}
		return
	}
	child.ResponseHash = hashText(response.Text)
	bundle, err := parseChildBundle(response.Text)
	if err != nil {
		markGAChild(&child, "failed", err.Error())
		results <- generateJobResult{Index: job.Index, Child: child, Status: "failed", Final: true}
		return
	}
	if err := writeChildBundle(repo, child, bundle); err != nil {
		markGAChild(&child, "failed", err.Error())
		results <- generateJobResult{Index: job.Index, Child: child, Status: "failed", Final: true}
		return
	}
	treeHash, err := currentSimulationTreeHash(repo, strings.TrimSuffix(normalizeRelPath(child.Path), "/"))
	if err != nil {
		markGAChild(&child, "failed", err.Error())
		results <- generateJobResult{Index: job.Index, Child: child, Status: "failed", Final: true}
		return
	}
	files, err := generatedChildFiles(repo, child)
	if err != nil {
		markGAChild(&child, "failed", err.Error())
		results <- generateJobResult{Index: job.Index, Child: child, Status: "failed", Final: true}
		return
	}
	child.DesignDeltaSummary = bundle.DesignDeltaSummary
	child.TreeHash = treeHash
	child.Files = files
	markGAChild(&child, "generated", "validated child tree")
	results <- generateJobResult{Index: job.Index, Child: child, Status: "generated", Final: true}
}

func parseChildBundle(text string) (childBundle, error) {
	clean := strings.TrimSpace(text)
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimPrefix(clean, "```")
	clean = strings.TrimSuffix(clean, "```")
	var bundle childBundle
	if err := json.Unmarshal([]byte(strings.TrimSpace(clean)), &bundle); err != nil {
		return childBundle{}, err
	}
	return bundle, nil
}

func buildGeneratePrompt(repo Repo, state GAState, child GAChild) (string, error) {
	// Intent: Generate children from complete local source bundles and in-run
	// fitness evidence so proposed sims can stand alone after materialization.
	// Source: DI-gijom
	var out strings.Builder
	out.WriteString("# GA Child Generation\n\n")
	out.WriteString("Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.\n")
	out.WriteString("Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.\n\n")
	fmt.Fprintf(&out, "- Run group ID: `%s`\n", state.RunGroupID)
	fmt.Fprintf(&out, "- Child ID: `%s`\n", child.ID())
	fmt.Fprintf(&out, "- Child path: `%s`\n", child.Path)
	fmt.Fprintf(&out, "- Operation: `%s`\n", child.Operation)
	fmt.Fprintf(&out, "- Parent IDs: `%s`\n\n", strings.Join(child.ParentIDs, ", "))
	out.WriteString("## Scenario Sample\n\n")
	for _, scenario := range state.ScenarioSample {
		fmt.Fprintf(&out, "- `%s` at `%s`\n", scenario.ScenarioID, scenario.Path)
	}
	out.WriteString("\n## Parent Source Documents\n\n")
	seen := map[string]bool{}
	for _, parentID := range child.ParentIDs {
		for _, scenarioState := range state.ScenarioSample {
			docs, err := sourceDocumentsForPrompt(repo, parentID, Scenario{ScenarioID: scenarioState.ScenarioID, Path: scenarioState.Path})
			if err != nil {
				return "", err
			}
			for _, doc := range docs {
				if seen[doc.Path] {
					continue
				}
				seen[doc.Path] = true
				fmt.Fprintf(&out, "### `%s`\n\n```markdown\n%s\n```\n\n", doc.Path, strings.TrimSpace(doc.Text))
			}
		}
	}
	out.WriteString("## Existing Fitness Evidence From This Run\n\n")
	for _, cell := range state.Cells {
		if !containsString(child.ParentIDs, cell.SimID) || cell.Status != "done" || cell.ResultPath == "" {
			continue
		}
		text, err := repo.ReadRel(cell.ResultPath)
		if err != nil {
			continue
		}
		fmt.Fprintf(&out, "### `%s`\n\n```json\n%s\n```\n\n", cell.ResultPath, strings.TrimSpace(text))
	}
	out.WriteString("## Required JSON Shape\n\n")
	fmt.Fprintf(&out, `{"child_id":%q,"design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}`, child.ID())
	out.WriteString("\n")
	return out.String(), nil
}

func writeChildBundle(repo Repo, child GAChild, bundle childBundle) error {
	// Intent: Treat the model response as a bounded transport envelope only; the
	// durable artifact is the materialized simulation tree under simulations/.
	// Source: DI-gijom
	if bundle.ChildID != child.ID() {
		return fmt.Errorf("bundle child_id %q does not match expected %q", bundle.ChildID, child.ID())
	}
	relChildPath := strings.TrimSuffix(normalizeRelPath(child.Path), "/")
	if relChildPath != filepath.ToSlash(filepath.Join("simulations", child.ID())) {
		return fmt.Errorf("child path must be simulations/%s/", child.ID())
	}
	if _, err := os.Stat(repo.Abs(relChildPath)); err == nil {
		return fmt.Errorf("child path already exists: %s", relChildPath)
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := validateChildBundle(bundle); err != nil {
		return err
	}
	tmpPath := repo.Path("simulations", "."+child.ID()+".tmp-"+fmt.Sprintf("%d", time.Now().UnixNano()))
	if err := os.MkdirAll(tmpPath, 0o755); err != nil {
		return err
	}
	for _, file := range bundle.Files {
		rel, err := childBundleRelPath(file.Path)
		if err != nil {
			return err
		}
		if err := writeFile(filepath.Join(tmpPath, rel), file.Content); err != nil {
			return err
		}
	}
	return os.Rename(tmpPath, repo.Abs(relChildPath))
}

func validateChildBundle(bundle childBundle) error {
	if strings.TrimSpace(bundle.DesignDeltaSummary) == "" {
		return fmt.Errorf("design_delta_summary is required")
	}
	seen := map[string]bool{}
	for _, file := range bundle.Files {
		rel, err := childBundleRelPath(file.Path)
		if err != nil {
			return err
		}
		seen[rel] = true
	}
	if !seen["README.md"] {
		return fmt.Errorf("child bundle must include README.md")
	}
	if !seen["QUESTION.md"] {
		return fmt.Errorf("child bundle must include QUESTION.md")
	}
	return nil
}

func childBundleRelPath(path string) (string, error) {
	clean := normalizeRelPath(path)
	if filepath.IsAbs(path) || strings.HasPrefix(clean, "..") || strings.HasPrefix(clean, "simulations/") || clean == "." {
		return "", fmt.Errorf("unsafe child bundle path %s", path)
	}
	return clean, nil
}

func generatedChildFiles(repo Repo, child GAChild) ([]SourceFile, error) {
	root := strings.TrimSuffix(normalizeRelPath(child.Path), "/")
	var files []SourceFile
	err := filepath.WalkDir(repo.Abs(root), func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		rel := repo.Rel(path)
		hash, err := sha256File(repo, rel)
		if err != nil {
			return err
		}
		files = append(files, SourceFile{Path: rel, SHA256: hash})
		return nil
	})
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})
	return files, err
}

func markGAChild(child *GAChild, status string, message string) {
	child.Status = status
	child.ValidationMessage = message
	child.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
}

func skipGAChildAfterFailure(child *GAChild) {
	message := strings.TrimSpace(child.ValidationMessage)
	if message == "" {
		message = "child generation failed"
	}
	markGAChild(child, "skipped", "skipped after child generation failure: "+message)
}

func containsString(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
