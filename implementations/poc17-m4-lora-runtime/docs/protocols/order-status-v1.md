# order_status_v1

## Status

Active POC17 behavior-simulator protocol specification. The exact bytes of this
document are hashed as CIDv1 raw sha2-256 to derive the pCID used in slot 0.
Editing this file intentionally changes that pCID. Source: `DI-dutah`.

## Abstract

`order_status_v1` carries bintags-shaped order assignment, status update, and
ACK messages over a PromiseGrid CBOR envelope. It preserves bintags' production
vocabulary while replacing comma-separated text packets with a pCID-owned CBOR
payload.

## pCID and envelope

Messages using this protocol MUST use:

```text
grid([42(pCID), payload])
```

Slot 0 MUST be CBOR tag 42 containing the binary CID for this exact spec
document. Slot 1 MUST be a CBOR byte string containing the payload below.

## Payload grammar

The payload byte string MUST decode to this positional CBOR array:

```text
[
  type: text,          ; "MSG" or "ACK"
  source: text,
  dest: text,
  counter: uint,
  order_number: text,
  status: text
]
```

`status` values used by the POC17 scenario are `created`, `cut`, `stripped`,
`soldered`, and `completed`.

## Sender behavior

A `MSG` sender promises an order assignment or status update from its local
role. An `ACK` sender promises it observed the message with the matching
counter, order number, and status.

## Receiver behavior

A receiver MUST ignore a `MSG` whose `dest` does not match its local radio
identity. A receiver MAY store exact bytes, update local display or database
evidence, send an ACK, or decline malformed messages.

## Examples

```text
grid([42(pCID), h'...'])
payload = ["MSG", "gateway-bob", "m4-ivan", 1, "BT-1042", "created"]
payload = ["ACK", "m4-ivan", "gateway-bob", 1, "BT-1042", "created"]
```
