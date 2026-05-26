# GA Child Generation

Return only JSON with keys `child_id`, `design_delta_summary`, and `files`.
Choose a descriptive `child_id` that starts with `SIM-vinag-child-` and ends with a kebab-case design slug. Do not use generic `ga-child`, `pending`, or ordinal-only names.
Each file path must be relative to the child simulation root. Include `README.md` and `QUESTION.md`.

Optimization goal: breed a child simulation from exactly two parent simulations, expected to score higher than its parent set on the same rubric and sampled scenarios.
Use the fitness evidence below as training feedback: preserve parent strengths, repair weaknesses, reduce risks, answer or route open questions, and keep changes to one to three bounded design deltas.
Do not merely summarize the parent. The child must make an explicit design move that should improve `fitness.normalized_0_100` while keeping the simulation standalone and auditable.

## Design Guardrails

- Treat `pCID` as Protocol CID: the pCID-named protocol spec may define payload shape, signature/proof encoding, refusal evidence, freeze-successor records, or capability promise-token records.
- For base-envelope children, avoid selector-shopping stacks such as `env_pCID`/`sig_pCID`/`payload_pCID`, generic `claim_header` or claim-card layers, generic claim cards, and universal `statement_capsule` wrappers.
- Do not ban higher-layer payload protocols from defining their own promise-accounting records. Put signed refusal versus silence/timeout, exact-byte local evidence, freeze successor records, transfer semantics, and capability-as-promise-token behavior inside the pCID-selected payload/specimen layer unless the scenario explicitly asks to test an envelope negative control.
- Keep outer signatures scoped to the current sender's own promise; no agent promises for another agent, and every receiver keeps local trust assessments.

- Run group ID: `ga-canary-20260525-goban-dalor-zukis-v4`
- Planned child ID prefix: `SIM-vinag-child-`
- Temporary child ID: `SIM-vinag-child-pending-0001`
- Temporary child path: `proposals/ga-canary-20260525-goban-dalor-zukis-v4/simulations/SIM-vinag-child-pending-0001/`
- Operation: `breed`
- Parent IDs: `SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig, SIM-dalor-grid-envelope-protocol-owned-signature-slot`

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

### `simulations/SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig/README.md`

```markdown
# SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig

This simulation is a direct specimen of the locked `TE-fikoj` /
`DI-sisak` outer-envelope direction. It does not reopen the fixed-three-slot
versus variable-arity decision. Instead, it makes that decision concrete with
one protocol-owned example:

```text
grid([42(pCID), payload, varsig])
```

For this specimen:

- slot `0` is always the tagged protocol selector `42(pCID)`;
- slot `1` is the primary payload anchor;
- slot `2` is this protocol's own `varsig` proof slot.

The broader family rule remains larger than this one specimen: PromiseGrid's
current direction is `grid([42(pCID), payload, ...])`, and the protocol named
by `pCID` defines whether later outer slots exist and what they mean. `SIM-zukis`
therefore tests one direct, PT-clean member of that family without turning
slot `2 = varsig` into universal envelope law. Source: `DI-sisak`; `DI-mabit`.

## Promise-Theory framing

The outer envelope helps a receiver interpret another agent's promise. It does
not command behavior, grant global permission, or decide trust centrally. In
this specimen, the current sender's `varsig` is evidence for the sender's own
scoped promise:

> "I promise these payload bytes and this outer-slot arrangement meet the
> protocol specification named by this `pCID`."

Each receiver still decides locally whether it recognizes the protocol, trusts
the sender, verifies the `varsig`, stores the bytes, relays them, or uses the
payload. Carriage is not semantic acceptance. Source: `DI-pagin`; `DI-sisak`;
`DI-mabit`.

## What this sim is testing

This sim tests whether a tagged selector in slot `0`, a stable payload anchor
in slot `1`, and one protocol-owned proof slot in slot `2` give PromiseGrid a
good balance of:

- DAG-CBOR / CID ecosystem interop;
- small deterministic outer parsing;
- protocol-owned evolution of later outer slots;
- clean separation between base-envelope promises and higher-layer promise
  accounting.

## Comparison targets

Primary comparison targets:

- `SIM-dalor-grid-envelope-protocol-owned-signature-slot`
- `SIM-pobod-grid-envelope-outer-promise-nested-signed-payload`
- `SIM-jufag-grid-envelope-quarantine-sig-pcid-outcomes`

`SIM-dalor` is the nearest fixed-three-slot neighbor. `SIM-pobod` pressures
explicit outer promise wording and nested signed payload structure. `SIM-jufag`
is the contrasting explicit-`sig_pcid` selector-shopping branch. Source:
`DI-mabit`.

## Boundaries

This sim does not declare that every PromiseGrid protocol must use slot `2`
as `varsig`. It only tests whether one direct specimen of the locked
`grid([42(pCID), payload, ...])` family performs well when the protocol named
by `pCID` chooses that shape for itself.
```

