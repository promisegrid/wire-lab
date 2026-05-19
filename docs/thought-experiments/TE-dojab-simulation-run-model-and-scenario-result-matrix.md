# TE-dojab: Simulation run model and scenario/result matrix

## TE ID

TE-dojab

## Status

decided, refined

First drafted 2026-05-19.

## Decision under test

DUT-dojab: How should wire-lab run simulations and compare simulation variants
against shared scenarios?

Steve's current lean is:

- root-level `scenarios/`;
- root-level
  `results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`;
- Codex-only runs for now;
- a 2D comparison matrix with scenarios on one axis and simulations on the
  other.

This TE tests whether that shape is enough, whether it conflicts with the
simulation-first apparatus/specimen boundary, and which pieces should remain
open for later multi-agent or coded-runner work.

## Short answer

The lean is basically right, with one constraint:

Root `scenarios/` and root `results/` should be treated as **wire-lab apparatus
for comparing simulations**, not as PromiseGrid world state and not as final
design authority.

Recommended first shape:

```text
scenarios/
  README.md
  <suite-id>/
    README.md
    <scenario-id>.md
    MATRIX.md

results/
  <sim-id>/
    <scenario-id>/
      <model-id>/
        <YYYYMMDD-HHMMSS>.md
```

Start with Codex-only runs. Do not build a multi-provider runner yet. Instead,
make each result path name the specific provider/model/reasoning configuration
and UTC run timestamp, while each result file still carries enough metadata that
later OpenAI, Anthropic, Perplexity, scripted, or human runs can be added without
changing the directory model. Source: `DI-miror`.

This keeps the first implementation small while preserving the important future
option: multiple agents can run the same `(scenario, sim)` cell later and record
separate model-specific run records or reviewer notes.

## Assumptions

- `simulations/` remains the home of candidate PromiseGrid design specimens and
  simulation-local world state.
- `protocols/wire-lab.d/` remains the harness apparatus home.
- Cross-sim comparison is a harness concern because it asks "which candidate
  performed better under the same pressure?", not "what does this sim contain?"
- Existing sim-local `SCENARIOS.md` files remain useful for a simulation's own
  local pressure tests.
- A root scenario suite is a reusable test packet, not a shared component that
  sims depend on internally.
- A result file is evidence. It does not itself settle a design decision. Any
  graduation still goes through DR, DI, frozen spec, or guide prose.

## Alternatives

### Alt A - Codex-only, root scenarios, root per-run results

