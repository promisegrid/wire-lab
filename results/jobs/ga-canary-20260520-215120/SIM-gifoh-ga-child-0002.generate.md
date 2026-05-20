# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260520-215120`
- Child ID: `SIM-gifoh-ga-child-0002`
- Child path: `simulations/SIM-gifoh-ga-child-0002/`
- Operation: `breed`
- Parent IDs: `SIM-zazit-chunk-feed-replication, SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid`

## Scenario Sample

- `conditional-release-geofencing-geofenced-group-dispatch` at `scenarios/conditional-release-geofencing-geofenced-group-dispatch/conditional-release-geofencing-geofenced-group-dispatch.md`
- `insurance-claims` at `scenarios/insurance-claims/insurance-claims.md`
- `bgp-class-routing-app-route-hijack` at `scenarios/bgp-class-routing-app-route-hijack/bgp-class-routing-app-route-hijack.md`

## Scenario Pressure

### `scenarios/conditional-release-geofencing-geofenced-group-dispatch/conditional-release-geofencing-geofenced-group-dispatch.md`

```markdown
# Geofenced group dispatch

## Scenario ID

conditional-release-geofencing-geofenced-group-dispatch

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zarud-conditional-release-geofencing/SCENARIOS.md`
- Source simulation: `SIM-zarud-conditional-release-geofencing/`
- Source row/title: Geofenced group dispatch
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zarud-conditional-release-geofencing/`.

## Setup

Alice permits content only for group members inside a stated region. Carol is a member but is outside the allowed region.

## Stimulus

Run the candidate simulation against this source test: Whether geofence checks are membership checks, per-message dispatch checks, fetch-policy checks, or storage constraints.

## Expected Pressure

The owner layer must explain both refusal and auditability without assuming a central location oracle.
```

### `scenarios/insurance-claims/insurance-claims.md`

```markdown
# Insurance Claims

## Scenario ID

insurance-claims

## Source / Provenance

- Source type: application seed
- Source path: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`
- Source row/title: Seed application catalog entry `insurance-claims`
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-midif`; `TODO-dadub`

## Purpose

Exercise PromiseGrid design candidates against insurance-claims application pressure:
Claim evidence, adjuster authority, fraud pressure, payments, and appeal promises.

## Setup

Alice depends on an outcome in the Insurance Claims domain. Bob makes promises about
claim evidence, adjuster authority, fraud pressure, payments, and appeal promises. Carol
needs enough evidence to rely on Bob's promise without having complete global state, and
Mallory may exploit stale, missing, or ambiguous evidence.

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

### `scenarios/bgp-class-routing-app-route-hijack/bgp-class-routing-app-route-hijack.md`

```markdown
# Route hijack

## Scenario ID

bgp-class-routing-app-route-hijack

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-punaz-bgp-class-routing-app/SCENARIOS.md`
- Source simulation: `SIM-punaz-bgp-class-routing-app/`
- Source row/title: Route hijack
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-punaz-bgp-class-routing-app/`.

## Setup

Mallory advertises a short attractive path to Carol but cannot deliver traffic or chunks.

## Stimulus

Run the candidate simulation against this source test: How peers detect failed promises and locally downgrade future route choices.

## Expected Pressure

PromiseGrid routing apps need hijack costs that are local and relationship-specific.
```

## Parent Simulation Documents

### `simulations/SIM-zazit-chunk-feed-replication/README.md`

```markdown
# SIM-zazit: Chunk feed replication

This simulation explores the turn-177 inversion that L5 feeds replicate CAS
chunks between sites while L6 CAS stores and resolves those chunks for higher
layers. It treats feeds as meaning-oblivious replication mechanisms rather
than group-message transports. Source: `DI-navod`.

## Question

What does an L5 feed need to advertise, request, and replicate when its unit of
movement is a CAS chunk or CAS object rather than a group-session message?
Source: `DI-navod`.

## Turn 177 pressure

Turn 177 corrected the model from "CAS below feeds that carry messages" to
"feeds below CAS that move chunks." That has practical consequences:

- A feed should not need to understand group-session message semantics.
- Sparse sites should be able to advertise only the chunks they have or want.
- A single feed family may serve group messages, files, and future content that
  share the CAS layer.
- Pull/keep/advertise decisions need inputs from peer-local promise accounting
  records without making L5 itself a central accounting authority.

