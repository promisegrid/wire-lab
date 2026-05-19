#!/usr/bin/env python3
"""
Validate result artifacts for required shape and matrix linkage.

Intent: Result validation must reject parser-generated prototype artifacts by
default so scripted plumbing tests cannot masquerade as design evidence.
Source: DI-moduf
"""

from __future__ import annotations

import argparse
import csv
from dataclasses import dataclass
from pathlib import Path
from typing import Iterable, List, Tuple

from matrix_common import REPO_ROOT, RESULTS_ROOT, TIMESTAMP_PLACEHOLDER, row_result_path


REQUIRED_HEADINGS = [
    "## Result ID",
    "## Scenario",
    "## Simulation",
    "## Runner",
    "## Prompt / Procedure",
    "## Observed Behavior",
    "## Verdict",
    "## Evidence Links",
    "## Open Questions",
    "## Handoff Target",
    "## Authority Boundary",
]


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Validate result artifact files.")
    parser.add_argument(
        "--manifest",
        default="",
        help="Optional manifest path; validates only rows in this manifest.",
    )
    parser.add_argument("--model", default="", help="Optional model filter.")
    parser.add_argument("--timestamp", default="", help="Optional timestamp filter.")
    parser.add_argument(
        "--strict-matrix",
        action="store_true",
        help="Require each result path to be referenced in its scenario MATRIX.md row.",
    )
    parser.add_argument(
        "--allow-prototype",
        action="store_true",
        help="Allow scripted prototype result files for plumbing checks.",
    )
    parser.add_argument(
        "--allow-missing",
        action="store_true",
        help="For manifest validation, skip missing result files instead of failing them.",
    )
    parser.add_argument("--max-errors", type=int, default=50, help="Maximum printed errors.")
    return parser.parse_args()


@dataclass(frozen=True)
class ManifestTargets:
    """Resolved manifest result targets plus unresolved/missing rows."""

    paths: List[Path]
    missing: List[Tuple[str, str]]


def resolve_manifest_targets(manifest_path: Path, allow_missing: bool) -> ManifestTargets:
    """Resolve concrete manifest rows and report missing cells explicitly.

    Intent: Full unattended validation must not claim success after silently
    filtering out missing result files from a 7000-cell manifest.
    Source: DI-nuhon
    """

    paths: List[Path] = []
    missing: List[Tuple[str, str]] = []
    with manifest_path.open() as handle:
        reader = csv.DictReader(handle)
        for row in reader:
            result_path = row_result_path(row)
            cell_id = row.get("cell_id", "") or (
                f"{row.get('run_group_id', '')}/"
                f"{row.get('sim_id', '')}/"
                f"{row.get('scenario_id', '')}/"
                f"{row.get('model_id', '')}"
            )
            if not result_path or TIMESTAMP_PLACEHOLDER in result_path:
                missing.append((cell_id, "manifest row has no concrete result_path/timestamp"))
                continue
            target = Path(result_path)
            if not target.is_absolute():
                target = (REPO_ROOT / target).resolve()
            if not target.exists():
                if not allow_missing:
                    missing.append((cell_id, f"missing result file: {target.relative_to(REPO_ROOT)}"))
                continue
            paths.append(target)
    return ManifestTargets(paths=paths, missing=missing)


def iter_targets(model: str, timestamp: str) -> Iterable[Path]:
    for path in RESULTS_ROOT.rglob("*.md"):
        if path == RESULTS_ROOT / "README.md":
            continue
        rel = path.relative_to(RESULTS_ROOT).parts
        if len(rel) != 4:
            continue
        sim_id, scenario_id, model_id, filename = rel
        ts = filename.replace(".md", "")
        if model and model_id != model:
            continue
        if timestamp and ts != timestamp:
            continue
        if not sim_id.startswith("SIM-"):
            continue
        if not scenario_id:
            continue
        yield path


def matrix_contains_result(result_path: Path) -> bool:
    rel = result_path.relative_to(REPO_ROOT)
    parts = rel.parts
    # results/<sim>/<scenario>/<model>/<ts>.md
    scenario_id = parts[2]
    matrix = REPO_ROOT / "scenarios" / scenario_id / "MATRIX.md"
    if not matrix.exists():
        return False
    return str(rel) in matrix.read_text()


def is_prototype_result(text: str) -> bool:
    """Return true for known scripted plumbing-test artifacts.

    Intent: Prototype detection is content-based so preserved old files remain
    auditable while validators and comparison reports exclude them by default.
    Source: DI-moduf
    """
    prototype_markers = [
        "Run mode: `scripted-doc-eval-blind`",
        "Runner/interface: `results/tools/run_matrix_batch.py`",
    ]
    return any(marker in text for marker in prototype_markers)


def validate_file(path: Path, strict_matrix: bool, allow_prototype: bool) -> List[str]:
    errors: List[str] = []
    text = path.read_text()
    lines = text.splitlines()

    if is_prototype_result(text) and not allow_prototype:
        errors.append(
            "prototype scripted result excluded by default; "
            "pass --allow-prototype for plumbing checks"
        )

    rel_parts = path.relative_to(RESULTS_ROOT).parts
    if len(rel_parts) != 4:
        errors.append("path shape must be results/<sim>/<scenario>/<model>/<timestamp>.md")
        return errors
    sim_id, scenario_id, model_id, filename = rel_parts
    timestamp = filename.replace(".md", "")

    if not lines or not lines[0].startswith("# Result: "):
        errors.append("missing # Result header")
    if f"/ {model_id} / {timestamp}" not in (lines[0] if lines else ""):
        errors.append("header does not match path model/timestamp")

    for heading in REQUIRED_HEADINGS:
        if heading not in text:
            errors.append(f"missing heading: {heading}")

    if "Evidence verdict:" not in text:
        errors.append("missing Evidence verdict line")

    model_line = next((line for line in lines if line.startswith("- Model ID: ")), "")
    if model_id not in model_line:
        errors.append("Model ID line does not match path model")

    if strict_matrix and not matrix_contains_result(path):
        errors.append("scenario matrix does not reference this result path")

    return errors


def main() -> int:
    args = parse_args()
    missing: List[Tuple[str, str]] = []
    if args.manifest:
        manifest_targets = resolve_manifest_targets(
            Path(args.manifest).resolve(),
            allow_missing=args.allow_missing,
        )
        targets = manifest_targets.paths
        missing = manifest_targets.missing
    else:
        targets = list(iter_targets(args.model, args.timestamp))
    if not targets and not missing and not (args.manifest and args.allow_missing):
        raise RuntimeError("No result files matched selection.")

    total = 0
    bad = len(missing)
    printed = 0
    for cell_id, issue in missing:
        if printed < args.max_errors:
            print(f"{cell_id}:")
            print(f"  - {issue}")
            printed += 1

    for path in sorted(targets):
        total += 1
        issues = validate_file(
            path,
            strict_matrix=args.strict_matrix,
            allow_prototype=args.allow_prototype,
        )
        if issues:
            bad += 1
            if printed < args.max_errors:
                rel = path.relative_to(REPO_ROOT)
                print(f"{rel}:")
                for issue in issues:
                    print(f"  - {issue}")
                printed += 1

    print(f"validated={total} failed={bad}")
    return 1 if bad else 0


if __name__ == "__main__":
    raise SystemExit(main())
