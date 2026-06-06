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

ID: DI-jinoz
Date: 2026-06-05 05:54:08
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Correct POC12 post-split evidence semantics so ordinary
non-commitment, `not_promised` acknowledgements, and local provider/runtime
decision failures are not treated as broken peer promises. Add minimal
app-local accounting shipment-update checkpointing so duplicate confirmations
for the same order/tracking/cost tuple are recorded as duplicate evidence
without repeatedly increasing relationship trust.
Intent: The clean post-split Docker run showed the right app/kernel process
shape but also showed false trust drift: some non-commitments were counted as
broken promises, one transient provider error was reported as broken evidence,
and repeated accounting confirmations inflated trust for the same shipment
checkpoint. Promise Theory requires distinguishing "I did not promise that"
from "I broke a promise I made"; trust should change only when local evidence
supports a kept, broken, repaired, malformed, or discovery promise outcome.
Constraints: Do not add new top-level acts, RPC verbs, pCIDs, durable queues,
retry/backoff machinery, kernel-owned trust, service-registry authority,
protocol-level idempotency layers, trust clamping, or changes to the envelope
shape. Keep duplicate shipment handling app-local to the accounting/fulfillment
evidence loop.
Affects: `implementations/poc12-production-progress/**`;
`implementations/poc12-production-progress/README.md`;
`DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-dunoz-poc12-production-progress.md`.

ID: DI-vujob
Date: 2026-06-05 06:25:02
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement the remaining promise-trust correctness improvements
in-place in POC12 rather than deferring them to POC13. Add an app-local
outstanding promise journal inside the POC12 runtime, record local
budget/capacity exhaustion as local resource evidence rather than peer trust
evidence, classify local transport/kernel send failures separately from
receiver non-commitment and explicit peer break/malformed acknowledgements,
deduplicate repeated live-agent promise text per target/protocol, add analyzer
counts that flag trust transitions adjacent to local resource exhaustion, and
add deterministic regression coverage for duplicate shipment updates across the
same app message path used by Docker.
Intent: POC12 is the current executable evidence for kernel/app boundaries,
shipping workflow, live autonomy, and local trust. Keeping these corrections in
POC12 makes the evidence model coherent before further POCs build on it: trust
should change because a locally tracked promise was kept, broken, malformed, or
accepted as discovery evidence, not because this agent ran out of capacity,
failed to connect to its own kernel, repeated a promise, or saw a receiver
decline an exchange it never promised to accept.
Constraints: Do not add new pCIDs, top-level action kinds, RPC verbs, durable
queues, protocol-level idempotency layers, kernel-owned trust, service-registry
authority, Docker network mutation, or envelope-shape changes. The promise
journal is in-memory and app-local for POC12. Naming is locked to
`promiseJournal`, `promiseRecord`, `promiseStatus`, `rememberOutstandingPromise`,
`resolveOutstandingPromise`, `recordLocalResourceExhaustion`, and
`suppressRepeatedPromise`, with supporting status constants named
`promiseStatusOutstanding`, `promiseStatusKept`, `promiseStatusBroken`,
`promiseStatusMalformed`, `promiseStatusNonCommitment`,
`promiseStatusDuplicate`, and `promiseStatusLocalFailure`; supporting helper
names are locked to `promiseRecordKey`, `promiseStatusOutcome`,
`promiseStatusFromOutcome`, `outcomeForSendError`, and `sendEventForError`.
Path scope is limited to
`implementations/poc12-production-progress/**`, `DEV-GUIDE-RESOURCES.md`, and
this TODO file.
Affects: `implementations/poc12-production-progress/**`;
`implementations/poc12-production-progress/README.md`;
`DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-dunoz-poc12-production-progress.md`.
Supersedes: DI-jinoz for the local resource exhaustion and outstanding-promise
journal scope only.

