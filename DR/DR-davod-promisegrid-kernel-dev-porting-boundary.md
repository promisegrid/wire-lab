# DR-davod - PromiseGrid kernel developer porting boundary

DR-ID: DR-davod
Date: 2026-05-08 17:25:12
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: What stable PromiseGrid infrastructure boundary should the Development Guide present to kernel developers who are porting PromiseGrid to a new platform or language?
Why this blocks progress: Wire-lab has a strong apparatus/specimen split and layered substrate model, but the guide still needs an explicit developer-facing porting target: which frozen specs, runtime expectations, conformance claims, and implementation obligations a port must satisfy.
Affects: `DEV-GUIDE-RESOURCES.md`; `/home/stevegt/lab/promisegrid-dev-guide/README.md`; `/home/stevegt/lab/promisegrid-dev-guide/TODO/TODO-binap-readme-outline-lock.md`; `protocols/wire-lab.d/specs/harness-spec-draft.md`; `transports/README.md`; `implementations/README.md`
Unblocks: `TODO-rozas.7`; `TODO-binap.7`
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision:
Linked DI: DI-zalak (wire-lab guide-resource guidance; DR remains open)
Related commits:
Last updated: 2026-05-08 20:32:35

## Event log

- 2026-05-08 17:25:12 — Opened from `TODO-rozas.9` so kernel-developer guide resources can cite an explicit DR for unsettled porting boundaries.
- 2026-05-08 20:32:35 — Narrowed current wire-lab answer for `FB-vitih`: the porting target is not the wire-lab harness. A PromiseGrid port should be framed as a runtime/implementation that hosts pCID-selected protocol handlers, implements the frozen binding/session/message specs it claims, and records those claims in implementation conformance records. This does not close the DR; the first required frozen spec set, runtime expectations, and implementation obligations remain open.