### `simulations/SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig/QUESTION.md`

```markdown
# Question

Can PromiseGrid use a CBOR outer envelope
`grid([42(pCID), payload, varsig])` as one direct specimen inside the broader
locked family `grid([42(pCID), payload, ...])`, where slot `0` is the tagged
protocol selector, slot `1` is the primary payload anchor, and slot `2` is a
protocol-owned `varsig` proof rather than a universal envelope law?

Source: `DI-sisak`; `DI-mabit`.
```

### `simulations/SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig/protocols/grid-envelope.d/specs/grid-envelope-draft.md`

```markdown
# Grid-envelope draft: tag-42 selector with protocol-owned slot-2 varsig

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Variant: `tag42-protocol-owned-slot2-varsig`.

## Scope

This spec defines one direct grid-envelope specimen for wire-lab comparison. It
is a specimen inside
`SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig`, not a harness rule
and not the canonical PromiseGrid envelope. It uses `pCID` only as Protocol
CID: the content identifier of the protocol specification document, never the
content identifier of payload bytes. Source: `DI-sisak`; `DI-mabit`.

This specimen implements one concrete member of the broader locked family:

```text
grid([42(pCID), payload, ...])
```

For this specimen, the protocol named by `pCID` chooses one later outer slot:

```text
grid([42(pCID), payload, varsig])
```

That is a protocol-owned choice for this specimen, not a universal requirement
that every PromiseGrid protocol use a third-slot `varsig`. Source: `DI-sisak`;
`DI-mabit`.

## Envelope Shape

The outer envelope shape is:

```text
[42(pCID), payload, varsig]
```

Slots are interpreted positionally:

- slot `0` is the tagged protocol selector, currently `42(pCID)`;
- slot `1` is opaque payload bytes until interpreted by the protocol named by
  `pCID`;
- slot `2` is this protocol's `varsig` proof over the signable view named by
  the same `pCID`.

The key design move under test is:

- PromiseGrid fixes the selector position and the primary payload anchor;
- the protocol named by `pCID` owns whether later outer slots exist;
- this specimen uses that freedom to place one `varsig` proof in slot `2`
  without introducing a second selector such as `sig_pcid`.

## Signable Bytes

The signable view for this specimen is the canonical bytes of:

```text
[42(pCID), payload]
```

The `varsig` in slot `2` is evidence over that exact prefix unless the protocol
named by `pCID` later refines associated-data rules more narrowly. This binds
both the tagged selector and the payload bytes without adding outer
selector-shopping machinery. Source: `DI-mabit`.

## Encoding

The outer envelope is a deterministic CBOR positional array. Slot `0` is the
tagged selector `42(pCID)`. Slots `1` and `2` are byte strings at the carrier
layer. The CBOR array header carries arity; this specimen does not add a second
arity field.

Small receivers do not need a full IPLD object model. To recover the selector
they need only:

- CBOR parsing;
- tag `42`;
- the following byte string;
- the leading `00` sentinel;
- CID parsing.

## Unknown pCID Policy

If a receiver lacks a handler for `pCID`, it may preserve or blind-carry the
exact outer bytes as uninterpreted evidence under local policy, but it MUST NOT
claim to parse the payload or verify the `varsig`.

This keeps the Promise Theory boundary explicit: bytes may survive as evidence,
but semantic acceptance remains local and protocol-dependent. Carriage is not
acceptance. Source: `DI-sisak`; `DI-mabit`.

## `varsig` Policy

This specimen has no separate `sig_pcid`, `env_pcid`, or `payload_pcid`. The
single `pCID` defines:

- what `varsig` encoding is valid in slot `2`;
- what signer binding and signer identity rules apply;
- whether freshness, delegation, threshold, or revocation semantics exist;
- whether any associated data beyond canonical `[42(pCID), payload]` bytes is
  required.

The universal envelope itself enforces only three things:

- slot `0` is the tagged selector;
- slot `1` is the primary payload anchor;
- later outer-slot roles are owned by the protocol named by `pCID`.

## Comparison Pressure

Compared with `SIM-dalor`, this specimen keeps a visible outer proof slot but
also makes the tagged selector `42(pCID)` part of the direct specimen.

Compared with `SIM-pobod`, this specimen keeps the outer shape smaller and
avoids pushing explicit nested structure into the base-envelope design.

Compared with `SIM-jufag`, this specimen removes `sig_pcid` and keeps one-pCID
discipline: one protocol selector names payload shape and slot-2 proof
semantics together. Source: `DI-mabit`.

## Open Questions

- Does one protocol-owned `varsig` slot preserve enough generic audit clarity
  without reintroducing selector shopping?
- Does this specimen outperform the fixed-three-slot `dalor` branch on the same
  scenario slice while remaining simpler than explicit-`sig_pcid` designs?
- Is slot `2 = varsig` a strong direct specimen of the broader
  `grid([42(pCID), payload, ...])` family, or does it still freeze too much
  proof structure too early?

## Non-Canonical Status

This draft does not declare a winning universal slot-2 rule. It exists to test
one direct member of the locked tagged-selector family against nearby
three-slot, nested-payload, and explicit-selector alternatives.
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

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/kernel-porting-boundary/openai-gpt-5.4-medium/20260525-152501.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=3 evolution_safety=3 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 envelope_discipline=3 kernel_implementation_promises=1 app_protocol_promise_semantics=1 risk_penalty=2
- Fitness: raw=41.00 normalized_0_100=63.08 confidence_0_1=0.86
- Rationale: Strong Promise-Theory framing and a clear envelope-layer boundary make this a credible wire specimen, but the scenario asks for kernel porting promises and this design intentionally does not specify those. It helps the porting boundary by keeping unsupported-pCID behavior and exact-byte evidence explicit, yet leaves most kernel/app obligations to higher-layer protocols and future guide decisions.
- Strengths:
  - Very clear local boundary between outer-envelope evidence and higher-layer promise accounting.
  - Good Promise-Theory vocabulary: sender intent, local trust, no command or global authority framing.
  - Explicit unknown-pCID behavior preserves exact bytes without pretending to parse or verify them.
  - ... 1 more
- Weaknesses:
  - Weak direct fit for a kernel-porting scenario because it does not enumerate kernel services or implementation promises.
  - No explicit app/kernel operation catalog, refusal record shape, adapter mapping, or host-assumption separation.
  - Audit clarity for generic peers is limited because proof-family semantics are hidden behind the pCID-specific protocol.
- Risks:
  - A first port team could mistake this envelope specimen for the kernel porting target and under-specify storage, dispatch, lifecycle, and evidence promises.
  - Using one pCID for both payload and proof semantics may make proof-family evolution or generic audit tooling harder.
  - The array shape is close to current envelope direction but not fully aligned with the more explicit grid([42(pCID), payload, ...]) direction.
- Open questions:
  - Is coupling payload shape and proof-family rules under one pCID evolution-safe enough for long-lived ports?
  - Should future kernel-facing guidance require a more explicit outer slot form closer to grid([42(pCID), payload, ...]) for audit clarity?
  - How should a kernel port advertise supported/unsupported pCIDs and exact-byte carriage when this envelope is only one layer of the stack?
- Authority boundary: Evidence only; evaluates an envelope-layer specimen and does not settle PromiseGrid kernel or porting-target design.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `promise-accounting-records-kept-storage-promise`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/promise-accounting-records-kept-storage-promise/openai-gpt-5.4-medium/20260525-152501.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=4 evolution_safety=3 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=5 promise_vocabulary=4 simplicity_durability=4 envelope_discipline=3 kernel_implementation_promises=1 app_protocol_promise_semantics=1 risk_penalty=2
- Fitness: raw=42.00 normalized_0_100=64.62 confidence_0_1=0.88
- Rationale: Strong Promise-Theory framing for an envelope-layer specimen: it clearly limits the outer promise to sender evidence about payload shape under a named protocol, keeps trust local, and avoids authority framing. For this storage-promise scenario, however, the design intentionally delegates promise accounting records, kept-storage semantics, and local trust updates to the payload protocol, so scenario fit is only partial rather than poor.
- Strengths:
  - Clear layer boundary: outer envelope handles protocol naming and proof carriage while higher-layer promise accounting stays in payload.
  - Promise-first wording is mostly correct and explicitly local.
  - Deterministic positional CBOR plus canonical signing of [pCID,payload] gives good byte-level audit evidence.
  - ... 1 more
- Weaknesses:
  - It does not model the scenario's kept-storage promise records or future pull/keep/advertise decisions at the app layer.
  - The universal mandatory third slot is a stronger envelope commitment than the current minimal-envelope direction prefers.
  - Proof evolution may require new protocol pCIDs for what some designs would treat as proof-profile changes.
- Risks:
  - A mandatory outer proof slot may be more envelope surface than necessary and could age less well than a smaller baseline.
  - Proof-family semantics hidden behind each protocol pCID may fragment generic tooling and reduce quick audit legibility.
  - The simulation does not itself define observable kept-storage outcome records, so scenario success depends on a higher-layer payload protocol.
- Open questions:
  - Is one protocol pCID enough for proof-family evolution and cross-protocol audit clarity over long horizons?
  - For kept-storage promise accounting, should all observable outcome records live entirely in the payload protocol while this envelope remains unchanged?
  - Does a mandatory outer proof slot justify its extra surface versus a smaller envelope direction for cases that do not need universal outer signatures?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `promise-accounting-records-refused-service`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/promise-accounting-records-refused-service/openai-gpt-5.4-medium/20260525-152501.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=3 evolution_safety=3 layer_boundary_clarity=5 failure_handling=2 implementation_plausibility=4 promise_vocabulary=4 simplicity_durability=4 envelope_discipline=3 kernel_implementation_promises=1 app_protocol_promise_semantics=1 risk_penalty=2
- Fitness: raw=39.00 normalized_0_100=60.00 confidence_0_1=0.84
- Rationale: Strong promise-first envelope specimen with clear local-trust and sender-evidence framing. It is weak on this scenario only because refused-service accounting is intentionally outside the envelope layer; the design can carry such records but does not define them. The boundary is explicit enough that this is a scenario-fit limitation, not a PT failure.
- Strengths:
  - Very clear layer boundary: outer envelope authenticates protocol-tagged payload bytes and leaves promise accounting to the payload protocol.
  - Good Promise-Theory wording: sender signature is evidence of the sender's own scoped promise, not authority or command.
  - Explicit unknown-pCID behavior supports sparse knowledge and local evidence retention.
  - ... 1 more
- Weaknesses:
  - Does not itself distinguish honest refusal from failure, timeout, or corruption.
  - Kernel/app promise semantics are largely delegated, so this scenario is only partially addressed at the claimed layer.
  - Proof-family evolution may require minting new protocol pCIDs more often than designs with a separate proof selector.
- Risks:
  - Scenario pressure may be under-served if reviewers mistake authenticated carriage for refusal semantics.
  - A mandatory universal proof slot adds cross-protocol overhead even though proof rules remain protocol-owned.
  - Audit clarity may suffer for generic peers because proof-family meaning is hidden behind the protocol pCID rather than a separate outer selector.
- Open questions:
  - Can a payload protocol carried in this envelope express signed refusal records distinctly from timeout, corruption, and ordinary failure without adding outer-envelope complexity?
  - Does coupling proof-family changes to the protocol pCID create too much pCID churn over time?
  - Will generic auditors have enough clarity when the outer proof slot is present but its exact verification semantics are owned by the payload protocol?
- Authority boundary: Evidence only; evaluates an envelope-layer specimen and does not settle higher-layer promise-accounting design.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `cas-object-model-dag-cbor-interop`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/cas-object-model-dag-cbor-interop/openai-gpt-5.4-medium/20260525-152501.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=4 evolution_safety=3 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 envelope_discipline=3 kernel_implementation_promises=1 app_protocol_promise_semantics=1 risk_penalty=2
- Fitness: raw=42.00 normalized_0_100=64.62 confidence_0_1=0.84
- Rationale: A strong promise-first envelope specimen: small, explicit, deterministic, and clear about local trust and unknown-pCID behavior. It partially helps the DAG-CBOR interop scenario through CBOR encoding and CID-aware slots, but it does not itself define DAG-CBOR payload/link/tag semantics, so scenario coverage is limited at this layer. The main rubric drag is the mandatory universal proof slot and the coupling of proof-family evolution to the protocol pCID.
- Strengths:
  - Very clear envelope-layer boundary.
  - Good Promise-Theory wording centered on sender-scoped intent.
  - Deterministic small CBOR array with direct signing of canonical [pCID,payload] bytes.
  - ... 1 more
- Weaknesses:
  - Does not itself answer the scenario's object-model questions about DAG-CBOR tags and link representation.
  - Mandatory universal proof slot is less disciplined than the leanest current envelope direction.
  - Limited kernel/app promise semantics because those are delegated out of scope.
- Risks:
  - Interop remains ambiguous unless carrier and payload protocols pin down DAG-CBOR-compatible CID/link/tag handling.
  - A mandatory outer proof slot may overconstrain simple content-addressed object transport.
  - Coupling payload shape and proof-family rules to one pCID may increase protocol churn during evolution.
- Open questions:
  - Can the carrier profile encode CID links and CBOR tags in a way that stays legible to DAG-CBOR/IPLD tooling while preserving the strict three-slot envelope?
  - Does making the outer proof slot mandatory create unnecessary friction for unsigned archival or pure-CAS payload cases?
  - Will proof-family evolution force too many new protocol pCIDs because proof rules are owned by the same pCID as payload shape?
- Authority boundary: Scores only the envelope-layer specimen against DAG-CBOR interop pressure; payload object semantics, storage semantics, and higher-layer trust accounting remain outside this simulation.

### `SIM-dalor-grid-envelope-protocol-owned-signature-slot` x `l6-cas-starting-profile-bakeoff-long-horizon-reprofile`

- Result path: `results/SIM-dalor-grid-envelope-protocol-owned-signature-slot/l6-cas-starting-profile-bakeoff-long-horizon-reprofile/openai-gpt-5.4-medium/20260525-152501.json`
- Scores: scenario_fit=2 promisegrid_alignment=4 auditability=4 evolution_safety=3 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=4 envelope_discipline=4 kernel_implementation_promises=2 app_protocol_promise_semantics=1 risk_penalty=2
- Fitness: raw=44.00 normalized_0_100=67.69 confidence_0_1=0.84
- Rationale: Strong promise-first envelope specimen with clear local trust boundaries, good auditability, and a compact deterministic artifact. It fits the envelope layer well but only indirectly addresses the scenario's long-horizon CAS reprofile pressure because migration semantics for old pointer objects and richer successor profiles are delegated to payload protocols. Main tradeoff: proof semantics are cleanly protocol-owned, but tying them to the same pCID as payload shape may increase reprofile churn and reduce generic audit legibility over time.
- Strengths:
  - Very clear layer boundary: outer envelope only names protocol, carries opaque payload, and carries proof bytes.
  - Excellent promise vocabulary for the envelope layer; the sender promise is scoped and local.
  - Unknown-pCID policy is locally trustworthy and preserves uninterpreted evidence without false claims.
  - ... 1 more
- Weaknesses:
  - Scenario fit is partial because the design does not itself explain how old pointer objects and raw chunks stay explainable across profile upgrades.
  - Kernel implementation promises are mostly absent.
  - App/payload migration semantics are intentionally delegated, leaving long-horizon reprofile behavior underspecified at this layer.
- Risks:
  - Payload-shape and proof-family evolution are coupled to one protocol pCID, which may force avoidable reminting during long-horizon reprofile.
  - The mandatory universal outer signature slot may harden one evidence pattern earlier than necessary.
  - Generic peers can preserve bytes but may have limited audit clarity when proof semantics are hidden behind unknown pCID-specific rules.
- Open questions:
  - Does coupling proof-family evolution to the payload protocol pCID create avoidable migration churn over long horizons?
  - Should the outer envelope align more directly with the current grid([42(pCID), payload, ...]) direction rather than a plain three-slot array reading?
  - Is a universal mandatory third proof slot worth the extra outer commitment for payloads that may want different evidence patterns?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig` x `kernel-porting-boundary`

