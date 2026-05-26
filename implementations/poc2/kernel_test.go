package main

import (
	"net"
	"testing"
)

func TestUnsupportedPCIDGetsObservation(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer func() {
		if closeErr := serverConn.Close(); closeErr != nil {
			t.Log(closeErr)
		}
	}()
	defer func() {
		if closeErr := clientConn.Close(); closeErr != nil {
			t.Log(closeErr)
		}
	}()
	protocolCID, cidErr := HelloProtocolCID()
	if cidErr != nil {
		t.Fatalf("HelloProtocolCID: %v", cidErr)
	}
	evidenceLog, cleanup, logErr := NewEvidenceLog("alice", "")
	if logErr != nil {
		t.Fatalf("NewEvidenceLog: %v", logErr)
	}
	defer func() {
		if cleanupErr := cleanup(); cleanupErr != nil {
			t.Fatalf("cleanup: %v", cleanupErr)
		}
	}()
	kernel := &Kernel{NodeName: "alice", ProtocolCID: protocolCID, EvidenceLog: evidenceLog}
	go kernel.refuseUnsupported(NewFrameConn(serverConn), []byte{0x01, 0x02}, nil)
	observationBytes, readErr := NewFrameConn(clientConn).ReadFrame()
	if readErr != nil {
		t.Fatalf("ReadFrame: %v", readErr)
	}
	observation, parseErr := ParseEnvelope(observationBytes)
	if parseErr != nil {
		t.Fatalf("ParseEnvelope: %v", parseErr)
	}
	kind, fields, kindErr := envelopeKind(observation)
	if kindErr != nil {
		t.Fatalf("envelopeKind: %v", kindErr)
	}
	if kind != "observation_v1" || fields["outcome"] != "refused" {
		t.Fatalf("unexpected observation: kind=%s fields=%#v", kind, fields)
	}
}
