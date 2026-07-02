// Package scenario builds deterministic POC18 collaboration scenarios.
//
// Intent: Keep higher-level POC fixture construction out of cmd packages while
// still using the same graph/store/workspace library paths as CLI code. Source:
// DI-guban
package scenario

import (
	"fmt"

	cidlib "github.com/ipfs/go-cid"

	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/graph"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/store"
	"promisegrid.dev/wire-lab/implementations/poc18-cas-git-replacement/workspace"
)

// Result records the exact CIDs produced by the deterministic collaboration
// scenario.
type Result struct {
	InitialSnapshotCID         string            `json:"initial_snapshot_cid"`
	InitialRootReferenceSetCID string            `json:"initial_root_reference_set_cid"`
	LineageNodeCID             string            `json:"lineage_node_cid"`
	RenameLabelNodeCID         string            `json:"rename_label_node_cid"`
	CopyLabelNodeCID           string            `json:"copy_label_node_cid"`
	RenameCopyDocsRefSetCID    string            `json:"rename_copy_docs_reference_set_cid"`
	RenameCopyRootRefSetCID    string            `json:"rename_copy_root_reference_set_cid"`
	RenameCopySnapshotCID      string            `json:"rename_copy_snapshot_cid"`
	BobDivergentRootRefSetCID  string            `json:"bob_divergent_root_reference_set_cid"`
	BobDivergentSnapshotCID    string            `json:"bob_divergent_snapshot_cid"`
	MergeRootRefSetCID         string            `json:"merge_root_reference_set_cid"`
	MergeSnapshotCID           string            `json:"merge_snapshot_cid"`
	MergeParentSnapshotCIDs    []string          `json:"merge_parent_snapshot_cids"`
	MergeBranchRefSetCID       string            `json:"merge_branch_reference_set_cid"`
	LogicalChangeRefSetCID     string            `json:"logical_change_reference_set_cid"`
	TestStatementCID           string            `json:"test_statement_cid"`
	AdoptionStatementCID       string            `json:"adoption_statement_cid"`
	ReviewThreadRefSetCID      string            `json:"review_thread_reference_set_cid"`
	ReviewAdoptionResult       string            `json:"review_adoption_result"`
	Counts                     map[string]int    `json:"counts"`
	LineageLabels              map[string]string `json:"lineage_labels"`
}

// Builder owns one deterministic scenario build against a sparse CAS.
type Builder struct {
	CAS *store.FileStore
}

// NewBuilder returns a scenario builder for one local sparse CAS.
func NewBuilder(cas *store.FileStore) *Builder {
	return &Builder{CAS: cas}
}

