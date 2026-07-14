# TODO-vumas: POC19 production shape

## Status

Planned. Owns the production-shaped successor to POC18 under
`implementations/poc19-production-shape/`. The first artifact is a design
document, not executable code. Source: `DI-lumir`; `DI-kodob`.

POC19 should turn the lessons from POC16, POC17, and POC18 into a single
`grid` binary that can run as a local PromiseGrid daemon/microkernel, expose the
VCS/CLI surface, fetch code and data from peers over TCP or WebSocket, and run
fetched apps from VCS/CAS state under local promise and capability-token
constraints. The stable installed binary is a small bootstrap seed for the
minimum microkernel: it self-updates only after local owner approval, then
fetches remaining microkernel modules before any app-specific modules. App,
agent, runtime executable, and data changes are adopted by CID-addressed CAS
roots, not by replacing the binary. Source: `DI-lumir`; `DI-kodob`; `DI-zitap`.

POC20 is now the parallel semantic-model track for promises as timeline
assertions, pure-function agents, CAS branches, local/group timelines, and
branch-aware token double-spend behavior. POC19 remains the production-shaped
plumbing path, but should avoid hidden token or app-run state that would prevent
future visible branch histories. Source: `DI-kakos`; `TE-lodom`;
`TODO-nudav`.

## Decision Intent Log

ID: DI-lumir
Date: 2026-07-07 20:18:23 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Plan POC19 as a production-shaped POC18 successor whose first artifact
is `implementations/poc19-production-shape/docs/DESIGN.md`; the future runtime
target is one binary named `grid` that can run as a daemon/client pair, act as a
VCS, fetch non-kernel code and data from VCS/CAS over TCP or WebSocket, and
execute fetched apps with WASI as the first execution profile.
Intent: POC18 proved important CAS/VCS, TCP, token, diagnostic, and sparse-sync
pieces, but it still splits production roles across `grid`, `poc-agent`,
`poc-sim`, `poc-analyze`, and collector binaries. POC19 should stop looking like
a harness and start looking like a deployable PromiseGrid node. Non-kernel code
should be installed by checking signed app reference sets into VCS/CAS, and
`grid run` or the daemon should load executable code, container images, WASI
modules, and data from that same substrate rather than from a separate package
manager or side channel. The design must preserve Promise Theory framing: every
inter-agent behavior is a voluntary promise over exact `grid()` messages, local
trust remains local, pCID remains a protocol selector, and resource control is a
local conditional capability promise.
Constraints: Preserve the POC superset discipline unless a later scoped DI
explicitly authorizes an exception; inherit POC18 CAS/VCS/continuous-sync/Git
bridge behavior; inherit POC16 pCID-owned arity, parser/builder role, secure
CWT/COSE token, encrypted-payload, and kernel-role lessons; inherit POC17
compact constrained-message and binary-CID/base32-printable discipline; keep one
top-level semantic action `promise`; use `grid([42(pCID),
...protocol-defined-slots])`; do not use pCID as peer address, app address,
operation code, route, repository name, or message type; use binary CIDs on wire
and CIDv1 base32 text when printable; make TCP and WebSocket equal first-class
transport targets for the design; make WASI the first execution profile before
native binary or container execution; keep observer/analyzer behavior as test
machinery rather than production trust infrastructure.
Affects: `protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`;
`protocols/wire-lab.d/TODO/TODO.md`;
`implementations/poc19-production-shape/docs/DESIGN.md`;
`implementations/README.md`; `DEV-GUIDE-RESOURCES.md`.

ID: DI-topab
Date: 2026-07-08 21:53:59 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Complete the `vumas.3` POC19 inheritance review inside the existing
POC19 design document instead of creating a separate artifact.
Intent: POC19 code generation should begin only after the design explicitly
checks that POC16, POC17, and POC18 lessons are preserved or assigned to later
POC19 tasks. Keeping the review in the design document makes it part of the
implementation contract rather than a side note.
Constraints: Do not start code generation; review only inherited behavior,
architecture lessons, and acceptance gaps. Treat unresolved implementation
choices as follow-ups for `vumas.4` or later tasks. Keep PromiseGrid vocabulary
promise-first and do not use pCID as an address, operation, route, app name, or
message type.
Affects: `implementations/poc19-production-shape/docs/DESIGN.md`;
`protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`;
`DEV-GUIDE-RESOURCES.md`.

