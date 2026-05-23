package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// PopulationSim describes one committed/tracked simulation available to GA
// selection.
//
// Intent: Keep ordinary GA population scans rooted in git-tracked simulation
// specimens, while generated untracked children stay pending until accepted.
// Source: DI-bagih
type PopulationSim struct {
	SimID    string
	Path     string
	Files    []string
	TreeHash string
}

func runInit(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	dryRun := fs.Bool("dry-run", false, "print tracked simulation population without writing GA state")
	modelID := fs.String("model", "", "model ID for read-only generation planning")
	runGroupID := fs.String("run-group-id", "", "run group ID for read-only generation planning")
	timestamp := fs.String("timestamp", "", "UTC result timestamp in YYYYMMDD-HHMMSS; defaults to current UTC time for non-dry-run init")
	shuffleSeed := fs.String("shuffle-seed", "", "deterministic shuffle seed for parent and scenario sampling")
	parentCount := fs.Int("parent-count", defaultParentCount, "number of parent sims to select for planning")
	scenarioCount := fs.Int("scenario-count", defaultScenarioCount, "number of root scenarios to sample uniformly")
	childCount := fs.Int("child-count", defaultChildCount, "number of planned child sims")
	maxPromotions := fs.Int("max-promotions", defaultMaxPromotions, "maximum accepted children for this planned generation")
	var includeSims stringListFlag
	var includeScenarios stringListFlag
	fs.Var(&includeSims, "include-sim", "simulation ID that must be included in selected parents; repeatable")
	fs.Var(&includeScenarios, "include-scenario", "scenario ID that must be included in the sampled scenarios; repeatable")
	if err := fs.Parse(args); err != nil {
		return err
	}
	repo, err := openRepo(*repoRoot)
	if err != nil {
		return err
	}
	population, err := discoverTrackedPopulation(repo)
	if err != nil {
		return err
	}
	if err := writeFormat(stdout, "population=%d\n", len(population)); err != nil {
		return err
	}
	for _, sim := range population {
		if err := writeFormat(stdout, "%s files=%d tree_hash=%s path=%s\n", sim.SimID, len(sim.Files), sim.TreeHash, sim.Path); err != nil {
			return err
		}
	}
	if *modelID == "" {
		if !*dryRun {
			return errUsage("init: -model is required")
		}
		return nil
	}
	scenarios, err := discoverScenarios(repo)
	if err != nil {
		return err
	}
	options := PlanOptions{
		RunGroupID:         *runGroupID,
		ModelID:            *modelID,
		ShuffleSeed:        *shuffleSeed,
		ParentCount:        *parentCount,
		ScenarioCount:      *scenarioCount,
		ChildCount:         *childCount,
		MaxPromotions:      *maxPromotions,
		IncludeSimIDs:      []string(includeSims),
		IncludeScenarioIDs: []string(includeScenarios),
	}
	plan, err := buildGenerationPlan(population, scenarios, options)
	if err != nil {
		return err
	}
	if !*dryRun {
		return writeInitialState(repo, stdout, plan, *timestamp)
	}
	return printGenerationPlan(stdout, plan)
}

func writeInitialState(repo Repo, stdout io.Writer, plan GenerationPlan, timestamp string) error {
	if plan.RunGroupID == "" {
		return errUsage("init: -run-group-id is required")
	}
	if timestamp == "" {
		timestamp = time.Now().UTC().Format("20060102-150405")
	}
	stateFile, err := statePath(repo, plan.RunGroupID)
	if err != nil {
		return err
	}
	if _, err := os.Stat(stateFile); err == nil {
		return fmt.Errorf("state file already exists: %s", repo.Rel(stateFile))
	} else if !os.IsNotExist(err) {
		return err
	}
	state, err := stateFromPlan(repo, plan, timestamp)
	if err != nil {
		return err
	}
	if err := writeGAStateAtomic(stateFile, state); err != nil {
		return err
	}
	if err := writeFormat(stdout, "state=%s population=%d parents=%d scenarios=%d children=%d cells=%d\n",
		repo.Rel(stateFile),
		len(state.Population),
		len(state.Parents),
		len(state.ScenarioSample),
		len(state.Children),
		len(state.Cells)); err != nil {
		return err
	}
	return printGenerationPlan(stdout, plan)
}

