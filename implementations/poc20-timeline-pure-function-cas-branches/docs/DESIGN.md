# POC20 timeline pure-function CAS branches design

## Status

Decision-complete pre-code design. This is the single human-readable POC20
design entrypoint, but protocol pCIDs require standalone hashable spec documents
under `docs/protocols/`. This document is not executable code and not a
production API. The canonical task record remains `TODO-nudav`. Source:
`DI-kakos`; `DI-bibah`; `DI-mokaz`; `DI-lamaz`; `DI-lulog`; `DI-kodob`;
`DI-ruvum`; `TE-lodom`; `TODO-nudav`.

## Locked direction

POC20 tests the semantic model that POC19 should not accidentally block:
promises as assertions about part of the universe on a timeline, agents as
deterministic pure-function servers over explicit content-addressed context,
decentralized CAS object chains as local or group timelines, and
capability-token double-spend as visible branch history rather than hidden
mutable projection state.

The first executable POC20 implementation must be a hybrid:

- reuse POC16 runtime lessons for pCID-selected parser/builder roles,
  grid-envelope handling, CWT/COSE token handling, proof checks, and
  length-framed TCP transport;
- reuse POC18 CAS lessons for CID handling, local filesystem CAS, parent-linked
  graph objects, sparse per-agent stores, diagnostic CBOR, and continuous sync
  behavior;
- do not implement cross-agent behavior through in-process fixture calls or
  shared mutable state.

## Protocol families

POC20 uses three protocol-family pCIDs. Message variants are payload semantics,
not separate pCIDs.

| family | slug spec | pCID alias | purpose |
| --- | --- | --- | --- |
| timeline | `docs/protocols/timeline-v1.md` | `bafkreihfdasban663gaabn7rtbxionkz7pnenf6t5uro27jd3hij5npqam.md` | local/group timelines, branch heads, merge/non-merge, shareability, root adoption/update, projection checkpoints, projection conflicts, replay decisions |
| pure function | `docs/protocols/pure-function-v1.md` | `bafkreigk2w2frnyh5dftaaftbarjdnl2gofzlepm5kfzmppe6cchui2l6i.md` | deterministic function/input/context/result promises, app/runtime root context, verification, disagreement, correction |
| capability token | `docs/protocols/capability-token-v1.md` | `bafkreie4prt2erwwwjyjltwm273oydhu7z22gbwixvnfdtrtaosegaeg6i.md` | CWT/COSE token issue, transfer, redemption, root fetch/execution access, local status, branch-visible double-spend |

All three use `grid([42(pCID), payload, proof])` for the first slice. Payloads
are CBOR maps because POC20 is testing semantic clarity more than constrained
device compactness. Parent links live in pCID-defined payload fields.

## Core state model

Every important POC20 event is a promise-shaped CAS object with parent links. An
agent's local CAS is its chronological event source and source of truth. Local
runtime indexes, projections, caches, JSON files, SQLite tables, or in-memory
summaries are disposable views rebuilt from local CAS. If a projection and local
CAS disagree, local CAS wins.

Local CAS is not automatically public. Some local CAS objects may remain private
forever, some may be sent only after encryption, and some may be plain-shareable.
Sending, withholding, or encrypting a CAS object is a separate local promise
decision. Source: `DI-mokaz`.

## Bootstrap roots and operator adoption

POC19's minimum microkernel rule has a POC20 timeline meaning. A bootstrap root
CID is a CID for a Merkle/CAS object graph containing app reference sets,
runtime profiles, executable objects, protocol specs, data roots, and update
metadata. The CID names bytes, not trust or authority.

Adopting a root is a local timeline promise. Alice may adopt root `A`, Bob may
adopt root `B`, and Carol and Dave may voluntarily promise to track root `C` as a
group. A later root update is another timeline event that supersedes or narrows
the earlier adoption event from the local promiser's vantage. No root adoption
event makes a promise on behalf of another agent.

If a pure-function result depends on fetched app code, runtime behavior, protocol
specs, model bytes, or data roots, the relevant root CIDs must be explicit in
the pure-function context object. Source: `DI-lulog`; `DI-kodob`.

Root decision bodies must preserve enough local context for replay: decision kind
such as adopted, rejected, superseded, rollback, or still-evaluating; current
root; candidate root; rollback root; closure summary CID; signer summary CID;
impact summary CID; accepted or rejected capability changes; approving local
identity or role; local decision state; and local reason. The record is still
only a local promise. It does not make the root safe for another agent or require
another agent to adopt it. Source: `DI-ruvum`.

## Projection conflicts, replay keys, and sensitive data

Projection conflicts are local timeline records, not hidden projection-only
state. A conflict record names the source system or peer, field or field group,
current owner or promiser if known, observed value, reviewed value if one exists,
projected value if one exists, raw supporting evidence CID or CIDs, local
decision, local reason, and resulting projection or action. The default posture
is warn-and-decide: retain the conflict, avoid silent overwrite, and let a local
promise resolve or leave the conflict open. Source: `DI-ruvum`.

Replay identity uses stable upstream facts, root events, or timeline events as
primary source keys. A generated action payload hash is only a secondary guard
that a retry is attempting the same approved action. If the same source key
appears with a different action hash, POC20 must require a correction event, new
local approval, or explicit sequence key instead of silently replaying the old
write. Source: `DI-ruvum`.

