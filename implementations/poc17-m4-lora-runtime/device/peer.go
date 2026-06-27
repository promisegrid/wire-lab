package device

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/artifact"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/protocol"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/radio"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/state"
)

// Peer is a non-M4 radio endpoint; it is not a hidden bridge for the device.
type Peer struct {
	Name    string
	Writer  *artifact.Writer
	Medium  *radio.Medium
	Storage *state.CAS
	Token   StorageToken
}

// ReceiveRadio records exact messages received over radio.
func (p *Peer) ReceiveRadio(packet radio.Packet) error {
	msg, err := protocol.Parse(packet.Bytes)
	if err != nil {
		contentCID, path, writeErr := p.Writer.RecordMalformed(packet.Bytes, "peer")
		if writeErr != nil {
			return writeErr
		}
		return p.Writer.WriteEvent(artifact.Event{Type: "peer_malformed_received", Actor: p.Name, Peer: packet.From, CID: contentCID, Path: path, Transport: "simulated_lora", Outcome: "review_only", Details: map[string]any{"error": err.Error()}})
	}
	contentCID, rel, err := p.Writer.RecordMessage(packet.Bytes)
	if err != nil {
		return err
	}
	if err := p.Writer.WriteEvent(artifact.Event{Type: "peer_envelope_received", Actor: p.Name, Peer: packet.From, PCID: msg.PCID, CID: contentCID, Path: rel, Transport: "simulated_lora", Outcome: "received"}); err != nil {
		return err
	}
	switch msg.ProtocolName {
	case protocol.ProtocolOrderStatus:
		return p.receiveOrderStatus(msg.Payload, contentCID)
	case protocol.ProtocolPeerStorage:
		return p.receivePeerStorage(msg.Payload, contentCID, packet.From)
	default:
		return nil
	}
}

