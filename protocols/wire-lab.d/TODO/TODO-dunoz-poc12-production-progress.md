# TODO-dunoz: POC12 production progress

## Decision Intent Log

ID: DI-timah
Date: 2026-06-04 06:13:01
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement `implementations/poc12-production-progress/` as a
POC11-derived live-LLM-first proof of concept that keeps the generic sparse-mesh
autonomy pressure and adds deterministic production agents for a postal scale,
UPS label printer, and accounting system. The workflow agent promises to weigh
packages, get shipping addresses, print UPS labels, and update accounting with
shipping cost and tracking number. POC12 keeps one pCID-selected signed CBOR
`grid([42(pCID), payload, proof])` protocol over length-framed TCP and keeps one
future-facing top-level semantic act: `promise`.
Intent: POC11 proved adaptive local TCP relationships and live autonomy are
credible but generic. POC12 should move toward production reality by making a
business workflow cross deterministic device/system agents while preserving
Promise Theory: every participant promises only for itself, all trust remains
local, no receiver is commanded, and workflow steps are pCID-owned promise
payload meanings rather than RPC verbs.
Constraints: Treat POC12 as executable evidence, not a stable shipping,
device, accounting, monitor, trust, provider, kernel, or workflow API. The
postal scale, UPS label printer, and accounting system are deterministic agents;
live LLM autonomy remains with the fulfillment and generic agents. Dynamic TCP
relationships mean local dial/accept promises, not Docker network mutation.
Do not commit secrets, provider outputs, Docker volume state, or local runtime
config.
Affects: `implementations/poc12-production-progress/**`;
`implementations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-dunoz-poc12-production-progress.md`;
`protocols/wire-lab.d/TODO/TODO.md`.

ID: DI-bikit
Date: 2026-06-04 15:58:44
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC12 must test pCID-based kernel routing by supporting multiple
protocol pCIDs in one node runtime: `relationship_v1`, `postal_scale_v1`,
`ups_label_v1`, and `accounting_v1`. The runtime parses slot 0
`42(pCID)`, selects a local handler that has promised to handle that protocol,
and records `unsupported_pcid` when no local handler exists. Message variants
inside each protocol remain payload meanings, not separate pCIDs per message
type.
Intent: A single POC-wide pCID only tests pCID checking. POC12 should test the
kernel behavior we care about: route by protocol identity to local app/device
handlers without creating a central service registry or RPC command surface.
Constraints: Keep one top-level semantic act (`promise`). Keep Docker networking
static. Dynamic TCP relationships remain local dial/accept promises. Do not make
the pCID a message-type selector; each pCID names a protocol spec, and the
payload carries protocol-owned meanings such as weigh package, print label, or
update accounting.
Affects: `implementations/poc12-production-progress/**`;
`protocols/wire-lab.d/TODO/TODO-dunoz-poc12-production-progress.md`;
`implementations/README.md`; `DEV-GUIDE-RESOURCES.md`.
Supersedes: DI-timah for the single-pCID runtime shape only.

ID: DI-parok
Date: 2026-06-04 16:16:35
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC12 fulfillment is a hybrid production-workflow agent: it executes
one concrete deterministic shipment sequence against `accounting`,
`postal_scale`, `ups_label_printer`, and `accounting` again before normal live
LLM relationship turns continue.
Intent: The first live Docker run proved that a prompt-only fulfillment agent
can talk about future shipment work without producing pCID-routed shipment
evidence. POC12 needs executable evidence that the kernel routes multiple pCIDs
through a realistic production workflow while still leaving later relationship
behavior to live agents.
Constraints: Keep the single top-level act as `promise`; keep workflow steps as
pCID-owned payload meanings; do not add RPC verbs or central routing authority;
do not make the device/system handlers live LLM actors; do not commit secrets,
provider output, Docker volume state, or local runtime config.
Affects: `implementations/poc12-production-progress/runtime/node.go`;
`implementations/poc12-production-progress/runtime/node_test.go`;
`implementations/poc12-production-progress/config.example.json`;
`implementations/poc12-production-progress/README.md`;
`protocols/wire-lab.d/TODO/TODO-dunoz-poc12-production-progress.md`.
Supersedes: DI-timah for prompt-only fulfillment workflow execution only.

ID: DI-gagok
Date: 2026-06-05 04:47:47
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Record the latest POC12 assessment in `DEV-GUIDE-RESOURCES.md`,
including what happened in the transactions, how the shipping workflow behaved,
where Promise Theory fit is strong or weak, whether impositions remain, how
local incentives operate, whether autonomy and dynamic TCP relationships are
working, and which next POCs should prepare the project for real production
use.
Intent: POC12 is now executable evidence rather than only a scaffold. The guide
needs a concrete run interpretation so future design work does not over-read
the POC as a stable API and does not miss the remaining production-readiness
gaps exposed by the latest run.
Constraints: Treat the latest run as evidence, not specification. Preserve
Promise Theory framing: local promises, local trust, no global authority, no
receiver command semantics, and pCID-owned payload meanings. Do not commit
provider logs, Docker volume state, secrets, or local runtime config.
Affects: `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-dunoz-poc12-production-progress.md`.

