# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Choose a descriptive `child_id` that starts with `SIM-virim-child-` and ends with a kebab-case design slug. Do not use generic `ga-child`, `pending`, or ordinal-only names.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

- Run group ID: `ga-canary-20260521-210902`
- Planned child ID prefix: `SIM-virim-child-`
- Temporary child ID: `SIM-virim-child-pending-0001`
- Temporary child path: `proposals/ga-canary-20260521-210902/simulations/SIM-virim-child-pending-0001/`
- Operation: `breed`
- Parent IDs: `SIM-nijuz-multi-embodiment-identity, SIM-kugap-live-sync-audit-split`

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

### `simulations/SIM-nijuz-multi-embodiment-identity/README.md`

```markdown
# SIM-nijuz: Multi-embodiment identity

This simulation is a provisional question home for `FB-zazon` and `FB-robif`:
one logical app spanning heterogeneous host embodiments while preserving one
protocol contract and one signing-key identity story. Source: `DI-ragaz`.

## Question

What worked example can show a browser tab and a long-running plugin/helper as
one logical PromiseGrid app without overstating unresolved identity, storage, or
binding decisions? Source: `DI-ragaz`.

## Decision Axes

- **One app, multiple embodiments:** each component claims only the part of the
  shared pCID-selected contract it actually implements.
- **Portable identity:** a user needs continuity across browser and plugin hosts
  without treating display names or local usernames as identity.
- **Key recipe status:** algorithm, rotation, handshake, and constrained-host
  storage guidance may be provisional and pivotable.
- **Host constraints:** browser tabs, Node helpers, plugins, and native daemons
  have different filesystem, process, and key-storage promises.
- **Auditability:** later peers must understand which component made which claim
  under which key and protocol version.

## Related Root Scenarios

- `scenarios/multi-embodiment-app-identity/multi-embodiment-app-identity.md`
- `scenarios/portable-signing-key-identity/portable-signing-key-identity.md`

## Boundaries

This simulation does not bless a permanent cryptographic primitive or storage
mechanism. It tests provisional guide language for conformance and identity
continuity across constrained hosts. Source: `DI-ragaz`.
```

### `simulations/SIM-nijuz-multi-embodiment-identity/QUESTION.md`

```markdown
# Question

How can one logical PromiseGrid app span a browser tab and a long-running plugin
host while publishing honest per-component conformance claims and a portable
signing-key identity story? Source: `DI-ragaz`.

Open decision points:

- What does each embodiment claim when one has only browser storage and another
  has a durable helper process?
- What provisional signing-key recipe, rotation lifecycle, and handshake shape is
  safe to teach before `DR-tuhaz` closes?
- How should the guide distinguish durable signing identity from display names,
  colors, and local user handles?
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

### `SIM-kugap-live-sync-audit-split` x `minimal-immutable-blob-app`

- Result path: `results/SIM-kugap-live-sync-audit-split/minimal-immutable-blob-app/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=1 promisegrid_alignment=4 auditability=3 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=3 risk_penalty=3
- Fitness: raw=19.00 normalized_0_100=54.00 confidence_0_1=0.83
- Rationale: Useful mainly as a boundary-discipline pattern: it cleanly separates live transport claims from durable audit publication, which helps preserve 100-year evidence and echoes the scenario's need to separate content identity from availability. But it is not a direct minimal blob design, so the core `hash in -> blob out`, authorization, retention, and cross-site retrieval promises remain underdefined.
- Strengths:
  - Explicitly separates protocol layers and avoids overclaiming that one mechanism solves another.
  - Keeps durable audit publication in scope, supporting later human and LLM inspection.
  - Provisional framing is evolution-safe: live and audit layers can change independently.
- Weaknesses:
  - Weak direct fit for the minimal immutable blob app; it targets real-time sync pressure.
  - Does not define who promises long-term blob availability, discovery, or read authorization.
  - Failure analysis centers on dropped frames and stale audit snapshots more than years-later blob retrieval loss.
- Risks:
  - Audit publication could be mistaken for a storage, availability, or authorization guarantee.
  - Applying this sim here could hide the need to specify whether a hash is only an address or also a capability.
  - If retention changes and caches are incomplete, Carol still lacks a clear retrieval path and local accounting recipe.
