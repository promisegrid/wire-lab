package transport

import (
	"context"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestPersistentSessionMatchesOutOfOrderParentReplies(t *testing.T) {
	// Intent: Persistent sessions must allow multiple outstanding PromiseGrid
	// messages and match replies by parent-linked exact message CID, not by frame
	// order or an RPC request number. Source: DI-vopab
	serverConn, clientConn := net.Pipe()
	defer closeTestConn(t, serverConn, "server")
	defer closeTestConn(t, clientConn, "client")
	setTestDeadline(t, serverConn)
	setTestDeadline(t, clientConn)
	serverFrames := NewFrameConn(serverConn)
	sessionEvents := newSessionEventLog()
	session := NewPersistentSession("test-out-of-order", NewFrameConn(clientConn), testParentExtractor, testResponseClassifier, nil, sessionEvents.record)
	defer closeTestSession(t, session)

	serverDone := make(chan error, 1)
	go func() {
		firstFrame, firstErr := serverFrames.ReadFrame()
		if firstErr != nil {
			serverDone <- firstErr
			return
		}
		secondFrame, secondErr := serverFrames.ReadFrame()
		if secondErr != nil {
			serverDone <- secondErr
			return
		}
		firstRequest := string(firstFrame)
		secondRequest := string(secondFrame)
		if !sameRequestSet(firstRequest, secondRequest, "first", "second") {
			serverDone <- errUnexpectedFrames(firstRequest, secondRequest)
			return
		}
		if writeErr := serverFrames.WriteFrame([]byte("ack:" + secondRequest)); writeErr != nil {
			serverDone <- writeErr
			return
		}
		serverDone <- serverFrames.WriteFrame([]byte("ack:" + firstRequest))
	}()

	firstResult := make(chan sessionTestResult, 1)
	secondResult := make(chan sessionTestResult, 1)
	go roundTripForTest(session, "first", []byte("first"), firstResult)
	go roundTripForTest(session, "second", []byte("second"), secondResult)

	if err := <-serverDone; err != nil {
		t.Fatalf("server exchange: %v", err)
	}
	assertSessionResult(t, <-firstResult, "ack:first")
	assertSessionResult(t, <-secondResult, "ack:second")
	if sessionEvents.count("persistent_session_reused") != 2 {
		t.Fatalf("persistent_session_reused count = %d, want 2", sessionEvents.count("persistent_session_reused"))
	}
}

func TestPersistentSessionWritesInboundResponseOnSameConnection(t *testing.T) {
	// Intent: A fresh inbound frame should be handled on the long-lived stream
	// and answered on that same stream while the read loop remains available for
	// unrelated replies. Source: DI-vopab
	serverConn, clientConn := net.Pipe()
	defer closeTestConn(t, serverConn, "server")
	defer closeTestConn(t, clientConn, "client")
	setTestDeadline(t, serverConn)
	setTestDeadline(t, clientConn)
	serverFrames := NewFrameConn(serverConn)
	session := NewPersistentSession("test-inbound", NewFrameConn(clientConn), testParentExtractor, testResponseClassifier, func(frameBytes []byte) ([]byte, error) {
		return []byte("ack:" + string(frameBytes)), nil
	}, nil)
	defer closeTestSession(t, session)

	if err := serverFrames.WriteFrame([]byte("fresh")); err != nil {
		t.Fatalf("write fresh frame: %v", err)
	}
	responseBytes, readErr := serverFrames.ReadFrame()
	if readErr != nil {
		t.Fatalf("read inbound response: %v", readErr)
	}
	if string(responseBytes) != "ack:fresh" {
		t.Fatalf("response = %q, want ack:fresh", string(responseBytes))
	}
}

func TestPersistentSessionDropsUnmatchedResponses(t *testing.T) {
	// Intent: An ACK-like frame whose parent does not match this session's
	// pending requests must not be handled as a fresh inbound promise, otherwise
	// two persistent peers can create an ACK-of-ACK loop. Source: DI-vopab
	serverConn, clientConn := net.Pipe()
	defer closeTestConn(t, serverConn, "server")
	defer closeTestConn(t, clientConn, "client")
	setTestDeadline(t, serverConn)
	setTestDeadline(t, clientConn)
	serverFrames := NewFrameConn(serverConn)
	sessionEvents := newSessionEventLog()
	handledFrames := make(chan []byte, 1)
	session := NewPersistentSession("test-unmatched-response", NewFrameConn(clientConn), testParentExtractor, testResponseClassifier, func(frameBytes []byte) ([]byte, error) {
		handledFrames <- append([]byte(nil), frameBytes...)
		return []byte("ack:" + string(frameBytes)), nil
	}, sessionEvents.record)
	defer closeTestSession(t, session)

	if err := serverFrames.WriteFrame([]byte("ack:unknown-request")); err != nil {
		t.Fatalf("write unmatched response: %v", err)
	}
	select {
	case handledFrame := <-handledFrames:
		t.Fatalf("unmatched response was handled as fresh frame: %s", string(handledFrame))
	case <-time.After(100 * time.Millisecond):
	}
	if sessionEvents.count("persistent_session_unmatched_response") != 1 {
		t.Fatalf("persistent_session_unmatched_response count = %d, want 1", sessionEvents.count("persistent_session_unmatched_response"))
	}
}

func TestPersistentSessionTreatsParentLinkedChildAsInbound(t *testing.T) {
	// Intent: A child promise can legitimately cite the original request as a
	// parent, but it is still a fresh promise rather than the response to that
	// request. Response classification must therefore gate pending delivery.
	// Source: DI-vopab
	serverConn, clientConn := net.Pipe()
	defer closeTestConn(t, serverConn, "server")
	defer closeTestConn(t, clientConn, "client")
	setTestDeadline(t, serverConn)
	setTestDeadline(t, clientConn)
	serverFrames := NewFrameConn(serverConn)
	handledFrames := make(chan []byte, 1)
	session := NewPersistentSession("test-parent-linked-child", NewFrameConn(clientConn), testParentExtractor, testResponseClassifier, func(frameBytes []byte) ([]byte, error) {
		handledFrames <- append([]byte(nil), frameBytes...)
		return []byte("ack:" + string(frameBytes)), nil
	}, nil)
	defer closeTestSession(t, session)

	serverDone := make(chan error, 1)
	go func() {
		requestFrame, readErr := serverFrames.ReadFrame()
		if readErr != nil {
			serverDone <- readErr
			return
		}
		if string(requestFrame) != "first" {
			serverDone <- errUnexpectedFrames(string(requestFrame), "missing second frame")
			return
		}
		if writeErr := serverFrames.WriteFrame([]byte("child:first")); writeErr != nil {
			serverDone <- writeErr
			return
		}
		childAck, ackErr := serverFrames.ReadFrame()
		if ackErr != nil {
			serverDone <- ackErr
			return
		}
		if string(childAck) != "ack:child:first" {
			serverDone <- errUnexpectedFrames(string(childAck), "ack:child:first")
			return
		}
		serverDone <- serverFrames.WriteFrame([]byte("ack:first"))
	}()

	result := make(chan sessionTestResult, 1)
	go roundTripForTest(session, "first", []byte("first"), result)
	if err := <-serverDone; err != nil {
		t.Fatalf("server exchange: %v", err)
	}
	select {
	case handledFrame := <-handledFrames:
		if string(handledFrame) != "child:first" {
			t.Fatalf("handled frame = %q, want child:first", string(handledFrame))
		}
	default:
		t.Fatalf("parent-linked child frame was not handled as inbound")
	}
	assertSessionResult(t, <-result, "ack:first")
}

type sessionTestResult struct {
	frameBytes []byte
	err        error
}

type sessionEventLog struct {
	mu     sync.Mutex
	counts map[string]int
}

func newSessionEventLog() *sessionEventLog {
	return &sessionEventLog{counts: make(map[string]int)}
}

func (eventLog *sessionEventLog) record(eventName, outcome, detail string) {
	eventLog.mu.Lock()
	defer eventLog.mu.Unlock()
	eventLog.counts[eventName]++
}

func (eventLog *sessionEventLog) count(eventName string) int {
	eventLog.mu.Lock()
	defer eventLog.mu.Unlock()
	return eventLog.counts[eventName]
}

func testParentExtractor(frameBytes []byte) ([]string, error) {
	frameText := string(frameBytes)
	switch {
	case strings.HasPrefix(frameText, "ack:"):
		return []string{strings.TrimPrefix(frameText, "ack:")}, nil
	case strings.HasPrefix(frameText, "child:"):
		return []string{strings.TrimPrefix(frameText, "child:")}, nil
	default:
		return nil, nil
	}
}

func testResponseClassifier(frameBytes []byte) (bool, error) {
	return strings.HasPrefix(string(frameBytes), "ack:"), nil
}

func roundTripForTest(session *PersistentSession, requestID string, frameBytes []byte, result chan<- sessionTestResult) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	responseBytes, err := session.RoundTrip(ctx, requestID, frameBytes)
	result <- sessionTestResult{frameBytes: responseBytes, err: err}
}

