# TODO-nahop: POC18 CAS-backed Git replacement

## Status

Planned. Owns a future `implementations/poc18-cas-git-replacement/` proof of
concept that treats PromiseGrid CAS as object storage and uses grid-envelope
parent links plus versioned reference-set promises for file, directory, branch,
release, logical-change, review, and workspace version chaining.

POC18 is a **protocol/CAS/kernel superset of POC16**, not a POC17 runtime
superset. It inherits POC16's CID, pCID, CAS, parser/builder, structured payload,
CWT/COSE token, encrypted-payload, lifecycle-token, raw-message, and analyzer
lessons. It does not inherit POC17's M4/LoRa runtime scope. Source: `DI-vilum`;
`DI-zuruj`.

`TE-vahoj` is the hard-gate POC18 superset architecture TE before code
scaffolding. It extends `TE-kopap` with today's session-log requirements:
sparse shared CAS, background chunk pulling, promise-shaped retrieval and
capability-token economics, UI/backend library boundaries, in-band
collaboration, and DevOps ordered replay. Source: `DI-fusir`.

## Decision Intent Log

ID: DI-vilum
Date: 2026-06-25 14:56:03 PDT
Status: superseded
Author: stevegt@t7a.org (Steve Traugott)
Decision: Plan POC18 as a Git-replacement proof of concept where CAS stores exact objects and grid-envelope parent links provide file, directory, branch, merge, review, and release version chaining.
Intent: Recent POCs have already proven sparse per-agent CAS stores, exact CBOR envelopes, parent links, pCID-selected parsing, local trust, token incentives, and Promise Theory-compatible kernel roles. POC18 should apply those pieces to the collaboration problem Git and GitHub currently solve poorly under LLM-scale workflows: storing content, naming versions, sharing branches, reviewing changes, merging work, preserving provenance, and deciding what each agent locally trusts. The design should preserve Git's strongest ideas, especially content-addressed objects and parent-linked history, while replacing global forge authority and repository-centered assumptions with local promises, sparse CAS replication, and signed branch/ref promises.
Constraints: Preserve one top-level semantic action `promise`; preserve `grid([42(pCID), ...protocol-defined-slots])`; keep pCID as protocol-spec selector, not branch name, file path, agent address, operation code, or repository name; use binary CIDs on wire and CIDv1 base32 text when printable; do not create a global branch authority, global merge authority, global CAS, global monitor, permission service, authorization service, or conformance service; treat branch heads, tags, releases, reviews, and merge decisions as voluntary signed promises by agents or groups; keep CAS stores sparse and peer-relative; use grid-envelope parent links for version graph edges unless a later TE/DI proves a payload-parent exception; decide explicitly whether POC18 is a strict runtime superset of POC17 or a protocol-focused branch with an approved scoped exception for M4/LoRa fidelity.
Affects: protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md; protocols/wire-lab.d/TODO/TODO.md; future implementations/poc18-cas-git-replacement/; DEV-GUIDE-RESOURCES.md; future docs/protocols or implementation-local protocol specs for CAS version graph, branch/ref promises, review promises, and workspace materialization.
Superseded by: DI-zuruj

ID: DI-zuruj
Date: 2026-06-26 11:38:43 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Refine POC18 around versioned reference-set promises as the shared abstraction for tags, directories, filenames, branches, releases, logical changes, reviews, workspaces, and fetch entrypoints; set POC16 as the POC18 superset baseline.
Intent: The POC18 discussion converged on the fact that a multi-target tag with labels is directory-shaped: both are scoped reference sets from names to target CIDs. Git's distinction between branches, tags, trees, and optional fetched tags is not the right root abstraction for PromiseGrid. POC18 should instead test versioned reference sets whose roles define validation and behavior. A branch is a reference set with role `branch`; a directory is a reference set with role `directory`; a filename is a directory-scoped label; a logical change is a reference set with role `logical_change`, replacing Jujutsu-style intrinsic change IDs with a PromiseGrid-native, content-addressed, promise-bearing object. This keeps Git-compatible snapshots and parent DAGs while making renames, aliases, review bundles, releases, and fetch/discovery first-class promises.
Constraints: Supersedes `DI-vilum` where it names branch/ref promises as the root abstraction or names POC17 as the superset decision point; do not make tags optional side data like Git's `--tags`; use reference-set promises as the primary fetch/discovery objects; filenames belong to directory/reference-set history, not file-history identity; use one main version-control pCID unless a later TE/DI proves a distinct protocol boundary; use restic's content-defined chunker library but store chunks and Merkle manifests in PromiseGrid CAS formats; start PromiseGrid-native but require content-and-DAG Git import/export before POC18 is complete; use a Go Git library for import/export; preserve binary CIDs on wire and CIDv1 base32 text when printable; do not introduce global namespace, branch, tag, merge, review, or forge authority.
Affects: protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md; docs/research/DN-rifir-poc18-versioned-reference-sets.md; protocols/wire-lab.d/TODO/TODO.md; future implementations/poc18-cas-git-replacement/; future implementation-local POC18 protocol specs.
Supersedes: DI-vilum

