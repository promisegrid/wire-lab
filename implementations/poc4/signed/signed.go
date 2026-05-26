package signed

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"

	"promisegrid.dev/wire-lab/implementations/poc4/kernel"
	"promisegrid.dev/wire-lab/implementations/poc4/lib"
	"promisegrid.dev/wire-lab/implementations/poc4/relay"
)

// SignedApp is an app agent that signs exact bytes as local evidence. It does
// not confer authority, permission, or global trust.
type SignedApp struct {
	NodeName   string
	AppName    string
	KernelAddr string
	Mode       string
}

// Run executes a bounded signed-app role for the demo.
func (signedApp SignedApp) Run() error {
	if signedApp.AppName == "" {
		signedApp.AppName = signedApp.NodeName + "-signed-app"
	}
	switch signedApp.Mode {
	case "serve":
		return signedApp.runServe()
	case "idle":
		fmt.Printf("%s made no signed promise in this bounded run\n", signedApp.AppName)
		return nil
	default:
		return fmt.Errorf("unknown signed mode %q", signedApp.Mode)
	}
}

func (signedApp SignedApp) runServe() error {
	// Intent: The signed app promises to answer one request while keeping
	// signature evidence at the app layer, not in relay or kernel routing.
	// Source: DI-bigub
	frameConn, registerErr := lib.RegisterReceiver(signedApp.KernelAddr, kernel.ReceiveProtocolCID(), signedApp.NodeName, signedApp.AppName, ProtocolCID(), "I promise to receive one sign_request_v1 and return signed exact-byte evidence if the requester promises to receive it.")
	if registerErr != nil {
		return registerErr
	}
	defer closeFrame(frameConn)
	envelope, kind, fields, frameBytes, readErr := lib.ReadEnvelope(frameConn)
	if readErr != nil {
		return readErr
	}
	if !envelope.ProtocolCID.Equal(ProtocolCID()) || kind != "sign_request_v1" {
		return fmt.Errorf("signed server got unexpected message kind %s", kind)
	}
	resultEnvelope, signErr := signedApp.resultEnvelope(fields, lib.HashExactBytes(frameBytes))
	if signErr != nil {
		return signErr
	}
	if _, sendErr := relay.SendInnerViaLocalRelay(signedApp.KernelAddr, signedApp.NodeName, signedApp.AppName, fields["from_node"], fields["from"], resultEnvelope, lib.HashExactBytes(frameBytes)); sendErr != nil {
		return sendErr
	}
	fmt.Printf("%s judged sign request from %s kept and returned signed evidence\n", signedApp.AppName, fields["from"])
	return nil
}

func (signedApp SignedApp) resultEnvelope(requestFields map[string]string, requestHash string) (lib.Envelope, error) {
	signableEnvelope, envelopeErr := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":         "signed_result_v1",
		"from":         signedApp.AppName,
		"from_node":    signedApp.NodeName,
		"to":           requestFields["from"],
		"to_node":      requestFields["from_node"],
		"request_hash": requestHash,
		"text":         requestFields["text"],
	})
	if envelopeErr != nil {
		return lib.Envelope{}, envelopeErr
	}
	return signEnvelope(signedApp.NodeName+":"+signedApp.AppName, signableEnvelope)
}

// NewSignRequest builds the unsigned request that another app can carry to a
// signed app through relay promises.
func NewSignRequest(fromNode string, fromApp string, targetApp string, text string) (lib.Envelope, error) {
	return lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":      "sign_request_v1",
		"from":      fromApp,
		"from_node": fromNode,
		"to":        targetApp,
		"text":      text,
	})
}

// VerifyEnvelope verifies the app-level proof slot and returns the decoded
// payload fields for local promise judgment.
func VerifyEnvelope(envelope lib.Envelope) (map[string]string, error) {
	kind, fields, kindErr := lib.EnvelopeKind(envelope)
	if kindErr != nil {
		return nil, kindErr
	}
	if !envelope.ProtocolCID.Equal(ProtocolCID()) || kind != "signed_result_v1" {
		return nil, fmt.Errorf("unexpected signed result kind %s", kind)
	}
	if len(envelope.ExtraSlots) != 1 {
		return nil, fmt.Errorf("signed result must carry exactly one proof slot")
	}
	signatureSlot, slotErr := lib.DecodeSignatureSlot(envelope.ExtraSlots[0])
	if slotErr != nil {
		return nil, slotErr
	}
	signableBytes, signableErr := lib.Envelope{ProtocolCID: envelope.ProtocolCID, Payload: envelope.Payload}.Bytes()
	if signableErr != nil {
		return nil, signableErr
	}
	if !ed25519.Verify(ed25519.PublicKey(signatureSlot.PublicKey), signableBytes, signatureSlot.Signature) {
		return nil, fmt.Errorf("signature did not verify")
	}
	return fields, nil
}

func signEnvelope(seedText string, signableEnvelope lib.Envelope) (lib.Envelope, error) {
	signableBytes, bytesErr := signableEnvelope.Bytes()
	if bytesErr != nil {
		return lib.Envelope{}, bytesErr
	}
	publicKey, privateKey := deterministicKeyPair(seedText)
	signature := ed25519.Sign(privateKey, signableBytes)
	return lib.Envelope{
		ProtocolCID: signableEnvelope.ProtocolCID,
		Payload:     signableEnvelope.Payload,
		ExtraSlots:  [][]byte{lib.EncodeSignatureSlot(lib.SignatureSlot{PublicKey: publicKey, Signature: signature})},
	}, nil
}

func deterministicKeyPair(seedText string) (ed25519.PublicKey, ed25519.PrivateKey) {
	seed := sha256.Sum256([]byte("poc4 demo key: " + seedText))
	privateKey := ed25519.NewKeyFromSeed(seed[:])
	publicKey := privateKey.Public().(ed25519.PublicKey)
	return publicKey, privateKey
}

func closeFrame(frameConn lib.FrameConn) {
	if closeErr := frameConn.Close(); closeErr != nil {
		fmt.Println(closeErr.Error())
	}
}
