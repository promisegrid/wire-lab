# TODO-tugoz: Grid-envelope signable-view and Gordian comparison sims

## Prior aliases

None. This TODO was minted after the proquint-handle migration.

## Status

Open. This TODO owns the next envelope-design sim batch: one explicit
signable-view specimen for a minimal outer grid, plus a small Gordian
comparison family. It is separate from `TODO-tapur`, which owns GA-runner
machinery and rescoring workflow.

## Decision Intent Log

ID: DI-dunat
Date: 2026-05-22 20:07:19
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Track the next envelope-design follow-on as one harness TODO that
creates two explicit work items under `simulations/`: (a) a standalone
specimen where the outer envelope stays `grid([pCID, payload])` and any
required signature/proof lives inside a pCID-defined payload with an explicit
unsigned/signable view; and (b) a small Gordian comparison batch that tests
payload/wrapper-focused Gordian ideas and a negative-control universal-envelope
variant without treating Gordian as current consensus.
Intent: The current sim corpus already probes nested signed payloads, wrapper
proofs, and outer attestations, but it does not isolate the exact
"sign-inside-payload with an explicit signable projection" rule, and Gordian
was explored in chat/subagent analysis without becoming repo-tracked simulation
work.
Constraints: Do not rewrite or mutate already-scored sim trees. Keep the
universal outer-envelope question open. Reuse existing scenarios first. Any
Gordian material added to `DEV-GUIDE-RESOURCES.md` must remain in an
alternatives/open-work section until scored evidence exists.
Affects: `protocols/wire-lab.d/TODO/TODO-tugoz-grid-envelope-signable-view-and-gordian-comparison-sims.md`;
`protocols/wire-lab.d/TODO/TODO.md`; future `simulations/SIM-*/` trees for the
signable-view and Gordian batch; future `simulations/README.md`; future
`DEV-GUIDE-RESOURCES.md`; future prior-art notes under `docs/` or
`docs/research/`.

ID: DI-kafot
Date: 2026-05-22 23:33:34
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Capture `tugoz.1` and `tugoz.2` as two non-normative notes under
`docs/research/`: one gap note comparing the nearest existing envelope sims,
and one prior-art note covering Ceramic, atproto, and UCAN. Do not treat either
note as a TE, frozen spec, or current consensus statement.
Intent: The next useful step is to make the unresolved signable-view and
Gordian comparison work citeable from repo-owned artifacts before creating new
sims. `docs/research/` is the right home for these notes because they are
reference material for future simulations, not decisions by themselves.
Constraints: Do not lock the winning signable-view rule or Gordian batch shape
in these notes. Keep the current outer-envelope question open. The notes may
recommend what to test next, but the actual sim-creation choice stays with the
later tugoz subtasks.
Affects: `protocols/wire-lab.d/TODO/TODO-tugoz-grid-envelope-signable-view-and-gordian-comparison-sims.md`;
`docs/research/grid-envelope-signable-view-gap-20260522.md`;
`docs/research/grid-envelope-signature-prior-art-20260522.md`.

ID: DI-nohir
Date: 2026-05-23 00:15:13
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement the next `tugoz` batch as six standalone simulations:
three signature-structure probes modeled after atproto, Ceramic, and UCAN, plus
three Gordian comparison sims (payload/wrapper, universal-envelope negative
control, and selective disclosure). Use the first scored slice
`portable-signing-key-identity`, `live-crdt-audit-publication`,
`multi-embodiment-app-identity`, `chunk-feed-replication-sparse-advertisement`,
`device-bound-agent-physical-effect`, and `minimal-immutable-blob-app`.
Intent: The user expanded the earlier single signable-view idea into a direct
prior-art comparison batch and chose the broader 6-scenario first score slice.
That makes the implementation target a concrete six-sim family rather than one
specimen plus deferred comparisons.
Constraints: Keep the universal outer-envelope question open. Do not rewrite
scored sim trees. Keep these sims out of the current-consensus section of
`DEV-GUIDE-RESOURCES.md` until scored evidence exists. The atproto-like sim is
the only one in this batch that directly tests "proof carried in a reserved
payload slot while signing a named signable projection". Use medium scoring for
the first focused slice.
Affects: `protocols/wire-lab.d/TODO/TODO-tugoz-grid-envelope-signable-view-and-gordian-comparison-sims.md`;
`simulations/SIM-pamap-grid-envelope-signable-view-atproto-like/`;
`simulations/SIM-jumav-grid-envelope-wrapper-proof-ceramic-like/`;
`simulations/SIM-mipum-grid-envelope-signed-body-envelope-ucan-like/`;
`simulations/SIM-fitin-gordian-payload-wrapper/`;
`simulations/SIM-suzuf-gordian-universal-envelope-negative-control/`;
`simulations/SIM-vizan-gordian-selective-disclosure/`;
`simulations/README.md`; `DEV-GUIDE-RESOURCES.md`; focused `results/` and
`results/jobs/` artifacts for the first score slice.

## Scope

This TODO covers the next standalone simulation batch for the grid-envelope
question space:

- one explicit signable-view sim where the universal outer envelope remains
  `[pCID, payload]`, the payload bytes are pCID-defined, and any mandatory proof
  is carried inside the payload contract with a clearly named unsigned/signable
  projection;
- one or more Gordian-oriented comparison sims that test whether Gordian-style
  subject/assertion/proof structure belongs in a PromiseGrid payload/wrapper
  family or is a poor fit for the universal outer envelope;
- a small repo-facing prior-art note that captures Ceramic / atproto / UCAN /
  adjacent lessons relevant to this batch.

This TODO does not lock a winning envelope and does not change the current
consensus wording in `DEV-GUIDE-RESOURCES.md` by itself.

## Predecessor context

- `DEV-GUIDE-RESOURCES.md` currently records a consensus toward a small
  positional outer grid and pCID-owned inner semantics.
- `SIM-janov`, `SIM-maraz`, `SIM-natim`, `SIM-riliz`, and `SIM-sajar` are the
  nearest existing specimens, but none is the exact signable-view specimen now
  desired.
- The Gordian and Ceramic/atproto comparisons were explored in chat and
  subagent summaries but are not yet tracked in a repo-owned TODO or sim batch.

## DF/DI decisions needed before implementation

- Signable-view exact shape: decide whether the first dedicated specimen uses a
  named `payload_without_sig` projection, a deterministic slot-elision rule, or
  another explicit signable-view rule. Source: `DI-dunat`.
- Gordian batch shape: decide whether the first comparison batch is two sims
  (payload/wrapper + universal-envelope negative control) or a slightly larger
  family that also tests disclosure/selective-reveal pressure. Source:
  `DI-dunat`.
- Prior-art capture home: decide whether the Ceramic / atproto / UCAN note
  belongs under `docs/research/`, `docs/thought-experiments/`, or a local
  sim-batch note. Source: `DI-dunat`.
- Scenario slice: decide which existing scenarios are the minimum useful first
  score slice for these new sims. Source: `DI-dunat`.

## Subtasks

- [x] tugoz.1 Write a concise gap note comparing `SIM-janov`, `SIM-maraz`,
  `SIM-natim`, `SIM-riliz`, and `SIM-sajar`, and pin down exactly what the new
  explicit signable-view sim must add. Implemented in
  `docs/research/grid-envelope-signable-view-gap-20260522.md`. Source:
  `DI-dunat`; `DI-kafot`.
- [x] tugoz.2 Capture the external prior-art note for Ceramic, atproto, UCAN,
  and adjacent signed-envelope systems in a repo file that future sims can cite.
  Implemented in `docs/research/grid-envelope-signature-prior-art-20260522.md`.
  Source: `DI-dunat`; `DI-kafot`.
- [x] tugoz.3 Create the signature-structure comparison sims under
  `simulations/`: atproto-like explicit signable view, Ceramic-like wrapper
  proof, and UCAN-like signed-body envelope. Source: `DI-dunat`; `DI-nohir`.
- [x] tugoz.4 Create the Gordian comparison sims under `simulations/`: a
  payload/wrapper specimen, a universal-envelope negative control, and a
  selective-disclosure specimen. Source: `DI-dunat`; `DI-nohir`.
- [x] tugoz.5 Register the new sims in `simulations/README.md` and list them as
  active alternatives, not consensus, in `DEV-GUIDE-RESOURCES.md` once they
  exist. Source: `DI-dunat`; `DI-nohir`.
- [ ] tugoz.6 Run a focused scored slice for the six new sims against
  `portable-signing-key-identity`, `live-crdt-audit-publication`,
  `multi-embodiment-app-identity`, `chunk-feed-replication-sparse-advertisement`,
  `device-bound-agent-physical-effect`, and `minimal-immutable-blob-app`.
  Attempted via `ga-tugoz-20260523-072345`, which produced 30 valid parent
  results and 6 failed `SIM-suzuf-gordian-universal-envelope-negative-control`
  cells because the provider emitted an 8-axis score object that omitted
  `promise_vocabulary` and `simplicity_durability`. Narrow retry
  `ga-tugoz-20260523-073108-suzuf` reproduced the same `SIM-suzuf`
  provider-output failure. Leave this item open until those cells are recovered
  or a scoring-policy exception is approved. Source: `DI-dunat`; `DI-nohir`.

## Validation and acceptance criteria

- The explicit signable-view specimen states exactly which bytes are signed and
  what rule reconstructs or projects the unsigned/signable view. Source:
  `DI-dunat`.
- The Gordian comparison batch keeps the universal outer-envelope question open
  and makes the trade-off against the minimal positional outer grid explicit.
  Source: `DI-dunat`.
- The prior-art note is specific enough that future decisions can cite it
  instead of re-deriving Ceramic / atproto / UCAN comparisons from chat memory.
  Source: `DI-dunat`.
- Any `DEV-GUIDE-RESOURCES.md` update from this TODO must keep these sims in an
  alternatives/open-work section until scored evidence exists. Source:
  `DI-dunat`.
