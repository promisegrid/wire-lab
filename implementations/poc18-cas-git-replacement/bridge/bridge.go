// Package bridge implements POC18's conventional Git compatibility seam.
//
// Intent: Git import/export/push/pull must share one conversion core, but this
// adapter must not let Git remotes become the native PromiseGrid authority
// model. Source: DI-dofoj; DI-harih; DI-fimap
package bridge

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	cidlib "github.com/ipfs/go-cid"

	pocchunk "promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/chunk"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/workspace"
)

// Operation names conventional Git bridge directions. These are compatibility
// adapter operations, not native PromiseGrid synchronization.
type Operation string

const (
	Import Operation = "import"
	Export Operation = "export"
	Pull   Operation = "pull"
	Push   Operation = "push"
)

// Mapping records one conversion seam between Git labels and native CIDs.
type Mapping struct {
	Operation Operation `json:"operation"`
	GitLabel  string    `json:"git_label"`
	GitHash   string    `json:"git_hash,omitempty"`
	GridRole  string    `json:"grid_role"`
	GridCID   string    `json:"grid_cid"`
	Outcome   string    `json:"outcome"`
}

// Result records one Git bridge operation and the local promises it emitted.
type Result struct {
	Operation         Operation           `json:"operation"`
	GitContext        string              `json:"git_context"`
	RepositoryPath    string              `json:"repository_path,omitempty"`
	RemoteURL         string              `json:"remote_url,omitempty"`
	WorktreePath      string              `json:"worktree_path,omitempty"`
	HeadSnapshotCID   string              `json:"head_snapshot_cid,omitempty"`
	HeadGitHash       string              `json:"head_git_hash,omitempty"`
	MappingMessageCID string              `json:"mapping_message_cid,omitempty"`
	Mappings          []Mapping           `json:"mappings"`
	LossRecords       []map[string]string `json:"loss_records,omitempty"`
	Counts            map[string]int      `json:"counts"`
}

// Adapter owns one local Git compatibility operation against one sparse CAS.
type Adapter struct {
	CAS      *store.FileStore
	Promiser string
	Promisee string
}

// NewAdapter returns a Git bridge adapter using the same CAS as native graph
// code.
func NewAdapter(cas *store.FileStore, promiser string, promisee string) *Adapter {
	return &Adapter{CAS: cas, Promiser: promiser, Promisee: promisee}
}

// ImportRepository converts conventional Git refs and commits into native
// PromiseGrid graph objects.
//
// Intent: Import is a compatibility adapter promise. It records how Git refs and
// commits were mapped into native snapshot/reference-set CIDs without treating
// those Git refs as native authority. Source: DI-fimap
func (adapter *Adapter) ImportRepository(repositoryPath string) (Result, error) {
	return adapter.importRepository(Import, repositoryPath, repositoryPath)
}

// PullRepository clones a conventional Git remote into worktreePath, then uses
// the same import conversion core as ImportRepository.
func (adapter *Adapter) PullRepository(remoteURL string, worktreePath string) (Result, error) {
	if remoteURL == "" || worktreePath == "" {
		return Result{}, fmt.Errorf("remote URL and worktree path are required")
	}
	if err := os.MkdirAll(filepath.Dir(worktreePath), 0o755); err != nil {
		return Result{}, err
	}
	if _, statErr := os.Stat(worktreePath); statErr == nil {
		return Result{}, fmt.Errorf("pull worktree already exists: %s", worktreePath)
	} else if !os.IsNotExist(statErr) {
		return Result{}, statErr
	}
	if _, cloneErr := git.PlainClone(worktreePath, false, &git.CloneOptions{URL: remoteURL, Tags: git.AllTags}); cloneErr != nil {
		return Result{}, cloneErr
	}
	result, importErr := adapter.importRepository(Pull, worktreePath, remoteURL)
	if importErr != nil {
		return Result{}, importErr
	}
	result.RemoteURL = remoteURL
	result.WorktreePath = worktreePath
	return result, nil
}

