# TE-vahoj: POC18 superset architecture

## TE ID

TE-vahoj

## Status

needs DF

## Decision under test

POC18 needs one implementation-sequencing and architecture decision before code
scaffolding begins:

> How should POC18 be implemented so it remains a POC16 superset, replaces
> Git/GitHub with PromiseGrid-native CAS/reference-set collaboration, preserves
> conventional Git interoperability, supports DevOps/root-filesystem use, and
> avoids accidentally narrowing the design to a local-only, Git-shaped, or
> RPC-shaped blind alley?

This TE is intentionally large. It supersets the POC18 planning thread that
started with "CAS as object storage and grid envelope parent links as file and
branch version chaining" and then accumulated requirements around Jujutsu,
Tangled, Rabin chunking, POSIX inode coverage, sparse CAS, promise-based
inter-agent messaging, porcelain/plumbing library boundaries, in-band
collaboration, and DevOps ordered replay.

The TE does not by itself authorize implementation. It narrows the alternatives
and identifies the remaining DF questions that must be locked before `nahop.2`
can scaffold code.

## Source corpus

This TE must be read with the following sources:

- `protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md`
- `docs/thought-experiments/TE-kopap-poc18-git-bridge-vs-native-cas.md`
- `docs/research/DN-rifir-poc18-versioned-reference-sets.md`
- `docs/research/DN-dopod-poc18-tangled-prior-art.md`
- `implementations/poc18-cas-git-replacement/docs/protocols/version-control.md`
- `implementations/poc18-cas-git-replacement/docs/turing-equiv.html`
- the current-session POC18 planning entries in `~/.codex/history.jsonl` under
  session `019dd713-8c32-78d3-8330-974286279de6`

The implementation-local paper copy
`implementations/poc18-cas-git-replacement/docs/turing-equiv.html` is the
working source for the DevOps/root-filesystem part of this TE. Its core finding
for POC18 is that self-administered hosts execute tools inside the target
environment they are changing, so changes can affect the tools and foundations
that apply later changes. Therefore test and production hosts need
deterministic, repeatable change order; once a change sequence is validated in
testing, production should replay the same sequence in the same order.

## Session-log coverage

The session log added the following requirements that must not be lost:

| Session item | Requirement carried into this TE |
|---|---|
| `history.jsonl:5918` | POC18 is a Git-replacement POC using CAS as object storage and grid-envelope parent links for version chaining. |
| `history.jsonl:5922` / `5924` | Git versus Jujutsu differences, especially change IDs versus commit IDs, must inform the model. |
| `history.jsonl:5925` | Instead of Jujutsu-style change IDs as raw fields, POC18 should evaluate versionable, multi-target tags/reference sets. |
| `history.jsonl:5926` | Tags/reference sets should be fetch/discovery roots; a branch is a role/use of the same mechanism. |
| `history.jsonl:5931` | POC18 is a POC16 superset, not a POC17 runtime superset; the reference-set reasoning belongs in a design note. |
| `history.jsonl:5933` | Tangled prior art must be reviewed; native PromiseGrid sync should not be explicit Git-style push/pull. |
| `history.jsonl:5947` | Rabin chunking for all files is load-bearing; Git import/export/push/pull are required bridge surfaces and should share conversion code. |
| `history.jsonl:5950` | POC18 must support all POSIX inode types: regular file, directory, symlink, hard link, char device, block device, FIFO, socket. |
| `history.jsonl:5962` / `5963` | Repo helper commands such as `tools/cid` and `tools/spec` do not belong inside pCID-derived protocol specs. |
| `history.jsonl:5966` | Sparse storage per node/agent, active background chunk pulling, CAS shared between repo-like views, and local-first blind-path risk must be tested. |
| `history.jsonl:5967` | All inter-agent messaging must be promise-based, likely with capability tokens; no RPC-shaped design drift. |
| `history.jsonl:5968` | Porcelain/plumbing communication must use the same library as plumbing; evaluate `grid()` messages for UI/backend and remote UI control; compare terminology. |
| `history.jsonl:5970` | POC18 must be usable as a DevOps/root-filesystem tool with patterns, triggers, large repositories, binary chunks, incremental ordered replay, and in-band large blobs. |
| `history.jsonl:5971` | Use the implementation-local copy of the ordering paper now stored at `implementations/poc18-cas-git-replacement/docs/turing-equiv.html`. |

