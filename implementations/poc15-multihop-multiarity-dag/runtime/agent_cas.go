package runtime

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/production"
	"promisegrid.dev/wire-lab/implementations/poc15-multihop-multiarity-dag/protocol"
)

const (
	agentCASKindMessage          = "message"
	agentCASKindMalformedMessage = "malformed_message"
	agentCASKindInternal         = "internal"
	agentCASKindEncrypted        = "encrypted"
	agentCASKindPeer             = "peer"
)

const (
	agentCASStorageProfileGenericBinary  = "generic_binary"
	agentCASStorageProfileTypedExtension = "typed_extension"
	agentCASStorageProfileCBORWrapper    = "cbor_wrapper"

	agentCASWrapperModeOriginalKey = "original_key"
	agentCASWrapperModeWrapperKey  = "wrapper_key"
	agentCASWrapperModeDualKey     = "dual_key"

	agentCASByteFormatBinary      = "binary"
	agentCASByteFormatCBOR        = "cbor"
	agentCASByteFormatCBORWrapper = "cbor_wrapper"
)

// agentCASObject is the local metadata this app keeps for bytes it voluntarily
// stores in its own sparse CAS.
// Intent: Agent-owned CAS state must be separate from the collector-owned raw
// message review store so POC15 can model incomplete per-agent stores, local
// retention promises, and peer storage incentives without creating a global CAS.
// Source: DI-manul; DI-fagog
type agentCASObject struct {
	CID               string   `json:"cid"`
	LogicalCID        string   `json:"logical_cid,omitempty"`
	StoredCID         string   `json:"stored_cid,omitempty"`
	RelativePath      string   `json:"relative_path,omitempty"`
	ByteFormat        string   `json:"byte_format,omitempty"`
	StorageProfile    string   `json:"storage_profile,omitempty"`
	WrapperMode       string   `json:"wrapper_mode,omitempty"`
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
	pinnedCID, pinnedErr := node.storeLocalCASObject([]byte(fmt.Sprintf("agent=%s\nrun=%s\nkind=pinned-local-state\n", node.Agent.Name, node.Config.RunID)), agentCASStoreOptions{
		Kind:         agentCASKindInternal,
		ProtocolName: "agent_cas_v1",
		Retention:    "pinned-run-local",
		Pinned:       true,
	})
	if pinnedErr != nil {
		node.record("agent_cas_internal_object_store_failed", "broken", "", pinnedErr.Error())
		return
	}
	temporaryCID, temporaryErr := node.storeLocalCASObject([]byte(fmt.Sprintf("agent=%s\nrun=%s\nkind=pressure-temporary\n", node.Agent.Name, node.Config.RunID)), agentCASStoreOptions{
		Kind:         agentCASKindInternal,
		ProtocolName: "agent_cas_v1",
		Retention:    "gc-pressure-candidate",
	})
	if temporaryErr != nil {
		node.record("agent_cas_internal_object_store_failed", "broken", "", temporaryErr.Error())
		return
	}
	missingParentCID := production.ContentCID([]byte("missing parent for sparse local DAG|" + node.Agent.Name + "|" + node.Config.RunID))
	node.indexMessageDAGObject(temporaryCID, "agent_cas_v1", []string{missingParentCID})
	ciphertextBytes, cleartextCID, encryptErr := node.encryptLocalCASBytes([]byte(fmt.Sprintf("agent=%s\nrun=%s\nkind=encrypted-local-secret\n", node.Agent.Name, node.Config.RunID)))
	if encryptErr != nil {
		node.record("agent_cas_encrypted_object_store_failed", "broken", "", encryptErr.Error())
		return
	}
	ciphertextCID, ciphertextErr := node.storeLocalCASObject(ciphertextBytes, agentCASStoreOptions{
		Kind:         agentCASKindEncrypted,
		ProtocolName: "agent_cas_v1",
		Retention:    "encrypted-run-local",
		Encrypted:    true,
		Pinned:       true,
	})
	if ciphertextErr != nil {
		node.record("agent_cas_encrypted_object_store_failed", "broken", "", ciphertextErr.Error())
		return
	}
	if ciphertextCID != cleartextCID {
		node.record("agent_cas_ciphertext_cid_selected", "kept", "", "ciphertext_cid="+ciphertextCID+" cleartext_cid="+cleartextCID)
		node.record("agent_cas_cleartext_cid_not_used", "kept", "", "encrypted local object is named by ciphertext CID")
	}
	node.record("agent_cas_local_roots_recorded", "kept", "", "pinned_cid="+pinnedCID+" encrypted_cid="+ciphertextCID)
}