ID: DI-kodob
Date: 2026-07-09 19:19:00 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Lock POC19's installed `grid` binary as the minimum microkernel. The
binary provides bootstrap, transport, local CAS verification, pCID parser
dispatch, local config, operator approval, runtime launch, and local
resource/capability roles. App code, agent code, runtime executable objects,
container/WASI/native artifacts, specs, and data are fetched by CID from
CAS/peers from an operator-adopted Merkle/root CID.
Intent: A person should install one simple binary on a laptop, Raspberry Pi,
server, or similar device and should not need to update that binary whenever app
or agent code changes. Code changes should be PromiseGrid data: a new root CID
is proposed, fetched, verified, and adopted only by local operator promise.
Constraints: Do not embed changing app/agent/runtime executables in the binary.
First-run config may name a bootstrap root CID; later root updates require local
operator approval and should be recorded as promise-shaped local state. The
binary may still evolve for true microkernel/protocol substrate changes, but not
for normal app/runtime updates.
Affects: `implementations/poc19-production-shape/docs/DESIGN.md`;
`protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/DESIGN.md`;
`DEV-GUIDE-RESOURCES.md`; `README.md`.

ID: DI-guhil
Date: 2026-07-10 14:26:37 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Extend POC19 root adoption into an operator-visible local flow with
candidate-root fetch, closure verification, local signer trust criteria,
requested host-capability review, impact summary, explicit local approval,
rollback root retention, and replay context. Treat every bridge to a local host
resource as a narrow capability promise made by the local node, and keep secrets
behind operation-scoped local services rather than ordinary CAS payloads.
Intent: POC19 should prove that app, agent, runtime, protocol, fixture, UI, or
policy changes are adopted as CID-addressed roots without silently widening local
resource promises or leaking secrets into config, CAS, diagnostics, logs, prompts,
or UI output.
Constraints: Use positive PromiseGrid vocabulary: local approval, local
capability promise, local signer trust criteria, local event, rollback root, and
replayable context. Do not add production UI scope, any globally privileged
update or signer role, secret plaintext in CAS, or sideband package-manager
semantics. This remains pre-code planning until `vumas.4` and implementation
tasks are separately started.
Affects: `protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`;
`implementations/poc19-production-shape/docs/DESIGN.md`;
`protocols/wire-lab.d/TODO/TODO-nudav-poc20-timeline-pure-function-cas-branches.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/DESIGN.md`;
`DEV-GUIDE-RESOURCES.md`; `README.md`.

ID: DI-romak
Date: 2026-07-11 18:15:37 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Supersede POC19 root and recovery wording so current planning uses
`Merkle/CAS root CID`, `prior Merkle/CAS root CID`, `corrective Merkle/CAS root
CID`, and explicit full-state restore attempts rather than `Merkle/root CID` or
`rollback root`.
Intent: POC19 should not imply that returning to a prior binary or root can undo
machine, peer, network, or physical side effects. The correct model is
history-preserving recovery: retain prior roots for replay, adopt corrective
roots when useful, and treat full-state restore as a narrow explicit operation
only when the affected program-and-data state is available.
Constraints: Do not rewrite historical DI bodies. Use this DI for current tasks,
README, guide prose, and downstream POC20/POC21 consistency. Retaining previous
binary or root material is recovery input, not a safe rollback promise.
Affects: `protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`;
`implementations/poc19-production-shape/docs/DESIGN.md`; `README.md`;
`DEV-GUIDE-RESOURCES.md`.
Supersedes: `DI-kodob` root-name wording only; `DI-guhil` rollback-root wording
only.

ID: DI-zitap
Date: 2026-07-11 12:12:01 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Narrow the stable installed `grid` binary to a small bootstrap seed for
the minimum microkernel. The stable binary contains only enough functionality to
detect a candidate new version of itself, fetch it by CID from locally trusted
peers using minimal PromiseGrid message handling, verify the exact bytes, ask an
owning agent or human for local approval, restart into the approved version, then
fetch the remaining modules that make up the minimum microkernel before fetching
application-specific modules.
Intent: The installed binary should be stable, small, and rarely changed. Normal
microkernel role changes, app changes, agent changes, runtime changes, protocol
spec changes, and data changes should travel as CID-addressed CAS modules or
roots instead of forcing binary replacement.
Constraints: Peer trust is local and relationship-relative. Self-update approval
is a local promise by an owning agent or human. Remaining microkernel modules are
fetched before app-specific modules. Do not treat any peer, root, or signer as a
global update role.
Affects: `protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`;
`implementations/poc19-production-shape/docs/DESIGN.md`; `README.md`;
`DEV-GUIDE-RESOURCES.md`; `protocols/wire-lab.d/TODO/TODO.md`.
Supersedes: `DI-kodob` only where it implied the stable installed binary itself
contains all local daemon/microkernel roles. `DI-kodob` remains active for the
operator-adopted CID-root update model.

