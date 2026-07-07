package pcid

import (
	"fmt"
	"sort"

	"promisegrid.dev/wire-lab/implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol"
)

const (
	RelationshipV1       = "relationship_v1"
	PostalScaleV1        = "postal_scale_v1"
	UPSLabelV1           = "ups_label_v1"
	AccountingV1         = "accounting_v1"
	PrinterPortV1        = "printer_port_v1"
	KernelReceiveV1      = "kernel_receive_v1"
	KernelTransportV1    = "kernel_transport_v1"
	LocalLifecycleV1     = "local_lifecycle_v1"
	ProductionShippingV1 = "production_shipping_v1"
	CASStorageV1         = "cas_storage_v1"
	CIDComputeV1         = "cid_compute_v1"
	IdentityKeyV1        = "identity_key_v1"
	RouteV1              = "route_v1"

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

type protocolRecord struct {
	name string
	pcid string
}

var protocolRecords = []protocolRecord{
	// accounting_v1; spec: accounting-v1.md
	{name: AccountingV1, pcid: "bafkreiamrh4cnmdpwp343xryjcaqvumbehsprqzdv74xxy4vnnj7cxjtci"},
	// cas_storage_v1; spec: cas-storage-v1.md
	{name: CASStorageV1, pcid: "bafkreibuvp6v3kqi6wdyrysfdppwr4vvgipbyp672t6eentkpj5swosu3y"},
	// cid_compute_v1; spec: cid-compute-v1.md
	{name: CIDComputeV1, pcid: "bafkreihadgz2qbcyqnkdofx4vye5izw6pyieg5jius7c7qguwuaj3oheym"},
	// encrypted_payload_v1; spec: encrypted-payload-v1.md
	{name: EncryptedPayloadV1, pcid: "bafkreid4gz4ccgmi4npdfh4r7kak2nnt4rifh5dxglvtodrlndvgfx5lua"},
	// identity_key_v1; spec: identity-key-v1.md
	{name: IdentityKeyV1, pcid: "bafkreiceftn52uv6adni7t4h6fcm6mpcgrpfwxglby7iu2fbd4tizymanu"},
	// kernel_receive_v1; spec: kernel-receive-v1.md
	{name: KernelReceiveV1, pcid: "bafkreie5ss65n2tovlzgobrt3wummlk5schbbrd5t36igl2ydpwrgktiym"},
	// kernel_transport_v1; spec: kernel-transport-v1.md
	{name: KernelTransportV1, pcid: "bafkreict6shy7n7jrjqkqv56kfjkr3oyzeqxhudp56ioxe6tincj5x3n5i"},
	// local_lifecycle_v1; spec: local-lifecycle-v1.md
	{name: LocalLifecycleV1, pcid: "bafkreidamxalqxl2depjwlzhwdvglpda57fkqy5hvnwiz6jow6tapungeu"},
	// map_payload_profile_v1; spec: map-payload-profile-v1.md
	{name: MapPayloadProfileV1, pcid: "bafkreiccyr5tf5xxmjkypsv6776rrk43iulnf64fhvcsi6skosm5decugq"},
	// message_shape_cose_payload_v1; spec: message-shape-cose-payload-v1.md
	{name: MessageShapeCOSEPayloadV1, pcid: "bafkreigl7pxqdt3cl4qmhl4ymcapjaopeshtjgmyaw2nj2tep2lvgzoauq"},
	// message_shape_cose_proof_v1; spec: message-shape-cose-proof-v1.md
	{name: MessageShapeCOSEProofV1, pcid: "bafkreihehmivruuanooi2qqtb65nrzsyxwtmmu3ougcoq2kxfudtfoe2uy"},
	// message_shape_envelope_parents_v1; spec: message-shape-envelope-parents-v1.md
	{name: MessageShapeEnvelopeParentsV1, pcid: "bafkreibjgum3ipyneq272fzw53udmxkdljlvuuohumpxkbfroljqozdake"},
	// message_shape_native_proof_v1; spec: message-shape-native-proof-v1.md
	{name: MessageShapeNativeProofV1, pcid: "bafkreihha2j7yl4stwq44pnaxyj4reaextg5uovjrb5tm7i5be2ro6asei"},
	// message_shape_payload_parents_v1; spec: message-shape-payload-parents-v1.md
	{name: MessageShapePayloadParentsV1, pcid: "bafkreibnsk5qulr7q2yyfw4heetarlzjdopp6ejt7z5udlbsrrmihewyea"},
	// message_shape_transport_v1; spec: message-shape-transport-v1.md
	{name: MessageShapeTransportV1, pcid: "bafkreif4e5npur2sqcxmos6vndeileztkjjvfdr3lgq2fyr3hwrxzqi6je"},
	// parser_builder_role_v1; spec: parser-builder-role-v1.md
	{name: ParserBuilderRoleV1, pcid: "bafkreibzakr6gkrk2zgdabahdcch2yst7dju3xmn377ho4m25kea5s5t24"},
	// postal_scale_v1; spec: postal-scale-v1.md
	{name: PostalScaleV1, pcid: "bafkreicnivt5fulp3yixamszjntwm3m6ooouo7chgpikqwteti63vqca5a"},
	// printer_port_v1; spec: printer-port-v1.md
	{name: PrinterPortV1, pcid: "bafkreiafpi76326levrz55dsejzsc2cjcr2xsngdz6u4g72jnltlxnqdve"},
	// production_shipping_v1; spec: production-shipping-v1.md
	{name: ProductionShippingV1, pcid: "bafkreiemnu7mezjqs5mdyiabut3lrc4t2z4yd5mpwkxeazrbr3lutk5sui"},
	// relationship_v1; spec: relationship-v1.md
	{name: RelationshipV1, pcid: "bafkreieqq5sjxsrsb64q5chm44rsznxckr2oqk3ax2zo6uiuh4wekj2l64"},
	// route_v1; spec: route-v1.md
	{name: RouteV1, pcid: "bafkreibfb23ds7h45wtwbs3lj4lynyii3gbwfnxdqjtpikzyuw5iwispo4"},
	// secure_capability_v1; spec: secure-capability-v1.md
	{name: SecureCapabilityV1, pcid: "bafkreif6e2tmwg62bjtgljuvido7kk5qcpjn5ryhgh6b6yjfetrqy7z4ve"},
	// ups_label_v1; spec: ups-label-v1.md
	{name: UPSLabelV1, pcid: "bafkreic6cju5snsuak6p6jweku2dicg24jfly2zt7v73fekksdrr4pqunq"},
}

// Registry is the POC16 kernel's local pCID table. It is not a central service
// registry; it is only the local mapping from known protocol-spec names to
// content-derived pCIDs used by this executable experiment.
// Intent: Test slot-0 pCID routing to app receive promises while keeping pCIDs
// as protocol-spec identities, not message-type selectors. The base32 pCID is
// authoritative; readable names are local labels only. Source: DI-galin;
// DI-vipih; DI-sazip
type Registry struct {
	byName map[string]protocol.ProtocolCID
	byCID  map[string]string
}

// NewRegistry returns the fixed protocol set for this POC. The pCIDs are
// hardcoded as canonical CIDv1 base32 strings. Tests verify those strings match
// the embedded markdown spec bytes under the implementation-local docs/protocols
// source of truth, but runtime dispatch does not derive authority from readable
// file names.
// Intent: Keep POC16 protocol specs implementation-local instead of maintaining a
// stale root-level duplicate corpus. Source: DI-magug; DI-sazip
func NewRegistry() Registry {
	registry := Registry{
		byName: make(map[string]protocol.ProtocolCID),
		byCID:  make(map[string]string),
	}
	for _, record := range protocolRecords {
		registry.register(record.name, protocol.MustProtocolCIDText(record.pcid))
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
