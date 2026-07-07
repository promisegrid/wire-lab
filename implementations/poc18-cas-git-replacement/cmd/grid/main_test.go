package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	pgrepo "promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/repo"
	pocsync "promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/sync"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/workspace"
)

func TestGridInitAndSnapshotDiscoverGridConfig(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	if err := run([]string{"init"}); err != nil {
		t.Fatalf("run(init) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("repo-root snapshot\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(README) error = %v", err)
	}
	outPath := filepath.Join(root, "snapshot.json")
	if err := run([]string{"snapshot", "-out", outPath}); err != nil {
		t.Fatalf("run(snapshot) error = %v", err)
	}
	result := readIngestResult(t, outPath)
	if result.SourceRoot != root {
		t.Fatalf("SourceRoot = %s, want %s", result.SourceRoot, root)
	}
	if result.StoreRoot != filepath.Join(root, ".grid", "cas") {
		t.Fatalf("StoreRoot = %s", result.StoreRoot)
	}
	if _, statErr := os.Stat(filepath.Join(root, ".grid", "cas", "objects")); statErr != nil {
		t.Fatalf("CAS objects dir missing: %v", statErr)
	}
}

func TestGridSnapshotDiscoversConfigFromSubdirectory(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	if err := run([]string{"init"}); err != nil {
		t.Fatalf("run(init) error = %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "docs", "api"), 0o755); err != nil {
		t.Fatalf("MkdirAll(docs) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("from nested cwd\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(guide) error = %v", err)
	}
	nested := filepath.Join(root, "docs", "api")
	t.Chdir(nested)
	outPath := filepath.Join(root, "nested-snapshot.json")
	if err := run([]string{"snapshot", "-out", outPath}); err != nil {
		t.Fatalf("run(snapshot nested) error = %v", err)
	}
	result := readIngestResult(t, outPath)
	if result.SourceRoot != root {
		t.Fatalf("SourceRoot from nested cwd = %s, want %s", result.SourceRoot, root)
	}
}

func TestGridStoreOverrideStillRequiresWorkspace(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	storeRoot := filepath.Join(root, "manual-cas")
	if err := run([]string{"snapshot", "-store", storeRoot, "-out", filepath.Join(root, "snapshot.json")}); err == nil {
		t.Fatalf("snapshot with -store and no -workspace succeeded")
	}
	source := filepath.Join(root, "source")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatalf("MkdirAll(source) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "README.md"), []byte("manual store\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(source README) error = %v", err)
	}
	outPath := filepath.Join(root, "manual-snapshot.json")
	if err := run([]string{"snapshot", "-store", storeRoot, "-workspace", source, "-out", outPath}); err != nil {
		t.Fatalf("snapshot with explicit store/workspace error = %v", err)
	}
	result := readIngestResult(t, outPath)
	if result.StoreRoot != storeRoot {
		t.Fatalf("StoreRoot = %s, want %s", result.StoreRoot, storeRoot)
	}
}

func TestGridSnapshotRecordsStateStatusAndLog(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	if err := run([]string{"init"}); err != nil {
		t.Fatalf("run(init) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("state-backed commands\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(README) error = %v", err)
	}
	outputRoot := t.TempDir()
	if err := run([]string{"snapshot", "-out", filepath.Join(outputRoot, "snapshot.json")}); err != nil {
		t.Fatalf("run(snapshot) error = %v", err)
	}
	state := readRepoState(t, filepath.Join(root, ".grid", "state.json"))
	if state.CurrentSnapshotCID == "" || state.CurrentBranchRefSetCID == "" || state.CurrentWorkspaceRefSetCID == "" {
		t.Fatalf("state after snapshot = %#v", state)
	}
	statusOut := filepath.Join(outputRoot, "status.json")
	if err := run([]string{"status", "-out", statusOut}); err != nil {
		t.Fatalf("run(status clean) error = %v", err)
	}
	statusReport := readStatusReport(t, statusOut)
	if !statusReport.Clean {
		t.Fatalf("status report = %#v, want clean", statusReport)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("changed\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(README changed) error = %v", err)
	}
	changedOut := filepath.Join(outputRoot, "status-changed.json")
	if err := run([]string{"status", "-out", changedOut}); err != nil {
		t.Fatalf("run(status changed) error = %v", err)
	}
	changedReport := readStatusReport(t, changedOut)
	assertCLIStatusEntry(t, changedReport, "README.md", workspace.StatusModified)
	logOut := filepath.Join(outputRoot, "log.json")
	if err := run([]string{"log", "-out", logOut}); err != nil {
		t.Fatalf("run(log) error = %v", err)
	}
	logged := readLogReport(t, logOut)
	if len(logged.Entries) != 1 || logged.Entries[0].SnapshotCID != state.CurrentSnapshotCID {
		t.Fatalf("log report = %#v, want current snapshot %s", logged, state.CurrentSnapshotCID)
	}
}

func TestGridStatusRequiresRecordedSnapshot(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	if err := run([]string{"init"}); err != nil {
		t.Fatalf("run(init) error = %v", err)
	}
	if err := run([]string{"status"}); err == nil {
		t.Fatalf("status without recorded snapshot succeeded")
	}
}

func TestGridTrackAndUntrackPathPolicy(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	if err := run([]string{"init"}); err != nil {
		t.Fatalf("run(init) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("tracked\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(README) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "local.log"), []byte("local-only\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(local.log) error = %v", err)
	}
	outputRoot := t.TempDir()
	if err := run([]string{"snapshot", "-out", filepath.Join(outputRoot, "snapshot-initial.json")}); err != nil {
		t.Fatalf("run(initial snapshot) error = %v", err)
	}
	if err := run([]string{"untrack", "local.log"}); err != nil {
		t.Fatalf("run(untrack) error = %v", err)
	}
	state := readRepoState(t, filepath.Join(root, ".grid", "state.json"))
	if len(state.UntrackedPaths) != 1 || state.UntrackedPaths[0] != "local.log" {
		t.Fatalf("state after untrack = %#v", state.UntrackedPaths)
	}
	untrackedOut := filepath.Join(outputRoot, "status-untracked.json")
	if err := run([]string{"status", "-out", untrackedOut}); err != nil {
		t.Fatalf("run(status after untrack) error = %v", err)
	}
	untrackedReport := readStatusReport(t, untrackedOut)
	assertCLIStatusEntry(t, untrackedReport, "local.log", workspace.StatusTrackingRemoved)
	assertCLIStatusFlags(t, untrackedReport, false, true)
	if err := run([]string{"snapshot", "-out", filepath.Join(outputRoot, "snapshot.json")}); err != nil {
		t.Fatalf("run(snapshot) error = %v", err)
	}
	statusOut := filepath.Join(outputRoot, "status.json")
	if err := run([]string{"status", "-out", statusOut}); err != nil {
		t.Fatalf("run(status clean) error = %v", err)
	}
	statusReport := readStatusReport(t, statusOut)
	if !statusReport.Clean {
		t.Fatalf("status with untracked local file = %#v, want clean", statusReport)
	}
	if err := run([]string{"track", "local.log"}); err != nil {
		t.Fatalf("run(track) error = %v", err)
	}
	trackedOut := filepath.Join(outputRoot, "status-tracked.json")
	if err := run([]string{"status", "-out", trackedOut}); err != nil {
		t.Fatalf("run(status after track) error = %v", err)
	}
	trackedReport := readStatusReport(t, trackedOut)
	assertCLIStatusEntry(t, trackedReport, "local.log", workspace.StatusTrackingAdded)
	assertCLIStatusFlags(t, trackedReport, false, true)
}

func TestGridSyncStateStartsOnlyFromSyncCommands(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	if err := run([]string{"init"}); err != nil {
		t.Fatalf("run(init) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("sync agent boundary\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(README) error = %v", err)
	}
	outputRoot := t.TempDir()
	if err := run([]string{"snapshot", "-out", filepath.Join(outputRoot, "snapshot.json")}); err != nil {
		t.Fatalf("run(snapshot) error = %v", err)
	}
	if err := run([]string{"status", "-out", filepath.Join(outputRoot, "status.json")}); err != nil {
		t.Fatalf("run(status) error = %v", err)
	}
	syncStatePath := filepath.Join(root, ".grid", "sync", "state.json")
	if _, statErr := os.Stat(syncStatePath); !os.IsNotExist(statErr) {
		t.Fatalf("status created sync-agent state: %v", statErr)
	}
	if err := run([]string{"sync", "status", "-out", filepath.Join(outputRoot, "sync-status.json")}); err != nil {
		t.Fatalf("run(sync status) error = %v", err)
	}
	if _, statErr := os.Stat(syncStatePath); !os.IsNotExist(statErr) {
		t.Fatalf("sync status created sync-agent state: %v", statErr)
	}
	syncOnceOut := filepath.Join(outputRoot, "sync-once.json")
	if err := run([]string{"sync", "once", "-out", syncOnceOut}); err != nil {
		t.Fatalf("run(sync once) error = %v", err)
	}
	if _, statErr := os.Stat(syncStatePath); statErr != nil {
		t.Fatalf("sync once did not create sync-agent state: %v", statErr)
	}
	state := readSyncAgentState(t, syncStatePath)
	if state.Agent != "alice" || state.LastReport == nil || state.LastReport.IdleReason != "no candidate peers" {
		t.Fatalf("sync-agent state = %#v", state)
	}
}

func readIngestResult(t *testing.T, path string) workspace.IngestResult {
	t.Helper()
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("ReadFile(%s) error = %v", path, readErr)
	}
	var result workspace.IngestResult
	if err := json.Unmarshal(content, &result); err != nil {
		t.Fatalf("Unmarshal(%s) error = %v", path, err)
	}
	return result
}

func readRepoState(t *testing.T, path string) pgrepo.State {
	t.Helper()
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("ReadFile(%s) error = %v", path, readErr)
	}
	var state pgrepo.State
	if err := json.Unmarshal(content, &state); err != nil {
		t.Fatalf("Unmarshal(%s) error = %v", path, err)
	}
	return state
}

func readStatusReport(t *testing.T, path string) workspace.StatusReport {
	t.Helper()
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("ReadFile(%s) error = %v", path, readErr)
	}
	var report workspace.StatusReport
	if err := json.Unmarshal(content, &report); err != nil {
		t.Fatalf("Unmarshal(%s) error = %v", path, err)
	}
	return report
}

func readLogReport(t *testing.T, path string) logReport {
	t.Helper()
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("ReadFile(%s) error = %v", path, readErr)
	}
	var report logReport
	if err := json.Unmarshal(content, &report); err != nil {
		t.Fatalf("Unmarshal(%s) error = %v", path, err)
	}
	return report
}

func readSyncAgentState(t *testing.T, path string) pocsync.AgentState {
	t.Helper()
	content, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("ReadFile(%s) error = %v", path, readErr)
	}
	var state pocsync.AgentState
	if err := json.Unmarshal(content, &state); err != nil {
		t.Fatalf("Unmarshal(%s) error = %v", path, err)
	}
	return state
}

func assertCLIStatusEntry(t *testing.T, report workspace.StatusReport, path string, status string) {
	t.Helper()
	for _, entry := range report.Entries {
		if entry.Path == path && entry.Status == status {
			return
		}
	}
	t.Fatalf("status entry %s/%s not found in %#v", path, status, report.Entries)
}

func assertCLIStatusFlags(t *testing.T, report workspace.StatusReport, contentDiff bool, trackedStatusDiff bool) {
	t.Helper()
	if report.ContentDiff != contentDiff || report.TrackedStatusDiff != trackedStatusDiff {
		t.Fatalf("report flags = content:%t tracking:%t, want content:%t tracking:%t", report.ContentDiff, report.TrackedStatusDiff, contentDiff, trackedStatusDiff)
	}
}
