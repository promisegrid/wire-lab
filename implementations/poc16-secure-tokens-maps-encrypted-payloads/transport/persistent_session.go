package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

// ParentExtractor returns exact request CIDs named by an incoming frame. The
// transport package treats those CIDs as opaque message IDs so PromiseGrid
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

const sessionFrameCloseTimeout = 250 * time.Millisecond

type sessionResult struct {
	frameBytes []byte
	err        error
}

type pendingFrame struct {
	requestID string
	result    chan sessionResult
}

// SessionTerminalReason names why a run-scoped persistent stream ended.
// Intent: POC16 needs strict terminal accounting for every session object so
// analyzer output can distinguish shutdown, peer EOF, and failed I/O without
// turning transport events into peer trust judgments. Source: DI-mapop
type SessionTerminalReason string

const (
	SessionTerminalReasonLocalClose      SessionTerminalReason = "local_close"
	SessionTerminalReasonRemoteEOF       SessionTerminalReason = "remote_eof"
	SessionTerminalReasonReadFailed      SessionTerminalReason = "read_failed"
	SessionTerminalReasonWriteFailed     SessionTerminalReason = "write_failed"
	SessionTerminalReasonProcessShutdown SessionTerminalReason = "process_shutdown"
)

// SessionStats is a snapshot of one session object's local counters.
// Intent: The counters are local transport accounting only; PromiseGrid message
// CIDs and parent links remain the semantic correlation mechanism. Source:
// DI-mapop
type SessionStats struct {
	SessionID      string
	Name           string
	FramesSent     int
	FramesReceived int
	Requests       int
	Responses      int
}

var persistentSessionCounter uint64

// PersistentSession carries many length-prefixed PromiseGrid frames over one TCP
// connection and correlates replies by parent-linked request message CID.
// Intent: POC16 tests production-shaped long-lived transport while keeping
// correlation in the message DAG instead of adding RPC request IDs. Source:
// DI-vopab
type PersistentSession struct {
	sessionID      string
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
	stats   SessionStats

	closeOnce sync.Once
	done      chan struct{}
}

