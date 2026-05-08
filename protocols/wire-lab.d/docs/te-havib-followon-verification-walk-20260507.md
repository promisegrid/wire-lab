# TE-havib follow-on verification walk

**Date:** 2026-05-07 (session continuation)
**Author:** stevegt-via-perplexity (bot)
**Trigger:** TODO-lilok (TE-havib follow-on: OQ-36.6 + tabletop walk).
Steve chose Alt-lilok.1.B — verification walk over close-as-superseded.
**Inputs:**

- `docs/thought-experiments/TE-havib-apparatus-vs-specimen-carve-out.md`
  (decided 2026-05-05; all 7 DFs locked; OQ-36.6 resolved in the negative)
- `protocols/wire-lab.d/docs/audit-20260503-015309-harness-spec-apparatus-vs-specimen.md`
  (committed 2026-05-03 at `4725b3e`)
- `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md`
  § "TE-havib follow-on" cluster (5 UTs)
- `protocols/wire-lab.d/TODO/TODO-lilok-te-36-followon-oq-and-tabletop-walk.md`

**Purpose:** verify, not re-decide. Walk each scenario and each UT
against the locked state of TE-havib and confirm or refute the lock.
Per Alt-lilok.1.B, this catches the case where the prose-level lock
fails to hold under scenario walk.

**Lock summary (target of the walk):**

| DF | Locked alt | What it commits |
|---|---|---|
| 36.1 | Alt-1.A Strict | Harness-spec mentions no candidate envelope, library API, struct shape, or named pCID. Apparatus-level invariants only. |
| 36.2 | Alt-2.A revised | Retire promise-stack as a separate hypothesis. `protocols/promise-stack.d/` is NOT created. |
| 36.3 | Alt-3.A | `protocols/grid-envelope.d/` is the home for `grid([pcid, payload])`, transport-agnostic. |
| 36.4 | Alt-4.A + PT-recast | Discourse vocabulary moves to `protocols/ppx-dr.d/`. Five named pCIDs collapse into two PT primitives: `conditional-promise` and `assessment` (with polarity). |
| 36.5 | Alt-5.C Both | Apparatus-level §1.3 summary stays in harness-spec; specimen-specific detail migrates (target shifts under DF-36.2 retirement). |
| 36.6 | Alt-6.D Lazy | Apparatus invariants only in §2.1; `protocols/trust-ledger.d/` deferred until a second candidate shape arrives. |
| 36.7 | Alt-7.A modified | §10 directory table: drop `promstack` row entirely, add `grid-envelope` row, lazy-defer `trust` row, apparatus rows stay. |

OQ-36.6 resolved in the negative (line 231): payload-recursion under
`grid <pcid> <payload>` is the universal nesting mechanism per
TE-lozip §3.1 + framing essay §3.1; promise-stack was invented
machinery at the wrong layer.

---

## Scenario walk

The audit memo (line 105) recommended six tabletop scenarios. TE-havib
shipped six (S1-S6, lines 48-82). I walk each against the locked
state to verify the lock holds.

### S1 — Alice the new contributor reads the spec corpus today

**Audit memo recommendation:** "(1) reading §1.1 today" — what does a
fresh contributor conclude from §1.1 in its current state.

**TE-havib S1 framing:** Alice opens README.md, follows the link to
harness-spec-draft.md, reads §1.1, concludes "the wire format is
`[]Promise`, encoded as a CBOR array, accessed via promstack."

**Test under DF-36.1 Alt-1.A Strict + DF-36.2 retirement:** after the
sweep, §1.1 contains no `Promise` struct, no `[]Promise` shape, no
CBOR-array prescription, no `promstack` API name. Apparatus-level
text only — "the harness studies multiple candidate envelopes; see
per-protocol directories under `protocols/<slug>.d/` for live specimens."

**Walks:** Alice reads §1.1 and finds an apparatus-level statement
plus pointers to `protocols/grid-envelope.d/` (DF-36.3) and a note
that promise-stack was retired (DF-36.2 lock annotation per DF-36.7's
§10 row). She does NOT come away with "the wire format is X." She
follows the per-protocol pointer to grid-envelope.d to learn the
current candidate envelope. **Lock holds.** S1 is the audit-memo
scenario (1) "reading §1.1 today" composed with scenario (2)
"reading §1.1 after carve-out" — TE-havib correctly merged them by
making each scenario a per-alternative test.

**Caveat:** the actual sweep edit to §1.1 has not yet landed. Per
DF-36.7's lock note, the sweep is part of TODO-vuhuj's leftover
sweep, not part of TE-havib itself. So today, on disk, §1.1 still
prescribes the old vocabulary. The lock is correct; the sweep is the
work item.

