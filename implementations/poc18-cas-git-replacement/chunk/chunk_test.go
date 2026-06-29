package chunk

import (
	"bytes"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

func TestStoreBytesRabinChunksAndReassembles(t *testing.T) {
	cas, openErr := store.Open(t.TempDir())
	if openErr != nil {
		t.Fatalf("Open() error = %v", openErr)
	}
	content := bytes.Repeat([]byte("abcdef0123456789"), 8192)
	stored, storeErr := StoreBytes(cas, content)
	if storeErr != nil {
		t.Fatalf("StoreBytes() error = %v", storeErr)
	}
	if len(stored.Manifest.Chunks) < 2 {
		t.Fatalf("chunk count = %d, want multiple chunks", len(stored.Manifest.Chunks))
	}
	manifestBytes, _, getErr := cas.Get(stored.ManifestCID)
	if getErr != nil {
		t.Fatalf("Get(manifest) error = %v", getErr)
	}
	manifest, decodeErr := DecodeManifest(manifestBytes)
	if decodeErr != nil {
		t.Fatalf("DecodeManifest() error = %v", decodeErr)
	}
	reassembled, reassembleErr := Reassemble(cas, manifest)
	if reassembleErr != nil {
		t.Fatalf("Reassemble() error = %v", reassembleErr)
	}
	if !bytes.Equal(reassembled, content) {
		t.Fatalf("reassembled content mismatch")
	}
}
