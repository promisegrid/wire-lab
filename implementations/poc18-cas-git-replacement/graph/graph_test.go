package graph

import (
	"strings"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

func TestMessageUsesGridTagAndTag42PCID(t *testing.T) {
	cas, openErr := store.Open(t.TempDir())
	if openErr != nil {
		t.Fatalf("Open() error = %v", openErr)
	}
	message, storeErr := StoreMessage(cas, nil, Payload{
		Promiser:    "alice",
		Promisee:    "bob",
		PromiseKind: "reference_set",
		PromiseBody: ReferenceSetBody("refset:test", "branch", "project:test", []any{}, []any{}),
	})
	if storeErr != nil {
		t.Fatalf("StoreMessage() error = %v", storeErr)
	}
	diagnostic, diagErr := Diagnostic(message.Bytes)
	if diagErr != nil {
		t.Fatalf("Diagnostic() error = %v", diagErr)
	}
	if !strings.Contains(diagnostic, "grid(") {
		t.Fatalf("diagnostic missing grid tag: %s", diagnostic)
	}
	if !strings.Contains(diagnostic, "42("+VersionControlPCIDText+")") {
		t.Fatalf("diagnostic missing pCID: %s", diagnostic)
	}
}
