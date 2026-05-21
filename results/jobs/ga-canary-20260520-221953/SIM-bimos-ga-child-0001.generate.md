# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260520-221953`
- Child ID: `SIM-bimos-ga-child-0001`
- Child path: `simulations/SIM-bimos-ga-child-0001/`
- Operation: `breed`
- Parent IDs: `SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes, SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload`

## Scenario Sample

- `iot-fleet-maintenance` at `scenarios/iot-fleet-maintenance/iot-fleet-maintenance.md`
- `bgp-routing` at `scenarios/bgp-routing/bgp-routing.md`
- `cas-object-type-binding-bakeoff-application-object-family` at `scenarios/cas-object-type-binding-bakeoff-application-object-family/cas-object-type-binding-bakeoff-application-object-family.md`

## Scenario Pressure

### `scenarios/iot-fleet-maintenance/iot-fleet-maintenance.md`

```markdown
# IoT Fleet Maintenance

## Scenario ID

iot-fleet-maintenance

## Source / Provenance

- Source type: application seed
- Source path: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`
- Source row/title: Seed application catalog entry `iot-fleet-maintenance`
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-midif`; `TODO-dadub`

## Purpose

Exercise PromiseGrid design candidates against iot-fleet-maintenance application
pressure: Device identity, maintenance history, firmware updates, telemetry, and access
control.

## Setup

Alice depends on an outcome in the IoT Fleet Maintenance domain. Bob makes promises
about device identity, maintenance history, firmware updates, telemetry, and access
control. Carol needs enough evidence to rely on Bob's promise without having complete
global state, and Mallory may exploit stale, missing, or ambiguous evidence.

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

### `scenarios/cas-object-type-binding-bakeoff-application-object-family/cas-object-type-binding-bakeoff-application-object-family.md`

```markdown
# Application object family

## Scenario ID

cas-object-type-binding-bakeoff-application-object-family

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-kohad-cas-object-type-binding-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-kohad-cas-object-type-binding-bakeoff/`
- Source row/title: Application object family
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-kohad-cas-object-type-binding-bakeoff/`.

## Setup

Ellen proposes a future application-level CAS object distinct from raw chunks, Merkle nodes, and pointer objects.

## Stimulus

Run the candidate simulation against this source test: Whether the chosen binding model leaves room for new object families without reinterpreting old bytes.

## Expected Pressure

The first type-binding rule should be extensible without changing old CIDs.
```

## Parent Simulation Documents

### `simulations/SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes/README.md`

```markdown
# SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes: Grid-envelope variant

This simulation is a standalone positional grid-envelope specimen. It tests the
combination `enc-dag-cbor`, `unknown-opaque`, and `sig-mandatory-opaque-bytes` without claiming
that this combination is the canonical PromiseGrid wire format. Source: `DI-fanah`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
```

### `simulations/SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes/QUESTION.md`

```markdown
# Question

Does a positional grid envelope using `enc-dag-cbor`, `unknown-opaque`, and
`sig-mandatory-opaque-bytes` satisfy the wire-lab harness scenarios better than the sibling
variants? Source: `DI-fanah`.
```

### `simulations/SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries
yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen created by `DI-fanah`.
```

### `simulations/SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid Envelope Variant Spec (DRAFT)

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `enc-dag-cbor` / `unknown-opaque` / `sig-mandatory-opaque-bytes`.
> Source: `DI-fanah`.

## Purpose

This spec defines one full positional grid-envelope candidate for wire-lab
comparison. It is a specimen inside `SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes`, not a harness rule and not the
canonical PromiseGrid wire format.

## Positional Envelope Shape

The envelope shape for this variant is:

```text
[pcid, payload, signature]
```

Slots are interpreted positionally:

- `pcid` identifies the protocol/spec/handler that interprets `payload`.
- `payload` is opaque bytes until interpreted by the handler named by `pcid`.
- `signature` is mandatory opaque bytes over the canonical unsigned prefix.

A `payload` may itself be the canonical bytes of another grid envelope when the
protocol named by `pcid` specifies recursive nesting. The outer grid envelope
does not prescribe the payload's internal organization beyond the bytes boundary.

## Encoding

This variant encodes the envelope as DAG-CBOR-compatible positional arrays. `pcid` and `sig_pcid`, when present, are DAG-CBOR Link values; `payload`, `signature`, and `sig_payload` are byte strings. The envelope remains positional: no map/object envelope fields are introduced. The canonical bytes for signing and hashing are the DAG-CBOR bytes of the exact positional array under this spec.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid`, it may store and forward the exact envelope bytes as opaque content, but interpretation fails with an explicit unsupported-pCID result. A receiver MUST NOT parse `payload` speculatively without the handler named by `pcid`.

## Signature and Authorship Policy

The third positional slot, `signature`, is mandatory opaque bytes. The bytes cover the canonical unsigned prefix `[pcid, payload]` under this variant's encoding. The envelope layer enforces presence and byte-string shape; signature algorithm, signer identity, and verification semantics are determined by the protocol ecosystem being tested with this variant.

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

### `simulations/SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/README.md`

```markdown
# SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload: Grid-envelope variant

