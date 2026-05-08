# DR-suhod - Partial fix for DR-mihip ppx/main merge-preview blockers

DR-ID: DR-suhod
Date: 2026-05-08
Asked by: stevegt+ppx@t7a.org (stevegt-via-perplexity)
State: open
Question: Are the partial fixes the bot has applied on `ppx/dr-mihip-merge-blockers-fix` (mechanical items 1, 2, 5-mechanical-only, 6, 7 of DR-mihip) acceptable for merge into `ppx/main`, and does Steve want to lock the four still-blocked semantic items (DR-mihip 2, 3, 4, and the rest of 5) via Decision-First framing now or defer them?
Why this blocks progress: DR-mihip names seven items that block a clean Codex merge of `origin/ppx/main` onto `main`. The bot has fixed the items it could fix mechanically without inventing decisions, but the remaining four are semantic and require Steve's locked DI. Without those four, Codex's merge-preview blockers do not all clear, even though `git diff --check` and `gofmt -l` are now clean.
Affects:
- `DR/DR-mihip-ppx-main-merge-preview.md` (the DR this responds to; lives on `origin/main`, will land on `ppx/main` when the next main→ppx/main sync runs).
- `docs/research/nested-vs-stacked-envelopes-20260504.md` (whitespace stripped).
- `protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md` (whitespace stripped).
- `protocols/wire-lab.d/TODO/pre149-audit-report-20260505.md` (whitespace stripped).
- `tools/spec/check.go` (gofmt'd).
- `tools/spec/manifest.go` (gofmt'd).
- `tools/mint-handle/corpus.go` (gofmt'd).
- `tools/spec/freeze.go` (`gitStage` now returns error, caller warns).
- `protocols/group-session.d/specs/group-session-draft.md` (TE/DR markdown links repathed `../docs` → `../../../docs`, `../DR` → `../../../DR`).
Unblocks: clean Codex merge preview on the validation/whitespace/gofmt axes; closes DR-mihip items 1, 2, 6, 7 outright; closes the unambiguous mechanical subset of item 5.
Waiting on: stevegt@t7a.org (Steve Traugott)
Decision: Pending — see "Open semantic items" below for the four DF intakes Steve must answer before DR-mihip can be fully closed.
Linked DI: none yet; DI(s) expected once Steve answers the DF intakes for the semantic items.
Related commits: to be appended on merge of `ppx/dr-mihip-merge-blockers-fix` into `ppx/main`.
Last updated: 2026-05-08

## Cross-link

This DR is the bot's authored response to **DR-mihip-ppx-main-merge-preview**
(`DR/DR-mihip-ppx-main-merge-preview.md`, asked by Steve on 2026-05-07,
State: open). DR-mihip enumerates seven blockers; this DR records which
ones the bot fixed mechanically and which ones remain blocked on Steve's
DF answers. DR-suhod and DR-mihip are intended to be read together: when
DR-suhod's open semantic items are answered and applied, DR-mihip can move
to `State: decided` and the merge preview re-runs clean.

Per AGENTS-ppx review style (DR-003, DI-001-20260428-195702), the bot does
not open a GitHub PR; this DR is the review-and-converge ask for the
`ppx/dr-mihip-merge-blockers-fix` branch.

## Items closed by this branch (mechanical)

### DR-mihip item 1 — Make validation clean

- `git diff --check origin/main..HEAD` is clean: trailing whitespace
  stripped from
  `docs/research/nested-vs-stacked-envelopes-20260504.md`,
  `protocols/wire-lab.d/TODO/TODO-lilar-session-replay-cleanup.md`,
  and `protocols/wire-lab.d/TODO/pre149-audit-report-20260505.md`.
  The first two stripped lines were markdown hard-line-break "  " sequences
  that exist nowhere else in the repo's markdown, so they were not load-
  bearing convention; removing them aligns with the rest of the corpus.
- `gofmt -l tools/spec/check.go tools/spec/manifest.go tools/mint-handle/corpus.go`
  is clean. The diffs were comment-list reflow (godoc indent rule) and
  struct-field column alignment. No semantic change.

### DR-mihip item 6 — Fix stale links and path references (mechanical subset)

In `protocols/group-session.d/specs/group-session-draft.md`, the markdown
links to `../docs/thought-experiments/...` and `../DR/...` were broken
because the spec file lives at depth `protocols/group-session.d/specs/`
and the relative paths only walked up one level. Repathed the four
markdown links (TE-hogus, TE-zalut, TE-junil, DR-009) to use
`../../../docs/...` / `../../../DR/...` so they resolve against the
repo root.

The remaining textual references in that file (`specs/MANIFEST.md`,
`specs/group-transport-draft.md`, `specs/transport-spec-draft.md`)
are NOT mechanical — they tie to the freeze-tool scope question
(DR-mihip item 3) and to the not-yet-existing `group-transport-draft.md`
file. They are left as-is and called out as a DF item below.

### DR-mihip item 7 — Handle `git add` errors in `tools/spec/freeze.go`

`gitStage` previously returned no value and discarded `Run()`'s error.
Per AGENTS.md Error Handling Policy ("In Go code, never ignore errors
with `_ = ...`; handle, propagate, or report errors explicitly"), this
violated the policy.

Change applied:

- `gitStage(repoRoot, path string) error` now returns the underlying
  `exec.Command(...).Run()` error.
- The two call sites in `cmdFreeze` warn to stderr with the affected path
  and the error string. The freeze ritual still completes, because the
  snapshot file and manifest are already written to the working tree;
  staging is a convenience whose failure is recoverable by the user
  running `git add` manually.
- The function's godoc comment now records the policy citation and the
  rationale, replacing the old "errors are ignored" wording without
  losing the convenience-not-correctness invariant. Comment Preservation
  Protocol satisfied.

### DR-mihip item 8 — Test note (record `go test ./...` result)

Run on the four Go modules in this repo (tools/mint-handle,
tools/migrate-handles, tools/spec, tools/sweep-citations) on
`ppx/dr-mihip-merge-blockers-fix` after applying items 1, 2, 6, 7:

    go version
    go version go1.24.4 linux/amd64

    (cd tools/mint-handle     && go test ./...)
    ok  github.com/promisegrid/wire-lab/tools/mint-handle	0.015s

    (cd tools/migrate-handles && go test ./...)
    ok  github.com/promisegrid/wire-lab/tools/migrate-handles	0.002s

    (cd tools/spec            && go test ./...)
    ok  github.com/promisegrid/wire-lab/tools/spec	0.003s

    (cd tools/sweep-citations && go test ./...)
    ok  github.com/promisegrid/wire-lab/tools/sweep-citations	0.004s

All four modules pass. The `go.mod` files declared `go 1.25` so the
toolchain auto-fetched go1.25.0 to run the build; the host install is
go1.24.4. No tests were skipped, no `-tags` flags were used. Network
access was available in this sandbox, so module downloads succeeded
(this contrasts with the restricted environment Codex reported in
DR-mihip item 7).

## Items still blocked on Steve's DF answers (semantic)

These four items cannot be closed without locking a DI. Per AGENTS.md
Decision-First protocol, the bot must not invent function names, header
semantics, freeze paths, or directory naming conventions on its own. DF
intakes follow the multiple-choice form AGENTS-ppx requires; each lists
the alternatives the bot can see today and explicitly asks Steve to pick
one or to add a fifth.

### DF intake A — DR-mihip item 2 — `Message-ID` semantics

Three repo locations describe `Message-ID` differently:

- `protocols/group-session.d/specs/group-session-draft.md` says
  `Message-ID` is dropped (the message CID is the identifier; a
  separate human identifier creates two competing identities).
- `transports/wire-lab-devs-draft/README.md` says `Message-ID` is
  retained.
- `protocols/group-session.d/TODO/TODO-bisur-group-transport-envelope.md`
  says `Message-ID` remains a human-oriented convenience field.

Question: which semantic should be canonical across the spec, the TODO/DI
trail, and transport docs?

- **A1.** `Message-ID` is dropped from the group-transport envelope. The
  CID-as-filename rule is the sole identifier. Update the transport
  README and TODO/DI trail to match the spec; record a superseding DI
  citing the "two-identities" reasoning already written into the spec.
- **A2.** `Message-ID` is retained as a human-readable convenience that
  is explicitly *not* the message identifier. Spec is the document that
  must change; record a superseding DI explaining why a non-identity
  human label is worth carrying, and update the spec's "dropped"
  language to "carried but not identity-bearing".
- **A3.** `Message-ID` is retained for transports that aren't
  group-transport (e.g., wire-lab-devs as currently written) but
  dropped specifically inside group-transport. Codify the
  per-transport divergence in a new DI; update both spec and transport
  README to cite it.
- **A4.** Defer: park all three locations behind a single open DR until
  group-transport has its own freeze, on the theory that the freeze
  itself is the right moment to lock identity semantics.

The bot's read-only sense is that A1 matches the strongest existing
argument in the spec, but the bot will not lock this without Steve's DI.

### DF intake B — DR-mihip item 3 — `tools/spec/manifest.go` scope

`tools/spec/manifest.go` hard-codes `SpecsDir = "protocols/wire-lab.d/specs"`
and `ManifestPath = "protocols/wire-lab.d/specs/MANIFEST.md"`. The
group-session spec text describes freezing group-session/group-transport
specs (which would live under `protocols/group-session.d/specs/`), so
the documented freeze command does not match the tool's actual scope.

Question: how should the tool and the documentation be reconciled?

- **B1.** Generalize the tool: take the protocol-directory as a flag
  (`-protocol protocols/group-session.d`) and freeze into that
  directory's `specs/MANIFEST.md`. Each protocol gets its own manifest;
  the existing wire-lab.d manifest stays put. The CLI grows one flag.
- **B2.** Keep the tool wire-lab-only and amend the group-session spec
  to say group-session has its own (yet-to-be-built) freeze tool, OR
  that group-session is frozen via the wire-lab manifest as a
  cross-protocol entry. Document explicitly that the current tool does
  not freeze group-session.
- **B3.** Move to a single repo-wide manifest (e.g., `specs/MANIFEST.md`
  or `MANIFEST.md` at the repo root) and update both the tool and every
  spec doc to point at that single manifest. The flat layout DI (DI-011-
  20260429-184454) was within `wire-lab.d/specs/`; a repo-wide flat
  layout would supersede that with a new DI.
- **B4.** Defer until group-session freezes; until then, keep the tool
  wire-lab-only and add a TODO/DR pointer in group-session-draft.md
  that says "freezing this spec requires generalizing the tool first".

### DF intake C — DR-mihip item 4 — Transport directory naming

Two conventions appear in the repo:

- The outer transport/spec direction in TE-junil and the wire-lab.d
  transport-spec draft uses `transports/<pcid>--<slug>/` once the
  protocol's pCID is minted.
- `transports/wire-lab-devs-draft/README.md` describes bootstrapping
  as `transports/wire-lab-devs-draft/` while drafted, and as
  `transports/wire-lab-devs-<pcid>/` once frozen — a different shape
  (suffix-with-pcid, no `--`, slug-first).

Question: which naming rule controls, and what is the controlling
DR/DI/TE citation?

- **C1.** `transports/<pcid>--<slug>/` is canonical for frozen
  transports; the `wire-lab-devs-draft` directory must be renamed to
  match this rule once frozen, and the README updated to cite the
  controlling DI from TE-junil (or a new DI if none exists yet).
  Drafts retain their `<slug>-draft` form.
- **C2.** `transports/<slug>-<pcid>/` (slug-first, hyphen, no `--`)
  is canonical; update the outer spec and TE-junil to match, and write
  a superseding DI explaining the choice (slug-first reads better in
  shell autocomplete; `--` is reserved for double-keyed dirs which
  the repo doesn't otherwise use).
- **C3.** Both forms are allowed, with the rule "drafts use `<slug>-draft`,
  frozen use `<pcid>--<slug>` or `<slug>-<pcid>` per the per-transport
  README's choice"; the controlling DI says the per-transport author
  picks one and cites it in their README.
- **C4.** Defer: park the question behind an open DR until a second
  transport (beyond wire-lab-devs) freezes and forces the convention
  to be tested by a real example.

The bot's read-only sense is that C1 matches the outer transport-spec
language most directly, but again will not lock without Steve's DI.

### DF intake D — DR-mihip item 5, semantic remainder

After the four markdown TE/DR links were repathed, three textual
references in `protocols/group-session.d/specs/group-session-draft.md`
are still stale-looking, but each one ties to a different upstream
semantic question:

- `specs/MANIFEST.md` (line 3) — tied to DF intake B (manifest scope).
- `specs/transport-spec-draft.md` (lines 11, 310) — the file exists at
  `protocols/wire-lab.d/specs/transport-spec-draft.md`, but the spec
  cites it as `specs/transport-spec-draft.md` as if relative to the
  group-session.d root. Tied to whether group-session has its own
  `specs/` (which it currently does, with only this draft in it) or
  citations should reach into wire-lab.d.
- `specs/group-transport-draft.md` (lines 5, 315) — file does not exist
  anywhere. Either group-transport will have its own draft file (in
  which case create the placeholder) or group-transport's content is
  considered already merged into group-session-draft.md (in which case
  drop the reference).

Question: how should these three references resolve?

- **D1.** After locking DF B (manifest scope), rewrite all three to
  point to whichever manifest/spec layout B picks. `group-transport-draft.md`
  becomes either a real placeholder file or is removed from the prose.
- **D2.** Treat each reference as a path-as-citation rather than a
  Markdown link; do not change the spec's wording, but add a paragraph
  near §1 saying "all `specs/<name>.md` references in this file
  resolve against the wire-lab.d root unless otherwise stated", then
  fix the missing `group-transport-draft.md` either by creating the
  placeholder or by collapsing the references.
- **D3.** Defer: leave all three references as-is, on the theory that
  they will all be rewritten when group-session freezes and the
  manifest/freeze-path question is locked anyway.

This intake's choice is constrained by DF B's choice; if DF B picks B1
or B3, D1 follows; if B picks B2 or B4, D3 is the lower-risk choice.

## Acceptance criteria for closing DR-mihip

- DR-mihip items 1, 2 (validation+gofmt halves only), 6, 7, 8 are closed
  by this branch's merge into `ppx/main`.
- DR-mihip items 2 (Message-ID), 3 (manifest scope), 4 (directory
  naming), and the semantic remainder of 6 are closed by Steve's DF
  answers above and the follow-up branches that apply the resulting
  DIs.

When all four DF intakes have locked DIs and the corresponding edits
have merged to `ppx/main`, the bot will append `State: decided` to
both DR-mihip and this DR-suhod, and Codex's merge preview should
re-run with no validation or semantic blockers.

## 2026-05-08 — DF intake A answered

Steve answered DF intake A with the `Message-ID:` drop semantics plus a
reader-side legacy compatibility allowance, recorded as
DI-012-20260508-033513 in
`protocols/group-session.d/TODO/TODO-bisur-group-transport-envelope.md`.

Locked result:

- Canonical group-session v0 writers MUST NOT emit `Message-ID:`.
- The message CID remains the only protocol identifier.
- Readers MAY accept exactly one legacy `Message-ID:` header before
  `Date:` and MUST ignore its value semantically.
- If the legacy header is present, its bytes remain part of the message
  file and therefore part of CID verification.

This closes the semantic part of DR-mihip item 2 once the linked spec and
transport README edits land. DF intakes B, C, and D remain open.

Correction to the `Unblocks` field above: the earlier statement that
item 2 closed outright referred only to the validation/gofmt half of
that item. The `Message-ID:` semantic half remained open until
DI-012-20260508-033513.

## 2026-05-08 — TE-nijab DFs answered; B/C/D reframed

Steve answered the three TE-nijab DFs that cover the remaining DR-suhod
B/C/D cluster. The locked answers are recorded in
`protocols/wire-lab.d/TODO/TODO-turog-te-41-group-session-freeze-procedure.md`:

- DF-nijab.1 -> 1.A, DI-026-20260508-054722: `transports/` is a
  lower-layer network/feed simulation surface, not a namespace owned by
  one frozen higher-layer protocol.
- DF-nijab.2 -> 2.A, DI-026-20260508-054723:
  `transports/wire-lab-devs-draft/` remains historical pre-layered
  specimen data; future migration or supersession is additive.
- DF-nijab.3 -> 3.B, DI-026-20260508-054724: freeze-doc cleanup is
  parked behind TODO-pipus/TE-43 instead of being executed now.

Effect on this DR: B/C/D are reframed, but not closed. The `tools/spec`
manifest-scope question is now downstream of the spec-vs-transport-data
boundary; the transport directory naming question is now a layered
path/metadata question; and stale spec references remain parked until
TODO-pipus/TE-43 designs the first CAS/feed migration. TODO-turog step 5
must not be executed as written.
