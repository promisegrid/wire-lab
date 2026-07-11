# POC19 production-shape design

## Status

Design draft. This is the first POC19 artifact. It is not executable code, not a
frozen protocol spec, and not a production API. Source: `DI-lumir`; `DI-kodob`;
`DI-guhil`; `DI-zitap`.

## Purpose

POC19 should make PromiseGrid feel less like a harness and more like something a
developer could drop onto a machine. The target is a single binary named `grid`
that can:

- run as a local PromiseGrid daemon/microkernel;
- expose the user-facing VCS and app CLI;
- fetch apps, modules, binaries, container images, and data from peers over TCP
  or WebSocket;
- store everything in sparse CAS and VCS/reference-set state;
- execute fetched code with fetched data under local promises and capability
  tokens;
- act as the stable bootstrap seed for the minimum microkernel, so microkernel
  role modules, app modules, agent modules, runtime executables, specs, and data
  are fetched by CID rather than shipped as ordinary binary updates;
- keep the Promise Theory model visible: agents cooperate by making promises,
  remembering results, and deciding locally whom to trust.

POC19 is not a rewrite of PromiseGrid as a conventional package manager,
orchestration engine, command endpoint, central gatekeeper, or forge. It is a
production-shaped composition of lessons already tested in earlier POCs.

POC20 now owns a parallel semantic-model track for promises as timeline
assertions, deterministic pure-function agents, CAS-backed local/group
timelines, and branch-aware capability-token double-spend behavior. POC19 should
continue as the production-shaped plumbing path, but it should avoid choices
that would make POC20's visible branch histories impossible. Source: `DI-kakos`;
`TE-lodom`.

## Inheritance from earlier POCs

POC19 should be a strict successor to POC18 unless a later scoped DI explicitly
records a non-superset exception. This inheritance review is part of the POC19
implementation contract: code generation should not begin until each inherited
lesson is either covered by this design or assigned to a later POC19 task.
Source: `DI-lumir`; `DI-topab`.

Reviewed anchors: POC16 `README.md`, `docs/MESSAGE-SHAPES.md`,
`docs/KERNEL-ROLES.md`, `docs/ROUTE-PROMISES.md`, and protocol specs; POC17
`README.md`, `CHANGELOG.md`, and protocol specs; POC18
`docs/protocols/version-control.md`, `cmd/grid`, `cmd/poc-sim`,
`cmd/poc-agent`, `cmd/poc-analyze`, and `protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md`.

