# TE-pokul: CBOR array extensibility, pCID-owned proof, and likely grid-envelope landing shape

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-pokul

## Status

needs DF

## Decision under test

Given current scored evidence, CBOR-array encoding, and the rule that `pCID`
names the protocol spec that can define payload shape and proof semantics, what
outer envelope shape should PromiseGrid treat as the most likely landing zone
for continued design work?

## Assumptions

- The outer envelope is a CBOR array, so the array header already provides a
  variable-arity escape hatch for future extension.
- A `pCID` is a Protocol CID: it names a protocol spec document, not a payload
  object. The protocol named by that `pCID` can define payload shape, signable
  view, proof encoding, proof verification, and sender-authorship semantics.
- If a three-slot envelope uses `[pCID, payload, proof]`, and the protocol named
  by `pCID` defines `proof` as the current sender's evidence over the message,
  then the third slot is already the current-sender statement. There is no
  separate Alt-C in that case.
- Promise Theory remains the governing interpretation layer: envelope artifacts
  carry evidence that agents use to make, recognize, remember, and evaluate
  promises. They do not grant authority or create non-local trust.
- Trust is local, relationship-specific, and assessed by each agent from its
  own evidence, direct history, referrals, or other locally accepted signals.
- The relevant currently scored specimens are:
  - `SIM-dalor-grid-envelope-protocol-owned-signature-slot`
  - `SIM-dutam-grid-envelope-fixed-header-variable-body`
  - `SIM-natim-grid-envelope-nested-payload-outer-attestation-multisig`
  - `SIM-pobod-grid-envelope-outer-promise-nested-signed-payload`
  - the related signable-view and payload-owned-proof family under
    `TODO-tugoz`
- The question is not whether a given shape can be made to work at all. The
  question is which default outer shape deserves the most design pressure as
  the likely long-term PromiseGrid wire-envelope baseline.

## Alternatives

- **Alt-A - Minimal two-slot outer grid.** Keep the universal outer envelope as
  `[pCID, payload]`. Any required proof, signable view, or nested attestation is
  defined by the protocol contract named by `pCID` and carried inside the
  protocol-defined payload.
- **Alt-B - Three-slot outer grid with protocol-defined sender proof.** Keep
  the outer envelope as `[pCID, payload, proof]`. The protocol named by `pCID`
  defines what the third slot means, including proof encoding, canonical signed
  bytes, and whether the proof is the current sender's promise/evidence for
  carrying this exact message.

## Scenario analysis

### Scenario 1 - 100-year evolution without shape churn

Alice wants the default outer envelope to survive algorithm churn and future
proof families for a century without repeated wire-shape redesign. Because the
envelope is a CBOR array, and because `pCID` points at a spec doc that can say
"the proof is varsig-encoded" or later "the proof uses some future format,"
future-proofing does not require extra visible outer fields by itself.

Alt-A handles this by pushing proof evolution into the protocol-defined payload.
Alt-B handles this by keeping one stable outer proof slot while letting the
protocol named by `pCID` define the proof family. Neither needs a separate
future-proofing mechanism outside CBOR plus the protocol spec.

### Scenario 2 - Generic relay, quarantine, and partial understanding

Bob is a relay or archival peer. He may not understand the inner payload spec,
but he still wants enough outer information to decide whether he can store,
forward, quarantine, or reject a message while keeping an evidence trail.

Alt-A is smallest, but a generic peer sees only a `pCID` plus opaque payload
bytes unless the payload convention is widely understood. Alt-B gives generic
peers a stable place to find proof bytes, even if the protocol still defines
the proof format. That can improve preservation and quarantine behavior without
turning the proof into global authority.

### Scenario 3 - Small-device and low-complexity implementation

Carol is implementing on a small device and wants minimal parser complexity and
long-lived positional stability.

Alt-A is the smallest parser surface: one CBOR array and two slots. Alt-B adds
one slot, but keeps the shape stable and still lets the `pCID`-named protocol
define proof details. The cost of Alt-B is not future-proofing complexity; it is
the requirement that every conforming message carry and preserve the third slot.

