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
