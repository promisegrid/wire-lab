# tools/

Corpus-maintenance tools. Read this file before running any of them.
You (the LLM driving the session) are the operator -- there is no
human in the loop pressing keys for these tools, so the runbooks below
exist for your benefit. Hazard notes are not optional.

## Inventory

| dir | purpose | run frequency |
| --- | --- | --- |
| `build-index/` | rebuild `docs/thought-experiments/README.md` and per-protocol `TODO.md` index tables from `migrate-handles/mapping.tsv` + git log mint dates | after any TE or TODO add/rename |
| `matrix-runner/` | generate manifests, run, checkpoint, validate, view, and compare root scenario/simulation results; preferred over legacy Python result tools | every matrix canary/full run |
| `migrate-handles/` | one-shot TE-39 migration: rename legacy `(TE|TODO)-<timestamp>-<slug>.md` to proquint form and inject `## Prior aliases` sections | done; mapping.tsv is now read-only authority |
| `mint-handle/` | allocate a new unique proquint handle for a new TE or TODO; collision-checks against the existing corpus | every new TE/TODO file |
| `spec/` | freeze draft spec docs into content-addressed snapshots, audit manifest, list snapshots; backs the harness-spec workflow | when freezing a new spec pCID |
| `sweep-citations/` | rewrite integer-alias and timestamp-alias references in body text to proquint form (Cat-2 mechanical sweep) | after any new TE/TODO add or rename, or whenever stale citations are discovered |

## sweep-citations runbook (READ THIS BEFORE RUNNING)

`sweep-citations` is a regex-based mechanical Cat-2 tool. It cannot
distinguish a *current pointer* (Cat-2 -- sweep) from a *historical
quotation* (Cat-1b -- leave) by surrounding text alone. The
classification is YOUR job, not the tool's.

### Hazard: destructive without `-n`

Running `go run . -r ../..` WITHOUT `-n` writes to disk immediately.
There is no confirmation prompt. If you forget `-n` and the run hits
a Cat-1b false positive, you will silently re-damage protected files
on every run. This has happened. Recovery is `git checkout -- .` if
you catch it before commit; if you commit, revert the commit.

**Rule: always run with `-n` first. Drop `-n` only after comparing
the flagged-files list against `tools/sweep-citations/CAT-1B-HANDOFF.md`
and confirming every flagged file is either (a) a true Cat-2 target
or (b) a documented Cat-1b that you will hand-edit (or skip).**

### Procedure

1. From `tools/sweep-citations/`, run the tests:

   ```
   go vet ./...
   go test -count=1 ./...
   ```

   All must pass. If they don't, stop and investigate; do not sweep
   with a broken tool.

2. Dry-run the sweep from the tool dir:

   ```
   go run . -n -r ../.. 2>/tmp/sweep-stderr.log 1>/tmp/sweep-stdout.log
   ```

   The first thing on stderr is the Cat-1b handoff (verbatim copy of
   `CAT-1B-HANDOFF.md`). Read it. It lists the known historical-
   quotation citations the regex will re-flag forever.

3. Compare the dry-run's flagged-files list against the handoff:

   - **Every flagged file is in the handoff:** the corpus is in
     steady state. There is nothing new to sweep. Stop.
   - **A flagged file is NOT in the handoff:** that file has a fresh
     Cat-2 stale citation. Inspect manually:
     - If it is a current pointer, it is Cat-2; let the tool sweep it.
     - If it is a historical quotation, it is a new Cat-1b false
       positive; add an entry to `CAT-1B-HANDOFF.md` (path, expected
       edit count, line numbers, matched tokens, why) BEFORE running
       without `-n`. Then sweep with `-n` again to confirm only the
       Cat-2 files remain unflagged-by-handoff. (The tool will still
       flag the new Cat-1b file even after you document it -- the
       handoff is for you, not the tool.)
   - **A flagged file in the handoff has a different edit count than
     the handoff records:** something changed in that file. Inspect
     before sweeping.

4. To actually write changes, you have two paths:

   - **Path A (preferred):** hand-edit only the Cat-2 files (the ones
     not in the handoff) using the dry-run output as your worklist.
     This keeps the Cat-1b files completely untouched and avoids the
     "ran without `-n`" hazard entirely.
   - **Path B:** drop `-n` to let the tool write, then immediately
     `git diff` and `git checkout --` each Cat-1b file listed in the
     handoff to revert the false-positive damage. Verify with
     `git status` that only Cat-2 files remain modified before you
     commit.

