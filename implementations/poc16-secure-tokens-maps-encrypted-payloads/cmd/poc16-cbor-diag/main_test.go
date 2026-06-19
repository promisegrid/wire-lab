package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
)

func TestDiagnosticDecoderRendersNestedEnvelopeCBOR(t *testing.T) {
	protocolCID := protocol.NewProtocolCID([]byte("poc16 diagnostic test protocol"))
	envelope, envelopeErr := protocol.NewEnvelope(protocolCID, map[string]string{
		"from":          "alice",
		"to":            "bob",
		"promise_about": "diagnostic_review",
	}, "alice")
	if envelopeErr != nil {
		t.Fatalf("NewEnvelope: %v", envelopeErr)
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		t.Fatalf("Envelope.Bytes: %v", bytesErr)
	}
	value, parseErr := parseDiagnosticValue(envelopeBytes)
	if parseErr != nil {
		t.Fatalf("parseDiagnosticValue: %v", parseErr)
	}
	rendered := value.render(diagnosticOptions{expandBytes: true, maxBytes: 128})
	for _, expected := range []string{
		"1735551332(",
		"42(h'00",
		"<< {",
		"\"from\": \"alice\"",
		"\"promise_about\": \"diagnostic_review\"",
		"\"signer\": \"alice\"",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("rendered diagnostic missing %q:\n%s", expected, rendered)
		}
	}
}

func TestRunFindsMessageByHashInRunRoot(t *testing.T) {
	runRoot := t.TempDir()
	messageCASDir := filepath.Join(runRoot, "message-cas")
	if mkdirErr := os.MkdirAll(messageCASDir, 0o755); mkdirErr != nil {
		t.Fatalf("MkdirAll: %v", mkdirErr)
	}
	rawEnvelope := []byte{0x83, 0x01, 0x02, 0x03}
	exactHash := protocol.HashExactBytes(rawEnvelope)
	artifactPath := filepath.Join(messageCASDir, exactHash+".cbor")
	if writeErr := os.WriteFile(artifactPath, rawEnvelope, 0o644); writeErr != nil {
		t.Fatalf("WriteFile artifact: %v", writeErr)
	}
	record := messageDAGRecord{
		Source:       "alice/agent:alice/stdout",
		Observer:     "alice",
		Direction:    "sent",
		Peer:         "bob",
		Protocol:     "relationship_v1",
		ExactSHA256:  exactHash,
		PromiseAbout: "diagnostic_review",
		SourceEvent:  "runtime.sent",
		Path:         filepath.Join("message-cas", exactHash+".cbor"),
	}
	recordBytes, marshalErr := json.Marshal(record)
	if marshalErr != nil {
		t.Fatalf("Marshal: %v", marshalErr)
	}
	if writeErr := os.WriteFile(filepath.Join(runRoot, "message-dag.jsonl"), append(recordBytes, '\n'), 0o644); writeErr != nil {
		t.Fatalf("WriteFile index: %v", writeErr)
	}
	var output bytes.Buffer
	if runErr := run([]string{"-run-root", runRoot, "-hash", exactHash}, &output); runErr != nil {
		t.Fatalf("run: %v", runErr)
	}
	rendered := output.String()
	for _, expected := range []string{
		"# metadata: source=alice/agent:alice/stdout observer=alice direction=sent peer=bob protocol=relationship_v1",
		"# bytes: 4",
		"[\n  1,\n  2,\n  3,\n]",
	} {
		if !strings.Contains(rendered, expected) {
			t.Fatalf("output missing %q:\n%s", expected, rendered)
		}
	}
}
