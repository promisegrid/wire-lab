# TE-gakoh: Local views over a promise/event hypergraph

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-gakoh

## Status

needs DF

## Decision under test

Can PromiseGrid be explained as a decentralized mainframe where each agent sees a
local, trust-filtered, file-like view over a shared promise/event hypergraph?

This TE captures the post-`TE-dunas` synthesis:

- PromiseGrid should recover the coherent shared affordances people lost when
  mainframe or organization-wide minicomputer environments gave way to
  URL-and-install workflows.
- PromiseGrid must not recover centralized authority, a global namespace, a
  single-system-image illusion, or a hidden trust oracle.
- Each agent's local view should be assembled from promises the agent can
  observe, recognize, and locally choose to trust.
- IPLD / DAG-CBOR can carry the durable graph substrate without owning
  PromiseGrid semantics.
- The long-horizon storage and computation model may be a Git-like
  content-addressed event graph or hypergraph.

This TE does not decide the final kernel boundary, pathname-reference format,
L6 CAS profile, or app API. It narrows the model that later sims and DFs should
test.

## Assumptions and constraints

- `TE-dunas` remains the prior-art influence TE. This TE does not rewrite it.
- `TE-pudiv` remains the same-grid app/kernel message-boundary TE. This TE
  builds on its open question rather than answering it.
- `DR-davod` remains open.
- The current outer envelope direction remains `grid([42(pCID), payload, ...])`.
- PromiseGrid follows Promise Theory: no agent promises on behalf of another
  agent; trust is local; promises are voluntary; evidence records support local
  assessment.
- In the broader Burgess / Semantic Spacetime reading, agents are not limited to
  humans or software. A byte chunk, atom, rock, service, or executable may be
  modeled as an agent for promises it can plausibly keep or embody.
- IPLD is treated as a content-addressed data substrate and interoperability
  ecosystem, not as the source of PromiseGrid semantics.

## Terms for this TE

- **Local view:** An agent's own file-like arrangement of resources, services,
  executables, promises, and evidence. A local view is not a global namespace.
- **Promise-bound reference:** A shareable reference that carries enough context
  for another agent to resolve an object or service in its own view: target CID
  or root, pCID, selector/path, version or frontier, promiser identity, promise
  body, and evidence/capability references as needed.
- **Promise/event hypergraph:** A durable content-addressed structure whose
  nodes and hyperedges describe agents, byte chunks, payloads, promises,
  requests, refusals, receipts, computations, observations, and trust updates.
- **File-like UX:** A developer or user experience in which shared data,
  executables, services, and evidence can be opened, listed, linked, and run as
  nearby named resources without a download/install ritual.

## Alternatives

### Alt 1 - Global single-system image

PromiseGrid presents one universal filesystem-like namespace. Alice's path and
Bob's path mean the same thing because "the system" defines a single global view.

**Easier:** This is the most mainframe-like experience. It is simple to explain
and convenient when everyone belongs to one administration domain.

**Harder:** It contradicts PromiseGrid's autonomy model. It hides who promised
what, erases local trust judgment, and implies a namespace authority.

**New obligations:** This would require a global naming, access-control, and
trust model. Those obligations are rejected for PromiseGrid.

### Alt 2 - Local file-like views over promise-bound references

Each agent has a local file-like view. Alice can call something
`/alice/work/foo`; Bob may map the same shared object as `/from/alice/foo` or
not map it at all. Alice shares a promise-bound reference, not a bare pathname.

**Easier:** Preserves mainframe-like convenience while keeping agency visible.
Alice and Bob can each arrange resources in locally meaningful ways.

**Harder:** Requires a reference object or message shape that carries enough
context for Bob to resolve Alice's intended resource without treating Alice's
path as global truth.

**New obligations:** Define how a shareable reference names root content, pCID,
selector/path, version/frontier, promise body, evidence, and optional
promise-as-capability-token behavior.

### Alt 3 - Message-only model

Everything meaningful is a pCID-selected message. Files, services, paths,
executables, and graph events are all views over messages.

**Easier:** Aligns strongly with `TE-pudiv` and the current envelope direction.
Every boundary crossing is explicit and evidence-friendly.

**Harder:** Stable data, executables, and long-lived resources become awkward to
explain if all of them are first described as messages. It overfits transport
and interaction to storage and UX.

**New obligations:** Define enough higher-level resource vocabulary that users
do not have to think in raw transport messages.

### Alt 4 - IPLD graph-first model

Everything meaningful is an IPLD object or graph node. Hyperedges are reified as
IPLD objects with links to participants and payloads.

**Easier:** Strong content-addressing, deterministic byte identity, link
traversal, bridgeability to IPLD / IPFS / Bluesky-adjacent tooling, and a
natural substrate for Git-like event history.

