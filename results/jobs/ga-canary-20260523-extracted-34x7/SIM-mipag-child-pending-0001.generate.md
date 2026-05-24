# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Choose a descriptive `child_id` that starts with `SIM-mipag-child-` and ends with a kebab-case design slug. Do not use generic `ga-child`, `pending`, or ordinal-only names.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260523-extracted-34x7`
- Planned child ID prefix: `SIM-mipag-child-`
- Temporary child ID: `SIM-mipag-child-pending-0001`
- Temporary child path: `proposals/ga-canary-20260523-extracted-34x7/simulations/SIM-mipag-child-pending-0001/`
- Operation: `breed`
- Parent IDs: `SIM-sutap-spec-sections-split-template, SIM-gibut-conditional-release-group-session-local`

## Scenario Sample

- `cas-object-type-binding-bakeoff-unknown-typed-object` at `scenarios/cas-object-type-binding-bakeoff-unknown-typed-object/cas-object-type-binding-bakeoff-unknown-typed-object.md`
- `conditional-release-geofencing-replay-outside-conditions` at `scenarios/conditional-release-geofencing-replay-outside-conditions/conditional-release-geofencing-replay-outside-conditions.md`
- `transport-family-bakeoff-per-hop-authorization-failure` at `scenarios/transport-family-bakeoff-per-hop-authorization-failure/transport-family-bakeoff-per-hop-authorization-failure.md`
- `peer-adoption-metadata-spec-answer-vocabulary-drift` at `scenarios/peer-adoption-metadata-spec-answer-vocabulary-drift/peer-adoption-metadata-spec-answer-vocabulary-drift.md`
- `spec-requirement-sections-100-year-review` at `scenarios/spec-requirement-sections-100-year-review/spec-requirement-sections-100-year-review.md`
- `group-session-freeze-promise-post-freeze-mutation-request` at `scenarios/group-session-freeze-promise-post-freeze-mutation-request/group-session-freeze-promise-post-freeze-mutation-request.md`
- `udp-feed-v0-conformance-session-layer-composition` at `scenarios/udp-feed-v0-conformance-session-layer-composition/udp-feed-v0-conformance-session-layer-composition.md`

## Scenario Pressure

### `scenarios/cas-object-type-binding-bakeoff-unknown-typed-object/cas-object-type-binding-bakeoff-unknown-typed-object.md`

```markdown
# Unknown typed object

## Scenario ID

cas-object-type-binding-bakeoff-unknown-typed-object

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kohad-cas-object-type-binding-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-kohad-cas-object-type-binding-bakeoff/`
- Source row/title: Unknown typed object
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kohad-cas-object-type-binding-bakeoff/`.

## Setup

Dave receives a CID whose codec he does not implement.

## Stimulus

Run the candidate simulation against this source test: Whether the peer can store, advertise, and forward the object opaquely while avoiding unsafe parsing.

## Expected Pressure

Type binding must define unknown-type behavior for long-lived mixed-version networks.
```

### `scenarios/conditional-release-geofencing-replay-outside-conditions/conditional-release-geofencing-replay-outside-conditions.md`

```markdown
# Replay outside conditions

## Scenario ID

conditional-release-geofencing-replay-outside-conditions

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zarud-conditional-release-geofencing/SCENARIOS.md`
- Source simulation: `SIM-zarud-conditional-release-geofencing/`
- Source row/title: Replay outside conditions
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zarud-conditional-release-geofencing/`.

## Setup

Mallory replays a valid old content reference to Dave outside the allowed audience or geography.

## Stimulus

Run the candidate simulation against this source test: Whether receivers, feeds, or group/session state detect stale or unauthorized reuse.

## Expected Pressure

Replay handling determines whether conditions must bind to recipients, epochs, locations, or session context.
```

### `scenarios/transport-family-bakeoff-per-hop-authorization-failure/transport-family-bakeoff-per-hop-authorization-failure.md`

```markdown
# Per-hop authorization failure

## Scenario ID

transport-family-bakeoff-per-hop-authorization-failure

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-narok-transport-family-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-narok-transport-family-bakeoff/`
- Source row/title: Per-hop authorization failure
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-narok-transport-family-bakeoff/`.

## Setup

Bob receives a ring message but is not authorized to forward it to Carol.

## Stimulus

Run the candidate simulation against this source test: Whether authorization failure breaks the ring, skips a hop, records refusal, or reconfigures membership.

## Expected Pressure

Ring semantics need a failure model before they can be compared with gossip.
```

### `scenarios/peer-adoption-metadata-spec-answer-vocabulary-drift/peer-adoption-metadata-spec-answer-vocabulary-drift.md`