ID: DI-dibut
Date: 2026-06-26 11:53:59 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add Tangled as a POC18 prior-art question and replace explicit Git-style push/pull operations with continuous peer DAG sync.
Intent: POC18 should learn from adjacent decentralized code-collaboration systems without copying Git's push/pull mental model. Tangled is relevant prior art because it combines decentralized Git hosting, self-hosted infrastructure, social coding, AT Protocol identity, round-based pull-request ideas, and Jujutsu change IDs. PromiseGrid should compare those ideas against versioned reference sets, local trust, and sparse CAS. Separately, POC18 should treat collaboration as continuous DAG synchronization among trusted peers: agents exchange reference-set promises, object availability promises, parent-linked CAS objects, and missing-object requests continuously rather than waiting for explicit `git push` or `git pull` commands.
Constraints: Do not adopt Tangled concepts without a focused review; do not reintroduce a global forge, appview, role-based access authority, or Git remote as the PromiseGrid authority model; preserve POC18's versioned reference-set root abstraction; continuous sync must remain voluntary and peer-relative, with each agent deciding what to advertise, request, retain, forward, verify, and trust; Git import/export remains compatibility work, not the native synchronization model.
Affects: protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md; docs/research/DN-rifir-poc18-versioned-reference-sets.md; future implementations/poc18-cas-git-replacement/; future POC18 analyzer gates.

ID: DI-dofoj
Date: 2026-06-27 14:50:16 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Refine POC18 to use Rabin chunking for all file content and require a conventional Git bridge covering import, export, push, and pull.
Intent: POC18 should differ from Git by natively handling large files in-band: every file becomes Rabin content-defined chunks plus PromiseGrid Merkle/index manifests in CAS, instead of treating large files as an out-of-band special case. POC18 should also interoperate with conventional Git in both directions. Import/export and push/pull should share the same conversion core between Git refs/objects and PromiseGrid reference sets, snapshots, manifests, chunks, and mapping records. Git push/pull remains a compatibility bridge to Git remotes; native PromiseGrid collaboration remains continuous peer DAG synchronization over local promises, sparse CAS, reference-set advertisements, object-availability promises, and missing-object requests.
Constraints: Use Rabin chunking as the design requirement for all files; verify any chosen Go library provides the required Rabin-style content-defined chunking before implementation lock; do not use Git LFS-style out-of-band large-file handling as the native POC18 model; do not make a Git remote, forge, appview, or push/pull endpoint the native PromiseGrid authority; keep Git push/pull bridge code separate from native peer DAG sync while sharing the Git-to-PromiseGrid and PromiseGrid-to-Git conversion core with import/export.
Affects: protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md; docs/research/DN-rifir-poc18-versioned-reference-sets.md; docs/research/DN-dopod-poc18-tangled-prior-art.md; DEV-GUIDE-RESOURCES.md; protocols/wire-lab.d/TODO/TODO.md; future implementations/poc18-cas-git-replacement/.
Supersedes: DI-zuruj and DI-dibut only where they limit required Git compatibility wording to import/export or imply conventional Git push/pull is not a required bridge surface.

ID: DI-radaj
Date: 2026-06-27 14:59:56 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Refine POC18 filesystem modeling to support every POSIX inode type: regular file, directory, symbolic link, hard link, character device, block device, FIFO, and socket.
Intent: POC18 should be a serious filesystem/versioning substrate, not only a text-file source-control demo. A PromiseGrid workspace needs to preserve enough filesystem meaning to version, sync, review, and materialize real POSIX trees. Regular files keep Rabin-chunked content manifests. Directories remain versioned reference sets. Symbolic links are node promises whose payload preserves link target bytes. Hard links are represented by multiple directory labels pointing at the same node identity or link-group target rather than by duplicated file contents. Character devices, block devices, FIFOs, and sockets are metadata-bearing node promises; POC18 records the promised node type and materialization constraints, not live kernel state or stream contents. Local materialization may refuse or adapt nodes the host cannot safely create.
Constraints: Do not collapse every POSIX object into a regular file blob; do not claim conventional Git can preserve inode types it cannot represent; Git bridge import/export/push/pull must preserve supported Git modes and must record explicit mapping, loss, refusal, or local non-commitment for POSIX inode types outside Git's normal tree model; device and socket materialization remains local capability/resource behavior, not global authority.
Affects: protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md; docs/research/DN-rifir-poc18-versioned-reference-sets.md; DEV-GUIDE-RESOURCES.md; protocols/wire-lab.d/TODO/TODO.md; future implementations/poc18-cas-git-replacement/.
Supersedes: DI-zuruj only where it implies POC18's filesystem model is limited to regular file and directory versions.

ID: DI-lidaj
Date: 2026-06-27 21:54:56 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement `nahop.3` as one normative implementation-local POC18 version-control protocol spec at `implementations/poc18-cas-git-replacement/docs/protocols/version-control.md`; do not use `-v1` or `_v1` because the pCID of the exact spec bytes is the version; create a same-directory CIDv1 base32 alias symlink after the spec bytes are finalized; use svgbob-safe ASCII diagrams; and document `grid([42(pCID), parents, payload, proof])` plus annotated examples as the POC18 version-control message shape.
Intent: POC18 needs an RFC-like pCID spec before code can safely use a version-control pCID. The spec must be comprehensive enough for a developer to implement from the document, must encompass the related design notes and thought experiments, and must avoid stale version-number vocabulary by treating the pCID as the version. Keeping the spec under the implementation-local `docs/protocols/` tree follows the POC16 pattern and prevents root docs from becoming a stale competing protocol authority.
Constraints: The spec is the single normative POC18 version-control pCID unless a later TE/DI proves a distinct protocol boundary; pCID remains a protocol selector, not a branch name, operation code, address, repository name, or message type; diagrams must be plain ASCII suitable for svgbob-style conversion; examples must annotate message slots and payload fields; Git import/export/push/pull remain bridge adapter behavior rather than native synchronization; native synchronization remains continuous peer DAG sync; the spec itself must not embed its own final pCID because the pCID is derived externally from exact bytes.
Affects: protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md; implementations/poc18-cas-git-replacement/docs/protocols/version-control.md; implementations/poc18-cas-git-replacement/docs/protocols/<base32-pCID>.md; DEV-GUIDE-RESOURCES.md.

