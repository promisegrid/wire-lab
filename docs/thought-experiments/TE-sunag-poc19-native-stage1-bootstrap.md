# TE-sunag: POC19 native stage1 bootstrap

## TE ID

TE-sunag

## Status

decided, refined

## Decision under test

POC19 needs a pre-code decision about the first stage1 runtime. Earlier POC19
planning kept the first stage1 runtime open: a WASI stage1 would prove portable
runtime behavior early, while a native/static stage1 would keep stage0 smaller
and provide direct host integration. The updated feedback in `/tmp/feedback.md`
and Steve's clarification resolve the decision pressure: a real node needs host
capability mediation for USB, serial devices, process monitoring and control,
platform checks, execution-cache materialization, and later WASI/WASM app
execution. The decision under test is therefore:

> Should the first POC19 stage1 bootstrap module be native/static, WASI, or a
> split sidecar design?

This TE evaluates that question before `vumas.5` scaffolding and `vumas.8`
runtime execution work proceed.

## Assumptions

- Stage0 is the installed `grid` binary and remains a small bootstrap seed.
- Stage0 must verify stage1 without depending on unverified stage1 trust-policy
  code.
- Stage0 can fetch exact CID-named descriptor and executable objects from local
  CAS or trusted peers.
- Stage1 is the rest of the minimum microkernel module set, not an application
  distribution unit.
- App/runtime modules remain CID-addressed CAS/VCS objects named by locally
  adopted roots.
- POC19 must preserve PromiseGrid vocabulary: local promises, local approval,
  host capability promises, local events, corrective roots, and replayable CAS
  objects.
- WASI/WASM remains important for portable app/runtime modules, even if it is
  not the first stage1 substrate.

## Alternatives

### Alt A: WASI stage1 first

Stage0 embeds enough WASI runtime support to load the first stage1 module as a
WASI module. The first proof starts a portable stage1 module before adding
native host adapters.

### Alt B: Native/static stage1 first

Stage0 fetches a native/static stage1 executable object by CID, verifies it,
materializes a runnable copy into a grid-owned execution cache, starts it with
the host process mechanism, and records readiness as a local event. Stage1 then
hosts WASI/WASM, OCI, and native app/runtime profiles.

### Alt C: WASI stage1 plus native sidecar

Stage0 starts a WASI stage1 module and a native sidecar for hostful behavior.
The sidecar owns USB, serial, process monitoring/control, platform checks, and
other host capability promises.

### Alt D: Stage0 owns host adapters directly

Stage0 remains the first executable and directly owns USB, process control,
device discovery, local platform checks, and later WASI loading.

## Scenario analysis

### Scenario 1: Alice installs `grid` on a normal Linux workstation

Alice installs one stable `grid` binary. On first run, stage0 reads local config,
trust anchors, and a bootstrap Merkle/CAS root CID. It fetches a descriptor and a
stage1 executable object by CID.

- Alt A makes the first stage1 module portable, but stage0 must contain enough
  WASI loader code before stage1 can start. If the node later needs local USB,
  process monitoring, or native platform checks, those capabilities must be added
  to stage0 or delegated to another native component.
- Alt B keeps stage0 narrow: fetch, verify, materialize, launch, and record
  readiness. Native stage1 can then promise hostful behavior and expose narrow
  capability promises to later WASI/WASM modules.
- Alt C adds coordination complexity before the first useful proof. The WASI
  module cannot perform hostful work without a native sidecar, and the sidecar
  becomes a second stage1-like object with its own launch, trust, and lifecycle
  surface.
- Alt D makes stage0 grow into the host adapter layer. That conflicts with the
  minimum-bootstrap rule because ordinary host capability evolution would require
  replacing the installed binary.

Alt B best matches the stage0/stage1 separation.

### Scenario 2: Bob runs a node that talks to a USB scale and label printer

Bob's node needs local device discovery, serial or USB access, and process
supervision. The local node may promise access to these resources only under
narrow host capability tokens.

- Alt A does not give the first stage1 module ambient host authority. That is a
  good sandboxing property for apps, but a poor substrate for the hostful module
  that must mediate local hardware.
- Alt B places hostful adapters in native stage1, where they can remain fetched,
  CID-verified modules rather than baked into stage0.
- Alt C can work, but makes the sidecar the real hostful stage1 while the WASI
  module becomes a coordination layer. That indirection is not useful for the
  first proof.
- Alt D works technically, but every new host adapter risks bloating stage0.

Alt B gives Bob a local host-capability surface without turning stage0 into a
hardware kernel.

### Scenario 3: Carol runs portable app modules

Carol wants portable WASI/WASM app modules so the same app can run on Linux,
macOS, Windows, and future constrained runtimes.

- Alt A proves portable stage1 early, but conflates the first microkernel
  bootstrap module with portable app execution.
- Alt B still supports WASI/WASM as first-class app/runtime profiles. Native
  stage1 provides the WASI host and controls which filesystem, network, device,
  secret, CPU, memory, time, and process promises are granted.