### Scenario 4 - Promise Theory and current-sender accountability

Dave wants the envelope to help local trust assessment without drifting into
authorization or command semantics. He especially cares about distinguishing
"the current sender promises X about this message" from "the nested payload
contains evidence that some earlier agent promised Y."

Alt-A can represent both, but only if the protocol-defined payload makes that
distinction explicit. Alt-B can represent it at the outer level if the protocol
named by `pCID` defines the third slot as the current sender's proof over the
message. That is the intended Alt-B semantics; it is not a separate alternative.

### Scenario 5 - Multi-legal-entity trust and selective cooperation

Alice, Bob, and Carol are separate legal entities. Alice decides whether to
send Bob data for storage or computation based on Bob's promise history and
current promises. Carol later evaluates Bob's returned results based on her own
local trust model.

No alternative creates trust. Alt-A and Alt-B both only provide evidence that
agents may use locally. Alt-B can make the current sender's promise easier to
preserve and inspect, but Alice, Bob, and Carol still decide locally whether
that evidence matters.

### Scenario 6 - Current scored evidence from the sim corpus

The recent PT-gated slices matter because they show which ideas survive contact
with the current scenarios. The strongest current specimens are not the most
minimal possible shapes; they are the ones that keep the envelope small while
making the sender's promise or evidence easy to identify.

The `dalor` three-slot specimen remains useful because it is the compact form of
Alt-B. Its weakness in recent scoring is not that Alt-B lacks the current-sender
semantics; those semantics are the intended meaning of Alt-B. Its weakness is
that some source prose confused `pCID` with a payload identifier and did not
state the Promise Theory meaning plainly enough.

## Conclusions

- Reject "visible outer extensibility for future-proofing" as a load-bearing
  argument. CBOR array encoding plus `pCID`-named protocol specs already provide
  the future-proofing escape hatch.
- Reject the earlier split between Alt-B and Alt-C. If `[pCID, payload, proof]`
  is specified so the third slot is the current sender's proof over this
  message, then Alt-B already has the semantics that were mistakenly described
  as Alt-C.
- Alt-A survives as the simplest plausible default: minimal outer grid,
  protocol-owned proof semantics, and no universal outer proof slot unless the
  protocol named by a specific `pCID` carries one inside its payload.
- Alt-B survives as the main competing default: one extra outer slot that every
  message preserves, with the protocol named by `pCID` defining the proof
  encoding and sender-evidence semantics.

## Recommendation

Treat the likely landing zone as a choice between Alt-A and Alt-B:

- Alt-A: `[pCID, payload]`
- Alt-B: `[pCID, payload, proof]`, where the protocol named by `pCID` defines
  proof encoding and the sender-proof meaning

The design question is not "do we need extra fields for future-proofing?" We do
not. The remaining question is whether the default envelope should always give
the current sender one stable outer proof slot, or whether that proof belongs
inside each protocol-defined payload when needed.

## Implications for open work

- Keep `SIM-dalor-grid-envelope-protocol-owned-signature-slot` breed-eligible,
  but supersede its confusing `pCID` wording before relying on it as evidence
  for Alt-B.
- Continue breeding and scoring corrected Alt-B successors against PT-clean
  parents such as `SIM-vuliv`, `SIM-konit`, `SIM-fonom`, `SIM-natim`, and
  `SIM-pobod`.
- Treat `SIM-dutam-grid-envelope-fixed-header-variable-body` as evidence for
  stable substrate shape, not yet evidence for a complete PromiseGrid
  interaction design.
- Focus follow-on sims and scoring pressure on the exact Alt-A vs Alt-B
  question: minimal two-slot outer grid versus three-slot outer grid with a
  protocol-defined sender proof.

## Decision status

`needs DF` - surviving alternatives are Alt-A and Alt-B. The earlier Alt-C was
not a separate alternative; it was the intended Promise Theory meaning of Alt-B.
The next DF question is: should the default PromiseGrid outer envelope remain
the minimal CBOR array `[pCID, payload]`, or should it use `[pCID, payload,
proof]` where the protocol named by `pCID` defines the proof encoding and the
current-sender proof semantics?
