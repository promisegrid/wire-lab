package poc13

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Event is one local evidence record emitted by a POC13 agent.
type Event struct {
	Observer string `json:"observer"`
	Event    string `json:"event"`
	Outcome  string `json:"outcome"`
	Peer     string `json:"peer,omitempty"`
	PCID     string `json:"pcid,omitempty"`
	Detail   string `json:"detail"`
}

// Runner executes one local agent's bounded POC13 script.
type Runner struct {
	Config   Config
	Agent    AgentConfig
	Registry Registry
	Decider  Decider
}

// NewRunner constructs an agent runner with the POC13 protocol registry.
func NewRunner(cfg Config, agent AgentConfig, decider Decider) Runner {
	return Runner{Config: cfg, Agent: agent, Registry: NewRegistry(), Decider: decider}
}

// Run emits one bounded set of CAS/compute promise evidence for this agent.
// Intent: POC13 keeps each agent autonomous and local while using deterministic
// Go validation for pCID, CID, signature, and cache-key evidence. Source:
// DI-notig
func (runner Runner) Run(ctx context.Context) error {
	if runner.Decider == nil {
		runner.Decider = LiveOrScriptedDecider{}
	}
	log, logErr := openAgentLog(runner.Config, runner.Agent.Name)
	if logErr != nil {
		return logErr
	}
	defer func() {
		closeErr := log.Close()
		if closeErr != nil {
			fmt.Fprintf(os.Stderr, "poc13: close agent log: %v\n", closeErr)
		}
	}()
	recorder := &Recorder{log: log, agent: runner.Agent, registry: runner.Registry}
	recorder.Record("app_receive_promise_registered", "kept", "", "", "cas_storage_v1 and cid_compute_v1 payloads are accepted only as local promises")
	decision, decisionErr := runner.Decider.Decide(ctx, runner.Config, runner.Agent, "Choose which local storage or compute promise you can make in this bounded POC13 turn.")
	if decisionErr != nil {
		recorder.Record("llm_decision_error", "non_commitment", "", "", decisionErr.Error())
	} else {
		recorder.Record("llm_decision_"+decision.Mode, "kept", "", "", decision.Text)
	}
	return runner.runRole(recorder)
}

func (runner Runner) runRole(recorder *Recorder) error {
	switch runner.Agent.Role {
	case "data_holder":
		return runner.runAlice(recorder)
	case "storage_peer":
		return runner.runBob(recorder)
	case "compute_peer":
		return runner.runCarol(recorder)
	case "cache_peer":
		return runner.runDave(recorder)
	case "context_peer":
		return runner.runEllen(recorder)
	case "replication_peer":
		return runner.runFrank(recorder)
	case "verifier_peer":
		return runner.runGrace(recorder)
	case "adversary_peer":
		return runner.runMallory(recorder)
	default:
		recorder.Record("role_not_promised", "non_commitment", "", "", "agent has no POC13 role script")
		return nil
	}
}

func (runner Runner) runAlice(recorder *Recorder) error {
	contentCID := ContentCID(sampleContentBytes())
	inputCID := ContentCID(sampleInputBytes())
	functionCID := ContentCID(sampleFunctionBytes())
	contextCID := ContentCID(sampleContextBytes())
	if err := recorder.RecordEnvelope(CASStorageV1, "bob", map[string]string{
		"act":           "promise",
		"variant":       "store_request",
		"promise_about": "store_content",
		"content_cid":   contentCID,
	}); err != nil {
		return err
	}
	recorder.Record("cas_store_requested", "kept", "bob", CASStorageV1, "Alice asks Bob for a storage promise over content_cid="+contentCID)
	return recorder.RecordEnvelope(CIDComputeV1, "carol", map[string]string{
		"act":           "promise",
		"variant":       "compute_request",
		"promise_about": "execute_function",
		"function_cid":  functionCID,
		"input_cid":     inputCID,
		"context_cid":   contextCID,
	})
}

func (runner Runner) runBob(recorder *Recorder) error {
	contentCID := ContentCID(sampleContentBytes())
	if !VerifyContentCID(sampleContentBytes(), contentCID) {
		return fmt.Errorf("sample content cid verification failed")
	}
	if err := recorder.RecordEnvelope(CASStorageV1, "alice", map[string]string{
		"act":           "promise",
		"variant":       "store_acceptance",
		"promise_about": "store_content",
		"content_cid":   contentCID,
		"retention":     "bounded-poc13-run",
	}); err != nil {
		return err
	}
	recorder.Record("cas_storage_promised", "kept", "alice", CASStorageV1, "Bob promises local storage for content_cid="+contentCID)
	recorder.Record("cas_retention_promised", "kept", "alice", CASStorageV1, "Bob promises retention only for the bounded POC13 run")
	recorder.Record("cas_serve_promised", "kept", "alice", CASStorageV1, "Bob promises to serve bytes to trusted peers during this run")
	return nil
}

func (runner Runner) runCarol(recorder *Recorder) error {
	functionCID := ContentCID(sampleFunctionBytes())
	inputCID := ContentCID(sampleInputBytes())
	contextCID := ContentCID(sampleContextBytes())
	resultBytes := computeSampleResult()
	resultCID := ContentCID(resultBytes)
	if err := recorder.RecordEnvelope(CIDComputeV1, "alice", map[string]string{
		"act":           "promise",
		"variant":       "compute_result",
		"promise_about": "execute_function",
		"function_cid":  functionCID,
		"input_cid":     inputCID,
		"context_cid":   contextCID,
		"result_cid":    resultCID,
	}); err != nil {
		return err
	}
	recorder.Record("cid_compute_promised", "kept", "alice", CIDComputeV1, "Carol promises compute only for the stated function, input, and context CIDs")
	recorder.Record("compute_result_promised", "kept", "alice", CIDComputeV1, "result_cid="+resultCID)
	return nil
}

