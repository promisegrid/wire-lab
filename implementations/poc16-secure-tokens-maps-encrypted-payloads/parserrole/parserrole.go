package parserrole

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

const parserRoleTimeout = 5 * time.Second

// ParserRole is the container-local process that owns pCID payload semantics
// between local apps and the transport kernel.
// Intent: POC16 needs a real runtime process for the parser/builder kernel role
// so the transport kernel stops projecting normal app payload fields such as
// `to`. The parser role may decode pCID-owned app payloads, but it promises only
// its own parsing, delivery, ACK, and non-commitment behavior. Source: DI-gazin
type ParserRole struct {
	Config        config.Config
	ContainerName string
	Protocols     pcid.Registry
	Parser        ProtocolParser
	Receivers     *AppReceiverRegistry

	kernelClient *KernelTransportClient
	appListener  net.Listener
	logFile      *os.File
	active       sync.WaitGroup
	mu           sync.Mutex
	appSessions  map[*transport.PersistentSession]string
}

// ProtocolParser decodes exact envelopes according to the pCID in slot 0.
// Intent: This role-local parser is where app addressing and business operation
// fields are projected; the transport kernel receives only exact-envelope
// transport requests from this process. Source: DI-gazin
type ProtocolParser struct {
	Protocols pcid.Registry
}

// ParsedMessage is the parser role's local projection of one exact envelope.
type ParsedMessage struct {
	Fields       map[string]string
	ExactHash    string
	ProtocolCID  protocol.ProtocolCID
	ProtocolName string
}

// AppReceiverRegistry remembers voluntary receive promises made by local app
// processes to this parser role.
// Intent: Local apps promise the parser role that they can receive specific
// pCIDs; the parser role then promises the transport kernel that it can receive
// exact envelopes for those pCIDs. Source: DI-gazin
type AppReceiverRegistry struct {
	mu        sync.Mutex
	receivers map[appReceiverKey]*appReceiver
}

type appReceiverKey struct {
	appName     string
	protocolCID string
}

type appReceiver struct {
	appName      string
	protocolName string
	protocolCID  protocol.ProtocolCID
	session      *transport.PersistentSession
}

// KernelTransportClient owns the parser role's persistent control stream to the
// local transport kernel.
// Intent: Parser roles submit kernel_transport_v1 promises for pCID receive
// registration and exact-envelope carriage; they never ask the kernel to inspect
// normal app payload semantics. Source: DI-gazin
type KernelTransportClient struct {
	containerName string
	parserName    string
	protocols     pcid.Registry
	session       *transport.PersistentSession
	record        func(eventName, outcome, peer, detail string)
}

// New returns an unstarted parser role for one container.
func New(cfg config.Config, containerName string) *ParserRole {
	registry := pcid.NewRegistry()
	return &ParserRole{
		Config:        cfg,
		ContainerName: containerName,
		Protocols:     registry,
		Parser:        ProtocolParser{Protocols: registry},
		Receivers:     NewAppReceiverRegistry(),
		appSessions:   make(map[*transport.PersistentSession]string),
	}
}

// NewAppReceiverRegistry returns an empty app receive-promise table.
func NewAppReceiverRegistry() *AppReceiverRegistry {
	return &AppReceiverRegistry{receivers: make(map[appReceiverKey]*appReceiver)}
}

// Run starts the app-facing listener and the parser/kernel control session until
// the container supervisor cancels the context.
func (role *ParserRole) Run(ctx context.Context) error {
	if err := role.openLog(); err != nil {
		return err
	}
	defer role.closeLog()
	if err := role.start(ctx); err != nil {
		return err
	}
	<-ctx.Done()
	role.stop()
	return nil
}

func (role *ParserRole) start(ctx context.Context) error {
	appPort, portFound := role.Config.ParserRoleAppPortForContainer(role.ContainerName)
	if !portFound {
		return fmt.Errorf("no parser role app port for container %s", role.ContainerName)
	}
	listener, listenErr := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(appPort)))
	if listenErr != nil {
		return listenErr
	}
	role.appListener = listener
	kernelAddress, addressFound := role.Config.KernelTransportAddressForContainer(role.ContainerName)
	if !addressFound {
		closeErr := listener.Close()
		if closeErr != nil {
			role.record("parser_role_listener_close_failed", "broken", "", closeErr.Error())
		}
		return fmt.Errorf("no kernel transport address for container %s", role.ContainerName)
	}
	kernelClient, clientErr := role.dialKernel(ctx, kernelAddress)
	if clientErr != nil {
		closeErr := listener.Close()
		if closeErr != nil {
			role.record("parser_role_listener_close_failed", "broken", "", closeErr.Error())
		}
		return clientErr
	}
	role.kernelClient = kernelClient
	role.record("parser_role_started", "kept", "kernel", fmt.Sprintf("app_port=%d kernel=%s", appPort, kernelAddress))
	role.active.Add(1)
	go func() {
		defer role.active.Done()
		role.acceptLoop(ctx, listener)
	}()
	return nil
}

