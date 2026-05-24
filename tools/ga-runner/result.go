package main

import (
	"encoding/json"
	"os"
)

const (
	resultSchemaV1  = "promisegrid.ga.result.v1"
	resultSchemaV2  = "promisegrid.ga.result.v2"
	resultSchemaV3  = "promisegrid.ga.result.v3"
	rubricVersionV1 = "ga-rubric-20260519-v1"
	rubricVersionV2 = "ga-rubric-20260522-v2"
	rubricVersionV3 = "ga-rubric-20260523-v3"
)

var rubricAxesV1 = []string{
	"scenario_fit",
	"promisegrid_alignment",
	"auditability",
	"evolution_safety",
	"layer_boundary_clarity",
	"failure_handling",
	"implementation_plausibility",
	"risk_penalty",
}

var rubricAxesV2 = []string{
	"scenario_fit",
	"promisegrid_alignment",
	"auditability",
	"evolution_safety",
	"layer_boundary_clarity",
	"failure_handling",
	"implementation_plausibility",
	"promise_vocabulary",
	"simplicity_durability",
	"risk_penalty",
}

var rubricAxesV3 = append([]string(nil), rubricAxesV2...)

var promiseTheoryRulesV1 = []string{
	"Agents are autonomous.",
	"A promise is a scoped declaration of intent.",
	"No agent can make a promise on behalf of another agent.",
	"Promises do not guarantee outcomes.",
	"Trust is a local assessment of whether a promise will be kept.",
	"Promises to receive or use are not equivalent to obligations, impositions, or promises to give.",
}

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
	Promotion    *PromotionInfo `json:"promotion,omitempty"`
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
	OutputContract    string  `json:"output_contract,omitempty"`
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

// PromotionInfo records how a scored proposal result was later promoted into a
// canonical `results/` home without claiming that the historical `source.*`
// proposal provenance changed.
type PromotionInfo struct {
	PromotionDI                string `json:"promotion_di,omitempty"`
	RunGroupID                 string `json:"run_group_id,omitempty"`
	OriginalChildSimID         string `json:"original_child_sim_id,omitempty"`
	FinalSimID                 string `json:"final_sim_id,omitempty"`
	OriginalProposalSimPath    string `json:"original_proposal_sim_path,omitempty"`
	OriginalProposalResultPath string `json:"original_proposal_result_path,omitempty"`
	CanonicalResultPath        string `json:"canonical_result_path,omitempty"`
	SourceProvenancePolicy     string `json:"source_provenance_policy,omitempty"`
}

type RubricInfo struct {
	RubricVersion      string            `json:"rubric_version"`
	ScoreScale         string            `json:"score_scale"`
	ScoreMeanings      map[string]string `json:"score_meanings"`
	Axes               []string          `json:"axes"`
	PromiseTheoryRules []string          `json:"promise_theory_rules,omitempty"`
	PromiseTheoryRefs  []string          `json:"promise_theory_references,omitempty"`
}

// FitnessScores is serialized as the durable score evidence. Required axes must
// be emitted even when the score is zero, because zero is a valid judgment
// while an absent JSON field is a schema defect.
//
// Intent: Preserve score-field presence after the SIM-suzuf rerun exposed
// zero-valued fields being omitted by JSON tags. Source: DI-vonot; DI-movur
type FitnessScores struct {
	ScenarioFit                int `json:"scenario_fit"`
	PromiseGridAlignment       int `json:"promisegrid_alignment"`
	Auditability               int `json:"auditability"`
	EvolutionSafety            int `json:"evolution_safety"`
	LayerBoundaryClarity       int `json:"layer_boundary_clarity"`
	FailureHandling            int `json:"failure_handling"`
	ImplementationPlausibility int `json:"implementation_plausibility"`
	PromiseVocabulary          int `json:"promise_vocabulary"`
	SimplicityDurability       int `json:"simplicity_durability"`
	RiskPenalty                int `json:"risk_penalty"`
}

type FitnessSummary struct {
	Raw              float64 `json:"raw"`
	Normalized0To100 float64 `json:"normalized_0_100"`
	Confidence0To1   float64 `json:"confidence_0_1"`
}

type PTGateRuleAssessment struct {
	Status string `json:"status"`
	Note   string `json:"note"`
}

type PTGate struct {
	Status                 string               `json:"status"`
	AutonomousAgents       PTGateRuleAssessment `json:"autonomous_agents"`
	ScopedIntent           PTGateRuleAssessment `json:"scoped_intent"`
	NoPromisesForOthers    PTGateRuleAssessment `json:"no_promises_for_others"`
	NoGuaranteedOutcomes   PTGateRuleAssessment `json:"no_guaranteed_outcomes"`
	LocalTrustAssessment   PTGateRuleAssessment `json:"local_trust_assessment"`
	AcceptUseNotObligation PTGateRuleAssessment `json:"accept_use_not_obligation"`
	Violations             []string             `json:"violations"`
}

type Assessment struct {
	Rationale         string   `json:"rationale"`
	Strengths         []string `json:"strengths"`
	Weaknesses        []string `json:"weaknesses"`
	Risks             []string `json:"risks"`
	OpenQuestions     []string `json:"open_questions"`
	AuthorityBoundary string   `json:"authority_boundary"`
	PTGate            PTGate   `json:"pt_gate,omitempty"`
}

func knownResultSchemas() []string {
	return []string{resultSchemaV1, resultSchemaV2, resultSchemaV3}
}

func isKnownResultSchema(schema string) bool {
	for _, candidate := range knownResultSchemas() {
		if schema == candidate {
			return true
		}
	}
	return false
}

