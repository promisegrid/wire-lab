package poc13

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const maxFrameBytes = 1 << 20

const (
	storageCapacity = 2
	computeCapacity = 1
	computePrice    = 5
)

// AgentState is one local agent's mutable evidence state for the bounded POC.
// Intent: POC13 keeps trust, credits, capabilities, and storage local to the
// observing agent instead of creating global authorities. Source: DI-fumol
type AgentState struct {
	Agent          AgentConfig
	Store          map[string][]byte
	Trust          map[string]int
	Credits        map[string]int
	Capabilities   map[string]string
	PendingCompute map[string]map[string]string
	mu             sync.Mutex
}

// OutboundPromise describes one signed promise envelope before it is serialized
// and length-framed onto TCP.
type OutboundPromise struct {
	From         string
	To           string
	ProtocolName string
	Fields       map[string]string
}

// RuntimeMessage is the verified pCID-owned payload plus transport metadata
// observed by a receiving agent.
type RuntimeMessage struct {
	From         string
	To           string
	ProtocolName string
	Fields       map[string]string
	ExactBytes   []byte
}

// FrameReader reads one length-prefixed raw grid envelope from TCP.
type FrameReader struct {
	reader io.Reader
}

// FrameWriter writes one length-prefixed raw grid envelope to TCP.
type FrameWriter struct {
	writer io.Writer
}

// TCPRuntime is the POC13 supervisor runtime for one Docker container.
// Intent: The supervisor owns the local TCP listener and peer sends so POC13
// proves real inter-agent transport without inventing a central router. Source:
// DI-fumol
type TCPRuntime struct {
	Config         Config
	Container      ContainerConfig
	Registry       Registry
	Decider        Decider
	states         map[string]*AgentState
	listener       net.Listener
	listenPort     int
	activeHandlers int
	lastActivity   time.Time
	mu             sync.Mutex
}

// NewTCPRuntime constructs one container-local POC13 runtime.
func NewTCPRuntime(cfg Config, container ContainerConfig, decider Decider) (*TCPRuntime, error) {
	if decider == nil {
		decider = LiveOrScriptedDecider{}
	}
	runtime := &TCPRuntime{
		Config:     cfg,
		Container:  container,
		Registry:   NewRegistry(),
		Decider:    decider,
		states:     make(map[string]*AgentState),
		listenPort: cfg.ListenPort,
	}
	for _, agentName := range container.Agents {
		agent, ok := cfg.Agent(agentName)
		if !ok {
			return nil, fmt.Errorf("unknown agent %s", agentName)
		}
		runtime.states[agentName] = &AgentState{
			Agent:          agent,
			Store:          make(map[string][]byte),
			Trust:          make(map[string]int),
			Credits:        make(map[string]int),
			Capabilities:   make(map[string]string),
			PendingCompute: make(map[string]map[string]string),
		}
	}
	return runtime, nil
}

// Run opens the local listener, records readiness, waits for peer readiness,
// runs local agents, then records done after local runtime quiescence.
// Intent: POC13 startup and shutdown should be evidence-driven instead of
// arbitrary fixed sleeps. Source: DI-mosil
func (runtime *TCPRuntime) Run(ctx context.Context) error {
	listener, listenErr := net.Listen("tcp", runtime.listenAddress())
	if listenErr != nil {
		return listenErr
	}
	runtime.listener = listener
	if tcpAddr, ok := listener.Addr().(*net.TCPAddr); ok {
		runtime.listenPort = tcpAddr.Port
	}
	serveErrs := make(chan error, 1)
	var handlerWG sync.WaitGroup
	go func() {
		serveErrs <- runtime.serve(ctx, &handlerWG)
	}()
	if err := runtime.recordRuntimeReadiness(); err != nil {
		return err
	}
	if err := runtime.waitForPeerReadiness(ctx); err != nil {
		return err
	}
	runErr := runtime.runLocalAgents(ctx)
	if waitErr := runtime.waitForRuntimeCompletion(ctx); waitErr != nil && runErr == nil {
		runErr = waitErr
	}
	if doneErr := runtime.recordRuntimeDone(); doneErr != nil && runErr == nil {
		runErr = doneErr
	}
	closeErr := listener.Close()
	handlerWG.Wait()
	serveErr := <-serveErrs
	if runErr != nil {
		return runErr
	}
	if closeErr != nil && !errors.Is(closeErr, net.ErrClosed) {
		return closeErr
	}
	if serveErr != nil && !errors.Is(serveErr, net.ErrClosed) {
		return serveErr
	}
	return nil
}

func (runtime *TCPRuntime) serve(ctx context.Context, handlerWG *sync.WaitGroup) error {
	for {
		conn, acceptErr := runtime.listener.Accept()
		if acceptErr != nil {
			return acceptErr
		}
		handlerWG.Add(1)
		runtime.mu.Lock()
		runtime.activeHandlers++
		runtime.lastActivity = time.Now()
		runtime.mu.Unlock()
		go func(accepted net.Conn) {
			defer func() {
				runtime.mu.Lock()
				runtime.activeHandlers--
				runtime.lastActivity = time.Now()
				runtime.mu.Unlock()
				handlerWG.Done()
			}()
			if err := runtime.handleConnection(ctx, accepted); err != nil {
				fmt.Fprintf(os.Stderr, "poc13-runtime: handle connection: %v\n", err)
			}
		}(conn)
	}
}

func (runtime *TCPRuntime) handleConnection(ctx context.Context, conn net.Conn) error {
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "poc13-runtime: close connection: %v\n", closeErr)
		}
	}()
	frameBytes, frameErr := (FrameReader{reader: conn}).ReadFrame()
	if frameErr != nil {
		return frameErr
	}
	runtime.markActivity()
	return runtime.handleEnvelope(ctx, frameBytes)
}

