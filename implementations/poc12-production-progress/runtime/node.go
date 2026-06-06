package runtime

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/config"
	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/decision"
	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/economy"
	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/pcid"
	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/production"
	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/protocol"
	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/relationship"
	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/transport"
)

const sendTimeout = 5 * time.Second
const shutdownDrainTimeout = 750 * time.Millisecond
const fulfillmentOrderID = "ORDER-1001"
const fulfillmentPackageID = "PKG-1001"
const duplicateShipmentEvidenceField = "field_duplicate_shipment_update"

// Node runs one local POC12 app process. A container may run several app
// processes, but each process keeps its own local relationship ledger, log, and
// live-LLM boundary while a separate container kernel handles byte routing.
// Intent: Apps are local processes that promise to handle pCIDs through their
// local kernel; the kernel does not own app trust or business workflow policy.
// Source: DI-galin
type Node struct {
	Config    config.Config
	Agent     config.AgentConfig
	Protocols pcid.Registry
	Decider   decision.Decider
	Monitor   decision.Monitor

	mu        sync.Mutex
	events    []decision.Event
	ledger    *relationship.Ledger
	evaluator economy.Evaluator
	logFile   *os.File
	budget    int
	capacity  int

	shipmentUpdates map[string]bool
	promiseJournal  map[string]promiseRecord

	activeHandlers sync.WaitGroup
	receiveConns   []transport.FrameConn
	stopping       bool
	drainRecorded  bool
}

type parsedMessage struct {
	Fields       map[string]string
	ExactHash    string
	ProtocolCID  protocol.ProtocolCID
	ProtocolName string
}

// promiseStatus is the app-local journal state for one promise this process has
// enough exact-byte evidence to track.
// Intent: POC12 needs promise-state words that distinguish local failure,
// non-commitment, duplicate evidence, and actual kept/broken outcomes before
// any peer trust update is considered. Source: DI-vujob
type promiseStatus string

const (
	promiseStatusOutstanding   promiseStatus = "outstanding"
	promiseStatusKept          promiseStatus = "kept"
	promiseStatusBroken        promiseStatus = "broken"
	promiseStatusMalformed     promiseStatus = "malformed"
	promiseStatusNonCommitment promiseStatus = "non_commitment"
	promiseStatusDuplicate     promiseStatus = "duplicate"
	promiseStatusLocalFailure  promiseStatus = "local_failure"
)

// promiseRecord is one app-local journal entry for promise evidence this app is
// currently tracking.
// Intent: POC12 applies peer trust only after local promise evidence is recorded
// in the app, never because the kernel, transport, or an unrelated local
// resource check says so. Source: DI-vujob
type promiseRecord struct {
	Key              string
	Fingerprint      string
	Peer             string
	ProtocolName     string
	ExactHash        string
	PromiseAbout     string
	PromiseText      string
	ExpectedEvidence string
	Status           promiseStatus
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
		Config:    cfg,
		Agent:     agent,
		Protocols: pcid.NewRegistry(),
		Decider:   decider,
		Monitor:   monitor,
		ledger:    relationship.NewLedger(peerNames, agent.InitialPeers, cfg.StrongTrustThreshold, cfg.WeakTrustThreshold, cfg.TrustDecayPerRound),
		evaluator: economy.Evaluator{},
		budget:    agent.Budget,
		capacity:  agent.Capacity,

		shipmentUpdates: make(map[string]bool),
		promiseJournal:  make(map[string]promiseRecord),
	}
}

// Run registers local receive promises with the container kernel, executes
// bounded autonomous turns, writes a done marker, and waits for the
// observer-only monitor report.
func (node *Node) Run(ctx context.Context) error {
	if err := node.openLog(); err != nil {
		return err
	}
	defer node.closeLog()
	if err := node.loadRelationshipState(); err != nil {
		return err
	}
	if err := node.registerReceivePromises(ctx); err != nil {
		return err
	}
	defer node.closeReceivePromises()
	time.Sleep(node.Config.StartupDelay())
	if err := node.runStartupWorkflow(ctx); err != nil {
		node.record("startup_workflow_failed", "broken", "", err.Error())
	}
	for turnIndex := 0; turnIndex < node.Config.MaxTurns && turnIndex < node.Config.MaxAgentCalls; turnIndex++ {
		if err := node.runTurn(ctx, turnIndex); err != nil {
			node.recordDecisionError(err)
		}
		time.Sleep(node.Config.TurnDelay())
	}
	if err := node.writeTurnsDoneMarker(); err != nil {
		return err
	}
	node.waitForShutdownGrace(ctx)
	// Intent: Stop receiving kernel-delivered frames before the local done
	// marker is written so `node_done` does not race with late app receipts.
	// Source: DI-galin
	node.closeReceivePromises()
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

func (node *Node) runStartupWorkflow(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	// Intent: Only the fulfillment agent owns the startup production workflow;
	// other agents keep their normal local turn behavior. Source: DI-parok
	if node.Agent.Kind != "fulfillment" {
		return nil
	}
	return node.runFulfillmentShipmentWorkflow()
}

func (node *Node) runFulfillmentShipmentWorkflow() error {
	// Intent: A prompt-only fulfillment agent can discuss shipping without
	// producing evidence. This deterministic startup sequence makes the
	// production workflow executable while later turns remain live/autonomous.
	// Source: DI-parok
	addressAck, addressErr := node.sendAndReceive("accounting", map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  "accounting",
		"turn":                "startup",
		"promise":             "I promise to receive accounting's local address evidence for this order and use it only for this shipment sequence.",
		"reason":              "fulfillment needs address evidence before it can promise label-print evidence",
		"field_promise_about": production.PromiseAddressLookup,
		"field_order_id":      fulfillmentOrderID,
	})
	if addressErr != nil {
		return fmt.Errorf("address lookup: %w", addressErr)
	}
	weightAck, weightErr := node.sendAndReceive("postal_scale", map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  "postal_scale",
		"turn":                "startup",
		"promise":             "I promise to receive postal_scale's local package weight evidence and use it only for this shipment sequence.",
		"reason":              "fulfillment needs local device weight evidence before label printing",
		"field_promise_about": production.PromiseWeighPackage,
		"field_package_id":    fulfillmentPackageID,
	})
	if weightErr != nil {
		return fmt.Errorf("package weighing: %w", weightErr)
	}
	labelAck, labelErr := node.sendAndReceive("ups_label_printer", map[string]string{
		"act":                    decision.ActPromise,
		"from":                   node.Agent.Name,
		"to":                     "ups_label_printer",
		"turn":                   "startup",
		"promise":                "I promise to receive UPS label evidence generated from this address and weight evidence and use it only for this shipment sequence.",
		"reason":                 "fulfillment has address and weight evidence and needs a label promise",
		"field_promise_about":    production.PromisePrintLabel,
		"field_package_id":       fulfillmentPackageID,
		"field_shipping_address": addressAck.Fields["field_shipping_address"],
		"field_weight_ounces":    weightAck.Fields["field_weight_ounces"],
	})
	if labelErr != nil {
		return fmt.Errorf("label printing: %w", labelErr)
	}
	accountingUpdateFields := map[string]string{
		"act":                   decision.ActPromise,
		"from":                  node.Agent.Name,
		"to":                    "accounting",
		"turn":                  "startup",
		"promise":               "I promise to report the shipment cost and tracking evidence I received back to accounting for this order.",
		"reason":                "fulfillment closes its shipment sequence by returning local label evidence to accounting",
		"field_promise_about":   production.PromiseShipmentUpdate,
		"field_order_id":        fulfillmentOrderID,
		"field_tracking_number": labelAck.Fields["field_tracking_number"],
		"field_cost_cents":      labelAck.Fields["field_cost_cents"],
	}
	_, updateErr := node.sendAndReceive("accounting", accountingUpdateFields)
	if updateErr != nil {
		return fmt.Errorf("accounting update: %w", updateErr)
	}
	duplicateUpdateAck, duplicateUpdateErr := node.sendAndReceive("accounting", accountingUpdateFields)
	if duplicateUpdateErr != nil {
		return fmt.Errorf("duplicate accounting update: %w", duplicateUpdateErr)
	}
	if duplicateUpdateAck.Fields[duplicateShipmentEvidenceField] != "true" {
		return fmt.Errorf("duplicate accounting update was not checkpointed")
	}
	node.record("fulfillment_workflow_completed", "kept", "accounting", "order_id="+fulfillmentOrderID+" package_id="+fulfillmentPackageID)
	return nil
}

