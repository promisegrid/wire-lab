package workspace

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	cidlib "github.com/ipfs/go-cid"

	pocchunk "promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/chunk"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
)

const (
	StatusClean       = "clean"
	StatusModified    = "modified"
	StatusMissing     = "missing"
	StatusUntracked   = "untracked"
	StatusTypeChanged = "type_changed"
	// StatusTrackingAdded means local policy currently tracks a path absent from
	// the latest recorded snapshot.
	StatusTrackingAdded = "tracking_added"
	// StatusTrackingRemoved means local policy currently excludes a path present
	// in the latest recorded snapshot.
	StatusTrackingRemoved = "tracking_removed"
)

// StatusReport summarizes a read-only comparison between a workspace and a
// retained snapshot.
type StatusReport struct {
	SourceRoot        string        `json:"source_root"`
	SnapshotCID       string        `json:"snapshot_cid"`
	Clean             bool          `json:"clean"`
	ContentDiff       bool          `json:"content_diff"`
	TrackedStatusDiff bool          `json:"tracked_status_diff"`
	Entries           []StatusEntry `json:"entries"`
}

// StatusEntry names one path whose current local state differs from a snapshot.
type StatusEntry struct {
	Path              string `json:"path"`
	Status            string `json:"status"`
	Type              string `json:"type,omitempty"`
	ContentDiff       bool   `json:"content_diff"`
	TrackedStatusDiff bool   `json:"tracked_status_diff"`
	Detail            string `json:"detail,omitempty"`
}

type expectedNode struct {
	nodeType string
	content  []byte
}

type currentNode struct {
	nodeType string
	content  []byte
}

// CompareSnapshot compares root against snapshotCID without writing workspace
// bytes, CAS objects, or repo state.
//
// Intent: `grid status` must answer from retained snapshot promises and the
// current filesystem without creating a new promise object or mutating the
// local CAS. Source: DI-bikif
func CompareSnapshot(cas *store.FileStore, snapshotCID cidlib.Cid, root string) (StatusReport, error) {
	return CompareSnapshotWithExcludedPaths(cas, snapshotCID, root, nil)
}

// CompareSnapshotWithExcludedPaths compares a workspace while applying
// repo-relative paths that local state excludes.
//
// Intent: `grid status` reports content drift separately from local tracking
// policy drift. `track` and `untrack` therefore become visible immediately
// without mutating the CAS or waiting for a new snapshot. Source: DI-jokav;
// DI-tuhoj
func CompareSnapshotWithExcludedPaths(cas *store.FileStore, snapshotCID cidlib.Cid, root string, excludedPaths []string) (StatusReport, error) {
	absRoot, absErr := filepath.Abs(root)
	if absErr != nil {
		return StatusReport{}, absErr
	}
	expected, expectedErr := loadExpectedSnapshot(cas, snapshotCID)
	if expectedErr != nil {
		return StatusReport{}, expectedErr
	}
	exclusions := newPathExclusionSet(excludedPaths)
	current, currentErr := scanCurrentWorkspace(absRoot)
	if currentErr != nil {
		return StatusReport{}, currentErr
	}
	entries := compareNodes(expected, current, exclusions)
	contentDiff, trackedStatusDiff := summarizeStatusEntries(entries)
	return StatusReport{
		SourceRoot:        absRoot,
		SnapshotCID:       store.CIDText(snapshotCID),
		Clean:             len(entries) == 0,
		ContentDiff:       contentDiff,
		TrackedStatusDiff: trackedStatusDiff,
		Entries:           entries,
	}, nil
}

func loadExpectedSnapshot(cas *store.FileStore, snapshotCID cidlib.Cid) (map[string]expectedNode, error) {
	content, _, getErr := cas.Get(snapshotCID)
	if getErr != nil {
		return nil, getErr
	}
	view, parseErr := graph.ParseEnvelope(content)
	if parseErr != nil {
		return nil, parseErr
	}
	kind, kindErr := view.PayloadKind()
	if kindErr != nil {
		return nil, kindErr
	}
	if kind != "snapshot" {
		return nil, fmt.Errorf("expected snapshot message, got %s", kind)
	}
	body, bodyErr := view.PayloadBody()
	if bodyErr != nil {
		return nil, bodyErr
	}
	if len(body) != 5 {
		return nil, fmt.Errorf("snapshot body must have five slots")
	}
	rootRefCID, cidErr := store.CIDFromLinkTag(body[1])
	if cidErr != nil {
		return nil, cidErr
	}
	expected := map[string]expectedNode{}
	if walkErr := loadExpectedReferenceSet(cas, rootRefCID, "", expected); walkErr != nil {
		return nil, walkErr
	}
	return expected, nil
}

