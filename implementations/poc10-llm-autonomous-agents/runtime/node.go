package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc10-llm-autonomous-agents/config"
	"promisegrid.dev/wire-lab/implementations/poc10-llm-autonomous-agents/decision"
	"promisegrid.dev/wire-lab/implementations/poc10-llm-autonomous-agents/protocol"
	"promisegrid.dev/wire-lab/implementations/poc10-llm-autonomous-agents/transport"
)

// Node runs one autonomous POC10 agent inside one container.
// Intent: The kernel/runtime owns transport, signing, config bounds, and local
// evidence while the LLM owns only the next local decision. Source: DI-pijan
type Node struct {
	Config      config.Config
	Agent       config.AgentConfig
	ProtocolCID protocol.ProtocolCID
	Decider     decision.Decider
	Monitor     decision.Monitor

	mu       sync.Mutex
	events   []decision.Event
	trust    map[string]int
	listener net.Listener
	logFile  *os.File
}

// NewNode creates a node with neutral local trust for every configured peer.
func NewNode(cfg config.Config, agent config.AgentConfig, decider decision.Decider, monitor decision.Monitor) *Node {
	trust := make(map[string]int, len(cfg.Agents))
	for _, peer := range cfg.Agents {
		if peer.Name != agent.Name {
			trust[peer.Name] = 0
		}
	}
	return &Node{
		Config:      cfg,
		Agent:       agent,
		ProtocolCID: protocol.NewProtocolCID([]byte("poc10 llm autonomous agents protocol v1")),
		Decider:     decider,
		Monitor:     monitor,
		trust:       trust,
	}
}

// Run executes bounded autonomous turns, writes a done marker, and waits for
// the observer-only monitor report.
func (node *Node) Run(ctx context.Context) error {
	if err := node.openLog(); err != nil {
		return err
	}
	defer func() {
		if node.logFile != nil {
			closeErr := node.logFile.Close()
			if closeErr != nil {
				fmt.Fprintf(os.Stderr, "close log: %v\n", closeErr)
			}
		}
	}()
	if node.Config.ListenPort > 0 {
		if err := node.startServer(ctx); err != nil {
			return err
		}
	} else {
		// Intent: Tests can exercise local decision/evidence behavior in sandboxes
		// that prohibit TCP listeners; validated runtime configs still require a
		// positive listen port. Source: DI-pijan
		node.record("server_skipped", "kept", "", "no listener for local-only test run")
	}
	defer func() {
		if node.listener != nil {
			closeErr := node.listener.Close()
			if closeErr != nil {
				fmt.Fprintf(os.Stderr, "close listener: %v\n", closeErr)
			}
		}
	}()
	time.Sleep(node.Config.StartupDelay())
	calls := 0
	for turn := 0; turn < node.Config.MaxTurns && calls < node.Config.MaxAgentCalls; turn++ {
		if err := node.runTurn(ctx, turn); err != nil {
			node.record("decision_error", "broken", "", err.Error())
		}
		calls++
		time.Sleep(node.Config.TurnDelay())
	}
	if err := node.writeDoneMarker(); err != nil {
		return err
	}
	if node.Agent.Name == node.Config.MonitorNode {
		if err := node.runMonitor(ctx); err != nil {
			node.record("monitor_error", "broken", "", err.Error())
		}
	}
	return node.waitForMonitor(ctx)
}

func (node *Node) runTurn(ctx context.Context, turn int) error {
	observation := node.observation(turn)
	rawDecision, decideErr := node.Decider.Decide(ctx, observation)
	if decideErr != nil {
		return decideErr
	}
	validDecision, validateErr := decision.Validate(node.Agent.Profile, rawDecision, node.Agent.Neighbors)
	if validateErr != nil {
		node.record("decision_rejected", "refused", rawDecision.Target, validateErr.Error())
		return nil
	}
	if validDecision.Target == "" {
		node.record("local_observation", "kept", "", validDecision.Reason)
		return nil
	}
	fields := decision.Fields(observation, validDecision)
	if sendErr := node.send(validDecision.Target, fields); sendErr != nil {
		node.adjustTrust(validDecision.Target, -1)
		node.record("send_failed", "broken", validDecision.Target, sendErr.Error())
		return nil
	}
	node.record("decision_sent", "kept", validDecision.Target, validDecision.Action+" / "+validDecision.Promise)
	return nil
}