// stateFromPlan turns the conservative generation preview into the durable GA
// state file that score/generate checkpoint after every expensive operation.
//
// Intent: Non-dry-run init is the boundary where a cheap plan becomes the
// authoritative state for pending children and score cells. Source: DI-gijom
func stateFromPlan(repo Repo, plan GenerationPlan, timestamp string) (GAState, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	var population []GAStateSim
	for _, sim := range plan.Population {
		population = append(population, GAStateSim{
			SimID:    sim.SimID,
			Path:     sim.Path,
			TreeHash: sim.TreeHash,
		})
	}
	var scenarios []GAStateScenario
	for _, scenario := range plan.Scenarios {
		hash, err := sha256File(repo, scenario.Path)
		if err != nil {
			return GAState{}, err
		}
		scenarios = append(scenarios, GAStateScenario{
			ScenarioID:   scenario.ScenarioID,
			Path:         scenario.Path,
			SamplePolicy: "uniform root scenario sample",
			SHA256:       hash,
		})
	}
	var parents []GAStateParent
	for _, parent := range plan.Parents {
		parents = append(parents, GAStateParent{
			SimID:     parent.SimID,
			Rationale: "uniform tracked-population sample",
		})
	}
	children, err := plannedGAChildren(repo, plan)
	if err != nil {
		return GAState{}, err
	}
	cells := plannedGACells(plan.RunGroupID, plan.ModelID, timestamp, plan.Parents, children, plan.Scenarios)
	return GAState{
		Schema:         stateSchemaV1,
		RunGroupID:     plan.RunGroupID,
		CreatedAt:      now,
		UpdatedAt:      now,
		RepoCommit:     repo.GitCommit(),
		ModelID:        plan.ModelID,
		Population:     population,
		ScenarioSample: scenarios,
		Parents:        parents,
		Children:       children,
		Cells:          cells,
	}, nil
}

func plannedGAChildren(repo Repo, plan GenerationPlan) ([]GAChild, error) {
	var children []GAChild
	used := map[string]bool{}
	for index, child := range plan.Children {
		childID, err := mintChildSimID(repo, plan.RunGroupID, index, used)
		if err != nil {
			return nil, err
		}
		used[childID] = true
		children = append(children, GAChild{
			ChildID:            childID,
			Path:               proposalChildSimulationPath(plan.RunGroupID, childID),
			ParentIDs:          child.ParentIDs,
			Operation:          child.Operation,
			Status:             "queued",
			DesignDeltaSummary: "",
		})
	}
	return children, nil
}

func plannedGACells(runGroupID string, modelID string, timestamp string, parents []PopulationSim, children []GAChild, scenarios []Scenario) []GACell {
	var cells []GACell
	ordinal := 1
	for _, parent := range parents {
		for _, scenario := range scenarios {
			cells = append(cells, newGACell(runGroupID, ordinal, parent.SimID, scenario.ScenarioID, modelID, filepath.ToSlash(filepath.Join("results", parent.SimID, scenario.ScenarioID, modelID, timestamp+".json"))))
			ordinal++
		}
	}
	for _, child := range children {
		for _, scenario := range scenarios {
			cells = append(cells, newGACell(runGroupID, ordinal, child.ID(), scenario.ScenarioID, modelID, proposalChildResultPath(runGroupID, child.ID(), scenario.ScenarioID, modelID, timestamp)))
			ordinal++
		}
	}
	return cells
}

func newGACell(runGroupID string, ordinal int, simID string, scenarioID string, modelID string, resultPath string) GACell {
	return GACell{
		CellID:     fmt.Sprintf("%s-%06d-%s--%s--%s", runGroupID, ordinal, simID, scenarioID, modelID),
		SimID:      simID,
		ScenarioID: scenarioID,
		ModelID:    modelID,
		ResultPath: resultPath,
		Status:     "queued",
	}
}

