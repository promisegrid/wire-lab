// Command migrate-handles performs the one-shot TE-39 migration: it
// renames every TE and TODO file from the legacy
// (TE|TODO)-<timestamp>-<slug>.md form to the new
// (TE|TODO)-<proquint>-<slug>.md form, and injects a
// "## Prior aliases" section into each file body recording the integer
// alias and timestamp the file used to carry.
//
// This tool is run exactly once per repo. After it runs, mint-handle is
// the steady-state tool for new files; this tool's job is done.
//
// Usage:
//
//	migrate-handles [-r REPO_ROOT] [-n] [-q]
//
// Flags:
//
//	-r   repo root (default: current directory)
//	-n   dry-run: print every action without changing the filesystem
//	-q   quiet: suppress per-file progress lines
//
// The integer-alias map is read from
// docs/thought-experiments/README.md (TE table) and
// protocols/wire-lab.d/TODO/TODO.md (master TODO cross-list). Files
// without an integer alias in either index are migrated with only a
// timestamp prior alias.
//
// For each file, the tool:
//  1. Computes the timestamp prior alias from the legacy filename
//     (TE-YYYYMMDD-HHMMSS form).
//  2. Looks up the integer prior alias from the indexes (e.g. TE-38).
//  3. Mints a fresh proquint-1 handle via the same scanCorpus + retry
//     loop as the mint-handle tool, but operates against an in-memory
//     "claimed" set so the 70-file batch run cannot self-collide.
//  4. Renames the file via `git mv` (so history follows).
//  5. Inserts a "## Prior aliases" section into the file body
//     immediately after the existing "## TE ID" or "## TODO ID"
//     section, or after the H1 if no such section exists.
//
// The renames and edits are staged but not committed; commit is left to
// the operator (the rename script).
package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// proquintCons / proquintVows are duplicated from tools/mint-handle/
// because Go modules don't share unexported symbols across modules.
// Keeping them aligned with mint-handle's tables is enforced by the
// vector tests both modules carry; if they ever diverge a vector test
// will fire.
const (
	proquintCons = "bdfghjklmnprstvz"
	proquintVows = "aiou"
)

func uint16ToProquint(n uint16) string {
	c1 := proquintCons[(n>>12)&0xF]
	v1 := proquintVows[(n>>10)&0x3]
	c2 := proquintCons[(n>>6)&0xF]
	v2 := proquintVows[(n>>4)&0x3]
	c3 := proquintCons[n&0xF]
	return string([]byte{c1, v1, c2, v2, c3})
}

// scanLiteralDirs and scanGlobs match tools/mint-handle/corpus.go.
var scanLiteralDirs = []string{
	"docs/thought-experiments",
}
var scanGlobs = []string{
	"protocols/*/TODO",
}

// handleFileRE matches the new (post-migration) filename form.
var handleFileRE = regexp.MustCompile(
	`^(?:TE|TODO)-(` +
		`[bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz]` +
		`-` +
		`[bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz]` +
		`|` +
		`[bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz]` +
		`)-`)

// legacyFileRE matches the pre-migration filename form. Captures:
//
//	[1] kind (TE|TODO)
//	[2] timestamp (YYYYMMDD-HHMMSS)
//	[3] slug
var legacyFileRE = regexp.MustCompile(
	`^(TE|TODO)-(\d{8}-\d{6})-(.+)\.md$`)

// indexLineRE matches a row in the README TE index or TODO master list.
// The README uses "| TE-1 | ... | [Title](path) |". The TODO master
// list uses "| TODO 1 | ... | [path](path) |". The separator between the
// kind and the integer is therefore either '-' (TE README) or whitespace
// (TODO master).
//
// Captures:
//
//	[1] integer alias (e.g. "38")
//	[2] filename (basename, possibly preceded by relative path segments)
var indexLineRE = regexp.MustCompile(
	`\|\s*(?:TE|TODO)[\s-]+(\d+)\s*\|.*?\(([^)]+\.md)\)`)