| Source | Inherited lesson | POC19 coverage | Gap status | Pre-code follow-up |
| --- | --- | --- | --- | --- |
| POC16 | pCID owns envelope arity, slot meaning, signable view, proof location, and payload interpretation; pCID is not a peer address, app address, operation, route, repository name, or message type. | Core design keeps `grid([42(pCID), ...protocol-defined-slots])` and puts route/app/operation semantics inside pCID-defined payloads. | covered | Preserve this as a regression check in `vumas.9`. |
| POC16 | Parser/builder roles are local kernel roles that receive exact slot-0 pCID bytes and deliver parsed protocol messages to apps. | Local daemon roles include a pCID parser/builder role separate from transport and app interface roles. | partially covered | `vumas.4` must lock whether this role is factored from POC18 code or implemented in a new shared core before `vumas.5` scaffolding. |
| POC16 | Protocol specs are first-class; printable pCIDs are CIDv1 base32 text and wire pCIDs are binary CID bytes. | App reference sets include pCID specs; object identity section preserves binary-on-wire/base32-printable discipline. | partially covered | `vumas.4` must produce the POC19 pCID inventory before code generation. |
| POC16 | Capability tokens are signed promises; local lifecycle, resource access, storage, and retrieval tokens use the CWT/COSE pattern rather than custom unsigned fields. | Design names CWT/COSE capability tokens for retrieval, storage, and runtime resources. | partially covered | `vumas.8` must define the app/runtime token profile before `grid run` executes fetched code. |
| POC16 | Encrypted payloads and COSE payload/proof variants are pCID-owned message shapes, not universal envelope rules. | Design preserves pCID-owned slot semantics but does not yet enumerate encrypted app-data profiles. | partially covered | Assign encrypted app-input and app-output profiles to `vumas.8`; keep any proof slot pCID-defined. |
| POC16 | Kernel roles are non-monolithic local promise surfaces for transport, lifecycle, storage, compute, device, key, app-interface, and resource-protection behavior. | Local daemon roles explicitly separate transport, parser/builder, CAS/VCS, app interface, execution runtime, local event journal, and key/token behavior. | covered | Keep role boundaries visible in `vumas.5` scaffold and `vumas.9` gates. |
| POC16 | Exact raw CBOR messages and sparse per-agent CAS are reviewable after a run. | Design keeps exact `grid()` bytes as durable boundary objects and requires diagnostics for raw CBOR/CAS objects. | covered | Preserve exact-message retention in `vumas.9`. |
| POC17 | Constrained agents use compact pCID-selected payloads, real pCID bytes in slot 0, binary CIDs on wire, and base32 CIDs when printable. | Design keeps compact pCID-owned payloads and the CID discipline across transports and app reference sets. | covered | Add a constrained-message fixture to `vumas.9` even if POC19's first runtime is not LoRa. |
| POC17 | Small agents should not need JSON-style self-description in every payload; the pCID spec carries the shared grammar. | Design avoids universal self-describing payload maps and keeps payload shape pCID-defined. | covered | Reject new generic map envelopes unless a pCID spec requires them. |
| POC17 | Radio and embedded slices expose real resource limits and incomplete peer storage; no peer is assumed to hold the full DAG. | Storage model states that CAS stores are partial and peer-relative. | covered | Include missing-object and partial-CAS cases in `vumas.9`. |
| POC18 | CAS/VCS is the durable substrate: Rabin chunks, POSIX node promises, reference sets, tags, snapshots, logical changes, review statements, and parent-linked exact messages. | Design makes apps and user data VCS/CAS objects and preserves POC18 command surface. | covered | Keep POC18 VCS behavior in `vumas.9` superset gates. |
| POC18 | `.grid` repo state remains the user-facing local repository shape while allowing a daemon-owned node CAS. | Storage model preserves `.grid/config.json`, `.grid/state.json`, and a configurable CAS locator. | covered | Implement in `vumas.6`. |
| POC18 | Native collaboration is continuous peer DAG sync over promise-shaped TCP messages; Git import/export, push, and pull are bridge adapters. | Network model keeps native fetching as peer DAG sync and treats Git push/pull as interoperability. | covered | Preserve no-sideband inter-agent transfer checks in `vumas.9`. |
| POC18 | Object retrieval uses signed CWT/COSE capability tokens and may transfer CAR payloads when a peer redeems a retrieval promise. | Design names retrieval/storage tokens and object transfer but has not frozen the durable object layout. | partially covered | Resolve token and object-transfer details under `vumas.8` and regression checks under `vumas.9`. |
| POC18 | Raw chunks, DAG-CBOR, GRID-CBOR, CAR files, and local store layout remain open storage-profile decisions. | Storage model explicitly calls this unresolved. | pre-code blocker | Resolve the object storage profile during `vumas.4` before durable POC19 stores are scaffolded. |
| POC18 | CLI commands must use shared core behavior rather than a parallel automation path. | Command surface and daemon/client sections require a shared core and daemon-backed normal operation. | covered | `vumas.4` must choose the factoring path that prevents CLI/core duplication. |

## Core design principles

### One binary, multiple local roles

The `grid` binary should contain both daemon and client behavior:

```text
                  stable installed executable: grid bootstrap seed

      +----------------------------------------------------------+
      |                                                          |
      |  built in:                    fetched microkernel modules|
      |  ---------                    ---------------------------|
      |  local config                 CLI/daemon surfaces        |
      |  CID verification             transport roles            |
      |  minimal grid() fetch         pCID parser/builder roles  |
      |  owner approval record        CAS/VCS and sync roles     |
      |  self-update restart          runtime/resource roles     |
      |                               key/token/event/trust roles|
      |                                                          |
      +----------------------------------------------------------+
```

