package lib

import (
	"crypto/sha256"
	"encoding/hex"
)

// ProtocolCID is the proof-of-concept Protocol CID for a spec document. It is a
// small CIDv1 raw-codec sha2-256 implementation for POC identity only.
//
// Intent: Keep pCID as Protocol CID and make exact protocol spec bytes select
// the message interpretation. Source: DI-horak.
type ProtocolCID struct {
	cidBytes []byte
	digest   [32]byte
}

// NewProtocolCID hashes a protocol spec and wraps the digest as CIDv1 raw-codec
// sha2-256 bytes: cid-version=1, multicodec=raw, multihash=sha2-256/32.
func NewProtocolCID(specBytes []byte) ProtocolCID {
	digest := sha256.Sum256(specBytes)
	cidBytes := make([]byte, 0, 36)
	cidBytes = append(cidBytes, 0x01)
	cidBytes = append(cidBytes, 0x55)
	cidBytes = append(cidBytes, 0x12)
	cidBytes = append(cidBytes, 0x20)
	cidBytes = append(cidBytes, digest[:]...)
	return ProtocolCID{cidBytes: cidBytes, digest: digest}
}

// NewProtocolCIDFromBytes rebuilds the POC CID value from binary CIDv1 bytes.
func NewProtocolCIDFromBytes(cidBytes []byte) ProtocolCID {
	copiedBytes := make([]byte, len(cidBytes))
	copy(copiedBytes, cidBytes)
	var digest [32]byte
	if len(copiedBytes) >= 36 {
		copy(digest[:], copiedBytes[len(copiedBytes)-32:])
	}
	return ProtocolCID{cidBytes: copiedBytes, digest: digest}
}

// Bytes returns the binary CIDv1 bytes without the DAG-CBOR tag-42 sentinel.
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

// String renders a stable debug form for evidence logs and receive promises.
func (protocolCID ProtocolCID) String() string {
	return "cidv1-raw-sha2-256:" + hex.EncodeToString(protocolCID.digest[:])
}

// Equal reports whether two Protocol CIDs identify the same spec bytes.
func (protocolCID ProtocolCID) Equal(other ProtocolCID) bool {
	return string(protocolCID.cidBytes) == string(other.cidBytes)
}
