# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Choose a descriptive `child_id` that starts with `SIM-jopuz-child-` and ends with a kebab-case design slug. Do not use generic `ga-child`, `pending`, or ordinal-only names.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260521-234037`
- Planned child ID prefix: `SIM-jopuz-child-`
- Temporary child ID: `SIM-jopuz-child-pending-0002`
- Temporary child path: `proposals/ga-canary-20260521-234037/simulations/SIM-jopuz-child-pending-0002/`
- Operation: `breed`
- Parent IDs: `SIM-zisan-device-bound-agent, SIM-funas-kernel-porting-boundary`

## Scenario Sample

- `minimal-immutable-blob-app` at `scenarios/minimal-immutable-blob-app/minimal-immutable-blob-app.md`
- `app-semantics-partial-conformance` at `scenarios/app-semantics-partial-conformance/app-semantics-partial-conformance.md`
- `device-bound-agent-physical-effect` at `scenarios/device-bound-agent-physical-effect/device-bound-agent-physical-effect.md`
- `live-crdt-audit-publication` at `scenarios/live-crdt-audit-publication/live-crdt-audit-publication.md`
- `multi-embodiment-app-identity` at `scenarios/multi-embodiment-app-identity/multi-embodiment-app-identity.md`
- `portable-signing-key-identity` at `scenarios/portable-signing-key-identity/portable-signing-key-identity.md`
- `kernel-porting-boundary` at `scenarios/kernel-porting-boundary/kernel-porting-boundary.md`

## Scenario Pressure

### `scenarios/minimal-immutable-blob-app/minimal-immutable-blob-app.md`

```markdown
# Minimal Immutable Blob App

## Scenario ID

minimal-immutable-blob-app

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-vopik - What CAS-facing guarantees are safe for a minimal immutable blob app?`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-tuhaz`; `DR-tumus`

## Purpose

Exercise candidate designs against the smallest useful CAS-facing app: Alice
writes immutable bytes and receives a hash; later Carol presents the hash and
expects to retrieve the same bytes.

## Setup

Alice uploads a blob through Bob's app. Bob stores or publishes a
content-addressed object and returns a hash. Carol later receives only the hash
and partial context. Mallory may withhold storage, replay stale availability
claims, or claim that possession of the hash is enough authorization.

## Stimulus

The original host changes retention policy, a peer cache is incomplete, and Carol
tries to read the blob from a different site years later.

## Expected Pressure

The candidate design must separate content identity from availability,
authorization, ingress, discovery, replication, and retention promises while
still preserving enough evidence for a 100-year audit trail.

## Scenario-Specific Evaluation Questions

- What exactly does `hash in -> blob out` promise, and who makes that promise?
- Is possession of the hash a read capability, an address, or an app-level
  convention?
- What local promise accounting records should Alice, Bob, and Carol keep when
  storage or retrieval fails?
```

### `scenarios/app-semantics-partial-conformance/app-semantics-partial-conformance.md`

```markdown
# App Semantics Partial Conformance

## Scenario ID

app-semantics-partial-conformance

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-dodos`, `FB-hisis`, `FB-kutub`, `FB-gomod`, and `FB-tahof`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-tuhaz`

## Purpose

Test whether an app can publish honest provisional semantics and conformance
claims without pretending to implement a final PromiseGrid app contract.

## Setup

Alice ships a useful first app slice. It uses local IDs internally, signs current
messages with an adapter-local carriage choice, and supports append/read behavior
but not the full merge, authority, capability-token, or break-witness story.
Carol needs to know what can interoperate. Mallory benefits if the app overclaims.

## Stimulus

Alice publishes a B-side conformance claim and Bob tries to interoperate using
only that claim, the draft spec path, and locally observed wire artifacts.

## Expected Pressure

The candidate design must distinguish local implementation shortcuts from
protocol-boundary identity and conformance, and must say which semantic claims
are provisional, blocked, or safe orientation.

## Scenario-Specific Evaluation Questions

- What wording makes a partial-conformance claim honest?
- Which identity is authoritative at the protocol boundary?
- How should provisional signature carriage and capability/witness language be
  described without freezing them?
```

### `scenarios/device-bound-agent-physical-effect/device-bound-agent-physical-effect.md`

