# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Choose a descriptive `child_id` that starts with `SIM-maraz-child-` and ends with a kebab-case design slug. Do not use generic `ga-child`, `pending`, or ordinal-only names.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260521-003110`
- Planned child ID prefix: `SIM-maraz-child-`
- Temporary child ID: `SIM-maraz-child-pending-0002`
- Temporary child path: `proposals/ga-canary-20260521-003110/simulations/SIM-maraz-child-pending-0002/`
- Operation: `breed`
- Parent IDs: `SIM-janov-grid-envelope-layer-pcid-nested-signed-payload, SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields`

## Scenario Sample

- `chunk-feed-replication-sparse-advertisement` at `scenarios/chunk-feed-replication-sparse-advertisement/chunk-feed-replication-sparse-advertisement.md`
- `conditional-release-geofencing-onward-restraint-chain` at `scenarios/conditional-release-geofencing-onward-restraint-chain/conditional-release-geofencing-onward-restraint-chain.md`
- `bgp-routing` at `scenarios/bgp-routing/bgp-routing.md`

## Scenario Pressure

### `scenarios/chunk-feed-replication-sparse-advertisement/chunk-feed-replication-sparse-advertisement.md`

```markdown
# Sparse advertisement

## Scenario ID

chunk-feed-replication-sparse-advertisement

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`
- Source simulation: `SIM-zazit-chunk-feed-replication/`
- Source row/title: Sparse advertisement
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zazit-chunk-feed-replication/`.

## Setup

Alice has a subset of chunks for a Merkle root; Bob has a different subset.

## Stimulus

Run the candidate simulation against this source test: Whether the feed advertises leaves, roots, pointer objects, frontiers, or compact summaries without assuming full replication.

## Expected Pressure

Feed specs must work when no site has all CAS objects.
```

### `scenarios/conditional-release-geofencing-onward-restraint-chain/conditional-release-geofencing-onward-restraint-chain.md`

```markdown
# Onward-restraint chain

## Scenario ID

conditional-release-geofencing-onward-restraint-chain

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zarud-conditional-release-geofencing/SCENARIOS.md`
- Source simulation: `SIM-zarud-conditional-release-geofencing/`
- Source row/title: Onward-restraint chain
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zarud-conditional-release-geofencing/`.

## Setup

Alice sends content to Bob only if Bob promises to forward it only to recipients who make the same promise. Bob wants to forward to Carol.

## Stimulus

Run the candidate simulation against this source test: Whether the recursive promise graph is represented at group-session, conditional-release, transport/feed, or CAS-object level.

## Expected Pressure

If the graph is central to dispatch semantics, group/session ownership gets stronger; if it composes across sessions, a separate family gets stronger.
```

### `scenarios/bgp-routing/bgp-routing.md`

```markdown
# BGP Routing

## Scenario ID

bgp-routing

## Source / Provenance

- Source type: application seed
- Source path: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`
- Source row/title: Seed application catalog entry `bgp-routing`
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-midif`; `TODO-dadub`

## Purpose

Exercise PromiseGrid design candidates against bgp-routing application pressure: Route-
like reachability promises, hijacks, leaks, stale paths, and sparse topology knowledge.

## Setup

Alice depends on an outcome in the BGP Routing domain. Bob makes promises about route-
like reachability promises, hijacks, leaks, stale paths, and sparse topology knowledge.
Carol needs enough evidence to rely on Bob's promise without having complete global
state, and Mallory may exploit stale, missing, or ambiguous evidence.

## Stimulus

A routine application event becomes contested or incomplete. Some relevant objects,
signatures, observations, or relationship records are delayed, partitioned, stale, or
disputed, and each peer must decide what to accept, retry, downgrade, or escalate using
only local evidence.

## Expected Pressure

The candidate simulation must show which promises, CAS objects, feeds, identity claims,
names, and promise accounting records are needed for this application pressure, while
avoiding hidden global state or a central authority that would make the result non-
comparable.
```

## Parent Simulation Documents

### `simulations/SIM-janov-grid-envelope-layer-pcid-nested-signed-payload/README.md`

