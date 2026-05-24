# TE-jimar: Kernel runtime portability boundary

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-jimar

## Status

needs DF

## Decision under test

What stable PromiseGrid infrastructure boundary should the Development Guide
describe for kernel/runtime porters when the same portable role may appear as a
local daemon, a microkernel dispatcher, a WASM or browser host boundary, a
container host, a header/library-only MCU role, a split mesh of local objects,
or no separately named runtime object at all?

This TE tests whether "kernel" should name one universal process shape or a
portable set of infrastructure roles and evidence promises. It does not settle
the final porting API, close `DR-davod`, or freeze any implementation profile.

## Assumptions and constraints

- PromiseGrid uses Burgess Promise Theory vocabulary: agents are autonomous,
  each agent promises only its own behavior, and no local kernel object may
  promise for another agent or become a global trust authority.
- Trust is local, per-agent, and relationship-specific. Alice may trust Bob to
  store or compute under one promise and not under another; Carol may assess the
  same evidence differently.
- A `pCID` is a Protocol CID, not a payload CID. The pCID names the spec that
  defines payload shape, canonical bytes, signature/proof encoding,
  interpretation, and handler behavior.
- Exact bytes matter. A peer that does not understand a pCID may still preserve
  bytes as evidence under local policy, but carriage is not semantic
  acceptance.
- Current wire-lab evidence points toward small pCID-selected envelopes and
  pCID-owned payload/proof semantics. `DEV-GUIDE-RESOURCES.md` currently treats
  Gordian-oriented structures as active payload/wrapper-family probes rather
  than as a universal dispatch surface.
- `DR-tumus` owns the concrete L6 CAS profile. Its current answerable packet
  recommends a minimal pointer/raw first profile with DAG-CBOR / CIDv1 bridge
  constraints, not an IPLD-first default. `TODO-kituj` must wait for that DR
  before TE-43 freezes broader CAS/profile prose.
- `SIM-lozuk` remains only one L6 starting-profile candidate under
  `SIM-bobud`; `SIM-vagak` remains the current minimal pointer/raw alternative,
  and `SIM-zazit` carries the sparse replication/feed pressure that any CAS
  profile must satisfy.
- Repo evidence supports the substance of the prompt-provided Fermat guidance
  only at profile level: IPLD, CAR, DAG-CBOR, and selectors are useful L6/CAS
  tooling candidates, but they are not a runtime mandate unless a frozen
  pCID-selected protocol requires them.
- TE-sihih's L5/L6/L7 model remains the current layer evidence: feeds move
  chunks over substrates, CAS stores and verifies chunks, and group/app
  protocols define message meaning.
- TE-vipir's protocol-stack tuple remains relevant, but later guidance must
  avoid treating one wire-lab harness layout as the portable PromiseGrid
  runtime.

## Alternatives

### Alt 1 - Classical local daemon/kernel

A PromiseGrid node runs a long-lived daemon or service. Applications call that
daemon through a local API. The daemon owns ingress, pCID dispatch, local
storage, key access, lifecycle, and evidence recording.

**Easier:** Familiar operational model; one supervision point; one place to
apply resource limits; straightforward native host integration.

**Harder:** Does not fit browser tabs, WASM-only hosts, MCU/header-only ports,
library embeddings, or environments where no separate daemon can exist. It also
invites guide prose that accidentally makes the daemon a local authority rather
than one agent's infrastructure.

**New obligations:** Define local API stability, daemon lifecycle,
multi-tenant isolation, upgrade policy, and recovery behavior.

### Alt 2 - Microkernel dispatcher with separate local services

The kernel is a small pCID dispatcher. Storage, key management, feed adapters,
execution sandboxes, and lifecycle services are separate local services behind
explicit interfaces.

**Easier:** Keeps the dispatch core small; supports replacement of CAS, key,
feed, and execution pieces; gives strong design pressure against monoliths.

**Harder:** Still assumes a runtime object and service graph. The word
"microkernel" may over-specify process topology for WASM, browsers, and MCUs.

**New obligations:** Define service boundaries, local authentication between
services, failure propagation, and evidence handoff across local components.

### Alt 3 - Host runtime as kernel boundary

The browser, WASM host, container runtime, OS sandbox manager, or language VM is
the practical kernel boundary. PromiseGrid code supplies handlers and evidence
logic inside that host rather than owning the whole host process.

**Easier:** Fits browser and WASM deployments; reuses existing sandbox,
lifecycle, and security mechanisms; keeps PromiseGrid portable across host
technology.

**Harder:** The host may not expose the exact storage, key, network, or
lifecycle primitives a native daemon expects. A guide that says "implement the
kernel" must not imply control over the browser or WASM engine.

**New obligations:** Specify which promises the PromiseGrid layer can make
inside the host and which promises remain host-provided assumptions assessed
locally by the agent.