```markdown
# Device-Bound Agent Physical Effect

## Scenario ID

device-bound-agent-physical-effect

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-nojit`, `FB-tisuf`, and `FB-tulit`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-tuhaz`; `DR-davod`

## Purpose

Exercise apps where Bob's local daemon controls a host-owned physical device and
Alice asks it to perform or report something through a grid protocol.

## Setup

Bob owns a label printer and a temperature sensor. He delegates the front-desk
machine to print shipping labels and the rack sensor to report readings. Alice
sends a print request. Carol audits the receipt. Mallory replays the request
after the label has already been printed.

## Stimulus

The daemon restarts after receiving the request but before all receipts propagate.
It must decide whether replay should print another label, report an already-done
effect, or emit a break-witness.

## Expected Pressure

The candidate design must handle owner/operator identity, host-driver
dependencies, non-idempotent physical effects, at-most-once posture, receipts,
break-witnesses, and 100-year interpretation of evidence after the device and
driver stack are gone.

## Scenario-Specific Evaluation Questions

- What evidence proves the owner delegated the device-bound agent?
- How is a physical effect deduplicated across replay and restart?
- What should a conformance claim say about CUPS, libusb, IIO, i2c, IPP, or
  vendor SDK dependencies?
```

### `scenarios/live-crdt-audit-publication/live-crdt-audit-publication.md`

```markdown
# Live CRDT Audit Publication

## Scenario ID

live-crdt-audit-publication

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-hurit` and `FB-nilat`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-tuhaz`; `group-session-draft.md`; `udp-binding-draft.md`

## Purpose

Test real-time app pressure where live state needs reliable, ordered, low-latency
frames, but durable PromiseGrid evidence may be published at milestones.

## Setup

Alice edits a shared document in a browser while Bob edits the same document in
Neovim. Their live CRDT sync needs sub-second in-order delivery. Carol audits
durable snapshots later through group-session messages that cite content-addressed
state. Mallory drops or reorders live frames and delays audit publication.

## Stimulus

The live channel partitions for thirty seconds, then reconnects. Alice and Bob
continue editing. The app emits an audit message at save time with a snapshot
reference and human-readable promise body.

## Expected Pressure

The candidate design must avoid pretending that best-effort datagrams or
git-paced group-session are the live transport, while showing how durable audit
evidence can still survive for 100-year review.

## Scenario-Specific Evaluation Questions

- Should live state be off-grid until a reliable binding exists, or should a
  future live pCID shape be sketched?
- What exact object does the audit message cite?
- How are live-channel conformance claims kept separate from audit-layer claims?
```

### `scenarios/multi-embodiment-app-identity/multi-embodiment-app-identity.md`

```markdown
# Multi-Embodiment App Identity

## Scenario ID

multi-embodiment-app-identity

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-zazon - heterogeneous embodiments but one app identity`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-tuhaz`; `DR-davod`

## Purpose

Exercise one logical app implemented by heterogeneous components that claim the
same pCID-selected contract without pretending each component implements every
part.

## Setup

Alice uses a browser tab with no filesystem access. Bob uses a Neovim plugin
with a long-running Node helper process. Both participate in the same app and
document. Carol reads their conformance claims. Mallory exploits ambiguity
between "same UX app" and "same protocol contract."

## Stimulus

The browser and plugin exchange data through a custom live channel and publish
separate conformance claims naming which parts of the shared contract they
implement.

## Expected Pressure

The candidate design must show how one app stays one app by honoring an explicit
protocol contract or family of contracts, while each embodiment states its own
runtime constraints and implementation scope.

## Scenario-Specific Evaluation Questions

- What makes the browser and plugin one app rather than two unrelated apps?
- What must each component's conformance claim include?
- How does the answer survive host replacement, browser storage loss, helper
  upgrades, and long-term audit?
```

### `scenarios/portable-signing-key-identity/portable-signing-key-identity.md`

```markdown
# Portable Signing-Key Identity

## Scenario ID

portable-signing-key-identity

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-robif - baseline signing-key identity across browser and plugin host`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-tuhaz`

## Purpose

Test a provisional identity recipe for one user appearing consistently across a
browser tab and a long-running plugin host.

## Setup

Alice has a browser profile and a plugin helper. The browser can use WebCrypto
and IndexedDB. The helper can use Node crypto and a filesystem. Bob observes only
signed protocol artifacts. Carol audits a later key rotation. Mallory injects a
display-name collision and attempts to confuse local usernames with durable
identity.

## Stimulus

Alice rotates from an old signing key to a new signing key and opens a new live
session from the other host embodiment.

## Expected Pressure

The candidate design must distinguish signing-key continuity from presentation
hints, define a pivotable v0 recipe without hardcoding forever cryptography, and
show how constrained-host storage changes the claim.

## Scenario-Specific Evaluation Questions

- Which parts of the key algorithm, rotation, and handshake can be provisional?
- What evidence links old and new keys?
- What should the guide say about browser key storage and XSS risk without
  overstating security?
```

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
the kernel/runtime terminology remains unsettled.

## Setup

Alice wants to port PromiseGrid infrastructure to a new host environment. Bob
offers a library, Carol offers a dispatcher/runtime, and Mallory claims that
copying the wire-lab harness is the porting target. The available specs are
drafts, and only some future frozen pCIDs will become obligations.

## Stimulus

Alice writes a first porting plan and a conformance claim. The plan must say what
it implements now, what draft evidence it follows, and what it defers until
`DR-davod` closes.