ID: DI-fusir
Date: 2026-06-28 19:38:10 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add `TE-vahoj` as the hard-gate POC18 superset architecture TE before `nahop.2` code-scaffolding decisions.
Intent: The POC18 plan grew beyond the earlier Git/native architecture TE. Today's session added implementation-sequencing risk, sparse shared CAS, active background chunk pulling, promise-based inter-agent retrieval, capability-token economics, UI/backend library boundaries, possible local `grid()` messages, in-band GitHub replacement surfaces, and DevOps ordered replay based on `implementations/poc18-cas-git-replacement/docs/turing-equiv.html`. These requirements must be tested together before code starts so POC18 does not drift into a local-only, Git-shaped, or RPC-shaped blind path.
Constraints: Keep `TE-kopap`, `DN-rifir`, `DN-dopod`, and the current POC18 spec as inputs; do not use the new TE to rewrite history or relax existing DIs; do not scaffold POC18 code until the TE's remaining DF questions are answered and locked; preserve pCID discipline, CID rendering rules, Promise Theory vocabulary, and POC16-superset scope.
Affects: docs/thought-experiments/TE-vahoj-poc18-superset-architecture.md; docs/thought-experiments/README.md; protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md; protocols/wire-lab.d/TODO/TODO.md; protocols/wire-lab.d/specs/harness-spec-draft.md; DEV-GUIDE-RESOURCES.md; implementations/poc18-cas-git-replacement/docs/turing-equiv.html; future implementations/poc18-cas-git-replacement/.

ID: DI-jifuj
Date: 2026-06-28 20:42:49 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Lock POC18's first implementation slice as Vahoj Alt-D plus Alt-G: build a deterministic local core first, but force sparse multi-agent CAS assumptions through interfaces, tests, and runtime paths from day one.
Intent: POC18 should not start as a local-only Git clone or as an overlarge multi-agent runtime. The useful first step is a core graph that can create, store, walk, diagnose, and materialize native PromiseGrid version-control objects while preserving sparse-CAS, missing-object, and peer-promise seams for the next slice. This lets implementation begin without forgetting the long-term decentralized sync model.
Constraints: Use `grid` for the user-facing CLI; use `poc-*` for non-production deterministic fixtures; use `grid-*` for production-shaped daemons and backend processes when they appear; use `/tmp/wire-lab-poc18-*` for generated runtime state in tests and clean runs; never run DevOps/root-filesystem tests on the host; keep DevOps replay tests container-only; keep pCID as protocol selector, not address, operation, path, repository, or branch name; keep native sync promise-shaped and peer-relative; do not implement real peer transport in the first slice.
Affects: implementations/poc18-cas-git-replacement/; protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md; docs/thought-experiments/TE-vahoj-poc18-superset-architecture.md; future POC18 run scripts and tests.

ID: DI-harih
Date: 2026-06-28 20:42:49 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Lock the POC18 CLI/core split after `TE-hikar`: documentation and code comments use `CLI/core`; package architecture uses layer packages `store`, `graph`, `workspace`, `sync`, and `bridge`; the first `grid` CLI balances familiar VCS words with PromiseGrid-native objects through `init`, `ingest`, `snapshot`, `checkout`, `refs`, and `diag`.
Intent: Vahoj settled the object model but not the command language. The follow-up CLI TE chooses commands that do not copy Git's staging or native push/pull assumptions, while still being learnable for developers familiar with Git, Jujutsu, Mercurial, SVN, and CVS. Layer packages better match implementation boundaries than one package per object type: `store` owns CIDs/CAS/chunks, `graph` owns protocol objects and envelopes, `workspace` owns scan/materialize, `sync` owns sparse retrieval promises, and `bridge` owns Git compatibility seams.
Constraints: Mention porcelain/plumbing only as a Git analogy; do not introduce a raw Jujutsu-style change-ID field; keep Git import/export/push/pull under bridge behavior, not native sync; keep all interagent communication promise-shaped; use exact CIDs as identities, binary on wire and CIDv1 base32 when printable; first-slice `grid` commands may be intentionally narrow but must route through the same core library as fixtures.
Affects: docs/thought-experiments/TE-hikar-poc18-grid-cli-command-model.md; implementations/poc18-cas-git-replacement/; protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md; future POC18 CLI docs.

## Core Hypothesis

POC18 should test this hypothesis:

> A PromiseGrid Git replacement can model source collaboration as versioned
> reference-set promises over CID-addressed CAS objects, with grid-envelope
> parent links carrying version ancestry and local agents deciding which promises
> to trust.

