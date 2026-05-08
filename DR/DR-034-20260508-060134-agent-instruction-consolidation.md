# DR-034 - Agent instruction consolidation

DR-ID: DR-034-20260508-060134
Date: 2026-05-08 06:01:34
Asked by: stevegt@t7a.org (Steve Traugott)
State: decided
Question: Should `AGENTS-ppx.md` be moved wholesale into `AGENTS.md`, or should repo-wide rules move to `AGENTS.md` while `AGENTS-ppx.md` remains a Perplexity-specific role overlay?
Why this blocks progress: `AGENTS-ppx.md` contains both repo-wide rules and Perplexity-only runtime mechanics. Leaving duplicated rules in role files creates multiple sources of truth, but moving the whole file into `AGENTS.md` would impose bot identity, cloud sandbox paths, private session logging, and `ppx/main` lifecycle rules on agents that should not follow them.
Affects:
- `AGENTS.md` — canonical repo-wide rules and glossary.
- `AGENTS-ppx.md` — Perplexity role overlay.
- `AGENTS-codex.md` — Codex role overlay.
- `protocols/wire-lab.d/TODO/TODO-milum-agent-instruction-consolidation.md` — parent TODO and DI.
- `protocols/wire-lab.d/TODO/TODO.md` — master cross-list gains TODO-milum.
Unblocks: safer agent onboarding, fewer duplicated protocol rules, clearer distinction between canonical repo rules and role-specific runtime procedures.
Waiting on: DI-034-20260508-060134
Decision: Do not move `AGENTS-ppx.md` wholesale. Move shared rules and definitions into `AGENTS.md`; keep `AGENTS-ppx.md` and `AGENTS-codex.md` as overlays that point to `AGENTS.md` for shared policy and retain only role-specific procedures.
Linked DI: DI-034-20260508-060134 in `protocols/wire-lab.d/TODO/TODO-milum-agent-instruction-consolidation.md`
Related commits: pending
Last updated: 2026-05-08 06:01:34

## Alternatives considered

- **Alt-A: move all `AGENTS-ppx.md` content into `AGENTS.md`.** Rejected because it would make Perplexity-only identity, credential, branch, and private logging rules appear repo-wide.
- **Alt-B: leave files as-is.** Rejected because duplicated TE editing and generic forbidden-action text already creates multiple sources of truth.
- **Alt-C: canonical root plus thin role overlays.** Selected. Repo-wide rules live in `AGENTS.md`; role overlays retain only environment, identity, lifecycle, and role-specific operational details.
