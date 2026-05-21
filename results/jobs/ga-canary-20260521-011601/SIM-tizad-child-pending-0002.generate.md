# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Choose a descriptive `child_id` that starts with `SIM-tizad-child-` and ends with a kebab-case design slug. Do not use generic `ga-child`, `pending`, or ordinal-only names.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260521-011601`
- Planned child ID prefix: `SIM-tizad-child-`
- Temporary child ID: `SIM-tizad-child-pending-0002`
- Temporary child path: `proposals/ga-canary-20260521-011601/simulations/SIM-tizad-child-pending-0002/`
- Operation: `breed`
- Parent IDs: `SIM-robot-app-semantics-conformance, SIM-kugap-live-sync-audit-split`

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

### `simulations/SIM-kugap-live-sync-audit-split/README.md`

```markdown
# SIM-kugap: Live sync and audit split

This simulation is a provisional question home for `FB-hurit` and `FB-nilat`:
apps that need low-latency reliable live state but still want durable
PromiseGrid audit publication. Source: `DI-ragaz`.

## Question

Should the guide describe a future reliable live-state protocol shape, an
off-grid live channel plus PromiseGrid audit publication pattern, or both as
explicitly provisional? Source: `DI-ragaz`.

## Decision Axes

- **Live transport needs:** reliable, in-order, frame-preserving delivery with
  sub-second latency for CRDT sync or presence.
- **Audit protocol needs:** durable, replayable publication of snapshots,
  milestones, receipts, or break-witnesses.
- **Layer separation:** live state and durable audit may use different pCIDs,
  timings, durability promises, and conformance claims.
- **Group-session role:** group-session may be an audit layer without being the
  live transport.
- **Failure behavior:** dropped live frames, stale audit snapshots, and partition
  recovery must be explicit.

## Related Root Scenario

- `scenarios/live-crdt-audit-publication/live-crdt-audit-publication.md`

## Boundaries

This simulation does not define a reliable low-latency binding. It tests whether
the guide can name the split honestly without implying that `udp-feed` or
`group-session` already solves live CRDT transport. Source: `DI-ragaz`.
```

### `simulations/SIM-kugap-live-sync-audit-split/QUESTION.md`

```markdown
# Question

How should PromiseGrid distinguish live-state protocol claims from durable audit
publication claims for real-time apps such as collaborative editors, presence,
multiplayer state, and telemetry? Source: `DI-ragaz`.

Open decision points:

- Is the guide-safe posture "live state stays off-grid for now, audit publishes
  to group-session" or a provisional future live protocol sketch?
- What does an audit message cite: snapshot CID, CRDT save blob, operation log,
  milestone receipt, or another object?
- How should conformance claims avoid conflating live channel behavior with
  group-session audit behavior?
```

## Compact Fitness Evidence From This Run

### `SIM-robot-app-semantics-conformance` x `minimal-immutable-blob-app`

- Result path: `results/SIM-robot-app-semantics-conformance/minimal-immutable-blob-app/openai-gpt-5.4-xhigh/20260521-011601.json`
- Scores: scenario_fit=4 promisegrid_alignment=4 auditability=3 evolution_safety=5 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=2
- Fitness: raw=30.00 normalized_0_100=75.00 confidence_0_1=0.81
- Rationale: Strong on the scenario's semantics and boundary questions: it clearly frames provisional app vocabulary, local-versus-wire identity, and honest partial conformance. It is weaker on concrete blob retrieval failure behavior, retention/replication promises, and durable peer-local audit records.
- Strengths:
  - Directly addresses app-facing semantics, conformance language, and capability vocabulary relevant to a minimal CAS app.
  - Excellent boundary clarity between local storage handles and protocol-level identity, and between adapter-local signing and frozen wire claims.
  - High evolution safety because it explicitly avoids premature stable app/API commitments and permits honest partial conformance.
