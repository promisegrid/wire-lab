package main

import (
	"encoding/hex"
	"encoding/json"
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
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)
	var request boundary.StdioRequest
	if err := decoder.Decode(&request); err != nil {
		return fmt.Errorf("decode request: %w", err)
	}
	if request.Type != "promise_request" || request.From == "" || request.To == "" {
		return fmt.Errorf("invalid stdio request")
	}
	registry := pcid.NewRegistry()
	fields := boundary.PromiseFields(
		request.From,
		request.To,
		boundary.PromiseAboutStdioBoundary,
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
	// Intent: The stdio worker writes exact envelope bytes as data, not commands;
	// the adapter decides locally whether to forward them through the kernel.
	// Source: DI-linof
	if err := encoder.Encode(boundary.StdioEnvelopeMessage{
		Type:     "outbound_envelope",
		From:     request.From,
		To:       request.To,
		Protocol: pcid.RelationshipV1,
		Hex:      hex.EncodeToString(envelopeBytes),
	}); err != nil {
		return fmt.Errorf("encode outbound envelope: %w", err)
	}
	var ack boundary.StdioAckMessage
	if err := decoder.Decode(&ack); err != nil {
		return fmt.Errorf("decode ack: %w", err)
	}
	if ack.Type != "ack_envelope" || ack.Hex == "" {
		return fmt.Errorf("invalid ack message")
	}
	ackBytes, decodeErr := hex.DecodeString(ack.Hex)
	if decodeErr != nil {
		return decodeErr
	}
	ackEnvelope, parseErr := protocol.ParseEnvelope(ackBytes)
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
	return encoder.Encode(boundary.StdioObservedMessage{
		Type:        "ack_observed",
		Outcome:     ackFields["outcome"],
		ExactSHA256: protocol.HashExactBytes(ackBytes),
	})
}
