# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Choose a descriptive `child_id` that starts with `SIM-lasuv-child-` and ends with a kebab-case design slug. Do not use generic `ga-child`, `pending`, or ordinal-only names.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

## Design Guardrails

- Treat `pCID` as Protocol CID: the pCID-named protocol spec may define payload shape, signature/proof encoding, refusal evidence, freeze-successor records, or capability promise-token records.
- For base-envelope children, avoid selector-shopping stacks such as `env_pCID`/`sig_pCID`/`payload_pCID`, generic `claim_header` or claim-card layers, generic claim cards, and universal `statement_capsule` wrappers.
- Do not ban higher-layer payload protocols from defining their own promise-accounting records. Put signed refusal versus silence/timeout, exact-byte local evidence, freeze successor records, transfer semantics, and capability-as-promise-token behavior inside the pCID-selected payload/specimen layer unless the scenario explicitly asks to test an envelope negative control.
- Keep outer signatures scoped to the current sender's own promise; no agent promises for another agent, and every receiver keeps local trust assessments.

- Run group ID: `ga-canary-20260525-goban-dalor-fovip-v4`
- Planned child ID prefix: `SIM-lasuv-child-`
- Temporary child ID: `SIM-lasuv-child-pending-0001`
- Temporary child path: `proposals/ga-canary-20260525-goban-dalor-fovip-v4/simulations/SIM-lasuv-child-pending-0001/`
- Operation: `breed`
- Parent IDs: `SIM-fovip-kernel-promise-boundary-port-contract, SIM-dalor-grid-envelope-protocol-owned-signature-slot`

## Scenario Sample

- `kernel-porting-boundary` at `scenarios/kernel-porting-boundary/kernel-porting-boundary.md`
- `promise-accounting-records-kept-storage-promise` at `scenarios/promise-accounting-records-kept-storage-promise/promise-accounting-records-kept-storage-promise.md`
- `promise-accounting-records-refused-service` at `scenarios/promise-accounting-records-refused-service/promise-accounting-records-refused-service.md`
- `cas-object-model-dag-cbor-interop` at `scenarios/cas-object-model-dag-cbor-interop/cas-object-model-dag-cbor-interop.md`
- `l6-cas-starting-profile-bakeoff-long-horizon-reprofile` at `scenarios/l6-cas-starting-profile-bakeoff-long-horizon-reprofile/l6-cas-starting-profile-bakeoff-long-horizon-reprofile.md`

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

### `scenarios/cas-object-model-dag-cbor-interop/cas-object-model-dag-cbor-interop.md`

```markdown
# DAG-CBOR interop

## Scenario ID

cas-object-model-dag-cbor-interop

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jomag-cas-object-model/SCENARIOS.md`
- Source simulation: `SIM-jomag-cas-object-model/`
- Source row/title: DAG-CBOR interop
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jomag-cas-object-model/`.

## Setup

Alice stores a Merkle node or pointer object using a DAG-CBOR-compatible representation.

## Stimulus

Run the candidate simulation against this source test: Whether CID links, byte strings, and tags stay compatible with IPFS / IPLD-style tooling without requiring those stacks.

## Expected Pressure

TE-43 must decide whether DAG-CBOR is the default object format or only one allowed profile.
```

### `scenarios/l6-cas-starting-profile-bakeoff-long-horizon-reprofile/l6-cas-starting-profile-bakeoff-long-horizon-reprofile.md`