This does not mean the kernel is monolithic. A PromiseGrid kernel remains a
local role/profile set. The production-shaped improvement is packaging: one
binary can bootstrap the roles needed on a normal machine by fetching the
microkernel modules named by approved CIDs, while future deployments can split
those roles into separate processes, browser workers, mobile sandboxes, firmware
functions, or hosted services without changing the PromiseGrid message model.

### Stable `grid` binary as the minimum microkernel

The installed `grid` binary is the stable bootstrap seed for the minimum
microkernel, not the whole microkernel and not the application distribution unit.
It should contain very little functionality and change rarely. Its first promise
is to provide just enough local identity/config handling, CID verification,
minimal PromiseGrid `grid()` message handling, peer fetch, local approval
recording, and restart handoff to bootstrap everything else.

On startup, the stable binary may detect a candidate new version of itself from
locally trusted peers. If an owning agent or human locally approves that
self-update, `grid` fetches the new binary by CID, verifies the exact bytes and
local signer trust criteria, installs the approved version, and restarts. After
restart, it fetches the remaining modules that make up the minimum microkernel:
transport adapters, pCID parser/builders, CAS/VCS, sync, app interface,
execution runtime, lifecycle/resource protection, host-capability, secret
service, key/token, local event journal, and trust-policy modules.

Application-specific modules come after that microkernel bootstrap. App, agent,
runtime executable, protocol-spec, and data changes are CID-addressed CAS/VCS
objects fetched from peers through PromiseGrid `grid()` messages and the CIDs
they name, rather than shipped as ordinary binary updates. The app bytes
themselves remain exact CAS objects. Source: `DI-zitap`; `DI-kodob`.

First run or local config may provide one bootstrap root CID. That root names a
Merkle/CAS object graph containing app reference sets, runtime profiles,
protocol specs, executable objects, data roots, and update metadata. The daemon
fetches the graph from peers, verifies exact CIDs, and asks for local operator
approval before adopting it as a local root. A later code or data update is a new
root CID plus a local approval promise, not a replacement of the `grid` binary.

The binary may still change for true substrate changes: transport support,
storage-profile changes, security fixes, parser/loader bugs, or other
microkernel behavior. Normal app, agent, runtime, and data changes should move
through CAS roots. Source: `DI-kodob`.

### Operator-visible root adoption

Root adoption is an explicit local operator flow, not an implied daemon side
effect. A candidate root is fetched by CID from configured peers, then evaluated
before it becomes the local adopted root. The evaluation must verify the
candidate closure, exact CIDs, local signer trust criteria, protocol/spec CIDs,
and requested host capabilities.

Before adoption, the daemon must produce a concise impact summary as a CAS object:
changed apps, agents, runtime profiles, protocol specs, policy bundles, signers,
current root, candidate root, requested capability changes, and rollback root.
The operator or local approving role then makes a local approval promise, local
rejection promise, or local still-evaluating promise. A later rollback is another
local root-adoption promise that returns to a retained prior root without
replacing the installed `grid` binary.

Replay must be able to identify which adopted runtime root originally produced a
retained artifact, even if that artifact is later displayed or interpreted
through a newer root. Source: `DI-guhil`.

### Daemon/client shape

The long-running role is `grid daemon`. It owns the local node identity, local
CAS, repo/config discovery, peer connections, TCP/WebSocket listeners, app
runtime supervision, token redemption, local event journal, and local trust
state.

The user still invokes `grid`, but the stable bootstrap binary may load or fetch
the CLI and daemon microkernel modules before those surfaces are available. When
a daemon module is running, CLI commands use the same PromiseGrid message
discipline over a local control stream. The first implementation can use loopback
TCP for that control stream because it is portable and exercises the same framing
as peer connections. Later Unix-socket, Windows-named-pipe, browser, or mobile
adapters are allowed if they preserve the same pCID-selected message semantics.

