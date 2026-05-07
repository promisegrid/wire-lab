# Pre-149 Dropped-Thread Audit
## Session ea135ce8, Turns 18-148 (2026-04-29 through 2026-05-02)

> **Status:** Reference document. Audit performed 2026-05-05 to address Steve's
> concern about lost threads from compressed context in the period before turn
> 149. Five drops identified (DT-1 through DT-5); each is captured as a
> `UT-pre149.dtN` entry in `TODO-lilar-session-replay-cleanup.md`
> with cross-reference to this report for full detail. This document is the
> primary source for the meta-observations (TE-dajot slot collision pattern;
> AGENTS edit freeze spillover) that would lose narrative coherence if
> inlined into the ledger.
>
> **Maintenance:** This is a snapshot, not a living document. If subsequent
> walks discover additional pre-149 drops, file them as additional
> `UT-pre149.dtN+1` entries with brief inline summaries; do NOT back-edit
> this report.


**Audit date:** 2026-05-05  
**Scope:** Conversation turns 18-148 (approximately 7,200 lines of conversation.md)  
**Sources consulted:**
- `past_session_contexts/sessions/2026-05-04_2026-05-10/ea135ce8/conversation.md` (turns 18-148)
- `wire-lab/protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md` (UT-* ledger)
- Git log spanning 2026-04-28 through 2026-05-03 (93 commits)
- Corpus snapshot (all committed TE, DI, DR, TODO, spec files at end of session)

---

## Methodology

For each substantive thread in turns 18-148, the auditor checked in order:
1. Is it in the UT-* ledger? (The ledger has NO pre-149 entries; all UT-* entries address turns 155+.)
2. Is it in a committed TE, DI, DR, TODO, spec, or AGENTS file?
3. Is it in the harness-spec-draft or group-session-draft?
4. If not in any of the above: DROPPED.

The UT-* ledger's absence of pre-149 entries is structural. The session-replay-cleanup TODO was filed after turn 154; it does not walk turns 18-148 at all. That means any thread dropped in that range is invisible to the ledger.

Conservatism rule applied: if a thread is even partially captured in a TE's ## Refinements, a locked DI, or a AGENTS-ppx.md rule, it is marked WRITTEN even if the coverage is thin.

---

## Coverage Summary

| Category      | Approximate count |
|---------------|-------------------|
| WRITTEN       | ~95 turns         |
| CONVERSATIONAL| ~35 turns         |
| DROPPED       | 5 threads         |
| PARKED        | 0 (TE-havib twig is turns 149+, out of scope) |
| REVERTED      | 0                 |

Total accounted: ~135 turns across 148-18+1 = 131 turns. (Some turns contribute multiple threads; a few turns are pure plumbing with no new content.)

No meta-finding is triggered: the count of genuine dropped threads is 5, well under the 15-thread alert threshold.

---

## Confirmed WRITTEN Items (selected high-value entries)

The following were verified in the corpus and are noted here as a positive control
record. Items not listed were either conversational or are captured in the same
commit cluster as an adjacent item.

**Turns 22-24 -- ppx/{twig} naming, email identities, ppx prefix discipline**
Committed in AGENTS-ppx.md (ppx/ branch prefix, stevegt+ppx@t7a.org / stevegt@t7a.org
identity rules, per-twig naming convention).

**Turn 32 -- Never force push rule**
In AGENTS-ppx.md "Things that are forbidden." One-time exception (turn 78 force push
after the ASCII-only cleanup) is acknowledged in that same section. Not a dropped
thread.

**Turn 39 -- pCID = protocol CID, not promise CID**
Captured in AGENTS-ppx.md Glossary. Also appears in TE-rujak, TE-lozip, harness-spec
preamble.

**Turn 39 -- "discrepancy" not "defect"**
Captured in TE-mirah, harness-spec sec.10a, AGENTS-ppx.md vocabulary rules.

**Turn 41 -- "I promise that..." claim wording**
Captured in TE-mirah (review reply as promise), DI-009 (group-session spec), and
harness-spec sec.10a. The basic first-person promise form is WRITTEN.

**Turn 42 -- Promises as assertions of state (past/present/future)**
Captured in TE-nibar (spec-doc-as-promise), TE-rujak (spec-doc-store), TE-lozip
(congruence-convergence), harness-spec sec.2. The foundational past/present/future
tense framing is WRITTEN.

