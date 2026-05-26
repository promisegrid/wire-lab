package signed

import (
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc3/kernel"
	"promisegrid.dev/wire-lab/implementations/poc3/lib"
)

// SignedApp is an app agent that signs canonical envelope bytes or verifies
// them as exact-byte evidence without making global trust claims.
type SignedApp struct {
	NodeName    string
	AppName     string
	KernelAddr  string
	Mode        string
	Destination string
	Text        string
}

// Run executes the signed app mode.
func (signedApp SignedApp) Run() error {
	if signedApp.AppName == "" {
		signedApp.AppName = signedApp.NodeName + "-signed-app"
	}
	switch signedApp.Mode {
	case "receive":
		return signedApp.runReceive()
	case "send":
		return signedApp.runSend()
	default:
		return fmt.Errorf("unknown signed mode %q", signedApp.Mode)
	}
}

func (signedApp SignedApp) runReceive() error {
	frameConn, dialErr := lib.DialFrameConn(signedApp.KernelAddr, 10*time.Second)
	if dialErr != nil {
		return dialErr
	}
	defer closeFrame(frameConn)
	if err := lib.WriteReceivePromise(frameConn, kernel.ReceiveProtocolCID(), signedApp.NodeName, signedApp.AppName, ProtocolCID(), "I promise to receive signed_note_v1 messages and judge them locally."); err != nil {
		return err
	}
	signedBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		return readErr
	}
	envelope, parseErr := lib.ParseEnvelope(signedBytes)
	if parseErr != nil {
		return parseErr
	}
	kind, fields, kindErr := lib.EnvelopeKind(envelope)
	if kindErr != nil {
		return kindErr
	}
	if !envelope.ProtocolCID.Equal(ProtocolCID()) || kind != "signed_note_v1" {
		return fmt.Errorf("signed receiver got unexpected message kind %s", kind)
	}
	if len(envelope.ExtraSlots) != 1 {
		return fmt.Errorf("signed envelope must carry exactly one proof slot")
	}
	signatureSlot, slotErr := lib.DecodeSignatureSlot(envelope.ExtraSlots[0])
	if slotErr != nil {
		return slotErr
	}
	signableBytes, signableErr := lib.Envelope{ProtocolCID: envelope.ProtocolCID, Payload: envelope.Payload}.Bytes()
	if signableErr != nil {
		return signableErr
	}
	if !ed25519.Verify(ed25519.PublicKey(signatureSlot.PublicKey), signableBytes, signatureSlot.Signature) {
		return fmt.Errorf("signature did not verify")
	}
	fmt.Printf("%s verified exact bytes from %s and judged signed note kept locally: %s\n", signedApp.AppName, fields["from"], fields["text"])
	return nil
}

func (signedApp SignedApp) runSend() error {
	frameConn, dialErr := lib.DialFrameConn(signedApp.KernelAddr, 10*time.Second)
	if dialErr != nil {
		return dialErr
	}
	defer closeFrame(frameConn)
	signableEnvelope, envelopeErr := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":      "signed_note_v1",
		"from":      signedApp.AppName,
		"from_node": signedApp.NodeName,
		"to":        signedApp.Destination,
		"text":      signedApp.Text,
	})
	if envelopeErr != nil {
		return envelopeErr
	}
	signableBytes, bytesErr := signableEnvelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	publicKey, privateKey := deterministicKeyPair(signedApp.NodeName + ":" + signedApp.AppName)
	signature := ed25519.Sign(privateKey, signableBytes)
	signedEnvelope := lib.Envelope{
		ProtocolCID: signableEnvelope.ProtocolCID,
		Payload:     signableEnvelope.Payload,
		ExtraSlots:  [][]byte{lib.EncodeSignatureSlot(lib.SignatureSlot{PublicKey: publicKey, Signature: signature})},
	}
	signedBytes, signedErr := signedEnvelope.Bytes()
	if signedErr != nil {
		return signedErr
	}
	if writeErr := frameConn.WriteFrame(signedBytes); writeErr != nil {
		return writeErr
	}
	fmt.Printf("%s made signed promise-message to %s: %s\n", signedApp.AppName, signedApp.Destination, signedApp.Text)
	return nil
}

func deterministicKeyPair(seedText string) (ed25519.PublicKey, ed25519.PrivateKey) {
	seed := sha256.Sum256([]byte("poc3 demo key: " + seedText))
	privateKey := ed25519.NewKeyFromSeed(seed[:])
	publicKey := privateKey.Public().(ed25519.PublicKey)
	return publicKey, privateKey
}

func closeFrame(frameConn lib.FrameConn) {
	if closeErr := frameConn.Close(); closeErr != nil {
		fmt.Println(closeErr.Error())
	}
}
