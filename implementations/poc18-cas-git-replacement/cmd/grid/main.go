// Command grid is the first POC18 CLI surface over the shared core library.
//
// Intent: The CLI must not reimplement core graph behavior separately from
// deterministic fixtures. Source: DI-harih
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	cidlib "github.com/ipfs/go-cid"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/bridge"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	pgrepo "promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/repo"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	pocsync "promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/sync"
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
	case "status":
		return status(args[1:])
	case "log":
		return logHistory(args[1:])
	case "track":
		return updateTracking("track", args[1:])
	case "untrack":
		return updateTracking("untrack", args[1:])
	case "checkout":
		return checkout(args[1:])
	case "refs", "diag":
		return diag(args[1:])
	case "git":
		return gitBridge(args[1:])
	case "sync":
		return syncCommand(args[1:])
	case "help":
		usage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func usage() {
	fmt.Println("usage: grid <init|ingest|snapshot|status|log|track|untrack|checkout|refs|diag|git|sync> [options]")
	fmt.Println()
	fmt.Println("first-slice commands:")
	fmt.Println("  init     - create .grid/config.json and the configured sparse local CAS")
	fmt.Println("  ingest   - scan a workspace and store graph objects")
	fmt.Println("  snapshot - scan a workspace, store graph objects, and record local state")
	fmt.Println("  status   - compare the current workspace with the recorded snapshot")
	fmt.Println("  log      - walk snapshot parent links from the recorded snapshot")
	fmt.Println("  track    - re-include repo paths in future snapshot/status handling")
	fmt.Println("  untrack  - exclude repo paths from future snapshot/status handling")
	fmt.Println("  checkout - materialize a snapshot from local CAS")
	fmt.Println("  refs     - render one retained reference/message object")
	fmt.Println("  diag     - render exact CBOR from a CAS CID or file")
	fmt.Println("  git      - run conventional Git bridge adapters: import, export, pull, push")
	fmt.Println("  sync     - inspect or run local sync-agent scheduling")
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
	scanner := workspace.NewScanner(cas, *promiser, *promisee)
	if command == "snapshot" {
		excludedPaths, exclusionErr := repoLocalExcludedPaths(repository, foundRepository, selectedWorkspace)
		if exclusionErr != nil {
			return exclusionErr
		}
		scanner.WithExcludedPaths(excludedPaths)
	}
	result, ingestErr := scanner.Ingest(selectedWorkspace)
	if ingestErr != nil {
		return ingestErr
	}
	if command == "snapshot" && foundRepository {
		if recordErr := recordRepoSnapshot(repository, selectedWorkspace, result); recordErr != nil {
			return recordErr
		}
	}
	return writeIngestResult(result, *outPath)
}

func recordRepoSnapshot(repository pgrepo.Repository, selectedWorkspace string, result workspace.IngestResult) error {
	absWorkspace, absErr := filepath.Abs(selectedWorkspace)
	if absErr != nil {
		return absErr
	}
	if absWorkspace != repository.Root {
		return nil
	}
	// Intent: A repo-local `grid snapshot` advances only this repo's local
	// convenience state. It does not claim global branch authority or update
	// state for an explicitly selected external workspace. Source: DI-bikif
	_, recordErr := repository.RecordSnapshot(result.SnapshotCID, result.BranchRefSetCID, result.WorkspaceRefSetCID)
	return recordErr
}

