#!/usr/bin/env python3
"""write_index.py - Apply build_index output by rewriting:
  - docs/thought-experiments/README.md  (replace existing index table + filename-convention prose)
  - protocols/<slug>.d/TODO/TODO.md     (rebuild index + brief preamble for each protocol)

This is one-shot, run from chunk D of the TE-39 (TE-mumuv) migration twig.
"""
import os, sys, re, importlib.util, subprocess

REPO = sys.argv[1] if len(sys.argv) > 1 else "/home/user/workspace/wire-lab"

# Load build_index module
spec = importlib.util.spec_from_file_location(
    "build_index", os.path.join(REPO, "tools/build-index/build_index.py"))
bi = importlib.util.module_from_spec(spec)
spec.loader.exec_module(bi)
bi.REPO = REPO
bi.MAP = os.path.join(REPO, "tools/migrate-handles/mapping.tsv")

# ---- Load mapping ----
rows = bi.load_mapping()
print(f"Loaded {len(rows)} mapping rows", file=sys.stderr)

# ---- Build TE README.md ----
te_table = bi.build_te_index(rows)

new_readme = f"""# Thought Experiments

Each thought experiment is a falsifiable mental run of a Wire Lab design choice. Each lives in its own file.

## Filename convention

```
TE-<proquint>-<slug>.md
```

The proquint handle (5 lowercase characters from the alphabet `bdfghjklmnprstvz` x `aiou` x `bdfghjklmnprstvz` x `aiou` x `bdfghjklmnprstvz`, per Wilkerson 2009) is the stable identifier of a TE. Handles are minted by `tools/mint-handle` from `time_ns -> SHA-256 -> first 2 bytes -> proquint-1`, with collision-retry against the directory glob. The handle is short, pronounceable, fork-stable, and assigned at file-creation time. The slug is a kebab-case rendering of the TE's title; it is informational and may be edited.

The proquint handle replaces both the integer alias (TE-1, TE-2, ...) and the timestamp slug (TE-YYYYMMDD-HHMMSS) used before the TE-39 (TE-mumuv) migration. The drafting timestamp is preserved in git history (the file's first commit) and surfaces in the index below as the **Mint date** column. Each migrated file carries a `## Prior aliases` section recording the integer + timestamp aliases it had before the migration.

## Index

{te_table}

The proquint handle is **both** the stable identifier and the display nickname. It is collision-free at mint time, fork-stable across branches (each fork mints its own handles; collisions at merge time are handled by re-minting), and short enough to use directly in prose ("per TE-titur S5"). DF / DI / DR descendant numbering still uses the handle root: DF-titur.1, DI-titur-..., DR-009 (DR has its own numbering scheme). Backward citations to integer aliases (e.g., "per TE-25 S5") remain valid; readers may consult the cited file's `## Prior aliases` section or the `Prior alias` column above to recover the integer.

Forward-pointers to a not-yet-drafted TE MUST use a thread-id (T-...) recorded in `OPEN-THREADS.md`. Naming an unminted *future* proquint is impossible (proquints are minted, not predicted), and naming a future *integer* is forbidden (the construction that produced the DT3 drift, locked closed by the 2026-05-07 Cat-3 Refinement on TE-titur).

## Editing policy

TE filenames are mostly immutable: once the proquint handle is minted, the file keeps that handle through any future title or content edits. The exception is the TE-39 corpus migration itself, which renamed 70 files in a single commit (`85766f0`); each migrated file's `## Prior aliases` section preserves the audit trail.

TE contents are edited under a categorized policy locked in [TE-dabol](TE-dabol-te-editing-policy-and-holistic-corpus.md) and refined by [TE-vudaf](TE-vudaf-editing-policy-tabletop.md). The locked DIs are `DI-020-20260502-213103` (categorized regimes), `DI-020-20260502-213104` (uniform applicability across all TE corpora), and `DI-020-20260502-213105` (holistic reading by default; single-TE reading only for obviously mechanical questions). The Cat-1 clause of `DI-020-20260502-213103` was superseded on 2026-05-02 by `DI-020-20260502-232651` (Cat-1a / Cat-1b split). Several Cat-3 navigational refinements appear in TE-dabol's `## Refinements` section. The canonical statement of the policy lives in `AGENTS.md` under "TE Editing Policy (Required)"; the seven categories in summary:

- **Cat-1a (current-pointer paths).** Mechanical sweep in place; no top-of-file note.
- **Cat-1b (historical-quotation paths).** Left untouched. Path references inside markdown blockquotes, attributed to another TE ("TE-N states ..."), in past tense, inside `## Refinements` sections, supersedence notes, or `Decision status` lines are Cat-1b. When in doubt, treat as Cat-1b.
- **Cat-2 (vocabulary updates).** Edit in place, with a top-of-file note pointing at the driving TE or TODO. The note must enumerate by ID every DI that lives in the affected TE, paired with an explicit promise that the rewrite preserves each DI's meaning. Mandatory pre-step: grep the corpus for the old term inside quotation contexts and classify each match Cat-2 (sweep) or Cat-2-historical (leave) before sweeping.
- **Cat-3 (navigational forward pointers).** Append a dated entry to the TE's `## Refinements` section (created if absent, placed after `## Decision status`). The TE body above is unchanged. No DI is filed.
- **Cat-4 (resolved-implication forward pointers).** Same shape as Cat-3, used when an Implications-and-future-work item resolves (a TODO filed; a DR opened; a downstream TE landed).
- **Cat-5 / Cat-6 / Cat-7 (substantive supersedence).** Not edits. Write a new TE that supersedes the old one and a new DI that supersedes the old DI. Update the older TE's `## Decision status` to `superseded by TE-<handle>` and its top-of-file `## Status` field to `superseded by TE-<handle> / DI-<id>`; otherwise leave the body untouched.

Every TE carries a top-of-file `## Status` field placed immediately after the TE ID line. Canonical values: `needs DF`, `decided`, `decided, refined`, `superseded by TE-<handle> / DI-<id>`, `withdrawn`. Legacy values preserved during the 2026-05-02 retrofit: `stub`, `open`, `recommended for immediate adoption`, `locked for the <protocol>`. New TEs prefer canonical values.

The corpus is read holistically by default: the TE corpus is one document with many facets, not a collection of independent essays. When any TE is in scope, the first move is to scan TE titles plus `## Status` fields and `Decision under test` sections across the corpus to find facets that share assumptions, vocabulary, or decisions. Single-TE reading is reserved for obviously mechanical questions (a single typo; a path that has demonstrably moved; a `## Status` field retrofit) and only after the holistic read has confirmed the question is mechanical.

Applicability is uniform across every TE corpus in this repository, whether the TE lives at the harness level (this directory) or inside a per-protocol directory (`protocols/<slug>.d/`). Per-protocol corpora may add stricter rules but may not relax these rules.

## Adding a new TE

1. Decide the title.
2. Render the title to kebab-case for the slug.
3. Mint a proquint handle: `cd tools/mint-handle && go run . -w 1`. The tool scans `docs/thought-experiments/` and `protocols/*/TODO/` for existing handles and retries until it finds a collision-free value.
4. Create `TE-<handle>-<slug>.md` in this directory. Include a top-of-file `## Status` field placed immediately after the TE ID line, with the appropriate initial value (`needs DF` for a TE in DF state; `decided` for a TE that locks DIs in the same commit). Use canonical values; reserve legacy values for the retrofit corpus.
5. While drafting, write any forward-pointers (to TEs that do not yet exist) using a thread-id from `OPEN-THREADS.md`. Naming an unminted future proquint is impossible (proquints are minted, not predicted); naming a future integer alias is forbidden (the DT3 drift class, locked closed by the 2026-05-07 Cat-3 Refinement on TE-titur).
6. Add a one-line summary to `../../protocols/wire-lab.d/specs/harness-spec-draft.md` Section 8 with a link.
7. Add the row to this index. The mint date is the date the file is committed (column auto-fills from `git log --format=%ad --date=short -- <path> | tail -1`); leave the `Prior alias` column blank for new TEs that never carried an integer alias.
8. Open a PR.

## Migration history

- **2026-05-07.** TE-mumuv (TE-39, formerly TE-39) locks proquint handles. Migration tool `tools/migrate-handles/` renamed 38 TEs and 32 TODOs in a single commit (`85766f0`); each renamed file gained a `## Prior aliases` section. Citation sweep `tools/sweep-citations/` rewrote 1,451 references across 94 files. The integer alias (TE-N) and the timestamp alias (TE-YYYYMMDD-HHMMSS) survive only as `Prior alias` entries in this index and in each file's `## Prior aliases` section.
- **2026-05-05.** TE-titur Cat-3 Refinement (chained on the same TE) split TE identity into a stable identifier (timestamp+slug filename) and a display nickname (integer TE-N). Superseded on 2026-05-07 by the proquint adoption above; the 2026-05-05 forward-pointer rule is reaffirmed and tightened in the 2026-05-07 Refinement.
- **2026-04-30.** TE-titur (formerly TE-25) locked the integer-alias display-nickname convention.
"""

