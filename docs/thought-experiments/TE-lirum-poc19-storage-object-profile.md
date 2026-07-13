# TE-lirum: POC19 storage object profile

## TE ID

TE-lirum

## Status

decided

## Decision under test

`vumas.17` asks how POC19 should resolve the storage object profile before
durable stores are scaffolded:

> Should POC19 preserve POC18's split `chunks/*.bin` and `objects/*.cbor`
> filesystem layout, move to one CID-keyed CAS namespace, split source-of-truth
> directories by profile, or make all retained objects DAG-CBOR?

This decision affects local filesystem layout, stage0 fetch/verify behavior,
stage1 store APIs, diagnostics, and POC18 extraction. It does not change the
outer wire envelope or pCID discipline.

## Source corpus

- `protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`
- `implementations/poc19-production-shape/docs/DESIGN.md`
- `docs/thought-experiments/TE-givul-poc18-chunk-storage-identity.md`
- `docs/thought-experiments/TE-nozal-dag-cbor-role-in-promisegrid.md`
- `protocols/wire-lab.d/TODO/TODO-nahop-poc18-cas-git-replacement.md`
- `implementations/poc18-cas-git-replacement/store/store.go`

## Locked inputs

- POC19 starts with hybrid staged extraction: fresh stage0, POC18 as CAS/VCS
  baseline, and reviewed behavior slices factored into production-shaped stage1
  roles. Source: `DI-nupag`.
- POC18 currently stores `kind == "chunk"` under `chunks/<cid>.bin` and other
  objects under `objects/<cid>.cbor`, but the CID profile used by `CIDForBytes`
  is CIDv1 `raw` over exact bytes for both. Therefore the `.cbor` suffix in
  POC18 is a local filename convention, not proof that the object CID uses
  `dag-cbor`.
- `TE-givul` recommends raw public chunk identity for deduplication, IPFS-like
  raw-block interop, and simple sparse retrieval, while keeping manifests and
  promise interpretation in CBOR/grid objects.
- `TE-nozal` recommends a hybrid: keep `grid([42(pCID), ...])` as the
  wire/app/kernel envelope, keep raw chunks and Markdown specs as raw CIDs, and
  move durable graph objects toward true DAG-CBOR where IPLD traversal and tools
  help.
- POC20 requires local CAS to be the source of truth; local indexes and
  projections must be rebuildable from CAS.

## Assumptions

- Stage0 should need only minimal CAS behavior: fetch exact bytes by CID, verify
  CID, read enough descriptor/root data to start stage1, and avoid knowing every
  durable graph profile.
- Stage1 may own richer storage profiles, DAG-CBOR validation, CAR import/export,
  diagnostics, and derived indexes.
- CAS stores remain sparse and peer-relative. No local store promises global
  completeness.
- Printable object names use CIDv1 base32 text. Wire messages carry binary CIDs.

## Alternatives

### Alt A: Preserve POC18 split layout as source-of-truth

POC19 would keep `chunks/<cid>.bin` for chunks and `objects/<cid>.cbor` for other
objects as the durable source-of-truth layout.

This is easy to extract from POC18 because the current `FileStore` already uses
that shape. It also makes raw chunks visually obvious in the filesystem.

The problem is semantic drift. A path under `objects/*.cbor` can still be CIDv1
`raw`; the suffix does not prove DAG-CBOR or even canonical CBOR unless the store
validates it. The layout also makes object class look authoritative even though
the source of truth is the CID-addressed bytes plus CAS-resident manifests and
promises.

### Alt B: One CID-keyed CAS namespace with derived views

POC19 would store every retained object as exact bytes under one CID-keyed
namespace. Profiles such as raw chunk, DAG-CBOR graph object, exact `grid()`
message, Markdown spec, CAR artifact, executable byte object, and encrypted
ciphertext object are recorded or derived from CAS-resident manifests,
protocol-specific records, and rebuildable local indexes.

This keeps the source-of-truth simple: an object is exact bytes named by a CID.
Stage0 can fetch and verify bytes without understanding every profile. Stage1 can
build convenient views such as chunks, messages, DAG-CBOR objects, CAR receipts,
and app descriptors from CAS facts.

The cost is that human filesystem browsing is less self-documenting unless
diagnostics and derived indexes are good. The implementation also must avoid
treating derived indexes as authoritative state.

### Alt C: Source-of-truth subdirectories by multicodec/profile

POC19 would store raw objects, DAG-CBOR objects, CAR artifacts, executable
objects, messages, specs, and encrypted objects in separate profile directories.

This makes local browsing clearer and can reduce accidental parser mistakes.

The problem is duplication of source-of-truth semantics in filesystem paths.
Objects may move between profile classifications as specs improve, and an object
can have multiple roles: a CAR artifact is exact bytes, a diagnostic artifact,
and a transfer package; a message can be both exact `grid()` bytes and part of a
message DAG. A profile directory layout is useful as a derived projection, but
not as the canonical store.

### Alt D: Make all retained objects true DAG-CBOR

POC19 would require every retained object, including messages and chunk records,
to be true DAG-CBOR with CIDv1 `dag-cbor` CIDs, except perhaps raw file chunks.

This maximizes IPLD tooling for graph traversal, selectors, and CAR workflows.

The problem is overreach. Raw chunks should retain raw byte identity. Wire
envelopes currently use the PromiseGrid `grid()` tag and pCID-owned slots, which
are not automatically valid DAG-CBOR blocks. CWT/COSE, encrypted payloads,
Markdown specs, executable binaries, container layers, and exact transfer
artifacts may be better represented as raw or profile-specific byte objects.

