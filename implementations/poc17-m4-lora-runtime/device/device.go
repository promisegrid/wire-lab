package device

import (
	"fmt"

	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/artifact"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/protocol"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/radio"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/state"
)

// Agent models the first M4-shaped PromiseGrid runtime loop.
type Agent struct {
	Name        string
	Status      string
	Battery     uint64
	RetryLimit  int
	CAS         *state.CAS
	Writer      *artifact.Writer
	Medium      *radio.Medium
	Peer        string
	OrderNumber string
	Counter     uint64
	Token       StorageToken
	MissingCIDs []string
}

// Snapshot is the fresh-agent restart handoff state retained outside volatile
// M4 memory in this behavior simulator.
type Snapshot struct {
	Name         string
	Status       string
	Battery      uint64
	RetryLimit   int
	Peer         string
	OrderNumber  string
	Counter      uint64
	Token        StorageToken
	MissingCIDs  []string
	RetainedCIDs []string
}

// Snapshot records the minimum state a fresh agent needs after restart.
func (a *Agent) Snapshot() Snapshot {
	return Snapshot{
		Name:         a.Name,
		Status:       a.Status,
		Battery:      a.Battery,
		RetryLimit:   a.RetryLimit,
		Peer:         a.Peer,
		OrderNumber:  a.OrderNumber,
		Counter:      a.Counter,
		Token:        a.Token,
		MissingCIDs:  append([]string(nil), a.MissingCIDs...),
		RetainedCIDs: a.CAS.CIDs(),
	}
}

// NewAgentFromSnapshot creates a fresh process with durable identity and token
// state but without pretending volatile CAS bytes survived the restart.
func NewAgentFromSnapshot(snapshot Snapshot, casLimit int, writer *artifact.Writer, medium *radio.Medium) *Agent {
	return &Agent{
		Name:        snapshot.Name,
		Status:      snapshot.Status,
		Battery:     snapshot.Battery,
		RetryLimit:  snapshot.RetryLimit,
		CAS:         state.NewCAS(casLimit),
		Writer:      writer,
		Medium:      medium,
		Peer:        snapshot.Peer,
		OrderNumber: snapshot.OrderNumber,
		Counter:     snapshot.Counter,
		Token:       snapshot.Token,
		MissingCIDs: append([]string(nil), snapshot.MissingCIDs...),
	}
}

// ReceiveRadio parses exact radio bytes and applies local promise judgment.
func (a *Agent) ReceiveRadio(packet radio.Packet) error {
	msg, err := protocol.Parse(packet.Bytes)
	if err != nil {
		contentCID, path, writeErr := a.Writer.RecordMalformed(packet.Bytes, "m4")
		if writeErr != nil {
			return writeErr
		}
		return a.Writer.WriteEvent(artifact.Event{
			Type:      "malformed_rejected",
			Actor:     a.Name,
			Peer:      packet.From,
			CID:       contentCID,
			Path:      path,
			Transport: "simulated_lora",
			Outcome:   "local_non_commitment",
			Details:   map[string]any{"error": err.Error()},
		})
	}
	artifactCID, rel, err := a.Writer.RecordMessage(packet.Bytes)
	if err != nil {
		return err
	}
	localCID, evicted := a.CAS.Put(packet.Bytes)
	if err := a.Writer.WriteEvent(artifact.Event{
		Type:      "cas_store",
		Actor:     a.Name,
		PCID:      msg.PCID,
		CID:       localCID,
		Path:      rel,
		Transport: "simulated_lora",
		Outcome:   "stored",
		Details:   map[string]any{"retained": a.CAS.Count(), "artifact_cid": artifactCID},
	}); err != nil {
		return err
	}
	for _, old := range evicted {
		if err := a.Writer.WriteEvent(artifact.Event{Type: "cas_gc", Actor: a.Name, CID: old, Outcome: "removed"}); err != nil {
			return err
		}
	}
	switch msg.ProtocolName {
	case protocol.ProtocolDeviceStatus:
		_, status, _, parents, err := protocol.ParseStatusPayload(msg.Payload)
		if err != nil {
			return fmt.Errorf("parse status payload: %w", err)
		}
		for _, parent := range parents {
			if !a.CAS.Has(parent) {
				if err := a.Writer.WriteEvent(artifact.Event{Type: "missing_parent", Actor: a.Name, CID: parent, Outcome: "sparse_store_normal"}); err != nil {
					return err
				}
			}
		}
		a.Status = status
		return a.promiseStatus(localCID)
	case protocol.ProtocolOrderStatus:
		return a.receiveOrderStatus(msg.Payload, localCID)
	case protocol.ProtocolLoRaLink:
		return a.Writer.WriteEvent(artifact.Event{Type: "lora_link_promise", Actor: a.Name, PCID: msg.PCID, Outcome: "accepted", Transport: "simulated_lora"})
	case protocol.ProtocolPeerStorage:
		return a.receivePeerStorage(msg.Payload, localCID, packet.From)
	default:
		return a.Writer.WriteEvent(artifact.Event{
			Type:      "unknown_pcid_non_commitment",
			Actor:     a.Name,
			PCID:      msg.PCID,
			CID:       localCID,
			Transport: "simulated_lora",
			Outcome:   "local_non_commitment",
		})
	}
}

