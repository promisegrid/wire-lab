# TODO 011 - TE-22 Spec-doc store layout and pCID machinery: drive to DI

Track the work to drive TE-22 (`docs/thought-experiments/TE-20260429-175530-spec-doc-store-and-pcid-machinery.md`) from `needs DF` to a set of decided DIs that lock the operational machinery for spec-doc storage, pCID computation, freezing, and cross-referencing in the wire-lab.

## Already-locked inputs

These were locked by Steve in chat 2026-04-29 before TE-22 was written. They are inputs to TE-22, not under DF in TE-22 itself:

- pCID format = CIDv1 (multibase + multihash + codec wrap).
- Layout = flat (`specs/<slug>-draft.md` and `specs/<slug>-{cidv1}.md` side-by-side; single `specs/MANIFEST.md`).
- Draft cross-refs = strict (a draft citing another spec MUST cite a frozen pCID, never another draft).
- Self-reference = external only (a frozen spec file does not contain its own pCID).

These four become DI entries in this file when work begins; they are decided and only need formal recording.

## Subtasks

- [x] 011.1 Steve answers DF-22.1 (hash input). LOCKED: 1.a (raw file bytes, no normalization). Chat 2026-04-29.
- [x] 011.1b Steve answers DF-22.2 (tooling language). LOCKED: 2.b (Go with `github.com/ipfs/go-cid`). Chat 2026-04-29.
- [x] 011.2 Steve answers DF-22.3 (freezing mechanic). LOCKED: 3.d (snapshot file + manifest entry, no git tag). Chat 2026-04-29.
- [x] 011.3 Steve answers DF-22.4 (manifest format). LOCKED: 4.d (single Markdown file with fenced YAML inside). Chat 2026-04-29.
- [x] 011.4 Steve answers DF-22.5 (freezing trigger and tool shape). LOCKED: same-binary variant — freeze and check are subcommands of one Go binary `tools/spec/spec`, not two separate programs. Manual freeze + scheduled CI audit (Alt-5.D's trigger half) is preserved; only the binary count is collapsed from two to one. Chat 2026-04-29.
- [x] 011.5 DIs written for all five DFs and the four already-locked inputs (CIDv1, flat layout, strict draft cross-refs, external-only self-reference). See Decision Intent Log below.
- [x] 011.6 Pin the CIDv1 parameter set explicitly in `specs/MANIFEST.md`: `multibase=base32`, `multihash=sha2-256`, `codec=raw`. This is provisional; if a future TE locks different parameters, capture as a new DI.
- [x] 011.7 Genesis freeze of `specs/harness-spec-draft.md`: move to `specs/harness-spec-draft.md`, mint pCID via `go run ./tools/spec freeze harness-spec`, snapshot to `specs/harness-spec-{cidv1}.md`, append manifest entry, update all in-repo references. Separate twig. (Genesis pCID: `bafkreigtaivld55rekcswfj26mo26e267m3ytzgflqb2qcclyiicpfzc6i`.)
- [x] 011.8 Implement `tools/spec/main.go` (single Go binary with `freeze`, `check`, `cid`, `ls` subcommands). Pin the Go module versions in `go.sum`. The bot's freeze ritual's bootstrap step is idempotent: checks `go version`, runs `sudo apt-get install -y golang` if Go is missing, then proceeds. This matters for ephemeral bot sandboxes that may not have Go preinstalled. Uses `github.com/ipfs/go-cid` and `github.com/multiformats/go-multihash`.
- [x] 011.9 Implement the CI audit step (`go run ./tools/spec check`). Performs manifest-vs-disk consistency, cross-ref-citation lint, and any other invariants that emerge during 011.8 implementation. Wire into CI; the same Go binary runs identically as a git pre-receive hook on a non-GitHub host. (Wiring into CI is deferred until after the genesis freeze lands on `ppx/main`; the binary itself is implemented and the audit passes locally.)
- [ ] 011.10 After genesis freeze lands, open a follow-on TE on peer-level adoption metadata (the wire shape of "I, peer P, promise to behave as pCID X with open-question answers Q7=yes, Q9=variant-B"). This is the missing half of TE-21 Alt-E that TE-22 did not address.

## Decision Intent Log

ID: DI-011-20260429-184453
Date: 2026-04-29 18:44:53
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: pCID format is CIDv1 (the IPFS Content Identifier v1 format, which is `multibase || version || codec || multihash`). All pCIDs in this repo are CIDv1, with default parameter set `multibase=base32`, `multihash=sha2-256`, `codec=raw`.
Intent: Adopt the canonical content-addressing format used by the IPFS ecosystem so that pCIDs are interoperable with any CID-aware tool. Chosen over CIDv0 because v1 is the modern default and supports multibase prefixing (so a pCID can be parsed without out-of-band agreement on encoding).
Constraints: Every spec doc's pCID is exactly the CIDv1 of its hashed input bytes. The parameter set is pinned in `specs/MANIFEST.md` prose; changing it requires a new DI. Tools that produce or verify pCIDs MUST use a canonical CIDv1 implementation (verified compatible: `github.com/ipfs/go-cid` and `py-multiformats-cid` both produce `bafkreifjjcie6lypi6ny7amxnfftagclbuxndqonfipmb64f2km2devei4` for the test vector `"hello world\n"`).
Affects: `specs/harness-spec-draft.md` §1's pCID note; `specs/MANIFEST.md`; the freeze and audit binary; future references to `pCID` everywhere in this repo.
Linked DR: none (chat-directed lock); pre-cursor TE: TE-22.

ID: DI-011-20260429-184454
Date: 2026-04-29 18:44:54
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Spec-doc store layout is FLAT. Drafts and frozen snapshots live side-by-side under `specs/`, named `specs/<slug>-draft.md` (mutable) and `specs/<slug>-{cidv1}.md` (immutable, filename encodes the pCID). A single global registry lives at `specs/MANIFEST.md`. No nested per-slug directories.
Intent: Keep the directory structure trivially scannable for a small but growing number of sibling specs (~10 over the lifetime of this repo). A flat layout means `ls specs/` shows everything at a glance; nested directories would add a layer of navigation that doesn't pay for itself at this scale.
Constraints: All spec files — drafts and frozen snapshots — live directly under `specs/`. The `specs/MANIFEST.md` file is the single global registry of frozen pCIDs and their metadata. If the spec count grows past the point where flat is unwieldy (~50?), this DI is the one to revisit.
Affects: where `specs/harness-spec-draft.md` migrates to (`specs/harness-spec-draft.md` and eventually `specs/harness-spec-{cidv1}.md`); how new spec docs are added; how the freeze and audit binary discovers files.
Linked DR: none; pre-cursor TE: TE-22.

ID: DI-011-20260429-184455
Date: 2026-04-29 18:44:55
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Cross-references between spec docs are STRICT: a draft spec that cites another spec MUST cite a frozen pCID, never another draft. A draft citing another draft is a lint failure.
Intent: Force the dependency graph between specs to grow only by reference to immutable artifacts. If draft A could cite draft B, and B then changed, A's promise would silently re-aim at a new target. By requiring frozen-pCID citations, every dependency edge is content-addressed and stable.
Constraints: The `tools/spec check` audit greps for cross-spec citations in `specs/*-draft.md` and verifies each cited pCID exists in the manifest as a `frozen` (or `superseded`) entry. A draft citing another draft fails the audit and blocks merge. A consequence: when starting a new spec from scratch, you can't cite a sibling spec until that sibling is frozen at least once.
Affects: how spec drafts are written; how the freeze ritual orders multi-spec changes; the `tools/spec check` audit logic.
Linked DR: none; pre-cursor TE: TE-22.

ID: DI-011-20260429-184456
Date: 2026-04-29 18:44:56
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Self-reference is EXTERNAL ONLY. A frozen spec file does not contain its own pCID anywhere in its content; the pCID lives only in the filename and in the manifest entry. The pCID is derived externally from the file's bytes.
Intent: Avoid the circular hash problem. If a frozen file contained its own pCID, the pCID would have to be computed before the file existed in its final form, which is impossible without a fixed-point trick. By keeping self-reference external, the freeze ritual is mechanical: hash the bytes, record the result in the filename and manifest. The file's content is unaffected by its own identity.
Constraints: Spec drafts MUST NOT contain a placeholder for their own pCID-to-be (`pcid: TBD`, `<self-pcid>`, etc.). If a draft needs to reference "this spec," it does so by slug name only (`harness-spec`) or by descriptive prose ("this document"). The `tools/spec check` audit greps for known self-reference patterns and fails if any are present.
Affects: how spec drafts are written; how the freeze ritual works (no placeholder substitution step); the audit's lint rules.
Linked DR: none; pre-cursor TE: TE-22.

ID: DI-011-20260429-184457
Date: 2026-04-29 18:44:57
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: DF-22.1 (hash input) = Alt-1.A: raw file bytes, no normalization. The CIDv1 hashes the literal bytes of `specs/<slug>-{cidv1}.md` exactly as written to disk: trailing newline, line endings, whitespace, all included.
Intent: Maximally honest hashing. The pCID hashes exactly what's on disk — no preprocessing, no hidden normalization, no formatter dependency, no parsing project. "What you see is what got hashed." Anyone with `sha256sum` and a multihash/multibase wrapper can reproduce the pCID without any out-of-band agreement on normalization rules. Trades editor-style robustness for total transparency.
Constraints: Operators must take care that their editors do not silently flip line endings, add BOMs, or change trailing-newline behavior between freeze runs. CRLF/LF flips, BOM additions, and trailing-whitespace changes all change the pCID. The audit may surface warnings about CRLF or BOM in drafts as advisory hints, but the canonical hash is over raw bytes regardless. If editor-style flapping becomes a real operational pain, a future DI may add Alt-1.D (raw bytes + machine-checked formatter); for now we accept the fragility for the transparency.
Affects: `tools/spec freeze` (no preprocessing step before hashing); `tools/spec check` (advisory warnings about CRLF/BOM allowed but not blocking); operator habits around editor configuration.
Linked DR: none; pre-cursor TE: TE-22 DF-22.1.

ID: DI-011-20260429-184458
Date: 2026-04-29 18:44:58
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: DF-22.2 (tooling language) = Alt-2.B: Go using `github.com/ipfs/go-cid`. The freeze and audit logic is implemented in Go; the canonical CIDv1 reference library handles multihash + multibase + codec wrapping.
Intent: Use the canonical CIDv1 reference implementation, the wider PromiseGrid Go ecosystem's native language (`grid-poc` is Go), and Steve's stated language preference. Avoid pure-bash byte-level encoding (brittle, one bug from silently wrong CIDs) and Python (adds a non-Go dependency to a Go-shaped project). Verified 2026-04-29: bot's sandbox can install Go via `sudo apt-get install -y golang`, and `go-cid` produces canonical CIDs matching `py-multiformats-cid` for the test vector.
Constraints: All freeze and audit logic lives in Go. Module versions are pinned in `go.sum`. Bot's freeze ritual bootstrap is idempotent: checks `go version`, runs `sudo apt-get install -y golang` if missing, then proceeds. CI on any host needs Go available on the runner. Migration to a different language is a new DI.
Affects: `tools/spec/` (Go module); CI configuration; bot's per-session bootstrap; what languages the wire-lab repo carries.
Linked DR: none; pre-cursor TE: TE-22 DF-22.2.

ID: DI-011-20260429-184459
Date: 2026-04-29 18:44:59
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: DF-22.3 (freezing mechanic) = Alt-3.D: snapshot file + manifest entry, no git tag. Freezing produces a new file at `specs/<slug>-{cidv1}.md` AND appends an entry to `specs/MANIFEST.md`. No git tag is created.
Intent: Keep all freeze metadata in-repo and host-portable. The on-disk filename carries the pCID for human readability; the manifest entry carries machine-readable status, supersedes-by, depends-on, and freezing-commit fields. Skip git tags because they duplicate the manifest's job and add a host-side concept that is confusing for readers who expect tags to mean "release."
Constraints: Every freeze MUST produce both the snapshot file and the manifest entry, atomically (in a single commit). The audit verifies that every frozen file in `specs/` has a manifest entry and vice versa; a mismatch fails the audit. Git tags MAY exist on freezing commits as bookmarks for human convenience, but they are not part of the protocol.
Affects: `tools/spec freeze` (must update both file and manifest in one commit); `tools/spec check` (manifest-vs-disk consistency check); `specs/MANIFEST.md` (the authoritative index of frozen artifacts).
Linked DR: none; pre-cursor TE: TE-22 DF-22.3.

ID: DI-011-20260429-184500
Date: 2026-04-29 18:45:00
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: DF-22.4 (manifest format) = Alt-4.D: single Markdown file with fenced YAML inside. `specs/MANIFEST.md` has prose at the top explaining what the manifest is and which CIDv1 parameters are pinned, followed by one fenced ` ```yaml ` code block containing the structured per-spec entries (pCID, slug, status, frozen-on, supersedes, superseded-by, depends-on, freezing-commit, notes). Tooling reads the YAML block; humans read the whole file.
Intent: Single source of truth, both audiences served. Markdown rendering on any git host shows both prose and YAML; no two-file sync. Status field is one of `frozen | superseded | draft-ahead`.
Constraints: There is exactly ONE fenced YAML block in `specs/MANIFEST.md`, and it is authoritative. Editing the YAML elsewhere (a copy in another file, or a YAML block elsewhere in the same file) is not editing the manifest. The `tools/spec` binary always edits the file in place. The audit verifies that the YAML block parses, that every entry has the required fields, and that internal cross-references (`supersedes`, `superseded_by`, `depends_on`) resolve to entries in the same manifest.
Affects: `specs/MANIFEST.md` shape; `tools/spec freeze` (appends to the YAML block); `tools/spec check` (parses the YAML block and validates).
Linked DR: none; pre-cursor TE: TE-22 DF-22.4.

ID: DI-011-20260429-184501
Date: 2026-04-29 18:45:01
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: DF-22.5 (freezing trigger and tool shape) = same-binary variant. The freezer and the checker are the SAME Go binary, with subcommands `spec freeze <slug>` and `spec check`. The trigger half is preserved from Alt-5.D: humans (or bots-acting-as-humans) run `spec freeze <slug>` when they decide a draft is freeze-worthy; CI runs `spec check` periodically as an audit. CI never freezes anything itself.
Intent: Eliminate a class of "audit and freeze disagree" bugs by sharing all the load-bearing logic (CIDv1 computation, manifest parsing, cross-ref resolution, format conventions) between the two operations. One binary, one source of truth, one place to fix bugs. Steve's preference for one binary over two (chat 2026-04-29).
Constraints: There is exactly ONE Go module under `tools/spec/`. Subcommands include at minimum `spec freeze <slug>` and `spec check`. Additional utility subcommands (`spec cid <file>`, `spec ls`) MAY exist but are not part of the protocol. CI invokes `go run ./tools/spec check` (or a built binary) on every push to `ppx/main` and `main`; failures block the merge. The same binary runs identically under GitHub Actions, a self-hosted runner, or a git pre-receive hook on a non-GitHub host.
Affects: `tools/spec/main.go` (single Go binary); CI configuration; the freeze ritual; the audit ritual; the operator's mental model (one binary, two main subcommands).
Linked DR: none; pre-cursor TE: TE-22 DF-22.5.

## Notes

- TE-22 carries the full alternative analysis (Alt-1.A-D, Alt-3.A-D, Alt-4.A-D, Alt-5.A-D) and six scenarios (S1-S6). This file does not duplicate that analysis; it tracks the decision-driving work.
- The recommended set is `(1.d, 2.b, 3.d, 4.d, 5.d)` per the amended TE. The four D answers survive migration off GitHub, keep the freeze act deliberate and auditable, and make the manifest a machine-walkable structure rather than just a human convenience. The B answer for tooling language is Go using the canonical `github.com/ipfs/go-cid`: it matches Steve's stated language preference, aligns with the wider PromiseGrid Go ecosystem, and uses the audited reference implementation. Bot-side execution is not a blocker because the bot's sandbox has passwordless `sudo` and can install Go on demand (verified 2026-04-29). Both `go-cid` and `py-multiformats-cid` agree on the test vector `"hello world\n"` — `bafkreifjjcie6lypi6ny7amxnfftagclbuxndqonfipmb64f2km2devei4` — confirming the canonical encoding.
- TE-22 is the operational follow-on to TE-21. TE-21 said *what* a spec doc is (a layered promise); TE-22 says *how* the repo handles such docs (freeze, hash, store, cite, supersede). The DI entries from both TEs should land in the same revision of `specs/harness-spec-draft.md`'s vocabulary section.
- Linked DR: to be created as DR-011 once subtasks 011.1-011.4 begin landing. For now, the TE itself stands as the open-question record.
- Companion TODOs: TODO 006 (DI-provenance backfill) and TODO 007 (DR backfill for §11) become more concrete after TE-22 locks because the genesis freeze produces the first content-addressed reference points for those backfills.
- The "External-only self-reference" lock has a subtle implication: the freeze ritual must produce a frozen file whose bytes do not reference its own pCID. The simplest implementation is to ensure no pCID-of-self placeholder exists in `specs/<slug>-draft.md` at freeze time. Document this as an explicit step in `tools/freeze-spec/main.go`.
