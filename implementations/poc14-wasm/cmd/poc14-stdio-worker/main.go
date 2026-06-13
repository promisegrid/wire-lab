package main

import (
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc14-wasm/boundary"
	"promisegrid.dev/wire-lab/implementations/poc14-wasm/pcid"
	"promisegrid.dev/wire-lab/implementations/poc14-wasm/protocol"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	requestFrameBytes, err := boundary.ReadCBORFrame(os.Stdin)
	if err != nil {
		return fmt.Errorf("decode request: %w", err)
	}
	request, err := boundary.ParseStdioCBORRequest(requestFrameBytes)
	if err != nil {
		return fmt.Errorf("parse request: %w", err)
	}
	if request.Type != "promise_request" || request.From == "" || request.To == "" {
		return fmt.Errorf("invalid stdio request")
	}
	registry := pcid.NewRegistry()
	fields := boundary.PromiseFields(
		request.From,
		request.To,
		boundary.PromiseAboutStdioAdapter,
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
	outboundFrameBytes, err := boundary.MarshalStdioCBOREnvelope(boundary.StdioCBOREnvelope{
		Type:          "outbound_envelope",
		From:          request.From,
		To:            request.To,
		Protocol:      pcid.RelationshipV1,
		EnvelopeBytes: envelopeBytes,
	})
	if err != nil {
		return fmt.Errorf("marshal outbound envelope: %w", err)
	}
	if err := boundary.WriteCBORFrame(os.Stdout, outboundFrameBytes); err != nil {
		return fmt.Errorf("encode outbound envelope: %w", err)
	}
	ackFrameBytes, err := boundary.ReadCBORFrame(os.Stdin)
	if err != nil {
		return fmt.Errorf("decode ack: %w", err)
	}
	ack, err := boundary.ParseStdioCBORAck(ackFrameBytes)
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
	eventFrameBytes, err := boundary.MarshalStdioCBOREvent(boundary.StdioCBOREvent{
		Type:        "ack_event",
		Outcome:     ackFields["outcome"],
		ExactSHA256: protocol.HashExactBytes(ack.EnvelopeBytes),
	})
	if err != nil {
		return fmt.Errorf("marshal ack event: %w", err)
	}
	return boundary.WriteCBORFrame(os.Stdout, eventFrameBytes)
}
