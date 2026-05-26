package echo

import (
	_ "embed"

	"promisegrid.dev/wire-lab/implementations/poc3/lib"
)

//go:embed specs/echo-draft.md
var echoSpecBytes []byte

// ProtocolCID returns the pCID for the poc3 echo protocol.
func ProtocolCID() lib.ProtocolCID {
	return lib.NewProtocolCID(echoSpecBytes)
}
