package eventstream

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/transport"
)

const (
	KindEvent           = "event"
	KindMessageArtifact = "message_artifact"
	KindCARArtifact     = "car_artifact"
	KindSupervisorDone  = "supervisor_done"
)

var stdoutMu sync.Mutex

// Event is an observer-only runtime fact. It is not a PromiseGrid promise and is
// never consumed by agents.
type Event struct {
	Observer string `json:"observer"`
	Event    string `json:"event"`
	Outcome  string `json:"outcome"`
	Peer     string `json:"peer,omitempty"`
	Detail   string `json:"detail,omitempty"`
}

// MessageArtifact carries one exact grid envelope copy for post-run review.
type MessageArtifact struct {
	Observer            string `json:"observer"`
	Direction           string `json:"direction"`
	Peer                string `json:"peer"`
	Protocol            string `json:"protocol"`
	ExactCID            string `json:"exact_cid"`
	ParentCID           string `json:"parent_cid,omitempty"`
	PromiseKind         string `json:"promise_kind,omitempty"`
	EnvelopeBytesBase64 string `json:"envelope_bytes_b64"`
}

// CARArtifact carries exact CARv1 bytes copied from an object_bytes promise.
type CARArtifact struct {
	Observer       string   `json:"observer"`
	Direction      string   `json:"direction"`
	Peer           string   `json:"peer"`
	MessageCID     string   `json:"message_cid"`
	CARCID         string   `json:"car_cid"`
	BlockCIDs      []string `json:"block_cids"`
	CARBytesBase64 string   `json:"car_bytes_b64"`
}

// Record is the one-way observer stream consumed by the collector.
type Record struct {
	Kind            string           `json:"kind"`
	Source          string           `json:"source"`
	Event           *Event           `json:"event,omitempty"`
	MessageArtifact *MessageArtifact `json:"message_artifact,omitempty"`
	CARArtifact     *CARArtifact     `json:"car_artifact,omitempty"`
	Detail          string           `json:"detail,omitempty"`
}

// Client sends observer-only records to the collector while also allowing the
// caller to write the same JSON line to stdout.
type Client struct {
	source string
	conn   *transport.Conn
	mu     sync.Mutex
}

// Dial connects to the observer collector.
func Dial(source string, address string, timeout time.Duration) (*Client, error) {
	deadline := time.Now().Add(timeout)
	var lastErr error
	for {
		conn, dialErr := transport.Dial(address, 750*time.Millisecond)
		if dialErr == nil {
			return &Client{source: source, conn: conn}, nil
		}
		lastErr = dialErr
		if time.Now().After(deadline) {
			return nil, fmt.Errorf("connect collector %s: %w", address, lastErr)
		}
		time.Sleep(250 * time.Millisecond)
	}
}

// Emit writes one record to stdout and forwards it to the collector.
//
// Intent: POC18 follows POC16's observer pattern: exact message copies are
// emitted out-of-band for analysis, but collector state is not visible to agents
// and cannot affect routing, trust, or promise outcomes. Source: DI-koriz
func (client *Client) Emit(record Record) error {
	if record.Source == "" {
		record.Source = client.source
	}
	encoded, marshalErr := json.Marshal(record)
	if marshalErr != nil {
		return marshalErr
	}
	stdoutMu.Lock()
	_, stdoutErr := fmt.Fprintln(os.Stdout, string(encoded))
	stdoutMu.Unlock()
	if stdoutErr != nil {
		return stdoutErr
	}
	client.mu.Lock()
	defer client.mu.Unlock()
	return client.conn.WriteFrame(encoded)
}

// Close closes the observer collector connection.
func (client *Client) Close() error {
	return client.conn.Close()
}

// Listen opens the collector listener.
func Listen(address string) (net.Listener, error) {
	return net.Listen("tcp", address)
}

// ReadRecord reads one observer record from a framed collector connection.
func ReadRecord(conn *transport.Conn) (Record, error) {
	frame, readErr := conn.ReadFrame()
	if readErr != nil {
		return Record{}, readErr
	}
	var record Record
	if unmarshalErr := json.Unmarshal(frame, &record); unmarshalErr != nil {
		return Record{}, unmarshalErr
	}
	return record, nil
}

// MessageArtifactFor builds an observer copy of one exact grid message.
func MessageArtifactFor(observer, direction, peer, protocol, promiseKind string, envelopeBytes []byte) MessageArtifact {
	return MessageArtifact{
		Observer:            observer,
		Direction:           direction,
		Peer:                peer,
		Protocol:            protocol,
		PromiseKind:         promiseKind,
		ExactCID:            store.CIDText(store.CIDForBytes(envelopeBytes)),
		EnvelopeBytesBase64: base64.StdEncoding.EncodeToString(envelopeBytes),
	}
}

// CARArtifactFor builds an observer copy of exact CAR bytes.
func CARArtifactFor(observer, direction, peer, messageCID string, carBytes []byte, blockCIDs []string) CARArtifact {
	return CARArtifact{
		Observer:       observer,
		Direction:      direction,
		Peer:           peer,
		MessageCID:     messageCID,
		CARCID:         store.CIDText(store.CIDForBytes(carBytes)),
		BlockCIDs:      append([]string(nil), blockCIDs...),
		CARBytesBase64: base64.StdEncoding.EncodeToString(carBytes),
	}
}
