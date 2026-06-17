package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/decision"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/pcid"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/protocol"
)

// RunSummary is the deterministic JSON report emitted by poc15-analyze.
// Intent: Keep POC15 live runs comparable without committing provider logs or
// Docker volume state. Source: DI-timah
type RunSummary struct {
	RunDir                         string                  `json:"run_dir"`
	TotalEvents                    int                     `json:"total_events"`
	EventCounts                    map[string]int          `json:"event_counts"`
	OutcomeCounts                  map[string]int          `json:"outcome_counts"`
	AgentCounts                    map[string]int          `json:"agent_counts"`
	FailureCounts                  map[string]int          `json:"failure_counts"`
	ProtocolCounts                 map[string]int          `json:"protocol_counts"`
	TrustDrivenChoiceCounts        map[string]int          `json:"trust_driven_choice_counts"`
	EconomicsCounts                map[string]int          `json:"economics_counts"`
	VerificationCounts             map[string]int          `json:"verification_counts"`
	ReplicaRecoveryCounts          map[string]int          `json:"replica_recovery_counts"`
	RPCDriftCounts                 map[string]int          `json:"rpc_drift_counts"`
	ShippingCounts                 map[string]int          `json:"shipping_counts"`
	RelationshipTransitionCounts   map[string]int          `json:"relationship_transition_counts"`
	DynamicTopologyCounts          map[string]int          `json:"dynamic_topology_counts"`
	LocalResourceCounts            map[string]int          `json:"local_resource_counts"`
	ResourceTrustCouplingCounts    map[string]int          `json:"resource_trust_coupling_counts"`
	DurabilityCounts               map[string]int          `json:"durability_counts"`
	RetentionCounts                map[string]int          `json:"retention_counts"`
	PressureCounts                 map[string]int          `json:"pressure_counts"`
	ReplayCounts                   map[string]int          `json:"replay_counts"`
	TrustCautionCounts             map[string]int          `json:"trust_caution_counts"`
	RuntimeAdapterEventCounts      map[string]int          `json:"runtime_adapter_event_counts"`
	DecentralizedMonitorCounts     map[string]int          `json:"decentralized_monitor_counts"`
	RouteCounts                    map[string]int          `json:"route_counts"`
	ArrayPayloadProtocolCounts     map[string]int          `json:"array_payload_protocol_counts"`
	AgentCASCounts                 map[string]int          `json:"agent_cas_counts"`
	MessageArtifactCount           int                     `json:"message_artifact_count"`
	MessageCASObjectCount          int                     `json:"message_cas_object_count"`
	MessageDAGRecordCount          int                     `json:"message_dag_record_count"`
	MessageDAGNodeCount            int                     `json:"message_dag_node_count"`
	MessageDAGParentLinkCount      int                     `json:"message_dag_parent_link_count"`
	MessageDAGRootCount            int                     `json:"message_dag_root_count"`
	MessageDAGReachableCount       int                     `json:"message_dag_reachable_count"`
	MessageDAGMaxDepth             int                     `json:"message_dag_max_depth"`
	MessageDAGMissingParentCount   int                     `json:"message_dag_missing_parent_count"`
	AckMessageMissingParentCount   int                     `json:"ack_message_missing_parent_count"`
	MessageArtifactDirectionCounts map[string]int          `json:"message_artifact_direction_counts"`
	MessageArtifactProtocolCounts  map[string]int          `json:"message_artifact_protocol_counts"`
	MessageArtifactBadPrefixCount  int                     `json:"message_artifact_bad_prefix_count"`
	MessageDAGParentLocationCounts map[string]int          `json:"message_dag_parent_location_counts"`
	MessageShapeSpecimenCounts     map[string]int          `json:"message_shape_specimen_counts"`
	MigrationCounts                map[string]int          `json:"migration_counts"`
	RestartCounts                  map[string]int          `json:"restart_counts"`
	ForbiddenVocabularyCounts      map[string]int          `json:"forbidden_vocabulary_counts,omitempty"`
	ScoreReport                    ScoreReport             `json:"score_report"`
	ProductionFitness              ProductionFitnessReport `json:"production_fitness"`
	MissingRequiredEventNames      []string                `json:"missing_required_event_names,omitempty"`
	MonitorReport                  *decision.MonitorReport `json:"monitor_report,omitempty"`
}

// ScoreReport gives the operator a fast POC15 fitness view in addition to the
// raw event counters.
// Intent: Future POCs must not be able to look healthy while silently dropping
// inherited transport, storage, compute, economics, trust, verification, or
// replica-recovery event records. Source: DI-sinur
type ScoreReport struct {
	Overall        int      `json:"overall"`
	Transport      int      `json:"transport"`
	Storage        int      `json:"storage"`
	Compute        int      `json:"compute"`
	Economics      int      `json:"economics"`
	Trust          int      `json:"trust"`
	Verification   int      `json:"verification"`
	Replica        int      `json:"replica"`
	Durability     int      `json:"durability"`
	Retention      int      `json:"retention"`
	Pressure       int      `json:"pressure"`
	Replay         int      `json:"replay"`
	RuntimeAdapter int      `json:"runtime_adapter"`
	Monitoring     int      `json:"monitoring"`
	Route          int      `json:"route"`
	AgentCAS       int      `json:"agent_cas"`
	Migration      int      `json:"migration"`
	Restart        int      `json:"restart"`
	Concerns       []string `json:"concerns,omitempty"`
}

// ProductionFitnessReport states the current POC-to-production gap in terms an
// operator can compare across clean runs.
// Intent: POC15 should produce a concise production-fitness summary from
// analyzer and monitor event records without pretending a successful POC run is
// production readiness. Source: DI-sihuz
type ProductionFitnessReport struct {
	Baseline           string   `json:"baseline"`
	ReadyForProduction bool     `json:"ready_for_production"`
	Verdict            string   `json:"verdict"`
	BlockingGaps       []string `json:"blocking_gaps,omitempty"`
}

// messageDAGRecord mirrors the collector's run-scoped message DAG index.
// Intent: The analyzer validates actual binary `.cbor` artifacts by hash so a
// clean run cannot pass by logging message-like events without retaining the
// messages themselves. Source: DI-tuhop
type messageDAGRecord struct {
	Source             string `json:"source"`
	Observer           string `json:"observer"`
	Direction          string `json:"direction"`
	Peer               string `json:"peer"`
	Protocol           string `json:"protocol"`
	ExactSHA256        string `json:"exact_sha256"`
	ParentExactSHA256  string `json:"parent_exact_sha256,omitempty"`
	ParentLinkLocation string `json:"parent_link_location,omitempty"`
	PromiseAbout       string `json:"promise_about,omitempty"`
	SourceEvent        string `json:"source_event,omitempty"`
	Path               string `json:"path"`
}

// AcceptanceCriteria describes the event-record gates for a clean POC15 regression
// run.
// Intent: Keep the current clean-run expectations executable so a wrong log
// path, missing shipping workflow, lost non-commitment suppression, or renewed
// resource/trust coupling fails loudly instead of producing a plausible but
// empty report. Source: DI-jidah
type AcceptanceCriteria struct {
	MinTotalEvents                    int
	RequireMonitorReport              bool
	MinMonitorScore                   int
	RequireFulfillmentWorkflow        bool
	RequireAccountingUpdated          bool
	RequireAccountingDuplicate        bool
	RequireNotPromisedSuppression     bool
	RequireNoResourceTrustCoupling    bool
	RequireAccountingDuplicateConfirm bool
	RequireCASComputeSuperset         bool
	RequireRuntimeAdapterSuperset     bool
	RequireDecentralizedMonitoring    bool
	RequireAllKnownArrayPayloads      bool
	RequireMixedVersionMigration      bool
	RequireRunInternalRestart         bool
	RequireEventVocabulary            bool
	RequireRawMessageArtifacts        bool
	RequireMessageShapeSpecimens      bool
	RequirePersistentSessions         bool
	MinScoreOverall                   int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: poc15-analyze RUN_DIR\n")
		os.Exit(2)
	}
	summary, err := analyzeRun(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "poc15-analyze: %v\n", err)
		os.Exit(1)
	}
	if err := validateSummary(summary, cleanRegressionCriteria()); err != nil {
		fmt.Fprintf(os.Stderr, "poc15-analyze: acceptance criteria failed: %v\n", err)
		os.Exit(1)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(summary); err != nil {
		fmt.Fprintf(os.Stderr, "poc15-analyze: encode summary: %v\n", err)
		os.Exit(1)
	}
}

