# TODO-binag: PromiseGrid kernel design resolution

## Prior aliases

This TODO was minted directly under the post-`TE-mumuv` proquint handle scheme.
No integer or timestamp alias exists.

## Status

Running. The resolution packet is being implemented in stages: the design note,
simulation evidence surface, scenario pressure, and canary focus file are
prepared first; `DR-davod` remains open until Steve runs the focused canary and
the resulting evidence is reviewed.

## Decision Intent Log

### DI-fidot

ID: DI-fidot
Date: 2026-05-25 13:22:27
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Resolve the current PromiseGrid kernel design packet by extending
`SIM-fovip` as the active evidence home, writing a standalone kernel
role/profile design note, running focused score evidence, and then deciding
`DR-davod`.
Intent: The kernel TEs have converged on a coherent model: kernel is a local
role/profile set, app/kernel operations are pCID-selected promise messages,
ports publish promises, assumptions, refusals, and evidence, and file-like UX is
a projection over promise/event history. The repo now needs one
implementation-ready resolution path instead of more scattered TE analysis.
Constraints: Do not create a successor sim unless `SIM-fovip` evidence fails. Do
not create another TE for this packet unless new alternatives emerge. Do not run
provider-backed canary/GA; Steve runs those commands. Do not rewrite scored
artifacts. Preserve Promise Theory vocabulary: local trust, autonomous
promisers, reciprocal promises, evidence, refusals, and make/break history.
Affects: `simulations/SIM-fovip-kernel-promise-boundary-port-contract/`;
`scenarios/kernel-porting-boundary/`;
`DR/DR-davod-promisegrid-kernel-dev-porting-boundary.md`;
`docs/research/DN-lujad-promisegrid-kernel-role-profile.md`;
`DEV-GUIDE-RESOURCES.md`; `protocols/wire-lab.d/specs/harness-spec-draft.md`;
kernel TEs under `docs/thought-experiments/`; `protocols/wire-lab.d/TODO/TODO.md`.

### DI-punuf

ID: DI-punuf
Date: 2026-06-05 19:46:28
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Record the non-monolithic kernel model as an explicit extension of
the existing role/profile synthesis. A PromiseGrid kernel is a local collection
of bounded promise-making roles, not necessarily one process. Those roles can be
implemented by one daemon, several local processes, cooperating objects, a WASM
host, MCU/library functions, or other runtime-specific profiles. POC12
`printer_port` is current executable evidence for a local hardware/resource
owner role that issues and redeems capability-promise tokens without making the
message kernel an authorization server.
Intent: The recent POC12 printer-port split makes the abstract role/profile
model concrete. Kernel prose should now say that dispatch, transport, lifecycle,
storage, compute, key/signing, hardware/resource access, and evidence recording
are role promises that may be carried by different local agents or objects,
while preserving the Promise Theory rule that no role commands another
autonomous agent.
Constraints: Do not create another TE for this unless a new architectural fork
appears. Do not freeze a final SDK, required pCID set, process topology, USB
API, token API, or production hardware contract. Do not describe these roles as
permission, authorization, service registry, policy enforcement, conformance, or
global trust authority. Keep POC12 as provisional evidence and keep `DR-davod`
open until the focused kernel evidence review is complete.
Affects: `docs/research/DN-lujad-promisegrid-kernel-role-profile.md`;
`DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-binag-promisegrid-kernel-design-resolution.md`.

## Basic principles to preserve

- Everything useful is a promise.
- Promise is the semantic unit; object/event/log/checkpoint is durable
  substrate; pCID-selected message is the boundary unit; file-like view is UX
  projection.
- PromiseGrid is a decentralized mainframe: agents get coherent file-like views
  over a shared promise/event hypergraph; trusted groups may maintain shared
  namespaces; coherence comes from reciprocal promises and local trust.
- No global namespace authority exists. An imposed global single-system image is
  rejected; voluntary group namespaces are valid only as promises inside a trust
  relationship.
- Local trust is the only trust. Every agent evaluates promises, evidence, and
  make/break history locally.
- No agent promises for another agent. A kernel, app, host, service, file, byte
  chunk, or peer promises only what it can plausibly keep or embody.
- Kernel is a role set, not a ruler. "Kernel" means local infrastructure roles
  that help an agent speak pCID-selected protocols, not a privileged authority.
- The stable interface is pCID discipline. A pCID names the protocol spec; the
  spec defines payload shape, canonical bytes, proof/signature encoding, and
  interpretation.
- App/kernel APIs are adapters over promises. Exposed app/kernel operations map
  to pCID-selected messages and evidence records.
- Ports publish promises, assumptions, refusals, and evidence policy. A credible
  port says what it supports, what it depends on, what it refuses, and what it
  records.
- Resources are promise-log projections. File/resource current state is a
  checkpoint over selected promise history; branches are different
  promise-history selections.
