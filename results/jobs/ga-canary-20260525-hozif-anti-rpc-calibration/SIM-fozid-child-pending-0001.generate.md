# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Choose a descriptive `child_id` that starts with `SIM-fozid-child-` and ends with a kebab-case design slug. Do not use generic `ga-child`, `pending`, or ordinal-only names.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

## Design Guardrails

- Treat `pCID` as Protocol CID: the pCID-named protocol spec may define payload shape, signature/proof encoding, refusal evidence, freeze-successor records, or capability promise-token records.
- For base-envelope children, avoid selector-shopping stacks such as `env_pCID`/`sig_pCID`/`payload_pCID`, generic `claim_header` or claim-card layers, generic claim cards, and universal `statement_capsule` wrappers.
- Do not ban higher-layer payload protocols from defining their own promise-accounting records. Put signed refusal versus silence/timeout, exact-byte local evidence, freeze successor records, transfer semantics, and capability-as-promise-token behavior inside the pCID-selected payload/specimen layer unless the scenario explicitly asks to test an envelope negative control.
- Keep outer signatures scoped to the current sender's own promise; no agent promises for another agent, and every receiver keeps local trust assessments.

- Run group ID: `ga-canary-20260525-hozif-anti-rpc-calibration`
- Planned child ID prefix: `SIM-fozid-child-`
- Temporary child ID: `SIM-fozid-child-pending-0001`
- Temporary child path: `proposals/ga-canary-20260525-hozif-anti-rpc-calibration/simulations/SIM-fozid-child-pending-0001/`
- Operation: `breed`
- Parent IDs: `SIM-fovip-kernel-promise-boundary-port-contract, SIM-dalor-grid-envelope-protocol-owned-signature-slot`

## Scenario Sample

- `kernel-porting-boundary` at `scenarios/kernel-porting-boundary/kernel-porting-boundary.md`
- `promise-accounting-records-kept-storage-promise` at `scenarios/promise-accounting-records-kept-storage-promise/promise-accounting-records-kept-storage-promise.md`
- `promise-accounting-records-refused-service` at `scenarios/promise-accounting-records-refused-service/promise-accounting-records-refused-service.md`

## Scenario Pressure

### `scenarios/kernel-porting-boundary/kernel-porting-boundary.md`

