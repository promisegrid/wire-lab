package signed

import (
	_ "embed"

	"promisegrid.dev/wire-lab/implementations/poc5/lib"
)

//go:embed specs/signed-draft.md
var signedSpec []byte

// ProtocolCID returns the content-derived pCID for signed-app messages.
func ProtocolCID() lib.ProtocolCID {
	return lib.NewProtocolCID(signedSpec)
}
