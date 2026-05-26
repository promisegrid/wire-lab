package echo

import (
	"fmt"

	"promisegrid.dev/wire-lab/implementations/poc5/kernel"
	"promisegrid.dev/wire-lab/implementations/poc5/lib"
	"promisegrid.dev/wire-lab/implementations/poc5/relay"
)

// EchoApp is an app agent that can ask for, or fulfill, one echo promise.
type EchoApp struct {
	NodeName   string
	AppName    string
	KernelAddr string
	Mode       string
	TargetNode string
	TargetApp  string
	Text       string
}

// Run executes one bounded echo role.
func (echoApp EchoApp) Run() error {
	if echoApp.AppName == "" {
		echoApp.AppName = echoApp.NodeName + "-echo-app"
	}
	switch echoApp.Mode {
	case "client":
		return echoApp.runClient()
	case "serve":
		return echoApp.runServe()
	default:
		return fmt.Errorf("unknown echo mode %q", echoApp.Mode)
	}
}

func (echoApp EchoApp) runClient() error {
	// Intent: The echo client makes the reciprocal receive promise first and
	// later judges the app result itself; relay and kernel evidence remain
	// transport evidence only.
	// Source: DI-rarim
	frameConn, registerErr := lib.RegisterReceiver(echoApp.KernelAddr, kernel.ReceiveProtocolCID(), echoApp.NodeName, echoApp.AppName, ProtocolCID(), "I promise to receive one echo_response_v1 and judge it locally.")
	if registerErr != nil {
		return registerErr
	}
	defer closeFrame(frameConn)
	requestEnvelope, requestErr := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":      "echo_request_v1",
		"from":      echoApp.AppName,
		"from_node": echoApp.NodeName,
		"to":        echoApp.TargetApp,
		"text":      echoApp.Text,
	})
	if requestErr != nil {
		return requestErr
	}
	requestBytes, sendErr := relay.SendInnerViaLocalRelay(echoApp.KernelAddr, echoApp.NodeName, echoApp.AppName, echoApp.TargetNode, echoApp.TargetApp, requestEnvelope, "")
	if sendErr != nil {
		return sendErr
	}
	_, kind, fields, _, readErr := lib.ReadEnvelope(frameConn)
	if readErr != nil {
		return readErr
	}
	if kind != "echo_response_v1" {
		return fmt.Errorf("echo client got unexpected kind %s", kind)
	}
	if fields["request_hash"] != lib.HashExactBytes(requestBytes) || fields["text"] != echoApp.Text {
		return fmt.Errorf("echo response did not match request")
	}
	fmt.Printf("%s judged echo kept from %s: %s\n", echoApp.AppName, fields["from"], fields["text"])
	return nil
}

func (echoApp EchoApp) runServe() error {
	frameConn, registerErr := lib.RegisterReceiver(echoApp.KernelAddr, kernel.ReceiveProtocolCID(), echoApp.NodeName, echoApp.AppName, ProtocolCID(), "I promise to receive one echo_request_v1 and return the same text if the requester promises to receive it.")
	if registerErr != nil {
		return registerErr
	}
	defer closeFrame(frameConn)
	envelope, kind, fields, requestBytes, readErr := lib.ReadEnvelope(frameConn)
	if readErr != nil {
		return readErr
	}
	if !envelope.ProtocolCID.Equal(ProtocolCID()) || kind != "echo_request_v1" {
		return fmt.Errorf("echo server got unexpected kind %s", kind)
	}
	responseEnvelope, responseErr := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":         "echo_response_v1",
		"from":         echoApp.AppName,
		"from_node":    echoApp.NodeName,
		"to":           fields["from"],
		"to_node":      fields["from_node"],
		"request_hash": lib.HashExactBytes(requestBytes),
		"text":         fields["text"],
	})
	if responseErr != nil {
		return responseErr
	}
	if _, sendErr := relay.SendInnerViaLocalRelay(echoApp.KernelAddr, echoApp.NodeName, echoApp.AppName, fields["from_node"], fields["from"], responseEnvelope, lib.HashExactBytes(requestBytes)); sendErr != nil {
		return sendErr
	}
	fmt.Printf("%s judged echo request from %s kept and returned text\n", echoApp.AppName, fields["from"])
	return nil
}

func closeFrame(frameConn lib.FrameConn) {
	if closeErr := frameConn.Close(); closeErr != nil {
		fmt.Println(closeErr.Error())
	}
}
