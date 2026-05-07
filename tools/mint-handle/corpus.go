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
// every directory matching any glob is also scanned. Missing directories are
// skipped so the tool works on both current main's root layout and ppx/main's
// protocolized TODO layout.
//
// Intent: Collision-check the shared TODO, TE, DR, and DI handle namespace
// across the layouts main and ppx/main need to merge. Source: DI-nisam
var scanLiteralDirs = []string{
	"docs/thought-experiments",
	"TODO",
	"DR",
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
//	<KIND>   = TE | TODO | DR
//	<HANDLE> = proquint-1 (5 chars, CVCVC, alphabets in proquint.go)
//	         | proquint-2 (5 + '-' + 5 = 11 chars)
//	<SLUG>   = lower-kebab-case description, may include digits
//
// The regex is anchored so it does not match files whose names happen to
// contain proquint-shaped substrings inside the slug.
var handleFileRE = regexp.MustCompile(
	`^(?:TE|TODO|DR)-(` +
		// proquint-2 must come first so the alternation prefers it
		// over the proquint-1 prefix it begins with.
		`[bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz]` +
		`-` +
		`[bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz]` +
		`|` +
		`[bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz]` +
		`)-`)

// diOwnerRE matches DI owner lines, not ordinary references to DI IDs. This
// distinction lets DRs and docs cite a DI without being treated as owning a
// second copy of that handle.
var diOwnerRE = regexp.MustCompile(`(?m)^ID:\s*DI-(` +
	`[bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz]` +
	`(?:-[bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz])?` +
	`)\s*$`)

// scanCorpus returns the set of handles currently in use across the
// configured directories rooted at repoRoot. It is the collision check set
// for mint.
//
// Files whose names do not match handleFileRE are ignored for filename-owned
// handles. Markdown files in TODO directories are additionally scanned for
// "ID: DI-<handle>" owner lines. Legacy integer and timestamp filenames are
// intentionally ignored because this task only changes newly-created
// artifacts.
func scanCorpus(repoRoot string) (map[string]string, error) {
	handles := make(map[string]string)
	dirs, err := scanDirs(repoRoot)
	if err != nil {
		return nil, err
	}
	for _, dir := range dirs {
		if err := scanDir(repoRoot, dir, handles); err != nil {
			return nil, err
		}
	}
	return handles, nil
}

// scanDirs resolves the configured literal directories and globbed
// directories. Missing directories are skipped because main and ppx/main do
// not yet share the same coordination layout.
func scanDirs(repoRoot string) ([]string, error) {
	dirs := make([]string, 0, len(scanLiteralDirs)+len(scanGlobs))
	for _, dir := range scanLiteralDirs {
		full := filepath.Join(repoRoot, dir)
		info, err := os.Stat(full)
		switch {
		case err == nil && info.IsDir():
			dirs = append(dirs, dir)
		case err == nil:
			return nil, fmt.Errorf("scan path %s is not a directory", dir)
		case os.IsNotExist(err):
			continue
		default:
			return nil, fmt.Errorf("stat %s: %w", dir, err)
		}
	}
	for _, glob := range scanGlobs {
		matches, err := filepath.Glob(filepath.Join(repoRoot, glob))
		if err != nil {
			return nil, fmt.Errorf("glob %s: %w", glob, err)
		}
		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				return nil, fmt.Errorf("stat glob match %s: %w", match, err)
			}
			if !info.IsDir() {
				continue
			}
			rel, err := filepath.Rel(repoRoot, match)
			if err != nil {
				return nil, fmt.Errorf("relpath %s: %w", match, err)
			}
			dirs = append(dirs, rel)
		}
	}
	return dirs, nil
}

// scanDir records filename-owned handles and DI-owner handles from one
// configured directory.
func scanDir(repoRoot, dir string, handles map[string]string) error {
	full := filepath.Join(repoRoot, dir)
	entries, err := os.ReadDir(full)
	if err != nil {
		return fmt.Errorf("scan %s: %w", dir, err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".md") {
			continue
		}
		relPath := filepath.Join(dir, name)
		if match := handleFileRE.FindStringSubmatch(name); match != nil {
			if err := rememberHandle(handles, match[1], relPath); err != nil {
				return err
			}
		}
		if isTODODir(dir) {
			if err := scanDIHandles(filepath.Join(full, name), relPath, handles); err != nil {
				return err
			}
		}
	}
	return nil
}

// isTODODir recognizes both the current root TODO layout and ppx/main's
// per-protocol TODO directories.
func isTODODir(dir string) bool {
	return dir == "TODO" || strings.HasSuffix(filepath.ToSlash(dir), "/TODO")
}

// scanDIHandles records the handles owned by Decision Intent Log entries in
// TODO files. A DI handle is global, so it must not duplicate any filename
// handle or any other DI owner.
func scanDIHandles(path, relPath string, handles map[string]string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", relPath, err)
	}
	for _, match := range diOwnerRE.FindAllStringSubmatch(string(data), -1) {
		if err := rememberHandle(handles, match[1], relPath+"#DI-"+match[1]); err != nil {
			return err
		}
	}
	return nil
}

// rememberHandle inserts one owned handle into the corpus and reports a clear
// error if the repo already contains a duplicate owner.
func rememberHandle(handles map[string]string, handle, owner string) error {
	if previous, duplicate := handles[handle]; duplicate {
		return fmt.Errorf(
			"corpus already contains duplicate handle %q in %s and %s",
			handle, previous, owner)
	}
	handles[handle] = owner
	return nil
}