// ExportSnapshot converts one native snapshot graph into a conventional Git
// repository.
func (adapter *Adapter) ExportSnapshot(snapshotCID cidlib.Cid, repositoryPath string) (Result, error) {
	return adapter.exportSnapshot(Export, snapshotCID, repositoryPath, repositoryPath, "")
}

// PushSnapshot exports one native snapshot through the same conversion core as
// ExportSnapshot, then pushes the resulting branch to a conventional Git remote.
func (adapter *Adapter) PushSnapshot(snapshotCID cidlib.Cid, remoteURL string, worktreePath string) (Result, error) {
	if remoteURL == "" || worktreePath == "" {
		return Result{}, fmt.Errorf("remote URL and worktree path are required")
	}
	result, exportErr := adapter.exportSnapshot(Push, snapshotCID, worktreePath, remoteURL, remoteURL)
	if exportErr != nil {
		return Result{}, exportErr
	}
	repository, openErr := git.PlainOpen(worktreePath)
	if openErr != nil {
		return Result{}, openErr
	}
	if _, remoteErr := repository.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{remoteURL}}); remoteErr != nil && !strings.Contains(remoteErr.Error(), "remote already exists") {
		return Result{}, remoteErr
	}
	pushErr := repository.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.RefSpec("+refs/heads/main:refs/heads/main")},
	})
	if pushErr != nil && !errors.Is(pushErr, git.NoErrAlreadyUpToDate) {
		return Result{}, pushErr
	}
	result.RemoteURL = remoteURL
	result.WorktreePath = worktreePath
	return result, nil
}

func (adapter *Adapter) importRepository(operation Operation, repositoryPath string, gitContext string) (Result, error) {
	if adapter.CAS == nil {
		return Result{}, fmt.Errorf("CAS is required")
	}
	repository, openErr := git.PlainOpen(repositoryPath)
	if openErr != nil {
		return Result{}, openErr
	}
	importer := &gitImporter{
		adapter:        adapter,
		repository:     repository,
		operation:      operation,
		gitContext:     gitContext,
		commitSnapshot: map[plumbing.Hash]cidlib.Cid{},
		treeReference:  map[plumbing.Hash]cidlib.Cid{},
		treeNode:       map[plumbing.Hash]cidlib.Cid{},
		counts:         map[string]int{},
	}
	if importErr := importer.importReferences(); importErr != nil {
		return Result{}, importErr
	}
	mappingCID, mappingErr := adapter.storeMappingMessage(operation, gitContext, importer.mappings, importer.lossRecords)
	if mappingErr != nil {
		return Result{}, mappingErr
	}
	return Result{
		Operation:         operation,
		GitContext:        gitContext,
		RepositoryPath:    repositoryPath,
		HeadSnapshotCID:   importer.headSnapshotCID,
		HeadGitHash:       importer.headGitHash,
		MappingMessageCID: store.CIDText(mappingCID),
		Mappings:          importer.mappings,
		LossRecords:       importer.lossRecords,
		Counts:            importer.counts,
	}, nil
}

func (adapter *Adapter) exportSnapshot(operation Operation, snapshotCID cidlib.Cid, repositoryPath string, gitContext string, remoteURL string) (Result, error) {
	if adapter.CAS == nil {
		return Result{}, fmt.Errorf("CAS is required")
	}
	repository, repoErr := openOrInitRepository(repositoryPath)
	if repoErr != nil {
		return Result{}, repoErr
	}
	exporter := &gitExporter{
		adapter:        adapter,
		repository:     repository,
		repositoryPath: repositoryPath,
		operation:      operation,
		gitContext:     gitContext,
		snapshotCommit: map[string]plumbing.Hash{},
		counts:         map[string]int{},
	}
	headHash, exportErr := exporter.exportSnapshot(snapshotCID)
	if exportErr != nil {
		return Result{}, exportErr
	}
	mappingCID, mappingErr := adapter.storeMappingMessage(operation, gitContext, exporter.mappings, exporter.lossRecords)
	if mappingErr != nil {
		return Result{}, mappingErr
	}
	return Result{
		Operation:         operation,
		GitContext:        gitContext,
		RepositoryPath:    repositoryPath,
		RemoteURL:         remoteURL,
		WorktreePath:      repositoryPath,
		HeadSnapshotCID:   store.CIDText(snapshotCID),
		HeadGitHash:       headHash.String(),
		MappingMessageCID: store.CIDText(mappingCID),
		Mappings:          exporter.mappings,
		LossRecords:       exporter.lossRecords,
		Counts:            exporter.counts,
	}, nil
}