// analyzeRun summarizes one POC15 run directory containing agent JSONL files and
// an optional monitor-report.json.
// Intent: The analyzer treats logs as event records to count, not as authority over
// agents; it does not mutate run state. Source: DI-timah
func analyzeRun(runDir string) (RunSummary, error) {
	logDir, resolveErr := resolveRunLogDir(runDir)
	if resolveErr != nil {
		return RunSummary{}, resolveErr
	}
	summary := RunSummary{
		RunDir:                         logDir,
		EventCounts:                    make(map[string]int),
		OutcomeCounts:                  make(map[string]int),
		AgentCounts:                    make(map[string]int),
		FailureCounts:                  make(map[string]int),
		ProtocolCounts:                 make(map[string]int),
		TrustDrivenChoiceCounts:        make(map[string]int),
		EconomicsCounts:                make(map[string]int),
		VerificationCounts:             make(map[string]int),
		ReplicaRecoveryCounts:          make(map[string]int),
		RPCDriftCounts:                 make(map[string]int),
		ShippingCounts:                 make(map[string]int),
		RelationshipTransitionCounts:   make(map[string]int),
		DynamicTopologyCounts:          make(map[string]int),
		LocalResourceCounts:            make(map[string]int),
		ResourceTrustCouplingCounts:    make(map[string]int),
		DurabilityCounts:               make(map[string]int),
		RetentionCounts:                make(map[string]int),
		PressureCounts:                 make(map[string]int),
		ReplayCounts:                   make(map[string]int),
		TrustCautionCounts:             make(map[string]int),
		RuntimeAdapterEventCounts:      make(map[string]int),
		DecentralizedMonitorCounts:     make(map[string]int),
		RouteCounts:                    make(map[string]int),
		ArrayPayloadProtocolCounts:     make(map[string]int),
		AgentCASCounts:                 make(map[string]int),
		MessageArtifactDirectionCounts: make(map[string]int),
		MessageArtifactProtocolCounts:  make(map[string]int),
		MessageDAGParentLocationCounts: make(map[string]int),
		MessageShapeSpecimenCounts:     make(map[string]int),
		MigrationCounts:                make(map[string]int),
		RestartCounts:                  make(map[string]int),
		ForbiddenVocabularyCounts:      make(map[string]int),
	}
	logPaths := jsonlLogPaths(logDir)
	sort.Strings(logPaths)
	for _, logPath := range logPaths {
		if err := summarizeLog(logPath, &summary); err != nil {
			return RunSummary{}, err
		}
	}
	report, reportErr := readMonitorReport(filepath.Join(logDir, "monitor-report.json"))
	if reportErr != nil {
		return RunSummary{}, reportErr
	}
	summary.MonitorReport = report
	if artifactErr := summarizeMessageArtifacts(logDir, &summary); artifactErr != nil {
		return RunSummary{}, artifactErr
	}
	countMonitorVocabulary(report, &summary)
	summary.MissingRequiredEventNames = missingRequiredEvents(summary)
	summary.ScoreReport = computeScores(summary)
	summary.ProductionFitness = computeProductionFitness(summary)
	return summary, nil
}

// resolveRunLogDir accepts either the concrete JSONL log directory or the
// parent run directory that contains the `run/` event-record subdirectory.
// Intent: Operators should not be able to accidentally analyze the parent path
// and receive a misleading zero-event success report. Source: DI-jidah
func resolveRunLogDir(inputPath string) (string, error) {
	if _, statErr := os.Stat(inputPath); statErr != nil {
		return "", statErr
	}
	if len(jsonlLogPaths(inputPath)) > 0 {
		return inputPath, nil
	}
	nestedRunPath := filepath.Join(inputPath, "run")
	if _, statErr := os.Stat(nestedRunPath); statErr == nil && len(jsonlLogPaths(nestedRunPath)) > 0 {
		return nestedRunPath, nil
	}
	return "", fmt.Errorf("%s contains no JSONL logs; pass the log directory or its parent run directory", inputPath)
}

func jsonlLogPaths(runDir string) []string {
	logPaths, globErr := filepath.Glob(filepath.Join(runDir, "*.jsonl"))
	if globErr != nil {
		return nil
	}
	filteredPaths := make([]string, 0, len(logPaths))
	for _, logPath := range logPaths {
		if filepath.Base(logPath) == "message-dag.jsonl" {
			// Intent: The message DAG index is raw-message review metadata, not
			// a decision.Event stream. It is validated separately by
			// summarizeMessageArtifacts. Source: DI-tuhop
			continue
		}
		filteredPaths = append(filteredPaths, logPath)
	}
	return filteredPaths
}

func cleanRegressionCriteria() AcceptanceCriteria {
	return AcceptanceCriteria{
		MinTotalEvents:                    1,
		RequireMonitorReport:              true,
		MinMonitorScore:                   4,
		RequireFulfillmentWorkflow:        true,
		RequireAccountingUpdated:          true,
		RequireAccountingDuplicate:        true,
		RequireNotPromisedSuppression:     true,
		RequireNoResourceTrustCoupling:    true,
		RequireAccountingDuplicateConfirm: true,
		RequireCASComputeSuperset:         true,
		RequireRuntimeAdapterSuperset:     true,
		RequireDecentralizedMonitoring:    true,
		RequireAllKnownArrayPayloads:      true,
		RequireMixedVersionMigration:      true,
		RequireRunInternalRestart:         true,
		RequireEventVocabulary:            true,
		RequireRawMessageArtifacts:        true,
		RequireMessageShapeSpecimens:      true,
		RequirePersistentSessions:         true,
		MinScoreOverall:                   4,
	}
}

