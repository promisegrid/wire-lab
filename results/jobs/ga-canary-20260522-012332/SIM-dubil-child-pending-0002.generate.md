# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Choose a descriptive `child_id` that starts with `SIM-dubil-child-` and ends with a kebab-case design slug. Do not use generic `ga-child`, `pending`, or ordinal-only names.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260522-012332`
- Planned child ID prefix: `SIM-dubil-child-`
- Temporary child ID: `SIM-dubil-child-pending-0002`
- Temporary child path: `proposals/ga-canary-20260522-012332/simulations/SIM-dubil-child-pending-0002/`
- Operation: `breed`
- Parent IDs: `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid, SIM-hagom-grid-envelope-enc-dag-cbor-unknown-opaque-sig-unsigned-v0`

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

### `simulations/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/README.md`

```markdown
# SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid: Grid-envelope variant

This simulation is a standalone positional grid-envelope specimen. It tests the
combination `enc-cbor`, `unknown-hard-reject`, and `sig-wrapper-pcid` without claiming
that this combination is the canonical PromiseGrid wire format. Source: `DI-fanah`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
```

### `simulations/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/QUESTION.md`

```markdown
# Question

Does a positional grid envelope using `enc-cbor`, `unknown-hard-reject`, and
`sig-wrapper-pcid` satisfy the wire-lab harness scenarios better than the sibling
variants? Source: `DI-fanah`.
```

### `simulations/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries
yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen created by `DI-fanah`.
```

### `simulations/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid Envelope Variant Spec (DRAFT)

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `enc-cbor` / `unknown-hard-reject` / `sig-wrapper-pcid`.
> Source: `DI-fanah`.

## Purpose

This spec defines one full positional grid-envelope candidate for wire-lab
comparison. It is a specimen inside `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid`, not a harness rule and not the
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

