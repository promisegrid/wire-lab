# device_status_v1

## Status

Active POC17 behavior-simulator protocol specification. The exact bytes of this
document are hashed as CIDv1 raw sha2-256 to derive the pCID used in slot 0.
Editing this file intentionally changes that pCID. Source: `DI-dutah`.

## Abstract

`device_status_v1` lets a constrained device promise a small local status report
to a radio peer. The protocol is meant for POC17 behavior evidence, not as a
final embedded-device API.

## pCID and envelope

Messages using this protocol MUST use the PromiseGrid CBOR envelope:

```text
grid([42(pCID), payload])
```

Slot 0 MUST be CBOR tag 42 containing the binary CID for this exact spec
document. Slot 1 MUST be a CBOR byte string containing the payload defined
below.

## Payload grammar

The payload byte string MUST decode to this positional CBOR array:

```text
[
  device_id: text,
  status: text,
  battery_percent: uint,
  parents: [text, ...]
]
```

`parents` contains locally known or missing parent CIDs as canonical CIDv1
base32 text with the multibase `b` prefix.

## Sender behavior

The sender promises only its own local status. It does not promise global device
truth, route availability, or peer storage.

## Receiver behavior

The receiver MAY store exact bytes, update local status evidence, note missing
parents, or ignore the message. A malformed payload is a local non-commitment.

## Examples

```text
grid([42(pCID), h'...'])
payload = ["m4-ivan", "ready", 87, ["bafkreid56qdzdvph2auwjicgnsiofpwk6vib7rdpmuiayqb4foyd7ru7ly"]]
```
