package protocol

import (
	"strings"
	"testing"
)

func TestMapBodyPayloadRoundTripKeepsCoreAndBodySeparate(t *testing.T) {
	// Intent: Flexible keyed details belong in slot-4's nested CBOR map, while
	// promiser/promisee/promise_about remain positional core slots. Source:
	// DI-mapah
	payloadBytes, err := MarshalRelationshipPayloadFields(map[string]string{
		"from":          "alice",
		"to":            "bob",
		"promise_about": "trust_update",
		"outcome":       "kept",
		"promise":       "Alice promises her local trust in Bob increased.",
		"reason":        "Bob kept the storage promise.",
		"turn":          "turn-0042",
		"subject":       "bob",
		"delta":         "+3",
	})
	if err != nil {
		t.Fatalf("marshal relationship payload: %v", err)
	}
	fields, err := RelationshipPayloadFields(payloadBytes)
	if err != nil {
		t.Fatalf("decode relationship payload: %v", err)
	}
	if fields["from"] != "alice" || fields["to"] != "bob" || fields["promise_about"] != "trust_update" {
		t.Fatalf("core fields were not preserved: %#v", fields)
	}
	if fields["subject"] != "bob" || fields["delta"] != "+3" {
		t.Fatalf("body fields were not projected for local compatibility: %#v", fields)
	}
}

func TestMapBodyPayloadRejectsReservedBodyKeys(t *testing.T) {
	// Intent: A malicious body map must not be able to overwrite local
	// compatibility projections for core promise attribution. Source: DI-mapah
	payloadBytes := mustMapBodyPayloadForTest(t, []mapBodyPayloadField{
		{key: "from", value: "alice"},
	})
	_, err := RelationshipPayloadFields(payloadBytes)
	if err == nil || !strings.Contains(err.Error(), "reserved") {
		t.Fatalf("reserved body key error = %v, want reserved-key rejection", err)
	}
}

func TestMapBodyPayloadRejectsDuplicateBodyKeys(t *testing.T) {
	// Intent: Duplicate CBOR map keys create ambiguous body semantics, so the
	// decoder rejects them instead of letting a later value win. Source: DI-mapah
	payloadBytes := mustMapBodyPayloadForTest(t, []mapBodyPayloadField{
		{key: "subject", value: "bob"},
		{key: "subject", value: "mallory"},
	})
	_, err := RelationshipPayloadFields(payloadBytes)
	if err == nil || !strings.Contains(err.Error(), "duplicated") {
		t.Fatalf("duplicate body key error = %v, want duplicate-key rejection", err)
	}
}

func TestMapBodyPayloadRejectsHistoricalPairListBodies(t *testing.T) {
	// Intent: POC16 deliberately does not preserve the old array-of-pairs body
	// shape because the replacement is a protocol-shape correction, not a
	// backwards-compatible wire migration. Source: DI-mapah
	writer := &cborWriter{}
	if err := writer.writeArrayHeader(5); err != nil {
		t.Fatalf("write payload header: %v", err)
	}
	for _, value := range []string{"mallory", "bob", "local_observation"} {
		if err := writer.writeString(value); err != nil {
			t.Fatalf("write core string: %v", err)
		}
	}
	if err := writer.writeArrayHeader(4); err != nil {
		t.Fatalf("write state header: %v", err)
	}
	for _, value := range []string{"kept", "old pair-list payload", "historical", "turn-old"} {
		if err := writer.writeString(value); err != nil {
			t.Fatalf("write state string: %v", err)
		}
	}
	if err := writer.writeArrayHeader(1); err != nil {
		t.Fatalf("write historical pair-list header: %v", err)
	}
	if err := writer.writeArrayHeader(2); err != nil {
		t.Fatalf("write historical pair entry header: %v", err)
	}
	if err := writer.writeString("subject"); err != nil {
		t.Fatalf("write historical pair key: %v", err)
	}
	if err := writer.writeString("bob"); err != nil {
		t.Fatalf("write historical pair value: %v", err)
	}
	_, err := RelationshipPayloadFields(writer.buffer.Bytes())
	if err == nil || !strings.Contains(err.Error(), "expected cbor major 5") {
		t.Fatalf("historical pair-list error = %v, want map-body rejection", err)
	}
}

func mustMapBodyPayloadForTest(t *testing.T, bodyFields []mapBodyPayloadField) []byte {
	t.Helper()
	writer := &cborWriter{}
	if err := writer.writeArrayHeader(5); err != nil {
		t.Fatalf("write payload header: %v", err)
	}
	for _, value := range []string{"mallory", "bob", "local_observation"} {
		if err := writer.writeString(value); err != nil {
			t.Fatalf("write core string: %v", err)
		}
	}
	if err := writer.writeArrayHeader(4); err != nil {
		t.Fatalf("write state header: %v", err)
	}
	for _, value := range []string{"kept", "map-body payload", "test fixture", "turn-test"} {
		if err := writer.writeString(value); err != nil {
			t.Fatalf("write state string: %v", err)
		}
	}
	if err := writer.writeMapHeader(len(bodyFields)); err != nil {
		t.Fatalf("write body map header: %v", err)
	}
	for _, bodyField := range bodyFields {
		if err := writer.writeString(bodyField.key); err != nil {
			t.Fatalf("write body key: %v", err)
		}
		if err := writer.writeString(bodyField.value); err != nil {
			t.Fatalf("write body value: %v", err)
		}
	}
	return writer.buffer.Bytes()
}