```markdown
# Long-horizon reprofile

## Scenario ID

l6-cas-starting-profile-bakeoff-long-horizon-reprofile

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-bobud-l6-cas-starting-profile-bakeoff/SCENARIOS.md`
- Source simulation: `SIM-bobud-l6-cas-starting-profile-bakeoff/`
- Source row/title: Long-horizon reprofile
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-bobud-l6-cas-starting-profile-bakeoff/`.

## Setup

A future implementation wants to replace the first profile with a richer object graph.

## Stimulus

Run the candidate simulation against this source test: Whether old pointer objects and raw chunks remain addressable and explainable after a later profile lands.

## Expected Pressure

The starting profile should avoid identity choices that become migration debt.
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

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `kernel-porting-boundary`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/kernel-porting-boundary/openai-gpt-5.4-medium/20260525-152502.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=3 evolution_safety=3 layer_boundary_clarity=4 failure_handling=3 implementation_plausibility=4 promise_vocabulary=4 simplicity_durability=4 envelope_discipline=2 kernel_implementation_promises=1 app_protocol_promise_semantics=1 risk_penalty=2
- Fitness: raw=37.00 normalized_0_100=56.92 confidence_0_1=0.85
- Rationale: Strong promise-first envelope specimen with clear sender-evidence semantics, exact-byte binding, and explicit unknown-pCID behavior. It scores lower for this scenario because kernel-porting-boundary asks for concrete kernel implementation promises, host assumptions, refusals, and port records that this envelope-level draft intentionally does not define. The universal mandatory outer proof slot also sits somewhat off the current smallest-envelope direction even though proof rules remain pCID-owned.
- Strengths:
  - Very clear envelope/payload boundary.
  - Promise vocabulary is mostly PT-correct and sender-scoped.
  - Exact-byte signing of [pCID, payload] aids evidence preservation and replay-safe audit.
  - ... 2 more
- Weaknesses:
  - Provides little of the kernel implementation promise record demanded by the scenario.
  - Does not enumerate app-facing storage/compute/send/receive/key/lifecycle promises or refusals.
  - Universal mandatory proof slot is less disciplined than the current minimal envelope direction.
  - ... 1 more
- Risks:
  - A mandatory universal third proof slot may overcommit the envelope layer before kernel/app boundaries freeze.
  - Generic peers may have weaker audit clarity because proof-family details are hidden behind pCID-specific logic.
  - Future proof-family changes may require minting new protocol pCIDs more often than desired.
- Open questions:
  - Does coupling payload shape and proof-family evolution to one pCID force unnecessary pCID churn?
  - Is a universal outer signature slot still too much envelope-level commitment for long-lived generic peers?
  - What minimum kernel-local promise record is needed to port this envelope shape without hiding the pCID-selected boundary behind host APIs?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `promise-accounting-records-kept-storage-promise`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/promise-accounting-records-kept-storage-promise/openai-gpt-5.4-medium/20260525-152502.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=3 evolution_safety=3 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 promise_vocabulary=4 simplicity_durability=4 envelope_discipline=2 kernel_implementation_promises=2 app_protocol_promise_semantics=2 risk_penalty=2
- Fitness: raw=39.00 normalized_0_100=60.00 confidence_0_1=0.84
- Rationale: Strong promise-first envelope specimen with clear local trust and boundary discipline, but this scenario asks for kept storage-promise accounting and the design intentionally leaves that to the payload protocol. As an envelope it is compact and plausible; as an answer to storage-promise evidence it is only partial.
- Strengths:
  - Very clear layer boundary between outer envelope evidence and higher-layer promise accounting.
  - Good Promise-Theory wording for sender-scoped intent and receiver-local trust.
  - Small deterministic artifact with explicit unknown-pCID behavior.
- Weaknesses:
  - Does not itself model storage promises, retrieval outcomes, or local promise-accounting records.
  - No kernel-level storage/dispatch/evidence promises are specified.
  - Mandatory outer signature slot is somewhat misaligned with the most minimal envelope direction.
- Risks:
  - Scenario pressure may be under-served unless higher-layer payload protocols define explicit kept-storage evidence records.
  - A mandatory universal outer proof slot may overconstrain envelopes that need lighter or differently scoped evidence.
  - Coupling proof-family semantics to the protocol pCID may make evolution and generic audit harder.
- Open questions:
  - How would a higher-layer payload protocol encode locally auditable evidence that Alice kept a storage promise and later served chunk C?
  - Does making the outer proof slot mandatory age well for payloads that want unsigned observation records or other lighter evidence forms?
  - Will proof-family evolution force too many new protocol pCIDs, reducing long-term audit clarity?
- Authority boundary: Envelope evidence only; storage-promise accounting, kept/failed outcome records, and trust updates belong to payload protocols and local peers.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `promise-accounting-records-refused-service`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/promise-accounting-records-refused-service/openai-gpt-5.4-medium/20260525-152502.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=3 evolution_safety=3 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 envelope_discipline=4 kernel_implementation_promises=1 app_protocol_promise_semantics=1 risk_penalty=2
- Fitness: raw=41.00 normalized_0_100=63.08 confidence_0_1=0.83
- Rationale: Strong PromiseGrid-style envelope specimen with clear promise-first semantics and explicit local trust boundaries. It is well aligned at the envelope layer, but this scenario asks for higher-layer refusal accounting, which the design intentionally leaves to the payload protocol. That delegation is legitimate, yet it limits direct scenario coverage.
- Strengths:
  - Very clear layer boundary: outer envelope authenticates protocol-bound payload bytes while leaving application promise accounting to the payload protocol.
  - Promise vocabulary is unusually strong and PT-correct for an envelope design.
  - Simple deterministic CBOR array with exact-byte signing supports durable evidence and small-device implementation.
  - ... 1 more
- Weaknesses:
  - Does not itself define how honest refusal is recorded differently from failure, corruption, or timeout.
  - Kernel implementation promises are largely absent because the specimen is not a kernel contract.
  - Mandatory outer proof slot adds some universal structure while still leaving proof semantics opaque to generic peers.
- Risks:
  - Refused-service evidence may be underspecified unless payload protocols consistently define signed refusal records distinct from timeout or corruption.
  - Generic auditors may have reduced clarity because proof-family semantics are hidden behind the payload protocol pCID rather than an explicit outer proof selector.
  - Proof-family changes may require minting new protocol pCIDs, creating migration and compatibility friction.
- Open questions:
  - Should a signed refusal record live entirely in the payload protocol, or should the envelope expose a more generic refusal/evidence affordance?
  - Does coupling proof-family evolution to the payload protocol pCID cause too much pCID churn over long-lived deployments?
  - How much audit clarity is lost for generic peers when the outer proof slot is mandatory but its verification family is defined only by the payload protocol?
- Authority boundary: Evidence only; this scores an envelope-layer specimen and does not settle higher-layer promise-accounting design.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `cas-object-model-dag-cbor-interop`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/cas-object-model-dag-cbor-interop/openai-gpt-5.4-medium/20260525-152502.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=4 envelope_discipline=3 kernel_implementation_promises=1 app_protocol_promise_semantics=1 risk_penalty=2
- Fitness: raw=44.00 normalized_0_100=67.69 confidence_0_1=0.84
- Rationale: Strong Promise-Theory-clean envelope specimen with clear layer boundaries, explicit unknown-pCID behavior, and good auditability for signed protocol-shaped bytes. It fits the DAG-CBOR interop scenario only partially because it can carry DAG-CBOR/IPLD-compatible payloads and CIDs but intentionally leaves object-model compatibility and profile commitments to the payload protocol rather than the envelope.
- Strengths:
  - Very clear envelope-layer scope and authority boundary.
  - Promise-first wording is unusually strong for an envelope design.
  - Small deterministic CBOR positional array is durable and implementation-friendly.
  - ... 2 more
- Weaknesses:
  - Does not itself specify CAS object-model behavior or DAG-CBOR compatibility commitments.
  - Universal signature slot is less aligned with fully protocol-owned post-slot-0 slot discipline.
  - Provides little kernel implementation promise detail.
  - ... 1 more
- Risks:
  - Proof semantics are coupled to the payload protocol pCID, which may force new pCIDs for proof-profile evolution.
  - A mandatory universal third slot may be more envelope overreach than the current direction prefers.
  - Generic peers may find proof-family interpretation less legible without a separate outer proof selector.
  - ... 1 more
- Open questions:
  - Is coupling proof-family evolution to the payload protocol pCID acceptable over long-lived interop timelines?
  - Should DAG-CBOR/IPLD compatibility be declared by payload protocols rather than inferred from the outer envelope?
  - Does a mandatory universal third slot reduce future envelope flexibility compared with fully protocol-owned post-slot-0 layouts?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design. Envelope-layer specimen only: payload protocol still owns CAS object semantics, DAG-CBOR profile choices, and higher-layer trust/accounting.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `l6-cas-starting-profile-bakeoff-long-horizon-reprofile`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/l6-cas-starting-profile-bakeoff-long-horizon-reprofile/openai-gpt-5.4-medium/20260525-152502.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=4 evolution_safety=3 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=5 promise_vocabulary=4 simplicity_durability=4 envelope_discipline=3 kernel_implementation_promises=1 app_protocol_promise_semantics=1 risk_penalty=2
- Fitness: raw=43.00 normalized_0_100=66.15 confidence_0_1=0.86
- Rationale: Strong envelope-layer specimen: small deterministic shape, explicit pCID-as-protocol binding, and clear local handling of unknown protocols. It uses largely correct promise-first framing for sender evidence at the envelope boundary. Scores are limited mainly because the scenario asks for long-horizon reprofiling of CAS/profile objects, which this design intentionally delegates to payload protocols, and because the mandatory universal outer signature slot plus proof-family coupling to pCID may create some evolution and audit tradeoffs.
- Strengths:
  - Very clear layer boundary between outer envelope evidence and higher-layer promise accounting.
  - Deterministic three-slot artifact is compact and durable.
  - Unknown-pCID behavior is local and conservative: preserve bytes without pretending to understand them.
  - ... 1 more
- Weaknesses:
  - Only partial fit for a long-horizon profile-migration scenario because object-graph continuity is left to payload protocols.
  - Uses a universal mandatory proof slot, which is somewhat less aligned with the current minimal-envelope direction.
  - Provides little kernel-level implementation promise detail.
  - ... 1 more
- Risks:
  - Proof-family changes may require minting new protocol pCIDs more often than a cleaner separation would.
  - Generic peers can preserve bytes but may have limited insight into proof semantics without protocol-specific handlers.
  - A mandatory universal outer proof slot may be more envelope than needed for content that only needs durable carriage.
- Open questions:
  - Does tying proof-family evolution to the protocol pCID create avoidable long-horizon migration churn?
  - When profiles are replaced, what payload-level records make old pointer objects and raw chunks still explainable to later auditors?
  - Is generic auditability weakened when the outer proof slot is explicit but its verification family is only discoverable through the protocol named by pCID?
- Authority boundary: Evidence only at the envelope layer; payload/profile migration, trust updates, and application semantics remain outside this specimen.

### `SIM-fovip-kernel-promise-boundary-port-contract` x `kernel-porting-boundary`

- Result path: `results/SIM-fovip-kernel-promise-boundary-port-contract/kernel-porting-boundary/openai-gpt-5.4-medium/20260525-152502.json`
- Scores: scenario_fit=5 promisegrid_alignment=5 auditability=5 evolution_safety=5 layer_boundary_clarity=5 failure_handling=5 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 envelope_discipline=4 kernel_implementation_promises=5 app_protocol_promise_semantics=3 risk_penalty=1
- Fitness: raw=59.00 normalized_0_100=90.77 confidence_0_1=0.90
- Rationale: Excellent fit for the scenario: it turns the porting boundary into an explicit kernel implementation promise record, separates host assumptions from local promises, keeps trust local, and makes app/kernel crossings pCID-selected messages. It is especially strong on boundary clarity, evidence, unsupported-feature disclosure, and promise-first framing. Minor deductions are for remaining draft-state ambiguity around exact record standardization and the relatively rich record surface for very small ports.
- Strengths:
  - Directly addresses the kernel-porting-boundary scenario and its runtime variants.
  - Very clear separation between kernel promises, host assumptions, and unsupported features.
  - Strong Promise-Theory framing: local promiser, local trust, no global certificate, no universal namespace authority.
  - ... 2 more
- Weaknesses:
  - Exact canonical encoding and interoperability expectations for the promise record remain provisional.
  - The record is broader than a tiny minimal artifact and may need profile-specific compression guidance.
  - Higher-layer app protocol semantics are intentionally delegated rather than fully modeled here.
- Risks:
  - Because the specimen is still a draft, different ports could publish similar-looking records with incompatible detail levels before a frozen pCID exists.
  - The record shape may be somewhat heavy for minimal MCU deployments unless a compact profile convention is nailed down.
  - Adapter and split-service evidence chains could become uneven across implementations if not further constrained.
- Open questions:
  - What exact on-wire or signed encoding, if any, should carry the kernel implementation promise record before freeze?
  - How should split-local-service ports normalize evidence handoff so each local promiser remains explicit without bloating records?
  - What minimum mandatory evidence retention is needed for tiny MCU profiles versus richer native ports?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-fovip-kernel-promise-boundary-port-contract` x `promise-accounting-records-kept-storage-promise`