This variant encodes the envelope as deterministic CBOR positional arrays. Slot values use definite-length encodings. `pcid` and `sig_pcid`, when present, are CIDv1 byte strings; `payload`, `signature`, and `sig_payload` are byte strings. The canonical bytes for signing and hashing are the deterministic CBOR bytes of the exact positional array under this spec.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid`, the envelope is rejected at the envelope layer. The receiver may keep local diagnostics, but it MUST NOT accept, store, or forward the message as a valid grid-envelope message under this variant.

## Signature and Authorship Policy

The base envelope has no fixed signature slot. Signatures, encryption, authorship, or hop evidence are represented by outer or inner grid envelopes whose own `pcid` selects the relevant signature or evidence protocol. This keeps the envelope shape minimal and tests whether pCID-selected wrapper protocols are enough for authorship and integrity.

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

## Compact Fitness Evidence From This Run

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `bgp-class-routing-app-conflicting-policies`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/bgp-class-routing-app-conflicting-policies/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=2 evolution_safety=1 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=15.00 normalized_0_100=43.00 confidence_0_1=0.88
- Rationale: Useful mainly as a carrier-layer boundary test, not as a direct answer to conflicting-policy routing. It keeps route policy out of the envelope and makes unknown-protocol failure visible, but it does not define the offers, local preference logic, or audit records this scenario needs. unknown-hard-reject plus wrapper-selected signatures also weakens migration safety.
- Strengths:
  - Clear separation between envelope dispatch and higher-layer routing policy.
  - Deterministic CBOR positional bytes support reproducible parse/sign/hash behavior.
  - Spec avoids a central pCID registry and makes unknown-pCID failure explicit.
- Weaknesses:
  - No model for Bob's competing offers or Alice/Carol policy-relative selection and refusal.
  - No peer-local accounting records for explaining route decisions to later auditors.
  - unknown-hard-reject is brittle for long-lived protocol and policy evolution.
- Risks:
  - New policy or signature wrapper pCIDs could split interoperability during migration.
  - Hard rejection may discard route evidence that a translator or auditor could still use.
  - The minimal envelope may be mistaken for an application-routing solution when it is only a carrier.
- Open questions:
  - What payload protocol carries route offers, refusals, and per-peer policy evidence?
  - How do older peers coexist with new wrapper pCIDs under unknown-hard-reject?
  - What local records should Alice, Bob, and Carol retain to audit divergent route choices over time?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `bgp-class-routing-app-honest-reachability-promise`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/bgp-class-routing-app-honest-reachability-promise/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=2 evolution_safety=1 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=13.00 normalized_0_100=37.00 confidence_0_1=0.82
- Rationale: Useful as a low-level carrier for route claims, but it does not itself define reachability promises, forwarding proof, or peer-local kept/broken accounting. Deterministic CBOR and explicit pCID dispatch help local recording and layer clarity, while unknown-pCID hard rejection and wrapper-only signatures weaken evolution and routing-evidence auditability.
- Strengths:
  - Deterministic CBOR gives stable bytes that Alice, Bob, and Carol can record locally.
  - Explicit pCID dispatch keeps the envelope/payload boundary clear.
  - The design does not assume a central pCID registry.
  - ... 1 more
- Weaknesses:
  - The spec does not define route-advertisement, forwarding, or promise-accounting semantics for this scenario.
  - Authorship and hop evidence depend on extra wrapper protocols instead of a fixed envelope field.
  - Hard rejection of unknown pCIDs makes mixed-version routing deployments brittle.
- Risks:
  - New pCIDs or signature wrappers can cause honest reachability messages to be dropped during upgrades.
  - Carol may lack enough standardized evidence to tell whether Bob kept the advertised promise.
  - Different wrapper choices may fragment audit practice across deployments.
- Open questions:
  - What payload and wrapper schema lets Alice, Bob, and Carol record a kept reachability promise in a comparable way?
  - Should older nodes be allowed to retain or relay unknown envelopes as opaque evidence instead of dropping them?
  - How are route claims linked to later forwarding outcomes across nested wrappers?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `bgp-class-routing-app-route-leak`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/bgp-class-routing-app-route-leak/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=2 promisegrid_alignment=2 auditability=2 evolution_safety=1 layer_boundary_clarity=4 failure_handling=3 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=20.00 normalized_0_100=50.00 confidence_0_1=0.76
- Rationale: Partial fit. This variant can carry signed routing-policy objects via payload and wrapper pCIDs, but it does not itself define policy scope, onward-restraint promises, or leak-proof evidence. Deterministic CBOR and clear dispatch boundaries help, while unknown-hard-reject is brittle for mixed-version routing and long-term evolution.
- Strengths:
  - Deterministic CBOR provides stable bytes for hashing and signatures.
  - Envelope dispatch is cleanly separated from routing-policy semantics.
  - Unknown pCID hard-reject makes unsupported policy/evidence formats fail visibly.
- Weaknesses:
  - The base envelope has no native structure for policy scope or route-leak evidence.
  - Auditability depends on extra wrapper and payload protocols, not the envelope alone.
  - Hard-reject on unknown pCIDs makes gradual rollout and future compatibility weak.
- Risks:
  - Peers may drop useful route-policy evidence when any required wrapper pCID is unknown.
  - Different signature-wrapper conventions could fragment interoperability.
  - Diagnostic retention for rejected envelopes may vary across implementations.
- Open questions:
  - Which payload and wrapper pCIDs would encode restricted transit offers and onward-restraint promises?
  - Can rejected unknown-pCID envelopes be preserved in a standard audit form without being treated as valid traffic?
  - How would this variant migrate safely as policy and signature wrappers change over decades?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `bgp-class-routing-app-sparse-knowledge`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/bgp-class-routing-app-sparse-knowledge/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=3 evolution_safety=1 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=21.00 normalized_0_100=52.50 confidence_0_1=0.74
- Rationale: Good wire-level substrate, weak scenario-complete answer. The minimal deterministic CBOR [pcid, payload] envelope preserves sparse-knowledge assumptions and clear layer boundaries, but BGP-class discovery lives entirely in payload protocols, and unknown-hard-reject makes mixed-version multi-hop route exchange brittle.
- Strengths:
  - The base envelope is compatible with sparse topology and sparse CAS because it does not require global graph knowledge and can carry content-addressed references inside payload protocols.
  - Deterministic CBOR bytes and explicit pCID dispatch give clear local evidence for hashing, signing, and parsing.
  - Layer separation is very clear: routing semantics, signatures, and hop evidence live in wrapper or payload protocols rather than the base envelope.
- Weaknesses:
  - Scenario coverage is indirect because the specimen does not define route advertisement, path selection, or acceptable-path policy.
  - unknown-hard-reject forbids accept/store/forward of unfamiliar routing or signature-wrapper pCIDs, making partial-upgrade multi-hop routing brittle.
  - Route provenance and audits depend on additional wrapper and payload specs beyond this draft.
- Risks:
  - Protocol-version drift can partition route discovery across sparsely connected peers when new routing or signature-wrapper pCIDs appear.
  - Long-term interpretability depends on durable availability of the pCID-linked specs/handlers; otherwise old envelopes become opaque.
  - Operational pressure to coordinate accepted pCIDs may recreate a de facto central compatibility authority.
- Open questions:
  - Is any opaque-retention or escrow path allowed for unknown pCIDs, or is hard drop mandatory even when later audit or replay would help?
  - How do peers advertise supported routing and signature-wrapper pCIDs in a sparse graph without introducing a central compatibility registry?
  - How would a routing payload encode acceptable-path policy and multi-hop discovery while preserving this strict layer split?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `cas-backed-group-session-envelope-independence`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/cas-backed-group-session-envelope-independence/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=3 evolution_safety=2 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=24.00 normalized_0_100=60.00 confidence_0_1=0.75
- Rationale: The enc-cbor / unknown-hard-reject / sig-wrapper-pcid specimen helps this scenario by keeping the envelope minimal and delegating meaning to the payload protocol. But it does not actually demonstrate cross-envelope invariance for a CAS-backed group session, and hard rejection plus wrapper-selected signature protocols can still make envelope choice matter in practice.
- Strengths:
  - Strong layer split: the envelope is just pCID dispatch plus opaque payload bytes.
  - Deterministic CBOR and explicit dispatch boundaries make bytes and rejection points locally auditable.
  - It avoids a built-in signature slot and explicitly does not claim a central registry or winning envelope.
- Weaknesses:
  - No concrete cross-variant group-session example demonstrates identical resolved payload semantics under different envelopes.
  - Unknown pCIDs are hard-rejected and not store/forwarded, which is brittle for migration and partial knowledge.
  - Authorship and integrity are pushed into wrapper protocols, so independence depends on extra specs not shown here.
- Risks:
  - Envelope choice may still leak into semantics through required signature-wrapper protocols and their validation rules.
  - Hard reject on unknown handlers can fragment mixed-version or long-lived group sessions.
  - The specimen is still draft and unfrozen, so its long-term audit and migration story is incomplete.
- Open questions:
  - Can the same CAS-backed group-session payload be wrapped by two envelope variants without changing application-level meaning or required audit evidence?
  - Should unknown envelopes be archivable or relayable for later decoding instead of being dropped outright?
  - How are signature-wrapper pCIDs versioned so authorship evidence stays comparable across envelope variants?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `cas-object-type-binding-bakeoff-unknown-typed-object`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/cas-object-type-binding-bakeoff-unknown-typed-object/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=2 promisegrid_alignment=2 auditability=3 evolution_safety=1 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=15.00 normalized_0_100=43.00 confidence_0_1=0.80
- Rationale: Good as a clear unsafe-parsing boundary specimen, but weak for this scenario because the only explicit unknown-type rule is hard rejection of unknown envelope pCIDs; CAS unknown-codec store/advertise/forward behavior is left unspecified.
- Strengths:
  - Deterministic CBOR [pcid, payload] gives a simple, auditable wire image.
  - Payload stays opaque until a known handler claims it, which helps avoid unsafe parsing.
  - Unknown outer pCIDs fail explicitly at the envelope boundary with optional local diagnostics.
- Weaknesses:
  - The spec does not define CAS-level opaque storage, advertisement, or forwarding for unknown object codecs.
  - Unknown outer pCIDs must not be accepted, stored, or forwarded as valid grid-envelope messages.
  - Wrapper-selected signatures make upgrades brittle when peers do not recognize newer wrapper types.
- Risks:
  - High mixed-version interoperability risk from hard-rejecting unknown outer types.
  - New signature wrappers may become non-relayable even when inner content would otherwise be harmless.
  - Implementations may invent inconsistent out-of-band blob handling to compensate for the missing opaque-forwarding story.
- Open questions:
  - Can a payload protocol on top of this envelope supply safe opaque handling for unknown-CID objects?
  - Is there a sanctioned quarantine/raw-blob retention path for unknown outer pCIDs without validating them as messages?
  - What migration path lets older nodes relay newer envelope or signature-wrapper types over long-lived deployments?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/cas-object-type-binding-bakeoff-raw-chunk-versus-pointer-bytes/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=4 evolution_safety=2 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=22.00 normalized_0_100=63.00 confidence_0_1=0.78
- Rationale: The variant helps this scenario by carrying an in-band pCID next to the bytes, so Bob does not need a local filename to know which protocol should interpret payload. But raw-versus-pointer semantics are still pushed into the payload protocol rather than defined by the envelope, and unknown-pCID hard rejection makes 100-year evolution brittle.
- Strengths:
  - In-band pCID ties interpretation to transported bytes instead of local paths.
  - Minimal deterministic CBOR [pCID, payload] gives a crisp and auditable layer boundary.
  - Unknown-pCID rejection prevents silent misparse.
- Weaknesses:
  - The spec does not itself define how raw payload bytes differ from pointer-object bytes inside payload.
  - The specimen is still draft and its own stable envelope identification is not yet frozen or minted.
  - Hard reject makes rollout of new pCIDs less graceful.
- Risks:
  - Older peers may drop envelopes that use future pCIDs instead of preserving them as opaque evidence.
  - Long-term audit depends on durable preservation of pCID-to-spec/handler material.
  - Bare CAS objects that bypass the envelope can still suffer raw-versus-pointer ambiguity.
- Open questions:
  - What CID codec or outer framing tells a CAS fetcher that stored bytes are this grid envelope?
  - Should unknown-pCID envelopes be storable or relayable as opaque bytes even when not accepted as valid envelopes?
  - Which payload protocol concretely binds raw content versus pointer-object semantics for this bakeoff?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `chunk-feed-replication-carrier-independence`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/chunk-feed-replication-carrier-independence/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=4 promisegrid_alignment=3 auditability=4 evolution_safety=2 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=22.00 normalized_0_100=68.00 confidence_0_1=0.76
- Rationale: Strong on carrier/layer separation: the minimal deterministic CBOR envelope can ride over many substrates, and the spec explicitly pushes ordering, forwarding, references, and auth wrappers out of the carrier. The main weakness is that carrier-independent chunk replication is only implied, not specified, and unknown-pCID hard reject at the envelope layer makes sparse-knowledge relays and long-term evolution brittle.
- Strengths:
  - Deterministic CBOR [pcid, payload] bytes are easy to move unchanged across diverse carriers.
  - Ordering, forwarding, external references, and auth are clearly delegated to payload or wrapper protocols, giving strong layer-boundary clarity.
  - The design does not require a central pCID registry.
- Weaknesses:
  - It does not define a concrete feed/chunk protocol that proves end-to-end carrier-independent replication.
  - Unknown-pCID hard reject stops generic grid-envelope peers from accepting, storing, or forwarding unfamiliar messages as valid.
  - Peer-local replication and accounting evidence is underspecified beyond local diagnostics.
- Risks:
  - Handler-preload requirements could fragment replication paths and weaken sparse-knowledge resilience.
  - Older peers may become dead ends instead of dumb carriers for opaque chunks.
  - Implementations may disagree on when a hop is merely transporting bytes versus acting as a hard-rejecting envelope receiver.
- Open questions:
  - Can unknown envelopes be archived or tunneled opaquely without violating the hard-reject rule?
  - What concrete feed/chunk and signature-wrapper stack demonstrates identical semantics over UDP, git, libp2p, IPFS, and ATPROTO-adjacent carriers?
  - What peer-local accounting artifact proves replication success or failure independent of carrier?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `chunk-feed-replication-corrupt-chunk`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/chunk-feed-replication-corrupt-chunk/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=3 evolution_safety=2 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=23.00 normalized_0_100=58.00 confidence_0_1=0.82
- Rationale: Deterministic CBOR, explicit pCID dispatch, and unknown-hard-reject make incompatibility visible, but for corrupt-chunk replication this specimen is mostly a clean transport boundary: CID/hash verification, retry, and peer-local accounting live in payload or wrapper protocols, not in the envelope.
- Strengths:
  - Deterministic CBOR gives stable bytes for hashing, signing, and later audit.
  - The [pcid, payload] split keeps CAS/body verification in the payload protocol instead of hiding it in the envelope.
  - Unknown pCIDs fail locally without requiring central enforcement or a central registry.
- Weaknesses:
  - The specimen does not define feed replication semantics, corrupt-chunk retry policy, or peer-local accounting records.
  - Base envelopes have no fixed signature slot, so sender attribution depends on extra sig-wrapper-pcid protocols.
  - Hard rejection of unknown pCIDs reduces graceful interop with newer protocols.
- Risks:
  - Older peers may drop or refuse to store future envelopes that they could otherwise preserve as evidence.
  - Different payload protocols may report corrupt chunks inconsistently because the envelope provides no common accounting shape.
  - Missing or inconsistent signature wrappers can weaken blame assignment for Mallory.
- Open questions:
  - Which payload protocol binds CID C to chunk bytes and records a local mismatch verdict?
  - What exact local accounting record is written when a peer sends corrupt bytes?
  - How should peers retry or switch sources after rejection without central coordination?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `chunk-feed-replication-duplicate-advertisement`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/chunk-feed-replication-duplicate-advertisement/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=3 evolution_safety=1 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=17.00 normalized_0_100=49.00 confidence_0_1=0.79
- Rationale: Partial fit: the variant can carry two separately signed advertisements for the same chunk and keeps reputation/accounting out of the feed itself, but it does not specify duplicate-offer semantics, Bob's peer-choice rules, or local promise-accounting records. The main liability for this scenario is unknown-hard-reject, which makes mixed-version wrapper or payload handling brittle.
- Strengths:
  - Minimal envelope keeps transport separate from reputation or promise-accounting policy.
  - Deterministic CBOR bytes and explicit pCID dispatch support local auditing.
  - Wrapper signatures can represent distinct Alice and Carol offers over the same chunk payload.
- Weaknesses:
  - No built-in chunk advertisement, deduplication, or peer-selection semantics.
  - Promise-accounting evidence is deferred to unspecified payload or wrapper protocols.
  - Unknown pCID handling is fail-closed, so future wrapper formats are dropped rather than tolerated.
- Risks:
  - Protocol evolution can break replication when Bob lacks a wrapper or payload handler.
  - Different wrapper conventions may make equivalent duplicate offers hard to compare.
  - Dropping unknown advertisements may hide usable chunk availability during migration.
- Open questions:
  - What payload protocol defines chunk offer identity, idempotence, and duplicate handling?
  - How does Bob record and compare local promise-accounting evidence for Alice versus Carol?
  - Can old peers preserve unknown advertisements for later reprocessing without violating hard-reject?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `conditional-release-geofencing-opaque-lower-layer-carriage`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/conditional-release-geofencing-opaque-lower-layer-carriage/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=2 evolution_safety=2 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=25.00 normalized_0_100=63.00 confidence_0_1=0.75
- Rationale: Strong on layer separation: the [pcid,payload] CBOR envelope preserves opaque bytes and pushes signatures/evidence into wrapper protocols, so lower layers need not parse condition semantics. Unknown-pCID hard reject also reduces accidental misinterpretation. But the key conditional-release proof and promise-accounting rules are mostly deferred to higher-layer protocols, and hard reject limits generic carriage when the outer pCID is unknown.
- Strengths:
  - Clear opaque boundary between envelope dispatch and payload semantics.
  - Deterministic CBOR plus wrapper-based signature/evidence protocols support stable layered verification.
  - Unknown-pCID hard reject makes unsupported cases fail visibly instead of being accepted ambiguously.
- Weaknesses:
  - Condition-object verification and promise-accounting rules are not defined in this specimen.
  - A node that lacks the outer pCID cannot store or forward the message as a valid envelope, limiting generic opaque carriage.
  - The base envelope has no mandatory signature, hop-evidence, or audit fields, and the spec is still draft.
- Risks:
  - Wrapper or payload protocol churn can strand older carriers because unknown top-level pCIDs are hard rejected.
  - Implementations may disagree on when opaque stored content is safe relative to local promises.
  - Audit trails may be uneven if nested signature conventions differ across payload protocols.
- Open questions:
  - What concrete wrapper binds encrypted content, condition reference, and releaser identity for this scenario?
  - Can relays or archive nodes retain unknown-top-level envelopes as uninterpreted bytes without counting them as accepted messages?
  - How does long-term pCID migration work without weakening the hard-reject safety boundary?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `conditional-release-geofencing-replay-outside-conditions`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/conditional-release-geofencing-replay-outside-conditions/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=3 evolution_safety=2 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=22.00 normalized_0_100=55.00 confidence_0_1=0.80
- Rationale: Deterministic CBOR, exact-byte dispatch, and unknown-hard-reject make the envelope auditable and keep unsafe unknown semantics from being silently accepted, but replay-outside-conditions is still pushed to wrapper/payload protocols and local state. This draft does not itself bind content to audience, geography, epoch, or session context.
- Strengths:
  - Unknown-hard-reject prevents permissive fallback when required condition or signature wrappers are not understood.
  - Deterministic CBOR bytes support exact local audit of what was replayed.
  - The minimal pCID envelope avoids central-registry dependence and keeps layer responsibility explicit.
- Weaknesses:
  - No built-in audience, location, expiry, nonce, or session binding for replay defense.
  - No fixed signature/authorship slot, so anti-replay evidence is not standardized at the envelope layer.
  - Draft/unfrozen semantics plus hard reject make long-lived evolution and mixed-version interoperability brittle.
- Risks:
  - A known valid old reference can still be replayed outside intended conditions unless higher-layer protocols and local state reject it.
  - Different deployments may encode conditional-release wrappers differently, reducing interoperability and audit comparability.
- Open questions:
  - Which pCID-defined wrapper or payload protocol must carry recipient, geography, epoch, and session constraints?
  - How should peers persist replay-detection state locally over long time spans?
  - Should conditional-release wrappers be mandatory-to-understand and frozen so omission or stripping is always unsafe?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `transport-family-bakeoff-per-hop-authorization-failure`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/transport-family-bakeoff-per-hop-authorization-failure/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=2 promisegrid_alignment=2 auditability=2 evolution_safety=1 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=14.00 normalized_0_100=40.00 confidence_0_1=0.85
- Rationale: This variant only partially fits the scenario: it cleanly separates envelope dispatch from per-hop authorization semantics, but it leaves Bob's refusal behavior to wrapper or payload protocols and hard-rejects unknown pCIDs, which is brittle for evolving authorization schemes.
- Strengths:
  - Clear boundary between envelope dispatch and higher-layer hop authorization logic.
  - Deterministic CBOR [pcid, payload] encoding is straightforward to implement and byte-audit.
  - Does not assume a central pCID registry and allows local rejection diagnostics.
- Weaknesses:
  - It does not define whether an authorization failure breaks the ring, skips a hop, records refusal, or reconfigures membership.
  - No fixed base-envelope signature or hop-evidence slot; auditors must understand wrapper-specific rules.
  - Unknown authorization or refusal wrapper pCIDs are hard-rejected rather than handled gracefully.
- Risks:
  - Protocol evolution can turn new authorization wrappers into transport breakage.
  - Refusal evidence may remain only local diagnostics instead of a portable cross-peer record.
  - Wrapper recursion can make long-term hop audit harder for humans and tools.
- Open questions:
  - What concrete wrapper or payload protocol records Bob's refusal to forward to Carol?
  - Should rejected unknown-pCID refusal evidence still be stored in a standard local audit form?
  - How would ring membership reconfiguration be represented without hidden global state?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `udp-feed-v0-conformance-malformed-datagram`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/udp-feed-v0-conformance-malformed-datagram/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=4 evolution_safety=2 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=20.00 normalized_0_100=57.00 confidence_0_1=0.74
- Rationale: This enc-cbor / unknown-hard-reject / sig-wrapper-pcid specimen is a partial fit for malformed-datagram UDP-feed pressure. It is strong on keeping bytes and dispatch identity separate from higher-layer meaning, but it does not actually define the UDP binding and leaves malformed-CBOR rejection behavior mostly implicit. Hard rejection of unknown pCIDs keeps failure visible but makes evolution less graceful.
- Strengths:
  - Clear envelope/payload boundary: [pcid, payload] with opaque payload bytes.
  - Deterministic CBOR canonical bytes support reproducible local audit evidence.
  - No central registry is assumed; dispatch failure is explicit at the pCID boundary.
- Weaknesses:
  - No UDP-feed handoff is specified, so the scenario match is indirect.
  - Malformed byte / invalid-CBOR handling is not normatively detailed.
  - unknown-hard-reject reduces mixed-version interoperability.
- Risks:
  - Implementations may diverge on how malformed datagrams are dropped, logged, or retained as evidence.
  - Future protocol rollout can fragment when older peers hard-reject new pCIDs.
- Open questions:
  - What exact UDP binding guarantees raw datagram bytes reach the envelope layer unchanged?
  - Should malformed or unknown envelopes be retained locally for audit even though they are invalid?
  - Can relays preserve unrecognized traffic without treating it as a valid grid-envelope message?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid` x `udp-feed-v0-conformance-session-layer-composition`