func main() {
	repoRoot := flag.String("r", ".", "repo root")
	dryRun := flag.Bool("n", false, "dry-run: print actions without changing the filesystem")
	quiet := flag.Bool("q", false, "suppress per-file progress lines")
	flag.Parse()

	root, err := filepath.Abs(*repoRoot)
	if err != nil {
		die("repo root: %v", err)
	}

	// Build integer-alias map from the two index files.
	intAlias := map[string]string{} // basename -> "TE-N" or "TODO-N"
	for _, idx := range []struct {
		path, kind string
	}{
		{"docs/thought-experiments/README.md", "TE"},
		{"protocols/wire-lab.d/TODO/TODO.md", "TODO"},
	} {
		readIndex(filepath.Join(root, idx.path), idx.kind, intAlias)
	}

	// Build the legacy file list (every file we plan to rename).
	type job struct {
		dir       string // relative directory
		oldName   string // basename, legacy form
		kind      string // TE | TODO
		timestamp string // YYYYMMDD-HHMMSS
		slug      string // hyphen-separated tail
		intAlias  string // e.g. "TE-38" or "" if not in the index
	}
	var jobs []job
	dirs := append([]string(nil), scanLiteralDirs...)
	for _, glob := range scanGlobs {
		matches, err := filepath.Glob(filepath.Join(root, glob))
		if err != nil {
			die("glob %s: %v", glob, err)
		}
		for _, m := range matches {
			rel, _ := filepath.Rel(root, m)
			dirs = append(dirs, rel)
		}
	}
	sort.Strings(dirs)
	for _, dir := range dirs {
		entries, err := os.ReadDir(filepath.Join(root, dir))
		if err != nil {
			die("scan %s: %v", dir, err)
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			m := legacyFileRE.FindStringSubmatch(e.Name())
			if m == nil {
				continue
			}
			j := job{
				dir:       dir,
				oldName:   e.Name(),
				kind:      m[1],
				timestamp: m[2],
				slug:      m[3],
			}
			if alias, ok := intAlias[e.Name()]; ok {
				j.intAlias = alias
			}
			jobs = append(jobs, j)
		}
	}
	sort.Slice(jobs, func(i, k int) bool {
		// Stable order: by timestamp (which corresponds to integer order
		// where aliases exist), then by oldName.
		if jobs[i].timestamp != jobs[k].timestamp {
			return jobs[i].timestamp < jobs[k].timestamp
		}
		return jobs[i].oldName < jobs[k].oldName
	})
	if !*quiet {
		fmt.Fprintf(os.Stderr, "migrate-handles: %d files to rename\n", len(jobs))
	}

	// Mint handles for every job. The "claimed" set starts empty (no
	// post-migration filenames exist yet) and grows as we mint.
	claimed := map[string]string{}
	seed := uint64(time.Now().UnixNano())
	for i := range jobs {
		var h string
		for attempt := 0; attempt < 1000; attempt++ {
			seed++
			var b [8]byte
			binary.BigEndian.PutUint64(b[:], seed)
			sum := sha256.Sum256(b[:])
			n := binary.BigEndian.Uint16(sum[:2])
			cand := uint16ToProquint(n)
			if _, taken := claimed[cand]; !taken {
				h = cand
				break
			}
		}
		if h == "" {
			die("could not mint handle for %s after 1000 attempts", jobs[i].oldName)
		}
		newName := fmt.Sprintf("%s-%s-%s.md", jobs[i].kind, h, jobs[i].slug)
		claimed[h] = filepath.Join(jobs[i].dir, newName)

		oldPath := filepath.Join(jobs[i].dir, jobs[i].oldName)
		newPath := filepath.Join(jobs[i].dir, newName)

		if !*quiet {
			alias := jobs[i].intAlias
			if alias == "" {
				alias = "(no integer alias)"
			}
			fmt.Fprintf(os.Stderr, "  %s  %s  %s -> %s\n",
				jobs[i].kind, h, alias, jobs[i].oldName)
		}

		if *dryRun {
			continue
		}

		// Inject the Prior aliases section into the file body BEFORE the
		// rename, so the in-place edit + git-mv yields a clean rename
		// detection plus a content edit on the new path.
		fullOld := filepath.Join(root, oldPath)
		body, err := os.ReadFile(fullOld)
		if err != nil {
			die("read %s: %v", oldPath, err)
		}
		newBody := injectPriorAliases(string(body), jobs[i].kind,
			jobs[i].intAlias, fmt.Sprintf("%s-%s", jobs[i].kind, jobs[i].timestamp))
		if err := os.WriteFile(fullOld, []byte(newBody), 0o644); err != nil {
			die("write %s: %v", oldPath, err)
		}

		// git mv from old to new.
		cmd := exec.Command("git", "mv", oldPath, newPath)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			die("git mv %s -> %s: %v\n%s", oldPath, newPath, err, out)
		}
	}

	// Emit the mapping file for downstream consumers (README rebuilder,
	// citation sweeper). One line per file: <handle> <kind> <intAlias>
	// <timestampAlias> <newPath>.
	if !*dryRun {
		mapPath := filepath.Join(root, "tools/migrate-handles/mapping.tsv")
		f, err := os.Create(mapPath)
		if err != nil {
			die("create mapping.tsv: %v", err)
		}
		w := bufio.NewWriter(f)
		fmt.Fprintln(w, "handle\tkind\tint_alias\tts_alias\tnew_path\told_path")
		for _, j := range jobs {
			// re-derive the new name from claimed (reverse map by path)
			var handle, newPath string
			for h, p := range claimed {
				if filepath.Base(p) == fmt.Sprintf("%s-%s-%s.md", j.kind, h, j.slug) &&
					filepath.Dir(p) == j.dir {
					handle = h
					newPath = p
					break
				}
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s-%s\t%s\t%s\n",
				handle, j.kind, j.intAlias, j.kind, j.timestamp, newPath,
				filepath.Join(j.dir, j.oldName))
		}
		w.Flush()
		f.Close()
		if !*quiet {
			fmt.Fprintf(os.Stderr, "wrote mapping to %s\n", mapPath)
		}
	}
}