func (adapter *Adapter) storeMappingMessage(operation Operation, gitContext string, mappings []Mapping, lossRecords []map[string]string) (cidlib.Cid, error) {
	mappingRows := make([]any, 0, len(mappings))
	for _, mapping := range mappings {
		gridCID, parseErr := store.ParseCIDText(mapping.GridCID)
		if parseErr != nil {
			return cidlib.Undef, parseErr
		}
		mappingRows = append(mappingRows, []any{
			mapping.GitLabel,
			mapping.GitHash,
			mapping.GridRole,
			store.LinkTag(gridCID),
			mapping.Outcome,
		})
	}
	lossRows := make([]any, 0, len(lossRecords))
	for _, lossRecord := range lossRecords {
		lossRows = append(lossRows, lossRecord)
	}
	payload := graph.Payload{
		Promiser:    adapter.Promiser,
		Promisee:    adapter.Promisee,
		PromiseKind: "git_bridge_mapping",
		PromiseBody: graph.GitBridgeMappingBody(
			string(operation),
			gitContext,
			mappingRows,
			lossRows,
			[]any{"compatibility bridge only", "native sync remains peer DAG synchronization"},
		),
	}
	message, storeErr := graph.StoreMessage(adapter.CAS, nil, payload)
	if storeErr != nil {
		return cidlib.Undef, storeErr
	}
	return message.CID, nil
}

type gitImporter struct {
	adapter         *Adapter
	repository      *git.Repository
	operation       Operation
	gitContext      string
	commitSnapshot  map[plumbing.Hash]cidlib.Cid
	treeReference   map[plumbing.Hash]cidlib.Cid
	treeNode        map[plumbing.Hash]cidlib.Cid
	mappings        []Mapping
	lossRecords     []map[string]string
	counts          map[string]int
	headSnapshotCID string
	headGitHash     string
}

func (importer *gitImporter) importReferences() error {
	iterator, refsErr := importer.repository.References()
	if refsErr != nil {
		return refsErr
	}
	return iterator.ForEach(func(reference *plumbing.Reference) error {
		if !importer.shouldImportReference(reference) {
			return nil
		}
		commit, commitErr := importer.commitForReference(reference)
		if commitErr != nil {
			importer.lossRecords = append(importer.lossRecords, map[string]string{
				"git_ref":  reference.Name().String(),
				"git_hash": reference.Hash().String(),
				"outcome":  "not_commit",
				"reason":   commitErr.Error(),
			})
			return nil
		}
		snapshotCID, convertErr := importer.convertCommit(commit)
		if convertErr != nil {
			return convertErr
		}
		refSetCID, refErr := importer.storeGitReference(reference, snapshotCID)
		if refErr != nil {
			return refErr
		}
		importer.mappings = append(importer.mappings, Mapping{
			Operation: importer.operation,
			GitLabel:  reference.Name().String(),
			GitHash:   reference.Hash().String(),
			GridRole:  "reference_set",
			GridCID:   store.CIDText(refSetCID),
			Outcome:   "mapped",
		})
		if reference.Name().IsBranch() || importer.headSnapshotCID == "" {
			importer.headSnapshotCID = store.CIDText(snapshotCID)
			importer.headGitHash = commit.Hash.String()
		}
		return nil
	})
}

func (importer *gitImporter) shouldImportReference(reference *plumbing.Reference) bool {
	if reference.Type() != plumbing.HashReference {
		return false
	}
	name := reference.Name()
	if name.IsBranch() || name.IsTag() {
		return true
	}
	nameText := name.String()
	return strings.HasPrefix(nameText, "refs/remotes/") && !strings.HasSuffix(nameText, "/HEAD")
}