The CLI should not grow a parallel implementation of graph, CAS, sync, or app
execution behavior. It should call the shared core directly only for bootstrap
or offline read-only operations where no daemon is available. Normal production
operation is daemon-backed.

### Everything crosses boundaries as promises

Every inter-agent message is a PromiseGrid message:

```text
grid([42(pCID), ...protocol-defined-slots])
```

Slot 0 is a protocol CID. It is not an agent address, app address, file path,
route, operation, command name, repository name, or message type. The pCID
chooses the protocol parser and slot grammar. Payloads or nested payloads carry
the local routing, app, resource, or VCS meaning defined by that protocol.

An implementation may expose ergonomic local methods, but those methods are
adapters. The durable boundary object is the exact `grid()` message and any CAS
objects it names.

### Apps are VCS/CAS objects

Non-kernel code is fetched from VCS/CAS. A grid app is installed by checking a
signed app reference set into VCS/CAS. There is no separate package-manager root
model for POC19.

An app reference set is a normal PromiseGrid reference set with role `app`. It
labels the app's code, runtime profile, input data roots, pCID specs, entrypoint
metadata, resource profile, and expected reciprocal promises.

App reference sets are usually reachable from an operator-adopted root CID. The
root is a convenience and update anchor, not a global authority. Alice may adopt
one root, Bob may adopt another, and a group may voluntarily converge on a shared
root by making reciprocal promises to track it. Source: `DI-kodob`.

Example labels in an app reference set:

```text
reference_set role: app

label                  target role             target
-----                  -----------             ------
runtime                runtime_profile         42(wasi_profile_cid)
module/main            executable_code         42(wasm_module_cid)
module/helper          executable_code         42(wasm_module_cid)
data/default           input_data_root         42(snapshot_or_refset_cid)
protocols/app          protocol_spec           42(app_protocol_pcid)
resources/requested    resource_profile        42(resource_profile_cid)
entrypoint/default     entrypoint_record       42(entrypoint_record_cid)
```

The reference set is the install object. A later design may add app-manifest
objects if they prove useful, but POC19 should not require a second app package
format before the reference-set model has been tested.

## Runtime architecture

### Local daemon roles

The daemon should make these local promises explicit:

| Role | Local promise surface |
|---|---|
| Transport | Listen, dial, frame, read, write, close, and remember direct TCP/WebSocket message results. |
| pCID parser/builder | Route exact slot-0 pCID bytes to local protocol parsers/builders without treating pCID as a destination or operation. |
| CAS/VCS | Store, retrieve, verify, retain, garbage-collect, and materialize sparse CAS/VCS objects. |
| Sync | Advertise selected reference sets, request missing parents, redeem retrieval tokens, and exchange CAR/object bytes. |
| Bootstrap/root adoption | Fetch root CID closures, verify exact CIDs, present adoption choices to the operator, and remember local root-adoption promises. |
| App interface | Resolve app reference sets, prepare local runtime inputs, deliver pCID-selected messages to local app processes or modules. |
| Execution runtime | Run WASI modules first; later run OCI containers and native binaries under stronger local resource promises. |
| Lifecycle/resource protection | Issue, narrow, revoke, or stop local CPU, memory, socket, storage, process, device, and time promises. |
| Host capability | Promise narrow access to local network targets, files, directories, devices, host functions, secret references, execution resources, and storage only when local terms are satisfied. |
| Secret service | Sign bytes, unwrap scoped keys, mint short-lived credentials, rotate or revoke material, and record denied attempts without exposing plaintext secrets as ordinary CAS payloads. |
| Key/token | Sign local promises, verify peer promises, issue CWT/COSE capability tokens, redeem local tokens, and reject replay locally. |
| Event journal | Record local events and exact artifact CIDs for later review without becoming a global monitor. |
| Trust policy | Choose peers, retention, forwarding, and execution willingness from local make/break history. |

