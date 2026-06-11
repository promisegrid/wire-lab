# TODO-gogug: POC13 Hardening and POC14 Superset Plan

## Decision Intent Log

ID: DI-sihuz
Date: 2026-06-11 09:40:58
Status: active
Decision: Harden POC13 around the latest clean-run findings and plan POC14 as the next explicit superset baseline.
Intent: The latest POC13 run passed, but its monitor and analyzer evidence exposed gaps that should become executable regression checks before the next POC: non-commitment evidence must be counted consistently, duplicate evidence must not hide refusals or cache misses, trust must stay on a bounded local scale, recovery caution must be analyzer-visible, live-agent autonomy should choose useful promise work without weakening protocol safety, and dynamic TCP relationships should be tested as actual send/receive reachability. POC14 should begin from POC13 as a regression baseline rather than accidentally omitting prior POC11/POC12/POC13 behavior.
Constraints: Keep the single top-level `promise` action; do not add RPC verbs, global trust, central routing authority, or permission/conformance language; preserve pCID-defined payload semantics; keep local trust per-agent and evidence-based; keep run state scoped to the current clean-run root; do not persist POC state across clean runs.
Affects: implementations/poc13-cas-compute-functions/; protocols/wire-lab.d/TODO/TODO.md; DEV-GUIDE-RESOURCES.md.

ID: DI-sifot
Date: 2026-06-11 14:35:43
Status: active
Decision: The POC14 plan must include heterogeneous process-boundary agents: one or more WASM agents in their own process and one or more agents that do all messaging through stdio.
Intent: POC14 should test PromiseGrid's portability and kernel/app boundary model beyond same-language TCP-only Go processes. WASM agents test sandboxed runtimes and host-call boundaries; stdio-only agents test whether the same promise envelope and local-kernel routing model can work for subprocesses that have no direct network API and communicate only through standard input/output.
Constraints: Keep the single top-level `promise` action; keep app trust and workflow judgment outside the kernel; do not turn stdio or WASM adapters into RPC command surfaces; preserve pCID-defined payload semantics and exact envelope evidence.
Affects: implementations/poc14-production-progress/docs/POC14-SUPERSET-PLAN.md.

ID: DI-fimoh
Date: 2026-06-11 14:54:40
Status: active
Decision: Rename the POC14 planning directory to `implementations/poc14-wasm/`.
Intent: The POC14 plan now has a concrete WASM/stdio portability emphasis, so the directory slug should name that axis directly instead of the broader production-progress label.
Constraints: Preserve the POC14 superset baseline from `DI-sihuz`; preserve the heterogeneous process-boundary requirements from `DI-sifot`; move the plan rather than duplicating it.
Affects: implementations/poc14-wasm/docs/POC14-SUPERSET-PLAN.md.

## Tasks

- [x] gogug.1 Fix POC13 evidence summary mismatch so saved evidence counts include all local non-commitment outcomes, not only receiver-side `not_promised` journal entries.
- [x] gogug.2 Separate duplicate evidence from non-commitment in promise resolution for cache misses, refusals, replay refusals, future-only repair, unsupported variants, and duplicate shipment checkpoints.
- [x] gogug.3 Decide and implement bounded local trust-scale saturation so trust values stay comparable across runs.
- [x] gogug.4 Add analyzer gates for `DI-fijov` trust-caution behavior.
- [x] gogug.5 Add latest clean-run narrative documentation for `poc13-demo`.
- [x] gogug.6 Add a production-fitness report derived from analyzer and monitor output.
- [x] gogug.7 Add tests for malformed or unsupported live-agent promises after trust caution.
- [x] gogug.8 Improve live-agent autonomy prompts and evidence so agents choose useful pCID-scoped promise work without losing protocol safety.
- [x] gogug.9 Add a true dynamic TCP topology experiment where direct links affect real send/receive reachability during the run.
- [x] gogug.10 Start POC14 planning as the next production-progress superset while keeping POC13 as the regression baseline.
