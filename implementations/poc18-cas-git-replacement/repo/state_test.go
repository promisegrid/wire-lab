package repo

import (
	"os"
	"path/filepath"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

func TestInitCreatesDefaultState(t *testing.T) {
	root := t.TempDir()
	repository, initErr := Init(root, "")
	if initErr != nil {
		t.Fatalf("Init() error = %v", initErr)
	}
	statePath := filepath.Join(root, GridDirName, StateFileName)
	if _, statErr := os.Stat(statePath); statErr != nil {
		t.Fatalf("default state was not created: %v", statErr)
	}
	state, loadErr := repository.LoadState()
	if loadErr != nil {
		t.Fatalf("LoadState() error = %v", loadErr)
	}
	if state.Version != 1 || state.CurrentBranch != defaultCurrentBranch {
		t.Fatalf("state = %#v", state)
	}
}

func TestRecordSnapshotUpdatesState(t *testing.T) {
	root := t.TempDir()
	repository, initErr := Init(root, "")
	if initErr != nil {
		t.Fatalf("Init() error = %v", initErr)
	}
	snapshotCID := store.CIDText(store.CIDForBytes([]byte("snapshot")))
	branchCID := store.CIDText(store.CIDForBytes([]byte("branch")))
	workspaceCID := store.CIDText(store.CIDForBytes([]byte("workspace")))
	state, recordErr := repository.RecordSnapshot(snapshotCID, branchCID, workspaceCID)
	if recordErr != nil {
		t.Fatalf("RecordSnapshot() error = %v", recordErr)
	}
	if state.CurrentSnapshotCID != snapshotCID {
		t.Fatalf("CurrentSnapshotCID = %s, want %s", state.CurrentSnapshotCID, snapshotCID)
	}
	loaded, loadErr := repository.LoadState()
	if loadErr != nil {
		t.Fatalf("LoadState() error = %v", loadErr)
	}
	if loaded.CurrentBranchRefSetCID != branchCID || loaded.CurrentWorkspaceRefSetCID != workspaceCID {
		t.Fatalf("loaded state = %#v", loaded)
	}
}

func TestRecordSnapshotRejectsInvalidCID(t *testing.T) {
	root := t.TempDir()
	repository, initErr := Init(root, "")
	if initErr != nil {
		t.Fatalf("Init() error = %v", initErr)
	}
	if _, recordErr := repository.RecordSnapshot("not-a-cid", "", ""); recordErr == nil {
		t.Fatalf("RecordSnapshot() accepted invalid CID")
	}
}

func TestTrackAndUntrackPathsUpdateLocalExclusions(t *testing.T) {
	root := t.TempDir()
	repository, initErr := Init(root, "")
	if initErr != nil {
		t.Fatalf("Init() error = %v", initErr)
	}
	untracked, untrackErr := repository.UntrackPaths([]string{"build/output.log", "tmp"})
	if untrackErr != nil {
		t.Fatalf("UntrackPaths() error = %v", untrackErr)
	}
	if len(untracked.UntrackedPaths) != 2 || untracked.UntrackedPaths[0] != "build/output.log" || untracked.UntrackedPaths[1] != "tmp" {
		t.Fatalf("untracked paths = %#v", untracked.UntrackedPaths)
	}
	tracked, trackErr := repository.TrackPaths([]string{"tmp"})
	if trackErr != nil {
		t.Fatalf("TrackPaths() error = %v", trackErr)
	}
	if len(tracked.UntrackedPaths) != 1 || tracked.UntrackedPaths[0] != "build/output.log" {
		t.Fatalf("tracked state = %#v", tracked.UntrackedPaths)
	}
}

func TestTrackAndUntrackRejectControlAndEscapingPaths(t *testing.T) {
	root := t.TempDir()
	repository, initErr := Init(root, "")
	if initErr != nil {
		t.Fatalf("Init() error = %v", initErr)
	}
	for _, path := range []string{".", ".grid/config.json", "../outside"} {
		if _, untrackErr := repository.UntrackPaths([]string{path}); untrackErr == nil {
			t.Fatalf("UntrackPaths(%q) succeeded", path)
		}
	}
}
