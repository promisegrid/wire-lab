# SIM-jumav-grid-envelope-wrapper-proof-ceramic-like: Ceramic-like wrapper-proof probe

This simulation is a standalone grid-envelope specimen. It keeps the universal
outer shape minimal as `grid([pcid, payload])`, but makes the payload a
proof-bearing wrapper that signs a linked content object rather than embedding a
signature in the exact same bytes being signed. Source: `DI-nohir`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.

Primary comparison targets: `SIM-gojot` and `SIM-maraz`.