## Expected Pressure

The candidate design must separate harness apparatus from porting target,
identify which binding/session/message/CAS/runtime obligations are provisional,
and preserve a clear path to future frozen-spec conformance.

## Scenario-Specific Evaluation Questions

- Should the guide say kernel, runtime, dispatcher, handler host, or library?
- What is the minimum viable porting target before final freeze?
- Which K1-K5 ingress, feed, CAS, session, and app-layer details should remain
  blocked versus provisional orientation?
```

## Parent Simulation Documents

### `simulations/SIM-zisan-device-bound-agent/README.md`

```markdown
# SIM-zisan: Device-bound agent

This simulation is a provisional question home for `FB-nojit`, `FB-tisuf`, and
`FB-tulit`: apps that expose a host-owned physical device, sensor, actuator, or
driver stack to grid peers. Source: `DI-ragaz`.

## Question

What can the guide safely teach about a device-bound PromiseGrid app before the
final app and kernel boundaries settle? Source: `DI-ragaz`.

## Decision Axes

- **Agent role:** a long-running userspace daemon acts as the protocol-facing
  voice of a host-owned device.
- **Owner and operator authority:** owner keys, delegated operator keys, and
  physical custody need explicit evidence.
- **Irreversible effects:** actuator requests need at-most-once posture,
  deduplication, receipts, and break-witnesses for non-idempotent outcomes.
- **Sensor versus actuator:** sensors and actuators share device-bound identity
  and driver-stack concerns; actuators add irreversible-effect pressure.
- **B-side conformance:** code may claim a grid protocol while honestly naming
  host dependencies such as CUPS, libusb, Linux IIO, i2c, IPP, or vendor SDKs.

## Related Root Scenario

- `scenarios/device-bound-agent-physical-effect/device-bound-agent-physical-effect.md`

## Boundaries

This simulation does not freeze a device-agent protocol or driver API. It tests
guide-safe patterns for provisional tutorials and conformance claims. Source:
`DI-ragaz`.
```

### `simulations/SIM-zisan-device-bound-agent/QUESTION.md`

```markdown
# Question

Can PromiseGrid guide prose describe device-bound agents, owner/operator
delegation, irreversible physical effects, and host-driver conformance as
provisional App Dev patterns while `DR-tuhaz` and `DR-davod` remain open? Source:
`DI-ragaz`.

Open decision points:

- What request, receipt, deduplication, and break-witness shape is needed for
  non-idempotent physical actions?
- How does an owner bind a signing key to a device and delegate operation?
- What should a conformance claim say when the implementation depends on a
  host-specific driver stack?
```

### `simulations/SIM-funas-kernel-porting-boundary/README.md`

```markdown
# SIM-funas: Kernel porting boundary

This simulation is a provisional question home for `FB-vitih`, `FB-mulum`, and
`FB-potin`: what "kernel developer" or "porter" should mean while `DR-davod`
remains open. Source: `DI-ragaz`.

## Question

What can the guide safely teach about the PromiseGrid porting target before the
stable kernel/runtime boundary is decided? Source: `DI-ragaz`.

## Decision Axes

- **Term of art:** kernel, runtime, dispatcher, handler host, or library surface.
- **Minimum viable port:** which frozen binding/session/message specs and
  conformance claims a port must implement first.
- **Runtime obligations:** handler dispatch, pCID routing, storage, key handling,
  ingress, feeds, CAS subtree, and implementation changelog evidence.
- **Layer boundaries:** what belongs to substrate/feed/group/CAS/app layers and
  what must not be taught as a monolithic harness clone.
- **Provisional versus blocked:** which K1-K5 ingress and runtime details are
  teachable orientation versus blocked settled instructions.

## Related Root Scenario

- `scenarios/kernel-porting-boundary/kernel-porting-boundary.md`

## Boundaries

This simulation does not define the final PromiseGrid porting API. It keeps the
kernel/runtime/dispatcher framing testable until `DR-davod` closes. Source:
`DI-ragaz`.
```

### `simulations/SIM-funas-kernel-porting-boundary/QUESTION.md`

```markdown
# Question

Which developer-facing porting boundary can PromiseGrid describe now without
turning the wire-lab harness, every simulation-local draft, or an unsettled
runtime surface into a false implementation target? Source: `DI-ragaz`.

Open decision points:

- Should the guide keep "kernel", shift to "runtime", or use a narrower
  dispatcher/handler-hosting term?
- What minimum set of frozen specs and conformance claims defines a first real
  port?
- Which K1-K5 ingress, CAS subtree, feed, session, and app-layer details are
  provisional orientation rather than stable porting obligations?
