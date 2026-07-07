// Command poc-event-collector receives observer-only POC18 run artifacts.
//
// Intent: POC18 agents must not coordinate through a shared filesystem. The
// collector receives one-way event streams after messages have already crossed
// agent TCP connections, then persists exact artifacts for post-run review.
// Source: DI-koriz
package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/eventstream"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/transport"
)

type collector struct {
	runDir       string
	expectedDone int
	eventsFile   *os.File
	mu           sync.Mutex
	doneSources  map[string]bool
	done         chan struct{}
	doneOnce     sync.Once
	active       sync.WaitGroup
}

type messageDAGRecord struct {
	Source      string `json:"source"`
	Observer    string `json:"observer"`
	Direction   string `json:"direction"`
	Peer        string `json:"peer"`
	Protocol    string `json:"protocol"`
	PromiseKind string `json:"promise_kind"`
	ExactCID    string `json:"exact_cid"`
	Path        string `json:"path"`
}

type carDAGRecord struct {
	Source     string   `json:"source"`
	Observer   string   `json:"observer"`
	Direction  string   `json:"direction"`
	Peer       string   `json:"peer"`
	MessageCID string   `json:"message_cid"`
	CARCID     string   `json:"car_cid"`
	BlockCIDs  []string `json:"block_cids"`
	Path       string   `json:"path"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "poc-event-collector: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	listenAddress := flag.String("listen", ":9200", "collector TCP listen address")
	runRoot := flag.String("run-root", "/run/poc18", "observer runtime root")
	runID := flag.String("run-id", "poc18-demo", "observer run id")
	expectedDone := flag.Int("expected", 7, "number of agent completion records expected")
	timeout := flag.Duration("timeout", 30*time.Second, "maximum wait for agent completions")
	flag.Parse()
	if *expectedDone <= 0 {
		return fmt.Errorf("expected agent count must be positive")
	}
	runCollector, collectorErr := newCollector(filepath.Join(*runRoot, *runID), *expectedDone)
	if collectorErr != nil {
		return collectorErr
	}
	defer runCollector.closeEventsFile()
	listener, listenErr := eventstream.Listen(*listenAddress)
	if listenErr != nil {
		return listenErr
	}
	defer closeListener(listener)
	fmt.Printf("poc-event-collector listening on %s run=%s expected=%d\n", *listenAddress, runCollector.runDir, *expectedDone)
	go runCollector.acceptLoop(listener)
	select {
	case <-runCollector.done:
	case <-time.After(*timeout):
		return fmt.Errorf("timed out waiting for %d agent completions", *expectedDone)
	}
	closeListener(listener)
	if waitErr := runCollector.waitForConnections(2 * time.Second); waitErr != nil {
		return waitErr
	}
	return nil
}

func newCollector(runDir string, expectedDone int) (*collector, error) {
	for _, dir := range []string{
		runDir,
		filepath.Join(runDir, "message-cas"),
		filepath.Join(runDir, "car-cas"),
	} {
		if mkdirErr := os.MkdirAll(dir, 0o755); mkdirErr != nil {
			return nil, mkdirErr
		}
	}
	eventsPath := filepath.Join(runDir, "events.jsonl")
	eventsFile, openErr := os.OpenFile(eventsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if openErr != nil {
		return nil, openErr
	}
	return &collector{
		runDir:       runDir,
		expectedDone: expectedDone,
		eventsFile:   eventsFile,
		doneSources:  map[string]bool{},
		done:         make(chan struct{}),
	}, nil
}

func (runCollector *collector) closeEventsFile() {
	if runCollector.eventsFile == nil {
		return
	}
	if closeErr := runCollector.eventsFile.Close(); closeErr != nil {
		fmt.Fprintf(os.Stderr, "close events file: %v\n", closeErr)
	}
}

func (runCollector *collector) acceptLoop(listener net.Listener) {
	for {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			if errors.Is(acceptErr, net.ErrClosed) {
				return
			}
			fmt.Fprintf(os.Stderr, "collector accept: %v\n", acceptErr)
			return
		}
		runCollector.active.Add(1)
		go func() {
			defer runCollector.active.Done()
			runCollector.handleConn(transport.Wrap(conn))
		}()
	}
}

func (runCollector *collector) waitForConnections(timeout time.Duration) error {
	// Intent: The final agent completion closes the run, but an accepted
	// event-stream handler may still be flushing exact artifact frames. Drain the
	// collector-side handlers without giving agents any feedback channel. Source:
	// DI-koriz
	done := make(chan struct{})
	go func() {
		runCollector.active.Wait()
		close(done)
	}()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-done:
		return nil
	case <-timer.C:
		return fmt.Errorf("collector connection drain exceeded %s", timeout)
	}
}

func (runCollector *collector) handleConn(conn *transport.Conn) {
	defer closeConn(conn)
	for {
		record, readErr := eventstream.ReadRecord(conn)
		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				return
			}
			fmt.Fprintf(os.Stderr, "collector read record: %v\n", readErr)
			return
		}
		if handleErr := runCollector.handleRecord(record); handleErr != nil {
			fmt.Fprintf(os.Stderr, "collector handle record: %v\n", handleErr)
			return
		}
	}
}

func (runCollector *collector) handleRecord(record eventstream.Record) error {
	switch record.Kind {
	case eventstream.KindEvent:
		if record.Event == nil {
			return fmt.Errorf("event record from %s is missing event", record.Source)
		}
		return runCollector.recordEvent(*record.Event)
	case eventstream.KindMessageArtifact:
		if record.MessageArtifact == nil {
			return fmt.Errorf("message artifact record from %s is missing artifact", record.Source)
		}
		return runCollector.recordMessageArtifact(record.Source, *record.MessageArtifact)
	case eventstream.KindCARArtifact:
		if record.CARArtifact == nil {
			return fmt.Errorf("CAR artifact record from %s is missing artifact", record.Source)
		}
		return runCollector.recordCARArtifact(record.Source, *record.CARArtifact)
	case eventstream.KindSupervisorDone:
		runCollector.recordDone(record.Source)
		return nil
	default:
		return fmt.Errorf("unknown event-stream record kind %q", record.Kind)
	}
}

func (runCollector *collector) recordEvent(event eventstream.Event) error {
	encoded, marshalErr := json.Marshal(event)
	if marshalErr != nil {
		return marshalErr
	}
	runCollector.mu.Lock()
	defer runCollector.mu.Unlock()
	if _, writeErr := runCollector.eventsFile.Write(append(encoded, '\n')); writeErr != nil {
		return writeErr
	}
	return nil
}

func (runCollector *collector) recordMessageArtifact(source string, artifact eventstream.MessageArtifact) error {
	rawBytes, decodeErr := base64.StdEncoding.DecodeString(artifact.EnvelopeBytesBase64)
	if decodeErr != nil {
		return decodeErr
	}
	actualCID := store.CIDText(store.CIDForBytes(rawBytes))
	if actualCID != artifact.ExactCID {
		return fmt.Errorf("message artifact CID mismatch source=%s record=%s actual=%s", source, artifact.ExactCID, actualCID)
	}
	path := filepath.Join(runCollector.runDir, "message-cas", artifact.ExactCID+".cbor")
	record := messageDAGRecord{
		Source:      source,
		Observer:    artifact.Observer,
		Direction:   artifact.Direction,
		Peer:        artifact.Peer,
		Protocol:    artifact.Protocol,
		PromiseKind: artifact.PromiseKind,
		ExactCID:    artifact.ExactCID,
		Path:        filepath.ToSlash(filepath.Join("message-cas", artifact.ExactCID+".cbor")),
	}
	return runCollector.writeArtifactAndIndex(path, rawBytes, filepath.Join(runCollector.runDir, "message-dag.jsonl"), record)
}

func (runCollector *collector) recordCARArtifact(source string, artifact eventstream.CARArtifact) error {
	rawBytes, decodeErr := base64.StdEncoding.DecodeString(artifact.CARBytesBase64)
	if decodeErr != nil {
		return decodeErr
	}
	actualCID := store.CIDText(store.CIDForBytes(rawBytes))
	if actualCID != artifact.CARCID {
		return fmt.Errorf("CAR artifact CID mismatch source=%s record=%s actual=%s", source, artifact.CARCID, actualCID)
	}
	path := filepath.Join(runCollector.runDir, "car-cas", artifact.CARCID+".car")
	record := carDAGRecord{
		Source:     source,
		Observer:   artifact.Observer,
		Direction:  artifact.Direction,
		Peer:       artifact.Peer,
		MessageCID: artifact.MessageCID,
		CARCID:     artifact.CARCID,
		BlockCIDs:  append([]string(nil), artifact.BlockCIDs...),
		Path:       filepath.ToSlash(filepath.Join("car-cas", artifact.CARCID+".car")),
	}
	return runCollector.writeArtifactAndIndex(path, rawBytes, filepath.Join(runCollector.runDir, "car-dag.jsonl"), record)
}

func (runCollector *collector) writeArtifactAndIndex(path string, content []byte, indexPath string, record any) error {
	indexBytes, marshalErr := json.Marshal(record)
	if marshalErr != nil {
		return marshalErr
	}
	runCollector.mu.Lock()
	defer runCollector.mu.Unlock()
	if writeErr := os.WriteFile(path, content, 0o644); writeErr != nil {
		return writeErr
	}
	indexFile, openErr := os.OpenFile(indexPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if openErr != nil {
		return openErr
	}
	defer closeFile(indexFile, "index")
	if _, writeErr := indexFile.Write(append(indexBytes, '\n')); writeErr != nil {
		return writeErr
	}
	return nil
}

func (runCollector *collector) recordDone(source string) {
	runCollector.mu.Lock()
	runCollector.doneSources[source] = true
	doneCount := len(runCollector.doneSources)
	expectedCount := runCollector.expectedDone
	runCollector.mu.Unlock()
	fmt.Printf("poc-event-collector done source=%s count=%d/%d\n", source, doneCount, expectedCount)
	if doneCount >= expectedCount {
		runCollector.doneOnce.Do(func() {
			close(runCollector.done)
		})
	}
}

func closeListener(listener net.Listener) {
	if listener == nil {
		return
	}
	if closeErr := listener.Close(); closeErr != nil && !errors.Is(closeErr, net.ErrClosed) {
		fmt.Fprintf(os.Stderr, "close listener: %v\n", closeErr)
	}
}

func closeConn(conn *transport.Conn) {
	if closeErr := conn.Close(); closeErr != nil {
		fmt.Fprintf(os.Stderr, "close collector connection: %v\n", closeErr)
	}
}

func closeFile(file *os.File, label string) {
	if closeErr := file.Close(); closeErr != nil {
		fmt.Fprintf(os.Stderr, "close %s file: %v\n", label, closeErr)
	}
}
