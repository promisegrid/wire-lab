# map_payload_profile_v1

## Status

Active POC16 specimen/profile protocol. The embedded Markdown bytes derive this
pCID. Source: `DI-bitug`; `DI-vulit`.

## Abstract

`map_payload_profile_v1` demonstrates that a pCID may choose a CBOR map payload
for self-documenting fields. Maps are permitted but not the default for all
PromiseGrid protocols; constrained protocols should prefer small positional
arrays.

## pCID and envelope

The specimen shape is:

```text
grid([42(pCID), payload, proof])
```

Slot 1 is a CBOR map whose key and value semantics are defined by this pCID.

## Promise Theory model

A map payload is still a promise payload. Field names do not create authority,
permissions, or universal semantics. The promiser promises only what the map says
it promises.

## Payload grammar

```text
payload = {
  "act": "promise",
  "from": text,
  "to": text,
  "promise_about": text,
  "outcome": text,
  "promise": text,
  "reason": text,
  ... pCID-defined text keys ...
}
```

POC16's generic map helper supports text-string maps and, for token internals,
limited CBOR booleans and signed integers under specific helpers. This profile is
for application-facing string maps.

## Sender behavior

A sender SHOULD use maps only when field self-documentation outweighs byte count
and constrained-device simplicity. It MUST still sign the exact envelope and set
clear promiser/promisee fields if local parser routing depends on them.

## Receiver and parser behavior

A receiver MUST parse the map under this pCID only. It MUST NOT assume every
other pCID uses maps. Unknown keys may be ignored or locally recorded; missing
required keys produce local non-commitment or malformed rejection depending on
receiver policy.

## Protocol state machine

```text
[map profile selected]
    | map has required keys
    v
[payload interpreted] --promise accepted locally--> [local event kept]
        | missing required key / wrong type
        v
[malformed or not promised]
```

## State, CAS, DAG, and retention

Map payloads are exact bytes like any other envelope payload and may be stored in
CAS. Receivers SHOULD store raw bytes if later exact interpretation matters.

## Security considerations

Maps are easier for humans and LLMs to read but easier to overgeneralize. A map
key named `permission` or `policy` has no authority unless the pCID explicitly
frames it as a voluntary local promise.

## Interoperability notes

This profile is useful for rich apps and LLM prompt context. IoT and 100-year
small-device protocols should prefer arrays.

## Examples

```text
grid([42(pCID),
  {"act":"promise", "from":"alice", "to":"bob", "promise_about":"demo",
   "outcome":"kept", "promise":"I promise this map is a profile specimen."},
  proof
])
```
