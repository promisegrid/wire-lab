# TODO 010 - TE-21 Spec doc as promise: drive to DI

Track the work to drive TE-21 (`docs/thought-experiments/TE-20260429-173520-spec-doc-as-promise.md`) from `needs DF` to a set of decided DIs that lock the spec-doc-as-promise vocabulary for the wire-lab.

## Subtasks

- [ ] 010.1 Steve answers DF-21.1 (layered Alt-E vs. doc-only Alt-D vs. other). Recommended: 1.a (Alt-E).
- [ ] 010.2 Steve answers DF-21.2 (assumption / open-questions / known-issues lists: required sections, best-practice, or required-when-present). Recommended: 2.a (required).
- [ ] 010.3 Steve answers DF-21.3 (peer adoption promises: structured metadata, commentary-only, or required-commentary-optional-structure). Recommended: 3.a (structured).
- [ ] 010.4 Steve answers DF-21.4 (does this TE imply rename or restructure of `protocols/wire-lab.d/specs/harness-spec-draft.md`). Recommended: 4.a (no, defer to TE-22).
- [ ] 010.5 Once 010.1-010.4 land, write a DI for each into this file.
- [ ] 010.6 If 010.2 lands as 2.a (required), update `protocols/wire-lab.d/specs/harness-spec-draft.md` to add or formalize the three normative sections: Assumptions, Open Questions (already present as \u00a711), Known Issues. Cross-link them to TE-21 and the relevant DI(s).
- [ ] 010.7 If 010.3 lands as 3.a (structured), surface the peer adoption metadata as a future TE (TE-23 placeholder) so the wire shape can be designed.
- [ ] 010.8 Open TE-22 (spec-doc-store layout) on a fresh twig once TE-21 / TODO 010 vocabulary is locked.

## Decision Intent Log

(No DI yet. DF answers from Steve will populate this section.)

## Notes

- TE-21 carries the full alternative analysis (Alt-A through Alt-E) and six scenarios (S1-S6). This file does not duplicate that analysis; it tracks the decision-driving work.
- The recommended set is `(1.a, 2.a, 3.a, 4.a)` per the TE. Reason: lock the layered framing fully, make the three lists structural, give peer-level adoption first-class promise machinery, defer layout questions to a follow-on TE.
- Linked DR: to be created as DR-010 once subtasks 010.1-010.4 begin landing. For now, the TE itself stands as the open-question record.
- Companion TODOs: TODO 006 (DI-provenance backfill) and TODO 007 (DR backfill for \u00a711) become more concrete after TE-21 locks; consider revising their scope statements after DF lands.
