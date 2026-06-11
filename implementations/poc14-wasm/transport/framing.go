package transport

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

// FrameConn carries one signed POC14 envelope per length-prefixed TCP frame.
// Intent: Keep transport as simple TCP plumbing while semantic content remains
// in the pCID-selected CBOR grid envelope. Source: DI-timah
type FrameConn struct {
	conn net.Conn
}

// NewFrameConn wraps an accepted TCP connection.
func NewFrameConn(conn net.Conn) FrameConn {
	return FrameConn{conn: conn}
}

// DialFrameConn dials a peer and returns a length-framed stream.
func DialFrameConn(address string, timeout time.Duration) (FrameConn, error) {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return FrameConn{}, err
	}
	return NewFrameConn(conn), nil
}

// WriteFrame writes exactly one bounded frame to the TCP stream.
func (frameConn FrameConn) WriteFrame(frameBytes []byte) error {
	if len(frameBytes) == 0 || len(frameBytes) > 16*1024*1024 {
		return fmt.Errorf("invalid frame length: %d", len(frameBytes))
	}
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(frameBytes)))
	if err := writeAll(frameConn.conn, header); err != nil {
		return err
	}
	return writeAll(frameConn.conn, frameBytes)
}

// ReadFrame reads exactly one bounded frame from the TCP stream.
func (frameConn FrameConn) ReadFrame() ([]byte, error) {
	header := make([]byte, 4)
	if _, err := io.ReadFull(frameConn.conn, header); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(header)
	if length == 0 || length > 16*1024*1024 {
		return nil, fmt.Errorf("invalid frame length: %d", length)
	}
	frameBytes := make([]byte, int(length))
	if _, err := io.ReadFull(frameConn.conn, frameBytes); err != nil {
		return nil, err
	}
	return frameBytes, nil
}

// Close closes the underlying TCP connection.
func (frameConn FrameConn) Close() error {
	return frameConn.conn.Close()
}

// writeAll handles short writes from TCP implementations and chaos tests.
// Intent: A partial kernel/app frame write should be retried locally rather than
// silently truncating a signed CBOR grid envelope. Source: DI-sunuf
func writeAll(writer io.Writer, frameBytes []byte) error {
	for len(frameBytes) > 0 {
		written, writeErr := writer.Write(frameBytes)
		if written > 0 {
			frameBytes = frameBytes[written:]
		}
		if writeErr != nil {
			return writeErr
		}
		if written == 0 {
			return io.ErrShortWrite
		}
	}
	return nil
}