// Build appends rename/copy, divergence, merge, review, test, and local adoption
// promises to an initial workspace ingest.
//
// Intent: Prove `nahop.12` through `nahop.14` as PromiseGrid-native messages:
// labels change in reference sets, divergent work is explicit, merges are
// multi-parent snapshots, and reviews/tests/adoptions remain local promises.
// Source: DI-guban
func (builder *Builder) Build(initial workspace.IngestResult) (Result, error) {
	if builder.CAS == nil {
		return Result{}, fmt.Errorf("scenario CAS is required")
	}
	cids, loadErr := builder.loadInitialCIDs(initial)
	if loadErr != nil {
		return Result{}, loadErr
	}
	rootTargets, rootErr := builder.referenceSetTargets(cids.rootRef)
	if rootErr != nil {
		return Result{}, rootErr
	}
	readmeNode, readmeErr := targetCID(rootTargets, "README.md")
	if readmeErr != nil {
		return Result{}, readmeErr
	}
	docsNode, docsErr := targetCID(rootTargets, "docs")
	if docsErr != nil {
		return Result{}, docsErr
	}
	linkNode, linkErr := targetCID(rootTargets, "README-link.md")
	if linkErr != nil {
		return Result{}, linkErr
	}
	pipeNode, pipeErr := targetCID(rootTargets, "build.pipe")
	if pipeErr != nil {
		return Result{}, pipeErr
	}
	docsRef, docsRefErr := builder.directoryReferenceCID(docsNode)
	if docsRefErr != nil {
		return Result{}, docsRefErr
	}
	docsTargets, docsTargetsErr := builder.referenceSetTargets(docsRef)
	if docsTargetsErr != nil {
		return Result{}, docsTargetsErr
	}
	largeNode, largeErr := targetCID(docsTargets, "large.bin")
	if largeErr != nil {
		return Result{}, largeErr
	}
	result := Result{
		InitialSnapshotCID:         initial.SnapshotCID,
		InitialRootReferenceSetCID: initial.RootReferenceSetCID,
		LineageNodeCID:             store.CIDText(readmeNode),
		LineageLabels:              map[string]string{},
		Counts:                     map[string]int{},
	}
	renameCopy, renameErr := builder.buildRenameCopy(cids.snapshot, cids.rootRef, docsRef, docsNode, readmeNode, largeNode, linkNode, pipeNode)
	if renameErr != nil {
		return Result{}, renameErr
	}
	result.RenameCopyDocsRefSetCID = store.CIDText(renameCopy.docsRefSet.CID)
	result.RenameCopyRootRefSetCID = store.CIDText(renameCopy.rootRefSet.CID)
	result.RenameCopySnapshotCID = store.CIDText(renameCopy.snapshot.CID)
	result.RenameLabelNodeCID = store.CIDText(readmeNode)
	result.CopyLabelNodeCID = store.CIDText(readmeNode)
	result.LineageLabels["docs/intro.md"] = store.CIDText(readmeNode)
	result.LineageLabels["docs/README-copy.md"] = store.CIDText(readmeNode)
	bobDivergence, bobErr := builder.buildBobDivergence(cids.snapshot, cids.rootRef, readmeNode, docsNode, linkNode, pipeNode)
	if bobErr != nil {
		return Result{}, bobErr
	}
	result.BobDivergentRootRefSetCID = store.CIDText(bobDivergence.rootRefSet.CID)
	result.BobDivergentSnapshotCID = store.CIDText(bobDivergence.snapshot.CID)
	merge, mergeErr := builder.buildMerge(cids.rootRef, renameCopy, bobDivergence, readmeNode, linkNode, pipeNode)
	if mergeErr != nil {
		return Result{}, mergeErr
	}
	result.MergeRootRefSetCID = store.CIDText(merge.rootRefSet.CID)
	result.MergeSnapshotCID = store.CIDText(merge.snapshot.CID)
	result.MergeParentSnapshotCIDs = []string{result.RenameCopySnapshotCID, result.BobDivergentSnapshotCID}
	refsets, refsetErr := builder.buildReferenceSets(cids, renameCopy, bobDivergence, merge)
	if refsetErr != nil {
		return Result{}, refsetErr
	}
	result.MergeBranchRefSetCID = store.CIDText(refsets.branch.CID)
	result.LogicalChangeRefSetCID = store.CIDText(refsets.logicalChange.CID)
	review, reviewErr := builder.buildReview(cids.reviewThread, merge.snapshot.CID)
	if reviewErr != nil {
		return Result{}, reviewErr
	}
	result.TestStatementCID = store.CIDText(review.testStatement.CID)
	result.AdoptionStatementCID = store.CIDText(review.adoptionStatement.CID)
	result.ReviewThreadRefSetCID = store.CIDText(review.thread.CID)
	result.ReviewAdoptionResult = "accepted_locally"
	result.Counts["rename_copy_lineage"] = 1
	result.Counts["divergent_snapshots"] = 2
	result.Counts["multi_parent_merge"] = 1
	result.Counts["review_statement:test_kept"] = 1
	result.Counts["review_statement:accepted_locally"] = 1
	return result, nil
}

