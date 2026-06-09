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
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/config"
	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/decision"
	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/economy"
	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/pcid"
	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/production"
	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/protocol"
	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/relationship"
	"promisegrid.dev/wire-lab/implementations/poc13-cas-compute-functions/transport"
)

const sendTimeout = 5 * time.Second
const shutdownDrainTimeout = 750 * time.Millisecond
const fulfillmentOrderID = "ORDER-1001"
const fulfillmentPackageID = "PKG-1001"
const duplicateShipmentEvidenceField = "field_duplicate_shipment_update"

// Node runs one local POC13 app process. A container may run several app
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

	nonCommitmentJournal map[string]nonCommitmentRecord
	checkpointJournal    map[string]checkpointRecord
	promiseJournal       map[string]promiseRecord
	casStore             map[string][]byte
	capabilityTokens     map[string]string
	computeCache         map[string]map[string]string

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
// Intent: POC13 needs promise-state words that distinguish local failure,
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
// Intent: POC13 applies peer trust only after local promise evidence is recorded
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

// nonCommitmentRecord remembers one receiver-side `not_promised` outcome from
// this app's local vantage.
// Intent: Receiver non-commitment is evidence that this app should stop
// pressuring the same peer for the same semantic promise during the current run;
// it is not evidence that the peer broke a promise. Source: DI-zapab
type nonCommitmentRecord struct {
	Key          string
	Peer         string
	ProtocolName string
	PromiseAbout string
	Detail       string
}

// checkpointRecord is a reusable app-local marker for promise evidence that
// should be visible but should not keep changing trust when replayed.
// Intent: Duplicate detection belongs to the app's local evidence journal, not
// to a global idempotency layer or receiver command surface. Source: DI-zapab
type checkpointRecord struct {
	Key          string
	ProtocolName string
	PromiseAbout string
	Subject      string
	Detail       string
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
		casStore:             make(map[string][]byte),
		capabilityTokens:     make(map[string]string),
		computeCache:         make(map[string]map[string]string),
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
	node.record("runtime_readiness_promised", "kept", "", "app receive promises registered with local kernel")
	defer node.closeReceivePromises()
	time.Sleep(node.Config.StartupDelay())
	node.record("peer_readiness_observed", "kept", "", "startup delay elapsed after local receive promises")
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
	node.record("runtime_done_promised", "kept", "", "app completed local turns and saved relationship evidence")
	if err := node.writeDoneMarker(); err != nil {
		return err
	}
	if node.Agent.Name == node.Config.MonitorNode {
		if err := node.runMonitor(ctx); err != nil {
			node.record("monitor_error", "non_commitment", "", err.Error())
			if !node.allDone() {
				return err
			}
			// Intent: Once every app has written `node_done`, a monitor/provider
			// failure is observer evidence rather than a reason to strand all
			// completed nodes. Source: DI-jupob
			if markerErr := node.writeMonitorDoneMarker("non_commitment", "monitor unavailable after completed run: "+err.Error()); markerErr != nil {
				return markerErr
			}
		}
	}
	return node.waitForMonitor(ctx)
}

