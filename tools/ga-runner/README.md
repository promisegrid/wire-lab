# GA Runner

`tools/ga-runner` is the PromiseGrid Wire Lab runner for GA/search work. It is
separate from the legacy Markdown matrix-runner path. Historical GA evidence
remains append-only `promisegrid.ga.result.v1`, new scoring writes
`promisegrid.ga.result.v2`, child proposals and child score evidence stage
under `proposals/<run-group-id>/`, and run state remains JSON at
`results/state/<run-group-id>.json`. Source: `DI-ramar`; `DI-zanon`;
`DI-ruzaj`; `DI-lirat`; `DI-roruj`.

## Current Commands

Run commands from this directory and pass `-repo-root ../..` when needed:

```bash
go run . validate -repo-root ../..

  go run . init -repo-root ../.. -dry-run \
    -model openai-gpt-5.4-medium \
    -run-group-id <run-group-id>

	  go run . init -repo-root ../.. \
	    -model openai-gpt-5.4-medium \
	    -run-group-id <run-group-id> \
	    -include-sim <SIM-id> \
	    -include-scenario <scenario-id>

  go run . audit -repo-root ../.. \
    -clean-envelope-count 6

  go run . backfill-init -repo-root ../.. \
    -run-group-id <run-group-id> \
    -staged-model-id openai-gpt-5.4-high \
    -staged-reasoning-effort high \
    -clean-envelope-count 6

  go run . compare-backfill -repo-root ../.. \
    -run-group-id <run-group-id>

  go run . score -repo-root ../.. \
    -run-group-id <run-group-id> \
    -target parents \
    -reasoning-effort medium \
    -reasoning-summary auto \
    -text-verbosity low \
    -service-tier flex \
    -workers 6 \
    -request-timeout 5m \
    -provider-max-attempts 2 \
    -provider-max-elapsed 12m \
    -stream=true \
    -stream-idle-timeout 2m \
    -stream-content-stdout=true \
    -skip-failed-cells \
    -max-run-cost-usd <budget>

  go run . generate -repo-root ../.. \
    -run-group-id <run-group-id> \
    -api-model gpt-5.4 \
    -reasoning-effort medium \
    -reasoning-summary auto \
    -text-verbosity low \
    -service-tier flex \
    -workers 1 \
    -request-timeout 15m \
    -provider-max-attempts 2 \
    -provider-max-elapsed 12m \
    -stream=true \
    -stream-idle-timeout 2m \
    -stream-content-stdout=true \
    -skip-failed-children \
    -max-run-cost-usd <budget>

go run . accept -repo-root ../.. \
  -run-group-id <run-group-id> \
  -child <SIM-id> \
  -result proposals/<run-group-id>/results/<SIM-id>/<scenario-id>/<model-id>/<YYYYMMDD-HHMMSS>.json \
  -reviewer-note '<why this child should be promoted>'

go run . cull -repo-root ../.. \
  -run-group-id <run-group-id> \
  -child <SIM-id> \
  -reason '<why this child is rejected>'
```