- Weaknesses:
  - Does not define concrete long-term retrieval, retention, replication, or discovery promises for the blob.
  - Does not specify the local promise-accounting records Alice, Bob, and Carol should keep when retrieval fails.
  - Leaves hash-as-address versus hash-as-capability semantics unresolved.
- Risks:
  - Provisional guide prose could still be misread as a stable app contract.
  - Unsettled signature and witness carriage may limit cross-site, years-later audit interoperability.
  - Implementers may still conflate content identity with availability or authorization if the boundary is not made explicit enough.
- Open questions:
  - Who exactly promises hash-to-blob retrieval, and with what retention or availability scope?
  - Is possession of a hash ever authorization, or only a reference plus app-level policy?
  - What minimal local records should each peer retain to audit stale availability claims or failed retrieval years later?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-robot-app-semantics-conformance` x `app-semantics-partial-conformance`

- Result path: `results/SIM-robot-app-semantics-conformance/app-semantics-partial-conformance/openai-gpt-5.4-xhigh/20260521-011601.json`
- Scores: scenario_fit=5 promisegrid_alignment=4 auditability=4 evolution_safety=5 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=2
- Fitness: raw=28.00 normalized_0_100=80.00 confidence_0_1=0.85
- Rationale: Direct match for the scenario: it explicitly frames honest partial conformance, separates local IDs from protocol-boundary identity, and marks signing/capability/witness surfaces as provisional rather than pretending a frozen app contract exists. It is weaker on concrete claim format, local evidence shape, and adversarial walkthroughs.
- Strengths:
  - Directly targets provisional app semantics and honest partial-conformance claims.
  - Clearly separates local implementation IDs and handles from spec-defined wire identity.
  - Treats signing, capability, witness, and policy surfaces as provisional or blocked pending upstream decisions.
- Weaknesses:
  - Does not define a concrete partial-conformance claim structure Bob can verify from wire artifacts.
  - Limited detail on Carol's local audit records and promise-accounting evidence.
  - Adversarial overclaim and interoperability failure modes are acknowledged more than operationalized.
- Risks:
  - Guide prose could be mistaken for a de facto stable app API if provisional labels are weak.
  - Adapter-local signature carriage could ossify before envelope and signature decisions freeze.
  - Ambiguity around protocol-boundary identity could lead to accidental overclaims of interoperability.
- Open questions:
  - What exact wording makes a B-side partial-conformance claim honest and non-overclaiming?
  - Which identity fields are authoritative at the protocol boundary versus explicitly local-only?
  - What minimum locally recordable evidence should Bob and Carol use to audit provisional signing and capability or witness behavior?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-robot-app-semantics-conformance` x `device-bound-agent-physical-effect`

- Result path: `results/SIM-robot-app-semantics-conformance/device-bound-agent-physical-effect/openai-gpt-5.4-xhigh/20260521-011601.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=2 evolution_safety=5 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=2 risk_penalty=3
- Fitness: raw=17.00 normalized_0_100=48.60 confidence_0_1=0.84
- Rationale: Moderate fit: the simulation is strong on provisional app-semantics hygiene—local-vs-wire identity, partial conformance, and dependency disclosure—but it does not yet define the delegation, effect identity, receipts, or break-witnesses needed for restart-safe device control.
- Strengths:
  - Clear local-vs-wire boundary helps separate host/device adapters from protocol claims.
  - Honest partial-conformance framing is well suited to printer/sensor apps with heavy local dependencies.
  - Explicitly provisional scope reduces premature lock-in while upstream app-contract decisions remain open.
- Weaknesses:
  - Delegation evidence among owner, operator host, and device-bound agent is not defined.
  - No protocol-visible effect identity or durable receipt model for replay-dedup after restart.
  - Long-horizon interpretation of physical-effect evidence after device/driver obsolescence is underspecified.
- Risks:
  - Teams may mistake local adapter behavior for stable wire semantics and overclaim conformance.
  - Mallory replay after daemon restart could trigger duplicate non-idempotent actions.
  - Adapter-local signing or witness conventions may fail 100-year audit needs.
- Open questions:
  - What artifact proves Bob delegated printer/sensor authority without a central registry?
  - What request or effect identifier survives restart and replay so Alice and Carol can verify at-most-once behavior?
  - What receipt or break-witness must be archived so later auditors can interpret the event after current drivers and vendor SDKs disappear?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-robot-app-semantics-conformance` x `live-crdt-audit-publication`

