# cas_storage_v1

## Status

Active POC16 CAS protocol. The embedded Markdown bytes derive this pCID. Source:
`DI-bitug`; `DI-gahuh`; `DI-manul`.

## Abstract

`cas_storage_v1` lets agents promise content-addressed storage, retrieval,
replication, sparse presence checks, storage reports, repair intent, and
unsupported-variant non-commitments. CAS objects are exact bytes addressed by
content CID. Every CAS is local and partial; there is no complete shared store.

## pCID and envelope

The active shape is:

```text
grid([42(pCID), payload, proof])
```

The pCID chooses this CAS payload grammar. Content CIDs inside the payload name
stored objects; the pCID itself names this protocol spec.

## Promise Theory model

A storage provider promises only its own storage or serving behavior. A requester
may promise reciprocal credits or bearer capability tokens. Trust in storage is
built from local verification of exact bytes, token redemption, and promise
history. Failing to store because no promise was made is not a broken promise.

## Payload grammar

The payload is the strict pCID-owned array profile:

```text
[promiser, promisee, promise_about, [outcome, promise_text, reason], body]
```

`promise_about` values and body slots are:

| promise_about | Body slots |
|---|---|
| `store_content` | `storage_status`, `exchange_id`, `content_cid`, `content_b64`, `credit_offer`, `units`, `capability_token`, `token_style`, `bearer_token`, `replica_peer`, `replica_token` |
| `serve_content` | `token_status`, `exchange_id`, `content_cid`, `content_b64`, `token`, `token_style`, `missing_object_probe` |
| `replicate_content` | `content_cid`, `exchange_id`, `content_b64`, `issuee`, `units`, `replica_token` |
| `serve_replica_content` | `token_status`, `exchange_id`, `content_cid`, `content_b64`, `token`, `missing_object_probe` |
| `replica_token_lifecycle` | `token_status`, `exchange_id`, `content_cid`, `token`, `bearer_token`, `token_style`, `issuer_peer`, `redeem_peer` |
| `present_storage_report` | `verdict`, `exchange_id`, `content_cid`, `content_b64` |
| `label_future_malformed_report` | `repair_status`, `exchange_id` |
| `unsupported_variant_probe` | `variant_status`, `exchange_id` |

## Sender behavior

A sender promising storage SHOULD include `content_cid` and either `content_b64`
for bytes being supplied or a token/reference enabling retrieval. A sender MUST
not claim another store has bytes unless forwarding that store's exact promise or
reporting only its own local observation. Bearer tokens may be offered as
reciprocal payment only if signed and unexpired.

## Receiver and parser behavior

A parser MUST reject unknown `promise_about` values, incorrect body lengths,
wrong state shape, non-text slots, and trailing CBOR. A receiver keeping a
storage promise MUST verify that decoded `content_b64` hashes to `content_cid`
before treating it as stored. Missing objects result in local non-commitment, not
proof of bad faith unless a prior storage promise is locally outstanding.

## Protocol state machine

```text
[no local bytes]
    | store_content kept + CID verified
    v
[locally stored] --serve_content with valid token--> [served once or according to scope]
      | replicate_content kept                         |
      v                                                | retention/GC promise expires
[replica promised] --serve_replica_content--> [replica served]
      |                                                v
      +-------------------------------> [locally collectable / GC candidate]
```

Corrupt bytes move to local malformed handling; repair promises are future-only.

## State, CAS, DAG, and retention

Each agent keeps its own filesystem CAS in POC16. Objects are partial and
run-scoped unless an agent explicitly promises longer retention. Parent message
links may form DAGs of storage requests, reports, and token redemptions. GC and
backpressure SHOULD be represented as local promises about retention capacity,
not silent deletion that contradicts a prior promise.

## Security considerations

Storage tokens must be signed, scoped, unexpired, and replay-checked according to
local state. `content_b64` is untrusted until CID-verified. A malicious peer can
offer high credits with malformed bytes; recipients should verify before local
trust increases.

## Interoperability notes

Production implementations may store CAS bytes in files, databases, object
stores, or flash. The protocol relies only on exact byte/CID verification and
local promise history.

## Examples

```text
grid([42(pCID),
  ["bob", "alice", "store_content",
    ["kept", "I promise to store this object for this run in exchange for Alice storage credits.",
     "capacity available"],
    ["stored", "cas-44", "bafkrei...", "SGVsbG8=", "alice-credit:2",
     "2", "", "", "", "frank", "replica-token-b64"]
  ], proof
])
```
