# TODO-juhub: Turns 149-208 chronological rewalk

## Prior aliases

This TODO is being filed after the TE-mumuv proquint migration locked
2026-05-07, so it is minted directly under `TODO-juhub`. No prior integer or
timestamp alias.

## Status

Open. This TODO owns the one-turn-at-a-time raw-log rewalk for turns 149-208.
`TODO-lilar` remains the historical original walk through turn 192, and
`TODO-jivam` remains the closure gate for overall recovery completion. No turn
may advance until the current turn has been analyzed, reported, and explicitly
approved.

## Decision Intent Log

ID: DI-nagat
Date: 2026-05-10 17:11:45
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: File `TODO-juhub` as the authoritative one-turn-at-a-time rewalk
ledger for turns 149-208. For each turn, read the raw session log, sweep all
later turn logs and later repo artifacts for whether that turn's questions,
decisions, or plans were later settled or changed, report the result, and stop
until explicitly approved to continue.
Intent: Rebuild confidence in the 149-208 recovery from raw evidence instead
of trusting earlier summaries, while keeping the historical walk and the
closure monitor distinct and readable.
Constraints: Do not use `TODO-jivam` as the per-turn ledger. Do not rewrite
existing `TODO-lilar` turn notes; if a correction is needed there, append a
provenance-bearing correction note. Do not advance more than one turn without
explicit approval. Every unresolved finding must be handed to a proper owner
artifact.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/TODO/TODO-jivam-turns-149-170-recovery-completion.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`protocols/wire-lab.d/TODO/TODO.md`;
`/home/stevegt/lab/session-logs/sessions/ea135ce8/149-turn.md` through
`/home/stevegt/lab/session-logs/sessions/ea135ce8/208-turn.md`.

ID: DI-pijun
Date: 2026-05-10 19:37:26
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Promote the turn-173 Usenet-lineage insight into two non-TE docs
with distinct roles: the research note keeps the historical-precedent
evidence, and a new simulation README explores the broader CAS /
Usenet-like / git-like design line. `group-session` is one current specimen in
that broader line, not the identity of the simulation itself.
Intent: Preserve the strongest historical analogy in a design-visible place
without pretending it is already a frozen protocol decision, while giving the
broader architectural framing a simulation-local home that can grow beyond the
current `group-session` specimen.
Constraints: Do not file a new TE for this promotion. Put the research-doc
wording in the Usenet section of
`docs/research/historical-networks-20260503.md`. Use the phrase
`content-addressed Usenet`. Create `simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md`
only; do not create a fuller simulation scaffold in this pass. Add the new
simulation to `simulations/README.md`.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`docs/research/historical-networks-20260503.md`;
`simulations/README.md`;
`simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md`.

ID: DI-gudap
Date: 2026-05-10 13:16:35
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Strengthen the `TODO-juhub` replay discipline so each turn
reconciliation also resolves or explicitly calls out its related loose ends.
Intent: Prevent the 149-208 replay from becoming a second narrative-only pass.
The rewalk must reduce relitigation by reconciling each turn's `TODO-lilar`
`UT-*` fallout against the current owner TODOs and by promoting any
load-bearing insight that is still trapped only in replay artifacts.
Constraints: Do not flip `TODO-lilar` checkboxes directly. Respect the
matrix-as-closure-index rule in `ut-verification-matrix-20260507.md`; close,
retire, or transfer loose ends only in the proper owner artifact. Update
`TODO-lilar` only with additive correction notes when its historical walk note
is actually wrong. Direct design/spec/research/simulation docs are touched only
when a specific turn exposes a load-bearing statement that is still missing
from those docs. Update `DEV-GUIDE-RESOURCES.md` only when such a cited or
relevant source changes.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/TODO/TODO-jivam-turns-149-170-recovery-completion.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`protocols/wire-lab.d/TODO/dropped-thread-disposition-20260506.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`;
current owner TODOs for the turns being replayed; directly implicated
spec/design/research/simulation docs when needed.

ID: DI-vanak
Date: 2026-05-10 13:36:01
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: During the `TODO-juhub` replay, Steve's single-word reply `turn`
counts as approval of the currently presented turn analysis and authorizes
rewriting that turn's note in the stronger `TODO-juhub` format before
advancing.
Intent: Remove repetitive confirmation chatter so the replay can proceed one
turn at a time while still preserving the explicit approval boundary.
Constraints: `turn` approves only the currently presented turn; it does not
skip turns. It authorizes rewriting the current turn's `TODO-juhub` note and
any already-described turn-local owner, correction-note, or direct-doc cleanup
required by that approved analysis. The replay still stops after the next turn
is presented for review.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
this session's `TODO-juhub` replay cadence.

ID: DI-nijod
Date: 2026-05-10 13:45:42
Status: superseded
Author: stevegt@t7a.org (Steve Traugott)
Decision: Every per-turn replay response must include an explicit `Work
pending` line near the bottom.
Intent: Make it immediately obvious whether a turn still leaves live `UT-*`,
owner TODO, spec-edit, DR/DI, or other substantive follow-on work, so the
replay steadily burns down loose ends instead of forcing Steve to infer status.
Constraints: `Work pending: yes` when any open `UT`, owner TODO item,
spec-edit, DR/DI, or other substantive work still stems from the turn.
`Work pending: no` only when the turn is fully reconciled and no live owner
work remains because the related loose ends are absent, closed, retired, or
transferred. This line is a report field; it does not itself close work.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
this session's `TODO-juhub` replay responses.

ID: DI-firap
Date: 2026-05-17 10:49:50
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Interpret `Work pending` as recovery-walkthrough pending work, not as
all downstream design work. If a turn's loose ends are fully captured in sim
questions for PromiseGrid design or TODOs for harness work, mark `Work pending:
no`.
Intent: The recovery walkthrough should close replay accountability once the
right downstream container exists. Running sims, answering PromiseGrid design
questions, or completing harness TODOs is follow-on work; it should not force
the same historical turn to remain pending after capture.
Constraints: `Work pending: yes` only when the recovery walkthrough still lacks
a sufficient downstream container or still needs a turn-local cleanup decision.
`Work pending: no` does not mean downstream design or harness work is done; it
means the loose end has been captured in the proper sim question or TODO owner.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
this session's `TODO-juhub` replay responses.
Supersedes: DI-nijod (`Work pending` yes/no semantics only)

ID: DI-vumir
Date: 2026-05-16 10:17:35
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Every per-turn replay final response must visibly include the full
`Turn N plain-English recap` text for the turn, not merely a path/line pointer
to the rewritten note.
Intent: Keep each `turn` handoff self-contained enough for Steve to confirm the
actual prompt/response/conclusion summary without opening the file, while still
allowing file references as supporting evidence.
Constraints: A path/line reference is allowed only as an addition, never as a
substitute for the recap text. The visible recap must be the full recap field
from the rewritten note or a faithful same-detail rendering of it. This rule
applies to every bare-`turn` response before the assistant yields.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
this session's `TODO-juhub` replay responses.

ID: DI-vaguf
Date: 2026-05-17 12:01:50
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Reconcile turn 178 by routing its remaining loose ends into existing
simulation, DR, resource, and replay-owner artifacts rather than opening a new
TE or creating new simulation roots.
Intent: Turn 178 mixed several design constraints with a deliberately terse bot
deferral. The replay should make those constraints visible without treating the
old 15-DF TE-sihih expansion as current scope. Sparse CAS, pull-decision
accounting, interop, promisebase prior art, CIDv1 object typing, group identity,
BGP-class app pressure, multi-repo site topology, and dogfood urgency each need
a downstream container so the recovery walkthrough can advance without forcing
the historical turn to stay open.
Constraints: Do not answer the open design questions in this pass. Do not create
new protocol trees, world fixtures, manifests, local TODO queues, or simulation
roots. Preserve the later TE-sihih scope contraction: turn-178 material is routed
to successor sims, DRs, and guide-resource notes instead of backloading the
broad capture narrative into TE-sihih.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`;
`DR/DR-napum-promisegrid-layperson-guide-claims.md`;
`simulations/SIM-zazit-chunk-feed-replication/`;
`simulations/SIM-rusap-promise-accounting-records/`;
`simulations/SIM-jurar-cas-backed-group-session/`;
`DEV-GUIDE-RESOURCES.md`.

ID: DI-lusum
Date: 2026-05-17 13:04:21
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Repair the turn-178 routing by adding explicit question homes and
TODO owners for the gaps found in review, while preserving append-only DR
history.
Intent: Turn-178 recovery is complete only when BGP-class app pressure, group
identity / anti-default-anonymity pressure, `pgmsg` tool naming,
collaborator-permission constraints, and promisebase human-readable reference
naming are discoverable from their downstream owners rather than only from
scenario rows or closed replay notes. The DR event log must keep historical
event wording and append terminology clarifications instead of rewriting older
events.
Constraints: Do not answer the open design questions in this pass. Do not
create new simulation roots. Route BGP and group identity into simulation
`QUESTION.md` files, route `pgmsg` / collaborator-permission meta work to a
root harness TODO, route promisebase reference naming separately from CBOR,
chunking, and CIDv1 object typing, and keep `Work pending` semantics tied to
whether each loose end has a downstream sim question or TODO owner.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`protocols/wire-lab.d/TODO/TODO-rohub-dogfood-tool-name-and-collaborator-permission.md`;
`protocols/wire-lab.d/TODO/TODO-kituj-te-43-promisebase-prior-art-adoption.md`;
`protocols/wire-lab.d/TODO/TODO.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`;
`DR/DR-napum-promisegrid-layperson-guide-claims.md`;
`simulations/SIM-rusap-promise-accounting-records/QUESTION.md`;
`simulations/SIM-jurar-cas-backed-group-session/QUESTION.md`;
`simulations/SIM-jomag-cas-object-model/QUESTION.md`;
`simulations/SIM-jomag-cas-object-model/SCENARIOS.md`;
`simulations/README.md`;
`DEV-GUIDE-RESOURCES.md`.
Supersedes: DI-vaguf (routing-completeness and append-only DR repair only)

ID: DI-tibis
Date: 2026-05-17 13:23:58
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Split the turn-178 BGP-class routing question and promisebase
reference-naming question into their own standalone simulations.
Intent: BGP-class routing is not merely an example row inside peer-local promise
accounting; it is its own application-pressure experiment. Promisebase
human-readable reference naming is not merely part of the L6 CAS object-model
bundle; it is its own reference-resolution experiment. Each needs an
independent simulation root so the competing design lineages can evolve without
being hidden under broader accounting or object-model sims.
Constraints: Do not answer either question in this pass. Keep `SIM-rusap`
focused on generic peer-local promise accounting records. Keep `SIM-jomag`
focused on CAS object bytes, object typing, pointer objects, CBOR, and chunking.
Route BGP-class routing to `SIM-punaz-bgp-class-routing-app/` and promisebase
reference naming to `SIM-ligan-promisebase-reference-naming/`.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`protocols/wire-lab.d/TODO/TODO-kituj-te-43-promisebase-prior-art-adoption.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`;
`DEV-GUIDE-RESOURCES.md`;
`simulations/README.md`;
`simulations/SIM-rusap-promise-accounting-records/README.md`;
`simulations/SIM-rusap-promise-accounting-records/QUESTION.md`;
`simulations/SIM-rusap-promise-accounting-records/SCENARIOS.md`;
`simulations/SIM-jomag-cas-object-model/README.md`;
`simulations/SIM-jomag-cas-object-model/QUESTION.md`;
`simulations/SIM-jomag-cas-object-model/SCENARIOS.md`;
`simulations/SIM-punaz-bgp-class-routing-app/`;
`simulations/SIM-ligan-promisebase-reference-naming/`.
Supersedes: DI-lusum (BGP and promisebase-reference simulation-home routing only)

ID: DI-vabij
Date: 2026-05-17 13:35:22
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Reconcile turn 179 by routing promise-economy mechanism agnosticism
to a standalone simulation, conditional release / geofencing to existing
TODO-ralud and SIM-zarud owners, and the promisebase wholesale-adoption /
cross-repo-scope lessons to existing process and promisebase-trajectory owners.
Intent: Turn 179 introduced a real PromiseGrid design question that was still
too broad for `SIM-rusap`: the protocol must be able to test a spectrum from
peer-local social assessment to capability-token marketplaces without baking in
one economics model. It also introduced a procedural failure: the assistant read
promisebase design docs as if they were implementation evidence and proposed
cross-repo promisebase documentation scope without authorization. The replay must
preserve those lessons without reviving the invalid wholesale-adoption pivot.
Constraints: Do not answer the economics-model question in this pass. Do not
make promisebase authoritative. Do not create cross-repo work. Keep
conditional-release / geofencing with TODO-ralud and SIM-zarud, keep
promisebase adoption and prototype-vs-canon with TODO-kituj / TODO-dozak, and
use `SIM-haros-promise-economy-spectrum/` for the standalone mechanism-spectrum
question.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`protocols/wire-lab.d/TODO/TODO-dozak-te-44-wire-lab-promisebase-merge-trajectory.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`;
`simulations/README.md`;
`simulations/SIM-haros-promise-economy-spectrum/`;
`DEV-GUIDE-RESOURCES.md`.

ID: DI-pidag
Date: 2026-05-17 13:57:29
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Correct turn-179 cleanup so `SIM-haros-promise-economy-spectrum/`
has a TODO owner and TODO-lilar remains historical evidence rather than the
closure authority for replay routing.
Intent: The turn-179 replay can close only if every loose end is routed to a
discoverable owner. A standalone simulation is not enough by itself, and directly
checking off `TODO-lilar` UT rows mutates the older replay artifact instead of
using the newer TODO-juhub / matrix owner-routing discipline.
Constraints: Keep `DI-vabij`'s design routing active. Do not answer the
promise-economy mechanism question in this pass. Do not make promisebase
authoritative. Do not create cross-repo work. Restore TODO-lilar UT checkboxes to
their historical open state and use TODO-juhub, TODO-rajig, successor TODOs, and
the UT matrix for closure evidence.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`protocols/wire-lab.d/TODO/TODO-rajig-promise-economy-spectrum.md`;
`protocols/wire-lab.d/TODO/TODO.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`;
`simulations/README.md`;
`simulations/SIM-haros-promise-economy-spectrum/`;
`DEV-GUIDE-RESOURCES.md`.

ID: DI-zotol
Date: 2026-05-17 14:02:41
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Route turn-180's promisebase correction lessons to existing process,
prior-art, prototype-canon, and replay-owner artifacts without creating new
PromiseGrid design work.
Intent: Turn 180 corrected the invalid turn-179 wholesale-adoption pivot by
making the assistant audit promisebase code and git history before making
architecture claims. The durable work is procedural and routing-oriented:
code-first / ground-truthed citation, apology-audit-invalidate-propose recovery,
promisebase-as-prototype rather than authority, and the independent-design path
that later turn 191 makes explicit.
Constraints: Do not treat promisebase as authoritative. Do not reopen the
turn-179 invalid DFs. Do not flip TODO-lilar UT checkboxes. Keep open
PromiseGrid design questions with TODO-kituj, TODO-dozak, TODO-rajig, TODO-ralud,
`DR-tumus`, TE-nizor, and related simulations; keep procedural lessons with
AGENTS-ppx B6/B7 and the replay matrix.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`AGENTS-ppx.md`;
`docs/thought-experiments/TE-nizor-pahah-implementation-sufficiency.md`;
`protocols/wire-lab.d/TODO/TODO-kituj-te-43-promisebase-prior-art-adoption.md`;
`protocols/wire-lab.d/TODO/TODO-dozak-te-44-wire-lab-promisebase-merge-trajectory.md`.

ID: DI-zarok
Date: 2026-05-17 14:05:30
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Route turn-181's promisebase `db/` audit findings to the existing
TE-43 / `DR-tumus` CAS-profile owner, the CAS-object and chunking simulations,
and the later promisebase regression-fix evidence.
Intent: Turn 181 produced useful code-first evidence about pitbase's implemented
CAS/Merkle/Rabin primitives, but it also introduced a wrong import path and an
over-strong working frame that wire-lab could adopt promisebase `db/` directly as
L6 substrate. The replay must preserve the evidence while routing unresolved
choices to the current owner artifacts and later corrections: correct repo path,
chunking parameters, type binding, pointer-object shape, test-status threshold,
and promisebase-as-prior-art rather than authority.
Constraints: Do not answer `DR-tumus` here. Do not make promisebase or pitbase
authoritative. Do not flip TODO-lilar UT checkboxes. Treat turn 181's
`github.com/t7a/pitbase/db` import-path claim as wrong; use
`github.com/stevegt/promisebase/db` when an import path is actually needed.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`protocols/wire-lab.d/TODO/TODO-kituj-te-43-promisebase-prior-art-adoption.md`;
`DR/DR-tumus-turn-177-l6-cas-adoption.md`;
`simulations/SIM-jomag-cas-object-model/`;
`simulations/SIM-gobaz-chunking-identity-bakeoff/`;
`simulations/SIM-bobud-l6-cas-starting-profile-bakeoff/`;
`simulations/SIM-kohad-cas-object-type-binding-bakeoff/`.

ID: DI-lumal
Date: 2026-05-17 14:08:32
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Route turn-182's promisebase regression-fix proposal lessons to the
existing replay, authorization, PAT, diagnostic-transparency, and later
turn-183/186 execution evidence instead of creating new PromiseGrid design work.
Intent: Turn 182 was the bridge between the turn-181 `db/` audit and the
turn-183 randStream fix. Its durable lessons are procedural: terse user replies
can confirm both a regression hypothesis and a prior concrete offer; fix
proposals should show the diagnostic reasoning in the same answer; and cross-repo
work proposals must separate local preparation, PAT-gated push/PR operations, and
user-push alternatives.
Constraints: Do not answer `DR-tumus` here. Do not flip TODO-lilar UT
checkboxes. Do not create promisebase work from this replay pass. Treat the
actual fix, test result, local-only commit, PAT grant, push, and status-visibility
concerns as later-turn evidence owned by turns 183, 185, 186, and 188.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`AGENTS-ppx.md`.

ID: DI-kegar
Date: 2026-05-17 14:14:05
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Route turn-183's promisebase randStream fix evidence to the existing
promisebase prior-art owner while creating a separate harness TODO for
local-only cross-repo work persistence before PAT-gated pushes.
Intent: Turn 183 completed useful cross-repo repair work and made the
promisebase `db/` chunker/Merkle path green, but its durable replay loose ends
split into two classes: PromiseGrid design evidence for TODO-kituj / `DR-tumus`,
and a harness/process risk that a local-only commit could be lost before
authentication allowed it to be pushed. The replay should preserve the evidence
without treating promisebase as authoritative, and should not leave the
persistence rule only in a closed turn note.
Constraints: Do not touch promisebase in this replay pass. Do not answer
`DR-tumus`. Do not flip TODO-lilar UT checkboxes. Route the `db/` test evidence
and adoption threshold to TODO-kituj / `DR-tumus`; route code-first external-repo
claims to AGENTS-ppx B7 and the replay notes; route local-only work persistence
to `TODO-fapev`.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/TODO/TODO-fapev-cross-repo-work-persistence.md`;
`protocols/wire-lab.d/TODO/TODO.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`protocols/wire-lab.d/TODO/TODO-kituj-te-43-promisebase-prior-art-adoption.md`;
`DR/DR-tumus-turn-177-l6-cas-adoption.md`;
`AGENTS-ppx.md`.

ID: DI-nulak
Date: 2026-05-17 14:20:14
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Route turn-184's promisebase full-code audit, RFC-1005 prior-art
pressure, and process defects to their current owners without treating the
turn-184 ten-DF list as locked design.
Intent: Turn 184 corrected the earlier promisebase-docs overreach by auditing
actual code, tests, RFCs, and ROADMAP state. The durable result is bounded
evidence: promisebase `db/` and `kv/fs` are useful implemented prior art,
several other promisebase surfaces were partially rotten or unfinished, RFC-1005
is a promise-economy prior-art seed, and the assistant still violated process by
listing ten DFs flat, surfacing a collaborator name in a question, and making an
opaque pattern-count claim. Each class needs a durable owner so the replay can
move on without losing the lesson or converting the audit into plan-of-record
PromiseGrid design.
Constraints: Do not touch promisebase in this replay pass. Do not answer
`DR-tumus`. Do not reopen TE-sihih's later-contracted scope. Do not flip
TODO-lilar UT checkboxes. Keep promisebase / pitbase as prior art unless the
later TE-43 / `DR-tumus` path decides otherwise. Route one-DF cadence to
AGENTS-ppx B5, collaborator propagation to TODO-rohub / AGENTS-ppx B3,
pattern-count reproducibility to AGENTS-ppx B7, promisebase partial-rot and
kv/fs dependency-target questions to TODO-kituj / `DR-tumus`, and RFC-1005
promise-economy prior art to TODO-rajig / `SIM-haros`.
Affects: `protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`protocols/wire-lab.d/TODO/TODO-rohub-dogfood-tool-name-and-collaborator-permission.md`;
`protocols/wire-lab.d/TODO/TODO-kituj-te-43-promisebase-prior-art-adoption.md`;
`protocols/wire-lab.d/TODO/TODO-rajig-promise-economy-spectrum.md`;
`simulations/SIM-haros-promise-economy-spectrum/QUESTION.md`;
`simulations/SIM-haros-promise-economy-spectrum/SCENARIOS.md`;
`DEV-GUIDE-RESOURCES.md`;
`AGENTS-ppx.md`.

