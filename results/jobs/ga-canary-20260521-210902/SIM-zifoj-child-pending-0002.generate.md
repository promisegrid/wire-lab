# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Choose a descriptive `child_id` that starts with `SIM-zifoj-child-` and ends with a kebab-case design slug. Do not use generic `ga-child`, `pending`, or ordinal-only names.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260521-210902`
- Planned child ID prefix: `SIM-zifoj-child-`
- Temporary child ID: `SIM-zifoj-child-pending-0002`
- Temporary child path: `proposals/ga-canary-20260521-210902/simulations/SIM-zifoj-child-pending-0002/`
- Operation: `breed`
- Parent IDs: `SIM-funas-kernel-porting-boundary, SIM-robot-app-semantics-conformance`

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

### `simulations/SIM-robot-app-semantics-conformance/README.md`

```markdown
# SIM-robot: App semantics and conformance

This simulation is a provisional question home for App Dev feedback items
`FB-dodos`, `FB-hisis`, `FB-kutub`, `FB-gomod`, and `FB-tahof`. It tests what the
guide can say about app semantics and honest conformance before `DR-tuhaz`
settles the stable app-developer contract. Source: `DI-ragaz`.

## Question

Which app-facing semantic patterns are safe as provisional guide prose, and which
must remain blocked until stronger upstream decisions land? Source: `DI-ragaz`.

## Decision Axes

- **Vocabulary status:** promise, assertion, authorship, forwarding,
  conformance, capability, and witness language.
- **Local versus wire identity:** local IDs and storage handles may exist, but
  protocol-boundary identity must be spec-defined.
- **Partial conformance:** useful first slices can be honest if they do not claim
  full draft-spec behavior.
- **Provisional signing:** current signature carriage may be adapter-local until
  grid-envelope/signature decisions freeze.
- **Policy surface:** ingress models and economic patterns may be orientation or
  blocked, depending on what a candidate spec claims.

## Related Root Scenario

- `scenarios/app-semantics-partial-conformance/app-semantics-partial-conformance.md`

## Boundaries

This simulation does not define a universal app SDK, handler ABI, capability
token standard, or final witness format. It exists to keep guide wording honest
while draft specs and DRs remain open. Source: `DI-ragaz`.
```

### `simulations/SIM-robot-app-semantics-conformance/QUESTION.md`

```markdown
# Question

How should PromiseGrid guide authors describe provisional app semantics,
local-versus-wire identity, signed artifacts, and partial conformance without
claiming a stable app API before `DR-tuhaz` closes? Source: `DI-ragaz`.

Open decision points:

- Which vocabulary can be taught as careful orientation rather than final
  protocol semantics?
- What does a valid partial-conformance claim say, and what must it not claim?
- How should signed apps proceed when signature carriage remains simulation-level
  evidence rather than a frozen contract?
```

## Compact Fitness Evidence From This Run

### `SIM-robot-app-semantics-conformance` x `minimal-immutable-blob-app`

- Result path: `results/SIM-robot-app-semantics-conformance/minimal-immutable-blob-app/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=3 evolution_safety=4 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=3 risk_penalty=2
- Fitness: raw=20.00 normalized_0_100=57.00 confidence_0_1=0.74
- Rationale: Moderate fit. The simulation is strong at setting honest provisional app-semantics boundaries, especially around local-versus-wire identity and partial conformance, but it does not yet provide the concrete CAS, retention, retrieval, and audit-record semantics that this minimal immutable blob scenario pressures.
- Strengths:
  - Explicitly avoids premature stable app-contract claims.
  - Clear local-versus-wire identity boundary helps prevent overclaiming blob-hash semantics.
  - Evolution-friendly framing supports honest partial conformance while upstream decisions remain open.
- Weaknesses:
  - Does not define concrete content-identity versus availability, replication, or retention promises.
  - Provides no specific local promise-accounting records for years-later retrieval failure.
  - Leaves signed/witnessed artifact shape provisional rather than durable scenario evidence.
- Risks:
  - A returned hash could be misread as a durable retrieval guarantee.
  - Possession of a hash could be mistaken for authorization if layers remain implicit.
  - Adapter-local conventions may ossify into de facto wire semantics before DRs close.
- Open questions:
  - What exact promise may Bob safely make when returning a hash?
  - Is a hash only content identity, or may any layer treat it as an address or capability?
  - What records should Alice, Bob, and Carol keep when retention or retrieval fails years later?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-robot-app-semantics-conformance` x `app-semantics-partial-conformance`

