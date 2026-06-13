date:            2026-06-11
version:         poc14-wasm
summary:         Adds POC14 as a POC13 superset with WASM, stdio, decentralized-monitoring, migration, and restart event records.
notes:           POC14 is scaffolded from POC13 and preserves inherited app/kernel, shipping, CAS, compute, trust, replay, pressure, and analyzer gates. It adds Peggy as a WASM-adapter app process, Victor as a stdio-worker agent behind a local adapter, exact-envelope forwarding through the same local kernel path, decentralized monitoring event records that avoid global observer assumptions, mixed-version pCID migration event records, and same-run restart/recovery event records. Source: `DI-sihuz`; `DI-sifot`; `DI-fimoh`; `DI-lulof`; `DI-linof`.
