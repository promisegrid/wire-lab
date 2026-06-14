# POC14 Superset Plan

POC14 should begin as a superset of POC13 unless a later DI explicitly narrows
scope. It should also test heterogeneous process boundaries beyond same-language
TCP-only Go agents. Source: `DI-sihuz`; `DI-sifot`; `DI-fimoh`;
`DI-lulof`; `DI-linof`.

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
- Add richer POC-local operators for analyzer/monitor comparison, failed-run
  triage, and explicit event export for review while keeping those tools out
  of the production trust model.

## Implemented First Slice

- Scaffolded `implementations/poc14-wasm/` from POC13 as a superset baseline
  with POC14-local module, command, Docker, run-root, and analyzer names.
- Added Peggy as a deterministic WASM-adapter app process. Peggy now executes an
  embedded no-import Fibonacci module with wazero, exchanges ordinary
  `relationship_v1` PromiseGrid envelopes through the local kernel, and keeps
  Alice's `cid_compute_v1` compute promise through the WASM export. Source:
  `DI-kimim`; `DI-sivis`.
- Added Victor as a deterministic stdio-adapter app. Victor's worker process
  emits and observes exact PromiseGrid envelope bytes as CBOR byte strings inside
  length-prefixed CBOR frames over stdin/stdout; the adapter forwards those bytes
  through the local kernel and delegates Alice's exact inbound `cid_compute_v1`
  envelope to the worker for a signed compute ACK. Source: `DI-kimim`;
  `DI-sivis`.
- Added analyzer gates for WASM, stdio, and decentralized monitoring event records
  so POC14 cannot pass by merely preserving the POC13 baseline.
- Added decentralized monitoring event records for local summaries, peer-carried
  attestations, bearer-token exchange-rate signals, topology signals, and
  voluntary gossip. Source: `DI-linof`; `DI-lulof`.
- Added hard local trust-line scenario event for permanent local distrust
  of Mallory and local rejection of Mallory-transit route candidates for Alice's
  own traffic. Source: `DI-kinaf`.
- Added app-local behavior for those hard boundaries: the relationship ledger
  persists permanent distrust, send gates reject ordinary future Mallory sends,
  and route-candidate checks reject Mallory as a transit hop. Source: `DI-dubih`.
- Added explicit mixed-version pCID migration event records and same-run restart
  recovery event records so those planned POC14 concerns have analyzer gates rather
  than remaining prose-only intentions. Source: `DI-linof`.

## Production Monitoring Scope

- Assume production agents are geographically distributed and owned by different
  legal entities. No process gets a global view of messages, local events,
  trust updates, exchange offers, or failures unless other agents voluntarily
  promise to disclose them.
- Treat POC analyzers and monitors as development-time local observers, not as
  production authorities. In production, an analyzer can only be an ordinary
  agent promising what it observed, what event it retained, and what local
  interpretation it made from that event.
- Do not design POC14 around a global health dashboard, global trust score,
  global exchange rate, or central audit trail. Those would contradict the
  PromiseGrid assumption that trust and event interpretation are local to
  each agent relationship.

## Decentralized Monitoring Candidates

- Local events summaries: each agent may promise signed summaries of its own
  keep/break observations, but peers decide locally whether those summaries are
  credible.
- Peer-carried attestations: Alice can carry Bob's signed promise history to
  Carol as event, while Carol remains free to discount it based on Carol's
  trust in Alice and Bob.
- Bearer-token exchange rates: the price peers offer for tokens issued by Alice,
  Bob, Carol, or Dave can act as a decentralized market signal for perceived
  trustworthiness, usefulness, scarcity, and recent promise behavior.
- Relationship topology: direct TCP links, relay willingness, storage
  replication choices, and compute-verification choices can be measured as local
  signals of relationship strength rather than as centrally managed routing
  policy.
- Voluntary gossip: agents may promise to share selected local observations with
  trusted peers, but recipients keep those reports as event, not as facts
  imposed by an authority.

## Non-Goals Until Decided

- Do not add global trust, centralized authorization, RPC command verbs, or
  cross-run persistent POC state by default.
- Do not drop POC13 behavior silently to make POC14 simpler.
- Do not treat monitor output as authority over agents; it remains observer
  event.
- Do not let WASM host calls or stdio adapters become RPC command channels; they
  should carry promises and exact envelope event, not external authority.