- Result path: `results/SIM-robot-app-semantics-conformance/live-crdt-audit-publication/openai-gpt-5.4-xhigh/20260521-011601.json`
- Scores: scenario_fit=4 promisegrid_alignment=4 auditability=3 evolution_safety=5 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=3 risk_penalty=2
- Fitness: raw=30.00 normalized_0_100=75.00 confidence_0_1=0.78
- Rationale: Good fit for this pressure: the simulation is explicitly about honest app semantics and partial conformance before DR-tuhaz, so it can cleanly separate low-latency live CRDT behavior from durable PromiseGrid audit publication. It scores lower where the reliable live binding, exact cited audit object, and frozen witness/signature contract are still unresolved.
- Strengths:
  - Directly supports honest partial-conformance claims instead of pretending the live channel is already a stable PromiseGrid wire contract.
  - Very strong local-versus-wire and live-versus-audit layer separation.
  - Conservative provisional stance is evolution-safe while upstream app and signature decisions remain open.
  - ... 1 more
- Weaknesses:
  - No concrete reliable ordered live binding is defined for the CRDT path.
  - The exact durable object the audit message should cite is still underspecified.
  - Final witness and signature carriage remain provisional rather than frozen.
  - ... 1 more
- Risks:
  - Readers could misread provisional guide prose as blessing an unfrozen live wire contract.
  - Partial-conformance language may let implementers overclaim interoperability across live and audit layers.
  - Adapter-local signing may drift from the eventual frozen envelope/signature model.
- Open questions:
  - Should live CRDT state stay explicitly off-grid until a reliable binding exists, or should a future live pCID shape be sketched now?
  - What exact CAS object should the audit message cite: a CRDT snapshot, an op-log checkpoint, a rendered document artifact, or something else?
  - How should adapter-local signatures migrate to a frozen envelope/signature contract without breaking old audit trails?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-robot-app-semantics-conformance` x `multi-embodiment-app-identity`

- Result path: `results/SIM-robot-app-semantics-conformance/multi-embodiment-app-identity/openai-gpt-5.4-xhigh/20260521-011601.json`
- Scores: scenario_fit=4 promisegrid_alignment=4 auditability=3 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=29.00 normalized_0_100=72.50 confidence_0_1=0.76
- Rationale: Good fit for the scenario's core issue: the simulation explicitly tackles local-vs-wire identity, honest partial conformance, and provisional app semantics across heterogeneous embodiments. It is strong on boundary discipline and avoiding overclaiming, but it remains intentionally provisional and does not yet freeze the durable same-app witness and long-term audit details this scenario pressures.
- Strengths:
  - Clear separation between host-local handles and protocol-boundary identity.
  - Supports scoped partial-conformance claims for different embodiments of one logical app.
  - Cautious provisional stance reduces premature claims of a stable universal app API.
- Weaknesses:
  - Does not yet define a frozen mechanism for proving that multiple embodiments are the same app.
  - Long-term witness/signature carriage remains unresolved.
  - Scenario-specific failure cases like storage loss, stale claims, and upgrade divergence are only partly covered.
- Risks:
  - Branding or shared UX could be mistaken for shared protocol identity.
  - Provisional guide language may be read as a stable contract before DR-tuhaz/DR-davod land.
  - Audit evidence may be too weak across host replacement or long time horizons until witness format freezes.
- Open questions:
  - What minimum fields must each embodiment's conformance claim include?
  - Is one app identity anchored by one contract identifier, a contract family, or another durable relation?
  - What durable evidence survives browser storage loss, helper upgrades, and long-term audit?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-robot-app-semantics-conformance` x `portable-signing-key-identity`

