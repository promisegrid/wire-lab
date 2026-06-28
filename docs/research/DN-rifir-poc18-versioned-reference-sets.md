# POC18 versioned reference sets

This design note records the POC18 tag and reference-set decision in plain
English. It is not a frozen protocol spec, not a TE, and not final SDK prose. It
is the rationale behind `TODO-nahop` and the current design synthesis for using
PromiseGrid CAS as a Git/GitHub replacement substrate. Source: `DI-zuruj`.

## Short version

POC18 should not treat tags as optional metadata like Git's `--tags` behavior,
and it should not treat `push` and `pull` as the native collaboration model. It
should treat **versioned reference sets** as the primary way agents discover,
sync, and fetch anything. Source: `DI-zuruj`; `DI-dibut`.

A versioned reference set is:

```text
namespace/context + labeled entries -> target CID(s)
```

with parent links, promiser identity, promise terms, and proof. A directory, a
branch, a release, a logical change, a review thread, and a workspace are all
role-specific reference sets.

## Why Git's split is not enough

Git has excellent low-level ideas: content-addressed objects, tree snapshots,
commit parents, cheap branches, and decentralized object exchange. It also has
awkward user-facing splits that POC18 should not copy blindly.

In Git:

- a branch is a moving ref to a commit;
- a tag is usually a separate ref to an object;
- tags are often fetched as optional extras;
- collaboration is typically expressed as explicit `push` and `pull` operations
  against named remotes;
- a tree maps filenames to blobs/subtrees;
- rename history is inferred from snapshots rather than stored as a first-class
  identity relationship.

That works, but it creates separate mechanisms for things that are all really
named references. PromiseGrid can use one cleaner mechanism: signed,
versionable, multi-target reference-set promises.

## What Jujutsu teaches us

Jujutsu separates the exact commit ID from the logical change ID. A commit ID
identifies the exact commit object. A change ID identifies the evolving logical
change across rewrites, amendments, and rebases.

That is the important lesson: users and reviewers need stable logical identity
for evolving work, while storage still needs immutable exact object identity.

POC18 should not copy Jujutsu's change ID literally as a random field inside a
commit-like object. A PromiseGrid-native logical change should be a versioned
reference-set promise. That gives the logical change its own CID, parent history,
promiser, promise terms, and proof.

## Why tags and directories converge

The key observation is that a multi-target tag with labels is directory-shaped.

A directory is:

```text
"README.md" -> regular-file node-version CID
"src"       -> directory-version CID
"LICENSE"  -> regular-file node-version CID
"latest"   -> symbolic-link node-version CID
```

A rich tag or reference bundle is:

```text
"current"      -> snapshot CID
"review"       -> review-thread CID
"tests"        -> test-result CID
"exported_git" -> git-commit CID
```

Both are scoped labeled references to CIDs. The difference is role and
validation, not mechanism.

That means POC18 should use **versioned reference sets** as the shared primitive:

- directory entries are reference-set labels;
- branch heads are reference-set labels;
- release artifacts are reference-set labels;
- logical-change versions are reference-set labels;
- review-thread entries are reference-set labels;
- workspace roots and local overrides are reference-set labels.

## Filenames are directory history

A file's name should not be a property of file history. It should be a property
of directory/reference-set history.

The Unix analogy is useful: a directory entry maps a filename to an inode-like
target. POC18 maps a filename label to a POSIX node-version or
directory-version CID.

This matters for renames:

```text
old directory:
  "README.md" -> regular-file node-version F1

new directory:
  "docs" -> directory D2

D2:
  "intro.md" -> regular-file node-version F1 or F2
```

If the content also changed, `intro.md` points to `F2`, and `F2` has parent
`F1`. The file lineage remains intact while directory labels change.

This avoids saying a file "is named README.md." More precisely:

> A directory reference-set promise says that the label `README.md` currently
> points at this POSIX node-version or directory-version CID.

## POSIX node types

POC18 should support every POSIX inode type: regular file, directory, symbolic
link, hard link, character device, block device, FIFO, and socket. Source:
`DI-radaj`.

