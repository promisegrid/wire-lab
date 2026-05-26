package storage

import (
	"fmt"

	"promisegrid.dev/wire-lab/implementations/poc4/kernel"
	"promisegrid.dev/wire-lab/implementations/poc4/lib"
	"promisegrid.dev/wire-lab/implementations/poc4/relay"
)

// StorageApp is an app agent that can promise local key-value storage or ask
// another storage app to keep and later return a value.
type StorageApp struct {
	NodeName   string
	AppName    string
	KernelAddr string
	Mode       string
	TargetNode string
	TargetApp  string
	Key        string
	Value      string
	values     map[string]string
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
	case "serve":
		return storageApp.runServe()
	default:
		return fmt.Errorf("unknown storage mode %q", storageApp.Mode)
	}
}

func (storageApp StorageApp) runClient() error {
	// Intent: Storage is exercised as a two-step promise sequence: first a
	// reciprocal confirmation promise for write, then a reciprocal result
	// promise for read. The client keeps local judgment over both outcomes.
	// Source: DI-bigub
	frameConn, registerErr := lib.RegisterReceiver(storageApp.KernelAddr, kernel.ReceiveProtocolCID(), storageApp.NodeName, storageApp.AppName, ProtocolCID(), "I promise to receive storage confirmations and read results for this bounded run.")
	if registerErr != nil {
		return registerErr
	}
	defer closeFrame(frameConn)
	storeRequest, storeErr := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":      "store_request_v1",
		"from":      storageApp.AppName,
		"from_node": storageApp.NodeName,
		"to":        storageApp.TargetApp,
		"key":       storageApp.Key,
		"value":     storageApp.Value,
	})
	if storeErr != nil {
		return storeErr
	}
	storeBytes, sendStoreErr := relay.SendInnerViaLocalRelay(storageApp.KernelAddr, storageApp.NodeName, storageApp.AppName, storageApp.TargetNode, storageApp.TargetApp, storeRequest, "")
	if sendStoreErr != nil {
		return sendStoreErr
	}
	if err := storageApp.expectStoreConfirm(frameConn, lib.HashExactBytes(storeBytes)); err != nil {
		return err
	}
	readRequest, readErr := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":      "read_request_v1",
		"from":      storageApp.AppName,
		"from_node": storageApp.NodeName,
		"to":        storageApp.TargetApp,
		"key":       storageApp.Key,
	})
	if readErr != nil {
		return readErr
	}
	readBytes, sendReadErr := relay.SendInnerViaLocalRelay(storageApp.KernelAddr, storageApp.NodeName, storageApp.AppName, storageApp.TargetNode, storageApp.TargetApp, readRequest, "")
	if sendReadErr != nil {
		return sendReadErr
	}
	return storageApp.expectReadResult(frameConn, lib.HashExactBytes(readBytes))
}

func (storageApp StorageApp) expectStoreConfirm(frameConn lib.FrameConn, requestHash string) error {
	_, kind, fields, _, readErr := lib.ReadEnvelope(frameConn)
	if readErr != nil {
		return readErr
	}
	if kind != "store_confirm_v1" || fields["request_hash"] != requestHash || fields["key"] != storageApp.Key {
		return fmt.Errorf("storage confirmation did not match store request")
	}
	fmt.Printf("%s judged storage confirmation kept from %s for key %s\n", storageApp.AppName, fields["from"], fields["key"])
	return nil
}

func (storageApp StorageApp) expectReadResult(frameConn lib.FrameConn, requestHash string) error {
	_, kind, fields, _, readErr := lib.ReadEnvelope(frameConn)
	if readErr != nil {
		return readErr
	}
	if kind != "read_result_v1" || fields["request_hash"] != requestHash || fields["key"] != storageApp.Key || fields["value"] != storageApp.Value {
		return fmt.Errorf("storage read result did not match local judgment")
	}
	fmt.Printf("%s judged storage read kept from %s: %s=%s\n", storageApp.AppName, fields["from"], fields["key"], fields["value"])
	return nil
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
