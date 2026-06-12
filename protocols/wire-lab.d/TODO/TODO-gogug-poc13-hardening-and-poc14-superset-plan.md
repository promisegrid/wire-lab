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

ID: DI-lulof
Date: 2026-06-11 15:03:55
Status: active
Decision: POC14 monitoring and analysis planning must assume production has no global observer and no global analyzer.
Intent: Real PromiseGrid agents may run in geographically diverse runtimes controlled by different legal entities, so each agent only has local evidence and voluntary peer disclosures. POC14 should therefore test decentralized monitoring signals such as local evidence exchange, peer-carried attestations, reciprocal promise histories, and bearer-token exchange rates that may reveal local market estimates of relative trustworthiness without creating a monitoring authority.
Constraints: Do not introduce global monitoring, global analysis, centralized trust scores, central exchange rates, or cross-agent authority. Treat analyzer and monitor components as POC/local-development tools unless explicitly reframed as ordinary agents making voluntary promises from their own local evidence.
Affects: implementations/poc14-wasm/docs/POC14-SUPERSET-PLAN.md.

ID: DI-linof
Date: 2026-06-11 16:07:05
Status: active
Decision: Implement POC14 by copying POC13 into `implementations/poc14-wasm/`, renaming POC-local commands and module paths to POC14, and adding deterministic WASM-boundary and stdio-boundary agent roles without dropping POC11/POC12/POC13 behavior.
Intent: POC14 needs executable evidence for heterogeneous process boundaries while remaining a strict superset baseline. The WASM role should exercise a separate process that treats sandboxed module bytes as locally hosted agent logic; the stdio role should exercise a separate worker process whose only application messaging path is stdin/stdout through a local adapter. Both roles must exchange ordinary PromiseGrid envelopes and record local evidence rather than introducing RPC commands or global analysis.
Constraints: Preserve one top-level semantic action `promise`; keep trust/workflow judgment in apps rather than the kernel; keep analyzer/monitor global views as POC-only development tools; preserve pCID-defined envelope semantics; use POC14 command names derived from existing POC13 command names plus `poc14-wasm-agent`, `poc14-stdio-adapter`, and `poc14-stdio-worker`; keep runtime state under the POC14 run root.
Affects: implementations/poc14-wasm/; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md.

ID: DI-kinaf
Date: 2026-06-11 17:16:42
Status: active
Decision: Add two POC14 regression scenarios: Alice locally decides to permanently distrust Mallory, and Alice locally promises that neither inbound nor outbound traffic should transit Mallory.
Intent: POC14 should exercise stronger local trust boundaries than temporary caution. Permanent distrust must remain Alice's local relationship decision, not punishment, enforcement, or global reputation. Transit exclusion must be Alice's local routing/peering promise about Alice's own traffic, not an imposed network-wide ban on Mallory.
Constraints: Preserve one top-level semantic action `promise`; do not add RPC verbs, permission, authorization, global policy, or central route enforcement; keep the scenario as local evidence and analyzer gates until a later POC implements real multi-hop route selection.
Affects: implementations/poc14-wasm/runtime/; implementations/poc14-wasm/cmd/poc14-analyze/; implementations/poc14-wasm/docs/; implementations/poc14-wasm/README.md.

ID: DI-dubih
Date: 2026-06-11 17:32:00
Status: active
Decision: Make the POC14 hard trust-boundary scenarios behavioral by adding app-local permanent distrust state and app-local transit exclusion state to the relationship ledger and send gate.
Intent: Permanent distrust should prevent future direct sends and narrow candidate-peer sends from Alice to Mallory until Alice makes a separate local state change. Transit exclusion should reject any route candidate for Alice-owned traffic that contains Mallory as a transit hop, while still treating the decision as Alice's local promise rather than global route enforcement.
Constraints: Keep the behavior local to the observing app; do not create a global ban, network authority, permission system, or RPC route command; persist the local state only in the run-scoped relationship snapshot; keep route behavior simple until a later POC implements real multi-hop forwarding.
Affects: implementations/poc14-wasm/relationship/; implementations/poc14-wasm/runtime/; implementations/poc14-wasm/cmd/poc14-analyze/; implementations/poc14-wasm/docs/; implementations/poc14-wasm/README.md.

ID: DI-vipih
Date: 2026-06-11 22:55:34
Status: active
Decision: Remove the vague `evidence_report_v1` pCID from POC14, move key rotation to `identity_key_v1`, and start new/reworked protocol payloads as pCID-owned CBOR arrays rather than universal `field_*` string maps.
Intent: The POC14 envelope should stay `grid([42(pCID), payload, proof])`, but payload shape belongs to the protocol spec named by pCID. `evidence_report_v1` blurred compute, storage, relationship, and identity semantics; key rotation is identity/key behavior, compute observations belong under `cid_compute_v1`, storage observations belong under `cas_storage_v1`, and relationship observations belong under `relationship_v1`. New protocol work should not extend the `field_*` map habit.
Constraints: Do not impose one universal payload shape across all pCIDs; keep existing legacy map payloads only as incremental POC scaffolding until their owning pCIDs are rewritten; keep kernel routing working during the transition; do not add relay forwarding or Peggy/Victor useful-work scope in this cleanup.
Affects: implementations/poc14-wasm/protocol/; implementations/poc14-wasm/pcid/; implementations/poc14-wasm/runtime/; implementations/poc14-wasm/cmd/poc14-analyze/; implementations/poc14-wasm/config.example.json; implementations/poc14-wasm/docs/; implementations/poc14-wasm/README.md; DEV-GUIDE-RESOURCES.md.

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
- [x] gogug.11 Scaffold `implementations/poc14-wasm/` from POC13 as a local POC14 module and command set.
- [x] gogug.12 Add Peggy as a WASM-boundary app process with exact-envelope evidence.
- [x] gogug.13 Add Victor as a stdio-worker app process behind a local adapter.
- [x] gogug.14 Keep WASM and stdio exchanges inside ordinary `grid([42(pCID), ...])` promise envelopes.
- [x] gogug.15 Preserve POC12 local app/kernel/device workflow in the POC14 scaffold.
- [x] gogug.16 Preserve POC13 CAS, compute, replica, verification, trust, replay, and pressure gates.
- [x] gogug.17 Add mixed-version pCID migration evidence and analyzer gates.
- [x] gogug.18 Add run-internal restart/recovery evidence and analyzer gates.
- [x] gogug.19 Replace global monitor assumptions with decentralized evidence-summary experiments.
- [x] gogug.20 Model bearer-token exchange rates as local trust/economic signals.
- [x] gogug.21 Add POC14 permanent local distrust scenario evidence and analyzer gate.
- [x] gogug.22 Add POC14 untrusted-transit exclusion scenario evidence and analyzer gate.
- [x] gogug.23 Make permanent distrust a persisted local relationship-ledger state.
- [x] gogug.24 Make untrusted-transit exclusion reject concrete local route candidates before send.
- [x] gogug.25 Remove `evidence_report_v1` from POC14 and replace key rotation with `identity_key_v1`.
- [x] gogug.26 Add pCID-owned CBOR-array payload support for the new identity/key protocol.
