package relay

import (
	_ "embed"

	"promisegrid.dev/wire-lab/implementations/poc5/lib"
)

//go:embed specs/relay-draft.md
var relaySpecBytes []byte

// ProtocolCID returns the pCID for poc5 relay wrappers.
func ProtocolCID() lib.ProtocolCID {
	return lib.NewProtocolCID(relaySpecBytes)
}
