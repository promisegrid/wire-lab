# POC19 production-shape design

## Status

Design draft. This is the first POC19 artifact. It is not executable code, not a
frozen protocol spec, and not a production API. Source: `DI-lumir`.

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
- keep the Promise Theory model visible: agents cooperate by making promises,
  remembering results, and deciding locally whom to trust.

POC19 is not a rewrite of PromiseGrid as a conventional package manager,
orchestration engine, command endpoint, central gatekeeper, or forge. It is a
production-shaped composition of lessons already tested in earlier POCs.

## Inheritance from earlier POCs

POC19 should be a strict successor to POC18 unless a later scoped DI explicitly
records a non-superset exception. The executable implementation should preserve
the useful POC18 behavior: sparse CAS, Rabin chunks, POSIX node promises,
reference sets, snapshots, Git bridge adapters, parent-linked exact messages,
TCP-carried object retrieval, CAR payloads, signed CWT/COSE tokens, diagnostic
CBOR rendering, and local sync-agent scheduling.

POC19 also inherits specific lessons from POC16 and POC17:

- From POC16: pCID owns arity, slot meaning, signable view, proof location, and
  payload interpretation. Parser/builder roles are local kernel roles, not
  global registries. Capability tokens are signed promises. Kernel roles are
  local promise surfaces for transport, lifecycle, storage, compute, device,
  key, app-interface, and resource-protection behavior.
- From POC17: constrained agents still use real pCID bytes in slot 0, binary CIDs
  on wire, CIDv1 base32 text when printable, compact payloads where the pCID
  spec permits them, and explicit resource limits. Small agents do not need
  heavy JSON-style self-description in every payload.
- From POC18: CAS/VCS state is the durable substrate; reference sets name
  directories, tags, branches, releases, logical changes, review threads,
  workspace roots, and app installations. Native collaboration is continuous
  peer DAG sync, not Git-style push/pull as the root model. Git import/export,
  push, and pull remain bridge adapters.

## Core design principles

### One binary, multiple local roles

The `grid` binary should contain both daemon and client behavior:

```text
                         same executable: grid

      +----------------------------------------------------------+
      |                                                          |
      |  CLI client modes              daemon/microkernel mode   |
      |  ----------------              -----------------------   |
      |  grid init                     grid daemon               |
      |  grid status                   transport roles           |
      |  grid snapshot                 pCID parser/builder roles |
      |  grid log                      CAS/VCS roles             |
      |  grid sync                     app runtime roles         |
      |  grid run                      lifecycle/resource roles  |
      |  grid git                      key/token roles           |
      |  grid diag                     local event journal       |
      |                                                          |
      +----------------------------------------------------------+
```

This does not mean the kernel is monolithic. A PromiseGrid kernel remains a
local role/profile set. The production-shaped improvement is packaging: one
binary can start the roles needed on a normal machine, while future deployments
can split those roles into separate processes, browser workers, mobile
sandboxes, firmware functions, or hosted services without changing the
PromiseGrid message model.

### Daemon/client shape

The long-running role is `grid daemon`. It owns the local node identity, local
CAS, repo/config discovery, peer connections, TCP/WebSocket listeners, app
runtime supervision, token redemption, local event journal, and local trust
state.

The CLI side is the same binary. When a daemon is running, CLI commands use the
same PromiseGrid message discipline over a local control stream. The first
implementation can use loopback TCP for that control stream because it is
portable and exercises the same framing as peer connections. Later Unix-socket,
Windows-named-pipe, browser, or mobile adapters are allowed if they preserve the
same pCID-selected message semantics.

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
| App interface | Resolve app reference sets, prepare local runtime inputs, deliver pCID-selected messages to local app processes or modules. |
| Execution runtime | Run WASI modules first; later run OCI containers and native binaries under stronger local resource promises. |
| Lifecycle/resource protection | Issue, narrow, revoke, or stop local CPU, memory, socket, storage, process, device, and time promises. |
| Key/token | Sign local promises, verify peer promises, issue CWT/COSE capability tokens, redeem local tokens, and reject replay locally. |
| Event journal | Record local events and exact artifact CIDs for later review without becoming a global monitor. |
| Trust policy | Choose peers, retention, forwarding, and execution willingness from local make/break history. |

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
and local-trust model. They differ only in the local runtime promises needed to
execute them safely.

