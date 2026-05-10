# TODO-zofan: Move DI records to first-class files

## Prior aliases

This TODO is being filed after the TE-mumuv proquint migration locked
2026-05-07, so it is minted directly under `TODO-zofan`. No prior integer or
timestamp alias.

## Status

Open. This TODO must run before removing DI-specific content scanning from
`tools/mint-handle`, because DIs currently live as inline records inside TODO
files and therefore do not yet have their own handle-bearing pathnames.

## Context

`tools/mint-handle` should be a generic occupied-name scanner. It should not
need to understand TODO, TE, DR, or DI semantics to avoid minting an already
visible handle. That requires every coordination artifact type that owns a
handle to expose that handle in a file or directory name.

TODO, TE, and DR records already do this through `TODO-<handle>-...`,
`TE-<handle>-...`, and `DR-<handle>-...` filenames. DI records do not: current
wire-lab DIs are inline `ID: DI-<handle>` entries inside TODO files. That
special case forced `tools/mint-handle` to parse DI body text, which is the
wrong abstraction.

The CSWG coordination repo provides the intended shape: a top-level `DI/`
directory with one file per DI, named `DI/DI-<handle>-<slug>.md`, plus TODO-file
references back to those DI files.

## Target shape

- `DI/README.md` documents the DI directory contract.
- `DI/_template.md` gives the one-DI-per-file template.
- `DI/DI-<handle>-<slug>.md` stores one append-only DI record.
- TODO files keep references to relevant DI files, not full inline DI bodies.
- DR files may link to resolving DI files using the same relative-link style.
- After DI filenames exist, `tools/mint-handle` can remove DI-specific content
  scanning and rely on whole-working-tree file/directory-name scanning.

## Subtasks

- [ ] zofan.1 Document the `DI/` directory contract using `~/lab/cswg/coordination/DI/README.md` as the model.
- [ ] zofan.2 Inventory every inline `ID: DI-...` record across TODO files.
- [ ] zofan.3 Create one `DI/DI-<handle>-<slug>.md` file per DI, preserving the fields, author identity, status, intent, constraints, affected paths, supersedence, and event history.
- [ ] zofan.4 Replace inline TODO DI bodies with references to the new DI files after those files exist.
- [ ] zofan.5 Update `AGENTS.md`, glossary text, and relevant docs so DIs are first-class files rather than inline TODO sections.
- [ ] zofan.6 Supersede or amend transitional `DI-sazud` wording after DI pathnames exist, so the decision record no longer normalizes DI-specific scanner logic.
- [ ] zofan.7 Remove DI-specific content scanning from `tools/mint-handle`; keep only whole-working-tree file/directory-name scanning, excluding `.git`.
- [ ] zofan.8 Validate DI links, handle discovery, documentation references, Go tests, and `errcheck`.
