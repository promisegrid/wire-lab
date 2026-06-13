package runtime

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"promisegrid.dev/wire-lab/implementations/poc14-wasm/boundary"
	"promisegrid.dev/wire-lab/implementations/poc14-wasm/pcid"
	"promisegrid.dev/wire-lab/implementations/poc14-wasm/protocol"
)

// runWASMBoundaryWorkflow records one standalone WASM-role process proving that
// it can keep PromiseGrid semantics outside the sandbox while still exchanging
// ordinary relationship_v1 promise envelopes through the local kernel.
// Intent: POC14 adds heterogeneous process-runtime adapter events without making
// WASM host calls into RPC commands or trust authorities. Source: DI-linof
func (node *Node) runWASMBoundaryWorkflow() error {
	if err := boundary.ValidateWASMModule(boundary.MinimalWASMModule); err != nil {
		return err
	}
	target, targetFound := node.firstConfiguredPeerByKind("stdio_agent")
	if !targetFound {
		target = "victor"
	}
	moduleHash := protocol.HashExactBytes(boundary.MinimalWASMModule)
	node.record("wasm_process_agent_started", "kept", "", "agent="+node.Agent.Name+" module_sha256="+moduleHash)
	node.record("wasm_module_header_validated", "kept", "", "module_magic=0061736d module_version=01000000")
	fields := boundary.PromiseFields(
		node.Agent.Name,
		target,
		boundary.PromiseAboutWASMAdapter,
		"Peggy promises that her local WASM-adapter process validated sandbox module bytes and will exchange only pCID-defined PromiseGrid envelopes.",
	)
	fields["field_wasm_module_sha256"] = moduleHash
	fields["field_protocol"] = pcid.RelationshipV1
	if _, err := node.sendAndReceive(target, fields); err != nil {
		return fmt.Errorf("wasm adapter promise: %w", err)
	}
	node.record("wasm_adapter_promise_sent", "kept", target, "pcid="+pcid.RelationshipV1+" module_sha256="+moduleHash)
	node.record("wasm_adapter_ack_received", "kept", target, "pcid="+pcid.RelationshipV1+" stdio peer accepted WASM-adapter event as a local promise")
	usefulTarget := "dave"
	usefulFields := boundary.PromiseFields(
		node.Agent.Name,
		usefulTarget,
		boundary.PromiseAboutWASMModuleUse,
		"Peggy promises Dave a reusable WASM module-validation observation: these module bytes have the WebAssembly magic/version header and can be judged before any sandbox host call is trusted.",
	)
	usefulFields["field_wasm_module_sha256"] = moduleHash
	usefulFields["field_protocol"] = pcid.RelationshipV1
	if _, err := node.sendAndReceive(usefulTarget, usefulFields); err != nil {
		return fmt.Errorf("wasm useful-work promise: %w", err)
	}
	// Intent: Peggy's WASM process should do useful PromiseGrid work, not merely
	// prove that a sandbox boundary exists. Source: DI-pamob
	node.record("wasm_useful_work_promised", "kept", usefulTarget, "pcid="+pcid.RelationshipV1+" module_sha256="+moduleHash)
	node.record("wasm_useful_work_ack_received", "kept", usefulTarget, "pcid="+pcid.RelationshipV1+" Dave received reusable module-validation event")
	return nil
}

