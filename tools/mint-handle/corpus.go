package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// scanDirs are the wire-lab directories scanned for existing handles.
// A file in any of these whose name matches handleFileRE is treated as
// owning its handle. The corpus union of those handles is the collision
// check set for new mints.
//
// Convention-based rather than registry-based: there is no separate
// HANDLES.md or registry file. Filenames are the ground truth. To add a
// new artifact type that also bears handles, add its directory here.
//
// See TE-39 for the lock decision and TE-vapoj (substrate-agnostic
// layered model) for the broader rationale.
var scanDirs = []string{
	"docs/thought-experiments",
	"protocols/wire-lab.d/TODO",
}

// handleFileRE matches a wire-lab handle-bearing filename and captures the
// handle in group 1.
//
// Filename grammar:
//   <KIND>-<HANDLE>-<SLUG>.md
// where:
//   <KIND>   = TE | TODO
//   <HANDLE> = proquint-1 (5 chars, CVCVC, alphabets in proquint.go)
//            | proquint-2 (5 + '-' + 5 = 11 chars)
//   <SLUG>   = lower-kebab-case description, may include digits
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
	for _, dir := range scanDirs {
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
