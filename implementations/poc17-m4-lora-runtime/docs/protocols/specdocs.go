package specdocs

import (
	"embed"
	"fmt"
	"sort"

	cidlib "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

const (
	DeviceStatusV1 = "device_status_v1"
	LoRaLinkV1     = "lora_link_v1"
	OrderStatusV1  = "order_status_v1"
	PeerStorageV1  = "peer_storage_v1"
	UnknownProbeV1 = "unknown_probe_v1"
)

const (
	DeviceStatusV1PCID = "bafkreihmlpwkc5ke2zbfvur4mylrikkainqlq4zr5lkwacjys5bgreekhu"
	LoRaLinkV1PCID     = "bafkreifwrfaveh4d2pkizq7pfidflk7cf2gtvfxkcroklizftudiav4r6i"
	OrderStatusV1PCID  = "bafkreidv6zmcfrdopsdpbomlwt75fklmbittdgyztzkaofndw7biorgqne"
	PeerStorageV1PCID  = "bafkreiblomxdawqalqqwdbwh6hlorzdejpkwfmrmt64kxbirohk75vp6dq"
	UnknownProbeV1PCID = "bafkreibtn5h2wf4lz2gbn5tbsjujfxpnadhsoda6x4uq7ok5pnqfuqqdya"
)

// docs embeds the exact POC17 protocol specs used to derive pCIDs.
//
// Intent: POC17 slot 0 must carry actual content-derived protocol CIDs, not
// readable placeholder names, while retaining local names for diagnostics and
// handler dispatch. Source: DI-dutah
//
//go:embed device-status-v1.md lora-link-v1.md order-status-v1.md peer-storage-v1.md unknown-probe-v1.md
var docs embed.FS

var protocolFiles = map[string]string{
	DeviceStatusV1: "device-status-v1.md",
	LoRaLinkV1:     "lora-link-v1.md",
	OrderStatusV1:  "order-status-v1.md",
	PeerStorageV1:  "peer-storage-v1.md",
	UnknownProbeV1: "unknown-probe-v1.md",
}

var protocolPCIDs = map[string]string{
	DeviceStatusV1: DeviceStatusV1PCID,
	LoRaLinkV1:     LoRaLinkV1PCID,
	OrderStatusV1:  OrderStatusV1PCID,
	PeerStorageV1:  PeerStorageV1PCID,
	UnknownProbeV1: UnknownProbeV1PCID,
}

// CID identifies one protocol spec as CIDv1 raw sha2-256.
type CID struct {
	name      string
	text      string
	parsed    cidlib.Cid
	tag42Data []byte
}

// NewCID derives the POC17 protocol CID from exact spec bytes.
func NewCID(name string, specBytes []byte) (CID, error) {
	hash, err := mh.Sum(specBytes, mh.SHA2_256, -1)
	if err != nil {
		return CID{}, fmt.Errorf("multihash spec bytes: %w", err)
	}
	parsed := cidlib.NewCidV1(cidlib.Raw, hash)
	tag42Data := make([]byte, 0, len(parsed.Bytes())+1)
	tag42Data = append(tag42Data, 0x00)
	tag42Data = append(tag42Data, parsed.Bytes()...)
	return CID{name: name, text: parsed.String(), parsed: parsed, tag42Data: tag42Data}, nil
}

// Name returns the local readable protocol name.
func (cid CID) Name() string {
	return cid.name
}

// Bytes returns the binary CIDv1 bytes without the DAG-CBOR tag-42 sentinel.
func (cid CID) Bytes() []byte {
	return append([]byte(nil), cid.parsed.Bytes()...)
}

// Tag42Data returns the DAG-CBOR tag-42 byte-string body for this CID.
func (cid CID) Tag42Data() []byte {
	return append([]byte(nil), cid.tag42Data...)
}

// String returns the canonical CIDv1 base32 text form.
func (cid CID) String() string {
	return cid.text
}

// Registry maps local protocol names to content-derived pCIDs.
type Registry struct {
	byName map[string]CID
	byText map[string]string
}

// NewRegistry derives every POC17 pCID from embedded spec bytes.
func NewRegistry() (Registry, error) {
	registry := Registry{
		byName: make(map[string]CID),
		byText: make(map[string]string),
	}
	for _, name := range Names() {
		specBytes, err := BytesFor(name)
		if err != nil {
			return Registry{}, err
		}
		cid, err := NewCID(name, specBytes)
		if err != nil {
			return Registry{}, err
		}
		if expected := protocolPCIDs[name]; expected != cid.String() {
			return Registry{}, fmt.Errorf("%s pCID drift: got %s want %s", name, cid.String(), expected)
		}
		registry.byName[name] = cid
		registry.byText[cid.String()] = name
	}
	return registry, nil
}

// MustRegistry returns the embedded registry or panics on programming errors.
func MustRegistry() Registry {
	registry, err := NewRegistry()
	if err != nil {
		panic(err)
	}
	return registry
}

// CIDFor returns the content-derived CID for a local protocol name.
func (registry Registry) CIDFor(name string) (CID, bool) {
	cid, ok := registry.byName[name]
	return cid, ok
}

// MustCID returns the content-derived CID for a local protocol name.
func (registry Registry) MustCID(name string) CID {
	cid, ok := registry.CIDFor(name)
	if !ok {
		panic(fmt.Sprintf("unknown POC17 protocol %q", name))
	}
	return cid
}

// NameForCID returns the readable local name for a CID text form.
func (registry Registry) NameForCID(cidText string) (string, bool) {
	name, ok := registry.byText[cidText]
	return name, ok
}

// BytesFor returns exact embedded spec bytes for one local protocol name.
func BytesFor(name string) ([]byte, error) {
	fileName, ok := protocolFiles[name]
	if !ok {
		return nil, fmt.Errorf("unknown POC17 protocol %q", name)
	}
	specBytes, err := docs.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return append([]byte(nil), specBytes...), nil
}

// Names returns protocol names in deterministic order.
func Names() []string {
	names := make([]string, 0, len(protocolFiles))
	for name := range protocolFiles {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