**Turns 44-46 -- "assertion/assert" not "burden/claim"**
Committed at 218a00e in the eradicate-burden twig. Harness-spec and TE-mirah both
use assertion vocabulary.

**Turn 50 -- Merge own branches to ppx/main, delete branch after merge**
Captured in AGENTS-ppx.md ppx/main lifecycle section.

**Turns 54-62 -- harness-spec as protocol doc, TE-nibar/22/23 cluster**
All committed: TE-nibar at dd69e3a area, TE-rujak, TE-lozip. Freeze ceremony (Alt-G
simrepo CHANGELOG approach) resolved by TE-zukug. CID computation tooling (DF-22.2
Go with go-cid, DF-22.5 freezer+checker same binary) captured in DI-011.

**Turn 66 -- hCID/fCID decision (per-protocol internal naming)**
Captured in TE-lozip DF-23.2 Alt-2.B: names are protocol-internal. Not a dropped
thread.

**Turns 68-72 -- congruence/convergence essay**
Committed at d784057. Essay at docs/essays/congruence-convergence-and-the-grid.md.

**Turns 82-101 -- channels->transports rename, Parents header, DAG message graph,
topology axes, T1.A/T2.A/T3.B/T4.A/Q1.A decisions**
Entire cluster committed. TE-hogus (numbering), TE-titur (numbering collision), TE-zalut
(transport types, committed with in-place rewrite), TE-junil (rename and axes,
committed at 5115f12). Decisions T1.A (Parents header), T2.A (space-separated),
T3.B (optional), T4.A (only), Q1.A (drop Kind), Q2 (cumulative-ack deferred to
OQ-100.3) all in DI-009 in TODO-bisur and group-session-draft spec.

**Turn 78 -- ASCII-only commit messages**
Locked in AGENTS-ppx.md "Things that are forbidden."

**Turns 107-120 -- protocols/ directory as simrepos, transports/{t-pcid}/{msg-pcid}
layering, UDP-binding binding layer, ns-3 simulation tooling**
Captured in TE-vipir (committed d4c82fc), TODO 018 (UDP-binding implementation),
TODO 019 (ns-3 harness scaffold, committed 18cc0f9).

**Turns 122-135 -- per-protocol TODO numbering, TE-magup/31/32/33 cluster**
All committed. TE-magup (renumbering, 9d037fc), TE-zukug (spec-doc inversion, c95c112),
TE-liviv (spec vs implementation split, 1522e5e), TE-potar (informative references,
e5f3c6a). OQ-100.3 (cumulative-ack) carried to TE-dajot and deferred notes in
TODO-013.

**Turns 139-148 -- TE editing policy, holistic corpus reading, tabletop TE**
TE-dabol committed at d841526, TE-vudaf committed at 6be8ea0. Three DIs locked
(DI-020-20260502-213103/-213104/-213105). TODO 020 filed with subtasks including
020.9 (the Alice/Bob/Carol tabletop scenario list). AGENTS edits explicitly HELD
pending 020.9 outcome (turn 143-144 choice A).

**Turn 111 -- 100-year goal**
TE-dajot filed as TE-dajot-100-year-goal-as-design-constraint.md.

**Turn 129-130 -- spec-CID in simrepo CHANGELOG (Alt-G, OQ-29.1 resolved)**
TE-zukug closes OQ-29.1 with the CHANGELOG inversion approach.

---

## DROPPED Threads

Five threads have no repo trace and no UT-* ledger entry.

---

### DT-1: Six-field tense-aware claim form and discourse protocol RFC spec

**Turns:** 41-43  
**Status:** DROPPED -- no TE, no DR, no DI, no TODO, no UT-* entry.

**What was discussed.** Turn 41 established "I promise that I..." as the claim form.
Turn 42 extended this: "in any claim, be clear about past, current, and/or future
states." The bot elaborated a six-field claim structure (promiser / tense /
state-description / condition / evidence / assessability criterion) and noted this
mapped to a tense-aware discourse protocol. The bot explicitly deferred writing it:
"Write the discourse-protocol document. Premature." It also said: "The right artifact,
when it's time, is one discourse-protocol spec document, prose-first, as if it's an
RFC. One spec, one (eventual) pCID."

