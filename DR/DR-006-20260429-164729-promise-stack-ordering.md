# DR-006 - Promise-stack ordering (TE-1)

DR-ID: DR-006-20260429-164729
Date: 2026-04-29 16:47:29
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: Given a stack of promise frames in a single message, what is the canonical evaluation order, what is the canonical wire-encoding order, and what is the placement convention for different assertion types?
Why this blocks progress: Senders, intermediate routers, and receivers all need a shared expectation about how a multi-frame `Promise` stack is read and acted upon. Without a shared convention, a v1 receiver and a v2 receiver can ratify the same wire bytes differently, partial-knowledge peers can act on stale assertions, and revocation/compliance frames added at one layer can be ignored by another. Locking a wire-and-evaluation convention is a prerequisite for harness conformance tests on §1.1 of `harness-spec.md`.
Affects: `harness-spec.md` §1.1 (Promise frame shape); `harness-spec.md` §2 (trust ledger, per-assertion-type); future TE-6 (capability-token revocation propagation); future TE-12 (zero-knowledge envelopes).
Unblocks: `TODO/005-te-promise-stack-ordering.md` subtasks; `harness-spec.md` §1.1 conformance tests; criticality-flag location decision; per-assertion-type position-convention authority.
Waiting on: stevegt@t7a.org (Steve Traugott) for DF answers DF-1.1 through DF-1.4 in `docs/thought-experiments/TE-20260427-180000-promise-stack-ordering.md`.

## Candidate alternatives considered

The full alternative analysis (Alt-A through Alt-E) lives in `docs/thought-experiments/TE-20260427-180000-promise-stack-ordering.md`. The recommended set is (1.1.a, 1.2.c, 1.3.a, 1.4.d):

- (a) Alt-E hybrid: peeling order outermost-first, plus `Project` available to receivers, plus per-assertion-type position-convention declared in the assertion-type spec (the pCID).
- (c) Criticality flag is a hybrid: assertion-type spec declares default; per-frame override is allowed.
- (a) Wire encoding writes outermost frame first.
- (d) Position convention is per-assertion-type — each assertion-type spec declares whether its position is normative.

## Decision

Pending. Awaiting DF answers from Steve. Once DF lands, this DR transitions to `decided` and emits DI entries into `TODO/005-te-promise-stack-ordering.md`.

## Linked DI

To be created in `TODO/005-te-promise-stack-ordering.md` after DF.

## Related commits

- `0ec32d9` Expand TE-1 (Promise-stack ordering) into full scenario form
- DR-006 first authored on `ppx/merge-all-20260429-164729`
