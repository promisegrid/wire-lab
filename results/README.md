# Root Results

Root `results/` entries are wire-lab comparison evidence for root scenarios.
They are not PromiseGrid node layout, production API, final design authority, or
simulation-local world state. Source: `DI-faros`; `DI-miror`; `DI-dimas`.

## Path Shapes

Legacy Markdown matrix evidence uses this path:

```text
results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md
```

GA/search fitness evidence uses this path:

```text
results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.json
```

- `<sim-id>` is the exact simulation directory name without a trailing slash.
- `<scenario-id>` is the stable root scenario ID.
- `<model-id>` is a path-safe, specific provider/model/reasoning slug, such as
  `openai-GPT-5.5-xhigh`. A generic interface label such as `codex` is not a
  valid model ID.
- `<YYYYMMDD-HHMMSS>` is the UTC run timestamp.

GA/search run state lives at `results/state/<run-group-id>.json` with schema
`promisegrid.ga.state.v1`. Source: `DI-zanon`; `DI-ruzaj`.

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

`tools/ga-runner` is the preferred runner for GA/search work. It validates JSON
fitness results, ignores old Markdown canary files, creates stateful generation
runs, scores cells through provider-backed reasoning, keeps pending child sims
state-bound, records accepted children without staging or committing them, and
explicitly culls rejected children. Source: `DI-ramar`; `DI-zanon`;
`DI-podot`; `DI-kofil`; `DI-ruzaj`; `DI-gijom`.

`tools/matrix-runner` and the preserved Python scripts remain useful for legacy
Markdown matrix canaries and historical comparisons. They are not the GA/search
contract for JSON fitness evidence. Source: `DI-lulom`; `DI-ruzaj`.

## Run Preflight

Before kicking off a legacy Markdown matrix run:

1. Follow `results/RUN-PROTOCOL.md`.
2. Generate a deterministic manifest with `tools/matrix-runner`.
3. Run a small API-backed canary with `tools/matrix-runner run`.
4. Review `results/state/<run-group-id>.json` token/cost fields and tune
   `-result-style`, `-max-output-tokens`, `-max-run-cost-usd`, and
   `-max-cell-estimate-usd`.
5. Validate canary output with `tools/matrix-runner validate`.
6. Run the full manifest only after canary validation and cost review pass.
Source: `DI-nugiv`.

The old Python scripts under `results/tools/` remain legacy/reference tools.
For legacy Markdown matrix runs, the Go matrix-runner is preferred because it
bundles local source document contents, checkpoints state, validates results,
and generates result views from the canonical result tree. Source: `DI-lulom`;
`DI-zamin`; `DI-ruzaj`.

Before a GA/search run:

1. Follow the GA/search section in `results/RUN-PROTOCOL.md`.
2. Use `tools/ga-runner init -dry-run` to preview the tracked population and
   conservative scenario/child plan.
3. Use non-dry-run `tools/ga-runner init` to create the checkpoint state, then
   run `tools/ga-runner score` and `tools/ga-runner generate` with explicit
   provider model, reasoning effort, and cost-budget flags.
4. Keep generated child sims uncommitted until `tools/ga-runner accept` records
   reviewed promotion evidence.
5. Use `tools/ga-runner cull` to delete rejected child sims and matching result
   trees only through the state-bound cleanup path. Source: `DI-ramar`;
   `DI-zusit`; `DI-podot`; `DI-kofil`; `DI-ruzaj`; `DI-gijom`.
