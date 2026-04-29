# TODO 011 - TE-22 Spec-doc store layout and pCID machinery: drive to DI

Track the work to drive TE-22 (`docs/thought-experiments/TE-20260429-175530-spec-doc-store-and-pcid-machinery.md`) from `needs DF` to a set of decided DIs that lock the operational machinery for spec-doc storage, pCID computation, freezing, and cross-referencing in the wire-lab.

## Already-locked inputs

These were locked by Steve in chat 2026-04-29 before TE-22 was written. They are inputs to TE-22, not under DF in TE-22 itself:

- pCID format = CIDv1 (multibase + multihash + codec wrap).
- Layout = flat (`specs/<slug>-draft.md` and `specs/<slug>-{cidv1}.md` side-by-side; single `specs/MANIFEST.md`).
- Draft cross-refs = strict (a draft citing another spec MUST cite a frozen pCID, never another draft).
- Self-reference = external only (a frozen spec file does not contain its own pCID).

These four become DI entries in this file when work begins; they are decided and only need formal recording.

## Subtasks

- [ ] 011.1 Steve answers DF-22.1 (hash input). Recommended: 1.d (raw bytes + machine-checked formatter).
- [ ] 011.2 Steve answers DF-22.3 (freezing mechanic). Recommended: 3.d (snapshot file + manifest entry, no git tag).
- [ ] 011.3 Steve answers DF-22.4 (manifest format). Recommended: 4.d (single Markdown file with fenced YAML inside).
- [ ] 011.4 Steve answers DF-22.5 (freezing trigger). Recommended: 5.d (manual ritual via `tools/freeze-spec.sh` + scheduled CI audit).
- [ ] 011.5 Once 011.1-011.4 land, write a DI for each into this file. Also write DIs for the four already-locked inputs (CIDv1, flat layout, strict draft cross-refs, external-only self-reference).
- [ ] 011.6 Pin the CIDv1 parameter set explicitly: multibase, multihash function, codec. Recommended starting point: `multibase=base32`, `multihash=sha2-256`, `codec=raw`. This is a sub-DF inside 011.1; if Steve has a different parameter preference, capture it as a DI.
- [ ] 011.7 Genesis freeze of `harness-spec.md`: move to `specs/harness-spec-draft.md`, mint pCID, snapshot to `specs/harness-spec-{cidv1}.md`, append manifest entry, update all in-repo references. Separate twig.
- [ ] 011.8 Implement `tools/freeze-spec.sh` and `tools/check-spec-format.sh` per the locked DF answers. Pin the formatter version inside the scripts.
- [ ] 011.9 Implement the CI audit step (manifest-vs-disk consistency, cross-ref-citation lint, format check). Wire it into whatever CI exists today; design the script so it can run as a git pre-receive hook on a non-GitHub host.
- [ ] 011.10 After genesis freeze lands, open a follow-on TE on peer-level adoption metadata (the wire shape of "I, peer P, promise to behave as pCID X with open-question answers Q7=yes, Q9=variant-B"). This is the missing half of TE-21 Alt-E that TE-22 did not address.

## Decision Intent Log

(No DI yet. DF answers from Steve will populate this section. Already-locked inputs from chat 2026-04-29 will be formalized as DIs once work on subtasks 011.1-011.5 begins.)

## Notes

- TE-22 carries the full alternative analysis (Alt-1.A-D, Alt-3.A-D, Alt-4.A-D, Alt-5.A-D) and six scenarios (S1-S6). This file does not duplicate that analysis; it tracks the decision-driving work.
- The recommended set is `all-D: (1.d, 3.d, 4.d, 5.d)` per the TE. Reason: each D answer survives migration off GitHub, keeps the freeze act deliberate and auditable, and makes the manifest a machine-walkable structure rather than just a human convenience.
- TE-22 is the operational follow-on to TE-21. TE-21 said *what* a spec doc is (a layered promise); TE-22 says *how* the repo handles such docs (freeze, hash, store, cite, supersede). The DI entries from both TEs should land in the same revision of `harness-spec.md`'s vocabulary section.
- Linked DR: to be created as DR-011 once subtasks 011.1-011.4 begin landing. For now, the TE itself stands as the open-question record.
- Companion TODOs: TODO 006 (DI-provenance backfill) and TODO 007 (DR backfill for §11) become more concrete after TE-22 locks because the genesis freeze produces the first content-addressed reference points for those backfills.
- The "External-only self-reference" lock has a subtle implication: the freeze ritual must produce a frozen file whose bytes do not reference its own pCID. The simplest implementation is to ensure no pCID-of-self placeholder exists in `specs/<slug>-draft.md` at freeze time. Document this as an explicit step in `tools/freeze-spec.sh`.