Use Codex as the only runner for now. Put cross-sim scenario suites under root
`scenarios/`. Put comparison results under root
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`. Each result
file names the scenario, sim, sim commit, runner/interface, specific model ID,
UTC run timestamp, prompt or procedure, observations, verdict, and open
questions. Source: `DI-miror`.

**Easier:** Very small first step. The matrix is visible from root. Humans and
agents can compare all sims without walking every sim tree. It matches Steve's
mental model directly.

**Harder:** It creates a second result location alongside simulation-local
`results/` discussed in earlier TEs. The files must be clearly labeled as
harness comparison results, or readers may think they are simulation-owned world
state.

**Obligations:** Define a minimal result template and a rule that root `results/`
is apparatus evidence. Link each result back to the sim path and scenario path.

### Alt B - Codex-only, root scenarios, sim-local results

Use root `scenarios/`, but write each cell result inside the tested simulation,
for example `simulations/<sim>/results/<suite>/<scenario-id>.md`.

**Easier:** Keeps evidence about a sim with the sim. This is closest to the
earlier simulation contract in TE-nizor and TE-mupoz.

**Harder:** The cross-sim matrix becomes scattered. A reader comparing 24
grid-envelope variants has to chase 24 separate result trees or rely on a
generated index.

**Obligations:** Maintain a root matrix index anyway. That means the root still
becomes the comparison surface, even if the result files are sim-local.

### Alt C - Multi-agent runner from the start

Write code now that calls multiple agents or providers, such as OpenAI,
Anthropic, Perplexity, scripted local agents, and possibly human-reviewed
runners. Each provider runs the same scenario/sim cell and the result matrix
compares both simulation behavior and runner disagreement.

**Easier:** Produces richer evidence and catches model-specific reasoning
failures earlier. This is attractive after the Perplexity/Claude degradation
concern.

**Harder:** It solves the wrong first problem. The result schema, prompt shape,
and scenario granularity are not stable yet. A multi-provider harness would add
credential handling, rate limits, model-version drift, cost, retries, and
provenance before the comparison format has earned its shape.

**Obligations:** Provider adapters, credential policy, prompt archives,
model-version capture, deterministic replay limits, and disagreement handling.

### Alt D - Coded deterministic runner first

Write a program that executes scenario fixtures against structured sim metadata
and produces result files automatically.

**Easier:** Repeatable, cheaper, and suitable for CI once sims have machine
readable specs and expected outputs.

**Harder:** Most current simulations are design-point documents, not executable
models. Forcing code too early will either flatten the design questions or
produce a weak parser over prose.

**Obligations:** A schema for sims, scenarios, verdicts, and assertions. This is
valuable later, but too early for the current corpus.

### Alt E - Scenario matrix only, no per-cell files

Keep one root `MATRIX.md` per suite. Each cell contains a short verdict and
links only to sim docs or notes.

**Easier:** Minimal file count. Easy to scan.

**Harder:** Cells become too small to preserve evidence. The matrix will either
turn into a huge prose document or silently compress away the reasoning that
future agents need to audit.

**Obligations:** Later extraction into result files when evidence grows. This is
almost guaranteed churn.

## Scenario analysis

### S1: Alice compares 24 grid-envelope variants

Alice wants to know how every grid-envelope variant behaves under the same
unknown-pCID and signature scenarios.

- **Alt A:** Alice opens one root suite matrix and each cell result is at a
  predictable path. The layout matches the comparison task.
- **Alt B:** Alice can find results, but must traverse many sim-local trees.
- **Alt C:** Alice gets richer model-disagreement evidence, but runner
  complexity distracts from the variant comparison.
- **Alt D:** Not ready unless the grid-envelope variants have structured
  executable specs.
- **Alt E:** The matrix is readable at first, then becomes too dense.

S1 favors Alt A.

### S2: Bob evolves one simulation independently

Bob changes one sim's design after a DR narrows the question. He wants to rerun
the shared scenarios without making that sim depend on a shared test fixture as
part of its own world.

- **Alt A:** Bob can rerun only that sim's result cells. The sim stays
  standalone; root results point to the sim commit tested.
- **Alt B:** Bob's sim contains its result history, which is convenient, but
  cross-sim comparison is less direct.
- **Alt C:** Multi-agent runs add cost for every sim iteration.
- **Alt D:** Good later if the sim has machine-readable claims.
- **Alt E:** Bob has nowhere to put detailed rerun evidence.

S2 favors Alt A or Alt B. The choice turns on comparison ergonomics; Alt A wins
if cross-sim comparison is the main job.

### S3: Carol wants OpenAI and Anthropic to judge the same cell

Carol suspects Codex is missing a failure mode. She wants Anthropic or another
agent to run the same scenario/sim cell independently.

- **Alt A:** Works because multiple model-specific run files can coexist under
  the same simulation/scenario path. The first runner can be Codex without
  baking the generic `codex` interface name into the model path.
- **Alt B:** Same, but evidence is scattered under sims.
- **Alt C:** Best long-term, but too much up-front machinery.
- **Alt D:** Does not answer LLM judgment disagreement.
- **Alt E:** Too little room for disagreement evidence.

S3 says Alt A is acceptable only if the result schema is runner- and
model-extensible from day one.

### S4: Dave reruns an old cell after a sim changes

Dave reruns scenario `S3` against `SIM-jurar` after the sim's CAS migration
question changes.

- **Alt A:** A per-model, per-UTC-timestamp run file avoids overwriting old
  evidence while keeping each run directly addressable.
- **Alt B:** Same risk, but local to the sim.
- **Alt C:** Runner identity makes the risk more obvious but not solved.
- **Alt D:** Deterministic CI can keep many run artifacts cheaply later.
- **Alt E:** The matrix cannot preserve rerun history well.

S4 refines Alt A: the first implementation should use one file per run, with the
run identified by the specific model slug and UTC timestamp in the path. Source:
`DI-miror`.

### S5: Ellen writes guide prose

Ellen wants to know which claims are stable enough for the PromiseGrid
Development Guide.

- **Alt A:** Root results are easy to cite as wire-lab evidence, but must say
  they are not normative.
- **Alt B:** Sim-local results are closer to the specimen but harder to survey.
- **Alt C:** Multi-agent disagreement could be useful, but only after result
  fields are comparable.
- **Alt D:** Strong once claims become executable.
- **Alt E:** Too compressed for guide provenance.

S5 favors Alt A with a clear graduation field: "handoff target: DR/DI/spec/guide
or none."

### S6: Mallory makes a result look authoritative

Mallory adds a result file saying a sim "wins" and tries to make readers treat
that as a design decision.

- **Alt A:** Root `results/` is visible, so it needs strong status language:
  result evidence is not a DI.
- **Alt B:** Sim-local placement makes the result look less global, but the same
  risk exists.
- **Alt C:** Multiple agents reduce single-run overclaiming but do not create
  authority.
- **Alt D:** Deterministic output still needs decision handoff.
- **Alt E:** Matrix-only verdicts are easiest to overclaim.

S6 requires every result file and matrix to state the authority boundary:
results inform DR/DI/spec work; they do not settle it.

### S7: Frank wants CI later

Frank wants a future CI job to detect when a sim stops satisfying a scenario.

- **Alt A:** Fine if scenario files and result metadata are structured enough to
  evolve. The first Markdown files can include a small front matter block later.
- **Alt B:** Also fine, but root discovery is harder.
- **Alt C:** Multi-agent CI is costly and nondeterministic.
- **Alt D:** Best eventual target for machine-checkable scenarios.
- **Alt E:** Weak foundation for automation.

S7 says Alt A should not preclude later structured scenario metadata and coded
runners.

## Recommended first contract

Adopt Alt A as the first implementation, with explicit future-proofing.

Root scenario suites:

```text
scenarios/
  README.md
  <suite-id>/
    README.md
    <scenario-id>.md
    MATRIX.md