func (role *ParserRole) dialKernel(ctx context.Context, kernelAddress string) (*KernelTransportClient, error) {
	frameConn, dialErr := transport.DialFrameConn(kernelAddress, parserRoleTimeout)
	if dialErr != nil {
		return nil, dialErr
	}
	client := &KernelTransportClient{
		containerName: role.ContainerName,
		parserName:    role.parserName(),
		protocols:     role.Protocols,
		record:        role.record,
	}
	sessionReady := make(chan struct{})
	session := transport.NewPersistentSession(
		"parser-kernel:"+role.ContainerName,
		frameConn,
		frameParentExactSHA256s,
		role.frameIsResponse,
		func(frameBytes []byte) ([]byte, error) {
			<-sessionReady
			return role.handleKernelFrame(frameBytes)
		},
		func(eventName, outcome, detail string) {
			role.record(eventName, outcome, "kernel", detail)
		},
	)
	client.session = session
	close(sessionReady)
	select {
	case <-ctx.Done():
		if closeErr := session.CloseWithReason(transport.SessionTerminalReasonProcessShutdown); closeErr != nil {
			role.record("parser_kernel_session_close_failed", "broken", "kernel", closeErr.Error())
		}
		return nil, ctx.Err()
	default:
		return client, nil
	}
}

func (role *ParserRole) acceptLoop(ctx context.Context, listener net.Listener) {
	for {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			select {
			case <-ctx.Done():
				return
			default:
				role.record("parser_role_accept_failed", "broken", "", acceptErr.Error())
				return
			}
		}
		role.active.Add(1)
		go func() {
			defer role.active.Done()
			role.handleAppConn(conn)
		}()
	}
}

func (role *ParserRole) handleAppConn(conn net.Conn) {
	frameConn := transport.NewFrameConn(conn)
	var session *transport.PersistentSession
	sessionReady := make(chan struct{})
	session = transport.NewPersistentSession(
		"app-parser:"+role.ContainerName+":"+conn.RemoteAddr().String(),
		frameConn,
		frameParentExactSHA256s,
		role.frameIsResponse,
		func(frameBytes []byte) ([]byte, error) {
			<-sessionReady
			return role.handleAppFrame(session, frameBytes)
		},
		func(eventName, outcome, detail string) {
			role.record(eventName, outcome, "app", detail)
		},
	)
	close(sessionReady)
	role.trackAppSession(session, conn.RemoteAddr().String())
	defer role.untrackAppSession(session)
	<-session.Done()
}

func (role *ParserRole) handleAppFrame(session *transport.PersistentSession, frameBytes []byte) ([]byte, error) {
	message, parseErr := role.Parser.Parse(frameBytes)
	if parseErr != nil {
		role.record("parser_role_app_frame_parse_failed", "broken", "", parseErr.Error())
		return nil, parseErr
	}
	if message.ProtocolName == pcid.KernelTransportV1 && message.Fields["transport_action"] == "app_receive_promise" {
		return nil, role.registerAppReceiver(session, message)
	}
	return role.routeOutboundAppEnvelope(frameBytes, message)
}

func (role *ParserRole) handleKernelFrame(frameBytes []byte) ([]byte, error) {
	message, parseErr := role.Parser.Parse(frameBytes)
	if parseErr != nil {
		role.record("parser_role_kernel_frame_parse_failed", "broken", "kernel", parseErr.Error())
		return nil, parseErr
	}
	role.record("parser_role_inbound_from_kernel", "kept", message.Fields["from"], "pcid="+message.ProtocolName+" exact_sha256="+message.ExactHash)
	return role.deliverToLocalApp(frameBytes, message), nil
}

