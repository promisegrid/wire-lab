package analyzer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/artifact"
)

// Summary is the analyzer's gate report.
type Summary struct {
	Events             int `json:"events"`
	RadioSends         int `json:"radio_sends"`
	RadioReceives      int `json:"radio_receives"`
	MessageArtifacts   int `json:"message_artifacts"`
	MalformedArtifacts int `json:"malformed_artifacts"`
	MTURefusals        int `json:"mtu_refusals"`
	NonCommitments     int `json:"non_commitments"`
	CASStores          int `json:"cas_stores"`
	CASGC              int `json:"cas_gc"`
	PeerStorage        int `json:"peer_storage"`
	FidelityNotices    int `json:"fidelity_notices"`
	OrderStatusEvents  int `json:"order_status_events"`
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
		case "peer_storage_promise":
			summary.PeerStorage++
		case "simulator_fidelity_notice":
			summary.FidelityNotices++
		case "order_status_received", "order_status_promise", "order_ack_received", "peer_order_status_received", "peer_order_ack_received":
			summary.OrderStatusEvents++
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
	if summary.PeerStorage == 0 {
		return summary, fmt.Errorf("missing peer-storage promise evidence")
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
	return summary, nil
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
