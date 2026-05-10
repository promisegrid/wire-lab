# Simulations

`simulations/` holds bounded experiment worlds for deriving PromiseGrid design
choices from wire-lab evidence. A simulation may contain candidate protocol
specimens, archived historical inputs, concrete world state, observations, and
results, but it is not itself the PromiseGrid development guide or a final
PromiseGrid node layout.

Root wire-lab apparatus stays rooted unless a later DI says otherwise. In the
current Mupoz split, `protocols/wire-lab.d/` remains the harness apparatus home,
while candidate PromiseGrid protocol trees and concrete specimens move into a
named simulation until results graduate through DR, DI, frozen specs, guide
prose, or a future PromiseGrid spec corpus. Source: `DI-pakid`; `DI-fakin`.

## Current simulations

| Simulation | Purpose | Status |
|---|---|---|
| `SIM-piloh-turns-149-208-recovery/` | Recovery/dogfood simulation for the turns 149-208 context-loss recovery slice, candidate protocol specimens, legacy proposal records, and wire-lab-devs transport evidence. | Active specimen boundary |