Regular files carry Rabin chunk manifests. Directories are reference sets.
Symbolic links are node promises whose payload preserves the link target bytes.
Hard links are multiple directory labels pointing at the same node identity or
link-group target, not duplicated content. Character devices and block devices
preserve type plus major/minor metadata. FIFOs and sockets preserve node type
and materialization constraints; POC18 does not pretend to store live stream
contents or kernel socket state in CAS.

Materialization is local. Alice may promise to materialize a device, FIFO, or
socket node on her machine only if local capabilities and policy allow it. If a
host cannot or should not create a node, it should record an explicit local
non-commitment, adaptation, or refusal rather than silently degrading the object
to a regular file.

## Branches are reference-set roles

A branch should not be a separate primitive. A branch is a reference set whose
role is `branch`.

For example:

```text
reference_set role=branch name=main:
  "head" -> snapshot/change-set CID
  "review" -> review-thread CID
  "latest_test" -> test-result CID
```

Moving the branch means publishing a new version of the reference set with a
parent link to the prior branch reference-set version. The old position remains
accessible by CID.

This is stronger than a Git branch ref because the movement itself is a signed,
versioned PromiseGrid object.

## Reference sets are sync entrypoints

In POC18, tags/reference sets are not a secondary thing to fetch after commits.
They are the way Alice starts synchronizing.

The flow is:

```text
Alice and Bob continuously exchange reference-set promises.
Alice locally judges which promises she trusts.
Alice decides which labeled target CIDs she wants.
Alice asks peers for target CIDs and parent-linked CAS objects she is missing.
Alice advertises the reference sets and CAS objects she is willing to share.
```

This works for branches, directories, releases, logical changes, reviews, and
workspaces. It also fits sparse CAS: Alice may know a reference set before she
has all target objects.

Native POC18 sync is therefore not "Alice pushes to Bob" or "Bob pulls from
Alice." It is continuous peer DAG reconciliation among agents who each decide
what to advertise, request, verify, retain, forward, or ignore. Git-style
`push`/`pull` remains required Git bridge behavior, but it should not be the
PromiseGrid-native mental model. Source: `DI-dibut`; `DI-dofoj`.

## Logical changes replace change IDs

For POC18, a logical change is a reference set with role `logical_change`.

Example:

```text
reference_set role=logical_change:
  "current" -> snapshot/change-set CID
  "previous" -> older snapshot/change-set CID
  "review" -> review-thread CID
  "tests" -> test-result CID
```

This covers most of what Jujutsu change IDs provide: stable identity for an
evolving change. It also adds PromiseGrid properties Jujutsu's intrinsic field
does not provide by itself:

- the logical change has a CID;
- the logical change has version history;
- the logical change is a promise object;
- the logical change can have multiple labeled targets;
- the logical change can be signed and judged locally;
- divergent current targets can be represented honestly.

## POSIX nodes and snapshots

POC18 should use both POSIX node lineage and snapshot/change-set lineage.

Node-version envelopes preserve logical node identity across edits, renames, hard
links, metadata changes, and copies. Snapshot/change-set envelopes preserve
Git-like project history, atomic multi-node changes, merge parents, branch
targets, and exportability to Git where Git can represent the node types.

This gives POC18:

- better rename/copy history than plain Git;
- Git-like project snapshots for interoperability;
- a clear place for review and test promises;
- a clear place for local adoption decisions.

## Git bridge

POC18 should start PromiseGrid-native, but it is not complete until it can
bridge to conventional Git repositories in both local and remote directions. It
must import from, export to, pull from, and push to real Git repositories while
preserving:

- file content;
- directory structure;
- supported POSIX node types, with explicit mapping or loss records when Git
  cannot represent a node type;
- branch and tag targets;
- parent DAG semantics.

The Git bridge should share one conversion core. Import and pull both read Git
refs and objects into PromiseGrid reference sets, snapshots, manifests, chunks,
and mapping records; export and push both materialize selected PromiseGrid
reference sets and CAS objects back into Git refs and objects. Import/export are
local filesystem edges. Push/pull add Git remote, auth, and transport edges.

