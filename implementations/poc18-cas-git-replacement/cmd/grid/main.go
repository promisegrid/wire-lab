// Command grid is the first POC18 CLI surface over the shared core library.
//
// Intent: The CLI must not reimplement core graph behavior separately from
// deterministic fixtures. Source: DI-harih
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/workspace"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "grid: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		usage()
		return nil
	}
	switch args[0] {
	case "init":
		return initStore(args[1:])
	case "ingest", "snapshot":
		return ingest(args[0], args[1:])
	case "checkout":
		return checkout(args[1:])
	case "refs", "diag":
		return diag(args[1:])
	case "help":
		usage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func usage() {
	fmt.Println("usage: grid <init|ingest|snapshot|checkout|refs|diag> [options]")
	fmt.Println()
	fmt.Println("first-slice commands:")
	fmt.Println("  init     - create a sparse local CAS")
	fmt.Println("  ingest   - scan a workspace and store graph objects")
	fmt.Println("  snapshot - first-slice alias for ingest while persistent workspace config is absent")
	fmt.Println("  checkout - materialize a snapshot from local CAS")
	fmt.Println("  refs     - render one retained reference/message object")
	fmt.Println("  diag     - render exact CBOR from a CAS CID or file")
}

func initStore(args []string) error {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	storeRoot := flags.String("store", "", "CAS store root")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *storeRoot == "" {
		return fmt.Errorf("-store is required")
	}
	if _, err := store.Open(*storeRoot); err != nil {
		return err
	}
	fmt.Printf("initialized store=%s\n", *storeRoot)
	return nil
}

func ingest(command string, args []string) error {
	flags := flag.NewFlagSet(command, flag.ContinueOnError)
	storeRoot := flags.String("store", "", "CAS store root")
	workspaceRoot := flags.String("workspace", "", "workspace path to scan")
	outPath := flags.String("out", "", "optional JSON result path")
	promiser := flags.String("promiser", "alice", "local promiser")
	promisee := flags.String("promisee", "bob", "local promisee")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *storeRoot == "" || *workspaceRoot == "" {
		return fmt.Errorf("-store and -workspace are required")
	}
	cas, openErr := store.Open(*storeRoot)
	if openErr != nil {
		return openErr
	}
	result, ingestErr := workspace.NewScanner(cas, *promiser, *promisee).Ingest(*workspaceRoot)
	if ingestErr != nil {
		return ingestErr
	}
	if *outPath != "" {
		if err := writeJSON(*outPath, result); err != nil {
			return err
		}
	}
	fmt.Printf("snapshot=%s root_refset=%s store=%s\n", result.SnapshotCID, result.RootReferenceSetCID, result.StoreRoot)
	return nil
}

func checkout(args []string) error {
	flags := flag.NewFlagSet("checkout", flag.ContinueOnError)
	storeRoot := flags.String("store", "", "CAS store root")
	snapshotText := flags.String("snapshot", "", "snapshot CID")
	outPath := flags.String("out", "", "checkout destination")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *storeRoot == "" || *snapshotText == "" || *outPath == "" {
		return fmt.Errorf("-store, -snapshot, and -out are required")
	}
	cas, openErr := store.Open(*storeRoot)
	if openErr != nil {
		return openErr
	}
	snapshotCID, parseErr := store.ParseCIDText(*snapshotText)
	if parseErr != nil {
		return parseErr
	}
	if err := workspace.MaterializeSnapshot(cas, snapshotCID, *outPath); err != nil {
		return err
	}
	fmt.Printf("checkout=%s snapshot=%s\n", *outPath, *snapshotText)
	return nil
}

func diag(args []string) error {
	flags := flag.NewFlagSet("diag", flag.ContinueOnError)
	storeRoot := flags.String("store", "", "CAS store root")
	cidText := flags.String("cid", "", "CID to render from CAS")
	filePath := flags.String("file", "", "CBOR file to render")
	if err := flags.Parse(args); err != nil {
		return err
	}
	content, loadErr := loadDiagnosticBytes(*storeRoot, *cidText, *filePath)
	if loadErr != nil {
		return loadErr
	}
	diagnostic, diagErr := graph.Diagnostic(content)
	if diagErr != nil {
		return diagErr
	}
	fmt.Print(diagnostic)
	return nil
}

func loadDiagnosticBytes(storeRoot, cidText, filePath string) ([]byte, error) {
	if filePath != "" {
		return os.ReadFile(filePath)
	}
	if storeRoot == "" || cidText == "" {
		return nil, fmt.Errorf("either -file or both -store and -cid are required")
	}
	cas, openErr := store.Open(storeRoot)
	if openErr != nil {
		return nil, openErr
	}
	objectCID, parseErr := store.ParseCIDText(cidText)
	if parseErr != nil {
		return nil, parseErr
	}
	content, _, getErr := cas.Get(objectCID)
	return content, getErr
}

func writeJSON(path string, value any) error {
	content, marshalErr := json.MarshalIndent(value, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	content = append(content, '\n')
	return os.WriteFile(path, content, 0o644)
}
