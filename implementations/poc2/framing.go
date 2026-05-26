package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const maxFrameBytes = 1 << 20

// FrameConn sends exact envelope bytes over TCP with a simple length prefix.
// The length prefix is transport framing only; the PromiseGrid boundary remains
// the enclosed CBOR grid message.
//
// Intent: Keep app/kernel and kernel/kernel transport boring so the POC tests
// pCID-selected promise messages, not stream parsing tricks. Source: DI-tijat.
type FrameConn struct {
	conn net.Conn
}

func NewFrameConn(conn net.Conn) FrameConn {
	return FrameConn{conn: conn}
}

func (frameConn FrameConn) WriteFrame(frameBytes []byte) error {
	if len(frameBytes) > maxFrameBytes {
		return fmt.Errorf("frame too large: %d", len(frameBytes))
	}
	var lengthBytes [4]byte
	binary.BigEndian.PutUint32(lengthBytes[:], uint32(len(frameBytes)))
	if _, err := frameConn.conn.Write(lengthBytes[:]); err != nil {
		return err
	}
	_, writeErr := frameConn.conn.Write(frameBytes)
	return writeErr
}

func (frameConn FrameConn) ReadFrame() ([]byte, error) {
	var lengthBytes [4]byte
	if _, err := io.ReadFull(frameConn.conn, lengthBytes[:]); err != nil {
		return nil, err
	}
	frameLength := binary.BigEndian.Uint32(lengthBytes[:])
	if frameLength > maxFrameBytes {
		return nil, fmt.Errorf("frame too large: %d", frameLength)
	}
	frameBytes := make([]byte, frameLength)
	if _, err := io.ReadFull(frameConn.conn, frameBytes); err != nil {
		return nil, err
	}
	return frameBytes, nil
}

func dialFrameConn(address string, timeout time.Duration) (FrameConn, error) {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for time.Now().Before(deadline) {
		conn, dialErr := net.DialTimeout("tcp", address, 250*time.Millisecond)
		if dialErr == nil {
			return NewFrameConn(conn), nil
		}
		lastErr = dialErr
		time.Sleep(100 * time.Millisecond)
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("deadline reached before dial attempt")
	}
	return FrameConn{}, fmt.Errorf("dial %s: %w", address, lastErr)
}

func (frameConn FrameConn) Close() error {
	return frameConn.conn.Close()
}
