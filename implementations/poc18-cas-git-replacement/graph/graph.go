// Package graph builds POC18 version-control promise messages and graph objects.
//
// Intent: Keep pCID-selected message shape, parent links, proofs, and payload
// semantics in one core library used by both `grid` and `poc-*` fixtures. Source:
// DI-harih
package graph

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/fxamacker/cbor/v2"
	cidlib "github.com/ipfs/go-cid"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

const (
	GridTagNumber          = uint64(0x67726964)
	VersionControlPCIDText = "bafkreicrikn3oqfumjnuvruw67h5ffvu6dyy7inz7h2rtm6s4qgwgz7oxu"
)

// Parent names one typed parent edge in the envelope parent slot.
type Parent struct {
	Role string
	CID  cidlib.Cid
}

// Payload is the six-slot pCID-owned promise payload.
type Payload struct {
	Promiser          string
	Promisee          string
	PromiseKind       string
	PromiseBody       any
	ReciprocalPromise any
	LocalConstraints  any
}

// Proof is the first-slice deterministic detached Ed25519 proof.
type Proof struct {
	Profile      string
	Signer       string
	Algorithm    string
	PublicKeyCID cidlib.Cid
	Signature    []byte
	CreatedAt    string
}

// Message is a complete grid([42(pCID), parents, payload, proof]) envelope.
type Message struct {
	ProtocolCID cidlib.Cid
	Parents     []Parent
	Payload     Payload
	Proof       Proof
}

// StoredMessage records one retained graph message.
type StoredMessage struct {
	CID   cidlib.Cid
	Bytes []byte
	Entry store.Entry
}

// VersionControlPCID returns the implementation-local POC18 protocol CID.
func VersionControlPCID() cidlib.Cid {
	parsedCID, parseErr := store.ParseCIDText(VersionControlPCIDText)
	if parseErr != nil {
		panic("version-control pCID constant should parse: " + parseErr.Error())
	}
	return parsedCID
}

// NewMessage signs one pCID-owned promise payload.
func NewMessage(parents []Parent, payload Payload) (Message, error) {
	if payload.Promiser == "" {
		return Message{}, fmt.Errorf("payload promiser is required")
	}
	if payload.PromiseKind == "" {
		return Message{}, fmt.Errorf("payload promise kind is required")
	}
	message := Message{ProtocolCID: VersionControlPCID(), Parents: append([]Parent(nil), parents...), Payload: payload}
	proof, proofErr := sign(message)
	if proofErr != nil {
		return Message{}, proofErr
	}
	message.Proof = proof
	return message, nil
}

// StoreMessage signs and stores one graph message under its message CID.
func StoreMessage(cas *store.FileStore, parents []Parent, payload Payload) (StoredMessage, error) {
	message, messageErr := NewMessage(parents, payload)
	if messageErr != nil {
		return StoredMessage{}, messageErr
	}
	messageBytes, bytesErr := message.Bytes()
	if bytesErr != nil {
		return StoredMessage{}, bytesErr
	}
	entry, putErr := cas.Put("message", messageBytes)
	if putErr != nil {
		return StoredMessage{}, putErr
	}
	messageCID, parseErr := store.ParseCIDText(entry.CID)
	if parseErr != nil {
		return StoredMessage{}, parseErr
	}
	return StoredMessage{CID: messageCID, Bytes: messageBytes, Entry: entry}, nil
}

// Bytes serializes the signed message using the POC18 protocol envelope shape.
func (message Message) Bytes() ([]byte, error) {
	return encodeGrid([]any{
		store.LinkTag(message.ProtocolCID),
		parentArray(message.Parents),
		message.Payload.Array(),
		message.Proof.Array(),
	})
}

// SignableBytes serializes grid([42(pCID), parents, payload]).
func (message Message) SignableBytes() ([]byte, error) {
	return encodeGrid([]any{
		store.LinkTag(message.ProtocolCID),
		parentArray(message.Parents),
		message.Payload.Array(),
	})
}

// Array returns the pCID-owned payload array.
func (payload Payload) Array() []any {
	reciprocal := payload.ReciprocalPromise
	if reciprocal == nil {
		reciprocal = []any{}
	}
	constraints := payload.LocalConstraints
	if constraints == nil {
		constraints = []any{}
	}
	return []any{
		payload.Promiser,
		payload.Promisee,
		payload.PromiseKind,
		payload.PromiseBody,
		reciprocal,
		constraints,
	}
}

// Array returns the pCID-owned proof array.
func (proof Proof) Array() []any {
	return []any{
		proof.Profile,
		proof.Signer,
		proof.Algorithm,
		store.LinkTag(proof.PublicKeyCID),
		proof.Signature,
		proof.CreatedAt,
	}
}

func parentArray(parents []Parent) []any {
	parentValues := make([]any, 0, len(parents))
	for _, parent := range parents {
		parentValues = append(parentValues, []any{parent.Role, store.LinkTag(parent.CID)})
	}
	return parentValues
}