The useful Git idea is not "central repo plus remote branches." The useful idea
is content-addressed objects plus parent-linked history. PromiseGrid should keep
that and add local trust, reciprocal storage/compute promises, sparse
replication, content-addressed reference sets, and pCID-defined CBOR envelopes.
Source: `DI-zuruj`.

## Versioned Reference Sets

POC18's root abstraction is the **versioned reference set**:

```text
reference_set_version:
  role
  namespace_or_context_cid
  entries:
    label -> target CID or target CID set
  parent_reference_set_versions
  promiser
  promise_terms
  proof
```

This is the shared mechanism behind names, tags, branches, directories, review
threads, releases, workspaces, and logical changes. The mechanism is shared; the
role-specific validation is not.

- A **directory** is a reference set with role `directory`. Its labels are
  filename/path-component dirents and its targets are POSIX node-version or
  directory-version CIDs.
- A **branch** is a reference set with role `branch`. It usually has a label such
  as `head` pointing at a snapshot/change-set CID.
- A **release** is a reference set with role `release`. It may label source,
  binary, SBOM, docs, signatures, and review/test artifacts.
- A **logical change** is a reference set with role `logical_change`. It replaces
  a Jujutsu-style intrinsic change ID by making the logical change itself a
  promise-bearing CAS object.
- A **review thread** is a reference set with role `review_thread`. It can point
  at comments, test results, requested-change promises, and accepted versions.
- A **workspace** is a reference set with role `workspace`. It can point at a root
  directory, toolchain, local overrides, dependencies, or materialization
  receipts.

Tags are not optional side objects. A tag/reference-set promise is what Alice
syncs from peers when she wants a branch, release, logical change, directory,
review, or workspace. She then continuously reconciles target CIDs, parent
reference-set versions, and parent-linked CAS objects with peers she locally
trusts. Source: `DI-zuruj`; `DI-dibut`.

## Git And Jujutsu Concepts To Preserve And Reframe

- **Blob/content:** Store every file in CAS using Rabin content-defined chunks.
  The normal file content root is a PromiseGrid Merkle manifest CID over chunk
  CIDs, not a restic repository object, Git blob special case, or out-of-band
  large-file artifact. Source: `DI-dofoj`.
- **Tree/directory:** Reframe as a reference set with role `directory`.
  Filenames are labels/dirents inside the directory's history, like Unix
  directory entries pointing at inodes. A file's logical history does not own one
  canonical path.
- **POSIX inode types:** Reframe filesystem entries as POSIX node promises.
  POC18 must preserve regular files, directories, symbolic links, hard links,
  character devices, block devices, FIFOs, and sockets. Device, FIFO, and socket
  objects preserve node metadata and local materialization constraints, not live
  kernel state or stream contents. Source: `DI-radaj`.
- **Commit/snapshot:** Reframe as a snapshot/change-set envelope that points at a
  root directory reference-set CID and carries parent snapshot links,
  author/promiser intent, local constraints, and review or merge context.
- **Branch:** Reframe as a reference-set role, not a special Git-like pointer.
  Branch fetch means fetching a branch-role reference-set promise, then fetching
  the target snapshot/change-set and required CAS objects.
- **Tag/release:** Reframe as a reference-set role. Tags are versionable,
  parent-linked, signed, and may have multiple labeled targets.
- **Jujutsu change ID:** Reframe as a `logical_change` reference set. The logical
  change is not merely a random field inside each version; it is a
  content-addressed promise object that versions can point back to.
- **Pull request/review:** Reframe as review-thread reference sets and signed
  review/test/adoption promises. There is no forge authority; each participant
  locally decides what review promises affect trust and adoption.
- **Clone/fetch/push:** Reframe native PromiseGrid collaboration as continuous
  peer DAG synchronization. Conventional Git clone/fetch/push/pull still must
  exist as a compatibility bridge to ordinary Git repositories and should share
  conversion code with import/export; natively, peers continuously promise which
  reference sets and CAS objects they have, what they are willing to retain or
  forward, which parent-linked objects they are missing, and under what local
  constraints. Source: `DI-dibut`; `DI-dofoj`.

## Architecture Targets

- Use per-agent sparse CAS stores. No agent is assumed to have all objects,
  reference sets, reviews, releases, branches, or history.
- Use grid envelopes as first-class CAS objects. A version, reference set, review,
  merge, release, or materialization promise is exact bytes with its own CID.
- Use envelope parent links as the version graph. File, directory, branch,
  logical-change, review, release, merge, and reference-set history should be
  traceable by walking parent links through CAS objects.
- Keep names in local or group namespaces. A filename, branch name, release name,
  or review-thread name is not authority; it is a label inside a signed
  reference-set promise.
- Keep all trust local. Alice may trust Bob's branch-role reference set based on
  Bob's make/break history, Carol's review promises, Dave's build promises, or
  her own local checks; no global merge or forge authority exists.
- Keep materialization separate from storage. Checking out a workspace means a
  local role promises to materialize a root directory reference-set CID into
  files, directories, links, devices, FIFOs, sockets, or other local resources
  under local constraints. Source: `DI-radaj`.
- Preserve pCID discipline. pCID selects the protocol parser/builder and slot
  semantics; names, roles, labels, paths, authors, repositories, and destinations
  live in pCID-defined payloads where needed.
