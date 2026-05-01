# TODO 013 - Transports carve-out

Carve the original combined `channels/` material into a thin outer wire-lab spec (`specs/transport-spec-draft.md`) plus a substantive first transport-protocol spec (`specs/group-transport-draft.md`), perform the `channels/` → `transports/` rename across the repo, and rename the TE-24 / DR-009 / TODO-012 trio to use the `group-transport-envelope` slug because their decisions are properties of the group-transport-protocol, not of the outer wire-lab transport-spec.

This TODO is the lock record for the carve-out. The substantive design decisions live in TE-26 (transport-protocol types and the four locked principles), TE-27 (transports rename and per-axis meta-rule), and TE-24 (group-transport envelope), all on `ppx/main`.

## Subtasks

- [x] 013.1 Rename `channels/` → `transports/`; rewrite `transports/README.md` as a thin pointer to the two specs.
- [x] 013.2 Create `specs/transport-spec-draft.md` (thin outer rule: `transports/<pcid>--<slug>/` keying convention, no `Transport:` header, code-as-handler, per-axis meta-rule).
- [x] 013.3 Create `specs/group-transport-draft.md` (substantive v0 contract for the group-transport-protocol class: N≥2, all-see-all, multi-writer DAG, append-only).
- [x] 013.4 Rewrite TE-24 in place under transport vocabulary; rename file to `TE-20260430-204108-group-transport-envelope.md`.
- [x] 013.5 Rewrite DR-009 to transport vocabulary; rename file to `DR-009-20260430-204108-group-transport-envelope.md`.
- [x] 013.6 Rewrite TODO 012 to transport vocabulary; rename file to `TODO/012-group-transport-envelope.md`; add freeze-gate subtasks (012.7, 012.8).
- [x] 013.7 Update `TODO/TODO.md` to reference the renamed TODO 012.
- [x] 013.8 Update `docs/thought-experiments/README.md` TE-24 row to use the new title and filename.
- [x] 013.9 Update `specs/harness-spec-draft.md` §8 TE-24 bullet to transport vocabulary; refresh TE-26/TE-27 cross-references; rewrite Open Question #19 to use transport vocabulary and reference both the outer transport-spec and the group-transport-draft.
- [x] 013.10 Run `tools/spec check`; commit; merge `--no-ff` to `ppx/main`; push; delete twig.

## Decision Intent Log

This TODO does not introduce a new DI; it executes the locked decisions of TE-26 (DF-26.7 Alt-7.C, DF-26.8 Alt-8.C, DF-26.10 Alt-10.A) and TE-27 (DF-27.1 Alt-1.A, DF-27.2 Alt-2.A, DF-27.4 per-axis meta-rule, DF-27.5 Alt-5.B), plus the locked group-transport DFs from this session: DF-T1 Alt-T1.A (`Parents:` header), DF-T2 Alt-T2.A (single line, space-separated), DF-T3 Alt-T3.B (always optional), DF-T4 Alt-T4.A (`Parents:` only; no `Prev-Message-CID:`), DF-T5 (ack in body, no `IHave:` header), DF-T6 Alt-T6.A (flat subdirectory layout), DF-Q1 Alt-Q1.A (drop `Kind:` header), DF-Q2 (cumulative-prefix/frontier ack deferred to its own future TE).

## Why a separate TODO

The carve-out crosses too many files for the existing TODO 012 to absorb cleanly:

- TODO 012 is the lock record for the group-transport envelope decision (DI-009-20260430-204108). Its subtasks are about that decision and its freeze gate.
- TODO 013 is the rename and spec-split mechanic. It does not introduce a new decision; it implements the consequences of TE-26 and TE-27.

Keeping the two records distinct lets the lock-of-decision (012) and the mechanical carve-out (013) be reasoned about independently.

## Notes

- The TE-24 / DR-009 / TODO-012 file renames went `grid-pcid-channel-carrier` → `group-transport-envelope` because those documents' load-bearing decisions are properties of the group-transport-protocol class, not of the outer wire-lab transport-spec. Future transport-protocol classes (ring, gossip, hub-mediated, large-N, ephemeral) will produce their own envelope decisions in their own spec docs, with their own TE / DR / TODO trios as needed.
- TE-25 (numbering collision) and TE-26 (transport-protocol types) retain "channel" in places: TE-25 is a historical record of the channels-branch reconciliation; TE-26 has a vocabulary-note prefix explaining the rewrite and otherwise uses transport vocabulary throughout. Both filenames retain "channel" because the timestamp slugs are content-addressable identifiers and renaming them would break the integer-anchors-on-first-drafted-timestamp invariant from TE-25.
- A future TE will address cumulative-prefix / frontier-style acknowledgement (Q2 from this session). v0 of the group-transport-protocol uses per-message body-level acknowledgement; that scheme will likely need a compact form once message volumes grow.
- Anticipated future transport-protocol specs per TE-27: ring, star, cluster-of-clusters, gossip; plus a future TE on transport-protocol migration semantics.