```markdown
# Spec answer vocabulary drift

## Scenario ID

peer-adoption-metadata-spec-answer-vocabulary-drift

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-dihiz-peer-adoption-metadata/SCENARIOS.md`
- Source simulation: `SIM-dihiz-peer-adoption-metadata/`
- Source row/title: Spec answer vocabulary drift
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-dihiz-peer-adoption-metadata/`.

## Setup

A later frozen spec supersedes pCID X and renames or removes Q9.

## Stimulus

Run the candidate simulation against this source test: Whether answer keys are spec-local, profile-local, globally named, or mapped through migration records.

## Expected Pressure

Adoption metadata must survive spec evolution without making old claims ambiguous.
```

### `scenarios/spec-requirement-sections-100-year-review/spec-requirement-sections-100-year-review.md`

```markdown
# 100-year review

## Scenario ID

spec-requirement-sections-100-year-review

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-ranib-spec-requirement-sections/SCENARIOS.md`
- Source simulation: `SIM-ranib-spec-requirement-sections/`
- Source row/title: 100-year review
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-ranib-spec-requirement-sections/`.

## Setup

Carol evaluates whether an old spec is still intelligible after tooling, hosts, and social assumptions changed.

## Stimulus

Run the candidate simulation against this source test: Whether a required long-horizon section preserves the assumptions needed for future readers.

## Expected Pressure

TE-dajot-style pressure may need a concrete spec-section home.
```

### `scenarios/group-session-freeze-promise-post-freeze-mutation-request/group-session-freeze-promise-post-freeze-mutation-request.md`

```markdown
# Post-freeze mutation request

## Scenario ID

group-session-freeze-promise-post-freeze-mutation-request

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-bohof-group-session-freeze-promise/SCENARIOS.md`
- Source simulation: `SIM-bohof-group-session-freeze-promise/`
- Source row/title: Post-freeze mutation request
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-bohof-group-session-freeze-promise/`.

## Setup

Carol proposes a breaking group-session change after freeze evidence exists.

## Stimulus

Run the candidate simulation against this source test: Whether the change becomes a new specimen, a superseding pCID, or an amendment to the old lineage.

## Expected Pressure

Freeze semantics must preserve independent evolution without mutating history.
```

### `scenarios/udp-feed-v0-conformance-session-layer-composition/udp-feed-v0-conformance-session-layer-composition.md`

```markdown
# Session-layer composition

## Scenario ID

udp-feed-v0-conformance-session-layer-composition

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kuful-udp-feed-v0-conformance/SCENARIOS.md`
- Source simulation: `SIM-kuful-udp-feed-v0-conformance/`
- Source row/title: Session-layer composition
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kuful-udp-feed-v0-conformance/`.

## Setup

A minimal group/session message rides above UDP-feed v0.

## Stimulus

Run the candidate simulation against this source test: Whether UDP-feed's API is sufficient for the next layer without leaking binding details.

## Expected Pressure

If composition is required, TODO-jodon's done criteria must include more than UDP round trip.
```

## Parent Simulation Documents

### `simulations/SIM-sutap-spec-sections-split-template/README.md`

```markdown
# SIM-sutap-spec-sections-split-template

This simulation turns the split-template spec-sections alternative from
`SIM-ranib-spec-requirement-sections` into a concrete candidate specimen. It
tests whether specs should carry concise normative hooks while longer
mental-model prose lives in companion documents. Source: `DI-fibuv`.

## Design Under Test

Each protocol spec promises compact required hooks for vocabulary, durability,
and implementation expectations. Companion docs may expand those hooks without
changing the pCID-owned normative surface.

## Boundaries

This simulation does not lock the final template. It tests whether a split
structure preserves pCID precision while allowing readable explanatory material
to evolve separately.
```

### `simulations/SIM-sutap-spec-sections-split-template/QUESTION.md`

```markdown
# Question

Can a split spec template keep pCID-owned requirements compact while letting
longer PromiseGrid mental-model prose evolve without changing the protocol
promise?

Source: `DI-fibuv`.
```

### `simulations/SIM-gibut-conditional-release-group-session-local/README.md`

```markdown
# SIM-gibut-conditional-release-group-session-local

This simulation turns the group-session-local conditional-release alternative
from `SIM-zarud-conditional-release-geofencing` into a concrete candidate
specimen. It tests whether onward-restraint, geofencing, and release conditions
belong inside the group/session semantics that deliver the content. Source:
`DI-fibuv`.

## Design Under Test

Group-session messages carry or reference release conditions, and each session
participant promises to preserve those conditions when dispatching content to
later recipients.

## Boundaries

This simulation does not modify the current group-session specimen. It tests
whether session-local ownership gives clear human semantics or overloads the
session layer with policy that should be separate.
```

