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
	RetryBudget int
	CAS         *state.CAS
	Writer      *artifact.Writer
	Medium      *radio.Medium
	Peer        string
	OrderNumber string
	Counter     uint64
}

// ReceiveRadio parses exact radio bytes and applies local promise judgment.
func (a *Agent) ReceiveRadio(packet radio.Packet) error {
	msg, err := protocol.Parse(packet.Bytes)
	if err != nil {
		hash, path, writeErr := a.Writer.RecordMalformed(packet.Bytes, "m4")
		if writeErr != nil {
			return writeErr
		}
		return a.Writer.WriteEvent(artifact.Event{
			Type:      "malformed_rejected",
			Actor:     a.Name,
			Peer:      packet.From,
			Hash:      hash,
			Path:      path,
			Transport: "simulated_lora",
			Outcome:   "local_non_commitment",
			Details:   map[string]any{"error": err.Error()},
		})
	}
	hash, rel, err := a.Writer.RecordMessage(packet.Bytes)
	if err != nil {
		return err
	}
	localHash, evicted := a.CAS.Put(packet.Bytes)
	if err := a.Writer.WriteEvent(artifact.Event{
		Type:      "cas_store",
		Actor:     a.Name,
		PCID:      msg.PCID,
		Hash:      localHash,
		Path:      rel,
		Transport: "simulated_lora",
		Outcome:   "stored",
		Details:   map[string]any{"retained": a.CAS.Count(), "artifact_hash": hash},
	}); err != nil {
		return err
	}
	for _, old := range evicted {
		if err := a.Writer.WriteEvent(artifact.Event{Type: "cas_gc", Actor: a.Name, Hash: old, Outcome: "removed"}); err != nil {
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
				if err := a.Writer.WriteEvent(artifact.Event{Type: "missing_parent", Actor: a.Name, Hash: parent, Outcome: "sparse_store_normal"}); err != nil {
					return err
				}
			}
		}
		a.Status = status
		return a.promiseStatus(localHash)
	case protocol.ProtocolOrderStatus:
		return a.receiveOrderStatus(msg.Payload, localHash)
	case protocol.ProtocolLoRaLink:
		return a.Writer.WriteEvent(artifact.Event{Type: "lora_link_promise", Actor: a.Name, PCID: msg.PCID, Outcome: "accepted", Transport: "simulated_lora"})
	default:
		return a.Writer.WriteEvent(artifact.Event{
			Type:      "unknown_pcid_non_commitment",
			Actor:     a.Name,
			PCID:      msg.PCID,
			Hash:      localHash,
			Transport: "simulated_lora",
			Outcome:   "local_non_commitment",
		})
	}
}

func (a *Agent) receiveOrderStatus(data []byte, hash string) error {
	payload, err := protocol.ParseOrderStatusPayload(data)
	if err != nil {
		return fmt.Errorf("parse order status payload: %w", err)
	}
	switch payload.Type {
	case "MSG":
		if payload.Dest != a.Name {
			return a.Writer.WriteEvent(artifact.Event{Type: "order_status_ignored", Actor: a.Name, Peer: payload.Source, PCID: protocol.MustPCIDForName(protocol.ProtocolOrderStatus), Hash: hash, Outcome: "wrong_destination"})
		}
		a.OrderNumber = payload.OrderNumber
		a.Status = payload.Status
		if err := a.Writer.WriteEvent(artifact.Event{
			Type:      "order_status_received",
			Actor:     a.Name,
			Peer:      payload.Source,
			PCID:      protocol.MustPCIDForName(protocol.ProtocolOrderStatus),
			Hash:      hash,
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
		return a.Writer.WriteEvent(artifact.Event{Type: "order_status_non_commitment", Actor: a.Name, PCID: protocol.MustPCIDForName(protocol.ProtocolOrderStatus), Hash: hash, Transport: "simulated_lora", Outcome: "unknown_order_message_type"})
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
			"retry_budget": a.RetryBudget,
			"status":       a.Status,
		},
	}); err != nil {
		return err
	}
	return a.Medium.Send(radio.Packet{From: a.Name, To: a.Peer, Bytes: raw, Label: "status-response"})
}

// PromisePeerStorage records that a constrained device asks a stronger peer for storage.
func (a *Agent) PromisePeerStorage(contentHash string) error {
	return a.Writer.WriteEvent(artifact.Event{
		Type:    "peer_storage_promise",
		Actor:   a.Name,
		Peer:    a.Peer,
		PCID:    protocol.MustPCIDForName(protocol.ProtocolPeerStorage),
		Hash:    contentHash,
		Outcome: "requested",
	})
}
