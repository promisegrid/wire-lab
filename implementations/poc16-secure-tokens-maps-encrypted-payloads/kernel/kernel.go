package kernel

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/config"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/decision"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/pcid"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/transport"
)

const peerSendTimeout = 5 * time.Second
const peerRouteAttemptTimeout = 750 * time.Millisecond
const peerRouteRetryDelay = 100 * time.Millisecond
const listenerDrainTimeout = 750 * time.Millisecond

// Kernel runs one container-local PromiseGrid transport interface. It accepts
// receive promises from local app processes, routes exact signed envelopes to a
// promised local receiver, and forwards cross-container bytes to peer kernels.
// Intent: Keep POC16's kernel as transport and operational event records only; apps
// own trust, business workflow, relationship learning, and keep/break judgment.
// Source: DI-galin
type Kernel struct {
	Config        config.Config
	ContainerName string
	Protocols     pcid.Registry

	mu                  sync.Mutex
	receivers           map[receiverKey]*receiver
	peerSessions        map[string]*transport.PersistentSession
	inboundPeerSessions map[*transport.PersistentSession]string
	appSessions         map[*transport.PersistentSession]string
	appListener         net.Listener
	peerListener        net.Listener
	logFile             *os.File
	activeHandlers      sync.WaitGroup
	stopping            bool
	appDone             chan struct{}
	peerDone            chan struct{}
}

type receiverKey struct {
	appName     string
	protocolCID string
}

type receiver struct {
	appName      string
	protocolName string
	protocolCID  protocol.ProtocolCID
	session      *transport.PersistentSession
}

type parsedEnvelope struct {
	fields       map[string]string
	payload      []byte
	exactHash    string
	protocolCID  protocol.ProtocolCID
	protocolName string
	signer       string
}

// New returns a kernel with an empty local receive-promise table.
func New(cfg config.Config, containerName string) *Kernel {
	return &Kernel{
		Config:              cfg,
		ContainerName:       containerName,
		Protocols:           pcid.NewRegistry(),
		receivers:           make(map[receiverKey]*receiver),
		peerSessions:        make(map[string]*transport.PersistentSession),
		inboundPeerSessions: make(map[*transport.PersistentSession]string),
		appSessions:         make(map[*transport.PersistentSession]string),
	}
}

// Run opens the loopback app listener and Docker-network peer listener until
// the supervisor cancels the container context.
func (kernel *Kernel) Run(ctx context.Context) error {
	if err := kernel.openLog(); err != nil {
		return err
	}
	defer kernel.closeLog()
	if err := kernel.startListeners(ctx); err != nil {
		return err
	}
	<-ctx.Done()
	kernel.closeListeners()
	kernel.closeReceivers()
	kernel.closeAppSessions()
	kernel.closePeerSessions()
	kernel.closeInboundPeerSessions()
	kernel.drainHandlers(ctx)
	return nil
}

func (kernel *Kernel) startListeners(ctx context.Context) error {
	appPort, appPortFound := kernel.Config.KernelAppPortForContainer(kernel.ContainerName)
	if !appPortFound {
		return fmt.Errorf("no app kernel port for container %s", kernel.ContainerName)
	}
	peerPort, peerPortFound := kernel.Config.KernelPeerPortForContainer(kernel.ContainerName)
	if !peerPortFound {
		return fmt.Errorf("no peer kernel port for container %s", kernel.ContainerName)
	}
	appListener, appErr := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(appPort)))
	if appErr != nil {
		return appErr
	}
	peerListener, peerErr := net.Listen("tcp", net.JoinHostPort("", strconv.Itoa(peerPort)))
	if peerErr != nil {
		closeErr := appListener.Close()
		if closeErr != nil {
			kernel.record("app_listener_close_failed", "broken", "", closeErr.Error())
		}
		return peerErr
	}
	kernel.appListener = appListener
	kernel.peerListener = peerListener
	kernel.appDone = make(chan struct{})
	kernel.peerDone = make(chan struct{})
	kernel.record("kernel_started", "kept", "", fmt.Sprintf("app_port=%d peer_port=%d", appPort, peerPort))
	go kernel.acceptLoop(ctx, appListener, kernel.appDone, kernel.handleAppConn)
	go kernel.acceptLoop(ctx, peerListener, kernel.peerDone, kernel.handlePeerConn)
	return nil
}

