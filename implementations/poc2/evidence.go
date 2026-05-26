package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// EvidenceRecord is a local observation. It is deliberately not an authority
// fact; another agent may ignore it or weigh it differently.
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

// EvidenceLog serializes local promise observations to stdout and, optionally,
// a JSONL file.
//
// Intent: Make make/break/refusal evidence visible without pretending there is a
// global trust ledger. Source: DI-ratij; DI-tijat.
type EvidenceLog struct {
	nodeName string
	writers  []io.Writer
	mutex    sync.Mutex
}

func NewEvidenceLog(nodeName string, evidencePath string) (*EvidenceLog, func() error, error) {
	writers := []io.Writer{os.Stdout}
	cleanup := func() error { return nil }
	if evidencePath != "" {
		file, createErr := os.Create(evidencePath)
		if createErr != nil {
			return nil, nil, createErr
		}
		writers = append(writers, file)
		cleanup = file.Close
	}
	return &EvidenceLog{nodeName: nodeName, writers: writers}, cleanup, nil
}

func (evidenceLog *EvidenceLog) Record(event string, boundary string, outcome string, protocolCID ProtocolCID, exactBytes []byte, detail string) error {
	record := EvidenceRecord{
		TimeUTC:       time.Now().UTC().Format(time.RFC3339Nano),
		Node:          evidenceLog.nodeName,
		Event:         event,
		Boundary:      boundary,
		Outcome:       outcome,
		ProtocolCID:   protocolCID.String(),
		ExactSHA256:   hashExactBytes(exactBytes),
		Detail:        detail,
		PromiseSource: "local observation; not global authority",
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

func hashExactBytes(exactBytes []byte) string {
	if len(exactBytes) == 0 {
		return ""
	}
	digest := sha256.Sum256(exactBytes)
	return hex.EncodeToString(digest[:])
}

func observationEnvelope(protocolCID ProtocolCID, from string, to string, outcome string, text string, evidenceBytes []byte) (Envelope, error) {
	return NewEnvelope(protocolCID, map[string]string{
		"kind":          "observation_v1",
		"from":          from,
		"to":            to,
		"outcome":       outcome,
		"text":          text,
		"evidence_hash": hashExactBytes(evidenceBytes),
	})
}

func envelopeKind(envelope Envelope) (string, map[string]string, error) {
	fields, payloadErr := envelope.PayloadFields()
	if payloadErr != nil {
		return "", nil, payloadErr
	}
	kind := fields["kind"]
	if kind == "" {
		return "", nil, fmt.Errorf("payload missing kind")
	}
	return kind, fields, nil
}