### `simulations/SIM-gibut-conditional-release-group-session-local/QUESTION.md`

```markdown
# Question

Can group-session semantics own conditional release and onward-restraint
promises without making lower-layer routing, storage, or audit behavior
ambiguous?

Source: `DI-fibuv`.
```

## Compact Fitness Evidence From This Run

### `SIM-gibut-conditional-release-group-session-local` x `cas-object-type-binding-bakeoff-unknown-typed-object`

- Result path: `results/SIM-gibut-conditional-release-group-session-local/cas-object-type-binding-bakeoff-unknown-typed-object/openai-gpt-5.4-medium/20260523-230210.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=2 evolution_safety=1 layer_boundary_clarity=1 failure_handling=1 implementation_plausibility=2 promise_vocabulary=2 simplicity_durability=1 risk_penalty=3
- Fitness: raw=15.00 normalized_0_100=30.00 confidence_0_1=0.84
- Rationale: This specimen is largely off-target for the unknown-typed-object bakeoff. It discusses conditional release inside group-session semantics, but does not clearly specify opaque handling of unknown CAS object types, safe non-parsing behavior, or mixed-version forwarding rules. It does use promise-oriented framing, yet the scenario pressures on type binding, storage/advertisement/forwarding behavior, and layer separation remain mostly unanswered.
- Strengths:
  - Uses explicit promise language about participants preserving onward conditions.
  - Acknowledges layer-boundary ambiguity as a core design question.
  - Could support local accounting of who promised to preserve release constraints.
- Weaknesses:
  - Does not define unknown-codec object behavior directly.
  - Does not clearly show opaque store/advertise/forward semantics.
  - Provides weak guidance for mixed-version evolution and safe fallback.
  - ... 1 more
- Risks:
  - Session semantics may be overloaded with policy that should live in separate payload or binding artifacts.
  - Peers may need to inspect unknown typed content or condition encodings to enforce release rules, creating unsafe parsing pressure.
  - Audit records may become ambiguous about whether failure came from transport, storage, policy, or unknown-type handling.
- Open questions:
  - Can session-local release conditions be attached to unknown-codec CAS objects without requiring semantic parsing by transit peers?
  - What exact promise lets Dave store and forward an unknown typed object opaquely while preserving any attached release conditions?
  - If a future codec changes condition encoding, which layer detects incompatibility and how is safe fallback recorded locally?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-gibut-conditional-release-group-session-local` x `conditional-release-geofencing-replay-outside-conditions`

- Result path: `results/SIM-gibut-conditional-release-group-session-local/conditional-release-geofencing-replay-outside-conditions/openai-gpt-5.4-medium/20260523-230210.json`
- Scores: scenario_fit=3 promisegrid_alignment=2 auditability=2 evolution_safety=2 layer_boundary_clarity=1 failure_handling=2 implementation_plausibility=3 promise_vocabulary=2 simplicity_durability=2 risk_penalty=4
- Fitness: raw=20.00 normalized_0_100=40.00 confidence_0_1=0.78
- Rationale: This specimen engages the scenario directly by placing release conditions inside group-session semantics, so replay outside the allowed audience or geography can be discussed in terms of session context and participant promises. But it scores down because replay control depends heavily on mutable session-local state, making off-session audit, durable evidence, and layer boundaries unclear. It is plausible to implement, yet weak on 100-year durable, payload-level promise artifacts and on clear handling of stale replays outside current session context.
- Strengths:
  - Directly targets replay pressure by binding release conditions to group/session delivery semantics.
  - Uses participant promises about onward restraint rather than only transport mechanics.
  - Implementable with session metadata, membership checks, and release-condition checks at dispatch time.
- Weaknesses:
  - Overloads the session layer with policy, weakening layer boundary clarity.
  - Does not clearly define durable evidence for stale or unauthorized replay detection outside active session context.
  - Promise vocabulary is session/policy-centric rather than explicit payload-promise-centric.
  - ... 1 more
- Risks:
  - Session-local conditions may create false confidence if old references are replayed where the original session context is unavailable or disputed.
  - Geofencing and audience checks may drift into implicit central authorities for membership, location, or policy validation.
  - Audit trails may be too context-dependent to support later human or LLM review of what was actually promised at release time.