These entries are not separate decisions. They are the input requirements this
TE must explain, compare, and route into follow-up DFs or implementation tasks.

## Existing locked inputs

- POC18 is a protocol/CAS/kernel superset of POC16, not a POC17 runtime
  continuation. Source: `DI-zuruj`.
- Versioned reference sets are the root abstraction for directories, filenames,
  tags, branches, releases, logical changes, review threads, workspace roots,
  and fetch/discovery. Source: `DI-zuruj`.
- Native sync is continuous peer DAG synchronization, not Git-style push/pull as
  the native mental model. Source: `DI-dibut`.
- POC18 uses Rabin content-defined chunking for all file content and stores
  chunks and manifests as PromiseGrid CAS objects. Source: `DI-dofoj`.
- POC18 must provide conventional Git bridge behavior for import, export, push,
  and pull through shared conversion code. Source: `DI-dofoj`.
- POC18 must represent every POSIX inode type. Source: `DI-radaj`.
- POC18 has one implementation-local version-control pCID spec at
  `implementations/poc18-cas-git-replacement/docs/protocols/version-control.md`.
  Source: `DI-lidaj`.

## PromiseGrid invariants

These invariants constrain every alternative:

- pCID selects a protocol spec and parser/builder, not an agent address,
  destination, operation, repository, branch, file path, message type, or
  workflow action.
- The wire/profile shape remains `grid([42(pCID), ...protocol-defined-slots])`.
  POC18's current spec shape is `grid([42(pCID), parents, payload, proof])`.
- All hashes used as object identities are CIDs. CIDs are binary on wire and
  CIDv1 base32 when printable.
- A PromiseGrid object helps an agent make, recognize, remember, and evaluate
  promises. It does not command another agent.
- No agent can promise on behalf of another agent. Branches, releases, reviews,
  merge statements, retention statements, and availability statements are local
  promises by the promiser.
- Trust is local and relationship-relative. There is no global merge authority,
  repository authority, tag authority, review authority, permission service,
  conformance service, global CAS, or global monitor.
- Inter-agent behavior must be promise-shaped. "Fetch", "serve", "store",
  "retain", "forward", "review", "merge", and "apply" are pCID-defined promise
  meanings or local implementation actions, not RPC commands.
- Capability tokens may be used as promises for future storage, retrieval,
  forwarding, compute verification, review, or archival behavior, but the token
  remains the issuer's promise, not a global permission object.

## Alternatives

### Alt A: Deterministic local-first

POC18 would first implement a deterministic single-process or local-process
slice: filesystem scan, Rabin chunking, local CAS, reference sets, snapshots,
materialization, and a CLI/UI. Docker multi-agent sync and LLM agents would come
after the local object model is stable.

This is the fastest route to correct data structures. It keeps the first
debugging loop tight, allows ordinary tests to prove chunking and reference
sets, and reduces network/autonomy noise while the object model is still being
formed.

The risk is serious: if the local slice treats a repo as a complete local
database, assumes all objects are present, or makes fetch/materialization
synchronous local function calls, later sparse multi-agent sync may be bolted on
awkwardly. The local-first path is acceptable only if the first slice is forced
to model missing objects, partial CAS, background pullers, and promise-shaped
retrieval from day one.

### Alt B: Docker sparse-network-first

POC18 would start with multiple agents, each with a partial filesystem CAS,
running in containers. The earliest demo would prove that a reference set can
arrive before all target objects, that missing chunks can be requested from
peers, and that each agent locally decides what to store, forward, retain, or
drop.

This directly attacks the sparse-CAS risk. It prevents a local-complete repo
assumption from hardening into code. It also exercises the POC16 inheritance
more honestly because POC16 already proved pCID-selected parser/builder roles,
filesystem CAS, exact CBOR messages, tokens, encrypted payloads, raw-message
retention, and analyzer gates.

The cost is implementation drag. If the object model is not nailed down first,
the Docker slice may produce network noise before basic files, POSIX node
objects, reference sets, and snapshots are correct. The system may become
harder to refactor because every small object-model correction touches process
boundaries and run orchestration.

### Alt C: LLM/autonomy-first

POC18 would introduce autonomous agents early. Alice, Bob, Carol, Dave, and
Ellen would make local choices about what to advertise, review, retain, fetch,
merge, and trust. A monitor/analyzer would inspect logs for PromiseGrid fitness.

This addresses the reason POC18 exists: Git/GitHub workflows are failing under
LLM-scale collaboration. Early autonomy can surface merge, review, issue,
alert, and chat requirements that a deterministic local implementation might
miss.