- Open questions:
  - For a pure blob app, what durable evidence should Bob publish: blob CID, storage receipt, retention promise, or all three?
  - Who owns the long-term availability and discovery promise after the original host changes policy?
  - What local promise-accounting records should Alice, Bob, and Carol keep when retrieval fails years later?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kugap-live-sync-audit-split` x `app-semantics-partial-conformance`

- Result path: `results/SIM-kugap-live-sync-audit-split/app-semantics-partial-conformance/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=4 promisegrid_alignment=4 auditability=4 evolution_safety=5 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=3 risk_penalty=2
- Fitness: raw=26.00 normalized_0_100=74.00 confidence_0_1=0.78
- Rationale: Good fit on the scenario's honesty pressure: the simulation explicitly separates provisional live transport from durable PromiseGrid audit publication and warns against overclaiming current live-sync support. It is weaker on concrete boundary identity, exact partial-conformance wording, and provisional signature/capability/witness details.
- Strengths:
  - Directly supports an honest B-side claim by splitting off-grid live sync from on-grid audit publication.
  - Makes layer boundaries explicit with distinct pCIDs, timings, durability promises, and conformance claims.
  - Calls out dropped live frames, stale audit snapshots, and partition recovery as first-class pressure.
- Weaknesses:
  - Does not yet define authoritative protocol-boundary identity clearly enough for third-party interop.
  - Leaves exact audit citation objects and claim wording unresolved.
  - Capability-token, break-witness, and signature-carriage language remain largely open.
- Risks:
  - Teams may still market the audit layer as full real-time PromiseGrid conformance.
  - Different live adapters could fragment interoperability behind a superficially similar audit story.
- Open questions:
  - What exact wording should a partial-conformance claim use for live-sync apps?
  - Which artifact should audit publication cite: snapshot CID, save blob, op log, milestone receipt, or break-witness?
  - How should local IDs and adapter-specific signature carriage map to the authoritative protocol boundary?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kugap-live-sync-audit-split` x `device-bound-agent-physical-effect`

- Result path: `results/SIM-kugap-live-sync-audit-split/device-bound-agent-physical-effect/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=1 promisegrid_alignment=3 auditability=3 evolution_safety=3 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=2 risk_penalty=4
- Fitness: raw=15.00 normalized_0_100=43.00 confidence_0_1=0.83
- Rationale: Strong as a boundary-pattern for separating live/device-facing action from durable audit publication, but weak as a direct answer to this scenario because delegation evidence, at-most-once physical-effect semantics, replay dedup, and break-witness specifics are mostly unspecified.
- Strengths:
  - Very clear live-vs-audit layer separation.
  - Explicitly avoids overclaiming current PromiseGrid audit layers as low-latency control transport.
  - Supports durable publication of receipts or related audit artifacts.
- Weaknesses:
  - Does not specify how owner delegation to a device-bound daemon is proven.
  - Does not define dedup or at-most-once handling for non-idempotent effects after restart/replay.
  - Does not say what receipt proves an actual physical effect versus a request or attempt.
- Risks:
  - A reboot plus replay could cause duplicate physical actions or ambiguous effect records.
  - Audit evidence may outlive the device/driver stack without preserving effect meaning.
  - Conformance claims could blur local CUPS/libusb/IIO-style dependencies with PromiseGrid guarantees.
- Open questions:
  - Should device commands stay entirely off-grid with only receipts or break-witnesses published on-grid?
  - What stable effect ID or receipt schema deduplicates a physical action across restart and delayed propagation?
  - How is operator/owner delegation attested so Carol can audit it after the host stack is gone?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kugap-live-sync-audit-split` x `live-crdt-audit-publication`

- Result path: `results/SIM-kugap-live-sync-audit-split/live-crdt-audit-publication/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=5 promisegrid_alignment=4 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=2
- Fitness: raw=32.00 normalized_0_100=80.00 confidence_0_1=0.82
- Rationale: Strong match for the scenario: it cleanly separates low-latency live sync from durable PromiseGrid audit publication and avoids overclaiming current transport support, but it remains provisional on recovery details and the exact audit object.
- Strengths:
  - Directly targets the live-sync versus durable-audit split the scenario is testing.
  - Very clear layer boundary: group-session can be audit publication without pretending to be the live CRDT transport.
  - Calls out partition, dropped-frame, stale-audit, and conformance-scope issues explicitly.
- Weaknesses:
  - Leaves the reliable live binding intentionally undefined.
  - Does not yet fix the exact object an audit message should cite.
  - Provides little detail on peer-local accounting and sparse-knowledge mechanics.
