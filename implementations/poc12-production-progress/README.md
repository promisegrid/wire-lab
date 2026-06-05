# poc12-production-progress

`poc12-production-progress` is executable POC evidence for live LLM autonomy plus
deterministic production device/system agents. It keeps the POC11 sparse mesh and
adds a shipping workflow: `fulfillment` weighs a package with `postal_scale`,
gets an address from `accounting`, prints a UPS label with `ups_label_printer`,
and updates `accounting` with cost and tracking evidence. Source: `DI-timah`;
`DI-bikit`; `DI-galin`.

## What This Tests

- Multiple app pCIDs through one local container kernel: `relationship_v1`,
  `postal_scale_v1`, `ups_label_v1`, and `accounting_v1`.
- Real app/kernel process boundary: each container runs one `poc12-kernel`
  process plus separate local app processes for relationship, fulfillment,
  postal scale, UPS label printer, and accounting roles.
- Kernel-style pCID routing: the kernel parses slot 0 `42(pCID)`, checks local
  app receive promises, and delivers exact bytes to the app process that
  promised the target pCID.
- Deterministic device/system agents for scale, label printer, and accounting
  alongside live LLM business/social agents.
- One top-level semantic act, `promise`; workflow steps are payload meanings,
  not RPC verbs.
- Explicit direct TCP relationship transition evidence:
  `direct_peer_added`, `direct_peer_removed`, and `direct_peer_unchanged`.
- Corrected promise evidence semantics: ordinary `not_promised` /
  `non_commitment` evidence is neutral rather than broken peer-promise
  evidence, provider/runtime decision failures stay local, and duplicate
  shipment updates are checkpointed without repeated trust inflation.
  Source: `DI-jinoz`.
- App-local promise journal: POC12 now records an outstanding promise before
  applying kept/broken/malformed trust evidence, separates `send_unavailable`
  from `send_not_promised`, suppresses repeated live-agent promises, and keeps
  local budget/capacity exhaustion out of peer trust. Source: `DI-vujob`.

## Protocol Shape

Every message is a signed CBOR envelope:

```text
grid([42(pCID), payload, proof])
```

The pCID identifies the protocol spec. Message variants such as
`weigh_package`, `address_lookup`, `print_label`, and `shipment_update` are
payload meanings inside their protocol, not separate pCIDs. Source: `DI-bikit`.

## Shipping Agents

- `poc12-fulfillment`: hybrid workflow coordinator that executes one deterministic
  startup shipment sequence across the production pCIDs, then continues normal
  live LLM relationship turns. It sends shipping pCIDs but receives only the
  relationship pCID in this POC.
- `poc12-postal-scale`: deterministic app for `postal_scale_v1`; promises
  package weight evidence only.
- `poc12-ups-label-printer`: deterministic app for `ups_label_v1`; promises
  label, cost, and tracking evidence only.
- `poc12-accounting`: deterministic app for `accounting_v1`; promises address
  lookup and shipment update evidence only.
- `poc12-relationship-agent`: generic live LLM relationship app for
  `relationship_v1`.
- `poc12-kernel`: container-local transport process that records operational
  routing evidence only; it does not own trust, workflow, device behavior, or
  promise judgment. Source: `DI-galin`.

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
or accept direct exchanges between app agents, carried by local kernels. The
kernel records only app receive-promise registration, exact-byte delivery,
peer-forwarding, unregistered pCID, and transport outcomes. Apps own trust,
keep/break/non-commitment judgment, relationship ledgers, workflow state, and
deterministic device/system behavior. The fulfillment startup sequence is a POC
guardrail so the run produces concrete pCID-routed shipment evidence instead of
relying on a live LLM to choose that sequence unaided. Source: `DI-timah`;
`DI-bikit`; `DI-parok`; `DI-galin`.

Post-split evidence semantics are intentionally conservative:
`not_promised` means the receiver did not promise the requested exchange, not
that the receiver broke a promise. Transient provider/runtime decision failures
are local app/runtime evidence, not peer trust evidence. The accounting app
keeps an app-local shipment checkpoint keyed by order, tracking number, and
cost; duplicate confirmations remain visible in logs but do not repeatedly
increase trust. Source: `DI-jinoz`.

The 2026-06-05 `poc12-jinoz-20260605-055916` validation run exited cleanly and
the analyzer reported 454 events: 433 kept, 20 non-commitment, and 1 broken
resource promise. Shipping evidence appeared exactly once for address lookup,
package weight, label printing, and accounting update/confirmation. The monitor
still flagged a real remaining concern: some trust changes are driven by local
budget/capacity exhaustion and should be separated from peer trust in a later
POC. Source: `DI-jinoz`.

`DI-vujob` keeps those follow-up corrections in POC12. The runtime now records
`promise_outstanding` before resolving local promise evidence, records
`promise_resolved` before applying trust effects, treats local send failures as
`send_unavailable`, treats receiver non-commitment as `send_not_promised`,
records `local_resource_exhausted` without changing peer trust, and suppresses
repeated live-agent promise text for the same target/protocol. The startup
shipping workflow intentionally repeats the accounting update once so Docker
runs exercise duplicate shipment checkpointing through the real app/kernel path.
The analyzer reports `local_resource_counts` and
`resource_trust_coupling_counts` so future runs can show whether local
budget/capacity state is leaking back into peer-trust transitions. The
2026-06-05 `poc12-vujob-20260605-063213` Docker run exited cleanly and analyzed
to 598 events: 569 kept, 29 non-commitment, 86 `promise_outstanding`, 86
`promise_resolved`, 7 `local_resource_exhausted`, 6 `send_not_promised`, 4
`promise_repeated_suppressed`, 1 `accounting_update_duplicate`, 1
`accounting_update_duplicate_confirmed`, and empty
`resource_trust_coupling_counts`. Source: `DI-vujob`.
