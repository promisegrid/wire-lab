package echo

import (
	"fmt"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc3/kernel"
	"promisegrid.dev/wire-lab/implementations/poc3/lib"
)

// EchoApp is an app agent that can promise to receive echo messages and choose
// to make a new echo response promise-message.
type EchoApp struct {
	NodeName    string
	AppName     string
	KernelAddr  string
	Mode        string
	Destination string
	Text        string
}

// Run executes the echo app mode.
func (echoApp EchoApp) Run() error {
	if echoApp.AppName == "" {
		echoApp.AppName = echoApp.NodeName + "-echo-app"
	}
	switch echoApp.Mode {
	case "serve":
		return echoApp.runServe()
	case "ask":
		return echoApp.runAsk()
	default:
		return fmt.Errorf("unknown echo mode %q", echoApp.Mode)
	}
}

func (echoApp EchoApp) runServe() error {
	receiveConn, registerErr := echoApp.registerReceiver("I promise to receive echo requests and may make a local echo response.")
	if registerErr != nil {
		return registerErr
	}
	defer closeFrame(receiveConn)
	requestBytes, readErr := receiveConn.ReadFrame()
	if readErr != nil {
		return readErr
	}
	requestEnvelope, parseErr := lib.ParseEnvelope(requestBytes)
	if parseErr != nil {
		return parseErr
	}
	kind, fields, kindErr := lib.EnvelopeKind(requestEnvelope)
	if kindErr != nil {
		return kindErr
	}
	if !requestEnvelope.ProtocolCID.Equal(ProtocolCID()) || kind != "echo_request_v1" {
		return fmt.Errorf("echo server got unexpected message kind %s", kind)
	}
	fmt.Printf("%s judged echo request from %s kept: %s\n", echoApp.AppName, fields["from"], fields["text"])
	return echoApp.sendMessage(fields["from_node"], "echo_response_v1", fields["text"])
}

func (echoApp EchoApp) runAsk() error {
	receiveConn, registerErr := echoApp.registerReceiver("I promise to receive echo responses for this bounded poc3 run.")
	if registerErr != nil {
		return registerErr
	}
	defer closeFrame(receiveConn)
	if sendErr := echoApp.sendMessage(echoApp.Destination, "echo_request_v1", echoApp.Text); sendErr != nil {
		return sendErr
	}
	responseBytes, readErr := receiveConn.ReadFrame()
	if readErr != nil {
		return readErr
	}
	responseEnvelope, parseErr := lib.ParseEnvelope(responseBytes)
	if parseErr != nil {
		return parseErr
	}
	kind, fields, kindErr := lib.EnvelopeKind(responseEnvelope)
	if kindErr != nil {
		return kindErr
	}
	if !responseEnvelope.ProtocolCID.Equal(ProtocolCID()) || kind != "echo_response_v1" {
		return fmt.Errorf("echo asker got unexpected message kind %s", kind)
	}
	fmt.Printf("%s judged echo response from %s kept: %s\n", echoApp.AppName, fields["from"], fields["text"])
	return nil
}

func (echoApp EchoApp) registerReceiver(text string) (lib.FrameConn, error) {
	frameConn, dialErr := lib.DialFrameConn(echoApp.KernelAddr, 10*time.Second)
	if dialErr != nil {
		return lib.FrameConn{}, dialErr
	}
	if err := lib.WriteReceivePromise(frameConn, kernel.ReceiveProtocolCID(), echoApp.NodeName, echoApp.AppName, ProtocolCID(), text); err != nil {
		closeFrame(frameConn)
		return lib.FrameConn{}, err
	}
	return frameConn, nil
}

func (echoApp EchoApp) sendMessage(destination string, kind string, text string) error {
	frameConn, dialErr := lib.DialFrameConn(echoApp.KernelAddr, 10*time.Second)
	if dialErr != nil {
		return dialErr
	}
	defer closeFrame(frameConn)
	envelope, envelopeErr := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":      kind,
		"from":      echoApp.AppName,
		"from_node": echoApp.NodeName,
		"to":        destination,
		"text":      text,
	})
	if envelopeErr != nil {
		return envelopeErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	if writeErr := frameConn.WriteFrame(envelopeBytes); writeErr != nil {
		return writeErr
	}
	fmt.Printf("%s made %s promise-message to %s: %s\n", echoApp.AppName, kind, destination, text)
	return nil
}

func closeFrame(frameConn lib.FrameConn) {
	if closeErr := frameConn.Close(); closeErr != nil {
		fmt.Println(closeErr.Error())
	}
}
