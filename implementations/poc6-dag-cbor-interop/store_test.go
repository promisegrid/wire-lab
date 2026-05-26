package poc6

import (
	"bytes"
	"testing"
)

func TestDAGCBORPointerRoundTrip(t *testing.T) {
	store := NewObjectStore()
	childNode, childErr := BuildMerkleNode("child", []byte("hello dag-cbor"))
	if childErr != nil {
		t.Fatalf("build child node: %v", childErr)
	}
	childObject, storeChildErr := store.StoreObject("Alice", childNode, "Alice promises to store child bytes exactly.")
	if storeChildErr != nil {
		t.Fatalf("store child: %v", storeChildErr)
	}

	pointerNode, pointerErr := BuildPointerObject("pointer", childObject.CID, []byte("Alice observed child storage."))
	if pointerErr != nil {
		t.Fatalf("build pointer: %v", pointerErr)
	}
	pointerObject, storePointerErr := store.StoreObject("Alice", pointerNode, "Alice promises to store pointer bytes exactly.")
	if storePointerErr != nil {
		t.Fatalf("store pointer: %v", storePointerErr)
	}

	decodedPointer, _, loadErr := store.LoadObject("Bob", pointerObject.CID, "Bob locally verifies pointer bytes by CID.")
	if loadErr != nil {
		t.Fatalf("load pointer: %v", loadErr)
	}
	targetCID, targetErr := VerifyLinkTarget(decodedPointer)
	if targetErr != nil {
		t.Fatalf("verify target: %v", targetErr)
	}
	if targetCID.String() != childObject.CID.String() {
		t.Fatalf("target cid mismatch: got %s want %s", targetCID, childObject.CID)
	}
}

func TestDAGCBORBytesContainTag42Link(t *testing.T) {
	store := NewObjectStore()
	childNode, childErr := BuildMerkleNode("child", []byte("link target"))
	if childErr != nil {
		t.Fatalf("build child node: %v", childErr)
	}
	childObject, storeChildErr := store.StoreObject("Alice", childNode, "Alice promises to store child bytes exactly.")
	if storeChildErr != nil {
		t.Fatalf("store child: %v", storeChildErr)
	}
	pointerNode, pointerErr := BuildPointerObject("pointer", childObject.CID, []byte("evidence"))
	if pointerErr != nil {
		t.Fatalf("build pointer: %v", pointerErr)
	}
	pointerBytes, encodeErr := EncodeDAGCBOR(pointerNode)
	if encodeErr != nil {
		t.Fatalf("encode pointer: %v", encodeErr)
	}

	expectedTag42CIDBytes := append([]byte{0xd8, 0x2a, 0x58, byte(len(childObject.CID.Bytes()) + 1), 0x00}, childObject.CID.Bytes()...)
	if !bytes.Contains(pointerBytes, expectedTag42CIDBytes) {
		t.Fatalf("encoded pointer does not contain DAG-CBOR tag-42 CID bytes %x in %x", expectedTag42CIDBytes, pointerBytes)
	}
}

func TestCIDStableAcrossRepeatedEncoding(t *testing.T) {
	childNode, childErr := BuildMerkleNode("child", []byte("stable"))
	if childErr != nil {
		t.Fatalf("build child node: %v", childErr)
	}
	firstBytes, firstEncodeErr := EncodeDAGCBOR(childNode)
	if firstEncodeErr != nil {
		t.Fatalf("first encode: %v", firstEncodeErr)
	}
	secondBytes, secondEncodeErr := EncodeDAGCBOR(childNode)
	if secondEncodeErr != nil {
		t.Fatalf("second encode: %v", secondEncodeErr)
	}
	if !bytes.Equal(firstBytes, secondBytes) {
		t.Fatalf("encoded bytes changed: %x != %x", firstBytes, secondBytes)
	}

	firstCID, firstCIDErr := CIDForDAGCBOR(firstBytes)
	if firstCIDErr != nil {
		t.Fatalf("first cid: %v", firstCIDErr)
	}
	secondCID, secondCIDErr := CIDForDAGCBOR(secondBytes)
	if secondCIDErr != nil {
		t.Fatalf("second cid: %v", secondCIDErr)
	}
	if firstCID.String() != secondCID.String() {
		t.Fatalf("cid changed: %s != %s", firstCID, secondCID)
	}
}

