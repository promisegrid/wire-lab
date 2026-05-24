# SIM-dalor-grid-envelope-protocol-owned-signature-slot: Protocol-owned outer signature-slot probe

This simulation is a standalone grid-envelope specimen. It tests a three-slot
outer envelope `grid([pCID, payload, signature])` where the outer third slot is
mandatory but the proof family is still owned by the protocol named by `pCID`
rather than by a
separate outer `sig_pcid`. The point is to test the smallest explicit
outer-signature design that still allows varsig-, multisig-, or other
pCID-defined proof rules. In Promise Theory terms, the current sender's
signature is evidence of the sender's promise that the payload bytes are shaped
according to the protocol specification named by the `pCID`; higher-layer
promise accounting remains inside the payload protocol. Source: `DI-kukuk`;
`DI-pozom`.

The local draft spec is
`protocols/grid-envelope.d/specs/grid-envelope-draft.md`.

Primary comparison targets: `SIM-kurim`, `SIM-jufag`, `SIM-pamap`, and
`SIM-jumav`.