- Result path: `results/SIM-robot-app-semantics-conformance/app-semantics-partial-conformance/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=5 promisegrid_alignment=4 auditability=3 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=1
- Fitness: raw=32.00 normalized_0_100=80.00 confidence_0_1=0.82
- Rationale: Strong match: the simulation is explicitly about honest provisional app semantics before DR-tuhaz, with clear separation between local IDs, protocol-boundary identity, adapter-local signing, and any final app contract. It loses points because it frames the right questions and boundaries but does not yet specify the exact claim text, locally auditable evidence, or migration/failure procedure Bob and Carol would use.
- Strengths:
  - Directly targets the scenario's partial-conformance pressure.
  - Very clear local-versus-wire boundary discipline.
  - Explicitly defers unstable signature, capability, and witness surfaces instead of freezing them early.
- Weaknesses:
  - No concrete partial-conformance claim template or verification checklist is defined.
  - Peer-local observable artifacts for Bob and Carol are implied rather than specified.
  - 100-year, sparse-knowledge, and no-central-authority pressures are not exercised explicitly in the sim docs.
- Risks:
  - Provisional guide language could ossify into a de facto app API before upstream decisions close.
  - Adapter-local signature carriage could be mistaken for stable wire-level conformance.
  - Local implementation IDs could leak into interoperability claims if the boundary is applied loosely.
- Open questions:
  - What exact wording or schema makes a partial-conformance claim honest and non-overclaiming?
  - What identity is authoritative at the protocol boundary once the app contract is settled?
  - What migration path moves from adapter-local signing and provisional witness language to a frozen contract?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-robot-app-semantics-conformance` x `device-bound-agent-physical-effect`

- Result path: `results/SIM-robot-app-semantics-conformance/device-bound-agent-physical-effect/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=2 promisegrid_alignment=2 auditability=3 evolution_safety=4 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=2 risk_penalty=3
- Fitness: raw=20.00 normalized_0_100=50.00 confidence_0_1=0.78
- Rationale: Useful for honest app/device-adapter conformance framing, especially local-vs-wire identity and provisional dependency claims, but it does not supply the delegated-authority, replay-safe at-most-once, receipt, or break-witness semantics this scenario needs.
- Strengths:
  - Explicitly separates local handles from protocol-boundary identity.
  - Encourages honest partial-conformance claims instead of implying a stable app API.
  - Keeps signing and witness semantics provisional, which is safer for later evolution.
- Weaknesses:
  - Does not define evidence that Bob delegated a device-bound agent to a host or daemon.
  - Does not specify deduplication or at-most-once handling for non-idempotent physical effects after restart or replay.
  - Leaves receipt and break-witness structure unresolved for Carol's audit.
  - ... 1 more
- Risks:
  - Provisional adapter behavior could be overread as enough to justify unsafe physical-effect execution.
  - Replay after restart could produce duplicate labels or ambiguous receipts.
  - Dependency claims around CUPS, libusb, IIO, i2c, IPP, or vendor SDKs could be overstated as PromiseGrid conformance.
- Open questions:
  - What durable receipt identifies one physical effect across replay and daemon restart?
  - How is delegated device authority expressed independently of local driver handles or process identity?
  - What minimum conformance language should cover host-driver dependencies without implying a frozen wire contract?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-robot-app-semantics-conformance` x `live-crdt-audit-publication`

- Result path: `results/SIM-robot-app-semantics-conformance/live-crdt-audit-publication/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=2 evolution_safety=4 layer_boundary_clarity=5 failure_handling=1 implementation_plausibility=2 risk_penalty=2
- Fitness: raw=18.00 normalized_0_100=58.00 confidence_0_1=0.74
- Rationale: This simulation fits the scenario mainly as a boundary-setting exercise: it supports honest separation between provisional live CRDT behavior and durable PromiseGrid audit publication, but it does not specify the reliable live binding, cited object shape, or witness details that the scenario needs.
- Strengths:
  - Strong partial-conformance framing supports claims like 'durable audit publication supported' without implying grid-native live sync.
  - Clear local-versus-wire identity language helps keep session-local CRDT handles from being mistaken for stable protocol identity.
  - Explicit provisional scope before DR-tuhaz improves evolution safety and leaves room for later transport and witness decisions.
- Weaknesses:
  - No concrete reliable ordered live binding or reconnection semantics for the partitioned live channel.
  - Does not define the exact content-addressed snapshot or witness object the save-time audit message should cite.
  - Leaves signature carriage and witness format unresolved at adapter/simulation level, limiting durable audit specificity.
- Risks:
  - Readers may overread orientation prose as approval of a specific live transport or future live pCID shape.
  - Adapter-local IDs or signatures could harden into de facto wire semantics before upstream decisions freeze.
  - Relying on prose-only separation between live and audit claims may be too weak for future conformance testing.
- Open questions:
  - Should guidance explicitly keep live CRDT sync off-grid until a reliable ordered binding exists, or sketch a future live pCID shape?
  - What exact content-addressed object should the audit message cite at save time?
  - How should provisional signing be phrased so audit-layer signatures do not imply live-channel conformance?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-robot-app-semantics-conformance` x `multi-embodiment-app-identity`

