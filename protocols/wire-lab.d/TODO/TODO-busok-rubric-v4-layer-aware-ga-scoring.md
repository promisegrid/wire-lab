# TODO-busok: Rubric V4 layer-aware GA scoring

## Prior aliases

This TODO was minted directly under the post-`TE-mumuv` proquint handle scheme.
No integer or timestamp alias exists.

## Status

Implemented pending commit and Steve-run rescore. Rubric V4 expands GA scoring
so envelope, kernel, and higher-layer protocol/app candidates are evaluated on
the PromiseGrid layer they actually claim to occupy. The first calibration
target is an additive rescore of `ga-canary-20260525-kernel-resolution`.

## Decision Intent Log

### DI-ripuz

ID: DI-ripuz
Date: 2026-05-25 15:01:45
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Expand GA scoring to Rubric V4 with explicit layer-aware coverage for
wire protocol envelopes, kernel implementation promises, and higher-layer/app
protocols, then rescore `ga-canary-20260525-kernel-resolution` additively.
Intent: The v3 rubric mostly scored the kernel canary correctly, but it does not
explicitly cover PromiseGrid's major design strata. It also risks penalizing
locally scoped kernel implementation promises as generic port/profile claims.
Rubric V4 should reflect `DN-jotob` for envelope scoring, reward local kernel
implementation promises, and score higher-layer/app protocols on Promise Theory
behavior.
Constraints: Do not overwrite existing score evidence. Do not silently change
`ga-rubric-20260523-v3`. Commit v4 source changes before any provider-backed v4
rescore. Steve runs canary/GA commands. Preserve historical TEs unless TE policy
allows a refinement or successor.
Affects: `tools/ga-runner/`; `protocols/wire-lab.d/TODO/TODO.md`;
`protocols/wire-lab.d/TODO/TODO-binag-promisegrid-kernel-design-resolution.md`;
`simulations/SIM-fovip-kernel-promise-boundary-port-contract/`;
`docs/research/DN-jotob-grid-envelope-tag42-variable-outer-slots.md`;
future `results/state/ga-canary-20260525-kernel-resolution-rubric-v4.json`.

## Locked direction

- Use a real new rubric/result version: `promisegrid.ga.result.v4` and
  `ga-rubric-20260525-v4`.
- Add score axes rather than only changing prompt prose.
- Additive rescore only; never overwrite committed v3 evidence.
- Rescore the `ga-canary-20260525-kernel-resolution` parent slice first.
- Replace active wording `port promise record` with
  `kernel implementation promises`.

## Proposed Rubric V4 axes

Keep existing v3 axes and add:

- `envelope_discipline`: rewards `DN-jotob` alignment: CBOR
  `grid([42(pCID), payload, ...])`, Protocol-CID semantics, protocol-owned later
  slots, local unknown-pCID handling, and no universal proof-slot overreach.
- `kernel_implementation_promises`: rewards explicit kernel implementation
  promises, host/runtime assumptions separated from promises, unsupported
  pCIDs/features, pCID adapter mapping, local evidence, voluntary namespace /
  reference / resource behavior, and no conformance authority.
- `app_protocol_promise_semantics`: rewards higher-layer/app protocols that
  model storage, computation, send/receive, reciprocal promises, selective
  sending, local trust, promise-as-capability-token behavior, and make/break
  evidence without command/permission framing.

## Subtasks

- [x] busok.1 Add this TODO and cross-list it in
  `protocols/wire-lab.d/TODO/TODO.md`.
- [x] busok.2 Add the DI above.
- [x] busok.3 Cite `DI-ripuz` from all non-trivial v4 scoring code/doc changes.
- [x] busok.4 Update active terminology from `port promise record` to
  `kernel implementation promises`.
- [x] busok.5 Implement `promisegrid.ga.result.v4` /
  `ga-rubric-20260525-v4` in `tools/ga-runner`.
- [x] busok.6 Add `DN-jotob` envelope guidance to the scoring prompt.
- [x] busok.7 Add kernel implementation promise guidance to the scoring prompt.
- [x] busok.8 Add higher-layer/app Promise Theory guidance to the scoring
  prompt.
- [x] busok.9 Update tests so v4 requires the three new axes while v1/v2/v3
  remain historically valid.
- [ ] busok.10 Commit v4 source/doc changes before any v4 scoring run.
- [x] busok.11 Configure the additive v4 kernel rescore using the same parent
  slice as `ga-canary-20260525-kernel-resolution`.
- [x] busok.12 Ask Steve to run `/tmp/run-busok-v4-kernel-rescore.sh` after
  committing this source/doc change set.
- [ ] busok.13 Compare v3 vs v4 drift: rank order, PT gate status, `SIM-fovip`
  stability, `SIM-funas` downgrade stability, and new-axis scores.
- [ ] busok.14 Decide whether v4 is safe for broader GA runs, needs prompt
  repair, or requires another calibration slice.

## Acceptance criteria

- Rubric V4 distinguishes envelope, kernel, and higher-layer/app fitness.
- `kernel implementation promises` is the active term in current kernel docs.
- Existing v3 evidence remains append-only and unmodified.
- The v4 kernel rescore is additive and comparable against committed v3
  evidence.
