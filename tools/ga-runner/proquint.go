package main

const (
	proquintCons = "bdfghjklmnprstvz"
	proquintVows = "aiou"
)

// uint16ToProquint turns random child-ID entropy into a pronounceable handle
// that matches the repo's existing TODO/TE/DR/DI handle style.
//
// Intent: Generated child simulation IDs should be human-reviewable paths, not
// opaque UUIDs, while remaining collision-resistant enough for one run.
// Source: DI-gijom
func uint16ToProquint(value uint16) string {
	buf := []byte{
		proquintCons[(value>>12)&0x0f],
		proquintVows[(value>>10)&0x03],
		proquintCons[(value>>6)&0x0f],
		proquintVows[(value>>4)&0x03],
		proquintCons[value&0x0f],
	}
	return string(buf)
}
