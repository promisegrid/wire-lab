// Package sync models POC18 sparse retrieval promises without implementing real
// peer transport in the first slice.
//
// Intent: Preserve the peer-relative continuous-DAG-sync seam while the initial
// code stays deterministic and local. Source: DI-dibut; DI-jifuj
package sync

import (
	"fmt"
	"sort"

	cidlib "github.com/ipfs/go-cid"

	pocchunk "promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/chunk"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
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

// Peer names one local agent and the sparse CAS that agent controls.
type Peer struct {
	Agent string
	CAS   *store.FileStore
}

// AvailabilityPromise records what a serving peer locally promises to provide.
type AvailabilityPromise struct {
	Promiser string          `json:"promiser"`
	Promisee string          `json:"promisee"`
	Objects  []MissingObject `json:"objects"`
	Terms    []string        `json:"terms"`
}

// RetrievedObject records one CID-verified byte transfer between sparse stores.
type RetrievedObject struct {
	CID  string `json:"cid"`
	Role string `json:"role"`
	Kind string `json:"kind"`
	From string `json:"from"`
	To   string `json:"to"`
	Size int64  `json:"size"`
}

// RetrievalReport is fixture evidence for one peer-relative retrieval attempt.
type RetrievalReport struct {
	Requester              string              `json:"requester"`
	Provider               string              `json:"provider"`
	Interest               Interest            `json:"interest"`
	InterestMessageCID     string              `json:"interest_message_cid"`
	Availability           AvailabilityPromise `json:"availability"`
	AvailabilityMessageCID string              `json:"availability_message_cid"`
	Retrieved              []RetrievedObject   `json:"retrieved"`
	AlreadyLocal           []MissingObject     `json:"already_local"`
	Missing                []MissingObject     `json:"missing"`
	Offer                  string              `json:"offer"`
}

// PlanInterest returns the missing-object interest for wanted CIDs.
func PlanInterest(cas *store.FileStore, promiser, promisee string, wanted map[string]cidlib.Cid) Interest {
	missing := []MissingObject{}
	for _, role := range sortedRoles(wanted) {
		objectCID := wanted[role]
		if !cas.Has(objectCID) {
			missing = append(missing, MissingObject{CID: store.CIDText(objectCID), Role: role})
		}
	}
	return Interest{Promiser: promiser, Promisee: promisee, Wanted: missing, Offer: "promise to receive and reciprocate later"}
}

// RetrieveGraph walks wanted roots, asks provider for locally missing objects,
// and retains only exact CID-verified bytes in requester.CAS.
//
// Intent: Prove sparse peer retrieval without introducing RPC semantics or a
// shared repository. The provider promises availability object-by-object, and
// the requester locally verifies every byte string against the requested CID
// before storing it. Source: DI-gozov
func RetrieveGraph(requester Peer, provider Peer, wanted map[string]cidlib.Cid, offer string) (RetrievalReport, error) {
	if requester.Agent == "" || provider.Agent == "" {
		return RetrievalReport{}, fmt.Errorf("requester and provider agent names are required")
	}
	if requester.CAS == nil || provider.CAS == nil {
		return RetrievalReport{}, fmt.Errorf("requester and provider CAS stores are required")
	}
	report := RetrievalReport{
		Requester: requester.Agent,
		Provider:  provider.Agent,
		Offer:     offer,
	}
	report.Interest = PlanInterest(requester.CAS, requester.Agent, provider.Agent, wanted)
	report.Interest.Offer = offer
	interestMessage, interestErr := storeInterestMessage(requester, provider, report.Interest)
	if interestErr != nil {
		return RetrievalReport{}, interestErr
	}
	report.InterestMessageCID = store.CIDText(interestMessage.CID)
	// Intent: Alice retains Bob's interest promise before answering so the
	// response parent link is a real local DAG edge, not a dangling test artifact.
	// Source: DI-gozov
	if _, putErr := provider.CAS.Put("message", interestMessage.Bytes); putErr != nil {
		return RetrievalReport{}, putErr
	}
	queue := initialQueue(wanted)
	seen := map[string]bool{}
	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]
		if seen[next.CID] {
			continue
		}
		seen[next.CID] = true
		objectCID, parseErr := store.ParseCIDText(next.CID)
		if parseErr != nil {
			return RetrievalReport{}, parseErr
		}
		content, entry, localErr := requester.CAS.Get(objectCID)
		if localErr == nil {
			report.AlreadyLocal = append(report.AlreadyLocal, next)
			references, referenceErr := referencesFromStoredObject(content, entry.Kind)
			if referenceErr != nil {
				return RetrievalReport{}, referenceErr
			}
			queue = append(queue, references...)
			continue
		}
		providerContent, providerEntry, providerErr := provider.CAS.Get(objectCID)
		if providerErr != nil {
			report.Missing = append(report.Missing, next)
			continue
		}
		report.Availability.Objects = append(report.Availability.Objects, next)
		storedEntry, putErr := requester.CAS.Put(providerEntry.Kind, providerContent)
		if putErr != nil {
			return RetrievalReport{}, putErr
		}
		if storedEntry.CID != next.CID {
			return RetrievalReport{}, fmt.Errorf("retrieved object CID mismatch: got %s want %s", storedEntry.CID, next.CID)
		}
		report.Retrieved = append(report.Retrieved, RetrievedObject{
			CID:  next.CID,
			Role: next.Role,
			Kind: providerEntry.Kind,
			From: provider.Agent,
			To:   requester.Agent,
			Size: providerEntry.Size,
		})
		references, referenceErr := referencesFromStoredObject(providerContent, providerEntry.Kind)
		if referenceErr != nil {
			return RetrievalReport{}, referenceErr
		}
		queue = append(queue, references...)
	}
	report.Availability.Promiser = provider.Agent
	report.Availability.Promisee = requester.Agent
	report.Availability.Terms = []string{offer, "serve only exact CID-verified bytes available in local sparse CAS"}
	availabilityMessage, availabilityErr := storeAvailabilityMessage(provider, requester, interestMessage.CID, report)
	if availabilityErr != nil {
		return RetrievalReport{}, availabilityErr
	}
	report.AvailabilityMessageCID = store.CIDText(availabilityMessage.CID)
	if _, putErr := requester.CAS.Put("message", availabilityMessage.Bytes); putErr != nil {
		return RetrievalReport{}, putErr
	}
	return report, nil
}

