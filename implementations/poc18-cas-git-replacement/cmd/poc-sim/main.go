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

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/workspace"
)

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
	storeRoot := filepath.Join(*runRoot, "alice-cas")
	checkoutRoot := filepath.Join(*runRoot, "bob-checkout")
	if err := createFixture(sourceRoot); err != nil {
		return err
	}
	cas, openErr := store.Open(storeRoot)
	if openErr != nil {
		return openErr
	}
	result, ingestErr := workspace.NewScanner(cas, "alice", "bob").Ingest(sourceRoot)
	if ingestErr != nil {
		return ingestErr
	}
	snapshotCID, parseErr := store.ParseCIDText(result.SnapshotCID)
	if parseErr != nil {
		return parseErr
	}
	if err := workspace.MaterializeSnapshot(cas, snapshotCID, checkoutRoot); err != nil {
		return err
	}
	result.CheckoutRoot = checkoutRoot
	resultPath := filepath.Join(*runRoot, "result.json")
	if err := writeJSON(resultPath, result); err != nil {
		return err
	}
	fmt.Printf("run_root=%s\n", *runRoot)
	fmt.Printf("store=%s\n", storeRoot)
	fmt.Printf("snapshot=%s\n", result.SnapshotCID)
	fmt.Printf("diagnostic_message=%s\n", result.DiagnosticMessageCID)
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
