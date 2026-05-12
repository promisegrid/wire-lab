# SIM-piloh: Turns 149-208 recovery

This simulation is the first Pahah/Mupoz recovery world. It gathers the
candidate PromiseGrid protocols, legacy proposal records, and concrete
wire-lab-devs transport evidence that were previously mixed into root-level
paths, then uses them as specimens for recovery and design decisions.

The simulation exists to support PromiseGrid design, not to document wire-lab as
the product. Guide writers should treat this directory as provenance and
experimental evidence. Settled PromiseGrid guide prose, frozen protocol specs,
or later DR/DI decisions can cite the results, but simulation-local paths are
not final PromiseGrid API layout. Source: `DI-fakin`.

## Contract

- `QUESTION.md` states the question under test.
- `concerns.md` maps recovery concerns to simulation inputs and expected result
  evidence.
- `protocol-set.md` lists protocol specimens and their source paths.
- `world/` holds active simulated world state.
- `archive/` holds historical material that the simulation preserves but does
  not treat as active world state.
- `events/`, `observations/`, and `results/` separate what happened, what was
  noticed, and what should feed later DR/DI/spec work.
- `decisions.md` records graduation and handoff status.
