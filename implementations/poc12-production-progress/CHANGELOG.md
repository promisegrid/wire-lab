date:            2026-06-01
version:         poc12-production-progress
summary:         Adds POC12 as a production-progress successor to POC11.
notes:           POC12 preserves signed CBOR grid envelopes, length-framed TCP, local trust evidence, and monitor-as-observer discipline while adding multiple pCIDs routed by the local runtime, deterministic shipping device/system agents, promise-only live decisions, local reciprocal economics, and adaptive direct TCP link promises. Source: `DI-timah`; `DI-bikit`.

date:            2026-06-05
version:         poc12-production-progress-local-kernel
summary:         Splits POC12 into one local kernel process per container and separate app command entrypoints.
notes:           POC12 now routes through `poc12-kernel`; app processes register receive promises for pCIDs and keep trust, workflow, and promise judgment outside the kernel. Source: `DI-galin`.