```

## Compact Fitness Evidence From This Run

### `SIM-zisan-device-bound-agent` x `minimal-immutable-blob-app`

- Result path: `results/SIM-zisan-device-bound-agent/minimal-immutable-blob-app/openai-gpt-5.4-xhigh/20260521-234037.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=3 evolution_safety=3 layer_boundary_clarity=2 failure_handling=2 implementation_plausibility=2 risk_penalty=3
- Fitness: raw=12.00 normalized_0_100=34.00 confidence_0_1=0.88
- Rationale: SIM-zisan is mainly a device-bound agent pattern, not a minimal CAS blob app. Its explicit evidence and boundary discipline are somewhat transferable, but it leaves the scenario's core blob promises largely unspecified.
- Strengths:
  - Emphasizes explicit local evidence such as receipts, witnesses, and delegation records.
  - Keeps protocol-facing behavior separate from host-specific dependencies and conformance claims.
  - Stays evolution-aware by treating boundaries as provisional rather than prematurely frozen.
- Weaknesses:
  - Does not define the core 'hash in -> blob out' promise or who is accountable for it.
  - Does not separate content identity from availability, authorization, replication, discovery, or retention promises.
  - Does not model years-later retrieval from partial caches or stale availability claims.
- Risks:
  - Applying device-custody patterns here could blur blob semantics with operator or host authority.
  - A weak tutorial derived from this sim might incorrectly treat hash possession as sufficient authorization.
  - Open device-agent boundary work could distract from the simpler blob-specific guarantees this scenario needs.
- Open questions:
  - Can the receipt and witness patterns be generalized into minimal CAS publish/read accounting?
  - What local records should Alice, Bob, and Carol keep when retention changes or retrieval fails years later?
  - Should minimal blob guidance explicitly forbid treating a hash as a read capability unless separately promised?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-zisan-device-bound-agent` x `app-semantics-partial-conformance`

- Result path: `results/SIM-zisan-device-bound-agent/app-semantics-partial-conformance/openai-gpt-5.4-xhigh/20260521-234037.json`
- Scores: scenario_fit=5 promisegrid_alignment=4 auditability=4 evolution_safety=5 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=28.00 normalized_0_100=80.00 confidence_0_1=0.83
- Rationale: Strong fit: SIM-zisan explicitly targets provisional app-dev guidance, honest B-side conformance, and protocol-versus-host boundary clarity for device-bound agents. It is less complete on generic app-semantic claim wording and still leaves key authority, signature-carriage, and witness details open.
- Strengths:
  - Explicitly addresses honest partial-conformance claims without pretending final app or kernel boundaries are settled.
  - Clearly separates protocol-facing agent identity from host-specific driver-stack dependencies.
  - Emphasizes auditable delegation and irreversible-effect handling through explicit evidence, receipts, deduplication, and break-witness concepts.
  - ... 1 more
- Weaknesses:
  - Does not yet give a concrete template for the published partial-conformance claim.
  - Leaves authoritative boundary identity and signature-carriage wording unresolved.
  - Covers device-bound concerns more strongly than generic append/read/merge app semantics.
- Risks:
  - Bob could be misled if adapter-local shortcuts are presented as stable interoperable semantics.
  - Ambiguity between device, owner, operator, and agent authority could blur who is making the protocol promise.
  - Host-specific dependencies and actuator side effects raise safety risk when conformance language is underspecified.
- Open questions:
  - What minimum fields and caveats must an honest partial-conformance claim include?
  - Which identity should be authoritative at the protocol boundary: device, owner, operator, or agent?
  - How should provisional signature carriage and capability or break-witness language be described without freezing them?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-zisan-device-bound-agent` x `device-bound-agent-physical-effect`

- Result path: `results/SIM-zisan-device-bound-agent/device-bound-agent-physical-effect/openai-gpt-5.4-xhigh/20260521-234037.json`
- Scores: scenario_fit=5 promisegrid_alignment=4 auditability=4 evolution_safety=5 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=3 risk_penalty=3
- Fitness: raw=32.00 normalized_0_100=80.00 confidence_0_1=0.80
- Rationale: Purpose-built for this scenario and especially strong on delegation, irreversible-effect, and host-driver boundary framing, but it is still more a well-aimed question home than a settled design: durable receipt/dedup/break-witness formats and 100-year authority/audit mechanics remain open.
- Strengths:
  - Near-exact match to the scenario's owner/operator, replay, and physical-effect pressure.
  - Explicitly names at-most-once posture, deduplication, receipts, and break-witnesses for non-idempotent actions.
  - Keeps PromiseGrid-facing agent behavior and host-specific driver dependencies clearly separated and provisional.
- Weaknesses:
  - No concrete crash-safe durable schema or local record flow is defined for requests, receipts, deduplication, and break-witnesses.
  - Sparse-knowledge and peer-local accounting behavior are not yet worked through in detail.
  - Conformance language for host-specific stacks is still posed as an open question.