func (role *ParserRole) registerAppReceiver(session *transport.PersistentSession, message ParsedMessage) error {
	appName := firstField(message.Fields, "app", "from")
	protocolName := firstField(message.Fields, "pcid", "protocol")
	if appName == "" || protocolName == "" {
		role.record("parser_role_app_receive_malformed", "malformed", appName, "receive promise requires app and pcid fields")
		return nil
	}
	protocolCID, known := role.Protocols.CID(protocolName)
	if !known || protocolName == pcid.KernelTransportV1 || protocolName == pcid.KernelReceiveV1 {
		role.record("parser_role_app_receive_rejected", "non_commitment", appName, "unknown or parser-internal pCID "+protocolName)
		return nil
	}
	role.Receivers.Register(appName, protocolName, protocolCID, session)
	role.record("parser_role_app_receive_registered", "kept", appName, "pcid="+protocolName+" promise="+message.Fields["promise"])
	return role.kernelClient.RegisterPCID(context.Background(), protocolName, protocolCID)
}

func (role *ParserRole) routeOutboundAppEnvelope(frameBytes []byte, message ParsedMessage) ([]byte, error) {
	target := strings.TrimSpace(message.Fields["to"])
	if target == "" {
		return role.notPromisedAck(message, "I promise the parser role could not route this envelope because the pCID-owned payload names no target app.")
	}
	targetContainer, containerFound := role.Config.ContainerForAgent(target)
	if !containerFound {
		return role.notPromisedAck(message, "I promise the parser role could not route this envelope because the target app is unknown.")
	}
	if targetContainer == role.ContainerName {
		return role.deliverToLocalApp(frameBytes, message), nil
	}
	role.record("parser_role_outbound_to_kernel", "kept", target, "pcid="+message.ProtocolName+" exact_sha256="+message.ExactHash)
	return role.kernelClient.CarryExactEnvelope(context.Background(), target, message.ProtocolName, frameBytes)
}

func (role *ParserRole) deliverToLocalApp(frameBytes []byte, message ParsedMessage) []byte {
	target := strings.TrimSpace(message.Fields["to"])
	if target == "" {
		ackBytes, ackErr := role.notPromisedAck(message, "I promise the parser role found no local target in the pCID-owned payload.")
		if ackErr != nil {
			role.record("parser_role_not_promised_ack_failed", "broken", message.Fields["from"], ackErr.Error())
			return []byte{0x00}
		}
		return ackBytes
	}
	receiver := role.Receivers.Lookup(target, message.ProtocolCID)
	if receiver == nil {
		role.record("parser_role_no_app_receiver", "non_commitment", target, "pcid="+message.ProtocolName)
		ackBytes, ackErr := role.notPromisedAck(message, "I promise no local app has promised this pCID to the parser role.")
		if ackErr != nil {
			role.record("parser_role_not_promised_ack_failed", "broken", message.Fields["from"], ackErr.Error())
			return []byte{0x00}
		}
		return ackBytes
	}
	ctx, cancel := context.WithTimeout(context.Background(), parserRoleTimeout)
	defer cancel()
	ackBytes, readErr := receiver.session.RoundTrip(ctx, message.ExactHash, frameBytes)
	if readErr != nil {
		role.record("parser_role_app_ack_failed", "broken", target, readErr.Error())
		ackBytes, ackErr := role.notPromisedAck(message, "I promise local parser-role delivery failed while waiting for the app ACK.")
		if ackErr != nil {
			role.record("parser_role_not_promised_ack_failed", "broken", message.Fields["from"], ackErr.Error())
			return []byte{0x00}
		}
		return ackBytes
	}
	role.record("parser_role_app_delivered", "kept", target, "pcid="+message.ProtocolName+" exact_sha256="+message.ExactHash)
	return ackBytes
}

func (client *KernelTransportClient) RegisterPCID(ctx context.Context, protocolName string, protocolCID protocol.ProtocolCID) error {
	sendCtx, cancel := context.WithTimeout(ctx, parserRoleTimeout)
	defer cancel()
	fields := map[string]string{
		"act":              decision.ActPromise,
		"from":             client.parserName,
		"to":               "kernel:" + client.containerName,
		"app":              client.parserName,
		"pcid":             protocolName,
		"pcid_cid":         protocolCID.String(),
		"promise_about":    "receive_pcid",
		"transport_action": "receive_pcid",
		"promise":          "I promise this parser role can receive exact envelopes for this pCID and forward them only to local apps that promised matching receive capability.",
		"reason":           "parser role pCID receive promise to local transport kernel",
	}
	controlBytes, buildErr := client.controlEnvelopeBytes(fields)
	if buildErr != nil {
		return buildErr
	}
	client.record("parser_role_kernel_receive_promise_sent", "kept", "kernel", "pcid="+protocolName)
	return client.session.Send(sendCtx, controlBytes)
}