```

Root result runs:

```text
results/
  <sim-id>/
    <scenario-id>/
      <model-id>/
        <YYYYMMDD-HHMMSS>.md
```

Minimal scenario fields:

- `Scenario ID`
- `Suite`
- `Purpose`
- `Applies to`
- `Actors`
- `Setup`
- `Stimulus`
- `Expected pressure`
- `Evaluation questions`
- `Authority boundary`

Minimal result fields:

- `Result ID`
- `Scenario`
- `Simulation`
- `Simulation commit`
- `Runner/interface`
- `Model ID`
- `Run timestamp UTC`
- `Prompt/procedure`
- `Observed behavior`
- `Verdict`
- `Evidence links`
- `Open questions`
- `Handoff target`
- `Authority boundary`

For now, `Runner/interface` should be `codex`. The path's `<model-id>` must name
the specific provider/model/reasoning configuration, for example
`openai-GPT-5.5-xhigh`; `codex` alone is not specific enough. Later runner
values can include `human`, `scripted`, `openai-agent`, `anthropic-agent`,
`perplexity-agent`, or `multi-agent-panel` without changing the root layout.
Source: `DI-miror`.

## What this TE rejects for now

- Writing a multi-provider agent harness before the scenario/result schema is
  tested by hand.
- Treating root `results/` as final design authority.
- Relying on a matrix-only record without per-cell evidence.
- Forcing every current simulation into machine-readable executable form before
  the prose design questions have stabilized.

## Open risks

- Root `results/` partially departs from earlier language that expected
  simulation-local `results/`. This is acceptable only if root `results/` is
  explicitly harness comparison apparatus, while sim-local results remain
  available for simulation-owned internal evidence if needed.
- One file per model/timestamp run may create many files after repeated reruns
  or multi-agent comparisons. Suite matrices should summarize or point to the
  relevant run files rather than duplicating their evidence.
- Root `scenarios/` can become too generic. Each suite should declare which sims
  it applies to and why.

## DF questions exposed

### DF-dojab.1 - Where do shared scenario suites live?

Locked answer: **1.A - Root `scenarios/`**. Source: `DI-faros`.

Surviving alternatives:

- **1.A - Root `scenarios/` (recommended).** Shared comparison apparatus at repo
  root.
- **1.B - `protocols/wire-lab.d/scenario-suites/`.** Harness-protocol-local
  apparatus, less visible but stricter.
- **1.C - Sim-local only.** Rejected for cross-sim comparison because it weakens
  apples-to-apples reuse.

### DF-dojab.2 - Where do cross-sim result runs live?

Locked answer: **2.A - Root per-run `results/`**:
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`, with UTC
timestamps and provider-prefixed specific model IDs. Source: `DI-miror`;
`DI-faros`.

