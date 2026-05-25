# TE-dunas: Prior-art influence on the PromiseGrid kernel

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-dunas

## Status

needs DF

## Decision under test

How much should four older, lower-patent-risk distributed operating-system
ideas influence the PromiseGrid kernel and app/kernel boundary?

The four systems considered here are:

1. V distributed kernel.
2. Amoeba distributed operating system.
3. Plan 9 / 9P.
4. GNU Hurd translators.

Spring is intentionally excluded from this TE. Modern multikernel work such as
Barrelfish is also excluded until patent risk is separately cleared.

## Assumptions and constraints

- This is not legal advice and does not perform patent clearance. It is an
  engineering filter: use older published systems as prior-art context, avoid
  copying modern patent-active designs, and do not copy implementation details.
- PromiseGrid is not trying to provide a transparent single-system image. It is
  a decentralized virtual machine spanning multiple agents and legal entities.
- Promise Theory remains the primary design constraint. No agent promises on
  behalf of another agent. Cooperation is voluntary. Trust is local, relative,
  relationship-specific, and updated from promise make/break evidence.
- The current outer message shape remains `grid([42(pCID), payload, ...])`.
- `TE-pudiv` keeps open whether local app/kernel operations are themselves
  pCID-selected grid messages or whether local APIs are adapters over those
  messages.

## Candidate influence levels

Each prior-art system can influence PromiseGrid at one of four levels:

1. **Ignore:** useful history, no direct design pressure.
2. **Vocabulary caution:** use it only to avoid known traps.
3. **Pattern influence:** borrow a pattern after reframing it as local promises.
4. **Simulation pressure:** create or adjust sims so candidate PromiseGrid
   designs must survive that system's strongest scenario.

No prior-art system in this TE is allowed to become authority for PromiseGrid.
Each influence must be translated into PromiseGrid's own terms: pCID-selected
messages, autonomous promisers, evidence records, local trust judgments, and
explicit refusals when a promise is not made.

## Prior art 1: V distributed kernel

### Relevant idea

V is the most directly relevant prior art for `TE-pudiv`: one message-oriented
kernel mechanism handled local and network interprocess communication. That
maps closely to the hypothesis that app/kernel operations and wire operations
may share the same PromiseGrid message shape.

### What PromiseGrid should borrow

- Strongly consider one message model for local and remote promise boundaries.
- Treat local service boundaries as real message/evidence boundaries, not as
  invisible implementation details.
- Keep the local fast path an optimization, not a separate semantic surface.
- Make every visible operation pCID-selected so later agents can replay why an
  app, kernel role, or peer believed a promise was made, kept, refused, or
  broken.

### What PromiseGrid should reject

- Do not inherit V's single-system transparency goal. PromiseGrid peers are not
  one administrative machine.
- Do not hide network boundaries when they cross trust or legal-entity
  boundaries. Alice must know when she is asking Bob, Carol, or a local adapter
  to promise something.
- Do not define one kernel authority that commands all services. Services are
  autonomous promisers.

### Influence level

**Simulation pressure.** V should strongly influence the next kernel-boundary
sim. The sim should test whether storage, compute, send, key-use, receipt,
refusal, and evidence operations work better as pCID-selected grid messages
across both local and remote boundaries.

## Prior art 2: Amoeba

### Relevant idea

Amoeba combined distributed microkernels, user-level servers, RPC, object-like
resources, and capabilities. It is useful for thinking about a small local
kernel plus services, and for comparing capability-style resource references
with PromiseGrid promise-as-capability-token ideas.

### What PromiseGrid should borrow

- Keep the kernel small and push policy into separately accountable services.
- Let services make explicit promises about resources they manage.
- Study capability references as one possible encoding of "Bob promises Alice
  that this token will be recognized for this resource under these conditions."
- Use server-authenticity and resource-reference lessons as scenarios for
  promise evidence: who promised the resource exists, who observed it, and what
  happens when the promise breaks.

### What PromiseGrid should reject

- Do not treat capabilities as global permission objects. In PromiseGrid, a
  capability-like token is useful only if it is tied to a promiser, a promise
  body, evidence, and local trust assessment.