ID: DI-lifub
Date: 2026-05-17 14:27:25
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Route turn-185's two-PAT operational pattern and credential-hygiene
loose ends into AGENTS-ppx B4 rather than leaving them only in replay notes.
Intent: Turn 185 correctly established that separate wire-lab and promisebase
PATs can be held and used in one session, but the answer under-specified three
durable safeguards: short expiry as part of fine-grained PAT guidance, actual
scope/expiry verification rather than filename-based read-only conventions, and
redaction of secret bytes from carry-over or handoff summaries. The credential
rule needs to say those things directly because later turns depend on the
two-token pattern.
Constraints: Do not record or echo PAT bytes. Do not touch `.secrets/` or any
runtime credential path in this replay pass. Do not flip TODO-lilar UT
checkboxes. Treat per-remote token separation, short practical expiry, actual
scope/expiry verification, filename suffixes as non-enforcement, and summary
redaction as the closed process rule for turn 185.
Affects: `AGENTS-ppx.md`;
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`.

ID: DI-zagus
Date: 2026-05-17 14:33:29
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Route turn-186's promisebase PAT application, randStream push,
redaction application, collaborator re-quote hazard, and GOTOOLCHAIN diagnostic
to process owners without starting new promisebase fix work.
Intent: Turn 186 completed the first end-to-end cross-repo fix in this replay
slice and proved the two-PAT pattern works, but it also exposed three durable
process lessons: redaction must be applied when a turn contains literal secret
bytes, collaborator-sensitive questions must be paraphrased when re-quoted, and
cross-repo build/toolchain diagnostics such as `GOTOOLCHAIN=auto` drift need a
discoverable owner. The replay should preserve those lessons without recording
secret bytes, touching promisebase, or letting deferred Docker/FUSE fixes distract
from the recovery walkthrough.
Constraints: Do not record or echo PAT bytes. Do not touch `.secrets/`, runtime
credential files, or promisebase. Do not flip TODO-lilar UT checkboxes. Route
redaction to AGENTS-ppx B4 and the turn-185 rule, collaborator re-quote handling
to TODO-rohub / AGENTS-ppx B3, the GOTOOLCHAIN diagnostic to `TODO-nasat` and
AGENTS-ppx B7, and deferred Docker/FUSE work to existing TODO-kituj / `DR-tumus`
routing.
Affects: `AGENTS-ppx.md`;
`protocols/wire-lab.d/TODO/TODO-juhub-turns-149-208-chronological-rewalk.md`;
`protocols/wire-lab.d/TODO/TODO-nasat-cross-repo-build-hazard-capture.md`;
`protocols/wire-lab.d/TODO/TODO.md`;
`protocols/wire-lab.d/TODO/TODO-rohub-dogfood-tool-name-and-collaborator-permission.md`;
`protocols/wire-lab.d/TODO/TODO-kituj-te-43-promisebase-prior-art-adoption.md`;
`protocols/wire-lab.d/docs/ut-verification-matrix-20260507.md`.

## Why a separate TODO

`TODO-lilar` already records the original historical walk through turn 192 and
explicitly says turns 193+ belong in a successor TODO. `TODO-jivam` is the
bounded recovery closure gate and should stay focused on completion criteria
instead of becoming a long per-turn ledger. `TODO-juhub` therefore owns the
rewalk mechanics while the older artifacts remain evidence and closure logic.

## Per-turn discipline

1. Read the raw session log for the current turn.
2. Sweep every later turn log and every later relevant artifact (`TE`, `DI`,
   `DR`, `TODO`, specs, matrix/disposition docs, essays, and research notes)
   for whether the current turn's questions, decisions, or plans were later
   settled, corrected, superseded, contradicted, or abandoned.
3. Collect every related loose end for that turn from `TODO-lilar`,
   `dropped-thread-disposition-20260506.md`,
   `ut-verification-matrix-20260507.md`, and the current owner TODOs.
4. Determine whether each loose end is already resolved, retired, transferred,
   still open, or missing a proper owner.
5. Compare the turn plus its loose ends against the existing replay artifacts
   and decide whether any historical correction note is actually needed.
6. Report the turn back to Steve using the fixed report format below, visibly
   include the full `Turn N plain-English recap` text, and stop; do not advance
   another turn in the same response.
7. A subsequent bare `turn` approves the next turn's `TODO-juhub` rewrite and
   any linked owner, correction-note, or direct-doc updates required by that
   turn.
8. Do not advance until each related loose end is either closed, retired, or
   transferred in its proper owner artifact, or explicitly called out here as
   still needing a named decision or work item.

## Interaction shorthand

- `turn` means: approve the currently presented turn, rewrite that turn's note
  in the stronger `TODO-juhub` format, perform any already-described turn-local
  cleanup authorized by the approved analysis, then report the completed turn
  with the full plain-English recap text visible and name the next chronological
  turn.
- The `Write needed? yes/no` line in the turn report is informational. During
  this replay, it does not require a separate confirmation from Steve before
  rewriting the current turn note in `TODO-juhub`.

## Turn report format

- `Turn N plain-English recap`: summarize the user's prompt and the assistant's
  response for the turn, conclusions reached during the turn, later updates or
  modifications to those conclusions found in later turns, and any loose ends or
  open questions that remained as of the end of the turn. The final response to
  Steve must include this full recap text inline; a filepath/line reference alone
  is not compliant.
- `Existing capture`
- `Gaps or contradictions`
- `Related UTs / owners`
- `Owner/doc cleanup`
- `Remaining decisions or work`
- `Work pending: yes|no` — use `no` when loose ends are fully captured in sim
  questions for PromiseGrid design or TODOs for harness work; use `yes` only
  when replay cleanup still needs a downstream container or turn-local decision.
- `Proposed disposition`
- `Write needed? yes/no`
- `Next`: name the next chronological turn; do not request separate approval
  beyond Steve's next bare `turn`.

## Loose-end backfill through turn 177

- `Turns 149-154` No related `UT-*` entries were filed for this block. No
  downstream owner or direct-doc cleanup is currently needed beyond the
  already-landed TE-editing-policy and TODO-020 artifacts.
- `Turn 155` `UT-155.a` and `UT-155.b` are now retired in `TODO-kugod` /
  `TODO-rivuk` under `DI-runuh`; no live turn-local owner work remains.
- `Turn 156` `UT-156.a` and `UT-156.b` are retired / resolved-retired.
  `UT-156.c` is resolved under `DI-lajod`, `DI-sujan`, and `DI-kinad`;
  broader TE-40 audit work continues under the later turn-159 rows.
- `Turn 157` `UT-157.b` is retired. `UT-157.a` and `UT-157.c` are resolved in
  the grid-envelope successor owner under `DI-joroh`.
- `Turn 158` `UT-158.b` is resolved; `UT-158.e` and `UT-158.g` are retired;
  `UT-158.a` is resolved-decomposed and `UT-158.h` is resolved-routed under
  `DI-sotuk`; `UT-158.c`, `UT-158.d`, and `UT-158.f` are closed for
  turn-158 scope under `DI-kinad` and `DI-fanah`. Broader TE-40 work begins
  with turn 159 and is now closed under `DI-mugar`.
- `Turn 159` `UT-159.a`, `UT-159.b`, and `UT-159.d` are resolved in
  `TODO-kugod` under `DI-mugar`; `UT-159.c` remains resolved-retired by the
  TE-havib follow-on verification walk.
- `Turn 160` `UT-160.a` is resolved by the TE-havib Cat-3 refinement and
  `TODO-kugod` / `DI-mugar` nine-item sweep closure. `UT-160.b` and
  `UT-160.c` are answered by the TE-havib follow-on verification path.
  `UT-160.d` is resolved by TE-havib's final all-seven-DF locked status.
  `TODO-lilok`'s former reopened harness-spec-sweep note is now reconciled
  through `TODO-kugod` / `DI-mugar`.
- `Turn 161` `UT-161.a` is answered by the TE-havib follow-on disposition.
  `UT-161.b` and `UT-161.c` are now captured as historical inputs by the
  2026-05-12 TE-havib refinement; they are no longer turn-local live work.
- `Turn 162` `UT-162.a` and `UT-162.b` are answered by the later TE-havib
  disposition path. The former `TODO-lilok` reopened sweep-handoff note is
  now closed through `TODO-kugod` / `DI-mugar`.
- `Turn 163` `UT-163.a` is resolved by the apparatus-level `§1.3` rewrite in
  `harness-spec-draft.md` under `DI-lajod` / `DI-mugar`. `UT-163.b` is closed
  as a future-process rule by `AGENTS-ppx.md` B1 / `DI-021-20260507-212249`;
  the commit-specific concerns it cross-referenced live in their own UT rows.
  `UT-cbf7f41-fallback` is retired by the final OQ-36.6 / DF-36.2 negative
  path plus the active `§1.3` apparatus rewrite.
- `Turn 164` `UT-164.a` and `UT-164.b` are historical corrections, not live
  implementation work. `UT-164.c` is resolved by `TODO-bisur` 012.7's
  four-message round-trip. `UT-164.d` is closed by sim-local `TODO-gapab` /
  `DI-rurab`. `UT-164.e` is closed by `TODO-gapab` / `DI-rurab` and
  `TODO-kakaz` / `DI-bomud`; rewrite-at-freeze remains rejected, and turn 164
  now has zero open successor work.
- `Turn 165` `UT-165.a` is closed as an observational privacy/slug lesson;
  `UT-165.b` is closed by the Steve-authored DI promise shape in `DI-rurab`;
  `UT-165.c` was already closed by the neutral-memory update; `UT-165.d` is
  closed by keeping OQ-G4 deferred while treating m000 as valid specimen
  evidence; `UT-165.e` is closed by the group-session example/freeze-gate
  cleanup under `DI-rurab`.
- `Turn 166` `UT-166.a` is closed for future-process purposes by the current
  decision-first protocol and `DI-vanak`; `UT-166.b`, `UT-166.c`, and
  `UT-166.e` are closed by active `DI-rurab` / specimen wording; `UT-166.d`
  is closed as historical git metadata that must not be rewritten.
- `Turn 167` `UT-167.a` is closed by active filename=CID docs; `UT-167.b` and
  `UT-167.c` are closed by `DI-rurab` branch-membership wording; `UT-167.d`
  is closed by the 2026-05-14 `.msg` corpus audit; `UT-167.e` is closed for
  future-process purposes by the current decision-first protocol and
  `DI-vanak`.
- `Turn 168` `UT-168.a` is closed by active `Message-ID:` compatibility
  wording; `UT-168.b`, `UT-168.c`, and `UT-168.d` are closed by `DI-rurab`,
  current §8/§9 wording, and the infrastructure/message-file distinction;
  `UT-168.e` is closed for future-process purposes by decision-first plus
  `DI-vanak`; `UT-168.f` is closed for active docs because current
  wire-lab-devs docs use post-turn-169 CIDs and stale CIDs are historical-only.
- `Turn 169` `UT-169.a`, `UT-169.b`, and `UT-169.e` are closed by
  `DI-012-20260508-033513` / `DI-rurab` and active group-session §4.3 / §4.7
  compatibility wording; `UT-169.c` is closed for future-process purposes by
  decision-first plus `DI-vanak`; `UT-169.d` is closed as historical branch
  metadata because active twig rules require `ppx/{twig}` kebab-case, not the
  historical `ppx/te-<utc>-<slug>` pattern.
- `Turn 170` `UT-170.a` is closed by TE-sihih / TODO-vunub landing the
  L5/L6/L7 substrate-agnostic model and by DR-nugog / `DI-fakin` superseding
  the original flat-versus-nested root `transports/` question for the current
  specimen. `UT-170.b`, `UT-170.c`, and `UT-170.d` are retired and checked off
  as historical naming / stale-summary / superseded-framing records.
- `Turn 171` `UT-171.a` is closed by TE-sihih / TODO-vunub Q-22.2
  plus `DI-rurab` keeping §9 as the current specimen's inline normative git
  binding; `UT-171.b` is closed by TODO-vunub Q-22.3 retracting the manifest
  field idea in favor of path-as-declaration; `UT-171.c` and `UT-171.d` are
  closed as recorded design-cadence lessons.
- `Turn 172` `UT-172.a` is closed as a framing-stability lesson;
  `UT-172.b` is closed by TE-sihih's forward vocabulary (`feed` with
  `substrate` as prose term); `UT-172.c` is closed by TODO-vunub Q-22.3's
  path-as-declaration retraction of per-instance feed manifests; `UT-172.d`
  is closed by TE-sihih's L5/L6/L7 replacement taxonomy; `UT-172.e` is closed
  by `DI-rurab` keeping §9 inline for the current specimen and deferring any
  future feed-spec extraction to successor work.
- `Turn 173` `UT-173.a` is closed by `DI-pijun`, the historical-networks note,
  and `SIM-hugoj` preserving the content-addressed-Usenet line as exploratory
  design evidence; `UT-173.b` is closed by TE-sihih's `feed` vocabulary;
  `UT-173.c` is closed as a cadence lesson refined by turn 174; `UT-173.d` is
  closed by the bounded negative-precedent check added to the research note;
  `UT-173.e` is closed by the research note's CAS-cardinality entry and
  TE-sihih's L5/L6/L7 split.
- `Turn 174` `UT-174.a` is closed by TE-sihih's forward vocabulary and the
  historical-networks current-reading note; `UT-174.b` is closed by the
  `udp-feed` simulation lineage; `UT-174.c` is closed by TODO-vunub Q-22.3's
  path-as-declaration retraction of per-instance feed manifests; `UT-174.d` is
  closed by TE-sihih's L5/L6/L7 model and `feed` vocabulary; `UT-174.e` is
  closed by the revised Pattern A/B/C current-reading note.
- `Turn 175` `UT-175.a` is closed/transferred by the no-rewrite freeze/CID
  decisions and grid-envelope simulation owners; `UT-175.b` is closed for
  TODO-vunub by TE-sihih's actual smaller landing scope and transferred
  migration/CAS work; `UT-175.c` is closed by active group-session header rules;
  `UT-175.d` is closed by the historical-networks / SIM-hugoj Usenet-lineage
  capture; `UT-175.e` is closed by TE-sihih's L5/L6/L7 model; `UT-175.f` is
  closed as a recorded reversal/cadence lesson; `UT-175.g` is closed by the
  TE-sihih and historical-networks `feed` / `substrate` convention. Successor
  work remains open in TODO-pipus, TODO-kituj / TE-43, and root/simulation layout
  owners such as DR-nugog / TE-domat / later simulation-structure TEs.
- `Turn 176` `UT-176.i` is closed as a corrected provenance defect with no data
  loss. `UT-176.a`, `UT-176.b`, `UT-176.c`, `UT-176.g`, and `UT-176.h` are
  resolved or transferred by TE-sihih / TODO-vunub's actual L5-L6-L7,
  slug-state, `group`, and no-exception L6-CAS-pointer decisions. `UT-176.d`,
  `UT-176.e`, and `UT-176.f` are transferred to successor owners: SIM-hugoj /
  historical-network framing, TE-domat / DR-nugog root layout, grid-envelope /
  body-semantics work, TODO-pipus operational migration, and TODO-kituj / TE-43
  concrete L6 CAS adoption.
- `Turn 177` `UT-177.a` and `UT-177.b` are closed for TODO-vunub by TE-sihih's
  actual L5/L6/L7 and smaller-scope landing. `UT-177.g` is closed as a positive
  cadence lesson. `UT-177.c`, `UT-177.d`, and `UT-177.h` are transferred to
  TODO-kituj / TE-43's concrete L6 CAS adoption work, with TODO-pipus retaining
  operational migration. `UT-177.e`, `UT-177.f`, and `UT-177.i` are transferred
  to TODO-kulih / TE-nibar spec-doc-shape work and TE-dajot 100-year-goal
  pressure testing.
- `Turn 178` `UT-178.a`, `UT-178.b`, and `UT-178.c` are captured by the
  sparse-CAS / chunk-feed / promise-accounting sims plus TODO-kituj and
  successor owner TODOs; `UT-178.d` is captured by standalone
  `SIM-punaz-bgp-class-routing-app`; `UT-178.h`, `UT-178.i`, and `UT-178.l`
  are split between TODO-kituj / `DR-tumus`, the CAS object-model / bakeoff
  sims, and standalone `SIM-ligan-promisebase-reference-naming`; `UT-178.j` is
  captured by the chunk-feed multi-repo site-topology scenario; `UT-178.e` is
  routed to `DR-napum`; `UT-178.f` is routed to `SIM-jurar` and `DR-napum`;
  `UT-178.g`, `UT-178.k`, and `UT-178.m` close as dogfood/procedural replay
  lessons.

## Subtasks

- [x] juhub.149 Turn 149 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn confirmed the Cat-1a/Cat-1b split recommendation; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.150 Turn 150 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn confirmed the cross-TE quotation-grep safeguard for future Cat-2 vocabulary sweeps; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.151 Turn 151 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn confirmed the Cat-2 cross-TE quotation-grep safeguard; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.152 Turn 152 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn confirmed the top-of-file `Status:` field rule and the immediate unblocking of the follow-on TE-policy sweeps; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.153 Turn 153 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn authorized execution of the already-unblocked TODO-020 rollout work; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.154 Turn 154 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn was a queue-status and recommendation checkpoint, not a new design turn; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.155 Turn 155 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn clarified the stalled DF-1.1 proposal but did not lock it; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.156 Turn 156 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn interrupted the old TE-famar DF path with a scope correction; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.157 Turn 157 raw-log rewalk plus later-turn and later-artifact sweep. **Done 2026-05-10.** Raw turn corrected the transport-vs-envelope confusion and proposed a higher-level envelope-shape TE; existing capture was correct; no later contradiction found. See `## Turn notes`.
- [x] juhub.158 Turn 158 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.159 Turn 159 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.160 Turn 160 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.161 Turn 161 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.162 Turn 162 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.163 Turn 163 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.164 Turn 164 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.165 Turn 165 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.166 Turn 166 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.167 Turn 167 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.168 Turn 168 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.169 Turn 169 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.170 Turn 170 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.171 Turn 171 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.172 Turn 172 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.173 Turn 173 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.174 Turn 174 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.175 Turn 175 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.176 Turn 176 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.177 Turn 177 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.178 Turn 178 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.179 Turn 179 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.180 Turn 180 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.181 Turn 181 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.182 Turn 182 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.183 Turn 183 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.184 Turn 184 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.185 Turn 185 raw-log rewalk plus later-turn and later-artifact sweep.
- [x] juhub.186 Turn 186 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.187 Turn 187 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.188 Turn 188 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.189 Turn 189 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.190 Turn 190 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.191 Turn 191 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.192 Turn 192 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.193 Turn 193 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.194 Turn 194 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.195 Turn 195 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.196 Turn 196 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.197 Turn 197 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.198 Turn 198 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.199 Turn 199 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.200 Turn 200 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.201 Turn 201 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.202 Turn 202 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.203 Turn 203 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.204 Turn 204 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.205 Turn 205 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.206 Turn 206 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.207 Turn 207 raw-log rewalk plus later-turn and later-artifact sweep.
- [ ] juhub.208 Turn 208 raw-log rewalk plus later-turn and later-artifact sweep.

