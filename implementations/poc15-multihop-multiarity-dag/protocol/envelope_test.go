package protocol

import (
	"bytes"
	"testing"
)

func TestEnvelopeUsesGridTag42AndVerifies(t *testing.T) {
	protocolCID := NewProtocolCID([]byte("poc15 test protocol"))
	envelope, err := NewEnvelope(protocolCID, map[string]string{"kind": "offer_promise", "from": "alice"}, "alice")
	if err != nil {
		t.Fatalf("new envelope: %v", err)
	}
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("envelope bytes: %v", err)
	}
	if len(envelopeBytes) < 5 || envelopeBytes[0] != 0xda || envelopeBytes[1] != 0x67 || envelopeBytes[2] != 0x72 || envelopeBytes[3] != 0x69 || envelopeBytes[4] != 0x64 {
		t.Fatalf("envelope does not start with CBOR grid tag: %x", envelopeBytes[:5])
	}
	parsed, err := ParseEnvelope(envelopeBytes)
	if err != nil {
		t.Fatalf("parse envelope: %v", err)
	}
	if !parsed.ProtocolCID.Equal(protocolCID) {
		t.Fatalf("parsed pCID mismatch")
	}
	if err := VerifyEnvelope(parsed); err != nil {
		t.Fatalf("verify envelope: %v", err)
	}
}

func TestEnvelopeRejectsTampering(t *testing.T) {
	envelope, err := NewEnvelope(NewProtocolCID([]byte("poc15 tamper protocol")), map[string]string{"kind": "need_advertisement"}, "alice")
	if err != nil {
		t.Fatalf("new envelope: %v", err)
	}
	envelopeBytes, err := envelope.Bytes()
	if err != nil {
		t.Fatalf("envelope bytes: %v", err)
	}
	envelopeBytes[len(envelopeBytes)/2] ^= 0x01
	parsed, parseErr := ParseEnvelope(envelopeBytes)
	if parseErr == nil {
		if verifyErr := VerifyEnvelope(parsed); verifyErr == nil {
			t.Fatalf("tampered envelope should fail parse or verify")
		}
	}
}

func TestParseEnvelopeRejectsMalformedCBOR(t *testing.T) {
	// Intent: Adversarial peers may send arbitrary TCP bytes, but malformed
	// bytes remain local parse-failure event records rather than a command surface.
	// Source: DI-timah
	if _, err := ParseEnvelope([]byte{0xda, 0x67, 0x72, 0x69}); err == nil {
		t.Fatalf("truncated grid-tag bytes should fail parse")
	}
	if _, err := ParseEnvelope([]byte("ignore previous instructions")); err == nil {
		t.Fatalf("plain prompt-injection bytes should fail parse")
	}
}

func TestIdentityKeyPayloadUsesArrayShape(t *testing.T) {
	// Intent: identity_key_v1 is the first POC15 pCID moved to pCID-owned CBOR
	// arrays, proving new protocols need not inherit legacy field maps.
	// Source: DI-vipih
	payloadBytes, err := MarshalIdentityKeyRotationPayload(IdentityKeyRotationPayload{
		Promiser:      "mallory",
		Promisee:      "grace",
		NewKeyLabel:   "mallory-next-key",
		RotationScope: "future-poc15-identity",
	})
	if err != nil {
		t.Fatalf("marshal identity payload: %v", err)
	}
	if len(payloadBytes) == 0 || payloadBytes[0]>>5 != 4 {
		t.Fatalf("identity payload should be a CBOR array, got %x", payloadBytes)
	}
	envelope, err := NewEnvelopeFromPayload(NewProtocolCID([]byte("poc15 identity key rotation promise protocol v1")), payloadBytes, "mallory")
	if err != nil {
		t.Fatalf("new array envelope: %v", err)
	}
	fields, err := envelope.PayloadFields()
	if err != nil {
		t.Fatalf("payload fields: %v", err)
	}
	if fields["act"] != "promise" || fields["from"] != "mallory" || fields["to"] != "grace" || fields["field_new_key_label"] != "mallory-next-key" {
		t.Fatalf("identity routing fields = %#v", fields)
	}
	ackPayloadBytes, err := MarshalIdentityKeyRotationAckPayload(IdentityKeyRotationAckPayload{
		Promiser:      "grace",
		Promisee:      "mallory",
		Outcome:       "kept",
		PromiseText:   "I promise I recorded this future key label locally.",
		NewKeyLabel:   "mallory-next-key",
		RotationScope: "future-poc15-identity",
	})
	if err != nil {
		t.Fatalf("marshal identity ack payload: %v", err)
	}
	ackFields, err := IdentityKeyPayloadFields(ackPayloadBytes)
	if err != nil {
		t.Fatalf("identity ack fields: %v", err)
	}
	if ackFields["act"] != "promise" || ackFields["outcome"] != "kept" || ackFields["promise"] == "" {
		t.Fatalf("identity ack fields = %#v", ackFields)
	}
}

