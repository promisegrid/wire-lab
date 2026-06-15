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

ID: DI-pohoh
Date: 2026-06-12 00:20:36
Status: active
Decision: Source the POC14 OpenAI API key from the host `OPENAI_API_KEY` environment variable through the `compose.yaml` secret definition instead of requiring a local `openai_api_key.txt` file.
Intent: The clean POC14 Docker run should use Compose's secret mount path already expected by the binaries while avoiding a repo-local secret file and avoiding command-line key exposure.
Constraints: Keep containers reading `/run/secrets/openai_api_key`; do not write the API key into config files, docs, command lines, or committed files; require the host shell to export `OPENAI_API_KEY` before running Compose.
Affects: implementations/poc14-wasm/compose.yaml; implementations/poc14-wasm/README.md; implementations/poc14-wasm/scripts/run-clean.sh.

ID: DI-gahuh
Date: 2026-06-12 00:59:56
Status: active
Decision: Continue the POC14 pCID-owned payload migration by encoding scripted `cas_storage_v1` and `cid_compute_v1` exchanges as protocol-owned CBOR arrays with runtime-only compatibility projections.
Intent: The POC should stop treating the legacy `field_*` string map as the apparent target shape for storage and compute protocols while preserving current handlers, analyzer evidence, and run comparability.
Constraints: Keep `grid([42(pCID), payload, proof])`; keep live LLM non-relationship field-map turns reframed unless they supply concrete protocol payloads; do not invent new top-level action kinds; do not make one universal payload shape across pCIDs; retain compatibility fields only at the runtime boundary.
Affects: implementations/poc14-wasm/protocol/; implementations/poc14-wasm/runtime/; implementations/poc14-wasm/cmd/poc14-analyze/; implementations/poc14-wasm/docs/; implementations/poc14-wasm/README.md.

ID: DI-pamob
Date: 2026-06-12 00:59:56
Status: active
Decision: Start POC15 as a planned successor focused on real multi-hop forwarding, promise-correct route exclusion, useful WASM/stdio work, and explicit non-monolithic kernel roles.
Intent: POC15 should be a superset-oriented successor to POC14 that exercises routing through voluntary peer promises and decomposes "kernel" into transport, app-boundary, routing, and local-resource roles rather than treating the kernel as one monolith or a global monitor.
Constraints: Do not implement production-grade route authority or global monitoring; all routing and exclusion choices remain local promises and local trust judgments; do not regress POC14 shipping, CAS, compute, WASM, stdio, trust, and analyzer evidence; keep POC15 under `implementations/poc15-*` planning until code work is explicitly started.
Affects: implementations/poc15-multihop/; implementations/poc14-wasm/docs/; implementations/poc14-wasm/README.md; DEV-GUIDE-RESOURCES.md.

ID: DI-kirat
Date: 2026-06-13 17:14:50
Status: active
Decision: Active POC14/POC15 PromiseGrid-facing vocabulary must use `event`, `promise`, and `outcome` instead of production-looking `evidence`, and must reserve `boundary` for real interface or trust-line prose rather than generic event categories.
Intent: `Evidence` is useful for human design/testing claims but can leak into production semantics as if agents exchange authoritative proof. POC14 and POC15 should instead show that agents make promises, local software records events and outcomes, and every trust judgment remains local. Runtime portability pressure from Peggy and Victor should be named as runtime adapter events, not boundary evidence.
Constraints: Preserve historical DI/TODO text as append-only records; do not add compatibility aliases for old POC14 run logs because clean-run JSONL files are resettable POC artifacts; keep `grid([42(pCID), payload, proof])`; do not add new top-level action kinds; keep production-facing code, active docs, analyzer output, and fresh run events free of `evidence` names.
Affects: implementations/poc14-wasm/; implementations/poc15-multihop/; DEV-GUIDE-RESOURCES.md; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md.

