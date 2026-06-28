# TE-kopap: POC18 Git bridge versus native CAS architecture

## TE ID

TE-kopap

## Status

decided

## Decision under test

POC18 needs to decide which architecture should guide implementation:

- Git bridge compatibility as the primary model.
- PromiseGrid-native CAS/version graph only.
- A hybrid where PromiseGrid-native reference-set promises are canonical and Git
  import/export/push/pull are required bridge behavior.

This TE tests those alternatives against the POC18 goal: replace Git/GitHub's
repository-centered collaboration model with sparse CAS, CID-addressed version
objects, grid-envelope parent links, versioned reference-set promises, local
trust, and conventional Git interoperability where needed.

## Existing locked inputs

- POC18 is a protocol/CAS/kernel superset of POC16, not a POC17 runtime
  continuation. Source: `DI-zuruj`.
- Versioned reference sets are the root abstraction for directories, filenames,
  branches, tags, releases, logical changes, review threads, workspace roots,
  and fetch/discovery entrypoints. Source: `DI-zuruj`.
- Native collaboration is continuous peer DAG synchronization, not Git-style
  push/pull as the native mental model. Source: `DI-dibut`.
- POC18 uses Rabin content-defined chunking for all file content and stores
  chunks and manifests as PromiseGrid CAS objects. Source: `DI-dofoj`.
- POC18 must provide conventional Git bridge behavior for import, export, push,
  and pull through shared conversion code. Source: `DI-dofoj`.
- POC18 must represent every POSIX inode type, including regular files,
  directories, symbolic links, hard links, character devices, block devices,
  FIFOs, and sockets. Source: `DI-radaj`.

## Assumptions

- A pCID is a protocol-spec selector, not a repository name, branch name, file
  path, agent address, operation code, or remote identity.
- Printable CIDs are CIDv1 base32; binary CIDs remain binary on the wire.
- Grid envelopes are first-class CAS objects and may carry parent links that
  connect version history.
- No agent can make another agent's promise. A branch, tag, review, release,
  storage, retention, or merge statement is the promiser's local promise, not a
  global authority.
- Each agent's CAS is partial. There is no complete global CAS and no complete
  global monitor.
- Git interoperability is adoption-critical, but Git's repository, remote,
  branch, tag, and forge assumptions are not automatically PromiseGrid
  assumptions.
- POC18 should be useful for small text repositories, large binary files,
  generated artifacts, reviews, releases, and mixed POSIX workspaces.

## Alternatives

### Alt A: Git bridge compatibility first

POC18 would treat conventional Git repositories and Git operations as the
primary architecture. PromiseGrid CAS objects would mostly mirror Git objects,
and native PromiseGrid behavior would be shaped around importing, exporting,
pushing, and pulling Git refs.

This makes migration easy because Alice can start with an existing Git
repository and Bob can continue using familiar Git tooling. It also lowers the
first-demo burden because Git already has commit, tree, blob, tag, and remote
semantics.

It also makes Git the gravitational center. Git cannot naturally represent all
POSIX inode types, it treats large files as awkward out-of-band concerns, and it
keeps branches/tags/refs tied to Git's remote-centered vocabulary. If POC18 is
implemented this way, PromiseGrid risks becoming a Git adapter rather than a
PromiseGrid-native collaboration substrate.

### Alt B: PromiseGrid-native CAS/version graph only

POC18 would ignore conventional Git bridge behavior at first. It would define
only PromiseGrid-native objects: Rabin chunks, manifests, POSIX node versions,
directory/reference-set versions, snapshots, logical changes, reviews, releases,
retention promises, and continuous peer DAG synchronization.

This is architecturally clean. Alice and Bob can exchange reference-set
promises without pretending there is a central repository or remote. The design
can represent POSIX nodes directly and can make large files native instead of
special.

The cost is adoption and verification. Existing Git repositories are the real
input corpus. Without import/export/push/pull, POC18 cannot prove that it can
meet developers where they already are. It would also remove one of the best
test oracles: roundtripping through Git exposes content, DAG, rename, mode, and
branch mapping bugs early.

### Alt C: Hybrid native core with required Git bridge

POC18 would make PromiseGrid-native reference-set promises canonical. The native
model is sparse CAS, Rabin chunks, POSIX node versions, snapshots, logical
changes, reviews, and continuous peer DAG sync. Git import/export/push/pull are
required bridge surfaces implemented through shared conversion code, not the
native authority model.

