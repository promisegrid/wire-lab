# production_shipping_v1

## Status

Active POC16 shipping/device workflow protocol. This spec consolidates earlier
single-function pCIDs and its exact embedded bytes derive the active pCID.
Source: `DI-bitug`; `DI-gazin`.

## Abstract

`production_shipping_v1` models a multi-agent shipping workflow: a fulfillment
agent coordinates package weighing, address lookup, UPS label creation,
printer-port capability-token issue/redeem, physical label printing, and
accounting shipment updates. The protocol is one family pCID; operation choice is
inside the pCID-owned payload, not in slot 0.

## pCID and envelope

The active shape is:

```text
grid([42(pCID), payload, proof])
```

The pCID identifies this shipping protocol family. Local app names, peer names,
orders, packages, and operation kinds are payload semantics.

## Promise Theory model

Each device or business app promises only what it controls. A postal scale may
promise a measured weight; an accounting app may promise an address or update; a
label app may promise a label and cost; a printer-port role may promise future
print attempts through capability tokens. No component commands another.

## Payload grammar

The payload is a strict pCID-owned CBOR array:

```text
payload = [
  promiser: text,
  promisee: text,
  promise_about: text,
  state: [outcome: text, promise_text: text, reason: text],
  body: [text, ...]
]
```

`promise_about` selects the body shape:

| promise_about | Body slots |
|---|---|
| `weigh_package` | `package_id`, `exchange_id`, `weight_ounces` |
| `address_lookup` | `order_id`, `exchange_id`, `shipping_address` |
| `print_label` | `package_id`, `exchange_id`, `shipping_address`, `weight_ounces`, `tracking_number`, `cost_cents`, `printer_spool_id` |
| `shipment_update` | `order_id`, `exchange_id`, `tracking_number`, `cost_cents`, `duplicate_shipment_update` |
| `issue_print_capability` | `issuee`, `exchange_id`, `print_capability_issuee`, `print_capability_token`, `print_capability_token_id`, `print_capability_scope`, `print_capability_max_bytes` |
| `redeem_print_capability` | `print_capability_issuee`, `exchange_id`, `print_capability_token`, `print_capability_token_id`, `print_capability_scope`, `print_capability_max_bytes`, `label_bytes_hex`, `printer_spool_id` |

All body slots are text in POC16. Empty text means absent or not promised for that
slot unless the specific promise requires a value.

## Sender behavior

A sender MUST choose exactly one listed `promise_about` value. It MUST set `from`
and `to` to the promiser and promisee used for local parser routing. It SHOULD
reuse `exchange_id` across all promises for one shipment. It MUST not use a
shipping promise to imply authorization; capability tokens are issuer promises
that may be redeemed according to their own terms.

## Receiver and parser behavior

A parser MUST reject unsupported `promise_about` values, wrong body length,
wrong state length, non-text slots, and trailing CBOR. A receiver MAY answer with
an ACK payload using the same `promise_about` and parent-linking to the request.
Duplicate shipment updates SHOULD be represented by `duplicate_shipment_update`,
not by pretending the original update was kept twice.

## Protocol state machine

```text
[start shipment]
    | weigh_package kept
    v
[weight known] --address_lookup kept--> [address known]
    |                                      |
    | issue_print_capability kept          v
    +-------------------------------> [print capability held]
                                           |
                                   print_label kept
                                           v
                                  [label ready]
                                           |
                                 redeem_print_capability kept
                                           v
                                  [label printed]
                                           |
                                  shipment_update kept
                                           v
                                  [shipment recorded]
```

Any step may produce local non-commitment or malformed rejection instead of
advancing. Repair is a new promise, not a retroactive change to the broken one.

## State, CAS, DAG, and retention

Agents SHOULD parent-link ACKs and downstream promises to the request or prior
workflow message. Label bytes MAY be stored in a local CAS under their content
CID. Business systems decide local retention and GC promises for address, cost,
and tracking data.

## Security considerations

Shipping addresses and label bytes are sensitive. Agents SHOULD send them only to
trusted peers or encrypt them under another pCID. Printer capability tokens MUST
be scoped, bounded, and single-purpose enough that redeeming a token cannot print
arbitrary future labels outside the issuer's promise.

## Interoperability notes

Earlier POC16 pCIDs `postal_scale_v1`, `ups_label_v1`, `accounting_v1`, and
`printer_port_v1` are retired specimens. Active POC16 traffic uses this protocol
family to avoid pCID-per-operation fragmentation.

## Examples

```text
grid([42(pCID),
  ["postal-scale", "fulfillment", "weigh_package",
    ["kept", "I promise package pkg-7 weighed 42.5 ounces on my local scale.",
     "scale stable"],
    ["pkg-7", "ship-123", "42.5"]
  ], proof
])
```
