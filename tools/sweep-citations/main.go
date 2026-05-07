// Command sweep-citations rewrites integer-alias and timestamp-alias
// references to TE/TODO files into the new proquint handle form, using
// the mapping produced by tools/migrate-handles/.
//
// Three rewrite classes:
//
//  1. Bare integer aliases ("TE-38", "TODO 22"): rewritten to the
//     proquint handle ("TE-vunub", "TODO-vunub" form -- KIND-handle).
//     For TE bodies the form is conventionally "TE-<handle>"; for TODO
//     bodies it is "TODO-<handle>".
//  2. Markdown link targets containing timestamped filenames
//     ([title](TE-YYYYMMDD-HHMMSS-slug.md)): rewritten to the new
//     proquint filename so the links keep resolving.
//  3. Bare timestamp identifiers in prose ("TE-20260506-184800",
//     "TODO-20260507-002306"): rewritten to "TE-<handle>" /
//     "TODO-<handle>" forms.
//
// Files explicitly skipped (by path prefix or filename pattern):
//
//   - protocols/wire-lab.d/TODO/pre*-audit-report-*.md
//   - protocols/wire-lab.d/TODO/dropped-thread-disposition-*.md
//   - proposals/  (frozen historical: contests, reviews, approved
//     bundles)
//   - protocols/wire-lab.d/specs/harness-spec-bafk*.md (frozen
//     CID-named spec snapshot)
//   - tools/migrate-handles/mapping.tsv (source of truth, untouched)
//   - tools/sweep-citations/  (this tool's own source)
//
// Within otherwise-eligible files, lines inside a "## Prior aliases"
// section are skipped (they are the migration's own historical
// record). The section starts at "## Prior aliases" and ends at the
// next "## " heading or end-of-file.
//
// Usage:
//
//	sweep-citations [-r REPO_ROOT] [-n] [-q]
//
// Flags mirror migrate-handles: -r repo root, -n dry-run, -q quiet.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// row is a parsed line of mapping.tsv.
type row struct {
	handle   string
	kind     string // TE | TODO
	intAlias string // e.g. TE-38, TODO-22, or "" if none
	tsAlias  string // e.g. TE-20260506-184800
	newPath  string // e.g. docs/thought-experiments/TE-vunub-...md
	oldPath  string // legacy path; "(new in this twig)" for TE-39
}

// rep is a single string-replacement instruction.
type rep struct {
	from string
	to   string
}

func main() {
	repoRoot := flag.String("r", ".", "repo root")
	dryRun := flag.Bool("n", false, "dry-run")
	quiet := flag.Bool("q", false, "suppress per-file progress")
	flag.Parse()

	root, err := filepath.Abs(*repoRoot)
	if err != nil {
		die("repo root: %v", err)
	}

	rows := loadMapping(filepath.Join(root, "tools/migrate-handles/mapping.tsv"))

	// Build replacement lists. They are sorted longest-first to avoid
	// prefix-match ambiguity (TE-3 vs TE-38).
	var intReps, tsReps []rep
	pathReps := map[string]string{}
	for _, r := range rows {
		if r.intAlias != "" {
			// "TE-38" -> "TE-vunub" ; "TODO-22" -> "TODO-vunub".
			// Mapping rows already encode the kind in intAlias so the
			// substitution preserves it.
			intReps = append(intReps, rep{r.intAlias, r.kind + "-" + r.handle})
		}
		if r.tsAlias != "" {
			tsReps = append(tsReps, rep{r.tsAlias, r.kind + "-" + r.handle})
		}
		// Path rewrites only for files that had an oldPath (not the
		// fresh TE-39 entry).
		if r.oldPath != "(new in this twig)" {
			pathReps[filepath.Base(r.oldPath)] = filepath.Base(r.newPath)
		}
	}

	// Sort integer reps so longer aliases come first (TE-38 before TE-3).
	sort.Slice(intReps, func(i, j int) bool {
		if len(intReps[i].from) != len(intReps[j].from) {
			return len(intReps[i].from) > len(intReps[j].from)
		}
		return intReps[i].from > intReps[j].from
	})

	// Walk the repo and rewrite each eligible file.
	var changed, scanned int
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if shouldSkipDir(path, root) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		if shouldSkipFile(rel) {
			return nil
		}
		scanned++
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		newBody, n := sweep(string(body), intReps, tsReps, pathReps)
		if n == 0 {
			return nil
		}
		if !*quiet {
			fmt.Fprintf(os.Stderr, "  %s  %d edits\n", rel, n)
		}
		if !*dryRun {
			if err := os.WriteFile(path, []byte(newBody), info.Mode()); err != nil {
				return err
			}
		}
		changed++
		return nil
	})
	if err != nil {
		die("walk: %v", err)
	}
	fmt.Fprintf(os.Stderr, "sweep-citations: scanned=%d, changed=%d, dry-run=%v\n",
		scanned, changed, *dryRun)
}

func loadMapping(path string) []row {
	f, err := os.Open(path)
	if err != nil {
		die("open mapping: %v", err)
	}
	defer f.Close()
	var rows []row
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 1<<20), 1<<20)
	first := true
	for sc.Scan() {
		line := sc.Text()
		if first {
			first = false
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 6 {
			continue
		}
		rows = append(rows, row{
			handle: fields[0], kind: fields[1], intAlias: fields[2],
			tsAlias: fields[3], newPath: fields[4], oldPath: fields[5],
		})
	}
	if err := sc.Err(); err != nil {
		die("read mapping: %v", err)
	}
	return rows
}

