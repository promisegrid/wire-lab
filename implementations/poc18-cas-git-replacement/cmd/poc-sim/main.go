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

	cidlib "github.com/ipfs/go-cid"

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
	Retrieval      pocsync.RetrievalReport `json:"retrieval"`
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
	aliceCAS, aliceOpenErr := store.Open(aliceStoreRoot)
	if aliceOpenErr != nil {
		return aliceOpenErr
	}
	result, ingestErr := workspace.NewScanner(aliceCAS, "alice", "bob").Ingest(sourceRoot)
	if ingestErr != nil {
		return ingestErr
	}
	bobCAS, bobOpenErr := store.Open(bobStoreRoot)
	if bobOpenErr != nil {
		return bobOpenErr
	}
	branchCID, branchErr := store.ParseCIDText(result.BranchRefSetCID)
	if branchErr != nil {
		return branchErr
	}
	retrieval, retrieveErr := pocsync.RetrieveGraph(
		pocsync.Peer{Agent: "bob", CAS: bobCAS},
		pocsync.Peer{Agent: "alice", CAS: aliceCAS},
		map[string]cidlib.Cid{"branch": branchCID},
		"storage_credit:2",
	)
	if retrieveErr != nil {
		return retrieveErr
	}
	snapshotCID, parseErr := store.ParseCIDText(result.SnapshotCID)
	if parseErr != nil {
		return parseErr
	}
	if err := workspace.MaterializeSnapshot(bobCAS, snapshotCID, checkoutRoot); err != nil {
		return err
	}
	result.CheckoutRoot = checkoutRoot
	combined := fixtureResult{
		IngestResult:   result,
		AliceStoreRoot: aliceStoreRoot,
		BobStoreRoot:   bobStoreRoot,
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
