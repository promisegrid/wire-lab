# PromiseGrid POC20 capability-token protocol v1

## Status

POC20 pre-code protocol specification. This document is intentionally standalone
so its exact bytes can be named by a pCID. The pCID is the CID of this document,
not the CID of any token, message, or payload.

Source: `DI-lamaz`; `DI-mokaz`; `DI-lulog`; `DI-kodob`; `TODO-nudav`;
`TE-lodom`.

## Purpose

`capability-token-v1` defines promises for issuing, transferring, redeeming, and
locally interpreting capability tokens as CAS timeline objects. It specifically
tests whether double-spend can be represented as branch-visible promise history
rather than hidden mutable projection state. It also provides the token vocabulary
for local promises to fetch, retain, decrypt, or execute objects from an
operator-adopted app/runtime root.

Tokens are issuer promises. A token does not command a resource owner. A token
gives a receiver bytes that an issuer promised to honor under pCID-defined terms,
and each recipient remains free to keep, ignore, redeem, discount, or distrust
the token based on local history.

## Token object

POC20 tokens are CWT/COSE-style signed CBOR objects stored as exact CAS bytes.
The token object's CID is referenced by protocol messages. The token signature is
the issuer's proof for the token itself; envelope proofs authenticate the
surrounding protocol message.

Minimum token claims:

| key | meaning |
| --- | --- |
| `issuer` | Agent that promises the token terms. |
| `subject` | Resource, service, branch family, or promise family the token is about. |
| `audience` | Intended redeemer or group, if holder-bound; omitted for bearer tokens. |
| `bearer` | Boolean; true means transferable by possession. |
| `not_before` | Optional earliest valid time or event CID. |
| `expires` | Optional expiry time or event CID. |
| `terms` | Map of pCID-defined service terms, cost, quantity, or constraints. |
| `token_family` | Identifier used to detect incompatible redemptions. |

For app/runtime roots, `subject` may name a root CID, app reference-set CID,
runtime profile CID, executable object CID, data root CID, or protocol-spec CID.
The token still does not command execution. It only records what the issuer
promises to honor when the token is presented.

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
| `promiser` | yes | Agent making this protocol promise. |
| `promisee` | no | Intended receiver or group, if any. |
| `parents` | yes | Parent record CIDs linking this message into a timeline. |
| `token` | yes | CID of the exact CWT/COSE token object. |
| `token_family` | yes | Family identifier used for branch-local double-spend analysis. |
| `branch_id` | yes | Local branch on which this token record is asserted. |
| `body` | yes | Record-type-specific map. |
| `reciprocal` | no | Requested reciprocal promise such as service, payment, storage, or receipt. |
| `local_constraints` | no | Promiser-local limits such as expiry, capacity, cost, or retention. |

Record types:

- `issue_promise`: issuer promises token terms and stores the token object in
  local CAS.
- `transfer_promise`: current holder promises to transfer a bearer token or
  identifies why a non-transferable token is not promised as transferable.
- `redeem_promise`: redeemer promises to present a token for service under local
  branch context.
- `redemption_result`: service agent promises how it locally handled the token:
  kept, not currently promised, expired, revoked, already seen on this branch, or
  conflicting with another branch.
- `status_promise`: issuer or holder promises current local status for a token
  without commanding other agents to agree.
- `merge_decision`: local decision to keep, reject, merge, compensate, or leave
  unmerged branches containing token redemptions.

## Fetch and execution capabilities

A token may promise access to root-graph retrieval, encrypted-object decryption,
storage/retention, or local runtime resources. Examples include:

- Alice promises Bob may fetch a selected app root closure from Alice's CAS.
- Alice promises a local runtime role may execute a selected WASI module under
  bounded CPU, memory, time, host-function, and network terms.
- Bob promises to retain a runtime root for Alice in return for a reciprocal
  payment or future service promise.

These are promises by the issuer about resources the issuer controls. They are
not globally granted access rights and not instructions to a runtime or peer.

## Double-spend model

A bearer token can be copied. POC20 does not pretend otherwise. If Mallory
presents the same token to Bob on branch `B` and Carol on branch `C`, both
redemption promises may exist as CAS objects. The conflict is visible because the
records name the same `token`, the same `token_family`, and incompatible parent
branches.

Whether that is a broken promise depends on token terms and local interpretation.
If the token terms promise one successful redemption per branch family, Bob and
Carol may lower local trust in Mallory when both branches are discovered. They do
not need to lower trust in Alice merely because a bearer token was copied.

## State machine

```text
        +------------------+
        | no token object  |
        +---------+--------+
                  |
                  | issue_promise
                  v
        +------------------+
        | token retained   |
        +----+--------+----+
             |        |
 transfer    |        | redeem_promise
             v        v
        +---------+  +----------------+
        | holder  |  | redemption seen|
        +----+----+  +-------+--------+
             |               |
             | second branch | redemption_result
             v               v
        +--------------------------+
        | branch conflict visible  |
        +-------------+------------+
                      |
                      | merge_decision
                      v
        +--------------------------+
        | local interpretation     |
        +--------------------------+
```

## Failure handling

- Missing token CID, missing token family, invalid proof, malformed CBOR, or
  unknown `record_type` makes the message not promised by this protocol.
- A token whose exact CWT/COSE bytes are unavailable is incomplete. Receivers may
  request the token object, keep an incomplete branch, or ignore it.
- Expired tokens are expired, not stale. Presenting an expired token as useful may
  affect trust in the presenter, not necessarily the issuer.
- Revocation is an issuer promise about issuer-local future behavior. It is not a
  global command to other agents.

## Security and privacy

Token signatures authenticate issuer token bytes. Envelope proofs authenticate
protocol messages about those token bytes. Neither proof creates a global
authority. Each receiver decides locally whether the issuer, presenter, service
agent, branch, and reciprocal promise are worth trusting.

## POC20 acceptance requirements

A POC20 implementation of this protocol must demonstrate:

- one signed CWT/COSE token object retained in local CAS;
- one token whose subject is an app/runtime/root object used by the first
  executable slice;
- one bearer transfer;
- two redemption promises on parallel branches;
- one local merge or non-merge decision after conflict discovery;
- projection rebuild that recovers token state from CAS rather than a hidden
  spent-token table.