```markdown
# Kernel Porting Boundary

## Scenario ID

kernel-porting-boundary

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-vitih`, `FB-mulum`, and `FB-potin`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-davod`

## Purpose

Exercise the developer-facing boundary for a first real PromiseGrid port while
the kernel/runtime terminology remains unsettled. Source: `DI-ragaz`;
`DI-fidot`.

## Setup

Alice wants to port PromiseGrid infrastructure to a new host environment. Bob
offers a native service, Carol offers a browser/WASM adapter, Dave offers an
MCU/header-only profile, Ellen offers split local services, and Mallory claims
that copying the wire-lab harness or one microkernel shape is the porting
target. The available specs are drafts, and only some future frozen pCIDs will
become obligations.

## Stimulus

Alice writes a first porting plan and a kernel implementation promise record.
The plan must say what it implements now, which pCID-selected messages it
exposes, what draft evidence it follows, what it refuses or cannot promise, what
host assumptions it depends on, what evidence it records, and what it defers
until `DR-davod` decides the guide-facing boundary. Source: `DI-ripuz`.

## Expected Pressure

The candidate design must separate harness apparatus from porting target,
identify which binding/session/message/CAS/runtime obligations are provisional,
and preserve a clear path to future frozen-spec implementation promises.

It must also show whether:

- app/kernel operations are pCID-selected `grid([42(pCID), payload, ...])`
  messages, even when local APIs wrap them;
- storage, compute, network send/receive, key use, device access, lifecycle,
  dispatch, refusal, receipt, evidence, namespace, reference, and checkpoint
  operations each name their promiser;
- host/runtime assumptions are separate from PromiseGrid promises;
- unsupported pCIDs and unsupported roles are direct refusals or non-promises;
- exact bytes are preserved where needed for proof, replay, unsupported-pCID
  carriage, or broken-promise evidence;
- voluntary group namespaces are reciprocal promises, not global truth;
- CID-rooted references let Alice share a resource that Bob maps into Bob's own
  local view;
- file/resource current state can be reconstructed as a checkpoint over a
  selected promise-log frontier.

## Scenario Variants

- **Native node:** Bob's service promises storage, dispatch, networking, keys,
  lifecycle, and evidence, but must name every pCID it supports and every role it
  does not promise.
- **Browser/WASM:** Carol's adapter depends on browser storage, network, clocks,
  and lifecycle. Carol can promise adapter behavior, but not that the browser
  will keep promises the browser has not made.
- **Mobile sandbox:** Dave promises work only while the OS offers foreground or
  background execution. Unavailable background work must be recorded as an
  unavailable promise, refusal, or host assumption rather than hidden success.
- **MCU/header-only:** Erin supports one actuator pCID, one bounded evidence
  store, and no general namespace service. The port is credible only if it says
  exactly what it cannot promise.
- **Split local services:** Frank separates dispatch, storage, keys, networking,
  and evidence among local promisers. The record must show which service makes
  each promise and how evidence moves between them.
- **Voluntary namespace:** Alice, Bob, and Carol maintain `/project/report` by
  reciprocal namespace promises. Mallory's lookalike namespace is rejected unless
  a local agent trusts Mallory's promise history.
- **CID-rooted reference:** Alice sends Bob a reference rooted at a CID with
  pCID, selector/path, frontier, promise body, and evidence. Bob chooses whether
  and where to mount it.
- **Checkpointed resource:** Grace reconstructs a file from old pCID specs,
  promise/event logs, and a selected frontier. A different branch may produce a
  different current file because it selects a different promise history.

## Scenario-Specific Evaluation Questions

- Should the guide say kernel, runtime, dispatcher, handler host, or library?
- What is the minimum viable porting target before final freeze?
- Which K1-K5 ingress, feed, CAS, session, and app-layer details should remain
  blocked versus provisional orientation?
- Does the candidate preserve local trust, autonomous promisers, and
  make/break evidence?
- Does the candidate avoid global permission, global namespace authority, and
  universal process-shape assumptions?
- Are local APIs faithful adapters over pCID-selected messages, or do they hide
  the promise boundary?
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

## Parent Simulation Documents

### `simulations/SIM-fovip-kernel-promise-boundary-port-contract/README.md`

```markdown
# SIM-fovip: Kernel promise boundary port contract

This simulation is the active evidence surface for the `DR-davod` kernel design
packet. It follows `TE-mazop`, which found that `TE-jimar` was not enough to
close the kernel-developer porting boundary, and now incorporates the follow-on
`TE-pudiv`, `TE-dunas`, and `TE-gakoh` questions. Source: `DI-funaf`;
`DI-fidot`.

## Question

What minimum PromiseGrid kernel implementation promises can be claimed across
rich native nodes, browser/WASM hosts, mobile sandboxes, MCU/header-only ports,
and split local service graphs without pretending that one process shape,
namespace, API, or prior-art pattern is universal?

## Candidate specimen

The specimen under test is a kernel implementation promise record: a local
implementation promises its own behavior by publishing:

- profile name and runtime class;
- supported pCIDs and unsupported-pCID behavior;
- app-facing promises for storage, compute, network send/receive, key use,
  device access, lifecycle, pCID dispatch, refusal, receipt, evidence,
  namespace, reference, and checkpoint behavior;
- host/runtime assumptions that the port depends on but does not itself promise;
- explicitly unsupported features;
- exact evidence records for kept, refused, unavailable, and broken promises;
- adapter promises when local APIs wrap pCID-selected grid messages;
- voluntary namespace promises when the port projects group namespaces;
- CID-rooted promise-bound reference behavior for cross-agent path sharing.

## Runtime classes

- **Native node:** Bob runs a local service with storage, network, keys, feed,
  CAS, dispatch, lifecycle, and evidence roles.
- **Browser/WASM:** Alice runs inside a host that owns storage, network, clocks,
  keys, and lifecycle.