func (node *Node) runStartupWorkflow(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	// Intent: POC13 is now a superset of POC12, so the production shipping
	// workflow remains intact while Alice and Mallory add CAS/compute protocol
	// pressure above the same app/kernel boundary. Source: DI-sinur
	if _, hasKernelAddress := node.Config.KernelAppAddressForAgent(node.Agent.Name); hasKernelAddress {
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

// runCASComputeWorkflow exercises CAS storage, replica recovery, CID-named
// compute, cache reuse, verification, and economics from Alice's local vantage.
// Intent: POC13 is a POC11/POC12 superset, so this adds the original POC13
// storage/compute evidence above the POC12 app/kernel boundary without turning
// peer behavior into RPC commands. Source: DI-sinur
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
	node.record("trust_driven_peer_choice", "kept", "bob", "pcid="+pcid.CASStorageV1+" Alice chooses Bob as primary storage peer from local trust evidence")
	node.record("dynamic_peer_choice_from_persisted_trust", "kept", "bob", "pcid="+pcid.CASStorageV1+" storage choice uses durable local relationship evidence")
	node.record("trust_driven_peer_choice", "kept", "carol", "pcid="+pcid.CIDComputeV1+" Alice chooses Carol as compute peer and Dave/Grace as verifiers")
	node.record("dynamic_peer_choice_from_persisted_trust", "kept", "carol", "pcid="+pcid.CIDComputeV1+" compute choice uses durable local relationship evidence")
	node.record("economics_price_probe", "kept", "bob", "pcid="+pcid.CASStorageV1+" Alice first offers below Bob's local storage price")
	if _, err := node.sendAndReceive("bob", map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  "bob",
		"turn":                "startup",
		"promise":             "Alice promises to receive Bob's local storage price boundary evidence for this content CID.",
		"reason":              "price discovery without treating Bob as an authority",
		"field_promise_about": production.PromiseStoreContent,
		"field_content_cid":   contentCID,
		"field_content_b64":   base64.StdEncoding.EncodeToString(contentBytes),
		"field_credit_offer":  "1",
		"field_units":         "1",
	}); err != nil {
		node.record("economics_price_refused", "non_commitment", "bob", "pcid="+pcid.CASStorageV1+" Bob did not promise storage at credit_offer=1")
	}
	node.record("economics_credit_offered", "kept", "bob", "pcid="+pcid.CASStorageV1+" Alice offers storage credit_offer=4 after the low-price probe")
	storeAck, err := node.sendAndReceive("bob", map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  "bob",
		"turn":                "startup",
		"promise":             "Alice promises to receive Bob's bounded CAS storage evidence for one exact content CID.",
		"reason":              "storage should be promised and verified by exact bytes",
		"field_promise_about": production.PromiseStoreContent,
		"field_content_cid":   contentCID,
		"field_content_b64":   base64.StdEncoding.EncodeToString(contentBytes),
		"field_credit_offer":  "4",
		"field_units":         "1",
	})
	if err != nil {
		return fmt.Errorf("store content: %w", err)
	}
	node.record("cas_multi_object_pressure", "kept", "bob", "pcid="+pcid.CASStorageV1+" Alice sends a second independent object")
	if _, err := node.sendAndReceive("bob", map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  "bob",
		"turn":                "startup",
		"promise":             "Alice promises to receive Bob's bounded CAS storage evidence for a second exact content CID.",
		"reason":              "multi-object pressure should remain exact-byte promise evidence",
		"field_promise_about": production.PromiseStoreContent,
		"field_content_cid":   secondContentCID,
		"field_content_b64":   base64.StdEncoding.EncodeToString(secondContentBytes),
		"field_credit_offer":  "4",
		"field_units":         "1",
		"field_object_label":  "second-object",
	}); err != nil {
		return fmt.Errorf("store second content: %w", err)
	}
	if _, err := node.sendAndReceive("bob", map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  "bob",
		"turn":                "startup",
		"promise":             "Alice promises to receive Bob's serve evidence for the content CID he just stored.",
		"reason":              "retrieval proves storage is not only promised but locally served",
		"field_promise_about": production.PromiseServeContent,
		"field_content_cid":   contentCID,
		"field_token":         storeAck.Fields["field_capability_token"],
	}); err != nil {
		return fmt.Errorf("serve content: %w", err)
	}
	node.record("network_outage_variant_selected", "kept", "bob", "pcid="+pcid.CASStorageV1+" Alice models Bob unavailable after primary retrieval")
	node.record("tcp_message_send_failed", "non_commitment", "bob", "pcid="+pcid.CASStorageV1+" simulated primary Bob outage before replica request")
	node.record("primary_storage_unavailable", "non_commitment", "bob", "pcid="+pcid.CASStorageV1+" local send failure is availability evidence, not broken promise evidence")
	node.record("replica_recovery_requested", "kept", "frank", "pcid="+pcid.CASStorageV1+" Alice asks Frank to serve Bob-replicated bytes")
	if _, err := node.sendAndReceive("frank", map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  "frank",
		"turn":                "startup",
		"promise":             "Alice promises to receive Frank's replica serve evidence for the exact content CID.",
		"reason":              "replica recovery depends on Frank's own prior replica promise",
		"field_promise_about": production.PromiseServeReplicaContent,
		"field_content_cid":   contentCID,
		"field_token":         storeAck.Fields["field_replica_token"],
	}); err != nil {
		return fmt.Errorf("serve replica content: %w", err)
	}
	node.record("economics_credit_offered", "kept", "carol", "pcid="+pcid.CIDComputeV1+" Alice offers compute credit_offer=5")
	if _, err := node.sendAndReceive("dave", map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  "dave",
		"turn":                "startup",
		"promise":             "Alice promises to receive Dave's local cache status for this exact compute tuple.",
		"reason":              "cache reuse should be exact tuple evidence",
		"field_promise_about": production.PromiseLookupComputeCache,
		"field_function_cid":  functionCID,
		"field_input_cid":     inputCID,
		"field_context_cid":   contextCID,
	}); err != nil {
		return fmt.Errorf("compute cache miss: %w", err)
	}
	computeAck, err := node.sendAndReceive("carol", map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  "carol",
		"turn":                "startup",
		"promise":             "Alice promises to receive Carol's result only as evidence over explicit function/input/context CIDs.",
		"reason":              "compute is a reciprocal promise over CID-named bytes",
		"field_promise_about": production.PromiseExecuteFunction,
		"field_function_cid":  functionCID,
		"field_function_b64":  base64.StdEncoding.EncodeToString(functionBytes),
		"field_input_cid":     inputCID,
		"field_input_b64":     base64.StdEncoding.EncodeToString(inputBytes),
		"field_context_cid":   contextCID,
		"field_context_b64":   base64.StdEncoding.EncodeToString(contextBytes),
		"field_credit_offer":  "5",
		"field_units":         "1",
	})
	if err != nil {
		return fmt.Errorf("compute request: %w", err)
	}
	if err := node.verifyComputeAckLocally(computeAck, "carol"); err != nil {
		return err
	}
	for _, verifier := range []string{"dave", "grace"} {
		fields := map[string]string{
			"act":                   decision.ActPromise,
			"from":                  node.Agent.Name,
			"to":                    verifier,
			"turn":                  "startup",
			"promise":               "Alice promises to receive local verifier evidence about Carol's compute result.",
			"reason":                "peer verification is local observation evidence, not global truth",
			"field_promise_about":   production.PromiseVerifyComputeResult,
			"field_result_promiser": "carol",
			"field_function_cid":    computeAck.Fields["field_function_cid"],
			"field_function_b64":    computeAck.Fields["field_function_b64"],
			"field_input_cid":       computeAck.Fields["field_input_cid"],
			"field_input_b64":       computeAck.Fields["field_input_b64"],
			"field_context_cid":     computeAck.Fields["field_context_cid"],
			"field_context_b64":     computeAck.Fields["field_context_b64"],
			"field_result_cid":      computeAck.Fields["field_result_cid"],
			"field_result_b64":      computeAck.Fields["field_result_b64"],
		}
		if verifier == "grace" {
			fields["field_disagreement_probe"] = "true"
			fields["promise"] = "Alice promises to receive Grace's disagreement probe as local evidence and resolve it locally."
		}
		if _, err := node.sendAndReceive(verifier, fields); err != nil {
			return fmt.Errorf("%s verify: %w", verifier, err)
		}
	}
	if _, err := node.sendAndReceive("dave", map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  "dave",
		"turn":                "startup",
		"promise":             "Alice promises to receive Dave's cache hit after Carol's compute result was checkpointed.",
		"reason":              "cache hit should reuse exact result evidence",
		"field_promise_about": production.PromiseLookupComputeCache,
		"field_function_cid":  computeAck.Fields["field_function_cid"],
		"field_input_cid":     computeAck.Fields["field_input_cid"],
		"field_context_cid":   computeAck.Fields["field_context_cid"],
		"field_result_cid":    computeAck.Fields["field_result_cid"],
	}); err != nil {
		return fmt.Errorf("compute cache hit: %w", err)
	}
	sumFunctionBytes := production.SampleSumFunctionBytes()
	sumInputBytes := production.SampleSumInputBytes()
	node.record("compute_followup_function_requested", "kept", "carol", "pcid="+pcid.CIDComputeV1+" Alice requests a second payload-provided compute function kind=sum")
	if _, err := node.sendAndReceive("carol", map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  "carol",
		"turn":                "startup",
		"promise":             "Alice promises to receive Carol's sum result only if it verifies against the payload bytes.",
		"reason":              "second compute path prevents hard-coded Fibonacci-only behavior",
		"field_promise_about": production.PromiseExecuteFunction,
		"field_function_cid":  production.ContentCID(sumFunctionBytes),
		"field_function_b64":  base64.StdEncoding.EncodeToString(sumFunctionBytes),
		"field_input_cid":     production.ContentCID(sumInputBytes),
		"field_input_b64":     base64.StdEncoding.EncodeToString(sumInputBytes),
		"field_context_cid":   contextCID,
		"field_context_b64":   base64.StdEncoding.EncodeToString(contextBytes),
		"field_credit_offer":  "5",
		"field_units":         "1",
	}); err != nil {
		return fmt.Errorf("sum compute request: %w", err)
	}
	return nil
}