- Open questions:
  - What durable artifact lets a receiver verify that this replayed content reference is still within allowed audience, geography, and epoch?
  - Can a peer audit an unauthorized replay later without access to full mutable group-session state?
  - How are release-condition changes, membership churn, and location-evidence formats evolved without breaking old evidence?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-gibut-conditional-release-group-session-local` x `transport-family-bakeoff-per-hop-authorization-failure`

- Result path: `results/SIM-gibut-conditional-release-group-session-local/transport-family-bakeoff-per-hop-authorization-failure/openai-gpt-5.4-medium/20260523-230210.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=2 evolution_safety=2 layer_boundary_clarity=1 failure_handling=1 implementation_plausibility=3 promise_vocabulary=3 simplicity_durability=2 risk_penalty=4
- Fitness: raw=20.00 normalized_0_100=40.00 confidence_0_1=0.76
- Rationale: The simulation is relevant because session-local release conditions can express Bob refusing onward delivery to Carol, but it does not define the ring-style per-hop failure model the scenario asks for. Its core uncertainty is exactly whether conditional release inside group/session semantics makes lower-layer behavior ambiguous, so it partially fits while leaving key comparison pressure unresolved.
- Strengths:
  - Uses explicit participant promises rather than pure central authorization framing.
  - Naturally represents onward-restraint as a local forwarding decision by Bob.
  - Implementation is plausible as an extension of group/session message semantics.
- Weaknesses:
  - Does not clearly specify what happens to the ring after a per-hop authorization failure.
  - Layer boundary between session policy and transport behavior is blurry.
  - Local audit evidence for refusal is not defined.
  - ... 1 more
- Risks:
  - Policy and transport concerns may be conflated, making forwarding failures hard to interpret.
  - Different participants may implement refusal, skip, or reconfiguration differently, fragmenting semantics.
  - Audit trails may be too weak to distinguish authorization denial from ordinary delivery failure.
- Open questions:
  - If Bob cannot forward to Carol, does the session stop, skip Carol, or require explicit membership change?
  - What local artifact records Bob's refusal so Carol or an auditor can distinguish denial from transport loss?
  - Which promises belong to session semantics versus lower transport/routing layers when authorization and forwarding diverge?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-gibut-conditional-release-group-session-local` x `peer-adoption-metadata-spec-answer-vocabulary-drift`

- Result path: `results/SIM-gibut-conditional-release-group-session-local/peer-adoption-metadata-spec-answer-vocabulary-drift/openai-gpt-5.4-medium/20260523-230210.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=2 evolution_safety=1 layer_boundary_clarity=1 failure_handling=1 implementation_plausibility=3 promise_vocabulary=3 simplicity_durability=2 risk_penalty=4
- Fitness: raw=17.00 normalized_0_100=34.00 confidence_0_1=0.85
- Rationale: The specimen is mostly about session-local conditional release, not adoption-metadata answer-key evolution. It has some PromiseGrid flavor because participants make local onward-restraint promises, but it does not explain how spec-local question identifiers, superseded pCIDs, or migration records remain unambiguous after vocabulary drift. The design also explicitly risks blurring session semantics with audit/policy concerns, which hurts this scenario.
- Strengths:
  - Uses explicit participant promises rather than pure central-registry framing.
  - Peer-local handling is at least plausible for enforcing conditional release inside a session.
- Weaknesses:
  - Does not directly address adoption metadata or spec-answer key naming.
  - No clear mechanism for spec-local vs profile-local keys or migration records.
  - Layer boundary is blurry: session semantics are asked to carry policy, release, and audit meaning.
  - ... 1 more
- Risks:
  - Old and new spec-answer terms can become ambiguous because no spec-evolution mapping artifact is defined.
  - Embedding policy into session semantics may make later audit and migration behavior inconsistent across peers.
  - Participants may promise onward restraint locally, but third parties may lack durable evidence tying an answer to the exact superseded spec vocabulary.
- Open questions:
  - Could release conditions carry explicit spec references or migration pointers without overloading session semantics further?
  - What locally auditable artifact would let Carol distinguish an old Q9 answer from a renamed successor question after spec change?
  - Should adoption/spec-answer metadata live in a separate payload promise layer instead of group-session semantics?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-gibut-conditional-release-group-session-local` x `spec-requirement-sections-100-year-review`

- Result path: `results/SIM-gibut-conditional-release-group-session-local/spec-requirement-sections-100-year-review/openai-gpt-5.4-medium/20260523-230210.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=1 evolution_safety=1 layer_boundary_clarity=1 failure_handling=1 implementation_plausibility=3 promise_vocabulary=2 simplicity_durability=1 risk_penalty=4
- Fitness: raw=14.00 normalized_0_100=28.00 confidence_0_1=0.83
- Rationale: The simulation names a relevant design question, but it is a weak fit for a 100-year spec-review scenario because it does not provide a concrete durable spec-section shape, artifact set, or future-reader-oriented explanation. Its main value here is exposing that embedding conditional release into group-session semantics risks long-term layer confusion and weak auditability.
- Strengths:
  - Explicitly tests whether conditional release belongs in session semantics.
  - Uses promise language at a high level by making participants responsible for preserving conditions.
  - Surfaces boundary tension between session meaning and lower-layer behavior.