- Result path: `results/SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig/kernel-porting-boundary/openai-gpt-5.4-medium/20260525-152501.json`
- Scores: scenario_fit=2 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=4 promise_vocabulary=5 simplicity_durability=5 envelope_discipline=5 kernel_implementation_promises=1 app_protocol_promise_semantics=2 risk_penalty=2
- Fitness: raw=48.00 normalized_0_100=73.85 confidence_0_1=0.85
- Rationale: PT-clean envelope specimen with strong tagged-selector discipline, explicit local-trust framing, and clear delegation of higher-layer semantics to the pCID-owned protocol. It fits the scenario only partially because the scenario asks for kernel porting promises, host assumptions, refusals, and service decomposition that this envelope sim does not attempt to define.
- Strengths:
  - Excellent envelope discipline: grid([42(pCID), payload, varsig]) with protocol-owned later-slot meaning.
  - Very clear layer boundary between envelope parsing and higher-layer/kernel semantics.
  - Strong Promise-Theory wording: scoped sender promise, local trust, no authority framing.
  - ... 2 more
- Weaknesses:
  - Does not provide the kernel implementation promise record the scenario explicitly asks for.
  - Does not map storage, compute, send/receive, key, lifecycle, dispatch, or evidence operations to named promisers.
  - Host/runtime assumptions for native, browser, mobile, MCU, or split-service ports are not specified.
  - ... 1 more
