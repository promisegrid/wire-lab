package poc13

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const maxFrameBytes = 1 << 20

const forceTCPUnreachableAddress = "127.0.0.1:9"

const (
	OutageVariantContainerStopped = "container_stopped"
	OutageVariantNetworkPartition = "network_partition"
)

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
	ComputeCache   map[string]map[string]string
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
			ComputeCache:   make(map[string]map[string]string),
		}
	}
	if err := runtime.loadPersistedTrustHistory(); err != nil {
		return nil, err
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
	if err := runtime.recordAgentEvent(agent.Name, "app_receive_promise_registered", "kept", "", "", "cas_storage_v1, cid_compute_v1, and evidence_report_v1 payloads are accepted only as local promises"); err != nil {
		return err
	}
	if len(state.Trust) > 0 {
		if err := runtime.recordAgentEvent(agent.Name, "persisted_trust_history_loaded", "kept", "", "", "loaded container-local trust history from "+runtime.persistedTrustHistoryPath()); err != nil {
			return err
		}
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
	secondContentBytes := sampleSecondContentBytes()
	secondContentCID := ContentCID(secondContentBytes)
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
	if err := runtime.recordDynamicPeerChoice(agent.Name, storagePeer, CASStorageV1, "storage choice used persisted local trust history before any live sends"); err != nil {
		return err
	}
	if err := runtime.recordTrustDrivenChoice(agent.Name, computePeer, CIDComputeV1, "Alice chooses Carol for compute and plans Dave/Grace verification from prior evidence"); err != nil {
		return err
	}
	if err := runtime.recordDynamicPeerChoice(agent.Name, computePeer, CIDComputeV1, "compute choice used persisted local trust history before any live sends"); err != nil {
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
	if err := runtime.recordAgentEvent(agent.Name, "cas_multi_object_pressure", "kept", storagePeer, CASStorageV1, "Alice sends a second independent object to test multi-object storage pressure"); err != nil {
		return err
	}
	if err := runtime.sendPromise(ctx, OutboundPromise{
		From:         agent.Name,
		To:           storagePeer,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":           "store_request",
			"promise_about":     "store_content",
			"content_cid":       secondContentCID,
			"content_b64":       base64.StdEncoding.EncodeToString(secondContentBytes),
			"credit_offer":      "4",
			"requested_token":   "serve-once",
			"decision_text":     decisionText,
			"promise_condition": "bob may decline locally",
			"object_label":      "second-object",
		},
	}); err != nil {
		return err
	}
	if err := runtime.recordEconomics(agent.Name, "economics_credit_offered", "kept", computePeer, CIDComputeV1, "Alice offers compute credit_offer=5 for payload-provided Fibonacci compute"); err != nil {
		return err
	}
	if err := runtime.sendPromise(ctx, OutboundPromise{
		From:         agent.Name,
		To:           "dave",
		ProtocolName: CIDComputeV1,
		Fields: map[string]string{
			"variant":       "compute_cache_lookup",
			"promise_about": "lookup_compute_cache",
			"function_cid":  ContentCID(functionBytes),
			"input_cid":     ContentCID(inputBytes),
			"context_cid":   ContentCID(contextBytes),
		},
	}); err != nil {
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
	if err := runtime.sendPromise(ctx, OutboundPromise{
		From:         agent.Name,
		To:           "grace",
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "trust_repair_promise",
			"promise_about": "label_future_malformed_evidence",
			"decision_text": "Mallory promises to label future malformed evidence explicitly; Grace remains free to distrust this until future evidence is kept.",
		},
	}); err != nil {
		return err
	}
	if err := runtime.sendUnknownProtocolPromise(ctx, agent.Name, "grace", map[string]string{
		"variant":       "mystery_storage_claim",
		"promise_about": "unknown_protocol_probe",
		"decision_text": decisionText,
	}); err != nil {
		return err
	}
	if err := runtime.sendPromise(ctx, OutboundPromise{
		From:         agent.Name,
		To:           "grace",
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "unsupported_storage_variant",
			"promise_about": "unsupported_variant_probe",
			"decision_text": decisionText,
		},
	}); err != nil {
		return err
	}
	if err := runtime.sendBadProofPromise(ctx, agent.Name, "grace", CASStorageV1, map[string]string{
		"variant":       "bad_proof_probe",
		"promise_about": "present_bad_proof_evidence",
		"decision_text": decisionText,
	}); err != nil {
		return err
	}
	if err := runtime.sendPromise(ctx, OutboundPromise{
		From:         agent.Name,
		To:           "grace",
		ProtocolName: EvidenceReportV1,
		Fields: map[string]string{
			"variant":        "key_rotation_promise",
			"promise_about":  "rotate_signing_key",
			"new_key_label":  "mallory-next-key",
			"rotation_scope": "future-poc13-evidence",
			"decision_text":  decisionText,
		},
	}); err != nil {
		return err
	}
	functionBytes := sampleFunctionBytes()
	inputBytes := []byte("n=7")
	contextBytes := sampleContextBytes()
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         agent.Name,
		To:           "carol",
		ProtocolName: CIDComputeV1,
		Fields: map[string]string{
			"variant":        "compute_request",
			"promise_about":  "execute_function",
			"function_cid":   ContentCID(functionBytes),
			"function_b64":   base64.StdEncoding.EncodeToString(functionBytes),
			"input_cid":      ContentCID(inputBytes),
			"input_b64":      base64.StdEncoding.EncodeToString(inputBytes),
			"context_cid":    ContentCID(contextBytes),
			"credit_offer":   "5",
			"capacity_probe": "true",
			"decision_text":  decisionText,
		},
	})
}