- Risks:
  - Restart or replay mistakes could duplicate an irreversible physical action.
  - Weak device-to-key binding or delegation evidence could undermine later authority audits.
  - Guide prose could overclaim PromiseGrid conformance while hiding host-specific dependencies.
- Open questions:
  - What durable record shape lets Carol distinguish done, replayed, and broken effects after restart and still interpret them after the device stack is gone?
  - How should owner binding, operator delegation, and physical custody be evidenced?
  - How should conformance be stated when CUPS, libusb, IIO, i2c, IPP, or vendor SDKs are required?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-zisan-device-bound-agent` x `live-crdt-audit-publication`

- Result path: `results/SIM-zisan-device-bound-agent/live-crdt-audit-publication/openai-gpt-5.4-xhigh/20260521-234037.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=2 evolution_safety=4 layer_boundary_clarity=3 failure_handling=2 implementation_plausibility=1 risk_penalty=4
- Fitness: raw=16.00 normalized_0_100=40.00 confidence_0_1=0.85
- Rationale: Useful mainly as boundary-setting evidence: the simulation is careful about provisional claims, delegation, and auditable receipts, but it does not actually cover live CRDT transport or the exact save-time audit object this scenario needs.
- Strengths:
  - Explicitly keeps unsettled protocol and host/kernel boundaries provisional.
  - Names conformance boundaries honestly between PromiseGrid-facing behavior and host-specific dependencies.
  - Uses receipts, deduplication, and break-witness ideas that transfer to durable audit publication.
- Weaknesses:
  - Does not model sub-second ordered live sync, reconnect, or partition recovery for CRDT peers.
  - Does not define what immutable object a save-time audit message should cite.
  - Focuses on device custody and physical effects rather than collaborative document state.
- Risks:
  - Could be misread as covering real-time collaboration when it only offers general boundary discipline.
  - May blur live-channel claims with durable audit-layer claims if reused without a separate binding.
  - Could encourage overclaiming unsupported transport properties from host-local or best-effort mechanisms.
- Open questions:
  - Should this simulation explicitly mark live CRDT transport out of scope until a reliable binding exists?
  - What CAS snapshot or state object should a save-time audit message cite?
  - Is separate conformance language for live transport and durable audit publication sufficient, or is a different simulation needed?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-zisan-device-bound-agent` x `multi-embodiment-app-identity`

- Result path: `results/SIM-zisan-device-bound-agent/multi-embodiment-app-identity/openai-gpt-5.4-xhigh/20260521-234037.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=3 evolution_safety=3 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=3 risk_penalty=3
- Fitness: raw=23.00 normalized_0_100=58.00 confidence_0_1=0.79
- Rationale: Partial fit. The simulation is strong on honest embodiment-scoped conformance claims, host/runtime boundary disclosure, and local evidence around delegated operation, but it is centered on device-bound agents rather than defining a durable shared app identity across heterogeneous embodiments. It helps more with what each component may claim than with why a browser tab and plugin remain one app over time.
- Strengths:
  - Supports scoped conformance claims instead of pretending every embodiment implements the full stack.
  - Makes agent, host, and dependency boundaries explicit, which helps audit runtime constraints.
  - Receipts, deduplication, and break-witness ideas provide useful local evidence patterns for disputed actions.
- Weaknesses:
  - Does not define the shared contract or contract family that makes multiple embodiments one app.
  - Focuses on device ownership and operator delegation more than browser/plugin continuity.
  - Does not directly address browser storage loss, host replacement, or helper upgrade lineage.
- Risks:
  - UX or branding sameness could be mistaken for protocol sameness.
  - App identity could be conflated with device, host, or operator-key identity.
  - Loose partial-conformance wording could enable overclaiming by a malicious or sloppy embodiment.
- Open questions:
  - What minimum fields must each embodiment claim to bind itself to the same app contract while declaring omitted capabilities?
  - How do embodiment records and state survive browser eviction, host replacement, and helper rotation without central authority?
  - What audit trail shows that heterogeneous embodiments remained compatible as contracts and dependencies evolve?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-zisan-device-bound-agent` x `portable-signing-key-identity`

- Result path: `results/SIM-zisan-device-bound-agent/portable-signing-key-identity/openai-gpt-5.4-xhigh/20260521-234037.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=3 evolution_safety=4 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=3 risk_penalty=3
- Fitness: raw=24.00 normalized_0_100=60.00 confidence_0_1=0.80
- Rationale: Partial fit. SIM-zisan overlaps on host-bound key custody, owner/operator delegation, and honest constrained-host conformance claims for a long-running helper, but it does not yet define user-level continuity across browser and helper, explicit old-to-new key linkage, or browser storage/XSS guidance.
- Strengths:
  - Explicit focus on owner/operator authority and device-key binding.
  - Clear boundary between protocol-facing agent and host-specific driver stack.
  - Provisional, non-frozen posture supports future algorithm and handshake evolution.