- Risks:
  - Porters may mistake the wire-envelope specimen for the full kernel porting target.
  - The visible slot-2 varsig could be overgeneralized into a universal rule despite the docs warning against that.
  - Lack of explicit kernel promise records may leave unsupported pCIDs, namespace behavior, and evidence movement underspecified during first-port planning.
- Open questions:
  - What minimum kernel promise record should sit above this envelope so ports can name supported pCIDs, refusals, and host assumptions explicitly?
  - Can local APIs wrap pCID-selected grid messages without hiding the promiser boundary or exact-byte evidence needs?
  - Does protocol-owned slot 2 remain flexible enough for future kernel evidence/refusal patterns without being mistaken for a universal envelope rule?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig` x `promise-accounting-records-kept-storage-promise`

- Result path: `results/SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig/promise-accounting-records-kept-storage-promise/openai-gpt-5.4-medium/20260525-152501.json`
- Scores: scenario_fit=2 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=5 envelope_discipline=5 kernel_implementation_promises=2 app_protocol_promise_semantics=1 risk_penalty=1
- Fitness: raw=50.00 normalized_0_100=76.92 confidence_0_1=0.87
- Rationale: Strong PT-clean envelope specimen with explicit tagged-selector discipline, small deterministic structure, and clear local unknown-pCID behavior. For this scenario, however, it mainly supplies transport/audit evidence and intentionally delegates storage-promise accounting semantics to the payload protocol, so scenario fit is only partial rather than poor.
- Strengths:
  - Excellent layer-boundary clarity: outer envelope versus higher-layer promise accounting is explicit.
  - Promise-first wording is strong and PT-correct at the claimed layer.
  - Envelope shape is compact, deterministic, and durable for small receivers.
  - ... 2 more
- Weaknesses:
  - Does not itself model a kept storage promise or Bob's local future pull/keep/advertise updates.
  - Provides little direct app-protocol guidance for storage outcomes beyond carrying signed bytes.
  - Kernel/local implementation promises are only lightly addressed compared with a fuller kernel-focused design.
- Risks:
  - Scenario pressure may be undersatisfied if evaluators expect direct storage-promise outcome records at this layer.
  - A protocol-owned varsig slot could still tempt over-reading as a universal proof pattern if the specimen boundary is forgotten.
  - Without a paired higher-layer record format, successful serve-after-store evidence remains underspecified for cross-peer comparison.
- Open questions:
  - Should a higher-layer pCID define a standard kept-storage observation record for successful later serve events?
  - Is one protocol-owned varsig slot enough generic evidence when storage-promise scenarios need richer make/break outcome records?
  - What minimal payload-level record shape best links exact served bytes to a prior storage promise without bloating the envelope?
- Authority boundary: Evidence only; this envelope specimen clarifies outer-wire promises and local trust boundaries, but storage-promise accounting and kept/served outcome records belong to higher-layer payload protocols and peer-local records.

### `SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig` x `promise-accounting-records-refused-service`

- Result path: `results/SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig/promise-accounting-records-refused-service/openai-gpt-5.4-medium/20260525-152501.json`
- Scores: scenario_fit=2 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=3 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=5 envelope_discipline=5 kernel_implementation_promises=2 app_protocol_promise_semantics=1 risk_penalty=1
- Fitness: raw=50.00 normalized_0_100=76.92 confidence_0_1=0.89
- Rationale: Strong PT-clean envelope specimen with excellent boundary clarity, vocabulary, and envelope discipline. For this scenario, however, it intentionally stops at transport/evidence framing and leaves refused-service accounting semantics to the payload protocol, so fit to the specific refusal-record question is only partial rather than poor.
- Strengths:
  - Excellent alignment with grid([42(pCID), payload, ...]) and protocol-owned later-slot roles.
  - Very clear PT framing: sender promises exact bytes/shape under a named protocol; receivers assess locally.
  - Small deterministic CBOR artifact with durable parsing and explicit unknown-pCID behavior.
  - ... 1 more
- Weaknesses:
  - Does not itself model refused service as distinct from failure, corruption, or timeout.
  - Kernel/app promise semantics are intentionally thin in this specimen.
  - Scenario success depends on a higher-layer payload protocol to define refusal records and trust-accounting consequences.
- Risks:
  - Scenario-specific refusal semantics may be underspecified unless paired with a higher-layer protocol for signed refusal records and local trust updates.
  - A visible slot-2 varsig could be overgeneralized by later readers into quasi-universal envelope law despite the draft's explicit boundary notes.
  - Audit trails can prove exact observed bytes, but not by themselves classify refusal versus timeout or corruption.
- Open questions:
  - Should a higher-layer pCID define a standard signed refusal record that binds the refused service context to the observed envelope bytes?
  - Is slot-2 varsig sufficient for later app protocols that need durable make/break/refusal evidence without adding selector-shopping?
  - What minimal payload-side pattern best distinguishes honest refusal from timeout, corruption, or non-observation while preserving local trust accounting?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig` x `cas-object-model-dag-cbor-interop`