// validateSummary checks the POC15 clean-run event-record gates without mutating
// logs or relationship state.
// Intent: The acceptance contract should catch regression events directly in
// analyzer output instead of requiring the operator to interpret counts by eye.
// Source: DI-jidah
func validateSummary(summary RunSummary, criteria AcceptanceCriteria) error {
	var failures []string
	if summary.TotalEvents < criteria.MinTotalEvents {
		failures = append(failures, fmt.Sprintf("total_events=%d below minimum %d", summary.TotalEvents, criteria.MinTotalEvents))
	}
	if criteria.RequireFulfillmentWorkflow && summary.EventCounts["fulfillment_workflow_completed"] != 1 {
		failures = append(failures, fmt.Sprintf("fulfillment_workflow_completed=%d want 1", summary.EventCounts["fulfillment_workflow_completed"]))
	}
	if criteria.RequireAccountingUpdated && summary.ShippingCounts["accounting_updated"] != 1 {
		failures = append(failures, fmt.Sprintf("accounting_updated=%d want 1", summary.ShippingCounts["accounting_updated"]))
	}
	if criteria.RequireAccountingDuplicate && summary.ShippingCounts["accounting_update_duplicate"] != 1 {
		failures = append(failures, fmt.Sprintf("accounting_update_duplicate=%d want 1", summary.ShippingCounts["accounting_update_duplicate"]))
	}
	if criteria.RequireAccountingDuplicateConfirm && summary.ShippingCounts["accounting_update_duplicate_confirmed"] != 1 {
		failures = append(failures, fmt.Sprintf("accounting_update_duplicate_confirmed=%d want 1", summary.ShippingCounts["accounting_update_duplicate_confirmed"]))
	}
	if criteria.RequireNotPromisedSuppression && summary.EventCounts["promise_not_promised_suppressed"] == 0 {
		failures = append(failures, "promise_not_promised_suppressed=0 want >0")
	}
	if criteria.RequireNoResourceTrustCoupling && len(summary.ResourceTrustCouplingCounts) != 0 {
		failures = append(failures, fmt.Sprintf("resource_trust_coupling_counts=%v want empty", summary.ResourceTrustCouplingCounts))
	}
	if criteria.RequireCASComputeSuperset {
		if len(summary.MissingRequiredEventNames) > 0 {
			failures = append(failures, "missing required events: "+strings.Join(summary.MissingRequiredEventNames, ", "))
		}
		if len(summary.RPCDriftCounts) > 0 {
			failures = append(failures, fmt.Sprintf("rpc_drift_counts=%v want empty", summary.RPCDriftCounts))
		}
		if summary.ProtocolCounts[pcid.CASStorageV1] == 0 {
			failures = append(failures, "cas_storage_v1 event missing")
		}
		if summary.ProtocolCounts[pcid.CIDComputeV1] == 0 {
			failures = append(failures, "cid_compute_v1 event missing")
		}
		// Intent: POC15 should prove that the vague generic-report pCID was
		// replaced by a narrow identity/key protocol without losing exact pCID
		// routing coverage. Source: DI-vipih
		if summary.ProtocolCounts[pcid.IdentityKeyV1] == 0 {
			failures = append(failures, "identity_key_v1 event missing")
		}
		if summary.ScoreReport.Overall < criteria.MinScoreOverall {
			failures = append(failures, fmt.Sprintf("score_report.overall=%d below minimum %d", summary.ScoreReport.Overall, criteria.MinScoreOverall))
		}
	}
	if criteria.RequireRuntimeAdapterSuperset && summary.ScoreReport.RuntimeAdapter < 5 {
		failures = append(failures, fmt.Sprintf("score_report.runtime_adapter=%d want 5", summary.ScoreReport.RuntimeAdapter))
	}
	if criteria.RequireDecentralizedMonitoring && summary.ScoreReport.Monitoring < 5 {
		failures = append(failures, fmt.Sprintf("score_report.monitoring=%d want 5", summary.ScoreReport.Monitoring))
	}
	if criteria.RequireAllKnownArrayPayloads {
		for _, protocolName := range requiredArrayPayloadProtocols() {
			if summary.ArrayPayloadProtocolCounts[protocolName] == 0 {
				failures = append(failures, "array payload missing for "+protocolName)
			}
		}
	}
	if criteria.RequireMixedVersionMigration && summary.ScoreReport.Migration < 5 {
		failures = append(failures, fmt.Sprintf("score_report.migration=%d want 5", summary.ScoreReport.Migration))
	}
	if criteria.RequireRunInternalRestart && summary.ScoreReport.Restart < 5 {
		failures = append(failures, fmt.Sprintf("score_report.restart=%d want 5", summary.ScoreReport.Restart))
	}
	if criteria.RequireEventVocabulary && len(summary.ForbiddenVocabularyCounts) > 0 {
		failures = append(failures, fmt.Sprintf("forbidden vocabulary counts: %v", summary.ForbiddenVocabularyCounts))
	}
	if criteria.RequireRawMessageArtifacts {
		failures = append(failures, rawMessageArtifactFailures(summary)...)
	}
	if criteria.RequireMessageShapeSpecimens {
		failures = append(failures, messageShapeSpecimenFailures(summary)...)
	}
	if criteria.RequirePersistentSessions {
		failures = append(failures, persistentSessionFailures(summary)...)
	}
	if criteria.RequireMonitorReport {
		if summary.MonitorReport == nil {
			failures = append(failures, "monitor_report missing")
		} else {
			failures = append(failures, monitorScoreFailures(*summary.MonitorReport, criteria.MinMonitorScore)...)
		}
	}
	if len(failures) > 0 {
		return errors.New(strings.Join(failures, "; "))
	}
	return nil
}

func monitorScoreFailures(report decision.MonitorReport, minimum int) []string {
	scores := map[string]int{
		"promise_theory_fit":      report.PromiseTheoryFit,
		"autonomy":                report.Autonomy,
		"protocol_validity":       report.ProtocolValidity,
		"local_trust_correctness": report.LocalTrustCorrectness,
		"imposition_avoidance":    report.ImpositionAvoidance,
	}
	names := make([]string, 0, len(scores))
	for name := range scores {
		names = append(names, name)
	}
	sort.Strings(names)
	var failures []string
	for _, name := range names {
		if scores[name] < minimum {
			failures = append(failures, fmt.Sprintf("%s=%d below minimum %d", name, scores[name], minimum))
		}
	}
	return failures
}

// summarizeMessageArtifacts verifies the observer-collected raw message CAS and
// message DAG index for one run.
// Intent: POC15 should fail loudly if it only records summaries about messages.
// The analyzer reads the collector-owned index, validates each binary `.cbor`
// artifact by exact hash, and counts directions/protocols without mutating run
// state or treating the artifacts as authority over any agent. Source: DI-tuhop
func summarizeMessageArtifacts(runDir string, summary *RunSummary) error {
	indexPath := filepath.Join(runDir, "message-dag.jsonl")
	indexFile, openErr := os.Open(indexPath)
	if openErr != nil {
		if errors.Is(openErr, os.ErrNotExist) {
			return nil
		}
		return openErr
	}
	defer func() {
		closeErr := indexFile.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "poc15-analyze: close %s: %v\n", indexPath, closeErr)
		}
	}()
	casObjects := make(map[string]bool)
	recordsByHash := make(map[string]messageDAGRecord)
	scanner := bufio.NewScanner(indexFile)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		var record messageDAGRecord
		if unmarshalErr := json.Unmarshal(scanner.Bytes(), &record); unmarshalErr != nil {
			return fmt.Errorf("%s:%d: %w", indexPath, lineNumber, unmarshalErr)
		}
		artifactPath, pathErr := artifactPathWithinRun(runDir, record.Path)
		if pathErr != nil {
			return fmt.Errorf("%s:%d: %w", indexPath, lineNumber, pathErr)
		}
		artifactBytes, readErr := os.ReadFile(artifactPath)
		if readErr != nil {
			return fmt.Errorf("%s:%d: read %s: %w", indexPath, lineNumber, artifactPath, readErr)
		}
		actualHash := protocol.HashExactBytes(artifactBytes)
		if actualHash != record.ExactSHA256 {
			return fmt.Errorf("%s:%d: artifact hash mismatch record=%s actual=%s", indexPath, lineNumber, record.ExactSHA256, actualHash)
		}
		if bytes.Contains(artifactBytes, obsoletePayloadPrefixBytes()) {
			summary.MessageArtifactBadPrefixCount++
		}
		summary.MessageDAGRecordCount++
		if record.ParentExactSHA256 != "" {
			summary.MessageDAGParentLinkCount++
			summary.MessageDAGParentLocationCounts[firstNonEmpty(record.ParentLinkLocation, "unspecified")]++
		} else if isAckArtifactDirection(record.Direction) {
			summary.AckMessageMissingParentCount++
		}
		summary.MessageArtifactDirectionCounts[record.Direction]++
		summary.MessageArtifactProtocolCounts[record.Protocol]++
		casObjects[record.ExactSHA256] = true
		recordsByHash[record.ExactSHA256] = record
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return scanErr
	}
	summarizeMessageDAGTraversal(recordsByHash, summary)
	summary.MessageArtifactCount = summary.MessageDAGRecordCount
	summary.MessageCASObjectCount = len(casObjects)
	return nil
}