### Alt 4 - Header/library-only kernel role

For MCU, Arduino, or embedded ports, there may be no daemon, no allocator-rich
runtime, and no separate process. A small library or header-only layer produces,
parses, verifies, dispatches, or records pCID messages as part of firmware.

**Easier:** Fits constrained devices and hardware-actuator deployments; avoids
forcing OS-like architecture where none exists.

**Harder:** Storage, key custody, crash recovery, and evidence persistence may
be minimal or delegated. The port may implement only one or two pCID-selected
protocols.

**New obligations:** Define profile vocabulary for partial infrastructure
promises, deterministic byte handling, and hardware guard behavior without
requiring a general-purpose runtime.

### Alt 5 - Split local role mesh

The portable "kernel" is a set of local roles: feed adapter, CAS store, key
store, pCID dispatcher, lifecycle supervisor, evidence recorder, hardware guard,
and optional execution host. A deployment may co-locate, split, omit, or
replace roles.

**Easier:** Matches heterogeneous deployment reality; names the roles a porter
must reason about without requiring one binary. It makes local trust and local
evidence boundaries explicit.

**Harder:** More vocabulary for guide writers and implementers. Without a small
profile system, it may become a vague checklist rather than a porting target.

**New obligations:** Define minimal profiles and per-role promises so Alice can
say exactly which roles her port implements and which it delegates.

### Alt 6 - Wire-protocol-only common denominator

PromiseGrid defines only pCID-selected messages, exact bytes, local promises,
and evidence records. Any runtime shape that can send, receive, preserve, and
interpret those messages is acceptable.

**Easier:** Maximally portable; avoids false daemon assumptions; aligns with
the common denominator across daemon, WASM, browser, MCU, split-object, and
future runtimes.

**Harder:** Too thin for developer guidance by itself. Porters still need to
know which local roles they must implement, stub, delegate, or explicitly not
promise.

**New obligations:** Pair the wire common denominator with named portability
profiles so guide prose remains actionable.

## Scenario analysis

### S1 - WASM app

Alice ships a PromiseGrid app as WASM. The WASM host controls filesystem access,
network access, clocks, sandbox lifecycle, and sometimes signing keys.

- **Alt 1** overfits. A daemon may exist outside the WASM module, but the WASM
  app cannot assume it.
- **Alt 2** works if the host exposes service imports, but the microkernel
  shape is host-specific.
- **Alt 3** is natural: the host boundary is the kernel-like surface.
- **Alt 4** is possible for tiny WASM modules that only parse and emit one
  pCID protocol.
- **Alt 5** lets the host provide some roles and the module provide others.
- **Alt 6** identifies the true invariant: pCID messages and exact bytes cross
  the boundary.

The surviving guide shape must say "which promises are made by the WASM module"
and "which promises are host assumptions," not "WASM implements a daemon."

### S2 - Native daemon managing multiple apps

Bob runs a native node with several local apps. The daemon manages ingress,
feeds, a CAS, key access, pCID handler dispatch, and app lifecycle.

- **Alt 1** is strongest here; it gives operators a clear service to supervise.
- **Alt 2** may be better as the daemon grows because keys, CAS, and execution
  can fail or upgrade independently.
- **Alt 5** captures both shapes: the roles may be co-located in one process or
  split later.
- **Alt 6** is necessary but insufficient; Bob still needs local operational
  boundaries.

The daemon is an implementation profile, not the definition of kernel for every
port.

### S3 - MCU/header-only port

Carol ports a single actuator protocol to an MCU. The firmware can parse a
small envelope, verify the specific pCID's proof, record a tiny evidence log,
and drive one pin through a hardware guard.

- **Alt 1** and **Alt 2** are inappropriate if they imply processes, services,
  dynamic storage, or general dispatch.
- **Alt 4** is the natural implementation shape.
- **Alt 5** helps name which roles are present: a tiny dispatcher, key verifier,
  hardware guard, and bounded evidence recorder.
- **Alt 6** provides the common denominator: the MCU must preserve exact pCID
  bytes for the protocols it claims to understand.

A guide that requires daemon features would exclude the small-device class that
the current design explicitly wants to preserve.

### S4 - Split-object runtime

Dave's local environment is a set of interdependent objects: one object watches
feeds, one stores CAS chunks, one holds keys, one dispatches pCIDs, and one
records evidence. There is no single "kernel object."

- **Alt 2** partially matches if the objects act like services.
- **Alt 5** matches better because it names role boundaries without assuming
  service/process mechanics.
- **Alt 6** keeps the inter-object boundary honest: only bytes and local
  promises cross role boundaries.

The guide should let Dave declare local role composition rather than force a
single class or process named `Kernel`.

### S5 - Remote storage or computation