**Harder:** IPLD does not supply PromiseGrid agency, trust, promiser identity,
local view semantics, or promise make/break interpretation. A pure graph model
can become storage-first and forget Promise Theory.

**New obligations:** Keep pCID-defined PromiseGrid semantics above IPLD. Treat
IPLD links as substrate links; treat promises, refusals, observations, and trust
updates as PromiseGrid protocol objects.

### Alt 5 - Promise-first object model

Everything useful is a promise. Files, messages, graph nodes, services,
executables, path entries, and byte chunks are promise-bearing objects or
evidence of promises.

**Easier:** Best matches Promise Theory and the Semantic Spacetime correction:
agents can be human, machine, mineral, byte chunks, atoms, or any modeled
autonomous element, so long as the promise body is something that agent can
plausibly keep or embody.

**Harder:** Too abstract by itself for app developers. Users still need file-like
views, durable object identifiers, message envelopes, and graph traversal.

**New obligations:** Define layer discipline: promise is the semantic unit;
object/event is the durable unit; message is the boundary-crossing unit;
file-like view is the UX projection.

## Scenario analysis

### S1 - Alice shares a pathname Bob can open

Alice's local view contains `/work/plan`. Bob cannot open that path as global
truth. Under Alt 1, Bob tries to open `/work/plan` directly, which assumes a
global namespace authority. Under Alt 2, Alice sends Bob a promise-bound
reference:

```text
local-name: "/work/plan"              ; Alice's local label, not global truth
target-root: <CID>                    ; content, service, or view root
target-pcid: <pCID>                   ; protocol for interpreting the target
selector: <path-or-query-within-root> ; optional local selection
frontier: <event-or-version-CID>      ; what Alice meant at send time
promiser: Alice or another agent
promise-body: "I promise this reference denotes the plan I intend to share"
evidence: <promise/event refs>
```

Bob then decides whether to trust the promiser and where to mount the reference
in Bob's own view, such as `/from/alice/plan`. The local pathname is UX. The
shareable object is the promise-bound reference.

Alt 2 survives. Alt 3 can carry the reference as a message. Alt 4 can store the
reference as IPLD. Alt 5 explains the semantics. Alt 1 fails the trust model.

### S2 - Shared executable without installation

Carol publishes an executable analysis tool. Dave wants to run it as though it
were a shared mainframe program. In PromiseGrid terms, Carol may promise what
the executable bytes are, what pCID defines the execution interface, and what
tests or attestations she knows. Dave's local compute service may promise to run
that executable under a resource profile and record results.

The mainframe-like experience is: Dave opens or runs a named resource. The
PromiseGrid reality is: Dave trusted a chain of promises enough to map the
executable into his local view and ask a local or remote compute agent to
promise execution.

This supports "no download/install ritual" without pretending that code runs
because a central system permits it.

### S3 - A byte chunk promises something

Ellen observes a byte chunk. Under a narrow software-agent-only reading, the
chunk cannot promise anything; only a human or service can promise about it.
Under the broader Semantic Spacetime reading, the chunk can be modeled as an
agent for simple promises it can plausibly keep or embody:

- "My bytes are exactly this sequence."
- "My CID is this value under this hash and codec."
- "When interpreted by this pCID, I have this shape."
- "I link to these children."

The chunk cannot promise that Bob will store it, Alice will keep it private, or
Carol's computation is honest. Those promises belong to Bob, Alice, and Carol.

This resolves the earlier overcorrection: PromiseGrid should allow non-human and
inert-seeming agents, but every promise body must still be scoped to what that
agent controls or embodies.

### S4 - IPLD representation of the hypergraph

Frank records a promise event involving Alice, Bob, a payload root, a compute
result, a refusal, and a trust update. IPLD can represent this by reifying the
hyperedge as an object:

```text
promise-event = {
  "type": "promise-event",
  "pCID": <protocol CID>,
  "promiser": <agent ref>,
  "promisee": <agent ref>,
  "body": <payload or link>,
  "participants": [<links...>],
  "evidence": [<links...>],
  "parents": [<prior event links...>]
}
```

The IPLD block graph remains a content-addressed DAG of encoded objects. The
logical PromiseGrid structure can still be a hypergraph because hyperedges are
ordinary objects that link to many participants. Logical cycles are represented
through append-only events, named frontiers, or later objects linking to earlier
objects; the immutable block layer does not need to contain cyclic bytes.

This is IPLD-compatible if PromiseGrid avoids relying on IPLD to define trust or
agency. IPLD carries links and bytes. pCID-selected protocols define meaning.

### S5 - Cross-legal-entity trust and broken promises