type initialCIDs struct {
	snapshot      cidlib.Cid
	rootRef       cidlib.Cid
	branch        cidlib.Cid
	logicalChange cidlib.Cid
	reviewThread  cidlib.Cid
}

type target struct {
	role string
	cid  cidlib.Cid
}

type renameCopyMessages struct {
	docsRefSet graph.StoredMessage
	docsNode   graph.StoredMessage
	rootRefSet graph.StoredMessage
	snapshot   graph.StoredMessage
}

type divergenceMessages struct {
	rootRefSet graph.StoredMessage
	snapshot   graph.StoredMessage
}

type mergeMessages struct {
	rootRefSet graph.StoredMessage
	snapshot   graph.StoredMessage
}

type referenceSetMessages struct {
	branch        graph.StoredMessage
	logicalChange graph.StoredMessage
}

type reviewMessages struct {
	testStatement     graph.StoredMessage
	adoptionStatement graph.StoredMessage
	thread            graph.StoredMessage
}

func (builder *Builder) loadInitialCIDs(initial workspace.IngestResult) (initialCIDs, error) {
	snapshotCID, snapshotErr := store.ParseCIDText(initial.SnapshotCID)
	if snapshotErr != nil {
		return initialCIDs{}, snapshotErr
	}
	rootRefCID, rootErr := store.ParseCIDText(initial.RootReferenceSetCID)
	if rootErr != nil {
		return initialCIDs{}, rootErr
	}
	branchCID, branchErr := store.ParseCIDText(initial.BranchRefSetCID)
	if branchErr != nil {
		return initialCIDs{}, branchErr
	}
	logicalChangeCID, changeErr := store.ParseCIDText(initial.LogicalChangeCID)
	if changeErr != nil {
		return initialCIDs{}, changeErr
	}
	reviewThreadCID, reviewErr := store.ParseCIDText(initial.ReviewThreadCID)
	if reviewErr != nil {
		return initialCIDs{}, reviewErr
	}
	return initialCIDs{
		snapshot:      snapshotCID,
		rootRef:       rootRefCID,
		branch:        branchCID,
		logicalChange: logicalChangeCID,
		reviewThread:  reviewThreadCID,
	}, nil
}