// runLocalAgents asks each local agent for a live/local promise judgment and
// lets the judgment gate its initial protocol behavior.
func (runtime *TCPRuntime) runLocalAgents(ctx context.Context) error {
	var wg sync.WaitGroup
	errs := make(chan error, len(runtime.Container.Agents))
	for _, agentName := range runtime.Container.Agents {
		state := runtime.states[agentName]
		wg.Add(1)
		go func(localState *AgentState) {
			defer wg.Done()
			if err := runtime.runLocalAgent(ctx, localState); err != nil {
				errs <- err
			}
		}(state)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func (runtime *TCPRuntime) runLocalAgent(ctx context.Context, state *AgentState) error {
	agent := state.Agent
	if err := runtime.recordAgentEvent(agent.Name, "app_receive_promise_registered", "kept", "", "", "cas_storage_v1 and cid_compute_v1 payloads are accepted only as local promises"); err != nil {
		return err
	}
	decision, decisionErr := runtime.Decider.Decide(ctx, runtime.Config, agent, fmt.Sprintf("Choose one local promise consistent with POC13 role %q and peers %q. Your text gates whether this agent sends protocol promises.", agent.Role, agent.Peers))
	if decisionErr != nil {
		return runtime.recordAgentEvent(agent.Name, "llm_decision_error", "non_commitment", "", "", decisionErr.Error())
	}
	if err := runtime.recordAgentEvent(agent.Name, "llm_decision_"+decision.Mode, "kept", "", "", decision.Text); err != nil {
		return err
	}
	if !decision.Promises() {
		return runtime.recordAgentEvent(agent.Name, "agent_non_commitment", "non_commitment", "", "", "local decision did not contain a voluntary promise")
	}
	switch agent.Role {
	case "data_holder":
		return runtime.runDataHolder(ctx, agent, decision.Text)
	case "adversary_peer":
		return runtime.runAdversaryPeer(ctx, agent, decision.Text)
	default:
		return runtime.recordAgentEvent(agent.Name, "agent_waiting_for_peer_promises", "kept", "", "", "agent waits for peer promise envelopes over TCP")
	}
}

func (runtime *TCPRuntime) runDataHolder(ctx context.Context, agent AgentConfig, decisionText string) error {
	contentBytes := sampleContentBytes()
	contentCID := ContentCID(contentBytes)
	functionBytes := sampleFunctionBytes()
	inputBytes := []byte("n=9")
	contextBytes := sampleContextBytes()
	// Intent: Alice's first actions now expose local trust-driven peer choice
	// and bounded economics before protocol sends, rather than treating Bob and
	// Carol as hard-coded RPC destinations. Source: DI-lupag
	storagePeer := runtime.preferredStoragePeer(agent.Name)
	computePeer := runtime.preferredComputePeer(agent.Name)
	if err := runtime.recordTrustDrivenChoice(agent.Name, storagePeer, CASStorageV1, "Alice chooses Bob as primary storage peer from local trust and plans Frank as replica fallback"); err != nil {
		return err
	}
	if err := runtime.recordTrustDrivenChoice(agent.Name, computePeer, CIDComputeV1, "Alice chooses Carol for compute and plans Dave/Grace verification from prior evidence"); err != nil {
		return err
	}
	if err := runtime.recordEconomics(agent.Name, "economics_price_probe", "kept", storagePeer, CASStorageV1, "Alice first offers low storage credit to test Bob's local price boundary"); err != nil {
		return err
	}
	if err := runtime.sendPromise(ctx, OutboundPromise{
		From:         agent.Name,
		To:           storagePeer,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":           "store_request",
			"promise_about":     "store_content",
			"content_cid":       contentCID,
			"content_b64":       base64.StdEncoding.EncodeToString(contentBytes),
			"credit_offer":      "1",
			"requested_token":   "serve-once",
			"decision_text":     decisionText,
			"promise_condition": "bob may decline locally",
		},
	}); err != nil {
		return err
	}
	if err := runtime.recordEconomics(agent.Name, "economics_credit_offered", "kept", storagePeer, CASStorageV1, "Alice offers storage credit_offer=4 after low-price probe"); err != nil {
		return err
	}
	if err := runtime.sendPromise(ctx, OutboundPromise{
		From:         agent.Name,
		To:           storagePeer,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":           "store_request",
			"promise_about":     "store_content",
			"content_cid":       contentCID,
			"content_b64":       base64.StdEncoding.EncodeToString(contentBytes),
			"credit_offer":      "4",
			"requested_token":   "serve-once",
			"decision_text":     decisionText,
			"promise_condition": "bob may decline locally",
		},
	}); err != nil {
		return err
	}
	if err := runtime.recordEconomics(agent.Name, "economics_credit_offered", "kept", computePeer, CIDComputeV1, "Alice offers compute credit_offer=5 for payload-provided Fibonacci compute"); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         agent.Name,
		To:           computePeer,
		ProtocolName: CIDComputeV1,
		Fields: map[string]string{
			"variant":           "compute_request",
			"promise_about":     "execute_function",
			"function_cid":      ContentCID(functionBytes),
			"function_b64":      base64.StdEncoding.EncodeToString(functionBytes),
			"input_cid":         ContentCID(inputBytes),
			"input_b64":         base64.StdEncoding.EncodeToString(inputBytes),
			"context_cid":       ContentCID(contextBytes),
			"credit_offer":      "5",
			"decision_text":     decisionText,
			"promise_condition": "carol may decline locally",
		},
	})
}

func (runtime *TCPRuntime) runAdversaryPeer(ctx context.Context, agent AgentConfig, decisionText string) error {
	claimedCID := ContentCID(sampleContentBytes())
	if err := runtime.sendPromise(ctx, OutboundPromise{
		From:         agent.Name,
		To:           "grace",
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "corrupt_bytes_offer",
			"promise_about": "present_storage_evidence",
			"content_cid":   claimedCID,
			"content_b64":   base64.StdEncoding.EncodeToString(corruptContentBytes()),
			"decision_text": decisionText,
		},
	}); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         agent.Name,
		To:           "grace",
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "trust_repair_promise",
			"promise_about": "label_future_malformed_evidence",
			"decision_text": "Mallory promises to label future malformed evidence explicitly; Grace remains free to distrust this until future evidence is kept.",
		},
	})
}

