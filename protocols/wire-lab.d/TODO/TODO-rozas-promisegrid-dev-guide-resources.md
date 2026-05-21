# TODO-rozas: PromiseGrid dev-guide writer resources

## Prior aliases

This TODO is being filed after the TE-mumuv proquint migration locked
2026-05-07, so it is minted directly under `TODO-rozas`. No prior
integer or timestamp alias.

## Status

Running. The resource map is implemented, but `DEV-GUIDE-RESOURCES.md`
remains a periodically maintained writer-resource file and `rozas.10`
stays open for the future post-stabilization review.

## Decision Intent Log

### DI-nunut

- ID: DI-nunut
- Date: 2026-05-08 17:15:32
- Status: active
- Author: stevegt@t7a.org (Steve Traugott)
- Decision: Wire-lab will provide a future top-level `DEV-GUIDE-RESOURCES.md` as a PromiseGrid-facing source map plus short writer notes for PromiseGrid Development Guide authors. The guide is about PromiseGrid, not wire-lab; wire-lab is experimental simulation evidence and design provenance. After the PromiseGrid Development Guide stabilizes, stable guide prose becomes the higher-level developer source of truth, except where the guide cites hashed, frozen protocol specs by pCID.
- Intent: Give human and LLM guide writers enough curated context to write accurate PromiseGrid guide prose while preventing experimental wire-lab mechanics from being mistaken for product commitments or final developer APIs.
- Constraints:
    - Do not edit `/home/stevegt/lab/promisegrid-dev-guide` under this TODO until a later task explicitly asks for cross-repo guide changes.
    - Keep `DEV-GUIDE-RESOURCES.md` focused on source mapping and writer notes; do not duplicate full guide prose or large excerpts.
    - Organize resources around the guide's locked audiences: Laypeople, App Devs, and Kernel Devs.
    - Preserve pCID authority: a frozen protocol spec cited by pCID remains authoritative for that protocol regardless of whether it is hosted in wire-lab or another repo.
    - Treat exploratory TEs, TODOs, and DRs as provenance once the guide has stabilized corresponding PromiseGrid claims.
- Affects: `docs/thought-experiments/TE-david-promisegrid-dev-guide-resources.md`; `docs/thought-experiments/README.md`; `protocols/wire-lab.d/TODO/TODO-rozas-promisegrid-dev-guide-resources.md`; `protocols/wire-lab.d/TODO/TODO.md`; future `DEV-GUIDE-RESOURCES.md`; future `README.md` pointer.

### DI-zalak

- ID: DI-zalak
- Date: 2026-05-08 20:32:35
- Status: active
- Author: stevegt@t7a.org (Steve Traugott)
- Decision: Wire-lab will answer PromiseGrid Development Guide feedback with an audience readiness matrix and likely normative citation paths in `DEV-GUIDE-RESOURCES.md`. The matrix may identify current settled design commitments for guide prose, but it must keep final app-developer APIs and kernel porting boundaries open where no frozen spec or DR/DI closure exists.
- Intent: Give guide writers a compact current answer to feedback items `FB-gigit`, `FB-rivot`, `FB-vitih`, `FB-mulaj`, and `FB-rigod` without over-promising unstable PromiseGrid APIs or treating wire-lab provenance as final guide authority.
- Constraints:
    - Do not modify `/home/stevegt/lab/promisegrid-dev-guide`.
    - Keep `DR-napum`, `DR-tuhaz`, and `DR-davod` open unless a future DI settles their final PromiseGrid decisions.
    - Use `DEV-GUIDE-RESOURCES.md` for writer guidance, not full guide prose.
    - Treat frozen pCID specs and implementation conformance records as likely future normative citations; treat TEs, TODOs, DRs, and harness docs as provenance unless the guide later cites them otherwise.
- Affects: `DEV-GUIDE-RESOURCES.md`; `DR/DR-napum-promisegrid-layperson-guide-claims.md`; `DR/DR-tuhaz-promisegrid-app-dev-guide-contract.md`; `DR/DR-davod-promisegrid-kernel-dev-porting-boundary.md`; `protocols/wire-lab.d/TODO/TODO-rozas-promisegrid-dev-guide-resources.md`.

### DI-kizav

ID: DI-kizav
Date: 2026-05-09 10:58:02
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Keep TODO-rozas running while `DEV-GUIDE-RESOURCES.md` remains under periodic maintenance and while `rozas.10` waits for PromiseGrid Development Guide prose to stabilize.
Intent: Prevent the TODO status from implying full closure when the file has both an ongoing AGENTS.md maintenance rule and a future post-stabilization review checkpoint.
Constraints: Keep `AGENTS.md` as the recurring instruction source for when `DEV-GUIDE-RESOURCES.md` must be updated. Keep `rozas.10` open as a future milestone rather than copying it into `AGENTS.md`. Do not edit `/home/stevegt/lab/promisegrid-dev-guide` under this status correction.
Affects: `protocols/wire-lab.d/TODO/TODO-rozas-promisegrid-dev-guide-resources.md`; `AGENTS.md`; `DEV-GUIDE-RESOURCES.md`.

### DI-ragaz

