# PromiseGrid kernel role/profile design

This design note explains the current PromiseGrid kernel direction in plain
English. It is not a TE, not a frozen protocol spec, and not final SDK prose. It
is the current design synthesis for resolving `DR-davod`: what guide writers can
say about kernel developers after the kernel TEs converged, and what still waits
for focused evidence and future frozen pCID specs. Source: `DI-fidot`;
`DI-punuf`; `DI-gopag`.

## Basic principles

### Everything useful is a promise

PromiseGrid protocol objects help agents make, recognize, remember, and evaluate
promises. Files, messages, graph nodes, services, executables, path entries, byte
chunks, app/kernel operations, port records, namespace entries, and evidence
records are useful because they are promises, evidence of promises, or durable
objects that help an agent evaluate promises.

The phrase should not be flattened into "everything is a file" or "everything is
a message." File-like views and messages are important, but promise is the
semantic unit.

## Layer discipline

The current layer model is:

- **Promise:** the semantic unit.
- **Object/event/log/checkpoint:** the durable history unit.
- **pCID-selected message:** the boundary-crossing unit.
- **File-like view:** the user and developer experience projection.

This lets PromiseGrid keep a mainframe-like developer experience without
pretending that one filesystem, runtime, or administrator owns the truth.

## Decentralized mainframe

PromiseGrid can be described as a decentralized mainframe if the metaphor is
used carefully.

The useful part of the mainframe analogy is that Alice, Bob, and Carol can share
data, executable resources, devices, and computation results as nearby named
resources rather than as unrelated downloads, URLs, installers, and services.
Dave can open or run a resource in a coherent local view, and that view can feel
like a shared environment.

The PromiseGrid difference is that coherence comes from promises and local
trust. There is no central namespace authority, global trust oracle, global
permission system, or universal runtime. Each agent builds its own view from
promises it can observe, recognize, and locally choose to trust.

Trusted groups can still maintain shared namespaces. Alice, Bob, and Carol may
all promise to publish, update, remember, and audit the same project namespace.
That namespace is real for their relationship group because each agent keeps
reciprocal promises and records evidence. It is not universal truth for everyone
else.

## Kernel means role/profile set

"Kernel" means a local PromiseGrid infrastructure role set that helps an agent
speak pCID-selected protocols and keep evidence about promises. It does not mean
one universal process shape.

A deployment may realize the role set as:

- a daemon on a rich native node;
- a microkernel-style service graph;
- a browser or WASM host adapter;
- a mobile sandbox profile;
- a tiny MCU/header-only library;
- a split local mesh of services;
- a future runtime shape not yet designed.

These are profiles. They are not competing definitions of PromiseGrid itself.

The kernel role set may include dispatch, transport, app lifecycle, storage,
compute, key use or signing, hardware/resource access, pCID handling, evidence
recording, namespace projection, and reference resolution. These are roles
because each one is a bounded local promise surface: what the role promises to
do, what it does not promise, what evidence it records, and what host/runtime
assumptions it depends on. Source: `DI-fidot`; `DI-punuf`.

In a split deployment, different local agents, processes, objects, adapters, or
library functions may make these promises separately. POC12's `printer_port`
role is the current concrete example: it owns a simulated local hardware
resource, promises scoped future printing by issuing a capability-promise token,
and later promises print evidence when the token is redeemed with bounded label
bytes. The message kernel only carries the pCID-selected bytes and records its
own transport evidence; it does not become a USB authority, permission server,
business workflow engine, or trust oracle. Source: `DI-pohaj`; `DI-vutok`;
`DI-punuf`.

## How to read the current packet

The current kernel evidence packet has three distinct roles:

- `DN-lujad` is this design note. It is the plain-English synthesis of the
  kernel role/profile model for guide writers and implementers.
- `DR-davod` is the open decision request. It owns the unresolved question of
  what a stable developer-facing porting boundary should say.
- `SIM-fovip` is the active simulation/evidence surface. It tests whether
  kernel implementation promises, app-facing storage/compute/network/device
  promises, host assumptions, unsupported behavior, pCID coverage, namespace and
  reference promises, and checkpoint evidence are enough to decide `DR-davod`.

POC12 adds executable pressure to this packet but does not close it. Its kernel
process proves that exact pCID-selected bytes can cross app/kernel and
kernel/kernel boundaries while app processes keep trust and promise judgment.
Its `printer_port` role proves that a "kernel" capability can be a separate
local resource promiser rather than part of one monolithic daemon. POC13 is the
next planned pressure because storage and compute are first-class kernel roles
that need the same promise-first treatment. Source: `DI-gopag`; `DI-punuf`;
`DI-bibom`.

## App/kernel boundary

The app/kernel boundary should be described as a promise boundary.

Exposed app/kernel operations are pCID-selected PromiseGrid messages shaped by
the current formal outer envelope rule:

```text
grid([42(pCID), ...protocol-defined-slots])
```

The recommended example profile remains `grid([42(pCID), payload, ...])` for
ordinary protocols. Source: `DI-rojij`; `DI-punam`.

Local APIs may exist for ergonomics, but they are adapters. An adapter's promise
is to faithfully emit, accept, preserve, or record the corresponding
pCID-selected message and its evidence.