func (kernel *Kernel) acceptLoop(ctx context.Context, listener net.Listener, done chan struct{}, handler func(net.Conn)) {
	defer close(done)
	for {
		conn, err := listener.Accept()
		if err != nil {
			if kernel.isStopping() || errors.Is(err, net.ErrClosed) {
				kernel.record("kernel_listener_closed", "kept", "", "listener closed during normal shutdown")
				return
			}
			select {
			case <-ctx.Done():
				return
			default:
				kernel.record("kernel_accept_failed", "broken", "", err.Error())
				return
			}
		}
		kernel.activeHandlers.Add(1)
		go func() {
			defer kernel.activeHandlers.Done()
			handler(conn)
		}()
	}
}

func (kernel *Kernel) handleAppConn(conn net.Conn) {
	frameConn := transport.NewFrameConn(conn)
	var session *transport.PersistentSession
	sessionReady := make(chan struct{})
	// Intent: One local app TCP stream now carries receive-promise registration,
	// outbound app requests, and inbound kernel deliveries. The session demuxes
	// replies by parent-linked exact message CIDs, not by RPC IDs. Source:
	// DI-vopab
	session = transport.NewPersistentSession(
		"kernel-app:"+kernel.ContainerName+":"+conn.RemoteAddr().String(),
		frameConn,
		kernel.frameParentExactSHA256s,
		kernel.frameIsResponse,
		func(frameBytes []byte) ([]byte, error) {
			<-sessionReady
			return kernel.handleAppSessionFrame(session, frameBytes)
		},
		func(eventName, outcome, detail string) {
			kernel.record(eventName, outcome, "app", detail)
		},
	)
	close(sessionReady)
	kernel.trackAppSession(session, conn.RemoteAddr().String())
	defer kernel.untrackAppSession(session)
	if kernel.isStopping() {
		if closeErr := session.CloseWithReason(transport.SessionTerminalReasonProcessShutdown); closeErr != nil {
			kernel.record("kernel_app_session_close_failed", "broken", conn.RemoteAddr().String(), closeErr.Error())
		}
	}
	<-session.Done()
}

func (kernel *Kernel) handlePeerConn(conn net.Conn) {
	frameConn := transport.NewFrameConn(conn)
	// Intent: A peer connection accepted during shutdown must not create an
	// unterminated persistent session after the kernel's close pass has already
	// collected active streams. Source: DI-vulit
	if kernel.isStopping() {
		if closeErr := frameConn.Close(); closeErr != nil {
			kernel.record("kernel_inbound_peer_session_close_failed", "broken", conn.RemoteAddr().String(), closeErr.Error())
		}
		return
	}
	// Intent: Peer kernels keep one TCP frame stream open long enough to carry
	// many exact envelopes in either direction; app trust remains outside this
	// transport reuse decision. Source: DI-vopab
	session := transport.NewPersistentSession(
		"kernel-peer-in:"+kernel.ContainerName+":"+conn.RemoteAddr().String(),
		frameConn,
		kernel.frameParentExactSHA256s,
		kernel.frameIsResponse,
		kernel.handlePeerSessionFrame,
		func(eventName, outcome, detail string) {
			kernel.record(eventName, outcome, "peer", detail)
		},
	)
	kernel.trackInboundPeerSession(session, conn.RemoteAddr().String())
	defer kernel.untrackInboundPeerSession(session)
	if kernel.isStopping() {
		if closeErr := session.CloseWithReason(transport.SessionTerminalReasonProcessShutdown); closeErr != nil {
			kernel.record("kernel_inbound_peer_session_close_failed", "broken", conn.RemoteAddr().String(), closeErr.Error())
		}
	}
	<-session.Done()
}

