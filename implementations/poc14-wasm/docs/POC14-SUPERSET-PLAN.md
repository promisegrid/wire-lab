# POC14 Superset Plan

POC14 should begin as a superset of POC13 unless a later DI explicitly narrows
scope. It should also test heterogeneous process boundaries beyond same-language
TCP-only Go agents. Source: `DI-sihuz`; `DI-sifot`; `DI-fimoh`.

## Baseline Promises To Preserve

- Keep POC11-style autonomous sparse-mesh relationship/economics pressure.
- Keep POC12-style separate app processes, local container kernels, pCID routing,
  device/system apps, and shipping workflow.
- Keep POC13 CAS storage, replica retrieval, CID-named compute, cache reuse,
  verification/disagreement, adversarial malformed pressure, run-scoped
  durability, retention/GC, backpressure, rate-limit, replay protection, bounded
  local trust, and dynamic topology gates.
- Keep one top-level semantic action: `promise`.

## Candidate POC14 Additions

- Replace more deterministic startup sequencing with LLM-chosen but schema-safe
  workflow promises.
- Add stronger peer bootstrapping and authenticated introduction promises.
- Add mixed-version pCID behavior and migration promises.
- Add run-internal crash/restart orchestration for multiple app processes, not
  only durable-state unit tests.
- Add one or more WASM agents that run in their own process. The point is to
  test sandboxed agent execution, host-call boundaries, and the ability to keep
  PromiseGrid message semantics outside the WASM runtime while still exchanging
  exact pCID-defined envelopes.
- Add one or more agents that run in their own process and do all messaging via
  stdio. The point is to test subprocess-style adapters where the agent has no
  direct network API and the local kernel or adapter is the only path for
  sending and receiving PromiseGrid envelopes.
- Add resource-isolation kernel roles for CPU, memory, filesystem, and device
  access without moving app trust into the kernel.
- Add richer production operators: analyzer/monitor comparison, failed-run
  triage, and explicit evidence export for review.

## Non-Goals Until Decided

- Do not add global trust, centralized authorization, RPC command verbs, or
  cross-run persistent POC state by default.
- Do not drop POC13 behavior silently to make POC14 simpler.
- Do not treat monitor output as authority over agents; it remains observer
  evidence.
- Do not let WASM host calls or stdio adapters become RPC command channels; they
  should carry promises and exact envelope evidence, not external authority.
