# printer_port_v1

## Status

Retired POC16 historical/specimen protocol. Its embedded bytes still derive a
pCID for historical prompt context and regression evidence, but active POC16
runtime traffic SHOULD NOT use this pCID. `production_shipping_v1` now carries printer-port capability issue and redemption promises. Source: `DI-bitug`;
`DI-gazin`.

## Abstract

printer_port_v1 was an earlier single-function shipping pCID. It is retained to document
what that pCID meant and why active POC16 consolidated these operations under
`production_shipping_v1`.

## pCID and envelope

The historical shape was:

```text
grid([42(pCID), payload, proof])
```

Slot 0 named this single-function spec. Active POC16 senders should use
`production_shipping_v1` instead.

## Promise Theory model

A printer-port role promises only scoped future access to the local print resource and later local print attempts. The protocol never authorized or commanded another shipping
component.

## Payload grammar

The historical payload used the strict pCID-owned array profile:

```text
[promiser, promisee, promise_about, [outcome, promise_text, reason], body]
```

Supported body shape:

| promise_about | Body slots |
|---|---|
| `issue_print_capability` | `issuee`, `exchange_id`, `print_capability_issuee`, `print_capability_token`, `print_capability_token_id`, `print_capability_scope`, `print_capability_max_bytes` |
| `redeem_print_capability` | `print_capability_issuee`, `exchange_id`, `print_capability_token`, `print_capability_token_id`, `print_capability_scope`, `print_capability_max_bytes`, `label_bytes_hex`, `printer_spool_id` |

## Sender behavior

Historical senders were expected to set `from`, `to`, `promise_about`, state, and
all body slots. New active senders MUST prefer `production_shipping_v1` for the
same operation.

## Receiver and parser behavior

Historical receivers parsed only the listed operation. Active POC16 parser roles
SHOULD treat this pCID as specimen context, not as a normal active receive
promise, unless a local test explicitly registers it.

## Protocol state machine

```text
[historical request]
    | issue_print_capability / redeem_print_capability promise kept
    v
[historical result]
    | active POC16 migration
    v
[use production_shipping_v1]
```

## State, CAS, DAG, and retention

Historical messages may remain in local CAS or run artifacts. New messages should
parent-link through `production_shipping_v1` when part of an active shipping DAG.

## Security considerations

Do not infer current support merely because this pCID is embedded. A peer that
has not promised to receive this retired pCID has not broken a promise by
ignoring it.

## Interoperability notes

This document lets old run artifacts be decoded. It is not a recommendation to
restore pCID-per-operation shipping design.

## Examples

```text
grid([42(pCID), ["device", "fulfillment", "issue_print_capability",
  ["kept", "I promise the historical operation result.", "historical specimen"],
  ["..."]], proof])
```