func loadExpectedReferenceSet(cas *store.FileStore, refCID cidlib.Cid, relDir string, expected map[string]expectedNode) error {
	content, _, getErr := cas.Get(refCID)
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
		relPath := joinStatusPath(relDir, label)
		if nodeErr := loadExpectedNode(cas, targetCID, relPath, expected); nodeErr != nil {
			return nodeErr
		}
	}
	return nil
}

func loadExpectedNode(cas *store.FileStore, nodeCID cidlib.Cid, relPath string, expected map[string]expectedNode) error {
	content, _, getErr := cas.Get(nodeCID)
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
		expected[relPath] = expectedNode{nodeType: nodeType}
		refCID, cidErr := store.CIDFromLinkTag(body[2])
		if cidErr != nil {
			return cidErr
		}
		return loadExpectedReferenceSet(cas, refCID, relPath, expected)
	case "regular":
		manifestCID, cidErr := store.CIDFromLinkTag(body[2])
		if cidErr != nil {
			return cidErr
		}
		manifestBytes, _, manifestErr := cas.Get(manifestCID)
		if manifestErr != nil {
			return manifestErr
		}
		manifest, decodeErr := pocchunk.DecodeManifest(manifestBytes)
		if decodeErr != nil {
			return decodeErr
		}
		fileBytes, reassembleErr := pocchunk.Reassemble(cas, manifest)
		if reassembleErr != nil {
			return reassembleErr
		}
		expected[relPath] = expectedNode{nodeType: nodeType, content: fileBytes}
	case "symlink":
		targetBytes, ok := body[2].([]byte)
		if !ok {
			return fmt.Errorf("symlink content must be bytes")
		}
		expected[relPath] = expectedNode{nodeType: nodeType, content: targetBytes}
	case "fifo", "char_device", "block_device", "socket":
		expected[relPath] = expectedNode{nodeType: nodeType}
	default:
		return fmt.Errorf("unsupported posix node type %s", nodeType)
	}
	return nil
}

func scanCurrentWorkspace(root string) (map[string]currentNode, error) {
	current := map[string]currentNode{}
	walkErr := filepath.WalkDir(root, func(path string, dirEntry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root {
			return nil
		}
		relPath, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		relPath = filepath.ToSlash(relPath)
		if relPath == ".grid" || strings.HasPrefix(relPath, ".grid/") {
			if dirEntry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		info, infoErr := dirEntry.Info()
		if infoErr != nil {
			return infoErr
		}
		nodeType := statusTypeForInfo(info)
		node := currentNode{nodeType: nodeType}
		switch nodeType {
		case "regular":
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr
			}
			node.content = content
		case "symlink":
			target, readErr := os.Readlink(path)
			if readErr != nil {
				return readErr
			}
			node.content = []byte(target)
		}
		current[relPath] = node
		return nil
	})
	return current, walkErr
}

func compareNodes(expected map[string]expectedNode, current map[string]currentNode, exclusions pathExclusionSet) []StatusEntry {
	entries := []StatusEntry{}
	for relPath, expectedNode := range expected {
		if exclusions.excludes(relPath) {
			continue
		}
		currentNode, found := current[relPath]
		if !found {
			entries = append(entries, StatusEntry{Path: relPath, Status: StatusMissing, Type: expectedNode.nodeType, ContentDiff: true})
			continue
		}
		if currentNode.nodeType != expectedNode.nodeType {
			entries = append(entries, StatusEntry{Path: relPath, Status: StatusTypeChanged, Type: currentNode.nodeType, ContentDiff: true, Detail: "snapshot type " + expectedNode.nodeType})
			continue
		}
		if !bytes.Equal(currentNode.content, expectedNode.content) {
			entries = append(entries, StatusEntry{Path: relPath, Status: StatusModified, Type: currentNode.nodeType, ContentDiff: true})
		}
	}
	for relPath, currentNode := range current {
		if exclusions.excludes(relPath) {
			continue
		}
		if _, found := expected[relPath]; !found {
			entries = append(entries, StatusEntry{Path: relPath, Status: StatusTrackingAdded, Type: currentNode.nodeType, TrackedStatusDiff: true, Detail: "local policy tracks path absent from snapshot"})
		}
	}
	entries = append(entries, trackingRemovedEntries(expected, exclusions)...)
	sort.Slice(entries, func(left, right int) bool {
		return entries[left].Path < entries[right].Path
	})
	return entries
}