ID: DI-pohaj
Date: 2026-06-05 18:50:50
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Extend POC12 in place with a separate `printer_port` local
kernel-role app in the existing `shipping` container. Add one new
`printer_port_v1` protocol pCID for printer-port capability-promise token issue
and redemption. The `ups_label_printer` app must ask `printer_port` for a
scoped future-print capability token, redeem that token with bounded label
bytes before returning label evidence to `fulfillment`, and treat the
printer-port outcome as local resource evidence rather than an external
authorization decision.
Intent: POC12 currently proves a message-routing kernel and separate app
processes, but it under-explores kernel roles for local hardware resources. A
printer-port app should demonstrate that a kernel can be a collection of narrow
local promise surfaces: the printer port promises bounded access to local USB or
spool-like printer hardware, while the UPS label app promises label workflow
behavior, and neither process becomes a monolithic authority.
Constraints: Keep the architecture in POC12, not POC13. Keep `printer_port` in
the existing `shipping` container. Keep one top-level semantic act,
`promise`. Do not add RPC verbs, permission language, a service registry,
kernel-owned trust, Docker network mutation, envelope-shape changes, durable
queues, or real USB dependencies. Use the approved names `printer_port`,
`poc12-printer-port`, `printer_port_v1`, `PromiseIssuePrintCapability`, and
`PromiseRedeemPrintCapability`. Path scope is limited to
`protocols/wire-lab.d/TODO/TODO-dunoz-poc12-production-progress.md`,
`implementations/poc12-production-progress/runtime/node.go`,
`implementations/poc12-production-progress/runtime/node_test.go`,
`implementations/poc12-production-progress/production/workflow.go`,
`implementations/poc12-production-progress/production/workflow_test.go`,
`implementations/poc12-production-progress/pcid/registry.go`,
`implementations/poc12-production-progress/pcid/registry_test.go`,
`implementations/poc12-production-progress/config/config.go`,
`implementations/poc12-production-progress/config.example.json`,
`implementations/poc12-production-progress/config/config_test.go`,
`implementations/poc12-production-progress/cmd/poc12-supervisor/main.go`,
`implementations/poc12-production-progress/cmd/poc12-printer-port/main.go`,
`implementations/poc12-production-progress/Dockerfile`,
`implementations/poc12-production-progress/cmd/poc12-analyze/main.go`,
`implementations/poc12-production-progress/cmd/poc12-analyze/main_test.go`,
`implementations/poc12-production-progress/README.md`, and
`DEV-GUIDE-RESOURCES.md`.
Affects: `implementations/poc12-production-progress/**`;
`DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-dunoz-poc12-production-progress.md`.

ID: DI-vutok
Date: 2026-06-05 18:52:53
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Lock the supporting names for the POC12 printer-port capability-token
implementation: `PrinterPortV1`, `PrintCapabilityScope`,
`PrintCapabilityMaxBytes`, `IssuePrintCapabilityToken`,
`ValidatePrintCapabilityToken`, `LabelBytesForShipment`,
`PrintLabelToLocalDevice`, `handlePrinterPortPromise`,
`requestPrinterPortCapability`, `redeemPrinterPortCapability`,
`field_print_capability_issuee`, `field_print_capability_token`,
`field_print_capability_token_id`, `field_print_capability_scope`,
`field_print_capability_max_bytes`, `field_label_bytes_hex`,
`field_printer_spool_id`, `printer_capability_issued`,
`printer_capability_received`, `printer_port_printed`,
`printer_port_print_confirmed`, `capabilityAck`, `printAck`, `token`,
`tokenID`, `labelBytes`, `scope`, `maxBytes`, `spoolID`, `printEvidence`,
`capabilityFields`, and `redemptionFields`.
Intent: The new printer-port kernel-role evidence needs small, explicit names
that describe promise-token issue, token validation, label bytes, local hardware
printing, and outcome observations without drifting into permission, RPC, or
authorization vocabulary.
Constraints: Use these names only inside the approved `DI-pohaj` POC12 scope.
Do not add unapproved helper names for behavior-bearing functions without
another decision lock. Keep event names local evidence labels, not protocol
commands.
Affects: `implementations/poc12-production-progress/**`;
`protocols/wire-lab.d/TODO/TODO-dunoz-poc12-production-progress.md`.

ID: DI-jupob
Date: 2026-06-05 20:02:47
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Fix the POC12 monitor lifecycle in place. Nodes must wait for
`monitor.done` using a timeout derived from the configured provider request
budget, monitor-call budget, turn delay, and shutdown grace rather than a
hard-coded 90 seconds. If the observer-only monitor fails after the protocol run
has otherwise completed, the monitor node must record `monitor_error` and still
write `monitor.done` as a non-authoritative marker so other agents do not treat
the monitor as a global authority or block forever.
Intent: The printer-port Docker run proved the shipping protocol path but
exited nonzero because early nodes timed out waiting for the monitor marker
before Dave could finish the observer report. The monitor is evidence, not a
governing participant; monitor latency or provider failure should not turn a
completed promise exchange into a failed protocol run.
Constraints: Keep the fix in POC12. Do not add pCIDs, top-level actions, RPC
verbs, durable queues, Docker network mutation, kernel-owned trust, or
envelope-shape changes. Do not commit provider logs, Docker volume state,
secrets, or local runtime config. Approved helper and event names are
`MonitorWaitTimeout`, `monitorWaitTimeout`, `writeMonitorDoneMarker`, and
`monitor_done`.
Affects: `implementations/poc12-production-progress/config/config.go`;
`implementations/poc12-production-progress/config/config_test.go`;
`implementations/poc12-production-progress/runtime/node.go`;
`implementations/poc12-production-progress/runtime/node_test.go`;
`implementations/poc12-production-progress/README.md`;
`DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-dunoz-poc12-production-progress.md`.

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
- [x] dunoz.16 Correct post-split POC12 non-commitment trust semantics,
  provider-error evidence classification, and duplicate shipment checkpoints.