**What was captured.** The basic "I promise that I..." first-person form is WRITTEN
in TE-mirah and harness-spec sec.10a. Mentions of past/present/future tense framing appear
in TE-nibar and TE-lozip (in the context of spec-doc temporal shape, not the six-field claim
form). The harness-spec sec.10a.2 names endorse-v1 / contest-v1 / counter-propose-v1 as
pCIDs but does not structure them around the six-field claim form.

**What is missing.** The six-field tense-aware structure was never filed as a TE, DI,
or TODO. The commitment to "eventually write one prose-first discourse-protocol spec
document" was made but was never filed as a backlog item. No UT-* entry exists.

**Recommended action.** Add to the TE-havib unpark checklist or create a standalone
TODO: "Discourse protocol RFC spec: single prose-first document covering the six-field
claim form (promiser / tense / state / condition / evidence / assessability) as a
formal protocol with a pCID. Predecessor to the endorse/contest/counter-propose
vocabulary currently in harness-spec sec.10a." Cross-reference: TE-havib DF-36.4
(discourse vocabulary home) was supposed to address where this vocabulary lives, but
DF-36.4 is still unlocked on the parked twig (UT-TE36-PARKED, UT-160.d).

---

### DT-2: Layer-pertinence indexing TE (proposed as TE-dajot, displaced)

**Turns:** 105-106  
**Status:** DROPPED -- no TE filed, no UT-* entry, only partial mitigation in TE-vipir layout.

**What was discussed.** Turn 105 (Steve): "I'm starting to wonder if we need to
somehow tag, index, or move the various TE, DR, DI, TODO etc. files according to
which protocol (pCID) they pertain to." The bot caught the concrete trigger: TE-hogus
had just been discovered to be layer-ambiguous. Bot proposed "Spawn a TE: TE-dajot
(anticipated): Indexing TE/DR/DI/TODO files by layer and protocol-pertinence." The
proposal included YAML frontmatter / docs/INDEX.md / "pertains-to" tagging as
candidate mechanisms.

**What happened instead.** Turn 107 shows Steve switching to critique of the
protocols/ tree layout. Turn 111 permanently claimed TE-dajot as the 100-year goal TE.
The protocols-as-simrepos layout (TE-vipir) addresses the concern partially: per-protocol
.d/ directories (wire-lab.d/TE/, wire-lab.d/DR/, etc.) provide natural containment of
TEs by protocol. But the explicit YAML frontmatter / docs/INDEX.md / cross-protocol
pertinence-tagging idea was never written as a TE or TODO.

**What is missing.** Whether the per-protocol .d/ directory layout fully satisfies the
original concern, or whether a cross-protocol index is still needed, was never formally
settled. No UT-* entry.

**Recommended action.** Evaluate whether TE-vipir's per-protocol .d/ layout closes the
concern. If yes, file a DI noting the resolution. If not, file a TE for "cross-protocol
TE/DR/DI pertinence tagging and indexing discipline."

---

### DT-3: Transport-protocol migration semantics TE

**Turns:** 87, 95, TE-junil anticipated work section  
**Status:** DROPPED -- anticipated as TE-dajot in TE-junil; slot displaced; no replacement filed; no UT-* entry.

**What was discussed.** Turn 87: "Migration between transport-protocol-pCIDs is
closing one transport and opening a new one... a real commitment with operational
implications." TE-junil (committed 5115f12) explicitly named an anticipated follow-on at
S7: "This question is non-trivial and TE-junil does not lock answers to it. It is deferred
to a future TE on transport-protocol migration semantics" and listed it as "TE-dajot
(anticipated): transport-protocol migration semantics -- what does it mean for a group
of participants to move from one transport-protocol-pCID to another?"

**What happened instead.** TE-dajot was claimed by the 100-year goal (turn 111).
TE-vipir (protocols-as-simrepos) became the next numbered TE. TE-junil's S7 open question
now points at a numbered slot that holds unrelated content. OQ-100.1 and OQ-100.2 in
TE-dajot (100-year) touch adjacent forking questions but do not answer S7 (what does
migration mean operationally for participants already using a transport-protocol-pCID).