func (node *Node) runTurn(ctx context.Context, turnIndex int) error {
	if node.Agent.Deterministic() {
		node.record("deterministic_agent_waiting", "kept", "", "deterministic production agent waits for pCID-routed promises")
		return nil
	}
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
	validDecision, validateErr := decision.ValidateObservedPromiseDecision(rawDecision, observation)
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
	if node.suppressRepeatedPromise(validDecision.Target, fields) {
		return nil
	}
	if resourceErr := node.checkLocalResourcePromise(fields); resourceErr != nil {
		node.recordLocalResourceExhaustion(validDecision.Target, fields, resourceErr.Error())
		return nil
	}
	if economicsDecision := node.evaluateEconomics(validDecision.Target, fields); !economicsDecision.PromiseWorthMaking {
		if economicsDecision.Reason == "budget exhausted" || economicsDecision.Reason == "capacity exhausted" {
			node.recordLocalResourceExhaustion(validDecision.Target, fields, economicsDecision.Reason)
			return nil
		}
		node.record("promise_withheld", "non_commitment", validDecision.Target, economicsDecision.Reason)
		return nil
	}
	if sendErr := node.send(validDecision.Target, fields); sendErr != nil {
		sendOutcome, updatesPeerTrust := outcomeForSendError(sendErr)
		if updatesPeerTrust {
			node.observeOutcome(validDecision.Target, sendOutcome)
			node.applyBrokenPromiseCost(validDecision.Target, fields, sendErr.Error())
		}
		sendEventName, sendEventOutcome := sendEventForError(sendErr)
		node.record(sendEventName, sendEventOutcome, validDecision.Target, sendErr.Error())
		return nil
	}
	node.spendLocalCapacity()
	node.record("promise_sent", "kept", validDecision.Target, validDecision.Promise)
	return nil
}

func (node *Node) registerReceivePromises(ctx context.Context) error {
	if node.Config.ListenPort <= 0 {
		node.record("app_kernel_registration_skipped", "kept", "", "no local kernel in unit-test config")
		return nil
	}
	kernelAddress, addressFound := node.Config.KernelAppAddressForAgent(node.Agent.Name)
	if !addressFound {
		node.record("app_kernel_registration_skipped", "kept", "", "no local kernel address for app")
		return nil
	}
	for _, protocolName := range node.Agent.Protocols() {
		if err := node.registerReceivePromise(ctx, kernelAddress, protocolName); err != nil {
			return err
		}
	}
	return nil
}

