package storage

import (
	_ "embed"

	"promisegrid.dev/wire-lab/implementations/poc5/lib"
)

//go:embed specs/storage-draft.md
var storageSpec []byte

// ProtocolCID returns the content-derived pCID for storage messages.
func ProtocolCID() lib.ProtocolCID {
	return lib.NewProtocolCID(storageSpec)
}
