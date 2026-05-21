# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Choose a descriptive `child_id` that starts with `SIM-natim-child-` and ends with a kebab-case design slug. Do not use generic `ga-child`, `pending`, or ordinal-only names.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260521-003110`
- Planned child ID prefix: `SIM-natim-child-`
- Temporary child ID: `SIM-natim-child-pending-0001`
- Temporary child path: `proposals/ga-canary-20260521-003110/simulations/SIM-natim-child-pending-0001/`
- Operation: `breed`
- Parent IDs: `SIM-janov-grid-envelope-layer-pcid-nested-signed-payload, SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs`

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

### `simulations/SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs/README.md`

```markdown
# SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs: Grid-envelope Cryptid Multisig signature/proof probe

This simulation is a standalone, non-child grid-envelope specimen. It tests
whether envelope signatures and proofs can use Cryptid's Multisig object model
as the signature/proof payload representation while keeping PromiseGrid's
signature-placement and verification questions unresolved. Source: `DI-sahiv`.

The upstream prior-art source is Cryptid's Multisig Specification v0.0.1,
currently marked pre-draft. This simulation treats that format as design input:
the Multisig object starts with the Multisig sigil `0x1239` encoded as varuint
`0xb924`, then carries a signing-codec sigil, optional message bytes, and a
counted sequence of attributes. Source: `DI-sahiv`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.

## Design Pressure

- **Detached versus combined:** the Multisig message field may be empty
  (detached) or present (combined), so the sim can compare signatures over
  outer envelope bytes, nested payload bytes, or in-object message bytes.
- **Envelope versus nested payload:** the same Multisig bytes can occupy an
  outer signature/proof slot or live inside a payload protocol's nested object.
- **Variable arity:** Multisig's counted attributes let one signature object
  carry extra codec-specific proof material without changing the outer envelope
  shape.
- **pCID interaction:** the envelope `pcid` still selects payload semantics, and
  a `sig_pcid` or payload schema decides how Multisig verification is invoked.
- **Unknown codecs:** generic tools that understand varuint and varbytes can
  skip unknown Multisig signing codecs or unknown attributes without claiming
  verification.
- **Threshold shares:** attributes such as `Scheme`, `Threshold`, `Limit`,
  `ShareIdentifier`, and `ThresholdData` let the sim test individual shares,
  accumulation, and final aggregate verification.
- **Verifier obligations:** verifiers must bind the exact message bytes, signing
  codec, verifying key material, required attributes, threshold policy, and
  payload interpretation before accepting an envelope.

## Non-Canonical Status

This simulation does not choose a final PromiseGrid envelope shape, does not
freeze a pCID, does not require Cryptid Multisig as a PromiseGrid dependency,
and does not supersede the existing positional, arity, nested-signature, or
generated child grid-envelope specimens. Source: `DI-sahiv`; open long-term
envelope pressure remains represented by `DR-009-20260430-204108`.
```

### `simulations/SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs/QUESTION.md`

```markdown
# Question

Can grid-envelope signatures and proofs use Cryptid's Multisig object model as
their signature/proof payload representation while preserving unresolved
PromiseGrid choices about detached versus combined signatures, outer versus
nested placement, variable arity, pCID binding, unknown-codec handling,
threshold-share aggregation, and verifier obligations? Source: `DI-sahiv`;
`DR-009-20260430-204108`.
```

### `simulations/SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs/protocols/grid-envelope.d/CHANGELOG.md`

```markdown
# CHANGELOG: grid-envelope

A-side CHANGELOG (per TE-liviv) for this simulation-local `grid-envelope`
protocol specimen.

This file records freeze events authored by the specimen maintainers. No entries
yet; this protocol specimen has not reached a first freeze.

This protocol tree is a simulation-local specimen created by `DI-sahiv`.
```

### `simulations/SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid-envelope draft: Cryptid Multisig signature/proof payloads

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `cryptid-multisig-signature-proofs`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a
specimen inside `SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs`,
not a harness rule and not the canonical PromiseGrid envelope. Source:
`DI-sahiv`.

The design under test is narrow: use Cryptid's Multisig object model as the
signature/proof payload representation while leaving the envelope-level
placement and verification choices open for simulation pressure. The Multisig
source is upstream prior art and is not treated as a frozen PromiseGrid
dependency. Source: `DI-sahiv`.

## Prior-Art Multisig Shape

Cryptid's pre-draft Multisig v0.0.1 encodes a codec-agnostic digital signature
object as:

```text
multisig_sigil signing_codec_sigil message attributes
```