func expectedResultSchemaMessage() string {
	return "schema must be " + resultSchemaV1 + ", " + resultSchemaV2 + ", or " + resultSchemaV3
}

func rubricVersionForSchema(schema string) string {
	if schema == resultSchemaV3 {
		return rubricVersionV3
	}
	if schema == resultSchemaV2 {
		return rubricVersionV2
	}
	return rubricVersionV1
}

func rubricAxesForSchema(schema string) []string {
	if schema == resultSchemaV3 {
		return append([]string(nil), rubricAxesV3...)
	}
	if schema == resultSchemaV2 {
		return append([]string(nil), rubricAxesV2...)
	}
	return append([]string(nil), rubricAxesV1...)
}

func rubricScoreMeaningsForSchema(schema string) map[string]string {
	meanings := map[string]string{
		"0":            "no fit or absent",
		"5":            "strong fit",
		"risk_penalty": "0 low risk, 5 severe risk",
	}
	if schema == resultSchemaV2 || schema == resultSchemaV3 {
		// Intent: Keep stored rubric descriptions aligned with the scorer's
		// layer-local Promise Theory interpretation, so envelope-layer promises
		// are not mislabeled as weak merely because higher-layer accounting lives
		// inside the payload protocol. Source: DI-pozom
		meanings["promise_vocabulary"] = "0 drifts into claims/profiles/central trust-ledger framing, 5 stays promise-first, layer-local, and pCID-specific"
		meanings["simplicity_durability"] = "0 overbuilt or fragile, 5 minimal, durable, and small-device-friendly under the 100-year goal"
	}
	return meanings
}

func rubricPromiseTheoryRulesForSchema(schema string) []string {
	if schema != resultSchemaV3 {
		return nil
	}
	return append([]string(nil), promiseTheoryRulesV1...)
}

func rubricPromiseTheoryReferencesForSchema(schema string) []string {
	if schema != resultSchemaV3 {
		return nil
	}
	return []string{
		"Mark Burgess, In Search of Certainty",
		"Mark Burgess, Promise Theory: Principles and Applications",
		"Mark Burgess, Thinking in Promises",
	}
}

func scoreAxesForResult(result FitnessResult) []string {
	if result.Schema != "" {
		return rubricAxesForSchema(result.Schema)
	}
	if len(result.Rubric.Axes) > 0 {
		return append([]string(nil), result.Rubric.Axes...)
	}
	return rubricAxesForSchema(resultSchemaV1)
}

func (scores FitnessScores) axisValue(name string) int {
	switch name {
	case "scenario_fit":
		return scores.ScenarioFit
	case "promisegrid_alignment":
		return scores.PromiseGridAlignment
	case "auditability":
		return scores.Auditability
	case "evolution_safety":
		return scores.EvolutionSafety
	case "layer_boundary_clarity":
		return scores.LayerBoundaryClarity
	case "failure_handling":
		return scores.FailureHandling
	case "implementation_plausibility":
		return scores.ImplementationPlausibility
	case "promise_vocabulary":
		return scores.PromiseVocabulary
	case "simplicity_durability":
		return scores.SimplicityDurability
	case "risk_penalty":
		return scores.RiskPenalty
	default:
		return 0
	}
}

const (
	ptGateStatusClean         = "pt_clean"
	ptGateStatusReframeNeeded = "pt_reframe_needed"
	ptGateStatusInvalid       = "pt_invalid"
	ptRuleStatusPass          = "pass"
	ptRuleStatusWarning       = "warning"
	ptRuleStatusFail          = "fail"
)

// applyPTGateScorePolicy enforces the PT gate before deterministic weighted
// fitness is computed so non-PT designs cannot float to the top on technical
// neatness alone.
//
// Intent: Make the lugag PT correction operational in score ranking instead of
// leaving it as prose-only reviewer guidance. Source: DI-movur
func applyPTGateScorePolicy(scores FitnessScores, gate PTGate) FitnessScores {
	switch gate.Status {
	case ptGateStatusInvalid:
		scores.PromiseGridAlignment = 0
		scores.PromiseVocabulary = 0
		if scores.RiskPenalty < 5 {
			scores.RiskPenalty = 5
		}
	case ptGateStatusReframeNeeded:
		if scores.PromiseGridAlignment > 2 {
			scores.PromiseGridAlignment = 2
		}
		if scores.PromiseVocabulary > 2 {
			scores.PromiseVocabulary = 2
		}
		if scores.RiskPenalty < 3 {
			scores.RiskPenalty = 3
		}
	}
	return scores
}

// deterministicFitnessSummary keeps v2 fitness normalization in the runner so
// rubric expansion changes comparison math exactly once in audited code instead
// of depending on provider-specific prose interpretation. Source: DI-roruj
func deterministicFitnessSummary(schema string, scores FitnessScores, confidence float64) FitnessSummary {
	axes := rubricAxesForSchema(schema)
	if len(axes) == 0 {
		return FitnessSummary{Confidence0To1: confidence}
	}
	raw := 0.0
	for _, axis := range axes {
		value := scores.axisValue(axis)
		if axis == "risk_penalty" {
			raw += float64(5 - value)
			continue
		}
		raw += float64(value)
	}
	maxRaw := float64(len(axes) * 5)
	normalized := 0.0
	if maxRaw > 0 {
		normalized = raw / maxRaw * 100
	}
	return FitnessSummary{
		Raw:              raw,
		Normalized0To100: normalized,
		Confidence0To1:   confidence,
	}
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