- **Mobile sandbox:** Dave can run only while the OS permits background work and
  network access.
- **MCU/header-only:** Carol supports one or two pCIDs, bounded evidence, and a
  hardware/device promise.
- **Split local services:** Ellen separates dispatch, storage, keys, networking,
  and evidence into local services with separate promises.

## Basic principles under test

- Kernel is a role/profile set, not a ruler.
- Everything useful is a promise: app/kernel operations, resources, namespaces,
  references, and kernel implementation promise records all help agents make or
  evaluate promises.
- The app/kernel boundary is a promise boundary; exposed operations are
  pCID-selected `grid([42(pCID), payload, ...])` messages, even when local APIs
  provide ergonomic adapters.
- A kernel implementation promise record is not a global certificate. Alice,
  Bob, Carol, and later agents evaluate the record and make/break history
  locally.
- Host assumptions are not implementation promises unless the host is also an
  explicit promiser.
- Voluntary group namespaces may exist inside trust relationships, but imposed
  universal namespaces are rejected.
- File-like resources are promise-log projections or checkpoints, not evidence
  that PromiseGrid is filesystem-first.

## Evidence axes

The simulation should let reviewers ask whether the candidate:

- names the local promiser for each storage, compute, network, key, device,
  lifecycle, dispatch, namespace, reference, and evidence promise;
- maps every exposed app/kernel operation to a pCID-selected message or explains
  why the operation is outside the PromiseGrid boundary;
- records exact bytes when needed for proof, replay, unsupported-pCID carriage,
  or broken-promise evidence;
- states host/runtime assumptions separately from the port's own promises;
- names unsupported pCIDs and unsupported roles directly;
- keeps trust local and avoids global permission, namespace, conformance, or
  policy authority;
- treats V, Amoeba, Plan 9, and Hurd as pattern pressure, not imported design
  authority;
- supports voluntary group namespaces and CID-rooted promise-bound references
  without treating Alice's local path as global truth;
- represents file/resource state as checkpoints or projections over selected
  promise/event log frontiers.

## Boundaries

This simulation does not close `DR-davod` and does not define a final
PromiseGrid kernel API. It tests whether kernel implementation promises give
guide writers enough evidence to discuss kernel developers without promising a
daemon, microkernel, browser host, mobile runtime, MCU library, namespace
protocol, or SDK as the single correct implementation shape.

The current envelope shape `grid([42(pCID), payload, ...])` is input evidence,
not a reopened decision.

## Related evidence

- `docs/research/DN-lujad-promisegrid-kernel-role-profile.md`
- `docs/thought-experiments/TE-jimar-kernel-runtime-portability-boundary.md`
- `docs/thought-experiments/TE-mazop-kernel-promise-boundary-and-minimum-port-contract.md`
- `docs/thought-experiments/TE-pudiv-app-kernel-grid-message-boundary.md`
- `docs/thought-experiments/TE-dunas-prior-art-influence-on-promisegrid-kernel.md`
- `docs/thought-experiments/TE-gakoh-local-views-over-promise-event-hypergraph.md`
- `DR/DR-davod-promisegrid-kernel-dev-porting-boundary.md`
- `scenarios/kernel-porting-boundary/kernel-porting-boundary.md`
- `simulations/SIM-funas-kernel-porting-boundary/`
```

### `simulations/SIM-fovip-kernel-promise-boundary-port-contract/QUESTION.md`

```markdown
# Question

Which minimum PromiseGrid kernel implementation promises can guide writers
describe without turning one runtime shape into a false universal kernel, and
without treating kernel, host, namespace, app API, or prior-art patterns as
external authority?

The answer must be concrete enough to test app-facing promises, host/runtime
assumptions, unsupported features, pCID coverage, exact-byte evidence records,
broken-promise handling, app/kernel pCID messages, voluntary group namespaces,
CID-rooted promise-bound references, and file/resource checkpoints over promise
logs across native, browser/WASM, mobile, MCU, and split local-service
deployments. Source: `DI-funaf`; `DI-fidot`; `DR-davod`; `TE-mazop`.
```

### `simulations/SIM-fovip-kernel-promise-boundary-port-contract/protocols/kernel-port.d/specs/kernel-port-contract-draft.md`

```markdown
# Kernel implementation promise record draft

