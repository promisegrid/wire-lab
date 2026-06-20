package main

import (
	"os"
	"path/filepath"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/decision"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/pcid"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
)

func TestAnalyzeRunSummarizesEventsAndMonitorReport(t *testing.T) {
	runDir := t.TempDir()
	writeFile(t, filepath.Join(runDir, "alice.jsonl"), ""+
		`{"observer":"alice","event":"promise_sent","outcome":"kept","peer":"bob","detail":"sent"}`+"\n"+
		`{"observer":"alice","event":"send_failed","outcome":"broken","peer":"bob","detail":"connection refused"}`+"\n"+
		`{"observer":"alice","event":"shipping_label_received","outcome":"kept","peer":"ups_label_printer","detail":"tracking"}`+"\n"+
		`{"observer":"alice","event":"printer_port_print_confirmed","outcome":"kept","peer":"printer_port","detail":"spool"}`+"\n"+
		`{"observer":"alice","event":"local_resource_exhausted","outcome":"non_commitment","peer":"bob","detail":"capacity exhausted"}`+"\n"+
		`{"observer":"alice","event":"direct_peer_unchanged","outcome":"kept","peer":"bob","detail":"outcome=non_commitment trust=0"}`+"\n"+
		`{"observer":"alice","event":"wasm_module_instantiated","outcome":"kept","peer":"victor","detail":"runtime=wazero"}`+"\n"+
		`{"observer":"alice","event":"wasm_adapter_ack_received","outcome":"kept","peer":"victor","detail":"pcid=relationship_v1"}`+"\n"+
		`{"observer":"alice","event":"stdio_cbor_ack_event","outcome":"kept","peer":"victor","detail":"exact_sha256=test"}`+"\n"+
		`{"observer":"alice","event":"bearer_token_exchange_rate_observed","outcome":"kept","peer":"grace","detail":"local market signal"}`+"\n"+
		`{"observer":"alice","event":"mixed_version_successor_pcid_selected","outcome":"kept","peer":"bob","detail":"current pCID"}`+"\n"+
		`{"observer":"alice","event":"run_internal_restart_recovery_observed","outcome":"kept","peer":"victor","detail":"same-run recovery"}`+"\n")
	writeFile(t, filepath.Join(runDir, "bob.jsonl"), ""+
		`{"observer":"bob","event":"decision_rejected","outcome":"malformed","peer":"alice","detail":"bad target"}`+"\n"+
		`{"observer":"bob","event":"direct_peer_added","outcome":"kept","peer":"alice","detail":"trust=2"}`+"\n")
	writeFile(t, filepath.Join(runDir, "monitor-report.json"), `{
  "promise_theory_fit": 4,
  "autonomy": 5,
  "protocol_validity": 3,
  "local_trust_correctness": 4,
  "imposition_avoidance": 5,
  "summary": "test",
  "concerns": ["synthetic"]
}`)

	summary, err := analyzeRun(runDir)
	if err != nil {
		t.Fatalf("analyze run: %v", err)
	}
	if summary.TotalEvents != 14 {
		t.Fatalf("total events = %d, want 14", summary.TotalEvents)
	}
	if summary.EventCounts["promise_sent"] != 1 || summary.FailureCounts["send_failed"] != 1 || summary.FailureCounts["decision_rejected"] != 1 {
		t.Fatalf("unexpected counts: %#v failures %#v", summary.EventCounts, summary.FailureCounts)
	}
	if summary.MonitorReport == nil || summary.MonitorReport.ProtocolValidity != 3 {
		t.Fatalf("monitor report not summarized: %#v", summary.MonitorReport)
	}
	if summary.ShippingCounts["shipping_label_received"] != 1 {
		t.Fatalf("shipping counts not summarized: %#v", summary.ShippingCounts)
	}
	if summary.ShippingCounts["printer_port_print_confirmed"] != 1 {
		t.Fatalf("printer-port counts not summarized: %#v", summary.ShippingCounts)
	}
	if summary.RelationshipTransitionCounts["direct_peer_added"] != 1 {
		t.Fatalf("relationship transitions not summarized: %#v", summary.RelationshipTransitionCounts)
	}
	if summary.LocalResourceCounts["local_resource_exhausted"] != 1 {
		t.Fatalf("local resource counts not summarized: %#v", summary.LocalResourceCounts)
	}
	if summary.ResourceTrustCouplingCounts["alice"] != 1 {
		t.Fatalf("resource/trust coupling not summarized: %#v", summary.ResourceTrustCouplingCounts)
	}
	if summary.RuntimeAdapterEventCounts["wasm_adapter_ack_received"] != 1 {
		t.Fatalf("runtime adapter event counts not summarized: %#v", summary.RuntimeAdapterEventCounts)
	}
	if summary.DecentralizedMonitorCounts["bearer_token_exchange_rate_observed"] != 1 {
		t.Fatalf("decentralized monitor counts not summarized: %#v", summary.DecentralizedMonitorCounts)
	}
	if summary.MigrationCounts["mixed_version_successor_pcid_selected"] != 1 {
		t.Fatalf("migration counts not summarized: %#v", summary.MigrationCounts)
	}
	if summary.RestartCounts["run_internal_restart_recovery_observed"] != 1 {
		t.Fatalf("restart counts not summarized: %#v", summary.RestartCounts)
	}
}

