package runtime

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/production"
)

const (
	agentCASKindMessage          = "message"
	agentCASKindMalformedMessage = "malformed_message"
	agentCASKindInternal         = "internal"
	agentCASKindEncrypted        = "encrypted"
	agentCASKindPeer             = "peer"
)

// agentCASObject is the local metadata this app keeps for bytes it voluntarily
// stores in its own sparse CAS.
// Intent: Agent-owned CAS state must be separate from the collector-owned raw
// message review store so POC15 can model incomplete per-agent stores, local
// retention promises, and peer storage incentives without creating a global CAS.
// Source: DI-manul
type agentCASObject struct {
	CID               string   `json:"cid"`
	Kind              string   `json:"kind"`
	Owner             string   `json:"owner"`
	SourcePeer        string   `json:"source_peer,omitempty"`
	ProtocolName      string   `json:"protocol_name,omitempty"`
	Retention         string   `json:"retention,omitempty"`
	SizeBytes         int      `json:"size_bytes"`
	Encrypted         bool     `json:"encrypted,omitempty"`
	Pinned            bool     `json:"pinned,omitempty"`
	Paid              bool     `json:"paid,omitempty"`
	ParentCIDs        []string `json:"parent_cids,omitempty"`
	MissingParentCIDs []string `json:"missing_parent_cids,omitempty"`
}

// agentMessageDAGNode is this app's optional local parent index over stored
// message bytes.
// Intent: Message DAG indexes are sparse local views over objects this app has
// actually retained; missing parents are normal local state, not a global
// consistency failure. Source: DI-manul
type agentMessageDAGNode struct {
	CID               string   `json:"cid"`
	ProtocolName      string   `json:"protocol_name,omitempty"`
	ParentCIDs        []string `json:"parent_cids,omitempty"`
	MissingParentCIDs []string `json:"missing_parent_cids,omitempty"`
}

type agentCASStoreOptions struct {
	Kind         string
	SourcePeer   string
	ProtocolName string
	Retention    string
	Encrypted    bool
	Pinned       bool
	Paid         bool
	ParentCIDs   []string
}

// recordAgentCASAccessEvents records the baseline promises every POC15 app makes
// about its own sparse CAS view.
// Intent: Each agent may store local bytes or rely on peer promises, but no
// agent starts with a complete object universe or shared run-store authority.
// Source: DI-manul
func (node *Node) recordAgentCASAccessEvents() {
	node.record("agent_cas_access_promised", "kept", "", "local sparse CAS available to "+node.Agent.Name+" for self-chosen bytes")
	node.record("agent_cas_store_incomplete", "kept", "", "local sparse CAS starts incomplete and may lack peer objects or parents")
	pinnedCID := node.storeLocalCASObject([]byte(fmt.Sprintf("agent=%s\nrun=%s\nkind=pinned-local-state\n", node.Agent.Name, node.Config.RunID)), agentCASStoreOptions{
		Kind:         agentCASKindInternal,
		ProtocolName: "agent_cas_v1",
		Retention:    "pinned-run-local",
		Pinned:       true,
	})
	temporaryCID := node.storeLocalCASObject([]byte(fmt.Sprintf("agent=%s\nrun=%s\nkind=pressure-temporary\n", node.Agent.Name, node.Config.RunID)), agentCASStoreOptions{
		Kind:         agentCASKindInternal,
		ProtocolName: "agent_cas_v1",
		Retention:    "gc-pressure-candidate",
	})
	missingParentCID := production.ContentCID([]byte("missing parent for sparse local DAG|" + node.Agent.Name + "|" + node.Config.RunID))
	node.indexMessageDAGObject(temporaryCID, "agent_cas_v1", []string{missingParentCID})
	ciphertextBytes, cleartextCID, encryptErr := node.encryptLocalCASBytes([]byte(fmt.Sprintf("agent=%s\nrun=%s\nkind=encrypted-local-secret\n", node.Agent.Name, node.Config.RunID)))
	if encryptErr != nil {
		node.record("agent_cas_encrypted_object_store_failed", "broken", "", encryptErr.Error())
		return
	}
	ciphertextCID := node.storeLocalCASObject(ciphertextBytes, agentCASStoreOptions{
		Kind:         agentCASKindEncrypted,
		ProtocolName: "agent_cas_v1",
		Retention:    "encrypted-run-local",
		Encrypted:    true,
		Pinned:       true,
	})
	if ciphertextCID != cleartextCID {
		node.record("agent_cas_ciphertext_cid_selected", "kept", "", "ciphertext_cid="+ciphertextCID+" cleartext_cid="+cleartextCID)
		node.record("agent_cas_cleartext_cid_not_used", "kept", "", "encrypted local object is named by ciphertext CID")
	}
	node.record("agent_cas_local_roots_recorded", "kept", "", "pinned_cid="+pinnedCID+" encrypted_cid="+ciphertextCID)
}