// summarizeMessageDAGTraversal walks parent links inside the retained
// raw-message index after collapsing duplicate exact-message artifacts into
// unique DAG nodes.
// Intent: POC15 should prove the retained CAS/index can be traversed as a DAG,
// not only counted as independent files. Missing parents or cycles are POC
// failures because they make later message-history review unreliable. Source:
// DI-kohuj
func summarizeMessageDAGTraversal(recordsByHash map[string]messageDAGRecord, summary *RunSummary) {
	summary.MessageDAGNodeCount = len(recordsByHash)
	depthMemo := make(map[string]int, len(recordsByHash))
	visiting := make(map[string]bool, len(recordsByHash))
	for exactHash, record := range recordsByHash {
		if record.ParentExactSHA256 == "" {
			summary.MessageDAGRootCount++
		}
		depth, ok := messageDAGDepth(exactHash, recordsByHash, depthMemo, visiting)
		if ok {
			summary.MessageDAGReachableCount++
			if depth > summary.MessageDAGMaxDepth {
				summary.MessageDAGMaxDepth = depth
			}
		} else {
			summary.MessageDAGMissingParentCount++
		}
	}
}

func messageDAGDepth(exactHash string, recordsByHash map[string]messageDAGRecord, depthMemo map[string]int, visiting map[string]bool) (int, bool) {
	if depth, ok := depthMemo[exactHash]; ok {
		return depth, true
	}
	record, ok := recordsByHash[exactHash]
	if !ok {
		return 0, false
	}
	if record.ParentExactSHA256 == "" {
		depthMemo[exactHash] = 1
		return 1, true
	}
	if visiting[exactHash] {
		return 0, false
	}
	visiting[exactHash] = true
	parentDepth, parentOK := messageDAGDepth(record.ParentExactSHA256, recordsByHash, depthMemo, visiting)
	delete(visiting, exactHash)
	if !parentOK {
		return 0, false
	}
	depth := parentDepth + 1
	depthMemo[exactHash] = depth
	return depth, true
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func artifactPathWithinRun(runDir, relativePath string) (string, error) {
	if relativePath == "" {
		return "", fmt.Errorf("message artifact path is empty")
	}
	if filepath.IsAbs(relativePath) {
		return "", fmt.Errorf("message artifact path %q must be relative", relativePath)
	}
	cleanPath := filepath.Clean(filepath.FromSlash(relativePath))
	if cleanPath == "." || cleanPath == ".." || strings.HasPrefix(cleanPath, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("message artifact path %q escapes run directory", relativePath)
	}
	return filepath.Join(runDir, cleanPath), nil
}

func rawMessageArtifactFailures(summary RunSummary) []string {
	var failures []string
	if summary.EventCounts["raw_message_artifact_emitted"] == 0 {
		failures = append(failures, "raw_message_artifact_emitted=0 want >0")
	}
	if summary.MessageArtifactCount == 0 {
		failures = append(failures, "message_artifact_count=0 want >0")
	}
	if summary.MessageCASObjectCount == 0 {
		failures = append(failures, "message_cas_object_count=0 want >0")
	}
	if summary.MessageArtifactBadPrefixCount != 0 {
		failures = append(failures, fmt.Sprintf("message_artifact_bad_prefix_count=%d want 0", summary.MessageArtifactBadPrefixCount))
	}
	if summary.MessageDAGMissingParentCount != 0 {
		failures = append(failures, fmt.Sprintf("message_dag_missing_parent_count=%d want 0", summary.MessageDAGMissingParentCount))
	}
	if summary.MessageDAGReachableCount != summary.MessageDAGNodeCount {
		failures = append(failures, fmt.Sprintf("message_dag_reachable_count=%d want node_count=%d", summary.MessageDAGReachableCount, summary.MessageDAGNodeCount))
	}
	if summary.MessageDAGMaxDepth < 2 {
		failures = append(failures, fmt.Sprintf("message_dag_max_depth=%d want >=2", summary.MessageDAGMaxDepth))
	}
	if summary.AckMessageMissingParentCount != 0 {
		failures = append(failures, fmt.Sprintf("ack_message_missing_parent_count=%d want 0", summary.AckMessageMissingParentCount))
	}
	for _, direction := range []string{"sent", "received", "ack_sent", "ack_received"} {
		if summary.MessageArtifactDirectionCounts[direction] == 0 {
			failures = append(failures, "message artifact direction missing: "+direction)
		}
	}
	for _, protocolName := range []string{pcid.KernelReceiveV1, pcid.CASStorageV1, pcid.CIDComputeV1, pcid.RouteV1, pcid.RelationshipV1} {
		if summary.MessageArtifactProtocolCounts[protocolName] == 0 {
			failures = append(failures, "message artifact protocol missing: "+protocolName)
		}
	}
	return failures
}

func isAckArtifactDirection(direction string) bool {
	// Intent: The raw-message index may use flow-specific ACK direction names,
	// but every retained ACK-like artifact should link back to its request in the
	// message DAG. Source: DI-vopab
	return direction == "ack_sent" || direction == "ack_received" || strings.HasSuffix(direction, "_ack")
}

func persistentSessionFailures(summary RunSummary) []string {
	// Intent: POC15 should fail loudly if transport silently regresses to
	// one-shot TCP or if ACK artifacts stop parent-linking the exact request
	// message CIDs used by persistent-session demux. Source: DI-vopab
	var failures []string
	for _, eventName := range []string{"persistent_session_opened", "persistent_session_reused", "persistent_session_closed"} {
		if summary.EventCounts[eventName] == 0 {
			failures = append(failures, eventName+"=0 want >0")
		}
	}
	if summary.AckMessageMissingParentCount != 0 {
		failures = append(failures, fmt.Sprintf("ack_message_missing_parent_count=%d want 0", summary.AckMessageMissingParentCount))
	}
	return failures
}

func messageShapeSpecimenFailures(summary RunSummary) []string {
	var failures []string
	for _, eventName := range []string{
		"message_shape_transport_specimen_emitted",
		"message_shape_native_proof_specimen_emitted",
		"message_shape_envelope_parent_specimen_emitted",
		"message_shape_payload_parent_specimen_emitted",
		"message_shape_cose_payload_specimen_emitted",
		"message_shape_cose_proof_specimen_emitted",
		"message_shape_cose_payload_verified",
		"message_shape_cose_proof_verified",
		"message_shape_cose_tamper_rejected",
		"kernel_role_profile_recorded",
	} {
		if summary.EventCounts[eventName] == 0 {
			failures = append(failures, eventName+"=0 want >0")
		}
	}
	if summary.MessageArtifactDirectionCounts["shape_specimen"] == 0 {
		failures = append(failures, "message artifact direction missing: shape_specimen")
	}
	if summary.MessageDAGParentLinkCount < 2 {
		failures = append(failures, fmt.Sprintf("message_dag_parent_link_count=%d want >=2", summary.MessageDAGParentLinkCount))
	}
	for _, location := range []string{"envelope", "payload"} {
		if summary.MessageDAGParentLocationCounts[location] == 0 {
			failures = append(failures, "message DAG parent location missing: "+location)
		}
	}
	for _, protocolName := range messageShapeSpecimenProtocols() {
		if summary.MessageArtifactProtocolCounts[protocolName] == 0 {
			failures = append(failures, "message shape artifact protocol missing: "+protocolName)
		}
	}
	return failures
}

// countMonitorVocabulary adds monitor-report vocabulary findings to the same
// run summary gate used for JSONL events.
// Intent: DI-kirat applies to POC development-tool output as well as agent
// events, because otherwise the retired term can re-enter guide prose through
// analyzer summaries. Source: DI-kirat
func countMonitorVocabulary(report *decision.MonitorReport, summary *RunSummary) {
	if report == nil {
		return
	}
	count := countForbiddenVocabulary(report.Summary)
	for _, concern := range report.Concerns {
		count += countForbiddenVocabulary(concern)
	}
	if count > 0 {
		summary.ForbiddenVocabularyCounts["monitor_report"] += count
	}
}

func countForbiddenVocabulary(values ...string) int {
	count := 0
	for _, value := range values {
		for _, term := range forbiddenVocabularyTerms() {
			count += strings.Count(strings.ToLower(value), term)
		}
	}
	return count
}

func knownProtocolNamesForAnalysis() []string {
	protocolNames := []string{pcid.KernelReceiveV1, pcid.CASStorageV1, pcid.CIDComputeV1, pcid.IdentityKeyV1, pcid.RouteV1, pcid.RelationshipV1, pcid.AccountingV1, pcid.UPSLabelV1, pcid.PostalScaleV1, pcid.PrinterPortV1}
	protocolNames = append(protocolNames, messageShapeSpecimenProtocols()...)
	return protocolNames
}

func requiredArrayPayloadProtocols() []string {
	// Intent: POC15's fresh clean-run traffic should exercise pCID-owned arrays
	// across relationship, kernel receive, device, CAS, compute, and identity
	// protocols so generic map payloads cannot remain the implicit target shape.
	// Source: DI-dirat; DI-pusak
	return []string{pcid.KernelReceiveV1, pcid.CASStorageV1, pcid.CIDComputeV1, pcid.IdentityKeyV1, pcid.RouteV1, pcid.RelationshipV1, pcid.AccountingV1, pcid.UPSLabelV1, pcid.PostalScaleV1, pcid.PrinterPortV1}
}

func obsoletePayloadPrefixBytes() []byte {
	// Intent: Keep the forbidden token out of source text while still checking
	// retained raw CBOR bytes for the old naming mistake. Source: DI-pusak
	return []byte("field" + "_")
}

func messageShapeSpecimenProtocols() []string {
	// Intent: Specimen pCIDs are analyzed as raw message artifacts, not as app
	// receive-promise protocols. They pressure outer arity, parent-link, and COSE
	// parsing without expanding normal app traffic. Source: DI-mosat
	return []string{
		pcid.MessageShapeTransportV1,
		pcid.MessageShapeNativeProofV1,
		pcid.MessageShapeEnvelopeParentsV1,
		pcid.MessageShapePayloadParentsV1,
		pcid.MessageShapeCOSEPayloadV1,
		pcid.MessageShapeCOSEProofV1,
	}
}

func isArrayPayloadEvent(eventName string) bool {
	switch eventName {
	case "pcid_owned_array_payload_sent", "pcid_owned_array_payload_received", "pcid_owned_array_ack_sent", "pcid_owned_array_ack_received":
		return true
	default:
		return false
	}
}

func forbiddenVocabularyTerms() []string {
	// Intent: Keep active POC15 run output from drifting back to retired
	// production-looking vocabulary while avoiding literal reintroduction of those
	// words in the codebase sweep itself. Source: DI-jofus
	return []string{"evi" + "dence", "boun" + "dary"}
}

func summarizeLog(logPath string, summary *RunSummary) error {
	logFile, openErr := os.Open(logPath)
	if openErr != nil {
		return openErr
	}
	defer func() {
		closeErr := logFile.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "poc15-analyze: close %s: %v\n", logPath, closeErr)
		}
	}()
	scanner := bufio.NewScanner(logFile)
	var previousEvent decision.Event
	hasPreviousEvent := false
	for scanner.Scan() {
		var event decision.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return fmt.Errorf("%s: %w", logPath, err)
		}
		if count := countForbiddenVocabulary(event.Observer, event.Event, event.Outcome, event.Peer, event.Detail); count > 0 {
			summary.ForbiddenVocabularyCounts["run_events"] += count
		}
		summary.TotalEvents++
		summary.EventCounts[event.Event]++
		summary.OutcomeCounts[event.Outcome]++
		summary.AgentCounts[event.Observer]++
		if event.Outcome != "kept" {
			summary.FailureCounts[event.Event]++
		}
		if strings.HasPrefix(event.Event, "message_shape_") {
			summary.MessageShapeSpecimenCounts[event.Event]++
		}
		for _, protocolName := range knownProtocolNamesForAnalysis() {
			if strings.Contains(event.Detail, "pcid="+protocolName) || strings.Contains(event.Detail, "protocol="+protocolName) {
				summary.ProtocolCounts[protocolName]++
				if isArrayPayloadEvent(event.Event) {
					summary.ArrayPayloadProtocolCounts[protocolName]++
				}
			}
		}
		if event.Event == "trust_driven_peer_choice" {
			summary.TrustDrivenChoiceCounts[event.Observer]++
		}
		if strings.HasPrefix(event.Event, "economics_") {
			summary.EconomicsCounts[event.Observer]++
		}
		if strings.Contains(event.Event, "verification") || strings.Contains(event.Event, "verified") || event.Event == "compute_result_locally_verified" || event.Event == "compute_result_peer_verified" {
			summary.VerificationCounts[event.Observer]++
		}
		if strings.HasPrefix(event.Event, "replica_recovery") || event.Event == "cas_replica_serve_promised" {
			summary.ReplicaRecoveryCounts[event.Observer]++
		}
		if isRPCDrift(event) {
			summary.RPCDriftCounts[event.Observer]++
		}
		if isShippingEvent(event.Event) {
			summary.ShippingCounts[event.Event]++
		}
		if isRelationshipTransition(event.Event) {
			summary.RelationshipTransitionCounts[event.Event]++
		}
		if isDynamicTopologyEvent(event.Event) {
			summary.DynamicTopologyCounts[event.Event]++
		}
		if isLocalResourceEvent(event) {
			summary.LocalResourceCounts[event.Event]++
		}
		if isDurabilityEvent(event.Event) {
			summary.DurabilityCounts[event.Event]++
		}
		if isRetentionEvent(event.Event) {
			summary.RetentionCounts[event.Event]++
		}
		if isPressureEvent(event.Event) {
			summary.PressureCounts[event.Event]++
		}
		if isReplayEvent(event.Event) {
			summary.ReplayCounts[event.Event]++
		}
		if isTrustCautionEvent(event.Event) {
			summary.TrustCautionCounts[event.Event]++
		}
		if isRuntimeAdapterEvent(event.Event) {
			summary.RuntimeAdapterEventCounts[event.Event]++
		}
		if isDecentralizedMonitorEvent(event.Event) {
			summary.DecentralizedMonitorCounts[event.Event]++
		}
		if isRouteEvent(event.Event) {
			summary.RouteCounts[event.Event]++
		}
		if isAgentCASEvent(event.Event) {
			summary.AgentCASCounts[event.Event]++
		}
		if isMigrationEvent(event.Event) {
			summary.MigrationCounts[event.Event]++
		}
		if isRestartEvent(event.Event) {
			summary.RestartCounts[event.Event]++
		}
		if hasPreviousEvent && isResourceTrustCoupling(previousEvent, event) {
			summary.ResourceTrustCouplingCounts[event.Observer]++
		}
		previousEvent = event
		hasPreviousEvent = true
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func missingRequiredEvents(summary RunSummary) []string {
	var missing []string
	for _, eventName := range requiredRegressionEvents() {
		if summary.EventCounts[eventName] == 0 {
			missing = append(missing, eventName)
		}
	}
	return missing
}