- Result path: `results/SIM-fovip-kernel-promise-boundary-port-contract/promise-accounting-records-kept-storage-promise/openai-gpt-5.4-medium/20260525-152502.json`
- Scores: scenario_fit=3 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=4 promise_vocabulary=4 simplicity_durability=3 envelope_discipline=4 kernel_implementation_promises=5 app_protocol_promise_semantics=2 risk_penalty=1
- Fitness: raw=51.00 normalized_0_100=78.46 confidence_0_1=0.84
- Rationale: Strong PT-clean kernel-layer design. It gives explicit app-facing storage and evidence promises, plus local kept/refused/broken record handling, which helps this scenario. But the scenario asks for Bob's local accounting and future pull/keep/advertise choices around an inter-agent storage promise, and that higher-layer storage semantics are only partly covered here and largely delegated to payload/app protocols.
- Strengths:
  - Clear app/kernel promise boundary with pCID-selected messages.
  - Explicit evidence policy for kept, refused, unavailable, and broken promises.
  - Strong separation of host assumptions from kernel promises.
  - ... 1 more
- Weaknesses:
  - Does not fully specify the higher-layer storage promise semantics this scenario centers on.
  - Future pull/keep/advertise decision logic for Bob is left implicit.
  - Record shape is somewhat broad for very small durable artifacts.
