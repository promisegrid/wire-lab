package analyzer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/artifact"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/protocol"
)

// Summary is the analyzer's gate report.
type Summary struct {
	Events              int `json:"events"`
	RadioSends          int `json:"radio_sends"`
	RadioReceives       int `json:"radio_receives"`
	MessageArtifacts    int `json:"message_artifacts"`
	MalformedArtifacts  int `json:"malformed_artifacts"`
	MTURefusals         int `json:"mtu_refusals"`
	NonCommitments      int `json:"non_commitments"`
	CASStores           int `json:"cas_stores"`
	CASGC               int `json:"cas_gc"`
	PeerStorageGrants   int `json:"peer_storage_grants"`
	PeerStoragePuts     int `json:"peer_storage_puts"`
	PeerStoragePutAcks  int `json:"peer_storage_put_acks"`
	PeerStorageGets     int `json:"peer_storage_gets"`
	PeerStorageGetAcks  int `json:"peer_storage_get_acks"`
	FidelityNotices     int `json:"fidelity_notices"`
	OrderStatusEvents   int `json:"order_status_events"`
	LifecycleIssued     int `json:"lifecycle_issued"`
	LifecycleInvoked    int `json:"lifecycle_invoked"`
	LifecycleFulfilled  int `json:"lifecycle_fulfilled"`
	LifecycleRejected   int `json:"lifecycle_rejected"`
	LifecycleArtifacts  int `json:"lifecycle_artifacts"`
	ResourcePromises    int `json:"resource_promises"`
	ResourceWithdrawals int `json:"resource_withdrawals"`
	ResourceSnapshots   int `json:"resource_snapshots"`
	DeviceRestarts      int `json:"device_restarts"`
	DeviceRecoveries    int `json:"device_recoveries"`
}

