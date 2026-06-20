package parserrole

import (
	"net"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/config"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/pcid"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/transport"
)

func TestFrameIsResponseRequiresParentLink(t *testing.T) {
	// Intent: Parser roles must not swallow fresh pCID-owned promises merely
	// because their payload discusses an outcome; only parent-linked ACK frames
	// are session responses. Source: DI-gazin
	role := New(config.Config{}, "test")
	fields := map[string]string{
		"act":           "promise",
		"from":          "alice",
		"to":            "bob",
		"promise_about": "relationship_repair",
		"outcome":       "non_commitment",
		"promise":       "I promise to describe my local repair outcome.",
		"reason":        "fresh promise with outcome vocabulary",
	}
	payloadBytes, arrayPayload, payloadErr := protocol.MarshalKnownArrayPayload(pcid.RelationshipV1, fields)
	if payloadErr != nil {
		t.Fatalf("marshal relationship payload: %v", payloadErr)
	}
	if !arrayPayload {
		t.Fatalf("relationship_v1 test payload was not encoded as an array")
	}
	freshEnvelope, freshErr := protocol.NewEnvelopeFromPayload(role.Protocols.MustCID(pcid.RelationshipV1), payloadBytes, "alice")
	if freshErr != nil {
		t.Fatalf("build fresh envelope: %v", freshErr)
	}
	freshBytes, freshBytesErr := freshEnvelope.Bytes()
	if freshBytesErr != nil {
		t.Fatalf("fresh envelope bytes: %v", freshBytesErr)
	}
	freshIsResponse, freshResponseErr := role.frameIsResponse(freshBytes)
	if freshResponseErr != nil {
		t.Fatalf("classify fresh envelope: %v", freshResponseErr)
	}
	if freshIsResponse {
		t.Fatalf("fresh outcome-bearing promise was misclassified as a response")
	}
	parentHash := protocol.HashExactBytes([]byte("parent promise"))
	responseEnvelope, responseErr := protocol.NewEnvelopeFromPayloadWithParents(role.Protocols.MustCID(pcid.RelationshipV1), payloadBytes, []string{parentHash}, "alice")
	if responseErr != nil {
		t.Fatalf("build response envelope: %v", responseErr)
	}
	responseBytes, responseBytesErr := responseEnvelope.Bytes()
	if responseBytesErr != nil {
		t.Fatalf("response envelope bytes: %v", responseBytesErr)
	}
	responseIsResponse, responseClassifyErr := role.frameIsResponse(responseBytes)
	if responseClassifyErr != nil {
		t.Fatalf("classify response envelope: %v", responseClassifyErr)
	}
	if !responseIsResponse {
		t.Fatalf("parent-linked outcome payload was not classified as a response")
	}
}

func TestAppReceiverRegistryReplacesReceiver(t *testing.T) {
	// Intent: A later local app receive promise should replace the earlier
	// parser-role receiver for the same app and pCID, because receive promises
	// are local routing state rather than immutable global registration. Source:
	// DI-gazin
	registry := NewAppReceiverRegistry()
	protocolCID := pcid.NewRegistry().MustCID(pcid.ProductionShippingV1)
	firstSession := &transport.PersistentSession{}
	secondSession := &transport.PersistentSession{}

	registry.Register("fulfillment", pcid.ProductionShippingV1, protocolCID, firstSession)
	registry.Register("fulfillment", pcid.ProductionShippingV1, protocolCID, secondSession)

	receiver := registry.Lookup("fulfillment", protocolCID)
	if receiver == nil {
		t.Fatalf("receiver missing after replacement")
	}
	if receiver.session != secondSession {
		t.Fatalf("receiver session was not replaced")
	}
}

func TestProtocolParserRejectsMalformedPayload(t *testing.T) {
	// Intent: Malformed pCID-owned payload bytes should be rejected by the
	// parser role, not projected into a fake local routing record. Source:
	// DI-gazin
	role := New(config.Config{}, "test")
	envelope, envelopeErr := protocol.NewEnvelopeFromPayload(role.Protocols.MustCID(pcid.ProductionShippingV1), []byte{0x81, 0x61, 0x78}, "alice")
	if envelopeErr != nil {
		t.Fatalf("build malformed-payload envelope: %v", envelopeErr)
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		t.Fatalf("malformed-payload envelope bytes: %v", bytesErr)
	}

	if _, parseErr := role.Parser.Parse(envelopeBytes); parseErr == nil {
		t.Fatalf("malformed production_shipping_v1 payload should be rejected")
	}
}

func TestRegisterAppReceiverRejectsParserInternalPCID(t *testing.T) {
	// Intent: Apps promise application pCID receive capability to the parser
	// role; parser/kernel control pCIDs must stay internal to the parser role and
	// transport kernel. Source: DI-gazin
	role := New(config.Config{}, "test")
	message := ParsedMessage{
		Fields: map[string]string{
			"app":  "alice",
			"pcid": pcid.KernelTransportV1,
		},
		ProtocolCID:  role.Protocols.MustCID(pcid.KernelTransportV1),
		ProtocolName: pcid.KernelTransportV1,
	}

	if registerErr := role.registerAppReceiver(nil, message); registerErr != nil {
		t.Fatalf("internal-pCID receive rejection should be local non-commitment: %v", registerErr)
	}
	if receiver := role.Receivers.Lookup("alice", role.Protocols.MustCID(pcid.KernelTransportV1)); receiver != nil {
		t.Fatalf("parser-internal pCID was registered as an app receiver")
	}
}