## Turn notes

### Turn 149 — 2026-05-03 00:05 UTC

- `Turn 149 summary` Steve's `yes` was a real decision, not a filler
  acknowledgement. The bot had just proposed splitting Cat-1 path renames into
  two cases: current live path references that can be updated mechanically, and
  historical quotations of old paths that must be preserved. Steve's approval
  locked that split.
- `Existing capture` `TODO-lilar` already records the turn correctly: the
  approval led to `DI-020-20260502-232651`, committed as `cd82c19`, merged as
  `d8c3e93`, pushed, then the conversation moved on to DF-35.2. Later policy
  artifacts (`TODO-dinub`, TE-dabol's Refinements, and the TE-editing-policy
  commits) still treat that Cat-1a/Cat-1b split as the live locked rule.
- `Gaps or contradictions` None found. Later turns and later artifacts do not
  retract, narrow, or correct the turn-149 decision. Turn 197's replay
  instructions even classify the early 149-153 block as straightforward
  confirmations of already-landed TODO-020 work.
- `Related UTs / owners` None. The replay/disposition chain begins its `UT-*`
  inventory at turn 155, so turn 149 has no downstream owner TODO to reconcile.
- `Owner/doc cleanup` None needed. The TE-editing-policy artifacts already carry
  this decision, and no later spec, research, simulation, or owner TODO is
  missing a load-bearing statement from this turn.
- `Remaining decisions or work` None.
- `Work pending` `no`
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 150 is next and remains
  pending approval.

### Turn 150 — 2026-05-03 00:23 UTC

- `Turn 150 summary` Steve's `yes` confirmed the second tightening of the TE
  editing policy. The issue here was not about paths anymore; it was about
  quoted words. If one TE quotes another TE's exact wording, and a later
  vocabulary sweep rewrites the quoted term everywhere, then the quote stops
  being historically true. This turn approved the rule that any future Cat-2
  vocabulary sweep must first grep the corpus for quoted uses of the old term
  and leave genuine historical quotations alone.
- `Existing capture` The old replay ledger already has the substantive outcome
  right even though it grouped the approval under the next line item's naming.
  The rule later landed as the Cat-3 Refinement recorded in TE-dabol and TODO
  020, with commit `04126ac` merged as `795a846`. Later artifacts explain the
  effect clearly: Cat-2 sweeps now have two required checks, not one. First,
  the sweeper must name the DIs whose meaning is unchanged. Second, the
  sweeper must grep for old-term-in-quotation contexts before rewriting.
- `Gaps or contradictions` No contradiction found. Later policy docs, TODO 020,
  TE-dabol's Refinements, and the replay notes all continue to describe this as
  a Cat-3 procedural tightening rather than a new superseding DI. The only
  nuance is that `TODO-lilar` phrases the turn boundary in a slightly shifted
  way, because the old walk grouped the DF confirmation block as a tight series
  of approvals. Substantively, the landing is still correct.
- `Related UTs / owners` None. The replay/disposition chain does not file any
  `UT-*` entries for turns 149-154, so there is no downstream owner TODO to
  reconcile for this turn.
- `Owner/doc cleanup` None needed. The TE-editing-policy artifacts already
  carry the rule family this turn confirmed, and no later owner or direct
  design/spec/research/simulation doc appears to be missing a load-bearing
  statement from this turn.
- `Remaining decisions or work` None.
- `Work pending` `no`
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 151 is next and remains
  pending approval.

### Turn 151 — 2026-05-03 00:28 UTC

- `Turn 151 summary` Steve's `yes` approved the next TE-editing safeguard:
  before any Cat-2 vocabulary sweep rewrites a term across the corpus, the bot
  must grep for uses of the old term inside quotation-like contexts and emit
  those matches for human review. In plain English, this is the rule that stops
  a vocabulary cleanup from silently rewriting quoted history and making an
  older TE appear to have said something it never said.
- `Existing capture` `TODO-lilar` already records the turn correctly as the
  confirmation of DF-35.3, the mandatory cross-TE quotation-grep step before a
  Cat-2 sweep. Later artifacts confirm the same result: the rule landed as a
  Cat-3 Refinement on TE-dabol, with the refinement text explaining that a
  Cat-2 sweeper now has a two-step protocol — enumerate unchanged DIs in the
  top-of-file note, then grep and classify quotation-context matches before
  rewriting. Later replay summaries keep the same interpretation.
- `Gaps or contradictions` None found. I did not find any later turn or later
  artifact that narrows, retracts, or contradicts the turn-151 decision. The
  later TE-35 summary still lists DF-35.3 as settled and merged, with the
  quotation-grep safeguard intact.
- `Related UTs / owners` None. The replay/disposition chain does not file any
  `UT-*` entries for turns 149-154, so there is no downstream owner TODO to
  reconcile for this turn.
- `Owner/doc cleanup` None needed. The TE-editing-policy corpus already carries
  this safeguard, and no later owner or direct design/spec/research/simulation
  doc appears to be missing a load-bearing statement from this turn.
- `Remaining decisions or work` None.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 152 is next and remains
  pending approval.

### Turn 152 — 2026-05-03 00:35 UTC

- `Turn 152 summary` Steve's `yes` closed the last open DF from TE-vudaf. The
  substance was the new top-of-file `Status:` field on every TE, placed where a
  reader sees it before reading the body. The problem being solved was stale
  supersedence discoverability: if the only "this TE is superseded" marker lives
  at the very bottom, a reader can miss it and act on stale reasoning. This
  turn approved the fix and immediately reframed the next work as execution
  tasks, not more policy debate.
- `Existing capture` `TODO-lilar` already records the turn correctly as the
  confirmation of DF-35.4 Alt-4.a: a uniform top-of-file `Status:` header field
  on every TE. The raw turn itself is mostly a summary table, but it clearly
  states the outcome: DF-35.4 landed as a Cat-3 Refinement on TE-dabol, the
  retrofit is subtask 020.10, and subtasks 020.5 / 020.6 / 020.7 / 020.10 moved
  from deferred to ready-to-execute. Later artifacts preserve exactly that
  framing. TE-dabol's Refinements define the field shape and purpose, and
  `TODO-dinub` records that 020.10 later added the field to all 35 existing TEs.
- `Gaps or contradictions` None found. The later corpus consistently treats
  this turn as the point where the TE-editing policy became fully settled: four
  DIs plus four Cat-3 Refinements, with the remaining work reduced to rollout
  and mechanical sweeps. No later artifact disputes the turn boundary or the
  meaning of the `Status:` field decision.
- `Related UTs / owners` None. The replay/disposition chain does not file any
  `UT-*` entries for turns 149-154, so there is no downstream owner TODO to
  reconcile for this turn.
- `Owner/doc cleanup` None needed. The TE-editing-policy corpus already carries
  this decision, and no later owner or direct design/spec/research/simulation
  doc appears to be missing a load-bearing statement from this turn.
- `Remaining decisions or work` None.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 153 is next and remains
  pending approval.

### Turn 153 — 2026-05-03 00:44 UTC

- `Turn 153 summary` Steve's `do it` was not another policy choice. It was
  authorization to execute the rollout work that turn 152 had just unblocked.
  In plain English, the TE-editing policy debate was over, and this turn is the
  bot reporting what it actually landed from that now-settled queue.
- `Existing capture` `TODO-lilar` already captures the important nuance
  correctly. The raw turn itself explicitly described two landed items: 020.7
  (the TE-famar `## Refinements` forward pointer) and 020.5 (the AGENTS
  rollout). But the same raw turn also said all four previously unblocked
  subtasks were now done, and the later TODO-020 state confirms the full batch
  included 020.6 (the Cat-1a/Cat-1b path-reference sweep) and 020.10 (the
  top-of-file `Status:` retrofit across all 35 TEs). So the raw turn's prose
  named two items in detail, while the later artifacts confirm the larger
  completion set.
- `Gaps or contradictions` None found. Later artifacts preserve the same
  interpretation: after this turn, only 020.8 remained open in TODO-020. I did
  not find any later correction saying the turn-153 batch was narrower than
  `TODO-lilar` records.
- `Related UTs / owners` None. The replay/disposition chain does not file any
  `UT-*` entries for turns 149-154, so there is no downstream owner TODO to
  reconcile for this turn.
- `Owner/doc cleanup` None needed. TODO-020 already carries the landing state
  of the rollout batch, and no later owner or direct
  design/spec/research/simulation doc appears to be missing a load-bearing
  statement from this turn.
- `Remaining decisions or work` None for turn 153 itself.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 154 is next and remains
  pending approval.

### Turn 154 — 2026-05-03 01:17 UTC

- `Turn 154 summary` Steve's `next?` did not open a new design thread by
  itself. The bot replied with a queue-status snapshot: here are the open
  TODOs, here is the recommended next move, and here is the suggested order.
  In plain English, this was a planning checkpoint between the now-finished
  TE-editing-policy rollout and the next substantive decision work.
- `Existing capture` `TODO-lilar` already records the turn correctly as a pure
  queue / recommendation boundary. The raw turn lists 020.8 as the small
  cleanup item, TODO-rivuk as the high-leverage next DF work, and several other
  open items ranging from DI/DR backfill to implementation scaffolding. The
  later replay notes correctly preserve two important nuances: this turn is not
  the start of TE-havib, and TODO-bisur was still visibly alive in the queue
  here with two open subtasks.
- `Gaps or contradictions` None found. Later notes explicitly correct the turn
  boundary in the same way: turn 154 is the queue-status turn, while the
  TE-havib / apparatus-vs-specimen sequence begins only after Steve's next
  response and the scope challenge that follows. I found no later artifact that
  reclassifies turn 154 as a substantive design decision.
- `Related UTs / owners` None. The replay/disposition chain does not file any
  `UT-*` entries for turns 149-154, so there is no downstream owner TODO to
  reconcile for this turn.
- `Owner/doc cleanup` None needed. This turn is a queue snapshot rather than a
  turn where a missing design/spec/research/simulation statement needs to be
  promoted, and there is no turn-local owner cleanup beyond the already-listed
  queue artifacts themselves.
- `Remaining decisions or work` None for turn 154 itself. The queue items named
  here remained open, but that is not a replay loose end created by this turn.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 155 is next and remains
  pending approval.

### Turn 155 — 2026-05-03 01:38 UTC

- `Turn 155 summary` Steve picked the TE-famar path from the queue, but he did
  not lock DF-1.1. Instead he pushed back on two missing pieces: what exactly
  `Project` means, and where any resulting decision would actually live. In
  plain English, this turn is the bot trying to rescue the old promise-stack
  DF flow by explaining the terms more concretely before asking for a lock.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer defines `Project(msg, predicate)` as a non-consuming, order-independent
  query over frames, contrasts it with `Peel(msg)`, and says the old intended
  lifecycle was: lock DI-005-1.1 in TODO-rivuk, later cite the resulting DIs
  from `harness-spec-draft.md §1.1`, then close DR-006. Later artifacts also
  preserve the crucial outcome: none of this actually locked, because turn 156
  immediately reframed the scope before Steve answered yes.
- `Gaps or contradictions` None found. The later corpus consistently treats two
  things as the important leftovers from this turn: first, DF-1.1 was never
  locked and the old TODO-rivuk queue was later superseded; second, the clearest
  `Project / Peel / Wrap` definitions in the corpus still live only in this
  conversation. The wording "projection mode" is also recognized later as a bad
  phrase that should not propagate into committed text; the real issue was
  whether `Project` is part of the spec contract.
- `Related UTs / owners` `UT-155.a` and `UT-155.b` are the turn-local loose
  ends. Both are now retired in `TODO-kugod`: DF-1.1 is no longer live
  promise-stack work, and the `Project / Peel / Wrap` vocabulary remains
  historical rather than active apparatus work.
- `Owner/doc cleanup` None needed now. The owner chain already retired the old
  TE-famar line correctly, and there is no direct-doc promotion needed from
  this turn unless that vocabulary is deliberately revived under a new active
  owner later.
- `Remaining decisions or work` None for turn 155 itself. The conversation text
  remains useful historical evidence, but there is no still-open replay owner
  action attached to this turn after the later retirement.
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 156 is next and remains
  pending approval.

### Turn 156 — 2026-05-03 01:39 UTC

- `Turn 156 summary` Steve cut off the old promise-stack DF flow before any
  lock happened. His point was that the topic under discussion was not a
  harness-level invariant; it was only one candidate wire-envelope design. In
  plain English, this turn is where the bot first admits the old TE-famar /
  TODO-rivuk framing is scoped wrong.
- `Existing capture` `TODO-lilar` already records the turn correctly as a
  mid-DF scope correction. The raw turn shows the bot responding by offering
  three structural choices for what to do with TE-famar's analysis: make it a
  harness-level default, demote it to a per-envelope concern, or split it into
  harness-level vocabulary plus per-envelope lockings. The later corpus also
  preserves the key outcome: none of those three options was actually locked,
  because turn 157 reframed the question again before Steve chose among them.
- `Gaps or contradictions` None found. Later artifacts preserve two critical
  follow-on facts consistently: first, the Option 1 / 2 / 3 menu was abandoned
  rather than answered; second, the bot's wording that the harness-spec should
  be "wire-envelope-agnostic" is itself later corrected as wrong. The later
  apparatus-vs-specimen framing says the harness is not envelope-agnostic; it is
  the apparatus that compares candidate envelopes and other layer choices.
- `Related UTs / owners` `UT-156.a` is retired and `UT-156.b` is
  resolved-retired under `TODO-rivuk` / `DI-runuh`; `UT-156.c` is now resolved
  in `TODO-kugod` under `DI-lajod`, which rewrites `harness-spec-draft.md §1.1`
  and `§1.3` at apparatus level and explicitly retracts the stale
  "wire-envelope-agnostic" wording.
- `Owner/doc cleanup` Done. `TODO-kugod` now records `UT-156.c` as resolved, and
  `harness-spec-draft.md` now states that the harness compares candidate
  envelopes rather than defining one canonical envelope. No correction note is
  needed in `TODO-lilar`; the former turn-158 slice is now closed, and the
  remaining broader sweep lives under the later `UT-159.*` rows in
  `TODO-kugod`.
- `Remaining decisions or work` None for turn `156` itself. The broader
  apparatus/specimen cleanup is still open, but it belongs to later turns and
  their owner items rather than to the turn-156 stale-claim residue.
- `Work pending` `no`
- `Proposed disposition` `already captured correctly`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 157 is next and remains
  pending approval.

### Turn 157 — 2026-05-03 01:42 UTC

- `Turn 157 summary` Steve corrected the bot again: transports and envelopes
  are different things, the goal is one transport-agnostic message envelope,
  and `grid([pcid, payload])` is only the current working hypothesis, not a
  proven answer. In plain English, this turn forces the bot to step back from
  the old promise-stack-only framing and admit that the envelope decision
  itself is still under study.
- `Existing capture` `TODO-lilar` already records the turn correctly. The raw
  turn names five candidate envelope shapes, explains why TE-famar jumped ahead
  by assuming the promise-stack family was already the right abstraction, and
  recommends "Reading 2": file a higher-level envelope-shape TE, treat TE-famar
  as misframed input to that larger decision, and gate the old TODO-rivuk DF
  queue behind the new envelope-shape work. Later owner cleanup also preserves
  the two load-bearing carry-forwards from this turn: the candidate-envelope
  inventory and the `grid([pcid, payload])` working-hypothesis prose move to the
  future grid-envelope successor work rather than staying attached to
  promise-stack.
- `Gaps or contradictions` None found. Later artifacts preserve this as a
  transitional correction, not as a final settled architecture. The five
  candidate envelopes and the Reading-2 recommendation are kept as historical
  evidence, but turn 158 immediately rejects the remaining assumption that the
  envelope belongs in the harness-spec as a single harness-wide decision. I did
  not find any later artifact claiming the turn-157 framing was itself the final
  answer.
- `Related UTs / owners` `UT-157.a` and `UT-157.c` are resolved in
  `TODO-kugod` by `DI-joroh`, with the candidate-envelope inventory and
  `grid([pcid, payload])` working-hypothesis prose captured in
  `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`.
  `UT-157.b` is retired under `TODO-rivuk` / `DI-runuh`.
- `Owner/doc cleanup` Done. The abandoned Reading-1/2/3 question is retired, the
  five candidate envelopes are captured in the grid-envelope successor owner,
  and the "working but not yet proven" grid framing is recorded as a candidate
  hypothesis rather than a canonical harness rule. No correction note is needed
  in `TODO-lilar`.
- `Remaining decisions or work` None for turn 157 itself. The later
  grid-envelope protocol directory/spec work that had remained under turn 158 /
  `UT-158.f` and `tujad.3` is now closed for turn-158 scope by `DI-fanah`.
- `Work pending` `no`
- `Proposed disposition` `reconciled after successor-owner capture`
- `Write needed? yes/no` `yes` for this rewalk update in `TODO-juhub`; `no`
  correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 158 is next and remains
  pending approval.

### Turn 158 — 2026-05-03 01:46 UTC

- `Turn 158 summary` This is the real apparatus-vs-specimen break. Steve
  pointed out that even calling the envelope "harness-wide" smuggles in the
  wrong assumption, because wire-lab exists to test multiple hypotheses at all
  layers, not to bake one answer into the harness. In plain English, this is
  the turn where the bot finally accepts that the harness-spec is lab
  apparatus and candidate envelopes are specimens under study.
- `Existing capture` `TODO-lilar` already captures the turn correctly as the
  foundational apparatus-vs-specimen reframe. The raw turn lays out a six-step
  sequence: audit the harness-spec, file the harness-level TE on the split,
  give each candidate envelope its own protocol home, sweep specimen material
  out of the harness-spec, reframe the old promise-stack TODO under protocol
  ownership, and file a parallel TODO for the grid-envelope hypothesis.
