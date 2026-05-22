# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Choose a descriptive `child_id` that starts with `SIM-nonol-child-` and ends with a kebab-case design slug. Do not use generic `ga-child`, `pending`, or ordinal-only names.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260522-012332`
- Planned child ID prefix: `SIM-nonol-child-`
- Temporary child ID: `SIM-nonol-child-pending-0001`
- Temporary child path: `proposals/ga-canary-20260522-012332/simulations/SIM-nonol-child-pending-0001/`
- Operation: `breed`
- Parent IDs: `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0, SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0`

## Scenario Sample

- `bgp-class-routing-app-conflicting-policies` at `scenarios/bgp-class-routing-app-conflicting-policies/bgp-class-routing-app-conflicting-policies.md`
- `bgp-class-routing-app-honest-reachability-promise` at `scenarios/bgp-class-routing-app-honest-reachability-promise/bgp-class-routing-app-honest-reachability-promise.md`
- `bgp-class-routing-app-route-leak` at `scenarios/bgp-class-routing-app-route-leak/bgp-class-routing-app-route-leak.md`
- `bgp-class-routing-app-sparse-knowledge` at `scenarios/bgp-class-routing-app-sparse-knowledge/bgp-class-routing-app-sparse-knowledge.md`
- `cas-backed-group-session-envelope-independence` at `scenarios/cas-backed-group-session-envelope-independence/cas-backed-group-session-envelope-independence.md`
- `cas-object-type-binding-bakeoff-unknown-typed-object` at `scenarios/cas-object-type-binding-bakeoff-unknown-typed-object/cas-object-type-binding-bakeoff-unknown-typed-object.md`
- `cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes` at `scenarios/cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes/cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes.md`
- `chunk-feed-replication-carrier-independence` at `scenarios/chunk-feed-replication-carrier-independence/chunk-feed-replication-carrier-independence.md`
- `chunk-feed-replication-corrupt-chunk` at `scenarios/chunk-feed-replication-corrupt-chunk/chunk-feed-replication-corrupt-chunk.md`
- `chunk-feed-replication-duplicate-advertisement` at `scenarios/chunk-feed-replication-duplicate-advertisement/chunk-feed-replication-duplicate-advertisement.md`
- `conditional-release-geofencing-opaque-lower-layer-carriage` at `scenarios/conditional-release-geofencing-opaque-lower-layer-carriage/conditional-release-geofencing-opaque-lower-layer-carriage.md`
- `conditional-release-geofencing-replay-outside-conditions` at `scenarios/conditional-release-geofencing-replay-outside-conditions/conditional-release-geofencing-replay-outside-conditions.md`
- `transport-family-bakeoff-per-hop-authorization-failure` at `scenarios/transport-family-bakeoff-per-hop-authorization-failure/transport-family-bakeoff-per-hop-authorization-failure.md`
- `udp-feed-v0-conformance-malformed-datagram` at `scenarios/udp-feed-v0-conformance-malformed-datagram/udp-feed-v0-conformance-malformed-datagram.md`
- `udp-feed-v0-conformance-session-layer-composition` at `scenarios/udp-feed-v0-conformance-session-layer-composition/udp-feed-v0-conformance-session-layer-composition.md`

## Scenario Pressure

### `scenarios/bgp-class-routing-app-conflicting-policies/bgp-class-routing-app-conflicting-policies.md`

```markdown
# Conflicting policies

## Scenario ID

bgp-class-routing-app-conflicting-policies

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-punaz-bgp-class-routing-app/SCENARIOS.md`
- Source simulation: `SIM-punaz-bgp-class-routing-app/`
- Source row/title: Conflicting policies
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-punaz-bgp-class-routing-app/`.

## Setup

Alice prefers paths that avoid a jurisdiction; Carol prefers cheapest path; Bob has both offers.

## Stimulus

Run the candidate simulation against this source test: Whether route choice can be policy-relative instead of globally best.

## Expected Pressure

A PromiseGrid routing app should support peer-specific preference and refusal semantics.
```

### `scenarios/bgp-class-routing-app-honest-reachability-promise/bgp-class-routing-app-honest-reachability-promise.md`

```markdown
# Honest reachability promise

## Scenario ID

bgp-class-routing-app-honest-reachability-promise

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-punaz-bgp-class-routing-app/SCENARIOS.md`
- Source simulation: `SIM-punaz-bgp-class-routing-app/`
- Source row/title: Honest reachability promise
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-punaz-bgp-class-routing-app/`.

## Setup

Alice advertises reachability to Carol through Bob, and Bob later forwards as promised.

## Stimulus

Run the candidate simulation against this source test: What Alice, Bob, and Carol each record locally after the path works.

## Expected Pressure

Route-like promises need observable kept/broken outcomes without a global route authority.
```

### `scenarios/bgp-class-routing-app-route-leak/bgp-class-routing-app-route-leak.md`

```markdown
# Route leak

## Scenario ID

bgp-class-routing-app-route-leak

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-punaz-bgp-class-routing-app/SCENARIOS.md`
- Source simulation: `SIM-punaz-bgp-class-routing-app/`
- Source row/title: Route leak
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-punaz-bgp-class-routing-app/`.

## Setup

Bob repeats Alice's restricted transit offer outside the intended policy context.

## Stimulus

Run the candidate simulation against this source test: Whether policy scope and onward-restraint promises can be represented and audited.

## Expected Pressure

Routing policy needs more than reachability; it needs promise scope and violation evidence.
```

### `scenarios/bgp-class-routing-app-sparse-knowledge/bgp-class-routing-app-sparse-knowledge.md`

```markdown
# Sparse knowledge

## Scenario ID

bgp-class-routing-app-sparse-knowledge

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-punaz-bgp-class-routing-app/SCENARIOS.md`
- Source simulation: `SIM-punaz-bgp-class-routing-app/`
- Source row/title: Sparse knowledge
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-punaz-bgp-class-routing-app/`.

## Setup

Alice knows only Bob and Carol; Bob knows Dave and Ellen; no peer has the whole graph.

## Stimulus

Run the candidate simulation against this source test: Whether multi-hop discovery can find acceptable paths without requiring full topology replication.

## Expected Pressure

BGP-class pressure must compose with sparse-CAS and sparse relationship knowledge.
```

### `scenarios/cas-backed-group-session-envelope-independence/cas-backed-group-session-envelope-independence.md`

```markdown
# Envelope independence

## Scenario ID

cas-backed-group-session-envelope-independence

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jurar-cas-backed-group-session/SCENARIOS.md`
- Source simulation: `SIM-jurar-cas-backed-group-session/`
- Source row/title: Envelope independence
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jurar-cas-backed-group-session/`.

## Setup

The message is wrapped by one candidate grid-envelope variant in one experiment and another in a different experiment.

## Stimulus

Run the candidate simulation against this source test: Whether group-session semantics depend only on resolved payload meaning, not on a chosen grid-envelope winner.

## Expected Pressure

This sim must not backdoor a preferred grid-envelope variant.
```

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

### `scenarios/cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes/cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes.md`

```markdown
# Raw chunk versus pointer bytes

## Scenario ID

cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kohad-cas-object-type-binding-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-kohad-cas-object-type-binding-bakeoff/`
- Source row/title: Raw chunk versus pointer bytes
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kohad-cas-object-type-binding-bakeoff/`.

## Setup

Alice receives bytes whose hash is known, but the local filename is missing.

## Stimulus

Run the candidate simulation against this source test: Whether CID codec identity alone tells Bob whether to parse the bytes as a pointer object or treat them as raw payload.

## Expected Pressure

Object type must survive transport without relying on local paths.
```

### `scenarios/chunk-feed-replication-carrier-independence/chunk-feed-replication-carrier-independence.md`

```markdown
# Carrier independence

## Scenario ID

chunk-feed-replication-carrier-independence

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`
- Source simulation: `SIM-zazit-chunk-feed-replication/`
- Source row/title: Carrier independence
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zazit-chunk-feed-replication/`.

## Setup

The same chunk exchange is attempted over UDP, git, libp2p, IPFS, or ATPROTO-adjacent carriers.

## Stimulus

Run the candidate simulation against this source test: Which semantics belong to the feed role and which are carrier mechanics.

## Expected Pressure

The simulation should preserve turn-177's claim that feeds move chunks independent of substrate.
```

### `scenarios/chunk-feed-replication-corrupt-chunk/chunk-feed-replication-corrupt-chunk.md`

