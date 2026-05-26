package kernel

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc3/lib"
)

// Receiver records one local app's voluntary promise to receive messages for a
// pCID. It is not a permission grant or service-registry entry.
type Receiver struct {
	AppName string
	Conn    lib.FrameConn
}

// Kernel is a local promise boundary. It accepts receive promises, routes exact
// pCID-selected message bytes, records its own evidence, and refuses locally.
//
// Intent: Test a multi-app kernel without making the kernel judge app-level
// promise keep/break status. Source: DI-horak.
type Kernel struct {
	NodeName           string
	AppListen          string
	PeerListen         string
	PeerAddress        string
	ReceiveProtocolCID lib.ProtocolCID
	EvidenceLog        *lib.EvidenceLog
	receiverMutex      sync.Mutex
	receivers          map[string]Receiver
}

// Run starts app and peer listeners until the context is canceled or a listener
// fails.
func (kernel *Kernel) Run(ctx context.Context) error {
	runContext, cancel := context.WithCancel(ctx)
	defer cancel()
	if kernel.receivers == nil {
		kernel.receivers = make(map[string]Receiver)
	}
	if kernel.ReceiveProtocolCID.String() == "cidv1-raw-sha2-256:0000000000000000000000000000000000000000000000000000000000000000" {
		kernel.ReceiveProtocolCID = ReceiveProtocolCID()
	}

	errorChannel := make(chan error, 2)
	if kernel.AppListen != "" {
		go kernel.listen(runContext, "app", kernel.AppListen, kernel.handleAppConnection, errorChannel)
	}
	if kernel.PeerListen != "" {
		go kernel.listen(runContext, "peer", kernel.PeerListen, kernel.handlePeerConnection, errorChannel)
	}

	select {
	case <-runContext.Done():
		return nil
	case listenErr := <-errorChannel:
		cancel()
		return listenErr
	}
}

func (kernel *Kernel) listen(ctx context.Context, boundary string, address string, handler func(context.Context, net.Conn), errorChannel chan<- error) {
	listener, listenErr := net.Listen("tcp", address)
	if listenErr != nil {
		errorChannel <- fmt.Errorf("listen %s %s: %w", boundary, address, listenErr)
		return
	}
	defer func() {
		if closeErr := listener.Close(); closeErr != nil {
			select {
			case <-ctx.Done():
				return
			default:
				errorChannel <- closeErr
			}
		}
	}()
	go func() {
		<-ctx.Done()
		// Intent: Listener close is local shutdown evidence, not a remote command
		// channel. Source: DI-horak.
		if closeErr := listener.Close(); closeErr != nil {
			return
		}
	}()
	for {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			select {
			case <-ctx.Done():
				return
			default:
				errorChannel <- fmt.Errorf("accept %s: %w", boundary, acceptErr)
				return
			}
		}
		go handler(ctx, conn)
	}
}

func (kernel *Kernel) handleAppConnection(ctx context.Context, conn net.Conn) {
	frameConn := lib.NewFrameConn(conn)
	frameBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		kernel.record("app_read", "app/kernel", "broken", kernel.ReceiveProtocolCID, nil, readErr.Error())
		closeFrame(frameConn)
		return
	}
	envelope, parseErr := lib.ParseEnvelope(frameBytes)
	if parseErr != nil {
		kernel.record("app_receive", "app/kernel", "refused", kernel.ReceiveProtocolCID, frameBytes, parseErr.Error())
		closeFrame(frameConn)
		return
	}
	if envelope.ProtocolCID.Equal(kernel.ReceiveProtocolCID) {
		kernel.handleReceivePromise(frameConn, frameBytes, envelope)
		return
	}
	kernel.handleOutboundMessage(ctx, frameConn, frameBytes, envelope)
}

func (kernel *Kernel) handlePeerConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			fmt.Println(closeErr.Error())
		}
	}()
	frameConn := lib.NewFrameConn(conn)
	frameBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		kernel.record("peer_read", "kernel/kernel", "broken", kernel.ReceiveProtocolCID, nil, readErr.Error())
		return
	}
	envelope, parseErr := lib.ParseEnvelope(frameBytes)
	if parseErr != nil {
		kernel.record("peer_receive", "kernel/kernel", "refused", kernel.ReceiveProtocolCID, frameBytes, parseErr.Error())
		return
	}
	kernel.record("peer_receive", "kernel/kernel", "kept", envelope.ProtocolCID, frameBytes, "received message for local delivery")
	kernel.deliverToLocalReceiver(ctx, frameBytes, envelope)
}