### Host capability promises and secret services

Fetched code does not widen local resource promises by naming a dependency,
policy object, or requested adapter. Every bridge to a local resource is a narrow
capability promise made by the local node. POC19 should model at least these
classes: named local network access, file or directory access, device access,
host-function access, secret-reference use, CPU, memory, time, sockets, process
count, and storage.

If a candidate root requests a capability that the local node has not already
promised, adoption must pause, fail locally, or require a separate local approval
promise before that root can become active. The root may still be retained in CAS
as a candidate without being adopted.

Secrets are not ordinary CAS payload. Config, CAS, diagnostics, logs, prompts,
and UI output may contain secret references, public key material, fingerprints,
policy CIDs, revocation events, and non-secret local events, but not plaintext
private keys, passphrases, API tokens, unwrapped data keys, or break-glass
credentials. A local secret service can promise scoped operations such as signing
bytes, unwrapping a scoped data key, minting a short-lived credential, rotating
material, revoking material, or recording a denied attempt. The response contains
only the narrow result needed for the operation. Source: `DI-guhil`.

### App execution profiles

POC19 should define three execution profiles but implement them in risk order.

1. **WASI profile.** First executable profile. Code is a WASM module stored by
   CID. Inputs are CAS objects, scalar arguments, or pCID-defined streams. The
   daemon provides only promised host functions: read named CAS inputs, write
   result objects, emit app messages, and use bounded clock/randomness if those
   are explicit context promises.
2. **OCI/container profile.** Code is an image or image-layer graph stored in
   CAS. Execution requires stronger local lifecycle/resource promises for CPU,
   memory, filesystem view, network, devices, and cleanup. Image layers are CAS
   objects, not external pulls hidden from PromiseGrid.
3. **Native-binary profile.** Code is a host binary stored in CAS. This is the
   highest-risk profile because host ABI, OS, CPU, filesystem, dynamic libraries,
   and process behavior matter. A daemon may refuse native execution while still
   storing and syncing the object graph.

The design goal is that all profiles use the same install, fetch, promise, CAS,
operator-adopted root, and local-trust model. They differ only in the local
runtime promises needed to execute them safely.

### `grid run`

`grid run <app-ref>` means: ask the local daemon to evaluate whether it currently
promises to execute the app named by `<app-ref>`. The app reference may be named
directly or found through an operator-adopted root CID. The daemon may promise,
decline, delay, or request missing objects from peers. If it runs the app, the
result is another CAS/VCS object or reference set, not a hidden local side
effect.

The output should include stable CIDs:

```text
app_ref=<cid>
run_record=<cid>
stdout=<cid or empty>
stderr=<cid or empty>
result=<cid or empty>
events=<cid>
outcome=kept|not_promised|broken|unavailable
```

The exact field names are implementation details for POC19, but the semantic
shape is fixed: the run returns local promise results and CID-addressed objects.

## Network model

### TCP and WebSocket parity

POC19 should treat TCP and WebSocket as equal transport targets. Both carry the
same exact PromiseGrid message bytes. WebSocket frames should use binary frames
for exact CBOR bytes, not JSON wrappers around protocol messages.

```text
                  exact PromiseGrid CBOR bytes
                             |
          +------------------+------------------+
          |                                     |
      TCP adapter                         WebSocket adapter
 length-framed stream                     binary WS frames
          |                                     |
       peer A                                peer B
```

Transport adapters make hop-local promises about delivery attempts, session
opening, session closing, and byte preservation. They do not decide app
semantics. The pCID parser/builder role receives exact bytes after transport.

### Fetch and sync flow

Native fetching is promise-shaped peer DAG sync. It should continue the POC18
message family:

1. `sync_interest`: Alice promises she is willing to receive selected missing
   objects and verify exact CIDs.
2. `object_availability`: Bob promises selected objects are available under
   stated terms and may include a signed retrieval capability token.