```markdown
# Corrupt chunk

## Scenario ID

chunk-feed-replication-corrupt-chunk

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`
- Source simulation: `SIM-zazit-chunk-feed-replication/`
- Source row/title: Corrupt chunk
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zazit-chunk-feed-replication/`.

## Setup

Mallory advertises or sends bytes whose hash does not match CID C.

## Stimulus

Run the candidate simulation against this source test: Whether rejection, accounting, and retry behavior are local enough to avoid central enforcement.

## Expected Pressure

Feed behavior must compose with CAS hash verification and peer-local accounting records.
```

### `scenarios/chunk-feed-replication-duplicate-advertisement/chunk-feed-replication-duplicate-advertisement.md`

```markdown
# Duplicate advertisement

## Scenario ID

chunk-feed-replication-duplicate-advertisement

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`
- Source simulation: `SIM-zazit-chunk-feed-replication/`
- Source row/title: Duplicate advertisement
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zazit-chunk-feed-replication/`.

## Setup

Alice and Carol both advertise chunk C.

## Stimulus

Run the candidate simulation against this source test: Whether duplicate offers are harmless and how Bob chooses between peers.

## Expected Pressure

Promise accounting can influence peer choice without making the feed a central reputation service.
```

### `scenarios/conditional-release-geofencing-opaque-lower-layer-carriage/conditional-release-geofencing-opaque-lower-layer-carriage.md`

```markdown
# Opaque lower-layer carriage

## Scenario ID

conditional-release-geofencing-opaque-lower-layer-carriage

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zarud-conditional-release-geofencing/SCENARIOS.md`
- Source simulation: `SIM-zarud-conditional-release-geofencing/`
- Source row/title: Opaque lower-layer carriage
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zarud-conditional-release-geofencing/`.

## Setup

Bob's node stores encrypted content whose condition vocabulary it cannot parse.

## Stimulus

Run the candidate simulation against this source test: Whether lower layers can safely carry opaque condition references while avoiding accidental promise violations.

## Expected Pressure

If opaque carriage is acceptable, the condition object must be verifiable without every layer understanding its semantics.
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

### `scenarios/udp-feed-v0-conformance-malformed-datagram/udp-feed-v0-conformance-malformed-datagram.md`

```markdown
# Malformed datagram

## Scenario ID

udp-feed-v0-conformance-malformed-datagram

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kuful-udp-feed-v0-conformance/SCENARIOS.md`
- Source simulation: `SIM-kuful-udp-feed-v0-conformance/`
- Source row/title: Malformed datagram
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kuful-udp-feed-v0-conformance/`.

## Setup

Bob receives arbitrary bytes that do not parse at higher layers.

## Stimulus

Run the candidate simulation against this source test: Whether UDP-feed passes bytes upward unchanged rather than inventing message semantics.

## Expected Pressure

Binding conformance must stay below session semantics.
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

### `simulations/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/README.md`

```markdown
# SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0: Grid-envelope variant

This simulation is a standalone positional grid-envelope specimen. It tests the
combination `enc-dag-cbor`, `unknown-opaque`, and `sig-unsigned-v0` without claiming
that this combination is the canonical PromiseGrid wire format. Source: `DI-fanah`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
```

### `simulations/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/QUESTION.md`

```markdown
# Question

Does a positional grid envelope using `enc-dag-cbor`, `unknown-opaque`, and
`sig-unsigned-v0` satisfy the wire-lab harness scenarios better than the sibling
variants? Source: `DI-fanah`.
```

### `simulations/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries
yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen created by `DI-fanah`.
```

### `simulations/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid Envelope Variant Spec (DRAFT)

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `enc-dag-cbor` / `unknown-opaque` / `sig-unsigned-v0`.
> Source: `DI-fanah`.

## Purpose

This spec defines one full positional grid-envelope candidate for wire-lab
comparison. It is a specimen inside `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0`, not a harness rule and not the
canonical PromiseGrid wire format.

## Positional Envelope Shape

The envelope shape for this variant is:

```text
[pcid, payload]
```

Slots are interpreted positionally:

- `pcid` identifies the protocol/spec/handler that interprets `payload`.
- `payload` is opaque bytes until interpreted by the handler named by `pcid`.

A `payload` may itself be the canonical bytes of another grid envelope when the
protocol named by `pcid` specifies recursive nesting. The outer grid envelope
does not prescribe the payload's internal organization beyond the bytes boundary.

## Encoding

This variant encodes the envelope as DAG-CBOR-compatible positional arrays. `pcid` and `sig_pcid`, when present, are DAG-CBOR Link values; `payload`, `signature`, and `sig_payload` are byte strings. The envelope remains positional: no map/object envelope fields are introduced. The canonical bytes for signing and hashing are the DAG-CBOR bytes of the exact positional array under this spec.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid`, it may store and forward the exact envelope bytes as opaque content, but interpretation fails with an explicit unsupported-pCID result. A receiver MUST NOT parse `payload` speculatively without the handler named by `pcid`.

## Signature and Authorship Policy

The base envelope has no signature requirement in v0. Authorship, integrity, and accountability must be supplied by the carrying transport, by the application protocol named in `pcid`, or by later successor variants. This tests the smallest possible positional envelope.

## Layering-Test Behavior

This variant answers the harness §1.3 layering scenarios as follows:

- Ordering disagreements are handled by the protocol named by `pcid`; the grid
  envelope only preserves the bytes and dispatch identity needed to make failures
  explicit.
- Forwarding, relay, or hop-local evidence is represented either by wrapper
  grid envelopes, by the payload protocol, or by the signature slots available in
  this variant.
- External or content-addressed body references live inside `payload` under the
  protocol named by `pcid`; the envelope can carry those payload bytes without
  understanding them.
- Incompatible interpretation rules fail visibly at the `pcid` dispatch boundary
  or under this variant's unknown-pCID policy.

## Non-Goals

This draft does not declare a winning envelope, does not define a central pCID
registry, does not freeze a final PromiseGrid signing scheme, and does not make
sibling grid-envelope variants obsolete.

## Freeze Gate

This draft can freeze only after at least one simulation run compares it against
sibling positional grid-envelope variants and Steve signs a merge/freeze promise
for this specific specimen.
```

### `simulations/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/README.md`

```markdown
# SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0: Grid-envelope variant

This simulation is a standalone positional grid-envelope specimen. It tests the
combination `enc-dag-cbor`, `unknown-hard-reject`, and `sig-unsigned-v0` without claiming
that this combination is the canonical PromiseGrid wire format. Source: `DI-fanah`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
```

### `simulations/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/QUESTION.md`

```markdown
# Question

Does a positional grid envelope using `enc-dag-cbor`, `unknown-hard-reject`, and
`sig-unsigned-v0` satisfy the wire-lab harness scenarios better than the sibling
variants? Source: `DI-fanah`.
```

### `simulations/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries
yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen created by `DI-fanah`.
```

### `simulations/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid Envelope Variant Spec (DRAFT)

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `enc-dag-cbor` / `unknown-hard-reject` / `sig-unsigned-v0`.
> Source: `DI-fanah`.

## Purpose

This spec defines one full positional grid-envelope candidate for wire-lab
comparison. It is a specimen inside `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0`, not a harness rule and not the
canonical PromiseGrid wire format.

## Positional Envelope Shape

The envelope shape for this variant is:

```text
[pcid, payload]
```

Slots are interpreted positionally:

- `pcid` identifies the protocol/spec/handler that interprets `payload`.
- `payload` is opaque bytes until interpreted by the handler named by `pcid`.

A `payload` may itself be the canonical bytes of another grid envelope when the
protocol named by `pcid` specifies recursive nesting. The outer grid envelope
does not prescribe the payload's internal organization beyond the bytes boundary.

## Encoding

This variant encodes the envelope as DAG-CBOR-compatible positional arrays. `pcid` and `sig_pcid`, when present, are DAG-CBOR Link values; `payload`, `signature`, and `sig_payload` are byte strings. The envelope remains positional: no map/object envelope fields are introduced. The canonical bytes for signing and hashing are the DAG-CBOR bytes of the exact positional array under this spec.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid`, the envelope is rejected at the envelope layer. The receiver may keep local diagnostics, but it MUST NOT accept, store, or forward the message as a valid grid-envelope message under this variant.

## Signature and Authorship Policy

The base envelope has no signature requirement in v0. Authorship, integrity, and accountability must be supplied by the carrying transport, by the application protocol named in `pcid`, or by later successor variants. This tests the smallest possible positional envelope.

## Layering-Test Behavior

This variant answers the harness §1.3 layering scenarios as follows:

- Ordering disagreements are handled by the protocol named by `pcid`; the grid
  envelope only preserves the bytes and dispatch identity needed to make failures
  explicit.
- Forwarding, relay, or hop-local evidence is represented either by wrapper
  grid envelopes, by the payload protocol, or by the signature slots available in
  this variant.
- External or content-addressed body references live inside `payload` under the
  protocol named by `pcid`; the envelope can carry those payload bytes without
  understanding them.
- Incompatible interpretation rules fail visibly at the `pcid` dispatch boundary
  or under this variant's unknown-pCID policy.

## Non-Goals

This draft does not declare a winning envelope, does not define a central pCID
registry, does not freeze a final PromiseGrid signing scheme, and does not make
sibling grid-envelope variants obsolete.

## Freeze Gate

This draft can freeze only after at least one simulation run compares it against
sibling positional grid-envelope variants and Steve signs a merge/freeze promise
for this specific specimen.
```

