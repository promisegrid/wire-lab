# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260520-205857`
- Child ID: `SIM-zukor-ga-child-0001`
- Child path: `simulations/SIM-zukor-ga-child-0001/`
- Operation: `breed`
- Parent IDs: `SIM-ludut-wire-lab-devs, SIM-kuful-udp-feed-v0-conformance`

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

### `simulations/SIM-kuful-udp-feed-v0-conformance/README.md`

```markdown
# SIM-kuful: UDP-Feed v0 Conformance

This simulation captures the protocol/specimen questions in `TODO-jodon`: what
minimal reference implementation, test-vector, simulation-artifact, and ns-3
harness surface proves UDP-feed v0 well enough to support later PromiseGrid
design work. It is a standalone design-point simulation, not an implementation
and not a frozen UDP-feed spec. Source: `DI-pukap`.

## Question

What conformance surface should UDP-feed v0 require before the lineage is useful
as a binding specimen: Go reference implementation, deterministic test vectors,
simulator artifact writer, ns-3 round trip, or a smaller/larger combination?
Source: `DI-pukap`; `TODO-jodon`.

## Candidate Shapes

- **Reference-first:** A minimal Go implementation establishes behavior, then
  test vectors and ns-3 scenarios validate it.
- **Vector-first:** Test vectors define the behavior before the implementation
  claims conformance.
- **Harness-first:** An ns-3 round trip proves the binding layer in realistic
  network conditions before polishing API shape.
- **Layer-composition proof:** UDP-feed is not done until a session protocol can
  ride above it and record conformance claims.

## Boundaries

This simulation does not write the Go implementation, test vectors, or ns-3
harness. It captures the conformance-design question so `TODO-jodon` can remain
implementation-focused while still exploring what evidence should count as a
usable v0 binding artifact. Source: `DI-pukap`.
```

### `simulations/SIM-kuful-udp-feed-v0-conformance/QUESTION.md`

```markdown
# Question

Which minimal evidence should prove UDP-feed v0: a Go reference implementation,
test vectors, simulator artifact output, ns-3 end-to-end round trip, and/or a
session-layer composition test? Source: `DI-pukap`; `TODO-jodon`.

Open decision points:

- Should test vectors be authored before or after the reference implementation?
- Is a loopback-only implementation sufficient for v0, or must ns-3 prove the
  binding through an emulated network?
- What simulation-artifact path and metadata prove promise 10 without becoming a
  production API?
- Which conformance claims belong in implementation `CHANGELOG.md` records?
- Does UDP-feed v0 need a group/session layer above it before it counts as a
  useful binding specimen?
```

### `simulations/SIM-kuful-udp-feed-v0-conformance/SCENARIOS.md`

```markdown
# UDP-Feed v0 Conformance Scenarios

These scenarios are evidence for `TODO-jodon`. They are not an implementation
and not a frozen UDP-feed spec. Source: `DI-pukap`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Loopback round trip | Alice sends a 612-byte payload to Bob over local UDP. | Whether the reference implementation preserves bytes and exposes the expected send/receive API. | A minimal reference may be enough for first v0 evidence if vectors lock the bytes. |
| Boundary-size payload | Alice sends exactly 1232 bytes, then 1233 bytes. | Whether the implementation honors the size promise and errors locally before oversize send. | Size behavior should be in vectors before wider conformance claims. |
| Malformed datagram | Bob receives arbitrary bytes that do not parse at higher layers. | Whether UDP-feed passes bytes upward unchanged rather than inventing message semantics. | Binding conformance must stay below session semantics. |
| Simulation artifact output | A simulator-mode send writes an artifact file for the transmitted bytes. | Whether artifact output proves promise 10 without becoming production behavior. | The artifact contract needs to be testable and explicitly scoped. |
| ns-3 two-node path | Alice and Bob communicate through an ns-3-emulated UDP network. | Whether the v0 reference survives non-loopback timing, interface, and packet-capture conditions. | ns-3 may be the evidence that separates a useful specimen from a local toy. |
| Session-layer composition | A minimal group/session message rides above UDP-feed v0. | Whether UDP-feed's API is sufficient for the next layer without leaking binding details. | If composition is required, TODO-jodon's done criteria must include more than UDP round trip. |

## Expected Outputs

- Evidence for TODO-jodon's done criteria and for whether TODO-bihon's ns-3
  scaffold is required before UDP-feed v0 is considered usable.
- A conformance checklist covering implementation API, vectors, artifact writer,
  ns-3 proof, and implementation conformance claims.
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

### `SIM-kuful-udp-feed-v0-conformance` x `bgp-class-routing-app-partition-and-stale-path`

- Result path: `results/SIM-kuful-udp-feed-v0-conformance/bgp-class-routing-app-partition-and-stale-path/openai-gpt-5.4-xhigh/20260520-205857.json`
- Scores: scenario_fit=0 promisegrid_alignment=2 auditability=2 evolution_safety=2 layer_boundary_clarity=5 failure_handling=0 implementation_plausibility=2 risk_penalty=2
- Fitness: raw=11.00 normalized_0_100=31.00 confidence_0_1=0.89
- Rationale: This simulation is a clear, lower-layer UDP-feed conformance study, not a routing-path aging study. It has strong boundary discipline and plausible conformance evidence, but it does not directly model stale promises, withdrawals, timeouts, or local path repair under partition.
- Strengths:
  - Explicitly keeps UDP-feed below session and routing semantics, which prevents overclaiming.
  - Defines concrete conformance surfaces such as vectors, artifact output, ns-3 round trip, and composition probing.
  - Could serve as precursor evidence for a later routing-layer simulation.
- Weaknesses:
  - No routing table, advertisement, withdrawal, or best-path decision model.
  - No peer-local record of stale-path detection or repair during intermittent connectivity.
  - Covers PromiseGrid failure pressure only indirectly at the transport/specimen layer.
- Risks:
  - Network realism from ns-3 could be mistaken for evidence of stale-path handling.
  - A composition test could overstate routing readiness unless timer, withdrawal, and aging semantics are added.
- Open questions:
  - Should a follow-on simulation layer route advertisement, withdrawal, and aging behavior over UDP-feed?
  - What local artifacts would prove stale-path repair without central convergence machinery?
  - Is partition fault injection better added here, or in a distinct routing-application simulation?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kuful-udp-feed-v0-conformance` x `udp-feed-v0-conformance-loopback-round-trip`

