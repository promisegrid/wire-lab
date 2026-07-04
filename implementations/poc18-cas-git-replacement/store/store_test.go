package store

import (
	"bytes"
	"os"
	"testing"
)

func TestFileStoreRoundTripUsesBase32CIDPaths(t *testing.T) {
	cas, openErr := Open(t.TempDir())
	if openErr != nil {
		t.Fatalf("Open() error = %v", openErr)
	}
	entry, putErr := cas.Put("message", []byte("hello"))
	if putErr != nil {
		t.Fatalf("Put() error = %v", putErr)
	}
	if entry.CID == "" || entry.CID[0] != 'b' {
		t.Fatalf("CID = %q, want CIDv1 base32 text", entry.CID)
	}
	if _, statErr := os.Stat(entry.Path); statErr != nil {
		t.Fatalf("stored path missing: %v", statErr)
	}
	parsedCID, parseErr := ParseCIDText(entry.CID)
	if parseErr != nil {
		t.Fatalf("ParseCIDText() error = %v", parseErr)
	}
	content, gotEntry, getErr := cas.Get(parsedCID)
	if getErr != nil {
		t.Fatalf("Get() error = %v", getErr)
	}
	if !bytes.Equal(content, []byte("hello")) {
		t.Fatalf("content = %q, want hello", string(content))
	}
	if gotEntry.CID != entry.CID {
		t.Fatalf("entry CID = %s, want %s", gotEntry.CID, entry.CID)
	}
}

func TestLinkTagRoundTrip(t *testing.T) {
	value := CIDForBytes([]byte("link target"))
	decoded, decodeErr := CIDFromLinkTag(LinkTag(value))
	if decodeErr != nil {
		t.Fatalf("CIDFromLinkTag() error = %v", decodeErr)
	}
	if !decoded.Equals(value) {
		t.Fatalf("decoded CID = %s, want %s", CIDText(decoded), CIDText(value))
	}
}

func TestFileStoreDeleteRemovesLocalEntryOnly(t *testing.T) {
	cas, openErr := Open(t.TempDir())
	if openErr != nil {
		t.Fatalf("Open() error = %v", openErr)
	}
	entry, putErr := cas.Put("chunk", []byte("pressure bytes"))
	if putErr != nil {
		t.Fatalf("Put() error = %v", putErr)
	}
	objectCID, parseErr := ParseCIDText(entry.CID)
	if parseErr != nil {
		t.Fatalf("ParseCIDText() error = %v", parseErr)
	}
	deleted, deleteErr := cas.Delete(objectCID)
	if deleteErr != nil {
		t.Fatalf("Delete() error = %v", deleteErr)
	}
	if deleted.CID != entry.CID {
		t.Fatalf("deleted CID = %s, want %s", deleted.CID, entry.CID)
	}
	if _, statErr := os.Stat(entry.Path); !os.IsNotExist(statErr) {
		t.Fatalf("deleted path stat error = %v, want not exist", statErr)
	}
	if cas.Has(objectCID) {
		t.Fatalf("Has(%s) = true after delete", entry.CID)
	}
}
