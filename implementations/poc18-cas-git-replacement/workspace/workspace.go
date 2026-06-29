// Package workspace scans and materializes local filesystem views over POC18
// graph objects.
//
// Intent: Keep storage separate from materialization. A checkout is a local
// promise to create files under local constraints, not global authority over a
// filesystem. Source: DI-radaj; DI-jifuj
package workspace

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"syscall"

	cidlib "github.com/ipfs/go-cid"

	pocchunk "promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/chunk"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

// IngestResult records the top-level CIDs and counts produced by one scan.
type IngestResult struct {
	StoreRoot            string         `json:"store_root"`
	SourceRoot           string         `json:"source_root"`
	CheckoutRoot         string         `json:"checkout_root,omitempty"`
	SnapshotCID          string         `json:"snapshot_cid"`
	RootReferenceSetCID  string         `json:"root_reference_set_cid"`
	WorkspaceRefSetCID   string         `json:"workspace_reference_set_cid"`
	BranchRefSetCID      string         `json:"branch_reference_set_cid"`
	LogicalChangeCID     string         `json:"logical_change_reference_set_cid"`
	ReviewThreadCID      string         `json:"review_thread_reference_set_cid"`
	ReleaseRefSetCID     string         `json:"release_reference_set_cid"`
	DiagnosticMessageCID string         `json:"diagnostic_message_cid"`
	Counts               map[string]int `json:"counts"`
}

// Scanner owns one local scan using one local promiser identity.
type Scanner struct {
	CAS      *store.FileStore
	Promiser string
	Promisee string
	counts   map[string]int
	hardlink map[string]cidlib.Cid
}

// NewScanner returns a scanner with deterministic POC defaults.
func NewScanner(cas *store.FileStore, promiser, promisee string) *Scanner {
	return &Scanner{
		CAS:      cas,
		Promiser: promiser,
		Promisee: promisee,
		counts:   map[string]int{},
		hardlink: map[string]cidlib.Cid{},
	}
}

// Ingest scans root and writes POC18 graph messages into the sparse CAS.
func (scanner *Scanner) Ingest(root string) (IngestResult, error) {
	absRoot, absErr := filepath.Abs(root)
	if absErr != nil {
		return IngestResult{}, absErr
	}
	rootRefCID, _, scanErr := scanner.scanDirectory(absRoot, ".")
	if scanErr != nil {
		return IngestResult{}, scanErr
	}
	snapshot := graph.Payload{
		Promiser:    scanner.Promiser,
		Promisee:    scanner.Promisee,
		PromiseKind: "snapshot",
		PromiseBody: graph.SnapshotBody(
			"snapshot:"+store.CIDText(rootRefCID),
			rootRefCID,
			nil,
			"snapshot from "+absRoot,
			[]any{"materialize only under local constraints"},
		),
	}
	snapshotMessage, snapshotErr := graph.StoreMessage(scanner.CAS, nil, snapshot)
	if snapshotErr != nil {
		return IngestResult{}, snapshotErr
	}
	scanner.counts["snapshot"]++
	workspaceCID, workspaceErr := scanner.storeRoleRefSet("refset:workspace", "workspace", "local:"+absRoot, "root", "snapshot", snapshotMessage.CID)
	if workspaceErr != nil {
		return IngestResult{}, workspaceErr
	}
	branchCID, branchErr := scanner.storeRoleRefSet("refset:main", "branch", "project:poc18-demo", "head", "snapshot", snapshotMessage.CID)
	if branchErr != nil {
		return IngestResult{}, branchErr
	}
	changeCID, changeErr := scanner.storeRoleRefSet("change:first-poc18-snapshot", "logical_change", "project:poc18-demo", "current", "snapshot", snapshotMessage.CID)
	if changeErr != nil {
		return IngestResult{}, changeErr
	}
	reviewCID, reviewErr := scanner.storeRoleRefSet("review:first-poc18-snapshot", "review_thread", "project:poc18-demo", "subject", "snapshot", snapshotMessage.CID)
	if reviewErr != nil {
		return IngestResult{}, reviewErr
	}
	releaseCID, releaseErr := scanner.storeRoleRefSet("release:first-poc18-snapshot", "release", "project:poc18-demo", "source", "snapshot", snapshotMessage.CID)
	if releaseErr != nil {
		return IngestResult{}, releaseErr
	}
	return IngestResult{
		StoreRoot:            scanner.CAS.Root,
		SourceRoot:           absRoot,
		SnapshotCID:          store.CIDText(snapshotMessage.CID),
		RootReferenceSetCID:  store.CIDText(rootRefCID),
		WorkspaceRefSetCID:   store.CIDText(workspaceCID),
		BranchRefSetCID:      store.CIDText(branchCID),
		LogicalChangeCID:     store.CIDText(changeCID),
		ReviewThreadCID:      store.CIDText(reviewCID),
		ReleaseRefSetCID:     store.CIDText(releaseCID),
		DiagnosticMessageCID: store.CIDText(snapshotMessage.CID),
		Counts:               scanner.counts,
	}, nil
}

