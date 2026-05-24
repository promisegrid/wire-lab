# TE-puzud: Codex local author identity

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-puzud

## Status

decided

## Decision under test

When local Codex work happens inside this clone, should the agent continue
impersonating Steve, use the actual configured local git author identity, or
adopt a separate Codex-only principal?

## Assumptions

- The existing `ppx` bot identity remains separately locked by `TE-lusut` and
  `DI-001-20260428-195700`.
- This question is about local Codex work in a human-operated clone, not about
  the hosted `ppx` bot workflow.
- Multiple authorized humans may work in the same local clone over time.
- Commit author metadata should reflect the actual accountable local human
  operator unless a stronger locked rule says otherwise.

## Alternatives

- **Alt-A — keep impersonating Steve.** Every local Codex session commits as
  Steve regardless of the configured local git identity.
- **Alt-B — use the actual configured local git author identity.**
  Local Codex sessions use whatever real local identity is configured in the
  clone, such as Steve or Jessica. The `ppx` bot identity remains separate.
- **Alt-C — create a separate Codex-only principal.** Local Codex sessions
  commit under a dedicated Codex identity distinct from both Steve and the
  `ppx` bot.

## Scenario analysis

### Scenario 1 — Steve operating locally

Alice is Steve. Under Alt-A and Alt-B, attribution is correct because Steve's
configured identity matches the actual operator. Alt-C creates a new extra
principal even though Steve remains the accountable human reviewer and committer.

### Scenario 2 — Jessica operating locally

Alice is Jessica, working in the same clone with explicit authorization.
Alt-A produces false authorship because the commit metadata says Steve did work
he did not do. Alt-B records Jessica accurately while preserving the existing
bot/human distinction. Alt-C avoids impersonating Steve, but still creates a
synthetic principal that hides the real accountable human.

### Scenario 3 — later audit of a regression

Carol later audits a regression. Alt-A makes it harder to tell which local
human actually made the change because every local Codex commit collapses into
Steve. Alt-B preserves the real human author while still distinguishing local
human work from `ppx` bot work. Alt-C distinguishes local work from bot work
but still obscures the actual human operator.

### Scenario 4 — interaction with existing bot identity rules

The repo already has a dedicated `ppx` bot identity and branch-prefix rule.
Alt-A unnecessarily reuses Steve's identity even when a different local human is
operating. Alt-B leaves the bot rule untouched and cleanly separates local human
authorship from bot authorship. Alt-C adds a third synthetic identity with no
clear benefit over recording the real local operator.

## Conclusions

- Alt-A is rejected because it falsifies authorship whenever the local operator
  is not Steve.
- Alt-C is rejected because it introduces another synthetic actor while hiding
  the real local accountable human.
- Alt-B survives because it preserves accurate local authorship and still leaves
  the `ppx` bot identity rule intact.

## Recommendation

Adopt Alt-B. Local Codex work should use the actual configured local git author
identity in the clone. The overlay should explicitly forbid using the `ppx` bot
identity for local work unless the work is actually being performed under the
bot workflow.

## Implications for open work

- `AGENTS-codex.md` should be updated to stop requiring Steve impersonation and
  instead require a real configured local git identity.
- Existing bot identity and branch-prefix decisions remain unchanged.

## Decision status

`decided` — locked by `DI-nufit` in
`protocols/wire-lab.d/TODO/TODO-tapaf-codex-local-author-identity.md`.