Turn 178 adds two more pressures to this simulation: every site should be sparse
by default, and some realistic experiments may put each simulated site or large
content corpus in a different repository rather than assuming all message bytes
live in wire-lab. Source: `DI-vaguf`.

This simulation owns the feed-side consequences of that inversion. The
`SIM-jomag` CAS object-model simulation owns object typing and chunking
parameters; this simulation asks how those objects move between sites. Source:
`DI-navod`.

## Decision axes

- **Advertisement unit:** individual chunk CID, Merkle root CID, pointer object
  CID, range/frontier summary, or a mix.
- **Request policy:** how a site decides which advertised chunks to pull under
  sparse-CAS assumptions.
- **Replication promise:** what a feed participant promises about timely
  storage, retransmission, integrity, and refusal.
- **Carrier fit:** how UDP, git, libp2p, IPFS, ATPROTO, or other substrates
  might carry the same feed role without becoming the protocol meaning.
- **Site topology:** whether a simulation site is a directory, a separate repo,
  a remote peer, or an imported fixture when testing sparse large-object flows.
- **Failure behavior:** duplicate advertisements, missing chunks, partial
  Merkle trees, corrupt chunks, and peer-specific refusal.

## Boundaries

This simulation is not the `udp-feed` spec and does not modify the existing
`SIM-ludaf-udp-feed` lineage. It is the turn-177 chunk-replication design-point
workspace that later feed specimens can adopt, reject, or compete against.
Source: `DI-navod`.
```

### `simulations/SIM-zazit-chunk-feed-replication/QUESTION.md`

```markdown
# Question

If L5 feeds move CAS chunks instead of group messages, what must a feed
advertise, request, replicate, and promise so sparse sites can converge without
assuming every site stores all CAS objects? Source: `DI-navod`.

Open decision points:

- Does a feed advertise leaves, roots, pointer objects, or compact frontiers?
- Which layer decides that an advertised chunk is worth pulling or retaining?
- What feed behavior remains substrate-neutral across UDP, git, libp2p, IPFS,
  ATPROTO, and future carriers?
- When do sparse sites need separate repos or external content corpora rather
  than directories inside one wire-lab checkout?
- How are corrupt, missing, duplicate, or refused chunks represented without
  leaking group-session semantics into L5?
```

### `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`

```markdown
# Chunk Feed Replication Scenarios

These scenarios make the turn-177 L5 feed inversion concrete: feeds move CAS
chunks or CAS objects between sparse sites, while L7 group semantics remain
above CAS. They are simulation inputs for TODO-kituj / TE-43 and TODO-pipus,
not a frozen feed wire format. Source: `DI-pator`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Sparse advertisement | Alice has a subset of chunks for a Merkle root; Bob has a different subset. | Whether the feed advertises leaves, roots, pointer objects, frontiers, or compact summaries without assuming full replication. | Feed specs must work when no site has all CAS objects. |
| Pull decision | Bob receives an advertisement for chunk C and has peer-local promise accounting records about Alice. | Which inputs decide whether Bob pulls, delays, refuses, or asks another peer. | The "decides" step needs an explicit cross-layer interface instead of a hand wave. |
| Multi-repo sparse site | Alice's site state, Bob's site state, and a large shared corpus live in separate repos or fixtures. | Whether feed promises and CAS object references remain meaningful when the harness orchestrates multiple storage roots. | Turn 178's multi-repo question should be explored without assuming one repo contains every site's content. |
| Partial Merkle fetch | Bob wants root R but only some children are locally available. | Whether the feed can request missing children without understanding group-session message semantics. | L5 should remain meaning-oblivious while still serving L6 CAS repair. |
| Corrupt chunk | Mallory advertises or sends bytes whose hash does not match CID C. | Whether rejection, accounting, and retry behavior are local enough to avoid central enforcement. | Feed behavior must compose with CAS hash verification and peer-local accounting records. |
| Duplicate advertisement | Alice and Carol both advertise chunk C. | Whether duplicate offers are harmless and how Bob chooses between peers. | Promise accounting can influence peer choice without making the feed a central reputation service. |
| Refusal or non-response | Alice refuses to send C or never answers. | Whether refusal is a normal observed outcome that can feed future local decisions. | The feed protocol needs space for refusal and timeout outcomes without treating every miss as corruption. |
| Carrier independence | The same chunk exchange is attempted over UDP, git, libp2p, IPFS, or ATPROTO-adjacent carriers. | Which semantics belong to the feed role and which are carrier mechanics. | The simulation should preserve turn-177's claim that feeds move chunks independent of substrate. |

