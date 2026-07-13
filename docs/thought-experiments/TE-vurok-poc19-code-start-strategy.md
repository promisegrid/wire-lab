# TE-vurok: POC19 code-start strategy

## TE ID

TE-vurok

## Status

decided

## Decision under test

`vumas.4` asks whether POC19 implementation should start by copying POC18
packages or by factoring shared packages into a new production-shaped module.
The practical decision is:

> What is the first code-generation path that preserves POC18's working CAS/VCS,
> sync, TCP, Git-bridge, CLI/core, and token lessons while moving POC19 toward a
> stable stage0 `grid` binary with fetched stage1 microkernel modules?

This TE evaluates the code-start strategy only. It does not choose the exact
stage1 runtime profile, final descriptor schema, or final package names.

## Source corpus

- `protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`
- `implementations/poc19-production-shape/docs/DESIGN.md`
- `protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md`
- `implementations/poc16-secure-tokens-maps-encrypted-payloads/`
- `protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md`
- `implementations/poc18-cas-git-replacement/`
- `docs/thought-experiments/TE-vahoj-poc18-superset-architecture.md`
- `protocols/wire-lab.d/TODO/TODO-nudav-poc20-timeline-pure-function-cas-branches.md`

## Locked inputs

- POC19 is the production-shaped successor to POC18 and must preserve POC18
  behavior unless a later scoped DI explicitly authorizes an exception. Source:
  `DI-lumir`.
- The installed `grid` program is a small stage0 bootstrap seed, not the whole
  microkernel. It fetches remaining microkernel modules by CID before app-specific
  modules. Source: `DI-zitap`.
- App, agent, runtime, spec, and data changes are adopted by Merkle/CAS root CIDs
  with local operator approval. Source: `DI-kodob`; `DI-romak`.
- POC18 code currently lives in its own Go module under
  `implementations/poc18-cas-git-replacement/`, with reusable layer packages
  such as `store`, `graph`, `workspace`, `sync`, `transport`, `economy`, `repo`,
  `bridge`, and `carbundle`, plus POC harness commands such as `poc-sim`,
  `poc-agent`, `poc-analyze`, and `poc-event-collector`.
- POC16 provides important runtime lessons that POC19 must not skip:
  pCID-selected parser/builder roles, pCID-owned payload maps or arrays, embedded
  protocol spec prose for LLM-backed agents, CWT/COSE capability-token profiles,
  encrypted payload behavior, unknown-pCID local non-commitments, and phased
  local lifecycle shutdown for app, parser, and transport roles. Source:
  `TODO-zugok`.
- POC19 should not duplicate CLI/core behavior. User-facing commands should use
  the same shared core or daemon-backed path as automation.
- POC20 requires local CAS to remain the source of truth, with derived indexes or
  projections rebuildable from CAS rather than authoritative.

## Assumptions

- The first POC19 executable slice should prove the stage0/stage1 bootstrap
  shape before adding the full daemon, WebSocket parity, WASI app execution, or
  POC21 DevOps behavior.
- POC18 remains useful implementation evidence, but its harness shape is not the
  desired production shape.
- Inter-agent behavior remains promise-shaped over exact `grid()` messages; no
  RPC shortcut is acceptable.
- pCID remains a protocol selector, not a route, peer address, app address,
  operation code, or message type.

## Alternatives

### Alt A: Copy POC18 packages wholesale into POC19

POC19 would copy the current POC18 module or most POC18 packages into
`implementations/poc19-production-shape/`, rename imports, and then start
editing the copy toward stage0/stage1.

This is fast at first. It preserves tests and behavior by file duplication, and
it gives implementers a working baseline without designing extraction seams up
front.

The cost is structural. POC18 intentionally contains POC harness binaries,
collector/analyzer machinery, local fixture assumptions, and command surfaces
that POC19 is trying to move beyond. Copying wholesale makes those shapes look
canonical. It also creates two diverging implementations of CAS/VCS, sync,
tokens, and diagnostics unless every later fix is ported twice.

### Alt B: Fresh-build POC19 without POC18 package reuse