func (runtime *TCPRuntime) handleEnvelope(ctx context.Context, exactBytes []byte) error {
	envelope, parseErr := ParseEnvelope(exactBytes)
	if parseErr != nil {
		return parseErr
	}
	if verifyErr := VerifyEnvelope(envelope); verifyErr != nil {
		return verifyErr
	}
	protocolName, known := runtime.Registry.Name(envelope.ProtocolCID)
	if !known {
		return fmt.Errorf("unknown protocol pCID %s", envelope.ProtocolCID.String())
	}
	fields, fieldsErr := envelope.PayloadFields()
	if fieldsErr != nil {
		return fieldsErr
	}
	message := RuntimeMessage{
		From:         fields["sender"],
		To:           fields["recipient"],
		ProtocolName: protocolName,
		Fields:       fields,
		ExactBytes:   exactBytes,
	}
	if message.Fields["act"] != "promise" {
		return fmt.Errorf("received top-level act %q, want promise", message.Fields["act"])
	}
	if _, ok := runtime.states[message.To]; !ok {
		return fmt.Errorf("recipient %s is not local to container %s", message.To, runtime.Container.Name)
	}
	if err := runtime.recordAgentEvent(message.To, "tcp_message_received", "kept", message.From, protocolName, "variant="+message.Fields["variant"]+" exact_sha256="+HashExactBytes(exactBytes)); err != nil {
		return err
	}
	if err := runtime.recordAgentEvent(message.To, "promise_envelope_validated", "kept", message.From, protocolName, "exact_sha256="+HashExactBytes(exactBytes)+" promise_about="+message.Fields["promise_about"]); err != nil {
		return err
	}
	switch message.Fields["variant"] {
	case "store_request":
		return runtime.handleStoreRequest(ctx, message)
	case "price_refusal":
		return runtime.handlePriceRefusal(message)
	case "store_acceptance":
		return runtime.handleStoreAcceptance(ctx, message)
	case "serve_request":
		return runtime.handleServeRequest(ctx, message)
	case "serve_response":
		return runtime.handleServeResponse(message)
	case "primary_unavailable_notice":
		return runtime.handlePrimaryUnavailableNotice(message)
	case "replica_available":
		return runtime.handleReplicaAvailable(ctx, message)
	case "replica_serve_request":
		return runtime.handleReplicaServeRequest(ctx, message)
	case "replicate_request":
		return runtime.handleReplicateRequest(ctx, message)
	case "replicate_acceptance":
		return runtime.handleReplicateAcceptance(ctx, message)
	case "compute_request":
		return runtime.handleComputeRequest(ctx, message)
	case "context_request":
		return runtime.handleContextRequest(ctx, message)
	case "context_response":
		return runtime.handleContextResponse(ctx, message)
	case "compute_result":
		return runtime.handleComputeResult(ctx, message)
	case "compute_verification_request":
		return runtime.handleComputeVerificationRequest(ctx, message)
	case "compute_verification_result":
		return runtime.handleComputeVerificationResult(message)
	case "corrupt_bytes_offer":
		return runtime.handleCorruptBytesOffer(message)
	case "trust_repair_promise":
		return runtime.handleTrustRepairPromise(message)
	default:
		return runtime.recordAgentEvent(message.To, "promise_variant_not_promised", "non_commitment", message.From, protocolName, "variant="+message.Fields["variant"])
	}
}

func (runtime *TCPRuntime) handleStoreRequest(ctx context.Context, message RuntimeMessage) error {
	contentBytes, decodeErr := base64.StdEncoding.DecodeString(message.Fields["content_b64"])
	if decodeErr != nil {
		return decodeErr
	}
	contentCID := message.Fields["content_cid"]
	if !VerifyContentCID(contentBytes, contentCID) {
		runtime.adjustTrust(message.To, message.From, -3)
		return runtime.recordAgentEvent(message.To, "cas_corrupt_bytes_rejected", "malformed", message.From, message.ProtocolName, "store_request bytes did not match content_cid="+contentCID)
	}
	creditOffer := parseCreditOffer(message.Fields["credit_offer"])
	if creditOffer < 3 {
		if err := runtime.recordEconomics(message.To, "economics_price_refused", "non_commitment", message.From, message.ProtocolName, fmt.Sprintf("storage credit_offer=%d below local opportunity_cost=3", creditOffer)); err != nil {
			return err
		}
		return runtime.sendPromise(ctx, OutboundPromise{
			From:         message.To,
			To:           message.From,
			ProtocolName: CASStorageV1,
			Fields: map[string]string{
				"variant":        "price_refusal",
				"promise_about":  "price_boundary",
				"content_cid":    contentCID,
				"minimum_credit": "3",
				"credit_offer":   message.Fields["credit_offer"],
			},
		})
	}
	state := runtime.states[message.To]
	state.mu.Lock()
	if len(state.Store) >= storageCapacity {
		state.mu.Unlock()
		return runtime.recordEconomics(message.To, "economics_capacity_refused", "non_commitment", message.From, message.ProtocolName, fmt.Sprintf("storage capacity exhausted capacity=%d", storageCapacity))
	}
	state.Store[contentCID] = append([]byte(nil), contentBytes...)
	state.Credits[message.From] += creditOffer
	state.mu.Unlock()
	runtime.adjustTrust(message.To, message.From, 1)
	token := runtime.issueCapabilityToken(message.To, message.From, contentCID)
	if err := runtime.recordAgentEvent(message.To, "cas_bytes_stored", "kept", message.From, message.ProtocolName, "stored content_cid="+contentCID); err != nil {
		return err
	}
	if err := runtime.recordAgentEvent(message.To, "cas_storage_promised", "kept", message.From, message.ProtocolName, "Bob promises local storage for content_cid="+contentCID); err != nil {
		return err
	}
	if err := runtime.recordAgentEvent(message.To, "cas_retention_promised", "kept", message.From, message.ProtocolName, "Bob promises bounded retention for content_cid="+contentCID); err != nil {
		return err
	}
	if err := runtime.recordAgentEvent(message.To, "cas_serve_promised", "kept", message.From, message.ProtocolName, "Bob promises token-scoped serving for content_cid="+contentCID); err != nil {
		return err
	}
	if err := runtime.recordEconomics(message.To, "economics_credit_accepted", "kept", message.From, message.ProtocolName, "accepted storage credit_offer="+message.Fields["credit_offer"]); err != nil {
		return err
	}
	if err := runtime.recordEconomics(message.To, "economics_capacity_reserved", "kept", message.From, message.ProtocolName, fmt.Sprintf("reserved storage slot capacity=%d", storageCapacity)); err != nil {
		return err
	}
	if err := runtime.recordEconomics(message.To, "economics_credits_earned", "kept", message.From, message.ProtocolName, fmt.Sprintf("Bob earned storage credits=%d capacity=%d", creditOffer, storageCapacity)); err != nil {
		return err
	}
	if err := runtime.recordAgentEvent(message.To, "capability_token_issued", "kept", message.From, message.ProtocolName, "token="+token+" content_cid="+contentCID); err != nil {
		return err
	}
	if err := runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           message.From,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "store_acceptance",
			"promise_about": "store_content",
			"content_cid":   contentCID,
			"token":         token,
			"retention":     "bounded-poc13-run",
		},
	}); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           "frank",
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "replicate_request",
			"promise_about": "replicate_content",
			"content_cid":   contentCID,
			"content_b64":   base64.StdEncoding.EncodeToString(contentBytes),
			"credit_offer":  "1",
			"requester":     message.From,
		},
	})
}

func (runtime *TCPRuntime) handleStoreAcceptance(ctx context.Context, message RuntimeMessage) error {
	state := runtime.states[message.To]
	state.mu.Lock()
	state.Capabilities[message.Fields["content_cid"]] = message.Fields["token"]
	state.mu.Unlock()
	runtime.adjustTrust(message.To, message.From, 1)
	if err := runtime.recordAgentEvent(message.To, "capability_token_received", "kept", message.From, message.ProtocolName, "token="+message.Fields["token"]+" content_cid="+message.Fields["content_cid"]); err != nil {
		return err
	}
	if err := runtime.recordEconomics(message.To, "economics_credits_spent", "kept", message.From, message.ProtocolName, "Alice spends storage credits=4 after Bob's promise is accepted"); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           message.From,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":              "serve_request",
			"promise_about":        "redeem_storage_capability",
			"content_cid":          message.Fields["content_cid"],
			"token":                message.Fields["token"],
			"simulate_unavailable": "true",
		},
	})
}

