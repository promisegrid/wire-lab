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

## Tasks

- [x] vumas.1 Lock the POC19 design-doc-first decision in `DI-lumir`.
- [x] vumas.2 Write the first POC19 production-shape design document at
  `implementations/poc19-production-shape/docs/DESIGN.md`.
- [x] vumas.3 Review the POC19 design against POC16, POC17, and POC18 to ensure
  inherited lessons are not lost. The inheritance matrix is now in
  `implementations/poc19-production-shape/docs/DESIGN.md`. Source: `DI-topab`.
- [ ] vumas.4 Decide whether POC19 implementation starts by copying POC18
  packages or by factoring shared packages into a new production-shaped module.
- [ ] vumas.5 Scaffold `implementations/poc19-production-shape/` with one stable
  `grid` bootstrap binary plus fetched daemon/client microkernel modules.
- [ ] vumas.6 Implement daemon-managed local CAS and VCS config discovery with
  POC18-compatible `.grid` repo state.
- [ ] vumas.7 Implement TCP and WebSocket transport adapters over the same exact
  PromiseGrid message framing.
- [ ] vumas.8 Implement signed app reference-set install and `grid run` WASI
  execution from VCS/CAS.
- [ ] vumas.9 Add analyzer/regression gates proving POC18 superset behavior,
  exact-message retention, TCP/WebSocket parity, and promise-first vocabulary.
- [ ] vumas.10 Run and archive a clean POC19 regression after implementation
  begins.
- [x] vumas.11 Lock the minimum-microkernel rule: one installed `grid` binary
  bootstraps from operator-adopted root CIDs, while app, agent, runtime
  executable, spec, and data changes are fetched from CAS/peers without replacing
  the binary. Source: `DI-kodob`.
- [x] vumas.12 Lock the operator-visible root-adoption flow: candidate-root fetch,
  closure/signature/spec/capability review, impact summary, explicit local
  approval, rollback root retention, and replay context. Source: `DI-guhil`.
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