```markdown
# SIM-janov-grid-envelope-layer-pcid-nested-signed-payload: Grid-envelope nested signed payload probe

This simulation is a standalone grid-envelope specimen. It tests a shared
layer/network pCID in the outer envelope, where that pCID defines a nested
payload structure containing the actual payload pCID, actual payload bytes, and
a signature over the nested payload. It does not claim that this layering is the
canonical PromiseGrid wire format. Source: `DI-joman`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
```

### `simulations/SIM-janov-grid-envelope-layer-pcid-nested-signed-payload/QUESTION.md`

```markdown
# Question

Can a grid envelope keep the outer shape minimal as `[pcid_a, payload_a]` while
letting `pcid_a` define a network/layer payload that all participating nodes
parse as `[pcid_b, payload_b, signature_b]`, or does the unsigned outer layer
make conformance and authorship too dependent on transport identity? Source:
`DI-joman`.
```

### `simulations/SIM-janov-grid-envelope-layer-pcid-nested-signed-payload/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries
yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen created by `DI-joman`.
```

### `simulations/SIM-janov-grid-envelope-layer-pcid-nested-signed-payload/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid-envelope draft: layer pCID with nested signed payload

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `layer-pcid-nested-signed-payload`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a
specimen inside `SIM-janov-grid-envelope-layer-pcid-nested-signed-payload`, not
a harness rule and not the canonical PromiseGrid envelope. Source: `DI-joman`.

## Envelope Shape

The outer envelope shape for this variant is:

```text
[pcid_a, payload_a]
```

Slots are interpreted positionally:

- `pcid_a` identifies a shared layer, ecosystem, or network protocol whose
  participants expect to parse `payload_a`.
- `payload_a` is opaque bytes at the outer envelope layer, but the `pcid_a`
  protocol defines its internal structure.

For this candidate, `pcid_a` defines `payload_a` as the canonical bytes of:

```text
[pcid_b, payload_b, signature_b]
```

Nested slots are interpreted by the `pcid_a` protocol:

- `pcid_b` identifies the actual payload protocol for the application data.
- `payload_b` is the actual application or evidence payload.
- `signature_b` is the Layer-A-required signature for the nested payload.

The nested signature covers the canonical bytes of `[pcid_b, payload_b]` in
this draft. That stricter coverage binds the actual payload bytes to their
payload protocol and avoids replaying the same bytes under a different
`pcid_b`.

## Encoding

Both the outer envelope and the `payload_a` interior are deterministic CBOR
positional arrays. `pcid_a` and `pcid_b` are CIDv1 byte strings. `payload_a`,
`payload_b`, and `signature_b` are byte strings at the layer that carries them.
The outer envelope's canonical bytes are the deterministic CBOR bytes of
`[pcid_a, payload_a]`; the nested signature's canonical bytes are the
deterministic CBOR bytes of `[pcid_b, payload_b]`.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid_a`, it cannot participate in the
Layer-A network semantics. It may preserve the full outer envelope bytes as
opaque evidence, but it MUST NOT claim that `payload_a` has the nested shape or
that `signature_b` verifies.

If a receiver understands `pcid_a` but lacks a handler for `pcid_b`, it can
verify the nested signature over `[pcid_b, payload_b]` if the signature scheme
is known through `pcid_a`, but it MUST mark the actual payload interpretation as
unsupported.

## Signature and Authorship Policy

The outer `[pcid_a, payload_a]` layer has no fixed signature slot. This is a
deliberate pressure point for the simulation: `pcid_a` promises that
`payload_a` has a parseable signed nested structure, but the envelope itself
does not prove who made that conformance promise.

Participants therefore rely on transport identity, local peer-adoption records,
or other surrounding context to decide which agent is promising that
`payload_a` conforms to `pcid_a`. The nested `signature_b` authenticates the
actual payload claim, but it does not by itself authenticate the outer
transport context or the sender's promise that `payload_a` is a valid Layer-A
message.

## Layering-Test Behavior

This variant tests whether a commonly shared `pcid_a` can act as the layer or
network contract that most participating nodes understand:

- Generic outer tooling only needs to parse `[pcid_a, payload_a]`.
- Layer-A nodes can parse `payload_a` and find `pcid_b`, `payload_b`, and
  `signature_b`.
- Variable arity is pushed inside `payload_a`, where `pcid_a` defines it,
  rather than changing the universal outer envelope shape.
- Verification is strong for the nested payload if `signature_b` covers
  `[pcid_b, payload_b]`, but weak for the unsigned outer conformance promise
  unless transport or peer context supplies authorship.

## Open Questions

- Is relying on transport identity for the `pcid_a` conformance promise
  acceptable, or does the outer layer need its own signature/proof slot?
- Should `signature_b` cover only `payload_b`, or should it cover `[pcid_b,
  payload_b]` as this draft proposes?
- Should `pcid_a` identify a broad network/ecosystem shared by most nodes, or
  should it be a narrower layer protocol adopted by a subset of peers?
- Does this pattern improve variable-arity flexibility enough to justify the
  extra nested parsing step?

## Non-Canonical Status

This draft does not declare a winning envelope, does not define a central pCID
registry, and does not constrain sibling simulations. It exists to compare a
layer-pCID nested signed payload against fixed-field and variable-outer-arity
alternatives.
```