- Weaknesses:
  - No concrete 100-year spec section or review-oriented structure is provided.
  - Audit trail, local records, and observable artifacts are underspecified.
  - Failure, adversarial pressure, and evolution mechanics are mostly absent.
  - ... 1 more
- Risks:
  - Session semantics may absorb policy that should live in a separate layer, making future interpretation brittle.
  - Participants are asked to preserve conditions, but the simulation does not define durable local evidence for whether they did so.
  - Geofencing and onward-restraint assumptions may age poorly as networks, organizations, and enforcement models change.
- Open questions:
  - What exact long-horizon spec section would preserve the assumptions of session-local conditional release for future readers?
  - What durable local audit artifacts would show which release conditions were promised and later propagated?
  - Can onward-restraint and geofencing be specified without making storage, routing, and replay semantics layer-ambiguous over long time horizons?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-gibut-conditional-release-group-session-local` x `group-session-freeze-promise-post-freeze-mutation-request`

- Result path: `results/SIM-gibut-conditional-release-group-session-local/group-session-freeze-promise-post-freeze-mutation-request/openai-gpt-5.4-medium/20260523-230210.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=2 evolution_safety=1 layer_boundary_clarity=1 failure_handling=1 implementation_plausibility=3 promise_vocabulary=2 simplicity_durability=2 risk_penalty=4
- Fitness: raw=16.00 normalized_0_100=32.00 confidence_0_1=0.81
- Rationale: This simulation is adjacent to the scenario but does not directly specify freeze lineage behavior. It places conditional-release obligations inside group-session semantics, which gives a concrete promise surface, yet the scenario asks how breaking post-freeze change is represented without mutating history. The design's own stated boundary/question admits likely ambiguity in routing, storage, and audit layers, so evolution safety and layer clarity are weak here.
- Strengths:
  - Uses explicit participant promises rather than central registry assumptions.
  - Peer-local onward-restraint semantics are at least implementable within session handling.
  - Directly explores whether ownership of release conditions at the session layer is coherent.
- Weaknesses:
  - Does not define what a post-freeze breaking mutation becomes in lineage terms.
  - Session layer is overloaded with policy, weakening boundary clarity.
  - Audit trail for frozen versus superseding terms is underspecified.
  - ... 1 more
- Risks:
  - Breaking changes may be expressed as session-level amendment instead of a clearly new frozen specimen.
  - Embedding release policy into session semantics can blur immutable content lineage versus mutable dispatch obligations.
  - Later auditors may be unable to distinguish historical frozen promises from post-freeze policy edits.
- Open questions:
  - If Carol proposes a breaking post-freeze change, does this design force a new frozen specimen or allow ambiguous amendment of the old session lineage?
  - What local artifact would let Alice or Carol audit that onward-restraint terms were preserved across pre-freeze and post-freeze variants?
  - Are release conditions session semantics, payload promises, or a separate policy layer referenced by frozen content IDs?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-gibut-conditional-release-group-session-local` x `udp-feed-v0-conformance-session-layer-composition`

- Result path: `results/SIM-gibut-conditional-release-group-session-local/udp-feed-v0-conformance-session-layer-composition/openai-gpt-5.4-medium/20260523-230210.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=2 evolution_safety=2 layer_boundary_clarity=1 failure_handling=1 implementation_plausibility=3 promise_vocabulary=3 simplicity_durability=2 risk_penalty=4
- Fitness: raw=21.00 normalized_0_100=42.00 confidence_0_1=0.74
- Rationale: Partial fit. The simulation directly pressures session-layer composition by placing conditional-release semantics inside group/session messages and asking whether routing, storage, and audit behavior become ambiguous. That matches the scenario's boundary question, but it does not specify a concrete UDP-feed v0 composition contract, so transport-conformance evidence remains thin. Promise-first intent is present, yet auditability, failure handling, and layer separation are weak because onward-restraint depends on later participant behavior and is hard to verify locally.
- Strengths:
  - Explicitly tests whether session-local semantics overload the next layer.
  - Uses participant promises rather than central registry assumptions.
  - Concrete enough to imagine implementation by attaching conditions to session messages or references.