func requiredRegressionEvents() []string {
	return []string{
		"runtime_readiness_promised",
		"peer_readiness_observed",
		"runtime_done_promised",
		"tcp_message_sent",
		"tcp_message_received",
		"tcp_message_send_failed",
		"persistent_session_opened",
		"persistent_session_reused",
		"persistent_session_closed",
		"network_outage_variant_selected",
		"cas_storage_promised",
		"cas_retention_promised",
		"cas_replication_promised",
		"cas_bytes_stored",
		"cas_bytes_retrieved",
		"cas_replica_stored",
		"cas_replica_serve_promised",
		"cas_multi_object_pressure",
		"primary_storage_unavailable",
		"replica_capability_token_issued",
		"replica_capability_token_received",
		"replica_capability_token_redeemed",
		"capability_token_issued",
		"capability_token_received",
		"capability_token_ttl_promised",
		"capability_token_ttl_observed",
		"capability_token_expired",
		"capability_token_revoked",
		"capability_token_renewal_requested",
		"capability_token_renewed",
		"cas_corrupt_bytes_rejected",
		"cas_corrupt_report_recorded",
		"compute_context_promised",
		"compute_function_executed",
		"compute_alternate_function_executed",
		"compute_followup_function_requested",
		"cid_compute_promised",
		"compute_result_promised",
		"compute_result_received",
		"compute_bad_result_promised",
		"compute_cache_checkpointed",
		"compute_cache_miss",
		"compute_cache_miss_observed",
		"compute_cache_hit",
		"compute_cache_reused",
		"compute_result_locally_verified",
		"compute_result_locally_rejected",
		"compute_result_peer_verified",
		"compute_verifier_disagreement",
		"compute_disagreement_resolved_locally",
		"compute_verification_received",
		// Intent: Compute verification event records now stay under cid_compute_v1;
		// they must not depend on the removed generic report pCID. Source:
		// DI-vipih
		"compute_verification_report_received",
		"trust_driven_peer_choice",
		"persisted_trust_history_loaded",
		"dynamic_peer_choice_from_persisted_trust",
		"economics_credit_accepted",
		"economics_price_refused",
		"economics_capacity_reserved",
		"economics_capacity_refused",
		"economics_credits_spent",
		"economics_credits_earned",
		"economics_payment_withheld",
		"replica_recovery_requested",
		"replica_recovery_succeeded",
		"trust_updated",
		"trust_repair_promise_recorded",
		// Intent: POC15 must keep DI-fijov recovery caution visible in analyzer
		// gates, not only in monitor prose, so malformed/broken events cannot
		// be followed by silent instant trust recovery. Source: DI-sihuz
		"trust_caution_recorded",
		"trust_caution_consumed",
		"trust_recovery_delayed",
		"trust_repair_future_only",
		"dynamic_tcp_topology_probe_started",
		"dynamic_tcp_topology_send_blocked",
		"dynamic_tcp_topology_send_succeeded",
		"unknown_pcid_not_promised",
		"promise_variant_not_promised",
		"bad_proof_sent",
		"bad_proof_rejected",
		"key_rotation_promise_recorded",
		"promise_envelope_validated",
		// Intent: CAS and compute have begun moving away from universal generic
		// maps; the regression must see pCID-owned array payloads in both
		// directions before the run is considered clean. Source: DI-gahuh;
		// DI-pusak
		"pcid_owned_array_payload_sent",
		"pcid_owned_array_payload_received",
		"pcid_owned_array_ack_sent",
		"pcid_owned_array_ack_received",
		// Intent: POC15's regression gate now covers run-scoped stores,
		// promise-shaped retention/GC, pressure promises, and replay protection
		// so those operational concerns cannot regress behind raw event counts.
		// Source: DI-sunuf
		"run_scoped_store_empty",
		"run_scoped_store_saved",
		"cas_run_store_saved",
		"event_run_store_saved",
		"compute_cache_run_store_saved",
		"retention_promise_recorded",
		"retention_until_promised",
		"delete_after_promised",
		"token_expiry_gc_promised",
		"disk_pressure_gc_promised",
		"superseded_checkpoint_gc_promised",
		"gc_object_retained",
		"gc_promise_ended",
		"gc_object_removed",
		"app_kernel_backpressure_promised",
		"app_kernel_rate_limit_promised",
		"backpressure_capacity_promised",
		"backpressure_capacity_refused",
		"send_rate_promised",
		"accept_rate_promised",
		"rate_limit_self_promise_recorded",
		"replay_window_promised",
		"replay_envelope_recorded",
		"capability_token_replay_rejected",
		"capability_token_replay_observed",
		"replay_probe_rejected",
		// Intent: POC15 must remain a POC13 superset while adding heterogeneous
		// WASM and stdio process-runtime adapter events. Source: DI-linof;
		// DI-kimim
		"wasm_process_agent_started",
		"wasm_module_instantiated",
		"wasm_export_called",
		"wasm_export_result_observed",
		"wasm_adapter_promise_sent",
		"wasm_adapter_ack_received",
		// Intent: Runtime-adapter usefulness now means Alice can ask Peggy to
		// keep a real cid_compute_v1 promise through WASM execution. Source:
		// DI-sivis
		"wasm_compute_request_promised",
		"wasm_compute_request_received",
		"wasm_compute_function_executed",
		"wasm_compute_result_promised",
		"wasm_compute_result_verified",
		// Intent: Peggy and Victor must perform useful routed promise work, not
		// only record that WASM/stdio runtime adapters exist. Source:
		// DI-pamob; DI-kimim
		"wasm_useful_work_promised",
		"wasm_useful_work_ack_received",
		"stdio_worker_started",
		"stdio_cbor_request_sent",
		"stdio_worker_envelope_received",
		"stdio_cbor_envelope_received",
		"stdio_adapter_kernel_forwarded",
		"stdio_cbor_ack_sent",
		"stdio_worker_ack_event",
		"stdio_cbor_ack_event",
		// Intent: Victor's stdio adapter must prove useful compute work by
		// delegating exact inbound cid_compute_v1 bytes to the worker and
		// returning the worker-signed ACK unchanged. Source: DI-sivis
		"stdio_compute_request_promised",
		"stdio_compute_worker_started",
		"stdio_compute_request_forwarded",
		"stdio_compute_worker_executed",
		"stdio_compute_ack_received",
		"stdio_compute_result_verified",
		"stdio_useful_work_promised",
		"stdio_useful_work_ack_received",
		// Intent: POC15 monitoring experiments must be decentralized because
		// production agents do not share a global observer. Source: DI-lulof
		"decentralized_monitoring_model_recorded",
		"local_event_summary_promised",
		"peer_carried_attestation_promised",
		"bearer_token_exchange_rate_observed",
		"relationship_topology_signal_observed",
		"voluntary_gossip_promised",
		// Intent: POC15 must test pCID evolution and same-run restart recovery
		// explicitly rather than assuming POC13 durability is enough. Source:
		// DI-linof
		"mixed_version_pcid_migration_promised",
		"mixed_version_legacy_pcid_observed",
		"mixed_version_successor_pcid_selected",
		"run_internal_restart_orchestration_promised",
		"run_internal_restart_checkpoint_promised",
		"run_internal_restart_recovery_observed",
		// Intent: POC15 should exercise permanent local distrust and
		// untrusted-transit exclusion as local promises and event records, not as global
		// bans or route enforcement. Source: DI-kinaf
		"permanent_distrust_decided",
		"permanent_distrust_future_repair_not_promised",
		"permanent_distrust_direct_peer_removed",
		"permanent_distrust_send_blocked",
		"transit_exclusion_promised",
		"input_transit_exclusion_recorded",
		"output_transit_exclusion_recorded",
		"transit_candidate_rejected",
		"transit_route_candidate_blocked",
		"transit_safe_route_selected",
		// Intent: POC15 now has an executable route_v1 slice where Alice,
		// Bob, Carol, and Dave build a route from voluntary neighboring
		// promises rather than a kernel-owned route table. Source: DI-lihir
		"route_setup_promise_made",
		"route_forward_promise_made",
		"route_forward_promise_kept",
		"route_reachability_promised",
		"route_reachability_confirmed",
		"route_payment_promised",
		"route_exclusion_promise_made",
		"route_exclusion_used_in_choice",
		"route_carried_message_sent",
		"route_carried_message_received",
		"route_carried_message_delivered",
		// Intent: DI-kohuj moves POC15 route pressure beyond one-shot route
		// specimens: normal route_v1 traffic must now exercise parent slots,
		// payload parent links, explicit route lifetime, asymmetric response
		// terms, reciprocal credit, and routed runtime compute. Source: DI-kohuj
		"route_multiarity_parent_slot_received",
		"route_payload_parent_link_received",
		"route_durability_promised",
		"route_reused_message_sent",
		"route_reused_message_delivered",
		"route_asymmetric_response_path_promised",
		"route_asymmetric_response_path_handled",
		"route_credit_offered",
		"route_credit_earned",
		"route_credit_spent",
		"route_carried_envelope_validated",
		"route_runtime_compute_message_sent",
		"route_runtime_compute_message_delivered",
		"wasm_routed_compute_result_verified",
		"stdio_routed_compute_result_verified",
		"transport_proof_comparison_recorded",
		// Intent: POC15 now requires agent-accessible sparse CAS, local message
		// DAG indexing, peer storage/retrieval promises, bearer storage tokens,
		// encrypted-object CIDs, and local GC records. Source: DI-manul
		"agent_cas_access_promised",
		"agent_cas_store_incomplete",
		"agent_cas_object_stored",
		"agent_cas_message_stored",
		"agent_cas_internal_object_stored",
		"agent_cas_encrypted_object_stored",
		"agent_cas_ciphertext_cid_selected",
		"agent_cas_cleartext_cid_not_used",
		"agent_cas_peer_storage_promised",
		"agent_cas_peer_retrieval_promised",
		"agent_cas_peer_object_stored",
		"agent_cas_sparse_object_missing",
		"agent_cas_retrieval_not_promised",
		"agent_cas_bearer_storage_token_issued",
		"agent_cas_bearer_storage_token_received",
		"agent_cas_bearer_storage_token_transferred",
		"agent_cas_bearer_storage_token_redeemed",
		"agent_cas_gc_object_retained",
		"agent_cas_gc_object_removed",
		"message_dag_node_indexed",
		"message_dag_missing_parent_recorded",
	}
}

