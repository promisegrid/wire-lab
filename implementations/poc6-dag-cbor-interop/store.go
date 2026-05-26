package poc6

import (
	"bytes"
	"fmt"

	cid "github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	mh "github.com/multiformats/go-multihash"
)

// EvidenceRecord is a local observation about one actor's promise activity.
// Intent: Keep poc6 aligned with Promise Theory by recording what Alice and Bob
// locally observe instead of treating successful storage or decoding as global
// truth. Source: DI-sagos
type EvidenceRecord struct {
	Actor   string
	Event   string
	Outcome string
	CID     string
	Detail  string
}

// StoredObject keeps the exact DAG-CBOR bytes and their CID together.
// Intent: Preserve exact bytes as the evidence boundary for DAG-CBOR interop,
// because the scenario is about whether CID links, byte strings, and tag-42
// encoding survive without an IPFS daemon. Source: DI-sagos
type StoredObject struct {
	CID   cid.Cid
	Bytes []byte
}

// ObjectStore is a tiny in-memory CAS for the POC.
// Intent: Keep the POC cheap and bounded: prove IPLD/DAG-CBOR byte behavior
// without adding networking, persistence, or kernel/app process shape.
// Source: DI-sagos
type ObjectStore struct {
	objects  map[string]StoredObject
	evidence []EvidenceRecord
}

// NewObjectStore creates an empty in-memory CAS for one bounded test run.
func NewObjectStore() *ObjectStore {
	return &ObjectStore{objects: map[string]StoredObject{}}
}

// StoreObject encodes a node as canonical DAG-CBOR, derives a CIDv1 DAG-CBOR
// CID from the exact bytes, and records the storing actor's local evidence.
func (store *ObjectStore) StoreObject(actor string, node datamodel.Node, promise string) (StoredObject, error) {
	encodedBytes, encodeErr := EncodeDAGCBOR(node)
	if encodeErr != nil {
		store.record(actor, "store", "broken", cid.Undef, promise+": "+encodeErr.Error())
		return StoredObject{}, encodeErr
	}

	objectCID, cidErr := CIDForDAGCBOR(encodedBytes)
	if cidErr != nil {
		store.record(actor, "store", "broken", cid.Undef, promise+": "+cidErr.Error())
		return StoredObject{}, cidErr
	}

	storedObject := StoredObject{CID: objectCID, Bytes: cloneBytes(encodedBytes)}
	store.objects[objectCID.String()] = storedObject
	store.record(actor, "store", "kept", objectCID, promise)
	return storedObject, nil
}

// LoadObject returns the exact bytes and decoded IPLD node for a locally known
// CID, recording whether the asking actor could verify the bytes locally.
func (store *ObjectStore) LoadObject(actor string, objectCID cid.Cid, detail string) (datamodel.Node, []byte, error) {
	storedObject, ok := store.objects[objectCID.String()]
	if !ok {
		err := fmt.Errorf("object %s not found", objectCID)
		store.record(actor, "load", "not-promised", objectCID, detail+": "+err.Error())
		return nil, nil, err
	}

	node, decodeErr := DecodeDAGCBOR(storedObject.Bytes)
	if decodeErr != nil {
		store.record(actor, "load", "broken", objectCID, detail+": "+decodeErr.Error())
		return nil, nil, decodeErr
	}

	store.record(actor, "load", "kept", objectCID, detail)
	return node, cloneBytes(storedObject.Bytes), nil
}

// Evidence returns a copy of the locally recorded evidence.
func (store *ObjectStore) Evidence() []EvidenceRecord {
	evidence := make([]EvidenceRecord, len(store.evidence))
	copy(evidence, store.evidence)
	return evidence
}

func (store *ObjectStore) record(actor string, event string, outcome string, objectCID cid.Cid, detail string) {
	cidText := ""
	if objectCID.Defined() {
		cidText = objectCID.String()
	}
	store.evidence = append(store.evidence, EvidenceRecord{
		Actor:   actor,
		Event:   event,
		Outcome: outcome,
		CID:     cidText,
		Detail:  detail,
	})
}

// BuildMerkleNode creates a small DAG-CBOR-compatible node with byte-string
// content. The map uses string keys because DAG-CBOR requires them.
func BuildMerkleNode(name string, payload []byte) (datamodel.Node, error) {
	builder := basicnode.Prototype.Any.NewBuilder()
	mapAssembler, err := builder.BeginMap(3)
	if err != nil {
		return nil, err
	}
	if err := assignStringEntry(mapAssembler, "kind", "merkle-node"); err != nil {
		return nil, err
	}
	if err := assignStringEntry(mapAssembler, "name", name); err != nil {
		return nil, err
	}
	if err := assignBytesEntry(mapAssembler, "payload", payload); err != nil {
		return nil, err
	}
	if err := mapAssembler.Finish(); err != nil {
		return nil, err
	}
	return builder.Build(), nil
}

