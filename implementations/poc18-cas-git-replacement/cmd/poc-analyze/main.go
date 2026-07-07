// Command poc-analyze checks deterministic POC18 fixture results.
//
// Intent: Analyzer checks remain non-production review aids; they must not imply
// a global monitor or authority. Source: DI-jifuj
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/bridge"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/carbundle"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/economy"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	pgrepo "promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/repo"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/retention"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/scenario"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	pocsync "promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/sync"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/workspace"
)

// fixtureResult mirrors poc-sim output so analyzer checks stay tied to the same
// deterministic fixture rather than becoming a production monitor.
//
// Intent: Check nahop.11 as local fixture evidence: separate peer CAS roots,
// complete CID-verified retrieval, and Bob's checkout from Bob's own CAS.
// Source: DI-gozov
type fixtureResult struct {
	workspace.IngestResult
	AliceStoreRoot     string                       `json:"alice_store_root"`
	BobStoreRoot       string                       `json:"bob_store_root"`
	CarolStoreRoot     string                       `json:"carol_store_root"`
	FrankStoreRoot     string                       `json:"frank_store_root"`
	MalloryStoreRoot   string                       `json:"mallory_store_root"`
	CarolSyncStatePath string                       `json:"carol_sync_state_path"`
	Scenario           scenario.Result              `json:"scenario"`
	Bridge             bridgeFixtureResult          `json:"bridge"`
	Retrieval          pocsync.RetrievalReport      `json:"retrieval"`
	ContinuousSync     pocsync.ContinuousSyncReport `json:"continuous_sync"`
	ScheduledSync      pocsync.SchedulerReport      `json:"scheduled_sync"`
	SyncAgentState     pocsync.AgentState           `json:"sync_agent_state"`
	Retention          retention.Report             `json:"retention"`
	RetentionPayment   economy.RedemptionReport     `json:"retention_payment"`
}

type bridgeFixtureResult struct {
	Export bridge.Result `json:"export"`
	Import bridge.Result `json:"import"`
	Push   bridge.Result `json:"push"`
	Pull   bridge.Result `json:"pull"`
}

// Report records local fixture checks.
type Report struct {
	RunRoot           string            `json:"run_root"`
	Pass              bool              `json:"pass"`
	Checks            map[string]string `json:"checks"`
	Objects           int               `json:"objects"`
	AliceObjects      int               `json:"alice_objects,omitempty"`
	BobObjects        int               `json:"bob_objects,omitempty"`
	CarolObjects      int               `json:"carol_objects,omitempty"`
	FrankObjects      int               `json:"frank_objects,omitempty"`
	MalloryObjects    int               `json:"mallory_objects,omitempty"`
	Retrieved         int               `json:"retrieved,omitempty"`
	Missing           int               `json:"missing,omitempty"`
	CollectorEvents   int               `json:"collector_events,omitempty"`
	CollectorMessages int               `json:"collector_messages,omitempty"`
	CollectorCARs     int               `json:"collector_cars,omitempty"`
}

type collectorEvent struct {
	Observer string `json:"observer"`
	Event    string `json:"event"`
	Outcome  string `json:"outcome"`
	Peer     string `json:"peer,omitempty"`
	Detail   string `json:"detail,omitempty"`
}

type collectorMessageRecord struct {
	Source      string `json:"source"`
	Observer    string `json:"observer"`
	Direction   string `json:"direction"`
	Peer        string `json:"peer"`
	Protocol    string `json:"protocol"`
	PromiseKind string `json:"promise_kind"`
	ExactCID    string `json:"exact_cid"`
	Path        string `json:"path"`
}

