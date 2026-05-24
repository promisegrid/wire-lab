# DR-pijum - Codex local author identity

DR-ID: DR-pijum
Date: 2026-05-23 23:27:04
Asked by: jessica@t7a.org (Jessica)
State: decided
Question: Should local Codex work in this clone keep impersonating Steve, use
the actual configured local git author identity, or adopt a separate Codex-only
principal identity?
Why this blocks progress: The current `AGENTS-codex.md` rule hard-codes Steve's
identity and blocks authorized local work by other humans in this clone. That
creates false attribution for commit metadata and prevents the approved
simulation work from proceeding under the correct local operator identity.
Affects: `AGENTS-codex.md`; `protocols/wire-lab.d/TODO/TODO-tapaf-codex-local-author-identity.md`.
Unblocks: TODO-tapaf implementation; subsequent local Codex work in this clone
under real local author identities.
Waiting on: DI-nufit
Decision: Use the actual local git author identity configured in the current
clone for local Codex work. Keep the existing `ppx` bot identity and Steve-only
`main` push posture unchanged.
Linked DI: DI-nufit
Related commits:
Last updated: 2026-05-23 23:27:04

## Event log

- 2026-05-23 23:27:04 — Opened and decided from chat after the local clone was
  found to be configured as `Jessica Traugott <jessica@t7a.org>` while
  `AGENTS-codex.md` still required impersonating Steve.

## Evidence

- `AGENTS-codex.md` currently says local Codex work must act as Steve and use
  Steve's git identity.
- `TE-lusut` and `TODO-dutaz` already lock a distinct identity for the `ppx`
  bot, so the open question here is local Codex attribution, not bot
  attribution.
- The user explicitly asked to allow multiple authors and confirmed that the
  local Codex rule should use the actual configured local git author identity
  while leaving the bot identity rules untouched.