func mintChildSimID(repo Repo, runGroupID string, index int, used map[string]bool) (string, error) {
	for attempt := 0; attempt < 1000; attempt++ {
		seed := fmt.Sprintf("%s-%d-%d-%d", runGroupID, index, attempt, time.Now().UnixNano())
		sum := sha256.Sum256([]byte(seed))
		handle := uint16ToProquint(uint16(sum[0])<<8 | uint16(sum[1]))
		childID := fmt.Sprintf("SIM-%s-child-pending-%04d", handle, index+1)
		if used[childID] {
			continue
		}
		proposalPath := strings.TrimSuffix(proposalChildSimulationPath(runGroupID, childID), "/")
		if _, err := os.Stat(repo.Path("simulations", childID)); err == nil {
			continue
		} else if err != nil && !os.IsNotExist(err) {
			return "", err
		}
		if _, err := os.Stat(repo.Abs(proposalPath)); os.IsNotExist(err) {
			return childID, nil
		} else if err != nil {
			return "", err
		}
	}
	return "", fmt.Errorf("could not mint unused child sim ID")
}

func printGenerationPlan(stdout io.Writer, plan GenerationPlan) error {
	if err := writeFormat(stdout, "plan run_group_id=%s model=%s parents=%d scenarios=%d children=%d max_promotions=%d parent_score_cells=%d child_score_cells=%d\n",
		plan.RunGroupID,
		plan.ModelID,
		len(plan.Parents),
		len(plan.Scenarios),
		len(plan.Children),
		plan.MaxPromotions,
		len(plan.ParentScoreCells),
		len(plan.ChildScoreCells)); err != nil {
		return err
	}
	for index, parent := range plan.Parents {
		if err := writeFormat(stdout, "parent %d %s tree_hash=%s\n", index+1, parent.SimID, parent.TreeHash); err != nil {
			return err
		}
	}
	for index, scenario := range plan.Scenarios {
		if err := writeFormat(stdout, "scenario %d %s path=%s\n", index+1, scenario.ScenarioID, scenario.Path); err != nil {
			return err
		}
	}
	for index, child := range plan.Children {
		if err := writeFormat(stdout, "child %d %s operation=%s parents=%s path=%s\n", index+1, child.ChildID, child.Operation, strings.Join(child.ParentIDs, ","), child.ResultPath); err != nil {
			return err
		}
	}
	return nil
}

func discoverTrackedPopulation(repo Repo) ([]PopulationSim, error) {
	output, err := repo.Git("ls-files", "-z", "--", "simulations")
	if err != nil {
		return nil, err
	}
	bySim := map[string][]string{}
	for _, rel := range splitNUL(output) {
		simID, ok := simIDFromTrackedPath(rel)
		if !ok {
			continue
		}
		abs := repo.Abs(rel)
		info, err := os.Stat(abs)
		if err != nil || info.IsDir() {
			continue
		}
		bySim[simID] = append(bySim[simID], filepath.ToSlash(rel))
	}
	var simIDs []string
	for simID := range bySim {
		simIDs = append(simIDs, simID)
	}
	sort.Strings(simIDs)
	var population []PopulationSim
	for _, simID := range simIDs {
		files := bySim[simID]
		sort.Strings(files)
		treeHash, err := trackedTreeHash(repo, files)
		if err != nil {
			return nil, err
		}
		population = append(population, PopulationSim{
			SimID:    simID,
			Path:     filepath.ToSlash(filepath.Join("simulations", simID)) + "/",
			Files:    files,
			TreeHash: treeHash,
		})
	}
	return population, nil
}

func splitNUL(output []byte) []string {
	raw := bytes.Split(output, []byte{0})
	var values []string
	for _, item := range raw {
		value := strings.TrimSpace(string(item))
		if value != "" {
			values = append(values, filepath.ToSlash(value))
		}
	}
	return values
}

