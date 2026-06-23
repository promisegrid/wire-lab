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
	return p.Writer.WriteEvent(artifact.Event{Type: "peer_envelope_received", Actor: p.Name, Peer: packet.From, PCID: msg.PCID, Hash: hash, Path: rel, Transport: "simulated_lora", Outcome: "received"})
}