`validate` reads only GA JSON fitness files and ignores old Markdown canary
results. `audit` inspects canonical historical `promisegrid.ga.result.v1`
results, reports exact-match vs root-contract drift, reports whether each sim
was audited from a live historical source, a canonical fallback, or an
unresolved source, and identifies hard-hit vocabulary families plus clean
grid-envelope calibration contenders. Promoted canonical results may keep their
historical `source.*` proposal paths; when those proposal trees are gone,
`audit` and `backfill-init` may compare against the current canonical
`simulations/<sim-id>/` tree, but only as an exact-byte fallback. If neither
source root exists, the record stays non-exact and is excluded from targeted
backfill instead of aborting the run. `backfill-init` turns that audit into a
targeted `promisegrid.ga.state.v1` state file for additive rubric-v2 rescoring.
`compare-backfill` turns a completed targeted v2 state into a reviewable
Markdown report under `results/reports/<run-group-id>-comparison.md`, comparing
each v2 cell against the latest exact-match canonical v1 result for the same
`sim_id` + `scenario_id`, preferring the same `runner.api_model` when
available. Same-model historical reruns collapse to the latest record and are
reported separately from true ambiguous historical pairings. The report
summarizes sim rank drift, family highlights, and the largest per-cell score
deltas before any broader rerun. Source: `DI-zuzup`; `DI-sirir`.
By default it preserves the historical source model lineage. When a new stage
should score under a different model ID or default effort, use
`-staged-model-id` and optionally `-staged-reasoning-effort` so the generated
state, cell `model_id`, result paths, derived `api_model`, and default
reasoning effort all match the new stage. That is the honest path for
multi-stage rescoring such as `high` triage followed by `xhigh` tie-breaks.
`init -dry-run` previews tracked
population and conservative generation sizing without writing state. Non-dry-run
`init` creates `promisegrid.ga.state.v1`. `init` can repeat `-include-sim` and
`-include-scenario` to guarantee focused coverage while filling remaining sample
slots by deterministic shuffle. `score` builds source-complete prompts, calls
the provider, writes validated JSON fitness results, and checkpoints usage/cost
metadata after each cell. `generate` builds compact child prompts from each
parent simulation tree once, scenario-specific pressure once, and summarized
fitness evidence, then uses completed parent fitness results to rank the
selected parent pool, applies deterministic linear-rank weighted high-parent
plus uniform random scored-parent selection, materializes strict file-bundle
responses under ignored review-stage
`proposals/<run-group-id>/simulations/<SIM-id>/` trees, and records
prompt/response hashes, file hashes, tree hashes, and cost metadata in state.
Matching child score evidence is written under
`proposals/<run-group-id>/results/<SIM-id>/` until a review/promotion pass gives
the child a final non-child simulation name and canonical result home.
`accept` records reviewed promotion evidence and prints proposal paths for the
later promotion workflow, but it does not run `git add` or commit. Follow
`PROMOTION.md` when Steve says `promote <child-proquint> ...`; under
`DI-higot`, scored child/result artifacts are append-only, so promotion must
not rewrite scored simulation files or scored JSON result bytes in place. A
future promotion may move or copy scored bytes unchanged into canonical homes,
but until that byte-identical path is locked accepted children may remain under
`proposals/` with review state and surrounding index/TODO provenance. `cull`
deletes only
state-selected generated child sim trees and matching result trees; use
`-dry-run` to print the deletion plan without changing files. Source:
`DI-pobus`; `DI-bagih`; `DI-zusit`; `DI-podot`; `DI-kofil`; `DI-ruzaj`;
`DI-gijom`; `DI-puhog`; `DI-dilaf`; `DI-fihof`; `DI-lirat`; `DI-dikoh`;
`DI-zadik`; `DI-higot`; `DI-roruj`; `DI-zobur`; `DI-hijub`; `DI-zuzup`.

Provider-backed `score` and `generate` always send an explicit service tier.
The default is `-service-tier flex`; `-service-tier default` is available when
standard processing is intentionally worth the cost, and `priority` is rejected.
OpenAI-compatible transient failures (`408`, `409`, `429`, `500`, `502`, `503`,
`504`, network timeouts, and request timeouts) retry with bounded exponential
backoff before the cell is checkpointed as failed. Source: `DI-mopob`;
`DI-tufud`.

Provider-backed `score` and `generate` are serial by default with `-workers 1`,
but sync canaries use bounded workers for higher throughput. The terminal
canary defaults to six scoring workers and one child-generation worker until a
successful generation phase provides evidence for parallel child writes. Scoring
uses `GA_CANARY_SCORE_REQUEST_TIMEOUT` with legacy `GA_CANARY_REQUEST_TIMEOUT`
as a fallback and a five-minute default; child generation uses
`GA_CANARY_GENERATE_REQUEST_TIMEOUT` with a fifteen-minute default. Each
cell/child retry window defaults to `-provider-max-attempts 2` and
`-provider-max-elapsed 12m`.
Set `GA_CANARY_INCLUDE_SIMS` and `GA_CANARY_INCLUDE_SCENARIOS` to comma- or
space-separated IDs when a canary must include a newly-added sim or a particular
scenario; the wrapper passes those through as repeatable `init` include flags.
`tools/ga-runner/run-canary.sh` also reads an optional `/tmp/canary-cells`
focus file with `sims:` / `scenarios:` sections, one selector per line,
ignores blank lines and `#` comments, resolves selectors by unique prefix
against current sim/scenario IDs, and merges the resolved IDs with any
`GA_CANARY_INCLUDE_*` values before the `init` step. The wrapper fails fast on
malformed, missing, or ambiguous focus-file selectors before any provider call.
Streaming is on by default with `-stream=true` and `-stream-idle-timeout 2m`, so
long Responses API calls log event progress and retry silent stalls. Concurrent
scoring reserves estimated cell cost before dispatch, so `-max-run-cost-usd`
remains a conservative launch budget rather than a best-effort warning. Source:
`DI-juzus`; `DI-tufud`; `DI-pivuj`; `DI-suzor`; `DI-guvif`; `DI-duzur`;
`DI-bataj`.