ID: DI-galin
Date: 2026-06-05 05:10:59
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Refactor POC12 so every app is a separate local process with its own
command entrypoint, every container runs one separate `poc12-kernel` process,
and apps communicate with that local kernel over loopback TCP using the same
length-framed `grid([42(pCID), ...])` envelope format used between kernels.
The kernel owns only transport, pCID parsing, local receive-promise
registration, exact-byte routing to app processes, peer forwarding, and
operational delivery evidence. Apps own trust, promise interpretation,
keep/break/non-commitment judgments, relationship ledgers, deterministic device
behavior, shipping workflow behavior, and provider-backed autonomy.
Intent: POC12 had drifted back toward a combined agent/kernel process and risked
making the kernel look like a trust ledger, service registry, RPC authority, or
business workflow engine. PromiseGrid apps are local processes that promise to
handle pCIDs; the kernel should only carry exact messages between local app
processes and peer kernels without pretending to command agents or judge trust.
Constraints: Do not add a generic `poc12-app` entrypoint. Do not add new
top-level acts, RPC verbs, idempotency layers, durable queues, retry/backoff
machinery, trust clamping, service-registry authority, or changes to pCID
semantics or envelope shape. Keep a single semantic act, `promise`, and keep
the local receive-promise protocol as promise content under one pCID.
Affects: `implementations/poc12-production-progress/**`;
`implementations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-dunoz-poc12-production-progress.md`.
Supersedes: DI-timah for the app/kernel process boundary; DI-bikit for where
pCID routing responsibility lives.

## Prior aliases

- None.

## Tasks

- [x] dunoz.1 Record POC12 implementation DI and TODO owner.
- [x] dunoz.2 Create the standalone `implementations/poc12-production-progress/`
  module from the POC11 scaffold.
- [x] dunoz.3 Rename commands, module path, docs, config, and compose surfaces
  from POC11 adaptive trust TCP to POC12 production progress.
- [x] dunoz.4 Add multi-pCID kernel routing for `relationship_v1`,
  `postal_scale_v1`, `ups_label_v1`, and `accounting_v1`.
- [x] dunoz.5 Add deterministic production agents for `postal_scale`,
  `ups_label_printer`, and `accounting`.
- [x] dunoz.6 Add the live fulfillment workflow agent and configure it to use
  the production agents while preserving local trust and receive promises.
- [x] dunoz.7 Add explicit direct-peer transition events for local TCP
  relationship additions, removals, and unchanged states.
- [x] dunoz.8 Add shipping workflow outcome checks for package weight, shipping
  address lookup, label/tracking/cost generation, and accounting update.
- [x] dunoz.9 Update analyzer summaries for shipping outcomes and relationship
  transitions.
- [x] dunoz.10 Add deterministic tests for production agents, workflow checks,
  receive-promise gating, relationship transitions, and analyzer output.
- [x] dunoz.11 Update implementation docs and development-guide resource notes
  with POC12 as provisional evidence only.
- [x] dunoz.12 Run Go validation and record implementation status.
- [x] dunoz.13 Make the fulfillment workflow produce concrete pCID-routed
  shipping evidence during live Docker validation.
- [x] dunoz.14 Record the latest POC12 assessment in guide resources.
- [x] dunoz.15 Split POC12 into one kernel process per container and separate
  local app processes with app-owned trust and promise judgment.

## Implementation status

Implemented in `implementations/poc12-production-progress/` with four local
pCID handlers, deterministic shipping device/system handlers, explicit direct
peer transition evidence, analyzer shipping/transition summaries, and unit tests.
Validation run on 2026-06-04 after `DI-parok`: `go test ./...`, `go vet
./...`, and `errcheck ./...` passed with `GOCACHE=/tmp/wire-lab-gocache`. The
corrected live Docker run exited cleanly and analyzer output reported 330 total
events with non-empty `shipping_counts`: address promised/received, package
weighed/received, label printed/received, and accounting updated/confirmed.
The guide assessment now records the transaction details, Promise Theory fit,
remaining imposition/alignment risks, incentives, autonomy status, dynamic TCP
relationship evidence, shipping-message details, and next POC pressure points.
The implementation now also has a separate `poc12-kernel` process per container
and separate local app command entrypoints. Apps register receive promises with
the local kernel over loopback TCP; the kernel routes exact pCID-selected
envelopes and records operational delivery evidence only. Apps retain trust,
workflow, device behavior, and promise judgment. Validation run on 2026-06-05
after `DI-galin`: `go test ./...` passed with
`GOCACHE=/tmp/wire-lab-gocache`; a fresh Docker run is still needed before the
post-split process model is treated as live-run evidence. Source: `DI-timah`;
`DI-bikit`; `DI-parok`; `DI-gagok`; `DI-galin`.