- Risks:
  - A port may publish good local evidence policy while still lacking a concrete inter-agent storage promise protocol.
  - Record breadth could encourage checklist-style implementations unless the minimum credible subset stays tight.
- Open questions:
  - What minimal higher-layer storage pCID is needed to express retention period, later retrieval, and Bob's pull/advertise trust updates?
  - What exact kept-promise evidence record should a port retain when storage success is observed through an adapter or delegated host surface?
- Authority boundary: Evidence only; evaluates a kernel-layer promise record and does not settle higher-layer storage protocol design.

### `SIM-fovip-kernel-promise-boundary-port-contract` x `promise-accounting-records-refused-service`

- Result path: `results/SIM-fovip-kernel-promise-boundary-port-contract/promise-accounting-records-refused-service/openai-gpt-5.4-medium/20260525-152502.json`
- Scores: scenario_fit=4 promisegrid_alignment=5 auditability=5 evolution_safety=5 layer_boundary_clarity=5 failure_handling=5 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=4 envelope_discipline=5 kernel_implementation_promises=5 app_protocol_promise_semantics=3 risk_penalty=1
- Fitness: raw=60.00 normalized_0_100=92.31 confidence_0_1=0.88
- Rationale: Strong fit at the claimed kernel layer. The specimen explicitly distinguishes refused, unavailable, kept, and broken promises in its evidence policy and makes refusal an honest local outcome rather than generic failure. It stays promise-first, keeps trust local, and cleanly separates kernel promises from host assumptions. Scenario coverage is slightly limited only because application-specific promise accounting semantics are intentionally delegated above the kernel boundary.
- Strengths:
  - Explicit evidence categories include refused, unavailable, kept, and broken promises.
  - Clear kernel/app and kernel/host boundary language supports honest refusal without authority framing.
  - Strong promise-first vocabulary with local trust and no global certificate semantics.
  - ... 1 more