// Analyze checks the first POC17 behavior evidence gates.
func Analyze(runDir string) (Summary, error) {
	events, err := readEvents(filepath.Join(runDir, "events.jsonl"))
	if err != nil {
		return Summary{}, err
	}
	var summary Summary
	for _, event := range events {
		summary.Events++
		if event.PCID == protocol.LocalLifecycleV1PCID && event.Transport == "simulated_lora" {
			return Summary{}, fmt.Errorf("local lifecycle token appeared on simulated LoRa path")
		}
		if strings.Contains(strings.ToLower(fmt.Sprint(event.Details)), "sigterm") || strings.Contains(strings.ToLower(fmt.Sprint(event.Details)), "sigkill") {
			return Summary{}, fmt.Errorf("clean lifecycle path used signal fallback")
		}
		legacyLimitWord := "bud" + "get"
		if strings.Contains(strings.ToLower(fmt.Sprint(event.Details)), legacyLimitWord) || strings.Contains(strings.ToLower(event.Outcome), legacyLimitWord) {
			return Summary{}, fmt.Errorf("event %s used legacy limit wording", event.Type)
		}
		switch event.Type {
		case "radio_send":
			summary.RadioSends++
			if event.Transport != "simulated_lora" {
				return Summary{}, fmt.Errorf("radio_send without simulated_lora transport")
			}
		case "radio_receive":
			summary.RadioReceives++
			if event.Transport != "simulated_lora" {
				return Summary{}, fmt.Errorf("radio_receive without simulated_lora transport")
			}
		case "peer_envelope_received":
			if event.Path != "" {
				summary.MessageArtifacts++
				if _, err := os.Stat(filepath.Join(runDir, event.Path)); err != nil {
					return Summary{}, fmt.Errorf("missing message artifact %s: %w", event.Path, err)
				}
			}
		case "malformed_rejected":
			summary.MalformedArtifacts++
			summary.NonCommitments++
			if _, err := os.Stat(filepath.Join(runDir, event.Path)); err != nil {
				return Summary{}, fmt.Errorf("missing malformed artifact %s: %w", event.Path, err)
			}
		case "peer_malformed_received":
			summary.MalformedArtifacts++
			if _, err := os.Stat(filepath.Join(runDir, event.Path)); err != nil {
				return Summary{}, fmt.Errorf("missing malformed artifact %s: %w", event.Path, err)
			}
		case "radio_mtu_refused":
			summary.MTURefusals++
		case "unknown_pcid_non_commitment":
			summary.NonCommitments++
		case "cas_store":
			summary.CASStores++
			if event.Path != "" {
				summary.MessageArtifacts++
				if _, err := os.Stat(filepath.Join(runDir, event.Path)); err != nil {
					return Summary{}, fmt.Errorf("missing message artifact %s: %w", event.Path, err)
				}
			}
		case "cas_gc":
			summary.CASGC++
		case "peer_storage_grant_sent", "peer_storage_grant_received":
			summary.PeerStorageGrants++
			if event.PCID != protocol.MustPCIDForName(protocol.ProtocolPeerStorage) {
				return Summary{}, fmt.Errorf("peer storage grant missing peer_storage pCID")
			}
		case "peer_storage_put_promised", "peer_storage_put_received":
			summary.PeerStoragePuts++
			if event.PCID != protocol.MustPCIDForName(protocol.ProtocolPeerStorage) {
				return Summary{}, fmt.Errorf("peer storage put missing peer_storage pCID")
			}
		case "peer_storage_put_accepted", "peer_storage_put_refused":
			summary.PeerStoragePutAcks++
			if event.PCID != protocol.MustPCIDForName(protocol.ProtocolPeerStorage) {
				return Summary{}, fmt.Errorf("peer storage put result missing peer_storage pCID")
			}
		case "peer_storage_get_promised", "peer_storage_get_received":
			summary.PeerStorageGets++
			if event.PCID != protocol.MustPCIDForName(protocol.ProtocolPeerStorage) {
				return Summary{}, fmt.Errorf("peer storage get missing peer_storage pCID")
			}
		case "peer_storage_get_fulfilled", "peer_storage_get_refused":
			summary.PeerStorageGetAcks++
			if event.PCID != protocol.MustPCIDForName(protocol.ProtocolPeerStorage) {
				return Summary{}, fmt.Errorf("peer storage get result missing peer_storage pCID")
			}
		case "simulator_fidelity_notice":
			summary.FidelityNotices++
		case "order_status_received", "order_status_promise", "order_ack_received", "peer_order_status_received", "peer_order_ack_received":
			summary.OrderStatusEvents++
		case "lifecycle_token_issued":
			summary.LifecycleIssued++
			if err := checkLifecycleArtifact(runDir, event); err != nil {
				return Summary{}, err
			}
			summary.LifecycleArtifacts++
			if detailString(event, "token_cid") == "" {
				return Summary{}, fmt.Errorf("lifecycle token issue missing token CID")
			}
		case "lifecycle_token_invoked":
			summary.LifecycleInvoked++
			if err := checkLifecycleArtifact(runDir, event); err != nil {
				return Summary{}, err
			}
			summary.LifecycleArtifacts++
		case "lifecycle_token_fulfilled":
			summary.LifecycleFulfilled++
			if err := checkLifecycleArtifact(runDir, event); err != nil {
				return Summary{}, err
			}
			summary.LifecycleArtifacts++
		case "lifecycle_token_rejected":
			summary.LifecycleRejected++
			if detailString(event, "reason") == "" {
				return Summary{}, fmt.Errorf("lifecycle rejection missing reason")
			}
		case "resource_access_promised":
			summary.ResourcePromises++
		case "resource_access_withdrawn":
			summary.ResourceWithdrawals++
			if detailString(event, "scope") != "host_local" {
				return Summary{}, fmt.Errorf("resource withdrawal missing host-local scope")
			}
			if event.Details["not_command_authority"] != true || event.Details["not_peer_trust_evidence"] != true {
				return Summary{}, fmt.Errorf("resource withdrawal drifted into authority or peer-trust evidence")
			}
		case "resource_limit_snapshot":
			summary.ResourceSnapshots++
		case "device_restart_started":
			summary.DeviceRestarts++
		case "device_recovery_loaded", "device_recovery_verified":
			summary.DeviceRecoveries++
		}
	}
	if summary.RadioSends == 0 || summary.RadioReceives == 0 {
		return summary, fmt.Errorf("missing radio send/receive evidence")
	}
	if summary.MessageArtifacts == 0 {
		return summary, fmt.Errorf("missing exact CBOR artifacts")
	}
	if summary.MalformedArtifacts == 0 || summary.MTURefusals == 0 || summary.NonCommitments == 0 {
		return summary, fmt.Errorf("missing failure-path evidence")
	}
	if summary.PeerStorageGrants < 2 || summary.PeerStoragePuts < 2 || summary.PeerStoragePutAcks < 2 || summary.PeerStorageGets < 2 || summary.PeerStorageGetAcks == 0 {
		return summary, fmt.Errorf("missing peer-storage grant/put/get evidence")
	}
	if summary.CASGC == 0 {
		return summary, fmt.Errorf("missing local CAS GC evidence")
	}
	if summary.FidelityNotices == 0 {
		return summary, fmt.Errorf("missing simulator fidelity notice")
	}
	if summary.OrderStatusEvents < 4 {
		return summary, fmt.Errorf("missing production-like order status evidence")
	}
	if summary.LifecycleIssued < 2 || summary.LifecycleInvoked < 2 || summary.LifecycleFulfilled < 2 {
		return summary, fmt.Errorf("missing host-local lifecycle token evidence")
	}
	if summary.LifecycleRejected < 5 {
		return summary, fmt.Errorf("missing lifecycle token rejection evidence")
	}
	if summary.ResourcePromises < 2 || summary.ResourceWithdrawals == 0 {
		return summary, fmt.Errorf("missing resource promise or withdrawal evidence")
	}
	if summary.ResourceSnapshots == 0 {
		return summary, fmt.Errorf("missing resource limit evidence")
	}
	if summary.DeviceRestarts == 0 || summary.DeviceRecoveries < 2 {
		return summary, fmt.Errorf("missing fresh-agent restart recovery evidence")
	}
	return summary, nil
}

func checkLifecycleArtifact(runDir string, event artifact.Event) error {
	if event.Path == "" || !strings.HasPrefix(event.Path, "lifecycle-cas/") {
		return fmt.Errorf("lifecycle event missing lifecycle artifact path")
	}
	if _, err := os.Stat(filepath.Join(runDir, event.Path)); err != nil {
		return fmt.Errorf("missing lifecycle artifact %s: %w", event.Path, err)
	}
	if detailString(event, "host_local_only") != "true" {
		return fmt.Errorf("lifecycle event missing host-local marker")
	}
	return nil
}

func detailString(event artifact.Event, key string) string {
	value, ok := event.Details[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	case bool:
		if typed {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprint(typed)
	}
}

func readEvents(path string) ([]artifact.Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open events: %w", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "close events file: %v\n", closeErr)
		}
	}()
	var events []artifact.Event
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var event artifact.Event
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			return nil, fmt.Errorf("decode event: %w", err)
		}
		events = append(events, event)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan events: %w", err)
	}
	return events, nil
}