**What is missing.** The S7 question is open. No TE has been filed. No TODO or UT-*
entry tracks it. TE-junil's anticipated-work pointer is stale.

**Recommended action.** File a new anticipated TE entry (the next available number
in the TE sequence is above TE-vudaf) for transport-protocol migration semantics. Update
TE-junil S7's "deferred to TE-dajot" note to point at the new anticipated number. Cross-
reference: TE-dajot OQ-100.1 (protocol forking) and OQ-100.2 (signing migration) are
adjacent topics that may merge into the same TE.

---

### DT-4: Conditional promise trust-ledger rule ("unfired conditionals do not move the ledger")

**Turns:** 41-42  
**Status:** DROPPED -- no TE, no DI, no harness-spec section captures this rule; no UT-* entry.

**What was discussed.** Turn 42 (Steve): "promises are assertions of state in the
past, present, or future, often conditional." The bot derived: "An unfired conditional
promise whose condition never fires is unassessed -- it is neither kept nor broken; it
does not move the trust ledger in either direction." This is a substantive trust-ledger
design rule with real implications for how the ledger counts interactions and computes
trust scores.

**What was captured.** The basic conditional promise framing appears in harness-spec
sec.2 and TE-nibar/TE-rujak in the sense that promises can be future-tense and conditional.
However, the specific rule that unfired conditionals are unassessed (neither kept nor
broken, no ledger movement) is not present in any committed document.

**Check performed.** grep for "unfired", "conditional.*ledger", "unassessed" across
all committed spec and TE files returned no matches. The harness-spec sec.2 TrustLedger
struct has fields `kept`, `broken`, `open_promises`, and `score` but no field or
prose that covers the unassessed/unfired case explicitly.

**What is missing.** The rule has architectural weight: it affects how open_promises
is treated when a condition never fires, how trust score computation handles the
case, and how break-witnesses interact with unfired conditionals. None of this is
written. No UT-* entry.

**Recommended action.** File as a DI candidate against harness-spec sec.2.2 or as an
OQ in TE-nibar/TE-rujak (via Cat-3 Refinements per the locked TE-dabol policy). The rule is:
"A conditional promise whose triggering condition never fires is unassessed: the
promiser's kept and broken counts do not change; the promise remains in open_promises
until the condition's evaluation deadline passes, at which point it is removed without
ledger movement." This needs a DF before locking.

---

### DT-5: Alice/Bob/Carol/Dave/Ellen/Frank convention in AGENTS-ppx.md

**Turns:** 145, 148  
**Status:** DROPPED from AGENTS-ppx.md specifically. Convention is implemented in TE-vudaf
but not enshrined in standing rules for future bot sessions. Minor; not a design
decision.

