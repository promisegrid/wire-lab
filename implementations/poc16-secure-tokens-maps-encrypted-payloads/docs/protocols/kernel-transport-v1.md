# kernel_transport_v1

## Status

Active POC16 local kernel-role protocol. This document is embedded into the
POC16 binary and hashed to derive the pCID for `kernel_transport_v1`; editing it
changes that pCID. Source: `DI-bitug`; `DI-gazin`.

## Abstract

`kernel_transport_v1` is the local protocol between a container's pCID-specific
parser role and that container's transport kernel. It is not an application
protocol, not a global registry, and not a permission surface. It lets a parser
promise that it can receive exact envelopes for a pCID and ask the local
transport kernel to carry exact envelope bytes toward a target that the parser
has already found in the pCID-owned application payload.

## pCID and envelope

The wire shape is:

```text
grid([42(pCID), payload, proof])
```

The payload belongs to this protocol, so the transport kernel MAY decode it. If
the payload embeds an application envelope, the transport kernel MUST treat that
embedded envelope as exact bytes except for generic envelope checks: grid tag,
slot-0 pCID, parent links, exact message CID, and proof validity.

## Promise Theory model

The parser role is the usual promiser. The local transport kernel is the usual
promisee. The parser promises local parser/kernel behavior; it does not promise
that a remote agent will accept, process, or keep any embedded application
promise. The kernel promises only local carriage or local non-commitment.

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

`promise_about` values defined by this spec are:

| promise_about | Required detail keys | Meaning |
|---|---|---|
| `receive_pcid` | `app_name`, `pcid_name`, `pcid_cid`, `transport_action` | Parser promises it can receive exact envelopes for the named application pCID and forward them only to matching local app receivers. |
| `carry_exact_envelope` | `target`, `target_protocol`, `target_exact_cid`, `envelope_b64`, `transport_action` | Parser asks the kernel to carry the exact base64 envelope bytes toward the named target. |

`transport_action` MUST match `promise_about` for these two actions. Empty
`target` or `envelope_b64` values make `carry_exact_envelope` not promised.

## Sender behavior

For `receive_pcid`, the parser MUST name the app and pCID it promises to handle.
For `carry_exact_envelope`, the parser MUST provide the target, embedded protocol
name, exact SHA-256 of the embedded envelope bytes, and base64 of the exact
embedded envelope bytes. The parser is responsible for application payload
parsing; the kernel is not asked to infer destinations from arbitrary app
payloads.

## Receiver and parser behavior

The kernel MUST reject malformed control payloads, undecodable base64, embedded
envelopes whose exact CID does not match `target_exact_cid`, unknown targets,
and missing peer endpoints. The kernel MAY respond with a not-promised ACK linked
to the parent message. The kernel MUST NOT decode embedded application payloads
to discover application semantics.

## Protocol state machine

```text
[parser connected]
    | receive_pcid
    v
[pCID receiver promised]
    | carry_exact_envelope + valid bytes + known target
    v
[carriage attempted] --peer/local delivery kept--> [ACK returned]
        | malformed / unknown / no route
        v
[not-promised ACK returned]
```

## State, CAS, DAG, and retention

The kernel MAY store exact control and embedded envelopes in local run-scoped
message artifacts. Parent links are carried in the embedded envelope and are not
rewritten by this protocol. Persistent sessions are local transport state, not
relationship trust.

## Security considerations

The base64 envelope is untrusted until parsed, verified, and hash-checked. A
parser that lies about `target_exact_cid` receives local non-commitment. The
kernel preserves voluntary local-promise routing by delegating application
payload addressing to the parser role.

## Interoperability notes

This protocol is local to a node/container boundary in POC16. Other runtimes may
replace it with function calls, device-driver queues, stdio, WASM host calls, or
another pCID, but the same promise split should hold: parser roles understand
app protocols; transport kernels carry exact envelopes.

## Examples

```text
grid([42(pCID),
  ["parser-a", "kernel-a", "carry_exact_envelope",
    ["kept", "I promise these exact bytes should be carried toward Bob.",
     "parser found Bob in production_shipping_v1 payload", "turn-010"],
    {"target": "bob", "target_protocol": "production_shipping_v1",
     "target_exact_cid": "9a...", "envelope_b64": "2GRncmlk...",
     "transport_action": "carry_exact_envelope"}
  ], proof
])
```