func (importer *gitImporter) commitForReference(reference *plumbing.Reference) (*object.Commit, error) {
	commit, commitErr := importer.repository.CommitObject(reference.Hash())
	if commitErr == nil {
		return commit, nil
	}
	tag, tagErr := importer.repository.TagObject(reference.Hash())
	if tagErr != nil {
		return nil, commitErr
	}
	return tag.Commit()
}

func (importer *gitImporter) convertCommit(commit *object.Commit) (cidlib.Cid, error) {
	if snapshotCID, ok := importer.commitSnapshot[commit.Hash]; ok {
		return snapshotCID, nil
	}
	parentCIDs := make([]cidlib.Cid, 0, len(commit.ParentHashes))
	parents := make([]graph.Parent, 0, len(commit.ParentHashes))
	for _, parentHash := range commit.ParentHashes {
		parentCommit, parentErr := importer.repository.CommitObject(parentHash)
		if parentErr != nil {
			return cidlib.Undef, parentErr
		}
		parentCID, convertErr := importer.convertCommit(parentCommit)
		if convertErr != nil {
			return cidlib.Undef, convertErr
		}
		parentCIDs = append(parentCIDs, parentCID)
		parents = append(parents, graph.Parent{Role: "previous_snapshot", CID: parentCID})
	}
	tree, treeErr := commit.Tree()
	if treeErr != nil {
		return cidlib.Undef, treeErr
	}
	rootRefCID, _, treeConvertErr := importer.convertTree(tree, ".")
	if treeConvertErr != nil {
		return cidlib.Undef, treeConvertErr
	}
	payload := graph.Payload{
		Promiser:    importer.adapter.Promiser,
		Promisee:    importer.adapter.Promisee,
		PromiseKind: "snapshot",
		PromiseBody: graph.SnapshotBody(
			"git:commit:"+commit.Hash.String(),
			rootRefCID,
			parentCIDs,
			strings.TrimSpace(commit.Message),
			[]any{"imported from Git commit", commit.Hash.String()},
		),
	}
	message, storeErr := graph.StoreMessage(importer.adapter.CAS, parents, payload)
	if storeErr != nil {
		return cidlib.Undef, storeErr
	}
	importer.commitSnapshot[commit.Hash] = message.CID
	importer.counts["git_commit"]++
	importer.counts["snapshot"]++
	importer.mappings = append(importer.mappings, Mapping{
		Operation: importer.operation,
		GitLabel:  "commit:" + commit.Hash.String(),
		GitHash:   commit.Hash.String(),
		GridRole:  "snapshot",
		GridCID:   store.CIDText(message.CID),
		Outcome:   "mapped",
	})
	return message.CID, nil
}

func (importer *gitImporter) convertTree(tree *object.Tree, relativePath string) (cidlib.Cid, cidlib.Cid, error) {
	if refCID, ok := importer.treeReference[tree.Hash]; ok {
		return refCID, importer.treeNode[tree.Hash], nil
	}
	entries := append([]object.TreeEntry(nil), tree.Entries...)
	sort.Slice(entries, func(left, right int) bool {
		return entries[left].Name < entries[right].Name
	})
	referenceEntries := []any{}
	for _, entry := range entries {
		childCID, childErr := importer.convertTreeEntry(entry, relativePath)
		if childErr != nil {
			return cidlib.Undef, cidlib.Undef, childErr
		}
		if childCID.Defined() {
			referenceEntries = append(referenceEntries, graph.ReferenceEntry(entry.Name, []any{graph.Target("node", childCID)}, []any{"git_mode:" + entry.Mode.String()}))
		}
	}
	refPayload := graph.Payload{
		Promiser:    importer.adapter.Promiser,
		Promisee:    importer.adapter.Promisee,
		PromiseKind: "reference_set",
		PromiseBody: graph.ReferenceSetBody("git:tree:"+tree.Hash.String(), "directory", "git:"+relativePath, referenceEntries, []any{"directory labels imported from Git tree"}),
	}
	refMessage, refErr := graph.StoreMessage(importer.adapter.CAS, nil, refPayload)
	if refErr != nil {
		return cidlib.Undef, cidlib.Undef, refErr
	}
	dirPayload := graph.Payload{
		Promiser:    importer.adapter.Promiser,
		Promisee:    importer.adapter.Promisee,
		PromiseKind: "posix_node",
		PromiseBody: graph.PosixNodeBody("git:tree-node:"+tree.Hash.String(), "directory", store.LinkTag(refMessage.CID), map[string]any{"git_hash": tree.Hash.String(), "mode": "drwxr-xr-x", "perm": "0755"}, []any{"materialize directory locally if safe"}),
	}
	dirMessage, dirErr := graph.StoreMessage(importer.adapter.CAS, []graph.Parent{{Role: "previous_reference_set", CID: refMessage.CID}}, dirPayload)
	if dirErr != nil {
		return cidlib.Undef, cidlib.Undef, dirErr
	}
	importer.treeReference[tree.Hash] = refMessage.CID
	importer.treeNode[tree.Hash] = dirMessage.CID
	importer.counts["reference_set:directory"]++
	importer.counts["posix_node:directory"]++
	return refMessage.CID, dirMessage.CID, nil
}