func (runtime *TCPRuntime) handlePriceRefusal(message RuntimeMessage) error {
	return runtime.recordEconomics(message.To, "economics_price_refused", "non_commitment", message.From, message.ProtocolName, "peer minimum_credit="+message.Fields["minimum_credit"]+" rejected credit_offer="+message.Fields["credit_offer"])
}

func (runtime *TCPRuntime) handleServeRequest(ctx context.Context, message RuntimeMessage) error {
	contentCID := message.Fields["content_cid"]
	if !runtime.redeemCapabilityToken(message.To, message.From, contentCID, message.Fields["token"]) {
		runtime.adjustTrust(message.To, message.From, -1)
		return runtime.recordAgentEvent(message.To, "capability_token_rejected", "malformed", message.From, message.ProtocolName, "token did not match issuer-local capability for content_cid="+contentCID)
	}
	state := runtime.states[message.To]
	state.mu.Lock()
	contentBytes, ok := state.Store[contentCID]
	state.mu.Unlock()
	if !ok {
		return runtime.recordAgentEvent(message.To, "cas_serve_not_promised", "non_commitment", message.From, message.ProtocolName, "content not present content_cid="+contentCID)
	}
	if err := runtime.recordAgentEvent(message.To, "capability_token_redeemed", "kept", message.From, message.ProtocolName, "token="+message.Fields["token"]+" content_cid="+contentCID); err != nil {
		return err
	}
	if message.Fields["simulate_unavailable"] == "true" {
		if err := runtime.recordAgentEvent(message.To, "primary_storage_unavailable", "non_commitment", message.From, message.ProtocolName, "Bob redeems Alice's token but does not currently promise immediate serving"); err != nil {
			return err
		}
		return runtime.sendPromise(ctx, OutboundPromise{
			From:         message.To,
			To:           message.From,
			ProtocolName: CASStorageV1,
			Fields: map[string]string{
				"variant":       "primary_unavailable_notice",
				"promise_about": "primary_storage_unavailable",
				"content_cid":   contentCID,
			},
		})
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           message.From,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "serve_response",
			"promise_about": "serve_content_bytes",
			"content_cid":   contentCID,
			"content_b64":   base64.StdEncoding.EncodeToString(contentBytes),
		},
	})
}

func (runtime *TCPRuntime) handlePrimaryUnavailableNotice(message RuntimeMessage) error {
	return runtime.recordAgentEvent(message.To, "primary_storage_unavailable", "non_commitment", message.From, message.ProtocolName, "Alice observes Bob unavailable for immediate retrieval and waits for replica evidence")
}

// handleReplicaAvailable lets Alice recover from Bob's current serving
// non-commitment by asking a trusted-enough replica peer.
// Intent: Replica recovery is a local trust decision by Alice, not an automatic
// global failover or availability guarantee. Source: DI-lupag
func (runtime *TCPRuntime) handleReplicaAvailable(ctx context.Context, message RuntimeMessage) error {
	replicaPeer := message.Fields["replica_peer"]
	contentCID := message.Fields["content_cid"]
	if runtime.trustScore(message.To, replicaPeer) < 0 {
		return runtime.recordAgentEvent(message.To, "replica_recovery_not_promised", "non_commitment", replicaPeer, message.ProtocolName, "local trust below replica recovery threshold")
	}
	if err := runtime.recordTrustDrivenChoice(message.To, replicaPeer, message.ProtocolName, "Alice chooses Frank for replica retrieval because Bob is unavailable and Frank's local trust is sufficient"); err != nil {
		return err
	}
	if err := runtime.recordAgentEvent(message.To, "replica_recovery_requested", "kept", replicaPeer, message.ProtocolName, "requesting replica content_cid="+contentCID); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           replicaPeer,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "replica_serve_request",
			"promise_about": "serve_replica_content",
			"content_cid":   contentCID,
		},
	})
}

func (runtime *TCPRuntime) handleReplicaServeRequest(ctx context.Context, message RuntimeMessage) error {
	contentCID := message.Fields["content_cid"]
	state := runtime.states[message.To]
	state.mu.Lock()
	contentBytes, ok := state.Store[contentCID]
	state.mu.Unlock()
	if !ok {
		return runtime.recordAgentEvent(message.To, "cas_replica_serve_not_promised", "non_commitment", message.From, message.ProtocolName, "replica content not present content_cid="+contentCID)
	}
	if err := runtime.recordAgentEvent(message.To, "cas_replica_serve_promised", "kept", message.From, message.ProtocolName, "Frank promises replica serving for content_cid="+contentCID); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           message.From,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "serve_response",
			"promise_about": "serve_content_bytes",
			"content_cid":   contentCID,
			"content_b64":   base64.StdEncoding.EncodeToString(contentBytes),
		},
	})
}

func (runtime *TCPRuntime) handleServeResponse(message RuntimeMessage) error {
	contentBytes, decodeErr := base64.StdEncoding.DecodeString(message.Fields["content_b64"])
	if decodeErr != nil {
		return decodeErr
	}
	contentCID := message.Fields["content_cid"]
	if !VerifyContentCID(contentBytes, contentCID) {
		runtime.adjustTrust(message.To, message.From, -3)
		return runtime.recordAgentEvent(message.To, "cas_corrupt_bytes_rejected", "malformed", message.From, message.ProtocolName, "serve_response bytes did not match content_cid="+contentCID)
	}
	runtime.adjustTrust(message.To, message.From, 1)
	if message.From == "frank" {
		if err := runtime.recordAgentEvent(message.To, "replica_recovery_succeeded", "kept", message.From, message.ProtocolName, "retrieved content from trusted replica content_cid="+contentCID); err != nil {
			return err
		}
	}
	return runtime.recordAgentEvent(message.To, "cas_bytes_retrieved", "kept", message.From, message.ProtocolName, "retrieved content_cid="+contentCID)
}

