package main

import (
	"fmt"
	"time"
)

// HelloApp is a tiny app agent. In receive mode it promises its local kernel it
// will receive hello messages. In send mode it sends a hello promise-message and
// waits for a local observation from its kernel.
type HelloApp struct {
	NodeName    string
	KernelAddr  string
	Mode        string
	Destination string
	Text        string
	ProtocolCID ProtocolCID
}

func (helloApp HelloApp) Run() error {
	frameConn, dialErr := dialFrameConn(helloApp.KernelAddr, 10*time.Second)
	if dialErr != nil {
		return dialErr
	}
	defer func() {
		if closeErr := frameConn.Close(); closeErr != nil {
			fmt.Println(closeErr.Error())
		}
	}()
	switch helloApp.Mode {
	case "receive":
		return helloApp.runReceive(frameConn)
	case "send":
		return helloApp.runSend(frameConn)
	default:
		return fmt.Errorf("unknown hello mode %q", helloApp.Mode)
	}
}

func (helloApp HelloApp) runReceive(frameConn FrameConn) error {
	promiseEnvelope, promiseErr := NewEnvelope(helloApp.ProtocolCID, map[string]string{
		"kind": "receive_promise_v1",
		"from": helloApp.NodeName + "-hello-app",
		"node": helloApp.NodeName,
		"text": "I promise to receive hello_v1 messages for this bounded poc2 run.",
	})
	if promiseErr != nil {
		return promiseErr
	}
	promiseBytes, bytesErr := promiseEnvelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	if writeErr := frameConn.WriteFrame(promiseBytes); writeErr != nil {
		return writeErr
	}
	helloBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		return readErr
	}
	helloEnvelope, parseErr := ParseEnvelope(helloBytes)
	if parseErr != nil {
		return parseErr
	}
	kind, fields, kindErr := envelopeKind(helloEnvelope)
	if kindErr != nil {
		return kindErr
	}
	if kind != "hello_v1" {
		return fmt.Errorf("receive mode got payload kind %s", kind)
	}
	fmt.Printf("%s app received from %s: %s\n", helloApp.NodeName, fields["from"], fields["text"])
	return nil
}

func (helloApp HelloApp) runSend(frameConn FrameConn) error {
	helloEnvelope, helloErr := NewEnvelope(helloApp.ProtocolCID, map[string]string{
		"kind": "hello_v1",
		"from": helloApp.NodeName + "-hello-app",
		"to":   helloApp.Destination,
		"text": helloApp.Text,
	})
	if helloErr != nil {
		return helloErr
	}
	helloBytes, bytesErr := helloEnvelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	if writeErr := frameConn.WriteFrame(helloBytes); writeErr != nil {
		return writeErr
	}
	observationBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		return readErr
	}
	observationEnvelope, parseErr := ParseEnvelope(observationBytes)
	if parseErr != nil {
		return parseErr
	}
	kind, fields, kindErr := envelopeKind(observationEnvelope)
	if kindErr != nil {
		return kindErr
	}
	if kind != "observation_v1" {
		return fmt.Errorf("send mode got payload kind %s", kind)
	}
	fmt.Printf("%s app observed %s: %s\n", helloApp.NodeName, fields["outcome"], fields["text"])
	if fields["outcome"] != "kept" {
		return fmt.Errorf("hello was not kept: %s", fields["text"])
	}
	return nil
}