func encodeGrid(slots []any) ([]byte, error) {
	return store.MarshalCBOR(cbor.Tag{Number: GridTagNumber, Content: slots})
}

func sign(message Message) (Proof, error) {
	signableBytes, signableErr := message.SignableBytes()
	if signableErr != nil {
		return Proof{}, signableErr
	}
	privateKey := DeterministicPrivateKey(message.Payload.Promiser)
	publicKey := DeterministicPublicKey(message.Payload.Promiser)
	publicKeyCID := store.CIDForBytes(publicKey)
	return Proof{
		Profile:      "ed25519_detached",
		Signer:       message.Payload.Promiser,
		Algorithm:    "ed25519",
		PublicKeyCID: publicKeyCID,
		Signature:    ed25519.Sign(privateKey, signableBytes),
		CreatedAt:    "poc18-deterministic",
	}, nil
}

// DeterministicPrivateKey is a POC-only reproducible signer.
func DeterministicPrivateKey(seedText string) ed25519.PrivateKey {
	seed := sha256.Sum256([]byte("poc18 protocol signer: " + seedText))
	return ed25519.NewKeyFromSeed(seed[:])
}

// DeterministicPublicKey returns the POC public key for a named promiser.
func DeterministicPublicKey(seedText string) ed25519.PublicKey {
	return DeterministicPrivateKey(seedText).Public().(ed25519.PublicKey)
}

// ChunkManifestBody returns the promise body for a chunk_manifest message.
func ChunkManifestBody(manifestCID cidlib.Cid, fileSize int64, chunkerName string, parameters any, chunks []any, contentDigest string) []any {
	return []any{store.LinkTag(manifestCID), fileSize, chunkerName, parameters, chunks, contentDigest}
}

// PosixNodeBody returns the promise body for a posix_node message.
func PosixNodeBody(nodeIdentity, nodeType string, content any, metadata any, materialization any) []any {
	return []any{nodeIdentity, nodeType, content, metadata, materialization}
}

// Target returns a reference-set target pair.
func Target(role string, targetCID cidlib.Cid) []any {
	return []any{role, store.LinkTag(targetCID)}
}

// ReferenceEntry returns a reference-set entry row.
func ReferenceEntry(label string, targets []any, terms any) []any {
	if terms == nil {
		terms = []any{}
	}
	return []any{label, targets, terms}
}

// ReferenceSetBody returns the promise body for a reference_set message.
func ReferenceSetBody(identity, role, namespace string, entries []any, promisedTerms any) []any {
	if promisedTerms == nil {
		promisedTerms = []any{}
	}
	return []any{identity, role, namespace, entries, promisedTerms}
}

// SnapshotBody returns the promise body for a snapshot message.
func SnapshotBody(identity string, rootDirectoryCID cidlib.Cid, parentSnapshots []cidlib.Cid, summary string, terms any) []any {
	parents := make([]any, 0, len(parentSnapshots))
	for _, parentCID := range parentSnapshots {
		parents = append(parents, store.LinkTag(parentCID))
	}
	if terms == nil {
		terms = []any{}
	}
	return []any{identity, store.LinkTag(rootDirectoryCID), parents, summary, terms}
}

// ObjectRow returns one availability or interest row naming a CID.
func ObjectRow(role string, objectCID cidlib.Cid, extra ...any) []any {
	row := []any{role, store.LinkTag(objectCID)}
	return append(row, extra...)
}

// ObjectAvailabilityBody returns the promise body for object_availability.
//
// Intent: Keep peer retrieval as an explicit promise that selected object CIDs
// are locally available, serveable, forwardable, missing, or not promised.
// Source: DI-gozov
func ObjectAvailabilityBody(scope string, objects []any, serviceTerms any) []any {
	if serviceTerms == nil {
		serviceTerms = []any{}
	}
	return []any{scope, objects, serviceTerms}
}

// SyncInterestBody returns the promise body for sync_interest.
//
// Intent: Keep missing-object retrieval as Bob's voluntary promise to receive
// selected CIDs under local constraints and reciprocal terms. Source: DI-gozov
func SyncInterestBody(scope string, wantedObjects []any, offerTerms any, refusalTerms any) []any {
	if offerTerms == nil {
		offerTerms = []any{}
	}
	if refusalTerms == nil {
		refusalTerms = []any{}
	}
	return []any{scope, wantedObjects, offerTerms, refusalTerms}
}

// EnvelopeView is a parsed POC18 message shape used by checkout and diagnostics.
type EnvelopeView struct {
	ProtocolCID cidlib.Cid
	Parents     []Parent
	Payload     []any
	Proof       []any
}