func (kernel *Kernel) handleReceivePromise(frameConn lib.FrameConn, frameBytes []byte, envelope lib.Envelope) {
	kind, fields, kindErr := lib.EnvelopeKind(envelope)
	if kindErr != nil {
		kernel.record("app_receive_promise", "app/kernel", "refused", envelope.ProtocolCID, frameBytes, kindErr.Error())
		closeFrame(frameConn)
		return
	}
	if kind != "receive_promise_v1" || fields["pcid"] == "" || fields["app"] == "" {
		kernel.record("app_receive_promise", "app/kernel", "refused", envelope.ProtocolCID, frameBytes, "invalid receive promise")
		closeFrame(frameConn)
		return
	}
	kernel.receiverMutex.Lock()
	kernel.receivers[fields["pcid"]] = Receiver{AppName: fields["app"], Conn: frameConn}
	kernel.receiverMutex.Unlock()
	kernel.record("app_receive_promise", "app/kernel", "kept", envelope.ProtocolCID, frameBytes, fields["app"]+" promised to receive "+fields["pcid"])
}

func (kernel *Kernel) handleOutboundMessage(ctx context.Context, appConn lib.FrameConn, frameBytes []byte, envelope lib.Envelope) {
	defer closeFrame(appConn)
	fields, fieldsErr := envelope.PayloadFields()
	if fieldsErr != nil {
		kernel.record("app_receive", "app/kernel", "refused", envelope.ProtocolCID, frameBytes, fieldsErr.Error())
		return
	}
	kernel.record("app_receive", "app/kernel", "kept", envelope.ProtocolCID, frameBytes, "received outbound app message")
	if fields["to"] == kernel.NodeName || kernel.PeerAddress == "" {
		kernel.deliverToLocalReceiver(ctx, frameBytes, envelope)
		return
	}
	peerConn, dialErr := lib.DialFrameConn(kernel.PeerAddress, 10*time.Second)
	if dialErr != nil {
		kernel.record("peer_send", "kernel/kernel", "broken", envelope.ProtocolCID, frameBytes, dialErr.Error())
		return
	}
	defer closeFrame(peerConn)
	if writeErr := peerConn.WriteFrame(frameBytes); writeErr != nil {
		kernel.record("peer_send", "kernel/kernel", "broken", envelope.ProtocolCID, frameBytes, writeErr.Error())
		return
	}
	kernel.record("peer_send", "kernel/kernel", "kept", envelope.ProtocolCID, frameBytes, "forwarded to "+kernel.PeerAddress)
}

func (kernel *Kernel) deliverToLocalReceiver(ctx context.Context, frameBytes []byte, envelope lib.Envelope) {
	kernel.receiverMutex.Lock()
	receiver, ok := kernel.receivers[envelope.ProtocolCID.String()]
	kernel.receiverMutex.Unlock()
	if !ok {
		kernel.record("app_deliver", "app/kernel", "not-promised", envelope.ProtocolCID, frameBytes, "no local app promised to receive this pCID")
		return
	}
	if writeErr := receiver.Conn.WriteFrame(frameBytes); writeErr != nil {
		kernel.record("app_deliver", "app/kernel", "broken", envelope.ProtocolCID, frameBytes, writeErr.Error())
		return
	}
	select {
	case <-ctx.Done():
		kernel.record("app_deliver", "app/kernel", "kept", envelope.ProtocolCID, frameBytes, "delivered to "+receiver.AppName+" during shutdown")
	default:
		kernel.record("app_deliver", "app/kernel", "kept", envelope.ProtocolCID, frameBytes, "delivered to "+receiver.AppName)
	}
}

func (kernel *Kernel) record(event string, boundary string, outcome string, protocolCID lib.ProtocolCID, exactBytes []byte, detail string) {
	if kernel.EvidenceLog == nil {
		return
	}
	if recordErr := kernel.EvidenceLog.Record(event, boundary, outcome, protocolCID, exactBytes, detail); recordErr != nil {
		fmt.Println(recordErr.Error())
	}
}

func closeFrame(frameConn lib.FrameConn) {
	if closeErr := frameConn.Close(); closeErr != nil {
		fmt.Println(closeErr.Error())
	}
}
