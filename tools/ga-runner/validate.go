package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var resultTimestampPattern = regexp.MustCompile(`^\d{8}-\d{6}$`)

type resultPathParts struct {
	SimID      string
	ScenarioID string
	ModelID    string
	Timestamp  string
}

func runValidate(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	resultPath := fs.String("result", "", "optional explicit JSON result path")
	model := fs.String("model", "", "optional model filter")
	timestamp := fs.String("timestamp", "", "optional timestamp filter")
	maxErrors := fs.Int("max-errors", 50, "maximum printed errors")
	if err := fs.Parse(args); err != nil {
		return err
	}
	repo, err := openRepo(*repoRoot)
	if err != nil {
		return err
	}
	targets, err := validateTargets(repo, *resultPath, *model, *timestamp)
	if err != nil {
		return err
	}
	sort.Strings(targets)
	failed := 0
	printed := 0
	for _, path := range targets {
		issues := validateResultFile(repo, path)
		if len(issues) == 0 {
			continue
		}
		failed++
		if printed < *maxErrors {
			if err := writeFormat(stdout, "%s:\n", repo.Rel(path)); err != nil {
				return err
			}
			for _, issue := range issues {
				if err := writeFormat(stdout, "  - %s\n", issue); err != nil {
					return err
				}
			}
			printed++
		}
	}
	if len(targets) == 0 {
		return fmt.Errorf("no JSON fitness result files matched selection")
	}
	if err := writeFormat(stdout, "validated=%d failed=%d\n", len(targets), failed); err != nil {
		return err
	}
	if failed > 0 {
		return fmt.Errorf("validation failed")
	}
	return nil
}

func validateTargets(repo Repo, explicitResult string, model string, timestamp string) ([]string, error) {
	if explicitResult != "" {
		return []string{repo.Abs(explicitResult)}, nil
	}
	return findResultFiles(repo, model, timestamp)
}