- Result path: `results/SIM-robot-app-semantics-conformance/portable-signing-key-identity/openai-gpt-5.4-xhigh/20260521-011601.json`
- Scores: scenario_fit=4 promisegrid_alignment=4 auditability=2 evolution_safety=5 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=2 risk_penalty=2
- Fitness: raw=27.00 normalized_0_100=68.00 confidence_0_1=0.76
- Rationale: Strong conceptual fit on provisional identity semantics, local-versus-wire identity, and provisional signing. Fitness stays moderate because the simulation does not yet define the concrete rotation proof, cross-host session evidence, or browser-storage/XSS guidance this scenario needs.
- Strengths:
  - Explicit separation of local storage handles from protocol-boundary identity.
  - Very strong evolution posture: a pivotable v0 is allowed without freezing forever cryptography.
  - Honest partial-conformance framing helps avoid overstating app identity guarantees before DR-tuhaz closes.
- Weaknesses:
  - No concrete artifact shape for proving continuity from the old signing key to the new one.
  - Browser and plugin-host storage caveats, including XSS exposure, are not spelled out.
  - Witness, envelope, and handshake details remain intentionally unresolved.
- Risks:
  - A provisional recipe could harden into a de facto identity contract too early.
  - Readers may still confuse display names or local usernames with durable identity if examples are sloppy.
  - Adapter-local signature carriage may be too weak for later cross-host audit and interoperability.
- Open questions:
  - What minimal signed artifact links old and new keys in a way Bob and Carol can verify locally?
  - What provisional cross-host handshake is safe to describe before envelope and witness formats freeze?
  - How should the guide discuss browser key storage and XSS risk without overstating security?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-robot-app-semantics-conformance` x `kernel-porting-boundary`

- Result path: `results/SIM-robot-app-semantics-conformance/kernel-porting-boundary/openai-gpt-5.4-xhigh/20260521-011601.json`
- Scores: scenario_fit=1 promisegrid_alignment=4 auditability=3 evolution_safety=4 layer_boundary_clarity=2 failure_handling=1 implementation_plausibility=1 risk_penalty=3
- Fitness: raw=18.00 normalized_0_100=45.00 confidence_0_1=0.76
- Rationale: Strong on honest provisional conformance boundaries, but largely off-target for defining the first kernel/runtime port boundary and separating harness apparatus from the real porting target.
- Strengths:
  - Explicitly limits claims to provisional guidance rather than frozen obligations.
  - Supports honest partial-conformance claims and local-versus-wire identity caution.
  - Keeps a migration path to later stable-spec conformance.
- Weaknesses:
  - Does not define a minimum viable kernel/runtime/dispatcher/library porting target.
  - Does not directly separate wire-lab harness apparatus from real port obligations.
  - Provides little guidance on ingress, feed, CAS, or session scope.
- Risks:
  - Readers may mistake app-level conformance language for infrastructure-port coverage.
  - Mallory's 'copy the harness' framing remains insufficiently rebutted.
- Open questions:
  - Can the simulation's partial-conformance language be safely reused for first-port claims without implying full kernel/runtime coverage?
  - Which terms and K1-K5 obligations should remain blocked versus provisional until DR-davod closes?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kugap-live-sync-audit-split` x `minimal-immutable-blob-app`

- Result path: `results/SIM-kugap-live-sync-audit-split/minimal-immutable-blob-app/openai-gpt-5.4-xhigh/20260521-011601.json`
- Scores: scenario_fit=1 promisegrid_alignment=3 auditability=3 evolution_safety=3 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=3 risk_penalty=2
- Fitness: raw=18.00 normalized_0_100=51.00 confidence_0_1=0.86
- Rationale: Clear and useful on separating durable audit publication from live transport, but only a weak direct fit for the minimal immutable blob app. It helps prevent conflating identity/audit with transport guarantees, yet it does not define the core blob-app promise of hash-to-bytes retrieval across time, sites, retention changes, and partial caches.
- Strengths:
  - Strong layer and conformance-boundary clarity.
  - Provides a durable audit-publication concept for long-term evidence about published snapshots or blobs.
  - Explicitly names stale audit and partition-style failure modes instead of assuming perfect sync.