func status(args []string) error {
	flags := flag.NewFlagSet("status", flag.ContinueOnError)
	storeRoot := flags.String("store", "", "CAS store root")
	workspaceRoot := flags.String("workspace", "", "workspace path to compare; defaults to discovered repo root")
	snapshotText := flags.String("snapshot", "", "snapshot CID; defaults to .grid/state.json current snapshot")
	outPath := flags.String("out", "", "optional JSON result path")
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
	selectedSnapshot, snapshotErr := snapshotDefault(*snapshotText, repository, foundRepository)
	if snapshotErr != nil {
		return snapshotErr
	}
	snapshotCID, parseErr := store.ParseCIDText(selectedSnapshot)
	if parseErr != nil {
		return parseErr
	}
	// Intent: `grid status` is a read-only local comparison against the recorded
	// snapshot; it must not create new promise objects just to answer whether the
	// workspace differs from the current local head. Source: DI-bikif
	excludedPaths, exclusionErr := repoLocalExcludedPaths(repository, foundRepository, selectedWorkspace)
	if exclusionErr != nil {
		return exclusionErr
	}
	report, compareErr := workspace.CompareSnapshotWithExcludedPaths(cas, snapshotCID, selectedWorkspace, excludedPaths)
	if compareErr != nil {
		return compareErr
	}
	return writeStatus(report, *outPath)
}

func repoLocalExcludedPaths(repository pgrepo.Repository, foundRepository bool, selectedWorkspace string) ([]string, error) {
	if !foundRepository {
		return nil, nil
	}
	sameWorkspace, pathErr := sameAbsPath(selectedWorkspace, repository.Root)
	if pathErr != nil {
		return nil, pathErr
	}
	if !sameWorkspace {
		return nil, nil
	}
	// Intent: Track/untrack exclusions are local repo convenience state. They
	// apply only when commands operate on this repo root, not when the user asks
	// to scan an unrelated workspace through an explicit override. Source: DI-jokav
	state, stateErr := repository.LoadState()
	if stateErr != nil {
		return nil, stateErr
	}
	return state.UntrackedPaths, nil
}

func snapshotDefault(snapshotText string, repository pgrepo.Repository, foundRepository bool) (string, error) {
	if snapshotText != "" {
		return snapshotText, nil
	}
	if !foundRepository {
		return "", fmt.Errorf("-snapshot is required when -store overrides repo discovery")
	}
	state, stateErr := repository.LoadState()
	if stateErr != nil {
		return "", stateErr
	}
	if state.CurrentSnapshotCID == "" {
		return "", fmt.Errorf("no current snapshot recorded; run grid snapshot or pass -snapshot")
	}
	return state.CurrentSnapshotCID, nil
}

func writeStatus(report workspace.StatusReport, outPath string) error {
	if outPath != "" {
		if err := writeJSON(outPath, report); err != nil {
			return err
		}
	}
	// Intent: Surface content drift and tracked-status drift as separate facts
	// in human output, matching the JSON report shape used by tests and future
	// porcelain. Source: DI-tuhoj
	if report.Clean {
		fmt.Printf("clean snapshot=%s workspace=%s content_diff=%t tracked_status_diff=%t\n", report.SnapshotCID, report.SourceRoot, report.ContentDiff, report.TrackedStatusDiff)
		return nil
	}
	fmt.Printf("changed snapshot=%s workspace=%s entries=%d content_diff=%t tracked_status_diff=%t\n", report.SnapshotCID, report.SourceRoot, len(report.Entries), report.ContentDiff, report.TrackedStatusDiff)
	for _, entry := range report.Entries {
		fmt.Printf("%s\t%s\t%s\tcontent_diff=%t\ttracked_status_diff=%t\n", entry.Status, entry.Type, entry.Path, entry.ContentDiff, entry.TrackedStatusDiff)
	}
	return nil
}

