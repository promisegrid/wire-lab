package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc11-adaptive-trust-tcp/config"
	"promisegrid.dev/wire-lab/implementations/poc11-adaptive-trust-tcp/decision"
	"promisegrid.dev/wire-lab/implementations/poc11-adaptive-trust-tcp/economy"
	"promisegrid.dev/wire-lab/implementations/poc11-adaptive-trust-tcp/protocol"
	"promisegrid.dev/wire-lab/implementations/poc11-adaptive-trust-tcp/relationship"
	"promisegrid.dev/wire-lab/implementations/poc11-adaptive-trust-tcp/transport"
)

const sendTimeout = 5 * time.Second

// Node runs one autonomous POC11 agent process. A container may run several
// Node processes, but each process keeps its own local relationship ledger,
// listener, log, and live-LLM boundary.
// Intent: POC11 tests agent autonomy and adaptive TCP adjacency without
// collapsing co-located agents into a shared authority. Source: DI-hotos
type Node struct {
	Config      config.Config
	Agent       config.AgentConfig
	ProtocolCID protocol.ProtocolCID
	Decider     decision.Decider
	Monitor     decision.Monitor

	mu        sync.Mutex
	events    []decision.Event
	ledger    *relationship.Ledger
	evaluator economy.Evaluator
	listener  net.Listener
	logFile   *os.File
	budget    int
	capacity  int
}

// NewNode creates a node with a private trust ledger for every configured peer.
func NewNode(cfg config.Config, agent config.AgentConfig, decider decision.Decider, monitor decision.Monitor) *Node {
	peerNames := make([]string, 0, len(cfg.Agents)-1)
	for _, peer := range cfg.Agents {
		if peer.Name != agent.Name {
			peerNames = append(peerNames, peer.Name)
		}
	}
	return &Node{
		Config:      cfg,
		Agent:       agent,
		ProtocolCID: protocol.NewProtocolCID([]byte("poc11 adaptive trust tcp protocol v1")),
		Decider:     decider,
		Monitor:     monitor,
		ledger:      relationship.NewLedger(peerNames, agent.InitialPeers, cfg.StrongTrustThreshold, cfg.WeakTrustThreshold, cfg.TrustDecayPerRound),
		evaluator:   economy.Evaluator{},
		budget:      agent.Budget,
		capacity:    agent.Capacity,
	}
}