This draft is a simulation-local specimen for `SIM-fovip`. It is not a frozen
PromiseGrid protocol spec. Its purpose is to make the `DR-davod` question
testable by forcing a port to publish the promises, assumptions, unsupported
features, and evidence records that make a first implementation credible. Source:
`DI-funaf`; `DI-fidot`.

## Role

A kernel implementation promise record says what a PromiseGrid implementation
promises to local apps and operators, what it depends on from the host runtime,
and what it explicitly does not promise.

The record is not a global certificate, permission, or authority. Each receiving
agent still assesses the record and the port's make/break history locally.

## Record shape

```text
kernel_implementation_promise_record = [
  record_pcid,
  port_identity,
  profile,
  supported_pcids,
  app_facing_promises,
  host_assumptions,
  unsupported_features,
  evidence_policy,
  adapter_promises,
  namespace_promises,
  reference_promises
]
```

## Fields

- `record_pcid` is the Protocol CID for this promise-record shape.
- `port_identity` names the local implementation or agent making the promises.
- `profile` names the runtime class: native node, browser/WASM, mobile sandbox,
  MCU/header-only, split local services, or another pCID-defined profile.
- `supported_pcids` lists the pCID-selected protocols the implementation promises to
  parse, dispatch, validate, or preserve.
- `app_facing_promises` states what the implementation promises for storage,
  compute, network send/receive, key use, device access, lifecycle, pCID
  dispatch, refusal, receipt, evidence recording, namespace projection,
  reference resolution, and resource checkpoint behavior.
- `host_assumptions` states what the port depends on from a browser, OS, mobile
  sandbox, language runtime, hardware platform, or local service graph.
- `unsupported_features` states what the port refuses or cannot perform.
- `evidence_policy` states what exact bytes and local records the implementation promises
  to keep for kept, refused, unavailable, and broken promises.
- `adapter_promises` states which local APIs wrap which pCID-selected messages
  and what evidence the adapter records.
- `namespace_promises` states whether the port projects voluntary group
  namespaces and which promisers maintain those namespace frontiers.
- `reference_promises` states how the port handles CID-rooted promise-bound
  references and local path mounting.

## Promise rules

- An implementation promises only its own behavior.
- An implementation may cite host/runtime assumptions, but it does not promise
  that the host will keep them unless the host is also an explicit promiser.
- Unsupported features must be named directly. They must not be hidden behind a
  generic "partial implementation" label.
- Evidence records are local. They help Alice, Bob, Carol, and future agents
  update their own trust judgments; they are not a global trust authority.
- Local APIs are adapters. If an API call crosses a PromiseGrid promise
  boundary, the record must identify the corresponding pCID-selected message or
  state that the operation is outside the PromiseGrid boundary.
- Voluntary group namespaces are promises among agents. The record must not
  describe a namespace as universal truth.
- File-like resources are projections over promises, logs, and checkpoints. The
  record must say what evidence is preserved for the selected frontier.

## Minimum credible first port

A first credible port is allowed to be small. It must still publish:

- at least one supported pCID or a bounded exact-byte carriage profile;
- clear unsupported-pCID behavior;
- app-facing promises for every operation it exposes;
- host/runtime assumptions for every operation it delegates;
- evidence records for kept, refused, unavailable, and broken promises;
- adapter, namespace, reference, and checkpoint promises where the port exposes
  those surfaces;
- an implementation promise record that can be compared with later behavior.

## Scenario pressure

The same record shape must be tested against:

- a native node with broad local services;
- a browser/WASM host with delegated storage, network, key, and lifecycle
  behavior;
- a mobile sandbox with restricted background execution;
- an MCU/header-only port with one pCID and bounded evidence;
- a split local service graph with multiple local promisers;
- a voluntary group namespace maintained by Alice, Bob, and Carol;
- a CID-rooted promise-bound reference from Alice that Bob mounts locally;
- a file/resource checkpoint reconstructed from a selected promise-log frontier.
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