// readIndex parses an index file (README.md or TODO.md) and adds entries
// to intAlias keyed by basename.
func readIndex(path, kind string, intAlias map[string]string) {
	data, err := os.ReadFile(path)
	if err != nil {
		die("read index %s: %v", path, err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		m := indexLineRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		// m[1] = integer, m[2] = filename (possibly relative). Take basename.
		bn := filepath.Base(m[2])
		intAlias[bn] = fmt.Sprintf("%s-%s", kind, m[1])
	}
}

// injectPriorAliases inserts a "## Prior aliases" section into body. It
// is placed immediately after the existing "## TE ID" / "## TODO ID"
// section (the section ends at the next blank line followed by a heading
// or end-of-file). If no ID section exists, the block is inserted after
// the H1.
func injectPriorAliases(body, kind, intAlias, tsAlias string) string {
	var lines []string
	if intAlias != "" {
		lines = []string{
			"## Prior aliases",
			"",
			fmt.Sprintf("Before the TE-39 proquint migration, this file was known as:"),
			"",
			fmt.Sprintf("- `%s` (integer alias)", intAlias),
			fmt.Sprintf("- `%s` (timestamp alias and pre-migration filename)", tsAlias),
			"",
		}
	} else {
		lines = []string{
			"## Prior aliases",
			"",
			"Before the TE-39 proquint migration, this file was known as:",
			"",
			fmt.Sprintf("- `%s` (timestamp alias and pre-migration filename;"+
				" no integer alias was assigned)", tsAlias),
			"",
		}
	}
	block := strings.Join(lines, "\n")

	// Find a section heading like "## TE ID" or "## TODO ID" (case-
	// insensitive on the kind). Insert AFTER that section.
	headRE := regexp.MustCompile(`(?im)^##\s+(TE|TODO)\s+ID\s*$`)
	loc := headRE.FindStringIndex(body)
	if loc == nil {
		// Fall back to inserting after the H1.
		h1RE := regexp.MustCompile(`(?m)^#\s+.+$`)
		h1 := h1RE.FindStringIndex(body)
		if h1 == nil {
			// No H1 either; just prepend.
			return block + "\n" + body
		}
		// Insert after the H1 plus its trailing blank line.
		insert := h1[1]
		// Skip any single newline that follows the H1.
		if insert < len(body) && body[insert] == '\n' {
			insert++
		}
		// Skip a possible blank line.
		if insert < len(body) && body[insert] == '\n' {
			insert++
		}
		return body[:insert] + block + "\n" + body[insert:]
	}

	// Find the end of the ID section: the next "## " heading, or EOF.
	rest := body[loc[1]:]
	nextHead := regexp.MustCompile(`(?m)^##\s+`).FindStringIndex(rest)
	var endOfSection int
	if nextHead == nil {
		endOfSection = len(body)
	} else {
		endOfSection = loc[1] + nextHead[0]
	}
	// Insert the block at endOfSection. The section's trailing whitespace
	// is preserved; the block carries its own trailing blank line.
	return body[:endOfSection] + block + "\n" + body[endOfSection:]
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "migrate-handles: "+format+"\n", args...)
	os.Exit(1)
}
