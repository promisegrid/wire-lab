# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Choose a descriptive `child_id` that starts with `SIM-jogoh-child-` and ends with a kebab-case design slug. Do not use generic `ga-child`, `pending`, or ordinal-only names.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260524-dalor-pobod-forced`
- Planned child ID prefix: `SIM-jogoh-child-`
- Temporary child ID: `SIM-jogoh-child-pending-0001`
- Temporary child path: `proposals/ga-canary-20260524-dalor-pobod-forced/simulations/SIM-jogoh-child-pending-0001/`
- Operation: `breed`
- Parent IDs: `SIM-pobod-grid-envelope-outer-promise-nested-signed-payload, SIM-dalor-grid-envelope-protocol-owned-signature-slot`

## Scenario Sample

- `promise-economy-spectrum-peer-local-assessment-only` at `scenarios/promise-economy-spectrum-peer-local-assessment-only/promise-economy-spectrum-peer-local-assessment-only.md`
- `promise-accounting-records-kept-storage-promise` at `scenarios/promise-accounting-records-kept-storage-promise/promise-accounting-records-kept-storage-promise.md`
- `promise-accounting-records-refused-service` at `scenarios/promise-accounting-records-refused-service/promise-accounting-records-refused-service.md`
- `conditional-release-geofencing-onward-restraint-chain` at `scenarios/conditional-release-geofencing-onward-restraint-chain/conditional-release-geofencing-onward-restraint-chain.md`
- `peer-adoption-metadata-relationship-scoped-adoption` at `scenarios/peer-adoption-metadata-relationship-scoped-adoption/peer-adoption-metadata-relationship-scoped-adoption.md`
- `group-session-freeze-promise-post-freeze-mutation-request` at `scenarios/group-session-freeze-promise-post-freeze-mutation-request/group-session-freeze-promise-post-freeze-mutation-request.md`
- `promise-economy-spectrum-transferable-promise-token` at `scenarios/promise-economy-spectrum-transferable-promise-token/promise-economy-spectrum-transferable-promise-token.md`
- `promise-economy-spectrum-permissioned-capability-token` at `scenarios/promise-economy-spectrum-permissioned-capability-token/promise-economy-spectrum-permissioned-capability-token.md`

## Scenario Pressure

### `scenarios/promise-economy-spectrum-peer-local-assessment-only/promise-economy-spectrum-peer-local-assessment-only.md`

```markdown
# Peer-local assessment only

## Scenario ID

promise-economy-spectrum-peer-local-assessment-only

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-haros-promise-economy-spectrum/SCENARIOS.md`
- Source simulation: `SIM-haros-promise-economy-spectrum/`
- Source row/title: Peer-local assessment only
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-haros-promise-economy-spectrum/`.

## Setup

Alice records that Bob kept storage promises and Mallory sent corrupt chunks.

## Stimulus

Run the candidate simulation against this source test: Whether pull/keep/advertise decisions can use local observations without token transfer.

## Expected Pressure

The base protocol should not require a token field if peer-local assessment is enough for some groups.
```

### `scenarios/promise-accounting-records-kept-storage-promise/promise-accounting-records-kept-storage-promise.md`

```markdown
# Kept storage promise

## Scenario ID

promise-accounting-records-kept-storage-promise

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-rusap-promise-accounting-records/SCENARIOS.md`
- Source simulation: `SIM-rusap-promise-accounting-records/`
- Source row/title: Kept storage promise
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-rusap-promise-accounting-records/`.

## Setup

Alice promises to store chunk C for a period; Bob later asks and Alice serves it.

## Stimulus

Run the candidate simulation against this source test: What Bob records locally and how that affects future pull / keep / advertise choices.

## Expected Pressure

Specs may need promise-vocabulary sections that name observable promise outcomes.
```

### `scenarios/promise-accounting-records-refused-service/promise-accounting-records-refused-service.md`

```markdown
# Refused service

## Scenario ID

promise-accounting-records-refused-service

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-rusap-promise-accounting-records/SCENARIOS.md`
- Source simulation: `SIM-rusap-promise-accounting-records/`
- Source row/title: Refused service
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-rusap-promise-accounting-records/`.

## Setup

Alice refuses to send C because of policy, cost, group context, or missing authorization.

## Stimulus

Run the candidate simulation against this source test: Whether refusal is recorded differently from failure, corruption, or timeout.

## Expected Pressure

Promise accounting must support honest refusal instead of treating every refusal as misbehavior.
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

### `scenarios/peer-adoption-metadata-relationship-scoped-adoption/peer-adoption-metadata-relationship-scoped-adoption.md`

```markdown
# Relationship-scoped adoption

## Scenario ID

peer-adoption-metadata-relationship-scoped-adoption

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-dihiz-peer-adoption-metadata/SCENARIOS.md`
- Source simulation: `SIM-dihiz-peer-adoption-metadata/`
- Source row/title: Relationship-scoped adoption
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-dihiz-peer-adoption-metadata/`.

## Setup

Bob follows one profile when talking to Alice and another when talking to Carol.

## Stimulus

Run the candidate simulation against this source test: Whether adoption promises bind to relationship context instead of only peer identity.

## Expected Pressure

Peer-local variation may be legitimate and should not look like equivocation unless it violates a promise.
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

### `scenarios/promise-economy-spectrum-transferable-promise-token/promise-economy-spectrum-transferable-promise-token.md`