func TestByteStringPayloadPreserved(t *testing.T) {
	payload := []byte{0x00, 0x01, 0x02, 0xfe, 0xff}
	childNode, childErr := BuildMerkleNode("bytes", payload)
	if childErr != nil {
		t.Fatalf("build child node: %v", childErr)
	}
	encodedBytes, encodeErr := EncodeDAGCBOR(childNode)
	if encodeErr != nil {
		t.Fatalf("encode child: %v", encodeErr)
	}
	decodedNode, decodeErr := DecodeDAGCBOR(encodedBytes)
	if decodeErr != nil {
		t.Fatalf("decode child: %v", decodeErr)
	}
	decodedPayload, lookupErr := LookupBytes(decodedNode, "payload")
	if lookupErr != nil {
		t.Fatalf("lookup payload: %v", lookupErr)
	}
	if !bytes.Equal(decodedPayload, payload) {
		t.Fatalf("payload changed: %x != %x", decodedPayload, payload)
	}
}

func TestLocalEvidenceRecords(t *testing.T) {
	store := NewObjectStore()
	childNode, childErr := BuildMerkleNode("child", []byte("evidence"))
	if childErr != nil {
		t.Fatalf("build child node: %v", childErr)
	}
	childObject, storeChildErr := store.StoreObject("Alice", childNode, "Alice promises to store child bytes exactly.")
	if storeChildErr != nil {
		t.Fatalf("store child: %v", storeChildErr)
	}
	if _, _, loadErr := store.LoadObject("Bob", childObject.CID, "Bob locally verifies child bytes by CID."); loadErr != nil {
		t.Fatalf("load child: %v", loadErr)
	}

	evidence := store.Evidence()
	if len(evidence) != 2 {
		t.Fatalf("evidence length got %d want 2: %#v", len(evidence), evidence)
	}
	if evidence[0].Actor != "Alice" || evidence[0].Event != "store" || evidence[0].Outcome != "kept" {
		t.Fatalf("unexpected Alice evidence: %#v", evidence[0])
	}
	if evidence[1].Actor != "Bob" || evidence[1].Event != "load" || evidence[1].Outcome != "kept" {
		t.Fatalf("unexpected Bob evidence: %#v", evidence[1])
	}
	if evidence[0].Detail == "" || evidence[1].Detail == "" {
		t.Fatalf("evidence details must not be empty: %#v", evidence)
	}
}

func TestMissingObjectRecordsNotPromised(t *testing.T) {
	store := NewObjectStore()
	childNode, childErr := BuildMerkleNode("missing", []byte("not stored here"))
	if childErr != nil {
		t.Fatalf("build child node: %v", childErr)
	}
	encodedBytes, encodeErr := EncodeDAGCBOR(childNode)
	if encodeErr != nil {
		t.Fatalf("encode child: %v", encodeErr)
	}
	objectCID, cidErr := CIDForDAGCBOR(encodedBytes)
	if cidErr != nil {
		t.Fatalf("cid child: %v", cidErr)
	}

	if _, _, loadErr := store.LoadObject("Bob", objectCID, "Bob asks for bytes Alice did not promise this store has."); loadErr == nil {
		t.Fatalf("load missing object unexpectedly succeeded")
	}
	evidence := store.Evidence()
	if len(evidence) != 1 {
		t.Fatalf("evidence length got %d want 1: %#v", len(evidence), evidence)
	}
	if evidence[0].Outcome != "not-promised" {
		t.Fatalf("missing object should be not-promised evidence: %#v", evidence[0])
	}
}
