# TODO 020 - TE-34 TE editing policy and holistic corpus reading: drive to DI

Track the work to drive TE-34 (`docs/thought-experiments/TE-20260502-212810-te-editing-policy-and-holistic-corpus.md`) from `needs DF` to a set of decided DIs that lock the TE editing policy, the per-protocol applicability rule, and the holistic-reading default. Once locked, sweep AGENTS files and back-update path references inside earlier TE files per the new policy's Cat-1 rule.

## Subtasks

- [ ] 020.1 Steve answers DF-34.1 (TE editing permissiveness). Recommended: 34.1.c (categorized: mechanical-in-place, navigational-append, substantive-supersede).
- [ ] 020.2 Steve answers DF-34.2 (per-protocol vs. uniform applicability). Recommended: 34.2.b (per-protocol; wire-lab harness locks its own; others adopt or override).
- [ ] 020.3 Steve answers DF-34.3 (holistic vs. single-TE reading by default). Recommended: 34.3.b (holistic by default; Alt-3.C is its day-to-day form).
- [ ] 020.4 Once 020.1-020.3 land, write a DI for each into this file.
- [ ] 020.5 Update `AGENTS.md`, `AGENTS-codex.md`, and `AGENTS-ppx.md` with a new "TE editing policy" section codifying the three regimes and the holistic-reading default. Cite TE-34 and the relevant DIs.
- [ ] 020.6 Sweep all TE files under `docs/thought-experiments/` for stale top-level path references (`specs/harness-spec-draft.md` -> `protocols/wire-lab.d/specs/harness-spec-draft.md`; old `TODO/0NN-...md` -> per-protocol `TODO-<timestamp>-...md`). Update in place per Cat-1; no top-of-file note needed because path renames are mechanical and self-explanatory in `git log -p`.
- [ ] 020.7 Add a "## Refinements" section to TE-1 (and any other earlier TE whose "Implications and future work" list now has resolved items) noting which downstream items have been filed since first drafting. This is a Cat-3 / Cat-4 edit per the new policy. Specifically for TE-1: note that the recommended set is awaiting DF in TODO 5 with a link, and that DR-006 is open and tracked there.
- [ ] 020.8 Update `docs/thought-experiments/README.md` (already done in this twig: editing-policy section added, path refs to `protocols/wire-lab.d/specs/harness-spec-draft.md` updated). Reconfirm after DF lands that the policy text matches the locked DIs.

## Decision Intent Log

(No DI yet. DF answers from Steve will populate this section. Expected DI IDs: DI-020-20260502-213103 through DI-020-20260502-213105 for DF-34.1 through DF-34.3.)

## Notes

- TE-34 carries the full alternative analysis. This file does not duplicate it; it tracks the decision-driving work.
- The recommended set is `(34.1.c, 34.2.b, 34.3.b)` per the TE.
- DF-34.1.c codifies what TE-24 and TE-26 already did: in-place rewrite with a top-of-file vocabulary note. It also formalizes the "## Refinements" append-only section convention which is new in this TE.
- The TODO 014 migration intentionally left stale paths inside TE files because the bot was applying the now-retired "do not back-edit" rule. Once DF-34.1 lands, subtask 020.6 sweeps those paths in place.
- DI-003 (proposals/pending/ contest artifacts immutable) is unaffected. The new TE editing policy applies to TE files only; contest / review artifacts remain byte-frozen.
- No linked DR. TE-34 is itself the deliberation; no separate DR is needed because the decision is fully framed inside the TE and locked by DI in this TODO. (Future TODOs may revisit this if the user prefers a DR-shaped artifact for editing-policy decisions; for now, the TE + DI pair is the audit object.)
