package device

import (
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/artifact"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/protocol"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/radio"
)

// Peer is a non-M4 radio endpoint; it is not a hidden bridge for the device.
type Peer struct {
	Name   string
	Writer *artifact.Writer
	Medium *radio.Medium
}

// ReceiveRadio records exact messages received over radio.
func (p *Peer) ReceiveRadio(packet radio.Packet) error {
	msg, err := protocol.Parse(packet.Bytes)
	if err != nil {
		hash, path, writeErr := p.Writer.RecordMalformed(packet.Bytes, "peer")
		if writeErr != nil {
			return writeErr
		}
		return p.Writer.WriteEvent(artifact.Event{Type: "peer_malformed_received", Actor: p.Name, Peer: packet.From, Hash: hash, Path: path, Transport: "simulated_lora", Outcome: "review_only", Details: map[string]any{"error": err.Error()}})
	}
	hash, rel, err := p.Writer.RecordMessage(packet.Bytes)
	if err != nil {
		return err
	}
	if err := p.Writer.WriteEvent(artifact.Event{Type: "peer_envelope_received", Actor: p.Name, Peer: packet.From, PCID: msg.PCID, Hash: hash, Path: rel, Transport: "simulated_lora", Outcome: "received"}); err != nil {
		return err
	}
	if msg.ProtocolName != protocol.ProtocolOrderStatus {
		return nil
	}
	payload, err := protocol.ParseOrderStatusPayload(msg.Payload)
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
			Hash:      hash,
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
			Hash:      hash,
			Transport: "simulated_lora",
			Outcome:   "acknowledged",
			Details: map[string]any{
				"counter": payload.Counter,
				"order":   payload.OrderNumber,
				"status":  payload.Status,
			},
		})
	default:
		return p.Writer.WriteEvent(artifact.Event{Type: "peer_order_status_non_commitment", Actor: p.Name, PCID: protocol.MustPCIDForName(protocol.ProtocolOrderStatus), Hash: hash, Transport: "simulated_lora", Outcome: "unknown_order_message_type"})
	}
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
