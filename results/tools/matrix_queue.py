#!/usr/bin/env python3
"""
Run or inspect an unattended full-matrix queue.

Intent: A large matrix run must be resumable and checkpointed cell-by-cell so
Steve does not have to press keys through thousands of prompts, while each
result file is still produced by an external LLM or human-authorized runner.
Committed scenario matrices were retired; the queue validates result files and
leaves result navigation to generated views. Source: DI-nuhon; DI-zamin
"""

from __future__ import annotations

import argparse
import json
import os
import shlex
import subprocess
from collections import Counter
from datetime import datetime, timezone
from pathlib import Path
from typing import Dict, List, Tuple

from matrix_common import (
    DEFAULT_JOB_ROOT,
    DEFAULT_STATE_ROOT,
    REPO_ROOT,
    MatrixCell,
    iter_selected,
    prompt_filename,
    prompt_for,
    read_manifest,
    repo_relative_path,
    slug,
)
from validate_results import validate_file


TERMINAL_STATUSES = {"done", "skipped"}


def utc_iso_timestamp() -> str:
    """Return an audit-friendly UTC timestamp for queue state events."""

    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Run or inspect a matrix queue.")
    subparsers = parser.add_subparsers(dest="command", required=True)

    run_parser = subparsers.add_parser("run", help="Run queued cells with an external LLM command.")
    add_common_queue_args(run_parser)
    run_parser.add_argument(
        "--runner-command",
        default="",
        help=(
            "External command template run once per cell. Placeholders include "
            "{prompt_path}, {result_path}, {cell_id}, {sim_id}, {scenario_id}, "
            "{model_id}, and {run_group_id}. The command is split with shlex, not a shell."
        ),
    )
    run_parser.add_argument("--start-index", type=int, default=0, help="Zero-based queue start index.")
    run_parser.add_argument("--limit", type=int, default=None, help="Maximum selected cells this invocation.")
    run_parser.add_argument(
        "--retry-failed",
        action="store_true",
        help="Retry cells previously marked failed instead of leaving them for review.",
    )
    run_parser.add_argument(
        "--rerun-done",
        action="store_true",
        help="Run cells already marked done. This should normally be avoided.",
    )
    run_parser.add_argument(
        "--no-matrix-update",
        action="store_true",
        help="Deprecated no-op; scenario MATRIX.md files were retired.",
    )
    run_parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Show selected work without writing prompts, state, or results.",
    )

    progress_parser = subparsers.add_parser("progress", help="Print queue status counts.")
    add_common_queue_args(progress_parser)
    return parser.parse_args()


def add_common_queue_args(parser: argparse.ArgumentParser) -> None:
    parser.add_argument("--manifest", required=True, help="Concrete matrix manifest CSV path.")
    parser.add_argument(
        "--state",
        default="",
        help="Queue state JSON path. Default: results/state/<run-group-id>.json",
    )
    parser.add_argument(
        "--job-dir",
        default="",
        help="Prompt directory. Default: results/jobs/<run-group-id>/",
    )


def default_state_path(cells: List[MatrixCell]) -> Path:
    """Return the default checkpoint path for a manifest's run group."""

    return DEFAULT_STATE_ROOT / f"{slug(cells[0].run_group_id)}.json"


def default_job_dir(cells: List[MatrixCell]) -> Path:
    """Return the default prompt directory for a manifest's run group."""

    return DEFAULT_JOB_ROOT / cells[0].run_group_id


def state_record_for(cell: MatrixCell) -> dict:
    """Create the persistent state payload for one queue cell."""

    return {
        "cell_id": cell.cell_id,
        "ordinal": cell.ordinal,
        "sim_id": cell.sim_id,
        "scenario_id": cell.scenario_id,
        "model_id": cell.model_id,
        "result_path": cell.result_path,
        "status": cell.status or "queued",
        "attempts": 0,
        "last_message": "",
        "updated_at": "",
    }