- Weaknesses:
  - Live-state concerns dominate a scenario that is really about minimal CAS retrieval semantics.
  - Does not specify who promises availability, discovery, replication, retention, or authorization for a blob addressed by hash.
  - Missing explicit peer-local promise-accounting records for Alice, Bob, and Carol when retrieval fails years later.
- Risks:
  - Audit publication of a CID could be mistaken for a promise that the blob will remain fetchable.
  - Readers may wrongly treat possession of a hash as sufficient read authority or discovery capability.
- Open questions:
  - What exact promise should Bob make when returning a hash for an immutable blob?
  - What receipts or local records should Alice, Bob, and Carol keep when retention policy changes or cross-site retrieval fails?
  - Should minimal blob apps use a simpler CAS/storage pattern distinct from the live-sync plus audit split?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kugap-live-sync-audit-split` x `app-semantics-partial-conformance`

- Result path: `results/SIM-kugap-live-sync-audit-split/app-semantics-partial-conformance/openai-gpt-5.4-xhigh/20260521-011601.json`
- Scores: scenario_fit=4 promisegrid_alignment=4 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=30.00 normalized_0_100=75.00 confidence_0_1=0.74
- Rationale: This simulation directly addresses honest partial conformance for real-time apps by separating provisional live-state transport from durable PromiseGrid audit publication. It remains incomplete for this scenario because it does not yet specify authoritative boundary identity, exact staged semantics claims, or how provisional signature/capability/witness language should be published.
- Strengths:
  - Supports honest early-slice claims by separating live transport from durable audit conformance.
  - Explicitly refuses to treat group-session or udp-feed as already solving reliable live CRDT transport.
  - Lets durable publication stay replayable and auditable while live transport evolves independently.
- Weaknesses:
  - Does not fully specify which identity or pCID is authoritative at the protocol boundary.
  - Does not yet turn app semantics into a staged conformance profile beyond the live-vs-audit split.
  - Leaves signature carriage and capability/witness language provisional and underspecified.
  - ... 1 more
- Risks:
  - Audit-only conformance may still be misread as full app interoperability.
  - Off-grid live transport can hide centralization or hidden-state assumptions.
  - Different layer-specific claims may drift and confuse auditors or implementers.
- Open questions:
  - Should current guidance allow only audit-layer conformance, or also a provisional future live protocol sketch?
  - When local IDs differ from published protocol claims, which identity or pCID is authoritative for Bob's interoperability check?
  - What exact audit object should claims cite: snapshot CID, CRDT save blob, op log, milestone receipt, or break-witness?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kugap-live-sync-audit-split` x `device-bound-agent-physical-effect`

- Result path: `results/SIM-kugap-live-sync-audit-split/device-bound-agent-physical-effect/openai-gpt-5.4-xhigh/20260521-011601.json`
- Scores: scenario_fit=3 promisegrid_alignment=3 auditability=3 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=4
- Fitness: raw=20.00 normalized_0_100=57.00 confidence_0_1=0.81
- Rationale: Moderate fit: the split cleanly matches host-local device control plus durable audit publication, and it clarifies conformance boundaries well. But the scenario's hardest requirements—delegated device authority, replay-safe at-most-once execution, and century-scale interpretation of physical-effect evidence—remain outside this simulation.
- Strengths:
  - Clean separation between PromiseGrid audit claims and host-driver/runtime dependencies such as CUPS, libusb, or IPP.
  - Plausible pattern for publishing durable receipts, milestones, or break-witnesses after host-local device activity.
  - Strong evolution posture because live/device protocols can change without collapsing the audit layer.