// runStdioBoundaryWorkflow starts a worker process whose only application
// messaging path is stdin/stdout. The adapter forwards the exact signed envelope
// through the same local kernel path every other app uses.
// Intent: POC14 tests subprocess adapters without giving stdio messages RPC
// semantics; the worker emits and observes PromiseGrid envelopes as bytes.
// Source: DI-linof
func (node *Node) runStdioBoundaryWorkflow(ctx context.Context) error {
	target, targetFound := node.firstConfiguredPeerByKind("wasm_agent")
	if !targetFound {
		target = "peggy"
	}
	command := exec.CommandContext(ctx, "poc14-stdio-worker")
	command.Stderr = os.Stderr
	stdin, stdinErr := command.StdinPipe()
	if stdinErr != nil {
		return stdinErr
	}
	stdout, stdoutErr := command.StdoutPipe()
	if stdoutErr != nil {
		return stdoutErr
	}
	if startErr := command.Start(); startErr != nil {
		return startErr
	}
	node.record("stdio_worker_started", "kept", target, "worker=poc14-stdio-worker")
	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	encoder := json.NewEncoder(stdin)
	request := boundary.StdioRequest{Type: "promise_request", From: node.Agent.Name, To: target}
	if err := encoder.Encode(request); err != nil {
		return node.finishStdioWorker(command, stdin, fmt.Errorf("write stdio request: %w", err))
	}
	outbound, outboundErr := readStdioEnvelope(scanner)
	if outboundErr != nil {
		return node.finishStdioWorker(command, stdin, outboundErr)
	}
	envelopeBytes, decodeErr := hex.DecodeString(outbound.Hex)
	if decodeErr != nil {
		return node.finishStdioWorker(command, stdin, decodeErr)
	}
	node.record("stdio_worker_envelope_received", "kept", target, "protocol="+outbound.Protocol+" exact_sha256="+protocol.HashExactBytes(envelopeBytes))
	_, ackBytes, sendErr := node.sendRawEnvelopeBytes(outbound.To, outbound.Protocol, envelopeBytes)
	if sendErr != nil {
		return node.finishStdioWorker(command, stdin, fmt.Errorf("stdio envelope forward: %w", sendErr))
	}
	node.record("stdio_adapter_kernel_forwarded", "kept", outbound.To, "pcid="+outbound.Protocol+" worker="+outbound.From)
	if err := encoder.Encode(boundary.StdioAckMessage{Type: "ack_envelope", Hex: hex.EncodeToString(ackBytes)}); err != nil {
		return node.finishStdioWorker(command, stdin, fmt.Errorf("write stdio ack: %w", err))
	}
	observed, observedErr := readStdioEvent(scanner)
	if observedErr != nil {
		return node.finishStdioWorker(command, stdin, observedErr)
	}
	node.record("stdio_worker_ack_event", observed.Outcome, target, "exact_sha256="+observed.ExactSHA256)
	usefulTarget := "dave"
	usefulFields := boundary.PromiseFields(
		node.Agent.Name,
		usefulTarget,
		boundary.PromiseAboutStdioWorkerUse,
		"Victor promises Dave that a stdio-only subprocess round-tripped an exact signed PromiseGrid envelope and observed the peer ACK bytes without direct network access.",
	)
	usefulFields["field_ack_exact_sha256"] = observed.ExactSHA256
	usefulFields["field_protocol"] = pcid.RelationshipV1
	if _, err := node.sendAndReceive(usefulTarget, usefulFields); err != nil {
		return node.finishStdioWorker(command, stdin, fmt.Errorf("stdio useful-work promise: %w", err))
	}
	// Intent: Victor's stdio worker should produce reusable relationship event records
	// about subprocess messaging, not only adapter plumbing logs. Source:
	// DI-pamob
	node.record("stdio_useful_work_promised", "kept", usefulTarget, "pcid="+pcid.RelationshipV1+" ack_exact_sha256="+observed.ExactSHA256)
	node.record("stdio_useful_work_ack_received", "kept", usefulTarget, "pcid="+pcid.RelationshipV1+" Dave received stdio subprocess round-trip event")
	return node.finishStdioWorker(command, stdin, nil)
}

func readStdioEnvelope(scanner *bufio.Scanner) (boundary.StdioEnvelopeMessage, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return boundary.StdioEnvelopeMessage{}, err
		}
		return boundary.StdioEnvelopeMessage{}, fmt.Errorf("stdio worker produced no envelope")
	}
	var outbound boundary.StdioEnvelopeMessage
	if err := json.Unmarshal(scanner.Bytes(), &outbound); err != nil {
		return boundary.StdioEnvelopeMessage{}, err
	}
	if outbound.Type != "outbound_envelope" {
		return boundary.StdioEnvelopeMessage{}, fmt.Errorf("stdio worker message type %q, want outbound_envelope", outbound.Type)
	}
	if outbound.Hex == "" || outbound.From == "" || outbound.To == "" || outbound.Protocol == "" {
		return boundary.StdioEnvelopeMessage{}, fmt.Errorf("stdio worker envelope message is incomplete")
	}
	return outbound, nil
}

