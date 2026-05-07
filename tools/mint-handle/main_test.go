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
	corpus := map[string]string{}
	got, err := mint(1, corpus, 0, false)
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	if len(got) != 5 {
		t.Errorf("mint width=1 length = %d, want 5 (%q)", len(got), got)
	}
}

// TestMintWidthTwo checks that proquint-2 mints are 11 chars (5+1+5) and
// contain exactly one hyphen.
func TestMintWidthTwo(t *testing.T) {
	corpus := map[string]string{}
	got, err := mint(2, corpus, 0, false)
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	if len(got) != 11 {
		t.Errorf("mint width=2 length = %d, want 11 (%q)", len(got), got)
	}
	if strings.Count(got, "-") != 1 {
		t.Errorf("mint width=2 hyphens = %d, want 1 (%q)",
			strings.Count(got, "-"), got)
	}
}

// TestMintAvoidsCollision pre-populates the corpus with the value that
// seed=42 would produce on first attempt, then mints with a clock-based
// retry. The returned handle must be different from the seeded one,
// proving that the retry path actually fires.
func TestMintAvoidsCollision(t *testing.T) {
	corpus := map[string]string{}
	first, err := mint(1, corpus, 42, false)
	if err != nil {
		t.Fatalf("seed mint: %v", err)
	}
	corpus[first] = "fake/path.md"
	second, err := mint(1, corpus, 42, false)
	if err != nil {
		t.Fatalf("retry mint: %v", err)
	}
	if second == first {
		t.Errorf("retry returned same handle %q despite collision", second)
	}
}

// TestMintDryRunCollidesErrors checks that dry-run mode reports a clear
// error when its single attempt hits a collision.
func TestMintDryRunCollidesErrors(t *testing.T) {
	first, err := mint(1, map[string]string{}, 42, true)
	if err != nil {
		t.Fatalf("seed dry-run: %v", err)
	}
	corpus := map[string]string{first: "fake/path.md"}
	_, err = mint(1, corpus, 42, true)
	if err == nil {
		t.Errorf("dry-run with seeded collision: want error, got nil")
	}
}

// TestScanCorpusEmpty verifies an empty repo scan returns an empty corpus
// (not a nil-map error). This pins the "fresh repo" boundary case.
func TestScanCorpusEmpty(t *testing.T) {
	root := t.TempDir()
	for _, dir := range scanDirs {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	corpus, err := scanCorpus(root)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(corpus) != 0 {
		t.Errorf("empty repo corpus size = %d, want 0", len(corpus))
	}
}

// TestScanCorpusFindsHandles seeds a temp repo with one TE and one TODO
// and verifies both handles are recovered with their relative paths.
func TestScanCorpusFindsHandles(t *testing.T) {
	root := t.TempDir()
	for _, dir := range scanDirs {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	teName := "TE-vapoj-substrate-agnostic-layered-model.md"
	todoName := "TODO-bahor-something-else.md"
	if err := os.WriteFile(filepath.Join(root, "docs/thought-experiments", teName),
		[]byte("# placeholder\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "protocols/wire-lab.d/TODO", todoName),
		[]byte("# placeholder\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	corpus, err := scanCorpus(root)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if got, want := len(corpus), 2; got != want {
		t.Fatalf("corpus size = %d, want %d (%v)", got, want, corpus)
	}
	if _, ok := corpus["vapoj"]; !ok {
		t.Errorf("missing handle %q in %v", "vapoj", corpus)
	}
	if _, ok := corpus["bahor"]; !ok {
		t.Errorf("missing handle %q in %v", "bahor", corpus)
	}
}

// TestScanCorpusIgnoresLegacyFilenames verifies that pre-migration
// filenames (timestamped, integer-numbered) are not parsed as handles.
// These files exist in the repo today and must be invisible to the mint
// collision check until they are renamed.
func TestScanCorpusIgnoresLegacyFilenames(t *testing.T) {
	root := t.TempDir()
	for _, dir := range scanDirs {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	legacy := []string{
		"TE-20260427-180000-promise-stack-ordering.md",
		"TE-20260506-184800-substrate-agnostic-layered-model.md",
		"TODO-20260507-002306-te-39-wire-lab-devs-migration.md",
		"README.md",
		"TODO.md",
	}
	for _, name := range legacy {
		dir := "docs/thought-experiments"
		if strings.HasPrefix(name, "TODO") || name == "TODO.md" {
			dir = "protocols/wire-lab.d/TODO"
		}
		if name == "README.md" {
			dir = "docs/thought-experiments"
		}
		if err := os.WriteFile(filepath.Join(root, dir, name),
			[]byte("# placeholder\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	corpus, err := scanCorpus(root)
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if len(corpus) != 0 {
		t.Errorf("legacy-only corpus = %d, want 0 (%v)", len(corpus), corpus)
	}
}

// TestScanCorpusDetectsDuplicateHandles guards against a future bug where
// two files end up with the same handle on disk. scanCorpus must surface
// this immediately so mint() does not silently mint a third colliding
// handle.
func TestScanCorpusDetectsDuplicateHandles(t *testing.T) {
	root := t.TempDir()
	for _, dir := range scanDirs {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	a := "TE-vapoj-something.md"
	b := "TODO-vapoj-something-else.md"
	if err := os.WriteFile(filepath.Join(root, "docs/thought-experiments", a),
		[]byte("# placeholder\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "protocols/wire-lab.d/TODO", b),
		[]byte("# placeholder\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := scanCorpus(root); err == nil {
		t.Errorf("duplicate handle: want error, got nil")
	}
}