func computeScores(summary RunSummary) ScoreReport {
	scores := ScoreReport{}
	// Intent: Transport fitness now includes persistent-session reuse and ACK
	// parent-link coverage so a run cannot pass by reopening one TCP connection
	// per message. Source: DI-vopab
	addScore(&scores.Transport, summary.EventCounts["tcp_message_sent"] > 0 && summary.EventCounts["tcp_message_received"] > 0 && summary.EventCounts["tcp_message_send_failed"] > 0 && summary.EventCounts["bad_proof_rejected"] > 0 && summary.EventCounts["persistent_session_opened"] > 0 && summary.EventCounts["persistent_session_reused"] > 0 && summary.AckMessageMissingParentCount == 0)
	addScore(&scores.Storage, summary.EventCounts["cas_bytes_stored"] > 0 && summary.EventCounts["cas_bytes_retrieved"] > 0 && summary.EventCounts["cas_multi_object_pressure"] > 0 && summary.EventCounts["capability_token_renewed"] > 0)
	addScore(&scores.Compute, summary.EventCounts["compute_function_executed"] > 0 && summary.EventCounts["compute_alternate_function_executed"] > 0 && summary.EventCounts["compute_result_received"] > 0 && summary.EventCounts["compute_bad_result_promised"] > 0 && summary.EventCounts["compute_cache_reused"] > 0)
	addScore(&scores.Economics, summary.EventCounts["economics_price_refused"] > 0 && summary.EventCounts["economics_capacity_refused"] > 0 && summary.EventCounts["economics_credits_spent"] > 0 && summary.EventCounts["economics_credits_earned"] > 0)
	addScore(&scores.Trust, summary.EventCounts["trust_updated"] > 0 && summary.EventCounts["trust_driven_peer_choice"] > 0 && summary.EventCounts["persisted_trust_history_loaded"] > 0 && summary.EventCounts["dynamic_peer_choice_from_persisted_trust"] > 0 && summary.EventCounts["trust_caution_recorded"] > 0 && summary.EventCounts["trust_recovery_delayed"] > 0 && summary.EventCounts["dynamic_tcp_topology_send_blocked"] > 0 && summary.EventCounts["dynamic_tcp_topology_send_succeeded"] > 0)
	addScore(&scores.Verification, summary.EventCounts["compute_result_locally_verified"] > 0 && summary.EventCounts["compute_result_locally_rejected"] > 0 && summary.EventCounts["compute_result_peer_verified"] > 0 && summary.EventCounts["compute_verifier_disagreement"] > 0 && summary.EventCounts["compute_disagreement_resolved_locally"] > 0 && summary.EventCounts["compute_verification_report_received"] > 0)
	addScore(&scores.Replica, summary.EventCounts["replica_recovery_succeeded"] > 0 && summary.EventCounts["replica_capability_token_redeemed"] > 0)
	// Intent: Analyzer scoring should treat durability, retention, pressure, and
	// replay as first-class POC15 fitness dimensions instead of burying them under
	// generic protocol validity. Source: DI-sunuf
	addScore(&scores.Durability, summary.EventCounts["run_scoped_store_saved"] > 0 && summary.EventCounts["cas_run_store_saved"] > 0 && summary.EventCounts["event_run_store_saved"] > 0 && summary.EventCounts["compute_cache_run_store_saved"] > 0)
	addScore(&scores.Retention, summary.EventCounts["retention_promise_recorded"] > 0 && summary.EventCounts["gc_object_retained"] > 0 && summary.EventCounts["gc_object_removed"] > 0)
	addScore(&scores.Pressure, summary.EventCounts["backpressure_capacity_promised"] > 0 && summary.EventCounts["backpressure_capacity_refused"] > 0 && summary.EventCounts["send_rate_promised"] > 0 && summary.EventCounts["accept_rate_promised"] > 0)
	addScore(&scores.Replay, summary.EventCounts["replay_window_promised"] > 0 && summary.EventCounts["replay_envelope_recorded"] > 0 && summary.EventCounts["capability_token_replay_rejected"] > 0 && summary.EventCounts["replay_probe_rejected"] > 0)
	// Intent: POC15 adds runtime-adapter and decentralized-monitoring dimensions
	// without relaxing inherited POC13 storage/compute/trust gates. Source:
	// DI-linof; DI-lulof
	addScore(&scores.RuntimeAdapter, summary.EventCounts["wasm_module_instantiated"] > 0 && summary.EventCounts["wasm_export_result_observed"] > 0 && summary.EventCounts["wasm_adapter_ack_received"] > 0 && summary.EventCounts["wasm_useful_work_promised"] > 0 && summary.EventCounts["wasm_compute_result_verified"] > 0 && summary.EventCounts["wasm_routed_compute_result_verified"] > 0 && summary.EventCounts["stdio_cbor_envelope_received"] > 0 && summary.EventCounts["stdio_cbor_ack_event"] > 0 && summary.EventCounts["stdio_adapter_kernel_forwarded"] > 0 && summary.EventCounts["stdio_useful_work_promised"] > 0 && summary.EventCounts["stdio_compute_result_verified"] > 0 && summary.EventCounts["stdio_routed_compute_result_verified"] > 0)
	addScore(&scores.Monitoring, summary.EventCounts["decentralized_monitoring_model_recorded"] > 0 && summary.EventCounts["bearer_token_exchange_rate_observed"] > 0 && summary.EventCounts["voluntary_gossip_promised"] > 0)
	addScore(&scores.Route, summary.EventCounts["route_setup_promise_made"] > 0 && summary.EventCounts["route_forward_promise_made"] >= 2 && summary.EventCounts["route_forward_promise_kept"] >= 2 && summary.EventCounts["route_reachability_confirmed"] > 0 && summary.EventCounts["route_carried_message_delivered"] > 0 && summary.EventCounts["route_multiarity_parent_slot_received"] > 0 && summary.EventCounts["route_payload_parent_link_received"] > 0 && summary.EventCounts["route_reused_message_delivered"] > 0 && summary.EventCounts["route_asymmetric_response_path_handled"] > 0 && summary.EventCounts["route_credit_earned"] > 0 && summary.ProtocolCounts[pcid.RouteV1] > 0)
	// Intent: Agent CAS fitness is distinct from the collector-owned raw message
	// CAS: the run must show sparse local stores, DAG indexing, peer storage,
	// bearer-token incentives, encrypted-object CIDs, and local GC. Source:
	// DI-manul
	addScore(&scores.AgentCAS, summary.EventCounts["agent_cas_access_promised"] > 0 && summary.EventCounts["agent_cas_store_incomplete"] > 0 && summary.EventCounts["agent_cas_message_stored"] > 0 && summary.EventCounts["message_dag_node_indexed"] > 0 && summary.EventCounts["agent_cas_peer_storage_promised"] > 0 && summary.EventCounts["agent_cas_bearer_storage_token_redeemed"] > 0 && summary.EventCounts["agent_cas_encrypted_object_stored"] > 0 && summary.EventCounts["agent_cas_gc_object_removed"] > 0)
	addScore(&scores.Migration, summary.EventCounts["mixed_version_pcid_migration_promised"] > 0 && summary.EventCounts["mixed_version_successor_pcid_selected"] > 0)
	addScore(&scores.Restart, summary.EventCounts["run_internal_restart_orchestration_promised"] > 0 && summary.EventCounts["run_internal_restart_recovery_observed"] > 0)
	scores.Overall = (scores.Transport + scores.Storage + scores.Compute + scores.Economics + scores.Trust + scores.Verification + scores.Replica + scores.Durability + scores.Retention + scores.Pressure + scores.Replay + scores.RuntimeAdapter + scores.Monitoring + scores.Route + scores.AgentCAS + scores.Migration + scores.Restart) / 17
	if len(summary.MissingRequiredEventNames) > 0 {
		scores.Concerns = append(scores.Concerns, "missing required events: "+strings.Join(summary.MissingRequiredEventNames, ", "))
	}
	if len(summary.RPCDriftCounts) > 0 {
		scores.Concerns = append(scores.Concerns, "RPC drift detected")
	}
	return scores
}