- [x] dunoz.17 Implement in-place POC12 promise journal, local resource
  exhaustion separation, repeated-promise suppression, and analyzer checks.
- [x] dunoz.18 Add printer-port kernel-role capability tokens for future label
  printing and local hardware access evidence.
- [x] dunoz.19 Fix observer-only monitor lifecycle after printer-port run.

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
`GOCACHE=/tmp/wire-lab-gocache`; at that point a fresh Docker run was still
needed before the post-split process model could be treated as live-run
evidence. Later `DI-jinoz` and `DI-vujob` runs satisfied that need. Source:
`DI-timah`; `DI-bikit`; `DI-parok`; `DI-gagok`; `DI-galin`; `DI-jinoz`;
`DI-vujob`.

After the clean post-split run exposed false trust drift, `DI-jinoz` corrected
POC12 evidence semantics: non-commitment and `not_promised` acknowledgements are
neutral unless a prior explicit promise was broken, local provider/runtime
decision failures are recorded as local non-commitment evidence, and duplicate
accounting shipment checkpoints remain visible without repeatedly increasing
trust. Validation run on 2026-06-05 after `DI-jinoz`: POC12 `go test ./...`
passed with `GOCACHE=/tmp/wire-lab-gocache`; `go vet ./...`, `errcheck ./...`,
and `git diff --check` also passed. The fresh Docker run
`poc12-jinoz-20260605-055916` exited cleanly, and analyzer output reported 454
events with 433 kept, 20 non-commitment, and 1 broken resource promise.
Shipping evidence appeared once for address, weight, label, and accounting
update/confirmation. Remaining concern: the monitor still saw some trust changes
that appear tied to local budget/capacity exhaustion rather than peer evidence.
Source: `DI-jinoz`.

`DI-vujob` keeps those remaining corrections in POC12. Validation run on
2026-06-05 after `DI-vujob`: POC12 `go test ./...`, `go vet ./...`,
`errcheck ./...`, and `git diff --check` passed with
`GOCACHE=/tmp/wire-lab-gocache`. The fresh Docker run
`poc12-vujob-20260605-063213` exited cleanly, and analyzer output reported 598
events with 569 kept and 29 non-commitment outcomes. The run recorded 86
`promise_outstanding` events, 86 `promise_resolved` events, 7
`local_resource_exhausted` events, 6 `send_not_promised` events, 4
`promise_repeated_suppressed` events, and no broken outcomes. Duplicate shipment
checkpointing was exercised through the real app/kernel path with one
`accounting_update_duplicate` and one `accounting_update_duplicate_confirmed`
event. Analyzer `resource_trust_coupling_counts` was empty, which is the
intended regression check that local capacity/budget exhaustion did not
immediately leak into peer-trust transitions. Source: `DI-vujob`.

`DI-pohaj` and `DI-vutok` add the printer-port kernel-role resource owner in
POC12. Validation run on 2026-06-05 after `DI-pohaj`: POC12 `go test ./...`,
`go vet ./...`, `errcheck ./...`, and `git diff --check` passed with
`GOCACHE=/tmp/wire-lab-gocache`. The implementation adds `printer_port_v1`,
`poc12-printer-port`, deterministic print capability token issue/redemption,
wrong-token rejection coverage, pCID routing registration, analyzer shipping
event recognition, README guidance, and guide-resource notes. The first fresh
printer-port Docker run proved the shipping path but exited nonzero because
early nodes timed out waiting for the observer-only monitor marker. `DI-jupob`
replaced the hard-coded monitor wait with a config-derived budget and lets a
completed run write `monitor.done` even if the observer report fails. Validation
run on 2026-06-05 after `DI-jupob`: POC12 `go test ./...`, `go vet ./...`, and
`errcheck ./...` passed with `GOCACHE=/tmp/wire-lab-gocache`. The fresh Docker
run `poc12-jupob-20260606-030610` exited cleanly; analyzer output reported 625
events, 595 kept, 30 non-commitment, 17 `node_done`, 17
`shutdown_grace_elapsed`, 1 `monitor_done`, all printer-port shipping counts,
and empty `resource_trust_coupling_counts`. Source: `DI-pohaj`; `DI-vutok`;
`DI-jupob`.
