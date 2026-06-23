package radio

import (
	"fmt"

	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/artifact"
)

// Packet carries exact bytes through the simulated radio path only.
type Packet struct {
	From  string
	To    string
	Bytes []byte
	Label string
}

// Receiver accepts a delivered radio packet.
type Receiver interface {
	ReceiveRadio(Packet) error
}

// Medium models RFM95/SX127x-shaped packet effects without judging promises.
type Medium struct {
	MTU       int
	writer    *artifact.Writer
	receivers map[string]Receiver
	links     map[string]bool
}

// NewMedium creates a deterministic radio transport model.
func NewMedium(mtu int, writer *artifact.Writer) *Medium {
	return &Medium{
		MTU:       mtu,
		writer:    writer,
		receivers: make(map[string]Receiver),
		links:     make(map[string]bool),
	}
}

// Register attaches a radio endpoint to the medium.
func (m *Medium) Register(name string, receiver Receiver) {
	m.receivers[name] = receiver
}

// SetReachable controls asymmetric reachability.
func (m *Medium) SetReachable(from, to string, ok bool) {
	m.links[from+"->"+to] = ok
}

// Send applies radio effects before delivering exact bytes.
func (m *Medium) Send(packet Packet, effects ...string) error {
	if len(packet.Bytes) > m.MTU {
		return m.writer.WriteEvent(artifact.Event{
			Type:      "radio_mtu_refused",
			Actor:     packet.From,
			Peer:      packet.To,
			Transport: "simulated_lora",
			Outcome:   "refused",
			Details: map[string]any{
				"bytes": len(packet.Bytes),
				"mtu":   m.MTU,
			},
		})
	}
	if ok, configured := m.links[packet.From+"->"+packet.To]; configured && !ok {
		return m.writer.WriteEvent(artifact.Event{
			Type:      "radio_asymmetric_unreachable",
			Actor:     packet.From,
			Peer:      packet.To,
			Transport: "simulated_lora",
			Outcome:   "not_delivered",
		})
	}
	for _, effect := range effects {
		switch effect {
		case "loss":
			return m.writer.WriteEvent(artifact.Event{Type: "radio_lost", Actor: packet.From, Peer: packet.To, Transport: "simulated_lora", Outcome: "lost"})
		case "delay":
			if err := m.writer.WriteEvent(artifact.Event{Type: "radio_delayed", Actor: packet.From, Peer: packet.To, Transport: "simulated_lora", Outcome: "delayed"}); err != nil {
				return err
			}
		case "duplicate":
			if err := m.deliver(packet, "duplicate"); err != nil {
				return err
			}
		case "replay":
			if err := m.deliver(packet, "replay"); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown radio effect %q", effect)
		}
	}
	return m.deliver(packet, "delivered")
}

func (m *Medium) deliver(packet Packet, outcome string) error {
	receiver, ok := m.receivers[packet.To]
	if !ok {
		return fmt.Errorf("no receiver for %s", packet.To)
	}
	if err := m.writer.WriteEvent(artifact.Event{
		Type:      "radio_send",
		Actor:     packet.From,
		Peer:      packet.To,
		Direction: "out",
		Transport: "simulated_lora",
		Outcome:   outcome,
		Details:   map[string]any{"bytes": len(packet.Bytes), "label": packet.Label},
	}); err != nil {
		return err
	}
	if err := receiver.ReceiveRadio(packet); err != nil {
		return err
	}
	return m.writer.WriteEvent(artifact.Event{
		Type:      "radio_receive",
		Actor:     packet.To,
		Peer:      packet.From,
		Direction: "in",
		Transport: "simulated_lora",
		Outcome:   outcome,
		Details:   map[string]any{"bytes": len(packet.Bytes), "label": packet.Label},
	})
}