## Compact Fitness Evidence From This Run

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `bgp-class-routing-app-conflicting-policies`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/bgp-class-routing-app-conflicting-policies/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=1 promisegrid_alignment=3 auditability=2 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=18.00 normalized_0_100=51.00 confidence_0_1=0.83
- Rationale: Useful as boundary evidence, not as a direct answer to the routing scenario. The envelope cleanly separates transport/container concerns from routing-policy semantics, but it does not itself represent peer-specific route preference, refusal, or accountable route choice. The unsigned base envelope also weakens auditability for disputed policy decisions.
- Strengths:
  - Very clear layer boundary: routing policy belongs in the payload protocol selected by pCID.
  - Unknown-pCID handling and opaque forwarding support sparse knowledge and avoid speculative parsing.
  - Simple DAG-CBOR positional form is plausible to implement and compare across sibling variants.
- Weaknesses:
  - Does not model Alice/Carol policy-relative route selection or refusal semantics.
  - Base envelope has no signature, so authorship and integrity are not inherent.
  - Opaque payload limits direct audit unless the payload protocol adds explicit evidence.
- Risks:
  - A routing app could overclaim support for conflicting policies while the envelope alone provides none.
  - Unsigned envelopes make route-policy disputes easier to spoof or harder to audit.
  - Long-term usability depends on preserving handler/spec knowledge for pCIDs.
- Open questions:
  - Which payload protocol carries peer-specific preference and refusal semantics for route choice?
  - What signature or receipt layer makes Bob's route offers accountable to Alice and Carol?
  - How are pCID interpretation rules preserved and discovered over long time horizons without central-registry drift?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `bgp-class-routing-app-honest-reachability-promise`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/bgp-class-routing-app-honest-reachability-promise/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=1 promisegrid_alignment=3 auditability=2 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=22.00 normalized_0_100=55.00 confidence_0_1=0.81
- Rationale: For the honest reachability scenario, this specimen is a good carrier but not a routing-promise design: it gives clean dispatch and opaque forwardability, yet leaves the actual reachability claim, forwarding evidence, and accountability to other layers.
- Strengths:
  - Very clear [pcid, payload] boundary with explicit no-speculative-parsing behavior.
  - Unknown-pCID opaque store/forward supports sparse knowledge and avoids central route-authority assumptions.
  - Canonical DAG-CBOR positional bytes are simple to implement and preserve as local evidence.
- Weaknesses:
  - No route advertisement, forwarding receipt, or kept/broken promise schema for Alice, Bob, and Carol.
  - Unsigned v0 envelope cannot by itself attribute or integrity-protect reachability claims.
  - Scenario outcome evidence depends on external payload, wrapper, or transport conventions.
- Risks:
  - Spoofed or replayed reachability messages may be mistaken for honest advertisements unless higher layers add signatures and anti-replay evidence.
  - Later auditors may retain bytes but lose durable sender or handler resolution if pCID/spec preservation is weak.
- Open questions:
  - What payload or wrapper protocol records Bob's forwarding and Carol's successful reachability observation locally?
  - How are authorship, integrity, and replay protection supplied for route-like promises in this variant?
  - How is long-term pCID/spec resolution preserved for 100-year audit without central registry dependence?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `bgp-class-routing-app-route-leak`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/bgp-class-routing-app-route-leak/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=2 promisegrid_alignment=2 auditability=1 evolution_safety=4 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=20.00 normalized_0_100=50.00 confidence_0_1=0.79
- Rationale: Useful as a minimal carriage layer, but a weak direct answer to route-leak pressure. The specimen gives clear pCID/payload boundaries, exact-byte preservation, and explicit unknown-pCID failure, yet policy scope, onward-restraint, and most audit evidence are left to payloads, wrappers, or transport. Because the base variant is unsigned, leak attribution and violation proof are fragile unless another layer supplies them.
- Strengths:
  - Clear layer boundary: the envelope carries bytes and dispatch identity without speculative parsing.
  - Opaque store-and-forward plus explicit unsupported-pCID behavior fits sparse-knowledge evolution.
  - Small DAG-CBOR positional form is straightforward to implement.
- Weaknesses:
  - No built-in signature or authorship, so envelope-level evidence for who leaked a restricted route is weak.
  - No native structure for policy scope or onward-restraint promises; the scenario depends almost entirely on payload design.
  - Unknown or opaque payloads reduce human and LLM auditability unless the handler remains available.
- Risks:
  - Deployments may confuse successful carriage with authorized policy propagation.
  - Unsigned payload/wrapper use can make route-leak evidence ambiguous, forgeable, or easy to strip.
  - Long-term audit depends on external payload specs and handler preservation rather than the envelope alone.
- Open questions:
  - Which payload protocol encodes restricted transit scope and violation evidence for this scenario?
  - Should routing-sensitive uses require a signed successor envelope or mandatory signed wrapper convention?
  - How will pCID resolution and handler availability be preserved for 100-year auditability?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `bgp-class-routing-app-sparse-knowledge`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/bgp-class-routing-app-sparse-knowledge/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=2 evolution_safety=3 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=14.00 normalized_0_100=48.00 confidence_0_1=0.79
- Rationale: For BGP-class sparse-knowledge pressure, this variant is mostly a clean substrate rather than a solution: unknown-pCID opaque forwarding and no central registry fit decentralized relay, but the envelope says nothing about discovery or acceptable-path evaluation and, being unsigned, does not by itself give durable route provenance.
- Strengths:
  - Explicit [pcid, payload] dispatch boundary keeps envelope and routing semantics separate.
  - Unsupported-pCID handling and opaque store/forward compose with sparse knowledge without requiring full graph or global understanding.
  - The draft explicitly avoids a central pCID registry, matching no-central-authority pressure.
  - ... 1 more
- Weaknesses:
  - It does not define multi-hop discovery, path acceptability, or route-selection logic.
  - Unsigned v0 gives no envelope-level provenance, integrity, or freshness evidence for route claims.
  - Peer-local audit can preserve bytes received but not who authored or endorsed the routing information.
  - ... 1 more
- Risks:
  - Mallory can forge, tamper with, or replay routing payloads unless another layer adds strong signatures and timestamps.
  - Sparse-knowledge peers may relay opaque route objects they cannot validate, spreading bad information.
  - Security may drift into transport-specific assumptions, weakening durable store/forward behavior.
- Open questions:
  - What payload or wrapper format carries signed path and policy evidence so Alice can judge routes with only partial graph knowledge?
  - Do relay wrappers need hop-local attestations or anti-replay data to make opaque multi-hop routing safe?
  - How should a future authenticated successor interoperate with stored unsigned v0 envelopes during migration?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `cas-backed-group-session-envelope-independence`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/cas-backed-group-session-envelope-independence/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=4 promisegrid_alignment=4 auditability=3 evolution_safety=4 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=5 risk_penalty=3
- Fitness: raw=26.00 normalized_0_100=74.00 confidence_0_1=0.78
- Rationale: Strong envelope/payload separation makes this a good fit for envelope-independence pressure: meaning lives under the pcid-selected payload handler, unknown envelopes stay opaque, and the spec explicitly avoids declaring a canonical winner. It is not a full pass because it does not define how a CAS-backed group session proves semantic equivalence across different outer-envelope bytes, and unsigned v0 pushes integrity and accountability into other layers.
- Strengths:
  - Meaning is explicitly owned by the pcid-selected payload protocol rather than the outer envelope.
  - Unknown-pCID handling preserves exact bytes and fails explicitly without speculative parsing.
  - The spec explicitly rejects a central registry and avoids treating this variant as the winner.