func (kernel *Kernel) handleAppSessionFrame(session *transport.PersistentSession, frameBytes []byte) ([]byte, error) {
	// Intent: The app-side kernel listener is now parser-facing. It may decode
	// only kernel_transport_v1 control promises; normal app payloads must already
	// have been parsed by the separate parser role. Source: DI-gazin
	message, parseErr := kernel.parseTransportEnvelope(frameBytes)
	if parseErr != nil {
		kernel.record("kernel_app_frame_parse_failed", "broken", "", parseErr.Error())
		return nil, parseErr
	}
	if message.protocolName != pcid.KernelTransportV1 {
		kernel.record("kernel_direct_app_payload_rejected", "non_commitment", message.signer, "pcid="+message.protocolName+" transport kernel accepts only "+pcid.KernelTransportV1+" on app-side listener")
		return kernel.notPromisedAck(message, "I promise the transport kernel did not accept this direct app payload because parser-role control is required."), nil
	}
	controlFields, fieldsErr := kernel.kernelTransportFields(message)
	if fieldsErr != nil {
		kernel.record("kernel_transport_control_parse_failed", "broken", message.signer, fieldsErr.Error())
		return nil, fieldsErr
	}
	message.fields = controlFields
	switch controlFields["transport_action"] {
	case "receive_pcid":
		kernel.registerReceiver(session, message)
		return nil, nil
	case "carry_exact_envelope":
		return kernel.routeParserCarriedEnvelope(message)
	default:
		kernel.record("kernel_transport_action_not_promised", "non_commitment", message.signer, "transport_action="+controlFields["transport_action"])
		return kernel.notPromisedAck(message, "I promise the transport kernel did not recognize this parser-role transport action."), nil
	}
}

func (kernel *Kernel) handlePeerSessionFrame(frameBytes []byte) ([]byte, error) {
	// Intent: A peer frame is only an exact signed envelope to deliver locally;
	// this handler does not infer global authority from the peer TCP session.
	// Source: DI-vopab
	message, parseErr := kernel.parseTransportEnvelope(frameBytes)
	if parseErr != nil {
		kernel.record("kernel_peer_frame_parse_failed", "broken", "", parseErr.Error())
		return nil, parseErr
	}
	return kernel.deliverToLocalParser(frameBytes, message), nil
}

func (kernel *Kernel) registerReceiver(session *transport.PersistentSession, message parsedEnvelope) {
	// Intent: A local parser role promises pCID receive capability over its
	// parser/kernel session. The kernel remembers that promise as transport state
	// but does not judge whether later payload promises are trustworthy. Source:
	// DI-vopab; DI-gazin
	appName := firstField(message.fields, "app", "from")
	protocolName := firstField(message.fields, "pcid", "protocol")
	if appName == "" || protocolName == "" {
		kernel.record("kernel_receive_promise_malformed", "malformed", appName, "receive promise requires app and pcid fields")
		return
	}
	protocolCID, known := kernel.Protocols.CID(protocolName)
	if !known || protocolName == pcid.KernelReceiveV1 || protocolName == pcid.KernelTransportV1 {
		kernel.record("kernel_receive_promise_rejected", "non_commitment", appName, "unknown or kernel-internal pCID "+protocolName)
		return
	}
	key := receiverKey{appName: appName, protocolCID: protocolCID.String()}
	kernel.mu.Lock()
	oldReceiver := kernel.receivers[key]
	kernel.receivers[key] = &receiver{
		appName:      appName,
		protocolName: protocolName,
		protocolCID:  protocolCID,
		session:      session,
	}
	kernel.mu.Unlock()
	if oldReceiver != nil && oldReceiver.session != session {
		if closeErr := oldReceiver.session.CloseWithReason(transport.SessionTerminalReasonLocalClose); closeErr != nil {
			kernel.record("kernel_receive_conn_replaced_close_failed", "broken", appName, closeErr.Error())
		}
	}
	kernel.record("kernel_transport_receive_registered", "kept", appName, "pcid="+protocolName+" promise="+message.fields["promise"])
}

