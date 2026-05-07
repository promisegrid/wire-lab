package main

// Proquint encoding per Wilkerson 2009, "Proquints: Identifiers that are
// Readable, Spellable, and Pronounceable" (https://arxiv.org/html/0901.4016).
//
// A proquint is a sequence of consonant-vowel quintuplets where each quint
// encodes 16 bits: four bits each from four consonants and two bits each
// from two vowels, in the pattern C V C V C.
//
//   consonants (16): b d f g h j k l m n p r s t v z
//   vowels      (4): a i o u
//
// A 16-bit unsigned value n encodes to:
//   c1 = CONS[(n >> 12) & 0xF]
//   v1 = VOWS[(n >> 10) & 0x3]
//   c2 = CONS[(n >>  6) & 0xF]
//   v2 = VOWS[(n >>  4) & 0x3]
//   c3 = CONS[ n        & 0xF]
//
// proquint-1 is one quint (5 chars, 16-bit space, 65,536 values).
// proquint-2 is two quints joined by '-' (10 chars + hyphen, 32-bit space,
// 4.29 billion values).
//
// We use proquint as the wire-lab handle encoding because it is short,
// pronounceable, and trivially derivable from any 16- or 32-bit hash output
// without a curated wordlist or central registry. See TE-vapoj (sub. agnostic
// layered model) for the broader rationale and TE-39 for the lock decision.

const (
	proquintCons = "bdfghjklmnprstvz"
	proquintVows = "aiou"
)

// uint16ToProquint encodes a 16-bit value as a 5-character proquint string.
func uint16ToProquint(n uint16) string {
	out := []byte{
		proquintCons[(n>>12)&0xF],
		proquintVows[(n>>10)&0x3],
		proquintCons[(n>>6)&0xF],
		proquintVows[(n>>4)&0x3],
		proquintCons[n&0xF],
	}
	return string(out)
}

// proquint1FromBytes returns the proquint-1 (5 chars) derived from the first
// two bytes of b. Caller must ensure len(b) >= 2.
func proquint1FromBytes(b []byte) string {
	n := uint16(b[0])<<8 | uint16(b[1])
	return uint16ToProquint(n)
}

// proquint2FromBytes returns the proquint-2 (5 + '-' + 5 = 11 chars) derived
// from the first four bytes of b. Caller must ensure len(b) >= 4.
func proquint2FromBytes(b []byte) string {
	q1 := proquint1FromBytes(b[0:2])
	q2 := proquint1FromBytes(b[2:4])
	return q1 + "-" + q2
}
