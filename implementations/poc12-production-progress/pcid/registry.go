package pcid

import (
	"fmt"
	"sort"

	"promisegrid.dev/wire-lab/implementations/poc12-production-progress/protocol"
)

const (
	RelationshipV1  = "relationship_v1"
	PostalScaleV1   = "postal_scale_v1"
	UPSLabelV1      = "ups_label_v1"
	AccountingV1    = "accounting_v1"
	KernelReceiveV1 = "kernel_receive_v1"
)

// Registry is the POC12 kernel's local pCID table. It is not a central service
// registry; it is only the local mapping from known protocol-spec names to
// content-derived pCIDs used by this executable experiment.
// Intent: Test slot-0 pCID routing to app receive promises while keeping pCIDs
// as protocol-spec identities, not message-type selectors. Source: DI-galin
type Registry struct {
	byName map[string]protocol.ProtocolCID
	byCID  map[string]string
}

// NewRegistry returns the fixed protocol set for this POC. The spec bytes are
// concise stand-ins for frozen specs; the resulting pCIDs still behave as
// content-derived protocol identities inside the envelope.
func NewRegistry() Registry {
	registry := Registry{
		byName: make(map[string]protocol.ProtocolCID),
		byCID:  make(map[string]string),
	}
	for _, entry := range []struct {
		name string
		spec string
	}{
		{RelationshipV1, "poc12 relationship trust discovery observation protocol v1"},
		{PostalScaleV1, "poc12 postal scale package weighing protocol v1"},
		{UPSLabelV1, "poc12 ups label printing cost tracking protocol v1"},
		{AccountingV1, "poc12 accounting address lookup shipment update protocol v1"},
		{KernelReceiveV1, "poc12 local app receive promise registration protocol v1"},
	} {
		registry.register(entry.name, protocol.NewProtocolCID([]byte(entry.spec)))
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
		panic(fmt.Sprintf("unknown POC12 pCID name %q", name))
	}
	return protocolCID
}

// Name returns the local protocol name for a parsed envelope pCID.
func (registry Registry) Name(protocolCID protocol.ProtocolCID) (string, bool) {
	name, ok := registry.byCID[protocolCID.String()]
	return name, ok
}

// Known reports whether a protocol name is in the POC12 registry.
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