Alice wants Bob to store encrypted chunks or run a computation. She must decide
whether to send Bob data before she can rely on returned storage or compute
evidence. Mallory may advertise a cheap runtime that makes broad claims.

- No alternative creates remote trust by architecture.
- **Alt 1** can hide trust decisions behind "the kernel sent it," which is
  dangerous if guide prose is sloppy.
- **Alt 2** and **Alt 5** can make local trust checks explicit before a feed
  adapter or computation role releases bytes.
- **Alt 6** is the strongest Promise Theory baseline: Bob's messages are only
  Bob's promises and evidence; Alice's trust assessment is local.

The kernel/runtime boundary must never imply that a local kernel globally
authorizes Bob. It can record Alice's local decision, preserve Bob's promises,
and enforce Alice's local policy.

### S6 - Hardware actuator protection

Ellen controls a door lock, valve, or robot arm. The PromiseGrid layer receives
a pCID message asking for physical action.

- **Alt 1** can centralize policy in a daemon on rich hosts.
- **Alt 4** is essential for embedded controllers.
- **Alt 5** names a hardware-guard role that can refuse unsafe action even if a
  dispatcher accepts the pCID message syntactically.
- **Alt 6** ensures the guard acts on exact bytes and pCID-defined semantics,
  not on ambient authority.

The hardware guard promises only its own behavior. It must not be described as
obeying commands from a remote agent or from a global kernel authority.

### S7 - Unknown pCID and mixed version

Frank receives an envelope whose pCID he does not support. Grace supports an
older version of the same family. The bytes may still matter for future audit
or relay.

- **Alt 1** and **Alt 2** need explicit unknown-pCID policy, otherwise the
  runtime may silently drop evidence or accidentally accept unsupported
  semantics.
- **Alt 3** depends on host resource limits and sandbox policy.
- **Alt 4** may hard-reject to protect constrained resources.
- **Alt 5** can make retention a CAS/evidence-recorder promise rather than a
  dispatcher promise.
- **Alt 6** states the invariant cleanly: retain or reject is local policy;
  unknown pCID carriage is never semantic acceptance.

This scenario supports a role/profile model with explicit unknown-pCID promises.

### S8 - Language translation

Alice writes a Rust implementation, Bob writes Go, Carol writes JavaScript, and
Dave writes C for an MCU. They need to interoperate without sharing a runtime
object model.

- **Alt 1** and **Alt 2** risk translating daemon/service architecture instead
  of protocol behavior.
- **Alt 3** is necessary for JavaScript/browser contexts.
- **Alt 4** is necessary for C/MCU contexts.
- **Alt 5** gives each language the same role vocabulary.
- **Alt 6** gives the interoperability anchor: same pCID spec, same canonical
  bytes, same pCID-defined payload and proof interpretation.

The guide should emphasize test vectors, exact bytes, and per-role promises
over class names or process topology.

### S9 - Crash, corruption, and incomplete writes

Mallory cuts power during a write, corrupts a CAS chunk, or crashes one local
role mid-dispatch.

- **Alt 1** has one recovery story but can lose role-specific evidence if all
  logs are hidden inside the daemon.
- **Alt 2** and **Alt 5** force each role to define its own evidence and
  recovery promises.
- **Alt 3** depends on host durability promises.
- **Alt 4** may only preserve a bounded ring log or monotonic counter.
- **Alt 6** says durable recovery must preserve exact bytes and local
  observations where the profile promises persistence.

The boundary should require implementations to state what evidence survives,
not require every runtime to use the same database, CAR file, or event log.

### S10 - Long-horizon future runtime

In 50 years, a future environment has no daemon concept, different hash
algorithms, different proof containers, and storage systems that look more like
content-addressed object fabrics than files.

- **Alt 1** and **Alt 2** may survive as historical profiles but should not be
  the definition.
- **Alt 3** will keep recurring because host runtimes change.
- **Alt 4** will remain necessary for durable small devices.
- **Alt 5** is durable if roles are promise/evidence boundaries rather than API
  objects.
- **Alt 6** is the common denominator: protocol CIDs, exact bytes, local
  promises, local trust, and evidence records.

This scenario rejects one universal process shape most strongly.

## Cross-cutting findings

### Kernel is a role set, not a global authority

The word "kernel" remains useful for guide audiences if it means "the local
infrastructure role set that helps an agent speak pCID-selected PromiseGrid
protocols." It is unsafe if it means a privileged global authority that
commands apps, devices, peers, or remote runtimes.

### Microkernel is useful design pressure

Microkernel thinking is valuable because it keeps dispatch small and separates
feed, CAS, key, lifecycle, hardware, and evidence responsibilities. It should
not become a mandate that every port expose the same daemon/service topology.

### pCID discipline is the stable interoperability surface

The portable surface is not a class name, daemon name, or host API. The stable
surface is: a pCID names a protocol spec; that spec defines payload shape,
canonical bytes, proof/signature encoding, and interpretation; agents locally
assess promises and evidence produced under that spec.

