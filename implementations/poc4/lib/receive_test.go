package lib

import (
	"net"
	"testing"
)

func TestWriteReceivePromiseCarriesPromisedProtocol(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	done := make(chan error, 1)
	go func() {
		frameConn := NewFrameConn(clientConn)
		err := WriteReceivePromise(frameConn, NewProtocolCID([]byte("receive spec")), "Alice", "alice-app", NewProtocolCID([]byte("app spec")), "I promise to receive app messages.")
		closeErr := frameConn.Close()
		if err != nil {
			done <- err
			return
		}
		done <- closeErr
	}()
	serverFrame := NewFrameConn(serverConn)
	frameBytes, readErr := serverFrame.ReadFrame()
	if readErr != nil {
		t.Fatalf("ReadFrame returned error: %v", readErr)
	}
	envelope, parseErr := ParseEnvelope(frameBytes)
	if parseErr != nil {
		t.Fatalf("ParseEnvelope returned error: %v", parseErr)
	}
	kind, fields, kindErr := EnvelopeKind(envelope)
	if kindErr != nil {
		t.Fatalf("EnvelopeKind returned error: %v", kindErr)
	}
	if kind != "receive_promise_v1" {
		t.Fatalf("kind = %q", kind)
	}
	if fields["node"] != "Alice" || fields["app"] != "alice-app" {
		t.Fatalf("unexpected receiver fields: %#v", fields)
	}
	if fields["pcid"] == "" {
		t.Fatalf("pcid field was empty")
	}
	if closeErr := serverFrame.Close(); closeErr != nil {
		t.Fatalf("Close returned error: %v", closeErr)
	}
	if err := <-done; err != nil {
		t.Fatalf("writer returned error: %v", err)
	}
}
