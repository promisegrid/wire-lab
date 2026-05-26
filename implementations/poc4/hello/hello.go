package hello

import (
	"fmt"

	"promisegrid.dev/wire-lab/implementations/poc4/kernel"
	"promisegrid.dev/wire-lab/implementations/poc4/lib"
	"promisegrid.dev/wire-lab/implementations/poc4/relay"
	"promisegrid.dev/wire-lab/implementations/poc4/signed"
)

// HelloApp is an app agent that can ask for signed evidence over a hello text.
type HelloApp struct {
	NodeName   string
	AppName    string
	KernelAddr string
	Mode       string
	TargetNode string
	TargetApp  string
	Text       string
}

// Run executes one bounded hello role.
func (helloApp HelloApp) Run() error {
	if helloApp.AppName == "" {
		helloApp.AppName = helloApp.NodeName + "-hello-app"
	}
	switch helloApp.Mode {
	case "ask-signed":
		return helloApp.runAskSigned()
	case "idle":
		fmt.Printf("%s made no hello network promise in this bounded run\n", helloApp.AppName)
		return nil
	default:
		return fmt.Errorf("unknown hello mode %q", helloApp.Mode)
	}
}

func (helloApp HelloApp) runAskSigned() error {
	// Intent: Alice's hello app promises to receive a signed result before it
	// asks Dave to sign, so the reciprocal receive promise stays at the app
	// layer instead of becoming kernel permission or RPC dispatch.
	// Source: DI-bigub
	frameConn, registerErr := lib.RegisterReceiver(helloApp.KernelAddr, kernel.ReceiveProtocolCID(), helloApp.NodeName, helloApp.AppName, signed.ProtocolCID(), "I promise to receive a signed_result_v1 and judge it locally.")
	if registerErr != nil {
		return registerErr
	}
	defer closeFrame(frameConn)
	requestEnvelope, requestErr := signed.NewSignRequest(helloApp.NodeName, helloApp.AppName, helloApp.TargetApp, helloApp.Text)
	if requestErr != nil {
		return requestErr
	}
	requestBytes, sendErr := relay.SendInnerViaLocalRelay(helloApp.KernelAddr, helloApp.NodeName, helloApp.AppName, helloApp.TargetNode, helloApp.TargetApp, requestEnvelope, "")
	if sendErr != nil {
		return sendErr
	}
	fmt.Printf("%s promised to receive signed evidence and asked %s/%s to sign hello text via relay path %s\n", helloApp.AppName, helloApp.TargetNode, helloApp.TargetApp, lib.HashExactBytes(requestBytes))
	resultEnvelope, _, _, resultBytes, readErr := lib.ReadEnvelope(frameConn)
	if readErr != nil {
		return readErr
	}
	fields, verifyErr := signed.VerifyEnvelope(resultEnvelope)
	if verifyErr != nil {
		return verifyErr
	}
	fmt.Printf("%s judged signed hello result kept from %s: %s (%s)\n", helloApp.AppName, fields["from"], fields["text"], lib.HashExactBytes(resultBytes))
	return nil
}

func closeFrame(frameConn lib.FrameConn) {
	if closeErr := frameConn.Close(); closeErr != nil {
		fmt.Println(closeErr.Error())
	}
}