func (runtime *TCPRuntime) handleReplicateRequest(ctx context.Context, message RuntimeMessage) error {
	contentBytes, decodeErr := base64.StdEncoding.DecodeString(message.Fields["content_b64"])
	if decodeErr != nil {
		return decodeErr
	}
	contentCID := message.Fields["content_cid"]
	if !VerifyContentCID(contentBytes, contentCID) {
		runtime.adjustTrust(message.To, message.From, -3)
		return runtime.recordAgentEvent(message.To, "cas_corrupt_bytes_rejected", "malformed", message.From, message.ProtocolName, "replicate_request bytes did not match content_cid="+contentCID)
	}
	state := runtime.states[message.To]
	state.mu.Lock()
	state.Store[contentCID] = append([]byte(nil), contentBytes...)
	state.Credits[message.From] += 1
	state.mu.Unlock()
	runtime.adjustTrust(message.To, message.From, 1)
	if err := runtime.recordAgentEvent(message.To, "cas_replica_stored", "kept", message.From, message.ProtocolName, "replica content_cid="+contentCID); err != nil {
		return err
	}
	if err := runtime.recordAgentEvent(message.To, "cas_replication_promised", "kept", message.From, message.ProtocolName, "Frank promises one local replica for content_cid="+contentCID); err != nil {
		return err
	}
	if err := runtime.recordEconomics(message.To, "economics_credit_accepted", "kept", message.From, message.ProtocolName, "accepted replication credit_offer="+message.Fields["credit_offer"]); err != nil {
		return err
	}
	if err := runtime.recordEconomics(message.To, "economics_credits_earned", "kept", message.From, message.ProtocolName, "Frank earned replication credits="+message.Fields["credit_offer"]); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           message.From,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "replicate_acceptance",
			"promise_about": "replicate_content",
			"content_cid":   contentCID,
			"requester":     message.Fields["requester"],
		},
	})
}

func (runtime *TCPRuntime) handleReplicateAcceptance(ctx context.Context, message RuntimeMessage) error {
	runtime.adjustTrust(message.To, message.From, 1)
	if err := runtime.recordAgentEvent(message.To, "cas_replication_confirmed", "kept", message.From, message.ProtocolName, "replica accepted content_cid="+message.Fields["content_cid"]); err != nil {
		return err
	}
	requester := message.Fields["requester"]
	if requester == "" {
		return nil
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           requester,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "replica_available",
			"promise_about": "replica_location",
			"content_cid":   message.Fields["content_cid"],
			"replica_peer":  message.From,
		},
	})
}

func (runtime *TCPRuntime) handleComputeRequest(ctx context.Context, message RuntimeMessage) error {
	state := runtime.states[message.To]
	requestID := message.Fields["function_cid"] + "|" + message.Fields["input_cid"]
	creditOffer := parseCreditOffer(message.Fields["credit_offer"])
	if creditOffer < computePrice {
		return runtime.recordEconomics(message.To, "economics_price_refused", "non_commitment", message.From, message.ProtocolName, fmt.Sprintf("compute credit_offer=%d below local opportunity_cost=%d", creditOffer, computePrice))
	}
	state.mu.Lock()
	if len(state.PendingCompute) >= computeCapacity {
		state.mu.Unlock()
		return runtime.recordEconomics(message.To, "economics_capacity_refused", "non_commitment", message.From, message.ProtocolName, fmt.Sprintf("compute capacity exhausted capacity=%d", computeCapacity))
	}
	state.PendingCompute[requestID] = message.Fields
	state.Credits[message.From] += creditOffer
	state.mu.Unlock()
	if err := runtime.recordEconomics(message.To, "economics_credit_accepted", "kept", message.From, message.ProtocolName, "accepted compute credit_offer="+message.Fields["credit_offer"]); err != nil {
		return err
	}
	if err := runtime.recordEconomics(message.To, "economics_capacity_reserved", "kept", message.From, message.ProtocolName, fmt.Sprintf("reserved compute slot capacity=%d", computeCapacity)); err != nil {
		return err
	}
	if err := runtime.recordEconomics(message.To, "economics_credits_earned", "kept", message.From, message.ProtocolName, fmt.Sprintf("Carol earned compute credits=%d capacity=%d", creditOffer, computeCapacity)); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           "ellen",
		ProtocolName: CIDComputeV1,
		Fields: map[string]string{
			"variant":       "context_request",
			"promise_about": "provide_compute_context",
			"request_id":    requestID,
			"context_cid":   message.Fields["context_cid"],
		},
	})
}

func (runtime *TCPRuntime) handleContextRequest(ctx context.Context, message RuntimeMessage) error {
	contextBytes := sampleContextBytes()
	contextCID := ContentCID(contextBytes)
	if contextCID != message.Fields["context_cid"] {
		return runtime.recordAgentEvent(message.To, "compute_context_not_promised", "non_commitment", message.From, message.ProtocolName, "requested context_cid not available")
	}
	if err := runtime.recordAgentEvent(message.To, "compute_context_promised", "kept", message.From, message.ProtocolName, "Ellen promises explicit context object context_cid="+contextCID); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           message.From,
		ProtocolName: CIDComputeV1,
		Fields: map[string]string{
			"variant":       "context_response",
			"promise_about": "provide_compute_context",
			"request_id":    message.Fields["request_id"],
			"context_cid":   contextCID,
			"context_b64":   base64.StdEncoding.EncodeToString(contextBytes),
		},
	})
}

func (runtime *TCPRuntime) handleContextResponse(ctx context.Context, message RuntimeMessage) error {
	contextBytes, decodeErr := base64.StdEncoding.DecodeString(message.Fields["context_b64"])
	if decodeErr != nil {
		return decodeErr
	}
	contextCID := message.Fields["context_cid"]
	if !VerifyContentCID(contextBytes, contextCID) {
		runtime.adjustTrust(message.To, message.From, -3)
		return runtime.recordAgentEvent(message.To, "compute_context_rejected", "malformed", message.From, message.ProtocolName, "context bytes did not match context_cid="+contextCID)
	}
	state := runtime.states[message.To]
	state.mu.Lock()
	pending := state.PendingCompute[message.Fields["request_id"]]
	delete(state.PendingCompute, message.Fields["request_id"])
	state.mu.Unlock()
	if pending == nil {
		return runtime.recordAgentEvent(message.To, "compute_request_missing", "non_commitment", message.From, message.ProtocolName, "no pending compute request_id="+message.Fields["request_id"])
	}
	functionBytes, functionErr := base64.StdEncoding.DecodeString(pending["function_b64"])
	if functionErr != nil {
		return functionErr
	}
	inputBytes, inputErr := base64.StdEncoding.DecodeString(pending["input_b64"])
	if inputErr != nil {
		return inputErr
	}
	resultBytes, executeErr := ExecuteFunction(functionBytes, inputBytes, contextBytes)
	if executeErr != nil {
		return executeErr
	}
	resultCID := ContentCID(resultBytes)
	if err := runtime.recordAgentEvent(message.To, "compute_function_executed", "kept", pending["sender"], message.ProtocolName, "function_cid="+pending["function_cid"]+" input_cid="+pending["input_cid"]+" result_cid="+resultCID); err != nil {
		return err
	}
	if err := runtime.recordAgentEvent(message.To, "cid_compute_promised", "kept", pending["sender"], message.ProtocolName, "Carol promises compute only for stated function/input/context CIDs"); err != nil {
		return err
	}
	if err := runtime.recordAgentEvent(message.To, "compute_result_promised", "kept", pending["sender"], message.ProtocolName, "result_cid="+resultCID); err != nil {
		return err
	}
	resultFields := map[string]string{
		"variant":       "compute_result",
		"promise_about": "execute_function",
		"function_cid":  pending["function_cid"],
		"function_b64":  pending["function_b64"],
		"input_cid":     pending["input_cid"],
		"input_b64":     pending["input_b64"],
		"context_cid":   contextCID,
		"context_b64":   base64.StdEncoding.EncodeToString(contextBytes),
		"result_cid":    resultCID,
		"result_b64":    base64.StdEncoding.EncodeToString(resultBytes),
	}
	if err := runtime.sendPromise(ctx, OutboundPromise{From: message.To, To: pending["sender"], ProtocolName: CIDComputeV1, Fields: resultFields}); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{From: message.To, To: "dave", ProtocolName: CIDComputeV1, Fields: resultFields})
}

