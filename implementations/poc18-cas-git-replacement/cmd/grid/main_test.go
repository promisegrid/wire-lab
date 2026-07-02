package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

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