// NewPersistentSession starts a read loop for one already-open TCP frame stream.
func NewPersistentSession(name string, frameConn FrameConn, extractParents ParentExtractor, isResponse ResponseClassifier, handleInbound InboundFrameHandler, record SessionRecorder) *PersistentSession {
	sessionID := newSessionID(name)
	session := &PersistentSession{
		sessionID:      sessionID,
		name:           name,
		frameConn:      frameConn,
		extractParents: extractParents,
		isResponse:     isResponse,
		handleInbound:  handleInbound,
		record:         record,
		pending:        make(map[string]pendingFrame),
		stats: SessionStats{
			SessionID: sessionID,
			Name:      name,
		},
		done: make(chan struct{}),
	}
	session.recordEvent("persistent_session_opened", "kept", session.detailString(""))
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
	session.recordSessionRequestStarted(requestID)
	session.recordEvent("persistent_session_reused", "kept", session.detailString("request="+requestID))
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

// Close closes the session, releases pending requests, and then attempts to
// close the underlying TCP stream.
func (session *PersistentSession) Close() error {
	return session.CloseWithReason(SessionTerminalReasonLocalClose)
}

// CloseWithReason closes the stream and records exactly one terminal reason.
func (session *PersistentSession) CloseWithReason(reason SessionTerminalReason) error {
	var closeErr error
	session.closeOnce.Do(func() {
		// Intent: Emit local lifecycle accounting before the underlying socket
		// close, because shutdown validation cares that this agent reached a
		// terminal state even if a busy TCP close stalls or returns late. Source:
		// DI-kunad
		session.mu.Lock()
		if !session.closed {
			session.closed = true
			for requestID, pending := range session.pending {
				delete(session.pending, requestID)
				pending.result <- sessionResult{err: fmt.Errorf("persistent session %s closed", session.name)}
			}
		}
		session.mu.Unlock()
		close(session.done)
		session.recordEvent("persistent_session_closed", "kept", session.detailString("reason="+string(reason)))
		session.recordSessionTerminal(reason)
		// Intent: Do not let one slow TCP close serialize or hide terminal records
		// for other persistent sessions during process shutdown. Local lifecycle
		// accounting is already complete above; the actual socket close is still
		// attempted and either its failure or deferral is recorded. Source:
		// DI-kunad; DI-nuriv
		closeResult := make(chan error, 1)
		go func() {
			closeResult <- session.frameConn.Close()
		}()
		timer := time.NewTimer(sessionFrameCloseTimeout)
		defer timer.Stop()
		select {
		case closeErr = <-closeResult:
			if closeErr != nil {
				session.recordEvent("persistent_session_close_failed", "broken", session.detailString(closeErr.Error()))
			}
		case <-timer.C:
			session.recordEvent("persistent_session_close_deferred", "non_commitment", session.detailString("underlying TCP close exceeded local close timeout"))
		}
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
					session.recordEvent("persistent_session_remote_closed", "kept", session.detailString(""))
					if closeErr := session.CloseWithReason(SessionTerminalReasonRemoteEOF); closeErr != nil {
						session.recordEvent("persistent_session_close_failed", "broken", session.detailString(closeErr.Error()))
					}
					return
				}
				session.recordEvent("persistent_session_read_failed", "broken", session.detailString(readErr.Error()))
				if closeErr := session.CloseWithReason(SessionTerminalReasonReadFailed); closeErr != nil {
					session.recordEvent("persistent_session_close_failed", "broken", session.detailString(closeErr.Error()))
				}
			}
			return
		}
		session.recordSessionFrameReceived(len(frameBytes))
		if session.isResponse != nil {
			responseFrame, classifyErr := session.isResponse(frameBytes)
			if classifyErr != nil {
				session.recordEvent("persistent_session_response_classify_failed", "broken", session.detailString(classifyErr.Error()))
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
			// parent-link-only demux, but production POC16 sessions require the
			// classifier so child promises parented to a request are not mistaken
			// for that request's response. Source: DI-vopab
			continue
		}
		if session.handleInbound == nil {
			session.recordEvent("persistent_session_unmatched_frame", "non_commitment", session.detailString(""))
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
			session.recordEvent("persistent_session_inbound_failed", "broken", session.detailString(handleErr.Error()))
			return
		}
		if len(responseBytes) == 0 {
			return
		}
		if writeErr := session.writeFrame(context.Background(), responseBytes); writeErr != nil {
			session.recordEvent("persistent_session_response_write_failed", "broken", session.detailString(writeErr.Error()))
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
		session.recordSessionResponseMatched(parentID)
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
	session.recordEvent("persistent_session_unmatched_response", "non_commitment", session.detailString(""))
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
			if closeErr := session.CloseWithReason(SessionTerminalReasonWriteFailed); closeErr != nil {
				session.recordEvent("persistent_session_close_failed", "broken", session.detailString(closeErr.Error()))
			}
		} else {
			session.recordSessionFrameSent(len(frameBytes))
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

func (session *PersistentSession) recordSessionFrameSent(frameBytes int) {
	session.mu.Lock()
	session.stats.FramesSent++
	stats := session.stats
	session.mu.Unlock()
	session.recordEvent("persistent_session_frame_sent", "kept", session.detailString(fmt.Sprintf("bytes=%d frames_sent=%d", frameBytes, stats.FramesSent)))
}

func (session *PersistentSession) recordSessionFrameReceived(frameBytes int) {
	session.mu.Lock()
	session.stats.FramesReceived++
	stats := session.stats
	session.mu.Unlock()
	session.recordEvent("persistent_session_frame_received", "kept", session.detailString(fmt.Sprintf("bytes=%d frames_received=%d", frameBytes, stats.FramesReceived)))
}

func (session *PersistentSession) recordSessionRequestStarted(requestID string) {
	session.mu.Lock()
	session.stats.Requests++
	stats := session.stats
	session.mu.Unlock()
	session.recordEvent("persistent_session_request_started", "kept", session.detailString(fmt.Sprintf("request=%s requests=%d", requestID, stats.Requests)))
}

func (session *PersistentSession) recordSessionResponseMatched(requestID string) {
	session.mu.Lock()
	session.stats.Responses++
	stats := session.stats
	session.mu.Unlock()
	session.recordEvent("persistent_session_response_matched", "kept", session.detailString(fmt.Sprintf("request=%s responses=%d", requestID, stats.Responses)))
}

func (session *PersistentSession) recordSessionTerminal(reason SessionTerminalReason) {
	session.mu.Lock()
	stats := session.stats
	session.mu.Unlock()
	detail := fmt.Sprintf("reason=%s frames_sent=%d frames_received=%d requests=%d responses=%d", reason, stats.FramesSent, stats.FramesReceived, stats.Requests, stats.Responses)
	session.recordEvent("persistent_session_terminal", "kept", session.detailString(detail))
}

func (session *PersistentSession) detailString(extra string) string {
	detail := "session=" + session.name + " session_id=" + session.sessionID
	if extra != "" {
		detail += " " + extra
	}
	return detail
}

func newSessionID(name string) string {
	// Intent: Session IDs are local log correlation names only. The counter
	// disambiguates reconnects with the same endpoint name without creating a
	// second hash-like printable identifier outside the PromiseGrid CID rules.
	// Source: DI-mapop; DI-sazip
	sequence := atomic.AddUint64(&persistentSessionCounter, 1)
	return fmt.Sprintf("session-%06d", sequence)
}

func isGracefulSessionClose(err error) bool {
	// Intent: A peer closing a run-scoped persistent TCP session during shutdown
	// is not a broken promise by either app; only unexpected frame/read failures
	// should be recorded as broken transport events. Source: DI-vopab
	return errors.Is(err, io.EOF)
}
