# TODO-kifok: POC21 grid DevOps

## Status

Planned. Owns the POC21 DevOps track under
`implementations/poc21-grid-devops/`. POC21 applies the POC18 CAS/VCS substrate,
the POC19 stage0/stage1 bootstrap shape, and the POC20 CAS timeline/replay model
to pull-based machine administration. It is a separate POC21 track, not extra
scope inside POC20. Source: `DI-zosol`; `DI-tanov`; `DI-dahaj`; `DI-moson`.

## Decision Intent Log

ID: DI-zosol
Date: 2026-07-11 16:39:34 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Create POC21 as `implementations/poc21-grid-devops/` with a harness
TODO, local design document, narrow POC20 guardrail, and guide-resource update.
Intent: DevOps machine administration is the first concrete application of the
CAS/VCS, stage0 bootstrap, and timeline/replay work, but it would overload POC20
if folded into the semantic-model POC. A separate POC21 keeps POC20 focused on
timeline semantics while giving the DevOps work a production-shaped planning
home.
Constraints: POC21 must build on POC18, POC19, and POC20 rather than replacing
them; the initial implementation artifact is planning only; no host-root
mutation is allowed; future code must preserve PromiseGrid vocabulary, exact
`grid()` CBOR over TCP for inter-agent communication, local trust, and
CAS-backed replay.
Affects: `protocols/wire-lab.d/TODO/TODO-kifok-poc21-grid-devops.md`;
`protocols/wire-lab.d/TODO/TODO.md`;
`implementations/poc21-grid-devops/docs/DESIGN.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/DESIGN.md`;
`DEV-GUIDE-RESOURCES.md`.

ID: DI-tanov
Date: 2026-07-11 16:39:34 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: The first POC21 execution target is a root-running `grid` process
inside disposable containers.
Intent: Running as root is load-bearing for DevOps, but host-root mutation is too
risky for a POC. Containers give enough real filesystem, process, trigger, and
restart behavior to test the design while keeping the host outside the mutation
target.
Constraints: Future POC21 clean runs must not mutate the host root filesystem;
runtime state is bounded under `/tmp/wire-lab-poc21-run/<run_id>/`; containers
may run as root only inside the disposable target environment; inter-agent
traffic must not use Docker volumes as a communication shortcut.
Affects: `implementations/poc21-grid-devops/docs/DESIGN.md`; future
`implementations/poc21-grid-devops/` scripts, compose files, runtime code, and
analyzer gates.

ID: DI-dahaj
Date: 2026-07-11 16:39:34 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Order the POC21 proof slices as actual container self-update first,
ordered configuration replay second, and package/container-image artifact
distribution third.
Intent: The hardest DevOps question is whether a minimal `grid` stage0 can fetch,
verify, swap, restart, and explain a replacement for itself. Ordered replay and
large artifact distribution depend on that bootstrap trust path, so they follow
after the self-update proof.
Constraints: The self-update proof must execute the replacement inside the target
container, not merely stage bytes; prior binary retention is recovery metadata,
not a promise of safe rollback; ordered replay must preserve tested order and
record each action/result as CAS events; package/image artifacts must travel
through in-band CAS/VCS objects.
Affects: `implementations/poc21-grid-devops/docs/DESIGN.md`; future POC21
runtime and analyzer acceptance criteria.

ID: DI-moson
Date: 2026-07-11 16:39:34 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC20 receives only a guardrail for POC21: its CAS event-source,
branch, root-decision, and replay semantics must remain sufficient for DevOps
ordered replay, while machine mutation stays in POC21.
Intent: POC20 is the semantic foundation. POC21 is the operational application.
This split prevents POC20 from absorbing root-machine mutation while still making
sure POC20 does not accidentally preclude the DevOps use case.
Constraints: POC20 remains pre-code semantic planning; POC21 may cite POC20 but
must not require POC20 to execute machine changes; both tracks must avoid
rollback claims and use corrective roll-forward/history-preserving language.
Affects: `implementations/poc20-timeline-pure-function-cas-branches/docs/DESIGN.md`;
`implementations/poc21-grid-devops/docs/DESIGN.md`; `DEV-GUIDE-RESOURCES.md`.

ID: DI-nafat
Date: 2026-07-11 18:15:37 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Simplify POC21 self-update replay requirements. Future POC21 should
record minimum replay facts in local CAS-backed machine-change journals rather
than locking a large universal field list: current running stage0 CID if known,
candidate stage0 CID, local decision state, and post-restart running CID/outcome.
Local reason, impact summary CID, and recovery notes are optional when useful.
Intent: The earlier field list over-specified a protocol that does not exist yet
and risked turning a local DevOps journal into a universal schema. POC21 needs
enough facts to replay and explain a self-update, but the exact payload shape
belongs to the future POC21 pCID, not the outer grid envelope or every protocol.
Constraints: Keep actual container self-update as the first POC21 proof. Do not
make previous-binary recovery metadata mandatory. Do not imply a safe rollback
promise; prior binaries and notes are recovery inputs when present.
Affects: `protocols/wire-lab.d/TODO/TODO-kifok-poc21-grid-devops.md`;
`implementations/poc21-grid-devops/docs/DESIGN.md`; `DEV-GUIDE-RESOURCES.md`.
Supersedes: `DI-dahaj` previous-binary recovery-metadata requirement only.