// storeLocalCASObject stores exact bytes in this app's filesystem-backed local
// sparse CAS and records the corresponding metadata.
// Intent: The same byte store can contain messages, local app state, encrypted
// blobs, or peer-served data; the pCID or local metadata explains what the bytes
// mean rather than a universal CAS schema, and DI-fagog keeps durable JSON as an
// index instead of a base64 byte pile. Source: DI-manul; DI-fagog
func (node *Node) storeLocalCASObject(objectBytes []byte, options agentCASStoreOptions) (string, error) {
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
	if err := node.writeLocalCASObjectFile(&objectRecord, objectBytes); err != nil {
		return "", err
	}
	node.mu.Lock()
	if node.agentCASStore == nil {
		node.agentCASStore = make(map[string]agentCASObject)
	}
	node.agentCASStore[objectCID] = objectRecord
	node.mu.Unlock()
	node.record("agent_cas_object_stored", "kept", options.SourcePeer, "cid="+objectCID+" stored_cid="+objectRecord.StoredCID+" kind="+objectKind+" bytes="+fmt.Sprintf("%d", len(objectBytes))+" profile="+objectRecord.StorageProfile+" path="+objectRecord.RelativePath)
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
	return objectCID, nil
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
	messageCID, storeErr := node.storeLocalCASObject(envelopeBytes, agentCASStoreOptions{
		Kind:         messageKind,
		SourcePeer:   peer,
		ProtocolName: protocolName,
		Retention:    "run-local-message",
		Pinned:       true,
		ParentCIDs:   parentCIDs,
	})
	if storeErr != nil {
		node.record("agent_cas_message_store_failed", "broken", peer, "protocol="+protocolName+" "+storeErr.Error())
		return
	}
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
	removedPath := ""
	for objectCID, objectRecord := range node.agentCASStore {
		if objectRecord.Pinned || objectRecord.Paid || objectRecord.Encrypted {
			retainedCount++
			continue
		}
		if removedCID == "" && objectRecord.Kind == agentCASKindInternal && objectRecord.Retention == "gc-pressure-candidate" {
			removedCID = objectCID
			removedPath = node.agentCASObjectPath(objectRecord)
			delete(node.agentCASStore, objectCID)
			delete(node.agentMessageDAG, objectCID)
			continue
		}
		retainedCount++
	}
	totalAfter := len(node.agentCASStore)
	node.mu.Unlock()
	node.record("agent_cas_gc_object_retained", "kept", "", "retained_objects="+fmt.Sprintf("%d", retainedCount)+" after_gc="+fmt.Sprintf("%d", totalAfter))
	if removedCID != "" {
		if removeErr := os.Remove(removedPath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
			node.record("agent_cas_gc_object_remove_failed", "broken", "", "cid="+removedCID+" "+removeErr.Error())
			return
		}
		node.record("agent_cas_gc_object_removed", "kept", "", "cid="+removedCID+" reason=local-pressure-candidate")
	}
}

// agentCASRootDir is the filesystem root for this app's own CAS object files.
// Intent: Agent CAS objects are normal per-agent run files, not observer-volume
// artifacts and not a shared global store. Source: DI-fagog
func (node *Node) agentCASRootDir() string {
	return filepath.Join(node.runDir(), "stores", node.Agent.Name, "cas-objects")
}

func (node *Node) agentCASObjectPath(objectRecord agentCASObject) string {
	return filepath.Join(node.runDir(), "stores", node.Agent.Name, filepath.FromSlash(objectRecord.RelativePath))
}

// safeCASFilenameCID turns a CID string into a portable filename stem.
// Intent: Metadata keeps the exact CID string; filenames only need stable local
// portability across Docker filesystems and developer machines. Source: DI-fagog
func safeCASFilenameCID(contentCID string) string {
	var builder strings.Builder
	for _, value := range contentCID {
		switch {
		case value >= 'a' && value <= 'z':
			builder.WriteRune(value)
		case value >= 'A' && value <= 'Z':
			builder.WriteRune(value)
		case value >= '0' && value <= '9':
			builder.WriteRune(value)
		case value == '-' || value == '_' || value == '.':
			builder.WriteRune(value)
		default:
			builder.WriteByte('_')
		}
	}
	if builder.Len() == 0 {
		return "empty-cid"
	}
	return builder.String()
}

