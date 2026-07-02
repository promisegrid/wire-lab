package scenario

import (
	"os"
	"path/filepath"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/workspace"
)

// TestBuildProducesLineageMergeAndReviewPromises verifies the deterministic
// scenario slice without relying on the command fixture.
func TestBuildProducesLineageMergeAndReviewPromises(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	if err := os.MkdirAll(filepath.Join(source, "docs"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "README.md"), []byte("hello lineage\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(README) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "docs", "large.bin"), []byte("large enough for scenario\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(large) error = %v", err)
	}
	if err := os.Symlink("README.md", filepath.Join(source, "README-link.md")); err != nil {
		t.Fatalf("Symlink() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "build.pipe"), []byte("fixture pipe placeholder\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(build.pipe) error = %v", err)
	}
	cas, openErr := store.Open(filepath.Join(root, "cas"))
	if openErr != nil {
		t.Fatalf("Open() error = %v", openErr)
	}
	initial, ingestErr := workspace.NewScanner(cas, "alice", "bob").Ingest(source)
	if ingestErr != nil {
		t.Fatalf("Ingest() error = %v", ingestErr)
	}
	result, buildErr := NewBuilder(cas).Build(initial)
	if buildErr != nil {
		t.Fatalf("Build() error = %v", buildErr)
	}
	if result.LineageNodeCID == "" || result.LineageNodeCID != result.RenameLabelNodeCID || result.LineageNodeCID != result.CopyLabelNodeCID {
		t.Fatalf("lineage CIDs not preserved: %#v", result)
	}
	if len(result.MergeParentSnapshotCIDs) != 2 {
		t.Fatalf("merge parents = %#v, want two", result.MergeParentSnapshotCIDs)
	}
	if result.ReviewAdoptionResult != "accepted_locally" {
		t.Fatalf("adoption result = %q", result.ReviewAdoptionResult)
	}
	for _, cidText := range []string{
		result.RenameCopySnapshotCID,
		result.BobDivergentSnapshotCID,
		result.MergeSnapshotCID,
		result.TestStatementCID,
		result.AdoptionStatementCID,
		result.ReviewThreadRefSetCID,
	} {
		objectCID, parseErr := store.ParseCIDText(cidText)
		if parseErr != nil {
			t.Fatalf("ParseCIDText(%s) error = %v", cidText, parseErr)
		}
		if !cas.Has(objectCID) {
			t.Fatalf("CAS missing scenario object %s", cidText)
		}
	}
}