// ParseEnvelope parses the POC18 envelope enough for local graph traversal.
func ParseEnvelope(content []byte) (EnvelopeView, error) {
	var raw any
	if err := store.UnmarshalCBOR(content, &raw); err != nil {
		return EnvelopeView{}, err
	}
	tag, ok := raw.(cbor.Tag)
	if !ok || tag.Number != GridTagNumber {
		return EnvelopeView{}, fmt.Errorf("message must be grid tag")
	}
	slots, ok := tag.Content.([]any)
	if !ok || len(slots) != 4 {
		return EnvelopeView{}, fmt.Errorf("grid message must have four slots")
	}
	protocolCID, cidErr := store.CIDFromLinkTag(slots[0])
	if cidErr != nil {
		return EnvelopeView{}, cidErr
	}
	parents, parentErr := parseParents(slots[1])
	if parentErr != nil {
		return EnvelopeView{}, parentErr
	}
	payload, ok := slots[2].([]any)
	if !ok || len(payload) != 6 {
		return EnvelopeView{}, fmt.Errorf("payload must be six-slot array")
	}
	proof, ok := slots[3].([]any)
	if !ok || len(proof) != 6 {
		return EnvelopeView{}, fmt.Errorf("proof must be six-slot array")
	}
	return EnvelopeView{ProtocolCID: protocolCID, Parents: parents, Payload: payload, Proof: proof}, nil
}

func parseParents(value any) ([]Parent, error) {
	rows, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("parents must be array")
	}
	parents := make([]Parent, 0, len(rows))
	for _, rowValue := range rows {
		row, rowOK := rowValue.([]any)
		if !rowOK || len(row) != 2 {
			return nil, fmt.Errorf("parent row must have two slots")
		}
		role, roleOK := row[0].(string)
		if !roleOK {
			return nil, fmt.Errorf("parent role must be text")
		}
		parentCID, cidErr := store.CIDFromLinkTag(row[1])
		if cidErr != nil {
			return nil, cidErr
		}
		parents = append(parents, Parent{Role: role, CID: parentCID})
	}
	return parents, nil
}

// PayloadKind returns the promise_kind slot from a parsed envelope.
func (view EnvelopeView) PayloadKind() (string, error) {
	kind, ok := view.Payload[2].(string)
	if !ok {
		return "", fmt.Errorf("payload kind must be text")
	}
	return kind, nil
}

// PayloadBody returns the promise_body slot from a parsed envelope.
func (view EnvelopeView) PayloadBody() ([]any, error) {
	body, ok := view.Payload[3].([]any)
	if !ok {
		return nil, fmt.Errorf("payload body must be array")
	}
	return body, nil
}

// Diagnostic renders a concise human-readable view of exact CBOR message bytes.
func Diagnostic(content []byte) (string, error) {
	var raw any
	if err := store.UnmarshalCBOR(content, &raw); err != nil {
		return "", err
	}
	var builder strings.Builder
	renderValue(&builder, raw, 0)
	builder.WriteByte('\n')
	return builder.String(), nil
}

func renderValue(builder *strings.Builder, value any, indent int) {
	prefix := strings.Repeat("  ", indent)
	switch typed := value.(type) {
	case cbor.Tag:
		if typed.Number == GridTagNumber {
			builder.WriteString("grid(")
			renderValue(builder, typed.Content, indent)
			builder.WriteString(")")
			return
		}
		if typed.Number == store.LinkTagNumber {
			if linkedCID, err := store.CIDFromLinkTag(typed); err == nil {
				builder.WriteString("42(")
				builder.WriteString(store.CIDText(linkedCID))
				builder.WriteString(")")
				return
			}
		}
		builder.WriteString(fmt.Sprintf("tag%d(", typed.Number))
		renderValue(builder, typed.Content, indent)
		builder.WriteString(")")
	case []any:
		builder.WriteString("[")
		if len(typed) > 0 {
			builder.WriteByte('\n')
		}
		for index, item := range typed {
			builder.WriteString(prefix)
			builder.WriteString("  ")
			renderValue(builder, item, indent+1)
			if index != len(typed)-1 {
				builder.WriteByte(',')
			}
			builder.WriteByte('\n')
		}
		if len(typed) > 0 {
			builder.WriteString(prefix)
		}
		builder.WriteString("]")
	case map[any]any:
		keys := make([]string, 0, len(typed))
		rendered := map[string]any{}
		for key, item := range typed {
			keyText := fmt.Sprint(key)
			keys = append(keys, keyText)
			rendered[keyText] = item
		}
		sort.Strings(keys)
		builder.WriteString("{")
		if len(keys) > 0 {
			builder.WriteByte('\n')
		}
		for index, key := range keys {
			builder.WriteString(prefix)
			builder.WriteString("  ")
			builder.WriteString(fmt.Sprintf("%q: ", key))
			renderValue(builder, rendered[key], indent+1)
			if index != len(keys)-1 {
				builder.WriteByte(',')
			}
			builder.WriteByte('\n')
		}
		if len(keys) > 0 {
			builder.WriteString(prefix)
		}
		builder.WriteString("}")
	case []byte:
		if len(typed) <= 16 {
			builder.WriteString("h'")
			builder.WriteString(hex.EncodeToString(typed))
			builder.WriteString("'")
			return
		}
		builder.WriteString("h'")
		builder.WriteString(hex.EncodeToString(typed[:16]))
		builder.WriteString("...'")
	default:
		builder.WriteString(fmt.Sprintf("%q", typed))
	}
}
