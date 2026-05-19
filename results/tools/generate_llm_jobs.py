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
import csv
import re
from dataclasses import dataclass
from pathlib import Path
from typing import List


REPO_ROOT = Path(__file__).resolve().parents[2]
DEFAULT_JOB_ROOT = REPO_ROOT / "results" / "jobs"
REQUIRED_FIELDS = {
    "run_group_id",
    "sim_id",
    "scenario_id",
    "model_id",
    "sim_path",
    "scenario_path",
}


@dataclass
class JobCell:
    run_group_id: str
    sim_id: str
    scenario_id: str
    model_id: str
    sim_path: str
    scenario_path: str


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
        default="<YYYYMMDD-HHMMSS>",
        help="Result timestamp placeholder or fixed timestamp for generated prompts.",
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


def slug(value: str) -> str:
    return re.sub(r"[^A-Za-z0-9_.-]+", "-", value).strip("-")


def read_manifest(path: Path) -> List[JobCell]:
    with path.open() as handle:
        reader = csv.DictReader(handle)
        if reader.fieldnames is None:
            raise RuntimeError("Manifest has no header.")
        missing = REQUIRED_FIELDS - set(reader.fieldnames)
        if missing:
            raise RuntimeError(f"Manifest missing required fields: {sorted(missing)}")
        rows = [
            JobCell(
                run_group_id=row["run_group_id"],
                sim_id=row["sim_id"],
                scenario_id=row["scenario_id"],
                model_id=row["model_id"],
                sim_path=row["sim_path"],
                scenario_path=row["scenario_path"],
            )
            for row in reader
        ]
    if not rows:
        raise RuntimeError("Manifest contains no rows.")
    return rows


def prompt_for(cell: JobCell, timestamp: str) -> str:
    result_path = f"results/{cell.sim_id}/{cell.scenario_id}/{cell.model_id}/{timestamp}.md"
    return f"""# LLM Matrix Cell Job

## Cell

- Run group ID: `{cell.run_group_id}`
- Simulation ID: `{cell.sim_id}`
- Scenario ID: `{cell.scenario_id}`
- Model ID: `{cell.model_id}`
- Intended result path: `{result_path}`

## Required Source Inputs

Read only source/design inputs before producing the verdict:

- `{cell.sim_path}README.md`
- `{cell.sim_path}QUESTION.md` if present
- local draft specs under `{cell.sim_path}` if present
- `{cell.scenario_path}`
- `results/RUN-PROTOCOL.md`

Do not read prior result files for this same sim/scenario cell before writing
the verdict. This job is blind with respect to prior results.

## Task

Evaluate the simulation against the scenario using deeper reasoning. Explain:

- what the simulation can actually cover,
- what obligations it pushes to another layer,
- where the scenario's 100-year, sparse-knowledge, no-central-authority,
  auditability, and migration pressures expose weaknesses,
- which open questions remain.

Write the result file at:

`{result_path}`

The result must follow the section contract in `results/RUN-PROTOCOL.md` and
must include:

- `- Run mode: llm-doc-eval-blind`
- a line starting with `Evidence verdict:`
- an explicit `Authority Boundary` section.
"""


def main() -> int:
    args = parse_args()
    cells = read_manifest(Path(args.manifest).resolve())
    if args.start_index:
        cells = cells[args.start_index :]
    if args.max_cells is not None:
        cells = cells[: args.max_cells]
    if not cells:
        raise RuntimeError("No cells selected after filters.")

    run_group_id = cells[0].run_group_id
    output_dir = Path(args.output_dir).resolve() if args.output_dir else DEFAULT_JOB_ROOT / run_group_id
    output_dir.mkdir(parents=True, exist_ok=True)

    index_path = output_dir / "INDEX.md"
    index_lines = [f"# LLM Jobs: {run_group_id}", ""]

    for index, cell in enumerate(cells, 1):
        filename = f"{index:05d}-{slug(cell.sim_id)}--{slug(cell.scenario_id)}--{slug(cell.model_id)}.md"
        path = output_dir / filename
        path.write_text(prompt_for(cell, args.timestamp))
        index_lines.append(f"- [{filename}]({filename})")

    index_path.write_text("\n".join(index_lines) + "\n")
    print(output_dir)
    print(f"jobs={len(cells)}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
