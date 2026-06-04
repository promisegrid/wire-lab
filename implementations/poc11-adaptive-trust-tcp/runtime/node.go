package runtime

import (
	"context"
	"encoding/json"
	"errors"
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
const shutdownDrainTimeout = 750 * time.Millisecond

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

	activeHandlers sync.WaitGroup
	stopping       bool
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
	if err := node.loadRelationshipState(); err != nil {
		return err
	}
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
	// Intent: Stop accepting new TCP frames before the local done marker is
	// written so `node_done` does not race with late receive receipts.
	// Source: DI-duhub
	node.closeListener()
	node.drainInflight(ctx)
	if err := node.saveRelationshipState(); err != nil {
		return err
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
		repairedDecision, repaired, repairErr := decision.RepairPromiseDecision(rawDecision, observation, validateErr)
		if repairErr != nil {
			node.observeOutcome(rawDecision.Target, relationship.OutcomeMalformed)
			node.record("decision_rejected", "malformed", rawDecision.Target, validateErr.Error())
			return nil
		}
		if repaired {
			node.record("decision_repaired", "kept", repairedDecision.Target, repairErrDetail(validateErr))
		}
		validDecision = repairedDecision
	}
	fields := decision.Fields(observation, validDecision)
	if resourceErr := node.checkLocalResourcePromise(fields); resourceErr != nil {
		node.observeOutcome(validDecision.Target, relationship.OutcomeBroken)
		node.applyBrokenPromiseCost(validDecision.Target, fields, resourceErr.Error())
		node.record("resource_promise_broken", "broken", validDecision.Target, resourceErr.Error())
		return nil
	}
	if economicsDecision := node.evaluateEconomics(validDecision.Target, fields); !economicsDecision.PromiseWorthMaking {
		node.observeOutcome(validDecision.Target, relationship.OutcomeNonCommitment)
		node.record("promise_withheld", "non_commitment", validDecision.Target, economicsDecision.Reason)
		return nil
	}
	if sendErr := node.send(validDecision.Target, fields); sendErr != nil {
		node.observeOutcome(validDecision.Target, relationship.OutcomeBroken)
		node.applyBrokenPromiseCost(validDecision.Target, fields, sendErr.Error())
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
			if node.isStopping() || errors.Is(err, net.ErrClosed) {
				node.record("listener_closed", "kept", "", "listener closed during normal shutdown")
				return
			}
			select {
			case <-ctx.Done():
				return
			default:
				node.record("accept_failed", "broken", "", err.Error())
				return
			}
		}
		node.activeHandlers.Add(1)
		go func() {
			defer node.activeHandlers.Done()
			node.handleConn(conn)
		}()
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
	if resourceErr := node.checkIncomingResourcePromise(fields); resourceErr != nil {
		node.observeOutcome(fromAgent, relationship.OutcomeBroken)
		node.applyBrokenPromiseCost(fromAgent, fields, resourceErr.Error())
		node.record("resource_promise_rejected", "broken", fromAgent, resourceErr.Error())
		node.writeAck(frameConn, fromAgent, "broken", "I promise I rejected this resource promise because local checks failed.")
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

func repairErrDetail(validateErr error) string {
	return "repaired common live decision formatting issue: " + validateErr.Error()
}

// checkLocalResourcePromise verifies that the local agent has enough current
// budget and capacity before making a storage or compute promise.
// Intent: Keep resource promises tied to locally fulfillable behavior rather
// than allowing the LLM to promise impossible storage/compute work.
// Source: DI-duhub
func (node *Node) checkLocalResourcePromise(fields map[string]string) error {
	resourceType := resourceField(fields)
	if resourceType == "" {
		return nil
	}
	requestedUnits := intField(fields, "field_units", "field_requested_units", "field_capacity", "field_capacity_mb")
	if requestedUnits <= 0 {
		return fmt.Errorf("resource promise for %s must declare positive units", resourceType)
	}
	node.mu.Lock()
	defer node.mu.Unlock()
	if requestedUnits > node.capacity {
		return fmt.Errorf("resource promise for %s asks %d units but local capacity is %d", resourceType, requestedUnits, node.capacity)
	}
	if requestedUnits > node.budget {
		return fmt.Errorf("resource promise for %s asks %d units but local budget is %d", resourceType, requestedUnits, node.budget)
	}
	return nil
}

// checkIncomingResourcePromise rejects inbound resource promises that cannot be
// safely interpreted by this bounded POC receiver.
// Intent: Treat malformed or extreme resource promises as local evidence about
// the sender, not as commands the receiver must obey. Source: DI-duhub
func (node *Node) checkIncomingResourcePromise(fields map[string]string) error {
	resourceType := resourceField(fields)
	if resourceType == "" {
		return nil
	}
	requestedUnits := intField(fields, "field_units", "field_requested_units", "field_capacity", "field_capacity_mb")
	if requestedUnits <= 0 {
		return fmt.Errorf("incoming resource promise for %s lacks positive units", resourceType)
	}
	if requestedUnits > 1000 {
		return fmt.Errorf("incoming resource promise for %s exceeds POC safety limit: %d units", resourceType, requestedUnits)
	}
	return nil
}

// applyBrokenPromiseCost spends locally posted stake/collateral when this node
// observes a broken promise with an explicit economic field.
// Intent: Make promise-breaking economically visible inside the POC without
// creating a central penalty authority. Source: DI-duhub
func (node *Node) applyBrokenPromiseCost(peerName string, fields map[string]string, detail string) {
	stakeAmount := intField(fields, "field_stake", "field_collateral", "stake", "collateral")
	if stakeAmount <= 0 {
		return
	}
	node.mu.Lock()
	if stakeAmount > node.budget {
		stakeAmount = node.budget
	}
	node.budget -= stakeAmount
	node.mu.Unlock()
	node.record("stake_forfeited", "broken", peerName, fmt.Sprintf("%s; forfeited %d local budget units", detail, stakeAmount))
}

// resourceField identifies the small set of resource-fulfillment promises this
// POC can check directly. Need advertisements such as "storage_need" are not
// treated as fulfillment promises because they do not claim local capacity.
// Intent: Keep storage/compute checks concrete and avoid misclassifying an
// agent's stated need as a promise to fulfill that need. Source: DI-duhub
func resourceField(fields map[string]string) string {
	for _, key := range []string{"field_resource", "field_resource_type", "resource", "field_promise_about"} {
		value := fields[key]
		if value == "storage" || value == "compute" {
			return value
		}
	}
	return ""
}

func intField(fields map[string]string, keys ...string) int {
	for _, key := range keys {
		value := fields[key]
		if value == "" {
			continue
		}
		parsedValue, parseErr := strconv.Atoi(value)
		if parseErr == nil {
			return parsedValue
		}
	}
	return 0
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
		node.setStopping()
		closeErr := node.listener.Close()
		if closeErr != nil && !errors.Is(closeErr, net.ErrClosed) {
			fmt.Fprintf(os.Stderr, "close listener: %v\n", closeErr)
		}
	}
}

func (node *Node) setStopping() {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.stopping = true
}

func (node *Node) isStopping() bool {
	node.mu.Lock()
	defer node.mu.Unlock()
	return node.stopping
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

// relationshipStatePath is the durable local memory file for this agent's
// private trust ledger.
// Intent: Keep relationship learning across POC11 runs without introducing a
// global trust database or shared authority. Source: DI-duhub
func (node *Node) relationshipStatePath() string {
	return filepath.Join(node.Config.RunRoot, "relationships", node.Agent.Name+".json")
}

// loadRelationshipState restores this agent's prior local trust snapshot if it
// exists; absence simply means this is the first run for that agent.
// Intent: Let multi-run POC11 experiments test relationship decay and repair
// over time while preserving local-only trust semantics. Source: DI-duhub
func (node *Node) loadRelationshipState() error {
	statePath := node.relationshipStatePath()
	stateBytes, readErr := os.ReadFile(statePath)
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) {
			return nil
		}
		return readErr
	}
	var state relationship.State
	if unmarshalErr := json.Unmarshal(stateBytes, &state); unmarshalErr != nil {
		return unmarshalErr
	}
	node.mu.Lock()
	node.ledger.ApplyState(state)
	node.mu.Unlock()
	node.record("relationship_state_loaded", "kept", "", "loaded durable local relationship snapshot")
	return nil
}

