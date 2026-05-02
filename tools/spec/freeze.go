package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// cmdFreeze mints a new pCID from the draft file, writes a snapshot file
// alongside, and appends a manifest entry. Per DI-011-20260429-184459, the
// snapshot file and the manifest entry are produced together; per
// DI-011-20260429-184501, both happen inside this single binary.
//
// The freeze ritual is deliberate: a human (or bot acting as human) decides
// the draft is freeze-worthy and runs `spec freeze <slug>`.
func cmdFreeze(slug string) error {
	if slug == "" {
		return fmt.Errorf("slug must be non-empty")
	}
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}
	draftPath := filepath.Join(repoRoot, SpecsDir, slug+"-draft.md")
	data, err := os.ReadFile(draftPath)
	if err != nil {
		return fmt.Errorf("read draft %q: %w", draftPath, err)
	}

	// Self-reference lint (DI-011-20260429-184456): drafts MUST NOT contain
	// a placeholder for their own pCID-to-be. Refuse to freeze if any known
	// placeholder pattern is present.
	if hits := scanSelfReferencePlaceholders(string(data)); len(hits) > 0 {
		return fmt.Errorf("draft contains self-reference placeholders that violate DI-011-20260429-184456: %s",
			strings.Join(hits, ", "))
	}

	pcid, err := computeCID(data)
	if err != nil {
		return err
	}

	// Snapshot file: literal copy of the draft bytes (DI-011-20260429-184457
	// says the bytes hashed are the bytes on disk; the snapshot file IS those
	// hashed bytes, byte-for-byte).
	snapPath := filepath.Join(repoRoot, SpecsDir, slug+"-"+pcid+".md")
	if _, err := os.Stat(snapPath); err == nil {
		// Already frozen at this exact pCID; nothing to do, but still append
		// a fresh manifest entry only if missing. Most likely this is a
		// re-run after a partial failure.
		fmt.Fprintf(os.Stderr, "spec freeze: snapshot %s already exists; checking manifest…\n", snapPath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("stat snapshot %q: %w", snapPath, err)
	} else {
		if err := os.WriteFile(snapPath, data, 0o644); err != nil {
			return fmt.Errorf("write snapshot %q: %w", snapPath, err)
		}
	}

	// Append manifest entry.
	m, err := readManifestOrEmpty(repoRoot)
	if err != nil {
		return err
	}
	if existing := m.findEntryByPCID(pcid); existing != nil {
		fmt.Fprintf(os.Stderr, "spec freeze: manifest already has entry for %s; nothing to do\n", pcid)
	} else {
		now := time.Now().UTC().Format("2006-01-02T15:04:05Z")
		entry := Entry{
			PCID:     pcid,
			Slug:     slug,
			Status:   StatusFrozen,
			FrozenOn: now,
		}
		// Link supersedes if there is a previous frozen entry for this slug.
		if prev := m.latestFrozen(slug); prev != nil {
			entry.Supersedes = prev.PCID
			prev.SupersededBy = pcid
			prev.Status = StatusSuperseded
		}
		// Try to record the freezing commit; if we're not in a git tree or
		// HEAD has no commits yet, leave it blank (the audit will allow this
		// since the field is optional).
		if commit := gitHeadCommit(repoRoot); commit != "" {
			entry.FreezingCommit = commit
		}
		m.Entries = append(m.Entries, entry)
		sortEntries(m.Entries)
		if err := writeManifest(repoRoot, m); err != nil {
			return fmt.Errorf("write manifest: %w", err)
		}
	}

	// Stage both files in git.
	gitStage(repoRoot, SpecsDir+"/"+slug+"-"+pcid+".md")
	gitStage(repoRoot, ManifestPath)

	fmt.Printf("Frozen: %s\n  pCID: %s\n  Snapshot: %s/%s-%s.md\n  Manifest: %s (entry appended)\n",
		slug, pcid, SpecsDir, slug, pcid, ManifestPath)
	return nil
}

