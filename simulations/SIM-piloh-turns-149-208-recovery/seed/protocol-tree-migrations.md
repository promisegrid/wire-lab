# Protocol Tree Migrations

Candidate protocol trees moved into this simulation as specimens. The move
preserves their repo-like shape and TODO queues because the simulation needs to
test protocol drafts as coherent artifacts, not as detached individual files.
Source: `DI-fakin`.

| Original path | New path | Method | Source commit | Notes |
|---|---|---|---|---|
| `protocols/group-session.d/` | `simulations/SIM-piloh-turns-149-208-recovery/protocols/group-session.d/` | `git mv` | `780f56525a8d528d3d5caf58ab18f9a7f41da892` | Draft protocol specimen; TODO queue moved with tree. |
| `protocols/udp-binding.d/` | `simulations/SIM-piloh-turns-149-208-recovery/protocols/udp-binding.d/` | `git mv` | `780f56525a8d528d3d5caf58ab18f9a7f41da892` | Draft protocol specimen; TODO queue moved with tree. |
| `protocols/ppx-dr.d/` | `simulations/SIM-piloh-turns-149-208-recovery/protocols/ppx-dr.d/` | `git mv` | `780f56525a8d528d3d5caf58ab18f9a7f41da892` | Legacy proposal/review protocol specimen; TODO queue moved with tree. |

The protocol specs remain drafts. Their simulation-local paths are provenance
for this experiment, not authoritative PromiseGrid API layout.