- `Gaps or contradictions` The insight itself stands; no later artifact reverts
  to the old claim that the harness-spec should define a single envelope
  specimen. The gap was bookkeeping: the six-step sequence and the parallel
  grid-hypothesis TODO were still represented as loose carry items even though
  later artifacts had decomposed, routed, or materialized them as successor
  simulations.
- `Related UTs / owners` `UT-158.a` is resolved-decomposed in `TODO-kugod`
  under `DI-sotuk`; `UT-158.b` is resolved by TE-havib DF-36.1; `UT-158.e`
  and `UT-158.g` are retired under `TODO-rivuk` / `DI-runuh`; `UT-158.h` is
  resolved-routed to `simulations/SIM-kurim-grid-envelope/TODO/TODO-tujad-grid-envelope-successor-owner.md`
  under `DI-sotuk`. `UT-158.c` and `UT-158.d` are resolved for turn-158 scope
  under `DI-kinad`; `UT-158.f` is resolved-transferred under `DI-fanah`.
- `Owner/doc cleanup` Done. `TODO-kugod` now has explicit `UT-158.a`,
  `UT-158.c`, `UT-158.d`, `UT-158.f`, and `UT-158.h` disposition rows, and
  `TODO-tujad` now closes `tujad.3` by pointing to the 24 standalone
  positional grid-envelope successor simulations.
  No correction note is needed in `TODO-lilar`.
- `Remaining decisions or work` None for turn 158 itself. Broader TE-40 audit
  work moved to later turn-159 rows and is now closed under `DI-mugar`.
- `Work pending` `no`
- `Proposed disposition` `reconciled after positional variant split`
- `Write needed? yes/no` `yes` for this rewalk update in `TODO-juhub` and the
  owner-routing updates in `TODO-kugod` / `TODO-tujad`; `no` correction note is
  needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 159 is next and remains
  pending approval.

### Turn 159 — 2026-05-03 01:50 UTC

- `Turn 159 summary` Steve confirmed the apparatus-vs-specimen reframe and
  told the bot to proceed with step 1 of the six-step plan. In plain English,
  this is the turn where the bot turns the reframe into a concrete inventory:
  it audits the harness-spec section by section and sorts what stays in the
  harness, what must move out as specimen-specific material, and what remains
  ambiguous.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The
  raw turn shows the bot produced an audit memo as an untracked working
  document, not yet a commit, and that memo listed the concrete specimen-
  bearing material in the harness-spec plus the ambiguous areas needing later
  TE or DF work. Later residual owner work in `TODO-kugod` carried this turn
  forward as `UT-159.a`, `UT-159.b`, and `UT-159.d`; those rows are now closed
  under `DI-mugar`.
- `Gaps or contradictions` None found. The only important timing nuance is
  that the audit memo existed in the workspace at the end of turn 159 but was
  not committed until Steve's turn-160 authorization. `TODO-lilar` already
  preserves that nuance, and later artifacts align with it.
- `Related UTs / owners` `UT-159.a` is resolved by the `DI-mugar`
  harness-spec sweep; `UT-159.b` is resolved by the `DI-huzor` feed-outer
  extraction plus the 2026-05-12 transport companion audit; `UT-159.c` is
  resolved-retired by the TE-havib follow-on verification walk; `UT-159.d` is
  resolved by the `DI-mugar` treatment of §1.3, §10, and §10a.
- `Owner/doc cleanup` Done. `TODO-kugod` closes `kugod.5`, `kugod.8`, and
  `kugod.9`; `TODO-lilok` no longer needs a separate harness-spec-sweep
  handoff; `TODO-jivam` and the UT verification matrix now point at this
  closure. No correction note is needed in `TODO-lilar`.
- `Remaining decisions or work` None for turn 159 itself. Later turns may
  still carry their own UTs, but the turn-159 audit/sweep residue is closed.
- `Work pending` `no`
- `Proposed disposition` `reconciled after apparatus/specimen sweep`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 160 is next and remains
  pending approval.

### Turn 160 — 2026-05-03 02:24 UTC

- `Turn 160 summary` Steve's `do it and continue` authorized two separate
  actions: first, commit the audit memo from turn 159, and second, continue
  into step 2 by drafting the apparatus-vs-specimen TE. In plain English, this
  is the turn where the work stops being just an audit and becomes a formal DF
  program on a dedicated TE twig.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  turn says the audit memo was committed at `4725b3e`, the new TE-havib work
  was drafted on twig `ppx/te-20260503-022446-apparatus-vs-specimen`, and the
  TE already contained seven drafted DFs plus six tabletop scenarios. Just as
  importantly, the answer text presented only `DF-36.1` to Steve in this turn,
  following the standing one-DF-at-a-time rule; the other six drafted DFs
  existed in the TE file but were not yet exposed in conversation.
- `Gaps or contradictions` None that overturn the existing capture. The later
  carry-forward items already preserve the important weaknesses introduced here:
  the audit count in the answer text said eight specimen-bearing items while
  the audit itself listed nine; the PT vocabulary collapse was baked into
  `DF-36.4` rather than framed as its own decision; the TE's six scenarios only
  partially align with the audit's recommended scenario set; and most of the
  seven drafted DFs remained unlocked at end-of-corpus.
- `Related UTs / owners` `UT-160.a` is resolved by the 2026-05-12 Cat-3
  refinement in TE-havib plus `TODO-kugod` / `DI-mugar`, which confirms the
  nine audit items were swept or explicitly retired. `UT-160.b` is answered by
  the TODO-lilok verification walk as procedural-meta, not a live DF split.
  `UT-160.c` is answered by the same verification walk as wrong on inspection.
  `UT-160.d` is resolved by the final TE-havib state: all seven DFs are locked
  after the Alt-B re-presentation path.
- `Owner/doc cleanup` Done. TE-havib now has a Cat-3 refinement that preserves
  the historical "eight" wording but points readers at the nine-item audit
  count and `DI-mugar` closure. `TODO-lilok` is closed; `TODO-lilar` remains
  append-only historical evidence and does not need a correction note.
- `Remaining decisions or work` None for turn 160 itself. Later turns still
  carry their own UTs, but the turn-160 count mismatch, PT/tabletop concerns,
  and end-of-corpus DF-lock concern are reconciled.
- `Work pending` `no`
- `Proposed disposition` `reconciled after TE-havib closure`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 161 is next and remains
  pending approval.

### Turn 161 — 2026-05-03 02:57 UTC

- `Turn 161 summary` Steve asked the key redundancy question: what is the
  actual difference between `promise-stack` and `grid-pcid-payload`? In plain
  English, this is where the bot is forced to explain whether these are truly
  two peer envelope hypotheses or whether one is really just a special case of
  the other.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  turn contains the most detailed side-by-side comparison in this whole replay
  slice: `promise-stack` is presented as a recursive CBOR array of promise
  frames whose trust and layering semantics live in the envelope itself, while
  `grid-pcid-payload` is presented as a thin dispatch envelope where the pCID
  is the top-level selector and the payload shape is left to the protocol it
  names. The most important point is the asymmetry the bot exposed: a
  `grid-pcid-payload` message can carry a promise-stack inside its payload, but
  a promise-stack message does not cleanly host `grid-pcid-payload` as a peer
  outer envelope. That asymmetry is the conceptual seed of `OQ-36.6`.
- `Gaps or contradictions` None that overturn the existing capture. The two
  main carry-forwards from this turn are now reconciled: the asymmetry is
  resolved by TE-havib DF-36.2's retirement of promise-stack as a separate
  envelope hypothesis, and the nine-axis comparison plus richer open-set
  assertion taxonomy are captured as historical inputs in the TE-havib
  2026-05-12 refinement.
- `Related UTs / owners` `UT-161.a` is answered by TE-havib DF-36.2 and the
  TODO-lilok verification path: promise-stack is retired as a separate envelope
  hypothesis, so the asymmetry concern is moot as a live OQ. `UT-161.b` and
  `UT-161.c` are captured by the 2026-05-12 TE-havib refinement as historical
  inputs rather than live downstream requirements.
- `Owner/doc cleanup` Done. TE-havib now records the nine-axis comparison's
  load-bearing result and the assertion-taxonomy examples. No correction note
  is needed in `TODO-lilar`.
- `Remaining decisions or work` None for turn 161 itself. Future assertion
  taxonomy work can still be filed if a later protocol needs it, but that would
  be new work rather than a turn-161 loose end.
- `Work pending` `no`
- `Proposed disposition` `reconciled after TE-havib refinement`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 162 is next and remains
  pending approval.

### Turn 162 — 2026-05-03 03:12 UTC

- `Turn 162 summary` Steve sharpened the redundancy concern from turn 161 into
  an explicit suspicion: promise-stack may be overcomplicated machinery
  invented from a misunderstanding of how nested messages already work inside
  `grid-pcid-payload`. He then gave two procedural instructions: note that
  concern for later, and keep going on `DF-36.5`. In plain English, this is
  the turn where the redundancy issue stops being just an implication and
  becomes an explicitly parked open question.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  turn's visible reply is minimal, but the substantive effect is preserved in
  later artifacts: the concern is recorded on the TE-havib twig as `OQ-36.6`,
  and the promise-stack-home path is later re-presented under Alt-B rather than
  left as a cleanly settled direction. The final TE-havib lock resolves the
  suspicion in the negative by retiring `promise-stack` as a separate protocol
  hypothesis instead of treating it as the active specimen home.
- `Gaps or contradictions` None remaining for turn 162. The two former
  residual concerns are now reconciled: `OQ-36.6` is visibly resolved in
  TE-havib, and `DF-36.2` was re-presented under Alt-B as Alt-2.A revised
  rather than left as a provisional promise-stack-home decision.
- `Related UTs / owners` `UT-162.a` is resolved by TE-havib's `OQ-36.6`
  negative-resolution path and the TODO-lilok verification walk. `UT-162.b` is
  resolved by the Alt-B re-presentation of `DF-36.2` and the final Alt-2.A
  revised lock.
- `Owner/doc cleanup` Done. TE-havib already carries the final `OQ-36.6` and
  `DF-36.2` resolution text; TODO-lilok is closed; the verification matrix now
  has a turn-162 closure pointer.
- `Remaining decisions or work` None for turn 162 itself.
- `Work pending` `no`
- `Proposed disposition` `reconciled after TE-havib final lock`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 163 is next and remains
  pending approval.

### Turn 163 — 2026-05-03 03:34 UTC

- `Turn 163 summary` Steve rejected the prior presentation of `DF-36.5` as
  unreadable and told the bot to format it better. In plain English, this is
  not just a meta-turn about presentation; it is the first readable,
  substantive walk of `DF-36.5`.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  turn re-presents the real decision about `§1.3` of `harness-spec`: whether
  the four simulator tests written in promise-stack vocabulary should stay in
  the harness envelope-agnostically, move wholesale to a specimen spec, or be
  split across the two abstraction levels. The bot recommends `Alt-5.C`: keep
  an apparatus-level summary in the harness and move specimen-specific details
  out. This is also the first DF in TE-havib where the parked `OQ-36.6`
  uncertainty is built explicitly into the recommendation: the apparatus-level
  summary survives either way, but the specimen-side destination depends on
  whether promise-stack later survives as a distinct specimen.
- `Gaps or contradictions` None remaining for turn 163. The turn's
  envelope-agnostic `§1.3` template no longer lives only in conversation:
  active `harness-spec-draft.md` now carries an apparatus-level
  layering-scenarios section under `DI-lajod`, later completed for the wider
  turn-159 audit by `DI-mugar`. The uncaptured-commit concern is also no longer
  a turn-local work item: its future-process lesson is captured in
  `AGENTS-ppx.md` B1, while any commit-specific residue is carried by the
  separately named UT rows.
- `Related UTs / owners` `UT-163.a` is resolved by the `§1.3` apparatus-level
  rewrite in `protocols/wire-lab.d/specs/harness-spec-draft.md`. `UT-163.b` is
  closed as a procedural rule by `AGENTS-ppx.md` B1 / `DI-021-20260507-212249`.
  `UT-cbf7f41-fallback` is retired for turn-163 purposes by OQ-36.6's negative
  resolution, DF-36.2's promise-stack retirement, and the active `§1.3`
  apparatus rewrite.
- `Owner/doc cleanup` Done. The verification matrix now has a turn-163 closure
  pointer. No `TODO-lilar` checkbox is flipped; the matrix remains the closure
  index.
- `Remaining decisions or work` None for turn 163 itself. Commit-specific rows
  not owned by this turn remain under their own UT identifiers if later replay
  reaches them.
- `Work pending` `no`
- `Proposed disposition` `reconciled after §1.3 rewrite and B1 transfer`
- `Write needed? yes/no` `yes` for this rewalk entry and checklist closure in
  `TODO-juhub`; `no` correction note is needed in `TODO-lilar`.
- `Next: wait for approval before turn N+1` Turn 164 is next and remains
  pending approval.

### Turn 164 — 2026-05-03 11:20 UTC

- `Turn 164 summary` This is the hard pivot away from the TE-havib DF walk and
  toward an urgent operational problem: Steve needs file-based transport
  working so he and another human collaborator can use the repo itself to
  collaborate. In plain English, the raw turn does not implement the bootstrap
  yet; it marks the urgency boundary and the point where the bot pauses to
  confirm scope before switching threads.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer is a pause, not a solution: the bot says it should not assume who the
  collaborator is and explicitly notes that pivoting away from TE-havib is a
  consequential choice. Later artifacts correctly treat this as the major
  thread boundary where TE-havib becomes parked and the transport-bootstrap
  work begins immediately afterward.
- `Gaps or contradictions` None that overturn the existing capture. The key
  later carry-forwards are now split by owner and closed: the "4-DF-locked"
  overclaim and TE-numan revert/merge discrepancy are historical corrections;
  the group-session round-trip coverage issue was later closed by the
  four-message `TODO-bisur` 012.7 artifact; membership-by-posting is closed by
  `DI-rurab`; freeze-rewrite is rejected by `DI-bomud`.
- `Related UTs / owners` `UT-164.a` is a historical correction against the
  current TE-havib status. `UT-164.b` is a historical correction; TE-numan
  scope was later settled by the TE-39 / TODO-lilar cross-cutting disposition.
  `UT-164.c` is resolved by `TODO-bisur` 012.7. `UT-164.d` is closed by
  sim-local `TODO-gapab` / `DI-rurab`. `UT-164.e` is closed by
  `TODO-gapab` / `DI-rurab` and `TODO-kakaz` / `DI-bomud`; the no-rewrite
  freeze boundary is now explicit in the feed-outer and wire-lab-devs docs.
- `Owner/doc cleanup` Done. The verification matrix has a turn-164 closure
  pointer. Rooted `TODO-turog` and `TODO-duvuk` remain historical coordination
  memory; active specimen-owned successor work is closed in the sim-local
  TODOs.
- `Remaining decisions or work` None for turn 164.
- `Work pending` `no`
- `Proposed disposition` `reconciled after successor closure`
- `Write needed? yes/no` `no`
- `Next` Turn 165 processed by the next chronological replay step.

### Turn 165 — 2026-05-03 11:21 UTC

- `Turn 165 summary` Steve disclosed that the collaborator is a real human
  using Claude in another clone and imposed a privacy constraint: the
  collaborator's name must not appear in committed docs. In plain English, this
  is the turn where the transport-bootstrap problem becomes a real coordination
  problem with an explicit anonymity rule attached to it.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer stacks four distinct moves: it records the collaborator context as a
  side effect, corrects the bot's architectural misread by admitting that a new
  git-file-transport TE was unnecessary because the existing transport and
  group-session drafts already cover the use case, claims the mistaken TE twig
  was "reverted cleanly," and then lays out the concrete bootstrap
  chicken-and-egg with three options and a recommendation for the draft-
  directory approach (`Option A`).
- `Gaps or contradictions` Captured and now resolved for turn-165 cleanup
  purposes. The "reverted cleanly" claim remains owned by `UT-164.b`, already
  reconciled as historical evidence. `UT-165.a` is closed as an observational
  privacy/slug lesson because the active specimen uses the generic
  `wire-lab-devs` slug and current group-session examples use neutral
  Alice/Bob prose. `UT-165.b` is closed by `DI-rurab`, which defines the
  interim `merge-group-transport-spec` shape as a Steve-authored DI until
  cryptographic promise tooling exists. `UT-165.c` was already closed by the
  neutral-memory update. `UT-165.d` is closed without a spec edit: OQ-G4 remains
  deferred, and m000 is one valid first-message pattern rather than a v0
  genesis-message mandate. `UT-165.e` is closed by the group-session spec
  cleanup that removed `codex-perplexity` examples and names the
  wire-lab-devs specimen as the freeze-gate evidence.
- `Related UTs / owners` `UT-165.a`, `UT-165.b`, `UT-165.d`, and
  `UT-165.e` are now checked off in `TODO-lilar`; `UT-165.c` was already
  checked off. No active owner TODO remains for turn 165.
- `Owner/doc cleanup` Done. `TODO-lilar` UT rows are closed; `TODO-juhub`
  carries the stronger turn note; `group-session-draft.md` and
  `wire-lab-devs-draft/README.md` already carry the resulting active wording
  from the `DI-rurab` / `DI-bomud` cleanup.
- `Remaining decisions or work` None for turn 165.
- `Work pending` `no`
- `Proposed disposition` `reconciled after lilar UT closure`
- `Write needed? yes/no` `no`
- `Next` Turn 166 is next.

### Turn 166 — 2026-05-03 11:29 UTC

- `Turn 166 summary` Steve corrected two things at once: the transport group is
  at least three developer agents, not a two-party collaboration, and the slug
  must be generic and identity-free rather than derived from people-names. In
  plain English, this is the turn where the bootstrap stops being a hypothetical
  plan and becomes the executed `wire-lab-devs` instance under urgency.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer updates the group model, proposes generic slug candidates, chooses
  `wire-lab-devs`, raises the membership-pinning issue from `§8`, and then
  reports the bootstrap as already live and pushed. Later timeline evidence
  confirms that the bootstrap commit and merge landed as part of this same
  response, not before it, so this is a high-effects execution turn rather than
  just a naming correction.
- `Gaps or contradictions` Captured and now resolved for turn-166 cleanup
  purposes. `UT-166.a` is closed for future-process purposes by the current
  decision-first protocol plus `DI-vanak`'s explicit replay approval shorthand;
  historical bootstrap commits remain preserved and are not rewritten.
  `UT-166.b` is closed by `DI-rurab`: active membership is the fixed configured
  set of exact `<author-id>/main` branches, so guessed actors are not enrolled
  by speculation. `UT-166.c` is closed by active specimen docs that use
  `stevegt-via-perplexity` as the committed `From:` identity. `UT-166.d` is
  closed as historical git metadata; rewriting the stale twig name would be a
  history rewrite. `UT-166.e` is closed by `DI-rurab`, which supersedes
  membership-by-posting with fixed configured branch membership, passive
  observer non-membership, and no self-admission from unknown branches.
- `Related UTs / owners` `UT-166.a`, `UT-166.b`, `UT-166.c`, `UT-166.d`, and
  `UT-166.e` are now checked off in `TODO-lilar`. Later turn rows that cite
  `UT-166.a` remain future-turn work only where their own turn-specific
  execution pattern still needs reconciliation.