3. `object_retrieval_redemption`: Alice redeems Bob's token and promises to
   evaluate the returned bytes locally.
4. `object_bytes`: Bob sends CAR/object bytes and promises they match the named
   CIDs.
5. `object_retention`: Alice, Bob, or a third peer may promise local retention
   under local terms.

The design should not reintroduce native Git push/pull as the core
synchronization model. Git push/pull remains a compatibility bridge to
conventional Git repositories.

## Storage model

### Node CAS and repo state

POC18 stores repo-local config and state under `.grid/`. POC19 should preserve
that user-facing shape while allowing the CAS locator to point at the daemon:

```text
workspace/
  .grid/
    config.json      # repo-local config and CAS locator
    state.json       # local VCS convenience state
```

The daemon owns the node's main sparse CAS. A single node CAS may serve many
repos and apps. Repo-local `.grid/config.json` should be able to point to:

- local file CAS for bootstrap or offline use;
- local daemon CAS for normal production-shaped use;
- later remote or delegated CAS profiles if a pCID and DI justify them.

CAS stores are partial. No node promises to have every object. Availability,
retention, and serving are separate promises.

### Object identity

POC19 should retain the current CID discipline:

- binary CIDs on wire;
- CIDv1 base32 text with `b` prefix when printable;
- exact message CIDs for parent-linked messages;
- CID-addressed app code, WASI modules, container layers, native binaries, data
  objects, app reference sets, run records, and result objects.

Open POC18 decisions about raw chunks, DAG-CBOR, GRID-CBOR, CAR files, and store
layout should be resolved before POC19 implementation freezes durable object
profiles. The POC19 design can proceed now because those decisions affect storage
layout, not the high-level production-shaped architecture.

## Example flows

### Flow 1: `grid run` fetches and runs a WASI app

```text
Alice CLI            Alice daemon             Bob daemon              WASI module
   |                      |                       |                       |
   | grid run app_ref     |                       |                       |
   |--------------------->|                       |                       |
   | local run promise    |                       |                       |
   |                      | resolve adopted root  |                       |
   |                      | missing module/data   |                       |
   |                      | sync_interest         |                       |
   |                      |---------------------->|                       |
   |                      | object_availability   |                       |
   |                      |<----------------------|                       |
   |                      | redemption            |                       |
   |                      |---------------------->|                       |
   |                      | object_bytes          |                       |
   |                      |<----------------------|                       |
   |                      | verify CIDs/proofs    |                       |
   |                      | issue local runtime capability promises        |
   |                      |---------------------------------------------->|
   |                      | execute with CAS inputs                         |
   |                      |<----------------------------------------------|
   |                      | store result/run-record CIDs                   |
   | run_record/result    |                       |                       |
   |<---------------------|                       |                       |
```

Important properties:

- Bob never promises Alice's code is safe to run. Bob only promises object
  availability and exact bytes under his local terms.
- Alice decides locally whether Bob is trusted enough as a source.
- Alice's adopted root identifies which app graph she is willing to consider,
  but it does not force execution.
- Alice's daemon decides locally whether to execute the WASI module under its
  resource promises.
- The module's outputs are CAS objects and local event records, not hidden
  daemon state.

### Flow 2: adopting an app/runtime root through VCS/CAS

```text
Carol creates:

  update root reference_set
    -> app reference_set
    -> runtime profile
    -> pCID specs
    -> update notes

  app reference_set
    -> runtime profile
    -> module CIDs
    -> data root CIDs
    -> pCID specs
    -> entrypoint record
    -> resource profile

Carol signs the root and app reference-set promises.

Alice receives it through normal peer sync.
Alice may:
  - retain it;
  - inspect it;
  - approve it as her local adopted root;
  - run it;
  - map it into her local namespace;
  - refuse to run it;
  - ask peers for missing objects.
```

There is no global installation root of truth. Installation is local adoption of
a signed root or app reference set and its reachable CAS closure. Operator
approval is a local promise by Alice, not a command to other nodes. Source:
`DI-kodob`.

