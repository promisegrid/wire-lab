# TODO-lugag: Promise Theory GA prompt, scenario, and rescore correction

## Prior aliases

None. This TODO was minted after the proquint-handle migration.

## Status

Open. This TODO owns the Promise Theory correction across GA child generation,
GA scoring and gating, root scenario quality and realism, proposal-child
disposition, and broad rescoring of previously scored cells. It coordinates with
`TODO-tapur`, `TODO-dadub`, `TODO-ralud`, and `TODO-rozas`, but remains a
separate harness TODO because the problem is not only a runner bug. It is a
cross-artifact correction to the way wire-lab models multi-agent behavior.
Source: `DI-lumit`.

`SIM-mipag-child-split-release-capsule-template` and
`SIM-dafod-child-dual-surface-contract-witness-ledger` are explicitly deferred.
They are not promotable as-is under this TODO's Promise Theory correction.
Source: `DI-safij`.

## Decision Intent Log

ID: DI-safij
Date: 2026-05-23 17:51:58
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Track the Promise Theory correction as a new harness TODO separate
from `TODO-tapur`, covering GA prompt hardening, PT-gated scoring, scenario
corpus realism and PT coverage, deferred child-promotion review, and the plan
to re-score previously scored cells under the corrected rules.
Intent: PromiseGrid is a decentralized virtual machine spanning nodes owned by
different legal entities. Agents cooperate by making reciprocal promises,
keeping or breaking them, and building local trust from relationship history.
Protocol objects and pCIDs help agents make, recognize, remember, and evaluate
promises; they do not command agents, grant permission, impose policy, or act
as trust authorities. The same Promise Theory rules apply to sending,
receiving, storing, computing, relaying, and trusting returned results.
Constraints: All trust is local, relative, and relationship-specific; there is
no global trust authority. Alice must build trust in Bob before sending Bob her
data for storage or computation, and before trusting computation results Bob
returns. Reciprocal promises are first-class. Scenario rewrites must improve
realism and PT framing without rewriting historical scored-result bytes in
place. `mipag` and `dafod` must not be promoted unchanged. Future rescoring must
record new evidence rather than mutating historical JSON artifacts.
Affects: `protocols/wire-lab.d/TODO/TODO-lugag-promise-theory-ga-prompt-scenario-and-rescore-correction.md`;
`protocols/wire-lab.d/TODO/TODO.md`; future `tools/ga-runner/`; future
`scenarios/`; future `simulations/`; future result comparison and rescore
artifacts.

ID: DI-lumit
Date: 2026-05-23 18:09:50
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Ground `TODO-lugag` explicitly in Mark Burgess's basic Promise Theory
rules and plan to carry those rules into GA scoring rubrics, GA child-generation
prompts, and future GA-runner code. Treat non-voluntary cooperation,
imposition-heavy framing, conformance-first language, and external-authority
reasoning as highly suspect in generation and scoring unless they are explicitly
decomposed into autonomous promises, accept/use promises, impositions, and local
assessments.
Intent: `TODO-lugag` should not merely push the repo toward a promise-first
style; it should align with Burgess's basic PT mechanics. The core rules include
autonomous agents, scoped declarations of intent, no promises on behalf of
other agents, no guarantee of outcome, and trust as a local assessment of
whether a promise will be kept. These are fundamental enough to function like
base laws for PromiseGrid design evaluation. The scorer and generator should be
made aware of those rules directly by including short references to Burgess's
work in prompts and rubric text.
Constraints: Keep `DI-safij`'s local-trust and reciprocal-promise direction.
Do not turn Burgess compatibility into a simplistic banned-word policy:
`contract`, `authority`, `authorization`, `permission`, and `conformance` may
still appear when analyzed or decomposed in PT terms, but score and generation
logic should be highly suspicious of designs that treat them as primitive
external force or global authority. Future prompts should cite Burgess's
foundational PT writings concisely enough to guide the model without bloating
prompt size.
Affects: `protocols/wire-lab.d/TODO/TODO-lugag-promise-theory-ga-prompt-scenario-and-rescore-correction.md`;
future `tools/ga-runner/score.go`; future `tools/ga-runner/generate.go`;
future `scenarios/README.md`; future GA tests and comparison reports.
Supersedes: DI-safij

