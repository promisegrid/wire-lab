# TODO-hifud: POC11 adaptive trust TCP

## Decision Intent Log

ID: DI-hotos
Date: 2026-06-03 07:32:07
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement `implementations/poc11-adaptive-trust-tcp/` as a
POC10-derived live-LLM proof of concept with twelve agents across eight
containers. POC11 keeps one pCID-selected signed CBOR
`grid([42(pCID), payload, proof])` protocol over length-framed TCP, keeps every
agent as a separate process, and uses only one future-facing top-level semantic
act: `promise`. Observation, refusal, repair, economics, trust updates, and
TCP-link choices are pCID-owned payload meanings or local interpretation under
`DI-mosoj`, not separate action kinds.
Intent: POC10 proved live LLMs can make useful local promises under Go-owned
protocol validation, but it did not stress long-running relationship decay,
repair, reciprocal economics, adversarial decisions, or trust-driven TCP
adjacency. POC11 should test those pressures without drifting back into RPC-like
verbs or global trust authority.
Constraints: Treat POC11 as executable evidence, not a final LLM-agent, trust,
economics, routing, transport, monitor, or kernel API. Do not import, vendor, or
shell out to `tools/ga-runner`. Keep API keys out of repo files. Keep monitor
output observer-only. Dynamic TCP link changes mean runtime dial/accept
adjacency, not Docker network mutation. Preserve Promise Theory: each agent
promises only for itself; all trust is local; no agent promises for another
agent.
Affects: `implementations/poc11-adaptive-trust-tcp/**`;
`implementations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-hifud-poc11-adaptive-trust-tcp.md`;
`protocols/wire-lab.d/TODO/TODO.md`.

ID: DI-horuh
Date: 2026-06-03 08:19:32
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Record the first live POC11 run as provisional evidence that
promise-only live LLM autonomy works well enough to study, but is not yet
production-grade. The guide may cite the run for strong Promise Theory fit,
local trust discipline, imposition avoidance, differentiated autonomous
promises, and signed CBOR/TCP plumbing; it must also cite the observed
weaknesses: malformed or rejected decisions, one invalid JSON decision, repeated
near-duplicate promises, late message/event ordering after done markers, and
shutdown listener-close events that look like failures.
Intent: The run should answer the user-facing assessment questions without
overclaiming. It demonstrates useful autonomy and PromiseGrid alignment, but it
also identifies the next POC hardening targets before production use.
Constraints: Treat the run as evidence from one nondeterministic live execution,
not a stable benchmark or final API. Do not commit secrets, provider outputs
that contain credentials, or Docker volume state. Keep the evidence local and
Promise Theory framed: agents promise only for themselves, no agent commands
another, and trust remains local.
Affects: `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-hifud-poc11-adaptive-trust-tcp.md`;
`implementations/poc11-adaptive-trust-tcp/README.md`.

ID: DI-mudar
Date: 2026-06-03 08:22:32
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC11 Compose must pass the provider API key by a secret file path,
not by `environment` passthrough or `env_file` key/value content. The live client
must first read the configured environment variable directly, then fall back to
`<APIKeyEnv>_FILE` and read the key from that file.
Intent: `docker compose config` renders environment and env-file values in some
Compose versions. POC11 should remain runnable while preventing config
inspection from printing the API key value.
Constraints: Do not commit the key file. Keep the config JSON non-secret. Do not
introduce a broad provider abstraction for this POC. Keep tests deterministic and
use temp files for file-secret coverage.
Affects: `implementations/poc11-adaptive-trust-tcp/.gitignore`;
`implementations/poc11-adaptive-trust-tcp/compose.yaml`;
`implementations/poc11-adaptive-trust-tcp/README.md`;
`implementations/poc11-adaptive-trust-tcp/decision/live.go`;
`implementations/poc11-adaptive-trust-tcp/decision/live_test.go`;
`protocols/wire-lab.d/TODO/TODO-hifud-poc11-adaptive-trust-tcp.md`.

ID: DI-duhub
Date: 2026-06-03 08:31:22
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement the first nine POC11 hardening items in place: clean
shutdown lifecycle, provider-structured decision requests, bounded
malformed-decision repair, idempotent done/monitor ordering with a short receive
drain, durable local relationship snapshots across runs, explicit 10+ round
trust-decay/repair tests, resource-fulfillment checks for storage/compute-like
promise payloads, stake/collateral cost accounting for broken promises, and
adversarial prompt-injection/malformed-CBOR tests.
Intent: The first live run showed good autonomy and Promise Theory fit, but only
medium protocol validity. These changes harden POC11 without broadening the
top-level action vocabulary or making a production API claim.
Constraints: Keep one top-level semantic act (`promise`). Keep trust local.
Keep storage/compute/resource/stake semantics as pCID-owned payload fields and
runtime-local checks, not new top-level verbs. Do not import `tools/ga-runner`.
Do not commit Docker volume state, provider outputs, or secrets. Keep changes
inside POC11 plus its TODO/docs unless a source-map update is directly needed.
Affects: `implementations/poc11-adaptive-trust-tcp/**`;
`protocols/wire-lab.d/TODO/TODO-hifud-poc11-adaptive-trust-tcp.md`;
`DEV-GUIDE-RESOURCES.md`.

