# message_shape_cose_proof_v1

## Status

Active POC16 message-shape specimen. The embedded Markdown bytes derive this
specimen pCID. Source: `DI-bitug`; `DI-mosat`.

## Abstract

message_shape_cose_proof_v1 specifies one concrete envelope arity specimen: `grid([42(pCID), payload, COSE_Sign1_detached])`. It exists so
POC16 can test that slot 0's pCID owns every later slot's meaning. It is not a
universal PromiseGrid envelope mandate.

## pCID and envelope

The exact shape is:

```text
grid([42(pCID), payload, COSE_Sign1_detached])
```

Slot 0 MUST be `42(pCID)` where the pCID is the CID of this spec document. The
outer item MUST be CBOR `grid` tag `0x67726964`. The CBOR array header gives the
arity. Receivers MUST use this pCID, not global assumptions, to interpret slots
1 and later.

## Promise Theory model

A message-shape specimen is a promise by the sender that the bytes are shaped
according to this profile. It does not require any peer to adopt the profile.
Peers may locally accept, reject, store, or ignore specimens.

## Payload grammar

The slot grammar is the envelope grammar itself:

| Slot | Meaning |
|---|---|
| 0 | DAG-CBOR tag 42 containing this protocol CID. |
| 1..n | As defined by this profile's shape above. |

When this profile names `payload`, the payload is a complete CBOR item. When it
names `parents`, the value is a CBOR array of tag-42 links to parent message CIDs.
When it names `proof`, the value is proof material defined by this profile.
When it names `COSE_Sign1`, the value MUST be CBOR tag 18 COSE_Sign1.

## Sender behavior

A sender MUST emit the exact arity in the shape, with no omitted or extra slots.
It SHOULD include parent links only where this profile says they live. It MUST
not rely on receivers applying another profile's signable view.

## Receiver and parser behavior

A generic parser may parse only the grid tag, array arity, slot-0 tag 42, and raw
CBOR slots. The profile-specific parser validates slot count and slot meanings.
Wrong arity, missing tag 42, malformed parent arrays, or invalid COSE/proof
material cause local malformed rejection.

## Protocol state machine

```text
[payload + detached COSE received] --detached signature verified--> [payload accepted]
        | wrong arity / bad proof / malformed slot
        v
[locally rejected]
```

## State, CAS, DAG, and retention

Specimen messages are exact byte artifacts. Parent links, when present, create
message-DAG edges but do not imply that any CAS has all ancestors.

## Security considerations

A profile that omits a proof slot relies on some other promise, such as transport
or session authentication. A profile that uses COSE must verify the protected
algorithm and signature before trusting payload bytes.

## Interoperability notes

These specimens document alternatives for future protocol specs. Production pCIDs
should choose one profile explicitly rather than relying on a hidden default.

## Examples

```text
grid([42(pCID), payload, COSE_Sign1_detached])
```