### `grid run`

`grid run <app-ref>` means: ask the local daemon to evaluate whether it currently
promises to execute the app named by `<app-ref>`. The daemon may promise,
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
- Alice's daemon decides locally whether to execute the WASI module under its
  resource promises.
- The module's outputs are CAS objects and local event records, not hidden
  daemon state.

### Flow 2: installing an app through VCS/CAS

```text
Carol creates:

  app reference_set
    -> runtime profile
    -> module CIDs
    -> data root CIDs
    -> pCID specs
    -> entrypoint record
    -> resource profile

Carol signs that reference_set promise.

Alice receives it through normal peer sync.
Alice may:
  - retain it;
  - inspect it;
  - run it;
  - map it into her local namespace;
  - refuse to run it;
  - ask peers for missing objects.
```

There is no separate installation root of truth. Installation is local adoption
of a signed app reference set and its reachable CAS closure.

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
| `grid init` | Create repo-local `.grid` config and choose a file or daemon CAS locator. |
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
- Alice may revoke a runtime capability if an app exceeds CPU, memory, socket,
  time, or storage terms.

None of these decisions are global. They are Alice's local relationship and
resource judgments.

## Production-shaped boundaries

### What should become real in POC19

- One `grid` binary with daemon and client modes.
- Real TCP and WebSocket transport adapters.
- Real sparse CAS shared across repos/apps on one node.
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

1. **Scaffold.** Create `implementations/poc19-production-shape/` with one
   `grid` binary and shared packages. Keep POC18 code available as the source
   baseline.
2. **Daemon baseline.** Move POC18 local CAS/VCS and TCP agent behavior behind
   `grid daemon`, with CLI commands using the daemon for normal operation.
3. **Transport parity.** Add WebSocket adapter that carries the same exact
   `grid()` bytes as TCP and proves parity in a clean run.
4. **App reference sets.** Add role `app` reference sets and resolve app code,
   data, pCID specs, runtime profile, and entrypoint labels from CAS.
5. **WASI run.** Implement `grid run` for a fetched WASI module with CID-verified
   inputs and CAS-recorded outputs.
6. **Capability promises.** Add local runtime capability tokens for CPU, memory,
   time, storage, network, and host-function promises.
7. **Container/native profiles.** Store and verify OCI image graphs and native
   binaries, then execute only under explicit local runtime promises.
8. **Regression gates.** Prove POC18 superset behavior, promise-shaped
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
- **Transport drifting into app semantics.** Mitigation: TCP/WebSocket adapters
  carry exact bytes; pCID parser/builder roles interpret protocol semantics.
- **Resource protection vocabulary drifting back to control language.**
  Mitigation: describe every runtime grant as a conditional capability promise
  made by the local resource owner.
- **Observer/analyzer assumptions leaking into production.** Mitigation: POC19
  uses a local event journal; any analyzer remains regression machinery.

## Acceptance criteria for the future executable POC19

The future implementation should not be considered complete until a clean run
shows:

- one `grid` binary starts a daemon and serves CLI commands;
- app code and data are fetched from peer CAS over TCP and WebSocket using exact
  PromiseGrid messages;
- at least one WASI app is installed through an app reference set and run from
  VCS/CAS;
- outputs are CAS objects with CIDs returned to the user;
- local runtime capability promises are issued and checked;
- sparse CAS remains partial and peer-relative;
- POC18 VCS behavior still works;
- diagnostics can render representative raw CBOR messages;
- regression checks prove no inter-agent path used simulated sideband transfer.