func (node *Node) startServer(ctx context.Context) error {
	address := fmt.Sprintf(":%d", node.Config.ListenPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	node.listener = listener
	go node.acceptLoop(ctx, listener)
	return nil
}

func (node *Node) acceptLoop(ctx context.Context, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				node.record("accept_failed", "broken", "", err.Error())
				return
			}
		}
		go node.handleConn(conn)
	}
}

func (node *Node) handleConn(conn net.Conn) {
	frameConn := transport.NewFrameConn(conn)
	defer func() {
		closeErr := frameConn.Close()
		if closeErr != nil {
			node.record("conn_close_failed", "broken", "", closeErr.Error())
		}
	}()
	frameBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		node.record("frame_read_failed", "broken", "", readErr.Error())
		return
	}
	fields, parseErr := node.parseAndRecord(frameBytes)
	if parseErr != nil {
		node.record("frame_parse_failed", "broken", "", parseErr.Error())
		return
	}
	ackFields := map[string]string{
		"kind":    "transport_ack",
		"from":    node.Agent.Name,
		"to":      fields["from"],
		"outcome": "kept",
		"detail":  "received valid signed grid envelope",
	}
	ack, ackErr := protocol.NewEnvelope(node.ProtocolCID, ackFields, node.Agent.Name)
	if ackErr != nil {
		node.record("ack_sign_failed", "broken", fields["from"], ackErr.Error())
		return
	}
	ackBytes, bytesErr := ack.Bytes()
	if bytesErr != nil {
		node.record("ack_bytes_failed", "broken", fields["from"], bytesErr.Error())
		return
	}
	if writeErr := frameConn.WriteFrame(ackBytes); writeErr != nil {
		node.record("ack_write_failed", "broken", fields["from"], writeErr.Error())
	}
}

func (node *Node) send(target string, fields map[string]string) error {
	envelope, envelopeErr := protocol.NewEnvelope(node.ProtocolCID, fields, node.Agent.Name)
	if envelopeErr != nil {
		return envelopeErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	frameConn, dialErr := transport.DialFrameConn(fmt.Sprintf("%s:%d", target, node.Config.ListenPort), 5*time.Second)
	if dialErr != nil {
		return dialErr
	}
	defer func() {
		closeErr := frameConn.Close()
		if closeErr != nil {
			node.record("send_close_failed", "broken", target, closeErr.Error())
		}
	}()
	if writeErr := frameConn.WriteFrame(envelopeBytes); writeErr != nil {
		return writeErr
	}
	ackBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		return readErr
	}
	ackFields, parseErr := node.parseAndRecord(ackBytes)
	if parseErr != nil {
		return parseErr
	}
	if ackFields["outcome"] != "kept" {
		return fmt.Errorf("ack outcome %q", ackFields["outcome"])
	}
	return nil
}

func (node *Node) parseAndRecord(frameBytes []byte) (map[string]string, error) {
	envelope, parseErr := protocol.ParseEnvelope(frameBytes)
	if parseErr != nil {
		return nil, parseErr
	}
	if !envelope.ProtocolCID.Equal(node.ProtocolCID) {
		return nil, fmt.Errorf("unexpected pCID %s", envelope.ProtocolCID.String())
	}
	if verifyErr := protocol.VerifyEnvelope(envelope); verifyErr != nil {
		return nil, verifyErr
	}
	fields, fieldsErr := envelope.PayloadFields()
	if fieldsErr != nil {
		return nil, fieldsErr
	}
	from := fields["from"]
	kind := fields["kind"]
	detail := kind + " from " + from + " exact_sha256=" + protocol.HashExactBytes(frameBytes)
	node.adjustTrust(from, trustDeltaForKind(kind))
	node.record("message_received", "kept", from, detail)
	return fields, nil
}

func (node *Node) observation(turn int) decision.Observation {
	node.mu.Lock()
	defer node.mu.Unlock()
	trust := make(map[string]int, len(node.trust))
	for key, value := range node.trust {
		trust[key] = value
	}
	recent := make([]decision.Event, len(node.events))
	copy(recent, node.events)
	if len(recent) > 12 {
		recent = recent[len(recent)-12:]
	}
	return decision.Observation{
		AgentName:        node.Agent.Name,
		Profile:          string(node.Agent.Profile),
		Persona:          node.Agent.Persona,
		Motivation:       node.Agent.Motivation,
		Turn:             turn,
		KnownPeers:       node.Config.AgentNames(),
		NeighborPeers:    node.Agent.Neighbors,
		LocalTrust:       trust,
		RecentEvents:     recent,
		AvailableActions: decision.AvailableActions(node.Agent.Profile),
	}
}

