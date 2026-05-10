# World

`world/` contains active simulated state for this recovery/dogfood experiment.
The subdirectories use the Pahah/Nizor vocabulary for bounded simulation state;
they are not final PromiseGrid node layout. Source: `DI-fakin`.

| Directory | Purpose |
|---|---|
| `sites/` | Local participant views when the simulation needs them. |
| `groups/` | Group-level manifests, rosters, frontiers, or semantic views when the simulation needs them. |
| `cas/` | Content-addressed storage sketches or fixtures when the simulation needs them. |
| `feeds/` | Feed/substrate observations when the simulation needs them. |
| `wires/` | Wire/carrier observations when the simulation needs them. |
| `transports/` | Current transport-message specimen data imported from the old root `transports/` tree. |
