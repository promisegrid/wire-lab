package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// ManifestPath is the path of the manifest relative to the repo root.
// ManifestPath is the relative path to the spec manifest, anchored at
// repo root. Per TE-29 + TE-32, the wire-lab harness's specs (including
// the manifest) live under protocols/wire-lab.d/specs/. The legacy
// top-level specs/ layout is no longer used.
const ManifestPath = "protocols/wire-lab.d/specs/MANIFEST.md"

// SpecsDir is the relative path to the wire-lab harness's specs
// directory (formerly top-level specs/, now under the harness's own
// protocol-as-simrepo directory per TE-29).
const SpecsDir = "protocols/wire-lab.d/specs"

// yamlFenceOpen and yamlFenceClose mark the boundaries of the authoritative
// YAML block in MANIFEST.md. The opening fence MUST be exactly "```yaml" on a
// line by itself; the closing fence MUST be exactly "```" on a line by itself.
// Per DI-011-20260429-184500, exactly one such fenced block exists in the file
// and is authoritative.
const (
	yamlFenceOpen  = "```yaml"
	yamlFenceClose = "```"
)

// Status is one of the locked status values for a manifest entry
// (per DI-011-20260429-184500).
type Status string

const (
	StatusFrozen      Status = "frozen"
	StatusSuperseded  Status = "superseded"
	StatusDraftAhead  Status = "draft-ahead"
)

// Entry is one frozen spec's record in the manifest.
//
// Field shape mirrors what TE-22 §Conclusions/Implications committed to:
// pCID, slug, status, frozen-on, supersedes, superseded-by, depends-on,
// freezing-commit, notes.
//
// YAML keys use snake_case to match the prose convention.
type Entry struct {
	PCID            string   `yaml:"pcid"`
	Slug            string   `yaml:"slug"`
	Status          Status   `yaml:"status"`
	FrozenOn        string   `yaml:"frozen_on"`
	Supersedes      string   `yaml:"supersedes,omitempty"`
	SupersededBy    string   `yaml:"superseded_by,omitempty"`
	DependsOn       []string `yaml:"depends_on,omitempty"`
	FreezingCommit  string   `yaml:"freezing_commit,omitempty"`
	Notes           string   `yaml:"notes,omitempty"`
}

// Manifest is the parsed, in-memory form of MANIFEST.md.
//
// The Prose field is split into preamble (everything before the YAML fence)
// and postamble (everything after the closing fence) so that round-trips
// preserve human-written content. Entries holds the parsed YAML payload.
type Manifest struct {
	Preamble  string
	Entries   []Entry
	Postamble string
}

// ManifestData is the YAML structure inside the fenced block.
//
// Wrapped in a struct so the fenced block parses as a YAML document, not a
// raw list. This matches the "single document with one entries: list" shape
// emitted by emitYAML below.
type ManifestData struct {
	Entries []Entry `yaml:"entries"`
}

// readManifest reads MANIFEST.md, splits out the YAML block, parses it, and
// returns the resulting Manifest. If the file does not exist, returns a
// Manifest with empty preamble/postamble/entries (caller decides whether
// that is an error).
func readManifest(repoRoot string) (*Manifest, error) {
	path := repoRoot + "/" + ManifestPath
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return parseManifest(string(raw))
}

// parseManifest splits a Markdown file's contents into preamble + YAML +
// postamble and parses the YAML. It enforces "exactly one fenced YAML block"
// (per DI-011-20260429-184500): zero or more than one is an error.
func parseManifest(s string) (*Manifest, error) {
	openIdx := indexLine(s, yamlFenceOpen)
	if openIdx == -1 {
		return nil, fmt.Errorf("manifest: no `%s` fence found; manifest must contain exactly one fenced YAML block", yamlFenceOpen)
	}
	// Find a "```" line strictly after the opening fence.
	postOpen := openIdx + len(yamlFenceOpen)
	if postOpen < len(s) && s[postOpen] == '\n' {
		postOpen++
	}
	closeRel := indexLine(s[postOpen:], yamlFenceClose)
	if closeRel == -1 {
		return nil, fmt.Errorf("manifest: opening `%s` fence has no matching closing `%s` fence", yamlFenceOpen, yamlFenceClose)
	}
	closeIdx := postOpen + closeRel
	// Reject a second `yamlFenceOpen` later in the file.
	rest := s[closeIdx+len(yamlFenceClose):]
	if indexLine(rest, yamlFenceOpen) != -1 {
		return nil, fmt.Errorf("manifest: more than one `%s` fence found; exactly one is allowed", yamlFenceOpen)
	}

	preamble := s[:openIdx]
	yamlBody := s[postOpen:closeIdx]
	postamble := s[closeIdx+len(yamlFenceClose):]
	// Trim a single leading newline off postamble (the one immediately after
	// the closing fence) so round-trips don't accumulate blank lines.
	if strings.HasPrefix(postamble, "\n") {
		postamble = postamble[1:]
	}

	var data ManifestData
	if strings.TrimSpace(yamlBody) != "" {
		if err := yaml.Unmarshal([]byte(yamlBody), &data); err != nil {
			return nil, fmt.Errorf("manifest: yaml parse: %w", err)
		}
	}
	return &Manifest{
		Preamble:  preamble,
		Entries:   data.Entries,
		Postamble: postamble,
	}, nil
}