The risk is highest. If the protocol bytes and object model are not stable, LLM
behavior can create plausible-looking but shallow artifacts. POC10 through POC13
already showed the danger of LLM text being recorded as evidence without
actually driving valid protocol behavior. POC18 should not repeat that mistake.
LLMs should choose local promise intent only after deterministic code owns the
valid object and message shapes.

### Alt D: Hybrid staged path

POC18 would implement deterministic local storage, plumbing, and UI first, but
the interfaces, tests, and runtime path layout would be sparse and multi-agent
from the beginning. A local run may use one process, but it must behave as if
CAS is partial, objects can be missing, background pullers exist, and retrieval
from another agent is promise-shaped. Docker follows as soon as the local core
can store, materialize, and walk a small graph. LLM agents follow after Docker
proves the sync and collaboration surfaces.

This alternative preserves the velocity of local-first development while
refusing the complete-local-repo assumption. It treats sparse storage as an
interface invariant rather than a later feature. It also gives the Git bridge a
stable native model to convert into and out of.

The cost is discipline. The first slice must deliberately include missing-object
tests, background retrieval abstractions, multi-agent path patterns, and
promise-shaped retrieval APIs even if there is only one local process at first.

### Alt E: Git-first bridge

POC18 would start from Git objects, refs, remotes, and push/pull behavior.
PromiseGrid objects would wrap or mirror Git concepts. Git compatibility would
be excellent from day one.

This path is tempting because it gives immediate real-world inputs and test
oracles. It also meets existing developers where they are.

It is rejected as a native baseline because it lets Git's object model and
forge vocabulary become the gravitational center. Git cannot faithfully model
all POSIX inode types, handles large binaries poorly without out-of-band
extensions, treats tags as optional side data, and carries remote-centered
branch/push/pull assumptions. POC18 can bridge to Git, but it must not become a
Git adapter with PromiseGrid branding.

### Alt F: Native-only PromiseGrid

POC18 would ignore conventional Git bridge work until much later. It would
define only PromiseGrid-native Rabin chunks, manifests, POSIX nodes, reference
sets, snapshots, reviews, releases, and continuous peer sync.

This is architecturally clean. It gives the least vocabulary contamination from
Git and forges.

It is rejected as the complete POC18 plan because Git interop is adoption- and
verification-critical. Existing Git repos are the practical corpus. Import,
export, push, and pull are needed both for migration and for testing whether
POC18's model can roundtrip ordinary content and DAG semantics.

### Alt G: Repos as views over shared CAS

POC18 would retain "repo" as a user-facing convenience: a workspace or project
view over shared CAS/reference-set state. The storage substrate is not a repo
database. One agent/node CAS may store chunks and messages for many repo-like
views; reference sets define which objects matter to a view.

This matches the user's stated concern: chunk storage should be shared between
repos, with one CAS per node or per agent. It also keeps familiar language
available for developers while preventing "repo" from becoming a protocol
authority.

The obligation is precision. Documentation and code must say "repo view",
"workspace view", or "project view" when discussing UX, and must say
CAS/reference-set when discussing protocol state.

### Alt H: No native repo concept

POC18 would avoid "repo" entirely in the native model. Users manipulate
workspaces, reference sets, snapshots, logical changes, and releases. "Repo" is
only a Git bridge term.

This is the cleanest protocol vocabulary. It avoids carrying Git assumptions
forward.

The downside is user comprehension. Existing Git users need migration language.
DevOps users also need to talk about "the root filesystem repository" or "the
configuration repository" even if the native substrate is a reference-set view.
Avoiding the word completely may make the system harder to explain without
improving correctness.

## Scenario analysis

### Scenario 1: Local edit with native object creation

Alice edits `README.md` and a small PNG asset. The implementation must Rabin
chunk both files, store chunk CIDs in CAS, write chunk manifests, create POSIX
node promises, update a directory reference set, and write a snapshot/change-set
with parent links.

Alt A handles this fastest. Alt B proves less about object correctness per unit
time because container orchestration distracts from chunk/manifests. Alt C is
premature. Alt D gets the benefit of Alt A while requiring Alice's local CAS to
be treated as partial even in tests. Alt E would store this as Git blobs and
trees too early. Alt F works but lacks bridge pressure. Alt G explains the
project as a repo view over a root directory reference set. Alt H explains it as
only a workspace/reference-set and may be less familiar.

