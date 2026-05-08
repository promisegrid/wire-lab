# TODO-kugod: TE-40 apparatus-vs-specimen completion + TE-famar closure

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-25` (integer alias)
- `TODO-20260507-002306` (timestamp alias and pre-migration filename)

## Status

Closed 2026-05-08. Cat-3 cascade work from TE-havib promise-stack
retirement has been applied to TE-famar, TE-muvuv, and TE-robub.
TODO-rivuk and DR-006 now point readers at the TE-havib DF-36.2
retirement instead of inviting stale promise-stack DF answers.

## Threads absorbed from OPEN-THREADS.md

### T-PROMSTACK-RETIRE-CASCADE (formerly OPEN-THREADS, opened 2026-05-06)

Cat-3 cascade from TE-havib promise-stack retirement. DF-36.2 Alt-2.A
revised: promise-stack retired as a separate hypothesis; payload-
recursion under per-protocol specs is the answer per TE-lozip § 3.1 +
framing essay § 3.1.

Scope: refine TE-famar, TE-muvuv, TE-robub with Cat-3 entries recasting
promise-stack vocabulary into payload-recursion vocabulary; verify no
other TE in the corpus references promise-stack as a separate layer;
confirm protocols/promise-stack.d/ remains absent from the tree.

Blocking: nothing today; pure refinement work.

Anchor: TE-havib § DF-36.2 Alt-2.A (revised); TE-lozip; framing essay § 3.1.
Disposition-file pointer: `dropped-thread-disposition-20260506.md`
§ TE-40 cluster (18 UTs).

## Question log

(Per AGENTS-ppx Question-logging discipline. No questions logged yet.)

## Decision Intent Log

ID: DI-runuh
Date: 2026-05-08 23:41:08
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Execute TODO-kugod as a narrow Cat-3/Cat-4 retirement cascade. Append forward pointers to TE-famar, TE-muvuv, and TE-robub; close TODO-rivuk and DR-006 as superseded/no-longer-applicable; do not create `protocols/promise-stack.d/`; do not sweep `protocols/wire-lab.d/specs/harness-spec-draft.md`; do not rewrite TE historical body text or TE status fields in this pass.
Intent: Preserve the historical promise-stack analyses while making the current corpus navigable after TE-havib DF-36.2 retired promise-stack as a separate hypothesis. Readers should follow payload-recursion under `grid <pcid> <payload>` via TE-havib, TE-lozip, and the congruence/convergence essay instead of answering stale TE-famar DF-1.* questions.
Constraints: This is a navigational cleanup only. Cat-3 refinements append to `## Refinements`; historical TE bodies remain untouched under TE-dabol/TE-vudaf Cat-1b/Cat-3 rules. DR-006 and TODO-rivuk may be updated to stop stale work queues, but no promise-stack ordering DIs are created.
Affects: `docs/thought-experiments/TE-famar-promise-stack-ordering.md`; `docs/thought-experiments/TE-muvuv-promise-stack-as-zero-knowledge-envelope.md`; `docs/thought-experiments/TE-robub-time-traveling-break-witness.md`; `protocols/wire-lab.d/TODO/TODO-rivuk-te-promise-stack-ordering.md`; `DR/DR-006-20260429-164729-promise-stack-ordering.md`; `protocols/wire-lab.d/TODO/TODO-kugod-te-40-apparatus-vs-specimen-completion.md`.