// storeLocalCASObject stores exact bytes in this app's local sparse CAS and
// records the corresponding metadata.
// Intent: The same byte store can contain messages, local app state, encrypted
// blobs, or peer-served data; the pCID or local metadata explains what the bytes
// mean rather than a universal CAS schema. Source: DI-manul
func (node *Node) storeLocalCASObject(objectBytes []byte, options agentCASStoreOptions) string {
	objectCID := production.ContentCID(objectBytes)
	objectKind := options.Kind
	if strings.TrimSpace(objectKind) == "" {
		objectKind = agentCASKindInternal
	}
	objectRecord := agentCASObject{
		CID:          objectCID,
		Kind:         objectKind,
		Owner:        node.Agent.Name,
		SourcePeer:   options.SourcePeer,
		ProtocolName: options.ProtocolName,
		Retention:    options.Retention,
		SizeBytes:    len(objectBytes),
		Encrypted:    options.Encrypted,
		Pinned:       options.Pinned,
		Paid:         options.Paid,
		ParentCIDs:   uniqueStrings(options.ParentCIDs),
	}
	node.mu.Lock()
	if node.casStore == nil {
		node.casStore = make(map[string][]byte)
	}
	if node.agentCASStore == nil {
		node.agentCASStore = make(map[string]agentCASObject)
	}
	node.casStore[objectCID] = append([]byte(nil), objectBytes...)
	node.agentCASStore[objectCID] = objectRecord
	node.mu.Unlock()
	node.record("agent_cas_object_stored", "kept", options.SourcePeer, "cid="+objectCID+" kind="+objectKind+" bytes="+fmt.Sprintf("%d", len(objectBytes)))
	switch objectKind {
	case agentCASKindMessage:
		node.record("agent_cas_message_stored", "kept", options.SourcePeer, "cid="+objectCID+" protocol="+options.ProtocolName)
	case agentCASKindMalformedMessage:
		node.record("agent_cas_malformed_message_stored", "kept", options.SourcePeer, "cid="+objectCID+" protocol="+options.ProtocolName)
	case agentCASKindInternal:
		node.record("agent_cas_internal_object_stored", "kept", options.SourcePeer, "cid="+objectCID+" retention="+options.Retention)
	case agentCASKindEncrypted:
		node.record("agent_cas_encrypted_object_stored", "kept", options.SourcePeer, "cid="+objectCID+" retention="+options.Retention)
	case agentCASKindPeer:
		node.record("agent_cas_peer_object_stored", "kept", options.SourcePeer, "cid="+objectCID+" protocol="+options.ProtocolName)
	}
	return objectCID
}

// recordAgentCASMessageArtifact mirrors exact message bytes into the local sparse
// CAS before the observer-only artifact line is emitted.
// Intent: Operators can still review raw messages through the collector, while
// agents also keep their own local sparse message stores and optional DAG indexes.
// Source: DI-manul
func (node *Node) recordAgentCASMessageArtifact(direction, peer, protocolName string, envelopeBytes []byte, fields map[string]string) {
	messageKind := agentCASKindMessage
	if strings.Contains(direction, "malformed") {
		messageKind = agentCASKindMalformedMessage
	}
	parentCIDs := parentCIDsFromFields(fields)
	messageCID := node.storeLocalCASObject(envelopeBytes, agentCASStoreOptions{
		Kind:         messageKind,
		SourcePeer:   peer,
		ProtocolName: protocolName,
		Retention:    "run-local-message",
		Pinned:       true,
		ParentCIDs:   parentCIDs,
	})
	if messageKind == agentCASKindMessage {
		node.indexMessageDAGObject(messageCID, protocolName, parentCIDs)
	}
}

// indexMessageDAGObject updates this agent's local parent index for one retained
// object.
// Intent: Missing parent CIDs are recorded as normal sparse-DAG state, because
// no agent is expected to have every object seen by every peer. Source: DI-manul
func (node *Node) indexMessageDAGObject(objectCID, protocolName string, parentCIDs []string) agentMessageDAGNode {
	cleanParents := uniqueStrings(parentCIDs)
	node.mu.Lock()
	if node.agentMessageDAG == nil {
		node.agentMessageDAG = make(map[string]agentMessageDAGNode)
	}
	if node.agentCASStore == nil {
		node.agentCASStore = make(map[string]agentCASObject)
	}
	missingParents := make([]string, 0)
	for _, parentCID := range cleanParents {
		if _, ok := node.agentCASStore[parentCID]; !ok {
			missingParents = append(missingParents, parentCID)
		}
	}
	indexNode := agentMessageDAGNode{
		CID:               objectCID,
		ProtocolName:      protocolName,
		ParentCIDs:        cleanParents,
		MissingParentCIDs: missingParents,
	}
	node.agentMessageDAG[objectCID] = indexNode
	if objectRecord, ok := node.agentCASStore[objectCID]; ok {
		objectRecord.MissingParentCIDs = append([]string(nil), missingParents...)
		node.agentCASStore[objectCID] = objectRecord
	}
	node.mu.Unlock()
	node.record("message_dag_node_indexed", "kept", "", "cid="+objectCID+" protocol="+protocolName+" parents="+fmt.Sprintf("%d", len(cleanParents)))
	for _, missingParentCID := range missingParents {
		node.record("message_dag_missing_parent_recorded", "kept", "", "cid="+objectCID+" missing_parent_cid="+missingParentCID)
	}
	if len(cleanParents) > 0 && len(missingParents) == 0 {
		node.record("message_dag_parents_available", "kept", "", "cid="+objectCID+" parent_count="+fmt.Sprintf("%d", len(cleanParents)))
	}
	return indexNode
}

