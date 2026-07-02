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

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/bridge"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	pgrepo "promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/repo"
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
	case "git":
		return gitBridge(args[1:])
	case "help":
		usage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func usage() {
	fmt.Println("usage: grid <init|ingest|snapshot|checkout|refs|diag|git> [options]")
	fmt.Println()
	fmt.Println("first-slice commands:")
	fmt.Println("  init     - create .grid/config.json and the configured sparse local CAS")
	fmt.Println("  ingest   - scan a workspace and store graph objects")
	fmt.Println("  snapshot - first-slice alias for ingest while persistent workspace config is absent")
	fmt.Println("  checkout - materialize a snapshot from local CAS")
	fmt.Println("  refs     - render one retained reference/message object")
	fmt.Println("  diag     - render exact CBOR from a CAS CID or file")
	fmt.Println("  git      - run conventional Git bridge adapters: import, export, pull, push")
}

func initStore(args []string) error {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	rootPath := flags.String("root", ".", "repo root for .grid/config.json")
	storeRoot := flags.String("store", "", "legacy file CAS path to write into .grid/config.json")
	casPath := flags.String("cas-path", "", "file CAS path to write into .grid/config.json")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *storeRoot != "" && *casPath != "" {
		return fmt.Errorf("-store and -cas-path cannot both be set")
	}
	selectedCASPath := *casPath
	if selectedCASPath == "" {
		selectedCASPath = *storeRoot
	}
	repository, initErr := pgrepo.Init(*rootPath, selectedCASPath)
	if initErr != nil {
		return initErr
	}
	casPathText := repository.ResolvePath(repository.Config.CAS.Path)
	fmt.Printf("initialized repo=%s config=%s cas=%s\n", repository.Root, repository.ConfigPath, casPathText)
	return nil
}

func openCLIStore(storeRoot string) (*store.FileStore, pgrepo.Repository, bool, error) {
	if storeRoot != "" {
		cas, openErr := store.Open(storeRoot)
		return cas, pgrepo.Repository{}, false, openErr
	}
	// Intent: Normal grid commands discover repo-local control state from the
	// current directory upward, while `--store` remains an explicit override for
	// tests and unusual workflows. Source: DI-pahor
	repository, discoverErr := pgrepo.Discover(".")
	if discoverErr != nil {
		return nil, pgrepo.Repository{}, false, discoverErr
	}
	cas, openErr := repository.OpenFileCAS()
	if openErr != nil {
		return nil, pgrepo.Repository{}, false, openErr
	}
	return cas, repository, true, nil
}

func openGitBridgeStore(storeRoot string) (*store.FileStore, error) {
	cas, _, _, openErr := openCLIStore(storeRoot)
	return cas, openErr
}

func workspaceDefault(storeRoot string, workspaceRoot string, repository pgrepo.Repository, foundRepository bool) (string, error) {
	if workspaceRoot != "" {
		return workspaceRoot, nil
	}
	if foundRepository {
		return repository.Root, nil
	}
	if storeRoot != "" {
		return "", fmt.Errorf("-workspace is required when -store overrides repo discovery")
	}
	return "", fmt.Errorf("-workspace is required")
}

func writeIngestResult(result workspace.IngestResult, outPath string) error {
	if outPath != "" {
		if err := writeJSON(outPath, result); err != nil {
			return err
		}
	}
	fmt.Printf("snapshot=%s root_refset=%s store=%s\n", result.SnapshotCID, result.RootReferenceSetCID, result.StoreRoot)
	return nil
}

func ingest(command string, args []string) error {
	flags := flag.NewFlagSet(command, flag.ContinueOnError)
	storeRoot := flags.String("store", "", "CAS store root")
	workspaceRoot := flags.String("workspace", "", "workspace path to scan; defaults to discovered repo root")
	outPath := flags.String("out", "", "optional JSON result path")
	promiser := flags.String("promiser", "alice", "local promiser")
	promisee := flags.String("promisee", "bob", "local promisee")
	if err := flags.Parse(args); err != nil {
		return err
	}
	cas, repository, foundRepository, openErr := openCLIStore(*storeRoot)
	if openErr != nil {
		return openErr
	}
	selectedWorkspace, workspaceErr := workspaceDefault(*storeRoot, *workspaceRoot, repository, foundRepository)
	if workspaceErr != nil {
		return workspaceErr
	}
	result, ingestErr := workspace.NewScanner(cas, *promiser, *promisee).Ingest(selectedWorkspace)
	if ingestErr != nil {
		return ingestErr
	}
	return writeIngestResult(result, *outPath)
}

func checkout(args []string) error {
	flags := flag.NewFlagSet("checkout", flag.ContinueOnError)
	storeRoot := flags.String("store", "", "CAS store root")
	snapshotText := flags.String("snapshot", "", "snapshot CID")
	outPath := flags.String("out", "", "checkout destination")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *snapshotText == "" || *outPath == "" {
		return fmt.Errorf("-snapshot and -out are required")
	}
	cas, _, _, openErr := openCLIStore(*storeRoot)
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
	if cidText == "" {
		return nil, fmt.Errorf("either -file or -cid is required")
	}
	cas, _, _, openErr := openCLIStore(storeRoot)
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

func gitBridge(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("git bridge command is required: import, export, pull, or push")
	}
	switch args[0] {
	case "import":
		return gitImport(args[1:])
	case "export":
		return gitExport(args[1:])
	case "pull":
		return gitPull(args[1:])
	case "push":
		return gitPush(args[1:])
	default:
		return fmt.Errorf("unknown git bridge command %q", args[0])
	}
}