### S2 — Bob the second-implementer attempts a clean-room rebuild

**Audit memo:** did not name this scenario. TE-havib added it.

**TE-havib S2 framing:** Bob writes a second `promstack`
implementation in a different language and needs an authoritative
spec-doc CID to reference per TE-liviv.

**Test under DF-36.2 retirement:** under DF-36.2 Alt-2.A revised,
`protocols/promise-stack.d/` is NOT created. There is no spec-doc
CID for "the promise-stack envelope" because the hypothesis was
retired. Bob's clean-room rebuild target shifts — he is now
implementing the grid-envelope hypothesis (DF-36.3 Alt-3.A), and the
spec-doc CID lives at `protocols/grid-envelope.d/specs/grid-envelope-draft.md`'s
frozen pCID once that spec is drafted.

**Walks:** Bob's CHANGELOG entry references `protocols/grid-envelope.d/specs/grid-envelope-draft.md@<pCID>`,
NOT harness-spec. Tightly scoped per TE-liviv. **Lock holds.** S2
correctly motivates DF-36.3's transport-agnostic home for the grid
envelope.

**Caveat:** scenario text says "writes a second `promstack`
implementation" — under DF-36.2 retirement, this scenario's
narrative is now historically obsolete (it presumes promise-stack is
the candidate). The lock is correct, but the scenario prose
references retired vocabulary. Cat-1b historical-quotation: leave it
in TE-havib (it explains why DF-36.2 retired promise-stack), but
note that the scenario-as-written is illustrative of the pre-lock
state.

### S3 — Carol the second envelope hypothesis

**Audit memo:** "(3) attempting to add a second envelope hypothesis."

**TE-havib S3 framing:** Carol arrives with `[pCID,
salted-canonical-CBOR-payload]` with mandatory canonicalization +
optional outer signature wrapper. Wants to file as peer hypothesis,
run through C1-C7 ToC scenarios, see how it scores.