- Weaknesses:
  - Does not define the UDP-feed v0 interface boundary the session layer relies on.
  - Local audit records for compliance and breach are underspecified.
  - Failure cases like refusal, stale policy, partial delivery, and malicious relay are not worked through.
- Risks:
  - Session semantics may absorb policy concerns that belong in a separate artifact or layer.
  - Onward-restraint promises may be unverifiable after forwarding, weakening local auditability.
  - Transport/session coupling may grow implicitly if release conditions depend on delivery or storage behavior.
- Open questions:
  - What exact UDP-feed v0 payload/API surface does the session layer consume without leaking transport binding details?
  - How does a peer record and later audit another participant's onward-restraint or release-condition breach using only local evidence?
  - Can release conditions evolve independently of session framing, or does this design hard-couple policy to the session layer?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-sutap-spec-sections-split-template` x `cas-object-type-binding-bakeoff-unknown-typed-object`

- Result path: `results/SIM-sutap-spec-sections-split-template/cas-object-type-binding-bakeoff-unknown-typed-object/openai-gpt-5.4-medium/20260523-230210.json`
- Scores: scenario_fit=1 promisegrid_alignment=3 auditability=4 evolution_safety=4 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=4 promise_vocabulary=2 simplicity_durability=4 risk_penalty=3
- Fitness: raw=29.00 normalized_0_100=58.00 confidence_0_1=0.79
- Rationale: This simulation is mainly about spec structure, not CAS type-binding behavior. It helps by separating stable normative hooks from evolving explanation, which is useful for long-lived unknown-type handling, but it does not itself specify the required opaque-forwarding and no-unsafe-parse behavior.
- Strengths:
  - Strong separation between normative protocol surface and explanatory material.
  - Good fit for long-term evolution and small durable spec artifacts.
  - Improves auditability of what the pCID actually owns.
- Weaknesses:
  - Does not directly solve the unknown typed object behavior in the scenario.
  - Failure handling for unsupported codecs is not concretely specified.
  - Promise-first payload wording is only partial, not fully explicit.
- Risks:
  - Normative hooks may be too thin to fully pin down unknown-type safety behavior.
  - Companion-doc evolution could reintroduce interpretation drift across implementations.
  - Readers may mistake explanatory guidance for the real protocol promise boundary.
- Open questions:
  - Can the compact normative hook explicitly require opaque store/advertise/forward behavior for unknown codecs without relying on companion prose?
  - How does the split template prevent drift between pCID-owned normative text and evolving explanatory docs for mixed-version peers?
  - What minimal normative language defines unsafe parsing avoidance and observable peer-local outcomes?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-sutap-spec-sections-split-template` x `conditional-release-geofencing-replay-outside-conditions`

- Result path: `results/SIM-sutap-spec-sections-split-template/conditional-release-geofencing-replay-outside-conditions/openai-gpt-5.4-medium/20260523-230210.json`
- Scores: scenario_fit=1 promisegrid_alignment=4 auditability=3 evolution_safety=4 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=3 promise_vocabulary=2 simplicity_durability=4 risk_penalty=2
- Fitness: raw=29.00 normalized_0_100=58.00 confidence_0_1=0.78
- Rationale: Good fit for PromiseGrid’s preference for compact durable normative artifacts and clearer separation between stable protocol promises and evolving explanation, but this simulation is a spec-template exercise rather than a concrete conditional-release mechanism. It offers little direct replay/geofencing handling beyond the possibility of specifying such hooks more cleanly.
- Strengths:
  - Compact pCID-owned normative surface supports long-term durability.
  - Clear split between normative hooks and explanatory material improves layer boundaries.
  - Evolution of mental-model prose can occur without changing the protocol promise.
- Weaknesses:
  - Does not directly solve replay detection or conditional-release enforcement.
  - Promise-first vocabulary is only partial; it is framed mainly as a spec template.
  - Failure handling and local audit evidence for stale or unauthorized reuse are not specified.
- Risks:
  - Security-critical replay or geofencing conditions may be under-specified if the compact normative surface is too terse.
  - Readers may rely on companion prose for behavior that should be normatively bound to the pCID.
  - Scenario pressure may be missed entirely because the template does not itself define detection, rejection, or audit records for unauthorized reuse.