```markdown
# Transferable promise token

## Scenario ID

promise-economy-spectrum-transferable-promise-token

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-haros-promise-economy-spectrum/SCENARIOS.md`
- Source simulation: `SIM-haros-promise-economy-spectrum/`
- Source row/title: Transferable promise token
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-haros-promise-economy-spectrum/`.

## Setup

Bob transfers an Alice-issued promise token to Carol.

## Stimulus

Run the candidate simulation against this source test: Whether provenance, permission, double-spend prevention, and refusal semantics stay local and auditable.

## Expected Pressure

Transferability creates obligations that the base protocol may need to defer.
```

### `scenarios/promise-economy-spectrum-permissioned-capability-token/promise-economy-spectrum-permissioned-capability-token.md`

```markdown
# Permissioned capability token

## Scenario ID

promise-economy-spectrum-permissioned-capability-token

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-haros-promise-economy-spectrum/SCENARIOS.md`
- Source simulation: `SIM-haros-promise-economy-spectrum/`
- Source row/title: Permissioned capability token
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-haros-promise-economy-spectrum/`.

## Setup

Alice issues Bob a non-transferable capability to fetch a chunk family.

## Stimulus

Run the candidate simulation against this source test: Whether permission and redemption can be scoped without implying exchange rates.

## Expected Pressure

Capability semantics may be useful while still avoiding marketplace behavior.
```

## Parent Simulation Documents

### `simulations/SIM-pobod-grid-envelope-outer-promise-nested-signed-payload/README.md`

```markdown
# SIM-pobod-grid-envelope-outer-promise-nested-signed-payload

This simulation is the Promise-Theory-first successor to
`SIM-janov-grid-envelope-layer-pcid-nested-signed-payload`. It keeps the nested
signed-payload idea but makes the outer-layer promise explicit: an agent
promises that the carried nested bytes match the pCID-defined structure and that
the nested signature is the evidence object for the inner payload claim. Source:
`DI-tavaz`.

## Design Under Test

The outer envelope remains small, but its semantics are restated in Promise
Theory terms:

- the sender promises that the outer pCID defines how to parse the nested body;
- the nested body carries the actual payload pCID, payload bytes, and signature;
- receivers decide locally whether they trust the sender enough to believe the
  outer promise and whether the nested signature is enough evidence for the
  inner payload claim;
- no unsigned outer field is treated as self-authenticating authority.

## Why this differs from `janov`

`janov` described the shape well but left the outer-layer promise implicit,
which made authorship and trust too dependent on unstated transport context.
`pobod` keeps the same structural pressure while making the promise boundary
explicit.

## Boundaries

This simulation does not declare the winning grid-envelope format. It only asks
whether a nested signed payload becomes more PT-clean once the outer-layer
promise is made explicit.
```

### `simulations/SIM-pobod-grid-envelope-outer-promise-nested-signed-payload/QUESTION.md`

```markdown
# Question

Can a grid envelope keep a small outer shape while explicitly modeling the
sender's outer-layer promise about nested structure and signed payload meaning,
so trust stays local and authorship does not collapse into implicit transport
authority?

Source: `DI-tavaz`.
```

### `simulations/SIM-dalor-grid-envelope-protocol-owned-signature-slot/README.md`

```markdown
# SIM-dalor-grid-envelope-protocol-owned-signature-slot: Protocol-owned outer signature-slot probe

This simulation is a standalone grid-envelope specimen. It tests a three-slot
outer envelope `grid([pCID, payload, signature])` where the outer third slot is
mandatory but the proof family is still owned by the protocol named by `pCID`
rather than by a
separate outer `sig_pcid`. The point is to test the smallest explicit
outer-signature design that still allows varsig-, multisig-, or other
pCID-defined proof rules. In Promise Theory terms, the current sender's
signature is evidence of the sender's promise that the payload bytes are shaped
according to the protocol specification named by the `pCID`; higher-layer
promise accounting remains inside the payload protocol. Source: `DI-kukuk`;
`DI-pozom`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.

Primary comparison targets: `SIM-kurim`, `SIM-jufag`, `SIM-pamap`, and
`SIM-jumav`.
```

### `simulations/SIM-dalor-grid-envelope-protocol-owned-signature-slot/QUESTION.md`

```markdown
# Question

Can PromiseGrid use a three-slot outer envelope
`grid([pCID, payload, signature])` where the third slot signs canonical
`[pCID, payload]` bytes and the protocol specification named by `pCID` defines
the proof family and verification rules, without needing a separate outer
`sig_pcid` selector? Source: `DI-kukuk`; `DI-pozom`.
```

### `simulations/SIM-dalor-grid-envelope-protocol-owned-signature-slot/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid-envelope draft: protocol-owned signature slot

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `protocol-owned-signature-slot`.

## Scope

This spec defines one grid-envelope candidate for wire-lab comparison. It is a
specimen inside `SIM-dalor-grid-envelope-protocol-owned-signature-slot`, not a
harness rule and not the canonical PromiseGrid envelope. It uses `pCID` only as
Protocol CID: the content identifier of the protocol specification document, not
the content identifier of the payload bytes. Source: `DI-kukuk`; `DI-pozom`.

## Envelope Shape

The outer envelope shape is:

```text
[pCID, payload, signature]
```

Slots are interpreted positionally:

- `pCID` identifies the protocol specification and the proof rules for the third
  slot.
- `payload` is opaque bytes until interpreted by the handler for the protocol
  named by `pCID`.
- `signature` is mandatory proof bytes over the canonical unsigned prefix.

This is the key design move under test: the proof remains a sibling outer slot,
but there is no separate outer proof-profile selector. If a protocol wants a
varsig, multisig, or other proof family, that choice is part of the protocol
named by `pCID`.

