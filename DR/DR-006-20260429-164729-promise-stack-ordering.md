# DR-006 - Promise-stack ordering (TE-famar)

DR-ID: DR-006-20260429-164729
Date: 2026-04-29 16:47:29
Asked by: stevegt@t7a.org (Steve Traugott)
State: closed
Question: Historical question: given a stack of promise frames in a single message, what is the canonical evaluation order, what is the canonical wire-encoding order, and what is the placement convention for different assertion types?
Why this blocks progress: No longer blocks progress as a standalone decision. TE-havib DF-36.2 Alt-2.A revised retired promise-stack as a separate hypothesis before the DF-1.* queue landed; payload recursion under `grid <pcid> <payload>` is the active framing.
Affects: Historical: `protocols/wire-lab.d/specs/harness-spec-draft.md` §1.1 (Promise frame shape); `protocols/wire-lab.d/specs/harness-spec-draft.md` §2 (trust ledger, per-assertion-type); future TE-fijub (capability-token revocation propagation); TE-muvuv (zero-knowledge envelopes). Current forward pointer: TE-havib DF-36.2, TE-lozip, and `docs/essays/congruence-convergence-and-the-grid.md` §3.1.
Unblocks: Closed by TODO-kugod DI-runuh. No promise-stack ordering DIs are emitted; downstream work should use payload-recursion/grid-envelope framing.
Waiting on: DI-runuh

## Candidate alternatives considered

The full alternative analysis (Alt-A through Alt-E) lives in `docs/thought-experiments/TE-famar-promise-stack-ordering.md`. The recommended set is (1.1.a, 1.2.c, 1.3.a, 1.4.d):

- (a) Alt-E hybrid: peeling order outermost-first, plus `Project` available to receivers, plus per-assertion-type position-convention declared in the assertion-type spec (the pCID).
- (c) Criticality flag is a hybrid: assertion-type spec declares default; per-frame override is allowed.
- (a) Wire encoding writes outermost frame first.
- (d) Position convention is per-assertion-type — each assertion-type spec declares whether its position is normative.

## Decision

Closed as superseded/no-longer-applicable on 2026-05-08. TE-havib DF-36.2
Alt-2.A revised retired promise-stack as a separate hypothesis, so the
DF-1.* promise-stack ordering questions should not be answered and should not
emit DIs. TODO-kugod DI-runuh records the closure cascade.

## Linked DI

DI-runuh in `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`.

## Related commits

- `0ec32d9` Expand TE-famar (Promise-stack ordering) into full scenario form
- `e0c00ff` Address TE-famar review conditions: add DR-006 and TODO-rivuk (DR-006 first authored)
- `5990e24` Merge ppx/merge-all-20260429-164729 (Converge branch reviews)

## Last updated

2026-05-08 23:41:08

## Resolution event — 2026-05-08

TODO-kugod implemented the Cat-3 promise-stack retirement cascade. TE-famar,
TE-muvuv, and TE-robub now carry forward-pointer refinements; TODO-rivuk is
closed as superseded; this DR is closed as no-longer-applicable rather than
decided. The original candidate alternatives above remain historical evidence
for why promise-stack ordering was considered.
