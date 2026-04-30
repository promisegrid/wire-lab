# TODO 006 - Backfill DI provenance into specs/harness-spec-draft.md settled statements

Originally listed as subtask 001.5 of `TODO/001-perplexity-computer-onboarding.md`. Split out into its own TODO so that TODO 001 closes cleanly.

`specs/harness-spec-draft.md` contains many settled statements (e.g. §1.1 promise-frame shape, §2 trust ledger, §3 currency, §10a invariants) that locked in earlier conversations without DR/DI records. Backfill DI provenance for those statements so the spec is auditable: each settled paragraph should be traceable to one or more DI entries that record who decided what, when, and why.

## Subtasks

- [ ] 006.1 Identify the settled statements in `specs/harness-spec-draft.md` that lack DI provenance. Produce a candidate list with file/line ranges.
- [ ] 006.2 For each candidate, write a DI entry capturing the decision in retrospective form (status: `recorded-after-the-fact`).
- [ ] 006.3 Add inline `(see DI-...)` references in `specs/harness-spec-draft.md` so future readers can navigate from spec to DI.
- [ ] 006.4 Decide whether retrospective DIs live in their own TODO file (e.g., `TODO/006-...`) or in a single archive file (e.g., `archive/retrospective-DIs.md`).

## Decision Intent Log

(No DI yet. This TODO drives provenance work; DIs land in either this file or the archive file once 006.4 decides.)

## Notes

- This work is a quality/provenance backfill, not a new design decision. Existing semantics in `specs/harness-spec-draft.md` are not changed.
- Track scope: ~50 candidate statements estimated; likely produces 20-40 DI entries after consolidation.