- Weaknesses:
  - Closer fit for sensor/reporting and audit publication than for non-idempotent actuation like printing.
  - Does not specify delegation proof, device identity binding, or owner/operator authority records.
  - No concrete replay-safe deduplication or at-most-once rule after daemon restart.
- Risks:
  - A replayed request could trigger duplicate physical effects if completion is not tied to a stable effect identifier.
  - Readers may overinterpret audit publication as a guarantee about device execution semantics the simulation does not define.
  - Long-term auditors may not be able to interpret evidence once the original hardware and driver stack are gone.
- Open questions:
  - Should a durable record bind command CID, delegation proof, device identity, and effect receipt or break-witness together?
  - When post-restart state is ambiguous, must the daemon default to break-witness instead of re-executing?
  - What conformance text cleanly separates PromiseGrid guarantees from host-local device and driver dependencies?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kugap-live-sync-audit-split` x `live-crdt-audit-publication`

- Result path: `results/SIM-kugap-live-sync-audit-split/live-crdt-audit-publication/openai-gpt-5.4-xhigh/20260521-011601.json`
- Scores: scenario_fit=5 promisegrid_alignment=4 auditability=3 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=2
- Fitness: raw=26.00 normalized_0_100=74.00 confidence_0_1=0.82
- Rationale: Strong scenario match: the simulation squarely addresses the need to separate sub-second live CRDT sync from durable PromiseGrid audit publication, but it remains provisional on the cited audit object, reconnection evidence, and live-binding path.
- Strengths:
  - Squarely targets the live-sync versus durable-audit split the scenario tests.
  - Preserves a durable 100-year audit publication story without overclaiming current transports as live sync.
  - Makes layer separation and conformance scoping explicit.
- Weaknesses:
  - No concrete reliable live binding or reconnection procedure is defined.
  - The exact object named by audit messages is unresolved.
  - Peer-local evidence and sparse-knowledge behavior are not yet spelled out.
- Risks:
  - Provisional wording could still be misread as approval of current transports for live CRDT sync.
  - Underspecified audit references could weaken replay, comparison, or break-witness publication.
  - The off-grid/on-grid bridge may create semantic drift between live state and durable promise accounting.
- Open questions:
  - Should the guide stop at off-grid live plus on-grid audit for now, or also sketch a future live pCID?
  - Should audit publication cite snapshots, save blobs, op logs, milestone receipts, or some combination?
  - What conformance language cleanly separates live-channel promises from audit-layer promises after partitions and reconnects?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kugap-live-sync-audit-split` x `multi-embodiment-app-identity`

- Result path: `results/SIM-kugap-live-sync-audit-split/multi-embodiment-app-identity/openai-gpt-5.4-xhigh/20260521-011601.json`
- Scores: scenario_fit=4 promisegrid_alignment=4 auditability=4 evolution_safety=3 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=2
- Fitness: raw=30.00 normalized_0_100=75.00 confidence_0_1=0.82
- Rationale: Strong partial fit. The simulation directly addresses the live-channel versus durable-audit split and the need for separate conformance claims, which is central to heterogeneous embodiments of one app. Its main gap is that it does not yet fully define the shared app-identity or contract-family mechanism that makes those embodiments unambiguously one app over time.
- Strengths:
  - Directly matches the scenario's custom live channel plus durable audit publication pressure.
  - Makes layer boundaries explicit: live transport and audit publication can have different pCIDs, timings, durability promises, and claims.
  - Supports later human/LLM audit better than a pure live-sync design by insisting on durable audit publication.
  - ... 1 more
- Weaknesses:
  - Does not yet specify the contract-family or identity anchor that proves the browser and plugin are the same app.
  - Per-embodiment conformance claim contents are still underspecified.
  - Recovery details for browser storage loss, host replacement, and helper turnover are only partially addressed.
- Risks:
  - Ambiguity between shared UX branding and shared protocol contract could still let Mallory misrepresent compatibility.
  - An off-grid live channel may hide central-service or trust assumptions unless explicitly constrained.
  - Live and audit layers may drift unless versioning and migration rules are made explicit.