- Risks:
  - Readers may still conflate provisional live guidance with actual PromiseGrid conformance.
  - Different apps may publish incompatible audit evidence until the cited-object shape is pinned down.
  - An off-grid live channel could become a sticky dependency without a clearer migration path.
- Open questions:
  - Should the guide standardize only off-grid live plus on-grid audit for now, or also sketch a future live pCID?
  - Should audit messages cite snapshots, CRDT save blobs, op-log checkpoints, receipts, or break-witnesses?
  - How should conformance language separate live-channel claims from audit-layer claims?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kugap-live-sync-audit-split` x `multi-embodiment-app-identity`

- Result path: `results/SIM-kugap-live-sync-audit-split/multi-embodiment-app-identity/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=4 promisegrid_alignment=4 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=3 risk_penalty=3
- Fitness: raw=24.00 normalized_0_100=69.00 confidence_0_1=0.80
- Rationale: Strong partial fit: the simulation cleanly separates embodiment-specific live transport from durable audit publication and supports honest per-component conformance claims, but it leaves the app-level contract family and recovery linkage underdefined, so 'one app' identity remains somewhat ambiguous.
- Strengths:
  - Separates live-channel claims from durable audit claims, reducing false full-stack conformance.
  - Supports scoped per-embodiment claims for browser, plugin, and helper components with different runtime limits.
  - Provides a durable audit path that can outlast embodiment churn and support long-term review.
- Weaknesses:
  - Does not yet define the contract family or manifest that binds heterogeneous embodiments into one app.
  - Peer-local accounting and compatibility rules across the custom live channel remain thin.
  - Minimum conformance-claim contents for each embodiment are not yet concrete.
- Risks:
  - Mallory can exploit loose 'same app' labeling if cross-embodiment linkage is not explicit.
  - Carol may mistake audit publication for proof of live transport guarantees such as ordering or reliability.
  - Upgrades or host replacement could create silent drift between embodiments without durable compatibility evidence.
- Open questions:
  - What app-level contract or descriptor links browser, plugin, and helper as one logical app?
  - What exact fields must each embodiment claim: runtime limits, supported subprotocols, durability scope, and excluded guarantees?
  - Which durable objects should be cited for recovery and long-term verification after storage loss or helper replacement?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kugap-live-sync-audit-split` x `portable-signing-key-identity`

- Result path: `results/SIM-kugap-live-sync-audit-split/portable-signing-key-identity/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=1 promisegrid_alignment=2 auditability=2 evolution_safety=3 layer_boundary_clarity=4 failure_handling=2 implementation_plausibility=2 risk_penalty=2
- Fitness: raw=14.00 normalized_0_100=40.00 confidence_0_1=0.85
- Rationale: Good as a boundary-setting partial answer: it clearly separates ephemeral live-session behavior from durable audit publication, which helps avoid treating a new live session as proof of durable identity continuity. But it does not define the core portable signing-key identity recipe across browser and helper, so it only weakly covers this scenario.
- Strengths:
  - Clear separation between live transport claims and durable audit claims.
  - Explicitly provisional posture supports future migration instead of freezing the wrong protocol shape.
  - Durable audit publication could later carry rotation receipts or break-witness evidence for audit.
- Weaknesses:
  - Portable cross-host signing-key identity is outside the simulation's main question.
  - No concrete evidence format links old and new signing keys.
  - Does not explain durable identity versus display names/local usernames, or browser storage/XSS caveats.
- Risks:
  - Readers could overread the live/audit split as an identity solution when it is only an architectural boundary.
  - Without a separate identity recipe, session re-entry from another host embodiment may be mistaken for durable key continuity.
- Open questions:
  - What durable PromiseGrid object should bind old and new signing keys across browser and helper?
  - How should live-session authentication reference durable identity without implying the live channel is PromiseGrid-native?
  - What minimum guide-safe statement about WebCrypto/IndexedDB, filesystem storage, and XSS risk belongs in v0?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-kugap-live-sync-audit-split` x `kernel-porting-boundary`

- Result path: `results/SIM-kugap-live-sync-audit-split/kernel-porting-boundary/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=1 promisegrid_alignment=4 auditability=3 evolution_safety=4 layer_boundary_clarity=3 failure_handling=2 implementation_plausibility=2 risk_penalty=3
- Fitness: raw=16.00 normalized_0_100=46.00 confidence_0_1=0.82
- Rationale: Useful as a guardrail on honest provisional conformance claims, especially separating off-grid live transport from durable PromiseGrid audit publication, but a weak fit for kernel-porting pressure because it does not define the first-port target, harness boundary, or runtime/library obligations.
- Strengths:
  - Clearly separates live-state transport claims from durable audit publication claims.
  - Explicitly avoids implying that current group-session or udp-feed already solves reliable live CRDT transport.
  - Preserves room for future frozen-spec conformance by treating the live layer as provisional.
