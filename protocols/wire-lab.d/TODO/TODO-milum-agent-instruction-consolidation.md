# TODO-milum: Agent instruction consolidation

## Prior aliases

This TODO is being filed after the TE-mumuv proquint migration locked
2026-05-07, so it is minted directly under `TODO-milum`. No prior
integer or timestamp alias.

For convenience in cross-references during the rollout window, the
master cross-list (`protocols/wire-lab.d/TODO/TODO.md`) carries this
file under the integer alias `TODO-34`.

## Status

Implemented. This TODO records the consolidation decision that keeps
`AGENTS.md` canonical for repo-wide rules and keeps `AGENTS-ppx.md` /
`AGENTS-codex.md` as role overlays.

## Decision Intent Log

### DI-034-20260508-060134

- ID: DI-034-20260508-060134
- Date: 2026-05-08 06:01:34
- Status: active
- Author: stevegt@t7a.org (Steve Traugott)
- Decision: `AGENTS.md` is the canonical home for repo-wide agent protocol, workflow, and vocabulary rules. `AGENTS-ppx.md` and `AGENTS-codex.md` are role overlays that may add stricter role-specific constraints, environment setup, branch lifecycle, credential handling, and private logging procedures, but must point back to `AGENTS.md` for shared rules rather than duplicating them.
- Intent: Reduce multiple sources of truth while preserving the important distinction between repo-wide protocol and role-specific operational mechanics.
- Constraints:
    - Do not move Perplexity-only runtime material into `AGENTS.md`: cloud sandbox paths, bot git identity, `ppx/main` lifecycle, session logging, private remote, PAT handling, redact-last, and Carry-J2 remain in `AGENTS-ppx.md`.
    - Do not move Codex-only review workflow into `AGENTS.md`; keep local-review and Steve-authored-work procedures in `AGENTS-codex.md`.
    - Historical references to `AGENTS-ppx.md` in old TEs, TODOs, DRs, and audit reports remain historical evidence and are not swept by this consolidation.
    - Carry-J2 remains ppx-only unless a later DI explicitly promotes selected bundles repo-wide.
- Affects: `AGENTS.md`; `AGENTS-ppx.md`; `AGENTS-codex.md`; `protocols/wire-lab.d/TODO/TODO.md`; `DR/DR-034-20260508-060134-agent-instruction-consolidation.md`.

## Subtasks

- [x] 034.1 Review `AGENTS-ppx.md`, `AGENTS-codex.md`, and `AGENTS.md` for duplication and role-specific content.
- [x] 034.2 Reject wholesale migration of `AGENTS-ppx.md` into `AGENTS.md`; choose canonical root plus role overlays.
- [x] 034.3 Promote repo-wide rules to `AGENTS.md`: instruction architecture, TE actor naming, DI authorship semantics, no-PR/no-force-push defaults, secret hygiene, and shared glossary.
- [x] 034.4 Thin role overlays by replacing duplicated TE editing policy, generic forbidden-action lists, and shared glossary text with pointers to `AGENTS.md`.
- [x] 034.5 Record the decision in this TODO and in the linked DR.