// writeLocalCASObjectFile writes one CAS object with the current agent's
// deterministic storage profile.
// Intent: POC15 deliberately mixes .bin, typed .cbor, and local CBOR wrapper
// files so the CAS contract is about pCID-defined exact bytes, not one universal
// file extension or one complete shared store. Source: DI-fagog
func (node *Node) writeLocalCASObjectFile(objectRecord *agentCASObject, objectBytes []byte) error {
	storageProfile := node.agentCASStorageProfileFor()
	wrapperMode := ""
	byteFormat := agentCASByteFormatBinary
	storedCID := objectRecord.CID
	storedBytes := append([]byte(nil), objectBytes...)
	extension := ".bin"
	if storageProfile == agentCASStorageProfileCBORWrapper {
		wrapperMode = node.agentCASWrapperModeFor()
		wrapperBytes, wrapperErr := marshalAgentCASWrapper(*objectRecord, objectBytes)
		if wrapperErr != nil {
			return wrapperErr
		}
		storedBytes = wrapperBytes
		storedCID = production.ContentCID(wrapperBytes)
		byteFormat = agentCASByteFormatCBORWrapper
		extension = ".cbor"
	} else if storageProfile == agentCASStorageProfileTypedExtension && isCompleteLocalCBORItem(objectBytes) {
		byteFormat = agentCASByteFormatCBOR
		extension = ".cbor"
	}
	relativePath := filepath.ToSlash(filepath.Join("cas-objects", safeCASFilenameCID(storedCID)+extension))
	objectRecord.LogicalCID = objectRecord.CID
	objectRecord.StoredCID = storedCID
	objectRecord.RelativePath = relativePath
	objectRecord.ByteFormat = byteFormat
	objectRecord.StorageProfile = storageProfile
	objectRecord.WrapperMode = wrapperMode
	if err := os.MkdirAll(node.agentCASRootDir(), 0o755); err != nil {
		return err
	}
	finalPath := node.agentCASObjectPath(*objectRecord)
	tempPath := finalPath + fmt.Sprintf(".%d.tmp", os.Getpid())
	if writeErr := os.WriteFile(tempPath, storedBytes, 0o644); writeErr != nil {
		if removeErr := os.Remove(tempPath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
			return fmt.Errorf("write CAS object %s failed, then temp cleanup failed: %w; cleanup: %v", objectRecord.CID, writeErr, removeErr)
		}
		return writeErr
	}
	if renameErr := os.Rename(tempPath, finalPath); renameErr != nil {
		if removeErr := os.Remove(tempPath); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
			return fmt.Errorf("rename CAS object %s failed, then temp cleanup failed: %w; cleanup: %v", objectRecord.CID, renameErr, removeErr)
		}
		return renameErr
	}
	return nil
}

// readLocalCASObject reads an exact object by the original content CID promised
// at the protocol layer, unwrapping local storage wrappers when necessary.
// Intent: Local wrapper files are a storage experiment only; peer CAS promises
// still serve the exact bytes named by the promised content CID. Source: DI-fagog
func (node *Node) readLocalCASObject(contentCID string) ([]byte, bool, error) {
	node.mu.Lock()
	objectRecord, found := node.localCASObjectRecordLocked(contentCID)
	node.mu.Unlock()
	if !found {
		return nil, false, nil
	}
	storedBytes, readErr := os.ReadFile(node.agentCASObjectPath(objectRecord))
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) {
			return nil, false, nil
		}
		return nil, true, readErr
	}
	if objectRecord.ByteFormat == agentCASByteFormatCBORWrapper {
		if contentCID == objectRecord.StoredCID && objectRecord.WrapperMode == agentCASWrapperModeWrapperKey {
			return storedBytes, true, nil
		}
		objectBytes, unwrapErr := unmarshalAgentCASWrapper(storedBytes)
		if unwrapErr != nil {
			return nil, true, unwrapErr
		}
		return objectBytes, true, nil
	}
	return storedBytes, true, nil
}

