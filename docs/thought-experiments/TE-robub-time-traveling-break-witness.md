# TE-robub: Time-traveling break-witness

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TE-13` (integer alias)
- `TE-20260427-181200` (timestamp alias and pre-migration filename)

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*


## Status

stub
Charlie discovers in year 5 that an event in year 2 was actually fraudulent. He emits a break-witness about a long-settled transaction. Does the network accept it? Does it re-litigate the past? Outcome: a policy on how far back break-witnesses can reach and how settlement becomes "final enough."

## Refinements

### 2026-05-08 — Reframe future break-witness work under payload recursion

TE-havib DF-36.2 Alt-2.A revised retired promise-stack as a separate
hypothesis. Future break-witness work should not assume break-witnesses are
outer promise-stack frames with special ordering semantics; it should ask how
a protocol identified by a pCID expresses a later promise or assessment about
an earlier payload, where both records may be nested under
`grid <pcid> <payload>`.

TODO-kugod records this Cat-3 cascade in DI-runuh. This stub remains a
historical placeholder for the time-traveling break-witness question, but its
active framing is now payload-recursion/grid-envelope work, not
promise-stack-envelope work.
