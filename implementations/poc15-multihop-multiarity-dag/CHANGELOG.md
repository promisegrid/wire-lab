date:            2026-06-14
version:         poc15-multihop-multiarity-dag
summary:         Starts executable POC15 as a mechanically renamed POC14 baseline under the multihop/multiarity/DAG planning directory.
notes:           POC15 is scaffolded from POC14 and preserves inherited app/kernel, shipping, CAS, compute, trust, replay, pressure, WASM, stdio, observer-collector, and analyzer gates. It now adds the first route_v1 multi-hop slice: Alice confirms and uses an Alice->Bob->Carol->Dave route through voluntary neighboring route promises. Peer-kernel startup dials retry briefly so clean runs fail on protocol behavior rather than transient Docker DNS readiness. Source: `DI-lutuv`; `DI-nivon`; `DI-lihir`.