**Test under DF-36.1 Alt-1.A Strict + DF-36.3 Alt-3.A:** Carol's
spec doc lives at `protocols/<carol-envelope>.d/specs/<slug>-draft.md`.
She does not fork harness-spec. She does not edit any registry of
envelopes (there is none under DF-36.7's locked table). The harness
discovers her spec via the harness's standard per-protocol
discovery mechanism (cross-reference §10's apparatus-row pointer at
the per-protocol-MANIFEST locations).

**Walks:** Carol's act is additive. She creates
`protocols/carol-envelope.d/specs/carol-envelope-draft.md`, scaffolds
a TODO and CHANGELOG, and the harness's existing apparatus
infrastructure (acceptance scenarios C1-C7, ingress models, agent
profiles) runs against her envelope by the same mechanism that runs
against grid-envelope. No file Steve owns alone needs editing.
**Lock holds.** S3 is structurally satisfied by the strict carve-out
+ per-protocol home pattern.

### S4 — Dave the promise-stack forker

**Audit memo:** "(4) attempting to fork the promise-stack hypothesis."

**TE-havib S4 framing:** Dave likes the promise-stack family but
wants a variant with mandatory innermost-signature instead of
TE-famar's still-open ordering. Per C-4 this is a normal fork.

**Test under DF-36.2 retirement:** under DF-36.2 Alt-2.A revised,
promise-stack is retired. Dave's fork target no longer exists as a
candidate envelope. The S4 scenario is now **moot** in its original
framing.

**Reframe:** the architecturally equivalent scenario is "Dave the
grid-envelope forker" — Dave wants a variant of `grid([pcid,
payload])` with a different inner-payload-recursion convention. Per
C-4, this is a normal fork. He creates
`protocols/grid-envelope-dave.d/specs/...` (or equivalent slug per
TE-vipir per-protocol-simrepo shape). Apparatus and other specimens
are untouched.

**Walks:** the original S4 narrative is invalidated by DF-36.2
retirement, BUT the structural property the scenario tested (the
unit of fork is the per-protocol directory, not the harness) is
preserved by the architecture and demonstrably holds for the
grid-envelope-as-fork-target rephrasing. **Lock holds; scenario
prose is now historical-quotation.**

### S5 — Mallory the bad-faith carve-out attacker

**Audit memo:** "(5) Mallory injecting a malicious envelope claim."

**TE-havib S5 framing:** Mallory submits a PR claiming "the
apparatus-vs-specimen carve-out makes §1.1 obsolete; here is a
deletion of §1.1 with no replacement."

**Test under DF-36.1 Alt-1.A Strict:** the editing-policy DI this
TE produces requires that harness-spec sweeps which delete content
must have a destination protocol directory ready to absorb the
content. Mallory's deletion-with-no-replacement violates this rule
and is rejected at review.

**Walks under DF-36.2 retirement specifically:** the DF-36.2 lock
explicitly drops the §1.1 promise-stack content — it doesn't
relocate it. Does this fail Mallory's test? No: the relocation rule
applies to **specimen** content. DF-36.2 found that §1.1 was
**non-specimen invented machinery** (re-implementing payload-recursion
at the wrong layer). The lock annotation in §10 points at the framing
essay §3.1 + TE-lozip §3.1 for the conceptual replacement. So the
content is replaced by pointers, not deleted into a void. Mallory's
deletion-with-no-replacement still fails review because there is
no equivalent pointer-to-explanation in her PR.

**Lock holds.** S5 correctly catches Mallory under the strict
carve-out, and the DF-36.2 retirement does not create a Mallory
loophole.

### S6 — Ellen the 30-years-later contributor

**Audit memo:** "(6) a 30-years-later contributor finding the
harness-spec but not the carved-out protocol specs."

**TE-havib S6 framing:** Ellen arrives in 2056 with no living mentor,
finds harness-spec-draft.md, reads it, tries to understand the wire
format. Multiple promise-stack and grid-envelope variants have been
forked and retired.

**Test under DF-36.1 Alt-1.A Strict:** harness-spec is honestly
silent on envelope identity — Ellen has to look at
`protocols/<various>/` to learn what is currently live.

**Walks:** Ellen reads harness-spec, finds apparatus-level claims
about scenarios, ingress models, ledger invariants, but no concrete
envelope. She follows §10's apparatus rows + per-protocol pointers
to discover the live specimens. She does NOT come away assuming §1.1
is canonical (because §1.1 contains apparatus claims, not envelope
prescriptions, after the sweep). **Lock holds.**

**S6 was the deciding scenario for DF-36.1's lock per line 100:**
"the S6/Ellen stale-canon failure mode under Alt-1.B and the
OQ-36.6-robustness argument were the deciding factors over the
original bot recommendation rationale." S6 is what forced Steve to
override the bot's earlier rationale.

---

## UT-by-UT verdicts

### UT-160.b — DF-36.4 PT-recast structure

**Original concern:** "DF-36.4 PT-recast collapsed five named pCIDs
into two PT primitives, baked into Alt-4.A rather than presented as
its own DF."

**Verdict (verification walk):** the PT-recast is itself the locked
content of DF-36.4 Alt-4.A — see the table at lines 144-151 mapping
each legacy pCID to the PT primitive (`conditional-promise` or
`assessment`-with-polarity). Steve confirmed Alt-4.A on 2026-05-03.
Promoting the PT-recast to its own DF would have meant presenting
"should the discourse vocabulary use PT primitives or keep the
five-pCID set?" as a separate decision — but that was not a live
choice, since Steve's directive in DF-36.4 framing was already to
use Promise Theory vocabulary throughout. Presenting it as its own
DF would have been ceremony for a decision Steve already gave by
naming PT in the DF-36.4 framing. **Lock holds; observation is
procedural-meta, not a live open question.**

### UT-160.c — Tabletop framing mismatch

**Original concern:** "TE-havib's six tabletop scenarios used
different framing than the audit memo's recommended Alice-through-
Mallory walk."

**Verdict (verification walk):** the audit memo's six recommended
scenarios are (1) reading §1.1 today, (2) reading §1.1 after
carve-out, (3) attempting to add a second envelope hypothesis, (4)
attempting to fork the promise-stack hypothesis, (5) Mallory
injecting a malicious envelope claim, (6) a 30-years-later
contributor. TE-havib's S1-S6 are Alice (reads spec today, tested
under each carve-out alt — covers memo's (1) and (2) folded), Bob
(second implementer — added value, not in memo), Carol (second
envelope hypothesis — memo's (3)), Dave (promise-stack forker —
memo's (4)), Mallory (bad-faith carve-out attacker — memo's (5)),
Ellen (30-years-later contributor — memo's (6)). Five of six map
directly; the sixth (memo (2)) is folded into S1's per-alternative
test rather than being a standalone scenario; TE-havib added S2 Bob
as additional value. **Concern was wrong on inspection. Lock holds.**

### UT-161.a — Envelope asymmetry observation

**Original concern:** "the asymmetry between the two envelope
hypotheses — grid-pcid-payload can host promise-stack as one
possible payload, but promise-stack cannot cleanly host grid-
pcid-payload at envelope level — is..." [truncated]

**Verdict (verification walk):** under DF-36.2 retirement,
promise-stack is no longer a candidate envelope. There is no second
envelope to be asymmetric with. The asymmetry observation was
itself the structural argument that drove DF-36.2 to retire
promise-stack: the asymmetry exists because promise-stack was
re-implementing payload-recursion at the wrong layer. The
observation is preserved in OQ-36.6's resolution (line 231) as
"invented machinery re-implementing payload-recursion at the wrong
layer." **Resolved by lock; no live question.**

### UT-162.a — OQ-36.6 unfinished

**Original concern:** "the deferred OQ-36.6 investigation is itself
a major unfinished thread."

**Verdict (verification walk):** TE-havib line 10 status: "OQ-36.6
resolved in the negative." Line 231 details the resolution citing
TE-lozip §3.1 + framing essay §3.1. Concern was true at the time
the cluster was captured (2026-05-06 ledger pass against pre-Alt-B
state), false against the current 2026-05-05 Alt-B locks. **Resolved
by lock.**

### UT-162.b — DF-36.2 provisional lock vulnerability

**Original concern:** "DF-36.2 is provisionally locked pending
OQ-36.6, but the conditionality may not survive visual inspection
of the TE file."

**Verdict (verification walk):** DF-36.2 was re-presented under
Alt-B disposition path on 2026-05-05 and locked at Alt-2.A revised
(retire promise-stack). The Alt-B re-presentation IS the visual
inspection the concern asked for; Steve confirmed retirement after
the re-presentation. Not provisional anymore. **Resolved by lock.**

---

## Conclusions

All six tabletop scenarios verify against the locked state. Five
match the audit memo recommendation directly; one (S2 Bob) is value-
add. All five UTs in the cluster are answered by the 2026-05-05
Alt-B locks. The verification walk catches no failure mode.

S2 and S4 reference promise-stack vocabulary in their narrative
prose, which is now historical-quotation (Cat-1b) under DF-36.2's
retirement. The scenarios remain valid as illustrations of the
structural property each tests; the narrative pCID names are stale
but kept for prose continuity per TE-dabol Cat-1b policy (this
scenario file's Refinement section, if added, would record the
recast).

The deciding scenario for DF-36.1's strict lock was S6 Ellen — the
walk re-confirms why: only Alt-1.A Strict structurally prevents the
stale-canon failure mode, and DF-36.2's retirement preserves the
S6 conclusion (Ellen looks at per-protocol directories to learn
live specimens; harness-spec is honestly silent on envelope
identity).

## Open follow-on items surfaced by the walk

None that affect the locks. Two procedural notes:

1. **The actual sweep edits to §1.1, §2.1, §7.1, §10a.2/.3/.6, §10
   directory table have not yet landed.** Per DF-36.7 lock annotation
   (line 207), they are part of TODO-vuhuj's leftover sweep, not
   part of TE-havib's lock. The verification walk is verifying the
   *decision*, not its realization in the spec text. The spec sweep
   remains as work — but it is properly tracked in TODO-vuhuj, not
   in TODO-lilok.

2. **The S2 Bob and S4 Dave scenario narrative prose references
   promise-stack vocabulary that is now retired.** Adding a brief
   `## Refinements` entry to TE-havib explaining that the narrative
   uses pre-DF-36.2-retirement vocabulary (Cat-1b historical-
   quotation) would close the loop with TE-dabol's Cat-3 navigational
   policy. This is a small TE-havib refinement, not a follow-on
   TODO.