Alice sees `/partners/bob/data` in her local view because Bob has a history of
keeping storage promises. Later Bob breaks a handling promise. Alice does not
need a global admin to revoke Bob. She records evidence, lowers local trust, and
her view stops projecting Bob's resources into convenient paths.

The system remains mainframe-like for trusted relationships and deliberately
less convenient for untrusted relationships. That is the point: convenience is
earned by promises and trust, not granted by central authority.

### S6 - Long-horizon Git-like event history

Grace reconstructs a project decades later. She has old pCID specs, IPLD/CID
objects, promise/event records, and local view manifests. Modern operating
systems and app packaging conventions have changed. Grace can still rebuild a
view by traversing the promise/event graph, selecting a frontier, and applying
her own trust policy.

The durable object is not a path or URL. It is the content-addressed
promise/event history plus pCID-defined interpretation.

## Cross-cutting findings

- "Decentralized mainframe" is compatible with Promise Theory if it describes
  the user experience, not a global authority model.
- Alice's pathname is local. Bob opens a promise-bound reference that Bob maps
  into Bob's own local view.
- "Everything is a promise" is the right semantic slogan, but it needs layer
  discipline:
  - semantic layer: promise;
  - durable substrate: content-addressed object/event/hyperedge;
  - boundary layer: pCID-selected message;
  - UX layer: file-like local view.
- IPLD can represent the durable substrate, including hypergraphs, by reifying
  hyperedges as linked objects.
- IPLD compatibility does not require PromiseGrid to become IPLD-defined.
  PromiseGrid pCIDs define agency, trust, promise bodies, evidence, and local
  view semantics.
- File-like UX remains valuable. PromiseGrid should not become filesystem-first,
  but it should let agents project trusted promise/event graph regions into
  file-like views.

## Conclusions

- Reject Alt 1 as the conceptual model: a global single-system image hides
  agency and conflicts with local trust.
- Keep Alt 2 as the strongest UX model: local file-like views over
  promise-bound references.
- Keep Alt 3 as the boundary-crossing mechanism: pCID-selected messages carry
  requests, promises, refusals, receipts, and evidence.
- Keep Alt 4 as the likely durable substrate: IPLD-compatible
  content-addressed objects can encode promise/event hypergraphs.
- Keep Alt 5 as the semantic model: everything useful is a promise, but each
  promise must be scoped to what the promiser can control or embody.

The compact framing is:

> PromiseGrid is a decentralized mainframe: each agent gets a coherent
> file-like view over a shared promise/event hypergraph, but that view is local,
> trust-filtered, and assembled from promises rather than imposed by a central
> namespace authority.

## Recommended next DF packet for DR-davod

Before deciding the stable kernel-developer porting boundary, answer:

1. Should guide prose describe PromiseGrid as a decentralized mainframe only if
   it also states that views are local and trust-filtered?
2. Should the app/kernel model include a first-class promise-bound reference
   object for sharing local pathnames across agents?
3. Should "everything useful is a promise" become the semantic slogan, with
   file/message/object/event/graph separated by layer?
4. Should IPLD-compatible objects be treated as the default representation for
   promise/event hypergraph storage, while pCID specs remain the source of
   meaning?
5. Should future sims test local file-like projections over promise/event graph
   frontiers, including broken-promise view decay?

## Implications for open work

- `DR-davod` remains open.
- `TE-dunas` should remain prior-art influence; this TE supplies the local-view
  and promise/event-hypergraph synthesis that followed from it.
- `TE-pudiv` should remain the same-grid app/kernel boundary question; this TE
  adds the file-like local-view and pathname-sharing pressure it should test.
- `TE-hirap`, `TE-vilot`, and `TE-gurov` remain apparatus/specimen guidance;
  this TE should not force every repo artifact to become a PromiseGrid message.
- `SIM-fovip` or a successor sim should test promise-bound references, local
  file-like projections, IPLD-compatible hyperedge objects, and trust-filtered
  view changes after broken promises.
- A later design note can explain the decentralized-mainframe metaphor in plain
  English after this TE's DF questions are answered or narrowed.
- A whitepaper should wait until at least one sim exercises the model.

## References

- Mark Burgess, "Spacetimes with Semantics" and Semantic Spacetime materials:
  agents, observer-local namespaces, adjacency, and local interpretation.
- Mark Burgess, Promise Theory materials: autonomous agents may be human,
  machine, mineral, or otherwise; no agent promises on behalf of another.
- IPLD Data Model and DAG-CBOR documentation: maps, lists, bytes, links, CID
  links, canonical DAG-CBOR encoding, and tag `42` link representation.

## Decision status

`needs DF` - this TE captures the decentralized-mainframe, local-view,
promise-first, IPLD-compatible hypergraph synthesis, but it does not decide the
final kernel boundary, reference object shape, L6 CAS profile, or guide prose.