func (scanner *Scanner) storeRoleRefSet(identity, role, namespace, label, targetRole string, targetCID cidlib.Cid) (cidlib.Cid, error) {
	entries := []any{graph.ReferenceEntry(label, []any{graph.Target(targetRole, targetCID)}, nil)}
	payload := graph.Payload{
		Promiser:    scanner.Promiser,
		Promisee:    scanner.Promisee,
		PromiseKind: "reference_set",
		PromiseBody: graph.ReferenceSetBody(identity, role, namespace, entries, []any{"local promise"}),
	}
	message, storeErr := graph.StoreMessage(scanner.CAS, nil, payload)
	if storeErr != nil {
		return cidlib.Undef, storeErr
	}
	scanner.counts["reference_set:"+role]++
	return message.CID, nil
}

func (scanner *Scanner) scanDirectory(absPath, relPath string) (cidlib.Cid, cidlib.Cid, error) {
	dirEntries, readErr := os.ReadDir(absPath)
	if readErr != nil {
		return cidlib.Undef, cidlib.Undef, readErr
	}
	sort.Slice(dirEntries, func(left, right int) bool {
		return dirEntries[left].Name() < dirEntries[right].Name()
	})
	referenceEntries := []any{}
	for _, dirEntry := range dirEntries {
		childAbs := filepath.Join(absPath, dirEntry.Name())
		childRel := filepath.Join(relPath, dirEntry.Name())
		info, infoErr := os.Lstat(childAbs)
		if infoErr != nil {
			return cidlib.Undef, cidlib.Undef, infoErr
		}
		var childNodeCID cidlib.Cid
		if info.IsDir() {
			_, nodeCID, childErr := scanner.scanDirectory(childAbs, childRel)
			if childErr != nil {
				return cidlib.Undef, cidlib.Undef, childErr
			}
			childNodeCID = nodeCID
		} else {
			nodeCID, childErr := scanner.storeNode(childAbs, childRel, info)
			if childErr != nil {
				return cidlib.Undef, cidlib.Undef, childErr
			}
			childNodeCID = nodeCID
		}
		referenceEntries = append(referenceEntries, graph.ReferenceEntry(dirEntry.Name(), []any{graph.Target("node", childNodeCID)}, nil))
	}
	refPayload := graph.Payload{
		Promiser:    scanner.Promiser,
		Promisee:    scanner.Promisee,
		PromiseKind: "reference_set",
		PromiseBody: graph.ReferenceSetBody("refset:directory:"+relPath, "directory", "workspace:"+relPath, referenceEntries, []any{"directory labels are local dirents"}),
	}
	refMessage, refErr := graph.StoreMessage(scanner.CAS, nil, refPayload)
	if refErr != nil {
		return cidlib.Undef, cidlib.Undef, refErr
	}
	scanner.counts["reference_set:directory"]++
	dirNodePayload := graph.Payload{
		Promiser:    scanner.Promiser,
		Promisee:    scanner.Promisee,
		PromiseKind: "posix_node",
		PromiseBody: graph.PosixNodeBody("node:directory:"+relPath, "directory", store.LinkTag(refMessage.CID), metadataForMode(fs.ModeDir|0o755), []any{"create directory locally if safe"}),
	}
	dirNode, nodeErr := graph.StoreMessage(scanner.CAS, []graph.Parent{{Role: "previous_reference_set", CID: refMessage.CID}}, dirNodePayload)
	if nodeErr != nil {
		return cidlib.Undef, cidlib.Undef, nodeErr
	}
	scanner.counts["posix_node:directory"]++
	return refMessage.CID, dirNode.CID, nil
}