func sortedRoles(wanted map[string]cidlib.Cid) []string {
	roles := make([]string, 0, len(wanted))
	for role := range wanted {
		roles = append(roles, role)
	}
	sort.Strings(roles)
	return roles
}

func initialQueue(wanted map[string]cidlib.Cid) []MissingObject {
	queue := make([]MissingObject, 0, len(wanted))
	for _, role := range sortedRoles(wanted) {
		queue = append(queue, MissingObject{CID: store.CIDText(wanted[role]), Role: role})
	}
	return queue
}

func storeInterestMessage(requester Peer, provider Peer, interest Interest) (graph.StoredMessage, error) {
	wantedRows, rowsErr := interestRows(interest.Wanted)
	if rowsErr != nil {
		return graph.StoredMessage{}, rowsErr
	}
	payload := graph.Payload{
		Promiser:    requester.Agent,
		Promisee:    provider.Agent,
		PromiseKind: "sync_interest",
		PromiseBody: graph.SyncInterestBody(
			"requested_roots_and_dependencies",
			wantedRows,
			[]any{interest.Offer},
			[]any{"provider remains free not to promise service"},
		),
		ReciprocalPromise: []any{interest.Offer},
		LocalConstraints:  []any{"verify exact CID before retaining bytes"},
	}
	return graph.StoreMessage(requester.CAS, nil, payload)
}

func storeAvailabilityMessage(provider Peer, requester Peer, interestCID cidlib.Cid, report RetrievalReport) (graph.StoredMessage, error) {
	availableRows, rowsErr := availabilityRows(report)
	if rowsErr != nil {
		return graph.StoredMessage{}, rowsErr
	}
	payload := graph.Payload{
		Promiser:    provider.Agent,
		Promisee:    requester.Agent,
		PromiseKind: "object_availability",
		PromiseBody: graph.ObjectAvailabilityBody(
			"response:"+store.CIDText(interestCID),
			availableRows,
			[]any{"service offered only for locally retained exact bytes", report.Offer},
		),
		ReciprocalPromise: []any{report.Offer},
		LocalConstraints:  []any{"no promise for objects absent from local sparse CAS"},
	}
	return graph.StoreMessage(provider.CAS, []graph.Parent{{Role: "responds_to", CID: interestCID}}, payload)
}

func interestRows(objects []MissingObject) ([]any, error) {
	rows := make([]any, 0, len(objects))
	for _, object := range objects {
		objectCID, parseErr := store.ParseCIDText(object.CID)
		if parseErr != nil {
			return nil, parseErr
		}
		rows = append(rows, graph.ObjectRow(object.Role, objectCID))
	}
	return rows, nil
}