ID: DI-rigob
Date: 2026-07-29 11:40:39 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC21 will implement one shared Grid language family with a common
data/declaration layer, a finite Gridfile journal profile, and a Turing-complete
`*.grid` program profile. Gridfile remains statically expandable and cannot turn
itself into the program profile through a local pragma. The general program
profile is part of POC21 rather than being deferred to another POC.
Intent: Configuration, ordered machine journals, agents, applications, parsers,
builders, planners, and pure-function services should share CID-native values,
imports, types, descriptors, and tooling without sacrificing the human-readable
and replayable properties that make Gridfile useful as a machine lifetime
journal.
Constraints: The common data layer must not perform effects. Gridfile must keep
finite ordered execution visible before machine mutation. General Grid programs
may be computationally universal, but their local and inter-agent effects remain
mediated by explicit kernel-role and capability promises. `DR-junaz` owns the
still-open canonical syntax, typing, and effect-system decision. `DR-lupiz` owns
the still-open source-header identity decision.
Affects: `protocols/wire-lab.d/TODO/TODO-kifok-poc21-grid-devops.md`;
`implementations/poc21-grid-devops/docs/DESIGN.md`;
`docs/research/DN-gagog-grid-language-profiles-and-runtime-descriptors.md`;
`DR/DR-junaz-canonical-grid-language-design.md`;
`DR/DR-lupiz-grid-source-shebang-identity.md`.

ID: DI-bigap
Date: 2026-07-29 11:40:39 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC21's installed stage0 `grid` binary will read only the common
effect-free data/declaration subset needed for bootstrap configuration. The
first full `*.grid` execution engine will be a Go AST interpreter fetched as a
stage1 module. Grid language semantics will be specified independently of that
engine so later bytecode, WASM, native, or constrained-device implementations
can execute the same language without redefining it.
Intent: Stage0 must remain small enough to establish local identity, owner
anchors, bootstrap roots, peers, and stage1 retrieval without executing
arbitrary configuration code. An AST interpreter is the least expensive route
to a real Turing-complete proof while keeping future runtime portability open.
Constraints: Stage0 configuration cannot loop, recurse, perform I/O, send
messages, or request local resources. Practical program execution may be bounded
by CPU, memory, time, and message capability promises even though the abstract
program profile is Turing complete. Exact language and runtime identity remain
open under `DR-lupiz`; canonical syntax and typing remain open under `DR-junaz`.
Affects: `protocols/wire-lab.d/TODO/TODO-kifok-poc21-grid-devops.md`;
`implementations/poc21-grid-devops/docs/DESIGN.md`;
`docs/research/DN-gagog-grid-language-profiles-and-runtime-descriptors.md`.

## Completed planning tasks

- [x] kifok.1 Create the POC21 planning artifact set.
- [x] kifok.2 Keep POC21 separate from POC20 while adding the POC20 guardrail.
- [x] kifok.3 Design the root-in-container safety boundary and runtime path
  bounds.
- [x] kifok.4 Design actual stage0 self-update from CAS as the first executable
  proof.
- [x] kifok.5 Design ordered DevOps replay with triggers, validation, and
  duplicate-target replay.
- [x] kifok.6 Design in-band package/container-image-style artifact distribution
  as a later POC21 slice.
- [x] kifok.7 Define future analyzer gates for host-root safety, TCP-only
  inter-agent communication, CAS rebuildability, and corrective-history language.
- [x] kifok.8 Update guide resources so POC21 is discoverable as the DevOps
  continuation of POC18, POC19, and POC20.
- [x] kifok.16 Record the shared Grid language family, finite Gridfile profile,
  Turing-complete `*.grid` profile, stage0 data subset, and first AST interpreter
  in `DN-gagog`. Source: `DI-rigob`; `DI-bigap`.
- [x] kifok.17 Run `TE-fakof` across pCID, ordinary language-spec CID, raw
  executable CID, runtime-descriptor CID, dual-CID, and separate-execution-object
  alternatives. The TE recommends a separate exact execution descriptor and
  remains `needs DF` under `DR-lupiz`.

## Future implementation tasks

- [ ] kifok.9 Implement the root-in-container clean-run harness without host-root
  mutation.
- [ ] kifok.10 Implement exact `grid()` CBOR over TCP between source and target
  containers.
- [ ] kifok.11 Implement actual stage0 self-update from CID-named CAS bytes inside
  the target container.
- [ ] kifok.12 Implement ordered machine-change journals with trigger and
  validation events.
- [ ] kifok.13 Implement duplicate-target replay from the same CAS-backed journal.
- [ ] kifok.14 Implement package/container-image-style artifact distribution as
  in-band CAS/VCS objects.
- [ ] kifok.15 Implement analyzer gates for safety, TCP-only communication,
  CAS rebuildability, and corrective-history language.
- [ ] kifok.18 Resolve `DR-lupiz` through the DF questions in `TE-fakof`, then
  record the selected source and execution identity model in a new DI before
  freezing source-header syntax.
- [ ] kifok.19 Run a dedicated TE for `DR-junaz` covering canonical syntax,
  static typing, effect/capability typing, evaluation order, content-addressed
  definitions, and the Gridfile/program shared grammar.
- [ ] kifok.20 Implement the common effect-free data/declaration reader in
  stage0 after the canonical grammar and source-header identity are locked.
- [ ] kifok.21 Implement the fetched stage1 AST interpreter for the full
  Turing-complete `*.grid` profile.
- [ ] kifok.22 Implement the finite Gridfile journal profile over the shared
  value, CID, descriptor, diagnostic, and import machinery; reject program-only
  constructs before machine mutation begins.
- [ ] kifok.23 Add analyzer gates for stage0 bootstrap safety, static Gridfile
  expansion, program-profile Turing-completeness proof, explicit resource
  bounds, source/runtime identity retention, and no pCID-as-language regression.
