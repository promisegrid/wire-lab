# Deterministic CBOR agreement

## Scenario ID

cas-object-model-deterministic-cbor-agreement

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jomag-cas-object-model/SCENARIOS.md`
- Source simulation: `SIM-jomag-cas-object-model/`
- Source row/title: Deterministic CBOR agreement
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jomag-cas-object-model/`.

## Applies To

- Primary source simulation: `SIM-jomag-cas-object-model/`
- Other simulations that claim to address `cas-object-model` or the same PromiseGrid
  pressure should record result runs under `results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`; generated views can summarize completed runs when needed.

## Actors

Use Alice, Bob, Carol, Dave, Ellen, Frank, and Mallory as needed when running
this scenario. Preserve any actor roles implied by the source row instead of
inventing a central coordinator.

## Setup

Alice and Bob encode the same pointer object with independent implementations.

## Stimulus

Run the candidate simulation against this source test: Whether map ordering, integer/string choices, tags, and byte-string boundaries produce identical bytes and therefore identical CIDs.

## Expected Pressure

TE-43 must lock a precise CBOR profile rather than saying "use CBOR" generically.

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

`results/<sim-id>/cas-object-model-deterministic-cbor-agreement/<model-id>/<YYYYMMDD-HHMMSS>.md`

## Authority Boundary

This scenario and its result runs are evidence only. Design authority still
graduates through DR, DI, frozen spec, or PromiseGrid Development Guide handoff.
