#!/usr/bin/env python3
"""
Generate a deterministic scenario-by-simulation-by-model matrix manifest.

The output manifest is CSV and is intended to drive full-matrix batch runs.
"""

from __future__ import annotations

import argparse
import csv
import random
from pathlib import Path
from typing import Iterable, List

from matrix_common import (
    default_cell_id,
    default_result_dir,
    default_result_path,
    default_result_template,
    utc_compact_timestamp,
)


REPO_ROOT = Path(__file__).resolve().parents[2]
SIM_ROOT = REPO_ROOT / "simulations"
SCENARIO_ROOT = REPO_ROOT / "scenarios"
RESULTS_ROOT = REPO_ROOT / "results"
DEFAULT_MANIFEST_DIR = RESULTS_ROOT / "manifests"


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Generate matrix manifest CSV.")
    parser.add_argument(
        "--models",
        required=True,
        help="Comma-separated model IDs (for example openai-gpt-5.3-codex-xhigh).",
    )
    parser.add_argument(
        "--run-group-id",
        default="",
        help="Optional run group ID. Default: current UTC timestamp.",
    )
    parser.add_argument(
        "--timestamp",
        default="",
        help=(
            "Optional concrete result timestamp. Default: the generated run group "
            "timestamp, or current UTC if --run-group-id is supplied."
        ),
    )
    parser.add_argument(
        "--output",
        default="",
        help="Optional output CSV path. Default: results/manifests/matrix-manifest-<run-group-id>.csv",
    )
    parser.add_argument(
        "--sim-glob",
        default="SIM-*",
        help="Simulation directory glob. Default: SIM-*",
    )
    parser.add_argument(
        "--shuffle-seed",
        type=int,
        default=None,
        help="Optional deterministic shuffle seed. If set, rows are shuffled before limiting.",
    )
    parser.add_argument(
        "--limit-cells",
        type=int,
        default=None,
        help="Optional cap on emitted rows (useful for canary runs).",
    )
    return parser.parse_args()


def discover_sim_ids(sim_glob: str) -> List[str]:
    sim_ids = [
        path.name
        for path in SIM_ROOT.glob(sim_glob)
        if path.is_dir() and path.name.startswith("SIM-")
    ]
    return sorted(sim_ids)


def discover_scenario_ids() -> List[str]:
    scenario_ids = [
        path.name
        for path in SCENARIO_ROOT.iterdir()
        if path.is_dir() and (path / f"{path.name}.md").exists()
    ]
    return sorted(scenario_ids)


def iter_rows(
    run_group_id: str,
    timestamp: str,
    sim_ids: Iterable[str],
    scenario_ids: Iterable[str],
    model_ids: Iterable[str],
) -> List[dict]:
    rows: List[dict] = []
    for sim_id in sim_ids:
        sim_path = f"simulations/{sim_id}/"
        for scenario_id in scenario_ids:
            scenario_path = f"scenarios/{scenario_id}/{scenario_id}.md"
            for model_id in model_ids:
                rows.append(
                    {
                        "run_group_id": run_group_id,
                        "sim_id": sim_id,
                        "scenario_id": scenario_id,
                        "model_id": model_id,
                        "sim_path": sim_path,
                        "scenario_path": scenario_path,
                        "result_dir": default_result_dir(sim_id, scenario_id, model_id),
                        "result_path_template": default_result_template(
                            sim_id, scenario_id, model_id
                        ),
                        "timestamp": timestamp,
                        "result_path": default_result_path(
                            sim_id, scenario_id, model_id, timestamp
                        ),
                        "status": "queued",
                    }
                )
    return rows


def assign_manifest_ordinals(run_group_id: str, rows: List[dict]) -> None:
    """Assign queue order after deterministic shuffle/limit operations.

    Intent: A full unattended run needs stable per-row IDs in the exact emitted
    queue order, including canary manifests that are shuffled or truncated.
    Source: DI-nuhon
    """

    for ordinal, row in enumerate(rows, 1):
        row["ordinal"] = str(ordinal)
        row["cell_id"] = default_cell_id(
            run_group_id,
            ordinal,
            row["sim_id"],
            row["scenario_id"],
            row["model_id"],
        )


def main() -> int:
    args = parse_args()
    generated_timestamp = utc_compact_timestamp()
    run_group_id = args.run_group_id or generated_timestamp
    timestamp = args.timestamp or (generated_timestamp if args.run_group_id else run_group_id)
    model_ids = [model.strip() for model in args.models.split(",") if model.strip()]
    if not model_ids:
        raise RuntimeError("At least one model ID is required.")

    sim_ids = discover_sim_ids(args.sim_glob)
    scenario_ids = discover_scenario_ids()
    if not sim_ids:
        raise RuntimeError(f"No simulations matched glob '{args.sim_glob}'.")
    if not scenario_ids:
        raise RuntimeError("No scenario entries found under scenarios/.")

    rows = iter_rows(run_group_id, timestamp, sim_ids, scenario_ids, model_ids)
    if args.shuffle_seed is not None:
        random.Random(args.shuffle_seed).shuffle(rows)
    if args.limit_cells is not None:
        rows = rows[: args.limit_cells]
    assign_manifest_ordinals(run_group_id, rows)

    if args.output:
        out_path = Path(args.output).resolve()
    else:
        out_path = DEFAULT_MANIFEST_DIR / f"matrix-manifest-{run_group_id}.csv"
    out_path.parent.mkdir(parents=True, exist_ok=True)

    fieldnames = [
        "run_group_id",
        "ordinal",
        "cell_id",
        "sim_id",
        "scenario_id",
        "model_id",
        "sim_path",
        "scenario_path",
        "result_dir",
        "result_path_template",
        "timestamp",
        "result_path",
        "status",
    ]
    with out_path.open("w", newline="") as handle:
        writer = csv.DictWriter(handle, fieldnames=fieldnames)
        writer.writeheader()
        writer.writerows(rows)

    print(out_path)
    print(
        f"rows={len(rows)} sims={len(sim_ids)} scenarios={len(scenario_ids)} "
        f"models={len(model_ids)} timestamp={timestamp}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
