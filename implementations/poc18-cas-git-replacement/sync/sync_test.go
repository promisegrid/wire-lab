package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	cidlib "github.com/ipfs/go-cid"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	pocsync "promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/sync"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/workspace"
)

// TestRetrieveGraphCopiesCIDVerifiedObjectsBetweenSparseStores proves the
// nahop.11 retrieval slice with two separate sparse CAS roots.
func TestRetrieveGraphCopiesCIDVerifiedObjectsBetweenSparseStores(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "alice-workspace")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "README.md"), []byte("hello from alice\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	aliceCAS, aliceErr := store.Open(filepath.Join(root, "alice-cas"))
	if aliceErr != nil {
		t.Fatalf("Open(alice) error = %v", aliceErr)
	}
	bobCAS, bobErr := store.Open(filepath.Join(root, "bob-cas"))
	if bobErr != nil {
		t.Fatalf("Open(bob) error = %v", bobErr)
	}
	result, ingestErr := workspace.NewScanner(aliceCAS, "alice", "bob").Ingest(source)
	if ingestErr != nil {
		t.Fatalf("Ingest() error = %v", ingestErr)
	}
	branchCID, branchErr := store.ParseCIDText(result.BranchRefSetCID)
	if branchErr != nil {
		t.Fatalf("ParseCIDText(branch) error = %v", branchErr)
	}
	report, retrieveErr := pocsync.RetrieveGraph(
		pocsync.Peer{Agent: "bob", CAS: bobCAS},
		pocsync.Peer{Agent: "alice", CAS: aliceCAS},
		map[string]cidlib.Cid{"branch": branchCID},
		"storage_credit:2",
	)
	if retrieveErr != nil {
		t.Fatalf("RetrieveGraph() error = %v", retrieveErr)
	}
	if len(report.Retrieved) == 0 {
		t.Fatalf("RetrieveGraph() retrieved no objects: %#v", report)
	}
	if len(report.Missing) != 0 {
		t.Fatalf("RetrieveGraph() missing = %#v, want none", report.Missing)
	}
	if report.InterestMessageCID == "" || report.AvailabilityMessageCID == "" {
		t.Fatalf("missing sync promise message CIDs: %#v", report)
	}
	interestCID, interestErr := store.ParseCIDText(report.InterestMessageCID)
	if interestErr != nil {
		t.Fatalf("ParseCIDText(interest) error = %v", interestErr)
	}
	if !bobCAS.Has(interestCID) {
		t.Fatalf("bob CAS does not contain sync_interest message %s", report.InterestMessageCID)
	}
	if !aliceCAS.Has(interestCID) {
		t.Fatalf("alice CAS does not contain received sync_interest message %s", report.InterestMessageCID)
	}
	availabilityCID, availabilityErr := store.ParseCIDText(report.AvailabilityMessageCID)
	if availabilityErr != nil {
		t.Fatalf("ParseCIDText(availability) error = %v", availabilityErr)
	}
	if !bobCAS.Has(availabilityCID) {
		t.Fatalf("bob CAS does not contain object_availability message %s", report.AvailabilityMessageCID)
	}
	if !aliceCAS.Has(availabilityCID) {
		t.Fatalf("alice CAS does not contain object_availability message %s", report.AvailabilityMessageCID)
	}
	snapshotCID, snapshotErr := store.ParseCIDText(result.SnapshotCID)
	if snapshotErr != nil {
		t.Fatalf("ParseCIDText(snapshot) error = %v", snapshotErr)
	}
	if !bobCAS.Has(snapshotCID) {
		t.Fatalf("bob CAS does not contain fetched snapshot %s", result.SnapshotCID)
	}
	checkout := filepath.Join(root, "bob-checkout")
	if err := workspace.MaterializeSnapshot(bobCAS, snapshotCID, checkout); err != nil {
		t.Fatalf("MaterializeSnapshot() error = %v", err)
	}
	content, readErr := os.ReadFile(filepath.Join(checkout, "README.md"))
	if readErr != nil {
		t.Fatalf("ReadFile(checkout) error = %v", readErr)
	}
	if string(content) != "hello from alice\n" {
		t.Fatalf("checkout README = %q", string(content))
	}
}