- Alt C supports WASI but creates two mandatory bootstrap runtimes.
- Alt D can load WASI, but again moves too much normal runtime evolution into the
  installed binary.

Alt B preserves portable app execution while giving the native node enough host
integration to run those modules safely.

### Scenario 4: Dave runs on macOS installed by Homebrew

Dave installs `grid` through Homebrew. Later, stage0 fetches stage1 from CAS.

- Alt A still requires stage0 to contain a WASI loader. Homebrew installed
  stage0, but Homebrew is not the approval mechanism for fetched stage1 payloads.
- Alt B lets Homebrew deliver the stable stage0 binary while stage0 applies
  PromiseGrid checks to fetched stage1: exact CID verification, signer or local
  trust criteria, operator approval, executable-cache materialization, and
  platform checks such as codesign or notarization where required.
- Alt C adds sidecar approval and lifecycle questions.
- Alt D would make Homebrew updates carry more microkernel host behavior than
  needed.

Alt B makes the distinction between package-manager installation and
PromiseGrid root adoption explicit.

### Scenario 5: Ellen uses iOS or iPadOS as a control surface

Ellen uses an iPad to control or observe a node. Platform policy does not allow
ordinary arbitrary fetched executable behavior after install.

- Alt A does not solve iOS/iPadOS restrictions if the fetched module changes app
  behavior.
- Alt B treats iOS/iPadOS as bundled clients or control surfaces unless a later
  signed-app distribution path is explicitly designed.
- Alt C has the same platform restriction problem as Alt A.
- Alt D does not improve the iOS/iPadOS distribution model.

Alt B avoids pretending that WASM makes iOS/iPadOS normal fetched-stage1 targets.

### Scenario 6: Frank examines failure and recovery

Frank's node fetches a bad stage1 candidate, a candidate with incomplete closure,
or a candidate requesting new host capabilities.

- Alt A must distinguish WASI loader failure, stage1 failure, and missing native
  sidecar capability if hostful behavior is needed.
- Alt B can make the first proof concrete: descriptor fetch, executable fetch,
  CID verification, local approval, capability check, execution-cache
  materialization, process launch, readiness event, and failure event.
- Alt C doubles the failure surface.
- Alt D again makes host adapter failures stage0 failures.

Alt B gives the cleanest regression gates.

## Conclusion

Alt B, native/static stage1 first, is the surviving alternative. POC19 should
lock the first stage1 proof as a native/static executable object fetched and
verified by stage0. Stage0 remains a small bootstrap seed; it does not contain
the whole microkernel, a broad trust-policy engine, USB/device adapters, process
supervision logic, or a mandatory WASI loader for the first proof.

The locked implementation direction is:

- stage0 owns configured bootstrap roots, pinned public keys or trust anchors,
  exact CID verification, local approval recording, minimal PromiseGrid fetch
  handling, self-update, and native stage1 launch;
- stage0 fetches a stage1 descriptor and executable object by CID;
- stage0 verifies exact bytes, signer or local trust criteria, platform
  constraints, and local capability requirements;
- stage0 materializes a runnable copy into a grid-owned execution cache rather
  than executing directly from immutable CAS;
- stage0 starts native/static stage1 through the host process mechanism such as
  `execve`, `CreateProcess`, or Go `os/exec`;
- stage0 passes minimal bootstrap facts such as descriptor CID, adopted root CID,
  config path, state directory, and a local control endpoint;
- stage1 reports readiness, and stage0 records the result as a local event;
- stage1 owns hostful adapters and exposes narrow host capability promises to
  later WASI/WASM, OCI, and native app/runtime modules.

WASI/WASM is not rejected. It remains a first-class portable app/runtime profile
under stage1. The rejected claim is only that the first stage1 bootstrap proof
should itself be WASI.

## Implications for open TODOs and pending DIs

- `vumas.5` should scaffold one stable stage0 binary plus a native/static stage1
  descriptor/executable proof.
- `vumas.8` should still include WASI/WASM app execution, but through stage1
  rather than by requiring stage0 to load WASI.
- `vumas.9` should gate that stage0 does not require unverified stage1
  trust-policy code, that native stage1 launches from an execution cache, and
  that WASI/WASM remains a later runtime profile.
- POC20 root-decision records should continue to model stage0 self-update and
  runtime-root adoption as local timeline promises, not rollback or package
  manager authority.
- DEV-GUIDE-RESOURCES should state that native stage1 is first, while WASI/WASM
  remains under stage1.

## Decision status

Locked by `DI-topiv` in
`protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`.

## Refinements

### 2026-07-13 — Launch-attempt first proof refinement

`DI-tuvub` supersedes only the readiness-required launch-record wording in
`DI-topiv` and in this TE's implementation implications. The native/static
stage1 conclusion remains locked. The first POC19 proof now requires stage0 to
record a launch-attempt local event with descriptor CID, executable CID, adopted
Merkle/CAS root CID when present, execution-cache path, platform, approval or
rejection outcome, and process-launch outcome. Readiness is optional supplemental
information if the stage1 process reports it before timeout.