func (importer *gitImporter) convertTreeEntry(entry object.TreeEntry, relativePath string) (cidlib.Cid, error) {
	childPath := filepath.Join(relativePath, entry.Name)
	switch entry.Mode {
	case filemode.Dir:
		childTree, treeErr := importer.repository.TreeObject(entry.Hash)
		if treeErr != nil {
			return cidlib.Undef, treeErr
		}
		_, dirNodeCID, convertErr := importer.convertTree(childTree, childPath)
		return dirNodeCID, convertErr
	case filemode.Regular, filemode.Executable, filemode.Deprecated:
		return importer.convertBlobNode(entry, childPath, false)
	case filemode.Symlink:
		return importer.convertBlobNode(entry, childPath, true)
	case filemode.Submodule:
		importer.lossRecords = append(importer.lossRecords, map[string]string{
			"git_path": childPath,
			"git_hash": entry.Hash.String(),
			"outcome":  "not_mapped",
			"reason":   "git submodule has no POSIX object equivalent in this POC18 bridge slice",
		})
		return cidlib.Undef, nil
	default:
		return cidlib.Undef, fmt.Errorf("unsupported Git file mode %s at %s", entry.Mode, childPath)
	}
}

func (importer *gitImporter) convertBlobNode(entry object.TreeEntry, relativePath string, symlink bool) (cidlib.Cid, error) {
	blob, blobErr := importer.repository.BlobObject(entry.Hash)
	if blobErr != nil {
		return cidlib.Undef, blobErr
	}
	reader, readerErr := blob.Reader()
	if readerErr != nil {
		return cidlib.Undef, readerErr
	}
	content, readErr := io.ReadAll(reader)
	closeErr := reader.Close()
	if readErr != nil {
		return cidlib.Undef, readErr
	}
	if closeErr != nil {
		return cidlib.Undef, closeErr
	}
	nodeType := "regular"
	nodeContent := any(nil)
	metadata := map[string]any{"git_hash": entry.Hash.String(), "git_mode": entry.Mode.String()}
	if symlink {
		nodeType = "symlink"
		nodeContent = content
		metadata["mode"] = "Lrwxrwxrwx"
		metadata["perm"] = "0777"
	} else {
		storedManifest, chunkErr := pocchunk.StoreBytes(importer.adapter.CAS, content)
		if chunkErr != nil {
			return cidlib.Undef, chunkErr
		}
		chunkRows := make([]any, 0, len(storedManifest.Manifest.Chunks))
		for _, chunkRef := range storedManifest.Manifest.Chunks {
			chunkCID, parseErr := store.ParseCIDText(chunkRef.CID)
			if parseErr != nil {
				return cidlib.Undef, parseErr
			}
			chunkRows = append(chunkRows, []any{chunkRef.Offset, chunkRef.Length, store.LinkTag(chunkCID)})
		}
		chunkPayload := graph.Payload{
			Promiser:    importer.adapter.Promiser,
			Promisee:    importer.adapter.Promisee,
			PromiseKind: "chunk_manifest",
			PromiseBody: graph.ChunkManifestBody(storedManifest.ManifestCID, storedManifest.Manifest.FileSize, "rabin", storedManifest.Manifest.Params, chunkRows, "git:"+entry.Hash.String()),
		}
		if _, msgErr := graph.StoreMessage(importer.adapter.CAS, nil, chunkPayload); msgErr != nil {
			return cidlib.Undef, msgErr
		}
		nodeContent = store.LinkTag(storedManifest.ManifestCID)
		metadata["mode"] = "-rw-r--r--"
		metadata["perm"] = "0644"
		if entry.Mode == filemode.Executable {
			metadata["mode"] = "-rwxr-xr-x"
			metadata["perm"] = "0755"
		}
		importer.counts["chunk"] += len(storedManifest.Manifest.Chunks)
		importer.counts["chunk_manifest"]++
	}
	payload := graph.Payload{
		Promiser:    importer.adapter.Promiser,
		Promisee:    importer.adapter.Promisee,
		PromiseKind: "posix_node",
		PromiseBody: graph.PosixNodeBody("git:blob:"+entry.Hash.String()+":"+relativePath, nodeType, nodeContent, metadata, []any{"materialize imported Git object locally if safe"}),
	}
	message, storeErr := graph.StoreMessage(importer.adapter.CAS, nil, payload)
	if storeErr != nil {
		return cidlib.Undef, storeErr
	}
	importer.counts["posix_node:"+nodeType]++
	return message.CID, nil
}