func (runtime *TCPRuntime) handleEnvelope(ctx context.Context, exactBytes []byte) error {
	envelope, parseErr := ParseEnvelope(exactBytes)
	if parseErr != nil {
		return parseErr
	}
	fields, fieldsErr := envelope.PayloadFields()
	if fieldsErr != nil {
		return fieldsErr
	}
	protocolName, known := runtime.Registry.Name(envelope.ProtocolCID)
	if !known {
		protocolName = envelope.ProtocolCID.String()
	}
	if verifyErr := VerifyEnvelope(envelope); verifyErr != nil {
		return runtime.handleBadProofEnvelope(fields, protocolName, verifyErr)
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
	if !known {
		return runtime.handleUnknownProtocolEnvelope(message)
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
	case "token_revocation_notice":
		return runtime.handleTokenRevocationNotice(ctx, message)
	case "replica_token_renewal_request":
		return runtime.handleReplicaTokenRenewalRequest(ctx, message)
	case "replica_token_renewal":
		return runtime.handleReplicaTokenRenewal(message)
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
	case "compute_cache_lookup":
		return runtime.handleComputeCacheLookup(ctx, message)
	case "compute_cache_response":
		return runtime.handleComputeCacheResponse(message)
	case "compute_verification_request":
		return runtime.handleComputeVerificationRequest(ctx, message)
	case "compute_verification_result":
		return runtime.handleComputeVerificationResult(message)
	case "evidence_report":
		return runtime.handleEvidenceReport(message)
	case "key_rotation_promise":
		return runtime.handleKeyRotationPromise(message)
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
	for _, outageVariant := range []string{OutageVariantContainerStopped, OutageVariantNetworkPartition} {
		if err := runtime.recordAgentEvent(message.To, "network_outage_variant_selected", "kept", message.From, message.ProtocolName, "Alice tests Bob retrieval under outage_variant="+outageVariant); err != nil {
			return err
		}
		sendErr := runtime.sendPromise(ctx, OutboundPromise{
			From:         message.To,
			To:           message.From,
			ProtocolName: CASStorageV1,
			Fields: map[string]string{
				"variant":        "serve_request",
				"promise_about":  "redeem_storage_capability",
				"content_cid":    message.Fields["content_cid"],
				"token":          message.Fields["token"],
				"outage_variant": outageVariant,
			},
		})
		if sendErr == nil {
			return nil
		}
		if err := runtime.recordAgentEvent(message.To, "primary_storage_unavailable", "non_commitment", message.From, message.ProtocolName, "Alice could not open a TCP path to Bob for outage_variant="+outageVariant+": "+sendErr.Error()); err != nil {
			return err
		}
	}
	return nil
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
	replicaToken := message.Fields["replica_token"]
	if runtime.trustScore(message.To, replicaPeer) < 0 {
		return runtime.recordAgentEvent(message.To, "replica_recovery_not_promised", "non_commitment", replicaPeer, message.ProtocolName, "local trust below replica recovery threshold")
	}
	state := runtime.states[message.To]
	state.mu.Lock()
	state.Capabilities[replicaCapabilityKey(replicaPeer, contentCID)] = replicaToken
	state.mu.Unlock()
	if err := runtime.recordAgentEvent(message.To, "replica_capability_token_received", "kept", replicaPeer, message.ProtocolName, "Frank-issued replica token received content_cid="+contentCID); err != nil {
		return err
	}
	if err := runtime.recordTokenLifecycle(message.To, "capability_token_ttl_observed", "kept", replicaPeer, message.ProtocolName, "Alice observes replica token ttl_uses="+message.Fields["ttl_uses"]+" content_cid="+contentCID); err != nil {
		return err
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
			"replica_token": replicaToken,
		},
	})
}

func (runtime *TCPRuntime) handleReplicaServeRequest(ctx context.Context, message RuntimeMessage) error {
	contentCID := message.Fields["content_cid"]
	if !runtime.redeemCapabilityToken(message.To, message.From, contentCID, message.Fields["replica_token"]) {
		runtime.adjustTrust(message.To, message.From, -1)
		return runtime.recordAgentEvent(message.To, "replica_capability_token_rejected", "malformed", message.From, message.ProtocolName, "replica token did not match Frank-local promise for content_cid="+contentCID)
	}
	state := runtime.states[message.To]
	state.mu.Lock()
	contentBytes, ok := state.Store[contentCID]
	state.mu.Unlock()
	if !ok {
		return runtime.recordAgentEvent(message.To, "cas_replica_serve_not_promised", "non_commitment", message.From, message.ProtocolName, "replica content not present content_cid="+contentCID)
	}
	if err := runtime.recordAgentEvent(message.To, "replica_capability_token_redeemed", "kept", message.From, message.ProtocolName, "Frank redeems replica token content_cid="+contentCID); err != nil {
		return err
	}
	if err := runtime.recordAgentEvent(message.To, "cas_replica_serve_promised", "kept", message.From, message.ProtocolName, "Frank promises replica serving for content_cid="+contentCID); err != nil {
		return err
	}
	if err := runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           message.From,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "serve_response",
			"promise_about": "serve_content_bytes",
			"content_cid":   contentCID,
			"content_b64":   base64.StdEncoding.EncodeToString(contentBytes),
		},
	}); err != nil {
		return err
	}
	if err := runtime.recordTokenLifecycle(message.To, "capability_token_expired", "non_commitment", message.From, message.ProtocolName, "replica token ttl exhausted after one serve content_cid="+contentCID); err != nil {
		return err
	}
	if err := runtime.recordTokenLifecycle(message.To, "capability_token_revoked", "non_commitment", message.From, message.ProtocolName, "Frank revokes spent replica token content_cid="+contentCID); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           message.From,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "token_revocation_notice",
			"promise_about": "replica_token_lifecycle",
			"content_cid":   contentCID,
			"replica_token": message.Fields["replica_token"],
			"reason":        "ttl_exhausted_after_one_serve",
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

