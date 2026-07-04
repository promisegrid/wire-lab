// Package retention models local POC18 object-retention promises and pressure GC.
//
// Intent: Retention is local promise accounting over a sparse CAS. A collector
// may withdraw local storage for unpromised bytes under pressure, but it does
// not command peers or assert global deletion. Source: DI-mivur
package retention

import (
	"fmt"
	"sort"

	"github.com/fxamacker/cbor/v2"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

// Target names one CID that a local agent promises to retain.
type Target struct {
	Role string `json:"role"`
	CID  string `json:"cid"`
}

// Promise is the deterministic fixture shape for one object_retention promise.
type Promise struct {
	Promiser    string   `json:"promiser"`
	Promisee    string   `json:"promisee"`
	Scope       string   `json:"scope"`
	Targets     []Target `json:"targets"`
	RetainUntil string   `json:"retain_until"`
	// Parents lets the retention promise preserve payment-redemption context in
	// the same message DAG as the retained objects. Intent: paid retention should
	// be auditable through CIDs rather than hidden fixture state. Source: DI-bidum
	Parents              []graph.Parent
	CollectionTerms      []string `json:"collection_terms"`
	ReciprocalEvidence   []any    `json:"reciprocal_evidence"`
	LocalConstraintTerms []string `json:"local_constraint_terms"`
}

// Report records one local retention and GC pass.
type Report struct {
	Promiser              string        `json:"promiser"`
	Promisee              string        `json:"promisee"`
	RetentionMessageCID   string        `json:"retention_message_cid"`
	Targets               []Target      `json:"targets"`
	ProtectedCIDs         []string      `json:"protected_cids"`
	MissingProtectedCIDs  []Target      `json:"missing_protected_cids"`
	BeforeObjects         int           `json:"before_objects"`
	AfterObjects          int           `json:"after_objects"`
	ProtectedObjects      int           `json:"protected_objects"`
	CollectedObjects      int           `json:"collected_objects"`
	PressureTargetObjects int           `json:"pressure_target_objects"`
	TemporaryObjectCID    string        `json:"temporary_object_cid,omitempty"`
	Collected             []store.Entry `json:"collected"`
	Kept                  []store.Entry `json:"kept"`
}

// PromiseAndCollect stores an object_retention promise, computes the local
// promised closure, and collects unpromised objects until pressure is relieved.
//
// Intent: The retention promise itself becomes a protected object along with the
// promised roots and their locally available dependencies. The pressure target
// is local storage policy, not a peer-visible deletion request. Source: DI-mivur
func PromiseAndCollect(cas *store.FileStore, promise Promise, pressureTargetObjects int) (Report, error) {
	if cas == nil {
		return Report{}, fmt.Errorf("CAS store is required")
	}
	message, messageErr := StorePromise(cas, promise)
	if messageErr != nil {
		return Report{}, messageErr
	}
	roots := append([]Target(nil), promise.Targets...)
	roots = append(roots, Target{Role: "retention_message", CID: store.CIDText(message.CID)})
	protected, missing, closureErr := promisedClosure(cas, roots)
	if closureErr != nil {
		return Report{}, closureErr
	}
	before, listErr := cas.List()
	if listErr != nil {
		return Report{}, listErr
	}
	// Intent: A non-positive target means "relieve pressure down to the promised
	// keep set" so diagnostics report a concrete object count rather than a magic
	// sentinel. Source: DI-mivur
	effectivePressureTarget := pressureTargetObjects
	if effectivePressureTarget <= 0 {
		effectivePressureTarget = countProtectedEntries(before, protected)
	}
	collected, collectErr := collectUnderPressure(cas, before, protected, effectivePressureTarget)
	if collectErr != nil {
		return Report{}, collectErr
	}
	after, afterErr := cas.List()
	if afterErr != nil {
		return Report{}, afterErr
	}
	kept := promisedEntries(after, protected)
	return Report{
		Promiser:              promise.Promiser,
		Promisee:              promise.Promisee,
		RetentionMessageCID:   store.CIDText(message.CID),
		Targets:               append([]Target(nil), promise.Targets...),
		ProtectedCIDs:         sortedCIDSet(protected),
		MissingProtectedCIDs:  missing,
		BeforeObjects:         len(before),
		AfterObjects:          len(after),
		ProtectedObjects:      len(kept),
		CollectedObjects:      len(collected),
		PressureTargetObjects: effectivePressureTarget,
		Collected:             collected,
		Kept:                  kept,
	}, nil
}

// StorePromise stores a signed object_retention promise message.
func StorePromise(cas *store.FileStore, promise Promise) (graph.StoredMessage, error) {
	rows, rowsErr := targetRows(promise.Targets)
	if rowsErr != nil {
		return graph.StoredMessage{}, rowsErr
	}
	payload := graph.Payload{
		Promiser:    promise.Promiser,
		Promisee:    promise.Promisee,
		PromiseKind: "object_retention",
		PromiseBody: graph.ObjectRetentionBody(
			promise.Scope,
			rows,
			promise.RetainUntil,
			stringRows(promise.CollectionTerms),
			anyRows(promise.ReciprocalEvidence),
		),
		ReciprocalPromise: anyRows(promise.ReciprocalEvidence),
		LocalConstraints:  stringRows(promise.LocalConstraintTerms),
	}
	return graph.StoreMessage(cas, promise.Parents, payload)
}

func targetRows(targets []Target) ([]any, error) {
	rows := make([]any, 0, len(targets))
	for _, target := range targets {
		targetCID, parseErr := store.ParseCIDText(target.CID)
		if parseErr != nil {
			return nil, parseErr
		}
		rows = append(rows, graph.ObjectRow(target.Role, targetCID, "retain"))
	}
	return rows, nil
}

func stringRows(values []string) []any {
	rows := make([]any, 0, len(values))
	for _, value := range values {
		rows = append(rows, value)
	}
	return rows
}

func anyRows(values []any) []any {
	if values == nil {
		return []any{}
	}
	return append([]any(nil), values...)
}

func promisedClosure(cas *store.FileStore, roots []Target) (map[string]bool, []Target, error) {
	protected := map[string]bool{}
	missing := []Target{}
	queue := append([]Target(nil), roots...)
	for len(queue) > 0 {
		next := queue[0]
		queue = queue[1:]
		if protected[next.CID] {
			continue
		}
		protected[next.CID] = true
		objectCID, parseErr := store.ParseCIDText(next.CID)
		if parseErr != nil {
			return nil, nil, parseErr
		}
		content, _, getErr := cas.Get(objectCID)
		if getErr != nil {
			missing = append(missing, next)
			continue
		}
		references, referenceErr := referencedCIDs(content)
		if referenceErr != nil {
			continue
		}
		for _, reference := range references {
			if !protected[reference.CID] {
				queue = append(queue, reference)
			}
		}
	}
	return protected, missing, nil
}

func referencedCIDs(content []byte) ([]Target, error) {
	if view, envelopeErr := graph.ParseEnvelope(content); envelopeErr == nil {
		return referencedEnvelopeCIDs(view), nil
	}
	var decoded any
	if err := store.UnmarshalCBOR(content, &decoded); err != nil {
		return nil, err
	}
	references := []Target{}
	collectReferences(decoded, &references)
	return references, nil
}

func referencedEnvelopeCIDs(view graph.EnvelopeView) []Target {
	references := make([]Target, 0, len(view.Parents))
	for _, parent := range view.Parents {
		references = append(references, Target{Role: "parent:" + parent.Role, CID: store.CIDText(parent.CID)})
	}
	if len(view.Payload) >= 4 {
		// Intent: Retention closure follows protocol payload links and envelope
		// parents, not pCID or proof links. The pCID and signer key are validation
		// context, not object bytes promised as retained content. Source: DI-mivur
		collectReferences(view.Payload[3], &references)
	}
	return references
}

func collectReferences(value any, references *[]Target) {
	switch typed := value.(type) {
	case cbor.Tag:
		if typed.Number == store.LinkTagNumber {
			if linkCID, cidErr := store.CIDFromLinkTag(typed); cidErr == nil {
				*references = append(*references, Target{Role: "linked_object", CID: store.CIDText(linkCID)})
			}
		}
		collectReferences(typed.Content, references)
	case []any:
		for _, child := range typed {
			collectReferences(child, references)
		}
	case map[any]any:
		for key, child := range typed {
			collectReferences(key, references)
			collectReferences(child, references)
		}
	}
}

func collectUnderPressure(cas *store.FileStore, entries []store.Entry, protected map[string]bool, pressureTargetObjects int) ([]store.Entry, error) {
	currentCount := len(entries)
	collected := []store.Entry{}
	for _, entry := range entries {
		if protected[entry.CID] {
			continue
		}
		if pressureTargetObjects > 0 && currentCount <= pressureTargetObjects {
			break
		}
		objectCID, parseErr := store.ParseCIDText(entry.CID)
		if parseErr != nil {
			return nil, parseErr
		}
		deleted, deleteErr := cas.Delete(objectCID)
		if deleteErr != nil {
			return nil, deleteErr
		}
		collected = append(collected, deleted)
		currentCount--
	}
	return collected, nil
}

func promisedEntries(entries []store.Entry, protected map[string]bool) []store.Entry {
	kept := []store.Entry{}
	for _, entry := range entries {
		if protected[entry.CID] {
			kept = append(kept, entry)
		}
	}
	return kept
}

func countProtectedEntries(entries []store.Entry, protected map[string]bool) int {
	count := 0
	for _, entry := range entries {
		if protected[entry.CID] {
			count++
		}
	}
	return count
}

func sortedCIDSet(values map[string]bool) []string {
	sorted := make([]string, 0, len(values))
	for value := range values {
		sorted = append(sorted, value)
	}
	sort.Strings(sorted)
	return sorted
}