ID: DI-movur
Date: 2026-05-23 18:16:08
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement `lugag.1`, `lugag.5`, and `lugag.7` in the scorer as a new
v3 score/result contract that keeps historical v1/v2 artifacts valid while
making Promise Theory fundamentals and the PT gate explicit. New scorer results
must carry the canonical Burgess-grounded PT rule set in the stored rubric and
must carry a structured `assessment.pt_gate` object with deterministic
statuses. The runner must apply the PT gate before computing weighted fitness,
so PT-invalid designs cannot score as promotable and PT-reframe-needed designs
are capped below PT-clean contenders.
Intent: `lugag` needs a first-class, machine-readable PT checklist and gate,
not only prose. Putting the rule set in the stored rubric makes later review
and rescoring auditable. Using a v3 result contract avoids retroactively
invalidating historical v2 score artifacts. Applying the gate before weighted
fitness makes the PT correction operational instead of advisory.
Constraints: Keep v1/v2 result validation intact. New score runs emit
`promisegrid.ga.result.v3`, `ga-rubric-20260523-v3`, and
`ga_score_payload_v3`. The PT gate must classify designs as `pt_clean`,
`pt_reframe_needed`, or `pt_invalid`. The canonical rule checklist must cover
autonomous agents, scoped intent, no promises on behalf of others, no
guaranteed outcomes, trust as local assessment, and the non-equivalence of
accept/use promises with obligations or impositions. The scorer prompt must
include concise Burgess reference notes without bloating the prompt.
Affects: `protocols/wire-lab.d/TODO/TODO-lugag-promise-theory-ga-prompt-scenario-and-rescore-correction.md`;
`tools/ga-runner/result.go`; `tools/ga-runner/score.go`;
`tools/ga-runner/validate.go`; `tools/ga-runner/ga_runner_test.go`
Supersedes: DI-lumit

ID: DI-tavaz
Date: 2026-05-23 19:18:00
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Rework PT-drifting scored sims by creating successor simulations
instead of rewriting the scored simulation directories in place. The immediate
successors are a selective-send/onward-promise rewrite for `gibut`, a
promise-evidence-record rewrite for `savak`, a pCID promise-evidence-index
rewrite for `tizad`, an outer-promise nested-signed-payload rewrite for
`janov`, and a fixed-header variable-body rewrite for `sajar`.
Intent: Historical scores must remain evidence about the exact bytes that were
scored. In-place rewrites would mix old scores with new semantics under one
`sim_id` and would undermine later audit, comparison, and GA parent-selection
honesty. Successor sims let wire-lab compare the old PT-drifting design against
the PT-clean rewrite directly.
Constraints: Do not rewrite the scored predecessor sim trees in place. New
successor sims should explicitly name the predecessor they replace and should
restate the design in Promise-Theory-first vocabulary. Promotion, breeding, and
later rescoring should prefer the successor lineage once evidence exists.
Affects: `protocols/wire-lab.d/TODO/TODO-lugag-promise-theory-ga-prompt-scenario-and-rescore-correction.md`;
`simulations/README.md`; new `simulations/SIM-fonom-conditional-release-selective-send-onward-promises/`;
new `simulations/SIM-vuliv-scoped-promise-evidence-records/`;
new `simulations/SIM-konit-pcid-promise-evidence-index/`;
new `simulations/SIM-pobod-grid-envelope-outer-promise-nested-signed-payload/`;
new `simulations/SIM-dutam-grid-envelope-fixed-header-variable-body/`
Supersedes: DI-movur