ID: DI-kimim
Date: 2026-06-13 14:08:38
Status: active
Decision: Upgrade POC14 Peggy and Victor in place so Peggy executes an embedded WebAssembly module with wazero, and Victor replaces JSON-plus-hex stdio control messages with length-prefixed CBOR frames carrying exact envelope bytes.
Intent: POC14 previously had useful process-boundary coverage, but Peggy only validated WASM header bytes and Victor used JSON wrappers around hex-encoded envelope bytes. The runtime-adapter slice should now show real WASM execution and binary CBOR subprocess I/O while staying PromiseGrid-correct: ordinary signed envelopes still carry promises, adapters do not become RPC command surfaces, and all trust/workflow judgment remains local to agents.
Constraints: Preserve one top-level semantic action `promise`; keep `grid([42(pCID), payload, proof])`; keep the WASM module deterministic and embedded for this POC; do not accept arbitrary user-provided WASM modules; keep stdio as local subprocess framing, not a new wire protocol; do not edit ignored `config.json`; keep runtime state scoped to clean-run roots and Docker volumes.
Affects: implementations/poc14-wasm/; DEV-GUIDE-RESOURCES.md; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md.

ID: DI-sivis
Date: 2026-06-13 22:24:04
Status: active
Decision: Extend the POC14 runtime-adapter slice so Alice can request real `cid_compute_v1` work from Peggy and Victor: Peggy keeps the promise by executing deterministic Fibonacci logic inside wazero, and Victor keeps the promise by delegating the exact inbound compute envelope to the stdio worker over binary CBOR frames.
Intent: Peggy and Victor should no longer merely report runtime-adapter event records. They should participate as ordinary compute peers under the existing compute pCID so the POC shows that heterogeneous runtimes can keep useful promises without adding RPC verbs, global authority, a new pCID, or a special-purpose adapter protocol.
Constraints: Preserve one top-level semantic action `promise`; keep `cid_compute_v1` payloads as pCID-owned CBOR arrays; keep Alice as the requester for both compute exchanges; keep the worker subprocess local to Victor; keep the WASM module deterministic and embedded; do not accept arbitrary user-provided WASM modules; do not edit ignored `config.json`; keep analyzer/monitor output as POC-only development tooling.
Affects: implementations/poc14-wasm/; DEV-GUIDE-RESOURCES.md; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md.

ID: DI-rofiz
Date: 2026-06-13 16:56:19
Status: active
Decision: Make `implementations/poc14-wasm/config.json` the committed canonical non-secret runtime config for POC14, remove the separate `config.example.json`, and stop ignoring `config.json`.
Intent: POC14 clean regression runs must exercise the committed runtime config. A separate ignored `config.json` can drift behind the committed template and cause container runs to validate stale local behavior instead of the current code and docs.
Constraints: Keep secrets out of `config.json`; continue sourcing `OPENAI_API_KEY` through Compose secrets; keep runtime Docker volume state ignored; supersede the `DI-kimim` and `DI-sivis` constraint that prohibited editing ignored `config.json` by making `config.json` tracked and non-secret.
Affects: implementations/poc14-wasm/config.json; implementations/poc14-wasm/config.example.json; implementations/poc14-wasm/.gitignore; implementations/poc14-wasm/Dockerfile; implementations/poc14-wasm/README.md; implementations/poc14-wasm/docs/IMPLEMENTATION-NOTES.md; implementations/poc14-wasm/config/config_test.go; implementations/poc14-wasm/scripts/run-clean.sh; DEV-GUIDE-RESOURCES.md.
Supersedes: DI-kimim; DI-sivis

ID: DI-kulik
Date: 2026-06-13 17:05:29
Status: active
Decision: Increase POC14 shutdown grace in committed `config.json` so deterministic apps keep receive promises open while slower live agents finish their turns.
Intent: The clean container run after `DI-rofiz` showed broken-pipe kernel delivery failures and shutdown grace timeouts when deterministic device/runtime apps closed before live agents completed. POC14 should treat shutdown coordination as part of the regression contract so protocol-validity scoring is not dominated by local process teardown races.
Constraints: Keep the run bounded; do not hide malformed/adversarial protocol probes; do not weaken analyzer gates; preserve app-local trust semantics and clean-run volume reset behavior.
Affects: implementations/poc14-wasm/config.json; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md.

ID: DI-rulul
Date: 2026-06-13 17:49:04
Status: active
Decision: Normalize live POC14 monitor report prose before analyzer gates inspect the report.
Intent: The monitor is a POC-only observer, not a production actor or protocol participant. Live monitor output can still choose retired production-looking vocabulary even when the run events and agent payloads are correct, so POC14 should normalize only monitor prose to the active event/promise/outcome vocabulary while preserving strict analyzer rejection for protocol traffic.
Constraints: Do not weaken forbidden-vocabulary gates; do not sanitize agent messages, payloads, envelopes, or runtime events to hide protocol drift; apply this only after structured monitor JSON is decoded and before the report is returned to writers/analyzers.
Affects: implementations/poc14-wasm/decision/live.go; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md.

