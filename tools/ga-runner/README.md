# GA Runner

`tools/ga-runner` is the PromiseGrid Wire Lab runner for GA/search work. It is
separate from the legacy Markdown matrix-runner path. GA/search fitness evidence
is JSON at `results/<sim-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.json`,
and run state is JSON at `results/state/<run-group-id>.json`. Source:
`DI-ramar`; `DI-zanon`; `DI-ruzaj`.

## Current Commands

Run commands from this directory and pass `-repo-root ../..` when needed:

```bash
go run . validate -repo-root ../..

go run . init -repo-root ../.. -dry-run \
  -model openai-gpt-5.3-codex-xhigh \
  -run-group-id <run-group-id>

go run . accept -repo-root ../.. \
  -run-group-id <run-group-id> \
  -child <SIM-id> \
  -result results/<SIM-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.json \
  -reviewer-note '<why this child should be promoted>'

go run . cull -repo-root ../.. \
  -run-group-id <run-group-id> \
  -child <SIM-id> \
  -reason '<why this child is rejected>'
```

`validate` reads only GA JSON fitness files and ignores old Markdown canary
results. `init -dry-run` previews tracked population and conservative generation
sizing without writing state. `accept` records reviewed promotion evidence and
prints paths to stage, but it does not run `git add` or commit. `cull` deletes
only state-selected generated child sim trees and matching result trees; use
`-dry-run` to print the deletion plan without changing files. Source:
`DI-pobus`; `DI-bagih`; `DI-zusit`; `DI-podot`; `DI-kofil`; `DI-ruzaj`.

## Planned Commands

These command names are reserved by the v1 contract but are not fully
implemented yet:

- `init` without `-dry-run`: create `promisegrid.ga.state.v1`.
- `score`: run one model over GA cells and write JSON fitness evidence.
- `generate`: create untracked child simulation trees under `simulations/SIM-*`.
- `progress`: summarize state counts, costs, child status, and culling status.

Do not treat generated children as accepted merely because they exist on disk.
Accepted children must pass through `accept`; rejected children should pass
through `cull`. Source: `DI-ramar`; `DI-zanon`; `DI-zohal`; `DI-podot`;
`DI-kofil`; `DI-ruzaj`.

## Legacy Boundary

`tools/matrix-runner` and `results/tools/` remain historical/canary Markdown
matrix tooling. They are not the GA/search JSON fitness contract, and GA runner
validation must not use old `.md` canary files as fitness evidence. Source:
`DI-lulom`; `DI-ramar`; `DI-pobus`; `DI-ruzaj`.