The terminal canary opts into `-reasoning-summary auto` and
`-stream-content-stdout=true` so stdout/logs show one no-newline progress dot
per `response.reasoning_summary_text.delta` event,
`response.reasoning_summary_part.done` event names and content while cells are
running. Delta dots are only liveness signals; reasoning-summary text-delta
content, text-delta event names, `response.reasoning_summary_part.added`
events, and `response.output_text.delta` events are not mirrored to the canary
content stream. Raw `score` and `generate` commands leave stdout stream content
off unless requested explicitly. Source: `DI-vadub`; `DI-babik`; `DI-vajut`;
`DI-sakam`; `DI-fupob`; `DI-ramun`.

Provider-backed `score` and `generate` omit `max_output_tokens` by default.
`-max-output-tokens` remains an explicit emergency fuse, but normal cost control
uses `-cost-estimate-output-tokens` only for preflight budget estimates. The
default provider text verbosity is `low`; the terminal canary scores with
medium reasoning by default and generates children with medium reasoning.
Escalate scoring to `xhigh` explicitly for tie-breaks and promotion-sensitive
comparisons. A provider response with
`status: incomplete` and `reason: max_output_tokens` is still treated as a
deterministic cap failure when an operator explicitly sets the cap. Source:
`DI-juzus`; `DI-zikag`; `DI-pulap`; `DI-nanor`.

OpenAI Structured Outputs remain a follow-up after the uncapped canary succeeds.
They may reduce JSON-shape retries and prompt boilerplate, but they add
provider-specific schema plumbing and do not directly reduce hidden reasoning
tokens. Source: `DI-pulap`.

Rubric-v2 scoring adds `promise_vocabulary` and `simplicity_durability` to the
historical axis set and recomputes `fitness.raw` / `fitness.normalized_0_100`
inside the runner with normal weighting so v2 comparisons do not depend on
provider-specific math. Historical v1 results remain valid evidence and are not
rewritten. Source: `DI-roruj`.

If `score` is run without `-api-model`, each selected cell uses the `api_model`
already stored in the state file, falling back to a provider-model derivation
from the durable `<model-id>` path when needed. This allows mixed-model
audit-first backfill states to preserve historical model lineage without
rewriting v1 evidence. Source: `DI-roruj`.

Child generation uses parent score evidence in two ways. First, `generate`
reranks the selected parent pool by average completed parent
`fitness.normalized_0_100` and rewrites queued child parent IDs as exactly two
distinct `breed` parents using a deterministic linear-rank weighted high parent
plus one deterministic uniform random scored non-high parent. The high-parent
weights are bounded as `n, n-1, ..., 1` over the ranked scored pool, so better
parents are favored without making the top parent mandatory. Second, the child
prompt tells the model to preserve parent strengths, repair weaknesses, reduce
risks, route open questions, and make bounded design deltas expected to improve
the same rubric score. If fewer than two viable parents exist, generation
records a failed or skipped child instead of creating a one-parent child.
Child-generation prompts do not embed complete parent result JSON; they include
compact result-path, score, fitness, rationale, strength, weakness, risk, and
open-question summaries to reduce timeout-prone prompt bulk while preserving the
feedback signal. Source: `DI-puhog`; `DI-sohus`; `DI-dilaf`.

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
Unreviewed children and their child score evidence live under ignored
`proposals/<run-group-id>/`. A promotion pass must first record acceptance and
promotion intent without rewriting scored artifact bytes. If a canonical
simulation/result home is needed, use only a byte-identical move/copy path or
leave the scored artifacts under `proposals/` until that path is locked.
Rejected children should pass through `cull`. The detailed operator procedure
is `PROMOTION.md`. Source: `DI-ramar`; `DI-zanon`; `DI-zohal`; `DI-podot`;
`DI-kofil`; `DI-ruzaj`; `DI-gijom`; `DI-fihof`; `DI-lirat`; `DI-dikoh`;
`DI-zadik`; `DI-higot`.

## Legacy Boundary

`tools/matrix-runner` and `results/tools/` remain historical/canary Markdown
matrix tooling. They are not the GA/search JSON fitness contract, and GA runner
validation must not use old `.md` canary files as fitness evidence. Source:
`DI-lulom`; `DI-ramar`; `DI-pobus`; `DI-ruzaj`.
