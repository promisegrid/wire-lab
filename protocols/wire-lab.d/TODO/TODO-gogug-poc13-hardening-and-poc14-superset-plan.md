# TODO-gogug: POC13 Hardening and POC14 Superset Plan

## Decision Intent Log

ID: DI-sihuz
Date: 2026-06-11 09:40:58
Status: active
Decision: Harden POC13 around the latest clean-run findings and plan POC14 as the next explicit superset baseline.
Intent: The latest POC13 run passed, but its monitor and analyzer evidence exposed gaps that should become executable regression checks before the next POC: non-commitment evidence must be counted consistently, duplicate evidence must not hide refusals or cache misses, trust must stay on a bounded local scale, recovery caution must be analyzer-visible, live-agent autonomy should choose useful promise work without weakening protocol safety, and dynamic TCP relationships should be tested as actual send/receive reachability. POC14 should begin from POC13 as a regression baseline rather than accidentally omitting prior POC11/POC12/POC13 behavior.
Constraints: Keep the single top-level `promise` action; do not add RPC verbs, global trust, central routing authority, or permission/conformance language; preserve pCID-defined payload semantics; keep local trust per-agent and evidence-based; keep run state scoped to the current clean-run root; do not persist POC state across clean runs.
Affects: implementations/poc13-cas-compute-functions/; protocols/wire-lab.d/TODO/TODO.md; DEV-GUIDE-RESOURCES.md.

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