- Replace Git-style push/pull as the native collaboration model with continuous
  peer DAG sync. Each agent's local sync role periodically or opportunistically
  compares reference-set heads, parent links, object availability, retention
  promises, and missing-object requests with selected peers, then locally decides
  what to advertise, request, store, verify, forward, or ignore. Conventional Git
  push/pull remains required bridge behavior, not native authority. Source:
  `DI-dibut`; `DI-dofoj`.
- Preserve POC superset discipline by inheriting POC16 protocol/CAS/kernel
  lessons. POC18 does not inherit POC17's M4/LoRa runtime scope. Source:
  `DI-zuruj`.

## Protocol Targets

- Use one main POC18 version-control pCID unless a later TE/DI proves a genuinely
  separate protocol boundary. The pCID should define versioned reference sets,
  file versions, snapshot/change-set envelopes, reviews, materialization, and
  peer fetch semantics.
- Add an RFC-like spec for POC18 version-control messages before implementation
  uses its pCID.
- Define exact CBOR message shapes. Prefer compact arrays for hot-path objects
  and maps for self-documenting review/release records when the pCID spec
  deliberately chooses them.
- Define reference-set roles and validation rules explicitly. Directory labels
  must be valid path components; branch-role sets should have a `head`; release
  sets may be append-only or locally immutable by promise; logical-change sets
  should point at current, prior, review, test, and superseded versions where
  useful.
- Define POSIX node types explicitly. Regular files point at Rabin chunk
  manifests; directories are reference sets; symbolic links preserve target
  bytes; hard links preserve shared node identity or link-group target;
  character/block devices preserve type plus major/minor metadata; FIFOs and
  sockets preserve node type and materialization constraints. Source: `DI-radaj`.
- Define parent-link slot rules explicitly. POC18 should prefer envelope-level
  parent links for version graph edges, with typed parent roles such as
  previous-node-version, previous-reference-set-version, previous-snapshot,
  merge-parent, review-parent, and supersedes.
- Define signature/proof placement explicitly. Reference-set, release, review,
  and adoption promises should be signed by the promiser. COSE or other proof
  shapes must be selected by the pCID spec rather than assumed globally.
- Define Git bridge requirements. POC18 starts PromiseGrid-native but is not
  complete until it can import from, export to, pull from, and push to a real Git
  repo while preserving content, directory structure, branch/tag targets, and
  parent DAG semantics. Import/export local filesystem edges and push/pull remote
  transport edges should share one conversion core between Git refs/objects and
  PromiseGrid reference sets, snapshots, manifests, chunks, and mapping records.
  Git bridge behavior must also record explicit mapping, loss, refusal, or local
  non-commitment for POSIX node types conventional Git cannot represent. Source:
  `DI-dofoj`; `DI-radaj`.
- Define continuous-sync payloads before implementation. The main version-control
  pCID should cover reference-set advertisements, object-availability promises,
  missing-object requests, retention/forwarding offers, and local non-commitments
  without treating any peer as a remote authority. Source: `DI-dibut`.

## CAS And Chunking Targets

- Store exact bytes under CID-derived paths using CIDv1 base32 text when
  printable. No bare SHA-256 hex strings should appear as object identities.
- Use Rabin content-defined chunking for all file content, including small text
  files and large binary files. Large files are native in-band CAS objects, not
  a special out-of-band transport or Git LFS-style exception. Source: `DI-dofoj`.
- Store chunks and Merkle/index manifests in PromiseGrid CAS formats, not restic
  repository format.
- Support at least these object classes:
  - raw chunks;
  - PromiseGrid Merkle manifests over chunk CIDs;
  - POSIX node-version envelopes for regular files, symbolic links, hard links,
    character devices, block devices, FIFOs, and sockets;
  - directory/reference-set envelopes;
  - snapshot/change-set envelopes;
  - logical-change, branch, release, workspace, and review reference sets;
  - review/test/materialization promise envelopes;
  - optional Git bridge mapping records for import, export, push, and pull.
- Keep object type discoverable from CID codec, pCID, or object bytes. Do not
  rely on filename extensions as authority.
- Support sparse retrieval. A fetched reference set may point at target CIDs
  Alice does not yet have; Alice can ask trusted peers to store or send missing
  objects.
- Support promise-based GC. Agents promise retention for selected reference sets,
  pinned heads, releases, paid storage, recent working sets, review windows, or
  peer agreements; they may drop unpromised objects under pressure.
- Support incentives. Agents may exchange bearer or non-transferable capability
  tokens for storage, forwarding, compute verification, review, or long-term
  archival promises.

## Scenarios

- **Single-author edit:** Alice edits `README.md`, Rabin-chunks its content,
  stores chunks and a Merkle manifest in CAS, emits a regular-file node-version
  envelope, then emits a snapshot/change-set envelope pointing at a root
  directory reference set.
- **Rename with history:** Alice renames `README.md` to `docs/intro.md`. The file
  version keeps its logical file lineage; the directory reference-set history
  changes the filename/dirent labels.
- **POSIX node tree:** Alice versions a workspace containing a regular file,
  directory, symbolic link, hard link, character device, block device, FIFO, and
  socket node. The CAS preserves each node promise, and local materialization
  records which nodes were created, adapted, or refused under host constraints.