- Weaknesses:
  - Not centered on same-user identity continuity across browser and plugin/helper embodiments.
  - No explicit rotation artifact or continuity proof for Carol to audit.
  - No concrete treatment of display-name collisions versus durable identity.
  - ... 1 more
- Risks:
  - Could conflate device/operator authority with durable personal identity.
  - Could overstate browser-held key security if XSS and profile compromise are not made explicit.
  - Could leave Bob with signed artifacts but no standard proof linking rotated keys.
- Open questions:
  - What signed object should link old and new keys across browser and helper sessions?
  - How should durable identity claims relate to device-bound or operator-delegated keys?
  - What exact guide language should compare WebCrypto/IndexedDB storage with helper filesystem storage under XSS and local compromise?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-zisan-device-bound-agent` x `kernel-porting-boundary`

- Result path: `results/SIM-zisan-device-bound-agent/kernel-porting-boundary/openai-gpt-5.4-xhigh/20260521-234037.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=3 evolution_safety=4 layer_boundary_clarity=3 failure_handling=2 implementation_plausibility=3 risk_penalty=3
- Fitness: raw=22.00 normalized_0_100=55.00 confidence_0_1=0.73
- Rationale: Moderate indirect fit. SIM-zisan is mainly about device-bound app agents, not the core PromiseGrid port, but it does help on provisional conformance, honest host-dependency claims, and the userspace-agent versus driver-stack boundary while kernel edges remain unsettled. It is weak on the scenario's core asks around minimum port target, runtime terminology, and K1-K5/CAS/session scope.
- Strengths:
  - Explicit provisional-conformance posture preserves room for later frozen specs.
  - Clear userspace-agent versus host-driver-stack boundary.
  - Realistic handling of host-specific dependencies in conformance claims.
- Weaknesses:
  - Primary focus is app/device integration, not kernel/runtime porting.
  - Does not clearly separate wire-lab harness apparatus from the true port target.
  - Little guidance on K1-K5 ingress, feed, CAS, and session obligations.
- Risks:
  - Alice could mistake a device-bound app integration for the core PromiseGrid port.
  - Mallory's claim that the harness is the port target is only weakly rebutted here.
  - Premature conformance language could harden before DR-davod settles boundaries.
- Open questions:
  - What is the minimum viable pre-freeze porting target below the device-bound agent layer?
  - Which obligations remain blocked versus provisional across ingress, feed, CAS, and session layers?
  - Which guide term should anchor this boundary: runtime, dispatcher, library, or kernel?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-funas-kernel-porting-boundary` x `minimal-immutable-blob-app`

- Result path: `results/SIM-funas-kernel-porting-boundary/minimal-immutable-blob-app/openai-gpt-5.4-xhigh/20260521-234037.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=3 evolution_safety=4 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=2 risk_penalty=2
- Fitness: raw=17.00 normalized_0_100=49.00 confidence_0_1=0.80
- Rationale: This simulation is useful mainly as a boundary-checking lens, not as a direct answer to the minimal immutable blob app. It helps separate CAS/app/runtime obligations and avoids teaching an unstable harness clone as the port target, but it does not specify the actual hash-to-blob promise, long-term availability behavior, or local failure accounting the scenario asks for.
- Strengths:
  - Explicitly distinguishes provisional guidance from blocked settled instructions.
  - Directly foregrounds CAS/app/runtime layer boundaries, which is central to this scenario.
  - Calls for frozen specs, conformance claims, and implementation changelog evidence.
- Weaknesses:
  - Does not define who promises `hash in -> blob out` or what that promise includes.
  - Does not explain discovery, retention, replication, or cross-site retrieval years later.
  - Provides little peer-local promise accounting for Alice, Bob, and Carol when retrieval fails.
- Risks:
  - Temporary porting guidance could ossify into an incorrect blob-contract expectation.
  - Implementers may conflate content identity with storage availability or read authorization.
- Open questions:
  - What minimum frozen spec set is enough for a first real immutable-blob app port?
  - Is possession of a hash only an address, or ever a capability/app convention?
  - What durable local records should survive incomplete-cache or stale-availability failures years later?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-funas-kernel-porting-boundary` x `app-semantics-partial-conformance`

- Result path: `results/SIM-funas-kernel-porting-boundary/app-semantics-partial-conformance/openai-gpt-5.4-xhigh/20260521-234037.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=3 evolution_safety=4 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=2 risk_penalty=3
- Fitness: raw=17.00 normalized_0_100=49.00 confidence_0_1=0.78
- Rationale: Best used as an upstream guardrail, not as a direct app-conformance scheme: it clearly separates stable versus provisional porting obligations and avoids freezing an unstable boundary, but it does not yet say how an app publishes an honest B-side partial-conformance claim that Bob and Carol can verify from protocol-boundary identity and local wire evidence alone.
- Strengths:
  - Explicitly distinguishes provisional guidance from blocked or stable obligations.
  - Strong attention to layer boundaries and to avoiding a false monolithic implementation target.
  - Connects early conformance claims to frozen specs rather than to sprawling draft behavior.
