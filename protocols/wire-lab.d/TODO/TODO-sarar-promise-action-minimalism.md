# TODO-sarar: Promise action minimalism

## Decision Intent Log

ID: DI-fitav
Date: 2026-06-03 06:59:06
Status: superseded
Author: stevegt@t7a.org (Steve Traugott)
Decision: Future PromiseGrid protocol, simulation, POC, scoring, generation, and
guide work must not invent workflow-specific top-level action kinds by default.
The default semantic action surface is `promise`, `refusal`, and `observation`.
Workflow labels such as repair, offer, counteroffer, accept, route, introduce,
redeem, transfer, store, compute, authorize, dispatch, grant, register, and
enforce belong in pCID-defined payload semantics or local interpretation unless a
future TE/DI proves a distinct wire-level role.
Intent: PromiseGrid must stay Promise-Theory-first. Agents make their own
promises, refuse locally, and observe local evidence. Turning every application
situation into a protocol verb drifts toward RPC, command dispatch, permission
systems, and spurious service vocabularies.
Constraints: Treat POC7 through POC10 action names as historical executable
evidence, not as the forward naming pattern. Do not rewrite scored/generated
artifacts in place. New POCs, sims, prompts, and guide prose must use the
minimal action rule unless a scoped TE/DI explicitly overrides it.
Affects: `AGENTS.md`; `DEV-GUIDE-RESOURCES.md`;
`tools/ga-runner/generate.go`; `tools/ga-runner/score.go`;
`tools/ga-runner/README.md`;
`protocols/wire-lab.d/TODO/TODO-sarar-promise-action-minimalism.md`;
`protocols/wire-lab.d/TODO/TODO.md`.

ID: DI-mosoj
Date: 2026-06-03 07:17:00
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Supersede the `DI-fitav` three-action vocabulary with a stricter
single-action rule. Future PromiseGrid-facing protocol semantics default to one
top-level semantic act: `promise`. Observation is a promise that the promiser
observed something from its local vantage. Refusal is either absence of a
promise, a promise not to do something, or a promise that the agent does not
currently promise the requested behavior. Workflow labels remain pCID-defined
payload semantics or local interpretation unless a future TE/DI proves a
distinct wire-level role.
Intent: PromiseGrid should not split ordinary Promise Theory semantics into
separate RPC-like action verbs. A single promise act is the cleaner abstraction:
agents can promise future behavior, local non-action, local belief, local
evidence, and observed events; other agents still make their own trust
assessments.
Constraints: Treat POC7 through POC10 action names, including refusal and
observation labels, as historical executable evidence rather than the forward
pattern. Do not rewrite scored/generated artifacts in place. Future POCs, sims,
prompts, and guide prose must default to `promise` unless a scoped TE/DI
explicitly overrides it.
Affects: `AGENTS.md`; `DEV-GUIDE-RESOURCES.md`;
`tools/ga-runner/generate.go`; `tools/ga-runner/score.go`;
`tools/ga-runner/result.go`; `tools/ga-runner/README.md`;
`tools/ga-runner/ga_runner_test.go`;
`protocols/wire-lab.d/TODO/TODO-sarar-promise-action-minimalism.md`;
`protocols/wire-lab.d/TODO/TODO.md`.
Supersedes: DI-fitav

## Status

Implemented. Owns the repo-wide Promise Action Minimalism rule that keeps future
PromiseGrid work from turning workflow situations into RPC-like action verbs.
The current rule is the single-action `promise` surface from `DI-mosoj`, which
supersedes the earlier `DI-fitav` three-action wording.

## Scope

- Canonicalize `promise` as the only default future-facing top-level semantic
  action.
- Treat observation, refusal, repair, economics, routing, introduction, storage,
  compute, token, and TCP-link changes as pCID-defined payload semantics plus
  local trust/evidence interpretation unless a later TE/DI proves a distinct
  wire-level role.
- Preserve POC7 through POC10 names as historical evidence; do not rewrite old
  scored/generated artifacts in place.
- Update the GA generator and scorer so future child sims are generated and
  scored against this rule.

## Subtasks

- [x] sarar.1 Record the locked repository-wide action-minimalism decision.
- [x] sarar.2 Add the canonical rule to `AGENTS.md`.
- [x] sarar.3 Update `DEV-GUIDE-RESOURCES.md` current-state guidance.
- [x] sarar.4 Update GA child-generation guardrails.
- [x] sarar.5 Update GA rubric/scoring prompt guidance.
- [x] sarar.6 Update GA runner documentation.
- [x] sarar.7 Supersede the three-action wording with the single-action
  `promise` rule under `DI-mosoj`.
