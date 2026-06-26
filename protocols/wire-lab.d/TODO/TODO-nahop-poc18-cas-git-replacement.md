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
  filename/path-component dirents and its targets are file-version or
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
fetches first when she wants a branch, release, logical change, directory,
review, or workspace. She then recursively asks peers for target CIDs, parent
reference-set versions, and parent-linked CAS objects she decides to retrieve.
Source: `DI-zuruj`.

## Git And Jujutsu Concepts To Preserve And Reframe

- **Blob/content:** Store bytes in CAS using restic's content-defined chunker.
  The normal file content root is a PromiseGrid Merkle manifest CID over chunk
  CIDs, not a restic repository object.
- **Tree/directory:** Reframe as a reference set with role `directory`.
  Filenames are labels/dirents inside the directory's history, like Unix
  directory entries pointing at inodes. A file's logical history does not own one
  canonical path.
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
- **Clone/fetch/push:** Reframe as reference-set promise exchange plus CAS object
  exchange. Peers promise which reference sets and CAS objects they have, what
  they are willing to retain or forward, and under what local constraints.

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
  files, directories, devices, or other local resources under local constraints.
- Preserve pCID discipline. pCID selects the protocol parser/builder and slot
  semantics; names, roles, labels, paths, authors, repositories, and destinations
  live in pCID-defined payloads where needed.
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
- Define parent-link slot rules explicitly. POC18 should prefer envelope-level
  parent links for version graph edges, with typed parent roles such as
  previous-file-version, previous-reference-set-version, previous-snapshot,
  merge-parent, review-parent, and supersedes.
- Define signature/proof placement explicitly. Reference-set, release, review,
  and adoption promises should be signed by the promiser. COSE or other proof
  shapes must be selected by the pCID spec rather than assumed globally.
- Define Git import/export requirements. POC18 starts PromiseGrid-native but is
  not complete until it can import and export a real Git repo while preserving
  content, directory structure, branch/tag targets, and parent DAG semantics.

## CAS And Chunking Targets

- Store exact bytes under CID-derived paths using CIDv1 base32 text when
  printable. No bare SHA-256 hex strings should appear as object identities.
- Use restic's content-defined chunker library for normal file content.
- Store chunks and Merkle/index manifests in PromiseGrid CAS formats, not restic
  repository format.
- Support at least these object classes:
  - raw chunks;
  - PromiseGrid Merkle manifests over chunk CIDs;
  - file-version envelopes;
  - directory/reference-set envelopes;
  - snapshot/change-set envelopes;
  - logical-change, branch, release, workspace, and review reference sets;
  - review/test/materialization promise envelopes;
  - optional Git import/export mapping records.
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

- **Single-author edit:** Alice edits `README.md`, restic-chunks its content,
  stores chunks and a Merkle manifest in CAS, emits a file-version envelope, then
  emits a snapshot/change-set envelope pointing at a root directory reference set.
- **Rename with history:** Alice renames `README.md` to `docs/intro.md`. The file
  version keeps its logical file lineage; the directory reference-set history
  changes the filename/dirent labels.
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
- **Git roundtrip:** Alice imports a real Git repo, maps Git branches/tags into
  reference-set promises, exports back to Git, and verifies content, directory
  structure, branch/tag targets, and DAG semantics.
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
- Verify Git roundtrip behavior by importing/exporting a real Git repo and
  checking content, directories, branch/tag targets, and DAG semantics.
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
- Logical-change reference sets replace Jujutsu-style intrinsic change IDs for
  POC18. Source: `DI-zuruj`.
- Use one main version-control pCID unless a later TE/DI proves a truly separate
  protocol boundary. Source: `DI-zuruj`.
- Use restic's content-defined chunker library, but store chunks and Merkle
  manifests in PromiseGrid CAS formats. Source: `DI-zuruj`.
- Build PromiseGrid-native first, but POC18 is not complete until content and DAG
  Git import/export roundtrip works. Source: `DI-zuruj`.
- Use a Go Git library for import/export. Source: `DI-zuruj`.
- Use Files + Snapshots: per-file lineage plus snapshot/change-set envelopes.
  Source: `DI-zuruj`.

Remaining:

- Lock implementation paths, command names, package names, runtime artifact
  paths, and generated CAS path patterns before scaffolding code.
- Choose the exact Go libraries for Git import/export and restic chunking.
- Decide whether LLM agents participate in the first implementation slice or
  whether POC18 begins deterministic and adds LLM-scale collaboration later.

## Subtasks

- [x] nahop.16 Create `docs/research/DN-rifir-poc18-versioned-reference-sets.md`
  describing the tag/reference-set decisions, the discussion that led to them,
  and the Git/Jujutsu reasoning behind them. Completed by `DN-rifir`.
- [ ] nahop.1 Run a TE comparing three designs: Git-compatible object import/export,
  PromiseGrid-native CAS/version graph only, and hybrid Git-import with
  PromiseGrid-native reference-set promises.
- [ ] nahop.2 Lock implementation paths, command names, package names, runtime
  artifact paths, and generated CAS path patterns before scaffolding
  `implementations/poc18-cas-git-replacement/`.
- [ ] nahop.3 Write RFC-like spec docs for the POC18 version-control pCID before
  using its pCID in code.
- [ ] nahop.4 Implement per-agent sparse CAS object stores with CIDv1 base32
  printable paths and binary CID values inside CBOR.
- [ ] nahop.5 Integrate restic content-defined chunking and PromiseGrid Merkle
  manifests for file content storage.
- [ ] nahop.6 Implement file-version envelopes with logical file identity, content
  root CIDs, and parent file-version links.
- [ ] nahop.7 Implement directory reference sets where filename labels point at
  file-version or directory-version CIDs.
- [ ] nahop.8 Implement branch, release, logical-change, review-thread, and
  workspace reference-set roles without global authority.
- [ ] nahop.9 Implement snapshot/change-set envelopes that compose root directory
  reference sets and parent snapshot links.
- [ ] nahop.10 Implement checkout/materialization from a workspace or snapshot
  reference set into a local workspace directory with explicit local promises.
- [ ] nahop.11 Implement peer fetch/retrieval of reference sets and missing CAS
  objects, including voluntary storage/forwarding promises and token incentives.
- [ ] nahop.12 Implement rename/copy scenarios that preserve file lineage while
  changing directory labels.
- [ ] nahop.13 Implement divergent branch/logical-change and multi-parent merge
  scenarios, including conflict-resolution promises.
- [ ] nahop.14 Implement review/test-result promises and local adoption decisions
  that can replace a GitHub pull-request approval flow without forge authority.
- [ ] nahop.15 Implement Git import/export content-and-DAG roundtrip with a Go Git
  library.
- [ ] nahop.17 Add promise-based retention and GC behavior for selected reference
  sets, release objects, paid storage, and unpromised objects under pressure.
- [ ] nahop.18 Add analyzer gates for CID correctness, sparse CAS, parent-chain
  integrity, reference-set signatures, multi-target labels, directory labels,
  logical-change reference sets, Git roundtrip, GC behavior, and anti-authority
  vocabulary.
- [ ] nahop.19 Add diagnostic rendering of representative raw CBOR messages for
  reference-set, file-version, directory, snapshot, review, merge,
  materialization, and peer-fetch flows.
- [ ] nahop.20 Run a clean deterministic POC18 scenario and archive exact commands,
  CAS object examples, reference-set walks, parent-chain walks, Git roundtrip
  output, and analyzer output.
