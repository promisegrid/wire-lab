# TE-hikar: POC18 grid CLI command model

TE ID: TE-hikar

## Status

decided

## Decision Under Test

What should the first POC18 `grid` CLI look like so it is familiar enough for
Git/Jujutsu users, but still matches PromiseGrid's message, promise, sparse CAS,
and versioned reference-set model?

## Assumptions

- POC18 is a POC16 protocol/CAS/kernel superset, not a POC17 runtime superset.
- The user-facing command is `grid`.
- Non-production deterministic fixture commands use `poc-*`.
- Production-shaped local daemons use `grid-*`.
- Documentation uses `CLI/core`, not porcelain/plumbing, except when explaining
  the Git analogy.
- Native synchronization is continuous peer DAG sync, not Git push/pull as the
  native collaboration model.
- The first implementation slice is deterministic core graph work: sparse CAS,
  Rabin chunks, POSIX nodes, reference sets, snapshots/change-sets, checkout, and
  promise-shaped retrieval seams.

## Alternatives

### Alt A: Git-compatible verbs first

Use commands such as `grid add`, `grid commit`, `grid branch`, `grid tag`,
`grid push`, and `grid pull`.

This is familiar. It lowers initial explanation cost for Git users. It is also
dangerous because these verbs carry Git's staging index, remote, mutable branch,
and optional tag assumptions. In POC18, a tag, branch, directory, logical change,
review thread, and workspace root are all versioned reference-set promises. A
native `push` or `pull` would suggest a remote authority or explicit transfer
moment even though native PromiseGrid sync is continuous and peer-relative.

### Alt B: Jujutsu-like verbs first

Use commands such as `grid snapshot`, `grid workspace`, `grid change`,
`grid refs`, `grid log`, `grid operation`, and `grid sync`.

This better matches modern distributed version-control lessons. Jujutsu's
working-copy snapshots, operation log, bookmarks, workspaces, and stable change
identity map more cleanly to POC18 than Git's index and branch vocabulary. The
risk is that copying Jujutsu too literally could reintroduce a raw change-ID
field or a repo-centric operation log instead of PromiseGrid reference sets and
message parent links.

### Alt C: PromiseGrid-native verbs only

Use commands such as `grid promise`, `grid refset`, `grid object`,
`grid materialize`, `grid retain`, and `grid advertise`.

This is precise, but it makes the first developer experience unnecessarily
abstract. It exposes protocol internals before users have a mental model for
ordinary work: initialize a workspace, store content, create a snapshot, inspect
references, and materialize a view.

### Alt D: Balanced native CLI

Use a small native command surface:

```text
grid init
grid ingest
grid snapshot
grid checkout
grid refs
grid diag
```

Then add future command groups only when the underlying objects are real:

```text
grid sync
grid change
grid review
grid merge
grid git import
grid git export
grid git pull
grid git push
```

`grid ingest` reads workspace bytes into sparse CAS and graph objects. `grid
snapshot` names a root directory reference set plus parent snapshot links.
`grid checkout` materializes a snapshot or workspace reference set locally.
`grid refs` inspects versioned reference sets. `grid diag` renders exact CBOR
messages for debugging.

Alt D keeps familiar nouns where they help (`checkout`, `snapshot`, `refs`) and
avoids misleading native `push`/`pull` until the Git bridge exists as a bridge,
not as the native sync model.

## Scenario Analysis

### Scenario 1: Alice starts a new workspace

With Alt A, Alice expects `add` and `commit`. That implies a staging index and a
commit object. POC18 has neither as its root abstraction. It has workspace bytes,
Rabin chunks, POSIX nodes, reference sets, and snapshot/change-set messages.

With Alt D, Alice can run `grid init`, then `grid ingest`, then `grid snapshot`.
The first slice may implement `ingest` as a scan that also emits a snapshot for
deterministic POC simplicity, but the vocabulary still leaves room to separate
ingest and snapshot once persistent workspace config exists.

### Scenario 2: Bob receives a sparse branch promise

Git-like `pull` suggests Bob asks a remote for all required objects. In POC18,
Bob may receive only a branch-role reference-set promise and then discover he is
missing chunks. He asks peers for missing objects through promise-shaped sync
interest. `grid refs` and future `grid sync` better describe this: inspect the
reference set, then continuously sync missing DAG pieces from peers Bob locally
trusts.

### Scenario 3: Carol revises one logical change repeatedly

Git-like `commit --amend` or Jujutsu-style raw change IDs are close but not
quite right. POC18 uses a `logical_change` reference set whose versions point at
exact snapshot CIDs. Future `grid change` should operate on that reference set.
The first slice should not invent a raw change-ID field.

### Scenario 4: Dave wants Git interoperability

Git command names are appropriate under an explicit Git bridge group:

```text
grid git import
grid git export
grid git pull
grid git push
```

Those commands should share one conversion core. They must not imply that Git
remotes, hidden refs, appviews, or forges are native PromiseGrid authorities.

### Scenario 5: Ellen diagnoses a protocol message

`grid diag` is useful for developers, but it is not a production authority. It
renders exact CBOR and CID relationships so humans can inspect whether pCID,
parent links, payload, and proof match the spec. This command belongs in the
first slice because POC18 is explicitly testing raw message shape and CAS paths.

## Conclusion

Choose Alt D for POC18.

The first implementation locks `grid` as the user-facing CLI name and exposes a
minimal balanced command set. It does not copy Git's staging or native push/pull
model. It borrows Jujutsu's respect for snapshots, workspaces, and change
evolution, but maps logical change identity to versioned reference sets rather
than raw change IDs. It keeps native PromiseGrid sync peer-relative and
promise-shaped.

## Implications For Open TODOs

- `nahop.2`: command, package, runtime-path, and terminology decisions are now
  lockable.
- `nahop.4` through `nahop.10`: first implementation should build core graph
  behavior before peer transport.
- `nahop.15`: Git push/pull remains bridge work under `grid git`, not native
  PromiseGrid synchronization.
- `nahop.19`: native sync remains continuous peer DAG sync and should not be
  exposed as Git-like push/pull.

## Decision Status

Locked by `DI-harih` in
`protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md`.