- Open questions:
  - Can the compact normative section force recipient/epoch/geography binding strongly enough to prevent replay outside conditions?
  - How will companion prose be prevented from becoming security-critical guidance that diverges from the pCID-owned normative hooks?
  - What minimum replay-detection artifacts must a split template require peers to record locally?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-sutap-spec-sections-split-template` x `transport-family-bakeoff-per-hop-authorization-failure`

- Result path: `results/SIM-sutap-spec-sections-split-template/transport-family-bakeoff-per-hop-authorization-failure/openai-gpt-5.4-medium/20260523-230210.json`
- Scores: scenario_fit=1 promisegrid_alignment=3 auditability=4 evolution_safety=4 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=4 promise_vocabulary=2 simplicity_durability=4 risk_penalty=3
- Fitness: raw=29.00 normalized_0_100=58.00 confidence_0_1=0.85
- Rationale: This simulation is mainly about spec-document structure, not transport behavior, so it fits the scenario weakly. Its main value here is that a compact normative surface plus separate explanatory prose could improve auditability and evolution safety if failure semantics are explicitly required. But the candidate does not itself define per-hop authorization handling, so it leaves the scenario's core pressure unresolved.
- Strengths:
  - Clear separation between normative protocol hooks and evolving explanatory prose.
  - Good durability/evolution story for long-lived specs.
  - Compact artifacts are plausible and auditable if the required hooks are explicit.
- Weaknesses:
  - Does not directly model ring forwarding or per-hop authorization failure outcomes.
  - Promise-first wording is only partial and remains template-level rather than actor/payload-level.
  - Failure handling for the scenario is mostly deferred to whatever future transport spec uses the template.
- Risks:
  - A split template may omit or under-specify critical authorization-failure behavior in the normative section.
  - Companion prose could drift from the pCID-owned hooks and confuse implementers about actual failure obligations.
  - Transport comparisons become hard if different specs use the template but encode failure semantics with inconsistent precision.
- Open questions:
  - Can the split template require an explicit failure-semantics section for per-hop authorization cases without bloating the pCID-owned surface?
  - How would Alice or Carol audit Bob's refusal locally if the normative hooks stay very compact?
  - Does the template force transport specs to distinguish break-ring, skip-hop, refusal-record, and membership-reconfiguration outcomes?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-sutap-spec-sections-split-template` x `peer-adoption-metadata-spec-answer-vocabulary-drift`

- Result path: `results/SIM-sutap-spec-sections-split-template/peer-adoption-metadata-spec-answer-vocabulary-drift/openai-gpt-5.4-medium/20260523-230210.json`
- Scores: scenario_fit=3 promisegrid_alignment=4 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 promise_vocabulary=3 simplicity_durability=4 risk_penalty=2
- Fitness: raw=36.00 normalized_0_100=72.00 confidence_0_1=0.78
- Rationale: The split-template sim helps on this scenario by separating compact pCID-owned normative hooks from changeable explanatory prose, which reduces vocabulary drift risk and sharpens what belongs to the protocol surface. But it is only an indirect fit: it does not itself define adoption-metadata answer-key semantics, migration records, or a concrete rule for renamed/removed fields such as Q9.
- Strengths:
  - Strong separation between normative surface and explanatory material.
  - Good fit for spec evolution because pCID-owned requirements stay compact and stable.
  - Clearer audit target for later readers than a blended spec/prose document.
- Weaknesses:
  - Does not directly specify whether answer keys are spec-local, profile-local, global, or mapped.
  - Limited failure-mode treatment for stale or conflicting old/new vocabulary.
  - Promise-first wording is present but not fully grounded in concrete payload-level promises.
- Risks:
  - Companion prose could reintroduce ambiguous terminology if the normative hook for answer-key identity is too thin.
  - Without explicit migration mechanics, old and new adoption claims may still be hard to compare across superseding specs.
  - Reviewers may overread the template as solving metadata evolution when it mainly structures where such rules would live.
- Open questions:
  - Does the split template explicitly require answer keys to be spec-local to the superseded/new pCID rather than globally named?
  - How are renamed or removed questions like Q9 carried forward: alias table, migration record, or explicit non-equivalence?
  - Which artifacts remain normative for auditors when companion prose evolves but adoption metadata must stay unambiguous?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-sutap-spec-sections-split-template` x `spec-requirement-sections-100-year-review`

- Result path: `results/SIM-sutap-spec-sections-split-template/spec-requirement-sections-100-year-review/openai-gpt-5.4-medium/20260523-230210.json`
- Scores: scenario_fit=4 promisegrid_alignment=4 auditability=4 evolution_safety=5 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 promise_vocabulary=3 simplicity_durability=4 risk_penalty=2
- Fitness: raw=38.00 normalized_0_100=76.00 confidence_0_1=0.84
- Rationale: The split-template design fits the scenario well because it explicitly separates compact pCID-owned normative hooks from longer explanatory prose, which improves boundary clarity and evolution safety. It supports long-horizon review by keeping the required protocol surface small and stable while allowing mental-model documents to evolve. The main weakness is that the scenario specifically pressures preservation of future-reader assumptions, and this simulation does not yet define a concrete durable home or audit rule for those assumptions when companion prose changes or disappears.
- Strengths:
  - Clear normative vs companion-doc boundary.
  - Strong evolution story: explanatory prose can change without changing the protocol promise.
  - Compact required hooks support durable, auditable long-term review artifacts.
- Weaknesses:
  - No explicit failure-mode treatment for stale or absent companion docs.
  - Promise-first wording is present but still secondary to spec-template framing.
  - The required 100-year review content is not yet pinned to a concrete mandatory section shape.
- Risks:
  - Future readers may lose critical assumptions if they live mainly in non-normative companion documents.
  - Normative and explanatory documents may drift, reducing long-term auditability.
  - The design does not yet show explicit handling for missing, stale, or contradictory companion prose.
- Open questions:
  - Should the 100-year assumptions section be normative inside the pCID-owned spec or only referenced from companion prose?
  - How is drift between compact normative hooks and evolving explanatory companions detected and audited over time?
  - What minimum long-horizon review checklist must remain in the normative surface to survive host, tooling, and social change?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-sutap-spec-sections-split-template` x `group-session-freeze-promise-post-freeze-mutation-request`

