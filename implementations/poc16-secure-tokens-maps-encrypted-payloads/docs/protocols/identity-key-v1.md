# identity_key_v1

## Status

Active POC16 key-continuity protocol. Its embedded Markdown bytes derive the
pCID. Source: `DI-bitug`; `DI-vipih`.

## Abstract

`identity_key_v1` carries promises about a promiser's future signing identity. It
supports key-rotation promises and acknowledgements. It is not a global identity
authority and does not tell other agents which key they must trust.

## pCID and envelope

The active shape is:

```text
grid([42(pCID), payload, proof])
```

The payload is a pCID-owned CBOR array, not the common pair or strict operation
array profile.

## Promise Theory model

The promiser promises that it intends to use or recognize a key label within a
stated scope. The promisee locally decides whether that continuity promise is
sufficient. No agent can rotate another agent's key.

## Payload grammar

Two payload shapes are defined:

```text
rotate_signing_key = [
  promiser: text,
  promisee: text,
  "rotate_signing_key",
  [new_key_label: text, rotation_scope: text]
]

rotate_signing_key_ack = [
  promiser: text,
  promisee: text,
  "rotate_signing_key_ack",
  [outcome: text, promise_text: text, new_key_label: text, rotation_scope: text]
]
```

All strings are required. `rotation_scope` SHOULD name the relationship, pCID
set, or runtime scope affected by the key promise.

## Sender behavior

A sender MUST sign the envelope with the current key that the receiver is expected
to recognize. It SHOULD parent-link to prior key promises where possible. ACKs
must say only whether the ACK sender locally accepts, rejects, or remembers the
rotation.

## Receiver and parser behavior

A parser MUST reject payload arrays that are not four slots, unknown promise kind
strings, wrong body lengths, non-text slots, or trailing CBOR. A receiver SHOULD
store the exact signed rotation message before changing local trust in future
signatures.

## Protocol state machine

```text
[current key locally recognized]
    | rotate_signing_key promise received
    v
[rotation pending] --ack kept / local acceptance--> [new key locally recognized]
        | signature mismatch / unacceptable scope
        v
[rotation not promised locally]
```

## State, CAS, DAG, and retention

Receivers SHOULD retain key-rotation parents as long as messages signed under the
new key remain relevant. CAS retention is local; losing key history may reduce an
agent's ability to verify old promises.

## Security considerations

Key rotation is vulnerable to replay, stripping, and ambiguous scope. Receivers
SHOULD verify parent chains and reject rotations outside locally trusted
relationships. This POC uses deterministic keys for repeatable tests; production
must use real key management.

## Interoperability notes

Future identity/key pCIDs may adopt COSE headers, DIDs, hardware keys, or
relationship-specific key-agreement promises. This spec only fixes the POC16
array grammar.

## Examples

```text
grid([42(pCID),
  ["alice", "bob", "rotate_signing_key", ["alice-key-2026q3", "relationship:alice-bob"]],
  proof
])
```