- Result path: `results/SIM-robot-app-semantics-conformance/multi-embodiment-app-identity/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=5 promisegrid_alignment=3 auditability=3 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=3 risk_penalty=2
- Fitness: raw=23.00 normalized_0_100=70.00 confidence_0_1=0.75
- Rationale: Strong scenario match: it directly addresses one app identity across heterogeneous embodiments through spec-defined wire identity and honest partial-conformance claims, but it remains provisional and under-specifies durable witness, peer-local evidence, and failure recovery.
- Strengths:
  - Anchors 'same app' in shared protocol contract semantics rather than shared UX or branding alone.
  - Explicitly separates local/runtime handles from protocol-boundary identity.
  - Supports embodiment-specific partial-conformance claims instead of forcing false full-conformance claims.
- Weaknesses:
  - Does not yet define the frozen shared contract or contract-family mechanism that would make cross-embodiment identity durable.
  - Witness and signature carriage remain provisional, limiting long-term audit strength.
  - Recovery and evidence behavior for storage loss, host replacement, helper upgrades, or channel failure are mostly implicit.
- Risks:
  - Mallory can exploit loose 'same app' language if compatible contract claims are not tightly scoped.
  - Adapter-local signing may not survive long-term audit or embodiment replacement.
  - Provisional guide prose could harden into de facto API promises before upstream DRs settle the boundary.
- Open questions:
  - What minimum fields must each embodiment publish in its conformance claim?
  - How is one app identity bound to a contract or contract family without central registration?
  - What durable audit artifact survives browser storage loss, helper upgrades, and future signature-envelope changes?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-robot-app-semantics-conformance` x `portable-signing-key-identity`

- Result path: `results/SIM-robot-app-semantics-conformance/portable-signing-key-identity/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=3 promisegrid_alignment=4 auditability=3 evolution_safety=5 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=3 risk_penalty=2
- Fitness: raw=23.00 normalized_0_100=66.00 confidence_0_1=0.78
- Rationale: Strong as provisional guide discipline: it clearly separates local handles from wire identity and keeps signing claims pivotable before DR-tuhaz. It is weaker as a direct answer to this scenario because it does not yet define the concrete rotation evidence, cross-host continuity artifact, or browser-helper security guidance the scenario asks for.
- Strengths:
  - Clear local-vs-wire identity boundary helps resist display-name and local-username confusion.
  - Provisional signing stance supports a pivotable v0 recipe instead of hardcoding forever cryptography.
  - Honest partial-conformance framing reduces overclaiming before the app contract freezes.
- Weaknesses:
  - No concrete artifact or witness chain links old and new keys for later audit.
  - No specific browser versus plugin-host storage and XSS guidance.
  - No settled handshake or session-continuity recipe across embodiments.
- Risks:
  - Adapter-local signatures could be mistaken for durable protocol identity.
  - Guide prose could still blur presentation hints with long-lived identity if wording stays loose.
- Open questions:
  - What minimal v0 record should bind old and new signing keys?
  - Which identity and signature-carriage details remain provisional until DR-tuhaz closes?
  - How should WebCrypto/IndexedDB versus helper filesystem storage risk be described without overstating security?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-robot-app-semantics-conformance` x `kernel-porting-boundary`

