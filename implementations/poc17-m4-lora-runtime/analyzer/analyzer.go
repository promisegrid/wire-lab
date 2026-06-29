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
	// Intent: Turn POC17's radio-only, CID-first, and promise-shaped claims into
	// explicit clean-run gates instead of accepting count-only evidence.
	// Source: DI-gidul; DI-mokit; DI-dutah
	gates := newGateState()
	for index, event := range events {
		summary.Events++
		if event.PCID == protocol.LocalLifecycleV1PCID && event.Transport == "simulated_lora" {
			return Summary{}, fmt.Errorf("local lifecycle token appeared on simulated LoRa path")
		}
		if err := gates.checkNoAuthorityDrift(event); err != nil {
			return Summary{}, err
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
			gates.noteRadioSend(event)
		case "radio_receive":
			summary.RadioReceives++
			if event.Transport != "simulated_lora" {
				return Summary{}, fmt.Errorf("radio_receive without simulated_lora transport")
			}
			gates.noteRadioReceive(event)
		case "peer_envelope_received":
			if event.Path != "" {
				summary.MessageArtifacts++
				if err := gates.checkMessageArtifact(runDir, event); err != nil {
					return Summary{}, err
				}
			}
		case "malformed_rejected":
			summary.MalformedArtifacts++
			summary.NonCommitments++
			if _, err := os.Stat(filepath.Join(runDir, event.Path)); err != nil {
				return Summary{}, fmt.Errorf("missing malformed artifact %s: %w", event.Path, err)
			}
			gates.malformedRejected = true
		case "peer_malformed_received":
			summary.MalformedArtifacts++
			if _, err := os.Stat(filepath.Join(runDir, event.Path)); err != nil {
				return Summary{}, fmt.Errorf("missing malformed artifact %s: %w", event.Path, err)
			}
		case "radio_mtu_refused":
			summary.MTURefusals++
			if event.Transport != "simulated_lora" {
				return Summary{}, fmt.Errorf("MTU refusal without simulated_lora transport")
			}
			gates.mtuRefused = true
		case "unknown_pcid_non_commitment":
			summary.NonCommitments++
			if event.Transport != "simulated_lora" {
				return Summary{}, fmt.Errorf("unknown pCID non-commitment without simulated_lora transport")
			}
			gates.unknownPCIDNonCommitment = true
		case "cas_store":
			summary.CASStores++
			if event.Path != "" {
				summary.MessageArtifacts++
				if err := gates.checkMessageArtifact(runDir, event); err != nil {
					return Summary{}, err
				}
			}
		case "cas_gc":
			summary.CASGC++
		case "radio_lost":
			gates.radioLost = true
		case "radio_delayed":
			gates.radioDelayed = true
		case "radio_asymmetric_unreachable":
			gates.asymmetricUnreachable = true
		case "peer_storage_grant_sent", "peer_storage_grant_received":
			summary.PeerStorageGrants++
			if event.PCID != protocol.MustPCIDForName(protocol.ProtocolPeerStorage) {
				return Summary{}, fmt.Errorf("peer storage grant missing peer_storage pCID")
			}
			if err := gates.notePeerStorage(event, index); err != nil {
				return Summary{}, err
			}
		case "peer_storage_put_promised", "peer_storage_put_received":
			summary.PeerStoragePuts++
			if event.PCID != protocol.MustPCIDForName(protocol.ProtocolPeerStorage) {
				return Summary{}, fmt.Errorf("peer storage put missing peer_storage pCID")
			}
			if err := gates.notePeerStorage(event, index); err != nil {
				return Summary{}, err
			}
		case "peer_storage_put_accepted", "peer_storage_put_refused":
			summary.PeerStoragePutAcks++
			if event.PCID != protocol.MustPCIDForName(protocol.ProtocolPeerStorage) {
				return Summary{}, fmt.Errorf("peer storage put result missing peer_storage pCID")
			}
			if err := gates.notePeerStorage(event, index); err != nil {
				return Summary{}, err
			}
		case "peer_storage_get_promised", "peer_storage_get_received":
			summary.PeerStorageGets++
			if event.PCID != protocol.MustPCIDForName(protocol.ProtocolPeerStorage) {
				return Summary{}, fmt.Errorf("peer storage get missing peer_storage pCID")
			}
			if err := gates.notePeerStorage(event, index); err != nil {
				return Summary{}, err
			}
		case "peer_storage_get_fulfilled", "peer_storage_get_refused":
			summary.PeerStorageGetAcks++
			if event.PCID != protocol.MustPCIDForName(protocol.ProtocolPeerStorage) {
				return Summary{}, fmt.Errorf("peer storage get result missing peer_storage pCID")
			}
			if err := gates.notePeerStorage(event, index); err != nil {
				return Summary{}, err
			}
		case "simulator_fidelity_notice":
			summary.FidelityNotices++
		case "order_status_received", "order_status_promise", "order_ack_received", "peer_order_status_received", "peer_order_ack_received":
			summary.OrderStatusEvents++
			if err := gates.noteOrderStatus(event); err != nil {
				return Summary{}, err
			}
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
	if err := gates.finish(); err != nil {
		return summary, err
	}
	return summary, nil
}

type gateState struct {
	radioSends                  map[string]bool
	radioReceives               map[string]bool
	orderStatuses               map[string]bool
	orderACKs                   map[string]bool
	peerStorageGrantSent        bool
	peerStorageGrantReceived    bool
	peerStoragePutPromised      bool
	peerStoragePutReceived      bool
	peerStoragePutAccepted      bool
	peerStorageGetPromised      bool
	peerStorageGetReceived      bool
	peerStorageGetFulfilled     bool
	peerStorageTokenCID         string
	peerStoragePutMessageCID    string
	peerStorageGetMessageCID    string
	peerStoragePutAcceptedAt    int
	peerStorageGetFulfilledAt   int
	peerStorageGetPromisedAt    int
	mtuRefused                  bool
	malformedRejected           bool
	unknownPCIDNonCommitment    bool
	radioLost                   bool
	radioDelayed                bool
	asymmetricUnreachable       bool
	messageArtifactProtocolCIDs map[string]bool
}

// newGateState tracks cross-event evidence that cannot be validated from one
// event row alone, such as request/result ordering and exact artifact pCIDs.
func newGateState() *gateState {
	return &gateState{
		radioSends:                  make(map[string]bool),
		radioReceives:               make(map[string]bool),
		orderStatuses:               make(map[string]bool),
		orderACKs:                   make(map[string]bool),
		peerStoragePutAcceptedAt:    -1,
		peerStorageGetFulfilledAt:   -1,
		peerStorageGetPromisedAt:    -1,
		messageArtifactProtocolCIDs: make(map[string]bool),
	}
}

func (g *gateState) noteRadioSend(event artifact.Event) {
	g.radioSends[detailString(event, "label")] = true
}

func (g *gateState) noteRadioReceive(event artifact.Event) {
	g.radioReceives[detailString(event, "label")] = true
}

func (g *gateState) checkMessageArtifact(runDir string, event artifact.Event) error {
	path := filepath.Join(runDir, event.Path)
	raw, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("missing message artifact %s: %w", event.Path, err)
	}
	artifactCID, err := protocol.CIDForBytes(raw)
	if err != nil {
		return fmt.Errorf("CID message artifact %s: %w", event.Path, err)
	}
	eventArtifactCID := detailString(event, "artifact_cid")
	if eventArtifactCID == "" {
		eventArtifactCID = event.CID
	}
	if eventArtifactCID != artifactCID {
		return fmt.Errorf("message artifact %s CID mismatch: event %s computed %s", event.Path, eventArtifactCID, artifactCID)
	}
	msg, err := protocol.Parse(raw)
	if err != nil {
		return fmt.Errorf("parse message artifact %s: %w", event.Path, err)
	}
	if event.PCID != "" && event.PCID != msg.PCID {
		return fmt.Errorf("message artifact %s pCID mismatch: event %s artifact %s", event.Path, event.PCID, msg.PCID)
	}
	if err := protocol.ValidateCIDText(msg.PCID); err != nil {
		return fmt.Errorf("message artifact %s has invalid pCID CID: %w", event.Path, err)
	}
	item, err := protocol.Decode(raw)
	if err != nil {
		return fmt.Errorf("decode message artifact %s: %w", event.Path, err)
	}
	if item.Tag == nil || item.Tag.Number != protocol.GridTag || len(item.Tag.Value.Array) < 2 {
		return fmt.Errorf("message artifact %s is not a grid envelope", event.Path)
	}
	slotZero := item.Tag.Value.Array[0].Tag
	if slotZero == nil || slotZero.Number != protocol.CIDTag || slotZero.Value.Bytes == nil {
		return fmt.Errorf("message artifact %s slot 0 is not tag-42 CID bytes", event.Path)
	}
	if slotZero.Value.Text != nil {
		return fmt.Errorf("message artifact %s slot 0 regressed to text pCID", event.Path)
	}
	g.messageArtifactProtocolCIDs[msg.PCID] = true
	return nil
}

