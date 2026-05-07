# DR-mihip - ppx/main merge preview blockers

DR-ID: DR-mihip
Date: 2026-05-07 21:22:13
Asked by: stevegt@t7a.org (Steve Traugott)
State: open
Question: Can Perplexity Computer make `origin/ppx/main` merge-ready by resolving the validation and protocol-consistency blockers found in Codex's merge preview before Codex merges `ppx/main` onto `main`?
Why this blocks progress: Codex previewed merging `origin/ppx/main` onto `main` and found the merge is mechanically likely to apply cleanly, but the branch still carries validation failures and contradictory protocol text. Merging now would move those contradictions into canonical `main` and make later cleanup harder to distinguish from integration fallout.
Affects:
- `origin/ppx/main` at preview tip `ccc33f2` (`Merge ppx/te-20260507-220000-disposition-pointer-ut-matrix: pointer from dropped-thread-disposition to UT verification matrix`).
- `origin/main` at preview tip `c646b06` (`Remove obsolete channel carrier artifacts`).
- `docs/research/nested-vs-stacked-envelopes-20260504.md`.
- `protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`.
- `protocols/wire-lab.d/TODO/pre149-audit-report-20260505.md`.
- `tools/spec/check.go`.
- `tools/spec/manifest.go`.
- `tools/mint-handle/corpus.go`.
- `protocols/group-session.d/specs/group-session-draft.md`.
- `protocols/group-session.d/TODO/TODO-bisur-group-transport-envelope.md`.
- `transports/wire-lab-devs-draft/README.md`.
- `tools/spec/freeze.go`.
Unblocks: clean Codex merge of `origin/ppx/main` onto `main`; follow-on review of ppx's protocolized TODO/TE/transport structure without known validation noise.
Waiting on: stevegt+ppx@t7a.org (stevegt-via-perplexity)
Decision: Pending. Perplexity Computer should either fix the blockers on `ppx/main` or reply with a DR/DI-backed explanation for any item it believes should not be fixed before merge.
Linked DI: none yet; expected follow-up DI(s) on `ppx/main` if any protocol semantics are intentionally changed.
Related commits:
- `c646b06` on `origin/main` — removes obsolete channel-carrier artifacts that previously caused merge conflicts.
- `ccc33f2` on `origin/ppx/main` — preview tip when these blockers were found.
Last updated: 2026-05-07 21:22:13 UTC

## Merge preview result

Codex ran a dry merge preview from the merge base of `origin/main` and
`origin/ppx/main`. The result had no content-conflict markers, no
`changed in both` conflicts, and no `added in both` conflicts. The reported
`removed in remote` entries were expected ppx-side deletions or migrations of
files that `main` had not changed.

The mechanical merge therefore looks likely to apply cleanly, but it should
not be treated as merge-ready until the validation and protocol blockers below
are resolved.

## Requested fixes before merge

1. **Make validation clean.**
   - `git diff --check origin/main..origin/ppx/main` reported trailing whitespace in:
     - `docs/research/nested-vs-stacked-envelopes-20260504.md`
     - `protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`
     - `protocols/wire-lab.d/TODO/pre149-audit-report-20260505.md`
   - `gofmt -l` on an exported `origin/ppx/main` tree reported:
     - `tools/spec/check.go`
     - `tools/spec/manifest.go`
     - `tools/mint-handle/corpus.go`

2. **Reconcile `Message-ID` semantics.**
   - `protocols/group-session.d/specs/group-session-draft.md` says `Message-ID` is dropped / absent.
   - `transports/wire-lab-devs-draft/README.md` says `Message-ID` is retained.
   - `protocols/group-session.d/TODO/TODO-bisur-group-transport-envelope.md` still has DI/TODO text saying `Message-ID` remains a human-oriented convenience field.
   - ppx should make one semantics true across the spec, TODO/DI trail, and transport docs, using a superseding DI if the meaning changes.

3. **Align freeze tooling with documented freeze paths.**
   - Group-session docs describe freezing group-session/group-transport specs.
   - `tools/spec/manifest.go` hard-codes `protocols/wire-lab.d/specs`, so the tool does not appear to operate on the documented group-session paths.
   - ppx should either generalize the tool or correct the docs so the documented command and actual tool scope match.

4. **Align transport directory naming.**
   - The outer transport/spec direction appears to use `transports/<pcid>--<slug>/`.
   - `transports/wire-lab-devs-draft/README.md` describes bootstrapping as `transports/wire-lab-devs-draft/` and later `transports/wire-lab-devs-<pcid>/`.
   - ppx should choose one naming rule and cite the controlling DR/DI/TE.

5. **Fix stale links and path references.**
   - `protocols/group-session.d/specs/group-session-draft.md` contains stale references such as `specs/MANIFEST.md`, `specs/group-transport-draft.md`, `specs/transport-spec-draft.md`, and relative links that appear to point at nonexistent sibling `docs/` or `DR/` directories from the spec file's current location.
   - ppx should update current-pointer paths while preserving historical quotation paths according to the TE editing policy.

6. **Handle `git add` errors in `tools/spec/freeze.go`.**
   - `tools/spec/freeze.go` currently ignores the return value from `exec.Command(...).Run()` in `gitStage`.
   - This violates the repo's error-handling policy for Go code. ppx should return and handle that error explicitly.

7. **Provide a clean test note.**
   - Codex could not complete exported-tree `go test ./...` validation because dependencies needed network fetches in the restricted environment.
   - ppx should rerun the relevant tests in its environment, report the exact commands, and include any required Go/toolchain assumptions.

## Acceptance criteria

- `git diff --check origin/main..origin/ppx/main` is clean.
- `gofmt -l` is clean for Go files added or changed on `ppx/main`.
- Relevant Go tests pass, or failures are documented as unrelated with exact output.
- `Message-ID`, freeze-tool scope, and transport directory naming each have one consistent, DI-backed story.
- Stale current-pointer paths are corrected without rewriting historical evidence.
- Codex can rerun the merge preview and see only expected delete/migration entries, with no validation blockers.
