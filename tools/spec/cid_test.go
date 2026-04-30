package main

import "testing"

// TestComputeCIDKnownVectors verifies that computeCID matches the canonical
// CIDv1 (base32 / sha2-256 / raw) for inputs whose expected CIDs are
// independently known.
//
// The "hello world\n" vector was independently produced by
// `github.com/ipfs/go-cid` and Python `py-multiformats-cid` during the
// TE-22 design review (chat 2026-04-29) and is the load-bearing test that
// pins the CIDv1 parameter set.
//
// The "" (empty) vector pins behavior on edge inputs.
//
// The "abc" vector is included because it cross-checks against the sha-256
// of "abc" being ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad,
// which is one of the most-cited sha-256 test vectors anywhere; if our CID
// disagrees, the multihash wrap is wrong.
func TestComputeCIDKnownVectors(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "hello-world-newline",
			in:   "hello world\n",
			want: "bafkreifjjcie6lypi6ny7amxnfftagclbuxndqonfipmb64f2km2devei4",
		},
		{
			name: "empty",
			in:   "",
			want: "bafkreihdwdcefgh4dqkjv67uzcmw7ojee6xedzdetojuzjevtenxquvyku",
		},
		{
			name: "abc",
			in:   "abc",
			want: "bafkreif2pall7dybz7vecqka3zo24irdwabwdi4wc55jznaq75q7eaavvu",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := computeCID([]byte(tc.in))
			if err != nil {
				t.Fatalf("computeCID(%q): %v", tc.in, err)
			}
			if got != tc.want {
				t.Errorf("computeCID(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// TestComputeCIDDeterminism verifies that hashing the same input twice yields
// the same CID. This is a property test guarding against future refactors
// that might accidentally introduce non-determinism (e.g., random salt,
// timestamp inclusion).
func TestComputeCIDDeterminism(t *testing.T) {
	in := []byte("the quick brown fox jumps over the lazy dog\n")
	a, err := computeCID(in)
	if err != nil {
		t.Fatal(err)
	}
	b, err := computeCID(in)
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Errorf("non-deterministic: %q != %q", a, b)
	}
}

// TestComputeCIDDistinguishes verifies that a one-byte change in the input
// produces a different CID. This guards against accidental input pre-processing
// (normalization, trimming) that would silently make different inputs hash the
// same.
func TestComputeCIDDistinguishes(t *testing.T) {
	a, _ := computeCID([]byte("abc"))
	b, _ := computeCID([]byte("abd"))
	if a == b {
		t.Errorf("collisions on single-byte change: both %q", a)
	}
}

// TestComputeCIDRespectsTrailingNewline verifies that adding or removing a
// trailing newline changes the CID. This is the load-bearing property of
// DI-011-20260429-184457 (raw bytes, no normalization): if the freeze tool
// silently stripped or added a trailing newline, the pCID would not honor
// "what you see is what got hashed."
func TestComputeCIDRespectsTrailingNewline(t *testing.T) {
	a, _ := computeCID([]byte("hello"))
	b, _ := computeCID([]byte("hello\n"))
	if a == b {
		t.Errorf("trailing newline ignored; raw-bytes hash should differ but both = %q", a)
	}
}

// TestComputeCIDRespectsCRLF verifies that LF and CRLF produce different
// CIDs. Same rationale as above: under DI-011-20260429-184457, a CRLF flip
// MUST change the pCID. An audit warning may flag CRLF as suspicious, but
// the hash itself does not normalize.
func TestComputeCIDRespectsCRLF(t *testing.T) {
	a, _ := computeCID([]byte("a\nb\n"))
	b, _ := computeCID([]byte("a\r\nb\r\n"))
	if a == b {
		t.Errorf("CRLF/LF normalization detected; raw-bytes hash should differ but both = %q", a)
	}
}