The important implementation lesson is that every local object constructor must
accept absent referenced objects as a first-class state. A snapshot can point at
a manifest whose chunks are not local yet; local materialization then promises
only what it can safely produce from available objects.

### Scenario 2: Sparse reference-set walk

Bob receives Alice's branch-role reference set. The `head` entry points at a
snapshot Bob lacks. The snapshot points at a root directory Bob lacks. The
directory points at two file nodes and a chunk manifest; Bob has only one chunk.

Alt A fails if it assumes branch fetch means all objects arrive together. Alt B
tests the sparse case naturally. Alt D must force this case into the first local
tests even before Docker. Bob's local background puller should notice missing
CIDs and ask selected peers for promise-shaped retrieval. Alice or Carol can
promise to serve exact bytes, promise retention for a period, require a token,
or decline.

There is no RPC call named "get object". There is a promise exchange: Bob
promises interest in receiving CID X under constraints; Alice promises to serve
CID X if her local trust/economics conditions are met; Alice fulfills by sending
bytes; Bob verifies bytes by CID and records local keep/break events.

### Scenario 3: CAS shared across repo-like views

Carol has two projects that vendor the same large dependency tarball. Dave has a
root-filesystem view that includes the same tarball as an install package. Ellen
has a release view that labels the tarball as `binary`.

If POC18 stores objects inside repo-private databases, it duplicates storage and
makes sparse sync harder. Alt G handles this cleanly: a node/agent CAS stores
chunks once, and multiple reference-set views label those chunks differently.
Alt H also handles it, but has less friendly migration language for Git users.

The TE outcome should require code paths to take a CAS root plus view/reference
set, not a repo-private object store. A "repo" may own configuration, ignore
rules, workspace roots, and UI defaults, but not object identity.

### Scenario 4: Rename, copy, hard link, and history

Alice renames `README.md` to `docs/intro.md`, copies it to
`docs/tutorial.md`, and has two directory entries that intentionally share one
hard-link identity.

Git follows renames heuristically after the fact. Jujutsu improves change
identity, but filenames are still tree entries. POC18's versioned reference-set
model says the directory owns labels; a file or POSIX node lineage is separate.
The rename is a directory-reference-set history change. The file node lineage
continues. The copy creates either a new node lineage sharing chunks or a
separate label pointing at the same content, depending on promised semantics.
The hard link is represented by multiple labels that intentionally share a node
identity or link-group target.

This supports the user's intuition that filenames and tags are the same family:
both are labels inside versioned reference sets. A multi-target tag with labels
is directory-shaped; a directory is a reference set whose role constrains labels
to path components and targets to node/directory CIDs.

### Scenario 5: Logical change across many revisions

Frank publishes a feature. He revises it six times after reviews by Grace and
Heidi. The logical identity of the feature must survive rebasing, squashing,
interdiffs, and final merge.

Jujutsu uses a change ID to connect revisions of the same logical change. POC18
should not copy that as a raw random field inside commits. The POC18 equivalent
is a `logical_change` reference set with its own CID history. Its entries can
label `current`, `submitted`, `reviewed_by_grace`, `reviewed_by_heidi`,
`accepted_by_frank`, `superseded`, and other exact target CIDs. The reference
set itself is versioned, signed, and parent-linked.

This is at least a functional superset of simple change IDs for POC18's goals:
it provides stable identity, preserves prior positions, supports multiple
targets and labels, and can carry review/adoption context. A raw change ID may
still be useful as an imported Jujutsu bridge field, but it is not the native
authority.

### Scenario 6: Review, merge availability, and pull-request replacement

Alice and Bob diverge. Carol notices both changes touch nearby code. Dave's
agent can compute that a clean merge is available now but may become complex if
Alice edits one more file. Ellen reviews Bob's change and promises local
acceptance of one version but not another.

POC18 should not recreate GitHub pull requests as forge-owned authority. It
should provide in-band reference sets for review threads, issue threads, merge
availability, requested changes, test results, comments, and adoption promises.
The "merge availability" message is a local promise or local observation by the
agent that computed it: "I promise that, according to my current local view and
tooling, these parents can be merged without conflict" or "I promise to alert
Alice that my local analysis predicts a future conflict if this divergence
continues." Other agents decide locally whether to trust that promise.

This is where LLMs matter, but only after deterministic object and message
shapes exist. LLM-based merge assistance should produce local promise intent
and candidate edits; implementation code must still create valid CAS objects,
reference sets, signatures, and parent links.