- Weaknesses:
  - No explicit rule shows how the same group-session message stays identical across different envelope variants.
  - Unsigned v0 weakens provenance and peer-local accountability at the envelope layer.
  - Canonical hashing of exact envelope bytes can still couple identifiers to the chosen wrapper.
- Risks:
  - Upper layers may accidentally content-address the outer DAG-CBOR envelope, reintroducing envelope dependence.
  - Transport-supplied authorship could become a de facto authority or vary between deployments.
  - Peers may be able to store and forward opaque bytes yet still lack enough evidence to audit payload meaning.
- Open questions:
  - What should a CAS-backed session bind to when two envelope variants carry the same resolved payload?
  - How should peers compare or deduplicate semantically equivalent cross-variant messages using only local evidence?
  - Is the layering note about signature slots meant to rely on wrappers or successor variants despite the base [pcid, payload] shape?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `cas-object-type-binding-bakeoff-unknown-typed-object`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/cas-object-type-binding-bakeoff-unknown-typed-object/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=3 evolution_safety=4 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=28.00 normalized_0_100=70.00 confidence_0_1=0.80
- Rationale: Good on the core safety behavior: it lets a peer store and forward exact bytes and return an explicit unsupported-pCID result without speculative parsing. But this is only a partial answer to the scenario, because the scenario is about unknown CID codec/object binding, while this envelope mostly defers that binding, advertisement, and provenance work to payload or transport layers.
- Strengths:
  - Explicit unknown-pCID policy allows exact-byte opaque store/forward.
  - MUST NOT parse unknown payloads speculatively, which fits the scenario's safety pressure.
  - The [pcid, payload] boundary is very clear and does not depend on a central registry.
- Weaknesses:
  - It handles unknown protocol dispatch more directly than unknown CID codec binding.
  - Safe advertisement/discovery behavior for unknown objects is not specified here.
  - Base v0 envelope is unsigned, so provenance and accountability are deferred.
- Risks:
  - If payload protocols do not supply durable type binding, interoperability may collapse into blind byte ferrying.
  - Unsigned opaque forwarding preserves bytes but not trustworthy authorship or relay evidence.
- Open questions:
  - How should an unknown-CID object be stored and re-advertised when the peer never learns its codec?
  - Is pCID intended to be the durable type-binding mechanism for CAS objects, or is a separate binding needed?
  - Which layer supplies durable signatures or promise-accounting records for forwarded opaque objects?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=4 promisegrid_alignment=3 auditability=3 evolution_safety=4 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=29.00 normalized_0_100=73.00 confidence_0_1=0.78
- Rationale: This variant answers the scenario by refusing to let DAG-CBOR codec identity decide semantic type: the codec only says "envelope container", while the in-band pCID tells Bob how to interpret the payload. When the content address is over the envelope bytes, the hash binds pCID and payload together without filenames, and unknown-pCID handling makes non-support explicit. But raw-versus-pointer semantics are still delegated to payload protocols, and the draft unsigned specimen is not yet a durable standalone PromiseGrid answer.
- Strengths:
  - Explicit pCID dispatch cleanly separates container codec from payload meaning.
  - Canonical envelope bytes can bind pCID and payload under one content address.
  - Unknown-pCID opaque forwarding supports sparse-knowledge relay and future evolution without speculative parsing.
- Weaknesses:
  - The envelope does not itself define whether the payload is a raw chunk or a pointer object.
  - The specimen is draft and its own pCID is not yet frozen or minted.
  - Unsigned v0 leaves provenance and accountability to transport or higher layers.
- Risks:
  - If bare payload bytes circulate outside the envelope, the type-binding benefit disappears.
  - pCID minting or resolution could become a hidden authority point or decay over long horizons.
  - Implementers may still overread DAG-CBOR codec identity and bypass the intended pCID boundary.
- Open questions:
  - Must every transport-visible CAS object be wrapped in a grid envelope for this guarantee to hold?
  - How are pCIDs minted, resolved, and archived over 100-year horizons without a central registry?
  - What signature or receipt layer is expected to provide durable audit evidence for pCID-bound objects?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `chunk-feed-replication-carrier-independence`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/chunk-feed-replication-carrier-independence/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=2 promisegrid_alignment=2 auditability=2 evolution_safety=4 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=16.00 normalized_0_100=46.00 confidence_0_1=0.81
- Rationale: This variant cleanly separates byte carriage from protocol interpretation, which helps carrier independence at the envelope layer, but it is too minimal and unsigned to define carrier-independent chunk-feed replication semantics on its own.
- Strengths:
  - Carrier-agnostic DAG-CBOR bytes can ride over many substrates without changing the outer envelope.
  - Explicit pCID dispatch and unknown-pCID store/forward support sparse knowledge and opaque forwarding.
  - The spec clearly pushes ordering and chunk/body semantics out of the envelope and into payload protocols or wrappers.
- Weaknesses:
  - It does not define feed-role or chunk-replication semantics, so the core scenario pressure is only partially covered.
  - Unsigned v0 outsources integrity and authorship to transport or payload, weakening true carrier independence.
  - Peer-local audit evidence is thin when payloads are opaque or handlers are unavailable.
- Risks:
  - Different carriers may accidentally become part of the trust model because the envelope itself provides no signature or integrity proof.
  - Tamper, spoof, or replay problems are easy to miss unless every payload protocol adds its own protections.
  - Long-term auditability may degrade if durable bytes survive but the needed handler or evidence layer does not.
- Open questions:
  - Which payload pCID, if any, standardizes chunk-feed semantics independently of carrier choice?
  - How are chunk integrity, authorship, and ordering proven when the carrier offers none of them?
  - Should wrapper envelopes or receipts carry the accountability evidence missing from unsigned v0?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `chunk-feed-replication-corrupt-chunk`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/chunk-feed-replication-corrupt-chunk/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=2 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=19.00 normalized_0_100=54.30 confidence_0_1=0.81
- Rationale: The envelope can preserve exact bytes and hand them to a payload/feed protocol that locally checks CID C, but it does not define that CID binding, the accounting records, retry behavior, or envelope-level attribution. It is stronger as a clean boundary specimen than as a complete corrupt-chunk design answer.
- Strengths:
  - Clear pcid/payload boundary with explicit unsupported-pCID failure.
  - Exact DAG-CBOR bytes can be stored and forwarded for later payload/CAS verification.
  - Supports sparse knowledge and evolution without a central registry via opaque unknown-pCID handling.
- Weaknesses:
  - Does not define feed semantics, CID commitments, peer-local accounting, or retry rules.
  - Unsigned v0 provides no envelope-level sender binding or accountability.
  - Opaque relays may forward bad bytes without being able to assess corruption.
- Risks:
  - Implementers may mistake successful envelope decoding for integrity unless payload-layer hash checks are mandatory.
  - Corrupt-chunk rejection can be local, but blame and anti-abuse evidence depend on transport or payload protocols.
  - Long-term audit trails are weaker because the base envelope carries no signature.
- Open questions:
  - Which payload protocol binds expected CID C to delivered bytes and records local reject/retry outcomes?
  - What transport or wrapper signature scheme preserves attribution when Mallory serves mismatched bytes?
  - Should a later signed envelope variant carry portable rejection or accounting evidence for corrupt-chunk cases?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `chunk-feed-replication-duplicate-advertisement`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/chunk-feed-replication-duplicate-advertisement/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=1 promisegrid_alignment=3 auditability=2 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=22.00 normalized_0_100=55.00 confidence_0_1=0.82
- Rationale: This variant keeps the wire layer narrow and non-centralized, but duplicate-advertisement behavior is almost entirely outside the envelope. Bob can receive and preserve multiple offers, yet the spec provides no built-in identity, freshness, or promise-accounting evidence to choose between Alice and Carol.
- Strengths:
  - Clear layer boundary keeps peer-choice and reputation logic out of the envelope and avoids hidden centralization.
  - Unknown-pCID opaque forwarding supports sparse knowledge, long-term preservation, and protocol evolution.
  - Minimal DAG-CBOR positional encoding is simple enough to implement and compare across variants.
- Weaknesses:
  - The envelope does not define advertiser identity, chunk-advertisement semantics, freshness, or peer-selection rules.
  - Unsigned v0 cannot prove that Alice and Carol are distinct accountable advertisers.
  - No native replay detection or local promise-accounting record format is provided at this layer.
- Risks:
  - Mallory can spoof or replay duplicate advertisements unless transport or payload authentication is added elsewhere.
  - Bob may overcount or mis-rank peers because envelope evidence alone is not accountable.
  - Forwarding can preserve offers without preserving trustworthy provenance.