- Weaknesses:
  - Centers kernel/runtime porter scope more than app-semantic partial conformance.
  - Leaves protocol-boundary identity for app claims underspecified.
  - Does not provide a concrete locally auditable claim/evidence pattern for Bob and Carol.
- Risks:
  - Readers may over-apply porting-boundary language as if it were an app contract.
  - Unsettled kernel/runtime terminology could accidentally freeze adapter-local carriage choices.
  - Ambiguous partial-claim wording leaves room for Mallory to exploit interoperability assumptions.
- Open questions:
  - What exact claim vocabulary separates append/read support from blocked merge, authority, capability-token, and break-witness behavior?
  - Which identity and signature carriage are authoritative at the protocol boundary versus purely local implementation detail?
  - What minimal locally observable artifacts should Bob and Carol inspect to verify a partial-conformance claim?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-funas-kernel-porting-boundary` x `device-bound-agent-physical-effect`

- Result path: `results/SIM-funas-kernel-porting-boundary/device-bound-agent-physical-effect/openai-gpt-5.4-xhigh/20260521-234037.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=2 evolution_safety=4 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=2 risk_penalty=2
- Fitness: raw=17.00 normalized_0_100=49.00 confidence_0_1=0.83
- Rationale: Useful mainly as boundary-setting evidence: it helps separate PromiseGrid port/runtime obligations from device-, driver-, and app-specific behavior, but it does not itself resolve delegation, replay-safe physical effects, receipts, or break-witness semantics.
- Strengths:
  - Directly pressures where CUPS/libusb/IIO-style dependencies belong relative to a PromiseGrid port.
  - Explicitly avoids teaching the harness as the false implementation target.
  - Strong on provisional-vs-blocked guidance and cross-layer separation.
- Weaknesses:
  - No concrete evidence model for owner delegation of a device-bound agent.
  - No explicit at-most-once or replay-after-restart handling for non-idempotent effects.
  - Receipts and 100-year audit interpretation of physical effects are underspecified.
- Risks:
  - If overextended, readers could put physical-effect deduplication in the wrong layer.
  - Premature guidance could accidentally freeze host-driver assumptions into the port target.
- Open questions:
  - Is effect deduplication/break-witness behavior app logic, runtime contract, or a separate conformance profile?
  - What durable receipt format should survive device and driver obsolescence?
  - How should conformance claims describe optional CUPS/libusb/IIO/vendor-SDK dependencies?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-funas-kernel-porting-boundary` x `live-crdt-audit-publication`

- Result path: `results/SIM-funas-kernel-porting-boundary/live-crdt-audit-publication/openai-gpt-5.4-xhigh/20260521-234037.json`
- Scores: scenario_fit=3 promisegrid_alignment=4 auditability=2 evolution_safety=5 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=3 risk_penalty=2
- Fitness: raw=22.00 normalized_0_100=63.00 confidence_0_1=0.76
- Rationale: Strong as boundary-setting evidence: it helps keep unsettled live CRDT transport out of the current teachable port target and separate from durable audit publication. Weak as a full scenario answer because it does not define the reliable live channel, the cited audit object, or reconnect behavior.
- Strengths:
  - Explicit provisional-versus-blocked framing reduces false implementation-target risk while the kernel/runtime boundary is unsettled.
  - Layer-boundary focus directly helps separate live-channel concerns from feed, group-session, CAS, and app-layer audit claims.
  - Minimum-port emphasis on frozen specs and conformance claims gives a plausible stable core for durable publication work.
- Weaknesses:
  - It is not a concrete design for sub-second ordered CRDT transport.
  - It does not specify the exact snapshot or CAS object an audit message should cite.
  - Partition, reordering, and delayed-publication handling are not worked through in scenario detail.
- Risks:
  - Porters may still treat provisional guidance or harness behavior as the stable port surface.
  - If the live-versus-audit boundary stays underspecified, teams may blur live-channel guarantees with durable audit guarantees.
- Open questions:
  - Should live CRDT state stay off-grid until a reliable binding is frozen, or should a provisional live pCID shape be sketched?
  - Should the audit publication cite a whole-document snapshot, a CRDT state vector, or a CAS subtree/manifest?
  - What conformance split cleanly separates live-channel reliability claims from durable group-session audit claims?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-funas-kernel-porting-boundary` x `multi-embodiment-app-identity`

