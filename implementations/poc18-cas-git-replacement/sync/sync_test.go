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

// TestRunContinuousDAGSyncPropagatesThenIdempotentlyRepeats proves that the
// nahop.19 continuous-sync slice moves a useful update once and then becomes a
// no-op when the receiver already retains the advertised DAG.
func TestRunContinuousDAGSyncPropagatesThenIdempotentlyRepeats(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "alice-workspace")
	if err := os.MkdirAll(source, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(source, "README.md"), []byte("continuous sync from alice\n"), 0o644); err != nil {
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
	report, syncErr := pocsync.RunContinuousDAGSync(
		pocsync.Peer{Agent: "alice", CAS: aliceCAS},
		pocsync.Peer{Agent: "bob", CAS: bobCAS},
		map[string]cidlib.Cid{"branch": branchCID},
		nil,
		pocsync.ContinuousSyncConfig{Rounds: 2, Offer: "continuous_storage_credit:1", RetainUntil: "2026-07-31T00:00:00Z"},
	)
	if syncErr != nil {
		t.Fatalf("RunContinuousDAGSync() error = %v", syncErr)
	}
	if report.Rounds != 2 || len(report.Directions) != 2 {
		t.Fatalf("continuous report shape = %#v, want two alice->bob rounds", report)
	}
	if !report.Directions[0].UsefulUpdate || len(report.Directions[0].Retrieval.Retrieved) == 0 {
		t.Fatalf("first continuous round did not retrieve useful objects: %#v", report.Directions[0])
	}
	if report.Directions[0].RetentionMessageCID == "" {
		t.Fatalf("first continuous round did not store retention promise: %#v", report.Directions[0])
	}
	if report.Directions[1].UsefulUpdate || len(report.Directions[1].Retrieval.Retrieved) != 0 || len(report.Directions[1].Retrieval.AlreadyLocal) == 0 {
		t.Fatalf("second continuous round was not idempotent: %#v", report.Directions[1])
	}
	if report.MissingObjects != 0 {
		t.Fatalf("continuous sync missing objects = %d, want none", report.MissingObjects)
	}
	retentionCID, retentionErr := store.ParseCIDText(report.Directions[0].RetentionMessageCID)
	if retentionErr != nil {
		t.Fatalf("ParseCIDText(retention) error = %v", retentionErr)
	}
	if !aliceCAS.Has(retentionCID) {
		t.Fatalf("alice CAS does not contain bob retention promise %s", report.Directions[0].RetentionMessageCID)
	}
	if !bobCAS.Has(retentionCID) {
		t.Fatalf("bob CAS does not contain own retention promise %s", report.Directions[0].RetentionMessageCID)
	}
	snapshotCID, snapshotErr := store.ParseCIDText(result.SnapshotCID)
	if snapshotErr != nil {
		t.Fatalf("ParseCIDText(snapshot) error = %v", snapshotErr)
	}
	if !bobCAS.Has(snapshotCID) {
		t.Fatalf("bob CAS does not contain continuously synced snapshot %s", result.SnapshotCID)
	}
}