### Scenario 7: In-band issues, TODOs, chat, and social workflow

Grace opens an issue. Heidi starts an in-band TODO list. Ivan starts a chat
thread because divergence appears. Judy runs an LLM-assisted review. None of
these should require Slack, Discord, GitHub Issues, GitHub Pull Requests, or a
forge database.

POC18 can model each as a reference-set role or thread role:

- `issue_thread` labels reports, reproductions, proposed fixes, status promises,
  and related code CIDs.
- `todo_thread` labels task promises, owners, dependencies, and completion
  promises.
- `chat_thread` labels messages tied to divergence, review, or release context.
- `review_thread` labels submissions, review comments, test promises, and
  acceptance/non-acceptance promises.

This TE leaves exact role names to DF unless the existing POC18 spec
already names them, but it should require that these collaboration artifacts are
in-band CAS/reference-set objects, not external SaaS records.

### Scenario 8: Tangled prior-art pressure

Mallory hosts a Tangled-like self-hosted service with Git/SSH compatibility,
ATProto identity, appview aggregation, and round-based pull requests. Alice's
community wants the social UX without inheriting Git/SSH push-pull, role-based
access control, appview authority, hidden Git refs, or raw Jujutsu change IDs as
native protocol concepts.

POC18 should adopt the lessons, not the mechanism. Self-hosting, migration,
social-code UX, round-based review, stable logical change identity, and
Git/SSH interop pressure matter. Native POC18 remains sparse CAS plus local
promises. An appview-like service can promise summaries or indexes; it does not
become truth. A Git/SSH endpoint can be a bridge; it does not become native
sync. ATProto identity can inform local trust and interop; it does not become a
global trust authority.

### Scenario 9: Conventional Git import/export/push/pull

Alice imports an existing Git repository with branches, tags, symlinks, merges,
and binary files. Later she exports a compatible view back to Git. Bob pulls
from a Git remote into PromiseGrid. Carol pushes a compatible projection back to
the Git remote.

POC18 is incomplete unless all four bridge operations exist. Import and pull
share the Git-to-PromiseGrid conversion path. Export and push share the
PromiseGrid-to-Git conversion path. Both paths should use the same mapping core
for Git refs, commits, trees, blobs, tags, modes, parent links, and notes about
loss or non-commitment.

Unsupported POSIX node types must not be silently lost. Device nodes, FIFOs,
sockets, xattrs, hard-link identity, and materialization constraints may need
explicit bridge records: "Git cannot represent this node; this export promises
only the compatible subset" or "this push declines to materialize this object in
Git." That is better than pretending Git is the schema authority.

### Scenario 10: DevOps live root modification

Alice runs POC18 as root on a live host. The view contains `/etc`, systemd unit
files, package artifacts, container images, service scripts, and generated
state. Alice applies one ordered change, runs a trigger to restart a daemon,
validates the daemon, then applies the next ordered change.

The paper in
`implementations/poc18-cas-git-replacement/docs/turing-equiv.html` makes this
more than a UX requirement. A self-administered host runs tools in the context
of the target it is changing. If a change modifies the tool, its dependencies,
the interpreter, libraries, filesystem, package database, service manager, or
network stack, later changes may behave differently. The least-cost safe path is
deterministic, repeatable change order: the sequence tested in staging should
replay in production in the same order.

POC18 therefore needs ordered change sequences as first-class CAS objects, not
just unordered desired-state convergence. A snapshot tells "what"; an ordered
change log tells "how this state was reached". A DevOps mode must support:

- include/exclude rules for root-sized trees;
- large binary/package/container-image storage in-band;
- incremental apply with checkpoints;
- triggers after selected changes;
- local validation promises;
- local refusal if a change is unsafe or a required capability token is absent;
- exact replay of the same ordered sequence on another host.

### Scenario 11: Root-filesystem scale

Bob versions a root filesystem with thousands of files and directories, package
caches, container layers, logs, generated files, and large binaries. He wants to
ignore some paths, include others, store large blobs in-band, and avoid
re-uploading unchanged binary chunks.

Rabin chunking must apply to all regular files, not just "large files". The CAS
deduplicates identical chunks across root views, project views, releases, and
package archives. Include/exclude rules are local materialization and scan
promises, not universal policy. A view may promise "I scan these paths and
exclude those paths"; another agent may accept, ignore, or distrust that view.