// emitYAML serializes the entries slice as YAML in deterministic key order
// (relying on yaml.v3's struct-field ordering). A trailing newline is included
// so the closing fence sits on its own line.
func emitYAML(entries []Entry) (string, error) {
	data := ManifestData{Entries: entries}
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return "", err
	}
	if err := enc.Close(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// writeManifest re-serializes the Manifest and writes it to disk. The output
// shape is: preamble + opening fence + YAML + closing fence + postamble. The
// preamble and postamble are written verbatim (round-trip stable).
func writeManifest(repoRoot string, m *Manifest) error {
	yamlText, err := emitYAML(m.Entries)
	if err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString(m.Preamble)
	b.WriteString(yamlFenceOpen)
	b.WriteString("\n")
	b.WriteString(yamlText)
	b.WriteString(yamlFenceClose)
	b.WriteString("\n")
	b.WriteString(m.Postamble)
	path := repoRoot + "/" + ManifestPath
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

// findEntryByPCID returns a pointer to the entry whose PCID matches, or nil.
func (m *Manifest) findEntryByPCID(pcid string) *Entry {
	for i := range m.Entries {
		if m.Entries[i].PCID == pcid {
			return &m.Entries[i]
		}
	}
	return nil
}

// latestFrozen returns the most-recent entry for a given slug whose status is
// "frozen". Used by `spec freeze` to decide what the supersedes link of a new
// entry should be, and by `spec check` to compute draft-ahead state.
//
// "Most recent" is determined by frozen_on string comparison, which works
// because timestamps are stored in ISO 8601 form (YYYY-MM-DDTHH:MM:SSZ).
func (m *Manifest) latestFrozen(slug string) *Entry {
	var latest *Entry
	for i := range m.Entries {
		e := &m.Entries[i]
		if e.Slug != slug || e.Status != StatusFrozen {
			continue
		}
		if latest == nil || e.FrozenOn > latest.FrozenOn {
			latest = e
		}
	}
	return latest
}

// sortEntries sorts entries deterministically: by slug ascending, then
// frozen_on ascending. This makes diffs to MANIFEST.md easy to read after
// each freeze operation.
func sortEntries(entries []Entry) {
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].Slug != entries[j].Slug {
			return entries[i].Slug < entries[j].Slug
		}
		return entries[i].FrozenOn < entries[j].FrozenOn
	})
}

// indexLine returns the byte offset of a line in s whose content (excluding
// the trailing newline) equals exactly `line`, or -1 if no such line exists.
// The line must be at the start of s or be preceded by a newline. This is
// stricter than strings.Index because we don't want to match a prefix inside
// a longer line (e.g., "```yaml-extension" should not match "```yaml").
func indexLine(s, line string) int {
	for i := 0; i <= len(s)-len(line); {
		// Check that we're at line-start.
		if i > 0 && s[i-1] != '\n' {
			// Skip to next line.
			next := strings.IndexByte(s[i:], '\n')
			if next == -1 {
				return -1
			}
			i += next + 1
			continue
		}
		if strings.HasPrefix(s[i:], line) {
			// And the matched range is followed by '\n' or end of string.
			end := i + len(line)
			if end == len(s) || s[end] == '\n' {
				return i
			}
		}
		next := strings.IndexByte(s[i:], '\n')
		if next == -1 {
			return -1
		}
		i += next + 1
	}
	return -1
}
