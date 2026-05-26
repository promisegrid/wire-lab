package signed

import (
	_ "embed"

	"promisegrid.dev/wire-lab/implementations/poc3/lib"
)

//go:embed specs/signed-draft.md
var signedSpecBytes []byte

// ProtocolCID returns the pCID for the poc3 signed protocol.
func ProtocolCID() lib.ProtocolCID {
	return lib.NewProtocolCID(signedSpecBytes)
}