ID: DI-jofus
Date: 2026-06-13 17:58:44
Status: active
Decision: Remove the remaining POC14 `boundary` vocabulary from active code, docs, config, package paths, function names, run output, and analyzer/monitor acceptance.
Intent: POC14 should not keep ambiguous boundary language after the vocabulary correction. Active POC14 concepts should say what they mean: runtime adapter, process interface, app/kernel interface, monitor scope, local trust line, or experiment scope. The analyzer should fail future fresh run output if this vocabulary drifts back.
Constraints: Preserve PromiseGrid envelope behavior, one top-level `promise` action, TCP/kernel message transport, WASM execution, stdio CBOR I/O, and clean-run reset behavior; use `git mv` for path renames; avoid reintroducing the retired word inside POC14 code by splitting analyzer/normalizer literals.
Affects: implementations/poc14-wasm/; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md.

ID: DI-dirat
Date: 2026-06-13 20:50:32
Status: active
Decision: Remove POC14 shared-volume agent coordination, move run analysis to an observer-only collector service, let supervisors exit naturally, and migrate remaining POC14 pCIDs to pCID-owned CBOR array payloads.
Intent: POC14 should not let agents coordinate through a Docker volume or in-run marker files. Agent/kernel processes should communicate only by PromiseGrid envelopes over TCP or local runtime-adapter interfaces, while POC-only analysis observes after the fact through an explicit collector that cannot send control messages back to agents. The same pass should finish the `field_*` payload migration so fresh wire payloads are pCID-owned CBOR arrays rather than universal field maps.
Constraints: Preserve one top-level semantic action `promise`; keep `grid([42(pCID), payload, proof])`; keep analyzer/monitor output as POC-only tooling; do not add global trust, route authority, RPC verbs, permission/conformance language, or a universal payload shape; keep secrets out of config and command lines; keep clean-run state resettable and scoped to the POC run.
Affects: implementations/poc14-wasm/; DEV-GUIDE-RESOURCES.md; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md.

ID: DI-podut
Date: 2026-06-14 18:38:30
Status: active
Decision: Rename POC15 to `implementations/poc15-multihop-multiarity-dag/` and expand its planned target from multihop routing alone to multihop forwarding, pCID-owned multiarity, raw-message CAS/DAG review, wire-visible parent links, COSE specimens, promise-based route economics, route durability, asymmetric routes, and transport-vs-envelope proof pressure.
Intent: POC15 should test the next set of unsettled PromiseGrid design questions together because they interact: multihop route promises need incentives and failure semantics; raw-message review needs exact bytes before parsing; message DAGs need parent links to exact envelope CIDs; pCID-owned slot vectors need specimens that vary arity and slot meaning; COSE may fit as payload or proof; and transport/session signatures prove direct-hop behavior but do not automatically replace durable object-level promise proofs.
Constraints: Preserve one top-level semantic action `promise`; keep `42(pCID)` as the slot-0 bootstrap; do not introduce route authority, global trust, permission/conformance vocabulary, RPC verbs, central monitoring, or a universal payload shape; keep retention run-scoped for POC clean runs; keep POC15 non-executable until the later implementation pass starts from an explicit executable-scope DI.
Affects: implementations/poc15-multihop-multiarity-dag/; docs/thought-experiments/TE-vakah-poc15-multihop-multiarity-dag.md; DEV-GUIDE-RESOURCES.md; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md.

ID: DI-lutuv
Date: 2026-06-14 19:04:10
Status: active
Decision: Start executable POC15 by copying the POC14 executable baseline into `implementations/poc15-multihop-multiarity-dag/`, renaming module paths, command names, Compose names, config run IDs, and scripts from POC14 to POC15 while preserving the POC15 planning docs already in that directory.
Intent: The first executable POC15 slice should be a mechanically validated superset baseline before adding new raw-CAS, multihop, multiarity, parent-link, or COSE behavior. Starting from POC14 keeps prior app/kernel, shipping, CAS, compute, WASM, stdio, event collector, and analyzer gates intact so later POC15 additions can be measured as deltas rather than accidental regressions.
Constraints: Preserve one top-level semantic action `promise`; do not add route authority, global monitoring, RPC verbs, permission/conformance vocabulary, or new protocol behavior in this scaffold slice; keep POC15 run state under `/run/poc15`; keep POC15 clean runs separate from POC14; preserve existing POC15 planning docs; keep secrets out of config and command lines.
Affects: implementations/poc15-multihop-multiarity-dag/; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md; DEV-GUIDE-RESOURCES.md.

