// Command poc-sim runs the deterministic non-production POC18 fixture.
//
// Intent: `poc-*` commands are fixtures, not production daemons. They exercise
// the same core library that `grid` uses. Source: DI-jifuj; DI-harih
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	cidlib "github.com/ipfs/go-cid"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/bridge"
	pgrepo "promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/repo"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/scenario"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	pocsync "promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/sync"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/workspace"
)

// fixtureResult records the deterministic multi-agent fixture state.
//
// Intent: The fixture must prove Bob's checkout comes from Bob's sparse CAS
// after CID-verified retrieval from Alice, not from a shared store shortcut.
// Source: DI-gozov
type fixtureResult struct {
	workspace.IngestResult
	AliceStoreRoot string                  `json:"alice_store_root"`
	BobStoreRoot   string                  `json:"bob_store_root"`
	Scenario       scenario.Result         `json:"scenario"`
	Bridge         bridgeFixtureResult     `json:"bridge"`
	Retrieval      pocsync.RetrievalReport `json:"retrieval"`
}

type bridgeFixtureResult struct {
	Export bridge.Result `json:"export"`
	Import bridge.Result `json:"import"`
	Push   bridge.Result `json:"push"`
	Pull   bridge.Result `json:"pull"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "poc-sim: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	runRoot := flag.String("run-root", "/tmp/wire-lab-poc18-run", "approved POC18 runtime root")
	flag.Parse()
	if err := resetRunRoot(*runRoot); err != nil {
		return err
	}
	sourceRoot := filepath.Join(*runRoot, "alice-workspace")
	aliceStoreRoot := filepath.Join(*runRoot, "alice-cas")
	bobStoreRoot := filepath.Join(*runRoot, "bob-cas")
	checkoutRoot := filepath.Join(*runRoot, "bob-checkout")
	if err := createFixture(sourceRoot); err != nil {
		return err
	}
	aliceRepository, repoErr := pgrepo.Init(sourceRoot, filepath.Join("..", "alice-cas"))
	if repoErr != nil {
		return repoErr
	}
	aliceCAS, aliceOpenErr := aliceRepository.OpenFileCAS()
	if aliceOpenErr != nil {
		return aliceOpenErr
	}
	result, ingestErr := workspace.NewScanner(aliceCAS, "alice", "bob").Ingest(sourceRoot)
	if ingestErr != nil {
		return ingestErr
	}
	scenarioResult, scenarioErr := scenario.NewBuilder(aliceCAS).Build(result)
	if scenarioErr != nil {
		return scenarioErr
	}
	mergeSnapshotCID, mergeSnapshotErr := store.ParseCIDText(scenarioResult.MergeSnapshotCID)
	if mergeSnapshotErr != nil {
		return mergeSnapshotErr
	}
	mergeBranchCID, mergeBranchErr := store.ParseCIDText(scenarioResult.MergeBranchRefSetCID)
	if mergeBranchErr != nil {
		return mergeBranchErr
	}
	reviewThreadCID, reviewThreadErr := store.ParseCIDText(scenarioResult.ReviewThreadRefSetCID)
	if reviewThreadErr != nil {
		return reviewThreadErr
	}
	result.SnapshotCID = scenarioResult.MergeSnapshotCID
	result.RootReferenceSetCID = scenarioResult.MergeRootRefSetCID
	result.BranchRefSetCID = scenarioResult.MergeBranchRefSetCID
	result.LogicalChangeCID = scenarioResult.LogicalChangeRefSetCID
	result.ReviewThreadCID = scenarioResult.ReviewThreadRefSetCID
	result.DiagnosticMessageCID = scenarioResult.MergeSnapshotCID
	for countName, countValue := range scenarioResult.Counts {
		result.Counts[countName] += countValue
	}
	if materializeErr := refreshWorkspaceToSnapshot(aliceCAS, mergeSnapshotCID, sourceRoot); materializeErr != nil {
		return materializeErr
	}
	// Intent: The deterministic fixture should leave Alice's workspace in the
	// same repo-local shape that normal `grid` users inspect: config points at the
	// selected sparse CAS, and state records Alice's current local head without
	// making it a global branch authority. Source: DI-kiram
	if _, recordErr := aliceRepository.RecordSnapshot(result.SnapshotCID, result.BranchRefSetCID, result.WorkspaceRefSetCID); recordErr != nil {
		return recordErr
	}
	bobCAS, bobOpenErr := store.Open(bobStoreRoot)
	if bobOpenErr != nil {
		return bobOpenErr
	}
	bridgeResult, bridgeErr := runGitBridgeFixture(*runRoot, aliceCAS, bobCAS)
	if bridgeErr != nil {
		return bridgeErr
	}
	retrieval, retrieveErr := pocsync.RetrieveGraph(
		pocsync.Peer{Agent: "bob", CAS: bobCAS},
		pocsync.Peer{Agent: "alice", CAS: aliceCAS},
		map[string]cidlib.Cid{"branch": mergeBranchCID, "review_thread": reviewThreadCID},
		"storage_credit:2",
	)
	if retrieveErr != nil {
		return retrieveErr
	}
	if err := workspace.MaterializeSnapshot(bobCAS, mergeSnapshotCID, checkoutRoot); err != nil {
		return err
	}
	result.CheckoutRoot = checkoutRoot
	combined := fixtureResult{
		IngestResult:   result,
		AliceStoreRoot: aliceStoreRoot,
		BobStoreRoot:   bobStoreRoot,
		Scenario:       scenarioResult,
		Bridge:         bridgeResult,
		Retrieval:      retrieval,
	}
	resultPath := filepath.Join(*runRoot, "result.json")
	if err := writeJSON(resultPath, combined); err != nil {
		return err
	}
	fmt.Printf("run_root=%s\n", *runRoot)
	fmt.Printf("alice_store=%s\n", aliceStoreRoot)
	fmt.Printf("bob_store=%s\n", bobStoreRoot)
	fmt.Printf("snapshot=%s\n", result.SnapshotCID)
	fmt.Printf("diagnostic_message=%s\n", result.DiagnosticMessageCID)
	fmt.Printf("retrieved_objects=%d missing_objects=%d\n", len(retrieval.Retrieved), len(retrieval.Missing))
	return nil
}

func runGitBridgeFixture(runRoot string, aliceCAS *store.FileStore, bobCAS *store.FileStore) (bridgeFixtureResult, error) {
	// Intent: The deterministic fixture covers Git bridge import/export/push/pull
	// as adapter behavior while native peer sync remains the separate retrieval
	// promise path. Source: DI-fimap
	gitWorkspace := filepath.Join(runRoot, "git-workspace")
	if err := createGitBridgeFixture(gitWorkspace); err != nil {
		return bridgeFixtureResult{}, err
	}
	gitIngest, ingestErr := workspace.NewScanner(aliceCAS, "alice", "bob").Ingest(gitWorkspace)
	if ingestErr != nil {
		return bridgeFixtureResult{}, ingestErr
	}
	gitSnapshotCID, parseErr := store.ParseCIDText(gitIngest.SnapshotCID)
	if parseErr != nil {
		return bridgeFixtureResult{}, parseErr
	}
	aliceBridge := bridge.NewAdapter(aliceCAS, "alice", "bob")
	exportResult, exportErr := aliceBridge.ExportSnapshot(gitSnapshotCID, filepath.Join(runRoot, "git-export"))
	if exportErr != nil {
		return bridgeFixtureResult{}, exportErr
	}
	importResult, importErr := aliceBridge.ImportRepository(exportResult.RepositoryPath)
	if importErr != nil {
		return bridgeFixtureResult{}, importErr
	}
	remotePath := filepath.Join(runRoot, "git-remote.git")
	if err := initBareMainRepository(remotePath); err != nil {
		return bridgeFixtureResult{}, err
	}
	pushResult, pushErr := aliceBridge.PushSnapshot(gitSnapshotCID, remotePath, filepath.Join(runRoot, "git-push-worktree"))
	if pushErr != nil {
		return bridgeFixtureResult{}, pushErr
	}
	pullResult, pullErr := bridge.NewAdapter(bobCAS, "bob", "alice").PullRepository(remotePath, filepath.Join(runRoot, "git-pull-worktree"))
	if pullErr != nil {
		return bridgeFixtureResult{}, pullErr
	}
	return bridgeFixtureResult{Export: exportResult, Import: importResult, Push: pushResult, Pull: pullResult}, nil
}

func createGitBridgeFixture(root string) error {
	if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("POC18 Git bridge fixture\n"), 0o644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(root, "docs", "bridge.txt"), []byte("export import push pull\n"), 0o644)
}

func refreshWorkspaceToSnapshot(cas *store.FileStore, snapshotCID cidlib.Cid, workspaceRoot string) error {
	if workspaceRoot == "" || filepath.Clean(workspaceRoot) == "/" {
		return fmt.Errorf("unsafe workspace root")
	}
	entries, readErr := os.ReadDir(workspaceRoot)
	if readErr != nil {
		return readErr
	}
	for _, entry := range entries {
		if entry.Name() == pgrepo.GridDirName {
			continue
		}
		// Intent: Alice's fixture workspace should match the local snapshot state
		// while preserving `.grid` repo-control files. Removing only non-control
		// entries prevents stale initial fixture labels from appearing as
		// untracked files after the final scenario head is recorded. Source:
		// DI-bamum
		if removeErr := os.RemoveAll(filepath.Join(workspaceRoot, entry.Name())); removeErr != nil {
			return removeErr
		}
	}
	return workspace.MaterializeSnapshot(cas, snapshotCID, workspaceRoot)
}

func initBareMainRepository(path string) error {
	repository, initErr := git.PlainInit(path, true)
	if initErr != nil {
		return initErr
	}
	return repository.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.NewBranchReferenceName("main")))
}

func resetRunRoot(runRoot string) error {
	if runRoot == "" || filepath.Clean(runRoot) == "/" {
		return fmt.Errorf("unsafe run root")
	}
	if err := os.RemoveAll(runRoot); err != nil {
		return err
	}
	return os.MkdirAll(runRoot, 0o755)
}

func createFixture(root string) error {
	if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
		return err
	}
	readme := []byte("PromiseGrid POC18\n\nThis file is chunked, stored in CAS, and materialized from a snapshot.\n")
	if err := os.WriteFile(filepath.Join(root, "README.md"), readme, 0o644); err != nil {
		return err
	}
	large := make([]byte, 96*1024)
	for index := range large {
		large[index] = byte('a' + index%26)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "large.bin"), large, 0o644); err != nil {
		return err
	}
	if err := os.Link(filepath.Join(root, "README.md"), filepath.Join(root, "README-hardlink.md")); err != nil {
		return err
	}
	if err := os.Symlink("README.md", filepath.Join(root, "README-link.md")); err != nil {
		return err
	}
	if err := syscall.Mkfifo(filepath.Join(root, "build.pipe"), 0o644); err != nil {
		return err
	}
	return nil
}

func writeJSON(path string, value any) error {
	content, marshalErr := json.MarshalIndent(value, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	content = append(content, '\n')
	return os.WriteFile(path, content, 0o644)
}
