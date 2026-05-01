# Transports

This directory is the wire-lab's simulation surface for the network being studied. Each subdirectory under `transports/` is one simulated transport instance: a particular way bytes-shaped-as-messages traverse a particular set of participants under a particular set of routing rules.

The wire-lab vocabulary is **transports** and **messages.** "Channel" is not used; if a logical-addressing layer above transports is later needed, a future TE introduces it. (See [TE-27](../docs/thought-experiments/TE-20260501-021921-transports-rename-and-axes-of-differentiation.md).)

## How transport directories are named

Transport directories are keyed by the pCID of the transport-protocol they speak:

```
transports/<pcid>--<slug>/
```

The pCID is canonical protocol identity; the slug is a human-readable suffix tools ignore. This convention is established by `specs/transport-spec-draft.md` and is invariant across all transport-protocols. (See [TE-26](../docs/thought-experiments/TE-20260430-215624-channel-transport-types-and-threaded-replies.md) DF-26.7 Alt-7.C.)

## Specs that govern this directory

- **`specs/transport-spec-draft.md`** is the *thin outer rule:* it defines the `transports/<pcid>--<slug>/` keying convention, the rule that messages do not declare their transport via a header, and the requirement that each transport-protocol-pCID names a spec defining the directory's interior.
- **`specs/group-transport-draft.md`** is the *first transport-protocol spec:* it defines the small-finite-closed-group transport-protocol (N≥2, all-see-all, multi-writer DAG of messages). The first concrete instance, when minted, is the Codex↔Perplexity case (N=2). The full v0 message contract — headers, parents, receipts, message kinds, canonical-bytes — lives in this spec.

A transport's interior layout, header set, parent semantics, receipt format, and message-kind vocabulary are properties of the *transport-protocol* that transport speaks, not of the wire-lab spec. The code that reads a transport's directory structure is the handler for that pCID.

## Status

No transport directories are created yet. The wire-lab is still in design mode; the first instance will appear when the group-transport-protocol-pCID is minted (frozen) and the first real Codex↔Perplexity exchange happens. Until then, this directory exists only as a design surface.

## Related design docs

- [TE-24: Group-transport envelope: `grid <pcid>` carrier, canonical bytes, and explicit promise body](../docs/thought-experiments/TE-20260430-204108-group-transport-envelope.md) — the v0 envelope choice for the group-transport-protocol; source document for the substantive group-transport contract.
- [TE-26: Transport-protocol types, pCID-keyed transport paths, and DAG message graphs](../docs/thought-experiments/TE-20260430-215624-channel-transport-types-and-threaded-replies.md) — the four locked principles for what's in `transports/`.
- [TE-27: Transports rename and axes of transport-protocol differentiation](../docs/thought-experiments/TE-20260501-021921-transports-rename-and-axes-of-differentiation.md) — the per-axis meta-rule for distinguishing transport-protocols and the rationale for the rename.
- [DR-009](../DR/DR-009-20260430-204108-group-transport-envelope.md) — the active decision request governing the group-transport envelope.
- [TODO 012](../TODO/012-group-transport-envelope.md) — the active TODO tracking group-transport envelope work.
- [TODO 013](../TODO/013-transports-carveout.md) — the carve-out TODO that performed the `channels/` → `transports/` rename and the spec carve-out.
