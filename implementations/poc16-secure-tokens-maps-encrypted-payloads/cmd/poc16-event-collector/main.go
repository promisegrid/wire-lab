package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/config"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/decision"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/eventstream"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/transport"
)

// collector keeps the POC-only global observation surface outside agent
// runtimes. Intent: Agents and kernels must not coordinate through a shared
// filesystem; this collector only receives supervisor-forwarded stdout events
// and writes analyzer inputs after the run. Source: DI-dirat
type collector struct {
	cfg config.Config

	mu          sync.Mutex
	events      []decision.Event
	doneSources map[string]bool
	logFile     *os.File
	done        chan struct{}
	doneOnce    sync.Once
}

// messageDAGRecord is the operator-facing index row for one raw envelope file.
// Intent: The index names the exact `.cbor` artifact and optional parent hash
// without embedding raw bytes, so humans and follow-on tools can review the
// message DAG after a run without treating the event log as the message itself.
// Source: DI-tuhop
type messageDAGRecord struct {
	Source             string `json:"source"`
	Observer           string `json:"observer"`
	Direction          string `json:"direction"`
	Peer               string `json:"peer"`
	Protocol           string `json:"protocol"`
	ExactSHA256        string `json:"exact_sha256"`
	ParentExactSHA256  string `json:"parent_exact_sha256,omitempty"`
	ParentLinkLocation string `json:"parent_link_location,omitempty"`
	PromiseAbout       string `json:"promise_about,omitempty"`
	SourceEvent        string `json:"source_event,omitempty"`
	Path               string `json:"path"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "config.json", "POC16 config path")
	flag.Parse()
	cfg, loadErr := config.Load(*configPath)
	if loadErr != nil {
		return loadErr
	}
	listenAddress, listenErr := collectorListenAddress(cfg.EventCollectorAddress)
	if listenErr != nil {
		return listenErr
	}
	runCollector := &collector{
		cfg:         cfg,
		doneSources: make(map[string]bool),
		done:        make(chan struct{}),
	}
	if openErr := runCollector.openLog(); openErr != nil {
		return openErr
	}
	defer runCollector.closeLog()
	listener, listenErr := eventstream.Listen(listenAddress)
	if listenErr != nil {
		return listenErr
	}
	defer closeListener(listener)
	fmt.Printf("poc16-event-collector listening on %s expecting %d supervisors\n", listenAddress, len(cfg.Containers))
	go runCollector.acceptLoop(listener)
	select {
	case <-runCollector.done:
	case <-time.After(cfg.MonitorWaitTimeout()):
		return fmt.Errorf("timed out waiting for %d supervisor completions", len(cfg.Containers))
	}
	closeListener(listener)
	return runCollector.writeMonitorReport(context.Background())
}

func collectorListenAddress(configuredAddress string) (string, error) {
	_, port, splitErr := net.SplitHostPort(configuredAddress)
	if splitErr != nil {
		return "", splitErr
	}
	if _, atoiErr := strconv.Atoi(port); atoiErr != nil {
		return "", atoiErr
	}
	return ":" + port, nil
}

func (runCollector *collector) openLog() error {
	runDir := filepath.Join(runCollector.cfg.RunRoot, runCollector.cfg.RunID)
	if mkdirErr := os.MkdirAll(runDir, 0o755); mkdirErr != nil {
		return mkdirErr
	}
	logPath := filepath.Join(runDir, "events.jsonl")
	logFile, openErr := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if openErr != nil {
		return openErr
	}
	runCollector.logFile = logFile
	return nil
}

func (runCollector *collector) closeLog() {
	if runCollector.logFile == nil {
		return
	}
	if closeErr := runCollector.logFile.Close(); closeErr != nil {
		fmt.Fprintf(os.Stderr, "close collector log: %v\n", closeErr)
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
		go runCollector.handleConn(transport.NewFrameConn(conn))
	}
}

func (runCollector *collector) handleConn(frameConn transport.FrameConn) {
	defer closeFrameConn(frameConn)
	for {
		record, readErr := eventstream.ReadRecord(frameConn)
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
	case eventstream.KindSupervisorDone:
		runCollector.recordSupervisorDone(record.Source)
		return nil
	default:
		return fmt.Errorf("unknown event-stream record kind %q", record.Kind)
	}
}

func (runCollector *collector) recordEvent(event decision.Event) error {
	encoded, marshalErr := json.Marshal(event)
	if marshalErr != nil {
		return marshalErr
	}
	runCollector.mu.Lock()
	defer runCollector.mu.Unlock()
	runCollector.events = append(runCollector.events, event)
	if _, writeErr := runCollector.logFile.Write(append(encoded, '\n')); writeErr != nil {
		return writeErr
	}
	return nil
}

// recordMessageArtifact persists exact PromiseGrid envelope bytes for later
// review by operators and tools.
// Intent: This observer-only write path is deliberately outside app/kernel
// behavior: artifacts cannot affect routing, trust, or ACK outcomes, but a clean
// run can now prove that real messages are inspectable after containers exit.
// Source: DI-tuhop
func (runCollector *collector) recordMessageArtifact(source string, artifact eventstream.MessageArtifact) error {
	rawBytes, decodeErr := base64.StdEncoding.DecodeString(artifact.EnvelopeBytesBase64)
	if decodeErr != nil {
		return fmt.Errorf("decode message artifact from %s: %w", source, decodeErr)
	}
	if artifact.ExactSHA256 == "" {
		return fmt.Errorf("message artifact from %s missing exact_sha256", source)
	}
	actualHash := protocol.HashExactBytes(rawBytes)
	if actualHash != artifact.ExactSHA256 {
		return fmt.Errorf("message artifact from %s hash mismatch: record=%s actual=%s", source, artifact.ExactSHA256, actualHash)
	}
	runDir := filepath.Join(runCollector.cfg.RunRoot, runCollector.cfg.RunID)
	casDir := filepath.Join(runDir, "message-cas")
	if mkdirErr := os.MkdirAll(casDir, 0o755); mkdirErr != nil {
		return mkdirErr
	}
	artifactPath := filepath.Join(casDir, artifact.ExactSHA256+".cbor")
	indexPath := filepath.Join(runDir, "message-dag.jsonl")
	indexRecord := messageDAGRecord{
		Source:             source,
		Observer:           artifact.Observer,
		Direction:          artifact.Direction,
		Peer:               artifact.Peer,
		Protocol:           artifact.Protocol,
		ExactSHA256:        artifact.ExactSHA256,
		ParentExactSHA256:  artifact.ParentExactSHA256,
		ParentLinkLocation: artifact.ParentLinkLocation,
		PromiseAbout:       artifact.PromiseAbout,
		SourceEvent:        artifact.SourceEvent,
		Path:               filepath.ToSlash(filepath.Join("message-cas", artifact.ExactSHA256+".cbor")),
	}
	indexBytes, marshalErr := json.Marshal(indexRecord)
	if marshalErr != nil {
		return marshalErr
	}
	runCollector.mu.Lock()
	defer runCollector.mu.Unlock()
	if writeErr := os.WriteFile(artifactPath, rawBytes, 0o644); writeErr != nil {
		return writeErr
	}
	indexFile, openErr := os.OpenFile(indexPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if openErr != nil {
		return openErr
	}
	defer func() {
		if closeErr := indexFile.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "close message DAG index: %v\n", closeErr)
		}
	}()
	if _, writeErr := indexFile.Write(append(indexBytes, '\n')); writeErr != nil {
		return writeErr
	}
	return nil
}

func (runCollector *collector) recordSupervisorDone(source string) {
	runCollector.mu.Lock()
	runCollector.doneSources[source] = true
	doneCount := len(runCollector.doneSources)
	expectedCount := len(runCollector.cfg.Containers)
	runCollector.mu.Unlock()
	fmt.Printf("poc16-event-collector supervisor done source=%s count=%d/%d\n", source, doneCount, expectedCount)
	if doneCount >= expectedCount {
		runCollector.doneOnce.Do(func() {
			close(runCollector.done)
		})
	}
}

func (runCollector *collector) writeMonitorReport(ctx context.Context) error {
	runCollector.mu.Lock()
	events := append([]decision.Event{}, runCollector.events...)
	runCollector.mu.Unlock()
	liveClient := decision.NewLiveClient(
		runCollector.cfg.ProviderBaseURL,
		runCollector.cfg.APIKeyEnv,
		runCollector.cfg.AgentModel,
		runCollector.cfg.MonitorModel,
		runCollector.cfg.ReasoningEffort,
		runCollector.cfg.ServiceTier,
		runCollector.cfg.RequestTimeout(),
	)
	report, evaluateErr := liveClient.Evaluate(ctx, events)
	if evaluateErr != nil {
		return evaluateErr
	}
	reportBytes, marshalErr := json.MarshalIndent(report, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	reportPath := filepath.Join(runCollector.cfg.RunRoot, runCollector.cfg.RunID, "monitor-report.json")
	return os.WriteFile(reportPath, append(reportBytes, '\n'), 0o644)
}

func closeListener(listener net.Listener) {
	if listener == nil {
		return
	}
	if closeErr := listener.Close(); closeErr != nil && !strings.Contains(closeErr.Error(), "use of closed network connection") {
		fmt.Fprintf(os.Stderr, "close collector listener: %v\n", closeErr)
	}
}

func closeFrameConn(frameConn transport.FrameConn) {
	if closeErr := frameConn.Close(); closeErr != nil {
		fmt.Fprintf(os.Stderr, "close collector connection: %v\n", closeErr)
	}
}