// shouldSkipDir keeps the walker out of frozen / non-prose subtrees.
func shouldSkipDir(path, root string) bool {
	rel, _ := filepath.Rel(root, path)
	switch rel {
	case ".git", "tools", "proposals", "session-logs":
		return true
	}
	// also skip nested .git inside worktrees, but the top-level .git
	// catch above handles the common case. Anything with .git in any
	// component is frozen for git's own purposes.
	if strings.Contains(rel, "/.git/") {
		return true
	}
	return false
}

// shouldSkipFile returns true for individual frozen files inside
// otherwise-walked directories.
func shouldSkipFile(rel string) bool {
	bn := filepath.Base(rel)
	switch {
	case strings.HasPrefix(bn, "pre18-audit-report-"),
		strings.HasPrefix(bn, "pre149-audit-report-"),
		strings.HasPrefix(bn, "dropped-thread-disposition-"),
		strings.HasPrefix(bn, "harness-spec-bafk"):
		return true
	}
	// Self-declared append-only files (per their own frontmatter/forward
	// pointer). TODO-lilar declares: "this file is append-only history and
	// is not retroactively rewritten" (forward pointer added 2026-05-06).
	switch filepath.Base(rel) {
	case "TODO-lilar-session-replay-cleanup.md":
		return true
	}
	// Index files are rebuilt in chunk D, not swept here.
	switch rel {
	case "docs/thought-experiments/README.md",
		"protocols/wire-lab.d/TODO/TODO.md",
		"protocols/group-session.d/TODO/TODO.md",
		"protocols/ppx-dr.d/TODO/TODO.md",
		"protocols/udp-binding.d/TODO/TODO.md":
		return true
	}
	return false
}

// sweep applies the three rewrite classes to body, skipping any "## Prior
// aliases" section. Returns the new body and the count of substitutions
// made.
func sweep(body string, intReps, tsReps []rep, pathReps map[string]string) (string, int) {
	// Split body into "regions" by Prior-aliases boundaries. Even-indexed
	// regions are sweepable; odd-indexed are inside a Prior-aliases
	// section and pass through verbatim.
	regions := splitByPriorAliases(body)
	var n int
	for i, region := range regions {
		if i%2 == 1 {
			continue
		}
		s, c := applyReps(region, intReps, tsReps, pathReps)
		regions[i] = s
		n += c
	}
	return strings.Join(regions, ""), n
}

// splitByPriorAliases returns alternating slices: regions[0] is before
// the first Prior aliases section; regions[1] is the first Prior aliases
// section (verbatim); regions[2] is between sections; etc.
func splitByPriorAliases(body string) []string {
	headRE := regexp.MustCompile(`(?m)^##\s+Prior aliases\s*$`)
	nextHeadRE := regexp.MustCompile(`(?m)^##\s+`)
	var out []string
	rest := body
	for {
		loc := headRE.FindStringIndex(rest)
		if loc == nil {
			out = append(out, rest)
			break
		}
		out = append(out, rest[:loc[0]])
		// section starts at loc[0]; find its end (next "## " heading or EOF).
		afterHead := rest[loc[1]:]
		nh := nextHeadRE.FindStringIndex(afterHead)
		var endRel int
		if nh == nil {
			endRel = len(rest)
		} else {
			endRel = loc[1] + nh[0]
		}
		out = append(out, rest[loc[0]:endRel])
		rest = rest[endRel:]
	}
	return out
}

func applyReps(s string, intReps, tsReps []rep, pathReps map[string]string) (string, int) {
	var n int
	// Class 1: integer aliases. Use word-boundary regex per replacement
	// to avoid matching inside larger tokens (TE-3 inside TE-38, etc.).
	for _, r := range intReps {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(r.from) + `\b`)
		s2 := re.ReplaceAllStringFunc(s, func(_ string) string {
			n++
			return r.to
		})
		s = s2
	}
	// Also handle "TODO N" form (space instead of hyphen). The mapping
	// stores "TODO-N" via the kind+integer concatenation; we also need
	// to catch "TODO 22" prose.
	for _, r := range intReps {
		if !strings.HasPrefix(r.from, "TODO-") {
			continue
		}
		spaceForm := strings.Replace(r.from, "TODO-", "TODO ", 1)
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(spaceForm) + `\b`)
		s = re.ReplaceAllStringFunc(s, func(_ string) string {
			n++
			return r.to
		})
	}
	// Class 2: timestamp aliases.
	for _, r := range tsReps {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(r.from) + `\b`)
		s = re.ReplaceAllStringFunc(s, func(_ string) string {
			n++
			return r.to
		})
	}
	// Class 3: legacy filename basenames in markdown link targets.
	// Only rewrite when the basename appears as a token (preceded by
	// '/' or '(' or whitespace, followed by ')' or whitespace). The
	// path itself may have any directory prefix; we replace just the
	// basename.
	for old, new := range pathReps {
		// Word-boundary won't help here because filenames contain '-'
		// and '.'. We use a lookalike: preceded by [^A-Za-z0-9_-] or
		// start-of-string, followed by [^A-Za-z0-9_-] or end.
		re := regexp.MustCompile(regexp.QuoteMeta(old))
		s = re.ReplaceAllStringFunc(s, func(_ string) string {
			n++
			return new
		})
	}
	return s, n
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "sweep-citations: "+format+"\n", args...)
	os.Exit(1)
}