func (node *Node) registerReceivePromise(ctx context.Context, kernelAddress, protocolName string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	receiveCID, knownReceiveCID := node.Protocols.CID(pcid.KernelReceiveV1)
	if !knownReceiveCID {
		return fmt.Errorf("missing kernel receive pCID")
	}
	targetCID, knownTargetCID := node.Protocols.CID(protocolName)
	if !knownTargetCID {
		return fmt.Errorf("unknown receive pCID %s", protocolName)
	}
	fields := map[string]string{
		"act":      decision.ActPromise,
		"from":     node.Agent.Name,
		"to":       "kernel",
		"app":      node.Agent.Name,
		"pcid":     protocolName,
		"promise":  "I promise to receive exact envelopes for this pCID and judge their promise content locally.",
		"reason":   "local app receive promise registration",
		"pcid_cid": targetCID.String(),
	}
	envelope, envelopeErr := protocol.NewEnvelope(receiveCID, fields, node.Agent.Name)
	if envelopeErr != nil {
		return envelopeErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	frameConn, dialErr := transport.DialFrameConn(kernelAddress, sendTimeout)
	if dialErr != nil {
		return dialErr
	}
	if writeErr := frameConn.WriteFrame(envelopeBytes); writeErr != nil {
		closeErr := frameConn.Close()
		if closeErr != nil {
			node.record("app_kernel_register_close_failed", "broken", "kernel", closeErr.Error())
		}
		return writeErr
	}
	node.mu.Lock()
	node.receiveConns = append(node.receiveConns, frameConn)
	node.mu.Unlock()
	node.activeHandlers.Add(1)
	go node.receiveLoop(protocolName, frameConn)
	node.record("app_receive_promise_sent", "kept", "kernel", "pcid="+protocolName+" kernel="+kernelAddress)
	return nil
}

func (node *Node) receiveLoop(protocolName string, frameConn transport.FrameConn) {
	defer node.activeHandlers.Done()
	for {
		frameBytes, readErr := frameConn.ReadFrame()
		if readErr != nil {
			if node.isStopping() {
				node.record("app_receive_loop_closed", "kept", "kernel", "pcid="+protocolName)
				return
			}
			node.record("app_receive_frame_read_failed", "broken", "kernel", readErr.Error())
			return
		}
		ackBytes, handleErr := node.handleFrame(frameBytes)
		if handleErr != nil {
			node.record("app_receive_frame_rejected", "broken", "kernel", handleErr.Error())
			return
		}
		if writeErr := frameConn.WriteFrame(ackBytes); writeErr != nil {
			node.record("app_receive_ack_write_failed", "broken", "kernel", writeErr.Error())
			return
		}
	}
}

func (node *Node) handleFrame(frameBytes []byte) ([]byte, error) {
	parsed, parseErr := node.parseEnvelope(frameBytes)
	if parseErr != nil {
		node.record("frame_parse_failed", "broken", "", parseErr.Error())
		return nil, parseErr
	}
	fields := parsed.Fields
	fromAgent := fields["from"]
	if !node.supportsProtocol(parsed.ProtocolName) {
		node.record("unsupported_pcid", "non_commitment", fromAgent, "no local app receive promise for "+parsed.ProtocolName)
		return node.newAckBytes(fromAgent, "not_promised", "I promise to remember that I did not promise to handle this pCID.", parsed.ProtocolCID, nil)
	}
	if fields["act"] != decision.ActPromise {
		node.observeOutcome(fromAgent, relationship.OutcomeMalformed)
		node.record("message_rejected", "malformed", fromAgent, "message act is not promise")
		return node.newAckBytes(fromAgent, "malformed", "I promise I rejected this non-promise message.", parsed.ProtocolCID, nil)
	}
	if !node.canAcceptFrom(fromAgent, fields) {
		node.record("message_not_promised", "non_commitment", fromAgent, "no current local promise to accept direct TCP exchange")
		return node.newAckBytes(fromAgent, "not_promised", "I promise to remember that I did not currently promise this direct exchange.", parsed.ProtocolCID, nil)
	}
	promiseID := node.rememberOutstandingPromise(fromAgent, parsed.ProtocolName, parsed.ExactHash, fields)
	if resourceErr := node.checkIncomingResourcePromise(fields); resourceErr != nil {
		node.resolveOutstandingPromise(promiseID, promiseStatusBroken, resourceErr.Error())
		node.observeOutcome(fromAgent, relationship.OutcomeBroken)
		node.applyBrokenPromiseCost(fromAgent, fields, resourceErr.Error())
		node.record("resource_promise_rejected", "broken", fromAgent, resourceErr.Error())
		return node.newAckBytes(fromAgent, "broken", "I promise I rejected this resource promise because local checks failed.", parsed.ProtocolCID, nil)
	}
	ackFields, handlerErr := node.handleProtocolPromise(parsed)
	if handlerErr != nil {
		node.resolveOutstandingPromise(promiseID, promiseStatusBroken, handlerErr.Error())
		node.observeOutcome(fromAgent, relationship.OutcomeBroken)
		node.applyBrokenPromiseCost(fromAgent, fields, handlerErr.Error())
		node.record("protocol_handler_rejected", "broken", fromAgent, handlerErr.Error())
		return node.newAckBytes(fromAgent, "broken", "I promise I rejected this protocol promise because local app checks failed.", parsed.ProtocolCID, nil)
	}
	acceptedAsCandidate := isLinkDiscoveryPromise(fields) && !node.canAccept(fromAgent)
	if evidenceUpdatesTrust(ackFields) {
		trustOutcome := outcomeForPromise(fields)
		node.resolveOutstandingPromise(promiseID, promiseStatusFromOutcome(trustOutcome), "accepted inbound promise")
		node.observeOutcome(fromAgent, trustOutcome)
	} else {
		node.resolveOutstandingPromise(promiseID, promiseStatusDuplicate, "duplicate evidence recorded without trust change")
	}
	eventName := "message_received"
	if acceptedAsCandidate {
		eventName = "candidate_message_received"
	}
	node.record(eventName, "kept", fromAgent, "received "+parsed.ProtocolName+" signed promise exact_sha256="+parsed.ExactHash)
	return node.newAckBytes(fromAgent, "kept", "I promise I received and recorded your signed promise message.", parsed.ProtocolCID, ackFields)
}

func (node *Node) newAckBytes(target, outcome, promiseText string, protocolCID protocol.ProtocolCID, extraFields map[string]string) ([]byte, error) {
	ackFields := map[string]string{
		"act":     decision.ActPromise,
		"from":    node.Agent.Name,
		"to":      target,
		"outcome": outcome,
		"promise": promiseText,
		"reason":  "transport acknowledgement expressed as local promise content",
	}
	for key, value := range extraFields {
		ackFields[key] = value
	}
	ack, ackErr := protocol.NewEnvelope(protocolCID, ackFields, node.Agent.Name)
	if ackErr != nil {
		node.record("ack_sign_failed", "broken", target, ackErr.Error())
		return nil, ackErr
	}
	ackBytes, bytesErr := ack.Bytes()
	if bytesErr != nil {
		node.record("ack_bytes_failed", "broken", target, bytesErr.Error())
		return nil, bytesErr
	}
	return ackBytes, nil
}

func (node *Node) send(target string, fields map[string]string) error {
	_, err := node.sendAndReceive(target, fields)
	return err
}

// sendAndReceive performs one signed promise exchange and returns the receiver's
// ACK evidence to the local caller.
// Intent: The fulfillment workflow needs concrete address, weight, label, and
// accounting evidence from pCID handlers while the app sends only through its
// local kernel rather than dialing peer app processes directly; each outbound
// promise is also journaled before its ACK can affect trust. Source: DI-galin;
// DI-vujob
func (node *Node) sendAndReceive(target string, fields map[string]string) (parsedMessage, error) {
	if !node.canDialTarget(target, fields) {
		return parsedMessage{}, fmt.Errorf("no local TCP promise to %s", target)
	}
	protocolName, protocolCID := node.protocolForFields(fields)
	fields["protocol"] = protocolName
	envelope, envelopeErr := protocol.NewEnvelope(protocolCID, fields, node.Agent.Name)
	if envelopeErr != nil {
		return parsedMessage{}, envelopeErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return parsedMessage{}, bytesErr
	}
	exactHash := protocol.HashExactBytes(envelopeBytes)
	promiseID := node.rememberOutstandingPromise(target, protocolName, exactHash, fields)
	kernelAddress, addressFound := node.Config.KernelAppAddressForAgent(node.Agent.Name)
	if !addressFound {
		node.resolveOutstandingPromise(promiseID, promiseStatusLocalFailure, "missing local kernel endpoint")
		return parsedMessage{}, fmt.Errorf("no local kernel endpoint for app %s", node.Agent.Name)
	}
	frameConn, dialErr := transport.DialFrameConn(kernelAddress, sendTimeout)
	if dialErr != nil {
		node.resolveOutstandingPromise(promiseID, promiseStatusLocalFailure, dialErr.Error())
		return parsedMessage{}, dialErr
	}
	defer node.closeFrameConn(frameConn, "send_close_failed", target)
	if writeErr := frameConn.WriteFrame(envelopeBytes); writeErr != nil {
		node.resolveOutstandingPromise(promiseID, promiseStatusLocalFailure, writeErr.Error())
		return parsedMessage{}, writeErr
	}
	ackBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		node.resolveOutstandingPromise(promiseID, promiseStatusLocalFailure, readErr.Error())
		return parsedMessage{}, readErr
	}
	ackMessage, parseErr := node.parseEnvelope(ackBytes)
	if parseErr != nil {
		node.resolveOutstandingPromise(promiseID, promiseStatusLocalFailure, parseErr.Error())
		return parsedMessage{}, parseErr
	}
	ackFields := ackMessage.Fields
	if ackFields["outcome"] != "kept" {
		ackOutcome, _ := outcomeForSendError(ackOutcomeError{outcome: ackFields["outcome"]})
		node.resolveOutstandingPromise(promiseID, promiseStatusFromOutcome(ackOutcome), "ack outcome "+ackFields["outcome"])
		return parsedMessage{}, ackOutcomeError{outcome: ackFields["outcome"]}
	}
	node.recordAckEvidence(target, ackMessage)
	if evidenceUpdatesTrust(ackFields) {
		trustOutcome := outcomeForPromise(fields)
		node.resolveOutstandingPromise(promiseID, promiseStatusFromOutcome(trustOutcome), "ack kept")
		node.observeOutcome(target, trustOutcome)
	} else {
		node.resolveOutstandingPromise(promiseID, promiseStatusDuplicate, "duplicate ack evidence recorded without trust change")
	}
	return ackMessage, nil
}

