package main

import (
	"strings"
	"testing"
)

func intReps() []rep {
	return []rep{
		{"TE-38", "TE-vunub"},
		{"TE-3", "TE-pubum"},
		{"TODO-22", "TODO-vunub"},
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