- Result path: `results/SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig/cas-object-model-dag-cbor-interop/openai-gpt-5.4-medium/20260525-152501.json`
- Scores: scenario_fit=4 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=4 envelope_discipline=5 kernel_implementation_promises=2 app_protocol_promise_semantics=2 risk_penalty=1
- Fitness: raw=53.00 normalized_0_100=81.54 confidence_0_1=0.85
- Rationale: Strong PT-clean envelope specimen. It matches the current grid([42(pCID), payload, ...]) direction, keeps one-pCID discipline, and gives explicit unknown-pCID behavior. For the DAG-CBOR interop scenario it performs well at the envelope/carrier layer, but it intentionally delegates richer CAS object-model semantics to the payload protocol, so scenario fit is high rather than perfect.
- Strengths:
  - Excellent layer-boundary clarity between envelope parsing and payload semantics.
  - Good DAG-CBOR/CID ecosystem fit via deterministic CBOR array, tag 42, and CID-oriented selector handling.
  - Explicit local behavior for unknown pCID preserves auditability without pretending semantic acceptance.
  - ... 1 more
- Weaknesses:
  - Does not itself define a CAS object model, so scenario pressure about default object format remains only partially answered.
  - Kernel-facing implementation promises are minimal because the sim stays at envelope scope.
  - Fixed use of slot 2 for this specimen may still be more structure than the minimal family rule strictly requires.
