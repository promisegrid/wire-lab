# TODO-fonuz - Proquint IDs for new coordination artifacts

Adopt a single proquint handle namespace for new TODO, TE, DR, and DI
artifacts on `main`, while preserving existing timestamp and integer IDs as
historical records.

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
Decision: New TODO, TE, DR, and DI artifacts on `main` use a single global proquint-1 handle namespace. TODO filenames use `TODO/TODO-<handle>-<slug>.md`; TE filenames use `docs/thought-experiments/TE-<handle>-<slug>.md`; DR filenames use `DR/DR-<handle>-<slug>.md`; DI entries use `ID: DI-<handle>`. Existing timestamp and integer IDs remain unchanged as historical records.
Intent: Remove the cross-fork collision risk created by timestamp and integer IDs while keeping current main stable enough to merge with `ppx/main`. A global handle namespace makes cross-references unambiguous across artifact kinds and lets `tools/mint-handle` serve as the local collision guard.
Constraints: The change applies only to newly-created artifacts. Dates remain metadata, not ID components. `tools/mint-handle` must scan current root paths and ppx's protocolized TODO paths so future merges do not create obvious handle collisions. Existing DR, TODO, TE, and DI references are not rewritten by this task.
Affects: `AGENTS.md`; `AGENTS-codex.md`; `AGENTS-ppx.md`; `TODO/TODO.md`; `TODO/TODO-fonuz-proquint-id-upgrade.md`; `tools/mint-handle/`.