- Do not adopt transparent processor pools or location-hiding as a goal.
  Selective sending is core: Alice may refuse to send data to Bob until she
  trusts Bob's prior promise history.
- Do not centralize naming, scheduling, or resource interpretation into a system
  authority.

### Influence level

**Pattern influence plus simulation pressure.** Amoeba should influence
small-kernel/service decomposition and capability-token scenarios, but every
borrowed pattern must be reframed as promises and peer-local trust evidence.

## Prior art 3: Plan 9 / 9P

### Relevant idea

Plan 9 is not a microkernel, but it is important distributed-OS prior art
because it presents local and remote resources through a common protocol and
per-process namespaces. It shows how a simple resource protocol can unify many
interfaces without forcing all resources into the kernel.

### What PromiseGrid should borrow

- Prefer one simple, durable message/resource protocol over many special-case
  APIs.
- Let each app or agent see a local view assembled from promises it accepts,
  instead of assuming a universal namespace.
- Make remote resources feel ordinary only after the local agent chooses to
  rely on the promiser and records evidence.
- Keep resource access inspectable and replayable.

### What PromiseGrid should reject

- Do not make "everything is a file" a PromiseGrid rule. PromiseGrid's common
  object is a promise-shaped message, not a filesystem node.
- Do not imply central administration or a universal namespace. Namespaces are
  local views produced by local trust and local promise acceptance.
- Do not hide who is promising a resource. The promiser must remain visible.

### Influence level

**Pattern influence.** Plan 9 should influence the simplicity and uniformity of
PromiseGrid app/kernel interfaces, but not the kernel shape itself. The useful
lesson is not "use files"; it is "use one small protocol and local views."

## Prior art 4: GNU Hurd translators

### Relevant idea

GNU Hurd translators are user-space servers attached to filesystem nodes. They
demonstrate that behavior traditionally buried in a kernel can be supplied by
replaceable user-space agents. Passive translators also show deferred startup
of services when a resource is first accessed.

### What PromiseGrid should borrow

- Treat many kernel-looking features as service promises made by replaceable
  local agents.
- Support deferred service startup as a promise: "I will start this service when
  this pCID-selected operation is requested, or I will record a refusal/failure."
- Let different local agents provide different resource views without changing
  the minimum kernel role.
- Use service restart and failure as first-class evidence events.

### What PromiseGrid should reject

- Do not tie PromiseGrid semantics to filesystem inodes, POSIX permissions, or
  root/non-root authority.
- Do not assume a single machine's access-control model maps to cross-entity
  trust.
- Do not let an adapter promise on behalf of the service it starts. The adapter
  may promise only its own startup, routing, and evidence-recording behavior.

### Influence level

**Pattern influence.** Hurd translators are strong prior art for modular
user-space services and deferred service activation. They should influence
runtime profiles and adapter sims, but not become a filesystem-first
architecture.

## Scenario analysis

### S1 - Same local/wire message boundary

Alice's app asks Bob's local PromiseGrid kernel role to store data, then later
asks Carol's remote peer to compute over a derivative. V pushes toward one
message/evidence model for both requests. Amoeba pushes toward resource tokens
and accountable services. Plan 9 pushes toward one simple resource interface.
Hurd pushes toward replaceable local service adapters.

The PromiseGrid result should be: each request is a pCID-selected promise
interaction with a visible promiser, explicit refusal path, and evidence record.
The local path may be optimized, but it should not erase who promised what.

### S2 - Cross-legal-entity trust boundary

Alice refuses to send private data to Bob until Bob has enough kept-promise
history. Classic distributed OS transparency is actively harmful here: hiding
the remote boundary would hide the trust decision. V and Plan 9 are useful only
after PromiseGrid keeps the promiser visible. Amoeba capabilities are useful
only if interpreted as Bob's promises, not as global permissions.

The PromiseGrid result should be: no transparent remote execution, no hidden
processor pool, and no remote resource access without local trust judgment.

### S3 - Browser/WASM host