POC19 would start with a clean module and write stage0, descriptors, CAS, VCS,
sync, transport, and CLI/daemon seams from scratch.

This maximizes architectural cleanliness and avoids carrying POC18 harness
artifacts forward. It also forces the team to name the production-shaped roles
instead of inheriting POC names.

The cost is regression risk. POC18 already contains real work for sparse CAS,
Rabin chunks, reference sets, Git bridge import/export/push/pull, CAR bundles,
token-shaped storage economics, continuous sync, TCP transport, exact CBOR
diagnostics, and repo-local `.grid` behavior. A fresh rewrite can easily drop one
of those lessons and pass a narrow stage0 demo while failing POC superset
discipline.

### Alt C: Hybrid staged extraction

POC19 starts with a small fresh stage0 scaffold and a production-shaped module
boundary, while treating POC18 as the source baseline for behavior and tests.
Reusable POC18 layer behavior is factored or copied in small, reviewed slices
only when it is placed behind POC19 roles: stage0 bootstrap, stage1 daemon/client
module, node CAS, repo state, sync, transport, Git bridge adapter, diagnostics,
and token/economy helpers.

This keeps POC19 from becoming a renamed POC18 harness while still preserving the
working implementation evidence. It makes the first proof small: stage0 reads
local config or a bootstrap root, fetches a descriptor and executable object by
CID, verifies exact bytes, starts one stage1 module, and records a readiness
result. After that, POC18 CAS/VCS/sync behavior can move behind fetched stage1
roles without duplicating CLI/core logic.

The cost is discipline. Extraction takes more deliberate sequencing than a raw
copy. The implementation must keep POC18 tests and behavior visible as regression
evidence, while refusing to preserve POC-only command names or fixture shortcuts
as production API.

## Scenario analysis

### Scenario 1: First stage0/stage1 proof

Alice installs one `grid` binary. The binary needs enough local config, CID,
CAS-fetch, descriptor, and process-start behavior to start a CID-named stage1
module.

- Alt A gives Alice too much old POC18 shape too early. The copied code already
  knows about POC18 CLI and harness roles, so stage0 can blur into a monolith.
- Alt B can keep stage0 clean, but risks implementing a demo-only CAS and
  descriptor path that later has to be replaced.
- Alt C fits the scenario. The stage0 scaffold is fresh and intentionally small,
  while POC18's store/CID/test patterns remain the reference for correctness.

### Scenario 2: Preserving POC18 VCS behavior

Bob expects POC19 to preserve POC18 behavior: `.grid` state, sparse CAS,
reference sets, Rabin chunks, Git bridge import/export/push/pull, diagnostic
rendering, and continuous peer sync.

- Alt A preserves behavior initially but creates a fork that will diverge from
  POC18 unless every fix is duplicated.
- Alt B makes regression likely because each feature has to be re-created.
- Alt C preserves behavior through targeted extraction plus regression gates.
  POC18 remains the behavioral baseline; POC19 decides which package seams become
  production roles.

### Scenario 3: Avoiding CLI/core duplication

Carol writes automation against `grid status`, `grid snapshot`, and future
`grid run`. The same behavior must be used by human CLI, daemon control, and
automation.

- Alt A risks inheriting POC18 command implementations as direct local logic and
  then adding a daemon path beside them.
- Alt B can design the seam cleanly but may re-solve problems POC18 already
  solved.
- Alt C makes the seam explicit: CLI code is thin, core behavior sits behind a
  shared production-shaped package or daemon module, and POC18 command code is
  mined for behavior rather than copied as the surface.

### Scenario 4: Adding TCP/WebSocket parity

Dave needs exact same `grid()` bytes over TCP and WebSocket. Transport adapters
must not know app semantics.

- Alt A may carry POC18 TCP assumptions into POC19 and leave WebSocket as a later
  bolt-on.
- Alt B can define parity cleanly but risks losing POC18's tested TCP framing and
  exact-message retention.
- Alt C extracts the existing exact-byte framing behavior into a transport role
  and adds WebSocket as an equal adapter to that role.

