# Root Results

Root `results/` entries are wire-lab comparison evidence for root scenarios.
They are not PromiseGrid node layout, production API, final design authority, or
simulation-local world state. Source: `DI-faros`; `DI-miror`; `DI-dimas`.

## Path Shape

Every real result run must use this path:

```text
results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md
```

- `<sim-id>` is the exact simulation directory name without a trailing slash.
- `<scenario-id>` is the stable root scenario ID.
- `<model-id>` is a path-safe, specific provider/model/reasoning slug, such as
  `openai-GPT-5.5-xhigh`. A generic interface label such as `codex` is not a
  valid model ID.
- `<YYYYMMDD-HHMMSS>` is the UTC run timestamp.

Do not create placeholder run files. A result file exists only after a real
human, Codex, scripted, or other model run has produced observations.

## Result Run Template

Use this template for each real run:

```markdown
# Result: <sim-id> / <scenario-id> / <model-id> / <YYYYMMDD-HHMMSS>

## Result ID

<sim-id>-<scenario-id>-<model-id>-<YYYYMMDD-HHMMSS>

## Scenario

- Scenario ID:
- Scenario path:

## Simulation

- Simulation ID:
- Simulation path:
- Simulation commit:

## Runner

- Runner/interface:
- Model ID:
- Run timestamp UTC:
- Operator:

## Prompt / Procedure

<Prompt, procedure, command, or manual run notes sufficient to audit what was
asked and how the run was performed.>

## Observed Behavior

<What happened when the simulation was evaluated against the scenario.>

## Verdict

<Evidence-oriented verdict. Do not declare final design authority here.>

## Evidence Links

- Scenario:
- Simulation docs:
- Supporting artifacts:

## Open Questions

- <Question 1>
- <Question 2>

## Handoff Target

<DR, DI, frozen spec, TODO, TE, guide prose, or none.>

## Authority Boundary

This result is evidence only. It does not settle the design by itself; design
authority graduates through DR, DI, frozen spec, or PromiseGrid Development
Guide handoff.
```

## Current State

Root result runs may be committed only when produced by an LLM or human reasoner
following `results/RUN-PROTOCOL.md`. Scripted prototype files are retained only
as plumbing-test artifacts and are excluded from evidence comparisons by default.
Source: `DI-moduf`.

`results/` is the only canonical store for simulation/scenario/model run
evidence. Scenario-side summary files are not committed; generate inspection
views with `tools/matrix-runner view` when needed. Source: `DI-zamin`.

## Run Preflight

Before kicking off a large matrix run:

1. Follow `results/RUN-PROTOCOL.md`.
2. Generate a deterministic manifest with `tools/matrix-runner`.
3. Run a small API-backed canary with `tools/matrix-runner run`.
4. Validate canary output with `tools/matrix-runner validate`.
5. Run the full manifest only after canary validation passes.

The old Python scripts under `results/tools/` remain legacy/reference tools;
the Go runner is preferred for unattended API-backed runs because it bundles
local source document contents, checkpoints state, validates results, and
generates result views from the canonical result tree. Source: `DI-lulom`;
`DI-zamin`.
