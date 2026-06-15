package eventstream

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/decision"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/transport"
)

const (
	KindEvent           = "event"
	KindMessageArtifact = "message_artifact"
	KindSupervisorDone  = "supervisor_done"
)

const dialRetryDelay = 250 * time.Millisecond

// Record is the observer-only event stream between a container supervisor and
// the POC collector. Intent: This stream is not PromiseGrid traffic and cannot
// affect agent behavior; it exists only so the development harness can analyze
// stdout events without giving agents a shared run volume. Source: DI-dirat
type Record struct {
	Kind            string           `json:"kind"`
	Source          string           `json:"source"`
	Event           *decision.Event  `json:"event,omitempty"`
	MessageArtifact *MessageArtifact `json:"message_artifact,omitempty"`
	Detail          string           `json:"detail,omitempty"`
}

// MessageArtifact carries exact PromiseGrid envelope bytes to the observer-only
// collector for run-scoped operator review.
// Intent: Agents must not coordinate through the Docker observer volume, but the
// POC harness still needs intact message bytes after the run. Supervisors forward
// these base64 transport records one way to the collector, which writes binary
// `.cbor` artifacts and a DAG-style index without influencing agents. Source:
// DI-tuhop
type MessageArtifact struct {
	Observer            string `json:"observer"`
	Direction           string `json:"direction"`
	Peer                string `json:"peer"`
	Protocol            string `json:"protocol"`
	ExactSHA256         string `json:"exact_sha256"`
	ParentExactSHA256   string `json:"parent_exact_sha256,omitempty"`
	PromiseAbout        string `json:"promise_about,omitempty"`
	SourceEvent         string `json:"source_event,omitempty"`
	EnvelopeBytesBase64 string `json:"envelope_bytes_b64"`
}

// Client is one supervisor's long-lived connection to the observer-only
// collector.
type Client struct {
	mu        sync.Mutex
	frameConn transport.FrameConn
}

// Dial connects to the collector with bounded retry because Compose startup
// order does not prove that the listening socket is ready.
func Dial(ctx context.Context, address string, timeout time.Duration) (*Client, error) {
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()
	var lastErr error
	for {
		frameConn, dialErr := transport.DialFrameConn(address, dialRetryDelay)
		if dialErr == nil {
			return &Client{frameConn: frameConn}, nil
		}
		lastErr = dialErr
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-deadline.C:
			return nil, fmt.Errorf("connect event collector %s: %w", address, lastErr)
		case <-time.After(dialRetryDelay):
		}
	}
}

// Send writes one framed JSON record to the collector.
func (client *Client) Send(record Record) error {
	recordBytes, marshalErr := json.Marshal(record)
	if marshalErr != nil {
		return marshalErr
	}
	client.mu.Lock()
	defer client.mu.Unlock()
	return client.frameConn.WriteFrame(recordBytes)
}

// Close closes the supervisor's collector connection.
func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	return client.frameConn.Close()
}

// ReadRecord reads one framed JSON record from an accepted collector connection.
func ReadRecord(frameConn transport.FrameConn) (Record, error) {
	frameBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		if readErr == io.EOF {
			return Record{}, io.EOF
		}
		return Record{}, readErr
	}
	var record Record
	if unmarshalErr := json.Unmarshal(frameBytes, &record); unmarshalErr != nil {
		return Record{}, unmarshalErr
	}
	return record, nil
}

// Listen opens the TCP listener used by the observer-only event collector.
func Listen(address string) (net.Listener, error) {
	return net.Listen("tcp", address)
}