- Risks:
  - The specimen may still overfit around a visible slot-2 varsig choice before broader envelope experience accumulates.
  - Interop success depends on payload protocols defining DAG-CBOR/CID link semantics clearly; the envelope alone does not settle TE-43.
  - Tooling may conflate protocol CID selection with payload/object-format semantics unless documentation stays sharp.
- Open questions:
  - Does binding the signable view to [42(pCID), payload] preserve enough audit clarity across heterogeneous DAG-CBOR/IPLD tooling?
  - Is slot 2 = varsig a durable specimen, or does it freeze proof structure earlier than necessary for long-term interop evolution?
  - How much DAG-CBOR object-model interoperability must live in payload protocols rather than the envelope to satisfy TE-43 cleanly?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

### `SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig` x `l6-cas-starting-profile-bakeoff-long-horizon-reprofile`

- Result path: `results/SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig/l6-cas-starting-profile-bakeoff-long-horizon-reprofile/openai-gpt-5.4-medium/20260525-152501.json`
- Scores: scenario_fit=3 promisegrid_alignment=5 auditability=4 evolution_safety=4 layer_boundary_clarity=5 failure_handling=4 implementation_plausibility=5 promise_vocabulary=5 simplicity_durability=5 envelope_discipline=5 kernel_implementation_promises=2 app_protocol_promise_semantics=1 risk_penalty=2
- Fitness: raw=51.00 normalized_0_100=78.46 confidence_0_1=0.87
- Rationale: Strong PT-clean envelope specimen that matches the current tag-42, protocol-owned-slot direction and stays small, explicit, and auditable. It fits the long-horizon reprofile scenario only partially because it deliberately delegates object-graph migration and old-pointer explainability to the payload protocol rather than modeling them itself.
- Strengths:
  - Excellent layer-boundary clarity between envelope promises and higher-layer payload semantics.
  - High envelope-discipline alignment with grid([42(pCID), payload, ...]) and protocol-owned later slots.
  - Good auditability via exact signable bytes over [42(pCID), payload] and explicit unknown-pCID behavior.
  - ... 1 more
