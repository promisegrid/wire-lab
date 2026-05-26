package storage

import (
	"fmt"

	"promisegrid.dev/wire-lab/implementations/poc5/kernel"
	"promisegrid.dev/wire-lab/implementations/poc5/lib"
	"promisegrid.dev/wire-lab/implementations/poc5/relay"
)

// StorageApp is an app agent that can promise local key-value storage or ask
// another storage app to keep and later return a value.
type StorageApp struct {
	NodeName     string
	AppName      string
	KernelAddr   string
	Mode         string
	TargetNode   string
	TargetApp    string
	FallbackNode string
	FallbackApp  string
	Key          string
	Value        string
	breakRead    bool
	values       map[string]string
}

// Run executes one bounded storage role.
func (storageApp StorageApp) Run() error {
	if storageApp.AppName == "" {
		storageApp.AppName = storageApp.NodeName + "-storage-app"
	}
	if storageApp.values == nil {
		storageApp.values = make(map[string]string)
	}
	switch storageApp.Mode {
	case "client":
		return storageApp.runClient()
	case "trust-client":
		return storageApp.runTrustClient()
	case "serve":
		return storageApp.runServe()
	case "serve-break":
		storageApp.breakRead = true
		return storageApp.runServe()
	default:
		return fmt.Errorf("unknown storage mode %q", storageApp.Mode)
	}
}

func (storageApp StorageApp) runTrustClient() error {
	// Intent: POC5 makes selective sending an Alice-local trust decision. Alice
	// first sends only low-sensitivity data to Bob, observes Bob break a promise,
	// lowers Alice's own trust in Bob, and then declines to send sensitive data
	// to Bob without asking any authority to decide for her. Source: DI-rarim.
	frameConn, registerErr := lib.RegisterReceiver(storageApp.KernelAddr, kernel.ReceiveProtocolCID(), storageApp.NodeName, storageApp.AppName, ProtocolCID(), "I promise to receive storage confirmations and read results so I can update my own local trust.")
	if registerErr != nil {
		return registerErr
	}
	defer closeFrame(frameConn)
	evidenceLog := lib.NewEvidenceLog(storageApp.NodeName, storageApp.AppName)
	localTrust := map[string]int{
		storageApp.TargetNode:   1,
		storageApp.FallbackNode: 2,
	}
	trustThreshold := 2
	targetKept, targetErr := storageApp.storageRoundTrip(frameConn, storageApp.TargetNode, storageApp.TargetApp, "poc5-probe-key", "poc5-probe-value")
	if targetErr != nil {
		return targetErr
	}
	if !targetKept {
		localTrust[storageApp.TargetNode]--
		if recordErr := evidenceLog.Record("trust_decreased", "app/app", "kept", ProtocolCID(), nil, storageApp.TargetNode+" local trust decreased after Alice observed a broken storage promise"); recordErr != nil {
			return recordErr
		}
	}
	if localTrust[storageApp.TargetNode] < trustThreshold {
		if recordErr := evidenceLog.Record("selective_send_declined", "app/app", "kept", ProtocolCID(), nil, "Alice declined to send sensitive storage data to "+storageApp.TargetNode+" because Alice's local trust threshold was not met"); recordErr != nil {
			return recordErr
		}
	} else {
		return fmt.Errorf("expected %s to fall below Alice's local trust threshold", storageApp.TargetNode)
	}
	if localTrust[storageApp.FallbackNode] < trustThreshold {
		return fmt.Errorf("fallback peer %s did not meet Alice's local trust threshold", storageApp.FallbackNode)
	}
	if recordErr := evidenceLog.Record("selective_send_chosen", "app/app", "kept", ProtocolCID(), nil, "Alice chose "+storageApp.FallbackNode+" for sensitive storage after local trust comparison"); recordErr != nil {
		return recordErr
	}
	fallbackKept, fallbackErr := storageApp.storageRoundTrip(frameConn, storageApp.FallbackNode, storageApp.FallbackApp, storageApp.Key, storageApp.Value)
	if fallbackErr != nil {
		return fallbackErr
	}
	if !fallbackKept {
		return fmt.Errorf("fallback storage promise from %s was not kept", storageApp.FallbackNode)
	}
	if recordErr := evidenceLog.Record("promise_kept", "app/app", "kept", ProtocolCID(), nil, storageApp.FallbackNode+" kept Alice's sensitive storage promise"); recordErr != nil {
		return recordErr
	}
	return nil
}

func (storageApp StorageApp) runClient() error {
	// Intent: Storage is exercised as a two-step promise sequence: first a
	// reciprocal confirmation promise for write, then a reciprocal result
	// promise for read. The client keeps local judgment over both outcomes.
	// Source: DI-rarim
	frameConn, registerErr := lib.RegisterReceiver(storageApp.KernelAddr, kernel.ReceiveProtocolCID(), storageApp.NodeName, storageApp.AppName, ProtocolCID(), "I promise to receive storage confirmations and read results for this bounded run.")
	if registerErr != nil {
		return registerErr
	}
	defer closeFrame(frameConn)
	kept, roundTripErr := storageApp.storageRoundTrip(frameConn, storageApp.TargetNode, storageApp.TargetApp, storageApp.Key, storageApp.Value)
	if roundTripErr != nil {
		return roundTripErr
	}
	if !kept {
		return fmt.Errorf("storage read result did not match local judgment")
	}
	return nil
}

