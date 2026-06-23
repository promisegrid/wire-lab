# unknown_probe_v1

## Status

Active negative-test protocol specification for POC17. The exact bytes of this
document are hashed as CIDv1 raw sha2-256 to derive the pCID used in slot 0.
Editing this file intentionally changes that pCID. Source: `DI-dutah`.

## Abstract

`unknown_probe_v1` exists so POC17 can prove that an agent may preserve exact
bytes while making no semantic commitment for a pCID it does not support.

## pCID and envelope

Messages using this protocol MUST use:

```text
grid([42(pCID), payload])
```

Slot 0 MUST be CBOR tag 42 containing the binary CID for this exact spec
document. Slot 1 MUST be a CBOR byte string.

## Payload grammar

The payload byte string is opaque. The current simulator uses the bytes
`probe`.

## Sender behavior

The sender promises only that it is sending a probe under this pCID.

## Receiver behavior

Receivers that do not support this pCID SHOULD record a local non-commitment
instead of accepting the payload meaning.

## Examples

```text
grid([42(pCID), h'70726f6265'])
payload bytes = "probe"
```
