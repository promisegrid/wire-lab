# peer_storage

## Status

Active POC17 behavior-simulator protocol specification. The exact bytes of this
document are hashed as CIDv1 raw sha2-256 to derive the pCID used in slot 0.
Editing this file intentionally changes that pCID. Source: `DI-gidul`.

## Abstract

`peer_storage` describes compact peer-storage promises between a constrained
device and a stronger peer. The stronger peer first issues a scoped bearer
capability token. The constrained device later presents that token when asking
the peer to put or get exact bytes by CID.

## pCID and envelope

Messages using this protocol use:

```text
grid([42(pCID), payload])
```

The pCID defines all payload slots. The payload does not carry a redundant
generic `promise` action or a redundant kind slot; receivers distinguish the
promise shape from slot count and slot types.

## Payload grammar

Bob grants Ivan a scoped storage capability:

```text
[
  issuer: text,
  holder: text,
  token: bytes,
  allowed_kinds: [text],
  max_content_bytes: uint,
  max_retained_objects: uint,
  retention_terms: text
]
```

Ivan asks Bob to retain exact bytes:

```text
[
  holder: text,
  issuer: text,
  token: bytes,
  42(content_cid),
  content: bytes,
  reason: text
]
```

Bob reports whether he accepted a put promise:

```text
[
  issuer: text,
  holder: text,
  42(token_cid),
  42(content_cid),
  accepted: uint,
  note: text
]
```

Ivan asks Bob to return exact bytes for a CID:

```text
[
  holder: text,
  issuer: text,
  token: bytes,
  42(content_cid),
  reason: text
]
```

Bob fulfills a get promise only by returning exact bytes:

```text
[
  issuer: text,
  holder: text,
  42(token_cid),
  42(content_cid),
  content: bytes
]
```

Bob refuses a get promise:

```text
[
  issuer: text,
  holder: text,
  42(token_cid),
  42(content_cid),
  reason: text
]
```

## Sender behavior

A `grant` message is Bob's promise that Ivan may later present the token for the
listed put/get kinds under the size, object-count, and retention terms. A `put`
message is Ivan's promise that the presented bytes match the CID and that Ivan is
redeeming Bob's token under its terms. A `get` message is Ivan's promise that it
is asking for bytes under that same token.

## Receiver behavior

Bob verifies token bytes, holder, allowed kind, content size, object count, and
content CID before accepting a put or fulfilling a get. Ivan verifies every
returned content byte string by recomputing its CID before retaining it locally.

## Security considerations

The POC17 token is a compact bearer capability for the 200-byte simulated LoRa
path. It is not a production CWT/COSE token. Full radio-visible CWT/COSE token
profiles require a later constrained-token decision. Source: `DI-gidul`.

## Examples

```text
put = ["m4-ivan", "gateway-bob", h'...', 42(content_cid), h'...', "cas_retention_limit"]
get = ["m4-ivan", "gateway-bob", h'...', 42(content_cid), "missing_parent_after_restart"]
```