- **Branch fetch:** Bob asks Alice for a branch-role reference set such as
  `main`. He retrieves its `head` target, then recursively asks for the snapshot,
  root directory, file versions, chunk manifests, and chunks he decides to trust.
- **Logical change revision:** Carol publishes a `logical_change` reference set
  whose `current` label points at a snapshot/change-set. Later revisions update
  the logical-change reference set while preserving older reference-set versions.
- **Divergent logical change:** Alice and Bob each publish current targets for
  the same logical-change role. The result is explicit divergence, not hidden
  branch confusion.
- **Directory-shaped tag:** A release reference set labels `source`, `binary`,
  `docs`, `sbom`, `tests`, and `signature`; those labels point at multiple CIDs
  under one versioned promise object.
- **Merge with conflict:** Dave resolves Alice/Bob divergence by creating a merge
  snapshot/change-set envelope with multiple parent links and a payload promising
  how the conflict was resolved.
- **Review flow:** Ellen reviews Dave's merge CID, publishes a review-thread
  reference set, attaches comments or test-result CIDs, and either promises local
  acceptance or promises requested changes.
- **Sparse missing object:** Frank receives a branch reference set but lacks one
  chunk. He asks Grace for storage/retrieval; Grace may promise to send it for a
  token or decline due to local capacity.
- **Malformed or malicious object:** Mallory sends bytes claiming to match a CID
  or publishes a misleading reference set. The recipient rejects mismatched bytes
  or distrusts Mallory locally without needing a global authority.
- **Git bridge roundtrip:** Alice imports or pulls a real Git repo, maps Git
  branches/tags into reference-set promises, exports or pushes back to Git, and
  verifies content, directory structure, branch/tag targets, and DAG semantics.
- **Continuous peer sync:** Alice and Bob do not run `push` or `pull`. Their
  local sync roles repeatedly exchange reference-set heads, object availability,
  missing-object requests, and retention promises. Each agent decides locally
  which parts of the DAG to fetch, verify, retain, forward, or ignore.
- **Forge replacement:** Alice and Bob collaborate without GitHub. Their local
  agents exchange reference-set promises, review promises, CAS objects, and trust
  updates directly or through chosen peers.

## Analyzer Targets

- Count CAS objects by class, CID, owner, retaining agent, and retrieval source.
- Verify every object identity is a CID, not a bare hash or pseudo-CID.
- Verify reference-set version history for directory, branch, release, logical
  change, review thread, and workspace roles.
- Verify filenames are labels inside directory reference sets, not file-history
  identity fields.
- Verify all POSIX inode types are represented: regular file, directory,
  symbolic link, hard link, character device, block device, FIFO, and socket.
- Verify logical-change reference sets replace Jujutsu-style intrinsic change IDs
  in the POC18 protocol model.
- Verify grid-envelope parent links form the expected file, reference-set,
  snapshot, branch, merge, review, release, and supersedence chains.
- Verify reference-set promises are signed by the promiser and never treated as
  global authority.
- Verify sparse-CAS behavior by requiring at least one reference walk with missing
  local objects resolved through a peer promise.
- Verify merge behavior by requiring at least one multi-parent snapshot/change-set
  envelope.
- Verify review behavior by requiring at least one review promise that affects a
  local adoption decision without acting as a global approval.
- Verify GC and retention behavior by requiring at least one promised-retained
  object and at least one locally dropped unpromised object.
- Verify Git bridge behavior by importing/exporting and pulling/pushing a real
  Git repo while checking content, directories, branch/tag targets, and DAG
  semantics, plus explicit loss/refusal records for inode types Git cannot
  represent.
- Verify continuous peer DAG sync by requiring at least one useful update to
  propagate without an explicit push/pull command and without any global remote
  authority. Source: `DI-dibut`.
- Verify anti-RPC and Promise Theory vocabulary: no permission, authorization,
  conformance, command/control, or policy-enforcement framing unless explicitly
  reframed as local promises and local trust.

## Locked Decisions And Remaining Questions

Locked:

- POC18 is a protocol/CAS/kernel superset of POC16, not a runtime superset of
  POC17. Source: `DI-zuruj`.
- Versioned reference sets are the root abstraction for directories, filenames,
  tags, branches, releases, logical changes, review threads, and workspace roots.
  Source: `DI-zuruj`.
- Tags are primary fetch/discovery objects, not optional side data like Git's
  `--tags`. Source: `DI-zuruj`.
- Filenames belong to directory/reference-set history, not file-history identity.
  Source: `DI-zuruj`.
- POC18 must support all POSIX inode types: regular file, directory, symbolic
  link, hard link, character device, block device, FIFO, and socket. Source:
  `DI-radaj`.
- Logical-change reference sets replace Jujutsu-style intrinsic change IDs for
  POC18. Source: `DI-zuruj`.
- Use one main version-control pCID unless a later TE/DI proves a truly separate
  protocol boundary. Source: `DI-zuruj`.
- Use Rabin content-defined chunking for all files, and store chunks and Merkle
  manifests in PromiseGrid CAS formats. Source: `DI-dofoj`.
- Build PromiseGrid-native first, but POC18 is not complete until content and DAG
  Git bridge roundtrip works for import, export, push, and pull. Source:
  `DI-dofoj`.