The object starts with Multisig sigil `0x1239`, encoded as varuint `0xb924`.
It then carries a signing-codec sigil, an optional `message` encoded as
varbytes, and a variable number of attributes encoded as a count followed by
attribute-id and varbytes pairs. Source: `DI-sahiv`.

This simulation recognizes these upstream attribute roles as pressure inputs:

- `SigData` for signature bytes;
- `PayloadEncoding` for the signed-message encoding sigil;
- `Scheme` for threshold-signing scheme;
- `Threshold` for the minimum share count required;
- `Limit` for total share count;
- `ShareIdentifier` for the share number or participant-local share label;
- `ThresholdData` for codec-specific threshold material;
- `AlgorithmName` for application-specific or non-standard algorithm naming.

## Envelope Shapes Under Test

The simulation keeps three placement modes alive instead of choosing one:

```text
[pcid, payload, sig_pcid, multisig]
[pcid, payload_with_nested_multisig]
[pcid, combined_multisig]
```

Slots are interpreted positionally only when the selected mode and its `pcid`
or `sig_pcid` define them:

- `pcid` identifies the payload protocol, handler, or proof-bearing schema.
- `payload` is opaque bytes until interpreted by the handler named by `pcid`.
- `sig_pcid`, when present, identifies the signature/proof protocol that says
  the fourth slot is a Cryptid-style Multisig object.
- `multisig` is the exact Multisig object bytes, not decoded fields projected
  into the envelope.

The first mode pressures explicit outer signature dispatch. The second mode
pressures nested payload ownership of signatures. The third mode pressures
combined Multisig objects where the signed message is carried inside the
Multisig `message` varbytes field rather than as a sibling envelope payload.
Source: `DI-sahiv`; `DR-009-20260430-204108`.

## Encoding

The envelope carrier for this specimen is deterministic CBOR positional arrays.
`pcid` and `sig_pcid`, when present, are CIDv1 byte strings or DAG-CBOR links as
defined by the concrete run profile. `payload` and `multisig` are byte strings.
The Multisig object itself keeps its own varuint and varbytes internal encoding;
the envelope does not translate Multisig attributes into CBOR fields.

Canonical envelope bytes are the deterministic CBOR bytes of the selected
outer array. Canonical Multisig bytes are the exact varuint/varbytes bytes
carried in the Multisig slot or nested payload. A verifier must never verify a
re-serialized approximation when the original bytes are available.

## Detached and Combined Signature Policy

A Multisig object with an empty message field is treated as detached. The
verifier must obtain the signed bytes from the selected envelope mode:

- outer explicit mode signs the canonical unsigned prefix `[pcid, payload]`
  unless `sig_pcid` defines stricter associated data;
- nested mode signs the nested bytes selected by the payload protocol named by
  `pcid`;
- share-collection mode signs the same byte string for every share before
  threshold aggregation.

A Multisig object with a non-empty message field is treated as combined. The
verifier must compare the embedded message bytes against the envelope-selected
payload binding before accepting the envelope. If the embedded message and
outer payload disagree, verification fails even if the Multisig cryptographic
check succeeds. Source: `DI-sahiv`.

## Unknown Codec and Attribute Policy

This specimen adopts a conservative unknown-codec rule:

- A receiver that does not understand the envelope `pcid` must not interpret
  `payload`, even if it can skip or parse Multisig framing.
- A receiver that understands `pcid` but not `sig_pcid` may preserve the exact
  Multisig bytes as opaque proof evidence, but must not claim verification.
- A receiver that understands Multisig varuint/varbytes framing but not the
  signing codec may skip the object, index its byte range, and preserve it for
  later tooling, but must not treat unknown `SigData` as valid.
- A receiver may ignore unknown non-critical attributes only if the signing
  codec or `sig_pcid` says they are non-critical; otherwise unknown attributes
  keep the verification result unsupported.

This policy intentionally separates structural skippability from cryptographic
acceptance. Skipping unknown signing codecs helps storage, relay, and future
audit; it is not a validity decision. Source: `DI-sahiv`.

## Threshold and Multi-Payload Pressure

Threshold runs use the Multisig attributes as follows:

- every share carries `Scheme`, `Threshold`, `Limit`, and `ShareIdentifier`
  values that must agree with the threshold policy named by `sig_pcid` or the
  payload schema;
- `SigData` carries the share or aggregate signature data as defined by the
  signing codec;
- `ThresholdData` carries codec-specific accumulation material when the codec
  needs more than raw signature bytes;
- `PayloadEncoding` records the signed-message encoding when the signature
  codec requires that value for replay-safe verification;