func TestDeliverToLocalAppReturnsNotPromisedWhenNoReceiver(t *testing.T) {
	// Intent: If no local app promised a pCID to the parser role, the parser role
	// returns a parent-linked non-commitment promise rather than imposing
	// delivery on any app. Source: DI-gazin
	role := New(config.Config{}, "test")
	frameBytes, message := buildProductionShippingEnvelope(t, role, "alice", "fulfillment")

	ackBytes := role.deliverToLocalApp(frameBytes, message)
	ackMessage, parseErr := role.Parser.Parse(ackBytes)
	if parseErr != nil {
		t.Fatalf("parse not-promised ACK: %v", parseErr)
	}
	if ackMessage.Fields["outcome"] != "not_promised" {
		t.Fatalf("ack outcome = %q, want not_promised", ackMessage.Fields["outcome"])
	}
	if ackMessage.Fields["parent_exact_sha256"] != message.ExactHash {
		t.Fatalf("ack parent = %q, want %q", ackMessage.Fields["parent_exact_sha256"], message.ExactHash)
	}
}

func TestDeliverToLocalAppReturnsAppAck(t *testing.T) {
	// Intent: Exact-envelope delivery should cross the parser/app session and
	// return the app's parent-linked ACK bytes unchanged to the parser role's
	// caller. Source: DI-gazin
	role := New(config.Config{}, "test")
	appConn, parserConn := net.Pipe()
	appSession := transport.NewPersistentSession(
		"test-app",
		transport.NewFrameConn(appConn),
		frameParentExactSHA256s,
		role.frameIsResponse,
		func(frameBytes []byte) ([]byte, error) {
			incoming, parseErr := role.Parser.Parse(frameBytes)
			if parseErr != nil {
				return nil, parseErr
			}
			ackFields := map[string]string{
				"act":           "promise",
				"from":          "fulfillment",
				"to":            incoming.Fields["from"],
				"outcome":       "kept",
				"promise_about": incoming.Fields["promise_about"],
				"promise":       "I promise the local app received the exact envelope.",
				"reason":        "test app ACK",
			}
			payloadBytes, _, payloadErr := protocol.MarshalKnownArrayPayload(pcid.ProductionShippingV1, ackFields)
			if payloadErr != nil {
				return nil, payloadErr
			}
			ackEnvelope, ackErr := protocol.NewEnvelopeFromPayloadWithParents(role.Protocols.MustCID(pcid.ProductionShippingV1), payloadBytes, []string{incoming.ExactHash}, "fulfillment")
			if ackErr != nil {
				return nil, ackErr
			}
			return ackEnvelope.Bytes()
		},
		func(eventName, outcome, detail string) {},
	)
	parserSession := transport.NewPersistentSession(
		"test-parser",
		transport.NewFrameConn(parserConn),
		frameParentExactSHA256s,
		role.frameIsResponse,
		nil,
		func(eventName, outcome, detail string) {},
	)
	defer closeTestSession(t, appSession)
	defer closeTestSession(t, parserSession)

	protocolCID := role.Protocols.MustCID(pcid.ProductionShippingV1)
	role.Receivers.Register("fulfillment", pcid.ProductionShippingV1, protocolCID, parserSession)
	frameBytes, message := buildProductionShippingEnvelope(t, role, "alice", "fulfillment")

	ackBytes := role.deliverToLocalApp(frameBytes, message)
	ackMessage, parseErr := role.Parser.Parse(ackBytes)
	if parseErr != nil {
		t.Fatalf("parse app ACK: %v", parseErr)
	}
	if ackMessage.Fields["outcome"] != "kept" {
		t.Fatalf("app ACK outcome = %q, want kept", ackMessage.Fields["outcome"])
	}
	if ackMessage.Fields["parent_exact_sha256"] != message.ExactHash {
		t.Fatalf("app ACK parent = %q, want %q", ackMessage.Fields["parent_exact_sha256"], message.ExactHash)
	}
}

func buildProductionShippingEnvelope(t *testing.T, role *ParserRole, from, to string) ([]byte, ParsedMessage) {
	t.Helper()
	fields := map[string]string{
		"act":           "promise",
		"from":          from,
		"to":            to,
		"promise_about": "address_lookup",
		"order_id":      "ORDER-1",
		"exchange_id":   "EX-1",
		"promise":       "I promise to request an address lookup through production_shipping_v1.",
		"reason":        "parser-role test",
	}
	payloadBytes, _, payloadErr := protocol.MarshalKnownArrayPayload(pcid.ProductionShippingV1, fields)
	if payloadErr != nil {
		t.Fatalf("marshal production_shipping_v1 payload: %v", payloadErr)
	}
	envelope, envelopeErr := protocol.NewEnvelopeFromPayload(role.Protocols.MustCID(pcid.ProductionShippingV1), payloadBytes, from)
	if envelopeErr != nil {
		t.Fatalf("build production_shipping_v1 envelope: %v", envelopeErr)
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		t.Fatalf("production_shipping_v1 envelope bytes: %v", bytesErr)
	}
	message, parseErr := role.Parser.Parse(envelopeBytes)
	if parseErr != nil {
		t.Fatalf("parse production_shipping_v1 envelope: %v", parseErr)
	}
	return envelopeBytes, message
}

func closeTestSession(t *testing.T, session *transport.PersistentSession) {
	t.Helper()
	if closeErr := session.CloseWithReason(transport.SessionTerminalReasonProcessShutdown); closeErr != nil {
		t.Fatalf("close test session: %v", closeErr)
	}
}