func (a *Agent) receiveOrderStatus(data []byte, contentCID string) error {
	payload, err := protocol.ParseOrderStatusPayload(data)
	if err != nil {
		return fmt.Errorf("parse order status payload: %w", err)
	}
	switch payload.Type {
	case "MSG":
		if payload.Dest != a.Name {
			return a.Writer.WriteEvent(artifact.Event{Type: "order_status_ignored", Actor: a.Name, Peer: payload.Source, PCID: protocol.MustPCIDForName(protocol.ProtocolOrderStatus), CID: contentCID, Outcome: "wrong_destination"})
		}
		a.OrderNumber = payload.OrderNumber
		a.Status = payload.Status
		if err := a.Writer.WriteEvent(artifact.Event{
			Type:      "order_status_received",
			Actor:     a.Name,
			Peer:      payload.Source,
			PCID:      protocol.MustPCIDForName(protocol.ProtocolOrderStatus),
			CID:       contentCID,
			Transport: "simulated_lora",
			Outcome:   "display_update_promised",
			Details: map[string]any{
				"counter": payload.Counter,
				"order":   payload.OrderNumber,
				"status":  payload.Status,
			},
		}); err != nil {
			return err
		}
		return a.sendOrderAck(payload)
	case "ACK":
		return a.Writer.WriteEvent(artifact.Event{
			Type:      "order_ack_received",
			Actor:     a.Name,
			Peer:      payload.Source,
			PCID:      protocol.MustPCIDForName(protocol.ProtocolOrderStatus),
			CID:       contentCID,
			Transport: "simulated_lora",
			Outcome:   "acknowledged",
			Details: map[string]any{
				"counter": payload.Counter,
				"order":   payload.OrderNumber,
				"status":  payload.Status,
			},
		})
	default:
		return a.Writer.WriteEvent(artifact.Event{Type: "order_status_non_commitment", Actor: a.Name, PCID: protocol.MustPCIDForName(protocol.ProtocolOrderStatus), CID: contentCID, Transport: "simulated_lora", Outcome: "unknown_order_message_type"})
	}
}

func (a *Agent) sendOrderAck(payload protocol.OrderStatusPayload) error {
	ackPayload, err := protocol.BuildOrderStatusPayload(protocol.OrderStatusPayload{
		Type:        "ACK",
		Source:      a.Name,
		Dest:        payload.Source,
		Counter:     payload.Counter,
		OrderNumber: payload.OrderNumber,
		Status:      payload.Status,
	})
	if err != nil {
		return err
	}
	raw, err := protocol.Build(protocol.Message{ProtocolName: protocol.ProtocolOrderStatus, Payload: ackPayload})
	if err != nil {
		return err
	}
	// Intent: Keep the bintags MSG/ACK behavior, but send the ACK as CBOR over the simulated LoRa path. Source: DI-mokit
	return a.Medium.Send(radio.Packet{From: a.Name, To: payload.Source, Bytes: raw, Label: "order-ack"})
}

// PromiseOrderStatus simulates a button-driven bintags order status update.
func (a *Agent) PromiseOrderStatus(status string) error {
	if a.OrderNumber == "" {
		return fmt.Errorf("cannot promise order status without an assigned order")
	}
	a.Counter++
	payload, err := protocol.BuildOrderStatusPayload(protocol.OrderStatusPayload{
		Type:        "MSG",
		Source:      a.Name,
		Dest:        a.Peer,
		Counter:     a.Counter,
		OrderNumber: a.OrderNumber,
		Status:      status,
	})
	if err != nil {
		return err
	}
	raw, err := protocol.Build(protocol.Message{ProtocolName: protocol.ProtocolOrderStatus, Payload: payload})
	if err != nil {
		return err
	}
	a.Status = status
	if err := a.Writer.WriteEvent(artifact.Event{
		Type:      "order_status_promise",
		Actor:     a.Name,
		Peer:      a.Peer,
		PCID:      protocol.MustPCIDForName(protocol.ProtocolOrderStatus),
		Transport: "simulated_lora",
		Outcome:   "promised",
		Details: map[string]any{
			"counter": a.Counter,
			"order":   a.OrderNumber,
			"status":  status,
		},
	}); err != nil {
		return err
	}
	return a.Medium.Send(radio.Packet{From: a.Name, To: a.Peer, Bytes: raw, Label: "order-status-" + status})
}

