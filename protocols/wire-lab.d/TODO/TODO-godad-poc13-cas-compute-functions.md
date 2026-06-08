# TODO-godad: POC13 CAS storage and CID-named compute functions

## Decision Intent Log

ID: DI-bibom
Date: 2026-06-06 01:12:50
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Define POC13 as `implementations/poc13-cas-compute-functions/`, a
future executable proof of concept for decentralized CAS storage promises and
CID-named function-call compute promises. POC13 should use two provisional
protocol pCID names, `cas_storage_v1` and `cid_compute_v1`, under the same
`grid([42(pCID), ...protocol-defined-slots])` envelope discipline. A pCID names
the protocol spec; message types and function CIDs live inside pCID-owned
payloads.
Intent: POC12 proves a production-ish app/kernel/device workflow, but it does
not yet prove PromiseGrid's core storage/compute shape. POC13 should focus on a
decentralized sparse CAS where content identity is separate from availability,
retention, replication, access, and serving promises, plus compute where code is
stored in CAS and called by CID. Pure results can be cached by exact function,
input, protocol, and context identity; impure work must externalize timestamp,
randomness, sensor reads, or other ambient inputs as explicit context objects so
the run is replayable and pure-after-the-fact.
Constraints: This TODO defines POC13 but does not implement the executable POC
yet. Do not add separate pCIDs for message types, global registries, central
storage authority, RPC verbs, permission/conformance framing, hidden ambient
inputs, or provider-backed runs. Keep storage, compute, capability, retention,
and cache behavior expressed as promises and local evidence.
Affects: `protocols/wire-lab.d/TODO/TODO-godad-poc13-cas-compute-functions.md`;
`protocols/wire-lab.d/TODO/TODO.md`;
`docs/research/DN-nuras-poc13-cas-compute-functions.md`;
`implementations/README.md`;
`DEV-GUIDE-RESOURCES.md`; storage/compute scenarios; future
`implementations/poc13-cas-compute-functions/**`.

ID: DI-notig
Date: 2026-06-06 02:04:59
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement POC13 as a self-contained Go proof of concept under
`implementations/poc13-cas-compute-functions/`, with commands
`poc13-supervisor`, `poc13-agent`, and `poc13-analyze`. The first executable
version uses LLM-backed local agents Alice, Bob, Carol, Dave, Ellen, Frank,
Grace, and Mallory across Docker containers, with deterministic protocol
validation in Go before any LLM decision is trusted.
Intent: POC13 should pressure-test CAS storage and CID-named compute as
PromiseGrid promises rather than RPC calls. The POC needs enough live autonomy
to expose local decision and malformed-input pressure, while keeping wire shape,
pCID handling, signatures, CID verification, cache keys, and analyzer gates
deterministic enough to diagnose.
Constraints: Keep POC13 self-contained and do not import POC12 code or extract
a shared library yet. Use only provisional protocol pCID names
`cas_storage_v1` and `cid_compute_v1`; message variants, content CIDs, function
CIDs, and context CIDs live inside pCID-owned payloads. Keep one top-level
semantic act, `promise`, and do not add global registries, central storage or
compute authority, hidden RPC verbs, permission/conformance framing, or stored
API keys.
Affects: `implementations/poc13-cas-compute-functions/**`;
`implementations/README.md`;
`DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-godad-poc13-cas-compute-functions.md`.