### `simulations/SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields/README.md`

```markdown
# SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields: Grid-envelope variable-arity probe

This simulation is a standalone grid-envelope specimen. It tests an outer
envelope where the first `pcid` defines how many fields follow and what each
field means. It does not claim that variable arity is the canonical
PromiseGrid wire format. Source: `DI-joman`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
```

### `simulations/SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields/QUESTION.md`

```markdown
# Question

Can a grid envelope safely use a first-slot `pcid` to define the number,
position, and type of all outer fields that follow it, while still satisfying
generic parsing, unknown-pCID handling, signing, audit, and 100-year
interoperability pressures? Source: `DI-joman`.
```

### `simulations/SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries
yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen created by `DI-joman`.
```

### `simulations/SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid-envelope draft: variable-arity pCID-defined fields

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `variable-arity-pcid-defined-fields`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a
specimen inside `SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields`,
not a harness rule and not the canonical PromiseGrid envelope. Source:
`DI-joman`.

## Envelope Shape

The envelope shape for this variant is:

```text
[pcid, field_1, field_2, ..., field_n]
```

Slots are interpreted positionally:

- `pcid` identifies the protocol/spec/handler that defines the remaining
  outer fields.
- `field_1` through `field_n` are interpreted only by the handler named by
  `pcid`.

The first `pcid` defines the arity, field order, field types, canonical
signature coverage, and any distinction between payload bytes, proof bytes,
routing evidence, or application data. Generic tooling can identify the first
slot and preserve the full encoded array, but it cannot know which later field
is payload or proof without the `pcid` handler.

## Encoding

This variant encodes the envelope as a deterministic CBOR positional array.
The first slot is a CIDv1 byte string. Later slots are typed by the `pcid`
spec; this draft does not impose a universal byte-string-only rule on them.
The canonical bytes for hashing and signing are the deterministic CBOR bytes of
the exact array accepted by the `pcid` handler.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid`, it may store or forward the full
envelope bytes as opaque evidence, but it MUST NOT interpret any later field,
claim that a payload field exists, claim that a signature verifies, or expose
application data as understood. Unknown-variable-arity envelopes are
unsupported byte evidence until their `pcid` spec is available.

## Signature and Authorship Policy

This variant has no universal signature slot. A signature exists only when the
`pcid`-named schema defines one. That schema must define which fields are
covered, whether the first `pcid` is signed as associated data, and how signer
identity is represented. A scenario should penalize schemas that make payload
or type-substitution attacks possible.

## Layering-Test Behavior

This variant intentionally moves more responsibility into the first `pcid` than
the fixed-field variants do:

- Ordering disagreements are handled by the `pcid` schema because field count
  and field order are not shared across all envelopes.
- Unknown-pCID behavior is conservative because generic tooling cannot know
  which fields are safe to inspect.
- Signature verification is schema-local, so generic envelope tooling cannot
  verify messages without the `pcid` handler.
- Evolution can add new per-pCID schemas without changing the outer envelope
  spec, but each schema must solve arity, typing, signature coverage, and
  migration itself.

## Open Questions

- Does variable outer arity make generic routing, indexing, and audit too weak
  compared with fixed `[pcid, payload]` or `[pcid, payload, sig_pcid,
  sig_payload]` envelopes?
- Should each `pcid` schema publish machine-readable arity and field-type
  metadata, or is prose plus handler code acceptable?
- What maximum field count or encoded-size limits are required to avoid denial
  of service from adversarial arrays?
- Can signature coverage be specified reliably enough when there is no shared
  signature slot?

## Non-Canonical Status

This draft does not declare a winning envelope, does not define a central pCID
registry, and does not constrain sibling simulations. It exists to let the
runner compare variable outer arity against fixed-field and nested-payload
alternatives.
```

