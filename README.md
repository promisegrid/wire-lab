# PromiseGrid Wire Lab

Wire Lab is the experimental workspace where PromiseGrid protocol, kernel,
runtime, storage, and application ideas are explored before they are written as
guide or production material. If you are returning after several months away,
the short version is: the project is still promise-first, but the recent work has
made the message format, sparse CAS substrate, executable POCs, and production
trajectory much more concrete.

PromiseGrid is a decentralized computing model based on promises between
agents. Agents do not command each other. They make local promises, exchange
messages, keep or break those promises, and update local trust from their own
observations. Wire Lab exists to find the message formats and runtime patterns
that make that practical across machines, organizations, legal entities,
constrained devices, browsers, containers, and long-lived shared systems.

## Current Direction In One Page

The strongest current direction is:

- messages use compact CBOR `grid(...)` envelopes;
- slot 0 carries `42(pCID)`, where the protocol CID identifies the protocol spec;
- the pCID-selected spec owns all following slots, payload shape, signable view,
  and proof semantics;
- pCID is a protocol selector, not a peer address, app address, repository name,
  message type, route, or operation code;
- protocol objects help agents make, recognize, remember, and evaluate promises;
- trust is local and relationship-specific, not global;
- CIDs, sparse CAS, parent-linked message/object DAGs, and signed capability
  tokens are becoming the durable substrate for storage, synchronization,
  collaboration, app installation, and execution.

The recent POC line moved from simple two-container hello messages to
containerized agents exchanging exact `grid()` CBOR messages over TCP, signed
CWT/COSE capability tokens, sparse per-agent CAS stores, CAR payloads for object
transfer, and diagnostic raw artifacts for later review. POC18 applies those
ideas to a CAS-backed version-control system. POC19 is now planned as the first
production-shaped pass: one `grid` binary as the minimum microkernel, with
daemon/client modes, VCS/CAS apps, WASI-first execution, CID-rooted app/runtime
updates, and equal TCP/WebSocket message transports. Source: `DI-kodob`.

The design is still provisional. The POCs are executable evidence and current
direction, not frozen public APIs.

## What This Repo Contains

- `simulations/` contains generated and hand-curated protocol experiments.
- `implementations/` contains executable POCs that test design candidates.
- `docs/` contains thought experiments, research notes, and design notes.
- `protocols/` contains TODOs, decision intent logs, decision requests, and
  protocol-level planning.
- `DEV-GUIDE-RESOURCES.md` is the detailed source map for people writing the
  PromiseGrid Development Guide.

The repo is intentionally not a polished product tree. It is a lab notebook plus
executable prototypes. The key is that important design changes are supposed to
leave a trail: TODOs, DIs, TEs, protocol specs, POC artifacts, and source-map
updates.

## Before The POCs: Simulations, GA, And Design Search

Before the recent executable POC sequence, much of the work happened in
simulations and genetic-algorithm-assisted design search. The goal was not to
let generated artifacts become the design. The goal was to compare candidate
protocol shapes, discover weaknesses, calibrate scoring rubrics, and surface
questions that needed human decisions.

That design-search phase clarified several durable rules:

- PromiseGrid must stay Promise Theory-first. Protocols help agents make,
  recognize, remember, and evaluate promises; they do not pretend to command
  independent agents.
- Trust is local. There is no global trust oracle, global branch owner, global
  monitor, global access gatekeeper, or universal namespace authority.
- A pCID is the CID of a protocol spec document. It is not a payload CID and not
  a message kind.
- The outer envelope should stay small and durable: a CBOR `grid(...)` wrapper
  whose slot 0 is the protocol selector.
- Any higher-layer routing, app identity, operation, or resource meaning belongs
  inside the pCID-defined payload or nested objects, not in the universal
  envelope.
- Observer and analyzer outputs are lab machinery. They can help us evaluate a
  POC run, but they are not production monitors or trust authorities.

The GA and simulation work also corrected vocabulary. We now avoid language
that makes PromiseGrid sound like a conventional method-call service,
commanding supervisor, or centralized policy system. The preferred language is
local promises, local events, voluntary cooperation, capability promises, sparse
CAS, exact bytes, parent links, and local trust updates.

## The POC Journey

### POC2 Through POC6: Message Shape, Local Kernel, Trust, And DAG-CBOR