func (runtime *TCPRuntime) handleTokenRevocationNotice(ctx context.Context, message RuntimeMessage) error {
	contentCID := message.Fields["content_cid"]
	if err := runtime.recordTokenLifecycle(message.To, "capability_token_revoked", "non_commitment", message.From, message.ProtocolName, "Alice records Frank's revocation reason="+message.Fields["reason"]+" content_cid="+contentCID); err != nil {
		return err
	}
	if err := runtime.recordTokenLifecycle(message.To, "capability_token_renewal_requested", "kept", message.From, message.ProtocolName, "Alice asks Frank for renewed replica token content_cid="+contentCID); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           message.From,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "replica_token_renewal_request",
			"promise_about": "replica_token_lifecycle",
			"content_cid":   contentCID,
			"reason":        "future-retrieval-option",
		},
	})
}

func (runtime *TCPRuntime) handleReplicaTokenRenewalRequest(ctx context.Context, message RuntimeMessage) error {
	contentCID := message.Fields["content_cid"]
	state := runtime.states[message.To]
	state.mu.Lock()
	_, ok := state.Store[contentCID]
	state.mu.Unlock()
	if !ok {
		return runtime.recordTokenLifecycle(message.To, "capability_token_renewal_refused", "non_commitment", message.From, message.ProtocolName, "Frank has no replica bytes for renewal content_cid="+contentCID)
	}
	token := runtime.issueCapabilityToken(message.To, message.From, contentCID+"|renewed")
	if err := runtime.recordTokenLifecycle(message.To, "capability_token_renewed", "kept", message.From, message.ProtocolName, "Frank renews replica token content_cid="+contentCID); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           message.From,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "replica_token_renewal",
			"promise_about": "replica_token_lifecycle",
			"content_cid":   contentCID,
			"replica_token": token,
			"ttl_uses":      "1",
		},
	})
}

func (runtime *TCPRuntime) handleReplicaTokenRenewal(message RuntimeMessage) error {
	state := runtime.states[message.To]
	state.mu.Lock()
	state.Capabilities[replicaCapabilityKey(message.From, message.Fields["content_cid"])] = message.Fields["replica_token"]
	state.mu.Unlock()
	return runtime.recordTokenLifecycle(message.To, "capability_token_renewed", "kept", message.From, message.ProtocolName, "Alice records renewed replica token content_cid="+message.Fields["content_cid"])
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
	requester := message.Fields["requester"]
	replicaToken := runtime.issueCapabilityToken(message.To, requester, contentCID)
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
	if requester != "" {
		if err := runtime.recordAgentEvent(message.To, "replica_capability_token_issued", "kept", requester, message.ProtocolName, "Frank issues replica token content_cid="+contentCID); err != nil {
			return err
		}
		if err := runtime.recordTokenLifecycle(message.To, "capability_token_ttl_promised", "kept", requester, message.ProtocolName, "replica token ttl_uses=1 content_cid="+contentCID); err != nil {
			return err
		}
	}
	if err := runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           message.From,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "replicate_acceptance",
			"promise_about": "replicate_content",
			"content_cid":   contentCID,
			"requester":     message.Fields["requester"],
		},
	}); err != nil {
		return err
	}
	if requester == "" {
		return nil
	}
	return runtime.sendReplicaAvailable(ctx, message.To, requester, contentCID, replicaToken)
}

func (runtime *TCPRuntime) handleReplicateAcceptance(ctx context.Context, message RuntimeMessage) error {
	runtime.adjustTrust(message.To, message.From, 1)
	return runtime.recordAgentEvent(message.To, "cas_replication_confirmed", "kept", message.From, message.ProtocolName, "replica accepted content_cid="+message.Fields["content_cid"])
}

