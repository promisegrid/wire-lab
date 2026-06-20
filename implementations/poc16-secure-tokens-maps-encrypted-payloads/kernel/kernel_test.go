package kernel

import (
	"bytes"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/config"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/pcid"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
)

func TestParseTransportEnvelopeDoesNotDecodeNormalPayload(t *testing.T) {
	// Intent: The transport kernel should be able to validate and classify the
	// envelope shell without projecting pCID-owned app payload fields such as
	// `to`; parser roles own that payload interpretation. Source: DI-gazin
	kernel := New(config.Config{}, "test")
	protocolCID := kernel.Protocols.MustCID(pcid.RelationshipV1)
	payloadBytes := []byte{0x83, 0x01, 0x02, 0x03}
	envelope, envelopeErr := protocol.NewEnvelopeFromPayload(protocolCID, payloadBytes, "alice")
	if envelopeErr != nil {
		t.Fatalf("build envelope: %v", envelopeErr)
	}
	frameBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		t.Fatalf("envelope bytes: %v", bytesErr)
	}
	message, parseErr := kernel.parseTransportEnvelope(frameBytes)
	if parseErr != nil {
		t.Fatalf("parse transport envelope: %v", parseErr)
	}
	if message.protocolName != pcid.RelationshipV1 {
		t.Fatalf("protocol name = %s, want %s", message.protocolName, pcid.RelationshipV1)
	}
	if len(message.fields) != 0 {
		t.Fatalf("transport parse projected app payload fields: %#v", message.fields)
	}
	if !bytes.Equal(message.payload, payloadBytes) {
		t.Fatalf("payload bytes changed during transport parse")
	}
}
