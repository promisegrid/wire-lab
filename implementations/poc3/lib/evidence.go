package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"
)

// EvidenceRecord is a local process observation. It is not a global fact and is
// not application-level keep/break judgment.
type EvidenceRecord struct {
	TimeUTC       string `json:"time_utc"`
	Node          string `json:"node"`
	Event         string `json:"event"`
	Boundary      string `json:"boundary"`
	Outcome       string `json:"outcome"`
	ProtocolCID   string `json:"pcid"`
	ExactSHA256   string `json:"exact_sha256,omitempty"`
	Detail        string `json:"detail"`
	PromiseSource string `json:"promise_source"`
}

// EvidenceLog serializes local kernel observations to stdout and optional
// writers. Apps make their own local promise judgments separately.
//
// Intent: Make kernel receive/deliver/refuse evidence visible without letting
// the kernel judge app-level promise status. Source: DI-horak.
type EvidenceLog struct {
	nodeName string
	writers  []io.Writer
	mutex    sync.Mutex
}

// NewEvidenceLog creates an evidence log that always writes to stdout.
func NewEvidenceLog(nodeName string, extraWriters ...io.Writer) *EvidenceLog {
	writers := []io.Writer{os.Stdout}
	writers = append(writers, extraWriters...)
	return &EvidenceLog{nodeName: nodeName, writers: writers}
}

// Record writes one local evidence record.
func (evidenceLog *EvidenceLog) Record(event string, boundary string, outcome string, protocolCID ProtocolCID, exactBytes []byte, detail string) error {
	record := EvidenceRecord{
		TimeUTC:       time.Now().UTC().Format(time.RFC3339Nano),
		Node:          evidenceLog.nodeName,
		Event:         event,
		Boundary:      boundary,
		Outcome:       outcome,
		ProtocolCID:   protocolCID.String(),
		ExactSHA256:   HashExactBytes(exactBytes),
		Detail:        detail,
		PromiseSource: "local kernel evidence; not app promise judgment",
	}
	encodedRecord, marshalErr := json.Marshal(record)
	if marshalErr != nil {
		return marshalErr
	}
	evidenceLog.mutex.Lock()
	defer evidenceLog.mutex.Unlock()
	for _, writer := range evidenceLog.writers {
		if _, err := writer.Write(append(encodedRecord, '\n')); err != nil {
			return err
		}
	}
	return nil
}

// HashExactBytes returns the sha256 hex digest for evidence correlation.
func HashExactBytes(exactBytes []byte) string {
	if len(exactBytes) == 0 {
		return ""
	}
	digest := sha256.Sum256(exactBytes)
	return hex.EncodeToString(digest[:])
}
