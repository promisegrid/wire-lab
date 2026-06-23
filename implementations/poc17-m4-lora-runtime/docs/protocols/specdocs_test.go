package specdocs

import (
	"strings"
	"testing"
)

func TestEmbeddedSpecsHaveRequiredShape(t *testing.T) {
	required := []string{
		"## Status",
		"## Abstract",
		"## pCID and envelope",
		"## Payload grammar",
		"## Sender behavior",
		"## Receiver behavior",
		"## Examples",
		"Source: `DI-",
		"grid([42(pCID), payload])",
	}
	for _, name := range Names() {
		specBytes, err := BytesFor(name)
		if err != nil {
			t.Fatal(err)
		}
		text := string(specBytes)
		for _, section := range required {
			if !strings.Contains(text, section) {
				t.Errorf("%s missing %q", name, section)
			}
		}
	}
}

func TestRegistryDerivesStableCIDText(t *testing.T) {
	registry := MustRegistry()
	for _, name := range Names() {
		cid := registry.MustCID(name)
		if !strings.HasPrefix(cid.String(), "bafkrei") {
			t.Fatalf("%s got non-raw-sha2-256 CID %s", name, cid.String())
		}
		if len(cid.Bytes()) != 36 {
			t.Fatalf("%s binary CID length = %d", name, len(cid.Bytes()))
		}
		if len(cid.Tag42Data()) != 37 || cid.Tag42Data()[0] != 0x00 {
			t.Fatalf("%s invalid tag42 CID data", name)
		}
		roundTripName, ok := registry.NameForCID(cid.String())
		if !ok || roundTripName != name {
			t.Fatalf("%s failed CID name round trip", name)
		}
	}
}