### CAS tooling is below the runtime mandate

`DR-tumus` keeps the first concrete L6 CAS profile open. Its current packet
leans toward minimal pointer/raw first, with DAG-CBOR and CIDv1 carried as
bridge constraints rather than as an IPLD-first runtime mandate. `SIM-lozuk`
keeps the IPFS/IPLD-aligned option alive, but `SIM-vagak` remains the current
minimal-profile recommendation and `SIM-zazit` keeps sparse replication
pressure active.

IPLD, CAR, selectors, DAG-CBOR, Rabin chunking, and Merkle assembly can be
excellent L6/CAS tools. A runtime profile may use them, and a frozen pCID may
require one of them. They should not be taught as mandatory properties of every
kernel/runtime port, they should not replace L5 feed promises, and IPLD schemas
should not be treated as validators of full PromiseGrid meaning because pCID
still owns protocol semantics.

Current scored `SIM-lozuk` evidence also records the long-horizon risk: IPLD /
CID / multicodec tooling may pull in registry or stale-tooling dependencies
that age poorly, and those dependencies can be especially expensive for
small-device or 100-year reimplementation targets. That risk belongs in the
portable-boundary discussion whenever IPLD-aligned profiles are mentioned.

### Gordian belongs in pCID-owned families unless evidence broadens

Gordian-style structures may be valuable as payload, wrapper, selective
disclosure, or proof families under specific pCIDs. Current guide-resource
evidence does not support making Gordian the universal dispatch surface for all
PromiseGrid runtimes.

## Conclusions

- Reject a single universal process shape as the PromiseGrid kernel definition.
- Reject daemon-only and microkernel-only guide prose as too narrow for WASM,
  browser, MCU, split-object, and future runtimes.
- Reject wire-protocol-only guidance as too thin for porters unless it is paired
  with named infrastructure roles and profiles.
- Keep the common denominator as pCID-selected messages, exact canonical bytes,
  local promises, local trust assessment, and evidence records.
- Treat "kernel" as a portable infrastructure role set/profile. A deployment may
  realize that profile as a daemon, microkernel service graph, host boundary,
  library/header-only role, split local mesh, or future runtime shape.
- Treat microkernel architecture as useful design pressure and one likely rich
  host profile, not as the universal binary architecture.

## Surviving alternatives

- **Alt 5 survives as the guide-facing structure:** a split local role mesh or
  co-located equivalent with explicit promises for feed, CAS, key, dispatch,
  lifecycle, evidence, execution, and hardware guard roles.
- **Alt 6 survives as the interoperability floor:** pCID messages, exact bytes,
  local promises, local trust, and evidence.
- **Alt 1, Alt 2, Alt 3, and Alt 4 survive as profiles:** daemon, microkernel,
  host-runtime, and header/library-only are implementation profiles of the role
  set, not definitions of kernel itself.

## Recommended DF question

Should `DR-davod` move toward guide prose that defines a PromiseGrid kernel as:

1. **Alt A - Role/profile definition:** "kernel" means the portable local
   infrastructure role set that lets an agent speak pCID-selected protocols,
   with daemon, microkernel, host-runtime, library-only, and split-object
   profiles as implementation shapes.
2. **Alt B - Runtime/process definition:** "kernel" means a concrete local
   runtime object or daemon that all ports must expose, with constrained hosts
   treated as exceptions.
3. **Alt C - Wire-only definition:** avoid "kernel" as a porting target and
   teach only pCID message interoperability plus implementation-specific
   evidence promises.

This TE recommends Alt A because it preserves the common denominator without
making guide prose too thin for porters.

## Implications for open work

- `DR-davod` remains open. TE-jimar supplies a narrowed DF question but does not
  close the stable kernel-developer porting boundary.
- `DEV-GUIDE-RESOURCES.md` should cite TE-jimar as current evidence for
  provisional Kernel Dev guidance and should continue to say settled porting
  instructions are blocked by `DR-davod`.
- `DR-tumus` remains the concrete owner for any first L6 CAS profile. TE-jimar
  may mention optional CAS tooling pressure, but it must not lock IPLD, CAR,
  selectors, multicodec rules, or promisebase-derived profile choices.
- Future guide prose should avoid saying "implement the daemon" unless it is
  explicitly describing a daemon profile.
- Future porting specs should define role/profile promises and exact byte test
  vectors before defining language-specific APIs.
- Any future L6/CAS profile may require IPLD, CAR, selectors, or Gordian-style
  payloads under a specific pCID, but the kernel/runtime boundary should not
  mandate them globally.

## Decision status

`needs DF` - this TE rejects one universal process shape and recommends
defining kernel as a portable infrastructure role/profile, but the final
`DR-davod` decision remains unmade.
