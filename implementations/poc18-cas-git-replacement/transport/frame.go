package transport

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const maxFrameSize = 64 << 20

// Conn owns one length-prefixed TCP stream between two PromiseGrid agents.
type Conn struct {
	conn net.Conn
}

// Dial opens one bounded TCP connection to a peer agent.
func Dial(address string, timeout time.Duration) (*Conn, error) {
	netConn, dialErr := net.DialTimeout("tcp", address, timeout)
	if dialErr != nil {
		return nil, dialErr
	}
	return &Conn{conn: netConn}, nil
}

// Wrap returns a frame connection around an accepted TCP socket.
func Wrap(conn net.Conn) *Conn {
	return &Conn{conn: conn}
}

// WriteFrame writes one complete grid message or observer record.
//
// Intent: POC18 remediation requires real TCP bytes between agents while keeping
// framing separate from protocol semantics. Length framing is local transport
// machinery; the payload remains a promise-shaped grid message. Source: DI-koriz
func (conn *Conn) WriteFrame(frame []byte) error {
	if len(frame) == 0 {
		return fmt.Errorf("frame must not be empty")
	}
	if len(frame) > maxFrameSize {
		return fmt.Errorf("frame size %d exceeds maximum %d", len(frame), maxFrameSize)
	}
	var prefix [4]byte
	binary.BigEndian.PutUint32(prefix[:], uint32(len(frame)))
	if _, writeErr := conn.conn.Write(prefix[:]); writeErr != nil {
		return writeErr
	}
	_, writeErr := conn.conn.Write(frame)
	return writeErr
}

// ReadFrame reads one complete length-prefixed frame.
func (conn *Conn) ReadFrame() ([]byte, error) {
	var prefix [4]byte
	if _, readErr := io.ReadFull(conn.conn, prefix[:]); readErr != nil {
		return nil, readErr
	}
	length := binary.BigEndian.Uint32(prefix[:])
	if length == 0 || length > maxFrameSize {
		return nil, fmt.Errorf("invalid frame size %d", length)
	}
	frame := make([]byte, int(length))
	_, readErr := io.ReadFull(conn.conn, frame)
	return frame, readErr
}

// Close releases the TCP connection.
func (conn *Conn) Close() error {
	return conn.conn.Close()
}