ID: DI-lasuh
Date: 2026-06-06 06:44:02
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Improve POC13's live-provider path without changing the POC13
architecture, pCID set, envelope shape, or role-script topology. The decider
must extract meaningful text from current Responses API shapes instead of
recording the placeholder `provider returned no output_text` when nested output
content is available. The analyzer must flag placeholder live decisions as an
evidence quality problem. The allowed implementation names for this change are
`ProviderResponse`, `ProviderOutput`, `ProviderContent`, `ResponseText`,
`PlaceholderLiveDecisionCounts`, `HasPlaceholderLiveDecision`, and local
variables derived from those names.
Intent: The latest live run proved that provider calls can happen, but the
recorded decision evidence was weak because the parser only checked a top-level
`output_text` field. POC13 should produce auditable live-local promise
judgments when a provider key is present, and should make placeholder decisions
visible as analyzer failures rather than silently treating them as equivalent to
autonomous behavior.
Constraints: Keep POC13 self-contained. Do not import GA-runner, POC12, or a
shared library. Do not store API keys. Do not add new protocol pCIDs, top-level
action kinds, RPC verbs, permission/conformance framing, or global authority.
Keep existing role scripts as bounded deterministic scenario evidence; this DI
improves live-decision evidence extraction and analyzer quality gates rather
than adding real transport or full autonomous negotiation.
Affects: `implementations/poc13-cas-compute-functions/decision.go`;
`implementations/poc13-cas-compute-functions/analyze.go`;
`implementations/poc13-cas-compute-functions/analyze_test.go`;
`implementations/poc13-cas-compute-functions/decision_test.go`;
`implementations/poc13-cas-compute-functions/README.md`;
`protocols/wire-lab.d/TODO/TODO-godad-poc13-cas-compute-functions.md`.

ID: DI-fumol
Date: 2026-06-06 07:03:14
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Expand POC13 in place from local evidence generation into a bounded
TCP-delivered protocol proof while preserving the existing directory
`implementations/poc13-cas-compute-functions/`, the existing commands, the
existing Docker Compose service set, the two existing provisional pCID names
`cas_storage_v1` and `cid_compute_v1`, and the envelope discipline
`grid([42(pCID), payload, proof])`. The POC13 supervisor becomes the primary
runtime: each container listens on TCP, sends length-framed raw signed grid
envelopes to peer containers, parses and verifies received envelopes before
delivery, and records local evidence per observing agent. LLM/local decision
text gates whether an agent acts, and signed payloads carry that decision text.
Concrete behavior must include real CAS store/serve/retrieve/replicate checks,
CID-named compute over payload-provided function/input/context bytes, local
trust updates, voluntary repair promises, credit-style economics promises, and
issuer-local capability-token issue/redeem evidence.
Intent: POC13 was useful protocol-shape evidence but too shallow: role scripts
recorded evidence without real transport, storage, retrieval, dynamic compute,
or trust/economy/capability pressure. The next version should still be a bounded
POC, not a final API, but it should make the storage/compute promises concrete
enough that analyzer gates prove peer messages crossed TCP, bytes were stored
and retrieved by CID, compute results came from payload-provided function/input
material, and trust/economics/capability records stayed local and voluntary.
Constraints: Keep POC13 self-contained; do not import POC12, GA-runner, or a
new shared library. Do not add new protocol pCIDs, top-level action kinds,
global authorities, permission/conformance framing, stored API keys, or final
storage/compute APIs. Keep TCP framing simple and local to POC13. The approved
new file path is
`implementations/poc13-cas-compute-functions/runtime.go`. Approved changed paths
are `implementations/poc13-cas-compute-functions/{analyze.go,analyze_test.go,cmd/poc13-supervisor/main.go,config.example.json,config.go,README.md,runtime.go,runtime_test.go,scenario_test.go}`;
`DEV-GUIDE-RESOURCES.md`; `implementations/README.md`;
`protocols/wire-lab.d/TODO/TODO.md`; and this TODO file. Approved runtime path
patterns are Docker-managed `poc13-run` volume paths under `/run/poc13/**`,
container-local TCP listeners on the configured POC13 port, and `/tmp/wire-lab-gocache`
for validation. Approved implementation names include `TCPRuntime`,
`AgentState`, `OutboundPromise`, `RuntimeMessage`, `FrameReader`,
`FrameWriter`, `ExecuteFunction`, `NewTCPRuntime`, `Run`, `runLocalAgents`,
`runLocalAgent`, `handleEnvelope`, `sendPromise`, `recordAgentEvent`,
`adjustTrust`, `issueCapabilityToken`, `redeemCapabilityToken`,
`containerForAgent`, `dialAddressForAgent`, `listenAddress`, `Promises`,
`StartupDelay`, `SettleDelay`, `ListenPort`, `StartupDelayMillis`, and
`SettleDelayMillis`, plus local variables derived from those names.
Affects: `implementations/poc13-cas-compute-functions/**`;
`implementations/README.md`;
`DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-godad-poc13-cas-compute-functions.md`;
`protocols/wire-lab.d/TODO/TODO.md`.

