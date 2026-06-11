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
	"sync"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc14-wasm/config"
	"promisegrid.dev/wire-lab/implementations/poc14-wasm/decision"
	"promisegrid.dev/wire-lab/implementations/poc14-wasm/pcid"
	"promisegrid.dev/wire-lab/implementations/poc14-wasm/protocol"
	"promisegrid.dev/wire-lab/implementations/poc14-wasm/transport"
)

const peerSendTimeout = 5 * time.Second
const listenerDrainTimeout = 750 * time.Millisecond

// Kernel runs one container-local PromiseGrid transport boundary. It accepts
// receive promises from local app processes, routes exact signed envelopes to a
// promised local receiver, and forwards cross-container bytes to peer kernels.
// Intent: Keep POC14's kernel as transport and operational evidence only; apps
// own trust, business workflow, relationship learning, and keep/break judgment.
// Source: DI-galin
type Kernel struct {
	Config        config.Config
	ContainerName string
	Protocols     pcid.Registry

	mu             sync.Mutex
	receivers      map[receiverKey]*receiver
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
	frameConn    transport.FrameConn
	mu           sync.Mutex
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
	kernel.drainHandlers(ctx)
	kernel.closeReceivers()
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
	frameBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		kernel.closeFrameConn(frameConn, "kernel_app_conn_close_failed", "")
		kernel.record("kernel_app_frame_read_failed", "broken", "", readErr.Error())
		return
	}
	message, parseErr := kernel.parseEnvelope(frameBytes)
	if parseErr != nil {
		kernel.closeFrameConn(frameConn, "kernel_app_conn_close_failed", "")
		kernel.record("kernel_app_frame_parse_failed", "broken", "", parseErr.Error())
		return
	}
	if message.protocolName == pcid.KernelReceiveV1 {
		kernel.registerReceiver(frameConn, message)
		return
	}
	defer kernel.closeFrameConn(frameConn, "kernel_app_conn_close_failed", message.fields["to"])
	ackBytes, routeErr := kernel.routeFromLocalApp(frameBytes, message)
	if routeErr != nil {
		kernel.record("kernel_route_failed", "broken", message.fields["to"], routeErr.Error())
		ackBytes = kernel.notPromisedAck(message, "I promise I could not route this exact envelope through the local kernel.")
	}
	if writeErr := frameConn.WriteFrame(ackBytes); writeErr != nil {
		kernel.record("kernel_app_ack_write_failed", "broken", message.fields["from"], writeErr.Error())
	}
}

func (kernel *Kernel) handlePeerConn(conn net.Conn) {
	frameConn := transport.NewFrameConn(conn)
	defer kernel.closeFrameConn(frameConn, "kernel_peer_conn_close_failed", "")
	frameBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		kernel.record("kernel_peer_frame_read_failed", "broken", "", readErr.Error())
		return
	}
	message, parseErr := kernel.parseEnvelope(frameBytes)
	if parseErr != nil {
		kernel.record("kernel_peer_frame_parse_failed", "broken", "", parseErr.Error())
		return
	}
	ackBytes := kernel.deliverToLocalApp(frameBytes, message)
	if writeErr := frameConn.WriteFrame(ackBytes); writeErr != nil {
		kernel.record("kernel_peer_ack_write_failed", "broken", message.fields["from"], writeErr.Error())
	}
}

func (kernel *Kernel) registerReceiver(frameConn transport.FrameConn, message parsedEnvelope) {
	appName := firstField(message.fields, "app", "field_app", "from")
	protocolName := firstField(message.fields, "pcid", "field_pcid", "protocol", "field_protocol")
	if appName == "" || protocolName == "" {
		kernel.closeFrameConn(frameConn, "kernel_receive_conn_close_failed", appName)
		kernel.record("kernel_receive_promise_malformed", "malformed", appName, "receive promise requires app and pcid fields")
		return
	}
	protocolCID, known := kernel.Protocols.CID(protocolName)
	if !known || protocolName == pcid.KernelReceiveV1 {
		kernel.closeFrameConn(frameConn, "kernel_receive_conn_close_failed", appName)
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
		frameConn:    frameConn,
	}
	kernel.mu.Unlock()
	if oldReceiver != nil {
		kernel.closeFrameConn(oldReceiver.frameConn, "kernel_receive_conn_replaced", appName)
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
	frameConn, dialErr := transport.DialFrameConn(net.JoinHostPort(hostName, strconv.Itoa(peerPort)), peerSendTimeout)
	if dialErr != nil {
		return nil, dialErr
	}
	defer kernel.closeFrameConn(frameConn, "kernel_peer_send_close_failed", target)
	if writeErr := frameConn.WriteFrame(frameBytes); writeErr != nil {
		return nil, writeErr
	}
	ackBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		return nil, readErr
	}
	kernel.record("kernel_peer_forwarded", "kept", target, "forwarded exact envelope to peer kernel")
	return ackBytes, nil
}

func (kernel *Kernel) deliverToLocalApp(frameBytes []byte, message parsedEnvelope) []byte {
	target := message.fields["to"]
	key := receiverKey{appName: target, protocolCID: message.protocolCID.String()}
	kernel.mu.Lock()
	receiver := kernel.receivers[key]
	kernel.mu.Unlock()
	if receiver == nil {
		kernel.record("kernel_unregistered_pcid", "non_commitment", target, "no local app promised pcid="+message.protocolName)
		return kernel.notPromisedAck(message, "I promise no local app has promised to receive this pCID.")
	}
	receiver.mu.Lock()
	defer receiver.mu.Unlock()
	if writeErr := receiver.frameConn.WriteFrame(frameBytes); writeErr != nil {
		kernel.record("kernel_app_deliver_failed", "broken", target, writeErr.Error())
		return kernel.notPromisedAck(message, "I promise local app delivery failed before the app could judge the message.")
	}
	ackBytes, readErr := receiver.frameConn.ReadFrame()
	if readErr != nil {
		kernel.record("kernel_app_ack_read_failed", "broken", target, readErr.Error())
		return kernel.notPromisedAck(message, "I promise local app delivery failed while waiting for app evidence.")
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
	fields, fieldsErr := envelope.PayloadFields()
	if fieldsErr != nil {
		return parsedEnvelope{}, fieldsErr
	}
	protocolName, known := kernel.Protocols.Name(envelope.ProtocolCID)
	if !known {
		protocolName = "unknown:" + envelope.ProtocolCID.String()
	}
	return parsedEnvelope{
		fields:       fields,
		exactHash:    protocol.HashExactBytes(frameBytes),
		protocolCID:  envelope.ProtocolCID,
		protocolName: protocolName,
	}, nil
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
	ack, ackErr := protocol.NewEnvelope(message.protocolCID, ackFields, "kernel:"+kernel.ContainerName)
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
	kernel.mu.Lock()
	receivers := make([]*receiver, 0, len(kernel.receivers))
	for _, registeredReceiver := range kernel.receivers {
		receivers = append(receivers, registeredReceiver)
	}
	kernel.receivers = make(map[receiverKey]*receiver)
	kernel.mu.Unlock()
	for _, registeredReceiver := range receivers {
		kernel.closeFrameConn(registeredReceiver.frameConn, "kernel_receive_conn_close_failed", registeredReceiver.appName)
	}
}

func (kernel *Kernel) closeFrameConn(frameConn transport.FrameConn, eventName, peerName string) {
	if err := frameConn.Close(); err != nil {
		kernel.record(eventName, "broken", peerName, err.Error())
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