In Promise Theory terms, the signature is not a command, permission, global
trust score, or promise made on behalf of anyone else. It is current-sender
evidence for the sender's own scoped promise: "I promise that these payload
bytes are shaped according to the protocol specification named by this `pCID`."
Each receiver still decides locally whether it trusts the sender, recognizes the
protocol, verifies the proof family, stores the bytes, relays the bytes, or uses
the payload.

## Signable Bytes

The signed bytes are the canonical bytes of:

```text
[pCID, payload]
```

This binds both the payload bytes and the protocol name without adding a second
protocol selector to the outer envelope.

## Encoding

The outer envelope is a deterministic CBOR positional array. `pCID` is a CIDv1
byte string or link as defined by the carrier profile. `payload` and
`signature` are byte strings at the carrier layer. The CBOR array header already
records arity, so this specimen does not add an outer arity field.

## Unknown pCID Policy

If a receiver lacks a handler for `pCID`, it may preserve or blind-carry the
exact outer bytes as uninterpreted evidence, but it MUST NOT claim to parse the
payload shape or verify the signature.

This specimen intentionally stays close to the current simple-envelope direction:
unsupported `pCID` means “bytes may survive, meaning does not.”

## Signature and Authorship Policy

This specimen has no universal outer `sig_pcid` slot. Instead, the protocol
named by `pCID` defines:

- whether the third slot is interpreted as varsig, multisig, or another proof
  family;
- signer binding and signer identity rules;
- delegation, freshness, revocation, and threshold semantics, if any;
- whether extra associated data beyond canonical `[pCID, payload]` bytes is
  required.

The outer envelope itself enforces only three things:

- there is a third proof slot;
- the signable baseline is canonical `[pCID, payload]`;
- proof semantics are owned by the protocol named by `pCID`.

## Comparison Pressure

Compared with `SIM-jufag`, this specimen removes the separate outer proof
selector and asks whether one protocol `pCID` is enough to name both payload
shape and proof semantics.

Compared with `SIM-pamap`, this specimen keeps the proof as an outer sibling
slot rather than moving it into the payload contract.

Compared with `SIM-jumav`, this specimen avoids wrapper indirection and signs
the carried payload directly.

Compared with `SIM-kurim`, this specimen adds one universal outer proof slot
while still trying to keep the outer surface small.

## Open Questions

- Is one protocol `pCID` enough to keep proof-family evolution legible, or does a
  separate outer proof selector age better?
- Does this design force too many proof-profile changes to mint new protocol
  `pCID`s?
- Do generic peers lose too much audit clarity when the outer slot is explicit
  but its proof family is defined by the protocol named by `pCID`?

## Non-Canonical Status

