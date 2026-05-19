# Tuition Attendance

## Scenario ID

tuition-attendance

## Source / Provenance

- Source type: application seed
- Source path: `protocols/wire-lab.d/TODO/TODO-dadub-root-scenario-skeleton-and-seed-catalog.md`
- Source row/title: Seed application catalog entry `tuition-attendance`
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-midif`; `TODO-dadub`

## Purpose

Exercise PromiseGrid design candidates against tuition-attendance application pressure:
Tuition payment, class enrollment, attendance evidence, refunds, and credential
eligibility.

## Applies To

- Any simulation that claims to support real application pressure rather than only
  protocol-internal mechanics.
- Result runs are recorded under `results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`; generated views can summarize completed runs when needed.

## Actors

- Alice depends on the application outcome.
- Bob makes or relays promises in the application workflow.
- Carol audits, verifies, or relies on the promised outcome.
- Mallory represents adversarial, captured, stale, or misleading evidence when
  that pressure is relevant to the run.

## Setup

Alice depends on an outcome in the Tuition Attendance domain. Bob makes promises about
tuition payment, class enrollment, attendance evidence, refunds, and credential
eligibility. Carol needs enough evidence to rely on Bob's promise without having
complete global state, and Mallory may exploit stale, missing, or ambiguous evidence.

## Stimulus

A routine application event becomes contested or incomplete. Some relevant objects,
signatures, observations, or relationship records are delayed, partitioned, stale, or
disputed, and each peer must decide what to accept, retry, downgrade, or escalate using
only local evidence.

## Expected Pressure

The candidate simulation must show which promises, CAS objects, feeds, identity claims,
names, and promise accounting records are needed for this application pressure, while
avoiding hidden global state or a central authority that would make the result non-
comparable.

## Overarching Goal Checks

- 100-year durability: The Tuition Attendance pressure should still be understandable after tools,
organizations, people, keys, policies, and infrastructure change over a century.
- Sparse, partial knowledge: Alice, Bob, Carol, and Mallory should each be able to act with partial state; no run
should assume one peer sees the whole graph or every CAS object.
- No central authority or registry: The run should not depend on a central naming, identity, currency, routing, governance,
or trust registry unless the run is explicitly testing that dependency as a failure
mode.
- Peer-local trust / promise accounting: The run should identify what each peer can observe and record locally about kept,
broken, refused, superseded, or ambiguous promises.
- Adversarial or failure pressure: The run should preserve contested, stale, missing, corrupt, captured, or adversarial
evidence pressure instead of reducing the application to a happy-path workflow.
- Human and LLM auditability: A later person or model should be able to reconstruct what was promised, what evidence
existed, what was missing, and why the result matters.
- Migration / evolution path: The run should expose how the answer changes when protocols, object shapes, keys, names,
organizations, or policies evolve.

## Evaluation Questions

- Which PromiseGrid promises and artifacts are needed to model this application
  pressure without collapsing it into a centralized workflow system?
- What can Alice, Bob, Carol, or Mallory observe and record locally after the
  stimulus?
- Which assumptions would make this application seed fail the 100-year,
  sparse-knowledge, or no-central-authority goals?
- What DR, DI, frozen spec, TODO, TE, or guide handoff would evidence from this
  run inform?

## Result Runs

Result runs live under:

`results/<sim-id>/tuition-attendance/<model-id>/<YYYYMMDD-HHMMSS>.md`

## Authority Boundary

This scenario and its result runs are evidence only. Design authority still
graduates through DR, DI, frozen spec, or PromiseGrid Development Guide handoff.