## Compact Fitness Evidence From This Run

### `SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields` x `chunk-feed-replication-sparse-advertisement`

- Result path: `results/SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields/chunk-feed-replication-sparse-advertisement/openai-gpt-5.4-xhigh/20260521-003110.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=1 evolution_safety=3 layer_boundary_clarity=4 failure_handling=3 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=19.00 normalized_0_100=47.50 confidence_0_1=0.84
- Rationale: Indirect fit only: this specimen defines an envelope, not a sparse chunk-feed advertisement protocol. Its conservative unknown-pCID rule and clear generic/schema-local boundary help safety, but the same choice makes leaves, frontiers, and summaries opaque to generic peers and later auditors.
- Strengths:
  - Conservative unknown-pCID handling prevents false interpretation and permits opaque retention or forwarding.
  - The generic vs. pCID-specific responsibility split is explicit, which improves layer-boundary clarity.
  - New per-pCID schemas can be added without changing the outer envelope, and deterministic CBOR keeps implementation plausible.
- Weaknesses:
  - It does not define how sparse chunk-replication advertisements represent leaves, roots, pointer objects, frontiers, or compact summaries.
  - Generic tooling cannot tell which outer fields are payload, proof, or advertisement summary without the pCID handler.
  - There is no universal signature slot, machine-readable arity metadata, or stated field-count/size limit.
- Risks:
  - Sparse peers may fail to interoperate if advertisements depend on unavailable or mismatched pCID schemas.
  - If old pCID schemas or handlers are lost, long-term evidence becomes hard to audit beyond opaque byte preservation.
  - Schema-local typing and signature rules, plus unresolved size limits, raise substitution and denial-of-service risk.
- Open questions:
  - Can a per-pCID sparse-advertisement schema expose roots, frontiers, or summaries with enough generic structure for routing and audit?
  - Should each pCID publish machine-readable arity, type, signature-coverage, and size-limit metadata?
  - What durable discovery path lets future peers recover the needed schema or handler for old envelopes?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields` x `conditional-release-geofencing-onward-restraint-chain`

- Result path: `results/SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields/conditional-release-geofencing-onward-restraint-chain/openai-gpt-5.4-xhigh/20260521-003110.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=1 evolution_safety=2 layer_boundary_clarity=2 failure_handling=3 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=11.00 normalized_0_100=40.00 confidence_0_1=0.86
- Rationale: Low fit for this scenario. This specimen only defines a variable-arity outer envelope; it could carry a conditional-release schema, but it does not define where the onward-restraint graph lives or how peers generically verify and audit recursive promises. Its main positive is conservative unknown-pCID handling, but meaning, signature coverage, and audit all become schema-local.
- Strengths:
  - Unknown-pCID handling is explicit and conservative: peers may store or forward opaque bytes without claiming understanding.
  - The draft avoids a central pCID registry and does not assume global state, which helps sparse-knowledge and no-central-authority goals.
  - Deterministic CBOR plus first-slot pCID dispatch is technically implementable and extensible for future specialized schemas.
- Weaknesses:
  - The simulation does not specify the recursive promise graph or whether it belongs at session, conditional-release, feed, or CAS-object level.
  - Generic tooling cannot know which outer fields are payload, proof, or signature, which weakens cross-peer audit of onward restraints.
  - There is no universal signature slot or coverage rule, so proof of Bob's and Carol's promises is schema-specific.
  - ... 1 more
- Risks:
  - Schema-local signature coverage could omit or ambiguously bind the restraint condition, enabling substitution or downgrade attacks.
  - If handlers are lost or diverge, restraint-chain evidence degrades into opaque bytes that later auditors cannot meaningfully interpret.
  - Unspecified maximum field count and encoded-size limits create denial-of-service and interoperability risk.
- Open questions:
  - Where should the recursive onward-restraint graph live in PromiseGrid: session state, conditional-release objects, feed entries, or separate CAS objects?
  - Should each pCID schema publish machine-readable arity, field-type, and signature-coverage metadata so generic tools can sanity-check envelopes?
  - What universal size and field-count limits are required for safe variable-arity envelopes?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields` x `bgp-routing`