This draft does not declare a winning envelope and does not constrain sibling
simulations. It exists to compare the protocol-owned three-slot outer proof idea
against minimal outer envelopes, explicit outer proof selectors, payload-owned
proofs, and wrapper-proof designs.
```

## Compact Fitness Evidence From This Run

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `promise-economy-spectrum-peer-local-assessment-only`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/promise-economy-spectrum-peer-local-assessment-only/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=3 promisegrid_alignment=5 auditability=4 evolution_safety=3 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=5 risk_penalty=2
- Fitness: raw=42.00 normalized_0_100=84.00 confidence_0_1=0.84
- Rationale: Strong PromiseGrid envelope specimen: small, explicit, signed, and locally assessable. It fits the scenario indirectly by avoiding token requirements and preserving local trust boundaries, but it does not itself define storage or advertise decision logic, so scenario fit is only moderate at this layer.
- Strengths:
  - Clear envelope-layer boundary: payload semantics and promise accounting are delegated explicitly.
  - Promise-first wording is strong and PT-correct at the claimed layer.
  - Small deterministic three-slot artifact is durable and plausible for constrained devices.
  - ... 1 more
- Weaknesses:
  - Scenario pressure about keep or advertise decisions is mostly deferred to payload protocols and local policy.
  - Generic auditability is weaker than designs with an explicit outer proof-profile selector.
  - Failure handling is limited to preserve-or-ignore behavior rather than richer envelope-level diagnostics.
- Risks:
  - Binding proof semantics to the same pCID as payload shape may force extra protocol churn when proof families evolve.
  - A visible outer signature slot may be overinterpreted as broader trust or identity assurance than the envelope actually promises.
  - Unknown-pCID handling preserves bytes but limits immediate audit or decision support for generic peers.
- Open questions:
  - Is coupling proof-family evolution to the payload protocol pCID durable enough for long-lived deployments?
  - Will generic peers have enough audit clarity when they can see a mandatory proof slot but must consult the pCID spec to know verification semantics?
  - What minimal payload-level evidence is needed so peers can make keep or advertise decisions from local observations alone?
- Authority boundary: Envelope-layer evidence only; peer-local trust/accounting and keep or advertise decisions remain payload-protocol and local-policy matters.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `promise-accounting-records-kept-storage-promise`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/promise-accounting-records-kept-storage-promise/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=2 promisegrid_alignment=5 auditability=3 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 risk_penalty=2
- Fitness: raw=37.00 normalized_0_100=74.00 confidence_0_1=0.85
- Rationale: Strong Promise-Theory-clean envelope specimen with explicit layer boundaries, but this scenario asks for higher-layer storage-promise accounting that the design intentionally delegates to the payload protocol. That delegation is clear and legitimate, yet it limits direct scenario coverage.
- Strengths:
  - Very clear envelope-layer scope and delegation boundary.
  - Promise-first wording is unusually strong and PT-correct for an envelope design.
  - Mandatory outer proof over canonical [pCID, payload] gives compact, auditable sender evidence.
  - ... 1 more
- Weaknesses:
  - Does not itself model Alice's storage promise, kept-status evidence, or Bob's local accounting records.
  - Scenario fit is limited because storage semantics live entirely in the payload protocol.
  - Protocol-owned proof semantics may reduce generic audit legibility compared with an explicit outer proof-profile selector.
- Risks:
  - A receiver could mistakenly treat the outer signature as evidence of storage success rather than only payload-shape/authorship evidence.
  - Proof-family evolution may require frequent new protocol pCIDs, creating long-term audit and migration friction.
  - Generic peers may have weaker cross-protocol audit clarity because proof semantics are not named by a separate outer selector.
- Open questions:
  - What payload-level artifact would let Bob record that Alice kept a storage promise over time?
  - Does coupling proof-family semantics to the payload pCID force too many new protocol pCIDs as signature schemes evolve?
  - Can generic auditors understand enough from the outer envelope alone when proof semantics are protocol-owned?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `promise-accounting-records-refused-service`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/promise-accounting-records-refused-service/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=1 promisegrid_alignment=4 auditability=4 evolution_safety=3 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=5 risk_penalty=2
- Fitness: raw=37.00 normalized_0_100=74.00 confidence_0_1=0.89
- Rationale: Strong Promise-Theory-clean envelope design with explicit local trust and excellent layer boundaries. It is small, deterministic, and auditable at the envelope layer, but this scenario asks for promise-accounting semantics that distinguish honest refusal from failure classes, and the candidate intentionally delegates that to the payload protocol. That keeps it promotable as an envelope specimen but limits direct scenario fit.
- Strengths:
  - Explicit promise-first framing: sender signs evidence that payload bytes are shaped according to the protocol named by pCID.
  - Very clear envelope/payload boundary; it does not pretend to solve higher-layer promise accounting.
  - Deterministic three-slot CBOR envelope is small, durable, and plausible for constrained devices.
  - ... 1 more
- Weaknesses:
  - Does not itself encode or classify refused service records.
  - Failure/refusal distinctions depend entirely on payload protocol design.
  - Proof-family semantics hidden behind pCID reduce generic cross-protocol audit legibility.
  - ... 1 more
- Risks:
  - Payload protocols may fail to record honest refusal distinctly, leaving scenario pressure unresolved above the envelope.
  - Protocol-owned proof semantics may create pCID churn when proof families evolve.
  - Generic peers have reduced audit clarity for unknown pCIDs because proof semantics are not outer-universal.
- Open questions:
  - How should payload protocols encode honest refusal distinctly from timeout, corruption, or transport failure?
  - Will proof-family evolution force excessive new pCIDs without a separate outer proof selector?
  - Can generic auditors compare refusal evidence across pCIDs when proof semantics are protocol-owned?
- Authority boundary: Evidence only; envelope-layer specimen, not final PromiseGrid or promise-accounting authority.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `conditional-release-geofencing-onward-restraint-chain`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/conditional-release-geofencing-onward-restraint-chain/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=3 evolution_safety=3 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=5 risk_penalty=2
- Fitness: raw=37.00 normalized_0_100=74.00 confidence_0_1=0.86
- Rationale: Strong Promise-Theory-clean envelope specimen with a clear layer boundary: it offers signed evidence that payload bytes are shaped according to the protocol named by pCID, while leaving conditional-release and recursive onward-restraint semantics to the payload protocol. That makes it a reasonable carrier for this scenario but not a direct solution to the scenario's core graph semantics.
- Strengths:
  - Very clear envelope/payload boundary.
  - Promise-first wording is explicit and PT-correct.
  - Deterministic small artifact with mandatory outer proof slot.
  - ... 1 more
- Weaknesses:
  - Does not itself represent the recursive restraint chain the scenario asks about.
  - Generic peers cannot fully audit signature semantics without fetching the pCID-defined protocol spec.
  - Proof semantics are coupled to payload protocol evolution rather than separated in the outer layer.
- Risks:
  - Scenario pressure may be under-served if payload protocols do not standardize onward-restraint accounting clearly.
  - Proof-family semantics hidden behind the payload pCID may reduce generic cross-protocol audit clarity.
  - Changing proof semantics may force new protocol pCIDs, creating coupling between payload evolution and signature policy.
- Open questions:
  - Does coupling proof-family evolution to the payload protocol pCID create too many new protocol IDs over time?
  - Can generic auditors understand enough from the outer envelope when signature verification rules are only discoverable via the payload protocol spec?
  - What payload-level promise/accounting shape best represents recursive onward-restraint chains while preserving local trust and sparse knowledge?
- Authority boundary: Evidence only; this envelope-layer design can carry but does not itself define onward-restraint or conditional-release semantics, which must live in the payload protocol and each peer's local trust assessment.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `peer-adoption-metadata-relationship-scoped-adoption`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/peer-adoption-metadata-relationship-scoped-adoption/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=2 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 risk_penalty=2
- Fitness: raw=39.00 normalized_0_100=78.00 confidence_0_1=0.84
- Rationale: Strong envelope-layer PromiseGrid candidate: small, explicit, deterministic, and promise-first. For this scenario, relationship-scoped adoption behavior is mostly outside the envelope layer and must be expressed by the payload protocol, so fit is only partial rather than poor. The design stays PT-clean because that delegation is explicit and local trust boundaries are preserved.
- Strengths:
  - Very clear layer boundary: the envelope authenticates protocol+payload bytes while leaving higher-layer adoption semantics to the payload protocol.
  - Excellent Promise Theory framing with sender-scoped intent and local receiver trust decisions.
  - Simple three-slot deterministic envelope is compact and durable for small-device and long-horizon use.
  - ... 1 more
- Weaknesses:
  - Does not itself model relationship-scoped adoption metadata, so scenario-specific behavior is only indirectly supported.
  - Outer proof semantics are less legible to generic tools because the proof family is protocol-owned rather than separately declared.
  - Failure semantics around stale or conflicting relationship-specific adoption claims depend entirely on payload protocols, not the envelope.
- Risks:
  - Scenario pressure may be underserved if evaluators expect relationship-scoped adoption semantics at the envelope layer.
  - Generic auditability is somewhat reduced because proof-family semantics are hidden behind the payload protocol pCID rather than an explicit outer selector.
  - Proof-family changes may force new protocol pCIDs, which could fragment comparable adoption metadata conventions across relationships.
- Open questions:
  - Does binding proof-family semantics to the payload protocol pCID make relationship-scoped adoption metadata harder for generic auditors to compare across peers?
  - Should any minimal outer hint exist for proof-family legibility without turning the envelope into a profile bundle?
  - How should payload protocols express relationship-scoped adoption promises while preserving the envelope's clean separation of transport evidence from higher-layer trust policy?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `group-session-freeze-promise-post-freeze-mutation-request`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/group-session-freeze-promise-post-freeze-mutation-request/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=3 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=5 risk_penalty=2
- Fitness: raw=42.00 normalized_0_100=84.00 confidence_0_1=0.87
- Rationale: Strong Promise-Theory-clean envelope specimen. For this scenario it does not itself define group-session freeze lineage rules, but its explicit pCID binding and immutable signed [pCID,payload] bytes support the safe answer that a post-freeze breaking change should travel as a new specimen or new pCID rather than mutate old history. Scenario fit is therefore moderate rather than high because the crucial freeze semantics live in the payload protocol, not the envelope.
- Strengths:
  - Very clear layer boundary: envelope binds protocol identity and payload bytes without pretending to solve higher-layer session semantics.
  - Excellent promise-first wording about sender evidence, local trust, and receiver autonomy.
  - Small deterministic three-slot artifact is durable and plausible for long-term storage and transport.
  - ... 1 more
- Weaknesses:
  - Does not by itself specify whether a post-freeze group-session change is a superseding pCID, a new specimen, or an amendment relation.
  - Outer auditability of proof semantics is weaker than designs with an explicit outer proof-profile selector.
  - Failure handling is mostly limited to unknown pCID handling; lineage conflict resolution is delegated upward.
- Risks:
  - Proof-family semantics are hidden behind the payload protocol pCID, which may reduce quick outer-layer audit clarity during lineage disputes.
  - If small proof-rule changes require minting new protocol pCIDs, evolution pressure could fragment protocol lineages.
  - A weak payload protocol could still mishandle post-freeze mutation even though the envelope preserves byte-level history correctly.
- Open questions:
  - For a frozen group-session protocol, should any post-freeze semantic change require a new pCID rather than an amendment under the old pCID?
  - Does protocol-owned proof semantics make later audit of freeze-era evidence too dependent on recovering the exact pCID spec bytes?
  - Should payload protocols carry an explicit supersedes/amends relation so post-freeze lineage stays legible without mutating historical specimens?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `promise-economy-spectrum-transferable-promise-token`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/promise-economy-spectrum-transferable-promise-token/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=2 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=5 risk_penalty=2
- Fitness: raw=41.00 normalized_0_100=82.00 confidence_0_1=0.88
- Rationale: Strong PromiseGrid envelope specimen: tiny, explicit, and promise-first at the claimed layer. For this transferable-promise-token scenario, however, most substantive token semantics are intentionally delegated to the payload protocol, so scenario fit is limited but not invalid.
- Strengths:
  - Very clear envelope/payload boundary.
  - Small deterministic three-slot artifact with good durability properties.
  - Explicit local handling of unknown pCIDs and blind carriage.
  - ... 1 more
- Weaknesses:
  - Does not itself provide provenance, permission, double-spend prevention, or transferable refusal semantics.
  - Failure handling is mostly limited to unknown-protocol preservation rather than token-specific dispute paths.
  - Auditability for transferable-token use depends heavily on payload protocol details.
- Risks:
  - A token protocol could smuggle authority-heavy or obligation-heavy semantics inside the payload while the clean outer envelope makes that easy to overlook.
  - Generic peers cannot fully audit provenance or transfer semantics without understanding the payload protocol named by the pCID.
  - Binding proof-family rules to the payload protocol may make long-term evolution and cross-protocol comparison less legible.
- Open questions:
  - Is generic audit clarity harmed when proof-family semantics are hidden behind the payload protocol pCID rather than an outer proof selector?
  - Would transferable-token payload protocols need frequent new pCIDs when proof or delegation rules evolve?
  - Can blinded carriage of unknown pCIDs preserve enough evidence for later token dispute audits?
- Authority boundary: Envelope-layer evidence only; payload protocol owns transferable-token provenance, permissions, double-spend prevention, and refusal semantics.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `promise-economy-spectrum-permissioned-capability-token`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/promise-economy-spectrum-permissioned-capability-token/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=3 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=5 risk_penalty=2
- Fitness: raw=42.00 normalized_0_100=84.00 confidence_0_1=0.85
- Rationale: Strong Promise-Theory-clean envelope specimen. It gives a small, explicit, signed carrier that could transport permissioned capability tokens well, but the actual permission, non-transferability, and redemption semantics are intentionally left to the payload protocol, so scenario fit is only partial at this layer.
- Strengths:
  - Very clear layer boundary: outer envelope only binds protocol name and payload bytes.
  - Explicit PT wording: current-sender evidence for a scoped promise about payload shape.
  - Deterministic small artifact: three-slot CBOR array, mandatory proof slot, canonical signed prefix.
  - ... 1 more
- Weaknesses:
  - Does not itself express permission scope, audience restriction, or redemption rules for the capability scenario.
  - Generic auditors cannot infer proof-family details from the envelope alone without resolving the pCID spec.
  - Failure handling beyond unknown-pCID carriage is mostly delegated to payload protocols.
- Risks:
  - Proof semantics hidden behind the same pCID reduce generic audit clarity for peers that do not already know the protocol.
  - Capability-token safety depends heavily on payload-level signer binding, freshness, and revocation rules that this envelope does not standardize.
  - If proof-family changes require new protocol pCIDs, long-lived capability ecosystems may face migration churn.
- Open questions:
  - How should a payload protocol bind a capability to Bob so forwarding the same envelope does not imply transferability?
  - What payload-level evidence handles freshness, revocation, or one-time redemption without introducing central authority assumptions?
  - Does proof-family evolution under one pCID force too-frequent protocol re-minting for long-lived capability schemes?
- Authority boundary: Envelope-layer evidence only; payload protocol must define capability issuance, non-transferability, redemption, and local trust decisions.

### `SIM-pobod-grid-envelope-outer-promise-nested-signed-payload` x `promise-economy-spectrum-peer-local-assessment-only`

- Result path: `results/SIM-pobod-grid-envelope-outer-promise-nested-signed-payload/promise-economy-spectrum-peer-local-assessment-only/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=3 promisegrid_alignment=5 auditability=4 evolution_safety=5 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=4 risk_penalty=1
- Fitness: raw=44.00 normalized_0_100=88.00 confidence_0_1=0.84
- Rationale: Strong PT-clean envelope design. It cleanly supports peer-local assessment by avoiding token requirements and by making trust in structure and signatures explicitly local, but it only indirectly addresses the scenario's higher-layer pull/keep/advertise behavior.
- Strengths:
  - Explicit outer-layer promise fixes the implicit-authorship weakness of the predecessor design.
  - Clear separation between envelope parsing semantics and inner payload meaning supports local trust assessment.
  - Small pCID-shaped envelope with nested signed payload is auditable and evolution-friendly.
  - ... 1 more
