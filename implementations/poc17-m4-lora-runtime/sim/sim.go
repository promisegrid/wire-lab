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
		Name:       "m4-ivan",
		Status:     "idle",
		Battery:    87,
		RetryLimit: cfg.RetryLimit,
		CAS:        state.NewCAS(cfg.LocalCASLimit),
		Writer:     writer,
		Medium:     medium,
		Peer:       "gateway-bob",
	}
	peer := &device.Peer{Name: "gateway-bob", Writer: writer, Medium: medium, Storage: state.NewCAS(2)}
	lifecycle := newLifecycleSupervisor(cfg.RunID, writer)
	if err := lifecycle.start("agent:"+m4.Name, "peer:"+peer.Name); err != nil {
		return err
	}
	medium.Register(m4.Name, m4)
	medium.Register(peer.Name, peer)
	medium.SetReachable(peer.Name, m4.Name, true)
	medium.SetReachable(m4.Name, peer.Name, true)
	// Intent: Resource-limit configuration is analyzer-visible, but POC17 does
	// not report activity usage unless the simulator actually measures it.
	// Source: DI-gidul; DI-rujod
	if err := device.EmitResourceLimitEvidence(writer, m4.Name, device.ResourceLimits{
		RAMBytes:          cfg.RAMByteLimit,
		FlashBytes:        cfg.FlashByteLimit,
		EnergyUnits:       cfg.EnergyUnitLimit,
		RadioAirtimeBytes: cfg.RadioAirtimeByteLimit,
		RetryCount:        uint64(cfg.RetryLimit),
		CASObjects:        uint64(cfg.LocalCASLimit),
	}); err != nil {
		return err
	}
	missingParentBytes := []byte("poc17 peer storage parent fixture")
	missingParentCID, err := protocol.CIDForBytes(missingParentBytes)
	if err != nil {
		return err
	}
	m4.MissingCIDs = append(m4.MissingCIDs, missingParentCID)
	if err := peer.GrantPeerStorage(m4.Name, 96, 2); err != nil {
		return err
	}
	if err := m4.PutPeerStorage(missingParentBytes, "cas_retention_limit"); err != nil {
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
		linkFrame, err := protocol.Build(protocol.Message{ProtocolName: protocol.ProtocolLoRaLink, Payload: []byte(fmt.Sprintf("link-margin-%d", i))})
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
	snapshot := m4.Snapshot()
	if err := writer.WriteEvent(artifact.Event{Type: "device_restart_started", Actor: m4.Name, Outcome: "fresh_agent_requested", Details: map[string]any{"retained_cids": snapshot.RetainedCIDs, "missing_cids": snapshot.MissingCIDs}}); err != nil {
		return err
	}
	m4 = device.NewAgentFromSnapshot(snapshot, cfg.LocalCASLimit, writer, medium)
	medium.Register(m4.Name, m4)
	if err := writer.WriteEvent(artifact.Event{Type: "device_recovery_loaded", Actor: m4.Name, Outcome: "token_and_missing_cids_loaded", Details: map[string]any{"token_cid": m4.Token.CID, "missing_cids": m4.MissingCIDs}}); err != nil {
		return err
	}
	medium.SetReachable(m4.Name, peer.Name, true)
	if err := m4.GetPeerStorage(missingParentCID, "missing_parent_after_restart"); err != nil {
		return err
	}
	if !m4.CAS.Has(missingParentCID) {
		return fmt.Errorf("fresh agent did not recover missing parent %s", missingParentCID)
	}
	if err := writer.WriteEvent(artifact.Event{Type: "device_recovery_verified", Actor: m4.Name, CID: missingParentCID, Outcome: "content_available_after_restart"}); err != nil {
		return err
	}
	return lifecycle.finish()
}