This simulation is a standalone positional grid-envelope specimen. It tests the
combination `enc-dag-cbor`, `unknown-hard-reject`, and `sig-mandatory-sig-pcid-payload` without claiming
that this combination is the canonical PromiseGrid wire format. Source: `DI-fanah`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
```

### `simulations/SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/QUESTION.md`

```markdown
# Question

Does a positional grid envelope using `enc-dag-cbor`, `unknown-hard-reject`, and
`sig-mandatory-sig-pcid-payload` satisfy the wire-lab harness scenarios better than the sibling
variants? Source: `DI-fanah`.
```

### `simulations/SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries
yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen created by `DI-fanah`.
```

### `simulations/SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid Envelope Variant Spec (DRAFT)

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `enc-dag-cbor` / `unknown-hard-reject` / `sig-mandatory-sig-pcid-payload`.
> Source: `DI-fanah`.

## Purpose

This spec defines one full positional grid-envelope candidate for wire-lab
comparison. It is a specimen inside `SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload`, not a harness rule and not the
canonical PromiseGrid wire format.

## Positional Envelope Shape

The envelope shape for this variant is:

```text
[pcid, payload, sig_pcid, sig_payload]
```

Slots are interpreted positionally:

- `pcid` identifies the protocol/spec/handler that interprets `payload`.
- `payload` is opaque bytes until interpreted by the handler named by `pcid`.
- `sig_pcid` identifies the signature or proof protocol.
- `sig_payload` is opaque bytes for the handler named by `sig_pcid`.

A `payload` may itself be the canonical bytes of another grid envelope when the
protocol named by `pcid` specifies recursive nesting. The outer grid envelope
does not prescribe the payload's internal organization beyond the bytes boundary.

## Encoding

This variant encodes the envelope as DAG-CBOR-compatible positional arrays. `pcid` and `sig_pcid`, when present, are DAG-CBOR Link values; `payload`, `signature`, and `sig_payload` are byte strings. The envelope remains positional: no map/object envelope fields are introduced. The canonical bytes for signing and hashing are the DAG-CBOR bytes of the exact positional array under this spec.

## Unknown pCID Policy

If a receiver lacks a handler for `pcid`, the envelope is rejected at the envelope layer. The receiver may keep local diagnostics, but it MUST NOT accept, store, or forward the message as a valid grid-envelope message under this variant.

## Signature and Authorship Policy

The third and fourth positional slots are mandatory. `sig_pcid` identifies the signature or proof protocol; `sig_payload` is opaque bytes interpreted by that signature protocol. The signature payload covers the canonical unsigned prefix `[pcid, payload]` under this variant's encoding unless `sig_pcid` publishes stricter rules.

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

### `SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` x `iot-fleet-maintenance`

- Result path: `results/SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/iot-fleet-maintenance/openai-gpt-5.4-xhigh/20260520-221953.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=3 evolution_safety=2 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=22.00 normalized_0_100=55.00 confidence_0_1=0.75
- Rationale: Strong as a strict signed wire envelope with clear dispatch boundaries, but only a partial fit for IoT fleet maintenance because device identity, maintenance history, freshness, and access-control logic live above the envelope. The main downside is upgrade brittleness from unknown pCID hard reject in long-lived mixed fleets.
- Strengths:
  - Clear separation of envelope dispatch, payload semantics, and signature semantics.
  - Mandatory signatures over canonical DAG-CBOR bytes improve integrity checks and later audit.
  - Does not require a central registry; peers can decide from local handler availability.
- Weaknesses:
  - Does not itself define the application objects, feeds, and promise accounting the scenario needs.
  - Unknown pCID hard reject reduces graceful interoperability across mixed protocol versions.
  - No built-in semantics for stale, delayed, disputed, or partially available maintenance evidence.
- Risks:
  - Older gateways or devices may drop future update or telemetry envelopes instead of safely relaying them.
  - Useful maintenance evidence may be lost or sidelined when receivers cannot interpret a new pCID.
  - Mandatory signature slots without a defined rotation and revocation story can create fleet lockout pressure.
- Open questions:
  - Which payload protocols and feeds carry device identity, maintenance history, firmware approval, telemetry, and access control?
  - Can unknown envelopes be archived as opaque evidence while keeping the hard-reject safety goal?
  - How are pCID and signature-protocol migrations handled across decades of intermittent, mixed-version fleets?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` x `bgp-routing`