func (node *Node) parseEnvelope(frameBytes []byte) (parsedMessage, error) {
	envelope, parseErr := protocol.ParseEnvelope(frameBytes)
	if parseErr != nil {
		return parsedMessage{}, parseErr
	}
	if verifyErr := protocol.VerifyEnvelope(envelope); verifyErr != nil {
		return parsedMessage{}, verifyErr
	}
	fields, fieldsErr := envelope.PayloadFields()
	if fieldsErr != nil {
		return parsedMessage{}, fieldsErr
	}
	if fields["from"] == "" {
		return parsedMessage{}, fmt.Errorf("payload from field is required")
	}
	protocolName, known := node.Protocols.Name(envelope.ProtocolCID)
	if !known {
		protocolName = "unknown:" + envelope.ProtocolCID.String()
	}
	return parsedMessage{
		Fields:       fields,
		ExactHash:    protocol.HashExactBytes(frameBytes),
		ProtocolCID:  envelope.ProtocolCID,
		ProtocolName: protocolName,
	}, nil
}

// protocolForFields chooses the pCID for an outbound promise from protocol or
// promise_about payload meaning. The pCID is still protocol identity, not a
// per-message-type selector; promise_about remains inside the pCID-owned body.
// Source: DI-bikit
func (node *Node) protocolForFields(fields map[string]string) (string, protocol.ProtocolCID) {
	for _, key := range []string{"field_protocol", "protocol"} {
		if protocolName := fields[key]; node.Protocols.Known(protocolName) {
			return protocolName, node.Protocols.MustCID(protocolName)
		}
	}
	switch fields["field_promise_about"] {
	case production.PromiseWeighPackage:
		return pcid.PostalScaleV1, node.Protocols.MustCID(pcid.PostalScaleV1)
	case production.PromiseAddressLookup, production.PromiseShipmentUpdate:
		return pcid.AccountingV1, node.Protocols.MustCID(pcid.AccountingV1)
	case production.PromisePrintLabel:
		return pcid.UPSLabelV1, node.Protocols.MustCID(pcid.UPSLabelV1)
	case production.PromiseIssuePrintCapability, production.PromiseRedeemPrintCapability:
		return pcid.PrinterPortV1, node.Protocols.MustCID(pcid.PrinterPortV1)
	default:
		return pcid.RelationshipV1, node.Protocols.MustCID(pcid.RelationshipV1)
	}
}

func (node *Node) supportsProtocol(protocolName string) bool {
	for _, supportedProtocol := range node.Agent.Protocols() {
		if protocolName == supportedProtocol {
			return true
		}
	}
	return false
}

func (node *Node) handleProtocolPromise(message parsedMessage) (map[string]string, error) {
	switch message.ProtocolName {
	case pcid.RelationshipV1:
		return nil, nil
	case pcid.PostalScaleV1:
		return node.handlePostalScalePromise(message.Fields)
	case pcid.UPSLabelV1:
		return node.handleUPSLabelPromise(message.Fields)
	case pcid.AccountingV1:
		return node.handleAccountingPromise(message.Fields)
	case pcid.PrinterPortV1:
		return node.handlePrinterPortPromise(message.Fields)
	default:
		return nil, fmt.Errorf("unsupported protocol %s", message.ProtocolName)
	}
}

func (node *Node) handlePostalScalePromise(fields map[string]string) (map[string]string, error) {
	if node.Agent.Kind != "postal_scale" {
		return nil, nil
	}
	if fields["field_promise_about"] != production.PromiseWeighPackage {
		return nil, fmt.Errorf("postal scale cannot handle promise_about=%q", fields["field_promise_about"])
	}
	packageID := firstStringField(fields, "field_package_id", "package_id")
	weightOunces, err := production.WeightForPackage(packageID)
	if err != nil {
		return nil, err
	}
	node.record("package_weighed", "kept", fields["from"], fmt.Sprintf("package_id=%s weight_ounces=%d", packageID, weightOunces))
	return map[string]string{
		"field_promise_about": production.PromiseWeighPackage,
		"field_package_id":    packageID,
		"field_weight_ounces": strconv.Itoa(weightOunces),
	}, nil
}

func (node *Node) handleUPSLabelPromise(fields map[string]string) (map[string]string, error) {
	if node.Agent.Kind != "ups_label_printer" {
		return nil, nil
	}
	if fields["field_promise_about"] != production.PromisePrintLabel {
		return nil, fmt.Errorf("ups label printer cannot handle promise_about=%q", fields["field_promise_about"])
	}
	packageID := firstStringField(fields, "field_package_id", "package_id")
	address := firstStringField(fields, "field_shipping_address", "shipping_address", "field_address")
	weightOunces := intField(fields, "field_weight_ounces", "weight_ounces")
	trackingNumber, costCents, err := production.LabelForShipment(packageID, address, weightOunces)
	if err != nil {
		return nil, err
	}
	capabilityAck, err := node.requestPrinterPortCapability()
	if err != nil {
		return nil, err
	}
	labelBytes, err := production.LabelBytesForShipment(map[string]string{
		"field_package_id":      packageID,
		"field_tracking_number": trackingNumber,
		"field_cost_cents":      strconv.Itoa(costCents),
	})
	if err != nil {
		return nil, err
	}
	printAck, err := node.redeemPrinterPortCapability(capabilityAck, labelBytes)
	if err != nil {
		return nil, err
	}
	node.record("shipping_label_printed", "kept", fields["from"], fmt.Sprintf("package_id=%s tracking_number=%s cost_cents=%d", packageID, trackingNumber, costCents))
	return map[string]string{
		"field_promise_about":    production.PromisePrintLabel,
		"field_package_id":       packageID,
		"field_tracking_number":  trackingNumber,
		"field_cost_cents":       strconv.Itoa(costCents),
		"field_printer_spool_id": printAck.Fields["field_printer_spool_id"],
	}, nil
}

