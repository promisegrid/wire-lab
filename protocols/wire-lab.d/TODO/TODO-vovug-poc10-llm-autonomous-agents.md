# TODO-vovug: POC10 LLM autonomous agents

## Status

Implemented. Owns
`implementations/poc10-llm-autonomous-agents/`, a seven-agent successor to POC9
that replaces scripted local strategy with live LLM-directed decisions while
preserving strict PromiseGrid message mechanics. Source: `DI-pijan`.

## Scope

- Treat POC10 as executable design evidence, not a final PromiseGrid LLM-agent
  API, SDK, discovery protocol, trust API, monitor standard, or economics
  protocol.
- Preserve the POC9 lessons: one protocol pCID for the whole protocol, signed
  CBOR `grid([42(pCID), payload, proof])` messages, length-framed TCP, sparse
  mesh, and local evidence only.
- Do not import, vendor, shell out to, or move code from `tools/ga-runner`.
- Use a repo-local config file in the POC10 directory. The config may name the
  API-key environment variable but must never contain the API key itself.
- Keep the monitor LLM as an observer/evaluator only; it must not control agents
  during the run.

## Subtasks

- [x] vovug.1 Record the locked POC10 implementation decision and cross-list this
  TODO. Done under `DI-pijan`.
- [x] vovug.2 Create the standalone `implementations/poc10-llm-autonomous-agents/`
  module with config template and ignored local config.
- [x] vovug.3 Implement fake-test and live-runtime LLM deciders without depending
  on `tools/ga-runner`.
- [x] vovug.4 Implement mixed autonomy profiles: structured action, structured
  payload, and freeform intent.
- [x] vovug.5 Implement the seven-agent runtime over signed CBOR grid envelopes
  and length-framed TCP.
- [x] vovug.6 Add monitor-as-observer evaluation after the run completes.
- [x] vovug.7 Update implementation/resource docs to describe POC10 as evidence,
  not final API.
- [x] vovug.8 Validate with fake-LLM tests, Go vet, errcheck, and a documented
  live-runtime Docker command.

## Decision Intent Log

ID: DI-pijan
Date: 2026-06-01 23:15:43
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement `implementations/poc10-llm-autonomous-agents/` as a
seven-agent successor to POC9. Runtime agents use live LLM decisions; tests use
fake LLM decisions. POC10 preserves strict PromiseGrid encoding and evidence
mechanics while allowing mixed LLM autonomy profiles. The runtime reads a
repo-local config file from the POC10 directory. The config may name the API-key
environment variable but must never contain the API key itself. POC10 must not
import, vendor, shell out to, or move code from `tools/ga-runner`.
Intent: POC8 and POC9 improved autonomy but still relied on deterministic local
plans. POC10 should test whether agents with local personas, motivations,
trust ledgers, and bounded action surfaces can make useful autonomous promises
without turning PromiseGrid into an RPC system or granting a monitor any
authority.
Constraints: Preserve one pCID for the whole POC protocol; preserve signed CBOR
`grid([42(pCID), payload, proof])`; keep all trust and valuation local to each
observing agent; keep monitor output observational only; reject secret-bearing
config fields; use fake LLMs for tests and live OpenAI-compatible Responses API
calls only at runtime; do not depend on `tools/ga-runner`.
Affects: `implementations/poc10-llm-autonomous-agents/**`;
`implementations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-vovug-poc10-llm-autonomous-agents.md`;
`protocols/wire-lab.d/TODO/TODO.md`.
