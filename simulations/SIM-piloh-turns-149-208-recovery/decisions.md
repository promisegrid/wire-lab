# Decisions and Graduation

This simulation is governed by `DI-pakid` and `DI-fakin`.

## Locked inputs

- `DF-mupoz.1 = 1.A`: keep apparatus rooted and move candidate protocol trees
  plus concrete specimens into the first simulation.
- `DF-mupoz.2 = 2.A`: migrate `transports/wire-lab-devs-draft/` with manifest
  and CID evidence.
- `DF-mupoz.3 = 3.A`: root `protocols/` contains only `wire-lab.d`.
- `DF-mupoz.4 = 4.C`: move all legacy proposal records under simulations.
- `DF-mupoz.5 = 5.A`: simulation results feed DR/DI/spec/dev-guide handoff.

## Graduation rule

Results from this simulation do not directly rewrite authoritative PromiseGrid
layout. A result graduates only when a later DR, DI, frozen spec, PromiseGrid
Development Guide change, or external PromiseGrid spec corpus records the
decision. Source: `DI-fakin`.