- Result path: `results/SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields/bgp-routing/openai-gpt-5.4-xhigh/20260521-003110.json`
- Scores: scenario_fit=2 promisegrid_alignment=2 auditability=1 evolution_safety=2 layer_boundary_clarity=2 failure_handling=2 implementation_plausibility=3 risk_penalty=4
- Fitness: raw=10.00 normalized_0_100=38.00 confidence_0_1=0.82
- Rationale: Relevant as a stress probe, but weak for bgp-routing under sparse knowledge: route, proof, and signature meaning all live behind the first pCID, so peers without the exact handler can preserve bytes yet cannot safely inspect or verify contested routing evidence.
- Strengths:
  - Deterministic CBOR encoding and whole-envelope preservation support peer-local byte retention.
  - Unknown-pCID handling is conservative: store or forward opaque bytes, but do not pretend to parse or verify them.
  - New per-pCID schemas can be added without changing the outer array shape or requiring a central registry.
- Weaknesses:
  - The specimen does not itself define route promises, path evidence, stale-path handling, leak detection, or local routing policy.
  - Generic tooling cannot tell which fields are routing data, proofs, or signatures without the pCID-specific handler.
  - No universal signature slot or coverage rule makes cross-schema audit and generic verification weak.
  - ... 1 more
- Risks:
  - Hijack or leak evidence may become unverifiable across peers that do not share the same handler set.
  - Missing field-count and encoded-size limits leaves denial-of-service room in adversarial control-plane traffic.
  - Schema-local signature coverage can allow substitution or partial-coverage mistakes that generic auditors cannot detect early.
  - ... 1 more
- Open questions:
  - Should PromiseGrid require a minimal fixed routing-visible field set even if payload schemas vary?
  - Can per-pCID schemas publish machine-readable arity and field metadata durable enough for 100-year audit?
  - What universal size or field-count bounds are needed before this shape is safe for routing workloads?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-janov-grid-envelope-layer-pcid-nested-signed-payload` x `chunk-feed-replication-sparse-advertisement`

- Result path: `results/SIM-janov-grid-envelope-layer-pcid-nested-signed-payload/chunk-feed-replication-sparse-advertisement/openai-gpt-5.4-xhigh/20260521-003110.json`
- Scores: scenario_fit=1 promisegrid_alignment=3 auditability=3 evolution_safety=3 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=17.00 normalized_0_100=49.00 confidence_0_1=0.86
- Rationale: Low direct fit: the specimen does not assume global completeness and cleanly defines a nested signed envelope, but it never defines the sparse chunk-feed advertisement contract itself. Here it is mainly a transport envelope for a future replication protocol, and the unsigned outer conformance promise leaves relay/authorship claims weak.
- Strengths:
  - The spec makes the outer layer, nested payload, and signature coverage explicit.
  - Minimal outer shape lets a future sparse-advertisement schema evolve under pcid_b without changing the universal envelope.
  - Deterministic CBOR, signature over [pcid_b, payload_b], and opaque handling of unknown pCIDs support later verification under partial knowledge.
- Weaknesses:
  - It never defines the sparse-advertisement objects themselves: roots, leaves, pointer objects, frontiers, or summaries.
  - A valid inner signature does not show that the current relay actually holds the missing chunks or can satisfy replication.
  - Outer-layer conformance/authorship depends on transport or peer context because the envelope itself is unsigned.
- Risks:
  - A stale or malicious relay can replay or reframe a signed inner advertisement at the outer layer without cryptographic proof of who made the Layer-A promise.
  - Long-term audits may preserve bytes and verify the inner signature yet still fail to establish who promised chunk availability at the time of exchange.
  - Transport-identity dependence is brittle across migration, relay changes, and long-lived evidence review.
- Open questions:
  - Should payload_b carry a dedicated signed sparse-advertisement schema with root IDs, frontier summaries, actor identity, freshness, and availability scope?
  - Does the outer layer need its own signature or attestation slot to make relay promises auditable?
  - How should the design distinguish partial possession, full service, and mere forwarding of an advertisement from another actor?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-janov-grid-envelope-layer-pcid-nested-signed-payload` x `conditional-release-geofencing-onward-restraint-chain`

