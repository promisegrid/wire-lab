package hello

import (
	"fmt"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc3/kernel"
	"promisegrid.dev/wire-lab/implementations/poc3/lib"
)

// HelloApp is an app agent that can promise to receive hello messages or make a
// hello promise-message for another node.
type HelloApp struct {
	NodeName    string
	AppName     string
	KernelAddr  string
	Mode        string
	Destination string
	Text        string
}

// Run executes the hello app mode.
func (helloApp HelloApp) Run() error {
	if helloApp.AppName == "" {
		helloApp.AppName = helloApp.NodeName + "-hello-app"
	}
	switch helloApp.Mode {
	case "receive":
		return helloApp.runReceive()
	case "send":
		return helloApp.runSend()
	default:
		return fmt.Errorf("unknown hello mode %q", helloApp.Mode)
	}
}

func (helloApp HelloApp) runReceive() error {
	frameConn, dialErr := lib.DialFrameConn(helloApp.KernelAddr, 10*time.Second)
	if dialErr != nil {
		return dialErr
	}
	defer closeFrame(frameConn)
	if err := lib.WriteReceivePromise(frameConn, kernel.ReceiveProtocolCID(), helloApp.NodeName, helloApp.AppName, ProtocolCID(), "I promise to receive hello_v1 messages for this bounded poc3 run."); err != nil {
		return err
	}
	helloBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		return readErr
	}
	envelope, parseErr := lib.ParseEnvelope(helloBytes)
	if parseErr != nil {
		return parseErr
	}
	kind, fields, kindErr := lib.EnvelopeKind(envelope)
	if kindErr != nil {
		return kindErr
	}
	if !envelope.ProtocolCID.Equal(ProtocolCID()) || kind != "hello_v1" {
		return fmt.Errorf("hello receiver got unexpected message kind %s", kind)
	}
	fmt.Printf("%s judged hello from %s kept: %s\n", helloApp.AppName, fields["from"], fields["text"])
	return nil
}

func (helloApp HelloApp) runSend() error {
	frameConn, dialErr := lib.DialFrameConn(helloApp.KernelAddr, 10*time.Second)
	if dialErr != nil {
		return dialErr
	}
	defer closeFrame(frameConn)
	envelope, envelopeErr := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":      "hello_v1",
		"from":      helloApp.AppName,
		"from_node": helloApp.NodeName,
		"to":        helloApp.Destination,
		"text":      helloApp.Text,
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
	fmt.Printf("%s made hello promise-message to %s: %s\n", helloApp.AppName, helloApp.Destination, helloApp.Text)
	return nil
}

func closeFrame(frameConn lib.FrameConn) {
	if closeErr := frameConn.Close(); closeErr != nil {
		fmt.Println(closeErr.Error())
	}
}