func availabilityRows(report RetrievalReport) ([]any, error) {
	rows := make([]any, 0, len(report.Retrieved)+len(report.Missing))
	for _, object := range report.Retrieved {
		objectCID, parseErr := store.ParseCIDText(object.CID)
		if parseErr != nil {
			return nil, parseErr
		}
		rows = append(rows, graph.ObjectRow(object.Role, objectCID, uint64(object.Size), "have"))
	}
	for _, object := range report.Missing {
		objectCID, parseErr := store.ParseCIDText(object.CID)
		if parseErr != nil {
			return nil, parseErr
		}
		rows = append(rows, graph.ObjectRow(object.Role, objectCID, uint64(0), "not_promised"))
	}
	return rows, nil
}

func referencesFromStoredObject(content []byte, kind string) ([]MissingObject, error) {
	if kind == "message" {
		view, envelopeErr := graph.ParseEnvelope(content)
		if envelopeErr != nil {
			return nil, envelopeErr
		}
		return referencesFromEnvelope(view)
	}
	if kind == "chunk_manifest" {
		manifest, manifestErr := pocchunk.DecodeManifest(content)
		if manifestErr != nil {
			return nil, manifestErr
		}
		return referencesFromManifest(manifest), nil
	}
	return nil, nil
}

func referencesFromManifest(manifest pocchunk.Manifest) []MissingObject {
	references := make([]MissingObject, 0, len(manifest.Chunks))
	for _, chunkRef := range manifest.Chunks {
		references = append(references, MissingObject{CID: chunkRef.CID, Role: "chunk"})
	}
	return references
}

func referencesFromEnvelope(view graph.EnvelopeView) ([]MissingObject, error) {
	references := make([]MissingObject, 0, len(view.Parents))
	for _, parent := range view.Parents {
		references = append(references, MissingObject{CID: store.CIDText(parent.CID), Role: "parent:" + parent.Role})
	}
	kind, kindErr := view.PayloadKind()
	if kindErr != nil {
		return nil, kindErr
	}
	body, bodyErr := view.PayloadBody()
	if bodyErr != nil {
		return nil, bodyErr
	}
	switch kind {
	case "reference_set":
		return append(references, referencesFromReferenceSet(body)...), nil
	case "snapshot":
		return append(references, referencesFromSnapshot(body)...), nil
	case "posix_node":
		return append(references, referencesFromPosixNode(body)...), nil
	case "chunk_manifest":
		return append(references, referencesFromChunkManifestMessage(body)...), nil
	default:
		return references, nil
	}
}

func referencesFromReferenceSet(body []any) []MissingObject {
	if len(body) != 5 {
		return nil
	}
	entries, ok := body[3].([]any)
	if !ok {
		return nil
	}
	references := []MissingObject{}
	for _, entryValue := range entries {
		entry, entryOK := entryValue.([]any)
		if !entryOK || len(entry) != 3 {
			continue
		}
		targets, targetsOK := entry[1].([]any)
		if !targetsOK {
			continue
		}
		for _, targetValue := range targets {
			target, targetOK := targetValue.([]any)
			if !targetOK || len(target) != 2 {
				continue
			}
			role, roleOK := target[0].(string)
			targetCID, cidErr := store.CIDFromLinkTag(target[1])
			if !roleOK || cidErr != nil {
				continue
			}
			references = append(references, MissingObject{CID: store.CIDText(targetCID), Role: role})
		}
	}
	return references
}

func referencesFromSnapshot(body []any) []MissingObject {
	if len(body) != 5 {
		return nil
	}
	references := []MissingObject{}
	if rootCID, cidErr := store.CIDFromLinkTag(body[1]); cidErr == nil {
		references = append(references, MissingObject{CID: store.CIDText(rootCID), Role: "root_directory"})
	}
	parentValues, ok := body[2].([]any)
	if !ok {
		return references
	}
	for _, parentValue := range parentValues {
		parentCID, cidErr := store.CIDFromLinkTag(parentValue)
		if cidErr == nil {
			references = append(references, MissingObject{CID: store.CIDText(parentCID), Role: "parent_snapshot"})
		}
	}
	return references
}

func referencesFromPosixNode(body []any) []MissingObject {
	if len(body) != 5 {
		return nil
	}
	nodeType, ok := body[1].(string)
	if !ok || (nodeType != "directory" && nodeType != "regular") {
		return nil
	}
	targetCID, cidErr := store.CIDFromLinkTag(body[2])
	if cidErr != nil {
		return nil
	}
	return []MissingObject{{CID: store.CIDText(targetCID), Role: nodeType + "_content"}}
}

func referencesFromChunkManifestMessage(body []any) []MissingObject {
	if len(body) != 6 {
		return nil
	}
	manifestCID, cidErr := store.CIDFromLinkTag(body[0])
	if cidErr != nil {
		return nil
	}
	return []MissingObject{{CID: store.CIDText(manifestCID), Role: "chunk_manifest_object"}}
}