func (kernel *Kernel) routeParserCarriedEnvelope(controlMessage parsedEnvelope) ([]byte, error) {
	target := controlMessage.fields["target"]
	if target == "" {
		return kernel.notPromisedAck(controlMessage, "I promise I could not route this parser-carried envelope because the parser named no target app."), nil
	}
	frameBytes, decodeErr := base64.StdEncoding.DecodeString(controlMessage.fields["envelope_b64"])
	if decodeErr != nil {
		kernel.record("kernel_transport_embedded_decode_failed", "broken", controlMessage.signer, decodeErr.Error())
		return kernel.notPromisedAck(controlMessage, "I promise the transport kernel could not decode the parser-supplied exact envelope bytes."), nil
	}
	embeddedMessage, parseErr := kernel.parseTransportEnvelope(frameBytes)
	if parseErr != nil {
		kernel.record("kernel_transport_embedded_parse_failed", "broken", controlMessage.signer, parseErr.Error())
		return kernel.notPromisedAck(controlMessage, "I promise the transport kernel could not parse the parser-supplied exact envelope bytes."), nil
	}
	if expectedHash := controlMessage.fields["embedded_exact_hash"]; expectedHash != "" && expectedHash != embeddedMessage.exactHash {
		kernel.record("kernel_transport_embedded_hash_mismatch", "broken", controlMessage.signer, "expected="+expectedHash+" actual="+embeddedMessage.exactHash)
		return kernel.notPromisedAck(controlMessage, "I promise the transport kernel could not carry an exact envelope whose bytes did not match the parser-supplied hash."), nil
	}
	kernel.record("kernel_transport_carry_requested", "kept", target, "pcid="+embeddedMessage.protocolName+" exact_sha256="+embeddedMessage.exactHash)
	targetContainer, containerFound := kernel.Config.ContainerForAgent(target)
	if !containerFound {
		return kernel.notPromisedAck(embeddedMessage, "I promise I could not route this parser-carried envelope because the target app is unknown."), nil
	}
	if targetContainer == kernel.ContainerName {
		return kernel.deliverToLocalParser(frameBytes, embeddedMessage), nil
	}
	return kernel.forwardToPeerKernel(target, frameBytes)
}

func (kernel *Kernel) forwardToPeerKernel(target string, frameBytes []byte) ([]byte, error) {
	hostName, peerPort, endpointFound := kernel.Config.KernelPeerEndpointForAgent(target)
	if !endpointFound {
		return nil, fmt.Errorf("no peer kernel endpoint for target %s", target)
	}
	endpoint := net.JoinHostPort(hostName, strconv.Itoa(peerPort))
	session, sessionErr := kernel.peerSessionForEndpoint(endpoint, target)
	if sessionErr != nil {
		return nil, sessionErr
	}
	ctx, cancel := context.WithTimeout(context.Background(), peerSendTimeout)
	defer cancel()
	ackBytes, readErr := session.RoundTrip(ctx, protocol.HashExactBytes(frameBytes), frameBytes)
	if readErr != nil {
		kernel.removePeerSession(endpoint, session, target)
		return nil, readErr
	}
	kernel.record("kernel_peer_forwarded", "kept", target, "forwarded exact envelope to peer kernel")
	return ackBytes, nil
}

func (kernel *Kernel) peerSessionForEndpoint(endpoint, target string) (*transport.PersistentSession, error) {
	// Intent: Reuse one outbound peer-kernel stream per endpoint so routing
	// tests exercise persistent TCP behavior while preserving exact message CIDs
	// as the only correlation IDs. Stop checks prevent shutdown races from
	// opening a stream after the shutdown close pass has already collected
	// sessions. Source: DI-vopab; DI-vulit
	kernel.mu.Lock()
	existingSession := kernel.peerSessions[endpoint]
	stopping := kernel.stopping
	kernel.mu.Unlock()
	if stopping {
		return nil, fmt.Errorf("kernel is stopping before opening peer session to %s", target)
	}
	if existingSession != nil {
		return existingSession, nil
	}
	frameConn, dialErr := kernel.dialPeerEndpoint(endpoint)
	if dialErr != nil {
		return nil, dialErr
	}
	session := transport.NewPersistentSession(
		"kernel-peer-out:"+kernel.ContainerName+":"+endpoint,
		frameConn,
		kernel.frameParentExactSHA256s,
		kernel.frameIsResponse,
		kernel.handlePeerSessionFrame,
		func(eventName, outcome, detail string) {
			kernel.record(eventName, outcome, target, detail)
		},
	)
	kernel.mu.Lock()
	defer kernel.mu.Unlock()
	if kernel.stopping {
		if closeErr := session.CloseWithReason(transport.SessionTerminalReasonProcessShutdown); closeErr != nil {
			kernel.record("kernel_peer_session_close_failed", "broken", target, closeErr.Error())
		}
		return nil, fmt.Errorf("kernel is stopping after opening peer session to %s", target)
	}
	if existingSession = kernel.peerSessions[endpoint]; existingSession != nil {
		if closeErr := session.CloseWithReason(transport.SessionTerminalReasonLocalClose); closeErr != nil {
			kernel.record("kernel_peer_duplicate_session_close_failed", "broken", target, closeErr.Error())
		}
		return existingSession, nil
	}
	kernel.peerSessions[endpoint] = session
	return session, nil
}

