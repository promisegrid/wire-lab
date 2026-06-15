package transport

import (
	"bytes"
	"encoding/binary"
	"net"
	"testing"
	"time"
)

func TestWriteAllRetriesPartialWrites(t *testing.T) {
	// Intent: Transport chaos should prove local retry behavior for short writes
	// without depending on a real flaky network. Source: DI-sunuf
	writer := &partialWriter{maxBytes: 2}
	if err := writeAll(writer, []byte("abcdef")); err != nil {
		t.Fatalf("write all with partial writer: %v", err)
	}
	if writer.String() != "abcdef" {
		t.Fatalf("partial writer captured %q, want abcdef", writer.String())
	}
}

func TestReadFrameAcceptsDelayedPartialBytes(t *testing.T) {
	// Intent: Delayed ACKs and byte-at-a-time TCP delivery should remain ordinary
	// frame parsing behavior, not malformed promise event records. Source: DI-sunuf
	serverConn, clientConn := net.Pipe()
	defer closeTestConn(t, serverConn, "server")
	defer closeTestConn(t, clientConn, "client")
	writeDone := make(chan error, 1)
	go func() {
		frameBytes := []byte("delayed-frame")
		header := make([]byte, 4)
		binary.BigEndian.PutUint32(header, uint32(len(frameBytes)))
		writeDone <- writeSlowBytes(serverConn, append(header, frameBytes...))
	}()
	gotFrame, readErr := NewFrameConn(clientConn).ReadFrame()
	if readErr != nil {
		t.Fatalf("read delayed frame: %v", readErr)
	}
	if string(gotFrame) != "delayed-frame" {
		t.Fatalf("frame = %q, want delayed-frame", string(gotFrame))
	}
	if writeErr := <-writeDone; writeErr != nil {
		t.Fatalf("write delayed frame: %v", writeErr)
	}
}

type partialWriter struct {
	bytes.Buffer
	maxBytes int
}

func (writer *partialWriter) Write(frameBytes []byte) (int, error) {
	if len(frameBytes) > writer.maxBytes {
		frameBytes = frameBytes[:writer.maxBytes]
	}
	return writer.Buffer.Write(frameBytes)
}

func writeSlowBytes(conn net.Conn, frameBytes []byte) error {
	for _, frameByte := range frameBytes {
		if _, err := conn.Write([]byte{frameByte}); err != nil {
			return err
		}
		time.Sleep(time.Millisecond)
	}
	return nil
}

func closeTestConn(t *testing.T, conn net.Conn, name string) {
	t.Helper()
	if err := conn.Close(); err != nil {
		t.Logf("close %s conn: %v", name, err)
	}
}
