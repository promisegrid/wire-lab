package protocol

import (
	"fmt"

	cidlib "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

// CIDForBytes returns a CIDv1 raw sha2-256 text identifier for exact bytes.
func CIDForBytes(data []byte) (string, error) {
	hash, err := mh.Sum(data, mh.SHA2_256, -1)
	if err != nil {
		return "", fmt.Errorf("multihash bytes: %w", err)
	}
	return cidlib.NewCidV1(cidlib.Raw, hash).String(), nil
}

// ValidateCIDText verifies canonical CIDv1 base32 text used in POC17 outputs.
func ValidateCIDText(text string) error {
	parsed, err := cidlib.Decode(text)
	if err != nil {
		return fmt.Errorf("decode CID: %w", err)
	}
	prefix := parsed.Prefix()
	if parsed.Version() != 1 || prefix.Codec != cidlib.Raw || prefix.MhType != mh.SHA2_256 {
		return fmt.Errorf("CID must be CIDv1 raw sha2-256")
	}
	if parsed.String() != text {
		return fmt.Errorf("CID text must be canonical base32")
	}
	return nil
}

// Tag42DataForCIDText returns DAG-CBOR tag-42 bytes for one CID text value.
func Tag42DataForCIDText(text string) ([]byte, error) {
	parsed, err := cidlib.Decode(text)
	if err != nil {
		return nil, fmt.Errorf("decode CID: %w", err)
	}
	data := make([]byte, 0, len(parsed.Bytes())+1)
	data = append(data, 0x00)
	data = append(data, parsed.Bytes()...)
	return data, nil
}

// CIDTextFromTag42Data returns canonical text for DAG-CBOR tag-42 CID bytes.
func CIDTextFromTag42Data(data []byte) (string, error) {
	if len(data) < 2 || data[0] != 0x00 {
		return "", fmt.Errorf("slot 0 tag 42 data must be DAG-CBOR CID bytes")
	}
	parsed, err := cidlib.Cast(data[1:])
	if err != nil {
		return "", fmt.Errorf("parse tag 42 CID: %w", err)
	}
	return parsed.String(), nil
}
