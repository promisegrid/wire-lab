package lifecycle

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/decision"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/eventstream"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/transport"
)

const (
	EnvAddress = "POC16_LIFECYCLE_ADDR"
	EnvProfile = "POC16_LIFECYCLE_PROFILE"

	RoleKindApp    = "app"
	RoleKindParser = "parser"
	RoleKindKernel = "kernel"
)

const lifecycleWaitPoll = 25 * time.Millisecond

// RoleOptions describe one local role's lifecycle promise surface.
//
// Intent: Lifecycle options are local process facts. They do not make the
// supervisor a PromiseGrid authority; they only let a role voluntarily issue a
// signed shutdown token to the local process supervisor. Source: DI-jafoj
type RoleOptions struct {
	Address          string
	SupervisorRoleID string
	RoleID           string
	RoleKind         string
	ChannelProfile   string
	RunID            string
	ProtocolCID      protocol.ProtocolCID
	ShutdownGrace    time.Duration
	ReadyAfterWork   bool
	Stdin            io.Reader
}

// ManagedRole owns the child side of local_lifecycle_v1 for one process.
type ManagedRole struct {
	options RoleOptions
	token   protocol.LifecycleToken
	conn    transport.FrameConn

	workCancel  context.CancelFunc
	invoked     chan protocol.LifecyclePayload
	started     bool
	mu          sync.Mutex
	invocations int64
}

// Supervisor owns local lifecycle tokens for one container supervisor.
type Supervisor struct {
	roleID      string
	runID       string
	protocolCID protocol.ProtocolCID
	grace       time.Duration
	listener    net.Listener
	record      func(event decision.Event)

	mu        sync.Mutex
	roles     map[string]*supervisedRole
	done      chan struct{}
	readyCh   chan struct{}
	fulfillCh chan struct{}
}

type supervisedRole struct {
	roleID         string
	roleKind       string
	channelProfile string
	token          protocol.LifecycleToken
	conn           transport.FrameConn
	stdin          io.WriteCloser
	ready          bool
	fulfilled      bool
}

// OptionsFromEnv creates lifecycle options from supervisor-provided process
// environment variables.
func OptionsFromEnv(supervisorRoleID, roleID, roleKind, runID string, protocolCID protocol.ProtocolCID, shutdownGrace time.Duration, readyAfterWork bool, stdin io.Reader) RoleOptions {
	return RoleOptions{
		Address:          os.Getenv(EnvAddress),
		SupervisorRoleID: supervisorRoleID,
		RoleID:           roleID,
		RoleKind:         roleKind,
		ChannelProfile:   os.Getenv(EnvProfile),
		RunID:            runID,
		ProtocolCID:      protocolCID,
		ShutdownGrace:    shutdownGrace,
		ReadyAfterWork:   readyAfterWork,
		Stdin:            stdin,
	}
}

// NewManagedRole prepares the role endpoint but does not start I/O.
func NewManagedRole(options RoleOptions) (*ManagedRole, error) {
	if strings.TrimSpace(options.Address) == "" {
		return nil, fmt.Errorf("lifecycle address is required")
	}
	if strings.TrimSpace(options.ChannelProfile) == "" {
		return nil, fmt.Errorf("lifecycle channel profile is required")
	}
	now := time.Now()
	token, tokenErr := protocol.IssueLifecycleToken(protocol.LifecycleTokenTerms{
		IssuerRoleID:   options.RoleID,
		AudienceRoleID: options.SupervisorRoleID,
		RunID:          options.RunID,
		RoleKind:       options.RoleKind,
		ChannelProfile: options.ChannelProfile,
		ProtocolCID:    options.ProtocolCID,
		GraceMillis:    options.ShutdownGrace.Milliseconds(),
		MaxInvocations: 1,
		IssuedAtUnix:   now.Unix(),
		NotBeforeUnix:  now.Add(-1 * time.Second).Unix(),
		ExpiresUnix:    now.Add(options.ShutdownGrace + time.Hour).Unix(),
		TokenID:        protocol.CIDForExactBytes([]byte(options.RunID + ":" + options.RoleID + ":" + options.ChannelProfile)),
		ShutdownTerms: []string{
			"stop starting new work",
			"drain accepted work",
			"close persistent sessions",
			"flush local events",
			"emit a role summary",
			"exit voluntarily",
		},
	})
	if tokenErr != nil {
		return nil, tokenErr
	}
	return &ManagedRole{
		options: options,
		token:   token,
		invoked: make(chan protocol.LifecyclePayload, 1),
	}, nil
}