func (kernel *Kernel) dialPeerEndpoint(endpoint string) (transport.FrameConn, error) {
	deadline := time.Now().Add(peerSendTimeout)
	for {
		frameConn, dialErr := transport.DialFrameConn(endpoint, peerRouteAttemptTimeout)
		if dialErr == nil {
			return frameConn, nil
		}
		if time.Now().Add(peerRouteRetryDelay).After(deadline) {
			return transport.FrameConn{}, dialErr
		}
		// Intent: Startup DNS/container ordering is a transport readiness issue,
		// not an event showing that the peer refused the promise. Retry only the
		// peer-kernel dial for a bounded window, then preserve the original
		// failure as a local route failure. Source: DI-nivon; DI-vopab
		time.Sleep(peerRouteRetryDelay)
	}
}

func (kernel *Kernel) removePeerSession(endpoint string, session *transport.PersistentSession, target string) {
	// Intent: A failed persistent peer stream is a local transport event only;
	// removing it permits a later fresh promise attempt without mutating app
	// trust directly. Source: DI-vopab
	kernel.mu.Lock()
	if kernel.peerSessions[endpoint] == session {
		delete(kernel.peerSessions, endpoint)
	}
	kernel.mu.Unlock()
	if closeErr := session.CloseWithReason(transport.SessionTerminalReasonLocalClose); closeErr != nil {
		kernel.record("kernel_peer_session_close_failed", "broken", target, closeErr.Error())
	}
}

func (kernel *Kernel) trackInboundPeerSession(session *transport.PersistentSession, remoteAddress string) {
	// Intent: Accepted peer streams are owned by the local kernel for lifecycle
	// accounting even though the remote peer initiated the TCP socket. Tracking
	// them here lets clean-run shutdown emit one terminal record per opened
	// inbound session without treating the connection as trust evidence. Source:
	// DI-fobuv
	kernel.mu.Lock()
	kernel.inboundPeerSessions[session] = remoteAddress
	kernel.mu.Unlock()
}

func (kernel *Kernel) untrackInboundPeerSession(session *transport.PersistentSession) {
	kernel.mu.Lock()
	delete(kernel.inboundPeerSessions, session)
	kernel.mu.Unlock()
}

func (kernel *Kernel) trackAppSession(session *transport.PersistentSession, remoteAddress string) {
	// Intent: The parser-role control stream is an app-side kernel session even
	// when no current receiver table entry points at it, so shutdown must own the
	// session directly rather than inferring ownership from routing state.
	// Source: DI-gazin
	kernel.mu.Lock()
	kernel.appSessions[session] = remoteAddress
	kernel.mu.Unlock()
}

func (kernel *Kernel) untrackAppSession(session *transport.PersistentSession) {
	kernel.mu.Lock()
	delete(kernel.appSessions, session)
	kernel.mu.Unlock()
}

func (kernel *Kernel) deliverToLocalParser(frameBytes []byte, message parsedEnvelope) []byte {
	parserName := "parser:" + kernel.ContainerName
	key := receiverKey{appName: parserName, protocolCID: message.protocolCID.String()}
	var receiver *receiver
	deadline := time.Now().Add(peerSendTimeout)
	readinessWaitRecorded := false
	for {
		kernel.mu.Lock()
		receiver = kernel.receivers[key]
		kernel.mu.Unlock()
		if receiver != nil {
			break
		}
		if time.Now().Add(peerRouteRetryDelay).After(deadline) {
			break
		}
		if !readinessWaitRecorded {
			// Intent: Peer frames may arrive while the local parser role is still
			// promising its receive pCIDs during container startup. Waiting here
			// removes startup-order false non-commitments without decoding normal
			// payload routing fields in the transport kernel. Source: DI-darur;
			// DI-gazin
			kernel.record("kernel_parser_receive_readiness_waited", "kept", parserName, "parser receive promise not registered yet for pcid="+message.protocolName)
			readinessWaitRecorded = true
		}
		time.Sleep(peerRouteRetryDelay)
	}
	if receiver == nil {
		kernel.record("kernel_unregistered_parser_pcid", "non_commitment", parserName, "no local parser promised pcid="+message.protocolName)
		return kernel.notPromisedAck(message, "I promise no local parser role has promised to receive this pCID.")
	}
	ctx, cancel := context.WithTimeout(context.Background(), peerSendTimeout)
	defer cancel()
	ackBytes, readErr := receiver.session.RoundTrip(ctx, message.exactHash, frameBytes)
	if readErr != nil {
		kernel.record("kernel_parser_ack_read_failed", "broken", parserName, readErr.Error())
		return kernel.notPromisedAck(message, "I promise local parser-role delivery failed while waiting for its ACK.")
	}
	kernel.record("kernel_parser_delivered", "kept", parserName, "pcid="+message.protocolName+" exact_sha256="+message.exactHash)
	return ackBytes
}

