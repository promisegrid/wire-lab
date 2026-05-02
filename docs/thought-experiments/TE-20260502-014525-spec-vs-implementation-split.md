# TE-32: Spec-side vs implementation-side split, and the `implementations/` top-level

**Status:** Draft. Amends TE-31 (does not retract). Refines TE-29.
**Drafted:** 2026-05-02 01:45 UTC

## Premise

TE-31 closed OQ-29.1 with Alt-G ("the simrepo's CHANGELOG names the
spec-doc-CID; no tree hash, no bundle"), but it conflated two
categories of file that both could plausibly live in a "simrepo":

- **Category A — Spec-side / design artifacts.** TEs, draft specs,
  design TODOs, DRs, per-protocol manifests. These are
  **WIP-toward-a-freeze**. They do not conform to anything; they
  *become* the spec they are heading toward.

- **Category B — Implementation-side / conformance artifacts.** Code,
  test vectors, conformance fixtures, simulator harness wiring. These
  **conform to** an already-frozen spec-doc-CID. They have an
  upstream; they make a promise about behaving consistently with that
  upstream.

TE-29 only contemplated category A under `protocols/<slug>.d/`. TE-31
introduced a CHANGELOG aimed at category B without saying where
category B lives, and used "simrepo" as if it covered both. This TE
fixes the conflation.

## The split

### Category A lives in `protocols/<slug>.d/` (unchanged from TE-29)

```
protocols/<slug>.d/
├── docs/thought-experiments/      (per-protocol TEs)
├── specs/<slug>-draft.md          (the WIP draft spec)
├── TODO/                          (per-protocol TODOs, per TE-30)
├── manifest.json                  (per-protocol release manifest)
└── CHANGELOG.md                   (freeze history of THIS spec's doc-CIDs)
```

When a freeze happens, the immutable sibling pair appears next to it
exactly as TE-29 §74-83 specified:

```
protocols/<slug>-<pcid>.md          (frozen spec doc, immutable)
protocols/<slug>-<pcid>.d/          (frozen design tree, immutable)
```

The live `<slug>.d/` keeps evolving toward the next freeze. The frozen
siblings never change.

### Category B lives in `implementations/<impl-name>/` (new top-level)

```
implementations/<impl-name>/
├── (whatever shape the implementation needs)
└── CHANGELOG.md                   (conformance-claim history)
```

Each implementation is a free-standing tree with its own internal
shape (Go module layout, ns-3 fixture scaffold, whatever). The only
universal requirement is a `CHANGELOG.md` at the implementation root
declaring which spec-doc-CIDs it implements.

`<impl-name>` is a human-readable slug chosen by the implementer.
Multiple implementations of the same protocol coexist as siblings:

```
implementations/
├── go-udp-binding-reference/
├── ns3-udp-binding-fixture/
├── go-group-session-reference/
└── go-ppx-dr-reference/
```

The mapping from implementation to protocol is not encoded in the
directory name (an implementation may implement multiple protocols);
it is encoded in the CHANGELOG entries.

### Why a top-level `implementations/`, not a subtree of `protocols/`

Three reasons.

1. **One implementation may implement many protocols.** A Go reference
   that ships UDP-binding plus group-session plus ppx-dr is one tree,
   not three. Putting it under `protocols/<slug>.d/impl/` would force
   either three copies or three symlinks, both ugly.

2. **External implementations have the same shape.** A third party's
   implementation lives in their own git repo with the same
   `CHANGELOG.md` convention. A top-level `implementations/` in this
   repo is just our local collection of implementations; the rest of
   the world's implementations live in their own equivalent
   directories elsewhere. The shape is location-independent.

3. **Spec-side and impl-side fork independently.** A spec may be
   frozen at v1 and v2 simultaneously (per C-4); independent
   implementations may target either. Coupling the directory shape
   would create a false hierarchy.

## CHANGELOG.md semantics, formally

Both categories have a `CHANGELOG.md`, but they say different things.
Same content-addressing primitive, different verbs, different
audiences.

### A-side CHANGELOG (freeze history of a spec)

Lives at `protocols/<slug>.d/CHANGELOG.md`. Each entry records a
freeze event: a moment when the WIP draft reached a state that the
maintainers wanted to publish as a stable doc-CID.

```
## 2026-06-15

- **event:** freeze
- **doc-cid:** bafkrei...spec-doc-cid-v1
- **title:** UDP-binding v1
- **notes:** First freeze of UDP-binding. Test vectors at
  bafkrei...vectors-pcid (referenced from spec doc body).

## 2026-09-03

- **event:** freeze
- **doc-cid:** bafkrei...spec-doc-cid-v2
- **title:** UDP-binding v2 (multi-stream support)
- **supersedes:** bafkrei...spec-doc-cid-v1 (not deprecated; both
  remain valid)
- **notes:** Adds multi-stream framing. v1 implementations remain
  conforming to v1.
```

`event` values:
- `freeze` — a doc-CID has been published.
- `withdraw` — a previously published doc-CID is being de-recommended
  (still resolvable, but the maintainers no longer endorse it).
- `note` — purely editorial; no doc-CID change.

The A-side CHANGELOG is **about the spec**, written by the spec
maintainers. It does not name implementations.

### B-side CHANGELOG (conformance-claim history of an implementation)

Lives at `implementations/<impl-name>/CHANGELOG.md`. Each entry
records a conformance claim: a moment when the implementation declared
its relationship to one or more spec-doc-CIDs.

```
## 2026-06-20

- **claim:** implements
- **spec:** bafkrei...spec-doc-cid-v1
- **scope:** full
- **notes:** Initial release. Passes test vectors referenced by
  the v1 doc.

## 2026-09-15

- **claim:** implements
- **spec:** bafkrei...spec-doc-cid-v2
- **scope:** full
- **breaking-change:** false
- **notes:** Migrated to v2 multi-stream support. Backward-compatible
  with v1 peers via runtime negotiation.

## 2027-04-01

- **claim:** deprecates
- **spec:** bafkrei...spec-doc-cid-v1
- **notes:** v1 is no longer tested; new connections refused. v2
  remains supported.
```

`claim` values (TE-31's set, unchanged):
- `implements`
- `partially-implements`
- `extends`
- `deprecates`

Plus the new **`breaking-change`** boolean (resolves TE-31 OQ-31.3 in
the affirmative — explicit flag is cleaner than inferring from a
preceding `deprecates`).

The B-side CHANGELOG is **about the implementation**, written by the
implementer. It names spec-doc-CIDs but does not name other
implementations.

### Symmetry

The two sides have a clean duality:

| | A-side | B-side |
|---|---|---|
| Lives in | `protocols/<slug>.d/CHANGELOG.md` | `implementations/<impl-name>/CHANGELOG.md` |
| Records | freeze events for a spec | conformance claims for an implementation |
| Names | spec doc-CIDs (its own) | spec doc-CIDs (upstream) |
| Authored by | spec maintainers | implementer |
| Promise verb | "I publish" | "I implement" |

Neither side names the other directly. Discovery from spec to
implementations (or vice versa) is a query over a known set of
CHANGELOGs, not a direct reference. This preserves TE-31's
universally-quantified promise property: a spec-doc-CID promises
behavior universally; an impl-CHANGELOG-entry accepts the promise on
behalf of one specific implementation.

## Amendment to TE-31

TE-31 stands; this TE does not retract it. Specifically:

- **The inversion thesis stands.** The reference goes
  implementation-to-spec, never spec-to-implementation. Both
  CHANGELOGs are consistent with that direction (the A-side names the
  spec's own previous CIDs; the B-side names upstream spec CIDs).

- **Alt-G stands as the freeze ceremony.** "A spec is frozen when its
  doc-CID is published and at least one CHANGELOG entry references it"
  needs one tiny refinement: the **A-side** CHANGELOG itself is the
  publication record. So the condition becomes:

  > A spec is frozen when its A-side CHANGELOG records a `freeze`
  > event with a doc-CID. B-side CHANGELOGs may then independently
  > claim conformance.

  This makes the freeze a one-sided act (the spec maintainers
  control it), which is structurally better than requiring an
  implementation to exist before the spec is "frozen" — RFC 768 was
  frozen on publication, not on first BSD release.

- **Where TE-31 said "simrepo," read "implementation."** Specifically
  these passages: lines 26-27 ("The simrepo points at the spec-doc
  CID"), lines 47-51 ("Each simrepo carries a CHANGELOG"), lines 142-
  149 (harness conformance check loop). The CHANGELOG-and-conformance
  language was always about the B-side; the A-side was conflated in
  by accident.

- **TE-31 OQ-31.2 (CHANGELOG location and format).** Refined here: the
  format is per side as specified above. Open sub-question: should the
  header block be YAML front-matter, fenced code, or HTML comment?
  Lean unchanged: fenced code with a `changelog-entry` info string.

- **TE-31 OQ-31.3 (breaking-change flag).** Resolved: yes, explicit
  flag on the B-side. See B-side schema above.

- **TE-31 OQ-31.4 (CHANGELOG entries as promises).** Stands; applies
  to both sides. Each entry is canonicalizable to bytes that have a
  pCID; v0 may skip cryptographic signing.

- **TE-31 OQ-31.5 (reverse index).** Refined: the index goes both
  directions — given a spec-doc-CID, find all B-side CHANGELOGs that
  reference it; given an implementation, find which specs it claims.
  Both indices are queries over known CHANGELOG sets; neither is
  embedded in the docs themselves.

## TODO 014 impact

TODO 014 (protocols-as-simulated-repos migration) was specified before
the A/B split was named. Its current steps either implicitly assume
everything goes under `protocols/<slug>.d/` (which would land
implementations there too) or are silent on category B. Steps to
add:

12. **Create top-level `implementations/` directory with stub
    `README.md`** explaining the A/B split, pointing at TE-32. Empty
    otherwise; no implementations exist yet at the time of migration.

13. **Update TODO 018 (UDP-binding v0 reference implementation)** to
    target `implementations/go-udp-binding-reference/` rather than
    any path under `protocols/udp-binding.d/`. The protocol's
    spec-side TODOs stay in `protocols/udp-binding.d/TODO/`; the
    implementation work is a separate concern that does not migrate
    with the protocol's design.

14. **Update TODO 019 (ns-3 harness scaffold)** similarly: target
    `implementations/ns3-harness-fixture/` (or whatever final slug
    we choose). Wire-lab harness reference impl is a B-side artifact,
    not a spec-side one.

15. **Add empty `CHANGELOG.md` stubs to each migrated
    `protocols/<slug>.d/`** with a placeholder explanatory header
    block. Real `freeze` events get added when the first frozen
    sibling appears.

These additions are mechanical and stay inside TODO 014's "happens
atomically with the migration" property.

## Surfaced questions

- **OQ-32.1: External implementations.** This repo's
  `implementations/` is just our local collection. Third parties have
  their own equivalent. Do we want a discovery mechanism — e.g. a
  `KNOWN_IMPLEMENTATIONS.md` registry that lists external
  implementation repos by URL or by some pCID-identified signed
  ledger? Lean: yes eventually, but not v0; the empty
  `implementations/README.md` can mention "external implementations
  exist; out of scope for v0 to discover them."

- **OQ-32.2: Implementation-of-multiple-protocols mapping.** A single
  Go binary might implement UDP-binding + group-session + ppx-dr.
  Should each of those CHANGELOG entries appear in one
  `implementations/go-reference/CHANGELOG.md`, or should the
  implementation tree be split into per-protocol subdirectories with
  per-protocol CHANGELOGs? Lean: one CHANGELOG with multiple
  `implements` entries (each entry already names which spec-doc-CID
  it claims; the protocol identity is fully captured by the doc-CID).

- **OQ-32.3: Test vectors as A or B.** Test vectors referenced from
  inside the spec doc by pCID are part of the doc's content (per the
  user's 2.C answer in TE-31). They are A-side. Test vectors written
  *by an implementation* to demonstrate its own behavior are B-side.
  No conflict, but worth naming to forestall confusion. Open: do we
  need a third top-level for "shared conformance vectors not yet
  promoted into a spec doc"? Lean: no; promote them into the spec
  doc when ready, otherwise keep them as drafts in the spec-side
  `<slug>.d/` until promoted.

- **OQ-32.4: Harness as protocol vs. harness as implementation.**
  TE-30 treated the wire-lab harness itself as a protocol
  (`protocols/wire-lab.d/`). The harness has both a design (A) and a
  reference implementation (B). The reference impl belongs in
  `implementations/wire-lab-harness-reference/`. Confirms TE-30's
  framing rather than complicating it.

- **OQ-32.5: Frozen-sibling tree size.** When a frozen
  `protocols/<slug>-<pcid>.d/` is created, what does it contain? Just
  the spec-side files, or also a snapshot of `implementations/`? Lean:
  spec-side only — the frozen sibling is the immutable freeze of the
  spec; implementations evolve independently and link back via their
  own CHANGELOGs.

## Closes / partially closes

- **Closes nothing previously open.** TE-32's role is to amend TE-31
  before downstream work (TODO 014, TODO 018, TODO 019) starts
  implementing the wrong shape.

- **Resolves TE-31 OQ-31.3** (breaking-change flag): yes, explicit
  flag on the B-side.

- **Refines TE-31 OQ-31.2** (CHANGELOG location and format): per side
  as specified above; format question (YAML / fenced / HTML comment)
  remains open.

- **Refines TE-31 OQ-31.5** (reverse index): bidirectional indices;
  neither embedded in docs themselves.

## Verdict

Adopt the A/B split. Add `implementations/` as a top-level. Patch
TODO 014 with steps 12-15. TE-31 stays on the books with the explicit
amendment that "simrepo" was always the B-side; the spec-side has its
own narrower CHANGELOG semantics (freeze history, not conformance
claim).
