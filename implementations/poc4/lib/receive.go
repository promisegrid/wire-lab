package lib

// WriteReceivePromise sends a local app promise to receive messages for a pCID.
// The kernel receive pCID is supplied by the kernel package.
func WriteReceivePromise(frameConn FrameConn, receiveProtocolCID ProtocolCID, nodeName string, appName string, protocolCID ProtocolCID, text string) error {
	promiseEnvelope, promiseErr := NewEnvelope(receiveProtocolCID, map[string]string{
		"kind": "receive_promise_v1",
		"app":  appName,
		"node": nodeName,
		"pcid": protocolCID.String(),
		"text": text,
	})
	if promiseErr != nil {
		return promiseErr
	}
	promiseBytes, bytesErr := promiseEnvelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	return frameConn.WriteFrame(promiseBytes)
}
