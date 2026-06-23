# kernel_receive_v1

## Status

Retired POC16 app/kernel registration protocol. Active POC16 uses
`kernel_transport_v1` between parser roles and transport kernels. This document
remains embedded for historical/specimen decoding. Source: `DI-bitug`;
`DI-gazin`.

## Abstract

`kernel_receive_v1` was the earlier protocol for local app receive-promise
registration. It is retained to document that app receive registration was a
promise, not a shared-volume side channel or RPC command.

## pCID and envelope

The historical shape was:

```text
grid([42(pCID), payload, proof])
```

## Promise Theory model

The local app promised that it could receive exact envelopes for a pCID. The
kernel promised only local registration or non-commitment. No remote peer was
obligated by this registration.

## Payload grammar

The payload is the pCID-owned map-body profile:

```text
payload = [
  promiser: text,
  promisee: text,
  promise_about: text,
  state: [outcome: text, promise_text: text, reason: text, turn: text],
  body: {detail_key: text => detail_value: text, ...}
]
```

All core slots are REQUIRED. `body` contains protocol-owned text key/value
details in a nested CBOR map namespace. A parser MUST reject non-arrays, wrong
array lengths, non-text core fields, non-map bodies, duplicate body keys,
reserved/core body keys, non-text body keys or values, or trailing CBOR bytes.

The main `promise_about` value was `receive_pcid`. Common details were
`app_name`, `pcid_name`, `pcid_cid`, and local endpoint information.

## Sender behavior

An app sender named the pCID it promised to receive and the local app identity.
Active POC16 parser roles now make equivalent receive promises through
`kernel_transport_v1`.

## Receiver and parser behavior

A historical kernel decoded the map-body payload and installed a local receive entry.
Active POC16 should not use this pCID for normal registration.

## Protocol state machine

```text
[app unregistered]
    | receive_pcid promise
    v
[app locally registered]
    | parser-role migration
    v
[use kernel_transport_v1 receive_pcid]
```

## State, CAS, DAG, and retention

Historical registration messages may remain in run artifacts. They should not be
used as evidence that current parser-role registration exists.

## Security considerations

A receive promise is not authorization. It only says a local app is willing to
receive exact envelopes for a pCID.

## Interoperability notes

Use this document to decode older artifacts; use `kernel_transport_v1` for active
POC16.

## Examples

```text
grid([42(pCID), ["alice-app", "kernel", "receive_pcid",
  ["kept", "I promise to receive relationship_v1 envelopes.", "app startup", "startup"],
  {"app_name": "alice-app", "pcid_name": "relationship_v1"}], proof])
```