func (a *Agent) promiseStatus(parent string) error {
	payload, err := protocol.BuildStatusPayload(a.Name, a.Status, a.Battery, []string{parent})
	if err != nil {
		return err
	}
	raw, err := protocol.Build(protocol.Message{ProtocolName: protocol.ProtocolDeviceStatus, Payload: payload})
	if err != nil {
		return err
	}
	if err := a.Writer.WriteEvent(artifact.Event{
		Type:    "device_status_promise",
		Actor:   a.Name,
		PCID:    protocol.MustPCIDForName(protocol.ProtocolDeviceStatus),
		Outcome: "promised",
		Details: map[string]any{
			"retry_limit": a.RetryLimit,
			"status":      a.Status,
		},
	}); err != nil {
		return err
	}
	return a.Medium.Send(radio.Packet{From: a.Name, To: a.Peer, Bytes: raw, Label: "status-response"})
}

// PutPeerStorage promises exact bytes to Bob under a Bob-issued capability.
func (a *Agent) PutPeerStorage(content []byte, reason string) error {
	if len(a.Token.Bytes) == 0 {
		return fmt.Errorf("missing peer-storage capability token")
	}
	contentCID, err := protocol.CIDForBytes(content)
	if err != nil {
		return err
	}
	payload, err := protocol.BuildPeerStoragePut(protocol.PeerStoragePayload{
		Holder:     a.Name,
		Issuer:     a.Peer,
		Token:      a.Token.Bytes,
		ContentCID: contentCID,
		Content:    content,
		Reason:     reason,
	})
	if err != nil {
		return err
	}
	raw, err := protocol.Build(protocol.Message{ProtocolName: protocol.ProtocolPeerStorage, Payload: payload})
	if err != nil {
		return err
	}
	messageCID, err := protocol.CIDForBytes(raw)
	if err != nil {
		return err
	}
	if err := a.Writer.WriteEvent(artifact.Event{
		Type:      "peer_storage_put_promised",
		Actor:     a.Name,
		Peer:      a.Peer,
		PCID:      protocol.MustPCIDForName(protocol.ProtocolPeerStorage),
		CID:       contentCID,
		Transport: "simulated_lora",
		Outcome:   "promised",
		Details:   map[string]any{"message_cid": messageCID, "token_cid": a.Token.CID, "reason": reason},
	}); err != nil {
		return err
	}
	return a.Medium.Send(radio.Packet{From: a.Name, To: a.Peer, Bytes: raw, Label: "peer-storage-put"})
}

// GetPeerStorage asks Bob to return exact bytes by CID under a prior token.
func (a *Agent) GetPeerStorage(contentCID string, reason string) error {
	if len(a.Token.Bytes) == 0 {
		return fmt.Errorf("missing peer-storage capability token")
	}
	payload, err := protocol.BuildPeerStorageGet(protocol.PeerStoragePayload{
		Holder:     a.Name,
		Issuer:     a.Peer,
		Token:      a.Token.Bytes,
		ContentCID: contentCID,
		Reason:     reason,
	})
	if err != nil {
		return err
	}
	raw, err := protocol.Build(protocol.Message{ProtocolName: protocol.ProtocolPeerStorage, Payload: payload})
	if err != nil {
		return err
	}
	messageCID, err := protocol.CIDForBytes(raw)
	if err != nil {
		return err
	}
	if err := a.Writer.WriteEvent(artifact.Event{
		Type:      "peer_storage_get_promised",
		Actor:     a.Name,
		Peer:      a.Peer,
		PCID:      protocol.MustPCIDForName(protocol.ProtocolPeerStorage),
		CID:       contentCID,
		Transport: "simulated_lora",
		Outcome:   "promised",
		Details:   map[string]any{"message_cid": messageCID, "token_cid": a.Token.CID, "reason": reason},
	}); err != nil {
		return err
	}
	return a.Medium.Send(radio.Packet{From: a.Name, To: a.Peer, Bytes: raw, Label: "peer-storage-get"})
}