- Result path: `results/SIM-funas-kernel-porting-boundary/multi-embodiment-app-identity/openai-gpt-5.4-xhigh/20260521-234037.json`
- Scores: scenario_fit=3 promisegrid_alignment=4 auditability=3 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=3 risk_penalty=3
- Fitness: raw=26.00 normalized_0_100=65.00 confidence_0_1=0.76
- Rationale: Good partial fit: it strongly supports one app across multiple embodiments through explicit conformance claims and layer separation, but it does not yet fully define the contract-family, failure evidence, and durable audit rules that make that identity unambiguous.
- Strengths:
  - Separates substrate/feed/group/CAS/app duties, so browser and plugin embodiments need not pretend to share a full runtime.
  - Centers frozen binding/session/message specs and conformance claims as the teachable core for first ports.
  - Keeps unsettled runtime details provisional, which helps host and helper evolution.
- Weaknesses:
  - Does not directly state what makes two heterogeneous embodiments one app rather than two coordinated apps.
  - Conformance-claim contents for partial implementations are not yet concrete enough for strong long-term audit.
  - Browser storage loss, helper replacement, and live-channel failure behavior are mostly outside the current framing.
- Risks:
  - Provisional porting guidance could be mistaken for a stable app-contract boundary.
  - Kernel/runtime/dispatcher ambiguity can produce incompatible embodiment claims.
  - Mallory can exploit loose claim semantics to blur UX sameness and protocol sameness.
- Open questions:
  - Is shared app identity anchored to one pCID-selected contract, a contract family, or another bundle of claims?
  - What minimum fields must each embodiment claim about implemented subsets, runtime constraints, and excluded obligations?
  - What local evidence should survive browser resets, helper upgrades, and host replacement for long-term audit?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-funas-kernel-porting-boundary` x `portable-signing-key-identity`

- Result path: `results/SIM-funas-kernel-porting-boundary/portable-signing-key-identity/openai-gpt-5.4-xhigh/20260521-234037.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=2 evolution_safety=4 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=2 risk_penalty=3
- Fitness: raw=21.00 normalized_0_100=52.50 confidence_0_1=0.79
- Rationale: Partial fit: this simulation is about provisional kernel/runtime porting boundaries rather than a full identity recipe, but it does help define which key-handling and storage duties belong in a first port and which signing-key continuity details should remain provisional.
- Strengths:
  - Strong provisional-vs-blocked framing matches the scenario's need for a pivotable v0 identity recipe.
  - Layer-boundary focus helps keep signing-key continuity distinct from display names and local usernames.
  - Explicit porter concerns include storage and key handling across different host embodiments.
- Weaknesses:
  - No concrete old-key/new-key linkage artifact is defined for Carol's audit.
  - No explicit browser-to-helper handshake or live-session continuity rule is given.
  - Adversarial handling for XSS, storage compromise, and Mallory's name-collision pressure is mostly implicit.
- Risks:
  - Porting-boundary guidance could be mistaken for a settled identity specification.
  - Guide language about browser key storage could overstate safety or understate XSS exposure.
  - Durable identity could drift back toward usernames or display labels if the boundary stays vague.
- Open questions:
  - What exact signed artifacts let Bob and Carol verify key continuity using only local evidence?
  - Which rotation and handshake details are stable protocol obligations versus host-local implementation choices?
  - What minimum frozen specs are enough for a first cross-embodiment signing-identity claim?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-funas-kernel-porting-boundary` x `kernel-porting-boundary`

- Result path: `results/SIM-funas-kernel-porting-boundary/kernel-porting-boundary/openai-gpt-5.4-xhigh/20260521-234037.json`
- Scores: scenario_fit=5 promisegrid_alignment=4 auditability=3 evolution_safety=5 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=3 risk_penalty=2
- Fitness: raw=29.00 normalized_0_100=73.00 confidence_0_1=0.72
- Rationale: Direct fit for the scenario and strong on provisional/frozen-boundary hygiene, but still a provisional question home rather than a concrete first-port contract.
- Strengths:
  - Directly targets the exact porting-boundary pressure and source feedback.
  - Separates wire-lab harness apparatus from the real porting target.
  - Makes provisional vs blocked guidance explicit until DR-davod closes.
  - ... 1 more
- Weaknesses:
  - Does not yet pin down the minimum frozen spec bundle for a first port.
  - Terminology and the exact runtime boundary remain unresolved.
  - Provides little detail on locally auditable provisional conformance evidence or failure handling.
- Risks:
  - Vague wording could let Mallory present the wire-lab harness as the de facto target.
  - Draft runtime obligations could be mistaken for stable API commitments.
  - Porters may over- or under-scope first implementations before the frozen minimum is named.
- Open questions:
  - Should the guide say kernel, runtime, dispatcher, handler host, or library?
  - What minimum frozen binding/session/message/CAS set and conformance claims define a first real port?
  - Which runtime, K1-K5 ingress, feed, CAS subtree, and app-layer details stay blocked versus provisional orientation?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

## Required JSON Shape

{"child_id":"SIM-jopuz-child-descriptive-design-slug","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