### `SIM-fovip-kernel-promise-boundary-port-contract` x `kernel-porting-boundary`

- Result path: `results/SIM-fovip-kernel-promise-boundary-port-contract/kernel-porting-boundary/openai-gpt-5.4-medium/20260526-003243.json`
- Scores: scenario_fit=5 promisegrid_alignment=5 auditability=5 evolution_safety=5 layer_boundary_clarity=5 failure_handling=5 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 envelope_discipline=4 kernel_implementation_promises=5 app_protocol_promise_semantics=3 risk_penalty=1
- Fitness: raw=59.00 normalized_0_100=90.77 confidence_0_1=0.92
- Rationale: Strong kernel-layer Promise Theory fit. The specimen squarely addresses the porting-boundary scenario by making a local implementation publish its own supported pCIDs, app-facing promises, host assumptions, refusals, and evidence behavior without treating the kernel as an authority, registry, or universal process shape.
- Strengths:
  - Directly targets the kernel porting boundary and separates harness pressure from the implementation target.
  - Excellent local-promiser accounting across storage, compute, networking, keys, device, lifecycle, dispatch, namespace, reference, and evidence surfaces.
  - Clear unsupported-feature and unsupported-pCID handling with explicit refusal/non-promise framing.
  - ... 3 more
- Weaknesses:
  - Still a draft specimen rather than a frozen protocol, so some terms and boundaries remain provisional.
  - The record is wide enough that profile-specific trimming may be needed for long-term simplicity on very small ports.
  - Higher-layer app protocol semantics are mostly delegated rather than fully modeled here, which is appropriate but leaves some integration details open.
- Risks:
  - The broad record could drift into a conformance-card or profile-catalog style artifact if later text stops emphasizing local promises and observations.
  - The distinction between the promise-record pCID and app/kernel message pCIDs must stay crisp to avoid selector confusion.
  - Split-service deployments may need tighter guidance so promiser naming stays explicit rather than collapsing into a kernel-global abstraction.
- Open questions:
  - Should namespace, reference, and checkpoint clauses stay in one kernel promise record pCID or split into separately deployable payload protocols later?
  - How should split-local-service ports name multiple local promisers without turning the record into a registry-like catalog?
  - What minimum exact-byte evidence retention is mandatory for tiny MCU profiles versus optional for richer hosts?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-fovip-kernel-promise-boundary-port-contract` x `promise-accounting-records-kept-storage-promise`

- Result path: `results/SIM-fovip-kernel-promise-boundary-port-contract/promise-accounting-records-kept-storage-promise/openai-gpt-5.4-medium/20260526-003243.json`
- Scores: scenario_fit=3 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 envelope_discipline=4 kernel_implementation_promises=5 app_protocol_promise_semantics=2 risk_penalty=1
- Fitness: raw=53.00 normalized_0_100=81.54 confidence_0_1=0.87
- Rationale: Strong promise-first kernel-layer design. It explicitly supports local evidence records for kept, refused, unavailable, and broken promises and cleanly separates host assumptions from kernel promises. For this storage-accounting scenario, it provides good kernel support for recording kept storage behavior, but it stops short of defining the higher-layer storage protocol semantics and Bob's trust-update logic in detail.
- Strengths:
  - Excellent local-trust and non-authoritative framing.
  - Clear app/kernel promise boundary via pCID-selected messages.
  - Explicit evidence-policy field for kept/refused/unavailable/broken outcomes.
  - ... 1 more
- Weaknesses:
  - Does not fully specify the application-level kept-storage observation protocol.
  - Bob's future pull/keep/advertise decision update remains implicit rather than modeled.
  - Scenario fit is limited by the candidate's deliberate kernel-layer scope.
- Risks:
  - Scenario pressure may be under-specified if reviewers expect the kernel record itself to define Bob's future advertise-choice semantics.
  - Storage promise outcomes could be split ambiguously between kernel evidence and app protocol payloads unless a sharper boundary rule is added.
- Open questions:
  - Should a kept-storage observation be a kernel evidence record, an app payload kind, or both?
  - What exact local record shape should Bob keep to distinguish served bytes, mere receipt, and longer-term storage trust?
  - How should future pull/keep/advertise choice updates be represented without turning the kernel into a trust authority?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-fovip-kernel-promise-boundary-port-contract` x `promise-accounting-records-refused-service`