func (a *Agent) receivePeerStorage(data []byte, messageCID string, from string) error {
	payload, err := protocol.ParsePeerStoragePayload(data)
	if err != nil {
		return fmt.Errorf("parse peer_storage payload: %w", err)
	}
	switch payload.Kind {
	case protocol.PeerStorageGrant:
		if payload.Holder != a.Name || payload.Issuer != from {
			return a.Writer.WriteEvent(artifact.Event{Type: "peer_storage_grant_refused", Actor: a.Name, Peer: from, PCID: protocol.MustPCIDForName(protocol.ProtocolPeerStorage), CID: messageCID, Transport: "simulated_lora", Outcome: "wrong_holder_or_issuer"})
		}
		tokenCID, err := protocol.CIDForBytes(payload.Token)
		if err != nil {
			return err
		}
		a.Token = StorageToken{Bytes: payload.Token, CID: tokenCID, Issuer: payload.Issuer, Holder: payload.Holder, MaxContentBytes: payload.MaxContentBytes, MaxRetainedObjects: payload.MaxRetainedObjects, RetentionTerms: payload.RetentionTerms}
		return a.Writer.WriteEvent(artifact.Event{Type: "peer_storage_grant_received", Actor: a.Name, Peer: from, PCID: protocol.MustPCIDForName(protocol.ProtocolPeerStorage), CID: messageCID, Transport: "simulated_lora", Outcome: "capability_recorded", Details: map[string]any{"token_cid": tokenCID, "allowed": payload.AllowedKinds, "max_content_bytes": payload.MaxContentBytes, "max_retained_objects": payload.MaxRetainedObjects}})
	case protocol.PeerStoragePutResult:
		eventType := "peer_storage_put_refused"
		outcome := "refused"
		if payload.Accepted {
			eventType = "peer_storage_put_accepted"
			outcome = "accepted"
		}
		return a.Writer.WriteEvent(artifact.Event{Type: eventType, Actor: a.Name, Peer: from, PCID: protocol.MustPCIDForName(protocol.ProtocolPeerStorage), CID: payload.ContentCID, Transport: "simulated_lora", Outcome: outcome, Details: map[string]any{"message_cid": messageCID, "related_message_cid": payload.RelatedMessageCID, "token_cid": payload.TokenCID, "note": payload.Reason}})
	case protocol.PeerStorageGetResult:
		contentCID, err := protocol.CIDForBytes(payload.Content)
		if err != nil {
			return err
		}
		if contentCID != payload.ContentCID {
			return a.Writer.WriteEvent(artifact.Event{Type: "peer_storage_content_rejected", Actor: a.Name, Peer: from, PCID: protocol.MustPCIDForName(protocol.ProtocolPeerStorage), CID: payload.ContentCID, Transport: "simulated_lora", Outcome: "cid_mismatch", Details: map[string]any{"computed_cid": contentCID, "message_cid": messageCID}})
		}
		localCID, evicted := a.CAS.Put(payload.Content)
		for _, old := range evicted {
			if err := a.Writer.WriteEvent(artifact.Event{Type: "cas_gc", Actor: a.Name, CID: old, Outcome: "removed"}); err != nil {
				return err
			}
		}
		return a.Writer.WriteEvent(artifact.Event{Type: "peer_storage_get_fulfilled", Actor: a.Name, Peer: from, PCID: protocol.MustPCIDForName(protocol.ProtocolPeerStorage), CID: localCID, Transport: "simulated_lora", Outcome: "content_verified", Details: map[string]any{"message_cid": messageCID, "related_message_cid": payload.RelatedMessageCID, "token_cid": payload.TokenCID}})
	case protocol.PeerStorageGetRefusal:
		return a.Writer.WriteEvent(artifact.Event{Type: "peer_storage_get_refused", Actor: a.Name, Peer: from, PCID: protocol.MustPCIDForName(protocol.ProtocolPeerStorage), CID: payload.ContentCID, Transport: "simulated_lora", Outcome: "refused", Details: map[string]any{"message_cid": messageCID, "related_message_cid": payload.RelatedMessageCID, "token_cid": payload.TokenCID, "reason": payload.Reason}})
	default:
		return a.Writer.WriteEvent(artifact.Event{Type: "peer_storage_ignored", Actor: a.Name, Peer: from, PCID: protocol.MustPCIDForName(protocol.ProtocolPeerStorage), CID: messageCID, Transport: "simulated_lora", Outcome: "not_for_agent"})
	}
}