func gitImport(args []string) error {
	flags := flag.NewFlagSet("git import", flag.ContinueOnError)
	storeRoot := flags.String("store", "", "CAS store root")
	repositoryPath := flags.String("repo", "", "Git repository path to import")
	outPath := flags.String("out", "", "optional JSON result path")
	promiser := flags.String("promiser", "alice", "local bridge promiser")
	promisee := flags.String("promisee", "bob", "local bridge promisee")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *repositoryPath == "" {
		return fmt.Errorf("-repo is required")
	}
	cas, openErr := openGitBridgeStore(*storeRoot)
	if openErr != nil {
		return openErr
	}
	// Intent: `grid git import` is a compatibility adapter over the shared bridge
	// package, not a separate CLI implementation of Git-to-grid conversion.
	// Source: DI-fimap
	result, importErr := bridge.NewAdapter(cas, *promiser, *promisee).ImportRepository(*repositoryPath)
	if importErr != nil {
		return importErr
	}
	return writeBridgeResult(result, *outPath)
}

func gitExport(args []string) error {
	flags := flag.NewFlagSet("git export", flag.ContinueOnError)
	storeRoot := flags.String("store", "", "CAS store root")
	snapshotText := flags.String("snapshot", "", "snapshot CID to export")
	repositoryPath := flags.String("repo", "", "Git repository path to write")
	outPath := flags.String("out", "", "optional JSON result path")
	promiser := flags.String("promiser", "alice", "local bridge promiser")
	promisee := flags.String("promisee", "bob", "local bridge promisee")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *snapshotText == "" || *repositoryPath == "" {
		return fmt.Errorf("-snapshot and -repo are required")
	}
	cas, openErr := openGitBridgeStore(*storeRoot)
	if openErr != nil {
		return openErr
	}
	snapshotCID, parseErr := store.ParseCIDText(*snapshotText)
	if parseErr != nil {
		return parseErr
	}
	result, exportErr := bridge.NewAdapter(cas, *promiser, *promisee).ExportSnapshot(snapshotCID, *repositoryPath)
	if exportErr != nil {
		return exportErr
	}
	return writeBridgeResult(result, *outPath)
}

func gitPull(args []string) error {
	flags := flag.NewFlagSet("git pull", flag.ContinueOnError)
	storeRoot := flags.String("store", "", "CAS store root")
	remoteURL := flags.String("remote", "", "Git remote URL or local path")
	worktreePath := flags.String("worktree", "", "temporary local clone path")
	outPath := flags.String("out", "", "optional JSON result path")
	promiser := flags.String("promiser", "alice", "local bridge promiser")
	promisee := flags.String("promisee", "bob", "local bridge promisee")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *remoteURL == "" || *worktreePath == "" {
		return fmt.Errorf("-remote and -worktree are required")
	}
	cas, openErr := openGitBridgeStore(*storeRoot)
	if openErr != nil {
		return openErr
	}
	result, pullErr := bridge.NewAdapter(cas, *promiser, *promisee).PullRepository(*remoteURL, *worktreePath)
	if pullErr != nil {
		return pullErr
	}
	return writeBridgeResult(result, *outPath)
}

func gitPush(args []string) error {
	flags := flag.NewFlagSet("git push", flag.ContinueOnError)
	storeRoot := flags.String("store", "", "CAS store root")
	snapshotText := flags.String("snapshot", "", "snapshot CID to push")
	remoteURL := flags.String("remote", "", "Git remote URL or local path")
	worktreePath := flags.String("worktree", "", "temporary local export path")
	outPath := flags.String("out", "", "optional JSON result path")
	promiser := flags.String("promiser", "alice", "local bridge promiser")
	promisee := flags.String("promisee", "bob", "local bridge promisee")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *snapshotText == "" || *remoteURL == "" || *worktreePath == "" {
		return fmt.Errorf("-snapshot, -remote, and -worktree are required")
	}
	cas, openErr := openGitBridgeStore(*storeRoot)
	if openErr != nil {
		return openErr
	}
	snapshotCID, parseErr := store.ParseCIDText(*snapshotText)
	if parseErr != nil {
		return parseErr
	}
	result, pushErr := bridge.NewAdapter(cas, *promiser, *promisee).PushSnapshot(snapshotCID, *remoteURL, *worktreePath)
	if pushErr != nil {
		return pushErr
	}
	return writeBridgeResult(result, *outPath)
}

func writeBridgeResult(result bridge.Result, outPath string) error {
	if outPath != "" {
		if err := writeJSON(outPath, result); err != nil {
			return err
		}
	}
	fmt.Printf("git_%s head_snapshot=%s head_git=%s mapping=%s\n", result.Operation, result.HeadSnapshotCID, result.HeadGitHash, result.MappingMessageCID)
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