func (runtime *TCPRuntime) handleComputeRequest(ctx context.Context, message RuntimeMessage) error {
	state := runtime.states[message.To]
	requestID := message.Fields["function_cid"] + "|" + message.Fields["input_cid"]
	creditOffer := parseCreditOffer(message.Fields["credit_offer"])
	if message.Fields["capacity_probe"] == "true" {
		return runtime.recordCapacityPressure(message.To, message.From, message.ProtocolName, fmt.Sprintf("Carol declines competing compute request from %s because capacity is reserved for prior local promises capacity=%d", message.From, computeCapacity))
	}
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
	if functionKind(pending["function_b64"]) != "fibonacci" {
		if err := runtime.recordAgentEvent(message.To, "compute_alternate_function_executed", "kept", pending["sender"], message.ProtocolName, "function_kind="+functionKind(pending["function_b64"])+" result_cid="+resultCID); err != nil {
			return err
		}
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
	badResultBytes := badComputeResultBytes(resultBytes)
	badResultFields := map[string]string{
		"variant":           "compute_result",
		"promise_about":     "execute_function",
		"function_cid":      pending["function_cid"],
		"function_b64":      pending["function_b64"],
		"input_cid":         pending["input_cid"],
		"input_b64":         pending["input_b64"],
		"context_cid":       contextCID,
		"context_b64":       base64.StdEncoding.EncodeToString(contextBytes),
		"result_cid":        ContentCID(badResultBytes),
		"result_b64":        base64.StdEncoding.EncodeToString(badResultBytes),
		"adversarial_probe": "true",
	}
	if err := runtime.recordAgentEvent(message.To, "compute_bad_result_promised", "malformed", pending["sender"], message.ProtocolName, "Carol sends one hash-valid but semantically wrong result for verifier pressure"); err != nil {
		return err
	}
	if err := runtime.sendPromise(ctx, OutboundPromise{From: message.To, To: pending["sender"], ProtocolName: CIDComputeV1, Fields: badResultFields}); err != nil {
		return err
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
		state := runtime.states[message.To]
		state.mu.Lock()
		state.ComputeCache[cacheKey] = copyFields(message.Fields)
		state.mu.Unlock()
		if err := runtime.recordAgentEvent(message.To, "compute_cache_checkpointed", "kept", message.From, message.ProtocolName, "Dave caches exact tuple cache_key="+cacheKey); err != nil {
			return err
		}
		if err := runtime.recordAgentEvent(message.To, "compute_cache_hit", "kept", "alice", message.ProtocolName, "Dave offers exact cache entry cache_key="+cacheKey); err != nil {
			return err
		}
		responseFields := copyFields(message.Fields)
		responseFields["variant"] = "compute_cache_response"
		responseFields["promise_about"] = "lookup_compute_cache"
		responseFields["cache_key"] = cacheKey
		responseFields["cache_status"] = "hit"
		return runtime.sendPromise(ctx, OutboundPromise{From: message.To, To: "alice", ProtocolName: message.ProtocolName, Fields: responseFields})
	}
	if err := runtime.recordAgentEvent(message.To, "compute_result_received", "kept", message.From, message.ProtocolName, "result_cid="+resultCID+" result="+string(resultBytes)); err != nil {
		return err
	}
	if runtime.states[message.To].Agent.Role != "data_holder" {
		return nil
	}
	verifyErr := runtime.verifyComputeResultLocally(message)
	if verifyErr != nil {
		runtime.adjustTrust(message.To, message.From, -2)
		if err := runtime.recordAgentEvent(message.To, "compute_result_locally_rejected", "malformed", message.From, message.ProtocolName, verifyErr.Error()); err != nil {
			return err
		}
		if err := runtime.recordEconomics(message.To, "economics_payment_withheld", "non_commitment", message.From, message.ProtocolName, "Alice withholds compute credits after local recompute rejected result_cid="+resultCID); err != nil {
			return err
		}
		for _, verifier := range []string{"dave", "grace"} {
			if err := runtime.sendComputeVerificationRequest(ctx, message, verifier); err != nil {
				return err
			}
		}
		return nil
	}
	if err := runtime.recordEconomics(message.To, "economics_credits_spent", "kept", message.From, message.ProtocolName, "Alice spends compute credits=5 after Carol's result verifies locally"); err != nil {
		return err
	}
	if functionKind(message.Fields["function_b64"]) == "fibonacci" {
		state := runtime.states[message.To]
		state.mu.Lock()
		_, followupSent := state.Capabilities["sum-compute-followup-sent"]
		if !followupSent {
			state.Capabilities["sum-compute-followup-sent"] = "true"
		}
		state.mu.Unlock()
		if !followupSent {
			sumFunctionBytes := sampleSumFunctionBytes()
			sumInputBytes := sampleSumInputBytes()
			contextBytes := sampleContextBytes()
			if err := runtime.recordAgentEvent(message.To, "compute_followup_function_requested", "kept", message.From, message.ProtocolName, "Alice requests second payload-provided compute function kind=sum"); err != nil {
				return err
			}
			if err := runtime.sendPromise(ctx, OutboundPromise{
				From:         message.To,
				To:           message.From,
				ProtocolName: CIDComputeV1,
				Fields: map[string]string{
					"variant":           "compute_request",
					"promise_about":     "execute_function",
					"function_cid":      ContentCID(sumFunctionBytes),
					"function_b64":      base64.StdEncoding.EncodeToString(sumFunctionBytes),
					"input_cid":         ContentCID(sumInputBytes),
					"input_b64":         base64.StdEncoding.EncodeToString(sumInputBytes),
					"context_cid":       ContentCID(contextBytes),
					"credit_offer":      "5",
					"decision_text":     "Alice promises to accept a sum result only if it verifies against the payload bytes.",
					"promise_condition": "carol may decline locally",
				},
			}); err != nil {
				return err
			}
		}
	}
	for _, verifier := range []string{"dave", "grace"} {
		if err := runtime.sendComputeVerificationRequest(ctx, message, verifier); err != nil {
			return err
		}
	}
	return nil
}

func (runtime *TCPRuntime) handleComputeCacheLookup(ctx context.Context, message RuntimeMessage) error {
	cacheKey := computeCacheKeyFromFields(message.ProtocolName, message.Fields)
	state := runtime.states[message.To]
	state.mu.Lock()
	cachedFields, ok := state.ComputeCache[cacheKey]
	state.mu.Unlock()
	if !ok {
		if err := runtime.recordAgentEvent(message.To, "compute_cache_miss", "non_commitment", message.From, message.ProtocolName, "Dave has no cache entry cache_key="+cacheKey); err != nil {
			return err
		}
		return runtime.sendPromise(ctx, OutboundPromise{
			From:         message.To,
			To:           message.From,
			ProtocolName: message.ProtocolName,
			Fields: map[string]string{
				"variant":       "compute_cache_response",
				"promise_about": "lookup_compute_cache",
				"cache_key":     cacheKey,
				"cache_status":  "miss",
			},
		})
	}
	if err := runtime.recordAgentEvent(message.To, "compute_cache_hit", "kept", message.From, message.ProtocolName, "Dave has exact cache entry cache_key="+cacheKey); err != nil {
		return err
	}
	responseFields := copyFields(cachedFields)
	responseFields["variant"] = "compute_cache_response"
	responseFields["promise_about"] = "lookup_compute_cache"
	responseFields["cache_key"] = cacheKey
	responseFields["cache_status"] = "hit"
	return runtime.sendPromise(ctx, OutboundPromise{From: message.To, To: message.From, ProtocolName: message.ProtocolName, Fields: responseFields})
}

func (runtime *TCPRuntime) handleComputeCacheResponse(message RuntimeMessage) error {
	if message.Fields["cache_status"] == "miss" {
		return runtime.recordAgentEvent(message.To, "compute_cache_miss_observed", "non_commitment", message.From, message.ProtocolName, "Alice observes cache miss cache_key="+message.Fields["cache_key"])
	}
	if err := runtime.recordAgentEvent(message.To, "compute_cache_reused", "kept", message.From, message.ProtocolName, "Alice reuses cached compute result cache_key="+message.Fields["cache_key"]); err != nil {
		return err
	}
	return runtime.recordAgentEvent(message.To, "compute_result_received", "kept", message.From, message.ProtocolName, "cache result_cid="+message.Fields["result_cid"])
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
		if recordErr := runtime.recordAgentEvent(message.To, "compute_result_peer_rejected", "malformed", message.Fields["result_promiser"], message.ProtocolName, err.Error()); recordErr != nil {
			return recordErr
		}
		return runtime.sendPromise(ctx, OutboundPromise{
			From:         message.To,
			To:           message.From,
			ProtocolName: EvidenceReportV1,
			Fields: map[string]string{
				"variant":            "evidence_report",
				"promise_about":      "local_compute_observation",
				"subject_pcid":       CIDComputeV1,
				"subject_peer":       message.Fields["result_promiser"],
				"subject_result_cid": message.Fields["result_cid"],
				"verdict":            "broken",
				"reason":             err.Error(),
			},
		})
	}
	if message.Fields["disagreement_probe"] == "true" {
		if err := runtime.recordAgentEvent(message.To, "compute_verifier_disagreement", "non_commitment", message.Fields["result_promiser"], message.ProtocolName, "Grace can recompute the result but locally withholds endorsement for disagreement pressure result_cid="+message.Fields["result_cid"]); err != nil {
			return err
		}
		return runtime.sendPromise(ctx, OutboundPromise{
			From:         message.To,
			To:           message.From,
			ProtocolName: EvidenceReportV1,
			Fields: map[string]string{
				"variant":                "evidence_report",
				"promise_about":          "local_compute_observation",
				"subject_pcid":           CIDComputeV1,
				"subject_peer":           message.Fields["result_promiser"],
				"subject_result_cid":     message.Fields["result_cid"],
				"verdict":                "disagree",
				"disagreement_probe":     "true",
				"local_resolution_basis": "alice-local-recompute-plus-dave-report",
			},
		})
	}
	if err := runtime.recordAgentEvent(message.To, "compute_result_peer_verified", "kept", message.From, message.ProtocolName, "verified result_cid="+message.Fields["result_cid"]); err != nil {
		return err
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           message.From,
		ProtocolName: EvidenceReportV1,
		Fields: map[string]string{
			"variant":            "evidence_report",
			"promise_about":      "local_compute_observation",
			"subject_pcid":       CIDComputeV1,
			"subject_peer":       message.Fields["result_promiser"],
			"subject_result_cid": message.Fields["result_cid"],
			"verdict":            "kept",
		},
	})
}

