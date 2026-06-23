package artifact

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Event is the durable simulator/analyzer evidence row.
type Event struct {
	Time      string         `json:"time"`
	Type      string         `json:"type"`
	Actor     string         `json:"actor,omitempty"`
	Peer      string         `json:"peer,omitempty"`
	Direction string         `json:"direction,omitempty"`
	PCID      string         `json:"pcid,omitempty"`
	Hash      string         `json:"hash,omitempty"`
	Path      string         `json:"path,omitempty"`
	Transport string         `json:"transport,omitempty"`
	Outcome   string         `json:"outcome,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
}

// Writer stores observer-only run evidence under the approved temp root.
type Writer struct {
	runDir string
}

// NewWriter resets nothing; callers own cleanup before constructing it.
func NewWriter(runDir string) (*Writer, error) {
	for _, dir := range []string{"message-cas", "malformed"} {
		if err := os.MkdirAll(filepath.Join(runDir, dir), 0o755); err != nil {
			return nil, fmt.Errorf("create artifact dir %s: %w", dir, err)
		}
	}
	return &Writer{runDir: runDir}, nil
}

// RunDir returns the artifact root.
func (w *Writer) RunDir() string { return w.runDir }

// RecordMessage keeps exact message bytes in the run-scoped CAS.
func (w *Writer) RecordMessage(raw []byte) (string, string, error) {
	hashBytes := sha256.Sum256(raw)
	hash := hex.EncodeToString(hashBytes[:])
	rel := filepath.Join("message-cas", hash+".cbor")
	if err := os.WriteFile(filepath.Join(w.runDir, rel), raw, 0o644); err != nil {
		return "", "", fmt.Errorf("write message artifact: %w", err)
	}
	return hash, rel, nil
}

// RecordMalformed stores malformed radio bytes for review.
func (w *Writer) RecordMalformed(raw []byte, label string) (string, string, error) {
	hashBytes := sha256.Sum256(raw)
	hash := hex.EncodeToString(hashBytes[:])
	rel := filepath.Join("malformed", label+"-"+hash+".bin")
	if err := os.WriteFile(filepath.Join(w.runDir, rel), raw, 0o644); err != nil {
		return "", "", fmt.Errorf("write malformed artifact: %w", err)
	}
	return hash, rel, nil
}

// WriteEvent appends one JSON evidence row.
func (w *Writer) WriteEvent(event Event) error {
	if event.Time == "" {
		event.Time = time.Now().UTC().Format(time.RFC3339Nano)
	}
	path := filepath.Join(w.runDir, "events.jsonl")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open events: %w", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			// The append already returned to the caller; report close failure
			// through stderr rather than hiding it.
			fmt.Fprintf(os.Stderr, "close events file: %v\n", closeErr)
		}
	}()
	encoded, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("encode event: %w", err)
	}
	if _, err := f.Write(append(encoded, '\n')); err != nil {
		return fmt.Errorf("write event: %w", err)
	}
	return nil
}
