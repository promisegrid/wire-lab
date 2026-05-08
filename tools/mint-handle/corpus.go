package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// scanLiteralDirs are the wire-lab directories scanned literally for
// existing handles. scanGlobs are glob patterns expanded against repoRoot;
// every directory matching any glob is also scanned. A file in any scanned
// directory whose name matches handleFileRE is treated as owning its
// handle. The corpus union of those handles is the collision check set for
// new mints.
//
// Convention-based rather than registry-based: there is no separate
// HANDLES.md or registry file. Filenames are the ground truth. To add a
// new artifact type that also bears handles, add its directory (or a
// matching glob) here.
//
// scanGlobs uses Go filepath.Glob semantics, which match path components
// individually. "protocols/*/TODO" therefore matches protocols/wire-lab.d/
// TODO, protocols/group-session.d/TODO, protocols/ppx-dr.d/TODO, etc., as
// new protocol subtrees are added.
//
// See TE-39 for the lock decision and TE-vapoj (substrate-agnostic
// layered model) for the broader rationale.
var scanLiteralDirs = []string{
	"docs/thought-experiments",
}

var scanGlobs = []string{
	"protocols/*/TODO",
}

// handleFileRE matches a wire-lab handle-bearing filename and captures the
// handle in group 1.
//
// Filename grammar:
//
//	<KIND>-<HANDLE>-<SLUG>.md
//
// where:
//
//	<KIND>   = TE | TODO
//	<HANDLE> = proquint-1 (5 chars, CVCVC, alphabets in proquint.go)
//	         | proquint-2 (5 + '-' + 5 = 11 chars)
//	<SLUG>   = lower-kebab-case description, may include digits
//
// The regex is anchored so it does not match files whose names happen to
// contain proquint-shaped substrings inside the slug.
var handleFileRE = regexp.MustCompile(
	`^(?:TE|TODO)-(` +
		// proquint-2 must come first so the alternation prefers it
		// over the proquint-1 prefix it begins with.
		`[bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz]` +
		`-` +
		`[bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz]` +
		`|` +
		`[bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz]` +
		`)-`)

// scanCorpus returns the set of handles currently in use across scanDirs,
// rooted at repoRoot. It is the collision check set for mint().
//
// Files whose names do not match handleFileRE are silently ignored. This
// includes README.md, TODO.md (the master cross-list), pre-migration files
// still using the integer/timestamp scheme, and any non-handle-bearing
// markdown.
//
// Returns an error if a scan directory does not exist; this surfaces typos
// in scanDirs rather than silently scanning nothing.
func scanCorpus(repoRoot string) (map[string]string, error) {
	handles := make(map[string]string)
	dirs := append([]string(nil), scanLiteralDirs...)
	for _, glob := range scanGlobs {
		matches, err := filepath.Glob(filepath.Join(repoRoot, glob))
		if err != nil {
			return nil, fmt.Errorf("glob %s: %w", glob, err)
		}
		for _, m := range matches {
			rel, err := filepath.Rel(repoRoot, m)
			if err != nil {
				return nil, fmt.Errorf("relpath %s: %w", m, err)
			}
			dirs = append(dirs, rel)
		}
	}
	for _, dir := range dirs {
		full := filepath.Join(repoRoot, dir)
		entries, err := os.ReadDir(full)
		if err != nil {
			return nil, fmt.Errorf("scan %s: %w", dir, err)
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if !strings.HasSuffix(name, ".md") {
				continue
			}
			m := handleFileRE.FindStringSubmatch(name)
			if m == nil {
				continue
			}
			h := m[1]
			if prev, dup := handles[h]; dup {
				return nil, fmt.Errorf(
					"corpus already contains duplicate handle %q in %s and %s",
					h, prev, filepath.Join(dir, name))
			}
			handles[h] = filepath.Join(dir, name)
		}
	}
	return handles, nil
}