ID: DI-ragaz
Date: 2026-05-20 17:49:46
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Route the remaining open PromiseGrid Development Guide feedback items
from `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md` into six
wire-lab question-home simulations plus root scenarios that can be run against
candidate designs.
Intent: The guide feedback now asks concrete App Dev and Kernel Dev questions
that should not live only in the external guide repo or in broad DR prose.
Simulation homes keep each question cluster independently evolvable, while root
scenarios make the pressure reusable across existing and future sims. The work
does not edit the guide repo and does not settle the blocking DRs.
Constraints: Keep the new sims as provisional question homes, not guide prose and
not PromiseGrid API commitments. Use `TODO-rozas` as the owner and
`DEV-GUIDE-RESOURCES.md` as the guide-writer map. `FB-gigit` remains covered by
the existing layperson readiness matrix and `DR-napum`; no new sim is required
for that narrative-status question. Do not mark `DR-tuhaz`, `DR-davod`,
`DR-tumus`, `DR-gabif`, `DR-robon`, or `DR-napum` closed.
Affects: `simulations/SIM-mikas-minimal-blob-app-contract/`;
`simulations/SIM-robot-app-semantics-conformance/`;
`simulations/SIM-zisan-device-bound-agent/`;
`simulations/SIM-kugap-live-sync-audit-split/`;
`simulations/SIM-nijuz-multi-embodiment-identity/`;
`simulations/SIM-funas-kernel-porting-boundary/`;
`scenarios/minimal-immutable-blob-app/`;
`scenarios/app-semantics-partial-conformance/`;
`scenarios/device-bound-agent-physical-effect/`;
`scenarios/live-crdt-audit-publication/`;
`scenarios/multi-embodiment-app-identity/`;
`scenarios/portable-signing-key-identity/`;
`scenarios/kernel-porting-boundary/`; `simulations/README.md`;
`DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-rozas-promisegrid-dev-guide-resources.md`.

## Context

The PromiseGrid Development Guide currently has three locked audiences:
Laypeople, App Devs, and Kernel Devs. Wire-lab needs to make its
PromiseGrid-relevant evidence easy to find, but it must not become the
guide itself. The guide should explain PromiseGrid. Wire-lab should explain
which experiments, DIs, DRs, protocol specs, and TODOs support or block the
guide's claims.

## Related DRs

- `DR/DR-napum-promisegrid-layperson-guide-claims.md` blocks settled Laypeople / Intro and Goals prose.
- `DR/DR-tuhaz-promisegrid-app-dev-guide-contract.md` blocks settled App Dev / How to write a grid app prose.
- `DR/DR-davod-promisegrid-kernel-dev-porting-boundary.md` blocks settled Kernel Dev / How to port the infrastructure prose.

## Feedback responses

- `FB-gigit` is answered at current wire-lab scope by the Laypeople row in `DEV-GUIDE-RESOURCES.md` and the 2026-05-08 `DR-napum` event.
- `FB-rivot` is answered at current wire-lab scope by the App Devs row in `DEV-GUIDE-RESOURCES.md` and the 2026-05-08 `DR-tuhaz` event.
- `FB-vitih` is answered at current wire-lab scope by the Kernel Devs row in `DEV-GUIDE-RESOURCES.md` and the 2026-05-08 `DR-davod` event.
- `FB-mulaj` is answered by the new audience readiness matrix in `DEV-GUIDE-RESOURCES.md`.
- `FB-rigod` is answered by the likely normative citation path in `DEV-GUIDE-RESOURCES.md`.
- `FB-vopik` is routed to `SIM-mikas` and
  `scenarios/minimal-immutable-blob-app/` for minimal CAS-facing app pressure.
- `FB-dodos`, `FB-hisis`, `FB-kutub`, `FB-gomod`, and `FB-tahof` are routed to
  `SIM-robot` and `scenarios/app-semantics-partial-conformance/` for
  app-vocabulary, local/wire identity, partial-conformance, provisional-signing,
  ingress, and policy/economy pressure.
- `FB-nojit`, `FB-tisuf`, and `FB-tulit` are routed to `SIM-zisan` and
  `scenarios/device-bound-agent-physical-effect/` for host-owned device agents,
  physical effects, owner/operator authority, and driver-stack conformance.
- `FB-hurit` and `FB-nilat` are routed to `SIM-kugap` and
  `scenarios/live-crdt-audit-publication/` for live-state versus durable-audit
  split pressure.
- `FB-zazon` and `FB-robif` are routed to `SIM-nijuz`,
  `scenarios/multi-embodiment-app-identity/`, and
  `scenarios/portable-signing-key-identity/` for multi-host app identity,
  per-component conformance claims, and portable signing-key recipes.
- `FB-vitih`, `FB-mulum`, and `FB-potin` are routed to `SIM-funas` and
  `scenarios/kernel-porting-boundary/` for kernel/runtime/dispatcher porting
  boundary pressure.

## Subtasks

- [x] rozas.1 Run TE-david with Alice, Bob, and Carol to test what guide writers need to see.
- [x] rozas.2 Lock the writer-resource decision in `DI-nunut`.
- [x] rozas.3 Create top-level `DEV-GUIDE-RESOURCES.md` as a source map plus writer notes, not guide prose.
- [x] rozas.4 Add a concise top-level `README.md` pointer to `DEV-GUIDE-RESOURCES.md`.
- [x] rozas.5 Map Laypeople guide needs to current wire-lab sources, stable claims, unsettled claims, and provenance.
- [x] rozas.6 Map App Dev guide needs to current wire-lab sources, stable claims, unsettled claims, and provenance.
- [x] rozas.7 Map Kernel Dev guide needs to current wire-lab sources, stable claims, unsettled claims, and provenance.
- [x] rozas.8 Add authority-transition notes explaining when stabilized guide prose supersedes exploratory wire-lab notes.
- [x] rozas.9 Identify unresolved guide-writing gaps and file or link DRs for any question that blocks a settled guide claim.
- [ ] rozas.10 After guide prose stabilizes, review `DEV-GUIDE-RESOURCES.md` and downgrade superseded wire-lab material to provenance pointers.
- [x] rozas.11 Route current open guide-feedback questions into provisional
  simulation homes and root scenarios so the guide can cite evidence pressure
  without treating open DRs as settled. Source: `DI-ragaz`.