// buildRenameCopy creates Alice's rename/copy branch. The important invariant is
// that `docs/intro.md` and `docs/README-copy.md` both point at the original
// README node CID; the filename labels change, but file lineage does not.
func (builder *Builder) buildRenameCopy(snapshotCID, rootRefCID, docsRefCID, docsNodeCID, readmeNodeCID, largeNodeCID, linkNodeCID, pipeNodeCID cidlib.Cid) (renameCopyMessages, error) {
	docsEntries := []any{
		graph.ReferenceEntry("README-copy.md", []any{graph.Target("node", readmeNodeCID)}, []any{"copy label promises same node lineage as README.md"}),
		graph.ReferenceEntry("intro.md", []any{graph.Target("node", readmeNodeCID)}, []any{"rename label promises same node lineage as README.md"}),
		graph.ReferenceEntry("large.bin", []any{graph.Target("node", largeNodeCID)}, nil),
	}
	docsRefSet, docsRefErr := builder.storeReferenceSet("alice", "bob", []graph.Parent{{Role: "previous_reference_set", CID: docsRefCID}}, "refset:directory:docs:rename-copy", "directory", "workspace:docs", docsEntries, []any{"I promise these labels preserve file lineage by reusing the same node CID."})
	if docsRefErr != nil {
		return renameCopyMessages{}, docsRefErr
	}
	docsNode, docsNodeErr := graph.StoreMessage(builder.CAS, []graph.Parent{{Role: "previous_node", CID: docsNodeCID}, {Role: "previous_reference_set", CID: docsRefSet.CID}}, graph.Payload{
		Promiser:    "alice",
		Promisee:    "bob",
		PromiseKind: "posix_node",
		PromiseBody: graph.PosixNodeBody("node:directory:docs", "directory", store.LinkTag(docsRefSet.CID), map[string]any{"mode": "drwxr-xr-x", "perm": "0755"}, []any{"materialize renamed/copied labels locally if safe"}),
	})
	if docsNodeErr != nil {
		return renameCopyMessages{}, docsNodeErr
	}
	rootEntries := []any{
		graph.ReferenceEntry("README-link.md", []any{graph.Target("node", linkNodeCID)}, nil),
		graph.ReferenceEntry("build.pipe", []any{graph.Target("node", pipeNodeCID)}, nil),
		graph.ReferenceEntry("docs", []any{graph.Target("node", docsNode.CID)}, []any{"directory version contains rename/copy labels"}),
	}
	rootRefSet, rootErr := builder.storeReferenceSet("alice", "bob", []graph.Parent{{Role: "previous_reference_set", CID: rootRefCID}}, "refset:directory:.:rename-copy", "directory", "workspace:.", rootEntries, []any{"I promise README.md moved to docs/intro.md and was copied as docs/README-copy.md without changing file lineage."})
	if rootErr != nil {
		return renameCopyMessages{}, rootErr
	}
	snapshot, snapshotErr := builder.storeSnapshot("alice", "bob", []graph.Parent{{Role: "previous_snapshot", CID: snapshotCID}}, "snapshot:rename-copy", rootRefSet.CID, []cidlib.Cid{snapshotCID}, "rename README.md to docs/intro.md and copy it as docs/README-copy.md", []any{"directory labels changed; node lineage did not"})
	if snapshotErr != nil {
		return renameCopyMessages{}, snapshotErr
	}
	return renameCopyMessages{docsRefSet: docsRefSet, docsNode: docsNode, rootRefSet: rootRefSet, snapshot: snapshot}, nil
}

// buildBobDivergence creates Bob's incompatible local revision. Bob keeps the
// original README.md label and adds a root README-copy.md label, producing a
// real divergent snapshot rather than rewriting Alice's promise.
func (builder *Builder) buildBobDivergence(snapshotCID, rootRefCID, readmeNodeCID, docsNodeCID, linkNodeCID, pipeNodeCID cidlib.Cid) (divergenceMessages, error) {
	rootEntries := []any{
		graph.ReferenceEntry("README-copy.md", []any{graph.Target("node", readmeNodeCID)}, []any{"copy label promises same node lineage as README.md"}),
		graph.ReferenceEntry("README-link.md", []any{graph.Target("node", linkNodeCID)}, nil),
		graph.ReferenceEntry("README.md", []any{graph.Target("node", readmeNodeCID)}, nil),
		graph.ReferenceEntry("build.pipe", []any{graph.Target("node", pipeNodeCID)}, nil),
		graph.ReferenceEntry("docs", []any{graph.Target("node", docsNodeCID)}, nil),
	}
	rootRefSet, rootErr := builder.storeReferenceSet("bob", "alice", []graph.Parent{{Role: "previous_reference_set", CID: rootRefCID}}, "refset:directory:.:bob-copy", "directory", "workspace:.", rootEntries, []any{"I promise this divergent directory keeps README.md and adds README-copy.md."})
	if rootErr != nil {
		return divergenceMessages{}, rootErr
	}
	snapshot, snapshotErr := builder.storeSnapshot("bob", "alice", []graph.Parent{{Role: "previous_snapshot", CID: snapshotCID}}, "snapshot:bob-copy", rootRefSet.CID, []cidlib.Cid{snapshotCID}, "divergent copy branch keeps README.md and adds README-copy.md", []any{"local divergent promise, not global truth"})
	if snapshotErr != nil {
		return divergenceMessages{}, snapshotErr
	}
	return divergenceMessages{rootRefSet: rootRefSet, snapshot: snapshot}, nil
}

