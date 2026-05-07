package main

import (
	"regexp"
	"strings"
	"testing"
)

// makeIntRep mirrors main.go's pattern construction for tests.
func makeIntRep(kind, intNum, handle string) intRep {
	pattern := `(^|[^\w-])` + kind + `[\s-]+0*` + intNum + `($|[^\w-])`
	return intRep{re: regexp.MustCompile(pattern), to: kind + "-" + handle}
}

func intReps() []intRep {
	return []intRep{
		makeIntRep("TE", "38", "vunub"),
		makeIntRep("TE", "3", "pubum"),
		makeIntRep("TODO", "22", "vunub"),
		makeIntRep("TODO", "18", "jodon"),
		makeIntRep("TODO", "14", "vuhuj"),
		makeIntRep("TE", "25", "titur"),
	}
}

func tsReps() []rep {
	return []rep{
		{"TE-20260506-184800", "TE-vunub"},
	}
}

func pathReps() map[string]string {
	return map[string]string{
		"TE-20260506-184800-substrate-agnostic-layered-model.md": "TE-vunub-substrate-agnostic-layered-model.md",
	}
}

// TestPrefixSafety verifies TE-3 inside TE-38 is not over-matched.
func TestPrefixSafety(t *testing.T) {
	in := "See TE-38 for the layered model and TE-3 for currency."
	got, n := sweep(in, intReps(), nil, nil)
	if !strings.Contains(got, "TE-vunub for the layered") {
		t.Errorf("TE-38 not rewritten: %q", got)
	}
	if !strings.Contains(got, "TE-pubum for currency") {
		t.Errorf("TE-3 not rewritten: %q", got)
	}
	if strings.Contains(got, "TE-pubum8") {
		t.Errorf("TE-3 over-matched into TE-38: %q", got)
	}
	if n != 2 {
		t.Errorf("substitution count = %d, want 2", n)
	}
}

// TestPriorAliasesSkipped verifies the Prior aliases section passes
// through verbatim.
func TestPriorAliasesSkipped(t *testing.T) {
	in := "# TE-foo\n\n## Prior aliases\n\n- TE-38 (integer alias)\n\n## Body\n\nNow we cite TE-38 in prose.\n"
	got, n := sweep(in, intReps(), nil, nil)
	if !strings.Contains(got, "- TE-38 (integer alias)") {
		t.Errorf("Prior aliases section was rewritten: %q", got)
	}
	if !strings.Contains(got, "Now we cite TE-vunub in prose") {
		t.Errorf("body citation not rewritten: %q", got)
	}
	if n != 1 {
		t.Errorf("substitution count = %d, want 1 (only the body cite)", n)
	}
}

// TestTODOSpaceForm verifies "TODO 22" prose is rewritten.
func TestTODOSpaceForm(t *testing.T) {
	in := "Tracked under TODO 22 in the master list."
	got, _ := sweep(in, intReps(), nil, nil)
	if !strings.Contains(got, "TODO-vunub") {
		t.Errorf("TODO 22 not rewritten: %q", got)
	}
}

// TestPathRewrite verifies markdown link targets are rewritten.
func TestPathRewrite(t *testing.T) {
	in := "See [TE-38](TE-20260506-184800-substrate-agnostic-layered-model.md)."
	got, _ := sweep(in, intReps(), nil, pathReps())
	if !strings.Contains(got, "TE-vunub-substrate-agnostic-layered-model.md") {
		t.Errorf("path not rewritten: %q", got)
	}
}

// TestTimestampPriorAliasSkipped verifies the timestamp form inside the
// Prior aliases section is preserved (it is part of the historical
// record).
func TestTimestampPriorAliasSkipped(t *testing.T) {
	in := "## Prior aliases\n\n- `TE-20260506-184800` (timestamp alias)\n\n## Body\n\nThe TE was first drafted at TE-20260506-184800.\n"
	got, _ := sweep(in, nil, tsReps(), nil)
	if !strings.Contains(got, "`TE-20260506-184800` (timestamp alias)") {
		t.Errorf("Prior aliases timestamp was rewritten: %q", got)
	}
	if !strings.Contains(got, "first drafted at TE-vunub") {
		t.Errorf("body timestamp not rewritten: %q", got)
	}
}

// TestSpaceFormVariants verifies "TODO 18", "TODO  18" (double space),
// "TODO\t18", "TE 25" all rewrite correctly.
func TestSpaceFormVariants(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"under TODO 18 today", "under TODO-jodon today"},
		{"under TODO  18 today", "under TODO-jodon today"},
		{"under TODO\t18 today", "under TODO-jodon today"},
		{"per TE 25 section S5", "per TE-titur section S5"},
		{"per TE-25 section S5", "per TE-titur section S5"},
	}
	for _, c := range cases {
		got, _ := sweep(c.in, intReps(), nil, nil)
		if got != c.want {
			t.Errorf("input %q -> got %q, want %q", c.in, got, c.want)
		}
	}
}

// TestLeadingZeroForm verifies "TODO 011", "TODO-014", "TODO-018"
// rewrite correctly (the leading zero is stripped via 0* in the regex).
func TestLeadingZeroForm(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"TODO 014 was the migration", "TODO-vuhuj was the migration"},
		{"TODO-014 was the migration", "TODO-vuhuj was the migration"},
		{"TODO 018 v0 reference", "TODO-jodon v0 reference"},
		{"TODO-018 v0 reference", "TODO-jodon v0 reference"},
	}
	for _, c := range cases {
		got, _ := sweep(c.in, intReps(), nil, nil)
		if got != c.want {
			t.Errorf("input %q -> got %q, want %q", c.in, got, c.want)
		}
	}
}

// TestNoOvermatchOnLargerInteger verifies "TODO-180" is not matched as
// "TODO-18" (the trailing 0 is part of the integer, not a separator).
func TestNoOvermatchOnLargerInteger(t *testing.T) {
	cases := []string{
		"see TODO-180 for context",
		"see TODO 180 for context",
		"see TODO-1800 for context",
		"see TE-380 for context", // would over-match TE-38 if guard fails
	}
	for _, c := range cases {
		got, n := sweep(c, intReps(), nil, nil)
		if got != c {
			t.Errorf("input %q was rewritten to %q (should be unchanged)", c, got)
		}
		if n != 0 {
			t.Errorf("input %q got %d substitutions, want 0", c, n)
		}
	}
}

// TestNoMatchOnTimestampDigits verifies the integer rewriter does not
// match digits inside ISO timestamps (the timestamp rewriter is the only
// path that should touch those).
func TestNoMatchOnTimestampDigits(t *testing.T) {
	in := "drafted at 20260427-180000 then refined"
	got, n := sweep(in, intReps(), nil, nil)
	if got != in {
		t.Errorf("input %q was rewritten to %q", in, got)
	}
	if n != 0 {
		t.Errorf("got %d substitutions, want 0", n)
	}
}

// TestEdgeOfStringBoundary verifies the leading/trailing boundary
// captures handle start-of-string and end-of-string correctly.
func TestEdgeOfStringBoundary(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"TODO 18", "TODO-jodon"},
		{"TODO-18", "TODO-jodon"},
		{"TE-38", "TE-vunub"},
		{"see TE 25", "see TE-titur"},
	}
	for _, c := range cases {
		got, _ := sweep(c.in, intReps(), nil, nil)
		if got != c.want {
			t.Errorf("input %q -> got %q, want %q", c.in, got, c.want)
		}
	}
}
