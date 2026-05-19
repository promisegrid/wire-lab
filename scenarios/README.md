# Root Scenarios

Root `scenarios/` entries are wire-lab comparison apparatus. They describe
pressure that multiple simulations can be run against; they are not PromiseGrid
node layout, production API, shared protocol components, or simulation-local
world state. Source: `DI-faros`; `DI-vabor`; `DI-dimas`.

## Directory Shape

Each root scenario entry uses this shape:

```text
scenarios/
  <entry-id>/
    README.md
    MATRIX.md
    <scenario-id>.md
```

- `<entry-id>` is stable kebab-case.
- Application entries use one `scenarios/<application-id>/` directory per
  application, such as `scenarios/bgp-routing/`.
- Mined simulation rows use one root scenario entry per source row, transformed
  for cross-simulation comparison and linked back to the source
  `simulations/.../SCENARIOS.md` row.
- `MATRIX.md` summarizes which simulations have been run against the entry and
  links to result-run files under `results/`.

## Scenario Entry Template

Use this template for each scenario markdown file:

```markdown
# <Scenario Title>

## Scenario ID

<stable-kebab-case-id>

## Source / Provenance

- Source type: application seed | mined simulation row | new harness scenario
- Source path:
- Source row/title:
- Source DI / TODO / TE:

## Purpose

<What design pressure this scenario applies.>

## Applies To

<Which simulations, protocol specimens, or candidate design families should be
compared against this scenario.>

## Actors

Use Alice, Bob, Carol, Dave, Ellen, Frank, and Mallory where named actors help
make promises, failures, or adversarial pressure concrete.

## Setup

<Initial state, promises, artifacts, sites, policies, and sparse knowledge.>

## Stimulus

<The event or request that exercises the scenario.>

## Expected Pressure

<What the scenario should force a candidate design to explain.>

## Overarching Goal Checks

- 100-year durability:
- Sparse, partial knowledge:
- No central authority or registry:
- Peer-local trust / promise accounting:
- Adversarial or failure pressure:
- Human and LLM auditability:
- Migration / evolution path:

## Evaluation Questions

- <Question 1>
- <Question 2>
- <Question 3>

## Result Runs

Result runs live under:

`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`

## Authority Boundary

This scenario and its result runs are evidence only. Design authority still
graduates through DR, DI, frozen spec, or PromiseGrid Development Guide handoff.
```

## Scenario Quality Gates

Every root scenario must explicitly state how it exercises the 100-year
PromiseGrid goal and other overarching PromiseGrid goals. A scenario may be
small, but it must not be narrow in a way that accidentally assumes away the
conditions PromiseGrid is meant to survive. Source: `DI-botup`.

At minimum, each scenario should address:

- **100-year durability:** Does the scenario still make sense after tools,
  organizations, keys, people, and infrastructure have changed?
- **Sparse, partial knowledge:** Does it avoid assuming any peer has the whole
  graph, all CAS objects, or globally complete state?
- **No central authority or registry:** Does it avoid relying on a central pCID,
  identity, trust, naming, routing, currency, or governance authority unless the
  point of the scenario is to test that failure mode?
- **Peer-local trust / promise accounting:** Does it show what Alice, Bob, Carol,
  or another peer can observe and record locally?
- **Adversarial or failure pressure:** Does it include Mallory, corruption,
  refusal, stale data, partition, capture, default, or another failure mode when
  relevant?
- **Human and LLM auditability:** Can a later person or model understand what
  was promised, what happened, and why the result matters?
- **Migration / evolution path:** Does it expose what happens when protocols,
  keys, names, object shapes, policies, or organizations evolve?

If a gate is not relevant, the scenario should say why instead of omitting it.

## Matrix Template

Each scenario entry's `MATRIX.md` should start with:

```markdown
# <Scenario Entry> Matrix

## Authority Boundary

This matrix summarizes evidence. It does not declare a winning design by itself.

| Simulation | Scenario | Latest result run | Status | Notes |
|---|---|---|---|---|
| `<sim-id>` | `<scenario-id>` | `results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md` | not-run |  |
```

## Population Plan

`TODO-dadub` owns the population plan. The first mining pass created one root
entry per existing sim-local scenario row under `DI-nanih`. The remaining
population work is to create one root entry per application seed.
