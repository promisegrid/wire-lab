# poc12-production-progress

`poc12-production-progress` is executable POC evidence for live LLM autonomy plus
deterministic production device/system agents. It keeps the POC11 sparse mesh and
adds a shipping workflow: `fulfillment` weighs a package with `postal_scale`,
gets an address from `accounting`, prints a UPS label with `ups_label_printer`,
and updates `accounting` with cost and tracking evidence. Source: `DI-timah`;
`DI-bikit`.

## What This Tests

- Multiple pCIDs in one runtime: `relationship_v1`, `postal_scale_v1`,
  `ups_label_v1`, and `accounting_v1`.
- Kernel-style pCID routing: the node parses slot 0 `42(pCID)` and dispatches to
  a local handler that promises to handle that protocol.
- Deterministic device/system agents for scale, label printer, and accounting
  alongside live LLM business/social agents.
- One top-level semantic act, `promise`; workflow steps are payload meanings,
  not RPC verbs.
- Explicit direct TCP relationship transition evidence:
  `direct_peer_added`, `direct_peer_removed`, and `direct_peer_unchanged`.

## Protocol Shape

Every message is a signed CBOR envelope:

```text
grid([42(pCID), payload, proof])
```

The pCID identifies the protocol spec. Message variants such as
`weigh_package`, `address_lookup`, `print_label`, and `shipment_update` are
payload meanings inside their protocol, not separate pCIDs. Source: `DI-bikit`.

## Shipping Agents

- `fulfillment`: hybrid workflow coordinator that executes one deterministic
  startup shipment sequence across the production pCIDs, then continues normal
  live LLM relationship turns.
- `postal_scale`: deterministic handler for `postal_scale_v1`; promises package
  weight evidence only.
- `ups_label_printer`: deterministic handler for `ups_label_v1`; promises label,
  cost, and tracking evidence only.
- `accounting`: deterministic handler for `accounting_v1`; promises address
  lookup and shipment update evidence only.

## Run

Copy the committed template and keep secrets out of it:

```sh
cp config.example.json config.json
printf '%s' "$OPENAI_API_KEY" > openai_api_key.txt
chmod 600 openai_api_key.txt
docker compose up --build --abort-on-container-exit
```

Summarize one completed run from inside the Compose volume:

```sh
docker compose run --rm --entrypoint /usr/local/bin/poc12-analyze dave /run/poc12/poc12-demo
```

`config.json`, `poc12.env`, `openai_api_key.txt`, provider outputs, and Docker
volume state are ignored and must not be committed.

## Current Limits

POC12 is provisional executable evidence, not a stable shipping, device,
accounting, kernel-routing, monitor, trust, provider, or workflow API. Docker
networking remains static; dynamic TCP relationships are local promises to dial
or accept direct exchanges. The fulfillment startup sequence is a POC guardrail
so the run produces concrete pCID-routed shipment evidence instead of relying on
a live LLM to choose that sequence unaided. Source: `DI-timah`; `DI-bikit`;
`DI-parok`.