- Weaknesses:
  - Does not answer the kernel vs runtime vs dispatcher vs library naming question.
  - Does not separate wire-lab harness apparatus from the real porting target.
  - Leaves K1-K5, CAS, session, and app-layer porting obligations too abstract for a first port plan.
- Risks:
  - Readers may misread audit publication support as end-to-end live sync support.
  - A first porter could still treat group-session or harness behavior as normative runtime surface.
  - Conformance claims may drift until DR-davod freezes the relevant obligations.
- Open questions:
  - What is the minimum viable porting target if live sync remains explicitly off-grid?
  - Which draft evidence can a first porter cite now for audit publication, and which claims must wait for frozen pCIDs?
  - Which terms should the guide bless for first-port claims: kernel, runtime, dispatcher, handler host, or library?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-nijuz-multi-embodiment-identity` x `minimal-immutable-blob-app`

- Result path: `results/SIM-nijuz-multi-embodiment-identity/minimal-immutable-blob-app/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=1 promisegrid_alignment=3 auditability=3 evolution_safety=4 layer_boundary_clarity=3 failure_handling=1 implementation_plausibility=2 risk_penalty=3
- Fitness: raw=14.00 normalized_0_100=40.00 confidence_0_1=0.79
- Rationale: Strong on honest multi-embodiment identity claims and migration-friendly key guidance, but it does not really answer the minimal blob app question about who promises hash-to-blob retrieval, retention, availability, or authorization years later.
- Strengths:
  - Explicitly requires each embodiment to claim only the contract surface it actually implements.
  - Separates durable signing identity from display names and local usernames.
  - Keeps key and storage guidance provisional, which helps future migration.
- Weaknesses:
  - Does not define the blob contract or the meaning of hash in and blob out.
  - Does not separate content identity from availability, authorization, replication, and retention promises.
  - Does not specify local promise-accounting records for long-delayed retrieval failure.
- Risks:
  - Could be misread as solving CAS guarantees when it mostly solves cross-host identity continuity.
  - Browser and helper storage differences may hide durability assumptions unless explicit retention promises are added.
  - A shared signing identity could be overextended into an authorization model for blob reads.
- Open questions:
  - Which embodiment, if any, promises long-term retention or replication for a returned hash?
  - Is possession of the hash only an address, or also a read capability in this design?
  - What evidence survives key rotation, host replacement, or policy change so Carol can audit a failed read years later?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-nijuz-multi-embodiment-identity` x `app-semantics-partial-conformance`

- Result path: `results/SIM-nijuz-multi-embodiment-identity/app-semantics-partial-conformance/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=5 promisegrid_alignment=4 auditability=4 evolution_safety=5 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 risk_penalty=2
- Fitness: raw=32.00 normalized_0_100=80.00 confidence_0_1=0.82
- Rationale: Strong fit: the simulation is explicitly about honest per-embodiment conformance claims and portable signing-key identity across heterogeneous hosts. It is strongest on boundary clarity and evolution safety, but remains a provisional question home, so unsupported semantics and local/adversarial verification behavior are still under-defined.
- Strengths:
  - Explicitly says each embodiment should claim only the part of the shared contract it actually implements.
  - Separates protocol-boundary signing identity from display names and local user handles.
  - Keeps key recipe, storage, and handshake guidance provisional and pivotable.
  - ... 1 more
- Weaknesses:
  - Still framed as a question home, not a full worked partial-conformance profile.
  - Does not yet spell out clear wording or schema for supported versus blocked semantics.
  - Little concrete treatment of Bob/Carol local evidence, Mallory misuse, or stale/incorrect claims.
  - ... 1 more
- Risks:
  - A 'one logical app' story could blur per-embodiment limits and encourage overclaiming.
  - Provisional signature-carriage or key-storage guidance could be mistaken for a final standard.
  - Interop may fail if unsupported semantics are not explicitly marked as blocked or absent.