- Use a Go Git library for the Git bridge. Source: `DI-zuruj`; `DI-dofoj`.
- Use POSIX nodes + snapshots: per-node lineage plus snapshot/change-set
  envelopes. Source: `DI-zuruj`; `DI-radaj`.
- Native synchronization is continuous peer DAG sync, not explicit Git-style
  push/pull. Source: `DI-dibut`.
- Tangled should influence POC18's self-hosting, migration, social-code UX,
  review-round, and stable-logical-change requirements, but POC18 should
  explicitly differ from Tangled's Git/SSH push-pull, appview aggregation,
  role-based access control, hidden Git refs, and raw Jujutsu change-ID
  mechanisms as native protocol concepts. Source: `DI-dibut`; `DN-dopod`.
- TE-kopap rejects a Git-first architecture baseline and a native-only
  architecture baseline. POC18's architecture baseline is a PromiseGrid-native
  reference-set/CAS core with required Git bridge import/export/push/pull as an
  adapter, not as the native authority model. Source: `TE-kopap`; `DI-zuruj`;
  `DI-dofoj`.
- TE-vahoj adds the post-TE-kopap hard-gate architecture and implementation
  sequencing analysis. Its recommended surviving path is deterministic local
  core first, but forced by tests and interfaces to preserve sparse multi-agent
  CAS assumptions from day one, with repo-like views over shared CAS and
  versioned reference sets. Its remaining DF questions must be answered before
  POC18 code scaffolding. Source: `TE-vahoj`; `DI-fusir`.
- TE-hikar locks the first `grid` CLI command model: balance familiar VCS
  terminology with PromiseGrid-native CAS, promise, and reference-set semantics;
  use `CLI/core` as the implementation split; use `poc-*` only for
  non-production deterministic fixtures. Source: `TE-hikar`; `DI-harih`.
- TE-givul analyzes whether POC18 chunks should remain raw `.bin` byte objects
  or move to CBOR `.cbor` objects. It recommends raw public chunk identity for
  deduplication, IPFS-like raw-block interop, and simple sparse retrieval, while
  keeping CBOR/grid messages responsible for manifests, promises, retention, and
  interpretation. It keeps encrypted/ciphertext chunk identity as a future
  explicit profile and leaves final DF open. Source: `TE-givul`; `DI-dofoj`.
- POC18 code scaffolding uses layer packages `store`, `graph`, `workspace`,
  `sync`, and `bridge`; generated first-slice runtime state lives only under
  `/tmp/wire-lab-poc18-*`; DevOps/root-filesystem tests are container-only and
  not part of the host first slice. Source: `DI-jifuj`; `DI-harih`.
- The POC18 version-control protocol spec is
  `implementations/poc18-cas-git-replacement/docs/protocols/version-control.md`,
  with pCID alias
  `bafkreicrikn3oqfumjnuvruw67h5ffvu6dyy7inz7h2rtm6s4qgwgz7oxu.md`. The
  filename intentionally has no `-v1` or `_v1` suffix because the pCID is the
  protocol version. Recompute it from the repo root with `cd tools/spec`, then
  `go run . cid ../../implementations/poc18-cas-git-replacement/docs/protocols/version-control.md`;
  in read-only Go-cache sandboxes, set `GOCACHE` and `GOMODCACHE` under `/tmp`.
  Source: `DI-lidaj`.

Remaining:

- Implement the first core graph slice under the locked `DI-jifuj` and
  `DI-harih` path, command, package, runtime, and CLI/core decisions.
- Choose the exact Go libraries for the Git bridge and Rabin chunking.
- Decide whether LLM agents participate in the first implementation slice or
  whether POC18 begins deterministic and adds LLM-scale collaboration later.

## Subtasks

- [x] nahop.16 Create `docs/research/DN-rifir-poc18-versioned-reference-sets.md`
  describing the tag/reference-set decisions, the discussion that led to them,
  and the Git/Jujutsu reasoning behind them. Completed by `DN-rifir`.
- [x] nahop.1 Run a TE comparing three designs: Git bridge compatibility,
  PromiseGrid-native CAS/version graph only, and hybrid Git bridge with
  PromiseGrid-native reference-set promises. Completed by
  `docs/thought-experiments/TE-kopap-poc18-git-bridge-vs-native-cas.md`.
- [x] nahop.2 Lock implementation paths, command names, package names, runtime
  artifact paths, and generated CAS path patterns before scaffolding
  `implementations/poc18-cas-git-replacement/`. Completed by `TE-hikar`,
  `DI-jifuj`, and `DI-harih`.
- [x] nahop.3 Write RFC-like spec docs for the POC18 version-control pCID before
  using its pCID in code. Completed by
  `implementations/poc18-cas-git-replacement/docs/protocols/version-control.md`;
  pCID alias
  `implementations/poc18-cas-git-replacement/docs/protocols/bafkreicrikn3oqfumjnuvruw67h5ffvu6dyy7inz7h2rtm6s4qgwgz7oxu.md`.
- [x] nahop.4 Implement per-agent sparse CAS object stores with CIDv1 base32
  printable paths and binary CID values inside CBOR. First slice completed by
  `implementations/poc18-cas-git-replacement/store/`. Source: `DI-jifuj`;
  `DI-harih`.
