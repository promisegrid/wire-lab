# TE-31: Spec-doc as upstream, simrepo as implementation — inverting the conformance reference

## Status

decided

(Originally marked Draft. Closes OQ-29.1. Drafted 2026-05-02 00:49 UTC. Supersedes within OQ-29.1 Alt-A through Alt-F enumerated in TE-29. The body's Verdict adopts the inversion.)

## Premise

TE-21 framed the spec doc as a promise. TE-22 sketched a content-addressed
store for spec docs. TE-29 promoted each protocol to a simulated repo
("simrepo") under `protocols/<slug>.d/` and deferred the freeze ceremony
question. The clarifying answers from the user are:

- **1.A:** the spec doc is RFC-shaped — a complete contract; hand someone
  the doc (or its pCID) and they can verify conformance.
- **2.C:** hybrid — small machine-checkable bits (state tables, schemas,
  ABNF) live inline in the doc; larger artifacts (test vectors, example
  pcaps) are companion files referenced **by pCID** from inside the doc.

Together, those answers force a structural inversion that TE-29's
Alt-A..Alt-F enumeration did not contemplate. This TE makes the
inversion explicit and uses it to close OQ-29.1.

## The inversion in one sentence

> The simrepo points at the spec-doc CID; the spec doc never points at
> the simrepo.

## The two mental models

### Old mental model (implicit in TE-29 Alt-A through Alt-F)

```
spec-doc ──(embeds)──▶ tree-hash ──(names)──▶ simrepo
```

The spec doc owns the bundle. To freeze a version, you compute a tree
hash over `protocols/<slug>.d/` and either embed it in the doc or
publish it alongside. The doc and the repo are coupled: one cannot
ship without referring to the other.

### New mental model (this TE)

```
simrepo ──(CHANGELOG entry)──▶ spec-doc CID
```

The simrepo owns its conformance claim. Each simrepo carries a
`CHANGELOG` whose entries each name a spec-doc CID and a conformance
verb. The spec doc itself is upstream and oblivious — it does not know
which simrepos exist, just as RFC 768 does not know which UDP stacks
exist.

## Why this is right

1. **It matches how RFCs work.** RFC 793 does not list TCP
   implementations. Each TCP implementation declares which RFC it
   implements. The reference goes implementation-to-spec, never
   spec-to-implementation. PromiseGrid gains nothing by inverting that.

2. **One spec-doc CID can have many implementations without
   republishing.** Under the old model, every new conforming simrepo
   would either need to be added to the doc's bundle (recomputing the
   tree hash and producing a new doc-CID) or live outside the bundle
   (in which case the bundle was never the source of truth anyway). The
   new model makes "many implementations of one spec" a first-class
   shape.

3. **An implementation can evolve its conformance claim cheaply.** A
   simrepo migrating from `implements: bafkrei...v1` to
   `implements: bafkrei...v2` writes one CHANGELOG entry. The spec
   author is not involved.

4. **Companion files travel with the doc automatically.** The 2.C
   answer puts test vectors and example pcaps as separate files
   referenced by pCID from inside the doc. Because the doc cites them
   by content address, those companion CIDs are part of the doc's
   transitive content. Anyone resolving the doc-CID can resolve its
   companions. No separate bundle hash is required.

5. **The spec doc becomes a proper promise object.** TE-21 said the
   spec doc is a promise. A promise that names its own implementations
   is not a promise — it is a release record. By keeping the doc
   ignorant of implementations, it stays a clean
   "I promise that any implementation conforming to the bytes of this
   document behaves as described" — i.e., a universally-quantified
   promise over implementations, which is exactly what an RFC is.

6. **It dissolves the bundle as a first-class object.** Under the old
   model, the unit of release was a bundle (doc + repo + tree hash).
   Under the new model, the unit of release is just a doc-CID; the
   simrepo's CHANGELOG is internal bookkeeping. Fewer first-class
   objects, less ceremony.

## The CHANGELOG format (sketch)

Each simrepo has a `CHANGELOG.md` at its root. Each entry has a
machine-readable header block followed by human prose:

```
## 2026-05-02

- **claim:** implements
- **spec:** bafkrei...spec-doc-cid-v1
- **scope:** full
- **notes:** First conforming release. Passes test vectors
  bafkrei...vectors-pcid (referenced from the spec doc itself).
```

`claim` is one of:

- `implements` — every machine-checkable bit of the spec doc passes.
- `partially-implements` — declare which sections; the rest are TODO.
- `extends` — implements the spec plus additional behavior; the
  extension MUST be backward-compatible with the spec.
- `deprecates` — the implementation now refuses to interoperate with
  this spec-doc-CID and points to a successor doc-CID via a separate
  `implements` entry.

Conformance verbs are deliberately limited. "Loosely inspired by" and
"based on" are not conformance claims; they are prose.

## What "frozen" means under the inversion

> **Alt-G:** The spec is frozen when its doc-CID is published and at
> least one simrepo CHANGELOG entry references it. There is no tree
> hash. There is no bundle. The doc-CID alone is the freeze.

This supersedes Alt-A through Alt-F as enumerated in TE-29 §326.
Alt-A..Alt-F all assumed the doc-and-repo bundle as the unit of
release; once the inversion is accepted, the bundle ceases to exist
and the freeze ceremony collapses to "publish the doc."

The doc-CID is, by content addressing, immutable. "Publish" here means
"resolvable from the spec-doc store described in TE-22." A draft doc
with a tentative pCID is not yet published; once any simrepo writes a
CHANGELOG entry naming that pCID, the pCID is effectively frozen
because the conformance claim now exists in the world.

## How the harness checks conformance

The harness (itself a protocol per TE-30) loads the simrepo, reads its
CHANGELOG, fetches each spec-doc-CID named under `implements`, parses
the doc's machine-checkable bits (state tables, schemas, ABNF, test
vector references), and runs them against the simrepo's code. The
harness reports `conforms`, `fails`, or `partial` per spec-doc-CID.

