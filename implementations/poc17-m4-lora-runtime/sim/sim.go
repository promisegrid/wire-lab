package sim

import (
	"fmt"
	"os"

	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/artifact"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/config"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/device"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/protocol"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/radio"
	"promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/state"
)

// Run executes the deterministic first POC17 behavior scenario.
func Run(cfg config.Config) error {
	if err := os.MkdirAll(cfg.RunDir(), 0o755); err != nil {
		return fmt.Errorf("create run dir: %w", err)
	}
	writer, err := artifact.NewWriter(cfg.RunDir())
	if err != nil {
		return err
	}
	if err := writer.WriteEvent(artifact.Event{
		Type:    "simulator_fidelity_notice",
		Actor:   "harness",
		Outcome: "behavior_evidence_only",
		Details: map[string]any{
			"not_firmware_proof": true,
			"source":             "DI-libis",
		},
	}); err != nil {
		return err
	}
	medium := radio.NewMedium(cfg.RadioMTUBytes, writer)
	m4 := &device.Agent{
		Name:        "m4-ivan",
		Status:      "idle",
		Battery:     87,
		RetryBudget: cfg.RetryBudget,
		CAS:         state.NewCAS(cfg.LocalCASLimit),
		Writer:      writer,
		Medium:      medium,
		Peer:        "gateway-bob",
	}
	peer := &device.Peer{Name: "gateway-bob", Writer: writer, Medium: medium}
	lifecycle := newLifecycleSupervisor(cfg.RunID, writer)
	if err := lifecycle.start("agent:"+m4.Name, "peer:"+peer.Name); err != nil {
		return err
	}
	medium.Register(m4.Name, m4)
	medium.Register(peer.Name, peer)
	medium.SetReachable(peer.Name, m4.Name, true)
	medium.SetReachable(m4.Name, peer.Name, true)
	missingParentCID, err := protocol.CIDForBytes([]byte("poc17 missing parent fixture"))
	if err != nil {
		return err
	}

	// Intent: Use bintags' order-number/status workflow as production-like traffic while retaining PromiseGrid CBOR envelopes. Source: DI-mokit
	orderPayload, err := protocol.BuildOrderStatusPayload(protocol.OrderStatusPayload{
		Type:        "MSG",
		Source:      peer.Name,
		Dest:        m4.Name,
		Counter:     1,
		OrderNumber: "BT-1042",
		Status:      "created",
	})
	if err != nil {
		return err
	}
	orderFrame, err := protocol.Build(protocol.Message{ProtocolName: protocol.ProtocolOrderStatus, Payload: orderPayload})
	if err != nil {
		return err
	}
	if err := medium.Send(radio.Packet{From: peer.Name, To: m4.Name, Bytes: orderFrame, Label: "order-reset-bt-1042"}, "duplicate"); err != nil {
		return err
	}
	for _, status := range []string{"cut", "stripped", "soldered", "completed"} {
		if err := m4.PromiseOrderStatus(status); err != nil {
			return err
		}
	}
	statusPayload, err := protocol.BuildStatusPayload(m4.Name, "ready", 87, []string{missingParentCID})
	if err != nil {
		return err
	}
	statusFrame, err := protocol.Build(protocol.Message{ProtocolName: protocol.ProtocolDeviceStatus, Payload: statusPayload})
	if err != nil {
		return err
	}
	if err := medium.Send(radio.Packet{From: peer.Name, To: m4.Name, Bytes: []byte{0xff, 0x01, 0x02}, Label: "malformed"}); err != nil {
		return err
	}
	unknownFrame, err := protocol.Build(protocol.Message{ProtocolName: protocol.ProtocolUnknownProbe, Payload: []byte("probe")})
	if err != nil {
		return err
	}
	if err := medium.Send(radio.Packet{From: peer.Name, To: m4.Name, Bytes: unknownFrame, Label: "unknown-pcid"}, "replay"); err != nil {
		return err
	}
	for i := 1; i <= 4; i++ {
		linkFrame, err := protocol.Build(protocol.Message{ProtocolName: protocol.ProtocolLoRaLink, Payload: []byte(fmt.Sprintf("link-budget-%d", i))})
		if err != nil {
			return err
		}
		if err := medium.Send(radio.Packet{From: peer.Name, To: m4.Name, Bytes: linkFrame, Label: fmt.Sprintf("link-%d", i)}); err != nil {
			return err
		}
	}
	oversized := make([]byte, cfg.RadioMTUBytes+1)
	if err := medium.Send(radio.Packet{From: peer.Name, To: m4.Name, Bytes: oversized, Label: "oversized"}); err != nil {
		return err
	}
	if err := medium.Send(radio.Packet{From: peer.Name, To: m4.Name, Bytes: statusFrame, Label: "lost"}, "loss"); err != nil {
		return err
	}
	if err := medium.Send(radio.Packet{From: peer.Name, To: m4.Name, Bytes: statusFrame, Label: "delayed"}, "delay"); err != nil {
		return err
	}
	medium.SetReachable(m4.Name, peer.Name, false)
	if err := medium.Send(radio.Packet{From: m4.Name, To: peer.Name, Bytes: statusFrame, Label: "asymmetric"}); err != nil {
		return err
	}
	if err := m4.PromisePeerStorage(missingParentCID); err != nil {
		return err
	}
	return lifecycle.finish()
}