The exported or pushed Git commit hashes do not need to match original hashes
unless the objects are byte-identical and the implementation can preserve them
naturally. The important first target is content and DAG fidelity.

Git bridge behavior is compatibility work. It should not force POC18 to preserve
Git's remote-centric push/pull model internally. Native PromiseGrid
collaboration remains continuous peer DAG synchronization over local promises,
reference-set advertisements, object-availability promises, and missing-object
requests. Source: `DI-dibut`; `DI-dofoj`.

## Tangled prior-art question

POC18 reviewed Tangled as adjacent prior art in
`docs/research/DN-dopod-poc18-tangled-prior-art.md`. Tangled matters because it
is a contemporary decentralized code-hosting project with Git compatibility,
self-hosted infrastructure, social coding goals, AT Protocol identity,
round-based pull-request ideas, and Jujutsu change-ID influence.

The question for POC18 is not "should PromiseGrid copy Tangled?" The question is:

> What lessons, if any, should POC18 learn from Tangled's split between hosted
> Git repositories, self-hosted infrastructure, social/appview aggregation,
> identity, review workflow, and Jujutsu-compatible change tracking?

The review conclusion is: POC18 should learn from Tangled's self-hosting,
migration, social-code UX, review rounds, stable logical-change identity, and
conventional Git/SSH interoperability pressure, but explicitly differ from
Tangled by keeping Git push/pull as bridge behavior rather than the native sync
model, and by keeping appview aggregation, role-based access control, hidden Git
refs, and raw Jujutsu change IDs out of the native PromiseGrid model unless a
later TE/DI narrows one of those choices.

Source: `DI-dibut`; `DI-dofoj`; `DN-dopod`.

## Chunking

POC18 should use Rabin content-defined chunking for all file content, but not
adopt restic's repository format as the CAS backend. A restic-derived Go chunker
is acceptable only if the implementation lock verifies it satisfies the Rabin
chunking requirement.

The CAS should contain PromiseGrid objects:

- raw chunks;
- Merkle manifests over chunk CIDs;
- POSIX node-version envelopes;
- directory/reference-set envelopes;
- snapshot/change-set envelopes;
- review/release/logical-change/workspace reference sets.

This keeps restic as a possible proven chunking component while preserving
PromiseGrid's CID, pCID, CBOR, parent-link, and promise semantics. It also means
large files are native in-band CAS content, not a special out-of-band path.
Source: `DI-dofoj`.

## Promise Theory interpretation

A reference set is a promise, not an authority.

When Alice publishes a branch-role reference set, she is promising something like:

> I currently treat this labeled target as my branch head under these local
> terms.

Bob may believe it, ignore it, mirror it, challenge it, or ask for the target
CAS objects. No global service decides that Alice's branch, tag, directory, or
release is authoritative.

This same framing applies to directories:

> I currently treat the label `README.md` in this directory context as pointing
> at this POSIX node-version CID.

That promise is useful only to agents that recognize and trust the promiser or
the group namespace enough to use it.

## Design consequences

- Use **versioned reference set** as the root abstraction in POC18 docs and code.
- Treat `branch`, `directory`, `release`, `logical_change`, `review_thread`, and
  `workspace` as roles.
- Let roles define validation. Same primitive does not mean same rules.
- Sync starts from reference-set promises.
- Native collaboration is continuous peer DAG synchronization, not Git-style
  push/pull; conventional Git push/pull still exists through the Git bridge.
- Filenames are directory/reference-set labels.
- POC18 supports all POSIX inode types as node-version or directory/reference-set
  promises.
- Logical-change reference sets replace Jujutsu-style intrinsic change IDs.
- File history owns content lineage, not names.
- Snapshot/change-set history owns atomic project state.
- CAS stores exact bytes and exact promise objects by CID; all file content uses
  Rabin chunks plus PromiseGrid manifests.

## Links

- TODO: `protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md`
- Decision: `DI-zuruj`
- Sync/prior-art decision: `DI-dibut`
- Git bridge/chunking refinement: `DI-dofoj`
- POSIX inode type refinement: `DI-radaj`
- Prior decision refined by this note: `DI-vilum`