- `Owner/doc cleanup` Done. `TODO-lilar` turn-166 UT rows are closed;
  stale future-row references to `UT-166.a` as "pending" were updated to point
  at the resolved process baseline; active group-session and wire-lab-devs
  docs already carry the `DI-rurab` membership and identity-safe wording.
- `Remaining decisions or work` None for turn 166.
- `Work pending` `no`
- `Proposed disposition` `reconciled after lilar UT closure`
- `Write needed? yes/no` `no`
- `Next` Turn 167 is next.

### Turn 167 — 2026-05-03 11:37 UTC

- `Turn 167 summary` Steve gave two directives at once: switch message files
  from `.msg` to `.txt`, and require members to fetch all branches but post
  only on their own `<author-id>/main`. In plain English, this is both a
  presentational cleanup turn and the first committed transport-binding turn
  that maps the abstract group-session protocol onto an actual Git branch
  discipline.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer shows that the bot treated both directives as immediately executable:
  it rewrote the docs, renamed the bootstrap message file without changing its
  bytes, added new non-normative `§9` to describe the per-author-branch
  binding, and authored `m001` as the first on-wire ratification of that
  branch-binding rule. This is the turn where the transport instance becomes
  "in flight" with an explicit Git discipline: read from all known branches,
  write only to your own branch, and let `Parents:` rather than branch topology
  carry ordering.
- `Gaps or contradictions` Captured and now resolved for turn-167 cleanup
  purposes. `UT-167.a` is closed because current active docs use filename = CID
  and no active sequential `m<N>-...` filename rule remains. `UT-167.b` is
  closed by `DI-rurab`, which rejects membership-by-posting for active wording
  and uses fixed configured `<author-id>/main` membership. `UT-167.c` is
  closed by the same `DI-rurab` wording: `{name}` now means author-id, not
  group name. `UT-167.d` is closed by the 2026-05-14 corpus audit for
  `\.msg\b`; remaining matches are historical replay/disposition notes or
  message-body evidence, not active spec guidance. `UT-167.e` is closed for
  future-process purposes by the current decision-first protocol plus
  `DI-vanak`; later execute-on-directive rows keep their own turn-specific
  reconciliation work.
- `Related UTs / owners` `UT-167.a`, `UT-167.b`, `UT-167.c`, `UT-167.d`, and
  `UT-167.e` are now checked off in `TODO-lilar`. The active group-session
  and wire-lab-devs docs already carry the filename, branch-membership, and
  passive-observer wording needed for this turn.
- `Owner/doc cleanup` Done. `TODO-lilar` turn-167 UT rows are closed;
  `TODO-juhub` carries this stronger note; the verification matrix has a
  turn-167 closure pointer. No active `.msg` spec-edit target was found.
- `Remaining decisions or work` None for turn 167.
- `Work pending` `no`
- `Proposed disposition` `reconciled after lilar UT closure and .msg audit`
- `Write needed? yes/no` `no`
- `Next` Turn 168 is next.

### Turn 168 — 2026-05-03 11:45 UTC

- `Turn 168 summary` Steve corrected the new branch-binding model from turn
  167 in two important ways: message files must not use sequential numbers, and
  before posting, each agent must first merge all observed messages from all
  branches into the directory on their own branch and push that merged state.
  In plain English, each participant's branch becomes "my replicated view of
  the message set, plus optionally a new post," and the message CID becomes the
  stable on-disk identity.
- `Existing capture` `TODO-lilar` captures the turn correctly. The bot treated
  both directives as executable: it rewrote filename rules to CID filenames,
  expanded §9 into a receive/merge/push/optionally-post cycle, renamed the
  existing bootstrap files to then-current CIDs, and authored the m2 on-wire
  ratification message. The turn fixes the global-sequence mistake from turn
  167 and defines the first actual replication model for the Git-backed
  specimen.
- `Gaps or contradictions` Captured and now resolved for active docs. The
  `Message-ID:` retention is closed by the active §4.3/§4.7 compatibility
  rule, §9's old non-normative status is closed by `DI-rurab`, passive-reader
  ambiguity is closed by current §8/§9.3, infrastructure/message-file
  boundaries are explicit enough for the current specimen, and stale turn-168
  CIDs are historical-only after turn 169's rehash.
- `Related UTs / owners` `UT-168.a` through `UT-168.f` are checked off in
  `TODO-lilar`. The only deeper rehash-continuity questions continue under the
  turn-169 rows, not as turn-168 work.
- `Owner/doc cleanup` Done. Updated the active `TODO-bisur` 012.7 note so it no
  longer calls §9 explicitly non-normative or keeps a stale open §9 OQ. No
  transport message files were changed.
- `Remaining decisions or work` None for turn 168.
- `Work pending` no.
- `Proposed disposition` `reconciled after lilar UT closure and active-doc CID/status audit`
- `Write needed? yes/no` `no` further turn-168 write is needed after this pass.
- `Next` Turn 169 is next.

### Turn 169 — 2026-05-03 11:54 UTC

- `Turn 169 summary` Steve asked whether `Message-ID:` is still needed now that
  the canonical identifier is the message CID. In plain English, the raw turn is
  a careful reasoning memo: the bot audits every plausible use of `Message-ID:`,
  concludes that it creates competing identity once filename = CID, and then
  reasons about compatibility for the already-authored bootstrap messages.
- `Existing capture` `TODO-lilar` captures the turn correctly. The raw reasoning
  recommended the conservative Path A (deprecate / preserve legacy bytes / ask
  Steve), while later repo history executed Path B (hard remove / rehash). The
  later active policy is no longer the unqualified Path-B-only state: current
  group-session §4.3 / §4.7, under `DI-012-20260508-033513`, says canonical
  writers MUST NOT emit `Message-ID:`, but readers MAY tolerate exactly one
  legacy pre-`Date:` header and MUST ignore its value semantically.
- `Gaps or contradictions` Captured and now resolved for active docs. The
  reasoning/action inversion is explicitly recorded; the strict-reader problem
  has an explicit legacy-header carve-out; the writer-side prohibition versus
  reader-side deprecation split is now deliberate; and the twig-name issue is
  historical metadata because active branch rules require only `ppx/{twig}` with
  a short kebab-case task phrase.
- `Related UTs / owners` `UT-169.a` through `UT-169.e` are checked off in
  `TODO-lilar`. `TODO-duvuk` is closed as historical coordination memory and no
  longer owns active execution for these Message-ID / CID-cascade rows.
- `Owner/doc cleanup` Done. Updated `TODO-duvuk` so its original
  T-FILENAME-CID-CASCADE scope is visibly historical and updated `TODO-jivam`
  so its live TE-42 closure condition no longer appears open. Historical quoted monitor snapshots remain unchanged as evidence.
- `Remaining decisions or work` None for turn 169. Later nested-body and
  grid-envelope cascade questions remain with their own later-turn rows, not as
  turn-169 work.
- `Work pending` no.
- `Proposed disposition` `reconciled after lilar UT closure and active Message-ID compatibility audit`
- `Write needed? yes/no` `no` further turn-169 write is needed after this pass.
- `Next` Turn 170 is next.

### Turn 170 — 2026-05-03 16:53 UTC

- `Turn 170 summary` Steve shifted back from execution to design review and
  asked whether the flat `transports/draft--wire-lab-devs/` layout should gain
  a protocol/grouping layer, especially if a second named group appears. In
  plain English, this opened the directory-axis question that later turns
  broadened into protocol identity, feed/substrate, site, CAS, and simulation
  placement.
- `Existing capture` `TODO-lilar` captures the turn correctly. The raw answer
  framed `DF-37.1` with three alternatives: protocol-slug nesting, recursive
  draft/pCID nesting, or status quo deferral. It recommended protocol-slug
  nesting, but no implementation happened and Steve did not answer that DF as
  posed.
- `Gaps or contradictions` Captured and now resolved. `DF-37.1` is not still an
  actionable open DF: turns 171-176 reframed the missing axis, TE-sihih / TODO-
  vunub landed the L5/L6/L7 substrate-agnostic model, TE-domat and DR-nugog
  reframed the root `transports/` / `groups/` question, and `DI-fakin` resolved
  the current specimen by moving it into a simulation world instead of choosing
  any of the original root flat/nested alternatives.
- `Related UTs / owners` `UT-170.a` through `UT-170.d` are checked off in
  `TODO-lilar`. `TODO-vunub` is now marked closed because TE-sihih is decided;
  DR-nugog is implemented for the current specimen and has an append-only note
  naming the current `SIM-ludut-wire-lab-devs` path after the later rusis split.
- `Owner/doc cleanup` Done. Updated TODO-vunub status, the master TODO row, the
  verification matrix owner/closure notes, and DR-nugog current-path provenance.
- `Remaining decisions or work` None for turn 170. Later root/reference-layout
  questions, if any, are downstream graduation questions rather than turn-170
  recovery loose ends.
- `Work pending` no.
- `Proposed disposition` `reconciled after TE-sihih / DR-nugog / simulation-path audit`
- `Write needed? yes/no` `no` further turn-170 write is needed after this pass.
- `Next` Turn 171 is next.

### Turn 171 — 2026-05-03 16:56 UTC

- `Turn 171 summary` Steve refined the turn-170 tree question from "should
  there be protocol grouping?" to "should there also be a separate path layer
  meaning git file transfer?" In plain English, this is where the replay first
  separates protocol identity from delivery substrate, before turn 172 expands
  the substrate axis beyond git.
- `Existing capture` `TODO-lilar` captures the turn correctly. The raw answer
  rejects a `git/` path layer for the current one-substrate case, keeps the
  protocol-slug recommendation from turn 170, and says substrate facts should
  become another axis or per-instance metadata if they later need first-class
  representation.
- `Gaps or contradictions` Captured and now resolved. The "§9 captures git"
  comfort is superseded by TE-sihih's L5 feed model, but the active
  group-session specimen is also intentionally allowed to keep §9 inline as the
  normative wire-lab-devs git binding under `DI-rurab`. The manifest-field hook
  is closed by TODO-vunub Q-22.3, which retracts manifest schema work in favor
  of TE-vipir path-as-declaration.
- `Related UTs / owners` `UT-171.a` through `UT-171.d` are checked off in
  `TODO-lilar`. `TODO-vunub` records the retired turn-171 substrate/manifest
  framing; cadence-only rows are preserved as lessons, not live implementation
  owners.
- `Owner/doc cleanup` Done. Updated TODO-vunub's retired-question and DI-log
  notes. No active spec or transport-message bytes were changed.
- `Remaining decisions or work` None for turn 171. Later feed-vocabulary,
  groups/transports, and nested-envelope questions remain with their own later
  turns, not with turn 171.
- `Work pending` no.
- `Proposed disposition` `reconciled after TE-sihih / TODO-vunub / DI-rurab audit`
- `Write needed? yes/no` `no` further turn-171 write is needed after this pass.
- `Next` Turn 172 is next.

### Turn 172 — 2026-05-03 17:01 UTC

- `Turn 172 summary` Steve blew up the narrow `git` question by listing
  multiple peer substrates: `rsync`, `unison`, `uucp`, `udp`, `svn`, `cvs`,
  and `git`. In plain English, this is the turn where the replay stops
  treating delivery substrate as a side detail and starts treating it as a
  first-class design axis.
- `Existing capture` `TODO-lilar` captures the turn correctly. The raw answer
  proposes a new TE/DF program because if one `group-session` can move over
  multiple byte-moving substrates, then git-specific rules cannot be treated as
  timeless core group semantics. The answer's specific words and layout are
  transitional: it says `binding`, sketches `bindings/` and `messages/`
  subdirectories, and imagines extracting `§9` into a separate git-specific
  spec.
- `Gaps or contradictions` Captured and now resolved. The load-bearing
  insight survives, but the exact turn-172 taxonomy does not. TE-sihih replaces
  `binding` with L5 `feed`, adds L6 CAS and L7 group as the citable layer
  split, retracts per-instance feed manifests in favor of TE-vipir
  path-as-declaration, and leaves current wire-lab-devs git rules inline in
  group-session `§9` under `DI-rurab`.
- `Related UTs / owners` `UT-172.a` through `UT-172.e` are checked off in
  `TODO-lilar`. `TODO-vunub` records the turn-172 proposal as retired by
  TE-sihih / Q-22.2 / Q-22.3 / Q-22.6, with future feed-spec extraction owned
  by successor work rather than this replay turn.
- `Owner/doc cleanup` Done. Updated `TODO-lilar`, `TODO-vunub`, this turn note,
  and the UT verification matrix. No active spec or transport-message bytes
  were changed.
- `Remaining decisions or work` None for turn 172. Later historical-analog,
  feed-vocabulary, CAS/site, and simulation-layout details belong to their own
  later turns and owner artifacts.
- `Work pending` no.
- `Proposed disposition` `reconciled after TE-sihih / TODO-vunub / DI-rurab audit`
- `Write needed? yes/no` `no` further turn-172 write is needed after this pass.
- `Next` Turn 173 is next.

### Turn 173 — 2026-05-03 17:06 UTC

- `Turn 173 plain-English recap` Steve did not answer the proposed TE directly.
  He asked whether the substrate-pluralism idea had precedent in practice,
  RFCs, and historical networks. The assistant answered with a precedent survey
  covering email over SMTP/UUCP/X.400, Usenet over NNTP/UUCP, FidoNet, CORBA
  GIOP/IIOP, SOAP/WSDL, modern pluggable-transport systems, and git itself. The
  turn's conclusion was that the broad concept is well precedented: message
  identity can stay stable while different substrates move the same bytes. The
  strongest plain-English insight was that `group-session` looks like a tiny
  content-addressed Usenet, but later turns modify the details: `binding` is
  rejected in favor of `feed`, the `bindings/` layout sketch is retracted, and
  the Usenet line is captured as exploratory design evidence rather than a
  frozen claim that PromiseGrid is Usenet. At the end of the turn, the open
  questions were naming, layout, negative counter-precedent, and git/CAS
  cardinality.
- `Existing capture` `TODO-lilar` captures the raw turn and its five residual
  rows. Later promotion work under `DI-pijun` puts the historical survey in
  `docs/research/historical-networks-20260503.md` and the broader design line in
  `simulations/SIM-hugoj-cas-usenetlike-gitlike/README.md`.
- `Gaps or contradictions` Captured and now resolved for turn-173 recovery. The
  exact Usenet analogy is now design-visible, the later vocabulary correction is
  recorded, and the negative-precedent gap has been narrowed by a bounded check
  in the research note. The analogy remains deliberately exploratory; future
  TEs may import or reject specific Usenet mechanisms one at a time.
- `Related UTs / owners` `UT-173.a` through `UT-173.e` are checked off in
  `TODO-lilar`. `TODO-vunub` records the turn-173 historical-precedent questions
  as retired by `DI-pijun`, the historical-networks note, `SIM-hugoj`, and
  TE-sihih's feed/CAS/layering decisions.
- `Owner/doc cleanup` Done. Updated `TODO-lilar`, `TODO-vunub`, the
  historical-networks research note, this turn note, and the UT verification
  matrix. No active spec or transport-message bytes were changed.
- `Remaining decisions or work` None for turn 173. Later feed-vocabulary,
  site/CAS, and simulation-layout decisions remain with their own later turns
  and owner artifacts.
- `Work pending` no.
- `Proposed disposition` `reconciled after DI-pijun / historical-networks / SIM-hugoj / TE-sihih audit`
- `Write needed? yes/no` `no` further turn-173 write is needed after this pass.
- `Next` Turn 174 is next.

### Turn 174 — 2026-05-03 17:13 UTC

- `Turn 174 plain-English recap` Steve did not accept the turn-173 proposal as
  ready to draft. He raised three objections at once: `binding` did not sound
  like real Usenet/email vocabulary, the idea seemed to conflict with OSI, and
  putting a substrate-spec file inside the messages directory felt inverted.
  The response accepted those objections: it rejected `binding` in favor of
  `feed`, reframed the relationship as alternative delivery/encapsulation
  rather than OSI-style vertical layering, and replaced the nested `bindings/`
  idea with an instance-root `INSTANCE.md` manifest. Later turns and owner
  artifacts keep the `feed` correction but modify the layout and layering
  conclusions: TE-sihih uses `feed` for the L5 protocol role and `substrate`
  for prose, TE-vipir / TODO-vunub Q-22.3 retract `INSTANCE.md` and
  `bindings/` feed declarations in favor of path-as-declaration, the temporary
  "horizontal encapsulation" wording is replaced by the L5/L6/L7 model, and
  the old `udp-binding` line is carried forward as the `udp-feed` simulation
  lineage. The loose ends at the end of the turn were naming, the
  `udp-binding` draft name, `INSTANCE.md` authority, the encapsulation term,
  and the Pattern-B mapping.
- `Existing capture` `TODO-lilar` already captures the turn correctly. The raw
  answer accepts all three objections and revises the proposal in three
  matching ways: `binding` is rejected in favor of `feed`; the relationship is
  reframed as alternative delivery/encapsulation rather than OSI-style vertical
  layering; and the nested `bindings/` idea is replaced with an instance-root
  `INSTANCE.md` manifest. This is also the turn where the bot explicitly
  admits that it imported `binding` from WSDL/CORBA/RPC lineage and that the
  more honest historical lineage for this problem is Usenet/FidoNet/email.
  Later residual notes correctly preserve the follow-on consequences: the
  historical-networks note stays accurate but needs vocabulary-aware reading,
  `udp-binding` now looks retroactively misnamed, `INSTANCE.md` may be
  overloading feed facts and governance facts, and the turn-173 Pattern-B map
  must change from `bindings/` to an instance-manifest shape.
- `Gaps or contradictions` None that overturn the existing capture. The main
  later limit is that two details changed after the turn: the
  encapsulation-not-layering framing was replaced by the L5/L6/L7 model, and
  the `INSTANCE.md` / `bindings/` mechanism was retracted by
  path-as-declaration. Those are later modifications to the turn's conclusions,
  not evidence that the raw turn was miscaptured.
- `Related UTs / owners` `UT-174.a`, `UT-174.b`, `UT-174.c`, `UT-174.d`, and
  `UT-174.e` are closed in `TODO-lilar`. Owner evidence now lives in
  TE-sihih, TE-vipir, `TODO-vunub` Q-22.3, `SIM-ludaf-udp-feed`, and the
  historical-networks research note.
- `Owner/doc cleanup` Done for this turn. The historical-networks note now
  carries the current-reading vocabulary and Pattern A/B/C mapping. TODO-vunub
  records that the turn-174 vocabulary/layout corrections are retired by the
  landed TE-sihih / TE-vipir combination. No active spec or transport-message
  bytes were changed.
- `Remaining decisions or work` None for turn 174. Later turn-175 and newer
  site/CAS/layering/procedural rows remain their own later-turn work.
- `Work pending` no.
- `Proposed disposition` `reconciled after TE-sihih / TE-vipir / TODO-vunub / historical-networks / udp-feed simulation audit`
- `Write needed? yes/no` `no` further turn-174 write is needed after this pass.
- `Next` Turn 175 is next.

### Turn 175 — 2026-05-03 17:25 UTC

