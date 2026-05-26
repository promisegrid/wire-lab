package lib

import (
	"crypto/sha256"
	"encoding/hex"
)

// ProtocolCID identifies the protocol spec document that defines payload and
// outer-slot semantics. It is not a payload CID.
//
// Intent: Preserve the pCID-as-protocol-selector rule across poc5's relay and
// app protocols. Source: DI-rarim.
type ProtocolCID struct {
	cidBytes []byte
	digest   [32]byte
}

// NewProtocolCID hashes spec bytes into the POC CIDv1 raw sha2-256 form.
func NewProtocolCID(specBytes []byte) ProtocolCID {
	digest := sha256.Sum256(specBytes)
	cidBytes := make([]byte, 0, 36)
	cidBytes = append(cidBytes, 0x01, 0x55, 0x12, 0x20)
	cidBytes = append(cidBytes, digest[:]...)
	return ProtocolCID{cidBytes: cidBytes, digest: digest}
}

// NewProtocolCIDFromBytes rebuilds a POC CID from binary CIDv1 bytes.
func NewProtocolCIDFromBytes(cidBytes []byte) ProtocolCID {
	copiedBytes := make([]byte, len(cidBytes))
	copy(copiedBytes, cidBytes)
	var digest [32]byte
	if len(copiedBytes) >= 36 {
		copy(digest[:], copiedBytes[len(copiedBytes)-32:])
	}
	return ProtocolCID{cidBytes: copiedBytes, digest: digest}
}

// Bytes returns binary CIDv1 bytes without the tag-42 sentinel.
func (protocolCID ProtocolCID) Bytes() []byte {
	cidBytes := make([]byte, len(protocolCID.cidBytes))
	copy(cidBytes, protocolCID.cidBytes)
	return cidBytes
}

// Tag42Bytes returns the DAG-CBOR tag-42 byte string payload.
func (protocolCID ProtocolCID) Tag42Bytes() []byte {
	tagBytes := make([]byte, 0, len(protocolCID.cidBytes)+1)
	tagBytes = append(tagBytes, 0x00)
	tagBytes = append(tagBytes, protocolCID.cidBytes...)
	return tagBytes
}

// String renders a stable debug and promise-record form.
func (protocolCID ProtocolCID) String() string {
	return "cidv1-raw-sha2-256:" + hex.EncodeToString(protocolCID.digest[:])
}

// Equal reports whether two pCIDs name the same protocol bytes.
func (protocolCID ProtocolCID) Equal(other ProtocolCID) bool {
	return string(protocolCID.cidBytes) == string(other.cidBytes)
}
