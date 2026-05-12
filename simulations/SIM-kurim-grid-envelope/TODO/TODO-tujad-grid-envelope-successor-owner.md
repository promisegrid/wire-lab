# TODO-tujad: Grid-envelope successor owner

## Prior aliases

None. This TODO is created directly as a sim-local successor owner under
`rusis.10`.

## Status

Open. This TODO is the concrete successor-owner record for TE-40
transferred-open rows `UT-157.a`, `UT-157.c`, and `UT-158.f`.
Successor-owner routing into this TODO was locked under `DI-mosor` in
`protocols/wire-lab.d/TODO/TODO-rusis-simulation-split-and-specimen-relocation.md`.
Seed anchors were established earlier under `DI-nijon`.

## Decision Intent Log

ID: DI-joroh
Date: 2026-05-12 08:44:53
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Resolve the turn-157 grid-envelope replay cleanup by capturing the
candidate envelope inventory and `grid([pcid, payload])` working hypothesis in
this sim-local successor owner, without locking a canonical envelope or creating
a protocol tree.
Intent: Turn 157 contains load-bearing design material that should not remain
only in replay memory. Capturing the alternatives and hypothesis here gives the
grid-envelope lineage a concrete owner for the open replay residue while
preserving the fact that the hypothesis remains unproven.
Constraints: Preserve `Env-1` through `Env-5` as candidate inventory. Do not
decide the final PromiseGrid envelope. Do not create `protocols/grid-envelope.d/`
or draft a grid-envelope spec in this cleanup pass. Leave turn-158 protocol/spec
work open under `tujad.3`.
Affects: `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`;
`protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`;
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`.

## Scope

This TODO owns grid-envelope follow-on that was previously parked under
`TODO-kugod` as "until grid-envelope successor exists":

- candidate envelope inventory ownership;
- `grid([pcid, payload])` working-hypothesis prose ownership;
- concrete successor planning for grid-envelope protocol directory/spec
  work in this lineage.

Anchor seed note:
`simulations/SIM-kurim-grid-envelope/seed/extraction-sources.md`.

## Candidate Envelope Inventory

Turn 157 named these candidate envelope alternatives as inputs to later
grid-envelope and envelope-bakeoff work. This inventory records the candidates;
it does not select a winner. Source: `DI-joroh`.

- `Env-1`: `grid([pcid, payload])`. A two-element CBOR array where the first
  element is a pCID identifying which protocol, handler, or assertion type
  interprets the second element. This is the current working hypothesis. A
  payload may itself be another `grid([pcid, payload])` value if recursion is
  needed.
- `Env-2`: Promise stack of grid frames. A CBOR sequence of
  `grid([pcid, payload])` frames where stack semantics apply at the sequence
  level and the grid shape applies at each frame. This candidate reconciles the
  grid hypothesis with the earlier TE-famar promise-stack work.
- `Env-3`: Bare CBOR with no shared envelope. Each protocol chooses its own
  message shape. This is maximally permissive, but it may leave the harness
  without a shared parser to exercise across candidate transports.
- `Env-4`: Capability-port triplet. A direct `(promiser, assertion, body)`
  structure with no grid indirection. This is closest to the older
  harness-spec `Promise` shape.
- `Env-5`: Tagged union over `Env-1` and `Env-2`. Single-frame messages use
  `grid([pcid, payload])`; multi-frame messages use a stack of grid frames; a
  top-level tag distinguishes the two cases.

## Grid Envelope Working Hypothesis

`grid([pcid, payload])` is the current working hypothesis for a
transport-agnostic message envelope, but turn 157 explicitly says it has not
been proven. This simulation owns that hypothesis as a candidate specimen to be
tested against alternatives, not as a settled harness rule. Source: `DI-joroh`.

The harness may use this inventory to construct later bakeoffs, but any final
canonical-envelope decision still needs its own TE/DF/DI path. Turn 158's
apparatus-vs-specimen correction remains controlling: the harness compares
candidate envelopes rather than declaring this one canonical in advance.

## Subtasks

- [x] tujad.1 Materialize the candidate envelope inventory owner record
  for `UT-157.a`.
- [x] tujad.2 Materialize the `grid([pcid, payload])`
  working-hypothesis owner record for `UT-157.c`.
- [ ] tujad.3 Define and track the concrete successor path for
  grid-envelope protocol directory/spec work for `UT-158.f`.
- [x] tujad.4 Back-link resulting decisions and artifacts to
  `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`.