ID: DI-gobul
Date: 2026-05-23 23:25:00
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Mark the five PT-drifting predecessor sims as non-breeding
`question-home` specimens by adding `SIM-META.json` files to their simulation
roots. The affected predecessors are `SIM-gibut-conditional-release-group-session-local`,
`SIM-savak-scoped-claim-card-audit-ledger`,
`SIM-tizad-scoped-conformance-citation-ledger`,
`SIM-janov-grid-envelope-layer-pcid-nested-signed-payload`, and
`SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields`.
Intent: Once PT-clean successor sims exist, the predecessor specimens should
remain available as scored historical evidence and question homes, but they
should stop silently re-entering default GA parent selection.
Constraints: Do not rewrite the predecessor design docs. Preserve their scored
history and keep them scoreable. Exclude them from default breeding through
standard `SIM-META.json` role handling rather than special-case runner code.
Affects: `protocols/wire-lab.d/TODO/TODO-lugag-promise-theory-ga-prompt-scenario-and-rescore-correction.md`;
`simulations/README.md`;
new `simulations/SIM-gibut-conditional-release-group-session-local/SIM-META.json`;
new `simulations/SIM-savak-scoped-claim-card-audit-ledger/SIM-META.json`;
new `simulations/SIM-tizad-scoped-conformance-citation-ledger/SIM-META.json`;
new `simulations/SIM-janov-grid-envelope-layer-pcid-nested-signed-payload/SIM-META.json`;
new `simulations/SIM-sajar-grid-envelope-variable-arity-pcid-defined-fields/SIM-META.json`
Supersedes: DI-tavaz

## Scope

This TODO covers:

- Promise-Theory-first GA child-generation prompts.
- Promise-Theory-gated GA scoring and review.
- Root scenario PT coverage, realism, actor usage, and trust-building quality.
- PT rework of currently deferred children and other PT-drifting sims.
- Planning and executing broad rescoring after the corrected rules land.

This TODO does not itself settle final wire formats, final result-schema field
names, or the final list of canonical successor sims. It is the coordination
owner that drives those later decisions to DI before implementation.

## Predecessor context

- `TODO-tapur` owns the current GA runner mechanics, rubric-v2 result contract,
  prompt bundling, scoring, generation, validation, acceptance, culling, and
  promotion procedure.
- `TODO-dadub` created the root `scenarios/` corpus and shared root scenario
  contract.
- `TODO-ralud` owns the conditional-release / geofencing architectural question,
  which must now be reframed in Promise Theory terms as selective sending plus
  reciprocal promises rather than imposed permission.
- `TODO-rozas` consumes mature evidence for dev-guide writer resources; it
  should not summarize non-PT-clean design drift as if it were consensus.

## Burgess grounding

The later implementation and review work under this TODO should stay grounded in
Mark Burgess's Promise Theory writings, especially the basic-rule statements in
`In Search of Certainty`, `Promise Theory: Principles and Applications`, `A
Treatise on Systems, Volume 2`, and the FAQ/trust notes. This grounding matters
because `lugag` is trying to correct not only vocabulary drift but also a
deeper conceptual drift back toward RPC, policy-enforcement, and command-style
thinking.

## Current problems

- GA child-generation prompts do not require the model to identify promising
  agents, reciprocal promises, make/break evidence, or what remains local trust
  judgment.
- The scorer still allows technically neat but conceptually wrong designs to
  score well even when they use permission, authorization, conformance,
  contract, or policy-enforcement framing as if those were external authority.
- PT principles are not being applied uniformly to send, receive, relay,
  storage, computation, result return, refusal, and trust updates.
- The TODO does not yet explicitly elevate Burgess's basic PT rules into a
  canonical checklist for GA prompts, rubrics, and later code.
- The current plan is suspicious of authority-style language, but does not yet
  say clearly enough that impositions, conformance claims, and other
  non-voluntary cooperation patterns should score poorly unless decomposed into
  autonomous promises and local assessments.
- Many scenarios still read as boilerplate rather than thought experiments:
  they under-use Alice/Bob/Carol/Mallory roles, have weak setup/stimulus/trust
  history, and often drift into non-PT vocabulary such as `authority`,
  `conformance`, and `access control`.
- Existing high-scoring or near-promoted artifacts may need to be reinterpreted
  or demoted once PT-corrected scoring exists.

## Subtasks