ID: DI-hohuf
Date: 2026-06-06 07:12:17
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add the POC13 clean Docker regression runner at
`implementations/poc13-cas-compute-functions/scripts/run-clean.sh` instead of
using an untracked `/tmp` script. The script resets the Docker Compose state,
rebuilds and runs POC13, and then runs the analyzer against
`/run/poc13/poc13-demo`.
Intent: POC13's expanded TCP runtime needs a repeatable repo-local command that
Steve can run directly, while preserving explicit exit-code handling and
avoiding hidden local helper scripts.
Constraints: Do not store secrets in the script. Do not use `|| true`. Keep the
script under the POC13 directory and keep Docker-managed runtime state in the
existing `poc13-run` volume under `/run/poc13/**`.
Affects: `implementations/poc13-cas-compute-functions/scripts/run-clean.sh`;
`protocols/wire-lab.d/TODO/TODO-godad-poc13-cas-compute-functions.md`.

ID: DI-mosil
Date: 2026-06-06 07:20:15
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Harden POC13 startup and shutdown by replacing fixed
`StartupDelayMillis` / `SettleDelayMillis` sleeps with explicit readiness and
done evidence plus bounded readiness and idle-completion gates. Each container
records a local readiness promise marker after its TCP listener is open, waits
for peer readiness markers before running local agents, and records a local done
promise only after local agents have finished and the runtime has observed an
idle period with no active TCP handlers. The analyzer must require readiness,
peer-readiness observation, and runtime-done evidence.
Intent: The expanded TCP POC should not depend on arbitrary fixed startup and
settle sleeps. Startup should wait for explicit peer readiness evidence, and
shutdown should be tied to runtime quiescence rather than a blind timeout, while
remaining bounded and diagnosable in Docker.
Constraints: Keep this as POC13-local runtime coordination, not a PromiseGrid
global authority or final API. Do not add new pCIDs or top-level action kinds.
Runtime marker files are local evidence in the existing Docker-managed run
volume, not protocol authority. Approved changed paths are
`implementations/poc13-cas-compute-functions/{analyze.go,analyze_test.go,config.example.json,config.go,runtime.go,scenario_test.go,README.md}`;
`DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-godad-poc13-cas-compute-functions.md`;
and `protocols/wire-lab.d/TODO/TODO.md`. Approved runtime path pattern:
`/run/poc13/<run_id>/runtime/*.ready` and
`/run/poc13/<run_id>/runtime/*.done`, created/overwritten by Docker containers
and removed by the existing clean-run volume reset. Approved implementation
names include `ReadinessTimeoutMillis`, `CompletionIdleMillis`,
`ReadinessTimeout`, `CompletionIdle`, `recordRuntimeReadiness`,
`waitForPeerReadiness`, `waitForRuntimeCompletion`, `recordRuntimeDone`,
`runtimeMarkerDir`, `runtimeMarkerPath`, `markActivity`, `activeHandlers`,
`lastActivity`, and local variables derived from those names.
Affects: `implementations/poc13-cas-compute-functions/**`;
`DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-godad-poc13-cas-compute-functions.md`.

