# DR-lotir - Gridfile plan and machine event history

DR-ID: DR-lotir
Date: 2026-08-13 19:35:34 PDT
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: Should POC21 formally distinguish a human-readable Gridfile that
describes the intended finite plan from the append-only, parent-linked CAS
history of actions actually attempted and results actually observed, and what
stable term should name the executed history? How should that history represent
run-once start/completion identity, indeterminate interruption, and recurring
entrypoint invocations over the same finite plan?
Why this blocks progress: Existing POC21 prose uses "Gridfile journal" and
"ordered machine journal" for both an intended plan and executed history. Code,
protocol specifications, analyzers, and replay tools cannot define identities or
acceptance checks safely while those two objects remain conflated.
Affects: `protocols/wire-lab.d/TODO/TODO-kifok-poc21-grid-devops.md`;
`implementations/poc21-grid-devops/docs/DESIGN.md`;
`implementations/poc21-grid-devops/docs/discussion/grid-language-handoff-20260813.md`;
`docs/research/DN-gagog-grid-language-profiles-and-runtime-descriptors.md`;
`DEV-GUIDE-RESOURCES.md`.
Unblocks: POC21 Gridfile specification, machine-history schema, run-once action
identity, recurring entrypoint synchronization, replay behavior, and analyzer
terminology.
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision:
Linked DI:
Related commits:
Last updated: 2026-08-13 19:35:34 PDT

## Event log

- 2026-08-13 19:35:34 PDT — Opened while promoting the Grid-language session
  handoff into retained repository documentation. The discussion recommends
  `Gridfile` for the intended plan and either `machine event journal` or
  `machine timeline` for executed CAS history, but no DI has selected the final
  terminology or object contract.
- 2026-08-13 19:35:34 PDT — A dedicated TE is required before DF because the
  distinction affects object identity, replay semantics, run-once detection,
  diagnostics, and existing POC21 vocabulary.
- 2026-08-13 19:35:34 PDT — The retained discussion also proposes start and
  completion events joined by a stable invocation CID, an explicit indeterminate
  state after interruption, and repeated external entrypoint requests over the
  same locally adopted finite graph. `DR-junaz` separately owns the eventual
  Gridfile syntax; this DR owns the event-history and invocation identities.