ID: DI-nivon
Date: 2026-06-14 19:26:00
Status: active
Decision: In POC15, peer-kernel forwarding retries transient peer endpoint dial failures for a short bounded window during clean-run startup.
Intent: The executable POC15 scaffold should fail on real protocol errors, not on Docker DNS/container startup ordering. Retrying only the peer-kernel dial preserves the promise semantics of the message: the sender still receives the receiver's own kept/not-promised/malformed ACK once the peer kernel is reachable, while a persistent transport failure remains a local kernel route failure.
Constraints: Do not retry app-level not-promised ACKs, malformed proofs, unsupported pCIDs, or peer trust decisions; keep the retry bounded by the existing peer send timeout; do not add route authority, central discovery, or global readiness coordination; preserve exact envelope bytes across retries.
Affects: implementations/poc15-multihop-multiarity-dag/kernel/kernel.go; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md.

ID: DI-lihir
Date: 2026-06-14 20:09:00
Status: active
Decision: Add the first executable POC15 `route_v1` slice as an Alice->Bob->Carol->Dave route setup chain using pCID-owned array payloads, app-level forwarding promises, parent exact-hash fields, and analyzer route gates.
Intent: POC15 needs real multi-hop behavior before adding more advanced route durability, asymmetric paths, raw-CAS DAGs, COSE, or variable outer arity. The first route slice should prove that neighboring apps voluntarily promise and keep one hop each, that Dave locally promises reachability, and that Alice only treats the route as usable after receiving the returned confirmation chain. This is not a kernel route authority and not a command to forward.
Constraints: Preserve one top-level semantic action `promise`; keep `route_v1` payload meanings under the pCID; do not add global route tables, route enforcement, permission vocabulary, RPC verbs, or a universal payload shape; route functions are named `runRoutePromiseWorkflow`, `handleRoutePromise`, `routePathParts`, `routePathIndex`, and `routeHopFields`; runtime state remains under `/run/poc15`; analyzer route counts are POC-only run review.
Affects: implementations/poc15-multihop-multiarity-dag/; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md; DEV-GUIDE-RESOURCES.md.

ID: DI-darur
Date: 2026-06-14 19:53:31
Status: active
Decision: POC15 peer kernels wait a short bounded window for a config-supported target app to register its local receive promise before synthesizing a kernel `not_promised` ACK.
Intent: The clean-run failure showed Peggy and Victor reaching Dave's peer kernel before Dave's app had registered `relationship_v1`, so the kernel produced a misleading transport non-commitment for a configured receive promise that was still starting. The kernel should wait only for pCIDs the target app is configured to support, preserving app-level semantic refusals while removing startup-order false negatives.
Constraints: Do not retry app-level kept/not-promised/malformed ACKs; do not wait for unknown target apps or pCIDs the target app is not configured to support; keep the wait bounded by the existing peer send timeout; preserve exact envelope bytes; do not add global readiness coordination or route authority.
Affects: implementations/poc15-multihop-multiarity-dag/kernel/kernel.go; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md.

ID: DI-daruf
Date: 2026-06-14 20:02:23
Status: active
Decision: POC15 monitor scoring instructions must rate deliberate malformed, adversarial, replay, and unsupported-pCID probes by whether agents and kernels contained them correctly, not by whether the probes appeared in the run.
Intent: The live monitor lowered protocol validity because the clean regression intentionally included malformed/adversarial cases. That makes the monitor gate punish successful negative tests. The monitor should treat protocol validity as high when valid envelopes are accepted, invalid probes are rejected or recorded as non-commitment/malformed, and trust remains local.
Constraints: Do not weaken analyzer hard gates for actual malformed accepted messages, RPC drift, forbidden vocabulary, or missing required events; do not sanitize agent traffic; keep the monitor observer-only and POC-only; preserve structured JSON monitor output.
Affects: implementations/poc15-multihop-multiarity-dag/decision/live.go; implementations/poc15-multihop-multiarity-dag/decision/live_test.go; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md.