The doc itself remains markdown. The harness extracts machine-checkable
content via documented fenced-code-block conventions (e.g. ``` ```abnf
```, ``` ```state-table ```, ``` ```vectors ``` referencing companion
pCIDs). The conventions are themselves part of the harness spec.

## Awkward cases handled

### Two simrepos disagree about the same doc-CID

Treat as a **discrepancy** (per the user's standing-rule vocabulary).
Surface via the harness; resolve by the spec-doc author publishing a
clarifying successor doc-CID, after which simrepos can migrate at
their own pace.

### A simrepo claims `implements: X` but does not

The harness detects this. The CHANGELOG claim is itself a promise; a
broken claim is a broken promise, surfaced exactly as any other
promise breach.

### A simrepo wants to claim conformance to multiple specs

Multiple `implements` entries in the CHANGELOG, possibly with `scope`
qualifiers. The harness checks each independently.

### The spec doc references companion files; one of them is wrong

The companion pCID is part of the doc's content. To fix it, the spec
author publishes a new doc-CID with corrected companion references.
The old doc-CID remains valid forever (content-addressed); simrepos
choose when to migrate.

### A breaking change between doc-CID v1 and v2

A simrepo migrating from v1 to v2 writes a CHANGELOG entry with
`claim: deprecates spec: bafkrei...v1` followed by
`claim: implements spec: bafkrei...v2`. Optionally a `breaking-change:
true` flag in the v2 entry. Downstreams pinned to v1 keep working
against the old doc-CID; downstreams that follow the simrepo's tip
must be ready for v2 semantics.

## Self-consistency with PromiseGrid's promise model

The inversion is structurally identical to PromiseGrid's own
promise/acceptance shape:

- The **spec-doc CID** is a promise: "I promise that any implementation
  conforming to the bytes of this document behaves as described."
- The **CHANGELOG entry** is an acceptance: "I promise that this code
  keeps the spec-doc-CID's promise."

The conformance reference is itself a promise relationship. That is
satisfyingly self-consistent — the meta-protocol (how we describe
protocols) uses the same primitive as the protocols it describes.

## Relationship to TE-22 (spec-doc store) and TE-30 (per-protocol TODOs)

TE-22 stands; the spec-doc store described there is exactly the
mechanism by which doc-CIDs become resolvable. The inversion does not
change the store, only the direction of reference at the simrepo
boundary.

TE-30 stands; per-protocol `TODO/` subtrees are orthogonal to the
conformance-reference direction. A simrepo's TODO/ tracks its work;
the CHANGELOG tracks its public conformance claims. They do not
overlap.

## Closes / partially closes

- **Closes OQ-29.1.** The freeze ceremony is Alt-G as defined above:
  publish the doc-CID; the simrepo CHANGELOG names it. No tree hash,
  no bundle.

## Surfaces

- **OQ-31.1: Harness self-application.** TE-30 made the harness itself
  a protocol (`protocols/wire-lab.d/`). Under the inversion, the
  harness has a spec-doc and a simrepo, and the simrepo's CHANGELOG
  names the harness-spec-doc-CID. That is meta but not paradoxical;
  the harness's own conformance claim can be checked by any other
  conforming harness implementation. Worth a smoke test once we have
  two harness implementations (e.g. Go reference + ns-3 fixture).
  Lean: works; document the recursion explicitly in the harness spec.

- **OQ-31.2: CHANGELOG location and format.** Proposal:
  `protocols/<slug>.d/CHANGELOG.md` at simrepo root, with the
  machine-readable header block sketched above. Open: should the
  header block be YAML front-matter, fenced code, or a `<!-- -->`
  comment? Lean: fenced code with a `changelog-entry` info string,
  consistent with how 2.C handles other machine-checkable bits.

- **OQ-31.3: Migrating across breaking changes.** Sketched above
  (`deprecates` then `implements` of successor). Open: do we want a
  `breaking-change: true` flag on the new entry, or is the
  `deprecates` of the old version enough signal? Lean: explicit flag,
  to make harness reporting cleaner.

- **OQ-31.4: CHANGELOG entries as promises.** A CHANGELOG entry is
  itself a promise (per the self-consistency argument above). Should
  each entry have its own pCID, signed by the simrepo's identity, so
  that a third party can verify the conformance claim without trusting
  the simrepo's git history? Lean: yes; v0 specs may skip signing,
  but the entry should be canonicalizable to bytes that have a pCID.

- **OQ-31.5: Spec-doc author's awareness of implementations.** The
  inversion says the doc is oblivious to implementations, but the
  spec-doc store (TE-22) could optionally maintain a reverse index
  for human convenience: "which CHANGELOGs name this doc-CID?" That
  index is metadata, not part of the doc's content, so it does not
  contaminate the doc's promise. Lean: build the reverse index as a
  query over a known set of simrepos; do not embed it in the doc.

## Why not just use git tags or releases?

Git tags point at trees, not at content-addressed promises. A git tag
named `v1.0` tells you nothing about whether the tagged tree conforms
to anything. The CHANGELOG-entry-naming-doc-CID approach makes the
conformance relationship explicit and content-addressable; git tags
remain useful as human-readable shortcuts but are not the conformance
mechanism.

## Why not the W3C "REC" model?

W3C RECs do name implementations (via implementation reports), which
is closer to the old mental model. PromiseGrid does not adopt that
because the spec-doc-as-promise framing (TE-21) prefers the doc to be
universally quantified over implementations; naming specific
implementations weakens the universal quantifier into an existential
one and forces the doc to be republished whenever the implementation
list changes.

## Implementation impact

TE-31 is design-only; no immediate code changes. Follow-on work:

1. **TODO 020 (forthcoming).** Define the CHANGELOG.md format
   precisely (resolve OQ-31.2 and OQ-31.3) and write the harness
   parser for it.
2. **TODO 014.** Continue as-is; the per-protocol simrepo migration
   does not depend on TE-31. Once 014 lands, each simrepo gets a
   stub CHANGELOG.md with the first `implements` entry pointing at
   the current draft spec-doc.
3. **Harness spec §8.** Add this TE to the bibliography (chronological
   order, after TE-30).
4. **TE-22 update.** Note that the spec-doc store gains an optional
   reverse index per OQ-31.5; do not block on it.

## Verdict

Adopt the inversion. Close OQ-29.1 with Alt-G. Defer the CHANGELOG
format details to TODO 020 and a follow-on TE if needed.
