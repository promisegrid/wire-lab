# implementations/

B-side trees per [TE-liviv](../docs/thought-experiments/TE-liviv-spec-vs-implementation-split.md).

This top-level directory holds reference implementations of PromiseGrid
protocols. Each implementation lives in its own subdirectory:

```
implementations/<impl-name>/
```

The `<impl-name>` slug is human-readable and chosen by the implementer.
Multiple implementations of the same protocol coexist as siblings; one
implementation may implement many protocols.

## CHANGELOG requirement

Every implementation directory MUST contain a `CHANGELOG.md` at its
root recording **conformance claims** against upstream spec doc-CIDs.
Each entry has the shape:

```changelog-entry
claim:           implements | partially-implements | extends | deprecates
spec:            bafkrei...spec-doc-cid
scope:           full | partial-section-N | etc.
breaking-change: true | false
notes:           prose
```

See [TE-zukug](../docs/thought-experiments/TE-zukug-spec-doc-inversion-and-conformance-changelog.md) for the inversion thesis (conformance reference goes implementation -> spec, never spec -> implementation) and [TE-liviv](../docs/thought-experiments/TE-liviv-spec-vs-implementation-split.md) for the spec-side vs implementation-side split.

## External implementations

This directory is just *our local* collection. Third parties' implementations live in their own external repos with the same `CHANGELOG.md` convention; the shape is location-independent. A registry of known external implementations is out of scope for v0 (see TE-liviv OQ-32.1).

## Currently empty

No implementations exist yet. TODO-jodon (UDP-binding v0 reference) and TODO-bihon (ns-3 harness fixture) will populate this directory when their work begins.