// buildMerge creates Dave's local merge promise with two snapshot parents. The
// merge records how Dave resolves the label conflict but does not claim to be a
// global merge decision for Alice, Bob, or any other agent.
func (builder *Builder) buildMerge(rootRefCID cidlib.Cid, renameCopy renameCopyMessages, divergence divergenceMessages, readmeNodeCID, linkNodeCID, pipeNodeCID cidlib.Cid) (mergeMessages, error) {
	rootEntries := []any{
		graph.ReferenceEntry("README-copy.md", []any{graph.Target("node", readmeNodeCID)}, []any{"kept Bob's root copy label"}),
		graph.ReferenceEntry("README-link.md", []any{graph.Target("node", linkNodeCID)}, nil),
		graph.ReferenceEntry("build.pipe", []any{graph.Target("node", pipeNodeCID)}, nil),
		graph.ReferenceEntry("docs", []any{graph.Target("node", renameCopy.docsNode.CID)}, []any{"kept Alice's rename/copy docs directory"}),
	}
	rootRefSet, rootErr := builder.storeReferenceSet("dave", "alice", []graph.Parent{
		{Role: "previous_reference_set", CID: renameCopy.rootRefSet.CID},
		{Role: "previous_reference_set", CID: divergence.rootRefSet.CID},
		{Role: "previous_reference_set", CID: rootRefCID},
	}, "refset:directory:.:merged", "directory", "workspace:.", rootEntries, []any{"I promise this directory resolves the Alice/Bob label conflict by keeping docs/intro.md and README-copy.md."})
	if rootErr != nil {
		return mergeMessages{}, rootErr
	}
	snapshot, snapshotErr := builder.storeSnapshot("dave", "alice", []graph.Parent{
		{Role: "previous_snapshot", CID: renameCopy.snapshot.CID},
		{Role: "previous_snapshot", CID: divergence.snapshot.CID},
	}, "snapshot:merge-rename-copy", rootRefSet.CID, []cidlib.Cid{renameCopy.snapshot.CID, divergence.snapshot.CID}, "merge Alice rename/copy and Bob divergent copy", []any{"conflict resolved locally; no global merge authority"})
	if snapshotErr != nil {
		return mergeMessages{}, snapshotErr
	}
	return mergeMessages{rootRefSet: rootRefSet, snapshot: snapshot}, nil
}

// buildReferenceSets publishes the merged branch head and logical-change
// history. The logical change is a versioned reference set, not a mutable forge
// pull request or an intrinsic change-ID field.
func (builder *Builder) buildReferenceSets(cids initialCIDs, renameCopy renameCopyMessages, divergence divergenceMessages, merge mergeMessages) (referenceSetMessages, error) {
	branch, branchErr := builder.storeReferenceSet("dave", "bob", []graph.Parent{{Role: "previous_reference_set", CID: cids.branch}}, "refset:main", "branch", "project:poc18-demo", []any{
		graph.ReferenceEntry("head", []any{graph.Target("snapshot", merge.snapshot.CID)}, []any{"local merged branch head promise"}),
	}, []any{"I promise this is my current merged branch head."})
	if branchErr != nil {
		return referenceSetMessages{}, branchErr
	}
	logicalChange, changeErr := builder.storeReferenceSet("carol", "dave", []graph.Parent{{Role: "previous_reference_set", CID: cids.logicalChange}}, "change:rename-copy-merge", "logical_change", "project:poc18-demo", []any{
		graph.ReferenceEntry("alice-round", []any{graph.Target("snapshot", renameCopy.snapshot.CID)}, []any{"Alice rename/copy revision"}),
		graph.ReferenceEntry("bob-round", []any{graph.Target("snapshot", divergence.snapshot.CID)}, []any{"Bob divergent copy revision"}),
		graph.ReferenceEntry("merged", []any{graph.Target("snapshot", merge.snapshot.CID)}, []any{"Dave merge revision"}),
	}, []any{"I promise these snapshots are revisions of the same logical change."})
	if changeErr != nil {
		return referenceSetMessages{}, changeErr
	}
	return referenceSetMessages{branch: branch, logicalChange: logicalChange}, nil
}

