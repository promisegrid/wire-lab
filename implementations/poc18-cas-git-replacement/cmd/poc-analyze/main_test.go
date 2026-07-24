package main

import (
	"os"
	"path/filepath"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
)

func TestAnalyzeCollectorFailsMissingArtifacts(t *testing.T) {
	// Intent: Empty collector artifacts should be an analyzer failure, not a clean
	// run with zero evidence. Source: DI-bovaf
	runRoot := t.TempDir()
	writeTestFile(t, filepath.Join(runRoot, "events.jsonl"), "")
	writeTestFile(t, filepath.Join(runRoot, "message-dag.jsonl"), "")
	writeTestFile(t, filepath.Join(runRoot, "car-dag.jsonl"), "")

	report, analyzeErr := analyzeRunRoot(runRoot)
	if analyzeErr != nil {
		t.Fatalf("analyze empty collector artifacts: %v", analyzeErr)
	}
	if report.Pass {
		t.Fatalf("empty collector artifacts unexpectedly passed: %#v", report.Checks)
	}
	if got := report.Checks["collector_message_artifacts"]; got != "missing" {
		t.Fatalf("collector_message_artifacts=%q, want missing", got)
	}
}

func TestVerifyMessageArtifactsFailsCIDMismatch(t *testing.T) {
	// Intent: The analyzer must recompute message CIDs from exact CBOR bytes so a
	// stale or forged collector index cannot satisfy `nahop.20`. Source: DI-bovaf
	runRoot := t.TempDir()
	message, messageErr := graph.NewMessage(nil, graph.Payload{
		Promiser:    "alice",
		Promisee:    "bob",
		PromiseKind: "sync_interest",
		PromiseBody: []any{"test-exchange", "latest", []any{}, []any{"promise to receive"}, []any{}},
	})
	if messageErr != nil {
		t.Fatalf("new message: %v", messageErr)
	}
	messageBytes, bytesErr := message.Bytes()
	if bytesErr != nil {
		t.Fatalf("message bytes: %v", bytesErr)
	}
	messageDir := filepath.Join(runRoot, "message-cas")
	if mkdirErr := os.MkdirAll(messageDir, 0o755); mkdirErr != nil {
		t.Fatalf("mkdir message dir: %v", mkdirErr)
	}
	messagePath := filepath.Join("message-cas", "mismatched.cbor")
	writeTestFile(t, filepath.Join(runRoot, messagePath), string(messageBytes))

	checks := map[string]string{}
	pass, verifyErr := verifyMessageArtifacts(runRoot, []collectorMessageRecord{{
		ExactCID: "bafkreiaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Path:     messagePath,
	}}, checks)
	if verifyErr != nil {
		t.Fatalf("verify message artifacts: %v", verifyErr)
	}
	if pass {
		t.Fatalf("CID mismatch unexpectedly passed")
	}
	if got := checks["message_cid_match_bafkreiaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"]; got != "missing" {
		t.Fatalf("CID mismatch check=%q, want missing", got)
	}
}

func writeTestFile(t *testing.T, path string, content string) {
	t.Helper()
	if mkdirErr := os.MkdirAll(filepath.Dir(path), 0o755); mkdirErr != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), mkdirErr)
	}
	if writeErr := os.WriteFile(path, []byte(content), 0o644); writeErr != nil {
		t.Fatalf("write %s: %v", path, writeErr)
	}
}