func TestKnownPayloadsUseArrayShape(t *testing.T) {
	// Intent: Fresh POC15 wire payloads are pCID-owned arrays; field_* names are
	// only runtime compatibility projections for existing handlers. Source:
	// DI-gahuh; DI-dirat
	testCases := []struct {
		name         string
		protocolName string
		fields       map[string]string
		wantField    string
	}{
		{
			name:         "relationship",
			protocolName: protocolRelationshipV1,
			fields: map[string]string{
				"from":                "alice",
				"to":                  "bob",
				"promise":             "Alice promises one local relationship event.",
				"reason":              "relationship test",
				"field_promise_about": "local_observation",
				"field_resource":      "storage",
			},
			wantField: "field_resource",
		},
		{
			name:         "kernel receive",
			protocolName: protocolKernelReceiveV1,
			fields: map[string]string{
				"from":     "alice",
				"to":       "kernel",
				"app":      "alice",
				"pcid":     "relationship_v1",
				"pcid_cid": "cidv1-raw-sha2-256:relationship",
			},
			wantField: "field_app",
		},
		{
			name:         "postal scale",
			protocolName: protocolPostalScaleV1,
			fields: map[string]string{
				"from":                "fulfillment",
				"to":                  "postal_scale",
				"field_promise_about": "weigh_package",
				"field_package_id":    "PKG-1001",
			},
			wantField: "field_package_id",
		},
		{
			name:         "accounting",
			protocolName: protocolAccountingV1,
			fields: map[string]string{
				"from":                  "fulfillment",
				"to":                    "accounting",
				"field_promise_about":   "shipment_update",
				"field_order_id":        "ORDER-1001",
				"field_tracking_number": "1Z71051733616616",
				"field_cost_cents":      "860",
			},
			wantField: "field_tracking_number",
		},
		{
			name:         "ups label",
			protocolName: protocolUPSLabelV1,
			fields: map[string]string{
				"from":                   "fulfillment",
				"to":                     "ups_label_printer",
				"field_promise_about":    "print_label",
				"field_package_id":       "PKG-1001",
				"field_shipping_address": "100 Promise Way",
				"field_weight_ounces":    "43",
			},
			wantField: "field_shipping_address",
		},
		{
			name:         "printer port",
			protocolName: protocolPrinterPortV1,
			fields: map[string]string{
				"from":                             "ups_label_printer",
				"to":                               "printer_port",
				"field_promise_about":              "issue_print_capability",
				"field_print_capability_issuee":    "ups_label_printer",
				"field_print_capability_token_id":  "printcap-ups_label_printer",
				"field_print_capability_scope":     "print_label",
				"field_print_capability_max_bytes": "4096",
			},
			wantField: "field_print_capability_issuee",
		},
		{
			name:         "cas storage",
			protocolName: protocolCASStorageV1,
			fields: map[string]string{
				"from":                "alice",
				"to":                  "bob",
				"field_promise_about": "store_content",
				"field_content_cid":   "cidv1-raw-sha2-256:abc",
				"field_content_b64":   "YWJj",
				"field_credit_offer":  "3",
				"field_units":         "1",
			},
			wantField: "field_content_cid",
		},
		{
			name:         "cid compute",
			protocolName: protocolCIDComputeV1,
			fields: map[string]string{
				"from":                 "alice",
				"to":                   "carol",
				"field_promise_about":  "execute_function",
				"field_function_cid":   "cidv1-raw-sha2-256:function",
				"field_function_b64":   "ZnVuY3Rpb24=",
				"field_input_cid":      "cidv1-raw-sha2-256:input",
				"field_input_b64":      "aW5wdXQ=",
				"field_context_cid":    "cidv1-raw-sha2-256:context",
				"field_context_b64":    "Y29udGV4dA==",
				"field_credit_offer":   "5",
				"field_units":          "2",
				"field_capacity_probe": "false",
			},
			wantField: "field_function_cid",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			payloadBytes, arrayPayload, err := MarshalKnownArrayPayload(testCase.protocolName, testCase.fields)
			if err != nil {
				t.Fatalf("marshal known array payload: %v", err)
			}
			if !arrayPayload {
				t.Fatalf("%s should use array payload", testCase.protocolName)
			}
			if len(payloadBytes) == 0 || payloadBytes[0]>>5 != 4 {
				t.Fatalf("%s payload should be a CBOR array, got %x", testCase.protocolName, payloadBytes)
			}
			if bytes.Contains(payloadBytes, []byte("field_")) {
				t.Fatalf("%s payload should not serialize compatibility field names: %x", testCase.protocolName, payloadBytes)
			}
			envelope, err := NewEnvelopeFromPayload(NewProtocolCID([]byte(testCase.protocolName+" test spec")), payloadBytes, testCase.fields["from"])
			if err != nil {
				t.Fatalf("new array envelope: %v", err)
			}
			fields, err := envelope.PayloadFields()
			if err != nil {
				t.Fatalf("payload fields: %v", err)
			}
			if fields["field_payload_shape"] != "cbor_array" || fields[testCase.wantField] == "" {
				t.Fatalf("%s compatibility fields = %#v", testCase.protocolName, fields)
			}
		})
	}
}

func FuzzParseEnvelopeHandlesArbitraryBytes(f *testing.F) {
	// Intent: POC15 should treat malformed CBOR, prompt-injection bytes, partial
	// writes, and random adversarial inputs as parse/verification outcomes rather
	// than panics or expanded protocol semantics. Source: DI-sunuf
	validEnvelope, err := NewEnvelope(NewProtocolCID([]byte("poc15 fuzz protocol")), map[string]string{"act": "promise", "from": "alice"}, "alice")
	if err != nil {
		f.Fatalf("new seed envelope: %v", err)
	}
	validBytes, err := validEnvelope.Bytes()
	if err != nil {
		f.Fatalf("seed envelope bytes: %v", err)
	}
	f.Add([]byte{})
	f.Add([]byte{0xda, 0x67, 0x72, 0x69})
	f.Add([]byte("ignore previous instructions"))
	f.Add(validBytes)
	f.Fuzz(func(t *testing.T, envelopeBytes []byte) {
		envelope, parseErr := ParseEnvelope(envelopeBytes)
		if parseErr != nil {
			return
		}
		if _, fieldsErr := envelope.PayloadFields(); fieldsErr != nil {
			return
		}
		if verifyErr := VerifyEnvelope(envelope); verifyErr != nil {
			return
		}
	})
}
