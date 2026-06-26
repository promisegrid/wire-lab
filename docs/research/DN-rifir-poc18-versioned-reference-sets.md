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
"README.md" -> file-version CID
"src"       -> directory-version CID
"LICENSE"  -> file-version CID
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
target. POC18 maps a filename label to a file-version or directory-version CID.

This matters for renames:

```text
old directory:
  "README.md" -> file-version F1

new directory:
  "docs" -> directory D2

D2:
  "intro.md" -> file-version F1 or F2
```

If the content also changed, `intro.md` points to `F2`, and `F2` has parent
`F1`. The file lineage remains intact while directory labels change.

This avoids saying a file "is named README.md." More precisely:

> A directory reference-set promise says that the label `README.md` currently
> points at this file-version or directory-version CID.

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
`push`/`pull` can remain as compatibility UI for import/export, but it should not
be the PromiseGrid-native mental model. Source: `DI-dibut`.

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

## Files and snapshots

POC18 should use both file lineage and snapshot/change-set lineage.

File-version envelopes preserve logical file identity across edits, renames, and
copies. Snapshot/change-set envelopes preserve Git-like project history, atomic
multi-file changes, merge parents, branch targets, and exportability to Git.

This gives POC18:

- better rename/copy history than plain Git;
- Git-like project snapshots for interoperability;
- a clear place for review and test promises;
- a clear place for local adoption decisions.

## Git roundtrip

POC18 should start PromiseGrid-native, but it is not complete until it can import
and export a real Git repository while preserving:

- file content;
- directory structure;
- branch and tag targets;
- parent DAG semantics.

The exported Git commit hashes do not need to match original hashes unless the
objects are byte-identical and the implementation can preserve them naturally.
The important first target is content and DAG fidelity.

Git roundtrip is compatibility work. It should not force POC18 to preserve Git's
remote-centric push/pull model internally. Source: `DI-dibut`.

## Tangled prior-art question

POC18 should review Tangled as adjacent prior art before implementation locks
its collaboration details. Tangled matters because it is a contemporary
decentralized code-hosting project with Git compatibility, self-hosted
infrastructure, social coding goals, AT Protocol identity, round-based pull
request ideas, and Jujutsu change-ID influence.

The question for POC18 is not "should PromiseGrid copy Tangled?" The question is:

> What lessons, if any, should POC18 learn from Tangled's split between hosted
> Git repositories, self-hosted infrastructure, social/appview aggregation,
> identity, review workflow, and Jujutsu-compatible change tracking?

Initial expectations:

- Tangled may teach useful lessons about migration from existing Git workflows.
- Tangled may teach useful lessons about social-code UX and round-based reviews.
- Tangled's knots/appview split may be useful contrast for PromiseGrid's local
  trust and no-global-authority model.
- Tangled's continued SSH/Git push/pull compatibility is likely compatibility
  pressure, not the native PromiseGrid sync model.
- Tangled's Jujutsu change-ID use should be compared against POC18 logical-change
  reference sets.

Source: `DI-dibut`.

## Chunking

POC18 should use restic's content-defined chunker library for file content, but
not adopt restic's repository format as the CAS backend.

The CAS should contain PromiseGrid objects:

- raw chunks;
- Merkle manifests over chunk CIDs;
- file-version envelopes;
- directory/reference-set envelopes;
- snapshot/change-set envelopes;
- review/release/logical-change/workspace reference sets.

This keeps restic as a proven chunking component while preserving PromiseGrid's
CID, pCID, CBOR, parent-link, and promise semantics.

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
> at this file-version CID.

That promise is useful only to agents that recognize and trust the promiser or
the group namespace enough to use it.

## Design consequences

- Use **versioned reference set** as the root abstraction in POC18 docs and code.
- Treat `branch`, `directory`, `release`, `logical_change`, `review_thread`, and
  `workspace` as roles.
- Let roles define validation. Same primitive does not mean same rules.
- Sync starts from reference-set promises.
- Native collaboration is continuous peer DAG synchronization, not Git-style
  push/pull.
- Filenames are directory/reference-set labels.
- Logical-change reference sets replace Jujutsu-style intrinsic change IDs.
- File history owns content lineage, not names.
- Snapshot/change-set history owns atomic project state.
- CAS stores exact bytes and exact promise objects by CID.

## Links

- TODO: `protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md`
- Decision: `DI-zuruj`
- Sync/prior-art decision: `DI-dibut`
- Prior decision refined by this note: `DI-vilum`
