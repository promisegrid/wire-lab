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
| `SIM-ludut-wire-lab-devs/` | Concrete wire-lab-devs world evidence and transport specimen state for replay and migration provenance. | Active concrete-world lineage |
| `SIM-rakot-group-session/` | Independent group-session lineage carrying session-envelope protocol drafts and TODO ownership. | Active protocol lineage |
| `SIM-ludaf-udp-feed/` | Independent UDP feed lineage (renamed from legacy `udp-binding` active-tree naming). | Active protocol lineage |
| `SIM-labit-feed-outer/` | Independent thin outer-feed lineage, including extracted feed-outer draft material. | Active protocol lineage |
| `SIM-kurim-grid-envelope/` | Independent `grid([pcid, payload])` lineage and successor-owner queue. | Active protocol lineage |
| `SIM-hugoj-cas-usenetlike-gitlike/` | Broad design exploration of a CAS + Usenet-like + git-like PromiseGrid lineage, with `group-session` treated as one current specimen rather than the whole subject. Source: `DI-pijun`. | Active design exploration |