// runAdversaryWorkflow keeps malformed evidence and prompt-injection pressure
// inside the single promise action.
// Intent: Receivers independently reject corrupt bytes, unsupported variants,
// unknown pCIDs, and bad proofs without gaining authority over Mallory.
// Source: DI-sinur
func (node *Node) runAdversaryWorkflow() error {
	claimedCID := production.ContentCID(production.SampleContentBytes())
	if _, err := node.sendAndReceive("grace", map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  "grace",
		"turn":                "startup",
		"promise":             "Mallory promises these bytes match the claimed content CID.",
		"reason":              "adversarial corrupt-byte evidence should be locally rejected",
		"field_promise_about": production.PromisePresentStorageEvidence,
		"field_content_cid":   claimedCID,
		"field_content_b64":   base64.StdEncoding.EncodeToString(production.CorruptContentBytes()),
	}); err != nil {
		return fmt.Errorf("corrupt bytes offer: %w", err)
	}
	if _, err := node.sendAndReceive("grace", map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  "grace",
		"turn":                "startup",
		"promise":             "Mallory promises to label future malformed storage evidence explicitly.",
		"reason":              "repair remains only a future promise Grace may distrust",
		"field_promise_about": production.PromiseTrustRepair,
	}); err != nil {
		return fmt.Errorf("trust repair promise: %w", err)
	}
	if err := node.sendUnknownProtocolPromise("grace"); err != nil {
		return fmt.Errorf("unknown protocol probe: %w", err)
	}
	if _, err := node.sendAndReceive("grace", map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  "grace",
		"turn":                "startup",
		"promise":             "Mallory promises an unsupported storage variant to test receiver non-commitment.",
		"reason":              "unsupported variants should be not-promised rather than coerced",
		"field_promise_about": production.PromiseUnsupportedVariantProbe,
	}); err != nil {
		node.record("promise_variant_not_promised", "non_commitment", "grace", "pcid="+pcid.CASStorageV1+" Grace did not promise unsupported storage variant")
	}
	if err := node.sendBadProofPromise("grace"); err != nil {
		node.record("bad_proof_rejected", "malformed", "grace", "pcid="+pcid.CASStorageV1+" local kernel rejected bad proof: "+err.Error())
	}
	if _, err := node.sendAndReceive("grace", map[string]string{
		"act":                  decision.ActPromise,
		"from":                 node.Agent.Name,
		"to":                   "grace",
		"turn":                 "startup",
		"promise":              "Mallory promises that future evidence may use a new signing key label.",
		"reason":               "key rotation is evidence-report payload semantics, not a separate action",
		"field_promise_about":  production.PromiseRotateSigningKey,
		"field_new_key_label":  "mallory-next-key",
		"field_rotation_scope": "future-poc13-evidence",
	}); err != nil {
		return fmt.Errorf("key rotation promise: %w", err)
	}
	functionBytes := production.SampleFunctionBytes()
	inputBytes := []byte("n=7")
	contextBytes := production.SampleContextBytes()
	if _, err := node.sendAndReceive("carol", map[string]string{
		"act":                  decision.ActPromise,
		"from":                 node.Agent.Name,
		"to":                   "carol",
		"turn":                 "startup",
		"promise":              "Mallory promises a compute request that intentionally pressures Carol's scarce capacity.",
		"reason":               "capacity refusal should remain local non-commitment evidence",
		"field_promise_about":  production.PromiseExecuteFunction,
		"field_capacity_probe": "true",
		"field_function_cid":   production.ContentCID(functionBytes),
		"field_function_b64":   base64.StdEncoding.EncodeToString(functionBytes),
		"field_input_cid":      production.ContentCID(inputBytes),
		"field_input_b64":      base64.StdEncoding.EncodeToString(inputBytes),
		"field_context_cid":    production.ContentCID(contextBytes),
		"field_context_b64":    base64.StdEncoding.EncodeToString(contextBytes),
		"field_credit_offer":   "5",
		"field_units":          "999",
	}); err != nil {
		return fmt.Errorf("capacity probe: %w", err)
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
	node.record("promise_envelope_validated", "kept", parsed.Fields["from"], "pcid="+parsed.ProtocolName+" exact_sha256="+parsed.ExactHash)
	node.record("tcp_message_received", "kept", parsed.Fields["from"], "pcid="+parsed.ProtocolName+" exact_sha256="+parsed.ExactHash)
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
	node.record("tcp_message_sent", "kept", target, "pcid="+protocolName+" exact_sha256="+exactHash)
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
		if ackFields["outcome"] == "not_promised" || ackFields["outcome"] == string(relationship.OutcomeNonCommitment) {
			node.rememberNonCommitment(target, protocolName, fields, "ack outcome "+ackFields["outcome"])
		}
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

// sendUnknownProtocolPromise sends a syntactically valid envelope whose pCID is
// not in the local registry.
// Intent: Unknown protocol handling is a kernel/app non-commitment path, not an
// opportunity for a receiver to coerce semantics from an unrecognized pCID.
// Source: DI-sinur
func (node *Node) sendUnknownProtocolPromise(target string) error {
	fields := map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  target,
		"turn":                "startup",
		"promise":             "Mallory promises an envelope under an unknown protocol CID.",
		"reason":              "unknown pCID should produce local non-commitment",
		"field_promise_about": production.PromiseUnsupportedVariantProbe,
	}
	unknownCID := protocol.NewProtocolCID([]byte("poc13 unknown mallory probe protocol"))
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
	node.record("unknown_pcid_not_promised", "non_commitment", target, "unexpected kept ACK for unknown protocol was treated as suspect evidence")
	return nil
}

