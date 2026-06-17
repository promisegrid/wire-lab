package runtime

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/config"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/decision"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/economy"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/eventstream"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/pcid"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/production"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/protocol"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/relationship"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/runtimeadapter"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/transport"
)

const sendTimeout = 5 * time.Second
const shutdownDrainTimeout = 750 * time.Millisecond
const fulfillmentOrderID = "ORDER-1001"
const fulfillmentPackageID = "PKG-1001"
const duplicateShipmentEventField = "duplicate_shipment_update"

// Node runs one local POC15 app process. A container may run several app
// processes, but each process keeps its own local relationship ledger, log, and
// live-LLM interface while a separate container kernel handles byte routing.
// Intent: Apps are local processes that promise to handle pCIDs through their
// local kernel; the kernel does not own app trust or business workflow policy.
// Source: DI-galin
type Node struct {
	Config    config.Config
	Agent     config.AgentConfig
	Protocols pcid.Registry
	Decider   decision.Decider
	Monitor   decision.Monitor

	mu sync.Mutex
	// stdoutMu keeps event and artifact JSON lines from interleaving when
	// receive loops emit large base64 artifact records concurrently. Intent:
	// Supervisors parse stdout line-by-line, so each line must remain intact as
	// observer-only harness input. Source: DI-tuhop
	stdoutMu  sync.Mutex
	events    []decision.Event
	ledger    *relationship.Ledger
	evaluator economy.Evaluator
	logFile   *os.File
	budget    int
	capacity  int

	nonCommitmentJournal map[string]nonCommitmentRecord
	checkpointJournal    map[string]checkpointRecord
	promiseJournal       map[string]promiseRecord
	eventOutcomeCounts   map[string]int
	agentCASStore        map[string]agentCASObject
	agentMessageDAG      map[string]agentMessageDAGNode
	capabilityTokens     map[string]string
	computeCache         map[string]map[string]string
	replayJournal        map[string]string
	exchangeCounter      int

	activeHandlers sync.WaitGroup
	appSession     *transport.PersistentSession
	stopping       bool
	drainRecorded  bool
}

type parsedMessage struct {
	Fields             map[string]string
	ExactHash          string
	RawBytes           []byte
	ProtocolCID        protocol.ProtocolCID
	ProtocolName       string
	ParentExactSHA256s []string
}

// protocolHandlerResult is the app-local handler result for one inbound
// pCID-owned promise.
// Intent: Most handlers return compatibility fields that the local app signs as
// a fresh ACK, while Victor's stdio compute handler returns exact ACK bytes that
// the worker signed after computing through stdin/stdout. Source: DI-sivis
type protocolHandlerResult struct {
	Fields   map[string]string
	AckBytes []byte
}

// promiseStatus is the app-local journal state for one promise this process has
// enough exact-byte event records to track.
// Intent: POC15 needs promise-state words that distinguish local failure,
// non-commitment, duplicate event records, and actual kept/broken outcomes before
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

// promiseRecord is one app-local journal entry for a promise event record this app is
// currently tracking.
// Intent: POC15 applies peer trust only after a local promise event is recorded
// in the app, never because the kernel, transport, or an unrelated local
// resource check says so. Source: DI-vujob
type promiseRecord struct {
	Key           string
	Fingerprint   string
	Peer          string
	ProtocolName  string
	ExactHash     string
	PromiseAbout  string
	PromiseText   string
	ExpectedEvent string
	Status        promiseStatus
}

// nonCommitmentRecord remembers one receiver-side `not_promised` outcome from
// this app's local vantage.
// Intent: Receiver non-commitment is an event record showing that this app should stop
// pressuring the same peer for the same semantic promise during the current run;
// it is not an event showing that the peer broke a promise. Source: DI-zapab
type nonCommitmentRecord struct {
	Key          string
	Peer         string
	ProtocolName string
	PromiseAbout string
	Detail       string
}

// checkpointRecord is a reusable app-local marker for a promise event record that
// should be visible but should not keep changing trust when replayed.
// Intent: Duplicate detection belongs to the app's local events journal, not
// to a global idempotency layer or receiver command surface. Source: DI-zapab
type checkpointRecord struct {
	Key          string
	ProtocolName string
	PromiseAbout string
	Subject      string
	Detail       string
}

// runScopedState is the restartable state this POC keeps only inside one run
// root. Intent: Per-agent sparse CAS metadata, message-DAG indexes, compute
// checkpoints, replay windows, and app-local event journals may survive a process
// restart within one experiment, but object bytes live in filesystem CAS files and
// the clean-run reset remains the experiment scope that prevents stale state from
// muddying the next POC15 run. Source: DI-sunuf; DI-manul; DI-fagog
type runScopedState struct {
	Version int `json:"version"`
	// CASObjects is read-only legacy migration input from earlier POC15 state
	// files. Intent: New saves omit base64 CAS bytes so durable-state.json remains
	// a small mutable root and index. Source: DI-fagog
	CASObjects           map[string]string              `json:"cas_objects_b64,omitempty"`
	CapabilityTokens     map[string]string              `json:"capability_tokens,omitempty"`
	ComputeCache         map[string]map[string]string   `json:"compute_cache,omitempty"`
	NonCommitmentJournal map[string]nonCommitmentRecord `json:"non_commitment_journal,omitempty"`
	CheckpointJournal    map[string]checkpointRecord    `json:"checkpoint_journal,omitempty"`
	PromiseJournal       map[string]promiseRecord       `json:"promise_journal,omitempty"`
	EventOutcomeCounts   map[string]int                 `json:"event_outcome_counts,omitempty"`
	ReplayJournal        map[string]string              `json:"replay_journal,omitempty"`
	AgentCASObjects      map[string]agentCASObject      `json:"agent_cas_objects,omitempty"`
	AgentMessageDAG      map[string]agentMessageDAGNode `json:"agent_message_dag,omitempty"`
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

		nonCommitmentJournal: make(map[string]nonCommitmentRecord),
		checkpointJournal:    make(map[string]checkpointRecord),
		promiseJournal:       make(map[string]promiseRecord),
		eventOutcomeCounts:   make(map[string]int),
		agentCASStore:        make(map[string]agentCASObject),
		agentMessageDAG:      make(map[string]agentMessageDAGNode),
		capabilityTokens:     make(map[string]string),
		computeCache:         make(map[string]map[string]string),
		replayJournal:        make(map[string]string),
	}
}

// Run registers local receive promises with the container kernel, executes
// bounded autonomous turns, and exits after a local receive-grace period.
// Intent: Earlier POC15 nodes wrote marker files and one node ran the monitor
// after all markers appeared; DI-dirat moves that observer lifecycle to the
// supervisor and collector so app processes no longer coordinate through the
// Docker observer volume. Source: DI-galin; DI-jupob; DI-dirat
func (node *Node) Run(ctx context.Context) error {
	if err := node.openLog(); err != nil {
		return err
	}
	defer node.closeLog()
	if err := node.loadRunScopedState(); err != nil {
		return err
	}
	if err := node.loadRelationshipState(); err != nil {
		return err
	}
	if err := node.registerReceivePromises(ctx); err != nil {
		return err
	}
	node.record("runtime_readiness_promised", "kept", "", "app receive promises registered with local kernel")
	defer node.closeReceivePromises()
	time.Sleep(node.Config.StartupDelay())
	node.record("peer_readiness_observed", "kept", "", "startup delay elapsed after local receive promises")
	node.recordAgentCASAccessEvents()
	if err := node.runStartupWorkflow(ctx); err != nil {
		node.record("startup_workflow_failed", "broken", "", err.Error())
	}
	for turnIndex := 0; turnIndex < node.Config.MaxTurns && turnIndex < node.Config.MaxAgentCalls; turnIndex++ {
		if err := node.runTurn(ctx, turnIndex); err != nil {
			node.recordDecisionError(err)
		}
		time.Sleep(node.Config.TurnDelay())
	}
	node.waitForShutdownGrace(ctx)
	// Intent: Stop receiving kernel-delivered frames before the local done
	// event record is written so process completion does not race with late app
	// receipts. Source: DI-galin; DI-dirat
	node.closeReceivePromises()
	node.drainInflight(ctx)
	node.recordRunScopedRetentionAndGC()
	if err := node.saveRunScopedState(); err != nil {
		return err
	}
	if err := node.saveRelationshipState(); err != nil {
		return err
	}
	node.record("runtime_done_promised", "kept", "", "app completed local turns and saved relationship and run-scoped event")
	node.record("node_done", "kept", "", "local process finished without observer-volume coordination")
	return nil
}

func (node *Node) runStartupWorkflow(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	// Intent: POC15 is now a superset of POC12, so the production shipping
	// workflow remains intact while Alice and Mallory add CAS/compute protocol
	// pressure above the same app/kernel interface, while POC15 adds WASM and
	// stdio adapter roles. Source: DI-sinur; DI-linof; DI-kimim
	if _, hasKernelAddress := node.Config.KernelAppAddressForAgent(node.Agent.Name); hasKernelAddress {
		switch node.Agent.Kind {
		case "wasm_agent":
			return node.runWASMAdapterWorkflow(ctx)
		case "stdio_agent":
			return node.runStdioAdapterWorkflow(ctx)
		}
		switch node.Agent.Name {
		case "alice":
			if err := node.runCASComputeWorkflow(); err != nil {
				return err
			}
		case "mallory":
			if err := node.runAdversaryWorkflow(); err != nil {
				return err
			}
		}
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
	// producing event records. This deterministic startup sequence makes the
	// production workflow executable while later turns remain live/autonomous.
	// Source: DI-parok
	addressAck, addressErr := node.sendAndReceive("accounting", map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "accounting",
		"turn":          "startup",
		"promise":       "I promise to receive accounting's local address event for this order and use it only for this shipment sequence.",
		"reason":        "fulfillment needs an address event record before it can promise a label-print event record",
		"promise_about": production.PromiseAddressLookup,
		"order_id":      fulfillmentOrderID,
	})
	if addressErr != nil {
		return fmt.Errorf("address lookup: %w", addressErr)
	}
	weightAck, weightErr := node.sendAndReceive("postal_scale", map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "postal_scale",
		"turn":          "startup",
		"promise":       "I promise to receive postal_scale's local package weight event and use it only for this shipment sequence.",
		"reason":        "fulfillment needs local device weight event before label printing",
		"promise_about": production.PromiseWeighPackage,
		"package_id":    fulfillmentPackageID,
	})
	if weightErr != nil {
		return fmt.Errorf("package weighing: %w", weightErr)
	}
	labelAck, labelErr := node.sendAndReceive("ups_label_printer", map[string]string{
		"act":              decision.ActPromise,
		"from":             node.Agent.Name,
		"to":               "ups_label_printer",
		"turn":             "startup",
		"promise":          "I promise to receive UPS label event generated from this address and weight event and use it only for this shipment sequence.",
		"reason":           "fulfillment has address and weight event and needs a label promise",
		"promise_about":    production.PromisePrintLabel,
		"package_id":       fulfillmentPackageID,
		"shipping_address": addressAck.Fields["shipping_address"],
		"weight_ounces":    weightAck.Fields["weight_ounces"],
	})
	if labelErr != nil {
		return fmt.Errorf("label printing: %w", labelErr)
	}
	accountingUpdateFields := map[string]string{
		"act":             decision.ActPromise,
		"from":            node.Agent.Name,
		"to":              "accounting",
		"turn":            "startup",
		"promise":         "I promise to report the shipment cost and tracking event I received back to accounting for this order.",
		"reason":          "fulfillment closes its shipment sequence by returning local label event to accounting",
		"promise_about":   production.PromiseShipmentUpdate,
		"order_id":        fulfillmentOrderID,
		"tracking_number": labelAck.Fields["tracking_number"],
		"cost_cents":      labelAck.Fields["cost_cents"],
	}
	_, updateErr := node.sendAndReceive("accounting", accountingUpdateFields)
	if updateErr != nil {
		return fmt.Errorf("accounting update: %w", updateErr)
	}
	duplicateUpdateAck, duplicateUpdateErr := node.sendAndReceive("accounting", accountingUpdateFields)
	if duplicateUpdateErr != nil {
		return fmt.Errorf("duplicate accounting update: %w", duplicateUpdateErr)
	}
	if duplicateUpdateAck.Fields[duplicateShipmentEventField] != "true" {
		return fmt.Errorf("duplicate accounting update was not checkpointed")
	}
	node.record("fulfillment_workflow_completed", "kept", "accounting", "order_id="+fulfillmentOrderID+" package_id="+fulfillmentPackageID)
	return nil
}

// runCASComputeWorkflow exercises CAS storage, replica recovery, CID-named
// compute, cache reuse, verification, and economics from Alice's local vantage.
// Intent: POC15 is a POC13 superset, so this preserves inherited
// storage/compute event records above the POC12 app/kernel interface without turning
// peer behavior into RPC commands, while adding sparse per-agent CAS and token
// incentive pressure. Source: DI-sinur; DI-linof; DI-manul
func (node *Node) runCASComputeWorkflow() error {
	contentBytes := production.SampleContentBytes()
	contentCID := production.ContentCID(contentBytes)
	secondContentBytes := production.SampleSecondContentBytes()
	secondContentCID := production.ContentCID(secondContentBytes)
	functionBytes := production.SampleFunctionBytes()
	inputBytes := production.SampleInputBytes()
	contextBytes := production.SampleContextBytes()
	functionCID := production.ContentCID(functionBytes)
	inputCID := production.ContentCID(inputBytes)
	contextCID := production.ContentCID(contextBytes)

	node.record("persisted_trust_history_loaded", "kept", "", "using app-local relationship snapshot before selecting CAS and compute peers")
	node.record("trust_driven_peer_choice", "kept", "bob", "pcid="+pcid.CASStorageV1+" Alice chooses Bob as primary storage peer from local trust event")
	node.record("dynamic_peer_choice_from_persisted_trust", "kept", "bob", "pcid="+pcid.CASStorageV1+" storage choice uses durable local relationship event")
	node.record("trust_driven_peer_choice", "kept", "carol", "pcid="+pcid.CIDComputeV1+" Alice chooses Carol as compute peer and Dave/Grace as verifiers")
	node.record("dynamic_peer_choice_from_persisted_trust", "kept", "carol", "pcid="+pcid.CIDComputeV1+" compute choice uses durable local relationship event")
	if err := node.runDynamicTCPTopologyWorkflow(); err != nil {
		return err
	}
	if err := node.recordPermanentDistrustAndTransitExclusionEvents(); err != nil {
		return err
	}
	if err := node.runRoutePromiseWorkflow(); err != nil {
		return err
	}
	if err := node.runMessageShapeSpecimenWorkflow(); err != nil {
		return err
	}
	node.recordMessageShapeSpecimenCoverage()
	node.recordDecentralizedMonitoringEvents()
	node.recordMixedVersionPCIDMigrationEvents()
	node.recordRunInternalRestartEvents()
	node.record("economics_price_probe", "kept", "bob", "pcid="+pcid.CASStorageV1+" Alice first offers below Bob's local storage price")
	if _, err := node.sendAndReceive("bob", map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "bob",
		"turn":          "startup",
		"promise":       "Alice promises to receive Bob's local storage price runtime adapter events for this content CID.",
		"reason":        "price discovery without treating Bob as an authority",
		"promise_about": production.PromiseStoreContent,
		"content_cid":   contentCID,
		"content_b64":   base64.StdEncoding.EncodeToString(contentBytes),
		"credit_offer":  "1",
		"units":         "1",
	}); err != nil {
		node.record("economics_price_refused", "non_commitment", "bob", "pcid="+pcid.CASStorageV1+" Bob did not promise storage at credit_offer=1")
	}
	node.record("economics_credit_offered", "kept", "bob", "pcid="+pcid.CASStorageV1+" Alice offers storage credit_offer=4 after the low-price probe")
	storeAck, err := node.sendAndReceive("bob", map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "bob",
		"turn":          "startup",
		"promise":       "Alice promises to receive Bob's bounded CAS storage event for one exact content CID.",
		"reason":        "storage should be promised and verified by exact bytes",
		"promise_about": production.PromiseStoreContent,
		"content_cid":   contentCID,
		"content_b64":   base64.StdEncoding.EncodeToString(contentBytes),
		"credit_offer":  "4",
		"units":         "1",
	})
	if err != nil {
		return fmt.Errorf("store content: %w", err)
	}
	node.record("cas_multi_object_pressure", "kept", "bob", "pcid="+pcid.CASStorageV1+" Alice sends a second independent object")
	secondStoreAck, err := node.sendAndReceive("bob", map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "bob",
		"turn":          "startup",
		"promise":       "Alice promises to receive Bob's bounded CAS storage event for a second exact content CID.",
		"reason":        "multi-object pressure should remain exact-byte promise event",
		"promise_about": production.PromiseStoreContent,
		"content_cid":   secondContentCID,
		"content_b64":   base64.StdEncoding.EncodeToString(secondContentBytes),
		"credit_offer":  "4",
		"units":         "1",
		"object_label":  "second-object",
		"token_style":   "bearer",
	})
	if err != nil {
		return fmt.Errorf("store second content: %w", err)
	}
	primaryToken := storeAck.Fields["capability_token"]
	if _, err := node.sendAndReceive("bob", map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "bob",
		"turn":          "startup",
		"promise":       "Alice promises to receive Bob's serve event for the content CID he just stored.",
		"reason":        "retrieval proves storage is not only promised but locally served",
		"promise_about": production.PromiseServeContent,
		"content_cid":   contentCID,
		"token":         primaryToken,
	}); err != nil {
		return fmt.Errorf("serve content: %w", err)
	}
	replayAck, replayErr := node.sendAndReceive("bob", map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "bob",
		"turn":          "startup",
		"promise":       "Alice promises to present the consumed serve-once token only as replay-protection test event.",
		"reason":        "a serve-once token should not create a second storage result after local redemption",
		"promise_about": production.PromiseServeContent,
		"content_cid":   contentCID,
		"token":         primaryToken,
	})
	if replayErr != nil || replayAck.Fields["token_status"] == "not_promised" {
		node.record("replay_probe_rejected", "non_commitment", "bob", "pcid="+pcid.CASStorageV1+" consumed serve-once token was not accepted as fresh event")
	} else {
		return fmt.Errorf("serve-once token replay unexpectedly produced fresh content")
	}
	node.record("network_outage_variant_selected", "kept", "bob", "pcid="+pcid.CASStorageV1+" Alice models Bob unavailable after primary retrieval")
	node.record("tcp_message_send_failed", "non_commitment", "bob", "pcid="+pcid.CASStorageV1+" simulated primary Bob outage before replica request")
	node.record("primary_storage_unavailable", "non_commitment", "bob", "pcid="+pcid.CASStorageV1+" local send failure is availability event, not broken promise event")
	node.record("replica_recovery_requested", "kept", "frank", "pcid="+pcid.CASStorageV1+" Alice asks Frank to serve Bob-replicated bytes")
	if _, err := node.sendAndReceive("frank", map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "frank",
		"turn":          "startup",
		"promise":       "Alice promises to receive Frank's replica serve event for the exact content CID.",
		"reason":        "replica recovery depends on Frank's own prior replica promise",
		"promise_about": production.PromiseServeReplicaContent,
		"content_cid":   contentCID,
		"token":         storeAck.Fields["replica_token"],
	}); err != nil {
		return fmt.Errorf("serve replica content: %w", err)
	}
	secondBearerToken := secondStoreAck.Fields["bearer_token"]
	if secondBearerToken != "" {
		node.record("agent_cas_bearer_storage_token_transferred", "kept", "frank", "pcid="+pcid.CASStorageV1+" issuer=bob content_cid="+secondContentCID)
		if _, err := node.sendAndReceive("frank", map[string]string{
			"act":           decision.ActPromise,
			"from":          node.Agent.Name,
			"to":            "frank",
			"turn":          "startup",
			"promise":       "Alice promises to transfer Bob's bearer storage token to Frank as payment for Frank's future storage work.",
			"reason":        "bearer storage tokens should be peer-held incentives rather than global authority",
			"promise_about": production.PromiseReplicaTokenLifecycle,
			"content_cid":   secondContentCID,
			"bearer_token":  secondBearerToken,
			"token_style":   "bearer",
			"token_status":  "transferred",
			"issuer_peer":   "bob",
			"redeem_peer":   "bob",
		}); err != nil {
			return fmt.Errorf("transfer bearer storage token: %w", err)
		}
	}
	missingObjectCID := production.ContentCID([]byte("poc15 missing sparse CAS object|" + node.Config.RunID))
	if _, err := node.sendAndReceive("frank", map[string]string{
		"act":                  decision.ActPromise,
		"from":                 node.Agent.Name,
		"to":                   "frank",
		"turn":                 "startup",
		"promise":              "Alice promises to treat Frank's missing-object response as a sparse CAS non-commitment, not as a broken storage promise.",
		"reason":               "sparse peer stores are expected to lack many objects",
		"promise_about":        production.PromiseServeReplicaContent,
		"content_cid":          missingObjectCID,
		"missing_object_probe": "true",
	}); err != nil {
		return fmt.Errorf("sparse CAS missing-object probe: %w", err)
	}
	node.record("economics_credit_offered", "kept", "carol", "pcid="+pcid.CIDComputeV1+" Alice offers compute credit_offer=5")
	if _, err := node.sendAndReceive("dave", map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "dave",
		"turn":          "startup",
		"promise":       "Alice promises to receive Dave's local cache status for this exact compute tuple.",
		"reason":        "cache reuse should be exact tuple event",
		"promise_about": production.PromiseLookupComputeCache,
		"function_cid":  functionCID,
		"input_cid":     inputCID,
		"context_cid":   contextCID,
	}); err != nil {
		return fmt.Errorf("compute cache miss: %w", err)
	}
	computeAck, err := node.sendAndReceive("carol", map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "carol",
		"turn":          "startup",
		"promise":       "Alice promises to receive Carol's result only as event over explicit function/input/context CIDs.",
		"reason":        "compute is a reciprocal promise over CID-named bytes",
		"promise_about": production.PromiseExecuteFunction,
		"function_cid":  functionCID,
		"function_b64":  base64.StdEncoding.EncodeToString(functionBytes),
		"input_cid":     inputCID,
		"input_b64":     base64.StdEncoding.EncodeToString(inputBytes),
		"context_cid":   contextCID,
		"context_b64":   base64.StdEncoding.EncodeToString(contextBytes),
		"credit_offer":  "5",
		"units":         "1",
	})
	if err != nil {
		return fmt.Errorf("compute request: %w", err)
	}
	if err := node.verifyComputeAckLocally(computeAck, "carol"); err != nil {
		return err
	}
	for _, verifier := range []string{"dave", "grace"} {
		fields := map[string]string{
			"act":             decision.ActPromise,
			"from":            node.Agent.Name,
			"to":              verifier,
			"turn":            "startup",
			"promise":         "Alice promises to receive a local verifier event record about Carol's compute result.",
			"reason":          "peer verification is local observation event, not global truth",
			"promise_about":   production.PromiseVerifyComputeResult,
			"result_promiser": "carol",
			"function_cid":    computeAck.Fields["function_cid"],
			"function_b64":    computeAck.Fields["function_b64"],
			"input_cid":       computeAck.Fields["input_cid"],
			"input_b64":       computeAck.Fields["input_b64"],
			"context_cid":     computeAck.Fields["context_cid"],
			"context_b64":     computeAck.Fields["context_b64"],
			"result_cid":      computeAck.Fields["result_cid"],
			"result_b64":      computeAck.Fields["result_b64"],
		}
		if verifier == "grace" {
			fields["disagreement_probe"] = "true"
			fields["promise"] = "Alice promises to receive Grace's disagreement probe as local events and resolve it locally."
		}
		if _, err := node.sendAndReceive(verifier, fields); err != nil {
			return fmt.Errorf("%s verify: %w", verifier, err)
		}
	}
	if _, err := node.sendAndReceive("dave", map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "dave",
		"turn":          "startup",
		"promise":       "Alice promises to receive Dave's cache hit after Carol's compute result was checkpointed.",
		"reason":        "cache hit should reuse exact result event",
		"promise_about": production.PromiseLookupComputeCache,
		"function_cid":  computeAck.Fields["function_cid"],
		"input_cid":     computeAck.Fields["input_cid"],
		"context_cid":   computeAck.Fields["context_cid"],
		"result_cid":    computeAck.Fields["result_cid"],
	}); err != nil {
		return fmt.Errorf("compute cache hit: %w", err)
	}
	sumFunctionBytes := production.SampleSumFunctionBytes()
	sumInputBytes := production.SampleSumInputBytes()
	node.record("compute_followup_function_requested", "kept", "dave", "pcid="+pcid.CIDComputeV1+" Alice requests a second payload-provided compute function kind=sum from a still-trusted compute peer")
	// Intent: After Alice locally observes malformed bad-result event from
	// Carol, the alternate-function coverage should follow trust and use Dave
	// rather than forcing another fresh compute promise to Carol. Source: DI-vahan
	if _, err := node.sendAndReceive("dave", map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "dave",
		"turn":          "startup",
		"promise":       "Alice promises to receive Dave's sum result only if it verifies against the payload bytes.",
		"reason":        "second compute path prevents hard-coded Fibonacci-only behavior",
		"promise_about": production.PromiseExecuteFunction,
		"function_cid":  production.ContentCID(sumFunctionBytes),
		"function_b64":  base64.StdEncoding.EncodeToString(sumFunctionBytes),
		"input_cid":     production.ContentCID(sumInputBytes),
		"input_b64":     base64.StdEncoding.EncodeToString(sumInputBytes),
		"context_cid":   contextCID,
		"context_b64":   base64.StdEncoding.EncodeToString(contextBytes),
		"credit_offer":  "5",
		"units":         "1",
	}); err != nil {
		return fmt.Errorf("sum compute request: %w", err)
	}
	if err := node.runRuntimeAdapterComputeWorkflow(functionBytes, inputBytes, contextBytes); err != nil {
		return err
	}
	return nil
}