func (client *KernelTransportClient) CarryExactEnvelope(ctx context.Context, target, protocolName string, envelopeBytes []byte) ([]byte, error) {
	roundTripCtx, cancel := context.WithTimeout(ctx, parserRoleTimeout)
	defer cancel()
	exactHash := protocol.HashExactBytes(envelopeBytes)
	fields := map[string]string{
		"act":                 decision.ActPromise,
		"from":                client.parserName,
		"to":                  "kernel:" + client.containerName,
		"target":              target,
		"embedded_pcid":       protocolName,
		"embedded_exact_hash": exactHash,
		"envelope_b64":        base64.StdEncoding.EncodeToString(envelopeBytes),
		"promise_about":       "carry_exact_envelope",
		"transport_action":    "carry_exact_envelope",
		"promise":             "I promise these exact envelope bytes should be carried toward the named target agent; the target was parsed by this parser role from the pCID-owned payload.",
		"reason":              "parser role exact-envelope transport request",
	}
	controlBytes, buildErr := client.controlEnvelopeBytes(fields)
	if buildErr != nil {
		return nil, buildErr
	}
	client.record("parser_role_kernel_carry_requested", "kept", target, "pcid="+protocolName+" exact_sha256="+exactHash)
	return client.session.RoundTrip(roundTripCtx, exactHash, controlBytes)
}

func (client *KernelTransportClient) controlEnvelopeBytes(fields map[string]string) ([]byte, error) {
	protocolCID := client.protocols.MustCID(pcid.KernelTransportV1)
	payloadBytes, _, payloadErr := protocol.MarshalKnownArrayPayload(pcid.KernelTransportV1, fields)
	if payloadErr != nil {
		return nil, payloadErr
	}
	envelope, envelopeErr := protocol.NewEnvelopeFromPayload(protocolCID, payloadBytes, client.parserName)
	if envelopeErr != nil {
		return nil, envelopeErr
	}
	return envelope.Bytes()
}

func (registry *AppReceiverRegistry) Register(appName, protocolName string, protocolCID protocol.ProtocolCID, session *transport.PersistentSession) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.receivers[appReceiverKey{appName: appName, protocolCID: protocolCID.String()}] = &appReceiver{
		appName:      appName,
		protocolName: protocolName,
		protocolCID:  protocolCID,
		session:      session,
	}
}

func (registry *AppReceiverRegistry) Lookup(appName string, protocolCID protocol.ProtocolCID) *appReceiver {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	return registry.receivers[appReceiverKey{appName: appName, protocolCID: protocolCID.String()}]
}

// Parse verifies the envelope and decodes slot 1 according to slot 0's pCID.
func (parser ProtocolParser) Parse(frameBytes []byte) (ParsedMessage, error) {
	envelope, parseErr := protocol.ParseEnvelope(frameBytes)
	if parseErr != nil {
		return ParsedMessage{}, parseErr
	}
	if verifyErr := protocol.VerifyEnvelope(envelope); verifyErr != nil {
		return ParsedMessage{}, verifyErr
	}
	protocolName, known := parser.Protocols.Name(envelope.ProtocolCID)
	if !known {
		protocolName = "unknown:" + envelope.ProtocolCID.String()
	}
	fields, fieldsErr := protocol.PayloadFieldsForProtocolName(protocolName, envelope.Payload)
	if fieldsErr != nil {
		return ParsedMessage{}, fieldsErr
	}
	fields["protocol"] = protocolName
	if len(envelope.ParentExactSHA256s) > 0 && fields["parent_exact_sha256"] == "" {
		fields["parent_exact_sha256"] = envelope.ParentExactSHA256s[0]
	}
	return ParsedMessage{
		Fields:       fields,
		ExactHash:    protocol.HashExactBytes(frameBytes),
		ProtocolCID:  envelope.ProtocolCID,
		ProtocolName: protocolName,
	}, nil
}

func (role *ParserRole) frameIsResponse(frameBytes []byte) (bool, error) {
	envelope, envelopeErr := protocol.ParseEnvelope(frameBytes)
	if envelopeErr != nil {
		return false, envelopeErr
	}
	if len(envelope.ParentExactSHA256s) == 0 {
		return false, nil
	}
	message, parseErr := role.Parser.Parse(frameBytes)
	if parseErr != nil {
		return false, parseErr
	}
	return strings.TrimSpace(message.Fields["outcome"]) != "", nil
}

func frameParentExactSHA256s(frameBytes []byte) ([]string, error) {
	envelope, parseErr := protocol.ParseEnvelope(frameBytes)
	if parseErr != nil {
		return nil, parseErr
	}
	return append([]string(nil), envelope.ParentExactSHA256s...), nil
}

