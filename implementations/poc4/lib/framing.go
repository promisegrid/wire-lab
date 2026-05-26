package lib

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

// FrameConn sends exact message bytes over a boring length-framed stream.
//
// Intent: Keep TCP framing out of the Promise Theory experiment while preserving
// exact-byte relay and app evidence. Source: DI-bigub.
type FrameConn struct {
	conn net.Conn
}

// NewFrameConn wraps a net.Conn.
func NewFrameConn(conn net.Conn) FrameConn {
	return FrameConn{conn: conn}
}

// DialFrameConn dials a framed TCP connection with a timeout.
func DialFrameConn(address string, timeout time.Duration) (FrameConn, error) {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return FrameConn{}, err
	}
	return NewFrameConn(conn), nil
}

// WriteFrame writes one length-prefixed frame.
func (frameConn FrameConn) WriteFrame(frameBytes []byte) error {
	if len(frameBytes) == 0 || len(frameBytes) > 16*1024*1024 {
		return fmt.Errorf("invalid frame length: %d", len(frameBytes))
	}
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(frameBytes)))
	if _, err := frameConn.conn.Write(header); err != nil {
		return err
	}
	_, err := frameConn.conn.Write(frameBytes)
	return err
}

// ReadFrame reads one length-prefixed frame.
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

// Close closes the underlying stream.
func (frameConn FrameConn) Close() error {
	return frameConn.conn.Close()
}