func (node *Node) adjustTrust(peer string, delta int) {
	if peer == "" || peer == node.Agent.Name {
		return
	}
	node.mu.Lock()
	defer node.mu.Unlock()
	node.trust[peer] += delta
}

func trustDeltaForKind(kind string) int {
	switch kind {
	case "transport_ack", "offer_promise", "need_advertisement", "introduction_promise", "route_promise":
		return 1
	case "refusal":
		return 0
	case "freeform_intent":
		return -1
	default:
		return 0
	}
}

func (node *Node) record(eventName, outcome, peer, detail string) {
	event := decision.Event{
		Observer: node.Agent.Name,
		Event:    eventName,
		Outcome:  outcome,
		Peer:     peer,
		Detail:   detail,
	}
	node.mu.Lock()
	node.events = append(node.events, event)
	node.mu.Unlock()
	encoded, err := json.Marshal(event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal event: %v\n", err)
		return
	}
	fmt.Println(string(encoded))
	if node.logFile != nil {
		if _, writeErr := node.logFile.Write(append(encoded, '\n')); writeErr != nil {
			fmt.Fprintf(os.Stderr, "write event: %v\n", writeErr)
		}
	}
}

func (node *Node) openLog() error {
	runDir := node.runDir()
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		return err
	}
	logPath := filepath.Join(runDir, node.Agent.Name+".jsonl")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	node.logFile = logFile
	return nil
}

func (node *Node) runDir() string {
	return filepath.Join(node.Config.RunRoot, node.Config.RunID)
}

func (node *Node) writeDoneMarker() error {
	donePath := filepath.Join(node.runDir(), node.Agent.Name+".done")
	if err := os.WriteFile(donePath, []byte("done\n"), 0o644); err != nil {
		return err
	}
	node.record("node_done", "kept", "", "wrote local done marker")
	return nil
}

func (node *Node) waitForMonitor(ctx context.Context) error {
	monitorDone := filepath.Join(node.runDir(), "monitor.done")
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	timeout := time.After(90 * time.Second)
	for {
		if _, err := os.Stat(monitorDone); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timed out waiting for monitor.done")
		case <-ticker.C:
		}
	}
}

func (node *Node) runMonitor(ctx context.Context) error {
	if err := node.waitForAllDone(ctx); err != nil {
		return err
	}
	events, readErr := node.readAllEvents()
	if readErr != nil {
		return readErr
	}
	report, evaluateErr := node.Monitor.Evaluate(ctx, events)
	if evaluateErr != nil {
		return evaluateErr
	}
	reportBytes, marshalErr := json.MarshalIndent(report, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	reportPath := filepath.Join(node.runDir(), "monitor-report.json")
	if err := os.WriteFile(reportPath, append(reportBytes, '\n'), 0o644); err != nil {
		return err
	}
	donePath := filepath.Join(node.runDir(), "monitor.done")
	if err := os.WriteFile(donePath, []byte("done\n"), 0o644); err != nil {
		return err
	}
	node.record("monitor_done", "kept", "", report.Summary)
	return nil
}

func (node *Node) waitForAllDone(ctx context.Context) error {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	timeout := time.After(90 * time.Second)
	for {
		allDone := true
		for _, agentName := range node.Config.AgentNames() {
			if _, err := os.Stat(filepath.Join(node.runDir(), agentName+".done")); err != nil {
				allDone = false
				break
			}
		}
		if allDone {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timed out waiting for all node done markers")
		case <-ticker.C:
		}
	}
}

func (node *Node) readAllEvents() ([]decision.Event, error) {
	var events []decision.Event
	for _, agentName := range node.Config.AgentNames() {
		logPath := filepath.Join(node.runDir(), agentName+".jsonl")
		logBytes, readErr := os.ReadFile(logPath)
		if readErr != nil {
			return nil, readErr
		}
		lines := splitLines(string(logBytes))
		for _, line := range lines {
			var event decision.Event
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				return nil, err
			}
			events = append(events, event)
		}
	}
	return events, nil
}

func splitLines(text string) []string {
	var lines []string
	start := 0
	for index, char := range text {
		if char == '\n' {
			if index > start {
				lines = append(lines, text[start:index])
			}
			start = index + 1
		}
	}
	if start < len(text) {
		lines = append(lines, text[start:])
	}
	return lines
}