// runRuntimeAdapterComputeWorkflow asks Peggy and Victor to keep ordinary
// cid_compute_v1 promises so runtime adapters prove useful compute work, not
// just process-interface existence.
// Intent: Alice is the requester for both exchanges; Peggy uses local WASM and
// Victor uses a stdio worker, but Alice sees both as normal PromiseGrid compute
// peers under the same pCID-owned payload contract. Source: DI-sivis
func (node *Node) runRuntimeAdapterComputeWorkflow(functionBytes, inputBytes, contextBytes []byte) error {
	targets := []struct {
		name        string
		eventPrefix string
		promiseText string
	}{
		{
			name:        "peggy",
			eventPrefix: "wasm",
			promiseText: "Alice promises to receive Peggy's WASM-kept compute result only as a cid_compute_v1 promise over exact function/input/context CIDs.",
		},
		{
			name:        "victor",
			eventPrefix: "stdio",
			promiseText: "Alice promises to receive Victor's stdio-worker-kept compute result only as a cid_compute_v1 promise over exact function/input/context CIDs.",
		},
	}
	for _, target := range targets {
		node.record(target.eventPrefix+"_compute_request_promised", "kept", target.name, "pcid="+pcid.CIDComputeV1+" Alice requests useful compute from runtime adapter")
		ack, err := node.sendAndReceive(target.name, map[string]string{
			"act":           decision.ActPromise,
			"from":          node.Agent.Name,
			"to":            target.name,
			"turn":          "startup",
			"promise":       target.promiseText,
			"reason":        "runtime adapter useful work must stay inside the existing compute pCID",
			"promise_about": production.PromiseExecuteFunction,
			"function_cid":  production.ContentCID(functionBytes),
			"function_b64":  base64.StdEncoding.EncodeToString(functionBytes),
			"input_cid":     production.ContentCID(inputBytes),
			"input_b64":     base64.StdEncoding.EncodeToString(inputBytes),
			"context_cid":   production.ContentCID(contextBytes),
			"context_b64":   base64.StdEncoding.EncodeToString(contextBytes),
			"credit_offer":  "5",
			"units":         "1",
		})
		if err != nil {
			return fmt.Errorf("%s runtime-adapter compute request: %w", target.name, err)
		}
		if err := node.verifyComputeAckLocally(ack, target.name); err != nil {
			return fmt.Errorf("%s runtime-adapter compute verify: %w", target.name, err)
		}
		node.record(target.eventPrefix+"_compute_result_verified", "kept", target.name, "pcid="+pcid.CIDComputeV1+" result_cid="+ack.Fields["result_cid"])
	}
	return nil
}

// runDynamicTCPTopologyWorkflow proves that relationship ledger changes affect
// real send/receive reachability through the local kernel.
// Intent: Dynamic TCP topology should be more than a log label: an ordinary
// promise to a removed direct peer is blocked locally, repair events can
// restore the relationship, and a later relationship promise then crosses the
// same app/kernel TCP path as other POC15 messages. Source: DI-sihuz
func (node *Node) runDynamicTCPTopologyWorkflow() error {
	target := "fulfillment"
	node.record("dynamic_tcp_topology_probe_started", "kept", target, "Alice tests whether local direct-peer changes affect actual kernel-routed sends")
	node.observeOutcome(target, relationship.OutcomeBroken)
	blockedFields := map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            target,
		"turn":          "startup",
		"promise":       "Alice promises a routine relationship observation after local trust dropped.",
		"reason":        "ordinary traffic should not cross a removed direct-peer relationship",
		"promise_about": "local_observation",
	}
	if _, err := node.sendAndReceive(target, blockedFields); err == nil {
		return fmt.Errorf("dynamic topology blocked send unexpectedly succeeded")
	}
	node.record("dynamic_tcp_topology_send_blocked", "non_commitment", target, "ordinary relationship promise was not sent because local direct TCP promise was removed")
	// Intent: One ordinary kept outcome after broken events must consume
	// caution without raising trust, giving the analyzer deterministic
	// `trust_recovery_delayed` event records independent of live-agent choices.
	// Source: DI-sihuz
	node.observeOutcome(target, relationship.OutcomeKept)
	for repairIndex := 0; repairIndex < 3; repairIndex++ {
		node.observeOutcome(target, relationship.OutcomeRepairKept)
	}
	restoredFields := map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            target,
		"turn":          "startup",
		"promise":       "Alice promises a routine relationship observation after local repair event restored direct reachability.",
		"reason":        "restored direct relationship should allow actual kernel-routed send/receive again",
		"promise_about": "local_observation",
	}
	if _, err := node.sendAndReceive(target, restoredFields); err != nil {
		return fmt.Errorf("dynamic topology restored send: %w", err)
	}
	node.record("dynamic_tcp_topology_send_succeeded", "kept", target, "restored direct TCP relationship carried an actual relationship_v1 promise")
	return nil
}

// runAdversaryWorkflow keeps malformed event and prompt-injection pressure
// inside the single promise action.
// Intent: Receivers independently reject corrupt bytes, unsupported variants,
// unknown pCIDs, and bad proofs without gaining authority over Mallory.
// Source: DI-sinur
func (node *Node) runAdversaryWorkflow() error {
	if err := node.sendUnknownProtocolPromise("grace"); err != nil {
		return fmt.Errorf("unknown protocol probe: %w", err)
	}
	unsupportedFields := map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "grace",
		"turn":          "startup",
		"promise":       "Mallory promises an unsupported storage variant to test receiver non-commitment.",
		"reason":        "unsupported variants should be not-promised rather than coerced",
		"promise_about": production.PromiseUnsupportedVariantProbe,
	}
	if _, err := node.sendAndReceive("grace", unsupportedFields); err != nil {
		node.record("promise_variant_not_promised", "non_commitment", "grace", "pcid="+pcid.CASStorageV1+" Grace did not promise unsupported storage variant")
	}
	node.shouldSuppressNonCommittedPromise("grace", unsupportedFields)
	if _, err := node.sendIdentityKeyRotationPromise("grace", "mallory-next-key", "future-poc15-identity"); err != nil {
		return fmt.Errorf("key rotation promise: %w", err)
	}
	functionBytes := production.SampleFunctionBytes()
	inputBytes := []byte("n=7")
	contextBytes := production.SampleContextBytes()
	if _, err := node.sendAndReceive("carol", map[string]string{
		"act":            decision.ActPromise,
		"from":           node.Agent.Name,
		"to":             "carol",
		"turn":           "startup",
		"promise":        "Mallory promises a compute request that intentionally pressures Carol's scarce capacity.",
		"reason":         "capacity refusal should remain a local non-commitment event record",
		"promise_about":  production.PromiseExecuteFunction,
		"capacity_probe": "true",
		"function_cid":   production.ContentCID(functionBytes),
		"function_b64":   base64.StdEncoding.EncodeToString(functionBytes),
		"input_cid":      production.ContentCID(inputBytes),
		"input_b64":      base64.StdEncoding.EncodeToString(inputBytes),
		"context_cid":    production.ContentCID(contextBytes),
		"context_b64":    base64.StdEncoding.EncodeToString(contextBytes),
		"credit_offer":   "5",
		"units":          "999",
	}); err != nil {
		return fmt.Errorf("capacity probe: %w", err)
	}
	claimedCID := production.ContentCID(production.SampleContentBytes())
	if _, err := node.sendAndReceive("grace", map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "grace",
		"turn":          "startup",
		"promise":       "Mallory promises these bytes match the claimed content CID.",
		"reason":        "adversarial corrupt-byte event should be locally rejected",
		"promise_about": production.PromisePresentStorageReport,
		"content_cid":   claimedCID,
		"content_b64":   base64.StdEncoding.EncodeToString(production.CorruptContentBytes()),
	}); err != nil {
		return fmt.Errorf("corrupt bytes offer: %w", err)
	}
	// Intent: The repair promise deliberately follows a malformed event. It is
	// narrow candidate traffic that Grace may choose to receive, but it records
	// future repair only and does not immediately restore trust. Source: DI-fijov
	if _, err := node.sendAndReceive("grace", map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            "grace",
		"turn":          "startup",
		"promise":       "Mallory promises to label future malformed storage event explicitly.",
		"reason":        "repair remains only a future promise Grace may distrust",
		"promise_about": production.PromiseLabelFutureMalformedReport,
	}); err != nil {
		return fmt.Errorf("trust repair promise: %w", err)
	}
	if err := node.sendBadProofPromise("grace"); err != nil {
		node.record("bad_proof_rejected", "malformed", "grace", "pcid="+pcid.CASStorageV1+" local kernel rejected bad proof: "+err.Error())
	}
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
	node.normalizeAutonomousPromiseFields(validDecision.Target, fields)
	if node.shouldSuppressNonCommittedPromise(validDecision.Target, fields) {
		return nil
	}
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
	frameConn, dialErr := transport.DialFrameConn(kernelAddress, sendTimeout)
	if dialErr != nil {
		return dialErr
	}
	// Intent: The app keeps one app/kernel session for both registering local
	// receive promises and sending exact envelopes. Replies are correlated by
	// parent-linked message CIDs, not by an RPC request table in payloads. Source:
	// DI-vopab
	session := transport.NewPersistentSession(
		"app-kernel:"+node.Agent.Name,
		frameConn,
		frameParentExactSHA256s,
		node.frameIsResponse,
		node.handleAppSessionFrame,
		func(eventName, outcome, detail string) {
			node.record(eventName, outcome, "kernel", detail)
		},
	)
	node.mu.Lock()
	node.appSession = session
	node.mu.Unlock()
	for _, protocolName := range node.Agent.Protocols() {
		if err := node.registerReceivePromise(ctx, session, kernelAddress, protocolName); err != nil {
			if closeErr := session.Close(); closeErr != nil {
				node.record("app_kernel_session_close_failed", "broken", "kernel", closeErr.Error())
			}
			return err
		}
	}
	return nil
}

