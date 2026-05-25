# PromiseGrid kernel role/profile design

This design note explains the current PromiseGrid kernel direction in plain
English. It is not a TE, not a frozen protocol spec, and not final SDK prose. It
is the current design synthesis for resolving `DR-davod`: what guide writers can
say about kernel developers after the kernel TEs converged, and what still waits
for focused evidence and future frozen pCID specs. Source: `DI-fidot`.

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

The kernel role set may include dispatch, storage, compute, networking, key use,
device access, lifecycle, pCID handling, evidence recording, namespace
projection, and reference resolution. In a split deployment, different local
agents may make these promises separately.

## App/kernel boundary

The app/kernel boundary should be described as a promise boundary.

Exposed app/kernel operations are pCID-selected PromiseGrid messages shaped by
the current outer envelope direction:

```text
grid([42(pCID), payload, ...])
```

Local APIs may exist for ergonomics, but they are adapters. An adapter's promise
is to faithfully emit, accept, preserve, or record the corresponding
pCID-selected message and its evidence.

This distinction matters because a local method call does not create authority.
An app does not command a kernel, and a kernel does not command an app. Each
side is an autonomous promiser that may keep, refuse, be unable to perform, or
break a promise. Evidence records make that visible.

## Port promise records

A credible PromiseGrid port publishes a local promise record. The record says
what the port promises, what it depends on, what it refuses or cannot do, and
what evidence it will keep.

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
record and the port's make/break history locally.

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
- `DR-davod` owns the kernel-developer porting-boundary decision.
- `TE-jimar` narrows kernel to role/profile rather than process shape.
- `TE-mazop` narrows the port promise record.
- `TE-pudiv` narrows app/kernel operations to pCID-selected messages with local
  APIs as adapters.
- `TE-dunas` narrows safe prior-art pressure.
- `TE-gakoh` narrows decentralized-mainframe, voluntary-namespace,
  CID-rooted-reference, and promise-log resource principles.
