#!/usr/bin/env python3
"""build_index.py - Rebuild docs/thought-experiments/README.md and per-protocol
TODO.md index sections from mapping.tsv + git log first-commit dates.

One-shot, run from chunk D of the TE-39 (TE-mumuv) migration twig.

Reads:
  tools/migrate-handles/mapping.tsv  (71 rows: handle/kind/int_alias/ts_alias/new_path/old_path)
  Each TE/TODO file's first line ("# TE-handle: Title" / "# TODO-handle: Title")
  git log --format=%ad --date=short -- <new_path> | tail -1   for mint date
  Plus: TE-mumuv (TE-39) was created on the twig itself; mint date = today (2026-05-07).

Writes:
  docs/thought-experiments/README.md  - rebuilt index with proquint primary key
  protocols/<slug>.d/TODO/TODO.md     - rebuilt index per protocol
"""
import os, sys, csv, subprocess, re
from collections import defaultdict

REPO = sys.argv[1] if len(sys.argv) > 1 else "/home/user/workspace/wire-lab"
MAP = os.path.join(REPO, "tools/migrate-handles/mapping.tsv")


def git_first_date(path):
    """First commit date for a file (YYYY-MM-DD)."""
    try:
        out = subprocess.run(
            ["git", "-C", REPO, "log", "--diff-filter=A", "--follow",
             "--format=%ad", "--date=short", "--", path],
            check=True, capture_output=True, text=True).stdout.strip().splitlines()
        if out:
            return out[-1]  # oldest = last
    except subprocess.CalledProcessError:
        pass
    return ""


def title_from_file(path):
    """Extract title from first '# TE-... :' or '# TODO-... :' line, or first '# ' header."""
    full = os.path.join(REPO, path)
    if not os.path.exists(full):
        return ""
    with open(full) as f:
        for line in f:
            m = re.match(r"^#\s+(?:TE|TODO)-[a-z]{5,}:\s*(.+?)\s*$", line)
            if m:
                return m.group(1).strip()
            m = re.match(r"^#\s+(.+?)\s*$", line)
            if m:
                # Remove the prefix if present
                t = m.group(1)
                t = re.sub(r"^(?:TE|TODO)-[a-z]{5,}:\s*", "", t)
                return t.strip()
    return ""


def load_mapping():
    rows = []
    with open(MAP) as f:
        reader = csv.DictReader(f, delimiter="\t")
        for r in reader:
            rows.append(r)
    return rows


def fmt_link(rel_path, basename_only=True):
    if basename_only:
        return os.path.basename(rel_path)
    return rel_path


def build_te_index(rows):
    """Build the new docs/thought-experiments/README.md index table for TEs."""
    te_rows = [r for r in rows if r["kind"] == "TE"]
    # Sort: by int alias numerically when present (TE-1..TE-38), then TE-mumuv (TE-39) last
    def keyfn(r):
        m = re.match(r"TE-(\d+)$", r["int_alias"])
        return (int(m.group(1)) if m else 99999)
    te_rows.sort(key=keyfn)

    lines = []
    lines.append("| Handle | Mint date | Title | Prior alias |")
    lines.append("|---|---|---|---|")
    for r in te_rows:
        handle = f"TE-{r['handle']}"
        new_path = r["new_path"]
        rel = os.path.relpath(new_path, "docs/thought-experiments")
        mint = git_first_date(new_path) or "(uncommitted)"
        title = title_from_file(new_path) or "(no title)"
        ts = r['ts_alias']
        if ts and not ts.endswith('-mint'):
            prior = f"`{r['int_alias']}` / `{ts}`"
        else:
            prior = f"`{r['int_alias']}` (newly minted in TE-39 twig)"
        lines.append(f"| [{handle}]({rel}) | {mint} | {title} | {prior} |")
    return "\n".join(lines)


def build_todo_index(rows, protocol_slug):
    """Build the per-protocol TODO.md index table."""
    todo_rows = [r for r in rows if r["kind"] == "TODO"
                 and f"protocols/{protocol_slug}.d/TODO" in r["new_path"]]
    def keyfn(r):
        m = re.match(r"TODO-(\d+)$", r["int_alias"])
        return (int(m.group(1)) if m else 99999)
    todo_rows.sort(key=keyfn)

    lines = []
    lines.append("| Handle | Mint date | Title | Prior alias |")
    lines.append("|---|---|---|---|")
    for r in todo_rows:
        handle = f"TODO-{r['handle']}"
        new_path = r["new_path"]
        rel = os.path.relpath(new_path, f"protocols/{protocol_slug}.d/TODO")
        mint = git_first_date(new_path) or "(uncommitted)"
        title = title_from_file(new_path) or "(no title)"
        ts = r.get('ts_alias', '')
        if ts and not ts.endswith('-mint'):
            prior = f"`{r['int_alias']}` / `{ts}`"
        else:
            prior = f"`{r['int_alias']}`"
        lines.append(f"| [{handle}]({rel}) | {mint} | {title} | {prior} |")
    return "\n".join(lines)


def main():
    rows = load_mapping()
    print(f"Loaded {len(rows)} mapping rows", file=sys.stderr)

    # TE index
    te_table = build_te_index(rows)
    print("\n=== TE INDEX ===\n")
    print(te_table)

    # Per-protocol TODO indexes
    protocols = sorted(set(
        re.match(r"protocols/([^/]+)\.d/", r["new_path"]).group(1)
        for r in rows if r["kind"] == "TODO" and r["new_path"].startswith("protocols/")
    ))
    for slug in protocols:
        print(f"\n=== TODO INDEX: {slug} ===\n")
        print(build_todo_index(rows, slug))


if __name__ == "__main__":
    main()