ID: DI-nanud
Date: 2026-06-04 05:23:41
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement the next ten POC11 fixes in place: bundled-target repair,
schema/prompt target constraints, coordinated shutdown grace, duplicate drain-log
removal, log-analysis tooling, updated clean-run documentation, deterministic
shutdown-race tests, live prompt fields for storage/compute/stake promises,
candidate-peer link discovery/formation, and a fresh Docker run comparison.
Intent: The hifud.8-hifud.15 run proved the hardened boundary works, but the
remaining weak points are operational: the LLM still sometimes chooses bundled
or non-direct targets, peers may close before lagging sends complete, and POC11
mostly demonstrates link removal rather than candidate-peer link formation.
Constraints: Keep one top-level semantic act (`promise`). Discovery, storage,
compute, stake, and shutdown behavior remain pCID-owned payload meanings or
runtime-local evidence, not new action kinds. Docker network membership remains
static; dynamic links mean local dial/accept promises. Keep provider outputs,
Docker volume state, and secrets out of git.
Affects: `implementations/poc11-adaptive-trust-tcp/**`;
`protocols/wire-lab.d/TODO/TODO-hifud-poc11-adaptive-trust-tcp.md`;
`DEV-GUIDE-RESOURCES.md`.

## Status

Implemented as an initial POC11 live-LLM-capable runtime and deterministic fake
test suite. The first live Docker run completed with monitor scores `4/5`
Promise Theory fit, `5/5` autonomy, `3/5` protocol validity, `5/5` local trust
correctness, and `5/5` imposition avoidance. Source: `DI-hotos`; `DI-horuh`.

## Scope

- Fork/adapt POC10's signed CBOR envelope, length-framed TCP, config, live/fake
  LLM decision boundary, runtime evidence, and observer-only monitor.
- Replace POC10's broad action vocabulary with one validated top-level
  `promise` act and pCID-owned payload meanings.
- Model directed local trust, broken-promise decay, repair by later kept
  promises, opportunity cost, reciprocal economics, adversarial/unreliable
  personas, and runtime adjacency choices.
- Use twelve agents across eight containers: Alice, Bob, Carol, and Dave as
  single-agent containers; Ellen+Frank, Grace+Heidi, Ivan+Judy, and
  Mallory+Oscar as co-located groups where each agent still runs as a separate
  process.

## Subtasks

- [x] hifud.1 Record the locked POC11 implementation decision and cross-list this
  TODO.
- [x] hifud.2 Create the standalone `implementations/poc11-adaptive-trust-tcp/`
  module with config template and ignored local config.
- [x] hifud.3 Adapt POC10 protocol and transport mechanics without depending on
  `tools/ga-runner`.
- [x] hifud.4 Implement single-action promise decision validation and fake/live
  deciders.
- [x] hifud.5 Implement local relationship, economics, and runtime adjacency
  behavior.
- [x] hifud.6 Add deterministic tests for invalid action kinds, observation and
  refusal as promise semantics, trust decay/repair, economics, and fake runtime.
- [x] hifud.7 Update implementation/resource docs and provide a user-run Docker
  command.
- [x] hifud.8 Treat normal listener shutdown as clean lifecycle evidence, not a
  broken promise.
- [x] hifud.9 Request structured live decision outputs and add bounded repair
  before final rejection.
- [x] hifud.10 Make done/monitor ordering idempotent and drain in-flight receipts
  before `node_done`.
- [x] hifud.11 Persist and reload local relationship snapshots across runs.
- [x] hifud.12 Add explicit 10+ round trust-decay/repair coverage.
- [x] hifud.13 Add resource-fulfillment and stake/collateral accounting checks.
- [x] hifud.14 Add adversarial prompt-injection and malformed-CBOR tests.
- [x] hifud.15 Update POC11 docs/source map with the hardening outcome.
- [x] hifud.16 Repair bundled target choices and tighten live target prompts.
- [x] hifud.17 Add coordinated shutdown grace and deterministic race tests.
- [x] hifud.18 Remove duplicate monitor `inflight_drained` evidence.
- [x] hifud.19 Add a POC11 log-analysis command for event and score summaries.
- [x] hifud.20 Record the latest clean run scores and remaining gaps.
- [x] hifud.21 Add provider target enum constraints from direct/candidate peers.
- [x] hifud.22 Exercise storage, compute, and stake fields in live prompt text.
- [x] hifud.23 Add candidate-peer link discovery and formation behavior.
- [x] hifud.24 Rerun POC11 and compare scores against the prior clean run.
- [x] hifud.25 Update docs/source map with hifud.16-hifud.24 outcomes.