// localCASObjectExists checks this agent's own sparse CAS metadata and backing
// file without assuming another agent has the same object.
// Intent: Sparse-store probes are local non-commitment checks, not global CAS
// reachability tests. Source: DI-fagog
func (node *Node) localCASObjectExists(contentCID string) bool {
	node.mu.Lock()
	objectRecord, found := node.localCASObjectRecordLocked(contentCID)
	node.mu.Unlock()
	if !found {
		return false
	}
	if _, statErr := os.Stat(node.agentCASObjectPath(objectRecord)); statErr != nil {
		return false
	}
	return true
}

// removeLocalCASObject removes one object from this agent's local CAS metadata
// and filesystem store.
// Intent: GC deletes only bytes this app controls; it does not revoke, command,
// or alter any peer's local CAS. Source: DI-fagog
func (node *Node) removeLocalCASObject(contentCID string) error {
	node.mu.Lock()
	objectRecord, found := node.localCASObjectRecordLocked(contentCID)
	if found {
		delete(node.agentCASStore, objectRecord.CID)
		delete(node.agentMessageDAG, objectRecord.CID)
	}
	node.mu.Unlock()
	if !found {
		return nil
	}
	removeErr := os.Remove(node.agentCASObjectPath(objectRecord))
	if removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
		return removeErr
	}
	return nil
}

func (node *Node) localCASObjectRecordLocked(contentCID string) (agentCASObject, bool) {
	if node.agentCASStore == nil {
		return agentCASObject{}, false
	}
	if objectRecord, ok := node.agentCASStore[contentCID]; ok {
		return objectRecord, true
	}
	for _, objectRecord := range node.agentCASStore {
		if objectRecord.StoredCID == contentCID {
			return objectRecord, true
		}
	}
	return agentCASObject{}, false
}

// agentCASStorageProfileFor deterministically spreads agents across storage
// profiles for each run ID.
// Intent: The POC gets mixed storage behavior without nondeterministic test
// drift, and no agent receives special global CAS privileges. Source: DI-fagog
func (node *Node) agentCASStorageProfileFor() string {
	return node.agentCASStorageProfileForAgent(node.Agent.Name)
}

func (node *Node) agentCASStorageProfileForAgent(agentName string) string {
	profiles := []string{
		agentCASStorageProfileGenericBinary,
		agentCASStorageProfileTypedExtension,
		agentCASStorageProfileCBORWrapper,
	}
	rankedNames := node.agentNamesRankedByRunHash()
	for rankIndex, rankedAgentName := range rankedNames {
		if rankedAgentName == agentName {
			return profiles[rankIndex%len(profiles)]
		}
	}
	return agentCASStorageProfileGenericBinary
}

// agentCASWrapperModeFor deterministically spreads wrapper-profile agents across
// local wrapper CID modes.
// Intent: Wrapper-key and dual-key pressure are tested as local storage metadata
// choices while protocol retrieval remains exact-byte by original CID. Source:
// DI-fagog
func (node *Node) agentCASWrapperModeFor() string {
	wrapperModes := []string{
		agentCASWrapperModeOriginalKey,
		agentCASWrapperModeWrapperKey,
		agentCASWrapperModeDualKey,
	}
	wrapperNames := make([]string, 0)
	for _, agentName := range node.agentNamesRankedByRunHash() {
		if node.agentCASStorageProfileForAgent(agentName) == agentCASStorageProfileCBORWrapper {
			wrapperNames = append(wrapperNames, agentName)
		}
	}
	for wrapperIndex, agentName := range wrapperNames {
		if agentName == node.Agent.Name {
			return wrapperModes[wrapperIndex%len(wrapperModes)]
		}
	}
	return agentCASWrapperModeOriginalKey
}

func (node *Node) agentNamesRankedByRunHash() []string {
	agentNames := make([]string, 0, len(node.Config.Agents))
	for _, agentConfig := range node.Config.Agents {
		agentNames = append(agentNames, agentConfig.Name)
	}
	sort.Slice(agentNames, func(leftIndex, rightIndex int) bool {
		leftDigest := sha256.Sum256([]byte(node.Config.RunID + "|" + agentNames[leftIndex]))
		rightDigest := sha256.Sum256([]byte(node.Config.RunID + "|" + agentNames[rightIndex]))
		if comparison := bytes.Compare(leftDigest[:], rightDigest[:]); comparison != 0 {
			return comparison < 0
		}
		return agentNames[leftIndex] < agentNames[rightIndex]
	})
	return agentNames
}

