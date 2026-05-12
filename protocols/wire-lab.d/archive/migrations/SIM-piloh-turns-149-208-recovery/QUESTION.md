# Question

Can wire-lab recover the turns 149-208 design context and continue deriving
PromiseGrid design choices by placing candidate protocol specimens, legacy
proposal records, and concrete transport evidence inside one bounded simulation
rather than preserving obsolete root paths?

## Success criteria

- Root `protocols/` contains only `wire-lab.d`.
- Candidate PromiseGrid protocol trees are available under this simulation as
  specimens with source-path and source-commit provenance.
- The legacy `proposals/` mechanism is archived here, while current guide-writer
  discovery runs through `DEV-GUIDE-RESOURCES.md` and the external guide
  feedback process.
- The wire-lab-devs transport messages remain byte-identical, CID-verifiable
  specimens under `world/transports/`.
- Simulation results feed DR/DI/spec/dev-guide handoff instead of directly
  becoming authoritative PromiseGrid layout. Source: `DI-fakin`.