## Expected Outputs

- A candidate list of L5 feed observations that can be recorded without group
  semantics leaking into the feed layer.
- A TODO-pipus migration constraint: successor group-session specimens should
  depend on chunk replication, not message-file transport.
- A TE-43 constraint: CAS object and chunking decisions must be usable by
  sparse feed replication.
```

### `simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/README.md`

```markdown
# SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid: Grid-envelope variant

This simulation is a standalone positional grid-envelope specimen. It tests the
combination `enc-dag-cbor`, `unknown-best-effort`, and `sig-wrapper-pcid` without claiming
that this combination is the canonical PromiseGrid wire format. Source: `DI-fanah`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.
```

### `simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/QUESTION.md`

```markdown
# Question

Does a positional grid envelope using `enc-dag-cbor`, `unknown-best-effort`, and
`sig-wrapper-pcid` satisfy the wire-lab harness scenarios better than the sibling
variants? Source: `DI-fanah`.
```

### `simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries
yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen created by `DI-fanah`.
```

### `simulations/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid Envelope Variant Spec (DRAFT)

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `enc-dag-cbor` / `unknown-best-effort` / `sig-wrapper-pcid`.
> Source: `DI-fanah`.

## Purpose

This spec defines one full positional grid-envelope candidate for wire-lab
comparison. It is a specimen inside `SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid`, not a harness rule and not the
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

If a receiver lacks a handler for `pcid`, it may expose `payload` bytes to generic tooling for inspection or salvage. Any such result MUST be marked unsupported and unverified; best-effort inspection does not count as interpretation under the missing `pcid` rules.

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

## Compact Fitness Evidence From This Run

### `SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid` x `conditional-release-geofencing-geofenced-group-dispatch`

- Result path: `results/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/conditional-release-geofencing-geofenced-group-dispatch/openai-gpt-5.4-xhigh/20260520-215120.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=2 evolution_safety=4 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=23.00 normalized_0_100=58.00 confidence_0_1=0.79
- Rationale: Strong as a minimal carrier and boundary marker, but mostly transport-only for this scenario: it can host a geofenced dispatch protocol without central registries, yet it does not define the membership, location, refusal, or audit evidence the scenario requires.
- Strengths:
  - Clear wire/application boundary: geofencing can live in higher-layer protocols without making the envelope depend on a central location oracle.
  - Wrapper-pCID signatures let higher layers attach signed authorization or refusal evidence without changing the base [pcid, payload] shape.
  - Unknown-pCID results must be marked unsupported and unverified, which helps auditors distinguish salvage from valid interpretation.
- Weaknesses:
  - No native structure for membership, region proofs, dispatch/fetch/storage policy, or peer-local refusal records.
  - Auditability for denial depends entirely on extra payload or wrapper protocols not defined in this specimen.
  - The draft, unfrozen pCID/spec state weakens current 100-year durability.
- Risks:
  - Unknown-best-effort payload exposure is uncomfortable for conditional release if tooling surfaces bytes before policy or confidentiality wrappers are validated.
  - Wrapper-only signatures can be stripped, ignored, or differently ordered, obscuring what policy actually governed release.
  - If pCID resolution becomes de facto centralized or poorly replicated, PromiseGrid's no-central-authority goal erodes.
- Open questions:
  - What higher-layer protocol carries peer-local proofs of membership and region without a central oracle?
  - How is a refusal recorded so auditors can tell denied dispatch from denied fetch or denied storage?
  - What guard stops generic tooling from exposing inner content before relevant wrappers are verified?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid` x `insurance-claims`