- Weaknesses:
  - The draft names refusal evidence but does not fully specify a concrete refusal-record wire shape in the provided text.
  - Scenario handling is kernel-centric and stops short of richer higher-layer promise-accounting workflows.
  - The record shape is somewhat broad, which may increase implementation variance across runtime classes.
- Risks:
  - Because the record is still a draft, different ports could encode refusal evidence incompatibly before a stable pCID-defined shape exists.
  - Consumers may misread host assumptions as guaranteed kernel behavior if implementations do not keep the separation sharp.
  - In split-service deployments, refusal attribution could become ambiguous if multiple local promisers emit overlapping evidence.
- Open questions:
  - What exact byte-level refusal record shape should ports emit so refusal stays distinguishable from unavailability, corruption, and broken promises across profiles?
  - How should split local-service ports correlate refusal evidence across dispatch, storage, network, and key promisers without implying a single authority?
  - What minimum app-facing mapping is needed so local apps can reliably consume refusal evidence without assuming higher-layer application semantics?
- Authority boundary: Evidence only; draft kernel-port promise record, not a final PromiseGrid kernel API or global trust authority.

### `SIM-fovip-kernel-promise-boundary-port-contract` x `cas-object-model-dag-cbor-interop`

- Result path: `results/SIM-fovip-kernel-promise-boundary-port-contract/cas-object-model-dag-cbor-interop/openai-gpt-5.4-medium/20260525-152502.json`
- Scores: scenario_fit=2 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 envelope_discipline=5 kernel_implementation_promises=5 app_protocol_promise_semantics=3 risk_penalty=2
- Fitness: raw=53.00 normalized_0_100=81.54 confidence_0_1=0.84
- Rationale: Strong PromiseGrid kernel-layer specimen. It is not a DAG-CBOR object-model design, so scenario fit is limited, but it cleanly supports this scenario by forcing explicit supported-pCID claims, exact-byte evidence policy, and unsupported behavior without pretending the kernel settles higher-layer CAS semantics. The main gap is that it does not decide whether DAG-CBOR is the default profile or define CID/tag/link semantics itself.
- Strengths:
  - Very explicit app/kernel promise boundary with pCID-selected grid messages.
  - Strong local evidence framing for kept, refused, unavailable, and broken promises.
  - Clear separation of implementation promises from host assumptions and unsupported features.
  - ... 2 more