- [x] nahop.5 Integrate Rabin content-defined chunking and PromiseGrid Merkle
  manifests for all file content storage, including large in-band files. First
  slice completed by `implementations/poc18-cas-git-replacement/chunk/`.
  Source: `DI-dofoj`; `DI-jifuj`.
- [x] nahop.6 Implement POSIX node-version envelopes with logical node identity,
  node type, content or metadata payload, and parent node-version links. First
  slice completed by `implementations/poc18-cas-git-replacement/workspace/` and
  `implementations/poc18-cas-git-replacement/graph/`. Source: `DI-radaj`;
  `DI-harih`.
- [x] nahop.7 Implement directory reference sets where filename labels point at
  POSIX node-version or directory-version CIDs, including hard-link labels that
  intentionally share a node identity. First slice completed by
  `implementations/poc18-cas-git-replacement/workspace/`. Source: `DI-zuruj`;
  `DI-harih`.
- [x] nahop.8 Implement branch, release, logical-change, review-thread, and
  workspace reference-set roles without global authority. First slice completed
  by `implementations/poc18-cas-git-replacement/workspace/`. Source:
  `DI-zuruj`; `DI-harih`.
- [x] nahop.9 Implement snapshot/change-set envelopes that compose root directory
  reference sets and parent snapshot links. First slice completed by
  `implementations/poc18-cas-git-replacement/graph/` and
  `implementations/poc18-cas-git-replacement/workspace/`. Source: `DI-zuruj`;
  `DI-harih`.
- [x] nahop.10 Implement checkout/materialization from a workspace or snapshot
  reference set into a local workspace directory with explicit local promises.
  First slice completed by `implementations/poc18-cas-git-replacement/workspace/`.
  Source: `DI-radaj`; `DI-jifuj`.
- [ ] nahop.11 Implement peer fetch/retrieval of reference sets and missing CAS
  objects, including voluntary storage/forwarding promises and token incentives.
- [ ] nahop.12 Implement rename/copy scenarios that preserve node lineage while
  changing directory labels.
- [ ] nahop.13 Implement divergent branch/logical-change and multi-parent merge
  scenarios, including conflict-resolution promises.
- [ ] nahop.14 Implement review/test-result promises and local adoption decisions
  that can replace a GitHub pull-request approval flow without forge authority.
- [ ] nahop.15 Implement the Git bridge content-and-DAG roundtrip with a Go Git
  library, covering import, export, push, and pull through shared conversion
  paths.
- [ ] nahop.17 Add promise-based retention and GC behavior for selected reference
  sets, release objects, paid storage, and unpromised objects under pressure.
- [x] nahop.18 Review Tangled as prior art and record whether POC18 should adopt,
  reject, or explicitly differ from Tangled's self-hosted knots, appview,
  AT Protocol identity, round-based PR flow, SSH/Git compatibility, and Jujutsu
  change-ID use. Completed by
  `docs/research/DN-dopod-poc18-tangled-prior-art.md`. Source: `DI-dibut`.
- [ ] nahop.19 Implement continuous peer DAG sync so agents exchange
  reference-set heads, object-availability promises, missing-object requests, and
  retention/forwarding promises without explicit native push/pull commands.
  Source: `DI-dibut`.
- [ ] nahop.20 Add analyzer gates for CID correctness, sparse CAS, parent-chain
  integrity, reference-set signatures, multi-target labels, directory labels,
  logical-change reference sets, POSIX inode type coverage, continuous sync,
  Rabin chunking for large in-band files, Git bridge roundtrip, GC behavior, and
  anti-authority vocabulary.
- [ ] nahop.21 Add diagnostic rendering of representative raw CBOR messages for
  reference-set, node-version, directory, snapshot, review, merge,
  materialization, and peer-fetch flows.
- [ ] nahop.22 Run a clean deterministic POC18 scenario and archive exact commands,
  CAS object examples, reference-set walks, parent-chain walks, Git bridge
  output, and analyzer output.
- [x] nahop.23 Run the hard-gate POC18 superset architecture TE covering today's
  full planning thread, sparse shared CAS, promise-shaped retrieval, UI/backend
  boundaries, in-band collaboration, DevOps ordered replay, and first-slice
  sequencing. Completed by
  `docs/thought-experiments/TE-vahoj-poc18-superset-architecture.md`. Source:
  `DI-fusir`.
- [x] nahop.24 Run a follow-up CLI TE to balance Git/Jujutsu/Tangled/Mercurial/
  SVN/CVS familiarity against PromiseGrid-native grid messages, promises, sparse
  CAS, and reference-set semantics before locking first-slice `grid` verbs.
  Completed by `docs/thought-experiments/TE-hikar-poc18-grid-cli-command-model.md`.
  Source: `DI-harih`.
- [x] nahop.25 Run a chunk storage identity TE comparing raw chunks, CBOR
  chunks, grid-wrapped chunk promises, encrypted chunks, IPFS/IPLD, Git, Ceph,
  and related storage systems. Completed by
  `docs/thought-experiments/TE-givul-poc18-chunk-storage-identity.md`.
- [ ] nahop.26 Lock the follow-up DF from `TE-givul`: whether public/plain chunk
  CIDs name raw bytes, whether raw chunks remain under `chunks/*.bin` or move
  under `objects/`, whether encrypted chunks use ciphertext CIDs, and whether
  standalone manifest objects remain distinct from `chunk_manifest` messages.