Surviving alternatives:

- **2.A - Root per-run `results/` (recommended).** Best for matrix comparison,
  rerun history, and model-specific evidence, with clear authority-boundary
  language.
- **2.B - Sim-local `simulations/<sim>/results/`.** Better locality, weaker
  comparison ergonomics.
- **2.C - Matrix-only.** Rejected as too lossy.

### DF-dojab.3 - What runner should be used first?

Locked answer: **3.A - Codex-only first**. Source: `DI-faros`.

Surviving alternatives:

- **3.A - Codex-only first (recommended).** Establish scenario and result shape
  before adding runner complexity.
- **3.B - Multi-agent from the start.** Useful later, premature now.
- **3.C - Coded deterministic runner first.** Useful after scenarios become
  structured enough to execute.

### DF-dojab.4 - How should reruns and later multi-agent runs be represented?

Locked answer: **4.A - One run file per model/timestamp from day one**:
one file per specific model and UTC run timestamp from day
one, at `results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`.
Source: `DI-miror`; `DI-faros`.

Surviving alternatives:

- **4.A - One run file per model/timestamp from day one (recommended).** More
  scalable and preserves rerun evidence without append-section ambiguity.
- **4.B - Append-only cell file first.** Smaller file count but makes reruns and
  model comparisons harder to address directly.
- **4.C - Overwrite latest result.** Rejected because it destroys evidence.

### DF-dojab.5 - What authority do results have?

Locked answer: **5.A - Evidence only**; results graduate through DR/DI/spec/guide
handoff. Source: `DI-faros`.

Surviving alternatives:

- **5.A - Evidence only (recommended).** Result files inform decisions but do
  not settle them.
- **5.B - Results can directly mark a sim as winning.** Rejected because it
  bypasses decision provenance.
- **5.C - Results are illustrative only.** Too weak; simulations need to produce
  actionable evidence.

## Implications for open TODOs, DRs, and DIs

- A follow-on TODO should define the first `scenarios/` and `results/` skeleton
  now that `DF-dojab.1` through `DF-dojab.5` are locked by `DI-faros`.
- `simulations/README.md` should later point readers to root scenario suites and
  root result matrices after the root skeleton exists.
- `DEV-GUIDE-RESOURCES.md` should treat result files as wire-lab evidence, not
  as final PromiseGrid guide prose.
- Existing sim-local `SCENARIOS.md` files remain valid as local pressure tests;
  cross-sim scenario suites should not overwrite them.

## Decision status

`decided, refined; MATRIX.md clause superseded by TE-lizuh / DI-zamin`.
`DF-dojab.1` through `DF-dojab.5` are locked by `DI-faros`: root `scenarios/`,
root `results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`,
Codex-only first runs, one result file per provider/model/reasoning slug and UTC
timestamp, and evidence-only authority. The result-path refinement remains
recorded by `DI-miror`. The scenario `MATRIX.md` layout detail in this TE's
recommended first contract is historical only; current apparatus treats
`results/` as canonical and derives matrix views from result files. Source:
`DI-zamin`.

## Refinements

### 2026-05-19 - Model-specific per-run result paths

Steve refined the result-tree recommendation from one result file per
simulation/scenario cell to one file per simulation, scenario, specific
provider/model/reasoning slug, and UTC timestamp:
`results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.md`. The generic
runner/interface label such as `codex` remains result metadata and must not be
used as the model ID. Source: `DI-miror`.

### 2026-05-19 - DF-dojab closure

Steve answered all five DF-dojab choices with the recommended alternatives:
root `scenarios/`, root per-run `results/`, Codex-only first execution, one
file per model/timestamp run from day one, and evidence-only result authority.
Those choices are locked by `DI-faros`; `DI-miror` remains the source for the
specific model-ID and UTC timestamp path refinement.

### 2026-05-18 - Scenario MATRIX.md committed state removed

TE-lizuh supersedes this TE's scenario `MATRIX.md` layout detail. Root
`scenarios/` and root per-run `results/` remain active, but committed
`scenarios/*/MATRIX.md` files are no longer source-of-truth artifacts. Matrix
views are generated from `results/` when needed. Source: `DI-zamin`.