func (runtime *TCPRuntime) handleComputeResult(ctx context.Context, message RuntimeMessage) error {
	resultBytes, decodeErr := base64.StdEncoding.DecodeString(message.Fields["result_b64"])
	if decodeErr != nil {
		return decodeErr
	}
	resultCID := message.Fields["result_cid"]
	if !VerifyContentCID(resultBytes, resultCID) {
		runtime.adjustTrust(message.To, message.From, -3)
		return runtime.recordAgentEvent(message.To, "compute_result_rejected", "malformed", message.From, message.ProtocolName, "result bytes did not match result_cid="+resultCID)
	}
	runtime.adjustTrust(message.To, message.From, 1)
	if runtime.states[message.To].Agent.Role == "cache_peer" {
		cacheKey := ComputeCacheKey(message.ProtocolName, message.Fields["function_cid"], message.Fields["input_cid"], message.Fields["context_cid"], resultCID)
		return runtime.recordAgentEvent(message.To, "compute_cache_checkpointed", "kept", message.From, message.ProtocolName, "Dave caches exact tuple cache_key="+cacheKey)
	}
	if err := runtime.recordAgentEvent(message.To, "compute_result_received", "kept", message.From, message.ProtocolName, "result_cid="+resultCID+" result="+string(resultBytes)); err != nil {
		return err
	}
	if runtime.states[message.To].Agent.Role != "data_holder" {
		return nil
	}
	if err := runtime.recordEconomics(message.To, "economics_credits_spent", "kept", message.From, message.ProtocolName, "Alice spends compute credits=5 after Carol's result arrives"); err != nil {
		return err
	}
	if err := runtime.verifyComputeResultLocally(message); err != nil {
		return err
	}
	for _, verifier := range []string{"dave", "grace"} {
		if err := runtime.sendPromise(ctx, OutboundPromise{
			From:         message.To,
			To:           verifier,
			ProtocolName: CIDComputeV1,
			Fields: map[string]string{
				"variant":       "compute_verification_request",
				"promise_about": "verify_compute_result",
				"function_cid":  message.Fields["function_cid"],
				"function_b64":  message.Fields["function_b64"],
				"input_cid":     message.Fields["input_cid"],
				"input_b64":     message.Fields["input_b64"],
				"context_cid":   message.Fields["context_cid"],
				"context_b64":   message.Fields["context_b64"],
				"result_cid":    message.Fields["result_cid"],
				"result_b64":    message.Fields["result_b64"],
			},
		}); err != nil {
			return err
		}
	}
	return nil
}

func (runtime *TCPRuntime) handleCorruptBytesOffer(message RuntimeMessage) error {
	contentBytes, decodeErr := base64.StdEncoding.DecodeString(message.Fields["content_b64"])
	if decodeErr != nil {
		return decodeErr
	}
	contentCID := message.Fields["content_cid"]
	if VerifyContentCID(contentBytes, contentCID) {
		return runtime.recordAgentEvent(message.To, "cas_verification_promised", "kept", message.From, message.ProtocolName, "Grace verified claimed bytes locally")
	}
	runtime.adjustTrust(message.To, message.From, -4)
	if err := runtime.recordAgentEvent(message.To, "cas_verification_promised", "kept", message.From, message.ProtocolName, "Grace promises to verify bytes against claimed CIDs from her local vantage"); err != nil {
		return err
	}
	if err := runtime.recordAgentEvent(message.To, "cas_corrupt_bytes_rejected", "malformed", message.From, message.ProtocolName, "presented bytes did not match content_cid="+contentCID); err != nil {
		return err
	}
	return runtime.recordAgentEvent(message.To, "cas_corrupt_evidence_recorded", "kept", message.From, message.ProtocolName, "Grace records corrupt-byte evidence locally")
}

func (runtime *TCPRuntime) handleComputeVerificationRequest(ctx context.Context, message RuntimeMessage) error {
	if err := runtime.verifyComputeResultLocally(message); err != nil {
		runtime.adjustTrust(message.To, message.From, -2)
		return runtime.recordAgentEvent(message.To, "compute_verification_rejected", "malformed", message.From, message.ProtocolName, err.Error())
	}
	if err := runtime.recordAgentEvent(message.To, "compute_result_peer_verified", "kept", message.From, message.ProtocolName, "verified result_cid="+message.Fields["result_cid"]); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           message.From,
		ProtocolName: CIDComputeV1,
		Fields: map[string]string{
			"variant":       "compute_verification_result",
			"promise_about": "verify_compute_result",
			"result_cid":    message.Fields["result_cid"],
			"verdict":       "kept",
		},
	})
}

func (runtime *TCPRuntime) handleComputeVerificationResult(message RuntimeMessage) error {
	runtime.adjustTrust(message.To, message.From, 1)
	return runtime.recordAgentEvent(message.To, "compute_verification_received", "kept", message.From, message.ProtocolName, "peer verified result_cid="+message.Fields["result_cid"])
}

func (runtime *TCPRuntime) handleTrustRepairPromise(message RuntimeMessage) error {
	return runtime.recordAgentEvent(message.To, "trust_repair_promise_recorded", "kept", message.From, message.ProtocolName, message.Fields["decision_text"])
}