- Open questions:
  - Which payload protocol defines chunk advertisement content, freshness, and dedupe behavior?
  - How does Bob bind duplicate advertisements to distinct peers and rank them using only local promise accounting?
  - For 100-year auditability, do signatures live in transport, payload, or a successor envelope variant?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `conditional-release-geofencing-opaque-lower-layer-carriage`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/conditional-release-geofencing-opaque-lower-layer-carriage/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=4 promisegrid_alignment=3 auditability=2 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=26.00 normalized_0_100=65.00 confidence_0_1=0.83
- Rationale: This variant matches the scenario well at the carriage boundary: lower layers can keep condition-bearing bytes opaque, preserve exact bytes, and fail explicitly on unknown pCIDs. It is only a partial end-to-end answer, though, because unsigned v0 does not itself make the condition object authentic, current, or auditable enough to avoid accidental promise violations.
- Strengths:
  - Unknown-pCID handling explicitly permits exact store-and-forward of opaque bytes without speculative parsing.
  - The [pcid, payload] boundary is very clear about what lower layers may carry versus what higher layers must interpret.
  - Canonical DAG-CBOR bytes and payload-level external/content-addressed references give higher layers a plausible place to attach verifiable condition material.
- Weaknesses:
  - The base envelope has no required signature or integrity protection.
  - The spec does not define how an opaque condition object becomes verifiable to carriers or later auditors.
  - Conditional-release semantics and peer-local promise accounting are left outside the envelope.
- Risks:
  - Opaque condition references can be replayed, swapped, or stale unless higher layers bind them cryptographically to the promised release conditions.
  - A node may be able to prove it carried bytes, but not that any release decision was actually authorized.
  - Long-term verification depends on preserving pCID interpretation material and whatever external authentication layer is used.
- Open questions:
  - What exact mechanism makes a geofencing condition reference verifiable when intermediary layers do not understand its semantics?
  - What local record should Bob keep to show he only stored or forwarded opaque material rather than authorizing release?
  - How are pCID handlers, spec snapshots, and authentication evidence preserved across long-term protocol evolution?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `conditional-release-geofencing-replay-outside-conditions`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/conditional-release-geofencing-replay-outside-conditions/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=1 evolution_safety=3 layer_boundary_clarity=5 failure_handling=1 implementation_plausibility=5 risk_penalty=5
- Fitness: raw=13.00 normalized_0_100=45.00 confidence_0_1=0.87
- Rationale: This specimen is a clear minimal wire envelope, but it is not a direct replay-defense design. Its unsigned two-slot form and unknown-opaque forwarding do not bind content to audience, geography, epoch, or session state, so a valid old envelope can be replayed verbatim unless the payload protocol, wrapper, or transport adds those checks. It helps by making the layer boundary explicit, not by solving the conditional-release problem itself.
- Strengths:
  - Very clear envelope/payload boundary with explicit pCID dispatch and no speculative parsing of unknown payloads.
  - Canonical DAG-CBOR bytes make higher-layer hashing, wrapping, and duplicate detection plausible.
  - Simple positional encoding is easy to implement and does not require central wire-layer authority.
- Weaknesses:
  - No envelope-level signature, freshness marker, recipient binding, geography binding, or session binding.
  - Replay detection is fully delegated to payload, wrapper, or transport layers.
  - Opaque forwarding preserves bytes but not policy meaning, limiting generic auditability for this scenario.
- Risks:
  - Mallory can replay a previously valid envelope outside the intended audience or geography unless higher layers keep local condition state.
  - If integrity or authorization context lives only in transport or session state, archived or relayed envelopes lose durable accountability.
  - Operators may confuse successful forwarding with authorized release.
- Open questions:
  - Should conditional-release constraints live in a signed payload protocol, a signed outer wrapper, or both?
  - What peer-local record lets Dave or Carol prove that an envelope was outside its allowed epoch, audience, or geography?
  - Should unknown-pCID store-and-forward be restricted for condition-bearing references to reduce replay spread?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `transport-family-bakeoff-per-hop-authorization-failure`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/transport-family-bakeoff-per-hop-authorization-failure/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=2 promisegrid_alignment=2 auditability=2 evolution_safety=3 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=3 risk_penalty=4
- Fitness: raw=13.00 normalized_0_100=37.00 confidence_0_1=0.86
- Rationale: This variant cleanly separates envelope dispatch from transport or payload policy and supports opaque forwarding, but it does not define ring membership, per-hop authorization proof, refusal recording, or reconfiguration behavior. For this scenario it mainly marks the layer where authorization failure must be handled rather than handling it itself, and the unsigned base envelope leaves accountability to outer layers.
- Strengths:
  - Clear `[pcid, payload]` boundary avoids speculative parsing and makes handler responsibility explicit.
  - Unknown-pCID opaque store-and-forward behavior supports sparse knowledge and future protocol evolution without a central registry.
  - Canonical DAG-CBOR bytes give higher layers a stable basis for hashing or signing.
- Weaknesses:
  - No native model for whether an unauthorized hop breaks the ring, skips a hop, records a refusal, or reconfigures membership.
  - Base `sig-unsigned-v0` provides no mandatory authorship or integrity evidence for who forwarded or refused.
  - The draft mentions optional signature-related fields despite the minimal two-slot envelope, leaving some behavior ambiguous.
- Risks:
  - Bob's refusal or an unauthorized forward may not be provable from the envelope alone.
  - Different transports or payload protocols may encode hop-authorization failure evidence incompatibly.
  - Users could mistake this minimal envelope for sufficient authorization evidence and lose peer-local accountability.
- Open questions:
  - What concrete wrapper or payload protocol records Bob's refusal so Carol can audit it later?
  - How is per-hop authorization bound to sender, recipient, and current ring membership over time?
  - Are signature-related slots actually part of this variant, or should all hop evidence live entirely outside the base envelope?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `udp-feed-v0-conformance-malformed-datagram`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/udp-feed-v0-conformance-malformed-datagram/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=3 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=26.00 normalized_0_100=65.00 confidence_0_1=0.76
- Rationale: Partial fit. The specimen is strong on layer separation: payload stays opaque, unknown pCID fails explicitly, and exact bytes can be retained, which matches the scenario's demand not to invent message semantics above a byte feed. But it is not a UDP-feed binding and it does not define how malformed outer bytes are surfaced or recorded, so malformed-datagram conformance is only partly evidenced.
- Strengths:
  - Payload remains opaque until a pCID handler is known.
  - Unknown pCID handling keeps exact bytes and returns an explicit unsupported-pCID outcome.
  - The positional DAG-CBOR envelope is simple to implement and evolve.
- Weaknesses:
  - The specimen does not define a UDP-feed transport binding.
  - Behavior for invalid outer bytes or wrong-shape arrays is unspecified.
  - Unsigned v0 leaves integrity and accountability to other layers.
- Risks:
  - Implementations may diverge on malformed outer datagram handling.
  - Unsigned opaque envelopes are easier to spoof, replay, or relay without attribution.
- Open questions:
  - What should a UDP-feed binding record locally when outer decoding fails before pCID dispatch?
  - Should malformed datagrams be passed upward unchanged, surfaced as parse-failure records, or both?
  - What bounds prevent opaque junk retention or forwarding abuse?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0` x `udp-feed-v0-conformance-session-layer-composition`

- Result path: `results/SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0/udp-feed-v0-conformance-session-layer-composition/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=2 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=20.00 normalized_0_100=63.00 confidence_0_1=0.73
- Rationale: Useful but incomplete fit: the envelope establishes a clean transport-agnostic [pcid, payload] boundary and explicit unsupported-pCID failure, which helps session-layer composition. But because sig-unsigned-v0 leaves integrity, authorship, ordering, and much of the real session contract outside the envelope, it does not by itself show that UDP-feed v0 is sufficient without extra binding assumptions.
- Strengths:
  - Transport-agnostic [pcid, payload] boundary keeps parsing separate from the carrier.
  - Unknown-pCID policy supports opaque store/forward and explicit failure under sparse knowledge.
  - No central pCID registry is required.
  - ... 1 more
- Weaknesses:
  - No built-in signature, sender binding, or durable accountability.
  - Ordering, relay evidence, and session semantics are deferred to payload or transport.
  - The specimen does not spell out what UDP-feed API guarantees the next layer actually needs.
- Risks:
  - Unsigned UDP-carried envelopes can be spoofed or replayed while still looking syntactically valid.
  - Deferring core semantics outward may force hidden transport coupling, defeating clean composition.