func (scanner *Scanner) storeNode(absPath, relPath string, info os.FileInfo) (cidlib.Cid, error) {
	if key, ok := hardlinkKey(info); ok {
		if existingCID, found := scanner.hardlink[key]; found {
			scanner.counts["hardlink_label"]++
			return existingCID, nil
		}
		nodeCID, nodeErr := scanner.storeFreshNode(absPath, relPath, info, "hardlink:"+key)
		if nodeErr != nil {
			return cidlib.Undef, nodeErr
		}
		scanner.hardlink[key] = nodeCID
		return nodeCID, nil
	}
	return scanner.storeFreshNode(absPath, relPath, info, "node:"+relPath)
}

func (scanner *Scanner) storeFreshNode(absPath, relPath string, info os.FileInfo, nodeIdentity string) (cidlib.Cid, error) {
	nodeType, content, contentErr := scanner.nodeContent(absPath, info)
	if contentErr != nil {
		return cidlib.Undef, contentErr
	}
	payload := graph.Payload{
		Promiser:    scanner.Promiser,
		Promisee:    scanner.Promisee,
		PromiseKind: "posix_node",
		PromiseBody: graph.PosixNodeBody(nodeIdentity, nodeType, content, metadataForInfo(info), []any{"materialize under local safety checks"}),
	}
	message, storeErr := graph.StoreMessage(scanner.CAS, nil, payload)
	if storeErr != nil {
		return cidlib.Undef, storeErr
	}
	scanner.counts["posix_node:"+nodeType]++
	return message.CID, nil
}

func (scanner *Scanner) nodeContent(absPath string, info os.FileInfo) (string, any, error) {
	mode := info.Mode()
	switch {
	case mode.IsRegular():
		storedManifest, chunkErr := pocchunk.StoreFile(scanner.CAS, absPath)
		if chunkErr != nil {
			return "", nil, chunkErr
		}
		chunkRows := make([]any, 0, len(storedManifest.Manifest.Chunks))
		for _, chunkRef := range storedManifest.Manifest.Chunks {
			chunkCID, parseErr := store.ParseCIDText(chunkRef.CID)
			if parseErr != nil {
				return "", nil, parseErr
			}
			chunkRows = append(chunkRows, []any{chunkRef.Offset, chunkRef.Length, store.LinkTag(chunkCID)})
		}
		chunkPayload := graph.Payload{
			Promiser:    scanner.Promiser,
			Promisee:    scanner.Promisee,
			PromiseKind: "chunk_manifest",
			PromiseBody: graph.ChunkManifestBody(storedManifest.ManifestCID, storedManifest.Manifest.FileSize, "rabin", storedManifest.Manifest.Params, chunkRows, ""),
		}
		if _, msgErr := graph.StoreMessage(scanner.CAS, nil, chunkPayload); msgErr != nil {
			return "", nil, msgErr
		}
		scanner.counts["chunk_manifest"]++
		scanner.counts["chunk"] += len(storedManifest.Manifest.Chunks)
		return "regular", store.LinkTag(storedManifest.ManifestCID), nil
	case mode&os.ModeSymlink != 0:
		target, linkErr := os.Readlink(absPath)
		if linkErr != nil {
			return "", nil, linkErr
		}
		return "symlink", []byte(target), nil
	case mode&os.ModeNamedPipe != 0:
		return "fifo", []any{"metadata_only"}, nil
	case mode&os.ModeSocket != 0:
		return "socket", []any{"metadata_only"}, nil
	case mode&os.ModeDevice != 0:
		if mode&os.ModeCharDevice != 0 {
			return "char_device", []any{"metadata_only"}, nil
		}
		return "block_device", []any{"metadata_only"}, nil
	default:
		return "", nil, fmt.Errorf("unsupported file mode %s for %s", mode, absPath)
	}
}