// buildReview creates Ellen's local test promise, Dave's local adoption promise,
// and a review-thread reference set that groups them with the merge subject.
func (builder *Builder) buildReview(previousReviewThreadCID, mergeSnapshotCID cidlib.Cid) (reviewMessages, error) {
	targets := []any{graph.ObjectRow("snapshot", mergeSnapshotCID)}
	testStatement, testErr := graph.StoreMessage(builder.CAS, []graph.Parent{{Role: "review_of", CID: mergeSnapshotCID}}, graph.Payload{
		Promiser:    "ellen",
		Promisee:    "dave",
		PromiseKind: "review_statement",
		PromiseBody: graph.ReviewStatementBody("local_test", targets, "I promise I ran the deterministic local checkout test for this exact merge snapshot.", "test_kept", nil),
		LocalConstraints: map[string]any{
			"authority": "local only",
		},
	})
	if testErr != nil {
		return reviewMessages{}, testErr
	}
	adoptionStatement, adoptionErr := graph.StoreMessage(builder.CAS, []graph.Parent{{Role: "review_of", CID: mergeSnapshotCID}, {Role: "responds_to", CID: testStatement.CID}}, graph.Payload{
		Promiser:    "dave",
		Promisee:    "alice",
		PromiseKind: "review_statement",
		PromiseBody: graph.ReviewStatementBody("local_adoption", targets, "I promise to adopt this merge in my local branch because Ellen's test promise and my own review are sufficient for me.", "accepted_locally", []any{graph.ObjectRow("test", testStatement.CID)}),
		LocalConstraints: map[string]any{
			"authority": "local only",
		},
	})
	if adoptionErr != nil {
		return reviewMessages{}, adoptionErr
	}
	thread, threadErr := builder.storeReferenceSet("ellen", "dave", []graph.Parent{{Role: "previous_reference_set", CID: previousReviewThreadCID}, {Role: "review_of", CID: mergeSnapshotCID}}, "review:rename-copy-merge", "review_thread", "project:poc18-demo", []any{
		graph.ReferenceEntry("subject", []any{graph.Target("snapshot", mergeSnapshotCID)}, nil),
		graph.ReferenceEntry("test", []any{graph.Target("review_statement", testStatement.CID)}, []any{"local deterministic test promise"}),
		graph.ReferenceEntry("adoption", []any{graph.Target("review_statement", adoptionStatement.CID)}, []any{"Dave's local adoption promise"}),
	}, []any{"I promise this review thread groups local review, test, and adoption statements."})
	if threadErr != nil {
		return reviewMessages{}, threadErr
	}
	return reviewMessages{testStatement: testStatement, adoptionStatement: adoptionStatement, thread: thread}, nil
}

// storeReferenceSet writes one signed reference_set promise using the shared
// graph package so scenario fixtures and CLI/core code use the same envelope.
func (builder *Builder) storeReferenceSet(promiser, promisee string, parents []graph.Parent, identity, role, namespace string, entries []any, terms any) (graph.StoredMessage, error) {
	return graph.StoreMessage(builder.CAS, parents, graph.Payload{
		Promiser:    promiser,
		Promisee:    promisee,
		PromiseKind: "reference_set",
		PromiseBody: graph.ReferenceSetBody(identity, role, namespace, entries, terms),
	})
}

// storeSnapshot writes one signed snapshot promise with explicit parent
// snapshot links in both envelope parents and the snapshot body.
func (builder *Builder) storeSnapshot(promiser, promisee string, parents []graph.Parent, identity string, rootRefCID cidlib.Cid, parentSnapshots []cidlib.Cid, summary string, terms any) (graph.StoredMessage, error) {
	return graph.StoreMessage(builder.CAS, parents, graph.Payload{
		Promiser:    promiser,
		Promisee:    promisee,
		PromiseKind: "snapshot",
		PromiseBody: graph.SnapshotBody(identity, rootRefCID, parentSnapshots, summary, terms),
	})
}

