package pcid

import (
	"fmt"
	"sort"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/specdocs"
)

const (
	RelationshipV1  = "relationship_v1"
	PostalScaleV1   = "postal_scale_v1"
	UPSLabelV1      = "ups_label_v1"
	AccountingV1    = "accounting_v1"
	PrinterPortV1   = "printer_port_v1"
	KernelReceiveV1 = "kernel_receive_v1"
	CASStorageV1    = "cas_storage_v1"
	CIDComputeV1    = "cid_compute_v1"
	IdentityKeyV1   = "identity_key_v1"
	RouteV1         = "route_v1"

	MessageShapeTransportV1       = "message_shape_transport_v1"
	MessageShapeNativeProofV1     = "message_shape_native_proof_v1"
	MessageShapeEnvelopeParentsV1 = "message_shape_envelope_parents_v1"
	MessageShapePayloadParentsV1  = "message_shape_payload_parents_v1"
	MessageShapeCOSEPayloadV1     = "message_shape_cose_payload_v1"
	MessageShapeCOSEProofV1       = "message_shape_cose_proof_v1"

	SecureCapabilityV1  = "secure_capability_v1"
	EncryptedPayloadV1  = "encrypted_payload_v1"
	ParserBuilderRoleV1 = "parser_builder_role_v1"
	MapPayloadProfileV1 = "map_payload_profile_v1"
)

// Registry is the POC16 kernel's local pCID table. It is not a central service
// registry; it is only the local mapping from known protocol-spec names to
// content-derived pCIDs used by this executable experiment.
// Intent: Test slot-0 pCID routing to app receive promises while keeping pCIDs
// as protocol-spec identities, not message-type selectors, and while replacing
// generic report pCIDs with narrower pCID-owned protocols. Source:
// DI-galin; DI-vipih
type Registry struct {
	byName map[string]protocol.ProtocolCID
	byCID  map[string]string
}

// NewRegistry returns the fixed protocol set for this POC. The pCIDs are
// derived from embedded markdown spec bytes so the runtime registry, prompt
// context, and docs/protocols symlink targets can converge on the same spec
// identity.
func NewRegistry() Registry {
	registry := Registry{
		byName: make(map[string]protocol.ProtocolCID),
		byCID:  make(map[string]string),
	}
	for _, name := range []string{
		RelationshipV1,
		PostalScaleV1,
		UPSLabelV1,
		AccountingV1,
		PrinterPortV1,
		KernelReceiveV1,
		CASStorageV1,
		CIDComputeV1,
		IdentityKeyV1,
		RouteV1,
		MessageShapeTransportV1,
		MessageShapeNativeProofV1,
		MessageShapeEnvelopeParentsV1,
		MessageShapePayloadParentsV1,
		MessageShapeCOSEPayloadV1,
		MessageShapeCOSEProofV1,
		SecureCapabilityV1,
		EncryptedPayloadV1,
		ParserBuilderRoleV1,
		MapPayloadProfileV1,
	} {
		specBytes, specErr := specdocs.BytesFor(name)
		if specErr != nil {
			panic(specErr)
		}
		registry.register(name, protocol.NewProtocolCID(specBytes))
	}
	return registry
}

func (registry Registry) register(name string, protocolCID protocol.ProtocolCID) {
	registry.byName[name] = protocolCID
	registry.byCID[protocolCID.String()] = name
}

// CID returns the content-derived pCID for a known protocol name.
func (registry Registry) CID(name string) (protocol.ProtocolCID, bool) {
	protocolCID, ok := registry.byName[name]
	return protocolCID, ok
}

// MustCID returns a known pCID or panics during POC-local programming errors.
func (registry Registry) MustCID(name string) protocol.ProtocolCID {
	protocolCID, ok := registry.CID(name)
	if !ok {
		panic(fmt.Sprintf("unknown POC16 pCID name %q", name))
	}
	return protocolCID
}

// Name returns the local protocol name for a parsed envelope pCID.
func (registry Registry) Name(protocolCID protocol.ProtocolCID) (string, bool) {
	name, ok := registry.byCID[protocolCID.String()]
	return name, ok
}

// Known reports whether a protocol name is in the POC16 registry.
func (registry Registry) Known(name string) bool {
	_, ok := registry.byName[name]
	return ok
}

// Names returns the known protocol names in deterministic order.
func (registry Registry) Names() []string {
	names := make([]string, 0, len(registry.byName))
	for name := range registry.byName {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