func readStdioEvent(scanner *bufio.Scanner) (boundary.StdioEventMessage, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return boundary.StdioEventMessage{}, err
		}
		return boundary.StdioEventMessage{}, fmt.Errorf("stdio worker produced no ack observation")
	}
	var observed boundary.StdioEventMessage
	if err := json.Unmarshal(scanner.Bytes(), &observed); err != nil {
		return boundary.StdioEventMessage{}, err
	}
	if observed.Type != "ack_event" {
		return boundary.StdioEventMessage{}, fmt.Errorf("stdio worker message type %q, want ack_event", observed.Type)
	}
	if observed.Outcome == "" || observed.ExactSHA256 == "" {
		return boundary.StdioEventMessage{}, fmt.Errorf("stdio worker ack observation is incomplete")
	}
	return observed, nil
}

func (node *Node) finishStdioWorker(command *exec.Cmd, stdin ioWriteCloser, cause error) error {
	closeErr := stdin.Close()
	if closeErr != nil && cause == nil {
		cause = closeErr
	}
	waitErr := command.Wait()
	if waitErr != nil && cause == nil {
		cause = waitErr
	}
	return cause
}

type ioWriteCloser interface {
	Close() error
}

// recordDecentralizedMonitoringEvents emits POC14 event records for analysis
// approaches that can work when no agent has a global production view.
// Intent: Monitoring-like behavior must be modeled as local promises, peer
// event records, exchange-rate signals, and voluntary disclosures. Source: DI-lulof
func (node *Node) recordDecentralizedMonitoringEvents() {
	node.record("decentralized_monitoring_model_recorded", "kept", "", "no global analyzer or monitor exists in the production trust model")
	node.record("local_event_summary_promised", "kept", "dave", "Alice promises a signed summary of Alice's own keep/break observations; Dave judges it locally")
	node.record("peer_carried_attestation_promised", "kept", "carol", "Alice carries Bob's signed storage event to Carol without making Bob authoritative")
	node.record("bearer_token_exchange_rate_observed", "kept", "grace", "Alice locally observes peer offers for bearer tokens as a market signal, not a global exchange rate")
	node.record("relationship_topology_signal_observed", "kept", "frank", "direct links, relay willingness, and replica choices are local relationship-strength signals")
	node.record("voluntary_gossip_promised", "kept", "ellen", "Ellen may promise selected local observations; recipients treat them as event records only")
}