ID: DI-nupag
Date: 2026-07-13 14:52:09 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Start POC19 implementation with hybrid staged extraction: create a
fresh stage0 `grid` scaffold, treat POC18 as the source baseline and regression
oracle, and factor or copy only reviewed POC18 behavior slices into
production-shaped POC19 stage1 roles.
Intent: POC19 must preserve POC18 CAS/VCS, sync, transport, Git-bridge, token,
diagnostic, and CLI/core lessons without copying POC18's POC-only harness shape
as the production architecture. It must also preserve POC16's runtime lessons
around pCID-selected parser/builder roles, embedded protocol specs, CWT/COSE
tokens, encrypted payloads, and local lifecycle shutdown. A fresh stage0 keeps
the installed binary small; targeted extraction preserves the working behavior
that proves the superset.
Constraints: Do not use long-term imports from
`implementations/poc18-cas-git-replacement/...` as the POC19 architecture. Do not
copy `poc-sim`, `poc-agent`, `poc-analyze`, or collector command shapes as the
production surface. Keep CLI behavior on the same shared core or daemon-backed
path as automation. Treat POC16 as the source baseline for parser/builder,
protocol-spec, CWT/COSE token, encrypted-payload, and lifecycle role boundaries.
Preserve POC18-superset regression gates before declaring POC19 complete.
Affects: `protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`;
`docs/thought-experiments/TE-vurok-poc19-code-start-strategy.md`;
`implementations/poc19-production-shape/docs/DESIGN.md`.

ID: DI-hofaz
Date: 2026-07-13 16:14:50 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC19 durable storage uses one CID-keyed CAS namespace as the source of
truth. Raw Rabin chunks remain CIDv1 `raw` exact-byte objects; true durable graph
objects may use CIDv1 `dag-cbor`; exact `grid()` messages, Markdown specs,
executables, encrypted bytes, and CAR artifacts remain exact byte objects under
their own CIDs. `chunks`, `messages`, `dag-cbor`, `car`, and similar views are
rebuildable projections, not authoritative source-of-truth directories.
Intent: Stage0 should be able to fetch and verify exact bytes by CID without
knowing every durable object profile. Stage1 can provide richer diagnostics,
DAG-CBOR validation, CAR import/export, and derived browsing views while keeping
CAS itself simple, sparse, and replayable from exact object bytes.
Constraints: Do not treat POC18's `chunks/*.bin` and `objects/*.cbor` layout as
the POC19 canonical store. Do not imply that a `.cbor` filename means a CIDv1
`dag-cbor` object. CAR files are transfer/archive packages and review artifacts,
not the authority over contained object identity. Local indexes and profile views
must be rebuildable from CAS-resident facts.
Affects: `protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`;
`docs/thought-experiments/TE-lirum-poc19-storage-object-profile.md`;
`implementations/poc19-production-shape/docs/DESIGN.md`;
`DEV-GUIDE-RESOURCES.md`.

