package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// cmdCheck runs the audit. Per DI-011-20260429-184501, this is the same
// binary as cmdFreeze; both share the manifest-parsing and CID-computing
// logic so they cannot disagree.
//
// Exit semantics: cmdCheck returns nil if the repo is clean, or an error
// summarizing all findings if anything is wrong. Advisory hints (e.g., CRLF
// in drafts) are printed but do not produce an error return.
//
// Audit invariants:
//
//  1. Every frozen file in `specs/<slug>-{cidv1}.md` has a matching manifest
//     entry, and the file's actual CIDv1 equals the filename's CIDv1.
//  2. Every manifest entry whose status is `frozen` or `superseded` has a
//     matching on-disk file at `specs/<slug>-{pCID}.md`.
//  3. Every cross-spec reference in any `specs/*-draft.md` cites a pCID that
//     is present in the manifest with status `frozen` or `superseded`.
//     A draft citing another draft is a failure.
//  4. No `specs/*-draft.md` contains a self-reference placeholder.
//  5. All `supersedes` and `superseded_by` links inside the manifest resolve
//     to other manifest entries.
//
// Advisory checks (warn, don't fail):
//
//  - CRLF line endings in `specs/*-draft.md` (DI-011-20260429-184457 trades
//    editor-style robustness for transparency, but a CRLF flip is still
//    almost always an accident, so we flag it).
//  - UTF-8 BOM at the start of any draft.
func cmdCheck() error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}
	specsDir := filepath.Join(repoRoot, "specs")

	// If specs/ doesn't exist, the audit is vacuously clean: there is nothing
	// to check. This handles the pre-genesis state where wire-lab has not yet
	// migrated harness-spec.md into specs/.
	if _, err := os.Stat(specsDir); os.IsNotExist(err) {
		fmt.Println("spec check: no specs/ directory; nothing to audit")
		return nil
	} else if err != nil {
		return err
	}

	m, err := readManifestOrEmpty(repoRoot)
	if err != nil {
		return err
	}

	var problems []string
	var warnings []string

	// Index manifest entries by pCID for fast lookup.
	byPCID := make(map[string]*Entry, len(m.Entries))
	for i := range m.Entries {
		byPCID[m.Entries[i].PCID] = &m.Entries[i]
	}

	// (1) frozen files on disk → manifest entries.
	frozenFiles, err := listFrozenFiles(specsDir)
	if err != nil {
		return err
	}
	for _, ff := range frozenFiles {
		actualCID, err := computeFileCID(ff.path)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: cannot compute CID: %v", ff.path, err))
			continue
		}
		if actualCID != ff.pcid {
			problems = append(problems, fmt.Sprintf(
				"%s: filename declares pCID %s but file content hashes to %s; the file has been modified after freeze, which violates immutability",
				ff.path, ff.pcid, actualCID))
		}
		entry, ok := byPCID[ff.pcid]
		if !ok {
			problems = append(problems, fmt.Sprintf("%s: no matching manifest entry for pCID %s", ff.path, ff.pcid))
			continue
		}
		if entry.Slug != ff.slug {
			problems = append(problems, fmt.Sprintf(
				"%s: filename declares slug %q but manifest entry for pCID %s has slug %q",
				ff.path, ff.slug, ff.pcid, entry.Slug))
		}
	}

	// (2) manifest entries → on-disk files.
	frozenByPCID := make(map[string]bool, len(frozenFiles))
	for _, ff := range frozenFiles {
		frozenByPCID[ff.pcid] = true
	}
	for _, e := range m.Entries {
		if e.Status == StatusFrozen || e.Status == StatusSuperseded {
			if !frozenByPCID[e.PCID] {
				problems = append(problems, fmt.Sprintf(
					"manifest: entry %s (slug %s) status %s expects file specs/%s-%s.md but file is missing",
					e.PCID, e.Slug, e.Status, e.Slug, e.PCID))
			}
		}
	}

	// (3) cross-references in drafts cite frozen pCIDs.
	// (4) drafts contain no self-reference placeholders.
	drafts, err := listDraftFiles(specsDir)
	if err != nil {
		return err
	}
	for _, d := range drafts {
		raw, err := os.ReadFile(d)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: cannot read: %v", d, err))
			continue
		}
		text := string(raw)

		// Self-reference placeholder lint.
		if hits := scanSelfReferencePlaceholders(text); len(hits) > 0 {
			problems = append(problems, fmt.Sprintf(
				"%s: contains self-reference placeholders that violate DI-011-20260429-184456: %s",
				d, strings.Join(hits, ", ")))
		}

		// Cross-reference lint.
		for _, ref := range extractSpecCrossRefs(text) {
			// A reference to specs/<slug>-draft.md is a draft-citing-draft
			// failure (DI-011-20260429-184455 forbids this).
			if strings.HasSuffix(ref, "-draft.md") {
				problems = append(problems, fmt.Sprintf(
					"%s: cites draft %q; drafts must cite frozen pCIDs only (DI-011-20260429-184455)",
					d, ref))
				continue
			}
			// Otherwise the reference must be `<slug>-<pCID>.md`. Verify the
			// pCID is in the manifest.
			pcid := pcidFromFilename(ref)
			if pcid == "" {
				problems = append(problems, fmt.Sprintf(
					"%s: cross-reference %q does not match expected `<slug>-<pCID>.md` form",
					d, ref))
				continue
			}
			entry, ok := byPCID[pcid]
			if !ok {
				problems = append(problems, fmt.Sprintf(
					"%s: cross-reference cites pCID %s which is not in the manifest", d, pcid))
				continue
			}
			if entry.Status != StatusFrozen && entry.Status != StatusSuperseded {
				problems = append(problems, fmt.Sprintf(
					"%s: cross-reference cites pCID %s with status %s (only frozen/superseded are allowed)",
					d, pcid, entry.Status))
			}
		}

		// Advisory: CRLF.
		if strings.Contains(text, "\r\n") {
			warnings = append(warnings, fmt.Sprintf(
				"%s: contains CRLF line endings; raw-bytes hashing means CRLF flips change the pCID",
				d))
		}
		// Advisory: UTF-8 BOM.
		if strings.HasPrefix(text, "\ufeff") {
			warnings = append(warnings, fmt.Sprintf(
				"%s: starts with a UTF-8 BOM; raw-bytes hashing includes the BOM in the pCID",
				d))
		}
	}

	// (5) supersedes/superseded_by links resolve.
	for _, e := range m.Entries {
		if e.Supersedes != "" {
			if _, ok := byPCID[e.Supersedes]; !ok {
				problems = append(problems, fmt.Sprintf(
					"manifest: entry %s supersedes %s but no such manifest entry exists",
					e.PCID, e.Supersedes))
			}
		}
		if e.SupersededBy != "" {
			if _, ok := byPCID[e.SupersededBy]; !ok {
				problems = append(problems, fmt.Sprintf(
					"manifest: entry %s superseded_by %s but no such manifest entry exists",
					e.PCID, e.SupersededBy))
			}
		}
	}

	// Print warnings (always non-fatal).
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}

	if len(problems) > 0 {
		for _, p := range problems {
			fmt.Fprintf(os.Stderr, "FAIL: %s\n", p)
		}
		return fmt.Errorf("%d audit problem(s) found", len(problems))
	}

	fmt.Printf("spec check: OK (%d frozen file(s), %d manifest entry(ies), %d draft(s))\n",
		len(frozenFiles), len(m.Entries), len(drafts))
	return nil
}

