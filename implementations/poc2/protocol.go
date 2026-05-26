package main

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
)

//go:embed protocols/hello.d/specs/hello-draft.md
var protocolSpecFS embed.FS

const helloSpecPath = "protocols/hello.d/specs/hello-draft.md"

// ProtocolCID is the proof-of-concept Protocol CID for a spec document. It is
// enough of CIDv1 to exercise the slot-0 PromiseGrid boundary without pulling in
// a full IPLD stack for this tiny POC.
//
// Intent: Keep pCID as Protocol CID and make the exact spec bytes determine the
// identifier used in grid slot 0. Source: DI-tijat.
type ProtocolCID struct {
	cidBytes []byte
	digest   [32]byte
}

// NewProtocolCID hashes a protocol spec and wraps the digest as CIDv1 raw-codec
// sha2-256 bytes: cid-version=1, multicodec=raw, multihash=sha2-256/32.
func NewProtocolCID(specBytes []byte) ProtocolCID {
	digest := sha256.Sum256(specBytes)
	cidBytes := make([]byte, 0, 36)
	cidBytes = append(cidBytes, 0x01) // CIDv1.
	cidBytes = append(cidBytes, 0x55) // multicodec raw.
	cidBytes = append(cidBytes, 0x12) // multihash sha2-256.
	cidBytes = append(cidBytes, 0x20) // 32-byte digest length.
	cidBytes = append(cidBytes, digest[:]...)
	return ProtocolCID{cidBytes: cidBytes, digest: digest}
}

// HelloProtocolCID returns the pCID for the local hello draft spec embedded in
// this implementation.
func HelloProtocolCID() (ProtocolCID, error) {
	specBytes, readErr := protocolSpecFS.ReadFile(helloSpecPath)
	if readErr != nil {
		return ProtocolCID{}, readErr
	}
	return NewProtocolCID(specBytes), nil
}

// Bytes returns the binary CIDv1 bytes, without the DAG-CBOR tag-42 leading 0x00
// link sentinel.
func (protocolCID ProtocolCID) Bytes() []byte {
	cidBytes := make([]byte, len(protocolCID.cidBytes))
	copy(cidBytes, protocolCID.cidBytes)
	return cidBytes
}

// Tag42Bytes returns the DAG-CBOR tag-42 byte string payload: a 0x00 sentinel
// followed by binary CID bytes.
func (protocolCID ProtocolCID) Tag42Bytes() []byte {
	tagBytes := make([]byte, 0, len(protocolCID.cidBytes)+1)
	tagBytes = append(tagBytes, 0x00)
	tagBytes = append(tagBytes, protocolCID.cidBytes...)
	return tagBytes
}

// String renders a stable debug form. It is intentionally not a full multibase
// CID string; the POC only needs deterministic logs and tests.
func (protocolCID ProtocolCID) String() string {
	return "cidv1-raw-sha2-256:" + hex.EncodeToString(protocolCID.digest[:])
}

// Equal reports whether two Protocol CIDs identify the same spec bytes.
func (protocolCID ProtocolCID) Equal(other ProtocolCID) bool {
	return string(protocolCID.cidBytes) == string(other.cidBytes)
}
