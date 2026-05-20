# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260520-205857`
- Child ID: `SIM-buvak-ga-child-0002`
- Child path: `simulations/SIM-buvak-ga-child-0002/`
- Operation: `breed`
- Parent IDs: `SIM-ludut-wire-lab-devs, SIM-ligan-promisebase-reference-naming`

## Scenario Sample

- `bgp-class-routing-app-partition-and-stale-path` at `scenarios/bgp-class-routing-app-partition-and-stale-path/bgp-class-routing-app-partition-and-stale-path.md`
- `udp-feed-v0-conformance-loopback-round-trip` at `scenarios/udp-feed-v0-conformance-loopback-round-trip/udp-feed-v0-conformance-loopback-round-trip.md`
- `cas-backed-group-session-missing-pointee` at `scenarios/cas-backed-group-session-missing-pointee/cas-backed-group-session-missing-pointee.md`

## Scenario Pressure

### `scenarios/bgp-class-routing-app-partition-and-stale-path/bgp-class-routing-app-partition-and-stale-path.md`

```markdown
# Partition and stale path

## Scenario ID

bgp-class-routing-app-partition-and-stale-path

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-punaz-bgp-class-routing-app/SCENARIOS.md`
- Source simulation: `SIM-punaz-bgp-class-routing-app/`
- Source row/title: Partition and stale path
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-punaz-bgp-class-routing-app/`.

## Setup

A formerly good path becomes unavailable during intermittent connectivity.

## Stimulus

Run the candidate simulation against this source test: How stale promises, timeouts, and withdrawal notices affect local decisions.

## Expected Pressure

Long-lived routing records need aging and repair without central convergence machinery.
```

### `scenarios/udp-feed-v0-conformance-loopback-round-trip/udp-feed-v0-conformance-loopback-round-trip.md`

```markdown
# Loopback round trip

## Scenario ID

udp-feed-v0-conformance-loopback-round-trip

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kuful-udp-feed-v0-conformance/SCENARIOS.md`
- Source simulation: `SIM-kuful-udp-feed-v0-conformance/`
- Source row/title: Loopback round trip
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kuful-udp-feed-v0-conformance/`.

## Setup

Alice sends a 612-byte payload to Bob over local UDP.

## Stimulus

Run the candidate simulation against this source test: Whether the reference implementation preserves bytes and exposes the expected send/receive API.

## Expected Pressure

A minimal reference may be enough for first v0 evidence if vectors lock the bytes.
```

### `scenarios/cas-backed-group-session-missing-pointee/cas-backed-group-session-missing-pointee.md`

```markdown
# Missing pointee

## Scenario ID

cas-backed-group-session-missing-pointee

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jurar-cas-backed-group-session/SCENARIOS.md`
- Source simulation: `SIM-jurar-cas-backed-group-session/`
- Source row/title: Missing pointee
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jurar-cas-backed-group-session/`.

## Setup

Bob sees pointer object CID Y but lacks root CID X or some child chunks.

## Stimulus

Run the candidate simulation against this source test: Whether the group view can show pending / unresolved state without treating the message as invalid.

## Expected Pressure

Sparse-CAS behavior must be a normal group-session state.
```

## Parent Simulation Documents

### `simulations/SIM-ludut-wire-lab-devs/README.md`

```markdown
# SIM-ludut: Wire-lab devs

This simulation carries the concrete `wire-lab-devs` developer-coordination
lineage as one candidate source of PromiseGrid design evidence. It keeps that
lineage available for comparison and evolution without treating its current
layout or mechanics as the final PromiseGrid product shape. Source:
`DI-limom`, `DI-rugig`.

Turn 193 is the urgency source for this lineage: Steve wanted the developer
group dogfooding message transport ASAP so he was not working solo. This
simulation is the durable home for that concrete dogfood evidence while protocol
and migration choices continue to evolve independently. Source: `DI-vuzot`.
```

### `simulations/SIM-ludut-wire-lab-devs/QUESTION.md`

```markdown
# Question

Does the `wire-lab-devs` developer-coordination lineage continue to produce
useful PromiseGrid design evidence when it is treated as one independent
simulation rather than as the default repo layout? Source: `DI-limom`,
`DI-rugig`.

Turn 193 adds the operational pressure: can this lineage support near-term
developer dogfooding without prematurely selecting a final root layout, CAS
migration shape, or production PromiseGrid group protocol? Source: `DI-vuzot`.
```

### `simulations/SIM-ludut-wire-lab-devs/seed/wire-lab-devs-draft-migration.md`

