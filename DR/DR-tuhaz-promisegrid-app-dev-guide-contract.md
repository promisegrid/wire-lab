# DR-tuhaz - PromiseGrid app developer guide contract

DR-ID: DR-tuhaz
Date: 2026-05-08 17:25:12
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: What stable PromiseGrid contract, API, protocol subset, or provisional fallback should the Development Guide present to app developers as "how to write a grid app"?
Why this blocks progress: Current wire-lab evidence supports pCID-selected protocols, spec-doc identity, and implementation conformance claims, but it does not yet freeze an app SDK, handler ABI, or message API that guide prose can present as final.
Affects: `DEV-GUIDE-RESOURCES.md`; `/home/stevegt/lab/promisegrid-dev-guide/README.md`; `/home/stevegt/lab/promisegrid-dev-guide/TODO/TODO-binap-readme-outline-lock.md`; `implementations/README.md`; `protocols/*/specs/*.md`
Unblocks: `TODO-rozas.6`; `TODO-binap.6`
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision:
Linked DI: DI-zalak (wire-lab guide-resource guidance; DR remains open)
Related commits:
Last updated: 2026-05-08 20:32:35

## Event log

- 2026-05-08 17:25:12 — Opened from `TODO-rozas.9` so app-developer guide resources can distinguish current design direction from settled developer-facing contracts.
- 2026-05-08 20:32:35 — Narrowed current wire-lab answer for `FB-rivot`: the minimum stable app-developer contract is currently a discipline, not a frozen SDK or ABI. App code should target explicit protocol specs, use pCIDs once specs freeze, let the selected spec define payload and handler semantics, and publish implementation conformance claims. This does not close the DR; final app SDK, handler ABI, message API, or protocol subset remains open.