// requestPrinterPortCapability asks the local printer-port kernel role for a
// bounded future-print promise token before any label bytes are presented.
// Intent: The UPS label app receives promise-token evidence from the local
// printer resource owner instead of assuming hardware access or treating the
// message kernel as an authorization service. Source: DI-pohaj; DI-vutok
func (node *Node) requestPrinterPortCapability() (parsedMessage, error) {
	tokenID := "printcap-" + node.Agent.Name
	capabilityFields := map[string]string{
		"act":                              decision.ActPromise,
		"from":                             node.Agent.Name,
		"to":                               "printer_port",
		"turn":                             "startup",
		"promise":                          "I promise to receive printer_port's scoped future-print capability token and use it only for bounded UPS label bytes.",
		"reason":                           "ups_label_printer needs local printer-port promise evidence before asking for hardware printing",
		"field_promise_about":              production.PromiseIssuePrintCapability,
		"field_print_capability_issuee":    node.Agent.Name,
		"field_print_capability_token_id":  tokenID,
		"field_print_capability_scope":     production.PrintCapabilityScope,
		"field_print_capability_max_bytes": strconv.Itoa(production.PrintCapabilityMaxBytes),
	}
	capabilityAck, err := node.sendAndReceive("printer_port", capabilityFields)
	if err != nil {
		return parsedMessage{}, err
	}
	return capabilityAck, nil
}

// redeemPrinterPortCapability presents bounded label bytes with the token that
// printer_port previously issued to this app.
// Intent: Hardware access is a reciprocal promise exchange with the local
// resource owner: the label app promises bounded bytes, and printer_port returns
// local print evidence if its own token is still recognizable. Source: DI-pohaj;
// DI-vutok
func (node *Node) redeemPrinterPortCapability(capabilityAck parsedMessage, labelBytes []byte) (parsedMessage, error) {
	redemptionFields := map[string]string{
		"act":                              decision.ActPromise,
		"from":                             node.Agent.Name,
		"to":                               "printer_port",
		"turn":                             "startup",
		"promise":                          "I promise to present only bounded UPS label bytes under this printer_port capability token and to receive printer_port's local print evidence.",
		"reason":                           "ups_label_printer has a scoped future-print token and now asks printer_port to write exact label bytes",
		"field_promise_about":              production.PromiseRedeemPrintCapability,
		"field_print_capability_issuee":    node.Agent.Name,
		"field_print_capability_token":     capabilityAck.Fields["field_print_capability_token"],
		"field_print_capability_token_id":  capabilityAck.Fields["field_print_capability_token_id"],
		"field_print_capability_scope":     capabilityAck.Fields["field_print_capability_scope"],
		"field_print_capability_max_bytes": capabilityAck.Fields["field_print_capability_max_bytes"],
		"field_label_bytes_hex":            hex.EncodeToString(labelBytes),
	}
	printAck, err := node.sendAndReceive("printer_port", redemptionFields)
	if err != nil {
		return parsedMessage{}, err
	}
	return printAck, nil
}

// handlePrinterPortPromise is the local printer-port resource owner's promise
// surface for future print tokens and bounded label-byte redemption.
// Intent: Keep hardware access as voluntary local promises by the agent that
// owns the port, while the kernel only transports exact bytes and the label app
// only receives explicit print evidence after token redemption. Source:
// DI-pohaj; DI-vutok
func (node *Node) handlePrinterPortPromise(fields map[string]string) (map[string]string, error) {
	if node.Agent.Kind != "printer_port" {
		return nil, nil
	}
	switch fields["field_promise_about"] {
	case production.PromiseIssuePrintCapability:
		token, err := production.IssuePrintCapabilityToken(fields)
		if err != nil {
			return nil, err
		}
		tokenID := firstStringField(fields, "field_print_capability_token_id")
		scope := firstStringField(fields, "field_print_capability_scope")
		if scope == "" {
			scope = production.PrintCapabilityScope
		}
		maxBytes := firstStringField(fields, "field_print_capability_max_bytes")
		if maxBytes == "" {
			maxBytes = strconv.Itoa(production.PrintCapabilityMaxBytes)
		}
		node.record("printer_capability_issued", "kept", fields["from"], fmt.Sprintf("token_id=%s scope=%s max_bytes=%s", tokenID, scope, maxBytes))
		return map[string]string{
			"field_promise_about":              production.PromiseIssuePrintCapability,
			"field_print_capability_issuee":    firstStringField(fields, "field_print_capability_issuee", "from"),
			"field_print_capability_token":     token,
			"field_print_capability_token_id":  tokenID,
			"field_print_capability_scope":     scope,
			"field_print_capability_max_bytes": maxBytes,
		}, nil
	case production.PromiseRedeemPrintCapability:
		spoolID, err := production.PrintLabelToLocalDevice(fields)
		if err != nil {
			return nil, err
		}
		printEvidence := firstStringField(fields, "field_label_bytes_hex")
		node.record("printer_port_printed", "kept", fields["from"], fmt.Sprintf("spool_id=%s label_hex_bytes=%d", spoolID, len(printEvidence)))
		return map[string]string{
			"field_promise_about":    production.PromiseRedeemPrintCapability,
			"field_printer_spool_id": spoolID,
		}, nil
	default:
		return nil, fmt.Errorf("printer_port cannot handle promise_about=%q", fields["field_promise_about"])
	}
}