- Result path: `results/SIM-robot-app-semantics-conformance/kernel-porting-boundary/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=4 evolution_safety=4 layer_boundary_clarity=2 failure_handling=2 implementation_plausibility=3 risk_penalty=3
- Fitness: raw=22.00 normalized_0_100=55.00 confidence_0_1=0.82
- Rationale: Useful for honest provisional conformance language, local-versus-wire boundaries, and stating what is deferred, but it is centered on app semantics under unsettled app-contract work rather than the kernel/runtime/harness port boundary this scenario probes.
- Strengths:
  - Explicitly supports honest partial-conformance claims instead of overstating draft-spec coverage.
  - Keeps local-versus-wire identity and other provisional semantics clearly bounded.
  - Has strong evolution posture by deferring stable claims until upstream decisions freeze.
- Weaknesses:
  - Focuses on app semantics more than on the first real infrastructure port target.
  - Does not clearly separate wire-lab harness apparatus from the actual porting surface.
  - Leaves kernel/runtime/dispatcher/library terminology and minimum binding/session/message/CAS obligations underspecified.
- Risks:
  - Port authors could mistake app-facing conformance guidance for sufficient kernel-porting guidance.
  - Partial-conformance language may hide missing ingress, feed, session, or CAS obligations.
  - Draft terminology could harden the wrong layer boundary before DR-davod closes.
- Open questions:
  - What minimum binding, session, message, CAS, and runtime surface counts as a valid first port before freeze?
  - Which guide term should be standardized for kernel, runtime, dispatcher, handler host, and library?
  - Which K1-K5 details should stay blocked versus provisional orientation in a first porting claim?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-funas-kernel-porting-boundary` x `minimal-immutable-blob-app`

- Result path: `results/SIM-funas-kernel-porting-boundary/minimal-immutable-blob-app/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=2 promisegrid_alignment=2 auditability=2 evolution_safety=5 layer_boundary_clarity=5 failure_handling=1 implementation_plausibility=2 risk_penalty=2
- Fitness: raw=22.00 normalized_0_100=55.00 confidence_0_1=0.79
- Rationale: Strong as boundary-setting evidence: it clearly separates provisional porting obligations from CAS/app semantics. Weak as a direct answer to the blob-app scenario because it does not define the retrieval promise, long-term availability story, or local failure records.
- Strengths:
  - Explicit focus on CAS/app/runtime layer separation and avoiding a monolithic harness clone.
  - Very strong on provisional-versus-stable guidance, which helps prevent freezing the wrong porting target.
  - Calls out storage, ingress, feeds, CAS subtree, and conformance claims as boundary questions a porter must not absorb blindly.
- Weaknesses:
  - Does not specify what "hash in -> blob out" promises or which actor makes that promise.
  - Does not address retention changes, incomplete caches, replication, or stale availability claims in the scenario.
  - Provides no concrete local promise-accounting records for Alice, Bob, and Carol when retrieval fails.
- Risks:
  - A porter could mistake this provisional runtime boundary work for a stable blob-app contract.
  - If CAS identity and storage availability stay under-specified, implementations may conflate address, capability, and authorization.
- Open questions:
  - Which minimal frozen specs are enough for a first port to support CAS lookup without implying durable storage guarantees?
  - Which blob semantics belong in app-layer contracts versus runtime or session conformance?
  - What evidence should a port record when content is identified correctly but unavailable years later?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-funas-kernel-porting-boundary` x `app-semantics-partial-conformance`

- Result path: `results/SIM-funas-kernel-porting-boundary/app-semantics-partial-conformance/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=3 evolution_safety=5 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=3 risk_penalty=3
- Fitness: raw=20.00 normalized_0_100=57.00 confidence_0_1=0.76
- Rationale: Strong on provisional-versus-blocked teaching and layer boundaries, but only a partial match for this scenario because it is centered on kernel/runtime porting rather than app-side partial-conformance claims.
- Strengths:
  - Clear emphasis on layer boundaries and avoiding a monolithic harness-clone target.
  - Strong evolution safety: it tries to teach only what is safe before the kernel/runtime boundary is settled.
  - Anchors early conformance discussion to frozen specs and explicit obligations.
- Weaknesses:
  - Does not directly define honest app-side partial-conformance wording.
  - Protocol-boundary identity for app semantics is not made explicit.
  - Signature carriage, capability, and witness language remain underspecified.
- Risks:
  - Kernel/runtime port guidance could be mistaken for an app contract.
  - Unsettled boundary decisions may freeze the wrong conformance story if taught too concretely.
  - Local implementation shortcuts may still leak into outward claims without stricter claim language.
- Open questions:
  - What exact wording makes an honest partial-conformance claim for an app that only supports a first slice?
  - Which identity is authoritative at the protocol boundary when the app uses local IDs internally?
  - Which current signature-carriage and capability/witness terms are safe orientation versus blocked?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-funas-kernel-porting-boundary` x `device-bound-agent-physical-effect`

