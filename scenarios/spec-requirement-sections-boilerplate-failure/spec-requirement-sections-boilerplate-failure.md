# Boilerplate failure

## Scenario ID

spec-requirement-sections-boilerplate-failure

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-ranib-spec-requirement-sections/SCENARIOS.md`
- Source simulation: `SIM-ranib-spec-requirement-sections/`
- Source row/title: Boilerplate failure
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-ranib-spec-requirement-sections/`.

## Applies To

- Primary source simulation: `SIM-ranib-spec-requirement-sections/`
- Other simulations that claim to address `spec-requirement-sections` or the same PromiseGrid
  pressure should record result runs under `results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`; generated views can summarize completed runs when needed.

## Actors

Use Alice, Bob, Carol, Dave, Ellen, Frank, and Mallory as needed when running
this scenario. Preserve any actor roles implied by the source row instead of
inventing a central coordinator.

## Setup

A spec includes all required sections but each is vague.

## Stimulus

Run the candidate simulation against this source test: Whether review/conformance gates can reject empty prose.

## Expected Pressure

Required sections need quality criteria or they become ceremony.

## Overarching Goal Checks

- 100-year durability: The result must explain whether this pressure remains
  understandable after tooling, organizations, keys, people, and infrastructure
  change.
- Sparse, partial knowledge: The run must not assume any peer has all CAS
  objects, all relationship records, or a complete global graph unless that
  assumption is explicitly the failure being tested.
- No central authority or registry: The run must avoid relying on a central pCID,
  identity, trust, naming, routing, currency, or governance authority unless the
  scenario is explicitly testing that dependency.
- Peer-local trust / promise accounting: The run should identify what each peer
  can observe and record locally after the stimulus.
- Adversarial or failure pressure: The run should preserve the source row's
  failure, refusal, corruption, ambiguity, scale, or migration pressure; add
  Mallory only when adversarial behavior makes that pressure clearer.
- Human and LLM auditability: The result should be readable enough for a later
  person or model to reconstruct what was promised, what happened, and why the
  evidence matters.
- Migration / evolution path: The run should expose whether protocol, object,
  key, policy, naming, or organizational evolution changes the answer.

## Evaluation Questions

- Which promises or protocol claims does the candidate simulation need in order
  to handle the setup and stimulus?
- What local observations can Alice, Bob, Carol, or another peer record after the
  run?
- Does the candidate design preserve the expected pressure without appealing to
  hidden global state or central authority?
- What DR, DI, frozen spec, TODO, TE, or guide handoff would this evidence inform?

## Result Runs

Result runs live under:

`results/<sim-id>/spec-requirement-sections-boilerplate-failure/<model-id>/<YYYYMMDD-HHMMSS>.md`

## Authority Boundary

This scenario and its result runs are evidence only. Design authority still
graduates through DR, DI, frozen spec, or PromiseGrid Development Guide handoff.