// recordPermanentDistrustAndTransitExclusionEvents records two POC14 local
// trust-boundary scenarios before later protocol work implements true multi-hop
// route selection.
// Intent: Alice can decide to permanently distrust Mallory, and Alice can promise
// that Alice's own inbound/outbound traffic should not transit Mallory, without
// pretending to impose a global ban or route policy on other agents. Source:
// DI-kinaf
func (node *Node) recordPermanentDistrustAndTransitExclusionEvents() error {
	node.markPermanentDistrustAndTransitExclusion("mallory")
	node.record("permanent_distrust_decided", "kept", "mallory", "Alice locally decides Mallory is permanently distrusted after repeated malformed/broken promise event")
	node.record("permanent_distrust_future_repair_not_promised", "non_commitment", "mallory", "Alice does not promise to consider future Mallory repair promises without a separate explicit local decision")
	node.record("permanent_distrust_direct_peer_removed", "kept", "mallory", "Alice removes Mallory from Alice's local direct-peer set for Alice-owned traffic")
	node.record("transit_exclusion_promised", "kept", "mallory", "Alice promises that Alice's input and output traffic should not use Mallory as a transit peer")
	node.record("input_transit_exclusion_recorded", "kept", "mallory", "Alice records that inbound traffic claiming Mallory as transit is not acceptable event for sensitive payload delivery")
	node.record("output_transit_exclusion_recorded", "kept", "mallory", "Alice records that outbound traffic candidates naming Mallory as transit are locally rejected")
	if _, err := node.sendAndReceive("mallory", map[string]string{
		"act":                 "promise",
		"from":                node.Agent.Name,
		"to":                  "mallory",
		"turn":                "startup",
		"promise":             "Alice would normally promise a low-risk discovery probe, but Alice's permanent distrust state blocks this local send.",
		"reason":              "permanent distrust should override ordinary candidate-peer repair or discovery",
		"field_promise_about": "link_discovery",
	}); err == nil {
		return fmt.Errorf("permanent distrust unexpectedly allowed send to mallory")
	}
	node.record("permanent_distrust_send_blocked", "non_commitment", "mallory", "Alice's local permanent-distrust state blocked a candidate discovery send before bytes left Alice")
	blockedRoute := []string{"alice", "mallory", "carol"}
	if node.routeAllowed(blockedRoute) {
		return fmt.Errorf("transit exclusion unexpectedly allowed route %v", blockedRoute)
	}
	node.record("transit_candidate_rejected", "non_commitment", "mallory", "Alice rejects route candidate alice->mallory->carol because Mallory is locally untrusted as transit")
	node.record("transit_route_candidate_blocked", "non_commitment", "mallory", "Alice's route-selection check rejected a candidate path containing Mallory as transit")
	safeRoute := []string{"alice", "frank", "carol"}
	if !node.routeAllowed(safeRoute) {
		return fmt.Errorf("transit exclusion unexpectedly rejected route %v", safeRoute)
	}
	node.record("transit_safe_route_selected", "kept", "frank", "Alice selects a non-Mallory transit candidate from Alice's own trusted peer events")
	return nil
}

// recordMixedVersionPCIDMigrationEvents records how an app can reason about a
// legacy pCID and the current pCID without any central registry.
// Intent: Mixed-version compatibility remains local event records and pCID-selected
// payload semantics, not a command to conform to a global version. Source:
// DI-linof
func (node *Node) recordMixedVersionPCIDMigrationEvents() {
	legacyCID := protocol.NewProtocolCID([]byte("poc13 relationship trust discovery observation protocol v1"))
	currentCID := node.Protocols.MustCID(pcid.RelationshipV1)
	node.record("mixed_version_pcid_migration_promised", "kept", "bob", "legacy_pcid="+legacyCID.String()+" current_pcid="+currentCID.String())
	node.record("mixed_version_legacy_pcid_observed", "kept", "bob", "Alice retains legacy pCID event as local context, not as current authority")
	node.record("mixed_version_successor_pcid_selected", "kept", "bob", "Alice selects current relationship_v1 pCID for new envelopes while preserving legacy event")
}

// recordRunInternalRestartEvents records the run-scoped restart contract for
// multiple app processes before the normal save/load code persists local state.
// Intent: POC14 should test crash/restart expectations as local promises by each
// process, not as cross-run state or supervisor authority. Source: DI-linof
func (node *Node) recordRunInternalRestartEvents() {
	node.record("run_internal_restart_orchestration_promised", "kept", "bob", "Alice promises a bounded same-run restart probe for storage state and boundary agents")
	node.record("run_internal_restart_checkpoint_promised", "kept", "peggy", "Bob and Peggy are expected to recover only run-scoped journals, not cross-run state")
	node.record("run_internal_restart_recovery_observed", "kept", "victor", "restart recovery is judged by local run-scoped event and exact envelope replay windows")
}

func (node *Node) firstConfiguredPeerByKind(kind string) (string, bool) {
	for _, peerName := range node.Agent.InitialPeers {
		peer, ok := node.Config.Agent(peerName)
		if ok && peer.Kind == kind {
			return peer.Name, true
		}
	}
	for _, peerName := range node.Agent.CandidatePeers {
		peer, ok := node.Config.Agent(peerName)
		if ok && peer.Kind == kind {
			return peer.Name, true
		}
	}
	return "", false
}