- Result path: `results/SIM-funas-kernel-porting-boundary/device-bound-agent-physical-effect/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=4 promisegrid_alignment=4 auditability=2 evolution_safety=4 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=3 risk_penalty=2
- Fitness: raw=21.00 normalized_0_100=60.00 confidence_0_1=0.77
- Rationale: Strong as a boundary test: the scenario directly pressures what a PromiseGrid port must standardize versus what should stay in host/device adapters and app logic. Weak as an end-to-end answer, because delegation proof, durable effect dedup, receipts, and break-witness behavior for restart/replay are still largely unspecified.
- Strengths:
  - Decision axes explicitly cover minimum viable port, runtime obligations, and conformance claims.
  - The simulation is well suited to separating core PromiseGrid duties from CUPS/libusb/IIO/i2c/IPP/vendor-SDK style host dependencies.
  - Its provisional-versus-blocked framing helps avoid freezing today's device stack into a false 100-year port target.
- Weaknesses:
  - It does not yet specify evidence that Bob delegated the device-bound agent.
  - It does not define durable effect IDs, receipts, or break-witnesses for replay after restart.
  - At-most-once handling for non-idempotent physical effects remains underdefined.
- Risks:
  - Implementers may push critical safety and accounting semantics above the boundary and end up with inconsistent device-agent behavior across ports.
  - Unsettled kernel/runtime terminology could still be mistaken for a final implementation target before DR-davod closes.
- Open questions:
  - Which durable local record or effect ID prevents duplicate printing after daemon restart?
  - What exact evidence proves Bob delegated control of the printer or sensor to a specific agent instance?
  - How should conformance text describe external CUPS/libusb/IIO/i2c/IPP/vendor SDK dependencies and their failure surfaces?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-funas-kernel-porting-boundary` x `live-crdt-audit-publication`

- Result path: `results/SIM-funas-kernel-porting-boundary/live-crdt-audit-publication/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=4 promisegrid_alignment=4 auditability=3 evolution_safety=5 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=3 risk_penalty=2
- Fitness: raw=25.00 normalized_0_100=71.00 confidence_0_1=0.80
- Rationale: Strong as a boundary-setting simulation for this scenario: it can clearly separate sub-second live CRDT transport from durable audit publication and keep unsettled live/runtime details out of the current core port target, but it does not yet specify the live binding, exact cited audit object, or detailed reconnect behavior.
- Strengths:
  - Directly addresses the scenario's main pressure to keep live transport claims separate from durable audit-layer claims.
  - Very strong on provisional-versus-stable guidance, which supports 100-year evolution without freezing the wrong porting API.
  - Supports durable, content-addressed audit publication without assuming central authority or globally complete state.
- Weaknesses:
  - Does not define the reliable ordered live channel needed for the CRDT session itself.
  - Leaves the exact snapshot or envelope cited by the save-time audit message underspecified.
  - Provides only limited actor-local failure and recovery detail for partition, reorder, and delayed publication.
- Risks:
  - Porters may mistake provisional boundary guidance for a settled kernel or runtime API.
  - Documentation could still blur optional live-channel behavior with core audit-layer conformance.
  - Developers may overfit to harness-era drafts instead of a minimal frozen-spec port target.
- Open questions:
  - Should live CRDT sync stay explicitly off-grid until a reliable binding is frozen?
  - What exact CAS object or snapshot manifest should the durable audit message cite?
  - How should live-channel conformance claims be versioned and reported separately from audit-layer claims?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-funas-kernel-porting-boundary` x `multi-embodiment-app-identity`