readme_path = os.path.join(REPO, "docs/thought-experiments/README.md")
with open(readme_path, "w") as f:
    f.write(new_readme)
print(f"Wrote {readme_path}", file=sys.stderr)

# ---- Build per-protocol TODO.md files ----
# wire-lab.d/TODO/TODO.md is the MASTER CROSS-LIST showing all protocols.
# The other three are per-protocol queues only.

def build_master_crosslist(rows):
    """wire-lab.d/TODO/TODO.md content: master with sections per protocol."""
    out = []
    out.append("# TODO master cross-list")
    out.append("")
    out.append("Master cross-listed queue per TE-magup. Lists every TODO across the")
    out.append("wire-lab. The wire-lab harness is itself a protocol; its TODOs")
    out.append("appear under the **wire-lab (harness)** section below alongside")
    out.append("every other protocol's per-protocol queue.")
    out.append("")
    out.append("Per TE-mumuv (TE-39, locked 2026-05-07), each TODO is addressable by")
    out.append("its proquint handle (TODO-<handle>). The integer alias (TODO-N) and")
    out.append("the timestamp alias (TODO-YYYYMMDD-HHMMSS) survive as `Prior alias`")
    out.append("entries here and in each file's `## Prior aliases` section. New TODOs")
    out.append("minted after that date carry only the proquint handle.")
    out.append("")
    sections = [
        ("wire-lab (harness)", "wire-lab", "./"),
        ("group-session", "group-session", "../../group-session.d/TODO/"),
        ("ppx-dr", "ppx-dr", "../../ppx-dr.d/TODO/"),
        ("udp-binding", "udp-binding", "../../udp-binding.d/TODO/"),
    ]
    for label, slug, prefix in sections:
        out.append(f"## {label}")
        out.append("")
        out.append("| Handle | Mint date | Title | Prior alias |")
        out.append("|---|---|---|---|")
        proto_rows = [r for r in rows if r["kind"] == "TODO"
                      and f"protocols/{slug}.d/TODO" in r["new_path"]]
        proto_rows.sort(key=lambda r: int(re.match(r"TODO-(\d+)", r["int_alias"]).group(1))
                        if re.match(r"TODO-(\d+)", r["int_alias"]) else 99999)
        for r in proto_rows:
            handle = f"TODO-{r['handle']}"
            new_path = r["new_path"]
            mint = bi.git_first_date(new_path) or "(uncommitted)"
            title = bi.title_from_file(new_path) or "(no title)"
            ts = r.get("ts_alias", "")
            if ts and not ts.endswith("-mint"):
                prior = f"`{r['int_alias']}` / `{ts}`"
            else:
                prior = f"`{r['int_alias']}`"
            link = prefix + os.path.basename(new_path)
            out.append(f"| [{handle}]({link}) | {mint} | {title} | {prior} |")
        out.append("")
    return "\n".join(out)


