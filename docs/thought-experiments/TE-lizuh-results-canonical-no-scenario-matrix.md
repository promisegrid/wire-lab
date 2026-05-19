# TE-lizuh: Results canonical, no committed scenario matrix

## TE ID

TE-lizuh

## Status

decided

First drafted 2026-05-18.

## Decision under test

DUT-lizuh: Should `scenarios/*/MATRIX.md` remain committed state, or should
`results/` be the only canonical record of simulation/scenario/model runs?

## Short answer

Delete committed scenario `MATRIX.md` files. Keep `results/<sim>/<scenario>/<model>/<timestamp>.md`
as the canonical evidence tree, and generate matrix views from `results/` when
humans or agents need scan tables. Source: `DI-zamin`.

This supersedes only the `MATRIX.md` scenario-layout clause in TE-dojab. The
TE-dojab decisions to keep root `scenarios/`, root per-run `results/`, and
evidence-only result authority remain active.

## Assumptions

- Result files already carry the full evidence payload: sim, scenario, model,
  timestamp, prompt/procedure, observed behavior, verdict, evidence links, open
  questions, and authority boundary.
- Scenario trees are input context; result trees are output evidence.
- A model prompt should not include previous result summaries for the same
  scenario, because that violates blind-run expectations.
- Generated matrix views are useful, but committed matrix snapshots are not
  authoritative.

## Alternatives

### Alt A - Keep committed `MATRIX.md`

Each scenario keeps a checked-in matrix of latest result links and short notes.

**Easier:** Humans can open a scenario directory and see recent evidence without
running a tool.

**Harder:** The matrix duplicates `results/`, carries placeholder `not-run`
rows, and can drift from the canonical result tree. It also tempts prompt
bundling to include previous result summaries as scenario context.

### Alt B - Delete committed matrices and generate views

`results/` is canonical. A runner command derives scenario/sim/model/latest
tables from result files on demand.

**Easier:** One source of truth, cleaner input/output boundary, and no matrix
drift. Validation checks result shape and source existence instead of requiring a
second committed link.

**Harder:** Humans need a command to see matrix-style summaries.

### Alt C - Keep matrices only as generated committed reports

Tools regenerate `MATRIX.md`, and humans commit the generated output.

**Easier:** Preserves browseability.

**Harder:** Still creates duplicate state and review churn. The generated files
are easy to mistake for input scenario context.

## Scenario analysis

### S1: Alice runs a new full matrix

With Alt A, every result write also mutates a scenario file, creating thousands
of scenario-side edits unrelated to scenario meaning. With Alt B, each cell
writes one result artifact and optional generated reports can be reviewed
separately. Alt B keeps the run audit cleaner.

### S2: Bob audits a single scenario

With Alt A, Bob sees a short summary quickly, but must still open result files
to audit evidence. With Alt B, Bob runs a generated view and then opens the same
result files. The only lost convenience is a pre-rendered table.

### S3: Carol prepares API prompt bundles

With Alt A, `MATRIX.md` sits inside the scenario tree but must be excluded from
the prompt to preserve blind evaluation. With Alt B, scenario directories are
input-only by construction, so prompt bundling is simpler and safer.

### S4: Dave validates repository consistency

With Alt A, validation has to prove result files and matrices agree. With Alt B,
validation proves result paths are well-shaped and that referenced sim/scenario
source files exist. The single-source validation is less fragile.

## Conclusion

Alt B is locked by `DI-zamin`: delete committed scenario matrices, stop updating
them during runs, and add generated result views instead. This is a source of
truth cleanup, not a loss of evidence.

## Implications for open TODOs, DRs, and DIs

- `TODO-dadub` owns the implementation because it owns root scenario/result
  apparatus.
- `tools/matrix-runner` should replace `update-matrix` with generated result
  views.
- `results/RUN-PROTOCOL.md`, `results/README.md`, `scenarios/README.md`, and
  legacy Python docs should stop treating `MATRIX.md` as current committed
  state.

## Decision status

`decided` — Steve chose Alt B on 2026-05-18. Locked by `DI-zamin`.