- Result path: `results/SIM-funas-kernel-porting-boundary/multi-embodiment-app-identity/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=3 risk_penalty=3
- Fitness: raw=21.00 normalized_0_100=60.00 confidence_0_1=0.76
- Rationale: Partial fit. The simulation is strong on separating app-layer contract identity from embodiment-specific runtime obligations via conformance claims and explicit layer boundaries, but it does not yet define the stable claim schema or recovery semantics needed to fully settle heterogeneous one-app identity.
- Strengths:
  - Explicit conformance-claim and minimum-port focus helps browser and plugin embodiments declare different implemented subsets.
  - Strong layer-boundary emphasis helps keep shared app contract distinct from host, storage, and runtime constraints.
  - Provisional-versus-blocked framing improves evolution safety while the kernel/runtime boundary remains open.
- Weaknesses:
  - It does not directly define what makes multiple embodiments one app at the app-identity layer.
  - Browser storage loss, host replacement, and helper-upgrade continuity are only indirectly covered.
  - Adversarial reading of ambiguous claims and sparse-knowledge evidence are not worked through in detail.
- Risks:
  - Ambiguous boundary language could let Mallory blur 'same UX app' and 'same protocol contract.'
  - Prematurely teaching an unstable runtime surface as the port target could freeze the wrong conformance expectations.
  - Without a stable claim schema, partial implementations may overclaim compatibility.
- Open questions:
  - What exact fields must each embodiment's conformance claim include to bind it to the shared app contract?
  - Is one pCID-selected contract enough, or is a contract family with required and optional roles needed?
  - What durable local records preserve auditability after browser state loss or helper replacement?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-funas-kernel-porting-boundary` x `portable-signing-key-identity`

- Result path: `results/SIM-funas-kernel-porting-boundary/portable-signing-key-identity/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=2 evolution_safety=4 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=2 risk_penalty=2
- Fitness: raw=17.00 normalized_0_100=49.00 confidence_0_1=0.74
- Rationale: Strong as a boundary-setting probe, but weak as a full portable-identity answer: it can frame what key handling, storage, and rotation guidance is provisional versus porting obligation, yet it does not define the concrete continuity evidence or host-specific security treatment the scenario needs.
- Strengths:
  - Directly focuses on what the guide can safely teach now, matching the scenario's provisional-v0 pressure.
  - Explicitly puts key handling, storage, and layer boundaries in scope, so browser/helper differences are at least discussable.
  - Provisional-versus-blocked framing supports future crypto and runtime evolution instead of freezing a premature target.
- Weaknesses:
  - No concrete mechanism for linking old and new signing keys for Bob or Carol to verify.
  - Does not clearly specify how durable signing identity differs from display names or local usernames.
  - Provides little explicit treatment of XSS, host compromise, or name-collision adversarial pressure.
- Risks:
  - Porters may infer unstable host-local key-storage practices as durable identity semantics.
  - Guide text could freeze the wrong security boundary before the kernel/runtime split is settled.
- Open questions:
  - Which identity and rotation claims belong in the port boundary versus a higher identity layer?
  - What minimal signed artifact should link rotated keys across browser and helper without a central registry?
  - How should browser storage and XSS caveats be stated without overstating browser trust?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-funas-kernel-porting-boundary` x `kernel-porting-boundary`

- Result path: `results/SIM-funas-kernel-porting-boundary/kernel-porting-boundary/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=5 promisegrid_alignment=4 auditability=3 evolution_safety=5 layer_boundary_clarity=4 failure_handling=3 implementation_plausibility=3 risk_penalty=2
- Fitness: raw=30.00 normalized_0_100=75.00 confidence_0_1=0.82
- Rationale: Strong direct fit to the scenario: it explicitly frames the porting-boundary question, separates provisional guidance from future frozen obligations, and resists treating the wire-lab harness as the porting target. Its main limitation is that it preserves the question space more than it defines a concrete first-port contract.
- Strengths:
  - Exact match to the kernel/runtime porting-boundary pressure.
  - Explicit provisional-vs-blocked framing supports later freeze without locking in a false API.
  - Strong emphasis on layer separation and avoiding monolithic harness-clone guidance.
- Weaknesses:
  - Does not yet define a concrete minimum viable port checklist.
  - Audit evidence and conformance mechanics are named but not operationalized in detail.
  - Sparse-knowledge, no-central-authority, and peer-local-accounting implications are mostly implicit.
- Risks:
  - Provisional terminology could still ossify into a misleading stable target.
  - Draft runtime or K1-K5 details could be over-read as current obligations.
  - 'Kernel' framing may bias implementers toward too-large or too-monolithic a boundary.
- Open questions:
  - Should the guide standardize on kernel, runtime, dispatcher, handler host, or library terminology?
  - Which frozen binding/session/message specs are the true minimum for a first real port?
  - Which runtime obligations are mandatory now versus deferred until DR-davod closes?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

## Required JSON Shape

{"child_id":"SIM-zifoj-child-descriptive-design-slug","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