- [x] lugag.1 Lock a canonical set of Burgess-grounded Promise Theory rules for
  PromiseGrid evaluation and generation. At minimum include: agents are
  autonomous; a promise is a scoped declaration of intent; no agent can make a
  promise on behalf of another; promises do not guarantee outcomes; trust is a
  local assessment of whether a promise will be kept; and promises to
  receive/use are not equivalent to impositions, obligations, or promises to
  give. Plan to make this rule set a first-class checklist in the scorer rubric,
  generation prompts, and future GA-runner code paths.
- [ ] lugag.2 Write a compact PT interaction checklist that applies the same
  discipline to send, receive, relay, store, retrieve, compute, return-result,
  forward, refuse, ignore, delay, supersede, and revoke behaviors. For each
  interaction class, require explicit answers to: which agent promises what,
  what evidence of keeping/breaking looks like, how local trust is updated, and
  what no artifact can decide on the agent's behalf. Explicitly distinguish
  promises, accept/use promises, impositions, and assessments so the scorer does
  not collapse them into one another.
- [ ] lugag.3 Lock the Promise Theory vocabulary and anti-vocabulary used by GA
  prompts, scoring, scenario docs, and promotion review. Define preferred terms
  such as promise, promiser, promisee, reciprocal promise, local trust, trust
  history, kept promise, broken promise, accept/use promise, assessment, and
  pCID-referenced promise; define suspect terms such as permission, who may,
  authorization, conformance, contract, policy enforcement, compliance, access
  control, and authority that require explicit PT reframing before they are
  allowed. Be highly suspicious of non-voluntary cooperation, imposition-heavy
  framing, and conformance-first reasoning in both scoring and generation.
- [ ] lugag.4 Harden GA child-generation prompts so every generated design must
  identify the promising agents, what each agent voluntarily promises, which
  promises are reciprocal, what make/break evidence records matter, and what
  remains local trust judgment. Require explicit treatment of send, receive,
  storage, and computation when relevant, and require the child to explain how
  Alice decides whether to trust Bob enough to send data and whether to trust
  computation results Bob returns. Add concise references to Burgess's PT work
  in the prompt so the generation model sees the basic autonomy/voluntary
  cooperation rules while it reasons.
- [x] lugag.5 Add an explicit Promise Theory gate before weighted scoring. The
  gate must classify designs as PT-clean, PT-reframe-needed, or PT-invalid.
  PT-invalid designs cannot be promotable regardless of technical neatness.
  PT-reframe-needed designs may survive only as question homes or rework
  candidates. PT-clean status is required for any design that can influence
  promotion, consensus, or guide-level summary.
- [ ] lugag.6 Harden the scorer rubric and penalties. Make
  `promise_vocabulary` close to a correctness axis rather than a soft style
  preference; make high `promisegrid_alignment` impossible unless the PT gate is
  clean; severely penalize external-authority language unless the scorer
  explains the local-promise interpretation; and penalize designs that treat
  storage or computation as neutral service calls rather than promises between
  agents with local trust consequences. Add concise references to Burgess's PT
  writings in the rubric so the scoring model is reminded of the governing
  principles rather than relying on ambient intuition.
- [x] lugag.7 Decide the PT-gate result shape and scoring-report contract. Lock
  whether the PT gate becomes an explicit structured field in score results, a
  required assessment subsection, or both. Ensure the promotion and compare
  tooling can key off the PT gate deterministically rather than inferring it
  from prose alone.
- [ ] lugag.8 Audit `scenarios/README.md` and strengthen the shared root
  scenario contract so every scenario inherits explicit PT requirements for
  reciprocal promises, local trust, send/receive/storage/computation discipline,
  observable make/break history, and trust updates. Tighten the actor guidance
  so Alice, Bob, Carol, Dave, Ellen, Frank, and Mallory are used as concrete
  cooperating or adversarial agents rather than as decorative placeholders. Add
  the Burgess-grounded PT rule set to the scenario contract so later scenario
  rewrites use the same fundamentals as the GA rubric and prompts.