```markdown
# wire-lab-devs-draft Migration

The `wire-lab-devs-draft` transport evidence moved into this simulation's world
so it can be replayed and evaluated as specimen data without preserving root
`transports/` as an active layout commitment. Source: `DI-fakin`.

| Field | Value |
|---|---|
| Original path | `transports/wire-lab-devs-draft/` |
| New path | `simulations/SIM-ludut-wire-lab-devs/world/transports/wire-lab-devs-draft/` |
| Method | `git mv` |
| Source commit | `780f56525a8d528d3d5caf58ab18f9a7f41da892` |
| CID parameters | CIDv1, raw codec, sha2-256 multihash, base32 multibase |
| Verification result | PASS on 2026-05-10: all four `bafk*.txt` filenames matched raw CIDv1 over file bytes after migration. |

## Verified message CIDs

| Message file | Verification |
|---|---|
| `bafkreia46vxsahmeicugfxmc7natorkstc3mdaz4r5d3zz46whjwpvqwta.txt` | PASS |
| `bafkreidef4b4qdc4xjvkjrern7jm4ta75q55ed2u2ilwcrkxqhn7n4fjce.txt` | PASS |
| `bafkreihhuejiefrqrm7zgw2jsdqc37lwmbvfkw5uqbnjx3wsobcxh3y7ni.txt` | PASS |
| `bafkreihnonvsf3vmcagukqcxwoh35255eduulvwwx3kax6ty4iidklk5vu.txt` | PASS |

The message files are not edited by this migration. Their body text may mention
old paths such as `transports/draft--wire-lab-devs/`; those references are
historical evidence and are preserved to keep CIDs stable.
```

### `simulations/SIM-ligan-promisebase-reference-naming/README.md`

```markdown
# SIM-ligan: Promisebase reference naming

This simulation explores promisebase's human-readable symbol pointing at a hash
problem as its own reference-resolution design question. It is split out from
the CAS object-model simulation so reference names can evolve independently from
CBOR profile, chunking, CIDv1 object typing, and pointer-object byte shape.
Source: `DI-tibis`.

## Question

How should PromiseGrid handle human-readable references to CAS roots learned
from promisebase without creating competing identities, mutable-name confusion,
or custom non-CID syntax?
Source: `DI-tibis`.

## Turn 178 pressure

Turn 178 surfaced promisebase as prior art and specifically called out a prior
reference problem: how should a human-readable symbol point at a hash, especially
when the earlier prototype did not use CIDs and invented custom syntax?

This simulation treats that as a separate question from the core CAS object
model. A CAS object model decides bytes, codecs, chunks, and pointer-object
identity. Reference naming decides whether and how humans get stable or mutable
names for roots without confusing the name with the CID it names. Source:
`DI-tibis`.

## Decision Axes

- **Layer home:** L6 reference objects, L7 group/session metadata, a separate
  reference protocol, or explicit rejection/deferment.
- **Mutability:** immutable labels, mutable refs with history, signed updates,
  or time-scoped aliases.
- **Identity boundary:** how to keep a reference name distinct from the pointer
  object CID and the root CID.
- **Conflict handling:** collisions, squatting, divergent local refs, and
  malicious updates.
- **Interop:** whether the shape should align with git refs, IPNS-like naming,
  AT-style handles, or avoid importing any existing naming model prematurely.

## Boundaries

This simulation does not choose the CBOR profile, chunking algorithm, CID codec
set, or pointer-object shape for L6 CAS. Those remain with
`SIM-jomag-cas-object-model` and TODO-kituj / `DR-tumus`. This simulation only
keeps the reference-symbol / hash-name-resolution question visible as its own
design lineage. Source: `DI-tibis`.
```

### `simulations/SIM-ligan-promisebase-reference-naming/QUESTION.md`

```markdown
# Question

How should PromiseGrid handle human-readable references to CAS roots learned
from promisebase without creating competing identities, mutable-name confusion,
or custom non-CID syntax?
Source: `DI-tibis`.

Open decision points:

- Does reference naming belong in L6 CAS objects, L7 group/session metadata, a
  separate reference protocol, or nowhere in the first profile?
- Are references immutable labels, mutable refs with signed history, local-only
  nicknames, or discoverable shared names?
- How does a reader distinguish the reference name, the reference object's CID,
  a pointer-object CID, and the target root CID?
- What happens when Alice and Bob use the same human-readable symbol for
  different roots, or when Mallory publishes a confusingly similar name?
- Which promisebase reference lessons should be adopted, wrapped, rejected, or
  preserved only as prototype history?
```

### `simulations/SIM-ligan-promisebase-reference-naming/SCENARIOS.md`