This scenario strongly supports Alt D plus Alt G: the local core must be fast
and deterministic, but storage cannot be repo-private or complete-local.

### Scenario 12: Triggers and local resources

Carol applies a configuration update that requires restarting `nginx`. Dave
applies a kernel update that requires reboot. Ellen applies an LLM-generated
mechanical edit to a config file but does not promise to restart anything.

Triggers are local promises around local resources. A trigger role can promise
to restart a daemon after receiving a matching ordered-change object and a valid
local capability token. It can decline if the token is absent, the host is not
in the required state, the service manager is unavailable, or local policy makes
the action unsafe. This is not authorization vocabulary; it is local resource
promise and local resource withdrawal.

POC18 should record trigger promises, trigger results, and validation outcomes
as parent-linked objects so a later audit can compare test and production
replay order.

### Scenario 13: Porcelain/plumbing library boundary

Frank writes a CLI command to create a branch-like reference set. Grace writes a
GUI to review a merge. Heidi writes an automation script. Ivan writes a remote
UI that controls a workspace view on another node.

Git historically has a split where porcelain commands and plumbing libraries do
not always share one coherent library boundary; external automation often needs
libgit or shelling out. POC18 should avoid repeating this mistake. The same
library that implements object creation, CAS access, reference-set mutation,
materialization, bridge conversion, and sync intent should be used by CLI,
GUI, automation, and tests.

This TE compares terminology:

- "porcelain/plumbing" is familiar to Git users but imports Git culture.
- "UI/backend" is clear but generic.
- "view/control plane" may fit PromiseGrid but can sound distributed-system
  heavy.
- "user surface / core library" is plain and avoids Git baggage.

The likely recommendation is to document both: use "user surface" and "core
library" in PromiseGrid-native docs, and mention "porcelain/plumbing" only as a
Git-user analogy.

### Scenario 14: Local `grid()` messages between UI and core

Judy's local GUI wants to ask a local backend to create a reference set. Ken's
remote UI wants to operate a workspace on a headless node. Lori's automation
wants to replay a DevOps change sequence on a fleet.

If UI/backend calls are ordinary Go method calls only, local development is
simple. If every UI/backend exchange is a `grid()` message, the same interface
can cross process and host boundaries. But forcing every internal call through
CBOR may slow development and produce unnecessary ceremony.

The surviving shape should likely be layered:

- the core library exposes typed methods for local code;
- every typed method has an equivalent pCID-owned promise/message shape;
- local UI can call the library directly;
- process/host boundaries use `grid()` messages;
- tests verify that direct calls and message calls converge on the same object
  bytes and promise semantics.

This avoids Git's automation split while preserving future remote UI control.

### Scenario 15: Promise economics for storage and forwarding

Mallory asks Bob to store a large release artifact. Bob has little disk. Carol
offers storage credits. Dave offers archival retention for a bearer token.
Ellen promises to forward missing chunks for peers she trusts.

POC18 should not rely on altruistic storage. It should test capability-token
promises and local economics:

- storage tokens promise future retention or serving;
- retrieval tokens promise future serving of exact CID bytes;
- forwarding tokens promise relay/forwarding behavior for scoped objects;
- review or compute tokens promise bounded review/verification effort;
- archival tokens promise longer retention windows.

Agents locally price these promises. Exchange rates, willingness to accept
tokens, and trust in token issuers are local. There is no central exchange or
global storage market. The token is evidence of the issuer's promise, not a
permission slip from an authority.

### Scenario 16: Retention and garbage collection

Alice pins a release. Bob promises to retain a chunk manifest for 30 days.
Carol stores a working-set object without a retention promise. Dave is under
capacity pressure and must collect unpromised objects.

GC decisions must be promise-based. An agent should retain objects it promised
to retain, objects locally pinned, objects needed by current views, objects paid
for by tokens, objects inside a review window, or objects it locally values.
Objects outside those scopes can be dropped under pressure. Dropping an object
is not a broken promise unless a prior outstanding retention promise covered
that object.

POC18 should model GC as local state over reference sets and tokens. Analyzer
gates should require at least one promised-retained object and one unpromised
dropped object.

### Scenario 17: Corruption and malicious peers

Mallory sends bytes that do not match the advertised CID, a malformed reference
set, a fake parent link, or a Git bridge mapping that silently drops device-node
metadata.

