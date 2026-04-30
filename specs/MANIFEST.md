# Spec-doc manifest

This file is the authoritative index of frozen spec docs in the wire-lab.
Drafts live alongside frozen snapshots under `specs/` (flat layout, per
DI-011-20260429-184454). Each frozen spec doc is content-addressed by its pCID
(a CIDv1 of the file's literal bytes, per DI-011-20260429-184457).

CIDv1 parameter set (per DI-011-20260429-184453):

- multibase: `base32`
- multihash: `sha2-256`
- codec:     `raw`

The single fenced YAML block below is the machine-readable record. It is
authoritative; humans read it via this file's Markdown rendering, machines
parse it via the `tools/spec` Go binary. Status values are
`frozen` (the canonical name of one frozen pCID), `superseded`
(an older frozen pCID whose successor has been frozen), or `draft-ahead`
(reserved; not currently emitted by the freeze tool but recognized by the
audit if a future tool emits it).

```yaml
entries:
  - pcid: bafkreigtaivld55rekcswfj26mo26e267m3ytzgflqb2qcclyiicpfzc6i
    slug: harness-spec
    status: frozen
    frozen_on: "2026-04-30T03:50:13Z"
    freezing_commit: 33920e004bd533a07e7b043a3d291d4d76d86ffa
```

See `docs/thought-experiments/TE-20260429-175530-spec-doc-store-and-pcid-machinery.md`
for the full reasoning behind the layout, hash input, manifest format, and
the single-binary freeze-and-check tool.