func (role *ParserRole) notPromisedAck(message ParsedMessage, promiseText string) ([]byte, error) {
	ackFields := map[string]string{
		"act":     decision.ActPromise,
		"from":    role.parserName(),
		"to":      message.Fields["from"],
		"outcome": "not_promised",
		"promise": promiseText,
		"reason":  "parser-role local non-commitment",
	}
	if promiseAbout := message.Fields["promise_about"]; promiseAbout != "" {
		ackFields["promise_about"] = promiseAbout
	} else {
		ackFields["promise_about"] = "local_observation"
	}
	payloadBytes, arrayPayload, payloadErr := protocol.MarshalKnownArrayPayload(message.ProtocolName, ackFields)
	if payloadErr == nil && arrayPayload {
		ack, ackErr := protocol.NewEnvelopeFromPayloadWithParents(message.ProtocolCID, payloadBytes, []string{message.ExactHash}, role.parserName())
		if ackErr != nil {
			return nil, ackErr
		}
		return ack.Bytes()
	}
	ack, ackErr := protocol.NewEnvelopeWithParents(message.ProtocolCID, ackFields, []string{message.ExactHash}, role.parserName())
	if ackErr != nil {
		return nil, ackErr
	}
	return ack.Bytes()
}

func (role *ParserRole) stop() {
	if role.appListener != nil {
		if closeErr := role.appListener.Close(); closeErr != nil {
			role.record("parser_role_listener_close_failed", "broken", "", closeErr.Error())
		}
	}
	if role.kernelClient != nil && role.kernelClient.session != nil {
		if closeErr := role.kernelClient.session.CloseWithReason(transport.SessionTerminalReasonProcessShutdown); closeErr != nil {
			role.record("parser_kernel_session_close_failed", "broken", "kernel", closeErr.Error())
		}
	}
	role.closeAppSessions()
	role.active.Wait()
}

func (role *ParserRole) trackAppSession(session *transport.PersistentSession, remoteAddress string) {
	// Intent: Parser roles own local app streams for lifecycle accounting even
	// when app receive promises are routing state. Shutdown closes these streams
	// explicitly so the analyzer can distinguish clean process exit from leaked
	// sessions. Source: DI-gazin
	role.mu.Lock()
	role.appSessions[session] = remoteAddress
	role.mu.Unlock()
}

func (role *ParserRole) untrackAppSession(session *transport.PersistentSession) {
	role.mu.Lock()
	delete(role.appSessions, session)
	role.mu.Unlock()
}

func (role *ParserRole) closeAppSessions() {
	// Intent: App/parser sessions are local process pipes for exact envelopes;
	// closing them at parser shutdown keeps POC16's parser role a real process
	// with complete terminal accounting. Source: DI-gazin
	role.mu.Lock()
	appSessions := make(map[*transport.PersistentSession]string, len(role.appSessions))
	for session, remoteAddress := range role.appSessions {
		appSessions[session] = remoteAddress
	}
	role.appSessions = make(map[*transport.PersistentSession]string)
	role.mu.Unlock()
	for session, remoteAddress := range appSessions {
		if session == nil {
			continue
		}
		if closeErr := session.CloseWithReason(transport.SessionTerminalReasonProcessShutdown); closeErr != nil {
			role.record("parser_app_session_close_failed", "broken", remoteAddress, closeErr.Error())
		}
	}
}

func (role *ParserRole) parserName() string {
	return "parser:" + role.ContainerName
}

func firstField(fields map[string]string, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(fields[key]); value != "" {
			return value
		}
	}
	return ""
}

func (role *ParserRole) record(eventName, outcome, peer, detail string) {
	event := decision.Event{
		Observer: role.parserName(),
		Event:    eventName,
		Outcome:  outcome,
		Peer:     peer,
		Detail:   detail,
	}
	encoded, err := json.Marshal(event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal parser role event: %v\n", err)
		return
	}
	fmt.Println(string(encoded))
	if role.logFile != nil {
		if _, writeErr := role.logFile.Write(append(encoded, '\n')); writeErr != nil {
			fmt.Fprintf(os.Stderr, "write parser role event: %v\n", writeErr)
		}
	}
}

func (role *ParserRole) openLog() error {
	runDir := filepath.Join(role.Config.RunRoot, role.Config.RunID)
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		return err
	}
	logPath := filepath.Join(runDir, "parser-"+role.ContainerName+".jsonl")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	role.logFile = logFile
	return nil
}

func (role *ParserRole) closeLog() {
	if role.logFile != nil {
		if err := role.logFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "close parser role log: %v\n", err)
		}
	}
}