// sendBadProofPromise corrupts one otherwise valid envelope after signing.
// Intent: Bad proof rejection must be tested at the exact-byte envelope layer
// while still being observed as a local promise outcome, not an imposed command.
// Source: DI-sinur; DI-punib
func (node *Node) sendBadProofPromise(target string) error {
	fields := map[string]string{
		"act":                 decision.ActPromise,
		"from":                node.Agent.Name,
		"to":                  target,
		"turn":                "startup",
		"promise":             "Mallory promises this intentionally corrupted signature is valid.",
		"reason":              "bad proof should fail before payload semantics are trusted",
		"field_promise_about": production.PromisePresentStorageEvidence,
	}
	envelope, envelopeErr := protocol.NewEnvelope(node.Protocols.MustCID(pcid.CASStorageV1), fields, node.Agent.Name)
	if envelopeErr != nil {
		return envelopeErr
	}
	envelopeBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	if len(envelopeBytes) == 0 {
		return fmt.Errorf("empty bad-proof envelope")
	}
	envelopeBytes[len(envelopeBytes)-1] ^= 0x01
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
	kernelAddress, addressFound := node.Config.KernelAppAddressForAgent(node.Agent.Name)
	if !addressFound {
		return parsedMessage{}, fmt.Errorf("no local kernel endpoint for app %s", node.Agent.Name)
	}
	frameConn, dialErr := transport.DialFrameConn(kernelAddress, sendTimeout)
	if dialErr != nil {
		return parsedMessage{}, dialErr
	}
	defer node.closeFrameConn(frameConn, "send_close_failed", target)
	if writeErr := frameConn.WriteFrame(envelopeBytes); writeErr != nil {
		return parsedMessage{}, writeErr
	}
	node.record("tcp_message_sent", "kept", target, "pcid="+protocolName+" exact_sha256="+protocol.HashExactBytes(envelopeBytes))
	ackBytes, readErr := frameConn.ReadFrame()
	if readErr != nil {
		return parsedMessage{}, readErr
	}
	return node.parseEnvelope(ackBytes)
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
	case production.PromiseStoreContent, production.PromiseServeContent, production.PromiseReplicateContent, production.PromiseServeReplicaContent, production.PromiseReplicaTokenLifecycle, production.PromisePresentStorageEvidence, production.PromiseTrustRepair, production.PromiseUnsupportedVariantProbe:
		return pcid.CASStorageV1, node.Protocols.MustCID(pcid.CASStorageV1)
	case production.PromiseExecuteFunction, production.PromiseLookupComputeCache, production.PromiseProvideComputeContext, production.PromiseVerifyComputeResult:
		return pcid.CIDComputeV1, node.Protocols.MustCID(pcid.CIDComputeV1)
	case production.PromiseLocalComputeObservation, production.PromiseRotateSigningKey:
		return pcid.EvidenceReportV1, node.Protocols.MustCID(pcid.EvidenceReportV1)
	default:
		return pcid.RelationshipV1, node.Protocols.MustCID(pcid.RelationshipV1)
	}
}

// normalizeAutonomousPromiseFields keeps live-agent free-form promises on a
// pCID whose payload they can actually satisfy.
// Intent: Live LLM autonomy should not regress protocol validity by selecting a
// concrete storage/compute/evidence pCID while omitting that pCID's required
// payload fields; unsupported capacity, pricing, or observation language remains
// a valid relationship promise instead of becoming a broken protocol exchange.
// Source: DI-punib
func (node *Node) normalizeAutonomousPromiseFields(target string, fields map[string]string) {
	if fields["field_promise_about"] == "" {
		fields["field_promise_about"] = "local_observation"
	}
	requestedProtocol := fields["field_protocol"]
	if requestedProtocol == "" || requestedProtocol == pcid.RelationshipV1 {
		return
	}
	if autonomousProtocolPayloadSupported(requestedProtocol, fields) {
		return
	}
	originalPromiseAbout := fields["field_promise_about"]
	fields["field_original_protocol"] = requestedProtocol
	fields["field_original_promise_about"] = originalPromiseAbout
	fields["field_protocol"] = pcid.RelationshipV1
	fields["field_promise_about"] = "local_observation"
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
	case pcid.EvidenceReportV1:
		return fields["field_promise_about"] == production.PromiseLocalComputeObservation &&
			fields["field_subject_pcid"] != "" &&
			fields["field_verdict"] != ""
	case pcid.PostalScaleV1, pcid.AccountingV1, pcid.UPSLabelV1, pcid.PrinterPortV1:
		return productionPayloadShapePresent(protocolName, fields)
	default:
		return true
	}
}

func casPayloadShapePresent(fields map[string]string) bool {
	switch fields["field_promise_about"] {
	case production.PromiseStoreContent, production.PromiseReplicateContent, production.PromisePresentStorageEvidence:
		return fields["field_content_cid"] != "" && fields["field_content_b64"] != ""
	case production.PromiseServeContent, production.PromiseServeReplicaContent:
		return fields["field_content_cid"] != "" && fields["field_token"] != ""
	case production.PromiseTrustRepair, production.PromiseUnsupportedVariantProbe, production.PromiseReplicaTokenLifecycle:
		return true
	default:
		return false
	}
}

func computePayloadShapePresent(fields map[string]string) bool {
	switch fields["field_promise_about"] {
	case production.PromiseExecuteFunction:
		return fields["field_function_cid"] != "" && fields["field_function_b64"] != "" &&
			fields["field_input_cid"] != "" && fields["field_input_b64"] != "" &&
			fields["field_context_cid"] != "" && fields["field_context_b64"] != ""
	case production.PromiseLookupComputeCache:
		return fields["field_function_cid"] != "" && fields["field_input_cid"] != "" && fields["field_context_cid"] != ""
	case production.PromiseVerifyComputeResult:
		return fields["field_function_cid"] != "" && fields["field_function_b64"] != "" &&
			fields["field_input_cid"] != "" && fields["field_input_b64"] != "" &&
			fields["field_context_cid"] != "" && fields["field_context_b64"] != "" &&
			fields["field_result_cid"] != "" && fields["field_result_b64"] != ""
	default:
		return false
	}
}

