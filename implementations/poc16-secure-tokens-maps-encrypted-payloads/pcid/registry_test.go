package pcid

import (
	"os"
	"path/filepath"
	"testing"

	specdocs "promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/docs/protocols"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
)

func TestRegistryNamesRoundTrip(t *testing.T) {
	registry := NewRegistry()
	for _, name := range registry.Names() {
		protocolCID, ok := registry.CID(name)
		if !ok {
			t.Fatalf("missing protocol %s", name)
		}
		roundTripName, ok := registry.Name(protocolCID)
		if !ok || roundTripName != name {
			t.Fatalf("round trip name = %q %v, want %q true", roundTripName, ok, name)
		}
	}
}

func TestRegistryPCIDsMatchSpecBytesAndSymlinks(t *testing.T) {
	registry := NewRegistry()
	for _, name := range registry.Names() {
		specBytes, bytesErr := specdocs.BytesFor(name)
		if bytesErr != nil {
			t.Fatalf("%s spec bytes: %v", name, bytesErr)
		}
		protocolCID := registry.MustCID(name)
		if got, want := protocolCID.String(), protocol.CIDForExactBytes(specBytes); got != want {
			t.Fatalf("%s hardcoded pCID = %s, want content CID %s", name, got, want)
		}
		fileName, ok := specdocs.FileFor(name)
		if !ok {
			t.Fatalf("%s missing canonical spec file", name)
		}
		linkPath := filepath.Join("..", "docs", "protocols", protocolCID.String()+".md")
		linkInfo, statErr := os.Lstat(linkPath)
		if statErr != nil {
			t.Fatalf("%s symlink stat: %v", linkPath, statErr)
		}
		if linkInfo.Mode()&os.ModeSymlink == 0 {
			t.Fatalf("%s must be a symlink to %s", linkPath, fileName)
		}
		target, readlinkErr := os.Readlink(linkPath)
		if readlinkErr != nil {
			t.Fatalf("%s readlink: %v", linkPath, readlinkErr)
		}
		if target != fileName {
			t.Fatalf("%s target = %s, want %s", linkPath, target, fileName)
		}
		targetBytes, readErr := os.ReadFile(linkPath)
		if readErr != nil {
			t.Fatalf("%s read target: %v", linkPath, readErr)
		}
		if string(targetBytes) != string(specBytes) {
			t.Fatalf("%s target bytes do not match embedded spec bytes", linkPath)
		}
	}
}
