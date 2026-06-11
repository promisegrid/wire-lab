package boundary

import (
	"bytes"
	"fmt"
)

const (
	PromiseAboutWASMBoundary    = "wasm_boundary_evidence"
	PromiseAboutStdioBoundary   = "stdio_boundary_evidence"
	PromiseAboutLocalSummary    = "local_evidence_summary"
	PromiseAboutPeerAttestation = "peer_carried_attestation"
	PromiseAboutExchangeRate    = "bearer_token_exchange_rate"
	PromiseAboutTopologySignal  = "relationship_topology_signal"
	PromiseAboutVoluntaryGossip = "voluntary_gossip"
)

// MinimalWASMModule is a valid empty WebAssembly module. POC14 uses it as a
// small sandbox-boundary fixture so the WASM-role process can prove that it is
// handling module bytes without adding a broad WASM runtime dependency yet.
// Intent: Exercise the process/runtime boundary called for by POC14 while
// keeping PromiseGrid semantics in ordinary pCID-defined envelopes. Source:
// DI-linof
var MinimalWASMModule = []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}

// ValidateWASMModule checks the stable magic and version bytes for the module
// fixture used by the POC14 WASM process-boundary agent.
func ValidateWASMModule(moduleBytes []byte) error {
	if len(moduleBytes) < len(MinimalWASMModule) {
		return fmt.Errorf("wasm module too short: %d bytes", len(moduleBytes))
	}
	if !bytes.Equal(moduleBytes[:4], MinimalWASMModule[:4]) {
		return fmt.Errorf("wasm module magic mismatch")
	}
	if !bytes.Equal(moduleBytes[4:8], MinimalWASMModule[4:8]) {
		return fmt.Errorf("wasm module version mismatch")
	}
	return nil
}

// PromiseFields returns the common relationship_v1 payload used by the
// heterogeneous boundary agents. The caller still signs and routes the envelope
// through the normal POC14 app/kernel path.
func PromiseFields(fromAgent, toAgent, promiseAbout, promiseText string) map[string]string {
	return map[string]string{
		"act":                 "promise",
		"from":                fromAgent,
		"to":                  toAgent,
		"turn":                "startup",
		"promise":             promiseText,
		"reason":              "heterogeneous runtime boundary evidence expressed as a local promise",
		"field_promise_about": promiseAbout,
	}
}

// StdioRequest is the adapter-to-worker request sent over stdin. It is not a
// wire protocol; it is the local subprocess adapter contract for this POC.
type StdioRequest struct {
	Type string `json:"type"`
	From string `json:"from"`
	To   string `json:"to"`
}

// StdioEnvelopeMessage is the worker-to-adapter outbound promise envelope. The
// worker only writes bytes to stdout; the adapter chooses whether to forward the
// exact envelope through the local kernel.
type StdioEnvelopeMessage struct {
	Type     string `json:"type"`
	From     string `json:"from"`
	To       string `json:"to"`
	Protocol string `json:"protocol"`
	Hex      string `json:"hex"`
}

// StdioAckMessage is the adapter-to-worker acknowledgement envelope after the
// adapter receives peer evidence through the local kernel.
type StdioAckMessage struct {
	Type string `json:"type"`
	Hex  string `json:"hex"`
}

// StdioObservedMessage is the worker's final stdout evidence after it parses
// the ACK envelope locally.
type StdioObservedMessage struct {
	Type        string `json:"type"`
	Outcome     string `json:"outcome"`
	ExactSHA256 string `json:"exact_sha256"`
}