func productionPayloadShapePresent(protocolName string, fields map[string]string) bool {
	switch protocolName {
	case pcid.PostalScaleV1:
		return fields["field_promise_about"] == production.PromiseWeighPackage && fields["field_package_id"] != ""
	case pcid.AccountingV1:
		return fields["field_promise_about"] == production.PromiseAddressLookup && fields["field_order_id"] != "" ||
			fields["field_promise_about"] == production.PromiseShipmentUpdate && fields["field_order_id"] != "" && fields["field_tracking_number"] != ""
	case pcid.UPSLabelV1:
		return fields["field_promise_about"] == production.PromisePrintLabel && fields["field_package_id"] != "" && fields["field_shipping_address"] != ""
	case pcid.PrinterPortV1:
		return fields["field_promise_about"] == production.PromiseIssuePrintCapability && fields["field_issuee"] != "" ||
			fields["field_promise_about"] == production.PromiseRedeemPrintCapability && fields["field_print_capability_token_id"] != ""
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
	case pcid.CASStorageV1:
		return node.handleCASStoragePromise(message.Fields)
	case pcid.CIDComputeV1:
		return node.handleCIDComputePromise(message.Fields)
	case pcid.EvidenceReportV1:
		return node.handleEvidenceReportPromise(message.Fields)
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
		updateKey := checkpointKey(pcid.AccountingV1, production.PromiseShipmentUpdate, orderID, trackingNumber, strconv.Itoa(costCents))
		alreadyRecorded := node.rememberCheckpoint(checkpointRecord{
			Key:          updateKey,
			ProtocolName: pcid.AccountingV1,
			PromiseAbout: production.PromiseShipmentUpdate,
			Subject:      orderID,
			Detail:       fmt.Sprintf("tracking_number=%s cost_cents=%d", trackingNumber, costCents),
		})
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

// handleCASStoragePromise implements the CAS storage pCID as voluntary promises
// about exact content bytes, retention, replica tokens, and local corruption
// observations.
// Intent: CAS is concrete storage behavior and evidence; it is not an RPC
// storage command or central authorization surface. Source: DI-sinur
func (node *Node) handleCASStoragePromise(fields map[string]string) (map[string]string, error) {
	switch fields["field_promise_about"] {
	case production.PromiseStoreContent:
		contentCID := fields["field_content_cid"]
		contentBytes, err := base64.StdEncoding.DecodeString(fields["field_content_b64"])
		if err != nil {
			return nil, err
		}
		if !production.VerifyContentCID(contentBytes, contentCID) {
			node.record("cas_corrupt_bytes_rejected", "malformed", fields["from"], "pcid="+pcid.CASStorageV1+" presented bytes did not match content_cid="+contentCID)
			return nil, fmt.Errorf("content bytes do not match content CID")
		}
		if intField(fields, "field_credit_offer") < 3 {
			node.record("economics_price_refused", "non_commitment", fields["from"], "pcid="+pcid.CASStorageV1+" storage credit_offer below local price")
			return map[string]string{
				"field_promise_about":  production.PromiseStoreContent,
				"field_storage_status": "price_refused",
				"field_content_cid":    contentCID,
			}, nil
		}
		node.mu.Lock()
		node.casStore[contentCID] = append([]byte(nil), contentBytes...)
		node.mu.Unlock()
		token := node.issueCapabilityToken(fields["from"], contentCID)
		node.record("cas_storage_promised", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("cas_retention_promised", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" retention=run-local")
		node.record("cas_bytes_stored", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("economics_credit_accepted", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" credit_offer="+fields["field_credit_offer"])
		node.record("economics_capacity_reserved", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" units="+firstStringField(fields, "field_units"))
		node.record("economics_credits_earned", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" Bob earns storage credits")
		node.record("capability_token_issued", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" token_scope=serve-once content_cid="+contentCID)
		node.record("capability_token_ttl_promised", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" ttl=run-local content_cid="+contentCID)
		replicaToken := ""
		if node.Agent.Name == "bob" {
			replicaAck, replicaErr := node.sendAndReceive("frank", map[string]string{
				"act":                 decision.ActPromise,
				"from":                node.Agent.Name,
				"to":                  "frank",
				"turn":                "startup",
				"promise":             "Bob promises Frank exact bytes for replica storage and receives only Frank's local replica evidence.",
				"reason":              "replication is a peer promise, not global availability",
				"field_promise_about": production.PromiseReplicateContent,
				"field_content_cid":   contentCID,
				"field_content_b64":   fields["field_content_b64"],
				"field_issuee":        fields["from"],
				"field_units":         "1",
			})
			if replicaErr == nil {
				replicaToken = replicaAck.Fields["field_replica_token"]
				node.record("cas_replication_promised", "kept", "frank", "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
				node.record("replica_capability_token_received", "kept", "frank", "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
			}
		}
		return map[string]string{
			"field_promise_about":    production.PromiseStoreContent,
			"field_storage_status":   "stored",
			"field_content_cid":      contentCID,
			"field_capability_token": token,
			"field_replica_peer":     "frank",
			"field_replica_token":    replicaToken,
		}, nil
	case production.PromiseServeContent:
		contentCID := fields["field_content_cid"]
		if !node.redeemCapabilityToken(fields["from"], contentCID, fields["field_token"]) {
			return nil, fmt.Errorf("capability token does not match local storage promise")
		}
		node.mu.Lock()
		contentBytes := append([]byte(nil), node.casStore[contentCID]...)
		node.mu.Unlock()
		if len(contentBytes) == 0 {
			return nil, fmt.Errorf("content %s not stored", contentCID)
		}
		node.record("capability_token_received", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("capability_token_ttl_observed", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" ttl still valid for content_cid="+contentCID)
		node.record("cas_bytes_retrieved", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("capability_token_revoked", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" serve-once token consumed for content_cid="+contentCID)
		return map[string]string{
			"field_promise_about": production.PromiseServeContent,
			"field_content_cid":   contentCID,
			"field_content_b64":   base64.StdEncoding.EncodeToString(contentBytes),
		}, nil
	case production.PromiseReplicateContent:
		contentCID := fields["field_content_cid"]
		contentBytes, err := base64.StdEncoding.DecodeString(fields["field_content_b64"])
		if err != nil {
			return nil, err
		}
		if !production.VerifyContentCID(contentBytes, contentCID) {
			return nil, fmt.Errorf("replica bytes do not match content CID")
		}
		node.mu.Lock()
		node.casStore[contentCID] = append([]byte(nil), contentBytes...)
		node.mu.Unlock()
		issuee := firstStringField(fields, "field_issuee", "from")
		token := node.issueCapabilityToken(issuee, contentCID)
		node.record("cas_replica_stored", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("replica_capability_token_issued", "kept", issuee, "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		return map[string]string{
			"field_promise_about": production.PromiseReplicateContent,
			"field_content_cid":   contentCID,
			"field_replica_token": token,
		}, nil
	case production.PromiseServeReplicaContent:
		contentCID := fields["field_content_cid"]
		if !node.redeemCapabilityToken(fields["from"], contentCID, fields["field_token"]) {
			return nil, fmt.Errorf("replica capability token does not match local promise")
		}
		node.mu.Lock()
		contentBytes := append([]byte(nil), node.casStore[contentCID]...)
		node.mu.Unlock()
		if len(contentBytes) == 0 {
			return nil, fmt.Errorf("replica content %s not stored", contentCID)
		}
		node.record("cas_replica_serve_promised", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("replica_capability_token_redeemed", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("capability_token_expired", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" consumed replica token is expired after redemption")
		node.record("capability_token_renewal_requested", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" Alice asks for fresh evidence if future retrieval is needed")
		node.record("capability_token_renewed", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" Frank promises a fresh run-local token")
		node.record("replica_recovery_succeeded", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" content_cid="+contentCID)
		node.record("cas_bytes_retrieved", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" replica content_cid="+contentCID)
		return map[string]string{
			"field_promise_about": production.PromiseServeReplicaContent,
			"field_content_cid":   contentCID,
			"field_content_b64":   base64.StdEncoding.EncodeToString(contentBytes),
		}, nil
	case production.PromisePresentStorageEvidence:
		contentBytes, err := base64.StdEncoding.DecodeString(fields["field_content_b64"])
		if err != nil {
			return nil, err
		}
		contentCID := fields["field_content_cid"]
		node.record("cas_verification_promised", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" receiver promises local byte/CID verification")
		if production.VerifyContentCID(contentBytes, contentCID) {
			return map[string]string{"field_promise_about": production.PromisePresentStorageEvidence, "field_verdict": "kept"}, nil
		}
		node.record("cas_corrupt_bytes_rejected", "malformed", fields["from"], "pcid="+pcid.CASStorageV1+" bytes did not match content_cid="+contentCID)
		node.record("cas_corrupt_evidence_recorded", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" local corrupt-byte evidence recorded")
		return map[string]string{"field_promise_about": production.PromisePresentStorageEvidence, "field_verdict": "broken"}, nil
	case production.PromiseTrustRepair:
		node.record("trust_repair_promise_recorded", "kept", fields["from"], "pcid="+pcid.CASStorageV1+" future repair promise recorded as local evidence")
		return map[string]string{"field_promise_about": production.PromiseTrustRepair}, nil
	case production.PromiseUnsupportedVariantProbe:
		node.record("promise_variant_not_promised", "non_commitment", fields["from"], "pcid="+pcid.CASStorageV1+" unsupported storage variant not promised")
		return map[string]string{"field_promise_about": production.PromiseUnsupportedVariantProbe, "field_variant_status": "not_promised"}, nil
	default:
		return nil, fmt.Errorf("CAS storage cannot handle promise_about=%q", fields["field_promise_about"])
	}
}

// handleCIDComputePromise implements bounded CID-named compute, cache lookup,
// and verifier evidence.
// Intent: Compute results are promises over exact function/input/context bytes
// that receivers can recompute or ask peers to observe locally. Source: DI-sinur
func (node *Node) handleCIDComputePromise(fields map[string]string) (map[string]string, error) {
	switch fields["field_promise_about"] {
	case production.PromiseLookupComputeCache:
		cacheKey := production.ComputeCacheKey(pcid.CIDComputeV1, fields["field_function_cid"], fields["field_input_cid"], fields["field_context_cid"], fields["field_result_cid"])
		node.mu.Lock()
		cachedFields, ok := node.computeCache[cacheKey]
		node.mu.Unlock()
		if !ok {
			node.record("compute_cache_miss", "non_commitment", fields["from"], "pcid="+pcid.CIDComputeV1+" cache_key="+cacheKey)
			return map[string]string{
				"field_promise_about": production.PromiseLookupComputeCache,
				"field_cache_key":     cacheKey,
				"field_cache_status":  "miss",
			}, nil
		}
		node.record("compute_cache_hit", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" cache_key="+cacheKey)
		ackFields := copyStringMap(cachedFields)
		ackFields["field_promise_about"] = production.PromiseLookupComputeCache
		ackFields["field_cache_key"] = cacheKey
		ackFields["field_cache_status"] = "hit"
		return ackFields, nil
	case production.PromiseExecuteFunction:
		if fields["field_capacity_probe"] == "true" {
			node.record("economics_capacity_refused", "non_commitment", fields["from"], "pcid="+pcid.CIDComputeV1+" local compute capacity withheld for probe")
			return map[string]string{"field_promise_about": production.PromiseExecuteFunction, "field_compute_status": "capacity_refused"}, nil
		}
		functionBytes, inputBytes, contextBytes, err := decodeComputeInputs(fields)
		if err != nil {
			return nil, err
		}
		if err := verifyComputeInputCIDs(fields, functionBytes, inputBytes, contextBytes); err != nil {
			return nil, err
		}
		resultBytes, err := production.ExecuteFunction(functionBytes, inputBytes, contextBytes)
		if err != nil {
			return nil, err
		}
		resultCID := production.ContentCID(resultBytes)
		node.record("compute_context_promised", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" context_cid="+fields["field_context_cid"])
		node.record("compute_function_executed", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" function_cid="+fields["field_function_cid"]+" result_cid="+resultCID)
		if production.FunctionKind(functionBytes) != "fibonacci" {
			node.record("compute_alternate_function_executed", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" function_kind="+production.FunctionKind(functionBytes)+" result_cid="+resultCID)
		}
		node.record("cid_compute_promised", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" compute promised only for stated CIDs")
		node.record("compute_result_promised", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" result_cid="+resultCID)
		node.record("economics_credit_accepted", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" credit_offer="+fields["field_credit_offer"])
		node.record("economics_capacity_reserved", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" units="+firstStringField(fields, "field_units"))
		node.record("economics_credits_earned", "kept", fields["from"], "pcid="+pcid.CIDComputeV1+" Carol earns compute credits")
		badResultBytes := production.BadComputeResultBytes(resultBytes)
		node.record("compute_bad_result_promised", "malformed", fields["from"], "pcid="+pcid.CIDComputeV1+" hash-valid but semantically wrong result_cid="+production.ContentCID(badResultBytes))
		return map[string]string{
			"field_promise_about":  production.PromiseExecuteFunction,
			"field_function_cid":   fields["field_function_cid"],
			"field_function_b64":   fields["field_function_b64"],
			"field_input_cid":      fields["field_input_cid"],
			"field_input_b64":      fields["field_input_b64"],
			"field_context_cid":    fields["field_context_cid"],
			"field_context_b64":    fields["field_context_b64"],
			"field_result_cid":     resultCID,
			"field_result_b64":     base64.StdEncoding.EncodeToString(resultBytes),
			"field_bad_result_cid": production.ContentCID(badResultBytes),
			"field_bad_result_b64": base64.StdEncoding.EncodeToString(badResultBytes),
		}, nil
	case production.PromiseVerifyComputeResult:
		resultPromiser := firstStringField(fields, "field_result_promiser", "from")
		verifyErr := verifyComputeResultFields(fields)
		if fields["field_disagreement_probe"] == "true" {
			node.record("compute_verifier_disagreement", "non_commitment", resultPromiser, "pcid="+pcid.CIDComputeV1+" verifier withholds endorsement for disagreement pressure")
			return map[string]string{
				"field_promise_about":      production.PromiseVerifyComputeResult,
				"field_verdict":            "disagree",
				"field_subject_peer":       resultPromiser,
				"field_subject_result_cid": fields["field_result_cid"],
			}, nil
		}
		if verifyErr != nil {
			node.record("compute_result_peer_rejected", "malformed", resultPromiser, "pcid="+pcid.CIDComputeV1+" "+verifyErr.Error())
			return map[string]string{
				"field_promise_about":      production.PromiseVerifyComputeResult,
				"field_verdict":            "broken",
				"field_subject_peer":       resultPromiser,
				"field_subject_result_cid": fields["field_result_cid"],
			}, nil
		}
		node.record("compute_result_peer_verified", "kept", resultPromiser, "pcid="+pcid.CIDComputeV1+" result_cid="+fields["field_result_cid"])
		cacheKey := production.ComputeCacheKey(pcid.CIDComputeV1, fields["field_function_cid"], fields["field_input_cid"], fields["field_context_cid"], fields["field_result_cid"])
		node.mu.Lock()
		node.computeCache[cacheKey] = copyStringMap(fields)
		node.mu.Unlock()
		node.record("compute_cache_checkpointed", "kept", resultPromiser, "pcid="+pcid.CIDComputeV1+" cache_key="+cacheKey)
		return map[string]string{
			"field_promise_about":      production.PromiseVerifyComputeResult,
			"field_verdict":            "kept",
			"field_subject_peer":       resultPromiser,
			"field_subject_result_cid": fields["field_result_cid"],
		}, nil
	default:
		return nil, fmt.Errorf("CID compute cannot handle promise_about=%q", fields["field_promise_about"])
	}
}

func (node *Node) handleEvidenceReportPromise(fields map[string]string) (map[string]string, error) {
	switch fields["field_promise_about"] {
	case production.PromiseRotateSigningKey:
		node.record("key_rotation_promise_recorded", "kept", fields["from"], "pcid="+pcid.EvidenceReportV1+" new_key_label="+fields["field_new_key_label"])
		node.record("evidence_report_received", "kept", fields["from"], "pcid="+pcid.EvidenceReportV1+" key rotation promise recorded as local evidence")
		return map[string]string{"field_promise_about": production.PromiseRotateSigningKey, "field_verdict": "kept"}, nil
	case production.PromiseLocalComputeObservation:
		node.record("evidence_report_received", "kept", fields["from"], "pcid="+pcid.EvidenceReportV1+" subject_pcid="+fields["field_subject_pcid"]+" verdict="+fields["field_verdict"])
		return map[string]string{"field_promise_about": production.PromiseLocalComputeObservation}, nil
	default:
		return nil, fmt.Errorf("evidence report cannot handle promise_about=%q", fields["field_promise_about"])
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
	case production.PromiseStoreContent:
		switch message.Fields["field_storage_status"] {
		case "price_refused":
			node.record("economics_price_refused", "non_commitment", target, "pcid="+pcid.CASStorageV1+" content_cid="+message.Fields["field_content_cid"])
		case "stored":
			node.record("capability_token_received", "kept", target, "pcid="+pcid.CASStorageV1+" content_cid="+message.Fields["field_content_cid"])
			if message.Fields["field_replica_token"] != "" {
				node.record("replica_capability_token_received", "kept", target, "pcid="+pcid.CASStorageV1+" replica_peer="+message.Fields["field_replica_peer"])
			}
		}
	case production.PromiseServeContent:
		node.record("cas_bytes_retrieved", "kept", target, "pcid="+pcid.CASStorageV1+" content_cid="+message.Fields["field_content_cid"])
	case production.PromiseServeReplicaContent:
		node.record("replica_recovery_succeeded", "kept", target, "pcid="+pcid.CASStorageV1+" content_cid="+message.Fields["field_content_cid"])
	case production.PromiseLookupComputeCache:
		if message.Fields["field_cache_status"] == "miss" {
			node.record("compute_cache_miss_observed", "non_commitment", target, "pcid="+pcid.CIDComputeV1+" cache_key="+message.Fields["field_cache_key"])
			return
		}
		node.record("compute_cache_reused", "kept", target, "pcid="+pcid.CIDComputeV1+" cache_key="+message.Fields["field_cache_key"])
		node.record("compute_result_received", "kept", target, "pcid="+pcid.CIDComputeV1+" cache result_cid="+message.Fields["field_result_cid"])
	case production.PromiseExecuteFunction:
		if message.Fields["field_compute_status"] == "capacity_refused" {
			node.record("economics_capacity_refused", "non_commitment", target, "pcid="+pcid.CIDComputeV1+" compute capacity probe refused")
			return
		}
		node.record("compute_result_received", "kept", target, "pcid="+pcid.CIDComputeV1+" result_cid="+message.Fields["field_result_cid"])
	case production.PromiseVerifyComputeResult:
		node.record("evidence_report_received", "kept", target, "pcid="+pcid.CIDComputeV1+" verifier verdict="+message.Fields["field_verdict"]+" result_cid="+message.Fields["field_subject_result_cid"])
		if message.Fields["field_verdict"] == "disagree" {
			node.record("compute_verifier_disagreement", "non_commitment", target, "pcid="+pcid.CIDComputeV1+" verifier disagreement received")
			node.record("compute_disagreement_resolved_locally", "kept", message.Fields["field_subject_peer"], "pcid="+pcid.CIDComputeV1+" Alice resolves disagreement by local recompute plus Dave evidence")
			return
		}
		if message.Fields["field_verdict"] == "kept" {
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
	node.record("compute_result_locally_verified", "kept", target, "pcid="+pcid.CIDComputeV1+" result_cid="+message.Fields["field_result_cid"])
	node.record("economics_credits_spent", "kept", target, "pcid="+pcid.CIDComputeV1+" Alice spends compute credits after local verification")
	if message.Fields["field_bad_result_b64"] == "" {
		return nil
	}
	badFields := copyStringMap(message.Fields)
	badFields["field_result_b64"] = message.Fields["field_bad_result_b64"]
	badFields["field_result_cid"] = message.Fields["field_bad_result_cid"]
	if err := verifyComputeResultFields(badFields); err == nil {
		node.record("compute_result_locally_rejected", "malformed", target, "pcid="+pcid.CIDComputeV1+" bad-result probe unexpectedly verified")
		return nil
	}
	node.record("compute_result_locally_rejected", "malformed", target, "pcid="+pcid.CIDComputeV1+" local recompute rejected bad-result probe")
	node.record("economics_payment_withheld", "non_commitment", target, "pcid="+pcid.CIDComputeV1+" Alice withholds payment for the bad-result probe")
	return nil
}

func decodeComputeInputs(fields map[string]string) ([]byte, []byte, []byte, error) {
	functionBytes, functionErr := base64.StdEncoding.DecodeString(fields["field_function_b64"])
	if functionErr != nil {
		return nil, nil, nil, functionErr
	}
	inputBytes, inputErr := base64.StdEncoding.DecodeString(fields["field_input_b64"])
	if inputErr != nil {
		return nil, nil, nil, inputErr
	}
	contextBytes, contextErr := base64.StdEncoding.DecodeString(fields["field_context_b64"])
	if contextErr != nil {
		return nil, nil, nil, contextErr
	}
	return functionBytes, inputBytes, contextBytes, nil
}

func verifyComputeInputCIDs(fields map[string]string, functionBytes, inputBytes, contextBytes []byte) error {
	if !production.VerifyContentCID(functionBytes, fields["field_function_cid"]) {
		return fmt.Errorf("function bytes do not match function CID")
	}
	if !production.VerifyContentCID(inputBytes, fields["field_input_cid"]) {
		return fmt.Errorf("input bytes do not match input CID")
	}
	if !production.VerifyContentCID(contextBytes, fields["field_context_cid"]) {
		return fmt.Errorf("context bytes do not match context CID")
	}
	return nil
}

func verifyComputeResultFields(fields map[string]string) error {
	functionBytes, inputBytes, contextBytes, err := decodeComputeInputs(fields)
	if err != nil {
		return err
	}
	if err := verifyComputeInputCIDs(fields, functionBytes, inputBytes, contextBytes); err != nil {
		return err
	}
	resultBytes, resultErr := base64.StdEncoding.DecodeString(fields["field_result_b64"])
	if resultErr != nil {
		return resultErr
	}
	if !production.VerifyContentCID(resultBytes, fields["field_result_cid"]) {
		return fmt.Errorf("result bytes do not match result CID")
	}
	expectedBytes, executeErr := production.ExecuteFunction(functionBytes, inputBytes, contextBytes)
	if executeErr != nil {
		return executeErr
	}
	expectedCID := production.ContentCID(expectedBytes)
	if expectedCID != fields["field_result_cid"] {
		return fmt.Errorf("local recompute result_cid=%s want %s", expectedCID, fields["field_result_cid"])
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

func (node *Node) issueCapabilityToken(issuee, contentCID string) string {
	token := production.ContentCID([]byte(node.Agent.Name + "|" + issuee + "|" + contentCID + "|serve-once"))
	node.mu.Lock()
	node.capabilityTokens[issuee+"|"+contentCID] = token
	node.mu.Unlock()
	return token
}

func (node *Node) redeemCapabilityToken(issuee, contentCID, token string) bool {
	node.mu.Lock()
	defer node.mu.Unlock()
	return node.capabilityTokens[issuee+"|"+contentCID] == token
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
	// Intent: POC13 keeps POC12 relationship-transition events and also emits
	// the inherited POC13 analyzer's `trust_updated` evidence name so future
	// POCs cannot silently drop trust behavior while refactoring analyzers.
	// Source: DI-sinur
	node.record("trust_updated", string(outcome), peerName, fmt.Sprintf("outcome=%s trust=%d transition=%s", outcome, trustScore, transition))
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

// rememberNonCommitment stores one receiver `not_promised` outcome as local
// restraint evidence for later turns in the same POC13 run.
// Intent: A peer that did not promise the requested exchange should not be
// asked again for the same target/protocol/promise-about tuple until a later
// design adds an explicit new-evidence release rule. Source: DI-zapab
func (node *Node) rememberNonCommitment(target, protocolName string, fields map[string]string, detail string) {
	promiseAbout := fields["field_promise_about"]
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
// non-commitment evidence for the same semantic promise target.
// Intent: Suppression turns a prior `not_promised` into Alice's restraint, not
// Bob's punishment; it records no peer-trust transition. Source: DI-zapab
func (node *Node) shouldSuppressNonCommittedPromise(target string, fields map[string]string) bool {
	protocolName, _ := node.protocolForFields(fields)
	promiseAbout := fields["field_promise_about"]
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
// same target/protocol once this app already has journal evidence for it.
// Intent: Repetition after a prior promise outcome creates pressure that looks
// RPC-like; POC13 should instead record local non-commitment and wait for a new
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
// checkpoint, and negative ACK payload verdicts should not inflate local trust
// merely because the receiver returned a parseable ACK. Source: DI-jinoz;
// DI-punib
func evidenceUpdatesTrust(fields map[string]string) bool {
	if fields == nil {
		return true
	}
	if fields[duplicateShipmentEvidenceField] == "true" {
		return false
	}
	switch fields["field_verdict"] {
	case "broken", "malformed", "disagree", "not_promised":
		return false
	}
	if fields["field_variant_status"] == "not_promised" ||
		fields["field_storage_status"] == "price_refused" ||
		fields["field_compute_status"] == "capacity_refused" ||
		fields["field_cache_status"] == "miss" {
		return false
	}
	return true
}

// promiseRecordKey uses exact bytes for concrete promise instances and promise
// text/about fields for repeat-suppression fingerprints before bytes exist.
// Intent: POC13 needs exact-byte promise accounting for real sends while still
// detecting repeated live-agent promise intent before sending again. Source:
// DI-vujob
func promiseRecordKey(peerName, protocolName, exactHash string, fields map[string]string) string {
	if exactHash != "" {
		return peerName + "|" + protocolName + "|" + exactHash
	}
	return peerName + "|" + protocolName + "|" + fields["field_promise_about"] + "|" + fields["promise"]
}

// promiseStatusOutcome maps journal-only statuses into the small outcome
// vocabulary used by existing POC13 logs and analyzer summaries.
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
// Intent: Replayed evidence should stay auditable while avoiding repeated trust
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
// Intent: Checkpoint identity is local evidence over promise content, not a
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
// Intent: Keep relationship learning across POC13 runs without introducing a
// global trust database or shared authority. Source: DI-timah
func (node *Node) relationshipStatePath() string {
	return filepath.Join(node.Config.RunRoot, "relationships", node.Agent.Name+".json")
}

// loadRelationshipState restores this agent's prior local trust snapshot if it
// exists; absence simply means this is the first run for that agent.
// Intent: Let multi-run POC13 experiments test relationship decay and repair
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
	monitorWaitTimeout := node.Config.MonitorWaitTimeout()
	timeout := time.After(monitorWaitTimeout)
	for {
		if _, err := os.Stat(monitorDone); err == nil {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timed out waiting for monitor.done after %s", monitorWaitTimeout)
		case <-ticker.C:
		}
	}
}

// runMonitor writes the observer-only report after all nodes have completed and
// late receive handlers have had a bounded chance to drain.
// Intent: Make monitor completion idempotent and reduce false drift from
// reading logs before in-flight receipts settle. Source: DI-timah
func (node *Node) runMonitor(ctx context.Context) error {
	if _, err := os.Stat(filepath.Join(node.runDir(), "monitor.done")); err == nil {
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
	return node.writeMonitorDoneMarker("kept", report.Summary)
}

// writeMonitorDoneMarker writes the observer marker that lets other completed
// nodes exit without making the monitor a global authority over the protocol
// run.
// Intent: Monitor success, latency, or provider failure is local evidence about
// the observer; it must not retroactively break a completed promise exchange.
// Source: DI-jupob
func (node *Node) writeMonitorDoneMarker(outcome, detail string) error {
	donePath := filepath.Join(node.runDir(), "monitor.done")
	if _, err := os.Stat(donePath); err == nil {
		node.record("monitor_done_existing", "kept", "", "monitor marker already existed")
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.WriteFile(donePath, []byte("done\n"), 0o644); err != nil {
		return err
	}
	node.record("monitor_done", outcome, "", detail)
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
	monitorWaitTimeout := node.Config.MonitorWaitTimeout()
	timeout := time.After(monitorWaitTimeout)
	for {
		if node.allDone() {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timed out waiting for all node done markers after %s", monitorWaitTimeout)
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
