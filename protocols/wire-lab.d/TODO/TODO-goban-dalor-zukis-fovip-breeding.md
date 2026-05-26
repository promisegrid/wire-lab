# TODO-goban: Dalor, Zukis, and Fovip breeding run

## Prior aliases

This TODO was minted directly under the post-`TE-mumuv` proquint handle scheme.
No integer or timestamp alias exists.

## Status

Configured pending commit and Steve-run execution. The configured GA breeding
slice uses `SIM-dalor-grid-envelope-protocol-owned-signature-slot` as the anchor
parent and tests it against both
`SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig` and
`SIM-fovip-kernel-promise-boundary-port-contract`.

## Decision Intent Log

### DI-kipoz

ID: DI-kipoz
Date: 2026-05-25 15:11:16
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Plan a dalor-anchored breeding run against `SIM-zukis` and `SIM-fovip`
before promoting any dalor-derived child.
Intent: `SIM-dalor` captures the compact protocol-owned proof-slot idea;
`SIM-zukis` captures the current `grid([42(pCID), payload, ...])` envelope
direction; `SIM-fovip` captures the kernel/app promise boundary pressure. The
next useful children should test whether the envelope shape and kernel boundary
compose cleanly without reintroducing external authority, universal
conformance, or non-promise command semantics.
Constraints: Steve runs canary/GA commands. Do not overwrite existing score
evidence. Keep provider-backed runs additive. Do not promote generated children
merely because they exist under `proposals/`. Do not propagate the legacy
`port promise record` wording from `SIM-fovip`; generated children should use
`kernel implementation promises` and Promise Theory vocabulary.
Affects: `simulations/SIM-dalor-grid-envelope-protocol-owned-signature-slot/`;
`simulations/SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig/`;
`simulations/SIM-fovip-kernel-promise-boundary-port-contract/`;
future `proposals/<run-group-id>/`; future `results/state/<run-group-id>.json`;
future `results/jobs/<run-group-id>/`.

### DI-dikat

ID: DI-dikat
Date: 2026-05-25 16:45:41
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Reject `SIM-lasuv-child-kernel-port-capability-outcome-ledger` from
the goban breeding proposals, mine only its local observation/evidence lesson,
and create `SIM-hozif-app-pcid-promises-local-observations` as a simpler
Vinag-derived successor rather than promoting `SIM-vinag-child-tag42-port-
observation-pairing` as scored.
Intent: The goban child review showed that Lasuv's kernel-port ledger,
per-`pCID` capability table, service-promiser vocabulary, and outcome-ledger
shape drift back toward RPC/service-registry/capability-table design. Vinag has
better envelope heritage and local-trust boundaries, but its service vocabulary
and supported-pCID lists should be reframed before canonical use. The durable
successor should say that apps promise local kernels which pCIDs they will
receive or handle, kernels promise best-effort delivery/observation behavior,
and all outcomes remain local observations rather than authority facts.
Constraints: Do not rewrite scored child bytes under `proposals/`. Use
`tools/ga-runner cull` for rejected Lasuv so the state records the destructive
review decision. Leave Vinag's scored proposal bytes unpromoted and unchanged;
the successor is a new human-authored simulation, not a byte-identical
promotion. Provider-backed recalibration is configured for Steve to run later,
not run by Codex.
Affects: `results/state/ga-canary-20260525-goban-dalor-fovip-v4.json`;
`proposals/ga-canary-20260525-goban-dalor-fovip-v4/`;
`simulations/SIM-hozif-app-pcid-promises-local-observations/`;
`protocols/wire-lab.d/TODO/TODO-goban-dalor-zukis-fovip-breeding.md`;
future `/tmp/canary-cells` or `/tmp/run-*` calibration commands.

### DI-bozid

ID: DI-bozid
Date: 2026-05-25 17:19:22
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Revise Hozif so promise and observation records are two payload
record kinds under one stable protocol pCID, not two separate pCID-selected
payload protocols.
Intent: A pCID names a protocol specification and should usually be stable
across many messages, closer to a network protocol version than to a per-message
type. Hozif's promise and observation records are tightly coupled halves of one
local app/kernel promise-accounting protocol; splitting them into separate pCIDs
would imply unnecessary churn and independent deployment surfaces.
Constraints: Preserve Hozif's anti-RPC, pro-promise correction from `DI-dikat`.
Keep the envelope shape `grid([42(pCID), payload, ...])`; only the Hozif
payload protocol structure changes. Do not run provider-backed recalibration
from Codex.
Affects: `simulations/SIM-hozif-app-pcid-promises-local-observations/`;
`protocols/wire-lab.d/TODO/TODO-goban-dalor-zukis-fovip-breeding.md`;
`/tmp/run-hozif-anti-rpc-calibration.sh`.

