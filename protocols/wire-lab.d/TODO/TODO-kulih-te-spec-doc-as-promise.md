# TODO-kulih - TE-nibar Spec doc as promise: drive to DI

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-10` (integer alias)
- `TODO-20260429-173837` (timestamp alias and pre-migration filename)

Track the work to drive TE-21 (`docs/thought-experiments/TE-20260429-173520-spec-doc-as-promise.md`) from `needs DF` to a set of decided DIs that lock the spec-doc-as-promise vocabulary for the wire-lab.

## Subtasks

- [ ] 010.1 Steve answers DF-21.1 (layered Alt-E vs. doc-only Alt-D vs. other). Recommended: 1.a (Alt-E).
- [ ] 010.2 Steve answers DF-21.2 (assumption / open-questions / known-issues lists: required sections, best-practice, or required-when-present). Recommended: 2.a (required).
- [ ] 010.3 Steve answers DF-21.3 (peer adoption promises: structured metadata, commentary-only, or required-commentary-optional-structure). Recommended: 3.a (structured).
- [ ] 010.4 Steve answers DF-21.4 (does this TE imply rename or restructure of `protocols/wire-lab.d/specs/harness-spec-draft.md`). Recommended: 4.a (no, defer to TE-rujak).
- [ ] 010.5 Once 010.1-010.4 land, write a DI for each into this file.
- [ ] 010.6 If 010.2 lands as 2.a (required), update `protocols/wire-lab.d/specs/harness-spec-draft.md` to add or formalize the three normative sections: Assumptions, Open Questions (already present as \u00a711), Known Issues. Cross-link them to TE-nibar and the relevant DI(s).
- [ ] 010.7 If 010.3 lands as 3.a (structured), surface the peer adoption metadata as a future TE (TE-lozip placeholder) so the wire shape can be designed.
- [ ] 010.8 Open TE-rujak (spec-doc-store layout) on a fresh twig once TE-nibar / TODO-kulih vocabulary is locked.
- [ ] 010.9 Decide whether protocol specs require layer-specific promise-vocabulary sections, 100-year pressure-test sections, and layperson-readable mental-model / easy-implementation summaries. This absorbs turn-177's `UT-177.e`, `UT-177.f`, and `UT-177.i` spec-shape fallout; TE-dajot remains the citable 100-year constraint while this TODO owns whether those obligations become required spec sections. `SIM-rusap-promise-accounting-records` now carries the simulation-facing peer-local promise-accounting and mental-model pressure that should inform this decision, without replacing TE-nibar's DF/DI path. Its `SCENARIOS.md` file adds concrete cases for kept promises, refusal, corruption, cross-layer decisions, sparse retention, identity rotation, and layperson explanation. `DR-robon` now owns the explicit decision request for this turn-177 spec-shape addition. Source: `DI-navod`; `DI-pator`; `DI-davov`.
- [ ] 010.10 Answer `DR-robon` / DF-robon.1 through DF-robon.4 before closing
  010.9. Source: `DI-majib`.
- [ ] 010.11 Use `simulations/SIM-ranib-spec-requirement-sections/` as the
  simulation question home for 010.9 / `DR-robon` spec-section pressure before
  deciding whether promise-vocabulary, 100-year, and
  layperson/easy-implementation sections are required. Source: `DI-pukap`.

## Question log

- 2026-05-17: `DR-robon` asks whether PromiseGrid protocol specs should
  require layer-specific promise-vocabulary sections, 100-year pressure-test
  sections, and layperson-readable mental-model / easy-implementation
  summaries. Source: `DI-davov`.
- 2026-05-17: `DR-robon` now has an unanswered next-DF packet and acceptance
  criteria for deciding 010.9. Source: `DI-majib`.
- 2026-05-17: `simulations/SIM-ranib-spec-requirement-sections/` now captures
  the simulation-facing version of the 010.9 / `DR-robon` question. This does
  not answer `DR-robon`. Source: `DI-pukap`.

## Decision Intent Log

(No DI yet. DF answers from Steve will populate this section.)

## Notes

- TE-nibar carries the full alternative analysis (Alt-A through Alt-E) and six scenarios (S1-S6). This file does not duplicate that analysis; it tracks the decision-driving work.
- The recommended set is `(1.a, 2.a, 3.a, 4.a)` per the TE. Reason: lock the layered framing fully, make the three lists structural, give peer-level adoption first-class promise machinery, defer layout questions to a follow-on TE.
- Linked DR: `DR-robon` covers the turn-177 010.9 addition. A broader TE-nibar
  DR may still be created once subtasks 010.1-010.4 begin landing.
- Companion TODOs: TODO-misul (DI-provenance backfill) and TODO-diliz (DR backfill for \u00a711) become more concrete after TE-nibar locks; consider revising their scope statements after DF lands.
