# POC14 Production Fitness

POC14 is executable POC evidence, not production software. Its purpose is to test
whether the POC13 superset can survive heterogeneous app processes,
decentralized monitoring, hard local distrust, and untrusted-transit exclusion.
Source: `DI-linof`; `DI-lulof`; `DI-kinaf`; `DI-dubih`.

## Current Fitness Claim

- POC14 should remain fit for continued POC work if `poc14-analyze` reports all
  inherited POC13 gates plus `boundary=5`, `monitoring=5`, `migration=5`, and
  `restart=5`.
- The 2026-06-11 `poc14-demo` clean Docker run passed those analyzer gates with
  2079 total events, all score dimensions at `5`, explicit behavioral evidence
  for permanent local distrust plus untrusted-transit exclusion, and empty
  `rpc_drift_counts` / `resource_trust_coupling_counts`.
- POC14 is not production-fit until a later implementation replaces POC-local
  whole-run analysis with ordinary agents that exchange local evidence promises.
- POC14's analyzer and monitor remain development-time observers. In real
  deployments, no agent can assume a global run log, global trust score, global
  exchange rate, or central monitor.

## Expected Blockers After First Run

- The WASM role validates a module fixture but does not yet execute arbitrary
  WASM application logic.
- The stdio role demonstrates exact envelope exchange through a worker process,
  but only for one deterministic startup promise.
- Decentralized monitoring signals are recorded as evidence, but agents do not
  yet adapt peer choice from token exchange rates or peer-carried attestations.
- Permanent distrust and transit exclusion now have app-local behavior: the
  ledger blocks future Mallory sends and rejects route candidates with Mallory as
  a transit hop. POC14 still does not implement real multi-hop forwarding, so it
  proves route-candidate rejection rather than end-to-end route execution.
- Mixed-version migration and restart recovery are currently evidenced as local
  promises and analyzer gates; later POCs should add stronger fault injection
  that kills and restarts multiple app processes during one run.
- Production deployment still needs packaging, credential handling, real device
  and storage backends, adversarial transport testing, and per-agent operations
  documentation.