ID: DI-tuhop
Date: 2026-06-14 20:26:41
Status: active
Decision: POC15 raw message review stores exact PromiseGrid envelope bytes through an observer-only artifact stream, with collector-owned run-scoped files under `/run/poc15/<run_id>/message-cas/<exact_sha256>.cbor` and a `message-dag.jsonl` index.
Intent: POC15 must let operators inspect actual messages, not only event records about messages. Because app containers must not coordinate through the observer Docker volume, apps emit raw envelope artifacts over stdout to their local supervisor, the supervisor forwards those records to the observer-only collector, and only the collector writes the shared run artifact files. The artifact stream records sent, received, ACK, and receive-promise envelopes by exact bytes so later DAG/parent-link work can compare message CIDs with pCID-owned payload parent fields.
Constraints: Preserve one top-level semantic action `promise`; do not let artifact storage affect app trust, routing, delivery, or peer behavior; do not mount the observer volume into app containers; keep retention scoped to one clean-run root; encode raw bytes as base64 only inside observer transport records, while persisted artifacts remain binary `.cbor`; use names `MessageArtifact`, `KindMessageArtifact`, `emitMessageArtifact`, `recordMessageArtifact`, `message-cas`, and `message-dag.jsonl`.
Affects: implementations/poc15-multihop-multiarity-dag/eventstream/; implementations/poc15-multihop-multiarity-dag/cmd/poc15-supervisor/; implementations/poc15-multihop-multiarity-dag/cmd/poc15-event-collector/; implementations/poc15-multihop-multiarity-dag/cmd/poc15-analyze/; implementations/poc15-multihop-multiarity-dag/runtime/; implementations/poc15-multihop-multiarity-dag/README.md; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md; DEV-GUIDE-RESOURCES.md.

ID: DI-bapif
Date: 2026-06-14 20:57:16
Status: active
Decision: POC15 includes a Go-based `poc15-cbor-diag` diagnostic decoder in the collector image for read-only inspection of retained raw CBOR message artifacts.
Intent: Human review should see CBOR diagnostic structure for exact message bytes, not only hashes, sizes, or event metadata; the tool must not become protocol behavior or agent-visible coordination.
Constraints: Read files only from an operator-supplied path or run-root query; do not mutate run state; keep the decoder in the POC15 module and final image; avoid adding third-party CBOR dependencies for this POC; render nested payload/proof byte strings when they contain valid CBOR; keep names `poc15-cbor-diag`, `diagnosticDecoder`, and `diagnosticValue`.
Affects: implementations/poc15-multihop-multiarity-dag/cmd/poc15-cbor-diag/main.go; implementations/poc15-multihop-multiarity-dag/cmd/poc15-cbor-diag/main_test.go; implementations/poc15-multihop-multiarity-dag/Dockerfile; implementations/poc15-multihop-multiarity-dag/README.md; DEV-GUIDE-RESOURCES.md.