// trackingRemovedEntries reports local exclusions that remove paths still
// present in the latest snapshot.
//
// Intent: `untrack` must be visible immediately as a tracked-status difference,
// but excluded bytes are not content drift because local policy says the next
// snapshot should omit them. Source: DI-tuhoj
func trackingRemovedEntries(expected map[string]expectedNode, exclusions pathExclusionSet) []StatusEntry {
	entries := []StatusEntry{}
	for excluded := range exclusions.paths {
		if exclusions.excludedByOtherPath(excluded) {
			continue
		}
		if node, found := expected[excluded]; found {
			entries = append(entries, StatusEntry{Path: excluded, Status: StatusTrackingRemoved, Type: node.nodeType, TrackedStatusDiff: true, Detail: "local policy excludes path present in snapshot"})
			continue
		}
		if snapshotHasDescendant(expected, excluded) {
			entries = append(entries, StatusEntry{Path: excluded, Status: StatusTrackingRemoved, Type: "directory", TrackedStatusDiff: true, Detail: "local policy excludes snapshot descendants"})
		}
	}
	return entries
}

// snapshotHasDescendant lets a directory exclusion report one parent-level
// tracked-status difference instead of one entry for every excluded child.
func snapshotHasDescendant(expected map[string]expectedNode, relPath string) bool {
	prefix := relPath + "/"
	for expectedPath := range expected {
		if strings.HasPrefix(expectedPath, prefix) {
			return true
		}
	}
	return false
}

// summarizeStatusEntries sets aggregate status flags from per-entry flags.
func summarizeStatusEntries(entries []StatusEntry) (bool, bool) {
	contentDiff := false
	trackedStatusDiff := false
	for _, entry := range entries {
		if entry.ContentDiff {
			contentDiff = true
		}
		if entry.TrackedStatusDiff {
			trackedStatusDiff = true
		}
	}
	return contentDiff, trackedStatusDiff
}

func statusTypeForInfo(info os.FileInfo) string {
	mode := info.Mode()
	switch {
	case mode.IsRegular():
		return "regular"
	case mode.IsDir():
		return "directory"
	case mode&os.ModeSymlink != 0:
		return "symlink"
	case mode&os.ModeNamedPipe != 0:
		return "fifo"
	case mode&os.ModeSocket != 0:
		return "socket"
	case mode&os.ModeDevice != 0 && mode&os.ModeCharDevice != 0:
		return "char_device"
	case mode&os.ModeDevice != 0:
		return "block_device"
	default:
		return "unknown"
	}
}

func joinStatusPath(relDir string, label string) string {
	if relDir == "" {
		return label
	}
	return relDir + "/" + label
}

type pathExclusionSet struct {
	paths map[string]bool
}

func newPathExclusionSet(paths []string) pathExclusionSet {
	set := pathExclusionSet{paths: map[string]bool{}}
	for _, path := range paths {
		cleaned := filepath.ToSlash(filepath.Clean(path))
		if cleaned == "." || cleaned == "" {
			continue
		}
		set.paths[cleaned] = true
	}
	return set
}

func (set pathExclusionSet) excludes(relPath string) bool {
	cleaned := filepath.ToSlash(filepath.Clean(relPath))
	if cleaned == "." || cleaned == "" {
		return false
	}
	for excluded := range set.paths {
		if cleaned == excluded || strings.HasPrefix(cleaned, excluded+"/") {
			return true
		}
	}
	return false
}

// excludedByOtherPath suppresses duplicate tracking removals when both a parent
// path and one of its descendants are excluded.
func (set pathExclusionSet) excludedByOtherPath(relPath string) bool {
	for excluded := range set.paths {
		if excluded == relPath {
			continue
		}
		if strings.HasPrefix(relPath, excluded+"/") {
			return true
		}
	}
	return false
}
