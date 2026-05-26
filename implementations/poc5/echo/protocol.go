package echo

import (
	_ "embed"

	"promisegrid.dev/wire-lab/implementations/poc5/lib"
)

//go:embed specs/echo-draft.md
var echoSpec []byte

// ProtocolCID returns the content-derived pCID for echo messages.
func ProtocolCID() lib.ProtocolCID {
	return lib.NewProtocolCID(echoSpec)
}
