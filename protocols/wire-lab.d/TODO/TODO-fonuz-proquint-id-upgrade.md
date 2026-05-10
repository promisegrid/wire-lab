# TODO-fonuz - Proquint IDs for new coordination artifacts

Adopt a single proquint handle namespace for new TODO, TE, DR, and DI
artifacts, while preserving existing timestamp and integer IDs as historical
records.

## Subtasks

- [x] fonuz.1 Lock the new-artifact ID decision in this TODO's Decision Intent Log.
- [x] fonuz.2 Import and extend `tools/mint-handle` for TODO, TE, DR, and DI handles.
- [x] fonuz.3 Update agent instructions so new artifacts use proquint IDs.
- [x] fonuz.4 Validate the mint tool, formatting, and stale-template search.

## Decision Intent Log

ID: DI-nisam
Date: 2026-05-07 21:09:47
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: New TODO, TE, DR, and DI artifacts use a single global proquint-1 handle namespace. TODO filenames use `protocols/<slug>.d/TODO/TODO-<handle>-<slug>.md`; TE filenames use `docs/thought-experiments/TE-<handle>-<slug>.md`; DR filenames use `DR/DR-<handle>-<slug>.md`; DI entries use `ID: DI-<handle>`. Existing timestamp and integer IDs remain unchanged as historical records.
Intent: Remove the cross-fork collision risk created by timestamp and integer IDs while keeping current main stable enough to merge with `ppx/main`. A global handle namespace makes cross-references unambiguous across artifact kinds and lets `tools/mint-handle` serve as the local collision guard.
Constraints: The change applies only to newly-created artifacts. Dates remain metadata, not ID components. `tools/mint-handle` must scan historical root paths and ppx's protocolized TODO paths so future merges do not create obvious handle collisions. Existing DR, TODO, TE, and DI references are not rewritten by this task.
Affects: `AGENTS.md`; `AGENTS-codex.md`; `AGENTS-ppx.md`; `protocols/wire-lab.d/TODO/TODO.md`; `protocols/wire-lab.d/TODO/TODO-fonuz-proquint-id-upgrade.md`; `tools/mint-handle/`.

ID: DI-sazud
Date: 2026-05-10 05:03:32
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: `tools/mint-handle` scans every file and directory name in the repo working tree, excluding only `.git` internals, and treats any proquint-shaped substring in those names as an occupied handle candidate. The scanner also scans readable regular files across the working tree for exact `ID: DI-<handle>` owner lines so DI handles remain in the same occupied set. Repeated occurrences of the same substring are allowed because ordinary repo names such as `manifest` can reserve the same proquint-shaped string in multiple places; the minting tool's job is to avoid occupied strings, not to certify every occurrence as a unique coordination owner.
Intent: Remove layout-specific blind spots from the local collision guard. A handle should not be minted if the same proquint-looking string already appears in any working-tree file or directory name, regardless of whether that path is currently a TODO, TE, DR, DI, protocol, simulation, tool, or ordinary support file.
Constraints: This decision is an occupied-set rule for minting, not a corpus-validity rule. The name scan is recursive and path-agnostic except for `.git`; content scanning remains limited to exact DI owner lines because arbitrary prose references to proquint-looking strings are not ownership claims. The decision updates and narrows the scanner implementation obligations from `DI-nisam`; it does not rewrite historical IDs.
Affects: `tools/mint-handle/corpus.go`; `tools/mint-handle/main.go`; `tools/mint-handle/main_test.go`; `protocols/wire-lab.d/TODO/TODO-fonuz-proquint-id-upgrade.md`.