This keeps PromiseGrid's vocabulary honest. A Git branch maps into a branch-role
reference-set promise. A Git tree maps into directory reference sets and POSIX
node versions. Git commits map into snapshot/change-set objects. Git remotes are
bridge endpoints; they are not global authorities. Unsupported POSIX node types
produce explicit mapping, loss, refusal, or local non-commitment records.

The cost is implementation discipline. POC18 must define the native object model
before the bridge can be correct, but it must also run Git roundtrip checks early
enough to prevent the native model from drifting into something impossible to
interop with. The bridge is mandatory compatibility work, not optional polish.

## Scenario analysis

### Scenario 1: Normal single-author edit

Alice edits a small text file and a large binary file in one workspace.

- Alt A can store the text file naturally, but the large binary pressures Git's
  normal blob model and tends to reintroduce Git LFS-style out-of-band thinking.
- Alt B stores both files as Rabin chunk manifests and CAS objects, which is the
  native POC18 goal, but cannot prove the text-file path still roundtrips with
  existing Git tools.
- Alt C stores both files natively and then proves bridge behavior by exporting
  the text-compatible parts to Git while preserving native large-file structure
  in PromiseGrid CAS.

Alt C gives the best long-term shape: native large-file handling without losing
Git compatibility pressure.

### Scenario 2: Import an existing Git repository

Bob has a conventional Git repository with branches, tags, merges, and symlinks.
He wants to import it into PromiseGrid, continue working natively, and later
push compatible state back to a Git remote.

- Alt A makes the first import straightforward but leaves Bob trapped in Git's
  conceptual model.
- Alt B cannot satisfy the push/pull requirement and gives Bob no clean migration
  path.
- Alt C maps Git refs, commits, trees, blobs, tags, and modes into native
  reference sets, snapshots, node versions, manifests, and mapping records. The
  same conversion core supports import and pull in one direction and export and
  push in the other.

Alt C is the only alternative that satisfies both migration and native evolution.

### Scenario 3: POSIX workspace beyond Git's model

Carol versions a workspace containing regular files, directories, symlinks,
hard links, a FIFO placeholder, and device-node metadata.

- Alt A either drops information, abuses Git blobs, or invents Git-specific side
  files. That makes Git the schema authority and hides loss.
- Alt B can represent each POSIX node type directly, including local
  materialization constraints.
- Alt C can represent each POSIX node type directly and can also state exactly
  what a Git bridge can or cannot preserve.

Alt C keeps the native model complete and makes Git limitations explicit instead
of silently corrupting the workspace.

### Scenario 4: Review and logical change identity

Dave sends Ellen a proposed change. He later revises it twice while preserving
the logical identity of the change. Ellen reviews each round, attaches test
results, and promises which version she locally accepts.

- Alt A tends to model this as pull requests, hidden refs, forge metadata, or
  Jujutsu-compatible Git overlays.
- Alt B models it cleanly as logical-change and review-thread reference sets,
  but has no Git-facing story for teams that still need roundtrips.
- Alt C models Dave's change as a logical-change reference set and Ellen's
  review as a review-thread reference set, while bridge code can materialize
  comparable Git branches or refs when needed.

Alt C learns from Git/Jujutsu/Tangled without adopting their native authority
model.

### Scenario 5: Sparse peers and missing objects

Frank receives a branch-role reference set from Alice but lacks several target
CAS objects. Grace has some chunks, Bob has the snapshot object, and Alice has
the root reference set.

- Alt A naturally asks a Git remote for missing objects, which re-centers the
  remote as an authority-like source.
- Alt B naturally lets Frank ask peers for missing CIDs and retain only what he
  locally chooses to keep.
- Alt C uses the native sparse CAS behavior for PromiseGrid peers and uses Git
  fetch only at bridge boundaries.

Alt C preserves sparse peer-relative storage while still allowing Git bridge
retrieval from conventional remotes.

### Scenario 6: Failure, corruption, and incomplete writes

Mallory sends Alice a malformed object, an object whose bytes do not match its
CID, or a reference set pointing at unavailable targets.

- Alt A inherits Git object verification for Git objects, but bridge-specific
  side metadata can become ambiguous.
- Alt B rejects mismatched CID bytes and records local non-commitment for missing
  targets.
- Alt C does the same native CID verification and additionally tests whether
  imported/exported Git objects preserve the same content and parent DAG.

Alt C gives two useful checks: native CID verification and Git roundtrip
verification.

### Scenario 7: Long-horizon migration

In ten years, Alice's team still uses some Git tools, Bob's team uses
PromiseGrid-native sync, and Carol's archival agent retains only release
reference sets and paid storage objects.

