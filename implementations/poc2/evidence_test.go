package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestObservationEnvelopeUsesLocalEvidenceHash(t *testing.T) {
	protocolCID, cidErr := HelloProtocolCID()
	if cidErr != nil {
		t.Fatalf("HelloProtocolCID: %v", cidErr)
	}
	observation, observationErr := observationEnvelope(protocolCID, "bob-kernel", "alice-kernel", "kept", "delivered", []byte("exact message bytes"))
	if observationErr != nil {
		t.Fatalf("observationEnvelope: %v", observationErr)
	}
	kind, fields, kindErr := envelopeKind(observation)
	if kindErr != nil {
		t.Fatalf("envelopeKind: %v", kindErr)
	}
	if kind != "observation_v1" {
		t.Fatalf("unexpected kind %q", kind)
	}
	if fields["outcome"] != "kept" {
		t.Fatalf("unexpected outcome %q", fields["outcome"])
	}
	if fields["evidence_hash"] == "" {
		t.Fatalf("missing evidence hash")
	}
}

func TestEvidenceRecordIsJSON(t *testing.T) {
	record := EvidenceRecord{Node: "alice", Event: "app_receive", Outcome: "kept", PromiseSource: "local observation; not global authority"}
	encodedRecord, marshalErr := json.Marshal(record)
	if marshalErr != nil {
		t.Fatalf("json marshal: %v", marshalErr)
	}
	if !strings.Contains(string(encodedRecord), "not global authority") {
		t.Fatalf("encoded record lost Promise Theory source: %s", string(encodedRecord))
	}
}