## Scenario analysis

### Scenario 1: Stage0 bootstrap fetch

Alice installs the small `grid` stage0 binary. It reads a configured root CID,
fetches missing bytes from local CAS or a trusted peer fixture, verifies exact
CID bytes, and starts a stage1 descriptor.

- Alt A makes stage0 know about old `chunks/` and `objects/` layout too early.
- Alt B gives stage0 the simplest rule: fetch exact bytes by CID and verify.
- Alt C pushes profile-routing policy into stage0 unnecessarily.
- Alt D requires stage0 to carry DAG-CBOR awareness even when fetching raw
  executable bytes.

### Scenario 2: Large-file VCS storage

Bob stores a large binary file as Rabin chunks plus manifests. Identical raw
chunks from unrelated files should deduplicate naturally.

- Alt A works for raw chunks, but keeps the old layout as source-of-truth.
- Alt B works best: raw chunks are exact raw CID objects and manifests explain
  how chunks reconstruct files.
- Alt C works if raw profile directories remain derived, but is brittle if the
  directory is treated as canonical.
- Alt D only works if raw chunks are exempted, which means it is not actually
  all-DAG-CBOR.

### Scenario 3: Durable graph object traversal

Carol wants IPLD tools to inspect reference sets, manifests, snapshots, app
descriptors, root-adoption records, and mapping records.

- Alt A hides DAG-CBOR intent because `.cbor` filenames may still use raw CIDs.
- Alt B allows true DAG-CBOR objects where useful while keeping the store generic.
- Alt C can help diagnostics but should remain a rebuildable profile view.
- Alt D is attractive for these graph objects, but over-applies the rule to raw
  and protocol-specific byte objects.

### Scenario 4: CAR transfer between peers

Dave redeems a retrieval token and Ellen sends a CARv1 package over TCP. Dave
verifies CAR structure and exact contained CIDs.

- Alt A can store the CAR as an object but does not clarify whether the CAR or
  its contained blocks are authoritative.
- Alt B is precise: retain the CAR bytes as an artifact if useful, but verify and
  retain contained CAS objects by their own CIDs.
- Alt C can derive a `car-artifacts` view, but the contained objects still belong
  in the one CAS namespace.
- Alt D is not relevant because CAR is a container format, not the identity of
  each contained object.

### Scenario 5: Diagnostics and human review

Frank asks `grid diag` to show which local objects are chunks, messages,
DAG-CBOR graph objects, encrypted ciphertext, CAR artifacts, and specs.

- Alt A looks simple in the filesystem but misleads if `.cbor` means raw-CID CBOR
  bytes rather than DAG-CBOR.
- Alt B requires better diagnostics, but those diagnostics can be rebuilt from
  CAS facts and validation results.
- Alt C gives convenient derived views; this is useful as a cache, not as source
  of truth.
- Alt D makes some diagnostics easier but cannot explain raw and encrypted bytes
  without exceptions.

### Scenario 6: POC20 timeline replay

Grace rebuilds local projections from CAS after deleting local indexes. The same
event stream should reconstruct object roles, retention promises, roots, and app
run records.

- Alt A risks treating filesystem placement as state that must be preserved.
- Alt B aligns with POC20: CAS bytes plus CAS-resident facts are source of truth;
  indexes are rebuildable.
- Alt C is acceptable only if subdirectories are projections.
- Alt D helps for DAG-CBOR graph objects but does not solve projection authority.

## Conclusion

Alt B, one CID-keyed CAS namespace with derived views, is the surviving
alternative.

POC19 should lock this storage profile:

- the source of truth is exact object bytes named by CID;
- raw Rabin chunks remain CIDv1 `raw` objects naming exact chunk bytes;
- durable graph objects may use CIDv1 `dag-cbor` when they are true DAG-CBOR and
  benefit from IPLD traversal or tooling;
- exact `grid()` messages remain stored by their exact bytes and message CIDs;
- Markdown protocol specs may remain raw pCID objects until a later spec-bundle
  decision;
- CAR files are transfer/archive packages and review artifacts, not the authority
  over contained object identity;
- local profile catalogs, chunk views, message-DAG views, CAR receipt views, and
  filesystem browsing aids are rebuildable projections.

Alt A is rejected because it turns POC18's convenient split layout into a
misleading canonical model. Alt C is rejected as a source-of-truth layout but
retained as a useful derived-view strategy. Alt D is rejected because it
over-applies DAG-CBOR to raw chunks, executable bytes, encrypted bytes, exact
wire messages, and transfer artifacts.

## Implications for open TODOs and pending DIs

- `vumas.17` can be marked complete by `DI-hofaz`.
- `vumas.5` stage0 scaffold should implement only exact CID fetch/verify and not
  profile-specific source-of-truth directories.
- `vumas.6` daemon-managed CAS should use the one-namespace profile and derive
  any `chunks`, `messages`, `dag-cbor`, or `car` views from CAS.
- POC18 extraction must verify old `chunks/*.bin` and `objects/*.cbor` entries by
  CID before retaining them in the POC19 layout.
- `nahop.26` and `nahop.28` remain historical POC18 follow-up DFs; POC19 now has
  its own locked storage-profile decision for new durable store work.

## Decision status

Locked by `DI-hofaz` in
`protocols/wire-lab.d/TODO/TODO-vumas-poc19-production-shape.md`.