// saveRelationshipState writes the local trust snapshot via a temporary file
// and rename so readers never see a partial JSON document.
// Intent: Persist relationship memory after each run while keeping incomplete
// writes from corrupting the next run's local evidence. Source: DI-duhub
func (node *Node) saveRelationshipState() error {
	statePath := node.relationshipStatePath()
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		return err
	}
	node.mu.Lock()
	state := node.ledger.Export()
	node.mu.Unlock()
	stateBytes, marshalErr := json.MarshalIndent(state, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	tempPath := statePath + ".tmp"
	if err := os.WriteFile(tempPath, append(stateBytes, '\n'), 0o644); err != nil {
		return err
	}
	if err := os.Rename(tempPath, statePath); err != nil {
		return err
	}
	node.record("relationship_state_saved", "kept", "", "saved durable local relationship snapshot")
	return nil
}

// writeDoneMarker records idempotent local completion for one agent.
// Intent: Re-running or restarting a bounded POC node should not turn an
// already-written completion marker into broken evidence. Source: DI-duhub
func (node *Node) writeDoneMarker() error {
	donePath := filepath.Join(node.runDir(), node.Agent.Name+".done")
	if _, err := os.Stat(donePath); err == nil {
		node.record("node_done_existing", "kept", "", "done marker already existed")
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
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

// runMonitor writes the observer-only report after all nodes have completed and
// late receive handlers have had a bounded chance to drain.
// Intent: Make monitor completion idempotent and reduce false drift from
// reading logs before in-flight receipts settle. Source: DI-duhub
func (node *Node) runMonitor(ctx context.Context) error {
	donePath := filepath.Join(node.runDir(), "monitor.done")
	if _, err := os.Stat(donePath); err == nil {
		node.record("monitor_done_existing", "kept", "", "monitor marker already existed")
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := node.waitForAllDone(ctx); err != nil {
		return err
	}
	node.drainInflight(ctx)
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
	if err := os.WriteFile(donePath, []byte("done\n"), 0o644); err != nil {
		return err
	}
	node.record("monitor_done", "kept", "", report.Summary)
	return nil
}

// drainInflight waits briefly for active receive handlers before done/monitor
// evidence is finalized.
// Intent: Preserve receipts that were already accepted without letting shutdown
// hang indefinitely. Source: DI-duhub
func (node *Node) drainInflight(ctx context.Context) {
	drained := make(chan struct{})
	go func() {
		node.activeHandlers.Wait()
		close(drained)
	}()
	timer := time.NewTimer(shutdownDrainTimeout)
	defer timer.Stop()
	select {
	case <-drained:
		node.record("inflight_drained", "kept", "", "all active receive handlers completed before done marker")
	case <-ctx.Done():
		node.record("inflight_drain_cancelled", "broken", "", ctx.Err().Error())
	case <-timer.C:
		node.record("inflight_drain_timeout", "non_commitment", "", "some receive handlers may still be running")
	}
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