def load_or_create_state(state_path: Path, manifest_path: Path, cells: List[MatrixCell]) -> dict:
    """Load existing state and add newly discovered manifest cells."""

    if state_path.exists():
        state = json.loads(state_path.read_text())
    else:
        state = {
            "manifest": manifest_path.relative_to(REPO_ROOT).as_posix()
            if manifest_path.is_relative_to(REPO_ROOT)
            else str(manifest_path),
            "run_group_id": cells[0].run_group_id,
            "created_at": utc_iso_timestamp(),
            "updated_at": "",
            "cells": {},
        }

    cell_state: Dict[str, dict] = state.setdefault("cells", {})
    for cell in cells:
        record = cell_state.setdefault(cell.cell_id, state_record_for(cell))
        record["ordinal"] = cell.ordinal
        record["sim_id"] = cell.sim_id
        record["scenario_id"] = cell.scenario_id
        record["model_id"] = cell.model_id
        record["result_path"] = cell.result_path
    return state


def save_state(state_path: Path, state: dict) -> None:
    """Persist checkpoint state atomically after each cell."""

    state["updated_at"] = utc_iso_timestamp()
    state_path.parent.mkdir(parents=True, exist_ok=True)
    tmp_path = state_path.with_suffix(state_path.suffix + ".tmp")
    tmp_path.write_text(json.dumps(state, indent=2, sort_keys=True) + "\n")
    tmp_path.replace(state_path)


def status_counts(state: dict) -> Counter:
    """Count queue statuses for progress output."""

    return Counter(record.get("status", "queued") for record in state.get("cells", {}).values())


def print_progress(state: dict) -> None:
    """Print stable status counts for humans and logs."""

    counts = status_counts(state)
    total = sum(counts.values())
    parts = [f"{status}={counts.get(status, 0)}" for status in sorted(counts)]
    print(f"total={total} " + " ".join(parts))


def write_prompt(job_dir: Path, cell: MatrixCell) -> Path:
    """Write or refresh the prompt file for a single cell."""

    job_dir.mkdir(parents=True, exist_ok=True)
    prompt_path = job_dir / prompt_filename(cell)
    prompt_path.write_text(prompt_for(cell))
    return prompt_path


def validate_result_path(result_path: Path) -> List[str]:
    """Validate one result file without requiring matrix linkage yet."""

    if not result_path.exists():
        return ["result file was not created"]
    return validate_file(result_path, strict_matrix=False, allow_prototype=False)


def command_argv(template: str, cell: MatrixCell, prompt_path: Path, result_path: Path) -> List[str]:
    """Expand a runner-command template into argv without using a shell."""

    values = {
        "prompt_path": str(prompt_path),
        "result_path": str(result_path),
        "cell_id": cell.cell_id,
        "sim_id": cell.sim_id,
        "scenario_id": cell.scenario_id,
        "model_id": cell.model_id,
        "run_group_id": cell.run_group_id,
    }
    return shlex.split(template.format(**values))


def runner_environment(cell: MatrixCell, prompt_path: Path, result_path: Path) -> dict:
    """Expose cell coordinates to external runners via environment variables."""

    env = os.environ.copy()
    env.update(
        {
            "WIRE_LAB_CELL_ID": cell.cell_id,
            "WIRE_LAB_RUN_GROUP_ID": cell.run_group_id,
            "WIRE_LAB_SIM_ID": cell.sim_id,
            "WIRE_LAB_SCENARIO_ID": cell.scenario_id,
            "WIRE_LAB_MODEL_ID": cell.model_id,
            "WIRE_LAB_PROMPT_PATH": str(prompt_path),
            "WIRE_LAB_RESULT_PATH": str(result_path),
        }
    )
    return env


def mark(record: dict, status: str, message: str = "") -> None:
    """Update a cell state record with audit metadata."""

    record["status"] = status
    record["last_message"] = message
    record["updated_at"] = utc_iso_timestamp()