func simIDFromTrackedPath(rel string) (string, bool) {
	parts := strings.Split(filepath.ToSlash(rel), "/")
	if len(parts) < 3 || parts[0] != "simulations" || !strings.HasPrefix(parts[1], "SIM-") {
		return "", false
	}
	return parts[1], true
}

func trackedTreeHash(repo Repo, files []string) (string, error) {
	digest := sha256.New()
	for _, rel := range files {
		if err := hashWrite(digest, []byte(rel)); err != nil {
			return "", err
		}
		if err := hashWrite(digest, []byte{0}); err != nil {
			return "", err
		}
		content, err := os.ReadFile(repo.Abs(rel))
		if err != nil {
			return "", err
		}
		fileHash := sha256.Sum256(content)
		if err := hashWrite(digest, []byte(hex.EncodeToString(fileHash[:]))); err != nil {
			return "", err
		}
		if err := hashWrite(digest, []byte{0}); err != nil {
			return "", err
		}
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}

func hashWrite(writer hash.Hash, data []byte) error {
	written, err := writer.Write(data)
	if err != nil {
		return err
	}
	if written != len(data) {
		return fmt.Errorf("hash write short count: wrote %d of %d", written, len(data))
	}
	return nil
}

// Intent: Turn the cheap vocabulary/source audit into a runnable targeted
// rubric-v2 state so operators can re-score the most affected sims before
// paying for any broader corpus backfill. Source: DI-roruj
func runBackfillInit(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("backfill-init", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	runGroupID := fs.String("run-group-id", "", "run group ID for the targeted v2 backfill state")
	timestamp := fs.String("timestamp", "", "UTC result timestamp in YYYYMMDD-HHMMSS; defaults to current UTC time")
	model := fs.String("model", "", "optional canonical source result model filter")
	cleanEnvelopeCount := fs.Int("clean-envelope-count", defaultBackfillCleanEnvelopeCount, "number of clean grid-envelope sims to keep as calibration targets")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *runGroupID == "" {
		return errUsage("backfill-init: -run-group-id is required")
	}
	repo, err := openRepo(*repoRoot)
	if err != nil {
		return err
	}
	if *timestamp == "" {
		*timestamp = time.Now().UTC().Format("20060102-150405")
	}
	stateFile, err := statePath(repo, *runGroupID)
	if err != nil {
		return err
	}
	if _, err := os.Stat(stateFile); err == nil {
		return fmt.Errorf("state file already exists: %s", repo.Rel(stateFile))
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	records, err := auditCanonicalV1Results(repo, *model, "")
	if err != nil {
		return err
	}
	selection := selectTargetedBackfill(records, *cleanEnvelopeCount)
	if len(selection.Records) == 0 {
		return fmt.Errorf("backfill-init: no exact-match hard-hit or clean-envelope calibration results were selected")
	}
	population, err := discoverTrackedPopulation(repo)
	if err != nil {
		return err
	}
	state, err := stateFromBackfillSelection(repo, *runGroupID, *timestamp, population, selection)
	if err != nil {
		return err
	}
	if err := writeGAStateAtomic(stateFile, state); err != nil {
		return err
	}
	if err := writeFormat(stdout,
		"state=%s parents=%d scenarios=%d cells=%d hard_hit_sims=%d clean_envelope_sims=%d\n",
		repo.Rel(stateFile),
		len(state.Parents),
		len(state.ScenarioSample),
		len(state.Cells),
		len(selection.HardHitSimIDs),
		len(selection.CleanEnvelopeSimIDs),
	); err != nil {
		return err
	}
	return nil
}

// Intent: Materialize a fresh GA state from audited historical results without
// mutating any v1 scored artifact, so targeted rubric-v2 rescoring starts from
// additive evidence only. Source: DI-roruj
func stateFromBackfillSelection(repo Repo, runGroupID string, timestamp string, population []PopulationSim, selection backfillSelection) (GAState, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	var statePopulation []GAStateSim
	for _, sim := range population {
		statePopulation = append(statePopulation, GAStateSim{
			SimID:    sim.SimID,
			Path:     sim.Path,
			TreeHash: sim.TreeHash,
		})
	}
	hardHit := map[string]bool{}
	for _, simID := range selection.HardHitSimIDs {
		hardHit[simID] = true
	}
	cleanEnvelope := map[string]bool{}
	for _, simID := range selection.CleanEnvelopeSimIDs {
		cleanEnvelope[simID] = true
	}
	parentCounts := map[string]int{}
	scenarioByID := map[string]GAStateScenario{}
	var cells []GACell
	modelIDs := map[string]bool{}
	for index, record := range selection.Records {
		parentCounts[record.Result.SimID]++
		modelIDs[record.Result.ModelID] = true
		if _, ok := scenarioByID[record.Result.ScenarioID]; !ok {
			hash, err := sha256File(repo, record.Result.Source.ScenarioPath)
			if err != nil {
				return GAState{}, err
			}
			scenarioByID[record.Result.ScenarioID] = GAStateScenario{
				ScenarioID:   record.Result.ScenarioID,
				Path:         record.Result.Source.ScenarioPath,
				SamplePolicy: "audit-first targeted rubric-v2 backfill from historical v1 results",
				SHA256:       hash,
			}
		}
		resultPath := filepath.ToSlash(filepath.Join("results", record.Result.SimID, record.Result.ScenarioID, record.Result.ModelID, timestamp+".json"))
		cell := newGACell(runGroupID, index+1, record.Result.SimID, record.Result.ScenarioID, record.Result.ModelID, resultPath)
		cell.Provider = record.Result.Runner.Provider
		if cell.Provider == "" {
			cell.Provider = "openai"
		}
		cell.APIModel = record.Result.Runner.APIModel
		if cell.APIModel == "" {
			cell.APIModel = deriveAPIModelFromModelID(record.Result.ModelID)
		}
		cell.ReasoningEffort = record.Result.Runner.ReasoningEffort
		if cell.ReasoningEffort == "" {
			cell.ReasoningEffort = defaultScoreReasoningEffort
		}
		cell.ValidationMessage = "queued from audit-first rubric-v2 backfill selection"
		cells = append(cells, cell)
	}
	var scenarioIDs []string
	for scenarioID := range scenarioByID {
		scenarioIDs = append(scenarioIDs, scenarioID)
	}
	sort.Strings(scenarioIDs)
	var scenarios []GAStateScenario
	for _, scenarioID := range scenarioIDs {
		scenarios = append(scenarios, scenarioByID[scenarioID])
	}
	var parentIDs []string
	for simID := range parentCounts {
		parentIDs = append(parentIDs, simID)
	}
	sort.Strings(parentIDs)
	var parents []GAStateParent
	for _, simID := range parentIDs {
		reason := "audit-first rubric-v2 backfill exact-match result set"
		switch {
		case hardHit[simID]:
			reason = fmt.Sprintf("audit-first rubric-v2 backfill hard_hit exact-match results=%d", parentCounts[simID])
		case cleanEnvelope[simID]:
			reason = fmt.Sprintf("audit-first rubric-v2 backfill clean grid-envelope calibration results=%d", parentCounts[simID])
		}
		parents = append(parents, GAStateParent{
			SimID:     simID,
			Rationale: reason,
		})
	}
	modelID := "mixed"
	if len(modelIDs) == 1 {
		for only := range modelIDs {
			modelID = only
		}
	}
	return GAState{
		Schema:         stateSchemaV1,
		RunGroupID:     runGroupID,
		CreatedAt:      now,
		UpdatedAt:      now,
		RepoCommit:     repo.GitCommit(),
		ModelID:        modelID,
		Population:     statePopulation,
		ScenarioSample: scenarios,
		Parents:        parents,
		Children:       nil,
		Cells:          cells,
	}, nil
}
