# TODO-sosoj: POC16 parser-role follow-up

## Status

Implemented. Owns the POC16 follow-up that corrects the gap between
`TE-ritig`'s parser/builder role conclusion and the current executable POC16
implementation. This TODO patches POC16 in place rather than creating a sibling
POC or deferring the work to POC17.

## Decision Intent Log

ID: DI-gazin
Date: 2026-06-19 12:10:30
Author: stevegt@t7a.org (Steve Traugott)
Status: active
Decision: Patch existing POC16 in place, implement real separate parser-role processes for the listener -> parser -> app flow, and consolidate active shipping/device workflow pCIDs now.
Intent: Current POC16 has app and peer TCP listeners, pCID-selected payload decoders, and parser/builder instrumentation, but the transport kernel still projects normal payload fields such as `to` and routes on them. The follow-up should make `TE-ritig` real by keeping the transport listener focused on exact envelope transport while pCID-specific parser/builder roles own pCID payload semantics, app delivery, local ACK/error promises, and backpressure promises.
Constraints: Preserve POC16 as a strict POC15 superset; preserve one top-level semantic action `promise`; preserve `grid([42(pCID), ...protocol-defined-slots])`; do not let the transport kernel decode normal app payload routing fields; do not make parser roles sign promises on behalf of apps; consolidate shipping/device operation pCIDs into a protocol-family pCID; keep old shipping/device pCID docs as historical/specimen evidence, not active runtime receive promises; keep specimen/profile pCIDs separate from active runtime pCIDs; run POC16 clean containers after implementation.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/; docs/protocols/; protocols/wire-lab.d/TODO/TODO-sosoj-poc16-parser-role-followup.md; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md; protocols/wire-lab.d/TODO/TODO.md; DEV-GUIDE-RESOURCES.md.

ID: DI-javuz
Date: 2026-06-19 18:14:20
Author: stevegt@t7a.org (Steve Traugott)
Status: active
Decision: Give POC16 child processes a longer SIGTERM grace window after introducing parser-role sessions.
Intent: The parser-role follow-up creates more independent persistent TCP sessions per container, and clean-run analysis needs those sessions to close and emit terminal records before any process is killed. The supervisor should still bound shutdown, but the bound must be long enough for kernels and parser roles to drain their local app, parser, and peer session closure records.
Constraints: Preserve fail-fast behavior after the bounded grace expires; do not weaken the analyzer's zero-unterminated-session gate; keep the change scoped to the POC16 supervisor shutdown path.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/cmd/poc16-supervisor/main.go; implementations/poc16-secure-tokens-maps-encrypted-payloads/scripts/run-clean.sh.

ID: DI-jozos
Date: 2026-06-19 19:59:56
Author: stevegt@t7a.org (Steve Traugott)
Status: active
Decision: Bound the observer-only monitor input by selecting representative POC16 events instead of sending the full clean-run event stream to the live monitor LLM.
Intent: Parser-role raw-message diagnostics substantially increase the number of runtime events, and the monitor is a POC-only observer rather than the source of analyzer truth. The collector should keep durable complete events and raw CBOR artifacts for operator review, but the live monitor prompt should receive a bounded cross-section that preserves critical events, failures, non-commitments, parser-role events, raw-message artifact markers, and shutdown/session signals without exceeding the model context window.
Constraints: Do not drop raw `events.jsonl`, `message-dag.jsonl`, or `.cbor` artifacts; do not weaken analyzer gates; do not let monitor sampling affect agent, kernel, parser, trust, ACK, or transport behavior; preserve exact run artifacts for later review.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/cmd/poc16-event-collector/main.go; implementations/poc16-secure-tokens-maps-encrypted-payloads/cmd/poc16-event-collector/main_test.go; protocols/wire-lab.d/TODO/TODO-sosoj-poc16-parser-role-followup.md.

## Current Gap

- Current POC16 implements separate local app and peer TCP listeners in the
  kernel.
- Current POC16 implements pCID-selected payload decoding through shared
  protocol helpers.
- Current POC16 records parser/builder role pressure events.
- Current POC16 does **not** yet implement a true process boundary for
  `network -> transport listener / pCID router -> pCID-specific parser role ->
  app`.
- Current POC16 still lets the transport kernel extract `to` from normal
  application payloads, which is too close to a universal routing projection.

## Implementation Plan

- [x] sosoj.1 Record this follow-up in `TODO-zugok` as the owner of the
  post-implementation parser-role correction.