The early POCs established that the same PromiseGrid message shape could be used
between apps and a local kernel role and between peer kernels. `poc2` proved a
minimal two-container app/kernel and kernel/kernel hello flow. `poc3` added
multiple apps per container. `poc4` moved multi-hop behavior into relay apps
rather than making the kernel a universal router. `poc5` showed local trust and
selective sending after a broken storage promise. `poc6` tested DAG-CBOR/IPLD
interop, CID links, and tag-42 link encoding.

The important lesson was that the kernel should not become a central brain. It
is better understood as a local role/profile set: transport, app interface,
resource access, lifecycle, storage, compute, key use, event recording, and
pCID-selected parsing/building. Different runtimes can collapse or split those
roles differently.

### POC7 Through POC11: Capability Tokens, Economy, Discovery, And Autonomy

`poc7` introduced signed bearer and non-transferable capability tokens as
promises by issuers. Agents could redeem tokens for storage or compute access,
track kept or broken promises, and make local economic decisions. `poc8` made
the promise economy richer with offers, counteroffers, local exchange-rate
signals, and reciprocal promises. `poc9` added sparse-mesh discovery, referrals,
route promises, and expired-token semantics. `poc10` brought in LLM-backed
autonomous agents while keeping protocol bytes owned by deterministic Go code.
`poc11` added multi-round relationship decay/repair, local economics, and
trust-correlated TCP links.

This era made two ideas concrete. First, a capability token is not a grant from
outside the promise relationship; it is an issuer's signed promise that a holder
may later redeem under stated terms. Second, agent autonomy is useful only when
bounded by protocol specs, exact bytes, local decisions, and clear promise
records.

### POC12 Through POC16: Production Workflows, CAS/Compute, WASM, Multihop, And Parser Roles

`poc12` introduced a more production-like shipping workflow: scale, UPS label
printer, accounting system, fulfillment app, local kernel, and printer-port
resource role. The printer-port work was important because it showed a
non-monolithic kernel role that owns a local hardware resource and issues future
printing capability promises.

`poc13` added decentralized CAS storage and CID-named compute. It tested real
store/serve/retrieve/replicate flows, named compute functions, verification,
caching, key rotation, bad-proof handling, token lifecycle, and local recovery.
It also reinforced that local durability inside a run matters for a POC, while
cross-run state should usually reset unless a scenario explicitly needs it.

`poc14` added heterogeneous runtime boundaries. Peggy became a real WASM role
and Victor used a stdio-only worker path. This exposed a major discipline point:
we should not claim fake or simulated boundaries are production-like. If a POC
says WASM or stdio works, it must actually cross that boundary.

`poc15` explored multihop, multiple pCID-owned envelope arities, parent links in
envelopes or payloads, raw CBOR diagnostics, and persistent TCP sessions. It
kept exact message bytes and parent-linked message DAGs visible for review.

`poc16` consolidated several wire/runtime lessons: pCID-selected parser/builder
roles, implementation-local RFC-like protocol specs, secure CWT/COSE-shaped
tokens, encrypted payload coverage, map-permitted payload profiles when a pCID
allows maps, binary CIDs on wire, CIDv1 base32 when printable, and clearer
lifecycle/resource promises.

### POC17: Constrained M4/LoRa-Shaped Runtime

`poc17-m4-lora-runtime` is a constrained-runtime behavior simulator. It does not
claim production Feather M4 or RFM95 radio fidelity. Its value is showing that
PromiseGrid messages can stay compact and pCID-selected under a small radio-like
MTU. It uses actual pCID bytes in slot 0, compact payloads, base32 CIDs in
printable logs, radio-visible peer storage promises, and restart recovery for
missing parent objects.

The POC17 lesson is that small devices do not need heavyweight universal
payloads. A constrained device can match the `grid` tag, parse slot 0, recognize
the pCID it supports, and apply the compact payload rules in that spec.

### POC18: CAS-Backed Version Control And Collaboration

`poc18-cas-git-replacement` applies the recent protocol work to version control
and collaboration. The core idea is that reference sets are the shared
abstraction for directories, filenames, branches, tags, releases, logical
changes, review threads, workspace roots, and app installations.

POC18 uses sparse filesystem CAS, Rabin content-defined chunking, POSIX node
promises, directory/reference-set promises, snapshots, multi-parent merges,
review statements, local adoption decisions, and a Git bridge using go-git.
Conventional Git import/export/push/pull are compatibility adapters, not the
native PromiseGrid sync model.