- Result path: `results/SIM-sutap-spec-sections-split-template/group-session-freeze-promise-post-freeze-mutation-request/openai-gpt-5.4-medium/20260523-230210.json`
- Scores: scenario_fit=3 promisegrid_alignment=4 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 promise_vocabulary=3 simplicity_durability=4 risk_penalty=2
- Fitness: raw=36.00 normalized_0_100=72.00 confidence_0_1=0.79
- Rationale: The split-template candidate fits this scenario indirectly: it strongly separates frozen pCID-owned normative hooks from mutable explanatory prose, which helps preserve history after freeze. That gives good evolution safety and boundary clarity for post-freeze requests. But the specimen does not explicitly define the decision rule for a breaking mutation request (new specimen, superseding pCID, or amendment lineage), so scenario coverage is only partial.
- Strengths:
  - Clear normative versus explanatory boundary.
  - Compact frozen surface supports durable audit after long time horizons.
  - Explicit goal of allowing explanatory evolution without mutating the protocol promise.
- Weaknesses:
  - Does not directly answer the scenario's new-specimen versus superseding-pCID versus amendment choice.
  - Failure/refusal and adversarial handling are not specified.
  - Uses some PromiseGrid language, but not strongly payload-promise phrasing.
- Risks:
  - Companion documents could be mistaken for normative updates after freeze.
  - Operators may blur frozen protocol promises and evolving mental-model prose.
  - Without explicit lineage rules, different peers may classify the same breaking change differently.
- Open questions:
  - Does any post-freeze normative hook change always require a new pCID, or can lineage amendments exist without history mutation?
  - How are readers prevented from treating evolved companion prose as retroactive protocol change?
  - What minimal audit record links frozen normative sections to later explanatory updates?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-sutap-spec-sections-split-template` x `udp-feed-v0-conformance-session-layer-composition`

- Result path: `results/SIM-sutap-spec-sections-split-template/udp-feed-v0-conformance-session-layer-composition/openai-gpt-5.4-medium/20260523-230210.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=4 evolution_safety=5 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=3 promise_vocabulary=3 simplicity_durability=4 risk_penalty=2
- Fitness: raw=33.00 normalized_0_100=66.00 confidence_0_1=0.80
- Rationale: Strong as a spec-structure candidate for compact, evolvable, auditable protocol requirements, but only weakly fitted to this scenario because it does not directly specify UDP-feed/session-layer composition behavior.
- Strengths:
  - Keeps normative surface compact and pCID-scoped.
  - Supports long-term evolution of explanatory material without changing the core promise.
  - Improves audit focus by separating required hooks from mental-model prose.
- Weaknesses:
  - Does not directly define session-layer composition behavior for UDP-feed v0.
  - Provides little explicit failure-handling guidance.
  - Promise wording is somewhat promise-aware but not strongly payload-promise-first.
- Risks:
  - Companion prose could carry important composition guidance that is not normatively anchored.
  - A split template may appear to solve layering while leaving binding-detail leakage underspecified.
  - Scenario-specific conformance could be unverifiable if the compact hooks are too abstract.
- Open questions:
  - Should the compact normative section explicitly state transport/session-layer non-leakage requirements?
  - What minimal conformance hooks would let a later session layer prove composition without relying on companion prose?
  - How are failure and boundary cases recorded when companion docs evolve faster than the pCID-owned normative text?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

## Required JSON Shape

{"child_id":"SIM-mipag-child-descriptive-design-slug","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
