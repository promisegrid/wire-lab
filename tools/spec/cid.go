package main

import (
	"fmt"
	"os"

	cidlib "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

// computeCID computes the CIDv1 of the given byte slice using the locked
// parameter set: multibase=base32, multihash=sha2-256, codec=raw.
//
// The returned string is the canonical text form of the CID, as produced by
// (cid.Cid).String(): a base32 multibase prefix 'b' followed by the lowercase
// base32 (RFC 4648, no padding) encoding of the binary CID.
//
// The hash is computed over the input bytes verbatim; no normalization,
// formatting, or text canonicalization is applied (per DI-011-20260429-184457).
func computeCID(data []byte) (string, error) {
	mhBuf, err := mh.Sum(data, mh.SHA2_256, -1)
	if err != nil {
		return "", fmt.Errorf("multihash: %w", err)
	}
	c := cidlib.NewCidV1(cidlib.Raw, mhBuf)
	return c.String(), nil
}

// computeFileCID reads a file and returns its CIDv1.
func computeFileCID(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return computeCID(data)
}

func cmdCid(file string) error {
	c, err := computeFileCID(file)
	if err != nil {
		return err
	}
	fmt.Println(c)
	return nil
}
