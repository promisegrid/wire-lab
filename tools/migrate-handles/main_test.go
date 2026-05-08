package main

import (
	"strings"
	"testing"
)

// TestInjectAfterTEID verifies the Prior aliases block lands immediately
// after the existing "## TE ID" section and before the next "## " heading.
func TestInjectAfterTEID(t *testing.T) {
	body := "# TE-1: Promise-stack ordering\n\n" +
		"*Thought experiment, part of...*\n\n" +
		"## TE ID\n\n" +
		"TE-20260427-180000\n\n" +
		"## Status\n\n" +
		"needs DF\n"
	got := injectPriorAliases(body, "TE", "TE-1", "TE-20260427-180000")
	if !strings.Contains(got, "## Prior aliases") {
		t.Fatalf("missing Prior aliases section:\n%s", got)
	}
	if !strings.Contains(got, "TE-1") || !strings.Contains(got, "TE-20260427-180000") {
		t.Errorf("missing alias values in output:\n%s", got)
	}
	idID := strings.Index(got, "## TE ID")
	idPrior := strings.Index(got, "## Prior aliases")
	idStatus := strings.Index(got, "## Status")
	if !(idID < idPrior && idPrior < idStatus) {
		t.Errorf("ordering wrong: TE ID=%d, Prior=%d, Status=%d\n%s",
			idID, idPrior, idStatus, got)
	}
}

// TestInjectAfterTODOID covers the TODO ID heading variant.
func TestInjectAfterTODOID(t *testing.T) {
	body := "# Thought work item - foo\n\n" +
		"## TODO ID\n\n" +
		"TODO-20260507-002306\n\n" +
		"## Status\n\nopen\n"
	got := injectPriorAliases(body, "TODO", "TODO-23", "TODO-20260507-002306")
	if !strings.Contains(got, "TODO-23") {
		t.Errorf("missing integer alias TODO-23:\n%s", got)
	}
	idID := strings.Index(got, "## TODO ID")
	idPrior := strings.Index(got, "## Prior aliases")
	idStatus := strings.Index(got, "## Status")
	if !(idID < idPrior && idPrior < idStatus) {
		t.Errorf("ordering wrong:\n%s", got)
	}
}

// TestInjectNoIntegerAlias covers the case where the file has no integer
// alias (e.g. some TEs in the index without a TE-N number).
func TestInjectNoIntegerAlias(t *testing.T) {
	body := "# TE - foo\n\n## TE ID\n\nTE-20260429-033208\n\n## Body\n\n..."
	got := injectPriorAliases(body, "TE", "", "TE-20260429-033208")
	if !strings.Contains(got, "no integer alias was assigned") {
		t.Errorf("expected no-integer-alias note:\n%s", got)
	}
	if strings.Contains(got, "(integer alias)") {
		t.Errorf("should not have integer-alias bullet:\n%s", got)
	}
}

// TestInjectFallbackToH1 covers files missing a "## TE ID" / "## TODO ID"
// section. The block is inserted after the H1 instead.
func TestInjectFallbackToH1(t *testing.T) {
	body := "# TE - foo\n\nSome prose.\n\n## Body\n\n..."
	got := injectPriorAliases(body, "TE", "TE-7", "TE-20260427-180600")
	idH1 := strings.Index(got, "# TE - foo")
	idPrior := strings.Index(got, "## Prior aliases")
	idBody := strings.Index(got, "## Body")
	if !(idH1 < idPrior && idPrior < idBody) {
		t.Errorf("ordering wrong:\n%s", got)
	}
}

// TestProquintVectors guards against the two-copy proquint table drifting
// out of sync with mint-handle's table.
func TestProquintVectors(t *testing.T) {
	cases := []struct {
		n    uint16
		want string
	}{
		{0x0000, "babab"},
		{0xFFFF, "zuzuz"},
	}
	for _, c := range cases {
		if got := uint16ToProquint(c.n); got != c.want {
			t.Errorf("uint16ToProquint(%#x) = %q, want %q", c.n, got, c.want)
		}
	}
}
