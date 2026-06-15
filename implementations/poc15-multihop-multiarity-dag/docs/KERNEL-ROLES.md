# POC15 Kernel Role Split Plan

POC15 should make kernel roles explicit without requiring every runtime to split
them into separate operating-system processes. The point is to make ownership and
promise boundaries visible. Source: `DI-pamob`.

## Roles To Split Or Name

1. **Transport role:** owns direct TCP framing, peer dialing/listening, send
   failure event, and exact-byte forwarding to a direct peer.
2. **App-boundary role:** owns local app registration, pCID receive promises, and
   local delivery queues.
3. **Routing role:** owns route-promise selection from local events and
   per-hop forwarding promises; it does not own a global route table.
4. **Local-resource role:** owns hardware/storage/compute resource promises such
   as printer-port capability tokens, CAS retention promises, or compute capacity
   promises.
5. **Event role:** owns local events journals and voluntary summary
   promises; it does not become a global monitor.

## First Executable Shape

The first POC15 implementation should prefer separate objects first, then split
processes only where the boundary matters:

- Keep one local transport process per container if that preserves POC14
  stability.
- Add a routing-role object with tests for path selection, route exclusion, and
  forwarding non-commitment.
- Keep `printer_port`, CAS, and compute as local-resource app roles with clear
  capability-token and capacity promises.
- Keep app receive registration as an app-boundary promise rather than making the
  transport role decide app semantics.
- Add analyzer event that records which roles were split and which were
  intentionally collapsed for the Docker runtime.

## Production Interpretation

A production node may implement these roles as one daemon, several daemons,
browser APIs, WASM host functions, firmware functions, or local objects. The
portable requirement is not the process layout; it is the promise boundary:
agents must be able to tell which local role promised transport, app delivery,
resource access, route selection, and event retention.