- Result path: `results/SIM-janov-grid-envelope-layer-pcid-nested-signed-payload/conditional-release-geofencing-onward-restraint-chain/openai-gpt-5.4-xhigh/20260521-003110.json`
- Scores: scenario_fit=2 promisegrid_alignment=2 auditability=2 evolution_safety=2 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=14.00 normalized_0_100=40.00 confidence_0_1=0.84
- Rationale: Useful mainly as a boundary probe: it strongly separates a minimal outer envelope from a signed inner payload, which suggests onward-restraint semantics would need to live in the inner payload or surrounding protocol, but it does not itself represent, enforce, or audit the recursive conditional-release promise chain.
- Strengths:
  - Clear outer/inner layer split with deterministic encoding and explicit signature coverage over [pcid_b, payload_b].
  - Explicit unknown-pCID behavior preserves opaque evidence without overclaiming interpretation.
  - Technically straightforward specimen for comparing whether restraint semantics belong in a universal envelope or a higher-layer payload.
- Weaknesses:
  - Does not define the recursive promise graph or say whether onward-restraint obligations live at group/session, transport/feed, or CAS-object level.
  - Outer conformance/authorship is unsigned, so who promised Layer-A validity can depend on transport context.
  - Lacks peer-local accounting rules for Alice, Bob, and Carol to prove acceptance, forwarding conditions, or violation of onward restraints.
- Risks:
  - A deployment could misread a valid inner signature as proof of Bob's forwarding promise when the outer release condition was never durably attested.
  - Transport-bound authorship weakens 100-year durability and migration across feeds, keys, or organizations.
  - Conditional-release enforcement may fragment across implementations because the envelope does not fix where the restraint graph is carried.
- Open questions:
  - Should the onward-restraint chain be a signed CAS object inside payload_b, or must some sender-bound proof also cover [pcid_a, payload_a]?
  - If Bob forwards to Carol, what exact locally recordable evidence shows Carol accepted the same restraint before release?
  - Can this variant support long-lived audit across unknown pcid_b versions and key rotation without relying on transport identity?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-janov-grid-envelope-layer-pcid-nested-signed-payload` x `bgp-routing`

- Result path: `results/SIM-janov-grid-envelope-layer-pcid-nested-signed-payload/bgp-routing/openai-gpt-5.4-xhigh/20260521-003110.json`
- Scores: scenario_fit=2 promisegrid_alignment=2 auditability=2 evolution_safety=3 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=20.00 normalized_0_100=50.00 confidence_0_1=0.81
- Rationale: Relevant as a routing-envelope subtest because it cleanly separates outer layer and signed inner payload, but it does not model BGP route propagation or policy, and the unsigned outer conformance/relay promise is a serious weakness under hijack/leak pressure.
- Strengths:
  - Clear two-layer envelope with deterministic CBOR and an inner signature over [pcid_b, payload_b].
  - Unknown-pCID policy preserves opaque evidence and avoids false claims of understanding.
  - Nested pcid_b supports payload evolution without changing the universal outer shape.
- Weaknesses:
  - Does not model route propagation, policy, withdrawals, or leak/hijack-specific promise accounting.
  - Unsigned outer [pcid_a, payload_a] leaves relay and conformance authorship tied to transport context rather than durable cryptographic evidence.
  - Local accept/retry/downgrade logic for stale or disputed routing evidence is unspecified.
- Risks:
  - A valid signed inner payload can be replayed in a misleading outer context, obscuring who endorsed propagation.
  - Carol may lack durable evidence to distinguish origin, relay, and leak during partitions or stale-state disputes.
- Open questions:
  - For routing, should outer propagation or endorsement be signed separately from the inner payload?
  - Can transport identity plus local peer-adoption records satisfy 100-year auditability for contested routing evidence?
  - Should routing-specific path or hop attestations live inside payload_b instead of relying on envelope context?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

## Required JSON Shape

{"child_id":"SIM-maraz-child-descriptive-design-slug","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
