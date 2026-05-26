package kernel

import (
	_ "embed"

	"promisegrid.dev/wire-lab/implementations/poc4/lib"
)

//go:embed specs/kernel-receive-draft.md
var receiveSpecBytes []byte

// ReceiveProtocolCID returns the pCID for local app receive promises.
func ReceiveProtocolCID() lib.ProtocolCID {
	return lib.NewProtocolCID(receiveSpecBytes)
}
