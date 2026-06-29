package workspace

import (
	"os"
	"path/filepath"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

func TestIngestAndMaterializeSnapshot(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	if err := os.MkdirAll(filepath.Join(source, "docs"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "README.md"), []byte("hello promise grid\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := os.Symlink("README.md", filepath.Join(source, "README-link.md")); err != nil {
		t.Fatalf("Symlink() error = %v", err)
	}
	cas, openErr := store.Open(filepath.Join(root, "cas"))
	if openErr != nil {
		t.Fatalf("Open() error = %v", openErr)
	}
	result, ingestErr := NewScanner(cas, "alice", "bob").Ingest(source)
	if ingestErr != nil {
		t.Fatalf("Ingest() error = %v", ingestErr)
	}
	snapshotCID, parseErr := store.ParseCIDText(result.SnapshotCID)
	if parseErr != nil {
		t.Fatalf("ParseCIDText() error = %v", parseErr)
	}
	checkout := filepath.Join(root, "checkout")
	if err := MaterializeSnapshot(cas, snapshotCID, checkout); err != nil {
		t.Fatalf("MaterializeSnapshot() error = %v", err)
	}
	content, readErr := os.ReadFile(filepath.Join(checkout, "README.md"))
	if readErr != nil {
		t.Fatalf("ReadFile(checkout) error = %v", readErr)
	}
	if string(content) != "hello promise grid\n" {
		t.Fatalf("checkout content = %q", string(content))
	}
	if result.Counts["reference_set:workspace"] == 0 || result.Counts["reference_set:branch"] == 0 {
		t.Fatalf("missing workspace/branch reference-set counts: %#v", result.Counts)
	}
}
