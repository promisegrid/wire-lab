package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

// ParentExtractor returns exact request hashes named by an incoming frame. The
// transport package treats those hashes as opaque message IDs so PromiseGrid
// message CIDs remain the only correlation layer.
type ParentExtractor func(frameBytes []byte) ([]string, error)

// ResponseClassifier distinguishes ACK/response frames from fresh incoming
// promise frames when no local pending request matches their parent links.
type ResponseClassifier func(frameBytes []byte) (bool, error)

// InboundFrameHandler handles a fresh incoming promise frame and returns the
// exact ACK/response frame that should be written on the same persistent session.
type InboundFrameHandler func(frameBytes []byte) ([]byte, error)

// SessionRecorder records local session events without making transport code
// depend on the POC runtime or kernel event packages.
type SessionRecorder func(eventName, outcome, detail string)

type sessionResult struct {
	frameBytes []byte
	err        error
}

type pendingFrame struct {
	requestID string
	result    chan sessionResult
}

// PersistentSession carries many length-prefixed PromiseGrid frames over one TCP
// connection and correlates replies by parent-linked request message CID.
// Intent: POC15 tests production-shaped long-lived transport while keeping
// correlation in the message DAG instead of adding RPC request IDs. Source:
// DI-vopab
type PersistentSession struct {
	name           string
	frameConn      FrameConn
	extractParents ParentExtractor
	isResponse     ResponseClassifier
	handleInbound  InboundFrameHandler
	record         SessionRecorder

	writeMu sync.Mutex

	mu      sync.Mutex
	pending map[string]pendingFrame
	closed  bool

	closeOnce sync.Once
	done      chan struct{}
}

// NewPersistentSession starts a read loop for one already-open TCP frame stream.
func NewPersistentSession(name string, frameConn FrameConn, extractParents ParentExtractor, isResponse ResponseClassifier, handleInbound InboundFrameHandler, record SessionRecorder) *PersistentSession {
	session := &PersistentSession{
		name:           name,
		frameConn:      frameConn,
		extractParents: extractParents,
		isResponse:     isResponse,
		handleInbound:  handleInbound,
		record:         record,
		pending:        make(map[string]pendingFrame),
		done:           make(chan struct{}),
	}
	session.recordEvent("persistent_session_opened", "kept", "session="+name)
	go session.readLoop()
	return session
}

// Send writes a fresh frame without registering a pending response. Receive
// promise registrations use this path because the promise itself establishes the
// future receive stream; no semantic ACK is required for registration.
func (session *PersistentSession) Send(ctx context.Context, frameBytes []byte) error {
	if err := session.checkOpen(); err != nil {
		return err
	}
	return session.writeFrame(ctx, frameBytes)
}

// RoundTrip writes a request frame and waits for the first response whose parent
// links include requestID.
func (session *PersistentSession) RoundTrip(ctx context.Context, requestID string, frameBytes []byte) ([]byte, error) {
	if requestID == "" {
		return nil, fmt.Errorf("persistent session request id is required")
	}
	pending := pendingFrame{requestID: requestID, result: make(chan sessionResult, 1)}
	if err := session.addPending(pending); err != nil {
		return nil, err
	}
	if err := session.writeFrame(ctx, frameBytes); err != nil {
		session.removePending(requestID)
		return nil, err
	}
	session.recordEvent("persistent_session_reused", "kept", "session="+session.name+" request="+requestID)
	select {
	case result := <-pending.result:
		if result.err != nil {
			return nil, result.err
		}
		return result.frameBytes, nil
	case <-ctx.Done():
		session.removePending(requestID)
		return nil, ctx.Err()
	case <-session.done:
		session.removePending(requestID)
		return nil, fmt.Errorf("persistent session %s closed before response", session.name)
	}
}

// Close closes the underlying TCP stream and releases all pending requests.
func (session *PersistentSession) Close() error {
	var closeErr error
	session.closeOnce.Do(func() {
		session.mu.Lock()
		if !session.closed {
			session.closed = true
			for requestID, pending := range session.pending {
				delete(session.pending, requestID)
				pending.result <- sessionResult{err: fmt.Errorf("persistent session %s closed", session.name)}
			}
		}
		session.mu.Unlock()
		closeErr = session.frameConn.Close()
		close(session.done)
		session.recordEvent("persistent_session_closed", "kept", "session="+session.name)
	})
	return closeErr
}

// Done closes when the persistent session read loop has ended. Intent: Kernel
// accept handlers wait on this channel instead of reading raw frames themselves,
// so the persistent-session object remains the sole framing/demux owner. Source:
// DI-vopab
func (session *PersistentSession) Done() <-chan struct{} {
	return session.done
}

