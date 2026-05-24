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
Source: `DI-safij`.

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

## Current problems

- GA child-generation prompts do not require the model to identify promising
  agents, reciprocal promises, make/break evidence, or what remains local trust
  judgment.
- The scorer still allows technically neat but conceptually wrong designs to
  score well even when they use permission, authorization, conformance,
  contract, or policy-enforcement framing as if those were external authority.
- PT principles are not being applied uniformly to send, receive, relay,
  storage, computation, result return, refusal, and trust updates.
- Many scenarios still read as boilerplate rather than thought experiments:
  they under-use Alice/Bob/Carol/Mallory roles, have weak setup/stimulus/trust
  history, and often drift into non-PT vocabulary such as `authority`,
  `conformance`, and `access control`.
- Existing high-scoring or near-promoted artifacts may need to be reinterpreted
  or demoted once PT-corrected scoring exists.

## Subtasks

- [ ] lugag.1 Lock the Promise Theory vocabulary and anti-vocabulary used by GA
  prompts, scoring, scenario docs, and promotion review. Define preferred terms
  such as promise, promiser, promisee, reciprocal promise, local trust, trust
  history, kept promise, broken promise, and pCID-referenced promise; define
  suspect terms such as permission, who may, authorization, conformance,
  contract, policy enforcement, compliance, access control, and authority that
  require explicit PT reframing before they are allowed.
- [ ] lugag.2 Write a compact PT interaction checklist that applies the same
  discipline to send, receive, relay, store, retrieve, compute, return-result,
  forward, refuse, ignore, delay, supersede, and revoke behaviors. For each
  interaction class, require explicit answers to: which agent promises what,
  what evidence of keeping/breaking looks like, how local trust is updated, and
  what no artifact can decide on the agent's behalf.
- [ ] lugag.3 Harden GA child-generation prompts so every generated design must
  identify the promising agents, what each agent voluntarily promises, which
  promises are reciprocal, what make/break evidence records matter, and what
  remains local trust judgment. Require explicit treatment of send, receive,
  storage, and computation when relevant, and require the child to explain how
  Alice decides whether to trust Bob enough to send data and whether to trust
  computation results Bob returns.
- [ ] lugag.4 Add an explicit Promise Theory gate before weighted scoring. The
  gate must classify designs as PT-clean, PT-reframe-needed, or PT-invalid.
  PT-invalid designs cannot be promotable regardless of technical neatness.
  PT-reframe-needed designs may survive only as question homes or rework
  candidates. PT-clean status is required for any design that can influence
  promotion, consensus, or guide-level summary.
- [ ] lugag.5 Harden the scorer rubric and penalties. Make
  `promise_vocabulary` close to a correctness axis rather than a soft style
  preference; make high `promisegrid_alignment` impossible unless the PT gate is
  clean; severely penalize external-authority language unless the scorer
  explains the local-promise interpretation; and penalize designs that treat
  storage or computation as neutral service calls rather than promises between
  agents with local trust consequences.
- [ ] lugag.6 Decide the PT-gate result shape and scoring-report contract. Lock
  whether the PT gate becomes an explicit structured field in score results, a
  required assessment subsection, or both. Ensure the promotion and compare
  tooling can key off the PT gate deterministically rather than inferring it
  from prose alone.
- [ ] lugag.7 Audit `scenarios/README.md` and strengthen the shared root
  scenario contract so every scenario inherits explicit PT requirements for
  reciprocal promises, local trust, send/receive/storage/computation discipline,
  observable make/break history, and trust updates. Tighten the actor guidance
  so Alice, Bob, Carol, Dave, Ellen, Frank, and Mallory are used as concrete
  cooperating or adversarial agents rather than as decorative placeholders.
- [ ] lugag.8 Audit the full root scenario corpus and classify each scenario as
  PT-clean, boilerplate-but-salvageable, or materially misframed. Record which
  scenarios misuse actors, omit trust-building, omit reciprocal promises, rely
  on hidden authority, misuse non-PT vocabulary, or lack concrete thought
  experiment framing.
- [ ] lugag.9 Define a stronger root scenario template for future rewrites.
  Require each scenario to include a realistic setup, prior relationship/trust
  history, explicit promises already in force, a concrete stimulus, at least one
  break/refusal/failure path, locally observable evidence, trust-update
  consequences, and scenario-specific open questions. Require the scenario to
  identify whether the interaction is about sending, receiving, storage,
  computation, or some combination.
- [ ] lugag.10 Add new PT-focused scenarios that specifically exercise
  cross-legal-entity trust, selective sending, remote storage promises, remote
  computation promises, trust in returned results, reciprocal data exchange,
  promise-as-capability-token behavior, trust decay and repair, and multi-hop
  local trust without central authority.
- [ ] lugag.11 Rework or add sims so the scenario matrix has PT-clean candidate
  designs for selective sending, onward-restraint, storage promises,
  computation-result promises, reciprocal value exchange, and
  promise-as-capability-token framing. Treat capability tokens as promise
  evidence or promise references, never as bearer authority or imposed
  permission.
- [ ] lugag.12 Rework deferred child proposals `mipag` and `dafod`. Preserve
  only the parts that survive PT correction, such as explicit local evidence,
  durable pCID-referenced promise language, narrow layer boundaries, and
  long-horizon auditability. Remove or replace any framing that makes capsules,
  contracts, conformance, witness ledgers, or similar artifacts sound like
  authority-bearing mechanisms.
- [ ] lugag.13 Define review and promotion policy under the PT gate. A child or
  sim with PT-invalid framing must not be promotable. PT-reframe-needed artifacts
  may be kept only as question homes, anti-pattern specimens, or inputs for
  rework. PT-clean status becomes a hard prerequisite for any artifact that
  influences current design-state summaries.
- [ ] lugag.14 Plan the rescoring campaign. Freeze the corrected PT gate,
  generator contract, and scorer contract before spending on rescoring. Start
  with a small calibration sample, then re-score promoted or near-promoted sims,
  then hard-hit vocabulary families, then the broader corpus. Preserve
  historical scored JSON bytes and record new PT-corrected evidence as new
  result files rather than rewriting old artifacts.
- [ ] lugag.15 Define the rescore comparison outputs. Produce reports that show
  rank drift, formerly high scorers rejected by PT rules, newly strong PT-clean
  contenders, and scenario-driven shifts caused by more realistic trust and
  reciprocity framing.
- [ ] lugag.16 Update operator and guide-facing docs only after PT-corrected
  evidence exists. Future `tools/ga-runner/README.md`, promotion procedure, and
  `DEV-GUIDE-RESOURCES.md` updates must describe PT-first generation and scoring
  rules and must not summarize non-PT-clean results as consensus.

## Validation and acceptance criteria

- A later implementation pass must prove generator prompts require named agents,
  voluntary promises, reciprocal promises where relevant, evidence records, and
  local trust judgment.
- A later implementation pass must prove the scorer rejects conventional
  access-control or contract-authority designs even when they are otherwise
  neat, and that it penalizes storage/computation designs that presume service
  trust by default.
- A later scenario rewrite pass must leave the root scenario corpus materially
  less boilerplate, more realistic, and more explicit about trust-building,
  reciprocal promises, evidence, and thought-experiment framing.
- `mipag` and `dafod` remain deferred until PT-corrected review and rework are
  complete.
- Rescoring must preserve historical scored artifacts and append new PT-clean
  evidence instead of mutating old bytes in place.