// recordAgentCASGCEvents applies one local run-end GC pass to this app's sparse
// CAS metadata.
// Intent: Retention and removal are promises made by the local agent about local
// bytes; paid, pinned, and encrypted objects are retained while pressure-tagged
// temporary objects may be removed before the run-scoped state is saved. Source:
// DI-manul
func (node *Node) recordAgentCASGCEvents() {
	node.mu.Lock()
	retainedCount := 0
	removedCID := ""
	for objectCID, objectRecord := range node.agentCASStore {
		if objectRecord.Pinned || objectRecord.Paid || objectRecord.Encrypted {
			retainedCount++
			continue
		}
		if removedCID == "" && objectRecord.Kind == agentCASKindInternal && objectRecord.Retention == "gc-pressure-candidate" {
			removedCID = objectCID
			delete(node.agentCASStore, objectCID)
			delete(node.agentMessageDAG, objectCID)
			delete(node.casStore, objectCID)
			continue
		}
		retainedCount++
	}
	totalAfter := len(node.agentCASStore)
	node.mu.Unlock()
	node.record("agent_cas_gc_object_retained", "kept", "", "retained_objects="+fmt.Sprintf("%d", retainedCount)+" after_gc="+fmt.Sprintf("%d", totalAfter))
	if removedCID != "" {
		node.record("agent_cas_gc_object_removed", "kept", "", "cid="+removedCID+" reason=local-pressure-candidate")
	}
}

func (node *Node) encryptLocalCASBytes(cleartextBytes []byte) ([]byte, string, error) {
	cleartextCID := production.ContentCID(cleartextBytes)
	keySeed := sha256.Sum256([]byte("poc15-agent-cas-key|" + node.Agent.Name))
	block, blockErr := aes.NewCipher(keySeed[:])
	if blockErr != nil {
		return nil, "", blockErr
	}
	aead, aeadErr := cipher.NewGCM(block)
	if aeadErr != nil {
		return nil, "", aeadErr
	}
	nonceSeed := sha256.Sum256(append([]byte("poc15-agent-cas-nonce|"+node.Agent.Name+"|"), cleartextBytes...))
	// Intent: The nonce is deterministic only so clean POC15 runs are
	// repeatable; production key management and nonce selection belong to a
	// later identity/encryption pCID. Source: DI-manul
	ciphertextBytes := aead.Seal(nil, nonceSeed[:aead.NonceSize()], cleartextBytes, []byte(node.Agent.Name))
	return ciphertextBytes, cleartextCID, nil
}

func parentCIDsFromFields(fields map[string]string) []string {
	if fields == nil {
		return nil
	}
	parentHashes := []string{
		fields["field_envelope_parent_exact_sha256"],
		fields["field_payload_parent_exact_sha256"],
		fields["field_parent_exact_sha256"],
		fields["parent_exact_sha256"],
	}
	parentCIDs := make([]string, 0, len(parentHashes))
	for _, parentHash := range parentHashes {
		parentCID := exactHashToContentCID(parentHash)
		if parentCID != "" {
			parentCIDs = append(parentCIDs, parentCID)
		}
	}
	return uniqueStrings(parentCIDs)
}

func exactHashToContentCID(exactHash string) string {
	cleanHash := strings.TrimSpace(exactHash)
	if len(cleanHash) != sha256.Size*2 {
		return ""
	}
	if _, decodeErr := hex.DecodeString(cleanHash); decodeErr != nil {
		return ""
	}
	return "cidv1-raw-sha2-256:" + cleanHash
}

func uniqueStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[string]bool, len(values))
	uniqueValues := make([]string, 0, len(values))
	for _, value := range values {
		cleanValue := strings.TrimSpace(value)
		if cleanValue == "" || seen[cleanValue] {
			continue
		}
		seen[cleanValue] = true
		uniqueValues = append(uniqueValues, cleanValue)
	}
	return uniqueValues
}

func copyAgentCASMapOrEmpty(fields map[string]agentCASObject) map[string]agentCASObject {
	copiedFields := make(map[string]agentCASObject, len(fields))
	for key, value := range fields {
		value.ParentCIDs = append([]string(nil), value.ParentCIDs...)
		value.MissingParentCIDs = append([]string(nil), value.MissingParentCIDs...)
		copiedFields[key] = value
	}
	return copiedFields
}

func copyAgentMessageDAGMapOrEmpty(fields map[string]agentMessageDAGNode) map[string]agentMessageDAGNode {
	copiedFields := make(map[string]agentMessageDAGNode, len(fields))
	for key, value := range fields {
		value.ParentCIDs = append([]string(nil), value.ParentCIDs...)
		value.MissingParentCIDs = append([]string(nil), value.MissingParentCIDs...)
		copiedFields[key] = value
	}
	return copiedFields
}