- Result path: `results/SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/bgp-routing/openai-gpt-5.4-xhigh/20260520-221953.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=4 evolution_safety=1 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=17.00 normalized_0_100=49.00 confidence_0_1=0.78
- Rationale: Useful as a strict, signed wrapper for route-like claims, but not a routing design. Canonical DAG-CBOR bytes, explicit pCID dispatch, and mandatory signatures help local verification; hijack/leak policy, freshness, and mixed-version rollout remain mostly outside the specimen or become brittle under unknown-hard-reject.
- Strengths:
  - Canonical DAG-CBOR bytes plus mandatory signature slots give stable local evidence for a route-like claim.
  - Signing covers both pCID and payload, reducing type-confusion and unsigned-wrapper ambiguity.
  - Unknown payload protocols fail visibly at the dispatch boundary instead of being silently reinterpreted.
  - ... 1 more
- Weaknesses:
  - No route/path/withdrawal/freshness semantics are defined here, so hijack, leak, and stale-route pressure is mostly pushed to the payload protocol.
  - unknown-hard-reject blocks graceful mixed-version interop and may discard unfamiliar routing evidence.
  - Long-term crypto/version migration behavior is under-specified.
- Risks:
  - Mixed deployments can black-hole newer route messages instead of preserving them for later validation.
  - A valid signature proves envelope authorship, not that a route export was authorized or still current.
  - If pCID or sig_pcid specs are not durably recoverable, 100-year auditability of old route evidence degrades.
- Open questions:
  - What is the exact receiver policy when sig_pcid is unknown, deprecated, or temporarily unavailable?
  - Can peers archive and later relay raw rejected envelopes as evidence without violating the hard-reject rule?
  - Which payload-level objects carry route freshness, withdrawals, and leak-policy evidence?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload` x `cas-object-type-binding-bakeoff-application-object-family`

- Result path: `results/SIM-sivus-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-sig-pcid-payload/cas-object-type-binding-bakeoff-application-object-family/openai-gpt-5.4-xhigh/20260520-221953.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=4 evolution_safety=3 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=22.00 normalized_0_100=63.00 confidence_0_1=0.74
- Rationale: Partial but useful fit: the in-band pCID gives a self-describing first binding rule that can name new application object families without reinterpreting prior envelope bytes, and the envelope/payload/signature split is clear and auditable. But this specimen stops at envelope dispatch rather than defining CAS object-family CID semantics directly, and unknown-hard-reject makes future-family rollout brittle.
- Strengths:
  - In-band pCID leaves room for new payload and application object families.
  - Canonical DAG-CBOR bytes plus mandatory signature slots support later audit of what was bound and signed.
  - Envelope, payload protocol, and signature protocol boundaries are explicit.
- Weaknesses:
  - The spec does not directly define CAS object-family rules or stable CID semantics for those families.
  - Unknown pCID values are hard-rejected, which weakens mixed-version evolution.
  - pCID publication, discovery, and long-term availability are not yet frozen.
- Risks:
  - Older peers may fracture off from future application object families because they must reject unknown pCIDs.
  - Long-term interpretability depends on durable recovery of the protocol specs named by pCID.
- Open questions:
  - Should unknown future application-family envelopes be archivable or relayable even when not locally interpretable?
  - Is the stable CAS type-binding meant to live at the outer envelope pCID, inside the payload, or both?
  - How are pCID lookup and long-term spec retrieval handled without creating a de facto authority?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes` x `iot-fleet-maintenance`

- Result path: `results/SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes/iot-fleet-maintenance/openai-gpt-5.4-xhigh/20260520-221953.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=2 evolution_safety=3 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=3 risk_penalty=3
- Fitness: raw=22.00 normalized_0_100=55.00 confidence_0_1=0.77
- Rationale: Strong at carrying exact signed bytes across sparse peers with explicit dispatch and failure boundaries, but too much IoT-maintenance meaning—identity, authorization, freshness, and maintenance-history semantics—lives outside the envelope for this specimen to solve the scenario on its own.
- Strengths:
  - Clear layer boundary: positional [pcid, payload, signature], opaque payload bytes, and no speculative parsing without the named handler.
  - Unknown-pCID store/forward behavior supports sparse knowledge and heterogeneous long-lived fleets.
  - Mandatory signature presence over canonical DAG-CBOR bytes is a sensible default for security-sensitive update and telemetry artifacts.
- Weaknesses:
  - The envelope does not define device identity, maintenance-log structure, firmware policy, or access-control semantics needed by the scenario.
  - Signature semantics are external and opaque, so peers cannot infer signer authority, algorithm agility, or revocation rules from the envelope alone.
  - The specimen is still a draft and references undefined sig_pcid/sig_payload fields, increasing ambiguity and implementation-drift risk.
- Risks:
  - Mandatory-but-opaque signature bytes can create false confidence if payload ecosystems do not pin signer and verification policy precisely.
  - Stale or replayed signed envelopes may look acceptable unless payload protocols add freshness and policy checks.
  - Unknown-envelope forwarding can preserve malicious or permanently uninterpretable artifacts until higher layers reject them.
- Open questions:
  - How are device identity, key rotation, revocation, and update authorization represented and audited locally?
  - What freshness or anti-replay mechanism protects telemetry, maintenance events, and firmware offers under partitioned conditions?
  - What durable recovery path lets a maintainer resolve a stored pCID and its signature rules decades later without central authority?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes` x `bgp-routing`