- Weaknesses:
  - Does not itself model the local observation records that drive pull/keep/advertise decisions in this scenario.
  - Failure handling for corrupt or stale storage evidence is mostly delegated to the nested payload protocol.
  - Scenario fit is partial because the design is intentionally envelope-layer, not a full promise-economy protocol.
- Risks:
  - Readers may still overinterpret the nested signature as global authority rather than local evidence.
  - Scenario-specific local reputation or storage-observation records are outside the envelope and must be supplied by the payload protocol.
- Open questions:
  - How should peer-local observations about kept or corrupt storage promises be represented in higher-layer payloads carried by this envelope?
  - What minimal payload conventions are needed so pull/keep/advertise decisions remain locally auditable without adding token semantics to the envelope?
- Authority boundary: Evidence only; evaluates envelope-layer promise semantics, not higher-layer promise accounting or token economics.

### `SIM-pobod-grid-envelope-outer-promise-nested-signed-payload` x `promise-accounting-records-kept-storage-promise`

- Result path: `results/SIM-pobod-grid-envelope-outer-promise-nested-signed-payload/promise-accounting-records-kept-storage-promise/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=3 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 risk_penalty=1
- Fitness: raw=41.00 normalized_0_100=82.00 confidence_0_1=0.87
- Rationale: Strong PT-clean envelope design with explicit outer-layer promise, local trust, and small durable artifacts. It fits this storage-promise-accounting scenario only partially because the scenario asks for Bob's local record of a kept storage promise, while this simulation intentionally stops at envelope/evidence carriage and delegates storage semantics and accounting vocabulary to the nested payload protocol.
- Strengths:
  - Explicit promise-first outer-layer wording.
  - Very clear layer boundary between envelope and payload semantics.
  - Good auditability from signed nested payload carriage.
  - ... 1 more
- Weaknesses:
  - Does not by itself define Bob's local kept-storage accounting record.
  - Failure and dispute handling for storage claims are mostly delegated upward.
  - Scenario-specific fit is limited because storage semantics are outside the claimed layer.
- Risks:
  - Reviewers may over-credit the envelope as if it solved storage-promise accounting by itself.
  - Without a standard payload receipt vocabulary, different peers may record incomparable kept-storage evidence.
  - Nested-signature evidence can still be misinterpreted as stronger authority than it actually provides.
- Open questions:
  - What exact nested payload shape would let Bob record a kept-storage observation in a durable promise-accounting form?
  - How should Bob distinguish evidence of chunk delivery from evidence of a longer-term keep promise at the payload layer?
  - What minimal receipt or observation vocabulary should be standardized above the envelope for storage scenarios?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-pobod-grid-envelope-outer-promise-nested-signed-payload` x `promise-accounting-records-refused-service`

