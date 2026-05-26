package hello

import (
	_ "embed"

	"promisegrid.dev/wire-lab/implementations/poc3/lib"
)

//go:embed specs/hello-draft.md
var helloSpecBytes []byte

// ProtocolCID returns the pCID for the poc3 hello protocol.
func ProtocolCID() lib.ProtocolCID {
	return lib.NewProtocolCID(helloSpecBytes)
}