func (node *Node) handleAccountingPromise(fields map[string]string) (map[string]string, error) {
	if node.Agent.Kind != "accounting" {
		return nil, nil
	}
	switch fields["field_promise_about"] {
	case production.PromiseAddressLookup:
		orderID := firstStringField(fields, "field_order_id", "order_id")
		address, err := production.AddressForOrder(orderID)
		if err != nil {
			return nil, err
		}
		node.record("shipping_address_promised", "kept", fields["from"], fmt.Sprintf("order_id=%s shipping_address=%s", orderID, address))
		return map[string]string{
			"field_promise_about":    production.PromiseAddressLookup,
			"field_order_id":         orderID,
			"field_shipping_address": address,
		}, nil
	case production.PromiseShipmentUpdate:
		orderID := firstStringField(fields, "field_order_id", "order_id")
		trackingNumber := firstStringField(fields, "field_tracking_number", "tracking_number")
		costCents := intField(fields, "field_cost_cents", "cost_cents")
		if err := production.ValidateAccountingUpdate(orderID, trackingNumber, costCents); err != nil {
			return nil, err
		}
		ackFields := map[string]string{
			"field_promise_about":   production.PromiseShipmentUpdate,
			"field_order_id":        orderID,
			"field_tracking_number": trackingNumber,
			"field_cost_cents":      strconv.Itoa(costCents),
		}
		updateKey := shipmentUpdateKey(orderID, trackingNumber, costCents)
		node.mu.Lock()
		alreadyRecorded := node.shipmentUpdates[updateKey]
		if !alreadyRecorded {
			node.shipmentUpdates[updateKey] = true
		}
		node.mu.Unlock()
		if alreadyRecorded {
			ackFields[duplicateShipmentEvidenceField] = "true"
			node.record("accounting_update_duplicate", "kept", fields["from"], fmt.Sprintf("order_id=%s tracking_number=%s cost_cents=%d", orderID, trackingNumber, costCents))
			return ackFields, nil
		}
		node.record("accounting_updated", "kept", fields["from"], fmt.Sprintf("order_id=%s tracking_number=%s cost_cents=%d", orderID, trackingNumber, costCents))
		return ackFields, nil
	default:
		return nil, fmt.Errorf("accounting cannot handle promise_about=%q", fields["field_promise_about"])
	}
}