func (kernel *Kernel) parseTransportEnvelope(frameBytes []byte) (parsedEnvelope, error) {
	// Intent: Transport routing inspects only the envelope shell: slot-0 pCID,
	// exact bytes, parent links, and proof signer. Normal slot-1 payload fields
	// remain opaque to this kernel unless the pCID is kernel_transport_v1.
	// Source: DI-gazin
	envelope, parseErr := protocol.ParseEnvelope(frameBytes)
	if parseErr != nil {
		return parsedEnvelope{}, parseErr
	}
	if verifyErr := protocol.VerifyEnvelope(envelope); verifyErr != nil {
		return parsedEnvelope{}, verifyErr
	}
	protocolName, known := kernel.Protocols.Name(envelope.ProtocolCID)
	if !known {
		protocolName = "unknown:" + envelope.ProtocolCID.String()
	}
	return parsedEnvelope{
		fields:       map[string]string{},
		payload:      append([]byte(nil), envelope.Payload...),
		exactHash:    protocol.HashExactBytes(frameBytes),
		protocolCID:  envelope.ProtocolCID,
		protocolName: protocolName,
		signer:       envelope.Proof.Signer,
	}, nil
}

func (kernel *Kernel) kernelTransportFields(message parsedEnvelope) (map[string]string, error) {
	// Intent: kernel_transport_v1 is the only slot-1 payload decoded by the
	// transport kernel. It is local parser/kernel control, not an app protocol.
	// Source: DI-gazin
	fields, fieldsErr := protocol.PayloadFieldsForProtocolName(pcid.KernelTransportV1, message.payload)
	if fieldsErr != nil {
		return nil, fieldsErr
	}
	fields["protocol"] = pcid.KernelTransportV1
	return fields, nil
}

func (kernel *Kernel) frameParentExactSHA256s(frameBytes []byte) ([]string, error) {
	// Intent: Persistent-session demux reads parent links from the signed
	// envelope itself, keeping correlation in the message DAG instead of a
	// separate transport header. Source: DI-vopab
	envelope, parseErr := protocol.ParseEnvelope(frameBytes)
	if parseErr != nil {
		return nil, parseErr
	}
	return append([]string(nil), envelope.ParentExactSHA256s...), nil
}

func (kernel *Kernel) frameIsResponse(frameBytes []byte) (bool, error) {
	// Intent: Kernel persistent sessions should drop unmatched ACK-like responses
	// instead of delivering them to parser roles or peer kernels as new promises.
	// A payload `outcome` field is ACK-like only when the envelope carries a
	// parent link to correlate; fresh promises may discuss outcomes as normal
	// pCID-owned content. This is a bounded demux check, not a route decision:
	// normal payload fields decoded here are never used to choose a target.
	// Source: DI-vopab; DI-gazin
	envelope, parseErr := protocol.ParseEnvelope(frameBytes)
	if parseErr != nil {
		return false, parseErr
	}
	if len(envelope.ParentExactSHA256s) == 0 {
		return false, nil
	}
	if verifyErr := protocol.VerifyEnvelope(envelope); verifyErr != nil {
		return false, verifyErr
	}
	protocolName, known := kernel.Protocols.Name(envelope.ProtocolCID)
	if !known {
		protocolName = "unknown:" + envelope.ProtocolCID.String()
	}
	fields, fieldsErr := fieldsForProtocolName(envelope, protocolName, known)
	if fieldsErr != nil {
		return false, fieldsErr
	}
	return strings.TrimSpace(fields["outcome"]) != "", nil
}

func fieldsForProtocolName(envelope protocol.Envelope, protocolName string, known bool) (map[string]string, error) {
	// Intent: The local kernel routes by the pCID from slot 0 before projecting
	// slot 1 into local compatibility fields, avoiding shape-compatible decoder
	// confusion between protocols. Source: DI-pusak
	if known {
		return protocol.PayloadFieldsForProtocolName(protocolName, envelope.Payload)
	}
	return envelope.PayloadFields()
}