// RunManaged runs work directly when no lifecycle supervisor is configured and
// otherwise wraps it in local_lifecycle_v1 token issue/invocation semantics.
func RunManaged(ctx context.Context, options RoleOptions, work func(context.Context) error) error {
	if strings.TrimSpace(options.Address) == "" {
		return work(ctx)
	}
	role, roleErr := NewManagedRole(options)
	if roleErr != nil {
		return roleErr
	}
	return role.Run(ctx, work)
}

// Run starts lifecycle I/O, runs the role, waits for token invocation, and then
// records fulfillment.
func (role *ManagedRole) Run(ctx context.Context, work func(context.Context) error) error {
	if strings.TrimSpace(role.options.Address) == "" {
		return work(ctx)
	}
	if startErr := role.start(ctx); startErr != nil {
		return startErr
	}
	workCtx, cancel := context.WithCancel(ctx)
	role.workCancel = cancel
	defer cancel()
	if !role.options.ReadyAfterWork {
		if readyErr := role.sendReady("server_started"); readyErr != nil {
			return readyErr
		}
	}
	workDone := make(chan error, 1)
	go func() {
		workDone <- work(workCtx)
	}()
	var workErr error
	if role.options.ReadyAfterWork {
		select {
		case workErr = <-workDone:
			if readyErr := role.sendReady("work_complete"); readyErr != nil && workErr == nil {
				workErr = readyErr
			}
		case <-ctx.Done():
			cancel()
			workErr = <-workDone
			if workErr == nil {
				workErr = ctx.Err()
			}
		}
	}
	invocation, invokeErr := role.waitInvocation(ctx)
	cancel()
	if !role.options.ReadyAfterWork {
		workErr = <-workDone
	}
	outcome := protocol.LifecycleOutcomeKept
	detail := "role honored lifecycle token"
	if invokeErr != nil {
		outcome = protocol.LifecycleOutcomeTimedOut
		detail = invokeErr.Error()
	} else if workErr != nil && !isContextShutdown(workErr) {
		outcome = protocol.LifecycleOutcomeBroken
		detail = workErr.Error()
	}
	if summaryErr := role.sendRoleSummary(outcome, detail); summaryErr != nil && workErr == nil {
		workErr = summaryErr
	}
	if fulfillErr := role.sendFulfillment(outcome, detail); fulfillErr != nil && workErr == nil {
		workErr = fulfillErr
	}
	if invocation.Kind != "" {
		role.record("local_lifecycle_invocation_observed", outcome, invocation.Promiser, "token_cid="+invocation.TokenCID+" reason="+invocation.Reason)
	}
	if isContextShutdown(workErr) {
		return nil
	}
	return workErr
}

// HandleInvocationFrame lets parser-path lifecycle messages trigger the same
// token verification path as TCP or stdio lifecycle channels.
func (role *ManagedRole) HandleInvocationFrame(frameBytes []byte) ([]byte, error) {
	if handleErr := role.handleInvocationFrame(frameBytes); handleErr != nil {
		return nil, handleErr
	}
	return nil, nil
}