// findResultFiles discovers only GA JSON result artifacts, including ignored
// proposal-stage child score evidence. Markdown canary files are intentionally
// invisible to GA-runner validation and future scoring logic.
//
// Intent: Keep old matrix-runner canaries as historical evidence without letting
// them contaminate GA fitness selection. Source: DI-ramar; DI-pobus; DI-lirat
func findResultFiles(repo Repo, model string, timestamp string) ([]string, error) {
	var paths []string
	for _, rootName := range []string{"results", "proposals"} {
		root := repo.Path(rootName)
		if info, err := os.Stat(root); err != nil || !info.IsDir() {
			if err != nil && !os.IsNotExist(err) {
				return nil, err
			}
			continue
		}
		err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() || filepath.Ext(path) != ".json" {
				return nil
			}
			parts, issues := parseResultPath(repo, path)
			if len(issues) > 0 {
				return nil
			}
			if model != "" && parts.ModelID != model {
				return nil
			}
			if timestamp != "" && parts.Timestamp != timestamp {
				return nil
			}
			paths = append(paths, path)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return paths, nil
}

func validateResultFile(repo Repo, path string) []string {
	var issues []string
	parts, pathIssues := parseResultPath(repo, path)
	issues = append(issues, pathIssues...)
	bytes, err := os.ReadFile(path)
	if err != nil {
		return append(issues, err.Error())
	}
	var result FitnessResult
	if err := json.Unmarshal(bytes, &result); err != nil {
		return append(issues, "invalid JSON: "+err.Error())
	}
	issues = append(issues, validateResultIdentity(repo, path, parts, result)...)
	issues = append(issues, validateRunner(result.Runner)...)
	issues = append(issues, validateSource(result.Source)...)
	issues = append(issues, validateRubric(result.Rubric)...)
	issues = append(issues, validateScores(result.Scores)...)
	issues = append(issues, validateFitness(result.Fitness)...)
	issues = append(issues, validateAssessment(result.Assessment)...)
	return issues
}

func parseResultPath(repo Repo, path string) (resultPathParts, []string) {
	rel := repo.Rel(path)
	relParts := strings.Split(rel, "/")
	var issues []string
	switch {
	case len(relParts) == 5 && relParts[0] == "results":
		return parseResultPathParts(relParts[1], relParts[2], relParts[3], relParts[4], issues)
	case len(relParts) == 7 && relParts[0] == "proposals" && relParts[2] == "results":
		return parseResultPathParts(relParts[3], relParts[4], relParts[5], relParts[6], issues)
	default:
		return resultPathParts{}, []string{"path shape must be results/<sim>/<scenario>/<model>/<timestamp>.json or proposals/<run-group>/results/<sim>/<scenario>/<model>/<timestamp>.json"}
	}
}

func parseResultPathParts(simID string, scenarioID string, modelID string, timestampFile string, issues []string) (resultPathParts, []string) {
	if filepath.Ext(timestampFile) != ".json" {
		issues = append(issues, "path extension must be .json")
	}
	if !strings.HasPrefix(simID, "SIM-") {
		issues = append(issues, "sim path component must start with SIM-")
	}
	timestamp := strings.TrimSuffix(timestampFile, ".json")
	if !resultTimestampPattern.MatchString(timestamp) {
		issues = append(issues, "timestamp must match YYYYMMDD-HHMMSS")
	}
	return resultPathParts{
		SimID:      simID,
		ScenarioID: scenarioID,
		ModelID:    modelID,
		Timestamp:  timestamp,
	}, issues
}

func validateResultIdentity(repo Repo, path string, parts resultPathParts, result FitnessResult) []string {
	var issues []string
	if result.Schema != resultSchemaV1 {
		issues = append(issues, "schema must be "+resultSchemaV1)
	}
	expectedID := strings.Join([]string{parts.SimID, parts.ScenarioID, parts.ModelID, parts.Timestamp}, "-")
	if result.ResultID != expectedID {
		issues = append(issues, "result_id does not match path-derived result ID")
	}
	if result.RunGroupID == "" {
		issues = append(issues, "run_group_id is required")
	}
	if result.CellID == "" {
		issues = append(issues, "cell_id is required")
	}
	if result.SimID != parts.SimID {
		issues = append(issues, "sim_id does not match path")
	}
	if result.ScenarioID != parts.ScenarioID {
		issues = append(issues, "scenario_id does not match path")
	}
	if result.ModelID != parts.ModelID {
		issues = append(issues, "model_id does not match path")
	}
	if result.TimestampUTC != parts.Timestamp {
		issues = append(issues, "timestamp_utc does not match path timestamp")
	}
	if result.ResultPath != repo.Rel(path) {
		issues = append(issues, "result_path does not match path")
	}
	return issues
}

func validateRunner(runner RunnerInfo) []string {
	if runner.Tool == "" {
		return []string{"runner.tool is required"}
	}
	return nil
}

func validateSource(source SourceInfo) []string {
	var issues []string
	if source.RepoCommit == "" {
		issues = append(issues, "source.repo_commit is required")
	}
	if source.SimPath == "" {
		issues = append(issues, "source.sim_path is required")
	}
	if source.ScenarioPath == "" {
		issues = append(issues, "source.scenario_path is required")
	}
	if len(source.RootContractPaths) == 0 {
		issues = append(issues, "source.root_contract_paths must not be empty")
	}
	if len(source.Files) == 0 {
		issues = append(issues, "source.files must not be empty")
	}
	for index, file := range source.Files {
		if file.Path == "" {
			issues = append(issues, fmt.Sprintf("source.files[%d].path is required", index))
		}
		if file.SHA256 == "" {
			issues = append(issues, fmt.Sprintf("source.files[%d].sha256 is required", index))
		}
	}
	if source.SimulationTreeHash == "" {
		issues = append(issues, "source.simulation_tree_hash is required")
	}
	return issues
}

func validateRubric(rubric RubricInfo) []string {
	var issues []string
	if rubric.RubricVersion == "" {
		issues = append(issues, "rubric.rubric_version is required")
	}
	if rubric.ScoreScale == "" {
		issues = append(issues, "rubric.score_scale is required")
	}
	if len(rubric.ScoreMeanings) == 0 {
		issues = append(issues, "rubric.score_meanings must not be empty")
	}
	if len(rubric.Axes) == 0 {
		issues = append(issues, "rubric.axes must not be empty")
	}
	return issues
}

func validateScores(scores FitnessScores) []string {
	scoreMap := map[string]int{
		"scenario_fit":                scores.ScenarioFit,
		"promisegrid_alignment":       scores.PromiseGridAlignment,
		"auditability":                scores.Auditability,
		"evolution_safety":            scores.EvolutionSafety,
		"layer_boundary_clarity":      scores.LayerBoundaryClarity,
		"failure_handling":            scores.FailureHandling,
		"implementation_plausibility": scores.ImplementationPlausibility,
		"risk_penalty":                scores.RiskPenalty,
	}
	var issues []string
	for name, value := range scoreMap {
		if value < 0 || value > 5 {
			issues = append(issues, fmt.Sprintf("scores.%s must be between 0 and 5", name))
		}
	}
	sort.Strings(issues)
	return issues
}

func validateFitness(fitness FitnessSummary) []string {
	var issues []string
	if fitness.Normalized0To100 < 0 || fitness.Normalized0To100 > 100 {
		issues = append(issues, "fitness.normalized_0_100 must be between 0 and 100")
	}
	if fitness.Confidence0To1 < 0 || fitness.Confidence0To1 > 1 {
		issues = append(issues, "fitness.confidence_0_1 must be between 0 and 1")
	}
	return issues
}

func validateAssessment(assessment Assessment) []string {
	var issues []string
	if assessment.Rationale == "" {
		issues = append(issues, "assessment.rationale is required")
	}
	if assessment.AuthorityBoundary == "" {
		issues = append(issues, "assessment.authority_boundary is required")
	}
	return issues
}