```markdown
# Promisebase Reference Naming Scenarios

These scenarios isolate promisebase's human-readable reference problem from the
broader CAS object-model work. They are design pressure, not a final naming
protocol. Source: `DI-tibis`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Immutable local nickname | Alice locally calls root CID X `spring-paper` without publishing that name. | Whether local-only convenience names can stay outside protocol identity. | The first profile may avoid reference protocol scope if local nicknames are enough. |
| Published mutable ref | Alice publishes `project/latest` first pointing at root X and later at root Y. | How signed update history, replay protection, and reader expectations work. | Mutable refs are not CAS roots; they need explicit update semantics if adopted. |
| Reference object as CAS object | Alice writes a CBOR reference object whose CID R points at root X and includes a human-readable label. | Whether the reference object is L6 CAS, L7 metadata, or a separate protocol object. | Reference-object identity must not be confused with the target root identity. |
| Name collision | Alice and Bob both publish `release` pointing at different roots. | Whether names are scoped by peer, group, pCID, site, or another authority. | Human-readable names need scope rules or they become a new central registry problem. |
| Maliciously similar name | Mallory publishes a visually confusing or policy-confusing name near Alice's ref. | Whether readers need normalization, warnings, or refusal policy. | Reference naming imports UI and social-trust risks that raw CIDs avoid. |
| Promisebase custom syntax migration | A promisebase-era reference uses non-CID custom syntax for a root. | Whether migration wraps it, rejects it, or maps it into CID-backed reference objects. | Prior-art adoption must be deliberate and not preserve known-bad syntax by accident. |

## Expected Outputs

- A decision packet for TODO-kituj `kituj.5` on whether reference naming belongs
  in L6, L7, a separate protocol, or a later deferral.
- A distinction table for reference name, reference-object CID, pointer-object
  CID, and target root CID.
- A prior-art disposition list for promisebase reference syntax and behavior.
```

## Compact Fitness Evidence From This Run

### `SIM-ludut-wire-lab-devs` x `bgp-class-routing-app-partition-and-stale-path`

- Result path: `results/SIM-ludut-wire-lab-devs/bgp-class-routing-app-partition-and-stale-path/openai-gpt-5.4-xhigh/20260520-205857.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=3 evolution_safety=4 layer_boundary_clarity=4 failure_handling=0 implementation_plausibility=1 risk_penalty=2
- Fitness: raw=13.00 normalized_0_100=37.00 confidence_0_1=0.80
- Rationale: The docs show strong evidence preservation, migration discipline, and explicit non-final boundaries, but they do not describe route aging, timeout or withdrawal handling, or peer-local repair under partition. Fit to this stale-path routing scenario is therefore weak and mostly indirect.
- Strengths:
  - CID-verified migration evidence makes the preserved transport artifacts traceable.
  - The simulation explicitly treats current layout and mechanics as non-final, which helps evolution safety.
  - Evidence and authority boundaries are clearly stated.
- Weaknesses:
  - No explicit route or path model is described.
  - No timeout, stale-promise aging, or withdrawal-notice behavior is specified.
  - The docs do not show what peers can record locally after a partition.
- Risks:
  - A reviewer could over-read developer dogfood transport evidence as proof of stale-path routing behavior.
  - Hidden manual coordination assumptions may exist because failure handling is not specified.
- Open questions:
  - What local record marks a once-good path as stale or withdrawn?
  - How would peers age out bad path information without central convergence machinery?
  - Can the preserved transport lineage be extended into explicit partition-and-repair tests?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-ludut-wire-lab-devs` x `udp-feed-v0-conformance-loopback-round-trip`

- Result path: `results/SIM-ludut-wire-lab-devs/udp-feed-v0-conformance-loopback-round-trip/openai-gpt-5.4-xhigh/20260520-205857.json`
- Scores: scenario_fit=1 promisegrid_alignment=3 auditability=3 evolution_safety=4 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=2 risk_penalty=3
- Fitness: raw=15.00 normalized_0_100=43.00 confidence_0_1=0.79
- Rationale: Good as migration and byte-stability evidence, weak as direct UDP loopback conformance evidence. The docs show preserved transport artifacts and CID-checked bytes, but not a UDP v0 reference implementation, the 612-byte loopback test, or an expected send/receive API.
- Strengths:
  - CID-checked migrated message files provide concrete byte-preservation evidence.
  - The simulation clearly separates lineage evidence from final PromiseGrid product or layout commitments.
  - Existing transport artifacts make near-term developer dogfooding more plausible than a purely hypothetical design.
- Weaknesses:
  - No explicit UDP loopback behavior or 612-byte round-trip proof is documented.
  - No send/receive API contract or conformance vector set is shown.
  - Failure, adversarial, and peer-local receipt accounting for network behavior are largely absent.
