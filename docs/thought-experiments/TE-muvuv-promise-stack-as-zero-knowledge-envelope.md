# TE-muvuv: Promise-stack as zero-knowledge envelope

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TE-12` (integer alias)
- `TE-20260427-181100` (timestamp alias and pre-migration filename)

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*


## Status

stub
Can a promise frame contain a ZK proof that "the inner payload satisfies predicate P" without revealing the payload? Pilot scenario: a regulator-style agent verifying compliance without seeing data. Outcome: a wire-level extension point or a "no, this needs out-of-band machinery" answer.

## Refinements

### 2026-05-08 — Reframe future ZK work under payload recursion

TE-havib DF-36.2 Alt-2.A revised retired promise-stack as a separate
hypothesis. Future zero-knowledge work should not assume ZK proofs are
promise-stack frames discovered by `Project`; it should ask how a protocol
identified by a pCID expresses "payload satisfies predicate P" when that
payload may itself be nested under `grid <pcid> <payload>`.

TODO-kugod records this Cat-3 cascade in DI-runuh. This stub remains a
historical placeholder for the ZK question, but its active framing is now
payload-recursion/grid-envelope work, not promise-stack-envelope work.