- Result path: `results/SIM-pobod-grid-envelope-outer-promise-nested-signed-payload/promise-accounting-records-refused-service/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=2 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 risk_penalty=1
- Fitness: raw=39.00 normalized_0_100=78.00 confidence_0_1=0.87
- Rationale: A strong Promise-Theory-clean envelope design: it makes the outer-layer promise explicit, keeps trust local, and avoids unsigned authority claims. For this scenario, however, it mainly provides a transport/evidence wrapper and does not itself define promise-accounting semantics that distinguish refusal from failure, corruption, or timeout. That limits scenario fit without making the design PT-invalid.
- Strengths:
  - Explicit outer-layer promise wording is PT-correct and locally assessable.
  - Clear layer boundary: envelope semantics are separated from payload meaning and accounting policy.
  - Signed nested payload supports durable audit evidence without relying on transport authority.
  - ... 1 more
- Weaknesses:
  - Does not by itself encode refusal-specific semantics or accounting records.
  - Scenario asks for differentiation among refusal, failure, corruption, and timeout, which is mostly delegated upward.
  - Auditability for refusal outcomes depends on payload-level conventions not specified here.
- Risks:
  - Implementers may overread the signed nested payload as solving refusal accounting when the distinction actually depends on a higher-layer payload protocol.
  - Different payload protocols may encode refusal inconsistently, reducing cross-peer audit comparability.
  - Absent a shared refusal record convention, later auditors may still struggle to distinguish honest refusal from silence.
- Open questions:
  - What higher-layer payload protocol cleanly encodes an explicit refusal so peers can record it distinctly from timeout, corruption, or execution failure?
  - Should the outer envelope reserve any minimal refusal/error discriminator, or should all such semantics remain entirely in the nested payload protocol?
  - What local promise-accounting record format would let Carol later audit that Alice honestly refused service rather than simply failed to respond?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-pobod-grid-envelope-outer-promise-nested-signed-payload` x `conditional-release-geofencing-onward-restraint-chain`

