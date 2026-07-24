package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cidlib "github.com/ipfs/go-cid"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/scenario"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/workspace"
)

func TestDiagnosticReportWritesRepresentativeFlows(t *testing.T) {
	runRoot := t.TempDir()
	aliceCAS, openErr := store.Open(filepath.Join(runRoot, "alice-cas"))
	if openErr != nil {
		t.Fatalf("open alice CAS: %v", openErr)
	}
	referenceSet := mustStoreMessage(t, aliceCAS, "reference_set", graph.ReferenceSetBody("refset:test", "branch", "project:test", nil, nil))
	directoryNode := mustStoreMessage(t, aliceCAS, "posix_node", graph.PosixNodeBody("node:directory:.", "directory", store.LinkTag(referenceSet.CID), nil, nil))
	nodeVersion := mustStoreMessage(t, aliceCAS, "posix_node", graph.PosixNodeBody("node:README.md", "regular", []byte("hello"), nil, nil))
	snapshot := mustStoreMessage(t, aliceCAS, "snapshot", graph.SnapshotBody("snapshot:test", directoryNode.CID, nil, "test snapshot", nil))
	mergeSnapshot := mustStoreMessage(t, aliceCAS, "snapshot", graph.SnapshotBody("snapshot:merge", directoryNode.CID, []cidlib.Cid{mustParseCID(t, snapshot.Entry.CID)}, "merge snapshot", nil))
	review := mustStoreMessage(t, aliceCAS, "review_statement", graph.ReviewStatementBody("local_test", []any{graph.ObjectRow("snapshot", mergeSnapshot.CID)}, "test kept", "test_kept", nil))
	if _, storeErr := graph.StoreMessage(aliceCAS, nil, graph.Payload{
		Promiser:    "alice",
		Promisee:    "bob",
		PromiseKind: "chunk_manifest",
		PromiseBody: graph.ChunkManifestBody(referenceSet.CID, 5, "rabin", nil, nil, ""),
	}); storeErr != nil {
		t.Fatalf("store chunk manifest: %v", storeErr)
	}
	result := fixtureDiagnosticResult{
		IngestResult:   workspace.IngestResult{StoreRoot: aliceCAS.Root},
		AliceStoreRoot: aliceCAS.Root,
		Scenario: scenario.Result{
			MergeBranchRefSetCID:  referenceSet.Entry.CID,
			RenameLabelNodeCID:    nodeVersion.Entry.CID,
			RenameCopySnapshotCID: snapshot.Entry.CID,
			TestStatementCID:      review.Entry.CID,
			MergeSnapshotCID:      mergeSnapshot.Entry.CID,
		},
	}
	writeJSONForTest(t, filepath.Join(runRoot, "result.json"), result)
	messageRun := filepath.Join(runRoot, "observer", "poc18-demo")
	if mkdirErr := os.MkdirAll(filepath.Join(messageRun, "message-cas"), 0o755); mkdirErr != nil {
		t.Fatalf("mkdir message-cas: %v", mkdirErr)
	}
	for _, kind := range []string{"sync_interest", "object_availability", "object_retrieval_redemption"} {
		message := mustStoreMessage(t, aliceCAS, kind, []any{"test-exchange", []any{}, []any{}})
		messagePath := filepath.Join(messageRun, "message-cas", message.Entry.CID+".cbor")
		if writeErr := os.WriteFile(messagePath, message.Bytes, 0o644); writeErr != nil {
			t.Fatalf("write message artifact: %v", writeErr)
		}
		record := messageDAGRecord{Source: "agent:alice", Observer: "alice", Direction: "sent", Peer: "bob", Protocol: graph.VersionControlPCIDText, PromiseKind: kind, ExactCID: message.Entry.CID, Path: filepath.ToSlash(filepath.Join("message-cas", message.Entry.CID+".cbor"))}
		appendJSONLineForTest(t, filepath.Join(messageRun, "message-dag.jsonl"), record)
	}
	outDir := filepath.Join(runRoot, "poc18-diagnostics")
	if reportErr := writeDiagnosticReport(runRoot, outDir); reportErr != nil {
		t.Fatalf("write diagnostic report: %v", reportErr)
	}
	for _, flow := range []string{"reference-set", "node-version", "snapshot", "review-statement", "merge-snapshot", "directory-node", "materialization-object", "sync-interest", "object-availability", "object-retrieval-redemption"} {
		content, readErr := os.ReadFile(filepath.Join(outDir, flow+".diag.txt"))
		if readErr != nil {
			t.Fatalf("read %s diagnostic: %v", flow, readErr)
		}
		if len(content) == 0 {
			t.Fatalf("%s diagnostic is empty", flow)
		}
		if !strings.Contains(string(content), "header_hex:") {
			t.Fatalf("%s diagnostic missing header hex breakdown: %s", flow, string(content))
		}
	}
	if _, statErr := os.Stat(filepath.Join(outDir, "index.json")); statErr != nil {
		t.Fatalf("index.json missing: %v", statErr)
	}
}

func mustStoreMessage(t *testing.T, cas *store.FileStore, promiseKind string, body any) graph.StoredMessage {
	t.Helper()
	message, storeErr := graph.StoreMessage(cas, nil, graph.Payload{
		Promiser:    "alice",
		Promisee:    "bob",
		PromiseKind: promiseKind,
		PromiseBody: body,
	})
	if storeErr != nil {
		t.Fatalf("store %s message: %v", promiseKind, storeErr)
	}
	return message
}

func mustParseCID(t *testing.T, cidText string) cidlib.Cid {
	t.Helper()
	parsedCID, parseErr := store.ParseCIDText(cidText)
	if parseErr != nil {
		t.Fatalf("parse CID: %v", parseErr)
	}
	return parsedCID
}

func writeJSONForTest(t *testing.T, path string, value any) {
	t.Helper()
	content, marshalErr := json.MarshalIndent(value, "", "  ")
	if marshalErr != nil {
		t.Fatalf("marshal JSON: %v", marshalErr)
	}
	content = append(content, '\n')
	if writeErr := os.WriteFile(path, content, 0o644); writeErr != nil {
		t.Fatalf("write JSON: %v", writeErr)
	}
}

func appendJSONLineForTest(t *testing.T, path string, value any) {
	t.Helper()
	content, marshalErr := json.Marshal(value)
	if marshalErr != nil {
		t.Fatalf("marshal JSON line: %v", marshalErr)
	}
	file, openErr := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if openErr != nil {
		t.Fatalf("open JSONL: %v", openErr)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			t.Fatalf("close JSONL: %v", closeErr)
		}
	}()
	if _, writeErr := file.Write(append(content, '\n')); writeErr != nil {
		t.Fatalf("write JSONL: %v", writeErr)
	}
}
