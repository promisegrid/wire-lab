# TE-22: Spec-doc store layout and pCID machinery

*Thought experiment, part of the [PromiseGrid Wire Lab](../../harness-spec.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-20260429-175530

(First drafted 2026-04-29 17:55:30 UTC.)

## Decision under test

Given that TE-21 (Alt-E) framed a spec doc as a layered promise (doc-level prediction + per-peer adoption), and given Steve has locked the following in chat 2026-04-29:

- pCID format = **CIDv1** (multibase + multihash + codec; standard IPLD CID).
- Cross-references between specs cite **the frozen pCID**, not a slug, not a slug-with-version.
- Layout = **flat** (`specs/<slug>-draft.md` and `specs/<slug>-{cidv1}.md` side-by-side, single `specs/MANIFEST.md`).
- Draft cross-refs = **strict**: a draft cross-reference must cite a frozen pCID; drafts may not cite other drafts.
- Self-reference = **external only**: a frozen spec file does not contain its own pCID; the pCID is derived externally from the file's bytes and recorded in the manifest.

… what is the operational machinery for freezing, hashing, storing, citing, and replacing spec docs in this repo?

The four still-open Decision Framing questions are:

- **DF-22.1 Hash input.** What bytes go into the CIDv1 hash — raw file bytes as written, normalized bytes, or some abstract syntax?
- **DF-22.3 Freezing mechanic.** What does the act of "freezing" produce on disk — a renamed copy, a git tag, both, or something else?
- **DF-22.4 Manifest format.** Is `specs/MANIFEST.md` machine-parseable structured data, prose with a table, or both?
- **DF-22.5 Freezing trigger.** Who/what initiates a freeze — manual ritual, CI hook, or chat command?

This TE is the operational follow-on to TE-21. TE-21 said *what a spec doc is* (a layered promise). TE-22 says *how the repo handles such docs in practice* — how the doc-level promise is captured as a content-addressed artifact, how peers cite which artifact they adopted, and how the next version of a spec relates to the previous one on disk.

## Assumptions

- A `Promise` in this repo is an autonomous speech act — an assertion of state in the past, present, or future, often conditional. (Carried from TE-21.)
- A `pCID` is a CIDv1 hash of a spec document's bytes. Two parties claim to "speak protocol pCID X" when each implements the rules in the document whose CIDv1 is X. (Carried from harness-spec.md §1 and TE-21.)
- A spec doc is a layered promise (TE-21 Alt-E): the doc itself promises future interop conditional on its assumptions/open-questions/known-issues lists, and each peer separately promises to behave as the doc says.
- The Wire Lab has, today, exactly one spec doc (`harness-spec.md`). It will likely grow to ~3-10 sibling spec docs over the lifetime of this repo (frame format, trust ledger, currency, eval rules, capability tokens, etc.). It will not grow to hundreds. Operational machinery should suit "a handful of long-lived spec families," not "an RFC index of thousands."
- The repo runs in git on GitHub today, but the design is meant to survive migration off GitHub. Anything that depends on GitHub-only features (Releases, Actions, Pages, branch protection) is a hazard for that migration.
- pCID-as-port-number means peers MUST be able to compute the pCID of any given spec file with no out-of-band agreement other than "use CIDv1 with parameters P." That parameter set must itself be pinned somewhere in the repo.
- "Frozen" means immutable on disk: a spec file whose content matches its filename's pCID will never be edited again. Edits move the content into a new draft and produce a new frozen file with a new pCID; the old file remains as historical evidence.
- The user has locked layout (flat), draft cross-refs (strict), and self-reference (external only). Those three are inputs to this TE, not under test.
- The CIDv1 format itself is locked. The remaining choice within CIDv1 is the parameter set — multibase, multihash function, multihash length, codec. This TE treats parameter selection as a sub-question of DF-22.1 (hash input) since both bear on what `pCID(spec)` deterministically returns.

## Alternatives

This TE walks five DFs sequentially. Each DF has its own alternatives.

DF-22.2 (tooling language) was added in an amendment after the bot verified in chat 2026-04-29 that pure-bash CIDv1 computation is brittle. The original DF list (1, 3, 4, 5) presupposed bash scripts; that presupposition is now itself under DF. The bot first guessed Go was not available, then on a follow-up check found it has passwordless `sudo` and successfully installed Go 1.24 + `github.com/ipfs/go-cid` and produced the canonical CID for the test vector `"hello world\n"` — `bafkreifjjcie6lypi6ny7amxnfftagclbuxndqonfipmb64f2km2devei4` — matching the result from Python `py-multiformats-cid` on the same input. Both libraries are therefore viable; the choice is now about preference, canonicality, and ergonomics, not feasibility.

### DF-22.1: Hash input

What bytes feed into the CIDv1 multihash?

#### Alt-1.A: Raw file bytes, no normalization

The CIDv1 hashes the literal byte sequence of `specs/<slug>-{cidv1}.md` exactly as written to disk. Trailing newline, line endings, whitespace — all part of the hash.

- **Easier:** trivially reproducible (`sha256sum file.md` style; no preprocessing). Anyone with the file gets the same hash. Tooling-independent.
- **Harder:** trivially fragile. CRLF vs LF flips the hash. Editor-trailing-newline flips the hash. Two operators with different editors who agree on content can disagree on hash.

#### Alt-1.B: Normalized text bytes (LF-only, single trailing newline, UTF-8 NFC)

The CIDv1 hashes a normalized form of the file: LF line endings, exactly one trailing newline, UTF-8 in NFC normalization, no BOM. Operators can write the file however they like; the hashing tool normalizes before hashing.

- **Easier:** robust against editor and OS differences. CRLF/LF and trailing-whitespace flips no longer flip the hash.
- **Harder:** the on-disk bytes don't equal the hashed bytes. To verify a frozen spec, you must run normalization first. The normalization rules become part of the protocol; if anyone normalizes differently, hashes diverge silently.

#### Alt-1.C: Abstract syntax (parsed Markdown AST or canonical JSON)

The CIDv1 hashes a parsed, re-serialized representation of the doc — e.g., the doc parsed to a Markdown AST, then re-serialized in a canonical form. Whitespace, ordering of frontmatter keys, link reference style, etc. all collapse to one canonical shape before hashing.

- **Easier:** robust against Markdown-syntax variation that doesn't change semantics (e.g., `[x](y)` vs reference-style links). Two operators using different Markdown styles for the same content can produce the same hash.
- **Harder:** big tooling commitment. Requires a parser, a canonicalizer, and agreement on the AST grammar. The grammar itself becomes a versioned protocol artifact. Bugs in parser/canonicalizer change the hash. Fundamentally trades concrete-file reproducibility for semantic robustness.

#### Alt-1.D: Raw bytes + machine-checked formatter

The CIDv1 hashes raw file bytes (Alt-1.A), but the repo enforces a single formatter (e.g., `prettier --parser=markdown`, or `tools/format-spec.sh`) on `specs/*.md`. The formatter is run before freeze and before any commit that would freeze. The formatter's behavior is what makes the bytes deterministic, not the hash function.

- **Easier:** keeps the hash function trivial (raw bytes) while still defending against editor-style flapping. Formatter is a single concrete tool the repo can lock to a version.
- **Harder:** formatter version is now part of the protocol. Two operators with different formatter versions can produce different bytes from the same input. Mitigated by pinning formatter version in `tools/freeze-spec.sh` and refusing to freeze if the working tree differs from a freshly-formatted run.

### DF-22.2: Tooling language for freeze and audit

Which language do `tools/freeze-spec.sh`-or-equivalent and `tools/check-spec-format.sh`-or-equivalent get written in? CIDv1 computation requires multihash header bytes + multibase encoding; the choice of language affects how robust, portable, and bot-runnable the freeze and audit rituals are.

#### Alt-2.A: Pure bash + sha256sum + base32 + xxd

Freeze and audit scripts are POSIX shell. Multihash header bytes are produced via `printf '\x12\x20'`, concatenated with the raw sha256, then base32-encoded with the `b` multibase prefix and lowercase, padding stripped.

- **Easier:** zero non-default dependencies on most Linux systems. Same shell that does git operations does the freeze.
- **Harder:** byte-level scripting in shell is brittle. Padding rules, lowercase conversion, and the `0x01 0x55` codec wrap are easy to get subtly wrong without a reference implementation. Hard to unit-test. The bot self-tested this in chat 2026-04-29: the encoding works but is one bug away from producing a wrong-but-plausible-looking CID.

#### Alt-2.B: Go (using `github.com/ipfs/go-cid`)

Freeze and audit scripts are Go programs that import the canonical reference implementation of CIDv1. Source lives at `tools/freeze-spec/main.go` and `tools/check-spec-format/main.go`; binaries are built on demand via `go run` or `go build`. No binary blob in git.

- **Easier:** the canonical CIDv1 reference. No risk of multihash-encoding bugs; the library has been audited by the IPFS community for years. Familiar to anyone in the content-addressable-storage ecosystem. Steve's stated language of preference (chat 2026-04-29). The PromiseGrid `grid-poc` repo is already Go, so this keeps the wire-lab toolchain coherent with the rest of the ecosystem.
- **Harder:** Go must be present on every machine that runs the freeze ritual. The bot's sandbox does not ship Go by default but `sudo apt-get install -y golang` works (verified 2026-04-29). The sandbox is ephemeral and may reset between sessions, so the freeze ritual's bootstrap step must idempotently install Go if missing. CI on a non-GitHub host needs Go available on the runner.

#### Alt-2.C: Python with `py-multiformats-cid` (or equivalent)

Freeze and audit scripts are Python 3 programs that import `multiformats_cid` (or `py-cid`, or `multiformats`). Python 3 is preinstalled on essentially every modern dev environment.

- **Easier:** Python is available everywhere, including the bot's sandbox by default (no install step). Verified 2026-04-29: `pip install py-multiformats-cid` works and produces a CID that matches the Go reference exactly for the same byte input. The library handles the multihash + multibase + codec wrap correctly. Scripts are short, readable, and unit-testable.
- **Harder:** adds a Python dependency for a repo that otherwise has no Python footprint and is otherwise Go-shaped (the wider PromiseGrid project is Go). The chosen library must itself be pinned (with a hash, ideally) because pip dependencies drift. Library quality varies — `py-multiformats-cid` works but is less mainstream than the Go reference; if it bitrots, the repo has to migrate. Conflicts with Steve's stated language preference (Go).

#### Alt-2.D: Go primary, Python cross-check (belt and suspenders)

Freeze and audit are Go programs (Alt-2.B). On developer machines that happen to have Python with `py-multiformats-cid` installed, an optional `tools/cross-check-cid.py` script can be run as a sanity test that reproduces the same CID for a given file. The Python cross-check is not part of the protocol; it's a manually-invoked second opinion that lets a developer catch a hypothetical bug in the Go library by comparing against an independent implementation.

- **Easier:** combines Steve's preferred language (Go) on the critical path with a free independent-implementation check. Catches library bugs the way two independent witnesses catch testimony bugs.
- **Harder:** redundant for a mature library like `go-cid`. Adds a Python dependency that is rarely exercised, which means it bitrots between exercises. For the scale this repo will reach (~10 specs over its lifetime), the cross-check is unlikely to ever fire.

### DF-22.3: Freezing mechanic

What does "freeze a spec" mean concretely on disk and in git?

#### Alt-3.A: Snapshot file only

Freezing copies `specs/<slug>-draft.md` to `specs/<slug>-{cidv1}.md`, computes the CIDv1, fills in the filename, commits it. The draft file remains in place for further editing if a next version is in progress; otherwise the draft can stay equal to the latest frozen content as a convenience. No git tag, no GitHub Release.

- **Easier:** completely git-host-agnostic. Survives migration off GitHub trivially. The on-disk file IS the artifact; nothing extra to track.
- **Harder:** git's per-commit identity (commit SHA) is independent of the file's pCID. To find "which commit froze this pCID," you need a manifest entry or a commit-message convention. There's no built-in pointer.

#### Alt-3.B: Git tag only

Freezing creates a git tag like `pcid/<cidv1>` pointing at the commit where the spec content matches that pCID. The spec file in the repo at that tag is the frozen artifact; no rename, no separate file. The draft file is the spec file.

- **Easier:** uses a git-native concept (tags). One source of truth: the file's content at the tagged commit.
- **Harder:** depends on git tags being preserved across hosting moves (mostly fine — tags are part of git). But it conflates the "draft / frozen" distinction with the "head / tagged-commit" distinction, which is confusing for human readers browsing the repo. The on-disk filename gives no signal that this content is frozen at some pCID. Every reader has to consult `git tag` to learn what's frozen.

#### Alt-3.C: Snapshot file + git tag (belt and suspenders)

Both Alt-3.A and Alt-3.B. The frozen file appears at `specs/<slug>-{cidv1}.md`, AND a git tag `pcid/<cidv1>` points at the freezing commit.

- **Easier:** human readers see the on-disk frozen file (Alt-3.A advantage); machines can also resolve `pcid/<cidv1>` via git (Alt-3.B advantage). Either path leads to the same artifact.
- **Harder:** two sources of truth must agree. The freeze ritual must produce both atomically. Diverges if someone tags without snapshotting, or snapshots without tagging.

#### Alt-3.D: Snapshot file + manifest entry, no git tag

Freezing copies the file (Alt-3.A) AND appends an entry to `specs/MANIFEST.md` recording the pCID, slug, freeze date, and freezing commit SHA. No git tag. The manifest is the index of frozen artifacts; the on-disk filename carries the pCID for human readability.

- **Easier:** the manifest is in-repo, version-controlled, fully portable across git hosts. No reliance on host-side metadata. The manifest entry can record richer metadata than a tag (status, supersedes, supersedes-by, freeze rationale).
- **Harder:** the manifest must be kept in sync with the filesystem. If someone adds a frozen file but forgets the manifest entry, machines miss it. If someone adds a manifest entry for a non-existent file, machines find a dangling reference. CI check needed: "every frozen file in `specs/` has a manifest entry, and vice versa."

### DF-22.4: Manifest format

How is `specs/MANIFEST.md` shaped?

#### Alt-4.A: Markdown table only

Single table with columns: pCID | Slug | Status | Frozen on | Supersedes | Notes. Human readers parse it visually; machines parse it by reading the table rows with a Markdown parser or a regex.

- **Easier:** lowest-tech possible. Renders nicely on any git host. No parser needed beyond Markdown.
- **Harder:** brittle for machine parsing if rows have unusual content (multi-line cells, escaped pipes). Schema is informal — adding a column is a coordination event.

#### Alt-4.B: Front-matter YAML block + prose

YAML frontmatter at the top of `specs/MANIFEST.md` carries the structured data; the rest of the file is human-readable prose explaining the manifest's role. Machines parse the YAML; humans read the prose.

- **Easier:** YAML is a clean machine-parseable format. Most Markdown tooling already handles frontmatter. The prose section can carry conventions and examples.
- **Harder:** typically frontmatter is at the top and small. A growing list of frozen specs in frontmatter starts to push the prose far down the page. The convention "frontmatter is metadata about the document" is mildly violated when frontmatter IS the document.

#### Alt-4.C: Separate machine file + Markdown index

`specs/MANIFEST.md` is human-readable Markdown with the table. `specs/manifest.yaml` (or `.json`) is the machine-readable mirror with the same data. Generation rule: one is the source of truth (probably the YAML/JSON), the other is generated from it; CI verifies they agree.

- **Easier:** humans read clean Markdown; machines read clean structured data. Each format optimized for its consumer.
- **Harder:** two files that must stay in sync. Source-of-truth choice is itself a sub-decision. Migration off GitHub requires whatever generates the Markdown to keep working.

#### Alt-4.D: Single Markdown file with fenced YAML inside

`specs/MANIFEST.md` has explanatory prose and one or more fenced ```yaml code blocks containing the structured data. Machines parse the YAML blocks; humans read the prose alongside. No separate file.

- **Easier:** single file, single source of truth, both audiences served. Markdown rendering on any host shows both prose and the YAML (as code blocks). The YAML can grow arbitrarily without distorting the document's shape.
- **Harder:** the convention "the YAML inside a code block in MANIFEST.md is authoritative" must be locked and tooling-supported. A reader who copies the YAML out and edits it has not edited the manifest. Mitigated by `tools/freeze-spec.sh` always editing the file in place.

### DF-22.5: Freezing trigger

Who or what initiates a freeze?

#### Alt-5.A: Manual ritual via `tools/freeze-spec.sh`

A human (or bot acting as human) decides a draft is ready, runs `tools/freeze-spec.sh <slug>`, and the script computes the pCID, copies the draft to `specs/<slug>-{cidv1}.md`, appends to the manifest, and stages a commit. The decision to freeze is a deliberate human (or agent-as-human) act.

- **Easier:** zero infrastructure. Survives any host migration. Aligns with this repo's "small, deliberate, reviewable acts" posture. Matches the existing `tools/` directory pattern (if any) or establishes one cleanly.
- **Harder:** depends on humans/bots remembering to run it. A draft can drift far from any frozen version if no one freezes for a while. Mitigated by adding "freeze" as a checkable item in the spec change workflow (TE-14).

#### Alt-5.B: CI hook on merge

Whenever a PR that modifies `specs/<slug>-draft.md` merges to main, CI computes a CIDv1 of the draft and freezes it as `specs/<slug>-{cidv1}.md` if the draft's content has changed since the last frozen version. Manifest entry is appended automatically.

- **Easier:** no human forgetting. Every meaningful change produces a freeze. Manifest stays current automatically.
- **Harder:** requires CI infrastructure on the host (GitHub Actions today; what about post-GitHub?). Every commit-merging-to-main becomes a freeze even when the change is editorial (typo fix), bloating the manifest with near-duplicates. Distinguishing "this change is freeze-worthy" from "this is a typo" needs human judgment that CI doesn't have.

#### Alt-5.C: Chat-or-PR command (`/freeze <slug>`)

Freezing is triggered by an explicit command — a comment on a PR, a chat command, or a commit message keyword. The command instructs CI (or a bot) to run the freeze ritual. Humans control the trigger; the mechanic is automated.

- **Easier:** combines Alt-5.A's deliberation with Alt-5.B's automation. The freeze act is explicit and traceable (the command appears in PR history).
- **Harder:** depends on the command-handling infrastructure. Less host-portable than Alt-5.A.

#### Alt-5.D: Manual ritual + scheduled CI audit

Humans run `tools/freeze-spec.sh` (Alt-5.A). CI runs a periodic audit (e.g., on every push) that checks: do `specs/MANIFEST.md`, the on-disk frozen files, and any cross-references all agree? If not, CI fails. CI does not freeze; it only audits.

- **Easier:** keeps the freeze act deliberate and host-independent. Adds a check that catches drift, broken cross-refs, missing manifest entries, etc. The audit logic can run anywhere git can run.
- **Harder:** still depends on humans choosing when to freeze. Relies on CI for discipline rather than for action.

## Scenarios

Six scenarios, each played against the alternatives. The bookkeeping convention follows the TE-21 pattern: each scenario calls out which DF-alternative combinations handle it best.

### S1 (genesis): freezing the first spec

Today: `harness-spec.md` is the only spec doc. It exists at the repo root, not under `specs/`. We want to mint its first pCID.

Steps under the locked layout (flat, draft + frozen side-by-side):
1. Move (or copy) `harness-spec.md` to `specs/harness-spec-draft.md`.
2. Compute pCID per DF-22.1.
3. Freeze per DF-22.3 → produce `specs/harness-spec-{cidv1}.md`.
4. Update manifest per DF-22.4.
5. Update all in-repo references to `harness-spec.md` to point at the new path(s).

- **DF-22.1:** Alt-1.A (raw bytes) or Alt-1.D (raw bytes + formatter) both work; Alt-1.B requires committing to normalization rules; Alt-1.C is overkill for one file. Alt-1.D is the most robust without becoming a parsing project.
- **DF-22.3:** Alt-3.A (snapshot only) or Alt-3.D (snapshot + manifest) both work; the latter is needed to make the manifest meaningful.
- **DF-22.4:** Any of Alt-4.A through Alt-4.D works for one entry. Alt-4.D scales best.
- **DF-22.5:** Alt-5.A (manual ritual) is the only viable choice for genesis — there is no prior automation.

S1 verdict: D-leaning answers across the board. Alt-1.D, Alt-3.D, Alt-4.D, Alt-5.A.

### S2 (typo fix): editing a draft after the spec is frozen

Frozen `specs/harness-spec-{X}.md` exists. Steve notices a typo. He fixes it in `specs/harness-spec-draft.md`. The frozen file is untouched. The draft is now slightly different from the latest frozen content.

Question: does the typo fix become a new pCID, or wait until a more substantive change accumulates?

- **DF-22.1:** Independent of which alternative, the typo changes the bytes/canonical form, so a freeze right now would produce a new pCID.
- **DF-22.3:** Alt-3.A or Alt-3.D handles either policy: freeze immediately = new file appears; defer = draft drifts. Neither alternative forces the choice.
- **DF-22.4:** Alt-4.A or Alt-4.D — the manifest must record both the prior pCID's status (now superseded?) and the new one. A "frozen | superseded-by | draft-ahead" set of statuses helps.
- **DF-22.5:** Alt-5.A (manual) lets Steve choose. Alt-5.B (CI on merge) would freeze every typo fix, bloating the manifest. Alt-5.C lets Steve trigger explicitly. Alt-5.D combines manual with audit (catches "you've drifted N commits without freezing").

S2 verdict: Alt-5.A (manual) or Alt-5.D (manual + audit) preserve human judgment about what's freeze-worthy. Alt-4.D's structured YAML can carry richer status fields.

### S3 (cross-spec citation): one spec cites another

Wire Lab grows a second spec, say `specs/trust-ledger-draft.md`. The trust-ledger draft wants to cite `harness-spec` at a specific frozen pCID (the locked-strict rule: cite frozen pCID, not slug, not draft).

How does the citation appear in `trust-ledger-draft.md`? How does a reader resolve it back to the actual file?

Option: the citation is `[harness-spec pCID bafy...](harness-spec-bafy....md)` — a relative link to the sibling file.

- **DF-22.1:** The citation is just bytes inside `trust-ledger-draft.md`; whatever Alt was chosen for the hash input applies. If `trust-ledger-draft.md` is later frozen, the citation is part of its hashed content, which is exactly the desired property — frozen specs hash-include the pCIDs of every other spec they depend on, creating an explicit dependency graph.
- **DF-22.3:** Alt-3.A is sufficient; Alt-3.D adds a manifest cross-link that lets a tool walk dependencies.
- **DF-22.4:** Alt-4.D's YAML can carry a `depends_on: [bafy...]` field per spec, making the dependency graph machine-readable.
- **DF-22.5:** Independent of trigger choice.

S3 verdict: Alt-4.D's structured YAML is doing real work here — it lets dependency graphs be checked mechanically, which becomes load-bearing as the number of specs grows.

### S4 (replacing a frozen spec): publishing the next pCID

Steve eventually freezes a major rev: `specs/harness-spec-draft.md` evolves substantially, and a new freeze produces `specs/harness-spec-{Y}.md`. The old `specs/harness-spec-{X}.md` remains on disk; it is not deleted, only marked superseded.

- **DF-22.1:** Alt-1.D minimizes drift between editor environments; the new pCID Y is the same regardless of who runs the freeze.
- **DF-22.3:** Alt-3.A leaves the old frozen file in place untouched (correct); Alt-3.D updates the manifest entry for X to set `superseded_by: Y` and adds a fresh entry for Y. The manifest is the only mutable record; the frozen files stay immutable.
- **DF-22.4:** Alt-4.D YAML cleanly records `supersedes: X` and `superseded_by: Y` fields; Alt-4.A's table can do it too but with weaker tooling support.
- **DF-22.5:** Alt-5.A or Alt-5.D — the freeze act is an explicit human event that aligns with publishing a major rev.

S4 verdict: the manifest's append-and-update pattern (Alt-4.D + Alt-3.D) handles supersession without ever touching the immutable frozen files. The on-disk timeline reads like an append-only log.

### S5 (host migration): the repo moves off GitHub

Hypothetical: tomorrow the repo migrates off GitHub to bare git on a self-hosted server. What machinery survives?

- **DF-22.1:** Alt-1.A, Alt-1.B, Alt-1.D all survive (pure file content + optional in-repo formatter). Alt-1.C survives if the AST parser is in-repo.
- **DF-22.3:** Alt-3.A and Alt-3.D survive cleanly (just files in git). Alt-3.B (git tags) survives because tags are part of git. Alt-3.C survives because both pieces are in git.
- **DF-22.4:** Alt-4.A, Alt-4.B, Alt-4.D survive (single file in git). Alt-4.C survives but has more pieces.
- **DF-22.5:** Alt-5.A survives trivially (a script in `tools/`). Alt-5.B does not survive without porting CI to whatever the new host offers. Alt-5.C does not survive without porting command infrastructure. Alt-5.D survives Alt-5.A; the audit half can be re-attached to whatever CI is available.

S5 verdict: the bot's "avoid GitHub lockin" rule is real. Alt-5.A (or Alt-5.D, with the audit treated as optional in the new home) is the only choice that fully survives migration. Alt-1.D, Alt-3.D, Alt-4.D all survive.

### S6 (chain-of-custody): tracing why pCID Y is "the right" successor of pCID X

Question: when Steve freezes pCID Y, how does someone reading the repo six months later know that Y is genuinely an evolution of X — not a parallel competing spec?

- **DF-22.1:** No alternative carries lineage information; the hash function is content-addressed, not provenance-aware.
- **DF-22.3:** Alt-3.A leaves only the freezing commit's diff to tell the story; Alt-3.D adds a manifest line linking Y to X. Alt-3.C's git tag adds a third pointer.
- **DF-22.4:** Alt-4.D's YAML can carry a structured `supersedes: X` field, making the chain explicit. A future tool can render the full lineage chain by walking the manifest.
- **DF-22.5:** Independent of trigger choice.

S6 verdict: chain-of-custody is the manifest's job. Alt-4.D shines here because the structure is machine-walkable; Alt-4.A works but requires conventions.

## Conclusions

Across S1-S6, a coherent recommended set emerges:

- **DF-22.1:** **Alt-1.D — raw bytes + machine-checked formatter.** The formatter is locked into the freeze and audit tools. Raw-byte hashing keeps the verification trivial (any CIDv1-aware tool reproduces the hash from the file as-is); the formatter eliminates editor-style flapping. Rejected: Alt-1.A is too fragile across editors; Alt-1.B silently diverges if normalization rules differ between operators; Alt-1.C is a parsing project we don't need yet. CIDv1 parameter set: `multibase=base32`, `multihash=sha2-256`, `codec=raw` (the spec is a byte stream). These three parameters are pinned in `specs/MANIFEST.md`.

- **DF-22.2:** **Alt-2.B — Go using `github.com/ipfs/go-cid`.** The freeze and audit programs are Go (`tools/freeze-spec/main.go`, `tools/check-spec-format/main.go`). Go is Steve's stated language preference and matches the wider PromiseGrid Go ecosystem (`grid-poc`). The library is the canonical CIDv1 reference; correctness is inherited rather than re-derived. The bot's sandbox installs Go on demand via `sudo apt-get install -y golang`; the freeze script's bootstrap is idempotent so sandbox resets do not break the ritual. Rejected: Alt-2.A (bash) is brittle for byte-level multihash assembly; Alt-2.C (Python) conflicts with stated language preference and adds a non-Go dependency to a Go-shaped project; Alt-2.D (cross-check) is unlikely to fire at this scale and bitrots between exercises.

- **DF-22.3:** **Alt-3.D — snapshot file + manifest entry, no git tag.** Snapshot file is the human-readable pointer (filename encodes pCID); manifest entry is the machine-readable record. Skip git tags: they duplicate the manifest's job, complicate the freeze ritual, and add a host-side concept that's confusing for readers who expect tags to mean "release." Rejected: Alt-3.A loses chain-of-custody (S6); Alt-3.B hides the frozen content behind tags (poor for browse-the-repo readers); Alt-3.C is double-bookkeeping with no clear win.

- **DF-22.4:** **Alt-4.D — single Markdown file with fenced YAML inside.** Prose at top explains what the manifest is; one fenced ```yaml block carries the structured per-spec entries (pCID, slug, status, frozen-on, supersedes, superseded-by, depends-on, freezing-commit, notes). Tooling reads the YAML block; humans read the whole file. Rejected: Alt-4.A is too brittle for machine parsing as the manifest grows; Alt-4.B distorts the document shape; Alt-4.C requires two-file sync.

- **DF-22.5:** **Alt-5.D — manual ritual + scheduled CI audit.** Humans (and bots-acting-as-humans) run `go run ./tools/freeze-spec <slug>` (or a built binary) when they decide a draft is freeze-worthy. CI's only role is to fail loudly if the manifest, the on-disk frozen files, and the cross-references disagree. Rejected: Alt-5.A alone is fine but skips the safety net; Alt-5.B over-freezes typo commits; Alt-5.C is host-bound. Alt-5.D matches the bot's "avoid GitHub lockin" rule because the audit is a Go program that can run anywhere with a Go toolchain; the trigger half (manual) doesn't need CI at all.

The full recommended set across the five DFs is **(1.d, 2.b, 3.d, 4.d, 5.d)** — four D answers and one B. Across the four original DFs, all-D wins; the addition of DF-22.2 (tooling language) introduces the one B because the canonical CIDv1 reference happens to live in the Go ecosystem (and Steve's stated language preference happens to be Go).

### Implications

- **Two-program pair in `tools/`.** `tools/freeze-spec/` is a Go module: it formats the draft, computes the CIDv1 via `github.com/ipfs/go-cid`, copies to `specs/<slug>-{cidv1}.md`, appends a manifest entry, stages a commit. `tools/check-spec-format/` is a sibling Go module used by the CI audit: it verifies that all draft files match the formatter's output, that every frozen file has a manifest entry and vice versa, and that every cross-reference cites a frozen pCID. Both programs pin the formatter version and the CIDv1 parameters as Go constants. A thin `tools/freeze-spec.sh` wrapper may exist purely as an ergonomic shim that runs `go run ./tools/freeze-spec "$@"` with the right working directory; it carries no logic of its own.

- **Three-state status field per spec entry.** Each manifest entry has `status: frozen | superseded | draft-ahead` (where "draft-ahead" means the draft has changed beyond this frozen version but no new freeze has been issued). This is the simplest status model that captures S2, S4, and S6.

- **Dependency graph is machine-readable.** Each manifest entry carries `depends_on: [<pCID>, ...]`. Walking the manifest produces the cross-spec dependency DAG. CI can check that every `depends_on` entry exists in the manifest and points at a `frozen` (not draft-ahead, not superseded) status — though `superseded` may be allowed with a warning, since older specs sometimes legitimately depend on older specs.

- **CIDv1 parameter set is pinned and visible.** The first lines of `specs/MANIFEST.md` (in prose, before the YAML block) state: "All pCIDs in this repo are CIDv1 with multibase=base32, multihash=sha2-256, codec=raw, computed over the formatted file bytes." This is the smallest piece of out-of-band agreement peers need.

- **External-only self-reference is honored automatically.** Frozen files do not contain their own pCID; the filename is the only place the pCID appears, and the filename is metadata. Drafts may contain a comment or front-matter line `pcid: TBD` to make their not-yet-frozen status visible, but this never conflicts with self-reference because the eventual frozen file has `pcid: TBD` replaced with… nothing. The freeze ritual strips the placeholder before hashing if present, OR (preferred) the placeholder is in the draft only, and the frozen snapshot is produced by removing it; either approach honors external-only self-reference.

- **Strict-draft cross-refs are honored by lint.** A draft file's cross-reference to another spec MUST cite a frozen pCID. CI audit script greps for cross-spec citations in `specs/*-draft.md` and verifies each cited pCID exists in the manifest as a frozen (or superseded) entry. A draft citing another draft fails the audit. This forces the dependency DAG to grow only by reference to frozen artifacts.

- **`harness-spec.md` migrates.** Per S1, it moves to `specs/harness-spec-draft.md`. The genesis freeze produces `specs/harness-spec-{cidv1}.md`. All in-repo references to `harness-spec.md` are updated. This is a separate TODO (TODO 011 below); it is NOT included in this TE's commit, because the migration is a discrete, reviewable act of its own.

- **Git tags remain available for human convenience.** The decision to skip git tags as a freeze artifact does not preclude humans from tagging meaningful commits for their own reasons (e.g., `v1-published` on the genesis-freeze commit). Such tags are not part of the protocol; they're just bookmarks.

- **Migration off GitHub is fully preserved.** Every machinery piece (freeze script, format script, audit script, manifest YAML, snapshot files) lives in the repo and works under bare git. The audit step that runs in CI today can run as a git pre-receive hook on a new host with no protocol change.

## Decision Framing questions

DF-22.1: Hash input.

- (a) Alt-1.A — raw file bytes, no normalization.
- (b) Alt-1.B — normalized text bytes (LF, NFC, single trailing newline).
- (c) Alt-1.C — abstract syntax (parsed Markdown AST canonical form).
- (d) Alt-1.D — raw bytes + machine-checked formatter (recommended).

DF-22.2: Tooling language.

- (a) Alt-2.A — pure bash + sha256sum + base32 + xxd.
- (b) Alt-2.B — Go (canonical `github.com/ipfs/go-cid`) (recommended; matches Steve's stated language preference and the wider PromiseGrid Go ecosystem).
- (c) Alt-2.C — Python with `py-multiformats-cid`.
- (d) Alt-2.D — Go primary, Python cross-check (belt and suspenders).

DF-22.3: Freezing mechanic.

- (a) Alt-3.A — snapshot file only.
- (b) Alt-3.B — git tag only.
- (c) Alt-3.C — snapshot file + git tag.
- (d) Alt-3.D — snapshot file + manifest entry, no git tag (recommended).

DF-22.4: Manifest format.

- (a) Alt-4.A — Markdown table only.
- (b) Alt-4.B — YAML frontmatter + prose.
- (c) Alt-4.C — separate machine file + Markdown index.
- (d) Alt-4.D — single Markdown file with fenced YAML inside (recommended).

DF-22.5: Freezing trigger.

- (a) Alt-5.A — manual ritual via `tools/freeze-spec.sh`.
- (b) Alt-5.B — CI hook on merge.
- (c) Alt-5.C — chat-or-PR command (`/freeze <slug>`).
- (d) Alt-5.D — manual ritual + scheduled CI audit (recommended).

The recommended set is **all-D except 2.b: (1.d, 2.b, 3.d, 4.d, 5.d)**. The four D answers survive migration off GitHub, keep the freeze act deliberate and auditable, and make the manifest a machine-walkable structure rather than just a human convenience. The B answer for tooling language is Go using the canonical `github.com/ipfs/go-cid` library: it matches Steve's stated language preference, aligns with the wider PromiseGrid Go ecosystem (`grid-poc`), and uses the audited reference implementation rather than a less-mainstream port. The bot's sandbox can `sudo apt-get install -y golang` on demand, so bot-side execution is not a blocker; the freeze ritual will include the install step idempotently.

Already-locked decisions, not under DF in this TE (carried from chat 2026-04-29):

- **pCID format = CIDv1.** Multibase + multihash + codec wrap.
- **Layout = flat.** `specs/<slug>-draft.md` and `specs/<slug>-{cidv1}.md` side-by-side; single `specs/MANIFEST.md` at the top of `specs/`.
- **Draft cross-refs = strict.** A draft citing another spec MUST cite a frozen pCID, never another draft.
- **Self-reference = external only.** A frozen spec file does not contain its own pCID; the pCID is in the filename and in the manifest, not in the document body.

## Decision status

`decided` — all five DFs locked by Steve in chat 2026-04-29. Locked answers are recorded as DI entries in `TODO/011-te-spec-doc-store-and-pcid-machinery.md`:

- **DF-22.1: 1.a** — raw file bytes, no normalization. (Steve preferred transparency over editor-style robustness; the formatter from the recommended Alt-1.D was dropped.)
- **DF-22.2: 2.b** — Go using `github.com/ipfs/go-cid`.
- **DF-22.3: 3.d** — snapshot file + manifest entry, no git tag.
- **DF-22.4: 4.d** — single Markdown file with fenced YAML inside.
- **DF-22.5: same-binary variant** — freezer and checker are subcommands of one Go binary `tools/spec`. The trigger half of Alt-5.D (manual freeze + scheduled CI audit) is preserved; the structural half (two separate programs) is collapsed to one.

Amendment history:

- 2026-04-29 (initial): TE drafted with four DFs (hash input, freezing mechanic, manifest format, freezing trigger). Recommended set: all-D.
- 2026-04-29 (amendment): DF-22.2 (tooling language) added after the bot verified that pure-bash CIDv1 is brittle. The original DFs assumed `tools/freeze-spec.sh` was a bash script; that assumption itself was elevated to a DF. The bot initially mis-reported Go as unavailable, then corrected itself after testing that passwordless `sudo` lets it install Go 1.24 and that `github.com/ipfs/go-cid` produces the canonical CID for the test vector `"hello world\n"` matching `py-multiformats-cid`. Revised recommended set: (1.d, 2.b, 3.d, 4.d, 5.d).
- 2026-04-29 (DF lock): Steve locked (1.a, 2.b, 3.d, 4.d, same-binary variant of 5). Two deviations from the bot's recommendation: (i) Alt-1.A over Alt-1.D — Steve preferred raw-byte transparency to the formatter's editor-style robustness; (ii) DF-22.5 same-binary variant — one Go binary with `freeze` and `check` subcommands instead of two separate programs.

## Implications for follow-on work

- **TODO 011 (now in progress):** DIs locked for all five DFs and the four already-decided inputs. Remaining work: implement `tools/spec/` (single Go binary with `freeze`, `check`, `cid`, `ls` subcommands), perform the genesis freeze of `harness-spec.md`, add `specs/MANIFEST.md`, wire `go run ./tools/spec check` into CI.

- **TODO 012 (provisional):** Add the CI audit step: `go run ./tools/spec check` performs format checks (advisory CRLF/BOM warnings on drafts), manifest-vs-disk consistency, cross-ref-citation lint, and self-reference lint. The audit runs on every push to ppx/main and main; failures block the merge. Because the audit is a Go program, it runs identically under GitHub Actions, a self-hosted runner, or a git pre-receive hook on a non-GitHub host.

- **TODO 010 (existing):** Drives TE-21 to DI. TE-21 + TE-22 together form the spec-doc-as-promise bundle: TE-21 says what a spec doc *is*; TE-22 says how the repo handles such docs. The DI entries from both TEs should land in the same revision of `harness-spec.md`'s vocabulary section.

- **Future TE (planned):** Peer-level adoption metadata. TE-21 Alt-E said each peer's adoption is a separate promise that can name which answers it chose for open questions. TE-22 makes the doc-side machinery concrete; the peer side is still open. A follow-on TE will work out the wire-level shape of an adoption promise (likely a small structured payload referencing `pcid` + `open_question_choices: {Q7: yes, Q9: variant-B}`).