// readManifestOrEmpty reads MANIFEST.md, or returns an empty Manifest if the
// file does not exist. The empty case allows `spec freeze` to bootstrap the
// manifest on a freshly-created repo. The empty Manifest is initialized with
// a default preamble explaining the file's role and the locked CIDv1
// parameters; that preamble is written on first save.
func readManifestOrEmpty(repoRoot string) (*Manifest, error) {
	m, err := readManifest(repoRoot)
	if err == nil {
		return m, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	return &Manifest{
		Preamble:  defaultPreamble,
		Entries:   nil,
		Postamble: defaultPostamble,
	}, nil
}

// defaultPreamble is the human-readable prose written above the YAML block on
// the first save of MANIFEST.md. Subsequent saves preserve whatever prose is
// already in place.
const defaultPreamble = `# Spec-doc manifest

This file is the authoritative index of frozen spec docs for the wire-lab
harness. Drafts live alongside frozen snapshots under
` + "`protocols/wire-lab.d/specs/`" + ` (flat layout per
DI-011-20260429-184454, anchored at the harness's protocol-as-simrepo
directory per TE-29 + TE-32). Each frozen spec doc is content-addressed by
its pCID (a CIDv1 of the file's literal bytes, per DI-011-20260429-184457).

CIDv1 parameter set (per DI-011-20260429-184453):

- multibase: ` + "`base32`" + `
- multihash: ` + "`sha2-256`" + `
- codec:     ` + "`raw`" + `

The single fenced YAML block below is the machine-readable record. It is
authoritative; humans read it via this file's Markdown rendering, machines
parse it via the ` + "`tools/spec`" + ` Go binary. Status values are
` + "`frozen`" + ` (the canonical name of one frozen pCID), ` + "`superseded`" + `
(an older frozen pCID whose successor has been frozen), or ` + "`draft-ahead`" + `
(reserved; not currently emitted by the freeze tool but recognized by the
audit if a future tool emits it).

`

const defaultPostamble = `
See ` + "`docs/thought-experiments/TE-20260429-175530-spec-doc-store-and-pcid-machinery.md`" + `
for the full reasoning behind the layout, hash input, manifest format, and
the single-binary freeze-and-check tool.
`

// scanSelfReferencePlaceholders looks for textual patterns that would violate
// the external-only self-reference rule (DI-011-20260429-184456). The list is
// conservative: only patterns that clearly mean "my own pCID goes here" are
// flagged. False positives in spec prose (e.g., a paragraph discussing the
// concept of self-reference) would be a problem, so we restrict to literal
// placeholder strings that no human would write in normative prose.
func scanSelfReferencePlaceholders(s string) []string {
	patterns := []string{
		"<self-pcid>",
		"<my-pcid>",
		"<own-pcid>",
		"pcid: TBD",
		"pcid: tbd",
	}
	var hits []string
	for _, p := range patterns {
		if strings.Contains(s, p) {
			hits = append(hits, p)
		}
	}
	return hits
}

// findRepoRoot walks up from the current working directory until it finds a
// directory containing a `.git` entry. The wire-lab harness's specs live
// under SpecsDir (= protocols/wire-lab.d/specs/) per TE-29 + TE-32, but
// the spec tool only needs the .git ancestor to anchor SpecsDir resolution;
// the specs directory itself may be absent in the bootstrap case.
//
// Used by all subcommands so they can be run from any subdirectory of the
// repo. Returns an absolute path.
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find repo root (no .git ancestor)")
		}
		dir = parent
	}
}

// gitHeadCommit returns the SHA of HEAD in the repo, or "" if it cannot be
// determined (e.g., empty repo, git not on PATH). This is best-effort; the
// freeze tool should not fail if this returns "".
func gitHeadCommit(repoRoot string) string {
	cmd := exec.Command("git", "-C", repoRoot, "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// gitStage runs `git add <path>` in the repo. Errors are ignored: staging is
// a convenience, not a correctness requirement. The freeze ritual still
// produces a correct working tree even if `git add` fails (the user can
// stage manually).
func gitStage(repoRoot, path string) {
	exec.Command("git", "-C", repoRoot, "add", path).Run()
}