func (runner Runner) runDave(recorder *Recorder) error {
	cacheKey := ComputeCacheKey(CIDComputeV1, ContentCID(sampleFunctionBytes()), ContentCID(sampleInputBytes()), ContentCID(sampleContextBytes()), ContentCID(computeSampleResult()))
	recorder.Record("compute_cache_checkpointed", "kept", "carol", CIDComputeV1, "Dave caches only the exact protocol/function/input/context/result tuple cache_key="+cacheKey)
	return nil
}

func (runner Runner) runEllen(recorder *Recorder) error {
	contextCID := ContentCID(sampleContextBytes())
	recorder.Record("compute_context_promised", "kept", "carol", CIDComputeV1, "Ellen promises an explicit timestamp/context object context_cid="+contextCID)
	return nil
}

func (runner Runner) runFrank(recorder *Recorder) error {
	contentCID := ContentCID(sampleContentBytes())
	recorder.Record("cas_replication_promised", "kept", "bob", CASStorageV1, "Frank promises one local replica for content_cid="+contentCID)
	return nil
}

func (runner Runner) runGrace(recorder *Recorder) error {
	recorder.Record("cas_verification_promised", "kept", "alice", CASStorageV1, "Grace promises to verify bytes against claimed CIDs from her local vantage")
	if VerifyContentCID(corruptContentBytes(), ContentCID(sampleContentBytes())) {
		return fmt.Errorf("corrupt bytes unexpectedly matched original content cid")
	}
	recorder.Record("cas_corrupt_evidence_recorded", "kept", "mallory", CASStorageV1, "Grace records Mallory's corrupt-byte evidence locally")
	return nil
}

func (runner Runner) runMallory(recorder *Recorder) error {
	claimedCID := ContentCID(sampleContentBytes())
	if VerifyContentCID(corruptContentBytes(), claimedCID) {
		return fmt.Errorf("corrupt bytes unexpectedly matched claimed cid")
	}
	recorder.Record("cas_corrupt_bytes_rejected", "malformed", "grace", CASStorageV1, "Mallory presented bytes that did not match claimed content_cid="+claimedCID)
	return nil
}

// Recorder writes JSONL evidence for one local agent.
type Recorder struct {
	log      *os.File
	agent    AgentConfig
	registry Registry
}

// Record appends one local evidence event.
func (recorder *Recorder) Record(eventName, outcome, peer, protocolName, detail string) {
	event := Event{Observer: recorder.agent.Name, Event: eventName, Outcome: outcome, Peer: peer, PCID: protocolName, Detail: detail}
	eventBytes, marshalErr := json.Marshal(event)
	if marshalErr != nil {
		fmt.Fprintf(os.Stderr, "poc13: marshal event: %v\n", marshalErr)
		return
	}
	if _, writeErr := recorder.log.Write(append(eventBytes, '\n')); writeErr != nil {
		fmt.Fprintf(os.Stderr, "poc13: write event: %v\n", writeErr)
	}
}

// RecordEnvelope builds, parses, verifies, and records one signed grid envelope.
func (recorder *Recorder) RecordEnvelope(protocolName, peer string, fields map[string]string) error {
	if fields["act"] != "promise" {
		return fmt.Errorf("top-level act must be promise")
	}
	protocolCID := recorder.registry.MustCID(protocolName)
	envelope, envelopeErr := NewEnvelope(protocolCID, fields, recorder.agent.Name)
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
	fieldsAfterParse, fieldsErr := parsed.PayloadFields()
	if fieldsErr != nil {
		return fieldsErr
	}
	if fieldsAfterParse["act"] != "promise" {
		return fmt.Errorf("parsed envelope act must be promise")
	}
	recorder.Record("promise_envelope_validated", "kept", peer, protocolName, "exact_sha256="+HashExactBytes(exactBytes)+" promise_about="+fieldsAfterParse["promise_about"])
	return nil
}

func openAgentLog(cfg Config, agentName string) (*os.File, error) {
	runDir := filepath.Join(cfg.RunRoot, cfg.RunID, "run")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		return nil, err
	}
	return os.OpenFile(filepath.Join(runDir, agentName+".jsonl"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
}

func sampleContentBytes() []byte {
	return []byte("poc13 sample invoice bytes")
}

func corruptContentBytes() []byte {
	return []byte("poc13 sample invoice bytes corrupted by mallory")
}

func sampleFunctionBytes() []byte {
	return []byte("poc13 function: fibonacci(n) v1")
}

func sampleInputBytes() []byte {
	return []byte("n=7")
}

func sampleContextBytes() []byte {
	return []byte("timestamp=2026-06-06T00:00:00Z;randomness=explicit-none")
}

func computeSampleResult() []byte {
	return []byte("fibonacci(7)=13")
}

// ComputeCacheKey names one exact local compute checkpoint.
func ComputeCacheKey(protocolName, functionCID, inputCID, contextCID, resultCID string) string {
	return ContentCID([]byte(protocolName + "|" + functionCID + "|" + inputCID + "|" + contextCID + "|" + resultCID))
}
