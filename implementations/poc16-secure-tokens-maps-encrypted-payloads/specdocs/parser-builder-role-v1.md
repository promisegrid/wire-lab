# parser_builder_role_v1

## Status

Active POC16 explanatory protocol. This document is embedded and hashed to derive
`parser_builder_role_v1`; editing it changes that pCID. Source: `DI-bitug`;
`DI-gazin`.

## Abstract

`parser_builder_role_v1` describes promises made by pCID-specific parser and
builder roles. A parser role decodes exact envelopes according to slot 0's pCID
and delivers parsed messages to local apps. A builder role takes an app's intent
and constructs exact envelopes for a pCID. This protocol documents the role split
without making pCID a destination address or RPC method.

## pCID and envelope

The active POC16 shape is:

```text
grid([42(pCID), payload, proof])
```

This pCID's payload is a relationship-style pair payload because it describes
local runtime promises and diagnostics rather than business data.

## Promise Theory model

Parser and builder roles are agents. They promise only their local behavior:
parse, build, reject, deliver, or decline. They do not command applications,
transports, or peers. An app promises what it wants to send or receive; the parser
or builder promises whether it can help with that pCID and exact byte shape.

## Payload grammar

The payload is the pCID-owned pair-payload profile:

```text
payload = [
  promiser: text,
  promisee: text,
  promise_about: text,
  state: [outcome: text, promise_text: text, reason: text, turn: text],
  details: [[key: text, value: text], ...]
]
```

All core slots are REQUIRED. `details` contains protocol-owned text key/value
pairs. A parser MUST reject non-arrays, wrong array lengths, non-text core
fields, malformed detail pairs, or trailing CBOR bytes.

Typical `promise_about` values are `parser_role_available`,
`builder_role_available`, `app_receive_promise`, `parsed_delivery`,
`builder_non_commitment`, and `malformed_payload`. Detail keys SHOULD include
`pcid_name`, `pcid_cid`, `app_name`, `exact_sha256`, `parent_exact_sha256`, and
`reason` when applicable.

## Sender behavior

A parser role announcing availability MUST name the pCID it can parse and the app
or local endpoint to which parsed messages may be delivered. A builder role MUST
name the pCID whose envelope it can construct. A role MUST NOT claim that another
role, app, kernel, or peer will process a message.

## Receiver and parser behavior

Receivers treat this protocol as local runtime coordination. Malformed payloads
are rejected. Unknown pCIDs or unsupported app receivers result in a local
non-commitment; they are not proof of peer misbehavior.

## Protocol state machine

```text
[role absent]
    | role promises pCID support
    v
[role available]
    | app sends/receives exact envelope
    v
[parse/build attempted] --valid--> [delivered or built]
        | malformed / unsupported pCID / no app receiver
        v
[local non-commitment]
```

## State, CAS, DAG, and retention

Parser/build events MAY be stored as local run events and exact envelopes MAY be
stored in a local CAS. The parser SHOULD preserve parent links from the incoming
envelope when producing ACKs.

## Security considerations

This protocol can expose local app names and supported pCIDs. Implementations
should avoid sharing it across legal-entity boundaries unless there is a local
relationship promise justifying that disclosure. Prompt-injection text in payload
strings MUST NOT alter parser validation rules.

## Interoperability notes

Production systems may implement parser/builder roles as processes, libraries,
WASM host functions, device drivers, or microcontroller firmware. The interface
must remain exact-envelope based.

## Examples

```text
grid([42(pCID),
  ["parser-a", "alice-app", "app_receive_promise",
    ["kept", "I promise to deliver relationship_v1 envelopes to Alice's local app.",
     "local app registered", "startup"],
    [["pcid_name", "relationship_v1"], ["app_name", "alice-app"]]
  ], proof
])
```