Projection rebuild is a proof obligation, not an implementation detail. At least
one realistic projection must be deleted and rebuilt from local CAS, then
compared with the prior view for the same root decisions, approvals, conflict
summaries, token status, replay indexes, and replayable actions. Useful
projection examples include review queues, root adoption status, conflict
summaries, token status, prior-action warnings, and replay indexes. Source:
`DI-ruvum`.

Sensitive payloads should not be embedded unnecessarily in broad replicated
timeline records. Ordinary non-sensitive facts may live in the main timeline.
Sensitive payloads should live in private or encrypted CAS namespaces with
narrower retention and sharing promises, and broad summaries should refer to
opaque handles, encrypted-object CIDs, or keyed commitments. Avoid plain hashes
of guessable sensitive data. Local deletion, redaction, or cryptographic erasure
is represented by tombstone or erasure local events when policy requires it.
Source: `DI-ruvum`.

## First executable scenario

The first executable slice is one unified run with Alice, Bob, Carol, Dave,
Ellen, and Mallory:

- Alice first adopts bootstrap root `A`, then later considers an update root `A2`
  after fetching and verifying the CAS closure; she records an impact summary,
  local decision reason, and rollback root.
- Alice issues a signed bearer capability token as a CAS object and records an
  `issue_promise` on her local timeline.
- Bob offers a pure-function service and promises a result for explicit function,
  input, and context CIDs, including the app/runtime root CIDs that affect the
  result.
- Dave independently verifies Bob's function result and records either a
  `verification_promise` or `disagreement_promise`.
- Carol and Dave voluntarily maintain a small group timeline branch for selected
  token and compute history.
- Mallory presents the same bearer token on two parallel branches, one to Bob and
  one to Carol.
- Ellen receives both branches and records a local merge, non-merge, or
  compensation promise without deciding for other agents.
- Ellen rebuilds a review projection from CAS, detects a conflict between a
  reviewed value and a projected value, and records a local conflict decision
  without silently overwriting either value.
- At least one sensitive payload is stored only as private or encrypted CAS data,
  while the broad timeline records carry only an opaque handle,
  encrypted-object CID, or keyed commitment.

The run succeeds only if the token conflict is visible as branch history in CAS
and can be explained without a hidden mutable spent-token table or
projection-only table as the source of truth.

## Future code paths and names

Future code generation should stay under
`implementations/poc20-timeline-pure-function-cas-branches/`.

Approved future package names:

- `protocol` for CBOR grid envelope, payload, proof, and pCID helpers;
- `store` for local filesystem CAS and CID rendering;
- `timeline` for event objects, branch heads, merge decisions, and projections;
- `function` for pure-function tuple objects and verification;
- `token` for CWT/COSE token objects and redemption records;
- `runtime` for POC16-derived parser/builder, transport, and local role wiring;
- `scenario` for the deterministic first run;
- `analyzer` for post-run checks that rebuild from CAS.

Approved future command names:

- `poc20-sim` for the deterministic clean run;
- `poc20-analyze` for acceptance checks;
- `poc20-cbor-diag` for raw CBOR rendering.

Approved runtime root pattern:

```text
/tmp/wire-lab-poc20-run/<run_id>/
```

Future code may create per-agent CAS under:

```text
/tmp/wire-lab-poc20-run/<run_id>/agents/<agent>/cas/objects/<cid>.cbor
/tmp/wire-lab-poc20-run/<run_id>/agents/<agent>/projections/*.json
/tmp/wire-lab-poc20-run/<run_id>/diagnostics/*.jsonl
```

Projection files are rebuildable and must be safe to delete between runs.

## Analyzer gates before code is complete

The future `poc20-analyze` must fail the run unless it proves:

- every cross-agent message was a promise-shaped grid CBOR message over TCP;
- each retained semantic record is present in at least one agent-local CAS;
- every branch head can be reached by following parent CIDs in local CAS;
- at least one bootstrap root adoption and one root update consideration are
  represented as local timeline promises;
- root-decision records include decision kind, current root, candidate root,
  rollback root, closure summary, signer summary, impact summary, capability
  changes, approving identity or role, decision state, and local reason;
- the non-CAS projection directory can be deleted and rebuilt from local CAS;
- the rebuilt projection reports the same root decisions, approvals, conflict
  summaries, token status, replay indexes, prior-action warnings, replayable
  actions, and pure-function tuple outcomes as the original projection;
- at least one projection conflict record is rebuilt from CAS and preserves the
  source or peer, field group, reviewed value, projected value, raw supporting
  evidence CIDs, local decision, local reason, and resulting projection or
  action;
- replay checks distinguish a stable source key from a generated action hash and
  require a new local decision when the action hash changes for the same source
  key;
- at least one private, encrypted-shareable, and plain-shareable CAS object is
  represented by `timeline-v1` shareability promises;
- broad timeline summaries avoid plaintext secrets and unnecessary sensitive
  personal payloads;
- double-spend appears as two branch records for the same token family, not as
  hidden mutation of a spent-token table;
- no analyzer result requires a global trust authority, global branch authority,
  hidden monitor, or shared mutable fixture state.

## Non-goals

- No final PromiseGrid token API, app API, storage profile, or production pCID
  registry.
- No requirement that POC19 wait for POC20.
- No pCID-per-message-kind fragmentation.
- No global app store, package registry, update authority, or root authority.
- No code generation until `TODO-nudav` records this pre-code plan as complete.