func TestAnalyzeRunAcceptsParentRunDirectory(t *testing.T) {
	parentDir := t.TempDir()
	logDir := filepath.Join(parentDir, "run")
	if err := os.Mkdir(logDir, 0o755); err != nil {
		t.Fatalf("make log dir: %v", err)
	}
	writeFile(t, filepath.Join(logDir, "alice.jsonl"), `{"observer":"alice","event":"promise_sent","outcome":"kept","peer":"bob","detail":"sent"}`+"\n")

	summary, err := analyzeRun(parentDir)
	if err != nil {
		t.Fatalf("analyze parent run dir: %v", err)
	}
	if summary.RunDir != logDir {
		t.Fatalf("summary run dir = %q, want %q", summary.RunDir, logDir)
	}
	if summary.TotalEvents != 1 {
		t.Fatalf("total events = %d, want 1", summary.TotalEvents)
	}
}

func TestAnalyzeRunRejectsDirectoryWithoutLogs(t *testing.T) {
	_, err := analyzeRun(t.TempDir())
	if err == nil {
		t.Fatalf("analyze without logs should fail")
	}
}

func TestValidateSummaryAcceptsCleanRegressionEvent(t *testing.T) {
	summary := cleanRegressionSummary()
	if err := validateSummary(summary, cleanRegressionCriteria()); err != nil {
		t.Fatalf("clean regression summary should pass: %v", err)
	}
}

func TestValidateSummaryRejectsMissingSuppression(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["promise_not_promised_suppressed"] = 0
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing suppression should fail")
	}
}

func TestValidateSummaryRejectsResourceTrustCoupling(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.ResourceTrustCouplingCounts["alice"] = 1
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("resource/trust coupling should fail")
	}
}

func TestValidateSummaryRejectsMissingAdapterEvent(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["stdio_cbor_ack_event"] = 0
	summary.ScoreReport = computeScores(summary)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing stdio adapter event should fail")
	}
}

func TestValidateSummaryRejectsForbiddenVocabulary(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.ForbiddenVocabularyCounts["run_events"] = 1
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("forbidden vocabulary should fail")
	}
}

func TestAnalyzeRunCountsForbiddenVocabulary(t *testing.T) {
	retiredWord := "evi" + "dence"
	retiredInterfaceWord := "boun" + "dary"
	runDir := t.TempDir()
	writeFile(t, filepath.Join(runDir, "alice.jsonl"), `{"observer":"alice","event":"promise_sent","outcome":"kept","peer":"bob","detail":"fresh `+retiredWord+` and `+retiredInterfaceWord+`"}`+"\n")
	writeFile(t, filepath.Join(runDir, "monitor-report.json"), `{"promise_theory_fit":5,"autonomy":5,"protocol_validity":5,"local_trust_correctness":5,"imposition_avoidance":5,"summary":"clean","concerns":["no issues"]}`)
	summary, err := analyzeRun(runDir)
	if err != nil {
		t.Fatalf("analyze run: %v", err)
	}
	if summary.ForbiddenVocabularyCounts["run_events"] < 2 {
		t.Fatalf("forbidden vocabulary not counted: %#v", summary.ForbiddenVocabularyCounts)
	}
}