func (node *Node) registerReceivePromise(ctx context.Context, session *transport.PersistentSession, kernelAddress, protocolName string) error {
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
	payloadBytes, _, payloadErr := protocol.MarshalKnownArrayPayload(pcid.KernelReceiveV1, fields)
	if payloadErr != nil {
		return payloadErr
	}
	envelope, envelopeErr := protocol.NewEnvelopeFromPayload(receiveCID, payloadBytes, node.Agent.Name)
	if envelopeErr != nil {
		return envelopeErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	node.emitMessageArtifact("receive_promise_sent", "kernel", pcid.KernelReceiveV1, envelopeBytes, fields)
	if writeErr := session.Send(ctx, envelopeBytes); writeErr != nil {
		return writeErr
	}
	node.record("pcid_owned_array_payload_sent", "kept", "kernel", "pcid="+pcid.KernelReceiveV1+" promise_about="+fields["promise_about"]+" exact_sha256="+protocol.HashExactBytes(envelopeBytes))
	node.record("app_receive_promise_sent", "kept", "kernel", "pcid="+protocolName+" kernel="+kernelAddress)
	node.record("app_kernel_backpressure_promised", "kept", "kernel", "pcid="+protocolName+" app promises bounded receive buffering through the local kernel")
	node.record("app_kernel_rate_limit_promised", "kept", "kernel", "pcid="+protocolName+" app promises to treat local kernel throughput as a bounded promise")
	return nil
}

func (node *Node) handleAppSessionFrame(frameBytes []byte) ([]byte, error) {
	// Intent: The app/kernel TCP stream is shared by receive promises and app
	// sends, so each inbound kernel delivery is handled as an ordinary local app
	// promise event while the persistent session owns byte correlation. Source:
	// DI-vopab
	node.activeHandlers.Add(1)
	defer node.activeHandlers.Done()
	return node.handleFrame(frameBytes)
}

func (node *Node) handleFrame(frameBytes []byte) ([]byte, error) {
	parsed, parseErr := node.parseEnvelope(frameBytes)
	if parseErr != nil {
		node.emitMessageArtifact("received_malformed", "", "unknown", frameBytes, nil)
		node.record("frame_parse_failed", "broken", "", parseErr.Error())
		node.recordMalformedFrameEvent(frameBytes, parseErr)
		return nil, parseErr
	}
	node.emitMessageArtifact("received", parsed.Fields["from"], parsed.ProtocolName, frameBytes, parsed.Fields)
	node.record("promise_envelope_validated", "kept", parsed.Fields["from"], "pcid="+parsed.ProtocolName+" exact_sha256="+parsed.ExactHash)
	node.record("tcp_message_received", "kept", parsed.Fields["from"], "pcid="+parsed.ProtocolName+" exact_sha256="+parsed.ExactHash)
	fields := parsed.Fields
	fromAgent := fields["from"]
	if isPcidOwnedArrayPayload(fields) {
		node.record("pcid_owned_array_payload_received", "kept", fromAgent, "pcid="+parsed.ProtocolName+" promise_about="+fields["promise_about"]+" exact_sha256="+parsed.ExactHash)
	}
	if node.rememberReplayEnvelope(fromAgent, parsed.ProtocolName, parsed.ExactHash) {
		return node.newAckBytes(fromAgent, "not_promised", "I promise to remember that I already saw this exact envelope and will not treat the replay as fresh promise event.", parsed.ProtocolCID, parsed.ExactHash, map[string]string{"replay_status": "not_promised"})
	}
	if !node.supportsProtocol(parsed.ProtocolName) {
		node.record("unsupported_pcid", "non_commitment", fromAgent, "no local app receive promise for "+parsed.ProtocolName)
		return node.newAckBytes(fromAgent, "not_promised", "I promise to remember that I did not promise to handle this pCID.", parsed.ProtocolCID, parsed.ExactHash, nil)
	}
	if fields["act"] != decision.ActPromise {
		node.observeOutcome(fromAgent, relationship.OutcomeMalformed)
		node.record("message_rejected", "malformed", fromAgent, "message act is not promise")
		return node.newAckBytes(fromAgent, "malformed", "I promise I rejected this non-promise message.", parsed.ProtocolCID, parsed.ExactHash, nil)
	}
	if !node.canAcceptFrom(fromAgent, fields) {
		node.record("message_not_promised", "non_commitment", fromAgent, "no current local promise to accept direct TCP exchange")
		return node.newAckBytes(fromAgent, "not_promised", "I promise to remember that I did not currently promise this direct exchange.", parsed.ProtocolCID, parsed.ExactHash, nil)
	}
	node.recordInboundPressurePromises(fromAgent, parsed.ProtocolName, fields)
	promiseID := node.rememberOutstandingPromise(fromAgent, parsed.ProtocolName, parsed.ExactHash, fields)
	if resourceErr := node.checkIncomingResourcePromise(fields); resourceErr != nil {
		node.resolveOutstandingPromise(promiseID, promiseStatusBroken, resourceErr.Error())
		node.observeOutcome(fromAgent, relationship.OutcomeBroken)
		node.applyBrokenPromiseCost(fromAgent, fields, resourceErr.Error())
		node.record("resource_promise_rejected", "broken", fromAgent, resourceErr.Error())
		return node.newAckBytes(fromAgent, "broken", "I promise I rejected this resource promise because local checks failed.", parsed.ProtocolCID, parsed.ExactHash, nil)
	}
	handlerResult, handlerErr := node.handleProtocolPromise(parsed)
	if handlerErr != nil {
		node.resolveOutstandingPromise(promiseID, promiseStatusBroken, handlerErr.Error())
		node.observeOutcome(fromAgent, relationship.OutcomeBroken)
		node.applyBrokenPromiseCost(fromAgent, fields, handlerErr.Error())
		node.record("protocol_handler_rejected", "broken", fromAgent, handlerErr.Error())
		return node.newAckBytes(fromAgent, "broken", "I promise I rejected this protocol promise because local app checks failed.", parsed.ProtocolCID, parsed.ExactHash, nil)
	}
	ackFields := handlerResult.Fields
	ackBytes := handlerResult.AckBytes
	if len(ackBytes) > 0 {
		ackMessage, parseErr := node.parseEnvelope(ackBytes)
		if parseErr != nil {
			node.resolveOutstandingPromise(promiseID, promiseStatusBroken, parseErr.Error())
			node.observeOutcome(fromAgent, relationship.OutcomeBroken)
			node.applyBrokenPromiseCost(fromAgent, fields, parseErr.Error())
			node.record("protocol_handler_rejected", "broken", fromAgent, parseErr.Error())
			return node.newAckBytes(fromAgent, "broken", "I promise I rejected this protocol promise because local app checks failed.", parsed.ProtocolCID, parsed.ExactHash, nil)
		}
		ackFields = ackMessage.Fields
		node.emitMessageArtifact("ack_sent", fromAgent, ackMessage.ProtocolName, ackBytes, ackFields)
	}
	acceptedAsCandidate := isLinkDiscoveryPromise(fields) && !node.canAccept(fromAgent)
	if eventUpdatesTrust(ackFields) {
		trustOutcome := outcomeForPromise(fields)
		node.resolveOutstandingPromise(promiseID, promiseStatusFromOutcome(trustOutcome), "accepted inbound promise")
		node.observeOutcome(fromAgent, trustOutcome)
	} else {
		node.resolveOutstandingPromise(promiseID, promiseStatusForNonTrustingEvent(ackFields), "non-mutating inbound event recorded without trust change")
	}
	eventName := "message_received"
	if acceptedAsCandidate {
		eventName = "candidate_message_received"
	}
	node.record(eventName, "kept", fromAgent, "received "+parsed.ProtocolName+" signed promise exact_sha256="+parsed.ExactHash)
	if len(ackBytes) > 0 {
		return ackBytes, nil
	}
	return node.newAckBytes(fromAgent, "kept", "I promise I received and recorded your signed promise message.", parsed.ProtocolCID, parsed.ExactHash, ackFields)
}

func (node *Node) newAckBytes(target, outcome, promiseText string, protocolCID protocol.ProtocolCID, parentExactSHA256 string, extraFields map[string]string) ([]byte, error) {
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
	// Intent: Handler event records may include copied request fields, but the ACK is
	// still a fresh promise by this local agent to the original sender. Source:
	// DI-gahuh
	ackFields["act"] = decision.ActPromise
	ackFields["from"] = node.Agent.Name
	ackFields["to"] = target
	ackFields["outcome"] = outcome
	ackFields["promise"] = promiseText
	ackFields["reason"] = "transport acknowledgement expressed as local promise content"
	if protocolCID.Equal(node.Protocols.MustCID(pcid.IdentityKeyV1)) && extraFields["promise_about"] == production.PromiseRotateSigningKey {
		// Intent: identity_key_v1 is the first POC15 protocol whose request and
		// ACK payloads are pCID-owned arrays rather than older generic maps.
		// Source: DI-vipih
		payloadBytes, payloadErr := protocol.MarshalIdentityKeyRotationAckPayload(protocol.IdentityKeyRotationAckPayload{
			Promiser:      node.Agent.Name,
			Promisee:      target,
			Outcome:       outcome,
			PromiseText:   promiseText,
			NewKeyLabel:   extraFields["new_key_label"],
			RotationScope: extraFields["rotation_scope"],
		})
		if payloadErr != nil {
			node.record("ack_sign_failed", "broken", target, payloadErr.Error())
			return nil, payloadErr
		}
		ack, ackErr := protocol.NewEnvelopeFromPayloadWithParents(protocolCID, payloadBytes, []string{parentExactSHA256}, node.Agent.Name)
		if ackErr != nil {
			node.record("ack_sign_failed", "broken", target, ackErr.Error())
			return nil, ackErr
		}
		ackBytes, bytesErr := ack.Bytes()
		if bytesErr != nil {
			node.record("ack_bytes_failed", "broken", target, bytesErr.Error())
			return nil, bytesErr
		}
		node.emitMessageArtifact("ack_sent", target, pcid.IdentityKeyV1, ackBytes, ackFields)
		node.record("pcid_owned_array_ack_sent", "kept", target, "pcid="+pcid.IdentityKeyV1+" promise_about="+ackFields["promise_about"])
		return ackBytes, nil
	}
	protocolName, protocolKnown := node.Protocols.Name(protocolCID)
	if protocolKnown {
		if ackFields["promise_about"] == "" && (protocolName == pcid.RelationshipV1 || protocolName == pcid.KernelReceiveV1) {
			ackFields["promise_about"] = "local_observation"
		}
		ack, arrayPayload, ackErr := node.buildEnvelopeFromFieldsWithParents(protocolName, protocolCID, ackFields, []string{parentExactSHA256})
		if ackErr == nil && arrayPayload {
			// Intent: Migrated pCIDs must keep ACKs in the same pCID-owned
			// positional payload family as their requests, rather than falling
			// back to a generic map acknowledgement. Source: DI-gahuh;
			// DI-dirat
			ackBytes, bytesErr := ack.Bytes()
			if bytesErr != nil {
				node.record("ack_bytes_failed", "broken", target, bytesErr.Error())
				return nil, bytesErr
			}
			node.emitMessageArtifact("ack_sent", target, protocolName, ackBytes, ackFields)
			node.record("pcid_owned_array_ack_sent", "kept", target, "pcid="+protocolName+" promise_about="+ackFields["promise_about"])
			return ackBytes, nil
		}
	}
	ack, ackErr := protocol.NewEnvelopeWithParents(protocolCID, ackFields, []string{parentExactSHA256}, node.Agent.Name)
	if ackErr != nil {
		node.record("ack_sign_failed", "broken", target, ackErr.Error())
		return nil, ackErr
	}
	ackBytes, bytesErr := ack.Bytes()
	if bytesErr != nil {
		node.record("ack_bytes_failed", "broken", target, bytesErr.Error())
		return nil, bytesErr
	}
	node.emitMessageArtifact("ack_sent", target, firstNonEmpty(protocolName, "unknown"), ackBytes, ackFields)
	return ackBytes, nil
}

func (node *Node) send(target string, fields map[string]string) error {
	_, err := node.sendAndReceive(target, fields)
	return err
}

// nextExchangeID assigns a sender-local event identity to each outbound
// promise. Intent: Semantic duplicates can still be recognized by pCID-owned
// fields, while exact-byte replay protection can reject a re-sent envelope whose
// bytes are literally identical. Source: DI-sunuf
func (node *Node) nextExchangeID(target, protocolName string) string {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.exchangeCounter++
	return fmt.Sprintf("%s-%s-%s-%06d", node.Agent.Name, target, protocolName, node.exchangeCounter)
}

// sendAndReceive performs one signed promise exchange and returns the receiver's
// ACK event record to the local caller.
// Intent: The fulfillment workflow needs concrete address, weight, label, and
// accounting event records from pCID handlers while the app sends only through its
// local kernel rather than dialing peer app processes directly; each outbound
// promise is also journaled before its ACK can affect trust. Source: DI-galin;
// DI-vujob
func (node *Node) sendAndReceive(target string, fields map[string]string) (parsedMessage, error) {
	if !node.canDialTarget(target, fields) {
		return parsedMessage{}, fmt.Errorf("no local TCP promise to %s", target)
	}
	protocolName, protocolCID := node.protocolForFields(fields)
	fields["protocol"] = protocolName
	fields["exchange_id"] = node.nextExchangeID(target, protocolName)
	envelope, arrayPayload, envelopeErr := node.buildEnvelopeFromFields(protocolName, protocolCID, fields)
	if envelopeErr != nil {
		return parsedMessage{}, envelopeErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return parsedMessage{}, bytesErr
	}
	exactHash := protocol.HashExactBytes(envelopeBytes)
	node.recordOutboundPressurePromises(target, protocolName, fields)
	promiseID := node.rememberOutstandingPromise(target, protocolName, exactHash, fields)
	session := node.appKernelSession()
	if session == nil {
		node.resolveOutstandingPromise(promiseID, promiseStatusLocalFailure, "missing local persistent app-kernel session")
		return parsedMessage{}, fmt.Errorf("no local persistent app-kernel session for app %s", node.Agent.Name)
	}
	node.emitMessageArtifact("sent", target, protocolName, envelopeBytes, fields)
	if arrayPayload {
		node.record("pcid_owned_array_payload_sent", "kept", target, "pcid="+protocolName+" promise_about="+fields["promise_about"]+" exact_sha256="+exactHash)
	}
	node.record("tcp_message_sent", "kept", target, "pcid="+protocolName+" exact_sha256="+exactHash)
	roundTripCtx, cancel := context.WithTimeout(context.Background(), sendTimeout)
	defer cancel()
	ackBytes, readErr := session.RoundTrip(roundTripCtx, exactHash, envelopeBytes)
	if readErr != nil {
		node.resolveOutstandingPromise(promiseID, promiseStatusLocalFailure, readErr.Error())
		return parsedMessage{}, readErr
	}
	ackMessage, parseErr := node.parseEnvelope(ackBytes)
	if parseErr != nil {
		node.emitMessageArtifact("ack_received_malformed", target, "unknown", ackBytes, nil)
		node.resolveOutstandingPromise(promiseID, promiseStatusLocalFailure, parseErr.Error())
		return parsedMessage{}, parseErr
	}
	node.emitMessageArtifact("ack_received", target, ackMessage.ProtocolName, ackBytes, ackMessage.Fields)
	ackFields := ackMessage.Fields
	if ackFields["outcome"] != "kept" {
		if ackFields["outcome"] == "not_promised" || ackFields["outcome"] == string(relationship.OutcomeNonCommitment) {
			node.rememberNonCommitment(target, protocolName, fields, "ack outcome "+ackFields["outcome"])
		}
		ackOutcome, _ := outcomeForSendError(ackOutcomeError{outcome: ackFields["outcome"]})
		node.resolveOutstandingPromise(promiseID, promiseStatusFromOutcome(ackOutcome), "ack outcome "+ackFields["outcome"])
		return parsedMessage{}, ackOutcomeError{outcome: ackFields["outcome"]}
	}
	if isPcidOwnedArrayPayload(ackFields) {
		node.record("pcid_owned_array_ack_received", "kept", target, "pcid="+ackMessage.ProtocolName+" promise_about="+ackFields["promise_about"]+" exact_sha256="+ackMessage.ExactHash)
	}
	node.recordAckEvent(target, ackMessage)
	if eventUpdatesTrust(ackFields) {
		trustOutcome := outcomeForPromise(fields)
		node.resolveOutstandingPromise(promiseID, promiseStatusFromOutcome(trustOutcome), "ack kept")
		node.observeOutcome(target, trustOutcome)
	} else {
		node.resolveOutstandingPromise(promiseID, promiseStatusForNonTrustingEvent(ackFields), "non-mutating ack event recorded without trust change")
	}
	return ackMessage, nil
}

// buildEnvelopeFromFields preserves legacy map payloads for pCIDs that have not
// been migrated, while encoding migrated storage/compute protocols as
// pCID-owned CBOR arrays.
// Intent: The kernel can still route by compatibility fields locally, but the
// wire payload for migrated pCIDs is no longer a universal named map.
// Source: DI-gahuh; DI-pusak
func (node *Node) buildEnvelopeFromFields(protocolName string, protocolCID protocol.ProtocolCID, fields map[string]string) (protocol.Envelope, bool, error) {
	return node.buildEnvelopeFromFieldsWithParents(protocolName, protocolCID, fields, nil)
}

func (node *Node) buildEnvelopeFromFieldsWithParents(protocolName string, protocolCID protocol.ProtocolCID, fields map[string]string, parentExactSHA256s []string) (protocol.Envelope, bool, error) {
	payloadBytes, arrayPayload, payloadErr := protocol.MarshalKnownArrayPayload(protocolName, fields)
	if payloadErr != nil {
		return protocol.Envelope{}, false, payloadErr
	}
	mergedParents := mergeParentExactSHA256s(parentExactSHA256s, envelopeParentExactSHA256sForFields(protocolName, fields))
	if arrayPayload {
		// Intent: ACKs and route-carried messages both use envelope parent links,
		// but ACK parents come from the exact request CID while route parents can
		// also be pCID-defined fields. Merge them without adding RPC IDs. Source:
		// DI-vopab
		envelope, envelopeErr := protocol.NewEnvelopeFromPayloadWithParents(protocolCID, payloadBytes, mergedParents, node.Agent.Name)
		return envelope, true, envelopeErr
	}
	envelope, envelopeErr := protocol.NewEnvelopeWithParents(protocolCID, fields, mergedParents, node.Agent.Name)
	return envelope, false, envelopeErr
}

func mergeParentExactSHA256s(parentLists ...[]string) []string {
	seenParents := make(map[string]bool)
	mergedParents := make([]string, 0)
	for _, parentList := range parentLists {
		for _, parentExactSHA256 := range parentList {
			parentExactSHA256 = strings.TrimSpace(parentExactSHA256)
			if parentExactSHA256 == "" || seenParents[parentExactSHA256] {
				continue
			}
			seenParents[parentExactSHA256] = true
			mergedParents = append(mergedParents, parentExactSHA256)
		}
	}
	return mergedParents
}

func envelopeParentExactSHA256sForFields(protocolName string, fields map[string]string) []string {
	if protocolName != pcid.RouteV1 {
		return nil
	}
	parentExactSHA256 := strings.TrimSpace(fields["envelope_parent_exact_sha256"])
	if parentExactSHA256 == "" {
		return nil
	}
	return []string{parentExactSHA256}
}

// sendUnknownProtocolPromise sends a syntactically valid envelope whose pCID is
// not in the local registry.
// Intent: Unknown protocol handling is a kernel/app non-commitment path, not an
// opportunity for a receiver to coerce semantics from an unrecognized pCID.
// Source: DI-sinur
func (node *Node) sendUnknownProtocolPromise(target string) error {
	fields := map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            target,
		"turn":          "startup",
		"promise":       "Mallory promises an envelope under an unknown protocol CID.",
		"reason":        "unknown pCID should produce local non-commitment",
		"promise_about": production.PromiseUnsupportedVariantProbe,
	}
	unknownCID := protocol.NewProtocolCID([]byte("poc15 unknown mallory probe protocol"))
	envelope, envelopeErr := protocol.NewEnvelope(unknownCID, fields, node.Agent.Name)
	if envelopeErr != nil {
		return envelopeErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	ackMessage, sendErr := node.sendRawEnvelope(target, "unknown_pcid_v1", envelopeBytes)
	if sendErr != nil {
		return sendErr
	}
	if ackMessage.Fields["outcome"] != "kept" {
		node.record("unknown_pcid_not_promised", "non_commitment", target, "unknown protocol CID was not promised by local receiver")
		return nil
	}
	node.record("unknown_pcid_not_promised", "non_commitment", target, "unexpected kept ACK for unknown protocol was treated as suspect event")
	return nil
}

// sendBadProofPromise corrupts one otherwise valid envelope after signing.
// Intent: Bad proof rejection must be tested at the exact-byte envelope layer
// while still being observed as a local promise outcome, not an imposed command.
// Source: DI-sinur; DI-punib
func (node *Node) sendBadProofPromise(target string) error {
	fields := map[string]string{
		"act":           decision.ActPromise,
		"from":          node.Agent.Name,
		"to":            target,
		"turn":          "startup",
		"promise":       "Mallory promises this intentionally corrupted signature is valid.",
		"reason":        "bad proof should fail before payload semantics are trusted",
		"promise_about": production.PromisePresentStorageReport,
	}
	payloadBytes, _, payloadErr := protocol.MarshalKnownArrayPayload(pcid.CASStorageV1, fields)
	if payloadErr != nil {
		return payloadErr
	}
	envelope, envelopeErr := protocol.NewEnvelopeFromPayload(node.Protocols.MustCID(pcid.CASStorageV1), payloadBytes, node.Agent.Name)
	if envelopeErr != nil {
		return envelopeErr
	}
	mutatedFields := copyStringMap(fields)
	mutatedFields["reason"] = "bad proof should fail because this parseable payload was changed after signing"
	mutatedPayload, _, mutatedPayloadErr := protocol.MarshalKnownArrayPayload(pcid.CASStorageV1, mutatedFields)
	if mutatedPayloadErr != nil {
		return mutatedPayloadErr
	}
	// Intent: Mutate the signable payload after proof generation so the receiver
	// can still parse the pCID and promiser fields but must reject the exact
	// envelope proof as a malformed event. Source: DI-sunuf
	envelope.Payload = mutatedPayload
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	node.record("bad_proof_sent", "kept", target, "pcid="+pcid.CASStorageV1+" sent intentionally corrupted signature")
	ackMessage, sendErr := node.sendRawEnvelope(target, pcid.CASStorageV1, envelopeBytes)
	if sendErr != nil {
		return sendErr
	}
	if ackMessage.Fields["outcome"] == "malformed" {
		node.record("bad_proof_rejected", "malformed", target, "pcid="+pcid.CASStorageV1+" receiver returned malformed outcome for intentionally corrupted proof")
		return nil
	}
	return fmt.Errorf("bad proof unexpectedly produced a readable ACK")
}

func (node *Node) sendRawEnvelope(target, protocolName string, envelopeBytes []byte) (parsedMessage, error) {
	message, _, err := node.sendRawEnvelopeBytes(target, protocolName, envelopeBytes)
	return message, err
}

// sendIdentityKeyRotationPromise sends identity_key_v1 using a pCID-owned CBOR
// array payload instead of the earlier generic map scaffold.
// Intent: Key rotation belongs to identity/key protocol semantics, and this is
// the first POC15 payload migrated under `DI-vipih`. Source: DI-vipih
func (node *Node) sendIdentityKeyRotationPromise(target, newKeyLabel, rotationScope string) (parsedMessage, error) {
	payloadBytes, payloadErr := protocol.MarshalIdentityKeyRotationPayload(protocol.IdentityKeyRotationPayload{
		Promiser:      node.Agent.Name,
		Promisee:      target,
		NewKeyLabel:   newKeyLabel,
		RotationScope: rotationScope,
	})
	if payloadErr != nil {
		return parsedMessage{}, payloadErr
	}
	envelope, envelopeErr := protocol.NewEnvelopeFromPayload(node.Protocols.MustCID(pcid.IdentityKeyV1), payloadBytes, node.Agent.Name)
	if envelopeErr != nil {
		return parsedMessage{}, envelopeErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return parsedMessage{}, bytesErr
	}
	message, sendErr := node.sendRawEnvelope(target, pcid.IdentityKeyV1, envelopeBytes)
	if sendErr != nil {
		return parsedMessage{}, sendErr
	}
	if message.Fields["outcome"] != "kept" {
		return parsedMessage{}, ackOutcomeError{outcome: message.Fields["outcome"]}
	}
	node.record("identity_key_rotation_ack_received", "kept", target, "pcid="+pcid.IdentityKeyV1+" new_key_label="+newKeyLabel)
	return message, nil
}

// sendRawEnvelopeBytes forwards an already-signed envelope and returns both the
// parsed ACK and the exact ACK bytes.
// Intent: POC15 stdio workers must observe the peer's original ACK envelope
// through stdin/stdout rather than a locally re-signed adapter copy. Source:
// DI-linof
func (node *Node) sendRawEnvelopeBytes(target, protocolName string, envelopeBytes []byte) (parsedMessage, []byte, error) {
	if !node.canDialTarget(target, nil) {
		return parsedMessage{}, nil, fmt.Errorf("no local TCP promise to %s", target)
	}
	session := node.appKernelSession()
	if session == nil {
		return parsedMessage{}, nil, fmt.Errorf("no local persistent app-kernel session for app %s", node.Agent.Name)
	}
	if sentMessage, parseErr := node.parseEnvelope(envelopeBytes); parseErr == nil {
		node.emitMessageArtifact("sent", target, sentMessage.ProtocolName, envelopeBytes, sentMessage.Fields)
		if isPcidOwnedArrayPayload(sentMessage.Fields) {
			node.record("pcid_owned_array_payload_sent", "kept", target, "pcid="+sentMessage.ProtocolName+" promise_about="+sentMessage.Fields["promise_about"]+" exact_sha256="+sentMessage.ExactHash)
		}
	} else {
		node.emitMessageArtifact("sent_malformed", target, protocolName, envelopeBytes, nil)
	}
	node.record("tcp_message_sent", "kept", target, "pcid="+protocolName+" exact_sha256="+protocol.HashExactBytes(envelopeBytes))
	roundTripCtx, cancel := context.WithTimeout(context.Background(), sendTimeout)
	defer cancel()
	ackBytes, readErr := session.RoundTrip(roundTripCtx, protocol.HashExactBytes(envelopeBytes), envelopeBytes)
	if readErr != nil {
		return parsedMessage{}, nil, readErr
	}
	message, parseErr := node.parseEnvelope(ackBytes)
	if parseErr != nil {
		node.emitMessageArtifact("ack_received_malformed", target, "unknown", ackBytes, nil)
		return parsedMessage{}, nil, parseErr
	}
	node.emitMessageArtifact("ack_received", target, message.ProtocolName, ackBytes, message.Fields)
	if isPcidOwnedArrayPayload(message.Fields) {
		node.record("pcid_owned_array_ack_received", "kept", target, "pcid="+message.ProtocolName+" promise_about="+message.Fields["promise_about"]+" exact_sha256="+message.ExactHash)
	}
	return message, ackBytes, nil
}

func (node *Node) parseEnvelope(frameBytes []byte) (parsedMessage, error) {
	envelope, parseErr := protocol.ParseEnvelope(frameBytes)
	if parseErr != nil {
		return parsedMessage{}, parseErr
	}
	if verifyErr := protocol.VerifyEnvelope(envelope); verifyErr != nil {
		return parsedMessage{}, verifyErr
	}
	protocolName, known := node.Protocols.Name(envelope.ProtocolCID)
	if !known {
		protocolName = "unknown:" + envelope.ProtocolCID.String()
	}
	fields, fieldsErr := fieldsForEnvelopeProtocol(envelope, protocolName, known)
	if fieldsErr != nil {
		return parsedMessage{}, fieldsErr
	}
	if fields["from"] == "" {
		return parsedMessage{}, fmt.Errorf("payload from field is required")
	}
	// Intent: pCID-owned array decoders expose only local compatibility fields;
	// the parser attaches the locally known protocol name for handlers that still
	// compare event records across legacy and migrated payloads. Source: DI-gahuh
	fields["protocol"] = protocolName
	if len(envelope.ParentExactSHA256s) > 0 {
		fields["envelope_parent_exact_sha256"] = envelope.ParentExactSHA256s[0]
		if fields["parent_exact_sha256"] == "" {
			fields["parent_exact_sha256"] = envelope.ParentExactSHA256s[0]
		}
		if fields["parent_link_location"] == "" {
			fields["parent_link_location"] = "envelope"
		}
	}
	return parsedMessage{
		Fields:             fields,
		ExactHash:          protocol.HashExactBytes(frameBytes),
		RawBytes:           append([]byte(nil), frameBytes...),
		ProtocolCID:        envelope.ProtocolCID,
		ProtocolName:       protocolName,
		ParentExactSHA256s: append([]string(nil), envelope.ParentExactSHA256s...),
	}, nil
}

func isPcidOwnedArrayPayload(fields map[string]string) bool {
	return fields["payload_protocol"] != ""
}

func fieldsForEnvelopeProtocol(envelope protocol.Envelope, protocolName string, known bool) (map[string]string, error) {
	// Intent: Slot 0 pCID owns slot 1 decoding; blindly trying compatible array
	// grammars can turn a relationship promise into a kernel-receive promise.
	// Source: DI-pusak
	if known {
		return protocol.PayloadFieldsForProtocolName(protocolName, envelope.Payload)
	}
	return envelope.PayloadFields()
}

func frameParentExactSHA256s(frameBytes []byte) ([]string, error) {
	envelope, parseErr := protocol.ParseEnvelope(frameBytes)
	if parseErr != nil {
		return nil, parseErr
	}
	return append([]string(nil), envelope.ParentExactSHA256s...), nil
}

func (node *Node) frameIsResponse(frameBytes []byte) (bool, error) {
	// Intent: Persistent app/kernel sessions must not treat an unmatched ACK as a
	// fresh peer promise. ACK-ness is pCID-decoded from local compatibility fields
	// while parent links remain the actual correlation key. Source: DI-vopab
	envelope, parseErr := protocol.ParseEnvelope(frameBytes)
	if parseErr != nil {
		return false, parseErr
	}
	protocolName, known := node.Protocols.Name(envelope.ProtocolCID)
	if !known {
		protocolName = "unknown:" + envelope.ProtocolCID.String()
	}
	fields, fieldsErr := fieldsForEnvelopeProtocol(envelope, protocolName, known)
	if fieldsErr != nil {
		return false, fieldsErr
	}
	return strings.TrimSpace(fields["outcome"]) != "", nil
}

func (node *Node) appKernelSession() *transport.PersistentSession {
	node.mu.Lock()
	defer node.mu.Unlock()
	return node.appSession
}

// recordMalformedFrameEvent extracts the promiser from a parseable but
// unverifiable envelope when possible. Intent: A bad proof should reduce only
// the local observer's trust in the identifiable promiser; malformed random
// bytes still remain unattributed parse-failure events. Source: DI-sunuf
func (node *Node) recordMalformedFrameEvent(frameBytes []byte, cause error) {
	envelope, parseErr := protocol.ParseEnvelope(frameBytes)
	if parseErr != nil {
		return
	}
	protocolName, known := node.Protocols.Name(envelope.ProtocolCID)
	if !known {
		protocolName = "unknown:" + envelope.ProtocolCID.String()
	}
	fields, fieldsErr := fieldsForEnvelopeProtocol(envelope, protocolName, known)
	if fieldsErr != nil {
		return
	}
	fromAgent := fields["from"]
	if fromAgent == "" {
		return
	}
	node.record("malformed_proof_observed", "malformed", fromAgent, "pcid="+protocolName+" exact_sha256="+protocol.HashExactBytes(frameBytes)+" error="+cause.Error())
	node.observeOutcome(fromAgent, relationship.OutcomeMalformed)
}

// protocolForFields chooses the pCID for an outbound promise from protocol or
// promise_about payload meaning. The pCID is still protocol identity, not a
// per-message-type selector; promise_about remains inside the pCID-owned body.
// Source: DI-bikit
func (node *Node) protocolForFields(fields map[string]string) (string, protocol.ProtocolCID) {
	if protocolName := fields["protocol"]; node.Protocols.Known(protocolName) {
		return protocolName, node.Protocols.MustCID(protocolName)
	}
	switch fields["promise_about"] {
	case production.PromiseWeighPackage:
		return pcid.PostalScaleV1, node.Protocols.MustCID(pcid.PostalScaleV1)
	case production.PromiseAddressLookup, production.PromiseShipmentUpdate:
		return pcid.AccountingV1, node.Protocols.MustCID(pcid.AccountingV1)
	case production.PromisePrintLabel:
		return pcid.UPSLabelV1, node.Protocols.MustCID(pcid.UPSLabelV1)
	case production.PromiseIssuePrintCapability, production.PromiseRedeemPrintCapability:
		return pcid.PrinterPortV1, node.Protocols.MustCID(pcid.PrinterPortV1)
	case production.PromiseStoreContent, production.PromiseServeContent, production.PromiseReplicateContent, production.PromiseServeReplicaContent, production.PromiseReplicaTokenLifecycle, production.PromisePresentStorageReport, production.PromiseLabelFutureMalformedReport, production.PromiseUnsupportedVariantProbe:
		return pcid.CASStorageV1, node.Protocols.MustCID(pcid.CASStorageV1)
	case production.PromiseExecuteFunction, production.PromiseLookupComputeCache, production.PromiseProvideComputeContext, production.PromiseVerifyComputeResult:
		return pcid.CIDComputeV1, node.Protocols.MustCID(pcid.CIDComputeV1)
	case production.PromiseRotateSigningKey:
		return pcid.IdentityKeyV1, node.Protocols.MustCID(pcid.IdentityKeyV1)
	default:
		return pcid.RelationshipV1, node.Protocols.MustCID(pcid.RelationshipV1)
	}
}

// normalizeAutonomousPromiseFields keeps live-agent free-form promises on a
// pCID whose payload they can actually satisfy.
// Intent: Live LLM autonomy should not regress protocol validity by selecting a
// concrete storage/compute pCID while omitting that pCID's required
// payload fields; unsupported capacity, pricing, or observation language remains
// a valid relationship promise instead of becoming a broken protocol exchange.
// Source: DI-punib
func (node *Node) normalizeAutonomousPromiseFields(target string, fields map[string]string) {
	if fields["promise_about"] == "" {
		fields["promise_about"] = "local_observation"
	}
	requestedProtocol := fields["protocol"]
	if requestedProtocol == "" || requestedProtocol == pcid.RelationshipV1 {
		return
	}
	if autonomousProtocolPayloadSupported(requestedProtocol, fields) {
		return
	}
	originalPromiseAbout := fields["promise_about"]
	fields["original_protocol"] = requestedProtocol
	fields["original_promise_about"] = originalPromiseAbout
	fields["protocol"] = pcid.RelationshipV1
	fields["promise_about"] = "local_observation"
	node.record("promise_reframed_for_pcid_fit", "kept", target, "original_protocol="+requestedProtocol+" original_promise_about="+originalPromiseAbout+" reframed_protocol="+pcid.RelationshipV1)
}

// autonomousProtocolPayloadSupported recognizes the concrete payloads that an
// LLM turn may send directly under non-relationship pCIDs.
// Intent: Scripted workflows still exercise full CAS, compute, and device
// protocols; live-agent turns must either supply the required pCID-owned fields
// or stay at relationship level. Source: DI-punib
func autonomousProtocolPayloadSupported(protocolName string, fields map[string]string) bool {
	switch protocolName {
	case pcid.CASStorageV1:
		return casPayloadShapePresent(fields)
	case pcid.CIDComputeV1:
		return computePayloadShapePresent(fields)
	case pcid.IdentityKeyV1:
		// Intent: Live LLM generic-map turns must not synthesize identity_key_v1
		// payloads; POC15 only sends this pCID through explicit array encoders.
		// Source: DI-vipih
		return false
	case pcid.PostalScaleV1, pcid.AccountingV1, pcid.UPSLabelV1, pcid.PrinterPortV1:
		return productionPayloadShapePresent(protocolName, fields)
	default:
		return true
	}
}

func casPayloadShapePresent(fields map[string]string) bool {
	switch fields["promise_about"] {
	case production.PromiseStoreContent, production.PromiseReplicateContent, production.PromisePresentStorageReport:
		return fields["content_cid"] != "" && fields["content_b64"] != ""
	case production.PromiseServeContent, production.PromiseServeReplicaContent:
		return fields["content_cid"] != "" && fields["token"] != ""
	case production.PromiseLabelFutureMalformedReport, production.PromiseUnsupportedVariantProbe, production.PromiseReplicaTokenLifecycle:
		return true
	default:
		return false
	}
}

func computePayloadShapePresent(fields map[string]string) bool {
	switch fields["promise_about"] {
	case production.PromiseExecuteFunction:
		return fields["function_cid"] != "" && fields["function_b64"] != "" &&
			fields["input_cid"] != "" && fields["input_b64"] != "" &&
			fields["context_cid"] != "" && fields["context_b64"] != ""
	case production.PromiseLookupComputeCache:
		return fields["function_cid"] != "" && fields["input_cid"] != "" && fields["context_cid"] != ""
	case production.PromiseVerifyComputeResult:
		return fields["function_cid"] != "" && fields["function_b64"] != "" &&
			fields["input_cid"] != "" && fields["input_b64"] != "" &&
			fields["context_cid"] != "" && fields["context_b64"] != "" &&
			fields["result_cid"] != "" && fields["result_b64"] != ""
	default:
		return false
	}
}

func productionPayloadShapePresent(protocolName string, fields map[string]string) bool {
	switch protocolName {
	case pcid.PostalScaleV1:
		return fields["promise_about"] == production.PromiseWeighPackage && fields["package_id"] != ""
	case pcid.AccountingV1:
		return fields["promise_about"] == production.PromiseAddressLookup && fields["order_id"] != "" ||
			fields["promise_about"] == production.PromiseShipmentUpdate && fields["order_id"] != "" && fields["tracking_number"] != ""
	case pcid.UPSLabelV1:
		return fields["promise_about"] == production.PromisePrintLabel && fields["package_id"] != "" && fields["shipping_address"] != ""
	case pcid.PrinterPortV1:
		return fields["promise_about"] == production.PromiseIssuePrintCapability && fields["print_capability_issuee"] != "" ||
			fields["promise_about"] == production.PromiseRedeemPrintCapability && fields["print_capability_token_id"] != ""
	default:
		return false
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

func (node *Node) handleProtocolPromise(message parsedMessage) (protocolHandlerResult, error) {
	switch message.ProtocolName {
	case pcid.RelationshipV1:
		return protocolHandlerResult{}, nil
	case pcid.PostalScaleV1:
		fields, err := node.handlePostalScalePromise(message.Fields)
		return protocolHandlerResult{Fields: fields}, err
	case pcid.UPSLabelV1:
		fields, err := node.handleUPSLabelPromise(message.Fields)
		return protocolHandlerResult{Fields: fields}, err
	case pcid.AccountingV1:
		fields, err := node.handleAccountingPromise(message.Fields)
		return protocolHandlerResult{Fields: fields}, err
	case pcid.PrinterPortV1:
		fields, err := node.handlePrinterPortPromise(message.Fields)
		return protocolHandlerResult{Fields: fields}, err
	case pcid.CASStorageV1:
		fields, err := node.handleCASStoragePromise(message.Fields)
		return protocolHandlerResult{Fields: fields}, err
	case pcid.CIDComputeV1:
		return node.handleCIDComputePromise(message)
	case pcid.IdentityKeyV1:
		fields, err := node.handleIdentityKeyPromise(message.Fields)
		return protocolHandlerResult{Fields: fields}, err
	case pcid.RouteV1:
		fields, err := node.handleRoutePromise(message)
		return protocolHandlerResult{Fields: fields}, err
	default:
		return protocolHandlerResult{}, fmt.Errorf("unsupported protocol %s", message.ProtocolName)
	}
}

func (node *Node) handlePostalScalePromise(fields map[string]string) (map[string]string, error) {
	if node.Agent.Kind != "postal_scale" {
		return nil, nil
	}
	if fields["promise_about"] != production.PromiseWeighPackage {
		return nil, fmt.Errorf("postal scale cannot handle promise_about=%q", fields["promise_about"])
	}
	packageID := firstStringField(fields, "package_id", "package_id")
	weightOunces, err := production.WeightForPackage(packageID)
	if err != nil {
		return nil, err
	}
	node.record("package_weighed", "kept", fields["from"], fmt.Sprintf("package_id=%s weight_ounces=%d", packageID, weightOunces))
	return map[string]string{
		"promise_about": production.PromiseWeighPackage,
		"package_id":    packageID,
		"weight_ounces": strconv.Itoa(weightOunces),
	}, nil
}

func (node *Node) handleUPSLabelPromise(fields map[string]string) (map[string]string, error) {
	if node.Agent.Kind != "ups_label_printer" {
		return nil, nil
	}
	if fields["promise_about"] != production.PromisePrintLabel {
		return nil, fmt.Errorf("ups label printer cannot handle promise_about=%q", fields["promise_about"])
	}
	packageID := firstStringField(fields, "package_id", "package_id")
	address := firstStringField(fields, "shipping_address", "shipping_address", "address")
	weightOunces := intField(fields, "weight_ounces", "weight_ounces")
	trackingNumber, costCents, err := production.LabelForShipment(packageID, address, weightOunces)
	if err != nil {
		return nil, err
	}
	capabilityAck, err := node.requestPrinterPortCapability()
	if err != nil {
		return nil, err
	}
	labelBytes, err := production.LabelBytesForShipment(map[string]string{
		"package_id":      packageID,
		"tracking_number": trackingNumber,
		"cost_cents":      strconv.Itoa(costCents),
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
		"promise_about":    production.PromisePrintLabel,
		"package_id":       packageID,
		"tracking_number":  trackingNumber,
		"cost_cents":       strconv.Itoa(costCents),
		"printer_spool_id": printAck.Fields["printer_spool_id"],
	}, nil
}

// requestPrinterPortCapability asks the local printer-port kernel role for a
// bounded future-print promise token before any label bytes are presented.
// Intent: The UPS label app receives promise-token event records from the local
// printer resource owner instead of assuming hardware access or treating the
// message kernel as an authorization service. Source: DI-pohaj; DI-vutok
func (node *Node) requestPrinterPortCapability() (parsedMessage, error) {
	tokenID := "printcap-" + node.Agent.Name
	capabilityFields := map[string]string{
		"act":                        decision.ActPromise,
		"from":                       node.Agent.Name,
		"to":                         "printer_port",
		"turn":                       "startup",
		"promise":                    "I promise to receive printer_port's scoped future-print capability token and use it only for bounded UPS label bytes.",
		"reason":                     "ups_label_printer needs local printer-port promise event before asking for hardware printing",
		"promise_about":              production.PromiseIssuePrintCapability,
		"print_capability_issuee":    node.Agent.Name,
		"print_capability_token_id":  tokenID,
		"print_capability_scope":     production.PrintCapabilityScope,
		"print_capability_max_bytes": strconv.Itoa(production.PrintCapabilityMaxBytes),
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
// local print event records if its own token is still recognizable. Source: DI-pohaj;
// DI-vutok
func (node *Node) redeemPrinterPortCapability(capabilityAck parsedMessage, labelBytes []byte) (parsedMessage, error) {
	redemptionFields := map[string]string{
		"act":                        decision.ActPromise,
		"from":                       node.Agent.Name,
		"to":                         "printer_port",
		"turn":                       "startup",
		"promise":                    "I promise to present only bounded UPS label bytes under this printer_port capability token and to receive printer_port's local print event record.",
		"reason":                     "ups_label_printer has a scoped future-print token and now asks printer_port to write exact label bytes",
		"promise_about":              production.PromiseRedeemPrintCapability,
		"print_capability_issuee":    node.Agent.Name,
		"print_capability_token":     capabilityAck.Fields["print_capability_token"],
		"print_capability_token_id":  capabilityAck.Fields["print_capability_token_id"],
		"print_capability_scope":     capabilityAck.Fields["print_capability_scope"],
		"print_capability_max_bytes": capabilityAck.Fields["print_capability_max_bytes"],
		"label_bytes_hex":            hex.EncodeToString(labelBytes),
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
// only receives explicit print event records after token redemption. Source:
// DI-pohaj; DI-vutok
func (node *Node) handlePrinterPortPromise(fields map[string]string) (map[string]string, error) {
	if node.Agent.Kind != "printer_port" {
		return nil, nil
	}
	switch fields["promise_about"] {
	case production.PromiseIssuePrintCapability:
		token, err := production.IssuePrintCapabilityToken(fields)
		if err != nil {
			return nil, err
		}
		tokenID := firstStringField(fields, "print_capability_token_id")
		scope := firstStringField(fields, "print_capability_scope")
		if scope == "" {
			scope = production.PrintCapabilityScope
		}
		maxBytes := firstStringField(fields, "print_capability_max_bytes")
		if maxBytes == "" {
			maxBytes = strconv.Itoa(production.PrintCapabilityMaxBytes)
		}
		node.record("printer_capability_issued", "kept", fields["from"], fmt.Sprintf("token_id=%s scope=%s max_bytes=%s", tokenID, scope, maxBytes))
		return map[string]string{
			"promise_about":              production.PromiseIssuePrintCapability,
			"print_capability_issuee":    firstStringField(fields, "print_capability_issuee", "from"),
			"print_capability_token":     token,
			"print_capability_token_id":  tokenID,
			"print_capability_scope":     scope,
			"print_capability_max_bytes": maxBytes,
		}, nil
	case production.PromiseRedeemPrintCapability:
		spoolID, err := production.PrintLabelToLocalDevice(fields)
		if err != nil {
			return nil, err
		}
		printEvent := firstStringField(fields, "label_bytes_hex")
		node.record("printer_port_printed", "kept", fields["from"], fmt.Sprintf("spool_id=%s label_hex_bytes=%d", spoolID, len(printEvent)))
		return map[string]string{
			"promise_about":    production.PromiseRedeemPrintCapability,
			"printer_spool_id": spoolID,
		}, nil
	default:
		return nil, fmt.Errorf("printer_port cannot handle promise_about=%q", fields["promise_about"])
	}
}

func (node *Node) handleAccountingPromise(fields map[string]string) (map[string]string, error) {
	if node.Agent.Kind != "accounting" {
		return nil, nil
	}
	switch fields["promise_about"] {
	case production.PromiseAddressLookup:
		orderID := firstStringField(fields, "order_id", "order_id")
		address, err := production.AddressForOrder(orderID)
		if err != nil {
			return nil, err
		}
		node.record("shipping_address_promised", "kept", fields["from"], fmt.Sprintf("order_id=%s shipping_address=%s", orderID, address))
		return map[string]string{
			"promise_about":    production.PromiseAddressLookup,
			"order_id":         orderID,
			"shipping_address": address,
		}, nil
	case production.PromiseShipmentUpdate:
		orderID := firstStringField(fields, "order_id", "order_id")
		trackingNumber := firstStringField(fields, "tracking_number", "tracking_number")
		costCents := intField(fields, "cost_cents", "cost_cents")
		if err := production.ValidateAccountingUpdate(orderID, trackingNumber, costCents); err != nil {
			return nil, err
		}
		ackFields := map[string]string{
			"promise_about":   production.PromiseShipmentUpdate,
			"order_id":        orderID,
			"tracking_number": trackingNumber,
			"cost_cents":      strconv.Itoa(costCents),
		}
		updateKey := checkpointKey(pcid.AccountingV1, production.PromiseShipmentUpdate, orderID, trackingNumber, strconv.Itoa(costCents))
		alreadyRecorded := node.rememberCheckpoint(checkpointRecord{
			Key:          updateKey,
			ProtocolName: pcid.AccountingV1,
			PromiseAbout: production.PromiseShipmentUpdate,
			Subject:      orderID,
			Detail:       fmt.Sprintf("tracking_number=%s cost_cents=%d", trackingNumber, costCents),
		})
		if alreadyRecorded {
			ackFields[duplicateShipmentEventField] = "true"
			node.record("accounting_update_duplicate", "kept", fields["from"], fmt.Sprintf("order_id=%s tracking_number=%s cost_cents=%d", orderID, trackingNumber, costCents))
			return ackFields, nil
		}
		node.record("accounting_updated", "kept", fields["from"], fmt.Sprintf("order_id=%s tracking_number=%s cost_cents=%d", orderID, trackingNumber, costCents))
		return ackFields, nil
	default:
		return nil, fmt.Errorf("accounting cannot handle promise_about=%q", fields["promise_about"])
	}
}

// handleCASStoragePromise implements the CAS storage pCID as voluntary promises
// about exact content bytes, retention, replica tokens, and local corruption
// observations.
// Intent: CAS is concrete storage behavior and event records; it is not an RPC
// storage command or central authorization surface. Sparse-store and bearer-token
// behavior remain pCID-owned promises between peers. Source: DI-sinur; DI-manul
func (node *Node) handleCASStoragePromise(fields map[string]string) (map[string]string, error) {
	switch fields["promise_about"] {
	case production.PromiseStoreContent:
		contentCID := fields["content_cid"]
		contentBytes, err := base64.StdEncoding.DecodeString(fields["content_b64"])
		if err != nil {
			return nil, err
		}
		if !production.VerifyContentCID(contentBytes, contentCID) {
			node.record("cas_corrupt_bytes_rejected", "malformed", fields["from"], "pcid="+pcid.CASStorageV1+" presented bytes did not match content_cid="+contentCID)
			return nil, fmt.Errorf("content bytes do not match content CID")
		}
		if intField(fields, "credit_offer") < 3 {
			node.record("economics_price_refused", "non_commitment", fields["from"], "pcid="+pcid.CASStorageV1+" storage credit_offer below local price")
			return map[string]string{
				"promise_about":  production.PromiseStoreContent,
				"storage_status": "price_refused",
				"content_cid":    contentCID,
			}, nil
		}
		if _, storeErr := node.storeLocalCASObject(contentBytes, agentCASStoreOptions{
			Kind:         agentCASKindPeer,
			SourcePeer:   fields["from"],
			ProtocolName: pcid.CASStorageV1,
			Retention:    "paid-run-local",
			Paid:         true,
		}); storeErr != nil {
			return nil, storeErr
		}
		token := node.issueCapabilityToken(fields["from"], contentCID)
		bearerToken := ""
		if fields["token_style"] == "bearer" {
			bearerToken = node.issueBearerStorageToken(contentCID)
			node.record("agent_cas_bearer_storage_token_issued", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" token_scope=bearer content_cid="+contentCID)
		}
		node.record("cas_storage_promised", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("agent_cas_peer_storage_promised", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID+" retention=paid-run-local")
		node.record("cas_retention_promised", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" retention=run-local")
		node.record("cas_bytes_stored", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("economics_credit_accepted", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" credit_offer="+fields["credit_offer"])
		node.record("economics_capacity_reserved", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" units="+firstStringField(fields, "units"))
		node.record("economics_credits_earned", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" Bob earns storage credits")
		node.record("capability_token_issued", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" token_scope=serve-once content_cid="+contentCID)
		node.record("capability_token_ttl_promised", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" ttl=run-local content_cid="+contentCID)
		replicaToken := ""
		if node.Agent.Name == "bob" {
			replicaAck, replicaErr := node.sendAndReceive("frank", map[string]string{
				"act":           decision.ActPromise,
				"from":          node.Agent.Name,
				"to":            "frank",
				"turn":          "startup",
				"promise":       "Bob promises Frank exact bytes for replica storage and receives only Frank's local replica event.",
				"reason":        "replication is a peer promise, not global availability",
				"promise_about": production.PromiseReplicateContent,
				"content_cid":   contentCID,
				"content_b64":   fields["content_b64"],
				"issuee":        fields["from"],
				"units":         "1",
			})
			if replicaErr == nil {
				replicaToken = replicaAck.Fields["replica_token"]
				node.record("cas_replication_promised", "kept", "frank", "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
				node.record("replica_capability_token_received", "kept", "frank", "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
			}
		}
		return map[string]string{
			"promise_about":    production.PromiseStoreContent,
			"storage_status":   "stored",
			"content_cid":      contentCID,
			"capability_token": token,
			"bearer_token":     bearerToken,
			"replica_peer":     "frank",
			"replica_token":    replicaToken,
		}, nil
	case production.PromiseServeContent:
		contentCID := fields["content_cid"]
		if fields["missing_object_probe"] == "true" {
			return node.handleSparseCASProbe(fields, production.PromiseServeContent), nil
		}
		if fields["token_style"] == "bearer" {
			if !node.redeemBearerStorageToken(contentCID, fields["token"]) {
				node.record("agent_cas_bearer_storage_token_rejected", "non_commitment", fields["from"], "pcid="+pcid.CASStorageV1+" bearer token not promised for content_cid="+contentCID)
				return map[string]string{
					"promise_about": production.PromiseServeContent,
					"content_cid":   contentCID,
					"token_status":  "not_promised",
				}, nil
			}
			node.record("agent_cas_bearer_storage_token_redeemed", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" bearer content_cid="+contentCID)
		} else if !node.redeemCapabilityToken(fields["from"], contentCID, fields["token"]) {
			node.record("capability_token_replay_rejected", "non_commitment", fields["from"], "pcid="+pcid.CASStorageV1+" consumed or unrecognized serve-once token for content_cid="+contentCID)
			return map[string]string{
				"promise_about": production.PromiseServeContent,
				"content_cid":   contentCID,
				"token_status":  "not_promised",
			}, nil
		}
		contentBytes, stored, readErr := node.readLocalCASObject(contentCID)
		if readErr != nil {
			return nil, readErr
		}
		if !stored || len(contentBytes) == 0 {
			return nil, fmt.Errorf("content %s not stored", contentCID)
		}
		node.record("capability_token_received", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("capability_token_ttl_observed", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" ttl still valid for content_cid="+contentCID)
		node.record("cas_bytes_retrieved", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("agent_cas_peer_retrieval_promised", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("capability_token_revoked", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" serve-once token consumed for content_cid="+contentCID)
		node.record("gc_promise_ended", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" serve-once token promise ended after redemption for content_cid="+contentCID)
		node.record("gc_object_removed", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" consumed capability token removed from run-scoped store for content_cid="+contentCID)
		return map[string]string{
			"promise_about": production.PromiseServeContent,
			"content_cid":   contentCID,
			"content_b64":   base64.StdEncoding.EncodeToString(contentBytes),
		}, nil
	case production.PromiseReplicateContent:
		contentCID := fields["content_cid"]
		contentBytes, err := base64.StdEncoding.DecodeString(fields["content_b64"])
		if err != nil {
			return nil, err
		}
		if !production.VerifyContentCID(contentBytes, contentCID) {
			return nil, fmt.Errorf("replica bytes do not match content CID")
		}
		if _, storeErr := node.storeLocalCASObject(contentBytes, agentCASStoreOptions{
			Kind:         agentCASKindPeer,
			SourcePeer:   fields["from"],
			ProtocolName: pcid.CASStorageV1,
			Retention:    "replica-run-local",
			Paid:         true,
		}); storeErr != nil {
			return nil, storeErr
		}
		issuee := firstStringField(fields, "issuee", "from")
		token := node.issueCapabilityToken(issuee, contentCID)
		node.record("cas_replica_stored", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("replica_capability_token_issued", "kept", issuee, "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		return map[string]string{
			"promise_about": production.PromiseReplicateContent,
			"content_cid":   contentCID,
			"replica_token": token,
		}, nil
	case production.PromiseServeReplicaContent:
		contentCID := fields["content_cid"]
		if fields["missing_object_probe"] == "true" {
			return node.handleSparseCASProbe(fields, production.PromiseServeReplicaContent), nil
		}
		if !node.redeemCapabilityToken(fields["from"], contentCID, fields["token"]) {
			node.record("capability_token_replay_rejected", "non_commitment", fields["from"], "pcid="+pcid.CASStorageV1+" consumed or unrecognized replica token for content_cid="+contentCID)
			return map[string]string{
				"promise_about": production.PromiseServeReplicaContent,
				"content_cid":   contentCID,
				"token_status":  "not_promised",
			}, nil
		}
		contentBytes, stored, readErr := node.readLocalCASObject(contentCID)
		if readErr != nil {
			return nil, readErr
		}
		if !stored || len(contentBytes) == 0 {
			return nil, fmt.Errorf("replica content %s not stored", contentCID)
		}
		node.record("cas_replica_serve_promised", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("agent_cas_peer_retrieval_promised", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" replica content_cid="+contentCID)
		node.record("replica_capability_token_redeemed", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("capability_token_expired", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" consumed replica token is expired after redemption")
		node.record("capability_token_renewal_requested", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" Alice asks for fresh event if future retrieval is needed")
		node.record("capability_token_renewed", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" Frank promises a fresh run-local token")
		node.record("gc_promise_ended", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" replica serve-once token promise ended after redemption for content_cid="+contentCID)
		node.record("gc_object_removed", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" consumed replica capability token removed from run-scoped store for content_cid="+contentCID)
		node.record("replica_recovery_succeeded", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("cas_bytes_retrieved", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" replica content_cid="+contentCID)
		return map[string]string{
			"promise_about": production.PromiseServeReplicaContent,
			"content_cid":   contentCID,
			"content_b64":   base64.StdEncoding.EncodeToString(contentBytes),
		}, nil
	case production.PromiseReplicaTokenLifecycle:
		contentCID := fields["content_cid"]
		if fields["token_status"] != "transferred" || fields["token_style"] != "bearer" {
			node.record("agent_cas_bearer_storage_token_rejected", "non_commitment", fields["from"], "pcid="+pcid.CASStorageV1+" token lifecycle was not a bearer transfer promise")
			return map[string]string{
				"promise_about": production.PromiseReplicaTokenLifecycle,
				"content_cid":   contentCID,
				"token_status":  "not_promised",
			}, nil
		}
		issuerPeer := firstStringField(fields, "issuer_peer", "redeem_peer")
		bearerToken := firstStringField(fields, "bearer_token", "token")
		node.record("agent_cas_bearer_storage_token_received", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" issuer="+issuerPeer+" content_cid="+contentCID)
		redeemAck, redeemErr := node.sendAndReceive(issuerPeer, map[string]string{
			"act":           decision.ActPromise,
			"from":          node.Agent.Name,
			"to":            issuerPeer,
			"turn":          "startup",
			"promise":       node.Agent.Name + " promises to present a bearer storage token and receive only the issuer's local serve event.",
			"reason":        "bearer token redemption should be voluntary peer storage work, not authority",
			"promise_about": production.PromiseServeContent,
			"content_cid":   contentCID,
			"token":         bearerToken,
			"token_style":   "bearer",
		})
		if redeemErr != nil {
			node.record("agent_cas_bearer_storage_token_rejected", "non_commitment", issuerPeer, "pcid="+pcid.CASStorageV1+" "+redeemErr.Error())
			return map[string]string{
				"promise_about": production.PromiseReplicaTokenLifecycle,
				"content_cid":   contentCID,
				"token_status":  "not_promised",
			}, nil
		}
		contentBytes, decodeErr := base64.StdEncoding.DecodeString(redeemAck.Fields["content_b64"])
		if decodeErr != nil {
			return nil, decodeErr
		}
		if !production.VerifyContentCID(contentBytes, contentCID) {
			return nil, fmt.Errorf("bearer token redeemed bytes do not match content CID")
		}
		if _, storeErr := node.storeLocalCASObject(contentBytes, agentCASStoreOptions{
			Kind:         agentCASKindPeer,
			SourcePeer:   issuerPeer,
			ProtocolName: pcid.CASStorageV1,
			Retention:    "bearer-token-paid-run-local",
			Paid:         true,
		}); storeErr != nil {
			return nil, storeErr
		}
		node.record("agent_cas_bearer_storage_token_redeemed", "kept", issuerPeer, "pcid="+pcid.CASStorageV1+" holder="+node.Agent.Name+" content_cid="+contentCID)
		return map[string]string{
			"promise_about": production.PromiseReplicaTokenLifecycle,
			"content_cid":   contentCID,
			"token_status":  "redeemed",
		}, nil
	case production.PromisePresentStorageReport:
		contentBytes, err := base64.StdEncoding.DecodeString(fields["content_b64"])
		if err != nil {
			return nil, err
		}
		contentCID := fields["content_cid"]
		node.record("cas_verification_promised", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" receiver promises local byte/CID verification")
		if production.VerifyContentCID(contentBytes, contentCID) {
			return map[string]string{"promise_about": production.PromisePresentStorageReport, "verdict": "kept"}, nil
		}
		node.record("cas_corrupt_bytes_rejected", "malformed", fields["from"], "pcid="+pcid.CASStorageV1+" bytes did not match content_cid="+contentCID)
		node.record("cas_corrupt_report_recorded", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" local corrupt-byte event recorded")
		// Intent: Corrupt CAS event is an observable malformed promise by an
		// identifiable peer, so it must enter the same local trust path as bad
		// proofs and bad compute results instead of remaining a log-only event.
		// Source: DI-fijov
		node.observeOutcome(fields["from"], relationship.OutcomeMalformed)
		return map[string]string{"promise_about": production.PromisePresentStorageReport, "verdict": "broken"}, nil
	case production.PromiseLabelFutureMalformedReport:
		node.record("trust_repair_promise_recorded", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" future repair promise recorded as local events")
		node.record("trust_repair_future_only", "non_commitment", fields["from"], "pcid="+pcid.CASStorageV1+" future repair promise is remembered but does not prove repair has already been kept")
		return map[string]string{"promise_about": production.PromiseLabelFutureMalformedReport, "repair_status": "future_only"}, nil
	case production.PromiseUnsupportedVariantProbe:
		node.record("promise_variant_not_promised", "non_commitment", fields["from"], "pcid="+pcid.CASStorageV1+" unsupported storage variant not promised")
		return map[string]string{"promise_about": production.PromiseUnsupportedVariantProbe, "variant_status": "not_promised"}, nil
	default:
		return nil, fmt.Errorf("CAS storage cannot handle promise_about=%q", fields["promise_about"])
	}
}

// handleCIDComputePromise implements bounded CID-named compute, cache lookup,
// and verifier event records.
// Intent: Compute results are promises over exact function/input/context bytes
// that receivers can recompute or ask peers to observe locally. Source: DI-sinur
func (node *Node) handleCIDComputePromise(message parsedMessage) (protocolHandlerResult, error) {
	fields := message.Fields
	switch fields["promise_about"] {
	case production.PromiseLookupComputeCache:
		cacheKey := production.ComputeCacheKey(pcid.CIDComputeV1, fields["function_cid"], fields["input_cid"], fields["context_cid"], fields["result_cid"])
		node.mu.Lock()
		cachedFields, ok := node.computeCache[cacheKey]
		node.mu.Unlock()
		if !ok {
			node.record("compute_cache_miss", "non_commitment", fields["from"], "pcid="+pcid.CIDComputeV1+" cache_key="+cacheKey)
			return protocolHandlerResult{Fields: map[string]string{
				"promise_about": production.PromiseLookupComputeCache,
				"cache_key":     cacheKey,
				"cache_status":  "miss",
			}}, nil
		}
		node.record("compute_cache_hit", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" cache_key="+cacheKey)
		ackFields := copyStringMap(cachedFields)
		ackFields["promise_about"] = production.PromiseLookupComputeCache
		ackFields["cache_key"] = cacheKey
		ackFields["cache_status"] = "hit"
		return protocolHandlerResult{Fields: ackFields}, nil
	case production.PromiseExecuteFunction:
		if fields["capacity_probe"] == "true" {
			node.record("economics_capacity_refused", "non_commitment", fields["from"], "pcid="+pcid.CIDComputeV1+" local compute capacity withheld for probe")
			node.record("backpressure_capacity_refused", "non_commitment", fields["from"], "pcid="+pcid.CIDComputeV1+" receiver keeps spare capacity by not promising this probe")
			return protocolHandlerResult{Fields: map[string]string{"promise_about": production.PromiseExecuteFunction, "compute_status": "capacity_refused"}}, nil
		}
		if node.Agent.Kind == "stdio_agent" {
			ackBytes, err := node.runStdioComputeWorker(message)
			if err != nil {
				return protocolHandlerResult{}, err
			}
			ackMessage, parseErr := node.parseEnvelope(ackBytes)
			if parseErr != nil {
				return protocolHandlerResult{}, parseErr
			}
			return protocolHandlerResult{Fields: ackMessage.Fields, AckBytes: ackBytes}, nil
		}
		resultBytes, err := node.executeComputeFunction(fields)
		if err != nil {
			return protocolHandlerResult{}, err
		}
		resultCID := production.ContentCID(resultBytes)
		inputs, inputErr := production.DecodeComputeInputs(fields)
		if inputErr != nil {
			return protocolHandlerResult{}, inputErr
		}
		node.record("compute_context_promised", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" context_cid="+fields["context_cid"])
		node.record("compute_function_executed", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" function_cid="+fields["function_cid"]+" result_cid="+resultCID)
		if production.FunctionKind(inputs.FunctionBytes) != "fibonacci" {
			node.record("compute_alternate_function_executed", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" function_kind="+production.FunctionKind(inputs.FunctionBytes)+" result_cid="+resultCID)
		}
		node.record("cid_compute_promised", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" compute promised only for stated CIDs")
		node.record("compute_result_promised", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" result_cid="+resultCID)
		node.record("economics_credit_accepted", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" credit_offer="+fields["credit_offer"])
		node.record("economics_capacity_reserved", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" units="+firstStringField(fields, "units"))
		node.record("economics_credits_earned", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" "+node.Agent.Name+" earns compute credits")
		if node.Agent.Kind == "wasm_agent" {
			return protocolHandlerResult{Fields: production.ExecuteComputePromiseFields(fields, resultBytes)}, nil
		}
		badResultBytes := production.BadComputeResultBytes(resultBytes)
		node.record("compute_bad_result_promised", "malformed", fields["from"], "pcid="+pcid.CIDComputeV1+" hash-valid but semantically wrong result_cid="+production.ContentCID(badResultBytes))
		return protocolHandlerResult{Fields: map[string]string{
			"promise_about":  production.PromiseExecuteFunction,
			"function_cid":   fields["function_cid"],
			"function_b64":   fields["function_b64"],
			"input_cid":      fields["input_cid"],
			"input_b64":      fields["input_b64"],
			"context_cid":    fields["context_cid"],
			"context_b64":    fields["context_b64"],
			"result_cid":     resultCID,
			"result_b64":     base64.StdEncoding.EncodeToString(resultBytes),
			"bad_result_cid": production.ContentCID(badResultBytes),
			"bad_result_b64": base64.StdEncoding.EncodeToString(badResultBytes),
		}}, nil
	case production.PromiseVerifyComputeResult:
		resultPromiser := firstStringField(fields, "result_promiser", "from")
		verifyErr := verifyComputeResultFields(fields)
		if fields["disagreement_probe"] == "true" {
			node.record("compute_verifier_disagreement", "non_commitment", resultPromiser, "pcid="+pcid.CIDComputeV1+" verifier withholds endorsement for disagreement pressure")
			return protocolHandlerResult{Fields: map[string]string{
				"promise_about":      production.PromiseVerifyComputeResult,
				"verdict":            "disagree",
				"subject_peer":       resultPromiser,
				"subject_result_cid": fields["result_cid"],
			}}, nil
		}
		if verifyErr != nil {
			node.record("compute_result_peer_rejected", "malformed", resultPromiser, "pcid="+pcid.CIDComputeV1+" "+verifyErr.Error())
			return protocolHandlerResult{Fields: map[string]string{
				"promise_about":      production.PromiseVerifyComputeResult,
				"verdict":            "broken",
				"subject_peer":       resultPromiser,
				"subject_result_cid": fields["result_cid"],
			}}, nil
		}
		node.record("compute_result_peer_verified", "kept", resultPromiser, "pcid="+pcid.CIDComputeV1+" result_cid="+fields["result_cid"])
		cacheKey := production.ComputeCacheKey(pcid.CIDComputeV1, fields["function_cid"], fields["input_cid"], fields["context_cid"], fields["result_cid"])
		node.mu.Lock()
		node.computeCache[cacheKey] = copyStringMap(fields)
		node.mu.Unlock()
		node.record("compute_cache_checkpointed", "kept", resultPromiser, "pcid="+pcid.CIDComputeV1+" cache_key="+cacheKey)
		return protocolHandlerResult{Fields: map[string]string{
			"promise_about":      production.PromiseVerifyComputeResult,
			"verdict":            "kept",
			"subject_peer":       resultPromiser,
			"subject_result_cid": fields["result_cid"],
		}}, nil
	default:
		return protocolHandlerResult{}, fmt.Errorf("CID compute cannot handle promise_about=%q", fields["promise_about"])
	}
}

// executeComputeFunction keeps one execute_function promise through the local
// runtime appropriate to this agent.
// Intent: Ordinary Go compute peers keep the existing production helper path,
// while Peggy proves useful WASM execution without changing cid_compute_v1
// semantics or introducing a runtime-specific pCID. Source: DI-sivis
func (node *Node) executeComputeFunction(fields map[string]string) ([]byte, error) {
	inputs, err := production.DecodeComputeInputs(fields)
	if err != nil {
		return nil, err
	}
	if err := production.VerifyComputeInputCIDs(fields, inputs); err != nil {
		return nil, err
	}
	if node.Agent.Kind != "wasm_agent" {
		return production.ExecuteFunction(inputs.FunctionBytes, inputs.InputBytes, inputs.ContextBytes)
	}
	if production.FunctionKind(inputs.FunctionBytes) != "fibonacci" {
		return nil, fmt.Errorf("WASM compute peer cannot handle function kind %q", production.FunctionKind(inputs.FunctionBytes))
	}
	n, inputErr := production.ParseFibonacciInput(inputs.InputBytes)
	if inputErr != nil {
		return nil, inputErr
	}
	result, runErr := runtimeadapter.RunWASMModule(context.Background(), runtimeadapter.MinimalWASMModule, uint64(n))
	if runErr != nil {
		return nil, runErr
	}
	resultBytes := production.FibonacciResultBytes(n, result.ExportValue, inputs.ContextBytes)
	expectedBytes, expectedErr := production.ExecuteFunction(inputs.FunctionBytes, inputs.InputBytes, inputs.ContextBytes)
	if expectedErr != nil {
		return nil, expectedErr
	}
	if production.ContentCID(resultBytes) != production.ContentCID(expectedBytes) {
		return nil, fmt.Errorf("WASM compute result CID mismatch")
	}
	moduleHash := protocol.HashExactBytes(runtimeadapter.MinimalWASMModule)
	node.record("wasm_compute_request_received", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" module_sha256="+moduleHash+" input_n="+fmt.Sprintf("%d", n))
	node.record("wasm_compute_function_executed", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" export="+result.ExportName+" value="+fmt.Sprintf("%d", result.ExportValue))
	node.record("wasm_compute_result_promised", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" result_cid="+production.ContentCID(resultBytes))
	return resultBytes, nil
}

// handleIdentityKeyPromise records the narrow identity/key promises currently
// modeled by POC15.
// Intent: Key rotation is identity protocol behavior, not a generic report pCID
// or a new top-level action kind. Source: DI-vipih
func (node *Node) handleIdentityKeyPromise(fields map[string]string) (map[string]string, error) {
	switch fields["promise_about"] {
	case production.PromiseRotateSigningKey:
		node.record("key_rotation_promise_recorded", "kept", fields["from"], "pcid="+pcid.IdentityKeyV1+" new_key_label="+fields["new_key_label"])
		return map[string]string{
			"promise_about":  production.PromiseRotateSigningKey,
			"new_key_label":  fields["new_key_label"],
			"rotation_scope": fields["rotation_scope"],
			"verdict":        "kept",
		}, nil
	default:
		return nil, fmt.Errorf("identity key cannot handle promise_about=%q", fields["promise_about"])
	}
}

func (node *Node) recordAckEvent(target string, message parsedMessage) {
	switch message.Fields["promise_about"] {
	case production.PromiseWeighPackage:
		node.record("package_weight_received", "kept", target, "weight_ounces="+message.Fields["weight_ounces"])
	case production.PromiseAddressLookup:
		node.record("shipping_address_received", "kept", target, "shipping_address="+message.Fields["shipping_address"])
	case production.PromisePrintLabel:
		node.record("shipping_label_received", "kept", target, "tracking_number="+message.Fields["tracking_number"]+" cost_cents="+message.Fields["cost_cents"])
	case production.PromiseShipmentUpdate:
		if message.Fields[duplicateShipmentEventField] == "true" {
			node.record("accounting_update_duplicate_confirmed", "kept", target, "tracking_number="+message.Fields["tracking_number"]+" cost_cents="+message.Fields["cost_cents"])
			return
		}
		node.record("accounting_update_confirmed", "kept", target, "tracking_number="+message.Fields["tracking_number"]+" cost_cents="+message.Fields["cost_cents"])
	case production.PromiseIssuePrintCapability:
		node.record("printer_capability_received", "kept", target, "token_id="+message.Fields["print_capability_token_id"]+" scope="+message.Fields["print_capability_scope"])
	case production.PromiseRedeemPrintCapability:
		node.record("printer_port_print_confirmed", "kept", target, "spool_id="+message.Fields["printer_spool_id"])
	case production.PromiseStoreContent:
		switch message.Fields["storage_status"] {
		case "price_refused":
			node.record("economics_price_refused", "non_commitment", target, "pcid="+pcid.CASStorageV1+" content_cid="+message.Fields["content_cid"])
		case "stored":
			node.record("capability_token_received", "kept", target, "pcid="+pcid.CASStorageV1+" content_cid="+message.Fields["content_cid"])
			if message.Fields["replica_token"] != "" {
				node.record("replica_capability_token_received", "kept", target, "pcid="+pcid.CASStorageV1+" replica_peer="+message.Fields["replica_peer"])
			}
		}
	case production.PromiseServeContent:
		if message.Fields["token_status"] == "not_promised" {
			node.record("capability_token_replay_observed", "non_commitment", target, "pcid="+pcid.CASStorageV1+" content_cid="+message.Fields["content_cid"])
			return
		}
		node.record("cas_bytes_retrieved", "kept", target, "pcid="+pcid.CASStorageV1+" content_cid="+message.Fields["content_cid"])
	case production.PromiseServeReplicaContent:
		if message.Fields["token_status"] == "not_promised" {
			node.record("capability_token_replay_observed", "non_commitment", target, "pcid="+pcid.CASStorageV1+" replica content_cid="+message.Fields["content_cid"])
			return
		}
		node.record("replica_recovery_succeeded", "kept", target, "pcid="+pcid.CASStorageV1+" content_cid="+message.Fields["content_cid"])
	case production.PromiseLookupComputeCache:
		if message.Fields["cache_status"] == "miss" {
			node.record("compute_cache_miss_observed", "non_commitment", target, "pcid="+pcid.CIDComputeV1+" cache_key="+message.Fields["cache_key"])
			return
		}
		node.record("compute_cache_reused", "kept", target, "pcid="+pcid.CIDComputeV1+" cache_key="+message.Fields["cache_key"])
		node.record("compute_result_received", "kept", target, "pcid="+pcid.CIDComputeV1+" cache result_cid="+message.Fields["result_cid"])
	case production.PromiseUnsupportedVariantProbe:
		if message.Fields["variant_status"] == "not_promised" {
			node.record("promise_variant_not_promised", "non_commitment", target, "pcid="+pcid.CASStorageV1+" unsupported variant was not promised")
			node.rememberNonCommitment(target, pcid.CASStorageV1, message.Fields, "ack variant_status not_promised")
		}
	case production.PromiseExecuteFunction:
		if message.Fields["compute_status"] == "capacity_refused" {
			node.record("economics_capacity_refused", "non_commitment", target, "pcid="+pcid.CIDComputeV1+" compute capacity probe refused")
			return
		}
		node.record("compute_result_received", "kept", target, "pcid="+pcid.CIDComputeV1+" result_cid="+message.Fields["result_cid"])
	case production.PromiseVerifyComputeResult:
		// Intent: Verification reports are cid_compute_v1 promises about a
		// compute result, not generic report messages. Source:
		// DI-vipih
		node.record("compute_verification_report_received", "kept", target, "pcid="+pcid.CIDComputeV1+" verifier verdict="+message.Fields["verdict"]+" result_cid="+message.Fields["subject_result_cid"])
		if message.Fields["verdict"] == "disagree" {
			node.record("compute_verifier_disagreement", "non_commitment", target, "pcid="+pcid.CIDComputeV1+" verifier disagreement received")
			node.record("compute_disagreement_resolved_locally", "kept", message.Fields["subject_peer"], "pcid="+pcid.CIDComputeV1+" Alice resolves disagreement by local recompute plus Dave event")
			return
		}
		if message.Fields["verdict"] == "kept" {
			node.record("compute_verification_received", "kept", target, "pcid="+pcid.CIDComputeV1+" peer verified result")
			return
		}
		node.record("compute_verification_received", "malformed", target, "pcid="+pcid.CIDComputeV1+" peer rejected result")
	}
}

func (node *Node) verifyComputeAckLocally(message parsedMessage, target string) error {
	if err := verifyComputeResultFields(message.Fields); err != nil {
		return err
	}
	node.record("compute_result_locally_verified", "kept", target, "pcid="+pcid.CIDComputeV1+" result_cid="+message.Fields["result_cid"])
	node.record("economics_credits_spent", "kept", target, "pcid="+pcid.CIDComputeV1+" Alice spends compute credits after local verification")
	if message.Fields["bad_result_b64"] == "" {
		return nil
	}
	badFields := copyStringMap(message.Fields)
	badFields["result_b64"] = message.Fields["bad_result_b64"]
	badFields["result_cid"] = message.Fields["bad_result_cid"]
	if err := verifyComputeResultFields(badFields); err == nil {
		node.record("compute_result_locally_rejected", "malformed", target, "pcid="+pcid.CIDComputeV1+" bad-result probe unexpectedly verified")
		// Intent: A peer that promises a semantically bad compute result creates
		// local negative trust event record even when the malformed bytes are only a
		// probe inside an otherwise parseable response. Source: DI-sunuf
		node.observeOutcome(target, relationship.OutcomeMalformed)
		return nil
	}
	node.record("compute_result_locally_rejected", "malformed", target, "pcid="+pcid.CIDComputeV1+" local recompute rejected bad-result probe")
	// Intent: Local recomputation is Alice's own event record that the compute
	// promiser exposed a malformed result candidate; this affects Alice's trust
	// in that promiser rather than any global authority. Source: DI-sunuf
	node.observeOutcome(target, relationship.OutcomeMalformed)
	node.record("economics_payment_withheld", "non_commitment", target, "pcid="+pcid.CIDComputeV1+" Alice withholds payment for the bad-result probe")
	return nil
}

func decodeComputeInputs(fields map[string]string) ([]byte, []byte, []byte, error) {
	inputs, err := production.DecodeComputeInputs(fields)
	if err != nil {
		return nil, nil, nil, err
	}
	return inputs.FunctionBytes, inputs.InputBytes, inputs.ContextBytes, nil
}

func verifyComputeInputCIDs(fields map[string]string, functionBytes, inputBytes, contextBytes []byte) error {
	return production.VerifyComputeInputCIDs(fields, production.ComputeInputs{
		FunctionBytes: functionBytes,
		InputBytes:    inputBytes,
		ContextBytes:  contextBytes,
	})
}

func verifyComputeResultFields(fields map[string]string) error {
	functionBytes, inputBytes, contextBytes, err := decodeComputeInputs(fields)
	if err != nil {
		return err
	}
	if err := verifyComputeInputCIDs(fields, functionBytes, inputBytes, contextBytes); err != nil {
		return err
	}
	resultBytes, resultErr := base64.StdEncoding.DecodeString(fields["result_b64"])
	if resultErr != nil {
		return resultErr
	}
	if !production.VerifyContentCID(resultBytes, fields["result_cid"]) {
		return fmt.Errorf("result bytes do not match result CID")
	}
	expectedBytes, executeErr := production.ExecuteFunction(functionBytes, inputBytes, contextBytes)
	if executeErr != nil {
		return executeErr
	}
	expectedCID := production.ContentCID(expectedBytes)
	if expectedCID != fields["result_cid"] {
		return fmt.Errorf("local recompute result_cid=%s want %s", expectedCID, fields["result_cid"])
	}
	return nil
}

func copyStringMap(fields map[string]string) map[string]string {
	copiedFields := make(map[string]string, len(fields))
	for key, value := range fields {
		copiedFields[key] = value
	}
	return copiedFields
}

func copyMapOrEmpty(fields map[string]string) map[string]string {
	copiedFields := make(map[string]string, len(fields))
	for key, value := range fields {
		copiedFields[key] = value
	}
	return copiedFields
}

func copyIntMapOrEmpty(fields map[string]int) map[string]int {
	copiedFields := make(map[string]int, len(fields))
	for key, value := range fields {
		copiedFields[key] = value
	}
	return copiedFields
}

func copyNestedMapOrEmpty(fields map[string]map[string]string) map[string]map[string]string {
	copiedFields := make(map[string]map[string]string, len(fields))
	for key, value := range fields {
		copiedFields[key] = copyMapOrEmpty(value)
	}
	return copiedFields
}

func copyNonCommitmentMapOrEmpty(fields map[string]nonCommitmentRecord) map[string]nonCommitmentRecord {
	copiedFields := make(map[string]nonCommitmentRecord, len(fields))
	for key, value := range fields {
		copiedFields[key] = value
	}
	return copiedFields
}

func copyCheckpointMapOrEmpty(fields map[string]checkpointRecord) map[string]checkpointRecord {
	copiedFields := make(map[string]checkpointRecord, len(fields))
	for key, value := range fields {
		copiedFields[key] = value
	}
	return copiedFields
}

func copyPromiseMapOrEmpty(fields map[string]promiseRecord) map[string]promiseRecord {
	copiedFields := make(map[string]promiseRecord, len(fields))
	for key, value := range fields {
		copiedFields[key] = value
	}
	return copiedFields
}

func (node *Node) issueCapabilityToken(issuee, contentCID string) string {
	token := production.ContentCID([]byte(node.Agent.Name + "|" + issuee + "|" + contentCID + "|serve-once"))
	node.mu.Lock()
	node.capabilityTokens[issuee+"|"+contentCID] = token
	node.mu.Unlock()
	return token
}

// issueBearerStorageToken records a transferable storage token promise by the
// local storage issuer.
// Intent: POC15 needs a minimal bearer-token incentive path where the current
// holder can redeem storage without treating the token as external authority.
// Source: DI-manul
func (node *Node) issueBearerStorageToken(contentCID string) string {
	token := production.ContentCID([]byte(node.Agent.Name + "|" + contentCID + "|bearer-storage"))
	node.mu.Lock()
	node.capabilityTokens["bearer|"+contentCID+"|"+token] = token
	node.mu.Unlock()
	return token
}

func (node *Node) redeemCapabilityToken(issuee, contentCID, token string) bool {
	node.mu.Lock()
	defer node.mu.Unlock()
	tokenKey := issuee + "|" + contentCID
	if node.capabilityTokens[tokenKey] != token {
		return false
	}
	delete(node.capabilityTokens, tokenKey)
	return true
}

// redeemBearerStorageToken consumes one transferable storage token from the
// issuer's local token table.
// Intent: Token redemption is an issuer-kept promise over local bytes; it does
// not impose storage, routing, or trust decisions on any other agent. Source:
// DI-manul
func (node *Node) redeemBearerStorageToken(contentCID, token string) bool {
	node.mu.Lock()
	defer node.mu.Unlock()
	tokenKey := "bearer|" + contentCID + "|" + token
	if node.capabilityTokens[tokenKey] != token {
		return false
	}
	delete(node.capabilityTokens, tokenKey)
	return true
}

// handleSparseCASProbe reports whether this agent currently has local bytes for
// a probed CID.
// Intent: Sparse CAS misses are ordinary non-commitment records, because an
// incomplete peer store has not broken a promise merely by lacking an object it
// did not promise to retain. Source: DI-manul
func (node *Node) handleSparseCASProbe(fields map[string]string, promiseAbout string) map[string]string {
	contentCID := fields["content_cid"]
	if node.localCASObjectExists(contentCID) {
		node.record("agent_cas_sparse_object_present", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		return map[string]string{
			"promise_about": promiseAbout,
			"content_cid":   contentCID,
			"token_status":  "present",
		}
	}
	node.record("agent_cas_sparse_object_missing", "non_commitment", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
	node.record("agent_cas_retrieval_not_promised", "non_commitment", fields["from"], "pcid="+pcid.CASStorageV1+" sparse store has no local bytes for content_cid="+contentCID)
	return map[string]string{
		"promise_about": promiseAbout,
		"content_cid":   contentCID,
		"token_status":  "not_promised",
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
	if fields["promise_about"] == "reciprocal_economics" {
		reciprocalValue = 4
	}
	offer := economy.Offer{
		Promiser:        node.Agent.Name,
		Promisee:        target,
		Resource:        fields["promise_about"],
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
// candidate peers are reachable only for explicit low-risk discovery or
// future-repair promises, not arbitrary traffic. Source: DI-timah; DI-fijov
func (node *Node) canDialTarget(peerName string, fields map[string]string) bool {
	node.mu.Lock()
	defer node.mu.Unlock()
	if node.ledger.PermanentlyDistrusted(peerName) {
		return false
	}
	if node.ledger.CanDial(peerName) {
		return true
	}
	return isLowRiskCandidatePromise(fields) && containsName(node.Agent.CandidatePeers, peerName)
}

// canAcceptFrom reports whether this node currently promises to accept one TCP
// exchange from a peer. Intent: Candidate-peer discovery is a narrow voluntary
// acceptance promise, and future repair after malformed events is similarly
// narrow. Neither creates broad permission or a global routing rule.
// Source: DI-timah; DI-fijov
func (node *Node) canAcceptFrom(peerName string, fields map[string]string) bool {
	node.mu.Lock()
	defer node.mu.Unlock()
	if node.ledger.PermanentlyDistrusted(peerName) {
		return false
	}
	if node.ledger.CanAccept(peerName) {
		return true
	}
	return isLowRiskCandidatePromise(fields) && containsName(node.Agent.CandidatePeers, peerName)
}

// markPermanentDistrustAndTransitExclusion updates Alice-local trust and route
// state for the hard trust-line scenario.
// Intent: The decision changes only this app's own future sends and path choices;
// it does not command other agents or publish global reputation. Source:
// DI-dubih
func (node *Node) markPermanentDistrustAndTransitExclusion(peerName string) {
	node.mu.Lock()
	defer node.mu.Unlock()
	node.ledger.PermanentlyDistrust(peerName)
	node.ledger.ExcludeTransit(peerName)
}

func (node *Node) routeAllowed(route []string) bool {
	node.mu.Lock()
	defer node.mu.Unlock()
	return node.ledger.RouteAllowed(route)
}

func (node *Node) observeOutcome(peerName string, outcome relationship.Outcome) {
	if peerName == "" || peerName == node.Agent.Name {
		return
	}
	node.mu.Lock()
	priorTrust := node.ledger.Trust(peerName)
	priorCaution := node.ledger.Caution(peerName)
	transition := node.ledger.ObserveOutcome(peerName, outcome)
	trustScore := node.ledger.Trust(peerName)
	cautionScore := node.ledger.Caution(peerName)
	node.mu.Unlock()
	node.record(string(transition), "kept", peerName, fmt.Sprintf("outcome=%s trust=%d caution=%d", outcome, trustScore, cautionScore))
	if cautionScore > priorCaution {
		node.record("trust_caution_recorded", string(outcome), peerName, fmt.Sprintf("outcome=%s trust=%d caution=%d", outcome, trustScore, cautionScore))
	}
	if priorCaution > cautionScore {
		node.record("trust_caution_consumed", string(outcome), peerName, fmt.Sprintf("outcome=%s trust=%d caution=%d prior_caution=%d", outcome, trustScore, cautionScore, priorCaution))
	}
	if priorCaution > cautionScore && trustScore == priorTrust {
		node.record("trust_recovery_delayed", "kept", peerName, fmt.Sprintf("outcome=%s trust=%d caution=%d prior_caution=%d", outcome, trustScore, cautionScore, priorCaution))
	}
	// Intent: POC15 keeps POC12 relationship-transition events and also emits
	// the inherited POC15 analyzer's `trust_updated` event name so future
	// POCs cannot silently drop trust behavior while refactoring analyzers.
	// Source: DI-sinur
	node.record("trust_updated", string(outcome), peerName, fmt.Sprintf("outcome=%s trust=%d caution=%d transition=%s", outcome, trustScore, cautionScore, transition))
}

// rememberOutstandingPromise adds one exact promise record to this app's local
// journal before later ACK or receive handling can resolve it.
// Intent: Trust changes should be explainable from a local promise-event
// record rather than from transport success, kernel routing, or local resource
// pressure alone. Source: DI-vujob
func (node *Node) rememberOutstandingPromise(peerName, protocolName, exactHash string, fields map[string]string) string {
	fingerprint := promiseRecordKey(peerName, protocolName, "", fields)
	record := promiseRecord{
		Key:           promiseRecordKey(peerName, protocolName, exactHash, fields),
		Fingerprint:   fingerprint,
		Peer:          peerName,
		ProtocolName:  protocolName,
		ExactHash:     exactHash,
		PromiseAbout:  fields["promise_about"],
		PromiseText:   fields["promise"],
		ExpectedEvent: fields["reason"],
		Status:        promiseStatusOutstanding,
	}
	node.mu.Lock()
	node.promiseJournal[record.Key] = record
	node.mu.Unlock()
	node.record("promise_outstanding", "kept", peerName, fmt.Sprintf("status=%s protocol=%s exact_sha256=%s promise_about=%s", record.Status, record.ProtocolName, record.ExactHash, record.PromiseAbout))
	return record.Key
}

// resolveOutstandingPromise records the local outcome of a previously journaled
// promise without itself deciding whether peer trust should change.
// Intent: Promise resolution events and trust mutation are deliberately
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
// Intent: Alice exhausting Alice's budget or capacity is an event about Alice's
// local state, not an event that Bob kept or broke a promise. Source: DI-vujob
func (node *Node) recordLocalResourceExhaustion(target string, fields map[string]string, detail string) {
	resourceName := resourceField(fields)
	if resourceName == "" {
		resourceName = fields["promise_about"]
	}
	node.record("local_resource_exhausted", "non_commitment", target, "resource="+resourceName+" detail="+detail)
}

// rememberNonCommitment stores one receiver `not_promised` outcome as local
// restraint event record for later turns in the same POC15 run.
// Intent: A peer that did not promise the requested exchange should not be
// asked again for the same target/protocol/promise-about tuple until a later
// design adds an explicit new-event release rule. Source: DI-zapab
func (node *Node) rememberNonCommitment(target, protocolName string, fields map[string]string, detail string) {
	promiseAbout := fields["promise_about"]
	record := nonCommitmentRecord{
		Key:          nonCommitmentKey(target, protocolName, promiseAbout),
		Peer:         target,
		ProtocolName: protocolName,
		PromiseAbout: promiseAbout,
		Detail:       detail,
	}
	node.mu.Lock()
	node.nonCommitmentJournal[record.Key] = record
	node.mu.Unlock()
}

// shouldSuppressNonCommittedPromise checks whether this app already has local
// non-commitment event record for the same semantic promise target.
// Intent: Suppression turns a prior `not_promised` into Alice's restraint, not
// Bob's punishment; it records no peer-trust transition. Source: DI-zapab
func (node *Node) shouldSuppressNonCommittedPromise(target string, fields map[string]string) bool {
	protocolName, _ := node.protocolForFields(fields)
	promiseAbout := fields["promise_about"]
	key := nonCommitmentKey(target, protocolName, promiseAbout)
	node.mu.Lock()
	record, exists := node.nonCommitmentJournal[key]
	node.mu.Unlock()
	if !exists {
		return false
	}
	node.record("promise_not_promised_suppressed", "non_commitment", target, "protocol="+record.ProtocolName+" promise_about="+record.PromiseAbout+" detail="+record.Detail)
	return true
}

// suppressRepeatedPromise avoids sending the same live-agent promise text to the
// same target/protocol once this app already has a journal event for it.
// Intent: Repetition after a prior promise outcome creates pressure that looks
// RPC-like; POC15 should instead record local non-commitment and wait for a new
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
	node.record("promise_repeated_suppressed", "non_commitment", target, "protocol="+protocolName+" promise_about="+fields["promise_about"])
	return true
}

func repairErrDetail(validateErr error) string {
	return "repaired common live decision formatting issue: " + validateErr.Error()
}

// recordDecisionError records LLM/provider/runtime failures as local event records,
// not as broken peer promises.
// Intent: A transient provider failure or runtime decision failure does not mean
// any peer broke a promise, so it should not enter peer trust as a broken
// event record. Source: DI-jinoz
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
// Intent: A receiver's `not_promised` ACK is an event record of non-commitment, not a
// broken peer promise; local transport failures are not peer events at all.
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

// eventUpdatesTrust reports whether an ACK payload event should change peer
// trust or merely be recorded as already-seen local events.
// Intent: Duplicate shipment-update confirmations should remain visible in logs
// without repeatedly increasing trust for the same order/tracking/cost
// checkpoint. Negative ACK payload verdicts and future-only repair promises
// should not inflate local trust merely because the receiver returned a
// parseable ACK. Source: DI-jinoz; DI-punib; DI-fijov
func eventUpdatesTrust(fields map[string]string) bool {
	if fields == nil {
		return true
	}
	if fields[duplicateShipmentEventField] == "true" {
		return false
	}
	switch fields["verdict"] {
	case "broken", "malformed", "disagree", "not_promised":
		return false
	}
	if fields["variant_status"] == "not_promised" ||
		fields["storage_status"] == "price_refused" ||
		fields["compute_status"] == "capacity_refused" ||
		fields["cache_status"] == "miss" ||
		fields["token_status"] == "not_promised" ||
		fields["repair_status"] == "future_only" ||
		fields["replay_status"] == "not_promised" {
		return false
	}
	return true
}

// promiseRecordKey uses exact bytes for concrete promise instances and promise
// text/about fields for repeat-suppression fingerprints before bytes exist.
// Intent: POC15 needs exact-byte promise accounting for real sends while still
// detecting repeated live-agent promise intent before sending again. Source:
// DI-vujob
func promiseRecordKey(peerName, protocolName, exactHash string, fields map[string]string) string {
	if exactHash != "" {
		return peerName + "|" + protocolName + "|" + exactHash
	}
	return peerName + "|" + protocolName + "|" + fields["promise_about"] + "|" + fields["promise"]
}

// promiseStatusOutcome maps journal-only statuses into the small outcome
// vocabulary used by existing POC15 logs and analyzer summaries.
// Intent: Local failures and non-commitments should remain non-commitment in
// reports, while duplicate events stay kept-but-non-mutating. Source:
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
// aligned at the interface where a real ACK or inbound promise is resolved.
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

// promiseStatusForNonTrustingEvent classifies ACK payload events that are
// visible but deliberately does not update peer trust.
// Intent: Refusals, cache misses, replay refusals, future-only repair, and
// unsupported variants are peer non-commitments rather than duplicate kept
// events; only semantic checkpoints such as duplicate shipment updates should
// resolve as duplicate events. Source: DI-sihuz
func promiseStatusForNonTrustingEvent(fields map[string]string) promiseStatus {
	if fields == nil {
		return promiseStatusDuplicate
	}
	if fields[duplicateShipmentEventField] == "true" {
		return promiseStatusDuplicate
	}
	switch fields["verdict"] {
	case "broken":
		return promiseStatusBroken
	case "malformed":
		return promiseStatusMalformed
	case "disagree", "not_promised":
		return promiseStatusNonCommitment
	}
	if fields["variant_status"] == "not_promised" ||
		fields["storage_status"] == "price_refused" ||
		fields["compute_status"] == "capacity_refused" ||
		fields["cache_status"] == "miss" ||
		fields["token_status"] == "not_promised" ||
		fields["repair_status"] == "future_only" ||
		fields["replay_status"] == "not_promised" {
		return promiseStatusNonCommitment
	}
	return promiseStatusDuplicate
}

// nonCommitmentKey keeps receiver non-commitment scoped to one peer, one
// protocol, and one protocol-owned promise meaning.
// Intent: A `not_promised` ACK should restrain only the semantic exchange that
// was actually declined, not all future cooperation with that peer. Source:
// DI-zapab
func nonCommitmentKey(peerName, protocolName, promiseAbout string) string {
	return strings.Join([]string{peerName, protocolName, promiseAbout}, "|")
}

// rememberCheckpoint records a reusable local checkpoint and reports whether it
// was already present.
// Intent: Replayed events should stay auditable while avoiding repeated trust
// mutation for the same app-local checkpoint. Source: DI-zapab
func (node *Node) rememberCheckpoint(record checkpointRecord) bool {
	node.mu.Lock()
	defer node.mu.Unlock()
	alreadySeen := node.checkpointAlreadySeen(record.Key)
	if !alreadySeen {
		node.checkpointJournal[record.Key] = record
	}
	return alreadySeen
}

// checkpointAlreadySeen assumes the caller holds node.mu and performs only the
// local map lookup.
// Intent: Keeping the lookup small makes checkpoint semantics reusable without
// turning them into a protocol-level idempotency authority. Source: DI-zapab
func (node *Node) checkpointAlreadySeen(key string) bool {
	_, alreadySeen := node.checkpointJournal[key]
	return alreadySeen
}

// checkpointKey builds stable local checkpoint identifiers from pCID-selected
// protocol meaning plus deterministic subject fields.
// Intent: Checkpoint identity is local event records over promise content, not a
// universal transaction ID or global command key. Source: DI-zapab
func checkpointKey(protocolName, promiseAbout string, parts ...string) string {
	keyParts := append([]string{protocolName, promiseAbout}, parts...)
	return strings.Join(keyParts, "|")
}

// isLinkDiscoveryPromise recognizes the pCID-owned payload meaning used for
// candidate-peer link formation.
// Intent: Link discovery is represented as promise content under the same
// top-level act, not as a separate protocol verb. Source: DI-timah
func isLinkDiscoveryPromise(fields map[string]string) bool {
	for _, key := range []string{"promise_about", "meaning", "intent", "link_intent"} {
		if fields[key] == decision.PromiseAboutLinkDiscovery {
			return true
		}
	}
	return false
}

// isLowRiskCandidatePromise recognizes promise meanings a peer may voluntarily
// receive even before ordinary direct-trust adjacency exists.
// Intent: A future-repair promise after malformed events is low risk enough to
// hear, but it still does not prove repair or permit arbitrary follow-up
// messages. Source: DI-fijov
func isLowRiskCandidatePromise(fields map[string]string) bool {
	return isLinkDiscoveryPromise(fields) || fields["promise_about"] == production.PromiseLabelFutureMalformedReport
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
	requestedUnits := intField(fields, "units", "requested_units", "capacity", "capacity_mb")
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
// Intent: Treat malformed or extreme resource promises as local events about
// the sender, not as commands the receiver must obey. Source: DI-timah
func (node *Node) checkIncomingResourcePromise(fields map[string]string) error {
	resourceType := resourceField(fields)
	if resourceType == "" {
		return nil
	}
	requestedUnits := intField(fields, "units", "requested_units", "capacity", "capacity_mb")
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
	stakeAmount := intField(fields, "stake", "collateral", "stake", "collateral")
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
	for _, key := range []string{"resource", "resource_type", "resource", "promise_about"} {
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

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
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
	if node.eventOutcomeCounts == nil {
		node.eventOutcomeCounts = make(map[string]int)
	}
	node.eventOutcomeCounts[outcome]++
	node.mu.Unlock()
	encoded, err := json.Marshal(event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal event: %v\n", err)
		return
	}
	node.stdoutMu.Lock()
	fmt.Println(string(encoded))
	node.stdoutMu.Unlock()
	if node.logFile != nil {
		if _, writeErr := node.logFile.Write(append(encoded, '\n')); writeErr != nil {
			fmt.Fprintf(os.Stderr, "write event: %v\n", writeErr)
		}
	}
}

// emitMessageArtifact sends exact PromiseGrid envelope bytes to the
// observer-only collector for post-run review.
// Intent: Operators need raw messages, not only event records about messages.
// Apps still do not share the observer Docker volume and cannot read these
// artifacts back; this one-way stdout record is harness plumbing that does not
// affect trust, routing, ACKs, or protocol semantics. Source: DI-tuhop
func (node *Node) emitMessageArtifact(direction, peer, protocolName string, envelopeBytes []byte, fields map[string]string) {
	if len(envelopeBytes) == 0 {
		return
	}
	node.recordAgentCASMessageArtifact(direction, peer, protocolName, envelopeBytes, fields)
	artifactProtocol := protocolName
	if strings.TrimSpace(artifactProtocol) == "" {
		artifactProtocol = "unknown"
	}
	artifact := eventstream.MessageArtifact{
		Observer:            node.Agent.Name,
		Direction:           direction,
		Peer:                peer,
		Protocol:            artifactProtocol,
		ExactSHA256:         protocol.HashExactBytes(envelopeBytes),
		EnvelopeBytesBase64: base64.StdEncoding.EncodeToString(envelopeBytes),
		SourceEvent:         "runtime." + direction,
	}
	if fields != nil {
		artifact.ParentExactSHA256 = firstNonEmpty(fields["envelope_parent_exact_sha256"], fields["payload_parent_exact_sha256"], fields["parent_exact_sha256"])
		artifact.ParentLinkLocation = firstNonEmpty(fields["parent_link_location"], parentLinkLocationFromFields(fields))
		artifact.PromiseAbout = fields["promise_about"]
	}
	if envelope, parseErr := protocol.ParseEnvelope(envelopeBytes); parseErr == nil && len(envelope.ParentExactSHA256s) > 0 {
		// Intent: ACK artifact records should derive envelope parent links from
		// the signed bytes themselves, not only from caller-supplied compatibility
		// fields. This keeps raw-message DAG review aligned with the demux data
		// that persistent sessions actually use. Source: DI-vopab
		artifact.ParentExactSHA256 = envelope.ParentExactSHA256s[0]
		artifact.ParentLinkLocation = "envelope"
	}
	record := eventstream.Record{
		Kind:            eventstream.KindMessageArtifact,
		Source:          "agent:" + node.Agent.Name,
		MessageArtifact: &artifact,
	}
	recordBytes, marshalErr := json.Marshal(record)
	if marshalErr != nil {
		node.record("raw_message_artifact_emit_failed", "broken", peer, marshalErr.Error())
		return
	}
	node.stdoutMu.Lock()
	fmt.Println(string(recordBytes))
	node.stdoutMu.Unlock()
	node.record("raw_message_artifact_emitted", "kept", peer, "direction="+direction+" pcid="+artifactProtocol+" exact_sha256="+artifact.ExactSHA256)
}

func parentLinkLocationFromFields(fields map[string]string) string {
	if fields["envelope_parent_exact_sha256"] != "" {
		return "envelope"
	}
	if fields["payload_parent_exact_sha256"] != "" {
		return "payload"
	}
	if fields["parent_exact_sha256"] != "" {
		return "payload"
	}
	return ""
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
	appSession := node.appSession
	node.appSession = nil
	node.mu.Unlock()
	if appSession != nil {
		if closeErr := appSession.Close(); closeErr != nil {
			node.record("app_receive_conn_close_failed", "broken", "kernel", closeErr.Error())
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

// waitForShutdownGrace leaves receive promises open briefly after active turns end.
// Intent: Peers that made late but still legitimate promises get a bounded
// chance to finish their exchange without consulting observer-owned run files.
// Source: DI-galin; DI-dirat
func (node *Node) waitForShutdownGrace(ctx context.Context) {
	graceDuration := node.Config.ShutdownGrace()
	if graceDuration <= 0 {
		return
	}
	timer := time.NewTimer(graceDuration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		node.record("shutdown_grace_timeout", "non_commitment", "", ctx.Err().Error())
	case <-timer.C:
		node.record("shutdown_grace_elapsed", "kept", "", "local receive grace elapsed before receive promises close")
	}
}

func (node *Node) runDir() string {
	return filepath.Join(node.Config.RunRoot, node.Config.RunID)
}

func (node *Node) runScopedStatePath() string {
	return filepath.Join(node.runDir(), "stores", node.Agent.Name, "durable-state.json")
}

// loadRunScopedState restores app-local state that belongs to this one POC15
// run. Intent: A process restart inside a run should not erase CAS files,
// sparse-CAS metadata, compute checkpoints, replay windows, or promise journals,
// but the clean-run reset can still delete the entire run root before the next
// experiment. Source: DI-sunuf; DI-manul; DI-fagog
func (node *Node) loadRunScopedState() error {
	statePath := node.runScopedStatePath()
	stateBytes, readErr := os.ReadFile(statePath)
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) {
			node.record("run_scoped_store_empty", "kept", "", "no run-scoped durable state exists yet for "+node.Agent.Name)
			return nil
		}
		return readErr
	}
	var state runScopedState
	if unmarshalErr := json.Unmarshal(stateBytes, &state); unmarshalErr != nil {
		return unmarshalErr
	}
	node.mu.Lock()
	node.capabilityTokens = copyMapOrEmpty(state.CapabilityTokens)
	node.computeCache = copyNestedMapOrEmpty(state.ComputeCache)
	node.nonCommitmentJournal = copyNonCommitmentMapOrEmpty(state.NonCommitmentJournal)
	node.checkpointJournal = copyCheckpointMapOrEmpty(state.CheckpointJournal)
	node.promiseJournal = copyPromiseMapOrEmpty(state.PromiseJournal)
	node.eventOutcomeCounts = copyIntMapOrEmpty(state.EventOutcomeCounts)
	node.replayJournal = copyMapOrEmpty(state.ReplayJournal)
	node.agentCASStore = copyAgentCASMapOrEmpty(state.AgentCASObjects)
	node.agentMessageDAG = copyAgentMessageDAGMapOrEmpty(state.AgentMessageDAG)
	node.mu.Unlock()
	migratedCount, migrateErr := node.migrateLegacyRunScopedCASObjects(state.CASObjects, state.AgentCASObjects)
	if migrateErr != nil {
		return migrateErr
	}
	node.record("run_scoped_store_loaded", "kept", "", fmt.Sprintf("agent_cas=%d agent_dag=%d tokens=%d compute=%d checkpoints=%d promises=%d replay=%d legacy_cas_migrated=%d", len(state.AgentCASObjects), len(state.AgentMessageDAG), len(state.CapabilityTokens), len(state.ComputeCache), len(state.CheckpointJournal), len(state.PromiseJournal), len(state.ReplayJournal), migratedCount))
	node.record("cas_run_store_loaded", "kept", "", fmt.Sprintf("agent_cas_objects=%d legacy_cas_migrated=%d", len(state.AgentCASObjects), migratedCount))
	node.record("event_run_store_loaded", "kept", "", fmt.Sprintf("promise_journal=%d checkpoint_journal=%d non_commitments=%d receiver_non_commitments=%d replay=%d", len(state.PromiseJournal), len(state.CheckpointJournal), state.EventOutcomeCounts[string(relationship.OutcomeNonCommitment)], len(state.NonCommitmentJournal), len(state.ReplayJournal)))
	node.record("compute_cache_run_store_loaded", "kept", "", fmt.Sprintf("compute_cache=%d", len(state.ComputeCache)))
	return nil
}

// migrateLegacyRunScopedCASObjects imports pre-filesystem CAS byte blobs into
// this agent's filesystem CAS and leaves new saves metadata-only.
// Intent: DI-fagog changes the durable-state split without losing restartability
// for any already-written POC15 run state that still has `cas_objects_b64`.
// Source: DI-fagog
func (node *Node) migrateLegacyRunScopedCASObjects(legacyCASObjects map[string]string, legacyMetadata map[string]agentCASObject) (int, error) {
	migratedCount := 0
	for contentCID, encodedBytes := range legacyCASObjects {
		contentBytes, decodeErr := base64.StdEncoding.DecodeString(encodedBytes)
		if decodeErr != nil {
			return migratedCount, fmt.Errorf("decode run-scoped CAS object %s: %w", contentCID, decodeErr)
		}
		if !production.VerifyContentCID(contentBytes, contentCID) {
			return migratedCount, fmt.Errorf("legacy CAS object %s bytes do not match content CID", contentCID)
		}
		legacyRecord := legacyMetadata[contentCID]
		storeOptions := agentCASStoreOptions{
			Kind:         firstNonEmpty(legacyRecord.Kind, agentCASKindPeer),
			SourcePeer:   legacyRecord.SourcePeer,
			ProtocolName: firstNonEmpty(legacyRecord.ProtocolName, "legacy_cas_objects_b64"),
			Retention:    firstNonEmpty(legacyRecord.Retention, "migrated-run-local"),
			Encrypted:    legacyRecord.Encrypted,
			Pinned:       legacyRecord.Pinned,
			Paid:         legacyRecord.Paid,
			ParentCIDs:   legacyRecord.ParentCIDs,
		}
		if _, storeErr := node.storeLocalCASObject(contentBytes, storeOptions); storeErr != nil {
			return migratedCount, storeErr
		}
		if len(legacyRecord.MissingParentCIDs) > 0 {
			node.mu.Lock()
			if migratedRecord, ok := node.agentCASStore[contentCID]; ok {
				migratedRecord.MissingParentCIDs = append([]string(nil), legacyRecord.MissingParentCIDs...)
				node.agentCASStore[contentCID] = migratedRecord
			}
			node.mu.Unlock()
		}
		migratedCount++
	}
	if migratedCount > 0 {
		node.record("agent_cas_legacy_json_bytes_migrated", "kept", "", fmt.Sprintf("objects=%d", migratedCount))
	}
	return migratedCount, nil
}

// saveRunScopedState writes restartable run state through a temporary file and
// rename. Intent: Partial writes should not corrupt the app's local events or
// sparse-CAS metadata if a process is killed mid-save, filesystem CAS bytes stay
// outside JSON, and the saved state remains scoped to this run root rather than
// becoming cross-run truth. Source: DI-sunuf; DI-manul; DI-fagog
func (node *Node) saveRunScopedState() error {
	state := node.exportRunScopedState()
	statePath := node.runScopedStatePath()
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		return err
	}
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
	node.record("run_scoped_store_saved", "kept", "", fmt.Sprintf("agent_cas=%d agent_dag=%d tokens=%d compute=%d checkpoints=%d promises=%d replay=%d", len(state.AgentCASObjects), len(state.AgentMessageDAG), len(state.CapabilityTokens), len(state.ComputeCache), len(state.CheckpointJournal), len(state.PromiseJournal), len(state.ReplayJournal)))
	node.record("cas_run_store_saved", "kept", "", fmt.Sprintf("agent_cas_objects=%d", len(state.AgentCASObjects)))
	node.record("event_run_store_saved", "kept", "", fmt.Sprintf("promise_journal=%d checkpoint_journal=%d non_commitments=%d receiver_non_commitments=%d replay=%d", len(state.PromiseJournal), len(state.CheckpointJournal), state.EventOutcomeCounts[string(relationship.OutcomeNonCommitment)], len(state.NonCommitmentJournal), len(state.ReplayJournal)))
	node.record("compute_cache_run_store_saved", "kept", "", fmt.Sprintf("compute_cache=%d", len(state.ComputeCache)))
	return nil
}

func (node *Node) exportRunScopedState() runScopedState {
	node.mu.Lock()
	defer node.mu.Unlock()
	return runScopedState{
		Version:              1,
		CapabilityTokens:     copyMapOrEmpty(node.capabilityTokens),
		ComputeCache:         copyNestedMapOrEmpty(node.computeCache),
		NonCommitmentJournal: copyNonCommitmentMapOrEmpty(node.nonCommitmentJournal),
		CheckpointJournal:    copyCheckpointMapOrEmpty(node.checkpointJournal),
		PromiseJournal:       copyPromiseMapOrEmpty(node.promiseJournal),
		EventOutcomeCounts:   copyIntMapOrEmpty(node.eventOutcomeCounts),
		ReplayJournal:        copyMapOrEmpty(node.replayJournal),
		AgentCASObjects:      copyAgentCASMapOrEmpty(node.agentCASStore),
		AgentMessageDAG:      copyAgentMessageDAGMapOrEmpty(node.agentMessageDAG),
	}
}

// recordRunScopedRetentionAndGC records retention and deletion as local promise
// event records. Intent: GC and backpressure are not central policy; each agent
// promises how it will keep, remove, or decline local objects under run-end,
// token-expiry, disk-pressure, and superseded-checkpoint conditions. Source:
// DI-sunuf
func (node *Node) recordRunScopedRetentionAndGC() {
	node.recordAgentCASGCEvents()
	state := node.exportRunScopedState()
	node.record("retention_promise_recorded", "kept", "", "run-scoped CAS, compute, token, replay, and journal state is retained only for this POC15 run")
	node.record("retention_until_promised", "kept", "", "retain_until=run_end_or_clean_reset")
	node.record("delete_after_promised", "kept", "", "delete_after=operator clean-run reset or explicit local GC promise condition")
	node.record("token_expiry_gc_promised", "kept", "", "serve-once capability tokens expire on redemption or run reset")
	node.record("disk_pressure_gc_promised", "kept", "", "local agent may remove retained objects before accepting new storage if disk pressure exceeds local promise budget")
	node.record("superseded_checkpoint_gc_promised", "kept", "", fmt.Sprintf("superseded checkpoints may be compacted locally checkpoint_count=%d", len(state.CheckpointJournal)))
	node.record("gc_object_retained", "kept", "", fmt.Sprintf("agent_cas=%d tokens=%d compute=%d checkpoints=%d promises=%d replay=%d", len(state.AgentCASObjects), len(state.CapabilityTokens), len(state.ComputeCache), len(state.CheckpointJournal), len(state.PromiseJournal), len(state.ReplayJournal)))
}

func (node *Node) recordRetentionPromiseBroken(subject, detail string) {
	node.record("retention_promise_broken", "broken", subject, detail)
}

// rememberReplayEnvelope keeps a local exact-byte replay window for received
// envelopes. Intent: Replays are not commands to punish a peer globally; they are
// local event records that the same exact promise bytes should not be counted as fresh
// promise event records again. Source: DI-sunuf
func (node *Node) rememberReplayEnvelope(peerName, protocolName, exactHash string) bool {
	node.recordReplayWindowPromise(protocolName)
	node.mu.Lock()
	priorEvent, replayed := node.replayJournal[exactHash]
	if !replayed {
		node.replayJournal[exactHash] = peerName + "|" + protocolName
	}
	node.mu.Unlock()
	if replayed {
		node.record("replay_envelope_rejected", "non_commitment", peerName, "pcid="+protocolName+" exact_sha256="+exactHash+" prior="+priorEvent)
		return true
	}
	node.record("replay_envelope_recorded", "kept", peerName, "pcid="+protocolName+" exact_sha256="+exactHash)
	return false
}

func (node *Node) recordReplayWindowPromise(protocolName string) {
	key := checkpointKey("run-replay-window", protocolName, node.Agent.Name)
	if node.rememberCheckpoint(checkpointRecord{
		Key:          key,
		ProtocolName: protocolName,
		PromiseAbout: "run_replay_window",
		Subject:      node.Agent.Name,
		Detail:       "exact-envelope replay window scoped to one POC15 run",
	}) {
		return
	}
	node.record("replay_window_promised", "kept", "", "pcid="+protocolName+" exact envelope hashes are remembered only inside this run")
}

// recordOutboundPressurePromises records sender-side rate promises before bytes
// leave the app. Intent: A sender voluntarily bounds its own pressure instead of
// treating a peer or kernel as something it can command. Source: DI-sunuf
func (node *Node) recordOutboundPressurePromises(target, protocolName string, fields map[string]string) {
	key := checkpointKey("outbound-pressure", protocolName, target, fields["promise_about"])
	if node.rememberCheckpoint(checkpointRecord{
		Key:          key,
		ProtocolName: protocolName,
		PromiseAbout: fields["promise_about"],
		Subject:      target,
		Detail:       "sender-side pressure promise",
	}) {
		return
	}
	node.record("send_rate_promised", "kept", target, "pcid="+protocolName+" sender promises bounded sends for promise_about="+fields["promise_about"])
	node.record("rate_limit_self_promise_recorded", "kept", target, "pcid="+protocolName+" sender keeps rate limiting as a self-promise, not receiver control")
}

// recordInboundPressurePromises records receiver-side capacity promises after the
// app has decided it currently accepts this peer. Intent: Receive capacity is
// the receiver's local promise interface, not permission granted by any external
// authority. Source: DI-sunuf
func (node *Node) recordInboundPressurePromises(peerName, protocolName string, fields map[string]string) {
	key := checkpointKey("inbound-pressure", protocolName, peerName, fields["promise_about"])
	if node.rememberCheckpoint(checkpointRecord{
		Key:          key,
		ProtocolName: protocolName,
		PromiseAbout: fields["promise_about"],
		Subject:      peerName,
		Detail:       "receiver-side pressure promise",
	}) {
		return
	}
	node.record("backpressure_capacity_promised", "kept", peerName, "pcid="+protocolName+" receiver promises only bounded capacity for promise_about="+fields["promise_about"])
	node.record("accept_rate_promised", "kept", peerName, "pcid="+protocolName+" receiver promises only bounded accept rate for this peer/protocol")
}

// relationshipStatePath is the durable local memory file for this agent's
// private trust ledger.
// Intent: Keep relationship learning across POC15 runs without introducing a
// global trust database or shared authority. Source: DI-timah
func (node *Node) relationshipStatePath() string {
	return filepath.Join(node.Config.RunRoot, "relationships", node.Agent.Name+".json")
}

// loadRelationshipState restores this agent's prior local trust snapshot if it
// exists; absence simply means this is the first run for that agent.
// Intent: Let multi-run POC15 experiments test relationship decay and repair
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
// writes from corrupting the next run's local events. Source: DI-timah
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

// drainInflight waits briefly for active receive handlers before process-exit
// event records are finalized.
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
		node.record("inflight_drained", "kept", "", "all active receive handlers completed before process exit")
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
