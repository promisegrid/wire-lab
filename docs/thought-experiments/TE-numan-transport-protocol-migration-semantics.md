# TE-numan: Transport-protocol migration invariants

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-numan

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TE-37` (integer alias)
- `TE-20260506-041241` (timestamp alias and pre-migration filename)

## Status

decided

## Decision under test

When a group of participants writes messages under one transport-protocol-pCID and later needs to operate under a different transport-protocol-pCID, the act of switching between them is a **migration**. This TE does not lock the migration's *operational shape* — close-old + open-new vs. overlap window vs. atomic-swap, the back-reference format, the disposition of pre-migration messages, the seal mechanics. Those are speculative without a concrete migration to design against and are explicitly deferred to a follow-on TE that will be written when (or if) a real migration becomes imminent.

Instead, this TE locks the **invariants** that any future migration-contract TE must satisfy: audit-trail reconstructibility, no silent rewrite of pre-migration history, and no unilateral abandonment of the old transport. These invariants are direct consequences of TE-zalut § S5 (frozen pCID contract is stable for the lifetime of the transport instance), TE-dajot's 100-year goal (the corpus must be reconstructible from on-disk artifacts and git history), and the audit-substrate framing locked across TE-nibar / TE-lozip / TE-dabol. They are load-bearing now, before any migration has happened, because the cost of retrofitting them onto migrations that already happened under no contract is much higher than the cost of locking them while the corpus is small.

The motivating fact is mostly housekeeping. TE-junil § S7 raised the migration question and deferred it as "anticipated TE-dajot." TE-dajot was then occupied by the 100-year-goal TE, leaving TE-junil's forward-pointer stale. `protocols/wire-lab.d/specs/transport-spec-draft.md` § OQ-2 raises the same question against an unnamed future TE. This TE is the new anchor for both pointers; a Cat-3 Refinement on TE-junil in the same twig dis-anchors the stale pointer. The transport-spec OQ-2 is updated to point here.

The migration question is fundamentally a PromiseGrid-level question (what happens when a group of participants must change the contract their transport is bound to), not a wire-lab artifact. The wire-lab is the simulation harness in which migrations would be observed; the on-disk encoding of a migration (directory keying, close-message format, back-reference layout) is wire-lab-specific. The semantic invariants this TE locks are PromiseGrid-level and would apply under any production substrate.

## Assumptions

- TE-zalut's four locked principles are in force: messages do not declare their transport; transport directories are keyed `transports/<pcid>--<slug>/`; each transport-protocol-pCID names its own spec; the code that reads the directory is the handler for that pCID.
- TE-zalut § S5 is in force: a frozen transport-protocol-pCID's contract is stable for the lifetime of the transport instance. Migration mints a new pCID; it does not edit the old contract. This TE accepts that and locks the invariants that follow.
- TE-junil's per-axis meta-rule (cardinality is a parameter except at extremes; visibility and routing topology warrant distinct pCIDs; etc.) is in force. The trigger conditions that legitimately motivate a migration are predominantly the axes that warrant distinct pCIDs. Parameter-level changes within the same pCID are not migrations and are out of scope.
- TE-dajot (the 100-year goal as a load-bearing design constraint) is in force: a reader who joins the corpus a year (or a century) later, with only on-disk contents and git history, must be able to reconstruct the design state. Migration must not break that property.
- TE-vipir (protocols-as-simulated-repos) is in force: each transport instance is a simrepo whose interior is governed by the spec named by its pCID. Migration is the act of switching which simrepo a group is writing under, not editing one in place.
- The audit-substrate framing locked across TE-nibar (spec-doc as promise), TE-lozip (congruence / convergence duality), and TE-dabol (TE editing policy) is in force. Pre-migration messages are part of the audit record; migration must respect that.
- Cooperative actors follow the alphabetical convention (Alice, Bob, Carol, Dave, Ellen, Frank). Mallory is the adversary. Steve is named explicitly only where his repo-owner role is load-bearing.
- Promise-theory vocabulary applies to the migration act itself: agreements among participants to perform a migration are conditional-promises and reciprocal-promises in the PT sense. The exact promise shape is downstream (operational-shape TE) and not locked here.
- Network and disk are reliable in this TE. Migration semantics under partition or data loss are out of scope.

### Threat / trust model

- Cooperative actors are honest and competent but may have stale information about whether a migration has happened. A participant who pulls before a migration completes may write into the old transport in good faith.
- Mallory has commit access (or can convince a maintainer to merge). She aims to use migration as a vector to (a) silently rewrite the past by editing pre-migration messages, (b) suppress pre-migration messages by abandoning the old transport without leaving any artifact pointing at it, or (c) split a group by creating a parallel transport that omits or edits the migration history.
- A reader joining the corpus after the migration has only the post-migration transport's contents, the pre-migration transport's contents (if preserved on disk), and the git history. The invariants this TE locks must let that reader determine that a migration happened and what its shape was, even if they cannot reconstruct the operational details (which the operational-shape TE will add later).
- The invariants must hold against Mallory in the absence of any operational-shape rules. They are the floor; the operational-shape TE adds the rest.

## Decision Forks (DFs)

This TE walks three invariant-DFs. Each is presented one at a time in conversation with the repo owner, with paragraphs describing the alternatives and a bot recommendation, before the multiple-choice question is asked.

- **DF-37.A — Audit-trail reconstructibility invariant.** Must any future migration contract guarantee that a reader with on-disk contents + git history can determine that a migration happened, when, and what the old / new transport-pCIDs are? *(LOCKED — Alt-A.1: hard requirement.)*
- **DF-37.B — No-silent-rewrite invariant.** Must any future migration contract forbid editing pre-migration message bytes? *(LOCKED — Alt-B.1: hard requirement; pre-migration messages are immutable.)*
- **DF-37.C — No-unilateral-abandonment invariant.** Must any future migration contract require an artifact (on disk, in the new transport, in the old transport, or somewhere reachable from both) that preserves the old transport's existence and frontier-CID, even when a participant wishes to forget the old transport? *(LOCKED — Alt-C.2: hard requirement with documented-erasure carve-out; the old transport's existence, pCID, and frontier-CID are preserved forever, but the old transport's contents may be intentionally erased under a recorded erasure-receipt.)*

Each invariant DF has a "yes / no / yes-but-narrowed" alternative shape rather than a multi-option operational menu. The DFs are ordered by independence: DF-37.A is the floor; DF-37.B is independent of A; DF-37.C is the most contentious and depends partly on the answers to A and B.

## Scenario analysis

These scenarios stress-test the three locked invariants under cooperative and adversarial play. They are deliberately invariant-focused; the operational shape of each scenario's migration is left abstract because that is downstream territory.

### S1 — Reader-from-the-future

Alice, Bob, and Carol completed a migration from `transports/<pcid-X>--alice-bob-carol/` to `transports/<pcid-Y>--alice-bob-carol-dave/` a year ago. Dave joined the new transport at migration time and has been an active participant since. Today Ellen clones the wire-lab repo for the first time and wants to understand the conversation history.

Under DF-37.A (hard reconstructibility), the corpus on disk plus git history must let Ellen determine that a migration happened from pCID-X to pCID-Y at the migration time, with no out-of-band knowledge. The operational-shape TE will lock *how* this is encoded (a closing-message in the old transport, an opening-message in the new transport, a back-reference, an entry in `transports/MIGRATIONS.md`, or some combination); this TE only requires that the encoding suffice for Ellen's reconstruction.

Under DF-37.B (immutability), the messages in `transports/<pcid-X>--alice-bob-carol/` that Ellen reads are byte-identical to the messages Alice / Bob / Carol read at migration time. No reframing has silently shifted their meaning.

Under DF-37.C (with documented-erasure carve-out), the old transport's directory either contains its messages (the default) or an erasure-receipt naming the pCID and frontier-CID-at-migration. Ellen can in either case determine the old transport's pCID and the migration's frontier; in the erasure case she sees the receipt rather than the messages.

The three invariants together give Ellen the reconstruction property at the cost of three on-disk requirements that any future operational-shape TE must satisfy.

### S2 — Mallory's silent rewrite

Mallory has commit access. She convinces herself she is justified in flipping the meaning of a pre-migration message that argued against a position she now wants to advance. She edits the bytes of that pre-migration message in `transports/<pcid-X>--alice-bob-carol/` and pushes the rewrite to the integration branch.

DF-37.B forbids this directly. The bytes-on-disk of pre-migration messages are immutable. Mallory's rewrite is a discrepancy detectable by `git log -p`, by content-hash check (the message's CID no longer matches its hash), and by spec-check tooling that verifies pre-migration message immutability. Cooperative actors reviewing Mallory's commit catch the rewrite at code-review time. If the rewrite slips past review, a reader's content-hash check still catches it later.

Mallory could try to launder the rewrite as a redaction. DF-37.B Alt-B.1 closes that route: the corpus rule is redact-by-annotation, not redact-by-edit. Annotations are new messages; they leave the original immutable.

### S3 — Mallory's silent abandonment

Mallory is the maintainer of the wire-lab clone that integrates with `ppx/main`. After a legitimate migration from pCID-X to pCID-Y, she runs `git rm -rf transports/<pcid-X>--alice-bob-carol/` and pushes. Post-migration history in `transports/<pcid-Y>--...` happens to lack any back-reference to pCID-X (because the operational-shape TE has not yet locked one, or because Mallory excluded it).

DF-37.C Alt-C.2 forbids this. The contract requires that some artifact preserves the old transport's pCID and frontier-CID-at-migration, persistent forever. Even if Mallory removes the old transport's directory, the contract requires the artifact to live in the new transport (a back-reference message), in a separate place (a MIGRATIONS.md entry), or in some other on-disk location reachable from a reader's clone of the post-migration corpus. Mallory's commit removing both the old transport directory and the back-reference is a discrepancy under DF-37.C; cooperative reviewers reject it.

If Mallory tries to relabel her abandonment as a documented erasure, the carve-out requires an erasure-receipt naming the old pCID and frontier-CID. A receipt that exists is not silent abandonment; it is a recorded erasure, which the audit substrate accepts. The contract design challenge is to ensure that erasure-receipts cannot be forged (e.g., that the authorizing participants are genuinely the participants of the old transport, not Mallory acting alone). That is operational-shape territory; this TE only requires that *some* receipt exists.

### S4 — Cooperative migration with stale pull

Alice, Bob, and Carol migrated from pCID-X to pCID-Y yesterday. Bob has been off-grid; he pulls the wire-lab repo today, but his clone is stale relative to integration: he sees only the pre-migration state of pCID-X. He writes a new message into `transports/<pcid-X>--alice-bob-carol/` in good faith, intending it as a follow-up to the most recent message he can see.

This scenario is not a violation of the invariants. DF-37.A is unaffected: Bob's late write does not erase the migration's existence. DF-37.B is unaffected: Bob is not editing pre-migration messages, he is writing a new one. DF-37.C is unaffected: the old transport's existence and frontier-CID are still preserved.

The scenario *does* surface a discrepancy that the operational-shape TE will need to handle: was Bob's message written before or after the migration's frontier-CID? If it was meant to be pre-migration but Bob's clone was stale, his message is now logically post-migration but written into the old transport's directory. Cooperative resolution probably involves Bob (or another participant) re-publishing his message into the new transport with a citation back to his late old-transport message. The invariants do not block this; they just ensure that whatever resolution the operational-shape TE picks, it cannot involve editing or erasing Bob's message silently.

## Conclusions

1. **The three invariant DFs are locked.** DF-37.A (audit-trail reconstructibility from disk + git history); DF-37.B (pre-migration message bytes immutable); DF-37.C (old transport's pCID, existence, and frontier-CID preserved forever, contents erasable only under a recorded erasure-receipt).

2. **Together the three invariants form the floor for any future migration-contract TE.** The operational-shape TE that eventually locks how migrations are encoded on disk (close-old + open-new vs. overlap window vs. atomic swap; back-reference format; seal mechanics; group-identity continuity; trigger-condition discipline; authorizing-promise shape) must satisfy all three. The operational-shape TE may add further constraints; it may not relax any of these three.

3. **The invariants are PromiseGrid-level, not wire-lab-specific.** The wire-lab is the simulation harness in which these invariants are first observable, but they are direct consequences of TE-zalut § S5, TE-dajot's 100-year goal, and the audit-substrate framing locked across TE-nibar / TE-lozip / TE-dabol / TE-vudaf. Any production substrate that calls itself PromiseGrid would inherit them with whatever encoding fits its on-disk-equivalent.

4. **The operational-shape TE is anticipated, not scheduled.** Locking it now would design against scenarios we are guessing at. It will be written when (or if) a concrete migration becomes imminent. Its scope is the seven DFs originally drafted in this TE's first scope (operational shape, back-reference, message disposition, authorizing promise, seal, group-identity continuity, trigger-condition discipline) plus whatever the imminent migration surfaces.

## Decision status

LOCKED:
- DF-37.A — Alt-A.1 (audit-trail reconstructibility from disk + git history is a hard requirement on any future migration contract).
- DF-37.B — Alt-B.1 (pre-migration message bytes are immutable; redaction handled through annotation, not editing).
- DF-37.C — Alt-C.2 (old transport's existence, pCID, and frontier-CID preserved forever; contents may be intentionally erased only under a recorded erasure-receipt).

Recorded principle: *transport-protocol migrations must produce reconstructible-from-disk audit trails; pre-migration message bytes are immutable; the old transport's existence and frontier-CID are preserved forever, with intentional erasure (if any) recorded as an erasure-receipt rather than effected silently.*

## Implications for follow-on work

- **Cat-3 Refinement on TE-junil** dis-anchoring the "TE-dajot (anticipated): transport-protocol migration semantics" forward-pointer in TE-junil's Implications-for-follow-on-work and § S7, replacing it with a forward-pointer to this TE (TE-numan). The Refinement is appended to TE-junil's `## Refinements` section per the TE editing policy locked in TE-dabol.
- **Update to `protocols/wire-lab.d/specs/transport-spec-draft.md`** resolving `OQ-2 (deferred)` against this TE's locked invariants, with a forward-pointer to the future operational-shape TE for the rest.
- **Future TE (slot-TBD): transport-protocol migration operational shape.** The substantive operational decisions (close-old + open-new vs. overlap window vs. atomic swap; back-reference format; disposition of pre-migration messages; seal mechanics; group-identity continuity; trigger-condition discipline; authorizing-promise shape) are deferred to that TE. It is anticipated, not scheduled, and will be written when a concrete migration is imminent.
