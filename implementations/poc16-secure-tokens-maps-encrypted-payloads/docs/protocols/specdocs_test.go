package specdocs

import (
	"strings"
	"testing"
)

// TestEmbeddedSpecsAreRFCComplete guards the pCID input documents from drifting
// back into short notes.
// Intent: Every embedded specdoc is hashed into its pCID and supplied to
// LLM-backed agents as implementation context, so missing protocol sections are
// runtime-relevant regressions, not merely documentation style issues. Source:
// DI-bitug
func TestEmbeddedSpecsAreRFCComplete(t *testing.T) {
	requiredSections := []string{
		"## Status",
		"## Abstract",
		"## pCID and envelope",
		"## Promise Theory model",
		"## Payload grammar",
		"## Sender behavior",
		"## Receiver and parser behavior",
		"## Protocol state machine",
		"## State, CAS, DAG, and retention",
		"## Security considerations",
		"## Interoperability notes",
		"## Examples",
	}
	for protocolName, fileName := range protocolFiles {
		specBytes, readErr := docs.ReadFile(fileName)
		if readErr != nil {
			t.Fatalf("%s read %s: %v", protocolName, fileName, readErr)
		}
		specText := string(specBytes)
		for _, requiredSection := range requiredSections {
			if !strings.Contains(specText, requiredSection) {
				t.Errorf("%s missing required section %q", fileName, requiredSection)
			}
		}
		if !strings.Contains(specText, "grid([42(pCID)") {
			t.Errorf("%s must spell out its grid([42(pCID), ...]) envelope shape", fileName)
		}
		if !strings.Contains(specText, "Source:") || !strings.Contains(specText, "`DI-") {
			t.Errorf("%s must cite DI provenance for normative spec text", fileName)
		}
		forbiddenPrefix := "field" + "_"
		if strings.Contains(specText, forbiddenPrefix) {
			t.Errorf("%s must not use legacy payload prefix vocabulary", fileName)
		}
		if strings.Count(specText, "\n") < 80 {
			t.Errorf("%s is too short to be an RFC-style protocol spec", fileName)
		}
	}
}
