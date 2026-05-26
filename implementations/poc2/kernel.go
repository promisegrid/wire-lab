package main

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// Kernel is a local role, not a ruler. It promises to receive pCID-selected
// messages, preserve evidence, forward when configured, deliver to a local app
// that has promised to receive, and refuse unsupported pCIDs locally.
//
// Intent: Test the kernel-as-promise-boundary model with executable code rather
// than final API claims. Source: DI-ratij; DI-tijat.
type Kernel struct {
	NodeName       string
	AppListen      string
	PeerListen     string
	PeerAddress    string
	ProtocolCID    ProtocolCID
	EvidenceLog    *EvidenceLog
	receiverMutex  sync.Mutex
	receiver       *FrameConn
	receiverName   string
	listenerCancel context.CancelFunc
}

func (kernel *Kernel) Run(ctx context.Context) error {
	runContext, cancel := context.WithCancel(ctx)
	kernel.listenerCancel = cancel
	defer cancel()

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
		// Intent: Closing the listener is how the test and container shutdown path
		// stop accepting promises without inventing a kernel command channel.
		// Source: DI-tijat.
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
	frameConn := NewFrameConn(conn)
	frameBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		if recordErr := kernel.EvidenceLog.Record("app_read", "app/kernel", "broken", kernel.ProtocolCID, nil, readErr.Error()); recordErr != nil {
			fmt.Println(recordErr.Error())
		}
		if closeErr := frameConn.Close(); closeErr != nil {
			fmt.Println(closeErr.Error())
		}
		return
	}
	envelope, parseErr := ParseEnvelope(frameBytes)
	if parseErr != nil || !envelope.ProtocolCID.Equal(kernel.ProtocolCID) {
		kernel.refuseUnsupported(frameConn, frameBytes, parseErr)
		return
	}
	kind, fields, kindErr := envelopeKind(envelope)
	if kindErr != nil {
		kernel.refuseWithObservation(frameConn, frameBytes, "refused", kindErr.Error())
		return
	}
	if recordErr := kernel.EvidenceLog.Record("app_receive", "app/kernel", "kept", kernel.ProtocolCID, frameBytes, "received "+kind); recordErr != nil {
		fmt.Println(recordErr.Error())
	}
	switch kind {
	case "receive_promise_v1":
		kernel.storeReceiver(fields["from"], frameConn)
	case "hello_v1":
		kernel.handleHelloFromLocalApp(ctx, frameConn, frameBytes, fields)
	default:
		kernel.refuseWithObservation(frameConn, frameBytes, "not-promised", "unsupported payload kind "+kind)
	}
}

func (kernel *Kernel) handlePeerConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			fmt.Println(closeErr.Error())
		}
	}()
	frameConn := NewFrameConn(conn)
	frameBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		if recordErr := kernel.EvidenceLog.Record("peer_read", "kernel/kernel", "broken", kernel.ProtocolCID, nil, readErr.Error()); recordErr != nil {
			fmt.Println(recordErr.Error())
		}
		return
	}
	envelope, parseErr := ParseEnvelope(frameBytes)
	if parseErr != nil || !envelope.ProtocolCID.Equal(kernel.ProtocolCID) {
		kernel.writePeerObservation(frameConn, frameBytes, "refused", "unsupported or unreadable pCID")
		return
	}
	kind, fields, kindErr := envelopeKind(envelope)
	if kindErr != nil {
		kernel.writePeerObservation(frameConn, frameBytes, "refused", kindErr.Error())
		return
	}
	if recordErr := kernel.EvidenceLog.Record("peer_receive", "kernel/kernel", "kept", kernel.ProtocolCID, frameBytes, "received "+kind); recordErr != nil {
		fmt.Println(recordErr.Error())
	}
	if kind != "hello_v1" {
		kernel.writePeerObservation(frameConn, frameBytes, "not-promised", "peer payload kind not handled: "+kind)
		return
	}
	outcome, detail := kernel.deliverToLocalReceiver(ctx, frameBytes, fields)
	kernel.writePeerObservation(frameConn, frameBytes, outcome, detail)
}

func (kernel *Kernel) storeReceiver(receiverName string, frameConn FrameConn) {
	kernel.receiverMutex.Lock()
	kernel.receiver = &frameConn
	kernel.receiverName = receiverName
	kernel.receiverMutex.Unlock()
	if recordErr := kernel.EvidenceLog.Record("app_receive_promise", "app/kernel", "kept", kernel.ProtocolCID, nil, receiverName+" promised to receive hello messages"); recordErr != nil {
		fmt.Println(recordErr.Error())
	}
}

