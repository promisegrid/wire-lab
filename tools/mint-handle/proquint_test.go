package main

import "testing"

// TestUint16ToProquintKnownVectors pins the proquint encoding to the
// canonical values from Wilkerson 2009 sec. 3, "Encoding Examples". These
// are the load-bearing vectors that lock the consonant/vowel ordering and
// bit-shift layout. If anyone ever changes the alphabets or bit packing,
// these tests will catch it.
func TestUint16ToProquintKnownVectors(t *testing.T) {
	cases := []struct {
		in   uint16
		want string
	}{
		// From the Wilkerson paper, table of IP-address example encodings.
		// IP 127.0.0.1 = 0x7F000001, which encodes as "lusab-babad".
		// We test each half independently:
		{0x7F00, "lusab"}, // first 16 bits of 0x7F000001
		{0x0001, "babad"}, // last  16 bits of 0x7F000001

		// IP 63.84.220.193 = 0x3F54DCC1 -> "gutih-tugad"
		{0x3F54, "gutih"},
		{0xDCC1, "tugad"},

		// Edge values
		{0x0000, "babab"}, // all zeros: cons[0]=b, vow[0]=a everywhere
		{0xFFFF, "zuzuz"}, // all ones:  cons[F]=z, vow[3]=u everywhere
	}
	for _, c := range cases {
		got := uint16ToProquint(c.in)
		if got != c.want {
			t.Errorf("uint16ToProquint(0x%04X) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestProquintLengths verifies the structural invariants: a proquint-1 is
// always exactly 5 characters and a proquint-2 is always exactly 11 (5 + '-'
// + 5). These are filename-shape contracts.
func TestProquintLengths(t *testing.T) {
	b := []byte{0x12, 0x34, 0x56, 0x78}
	if got := proquint1FromBytes(b); len(got) != 5 {
		t.Errorf("proquint1 length = %d, want 5 (%q)", len(got), got)
	}
	if got := proquint2FromBytes(b); len(got) != 11 {
		t.Errorf("proquint2 length = %d, want 11 (%q)", len(got), got)
	}
}

// TestProquintAlphabetMembership verifies every character in a generated
// proquint is in the documented consonant or vowel alphabet, in the right
// position (CVCVC). This guards against off-by-one bit-shift bugs.
func TestProquintAlphabetMembership(t *testing.T) {
	cons := proquintCons
	vows := proquintVows
	contains := func(s string, c byte) bool {
		for i := 0; i < len(s); i++ {
			if s[i] == c {
				return true
			}
		}
		return false
	}
	// Sweep the full 16-bit space for proquint-1.
	for n := 0; n < 65536; n++ {
		q := uint16ToProquint(uint16(n))
		if len(q) != 5 {
			t.Fatalf("n=%d: length=%d", n, len(q))
		}
		positions := []struct {
			pos      int
			alphabet string
			kind     string
		}{
			{0, cons, "cons"},
			{1, vows, "vows"},
			{2, cons, "cons"},
			{3, vows, "vows"},
			{4, cons, "cons"},
		}
		for _, p := range positions {
			if !contains(p.alphabet, q[p.pos]) {
				t.Fatalf("n=%d q=%q: position %d byte %q not in %s alphabet",
					n, q, p.pos, q[p.pos], p.kind)
			}
		}
	}
}

// TestProquintFullSpaceUniqueness verifies that the 65,536 inputs in the
// uint16 range produce 65,536 distinct proquint strings (i.e. the encoding
// is a bijection on its domain). Without this property the collision check
// in mint() would be defending against false collisions.
func TestProquintFullSpaceUniqueness(t *testing.T) {
	seen := make(map[string]uint16, 65536)
	for n := 0; n < 65536; n++ {
		q := uint16ToProquint(uint16(n))
		if prev, dup := seen[q]; dup {
			t.Fatalf("collision: %q from n=%d and n=%d", q, prev, n)
		}
		seen[q] = uint16(n)
	}
	if len(seen) != 65536 {
		t.Errorf("unique proquints = %d, want 65536", len(seen))
	}
}
