# relationship_v1

## Status

Active POC16 protocol-family specification. This document is embedded into the
POC16 binary and its exact bytes are hashed to derive the pCID for
`relationship_v1`; editing this file intentionally changes that pCID. Source:
`DI-bitug`.

## Abstract

`relationship_v1` carries promises about local relationships: trust changes,
peer introductions, local event summaries, repair intent, and non-commitments.
It is the default human-readable relationship protocol for POC16 agents. It is
not an identity authority, reputation oracle, permission system, or command
surface.

## pCID and envelope

Messages using this protocol MUST be CBOR `grid` messages with the outer CBOR tag
`0x67726964` (decimal `1735551332`). Slot 0 MUST be DAG-CBOR tag 42 containing
the pCID of this spec document. The normal active shape is:

```text
grid([42(pCID), payload, proof])
```

The pCID selects this protocol's payload parser. The pCID is not a destination
address, app name, peer identity, or RPC method. The proof slot is the POC16
native proof profile used by the generic envelope builder; production successors
may choose a different pCID with COSE or transport-layer proofs.

## Promise Theory model

The promiser is the agent named in payload slot 0. The promisee is the agent or
local audience named in payload slot 1. No payload can make a promise on behalf
of another agent. Trust changes are local judgments by the receiving or observing
agent. A promise that Alice trusts Bob is Alice's promise about Alice's local
state; it is not a global fact and it does not require Bob to behave differently.

## Payload grammar

The payload is a pCID-owned CBOR array using the POC16 pair-payload profile:

```text
payload = [
  promiser: text,
  promisee: text,
  promise_about: text,
  state: [outcome: text, promise_text: text, reason: text, turn: text],
  details: [[key: text, value: text], ...]
]
```

Core slots are positional and REQUIRED. `promise_about` defaults to
`local_observation` only in compatibility builders that are given no value; new
senders SHOULD set it explicitly. `details` is sorted by key in current POC16
encoders, but receivers MUST NOT infer semantic priority from sort order. Common
`promise_about` values include `trust_update`, `peer_introduction`,
`repair_promise`, `local_event`, `route_preference`, and
`relationship_summary`; the protocol permits other relationship-local values as
long as they remain promises made by the promiser.

## Sender behavior

A sender MUST describe only its own promise, local event, local trust update, or
local non-commitment. If it reports another agent's words or bytes, it promises
only that it observed or stored those bytes from its own vantage. It MUST NOT say
that another agent has promised unless it is forwarding exact signed bytes from
that agent.

## Receiver and parser behavior

A parser MUST verify the outer grid/proof, decode the pair-payload slots, and
reject malformed arrays, non-text core fields, details entries that are not
2-slot text pairs, or trailing CBOR bytes. A receiver MAY update local trust, keep
an event, ignore the message, or answer with a local promise. Unsupported
`promise_about` values are not protocol errors; they are relationship meanings
the receiver may decline to use.

## Protocol state machine

```text
[unrelated]
    | receive promise / local event
    v
[locally known] --kept promise--> [trust increased]
      |  broken promise                  ^
      v                                  |
[trust decreased] --repair promise kept--+
      |
      | permanent distrust threshold reached
      v
[locally avoided]
```

The terminal states are local. Alice's `locally avoided` state for Mallory does
not require Bob or Carol to avoid Mallory.

## State, CAS, DAG, and retention

Agents MAY store exact relationship envelopes in their local CAS and MAY link
responses to parent message CIDs. No CAS is complete. Retention and garbage
collection are local promises by the storing agent. Relationship summaries MUST
NOT replace exact signed parent messages when exact provenance matters.

## Security considerations

Relationship payloads can affect local trust, so forged signatures, replayed
old promises, and ambiguous quotations are security-sensitive. Receivers SHOULD
correlate trust updates with exact signed parent messages where available and
SHOULD treat expired, superseded, or contextless summaries as weaker evidence.
Relationship messages MUST NOT be interpreted as permissions, authorizations, or
policy-enforcement commands.

## Interoperability notes

The pair-payload profile is intentionally array-based for constrained devices
while preserving extensibility through `details`. A production implementation can
translate `details` into local structs, maps, or database rows, but the wire
payload remains the positional CBOR array above.

## Examples

```text
grid([42(pCID),
  ["alice", "alice", "trust_update",
    ["kept", "I promise my local trust in Bob increased after he returned the promised bytes.",
     "Bob kept a storage promise", "turn-0042"],
    [["subject", "bob"], ["delta", "+3"], ["parent_exact_sha256", "abc123..."]]
  ],
  proof
])
```

A malformed example is a details entry like `["subject"]`; receivers MUST reject
it because every detail entry is a two-slot `[key, value]` pair.
