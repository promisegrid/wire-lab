package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

const (
	defaultParentCount   = 3
	defaultScenarioCount = 5
	defaultChildCount    = 4
	defaultMaxPromotions = 2
)

// GenerationPlan is a read-only preview of a conservative GA generation.
//
// Intent: Size GA work before spending model tokens or writing state so early
// generations stay small, deterministic, and auditable. Source: DI-zusit
type GenerationPlan struct {
	RunGroupID       string
	ModelID          string
	Parents          []PopulationSim
	Scenarios        []Scenario
	Children         []PlannedChild
	ParentScoreCells []PlannedCell
	ChildScoreCells  []PlannedCell
	MaxPromotions    int
}

type Scenario struct {
	ScenarioID string
	Path       string
}

type PlannedChild struct {
	ChildID    string
	Operation  string
	ParentIDs  []string
	ResultPath string
}

type PlannedCell struct {
	SimID      string
	ScenarioID string
	ModelID    string
	ResultPath string
}

type PlanOptions struct {
	RunGroupID    string
	ModelID       string
	ShuffleSeed   string
	ParentCount   int
	ScenarioCount int
	ChildCount    int
	MaxPromotions int
}

func defaultPlanOptions() PlanOptions {
	return PlanOptions{
		ParentCount:   defaultParentCount,
		ScenarioCount: defaultScenarioCount,
		ChildCount:    defaultChildCount,
		MaxPromotions: defaultMaxPromotions,
	}
}

func discoverScenarios(repo Repo) ([]Scenario, error) {
	entries, err := os.ReadDir(repo.Path("scenarios"))
	if err != nil {
		return nil, err
	}
	var scenarios []Scenario
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		id := entry.Name()
		path := filepath.ToSlash(filepath.Join("scenarios", id, id+".md"))
		if info, err := os.Stat(repo.Abs(path)); err == nil && !info.IsDir() {
			scenarios = append(scenarios, Scenario{ScenarioID: id, Path: path})
		}
	}
	sort.Slice(scenarios, func(i, j int) bool {
		return scenarios[i].ScenarioID < scenarios[j].ScenarioID
	})
	return scenarios, nil
}

func buildGenerationPlan(population []PopulationSim, scenarios []Scenario, options PlanOptions) (GenerationPlan, error) {
	if err := validatePlanOptions(options); err != nil {
		return GenerationPlan{}, err
	}
	if options.ModelID == "" {
		return GenerationPlan{}, fmt.Errorf("model is required for generation planning")
	}
	if len(population) == 0 {
		return GenerationPlan{}, fmt.Errorf("no tracked simulation population available")
	}
	if len(scenarios) == 0 {
		return GenerationPlan{}, fmt.Errorf("no scenarios available")
	}
	parents, err := selectParents(population, options.ParentCount, options.ShuffleSeed)
	if err != nil {
		return GenerationPlan{}, err
	}
	sample, err := selectScenarios(scenarios, options.ScenarioCount, options.ShuffleSeed)
	if err != nil {
		return GenerationPlan{}, err
	}
	children := planChildren(parents, options.ChildCount)
	parentCells := planCellsFromParents(parents, sample, options.ModelID)
	childCells := planCellsFromChildren(children, sample, options.ModelID)
	return GenerationPlan{
		RunGroupID:       options.RunGroupID,
		ModelID:          options.ModelID,
		Parents:          parents,
		Scenarios:        sample,
		Children:         children,
		ParentScoreCells: parentCells,
		ChildScoreCells:  childCells,
		MaxPromotions:    options.MaxPromotions,
	}, nil
}

func validatePlanOptions(options PlanOptions) error {
	if options.ParentCount <= 0 {
		return fmt.Errorf("parent-count must be positive")
	}
	if options.ScenarioCount <= 0 {
		return fmt.Errorf("scenario-count must be positive")
	}
	if options.ChildCount <= 0 {
		return fmt.Errorf("child-count must be positive")
	}
	if options.MaxPromotions < 0 {
		return fmt.Errorf("max-promotions must be non-negative")
	}
	if options.MaxPromotions > options.ChildCount {
		return fmt.Errorf("max-promotions cannot exceed child-count")
	}
	return nil
}

func selectParents(population []PopulationSim, count int, seedText string) ([]PopulationSim, error) {
	shuffled := append([]PopulationSim(nil), population...)
	if err := shuffleBySeed(seedText, len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}); err != nil {
		return nil, err
	}
	if count > len(shuffled) {
		count = len(shuffled)
	}
	return shuffled[:count], nil
}

func selectScenarios(scenarios []Scenario, count int, seedText string) ([]Scenario, error) {
	shuffled := append([]Scenario(nil), scenarios...)
	if err := shuffleBySeed(seedText, len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}); err != nil {
		return nil, err
	}
	if count > len(shuffled) {
		count = len(shuffled)
	}
	return shuffled[:count], nil
}

func shuffleBySeed(seedText string, size int, swap func(i, j int)) error {
	if seedText == "" {
		return nil
	}
	seed, err := strconv.ParseInt(seedText, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid shuffle seed: %w", err)
	}
	rand.New(rand.NewSource(seed)).Shuffle(size, swap)
	return nil
}

func planChildren(parents []PopulationSim, childCount int) []PlannedChild {
	var children []PlannedChild
	for index := 0; index < childCount; index++ {
		child := PlannedChild{
			ChildID: fmt.Sprintf("planned-child-%04d", index+1),
		}
		if index%2 == 1 && len(parents) >= 2 {
			child.Operation = "crossover"
			child.ParentIDs = []string{
				parents[index%len(parents)].SimID,
				parents[(index+1)%len(parents)].SimID,
			}
		} else {
			child.Operation = "mutation"
			child.ParentIDs = []string{parents[index%len(parents)].SimID}
		}
		child.ResultPath = fmt.Sprintf("simulations/SIM-<handle>-%s/", child.ChildID)
		children = append(children, child)
	}
	return children
}

func planCellsFromParents(parents []PopulationSim, scenarios []Scenario, modelID string) []PlannedCell {
	var cells []PlannedCell
	for _, parent := range parents {
		for _, scenario := range scenarios {
			cells = append(cells, PlannedCell{
				SimID:      parent.SimID,
				ScenarioID: scenario.ScenarioID,
				ModelID:    modelID,
				ResultPath: plannedResultPath(parent.SimID, scenario.ScenarioID, modelID),
			})
		}
	}
	return cells
}

func planCellsFromChildren(children []PlannedChild, scenarios []Scenario, modelID string) []PlannedCell {
	var cells []PlannedCell
	for _, child := range children {
		for _, scenario := range scenarios {
			cells = append(cells, PlannedCell{
				SimID:      child.ChildID,
				ScenarioID: scenario.ScenarioID,
				ModelID:    modelID,
				ResultPath: plannedResultPath(child.ChildID, scenario.ScenarioID, modelID),
			})
		}
	}
	return cells
}

func plannedResultPath(simID string, scenarioID string, modelID string) string {
	return fmt.Sprintf("results/%s/%s/%s/<YYYYMMDD-HHMMSS>.json", simID, scenarioID, modelID)
}
