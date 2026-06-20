package main

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/config"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/decision"
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

func TestCompactMonitorEventsBoundsInputAndPreservesCriticalSignals(t *testing.T) {
	// Intent: The live monitor prompt must stay bounded even when parser-role
	// raw diagnostics greatly increase event volume, while critical malformed and
	// parser-flow signals remain visible to the monitor. Source: DI-jozos
	events := make([]decision.Event, 0, 600)
	longDetail := strings.Repeat("x", monitorDetailLimit+20)
	for index := 0; index < 600; index++ {
		events = append(events, decision.Event{
			Observer: "agent",
			Event:    "ordinary_kept_event",
			Outcome:  "kept",
			Detail:   longDetail,
		})
	}
	events[200] = decision.Event{
		Observer: "parser:alice",
		Event:    "parser_role_payload_parsed",
		Outcome:  "kept",
		Detail:   longDetail,
	}
	events[300] = decision.Event{
		Observer: "alice",
		Event:    "raw_message_artifact_emitted",
		Outcome:  "kept",
		Detail:   "direction=app_to_parser pcid=relationship_v1 exact_sha256=abc",
	}
	events[400] = decision.Event{
		Observer: "bob",
		Event:    "local_resource_exhausted",
		Outcome:  "non_commitment",
		Detail:   "capacity exhausted",
	}
	compacted := compactMonitorEvents(events)
	if len(compacted) > monitorEventLimit {
		t.Fatalf("compacted event count = %d, want <= %d", len(compacted), monitorEventLimit)
	}
	if compacted[0].Event != "monitor_input_compacted" {
		t.Fatalf("first compacted event = %q, want monitor_input_compacted", compacted[0].Event)
	}
	for _, requiredEvent := range []string{"parser_role_payload_parsed", "raw_message_artifact_emitted", "local_resource_exhausted"} {
		if !hasCompactedEvent(compacted, requiredEvent) {
			t.Fatalf("compacted monitor input missing %s", requiredEvent)
		}
	}
	for _, event := range compacted {
		if len(event.Detail) > monitorDetailLimit+len("...[truncated for monitor prompt]") {
			t.Fatalf("event detail was not bounded: %q", event.Detail)
		}
	}
}

func hasCompactedEvent(events []decision.Event, eventName string) bool {
	for _, event := range events {
		if event.Event == eventName {
			return true
		}
	}
	return false
}