- Open questions:
  - What UDP-feed v0 API surface beyond raw byte delivery is required for clean session-layer composition?
  - Is a signed wrapper or mandatory signed payload needed before this variant meets PromiseGrid audit/accounting goals?
  - How are MTU, fragmentation, deduplication, and replay handled without leaking UDP specifics into the payload contract?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `bgp-class-routing-app-conflicting-policies`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/bgp-class-routing-app-conflicting-policies/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=2 evolution_safety=1 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=18.00 normalized_0_100=45.00 confidence_0_1=0.86
- Rationale: This specimen is mostly a minimal wire envelope, not a routing-policy design. It partially aligns with PromiseGrid by keeping the envelope small, payload-driven, and non-centralized, but Alice/Carol-specific route preference and refusal semantics live entirely above the envelope. DAG-CBOR bytes and explicit unknown-pCID rejection help clarity, while unsigned v0 and hard-reject behavior make durable auditing and long-term evolution weak for this scenario.
- Strengths:
  - Clear envelope/payload boundary: routing policy can live in the application protocol instead of being baked into the wire wrapper.
  - No central routing truth or central pCID registry is required at the envelope layer.
  - Canonical DAG-CBOR positional encoding is simple to implement and inspect.
  - ... 1 more
- Weaknesses:
  - No native model for peer-specific route preference, refusal, or policy-relative selection.
  - Unsigned v0 gives no envelope-level authorship or integrity for route offers.
  - Hard reject on unknown pCID is brittle for protocol evolution and mixed-version operation.
  - ... 1 more
- Risks:
  - Spoofed or replayed route-policy messages are possible unless the payload or transport adds strong signatures.
  - Version skew can cause outright rejection instead of graceful coexistence during upgrades.
  - Thin envelope evidence may be insufficient for later dispute resolution about conflicting policies.
- Open questions:
  - What routing payload format would let Bob record and prove distinct offers, refusals, and local accounting for Alice and Carol?
  - Is payload- or transport-level signing enough for durable routing audits, or is envelope-level signing needed?
  - Can relays safely carry future routing pCIDs under hard-reject rules, or will upgrades fragment the network?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `bgp-class-routing-app-honest-reachability-promise`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/bgp-class-routing-app-honest-reachability-promise/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=2 promisegrid_alignment=2 auditability=2 evolution_safety=1 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=13.00 normalized_0_100=37.00 confidence_0_1=0.80
- Rationale: As a minimal wire envelope, this variant can carry routing payloads and lets Alice, Bob, and Carol log canonical bytes plus the pCID dispatch point locally, but it does not itself define or authenticate reachability promises, forwarding receipts, or kept/broken outcomes. DAG-CBOR and no central registry help sparse local operation, while unsigned v0 and unknown-hard-reject weaken auditability, migration, and long-term mixed-version reachability.
- Strengths:
  - Very clear envelope/payload boundary: routing meaning stays in the payload protocol named by pCID.
  - Canonical DAG-CBOR bytes support stable local logging and hashing.
  - Unknown pCID fails closed instead of silently misinterpreting a route message.
- Weaknesses:
  - No envelope-level signature, authorship, or integrity for route advertisements or forwarding evidence.
  - No built-in schema for reachability promises, relays, receipts, or peer-local promise accounting.
  - Receivers that lack a handler must not accept, store, or forward the message as valid, limiting interoperability.
- Risks:
  - Spoofed or replayed reachability claims can look well-formed unless every payload or transport adds its own protection.
  - Hard reject on unknown pCID can cause upgrade partitions and dropped reachability during migration.
  - Later auditors may not be able to prove who promised or forwarded what if transport-side evidence is missing.
- Open questions:
  - Which payload protocol provides signed reachability advertisements and forwarding receipts for this scenario?
  - How do peers preserve and rediscover pCID meaning over decades without introducing a central authority?
  - Should relays be allowed to preserve unknown envelopes as opaque bytes during migration instead of hard-rejecting them?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `bgp-class-routing-app-route-leak`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/bgp-class-routing-app-route-leak/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=1 evolution_safety=1 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=16.00 normalized_0_100=40.00 confidence_0_1=0.78
- Rationale: The envelope cleanly separates dispatch from payload semantics, but the route-leak scenario needs durable policy-scope, provenance, and violation evidence. This unsigned minimal carrier can transport such evidence only if another layer defines it, so it underfits the scenario while making its boundary explicit.
- Strengths:
  - Explicit pCID dispatch gives a clear envelope/payload boundary.
  - Canonical DAG-CBOR bytes support stable local byte-level recording.
  - Unknown pCID hard-reject makes incompatible interpretation fail visibly.
  - ... 1 more
- Weaknesses:
  - The envelope itself does not represent policy scope, onward-restraint, identity, or relay-chain evidence.
  - Unsigned v0 provides no built-in authorship or accountability for leak claims.
  - Unknown-hard-reject makes mixed-version evolution brittle.
  - ... 1 more
- Risks:
  - Bob or Mallory can replay or fabricate route offers unless transport or payload adds strong provenance.
  - Carol may be unable to prove whether a scope violation occurred from envelope evidence alone.
  - Policy upgrades can partition deployments and cause older peers to discard new evidence because unknown pCIDs are hard-rejected.
- Open questions:
  - Where should route-policy scope and onward-restraint evidence live: payload protocol, wrapper envelope, or transport receipt?
  - What durable peer-local artifact replaces envelope signatures for 100-year auditability in adversarial routing cases?
  - How do future auditors recover pCID meaning and version transitions when this draft is not frozen and unknown pCIDs are hard-rejected?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `bgp-class-routing-app-sparse-knowledge`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/bgp-class-routing-app-sparse-knowledge/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=2 evolution_safety=1 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=17.00 normalized_0_100=43.00 confidence_0_1=0.83
- Rationale: As a minimal transport envelope, this variant can carry routing payloads without needing global topology knowledge, but it does not itself define the multi-hop discovery or path-evaluation behavior this scenario asks for. Hard rejection of unknown pCIDs and an unsigned base envelope make mixed-version sparse routing and durable audit dependent on higher layers.
- Strengths:
  - Small two-slot [pcid, payload] envelope keeps dispatch and payload responsibility relatively clear.
  - Canonical DAG-CBOR bytes aid deterministic parsing and hashing.
  - Does not impose a central registry or full-topology replication at the envelope layer.
- Weaknesses:
  - No route discovery, path selection, sparse-CAS handling, or promise-accounting behavior is specified.
  - Unknown pCIDs are hard-rejected instead of being retainable or relayable as opaque messages.
  - Unsigned v0 provides no built-in origin or integrity evidence.
- Risks:
  - Mixed-version peers may drop newer routing messages and stall discovery.
  - Route claims can be spoofed or later disputed unless the payload or transport adds strong signatures.
  - Future auditors may lose first-class access to unknown envelopes if only rejection diagnostics remain.
- Open questions:
  - Can one stable routing payload pCID carry enough signed provenance to make the unsigned outer envelope acceptable?
  - Should unknown envelopes be storable or forwardable as opaque bytes for migration and audit?
  - What CAS objects or feeds would hold route advertisements and the local evidence Alice or Carol uses to justify a path?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `cas-backed-group-session-envelope-independence`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/cas-backed-group-session-envelope-independence/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=4 promisegrid_alignment=3 auditability=3 evolution_safety=2 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=21.00 normalized_0_100=60.00 confidence_0_1=0.74
- Rationale: Moderately good fit: the spec keeps resolved payload meaning above the wrapper and explicitly avoids claiming canonical-winner status, but it does not prove cross-envelope group-session equivalence, and its unsigned v0 plus unknown-hard-reject stance weakens auditability and migration.
- Strengths:
  - Resolved meaning is delegated to the `pcid` payload handler rather than to envelope-specific semantics.
  - The specimen explicitly avoids claiming canonical-winner status and avoids requiring a central registry.
  - Simple DAG-CBOR positional encoding makes the boundary between wrapper and payload easy to implement and compare.
- Weaknesses:
  - It does not show a concrete group-session message wrapped equivalently across sibling envelope variants.
  - Unsigned v0 leaves integrity, authorship, and accountability outside the envelope.
  - Still a draft specimen, and its unknown-hard-reject stance makes long-horizon migration and mixed-version coexistence brittle.
- Risks:
  - If higher layers hash or sign the full wrapper bytes, envelope choice becomes semantically sticky and independence is undermined.
  - Receivers without the needed `pcid` handler cannot preserve messages as valid traffic under this variant, which can privilege already-supported envelopes.
  - Disputed group-session events may lack enough envelope-level evidence for later audit.