**What was discussed.** Turn 145 (Steve): "use Alice, Bob, Carol, Dave, Ellen, etc.
characters in TEs when characters are needed or useful." The bot acknowledged this
as a going-forward rule: "Saved. Going forward, TE tabletops on this repo use Alice /
Bob / Carol / Dave / Ellen / Frank / Grace / Heidi / Ivan / Judy / etc. for cooperative
actors, with Mallory reserved for adversaries." TE-vudaf (committed 6be8ea0) correctly
implements this convention in its scenario actors section (line 27: "Actors in the
tabletop follow the cryptography-literature alphabetical convention").

**What is missing.** AGENTS-ppx.md does not contain this convention. The bot said
"Going forward" -- implying a standing rule -- but the AGENTS edits for TE-dabol policy
were explicitly held (turn 143-144 choice A), and the Alice/Bob rule was never
separately committed to AGENTS-ppx.md as its own change. A future bot session
will not see this rule in AGENTS-ppx.md.

**Why this was held.** The AGENTS edit freeze from turn 143-144 was specifically
for the TE-dabol editing policy rules; the Alice/Bob convention is independent and
would not have been blocked by that freeze. It was simply not committed separately.

**Recommended action.** Add one bullet to AGENTS-ppx.md "Conventions" or
"Tabletop scenarios" section: "When naming actors in TE tabletop scenarios, use
the cryptography-literature alphabetical convention: Alice, Bob, Carol, Dave,
Ellen, Frank, ... for cooperative actors; Mallory for adversaries; Steve when
his role as repo owner is load-bearing." This is a Cat-1 AGENTS edit, no DF
needed.

---

## Items Confirmed Not Dropped (possible false alarms)

**Turn 36 -- Subagent simulation of peers capability discussion**
CONVERSATIONAL. Steve asked whether the bot could simulate multiple peer agents in
a session. The bot answered descriptively. No design decision was made.

**Turn 55 -- "pCIDs change when uncertainty changes"**
WRITTEN. Captured in TE-nibar's spec-doc-as-promise framing: the pCID hash over the
spec content means any change to the spec produces a new pCID. The conceptual point
is encoded in the freeze machinery.

**Turn 66 -- Freeze ceremony deep-dive (chicken-and-egg)**
WRITTEN at the decision level. TE-zukug closes OQ-29.1 with Alt-G (simrepo CHANGELOG
approach). The multi-alternative analysis was only in conversation, but the decision
is captured.

**Turns 63-64 -- Hand-computed CID discussion**
CONVERSATIONAL clarification about tooling. No design decision.

**Turns 99-101 -- Cumulative-prefix ack (Q2)**
WRITTEN as a deferral. OQ-100.3 in TE-dajot (100-year) and a deferred note in TODO-013.
The decision was to defer, and the deferral is recorded.

**Turn 113 -- "Proposals should move into transports"**
WRITTEN. Captured in TODO 016 (BLOCKED) and TE-vipir (ppx-dr as simulated protocol
under the protocols/ layout). The proposal is tracked.

**Turn 110 -- git commit hash in spec doc before computing pCID**
WRITTEN at the decision level by TE-zukug Alt-G (CHANGELOG inversion). The
multi-alternative freeze ceremony analysis was conversational, but the chosen
alternative is recorded.

---

## Meta-Observations

**UT-* ledger gap.** The session-replay-cleanup TODO (the UT-* ledger) was filed
starting at turn 154 and walks forward from there. It has zero entries covering
turns 18-148. The present audit is the first systematic coverage of that range.
If the ledger is later extended backward, this report should be the primary source.

**TE-dajot slot collision caused two dropped threads (DT-2 and DT-3).** Both the
layer-pertinence indexing TE and the transport migration TE were anticipated as
TE-dajot before that slot was claimed by the 100-year goal. When a numbered slot is
claimed by a different topic, anticipated follow-ons pointing at that slot are
silently orphaned. Recommendation: when a numbered slot is reassigned, the prior
anticipation references should be updated in the same commit (Cat-3 Refinement on the
TEs that contained the stale anticipated-TE pointers).

**AGENTS edit freeze spillover (DT-5).** The decision in turn 143-144 to hold
all AGENTS edits pending the 020.9 tabletop TE was correct for the TE-dabol editing
policy rules. However, it inadvertently captured an independent unrelated rule
(Alice/Bob/Carol convention) in the same freeze. Independent AGENTS additions should
be committed separately rather than bundled with policy updates under a hold.

**TE-havib twig (out of scope but relevant).** The parked TE-havib twig contains DF-36.4
(discourse vocabulary home), which is the intended resolution vehicle for DT-1 (the
discourse protocol RFC spec and six-field claim form). Resolution of DT-1 is therefore
gated on the TE-havib unpark process (UT-TE36-PARKED in the ledger).

---

## Summary Table

| ID   | Thread                                              | Turns   | Recommended action                                     |
|------|-----------------------------------------------------|---------|--------------------------------------------------------|
| DT-1 | Six-field tense-aware claim form + discourse RFC    | 41-43   | Add TODO; cross-ref TE-havib DF-36.4 unpark              |
| DT-2 | Layer-pertinence indexing TE (TE-dajot displaced)      | 105-106 | Evaluate TE-vipir sufficiency; file TE or DI              |
| DT-3 | Transport migration semantics TE (TE-junil S7)         | 87, 95  | File new anticipated TE; update TE-junil S7 pointer       |
| DT-4 | Conditional promise / unfired conditionals rule     | 41-42   | File as OQ / DI candidate against harness-spec sec.2.2   |
| DT-5 | Alice/Bob convention not in AGENTS-ppx.md           | 145     | One-bullet AGENTS-ppx.md edit, Cat-1, no DF needed    |

---

*Report generated by audit subagent. No wire-lab files were modified. No commits were made.*
