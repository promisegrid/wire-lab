#!/usr/bin/env python3
"""
Generate LLM evaluation job prompts from a matrix manifest.

This tool does not write result verdict files. It writes prompt files that can be
fed to Codex or another LLM runner. Final result files must be produced by an
LLM/human evaluator following results/RUN-PROTOCOL.md.

Intent: Keep batch tooling limited to work preparation so root result evidence is
created by deeper LLM or human reasoning, not mechanical parsing.
Source: DI-moduf
"""

from __future__ import annotations

import argparse
from pathlib import Path

from matrix_common import (
    DEFAULT_JOB_ROOT,
    TIMESTAMP_PLACEHOLDER,
    iter_selected,
    prompt_filename,
    prompt_for,
    read_manifest,
    with_timestamp,
)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Generate LLM job prompts from a manifest.")
    parser.add_argument("--manifest", required=True, help="Matrix manifest CSV path.")
    parser.add_argument(
        "--output-dir",
        default="",
        help="Optional output directory. Default: results/jobs/<run_group_id>/",
    )
    parser.add_argument(
        "--timestamp",
        default=TIMESTAMP_PLACEHOLDER,
        help=(
            "Optional fixed timestamp override. Default: use concrete manifest "
            "result_path/timestamp fields, or retain the placeholder for old manifests."
        ),
    )
    parser.add_argument(
        "--max-cells",
        type=int,
        default=None,
        help="Optional cap on generated prompts.",
    )
    parser.add_argument(
        "--start-index",
        type=int,
        default=0,
        help="Optional start index for resumable prompt generation.",
    )
    return parser.parse_args()


def main() -> int:
    args = parse_args()
    cells = read_manifest(Path(args.manifest).resolve())
    cells = [with_timestamp(cell, args.timestamp) for cell in cells]
    cells = iter_selected(cells, args.start_index, args.max_cells)

    run_group_id = cells[0].run_group_id
    output_dir = Path(args.output_dir).resolve() if args.output_dir else DEFAULT_JOB_ROOT / run_group_id
    output_dir.mkdir(parents=True, exist_ok=True)

    index_path = output_dir / "INDEX.md"
    index_lines = [f"# LLM Jobs: {run_group_id}", ""]

    for cell in cells:
        filename = prompt_filename(cell)
        path = output_dir / filename
        path.write_text(prompt_for(cell))
        index_lines.append(f"- [{filename}]({filename})")

    index_path.write_text("\n".join(index_lines) + "\n")
    print(output_dir)
    print(f"jobs={len(cells)}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