Native PromiseGrid collaboration is continuous peer DAG synchronization. Agents
advertise selected reference-set heads, request missing objects, promise object
availability, redeem capability tokens, transfer CAR/object bytes, verify exact
CIDs, and decide locally what to retain. Recent POC18 work moved those
inter-agent flows out of simulated in-process transfer and into Docker
containers over persistent length-framed TCP. Diagnostic CBOR rendering now
keeps representative exact messages reviewable.

POC18 also introduced a usable first `grid` CLI surface: `init`, `snapshot`,
`status`, `log`, `track`, `untrack`, `checkout`, `refs`, `diag`, `git`, and
`sync`. Repo-local `.grid/config.json` and `.grid/state.json` are user-facing
convenience state; the deeper storage model is sparse CAS plus signed promise
objects.

### POC19: Production Shape

`poc19-production-shape` is currently a design artifact, not implementation
code. Its goal is to turn POC18's CAS/VCS/TCP work plus POC16/POC17 runtime
lessons into a deployable-node shape.

The planned target is one binary named `grid`. That installed binary is the
minimum microkernel and local loader, not the app distribution package. The same
binary should be able to run as `grid daemon`, expose CLI commands, act as a VCS,
fetch non-kernel code and data from peers over TCP or WebSocket, and execute
fetched apps from VCS/CAS. First run or local config may name a bootstrap Merkle
root CID; later app/runtime root changes are fetched as exact CAS graphs and
adopted only after local operator approval. The first execution profile is WASI.
OCI container images and native binaries are later, higher-risk profiles that
require stronger local lifecycle and resource promises. Source: `DI-kodob`.

A grid app is installed by checking a signed app reference set into VCS/CAS.
`grid run` asks the local daemon whether it currently promises to execute that
app under local resource constraints. The daemon may fetch missing code/data
objects from peers, verify exact CIDs, issue local runtime capability promises,
run the app, and store result objects back in CAS.

## The Current Wire And Storage Model

The current outer message direction is:

```text
grid([42(pCID), ...protocol-defined-slots])
```

The `grid(...)` wrapper is a CBOR tag for PromiseGrid messages. Slot 0 is a
tag-42 CID link to the protocol spec document. The pCID-named spec defines every
following slot: whether there are parents, a payload, a proof, a COSE object,
CAR bytes, or another pCID-owned shape.

This is the key boundary rule:

```text
pCID selects the protocol parser.
The payload carries protocol meaning.
The pCID is not an address, method, route, or operation.
```

CIDs are binary on the wire and CIDv1 base32 text when printed. Exact message
CIDs identify exact CBOR envelope bytes. Parent links point at exact message or
object CIDs. CAS stores are sparse: each agent keeps only what it chooses or has
promised to retain. A CID names bytes; it does not by itself promise
availability, retention, access by a holder, or trustworthiness.

Capability tokens are signed promises. In recent POCs they are CWT/COSE-shaped
bytes stored and transferred as exact objects. A storage token, retrieval token,
runtime token, or lifecycle token is useful because the issuer made a promise
that another agent may locally decide to rely on. Redemption, replay rejection,
expiry, revocation, and trust updates remain local to the agents involved.

## How To Re-Enter The Work

For the returning team member, the practical concept map is:

1. Read this README for the current story.
2. Read `DEV-GUIDE-RESOURCES.md` for the detailed source map and current design
   state.
3. Read `implementations/README.md` for the POC catalog.
4. Read `implementations/poc19-production-shape/docs/DESIGN.md` for the current
   production-shape target.
5. Drill into POC-specific TODOs and protocol specs only when you need exact
   implementation detail.

The most useful mental model is: wire-lab is not trying to perfect each POC in
isolation. Each POC moves one or more ideas from vague to executable. The durable
results are the design constraints, exact artifacts, and decision records that
survive into later POCs.

## What Is Not Final

The following are still active design areas:

- the final SDK and stable app/kernel boundary;
- the final production pCID set;
- the exact POC19 daemon/client implementation shape;
- the durable object-store profile for raw chunks, DAG-CBOR, GRID-CBOR, and CAR
  storage;
- the final WASI host ABI and runtime token shape;
- production identity, key management, and recovery;
- production browser/mobile/MCU profiles;
- public PromiseGrid Development Guide prose.

When in doubt, treat executable POCs as evidence and DIs/TEs/TODOs as the design
trail. Do not treat a POC-local analyzer, fixture, monitor, collector, or
simulator as production architecture unless a later DI explicitly says so.

## License

GPL-3.0, matching the rest of PromiseGrid.