- Result path: `results/SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes/bgp-routing/openai-gpt-5.4-xhigh/20260520-221953.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=2 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=24.00 normalized_0_100=60.00 confidence_0_1=0.75
- Rationale: Strong as an envelope substrate, limited as a direct bgp-routing answer: it preserves signature-bearing opaque route objects across sparse peers and makes interpretation boundaries explicit, but routing trust, freshness, and leak/hijack semantics sit above the envelope.
- Strengths:
  - Clear pCID dispatch boundary and explicit ban on speculative parsing preserve layer separation.
  - Canonical DAG-CBOR bytes plus a mandatory signature slot give stable byte-level evidence.
  - Unknown-pCID store/forward supports sparse peers and future protocol variants without a central registry.
- Weaknesses:
  - It does not define route announcements, withdrawals, freshness, identity claims, or promise-accounting records for bgp-routing.
  - Signer identity and verification rules are outside the envelope, so Carol cannot rely on envelope data alone.
  - The draft mentions sig_pcid/sig_payload despite a 3-slot envelope, leaving signature-layer details unsettled.
- Risks:
  - Hijack or leak evidence can stay ambiguous unless the payload protocol carries explicit path and policy attestations.
  - Stale or replayed route objects may circulate because freshness and expiry are not envelope-level concepts here.
  - Long-term audit depends on durable access to the pCID-linked handler/spec; otherwise payload and signature remain opaque.
- Open questions:
  - How does a verifier learn the signature algorithm and signer/key binding when signature is only opaque bytes?
  - How are route freshness, withdrawals, and replay protection represented for BGP-like payloads?
  - Are nested wrapper envelopes enough for multi-hop AS-path evidence, or is native multi-signature support needed?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes` x `cas-object-type-binding-bakeoff-application-object-family`

- Result path: `results/SIM-ruzil-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-opaque-bytes/cas-object-type-binding-bakeoff-application-object-family/openai-gpt-5.4-xhigh/20260520-221953.json`
- Scores: scenario_fit=4 promisegrid_alignment=4 auditability=3 evolution_safety=4 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=4 risk_penalty=2
- Fitness: raw=31.00 normalized_0_100=78.00 confidence_0_1=0.74
- Rationale: This variant gives a plausible extensible first type-binding rule for future application object families: `pcid` selects the interpreter, `payload` stays opaque, and unknown families can be preserved without reinterpreting old bytes. Its main limitation is that it stops at envelope dispatch and leaves the actual CAS object-family, CID, and long-term pCID/signature governance story to higher layers.
- Strengths:
  - Explicit `pcid` dispatch cleanly separates old and new object families.
  - Unknown-`pcid` handling preserves exact bytes and fails visibly instead of speculatively parsing.
  - Deterministic DAG-CBOR positional encoding keeps the boundary simple and implementable.
- Weaknesses:
  - It is a wire-envelope specimen, not a complete CAS object-family binding model.
  - Opaque payload and signature fields reduce semantic auditability without the referenced handler/spec.
  - The draft still leaves pCID and signature lifecycle details unresolved, so freeze-ready confidence is limited.
- Risks:
  - Long-term meaning depends on durable decentralized publication and recovery of pCID-linked specs.
  - Mandatory opaque signature bytes may be an awkward fit for some future application object families.
  - Different ecosystems could drift on verification or family conventions while sharing the same outer shape.
- Open questions:
  - Should application-family identity live entirely in outer `pcid`, inside payload rules, or in both places?
  - What bytes are intended to be the stable CAS identity for comparison: outer envelope, inner payload, or a higher-level object?
  - How should unsigned, multi-signed, or delegated-authorship application objects fit this mandatory-signature variant?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

## Required JSON Shape

{"child_id":"SIM-bimos-ga-child-0001","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