// Run executes bounded autonomous turns, writes a done marker, and waits for
// the observer-only monitor report.
func (node *Node) Run(ctx context.Context) error {
	if err := node.openLog(); err != nil {
		return err
	}
	defer node.closeLog()
	if err := node.maybeStartServer(ctx); err != nil {
		return err
	}
	defer node.closeListener()
	time.Sleep(node.Config.StartupDelay())
	for turnIndex := 0; turnIndex < node.Config.MaxTurns && turnIndex < node.Config.MaxAgentCalls; turnIndex++ {
		if err := node.runTurn(ctx, turnIndex); err != nil {
			node.record("decision_error", "broken", "", err.Error())
		}
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

func (node *Node) runTurn(ctx context.Context, turnIndex int) error {
	node.decayRelationships()
	observation := node.observation(turnIndex)
	if len(observation.DirectPeers) == 0 {
		node.record("local_non_commitment", "non_commitment", "", "no direct peer currently has enough local trust for a TCP promise")
		return nil
	}
	rawDecision, decideErr := node.Decider.Decide(ctx, observation)
	if decideErr != nil {
		return decideErr
	}
	validDecision, validateErr := decision.ValidatePromiseDecision(rawDecision, observation.DirectPeers)
	if validateErr != nil {
		node.observeOutcome(rawDecision.Target, relationship.OutcomeMalformed)
		node.record("decision_rejected", "malformed", rawDecision.Target, validateErr.Error())
		return nil
	}
	fields := decision.Fields(observation, validDecision)
	if economicsDecision := node.evaluateEconomics(validDecision.Target, fields); !economicsDecision.PromiseWorthMaking {
		node.observeOutcome(validDecision.Target, relationship.OutcomeNonCommitment)
		node.record("promise_withheld", "non_commitment", validDecision.Target, economicsDecision.Reason)
		return nil
	}
	if sendErr := node.send(validDecision.Target, fields); sendErr != nil {
		node.observeOutcome(validDecision.Target, relationship.OutcomeBroken)
		node.record("send_failed", "broken", validDecision.Target, sendErr.Error())
		return nil
	}
	node.spendLocalCapacity()
	node.record("promise_sent", "kept", validDecision.Target, validDecision.Promise)
	return nil
}

func (node *Node) maybeStartServer(ctx context.Context) error {
	listenPort, portFound := node.Config.ListenPortFor(node.Agent.Name)
	if !portFound || listenPort <= 0 {
		node.record("server_skipped", "kept", "", "no listener for local-only test run")
		return nil
	}
	address := net.JoinHostPort("", strconv.Itoa(listenPort))
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
	defer node.closeFrameConn(frameConn, "conn_close_failed", "")
	frameBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		node.record("frame_read_failed", "broken", "", readErr.Error())
		return
	}
	fields, exactHash, parseErr := node.parseEnvelope(frameBytes)
	if parseErr != nil {
		node.record("frame_parse_failed", "broken", "", parseErr.Error())
		return
	}
	fromAgent := fields["from"]
	if fields["act"] != decision.ActPromise {
		node.observeOutcome(fromAgent, relationship.OutcomeMalformed)
		node.record("message_rejected", "malformed", fromAgent, "message act is not promise")
		return
	}
	if !node.canAccept(fromAgent) {
		node.observeOutcome(fromAgent, relationship.OutcomeNonCommitment)
		node.record("message_not_promised", "non_commitment", fromAgent, "no current local promise to accept direct TCP exchange")
		node.writeAck(frameConn, fromAgent, "not_promised", "I promise to remember that I did not currently promise this direct exchange.")
		return
	}
	node.observeOutcome(fromAgent, relationship.OutcomeKept)
	node.record("message_received", "kept", fromAgent, "received valid signed promise exact_sha256="+exactHash)
	node.writeAck(frameConn, fromAgent, "kept", "I promise I received and recorded your signed promise message.")
}

func (node *Node) writeAck(frameConn transport.FrameConn, target, outcome, promiseText string) {
	ackFields := map[string]string{
		"act":     decision.ActPromise,
		"from":    node.Agent.Name,
		"to":      target,
		"outcome": outcome,
		"promise": promiseText,
		"reason":  "transport acknowledgement expressed as local promise content",
	}
	ack, ackErr := protocol.NewEnvelope(node.ProtocolCID, ackFields, node.Agent.Name)
	if ackErr != nil {
		node.record("ack_sign_failed", "broken", target, ackErr.Error())
		return
	}
	ackBytes, bytesErr := ack.Bytes()
	if bytesErr != nil {
		node.record("ack_bytes_failed", "broken", target, bytesErr.Error())
		return
	}
	if writeErr := frameConn.WriteFrame(ackBytes); writeErr != nil {
		node.record("ack_write_failed", "broken", target, writeErr.Error())
	}
}

func (node *Node) send(target string, fields map[string]string) error {
	if !node.canDial(target) {
		return fmt.Errorf("no local direct TCP promise to %s", target)
	}
	envelope, envelopeErr := protocol.NewEnvelope(node.ProtocolCID, fields, node.Agent.Name)
	if envelopeErr != nil {
		return envelopeErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	hostName, listenPort, endpointFound := node.Config.EndpointFor(target)
	if !endpointFound {
		return fmt.Errorf("no endpoint for target %s", target)
	}
	address := net.JoinHostPort(hostName, strconv.Itoa(listenPort))
	frameConn, dialErr := transport.DialFrameConn(address, sendTimeout)
	if dialErr != nil {
		return dialErr
	}
	defer node.closeFrameConn(frameConn, "send_close_failed", target)
	if writeErr := frameConn.WriteFrame(envelopeBytes); writeErr != nil {
		return writeErr
	}
	ackBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		return readErr
	}
	ackFields, _, parseErr := node.parseEnvelope(ackBytes)
	if parseErr != nil {
		return parseErr
	}
	if ackFields["outcome"] != "kept" {
		node.observeOutcome(target, relationship.OutcomeNonCommitment)
		return fmt.Errorf("ack outcome %q", ackFields["outcome"])
	}
	node.observeOutcome(target, relationship.OutcomeKept)
	return nil
}

func (node *Node) parseEnvelope(frameBytes []byte) (map[string]string, string, error) {
	envelope, parseErr := protocol.ParseEnvelope(frameBytes)
	if parseErr != nil {
		return nil, "", parseErr
	}
	if !envelope.ProtocolCID.Equal(node.ProtocolCID) {
		return nil, "", fmt.Errorf("unexpected pCID %s", envelope.ProtocolCID.String())
	}
	if verifyErr := protocol.VerifyEnvelope(envelope); verifyErr != nil {
		return nil, "", verifyErr
	}
	fields, fieldsErr := envelope.PayloadFields()
	if fieldsErr != nil {
		return nil, "", fieldsErr
	}
	if fields["from"] == "" {
		return nil, "", fmt.Errorf("payload from field is required")
	}
	return fields, protocol.HashExactBytes(frameBytes), nil
}

func (node *Node) observation(turnIndex int) decision.Observation {
	node.mu.Lock()
	defer node.mu.Unlock()
	recentEvents := make([]decision.Event, len(node.events))
	copy(recentEvents, node.events)
	if len(recentEvents) > 16 {
		recentEvents = recentEvents[len(recentEvents)-16:]
	}
	return decision.Observation{
		AgentName:      node.Agent.Name,
		Persona:        node.Agent.Persona,
		Motivation:     node.Agent.Motivation,
		Turn:           turnIndex,
		KnownPeers:     node.Config.AgentNames(),
		DirectPeers:    node.ledger.DirectPeers(),
		CandidatePeers: node.Config.CandidatePeersFor(node.Agent),
		LocalTrust:     node.ledger.Snapshot(),
		Budget:         node.budget,
		Capacity:       node.capacity,
		Adversarial:    node.Agent.Adversarial,
		RecentEvents:   recentEvents,
		RequiredAct:    decision.ActPromise,
	}
}

func (node *Node) evaluateEconomics(target string, fields map[string]string) economy.Decision {
	node.mu.Lock()
	defer node.mu.Unlock()
	reciprocalValue := 1
	if fields["field_promise_about"] == "reciprocal_economics" {
		reciprocalValue = 4
	}
	offer := economy.Offer{
		Promiser:        node.Agent.Name,
		Promisee:        target,
		Resource:        fields["field_promise_about"],
		PromisedValue:   2,
		ReciprocalValue: reciprocalValue,
		OpportunityCost: 1,
		Trust:           node.ledger.Trust(target),
		Budget:          node.budget,
		Capacity:        node.capacity,
	}
	return node.evaluator.Decide(offer)
}

func (node *Node) spendLocalCapacity() {
	node.mu.Lock()
	defer node.mu.Unlock()
	if node.budget > 0 {
		node.budget--
	}
	if node.capacity > 0 {
		node.capacity--
	}
}

func (node *Node) decayRelationships() {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.ledger.DecayRound()
}

func (node *Node) canDial(peerName string) bool {
	node.mu.Lock()
	defer node.mu.Unlock()
	return node.ledger.CanDial(peerName)
}

func (node *Node) canAccept(peerName string) bool {
	node.mu.Lock()
	defer node.mu.Unlock()
	return node.ledger.CanAccept(peerName)
}

func (node *Node) observeOutcome(peerName string, outcome relationship.Outcome) {
	if peerName == "" || peerName == node.Agent.Name {
		return
	}
	node.mu.Lock()
	defer node.mu.Unlock()
	node.ledger.ObserveOutcome(peerName, outcome)
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

func (node *Node) closeLog() {
	if node.logFile != nil {
		closeErr := node.logFile.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "close log: %v\n", closeErr)
		}
	}
}

func (node *Node) closeListener() {
	if node.listener != nil {
		closeErr := node.listener.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "close listener: %v\n", closeErr)
		}
	}
}

func (node *Node) closeFrameConn(frameConn transport.FrameConn, eventName, peerName string) {
	closeErr := frameConn.Close()
	if closeErr != nil {
		node.record(eventName, "broken", peerName, closeErr.Error())
	}
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
		if node.allDone() {
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

func (node *Node) allDone() bool {
	for _, agentName := range node.Config.AgentNames() {
		if _, err := os.Stat(filepath.Join(node.runDir(), agentName+".done")); err != nil {
			return false
		}
	}
	return true
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
	startIndex := 0
	for charIndex, charValue := range text {
		if charValue == '\n' {
			if charIndex > startIndex {
				lines = append(lines, text[startIndex:charIndex])
			}
			startIndex = charIndex + 1
		}
	}
	if startIndex < len(text) {
		lines = append(lines, text[startIndex:])
	}
	return lines
}