5. Commit the sweep result on a `ppx/{twig}` branch with a subject
   that names the trigger (new TE, new TODO, citation drift, etc.).
   Do not bundle the sweep with unrelated changes.

### Cat-1b classification rule (DI-020-20260502-232651)

A citation is Cat-1b (leave) if any of these apply:

- It lives inside a Refinement section that intentionally contrasts
  old/new aliases for pedagogical reasons.
- It is a quoted commit message or other quoted historical artifact.
- It is a `Source:` provenance line citing the alias the file
  originally carried at authoring time.
- It is a speculative prediction whose proquint substitution would
  mislead future readers (e.g. "originally hedged as 'Likely TODO 020'
  before that slot was claimed").

Otherwise it is Cat-1a (current pointer in surviving prose) or Cat-2
(mechanical reference) -- sweep it.

### Currently documented Cat-1b false positives (as of 2026-05-07)

See `tools/sweep-citations/CAT-1B-HANDOFF.md` for full per-entry
detail. Summary:

1. `AGENTS.md` -- 1 edit (vocabulary-pedagogy contrast)
2. `AGENTS-ppx.md` -- 1 edit (indented literal commit-message body)
3. `docs/thought-experiments/TE-titur-...md` -- 6 edits
   (chunk-D Refinement explaining proquint identity rules)
4. `docs/thought-experiments/TE-vudaf-...md` -- 1 edit
   (historical-narrative paragraph about path-rename event)
5. `protocols/wire-lab.d/TODO/TODO-bihon-...md` -- 1 edit
   (`Source:` provenance line)
6. `protocols/udp-binding.d/TODO/TODO-jodon-...md` -- 3 edits
   (`Source:` provenance + parenthetical historical hedge)

## mint-handle runbook

For every new TE or TODO file:

```
cd tools/mint-handle
go run . -r ../..
```

Prints a fresh proquint that does not collide with any existing
handle in the corpus. Width is locked at proquint-1 (5 chars) per
TE-mumuv (TE-39). Use `-w 2` only if Steve explicitly authorizes a
width bump.

## build-index runbook

After any TE or TODO add, rename, or mapping.tsv edit:

```
python3 tools/build-index/build_index.py
```

Regenerates `docs/thought-experiments/README.md` and the per-protocol
`TODO.md` index tables from `tools/migrate-handles/mapping.tsv` plus
git log first-commit dates. The README and `TODO.md` master indices
are sweep-skipped (sweep-citations excludes them) precisely because
this tool owns them.

## matrix-runner runbook

`matrix-runner` is the preferred root result runner. It replaces the legacy
Python result scripts for normal operation but does not delete them. Result
views are generated from `results/`; scenario-side matrices are not committed.
Source: `DI-lulom`; `DI-zamin`.

Before a full run:

```
cd tools/matrix-runner
GOCACHE=/tmp/wire-lab-gocache go test ./...
```

Generate and run a canary first:

```
go run . manifest -repo-root ../.. -models openai-gpt-5.3-codex-xhigh -run-group-id canary-<ts> -timestamp <ts> -shuffle-seed 42 -limit-cells 3
go run . run -repo-root ../.. -manifest results/manifests/matrix-manifest-canary-<ts>.csv -provider openai -api-model <openai-api-model> -reasoning-effort xhigh -result-style concise -max-output-tokens 6000 -max-run-cost-usd <budget> -max-cell-estimate-usd <cell-cap>
go run . validate -repo-root ../.. -manifest results/manifests/matrix-manifest-canary-<ts>.csv
go run . view -repo-root ../.. -model openai-gpt-5.3-codex-xhigh
```

Only start the full manifest after canary validation passes and the state file's
token/cost totals are acceptable. Source: `DI-nugiv`.

## spec runbook

See `protocols/wire-lab.d/specs/` and the harness-spec workflow.
The `spec` binary's behaviour is locked by TODO 011 DIs
(DI-011-20260429-184457 through -184501).

## migrate-handles

**Done. Do not run again.** `tools/migrate-handles/mapping.tsv` is the
read-only authority for integer/timestamp -> proquint mappings. If a
new mapping is needed (e.g. for a previously un-migrated alias), edit
mapping.tsv directly with a clear commit message; do not re-run the
one-shot tool.
