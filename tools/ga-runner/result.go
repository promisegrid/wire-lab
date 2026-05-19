package main

import (
	"encoding/json"
	"os"
)

const resultSchemaV1 = "promisegrid.ga.result.v1"

// FitnessResult is the machine-readable evidence shape for one GA runner cell.
//
// Intent: Make fitness the result artifact itself, not a parallel score tree or
// a post-hoc parser over Markdown canaries. Source: DI-zanon; DI-pobus
type FitnessResult struct {
	Schema       string         `json:"schema"`
	ResultID     string         `json:"result_id"`
	RunGroupID   string         `json:"run_group_id"`
	CellID       string         `json:"cell_id"`
	SimID        string         `json:"sim_id"`
	ScenarioID   string         `json:"scenario_id"`
	ModelID      string         `json:"model_id"`
	TimestampUTC string         `json:"timestamp_utc"`
	ResultPath   string         `json:"result_path"`
	Runner       RunnerInfo     `json:"runner"`
	Source       SourceInfo     `json:"source"`
	Rubric       RubricInfo     `json:"rubric"`
	Scores       FitnessScores  `json:"scores"`
	Fitness      FitnessSummary `json:"fitness"`
	Assessment   Assessment     `json:"assessment"`
}

type RunnerInfo struct {
	Tool              string  `json:"tool"`
	Provider          string  `json:"provider,omitempty"`
	APIModel          string  `json:"api_model,omitempty"`
	ReasoningEffort   string  `json:"reasoning_effort,omitempty"`
	ServiceTier       string  `json:"service_tier,omitempty"`
	ServedServiceTier string  `json:"served_service_tier,omitempty"`
	RequestID         string  `json:"request_id,omitempty"`
	ResponseID        string  `json:"response_id,omitempty"`
	InputTokens       int     `json:"input_tokens,omitempty"`
	CachedInputTokens int     `json:"cached_input_tokens,omitempty"`
	OutputTokens      int     `json:"output_tokens,omitempty"`
	CostUSD           float64 `json:"cost_usd,omitempty"`
}

type SourceInfo struct {
	RepoCommit         string       `json:"repo_commit"`
	SimPath            string       `json:"sim_path"`
	ScenarioPath       string       `json:"scenario_path"`
	RootContractPaths  []string     `json:"root_contract_paths"`
	Files              []SourceFile `json:"files"`
	SimulationTreeHash string       `json:"simulation_tree_hash"`
}

type SourceFile struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type RubricInfo struct {
	RubricVersion string            `json:"rubric_version"`
	ScoreScale    string            `json:"score_scale"`
	ScoreMeanings map[string]string `json:"score_meanings"`
	Axes          []string          `json:"axes"`
}

type FitnessScores struct {
	ScenarioFit                int `json:"scenario_fit"`
	PromiseGridAlignment       int `json:"promisegrid_alignment"`
	Auditability               int `json:"auditability"`
	EvolutionSafety            int `json:"evolution_safety"`
	LayerBoundaryClarity       int `json:"layer_boundary_clarity"`
	FailureHandling            int `json:"failure_handling"`
	ImplementationPlausibility int `json:"implementation_plausibility"`
	RiskPenalty                int `json:"risk_penalty"`
}

type FitnessSummary struct {
	Raw              float64 `json:"raw"`
	Normalized0To100 float64 `json:"normalized_0_100"`
	Confidence0To1   float64 `json:"confidence_0_1"`
}

type Assessment struct {
	Rationale         string   `json:"rationale"`
	Strengths         []string `json:"strengths"`
	Weaknesses        []string `json:"weaknesses"`
	Risks             []string `json:"risks"`
	OpenQuestions     []string `json:"open_questions"`
	AuthorityBoundary string   `json:"authority_boundary"`
}

// writeFitnessResultAtomic writes a JSON result through a same-directory temp
// file so readers never see a partially-written fitness artifact.
//
// Intent: Future scoring runs can checkpoint result files safely without
// exposing half-written JSON to validation or culling. Source: DI-pobus
func writeFitnessResultAtomic(path string, result FitnessResult) error {
	bytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	if err := ensureParent(path); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(bytes, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