- [ ] lugag.9 Audit the full root scenario corpus and classify each scenario as
  PT-clean, boilerplate-but-salvageable, or materially misframed. Record which
  scenarios misuse actors, omit trust-building, omit reciprocal promises, rely
  on hidden authority, misuse non-PT vocabulary, or lack concrete thought
  experiment framing. Include a realism pass: find where scenarios use actors as
  generic placeholders instead of as agents with actual relationships, history,
  and reasons to trust or distrust one another.
- [ ] lugag.10 Define a stronger root scenario template for future rewrites.
  Require each scenario to include a realistic setup, prior relationship/trust
  history, explicit promises already in force, a concrete stimulus, at least one
  break/refusal/failure path, locally observable evidence, trust-update
  consequences, and scenario-specific open questions. Require the scenario to
  identify whether the interaction is about sending, receiving, storage,
  computation, or some combination. Require enough thought-experiment framing
  that a scorer can tell why Alice would trust Bob, distrust Bob, or demand
  reciprocal promises from Carol.
- [ ] lugag.11 Add new PT-focused scenarios that specifically exercise
  cross-legal-entity trust, selective sending, remote storage promises, remote
  computation promises, trust in returned results, reciprocal data exchange,
  promise-as-capability-token behavior, trust decay and repair, and multi-hop
  local trust without central authority.
- [ ] lugag.12 Rework or add sims so the scenario matrix has PT-clean candidate
  designs for selective sending, onward-restraint, storage promises,
  computation-result promises, reciprocal value exchange, and
  promise-as-capability-token framing. Treat capability tokens as promise
  evidence or promise references, never as bearer authority or imposed
  permission.
- [ ] lugag.13 Rework deferred child proposals `mipag` and `dafod`. Preserve
  only the parts that survive PT correction, such as explicit local evidence,
  durable pCID-referenced promise language, narrow layer boundaries, and
  long-horizon auditability. Remove or replace any framing that makes capsules,
  contracts, conformance, witness ledgers, or similar artifacts sound like
  authority-bearing mechanisms.
- [ ] lugag.14 Define review and promotion policy under the PT gate. A child or
  sim with PT-invalid framing must not be promotable. PT-reframe-needed artifacts
  may be kept only as question homes, anti-pattern specimens, or inputs for
  rework. PT-clean status becomes a hard prerequisite for any artifact that
  influences current design-state summaries.
- [ ] lugag.15 Plan the rescoring campaign. Freeze the corrected PT gate,
  generator contract, and scorer contract before spending on rescoring. Start
  with a small calibration sample, then re-score promoted or near-promoted sims,
  then hard-hit vocabulary families, then the broader corpus. Preserve
  historical scored JSON bytes and record new PT-corrected evidence as new
  result files rather than rewriting old artifacts.
- [ ] lugag.16 Define the rescore comparison outputs. Produce reports that show
  rank drift, formerly high scorers rejected by PT rules, newly strong PT-clean
  contenders, and scenario-driven shifts caused by more realistic trust and
  reciprocity framing.
- [ ] lugag.17 Update operator and guide-facing docs only after PT-corrected
  evidence exists. Future `tools/ga-runner/README.md`, promotion procedure, and
  `DEV-GUIDE-RESOURCES.md` updates must describe PT-first generation and scoring
  rules and must not summarize non-PT-clean results as consensus.

## Validation and acceptance criteria

- A later implementation pass must prove generator prompts require named agents,
  voluntary promises, reciprocal promises where relevant, evidence records,
  local trust judgment, and the Burgess-grounded PT rule set.
- A later implementation pass must prove the scorer rejects conventional
  access-control or contract-authority designs even when they are otherwise
  neat, that it is highly suspicious of non-voluntary cooperation and
  imposition-heavy framing, and that it penalizes storage/computation designs
  that presume service trust by default.
- A later scenario rewrite pass must leave the root scenario corpus materially
  less boilerplate, more realistic, and more explicit about trust-building,
  reciprocal promises, evidence, and thought-experiment framing.
- `mipag` and `dafod` remain deferred until PT-corrected review and rework are
  complete.
- Rescoring must preserve historical scored artifacts and append new PT-clean
  evidence instead of mutating old bytes in place.