- Alt A keeps everyone bound to Git's long-term object and workflow limits.
- Alt B may be technically cleaner, but it cuts off Git users too early.
- Alt C allows migration by degrees: Git remains an adapter while native
  reference-set promises become the durable model.

Alt C is the best migration path because it does not require a flag day.

### Scenario 8: Trust-boundary and forge replacement

Ellen trusts Alice's release promise but not Mallory's review summary. Frank
trusts Bob's storage promises but not Bob's merge choices. No one wants a GitHub
or appview server to decide truth.

- Alt A tends to preserve forge-shaped concepts such as central remotes,
  permission-like merge state, and hosted review authority.
- Alt B keeps all trust local, but lacks the bridge path needed for incremental
  adoption.
- Alt C keeps merge, release, review, storage, and retention as local promises
  while exposing Git-compatible views only as bridge materialization.

Alt C best matches Promise Theory while still being practical.

### Scenario 9: Scale effects

Grace's workspace has many small files, a few multi-gigabyte assets, generated
artifacts, and frequent LLM-produced review iterations.

- Alt A scales poorly for large in-band files and can make frequent generated
  artifacts feel like repository bloat.
- Alt B scales better because Rabin chunks, sparse CAS, retention promises, and
  reference-set sync can avoid fetching everything.
- Alt C keeps those native scale advantages and uses Git bridge behavior only
  where conventional tooling needs it.

Alt C preserves the scale behavior POC18 is meant to test.

## Comparison summary

| Criterion | Alt A: Git-first | Alt B: native-only | Alt C: hybrid |
|---|---|---|---|
| Git migration | Strong | Weak | Strong |
| PromiseGrid-native semantics | Weak | Strong | Strong |
| Large files in-band | Weak | Strong | Strong |
| POSIX inode coverage | Weak | Strong | Strong |
| Local trust | Medium | Strong | Strong |
| Sparse CAS | Medium | Strong | Strong |
| Conventional push/pull | Strong | Weak | Strong as bridge |
| Risk of stale Git vocabulary leak | High | Low | Medium, manageable |
| First implementation complexity | Medium | Medium | High |
| Long-term architecture fit | Weak | Medium | Strong |

## Rejected alternatives

- Reject Alt A as the architecture baseline. It is useful bridge pressure, but
  it would make Git concepts the native shape and would undermine POC18's
  reference-set, sparse-CAS, large-file, and local-trust goals.
- Reject Alt B as the complete POC18 plan. It is the cleanest native model, but
  it fails the explicit Git import/export/push/pull requirement and makes
  adoption unnecessarily hard.

## Surviving alternative

Alt C survives: POC18 should implement a PromiseGrid-native core with required
Git bridge compatibility.

The native core is the source of truth:

- Rabin chunks and manifests for all file content.
- POSIX node-version objects for every POSIX inode type.
- Directory/reference-set versions for names, tags, branches, releases, reviews,
  logical changes, and workspaces.
- Snapshot/change-set objects with parent links.
- Continuous peer DAG sync through voluntary promises and sparse CAS.

The Git bridge is a required adapter:

- Import and pull read Git refs and objects into native PromiseGrid objects.
- Export and push materialize compatible native state back into Git refs and
  objects.
- One shared conversion core handles both directions.
- Unsupported or lossy mappings are explicit, not silent.

## Decision status

No new DF question is needed for the high-level architecture. Existing DIs
already lock the hybrid direction: `DI-zuruj` locks native reference sets,
`DI-dibut` locks continuous peer DAG sync, `DI-dofoj` locks Git bridge
import/export/push/pull, and `DI-radaj` locks full POSIX inode coverage.

Remaining DF work is implementation-level:

- Choose Git and Rabin chunking libraries.
- Lock implementation paths, command names, package names, and runtime path
  patterns.
- Decide whether the first POC18 slice includes LLM agents or stays
  deterministic.

## Implications for open TODOs and pending DIs

- `nahop.1` is complete after this TE.
- `nahop.2` remains the next blocker before code scaffolding because path,
  package, command, runtime, and generated-CAS patterns still need to be locked.
- `nahop.3` should write the RFC-like POC18 pCID spec around Alt C, not around
  Git-first or native-only framing.
- `nahop.15` remains required bridge work, but it should not be allowed to
  distort native object semantics.
- `nahop.19` remains required native sync work, and it should not be implemented
  as Git push/pull with different names.
- Analyzer gates in `nahop.20` should verify both halves: native reference-set
  behavior and Git bridge roundtrip behavior.
