package main

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/config"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/eventstream"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
)

func TestCollectorListenAddressUsesConfiguredPort(t *testing.T) {
	listenAddress, err := collectorListenAddress("event-collector:9200")
	if err != nil {
		t.Fatalf("collector listen address: %v", err)
	}
	if listenAddress != ":9200" {
		t.Fatalf("listen address = %q, want :9200", listenAddress)
	}
}

func TestRecordMessageArtifactWritesBinaryCASAndIndex(t *testing.T) {
	rawEnvelope := []byte{0xd9, 0x67, 0x72, 0x69, 0x64, 0x83, 0x01, 0x02, 0x03}
	exactHash := protocol.HashExactBytes(rawEnvelope)
	runCollector := &collector{
		cfg: config.Config{
			RunRoot: t.TempDir(),
			RunID:   "test-run",
		},
	}
	artifact := eventstream.MessageArtifact{
		Observer:            "alice",
		Direction:           "sent",
		Peer:                "bob",
		Protocol:            "route_v1",
		ExactSHA256:         exactHash,
		EnvelopeBytesBase64: base64.StdEncoding.EncodeToString(rawEnvelope),
	}
	if err := runCollector.recordMessageArtifact("alice/agent:alice/stdout", artifact); err != nil {
		t.Fatalf("record message artifact: %v", err)
	}
	artifactPath := filepath.Join(runCollector.cfg.RunRoot, runCollector.cfg.RunID, "message-cas", exactHash+".cbor")
	storedBytes, readErr := os.ReadFile(artifactPath)
	if readErr != nil {
		t.Fatalf("read artifact: %v", readErr)
	}
	if string(storedBytes) != string(rawEnvelope) {
		t.Fatalf("stored artifact bytes = %x, want %x", storedBytes, rawEnvelope)
	}
	indexBytes, readIndexErr := os.ReadFile(filepath.Join(runCollector.cfg.RunRoot, runCollector.cfg.RunID, "message-dag.jsonl"))
	if readIndexErr != nil {
		t.Fatalf("read index: %v", readIndexErr)
	}
	if !strings.Contains(string(indexBytes), `"path":"message-cas/`+exactHash+`.cbor"`) {
		t.Fatalf("index does not name artifact path: %s", string(indexBytes))
	}
}
