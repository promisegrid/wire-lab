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
    -model openai-gpt-5.4-xhigh \
    -run-group-id <run-group-id>

  go run . init -repo-root ../.. \
    -model openai-gpt-5.4-xhigh \
    -run-group-id <run-group-id>

  go run . score -repo-root ../.. \
    -run-group-id <run-group-id> \
    -target parents \
    -api-model gpt-5.4 \
    -reasoning-effort xhigh \
    -service-tier flex \
    -workers 3 \
    -request-timeout 5m \
    -provider-max-attempts 2 \
    -provider-max-elapsed 6m \
    -max-output-tokens 12000 \
    -skip-failed-cells \
    -max-run-cost-usd <budget>

  go run . generate -repo-root ../.. \
    -run-group-id <run-group-id> \
    -api-model gpt-5.4 \
    -reasoning-effort xhigh \
    -service-tier flex \
    -workers 1 \
    -request-timeout 5m \
    -provider-max-attempts 2 \
    -provider-max-elapsed 6m \
    -max-output-tokens 16000 \
    -skip-failed-children \
    -max-run-cost-usd <budget>

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
sizing without writing state. Non-dry-run `init` creates
`promisegrid.ga.state.v1`. `score` builds source-complete prompts, calls the
provider, writes validated JSON fitness results, and checkpoints usage/cost
metadata after each cell. `generate` builds source-complete child prompts,
materializes strict file-bundle responses as untracked `simulations/SIM-*`
trees, and records prompt/response hashes, file hashes, tree hashes, and cost
metadata in state. `accept` records reviewed promotion evidence and prints paths
to stage, but it does not run `git add` or commit. `cull` deletes only
state-selected generated child sim trees and matching result trees; use
`-dry-run` to print the deletion plan without changing files. Source:
`DI-pobus`; `DI-bagih`; `DI-zusit`; `DI-podot`; `DI-kofil`; `DI-ruzaj`;
`DI-gijom`.

Provider-backed `score` and `generate` always send an explicit service tier.
The default is `-service-tier flex`; `-service-tier default` is available when
standard processing is intentionally worth the cost, and `priority` is rejected.
Flex `429` resource-unavailable responses and request timeouts retry with
bounded exponential backoff before the cell is checkpointed as failed. Source:
`DI-mopob`.

Provider-backed `score` and `generate` are serial by default with `-workers 1`,
but sync canaries may use bounded workers for higher throughput. Each provider
attempt defaults to `-request-timeout 5m`; each cell/child retry window defaults
to `-provider-max-attempts 2` and `-provider-max-elapsed 6m`. Concurrent scoring
reserves estimated cell cost before dispatch, so `-max-run-cost-usd` remains a
conservative launch budget rather than a best-effort warning. Source:
`DI-juzus`.

For xhigh canary runs, the wrapper passes explicit output caps with
`GA_CANARY_SCORE_MAX_OUTPUT_TOKENS` defaulting to `12000` and
`GA_CANARY_GENERATE_MAX_OUTPUT_TOKENS` defaulting to `16000`. A provider
response with `status: incomplete` and `reason: max_output_tokens` is treated as
a deterministic cap failure rather than retried with the same cap. Source:
`DI-juzus`; `DI-zikag`.

For unattended canary-style runs, `score -skip-failed-cells` and
`generate -skip-failed-children` preserve per-cell or per-child failure messages
as `skipped` state entries while allowing later phases to continue. The terminal
canary uses these flags so a single unusable provider response does not prevent
child generation or child scoring from exercising the full loop. Source:
`DI-zikag`.

## Planned Commands

These command names are reserved by the v1 contract but are not fully
implemented yet:

- `progress`: summarize state counts, costs, child status, and culling status.

Do not treat generated children as accepted merely because they exist on disk.
Accepted children must pass through `accept`; rejected children should pass
through `cull`. Source: `DI-ramar`; `DI-zanon`; `DI-zohal`; `DI-podot`;
`DI-kofil`; `DI-ruzaj`; `DI-gijom`.

## Legacy Boundary

`tools/matrix-runner` and `results/tools/` remain historical/canary Markdown
matrix tooling. They are not the GA/search JSON fitness contract, and GA runner
validation must not use old `.md` canary files as fitness evidence. Source:
`DI-lulom`; `DI-ramar`; `DI-pobus`; `DI-ruzaj`.