func (g *gateState) noteOrderStatus(event artifact.Event) error {
	if event.Transport != "simulated_lora" {
		return fmt.Errorf("%s without simulated_lora transport", event.Type)
	}
	if event.PCID != protocol.MustPCIDForName(protocol.ProtocolOrderStatus) {
		return fmt.Errorf("%s missing order_status pCID", event.Type)
	}
	order := detailString(event, "order")
	status := detailString(event, "status")
	if order != "BT-1042" || status == "" || detailString(event, "counter") == "" {
		return fmt.Errorf("%s missing production-like order details", event.Type)
	}
	switch event.Type {
	case "order_status_received", "order_status_promise", "peer_order_status_received":
		g.orderStatuses[status] = true
	case "order_ack_received", "peer_order_ack_received":
		g.orderACKs[status] = true
	}
	return nil
}

func (g *gateState) notePeerStorage(event artifact.Event, index int) error {
	if event.Transport != "simulated_lora" {
		return fmt.Errorf("%s without simulated_lora transport", event.Type)
	}
	tokenCID := detailString(event, "token_cid")
	if tokenCID != "" {
		if err := protocol.ValidateCIDText(tokenCID); err != nil {
			return fmt.Errorf("%s has invalid token CID: %w", event.Type, err)
		}
		if g.peerStorageTokenCID == "" {
			g.peerStorageTokenCID = tokenCID
		} else if g.peerStorageTokenCID != tokenCID {
			return fmt.Errorf("peer_storage token CID changed from %s to %s", g.peerStorageTokenCID, tokenCID)
		}
	}
	switch event.Type {
	case "peer_storage_grant_sent":
		g.peerStorageGrantSent = true
	case "peer_storage_grant_received":
		g.peerStorageGrantReceived = true
	case "peer_storage_put_promised":
		g.peerStoragePutPromised = true
		g.peerStoragePutMessageCID = detailString(event, "message_cid")
		if g.peerStoragePutMessageCID == "" {
			return fmt.Errorf("peer_storage put promise missing message CID")
		}
	case "peer_storage_put_received":
		g.peerStoragePutReceived = true
		if detailString(event, "message_cid") != g.peerStoragePutMessageCID {
			return fmt.Errorf("peer_storage put receive did not reference Ivan's put message")
		}
	case "peer_storage_put_accepted":
		g.peerStoragePutAccepted = true
		g.peerStoragePutAcceptedAt = index
		if detailString(event, "related_message_cid") != g.peerStoragePutMessageCID {
			return fmt.Errorf("peer_storage put acknowledgement did not reference Ivan's put message")
		}
	case "peer_storage_get_promised":
		g.peerStorageGetPromised = true
		g.peerStorageGetPromisedAt = index
		g.peerStorageGetMessageCID = detailString(event, "message_cid")
		if g.peerStorageGetMessageCID == "" {
			return fmt.Errorf("peer_storage get promise missing message CID")
		}
	case "peer_storage_get_received":
		g.peerStorageGetReceived = true
		if detailString(event, "message_cid") != g.peerStorageGetMessageCID {
			return fmt.Errorf("peer_storage get receive did not reference Ivan's get message")
		}
	case "peer_storage_get_fulfilled":
		g.peerStorageGetFulfilled = true
		g.peerStorageGetFulfilledAt = index
		if detailString(event, "related_message_cid") != g.peerStorageGetMessageCID {
			return fmt.Errorf("peer_storage get fulfillment did not reference Ivan's get message")
		}
		if event.Outcome != "content_verified" {
			return fmt.Errorf("peer_storage fulfillment did not verify returned bytes by CID")
		}
	case "peer_storage_put_refused", "peer_storage_get_refused":
		if event.Outcome != "refused" {
			return fmt.Errorf("%s must be a refusal outcome", event.Type)
		}
	}
	return nil
}