func (kernel *Kernel) notPromisedAck(message parsedEnvelope, promiseText string) []byte {
	target := firstField(message.fields, "from")
	if target == "" {
		target = message.signer
	}
	ackFields := map[string]string{
		"act":     decision.ActPromise,
		"from":    "kernel:" + kernel.ContainerName,
		"to":      target,
		"outcome": "not_promised",
		"promise": promiseText,
		"reason":  "kernel transport non-commitment expressed as local promise content",
	}
	if promiseAbout := message.fields["promise_about"]; promiseAbout != "" {
		ackFields["promise_about"] = promiseAbout
	} else if message.protocolName == pcid.RelationshipV1 || message.protocolName == pcid.KernelTransportV1 {
		ackFields["promise_about"] = "local_observation"
	}
	ack, ackErr := kernel.newAckEnvelope(message, ackFields)
	if ackErr != nil {
		kernel.record("kernel_not_promised_ack_sign_failed", "broken", target, ackErr.Error())
		return []byte{0x00}
	}
	ackBytes, bytesErr := ack.Bytes()
	if bytesErr != nil {
		kernel.record("kernel_not_promised_ack_bytes_failed", "broken", target, bytesErr.Error())
		return []byte{0x00}
	}
	return ackBytes
}

func (kernel *Kernel) newAckEnvelope(message parsedEnvelope, ackFields map[string]string) (protocol.Envelope, error) {
	// Intent: Kernel ACKs are transport non-commitment promises, but their payload
	// still belongs to the original pCID when that pCID has a migrated array
	// encoder. Persistent-session ACKs also parent-link the request exact hash so
	// peers can correlate responses by message DAG rather than RPC IDs. Source:
	// DI-dirat; DI-vopab
	if payloadBytes, arrayPayload, payloadErr := protocol.MarshalKnownArrayPayload(message.protocolName, ackFields); payloadErr == nil && arrayPayload {
		return protocol.NewEnvelopeFromPayloadWithParents(message.protocolCID, payloadBytes, []string{message.exactHash}, "kernel:"+kernel.ContainerName)
	}
	return protocol.NewEnvelopeWithParents(message.protocolCID, ackFields, []string{message.exactHash}, "kernel:"+kernel.ContainerName)
}

func firstField(fields map[string]string, keys ...string) string {
	for _, key := range keys {
		if value := fields[key]; value != "" {
			return value
		}
	}
	return ""
}

func (kernel *Kernel) record(eventName, outcome, peer, detail string) {
	event := decision.Event{
		Observer: "kernel:" + kernel.ContainerName,
		Event:    eventName,
		Outcome:  outcome,
		Peer:     peer,
		Detail:   detail,
	}
	encoded, err := json.Marshal(event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal kernel event: %v\n", err)
		return
	}
	fmt.Println(string(encoded))
	if kernel.logFile != nil {
		if _, writeErr := kernel.logFile.Write(append(encoded, '\n')); writeErr != nil {
			fmt.Fprintf(os.Stderr, "write kernel event: %v\n", writeErr)
		}
	}
}

func (kernel *Kernel) openLog() error {
	runDir := filepath.Join(kernel.Config.RunRoot, kernel.Config.RunID)
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		return err
	}
	logPath := filepath.Join(runDir, "kernel-"+kernel.ContainerName+".jsonl")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	kernel.logFile = logFile
	return nil
}

func (kernel *Kernel) closeLog() {
	if kernel.logFile != nil {
		if err := kernel.logFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "close kernel log: %v\n", err)
		}
	}
}

func (kernel *Kernel) closeListeners() {
	kernel.setStopping()
	if kernel.appListener != nil {
		if err := kernel.appListener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			fmt.Fprintf(os.Stderr, "close app listener: %v\n", err)
		}
	}
	if kernel.peerListener != nil {
		if err := kernel.peerListener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			fmt.Fprintf(os.Stderr, "close peer listener: %v\n", err)
		}
	}
	kernel.waitForListener(kernel.appDone, "app")
	kernel.waitForListener(kernel.peerDone, "peer")
}

func (kernel *Kernel) waitForListener(done chan struct{}, listenerName string) {
	if done == nil {
		return
	}
	timer := time.NewTimer(listenerDrainTimeout)
	defer timer.Stop()
	select {
	case <-done:
	case <-timer.C:
		kernel.record("kernel_listener_close_timeout", "non_commitment", "", listenerName+" listener did not close before timeout")
	}
}