- Weaknesses:
  - Does not itself specify how old and new profile objects remain mutually explainable after reprofiling.
  - Kernel implementation promises are mostly absent beyond envelope handling behavior.
  - App-level promise semantics for storage, transfer, capability behavior, and trust updates are intentionally out of scope.
- Risks:
  - Scenario pressure around old pointer objects, raw chunks, and richer replacement profiles is only indirectly addressed at the envelope layer.
  - A visible slot-2 varsig specimen may still freeze proof structure earlier than necessary if later protocols want different outer evidence layouts.
  - The draft is not frozen and the pCID is not yet minted, so long-horizon stability is not yet demonstrated.
- Open questions:
  - Does fixing this specimen to slot 2 = varsig create future migration debt if later profiles want different outer-slot evidence shapes?
  - Can old payloads and pointer objects remain explainable across reprofiling purely through stable pCID-scoped payload protocols, or is extra higher-layer evidence needed?
  - How should long-horizon receivers compare and trust old versus new profile pCIDs when both remain addressable?
- Authority boundary: Evidence only; does not settle PromiseGrid design.

## Required JSON Shape

{"child_id":"SIM-vinag-child-descriptive-design-slug","design_delta_summary":"one to three bounded design deltas","files":[{"path":"README.md","content":"# ..."},{"path":"QUESTION.md","content":"# ..."}]}
