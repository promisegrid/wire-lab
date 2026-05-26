package hello

import (
	_ "embed"

	"promisegrid.dev/wire-lab/implementations/poc5/lib"
)

//go:embed specs/hello-draft.md
var helloSpec []byte

// ProtocolCID returns the content-derived pCID for hello messages.
func ProtocolCID() lib.ProtocolCID {
	return lib.NewProtocolCID(helloSpec)
}
