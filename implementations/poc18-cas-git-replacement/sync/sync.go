// Package sync models POC18 sparse retrieval promises without implementing real
// peer transport in the first slice.
//
// Intent: Preserve the peer-relative continuous-DAG-sync seam while the initial
// code stays deterministic and local. Source: DI-dibut; DI-jifuj
package sync

import (
	cidlib "github.com/ipfs/go-cid"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

// MissingObject records a CID that a local graph walk wants but the sparse CAS
// does not currently retain.
type MissingObject struct {
	CID  string `json:"cid"`
	Role string `json:"role"`
}

// Interest is the local shape that will become a sync_interest promise in a
// later transport slice.
type Interest struct {
	Promiser string          `json:"promiser"`
	Promisee string          `json:"promisee"`
	Wanted   []MissingObject `json:"wanted"`
	Offer    string          `json:"offer"`
}

// PlanInterest returns the missing-object interest for wanted CIDs.
func PlanInterest(cas *store.FileStore, promiser, promisee string, wanted map[string]cidlib.Cid) Interest {
	missing := []MissingObject{}
	for role, objectCID := range wanted {
		if !cas.Has(objectCID) {
			missing = append(missing, MissingObject{CID: store.CIDText(objectCID), Role: role})
		}
	}
	return Interest{Promiser: promiser, Promisee: promisee, Wanted: missing, Offer: "promise to receive and reciprocate later"}
}
