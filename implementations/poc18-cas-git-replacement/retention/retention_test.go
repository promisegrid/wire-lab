package retention

import (
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

func TestPromiseAndCollectKeepsPromisedObjectsAndCollectsUnpromised(t *testing.T) {
	cas, openErr := store.Open(t.TempDir())
	if openErr != nil {
		t.Fatalf("Open() error = %v", openErr)
	}
	promisedMessage, messageErr := graph.StoreMessage(cas, nil, graph.Payload{
		Promiser:    "alice",
		Promisee:    "bob",
		PromiseKind: "reference_set",
		PromiseBody: graph.ReferenceSetBody("refset:test", "release", "test", []any{}, []any{"retain for test"}),
	})
	if messageErr != nil {
		t.Fatalf("StoreMessage() error = %v", messageErr)
	}
	tempEntry, tempErr := cas.Put("chunk", []byte("unpromised pressure bytes"))
	if tempErr != nil {
		t.Fatalf("Put(temp) error = %v", tempErr)
	}
	report, retentionErr := PromiseAndCollect(cas, Promise{
		Promiser:             "frank",
		Promisee:             "alice",
		Scope:                "test-retention",
		Targets:              []Target{{Role: "release", CID: store.CIDText(promisedMessage.CID)}},
		RetainUntil:          "2026-07-31T00:00:00Z",
		CollectionTerms:      []string{"collect unpromised under pressure"},
		ReciprocalEvidence:   []any{"paid_storage_credit:1"},
		LocalConstraintTerms: []string{"local-only test"},
	}, 0)
	if retentionErr != nil {
		t.Fatalf("PromiseAndCollect() error = %v", retentionErr)
	}
	if report.RetentionMessageCID == "" || report.CollectedObjects == 0 {
		t.Fatalf("report = %#v, want retention message and collected object", report)
	}
	if !cas.Has(promisedMessage.CID) {
		t.Fatalf("promised object was collected")
	}
	retentionCID, parseErr := store.ParseCIDText(report.RetentionMessageCID)
	if parseErr != nil {
		t.Fatalf("ParseCIDText(retention) error = %v", parseErr)
	}
	if !cas.Has(retentionCID) {
		t.Fatalf("retention promise message was collected")
	}
	tempCID, tempParseErr := store.ParseCIDText(tempEntry.CID)
	if tempParseErr != nil {
		t.Fatalf("ParseCIDText(temp) error = %v", tempParseErr)
	}
	if cas.Has(tempCID) {
		t.Fatalf("unpromised pressure object survived GC")
	}
}
