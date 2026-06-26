# peer_storage_v1

## Status

Active POC17 behavior-simulator protocol specification. The exact bytes of this
document are hashed as CIDv1 raw sha2-256 to derive the pCID used in slot 0.
Editing this file intentionally changes that pCID. Source: `DI-dutah`.

## Abstract

`peer_storage_v1` describes a constrained device's promise to ask a stronger
peer for storage help when local CAS or radio-frame limits are tight.

## pCID and envelope

Messages using this protocol SHOULD use:

```text
grid([42(pCID), payload])
```

The current POC17 simulator records the peer-storage promise as event evidence
rather than sending a full peer-storage envelope. The pCID still names the
intended protocol surface.

## Payload grammar

When carried as a message, the payload byte string SHOULD decode to:

```text
[
  requester: text,
  peer: text,
  content_cid: text,
  reason: text
]
```

## Sender behavior

The sender promises it is asking for storage help. It does not promise the peer
will accept, retain, or serve the content.

## Receiver behavior

The receiver MAY accept, refuse, defer, or ignore the storage request according
to local policy.

## Examples

```text
payload = ["m4-ivan", "gateway-bob", "bafkreid56qdzdvph2auwjicgnsiofpwk6vib7rdpmuiayqb4foyd7ru7ly", "local CAS pressure"]
```