func (role *ManagedRole) start(ctx context.Context) error {
	frameConn, dialErr := transport.DialFrameConn(role.options.Address, role.options.ShutdownGrace)
	if dialErr != nil {
		return dialErr
	}
	role.conn = frameConn
	role.started = true
	tokenFrame, tokenFrameErr := protocol.EncodeLifecycleMessage(role.options.ProtocolCID, protocol.NewLifecycleTokenIssuedPayload(role.token))
	if tokenFrameErr != nil {
		return tokenFrameErr
	}
	if writeErr := frameConn.WriteFrame(tokenFrame); writeErr != nil {
		return writeErr
	}
	role.record("local_lifecycle_token_issued", "kept", role.options.SupervisorRoleID, "token_cid="+role.token.CID+" channel_profile="+role.options.ChannelProfile)
	go role.readTCPInvocations()
	if role.options.ChannelProfile == protocol.LifecycleChannelStdio && role.options.Stdin != nil {
		go role.readStdinInvocations(role.options.Stdin)
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (role *ManagedRole) readTCPInvocations() {
	for {
		frameBytes, readErr := role.conn.ReadFrame()
		if readErr != nil {
			return
		}
		if handleErr := role.handleInvocationFrame(frameBytes); handleErr != nil {
			role.record("local_lifecycle_invocation_rejected", protocol.LifecycleOutcomeRejected, role.options.SupervisorRoleID, handleErr.Error())
		}
	}
}

func (role *ManagedRole) readStdinInvocations(reader io.Reader) {
	frameBytes, readErr := readFrame(reader)
	if readErr != nil {
		role.record("local_lifecycle_stdio_read_failed", "broken", role.options.SupervisorRoleID, readErr.Error())
		return
	}
	if handleErr := role.handleInvocationFrame(frameBytes); handleErr != nil {
		role.record("local_lifecycle_invocation_rejected", protocol.LifecycleOutcomeRejected, role.options.SupervisorRoleID, handleErr.Error())
	}
}

func (role *ManagedRole) handleInvocationFrame(frameBytes []byte) error {
	payload, decodeErr := protocol.DecodeLifecycleMessage(frameBytes, role.options.ProtocolCID)
	if decodeErr != nil {
		return decodeErr
	}
	if payload.Kind != protocol.LifecycleKindTokenInvoked {
		return fmt.Errorf("lifecycle invocation got kind %q", payload.Kind)
	}
	tokenBytes, tokenErr := payload.TokenBytes()
	if tokenErr != nil {
		return tokenErr
	}
	role.mu.Lock()
	priorInvocations := role.invocations
	role.invocations++
	role.mu.Unlock()
	_, verifyErr := protocol.VerifyLifecycleToken(tokenBytes, role.token.Terms, time.Now(), priorInvocations)
	if verifyErr != nil {
		return verifyErr
	}
	role.record("local_lifecycle_token_verified", "kept", payload.Promiser, "token_cid="+payload.TokenCID+" channel_profile="+payload.ChannelProfile)
	select {
	case role.invoked <- payload:
	default:
		return fmt.Errorf("lifecycle token already invoked")
	}
	if role.workCancel != nil {
		role.workCancel()
	}
	return nil
}

func (role *ManagedRole) sendReady(detail string) error {
	frameBytes, frameErr := protocol.EncodeLifecycleMessage(role.options.ProtocolCID, protocol.NewLifecycleReadyPayload(role.token, detail))
	if frameErr != nil {
		return frameErr
	}
	if writeErr := role.conn.WriteFrame(frameBytes); writeErr != nil {
		return writeErr
	}
	role.record("local_lifecycle_ready", "kept", role.options.SupervisorRoleID, "token_cid="+role.token.CID+" detail="+detail)
	return nil
}

func (role *ManagedRole) waitInvocation(ctx context.Context) (protocol.LifecyclePayload, error) {
	select {
	case invocation := <-role.invoked:
		return invocation, nil
	case <-ctx.Done():
		return protocol.LifecyclePayload{}, ctx.Err()
	}
}

func (role *ManagedRole) sendRoleSummary(outcome, detail string) error {
	role.record("local_lifecycle_role_summary", outcome, role.options.SupervisorRoleID, "role_id="+role.options.RoleID+" role_kind="+role.options.RoleKind+" channel_profile="+role.options.ChannelProfile+" detail="+detail)
	return nil
}

func (role *ManagedRole) sendFulfillment(outcome, detail string) error {
	frameBytes, frameErr := protocol.EncodeLifecycleMessage(role.options.ProtocolCID, protocol.NewLifecycleFulfilledPayload(role.token, outcome, detail))
	if frameErr != nil {
		return frameErr
	}
	if role.started {
		if writeErr := role.conn.WriteFrame(frameBytes); writeErr != nil {
			return writeErr
		}
	}
	role.record("local_lifecycle_token_fulfilled", outcome, role.options.SupervisorRoleID, "token_cid="+role.token.CID+" detail="+detail)
	return nil
}

func (role *ManagedRole) record(eventName, outcome, peer, detail string) {
	recordEvent(role.options.RoleID, eventName, outcome, peer, detail)
}

// NewSupervisor starts the supervisor-side lifecycle listener.
func NewSupervisor(roleID, runID string, protocolCID protocol.ProtocolCID, grace time.Duration, record func(event decision.Event)) (*Supervisor, error) {
	listener, listenErr := net.Listen("tcp", "127.0.0.1:0")
	if listenErr != nil {
		return nil, listenErr
	}
	supervisor := &Supervisor{
		roleID:      roleID,
		runID:       runID,
		protocolCID: protocolCID,
		grace:       grace,
		listener:    listener,
		record:      record,
		roles:       make(map[string]*supervisedRole),
		done:        make(chan struct{}),
		readyCh:     make(chan struct{}, 1),
		fulfillCh:   make(chan struct{}, 1),
	}
	go supervisor.acceptLoop()
	return supervisor, nil
}

// Address returns the loopback address child roles use for lifecycle streams.
func (supervisor *Supervisor) Address() string {
	return supervisor.listener.Addr().String()
}

// EnvFor returns the child process environment entries for one role.
func (supervisor *Supervisor) EnvFor(channelProfile string) []string {
	return []string{
		EnvAddress + "=" + supervisor.Address(),
		EnvProfile + "=" + channelProfile,
	}
}

// RegisterStdin stores the stdin writer for a stdio-profile role.
func (supervisor *Supervisor) RegisterStdin(roleID string, writer io.WriteCloser) {
	supervisor.mu.Lock()
	defer supervisor.mu.Unlock()
	role := supervisor.ensureRoleLocked(roleID)
	role.stdin = writer
}

// WaitReady waits until every listed role has issued a token and reported
// readiness.
func (supervisor *Supervisor) WaitReady(ctx context.Context, roleIDs []string) error {
	return supervisor.waitFor(ctx, roleIDs, func(role *supervisedRole) bool {
		return role.ready
	}, supervisor.readyCh, "ready")
}

// Invoke presents the stored token back to its issuer over the selected channel
// profile.
func (supervisor *Supervisor) Invoke(ctx context.Context, roleID, reason string, parserPathAddress string) error {
	supervisor.mu.Lock()
	role := supervisor.roles[roleID]
	supervisor.mu.Unlock()
	if role == nil {
		return fmt.Errorf("no lifecycle token for %s", roleID)
	}
	deadline := time.Now().Add(supervisor.grace)
	payload := protocol.NewLifecycleInvocationPayload(supervisor.roleID, role.token, reason, deadline)
	frameBytes, frameErr := protocol.EncodeLifecycleMessage(supervisor.protocolCID, payload)
	if frameErr != nil {
		return frameErr
	}
	supervisor.recordEvent(supervisor.roleID, "local_lifecycle_token_invoked", "kept", roleID, "token_cid="+role.token.CID+" channel_profile="+role.channelProfile+" reason="+reason)
	switch role.channelProfile {
	case protocol.LifecycleChannelTCP:
		return role.conn.WriteFrame(frameBytes)
	case protocol.LifecycleChannelParserPath:
		return supervisor.invokeParserPath(ctx, parserPathAddress, frameBytes)
	case protocol.LifecycleChannelStdio:
		if role.stdin == nil {
			return fmt.Errorf("stdio lifecycle role %s has no stdin writer", roleID)
		}
		return writeFrame(role.stdin, frameBytes)
	default:
		return fmt.Errorf("unsupported lifecycle channel profile %q", role.channelProfile)
	}
}

// WaitFulfilled waits until every listed role reports token fulfillment.
func (supervisor *Supervisor) WaitFulfilled(ctx context.Context, roleIDs []string) error {
	return supervisor.waitFor(ctx, roleIDs, func(role *supervisedRole) bool {
		return role.fulfilled
	}, supervisor.fulfillCh, "fulfilled")
}

// Close closes the supervisor listener.
func (supervisor *Supervisor) Close() error {
	close(supervisor.done)
	return supervisor.listener.Close()
}

func (supervisor *Supervisor) acceptLoop() {
	for {
		conn, acceptErr := supervisor.listener.Accept()
		if acceptErr != nil {
			select {
			case <-supervisor.done:
				return
			default:
				supervisor.recordEvent(supervisor.roleID, "local_lifecycle_accept_failed", "broken", "", acceptErr.Error())
				return
			}
		}
		go supervisor.handleConn(conn)
	}
}

func (supervisor *Supervisor) handleConn(conn net.Conn) {
	frameConn := transport.NewFrameConn(conn)
	frameBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		supervisor.recordEvent(supervisor.roleID, "local_lifecycle_token_read_failed", "broken", "", readErr.Error())
		return
	}
	payload, decodeErr := protocol.DecodeLifecycleMessage(frameBytes, supervisor.protocolCID)
	if decodeErr != nil {
		supervisor.recordEvent(supervisor.roleID, "local_lifecycle_token_decode_failed", "broken", "", decodeErr.Error())
		return
	}
	if payload.Kind != protocol.LifecycleKindTokenIssued {
		supervisor.recordEvent(supervisor.roleID, "local_lifecycle_unexpected_first_message", "broken", payload.Promiser, "kind="+payload.Kind)
		return
	}
	tokenBytes, tokenErr := payload.TokenBytes()
	if tokenErr != nil {
		supervisor.recordEvent(supervisor.roleID, "local_lifecycle_token_decode_failed", "broken", payload.Promiser, tokenErr.Error())
		return
	}
	expected := protocol.LifecycleTokenTerms{
		IssuerRoleID:   payload.RoleID,
		AudienceRoleID: supervisor.roleID,
		RunID:          supervisor.runID,
		RoleKind:       payload.RoleKind,
		ChannelProfile: payload.ChannelProfile,
		ProtocolCID:    supervisor.protocolCID,
	}
	terms, verifyErr := protocol.VerifyLifecycleToken(tokenBytes, expected, time.Now(), 0)
	if verifyErr != nil {
		supervisor.recordEvent(supervisor.roleID, "local_lifecycle_token_rejected", protocol.LifecycleOutcomeRejected, payload.RoleID, verifyErr.Error())
		return
	}
	token := protocol.LifecycleToken{
		Terms:        terms,
		COSEBytes:    tokenBytes,
		CID:          payload.TokenCID,
		COSEBase64:   payload.TokenCOSEBase64,
		PublicBase64: payload.PublicKeyBase64,
	}
	supervisor.mu.Lock()
	role := supervisor.ensureRoleLocked(payload.RoleID)
	role.roleID = payload.RoleID
	role.roleKind = payload.RoleKind
	role.channelProfile = payload.ChannelProfile
	role.token = token
	role.conn = frameConn
	supervisor.mu.Unlock()
	supervisor.recordEvent(supervisor.roleID, "local_lifecycle_token_issued", "kept", payload.RoleID, "token_cid="+payload.TokenCID+" channel_profile="+payload.ChannelProfile)
	supervisor.readRoleMessages(frameConn)
}

func (supervisor *Supervisor) readRoleMessages(frameConn transport.FrameConn) {
	for {
		frameBytes, readErr := frameConn.ReadFrame()
		if readErr != nil {
			return
		}
		payload, decodeErr := protocol.DecodeLifecycleMessage(frameBytes, supervisor.protocolCID)
		if decodeErr != nil {
			supervisor.recordEvent(supervisor.roleID, "local_lifecycle_message_decode_failed", "broken", "", decodeErr.Error())
			continue
		}
		switch payload.Kind {
		case protocol.LifecycleKindReady:
			supervisor.markReady(payload)
		case protocol.LifecycleKindTokenFulfilled:
			supervisor.markFulfilled(payload)
		default:
			supervisor.recordEvent(supervisor.roleID, "local_lifecycle_unexpected_message", "non_commitment", payload.RoleID, "kind="+payload.Kind)
		}
	}
}

func (supervisor *Supervisor) markReady(payload protocol.LifecyclePayload) {
	supervisor.mu.Lock()
	role := supervisor.ensureRoleLocked(payload.RoleID)
	role.ready = true
	supervisor.mu.Unlock()
	supervisor.recordEvent(supervisor.roleID, "local_lifecycle_ready", "kept", payload.RoleID, "token_cid="+payload.TokenCID+" detail="+payload.Detail)
	notify(supervisor.readyCh)
}

func (supervisor *Supervisor) markFulfilled(payload protocol.LifecyclePayload) {
	supervisor.mu.Lock()
	role := supervisor.ensureRoleLocked(payload.RoleID)
	role.fulfilled = payload.Outcome == protocol.LifecycleOutcomeKept
	supervisor.mu.Unlock()
	supervisor.recordEvent(supervisor.roleID, "local_lifecycle_token_fulfilled", payload.Outcome, payload.RoleID, "token_cid="+payload.TokenCID+" detail="+payload.Detail)
	notify(supervisor.fulfillCh)
}

func (supervisor *Supervisor) waitFor(ctx context.Context, roleIDs []string, predicate func(*supervisedRole) bool, notifyCh <-chan struct{}, label string) error {
	ticker := time.NewTicker(lifecycleWaitPoll)
	defer ticker.Stop()
	for {
		if supervisor.allRoles(roleIDs, predicate) {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("lifecycle wait for %s: %w", label, ctx.Err())
		case <-notifyCh:
		case <-ticker.C:
		}
	}
}

func (supervisor *Supervisor) allRoles(roleIDs []string, predicate func(*supervisedRole) bool) bool {
	supervisor.mu.Lock()
	defer supervisor.mu.Unlock()
	for _, roleID := range roleIDs {
		role := supervisor.roles[roleID]
		if role == nil || !predicate(role) {
			return false
		}
	}
	return true
}

func (supervisor *Supervisor) ensureRoleLocked(roleID string) *supervisedRole {
	role := supervisor.roles[roleID]
	if role == nil {
		role = &supervisedRole{roleID: roleID}
		supervisor.roles[roleID] = role
	}
	return role
}

func (supervisor *Supervisor) invokeParserPath(ctx context.Context, address string, frameBytes []byte) error {
	if strings.TrimSpace(address) == "" {
		return fmt.Errorf("parser-path lifecycle invocation needs parser address")
	}
	frameConn, dialErr := transport.DialFrameConn(address, supervisor.grace)
	if dialErr != nil {
		return dialErr
	}
	defer func() {
		if closeErr := frameConn.Close(); closeErr != nil {
			supervisor.recordEvent(supervisor.roleID, "local_lifecycle_parser_path_close_failed", "broken", address, closeErr.Error())
		}
	}()
	done := make(chan error, 1)
	go func() {
		done <- frameConn.WriteFrame(frameBytes)
	}()
	select {
	case writeErr := <-done:
		return writeErr
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (supervisor *Supervisor) recordEvent(observer, eventName, outcome, peer, detail string) {
	if supervisor.record == nil {
		return
	}
	supervisor.record(decision.Event{
		Observer: observer,
		Event:    eventName,
		Outcome:  outcome,
		Peer:     peer,
		Detail:   detail,
	})
}

func recordEvent(observer, eventName, outcome, peer, detail string) {
	event := decision.Event{
		Observer: observer,
		Event:    eventName,
		Outcome:  outcome,
		Peer:     peer,
		Detail:   detail,
	}
	if _, err := eventstream.WriteStdoutJSON(event); err != nil {
		fmt.Fprintf(os.Stderr, "lifecycle event stdout write failed: %v\n", err)
	}
}

func notify(channel chan<- struct{}) {
	select {
	case channel <- struct{}{}:
	default:
	}
}

func writeFrame(writer io.Writer, frameBytes []byte) error {
	if len(frameBytes) == 0 || len(frameBytes) > 16*1024*1024 {
		return fmt.Errorf("invalid lifecycle frame length: %d", len(frameBytes))
	}
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(frameBytes)))
	if err := writeAll(writer, header); err != nil {
		return err
	}
	return writeAll(writer, frameBytes)
}

func readFrame(reader io.Reader) ([]byte, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(header)
	if length == 0 || length > 16*1024*1024 {
		return nil, fmt.Errorf("invalid lifecycle frame length: %d", length)
	}
	frameBytes := make([]byte, int(length))
	if _, err := io.ReadFull(reader, frameBytes); err != nil {
		return nil, err
	}
	return frameBytes, nil
}

func writeAll(writer io.Writer, frameBytes []byte) error {
	for len(frameBytes) > 0 {
		written, writeErr := writer.Write(frameBytes)
		if written > 0 {
			frameBytes = frameBytes[written:]
		}
		if writeErr != nil {
			return writeErr
		}
		if written == 0 {
			return io.ErrShortWrite
		}
	}
	return nil
}

func isContextShutdown(err error) bool {
	return err == nil || err == context.Canceled || err == context.DeadlineExceeded
}
