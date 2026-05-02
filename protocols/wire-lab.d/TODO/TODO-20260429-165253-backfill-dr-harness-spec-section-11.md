# TODO 007 - Backfill DR records for protocols/wire-lab.d/specs/harness-spec-draft.md §11 open questions

Originally listed as subtask 001.6 of `protocols/wire-lab.d/TODO/TODO-20260429-030146-perplexity-computer-onboarding.md`. Split out into its own TODO so that TODO 001 closes cleanly.

`protocols/wire-lab.d/specs/harness-spec-draft.md` §11 lists open design questions about the Wire Lab. Each is a real Decision Request that deserves a DR file under `DR/` so that the question is citable, has an owner, has a state (`open`, `decided`, `superseded`), and tracks alternatives considered. Backfill those DRs.

## Subtasks

- [ ] 007.1 Enumerate the open questions in `protocols/wire-lab.d/specs/harness-spec-draft.md` §11. Produce a candidate list with section anchors.
- [ ] 007.2 For each, write a `DR/DR-NNN-<timestamp>-<slug>.md` with `State: open` and the question phrased per the AGENTS.md DR template.
- [ ] 007.3 Cross-link `protocols/wire-lab.d/specs/harness-spec-draft.md` §11 entries to their DR files so the spec is navigable from question to DR.
- [ ] 007.4 Decide whether §11 should remain a question list at all once each question has its own DR, or whether §11 becomes an index pointing into `DR/`.

## Decision Intent Log

(No DI yet. This TODO drives DR backfill; DIs land here as those questions transition to `decided`.)

## Notes

- This work does not answer the open questions; it just gives each one a DR file so future answers have a citable anchor.
- The 18-question count cited in the original 001.6 entry was approximate; 007.1 produces the authoritative count.
