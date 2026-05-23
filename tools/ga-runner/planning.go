package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	defaultParentCount   = 3
	defaultScenarioCount = 5
	defaultChildCount    = 4
	defaultMaxPromotions = 2
	childOperationBreed  = "breed"
)

// GenerationPlan is a read-only preview of a conservative GA generation.
//
// Intent: Size GA work before spending model tokens or writing state so early
// generations stay small, deterministic, and auditable. Source: DI-zusit
type GenerationPlan struct {
	RunGroupID       string
	ModelID          string
	Population       []PopulationSim
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
	RunGroupID         string
	ModelID            string
	ShuffleSeed        string
	ParentCount        int
	ScenarioCount      int
	ChildCount         int
	MaxPromotions      int
	IncludeSimIDs      []string
	IncludeScenarioIDs []string
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
	parents, err := selectParents(population, options.ParentCount, options.ShuffleSeed, options.IncludeSimIDs)
	if err != nil {
		return GenerationPlan{}, err
	}
	if len(parents) < 2 {
		return GenerationPlan{}, fmt.Errorf("at least two parent sims are required for breed child planning")
	}
	sample, err := selectScenarios(scenarios, options.ScenarioCount, options.ShuffleSeed, options.IncludeScenarioIDs)
	if err != nil {
		return GenerationPlan{}, err
	}
	children := planChildren(parents, options.ChildCount)
	parentCells := planCellsFromParents(parents, sample, options.ModelID)
	childCells := planCellsFromChildren(options.RunGroupID, children, sample, options.ModelID)
	return GenerationPlan{
		RunGroupID:       options.RunGroupID,
		ModelID:          options.ModelID,
		Population:       append([]PopulationSim(nil), population...),
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
	if options.ParentCount < 2 {
		return fmt.Errorf("parent-count must be at least 2 for breed child planning")
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

func selectParents(population []PopulationSim, count int, seedText string, includeIDs []string) ([]PopulationSim, error) {
	// Intent: Focused canaries must include explicitly named new or suspect sims,
	// while keeping machine-tagged negative controls out of the default GA parent
	// pool unless the operator names them explicitly. Source: DI-duzur; DI-kuzag
	byID := map[string]PopulationSim{}
	for _, sim := range population {
		byID[sim.SimID] = sim
	}
	var selected []PopulationSim
	selectedIDs := map[string]bool{}
	for _, id := range uniqueNonEmptyStrings(includeIDs) {
		item, ok := byID[id]
		if !ok {
			return nil, fmt.Errorf("included simulation %q was not discovered", id)
		}
		item.ExplicitInclude = true
		selected = append(selected, item)
		selectedIDs[id] = true
	}
	var eligible []PopulationSim
	for _, sim := range population {
		if selectedIDs[sim.SimID] {
			continue
		}
		if sim.Role == simRoleNegativeCtl {
			continue
		}
		eligible = append(eligible, sim)
	}
	shuffled := append([]PopulationSim(nil), eligible...)
	if err := shuffleBySeed(seedText, len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}); err != nil {
		return nil, err
	}
	target := count
	if len(selected) > target {
		target = len(selected)
	}
	if maxSelectable := len(selected) + len(shuffled); target > maxSelectable {
		target = maxSelectable
	}
	for _, sim := range shuffled {
		if len(selected) >= target {
			break
		}
		selected = append(selected, sim)
	}
	return selected, nil
}

func selectScenarios(scenarios []Scenario, count int, seedText string, includeIDs []string) ([]Scenario, error) {
	// Intent: Focused canaries must include explicitly named scenarios while
	// keeping random scenario pressure for apples-to-apples breadth. Source:
	// DI-duzur
	return selectByID(scenarios, count, seedText, includeIDs, func(scenario Scenario) string {
		return scenario.ScenarioID
	}, "scenario")
}

func selectByID[T any](items []T, count int, seedText string, includeIDs []string, itemID func(T) string, label string) ([]T, error) {
	byID := map[string]T{}
	for _, item := range items {
		byID[itemID(item)] = item
	}
	var selected []T
	selectedIDs := map[string]bool{}
	for _, id := range uniqueNonEmptyStrings(includeIDs) {
		item, ok := byID[id]
		if !ok {
			return nil, fmt.Errorf("included %s %q was not discovered", label, id)
		}
		selected = append(selected, item)
		selectedIDs[id] = true
	}
	shuffled := append([]T(nil), items...)
	if err := shuffleBySeed(seedText, len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}); err != nil {
		return nil, err
	}
	target := count
	if len(selected) > target {
		target = len(selected)
	}
	if target > len(items) {
		target = len(items)
	}
	for _, item := range shuffled {
		if len(selected) >= target {
			break
		}
		id := itemID(item)
		if selectedIDs[id] {
			continue
		}
		selected = append(selected, item)
		selectedIDs[id] = true
	}
	return selected, nil
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

func uniqueNonEmptyStrings(values []string) []string {
	seen := map[string]bool{}
	var unique []string
	for _, value := range values {
		clean := strings.TrimSpace(value)
		if clean == "" || seen[clean] {
			continue
		}
		seen[clean] = true
		unique = append(unique, clean)
	}
	return unique
}

func planChildren(parents []PopulationSim, childCount int) []PlannedChild {
	// Intent: LLM-based child generation uses one two-parent breed operator
	// instead of pretending that prompt synthesis is byte-level mutation or
	// crossover. Source: DI-sohus
	var children []PlannedChild
	for index := 0; index < childCount; index++ {
		child := PlannedChild{
			ChildID:   fmt.Sprintf("planned-child-%04d", index+1),
			Operation: childOperationBreed,
		}
		child.ParentIDs = plannedBreedParentIDs(parents, index)
		child.ResultPath = "proposals/<run-group-id>/simulations/SIM-<handle>-child-<descriptive-slug>/"
		children = append(children, child)
	}
	return children
}

func plannedBreedParentIDs(parents []PopulationSim, index int) []string {
	if len(parents) < 2 {
		return nil
	}
	first := parents[index%len(parents)].SimID
	for offset := 1; offset < len(parents); offset++ {
		second := parents[(index+offset)%len(parents)].SimID
		if second != first {
			return []string{first, second}
		}
	}
	return nil
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

func planCellsFromChildren(runGroupID string, children []PlannedChild, scenarios []Scenario, modelID string) []PlannedCell {
	var cells []PlannedCell
	for _, child := range children {
		for _, scenario := range scenarios {
			cells = append(cells, PlannedCell{
				SimID:      child.ChildID,
				ScenarioID: scenario.ScenarioID,
				ModelID:    modelID,
				ResultPath: plannedChildResultPath(runGroupID, child.ChildID, scenario.ScenarioID, modelID),
			})
		}
	}
	return cells
}

func plannedResultPath(simID string, scenarioID string, modelID string) string {
	return fmt.Sprintf("results/%s/%s/%s/<YYYYMMDD-HHMMSS>.json", simID, scenarioID, modelID)
}

func plannedChildResultPath(runGroupID string, simID string, scenarioID string, modelID string) string {
	if runGroupID == "" {
		runGroupID = "<run-group-id>"
	}
	return fmt.Sprintf("proposals/%s/results/%s/%s/%s/<YYYYMMDD-HHMMSS>.json", runGroupID, simID, scenarioID, modelID)
}