func (storageApp StorageApp) storageRoundTrip(frameConn lib.FrameConn, targetNode string, targetApp string, key string, value string) (bool, error) {
	storeRequest, storeErr := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":      "store_request_v1",
		"from":      storageApp.AppName,
		"from_node": storageApp.NodeName,
		"to":        targetApp,
		"key":       key,
		"value":     value,
	})
	if storeErr != nil {
		return false, storeErr
	}
	storeBytes, sendStoreErr := relay.SendInnerViaLocalRelay(storageApp.KernelAddr, storageApp.NodeName, storageApp.AppName, targetNode, targetApp, storeRequest, "")
	if sendStoreErr != nil {
		return false, sendStoreErr
	}
	if err := storageApp.expectStoreConfirm(frameConn, lib.HashExactBytes(storeBytes), key); err != nil {
		return false, err
	}
	readRequest, readErr := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":      "read_request_v1",
		"from":      storageApp.AppName,
		"from_node": storageApp.NodeName,
		"to":        targetApp,
		"key":       key,
	})
	if readErr != nil {
		return false, readErr
	}
	readBytes, sendReadErr := relay.SendInnerViaLocalRelay(storageApp.KernelAddr, storageApp.NodeName, storageApp.AppName, targetNode, targetApp, readRequest, "")
	if sendReadErr != nil {
		return false, sendReadErr
	}
	return storageApp.expectReadResult(frameConn, lib.HashExactBytes(readBytes), key, value)
}

func (storageApp StorageApp) expectStoreConfirm(frameConn lib.FrameConn, requestHash string, key string) error {
	_, kind, fields, _, readErr := lib.ReadEnvelope(frameConn)
	if readErr != nil {
		return readErr
	}
	if kind != "store_confirm_v1" || fields["request_hash"] != requestHash || fields["key"] != key {
		return fmt.Errorf("storage confirmation did not match store request")
	}
	fmt.Printf("%s judged storage confirmation kept from %s for key %s\n", storageApp.AppName, fields["from"], fields["key"])
	return nil
}

func (storageApp StorageApp) expectReadResult(frameConn lib.FrameConn, requestHash string, key string, value string) (bool, error) {
	_, kind, fields, resultBytes, readErr := lib.ReadEnvelope(frameConn)
	if readErr != nil {
		return false, readErr
	}
	if kind != "read_result_v1" || fields["request_hash"] != requestHash || fields["key"] != key {
		return false, fmt.Errorf("storage read result did not match request")
	}
	if fields["value"] != value {
		evidenceLog := lib.NewEvidenceLog(storageApp.NodeName, storageApp.AppName)
		if recordErr := evidenceLog.Record("promise_broken", "app/app", "broken", ProtocolCID(), resultBytes, fields["from"]+" returned value "+fields["value"]+" when Alice expected "+value); recordErr != nil {
			return false, recordErr
		}
		fmt.Printf("%s judged storage read broken from %s: expected %s=%s got %s=%s\n", storageApp.AppName, fields["from"], key, value, fields["key"], fields["value"])
		return false, nil
	}
	fmt.Printf("%s judged storage read kept from %s: %s=%s\n", storageApp.AppName, fields["from"], fields["key"], fields["value"])
	return true, nil
}

func (storageApp StorageApp) runServe() error {
	frameConn, registerErr := lib.RegisterReceiver(storageApp.KernelAddr, kernel.ReceiveProtocolCID(), storageApp.NodeName, storageApp.AppName, ProtocolCID(), "I promise to receive one store request and one read request for this bounded run.")
	if registerErr != nil {
		return registerErr
	}
	defer closeFrame(frameConn)
	for messageIndex := 0; messageIndex < 2; messageIndex++ {
		if err := storageApp.handleOneStorageMessage(frameConn); err != nil {
			return err
		}
	}
	return nil
}

func (storageApp StorageApp) handleOneStorageMessage(frameConn lib.FrameConn) error {
	envelope, kind, fields, requestBytes, readErr := lib.ReadEnvelope(frameConn)
	if readErr != nil {
		return readErr
	}
	if !envelope.ProtocolCID.Equal(ProtocolCID()) {
		return fmt.Errorf("storage server got unsupported pCID")
	}
	switch kind {
	case "store_request_v1":
		storageApp.values[fields["key"]] = fields["value"]
		return storageApp.reply(fields, requestBytes, "store_confirm_v1", fields["value"])
	case "read_request_v1":
		value, ok := storageApp.values[fields["key"]]
		if !ok {
			return fmt.Errorf("storage key %s was not promised by this app", fields["key"])
		}
		if storageApp.breakRead {
			return storageApp.reply(fields, requestBytes, "read_result_v1", value+"-broken")
		}
		return storageApp.reply(fields, requestBytes, "read_result_v1", value)
	default:
		return fmt.Errorf("storage server got unexpected kind %s", kind)
	}
}

func (storageApp StorageApp) reply(fields map[string]string, requestBytes []byte, kind string, value string) error {
	responseEnvelope, responseErr := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":         kind,
		"from":         storageApp.AppName,
		"from_node":    storageApp.NodeName,
		"to":           fields["from"],
		"to_node":      fields["from_node"],
		"request_hash": lib.HashExactBytes(requestBytes),
		"key":          fields["key"],
		"value":        value,
	})
	if responseErr != nil {
		return responseErr
	}
	if _, sendErr := relay.SendInnerViaLocalRelay(storageApp.KernelAddr, storageApp.NodeName, storageApp.AppName, fields["from_node"], fields["from"], responseEnvelope, lib.HashExactBytes(requestBytes)); sendErr != nil {
		return sendErr
	}
	fmt.Printf("%s judged %s from %s kept for key %s\n", storageApp.AppName, kind, fields["from"], fields["key"])
	return nil
}

func closeFrame(frameConn lib.FrameConn) {
	if closeErr := frameConn.Close(); closeErr != nil {
		fmt.Println(closeErr.Error())
	}
}