## Disposition of TODO-lilok

Per Alt-lilok.1.B the verification walk has run. Result: no live
work surfaces. TODO-lilok closes as **verified-superseded** —
distinct from "close-as-superseded" (Alt-lilok.1.A) only in that
the verification was actually run rather than asserted. The five
UTs and the cluster in `dropped-thread-disposition-20260506.md`
absorb this memo as their close-out.

---

## Sources

- TE-havib: `docs/thought-experiments/TE-havib-apparatus-vs-specimen-carve-out.md`
  (decided 2026-05-05; commits on twig + Alt-B re-presentation chain)
- Audit memo: `protocols/wire-lab.d/docs/audit-20260503-015309-harness-spec-apparatus-vs-specimen.md`
  at commit `4725b3e`
- TE-lozip §3.1: `docs/thought-experiments/TE-lozip-congruence-convergence-duality-and-pcid-framing.md`
- Framing essay §3.1: `docs/essays/congruence-convergence-and-the-grid.md`
- Nested-vs-stacked-envelopes research: `docs/research/nested-vs-stacked-envelopes-20260504.md`
- TE-dabol Cat-1b policy: `docs/thought-experiments/TE-dabol-frozen-history-and-editing-policy.md`
- DI-020-20260502-232651 (Cat-1a/Cat-1b classification rule)
- TODO-lilok: `protocols/wire-lab.d/TODO/TODO-lilok-te-36-followon-oq-and-tabletop-walk.md`
- Dropped-thread disposition file: `protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md`