- `AlgorithmName` is advisory unless the selected `sig_pcid` makes it part of
  the verification policy.

The simulation should penalize designs that let two shares over different
message bytes, different `pcid` values, or different threshold policies
aggregate as if they belonged to the same proof.

## pCID Interaction

The envelope `pcid` continues to name the payload protocol. Multisig does not
replace pCID-selected payload semantics. A `sig_pcid`, nested payload schema, or
future frozen profile must define:

- whether the signed bytes include `pcid`, `payload`, both, or an enclosing
  transcript;
- whether `PayloadEncoding` duplicates, complements, or constrains the pCID;
- which Multisig signing codecs are acceptable for the payload protocol;
- how verifier key material is found, authenticated, rotated, and revoked;
- whether unknown attributes are reject, ignore, quarantine, or relay-only
  evidence.

Until those choices are frozen, this specimen is evidence for verifier
obligations rather than a PromiseGrid validity rule. Source: `DI-sahiv`;
`DR-009-20260430-204108`.

## Verifier Obligations

A verifier accepts an envelope under this specimen only after all of the
following hold:

- the envelope shape is recognized by `pcid`, `sig_pcid`, or a local run
  profile;
- the exact signed byte string is determined without ambiguity;
- detached and combined-message bindings agree with the selected envelope mode;
- the signing codec is understood and allowed by the selected signature policy;
- required attributes are present, canonical, non-duplicated unless the codec
  permits duplicates, and internally consistent;
- threshold shares all bind to the same message, pCID context, threshold policy,
  and signer set before aggregation;
- unknown codecs or critical attributes produce an unsupported or quarantined
  result rather than a successful verification;
- local audit records retain exact envelope bytes, exact Multisig bytes, the
  verification profile, and the reason for accept, reject, unsupported, or
  quarantine.

## Scenario Pressure Notes

### Normal detached signature

Alice sends `[pcid, payload, sig_pcid, multisig]` where `multisig` has an empty
message field. Bob verifies that the signature covers the canonical unsigned
prefix `[pcid, payload]`, not only `payload`, so replay under a different
payload protocol fails.

### Combined signature mismatch

Carol sends `[pcid, payload, sig_pcid, multisig]` where the Multisig message
field contains bytes that differ from `payload`. Dave must reject the envelope
because the embedded message and envelope-selected payload binding disagree.

### Unknown signing codec

Ellen understands the envelope `pcid` and Multisig framing but lacks the
signing-codec implementation. Ellen may keep and relay the exact proof bytes as
unsupported evidence, but she must not mark the envelope verified.

### Threshold share collection

Frank receives three BLS-style shares that claim a threshold of three out of
four. He aggregates only shares whose message bytes, `pcid` binding, scheme,
threshold, limit, and signer-set policy match; mixed-context shares remain
separate unsupported evidence.

### Nested payload proof

Alice sends `[pcid, payload_with_nested_multisig]`. Bob can verify only after
the `pcid` payload schema identifies the nested Multisig byte range and the
exact bytes the nested proof signs.

## Non-Goals

This draft does not:

- declare a winning envelope;
- freeze a pCID;
- require a central pCID registry;
- decide detached versus combined signatures globally;
- decide whether signatures live in the envelope or nested payload;
- define final PromiseGrid key discovery, revocation, freshness, or authority;
- claim that Cryptid Multisig is stable or normative for PromiseGrid.

## Freeze Gate

This draft can freeze only after simulation runs compare it against sibling
grid-envelope signature-placement specimens, at least one verifier profile
fully specifies pCID binding and unknown-codec behavior, and a maintainer signs
a merge/freeze promise for this specific specimen.
```

## Compact Fitness Evidence From This Run

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

### `SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs` x `chunk-feed-replication-sparse-advertisement`

- Result path: `results/SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs/chunk-feed-replication-sparse-advertisement/openai-gpt-5.4-xhigh/20260521-003110.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=3 evolution_safety=3 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=2 risk_penalty=3
- Fitness: raw=14.00 normalized_0_100=40.00 confidence_0_1=0.89
- Rationale: Good envelope-level proof specimen, but it does not model chunk feeds, Merkle advertisements, or partial-replica discovery. For this scenario it is only an indirect building block for signing sparse advertisements, not a direct answer.
- Strengths:
  - Exact-byte retention and local verification records support later audit under sparse knowledge.
  - Unknown-codec and unsupported/quarantine rules avoid claiming verification from partial understanding.
  - Clear separation of payload `pcid`, signature policy, and proof bytes would help layer a future sparse-advertisement protocol.