// BuildPointerObject creates a node that links to another CID and carries
// byte-string evidence. This is the scenario's "pointer object" shape.
func BuildPointerObject(name string, target cid.Cid, evidence []byte) (datamodel.Node, error) {
	builder := basicnode.Prototype.Any.NewBuilder()
	mapAssembler, err := builder.BeginMap(4)
	if err != nil {
		return nil, err
	}
	if err := assignStringEntry(mapAssembler, "kind", "pointer-object"); err != nil {
		return nil, err
	}
	if err := assignStringEntry(mapAssembler, "name", name); err != nil {
		return nil, err
	}
	if err := assignLinkEntry(mapAssembler, "target", cidlink.Link{Cid: target}); err != nil {
		return nil, err
	}
	if err := assignBytesEntry(mapAssembler, "evidence", evidence); err != nil {
		return nil, err
	}
	if err := mapAssembler.Finish(); err != nil {
		return nil, err
	}
	return builder.Build(), nil
}

// EncodeDAGCBOR serializes the node with DAG-CBOR link support and canonical
// RFC7049 map sorting so repeated encodes produce stable bytes.
func EncodeDAGCBOR(node datamodel.Node) ([]byte, error) {
	var buffer bytes.Buffer
	encodeErr := dagcbor.EncodeOptions{
		AllowLinks:  true,
		MapSortMode: codec.MapSortMode_RFC7049,
	}.Encode(node, &buffer)
	if encodeErr != nil {
		return nil, encodeErr
	}
	return buffer.Bytes(), nil
}

// DecodeDAGCBOR decodes DAG-CBOR bytes using the schema-free IPLD basic node
// representation, including tag-42 CID links.
func DecodeDAGCBOR(encodedBytes []byte) (datamodel.Node, error) {
	builder := basicnode.Prototype.Any.NewBuilder()
	decodeErr := dagcbor.Decode(builder, bytes.NewReader(encodedBytes))
	if decodeErr != nil {
		return nil, decodeErr
	}
	return builder.Build(), nil
}

// CIDForDAGCBOR derives a CIDv1 DAG-CBOR CID using sha2-256 over exact bytes.
func CIDForDAGCBOR(encodedBytes []byte) (cid.Cid, error) {
	return cid.Prefix{
		Version:  1,
		Codec:    cid.DagCBOR,
		MhType:   mh.SHA2_256,
		MhLength: -1,
	}.Sum(encodedBytes)
}

func assignStringEntry(mapAssembler datamodel.MapAssembler, key string, value string) error {
	if err := mapAssembler.AssembleKey().AssignString(key); err != nil {
		return err
	}
	return mapAssembler.AssembleValue().AssignString(value)
}

func assignBytesEntry(mapAssembler datamodel.MapAssembler, key string, value []byte) error {
	if err := mapAssembler.AssembleKey().AssignString(key); err != nil {
		return err
	}
	return mapAssembler.AssembleValue().AssignBytes(cloneBytes(value))
}

func assignLinkEntry(mapAssembler datamodel.MapAssembler, key string, value datamodel.Link) error {
	if err := mapAssembler.AssembleKey().AssignString(key); err != nil {
		return err
	}
	return mapAssembler.AssembleValue().AssignLink(value)
}

func cloneBytes(value []byte) []byte {
	clonedValue := make([]byte, len(value))
	copy(clonedValue, value)
	return clonedValue
}

// VerifyLinkTarget returns the CID in the pointer object's target link.
func VerifyLinkTarget(pointerNode datamodel.Node) (cid.Cid, error) {
	targetNode, lookupErr := pointerNode.LookupByString("target")
	if lookupErr != nil {
		return cid.Undef, lookupErr
	}
	targetLink, linkErr := targetNode.AsLink()
	if linkErr != nil {
		return cid.Undef, linkErr
	}
	cidLink, ok := targetLink.(cidlink.Link)
	if !ok {
		return cid.Undef, fmt.Errorf("target link is %T, not cidlink.Link", targetLink)
	}
	return cidLink.Cid, nil
}

// LookupBytes returns one byte-string field from a decoded IPLD map.
func LookupBytes(node datamodel.Node, key string) ([]byte, error) {
	valueNode, lookupErr := node.LookupByString(key)
	if lookupErr != nil {
		return nil, lookupErr
	}
	valueBytes, bytesErr := valueNode.AsBytes()
	if bytesErr != nil {
		return nil, bytesErr
	}
	return cloneBytes(valueBytes), nil
}