### Flow 3: WebSocket peer

```text
Browser/mobile peer        Alice daemon
       |                       |
       | binary WS grid bytes  |
       |---------------------->|
       | binary WS grid bytes  |
       |<----------------------|
```

The WebSocket adapter changes only framing and reachability. The pCID-selected
message bytes, CAS objects, CIDs, token checks, and local trust decisions are the
same as TCP.

## Command surface

POC19 should preserve and extend the POC18 `grid` surface:

| Command | Production-shaped meaning |
|---|---|
| `grid daemon` | Start the local long-running PromiseGrid role set. |
| `grid init` | Create repo-local `.grid` config, choose a file or daemon CAS locator, and optionally record a bootstrap root CID. |
| `grid snapshot` | Record workspace state as CAS/VCS promises. |
| `grid status` | Compare workspace state against local recorded VCS state. |
| `grid log` | Walk local parent-linked snapshot history. |
| `grid track` / `grid untrack` | Mutate local path inclusion policy, not Git staging. |
| `grid sync` | Inspect or run local peer sync through the daemon. |
| `grid run` | Execute an app reference set from VCS/CAS under local promises. |
| `grid git` | Bridge import/export/push/pull with conventional Git repositories. |
| `grid diag` | Render exact CBOR/CAS objects for local review. |

The command surface should avoid implying a central repository or global
namespace. It is acceptable for commands to be familiar, but their meaning is
PromiseGrid-native.

## Local trust and resource decisions

The daemon keeps local records of promises made, promises relied on, results,
unavailable peers, broken reciprocal terms, token redemption, replay rejection,
and resource exhaustion. These records are local inputs to future decisions.

Examples:

- Alice may fetch app objects from Bob because Bob has kept storage promises.
- Alice may avoid Mallory because Mallory has repeatedly advertised unavailable
  objects.
- Alice may run a WASI app but refuse a native-binary profile for the same app.
- Alice may keep an app installed but decline to map it into a convenient local
  namespace.
- Alice may approve a new app/runtime root CID for future runs or keep the old
  root if the new root's source, terms, or closure look weak.
- Alice may revoke a runtime capability if an app exceeds CPU, memory, socket,
  time, or storage terms.

None of these decisions are global. They are Alice's local relationship and
resource judgments.

## Production-shaped boundaries

### What should become real in POC19

- One `grid` binary with daemon and client modes.
- Real TCP and WebSocket transport adapters.
- Real sparse CAS shared across repos/apps on one node.
- Real bootstrap from an operator-provided root CID.
- Real app reference-set resolution from VCS/CAS.
- Real WASI execution from fetched module bytes.
- Real CWT/COSE capability tokens for retrieval, storage, and local runtime
  resources.
- Real exact-message retention and diagnostic rendering.
- Real local event journal that is useful without a global observer.

### What remains POC-local or later

- Production key management and identity recovery.
- Full OCI/native sandboxing hardening.
- Final store profile choices for DAG-CBOR, GRID-CBOR, raw chunks, and CAR
  storage.
- Cross-platform service installation.
- Rich UI, app marketplace, or group namespace UX.
- Final external PromiseGrid Development Guide prose.

## Implementation phases after this design

1. **Scaffold.** Create `implementations/poc19-production-shape/` with one stable
   `grid` bootstrap binary, a fetched microkernel-module set, and shared
   packages. Keep POC18 code available as the source baseline.
2. **Daemon baseline.** Move POC18 local CAS/VCS and TCP agent behavior behind
   `grid daemon`, with CLI commands using the daemon for normal operation.
3. **Transport parity.** Add WebSocket adapter that carries the same exact
   `grid()` bytes as TCP and proves parity in a clean run.
4. **Bootstrap roots.** Add config/first-run support for an operator-provided
   root CID, candidate-root closure verification, impact summaries, local
   approval promises, rollback roots, and replay context.
5. **App reference sets.** Add role `app` reference sets and resolve app code,
   data, pCID specs, runtime profile, and entrypoint labels from CAS.