// referenceSetTargets parses a local reference_set message and returns the first
// target for each label. POC18 reference sets can be multi-target, but these
// deterministic fixtures use one target per directory label.
func (builder *Builder) referenceSetTargets(refCID cidlib.Cid) (map[string]target, error) {
	content, _, getErr := builder.CAS.Get(refCID)
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
	if kind != "reference_set" {
		return nil, fmt.Errorf("expected reference_set, got %s", kind)
	}
	body, bodyErr := view.PayloadBody()
	if bodyErr != nil {
		return nil, bodyErr
	}
	if len(body) != 5 {
		return nil, fmt.Errorf("reference_set body must have five slots")
	}
	entries, ok := body[3].([]any)
	if !ok {
		return nil, fmt.Errorf("reference_set entries must be array")
	}
	targets := map[string]target{}
	for _, entryValue := range entries {
		entry, entryOK := entryValue.([]any)
		if !entryOK || len(entry) != 3 {
			return nil, fmt.Errorf("reference_set entry must have three slots")
		}
		label, labelOK := entry[0].(string)
		if !labelOK {
			return nil, fmt.Errorf("reference_set label must be text")
		}
		targetRows, rowsOK := entry[1].([]any)
		if !rowsOK || len(targetRows) == 0 {
			return nil, fmt.Errorf("reference_set target list must be non-empty")
		}
		targetRow, rowOK := targetRows[0].([]any)
		if !rowOK || len(targetRow) != 2 {
			return nil, fmt.Errorf("reference_set target row must have two slots")
		}
		role, roleOK := targetRow[0].(string)
		if !roleOK {
			return nil, fmt.Errorf("reference_set target role must be text")
		}
		targetCID, cidErr := store.CIDFromLinkTag(targetRow[1])
		if cidErr != nil {
			return nil, cidErr
		}
		targets[label] = target{role: role, cid: targetCID}
	}
	return targets, nil
}

// directoryReferenceCID resolves a POSIX directory node to the reference_set CID
// that holds the directory's labels.
func (builder *Builder) directoryReferenceCID(nodeCID cidlib.Cid) (cidlib.Cid, error) {
	content, _, getErr := builder.CAS.Get(nodeCID)
	if getErr != nil {
		return cidlib.Undef, getErr
	}
	view, parseErr := graph.ParseEnvelope(content)
	if parseErr != nil {
		return cidlib.Undef, parseErr
	}
	kind, kindErr := view.PayloadKind()
	if kindErr != nil {
		return cidlib.Undef, kindErr
	}
	if kind != "posix_node" {
		return cidlib.Undef, fmt.Errorf("expected posix_node, got %s", kind)
	}
	body, bodyErr := view.PayloadBody()
	if bodyErr != nil {
		return cidlib.Undef, bodyErr
	}
	if len(body) != 5 {
		return cidlib.Undef, fmt.Errorf("posix_node body must have five slots")
	}
	nodeType, typeOK := body[1].(string)
	if !typeOK || nodeType != "directory" {
		return cidlib.Undef, fmt.Errorf("expected directory node")
	}
	return store.CIDFromLinkTag(body[2])
}

// targetCID returns the node CID for a directory label and rejects unexpected
// target roles so scenario mistakes do not silently create false lineage.
func targetCID(targets map[string]target, label string) (cidlib.Cid, error) {
	found, ok := targets[label]
	if !ok {
		return cidlib.Undef, fmt.Errorf("missing reference-set label %s", label)
	}
	if found.role != "node" {
		return cidlib.Undef, fmt.Errorf("label %s role = %s, want node", label, found.role)
	}
	return found.cid, nil
}