func assertSessionResult(t *testing.T, result sessionTestResult, want string) {
	t.Helper()
	if result.err != nil {
		t.Fatalf("round trip: %v", result.err)
	}
	if string(result.frameBytes) != want {
		t.Fatalf("round trip response = %q, want %s", string(result.frameBytes), want)
	}
}

func closeTestSession(t *testing.T, session *PersistentSession) {
	t.Helper()
	if err := session.Close(); err != nil {
		t.Logf("close persistent session: %v", err)
	}
}

func setTestDeadline(t *testing.T, conn net.Conn) {
	t.Helper()
	if err := conn.SetDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatalf("set test connection deadline: %v", err)
	}
}

func errUnexpectedFrames(firstFrame, secondFrame string) error {
	return &unexpectedFramesError{firstFrame: firstFrame, secondFrame: secondFrame}
}

func sameRequestSet(firstFrame, secondFrame, wantFirst, wantSecond string) bool {
	return (firstFrame == wantFirst && secondFrame == wantSecond) || (firstFrame == wantSecond && secondFrame == wantFirst)
}

type unexpectedFramesError struct {
	firstFrame  string
	secondFrame string
}

func (err *unexpectedFramesError) Error() string {
	return "unexpected request frames: " + err.firstFrame + ", " + err.secondFrame
}
