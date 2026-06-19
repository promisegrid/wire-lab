package main

import (
	"bytes"
	"encoding/base64"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/pcid"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/production"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/runtimeadapter"
)

func TestRunWithIOComputesCIDPromise(t *testing.T) {
	registry := pcid.NewRegistry()
	functionBytes := production.SampleFunctionBytes()
	inputBytes := production.SampleInputBytes()
	contextBytes := production.SampleContextBytes()
	fields := map[string]string{
		"act":           "promise",
		"from":          "alice",
		"to":            "victor",
		"turn":          "test",
		"promise":       "Alice promises to receive Victor's stdio compute result.",
		"reason":        "test stdio compute promise",
		"promise_about": production.PromiseExecuteFunction,
		"function_cid":  production.ContentCID(functionBytes),
		"function_b64":  base64.StdEncoding.EncodeToString(functionBytes),
		"input_cid":     production.ContentCID(inputBytes),
		"input_b64":     base64.StdEncoding.EncodeToString(inputBytes),
		"context_cid":   production.ContentCID(contextBytes),
		"context_b64":   base64.StdEncoding.EncodeToString(contextBytes),
		"credit_offer":  "5",
		"units":         "1",
	}
	payloadBytes, _, payloadErr := protocol.MarshalKnownArrayPayload(pcid.CIDComputeV1, fields)
	if payloadErr != nil {
		t.Fatalf("marshal request payload: %v", payloadErr)
	}
	envelope, envelopeErr := protocol.NewEnvelopeFromPayload(registry.MustCID(pcid.CIDComputeV1), payloadBytes, "alice")
	if envelopeErr != nil {
		t.Fatalf("new request envelope: %v", envelopeErr)
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		t.Fatalf("request bytes: %v", bytesErr)
	}
	requestBytes, requestErr := runtimeadapter.MarshalStdioCBOREnvelope(runtimeadapter.StdioCBOREnvelope{
		Type:          "compute_request",
		From:          "alice",
		To:            "victor",
		Protocol:      pcid.CIDComputeV1,
		EnvelopeBytes: envelopeBytes,
	})
	if requestErr != nil {
		t.Fatalf("marshal stdio request: %v", requestErr)
	}
	var input bytes.Buffer
	if err := runtimeadapter.WriteCBORFrame(&input, requestBytes); err != nil {
		t.Fatalf("write stdio request: %v", err)
	}
	var output bytes.Buffer
	if err := runWithIO(&input, &output); err != nil {
		t.Fatalf("run worker: %v", err)
	}
	ackFrameBytes, frameErr := runtimeadapter.ReadCBORFrame(&output)
	if frameErr != nil {
		t.Fatalf("read ack frame: %v", frameErr)
	}
	ack, ackErr := runtimeadapter.ParseStdioCBORAck(ackFrameBytes)
	if ackErr != nil {
		t.Fatalf("parse ack: %v", ackErr)
	}
	if ack.Type != "compute_ack_envelope" {
		t.Fatalf("ack type = %q, want compute_ack_envelope", ack.Type)
	}
	ackEnvelope, parseErr := protocol.ParseEnvelope(ack.EnvelopeBytes)
	if parseErr != nil {
		t.Fatalf("parse ack envelope: %v", parseErr)
	}
	if verifyErr := protocol.VerifyEnvelope(ackEnvelope); verifyErr != nil {
		t.Fatalf("verify ack envelope: %v", verifyErr)
	}
	ackFields, fieldsErr := ackEnvelope.PayloadFields()
	if fieldsErr != nil {
		t.Fatalf("ack fields: %v", fieldsErr)
	}
	if ackFields["from"] != "victor" || ackFields["to"] != "alice" || ackFields["result_cid"] == "" {
		t.Fatalf("unexpected ack fields: %#v", ackFields)
	}
	expectedBytes, expectedErr := production.ExecuteFunction(functionBytes, inputBytes, contextBytes)
	if expectedErr != nil {
		t.Fatalf("expected compute: %v", expectedErr)
	}
	if ackFields["result_cid"] != production.ContentCID(expectedBytes) {
		t.Fatalf("result cid = %s, want %s", ackFields["result_cid"], production.ContentCID(expectedBytes))
	}
}