ID: DI-mosat
Date: 2026-06-14 21:51:28
Status: active
Decision: Add the next executable POC15 slice as run-local message-shape specimens and analyzer gates for pCID-owned outer arity, wire-visible parent links, COSE-as-payload, COSE-as-proof, and native-proof comparison without changing normal app traffic away from `grid([42(pCID), payload, proof])`.
Intent: The first POC15 run proved a route slice and raw-message retention, but the POC still only exercised one envelope shape. The next slice should emit exact CBOR specimens into the same observer-only raw-message CAS so operators can inspect the actual bytes, and the analyzer should fail if those specimens disappear. This closes the documented multiarity/parent/COSE coverage gap without pretending the normal app/kernel API has settled on every shape yet.
Constraints: Preserve one top-level semantic action `promise`; preserve `42(pCID)` in slot 0; keep normal app traffic on the existing signed three-slot envelope until a later DI changes app transport; do not add route authority, global trust, permission/conformance vocabulary, RPC verbs, central monitoring, or a universal payload shape; keep specimen artifacts run-scoped and observer-only; use names `GridSlot`, `ByteStringGridSlot`, `RawCBORGridSlot`, `EncodeGridMessage`, `ParseGridMessage`, `GridMessage`, `COSESign1`, `EncodeCOSESign1`, `VerifyCOSESign1`, `MessageShapeSpecimen`, and `runMessageShapeSpecimenWorkflow`.
Affects: implementations/poc15-multihop-multiarity-dag/protocol/; implementations/poc15-multihop-multiarity-dag/pcid/registry.go; implementations/poc15-multihop-multiarity-dag/runtime/; implementations/poc15-multihop-multiarity-dag/cmd/poc15-analyze/; implementations/poc15-multihop-multiarity-dag/docs/; implementations/poc15-multihop-multiarity-dag/README.md; DEV-GUIDE-RESOURCES.md; protocols/wire-lab.d/TODO/TODO-gogug-poc13-hardening-and-poc14-superset-plan.md.

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
- [x] gogug.27 Source the POC14 OpenAI key from host environment through Compose secrets.
- [x] gogug.28 Convert scripted `cas_storage_v1` request and ACK payloads to pCID-owned CBOR arrays.
- [x] gogug.29 Convert scripted `cid_compute_v1` request and ACK payloads to pCID-owned CBOR arrays.
- [x] gogug.30 Add analyzer gates for migrated pCID-owned array payload evidence.
- [x] gogug.31 Update POC14 run docs with latest `poc14-demo` clean-run metrics.
- [x] gogug.32 Make Peggy and Victor perform useful routed work, not only boundary evidence.
- [x] gogug.33 Draft POC15 multi-hop forwarding plan under `implementations/poc15-multihop/`.
- [x] gogug.34 Add POC15 route-exclusion plan based on peer promises rather than omniscient filtering.
- [x] gogug.35 Document non-monolithic kernel roles using POC14/POC15 evidence.
- [x] gogug.36 Add production-fitness follow-up plan for protocol validity, local trust correctness, and Promise Theory fit blockers.
- [x] gogug.37 Rename active POC14/POC15 vocabulary from evidence/boundary categories to event/promise/outcome and runtime adapter terms.
- [x] gogug.38 Upgrade Peggy to real wazero WASM execution and Victor to binary CBOR stdio I/O.
- [x] gogug.39 Make Peggy and Victor keep real `cid_compute_v1` compute promises for Alice through WASM and stdio worker execution.
- [x] gogug.40 Make committed POC14 `config.json` canonical and remove stale ignored example/local split.
- [x] gogug.41 Increase POC14 shutdown grace to avoid deterministic app teardown racing slower live-agent sends.
- [x] gogug.42 Normalize live monitor report prose before analyzer gates inspect POC14 reports.
- [x] gogug.43 Remove remaining POC14 boundary vocabulary and add analyzer regression coverage.
- [x] gogug.44 Remove shared Docker volume and marker-file coordination from POC14 agents and kernels.
- [x] gogug.45 Add an observer-only POC14 event collector service and route analyzer input through it.
- [x] gogug.46 Change POC14 clean runs to natural container exit instead of abort-on-container-exit.
- [x] gogug.47 Migrate all remaining POC14 pCIDs to pCID-owned CBOR array payloads with runtime compatibility projections only.
- [x] gogug.48 Update POC14 docs, analyzer gates, and run scripts for the no-shared-volume array-payload regression.
- [x] gogug.49 Rename the POC15 planning directory to `implementations/poc15-multihop-multiarity-dag/`.
- [x] gogug.50 Write TE-vakah for POC15 multihop, multiarity, raw-message CAS/DAG, parent links, COSE, and proof layering.
- [x] gogug.51 Expand POC15 docs for promise-based route setup, incentives, durability, asymmetric routes, and failure semantics.
- [x] gogug.52 Add POC15 message-shape planning docs for pCID-owned arity, parent-link locations, COSE specimens, and raw artifact CAS.
- [x] gogug.53 Implement executable POC15 from the POC14 baseline under `implementations/poc15-multihop-multiarity-dag/`.
- [x] gogug.54 Add POC15 clean-run analyzer gates for `route_v1` multi-hop setup and carried-message delivery.
- [ ] gogug.55 Add POC15 clean-run analyzer gates for multiarity, parent DAGs, COSE validation, advanced route economics, durable/asymmetric routes, and useful routed WASM/stdio work.
- [x] gogug.56 Make peer-kernel delivery wait for configured app receive-promise registration before reporting startup non-commitment.
- [x] gogug.57 Harden POC15 monitor scoring instructions so intentional negative probes are scored by containment rather than presence.
- [x] gogug.58 Add POC15 raw message CAS/DAG retention and analyzer gates for operator review of exact envelope bytes.
- [x] gogug.59 Add POC15 `poc15-cbor-diag` to the collector image for operator diagnostic-format inspection of raw CBOR artifacts.
- [x] gogug.60 Add POC15 run-local multiarity, parent-link, COSE payload/proof, native-proof, and COSE tamper-rejection specimen gates.