- Result path: `results/SIM-pobod-grid-envelope-outer-promise-nested-signed-payload/conditional-release-geofencing-onward-restraint-chain/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=2 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=5 risk_penalty=2
- Fitness: raw=41.00 normalized_0_100=82.00 confidence_0_1=0.84
- Rationale: Strong PT-clean envelope-layer design. It explicitly states the sender's outer-layer promise about byte shape and treats nested signatures as evidence, preserving local trust and clear layer boundaries. For this scenario, however, the recursive onward-restraint chain is mostly delegated to the payload protocol, so fit is limited but not invalid at the claimed layer.
- Strengths:
  - Explicit promise-first restatement of the outer envelope semantics.
  - Clear separation between envelope parsing/authorship evidence and higher-layer payload meaning.
  - Small signed-payload carrier is simple, durable, and plausible for constrained devices.
  - ... 1 more
- Weaknesses:
  - Does not itself model the recursive promise graph central to this scenario.
  - Limited direct handling of breach, refusal, or onward-restraint failure evidence.
  - Scenario semantics depend heavily on a separate payload family that is not defined here.
- Risks:
  - Users may over-read the envelope as enforcing conditional release when it only carries evidence for higher-layer policy claims.
  - Without a well-specified payload protocol, onward-restraint chains may become inconsistent or unauditable across hops.
  - Forwarding-chain semantics such as geofencing, expiry, and breach evidence remain underspecified here.
- Open questions:
  - Which payload protocol would express Bob's onward-restraint promise to Carol without collapsing local trust into a global policy graph?
  - How would Alice or Carol record and audit breaches of the forwarding-chain condition using only local evidence?
  - How are geofencing predicates, expiry, and evolution of policy terms named and signed inside the nested payload?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-pobod-grid-envelope-outer-promise-nested-signed-payload` x `peer-adoption-metadata-relationship-scoped-adoption`

- Result path: `results/SIM-pobod-grid-envelope-outer-promise-nested-signed-payload/peer-adoption-metadata-relationship-scoped-adoption/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=3 promisegrid_alignment=5 auditability=4 evolution_safety=5 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 risk_penalty=1
- Fitness: raw=42.00 normalized_0_100=84.00 confidence_0_1=0.84
- Rationale: Strong PT-clean envelope design: it makes the outer promise explicit, keeps trust local, and preserves a clean layer boundary. For this scenario it is only a partial fit, because relationship-scoped adoption behavior lives in the nested payload protocol rather than the envelope itself.
- Strengths:
  - Explicit outer-layer promise wording matches Promise Theory well.
  - Clear separation between envelope parsing promises and inner payload meaning.
  - Small signed-payload structure is auditable and evolution-friendly.
  - ... 1 more
- Weaknesses:
  - Relationship-scoped adoption is not expressed directly at the envelope layer.
  - Failure semantics beyond local refusal and skepticism are not deeply modeled.
  - The design depends on payload protocols to prevent ambiguity between peer identity and relationship context.
- Risks:
  - Scenario-specific relationship scoping may be under-specified if the payload protocol does not bind claims to a concrete peer relationship.
  - Auditors may over-attribute higher-layer adoption semantics to the envelope because the envelope only transports and frames them.
  - If relationship context leaks into unsigned transport assumptions, local trust can collapse back into implicit authority.
- Open questions:
  - How should the nested payload name relationship context durably without introducing central identity assumptions?
  - Should relationship-scoped adoption require recipient-bound signatures or separate payload instances per peer relationship?
  - What local evidence best distinguishes legitimate per-relationship variation from equivocation across peers?
- Authority boundary: Evidence only at the envelope layer; relationship-scoped adoption semantics remain a payload-protocol and local-trust matter.

### `SIM-pobod-grid-envelope-outer-promise-nested-signed-payload` x `group-session-freeze-promise-post-freeze-mutation-request`

- Result path: `results/SIM-pobod-grid-envelope-outer-promise-nested-signed-payload/group-session-freeze-promise-post-freeze-mutation-request/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=3 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=4 risk_penalty=1
- Fitness: raw=42.00 normalized_0_100=84.00 confidence_0_1=0.85
- Rationale: Strong PT-clean envelope design with explicit outer-layer promises, good local trust framing, and clear delegation to payload protocols. For this freeze-mutation scenario it helps preserve signed evidence and evolution boundaries, but it does not itself decide whether a post-freeze change is a new specimen, superseding pCID, or amendment; that remains payload/application work.
- Strengths:
  - Explicit promise-first outer envelope semantics.
  - Very clear layer boundary between envelope parsing and payload meaning.
  - Nested signed payload improves local auditability and durable evidence carriage.
  - ... 2 more