### Scenario 5: POC20 timeline compatibility

Ellen wants later POC20 semantics: CAS as the chronological event source,
rebuildable projections, branch-visible roots, and branch-aware token behavior.

- Alt A may preserve hidden POC18 fixture state that is not event-source shaped.
- Alt B can design from the POC20 model but may ignore POC18's proven sparse CAS
  and object-transfer behavior.
- Alt C is the best fit if every extracted POC18 subsystem is required to expose
  source facts as CAS objects and rebuildable views rather than hidden mutable
  state.

### Scenario 6: POC21 DevOps inheritance

Frank later uses POC21 to apply root-in-container machine changes using POC19's
stage0 and CAS/VCS substrate.

- Alt A risks dragging collector/analyzer or fixture assumptions into DevOps.
- Alt B risks producing a stage0 proof that cannot run the richer POC18-derived
  VCS/CAS behavior needed by DevOps.
- Alt C gives POC21 a staged substrate: stage0 is narrow, stage1 roles are
  fetched by CID, and CAS/VCS behavior is production-shaped before root-like
  machine changes are attempted.

### Scenario 7: POC16 parser, token, and lifecycle lessons

Grace runs a mixed app/kernel node where the transport listener sees slot-0 pCID
bytes, the selected parser role owns payload decoding, and app/runtime tokens use
CWT/COSE-shaped promises. During shutdown, app sessions, parser sessions, and
transport sessions must close in dependency order so terminal lifecycle records
are not lost.

- Alt A risks copying POC18 without the POC16 parser/builder split and lifecycle
  shutdown lessons, because POC18 is stronger on CAS/VCS than on the POC16
  runtime boundary corrections.
- Alt B can implement the POC16 model cleanly, but would have to recreate POC18
  CAS/VCS behavior in parallel.
- Alt C lets POC19 start with fresh stage0 and parser/builder role boundaries
  informed by POC16, while extracting POC18 CAS/VCS behavior only behind those
  production-shaped boundaries.

## Conclusion

Alt C, hybrid staged extraction, is the surviving alternative.

Alt A is rejected because copying POC18 wholesale would preserve too much
non-production harness shape and create duplicated implementations. Alt B is
rejected because a fresh rewrite would likely regress POC18 behavior and delay
the production-shaped proof. Alt C preserves the POC18 behavioral baseline while
forcing POC19 to start with a small fresh stage0 and production-shaped stage1
role boundaries. POC16 strengthens the same conclusion: stage0 and stage1 should
be shaped around pCID-selected parser/builder roles, CWT/COSE token promises,
embedded spec docs, encrypted-payload constraints, and local lifecycle shutdown
from the start, not retrofitted after copying POC18.

The locked implementation direction should be:

- start POC19 with a fresh stage0 `grid` scaffold;
- treat POC18 as a source baseline and regression oracle, not a long-term import
  dependency;
- extract or copy only small behavior slices into POC19 roles after reviewing
  each slice against the stage0/stage1 architecture;
- keep POC-only harness commands out of the production command surface;
- require POC18-superset gates before considering POC19 complete.

## Implications for open TODOs and pending DIs

- `vumas.4` can be marked complete by a DI that locks hybrid staged extraction.
- `vumas.5` should scaffold the fresh stage0 `grid` binary and descriptor
  fixture first, keeping parser/builder role seams visible.
- `vumas.6` should move POC18-compatible `.grid` and sparse-CAS behavior behind
  POC19 node-CAS/repo-state roles.
- `vumas.7` should extract exact-byte transport framing before adding WebSocket.
- `vumas.8` should add WASI app execution only after descriptor/stage1 readiness
  is working, and should preserve POC16 CWT/COSE token and encrypted-payload
  lessons for app/runtime capabilities.
- `vumas.9` should include POC18-superset gates that catch dropped CAS/VCS,
  Git-bridge, continuous-sync, token, diagnostic, and TCP behavior.
- `vumas.16` should include the POC16 pCID parser/builder inventory, not only
  POC18 VCS/CAS pCIDs.

## Decision status

Locked by `DI-nupag` in
`protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`.
