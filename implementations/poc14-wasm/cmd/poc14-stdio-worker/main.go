package main

import (
	"fmt"
	"io"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc14-wasm/pcid"
	"promisegrid.dev/wire-lab/implementations/poc14-wasm/production"
	"promisegrid.dev/wire-lab/implementations/poc14-wasm/protocol"
	"promisegrid.dev/wire-lab/implementations/poc14-wasm/runtimeadapter"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	return runWithIO(os.Stdin, os.Stdout)
}

// runWithIO executes one stdio worker exchange against explicit streams so the
// same binary behavior can be tested without Docker.
// Intent: Victor's stdio worker must remain a subprocess-I/O PromiseGrid adapter
// rather than an in-process test fake or RPC command handler. Source: DI-sivis
func runWithIO(input io.Reader, output io.Writer) error {
	requestFrameBytes, err := runtimeadapter.ReadCBORFrame(input)
	if err != nil {
		return fmt.Errorf("decode request: %w", err)
	}
	computeRequest, computeParseErr := runtimeadapter.ParseStdioCBOREnvelope(requestFrameBytes)
	if computeParseErr == nil && computeRequest.Type == "compute_request" {
		return runComputeRequest(computeRequest, output)
	}
	request, err := runtimeadapter.ParseStdioCBORRequest(requestFrameBytes)
	if err != nil {
		return fmt.Errorf("parse request: %w", err)
	}
	if request.Type != "promise_request" || request.From == "" || request.To == "" {
		return fmt.Errorf("invalid stdio request")
	}
	registry := pcid.NewRegistry()
	fields := runtimeadapter.PromiseFields(
		request.From,
		request.To,
		runtimeadapter.PromiseAboutStdioAdapter,
		"Victor promises that this worker process sends and receives PromiseGrid envelopes only through stdio.",
	)
	fields["field_protocol"] = pcid.RelationshipV1
	envelope, envelopeErr := protocol.NewEnvelope(registry.MustCID(pcid.RelationshipV1), fields, request.From)
	if envelopeErr != nil {
		return envelopeErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	// Intent: The stdio worker writes exact envelope bytes inside a CBOR byte
	// string, not JSON text or RPC commands; the adapter decides locally
	// whether to forward them through the kernel. Source: DI-linof; DI-kimim
	outboundFrameBytes, err := runtimeadapter.MarshalStdioCBOREnvelope(runtimeadapter.StdioCBOREnvelope{
		Type:          "outbound_envelope",
		From:          request.From,
		To:            request.To,
		Protocol:      pcid.RelationshipV1,
		EnvelopeBytes: envelopeBytes,
	})
	if err != nil {
		return fmt.Errorf("marshal outbound envelope: %w", err)
	}
	if err := runtimeadapter.WriteCBORFrame(output, outboundFrameBytes); err != nil {
		return fmt.Errorf("encode outbound envelope: %w", err)
	}
	ackFrameBytes, err := runtimeadapter.ReadCBORFrame(input)
	if err != nil {
		return fmt.Errorf("decode ack: %w", err)
	}
	ack, err := runtimeadapter.ParseStdioCBORAck(ackFrameBytes)
	if err != nil {
		return fmt.Errorf("parse ack: %w", err)
	}
	if ack.Type != "ack_envelope" || len(ack.EnvelopeBytes) == 0 {
		return fmt.Errorf("invalid ack message")
	}
	ackEnvelope, parseErr := protocol.ParseEnvelope(ack.EnvelopeBytes)
	if parseErr != nil {
		return parseErr
	}
	if verifyErr := protocol.VerifyEnvelope(ackEnvelope); verifyErr != nil {
		return verifyErr
	}
	ackFields, fieldsErr := ackEnvelope.PayloadFields()
	if fieldsErr != nil {
		return fieldsErr
	}
	eventFrameBytes, err := runtimeadapter.MarshalStdioCBOREvent(runtimeadapter.StdioCBOREvent{
		Type:        "ack_event",
		Outcome:     ackFields["outcome"],
		ExactSHA256: protocol.HashExactBytes(ack.EnvelopeBytes),
	})
	if err != nil {
		return fmt.Errorf("marshal ack event: %w", err)
	}
	return runtimeadapter.WriteCBORFrame(output, eventFrameBytes)
}

// runComputeRequest parses Alice's exact compute envelope, keeps the compute
// promise locally, and writes a signed ACK envelope to stdout.
// Intent: The worker, not the adapter, performs Victor's useful compute work
// while preserving the same cid_compute_v1 promise semantics. Source: DI-sivis
func runComputeRequest(request runtimeadapter.StdioCBOREnvelope, output io.Writer) error {
	if request.From == "" || request.To == "" || request.Protocol != pcid.CIDComputeV1 || len(request.EnvelopeBytes) == 0 {
		return fmt.Errorf("invalid stdio compute request")
	}
	registry := pcid.NewRegistry()
	envelope, parseErr := protocol.ParseEnvelope(request.EnvelopeBytes)
	if parseErr != nil {
		return parseErr
	}
	if verifyErr := protocol.VerifyEnvelope(envelope); verifyErr != nil {
		return verifyErr
	}
	fields, fieldsErr := envelope.PayloadFields()
	if fieldsErr != nil {
		return fieldsErr
	}
	if fields["field_promise_about"] != production.PromiseExecuteFunction || fields["from"] != request.From || fields["to"] != request.To {
		return fmt.Errorf("stdio compute request payload does not match request frame")
	}
	inputs, decodeErr := production.DecodeComputeInputs(fields)
	if decodeErr != nil {
		return decodeErr
	}
	if verifyErr := production.VerifyComputeInputCIDs(fields, inputs); verifyErr != nil {
		return verifyErr
	}
	resultBytes, executeErr := production.ExecuteFunction(inputs.FunctionBytes, inputs.InputBytes, inputs.ContextBytes)
	if executeErr != nil {
		return executeErr
	}
	ackFields := production.ExecuteComputePromiseFields(fields, resultBytes)
	ackFields["act"] = "promise"
	ackFields["from"] = request.To
	ackFields["to"] = request.From
	ackFields["outcome"] = "kept"
	ackFields["promise"] = "Victor promises that the stdio worker computed this cid_compute_v1 result from exact function/input/context bytes."
	ackFields["reason"] = "stdio worker compute result returned as a signed pCID-owned array ACK"
	payloadBytes, _, payloadErr := protocol.MarshalKnownArrayPayload(pcid.CIDComputeV1, ackFields)
	if payloadErr != nil {
		return payloadErr
	}
	ackEnvelope, envelopeErr := protocol.NewEnvelopeFromPayload(registry.MustCID(pcid.CIDComputeV1), payloadBytes, request.To)
	if envelopeErr != nil {
		return envelopeErr
	}
	ackBytes, bytesErr := ackEnvelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	ackFrameBytes, marshalErr := runtimeadapter.MarshalStdioCBORAck(runtimeadapter.StdioCBORAck{Type: "compute_ack_envelope", EnvelopeBytes: ackBytes})
	if marshalErr != nil {
		return marshalErr
	}
	// Intent: The worker returns a real signed compute ACK envelope as binary
	// CBOR-framed bytes; Victor's adapter forwards those exact bytes back to
	// Alice instead of translating them into an RPC result. Source: DI-sivis
	return runtimeadapter.WriteCBORFrame(output, ackFrameBytes)
}
