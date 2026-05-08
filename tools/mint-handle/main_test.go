package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMintWidthOne checks that proquint-1 mints are 5 chars and pass the
// proquint shape regex.
func TestMintWidthOne(t *testing.T) {
	got, err := mint(1, map[string]string{}, 0, false)
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	if len(got) != 5 {
		t.Fatalf("len(%q) = %d, want 5", got, len(got))
	}
	if !handleFileRE.MatchString("TE-" + got + "-x.md") {
		t.Fatalf("mint %q does not match handle regex", got)
	}
}

// TestMintWidthTwo checks that proquint-2 mints are 11 chars (5+1+5) and
// contain exactly one hyphen.
func TestMintWidthTwo(t *testing.T) {
	got, err := mint(2, map[string]string{}, 0, false)
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	if len(got) != 11 {
		t.Fatalf("len(%q) = %d, want 11", got, len(got))
	}
	if strings.Count(got, "-") != 1 {
		t.Fatalf("hyphen count in %q = %d, want 1", got, strings.Count(got, "-"))
	}
}

// TestMintAvoidsCollision pre-populates the corpus with the value that
// seed=42 would produce on first attempt, then mints with a clock-based
// retry. The returned handle must be different from the seeded one, proving
// that the retry path actually fires.
func TestMintAvoidsCollision(t *testing.T) {
	corpus := map[string]string{}
	first, err := mint(1, corpus, 42, false)
	if err != nil {
		t.Fatalf("seed mint: %v", err)
	}
	corpus[first] = "fake/path.md"
	got, err := mint(1, corpus, 42, false)
	if err != nil {
		t.Fatalf("mint after collision: %v", err)
	}
	if got == first {
		t.Fatalf("mint returned colliding handle %q", got)
	}
}

// TestMintDryRunCollidesErrors verifies dry-run is deterministic and reports
// a collision instead of retrying with wall-clock entropy.
func TestMintDryRunCollidesErrors(t *testing.T) {
	first, err := mint(1, map[string]string{}, 42, true)
	if err != nil {
		t.Fatalf("seed dry-run: %v", err)
	}
	_, err = mint(1, map[string]string{first: "fake/path.md"}, 42, true)
	if err == nil {
		t.Errorf("dry-run with seeded collision: want error, got nil")
	}
}

// scanDirsForTests is the combined root-layout-plus-protocol-layout set used
// by tests to pre-create representative directories.
var scanDirsForTests = []string{
	"docs/thought-experiments",
	"TODO",
	"DR",
	"protocols/wire-lab.d/TODO",
}

// TestScanCorpusEmpty verifies an empty repo scan returns an empty corpus
// instead of a nil-map error.
func TestScanCorpusEmpty(t *testing.T) {
	root := t.TempDir()
	createScanDirs(t, root)
	corpus, err := scanCorpus(root)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(corpus) != 0 {
		t.Fatalf("corpus size = %d, want 0", len(corpus))
	}
}

// TestScanCorpusAllowsMissingLayouts proves the scanner works when only one
// of the root or protocol TODO layouts exists.
func TestScanCorpusAllowsMissingLayouts(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "TODO"), 0o755); err != nil {
		t.Fatal(err)
	}
	corpus, err := scanCorpus(root)
	if err != nil {
		t.Fatalf("scan root-only layout: %v", err)
	}
	if len(corpus) != 0 {
		t.Errorf("root-only corpus size = %d, want 0", len(corpus))
	}
}

// TestScanCorpusFindsHandles seeds a temp repo with TE, TODO, DR, and DI
// owners and verifies all handles are recovered with relative owner paths.
func TestScanCorpusFindsHandles(t *testing.T) {
	root := t.TempDir()
	createScanDirs(t, root)
	writeFile(t, root, "docs/thought-experiments/TE-vapoj-substrate.md", "# placeholder\n")
	writeFile(t, root, "TODO/TODO-bahor-something.md", "# placeholder\n\nID: DI-nisam\n")
	writeFile(t, root, "DR/DR-fonuz-question.md", "# placeholder\n")
	writeFile(t, root, "protocols/wire-lab.d/TODO/TODO-lusab-protocol.md", "# placeholder\n")

	corpus, err := scanCorpus(root)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	for _, handle := range []string{"vapoj", "bahor", "nisam", "fonuz", "lusab"} {
		if _, ok := corpus[handle]; !ok {
			t.Errorf("missing handle %q in %v", handle, corpus)
		}
	}
}

// TestScanCorpusIgnoresLegacyFilenames verifies that pre-upgrade filenames
// are not parsed as handles.
func TestScanCorpusIgnoresLegacyFilenames(t *testing.T) {
	root := t.TempDir()
	createScanDirs(t, root)
	legacy := []string{
		"docs/thought-experiments/TE-20260427-180000-promise-stack-ordering.md",
		"TODO/005-te-promise-stack-ordering.md",
		"DR/DR-006-20260429-164729-promise-stack-ordering.md",
		"TODO/TODO.md",
	}
	for _, path := range legacy {
		writeFile(t, root, path, "# placeholder\n")
	}
	corpus, err := scanCorpus(root)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(corpus) != 0 {
		t.Fatalf("legacy corpus size = %d, want 0 (%v)", len(corpus), corpus)
	}
}

// TestScanCorpusDetectsDuplicateHandles guards against a future bug where
// two owners end up with the same handle on disk.
func TestScanCorpusDetectsDuplicateHandles(t *testing.T) {
	root := t.TempDir()
	createScanDirs(t, root)
	writeFile(t, root, "docs/thought-experiments/TE-vapoj-something.md", "# placeholder\n")
	writeFile(t, root, "TODO/TODO-vapoj-something-else.md", "# placeholder\n")
	if _, err := scanCorpus(root); err == nil {
		t.Errorf("duplicate filename handle: want error, got nil")
	}
}

// TestScanCorpusDetectsDIDuplicateHandles proves DI owner lines participate
// in the same global namespace as filename-owned handles.
func TestScanCorpusDetectsDIDuplicateHandles(t *testing.T) {
	root := t.TempDir()
	createScanDirs(t, root)
	writeFile(t, root, "DR/DR-vapoj-question.md", "# placeholder\n")
	writeFile(t, root, "TODO/TODO-bahor-something.md", "# placeholder\n\nID: DI-vapoj\n")
	if _, err := scanCorpus(root); err == nil {
		t.Errorf("duplicate DI handle: want error, got nil")
	}
}

func createScanDirs(t *testing.T, root string) {
	t.Helper()
	for _, dir := range scanDirsForTests {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
}

func writeFile(t *testing.T, root, path, contents string) {
	t.Helper()
	full := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
