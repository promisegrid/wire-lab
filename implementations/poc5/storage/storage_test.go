package storage

import (
	"net"
	"testing"

	"promisegrid.dev/wire-lab/implementations/poc5/lib"
)

func TestStorageClientJudgesStoreConfirmation(t *testing.T) {
	storageApp := StorageApp{AppName: "carol-storage-client", Key: "poc5-key", Value: "poc5-value"}
	envelope, err := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":         "store_confirm_v1",
		"from":         "dave-storage-app",
		"request_hash": "store-hash",
		"key":          "poc5-key",
		"value":        "poc5-value",
	})
	if err != nil {
		t.Fatalf("NewEnvelope returned error: %v", err)
	}
	frameConn, done := testFrameWithEnvelope(t, envelope)
	if err := storageApp.expectStoreConfirm(frameConn, "store-hash", "poc5-key"); err != nil {
		t.Fatalf("expectStoreConfirm returned error: %v", err)
	}
	if err := frameConn.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	if err := <-done; err != nil {
		t.Fatalf("writer returned error: %v", err)
	}
}

func TestStorageClientJudgesReadResult(t *testing.T) {
	storageApp := StorageApp{AppName: "carol-storage-client", Key: "poc5-key", Value: "poc5-value"}
	envelope, err := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":         "read_result_v1",
		"from":         "dave-storage-app",
		"request_hash": "read-hash",
		"key":          "poc5-key",
		"value":        "poc5-value",
	})
	if err != nil {
		t.Fatalf("NewEnvelope returned error: %v", err)
	}
	frameConn, done := testFrameWithEnvelope(t, envelope)
	kept, err := storageApp.expectReadResult(frameConn, "read-hash", "poc5-key", "poc5-value")
	if err != nil {
		t.Fatalf("expectReadResult returned error: %v", err)
	}
	if !kept {
		t.Fatalf("expectReadResult did not judge matching read result kept")
	}
	if err := frameConn.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	if err := <-done; err != nil {
		t.Fatalf("writer returned error: %v", err)
	}
}

func TestStorageClientJudgesBrokenReadResultLocally(t *testing.T) {
	storageApp := StorageApp{NodeName: "alice", AppName: "alice-trust-client", Key: "poc5-key", Value: "poc5-value"}
	envelope, err := lib.NewEnvelope(ProtocolCID(), map[string]string{
		"kind":         "read_result_v1",
		"from":         "bob-storage-app",
		"request_hash": "read-hash",
		"key":          "poc5-key",
		"value":        "wrong-value",
	})
	if err != nil {
		t.Fatalf("NewEnvelope returned error: %v", err)
	}
	frameConn, done := testFrameWithEnvelope(t, envelope)
	kept, err := storageApp.expectReadResult(frameConn, "read-hash", "poc5-key", "poc5-value")
	if err != nil {
		t.Fatalf("expectReadResult returned error: %v", err)
	}
	if kept {
		t.Fatalf("expectReadResult judged mismatched read result kept")
	}
	if err := frameConn.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	if err := <-done; err != nil {
		t.Fatalf("writer returned error: %v", err)
	}
}

func testFrameWithEnvelope(t *testing.T, envelope lib.Envelope) (lib.FrameConn, <-chan error) {
	t.Helper()
	clientConn, serverConn := net.Pipe()
	done := make(chan error, 1)
	go func() {
		frameConn := lib.NewFrameConn(serverConn)
		envelopeBytes, bytesErr := envelope.Bytes()
		if bytesErr != nil {
			done <- bytesErr
			return
		}
		writeErr := frameConn.WriteFrame(envelopeBytes)
		closeErr := frameConn.Close()
		if writeErr != nil {
			done <- writeErr
			return
		}
		done <- closeErr
	}()
	return lib.NewFrameConn(clientConn), done
}