- Open questions:
  - What exact B-side wording makes a partial-conformance claim honest without implying full app semantics?
  - What should each embodiment claim when one has only browser storage and another has durable helper state?
  - What wire-visible artifact binds browser and helper embodiments to one authoritative signing identity across rotation?
  - ... 1 more
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-nijuz-multi-embodiment-identity` x `device-bound-agent-physical-effect`

- Result path: `results/SIM-nijuz-multi-embodiment-identity/device-bound-agent-physical-effect/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=2 promisegrid_alignment=3 auditability=3 evolution_safety=4 layer_boundary_clarity=4 failure_handling=1 implementation_plausibility=3 risk_penalty=3
- Fitness: raw=22.00 normalized_0_100=55.00 confidence_0_1=0.78
- Rationale: Partial fit: the simulation is strong on multi-embodiment identity, per-component conformance claims, and long-term interpretability, but it does not yet define delegation proof, physical-effect receipts, or replay/restart handling for at-most-once behavior.
- Strengths:
  - Clear per-embodiment boundary: each component claims only the contract slice it actually implements.
  - Portable signing-key identity across browser and helper embodiments supports continuity across hosts without relying on local names.
  - The simulation explicitly keeps key, rotation, handshake, and storage choices provisional, which helps avoid prematurely freezing unsafe binding assumptions.
- Weaknesses:
  - No explicit model for non-idempotent physical effects or at-most-once execution.
  - No receipt, break-witness, or peer-local replay-dedup story across daemon restart.
  - Device delegation evidence and driver or SDK dependency claims are not specified.
- Risks:
  - If applied as-is to device control, replay after restart could trigger duplicate physical actions.
  - Provisional binding and storage choices could leave long-lived delegation evidence ambiguous.
- Open questions:
  - What locally auditable artifact proves Bob delegated the printer or sensor to a specific daemon embodiment?
  - What durable operation ID and receipt rules let the daemon report "already done" versus emit a break-witness after restart?
  - How should per-embodiment conformance claims encode dependencies on CUPS, libusb, IIO, IPP, or vendor SDK layers?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-nijuz-multi-embodiment-identity` x `live-crdt-audit-publication`

- Result path: `results/SIM-nijuz-multi-embodiment-identity/live-crdt-audit-publication/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=3 promisegrid_alignment=4 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=3 risk_penalty=2
- Fitness: raw=23.00 normalized_0_100=66.00 confidence_0_1=0.74
- Rationale: Strong on multi-embodiment identity and honest per-component conformance, which helps separate live CRDT transport from durable audit publication, but it does not yet define the live binding, save-time cited object, or partition recovery flow.
- Strengths:
  - Directly fits a browser-plus-helper/plugin embodiment split close to the scenario's browser and editor setup.
  - Explicit per-component conformance claims support clean separation between live-channel promises and audit-layer promises.
  - Provisional key and binding guidance is evolution-friendly and avoids premature lock-in.
- Weaknesses:
  - Does not specify the reliable ordered low-latency transport needed for live CRDT sync.
  - Does not state exactly what content-addressed snapshot or merge artifact the audit message should cite.
  - Partition, reordering, and delayed publication handling remain underspecified.
- Risks:
  - One logical app identity could obscure which embodiment actually made which claim or state transition.
  - A provisional live binding sketch could be mistaken for a settled PromiseGrid transport commitment.
  - Thin audit publication may be insufficient for contested long-term review of save-time state.
- Open questions:
  - Should live CRDT traffic remain explicitly off-grid until a reliable binding exists?
  - What exact CAS snapshot or lineage should the save-time audit message reference?
  - How should key rotation and post-partition reconciliation work across browser and helper embodiments?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-nijuz-multi-embodiment-identity` x `multi-embodiment-app-identity`

- Result path: `results/SIM-nijuz-multi-embodiment-identity/multi-embodiment-app-identity/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=5 promisegrid_alignment=4 auditability=4 evolution_safety=5 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=3 risk_penalty=3
- Fitness: raw=31.00 normalized_0_100=77.50 confidence_0_1=0.79
- Rationale: Direct hit on the scenario: it defines one logical app through a shared pCID-selected contract while keeping embodiment claims honest and separately scoped. It is strongest on boundary clarity and evolution safety, but it remains more of a question frame than a complete mechanism for key binding, recovery, and audit evidence.
- Strengths:
  - Makes browser and plugin 'one app' by shared contract rather than UX branding or host identity.
  - Requires each embodiment to claim only the contract subset it actually implements.
  - Separates durable signing identity from display names and local usernames.
  - ... 1 more
- Weaknesses:
  - Does not yet specify a concrete attestation or binding scheme across browser and helper embodiments.
  - Failure and recovery behavior for browser storage loss, host replacement, and helper upgrades is underdefined.
  - Peer-local promise accounting and sparse-knowledge evidence are implied more than shown.