func (runtime *TCPRuntime) handleComputeVerificationResult(message RuntimeMessage) error {
	runtime.adjustTrust(message.To, message.From, 1)
	return runtime.recordAgentEvent(message.To, "compute_verification_received", "kept", message.From, message.ProtocolName, "peer verified result_cid="+message.Fields["result_cid"])
}

// handleEvidenceReport applies a verifier's local observation report to the
// receiving agent's evidence log and trust state.
// Intent: Evidence reports are promises about local observations; they inform
// Alice's local trust without becoming global truth. Source: DI-nisaz
func (runtime *TCPRuntime) handleEvidenceReport(message RuntimeMessage) error {
	verdict := message.Fields["verdict"]
	outcome := "kept"
	if verdict == "disagree" {
		outcome = "non_commitment"
	} else if verdict != "kept" {
		outcome = "malformed"
		runtime.adjustTrust(message.To, message.Fields["subject_peer"], -1)
	} else {
		runtime.adjustTrust(message.To, message.From, 1)
	}
	if err := runtime.recordAgentEvent(message.To, "evidence_report_received", outcome, message.From, message.ProtocolName, "subject_pcid="+message.Fields["subject_pcid"]+" verdict="+verdict+" result_cid="+message.Fields["subject_result_cid"]); err != nil {
		return err
	}
	if verdict == "disagree" {
		if err := runtime.recordAgentEvent(message.To, "compute_verifier_disagreement", "non_commitment", message.From, message.ProtocolName, "Alice receives disagreement from verifier="+message.From+" result_cid="+message.Fields["subject_result_cid"]); err != nil {
			return err
		}
		return runtime.recordAgentEvent(message.To, "compute_disagreement_resolved_locally", "kept", message.Fields["subject_peer"], message.ProtocolName, "Alice resolves by local recompute plus other peer evidence; no global authority")
	}
	if message.Fields["subject_pcid"] == CIDComputeV1 {
		return runtime.recordAgentEvent(message.To, "compute_verification_received", outcome, message.From, message.ProtocolName, "peer verdict="+verdict+" result_cid="+message.Fields["subject_result_cid"])
	}
	return nil
}

