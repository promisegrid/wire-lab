# POC16 Kernel Role Split Plan

POC16 should make kernel roles explicit without requiring every runtime to split
them into separate operating-system processes. The point is to make ownership and
promise boundaries visible. Source: `DI-pamob`.

## Roles To Split Or Name

1. **Transport role:** owns direct TCP framing, peer dialing/listening, send
   failure event, and exact-byte forwarding to a direct peer.
2. **App-interface role:** owns local app registration, pCID receive promises, and
   local delivery queues.
3. **Routing role:** owns route-promise selection from local events and
   per-hop forwarding promises; it does not own a global route table.
4. **Local-resource role:** owns hardware/storage/compute resource promises such
   as printer-port capability tokens, CAS retention promises, or compute capacity
   promises.
5. **Lifecycle/resource-protection role:** owns conditional capability promises
   for process lifetime, CPU, RAM, sockets, devices, storage, quotas, quiesce,
   drain, and shutdown. It may withdraw, throttle, or terminate access to keep
   its own resource-protection promises, but it does not command apps or promise
   on their behalf. Source: `TE-ragin`; `DI-vuruz`.
6. **Event role:** owns local events journals and voluntary summary
   promises; it does not become a global monitor.

## First Executable Shape

The first POC16 implementation should prefer separate objects first, then split
processes only where the interface matters:

- Keep one local transport process per container if that preserves POC14
  stability.
- Add a routing-role object with tests for path selection, route exclusion, and
  forwarding non-commitment.
- Keep `printer_port`, CAS, and compute as local-resource app roles with clear
  capability-token and capacity promises.
- Treat the local supervisor as a lifecycle/resource-protection role when it
  controls local process lifetime or resource quotas; shutdown by signal or kill
  is a host mechanism for withdrawing a conditional capability promise, not
  authorization enforcement.
- Keep app receive registration as an app-interface promise rather than making the
  transport role decide app semantics.
- Add analyzer event that records which roles were split and which were
  intentionally collapsed for the Docker runtime.

## Production Interpretation

A production node may implement these roles as one daemon, several daemons,
browser APIs, WASM host functions, firmware functions, or local objects. The
portable requirement is not the process layout; it is the promise interface:
agents must be able to tell which local role promised transport, app delivery,
resource access, lifecycle/resource protection, route selection, and event
retention.