- Result path: `results/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/insurance-claims/openai-gpt-5.4-xhigh/20260520-215120.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=2 evolution_safety=3 layer_boundary_clarity=4 failure_handling=3 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=18.00 normalized_0_100=51.00 confidence_0_1=0.80
- Rationale: Useful as a minimal wire-layer carrier for contested claim artifacts: it gives explicit pCID dispatch, canonical DAG-CBOR bytes, and a visible downgrade path for unknown protocols. But it only partially meets insurance-claims pressure because adjuster authority, fraud evidence, payments, appeals, names, and promise accounting are left to unspecified payload and wrapper protocols. Wrapper-only signatures and best-effort unknown-pCID handling weaken direct auditability for high-stakes disputes.
- Strengths:
  - Clear [pcid, payload] boundary keeps transport semantics separate from claim-specific payload rules.
  - Unknown-pCID handling explicitly marks salvage as unsupported and unverified instead of silently treating it as valid evidence.
  - DAG-CBOR canonical bytes plus pCID-selected wrappers are straightforward to implement and do not add a central registry dependency at the envelope layer.
- Weaknesses:
  - Does not define insurance-claims objects or local accounting for evidence, adjuster authority, fraud review, payment, or appeals.
  - Base envelope has no fixed signature/authorship slot, so provenance inspection depends on extra wrapper conventions.
  - The specimen is still draft and the spec pCID is not yet minted, which hurts 100-year durability.
- Risks:
  - Operators or tools may overtrust best-effort inspection of unsupported payloads during contested claims.
  - Different signature-wrapper conventions could fragment provenance verification across peers.
  - Long-retention claim records may become hard to interpret if pCID/spec preservation is weak.
- Open questions:
  - How are adjuster authority, delegation, and revocation represented and audited locally?
  - How are claim evidence, payment decisions, and appeal records linked into peer-local promise accounting?
  - How should peers distinguish stale or malicious claim data from merely unsupported unknown-pCID payloads?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid` x `bgp-class-routing-app-route-hijack`

- Result path: `results/SIM-johum-grid-envelope-enc-dag-cbor-unknown-best-effort-sig-wrapper-pcid/bgp-class-routing-app-route-hijack/openai-gpt-5.4-xhigh/20260520-215120.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=3 evolution_safety=3 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=20.00 normalized_0_100=50.00 confidence_0_1=0.83
- Rationale: The specimen gives a clean, auditable wire envelope with explicit dispatch and unsupported/unverified handling, but it does not define the signed route claims, delivery-failure evidence, or peer-local downgrade rules that a route-hijack scenario needs. It is a possible substrate, not a sufficient routing answer.
- Strengths:
  - Clear separation between envelope, wrapper, and payload responsibilities.
  - Canonical DAG-CBOR bytes help preserve what was actually advertised for later audit.
  - Unknown-pCID best-effort inspection must be marked unsupported and unverified.
  - ... 1 more
- Weaknesses:
  - No route-advertisement, delivery-proof, or relationship-specific penalty semantics are defined.
  - No fixed signature/authorship slot; accountability depends on extra wrapper protocols.
  - Failure handling is clearer for dispatch and unknown-pCID cases than for undelivered traffic or chunks.
  - ... 1 more
- Risks:
  - Best-effort inspection of unknown payloads can create ambiguous incident evidence during hijack handling.
  - Wrapper-specific signature schemes may fragment route-authenticity checks across peers.
  - Without higher-layer promise accounting, Carol may not have a consistent local rule for downgrading Mallory.
- Open questions:
  - Which wrapper or payload protocol carries signed route advertisements and failed-delivery evidence?
  - How are local, relationship-specific hijack costs recorded and applied to future route choice?
  - What peer-local evidence distinguishes malicious hijack from transient outage or path churn?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-zazit-chunk-feed-replication` x `conditional-release-geofencing-geofenced-group-dispatch`

- Result path: `results/SIM-zazit-chunk-feed-replication/conditional-release-geofencing-geofenced-group-dispatch/openai-gpt-5.4-xhigh/20260520-215120.json`
- Scores: scenario_fit=1 promisegrid_alignment=4 auditability=1 evolution_safety=3 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=2 risk_penalty=3
- Fitness: raw=14.00 normalized_0_100=40.00 confidence_0_1=0.85
- Rationale: SIM-zazit is useful mainly as a boundary-setting lower-layer study: it strongly suggests geofence policy should not live inside L5 chunk replication, and it gives local refusal/request/accounting hooks. But it does not model membership, location evidence, or owner-layer release and audit rules, so it only weakly covers geofenced group dispatch.
- Strengths:
  - Clearly keeps L5 feeds meaning-oblivious, reducing pressure to encode group or geofence semantics in chunk transport.
  - Fits PromiseGrid sparse-site and peer-local-accounting goals better than any central geofence oracle approach.
  - Treats refusal and non-response as normal observable outcomes that higher layers could record.
