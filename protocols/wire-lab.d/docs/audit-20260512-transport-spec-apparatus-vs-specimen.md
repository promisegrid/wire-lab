# Audit: transport-spec-draft.md — apparatus vs. specimen

**Date:** 2026-05-12
**Target:** `protocols/wire-lab.d/specs/transport-spec-draft.md`
**Status:** Companion audit for turn-159 cleanup.
**Source:** `DI-mugar`.

## Framing

Turn 159 flagged `transport-spec-draft.md` as a likely companion target for
the same apparatus-vs-specimen split applied to `harness-spec-draft.md`. That
concern was valid when the original harness audit was written, because the
rooted transport draft still mixed thin outer transport/feed specimen material
with rooted wire-lab governance residue.

Later `rusis.11` work already performed the transport split under `DI-huzor`.
The specimen-side outer-feed material now lives in
`simulations/SIM-labit-feed-outer/protocols/feed-outer.d/specs/feed-outer-draft.md`.
The rooted `transport-spec-draft.md` now says it retains only the
apparatus/governance residue that remains after extraction.

## Current classification

| Section | Classification | Disposition |
|---|---|---|
| `## Purpose` | Apparatus/governance residue | States that specimen-side outer-feed material was extracted under `rusis.11`; no active specimen contract remains rooted here. |
| `## Sources` | Apparatus provenance | Records which TEs and DRs informed transport/feed history; TE-hogus is explicitly not a constraint on the rooted outer spec. |
| `## The per-axis meta-rule (TE-junil)` | Apparatus/meta-rule | Describes when distinct transport-protocol pCIDs are warranted; it is a classification rule, not a concrete transport specimen. |
| `## Open questions` | Apparatus governance | Tracks unresolved discovery/migration questions for future TE/DR work. |
| `## Freeze gate` | Apparatus governance | Describes how this rooted draft would freeze if it remains useful; does not make a specimen canonical. |

## Conclusion

No additional `transport-spec-draft.md` edit is required for turn-159 scope.
The companion audit item `UT-159.b` can close because the earlier suspicion has
been resolved by the `DI-huzor` extraction plus this verification note.

This audit does not freeze `transport-spec-draft.md`, does not make
feed-outer canonical, and does not create a new rooted specimen home.