func (runtime *TCPRuntime) handleTrustRepairPromise(message RuntimeMessage) error {
	return runtime.recordAgentEvent(message.To, "trust_repair_promise_recorded", "kept", message.From, message.ProtocolName, message.Fields["decision_text"])
}

func (runtime *TCPRuntime) handleKeyRotationPromise(message RuntimeMessage) error {
	return runtime.recordAgentEvent(message.To, "key_rotation_promise_recorded", "kept", message.From, message.ProtocolName, "peer promises future key label="+message.Fields["new_key_label"]+" scope="+message.Fields["rotation_scope"])
}

// handleBadProofEnvelope records a locally rejected proof without accepting the
// payload as a kept promise.
// Intent: Proof failures should become receiver-local evidence, not silent
// crashes, global invalidity judgments, or accepted protocol messages. Source:
// DI-kikoj
func (runtime *TCPRuntime) handleBadProofEnvelope(fields map[string]string, protocolName string, verifyErr error) error {
	recipient := fields["recipient"]
	sender := fields["sender"]
	if _, ok := runtime.states[recipient]; !ok {
		return verifyErr
	}
	runtime.adjustTrust(recipient, sender, -2)
	return runtime.recordAgentEvent(recipient, "bad_proof_rejected", "malformed", sender, protocolName, "signature verification failed for variant="+fields["variant"]+": "+verifyErr.Error())
}

// handleUnknownProtocolEnvelope records local non-commitment for a valid
// envelope whose pCID is not in this runtime's local protocol table.
// Intent: Unknown pCIDs should not become hard authority failures; a receiver
// can simply promise that it does not currently promise that protocol. Source:
// DI-nisaz
func (runtime *TCPRuntime) handleUnknownProtocolEnvelope(message RuntimeMessage) error {
	return runtime.recordAgentEvent(message.To, "unknown_pcid_not_promised", "non_commitment", message.From, message.ProtocolName, "unknown protocol pCID variant="+message.Fields["variant"])
}

// sendComputeVerificationRequest asks one verifier to recompute a result from
// the exact pCID-owned function, input, context, and result bytes.
// Intent: Verification is another voluntary promise exchange; Dave and Grace
// report local observations instead of becoming compute authorities. Source:
// DI-nisaz
func (runtime *TCPRuntime) sendComputeVerificationRequest(ctx context.Context, message RuntimeMessage, verifier string) error {
	fields := map[string]string{
		"variant":         "compute_verification_request",
		"promise_about":   "verify_compute_result",
		"function_cid":    message.Fields["function_cid"],
		"function_b64":    message.Fields["function_b64"],
		"input_cid":       message.Fields["input_cid"],
		"input_b64":       message.Fields["input_b64"],
		"context_cid":     message.Fields["context_cid"],
		"context_b64":     message.Fields["context_b64"],
		"result_cid":      message.Fields["result_cid"],
		"result_b64":      message.Fields["result_b64"],
		"result_promiser": message.From,
	}
	if verifier == "grace" && message.Fields["adversarial_probe"] != "true" {
		fields["disagreement_probe"] = "true"
	}
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         message.To,
		To:           verifier,
		ProtocolName: CIDComputeV1,
		Fields:       fields,
	})
}

// sendReplicaAvailable lets the replica peer announce its own serving promise
// and token, rather than having Bob promise on Frank's behalf.
// Intent: Replica access must be backed by Frank's own token promise before
// Alice treats Frank as a retrieval option. Source: DI-nisaz
func (runtime *TCPRuntime) sendReplicaAvailable(ctx context.Context, from, to, contentCID, replicaToken string) error {
	return runtime.sendPromise(ctx, OutboundPromise{
		From:         from,
		To:           to,
		ProtocolName: CASStorageV1,
		Fields: map[string]string{
			"variant":       "replica_available",
			"promise_about": "replica_location",
			"content_cid":   contentCID,
			"replica_peer":  from,
			"replica_token": replicaToken,
			"ttl_uses":      "1",
		},
	})
}