func TestValidateSummaryRejectsMissingDecentralizedMonitoring(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["bearer_token_exchange_rate_observed"] = 0
	summary.ScoreReport = computeScores(summary)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing decentralized monitoring event should fail")
	}
}

func TestValidateSummaryRejectsMissingMigrationEvent(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["mixed_version_successor_pcid_selected"] = 0
	summary.ScoreReport = computeScores(summary)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing migration event should fail")
	}
}

func TestValidateSummaryRejectsMissingRestartEvent(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["run_internal_restart_recovery_observed"] = 0
	summary.ScoreReport = computeScores(summary)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing restart event should fail")
	}
}

func TestValidateSummaryRejectsMissingPermanentDistrustEvent(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["permanent_distrust_decided"] = 0
	summary.MissingRequiredEventNames = missingRequiredEvents(summary)
	summary.ScoreReport = computeScores(summary)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing permanent distrust event should fail")
	}
}

func TestValidateSummaryRejectsMissingTransitExclusionEvent(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["transit_safe_route_selected"] = 0
	summary.MissingRequiredEventNames = missingRequiredEvents(summary)
	summary.ScoreReport = computeScores(summary)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing transit exclusion event should fail")
	}
}

func TestValidateSummaryRejectsMissingMigratedArrayPayloadEvent(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["pcid_owned_array_payload_received"] = 0
	summary.MissingRequiredEventNames = missingRequiredEvents(summary)
	summary.ScoreReport = computeScores(summary)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing migrated array payload event should fail")
	}
}

func TestValidateSummaryRejectsMissingArrayPayloadProtocol(t *testing.T) {
	summary := cleanRegressionSummary()
	delete(summary.ArrayPayloadProtocolCounts, pcid.RelationshipV1)
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing relationship array payload coverage should fail")
	}
}

func TestAnalyzeRunSummarizesRawMessageArtifacts(t *testing.T) {
	runDir := t.TempDir()
	rawEnvelope := []byte{0xd9, 0x67, 0x72, 0x69, 0x64, 0x83, 0x01, 0x02, 0x03}
	exactHash := protocol.HashExactBytes(rawEnvelope)
	writeFile(t, filepath.Join(runDir, "alice.jsonl"), `{"observer":"alice","event":"raw_message_artifact_emitted","outcome":"kept","peer":"bob","detail":"direction=sent pcid=route_v1 exact_sha256=`+exactHash+`"}`+"\n")
	writeFile(t, filepath.Join(runDir, "monitor-report.json"), `{"promise_theory_fit":5,"autonomy":5,"protocol_validity":5,"local_trust_correctness":5,"imposition_avoidance":5,"summary":"clean","concerns":[]}`)
	messageCASDir := filepath.Join(runDir, "message-cas")
	if err := os.Mkdir(messageCASDir, 0o755); err != nil {
		t.Fatalf("make message CAS dir: %v", err)
	}
	writeBytes(t, filepath.Join(messageCASDir, exactHash+".cbor"), rawEnvelope)
	writeFile(t, filepath.Join(runDir, "message-dag.jsonl"), `{"source":"alice/agent:alice/stdout","observer":"alice","direction":"sent","peer":"bob","protocol":"route_v1","exact_sha256":"`+exactHash+`","path":"message-cas/`+exactHash+`.cbor"}`+"\n")

	summary, err := analyzeRun(runDir)
	if err != nil {
		t.Fatalf("analyze run: %v", err)
	}
	if summary.MessageArtifactCount != 1 || summary.MessageCASObjectCount != 1 || summary.MessageDAGRecordCount != 1 || summary.MessageDAGNodeCount != 1 {
		t.Fatalf("message artifact counts not summarized: %#v", summary)
	}
	if summary.MessageArtifactDirectionCounts["sent"] != 1 {
		t.Fatalf("message direction counts not summarized: %#v", summary.MessageArtifactDirectionCounts)
	}
	if summary.MessageArtifactProtocolCounts[pcid.RouteV1] != 1 {
		t.Fatalf("message protocol counts not summarized: %#v", summary.MessageArtifactProtocolCounts)
	}
}

