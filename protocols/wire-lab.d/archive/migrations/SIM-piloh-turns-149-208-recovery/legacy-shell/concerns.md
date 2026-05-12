# Concerns Matrix

This matrix keeps the recovery goal visible while the simulation fills in more
evidence. It is intentionally about PromiseGrid design recovery, not about
teaching developers to implement wire-lab's temporary directory layout.

| Concern | Simulation evidence | Expected result |
|---|---|---|
| Turns 149-208 recovery must stay open until every source turn is accounted for. | `protocols/wire-lab.d/TODO/TODO-jivam-turns-149-170-recovery-completion.md`; this simulation's future observations. | Recovery closure evidence or successor TODOs before TODO-jivam closes. |
| Apparatus and specimens must not be confused. | `protocols/wire-lab.d/` remains rooted; candidate protocols move under `protocols/` inside this simulation. | Guide writers and agents can distinguish harness apparatus from PromiseGrid candidates. |
| Legacy proposal records are no longer the active guide-feedback mechanism. | `archive/proposals/`; `DEV-GUIDE-RESOURCES.md`; `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md` as external process context. | Dev-guide writers find current resources without treating proposal records as live workflow. |
| Transport evidence must remain verifiable. | `world/transports/wire-lab-devs-draft/`; `seed/wire-lab-devs-draft-migration.md`. | Message filenames still match raw CIDv1 hashes after migration. |
| Simulation evidence must graduate through provenance. | `decisions.md`; later DR/DI/spec/dev-guide records. | No simulation result silently becomes an authoritative PromiseGrid spec. |

Source: `DI-fakin`.
