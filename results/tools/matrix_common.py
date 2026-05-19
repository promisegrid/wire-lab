#!/usr/bin/env python3
"""
Shared helpers for scenario-by-simulation matrix tooling.

Intent: Keep full-matrix runs deterministic and resumable without letting
tooling synthesize design verdict prose. Manifest, prompt, queue, and
validation helpers may coordinate work; result substance still comes from an
LLM or human evaluator reading the cell inputs.
Source: DI-nuhon
"""

from __future__ import annotations

import csv
import re
from dataclasses import dataclass, replace
from datetime import datetime, timezone
from pathlib import Path
from typing import Iterable, List, Optional


REPO_ROOT = Path(__file__).resolve().parents[2]
RESULTS_ROOT = REPO_ROOT / "results"
DEFAULT_JOB_ROOT = RESULTS_ROOT / "jobs"
DEFAULT_STATE_ROOT = RESULTS_ROOT / "state"
TIMESTAMP_PLACEHOLDER = "<YYYYMMDD-HHMMSS>"

REQUIRED_MANIFEST_FIELDS = {
    "run_group_id",
    "sim_id",
    "scenario_id",
    "model_id",
    "sim_path",
    "scenario_path",
}


@dataclass(frozen=True)
class MatrixCell:
    """One concrete matrix cell from a manifest.

    Intent: Carry both human-readable cell coordinates and concrete run paths
    so queue runners do not have to infer missing timestamps while unattended.
    Source: DI-nuhon
    """

    run_group_id: str
    sim_id: str
    scenario_id: str
    model_id: str
    sim_path: str
    scenario_path: str
    result_dir: str
    result_path_template: str
    timestamp: str
    result_path: str
    status: str
    ordinal: int
    cell_id: str


def utc_compact_timestamp() -> str:
    """Return the compact UTC timestamp used for result file names."""

    return datetime.now(timezone.utc).strftime("%Y%m%d-%H%M%S")


def slug(value: str) -> str:
    """Convert a matrix coordinate into a path-safe slug fragment."""

    return re.sub(r"[^A-Za-z0-9_.-]+", "-", value).strip("-")


def repo_relative_path(path_text: str) -> Path:
    """Resolve a repo-relative or absolute path without changing its target."""

    path = Path(path_text)
    if path.is_absolute():
        return path
    return (REPO_ROOT / path).resolve()


def default_cell_id(run_group_id: str, ordinal: int, sim_id: str, scenario_id: str, model_id: str) -> str:
    """Create a stable cell ID for queue state and prompt filenames."""

    return (
        f"{slug(run_group_id)}-{ordinal:06d}-"
        f"{slug(sim_id)}--{slug(scenario_id)}--{slug(model_id)}"
    )


def default_result_dir(sim_id: str, scenario_id: str, model_id: str) -> str:
    """Return the locked result directory for a matrix cell."""

    return f"results/{sim_id}/{scenario_id}/{model_id}/"


def default_result_path(sim_id: str, scenario_id: str, model_id: str, timestamp: str) -> str:
    """Return the locked concrete result path for a matrix cell."""

    return f"{default_result_dir(sim_id, scenario_id, model_id)}{timestamp}.md"


def default_result_template(sim_id: str, scenario_id: str, model_id: str) -> str:
    """Return the historical placeholder template for compatibility."""

    return default_result_path(sim_id, scenario_id, model_id, TIMESTAMP_PLACEHOLDER)


def row_result_path(row: dict) -> Optional[str]:
    """Resolve a manifest row to a concrete result path when possible.

    Intent: Manifest validation and queue progress must not silently accept a
    placeholder-only path as if a result exists; unattended runs need a concrete
    destination per cell.
    Source: DI-nuhon
    """

    result_path = row.get("result_path", "").strip()
    if result_path and TIMESTAMP_PLACEHOLDER not in result_path:
        return result_path

    timestamp = row.get("timestamp", "").strip()
    template = row.get("result_path_template", "").strip()
    if timestamp and template:
        return template.replace(TIMESTAMP_PLACEHOLDER, timestamp)

    if timestamp:
        return default_result_path(
            row.get("sim_id", "").strip(),
            row.get("scenario_id", "").strip(),
            row.get("model_id", "").strip(),
            timestamp,
        )
    return None