- Weaknesses:
  - No representation of group membership, regional claims, or geofence policy objects.
  - No owner-layer audit trail for why Carol was denied despite being a member.
  - Does not decide whether enforcement belongs at dispatch, fetch, storage, pointer-object, or key-release time.
- Risks:
  - If feed-level refusal is mistaken for complete enforcement, chunks may leak before policy is checked.
  - Without a non-central way to bind location evidence to release decisions, refusals remain hard to audit.
  - Previously replicated chunks create unresolved retention and eviction risk when a peer moves out of region.
- Open questions:
  - What higher-layer artifact binds region policy to specific CAS roots, pointer objects, or keys?
  - What peer-local evidence can audit Alice's location-conditioned refusal without a central location oracle?
  - Should geofenced release be enforced by suppressing advertisements, refusing chunk transfer, withholding pointers, or withholding decryption keys?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-zazit-chunk-feed-replication` x `insurance-claims`

- Result path: `results/SIM-zazit-chunk-feed-replication/insurance-claims/openai-gpt-5.4-xhigh/20260520-215120.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=2 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=3 risk_penalty=2
- Fitness: raw=21.00 normalized_0_100=60.00 confidence_0_1=0.82
- Rationale: Useful as supporting infrastructure for insurance-claims evidence movement: it can transport and verify disputed claim-evidence bytes under sparse, adversarial conditions, but it does not by itself model or decide claim authority, fraud, payments, appeals, or naming/identity semantics.
- Strengths:
  - Excellent L5/L6/L7 boundary clarity: chunk replication stays separate from claim semantics.
  - Strong PromiseGrid fit on sparse knowledge, peer-local accounting, and no central replication authority.
  - Carrier-neutral and multi-repo framing improves migration and evolution safety.
  - ... 1 more
- Weaknesses:
  - Covers evidence transport only; it does not model claimant/adjuster authority, payments, appeals, or fraud decisions.
  - Human-auditable claim outcome reasoning is weak because L5 is intentionally meaning-oblivious.
  - The cross-layer policy for deciding which advertised chunks to pull or retain is still an open design point.
- Risks:
  - May be over-read as insurance-claims coverage when it only addresses the replication substrate.
  - Higher-layer identity, authorization, and dispute semantics could remain under-specified if this feed design is used without a companion claims model.
- Open questions:
  - What higher-layer objects bind replicated chunks to claim, claimant, adjuster, and appeal/payment authority?
  - How should peer-local promise accounting influence pull, retention, and escalation for disputed claim evidence?
  - What minimal audit bundle would let Carol review a contested claim without hidden global state?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-zazit-chunk-feed-replication` x `bgp-class-routing-app-route-hijack`

- Result path: `results/SIM-zazit-chunk-feed-replication/bgp-class-routing-app-route-hijack/openai-gpt-5.4-xhigh/20260520-215120.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=2 risk_penalty=3
- Fitness: raw=27.00 normalized_0_100=68.00 confidence_0_1=0.83
- Rationale: Indirect fit: the simulation does not model BGP-style path advertisement, but its chunk-offer failure cases and peer-local promise accounting provide a plausible local analogue for detecting failed promises and downgrading future peer choices.
- Strengths:
  - Uses peer-local promise accounting instead of central reputation.
  - Covers corrupt, duplicate, refused, and missing chunk exchanges as locally observable failures.
  - Maintains clear L5/L6/higher-layer boundaries and carrier neutrality.
- Weaknesses:
  - No explicit route, path-vector, or multi-hop routing model.
  - Chunk-source peer choice is only a partial stand-in for route choice.
  - Does not show hijack propagation or convergence across relationships.
- Risks:
  - Could overclaim routing-app coverage if feed-level fetch failures are treated as full hijack handling.
  - Unclear how routing policy would consume and weight these local failure records.
- Open questions:
  - What cross-layer interface would let a routing app use feed-level delivery failures without collapsing layer boundaries?
  - How should local penalties compose when a false path claim spans multiple peers rather than one chunk source?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

## Required JSON Shape

{"child_id":"SIM-gifoh-ga-child-0002","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