- Result path: `results/SIM-dorut-grid-envelope-enc-cbor-unknown-hard-reject-sig-wrapper-pcid/udp-feed-v0-conformance-session-layer-composition/openai-gpt-5.4-xhigh/20260522-012332.json`
- Scores: scenario_fit=4 promisegrid_alignment=3 auditability=3 evolution_safety=2 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=27.00 normalized_0_100=67.50 confidence_0_1=0.76
- Rationale: Strong on composition boundaries: the envelope stays transport-neutral and assigns ordering, relay, and auth semantics to pCID-selected payload or wrapper protocols. That fits the session-layer composition pressure, but unknown-hard-reject and wrapper-only signatures make long-term evolution brittle.
- Strengths:
  - Clear separation between envelope dispatch and higher-layer/session semantics.
  - Deterministic CBOR exact-byte rules support reproducible hashing, signing, and audit trails.
  - pCID-selected wrappers avoid baking one signature or relay scheme into the base envelope.
- Weaknesses:
  - The docs do not define the actual UDP-feed v0 API surface, so the scenario is answered indirectly.
  - Unknown-pCID hard rejection weakens forward compatibility and sparse-knowledge interoperability.
  - Authorship and integrity depend on wrapper conventions rather than a shared base-envelope field.
- Risks:
  - Version or handler skew can cause legitimate future session messages to be dropped.
  - Long-term archives may keep bytes but lose interpretability if wrapper or pCID knowledge decays.
- Open questions:
  - What minimum UDP-feed metadata must be exposed with datagram bytes for this envelope to support a session layer cleanly?
  - Can peers quarantine or archive unknown-pCID envelopes safely without violating the hard-reject rule?
  - Which signature-wrapper profile should a minimal group/session layer assume for interoperable authentication?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

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

## Required JSON Shape

{"child_id":"SIM-dubil-child-descriptive-design-slug","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