The recipient verifies exact bytes by CID. Mismatched bytes are local malformed
events. Missing parents are sparse state, not proof of bad faith. A fake claim
to have stored or reviewed something can lower local trust if it breaks a prior
promise. A bridge mapping that loses data must be explicit. Nothing requires a
global ban list or conformance authority.

This scenario tests that POC18 keeps Promise Theory semantics even when the
workflow resembles source-control security.

### Scenario 18: Long-horizon migration and ecosystem interop

In ten years, Alice still uses Git tooling for some projects. Bob uses
PromiseGrid-native sync. Carol uses a Tangled-like social code host. Dave uses
ATProto identity. Ellen archives release reference sets. Frank runs root
filesystem DevOps views. Grace runs old hardware that only understands compact
CBOR profiles.

POC18 should not require a flag day. Git bridge behavior supports existing
tools. Reference sets and CIDs support native sync. ATProto/Bluesky-style
identity can inform local trust and interop without becoming authority. Compact
pCID-selected payloads and binary CIDs keep constrained devices plausible.
Because all object identities are CIDs and specs are pCID-selected, old objects
remain interpretable if their specs are retained.

This scenario rejects both Git-first and native-only extremes.

## Comparison summary

| Criterion | A Local-first | B Docker-first | C LLM-first | D Hybrid staged | E Git-first | F Native-only | G Repo views | H No repo |
|---|---|---|---|---|---|---|---|---|
| First implementation speed | Strong | Medium | Weak | Strong | Medium | Medium | N/A | N/A |
| Sparse-CAS pressure | Weak unless forced | Strong | Medium | Strong | Weak | Strong | Strong | Strong |
| Object-model correctness | Strong | Medium | Weak | Strong | Medium | Strong | Strong | Strong |
| Avoids Git vocabulary gravity | Strong | Strong | Medium | Strong | Weak | Strong | Medium | Strong |
| Git migration | Medium | Medium | Weak | Strong | Strong | Weak | Strong | Medium |
| PromiseGrid semantics | Medium if disciplined | Strong | Risky | Strong | Weak | Strong | Strong | Strong |
| DevOps ordered replay | Strong if modeled early | Medium | Weak | Strong | Weak | Strong | Strong | Medium |
| LLM workflow realism | Deferred | Deferred | Strong | Deferred but safe | Medium | Medium | N/A | N/A |
| Risk of blind path | High if local-complete | Medium | High | Low | High | Medium | Low | Medium |

## Rejected alternatives

- Reject Alt C as the first implementation path. LLM autonomy is important, but
  POC18 should not let LLM text outrun protocol-valid object creation.
- Reject Alt E as the native architecture. Git compatibility is required as a
  bridge, not as the source of truth.
- Reject Alt F as the complete POC18 plan. Native-only is clean but fails the
  explicit Git import/export/push/pull requirement.
- Reject a plain Alt A local-first path unless it is modified to include sparse
  CAS, missing-object, background-puller, and promise-shaped retrieval tests
  from the first slice.
- Reject a plain Alt H vocabulary path if it makes migration harder. "Repo" can
  survive as a user-facing view word as long as it is not a protocol authority.

## Surviving architecture

The surviving architecture is **Alt D plus Alt G**:

> Build a deterministic local core first, but design it from the first line as
> sparse multi-agent CAS/reference-set infrastructure. Treat repos as
> user-facing views over shared CAS and versioned reference sets, not as
> authoritative storage databases. Add Docker multi-agent sync immediately after
> the local core can create, walk, materialize, and diagnose a small graph. Add
> LLM autonomy after deterministic and Docker paths prove the protocol surface.

The first slice should therefore include:

- one shared node/agent CAS abstraction;
- repo/workspace/project views over reference sets;
- Rabin chunking and manifests for all regular files;
- POSIX node objects for every inode type;
- directory reference sets with filename labels;
- branch/tag/logical-change/review/workspace reference-set roles;
- snapshot/change-set objects with parent links;
- missing-object states and background retrieval abstraction;
- promise-shaped retrieval/retention/forwarding interfaces, even if fulfilled
  locally in the first run;
- Git bridge conversion seams, even if full push/pull comes later;
- DevOps ordered-change objects, triggers, and replay hooks as first-class
  planned surfaces, not afterthoughts.

## Implications for `nahop` subtasks

- `nahop.2`: must lock paths, command names, package names, runtime paths, and
  generated CAS path patterns for a deterministic local core that already
  models sparse CAS and repo views.