// sendUnknownProtocolPromise emits a syntactically valid signed envelope with a
// pCID that this POC registry intentionally does not know.
// Intent: POC13 needs an executable unknown-pCID case to prove receivers choose
// local non-commitment rather than crashing or inventing authority. Source:
// DI-nisaz
func (runtime *TCPRuntime) sendUnknownProtocolPromise(ctx context.Context, from, to string, fields map[string]string) error {
	copiedFields := make(map[string]string, len(fields)+3)
	for key, value := range fields {
		copiedFields[key] = value
	}
	copiedFields["act"] = "promise"
	copiedFields["sender"] = from
	copiedFields["recipient"] = to
	protocolCID := NewProtocolCID([]byte("poc13 unknown protocol probe v1"))
	envelope, envelopeErr := NewEnvelope(protocolCID, copiedFields, from)
	if envelopeErr != nil {
		return envelopeErr
	}
	exactBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	conn, dialErr := runtime.dialWithRetry(ctx, OutboundPromise{From: from, To: to, ProtocolName: protocolCID.String(), Fields: copiedFields})
	if dialErr != nil {
		return runtime.recordAgentEvent(from, "tcp_message_send_failed", "non_commitment", to, protocolCID.String(), dialErr.Error())
	}
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "poc13-runtime: close unknown-protocol connection: %v\n", closeErr)
		}
	}()
	if err := (FrameWriter{writer: conn}).WriteFrame(exactBytes); err != nil {
		return err
	}
	runtime.markActivity()
	return runtime.recordAgentEvent(from, "tcp_message_sent", "kept", to, protocolCID.String(), "variant="+copiedFields["variant"]+" exact_sha256="+HashExactBytes(exactBytes))
}

// sendBadProofPromise emits a signed envelope whose proof bytes are corrupted
// after signing while leaving the payload parseable.
// Intent: POC13 needs bad-proof evidence at the receiver boundary without
// treating malformed signatures as accepted promises. Source: DI-kikoj
func (runtime *TCPRuntime) sendBadProofPromise(ctx context.Context, from, to, protocolName string, fields map[string]string) error {
	copiedFields := make(map[string]string, len(fields)+3)
	for key, value := range fields {
		copiedFields[key] = value
	}
	copiedFields["act"] = "promise"
	copiedFields["sender"] = from
	copiedFields["recipient"] = to
	envelope, envelopeErr := NewEnvelope(runtime.Registry.MustCID(protocolName), copiedFields, from)
	if envelopeErr != nil {
		return envelopeErr
	}
	if len(envelope.Proof.Signature) == 0 {
		return fmt.Errorf("cannot corrupt empty signature")
	}
	envelope.Proof.Signature[0] ^= 0x01
	exactBytes, bytesErr := envelope.Bytes()
	if bytesErr != nil {
		return bytesErr
	}
	conn, dialErr := runtime.dialWithRetry(ctx, OutboundPromise{From: from, To: to, ProtocolName: protocolName, Fields: copiedFields})
	if dialErr != nil {
		return runtime.recordAgentEvent(from, "tcp_message_send_failed", "non_commitment", to, protocolName, dialErr.Error())
	}
	defer func() {
		closeErr := conn.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "poc13-runtime: close bad-proof connection: %v\n", closeErr)
		}
	}()
	if err := (FrameWriter{writer: conn}).WriteFrame(exactBytes); err != nil {
		return err
	}
	runtime.markActivity()
	return runtime.recordAgentEvent(from, "bad_proof_sent", "malformed", to, protocolName, "sent intentionally corrupted signature exact_sha256="+HashExactBytes(exactBytes))
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
	conn, dialErr := runtime.dialWithRetry(ctx, outbound)
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

