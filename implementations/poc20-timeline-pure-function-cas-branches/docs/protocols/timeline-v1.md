# PromiseGrid POC20 timeline protocol v1

## Status

POC20 pre-code protocol specification. This document is intentionally standalone
so its exact bytes can be named by a pCID. The pCID is the CID of this document,
not the CID of any message or payload.

Source: `DI-lamaz`; `DI-mokaz`; `DI-lulog`; `DI-kodob`; `TODO-nudav`;
`TE-lodom`.

## Purpose

`timeline-v1` defines promise-shaped records for local timelines, group
timeline agreements, branch heads, merge decisions, shareability decisions,
bootstrap root adoption, root updates, and projection checkpoints. It exists so
POC20 can test the rule that local CAS is the chronological event source and
source of truth, while every JSON file, SQLite table, in-memory map, or other
local view is only a rebuildable projection.

This protocol does not define global truth, global branch authority, global
trust, or global monitoring. Every record is a promise by its promiser from that
promiser's local vantage.

## Envelope

Messages use this CBOR grid envelope:

```text
grid([42(pCID), payload, proof])
```

- `pCID` is the protocol CID of this specification document.
- `payload` is a CBOR map as defined below.
- `proof` is a pCID-owned proof object. For POC20, it is a single-signer
  COSE_Sign1-style proof over the pCID-defined signable view
  `grid([42(pCID), payload])`.

Parent links are payload fields in this protocol. They are not universal
envelope slots.

## Payload map

The payload is a CBOR map. Keys are UTF-8 strings. CIDs are encoded as CBOR
tag-42 links when carried in CBOR and as canonical CIDv1 base32 text when
rendered in diagnostics.

Common keys:

| key | required | meaning |
| --- | --- | --- |
| `record_type` | yes | One of the record types below. |
| `promiser` | yes | Local identifier of the agent making this promise. |
| `promisee` | no | Intended receiver or group, if the promiser chooses to name one. |
| `timeline_id` | yes | Local or group timeline identifier chosen by the promiser or group. |
| `branch_id` | yes | Branch identifier within `timeline_id`. |
| `parents` | yes | Array of parent record CIDs. Empty only for a genesis record. |
| `subject` | yes | Map naming what part of the universe the promise is about. |
| `body` | yes | Record-type-specific promise body. |
| `valid_time` | no | Map describing the time or interval the promise is about. |
| `context` | no | Array or map of context object CIDs the promiser wants receivers to consider. |
| `reciprocal` | no | Promise requested or expected from another agent, without making that promise on their behalf. |
| `local_constraints` | no | Promiser-local limits such as retention, sharing, cost, capacity, or expiry. |

Record types:

- `timeline_event`: a promise that something was asserted, observed, retained,
  withheld, or corrected at a point in the promiser's local timeline.
- `branch_head`: a promise that a named branch currently has a given head CID
  from the promiser's local vantage.
- `group_timeline_promise`: a voluntary promise to maintain, interpret, or sync
  a named group branch with one or more peers.
- `merge_decision`: a local promise to keep, reject, merge, compensate, or leave
  unmerged one or more branches.
- `shareability_promise`: a local promise classifying one CAS object or branch
  as `private`, `encrypted_shareable`, or `plain_shareable`.
- `root_adoption`: a local promise that the promiser currently adopts a
  bootstrap, app, agent, runtime, protocol-spec, or data root CID for some local
  purpose.
- `root_update`: a local promise that a later root CID supersedes, narrows,
  forks, or is rejected relative to an earlier adopted root.
- `projection_checkpoint`: a promise that a derived local projection was rebuilt
  from a named CAS head at a named time.

## Shareability classes

Shareability is a promise, not a property implied by possession:

- `private`: the promiser does not currently promise to send the object.
- `encrypted_shareable`: the promiser may send ciphertext or encrypted bundles
  under pCID-defined recipient promises.
- `plain_shareable`: the promiser may send the exact object bytes to selected
  peers under local promises.

The existence of an object in local CAS never implies sendability.

## Bootstrap roots and updates

A root CID names a Merkle/CAS object graph. It does not by itself promise trust,
safe execution, compatibility, freshness, or authority. A `root_adoption` record
is the promiser's local promise that the named root is the root the promiser is
currently willing to use for a stated purpose, such as app discovery, runtime
profile lookup, protocol-spec lookup, or default data roots.

A `root_update` record links a prior adopted root to a later candidate or adopted
root. The update body states whether the promiser locally accepts, rejects,
forks, narrows, or is still evaluating the candidate root. Operator approval is
represented as a local promise by the operator or local node role. No update
record commands another agent to adopt the same root.

## Behavior

An agent that receives a valid `timeline-v1` message may store the exact envelope
bytes in local CAS, store the payload object in local CAS, or ignore the message.
Keeping a record does not mean accepting it as truth. It means the receiver has
retained a promise made by the promiser.

Branch heads and adopted roots are local promises. If Alice and Bob both promise
a branch head or app/runtime root for the same group, neither head is globally
authoritative. Receivers compare parents, promisers, proofs, and local trust
history before deciding what to keep, merge, or ignore.

Projection checkpoints are diagnostic promises about rebuildable local views.
They do not create source-of-truth state. If a projection checkpoint conflicts
with replayed local CAS, local CAS wins.

## State machine

```text
        +------------------+
        | no local object  |
        +---------+--------+
                  |
                  | receive or create valid promise
                  v
        +------------------+
        | retained in CAS  |
        +----+--------+----+
             |        |
             |        | shareability_promise
             |        v
             |  +----------------------+
             |  | share class recorded |
             |  +----------+-----------+
             |             |
             | branch_head | merge_decision
             v             v
        +------------------------+
        | local branch projected |
        +-----------+------------+
                    |
                    | delete projection and replay CAS
                    v
        +------------------------+
        | projection rebuilt     |
        +------------------------+
```

## Failure handling

- Malformed CBOR, missing required keys, invalid CIDs, invalid proof, or unknown
  `record_type` makes the message not promised by this protocol.
- A receiver may keep malformed bytes as private local CAS material for debugging
  only if local policy allows; malformed bytes are not a valid timeline record.
- A received branch with missing parents is incomplete. The receiver may request
  missing objects, keep the branch as partial, or ignore it.
- A received root with missing closure objects is incomplete. The receiver may
  request missing objects, retain only the root record, or decline adoption.
- A branch conflict is represented by parent-linked records, not by hidden
  mutation of a spent-token table or other projection.

## Security and privacy

The proof authenticates the promiser's message bytes under this protocol. It does
not prove that the promise is true, useful, current, complete, or globally
accepted. Privacy is local: an agent may retain private CAS objects forever and
may refuse to send objects it does not currently promise to share.

## POC20 acceptance requirements

A POC20 implementation of this protocol must demonstrate:

- at least one local timeline per participating agent;
- at least one voluntary group timeline promise;
- at least one local root-adoption record and one local root-update record;
- at least one branch conflict and local merge/non-merge decision;
- at least one private, encrypted-shareable, and plain-shareable CAS object;
- deletion and rebuild of at least one projection from local CAS.