func metadataForInfo(info os.FileInfo) map[string]any {
	metadata := metadataForMode(info.Mode())
	metadata["size"] = info.Size()
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		metadata["dev"] = strconv.FormatUint(uint64(stat.Dev), 10)
		metadata["ino"] = strconv.FormatUint(stat.Ino, 10)
		metadata["nlink"] = uint64(stat.Nlink)
	}
	return metadata
}

func metadataForMode(mode fs.FileMode) map[string]any {
	return map[string]any{
		"mode": mode.String(),
		"perm": fmt.Sprintf("%#o", mode.Perm()),
	}
}

func hardlinkKey(info os.FileInfo) (string, bool) {
	if !info.Mode().IsRegular() {
		return "", false
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok || stat.Nlink < 2 {
		return "", false
	}
	return strconv.FormatUint(uint64(stat.Dev), 10) + ":" + strconv.FormatUint(stat.Ino, 10), true
}

// MaterializeSnapshot checks out snapshotCID into dest using only local CAS
// bytes. Unsupported POSIX node types are recorded in a local note rather than
// silently pretending they were created.
func MaterializeSnapshot(cas *store.FileStore, snapshotCID cidlib.Cid, dest string) error {
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	content, _, getErr := cas.Get(snapshotCID)
	if getErr != nil {
		return getErr
	}
	view, parseErr := graph.ParseEnvelope(content)
	if parseErr != nil {
		return parseErr
	}
	kind, kindErr := view.PayloadKind()
	if kindErr != nil {
		return kindErr
	}
	if kind != "snapshot" {
		return fmt.Errorf("expected snapshot message, got %s", kind)
	}
	body, bodyErr := view.PayloadBody()
	if bodyErr != nil {
		return bodyErr
	}
	if len(body) != 5 {
		return fmt.Errorf("snapshot body must have five slots")
	}
	rootRefCID, cidErr := store.CIDFromLinkTag(body[1])
	if cidErr != nil {
		return cidErr
	}
	materializer := &materializer{cas: cas, seenNodes: map[string]string{}, notes: []map[string]string{}}
	if err := materializer.materializeReferenceSet(rootRefCID, dest); err != nil {
		return err
	}
	return materializer.writeNotes(dest)
}

type materializer struct {
	cas       *store.FileStore
	seenNodes map[string]string
	notes     []map[string]string
}

func (materializer *materializer) materializeReferenceSet(refCID cidlib.Cid, dest string) error {
	content, _, getErr := materializer.cas.Get(refCID)
	if getErr != nil {
		return getErr
	}
	view, parseErr := graph.ParseEnvelope(content)
	if parseErr != nil {
		return parseErr
	}
	kind, kindErr := view.PayloadKind()
	if kindErr != nil {
		return kindErr
	}
	if kind != "reference_set" {
		return fmt.Errorf("expected reference_set, got %s", kind)
	}
	body, bodyErr := view.PayloadBody()
	if bodyErr != nil {
		return bodyErr
	}
	if len(body) != 5 {
		return fmt.Errorf("reference_set body must have five slots")
	}
	entries, ok := body[3].([]any)
	if !ok {
		return fmt.Errorf("reference_set entries must be array")
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}
	for _, entryValue := range entries {
		entry, entryOK := entryValue.([]any)
		if !entryOK || len(entry) != 3 {
			return fmt.Errorf("reference_set entry must have three slots")
		}
		label, labelOK := entry[0].(string)
		if !labelOK {
			return fmt.Errorf("reference_set label must be text")
		}
		targets, targetsOK := entry[1].([]any)
		if !targetsOK || len(targets) == 0 {
			return fmt.Errorf("reference_set target list is empty")
		}
		targetRow, targetOK := targets[0].([]any)
		if !targetOK || len(targetRow) != 2 {
			return fmt.Errorf("reference_set target row must have two slots")
		}
		targetCID, cidErr := store.CIDFromLinkTag(targetRow[1])
		if cidErr != nil {
			return cidErr
		}
		if err := materializer.materializeNode(targetCID, filepath.Join(dest, label)); err != nil {
			return err
		}
	}
	return nil
}

func (materializer *materializer) materializeNode(nodeCID cidlib.Cid, dest string) error {
	content, _, getErr := materializer.cas.Get(nodeCID)
	if getErr != nil {
		return getErr
	}
	view, parseErr := graph.ParseEnvelope(content)
	if parseErr != nil {
		return parseErr
	}
	kind, kindErr := view.PayloadKind()
	if kindErr != nil {
		return kindErr
	}
	if kind != "posix_node" {
		return fmt.Errorf("expected posix_node, got %s", kind)
	}
	nodeText := store.CIDText(nodeCID)
	if firstPath, ok := materializer.seenNodes[nodeText]; ok {
		if err := os.Link(firstPath, dest); err != nil {
			return err
		}
		return nil
	}
	body, bodyErr := view.PayloadBody()
	if bodyErr != nil {
		return bodyErr
	}
	if len(body) != 5 {
		return fmt.Errorf("posix_node body must have five slots")
	}
	nodeType, typeOK := body[1].(string)
	if !typeOK {
		return fmt.Errorf("posix_node type must be text")
	}
	switch nodeType {
	case "directory":
		refCID, cidErr := store.CIDFromLinkTag(body[2])
		if cidErr != nil {
			return cidErr
		}
		if err := materializer.materializeReferenceSet(refCID, dest); err != nil {
			return err
		}
	case "regular":
		manifestCID, cidErr := store.CIDFromLinkTag(body[2])
		if cidErr != nil {
			return cidErr
		}
		manifestBytes, _, getManifestErr := materializer.cas.Get(manifestCID)
		if getManifestErr != nil {
			return getManifestErr
		}
		manifest, decodeErr := pocchunk.DecodeManifest(manifestBytes)
		if decodeErr != nil {
			return decodeErr
		}
		fileBytes, reassembleErr := pocchunk.Reassemble(materializer.cas, manifest)
		if reassembleErr != nil {
			return reassembleErr
		}
		if err := os.WriteFile(dest, fileBytes, 0o644); err != nil {
			return err
		}
		materializer.seenNodes[nodeText] = dest
	case "symlink":
		targetBytes, ok := body[2].([]byte)
		if !ok {
			return fmt.Errorf("symlink content must be bytes")
		}
		if err := os.Symlink(string(targetBytes), dest); err != nil {
			return err
		}
	case "fifo":
		if err := syscall.Mkfifo(dest, 0o644); err != nil {
			return err
		}
	case "char_device", "block_device", "socket":
		materializer.notes = append(materializer.notes, map[string]string{
			"path":      dest,
			"node_type": nodeType,
			"outcome":   "local_non_commitment_metadata_only",
		})
	default:
		return fmt.Errorf("unsupported posix node type %s", nodeType)
	}
	return nil
}

func (materializer *materializer) writeNotes(dest string) error {
	if len(materializer.notes) == 0 {
		return nil
	}
	notePath := filepath.Join(dest, ".poc18-materialization-notes.json")
	content, marshalErr := json.MarshalIndent(materializer.notes, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	content = append(content, '\n')
	return os.WriteFile(notePath, content, 0o644)
}
