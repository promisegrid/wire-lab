package device

import "promisegrid.dev/wire-lab/implementations/poc17-m4-lora-runtime/artifact"

// ResourceLimits captures simulator-visible limits without pretending to model
// exact M4 hardware counters.
type ResourceLimits struct {
	RAMBytes          uint64
	FlashBytes        uint64
	EnergyUnits       uint64
	RadioAirtimeBytes uint64
	RetryCount        uint64
	CASObjects        uint64
}

// EmitResourceLimitEvidence writes configured limits only; activity counters
// must come from real simulator measurements before they are reported.
func EmitResourceLimitEvidence(writer *artifact.Writer, actor string, limits ResourceLimits) error {
	// Intent: Show resource-limit configuration without fabricating activity
	// usage values. Source: DI-gidul; DI-rujod
	snapshot := map[string]any{
		"ram_byte_limit":           limits.RAMBytes,
		"flash_byte_limit":         limits.FlashBytes,
		"energy_unit_limit":        limits.EnergyUnits,
		"radio_airtime_byte_limit": limits.RadioAirtimeBytes,
		"retry_limit":              limits.RetryCount,
		"cas_object_limit":         limits.CASObjects,
	}
	return writer.WriteEvent(artifact.Event{Type: "resource_limit_snapshot", Actor: actor, Outcome: "configured", Details: snapshot})
}