func (kernel *Kernel) handleHelloFromLocalApp(ctx context.Context, appConn FrameConn, frameBytes []byte, fields map[string]string) {
	destination := fields["to"]
	if destination == kernel.NodeName || kernel.PeerAddress == "" {
		outcome, detail := kernel.deliverToLocalReceiver(ctx, frameBytes, fields)
		kernel.refuseWithObservation(appConn, frameBytes, outcome, detail)
		return
	}
	peerConn, dialErr := dialFrameConn(kernel.PeerAddress, 10*time.Second)
	if dialErr != nil {
		kernel.refuseWithObservation(appConn, frameBytes, "refused", dialErr.Error())
		return
	}
	defer func() {
		if closeErr := peerConn.Close(); closeErr != nil {
			fmt.Println(closeErr.Error())
		}
	}()
	if writeErr := peerConn.WriteFrame(frameBytes); writeErr != nil {
		kernel.refuseWithObservation(appConn, frameBytes, "broken", writeErr.Error())
		return
	}
	if recordErr := kernel.EvidenceLog.Record("peer_send", "kernel/kernel", "kept", kernel.ProtocolCID, frameBytes, "forwarded hello to "+kernel.PeerAddress); recordErr != nil {
		fmt.Println(recordErr.Error())
	}
	ackBytes, readErr := peerConn.ReadFrame()
	if readErr != nil {
		kernel.refuseWithObservation(appConn, frameBytes, "timed-out", readErr.Error())
		return
	}
	if writeErr := appConn.WriteFrame(ackBytes); writeErr != nil {
		fmt.Println(writeErr.Error())
	}
}

func (kernel *Kernel) deliverToLocalReceiver(ctx context.Context, frameBytes []byte, fields map[string]string) (string, string) {
	kernel.receiverMutex.Lock()
	receiverConn := kernel.receiver
	receiverName := kernel.receiverName
	kernel.receiverMutex.Unlock()
	if receiverConn == nil {
		return "not-promised", "no local app promised to receive hello messages"
	}
	if writeErr := receiverConn.WriteFrame(frameBytes); writeErr != nil {
		return "broken", writeErr.Error()
	}
	detail := fmt.Sprintf("delivered hello from %s to %s app %s", fields["from"], kernel.NodeName, receiverName)
	if recordErr := kernel.EvidenceLog.Record("app_deliver", "app/kernel", "kept", kernel.ProtocolCID, frameBytes, detail); recordErr != nil {
		fmt.Println(recordErr.Error())
	}
	select {
	case <-ctx.Done():
		return "kept", detail
	default:
		return "kept", detail
	}
}

func (kernel *Kernel) refuseUnsupported(frameConn FrameConn, frameBytes []byte, parseErr error) {
	detail := "unsupported pCID"
	if parseErr != nil {
		detail = parseErr.Error()
	}
	kernel.refuseWithObservation(frameConn, frameBytes, "refused", detail)
}

func (kernel *Kernel) refuseWithObservation(frameConn FrameConn, frameBytes []byte, outcome string, detail string) {
	if recordErr := kernel.EvidenceLog.Record("app_observation", "app/kernel", outcome, kernel.ProtocolCID, frameBytes, detail); recordErr != nil {
		fmt.Println(recordErr.Error())
	}
	observation, observationErr := observationEnvelope(kernel.ProtocolCID, kernel.NodeName+"-kernel", "local-app", outcome, detail, frameBytes)
	if observationErr != nil {
		fmt.Println(observationErr.Error())
		return
	}
	observationBytes, bytesErr := observation.Bytes()
	if bytesErr != nil {
		fmt.Println(bytesErr.Error())
		return
	}
	if writeErr := frameConn.WriteFrame(observationBytes); writeErr != nil {
		fmt.Println(writeErr.Error())
	}
	if closeErr := frameConn.Close(); closeErr != nil {
		fmt.Println(closeErr.Error())
	}
}

func (kernel *Kernel) writePeerObservation(frameConn FrameConn, frameBytes []byte, outcome string, detail string) {
	if recordErr := kernel.EvidenceLog.Record("peer_observation", "kernel/kernel", outcome, kernel.ProtocolCID, frameBytes, detail); recordErr != nil {
		fmt.Println(recordErr.Error())
	}
	observation, observationErr := observationEnvelope(kernel.ProtocolCID, kernel.NodeName+"-kernel", "peer-kernel", outcome, detail, frameBytes)
	if observationErr != nil {
		fmt.Println(observationErr.Error())
		return
	}
	observationBytes, bytesErr := observation.Bytes()
	if bytesErr != nil {
		fmt.Println(bytesErr.Error())
		return
	}
	if writeErr := frameConn.WriteFrame(observationBytes); writeErr != nil {
		fmt.Println(writeErr.Error())
	}
}
