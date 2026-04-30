# Thought Experiments

Each thought experiment is a falsifiable mental run of a Wire Lab design choice. Each lives in its own file, addressable by content hash (its pCID).

## Filename convention

```
TE-YYYYMMDD-HHMMSS-some-phrase.md
```

The timestamp is the moment the TE was first drafted (or, for TEs that pre-date this convention, the revision of `specs/harness-spec-draft.md` in which they first appeared). The slug is a kebab-case rendering of the TE's title. Files are not renamed when the experiment is later refined — the timestamp pins origin, not last-edited.

## Index

| # | First drafted | Title |
|---|---|---|
| TE-1 | 2026-04-27 18:00 | [Promise-stack ordering](TE-20260427-180000-promise-stack-ordering.md) |
| TE-2 | 2026-04-27 18:01 | [Trust-ledger merge after partition](TE-20260427-180100-trust-ledger-merge-after-partition.md) |
| TE-3 | 2026-04-27 18:02 | [Currency exchange-rate equilibration](TE-20260427-180200-currency-exchange-rate-equilibration.md) |
| TE-4 | 2026-04-27 18:03 | [Sybil under double auction](TE-20260427-180300-sybil-under-double-auction.md) |
| TE-5 | 2026-04-27 18:04 | [Kernel-as-handler vs. classical kernel](TE-20260427-180400-kernel-as-handler-vs-classical-kernel.md) |
| TE-6 | 2026-04-27 18:05 | [Capability-token revocation propagation](TE-20260427-180500-capability-token-revocation-propagation.md) |
| TE-7 | 2026-04-27 18:06 | [Human-novice onboarding under K4](TE-20260427-180600-human-novice-onboarding-under-k4.md) |
| TE-8 | 2026-04-27 18:07 | [Generational handoff](TE-20260427-180700-generational-handoff.md) |
| TE-9 | 2026-04-27 18:08 | [Two communities, two pCIDs, same intent](TE-20260427-180800-two-communities-two-pcids-same-intent.md) |
| TE-10 | 2026-04-27 18:09 | [Slow-mover survival](TE-20260427-180900-slow-mover-survival.md) |
| TE-11 | 2026-04-27 18:10 | [Ostrom's principles audit](TE-20260427-181000-ostroms-principles-audit.md) |
| TE-12 | 2026-04-27 18:11 | [Promise-stack as zero-knowledge envelope](TE-20260427-181100-promise-stack-as-zero-knowledge-envelope.md) |
| TE-13 | 2026-04-27 18:12 | [Time-traveling break-witness](TE-20260427-181200-time-traveling-break-witness.md) |
| TE-14 | 2026-04-28 08:00 | [A harness-spec change walks through the unified flow](TE-20260428-080000-harness-spec-change-walks-through-unified-flow.md) |
| TE-15 | 2026-04-28 09:45 | [Should this design become `promisegrid/promisegrid/README.md`?](TE-20260428-094500-should-this-design-become-promisegrid-readme.md) |
| TE-16 | 2026-04-29 03:32 | [Durable review feedback as contest artifact](TE-20260429-033208-review-feedback-as-contest-artifact.md) |
| TE-17 | 2026-04-29 16:22 | [Review reply as promise](TE-20260429-162212-review-reply-as-promise.md) |
| TE-18 | 2026-04-29 16:51 | [Bot identity and branch prefix](TE-20260429-165101-bot-identity-and-branch-prefix.md) |
| TE-19 | 2026-04-29 16:51 | [Branch-protection posture for `main`](TE-20260429-165102-branch-protection-posture.md) |
| TE-20 | 2026-04-29 16:51 | [Bot review style](TE-20260429-165103-bot-review-style.md) |
| TE-21 | 2026-04-29 17:35 | [Spec doc as promise](TE-20260429-173520-spec-doc-as-promise.md) |
| TE-22 | 2026-04-29 17:55 | [Spec-doc store layout and pCID machinery](TE-20260429-175530-spec-doc-store-and-pcid-machinery.md) |
| TE-23 | 2026-04-30 06:43 | [Congruence/convergence duality and pCID framing](TE-20260430-064307-congruence-convergence-duality-and-pcid-framing.md) |
| TE-25 | 2026-04-30 21:34 | [TE-21 numbering collision and harness-spec path](TE-20260430-213447-te-numbering-collision-and-harness-spec-path.md) |

The TE numbers (TE-1, TE-2, …) are stable identifiers used in the harness-spec; the timestamp slug is what makes the file content-addressable and chronologically sortable on disk.

## Adding a new TE

1. Decide the title.
2. Pick a UTC timestamp — typically `date -u +%Y%m%d-%H%M%S`.
3. Render the title to kebab-case for the slug.
4. Create `TE-YYYYMMDD-HHMMSS-slug.md` in this directory.
5. Add a one-line summary to `../../specs/harness-spec-draft.md` §8 with a link.
6. Add the row to this index.
7. Open a PR.
