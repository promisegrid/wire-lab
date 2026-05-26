package relay

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc4/kernel"
	"promisegrid.dev/wire-lab/implementations/poc4/lib"
)

// RelayApp is an app-level promiser. It owns neighbor TCP and makes local
// next-hop relay promises; it is not a kernel router.
type RelayApp struct {
	NodeName    string
	AppName     string
	KernelAddr  string
	ListenAddr  string
	RouteTable  map[string]string
	EvidenceLog *lib.EvidenceLog
}

// Run registers the relay with the local kernel and starts local plus remote
// relay handling loops.
func (relayApp RelayApp) Run(ctx context.Context) error {
	if relayApp.AppName == "" {
		relayApp.AppName = relayApp.NodeName + "-relay-app"
	}
	if relayApp.EvidenceLog == nil {
		relayApp.EvidenceLog = lib.NewEvidenceLog(relayApp.NodeName, relayApp.AppName)
	}
	receiveConn, registerErr := relayApp.registerReceiver()
	if registerErr != nil {
		return registerErr
	}
	defer closeFrame(receiveConn)
	errorChannel := make(chan error, 2)
	go relayApp.listenRemote(ctx, errorChannel)
	go relayApp.handleLocalLoop(ctx, receiveConn, errorChannel)
	select {
	case <-ctx.Done():
		return nil
	case err := <-errorChannel:
		return err
	}
}

// Wrap builds a relay-pCID envelope carrying exact inner envelope bytes.
func Wrap(originNode string, originApp string, targetNode string, targetApp string, innerBytes []byte, requestHash string) (lib.Envelope, error) {
	return lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":         "relay_forward_v1",
		"origin_node":  originNode,
		"origin_app":   originApp,
		"target_node":  targetNode,
		"target_app":   targetApp,
		"inner_hex":    lib.HexBytes(innerBytes),
		"request_hash": requestHash,
	})
}

// SendViaLocalRelay submits a relay wrapper to the local kernel for delivery to
// the local relay app.
func SendViaLocalRelay(kernelAddr string, wrapper lib.Envelope) error {
	wrapperBytes, bytesErr := wrapper.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	frameConn, dialErr := lib.DialFrameConn(kernelAddr, 10*time.Second)
	if dialErr != nil {
		return dialErr
	}
	defer closeFrame(frameConn)
	return frameConn.WriteFrame(wrapperBytes)
}

// SendInnerViaLocalRelay wraps one exact app envelope in a relay promise and
// submits the wrapper through the local kernel to the local relay app.
func SendInnerViaLocalRelay(kernelAddr string, originNode string, originApp string, targetNode string, targetApp string, inner lib.Envelope, requestHash string) ([]byte, error) {
	innerBytes, bytesErr := inner.Bytes()
	if bytesErr != nil {
		return nil, bytesErr
	}
	wrapper, wrapErr := Wrap(originNode, originApp, targetNode, targetApp, innerBytes, requestHash)
	if wrapErr != nil {
		return nil, wrapErr
	}
	if sendErr := SendViaLocalRelay(kernelAddr, wrapper); sendErr != nil {
		return innerBytes, sendErr
	}
	return innerBytes, nil
}

func (relayApp RelayApp) registerReceiver() (lib.FrameConn, error) {
	frameConn, dialErr := lib.DialFrameConn(relayApp.KernelAddr, 10*time.Second)
	if dialErr != nil {
		return lib.FrameConn{}, dialErr
	}
	if err := lib.WriteReceivePromise(frameConn, kernel.ReceiveProtocolCID(), relayApp.NodeName, relayApp.AppName, ProtocolCID(), "I promise to receive relay wrappers and try locally promised next hops."); err != nil {
		closeFrame(frameConn)
		return lib.FrameConn{}, err
	}
	return frameConn, nil
}

func (relayApp RelayApp) listenRemote(ctx context.Context, errorChannel chan<- error) {
	listener, listenErr := net.Listen("tcp", relayApp.ListenAddr)
	if listenErr != nil {
		errorChannel <- listenErr
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
				errorChannel <- acceptErr
				return
			}
		}
		go relayApp.handleRemoteConn(ctx, conn)
	}
}

func (relayApp RelayApp) handleLocalLoop(ctx context.Context, receiveConn lib.FrameConn, errorChannel chan<- error) {
	for {
		frameBytes, readErr := receiveConn.ReadFrame()
		if readErr != nil {
			select {
			case <-ctx.Done():
				return
			default:
				errorChannel <- readErr
				return
			}
		}
		if err := relayApp.processRelayFrame(ctx, frameBytes, "local-kernel"); err != nil {
			relayApp.record("relay_process", "relay/local", "broken", frameBytes, err.Error())
		}
	}
}