- Weaknesses:
  - No mechanism for advertising leaves, roots, pointer objects, frontiers, or compact chunk summaries.
  - No protocol for peers with disjoint chunk subsets to compare availability or request missing CAS objects.
  - Key discovery, binding policy, and accepted signature profiles remain unresolved.
- Risks:
  - If reused without a frozen sparse-advertisement profile, signed claims could replay or misbind across payload contexts.
  - Stored opaque proof bytes may be mistaken for verified chunk-availability evidence by downstream tooling.
- Open questions:
  - What payload schema should encode sparse chunk advertisements when no peer has the full Merkle tree?
  - What exact bytes must be signed so root, frontier, and chunk claims cannot replay under a different `pcid` or root?
  - How should partial-replica peers merge or challenge availability claims without central indexes?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs` x `conditional-release-geofencing-onward-restraint-chain`

- Result path: `results/SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs/conditional-release-geofencing-onward-restraint-chain/openai-gpt-5.4-xhigh/20260521-003110.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=2 evolution_safety=3 layer_boundary_clarity=4 failure_handling=3 implementation_plausibility=3 risk_penalty=4
- Fitness: raw=19.00 normalized_0_100=47.50 confidence_0_1=0.90
- Rationale: Useful only as a cryptographic substrate: its outer-vs-nested proof placement and exact-byte verification could carry signed conditional-release artifacts, but it does not define the recursive onward-restraint promise graph, geofencing policy, or forwarding semantics this scenario is testing.
- Strengths:
  - Outer-slot versus nested-payload proof placement gives a hook for later transport-level or CAS-object-level policy evidence.
  - Exact-byte retention and local accept/reject/unsupported records help peer-local audit of proof material.
  - Conservative unknown-codec handling is evolution-friendly and avoids false verification claims.
- Weaknesses:
  - No model of Alice/Bob/Carol onward-promise chaining or recipient assent.
  - Does not decide whether the restraint graph lives at group-session, conditional-release, transport/feed, or CAS-object level.
  - Key discovery, revocation, and policy binding remain unresolved.
- Risks:
  - A reviewer could mistake proof carriage for actual enforcement of onward-restraint rules.
  - Different pcid/sig_pcid profiles could bind different policy bytes and diverge on what was promised.
  - The upstream Cryptid Multisig shape is still pre-draft.
- Open questions:
  - Should the onward-restraint graph be a signed payload/CAS object carried by this envelope, or a separate session/feed protocol?
  - How are recipient identity, geography conditions, and re-forward promises bound into the signed bytes?
  - What local evidence lets Alice later audit that Carol accepted the same restraint before Bob forwarded content?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs` x `bgp-routing`

- Result path: `results/SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs/bgp-routing/openai-gpt-5.4-xhigh/20260521-003110.json`
- Scores: scenario_fit=1 promisegrid_alignment=3 auditability=4 evolution_safety=3 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=3 risk_penalty=4
- Fitness: raw=18.00 normalized_0_100=45.00 confidence_0_1=0.82
- Rationale: Strong as a low-level signed-proof envelope specimen, but only a partial fit for BGP routing. It helps with exact-byte binding, local audit, unknown-codec quarantine, and threshold-consistency checks, yet it does not model routing promises, path semantics, freshness, or leak/hijack policy.
- Strengths:
  - Explicit binding of payload, pCID context, and detached/combined signature rules reduces replay and ambiguity.
  - Requires local retention of exact envelope bytes, exact proof bytes, verification profile, and accept/reject/unsupported reasons.
  - Conservative handling of unknown codecs/attributes and mixed-context threshold shares fits sparse-knowledge auditing.
- Weaknesses:
  - No route object model, path/neighbor semantics, or promise-accounting model for hijacks, leaks, and contested reachability.
  - Freshness, revocation, key discovery, and signer authorization are explicitly unresolved.
  - Non-canonical draft built on pre-draft prior art, so deployment assumptions remain unstable.
- Risks:
  - Valid-looking but stale route evidence could be replayed unless a routing profile adds freshness and policy checks.
  - Without route-specific signer-set and path semantics, proofs may be attributable but still insufficient to judge leaks or hijacks.
  - Key-management completion could drift toward brittle or centralized trust anchors.
- Open questions:
  - What routing payload schema binds prefix, path, policy, and freshness into the signed transcript?
  - How are route-announcer keys discovered, rotated, revoked, and audited without a central registry?
  - What local accept/downgrade/quarantine policy should apply when proofs are structurally intact but stale, partial, or codec-unsupported?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

## Required JSON Shape

{"child_id":"SIM-natim-child-descriptive-design-slug","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