### DI-ditan

ID: DI-ditan
Date: 2026-05-25 17:34:31
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Cull `SIM-vinag-child-tag42-port-observation-pairing` from the goban
proposal tree instead of promoting it.
Intent: Vinag's useful envelope and audit pressure has been carried forward
into Hozif, while Vinag itself still uses two tightly coupled pCID-selected
payload protocols plus service/dispatch/support vocabulary that conflicts with
the anti-RPC and pCID-stability refinements in `DI-sitim` and `DI-gakij`.
Keeping Vinag under `proposals/` would invite stale promotion review based on
pre-refinement scores.
Constraints: Use `tools/ga-runner cull` so the rejection is recorded in the
run state. Do not rewrite Vinag scored child bytes. Keep Hozif as the successor
simulation path.
Affects: `results/state/ga-canary-20260525-goban-dalor-zukis-v4.json`;
`proposals/ga-canary-20260525-goban-dalor-zukis-v4/`;
`protocols/wire-lab.d/TODO/TODO-goban-dalor-zukis-fovip-breeding.md`.

### DI-pavub

ID: DI-pavub
Date: 2026-05-25 18:15:22
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Cull `SIM-fozid-child-slim-port-commitment-service-outcomes` as a
proposal and create `SIM-kifas-app-kernel-surface-promises-local-observations`
as a Hozif-derived successor instead of promoting Fozid.
Intent: Fozid usefully exposed slimmer per-surface promiser entries and
kept/refused/unavailable/broken/timed-out distinctions, but it also regressed
toward service/request/dispatch/authorization vocabulary and an optional
`grid([pCID, payload, signature])` carriage shape. The successor should keep
Hozif's one-stable-pCID, `grid([42(pCID), payload, ...])`, local-observation,
and anti-RPC structure while mining only Fozid's useful surface-level promise
and outcome distinctions.
Constraints: Use `tools/ga-runner cull` so the rejected proposal is recorded in
the run state. Do not promote or rewrite Fozid proposal bytes. Keep the
successor human-authored under canonical `simulations/` and avoid service
registry, capability table, dispatcher, authorization, permission, conformance,
or external-authority semantics. Provider-backed scoring remains Steve-run.
Affects: `results/state/ga-canary-20260525-hozif-anti-rpc-calibration.json`;
`proposals/ga-canary-20260525-hozif-anti-rpc-calibration/`;
`simulations/SIM-kifas-app-kernel-surface-promises-local-observations/`;
`protocols/wire-lab.d/TODO/TODO-goban-dalor-zukis-fovip-breeding.md`;
future `/tmp/run-*` calibration commands.

## Breeding intent

The run should test three questions:

1. Can `SIM-dalor`'s compact mandatory sender-proof slot merge with
   `SIM-zukis`'s `42(pCID)` slot-0 ecosystem-compatible selector without
   reopening the settled `grid([42(pCID), payload, ...])` direction?
2. Can `SIM-dalor`'s sender promise evidence be used at the app/kernel boundary
   described by `SIM-fovip`, where exposed app/kernel operations are
   pCID-selected grid messages rather than RPC commands?
3. Can the resulting children keep trust local and promise-first across wire
   envelope parsing, kernel implementation promises, storage, computation,
   send/receive, refusal, unsupported-pCID behavior, and evidence recording?

## Required parent set

- Anchor parent:
  `SIM-dalor-grid-envelope-protocol-owned-signature-slot`
- Envelope-family parent:
  `SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig`
- Kernel-boundary parent:
  `SIM-fovip-kernel-promise-boundary-port-contract`

## Candidate breeding shapes

- Pair A: `SIM-dalor` × `SIM-zukis`.
  This should produce envelope-shape children that reconcile dalor's explicit
  sender-proof slot with zukis's tagged Protocol-CID selector and
  protocol-owned later slots.
- Pair B: `SIM-dalor` × `SIM-fovip`.
  This should produce cross-layer children that test whether the same grid
  envelope can carry app/kernel messages without confusing kernel promises with
  authority, permissions, or conformance claims.
- Optional pair C: `SIM-zukis` × `SIM-fovip`.
  Use only if the runner configuration or manual review needs a cleaner
  baseline for the current envelope direction at the kernel boundary before
  comparing dalor-derived children.

## Scoring requirements

- Use the current PT-gated scorer.
- Prefer the forthcoming Rubric V4 once `TODO-busok` lands, because this run is
  explicitly cross-layer: envelope plus kernel plus higher-layer behavior.
- If Rubric V4 is not ready, run with the current scorer and mark the run as
  provisional pending v4 rescore.