func (importer *gitImporter) storeGitReference(reference *plumbing.Reference, snapshotCID cidlib.Cid) (cidlib.Cid, error) {
	refRole := "git_ref"
	if reference.Name().IsBranch() {
		refRole = "branch"
	} else if reference.Name().IsTag() {
		refRole = "tag"
	} else if strings.HasPrefix(reference.Name().String(), "refs/remotes/") {
		refRole = "remote_tracking"
	}
	payload := graph.Payload{
		Promiser:    importer.adapter.Promiser,
		Promisee:    importer.adapter.Promisee,
		PromiseKind: "reference_set",
		PromiseBody: graph.ReferenceSetBody(
			"git-ref:"+reference.Name().String(),
			refRole,
			"git:"+importer.gitContext,
			[]any{graph.ReferenceEntry(reference.Name().Short(), []any{graph.Target("snapshot", snapshotCID)}, []any{"imported Git ref", reference.Hash().String()})},
			[]any{"Git ref compatibility label, not native forge authority"},
		),
	}
	message, storeErr := graph.StoreMessage(importer.adapter.CAS, nil, payload)
	if storeErr != nil {
		return cidlib.Undef, storeErr
	}
	importer.counts["reference_set:"+refRole]++
	return message.CID, nil
}

type gitExporter struct {
	adapter        *Adapter
	repository     *git.Repository
	repositoryPath string
	operation      Operation
	gitContext     string
	snapshotCommit map[string]plumbing.Hash
	mappings       []Mapping
	lossRecords    []map[string]string
	counts         map[string]int
}

func (exporter *gitExporter) exportSnapshot(snapshotCID cidlib.Cid) (plumbing.Hash, error) {
	snapshotText := store.CIDText(snapshotCID)
	if commitHash, ok := exporter.snapshotCommit[snapshotText]; ok {
		return commitHash, nil
	}
	parts, partsErr := exporter.snapshotParts(snapshotCID)
	if partsErr != nil {
		return plumbing.ZeroHash, partsErr
	}
	parentHashes := make([]plumbing.Hash, 0, len(parts.parentCIDs))
	for _, parentCID := range parts.parentCIDs {
		parentHash, parentErr := exporter.exportSnapshot(parentCID)
		if parentErr != nil {
			return plumbing.ZeroHash, parentErr
		}
		parentHashes = append(parentHashes, parentHash)
	}
	if err := cleanWorktree(exporter.repositoryPath); err != nil {
		return plumbing.ZeroHash, err
	}
	if err := workspace.MaterializeSnapshot(exporter.adapter.CAS, snapshotCID, exporter.repositoryPath); err != nil {
		return plumbing.ZeroHash, err
	}
	worktree, worktreeErr := exporter.repository.Worktree()
	if worktreeErr != nil {
		return plumbing.ZeroHash, worktreeErr
	}
	if addErr := worktree.AddWithOptions(&git.AddOptions{All: true}); addErr != nil {
		return plumbing.ZeroHash, addErr
	}
	commitHash, commitErr := worktree.Commit(parts.summary, &git.CommitOptions{
		Author:            bridgeSignature(exporter.adapter.Promiser),
		Committer:         bridgeSignature(exporter.adapter.Promiser),
		Parents:           parentHashes,
		AllowEmptyCommits: true,
	})
	if commitErr != nil {
		return plumbing.ZeroHash, commitErr
	}
	exporter.snapshotCommit[snapshotText] = commitHash
	exporter.counts["git_commit"]++
	exporter.mappings = append(exporter.mappings, Mapping{
		Operation: exporter.operation,
		GitLabel:  "commit:" + commitHash.String(),
		GitHash:   commitHash.String(),
		GridRole:  "snapshot",
		GridCID:   snapshotText,
		Outcome:   "mapped",
	})
	return commitHash, nil
}

