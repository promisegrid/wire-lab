package workspace

import (
	"os"
	"path/filepath"
	"testing"

	cidlib "github.com/ipfs/go-cid"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

func TestCompareSnapshotReportsCleanWorkspace(t *testing.T) {
	root, cas, snapshotCID := createStatusFixture(t)
	report, compareErr := CompareSnapshot(cas, snapshotCID, root)
	if compareErr != nil {
		t.Fatalf("CompareSnapshot() error = %v", compareErr)
	}
	if !report.Clean || len(report.Entries) != 0 {
		t.Fatalf("report = %#v, want clean", report)
	}
}

func TestCompareSnapshotReportsModifiedMissingAndUntracked(t *testing.T) {
	root, cas, snapshotCID := createStatusFixture(t)
	if writeErr := os.WriteFile(filepath.Join(root, "README.md"), []byte("changed\n"), 0o644); writeErr != nil {
		t.Fatalf("WriteFile(README) error = %v", writeErr)
	}
	if removeErr := os.Remove(filepath.Join(root, "docs", "guide.md")); removeErr != nil {
		t.Fatalf("Remove(guide) error = %v", removeErr)
	}
	if writeErr := os.WriteFile(filepath.Join(root, "extra.txt"), []byte("untracked\n"), 0o644); writeErr != nil {
		t.Fatalf("WriteFile(extra) error = %v", writeErr)
	}
	report, compareErr := CompareSnapshot(cas, snapshotCID, root)
	if compareErr != nil {
		t.Fatalf("CompareSnapshot() error = %v", compareErr)
	}
	assertStatusEntry(t, report, "README.md", StatusModified)
	assertStatusEntry(t, report, "docs/guide.md", StatusMissing)
	assertStatusEntry(t, report, "extra.txt", StatusTrackingAdded)
	assertStatusFlags(t, report, true, true)
}

func TestCompareSnapshotReportsTypeChanged(t *testing.T) {
	root, cas, snapshotCID := createStatusFixture(t)
	readmePath := filepath.Join(root, "README.md")
	if removeErr := os.Remove(readmePath); removeErr != nil {
		t.Fatalf("Remove(README) error = %v", removeErr)
	}
	if mkdirErr := os.Mkdir(readmePath, 0o755); mkdirErr != nil {
		t.Fatalf("Mkdir(README) error = %v", mkdirErr)
	}
	report, compareErr := CompareSnapshot(cas, snapshotCID, root)
	if compareErr != nil {
		t.Fatalf("CompareSnapshot() error = %v", compareErr)
	}
	assertStatusEntry(t, report, "README.md", StatusTypeChanged)
}

func TestCompareSnapshotWithExcludedPathsIgnoresExcludedLocalOnlyPath(t *testing.T) {
	root, cas, snapshotCID := createStatusFixture(t)
	if writeErr := os.WriteFile(filepath.Join(root, "local.log"), []byte("local-only\n"), 0o644); writeErr != nil {
		t.Fatalf("WriteFile(local.log) error = %v", writeErr)
	}
	report, compareErr := CompareSnapshotWithExcludedPaths(cas, snapshotCID, root, []string{"local.log"})
	if compareErr != nil {
		t.Fatalf("CompareSnapshotWithExcludedPaths() error = %v", compareErr)
	}
	if !report.Clean {
		t.Fatalf("report = %#v, want clean after exclusions", report)
	}
}

func TestCompareSnapshotWithExcludedPathsReportsTrackingRemoved(t *testing.T) {
	root, cas, snapshotCID := createStatusFixture(t)
	report, compareErr := CompareSnapshotWithExcludedPaths(cas, snapshotCID, root, []string{"docs"})
	if compareErr != nil {
		t.Fatalf("CompareSnapshotWithExcludedPaths() error = %v", compareErr)
	}
	assertStatusEntry(t, report, "docs", StatusTrackingRemoved)
	assertStatusFlags(t, report, false, true)
}

func createStatusFixture(t *testing.T) (string, *store.FileStore, cidlib.Cid) {
	t.Helper()
	root := filepath.Join(t.TempDir(), "workspace")
	if mkdirErr := os.MkdirAll(filepath.Join(root, "docs"), 0o755); mkdirErr != nil {
		t.Fatalf("MkdirAll(workspace) error = %v", mkdirErr)
	}
	if writeErr := os.WriteFile(filepath.Join(root, "README.md"), []byte("hello\n"), 0o644); writeErr != nil {
		t.Fatalf("WriteFile(README) error = %v", writeErr)
	}
	if writeErr := os.WriteFile(filepath.Join(root, "docs", "guide.md"), []byte("guide\n"), 0o644); writeErr != nil {
		t.Fatalf("WriteFile(guide) error = %v", writeErr)
	}
	cas, openErr := store.Open(filepath.Join(t.TempDir(), "cas"))
	if openErr != nil {
		t.Fatalf("store.Open() error = %v", openErr)
	}
	result, ingestErr := NewScanner(cas, "alice", "bob").Ingest(root)
	if ingestErr != nil {
		t.Fatalf("Ingest() error = %v", ingestErr)
	}
	snapshotCID, parseErr := store.ParseCIDText(result.SnapshotCID)
	if parseErr != nil {
		t.Fatalf("ParseCIDText(snapshot) error = %v", parseErr)
	}
	return root, cas, snapshotCID
}

func assertStatusEntry(t *testing.T, report StatusReport, path string, status string) {
	t.Helper()
	for _, entry := range report.Entries {
		if entry.Path == path && entry.Status == status {
			return
		}
	}
	t.Fatalf("status entry %s/%s not found in %#v", path, status, report.Entries)
}

func assertStatusFlags(t *testing.T, report StatusReport, contentDiff bool, trackedStatusDiff bool) {
	t.Helper()
	if report.ContentDiff != contentDiff || report.TrackedStatusDiff != trackedStatusDiff {
		t.Fatalf("report flags = content:%t tracking:%t, want content:%t tracking:%t", report.ContentDiff, report.TrackedStatusDiff, contentDiff, trackedStatusDiff)
	}
}