def should_skip_record(record: dict, args: argparse.Namespace) -> Tuple[bool, str]:
    """Decide whether a state record should be left untouched this run."""

    status = record.get("status", "queued")
    if status == "done" and not args.rerun_done:
        return True, "already done"
    if status == "failed" and not args.retry_failed:
        return True, "failed; pass --retry-failed to retry"
    if status in TERMINAL_STATUSES and not args.rerun_done:
        return True, f"terminal status {status}"
    return False, ""


def run_queue(args: argparse.Namespace) -> int:
    manifest_path = Path(args.manifest).resolve()
    cells = read_manifest(manifest_path)
    cells = iter_selected(cells, args.start_index, args.limit)
    state_path = Path(args.state).resolve() if args.state else default_state_path(cells)
    job_dir = Path(args.job_dir).resolve() if args.job_dir else default_job_dir(cells)

    state = load_or_create_state(state_path, manifest_path, cells)
    processed = 0
    failed = 0

    for cell in cells:
        record = state["cells"][cell.cell_id]
        skip, reason = should_skip_record(record, args)
        if skip:
            print(f"skip: {cell.cell_id}: {reason}")
            continue

        status = run_cell(
            cell,
            record,
            state,
            state_path,
            job_dir,
            args.runner_command,
            dry_run=args.dry_run,
        )
        processed += 1
        if status == "failed":
            failed += 1
        if not args.dry_run:
            save_state(state_path, state)
        print_progress(state)

    if not args.dry_run:
        save_state(state_path, state)
    print(f"processed={processed} failed={failed} state={state_path}")
    return 1 if failed else 0


def run_cell(
    cell: MatrixCell,
    record: dict,
    state: dict,
    state_path: Path,
    job_dir: Path,
    runner_command: str,
    dry_run: bool,
) -> str:
    """Run one cell and persist `running` before external LLM launch.

    Intent: If the process is interrupted while an LLM call is in flight, the
    checkpoint should name the active cell instead of leaving it indistinguishably
    queued. Source: DI-bujiv
    """

    result_path = repo_relative_path(cell.result_path)

    if dry_run:
        print(f"dry-run: {cell.cell_id} -> {cell.result_path}")
        return record.get("status", "queued")

    prompt_path = write_prompt(job_dir, cell)
    existing_issues = validate_result_path(result_path)
    if not existing_issues:
        mark(record, "done", "existing valid result")
        return "done"

    if not runner_command:
        mark(record, "failed", "missing --runner-command and no valid existing result")
        return "failed"

    record["attempts"] = int(record.get("attempts", 0)) + 1
    mark(record, "running", "runner command started")
    save_state(state_path, state)

    argv = command_argv(runner_command, cell, prompt_path, result_path)
    completed = subprocess.run(
        argv,
        cwd=REPO_ROOT,
        env=runner_environment(cell, prompt_path, result_path),
        check=False,
    )
    if completed.returncode != 0:
        mark(record, "failed", f"runner exited {completed.returncode}")
        return "failed"

    issues = validate_result_path(result_path)
    if issues:
        mark(record, "failed", "; ".join(issues))
        return "failed"

    mark(record, "done", "validated result")
    return "done"


def show_progress(args: argparse.Namespace) -> int:
    manifest_path = Path(args.manifest).resolve()
    cells = read_manifest(manifest_path)
    state_path = Path(args.state).resolve() if args.state else default_state_path(cells)
    state = load_or_create_state(state_path, manifest_path, cells)
    print_progress(state)
    print(f"state={state_path}")
    return 0


def main() -> int:
    args = parse_args()
    if args.command == "run":
        return run_queue(args)
    if args.command == "progress":
        return show_progress(args)
    raise RuntimeError(f"Unknown command: {args.command}")


if __name__ == "__main__":
    raise SystemExit(main())