func (session *PersistentSession) readLoop() {
	for {
		frameBytes, readErr := session.frameConn.ReadFrame()
		if readErr != nil {
			if !session.isClosed() {
				if isGracefulSessionClose(readErr) {
					session.recordEvent("persistent_session_remote_closed", "kept", "session="+session.name)
					if closeErr := session.Close(); closeErr != nil {
						session.recordEvent("persistent_session_close_failed", "broken", "session="+session.name+" "+closeErr.Error())
					}
					return
				}
				session.recordEvent("persistent_session_read_failed", "broken", "session="+session.name+" "+readErr.Error())
				if closeErr := session.Close(); closeErr != nil {
					session.recordEvent("persistent_session_close_failed", "broken", "session="+session.name+" "+closeErr.Error())
				}
			}
			return
		}
		if session.isResponse != nil {
			responseFrame, classifyErr := session.isResponse(frameBytes)
			if classifyErr != nil {
				session.recordEvent("persistent_session_response_classify_failed", "broken", "session="+session.name+" "+classifyErr.Error())
			}
			if responseFrame {
				if session.deliverPending(frameBytes) {
					continue
				}
				session.dropUnmatchedResponse()
				continue
			}
		} else if session.deliverPending(frameBytes) {
			// Intent: Legacy tests and adapters without a response classifier keep
			// parent-link-only demux, but production POC15 sessions require the
			// classifier so child promises parented to a request are not mistaken
			// for that request's response. Source: DI-vopab
			continue
		}
		if session.handleInbound == nil {
			session.recordEvent("persistent_session_unmatched_frame", "non_commitment", "session="+session.name)
			continue
		}
		session.handleFreshFrame(frameBytes)
	}
}

func (session *PersistentSession) handleFreshFrame(frameBytes []byte) {
	copiedFrame := append([]byte(nil), frameBytes...)
	go func() {
		// Intent: A fresh inbound frame may cause the app to send another promise
		// over the same persistent session before it can ACK the inbound frame.
		// Running handlers outside the read loop keeps response demux live and
		// avoids reintroducing one TCP connection per nested message. Source:
		// DI-vopab
		responseBytes, handleErr := session.handleInbound(copiedFrame)
		if handleErr != nil {
			session.recordEvent("persistent_session_inbound_failed", "broken", "session="+session.name+" "+handleErr.Error())
			return
		}
		if len(responseBytes) == 0 {
			return
		}
		if writeErr := session.writeFrame(context.Background(), responseBytes); writeErr != nil {
			session.recordEvent("persistent_session_response_write_failed", "broken", "session="+session.name+" "+writeErr.Error())
		}
	}()
}

func (session *PersistentSession) deliverPending(frameBytes []byte) bool {
	if session.extractParents == nil {
		return false
	}
	parentIDs, parentErr := session.extractParents(frameBytes)
	if parentErr != nil {
		return false
	}
	for _, parentID := range parentIDs {
		pending, found := session.takePending(parentID)
		if !found {
			continue
		}
		pending.result <- sessionResult{frameBytes: append([]byte(nil), frameBytes...)}
		return true
	}
	return false
}

func (session *PersistentSession) dropUnmatchedResponse() {
	// Intent: ACK-like frames whose parent links do not match a local pending
	// request are terminal unmatched responses, not new incoming promises to ACK
	// again. This prevents persistent sessions from creating ACK-of-ACK storms
	// while preserving message-CID correlation as the only demux key. Source:
	// DI-vopab
	session.recordEvent("persistent_session_unmatched_response", "non_commitment", "session="+session.name)
}

func (session *PersistentSession) writeFrame(ctx context.Context, frameBytes []byte) error {
	if ctx == nil {
		ctx = context.Background()
	}
	done := make(chan error, 1)
	go func() {
		session.writeMu.Lock()
		defer session.writeMu.Unlock()
		done <- session.frameConn.WriteFrame(frameBytes)
	}()
	select {
	case err := <-done:
		if err != nil {
			if closeErr := session.Close(); closeErr != nil {
				session.recordEvent("persistent_session_close_failed", "broken", "session="+session.name+" "+closeErr.Error())
			}
		}
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-session.done:
		return fmt.Errorf("persistent session %s closed before write", session.name)
	}
}

func (session *PersistentSession) addPending(pending pendingFrame) error {
	session.mu.Lock()
	defer session.mu.Unlock()
	if session.closed {
		return fmt.Errorf("persistent session %s is closed", session.name)
	}
	if _, exists := session.pending[pending.requestID]; exists {
		return fmt.Errorf("persistent session %s already has pending request %s", session.name, pending.requestID)
	}
	session.pending[pending.requestID] = pending
	return nil
}

func (session *PersistentSession) removePending(requestID string) {
	session.mu.Lock()
	defer session.mu.Unlock()
	delete(session.pending, requestID)
}

func (session *PersistentSession) takePending(requestID string) (pendingFrame, bool) {
	session.mu.Lock()
	defer session.mu.Unlock()
	pending, found := session.pending[requestID]
	if found {
		delete(session.pending, requestID)
	}
	return pending, found
}

func (session *PersistentSession) checkOpen() error {
	session.mu.Lock()
	defer session.mu.Unlock()
	if session.closed {
		return fmt.Errorf("persistent session %s is closed", session.name)
	}
	return nil
}

func (session *PersistentSession) isClosed() bool {
	session.mu.Lock()
	defer session.mu.Unlock()
	return session.closed
}

func (session *PersistentSession) recordEvent(eventName, outcome, detail string) {
	if session.record != nil {
		session.record(eventName, outcome, detail)
	}
}

func isGracefulSessionClose(err error) bool {
	// Intent: A peer closing a run-scoped persistent TCP session during shutdown
	// is not a broken promise by either app; only unexpected frame/read failures
	// should be recorded as broken transport events. Source: DI-vopab
	return errors.Is(err, io.EOF)
}