func updateTracking(command string, args []string) error {
	flags := flag.NewFlagSet(command, flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() == 0 {
		return fmt.Errorf("%s requires at least one path", command)
	}
	repository, discoverErr := pgrepo.Discover(".")
	if discoverErr != nil {
		return discoverErr
	}
	paths, pathErr := repoRelativeCLIPaths(repository, flags.Args())
	if pathErr != nil {
		return pathErr
	}
	var state pgrepo.State
	var updateErr error
	verb := "tracked"
	switch command {
	case "track":
		state, updateErr = repository.TrackPaths(paths)
	case "untrack":
		verb = "untracked"
		state, updateErr = repository.UntrackPaths(paths)
	default:
		return fmt.Errorf("unknown tracking command %q", command)
	}
	if updateErr != nil {
		return updateErr
	}
	for _, path := range paths {
		fmt.Printf("%s %s\n", verb, path)
	}
	fmt.Printf("state=%s untracked_paths=%d\n", filepath.Join(repository.GridDir, pgrepo.StateFileName), len(state.UntrackedPaths))
	return nil
}

func repoRelativeCLIPaths(repository pgrepo.Repository, inputs []string) ([]string, error) {
	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		return nil, cwdErr
	}
	paths := make([]string, 0, len(inputs))
	for _, input := range inputs {
		path, pathErr := repoRelativeCLIPath(repository, cwd, input)
		if pathErr != nil {
			return nil, pathErr
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func repoRelativeCLIPath(repository pgrepo.Repository, cwd string, input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("path is required")
	}
	target := input
	if !filepath.IsAbs(target) {
		target = filepath.Join(cwd, target)
	}
	absTarget, absErr := filepath.Abs(target)
	if absErr != nil {
		return "", absErr
	}
	relPath, relErr := filepath.Rel(repository.Root, absTarget)
	if relErr != nil {
		return "", relErr
	}
	slashPath := filepath.ToSlash(relPath)
	if filepath.IsAbs(relPath) || slashPath == ".." || strings.HasPrefix(slashPath, "../") {
		return "", fmt.Errorf("path must stay inside repo root %s: %s", repository.Root, input)
	}
	return pgrepo.NormalizeStatePath(relPath)
}

func sameAbsPath(left string, right string) (bool, error) {
	absLeft, leftErr := filepath.Abs(left)
	if leftErr != nil {
		return false, leftErr
	}
	absRight, rightErr := filepath.Abs(right)
	if rightErr != nil {
		return false, rightErr
	}
	return filepath.Clean(absLeft) == filepath.Clean(absRight), nil
}

type logReport struct {
	StartSnapshotCID string     `json:"start_snapshot_cid"`
	Entries          []logEntry `json:"entries"`
}

type logEntry struct {
	SnapshotCID         string   `json:"snapshot_cid"`
	Promiser            string   `json:"promiser"`
	Promisee            string   `json:"promisee"`
	Summary             string   `json:"summary"`
	RootReferenceSetCID string   `json:"root_reference_set_cid"`
	ParentSnapshotCIDs  []string `json:"parent_snapshot_cids"`
}

func logHistory(args []string) error {
	flags := flag.NewFlagSet("log", flag.ContinueOnError)
	storeRoot := flags.String("store", "", "CAS store root")
	snapshotText := flags.String("snapshot", "", "snapshot CID; defaults to .grid/state.json current snapshot")
	limit := flags.Int("limit", 20, "maximum snapshots to print; 0 means unlimited")
	outPath := flags.String("out", "", "optional JSON result path")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *limit < 0 {
		return fmt.Errorf("-limit cannot be negative")
	}
	cas, repository, foundRepository, openErr := openCLIStore(*storeRoot)
	if openErr != nil {
		return openErr
	}
	selectedSnapshot, snapshotErr := snapshotDefault(*snapshotText, repository, foundRepository)
	if snapshotErr != nil {
		return snapshotErr
	}
	snapshotCID, parseErr := store.ParseCIDText(selectedSnapshot)
	if parseErr != nil {
		return parseErr
	}
	// Intent: `grid log` walks exact snapshot parent links from local state or an
	// explicit snapshot override. This keeps history inspection local and
	// CID-grounded without adding a global branch/log authority. Source: DI-bikif
	report, logErr := snapshotLog(cas, snapshotCID, *limit)
	if logErr != nil {
		return logErr
	}
	return writeLog(report, *outPath)
}

func snapshotLog(cas *store.FileStore, startSnapshot cidlib.Cid, limit int) (logReport, error) {
	report := logReport{StartSnapshotCID: store.CIDText(startSnapshot)}
	queue := []cidlib.Cid{startSnapshot}
	visited := map[string]bool{}
	for len(queue) > 0 {
		if limit > 0 && len(report.Entries) >= limit {
			break
		}
		next := queue[0]
		queue = queue[1:]
		nextText := store.CIDText(next)
		if visited[nextText] {
			continue
		}
		visited[nextText] = true
		entry, parents, entryErr := snapshotLogEntry(cas, next)
		if entryErr != nil {
			return logReport{}, entryErr
		}
		report.Entries = append(report.Entries, entry)
		queue = append(queue, parents...)
	}
	return report, nil
}

func snapshotLogEntry(cas *store.FileStore, snapshotCID cidlib.Cid) (logEntry, []cidlib.Cid, error) {
	content, _, getErr := cas.Get(snapshotCID)
	if getErr != nil {
		return logEntry{}, nil, getErr
	}
	view, parseErr := graph.ParseEnvelope(content)
	if parseErr != nil {
		return logEntry{}, nil, parseErr
	}
	kind, kindErr := view.PayloadKind()
	if kindErr != nil {
		return logEntry{}, nil, kindErr
	}
	if kind != "snapshot" {
		return logEntry{}, nil, fmt.Errorf("expected snapshot message, got %s", kind)
	}
	body, bodyErr := view.PayloadBody()
	if bodyErr != nil {
		return logEntry{}, nil, bodyErr
	}
	if len(body) != 5 {
		return logEntry{}, nil, fmt.Errorf("snapshot body must have five slots")
	}
	rootRefCID, rootErr := store.CIDFromLinkTag(body[1])
	if rootErr != nil {
		return logEntry{}, nil, rootErr
	}
	summary, summaryOK := body[3].(string)
	if !summaryOK {
		return logEntry{}, nil, fmt.Errorf("snapshot summary must be text")
	}
	parents, parentErr := snapshotParentCIDs(view, body)
	if parentErr != nil {
		return logEntry{}, nil, parentErr
	}
	entry := logEntry{
		SnapshotCID:         store.CIDText(snapshotCID),
		Promiser:            fmt.Sprint(view.Payload[0]),
		Promisee:            fmt.Sprint(view.Payload[1]),
		Summary:             summary,
		RootReferenceSetCID: store.CIDText(rootRefCID),
		ParentSnapshotCIDs:  cidTexts(parents),
	}
	return entry, parents, nil
}

func snapshotParentCIDs(view graph.EnvelopeView, body []any) ([]cidlib.Cid, error) {
	parents := []cidlib.Cid{}
	seen := map[string]bool{}
	payloadParents, ok := body[2].([]any)
	if !ok {
		return nil, fmt.Errorf("snapshot parent list must be array")
	}
	for _, parentValue := range payloadParents {
		parentCID, cidErr := store.CIDFromLinkTag(parentValue)
		if cidErr != nil {
			return nil, cidErr
		}
		parentText := store.CIDText(parentCID)
		if !seen[parentText] {
			parents = append(parents, parentCID)
			seen[parentText] = true
		}
	}
	for _, parent := range view.Parents {
		if parent.Role != "previous_snapshot" {
			continue
		}
		parentText := store.CIDText(parent.CID)
		if !seen[parentText] {
			parents = append(parents, parent.CID)
			seen[parentText] = true
		}
	}
	return parents, nil
}

func cidTexts(cids []cidlib.Cid) []string {
	texts := make([]string, 0, len(cids))
	for _, value := range cids {
		texts = append(texts, store.CIDText(value))
	}
	return texts
}

func writeLog(report logReport, outPath string) error {
	if outPath != "" {
		if err := writeJSON(outPath, report); err != nil {
			return err
		}
	}
	for _, entry := range report.Entries {
		fmt.Printf("%s %s\n", entry.SnapshotCID, entry.Summary)
		fmt.Printf("  promiser=%s promisee=%s root_refset=%s\n", entry.Promiser, entry.Promisee, entry.RootReferenceSetCID)
		if len(entry.ParentSnapshotCIDs) > 0 {
			fmt.Printf("  parents=%v\n", entry.ParentSnapshotCIDs)
		}
	}
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

func syncCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("sync command is required: status or once")
	}
	switch args[0] {
	case "status":
		return syncStatus(args[1:])
	case "once":
		return syncOnce(args[1:])
	default:
		return fmt.Errorf("unknown sync command %q", args[0])
	}
}

func syncStatus(args []string) error {
	flags := flag.NewFlagSet("sync status", flag.ContinueOnError)
	agent := flags.String("agent", "alice", "local sync-agent name")
	outPath := flags.String("out", "", "optional JSON result path")
	if err := flags.Parse(args); err != nil {
		return err
	}
	repository, discoverErr := pgrepo.Discover(".")
	if discoverErr != nil {
		return discoverErr
	}
	// Intent: `grid sync status` may inspect an absent local sync-agent checkpoint
	// but must not create `.grid/sync` just to answer. Source: DI-fakop
	state, loadErr := pocsync.LoadAgentState(syncAgentStatePath(repository), *agent)
	if loadErr != nil {
		return loadErr
	}
	if *outPath != "" {
		if err := writeJSON(*outPath, state); err != nil {
			return err
		}
	}
	chosen := ""
	if state.LastReport != nil {
		chosen = state.LastReport.ChosenPeer
	}
	fmt.Printf("sync_agent_state=%s agent=%s trusted_peers=%d last_chosen=%s\n", syncAgentStatePath(repository), state.Agent, len(state.Trust), chosen)
	return nil
}

func syncOnce(args []string) error {
	flags := flag.NewFlagSet("sync once", flag.ContinueOnError)
	agent := flags.String("agent", "alice", "local sync-agent name")
	peerName := flags.String("peer", "", "candidate peer agent name")
	peerStore := flags.String("peer-store", "", "candidate peer CAS root")
	marketSignal := flags.Int("market-signal", 0, "fallback market signal when local graph evidence is absent")
	bearerTokenFile := flags.String("bearer-token-file", "", "optional bearer token bytes to redeem for sync capability")
	bearerIssuer := flags.String("bearer-issuer", "", "issuer expected for bearer token verification")
	paymentScope := flags.String("payment-scope", "", "expected bearer-token scope")
	paymentObjectCID := flags.String("payment-object", "", "expected bearer-token object CID; defaults to first advertised head")
	paymentValue := flags.Int64("payment-value", 1, "expected bearer-token value")
	paymentUnit := flags.String("payment-unit", "sync_credit", "expected bearer-token unit")
	capability := flags.String("capability", "storage", "capability requested from selected peer")
	rounds := flags.Int("rounds", 1, "continuous sync rounds to run")
	outPath := flags.String("out", "", "optional JSON report path")
	var heads headFlags
	flags.Var(&heads, "head", "advertised head in role=cid form; may repeat")
	if err := flags.Parse(args); err != nil {
		return err
	}
	cas, repository, foundRepository, openErr := openCLIStore("")
	if openErr != nil {
		return openErr
	}
	if !foundRepository {
		return fmt.Errorf("grid sync once requires a discovered repo")
	}
	statePath := syncAgentStatePath(repository)
	state, stateErr := pocsync.LoadAgentState(statePath, *agent)
	if stateErr != nil {
		return stateErr
	}
	local := pocsync.Peer{Agent: *agent, CAS: cas}
	candidates, candidatesErr := syncCandidatesFromFlags(*peerName, *peerStore, heads, *marketSignal, *bearerTokenFile, *bearerIssuer, *paymentScope, *paymentObjectCID, *paymentValue, *paymentUnit, *capability)
	if candidatesErr != nil {
		return candidatesErr
	}
	nextState, report, runErr := pocsync.RunScheduledSync(local, candidates, state, pocsync.SchedulerConfig{Rounds: *rounds, Now: time.Date(2026, 7, 4, 12, 0, 0, 0, time.UTC)})
	if runErr != nil {
		return runErr
	}
	if saveErr := pocsync.SaveAgentState(statePath, nextState); saveErr != nil {
		return saveErr
	}
	if *outPath != "" {
		if err := writeJSON(*outPath, report); err != nil {
			return err
		}
	}
	fmt.Printf("sync_once_state=%s chosen=%s decisions=%d redeemed=%d missing=%d\n", statePath, report.ChosenPeer, len(report.Decisions), len(report.CapabilityRedemptions), report.ContinuousSync.MissingObjects)
	return nil
}

type headFlags map[string]cidlib.Cid

func (flags *headFlags) String() string {
	if flags == nil || len(*flags) == 0 {
		return ""
	}
	roles := make([]string, 0, len(*flags))
	for role := range *flags {
		roles = append(roles, role)
	}
	sort.Strings(roles)
	parts := make([]string, 0, len(roles))
	for _, role := range roles {
		parts = append(parts, role+"="+store.CIDText((*flags)[role]))
	}
	return strings.Join(parts, ",")
}

func (flags *headFlags) Set(value string) error {
	if *flags == nil {
		*flags = headFlags{}
	}
	role, cidText, found := strings.Cut(value, "=")
	if !found || role == "" || cidText == "" {
		return fmt.Errorf("-head must be role=cid")
	}
	objectCID, parseErr := store.ParseCIDText(cidText)
	if parseErr != nil {
		return parseErr
	}
	(*flags)[role] = objectCID
	return nil
}

func syncCandidatesFromFlags(peerName string, peerStore string, heads headFlags, marketSignal int, bearerTokenFile string, bearerIssuer string, paymentScope string, paymentObjectCID string, paymentValue int64, paymentUnit string, capability string) ([]pocsync.CandidatePeer, error) {
	if peerName == "" && peerStore == "" && len(heads) == 0 {
		return nil, nil
	}
	if peerName == "" || peerStore == "" || len(heads) == 0 {
		return nil, fmt.Errorf("-peer, -peer-store, and at least one -head are required together")
	}
	peerCAS, peerErr := store.Open(peerStore)
	if peerErr != nil {
		return nil, peerErr
	}
	bearerBytes, bearerErr := readOptionalFile(bearerTokenFile)
	if bearerErr != nil {
		return nil, bearerErr
	}
	if len(bearerBytes) > 0 {
		if bearerIssuer == "" {
			return nil, fmt.Errorf("-bearer-issuer is required with -bearer-token-file")
		}
		if paymentScope == "" {
			paymentScope = "grid-sync:" + peerName
		}
		if paymentObjectCID == "" {
			paymentObjectCID = firstHeadCID(heads)
		}
	}
	return []pocsync.CandidatePeer{{
		Peer:             pocsync.Peer{Agent: peerName, CAS: peerCAS},
		Heads:            map[string]cidlib.Cid(heads),
		MarketSignal:     marketSignal,
		BearerTokenBytes: bearerBytes,
		BearerIssuer:     bearerIssuer,
		PaymentScope:     paymentScope,
		PaymentObjectCID: paymentObjectCID,
		PaymentValue:     paymentValue,
		PaymentUnit:      paymentUnit,
		Capability:       capability,
	}}, nil
}

func readOptionalFile(path string) ([]byte, error) {
	if path == "" {
		return nil, nil
	}
	file, openErr := os.Open(path)
	if openErr != nil {
		return nil, openErr
	}
	deferErr := file.Close
	content, readErr := io.ReadAll(file)
	if closeErr := deferErr(); closeErr != nil && readErr == nil {
		return nil, closeErr
	}
	return content, readErr
}

func firstHeadCID(heads headFlags) string {
	roles := make([]string, 0, len(heads))
	for role := range heads {
		roles = append(roles, role)
	}
	sort.Strings(roles)
	if len(roles) == 0 {
		return ""
	}
	return store.CIDText(heads[roles[0]])
}

func syncAgentStatePath(repository pgrepo.Repository) string {
	return filepath.Join(repository.GridDir, "sync", "state.json")
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
