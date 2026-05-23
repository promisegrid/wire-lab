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

- [ ] tugoz.1 Write a concise gap note comparing `SIM-janov`, `SIM-maraz`,
  `SIM-natim`, `SIM-riliz`, and `SIM-sajar`, and pin down exactly what the new
  explicit signable-view sim must add. Source: `DI-dunat`.
- [ ] tugoz.2 Capture the external prior-art note for Ceramic, atproto, UCAN,
  and adjacent signed-envelope systems in a repo file that future sims can cite.
  Source: `DI-dunat`.
- [ ] tugoz.3 Create the standalone explicit signable-view simulation under
  `simulations/`, with `README.md`, `QUESTION.md`, and local protocol/spec files
  as needed. Source: `DI-dunat`.
- [ ] tugoz.4 Create the Gordian comparison sims under `simulations/`, keeping
  at least one payload/wrapper-focused specimen and one negative-control
  universal-envelope specimen. Source: `DI-dunat`.
- [ ] tugoz.5 Register the new sims in `simulations/README.md` and list them as
  active alternatives, not consensus, in `DEV-GUIDE-RESOURCES.md` once they
  exist. Source: `DI-dunat`.
- [ ] tugoz.6 Run a focused scored slice for the new sims against the chosen
  existing scenarios after the audit/backfill path is unblocked. Source:
  `DI-dunat`.

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