// frozenFile names a discovered `specs/<slug>-<pcid>.md` snapshot.
type frozenFile struct {
	path string
	slug string
	pcid string
}

// frozenNameRegexp recognizes `<slug>-<pcid>.md` where pcid is a CIDv1 in
// base32 form: starts with 'b' followed by lowercase letters and digits 2-7.
// Slugs are kebab-case lowercase letters/digits/hyphens; the LAST hyphen
// before the pcid component is the boundary.
var frozenNameRegexp = regexp.MustCompile(`^([a-z0-9][a-z0-9-]*?)-(b[a-z2-7]{50,})\.md$`)

// listFrozenFiles enumerates `specs/*.md` excluding drafts and MANIFEST.md,
// returning those whose filenames match the frozen-snapshot pattern.
func listFrozenFiles(specsDir string) ([]frozenFile, error) {
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return nil, err
	}
	var out []frozenFile
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == "MANIFEST.md" {
			continue
		}
		if strings.HasSuffix(name, "-draft.md") {
			continue
		}
		m := frozenNameRegexp.FindStringSubmatch(name)
		if m == nil {
			continue
		}
		out = append(out, frozenFile{
			path: filepath.Join(specsDir, name),
			slug: m[1],
			pcid: m[2],
		})
	}
	return out, nil
}

// listDraftFiles enumerates `specs/*-draft.md`.
func listDraftFiles(specsDir string) ([]string, error) {
	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if !strings.HasSuffix(e.Name(), "-draft.md") {
			continue
		}
		out = append(out, filepath.Join(specsDir, e.Name()))
	}
	return out, nil
}

// extractSpecCrossRefs finds substrings of the form `<slug>-<something>.md`
// in the text where the surrounding context suggests a sibling-spec
// reference. We adopt a deliberately simple heuristic: any reference that
// looks like a relative path inside `specs/` (either bare like
// `harness-spec-bafk....md` or prefixed with `specs/` or `./` or `../`).
//
// This is a lint, not a parser, so false positives are acceptable as long as
// the canonical citation form (a relative-path Markdown link) is reliably
// detected.
func extractSpecCrossRefs(text string) []string {
	// Match Markdown link targets pointing at `*.md` files that look like
	// they live in `specs/`. We accept three shapes inside `(...)`:
	//
	//   specs/<name>.md
	//   ./<name>.md     (when the citing draft is also under specs/)
	//   <name>.md       (same)
	//
	// And we further require the bare filename to look like one of:
	//
	//   <slug>-draft.md           -> flagged by caller as a draft-citation
	//   <slug>-<base32-cid>.md    -> the legitimate frozen-pcid form
	//
	// Anything else is ignored (it could be a link to a non-spec doc).
	re := regexp.MustCompile(`\]\(((?:specs/|\./|\.\./|)([a-z0-9][a-z0-9-]*?)\.md)(?:#[^)]*)?\)`)
	hits := re.FindAllStringSubmatch(text, -1)
	var out []string
	for _, h := range hits {
		full := h[1]
		// Reduce to bare filename for caller (so caller's pcidFromFilename
		// and -draft.md suffix check work uniformly).
		base := full
		if i := strings.LastIndex(base, "/"); i >= 0 {
			base = base[i+1:]
		}
		// Only collect references that look like spec docs:
		// either `<slug>-draft.md` or `<slug>-<cid>.md`.
		if strings.HasSuffix(base, "-draft.md") {
			out = append(out, base)
		} else if frozenNameRegexp.MatchString(base) {
			out = append(out, base)
		}
	}
	return out
}

// pcidFromFilename pulls the trailing pcid out of a `<slug>-<pcid>.md` name,
// or returns "" if the name doesn't match.
func pcidFromFilename(name string) string {
	m := frozenNameRegexp.FindStringSubmatch(name)
	if m == nil {
		return ""
	}
	return m[2]
}
