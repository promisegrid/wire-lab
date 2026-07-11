# PromiseGrid POC20 pure-function protocol v1

## Status

POC20 pre-code protocol specification. This document is intentionally standalone
so its exact bytes can be named by a pCID. The pCID is the CID of this document,
not the CID of any function, input, output, message, or payload.

Source: `DI-lamaz`; `DI-mokaz`; `DI-lulog`; `DI-kodob`; `DI-ruvum`;
`TODO-nudav`; `TE-lodom`.

## Purpose

`pure-function-v1` defines promises about deterministic computation over explicit
content-addressed inputs. It lets an agent promise that for a function CID,
input CID, and context CID, it will return or did return a result CID. It also
defines local verification, disagreement, correction, and non-commitment records
without creating a compute authority. When the computation depends on fetched app
code, runtime executable objects, protocol specs, model bytes, or data roots,
those root CIDs are part of context.

## Envelope

Messages use this CBOR grid envelope:

```text
grid([42(pCID), payload, proof])
```

- `pCID` is the protocol CID of this specification document.
- `payload` is a CBOR map as defined below.
- `proof` is a single-signer COSE_Sign1-style proof over
  `grid([42(pCID), payload])`.

## Payload map

The payload is a CBOR map. Keys are UTF-8 strings. CIDs are encoded as CBOR
tag-42 links in CBOR and canonical CIDv1 base32 in diagnostics.

Common keys:

| key | required | meaning |
| --- | --- | --- |
| `record_type` | yes | One of the record types below. |
| `promiser` | yes | Agent making the promise. |
| `promisee` | no | Agent or group the promise is addressed to, if any. |
| `parents` | yes | Parent record CIDs linking this record into a timeline. |
| `function` | yes | CID of the function definition or executable object. |
| `input` | yes | CID of the explicit input object. |
| `context` | yes | CID of the explicit context object. |
| `result` | conditional | CID of the promised or observed result, required for result records. |
| `body` | yes | Record-type-specific map. |
| `reciprocal` | no | Requested reciprocal promise, such as payment, verification, retention, or receipt. |
| `local_constraints` | no | Promiser-local bounds such as time, cost, capacity, or supported runtime. |

Record types:

- `result_promise`: the promiser promises that tuple `(function, input, context)`
  maps to `result`.
- `compute_offer`: the promiser promises willingness to compute a tuple under
  stated local constraints if a reciprocal promise is acceptable.
- `non_commitment`: the promiser states that it does not currently promise to
  compute or verify the tuple.
- `verification_promise`: the promiser promises it recomputed or checked the
  tuple and locally observed the stated result.
- `disagreement_promise`: the promiser promises it got a different result or
  could not reproduce the tuple under the stated context.
- `correction_promise`: the promiser narrows or supersedes an earlier result
  promise while preserving the earlier record in the CAS timeline.

## Determinism rule

For the same `function`, `input`, and `context`, a pure-function agent promises
the same `result`. Any clock reads, random values, sensor samples, exchange-rate
quotes, peer promises, model versions, runtime versions, app root CIDs, runtime
root CIDs, protocol-spec CIDs, data root CIDs, or hardware-specific facts that
affect output must be represented inside `context`.

Replay source keys and generated action hashes are not hidden compute inputs. If
they affect the function result, retry behavior, or downstream action being
promised, they must appear in `context` as ordinary CAS objects. A changed action
hash for the same source key is a different replay situation unless a later local
timeline record proves a correction or explicit sequence extension.

If Bob returns two different result CIDs for the same tuple, receivers do not
need a central judge. Each receiver can retain both records and make a local
trust decision.

## Runtime and app roots in context

The installed microkernel binary is not enough context for a pure-function
result if fetched code or runtime objects affect the result. The context object
must name the relevant app root, runtime root, executable object, protocol-spec
root, model root, and data root CIDs. A later app/runtime update is therefore a
different context unless the function spec explicitly proves that the update is
irrelevant.

Operator approval of a root is a timeline promise, not a hidden compute input.
If a result depends on Alice's adopted root rather than Bob's adopted root, that
root CID must be explicit in the context.

## Behavior

An agent may offer compute, decline to promise compute, promise a result, verify
another result, or correct its own earlier result. None of these records makes a
promise on behalf of another agent.

Result objects and context objects are ordinary CAS objects. Their shareability
is not implied by the result promise. A result promise may name private context
objects that the promiser does not currently promise to send.

## State machine

```text
        +------------------+
        | tuple unknown    |
        +---+----------+---+
            |          |
 compute_offer     non_commitment
            |          |
            v          v
        +------------------+
        | tuple considered |
        +---------+--------+
                  |
                  | result_promise
                  v
        +------------------+
        | result promised  |
        +----+--------+----+
             |        |
 verification|        | disagreement
             v        v
        +---------+  +----------------+
        | checked |  | conflict kept  |
        +----+----+  +-------+--------+
             |               |
             | correction    | correction or branch merge
             v               v
        +--------------------------+
        | corrected timeline state |
        +--------------------------+
```

## Failure handling

- Missing tuple CIDs, invalid proof, invalid CID encoding, malformed CBOR, or
  unknown `record_type` makes the message not promised by this protocol.
- A result without a locally available function, input, context, or result object
  may be retained as incomplete and repaired by requesting missing CAS objects.
- A context object that names unavailable app/runtime roots is incomplete. A
  receiver may request those roots, retain the record as partial, or decline to
  verify the tuple.
- A context object that omits a replay source key or action hash that affected
  the promised result is incomplete for deterministic replay.
- A verification record is local to the verifier. It is not a global verdict.

## Security and privacy

The proof authenticates the promiser's exact message bytes. It does not prove the
function was actually run, that the result is true, or that the result should be
trusted. Those judgments remain local and relationship-relative.

## POC20 acceptance requirements

A POC20 implementation of this protocol must demonstrate:

- one result promise over explicit function, input, and context CIDs;
- one independent verification promise by another agent;
- one disagreement or correction branch;
- local CAS retention of the tuple objects needed to replay the local decision;
- no ambient timestamp, randomness, sensor, model, peer quote, or exchange-rate
  input outside the context object;
- at least one result whose context explicitly names the app/runtime root CIDs
  used to produce or verify the result;
- at least one result whose context names a replay source key and generated
  action hash when those facts affect retry or downstream-action behavior.