- `Turn 175 plain-English recap` Steve delivered four corrections instead of
  answering the turn-174 questions: `grid <pcid>` should be understood as the
  grid envelope, not a carrier line; both `feed` and `substrate` were acceptable
  words; `transports/` looked misnamed; and the simulation was missing first-
  class sites and decentralized CAS. The response accepted the corrections and
  sketched a much larger model: grid envelope plus headers, forum/feed/site/CAS
  vocabulary, and a vertical layer stack where named message spaces sit above
  feeds and CAS. Later updates modify several conclusions: `forum` becomes
  forward path vocabulary `group`; TE-sihih lands as a smaller vocabulary /
  L5-L6-L7 / L6-CAS-pointer TE rather than the full 10-DF migration TE; historic
  message bytes are not rewritten for vocabulary sweeps; operational migration
  and the first concrete L6 CAS spec move to TODO-pipus and TODO-kituj / TE-43;
  root `transports/` versus `groups/` work moves through DR-nugog, TE-domat, and
  later simulation-structure TEs.
- `Existing capture` `TODO-lilar` captured the turn's core correctly: Steve's
  prompt is a locked correction batch, not an answer to the previous bot
  question, and the bot's response expanded TE-sihih scope by adding grid
  envelope vocabulary, sites, CAS, and a layered model.
- `Gaps or contradictions` The existing capture needed later-owner context.
  The response's predicted second CID cascade did not become the active plan;
  later freeze/CID decisions keep historical specimen bytes immutable. The
  response's `forums/` recommendation was replaced by `groups`. The response's
  broad TE-sihih scope was later contracted, with migration/CAS adoption split
  into successor owners. The turn also silently reversed turn 174's
  encapsulation-not-layering frame; the rewalk now records that reversal
  explicitly.
- `Related UTs / owners` `UT-175.a` through `UT-175.g` are checked off in
  `TODO-lilar` as resolved/transferred. Current owner evidence is split across
  TODO-vunub / TE-sihih, TODO-pipus, TODO-kituj / TE-43, TODO-turog /
  TODO-duvuk plus their sim-local successors TODO-gapab / TODO-kakaz, DR-nugog /
  TE-domat / later simulation-structure TEs, SIM-kurim grid-envelope variants,
  `docs/research/historical-networks-20260503.md`, and SIM-hugoj.
- `Owner/doc cleanup` Done for this pass. TODO-vunub now states the turn-175
  corrections are retired or transferred for that owner; TODO-pipus and
  TODO-kituj no longer claim they are blocked merely on TE-sihih landing; the
  historical-networks naming table now says `grid envelope`, not `carrier line`;
  the UT matrix has a turn-175 closure pointer. No transport message bytes were
  edited.
- `Remaining decisions or work` TODO-pipus remains open for the operational
  wire-lab-devs migration. TODO-kituj / TE-43 remains open for the first
  concrete L6 CAS / promisebase adoption decision. DR-nugog / TE-domat and later
  simulation-structure TEs still own root/simulation layout graduation. The
  grid-envelope simulations remain competing specimens rather than a chosen
  canonical envelope.
- `Work pending` yes.
- `Proposed disposition` `resolved/transferred after TE-sihih scope contraction, no-rewrite freeze/CID decisions, and simulation-owner split`
- `Write needed? yes/no` `no` further turn-175 write is needed after this pass.
- `Next` Turn 176 is next.

### Turn 176 — 2026-05-03 18:05 UTC

- `Turn 176 plain-English recap` Steve accepted some of the previous layering
  direction but corrected the details: PromiseGrid layers should not be numbered
  below ordinary networking layers; `groups` is better than `forums`; 1:1
  ephemeral flows still need to fit; PromiseGrid-over-Usenet and
  Usenet-over-PromiseGrid both matter; nested grid envelopes are expected;
  draft/frozen names should be `<slug>-draft` and `<slug>-<cid>`; CAS content
  might eventually use symlinks, Merkle roots, and Rabin chunking; and the broad
  strategic analogy is that PromiseGrid can become a content-addressed
  Usenet-plus-git successor. The response expanded TE-sihih again into a large
  multi-DF plan. Later updates narrow that expansion: TE-sihih lands as a
  smaller L5/L6/L7 vocabulary and L6-CAS-pointer decision, TE-domat owns the
  later root `groups/` versus `transports/` reconciliation, SIM-hugoj and the
  historical-networks note capture the Usenet/git analogy, and concrete CAS /
  migration work moves to TODO-kituj / TE-43 and TODO-pipus.
- `Existing capture` `TODO-lilar` captured the turn's substance well. It
  preserved the seven corrections, the symlink/CAS/chunking considerations, the
  Usenet-plus-git framing, the body-can-be-anything realization, and the bot's
  risky TE-sihih scope growth. TE-domat later also captured the turn's `groups`
  vocabulary, 1:1-flow test, nested-envelope allowance, slug-state naming, and
  mapping from the old transport path to a group-session path.
- `Gaps or contradictions` The capture needed later-owner reconciliation. The
  turn-176 response treated several ideas as TE-sihih scope, but TE-sihih later
  explicitly contracted to L5/L6/L7 vocabulary, L6 CAS subtree, and 100-year-goal
  citation work. The literal root rename from `transports/` to `groups/` is not
  settled; TE-domat rejects a naive whole-tree rename and keeps the decision open
  for DF. Nested-envelope strict-reader behavior is not settled by turn 176; it
  moved out of TODO-duvuk into later body-semantics / grid-envelope successors.
  The "foundational promise" wording also needs assessment-oriented Promise
  Theory phrasing in successor CAS work rather than a central or constitutive
  harness rule.
- `Related UTs / owners` `UT-176.a` through `UT-176.i` are checked off in
  `TODO-lilar` as resolved, closed, or transferred. Current owner evidence lives
  in TE-sihih / TODO-vunub for L5/L6/L7, slug-state, group vocabulary, and
  no-exception L6 CAS pointers; TE-domat / DR-nugog for the root `groups/` /
  `transports/` split; SIM-hugoj and `docs/research/historical-networks-20260503.md`
  for the content-addressed-Usenet/git analogy; TODO-pipus for operational
  migration; TODO-kituj / TE-43 for concrete L6 CAS adoption; and grid-envelope /
  body-semantics successor work for nested envelopes.
- `Owner/doc cleanup` Done for this pass. TODO-vunub now says turn-176 items are
  retired for its closed scope or transferred to successor owners. The UT matrix
  has a turn-176 closure / transfer pointer. No transport message bytes or TE
  bodies were edited.
- `Remaining decisions or work` TODO-pipus remains open for operational migration
  from pre-CAS inline specimens to CAS / pointer / feed shape. TODO-kituj / TE-43
  remains open for the first concrete L6 CAS / promisebase adoption decision.
  TE-domat / DR-nugog still own root layout and the `groups/` / `transports/`
  split. Nested-envelope body semantics and grid-envelope variant convergence
  remain successor simulation/spec work.
- `Work pending` yes.
- `Proposed disposition` `resolved/transferred after TE-sihih scope contraction, TE-domat root-layout analysis, and successor CAS/migration/body-semantics owner split`
- `Write needed? yes/no` `no` further turn-176 write is needed after this pass.
- `Next` Turn 177 is next.

### Turn 177 — 2026-05-03 18:34 UTC

- `Turn 177 plain-English recap` Steve corrected and expanded the layered model
  again. He preferred L5/L6/L7 names, asked for a concrete explanation of sites,
  made the promise economy foundational at every layer, made the 100-year goal
  explicit, corrected JSON to CBOR, asked whether feeds and CAS were inverted,
  preferred pointer files over filesystem symlinks, and asked whether Rabin
  chunking should be included now. The response did the right thing on the major
  structural question: it gamed out both layer orders and conceded the inversion.
  The corrected model became L7 groups define message semantics, L6 CAS stores
  and resolves chunks, and L5 feeds advertise/request/replicate chunks between
  sites. Later artifacts mostly keep that model: TE-sihih lands L5/L6/L7 and
  no-exception CBOR pointers into L6 CAS, TE-domat cites turn 177 as the reason a
  naive `transports/` to `groups/` rename is wrong, TODO-kituj / TE-43 keeps the
  concrete CBOR, chunking, CIDv1-codec, pointer-file, and promisebase questions,
  and TODO-kulih / TE-nibar is the better owner for any rule that every spec must
  expose promise vocabulary, 100-year pressure tests, or easy mental models.
- `Existing capture` `TODO-lilar` captured the turn's core correctly: the layer
  inversion was explicit, the promise economy and 100-year framing became
  load-bearing, CBOR displaced JSON, pointer files displaced symlinks, Rabin /
  FastCDC chunking entered active design pressure, and TE-sihih scope growth was
  becoming unmanageable.
- `Gaps or contradictions` The existing capture needed current-owner routing.
  The turn-177 response grew TE-sihih to 15 DFs and proposed stuffing migration,
  concrete CAS encoding, promise vocabulary, and easy mental-model obligations
  into the same TE. Later work contradicts that scope: TE-sihih lands smaller,
  TODO-kituj / TE-43 owns concrete L6 CAS choices, TODO-pipus owns operational
  migration, and TODO-kulih / TE-nibar is the spec-shape owner for per-spec
  promise-vocabulary / mental-model requirements. Turn 178 also refines
  `UT-177.h`: CIDv1 codec fields are the right object-type discriminator for
  raw chunks, Merkle nodes, and pointer objects.
- `Related UTs / owners` `UT-177.a` through `UT-177.i` are checked off in
  `TODO-lilar` as resolved, closed, or transferred. TE-sihih / TODO-vunub owns
  the layer-order resolution and scope contraction; TODO-kituj / TE-43 owns
  deterministic CBOR, tags, chunking parameters, CIDv1 codec/object typing, and
  promisebase prior-art adoption; TODO-pipus owns the eventual migration from
  pre-CAS inline specimens to pointer-and-CAS form; TODO-kulih / TE-nibar owns
  spec-doc-shape questions such as promise-vocabulary and easy mental-model
  sections; TE-dajot remains the citable 100-year-goal constraint; TE-domat /
  DR-nugog owns the root `groups/` / `transports/` split.
- `Owner/doc cleanup` Done for this pass. TODO-vunub now records that turn-177
  items are retired for its closed scope or transferred to successor owners.
  TODO-kituj now explicitly carries the turn-177 concrete CAS obligations.
  TODO-kulih now carries the turn-177 per-spec promise-vocabulary / 100-year /
  mental-model obligations. The UT matrix has a turn-177 closure / transfer
  pointer. A later simulation backfill added standalone charters for CAS object
  modeling, chunk-feed replication, CAS-backed group-session, and promise
  accounting records so the turn's conclusions can evolve inside simulations
  without contaminating the grid-envelope variants or rewriting historical
  message bytes. A second simulation pass added `SCENARIOS.md` matrices to make
  the turn-177 pressure concrete for successor owners while leaving decisions in
  TODO-kituj / TE-43, TODO-pipus, and TODO-kulih / TE-nibar. No transport
  message bytes or TE bodies were edited. A third cleanup pass filed scoped DRs:
  `DR-tumus` for concrete L6 CAS adoption, `DR-gabif` for additive CAS-backed
  group-session migration, and `DR-robon` for turn-177 spec-shape requirements.
  A fourth pass added unanswered next-DF packets and acceptance criteria to
  those DRs, then added owner TODO subtasks so the next cleanup does not need to
  re-derive the question shape. A fifth pass routed `DR-tumus` DF-tumus.1
  through DF-tumus.3 through three standalone bakeoff simulations for starting
  profile, CAS object type binding, and chunking identity after Steve asked for
  sims instead of direct answers. A sixth pass synthesized those bakeoff sims
  back into `DR-tumus` as a final answerable packet with recommended defaults,
  but did not answer or close the DR. Source: `DI-navod`; `DI-pator`;
  `DI-davov`; `DI-majib`; `DI-bukoh`; `DI-molah`.
- `Remaining decisions or work` TODO-kituj / TE-43 remains open for concrete L6
  CAS adoption: deterministic CBOR profile, allowed tags, chunking algorithm and
  parameters, CIDv1 codec/object typing, and promisebase prior-art stance.
  TODO-pipus remains open for operational migration to pointer-and-CAS shape.
  TODO-kulih / TE-nibar remains open for spec-doc structure decisions, now
  including promise vocabulary, 100-year pressure-test, and easy mental-model
  obligations. The new `DI-navod` simulations make these workstreams visible as
  specimen pressure, and `DI-pator` scenario matrices make them actionable, but
  the concrete TE/DF/DI decisions remain open through `DR-tumus`, `DR-gabif`,
  and `DR-robon`; `DI-majib` makes each DR answerable with explicit DF packets.
  `DI-bukoh` added the evidence step for `DR-tumus` DF-tumus.1 through
  DF-tumus.3, and `DI-molah` removes that sim-review blocker by synthesizing
  the bakeoffs into the current `DR-tumus` packet. For recovery-walkthrough
  purposes, these loose ends are captured in sim questions or TODO owners even
  though the downstream design work remains open. TE-domat / DR-nugog remains
  open for root layout.
- `Work pending` no.
- `Proposed disposition` `resolved/transferred after TE-sihih scope contraction, TE-43 CAS-owner routing, and TODO-kulih spec-shape routing`
- `Write needed? yes/no` `no` further turn-177 write is needed after this pass.
- `Next` Turn 178 is next.

### Turn 178 — 2026-05-03 19:14 UTC

- `Turn 178 plain-English recap` Steve gave the largest combined confirmation
  and expansion of the L5/L6/L7 discussion. He made sparse CAS a foundational
  assumption: no site should be assumed to have all objects, and every protocol
  or simulation must work when each site has only a subset. He pointed out that
  the word "decides" in "Bob's CAS decides it wants the chunk" is load-bearing:
  that is where consensus, promises, trust relationships, costs, refusal,
  policy, and incentives enter the system. He also questioned whether promise
  accounting belongs independently at each layer or mostly at L7 with L5/L6
  taking directions about what to pull, keep, and advertise. He floated a
  BGP-class future application: could a PromiseGrid app replace BGP's vulnerable
  inter-domain routing trust model? He confirmed that the capture-resistance
  narrative from turn 177 should lead the explanatory story: email, web, and
  social systems were captured partly because their upper layers lacked
  decentralized promise economies. He also confirmed the simple layperson model:
  sites make promises, keep them or do not, peers decide whom to trust, bad
  actors get cut off by peers, and no central authority is needed. He stated a
  design stance against pure anonymity as the default: group members know each
  other in some way, trust relationships exist, and this has implications for
  group identity. He asked whether to switch messages to CBOR and chunking now,
  under dogfood pressure from the developer/agent team. He surfaced promisebase
  as prior art for chunking, Merkle trees, streaming bytes into CAS, FUSE/container
  attempts, and reference/hash naming problems. He asked whether CIDv1 codec
  fields can distinguish Merkle roots from raw content hashes; the answer is the
  design direction now carried by TE-43 / DR-tumus and the CAS object-type sims:
  use codec / multicodec identity rather than filename suffixes. He required
  continued optional interop with libp2p, IPFS, and ATPROTO. Finally, he raised
  the multi-repo question: actual message content, especially large test files,
  may need to live in separate repos or one repo per simulated site rather than
  always in wire-lab. The assistant response was unusually terse: it recorded
  memory and asked only three meta-questions about reading promisebase, keeping
  TE-sihih together or split, and whether `pgmsg` was an acceptable tool name.
  Later work changes the routing: TE-sihih landed much smaller, so turn-178's
  broad material is now successor pressure for TODO-kituj / DR-tumus, TODO-pipus
  / DR-gabif, TODO-kulih / DR-robon, DR-napum, and the relevant simulations
  rather than a mandate to stuff everything into TE-sihih.
- `Existing capture` `TODO-lilar` captured the raw turn in detail and correctly
  flagged that the bot's response was a tactical deferral rather than an answer.
  Later artifacts now capture most of the substance: TODO-vunub records sparse
  CAS as a foundational invariant; `SIM-zazit` carries sparse chunk-feed and
  "decides" pressure; `SIM-rusap` carries peer-local promise accounting and
  cross-layer decision pressure; `SIM-punaz` carries the BGP-class routing app
  question; TODO-kituj / `DR-tumus`, `SIM-jomag`, `SIM-bobud`, `SIM-kohad`, and
  `SIM-gobaz` own promisebase prior art, interop, chunking, CIDv1 object typing,
  and pointer-object shape; `SIM-ligan` carries the separately-routed
  promisebase reference-naming question; `SIM-jurar` owns CAS-backed
  group-session identity pressure through an explicit question; TODO-rohub owns
  `pgmsg` tool-name and collaborator-permission meta-questions; DR-napum owns
  final public layperson claims.
- `Gaps or contradictions` The historical turn note was written before the
  later TE-sihih scope contraction, before the turn-177 simulation backfill, and
  before the protocol/specimen sim-question sweep. It therefore treated many
  questions as "TE-sihih must decide" even though later work deliberately routes
  them to TE-43 / DR-tumus, TODO-pipus / DR-gabif, TODO-kulih / DR-robon,
  DR-napum, or simulation questions. Review found that BGP-class routing app
  pressure and group identity / anti-anonymity pressure were visible in scenario
  rows but not in `QUESTION.md` homes, that promisebase reference naming was
  being conflated with CBOR / chunking / CID object typing, and that `UT-178.k`
  still needed an explicit `pgmsg` / collaborator-permission TODO owner. The
  next correction was that BGP-class routing and promisebase reference naming
  should not merely be questions inside broader sims; `DI-tibis` splits them
  into standalone simulations.
- `Related UTs / owners` `UT-178.a` is captured by TODO-vunub's sparse-CAS
  invariant and the sparse scenarios in `SIM-zazit` / `SIM-rusap`.
  `UT-178.b` and `UT-178.c` are captured by the `SIM-zazit` pull-decision
  scenario and the `SIM-rusap` cross-layer decision/accounting scenarios.
  `UT-178.d` is captured as BGP-class routing app pressure in standalone
  `SIM-punaz-bgp-class-routing-app`.
  `UT-178.e` is routed to DR-napum and DEV-GUIDE-RESOURCES rather than to the
  now-smaller TE-sihih. `UT-178.f` is captured by `SIM-jurar`'s explicit
  group-membership / identity question and DR-napum public-claims caution.
  `UT-178.g` is closed as dogfood pressure already reflected in TE-pahah /
  TE-nizor and the `SIM-ludut-wire-lab-devs` world, with concrete tool-name and
  collaborator-permission follow-up routed to TODO-rohub. `UT-178.h`,
  `UT-178.i`, and `UT-178.l` are split: CBOR, chunking, CIDv1 object typing,
  pointer-object shape, and interop are captured by TODO-kituj / DR-tumus and
  the CAS object-model / bakeoff sims, while promisebase human-readable
  reference-symbol / hash-name-resolution naming is separately routed to
  TODO-kituj `kituj.5` and standalone
  `SIM-ligan-promisebase-reference-naming`. `UT-178.j` is captured by the new
  multi-repo sparse-site scenario in `SIM-zazit`.
  `UT-178.k` is closed only after routing its meta-questions: promisebase review
  to TODO-kituj / DR-tumus and later promisebase work, TE-sihih scope to
  TODO-vunub and successor owners, and `pgmsg` / collaborator-permission
  follow-up to TODO-rohub. `UT-178.m` is closed as a caught-before-commit
  edit-defect lesson.
- `Owner/doc cleanup` Done. Added `DI-vaguf`; updated the three existing
  simulation surfaces that naturally own the turn-178 pressure; updated
  DR-napum and DEV-GUIDE-RESOURCES for the capture narrative / layperson claim
  routing; added a turn-178 closure pointer to the UT matrix; closed the
  `TODO-lilar` turn-178 UT rows with an additive closure summary, then repaired
  the review gaps under `DI-lusum` by adding explicit `QUESTION.md`, TODO, and
  reference-naming owner entries; then split BGP-class routing and promisebase
  reference naming into standalone sims under `DI-tibis`.
