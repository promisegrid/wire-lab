# TE-18: Bot identity and branch prefix

*Thought experiment, part of the [PromiseGrid Wire Lab](../../specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260429-165101

(First drafted 2026-04-29 16:51:01 UTC, retrospective-style: this TE was authored after the corresponding DI was already locked, in response to the conditions on `review-20260429-162212-steve-traugott.md` for `ppx/dr-001-bootstrap`.)

## Status

decided

## Decision under test

When an LLM-driven agent operating on a human's behalf participates in a mob-consensus-style repo, what identity should the agent commit under, and what branch prefix should its proposal branches use?

## Setup

- The repo is `github.com/promisegrid/wire-lab`. The human principal is Steve Traugott (`stevegt@t7a.org`).
- The agent is Perplexity Computer, an LLM-driven assistant that pushes proposal branches the human reviews.
- Mob-consensus convention (see `https://github.com/stevegt/mob-consensus`) derives a branch prefix from the committer's email local-part: `<local-part>/{twig}`.
- PromiseGrid is the long-term replacement for GitHub. Anything that ties identity or workflow to a forge-specific feature is migration debt.

## Alternatives

- **Alt-A (literal mob-consensus):** Bot identity is `stevegt+ppx@t7a.org (Perplexity Computer)`. Branch prefix is the literal local-part: `stevegt+ppx/{twig}`. **Easier:** maximally consistent with the documented mob-consensus convention. **Harder:** the literal `stevegt+ppx/{twig}` prefix is verbose and visually noisy. **New obligation:** none beyond mob-consensus.

- **Alt-B (short prefix):** Bot identity is `stevegt+ppx@t7a.org (stevegt-via-perplexity)`. Branch prefix is shortened to `ppx/{twig}` (a concession against literal mob-consensus). **Easier:** branch names are short and memorable; subaddressing still routes the bot's mail to Steve's inbox; the parenthetical name still attributes the actor unambiguously. **Harder:** small departure from mob-consensus literal convention; future reviewers must read `AGENTS-ppx.md` and `DR-001` to learn that `ppx/` means Perplexity-on-behalf-of-Steve.

- **Alt-C (separate principal identity):** Bot identity is `perplexity-computer@example` (a wholly separate principal). **Easier:** strong separation between bot and human. **Harder:** the bot is *not* a separate principal — it acts on Steve's behalf with Steve's accountability. A separate identity invites the misimpression that the bot has independent authority, which contradicts the agent's actual promise (it commits to acting under Steve's review).

- **Alt-D (no distinct identity):** Bot commits as `stevegt@t7a.org (Steve Traugott)`. **Easier:** simplest. **Harder:** loses the audit trail of which commits were drafted by the agent versus Steve. When Steve later reviews accountability for a regression, the distinction matters.

## Selection

Alt-B. Bot identity is `stevegt+ppx@t7a.org (stevegt-via-perplexity)`; branch prefix is `ppx/{twig}`. This was selected by Steve in chat on 2026-04-28; the locked decision is recorded as `DI-001-20260428-195700` in `protocols/wire-lab.d/TODO/TODO-20260429-030146-perplexity-computer-onboarding.md`, and the supporting DR is `DR/DR-001-20260428-195700-bot-identity.md`.

## Decision status

`decided` — DI-001-20260428-195700 is already active. This TE is retrospective and exists to satisfy the TE-protocol requirement that non-trivial decisions have a corresponding TE artifact.
