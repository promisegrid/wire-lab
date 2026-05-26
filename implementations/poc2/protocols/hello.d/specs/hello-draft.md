# poc2 hello protocol draft

> Status: DRAFT. Proof-of-concept protocol for `implementations/poc2`.

This protocol is intentionally tiny. Its Protocol CID is the CIDv1 raw-codec
SHA-256 digest of this spec document. Messages that use this pCID carry CBOR map
payloads under the current PromiseGrid envelope direction:

```text
grid([42(pCID), payload, ...])
```

## Payload kinds

### `receive_promise_v1`

An app promises its local kernel that it will receive hello messages for a bounded
local run.

Fields:

- `kind`: `receive_promise_v1`
- `from`: app agent name
- `node`: local node name
- `text`: human-readable promise text

### `hello_v1`

An app sends a hello promise-message to another node.

Fields:

- `kind`: `hello_v1`
- `from`: sender agent name
- `to`: destination node name
- `text`: hello text

### `observation_v1`

A kernel records and returns an observer-local outcome for one message attempt.

Fields:

- `kind`: `observation_v1`
- `from`: observer agent name
- `to`: local promisee or sender agent name
- `outcome`: `kept`, `refused`, or `not-promised`
- `text`: observer-local explanation
- `evidence_hash`: SHA-256 hex of exact envelope bytes when available

## Promise Theory discipline

No message in this protocol commands another agent. Kernels and apps each promise
only their own local behavior. Evidence records are local observations, not global
authority.