type collectorCARRecord struct {
	Source     string   `json:"source"`
	Observer   string   `json:"observer"`
	Direction  string   `json:"direction"`
	Peer       string   `json:"peer"`
	MessageCID string   `json:"message_cid"`
	CARCID     string   `json:"car_cid"`
	BlockCIDs  []string `json:"block_cids"`
	Path       string   `json:"path"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "poc-analyze: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	runRoot := flag.String("run-root", "/tmp/wire-lab-poc18-run", "approved POC18 runtime root")
	flag.Parse()
	report, analyzeErr := analyzeRunRoot(*runRoot)
	if analyzeErr != nil {
		return analyzeErr
	}
	reportPath := filepath.Join(*runRoot, "analysis.json")
	if err := writeJSON(reportPath, report); err != nil {
		return err
	}
	fmt.Printf("pass=%t objects=%d report=%s\n", report.Pass, report.Objects, reportPath)
	if !report.Pass {
		return fmt.Errorf("analysis failed")
	}
	return nil
}

func analyzeRunRoot(runRoot string) (Report, error) {
	resultPath := filepath.Join(runRoot, "result.json")
	if _, statErr := os.Stat(resultPath); statErr == nil {
		result, resultErr := readResult(resultPath)
		if resultErr != nil {
			return Report{}, resultErr
		}
		return analyze(runRoot, result)
	} else if !os.IsNotExist(statErr) {
		return Report{}, statErr
	}
	// Intent: The same analyzer command now validates the Docker/TCP remediation
	// path by reading only observer-collected artifacts. No agent can read this
	// collector output during the run. Source: DI-koriz
	return analyzeCollector(runRoot)
}

func readResult(path string) (fixtureResult, error) {
	var result fixtureResult
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		return fixtureResult{}, readErr
	}
	if err := json.Unmarshal(content, &result); err != nil {
		return fixtureResult{}, err
	}
	return result, nil
}

func analyze(runRoot string, result fixtureResult) (Report, error) {
	checks := map[string]string{}
	requiredCounts := []string{
		"chunk",
		"chunk_manifest",
		"posix_node:regular",
		"posix_node:directory",
		"posix_node:symlink",
		"posix_node:fifo",
		"hardlink_label",
		"reference_set:directory",
		"reference_set:branch",
		"reference_set:workspace",
		"reference_set:logical_change",
		"reference_set:review_thread",
		"reference_set:release",
		"snapshot",
	}
	pass := true
	for _, countName := range requiredCounts {
		if result.Counts[countName] > 0 {
			checks[countName] = "kept"
			continue
		}
		checks[countName] = "missing"
		pass = false
	}
	// Intent: The merged scenario intentionally moves README.md to docs/intro.md
	// and keeps a root copy label, so checkout validation follows labels rather
	// than treating the old filename as file identity. Source: DI-guban
	if _, err := os.Stat(filepath.Join(result.CheckoutRoot, "docs", "intro.md")); err == nil {
		checks["checkout_renamed_intro"] = "kept"
	} else {
		checks["checkout_renamed_intro"] = "missing"
		pass = false
	}
	if _, err := os.Stat(filepath.Join(result.CheckoutRoot, "README-copy.md")); err == nil {
		checks["checkout_root_copy"] = "kept"
	} else {
		checks["checkout_root_copy"] = "missing"
		pass = false
	}
	if _, err := os.Lstat(filepath.Join(result.CheckoutRoot, "README-link.md")); err == nil {
		checks["checkout_symlink"] = "kept"
	} else {
		checks["checkout_symlink"] = "missing"
		pass = false
	}
	if result.AliceStoreRoot != "" && result.BobStoreRoot != "" && result.AliceStoreRoot != result.BobStoreRoot {
		checks["separate_sparse_peer_cas"] = "kept"
	} else {
		checks["separate_sparse_peer_cas"] = "missing"
		pass = false
	}
	// Intent: The deterministic fixture must now prove that Alice's workspace is
	// a real repo-local grid workspace, not only a plain fixture directory passed
	// to lower-level scanner code. Source: DI-kiram
	if repoPass, repoErr := checkAliceGridRepo(result, checks); repoErr != nil {
		return Report{}, repoErr
	} else if !repoPass {
		pass = false
	}
	if len(result.Retrieval.Retrieved) > 0 {
		checks["peer_retrieval"] = "kept"
	} else {
		checks["peer_retrieval"] = "missing"
		pass = false
	}
	if len(result.Retrieval.Missing) == 0 {
		checks["peer_retrieval_complete"] = "kept"
	} else {
		checks["peer_retrieval_complete"] = "missing"
		pass = false
	}
	if result.Retrieval.InterestMessageCID != "" && result.Retrieval.AvailabilityMessageCID != "" {
		checks["sync_promise_messages"] = "kept"
	} else {
		checks["sync_promise_messages"] = "missing"
		pass = false
	}
	if result.Scenario.LineageNodeCID != "" && result.Scenario.LineageNodeCID == result.Scenario.RenameLabelNodeCID && result.Scenario.LineageNodeCID == result.Scenario.CopyLabelNodeCID {
		checks["rename_copy_preserves_node_lineage"] = "kept"
	} else {
		checks["rename_copy_preserves_node_lineage"] = "missing"
		pass = false
	}
	if result.Scenario.BobDivergentSnapshotCID != "" && result.Scenario.RenameCopySnapshotCID != "" && result.Scenario.BobDivergentSnapshotCID != result.Scenario.RenameCopySnapshotCID {
		checks["divergent_snapshots"] = "kept"
	} else {
		checks["divergent_snapshots"] = "missing"
		pass = false
	}
	if len(result.Scenario.MergeParentSnapshotCIDs) == 2 && result.Scenario.MergeSnapshotCID != "" {
		checks["multi_parent_merge"] = "kept"
	} else {
		checks["multi_parent_merge"] = "missing"
		pass = false
	}
	if result.Scenario.TestStatementCID != "" && result.Scenario.AdoptionStatementCID != "" && result.Scenario.ReviewAdoptionResult == "accepted_locally" {
		checks["review_test_local_adoption"] = "kept"
	} else {
		checks["review_test_local_adoption"] = "missing"
		pass = false
	}
	// Intent: The bridge fixture proves conventional Git import/export/push/pull
	// remain adapter operations with explicit mapping promises, not native sync or
	// forge authority. Source: DI-fimap
	if bridgeResultComplete(result.Bridge.Export) && bridgeResultComplete(result.Bridge.Import) && bridgeResultComplete(result.Bridge.Push) && bridgeResultComplete(result.Bridge.Pull) {
		checks["git_bridge_import_export_push_pull"] = "kept"
	} else {
		checks["git_bridge_import_export_push_pull"] = "missing"
		pass = false
	}
	aliceCAS, aliceOpenErr := store.Open(result.AliceStoreRoot)
	if aliceOpenErr != nil {
		return Report{}, aliceOpenErr
	}
	aliceEntries, aliceListErr := aliceCAS.List()
	if aliceListErr != nil {
		return Report{}, aliceListErr
	}
	bobCAS, bobOpenErr := store.Open(result.BobStoreRoot)
	if bobOpenErr != nil {
		return Report{}, bobOpenErr
	}
	bobEntries, bobListErr := bobCAS.List()
	if bobListErr != nil {
		return Report{}, bobListErr
	}
	carolCAS, carolOpenErr := store.Open(result.CarolStoreRoot)
	if carolOpenErr != nil {
		return Report{}, carolOpenErr
	}
	carolEntries, carolListErr := carolCAS.List()
	if carolListErr != nil {
		return Report{}, carolListErr
	}
	frankCAS, frankOpenErr := store.Open(result.FrankStoreRoot)
	if frankOpenErr != nil {
		return Report{}, frankOpenErr
	}
	frankEntries, frankListErr := frankCAS.List()
	if frankListErr != nil {
		return Report{}, frankListErr
	}
	malloryCAS, malloryOpenErr := store.Open(result.MalloryStoreRoot)
	if malloryOpenErr != nil {
		return Report{}, malloryOpenErr
	}
	malloryEntries, malloryListErr := malloryCAS.List()
	if malloryListErr != nil {
		return Report{}, malloryListErr
	}
	if retentionPass, retentionErr := checkRetentionGC(result, frankCAS, checks); retentionErr != nil {
		return Report{}, retentionErr
	} else if !retentionPass {
		pass = false
	}
	if paymentPass, paymentErr := checkRetentionPayment(result, aliceCAS, frankCAS, checks); paymentErr != nil {
		return Report{}, paymentErr
	} else if !paymentPass {
		pass = false
	}
	if continuousPass, continuousErr := checkContinuousSync(result, bobCAS, carolCAS, checks); continuousErr != nil {
		return Report{}, continuousErr
	} else if !continuousPass {
		pass = false
	}
	if scheduledPass, scheduledErr := checkScheduledSync(result, bobCAS, carolCAS, checks); scheduledErr != nil {
		return Report{}, scheduledErr
	} else if !scheduledPass {
		pass = false
	}
	snapshotCID, snapshotErr := store.ParseCIDText(result.SnapshotCID)
	if snapshotErr != nil {
		return Report{}, snapshotErr
	}
	if bobCAS.Has(snapshotCID) {
		checks["bob_has_snapshot"] = "kept"
	} else {
		checks["bob_has_snapshot"] = "missing"
		pass = false
	}
	if result.Retrieval.InterestMessageCID != "" {
		interestCID, interestErr := store.ParseCIDText(result.Retrieval.InterestMessageCID)
		if interestErr != nil {
			return Report{}, interestErr
		}
		if bobCAS.Has(interestCID) {
			checks["bob_has_sync_interest"] = "kept"
		} else {
			checks["bob_has_sync_interest"] = "missing"
			pass = false
		}
		if aliceCAS.Has(interestCID) {
			checks["alice_has_received_sync_interest"] = "kept"
		} else {
			checks["alice_has_received_sync_interest"] = "missing"
			pass = false
		}
	}

	if result.Retrieval.AvailabilityMessageCID != "" {
		availabilityCID, availabilityErr := store.ParseCIDText(result.Retrieval.AvailabilityMessageCID)
		if availabilityErr != nil {
			return Report{}, availabilityErr
		}
		if bobCAS.Has(availabilityCID) {
			checks["bob_has_object_availability"] = "kept"
		} else {
			checks["bob_has_object_availability"] = "missing"
			pass = false
		}
		if aliceCAS.Has(availabilityCID) {
			checks["alice_has_object_availability"] = "kept"
		} else {
			checks["alice_has_object_availability"] = "missing"
			pass = false
		}
	}
	for checkName, cidText := range map[string]string{
		"alice_has_git_export_mapping": result.Bridge.Export.MappingMessageCID,
		"alice_has_git_import_mapping": result.Bridge.Import.MappingMessageCID,
		"alice_has_git_push_mapping":   result.Bridge.Push.MappingMessageCID,
		"bob_has_git_pull_mapping":     result.Bridge.Pull.MappingMessageCID,
		"bob_has_git_pull_snapshot":    result.Bridge.Pull.HeadSnapshotCID,
	} {
		if cidText == "" {
			checks[checkName] = "missing"
			pass = false
			continue
		}
		objectCID, objectErr := store.ParseCIDText(cidText)
		if objectErr != nil {
			return Report{}, objectErr
		}
		targetCAS := aliceCAS
		if checkName == "bob_has_git_pull_mapping" || checkName == "bob_has_git_pull_snapshot" {
			targetCAS = bobCAS
		}
		if targetCAS.Has(objectCID) {
			checks[checkName] = "kept"
		} else {
			checks[checkName] = "missing"
			pass = false
		}
	}
	for checkName, cidText := range map[string]string{
		"bob_has_rename_copy_snapshot":   result.Scenario.RenameCopySnapshotCID,
		"bob_has_bob_divergent_snapshot": result.Scenario.BobDivergentSnapshotCID,
		"bob_has_merge_snapshot":         result.Scenario.MergeSnapshotCID,
		"bob_has_test_statement":         result.Scenario.TestStatementCID,
		"bob_has_adoption_statement":     result.Scenario.AdoptionStatementCID,
		"bob_has_review_thread":          result.Scenario.ReviewThreadRefSetCID,
	} {
		if cidText == "" {
			checks[checkName] = "missing"
			pass = false
			continue
		}
		objectCID, objectErr := store.ParseCIDText(cidText)
		if objectErr != nil {
			return Report{}, objectErr
		}
		if bobCAS.Has(objectCID) {
			checks[checkName] = "kept"
		} else {
			checks[checkName] = "missing"
			pass = false
		}
	}
	return Report{
		RunRoot:        runRoot,
		Pass:           pass,
		Checks:         checks,
		Objects:        len(bobEntries),
		AliceObjects:   len(aliceEntries),
		BobObjects:     len(bobEntries),
		CarolObjects:   len(carolEntries),
		FrankObjects:   len(frankEntries),
		MalloryObjects: len(malloryEntries),
		Retrieved:      len(result.Retrieval.Retrieved),
		Missing:        len(result.Retrieval.Missing),
	}, nil
}

func analyzeCollector(runRoot string) (Report, error) {
	checks := map[string]string{}
	pass := true
	events, eventsErr := readJSONL[collectorEvent](filepath.Join(runRoot, "events.jsonl"))
	if eventsErr != nil {
		return Report{}, eventsErr
	}
	messages, messagesErr := readJSONL[collectorMessageRecord](filepath.Join(runRoot, "message-dag.jsonl"))
	if messagesErr != nil {
		return Report{}, messagesErr
	}
	cars, carsErr := readJSONL[collectorCARRecord](filepath.Join(runRoot, "car-dag.jsonl"))
	if carsErr != nil {
		return Report{}, carsErr
	}
	if len(events) > 0 {
		checks["collector_events"] = "kept"
	} else {
		checks["collector_events"] = "missing"
		pass = false
	}
	if len(messages) > 0 {
		checks["collector_message_artifacts"] = "kept"
	} else {
		checks["collector_message_artifacts"] = "missing"
		pass = false
	}
	if len(cars) > 0 {
		checks["collector_car_artifacts"] = "kept"
	} else {
		checks["collector_car_artifacts"] = "missing"
		pass = false
	}
	// Intent: The collector-mode analyzer now checks the runtime signals that
	// prove POC18's remediation is real TCP session behavior, not fixture-only
	// object copying. These checks remain post-run diagnostics and do not become
	// a production monitor. Source: DI-biruf
	if requiredEventsPresent(events, checks, []string{
		"session_opened",
		"session_closed",
		"exchange_started",
		"exchange_completed",
		"scheduler_tcp_sync",
		"storage_payment_token_issued",
		"storage_payment_redemption",
		"storage_capability_token_issued",
		"dag_closure_missing",
		"dag_repair_completed",
		"trust_peer_choice",
		"workspace_scenario_rich",
	}) {
		checks["collector_runtime_remediation_events"] = "kept"
	} else {
		checks["collector_runtime_remediation_events"] = "missing"
		pass = false
	}
	completedAgents := map[string]bool{}
	agentFailures := []string{}
	for _, event := range events {
		if event.Event == "agent_completed" && event.Outcome == "kept" {
			completedAgents[event.Observer] = true
		}
		if event.Event == "agent_failed" {
			agentFailures = append(agentFailures, event.Observer)
		}
	}
	for _, agentName := range []string{"alice", "bob", "carol", "dave", "ellen", "frank", "mallory"} {
		checkName := "agent_completed_" + agentName
		if completedAgents[agentName] {
			checks[checkName] = "kept"
		} else {
			checks[checkName] = "missing"
			pass = false
		}
	}
	if len(agentFailures) == 0 {
		checks["no_agent_failed_events"] = "kept"
	} else {
		checks["no_agent_failed_events"] = strings.Join(agentFailures, ",")
		pass = false
	}
	if collectorPromiseKindsPresent(messages, checks) {
		checks["tcp_promise_flow"] = "kept"
	} else {
		checks["tcp_promise_flow"] = "missing"
		pass = false
	}
	if collectorDirectionsPresent(messages) {
		checks["sent_and_received_artifacts"] = "kept"
	} else {
		checks["sent_and_received_artifacts"] = "missing"
		pass = false
	}
	if messagesPass, verifyErr := verifyMessageArtifacts(runRoot, messages, checks); verifyErr != nil {
		return Report{}, verifyErr
	} else if !messagesPass {
		pass = false
	}
	if carsPass, verifyErr := verifyCARArtifacts(runRoot, cars, checks); verifyErr != nil {
		return Report{}, verifyErr
	} else if !carsPass {
		pass = false
	}
	return Report{
		RunRoot:           runRoot,
		Pass:              pass,
		Checks:            checks,
		Objects:           len(messages),
		CollectorEvents:   len(events),
		CollectorMessages: len(messages),
		CollectorCARs:     len(cars),
	}, nil
}

func requiredEventsPresent(events []collectorEvent, checks map[string]string, required []string) bool {
	present := map[string]bool{}
	for _, event := range events {
		if event.Outcome == "kept" {
			present[event.Event] = true
		}
	}
	pass := true
	for _, eventName := range required {
		checkName := "event_" + eventName
		if present[eventName] {
			checks[checkName] = "kept"
			continue
		}
		checks[checkName] = "missing"
		pass = false
	}
	return pass
}

func collectorPromiseKindsPresent(messages []collectorMessageRecord, checks map[string]string) bool {
	present := map[string]bool{}
	for _, message := range messages {
		present[message.PromiseKind] = true
	}
	pass := true
	for _, promiseKind := range []string{"sync_interest", "object_availability", "object_retrieval_redemption", "object_bytes"} {
		checkName := "promise_kind_" + promiseKind
		if present[promiseKind] {
			checks[checkName] = "kept"
		} else {
			checks[checkName] = "missing"
			pass = false
		}
	}
	return pass
}

func collectorDirectionsPresent(messages []collectorMessageRecord) bool {
	directions := map[string]bool{}
	for _, message := range messages {
		directions[message.Direction] = true
	}
	return directions["sent"] && directions["received"]
}

func verifyMessageArtifacts(runRoot string, messages []collectorMessageRecord, checks map[string]string) (bool, error) {
	pass := true
	for _, message := range messages {
		content, readErr := os.ReadFile(filepath.Join(runRoot, filepath.FromSlash(message.Path)))
		if readErr != nil {
			return false, readErr
		}
		actualCID := store.CIDText(store.CIDForBytes(content))
		if actualCID != message.ExactCID {
			checks["message_cid_match_"+message.ExactCID] = "missing"
			pass = false
			continue
		}
		if _, parseErr := graph.ParseEnvelope(content); parseErr != nil {
			checks["message_parse_"+message.ExactCID] = "missing"
			pass = false
			continue
		}
		checks["message_cid_match_"+message.ExactCID] = "kept"
	}
	return pass, nil
}

func verifyCARArtifacts(runRoot string, cars []collectorCARRecord, checks map[string]string) (bool, error) {
	pass := true
	for _, carRecord := range cars {
		content, readErr := os.ReadFile(filepath.Join(runRoot, filepath.FromSlash(carRecord.Path)))
		if readErr != nil {
			return false, readErr
		}
		actualCID := store.CIDText(store.CIDForBytes(content))
		if actualCID != carRecord.CARCID {
			checks["car_cid_match_"+carRecord.CARCID] = "missing"
			pass = false
			continue
		}
		if len(carRecord.BlockCIDs) == 0 {
			checks["car_blocks_"+carRecord.CARCID] = "missing"
			pass = false
			continue
		}
		if verifyErr := carbundle.VerifyStandard(content, carRecord.BlockCIDs, carRecord.BlockCIDs); verifyErr != nil {
			checks["car_standard_"+carRecord.CARCID] = "missing"
			pass = false
			continue
		}
		checks["car_cid_match_"+carRecord.CARCID] = "kept"
		checks["car_standard_"+carRecord.CARCID] = "kept"
	}
	return pass, nil
}

func readJSONL[T any](path string) ([]T, error) {
	file, openErr := os.Open(path)
	if openErr != nil {
		return nil, openErr
	}
	defer closeReadFile(file, path)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), 32*1024*1024)
	values := []T{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var value T
		if unmarshalErr := json.Unmarshal([]byte(line), &value); unmarshalErr != nil {
			return nil, unmarshalErr
		}
		values = append(values, value)
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return nil, scanErr
	}
	return values, nil
}

func closeReadFile(file *os.File, path string) {
	if closeErr := file.Close(); closeErr != nil {
		fmt.Fprintf(os.Stderr, "close %s: %v\n", path, closeErr)
	}
}

func checkScheduledSync(result fixtureResult, bobCAS *store.FileStore, carolCAS *store.FileStore, checks map[string]string) (bool, error) {
	// Intent: The scheduler gate verifies local sync-agent behavior, not a global
	// scheduler. Bob should win from Carol's retained graph evidence, while Mallory
	// remains only market-signal evidence. Source: DI-fakop
	pass := true
	if result.ScheduledSync.ChosenPeer == "bob" {
		checks["scheduled_sync_chose_bob"] = "kept"
	} else {
		checks["scheduled_sync_chose_bob"] = "missing"
		pass = false
	}
	if scheduledDecisionHasLocalGraph(result.ScheduledSync, "bob") {
		checks["scheduled_sync_bob_from_local_graph"] = "kept"
	} else {
		checks["scheduled_sync_bob_from_local_graph"] = "missing"
		pass = false
	}
	if scheduledDecisionUsedMarketFallback(result.ScheduledSync, "mallory") {
		checks["scheduled_sync_market_fallback_for_mallory"] = "kept"
	} else {
		checks["scheduled_sync_market_fallback_for_mallory"] = "missing"
		pass = false
	}
	if len(result.ScheduledSync.CapabilityRedemptions) > 0 && result.ScheduledSync.CapabilityRedemptions[0].Redeemed && result.ScheduledSync.CapabilityRedemptions[0].Capability == "storage" {
		checks["scheduled_sync_capability_redeemed"] = "kept"
	} else {
		checks["scheduled_sync_capability_redeemed"] = "missing"
		pass = false
	}
	if len(result.ScheduledSync.CapabilityRedemptions) > 0 {
		redemption := result.ScheduledSync.CapabilityRedemptions[0]
		bearerCID, bearerErr := store.ParseCIDText(redemption.BearerTokenCID)
		if bearerErr != nil {
			return false, bearerErr
		}
		capabilityCID, capabilityErr := store.ParseCIDText(redemption.CapabilityTokenCID)
		if capabilityErr != nil {
			return false, capabilityErr
		}
		if bobCAS.Has(bearerCID) && carolCAS.Has(bearerCID) {
			checks["scheduled_sync_bearer_token_exchanged"] = "kept"
		} else {
			checks["scheduled_sync_bearer_token_exchanged"] = "missing"
			pass = false
		}
		if bobCAS.Has(capabilityCID) && carolCAS.Has(capabilityCID) {
			checks["scheduled_sync_capability_token_exchanged"] = "kept"
		} else {
			checks["scheduled_sync_capability_token_exchanged"] = "missing"
			pass = false
		}
	}
	if result.SyncAgentState.LastReport != nil && result.SyncAgentState.LastReport.ChosenPeer == "bob" && len(result.SyncAgentState.RedeemedBearerTokenCIDs) > 0 {
		checks["scheduled_sync_state_checkpointed"] = "kept"
	} else {
		checks["scheduled_sync_state_checkpointed"] = "missing"
		pass = false
	}
	if result.CarolSyncStatePath != "" {
		if _, statErr := os.Stat(result.CarolSyncStatePath); statErr == nil {
			checks["scheduled_sync_state_file_written"] = "kept"
		} else {
			checks["scheduled_sync_state_file_written"] = "missing"
			pass = false
		}
	}
	if result.ScheduledSync.ContinuousSync.MissingObjects == 0 {
		checks["scheduled_sync_no_missing"] = "kept"
	} else {
		checks["scheduled_sync_no_missing"] = "missing"
		pass = false
	}
	return pass, nil
}

func scheduledDecisionHasLocalGraph(report pocsync.SchedulerReport, peer string) bool {
	for _, decision := range report.Decisions {
		if decision.Peer == peer && decision.Result == "chosen" && decision.Trust.KeptPromises > 0 && !decision.Trust.MarketSignalUsed {
			return true
		}
	}
	return false
}

func scheduledDecisionUsedMarketFallback(report pocsync.SchedulerReport, peer string) bool {
	for _, decision := range report.Decisions {
		if decision.Peer == peer && decision.Trust.MarketSignalUsed {
			return true
		}
	}
	return false
}

func checkContinuousSync(result fixtureResult, bobCAS *store.FileStore, carolCAS *store.FileStore, checks map[string]string) (bool, error) {
	// Intent: Continuous sync checks prove repeated peer DAG exchange converges
	// through promise messages and sparse CAS verification. They are local fixture
	// gates only, not production-wide monitoring. Source: DI-rudos
	pass := true
	if result.CarolStoreRoot != "" && result.CarolStoreRoot != result.AliceStoreRoot && result.CarolStoreRoot != result.BobStoreRoot && result.CarolStoreRoot != result.FrankStoreRoot {
		checks["continuous_sync_separate_carol_cas"] = "kept"
	} else {
		checks["continuous_sync_separate_carol_cas"] = "missing"
		pass = false
	}
	if result.ContinuousSync.UsefulUpdates > 0 {
		checks["continuous_sync_useful_update"] = "kept"
	} else {
		checks["continuous_sync_useful_update"] = "missing"
		pass = false
	}
	if result.ContinuousSync.MissingObjects == 0 {
		checks["continuous_sync_no_missing"] = "kept"
	} else {
		checks["continuous_sync_no_missing"] = "missing"
		pass = false
	}
	if continuousSecondRoundIsIdempotent(result.ContinuousSync) {
		checks["continuous_sync_second_round_idempotent"] = "kept"
	} else {
		checks["continuous_sync_second_round_idempotent"] = "missing"
		pass = false
	}
	if continuousRetentionExchanged(result.ContinuousSync, bobCAS, carolCAS) {
		checks["continuous_sync_retention_message_exchanged"] = "kept"
	} else {
		checks["continuous_sync_retention_message_exchanged"] = "missing"
		pass = false
	}
	for checkName, cidText := range map[string]string{
		"carol_has_merge_branch":   result.Scenario.MergeBranchRefSetCID,
		"carol_has_review_thread":  result.Scenario.ReviewThreadRefSetCID,
		"carol_has_merge_snapshot": result.Scenario.MergeSnapshotCID,
	} {
		objectCID, objectErr := store.ParseCIDText(cidText)
		if objectErr != nil {
			return false, objectErr
		}
		if carolCAS.Has(objectCID) {
			checks[checkName] = "kept"
		} else {
			checks[checkName] = "missing"
			pass = false
		}
	}
	return pass, nil
}

func continuousSecondRoundIsIdempotent(report pocsync.ContinuousSyncReport) bool {
	for _, direction := range report.Directions {
		if direction.Round == 2 && !direction.UsefulUpdate && len(direction.Retrieval.Retrieved) == 0 && len(direction.Retrieval.AlreadyLocal) > 0 {
			return true
		}
	}
	return false
}

func continuousRetentionExchanged(report pocsync.ContinuousSyncReport, bobCAS *store.FileStore, carolCAS *store.FileStore) bool {
	for _, direction := range report.Directions {
		if direction.RetentionMessageCID == "" {
			continue
		}
		retentionCID, retentionErr := store.ParseCIDText(direction.RetentionMessageCID)
		if retentionErr != nil {
			return false
		}
		if direction.Provider == "bob" && direction.Receiver == "carol" && bobCAS.Has(retentionCID) && carolCAS.Has(retentionCID) {
			return true
		}
	}
	return false
}

func checkRetentionGC(result fixtureResult, frankCAS *store.FileStore, checks map[string]string) (bool, error) {
	pass := true
	if result.FrankStoreRoot != "" && result.FrankStoreRoot != result.AliceStoreRoot && result.FrankStoreRoot != result.BobStoreRoot {
		checks["retention_separate_frank_cas"] = "kept"
	} else {
		checks["retention_separate_frank_cas"] = "missing"
		pass = false
	}
	if result.Retention.RetentionMessageCID != "" {
		retentionCID, retentionErr := store.ParseCIDText(result.Retention.RetentionMessageCID)
		if retentionErr != nil {
			return false, retentionErr
		}
		if frankCAS.Has(retentionCID) {
			checks["retention_message_kept"] = "kept"
		} else {
			checks["retention_message_kept"] = "missing"
			pass = false
		}
	} else {
		checks["retention_message_kept"] = "missing"
		pass = false
	}
	if len(result.Retention.MissingProtectedCIDs) == 0 {
		checks["retention_promised_closure_available"] = "kept"
	} else {
		checks["retention_promised_closure_available"] = "missing"
		pass = false
	}
	for _, target := range result.Retention.Targets {
		targetCID, targetErr := store.ParseCIDText(target.CID)
		if targetErr != nil {
			return false, targetErr
		}
		checkName := "retention_kept_" + target.Role
		if frankCAS.Has(targetCID) {
			checks[checkName] = "kept"
		} else {
			checks[checkName] = "missing"
			pass = false
		}
	}
	if result.Retention.CollectedObjects > 0 {
		checks["retention_collected_unpromised"] = "kept"
	} else {
		checks["retention_collected_unpromised"] = "missing"
		pass = false
	}
	if result.Retention.TemporaryObjectCID != "" {
		tempCID, tempErr := store.ParseCIDText(result.Retention.TemporaryObjectCID)
		if tempErr != nil {
			return false, tempErr
		}
		if !frankCAS.Has(tempCID) {
			checks["retention_temp_pressure_object_collected"] = "kept"
		} else {
			checks["retention_temp_pressure_object_collected"] = "missing"
			pass = false
		}
	}
	return pass, nil
}

func checkRetentionPayment(result fixtureResult, aliceCAS *store.FileStore, frankCAS *store.FileStore, checks map[string]string) (bool, error) {
	// Intent: The clean run must prove retention economics are now signed,
	// spendable token promises. These checks remain local fixture gates, not a
	// global monitor or payment authority. Source: DI-bidum
	pass := true
	payment := result.RetentionPayment
	if payment.Issued && payment.SignatureVerified && payment.Transferable && payment.Value == 5 && payment.Unit == "storage_credit" && payment.ObjectCID == result.ReleaseRefSetCID {
		checks["retention_payment_token_signed"] = "kept"
	} else {
		checks["retention_payment_token_signed"] = "missing"
		pass = false
	}
	if payment.Redeemed && payment.TokenStoredByRedeemer {
		checks["retention_payment_token_redeemed"] = "kept"
	} else {
		checks["retention_payment_token_redeemed"] = "missing"
		pass = false
	}
	if payment.ReplayRejected {
		checks["retention_payment_replay_rejected"] = "kept"
	} else {
		checks["retention_payment_replay_rejected"] = "missing"
		pass = false
	}
	tokenCID, tokenErr := store.ParseCIDText(payment.TokenCID)
	if tokenErr != nil {
		return false, tokenErr
	}
	if payment.TokenStoredByIssuer && aliceCAS.Has(tokenCID) {
		checks["retention_payment_token_in_alice_cas"] = "kept"
	} else {
		checks["retention_payment_token_in_alice_cas"] = "missing"
		pass = false
	}
	if frankCAS.Has(tokenCID) {
		checks["retention_payment_token_in_frank_cas"] = "kept"
	} else {
		checks["retention_payment_token_in_frank_cas"] = "missing"
		pass = false
	}
	redemptionCID, redemptionErr := store.ParseCIDText(payment.RedemptionMessageCID)
	if redemptionErr != nil {
		return false, redemptionErr
	}
	if frankCAS.Has(redemptionCID) {
		checks["retention_payment_redemption_message_kept"] = "kept"
	} else {
		checks["retention_payment_redemption_message_kept"] = "missing"
		pass = false
	}
	if cidListContains(result.Retention.ProtectedCIDs, payment.TokenCID) && cidListContains(result.Retention.ProtectedCIDs, payment.RedemptionMessageCID) {
		checks["retention_payment_protected_by_gc"] = "kept"
	} else {
		checks["retention_payment_protected_by_gc"] = "missing"
		pass = false
	}
	return pass, nil
}

func cidListContains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func checkAliceGridRepo(result fixtureResult, checks map[string]string) (bool, error) {
	pass := true
	repository, discoverErr := pgrepo.Discover(result.SourceRoot)
	if discoverErr != nil {
		checks["alice_grid_repo_discovery"] = "missing"
		return false, nil
	}
	checks["alice_grid_repo_discovery"] = "kept"
	if filepath.Clean(repository.ResolvePath(repository.Config.CAS.Path)) == filepath.Clean(result.AliceStoreRoot) {
		checks["alice_grid_repo_cas_locator"] = "kept"
	} else {
		checks["alice_grid_repo_cas_locator"] = "missing"
		pass = false
	}
	state, stateErr := repository.LoadState()
	if stateErr != nil {
		return false, stateErr
	}
	if state.CurrentSnapshotCID == result.SnapshotCID {
		checks["alice_grid_state_snapshot"] = "kept"
	} else {
		checks["alice_grid_state_snapshot"] = "missing"
		pass = false
	}
	if state.CurrentBranchRefSetCID == result.BranchRefSetCID {
		checks["alice_grid_state_branch"] = "kept"
	} else {
		checks["alice_grid_state_branch"] = "missing"
		pass = false
	}
	cas, casErr := repository.OpenFileCAS()
	if casErr != nil {
		return false, casErr
	}
	snapshotCID, snapshotErr := store.ParseCIDText(state.CurrentSnapshotCID)
	if snapshotErr != nil {
		return false, snapshotErr
	}
	// Intent: Analyzer coverage should catch the fixture mismatch where Alice's
	// state names a final snapshot but her filesystem still contains stale initial
	// fixture labels. Source: DI-bamum
	statusReport, statusErr := workspace.CompareSnapshot(cas, snapshotCID, result.SourceRoot)
	if statusErr != nil {
		return false, statusErr
	}
	if statusReport.Clean {
		checks["alice_grid_workspace_clean"] = "kept"
	} else {
		checks["alice_grid_workspace_clean"] = "missing"
		pass = false
	}
	return pass, nil
}

func bridgeResultComplete(result bridge.Result) bool {
	return result.HeadSnapshotCID != "" &&
		result.HeadGitHash != "" &&
		result.MappingMessageCID != "" &&
		len(result.Mappings) > 0
}

func writeJSON(path string, value any) error {
	content, marshalErr := json.MarshalIndent(value, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	content = append(content, '\n')
	return os.WriteFile(path, content, 0o644)
}
