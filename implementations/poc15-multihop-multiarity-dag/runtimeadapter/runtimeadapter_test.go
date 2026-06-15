package runtimeadapter

import (
	"bytes"
	"context"
	"testing"
)

func TestValidateWASMModuleAcceptsFixture(t *testing.T) {
	if err := ValidateWASMModule(MinimalWASMModule); err != nil {
		t.Fatalf("validate wasm fixture: %v", err)
	}
}

func TestValidateWASMModuleRejectsBadMagic(t *testing.T) {
	moduleBytes := append([]byte(nil), MinimalWASMModule...)
	moduleBytes[0] = 0xff
	if err := ValidateWASMModule(moduleBytes); err == nil {
		t.Fatalf("bad wasm magic should be rejected")
	}
}

func TestRunWASMModuleExecutesExport(t *testing.T) {
	result, err := RunWASMModule(context.Background(), MinimalWASMModule, ExpectedWASMInput)
	if err != nil {
		t.Fatalf("run wasm module: %v", err)
	}
	if result.ExportName != WASMExportName {
		t.Fatalf("export name = %q, want %q", result.ExportName, WASMExportName)
	}
	if result.InputValue != ExpectedWASMInput {
		t.Fatalf("input value = %d, want %d", result.InputValue, ExpectedWASMInput)
	}
	if result.ExportValue != ExpectedWASMResult {
		t.Fatalf("export value = %d, want %d", result.ExportValue, ExpectedWASMResult)
	}
}

func TestRunWASMModuleRejectsInvalidWASM(t *testing.T) {
	moduleBytes := append([]byte(nil), MinimalWASMModule...)
	moduleBytes[len(moduleBytes)-1] = 0xff
	if _, err := RunWASMModule(context.Background(), moduleBytes, ExpectedWASMInput); err == nil {
		t.Fatalf("invalid wasm body should be rejected")
	}
}

func TestPromiseFieldsUsesSinglePromiseAction(t *testing.T) {
	fields := PromiseFields("victor", "peggy", PromiseAboutStdioAdapter, "Victor promises one stdio-carried envelope.")
	if fields["act"] != "promise" {
		t.Fatalf("act = %q, want promise", fields["act"])
	}
	if fields["field_promise_about"] != PromiseAboutStdioAdapter {
		t.Fatalf("field_promise_about = %q, want %q", fields["field_promise_about"], PromiseAboutStdioAdapter)
	}
}

func TestStdioCBORFrameRoundTripPreservesEnvelopeBytes(t *testing.T) {
	envelopeBytes := []byte{0xd9, 0x67, 0x72, 0x69, 0x64, 0x83, 0xd8, 0x2a}
	encodedEnvelope, err := MarshalStdioCBOREnvelope(StdioCBOREnvelope{
		Type:          "outbound_envelope",
		From:          "victor",
		To:            "peggy",
		Protocol:      "relationship_v1",
		EnvelopeBytes: envelopeBytes,
	})
	if err != nil {
		t.Fatalf("marshal envelope: %v", err)
	}
	var frameBuffer bytes.Buffer
	if err := WriteCBORFrame(&frameBuffer, encodedEnvelope); err != nil {
		t.Fatalf("write frame: %v", err)
	}
	readFrameBytes, err := ReadCBORFrame(&frameBuffer)
	if err != nil {
		t.Fatalf("read frame: %v", err)
	}
	decodedEnvelope, err := ParseStdioCBOREnvelope(readFrameBytes)
	if err != nil {
		t.Fatalf("parse envelope: %v", err)
	}
	if !bytes.Equal(decodedEnvelope.EnvelopeBytes, envelopeBytes) {
		t.Fatalf("envelope bytes changed: %x != %x", decodedEnvelope.EnvelopeBytes, envelopeBytes)
	}
}

func TestStdioCBORMessagesRoundTrip(t *testing.T) {
	requestBytes, err := MarshalStdioCBORRequest(StdioCBORRequest{Type: "promise_request", From: "victor", To: "peggy"})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	request, err := ParseStdioCBORRequest(requestBytes)
	if err != nil {
		t.Fatalf("parse request: %v", err)
	}
	if request.From != "victor" || request.To != "peggy" {
		t.Fatalf("unexpected request: %#v", request)
	}
	ackBytes, err := MarshalStdioCBORAck(StdioCBORAck{Type: "ack_envelope", EnvelopeBytes: []byte{0x01, 0x02, 0x03}})
	if err != nil {
		t.Fatalf("marshal ack: %v", err)
	}
	ack, err := ParseStdioCBORAck(ackBytes)
	if err != nil {
		t.Fatalf("parse ack: %v", err)
	}
	if !bytes.Equal(ack.EnvelopeBytes, []byte{0x01, 0x02, 0x03}) {
		t.Fatalf("unexpected ack bytes: %x", ack.EnvelopeBytes)
	}
	eventBytes, err := MarshalStdioCBOREvent(StdioCBOREvent{Type: "ack_event", Outcome: "kept", ExactSHA256: "abc123"})
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}
	event, err := ParseStdioCBOREvent(eventBytes)
	if err != nil {
		t.Fatalf("parse event: %v", err)
	}
	if event.Outcome != "kept" || event.ExactSHA256 != "abc123" {
		t.Fatalf("unexpected event: %#v", event)
	}
}