- Risks:
  - Ambiguous cross-embodiment identity could let Mallory overread 'same app' as stronger shared authority than warranted.
  - Portable signing identity may encourage unsafe key export, duplication, or confusing rekey flows on constrained hosts.
  - Long-term audits may become ambiguous if embodiment changes are not recorded distinctly from identity continuity.
- Open questions:
  - What exact fields must each embodiment's conformance claim include?
  - How are browser and helper keys linked to one logical app identity without a central registry?
  - What local records survive browser storage loss or helper replacement well enough for later audit?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-nijuz-multi-embodiment-identity` x `portable-signing-key-identity`

- Result path: `results/SIM-nijuz-multi-embodiment-identity/portable-signing-key-identity/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=5 promisegrid_alignment=4 auditability=4 evolution_safety=5 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 risk_penalty=3
- Fitness: raw=32.00 normalized_0_100=80.00 confidence_0_1=0.82
- Rationale: Direct fit to the scenario: the simulation is explicitly about one logical app spanning browser and helper/plugin embodiments while preserving durable signing-key identity, honest per-embodiment claims, auditability, and a pivotable v0 key recipe. The score is capped by intentionally unresolved details around cross-host binding, concrete rotation evidence, and browser-storage/XSS guidance.
- Strengths:
  - Targets the exact browser-plus-helper identity continuity problem.
  - Makes embodiment boundaries explicit by limiting each component to honest partial conformance claims.
  - Strong evolution posture: key algorithm, rotation, handshake, and storage guidance are explicitly provisional and pivotable.
  - ... 1 more
- Weaknesses:
  - Does not yet specify the concrete evidence chain linking old and new signing keys.
  - Cross-embodiment binding remains underspecified beyond the high-level identity story.
  - Failure handling for XSS, stolen browser state, and confusing local usernames is acknowledged more than fully designed.
- Risks:
  - A provisional recipe could be misread as a permanent identity standard before DR-tuhaz closes.
  - Browser key storage guidance could overstate security if XSS and host-compromise limits are not stated sharply.
  - Users or implementers may still conflate display names or local handles with durable signing identity without stricter guide language.
- Open questions:
  - What exact signed artifacts prove rotation continuity from old key to new key for later auditors?
  - How should the guide describe browser-side key storage and XSS risk without implying stronger guarantees than exist?
  - What handshake or attestation shape best links browser and helper embodiments as one logical app without relying on central authority?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-nijuz-multi-embodiment-identity` x `kernel-porting-boundary`

- Result path: `results/SIM-nijuz-multi-embodiment-identity/kernel-porting-boundary/openai-gpt-5.4-xhigh/20260521-210902.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=4 evolution_safety=4 layer_boundary_clarity=2 failure_handling=2 implementation_plausibility=4 risk_penalty=2
- Fitness: raw=20.00 normalized_0_100=57.10 confidence_0_1=0.64
- Rationale: Useful adjacent evidence, not a full answer: SIM-nijuz helps with honest per-embodiment claims and pre-freeze portability guidance, but it leaves the core scenario question - the kernel/runtime/dispatcher porting surface and harness-vs-port separation - mostly unstated.
- Strengths:
  - It explicitly wants later peers to know which component made which claim under which key and protocol version.
  - Each embodiment claims only the part of the shared contract it actually implements.
  - Durable signing identity is separated from display names and local host handles.
  - ... 1 more
- Weaknesses:
  - The docs do not explicitly separate the wire-lab harness from the real porting target.
  - They do not map ingress, feed, CAS, session, and runtime obligations into blocked versus provisional buckets.
  - The focus is multi-embodiment app identity more than first-port kernel/runtime terminology.
- Risks:
  - Readers may overread app-level multi-embodiment guidance as settling the kernel/runtime boundary.
  - Provisional identity and storage guidance could harden into de facto obligations before the relevant design reviews close.
  - Browser/helper continuity may hide host-specific persistence, trust, and session-boundary gaps.
- Open questions:
  - What is the minimum viable porting target that can be claimed now without freezing the wrong surface?
  - Which ingress, feed, CAS, session, and app-layer details should stay blocked versus be taught as provisional orientation?
  - Which term - kernel, runtime, dispatcher, handler host, or library - should name the conformance boundary?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

## Required JSON Shape

{"child_id":"SIM-virim-child-descriptive-design-slug","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