6. **WASI run.** Implement `grid run` for a fetched WASI module with CID-verified
   inputs and CAS-recorded outputs.
7. **Capability promises.** Add local runtime capability tokens for CPU, memory,
   time, storage, network, filesystem, device, secret-reference, and host-function
   promises.
8. **Secret service.** Add operation-scoped local secret service promises for
   signing, unwrap, short-lived credential minting, rotation, revocation, and
   denial events without placing plaintext secrets in broad CAS/log/prompt/UI
   paths.
9. **Container/native profiles.** Store and verify OCI image graphs and native
   binaries, then execute only under explicit local runtime promises.
10. **Regression gates.** Prove POC18 superset behavior, promise-shaped
   inter-agent traffic, TCP/WebSocket parity, exact-CBOR diagnostics, and local
   event records.

## Design risks

- **Single binary becoming monolithic.** Mitigation: keep local roles explicit
  and testable even when packaged together.
- **CLI reimplementing core logic.** Mitigation: one shared core library;
  daemon-backed normal operation; direct CLI logic only for bootstrap/offline
  reads.
- **Execution drifting into package-manager behavior.** Mitigation: apps are
  reference sets in VCS/CAS; all code/data objects are CID-addressed.
- **Minimum binary drifting into app bundle.** Mitigation: regression gates should
  prove app/runtime updates change adopted root CIDs and CAS objects, not the
  installed `grid` binary.
- **Root adoption becoming invisible.** Mitigation: root updates must produce
  candidate-root impact summaries, explicit local approval events, retained
  rollback roots, and replayable original-root context.
- **Fetched code widening local resources.** Mitigation: every local bridge is a
  narrow host-capability promise, and new requested capabilities pause, fail, or
  require separate local approval before adoption.
- **Secrets leaking into CAS or diagnostics.** Mitigation: use operation-scoped
  secret services and keep plaintext secrets out of config, CAS, logs, prompts,
  diagnostics, and UI output.
- **Transport drifting into app semantics.** Mitigation: TCP/WebSocket adapters
  carry exact bytes; pCID parser/builder roles interpret protocol semantics.
- **Resource protection vocabulary drifting back to control language.**
  Mitigation: describe every runtime grant as a conditional capability promise
  made by the local resource owner.
- **Observer/analyzer assumptions leaking into production.** Mitigation: POC19
  uses a local event journal; any analyzer remains regression machinery.
- **POC19 hiding semantic timeline state.** Mitigation: keep token ledgers,
  app-run records, parent links, and CAS object transfers explainable as durable
  objects so POC20 can later test branch-based timeline semantics without
  undoing POC19 storage and runtime choices. Source: `DI-kakos`; `TE-lodom`.

## Acceptance criteria for the future executable POC19

The future implementation should not be considered complete until a clean run
shows:

- one stable `grid` bootstrap binary can fetch microkernel modules, then start a
  daemon module and serve CLI commands;
- app code and data are fetched from peer CAS over TCP and WebSocket using exact
  PromiseGrid messages;
- first run or config can name a bootstrap root CID, and later app/runtime
  updates are adopted by operator-approved root CIDs without replacing the
  `grid` binary;
- candidate roots with incomplete closures, locally untrusted signatures, or
  unapproved requested host capabilities do not become active roots silently;
- root adoption produces an impact summary CID and rollback root, and rollback
  returns to the prior root without replacing the installed binary;
- retained artifacts can be replayed with the runtime root that originally
  produced them;
- plaintext secrets do not appear in config, CAS, diagnostics, logs, prompts, or
  UI output;
- at least one WASI app is installed through an app reference set and run from
  VCS/CAS;
- outputs are CAS objects with CIDs returned to the user;
- local runtime capability promises are issued and checked;
- sparse CAS remains partial and peer-relative;
- POC18 VCS behavior still works;
- diagnostics can render representative raw CBOR messages;
- regression checks prove no inter-agent path used simulated sideband transfer.