- `nahop.4`: per-agent sparse CAS stores must be filesystem-based from the
  start, shared by multiple repo/workspace views, and keyed by CIDv1 base32
  printable paths.
- `nahop.5`: Rabin chunking applies to every regular file; large-file support is
  native, not a special mode.
- `nahop.6`: POSIX node envelopes must cover all inode types and local
  materialization constraints.
- `nahop.7`: directory reference sets own filename labels; file/node history
  does not own a canonical path.
- `nahop.8`: branch, release, tag, logical-change, issue, review-thread,
  todo-thread, chat-thread, and workspace roles should be evaluated as
  reference-set roles or role extensions before separate pCIDs are introduced.
- `nahop.9`: snapshot/change-set objects must support both state and ordered
  replay context where DevOps scenarios require it.
- `nahop.10`: materialization must be local-promise-shaped and safe for root and
  non-root workspaces.
- `nahop.11`: peer retrieval must be promise-shaped, with capability-token and
  local non-commitment paths.
- `nahop.12`: rename/copy tests must verify directory-label history and node
  lineage.
- `nahop.13`: merge scenarios must include multi-parent snapshots and local
  merge-availability promises.
- `nahop.14`: reviews must be in-band reference-set promises, not forge
  approvals.
- `nahop.15`: import/export/push/pull must share bridge conversion paths.
- `nahop.17`: retention and GC must be promise-based and token-aware.
- `nahop.19`: continuous sync must exchange reference-set heads,
  availability/retention promises, and missing-object interests without native
  push/pull commands.
- `nahop.20`: analyzer gates must include sparse CAS, shared CAS across views,
  Rabin all-file chunking, POSIX inode coverage, ordered DevOps replay,
  promise-shaped retrieval, anti-RPC vocabulary, and Git bridge roundtrip.
- `nahop.21`: diagnostics must render raw CBOR messages for reference sets,
  snapshots, missing-object interests, retention promises, trigger/replay
  promises, and Git bridge mappings.
- `nahop.22`: clean deterministic runs must archive commands, CAS examples,
  reference walks, missing-object walks, Git bridge output, DevOps replay output,
  and analyzer output.

## New follow-up questions for DF

Before code scaffolding, POC18 still needs DF answers for:

1. Exact first-slice runtime shape: single binary with subcommands, multiple
   local processes, or one deterministic simulator command that owns all roles.
2. Exact path bundle: command paths, package names, runtime run root, CAS path
   layout, workspace view path layout, and generated diagnostics path layout.
3. Exact terminology: "repo view", "workspace view", "project view", or another
   user-facing term for views over shared CAS/reference sets.
4. Exact local API boundary: typed library calls only, `grid()` messages between
   local UI/backend processes, or both with equivalence tests.
5. Exact Git and Rabin libraries, with verification that the chunker implements
   the required content-defined Rabin behavior.
6. Exact minimum DevOps slice: whether ordered-change and trigger objects appear
   in the first deterministic run or are locked as second-slice work before
   Docker.

## Decision status

This TE is **needs DF**.

The recommendation is Alt D plus Alt G: deterministic local core first, forced
to preserve sparse multi-agent CAS assumptions from day one, with repo-like
views over shared CAS and versioned reference sets. Docker multi-agent sync
should follow immediately after basic storage/plumbing/UI work. LLM autonomy
should follow only after deterministic and Docker runs prove valid protocol
behavior.

No POC18 code scaffolding should start until the DF questions above are answered
and the locked implementation decision is recorded in `TODO-nahop`.

## Implications for open TODOs and pending DIs

- `TODO-nahop` should record this TE as the hard-gate architecture review for
  the post-`TE-kopap` implementation sequencing question.
- `TE-kopap` remains valid for rejecting Git-first and native-only as overall
  architecture baselines. `TE-vahoj` adds the first-slice sequencing, sparse CAS,
  DevOps, UI/backend, and session-log-superset analysis.
- `DN-rifir` remains the reference-set design note. `TE-vahoj` does not replace
  it; it tests how to implement its conclusions without local-first or
  Git-shaped drift.
- `DN-dopod` remains the Tangled prior-art review. `TE-vahoj` routes its lessons
  into implementation sequencing and scenario gates.
- The implementation-local POC18 version-control spec may need refinement after
  DF, especially around ordered DevOps replay, thread/reference-set roles for
  issues/TODOs/chat, and promise-shaped retrieval/economics. Such refinements
  must not add repo-local helper-command documentation to the pCID-derived spec.