- [x] sosoj.2 Update `TE-ritig` status/refinement under the TE editing policy so
  readers can see that the first POC16 implementation only partially satisfied
  the parser-role conclusion and this TODO owns the executable correction.
- [x] sosoj.3 Add a real parser-role process, with command name
  `cmd/poc16-parser-role`, package name `parserrole`, and primary types
  `ParserRole`, `ProtocolParser`, `AppReceiverRegistry`, and
  `KernelTransportClient`.
- [x] sosoj.4 Change the transport kernel so normal peer/app frames parse only
  the grid tag, slot-0 pCID, parent links, exact message hash, and proof
  validity.
- [x] sosoj.5 Permit transport-kernel payload decoding only for the
  kernel-handled control protocol `kernel_transport_v1`.
- [x] sosoj.6 Rename or supersede active use of `kernel_receive_v1` with
  `kernel_transport_v1`, covering parser-role registration and exact-envelope
  send requests.
- [x] sosoj.7 Change apps so they register receive promises with their local
  parser role, not directly with the transport kernel.
- [x] sosoj.8 Change parser roles so they register pCID receive promises with
  the local transport kernel and receive exact envelopes from it.
- [x] sosoj.9 Change parser roles so they decode pCID-owned payloads, choose
  local app receivers, forward exact envelope bytes to apps, and return app ACKs
  to the transport kernel.
- [x] sosoj.10 Change outbound apps so they submit exact signed envelopes to the
  parser role; the parser role validates pCID-owned local routing semantics and
  asks the transport kernel to carry exact bytes.
- [x] sosoj.11 Ensure parser/builder roles never sign promises on behalf of apps;
  they may promise only their own parsing, building, delivery, queueing, ACK, and
  non-commitment behavior.
- [x] sosoj.12 Add `production_shipping_v1` as the active protocol-family pCID
  for package weighing, address lookup, print capability issue/redeem, label
  printing, and shipment update promises.
- [x] sosoj.13 Remove active normal-traffic use of `postal_scale_v1`,
  `ups_label_v1`, `accounting_v1`, and `printer_port_v1`; retain those docs only
  as historical/specimen evidence.
- [x] sosoj.14 Split pCID inventory reporting into active runtime pCIDs,
  specimen/profile pCIDs, and deliberate unknown negative-path pCIDs.
- [x] sosoj.15 Target active runtime pCIDs after consolidation:
  `kernel_transport_v1`, `relationship_v1`, `route_v1`,
  `production_shipping_v1`, `cas_storage_v1`, `cid_compute_v1`,
  `identity_key_v1`, `secure_capability_v1`, `encrypted_payload_v1`, and
  `parser_builder_role_v1`.
- [x] sosoj.16 Keep message-shape specimens and map-profile examples out of
  normal receive promises and normal routing targets.
- [x] sosoj.17 Add analyzer gates for zero normal-payload decodes by the
  transport kernel, positive parser-role deliveries, positive parser-role ACKs,
  absence of old shipping/device pCIDs from normal traffic, active runtime pCID
  count of ten, and preserved POC15 superset coverage.
- [x] sosoj.18 Add unit tests for parser-role registration, exact-envelope
  delivery, ACK return, malformed payload rejection, local non-commitment when no
  app promised the pCID, receiver replacement, and transport-kernel refusal to
  decode normal app payloads.
- [x] sosoj.19 Update POC16 README, protocol docs, raw diagnostic examples, and
  `DEV-GUIDE-RESOURCES.md` to describe the corrected role boundary without
  presenting it as final PromiseGrid API.
- [x] sosoj.20 Run `go test ./...`, `errcheck ./...`, and
  `./scripts/run-clean.sh` from the POC16 implementation directory after the
  follow-up is implemented.
- [x] sosoj.21 Bound the observer-only monitor LLM input while preserving the
  complete event log, message DAG index, and raw `.cbor` artifacts for operator
  review.

## Acceptance Criteria

- The transport kernel does not decode normal application payloads to discover
  `to`, operation, promise body, or route semantics.
- Parser-role processes are actual runtime processes, not only helper functions
  or analyzer events.
- Apps communicate with parser roles for pCID-owned local semantics.
- Parser roles communicate with the transport kernel for exact envelope
  transport.
- Shipping/device workflow traffic uses `production_shipping_v1` as the active
  protocol-family pCID.
- Analyzer output distinguishes active runtime pCIDs from specimen/profile pCIDs.
- POC16 still passes its POC15 superset gates and clean container run.