// sendPromise signs, validates, and sends one raw grid envelope over TCP.
func (runtime *TCPRuntime) sendPromise(ctx context.Context, outbound OutboundPromise) error {
	fields := make(map[string]string, len(outbound.Fields)+3)
	for key, value := range outbound.Fields {
		fields[key] = value
	}
	fields["act"] = "promise"
	fields["sender"] = outbound.From
	fields["recipient"] = outbound.To
	protocolCID := runtime.Registry.MustCID(outbound.ProtocolName)
	envelope, envelopeErr := NewEnvelope(protocolCID, fields, outbound.From)
	if envelopeErr != nil {
		return envelopeErr
	}
	exactBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	parsed, parseErr := ParseEnvelope(exactBytes)
	if parseErr != nil {
		return parseErr
	}
	if verifyErr := VerifyEnvelope(parsed); verifyErr != nil {
		return verifyErr
	}
	if err := runtime.recordAgentEvent(outbound.From, "promise_envelope_validated", "kept", outbound.To, outbound.ProtocolName, "exact_sha256="+HashExactBytes(exactBytes)+" promise_about="+fields["promise_about"]); err != nil {
		return err
	}
	conn, dialErr := runtime.dialWithRetry(ctx, outbound.To)
	if dialErr != nil {
		recordErr := runtime.recordAgentEvent(outbound.From, "tcp_message_send_failed", "non_commitment", outbound.To, outbound.ProtocolName, dialErr.Error())
		if recordErr != nil {
			return recordErr
		}
		return dialErr
	}
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "poc13-runtime: close outbound connection: %v\n", closeErr)
		}
	}()
	if err := (FrameWriter{writer: conn}).WriteFrame(exactBytes); err != nil {
		return err
	}
	runtime.markActivity()
	return runtime.recordAgentEvent(outbound.From, "tcp_message_sent", "kept", outbound.To, outbound.ProtocolName, "variant="+fields["variant"]+" exact_sha256="+HashExactBytes(exactBytes))
}

func (runtime *TCPRuntime) dialWithRetry(ctx context.Context, agentName string) (net.Conn, error) {
	address, addressErr := runtime.dialAddressForAgent(agentName)
	if addressErr != nil {
		return nil, addressErr
	}
	var lastErr error
	for attempt := 0; attempt < 20; attempt++ {
		dialer := net.Dialer{Timeout: 500 * time.Millisecond}
		conn, dialErr := dialer.DialContext(ctx, "tcp", address)
		if dialErr == nil {
			return conn, nil
		}
		lastErr = dialErr
		time.Sleep(100 * time.Millisecond)
	}
	return nil, lastErr
}

func (runtime *TCPRuntime) recordRuntimeReadiness() error {
	if err := os.MkdirAll(runtime.runtimeMarkerDir(), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(runtime.runtimeMarkerPath("ready"), []byte(time.Now().UTC().Format(time.RFC3339Nano)+"\n"), 0o644); err != nil {
		return err
	}
	for _, agentName := range runtime.Container.Agents {
		if err := runtime.recordAgentEvent(agentName, "runtime_readiness_promised", "kept", "", "", "container "+runtime.Container.Name+" promises its TCP listener is open"); err != nil {
			return err
		}
	}
	return nil
}

func (runtime *TCPRuntime) waitForPeerReadiness(ctx context.Context) error {
	deadline := time.Now().Add(runtime.Config.ReadinessTimeout())
	observed := make(map[string]bool)
	for {
		allReady := true
		for _, container := range runtime.Config.Containers {
			readyPath := filepath.Join(runtime.runtimeMarkerDir(), container.Name+".ready")
			if _, statErr := os.Stat(readyPath); statErr != nil {
				allReady = false
				continue
			}
			if !observed[container.Name] {
				observed[container.Name] = true
				for _, agentName := range runtime.Container.Agents {
					if err := runtime.recordAgentEvent(agentName, "peer_readiness_observed", "kept", "", "", "observed container "+container.Name+" readiness marker"); err != nil {
						return err
					}
				}
			}
		}
		if allReady {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for peer readiness after %s", runtime.Config.ReadinessTimeout())
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(50 * time.Millisecond):
		}
	}
}

func (runtime *TCPRuntime) waitForRuntimeCompletion(ctx context.Context) error {
	idle := runtime.Config.CompletionIdle()
	deadline := time.Now().Add(runtime.Config.Timeout())
	runtime.markActivity()
	for {
		runtime.mu.Lock()
		activeHandlers := runtime.activeHandlers
		quietFor := time.Since(runtime.lastActivity)
		runtime.mu.Unlock()
		if activeHandlers == 0 && quietFor >= idle {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out waiting for runtime completion: active_handlers=%d quiet_for=%s idle=%s", activeHandlers, quietFor, idle)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(25 * time.Millisecond):
		}
	}
}

func (runtime *TCPRuntime) recordRuntimeDone() error {
	if err := os.MkdirAll(runtime.runtimeMarkerDir(), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(runtime.runtimeMarkerPath("done"), []byte(time.Now().UTC().Format(time.RFC3339Nano)+"\n"), 0o644); err != nil {
		return err
	}
	for _, agentName := range runtime.Container.Agents {
		if err := runtime.recordAgentEvent(agentName, "runtime_done_promised", "kept", "", "", "container "+runtime.Container.Name+" promises local runtime is idle and done"); err != nil {
			return err
		}
	}
	return nil
}

func (runtime *TCPRuntime) runtimeMarkerDir() string {
	return filepath.Join(runtime.Config.RunRoot, runtime.Config.RunID, "runtime")
}

func (runtime *TCPRuntime) runtimeMarkerPath(kind string) string {
	return filepath.Join(runtime.runtimeMarkerDir(), runtime.Container.Name+"."+kind)
}

func (runtime *TCPRuntime) markActivity() {
	runtime.mu.Lock()
	runtime.lastActivity = time.Now()
	runtime.mu.Unlock()
}

func (runtime *TCPRuntime) recordAgentEvent(agentName, eventName, outcome, peer, protocolName, detail string) error {
	logFile, openErr := openAgentLog(runtime.Config, agentName)
	if openErr != nil {
		return openErr
	}
	defer func() {
		closeErr := logFile.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "poc13-runtime: close %s log: %v\n", agentName, closeErr)
		}
	}()
	recorder := &Recorder{log: logFile, agent: AgentConfig{Name: agentName}, registry: runtime.Registry}
	recorder.Record(eventName, outcome, peer, protocolName, detail)
	return nil
}

func (runtime *TCPRuntime) recordTrustDrivenChoice(agentName, peer, protocolName, detail string) error {
	return runtime.recordAgentEvent(agentName, "trust_driven_peer_choice", "kept", peer, protocolName, detail)
}

// recordEconomics writes one local economics evidence record.
// Intent: POC13 economics remain local promise evidence about capacity,
// opportunity cost, and credits; there is no central price authority. Source:
// DI-lupag
func (runtime *TCPRuntime) recordEconomics(agentName, eventName, outcome, peer, protocolName, detail string) error {
	return runtime.recordAgentEvent(agentName, eventName, outcome, peer, protocolName, detail)
}

