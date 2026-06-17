package kernel

import (
	"context"
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

	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/config"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/decision"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/pcid"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/protocol"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/transport"
)

const peerSendTimeout = 5 * time.Second
const peerRouteAttemptTimeout = 750 * time.Millisecond
const peerRouteRetryDelay = 100 * time.Millisecond
const listenerDrainTimeout = 750 * time.Millisecond

// Kernel runs one container-local PromiseGrid transport interface. It accepts
// receive promises from local app processes, routes exact signed envelopes to a
// promised local receiver, and forwards cross-container bytes to peer kernels.
// Intent: Keep POC15's kernel as transport and operational event records only; apps
// own trust, business workflow, relationship learning, and keep/break judgment.
// Source: DI-galin
type Kernel struct {
	Config        config.Config
	ContainerName string
	Protocols     pcid.Registry

	mu             sync.Mutex
	receivers      map[receiverKey]*receiver
	peerSessions   map[string]*transport.PersistentSession
	appListener    net.Listener
	peerListener   net.Listener
	logFile        *os.File
	activeHandlers sync.WaitGroup
	stopping       bool
	appDone        chan struct{}
	peerDone       chan struct{}
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
	exactHash    string
	protocolCID  protocol.ProtocolCID
	protocolName string
}

// New returns a kernel with an empty local receive-promise table.
func New(cfg config.Config, containerName string) *Kernel {
	return &Kernel{
		Config:        cfg,
		ContainerName: containerName,
		Protocols:     pcid.NewRegistry(),
		receivers:     make(map[receiverKey]*receiver),
		peerSessions:  make(map[string]*transport.PersistentSession),
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
	kernel.closePeerSessions()
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
	<-session.Done()
}

func (kernel *Kernel) handlePeerConn(conn net.Conn) {
	frameConn := transport.NewFrameConn(conn)
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
	<-session.Done()
}

func (kernel *Kernel) handleAppSessionFrame(session *transport.PersistentSession, frameBytes []byte) ([]byte, error) {
	// Intent: Kernel receive registrations update only the local routing table;
	// all other app frames are routed as exact promise envelopes on the same
	// persistent session. Source: DI-vopab
	message, parseErr := kernel.parseEnvelope(frameBytes)
	if parseErr != nil {
		kernel.record("kernel_app_frame_parse_failed", "broken", "", parseErr.Error())
		return nil, parseErr
	}
	if message.protocolName == pcid.KernelReceiveV1 {
		kernel.registerReceiver(session, message)
		return nil, nil
	}
	ackBytes, routeErr := kernel.routeFromLocalApp(frameBytes, message)
	if routeErr != nil {
		kernel.record("kernel_route_failed", "broken", message.fields["to"], routeErr.Error())
		ackBytes = kernel.notPromisedAck(message, "I promise I could not route this exact envelope through the local kernel.")
	}
	return ackBytes, nil
}

func (kernel *Kernel) handlePeerSessionFrame(frameBytes []byte) ([]byte, error) {
	// Intent: A peer frame is only an exact signed envelope to deliver locally;
	// this handler does not infer global authority from the peer TCP session.
	// Source: DI-vopab
	message, parseErr := kernel.parseEnvelope(frameBytes)
	if parseErr != nil {
		kernel.record("kernel_peer_frame_parse_failed", "broken", "", parseErr.Error())
		return nil, parseErr
	}
	return kernel.deliverToLocalApp(frameBytes, message), nil
}

func (kernel *Kernel) registerReceiver(session *transport.PersistentSession, message parsedEnvelope) {
	// Intent: A local app promises pCID receive capability over its existing
	// app/kernel session. The kernel remembers that promise as routing state but
	// does not judge whether later payload promises are trustworthy. Source:
	// DI-vopab
	appName := firstField(message.fields, "app", "from")
	protocolName := firstField(message.fields, "pcid", "protocol")
	if appName == "" || protocolName == "" {
		kernel.record("kernel_receive_promise_malformed", "malformed", appName, "receive promise requires app and pcid fields")
		return
	}
	protocolCID, known := kernel.Protocols.CID(protocolName)
	if !known || protocolName == pcid.KernelReceiveV1 {
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
		if closeErr := oldReceiver.session.Close(); closeErr != nil {
			kernel.record("kernel_receive_conn_replaced_close_failed", "broken", appName, closeErr.Error())
		}
	}
	kernel.record("app_receive_promise_registered", "kept", appName, "pcid="+protocolName+" promise="+message.fields["promise"])
}

func (kernel *Kernel) routeFromLocalApp(frameBytes []byte, message parsedEnvelope) ([]byte, error) {
	target := message.fields["to"]
	if target == "" {
		return kernel.notPromisedAck(message, "I promise I could not route this envelope because it names no target app."), nil
	}
	targetContainer, containerFound := kernel.Config.ContainerForAgent(target)
	if !containerFound {
		return kernel.notPromisedAck(message, "I promise I could not route this envelope because the target app is unknown."), nil
	}
	if targetContainer == kernel.ContainerName {
		return kernel.deliverToLocalApp(frameBytes, message), nil
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
	// as the only correlation IDs. Source: DI-vopab
	kernel.mu.Lock()
	existingSession := kernel.peerSessions[endpoint]
	kernel.mu.Unlock()
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
	if existingSession = kernel.peerSessions[endpoint]; existingSession != nil {
		if closeErr := session.Close(); closeErr != nil {
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
	if closeErr := session.Close(); closeErr != nil {
		kernel.record("kernel_peer_session_close_failed", "broken", target, closeErr.Error())
	}
}

func (kernel *Kernel) deliverToLocalApp(frameBytes []byte, message parsedEnvelope) []byte {
	target := message.fields["to"]
	key := receiverKey{appName: target, protocolCID: message.protocolCID.String()}
	targetAgent, agentFound := kernel.Config.Agent(target)
	if !agentFound {
		kernel.record("kernel_target_app_unknown", "non_commitment", target, "no configured local app named "+target)
		return kernel.notPromisedAck(message, "I promise no configured local app matches this target.")
	}
	configuredReceivePromise := false
	for _, protocolName := range targetAgent.Protocols() {
		if protocolName == message.protocolName {
			configuredReceivePromise = true
			break
		}
	}
	if !configuredReceivePromise {
		kernel.record("kernel_target_pcid_not_configured", "non_commitment", target, "target app is not configured to promise pcid="+message.protocolName)
		return kernel.notPromisedAck(message, "I promise the configured target app does not promise to receive this pCID.")
	}
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
			// Intent: A target app listed as supporting this pCID may still be
			// registering its receive promise during container startup. Waiting here
			// removes startup-order false non-commitments without retrying any
			// app-level semantic ACK. Source: DI-darur
			kernel.record("kernel_app_receive_readiness_waited", "kept", target, "configured receive promise not registered yet for pcid="+message.protocolName)
			readinessWaitRecorded = true
		}
		time.Sleep(peerRouteRetryDelay)
	}
	if receiver == nil {
		kernel.record("kernel_unregistered_pcid", "non_commitment", target, "no local app promised pcid="+message.protocolName)
		return kernel.notPromisedAck(message, "I promise no local app has promised to receive this pCID.")
	}
	ctx, cancel := context.WithTimeout(context.Background(), peerSendTimeout)
	defer cancel()
	ackBytes, readErr := receiver.session.RoundTrip(ctx, message.exactHash, frameBytes)
	if readErr != nil {
		kernel.record("kernel_app_ack_read_failed", "broken", target, readErr.Error())
		return kernel.notPromisedAck(message, "I promise local app delivery failed while waiting for app event.")
	}
	kernel.record("kernel_app_delivered", "kept", target, "pcid="+message.protocolName+" exact_sha256="+message.exactHash)
	return ackBytes
}

func (kernel *Kernel) parseEnvelope(frameBytes []byte) (parsedEnvelope, error) {
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
	fields, fieldsErr := fieldsForProtocolName(envelope, protocolName, known)
	if fieldsErr != nil {
		return parsedEnvelope{}, fieldsErr
	}
	return parsedEnvelope{
		fields:       fields,
		exactHash:    protocol.HashExactBytes(frameBytes),
		protocolCID:  envelope.ProtocolCID,
		protocolName: protocolName,
	}, nil
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
	// instead of delivering them to apps or peer kernels as new promises. This is
	// a pCID-decoded transport demux check; apps still own semantic judgment.
	// Source: DI-vopab
	message, parseErr := kernel.parseEnvelope(frameBytes)
	if parseErr != nil {
		return false, parseErr
	}
	return strings.TrimSpace(message.fields["outcome"]) != "", nil
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
	ackFields := map[string]string{
		"act":     decision.ActPromise,
		"from":    "kernel:" + kernel.ContainerName,
		"to":      message.fields["from"],
		"outcome": "not_promised",
		"promise": promiseText,
		"reason":  "kernel transport non-commitment expressed as local promise content",
	}
	if promiseAbout := message.fields["promise_about"]; promiseAbout != "" {
		ackFields["promise_about"] = promiseAbout
	} else if message.protocolName == pcid.RelationshipV1 {
		ackFields["promise_about"] = "local_observation"
	}
	ack, ackErr := kernel.newAckEnvelope(message, ackFields)
	if ackErr != nil {
		kernel.record("kernel_not_promised_ack_sign_failed", "broken", message.fields["from"], ackErr.Error())
		return []byte{0x00}
	}
	ackBytes, bytesErr := ack.Bytes()
	if bytesErr != nil {
		kernel.record("kernel_not_promised_ack_bytes_failed", "broken", message.fields["from"], bytesErr.Error())
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
	// Intent: Multiple pCID receive promises can share one app/kernel session,
	// so shutdown closes each session once rather than once per protocol. Source:
	// DI-vopab
	kernel.mu.Lock()
	receiverSessions := make(map[*transport.PersistentSession]string)
	for _, registeredReceiver := range kernel.receivers {
		receiverSessions[registeredReceiver.session] = registeredReceiver.appName
	}
	kernel.receivers = make(map[receiverKey]*receiver)
	kernel.mu.Unlock()
	for registeredSession, appName := range receiverSessions {
		if registeredSession == nil {
			continue
		}
		if closeErr := registeredSession.Close(); closeErr != nil {
			kernel.record("kernel_receive_session_close_failed", "broken", appName, closeErr.Error())
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
		if closeErr := session.Close(); closeErr != nil {
			kernel.record("kernel_peer_session_close_failed", "broken", endpoint, closeErr.Error())
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