ID: DI-topiv
Date: 2026-07-13 16:44:39 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC19's first stage1 bootstrap proof uses a native/static executable
object, not a WASI stage1 module. Stage0 remains a small bootstrap seed that
fetches a descriptor and native/static stage1 executable by CID, verifies exact
bytes, signer or local trust criteria, platform constraints, and local capability
requirements, materializes a runnable copy into a grid-owned execution cache,
starts stage1 through the host process mechanism, passes minimal bootstrap facts,
and records readiness as a local event. WASI/WASM remains a first-class portable
app/runtime profile under stage1.
Intent: Hostful nodes need USB, serial/device access, process monitoring and
control, platform execution checks, and local host capability mediation before
portable app modules can be useful. Putting those roles in native/static stage1
keeps stage0 small while avoiding a circular dependency on unverified stage1
trust-policy code or a premature WASI loader in stage0.
Constraints: Do not require stage0 to contain a WASI loader for the first proof.
Do not move ordinary host adapters, device access, or process supervision into
stage0. Homebrew, signed installers, and other package managers may distribute
stage0, but they do not approve fetched stage1 CIDs. iOS and iPadOS are bundled
clients or control surfaces unless a later signed-app distribution path is
explicitly designed. Stage0 self-update remains separate from runtime-root
adoption, and retaining prior binaries or roots must not imply true rollback.
Affects: `protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`;
`docs/thought-experiments/TE-sunag-poc19-native-stage1-bootstrap.md`;
`implementations/poc19-production-shape/docs/DESIGN.md`;
`implementations/poc20-timeline-pure-function-cas-branches/docs/DESIGN.md`;
`DEV-GUIDE-RESOURCES.md`.

ID: DI-bugik
Date: 2026-07-13 17:13:56 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC19 uses pCID as the local kernel/stage1 protocol selector for
routing exact `grid([42(pCID), ...])` messages to registered handler(s) that
promise to parse or build that pCID. A pCID is still not a peer address, app
address, route, operation, command name, repository name, RPC method, or message
kind. Handler-owned parsing interprets the remaining slots and any payload-local
destination, app, route, operation, repository, or VCS semantics.
Intent: POC19 code generation needs a concrete handler inventory before parser
scaffolding so stage0 and native/static stage1 can stay small and deterministic
without regressing into pCID-as-address or pCID-per-message-kind designs.
Constraints: Stage0 recognizes only the minimal bootstrap/fetch/self-update
surface needed to fetch and launch native/static stage1. Native/static stage1
owns the main handler registry. Protocol families should stay coarse enough that
payload variants remain payload semantics. Standalone spec docs and base32 pCID
aliases are required before a handler becomes an active POC19 runtime handler.
Affects: `protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`;
`implementations/poc19-production-shape/docs/DESIGN.md`;
`DEV-GUIDE-RESOURCES.md`.

ID: DI-tuvub
Date: 2026-07-13 17:27:53 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC19's first native/static stage1 proof requires a launch-attempt
local event, not a mandatory readiness result. Stage0 fetches a descriptor and
executable object by CID, verifies exact bytes and local capability/trust
criteria, materializes the executable into an execution cache, starts it through
a portable host process launch mechanism, and records descriptor CID, executable
CID, adopted Merkle/CAS root CID when present, execution-cache path, platform,
approval or rejection outcome, and process-launch outcome. Readiness is optional
supplemental information if the stage1 process reports it before timeout.
Intent: The first proof should show that stage0 can perform a narrow
CID-verified launch handoff to fetched native/static stage1 code without forcing
a stage1 control protocol or readiness handshake into the minimum bootstrap
binary.
Constraints: Use portable process spawn semantics such as Go `os/exec`, mapping
to `CreateProcess`-style behavior on Windows and normal process launch on Unix;
do not require Unix same-PID `execve`. Do not add a WASI loader, broad
trust-policy engine, device/process adapter layer, or app-specific runtime logic
to stage0 for this proof. Do not treat missing readiness as failure when the
launch-attempt local event records the process-launch outcome.
Affects: `protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`;
`docs/thought-experiments/TE-sunag-poc19-native-stage1-bootstrap.md`;
`implementations/poc19-production-shape/docs/DESIGN.md`;
`DEV-GUIDE-RESOURCES.md`.
Supersedes: `DI-topiv` launch-record/readiness-required wording only.

## Tasks

- [x] vumas.1 Lock the POC19 design-doc-first decision in `DI-lumir`.
- [x] vumas.2 Write the first POC19 production-shape design document at
  `implementations/poc19-production-shape/docs/DESIGN.md`.
- [x] vumas.3 Review the POC19 design against POC16, POC17, and POC18 to ensure
  inherited lessons are not lost. The inheritance matrix is now in
  `implementations/poc19-production-shape/docs/DESIGN.md`. Source: `DI-topab`.
- [x] vumas.4 Decide whether POC19 implementation starts by copying POC18
  packages or by factoring shared packages into a new production-shaped module.
  `TE-vurok` rejects both wholesale copy and fresh rewrite; POC19 starts with
  hybrid staged extraction. Source: `TE-vurok`; `DI-nupag`.