- Weaknesses:
  - Does not by itself answer the scenario's core lineage decision after freeze.
  - Failure handling for conflicting post-freeze claims is only partial at this layer.
  - Audit conclusions depend on payload protocol clarity, not envelope semantics alone.
- Risks:
  - Consumers may overread the envelope as settling freeze semantics that are actually delegated to the payload protocol.
  - Competing post-freeze payloads may still be hard to reconcile without explicit lineage rules above the envelope layer.
  - If outer and inner semantics are poorly documented, auditors may confuse structure validity with legitimacy of mutation.
- Open questions:
  - Which payload-level protocol names freeze state, supersession, or amendment lineage for post-freeze changes?
  - How should receivers distinguish a legitimate superseding specimen from an attempted mutation of frozen history?
  - What minimal local audit record should peers keep alongside the nested signed payload to compare competing post-freeze claims?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-pobod-grid-envelope-outer-promise-nested-signed-payload` x `promise-economy-spectrum-transferable-promise-token`

- Result path: `results/SIM-pobod-grid-envelope-outer-promise-nested-signed-payload/promise-economy-spectrum-transferable-promise-token/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 risk_penalty=2
- Fitness: raw=37.00 normalized_0_100=74.00 confidence_0_1=0.82
- Rationale: Strong PT-clean envelope design with explicit outer-layer promise, good auditability, and durable small-artifact shape. It fits this transferable-token scenario only partially because provenance, permission, double-spend prevention, and refusal semantics live above the envelope layer.
- Strengths:
  - Explicit promise-first outer envelope semantics.
  - Clear separation between envelope evidence and payload-level token logic.
  - Signed nested payload supports later local audit.
  - ... 1 more
- Weaknesses:
  - Does not itself model transfer authorization or redemption semantics.
  - No native anti-double-spend or revocation/accounting mechanism at this layer.
  - Scenario-specific behavior depends heavily on unspecified payload protocol choices.
- Risks:
  - Operators may overread the signed envelope as a complete transferable-token protocol.
  - Without a higher-layer spend/accounting protocol, double-spend disputes remain unresolved.
  - Different payload protocols may encode transfer semantics incompatibly, weakening interoperability.
- Open questions:
  - What payload protocol records transfer provenance and anti-double-spend evidence without requiring global state?
  - How are refusal and non-acceptance expressed so receipt of a transferred token is not misread as obligation?
  - What local evidence should Carol retain to judge whether Alice's original promise remains meaningful after transfer?
- Authority boundary: Evidence only; evaluates an envelope-layer design and leaves transferable-token accounting, spend control, and refusal policy to higher-layer payload protocols and local trust decisions.

### `SIM-pobod-grid-envelope-outer-promise-nested-signed-payload` x `promise-economy-spectrum-permissioned-capability-token`

- Result path: `results/SIM-pobod-grid-envelope-outer-promise-nested-signed-payload/promise-economy-spectrum-permissioned-capability-token/openai-gpt-5.4-medium/20260524-190507.json`
- Scores: scenario_fit=3 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 risk_penalty=1
- Fitness: raw=41.00 normalized_0_100=82.00 confidence_0_1=0.84
- Rationale: Strong PT-clean envelope-layer design. It explicitly states the sender's outer-layer promise about nested bytes and leaves capability-token meaning to the signed payload protocol, which is the right boundary for this simulation. It fits the scenario only partially because permission and redemption semantics are mostly delegated inward, but that delegation is explicit and locally auditable rather than authority-led.
- Strengths:
  - Excellent layer-boundary clarity between envelope parsing promises and payload semantics.
  - Very strong Promise-Theory vocabulary at the claimed layer.
  - Small signed-evidence-oriented shape is durable, auditable, and plausible for constrained devices.
- Weaknesses:
  - Does not by itself define capability-token permission or redemption behavior.
  - Failure handling for revocation, expiry, and replay is mostly deferred to payload design.
  - Scenario fit is limited because marketplace-avoidance semantics are not expressed at the envelope layer.
- Risks:
  - Capability semantics may be underspecified if payload protocols do not clearly define non-transferability and redemption scope.
  - Different payload protocols could produce inconsistent local interpretations of permission without stronger conventions.
  - Replay or stale-token risks remain unless the nested payload protocol carries freshness and revocation evidence.
- Open questions:
  - What payload-level promise shape best expresses non-transferable capability scope and redemption semantics without slipping into exchange-rate framing?
  - How should expiry, revocation, and replay resistance be evidenced inside the nested signed payload for long-lived offline audit?
  - What minimal payload conventions let different peers locally assess capability redemption while preserving the small outer envelope?
- Authority boundary: Evidence only; envelope-layer design; capability semantics and trust updates remain payload-protocol-local.

## Required JSON Shape

{"child_id":"SIM-jogoh-child-descriptive-design-slug","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