- `Remaining decisions or work` Downstream design remains open, but it is now
  containerized: `DR-tumus` / TODO-kituj for concrete L6 CAS adoption,
  `DR-gabif` / TODO-pipus for additive CAS-backed group-session migration,
  `DR-robon` / TODO-kulih for spec-section requirements, DR-napum for final
  public layperson claims, TODO-rohub for dogfood message-tool naming and
  collaborator permission, and the sim questions for BGP-class app pressure
  (`SIM-punaz`), group identity, promisebase reference naming (`SIM-ligan`),
  sparse chunk-feed decisions, and multi-repo site topology. For
  recovery-walkthrough purposes, turn 178 has no uncaptured loose end.
- `Work pending` no.
- `Proposed disposition` `resolved/transferred after sparse-CAS, promise-accounting, promisebase, interop, guide-claim, and multi-repo simulation-owner routing`
- `Write needed? yes/no` `no` further turn-178 write is needed after this pass.
- `Next` Turn 179 is next.

### Turn 179 — 2026-05-03 19:49 UTC

- `Turn 179 plain-English recap` Steve confirmed that discovery and transport of
  chunks are both promise-mediated in a promise economy. He added that the design
  needs multi-hop discovery and fetch, geofencing, permissioned groups, and
  conditional release, with the concrete example that Alice may send a chunk to
  Bob only if Bob promises not to send it outside group X. He then widened the
  promise-economy design space: "promise economy" might mean a real economy with
  promises issued as capability tokens, possibly transferable or permissioned,
  possibly with floating exchange rates where everyone is effectively their own
  central bank. He explicitly warned that this is only one possible model, must
  be tested in wire-lab, and could degenerate into cryptocurrency toxicity. He
  also asked to be reminded to connect the assistant to the promisebase GitHub
  repo with a PAT, while asking whether the assistant could pull it read-only for
  now. The assistant responded by reading promisebase but over-trusted design
  docs rather than implementation code. It pivoted TE-sihih toward "adopt
  promisebase wholesale," proposed eleven DFs around promisebase message format,
  layering, gitlike refs, grid CIDs, module dependency, migration, and `pgmsg`,
  and suggested PromiseGrid narrative material should live in promisebase docs.
  Later turns correct most of that response. Turn 180 identifies the doc-vs-code
  mistake and forces a git-log/code audit. Turns 181-184 show that only the lower
  CAS/Merkle/Rabin pieces are real enough to consider, while much of the
  promisebase layering and message-format material is design-only. Turn 191 locks
  the stance that promisebase is prototype evidence and wire-lab should be
  preferred in design conflicts. The valid turn-179 conclusions that remain are:
  conditional release and geofencing need their own owner, promise-economy
  mechanism neutrality needs testing, read-only/PAT and cross-repo work need
  explicit authorization, and promisebase cannot be treated as authority without
  code-first evidence.
- `Existing capture` `TODO-lilar` captured the raw turn in detail and flagged
  seven related UT rows. Later artifacts now capture the substance: TODO-ralud
  and `SIM-zarud` own conditional release, onward-restraint, geofencing, and
  recursive promise-graph ownership; TODO-rajig and `SIM-haros` own
  promise-economy mechanism neutrality; AGENTS-ppx B7 owns the
  ground-truth-before-citation process lesson; TODO-kituj and `DR-tumus` own
  concrete promisebase prior-art adoption; TODO-dozak owns wire-lab /
  promisebase merge-trajectory and scope-boundary questions; and TE-nizor records
  the later synthesis that promisebase is prototype evidence, not authority.
- `Gaps or contradictions` The historical turn-179 response was wrong in the
  specific way Steve later identified: it read promisebase docs as if they were
  implementation evidence and tried to turn wire-lab into a downstream adopter of
  a mostly design-only stack. It also bundled previously locked vocabulary fixes
  into a fresh DF and proposed cross-repo promisebase documentation scope without
  authorization. The remaining uncaptured design gap before this pass was the
  promise-economy mechanism spectrum: it was too specific and risky to hide under
  generic promise accounting records. `DI-vabij` fixes the simulation routing by
  adding `SIM-haros-promise-economy-spectrum`; `DI-pidag` fixes the owner routing
  by adding TODO-rajig and restoring TODO-lilar to historical evidence status.
- `Related UTs / owners` `UT-179.a` is closed for recovery by the later
  turn-180 / turn-184 correction path plus AGENTS-ppx B7 and TE-nizor C8: future
  external-repo claims need code-first evidence, not doc-only adoption. `UT-179.b`
  is captured by TODO-rajig and `SIM-haros-promise-economy-spectrum`, which ask
  how to test social assessment, reciprocal promises, capability tokens,
  transferability, floating exchange rates, and cryptocurrency-toxicity failure
  modes without forcing one economics model into the base protocol. `UT-179.c`
  and `UT-179.d` are captured by TODO-ralud and
  `SIM-zarud-conditional-release-geofencing`.
  `UT-179.e` is closed against TODO-vunub / TE-sihih's actual small landing and
  foundational vocabulary invariants: the vocabulary corrections from turns 175
  and 176 are background, not a reopened bundled DF. `UT-179.f` is closed as a
  historical cadence observation. `UT-179.g` is captured by TODO-dozak's
  merge-trajectory / scope-boundary owner and by existing cross-repo authorization
  discipline; no promisebase edit is authorized by the old sketch.
- `Owner/doc cleanup` Done. Added `DI-vabij`; added `DI-pidag`; added
  `SIM-haros-promise-economy-spectrum`; added TODO-rajig as the TODO owner for
  that sim; updated `simulations/README.md` and `DEV-GUIDE-RESOURCES.md`; added a
  scope-boundary subtask and question-log entry to TODO-dozak; added a turn-179
  closure pointer to the UT matrix; and restored `TODO-lilar` UT checkboxes as
  historical evidence rather than closure authority.
- `Remaining decisions or work` Downstream work remains open, but it is captured:
  TODO-rajig / `SIM-haros` hold promise-economy mechanism-neutrality questions;
  TODO-ralud / `SIM-zarud` hold conditional-release and geofencing ownership
  questions; TODO-kituj / `DR-tumus` hold concrete promisebase prior-art adoption;
  TODO-dozak holds wire-lab / promisebase merge or independence decisions; and
  AGENTS-ppx carries the code-first / authorization process lessons. For
  recovery-walkthrough purposes, turn 179 has no uncaptured loose end.
- `Work pending` no.
- `Proposed disposition` `resolved/transferred after code-first correction, promise-economy-spectrum simulation routing, conditional-release/geofencing routing, and promisebase scope-boundary routing`
- `Write needed? yes/no` `no` further turn-179 write is needed after this pass.
- `Next` Turn 180 is next.

### Turn 180 — 2026-05-03 20:21 UTC

- `Turn 180 plain-English recap` Steve stopped the turn-179 promisebase pivot with
  the correction: the assistant had read promisebase markdown files as if they
  were implemented code, and those files might be refactoring plans rather than
  finished work. The assistant acknowledged the mistake, audited promisebase more
  carefully, and separated implemented pieces from design-only claims. The
  corrected findings were that promisebase had real lower-level storage pieces
  such as filesystem key-value storage and some command/client/server scaffolding,
  but the five-element CBOR message format, layered `KV -> HashKV -> Message ->
  Refs -> Streaming` model, gitlike refs, grid CIDs, hypergraph refs, and much of
  the federation story were design plans rather than code wire-lab could adopt.
  The assistant invalidated the turn-179 "adopt promisebase wholesale" framing,
  listed which proposed DFs were no longer safe, and proposed two immediate
  recovery paths: treat promisebase docs as shared vocabulary only, or read the
  actual code and RFC series before deciding anything. It recommended reading
  code first. Later turns refine this correction. Turn 181 follows the code-first
  path but contains a wrong import-path claim that is later captured separately.
  Turns 183-184 show that `db/` and `kv/fs/` are the useful implemented evidence,
  while other promisebase surfaces remain broken or design-only. Turn 191 adds the
  missing third path that turn 180 should have surfaced: wire-lab may design
  independently and converge with promisebase only when useful, with wire-lab
  preferred in conflicts. The loose ends from turn 180 are therefore procedural
  and routing issues, not new protocol decisions by themselves: require
  code-first / ground-truthed evidence for external-repo claims; use the
  apologize-audit-invalidate-propose pattern after structural errors; do not
  describe promisebase as authority; preserve independent design with optional
  convergence as an available path; and correct the earlier replay summary that
  misdescribed turn 180 as context compression rather than context anchoring.
- `Existing capture` `TODO-lilar` records four related rows: `UT-180.a` for the
  structurally bad turn-179 pivot and corrective procedure, `UT-180.b` for the
  missing independent-design path, `UT-180.c` for the bad pending-line summary
  about context compression, and `UT-180.d` for the reusable
  apologize-audit-invalidate-propose recovery pattern. The UT matrix records that
  B6 and B7 landed in `AGENTS-ppx.md` under `DI-021-20260507-212254` and
  `DI-021-20260507-212255`. TE-nizor C8 records the later
  promisebase-as-prototype synthesis. TODO-kituj owns concrete promisebase /
  pitbase prior-art adoption and `DR-tumus`; TODO-dozak owns the separate
  wire-lab / promisebase convergence-or-independence question.
- `Gaps or contradictions` The raw turn correctly invalidated the wholesale
  promisebase adoption path, but it still framed recovery as two promisebase-
  centered options: shared vocabulary or read-code-first. Later turn 191 supplies
  the omitted independent-design path and makes promisebase prototype evidence,
  not authority. The old `TODO-lilar` pending-line characterization that turn 180
  caused context compression is also wrong: turn 180 re-anchored the replay by
  demoting invalid turn-179 claims and re-elevating the turn-170s design points.
- `Related UTs / owners` `UT-180.a` is captured by TODO-kituj / `DR-tumus`,
  TE-nizor C8, AGENTS-ppx B7, and the turn-179 / turn-180 replay notes. `UT-180.b`
  is captured by TODO-dozak's convergence-or-independence subtask and TE-nizor C8,
  with turn 191 supplying the stronger prefer-wire-lab rule. `UT-180.c` is closed
  for replay purposes by this corrected turn-180 note. `UT-180.d` is captured by
  AGENTS-ppx B6 and `DI-021-20260507-212254`.
- `Owner/doc cleanup` Done. Added `DI-zotol`; added this turn-180 report; marked
  `juhub.180` complete; and added a turn-180 transfer pointer to the UT
  verification matrix. No `TODO-lilar` UT checkbox was flipped.
- `Remaining decisions or work` Downstream work remains open, but it is captured:
  TODO-kituj / `DR-tumus` own concrete promisebase prior-art adoption; TODO-dozak
  owns whether wire-lab and promisebase converge, merge, or remain independent;
  TE-nizor records the promisebase-as-prototype / prefer-wire-lab synthesis; and
  AGENTS-ppx B6/B7 own the reusable process lessons. For recovery-walkthrough
  purposes, turn 180 has no uncaptured loose end.
- `Work pending` no.
- `Proposed disposition` `resolved/transferred after code-first correction, prototype-not-authority routing, independent-design-path routing, and process-rule capture`
- `Write needed? yes/no` `no` further turn-180 write is needed after this pass.
- `Next` Turn 181 is next.

### Turn 181 — 2026-05-04 05:00 UTC

- `Turn 181 plain-English recap` Steve selected the code-first path from turn 180
  and narrowed the scope with the instruction to read the code and focus on the
  `db/` directory. The assistant audited promisebase `db/` and reported that it
  is pitbase: a real content-addressable, deduplicating Merkle storage layer with
  blocks, trees, streams, path handling, WORM files, and Rabin content-defined
  chunking. The audit found concrete pieces: an on-disk `Db` rooted in `block/`,
  `tree/`, and `stream/`; WORM files whose hash covers a class header plus body;
  blocks as leaf bytes; trees as newline-separated child paths; streams as
  symlinks to tree roots rewritten atomically; path parsing and 3-hex fanout;
  Rabin chunking defaults of 512 KiB minimum and 8 MiB maximum; and an object
  interface shared by blocks and trees. It also found important absences:
  pitbase does not implement CIDv1 / multibase / multihash codec fields,
  deterministic CBOR or DAG-CBOR tree objects, the five-element CBOR message
  format, networking, group semantics, or content-addressed pointer files. The
  assistant reported 15 of 17 tests passing, with `TestPutStreamBig` and
  `TestPutStreamSmall` failing in the chunker / Merkle-builder integration path.
  It ended by offering to dig into `chunker_test.go:106`. Later turns update this
  turn in several ways. Turn 182 confirms the failing tests were regressions;
  turn 183 fixes them in promisebase and gets the relevant `db/` path to 17/17
  green; turn 184 broadens the audit beyond `db/` and shows that `db/` and
  `kv/fs/` are the useful implemented evidence while several other promisebase
  surfaces are broken or design-only; turn 191 makes promisebase prototype
  evidence rather than authority; and the current `DR-tumus` packet recommends
  using promisebase / pitbase as prior art only for the first L6 CAS spec rather
  than adopting it as substrate before the PromiseGrid pointer/raw profile is
  decided. The turn-181 answer also contained a concrete error: it suggested
  `github.com/t7a/pitbase/db`, but the correct repository path is
  `github.com/stevegt/promisebase/db`.
- `Existing capture` `TODO-lilar` records six related rows. `UT-181.a` flags the
  bad import path. `UT-181.b` flags the chunking-parameter mismatch between
  pitbase Rabin defaults and the earlier FastCDC-style small-chunk proposal.
  `UT-181.c` flags the two failing tests in the exact `PutStream` path wire-lab
  would care about. `UT-181.d` flags pitbase's class-header hashing as analogous
  to CIDv1 codec type binding. `UT-181.e` flags pitbase streams-as-symlinks as
  prior art for the pointer-file design. `UT-181.f` records the positive cadence
  lesson that the assistant offered a concrete next step and Steve's next terse
  response implicitly accepted it. TODO-kituj already owns the concrete TE-43
  promisebase / pitbase prior-art adoption question; `DR-tumus` is open with
  answerable DFs for starting profile, type binding, chunking lock scope, and
  promisebase stance; `SIM-jomag`, `SIM-gobaz`, `SIM-bobud`, and `SIM-kohad`
  contain the relevant simulation pressure.
- `Gaps or contradictions` The main contradiction is that turn 181's working
  frame was still too adoption-heavy. It correctly distinguished real `db/` code
  from design-only promisebase docs, but it treated direct `db/` dependency as a
  likely L6 substrate path. Later turn 191 and `DR-tumus` narrow that: pitbase is
  useful evidence, but the first PromiseGrid CAS profile should decide pointer
  objects, type binding, chunking scope, and adoption stance before treating
  promisebase as a substrate. The wrong import path must also be suppressed in
  later TE/DR prose.
- `Related UTs / owners` `UT-181.a` is routed to TODO-kituj's import-path
  correction and this turn note's explicit correction. `UT-181.b`, `UT-181.d`,
  and `UT-181.e` are routed to TODO-kituj / `DR-tumus` plus `SIM-jomag`,
  `SIM-gobaz`, `SIM-bobud`, and `SIM-kohad`. `UT-181.c` is satisfied for replay
  by the turn-183 promisebase fix and remains as evidence in TODO-kituj /
  `DR-tumus` for deciding what test-status threshold is enough before adoption.
  `UT-181.f` is a recorded cadence lesson paired with `UT-182.c`; it does not need
  a new owner before turn 182 is processed.
- `Owner/doc cleanup` Done. Added `DI-zarok`; added this turn-181 report; marked
  `juhub.181` complete; and added a turn-181 transfer pointer to the UT
  verification matrix. No `TODO-lilar` UT checkbox was flipped.
- `Remaining decisions or work` Downstream work remains open, but it is captured:
  TODO-kituj / `DR-tumus` own the CAS profile, chunking, type-binding,
  pointer-object, test-threshold, and promisebase / pitbase stance decisions; the
  CAS and chunking simulations preserve the design pressure; turn 183 supplies
  the regression-fix evidence; and turn 182 will separately process the
  implicit-yes cadence lesson. For recovery-walkthrough purposes, turn 181 has no
  uncaptured loose end.
- `Work pending` no.
- `Proposed disposition` `resolved/transferred after promisebase db audit routing, import-path correction, CAS-profile owner routing, and later regression-fix evidence`
- `Write needed? yes/no` `no` further turn-181 write is needed after this pass.
- `Next` Turn 182 is next.

### Turn 182 — 2026-05-04 05:05 UTC

- `Turn 182 plain-English recap` Steve replied to the turn-181 failing-test report
  with the short statement "That used to work." In context, that did two things:
  it confirmed the `TestPutStreamBig` and `TestPutStreamSmall` failures were
  regressions rather than known old failures, and it implicitly accepted the
  assistant's prior offer to dig into `chunker_test.go:106`. The assistant's
  response jumped ahead to a fix proposal: it offered to apply a one-file,
  roughly six-line change in `db/stream_test.go`, on a
  `ppx/fix-randstream-go120` twig in promisebase, verify the tests green, and
  either open a PR or leave the commit ready for Steve to push. It also offered a
  defer path: merely note the issue in the TE-sihih draft as a known-fixable
  regression if Steve wanted to postpone promisebase write access and scope
  decisions. Later turns fill in what turn 182 omitted. Turn 183 shows the actual
  diagnosis: Go 1.20 made `math/rand.Seed` a no-op for the global generator, and
  promisebase's Go 1.15-to-1.24 toolchain bump in commit `9a5634f` silently broke
  tests that expected deterministic global-rand output. The fix was to replace
  `rand.Seed(42)` plus global `rand.Read` with a per-`randStream` `*rand.Rand`
  seeded from `rand.NewSource(42)`, making the tests hermetic. Turn 183 applies
  the fix locally and reports 17/17 tests green; turn 186 later gets PAT access
  and pushes the work; turn 188 surfaces that the push confirmation was not
  visible enough. The loose ends from turn 182 are therefore process lessons, not
  PromiseGrid design questions: show the diagnostic when proposing a fix, state
  cross-repo auth requirements precisely, and recognize Steve's terse replies as
  possible implicit approval only when they answer a concrete offer already on
  the table.
- `Existing capture` `TODO-lilar` records three related rows: `UT-182.a` for the
  missing diagnostic in the turn-182 answer, `UT-182.b` for the loose "open a PR /
  leave the commit ready" phrasing before a promisebase PAT existed, and
  `UT-182.c` for the positive implicit-yes / regression-confirmation cadence
  lesson. Later rows also cover the actual fix (`UT-183.c`), local-only commit
  persistence risk (`UT-183.a`), PAT handling (`UT-185.*` / `UT-186.*`), and the
  foreground DONE visibility gap (`UT-188.a`). AGENTS-ppx B1 covers separate
  authorization from execution when a turn says "ask Steve"; B4 covers PAT scope
  and expiry checks; B2 covers foreground DONE reporting; B6/B7 cover structural
  error recovery and ground-truthed evidence.
- `Gaps or contradictions` The assistant's turn-182 fix proposal was probably
  correct, but it was under-narrated: readers could not see how the assistant got
  from "tests fail" to "one-file randStream fix" until turn 183. The offer also
  bundled operations with different authorization requirements: preparing a local
  commit, pushing to promisebase, and opening a PR are not equivalent when no
  promisebase PAT exists yet. The replay should preserve that distinction without
  retroactively pretending turn 182 already said it cleanly.