def cell_from_row(row: dict, ordinal: int) -> MatrixCell:
    """Normalize old and new manifest rows into a MatrixCell."""

    sim_id = row["sim_id"].strip()
    scenario_id = row["scenario_id"].strip()
    model_id = row["model_id"].strip()
    row_ordinal = int(row.get("ordinal", "") or ordinal)
    timestamp = row.get("timestamp", "").strip()
    result_dir = row.get("result_dir", "").strip() or default_result_dir(sim_id, scenario_id, model_id)
    result_template = (
        row.get("result_path_template", "").strip()
        or default_result_template(sim_id, scenario_id, model_id)
    )
    result_path = row_result_path(row) or default_result_path(
        sim_id, scenario_id, model_id, timestamp or TIMESTAMP_PLACEHOLDER
    )
    cell_id = row.get("cell_id", "").strip() or default_cell_id(
        row["run_group_id"].strip(), row_ordinal, sim_id, scenario_id, model_id
    )

    return MatrixCell(
        run_group_id=row["run_group_id"].strip(),
        sim_id=sim_id,
        scenario_id=scenario_id,
        model_id=model_id,
        sim_path=row["sim_path"].strip(),
        scenario_path=row["scenario_path"].strip(),
        result_dir=result_dir,
        result_path_template=result_template,
        timestamp=timestamp,
        result_path=result_path,
        status=row.get("status", "").strip() or "queued",
        ordinal=row_ordinal,
        cell_id=cell_id,
    )


def read_manifest(path: Path) -> List[MatrixCell]:
    """Read a matrix manifest and preserve row order as queue order."""

    with path.open(newline="") as handle:
        reader = csv.DictReader(handle)
        if reader.fieldnames is None:
            raise RuntimeError("Manifest has no header.")
        missing = REQUIRED_MANIFEST_FIELDS - set(reader.fieldnames)
        if missing:
            raise RuntimeError(f"Manifest missing required fields: {sorted(missing)}")
        cells = [cell_from_row(row, index) for index, row in enumerate(reader, 1)]
    if not cells:
        raise RuntimeError("Manifest contains no rows.")
    return cells


def with_timestamp(cell: MatrixCell, timestamp: str) -> MatrixCell:
    """Return a cell with a concrete result timestamp override."""

    if not timestamp or timestamp == TIMESTAMP_PLACEHOLDER:
        return cell
    return replace(
        cell,
        timestamp=timestamp,
        result_path=default_result_path(cell.sim_id, cell.scenario_id, cell.model_id, timestamp),
    )


def prompt_filename(cell: MatrixCell) -> str:
    """Return the deterministic prompt filename for a queue cell."""

    return f"{cell.ordinal:05d}-{slug(cell.sim_id)}--{slug(cell.scenario_id)}--{slug(cell.model_id)}.md"


def prompt_for(cell: MatrixCell) -> str:
    """Build the blind LLM evaluation prompt for a concrete matrix cell."""

    return f"""# LLM Matrix Cell Job

## Cell

- Cell ID: `{cell.cell_id}`
- Run group ID: `{cell.run_group_id}`
- Queue ordinal: `{cell.ordinal}`
- Simulation ID: `{cell.sim_id}`
- Scenario ID: `{cell.scenario_id}`
- Model ID: `{cell.model_id}`
- Intended result path: `{cell.result_path}`

## Required Source Inputs

Read only source/design inputs before producing the verdict:

- `{cell.sim_path}README.md`
- `{cell.sim_path}QUESTION.md` if present
- local draft specs under `{cell.sim_path}` if present
- `scenarios/README.md`
- `scenarios/{cell.scenario_id}/README.md`
- local scenario docs under `scenarios/{cell.scenario_id}/` if present
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

`{cell.result_path}`

This is an unattended batch cell. Do not ask for confirmation, do not wait for
interactive approval, and do not write a conversational answer instead of the
result file.

The result must follow the section contract in `results/RUN-PROTOCOL.md` and
must include:

- `- Run mode: llm-doc-eval-blind`
- a line starting with `Evidence verdict:`
- an explicit `Authority Boundary` section.
"""


def iter_selected(cells: Iterable[MatrixCell], start_index: int, max_cells: Optional[int]) -> List[MatrixCell]:
    """Apply stable start/limit slicing used by job and queue tools."""

    selected = list(cells)
    if start_index:
        selected = selected[start_index:]
    if max_cells is not None:
        selected = selected[:max_cells]
    if not selected:
        raise RuntimeError("No cells selected after filters.")
    return selected
