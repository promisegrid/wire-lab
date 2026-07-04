// Command poc-analyze checks deterministic POC18 fixture results.
//
// Intent: Analyzer checks remain non-production review aids; they must not imply
// a global monitor or authority. Source: DI-jifuj
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/bridge"
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
	AliceStoreRoot string                  `json:"alice_store_root"`
	BobStoreRoot   string                  `json:"bob_store_root"`
	FrankStoreRoot string                  `json:"frank_store_root"`
	Scenario       scenario.Result         `json:"scenario"`
	Bridge         bridgeFixtureResult     `json:"bridge"`
	Retrieval      pocsync.RetrievalReport `json:"retrieval"`
	Retention      retention.Report        `json:"retention"`
}

type bridgeFixtureResult struct {
	Export bridge.Result `json:"export"`
	Import bridge.Result `json:"import"`
	Push   bridge.Result `json:"push"`
	Pull   bridge.Result `json:"pull"`
}

// Report records local fixture checks.
type Report struct {
	RunRoot      string            `json:"run_root"`
	Pass         bool              `json:"pass"`
	Checks       map[string]string `json:"checks"`
	Objects      int               `json:"objects"`
	AliceObjects int               `json:"alice_objects"`
	BobObjects   int               `json:"bob_objects"`
	FrankObjects int               `json:"frank_objects"`
	Retrieved    int               `json:"retrieved"`
	Missing      int               `json:"missing"`
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
	result, resultErr := readResult(filepath.Join(*runRoot, "result.json"))
	if resultErr != nil {
		return resultErr
	}
	report, analyzeErr := analyze(*runRoot, result)
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
	frankCAS, frankOpenErr := store.Open(result.FrankStoreRoot)
	if frankOpenErr != nil {
		return Report{}, frankOpenErr
	}
	frankEntries, frankListErr := frankCAS.List()
	if frankListErr != nil {
		return Report{}, frankListErr
	}
	if retentionPass, retentionErr := checkRetentionGC(result, frankCAS, checks); retentionErr != nil {
		return Report{}, retentionErr
	} else if !retentionPass {
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
		RunRoot:      runRoot,
		Pass:         pass,
		Checks:       checks,
		Objects:      len(bobEntries),
		AliceObjects: len(aliceEntries),
		BobObjects:   len(bobEntries),
		FrankObjects: len(frankEntries),
		Retrieved:    len(result.Retrieval.Retrieved),
		Missing:      len(result.Retrieval.Missing),
	}, nil
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
