package kernel

import (
	"context"
	"fmt"
	"net"
	"sync"

	"promisegrid.dev/wire-lab/implementations/poc5/lib"
)

// Receiver records one local app's promise to receive a pCID.
type Receiver struct {
	AppName string
	Conn    lib.FrameConn
}

// Kernel is a local app/kernel boundary. It accepts receive promises, delivers
// exact local message bytes, refuses unsupported local pCIDs, and records local
// kernel evidence only.
//
// Intent: In poc5, relays own neighbor transport. The kernel stays a local
// promise boundary rather than a multi-hop router. Source: DI-rarim.
type Kernel struct {
	NodeName           string
	AppListen          string
	ReceiveProtocolCID lib.ProtocolCID
	EvidenceLog        *lib.EvidenceLog
	receiverMutex      sync.Mutex
	receivers          map[string]Receiver
}

// Run starts the local app listener.
func (kernel *Kernel) Run(ctx context.Context) error {
	runContext, cancel := context.WithCancel(ctx)
	defer cancel()
	if kernel.receivers == nil {
		kernel.receivers = make(map[string]Receiver)
	}
	if kernel.ReceiveProtocolCID.String() == "cidv1-raw-sha2-256:0000000000000000000000000000000000000000000000000000000000000000" {
		kernel.ReceiveProtocolCID = ReceiveProtocolCID()
	}
	errorChannel := make(chan error, 1)
	go kernel.listen(runContext, errorChannel)
	select {
	case <-runContext.Done():
		return nil
	case listenErr := <-errorChannel:
		cancel()
		return listenErr
	}
}

func (kernel *Kernel) listen(ctx context.Context, errorChannel chan<- error) {
	listener, listenErr := net.Listen("tcp", kernel.AppListen)
	if listenErr != nil {
		errorChannel <- fmt.Errorf("listen app %s: %w", kernel.AppListen, listenErr)
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
		// Intent: Kernel shutdown is local process lifecycle, not a remote control
		// surface. Source: DI-rarim.
		if closeErr := listener.Close(); closeErr != nil {
			fmt.Println(closeErr.Error())
		}
	}()
	for {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			select {
			case <-ctx.Done():
				return
			default:
				errorChannel <- fmt.Errorf("accept app: %w", acceptErr)
				return
			}
		}
		go kernel.handleAppConnection(ctx, conn)
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
	defer closeFrame(frameConn)
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

// DeliverEnvelope writes exact bytes to the app that promised the envelope pCID.
func (kernel *Kernel) DeliverEnvelope(ctx context.Context, envelopeBytes []byte) error {
	envelope, parseErr := lib.ParseEnvelope(envelopeBytes)
	if parseErr != nil {
		return parseErr
	}
	kernel.deliverToLocalReceiver(ctx, envelopeBytes, envelope)
	return nil
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
