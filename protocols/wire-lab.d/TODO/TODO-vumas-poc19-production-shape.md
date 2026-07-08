# TODO-vumas: POC19 production shape

## Status

Planned. Owns the production-shaped successor to POC18 under
`implementations/poc19-production-shape/`. The first artifact is a design
document, not executable code. Source: `DI-lumir`.

POC19 should turn the lessons from POC16, POC17, and POC18 into a single
`grid` binary that can run as a local PromiseGrid daemon/microkernel, expose the
VCS/CLI surface, fetch code and data from peers over TCP or WebSocket, and run
fetched apps from VCS/CAS state under local promise and capability-token
constraints. Source: `DI-lumir`.

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

## Tasks

- [x] vumas.1 Lock the POC19 design-doc-first decision in `DI-lumir`.
- [x] vumas.2 Write the first POC19 production-shape design document at
  `implementations/poc19-production-shape/docs/DESIGN.md`.
- [ ] vumas.3 Review the POC19 design against POC16, POC17, and POC18 to ensure
  inherited lessons are not lost.
- [ ] vumas.4 Decide whether POC19 implementation starts by copying POC18
  packages or by factoring shared packages into a new production-shaped module.
- [ ] vumas.5 Scaffold `implementations/poc19-production-shape/` with one
  `grid` binary that has daemon/client modes.
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