ID: DI-lupag
Date: 2026-06-07 07:27:39
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Extend POC13 in place with the six selected next-step improvements:
retrieve from Frank when Bob is unavailable and Frank is locally trusted enough;
verify Carol's compute result through Alice-local recomputation and Dave/Grace
verification evidence; add richer bounded economics with capacity limits,
opportunity cost, price refusal, credits spent, and credits earned; add
trust-driven peer choice evidence where Alice chooses storage and compute peers
from local trust and prior evidence; add per-agent narrative documentation with
the exact message sequence and incentives; and add analyzer score/report fields
in addition to raw event counts.
Intent: The TCP POC should now demonstrate not only that storage/compute
messages can move, but that agents make local choices based on trust and
economics, can recover from a peer outage by asking a replica peer, can verify
compute evidence, and can produce a readable run report for guide work.
Constraints: Keep POC13 self-contained and bounded. Do not add new protocol
pCIDs, top-level action kinds, global trust, permissions, or central pricing.
Treat Bob's outage as a scenario event and local non-commitment, not a global
failure. Keep economics as local promise/evidence fields inside existing
pCID-owned payloads. Approved changed paths are
`implementations/poc13-cas-compute-functions/{analyze.go,analyze_test.go,runtime.go,runtime_test.go,README.md}`;
`implementations/poc13-cas-compute-functions/docs/RUN-NARRATIVE.md`;
`DEV-GUIDE-RESOURCES.md`; `implementations/README.md`;
`protocols/wire-lab.d/TODO/TODO-godad-poc13-cas-compute-functions.md`; and
`protocols/wire-lab.d/TODO/TODO.md`. Approved implementation names include
`ScoreReport`, `Narrative`, `scores`, `computeScores`, `addScore`,
`TrustDrivenChoiceCounts`, `EconomicsCounts`, `VerificationCounts`,
`ReplicaRecoveryCounts`, `recordTrustDrivenChoice`, `recordEconomics`,
`handleReplicaServeRequest`, `handleComputeVerificationRequest`,
`handleComputeVerificationResult`, `verifyComputeResultLocally`,
`computePrice`, `storageCapacity`, `computeCapacity`, `preferredStoragePeer`,
`preferredComputePeer`, and local variables derived from those names.
Affects: `implementations/poc13-cas-compute-functions/**`;
`DEV-GUIDE-RESOURCES.md`;
`implementations/README.md`;
`protocols/wire-lab.d/TODO/TODO-godad-poc13-cas-compute-functions.md`;
`protocols/wire-lab.d/TODO/TODO.md`.

## Prior aliases

- None.

## Status

Implemented as first executable evidence. POC13 remains provisional and does
not freeze final storage, compute, cache, provider, kernel, or app APIs.

## Tasks

- [x] godad.1 Record the POC13 definition DI and cross-list this TODO.
- [x] godad.2 Write `docs/research/DN-nuras-poc13-cas-compute-functions.md`.
- [x] godad.3 Document decentralized CAS storage promises: content identity,
  availability, retention, replication, access, serving, corrupt-byte evidence,
  and partial availability.
- [x] godad.4 Document CID-named compute promises: `function_cid`, inputs,
  declared context objects, result cache keys, pure-after-the-fact impure calls,
  and broken/malformed result evidence.
- [x] godad.5 Expand or add scenarios for CAS storage and CID-named function
  compute pressure before implementing the executable POC.
- [x] godad.6 Update `DEV-GUIDE-RESOURCES.md` and `implementations/README.md`
  so POC13 is visible as storage/compute evidence, not a stable API.
- [x] godad.7 Implement `implementations/poc13-cas-compute-functions/` as a
  self-contained first executable POC with analyzer gates.
- [x] godad.8 Improve live decision extraction and analyzer gates so provider
  runs record meaningful local promise judgments instead of silent placeholders.
- [x] godad.9 Expand POC13 into a TCP-delivered bounded runtime with concrete
  CAS storage/retrieval/replication, dynamic compute, local trust, voluntary
  repair, credit economics, and capability-token evidence.
- [x] godad.10 Replace fixed POC13 startup/settle sleeps with explicit
  readiness and done evidence plus bounded readiness and idle-completion gates.
- [x] godad.11 Add selected POC13 improvements: replica recovery, compute
  verification, richer economics, trust-driven peer choice, narrative docs, and
  analyzer score/report fields.

## Acceptance criteria

- POC13's storage and compute protocols keep one top-level semantic act:
  `promise`.
- `cas_storage_v1` and `cid_compute_v1` are protocol pCID names, not message
  type names.
- Function code identity is payload-level CAS identity, not envelope pCID
  identity.
- Content identity never implies availability, retention, access, or serving
  authority.
- Impure compute inputs are explicit context promises or context objects so
  result evidence can be replayed and audited locally.