func (kernel *Kernel) closeReceivers() {
	// Intent: The receiver table is routing state promised by parser roles; the
	// actual parser/kernel TCP sessions are closed by closeAppSessions so a
	// session without a current receiver entry still records terminal state.
	// Source: DI-gazin
	kernel.mu.Lock()
	kernel.receivers = make(map[receiverKey]*receiver)
	kernel.mu.Unlock()
}

func (kernel *Kernel) closeAppSessions() {
	// Intent: App-side sessions are now parser-role control streams, not normal
	// app payload streams. Closing every tracked control stream during shutdown
	// keeps persistent-session accounting complete without decoding app payloads.
	// Source: DI-gazin
	kernel.mu.Lock()
	appSessions := make(map[*transport.PersistentSession]string, len(kernel.appSessions))
	for session, remoteAddress := range kernel.appSessions {
		appSessions[session] = remoteAddress
	}
	kernel.appSessions = make(map[*transport.PersistentSession]string)
	kernel.mu.Unlock()
	for session, remoteAddress := range appSessions {
		if session == nil {
			continue
		}
		if closeErr := session.CloseWithReason(transport.SessionTerminalReasonProcessShutdown); closeErr != nil {
			kernel.record("kernel_app_session_close_failed", "broken", remoteAddress, closeErr.Error())
		}
	}
}

func (kernel *Kernel) closePeerSessions() {
	// Intent: Peer-kernel sessions are run-scoped transport promises; closing
	// them at shutdown prevents durable TCP state from leaking into later POC
	// runs. Source: DI-vopab
	kernel.mu.Lock()
	peerSessions := make(map[string]*transport.PersistentSession, len(kernel.peerSessions))
	for endpoint, session := range kernel.peerSessions {
		peerSessions[endpoint] = session
	}
	kernel.peerSessions = make(map[string]*transport.PersistentSession)
	kernel.mu.Unlock()
	for endpoint, session := range peerSessions {
		if session == nil {
			continue
		}
		if closeErr := session.CloseWithReason(transport.SessionTerminalReasonProcessShutdown); closeErr != nil {
			kernel.record("kernel_peer_session_close_failed", "broken", endpoint, closeErr.Error())
		}
	}
}

func (kernel *Kernel) closeInboundPeerSessions() {
	// Intent: Inbound peer-kernel sessions are the accepted side of the same
	// run-scoped transport promise as outbound peer sessions. Closing them during
	// kernel shutdown makes terminal accounting symmetric without changing app
	// trust or routing decisions. Source: DI-fobuv
	kernel.mu.Lock()
	inboundPeerSessions := make(map[*transport.PersistentSession]string, len(kernel.inboundPeerSessions))
	for session, remoteAddress := range kernel.inboundPeerSessions {
		inboundPeerSessions[session] = remoteAddress
	}
	kernel.inboundPeerSessions = make(map[*transport.PersistentSession]string)
	kernel.mu.Unlock()
	for session, remoteAddress := range inboundPeerSessions {
		if session == nil {
			continue
		}
		if closeErr := session.CloseWithReason(transport.SessionTerminalReasonProcessShutdown); closeErr != nil {
			kernel.record("kernel_inbound_peer_session_close_failed", "broken", remoteAddress, closeErr.Error())
		}
	}
}

func (kernel *Kernel) drainHandlers(ctx context.Context) {
	drained := make(chan struct{})
	go func() {
		kernel.activeHandlers.Wait()
		close(drained)
	}()
	timer := time.NewTimer(listenerDrainTimeout)
	defer timer.Stop()
	select {
	case <-drained:
		kernel.record("kernel_handlers_drained", "kept", "", "all active handlers completed before kernel shutdown")
	case <-ctx.Done():
		kernel.record("kernel_handler_drain_cancelled", "non_commitment", "", ctx.Err().Error())
	case <-timer.C:
		kernel.record("kernel_handler_drain_timeout", "non_commitment", "", "some kernel handlers may still be running")
	}
}

func (kernel *Kernel) setStopping() {
	kernel.mu.Lock()
	defer kernel.mu.Unlock()
	kernel.stopping = true
}

func (kernel *Kernel) isStopping() bool {
	kernel.mu.Lock()
	defer kernel.mu.Unlock()
	return kernel.stopping
}
