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

type rawResultPresence struct {
	Schema     string                     `json:"schema"`
	Scores     map[string]json.RawMessage `json:"scores"`
	Assessment map[string]json.RawMessage `json:"assessment"`
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

func findCanonicalResultFiles(repo Repo, model string, timestamp string) ([]string, error) {
	return findResultFilesUnderRoots(repo, []string{"results"}, model, timestamp)
}

// findResultFiles discovers only GA JSON result artifacts, including ignored
// proposal-stage child score evidence. Markdown canary files are intentionally
// invisible to GA-runner validation and future scoring logic.
//
// Intent: Keep old matrix-runner canaries as historical evidence without letting
// them contaminate GA fitness selection. Source: DI-ramar; DI-pobus; DI-lirat
func findResultFiles(repo Repo, model string, timestamp string) ([]string, error) {
	return findResultFilesUnderRoots(repo, []string{"results", "proposals"}, model, timestamp)
}

func findResultFilesUnderRoots(repo Repo, rootNames []string, model string, timestamp string) ([]string, error) {
	var paths []string
	for _, rootName := range rootNames {
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
	var presence rawResultPresence
	if err := json.Unmarshal(bytes, &presence); err != nil {
		return append(issues, "invalid JSON presence scan: "+err.Error())
	}
	issues = append(issues, validateResultIdentity(repo, path, parts, result)...)
	issues = append(issues, validateRunner(result.Runner)...)
	issues = append(issues, validateSource(result.Source)...)
	issues = append(issues, validateRubric(result.Schema, result.Rubric)...)
	issues = append(issues, validateScores(result.Schema, presence.Scores, result.Scores)...)
	issues = append(issues, validateFitness(result.Fitness)...)
	issues = append(issues, validateAssessmentPresence(result.Schema, presence.Assessment)...)
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
	if !isKnownResultSchema(result.Schema) {
		issues = append(issues, expectedResultSchemaMessage())
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
	var issues []string
	if runner.Tool == "" {
		issues = append(issues, "runner.tool is required")
	}
	if runner.OutputContract != "" {
		if _, err := normalizeOutputContract(runner.OutputContract); err != nil {
			issues = append(issues, "runner."+err.Error())
		}
	}
	return issues
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

func validateRubric(schema string, rubric RubricInfo) []string {
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
	if schema == resultSchemaV2 {
		if rubric.RubricVersion != rubricVersionV2 {
			issues = append(issues, "rubric.rubric_version must be "+rubricVersionV2+" for "+resultSchemaV2)
		}
		expectedAxes := rubricAxesForSchema(schema)
		if strings.Join(rubric.Axes, ",") != strings.Join(expectedAxes, ",") {
			issues = append(issues, "rubric.axes must match "+strings.Join(expectedAxes, ",")+" for "+resultSchemaV2)
		}
	}
	if schema == resultSchemaV3 {
		if rubric.RubricVersion != rubricVersionV3 {
			issues = append(issues, "rubric.rubric_version must be "+rubricVersionV3+" for "+resultSchemaV3)
		}
		expectedAxes := rubricAxesForSchema(schema)
		if strings.Join(rubric.Axes, ",") != strings.Join(expectedAxes, ",") {
			issues = append(issues, "rubric.axes must match "+strings.Join(expectedAxes, ",")+" for "+resultSchemaV3)
		}
		if strings.Join(rubric.PromiseTheoryRules, "\n") != strings.Join(rubricPromiseTheoryRulesForSchema(schema), "\n") {
			issues = append(issues, "rubric.promise_theory_rules must match the canonical PT rule list for "+resultSchemaV3)
		}
		if len(rubric.PromiseTheoryRefs) == 0 {
			issues = append(issues, "rubric.promise_theory_references must not be empty for "+resultSchemaV3)
		}
	}
	// Intent: Keep V4 validation append-only and explicit so new layer-aware
	// results cannot masquerade as V3 or omit the envelope/kernel/app axes that
	// make the rubric expansion meaningful. Source: DI-ripuz
	if schema == resultSchemaV4 {
		if rubric.RubricVersion != rubricVersionV4 {
			issues = append(issues, "rubric.rubric_version must be "+rubricVersionV4+" for "+resultSchemaV4)
		}
		expectedAxes := rubricAxesForSchema(schema)
		if strings.Join(rubric.Axes, ",") != strings.Join(expectedAxes, ",") {
			issues = append(issues, "rubric.axes must match "+strings.Join(expectedAxes, ",")+" for "+resultSchemaV4)
		}
		if strings.Join(rubric.PromiseTheoryRules, "\n") != strings.Join(rubricPromiseTheoryRulesForSchema(schema), "\n") {
			issues = append(issues, "rubric.promise_theory_rules must match the canonical PT rule list for "+resultSchemaV4)
		}
		if len(rubric.PromiseTheoryRefs) == 0 {
			issues = append(issues, "rubric.promise_theory_references must not be empty for "+resultSchemaV4)
		}
	}
	return issues
}

func validateScores(schema string, rawScores map[string]json.RawMessage, scores FitnessScores) []string {
	scoreMap := map[string]int{
		"scenario_fit":                   scores.ScenarioFit,
		"promisegrid_alignment":          scores.PromiseGridAlignment,
		"auditability":                   scores.Auditability,
		"evolution_safety":               scores.EvolutionSafety,
		"layer_boundary_clarity":         scores.LayerBoundaryClarity,
		"failure_handling":               scores.FailureHandling,
		"implementation_plausibility":    scores.ImplementationPlausibility,
		"promise_vocabulary":             scores.PromiseVocabulary,
		"simplicity_durability":          scores.SimplicityDurability,
		"envelope_discipline":            scores.EnvelopeDiscipline,
		"kernel_implementation_promises": scores.KernelImplementationPromises,
		"app_protocol_promise_semantics": scores.AppProtocolPromiseSemantics,
		"risk_penalty":                   scores.RiskPenalty,
	}
	var issues []string
	for name, value := range scoreMap {
		if value < 0 || value > 5 {
			issues = append(issues, fmt.Sprintf("scores.%s must be between 0 and 5", name))
		}
	}
	if schema == resultSchemaV2 {
		for _, field := range []string{"promise_vocabulary", "simplicity_durability"} {
			if _, ok := rawScores[field]; !ok {
				issues = append(issues, "scores."+field+" is required for "+resultSchemaV2)
			}
		}
	}
	if schema == resultSchemaV3 {
		for _, field := range []string{"promise_vocabulary", "simplicity_durability"} {
			if _, ok := rawScores[field]; !ok {
				issues = append(issues, "scores."+field+" is required for "+resultSchemaV3)
			}
		}
	}
	// Intent: Require every V4 axis to be present even when a score is zero, so
	// additive rescoring distinguishes absent provider output from an explicit
	// low judgment for a layer-specific axis. Source: DI-ripuz
	if schema == resultSchemaV4 {
		for _, field := range []string{
			"promise_vocabulary",
			"simplicity_durability",
			"envelope_discipline",
			"kernel_implementation_promises",
			"app_protocol_promise_semantics",
		} {
			if _, ok := rawScores[field]; !ok {
				issues = append(issues, "scores."+field+" is required for "+resultSchemaV4)
			}
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

func validateAssessmentPresence(schema string, rawAssessment map[string]json.RawMessage) []string {
	var issues []string
	if schema == resultSchemaV3 || schema == resultSchemaV4 {
		if _, ok := rawAssessment["pt_gate"]; !ok {
			issues = append(issues, "assessment.pt_gate is required for "+schema)
		}
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
	if assessment.PTGate.Status != "" {
		issues = append(issues, validatePTGate(assessment.PTGate)...)
	}
	return issues
}

func validatePTGate(gate PTGate) []string {
	var issues []string
	switch gate.Status {
	case ptGateStatusClean, ptGateStatusReframeNeeded, ptGateStatusInvalid:
	default:
		issues = append(issues, "assessment.pt_gate.status must be pt_clean, pt_reframe_needed, or pt_invalid")
	}
	rules := map[string]PTGateRuleAssessment{
		"autonomous_agents":         gate.AutonomousAgents,
		"scoped_intent":             gate.ScopedIntent,
		"no_promises_for_others":    gate.NoPromisesForOthers,
		"no_guaranteed_outcomes":    gate.NoGuaranteedOutcomes,
		"local_trust_assessment":    gate.LocalTrustAssessment,
		"accept_use_not_obligation": gate.AcceptUseNotObligation,
	}
	for name, rule := range rules {
		switch rule.Status {
		case ptRuleStatusPass, ptRuleStatusWarning, ptRuleStatusFail:
		default:
			issues = append(issues, "assessment.pt_gate."+name+".status must be pass, warning, or fail")
		}
		if rule.Note == "" {
			issues = append(issues, "assessment.pt_gate."+name+".note is required")
		}
	}
	sort.Strings(issues)
	return issues
}

const defaultBackfillCleanEnvelopeCount = 6

type auditRecord struct {
	Path              string
	Result            FitnessResult
	ExactMatch        bool
	RootContractDrift bool
	SourceResolution  string
	VocabularyStatus  string
	VocabularyReasons []string
}

type auditSummary struct {
	SimID             string
	ResultCount       int
	ExactMatchCount   int
	AverageNormalized float64
	SourceResolution  string
	VocabularyStatus  string
	RootContractDrift bool
	Models            []string
}

type backfillSelection struct {
	Records             []auditRecord
	HardHitSimIDs       []string
	CleanEnvelopeSimIDs []string
}

func runAudit(args []string, stdout io.Writer) error {
	// Intent: Make rubric-v2 backfill selection cheap and reviewable before any
	// new provider calls spend money on a broad rescore. Source: DI-roruj
	fs := flag.NewFlagSet("audit", flag.ContinueOnError)
	repoRoot := commonRepoFlag(fs)
	model := fs.String("model", "", "optional canonical result model filter")
	timestamp := fs.String("timestamp", "", "optional canonical result timestamp filter")
	cleanEnvelopeCount := fs.Int("clean-envelope-count", defaultBackfillCleanEnvelopeCount, "number of clean grid-envelope sims to surface as calibration targets")
	if err := fs.Parse(args); err != nil {
		return err
	}
	repo, err := openRepo(*repoRoot)
	if err != nil {
		return err
	}
	records, err := auditCanonicalV1Results(repo, *model, *timestamp)
	if err != nil {
		return err
	}
	selection := selectTargetedBackfill(records, *cleanEnvelopeCount)
	summaries := summarizeAuditRecords(records)
	exactMatches := 0
	rootDrift := 0
	historicalSources := 0
	canonicalFallbackSources := 0
	missingSources := 0
	for _, record := range records {
		if record.ExactMatch {
			exactMatches++
		}
		if record.RootContractDrift {
			rootDrift++
		}
		switch record.SourceResolution {
		case auditSourceResolutionHistorical:
			historicalSources++
		case auditSourceResolutionCanonicalFallback:
			canonicalFallbackSources++
		default:
			missingSources++
		}
	}
	if err := writeFormat(stdout, "audited=%d exact_match=%d mismatch=%d root_contract_drift=%d source_historical=%d source_canonical_fallback=%d source_missing=%d hard_hit_sims=%d clean_envelope_candidates=%d\n",
		len(records),
		exactMatches,
		len(records)-exactMatches,
		rootDrift,
		historicalSources,
		canonicalFallbackSources,
		missingSources,
		len(selection.HardHitSimIDs),
		len(selection.CleanEnvelopeSimIDs)); err != nil {
		return err
	}
	if len(selection.HardHitSimIDs) > 0 {
		if err := writeLine(stdout, "hard_hit sims:"); err != nil {
			return err
		}
		for _, simID := range selection.HardHitSimIDs {
			summary := summaries[simID]
			if err := writeAuditSummaryLine(stdout, summary); err != nil {
				return err
			}
		}
	}
	if len(selection.CleanEnvelopeSimIDs) > 0 {
		if err := writeLine(stdout, "clean envelope calibration sims:"); err != nil {
			return err
		}
		for _, simID := range selection.CleanEnvelopeSimIDs {
			summary := summaries[simID]
			if err := writeAuditSummaryLine(stdout, summary); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeAuditSummaryLine(stdout io.Writer, summary auditSummary) error {
	return writeFormat(stdout,
		"  %s results=%d exact_match=%d avg_normalized_0_100=%.2f source_resolution=%s vocab=%s root_contract_drift=%t models=%s\n",
		summary.SimID,
		summary.ResultCount,
		summary.ExactMatchCount,
		summary.AverageNormalized,
		summary.SourceResolution,
		summary.VocabularyStatus,
		summary.RootContractDrift,
		strings.Join(summary.Models, ","),
	)
}

func auditCanonicalV1Results(repo Repo, model string, timestamp string) ([]auditRecord, error) {
	paths, err := findCanonicalResultFiles(repo, model, timestamp)
	if err != nil {
		return nil, err
	}
	var records []auditRecord
	for _, path := range paths {
		result, err := readFitnessResult(path)
		if err != nil {
			return nil, err
		}
		if result.Schema != resultSchemaV1 {
			continue
		}
		sourceResolution := resolveAuditSimulationSource(repo, repo.Rel(path), result)
		exactMatch, rootContractDrift := auditResultSourceMatch(repo, result, sourceResolution)
		vocabularyStatus, vocabularyReasons, err := auditSimulationVocabulary(repo, sourceResolution, result.SimID)
		if err != nil {
			return nil, err
		}
		records = append(records, auditRecord{
			Path:              repo.Rel(path),
			Result:            result,
			ExactMatch:        exactMatch,
			RootContractDrift: rootContractDrift,
			SourceResolution:  sourceResolution.Mode,
			VocabularyStatus:  vocabularyStatus,
			VocabularyReasons: vocabularyReasons,
		})
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("no canonical %s results matched selection", resultSchemaV1)
	}
	sort.Slice(records, func(i, j int) bool {
		return records[i].Path < records[j].Path
	})
	return records, nil
}

const (
	auditSourceResolutionHistorical        = "historical"
	auditSourceResolutionCanonicalFallback = "canonical_fallback"
	auditSourceResolutionMissing           = "missing"
	auditSourceResolutionMixed             = "mixed"
)

type auditSourceState struct {
	Mode              string
	HistoricalSimPath string
	ActiveSimPath     string
	CanonicalSimPath  string
}

// resolveAuditSimulationSource keeps historical source provenance intact while
// letting audit/backfill compare promoted canonical results against the current
// canonical sim tree when the original proposal tree has been deleted.
//
// Intent: Preserve append-only historical `source.*` fields while unblocking
// audit/backfill on promoted canonical results that intentionally still point at
// deleted proposal paths. Source: DI-zobur
func resolveAuditSimulationSource(repo Repo, auditedResultPath string, result FitnessResult) auditSourceState {
	historicalSimPath := strings.TrimSuffix(normalizeRelPath(result.Source.SimPath), "/")
	if info, err := os.Stat(repo.Abs(historicalSimPath)); err == nil && info.IsDir() {
		return auditSourceState{
			Mode:              auditSourceResolutionHistorical,
			HistoricalSimPath: historicalSimPath,
			ActiveSimPath:     historicalSimPath,
		}
	}
	if result.Promotion == nil {
		return auditSourceState{
			Mode:              auditSourceResolutionMissing,
			HistoricalSimPath: historicalSimPath,
		}
	}
	if normalizeRelPath(result.Promotion.CanonicalResultPath) != normalizeRelPath(auditedResultPath) {
		return auditSourceState{
			Mode:              auditSourceResolutionMissing,
			HistoricalSimPath: historicalSimPath,
		}
	}
	if !strings.HasPrefix(historicalSimPath, "proposals/") {
		return auditSourceState{
			Mode:              auditSourceResolutionMissing,
			HistoricalSimPath: historicalSimPath,
		}
	}
	finalSimID := strings.TrimSpace(result.Promotion.FinalSimID)
	if finalSimID == "" {
		finalSimID = result.SimID
	}
	canonicalSimPath := filepath.ToSlash(filepath.Join("simulations", finalSimID))
	if info, err := os.Stat(repo.Abs(canonicalSimPath)); err == nil && info.IsDir() {
		return auditSourceState{
			Mode:              auditSourceResolutionCanonicalFallback,
			HistoricalSimPath: historicalSimPath,
			ActiveSimPath:     canonicalSimPath,
			CanonicalSimPath:  canonicalSimPath,
		}
	}
	return auditSourceState{
		Mode:              auditSourceResolutionMissing,
		HistoricalSimPath: historicalSimPath,
	}
}

// auditPathRelativeToRoot preserves the relative suffix of a historical
// sim-root file so canonical fallback can compare the same logical file under a
// different root.
func auditPathRelativeToRoot(path string, root string) (string, bool) {
	cleanPath := normalizeRelPath(path)
	cleanRoot := strings.TrimSuffix(normalizeRelPath(root), "/")
	prefix := cleanRoot + "/"
	if !strings.HasPrefix(cleanPath, prefix) {
		return "", false
	}
	return strings.TrimPrefix(cleanPath, prefix), true
}

// remapAuditSourcePath rewrites only sim-root source paths during canonical
// fallback. Root contracts and scenario files keep their historical stored
// paths.
func remapAuditSourcePath(path string, source auditSourceState) string {
	cleanPath := normalizeRelPath(path)
	if source.Mode != auditSourceResolutionCanonicalFallback {
		return cleanPath
	}
	relativePath, ok := auditPathRelativeToRoot(cleanPath, source.HistoricalSimPath)
	if !ok {
		return cleanPath
	}
	return filepath.ToSlash(filepath.Join(source.CanonicalSimPath, relativePath))
}

func auditResultSourceMatch(repo Repo, result FitnessResult, source auditSourceState) (bool, bool) {
	// Intent: Gate targeted backfill on current sim/scenario bytes while
	// reporting root-contract drift separately, because rubric-v2 updates will
	// naturally change shared scoring docs without changing the candidate sim
	// itself. Promoted canonical results may compare against the current canonical
	// sim tree when historical proposal paths are gone, but only as an exact-byte
	// fallback. Source: DI-roruj; DI-zobur
	rootContracts := map[string]bool{}
	for _, path := range result.Source.RootContractPaths {
		rootContracts[normalizeRelPath(path)] = true
	}
	exactMatch := true
	rootContractDrift := false
	for _, file := range result.Source.Files {
		storedPath := normalizeRelPath(file.Path)
		currentPath := remapAuditSourcePath(storedPath, source)
		hash, err := sha256File(repo, currentPath)
		if err != nil {
			if rootContracts[storedPath] {
				rootContractDrift = true
				continue
			}
			exactMatch = false
			continue
		}
		if hash == file.SHA256 {
			continue
		}
		if rootContracts[storedPath] {
			rootContractDrift = true
			continue
		}
		exactMatch = false
	}
	if source.Mode == auditSourceResolutionMissing {
		return false, rootContractDrift
	}
	currentTreeHash, err := currentSimulationTreeHash(repo, source.ActiveSimPath)
	if err != nil || currentTreeHash != result.Source.SimulationTreeHash {
		exactMatch = false
	}
	return exactMatch, rootContractDrift
}

func auditSimulationVocabulary(repo Repo, source auditSourceState, simID string) (string, []string, error) {
	// Intent: Classify only the current sim's own vocabulary drift so the audit
	// can target clearly affected families first and leave broader clean sims for
	// calibration sampling. Promoted canonical results audit the current canonical
	// sim docs when the original proposal tree no longer exists. Source: DI-roruj;
	// DI-zobur
	if source.Mode == auditSourceResolutionMissing {
		return "soft_hit", []string{"simulation source unresolved for current audit"}, nil
	}
	paths, err := simulationAuditPaths(repo, source.ActiveSimPath)
	if err != nil {
		return "", nil, err
	}
	docs, err := sourceDocumentsFromPaths(repo, paths)
	if err != nil {
		return "", nil, err
	}
	var body strings.Builder
	for _, doc := range docs {
		body.WriteString("\n")
		body.WriteString(strings.ToLower(doc.Text))
	}
	text := body.String()
	simIDLower := strings.ToLower(simID)
	isUDPConformanceFixture := strings.Contains(simIDLower, "udp-feed-v0-conformance")
	var hardReasons []string
	if strings.Contains(simIDLower, "claim-card") {
		hardReasons = append(hardReasons, "sim-id uses claim-card artifact vocabulary")
	}
	if strings.Contains(simIDLower, "boundary-claim") {
		hardReasons = append(hardReasons, "sim-id uses boundary-claim artifact vocabulary")
	}
	if strings.Contains(simIDLower, "conformance-citation") {
		hardReasons = append(hardReasons, "sim-id uses conformance-citation artifact vocabulary")
	}
	if strings.Contains(simIDLower, "conformance") && !isUDPConformanceFixture {
		hardReasons = append(hardReasons, "sim-id uses conformance-family naming outside the allowed UDP fixture")
	}
	for phrase, reason := range map[string]string{
		"claim card":            "docs use claim-card artifact vocabulary",
		"claim cards":           "docs use claim-card artifact vocabulary",
		"claim_header":          "docs use generic claim-header vocabulary",
		"claim header":          "docs use generic claim-header vocabulary",
		"statement_capsule":     "docs use universal statement-capsule vocabulary",
		"statement capsule":     "docs use universal statement-capsule vocabulary",
		"conformance bundle":    "docs use conformance-bundle artifact vocabulary",
		"profile support claim": "docs use profile-support-claim vocabulary",
		"port claim":            "docs use port-claim vocabulary",
		"trust ledger":          "docs use trust-ledger vocabulary",
		"boundary claim":        "docs use boundary-claim artifact vocabulary",
		"conformance citation":  "docs use conformance-citation artifact vocabulary",
	} {
		if strings.Contains(text, phrase) {
			hardReasons = append(hardReasons, reason)
		}
	}
	hasEnvSelector := strings.Contains(text, "env_pcid")
	hasSigSelector := strings.Contains(text, "sig_pcid")
	hasPayloadSelector := strings.Contains(text, "payload_pcid")
	if hasEnvSelector && (hasSigSelector || hasPayloadSelector) {
		hardReasons = append(hardReasons, "docs use envelope selector-stack vocabulary")
	}
	if hasSigSelector && hasPayloadSelector && (strings.Contains(text, "claim_header") || strings.Contains(text, "claim header") || strings.Contains(text, "statement_capsule") || strings.Contains(text, "statement capsule")) {
		hardReasons = append(hardReasons, "docs use nested selector-stack vocabulary")
	}
	if len(hardReasons) > 0 {
		return "hard_hit", uniqueStrings(hardReasons), nil
	}
	cleaned := text
	for _, allowed := range []string{
		"payload conforms to the protocol specification referred to by this pcid",
		"payload conforms to pcid",
		"alice promises this payload meets the protocol specification referred to by this pcid",
		"udp-feed v0 conformance",
	} {
		cleaned = strings.ReplaceAll(cleaned, allowed, " ")
	}
	var softReasons []string
	for _, pattern := range []struct {
		substring string
		reason    string
	}{
		{substring: "claim", reason: "docs still use claim vocabulary"},
		{substring: "conformance", reason: "docs still use conformance vocabulary"},
		{substring: "profile", reason: "docs still use profile vocabulary"},
		{substring: "trust ledger", reason: "docs still use trust-ledger vocabulary"},
	} {
		if strings.Contains(cleaned, pattern.substring) {
			softReasons = append(softReasons, pattern.reason)
		}
	}
	if len(softReasons) > 0 {
		return "soft_hit", uniqueStrings(softReasons), nil
	}
	return "clean", nil, nil
}

func simulationAuditPaths(repo Repo, simPath string) ([]string, error) {
	cleanSimPath := strings.TrimSuffix(normalizeRelPath(simPath), "/")
	paths := []string{filepath.ToSlash(filepath.Join(cleanSimPath, "README.md"))}
	questionPath := filepath.ToSlash(filepath.Join(cleanSimPath, "QUESTION.md"))
	if info, err := os.Stat(repo.Abs(questionPath)); err == nil && !info.IsDir() {
		paths = append(paths, questionPath)
	}
	localSim, err := localMarkdownFiles(repo, cleanSimPath, map[string]bool{
		"README.md":   true,
		"QUESTION.md": true,
	})
	if err != nil {
		return nil, err
	}
	paths = append(paths, localSim...)
	return uniqueStrings(paths), nil
}

func summarizeAuditRecords(records []auditRecord) map[string]auditSummary {
	type aggregate struct {
		summary   auditSummary
		sum       float64
		modelSeen map[string]bool
	}
	aggregates := map[string]*aggregate{}
	for _, record := range records {
		item, ok := aggregates[record.Result.SimID]
		if !ok {
			item = &aggregate{
				summary: auditSummary{
					SimID:            record.Result.SimID,
					SourceResolution: record.SourceResolution,
					VocabularyStatus: record.VocabularyStatus,
				},
				modelSeen: map[string]bool{},
			}
			aggregates[record.Result.SimID] = item
		}
		item.summary.ResultCount++
		if record.ExactMatch {
			item.summary.ExactMatchCount++
			item.sum += record.Result.Fitness.Normalized0To100
		}
		if vocabularySeverity(record.VocabularyStatus) > vocabularySeverity(item.summary.VocabularyStatus) {
			item.summary.VocabularyStatus = record.VocabularyStatus
		}
		item.summary.SourceResolution = mergeAuditSourceResolution(item.summary.SourceResolution, record.SourceResolution)
		if record.RootContractDrift {
			item.summary.RootContractDrift = true
		}
		if !item.modelSeen[record.Result.ModelID] {
			item.modelSeen[record.Result.ModelID] = true
			item.summary.Models = append(item.summary.Models, record.Result.ModelID)
		}
	}
	summaries := map[string]auditSummary{}
	for simID, item := range aggregates {
		sort.Strings(item.summary.Models)
		if item.summary.ExactMatchCount > 0 {
			item.summary.AverageNormalized = item.sum / float64(item.summary.ExactMatchCount)
		}
		summaries[simID] = item.summary
	}
	return summaries
}

// mergeAuditSourceResolution keeps grouped per-sim audit summaries readable if
// the corpus ever contains mixed historical and canonical-fallback records for
// the same canonical sim ID.
func mergeAuditSourceResolution(current string, next string) string {
	if current == "" {
		return next
	}
	if current == next {
		return current
	}
	return auditSourceResolutionMixed
}

func selectTargetedBackfill(records []auditRecord, cleanEnvelopeCount int) backfillSelection {
	// Intent: Spend the first v2 rescoring budget on the sims most likely to move
	// under the new vocabulary rules, then add a small clean envelope slice so
	// the rerun still shows whether the top wire-format contenders stay stable.
	// When the historical corpus has repeated exact-match v1 rows for the same
	// sim/scenario pair, future targeted backfill states should queue one
	// deterministic winner rather than emit repeated v2 cells that all point at
	// the same result path. Source: DI-roruj; DI-guhar
	summaries := summarizeAuditRecords(records)
	var hardHit []auditSummary
	var cleanEnvelope []auditSummary
	for _, summary := range summaries {
		if summary.ExactMatchCount == 0 {
			continue
		}
		if summary.VocabularyStatus == "hard_hit" {
			hardHit = append(hardHit, summary)
			continue
		}
		if summary.VocabularyStatus == "clean" && strings.Contains(strings.ToLower(summary.SimID), "grid-envelope") {
			cleanEnvelope = append(cleanEnvelope, summary)
		}
	}
	sort.Slice(hardHit, func(i, j int) bool {
		return hardHit[i].SimID < hardHit[j].SimID
	})
	sort.Slice(cleanEnvelope, func(i, j int) bool {
		if cleanEnvelope[i].AverageNormalized != cleanEnvelope[j].AverageNormalized {
			return cleanEnvelope[i].AverageNormalized > cleanEnvelope[j].AverageNormalized
		}
		return cleanEnvelope[i].SimID < cleanEnvelope[j].SimID
	})
	if cleanEnvelopeCount > 0 && len(cleanEnvelope) > cleanEnvelopeCount {
		cleanEnvelope = cleanEnvelope[:cleanEnvelopeCount]
	}
	selectedSimIDs := map[string]bool{}
	selection := backfillSelection{}
	for _, summary := range hardHit {
		selectedSimIDs[summary.SimID] = true
		selection.HardHitSimIDs = append(selection.HardHitSimIDs, summary.SimID)
	}
	for _, summary := range cleanEnvelope {
		selectedSimIDs[summary.SimID] = true
		selection.CleanEnvelopeSimIDs = append(selection.CleanEnvelopeSimIDs, summary.SimID)
	}
	for _, record := range records {
		if record.ExactMatch && selectedSimIDs[record.Result.SimID] {
			selection.Records = append(selection.Records, record)
		}
	}
	selection.Records = dedupeBackfillRecordsByPair(selection.Records)
	sort.Slice(selection.Records, func(i, j int) bool {
		return selection.Records[i].Path < selection.Records[j].Path
	})
	return selection
}

// dedupeBackfillRecordsByPair keeps one deterministic historical v1 candidate
// for each targeted sim/scenario pair before `backfill-init` materializes new
// v2 cells.
//
// Intent: Future backfill states should not queue duplicate work for the same
// v2 result home just because the historical corpus kept multiple exact-match
// v1 rows for that pair. Source: DI-guhar
func dedupeBackfillRecordsByPair(records []auditRecord) []auditRecord {
	byPair := map[string]auditRecord{}
	for _, record := range records {
		key := backfillPairKey(record.Result.SimID, record.Result.ScenarioID)
		current, ok := byPair[key]
		if !ok || preferComparisonAuditRecord(record, current) {
			byPair[key] = record
		}
	}
	var deduped []auditRecord
	for _, record := range byPair {
		deduped = append(deduped, record)
	}
	return deduped
}

func vocabularySeverity(status string) int {
	switch status {
	case "hard_hit":
		return 2
	case "soft_hit":
		return 1
	default:
		return 0
	}
}