- Weaknesses:
  - Does not define the DAG-CBOR object model or decide default-vs-optional profile status.
  - Leaves CID-link, byte-string, and tag compatibility semantics to higher-layer protocols.
  - Scenario fit is partial because the specimen is about kernel implementation promises, not CAS encoding rules.
- Risks:
  - Readers may over-credit the kernel promise record as solving CAS object-model interop when it mostly defines implementation disclosure and evidence boundaries.
  - Different ports could claim broad interoperability while differing on exact byte preservation, tag handling, or adapter normalization unless a higher-layer pCID tightens those rules.
  - If supported-pCID declarations stay too coarse, DAG-CBOR compatibility claims could become vague marketing rather than auditable local promises.
- Open questions:
  - How should a kernel promise record name DAG-CBOR support: exact-byte carriage only, parse-and-dispatch support, or a specific higher-layer object-model pCID?
  - What exact evidence must a port keep to show CID-link, byte-string, and tag compatibility across adapters without silently normalizing bytes?
  - Does DAG-CBOR belong as a default payload profile above the kernel boundary, or only as one optional pCID-owned object encoding?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-fovip-kernel-promise-boundary-port-contract` x `l6-cas-starting-profile-bakeoff-long-horizon-reprofile`

- Result path: `results/SIM-fovip-kernel-promise-boundary-port-contract/l6-cas-starting-profile-bakeoff-long-horizon-reprofile/openai-gpt-5.4-medium/20260525-152502.json`
- Scores: scenario_fit=2 promisegrid_alignment=5 auditability=5 evolution_safety=4 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=3 envelope_discipline=5 kernel_implementation_promises=5 app_protocol_promise_semantics=2 risk_penalty=1
- Fitness: raw=53.00 normalized_0_100=81.54 confidence_0_1=0.84
- Rationale: Strong Promise-Grid-aligned kernel-layer specimen: explicit local implementation promises, host assumptions, unsupported pCID behavior, exact-byte evidence, adapters, namespaces, references, and checkpoint behavior. It handles evolution pressure well at the kernel boundary, but the scenario's core question about old pointer objects and raw chunks after later profile replacement is mostly delegated to higher-layer CAS/profile protocols.
- Strengths:
  - Excellent kernel-layer promise framing with explicit local promiser identity.
  - Clear separation of implementation promises from host/runtime assumptions.
  - Strong evidence and refusal/broken-promise accounting supports audit and trust updates.
  - ... 2 more
- Weaknesses:
  - Scenario fit is limited because old CAS object semantics are not defined here.
  - Does not itself specify how obsolete pointer objects remain explainable after reprofiling.
  - Record shape is larger than the smallest durable artifact and may need tight profiling for MCU-class ports.
- Risks:
  - Kernel promise record may be mistaken for solving object-profile migration when it mainly provides substrate/evidence hooks.
  - Record shape is somewhat broad, which could invite profile sprawl across constrained ports.
  - Long-horizon interpretability still depends on durable higher-layer pCID specs and local evidence retention policies.
- Open questions:
  - Which kernel promise, if any, commits to preserving exact old profile bytes versus only best-effort carriage/checkpoint retention across reprofiling?
  - How are old pointer-object and raw-chunk explanations surfaced to apps when a later richer profile supersedes the first one?
  - What minimum reference/checkpoint evidence is durable enough for tiny ports while still supporting long-horizon explainability?
- Authority boundary: Scores the kernel-port promise record at the kernel boundary only; CAS object/profile migration semantics remain higher-layer pCID-owned protocol work and local trust judgments.

## Required JSON Shape

{"child_id":"SIM-lasuv-child-descriptive-design-slug","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