- `Related UTs / owners` `UT-182.a` is captured by this turn note, the turn-183
  diagnostic, and the replay/matrix lesson that fix proposals should show root
  cause in the same answer. `UT-182.b` is captured by AGENTS-ppx B1/B4 plus the
  later PAT and push-status turns; the precise lesson is recorded here and in the
  matrix. `UT-182.c` is captured as a positive cadence lesson paired with
  `UT-181.f`; it does not require a new durable rule beyond the per-turn replay
  because it applies only when a prior concrete offer is already on the table and
  the reply is not a no or redirect.
- `Owner/doc cleanup` Done. Added `DI-lumal`; added this turn-182 report; marked
  `juhub.182` complete; and added a turn-182 transfer pointer to the UT
  verification matrix. No `TODO-lilar` UT checkbox was flipped.
- `Remaining decisions or work` Downstream process and cross-repo work remains
  open in later turns, but it is captured: turn 183 owns the actual randStream fix
  evidence and local-only persistence concern; turns 185-186 own PAT handling and
  push mechanics; turn 188 owns foreground DONE visibility. For
  recovery-walkthrough purposes, turn 182 has no uncaptured loose end.
- `Work pending` no.
- `Proposed disposition` `resolved/transferred after diagnostic-transparency, auth-requirement, implicit-yes cadence, and later-fix routing`
- `Write needed? yes/no` `no` further turn-182 write is needed after this pass.
- `Next` Turn 183 is next.

### Turn 183 — 2026-05-04 05:27 UTC

- `Turn 183 plain-English recap` Steve said "Do it," selecting the turn-182
  option to apply the promisebase randStream fix locally. The assistant applied
  the fix on the promisebase twig `ppx/fix-randstream-go120`, and turn 183 is
  where the previously missing diagnostic finally became visible: after
  promisebase moved from Go 1.15-era behavior to a Go 1.24 toolchain path,
  `math/rand.Seed(42)` no longer made the global random generator deterministic
  for the tests, so `TestPutStreamBig` and `TestPutStreamSmall` failed. The fix
  replaced the old `rand.Seed(42)` plus global `rand.Read` pattern with a
  per-`randStream` `*rand.Rand` seeded from `rand.NewSource(42)`, which made the
  test stream hermetic. The assistant reported the one-file patch, the local
  commit `d98b5d3`, and all 17 promisebase `db/` tests passing. Because no
  promisebase PAT existed yet, the commit was only local at this point; turn 186
  later grants the PAT and pushes/merges the fix, and turn 188 later shows that
  the push confirmation needed to be more visible. Turn 183 also recommended
  reading more of shipped promisebase before drafting TE-sihih, rather than
  relying on old `x/*.md` notes; turn 184 accepts that recommendation and then
  finds broader partial rot outside `db/`. The lasting conclusions are therefore
  bounded: the `db/` chunker/Merkle path was green after the fix, promisebase
  should still be treated as prior art rather than canon, external-repo
  architecture claims require code-first verification, and cross-repo local-only
  commits need a persistence rule before PAT-gated push work can be considered
  operationally safe.
- `Existing capture` `TODO-lilar` records `UT-183.a` for the local-only
  promisebase commit with no session-survivable persistence path, `UT-183.b` for
  the read-shipped-code-before-claims procedural default, and `UT-183.c` for the
  17/17 `db/` test result plus the unresolved dependency-test threshold. TODO-kituj
  already records the 17/17 green `db/` evidence as prior-art input for TE-43, and
  `DR-tumus` asks the later promisebase / pitbase stance question without
  treating promisebase as authoritative. AGENTS-ppx B7 covers ground-truthing
  before citation, and the new `TODO-fapev` owns the missing local-only
  cross-repo persistence rule.
- `Gaps or contradictions` The fix result is real evidence, but it applies only
  to promisebase `db/` tests, not to the whole promisebase repo. Turn 184 later
  finds Docker SDK, FUSE, server, and daemon uncertainty, so turn 183 must not be
  read as "promisebase is green" or "wire-lab should adopt promisebase." The other
  gap is procedural: a local-only cross-repo commit waiting on PAT access could
  have been lost on reset, and that risk was not owned by the existing PAT or
  foreground-DONE rules.
- `Related UTs / owners` `UT-183.a` is routed to `TODO-fapev`, which will decide
  whether the durable mechanism is a private patch file, bundle, private mirror,
  or another approved persistence path. `UT-183.b` is captured by AGENTS-ppx B7
  and this replay note as the code-first external-repo evidence rule. `UT-183.c`
  is routed to TODO-kituj / `DR-tumus`: the `db/` path is green evidence, but the
  broader promisebase adoption stance and required test-status threshold remain
  downstream decisions.
- `Owner/doc cleanup` Done. Added `DI-kegar`; created `TODO-fapev`; indexed
  `TODO-fapev` in the master TODO list; added this turn-183 report; marked
  `juhub.183` complete; and added a turn-183 transfer pointer to the UT
  verification matrix. No `TODO-lilar` UT checkbox was flipped.
- `Remaining decisions or work` Downstream work remains open, but it is captured:
  `TODO-fapev` owns the harness persistence mechanism for local-only cross-repo
  commits pending auth; TODO-kituj / `DR-tumus` own the PromiseGrid-facing
  promisebase / pitbase stance and test-threshold decisions; AGENTS-ppx B7 owns
  the general ground-truthing rule for external artifact claims. For
  recovery-walkthrough purposes, turn 183 has no uncaptured loose end.
- `Work pending` no.
- `Proposed disposition` `resolved/transferred after randStream fix evidence, code-first external-repo rule, local-only persistence owner, and TE-43 prior-art routing`
- `Write needed? yes/no` `no` further turn-183 write is needed after this pass.
- `Next` Turn 184 is next.

### Turn 184 — 2026-05-04 05:29 UTC

- `Turn 184 plain-English recap` Steve answered "Do 2," choosing the turn-183
  recommendation to read the rest of promisebase code before drafting TE-sihih.
  The assistant then produced the largest promisebase audit in this replay slice:
  it surveyed `db/`, `kv/fs/`, `cmd/pb/`, `cmd/pitd/`, `fuse/`, `server/`, the
  RFC series, the `x/` design notes, ROADMAP items, current vocabulary, sensible
  wire-lab / promisebase coupling, near-term promisebase fixes, and a ten-item
  TE-sihih DF list. The useful technical conclusions were bounded. `db/` was the
  mature CAS / Rabin / Merkle / stream module and was green after the randStream
  fix; `kv/fs/` was a small working KV layer and looked like the intended new
  bottom of a future refactor; `cmd/pb/` contained useful exec-by-CID prior art
  but was broken against newer Docker SDK APIs; `fuse/`, `server/`, and daemon
  surfaces showed unfinished or failing work; current promisebase code used
  algo/hexhash rather than CIDv1; and the "promise economy" itself was not
  implemented in promisebase, with RFC-1005's test tree CID + executable tree CID
  + cache-on-pass idea serving only as prior-art seed. The audit reframed the
  earlier invalid wholesale-adoption idea into a more defensible split:
  promisebase supplies prior art for local CAS/KV/execution pieces, while wire-lab
  owns federation, sparse-CAS feeds, groups, promise economics, partition
  tolerance, multi-hop discovery, conditional release, and protocol simulations.
  The assistant then listed ten TE-sihih decision areas, but did not actually
  walk any one DF with alternatives, consideration paragraphs, and a
  recommendation. It also asked whether promisebase should be treated as a solo
  Steve project or as having a broader active community; turn 187 later answers
  that Steve is solo in promisebase right now. Later artifacts further narrow the
  conclusions: TE-sihih lands as a smaller substrate-agnostic layered-model TE,
  TODO-kituj / `DR-tumus` own concrete promisebase / pitbase adoption, TODO-dozak
  owns merge-trajectory questions, TODO-rajig / `SIM-haros` own promise-economy
  mechanism pressure, and AGENTS-ppx now owns the process defects exposed by the
  turn.
- `Existing capture` `TODO-lilar` records seven turn-184 rows: `UT-184.a` for the
  flat ten-DF list violating one-DF-at-a-time discipline; `UT-184.b` for the
  collaborator-name propagation risk; `UT-184.c` for the undefined "grokker
  boilerplate" count; `UT-184.d` for the `kv/fs` refactor / dependency-target
  ambiguity; `UT-184.e` for Docker SDK rot in `cmd/pb`; `UT-184.f` for RFC-1005
  as promise-economy prior art; and `UT-184.g` for the broader partial-rot finding
  outside `db/` and `kv/fs`. Existing AGENTS-ppx B5 covers one-DF-at-a-time
  discipline; AGENTS-ppx B3 and TODO-rohub cover collaborator permission;
  AGENTS-ppx B7 now explicitly covers reproducible pattern-count claims;
  TODO-kituj / `DR-tumus` cover promisebase prior-art adoption, partial rot,
  `kv/fs` layering, and test-threshold consequences; TODO-rajig and `SIM-haros`
  now include the RFC-1005 prior-art pressure.
- `Gaps or contradictions` Turn 184 was the right corrective audit after the
  turn-179 promisebase-doc overread, but it still mixed several abstraction
  levels. The ten-DF list was a scope inventory, not actual DF walking. The audit
  proved that some local promisebase components are useful prior art; it did not
  prove that promisebase is ready to be a PromiseGrid substrate or that wire-lab
  should depend on all of it. The collaborator question should have been phrased
  without naming a protected collaborator, and the `x/discussion.md` count needed
  the counted pattern defined.
- `Related UTs / owners` `UT-184.a` is closed for replay by AGENTS-ppx B5 and the
  later TE-sihih contraction in TODO-vunub; no flat DF list from turn 184 is
  treated as a locked decision. `UT-184.b` is routed to TODO-rohub and AGENTS-ppx
  B3. `UT-184.c` is routed to the expanded AGENTS-ppx B7 pattern-count rule.
  `UT-184.d`, `UT-184.e`, and `UT-184.g` are routed to TODO-kituj / `DR-tumus`
  as promisebase / pitbase prior-art, dependency-target, partial-rot, and
  test-threshold inputs; later turn 187 explicitly holds other promisebase fix
  work rather than making it the next replay task. `UT-184.f` is routed to
  TODO-rajig and `SIM-haros`, where RFC-1005's test-driven fabric is treated as
  a prior-art seed to compare against promise-economy mechanisms, not as settled
  PromiseGrid economics.
- `Owner/doc cleanup` Done. Added `DI-nulak`; updated AGENTS-ppx B7; updated
  TODO-rohub, TODO-kituj, TODO-rajig, `SIM-haros`, and DEV-GUIDE-RESOURCES.md;
  added this turn-184 report; marked `juhub.184` complete; and added a turn-184
  transfer pointer to the UT verification matrix. No `TODO-lilar` UT checkbox was
  flipped.
- `Remaining decisions or work` Downstream work remains open, but it is captured:
  TODO-kituj / `DR-tumus` own the concrete promisebase / pitbase adoption and
  test-threshold questions; TODO-rajig / `SIM-haros` own RFC-1005 as
  promise-economy prior-art pressure; TODO-rohub owns collaborator permission;
  AGENTS-ppx owns the process rules. For recovery-walkthrough purposes, turn 184
  has no uncaptured loose end.
- `Work pending` no.
- `Proposed disposition` `resolved/transferred after code-first promisebase audit, process-rule routing, RFC-1005 sim capture, and TE-43 prior-art routing`
- `Write needed? yes/no` `no` further turn-184 write is needed after this pass.
- `Next` Turn 185 is next.

### Turn 185 — 2026-05-04 14:28 UTC

- `Turn 185 plain-English recap` Steve asked whether the assistant could hold and
  use two different GitHub PATs in the same session, one for wire-lab and one for
  promisebase. The assistant answered yes and described a concrete operating
  pattern: keep separate secret files or secret names for each repo, use the
  correct token only for the matching remote, invoke the Git credential helper
  per command rather than installing a shared credential, and keep the same bot
  identity while separating the token material by repo. The assistant also
  recommended fine-grained PATs scoped narrowly to the target repo, especially for
  promisebase where the immediate need was to push the already-prepared
  randStream fix. The answer was directionally right and turn 186 uses the
  pattern successfully, but it missed three durable safeguards. First, when
  recommending a fine-grained PAT, the assistant should also recommend a short
  practical expiry so leaked or echoed context becomes stale quickly. Second,
  labeling a file `readonly` is only a convention; actual token scope and expiry
  need to be verified before a write operation. Third, the assistant acknowledged
  that a wire-lab PAT had appeared in carry-over context, but did not state the
  redaction mechanism that should apply to all summaries, handoffs, walk notes,
  and commit messages. Later turn 186 confirms the redaction lesson by recording
  the promisebase PAT only as a secret marker, not as literal bytes. The durable
  conclusion is process-level: multiple repo-scoped credentials are allowed only
  with per-remote separation, shortest-practical expiry, actual scope/expiry
  verification, and strict secret redaction.
- `Existing capture` `TODO-lilar` records `UT-185.a` for missing short-expiry
  guidance, `UT-185.b` for the non-enforced read-only filename convention, and
  `UT-185.c` for PAT bytes echoing into carry-over summaries without a stated
  redaction mechanism. AGENTS-ppx B4 already covered secret redaction and
  scope/expiry checks, but it did not yet say per-remote separation,
  shortest-practical expiry, summary redaction before emission, or filename
  suffixes as non-enforcement as explicitly as turn 185 requires.
- `Gaps or contradictions` The assistant's operational answer worked in practice,
  because turn 186 uses the promisebase PAT without intentionally crossing it
  with the wire-lab token. The gap was defensive rigor, not functional capability:
  token separation must be enforced by command structure and verification, not by
  memory or filenames, and persistent context must never carry literal token
  bytes forward.
- `Related UTs / owners` `UT-185.a`, `UT-185.b`, and `UT-185.c` are all routed to
  AGENTS-ppx B4 as a closed process rule. `UT-186.a` will later be processed as
  the first concrete application of this redaction rule, but it does not need a
  separate turn-185 owner.
- `Owner/doc cleanup` Done. Added `DI-lifub`; expanded AGENTS-ppx B4; added this
  turn-185 report; marked `juhub.185` complete; and added a turn-185 transfer
  pointer to the UT verification matrix. No `TODO-lilar` UT checkbox was flipped.
- `Remaining decisions or work` No downstream decision remains for turn 185:
  credential storage path choices and actual PAT use are runtime operations, and
  the durable policy is now captured in AGENTS-ppx B4. For
  recovery-walkthrough purposes, turn 185 has no uncaptured loose end.
- `Work pending` no.
- `Proposed disposition` `resolved/transferred after two-PAT process rule, short-expiry guidance, scope verification, and redaction-rule capture`
- `Write needed? yes/no` `no` further turn-185 write is needed after this pass.
- `Next` Turn 186 is next.

### Turn 186 — 2026-05-04 14:29 UTC

- `Turn 186 plain-English recap` Steve supplied the promisebase PAT; the replay
  records it only as a secret marker, not as literal token bytes. The assistant
  saved it as the promisebase-specific credential, used it only for the
  promisebase remote, and kept it separate from the wire-lab credential. The
  assistant then pushed and fast-forwarded the promisebase randStream fix: the
  local `ppx/fix-randstream-go120` work from turn 183 landed on promisebase main
  at commit `d98b5d3`, the twig was deleted locally and remotely, and `db/` tests
  were green on main. This completed the first end-to-end cross-repo contribution
  in the replay slice: failing tests found in turn 181, regression confirmed in
  turn 182, fix prepared locally in turn 183, two-PAT handling established in
  turn 185, and push/merge/delete completed in turn 186. The assistant also
  noticed a side-channel build hazard: `GOTOOLCHAIN=auto` plus a newer installed
  Go toolchain can lazily rewrite `go.mod` with a higher `toolchain` line during
  `go test`; the assistant checked out that drift before pushing so promisebase
  main kept Steve's existing Go version/toolchain setting. Finally, the assistant
  pivoted back toward TE-sihih by re-asking the solo-vs-broader-community question
  from turn 184 and by asking whether to tackle Docker SDK / FUSE rot or hold it.
  Turn 187 answers both: Steve is solo in promisebase right now, and other
  promisebase work should be held while returning to TE-sihih. The durable
  conclusions are process-level: PAT redaction was correctly applied, foreground
  cross-repo completion needed clearer reporting, collaborator-sensitive questions
  must be paraphrased when re-quoted, and operational build diagnostics need a
  durable owner.
- `Existing capture` `TODO-lilar` records `UT-186.a` for first application of the
  PAT-redaction discipline, `UT-186.b` for the GOTOOLCHAIN auto-bump diagnostic
  not yet having a persistent home, and `UT-186.c` for the collaborator-name
  propagation caused by re-quoting the turn-184 question. AGENTS-ppx B4 now owns
  secret redaction and per-remote credential separation; TODO-rohub and
  AGENTS-ppx B3 own collaborator permission / non-mention constraints; the new
  `TODO-nasat` owns cross-repo build-hazard capture; and TODO-kituj / `DR-tumus`
  already own the Docker SDK, FUSE, partial-rot, and promisebase-readiness
  consequences.
- `Gaps or contradictions` The cross-repo fix was successful, but the assistant
  mixed too many things into one response: credential storage, push/merge/delete,
  db test status, GOTOOLCHAIN drift, TE-sihih pivot, and two new questions. Turn
  188 later shows that Steve did not register the push status, so turn 186's
  completion should have been foregrounded with a clear DONE line. The
  GOTOOLCHAIN observation was useful, but it needed a durable owner. The
  collaborator question should have been paraphrased instead of re-quoting the
  protected name from the earlier question.
- `Related UTs / owners` `UT-186.a` is captured by AGENTS-ppx B4 and the redacted
  turn-186 replay note. `UT-186.b` is routed to `TODO-nasat` and the expanded
  AGENTS-ppx B7 cross-repo hazard-capture rule. `UT-186.c` is routed to TODO-rohub
  and AGENTS-ppx B3. The deferred Docker SDK / FUSE fixes remain routed to
  TODO-kituj / `DR-tumus`; turn 187's "hold other promisebase work" confirms they
  should not interrupt the recovery walkthrough.
- `Owner/doc cleanup` Done. Added `DI-zagus`; created `TODO-nasat`; indexed
  `TODO-nasat`; expanded AGENTS-ppx B7; added this turn-186 report; marked
  `juhub.186` complete; and added a turn-186 transfer pointer to the UT
  verification matrix. No `TODO-lilar` UT checkbox was flipped.
- `Remaining decisions or work` Downstream work remains open, but it is captured:
  `TODO-nasat` owns where cross-repo build hazards should live permanently;
  TODO-rohub owns collaborator permission; AGENTS-ppx B4/B7 own credential
  redaction and hazard-capture process; TODO-kituj / `DR-tumus` own the
  promisebase partial-rot and adoption-readiness questions. For
  recovery-walkthrough purposes, turn 186 has no uncaptured loose end.
- `Work pending` no.
- `Proposed disposition` `resolved/transferred after redacted PAT handling, promisebase fix push evidence, GOTOOLCHAIN hazard owner, and collaborator re-quote routing`
- `Write needed? yes/no` `no` further turn-186 write is needed after this pass.
- `Next` Turn 187 is next.