func (node *Node) recordAckEvidence(target string, message parsedMessage) {
	switch message.Fields["field_promise_about"] {
	case production.PromiseWeighPackage:
		node.record("package_weight_received", "kept", target, "weight_ounces="+message.Fields["field_weight_ounces"])
	case production.PromiseAddressLookup:
		node.record("shipping_address_received", "kept", target, "shipping_address="+message.Fields["field_shipping_address"])
	case production.PromisePrintLabel:
		node.record("shipping_label_received", "kept", target, "tracking_number="+message.Fields["field_tracking_number"]+" cost_cents="+message.Fields["field_cost_cents"])
	case production.PromiseShipmentUpdate:
		if message.Fields[duplicateShipmentEvidenceField] == "true" {
			node.record("accounting_update_duplicate_confirmed", "kept", target, "tracking_number="+message.Fields["field_tracking_number"]+" cost_cents="+message.Fields["field_cost_cents"])
			return
		}
		node.record("accounting_update_confirmed", "kept", target, "tracking_number="+message.Fields["field_tracking_number"]+" cost_cents="+message.Fields["field_cost_cents"])
	case production.PromiseIssuePrintCapability:
		node.record("printer_capability_received", "kept", target, "token_id="+message.Fields["field_print_capability_token_id"]+" scope="+message.Fields["field_print_capability_scope"])
	case production.PromiseRedeemPrintCapability:
		node.record("printer_port_print_confirmed", "kept", target, "spool_id="+message.Fields["field_printer_spool_id"])
	}
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
		SupportedPCIDs: node.Agent.Protocols(),
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
	beforePeers := stringSet(node.ledger.DirectPeers())
	node.ledger.DecayRound()
	afterPeers := stringSet(node.ledger.DirectPeers())
	node.mu.Unlock()
	for peerName := range beforePeers {
		if !afterPeers[peerName] {
			node.record(string(relationship.TransitionRemoved), "kept", peerName, "relationship decay crossed weak threshold")
		}
	}
	for peerName := range afterPeers {
		if !beforePeers[peerName] {
			node.record(string(relationship.TransitionAdded), "kept", peerName, "relationship decay/reconfiguration crossed strong threshold")
		}
	}
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

// canDialTarget reports whether this node currently promises to initiate one TCP
// exchange with a peer. Intent: Existing direct peers remain the ordinary path;
// candidate peers are reachable only for explicit low-risk link-discovery
// promises, not arbitrary traffic. Source: DI-timah
func (node *Node) canDialTarget(peerName string, fields map[string]string) bool {
	node.mu.Lock()
	defer node.mu.Unlock()
	if node.ledger.CanDial(peerName) {
		return true
	}
	return isLinkDiscoveryPromise(fields) && containsName(node.Agent.CandidatePeers, peerName)
}

// canAcceptFrom reports whether this node currently promises to accept one TCP
// exchange from a peer. Intent: Candidate-peer discovery is a narrow voluntary
// acceptance promise, not broad permission or a global routing rule.
// Source: DI-timah
func (node *Node) canAcceptFrom(peerName string, fields map[string]string) bool {
	node.mu.Lock()
	defer node.mu.Unlock()
	if node.ledger.CanAccept(peerName) {
		return true
	}
	return isLinkDiscoveryPromise(fields) && containsName(node.Agent.CandidatePeers, peerName)
}

func (node *Node) observeOutcome(peerName string, outcome relationship.Outcome) {
	if peerName == "" || peerName == node.Agent.Name {
		return
	}
	node.mu.Lock()
	transition := node.ledger.ObserveOutcome(peerName, outcome)
	trustScore := node.ledger.Trust(peerName)
	node.mu.Unlock()
	node.record(string(transition), "kept", peerName, fmt.Sprintf("outcome=%s trust=%d", outcome, trustScore))
}

// rememberOutstandingPromise adds one exact promise record to this app's local
// journal before later ACK or receive handling can resolve it.
// Intent: Trust changes should be explainable from a local promise-evidence
// record rather than from transport success, kernel routing, or local resource
// pressure alone. Source: DI-vujob
func (node *Node) rememberOutstandingPromise(peerName, protocolName, exactHash string, fields map[string]string) string {
	fingerprint := promiseRecordKey(peerName, protocolName, "", fields)
	record := promiseRecord{
		Key:              promiseRecordKey(peerName, protocolName, exactHash, fields),
		Fingerprint:      fingerprint,
		Peer:             peerName,
		ProtocolName:     protocolName,
		ExactHash:        exactHash,
		PromiseAbout:     fields["field_promise_about"],
		PromiseText:      fields["promise"],
		ExpectedEvidence: fields["reason"],
		Status:           promiseStatusOutstanding,
	}
	node.mu.Lock()
	node.promiseJournal[record.Key] = record
	node.mu.Unlock()
	node.record("promise_outstanding", "kept", peerName, fmt.Sprintf("status=%s protocol=%s exact_sha256=%s promise_about=%s", record.Status, record.ProtocolName, record.ExactHash, record.PromiseAbout))
	return record.Key
}

// resolveOutstandingPromise records the local outcome of a previously journaled
// promise without itself deciding whether peer trust should change.
// Intent: Promise resolution evidence and trust mutation are deliberately
// separate so duplicate, local-failure, and non-commitment cases stay visible
// without being treated as broken peer promises. Source: DI-vujob
func (node *Node) resolveOutstandingPromise(recordKey string, status promiseStatus, detail string) {
	if recordKey == "" {
		return
	}
	node.mu.Lock()
	record, exists := node.promiseJournal[recordKey]
	if exists {
		record.Status = status
		node.promiseJournal[recordKey] = record
	}
	node.mu.Unlock()
	if !exists {
		node.record("promise_resolution_unmatched", promiseStatusOutcome(status), "", "status="+string(status)+" detail="+detail)
		return
	}
	node.record("promise_resolved", promiseStatusOutcome(status), record.Peer, fmt.Sprintf("status=%s protocol=%s exact_sha256=%s promise_about=%s detail=%s", status, record.ProtocolName, record.ExactHash, record.PromiseAbout, detail))
}

// recordLocalResourceExhaustion records this app's own inability or refusal to
// spend local resources without changing trust in the target peer.
// Intent: Alice exhausting Alice's budget or capacity is evidence about Alice's
// local state, not evidence that Bob kept or broke a promise. Source: DI-vujob
func (node *Node) recordLocalResourceExhaustion(target string, fields map[string]string, detail string) {
	resourceName := resourceField(fields)
	if resourceName == "" {
		resourceName = fields["field_promise_about"]
	}
	node.record("local_resource_exhausted", "non_commitment", target, "resource="+resourceName+" detail="+detail)
}

// suppressRepeatedPromise avoids sending the same live-agent promise text to the
// same target/protocol once this app already has journal evidence for it.
// Intent: Repetition after a prior promise outcome creates pressure that looks
// RPC-like; POC12 should instead record local non-commitment and wait for a new
// promise meaning. Source: DI-vujob
func (node *Node) suppressRepeatedPromise(target string, fields map[string]string) bool {
	protocolName, _ := node.protocolForFields(fields)
	fingerprint := promiseRecordKey(target, protocolName, "", fields)
	node.mu.Lock()
	repeated := false
	for _, record := range node.promiseJournal {
		if record.Fingerprint == fingerprint && record.Status != promiseStatusLocalFailure {
			repeated = true
			break
		}
	}
	node.mu.Unlock()
	if !repeated {
		return false
	}
	node.record("promise_repeated_suppressed", "non_commitment", target, "protocol="+protocolName+" promise_about="+fields["field_promise_about"])
	return true
}

func repairErrDetail(validateErr error) string {
	return "repaired common live decision formatting issue: " + validateErr.Error()
}

// recordDecisionError records LLM/provider/runtime failures as local evidence,
// not as broken peer promises.
// Intent: A transient provider failure or runtime decision failure does not mean
// any peer broke a promise, so it should not enter peer trust as broken
// evidence. Source: DI-jinoz
func (node *Node) recordDecisionError(err error) {
	node.record("decision_error", "non_commitment", "", "local provider/runtime error: "+err.Error())
}

type ackOutcomeError struct {
	outcome string
}

func (err ackOutcomeError) Error() string {
	return fmt.Sprintf("ack outcome %q", err.outcome)
}

// outcomeForSendError converts a transport or ACK failure into the peer-trust
// outcome it actually supports.
// Intent: A receiver's `not_promised` ACK is evidence of non-commitment, not a
// broken peer promise; local transport failures are not peer evidence at all.
// Source: DI-jinoz; DI-vujob
func outcomeForSendError(err error) (relationship.Outcome, bool) {
	var ackErr ackOutcomeError
	if !errors.As(err, &ackErr) {
		return relationship.OutcomeNonCommitment, false
	}
	switch ackErr.outcome {
	case "not_promised", string(relationship.OutcomeNonCommitment):
		return relationship.OutcomeNonCommitment, false
	case string(relationship.OutcomeMalformed):
		return relationship.OutcomeMalformed, true
	default:
		return relationship.OutcomeBroken, true
	}
}

// sendEventForError names send failures without collapsing local transport
// failure, receiver non-commitment, and explicit malformed/broken ACKs.
// Intent: Analyzer output and logs should show why a send did not complete
// without implying a peer broke a promise it never made. Source: DI-vujob
func sendEventForError(err error) (string, string) {
	var ackErr ackOutcomeError
	if !errors.As(err, &ackErr) {
		return "send_unavailable", "non_commitment"
	}
	switch ackErr.outcome {
	case "not_promised", string(relationship.OutcomeNonCommitment):
		return "send_not_promised", "non_commitment"
	case string(relationship.OutcomeMalformed):
		return "send_failed", string(relationship.OutcomeMalformed)
	default:
		return "send_failed", string(relationship.OutcomeBroken)
	}
}

// outcomeForPromise maps a kept payload to the local trust effect it should have.
// Intent: Successful candidate discovery can form a direct relationship while
// ordinary kept promises keep the previous incremental trust behavior.
// Source: DI-timah
func outcomeForPromise(fields map[string]string) relationship.Outcome {
	if isLinkDiscoveryPromise(fields) {
		return relationship.OutcomeDiscoveryKept
	}
	return relationship.OutcomeKept
}

// evidenceUpdatesTrust reports whether ACK payload evidence should change peer
// trust or merely be recorded as already-seen local evidence.
// Intent: Duplicate shipment-update confirmations should remain visible in logs
// without repeatedly increasing trust for the same order/tracking/cost
// checkpoint. Source: DI-jinoz
func evidenceUpdatesTrust(fields map[string]string) bool {
	return fields == nil || fields[duplicateShipmentEvidenceField] != "true"
}

// promiseRecordKey uses exact bytes for concrete promise instances and promise
// text/about fields for repeat-suppression fingerprints before bytes exist.
// Intent: POC12 needs exact-byte promise accounting for real sends while still
// detecting repeated live-agent promise intent before sending again. Source:
// DI-vujob
func promiseRecordKey(peerName, protocolName, exactHash string, fields map[string]string) string {
	if exactHash != "" {
		return peerName + "|" + protocolName + "|" + exactHash
	}
	return peerName + "|" + protocolName + "|" + fields["field_promise_about"] + "|" + fields["promise"]
}

// promiseStatusOutcome maps journal-only statuses into the small outcome
// vocabulary used by existing POC12 logs and analyzer summaries.
// Intent: Local failures and non-commitments should remain non-commitment in
// reports, while duplicate evidence stays kept-but-non-mutating. Source:
// DI-vujob
func promiseStatusOutcome(status promiseStatus) string {
	switch status {
	case promiseStatusBroken:
		return string(relationship.OutcomeBroken)
	case promiseStatusMalformed:
		return string(relationship.OutcomeMalformed)
	case promiseStatusNonCommitment, promiseStatusLocalFailure:
		return string(relationship.OutcomeNonCommitment)
	default:
		return string(relationship.OutcomeKept)
	}
}

// promiseStatusFromOutcome keeps relationship outcomes and journal statuses
// aligned at the boundary where a real ACK or inbound promise is resolved.
// Intent: The journal records the same kept/broken/malformed/non-commitment
// distinction that peer trust code uses, but the journal record happens first.
// Source: DI-vujob
func promiseStatusFromOutcome(outcome relationship.Outcome) promiseStatus {
	switch outcome {
	case relationship.OutcomeBroken:
		return promiseStatusBroken
	case relationship.OutcomeMalformed:
		return promiseStatusMalformed
	case relationship.OutcomeNonCommitment:
		return promiseStatusNonCommitment
	default:
		return promiseStatusKept
	}
}

func shipmentUpdateKey(orderID, trackingNumber string, costCents int) string {
	return fmt.Sprintf("%s|%s|%d", orderID, trackingNumber, costCents)
}

// isLinkDiscoveryPromise recognizes the pCID-owned payload meaning used for
// candidate-peer link formation.
// Intent: Link discovery is represented as promise content under the same
// top-level act, not as a separate protocol verb. Source: DI-timah
func isLinkDiscoveryPromise(fields map[string]string) bool {
	for _, key := range []string{"field_promise_about", "field_meaning", "field_intent", "field_link_intent"} {
		if fields[key] == decision.PromiseAboutLinkDiscovery {
			return true
		}
	}
	return false
}

func containsName(names []string, wantedName string) bool {
	for _, name := range names {
		if name == wantedName {
			return true
		}
	}
	return false
}

func stringSet(names []string) map[string]bool {
	set := make(map[string]bool, len(names))
	for _, name := range names {
		set[name] = true
	}
	return set
}

// checkLocalResourcePromise verifies that the local agent has enough current
// budget and capacity before making a storage or compute promise.
// Intent: Keep resource promises tied to locally fulfillable behavior rather
// than allowing the LLM to promise impossible storage/compute work.
// Source: DI-timah
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
// the sender, not as commands the receiver must obey. Source: DI-timah
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
// creating a central penalty authority. Source: DI-timah
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
// agent's stated need as a promise to fulfill that need. Source: DI-timah
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

func firstStringField(fields map[string]string, keys ...string) string {
	for _, key := range keys {
		if value := fields[key]; value != "" {
			return value
		}
	}
	return ""
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

func (node *Node) closeReceivePromises() {
	node.setStopping()
	node.mu.Lock()
	receiveConns := append([]transport.FrameConn{}, node.receiveConns...)
	node.receiveConns = nil
	node.mu.Unlock()
	for _, frameConn := range receiveConns {
		node.closeFrameConn(frameConn, "app_receive_conn_close_failed", "kernel")
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

// waitForShutdownGrace leaves receive promises open briefly after active turns end.
// Intent: Peers that made late but still legitimate promises get a bounded
// chance to finish their exchange before this node writes `node_done`.
// Source: DI-galin
func (node *Node) waitForShutdownGrace(ctx context.Context) {
	graceDuration := node.Config.ShutdownGrace()
	if graceDuration <= 0 {
		return
	}
	if err := node.waitForAllTurnsDone(ctx, graceDuration); err != nil {
		node.record("shutdown_grace_timeout", "non_commitment", "", err.Error())
		return
	}
	node.record("shutdown_grace_elapsed", "kept", "", "all agents reached turns_done before receive promises close")
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
// Intent: Keep relationship learning across POC12 runs without introducing a
// global trust database or shared authority. Source: DI-timah
func (node *Node) relationshipStatePath() string {
	return filepath.Join(node.Config.RunRoot, "relationships", node.Agent.Name+".json")
}

// loadRelationshipState restores this agent's prior local trust snapshot if it
// exists; absence simply means this is the first run for that agent.
// Intent: Let multi-run POC12 experiments test relationship decay and repair
// over time while preserving local-only trust semantics. Source: DI-timah
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
// writes from corrupting the next run's local evidence. Source: DI-timah
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
// already-written completion marker into broken evidence. Source: DI-timah
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

// writeTurnsDoneMarker records that this app has finished making active turn
// decisions but is still keeping its receive promises open for peers that are
// finishing their own planned sends.
// Intent: Coordinate shutdown on agent-turn completion instead of letting early
// finishers close receive promises while slower peers are still making
// promises. Source: DI-galin
func (node *Node) writeTurnsDoneMarker() error {
	turnsDonePath := filepath.Join(node.runDir(), node.Agent.Name+".turns_done")
	if _, err := os.Stat(turnsDonePath); err == nil {
		node.record("turns_done_existing", "kept", "", "turns-done marker already existed")
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.WriteFile(turnsDonePath, []byte("turns_done\n"), 0o644); err != nil {
		return err
	}
	node.record("turns_done", "kept", "", "finished active turns and kept receive promises open for shutdown grace")
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
// reading logs before in-flight receipts settle. Source: DI-timah
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
// hang indefinitely. Source: DI-timah
func (node *Node) drainInflight(ctx context.Context) {
	if !node.markDrainStarted() {
		return
	}
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

func (node *Node) markDrainStarted() bool {
	node.mu.Lock()
	defer node.mu.Unlock()
	if node.drainRecorded {
		return false
	}
	node.drainRecorded = true
	return true
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

func (node *Node) waitForAllTurnsDone(ctx context.Context, timeoutDuration time.Duration) error {
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	timer := time.NewTimer(timeoutDuration)
	defer timer.Stop()
	for {
		if node.allTurnsDone() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return fmt.Errorf("timed out waiting for all turns_done markers")
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

func (node *Node) allTurnsDone() bool {
	for _, agentName := range node.Config.AgentNames() {
		if _, err := os.Stat(filepath.Join(node.runDir(), agentName+".turns_done")); err != nil {
			return false
		}
	}
	return true
}

func (node *Node) readAllEvents() ([]decision.Event, error) {
	var events []decision.Event
	// Intent: The observer-only monitor should see app-local evidence and
	// kernel-local operational evidence without giving the kernel authority over
	// trust interpretation. Source: DI-galin
	logPaths, globErr := filepath.Glob(filepath.Join(node.runDir(), "*.jsonl"))
	if globErr != nil {
		return nil, globErr
	}
	sort.Strings(logPaths)
	for _, logPath := range logPaths {
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