// GrantPeerStorage sends Bob's compact storage capability to a holder before
// that holder can promise put/get work.
func (p *Peer) GrantPeerStorage(holder string, maxContentBytes uint64, maxRetainedObjects uint64) error {
	token := makeStorageToken(p.Name, holder)
	tokenCID, err := protocol.CIDForBytes(token)
	if err != nil {
		return err
	}
	p.Token = StorageToken{
		Bytes:              token,
		CID:                tokenCID,
		Issuer:             p.Name,
		Holder:             holder,
		MaxContentBytes:    maxContentBytes,
		MaxRetainedObjects: maxRetainedObjects,
		RetentionTerms:     "retain_until_replaced",
	}
	payload, err := protocol.BuildPeerStorageGrant(protocol.PeerStoragePayload{
		Issuer:             p.Name,
		Holder:             holder,
		Token:              token,
		AllowedKinds:       []string{"put", "get"},
		MaxContentBytes:    maxContentBytes,
		MaxRetainedObjects: maxRetainedObjects,
		RetentionTerms:     p.Token.RetentionTerms,
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
	if err := p.Writer.WriteEvent(artifact.Event{Type: "peer_storage_grant_sent", Actor: p.Name, Peer: holder, PCID: protocol.MustPCIDForName(protocol.ProtocolPeerStorage), CID: messageCID, Transport: "simulated_lora", Outcome: "capability_promised", Details: map[string]any{"token_cid": tokenCID, "max_content_bytes": maxContentBytes, "max_retained_objects": maxRetainedObjects}}); err != nil {
		return err
	}
	return p.Medium.Send(radio.Packet{From: p.Name, To: holder, Bytes: raw, Label: "peer-storage-grant"})
}

func (p *Peer) receiveOrderStatus(data []byte, contentCID string) error {
	payload, err := protocol.ParseOrderStatusPayload(data)
	if err != nil {
		return err
	}
	switch payload.Type {
	case "MSG":
		if err := p.Writer.WriteEvent(artifact.Event{
			Type:      "peer_order_status_received",
			Actor:     p.Name,
			Peer:      payload.Source,
			PCID:      protocol.MustPCIDForName(protocol.ProtocolOrderStatus),
			CID:       contentCID,
			Transport: "simulated_lora",
			Outcome:   "database_update_promised",
			Details: map[string]any{
				"counter": payload.Counter,
				"order":   payload.OrderNumber,
				"status":  payload.Status,
			},
		}); err != nil {
			return err
		}
		return p.sendOrderAck(payload)
	case "ACK":
		return p.Writer.WriteEvent(artifact.Event{
			Type:      "peer_order_ack_received",
			Actor:     p.Name,
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
		return p.Writer.WriteEvent(artifact.Event{Type: "peer_order_status_non_commitment", Actor: p.Name, PCID: protocol.MustPCIDForName(protocol.ProtocolOrderStatus), CID: contentCID, Transport: "simulated_lora", Outcome: "unknown_order_message_type"})
	}
}

func (p *Peer) receivePeerStorage(data []byte, messageCID string, from string) error {
	payload, err := protocol.ParsePeerStoragePayload(data)
	if err != nil {
		return fmt.Errorf("parse peer_storage payload: %w", err)
	}
	switch payload.Kind {
	case protocol.PeerStoragePut:
		return p.receivePeerStoragePut(payload, messageCID, from)
	case protocol.PeerStorageGet:
		return p.receivePeerStorageGet(payload, messageCID, from)
	default:
		return p.Writer.WriteEvent(artifact.Event{Type: "peer_storage_ignored", Actor: p.Name, Peer: from, PCID: protocol.MustPCIDForName(protocol.ProtocolPeerStorage), CID: messageCID, Transport: "simulated_lora", Outcome: "not_for_peer"})
	}
}

func (p *Peer) receivePeerStoragePut(payload protocol.PeerStoragePayload, messageCID string, from string) error {
	if reason := p.validateStorageToken(payload.Holder, payload.Token); reason != "" {
		return p.sendPutResult(from, payload.ContentCID, messageCID, false, reason)
	}
	contentCID, err := protocol.CIDForBytes(payload.Content)
	if err != nil {
		return err
	}
	if contentCID != payload.ContentCID {
		return p.sendPutResult(from, payload.ContentCID, messageCID, false, "content_cid_mismatch")
	}
	if uint64(len(payload.Content)) > p.Token.MaxContentBytes {
		return p.sendPutResult(from, payload.ContentCID, messageCID, false, "content_too_large")
	}
	if p.Storage == nil {
		p.Storage = state.NewCAS(int(p.Token.MaxRetainedObjects))
	}
	localCID, evicted := p.Storage.Put(payload.Content)
	for _, old := range evicted {
		if err := p.Writer.WriteEvent(artifact.Event{Type: "peer_storage_gc", Actor: p.Name, CID: old, Outcome: "removed"}); err != nil {
			return err
		}
	}
	if err := p.Writer.WriteEvent(artifact.Event{Type: "peer_storage_put_received", Actor: p.Name, Peer: from, PCID: protocol.MustPCIDForName(protocol.ProtocolPeerStorage), CID: localCID, Transport: "simulated_lora", Outcome: "content_retained", Details: map[string]any{"message_cid": messageCID, "token_cid": p.Token.CID, "retained": p.Storage.Count()}}); err != nil {
		return err
	}
	return p.sendPutResult(from, localCID, messageCID, true, p.Token.RetentionTerms)
}

func (p *Peer) receivePeerStorageGet(payload protocol.PeerStoragePayload, messageCID string, from string) error {
	if reason := p.validateStorageToken(payload.Holder, payload.Token); reason != "" {
		return p.sendGetRefusal(from, payload.ContentCID, messageCID, reason)
	}
	if p.Storage == nil {
		return p.sendGetRefusal(from, payload.ContentCID, messageCID, "not_retained")
	}
	content, ok := p.Storage.Get(payload.ContentCID)
	if !ok {
		return p.sendGetRefusal(from, payload.ContentCID, messageCID, "not_retained")
	}
	if err := p.Writer.WriteEvent(artifact.Event{Type: "peer_storage_get_received", Actor: p.Name, Peer: from, PCID: protocol.MustPCIDForName(protocol.ProtocolPeerStorage), CID: payload.ContentCID, Transport: "simulated_lora", Outcome: "content_available", Details: map[string]any{"message_cid": messageCID, "token_cid": p.Token.CID}}); err != nil {
		return err
	}
	payloadBytes, err := protocol.BuildPeerStorageGetResult(protocol.PeerStoragePayload{Issuer: p.Name, Holder: from, TokenCID: p.Token.CID, RelatedMessageCID: messageCID, ContentCID: payload.ContentCID, Content: content})
	if err != nil {
		return err
	}
	raw, err := protocol.Build(protocol.Message{ProtocolName: protocol.ProtocolPeerStorage, Payload: payloadBytes})
	if err != nil {
		return err
	}
	return p.Medium.Send(radio.Packet{From: p.Name, To: from, Bytes: raw, Label: "peer-storage-get-fulfill"})
}

func (p *Peer) sendPutResult(to string, contentCID string, relatedMessageCID string, accepted bool, note string) error {
	payload, err := protocol.BuildPeerStoragePutResult(protocol.PeerStoragePayload{Issuer: p.Name, Holder: to, TokenCID: p.Token.CID, RelatedMessageCID: relatedMessageCID, ContentCID: contentCID, Accepted: accepted, Reason: note})
	if err != nil {
		return err
	}
	raw, err := protocol.Build(protocol.Message{ProtocolName: protocol.ProtocolPeerStorage, Payload: payload})
	if err != nil {
		return err
	}
	eventType := "peer_storage_put_refused"
	outcome := "refused"
	if accepted {
		eventType = "peer_storage_put_accepted"
		outcome = "accepted"
	}
	if err := p.Writer.WriteEvent(artifact.Event{Type: eventType, Actor: p.Name, Peer: to, PCID: protocol.MustPCIDForName(protocol.ProtocolPeerStorage), CID: contentCID, Transport: "simulated_lora", Outcome: outcome, Details: map[string]any{"related_message_cid": relatedMessageCID, "token_cid": p.Token.CID, "note": note}}); err != nil {
		return err
	}
	return p.Medium.Send(radio.Packet{From: p.Name, To: to, Bytes: raw, Label: "peer-storage-put-result"})
}

func (p *Peer) sendGetRefusal(to string, contentCID string, relatedMessageCID string, reason string) error {
	payload, err := protocol.BuildPeerStorageGetRefusal(protocol.PeerStoragePayload{Issuer: p.Name, Holder: to, TokenCID: p.Token.CID, RelatedMessageCID: relatedMessageCID, ContentCID: contentCID, Reason: reason})
	if err != nil {
		return err
	}
	raw, err := protocol.Build(protocol.Message{ProtocolName: protocol.ProtocolPeerStorage, Payload: payload})
	if err != nil {
		return err
	}
	if err := p.Writer.WriteEvent(artifact.Event{Type: "peer_storage_get_refused", Actor: p.Name, Peer: to, PCID: protocol.MustPCIDForName(protocol.ProtocolPeerStorage), CID: contentCID, Transport: "simulated_lora", Outcome: "refused", Details: map[string]any{"related_message_cid": relatedMessageCID, "token_cid": p.Token.CID, "reason": reason}}); err != nil {
		return err
	}
	return p.Medium.Send(radio.Packet{From: p.Name, To: to, Bytes: raw, Label: "peer-storage-get-refuse"})
}

func (p *Peer) validateStorageToken(holder string, token []byte) string {
	if holder != p.Token.Holder {
		return "wrong_holder"
	}
	if !bytes.Equal(token, p.Token.Bytes) {
		return "invalid_token"
	}
	return ""
}

func makeStorageToken(issuer string, holder string) []byte {
	sum := sha256.Sum256([]byte("poc17 peer_storage capability:" + issuer + ":" + holder))
	return sum[:16]
}

func (p *Peer) sendOrderAck(payload protocol.OrderStatusPayload) error {
	ackPayload, err := protocol.BuildOrderStatusPayload(protocol.OrderStatusPayload{
		Type:        "ACK",
		Source:      p.Name,
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
	// Intent: Model the bintags gateway ACK without letting the gateway become a hidden reliable bridge. Source: DI-mokit
	return p.Medium.Send(radio.Packet{From: p.Name, To: payload.Source, Bytes: raw, Label: "order-ack"})
}