func TestAnalyzeRunCountsBadPayloadPrefixInRawMessageArtifacts(t *testing.T) {
	// Intent: POC16 keeps raw CBOR messages for later review, so the analyzer
	// must inspect the retained bytes and reject any return to the retired
	// prefixed payload vocabulary. Source: DI-pusak
	runDir := t.TempDir()
	rawEnvelope := []byte{0xd9, 0x67, 0x72, 0x69, 0x64, 0x83, 0x01, 0x02, 0x03}
	rawEnvelope = append(rawEnvelope, obsoletePayloadPrefixBytes()...)
	exactHash := protocol.HashExactBytes(rawEnvelope)
	writeFile(t, filepath.Join(runDir, "alice.jsonl"), `{"observer":"alice","event":"raw_message_artifact_emitted","outcome":"kept","peer":"bob","detail":"direction=sent pcid=route_v1 exact_sha256=`+exactHash+`"}`+"\n")
	writeFile(t, filepath.Join(runDir, "monitor-report.json"), `{"promise_theory_fit":5,"autonomy":5,"protocol_validity":5,"local_trust_correctness":5,"imposition_avoidance":5,"summary":"clean","concerns":[]}`)
	messageCASDir := filepath.Join(runDir, "message-cas")
	if err := os.Mkdir(messageCASDir, 0o755); err != nil {
		t.Fatalf("make message CAS dir: %v", err)
	}
	artifactPath := filepath.Join(messageCASDir, exactHash+".cbor")
	writeBytes(t, artifactPath, rawEnvelope)
	writeFile(t, filepath.Join(runDir, "message-dag.jsonl"), `{"cid":"cid://message","path":"message-cas/`+exactHash+`.cbor","exact_sha256":"`+exactHash+`","direction":"sent","protocol":"`+pcid.RouteV1+`","parents":[]}`+"\n")

	summary, err := analyzeRun(runDir)
	if err != nil {
		t.Fatalf("analyze run: %v", err)
	}
	if summary.MessageArtifactBadPrefixCount != 1 {
		t.Fatalf("bad prefix not counted: %#v", summary)
	}
	if failures := rawMessageArtifactFailures(summary); len(failures) == 0 {
		t.Fatalf("bad prefix should fail raw message artifact checks")
	}
}

func TestRawMessageArtifactFailuresCompareReachabilityToUniqueNodes(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.MessageDAGRecordCount = 4
	summary.MessageDAGNodeCount = 2
	summary.MessageDAGReachableCount = 2

	if failures := rawMessageArtifactFailures(summary); len(failures) != 0 {
		t.Fatalf("duplicate artifact records should not fail unique-node DAG reachability: %v", failures)
	}
}

func TestValidateSummaryRejectsMissingRawMessageArtifacts(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.MessageArtifactCount = 0
	summary.MessageCASObjectCount = 0
	summary.MessageDAGRecordCount = 0
	summary.MessageArtifactDirectionCounts = map[string]int{}
	summary.MessageArtifactProtocolCounts = map[string]int{}
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing raw message artifacts should fail")
	}
}

func TestValidateSummaryRejectsMissingMessageShapeSpecimens(t *testing.T) {
	summary := cleanRegressionSummary()
	summary.EventCounts["message_shape_cose_proof_verified"] = 0
	err := validateSummary(summary, cleanRegressionCriteria())
	if err == nil {
		t.Fatalf("missing message shape specimen coverage should fail")
	}
}

