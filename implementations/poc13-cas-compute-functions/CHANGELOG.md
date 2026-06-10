date:            2026-06-01
version:         poc13-cas-compute-functions
summary:         Adds POC13 as a production-progress successor to POC11.
notes:           POC13 preserves signed CBOR grid envelopes, length-framed TCP, local trust evidence, and monitor-as-observer discipline while adding multiple pCIDs routed by the local runtime, deterministic shipping device/system agents, promise-only live decisions, local reciprocal economics, and adaptive direct TCP link promises. Source: `DI-timah`; `DI-bikit`.

date:            2026-06-05
version:         poc13-cas-compute-functions-local-kernel
summary:         Splits POC13 into one local kernel process per container and separate app command entrypoints.
notes:           POC13 now routes through `poc13-kernel`; app processes register receive promises for pCIDs and keep trust, workflow, and promise judgment outside the kernel. Source: `DI-galin`.

date:            2026-06-08
version:         poc13-cas-compute-functions-superset-repair
summary:         Repairs POC13 into a strict POC11/POC12 superset while preserving CAS and compute pressure.
notes:           POC13 now keeps the POC12 app/kernel shipping architecture, POC11/POC12 relationship and analyzer lessons, and the POC13 CAS/compute/evidence pCIDs in one executable baseline. Future POCs are superset-by-default unless a scoped DI records an explicit exception. Source: `DI-sinur`.

date:            2026-06-09
version:         poc13-cas-compute-functions-run-scoped-durability
summary:         Adds run-scoped durability, retention/GC, backpressure, rate-limit, replay, recovery, and chaos evidence.
notes:           POC13 now persists CAS/evidence state under the current run root, keeps clean-run reset as the experiment boundary, records retention/GC/backpressure/rate-limit as local promises, rejects exact-envelope and capability-token replays as non-commitment evidence, expands analyzer score gates, covers recovery plus CBOR/TCP chaos cases, and routes Alice's alternate compute promise to Dave after Carol's malformed result evidence. Source: `DI-sunuf`; `DI-vahan`.