def build_per_protocol(rows, slug, label):
    """Per-protocol TODO.md (group-session, ppx-dr, udp-binding)."""
    out = []
    out.append(f"# TODO queue: {slug}")
    out.append("")
    out.append(f"Per-protocol TODO queue (per TE-magup). Items in this file touch")
    out.append(f"only files under `protocols/{slug}.d/`. Anything broader is")
    out.append("harness-level and lives at `protocols/wire-lab.d/TODO/TODO.md`.")
    out.append("")
    out.append("Per TE-mumuv (TE-39, locked 2026-05-07), each TODO is addressable")
    out.append("by its proquint handle (TODO-<handle>). Prior integer / timestamp")
    out.append("aliases survive in the `Prior alias` column and in each file's")
    out.append("`## Prior aliases` section.")
    out.append("")
    out.append("## Index")
    out.append("")
    out.append(bi.build_todo_index(rows, slug))
    return "\n".join(out) + "\n"


# Master cross-list
master_path = os.path.join(REPO, "protocols/wire-lab.d/TODO/TODO.md")
with open(master_path, "w") as f:
    f.write(build_master_crosslist(rows))
print(f"Wrote {master_path}", file=sys.stderr)

# Per-protocol queues
for slug in ["group-session", "ppx-dr", "udp-binding"]:
    path = os.path.join(REPO, f"protocols/{slug}.d/TODO/TODO.md")
    if not os.path.exists(os.path.dirname(path)):
        continue
    with open(path, "w") as f:
        f.write(build_per_protocol(rows, slug, slug))
    print(f"Wrote {path}", file=sys.stderr)

print("Done.", file=sys.stderr)