- [ ] vumas.5 Scaffold `implementations/poc19-production-shape/` with one stable
  `grid` stage0 bootstrap binary plus a fetched native/static stage1 descriptor
  and executable proof. Stage0 must verify exact CIDs, local trust criteria,
  platform constraints, and local capability requirements before materializing
  stage1 into an execution cache, launching it through a portable host process
  mechanism, and recording a launch-attempt local event. Source: `TE-sunag`;
  `DI-topiv`; `DI-tuvub`.
- [ ] vumas.6 Implement daemon-managed local CAS and VCS config discovery with
  POC18-compatible `.grid` repo state.
- [ ] vumas.7 Implement TCP and WebSocket transport adapters over the same exact
  PromiseGrid message framing.
- [ ] vumas.8 Implement signed app reference-set install and `grid run` WASI/WASM
  execution from VCS/CAS under native/static stage1. Stage0 must not need a WASI
  loader for the first proof. Source: `TE-sunag`; `DI-topiv`.
- [ ] vumas.9 Add analyzer/regression gates proving POC18 superset behavior,
  exact-message retention, TCP/WebSocket parity, promise-first vocabulary,
  native/static stage1 launch from an execution cache, no dependency on
  unverified stage1 trust-policy code, platform/package-manager separation,
  stage0 self-update separation, required launch-attempt local-event fields,
  optional readiness if reported before timeout, failure reporting, and no true
  rollback claim. Source: `TE-sunag`; `DI-topiv`; `DI-tuvub`.
- [ ] vumas.10 Run and archive a clean POC19 regression after implementation
  begins.
- [x] vumas.11 Lock the minimum-microkernel rule: one installed `grid` binary
  bootstraps from operator-adopted Merkle/CAS root CIDs, while app, agent, runtime
  executable, spec, and data changes are fetched from CAS/peers without replacing
  the binary. Source: `DI-kodob`; `DI-romak`.
- [x] vumas.12 Lock the operator-visible root-adoption flow: candidate-root fetch,
  closure/signature/spec/capability review, impact summary, explicit local
  approval, prior/corrective Merkle/CAS root CIDs, and replay context. Source:
  `DI-guhil`; `DI-romak`.
- [x] vumas.13 Lock host capabilities as narrow local promises covering network,
  filesystem, device, host-function, secret-reference, execution-resource, and
  storage access. Source: `DI-guhil`.
- [x] vumas.14 Lock operation-scoped secret services as the POC19 direction:
  callers request narrow signing, unwrap, mint, rotate, revoke, or denial records;
  plaintext secrets do not become ordinary CAS/config/log/prompt/UI payloads.
  Source: `DI-guhil`.
- [x] vumas.15 Narrow the stable installed `grid` binary to a small bootstrap seed
  that can self-update with local owner approval, restart, fetch remaining
  microkernel modules, and only then fetch application-specific modules. Source:
  `DI-zitap`.
- [x] vumas.16 Produce the POC19 pCID inventory before parser scaffolding or app
  reference-set execution proceeds beyond the stage0 bootstrap proof. The POC19
  design now inventories stage0-minimal, native stage1/kernel-role, app/runtime,
  and bridge/interoperability handler families, and requires standalone spec docs
  plus base32 pCID aliases before a handler becomes active. Source: `DI-nupag`;
  `DI-bugik`.
- [x] vumas.17 Resolve the POC19 storage object profile for raw chunks,
  DAG-CBOR, GRID-CBOR, CAR files, and local store layout before durable POC19
  stores are scaffolded. `TE-lirum` locks one CID-keyed CAS namespace as source
  of truth, with raw chunks as CIDv1 `raw`, true graph objects optionally
  CIDv1 `dag-cbor`, CAR as transfer/archive packaging, and profile views as
  rebuildable projections. Source: `TE-lirum`; `DI-hofaz`.
- [x] vumas.18 Lock native/static stage1 bootstrap as the first POC19 runtime
  proof. `TE-sunag` rejects WASI-first stage1 for the first proof because
  hostful nodes need USB/device, process monitoring/control, platform checks, and
  host capability mediation in fetched stage1 rather than in stage0. WASI/WASM
  remains a portable app/runtime profile under stage1. Source: `TE-sunag`;
  `DI-topiv`.