func computeProductionFitness(summary RunSummary) ProductionFitnessReport {
	report := ProductionFitnessReport{
		Baseline: "POC15 regression baseline for POC15; not production software",
		Verdict:  "fit for continued POC work, not production deployment",
	}
	if summary.ScoreReport.Overall < 5 {
		report.BlockingGaps = append(report.BlockingGaps, fmt.Sprintf("score_report.overall=%d below 5", summary.ScoreReport.Overall))
	}
	if summary.MonitorReport == nil {
		report.BlockingGaps = append(report.BlockingGaps, "monitor_report missing")
	} else {
		for _, failure := range monitorScoreFailures(*summary.MonitorReport, 5) {
			report.BlockingGaps = append(report.BlockingGaps, failure)
		}
	}
	if len(summary.MissingRequiredEventNames) > 0 {
		report.BlockingGaps = append(report.BlockingGaps, "missing required regression events")
	}
	if len(summary.RPCDriftCounts) > 0 {
		report.BlockingGaps = append(report.BlockingGaps, "RPC drift terms detected")
	}
	report.BlockingGaps = append(report.BlockingGaps, rawMessageArtifactFailures(summary)...)
	report.ReadyForProduction = len(report.BlockingGaps) == 0
	if report.ReadyForProduction {
		report.Verdict = "production-candidate event complete for current POC scope"
	}
	return report
}

