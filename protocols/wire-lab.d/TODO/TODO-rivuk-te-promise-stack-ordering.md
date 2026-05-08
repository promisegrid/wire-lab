# TODO-rivuk - TE-famar Promise-stack ordering: drive to DI

## Status

Closed as superseded/no-longer-applicable on 2026-05-08 by TE-havib
DF-36.2 Alt-2.A revised and TODO-kugod DI-runuh. Do not answer DF-1.1
through DF-1.4 as live questions; promise-stack was retired as a separate
hypothesis and future work follows payload recursion under
`grid <pcid> <payload>`.

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-5` (integer alias)
- `TODO-20260429-164955` (timestamp alias and pre-migration filename)

Historical scope: track the work to drive TE-famar from `needs DF` to a
set of decided DIs that would have locked the wire-and-evaluation convention
for multi-frame promise stacks. That scope was superseded before the DF queue
landed.

## Subtasks

- [x] 005.1 Steve answers DF-1.1 (peeling vs. projection vs. Alt-E hybrid). Superseded by TE-havib DF-36.2; no promise-stack DI created.
- [x] 005.2 Steve answers DF-1.2 (criticality-flag location: per-frame, per-assertion-type spec, or hybrid). Superseded by TE-havib DF-36.2; no promise-stack DI created.
- [x] 005.3 Steve answers DF-1.3 (wire-encoding direction: outermost-first vs. innermost-first). Superseded by TE-havib DF-36.2; no promise-stack DI created.
- [x] 005.4 Steve answers DF-1.4 (position-convention authority: who declares whether an assertion-type has a normative position). Superseded by TE-havib DF-36.2; no promise-stack DI created.
- [x] 005.5 Once 005.1-005.4 land, write a DI for each into this file. Replaced by TODO-kugod DI-runuh, which records closure without promise-stack ordering DIs.
- [x] 005.6 Update `protocols/wire-lab.d/specs/harness-spec-draft.md` §1.1 to reference the locked DIs and the position-convention rule. Not applicable after promise-stack retirement; no harness-spec sweep in this pass.
- [x] 005.7 Update `DR/DR-006-20260429-164729-promise-stack-ordering.md` from `open` to `decided`. Completed as closed/superseded, not decided, because the old question is no longer applicable.

## Decision Intent Log

No promise-stack ordering DI was created. The controlling closure record is
DI-runuh in `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`.

## Notes

- TE-famar carries the full historical alternative analysis (Alt-A through Alt-E) and the six scenarios (S1-S6) that drove the old recommended set.
- The old recommended set was `(1.1.a, 1.2.c, 1.3.a, 1.4.d)` per TE-famar. That set is preserved as historical analysis, not as an active decision queue.
- Linked DR: `DR/DR-006-20260429-164729-promise-stack-ordering.md`.
- Current forward pointer: TE-havib DF-36.2 Alt-2.A revised, TE-lozip, and `docs/essays/congruence-convergence-and-the-grid.md` §3.1.