- Open questions:
  - Can a CAS-backed group-session payload keep the same higher-layer identity across sibling envelope variants?
  - How are stable `pcid` values coordinated across envelope variants without reintroducing hidden central authority?
  - Should unknown-`pcid` envelopes be retainable as opaque bytes for later reinterpretation rather than hard-rejected?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `cas-object-type-binding-bakeoff-unknown-typed-object`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/cas-object-type-binding-bakeoff-unknown-typed-object/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=2 promisegrid_alignment=2 auditability=2 evolution_safety=1 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=5 risk_penalty=4
- Fitness: raw=16.00 normalized_0_100=46.00 confidence_0_1=0.85
- Rationale: This specimen is an envelope-layer answer, not a full CAS object-binding design. Its explicit hard reject on unknown pCID safely avoids unsafe parsing, but it does not satisfy the scenario's main requirement to let mixed-version peers store, advertise, and forward unknown typed objects opaquely. Strong layering clarity and easy implementation are offset by poor evolution safety and weak auditability from unsigned v0.
- Strengths:
  - Explicit unknown-pCID reject rule prevents unsafe envelope interpretation.
  - Envelope/payload boundary is crisp; payload bytes stay opaque until handler dispatch.
  - No central pCID registry is assumed, and the positional DAG-CBOR format is straightforward to implement.
- Weaknesses:
  - Rejects rather than opaquely storing, advertising, or forwarding unknown types.
  - Defines unknown behavior for pCID dispatch, not for unknown CID codecs themselves.
  - Unsigned v0 leaves weak provenance and audit evidence.
- Risks:
  - Hard reject on future types can fragment mixed-version networks and force coordinated upgrades.
  - Depending on external transport or payload integrity leaves envelope-layer accountability ambiguous.
- Open questions:
  - Can a known wrapper pCID relay unknown inner typed objects while keeping the outer envelope valid?
  - How are pCID handlers discovered and archived over decades without de facto central coordination?
  - What durable signed evidence records why a peer rejected an unknown object?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=3 evolution_safety=2 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=20.00 normalized_0_100=57.00 confidence_0_1=0.72
- Rationale: This variant gives a path-independent answer by carrying an in-band `pcid` inside the hashed bytes, so Bob need not rely on filenames, but it does not make CID codec identity alone sufficient and leaves raw-versus-pointer semantics to the payload protocol. The layer split is very clear, yet unsigned v0 and unknown-hard-reject limit long-term evolution and accountability.
- Strengths:
  - In-band `pcid` binds interpretation to transported bytes rather than filenames or local paths.
  - Minimal `[pcid, payload]` structure makes the envelope/payload boundary explicit.
  - Unknown `pcid` hard rejection prevents silent misparse and type confusion.
- Weaknesses:
  - The scenario is only answered if the hashed object is the full envelope; bare payload bytes remain ambiguous here.
  - Raw-chunk versus pointer-object semantics are delegated to the payload protocol, not defined by the envelope itself.
  - Unsigned v0 leaves provenance and accountability to other layers.
- Risks:
  - Hard rejection of unknown `pcid` can block relay or archival of future object types.
  - 100-year interpretability depends on durable recovery of `pcid` specs/handlers.
  - If transport and payload layers do not add signatures, Mallory can spoof or replay envelopes.
- Open questions:
  - Should unknown `pcid` envelopes be storable or quarantinable even when not accepted?
  - Should raw-chunk versus pointer-object typing live in distinct payload `pcid`s, inner CAS codecs, or both?
  - What durable mechanism lets future peers recover `pcid` meaning without a central registry?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `chunk-feed-replication-carrier-independence`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/chunk-feed-replication-carrier-independence/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=2 promisegrid_alignment=2 auditability=3 evolution_safety=1 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=5 risk_penalty=4
- Fitness: raw=22.00 normalized_0_100=55.00 confidence_0_1=0.76
- Rationale: Strong as a carrier-neutral dispatch envelope, but only a partial answer to carrier-independent chunk-feed replication: it cleanly separates payload semantics from transport while leaving feed semantics, integrity, and authorship outside the envelope and hard-rejecting unknown pCIDs.
- Strengths:
  - Very clear split between envelope dispatch and payload/feed semantics, helping keep carrier mechanics out of the wire format.
  - Simple DAG-CBOR positional encoding is plausible to ship over many substrates.
  - Failure on unknown pCID is explicit rather than silent.
- Weaknesses:
  - The simulation does not define feed-role or chunk-replication semantics, so scenario coverage is limited.
  - Unsigned v0 pushes integrity, authorship, and accountability to carriers or payload protocols.
  - Hard reject on unknown pCID blocks generic store-and-forward or archival behavior across mixed-version peers.
- Risks:
  - Carrier-specific trust and integrity behavior can leak back into the design, weakening true substrate independence.
  - Protocol evolution may fragment replication because older peers cannot accept or relay newer envelopes.
- Open questions:
  - May blind relays or archives retain raw envelope bytes without violating the hard-reject rule?
  - What feed-level signing and promise-accounting layer would make the same chunk exchange auditable across UDP, git, libp2p, IPFS, and ATPROTO-like carriers?
  - How should pCID migration work over long time horizons when unknown formats are rejected at the envelope layer?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `chunk-feed-replication-corrupt-chunk`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/chunk-feed-replication-corrupt-chunk/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=2 evolution_safety=1 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=17.00 normalized_0_100=43.00 confidence_0_1=0.83
- Rationale: This specimen preserves local failure boundaries and avoids central enforcement, but it is mostly a boundary test rather than a full corrupt-chunk replication design. It can carry chunk-feed or CAS payloads and clearly leaves CID verification to the payload layer, yet it does not itself define corrupt-chunk rejection, peer-local accounting, or retry behavior. DAG-CBOR canonical bytes and explicit dispatch boundaries help, while unsigned v0 and unknown-hard-reject weaken auditability and long-term evolution for this scenario.
- Strengths:
  - Makes the envelope-vs-payload boundary explicit, including that content-addressed body references live inside the payload protocol.
  - Unknown protocol handling fails locally and visibly at the dispatch boundary, not by central enforcement.
  - Simple DAG-CBOR positional encoding is plausible to implement and gives deterministic bytes for hashing.
- Weaknesses:
  - Does not specify CID-mismatch handling, retry policy, or peer-local replication accounting.
  - Unsigned v0 provides no built-in authorship or integrity evidence for blaming a sender of corrupt data.
  - Hard rejection of unknown protocol IDs is brittle for long-lived migration.
- Risks:
  - Peers may detect bad bytes but still lack durable evidence about which peer supplied them after transport context is lost.
  - Future protocol evolution can strand old replicas that only know the hard-reject behavior.
  - Accountability may drift into ad hoc transport or central mechanisms, undercutting the scenario's local-enforcement goal.
- Open questions:
  - Which payload or wrapper protocol records CID mismatch events, retries, and peer-local trust/accounting state?
  - Should adversarial replication require a signed wrapper even if the base envelope remains unsigned?
  - How should newer protocol IDs roll out without making older replicas fail closed indefinitely?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `chunk-feed-replication-duplicate-advertisement`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/chunk-feed-replication-duplicate-advertisement/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=2 evolution_safety=1 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=5 risk_penalty=3
- Fitness: raw=15.00 normalized_0_100=43.00 confidence_0_1=0.84
- Rationale: This envelope specimen mainly demonstrates a clean boundary: duplicate chunk advertisements and Bob's peer-choice logic belong in the payload protocol or transport, not in the envelope. That keeps the wire format simple and non-centralized, but the scenario's main pressure—attributable duplicate offers and peer-local promise accounting—is largely unresolved here. Canonical DAG-CBOR bytes help local recording, while unknown-hard-reject and unsigned v0 weaken migration and provenance.
- Strengths:
  - Very clear envelope-versus-payload boundary.
  - Simple positional DAG-CBOR encoding is easy to implement and record locally.
  - Does not require a central registry or reputation service at the envelope layer.
- Weaknesses:
  - Does not define chunk advertisement semantics, duplicate handling, or Bob's peer-selection rules.
  - Unsigned v0 does not bind an advertisement to Alice or Carol at the envelope layer.
  - Unknown hard reject makes mixed-version and long-lived interoperability brittle.
- Risks:
  - Bob may rely on transport-specific or ad hoc trust signals that are not durable PromiseGrid evidence.
  - Forwarded or cached advertisements can lose trustworthy provenance unless higher layers add signatures or promise records.
  - Future advertisement variants may be dropped outright by older peers because unknown pCIDs are hard-rejected.
