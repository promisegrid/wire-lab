package fibonacci

import (
	_ "embed"

	"promisegrid.dev/wire-lab/implementations/poc5/lib"
)

//go:embed specs/fibonacci-draft.md
var fibonacciSpec []byte

// ProtocolCID returns the content-derived pCID for fibonacci messages.
func ProtocolCID() lib.ProtocolCID {
	return lib.NewProtocolCID(fibonacciSpec)
}