func cleanRegressionSummary() RunSummary {
	eventCounts := map[string]int{
		"fulfillment_workflow_completed":                 1,
		"promise_not_promised_suppressed":                1,
		"raw_message_artifact_emitted":                   1,
		"message_shape_transport_specimen_emitted":       1,
		"message_shape_native_proof_specimen_emitted":    1,
		"message_shape_envelope_parent_specimen_emitted": 1,
		"message_shape_payload_parent_specimen_emitted":  1,
		"message_shape_cose_payload_specimen_emitted":    1,
		"message_shape_cose_proof_specimen_emitted":      1,
		"message_shape_cose_payload_verified":            1,
		"message_shape_cose_proof_verified":              1,
		"message_shape_cose_tamper_rejected":             1,
		"kernel_role_profile_recorded":                   1,
	}
	for _, eventName := range requiredRegressionEvents() {
		eventCounts[eventName] = 1
	}
	eventCounts["poc15_superset_named_agent_preserved"] = 19
	eventCounts["poc16_protocol_spec_doc_recorded"] = 10
	eventCounts["llm_spec_context_embedded"] = 10
	eventCounts["llm_spec_context_cid_recorded"] = 10
	summary := RunSummary{
		RunDir:        "test/run",
		TotalEvents:   len(eventCounts),
		EventCounts:   eventCounts,
		OutcomeCounts: map[string]int{},
		AgentCounts:   map[string]int{},
		FailureCounts: map[string]int{},
		ProtocolCounts: map[string]int{
			pcid.KernelTransportV1:    1,
			pcid.CASStorageV1:         1,
			pcid.CIDComputeV1:         1,
			pcid.IdentityKeyV1:        1,
			pcid.RouteV1:              1,
			pcid.RelationshipV1:       1,
			pcid.ProductionShippingV1: 1,
			pcid.SecureCapabilityV1:   1,
			pcid.EncryptedPayloadV1:   1,
			pcid.ParserBuilderRoleV1:  1,
			pcid.MapPayloadProfileV1:  1,
		},
		ArrayPayloadProtocolCounts: map[string]int{},
		MessageArtifactCount:       4,
		MessageCASObjectCount:      4,
		MessageDAGRecordCount:      4,
		MessageDAGNodeCount:        4,
		MessageDAGParentLinkCount:  2,
		MessageDAGRootCount:        2,
		MessageDAGReachableCount:   4,
		MessageDAGMaxDepth:         2,
		MessageArtifactDirectionCounts: map[string]int{
			"sent":           1,
			"received":       1,
			"ack_sent":       1,
			"ack_received":   1,
			"shape_specimen": 1,
		},
		MessageArtifactProtocolCounts: map[string]int{
			pcid.KernelTransportV1:             1,
			pcid.CASStorageV1:                  1,
			pcid.CIDComputeV1:                  1,
			pcid.RouteV1:                       1,
			pcid.RelationshipV1:                1,
			pcid.ProductionShippingV1:          1,
			pcid.MessageShapeTransportV1:       1,
			pcid.MessageShapeNativeProofV1:     1,
			pcid.MessageShapeEnvelopeParentsV1: 1,
			pcid.MessageShapePayloadParentsV1:  1,
			pcid.MessageShapeCOSEPayloadV1:     1,
			pcid.MessageShapeCOSEProofV1:       1,
			pcid.SecureCapabilityV1:            1,
			pcid.EncryptedPayloadV1:            1,
			pcid.ParserBuilderRoleV1:           1,
			pcid.MapPayloadProfileV1:           1,
		},
		MessageDAGParentLocationCounts: map[string]int{
			"envelope": 1,
			"payload":  1,
		},
		MessageShapeSpecimenCounts: map[string]int{},
		ShippingCounts: map[string]int{
			"accounting_updated":                    1,
			"accounting_update_duplicate":           1,
			"accounting_update_duplicate_confirmed": 1,
		},
		RelationshipTransitionCounts: map[string]int{},
		LocalResourceCounts:          map[string]int{},
		ResourceTrustCouplingCounts:  map[string]int{},
		ForbiddenVocabularyCounts:    map[string]int{},
		CASRetrievalCounts:           map[string]int{},
		TokenSecurityCounts:          map[string]int{},
		PersistentSessionCounts:      map[string]int{},
		PersistentSessionOpenCounts: map[string]int{
			"test-session": 1,
		},
		PersistentSessionTerminalCounts: map[string]int{
			"test-session": 1,
		},
		PersistentSessionTerminalReasonCounts: map[string]int{
			"process_shutdown": 1,
		},
		MonitorReport: &decision.MonitorReport{
			PromiseTheoryFit:      5,
			Autonomy:              5,
			ProtocolValidity:      4,
			LocalTrustCorrectness: 5,
			ImpositionAvoidance:   5,
		},
	}
	for _, agentName := range []string{"alice", "bob", "carol", "dave", "ellen", "frank", "grace", "heidi", "ivan", "judy", "mallory", "oscar", "fulfillment", "postal_scale", "ups_label_printer", "printer_port", "accounting", "peggy", "victor"} {
		summary.AgentCounts[agentName] = 1
	}
	for _, protocolName := range requiredArrayPayloadProtocols() {
		summary.ArrayPayloadProtocolCounts[protocolName] = 1
	}
	summary.ScoreReport = computeScores(summary)
	return summary
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func writeBytes(t *testing.T, path string, content []byte) {
	t.Helper()
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