- Result path: `results/SIM-kuful-udp-feed-v0-conformance/udp-feed-v0-conformance-loopback-round-trip/openai-gpt-5.4-xhigh/20260520-205857.json`
- Scores: scenario_fit=4 promisegrid_alignment=2 auditability=4 evolution_safety=3 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=3 risk_penalty=2
- Fitness: raw=22.00 normalized_0_100=63.00 confidence_0_1=0.79
- Rationale: Directly fits the loopback conformance pressure and keeps the UDP-feed/session boundary clean, but it is still a planning artifact for what evidence should count rather than a concrete conformance specimen, so broader PromiseGrid goal checks and executable proof remain limited.
- Strengths:
  - Direct coverage of loopback, vectors, ns-3, and composition pressures in one conformance decision space.
  - Strong boundary hygiene: UDP-feed stays below session semantics, and specimen artifacts are kept distinct from production API.
  - The question, scenario matrix, and open decisions make the evidence gap easy to audit.
- Weaknesses:
  - No concrete send/receive API, byte-level vector set, or artifact schema is defined for the 612-byte round trip.
  - The simulation does not explicitly walk the 100-year, sparse-knowledge, or no-central-authority checks for this scenario.
  - Failure handling is acknowledged, but not yet resolved into a single conformance contract.
- Risks:
  - Loopback-only success could be mistaken for sufficient v0 evidence before ns-3 or composition proof exists.
  - Because this is a design-of-evidence simulation, readers may overread planning text as normative protocol behavior.
- Open questions:
  - Must vectors lead implementation, or can the reference implementation define the first specimen?
  - Is loopback enough for v0, or is an ns-3 round trip required?
  - What artifact metadata proves promise 10 without creating a production API?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kuful-udp-feed-v0-conformance` x `cas-backed-group-session-missing-pointee`

- Result path: `results/SIM-kuful-udp-feed-v0-conformance/cas-backed-group-session-missing-pointee/openai-gpt-5.4-xhigh/20260520-205857.json`
- Scores: scenario_fit=1 promisegrid_alignment=1 auditability=1 evolution_safety=1 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=1 risk_penalty=3
- Fitness: raw=12.00 normalized_0_100=30.00 confidence_0_1=0.89
- Rationale: This simulation is mainly a UDP-feed v0 transport-conformance design. It helps clarify that missing CAS pointees are an upper-layer concern, but it does not model CAS objects, unresolved group state, or the local evidence needed to judge this scenario directly.
- Strengths:
  - Clear transport/session boundary reduces pressure to smuggle group semantics into UDP-feed.
  - Vectors, artifacts, ns-3, and composition smoke-test ideas provide auditable lower-layer evidence.
- Weaknesses:
  - No CAS-backed session model, pointer fixtures, or rules for pending versus invalid group state.
  - Sparse-knowledge behavior is not exercised as a normal steady-state condition.
  - Peer-local promise accounting for unresolved pointees is absent.
- Risks:
  - A successful transport or composition check could be misread as evidence that missing-pointee handling is solved.
  - Overclaiming from this layer would blur the boundary between byte carriage and group-session semantics.
- Open questions:
  - Should UDP-feed v0 conformance include a concrete CAS-backed session fixture, or stay strictly below that layer?
  - What artifact or test vector would prove pending/unresolved state without treating the message as invalid?
  - If composition is required, where should unresolved-pointee behavior be specified?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

## Required JSON Shape

{"child_id":"SIM-zukor-ga-child-0001","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