// verifyComputeResultLocally recomputes the signed payload material from the
// observer's local vantage.
// Intent: Alice, Dave, and Grace verify Carol's compute result using the same
// pCID-owned payload bytes instead of trusting Carol as an authority. Source:
// DI-lupag
func (runtime *TCPRuntime) verifyComputeResultLocally(message RuntimeMessage) error {
	functionBytes, functionErr := base64.StdEncoding.DecodeString(message.Fields["function_b64"])
	if functionErr != nil {
		return functionErr
	}
	inputBytes, inputErr := base64.StdEncoding.DecodeString(message.Fields["input_b64"])
	if inputErr != nil {
		return inputErr
	}
	contextBytes, contextErr := base64.StdEncoding.DecodeString(message.Fields["context_b64"])
	if contextErr != nil {
		return contextErr
	}
	expectedBytes, executeErr := ExecuteFunction(functionBytes, inputBytes, contextBytes)
	if executeErr != nil {
		return executeErr
	}
	expectedCID := ContentCID(expectedBytes)
	if expectedCID != message.Fields["result_cid"] {
		return fmt.Errorf("local recompute result_cid=%s want %s", expectedCID, message.Fields["result_cid"])
	}
	return runtime.recordAgentEvent(message.To, "compute_result_locally_verified", "kept", message.From, message.ProtocolName, "local recompute matched result_cid="+expectedCID)
}

func (runtime *TCPRuntime) preferredStoragePeer(agentName string) string {
	if runtime.trustScore(agentName, "bob") >= runtime.trustScore(agentName, "frank") {
		return "bob"
	}
	return "frank"
}

func (runtime *TCPRuntime) preferredComputePeer(agentName string) string {
	if runtime.trustScore(agentName, "carol") >= 0 {
		return "carol"
	}
	return "dave"
}

func (runtime *TCPRuntime) adjustTrust(observer, peer string, delta int) {
	state := runtime.states[observer]
	if state == nil {
		return
	}
	state.mu.Lock()
	state.Trust[peer] += delta
	currentTrust := state.Trust[peer]
	state.mu.Unlock()
	outcome := "kept"
	if delta < 0 {
		outcome = "broken"
	}
	if err := runtime.recordAgentEvent(observer, "trust_updated", outcome, peer, "", fmt.Sprintf("delta=%d local_trust=%d", delta, currentTrust)); err != nil {
		fmt.Fprintf(os.Stderr, "poc13-runtime: record trust update: %v\n", err)
	}
}

func (runtime *TCPRuntime) trustScore(observer, peer string) int {
	state := runtime.states[observer]
	if state == nil {
		return 0
	}
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.Trust[peer]
}

func parseCreditOffer(value string) int {
	var creditOffer int
	if _, err := fmt.Sscanf(value, "%d", &creditOffer); err != nil {
		return 0
	}
	return creditOffer
}

func (runtime *TCPRuntime) issueCapabilityToken(issuer, issuee, contentCID string) string {
	token := ContentCID([]byte(issuer + "|" + issuee + "|" + contentCID + "|serve-once"))
	state := runtime.states[issuer]
	state.mu.Lock()
	state.Capabilities[issuee+"|"+contentCID] = token
	state.mu.Unlock()
	return token
}

func (runtime *TCPRuntime) redeemCapabilityToken(issuer, issuee, contentCID, token string) bool {
	state := runtime.states[issuer]
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.Capabilities[issuee+"|"+contentCID] == token
}

func (runtime *TCPRuntime) containerForAgent(agentName string) (ContainerConfig, bool) {
	for _, container := range runtime.Config.Containers {
		for _, localAgentName := range container.Agents {
			if localAgentName == agentName {
				return container, true
			}
		}
	}
	return ContainerConfig{}, false
}

func (runtime *TCPRuntime) dialAddressForAgent(agentName string) (string, error) {
	container, ok := runtime.containerForAgent(agentName)
	if !ok {
		return "", fmt.Errorf("unknown target agent %s", agentName)
	}
	if container.Name == runtime.Container.Name {
		return fmt.Sprintf("127.0.0.1:%d", runtime.listenPort), nil
	}
	if runtime.Config.ListenPort == 0 {
		return "", fmt.Errorf("remote listen_port cannot be zero for target %s", agentName)
	}
	return fmt.Sprintf("%s:%d", container.Name, runtime.Config.ListenPort), nil
}

func (runtime *TCPRuntime) listenAddress() string {
	return fmt.Sprintf(":%d", runtime.Config.ListenPort)
}

// ReadFrame reads one bounded length-prefixed envelope.
func (frameReader FrameReader) ReadFrame() ([]byte, error) {
	var lengthBytes [4]byte
	if _, err := io.ReadFull(frameReader.reader, lengthBytes[:]); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBytes[:])
	if length == 0 || length > maxFrameBytes {
		return nil, fmt.Errorf("invalid frame length %d", length)
	}
	frameBytes := make([]byte, length)
	if _, err := io.ReadFull(frameReader.reader, frameBytes); err != nil {
		return nil, err
	}
	return frameBytes, nil
}

// WriteFrame writes one bounded length-prefixed envelope.
func (frameWriter FrameWriter) WriteFrame(frameBytes []byte) error {
	if len(frameBytes) == 0 || len(frameBytes) > maxFrameBytes {
		return fmt.Errorf("invalid frame length %d", len(frameBytes))
	}
	var lengthBytes [4]byte
	binary.BigEndian.PutUint32(lengthBytes[:], uint32(len(frameBytes)))
	if _, err := frameWriter.writer.Write(lengthBytes[:]); err != nil {
		return err
	}
	_, err := frameWriter.writer.Write(frameBytes)
	return err
}

// ExecuteFunction runs the bounded POC13 compute language over payload-provided
// function, input, and context bytes.
// Intent: POC13 should prove CID-named compute over dynamic payload material
// without executing arbitrary host code or hiding ambient inputs. Source:
// DI-fumol
func ExecuteFunction(functionBytes, inputBytes, contextBytes []byte) ([]byte, error) {
	functionText := strings.TrimSpace(string(functionBytes))
	inputText := strings.TrimSpace(string(inputBytes))
	if !strings.Contains(functionText, "fibonacci") {
		return nil, fmt.Errorf("unsupported function source %q", functionText)
	}
	if !strings.HasPrefix(inputText, "n=") {
		return nil, fmt.Errorf("unsupported input %q", inputText)
	}
	var n int
	if _, err := fmt.Sscanf(inputText, "n=%d", &n); err != nil {
		return nil, err
	}
	if n < 0 || n > 40 {
		return nil, fmt.Errorf("fibonacci input out of bounded POC range: %d", n)
	}
	value := fibonacci(n)
	return []byte(fmt.Sprintf("fibonacci(%d)=%d;context_cid=%s", n, value, ContentCID(contextBytes))), nil
}

func fibonacci(n int) int {
	if n < 2 {
		return n
	}
	previous, current := 0, 1
	for index := 2; index <= n; index++ {
		previous, current = current, previous+current
	}
	return current
}