- Result path: `results/SIM-fovip-kernel-promise-boundary-port-contract/promise-accounting-records-refused-service/openai-gpt-5.4-medium/20260526-003243.json`
- Scores: scenario_fit=5 promisegrid_alignment=5 auditability=5 evolution_safety=4 layer_boundary_clarity=5 failure_handling=5 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 envelope_discipline=4 kernel_implementation_promises=5 app_protocol_promise_semantics=3 risk_penalty=1
- Fitness: raw=58.00 normalized_0_100=89.23 confidence_0_1=0.88
- Rationale: Strong PT-clean kernel-layer design. It directly addresses the scenario by requiring distinct evidence for refused, unavailable, and broken promises, while keeping trust local and separating host assumptions from kernel promises. Main limits are that it is a kernel promise record rather than a full higher-layer promise-accounting protocol, and the in-record `record_pcid` could blur payload-vs-envelope selection if not carefully specified.
- Strengths:
  - Explicit evidence policy names refused, unavailable, kept, and broken outcomes separately.
  - Very clear app/kernel promise boundary with pCID-selected messages.
  - Strong local-trust, non-authoritative framing.
  - ... 2 more
- Weaknesses:
  - Not itself a full app-level promise-accounting protocol.
  - Record shape is somewhat broad and may drift toward a feature catalog.
  - Envelope discipline is mostly asserted rather than fully exemplified in the record specimen.
- Risks:
  - Catalog-like record growth could become bulky across many runtimes and features.
  - `record_pcid` inside the record may invite redundant selector patterns or pCID confusion if not tightly scoped.
  - Scenario coverage depends on later payload protocols to standardize concrete refusal-record bytes and semantics.
- Open questions:
  - Should refusal evidence be one stable pCID-owned payload family with kind variants for refused/unavailable/broken outcomes, rather than only a field in an implementation record?
  - Is `record_pcid` inside the record payload intentionally distinct from the outer message pCID, or should that redundancy be removed to avoid selector drift?
  - What minimum exact-byte refusal record is required on the smallest MCU/header-only ports?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `kernel-porting-boundary`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/kernel-porting-boundary/openai-gpt-5.4-medium/20260526-003243.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=5 envelope_discipline=3 kernel_implementation_promises=1 app_protocol_promise_semantics=1 risk_penalty=2
- Fitness: raw=45.00 normalized_0_100=69.23 confidence_0_1=0.86
- Rationale: A strong Promise-Theory-clean envelope specimen: small, explicit, and disciplined about pCID meaning, signed bytes, and unknown-pCID behavior. It fits the kernel-porting scenario only partially because it does not define the kernel implementation promise record, local service promisers, refusals, or app/kernel message families beyond the envelope boundary.
- Strengths:
  - Very clear layer boundary: envelope only, with higher-layer promise accounting left to payload protocols.
  - Excellent promise-first wording: sender signature as evidence of the sender's own scoped promise.
  - Small deterministic artifact likely to age well and fit constrained devices.
  - ... 1 more
- Weaknesses:
  - Does not provide the kernel implementation promises the scenario explicitly asks Alice to write.
  - Does not name local promisers for storage, dispatch, networking, evidence, namespace, or lifecycle behavior.
  - Envelope discipline is mixed because the universal mandatory proof slot pushes against the current no-proof-slot-overreach direction.
- Risks:
  - Mandatory universal outer proof slot may be too opinionated for the current envelope direction.
  - Proof-family evolution may couple payload and proof changes too tightly to one protocol pCID.
  - Because the scenario asks for kernel-porting commitments, readers may overread this envelope specimen as a fuller port target than it is.