- Open questions:
  - What payload protocol carries chunk offers and the peer-local promise accounting Bob uses to choose between Alice and Carol?
  - How are duplicate advertisements authenticated and persisted for later audit without creating a central reputation service?
  - What migration path lets newer advertisement handlers coexist with older peers under unknown-hard-reject?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `conditional-release-geofencing-opaque-lower-layer-carriage`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/conditional-release-geofencing-opaque-lower-layer-carriage/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=2 evolution_safety=3 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=5 risk_penalty=3
- Fitness: raw=21.00 normalized_0_100=60.00 confidence_0_1=0.79
- Rationale: Strong as a lower-layer opaque carrier because payload stays opaque and pcid cleanly marks the interpretation boundary, but only a partial fit for this scenario: the envelope itself does not make conditional-release references verifiable, and the unsigned v0 base leaves integrity and accountability to payload or transport layers.
- Strengths:
  - Opaque payload bytes let lower layers carry condition references without parsing their semantics.
  - Explicit pcid dispatch and unknown-pcid hard reject make interpretation failures visible instead of silent.
  - Canonical DAG-CBOR bytes are a practical base for higher-layer hashing, CAS binding, or nested signatures.
- Weaknesses:
  - The base envelope has no native signature, integrity, or authorship proof.
  - Scenario-critical verification of the condition object is delegated entirely to payload or transport protocols.
  - A node without the relevant pcid handler cannot accept, store, or forward the message as a valid envelope, which limits generic opaque carriage.
- Risks:
  - If payload or transport does not add cryptographic binding, opaque condition references can be swapped, replayed, or misbound.
  - Operators may overestimate what the envelope itself guarantees and treat unauthenticated carriage as sufficient promise evidence.
  - Hard reject on unknown pcid can create brittle interop or archival gaps during long-term evolution.
- Open questions:
  - What payload-level structure lets Bob verify a condition object's integrity without understanding its geofencing semantics?
  - Should archival or relay peers be allowed to retain unknown-pcid envelopes as opaque bytes without treating them as valid messages?
  - Is a signed or content-addressed condition object mandatory for conditional-release use of this variant?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `conditional-release-geofencing-replay-outside-conditions`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/conditional-release-geofencing-replay-outside-conditions/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=1 evolution_safety=2 layer_boundary_clarity=5 failure_handling=1 implementation_plausibility=3 risk_penalty=5
- Fitness: raw=10.00 normalized_0_100=29.00 confidence_0_1=0.86
- Rationale: Clear as a minimal carrier, but weak for this scenario. The unsigned [pcid, payload] envelope has no native recipient, geography, epoch, or session binding, so replay protection is deferred to unspecified upper layers and local audit evidence stays thin.
- Strengths:
  - Very clear positional dispatch boundary.
  - Unknown pCID hard reject makes incompatible interpretation fail visibly.
  - Canonical DAG-CBOR bytes are simple to implement and can support higher-layer proofs.
- Weaknesses:
  - No envelope-level signature, freshness, audience binding, or geofence binding.
  - Replay detection is delegated to unspecified payload, transport, or session state.
  - Peer-local records can capture received bytes but not authoritative release conditions.
- Risks:
  - A previously valid content reference can be replayed unchanged unless upper layers block reuse.
  - Dave may be unable to prove later that reuse violated audience or geography conditions.
  - Teams may compensate with transport-local or central policy checks that weaken 100-year and no-central-authority goals.
- Open questions:
  - What payload or wrapper protocol binds content references to recipients, geography, and time windows?
  - How do peers record replay decisions locally without a central geofence authority?
  - What is the upgrade path to signed, freshness-carrying variants for hard-reject peers?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `transport-family-bakeoff-per-hop-authorization-failure`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/transport-family-bakeoff-per-hop-authorization-failure/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=1 evolution_safety=1 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=15.00 normalized_0_100=37.50 confidence_0_1=0.85
- Rationale: Useful mainly as boundary evidence: the specimen clearly keeps per-hop authorization out of the envelope, but therefore does not itself decide whether Bob records refusal, skips Carol, breaks the ring, or reconfigures membership. The DAG-CBOR positional form is simple, yet unsigned v0 plus unknown-hard-reject leave weak auditability and brittle evolution under this scenario.
- Strengths:
  - Clear envelope/payload boundary keeps hop authorization responsibility explicit.
  - Simple DAG-CBOR positional encoding is straightforward to implement and compare.
  - Unknown-format hard reject avoids silent misinterpretation.
- Weaknesses:
  - Per-hop authorization failure semantics are deferred to transport or payload layers, so the scenario's core question is mostly unanswered.
  - Unsigned v0 provides no envelope-level authorship or integrity for forwarding or refusal evidence.
  - Hard reject reduces compatibility and migration safety when peers lack a handler.
- Risks:
  - A ring can fail or fragment without a durable shared record of why Bob did not forward.
  - Unauthorized-forwarding disputes cannot be independently audited from envelope evidence alone.
  - Long-horizon version drift or handler loss can turn hard reject into interoperability breakage.
- Open questions:
  - Which layer emits Bob's durable peer-local refusal record?
  - How would Carol distinguish authorized refusal, accidental omission, and malicious drop without signed hop evidence?
  - Is there a forward-compatible path for new pCIDs in long-lived rings under unknown-hard-reject?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `udp-feed-v0-conformance-malformed-datagram`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/udp-feed-v0-conformance-malformed-datagram/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=3 promisegrid_alignment=2 auditability=3 evolution_safety=1 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=23.00 normalized_0_100=58.00 confidence_0_1=0.78
- Rationale: Indirect fit: this draft is not a UDP-feed binding, but it keeps meaning above the transport/envelope boundary by treating payload as opaque bytes and rejecting at the envelope parse/pCID boundary instead of inventing session semantics. Fitness is capped by missing UDP pass-through detail, evolution-hostile unknown-hard-reject, and unsigned-v0 audit gaps.
- Strengths:
  - Clear opaque-bytes payload boundary helps avoid inventing session semantics in lower layers.
  - Reject-at-envelope behavior keeps failure local and explicit.
  - Compact DAG-CBOR positional encoding is straightforward to implement.
- Weaknesses:
  - The draft does not explicitly specify a UDP-feed binding or byte-for-byte datagram pass-through.
  - Unknown-hard-reject is brittle for long-lived evolution and relay/storage of future protocols.
  - Unsigned v0 weakens provenance and durable audit evidence.
- Risks:
  - Injected malformed or unsigned traffic is rejectable but difficult to attribute afterward.
  - Hard rejection may discard bytes that would be useful for later migration or forensic review.
- Open questions:
  - What exact local record or quarantine behavior is required for non-DAG-CBOR or wrong-shape datagrams?
  - What transport binding guarantees UDP-feed passes raw bytes upward unchanged before envelope parsing?
  - Why does the draft mention sig_pcid/signature fields in an unsigned-v0 two-slot envelope, and what is the intended wire shape?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0` x `udp-feed-v0-conformance-session-layer-composition`

- Result path: `results/SIM-hiviv-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-unsigned-v0/udp-feed-v0-conformance-session-layer-composition/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=2 evolution_safety=2 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=19.00 normalized_0_100=54.00 confidence_0_1=0.78
- Rationale: Strong on transport-agnostic layer separation, but only a partial answer to the UDP-feed session-composition scenario because session behavior, integrity, and long-term evolution are left outside the envelope.
- Strengths:
  - Clear transport-agnostic [pcid,payload] boundary; the next layer need not know UDP-specific details.
  - No central registry or global state is assumed; receivers act from local handler knowledge.
  - Explicit dispatch and canonical DAG-CBOR bytes make parsing and failure evidence straightforward.
- Weaknesses:
  - It does not define session semantics or the exact UDP-feed API contract, so scenario coverage is only partial.
  - Unsigned v0 leaves authorship and integrity outside the envelope, limiting peer-local auditability.
  - Unknown-hard-reject is brittle for mixed-version rollout and generic store/forward behavior.
- Risks:
  - Unsigned envelopes invite spoofing or replay if transport or payload does not add protection.
  - New pCIDs may cause hard interop failures in long-lived mixed deployments.
  - Depending on carrier security can leak transport-specific assumptions into upper layers.
- Open questions:
  - Can intermediaries preserve unknown pCIDs opaquely without violating the reject rule?
  - What minimal UDP-feed API signals, beyond raw datagrams, does the next layer actually need?
  - Should integrity/authorship be mandatory in the envelope for session-layer use?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

## Required JSON Shape

{"child_id":"SIM-nonol-child-descriptive-design-slug","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