Carol's browser app has no traditional kernel. A host adapter can promise to
store bytes, send messages, use WebCrypto, or report failure. Hurd-style
replaceable services fit this case well. Plan 9-style uniformity helps keep the
surface small. V-style same-message IPC helps evidence. Amoeba-style
capabilities are risky unless they remain promises made by named agents.

The PromiseGrid result should be: host APIs are adapters over pCID-selected
messages and evidence records, not normative authority.

### S4 - MCU or header-only port

Dave's small device supports one actuator pCID and one evidence pCID. V's
uniform message idea is helpful because the device can parse the same shape at
local and wire boundaries. Plan 9 and Hurd are less directly applicable because
there may be no namespace or process model. Amoeba's small-kernel/resource
server split is useful only if the device has enough runtime to split roles.

The PromiseGrid result should be: the minimum kernel role can shrink to exact
message parsing, pCID dispatch, one or two local promises, and bounded evidence.

### S5 - Long-horizon port decades later

Frank reimplements a PromiseGrid port after current OSes, browsers, and
language runtimes have changed. V's uniform local/network message lesson and
Plan 9's small-protocol lesson age well. Hurd's "services can be ordinary
programs" lesson also ages well. Amoeba's transparent distributed-system
ambition ages poorly for PromiseGrid because it hides agency and trust.

The PromiseGrid result should be: durable pCID specs and evidence records
matter more than any inherited OS mechanism.

## Cross-cutting findings

- V is the strongest direct influence on `TE-pudiv`: same semantic message at
  local and wire boundaries.
- Amoeba is useful for small-kernel/service decomposition and capability-token
  pressure, but its capabilities must be PromiseGrid promises, not permissions.
- Plan 9 is useful for simple uniform resource interfaces and local views, but
  PromiseGrid's common surface is a message/promise protocol, not files.
- Hurd is useful for replaceable user-space services and deferred activation,
  but PromiseGrid must avoid filesystem/POSIX authority assumptions.
- All four prior-art systems become unsafe if they encourage transparency that
  hides who promised what across a trust boundary.

## Conclusions

- Use **V** as a strong simulation pressure for same-grid app/kernel messages.
- Use **Amoeba** as capability/service decomposition pressure, but heavily
  reframe it through promises and local trust.
- Use **Plan 9 / 9P** as a simplicity and uniform-interface pressure, not as a
  filesystem mandate.
- Use **GNU Hurd translators** as modular-service and deferred-start pressure,
  not as a POSIX namespace mandate.
- Skip Spring and modern multikernel systems for current design influence until
  patent risk is separately cleared.

## Recommended next DF packet for DR-davod

Before deciding the stable kernel-developer porting boundary, answer:

1. Should the next kernel-boundary sim explicitly include V-style same-message
   local/wire operations as the strongest candidate?
2. Should capability-like tokens be tested only as promises made by named
   promisers, never as global permission objects?
3. Should app/kernel guide prose describe local APIs as adapters over a small
   pCID-selected message/resource protocol?
4. Should runtime profiles allow Hurd-like replaceable local services and
   deferred service startup, with startup/refusal/failure as evidence events?
5. Should the prior-art filter exclude Spring and modern multikernel details
   until a later patent-clearance task exists?

## Implications for open work

- `DR-davod` remains open.
- `TE-pudiv` should use this TE as prior-art support for the same-grid
  app/kernel hypothesis.
- `SIM-fovip` or a successor sim should test V-style same-message boundaries,
  Amoeba-style promise tokens, Plan 9-style uniform resource views, and
  Hurd-style replaceable/deferred local services.
- `DEV-GUIDE-RESOURCES.md` should describe these as influence guardrails, not
  as imported designs.

## Decision status

`needs DF` - this TE narrows safe prior-art influence for kernel work but does
not decide the final PromiseGrid kernel boundary.

## Refinements

### 2026-05-25 - Decentralized-mainframe local-view synthesis

`TE-gakoh` follows this TE with the broader synthesis that arose after the
prior-art analysis: PromiseGrid can preserve mainframe-like shared affordances
without a single-system-image authority by giving each agent a local,
trust-filtered, file-like view over a shared promise/event hypergraph. This is a
forward pointer only; the analysis above remains unchanged historical evidence.
Source: `DI-kuvum`.