func addScore(score *int, kept bool) {
	if kept {
		*score = 5
		return
	}
	*score = 1
}

func isRPCDrift(event decision.Event) bool {
	text := strings.ToLower(event.Event + " " + event.Detail)
	for _, badTerm := range []string{"rpc", "command", "permission", "authorization", "conformance", "enforce"} {
		if strings.Contains(text, badTerm) {
			return true
		}
	}
	return false
}

func isShippingEvent(eventName string) bool {
	switch eventName {
	case "package_weighed", "package_weight_received", "shipping_address_promised", "shipping_address_received", "shipping_label_printed", "shipping_label_received", "accounting_updated", "accounting_update_confirmed", "accounting_update_duplicate", "accounting_update_duplicate_confirmed", "printer_capability_issued", "printer_capability_received", "printer_port_printed", "printer_port_print_confirmed":
		return true
	default:
		return false
	}
}

func isRelationshipTransition(eventName string) bool {
	switch eventName {
	case "direct_peer_added", "direct_peer_removed", "direct_peer_unchanged":
		return true
	default:
		return false
	}
}

func isDynamicTopologyEvent(eventName string) bool {
	return strings.HasPrefix(eventName, "dynamic_tcp_topology_")
}

func isDurabilityEvent(eventName string) bool {
	return strings.Contains(eventName, "run_store") || strings.HasPrefix(eventName, "run_scoped_store_")
}

func isRetentionEvent(eventName string) bool {
	return strings.Contains(eventName, "retention") || strings.HasPrefix(eventName, "gc_") || strings.Contains(eventName, "expiry_gc") || strings.Contains(eventName, "delete_after") || strings.Contains(eventName, "disk_pressure")
}

func isPressureEvent(eventName string) bool {
	return strings.Contains(eventName, "backpressure") || strings.Contains(eventName, "rate_limit") || eventName == "send_rate_promised" || eventName == "accept_rate_promised"
}

func isReplayEvent(eventName string) bool {
	return strings.Contains(eventName, "replay")
}

func isTrustCautionEvent(eventName string) bool {
	return strings.HasPrefix(eventName, "trust_caution_") || strings.HasPrefix(eventName, "permanent_distrust_") || eventName == "trust_recovery_delayed" || eventName == "trust_repair_future_only"
}

func isRuntimeAdapterEvent(eventName string) bool {
	return strings.HasPrefix(eventName, "wasm_") || strings.HasPrefix(eventName, "stdio_")
}

func isDecentralizedMonitorEvent(eventName string) bool {
	switch eventName {
	case "decentralized_monitoring_model_recorded",
		"local_event_summary_promised",
		"peer_carried_attestation_promised",
		"bearer_token_exchange_rate_observed",
		"relationship_topology_signal_observed",
		"voluntary_gossip_promised",
		"transit_exclusion_promised",
		"input_transit_exclusion_recorded",
		"output_transit_exclusion_recorded",
		"transit_candidate_rejected",
		"transit_route_candidate_blocked",
		"transit_safe_route_selected":
		return true
	default:
		return false
	}
}

func isRouteEvent(eventName string) bool {
	return strings.HasPrefix(eventName, "route_")
}

func isAgentCASEvent(eventName string) bool {
	return strings.HasPrefix(eventName, "agent_cas_") || strings.HasPrefix(eventName, "message_dag_")
}

func isMigrationEvent(eventName string) bool {
	return strings.HasPrefix(eventName, "mixed_version_")
}

func isRestartEvent(eventName string) bool {
	return strings.HasPrefix(eventName, "run_internal_restart_")
}

// isLocalResourceEvent recognizes event records about the observing app's own
// budget/capacity state.
// Intent: Analyzer output should make it visible if local resource exhaustion is
// accidentally coupled back into peer trust transitions. Source: DI-vujob
func isLocalResourceEvent(event decision.Event) bool {
	if event.Event == "local_resource_exhausted" {
		return true
	}
	if event.Event == "promise_withheld" && (strings.Contains(event.Detail, "budget exhausted") || strings.Contains(event.Detail, "capacity exhausted")) {
		return true
	}
	if event.Event == "resource_promise_broken" && strings.Contains(event.Detail, "local") {
		return true
	}
	return false
}

// isResourceTrustCoupling reports suspicious adjacency between local resource
// exhaustion and peer relationship transitions for one observer.
// Intent: This is a regression tripwire; local budget/capacity state should not
// immediately mutate peer trust without intervening promise event records. Source:
// DI-vujob
func isResourceTrustCoupling(previousEvent, event decision.Event) bool {
	if previousEvent.Observer != event.Observer {
		return false
	}
	return (isRelationshipTransition(previousEvent.Event) && isLocalResourceEvent(event)) ||
		(isLocalResourceEvent(previousEvent) && isRelationshipTransition(event.Event))
}

func readMonitorReport(reportPath string) (*decision.MonitorReport, error) {
	reportBytes, readErr := os.ReadFile(reportPath)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return nil, nil
		}
		return nil, readErr
	}
	var report decision.MonitorReport
	if err := json.Unmarshal(reportBytes, &report); err != nil {
		return nil, err
	}
	return &report, nil
}