func (runtime *TCPRuntime) dialWithRetry(ctx context.Context, outbound OutboundPromise) (net.Conn, error) {
	address, addressErr := runtime.dialAddressForPromise(outbound)
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

func (runtime *TCPRuntime) persistedTrustHistoryPath() string {
	return filepath.Join(runtime.runtimeMarkerDir(), runtime.Container.Name+".trust-history.json")
}

// loadPersistedTrustHistory seeds and reads one container-local trust file.
// Intent: Dynamic peer choice should use explicit local evidence under the run
// volume, while remaining local to each container and never becoming a global
// trust authority. Source: DI-kikoj
func (runtime *TCPRuntime) loadPersistedTrustHistory() error {
	if err := os.MkdirAll(runtime.runtimeMarkerDir(), 0o755); err != nil {
		return err
	}
	trustHistoryPath := runtime.persistedTrustHistoryPath()
	if _, statErr := os.Stat(trustHistoryPath); statErr != nil {
		if !errors.Is(statErr, os.ErrNotExist) {
			return statErr
		}
		historyBytes, marshalErr := json.MarshalIndent(runtime.defaultPersistedTrustHistory(), "", "  ")
		if marshalErr != nil {
			return marshalErr
		}
		if writeErr := os.WriteFile(trustHistoryPath, append(historyBytes, '\n'), 0o644); writeErr != nil {
			return writeErr
		}
	}
	historyBytes, readErr := os.ReadFile(trustHistoryPath)
	if readErr != nil {
		return readErr
	}
	var history map[string]map[string]int
	if unmarshalErr := json.Unmarshal(historyBytes, &history); unmarshalErr != nil {
		return unmarshalErr
	}
	for agentName, trustByPeer := range history {
		state := runtime.states[agentName]
		if state == nil {
			continue
		}
		for peer, score := range trustByPeer {
			state.Trust[peer] = score
		}
	}
	return nil
}

func (runtime *TCPRuntime) defaultPersistedTrustHistory() map[string]map[string]int {
	history := make(map[string]map[string]int)
	for agentName := range runtime.states {
		history[agentName] = map[string]int{"alice": 1, "bob": 1, "carol": 1, "dave": 1, "ellen": 1, "frank": 1, "grace": 1}
	}
	if _, ok := history["alice"]; ok {
		history["alice"]["bob"] = 4
		history["alice"]["frank"] = 3
		history["alice"]["carol"] = 4
		history["alice"]["dave"] = 2
		history["alice"]["grace"] = 2
	}
	return history
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

func (runtime *TCPRuntime) recordDynamicPeerChoice(agentName, peer, protocolName, detail string) error {
	return runtime.recordAgentEvent(agentName, "dynamic_peer_choice_from_persisted_trust", "kept", peer, protocolName, detail)
}

// recordEconomics writes one local economics evidence record.
// Intent: POC13 economics remain local promise evidence about capacity,
// opportunity cost, and credits; there is no central price authority. Source:
// DI-lupag
func (runtime *TCPRuntime) recordEconomics(agentName, eventName, outcome, peer, protocolName, detail string) error {
	return runtime.recordAgentEvent(agentName, eventName, outcome, peer, protocolName, detail)
}

func (runtime *TCPRuntime) recordTokenLifecycle(agentName, eventName, outcome, peer, protocolName, detail string) error {
	return runtime.recordAgentEvent(agentName, eventName, outcome, peer, protocolName, detail)
}

// recordCapacityPressure records a local refusal when a peer asks for work
// that would exceed the bounded POC capacity.
// Intent: Capacity refusal is an agent's voluntary non-commitment, not an
// externally enforced quota or global scheduler decision. Source: DI-nisaz
func (runtime *TCPRuntime) recordCapacityPressure(agentName, peer, protocolName, detail string) error {
	return runtime.recordEconomics(agentName, "economics_capacity_refused", "non_commitment", peer, protocolName, detail)
}

func computeCacheKeyFromFields(protocolName string, fields map[string]string) string {
	return ComputeCacheKey(protocolName, fields["function_cid"], fields["input_cid"], fields["context_cid"], fields["result_cid"])
}

func copyFields(fields map[string]string) map[string]string {
	copiedFields := make(map[string]string, len(fields))
	for key, value := range fields {
		copiedFields[key] = value
	}
	return copiedFields
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

// badComputeResultBytes creates a hash-valid but semantically wrong result.
// Intent: The verifier path must catch results whose bytes match their own CID
// but do not match the promised function/input/context semantics. Source:
// DI-nisaz
func badComputeResultBytes(correctResultBytes []byte) []byte {
	return append(append([]byte(nil), correctResultBytes...), []byte(";bad-poc13-result")...)
}

// replicaCapabilityKey stores Alice's local copy of a token issued by one
// replica peer for one content CID.
// Intent: Received tokens stay local to the observing agent and do not become a
// global access list or central capability registry. Source: DI-nisaz
func replicaCapabilityKey(replicaPeer, contentCID string) string {
	return replicaPeer + "|" + contentCID
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

func outageAddressForVariant(outageVariant string) string {
	switch outageVariant {
	case OutageVariantContainerStopped:
		return forceTCPUnreachableAddress
	case OutageVariantNetworkPartition:
		return "127.0.0.1:10"
	default:
		return ""
	}
}

// dialAddressForPromise optionally routes one scenario promise to a deliberately
// closed local TCP port so failure is observed at the transport layer.
// Intent: Bob's simulated outage should be a real dial failure in the POC, not
// a cooperative application-level "unavailable" reply. Source: DI-nisaz;
// DI-kikoj
func (runtime *TCPRuntime) dialAddressForPromise(outbound OutboundPromise) (string, error) {
	if outageAddress := outageAddressForVariant(outbound.Fields["outage_variant"]); outageAddress != "" {
		return outageAddress, nil
	}
	return runtime.dialAddressForAgent(outbound.To)
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
	switch functionKind(base64.StdEncoding.EncodeToString(functionBytes)) {
	case "fibonacci":
		if !strings.HasPrefix(inputText, "n=") {
			return nil, fmt.Errorf("unsupported fibonacci input %q", inputText)
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
	case "sum":
		if !strings.HasPrefix(inputText, "values=") {
			return nil, fmt.Errorf("unsupported sum input %q", inputText)
		}
		rawValues := strings.TrimPrefix(inputText, "values=")
		total := 0
		for _, rawValue := range strings.Split(rawValues, ",") {
			value, err := strconv.Atoi(strings.TrimSpace(rawValue))
			if err != nil {
				return nil, err
			}
			total += value
		}
		return []byte(fmt.Sprintf("sum(%s)=%d;context_cid=%s", rawValues, total, ContentCID(contextBytes))), nil
	default:
		return nil, fmt.Errorf("unsupported function source %q", functionText)
	}
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

func functionKind(functionB64 string) string {
	functionBytes, err := base64.StdEncoding.DecodeString(functionB64)
	if err != nil {
		return "unknown"
	}
	functionText := strings.TrimSpace(string(functionBytes))
	if strings.Contains(functionText, "fibonacci") {
		return "fibonacci"
	}
	if strings.Contains(functionText, "sum") {
		return "sum"
	}
	return "unknown"
}
