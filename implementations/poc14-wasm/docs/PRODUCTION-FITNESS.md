# POC14 Production Fitness

POC14 is executable POC event, not production software. Its purpose is to test
whether the POC13 superset can survive heterogeneous app processes,
decentralized monitoring, hard local distrust, and untrusted-transit exclusion.
Source: `DI-linof`; `DI-lulof`; `DI-kinaf`; `DI-dubih`.

## Current Fitness Claim

- POC14 should remain fit for continued POC work if `poc14-analyze` reports all
  inherited POC13 gates plus `runtime_adapter=5`, `monitoring=5`, `migration=5`, and
  `restart=5`.
- The 2026-06-12 `poc14-demo` clean Docker run passed those analyzer gates with
  2092 total events, all score dimensions at `5`, explicit behavioral event
  for permanent local distrust plus untrusted-transit exclusion, identity-key
  cleanup event, and empty `rpc_drift_counts` /
  `resource_trust_coupling_counts`.
- The next clean run must satisfy the new `DI-gahuh` migrated-payload gates for
  `cas_storage_v1` and `cid_compute_v1` array requests and ACKs.
- POC14 is not production-fit until a later implementation replaces POC-local
  whole-run analysis with ordinary agents that exchange local events promises.
- POC14's analyzer and monitor remain development-time observers. In real
  deployments, no agent can assume a global run log, global trust score, global
  exchange rate, or central monitor.

## Expected Blockers After First Run

- The WASM role validates a module fixture and now sends a useful relationship
  promise about reusable module-validation event, but it still does not
  execute arbitrary WASM application logic.
- The stdio role demonstrates exact envelope exchange through a worker process,
  and now sends a useful relationship promise about subprocess round-trip
  event, but it still covers one deterministic subprocess fixture.
- Decentralized monitoring signals are recorded as event, but agents do not
  yet adapt peer choice from token exchange rates or peer-carried attestations.
- Permanent distrust and transit exclusion now have app-local behavior: the
  ledger blocks future Mallory sends and rejects route candidates with Mallory as
  a transit hop. POC14 still does not implement real multi-hop forwarding, so it
  proves route-candidate rejection rather than end-to-end route execution.
- Mixed-version migration and restart recovery are currently shown as local
  promises and analyzer gates; later POCs should add stronger fault injection
  that kills and restarts multiple app processes during one run.
- Production deployment still needs packaging, credential handling, real device
  and storage backends, adversarial transport testing, and per-agent operations
  documentation.
- Production deployment also needs real multi-hop forwarding with route
  exclusion based on peer promises rather than a POC-local route filter. POC15 is
  the planned successor for that work. Source: `DI-pamob`.

## Follow-Up Plan

- **Protocol validity:** Keep migrating concrete pCIDs away from legacy field
  maps into pCID-owned CBOR arrays; reject any new universal payload shape unless
  a pCID-specific TE/DI locks it. Source: `DI-gahuh`.
- **Local trust correctness:** Add POC15 forwarding event that distinguishes
  direct peer non-commitment, forwarder non-commitment, final receiver
  non-commitment, transport unavailability, and a broken forwarding promise by
  the hop that actually promised forwarding.
- **Promise Theory fit:** Route exclusion must be based on voluntary peer
  promises and local keep/break history, not on commands, global bans,
  permission, authorization, conformance, or route enforcement. Source:
  `DI-pamob`.
- **Kernel model:** Continue documenting kernel roles as local promise surfaces
  that can be split or collapsed by runtime: transport, app boundary, pCID
  routing, local resource, adapter, and event roles.
