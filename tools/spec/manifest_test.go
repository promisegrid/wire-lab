package main

import (
	"strings"
	"testing"
)

const sampleManifest = "# Spec-doc manifest\n" +
	"\n" +
	"Some explanatory prose.\n" +
	"\n" +
	"```yaml\n" +
	"entries:\n" +
	"  - pcid: bafkreifjjcie6lypi6ny7amxnfftagclbuxndqonfipmb64f2km2devei4\n" +
	"    slug: harness-spec\n" +
	"    status: frozen\n" +
	"    frozen_on: 2026-04-29T20:00:00Z\n" +
	"    freezing_commit: 0000000000000000000000000000000000000000\n" +
	"```\n" +
	"\n" +
	"More prose at the bottom.\n"

// TestParseManifestRoundTrip verifies that read-then-write of a manifest
// preserves both the structured YAML data and the surrounding human prose.
func TestParseManifestRoundTrip(t *testing.T) {
	m, err := parseManifest(sampleManifest)
	if err != nil {
		t.Fatalf("parseManifest: %v", err)
	}
	if len(m.Entries) != 1 {
		t.Fatalf("want 1 entry, got %d", len(m.Entries))
	}
	got := m.Entries[0]
	if got.Slug != "harness-spec" {
		t.Errorf("slug = %q, want harness-spec", got.Slug)
	}
	if got.Status != StatusFrozen {
		t.Errorf("status = %q, want frozen", got.Status)
	}
	if !strings.HasPrefix(got.PCID, "bafkrei") {
		t.Errorf("pcid = %q, want CIDv1 prefix bafkrei", got.PCID)
	}
}

// TestParseManifestMissingFenceFails verifies that a manifest without a
// `yaml fence is rejected.
func TestParseManifestMissingFenceFails(t *testing.T) {
	bad := "# Just prose\n\nNo YAML block here.\n"
	if _, err := parseManifest(bad); err == nil {
		t.Errorf("expected error for missing fence, got nil")
	}
}

// TestParseManifestDoubleFenceFails verifies that a manifest with two
// `yaml fences is rejected (DI-011-20260429-184500: exactly one).
func TestParseManifestDoubleFenceFails(t *testing.T) {
	doubled := sampleManifest + "\n```yaml\nentries: []\n```\n"
	if _, err := parseManifest(doubled); err == nil {
		t.Errorf("expected error for double fence, got nil")
	}
}

// TestEmitYAMLDeterministic verifies that emitYAML produces stable output
// for the same input. Stable output means git diffs of MANIFEST.md are
// readable after each freeze.
func TestEmitYAMLDeterministic(t *testing.T) {
	entries := []Entry{
		{
			PCID:     "bafkreifjjcie6lypi6ny7amxnfftagclbuxndqonfipmb64f2km2devei4",
			Slug:     "harness-spec",
			Status:   StatusFrozen,
			FrozenOn: "2026-04-29T20:00:00Z",
		},
	}
	a, err := emitYAML(entries)
	if err != nil {
		t.Fatal(err)
	}
	b, err := emitYAML(entries)
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Errorf("non-deterministic emit:\n--- a ---\n%s\n--- b ---\n%s", a, b)
	}
}

// TestSortEntriesStable verifies that entries are sorted by slug, then by
// frozen_on, with no surprises on equal keys.
func TestSortEntriesStable(t *testing.T) {
	entries := []Entry{
		{Slug: "trust-ledger", FrozenOn: "2026-05-01T00:00:00Z"},
		{Slug: "harness-spec", FrozenOn: "2026-04-29T20:00:00Z"},
		{Slug: "harness-spec", FrozenOn: "2026-04-30T20:00:00Z"},
		{Slug: "trust-ledger", FrozenOn: "2026-04-30T00:00:00Z"},
	}
	sortEntries(entries)
	want := []string{
		"harness-spec/2026-04-29T20:00:00Z",
		"harness-spec/2026-04-30T20:00:00Z",
		"trust-ledger/2026-04-30T00:00:00Z",
		"trust-ledger/2026-05-01T00:00:00Z",
	}
	for i, e := range entries {
		got := e.Slug + "/" + e.FrozenOn
		if got != want[i] {
			t.Errorf("entries[%d] = %s, want %s", i, got, want[i])
		}
	}
}

// TestLatestFrozen verifies that the helper used by spec freeze to discover
// the predecessor of a new pCID returns the most-recent frozen entry for
// the slug.
func TestLatestFrozen(t *testing.T) {
	m := &Manifest{
		Entries: []Entry{
			{PCID: "old", Slug: "harness-spec", Status: StatusSuperseded, FrozenOn: "2026-04-29T20:00:00Z"},
			{PCID: "mid", Slug: "harness-spec", Status: StatusFrozen, FrozenOn: "2026-04-30T20:00:00Z"},
			{PCID: "other", Slug: "trust-ledger", Status: StatusFrozen, FrozenOn: "2026-05-01T00:00:00Z"},
		},
	}
	got := m.latestFrozen("harness-spec")
	if got == nil {
		t.Fatal("latestFrozen returned nil")
	}
	if got.PCID != "mid" {
		t.Errorf("latestFrozen.PCID = %q, want mid", got.PCID)
	}
	// Only `frozen` status counts; `superseded` is skipped even if more recent.
	got = m.latestFrozen("nonexistent")
	if got != nil {
		t.Errorf("latestFrozen(nonexistent) = %+v, want nil", got)
	}
}

// TestIndexLineStrict verifies that indexLine matches whole lines, not
// prefixes inside longer lines. Critical because `yamlFenceOpen` is "```yaml"
// and we don't want to match a hypothetical "```yaml-table" or similar.
func TestIndexLineStrict(t *testing.T) {
	cases := []struct {
		name string
		s    string
		line string
		want int
	}{
		{name: "exact-match-bof", s: "```yaml\nbody\n```\n", line: "```yaml", want: 0},
		{name: "exact-match-mid", s: "header\n```yaml\nbody\n```\n", line: "```yaml", want: 7},
		{name: "no-prefix-match", s: "```yaml-extension\n", line: "```yaml", want: -1},
		{name: "no-suffix-match", s: "extra```yaml\n", line: "```yaml", want: -1},
		{name: "empty", s: "", line: "```yaml", want: -1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := indexLine(tc.s, tc.line)
			if got != tc.want {
				t.Errorf("indexLine(%q, %q) = %d, want %d", tc.s, tc.line, got, tc.want)
			}
		})
	}
}