- Prior art is pressure, not authority. Borrow V, Amoeba, Plan 9, and Hurd
  patterns only after reframing them as local promises.

## Related kernel artifacts

- `DR/DR-davod-promisegrid-kernel-dev-porting-boundary.md` owns the stable
  kernel-developer porting-boundary decision.
- `docs/thought-experiments/TE-jimar-kernel-runtime-portability-boundary.md`
  narrows "kernel" toward role/profile definition.
- `docs/thought-experiments/TE-mazop-kernel-promise-boundary-and-minimum-port-contract.md`
  narrows the minimum credible kernel implementation promises.
- `docs/thought-experiments/TE-pudiv-app-kernel-grid-message-boundary.md`
  narrows app/kernel operations toward pCID-selected grid-message semantics.
- `docs/thought-experiments/TE-dunas-prior-art-influence-on-promisegrid-kernel.md`
  narrows safe prior-art influence.
- `docs/thought-experiments/TE-gakoh-local-views-over-promise-event-hypergraph.md`
  narrows the local-view, voluntary-namespace, CID-rooted-reference, and
  promise-log resource model.
- `simulations/SIM-fovip-kernel-promise-boundary-port-contract/` is the active
  evidence home.
- `simulations/SIM-funas-kernel-porting-boundary/` is the legacy question home.

## Subtasks

- [x] binag.1 Add this TODO file and cross-list it in
  `protocols/wire-lab.d/TODO/TODO.md`.
- [x] binag.2 Record `DI-fidot` in this TODO and cite it from the non-trivial
  design-document and evidence-surface changes made under this packet.
- [x] binag.3 Expand `SIM-fovip` so it becomes the active evidence home for all
  open kernel questions: kernel as role/profile set, app/kernel pCID grid
  messages, kernel implementation promise records, operation coverage, separated
  host/runtime assumptions, unsupported-pCID behavior, voluntary namespaces,
  CID-rooted references, and file/resource checkpoints over promise logs.
- [x] binag.4 Update the `SIM-fovip` simulation-local kernel-port draft to prefer
  "kernel implementation promises" in prose while preserving the historical
  filename and making clear the record is not a global certificate.
- [x] binag.5 Expand `scenarios/kernel-porting-boundary/` with Alice/Bob/Carol
  pressure for native, browser/WASM, mobile, MCU/header-only, split-service,
  broken-promise, voluntary-namespace, CID-rooted-reference, and checkpoint/log
  cases.
- [x] binag.6 Write
  `docs/research/DN-lujad-promisegrid-kernel-role-profile.md` as the
  plain-English design synthesis.
- [x] binag.7 Configure `/tmp/canary-cells` and provide staged canary commands
  for Steve to run.
- [ ] binag.8 Review the focused canary output after Steve runs it.
  Confirm `SIM-fovip` does not depend on central authority, global permission,
  or universal process shape; confirm open questions are implementation details,
  not blockers for `DR-davod`.
- [ ] binag.9 Decide `DR-davod` after focused evidence review. Set `State:
  decided` only if evidence passes, and keep guide handoff separate from the
  design decision.
- [ ] binag.10 Update TE statuses and refinements after the decision. Mark
  `TE-jimar`, `TE-mazop`, `TE-pudiv`, `TE-dunas`, and `TE-gakoh` decided by
  `DI-fidot`, and add a `TE-jikaf` refinement saying K1/K3 are profile cases
  under the role/profile model.
- [x] binag.11 Update guide-resource routing so guide writers can find
  `DN-lujad`, the active `SIM-fovip` evidence home, and the legacy `SIM-funas`
  question home.
- [x] binag.12 Update harness/spec summaries so the active kernel packet points
  at `DN-lujad` and remains evidence rather than final PromiseGrid API.
- [ ] binag.13 Validate after canary review and DR decision: search for stale
  unresolved `needs DF` language in kernel docs, run `git diff --check`, search
  for authority drift, and verify no deleted false-split reference IDs remain.
- [ ] binag.14 Commit with explicit file staging and an AGENTS-compliant body
  after final validation.
- [x] binag.15 Fold POC12 printer-port evidence into the non-monolithic
  kernel role/profile synthesis without freezing a final kernel API.

## Acceptance criteria

- `SIM-fovip` covers every open decision in `TE-jimar`, `TE-mazop`, `TE-pudiv`,
  `TE-dunas`, and `TE-gakoh`.
- `DN-lujad` gives a plain-English design synthesis suitable for future
  PromiseGrid guide writers.
- `DR-davod` can be decided without claiming a final SDK, final pCID set, or
  final implementation topology.
- No document presents daemon, microkernel, host runtime, library-only, or split
  mesh as the universal kernel shape.
- No document describes trust, permission, namespace, conformance, or policy as
  external authority.
- No provider-backed canary/GA is run by the agent.
- No scored artifacts are rewritten.