func (relayApp RelayApp) handleRemoteConn(ctx context.Context, conn net.Conn) {
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			fmt.Println(closeErr.Error())
		}
	}()
	frameConn := lib.NewFrameConn(conn)
	frameBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		relayApp.record("relay_read", "relay/relay", "broken", nil, readErr.Error())
		return
	}
	if err := relayApp.processRelayFrame(ctx, frameBytes, "neighbor"); err != nil {
		relayApp.record("relay_process", "relay/relay", "broken", frameBytes, err.Error())
	}
}

func (relayApp RelayApp) processRelayFrame(ctx context.Context, frameBytes []byte, source string) error {
	envelope, parseErr := lib.ParseEnvelope(frameBytes)
	if parseErr != nil {
		return parseErr
	}
	if !envelope.ProtocolCID.Equal(ProtocolCID()) {
		return fmt.Errorf("relay got unsupported pCID %s", envelope.ProtocolCID.String())
	}
	kind, fields, kindErr := lib.EnvelopeKind(envelope)
	if kindErr != nil {
		return kindErr
	}
	if kind != "relay_forward_v1" {
		return fmt.Errorf("relay got unsupported kind %s", kind)
	}
	targetNode := fields["target_node"]
	innerBytes, innerErr := lib.ParseHexBytes(fields["inner_hex"])
	if innerErr != nil {
		return innerErr
	}
	relayApp.record("relay_receive", "relay/"+source, "kept", frameBytes, "received wrapper for "+targetNode+"/"+fields["target_app"])
	if targetNode == relayApp.NodeName {
		if err := relayApp.deliverLocal(ctx, innerBytes, fields); err != nil {
			return err
		}
		return nil
	}
	nextHop, ok := relayApp.RouteTable[targetNode]
	if !ok {
		relayApp.record("relay_forward", "relay/relay", "not-promised", frameBytes, "no local route promise for "+targetNode)
		return nil
	}
	peerConn, dialErr := lib.DialFrameConn(nextHop, 10*time.Second)
	if dialErr != nil {
		return dialErr
	}
	defer closeFrame(peerConn)
	if writeErr := peerConn.WriteFrame(frameBytes); writeErr != nil {
		return writeErr
	}
	fmt.Printf("%s promised hop toward %s via %s for %s\n", relayApp.AppName, targetNode, nextHop, lib.HashExactBytes(frameBytes))
	relayApp.record("relay_forward", "relay/relay", "kept", frameBytes, "forwarded toward "+targetNode+" via "+nextHop)
	return nil
}

func (relayApp RelayApp) deliverLocal(ctx context.Context, innerBytes []byte, fields map[string]string) error {
	frameConn, dialErr := lib.DialFrameConn(relayApp.KernelAddr, 10*time.Second)
	if dialErr != nil {
		return dialErr
	}
	defer closeFrame(frameConn)
	if writeErr := frameConn.WriteFrame(innerBytes); writeErr != nil {
		return writeErr
	}
	select {
	case <-ctx.Done():
		relayApp.record("relay_deliver", "relay/kernel", "kept", innerBytes, "delivered during shutdown to "+fields["target_app"])
	default:
		relayApp.record("relay_deliver", "relay/kernel", "kept", innerBytes, "delivered to "+fields["target_app"])
	}
	fmt.Printf("%s completed final relay delivery to %s/%s for %s\n", relayApp.AppName, fields["target_node"], fields["target_app"], lib.HashExactBytes(innerBytes))
	return nil
}

func (relayApp RelayApp) record(event string, boundary string, outcome string, exactBytes []byte, detail string) {
	if relayApp.EvidenceLog == nil {
		return
	}
	if recordErr := relayApp.EvidenceLog.Record(event, boundary, outcome, ProtocolCID(), exactBytes, detail); recordErr != nil {
		fmt.Println(recordErr.Error())
	}
}

// ParseRoutes turns "target=addr,target=addr" into local relay promises.
func ParseRoutes(routeText string) (map[string]string, error) {
	routes := make(map[string]string)
	if strings.TrimSpace(routeText) == "" {
		return routes, nil
	}
	for _, entry := range strings.Split(routeText, ",") {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid route promise %q", entry)
		}
		routes[parts[0]] = parts[1]
	}
	return routes, nil
}

func closeFrame(frameConn lib.FrameConn) {
	if closeErr := frameConn.Close(); closeErr != nil {
		fmt.Println(closeErr.Error())
	}
}
