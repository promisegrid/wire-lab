#!/usr/bin/env python3
"""
Compatibility tombstone for retired scenario MATRIX.md updates.

Intent: Prevent old scripts from silently recreating duplicate scenario-side
result state after `results/` became the only canonical result evidence tree.
Use `tools/matrix-runner view` to generate read-only result views instead.
Source: DI-zamin
"""

from __future__ import annotations

import argparse
from pathlib import Path
from matrix_common import REPO_ROOT


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Retired scenario matrix updater.")
    parser.add_argument(
        "--result",
        action="append",
        required=True,
        help="Retired compatibility argument. May be passed more than once.",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Accepted for compatibility; no matrix updates are available.",
    )
    return parser.parse_args()


def update_matrix_for_result(result_path: Path, dry_run: bool = False) -> Path:
    """Reject obsolete matrix updates while preserving import compatibility."""

    rel = result_path
    if result_path.is_absolute():
        try:
            rel = result_path.relative_to(REPO_ROOT)
        except ValueError:
            rel = result_path
    raise RuntimeError(
        "scenario MATRIX.md files were retired by DI-zamin; "
        f"use `tools/matrix-runner view` for result navigation instead of updating {rel}"
    )


def main() -> int:
    args = parse_args()
    for result in args.result:
        result_path = Path(result)
        if not result_path.is_absolute():
            result_path = (REPO_ROOT / result_path).resolve()
        update_matrix_for_result(result_path, dry_run=args.dry_run)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