## First Live Run Assessment

The first live POC11 Docker run completed successfully: all twelve agents wrote
`node_done`, and Dave wrote `monitor_done`. The monitor summary judged the run
as strongly Promise Theory aligned: most agents made voluntary, self-scoped,
capacity-bounded promises, kept trust local, and avoided claiming control over
others. Source: `DI-horuh`.

Transactions were real signed CBOR-grid messages over length-framed TCP, but
the economic and trust values remain POC-local. Alice and Bob exchanged
low-risk storage/trust-building promises; Bob and Carol exchanged storage and
compute willingness; Dave initiated evidence-gathering promises with Carol,
Ellen, and Heidi; Ellen made observation/introduction promises scoped to what
she directly observed; Frank made bounded relay-like promises without claiming
delivery control; Grace/Judy/Heidi/Ivan made observation, reciprocal-value, and
repair-oriented promises; Mallory and Oscar tested malformed/adversarial
pressure and later produced self-scoped promise-only messages. Source:
`DI-horuh`.

The main failures were protocol-hygiene failures, not Promise Theory failures:
several LLM decisions were rejected for missing targets, non-direct targets, or
self-targeting; Grace produced one invalid JSON response; several agents sent
near-duplicate promises; some message receipts happened after `node_done`; and
normal listener shutdown was recorded as `accept_failed`. These are the next
production-readiness targets. Source: `DI-horuh`.

## Hardening Outcome

The hifud.8-hifud.15 pass keeps POC11 on the same single top-level `promise`
action while tightening the runtime boundary around live LLM autonomy. Live
provider calls now ask for strict structured output; the decision `fields`
boundary is a strict key/value list at the provider edge and is converted back
to the existing pCID-owned runtime map before signing. Go validation remains the
semantic authority for POC behavior: it rejects non-`promise` action kinds,
non-direct or self targets, authority/RPC language, and prompt-injection
language, and it permits only one bounded repair attempt for missing common
shape fields. Source: `DI-duhub`.

The runtime now treats normal listener close as clean shutdown evidence, stops
accepting new frames before draining active receive handlers and writing
`node_done`, treats done/monitor markers as idempotent, and persists local
relationship snapshots under the run root for restart-durable trust experiments.
Storage/compute fulfillment promises are checked against local budget/capacity
before sending, inbound extreme resource promises are rejected locally, and
stake/collateral fields can spend local budget after broken-promise evidence.
Source: `DI-duhub`.

## Targeting, Discovery, And Shutdown Follow-up

The hifud.16-hifud.24 pass tightened the live target boundary and added a narrow
candidate-link formation path without introducing new top-level action kinds.
Provider-side structured output now enumerates locally visible target names; the
prompt says ordinary promises target one `direct_peers` name, candidate targets
are valid only for `promise_about=link_discovery`, and bundled target strings are
repaired to the first locally valid target when that repair is safe. Go remains
the final validator. Source: `DI-nanud`.

Candidate-peer discovery is now executable behavior: a node may dial a candidate
peer only for an explicit link-discovery promise, a receiver may accept a
candidate sender only for the same low-risk promise meaning, and a kept discovery
promise records `discovery_kept` local trust that can form a direct peer. This
is still local trust and static Docker networking; it is not a global peering
authority. Source: `DI-nanud`.

The follow-up Docker run completed with all containers exiting `0`. Compared to
the previous clean run, monitor scores stayed at autonomy `5/5`, improved
protocol validity from `3/5` to `4/5`, held imposition avoidance at `4/5`, and
settled at Promise Theory fit `4/5` and local trust correctness `4/5` after the
monitor penalized duplicate-ish commitments and two unaligned send attempts.
The new analyzer counted `143` events, no `decision_rejected` events, no
refused-connection transport failures, `2` `send_failed` events caused by
explicit `not_promised` acknowledgements, `2` `message_not_promised` events,
`2` `promise_withheld` events, `12` `turns_done` events, `12`
`shutdown_grace_elapsed` events, and `12` `inflight_drained` events. Remaining
gaps are sender-side alignment with current acceptance promises, duplicate-ish
promise churn, more substantive storage/compute fulfillment evidence, and
getting live agents to exercise candidate discovery more often. Source:
`DI-nanud`.
