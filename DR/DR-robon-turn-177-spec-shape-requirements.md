# DR-robon - Turn-177 spec-shape requirements

DR-ID: DR-robon
Date: 2026-05-17 09:28:02
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: Should PromiseGrid protocol specs require layer-specific promise-vocabulary sections, 100-year pressure-test sections, and layperson-readable mental-model / easy-implementation summaries, and if so which specs and layers must include them?
Why this blocks progress: Turn 177 promoted promise economy, 100-year survival, easy mental models, and easy implementations from background motivation into design constraints. TODO-kulih / TE-nibar owns spec-doc shape, but the turn-177 obligations are not yet locked as required spec sections, recommended guide prose, or deferred notes.
Affects: `protocols/wire-lab.d/TODO/TODO-kulih-te-spec-doc-as-promise.md`; `simulations/SIM-rusap-promise-accounting-records/`; `docs/thought-experiments/TE-nibar-spec-doc-as-promise.md`; `docs/thought-experiments/TE-dajot-100-year-goal-as-design-constraint.md`; `DEV-GUIDE-RESOURCES.md`.
Unblocks: TODO-kulih 010.9; future protocol-spec templates; guide-writer treatment of promise-accounting and 100-year claims.
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision:
Linked DI: DI-navod; DI-pator; DI-davov
Related commits:
Last updated: 2026-05-17 09:28:02

## Event log

- 2026-05-17 09:28:02 — Opened during turn-177 cleanup so spec-shape requirements have an explicit DR rather than remaining embedded in TODO-kulih 010.9 and simulation scenarios.

## Evidence

- Turn 177 made the promise economy foundational at every layer and tied that to the 100-year goal, easy mental models, and easy implementations.
- `simulations/SIM-rusap-promise-accounting-records/SCENARIOS.md` records concrete pressure cases: kept promises, refusal, corruption, cross-layer decisions, sparse retention, identity rotation, and layperson explanation.
- `TODO-kulih` owns TE-nibar's spec-doc-as-promise decision path and now includes 010.9 for turn-177 promise-vocabulary / 100-year / mental-model obligations.
- `DEV-GUIDE-RESOURCES.md` already warns that promise accounting is peer-local and not a central or harness-owned ledger.

## Candidate decisions

- **Alt-A: required sections for every protocol spec.** Every protocol spec must include promise vocabulary, 100-year pressure test, and layperson/easy-implementation summary sections.
- **Alt-B: required when materially relevant.** Specs must include these sections when they define promises, peer behavior, retention, replication, identity, or cross-layer decisions; purely mechanical specs may cite an explicit non-applicability note.
- **Alt-C: guide-level only for now.** Keep these as guide-writer obligations and simulation pressure, but do not require every protocol spec to carry them until TE-nibar locks broader spec-doc shape.

## Notes

This DR does not settle TE-nibar. It narrows the turn-177 addition that TODO-kulih 010.9 must answer.
