# lora_link_v1

## Status

Active POC17 behavior-simulator protocol specification. The exact bytes of this
document are hashed as CIDv1 raw sha2-256 to derive the pCID used in slot 0.
Editing this file intentionally changes that pCID. Source: `DI-dutah`.

## Abstract

`lora_link_v1` carries small link-observation promises between simulated LoRa
peers. It is used to create radio-path evidence and CAS pressure in POC17.

## pCID and envelope

Messages using this protocol MUST use:

```text
grid([42(pCID), payload])
```

Slot 0 MUST be CBOR tag 42 containing the binary CID for this exact spec
document. Slot 1 MUST be a CBOR byte string.

## Payload grammar

The payload byte string is protocol-owned opaque bytes. In the current POC17
scenario those bytes are UTF-8 labels such as `link-budget-1`.

## Sender behavior

The sender promises it is reporting a local link observation or link-budget
sample. It does not promise delivery quality beyond its local vantage.

## Receiver behavior

The receiver MAY accept the link observation as local evidence, store the exact
message bytes, or ignore it.

## Examples

```text
grid([42(pCID), h'6c696e6b2d6275646765742d31'])
payload bytes = "link-budget-1"
```