- Risks:
  - File-level byte preservation could be mistaken for UDP conformance if the boundary is not kept explicit.
  - Urgent dogfood pressure could freeze an accidental transport or interface boundary before it is fully specified.
- Open questions:
  - Does the migrated transport include executable UDP send/receive code or only specimen files?
  - What exact API and fixed vectors define udp-feed-v0 conformance for replay across implementations?
  - How are timeout, truncation, bind failure, duplication, or reordering recorded as durable local evidence?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-ludut-wire-lab-devs` x `cas-backed-group-session-missing-pointee`

- Result path: `results/SIM-ludut-wire-lab-devs/cas-backed-group-session-missing-pointee/openai-gpt-5.4-xhigh/20260520-205857.json`
- Scores: scenario_fit=1 promisegrid_alignment=3 auditability=3 evolution_safety=4 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=2 risk_penalty=2
- Fitness: raw=16.00 normalized_0_100=46.00 confidence_0_1=0.68
- Rationale: This lineage is auditable, evolution-friendly transport evidence, but on the provided docs it is mostly orthogonal to the missing-pointee scenario: it preserves CID-named messages and migration history yet does not define a group view, pointer graph, or pending/unresolved treatment for absent CAS content.
- Strengths:
  - Concrete dogfood transport evidence is preserved as a separate simulation instead of hard-coding current repo layout as the final product shape.
  - The migration record is auditable through source path, source commit, CID parameters, and explicit post-move CID verification.
  - The docs clearly preserve room to evolve transport, CAS layout, and group protocol choices later.
- Weaknesses:
  - No explicit group-session object model, pointer/root/chunk structure, or group-view semantics are described.
  - The docs do not show how partially available CAS content stays visible as pending or unresolved rather than being treated as invalid.
  - Peer-local promise accounting and later audit records for Bob or Carol under missing-content conditions are not specified.
  - ... 1 more
- Risks:
  - Readers may over-read preserved transport evidence as coverage of group-session semantics that are not actually specified.
  - A real implementation derived too directly from this lineage could drop or invalidate partially fetched messages instead of surfacing sparse-CAS state.
  - Near-term dogfood shortcuts could harden into de facto protocol expectations before missing-pointee behavior is designed.
- Open questions:
  - What object or feed structure links the CID-named message files into a group session view?
  - How should Bob record and display pointer CID Y when root CID X or child chunks are missing?
  - What durable local evidence should Carol inspect to distinguish pending, permanently missing, and maliciously withheld content?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-ligan-promisebase-reference-naming` x `bgp-class-routing-app-partition-and-stale-path`

- Result path: `results/SIM-ligan-promisebase-reference-naming/bgp-class-routing-app-partition-and-stale-path/openai-gpt-5.4-xhigh/20260520-205857.json`
- Scores: scenario_fit=0 promisegrid_alignment=2 auditability=2 evolution_safety=4 layer_boundary_clarity=5 failure_handling=1 implementation_plausibility=1 risk_penalty=3
- Fitness: raw=12.00 normalized_0_100=34.00 confidence_0_1=0.90
- Rationale: This simulation is a naming-boundary study, not a routing-path design. It is strong on separating human-readable references from CAS identity and on migration safety, but it does not define path aging, withdrawal, timeout, or partition-repair behavior. That makes it useful mainly as negative scope evidence for this stale-path routing scenario.
- Strengths:
  - Very clear distinction between reference names, reference objects, and target/root identity.
  - Strong layer-boundary discipline; naming is explicitly separated from CAS object-model choices.
  - Good evolution posture through independent simulation scope and explicit legacy-syntax migration pressure.
- Weaknesses:
  - Does not model route selection, stale path aging, withdrawal notices, or local repair after partition.
  - Failure cases center on naming collisions and misleading names rather than routing liveness.
  - Provides little peer-local accounting for how Alice or Bob should decide when a formerly good path disappears.
- Risks:
  - Mutable reference history could be mistaken for adequate routing freshness semantics.
  - Reusing naming objects for route state would blur identity semantics with reachability semantics.
  - Partition behavior is left undefined for this scenario, so stale local decisions could persist.
- Open questions:
  - If PromiseGrid needs route-like mutable announcements, should they live in a separate protocol rather than reference naming?
  - What expiry, withdrawal, or freshness proofs would be required for any dynamic state under intermittent connectivity?
  - How should peer-local evidence distinguish stale reference history from currently usable path state?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

## Required JSON Shape

{"child_id":"SIM-buvak-ga-child-0002","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