- Open questions:
  - Does a mandatory outer proof slot overreach the current envelope direction for minimal kernels and MCU-class ports?
  - Will proof-family evolution force too many new protocol pCIDs because proof semantics are owned by the same pCID as payload shape?
  - Is generic audit clarity good enough when outer proof bytes are visible but verification semantics remain protocol-specific?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `promise-accounting-records-kept-storage-promise`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/promise-accounting-records-kept-storage-promise/openai-gpt-5.4-medium/20260526-003243.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=3 evolution_safety=3 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=4 envelope_discipline=3 kernel_implementation_promises=2 app_protocol_promise_semantics=2 risk_penalty=1
- Fitness: raw=46.00 normalized_0_100=70.77 confidence_0_1=0.84
- Rationale: Strong promise-first envelope specimen with clear layer boundaries and disciplined pCID use. It fits the storage-promise-accounting scenario only indirectly because it intentionally leaves storage promises, kept/broken evidence, and future trust updates to the payload protocol. The mandatory outer signature slot is plausible and compact, but it slightly weakens envelope neutrality and makes proof-family audit/evolution more coupled to the protocol pCID.
- Strengths:
  - Very clear envelope-layer scope and explicit delegation of higher-layer promise accounting to payload protocols.
  - Good Promise Theory wording: sender signature as evidence of the sender's own scoped promise.
  - Small deterministic artifact with no selector-shopping stack and correct pCID-as-protocol-spec framing.
  - ... 1 more
- Weaknesses:
  - Does not itself model kept storage promises, local storage observations, or trust updates for this scenario.
  - Generic auditability is reduced because proof-family details are hidden behind the protocol pCID.
  - Kernel implementation promises are largely unspecified.
- Risks:
  - Proof-family semantics are less legible to generic auditors because they are owned by the protocol pCID rather than a separately named outer proof selector.
  - Changes in proof rules may require minting new protocol pCIDs more often than a more decoupled design.
  - A universal mandatory proof slot may be broader than needed for all envelope use cases.
- Open questions:
  - Does coupling proof-family evolution to the protocol pCID force unnecessary pCID churn over time?
  - Is the mandatory outer signature slot too strong for envelope universality compared with more minimal envelopes?
  - Do generic auditors lose too much clarity when proof-family rules are defined only by the protocol named by pCID?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `promise-accounting-records-refused-service`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/promise-accounting-records-refused-service/openai-gpt-5.4-medium/20260526-003243.json`
- Scores: scenario_fit=1 promisegrid_alignment=4 auditability=3 evolution_safety=3 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=5 envelope_discipline=4 kernel_implementation_promises=2 app_protocol_promise_semantics=2 risk_penalty=2
- Fitness: raw=45.00 normalized_0_100=69.23 confidence_0_1=0.89
- Rationale: Strong promise-first envelope specimen with clear layer boundaries and good alignment to local trust and protocol-owned semantics. It scores poorly on this scenario because refused-service accounting is intentionally delegated to the payload protocol, so the design itself does not explain how refusal is recorded distinctly from failure, corruption, or timeout.
- Strengths:
  - Very clear envelope/payload boundary.
  - Promise-first wording is unusually clean for an envelope design.
  - Small deterministic three-slot artifact is easy to implement and durable.
  - ... 1 more
- Weaknesses:
  - Does not itself model promise-accounting records for honest refusal.
  - Kernel implementation promises are mostly implicit rather than spelled out.
  - Proof semantics hidden behind the protocol pCID may be less legible to generic auditors than a more explicit outer proof selector.
- Risks:
  - Scenario pressure may be under-served if higher-layer payload protocols do not standardize explicit refusal records.
  - Coupling payload shape and proof-family evolution under one pCID could reduce audit legibility or force extra protocol churn.
  - Generic peers may preserve bytes but still lack enough context to explain proof semantics without protocol-specific tooling.
- Open questions:
  - Should refused-service, timeout, corruption, and broken-promise records be payload kind variants under one stable higher-layer promise-accounting pCID carried inside this envelope?
  - Does tying proof-family semantics to the same pCID as payload shape create too much pCID churn when signature rules evolve?
  - Is generic auditability adequate when the outer proof slot is explicit but verification rules are only discoverable via the protocol named by pCID?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

## Required JSON Shape

{"child_id":"SIM-fozid-child-descriptive-design-slug","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