- Score against scenarios that include:
  - small-device CBOR parsing of `42(pCID)`;
  - app/kernel messages using the same grid shape as wire messages;
  - storage and computation promises, not only send/receive promises;
  - unsupported pCID behavior;
  - broken or refused kernel implementation promises;
  - local trust updates after make/break history;
  - long-horizon ecosystem compatibility with CID/IPLD/Bluesky-style objects.

## Expected child properties

A strong child should:

- preserve slot `0` as `42(pCID)` unless it explicitly justifies a stronger
  alternative under current envelope evidence;
- treat `pCID` as the Protocol CID of the spec document, never as a payload CID;
- let the protocol named by `pCID` define slot `1` and later-slot roles;
- make sender signatures evidence of the sender's own promise over exact bytes;
- keep higher-layer promise accounting inside the payload protocol or
  app/kernel protocol, not in a universal envelope layer;
- describe app/kernel interaction as pCID-selected grid messages and ergonomic
  local adapters, not as traditional RPC authority;
- name the agent making each storage, compute, send, receive, key, lifecycle,
  namespace, reference, and evidence promise;
- separate host/runtime assumptions from kernel implementation promises;
- make all trust local and relationship-scoped.

## Rejection triggers

Reject or rework generated children that:

- use `permission`, `authorization`, `contract`, `policy enforcement`, or
  `conformance` as external authority rather than local promises and local
  trust updates;
- treat the kernel as a ruler, central arbiter, global namespace authority, or
  mandatory single process shape;
- move promise-accounting semantics into the universal envelope layer;
- use `pCID` to mean payload CID;
- reopen the `42(pCID)` envelope decision without new evidence;
- promote `sig_pCID` selector-shopping over protocol-owned later slots without
  an explicit, scored reason;
- propagate `port promise record` as the preferred term instead of
  `kernel implementation promises`.

## Subtasks

- [x] goban.1 Add this TODO and cross-list it in
  `protocols/wire-lab.d/TODO/TODO.md`.
- [x] goban.2 Record `DI-kipoz` for the dalor-anchored breeding intent.
- [x] goban.3 Confirm whether Rubric V4 from `TODO-busok` is ready enough for
  this cross-layer run.
- [x] goban.4 If Rubric V4 is not ready, mark the run provisional and plan an
  additive v4 rescore. V4 is ready, so no provisional fallback is needed.
- [x] goban.5 Configure a small canary/GA slice with Pair A
  (`SIM-dalor` × `SIM-zukis`).
- [x] goban.6 Configure a small canary/GA slice with Pair B
  (`SIM-dalor` × `SIM-fovip`).
- [x] goban.7 Decide whether Pair C (`SIM-zukis` × `SIM-fovip`) is needed as a
  baseline. Defer Pair C until Pair A/B results show a baseline gap.
- [x] goban.8 Give Steve exact commands to run; do not run provider-backed
  canary/GA commands from Codex: `/tmp/run-goban-dalor-breeding.sh`.
- [ ] goban.9 Review generated children for PT fit, envelope correctness,
  kernel-boundary correctness, and vocabulary regressions.
- [x] goban.10 Cull weak children before promotion review. Lasuv is rejected
  under `DI-dikat`; Vinag is mined into a new successor rather than promoted
  as scored.
- [ ] goban.11 If a child survives, record acceptance/promotion intent through
  the GA promotion process before moving anything into canonical
  `simulations/`.
- [x] goban.12 Create a simple promise-first Vinag successor that avoids
  service-registry and capability-table vocabulary.
- [x] goban.13 Configure a small Steve-run recalibration slice for the successor
  after the rubric refinement and culling changes are committed:
  `/tmp/run-hozif-anti-rpc-calibration.sh`.
- [x] goban.14 Revise Hozif to use one stable payload protocol pCID with
  promise and observation record kinds inside the payload.
- [x] goban.15 Cull Vinag after Hozif superseded it and Rubric V4 gained
  anti-RPC plus pCID-stability refinements.
- [x] goban.16 Cull Fozid as a generated proposal after review showed
  service/request/dispatch/authorization vocabulary and envelope-shape drift.
- [x] goban.17 Create a Hozif-derived successor that keeps per-surface
  promisers and outcome distinctions without Fozid's service/RPC vocabulary.

## Acceptance criteria

- A configured run exists that deliberately breeds `SIM-dalor` with
  `SIM-zukis` and `SIM-fovip`.
- Steve has exact commands for the run.
- Generated children are evaluated under PT-gated scoring and, when available,
  Rubric V4 layer-aware scoring.
- Any surviving child clearly composes the current envelope direction with
  kernel implementation promises.
- No child is promoted without explicit review, DI-backed intent, accepted score
  evidence, and canonical promotion steps.
