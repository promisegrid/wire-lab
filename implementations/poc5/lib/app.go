package lib

import "time"

// RegisterReceiver opens a long-lived local kernel connection and promises to
// receive messages for one pCID.
func RegisterReceiver(kernelAddr string, receiveProtocolCID ProtocolCID, nodeName string, appName string, protocolCID ProtocolCID, text string) (FrameConn, error) {
	frameConn, dialErr := DialFrameConn(kernelAddr, 10*time.Second)
	if dialErr != nil {
		return FrameConn{}, dialErr
	}
	if err := WriteReceivePromise(frameConn, receiveProtocolCID, nodeName, appName, protocolCID, text); err != nil {
		if closeErr := frameConn.Close(); closeErr != nil {
			return FrameConn{}, closeErr
		}
		return FrameConn{}, err
	}
	return frameConn, nil
}

// SendEnvelopeToKernel submits one envelope to the local kernel.
func SendEnvelopeToKernel(kernelAddr string, envelope Envelope) ([]byte, error) {
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return nil, bytesErr
	}
	frameConn, dialErr := DialFrameConn(kernelAddr, 10*time.Second)
	if dialErr != nil {
		return nil, dialErr
	}
	writeErr := frameConn.WriteFrame(envelopeBytes)
	closeErr := frameConn.Close()
	if writeErr != nil {
		return envelopeBytes, writeErr
	}
	if closeErr != nil {
		return envelopeBytes, closeErr
	}
	return envelopeBytes, nil
}

// ReadEnvelope reads one framed envelope and decodes its fields.
func ReadEnvelope(frameConn FrameConn) (Envelope, string, map[string]string, []byte, error) {
	frameBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		return Envelope{}, "", nil, nil, readErr
	}
	envelope, parseErr := ParseEnvelope(frameBytes)
	if parseErr != nil {
		return Envelope{}, "", nil, frameBytes, parseErr
	}
	kind, fields, kindErr := EnvelopeKind(envelope)
	if kindErr != nil {
		return Envelope{}, "", nil, frameBytes, kindErr
	}
	return envelope, kind, fields, frameBytes, nil
}