type snapshotParts struct {
	parentCIDs []cidlib.Cid
	summary    string
}

func (exporter *gitExporter) snapshotParts(snapshotCID cidlib.Cid) (snapshotParts, error) {
	content, _, getErr := exporter.adapter.CAS.Get(snapshotCID)
	if getErr != nil {
		return snapshotParts{}, getErr
	}
	view, parseErr := graph.ParseEnvelope(content)
	if parseErr != nil {
		return snapshotParts{}, parseErr
	}
	kind, kindErr := view.PayloadKind()
	if kindErr != nil {
		return snapshotParts{}, kindErr
	}
	if kind != "snapshot" {
		return snapshotParts{}, fmt.Errorf("expected snapshot, got %s", kind)
	}
	body, bodyErr := view.PayloadBody()
	if bodyErr != nil {
		return snapshotParts{}, bodyErr
	}
	if len(body) != 5 {
		return snapshotParts{}, fmt.Errorf("snapshot body must have five slots")
	}
	parentValues, ok := body[2].([]any)
	if !ok {
		return snapshotParts{}, fmt.Errorf("snapshot parents must be array")
	}
	parents := make([]cidlib.Cid, 0, len(parentValues))
	for _, parentValue := range parentValues {
		parentCID, parentErr := store.CIDFromLinkTag(parentValue)
		if parentErr != nil {
			return snapshotParts{}, parentErr
		}
		parents = append(parents, parentCID)
	}
	summary, _ := body[3].(string)
	if strings.TrimSpace(summary) == "" {
		summary = "PromiseGrid snapshot " + store.CIDText(snapshotCID)
	}
	return snapshotParts{parentCIDs: parents, summary: summary}, nil
}

func openOrInitRepository(repositoryPath string) (*git.Repository, error) {
	if repositoryPath == "" {
		return nil, fmt.Errorf("repository path is required")
	}
	if err := os.MkdirAll(repositoryPath, 0o755); err != nil {
		return nil, err
	}
	repository, openErr := git.PlainOpen(repositoryPath)
	if openErr == nil {
		return repository, nil
	}
	if !errors.Is(openErr, git.ErrRepositoryNotExists) {
		return nil, openErr
	}
	repository, initErr := git.PlainInit(repositoryPath, false)
	if initErr != nil {
		return nil, initErr
	}
	if setErr := repository.Storer.SetReference(plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.NewBranchReferenceName("main"))); setErr != nil {
		return nil, setErr
	}
	return repository, nil
}

func cleanWorktree(repositoryPath string) error {
	entries, readErr := os.ReadDir(repositoryPath)
	if readErr != nil {
		return readErr
	}
	for _, entry := range entries {
		if entry.Name() == ".git" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(repositoryPath, entry.Name())); err != nil {
			return err
		}
	}
	return nil
}

func bridgeSignature(agent string) *object.Signature {
	if agent == "" {
		agent = "promisegrid"
	}
	return &object.Signature{
		Name:  agent,
		Email: agent + "@poc18.promisegrid.local",
		When:  time.Unix(1_800_000_000, 0).UTC(),
	}
}
