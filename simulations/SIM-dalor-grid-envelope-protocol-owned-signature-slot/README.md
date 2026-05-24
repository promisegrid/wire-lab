# SIM-dalor-grid-envelope-protocol-owned-signature-slot: Protocol-owned outer signature-slot probe

This simulation is a standalone grid-envelope specimen. It tests a three-slot
outer envelope `grid([pcid, payload, signature])` where the outer third slot is
mandatory but the proof family is still owned by `pcid` rather than by a
separate outer `sig_pcid`. The point is to test the smallest explicit
outer-signature design that still allows varsig-, multisig-, or other
pCID-defined proof rules. Source: `DI-kukuk`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.

Primary comparison targets: `SIM-kurim`, `SIM-jufag`, `SIM-pamap`, and
`SIM-jumav`.