- Open questions:
  - What shared contract or family identifier anchors 'one app' across heterogeneous embodiments?
  - What exact fields must each embodiment's conformance claim publish?
  - What durable records let Carol audit continuity after browser storage loss, helper replacement, or long-gap replay?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kugap-live-sync-audit-split` x `portable-signing-key-identity`

- Result path: `results/SIM-kugap-live-sync-audit-split/portable-signing-key-identity/openai-gpt-5.4-xhigh/20260521-011601.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=2 evolution_safety=2 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=3 risk_penalty=4
- Fitness: raw=12.00 normalized_0_100=43.00 confidence_0_1=0.88
- Rationale: This simulation is strong on separating live-session transport from durable audit publication, which is adjacent to the scenario's need to avoid identity conflation, but it does not actually specify a portable signing-key identity recipe, key-rotation evidence chain, or browser/helper storage guidance.
- Strengths:
  - Clear live-state versus durable-audit layer split.
  - Explicitly provisional posture avoids overstating current protocol claims.
  - Names concrete live/audit failure pressures such as dropped frames, stale snapshots, and partition recovery.
- Weaknesses:
  - Does not define a browser-plus-helper portable signing-key identity model.
  - Does not specify evidence linking old and new signing keys for Carol's audit.
  - Does not address display-name collisions, local usernames, browser key storage, or XSS risk.
- Risks:
  - Could be misread as answering identity continuity when it only addresses live-versus-audit layering.
  - Cross-host embodiment continuity may be left implicit, leading to inconsistent implementations.
  - Durable audit publication without explicit rotation semantics could create ambiguous long-term evidence.
- Open questions:
  - What durable claim or receipt should bind old and new signing keys across host embodiments?
  - How should browser-held and helper-held key storage be described without making false security claims?
  - What provisional handshake separates live-session presence from durable actor identity and presentation hints?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kugap-live-sync-audit-split` x `kernel-porting-boundary`

- Result path: `results/SIM-kugap-live-sync-audit-split/kernel-porting-boundary/openai-gpt-5.4-xhigh/20260521-011601.json`
- Scores: scenario_fit=1 promisegrid_alignment=3 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=2 risk_penalty=3
- Fitness: raw=23.00 normalized_0_100=58.00 confidence_0_1=0.80
- Rationale: Useful as a boundary-setting example, not as a primary answer to kernel porting. It clearly separates live transport from durable audit and avoids false conformance claims, but it leaves the minimum porting target, harness-vs-product boundary, and K-layer obligations largely unresolved.
- Strengths:
  - Explicitly separates low-latency live-state transport from durable PromiseGrid audit publication.
  - Preserves future frozen-spec conformance by refusing to overclaim group-session or udp-feed as live CRDT transport.
  - Keeps provisional claims and evolution pressure visible through separate pCIDs, timings, and conformance surfaces.
- Weaknesses:
  - Only tangentially addresses the kernel/runtime/dispatcher naming and first-port scope questioned by this scenario.
  - Does not define a minimum viable porting target or clearly separate harness apparatus from real port obligations.
  - Leaves K1-K5 ingress, feed, CAS, session, and app-layer blocking versus provisional scope mostly unmapped.
- Risks:
  - Porters may misread the off-grid live channel pattern as part of the core PromiseGrid port boundary.
  - Conformance claims can remain ambiguous until frozen terminology and obligations land through DR-davod.
  - App-specific live-sync concerns may distract from lower-layer portability work.
- Open questions:
  - Should the guide name the porting target kernel, runtime, dispatcher, handler host, or library?
  - What is the minimum viable first port before DR-davod closes?
  - Which K1-K5 ingress, feed, CAS, session, and app-layer details stay blocked versus provisional?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

## Required JSON Shape

{"child_id":"SIM-tizad-child-descriptive-design-slug","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
