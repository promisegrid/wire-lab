# SIM-pamap-grid-envelope-signable-view-atproto-like: Atproto-like payload signable-view probe

This simulation is a standalone grid-envelope specimen. It keeps the universal
outer shape minimal as `grid([pcid, payload])`, but makes the payload contract
itself carry a reserved proof slot and a named explicit signable view. The
point is to test the simplest "signature carried inside the payload contract"
variant without adding universal outer proof slots. Source: `DI-nohir`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.

Primary comparison targets: `SIM-gojot`, `SIM-riliz`, and `SIM-janov`.
