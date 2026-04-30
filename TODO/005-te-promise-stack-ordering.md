# TODO 005 - TE-1 Promise-stack ordering: drive to DI

Track the work to drive TE-1 (`docs/thought-experiments/TE-20260427-180000-promise-stack-ordering.md`) from `needs DF` to a set of decided DIs that lock the wire-and-evaluation convention for multi-frame promise stacks.

## Subtasks

- [ ] 005.1 Steve answers DF-1.1 (peeling vs. projection vs. Alt-E hybrid). Recommended: 1.1.a (Alt-E hybrid). Recording form: a DI entry in this file.
- [ ] 005.2 Steve answers DF-1.2 (criticality-flag location: per-frame, per-assertion-type spec, or hybrid). Recommended: 1.2.c (hybrid).
- [ ] 005.3 Steve answers DF-1.3 (wire-encoding direction: outermost-first vs. innermost-first). Recommended: 1.3.a (outermost-first).
- [ ] 005.4 Steve answers DF-1.4 (position-convention authority: who declares whether an assertion-type has a normative position). Recommended: 1.4.d (per-assertion-type — each assertion's spec declares it).
- [ ] 005.5 Once 005.1-005.4 land, write a DI for each into this file.
- [ ] 005.6 Update `specs/harness-spec-draft.md` §1.1 to reference the locked DIs and the position-convention rule.
- [ ] 005.7 Update `DR/DR-006-20260429-164729-promise-stack-ordering.md` from `open` to `decided`.

## Decision Intent Log

(No DI yet. DF answers from Steve will populate this section.)

## Notes

- TE-1 carries the full alternative analysis (Alt-A through Alt-E) and the six scenarios (S1-S6) that drove the recommended set. This file does not duplicate that analysis; it tracks the decision-driving work.
- The recommended set is `(1.1.a, 1.2.c, 1.3.a, 1.4.d)` per the TE. Reason: maximally additive, leaves room for assertion-type specialization, doesn't lock down what isn't yet known.
- Linked DR: `DR/DR-006-20260429-164729-promise-stack-ordering.md`.