func (g *gateState) checkNoAuthorityDrift(event artifact.Event) error {
	if event.Transport != "" && event.Transport != "simulated_lora" {
		if event.Transport != "host_local" {
			return fmt.Errorf("unexpected transport %q on %s", event.Transport, event.Type)
		}
	}
	for _, text := range eventStrings(event) {
		lower := strings.ToLower(text)
		for _, forbidden := range []string{"host bridge", "global trust", "authorization service", "permission service", "command authority", "monitor authority", "rpc controller"} {
			if strings.Contains(lower, forbidden) {
				return fmt.Errorf("authority-drift wording %q in %s", forbidden, event.Type)
			}
		}
	}
	return nil
}

func eventStrings(event artifact.Event) []string {
	values := []string{event.Outcome}
	for _, value := range event.Details {
		switch typed := value.(type) {
		case string:
			values = append(values, typed)
		case []any:
			for _, child := range typed {
				if text, ok := child.(string); ok {
					values = append(values, text)
				}
			}
		}
	}
	return values
}

func (g *gateState) finish() error {
	for _, label := range []string{"peer-storage-grant", "peer-storage-put", "peer-storage-put-result", "peer-storage-get", "peer-storage-get-fulfill", "order-reset-bt-1042", "order-ack"} {
		if !g.radioSends[label] || !g.radioReceives[label] {
			return fmt.Errorf("missing radio-only send/receive evidence for %s", label)
		}
	}
	for _, status := range []string{"created", "cut", "stripped", "soldered", "completed"} {
		if !g.orderStatuses[status] {
			return fmt.Errorf("missing order status %s", status)
		}
		if !g.orderACKs[status] {
			return fmt.Errorf("missing order ACK %s", status)
		}
	}
	if !g.peerStorageGrantSent || !g.peerStorageGrantReceived {
		return fmt.Errorf("missing Bob-issued peer_storage capability grant")
	}
	if g.peerStorageTokenCID == "" {
		return fmt.Errorf("missing peer_storage token CID")
	}
	if !g.peerStoragePutPromised || !g.peerStoragePutReceived || !g.peerStoragePutAccepted {
		return fmt.Errorf("missing peer_storage put promise, receipt, or acknowledgement")
	}
	if !g.peerStorageGetPromised || !g.peerStorageGetReceived || !g.peerStorageGetFulfilled {
		return fmt.Errorf("missing peer_storage get promise, receipt, or fulfillment")
	}
	if g.peerStoragePutAcceptedAt >= g.peerStorageGetFulfilledAt {
		return fmt.Errorf("peer_storage get fulfilled before put acknowledgement")
	}
	if g.peerStorageGetPromisedAt >= g.peerStorageGetFulfilledAt {
		return fmt.Errorf("peer_storage fulfilled before Ivan promised get")
	}
	for _, pcid := range []string{
		protocol.MustPCIDForName(protocol.ProtocolOrderStatus),
		protocol.MustPCIDForName(protocol.ProtocolPeerStorage),
		protocol.MustPCIDForName(protocol.ProtocolDeviceStatus),
		protocol.MustPCIDForName(protocol.ProtocolLoRaLink),
	} {
		if !g.messageArtifactProtocolCIDs[pcid] {
			return fmt.Errorf("missing exact CBOR artifact for pCID %s", pcid)
		}
	}
	if !g.mtuRefused || !g.malformedRejected || !g.unknownPCIDNonCommitment || !g.radioLost || !g.radioDelayed || !g.asymmetricUnreachable {
		return fmt.Errorf("missing expected refusal/loss/asymmetric workaround evidence")
	}
	return nil
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
