package kernel

import (
	"bytes"
	"context"
	"net"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc3/lib"
)

func TestKernelDeliversToPromisedReceiver(t *testing.T) {
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
	protocolCID := lib.NewProtocolCID([]byte("test delivery protocol"))
	evidenceBuffer := &bytes.Buffer{}
	testKernel := &Kernel{
		NodeName:    "bob",
		EvidenceLog: lib.NewEvidenceLog("bob", evidenceBuffer),
		receivers: map[string]Receiver{
			protocolCID.String(): {AppName: "bob-test-app", Conn: lib.NewFrameConn(serverConn)},
		},
	}
	envelope, envelopeErr := lib.NewEnvelope(protocolCID, map[string]string{
		"kind": "test_v1",
		"to":   "bob",
		"text": "hello",
	})
	if envelopeErr != nil {
		t.Fatalf("NewEnvelope: %v", envelopeErr)
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		t.Fatalf("Envelope.Bytes: %v", bytesErr)
	}
	delivered := make(chan []byte, 1)
	go func() {
		frameBytes, readErr := lib.NewFrameConn(clientConn).ReadFrame()
		if readErr != nil {
			t.Errorf("ReadFrame: %v", readErr)
			return
		}
		delivered <- frameBytes
	}()
	testKernel.deliverToLocalReceiver(context.Background(), envelopeBytes, envelope)
	if string(<-delivered) != string(envelopeBytes) {
		t.Fatalf("delivered bytes did not match")
	}
	if !bytes.Contains(evidenceBuffer.Bytes(), []byte("app_deliver")) {
		t.Fatalf("missing app_deliver evidence: %s", evidenceBuffer.String())
	}
}

func TestKernelRecordsNotPromisedWithoutReceiver(t *testing.T) {
	protocolCID := lib.NewProtocolCID([]byte("unsupported local protocol"))
	evidenceBuffer := &bytes.Buffer{}
	testKernel := &Kernel{
		NodeName:    "bob",
		EvidenceLog: lib.NewEvidenceLog("bob", evidenceBuffer),
		receivers:   map[string]Receiver{},
	}
	envelope, envelopeErr := lib.NewEnvelope(protocolCID, map[string]string{
		"kind": "test_v1",
		"to":   "bob",
		"text": "hello",
	})
	if envelopeErr != nil {
		t.Fatalf("NewEnvelope: %v", envelopeErr)
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		t.Fatalf("Envelope.Bytes: %v", bytesErr)
	}
	testKernel.deliverToLocalReceiver(context.Background(), envelopeBytes, envelope)
	if !bytes.Contains(evidenceBuffer.Bytes(), []byte("not-promised")) {
		t.Fatalf("missing not-promised evidence: %s", evidenceBuffer.String())
	}
}