func marshalAgentCASWrapper(objectRecord agentCASObject, objectBytes []byte) ([]byte, error) {
	buffer := &bytes.Buffer{}
	if err := writeLocalCBORArrayHeader(buffer, 6); err != nil {
		return nil, err
	}
	for _, value := range []string{"agent_cas_wrapper_v1", objectRecord.CID, objectRecord.Kind, objectRecord.ProtocolName, objectRecord.Owner} {
		if err := writeLocalCBORString(buffer, value); err != nil {
			return nil, err
		}
	}
	if err := writeLocalCBORBytes(buffer, objectBytes); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func unmarshalAgentCASWrapper(wrapperBytes []byte) ([]byte, error) {
	reader := &localCBORReader{data: wrapperBytes}
	arrayLength, arrayErr := reader.readTypeAndLength(4)
	if arrayErr != nil {
		return nil, arrayErr
	}
	if arrayLength != 6 {
		return nil, fmt.Errorf("agent CAS wrapper array length = %d, want 6", arrayLength)
	}
	label, labelErr := reader.readString()
	if labelErr != nil {
		return nil, labelErr
	}
	if label != "agent_cas_wrapper_v1" {
		return nil, fmt.Errorf("agent CAS wrapper label = %q", label)
	}
	for index := 0; index < 4; index++ {
		if _, stringErr := reader.readString(); stringErr != nil {
			return nil, stringErr
		}
	}
	objectBytes, bytesErr := reader.readBytes()
	if bytesErr != nil {
		return nil, bytesErr
	}
	if reader.offset != len(reader.data) {
		return nil, fmt.Errorf("trailing agent CAS wrapper bytes: %d", len(reader.data)-reader.offset)
	}
	return objectBytes, nil
}

func isCompleteLocalCBORItem(objectBytes []byte) bool {
	if _, gridErr := protocol.ParseGridMessage(objectBytes); gridErr == nil {
		return true
	}
	reader := &localCBORReader{data: objectBytes}
	if skipErr := reader.skipItem(); skipErr != nil {
		return false
	}
	return reader.offset == len(reader.data)
}

func writeLocalCBORArrayHeader(buffer *bytes.Buffer, length int) error {
	return writeLocalCBORTypeAndLength(buffer, 4, uint64(length))
}

func writeLocalCBORBytes(buffer *bytes.Buffer, value []byte) error {
	if err := writeLocalCBORTypeAndLength(buffer, 2, uint64(len(value))); err != nil {
		return err
	}
	_, writeErr := buffer.Write(value)
	return writeErr
}

func writeLocalCBORString(buffer *bytes.Buffer, value string) error {
	if !utf8.ValidString(value) {
		return fmt.Errorf("invalid utf8 string")
	}
	if err := writeLocalCBORTypeAndLength(buffer, 3, uint64(len(value))); err != nil {
		return err
	}
	_, writeErr := buffer.WriteString(value)
	return writeErr
}

func writeLocalCBORTypeAndLength(buffer *bytes.Buffer, major byte, length uint64) error {
	prefix := major << 5
	switch {
	case length < 24:
		return buffer.WriteByte(prefix | byte(length))
	case length <= 0xff:
		if err := buffer.WriteByte(prefix | 24); err != nil {
			return err
		}
		return buffer.WriteByte(byte(length))
	case length <= 0xffff:
		if err := buffer.WriteByte(prefix | 25); err != nil {
			return err
		}
		return binary.Write(buffer, binary.BigEndian, uint16(length))
	case length <= 0xffffffff:
		if err := buffer.WriteByte(prefix | 26); err != nil {
			return err
		}
		return binary.Write(buffer, binary.BigEndian, uint32(length))
	default:
		if err := buffer.WriteByte(prefix | 27); err != nil {
			return err
		}
		return binary.Write(buffer, binary.BigEndian, length)
	}
}

type localCBORReader struct {
	data   []byte
	offset int
}

func (reader *localCBORReader) readByte() (byte, error) {
	if reader.offset >= len(reader.data) {
		return 0, fmt.Errorf("unexpected end of local cbor data")
	}
	value := reader.data[reader.offset]
	reader.offset++
	return value, nil
}

func (reader *localCBORReader) readTypeAndLength(expectedMajor byte) (uint64, error) {
	initial, err := reader.readByte()
	if err != nil {
		return 0, err
	}
	major := initial >> 5
	if major != expectedMajor {
		return 0, fmt.Errorf("expected cbor major %d, got %d", expectedMajor, major)
	}
	return reader.readAdditionalLength(initial & 0x1f)
}

func (reader *localCBORReader) readAdditionalLength(additional byte) (uint64, error) {
	switch additional {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23:
		return uint64(additional), nil
	case 24:
		value, readErr := reader.readByte()
		return uint64(value), readErr
	case 25:
		if reader.offset+2 > len(reader.data) {
			return 0, fmt.Errorf("truncated uint16")
		}
		value := binary.BigEndian.Uint16(reader.data[reader.offset : reader.offset+2])
		reader.offset += 2
		return uint64(value), nil
	case 26:
		if reader.offset+4 > len(reader.data) {
			return 0, fmt.Errorf("truncated uint32")
		}
		value := binary.BigEndian.Uint32(reader.data[reader.offset : reader.offset+4])
		reader.offset += 4
		return uint64(value), nil
	case 27:
		if reader.offset+8 > len(reader.data) {
			return 0, fmt.Errorf("truncated uint64")
		}
		value := binary.BigEndian.Uint64(reader.data[reader.offset : reader.offset+8])
		reader.offset += 8
		return value, nil
	default:
		return 0, fmt.Errorf("unsupported cbor additional information %d", additional)
	}
}

func (reader *localCBORReader) readString() (string, error) {
	length, err := reader.readTypeAndLength(3)
	if err != nil {
		return "", err
	}
	if reader.offset+int(length) > len(reader.data) {
		return "", fmt.Errorf("truncated string")
	}
	value := string(reader.data[reader.offset : reader.offset+int(length)])
	reader.offset += int(length)
	if !utf8.ValidString(value) {
		return "", fmt.Errorf("invalid utf8 string")
	}
	return value, nil
}

func (reader *localCBORReader) readBytes() ([]byte, error) {
	length, err := reader.readTypeAndLength(2)
	if err != nil {
		return nil, err
	}
	if reader.offset+int(length) > len(reader.data) {
		return nil, fmt.Errorf("truncated byte string")
	}
	value := make([]byte, int(length))
	copy(value, reader.data[reader.offset:reader.offset+int(length)])
	reader.offset += int(length)
	return value, nil
}

func (reader *localCBORReader) skipItem() error {
	initial, err := reader.readByte()
	if err != nil {
		return err
	}
	major := initial >> 5
	additional := initial & 0x1f
	switch major {
	case 0, 1:
		_, lengthErr := reader.readAdditionalLength(additional)
		return lengthErr
	case 2, 3:
		length, lengthErr := reader.readAdditionalLength(additional)
		if lengthErr != nil {
			return lengthErr
		}
		if reader.offset+int(length) > len(reader.data) {
			return fmt.Errorf("truncated cbor item")
		}
		reader.offset += int(length)
		return nil
	case 4:
		length, lengthErr := reader.readAdditionalLength(additional)
		if lengthErr != nil {
			return lengthErr
		}
		for index := uint64(0); index < length; index++ {
			if err := reader.skipItem(); err != nil {
				return err
			}
		}
		return nil
	case 5:
		length, lengthErr := reader.readAdditionalLength(additional)
		if lengthErr != nil {
			return lengthErr
		}
		for index := uint64(0); index < length*2; index++ {
			if err := reader.skipItem(); err != nil {
				return err
			}
		}
		return nil
	case 6:
		if _, lengthErr := reader.readAdditionalLength(additional); lengthErr != nil {
			return lengthErr
		}
		return reader.skipItem()
	case 7:
		switch additional {
		case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23:
			return nil
		case 24:
			if reader.offset+1 > len(reader.data) {
				return fmt.Errorf("truncated cbor simple value")
			}
			reader.offset++
			return nil
		case 25:
			if reader.offset+2 > len(reader.data) {
				return fmt.Errorf("truncated cbor float16")
			}
			reader.offset += 2
			return nil
		case 26:
			if reader.offset+4 > len(reader.data) {
				return fmt.Errorf("truncated cbor float32")
			}
			reader.offset += 4
			return nil
		case 27:
			if reader.offset+8 > len(reader.data) {
				return fmt.Errorf("truncated cbor float64")
			}
			reader.offset += 8
			return nil
		default:
			return fmt.Errorf("unsupported cbor simple additional information %d", additional)
		}
	default:
		return fmt.Errorf("unsupported cbor major type %d", major)
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