This distinction matters because a local method call does not create authority.
An app does not command a kernel, and a kernel does not command an app. Each
side is an autonomous promiser that may keep, refuse, be unable to perform, or
break a promise. Evidence records make that visible.

## Kernel implementation promise records

A credible PromiseGrid implementation publishes local kernel implementation
promises. The record says what the implementation promises, what it depends on,
what it refuses or cannot do, and what evidence it will keep. Source:
`DI-ripuz`.

The minimum record includes:

- the port identity or local promiser;
- the runtime profile;
- supported pCIDs and unsupported-pCID behavior;
- app-facing promises for storage, compute, network send/receive, key use,
  device access, lifecycle, pCID dispatch, refusal, receipt, evidence,
  namespace, reference, and checkpoint behavior;
- host/runtime assumptions that are not PromiseGrid promises unless the host is
  also an explicit promiser;
- unsupported features and roles;
- exact-byte and local-record evidence policy for kept, refused, unavailable,
  and broken promises.

The record is not a global conformance certificate. Each receiver evaluates the
record and the implementation's make/break history locally.

## Namespaces and references

Alice's local pathname is local. Bob should not treat Alice's string path as
universal truth.

The portable object is a CID-rooted promise-bound reference. A reference should
carry enough context for Bob to decide locally whether to trust and how to map it
into his own view:

```text
root: <CID>
path-or-selector: <pCID-defined path or selector>
target-pcid: <Protocol CID>
frontier: <event/checkpoint CID>
promiser: <agent ref>
promise-body: <what the promiser promises this reference means>
evidence: <promise/event refs>
```

Bob may mount Alice's reference privately, decline it, or map it into a
voluntary group namespace if the group has reciprocal namespace promises.

## Resources as promise-log projections

A file's current content can be understood as a checkpoint over selected promise
history. The same model applies to directories, devices, named pipes, services,
executables, and computation results.

Examples:

- a file resource promises bytes or a checkpoint;
- a directory resource promises name-to-reference entries;
- a device resource promises observations or refusals;
- a stream resource promises ordered messages;
- an executable resource promises invocation behavior;
- a computation result promises how it was produced and what evidence supports
  it.

Branches are different promise-history selections. If Bob breaks a storage
promise, Alice may stop projecting Bob-maintained resources into her convenient
view while preserving evidence of the break.

## Storage and compute pressure

Storage should be described as decentralized sparse CAS promises, not as a
global storage service. A CID can identify bytes or an object, but it does not by
itself promise availability, retention, replication, access, or serving. Alice
may trust Bob to retain and serve an object only after Bob makes that promise and
Alice accumulates local evidence that Bob keeps similar promises.

Compute should be described as pCID-selected function-call promises whose code
identity lives inside the payload. The pCID names the compute protocol, while a
payload-level `function_cid` names the CAS-stored function code. Pure compute
results can be cached by exact function CID, input CIDs or scalar values,
declared context CIDs, ABI/version, and protocol pCID. Impure work should make
timestamp, randomness, sensor reads, network observations, or other ambient
inputs explicit context objects, so the completed run is replayable and
pure-after-the-fact.

This pressure is assigned to `TODO-godad` / planned POC13. It should inform
`SIM-fovip` and `DR-davod`, but it does not freeze final storage, CAS,
function-call, cache, or host API shapes. Source: `DI-bibom`; `DI-gopag`.

## Prior-art pressure

Prior art should influence PromiseGrid only after being reframed as local
promises:

- V is useful pressure for using the same semantic message at local and wire
  boundaries.
- Amoeba is useful pressure for small services and promise-token patterns, not
  global permissions.
- Plan 9 / 9P is useful pressure for simple uniform local views, not a universal
  filesystem authority.
- GNU Hurd is useful pressure for replaceable local services and deferred
  startup, not POSIX namespace authority.

Spring and modern multikernel details remain out of current influence until a
separate patent-clearance task exists.

## What remains open

This design note narrows the kernel model, but it does not freeze:

- the final SDK;
- the final required pCID set for first production ports;
- the exact namespace protocol;
- the exact CID-rooted reference protocol;
- the exact file/resource checkpoint protocol;
- the L6 CAS profile owned by `DR-tumus`;
- guide prose in the external PromiseGrid Development Guide.

`DR-davod` should remain open until focused `SIM-fovip` evidence is reviewed.
If that evidence passes, `DR-davod` can be marked decided with this role/profile
model while guide handoff remains a separate task.

## Source trail

- `DI-fidot` records the design-resolution path.
- `DI-punuf` records the non-monolithic kernel-role synthesis after POC12 made a
  local hardware/resource role concrete.
- `DI-gopag` records the static packet-map and storage/compute pressure pass.
- `DI-bibom` records the planned POC13 CAS storage and CID-named compute shape.
- `DR-davod` owns the kernel-developer porting-boundary decision.
- `TE-jimar` narrows kernel to role/profile rather than process shape.
- `TE-mazop` narrows the minimum credible kernel implementation promises.
- `TE-pudiv` narrows app/kernel operations to pCID-selected messages with local
  APIs as adapters.
- `TE-dunas` narrows safe prior-art pressure.
- `TE-gakoh` narrows decentralized-mainframe, voluntary-namespace,
  CID-rooted-reference, and promise-log resource principles.
