#!/usr/bin/env python3
"""
Compare verdict drift between two result-model corpora.

Intent: Model comparisons must summarize real LLM/human result evidence and
exclude scripted prototype plumbing outputs by default. Source: DI-moduf

Usage:
  python3 results/comparisons/compare_model_results.py \
    --old-model openai-gpt-5.5-xhigh \
    --new-model openai-gpt-5.3-codex-xhigh

By default the tool picks the latest timestamp per (sim, scenario, model) cell.
You can pin timestamps with --old-ts / --new-ts.
"""

from __future__ import annotations

import argparse
from collections import defaultdict
from dataclasses import dataclass
from pathlib import Path
from typing import Dict, Iterable, List, Tuple


REPO_ROOT = Path(__file__).resolve().parents[2]
RESULTS_ROOT = REPO_ROOT / "results"
DEFAULT_OUTPUT_DIR = RESULTS_ROOT / "comparisons"


SCORE_RULES: List[Tuple[str, int]] = [
    ("strongest fit", 6),
    ("best migration fit", 6),
    ("best fit", 6),
    ("strong fit", 5),
    ("good partial fit", 4),
    ("good strict baseline", 4),
    ("partial fit", 3),
    ("partial guardrail", 3),
    ("partial with privacy risk", 2),
    ("partial but brittle", 2),
    ("weak-to-partial fit", 2),
    ("weak fit", 1),
    ("poor standalone fit", 0),
    ("poor fit", 0),
    ("negative-control", 1),
]


@dataclass(frozen=True)
class CellKey:
    sim: str
    scenario: str


@dataclass
class ResultCell:
    key: CellKey
    model: str
    timestamp: str
    path: Path
    verdict: str
    verdict_line: int
    score: int


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Compare model-result verdict drift.")
    parser.add_argument("--old-model", required=True, help="Baseline model slug.")
    parser.add_argument("--new-model", required=True, help="Comparison model slug.")
    parser.add_argument(
        "--old-ts",
        default="latest",
        help="Baseline timestamp (YYYYMMDD-HHMMSS) or 'latest' per cell.",
    )
    parser.add_argument(
        "--new-ts",
        default="latest",
        help="Comparison timestamp (YYYYMMDD-HHMMSS) or 'latest' per cell.",
    )
    parser.add_argument(
        "--output",
        default="",
        help="Optional output markdown path. Default: results/comparisons/<old>_vs_<new>_<timestamp>.md",
    )
    parser.add_argument(
        "--include-prototype",
        action="store_true",
        help="Include scripted prototype plumbing outputs in the comparison.",
    )
    return parser.parse_args()


def score_verdict(verdict: str) -> int:
    lower = verdict.lower()
    for phrase, value in SCORE_RULES:
        if phrase in lower:
            return value
    if "strong" in lower:
        return 5
    if "good" in lower:
        return 4
    if "partial" in lower:
        return 3
    if "weak" in lower:
        return 1
    if "poor" in lower:
        return 0
    return 2


def extract_verdict(path: Path) -> Tuple[str, int]:
    for line_no, line in enumerate(path.read_text().splitlines(), 1):
        if line.startswith("Evidence verdict:"):
            verdict = line.removeprefix("Evidence verdict:").strip().rstrip(".")
            return verdict, line_no
    raise RuntimeError(f"Missing verdict line in {path}")


def is_prototype_result(path: Path) -> bool:
    """Return true for known scripted plumbing-test artifacts.

    Intent: Latest-result selection must not accidentally prefer a preserved
    prototype file over an older real reasoning result. Source: DI-moduf
    """
    text = path.read_text()
    prototype_markers = [
        "Run mode: `scripted-doc-eval-blind`",
        "Runner/interface: `results/tools/run_matrix_batch.py`",
    ]
    return any(marker in text for marker in prototype_markers)


def iter_model_files(model: str, include_prototype: bool) -> Iterable[Path]:
    for path in RESULTS_ROOT.rglob("*.md"):
        if path == RESULTS_ROOT / "README.md":
            continue
        parts = path.relative_to(RESULTS_ROOT).parts
        if len(parts) != 4:
            continue
        # results/<sim>/<scenario>/<model>/<timestamp>.md
        if parts[2] == model:
            if is_prototype_result(path) and not include_prototype:
                continue
            yield path


def index_model_cells(model: str, include_prototype: bool) -> Dict[CellKey, Dict[str, Path]]:
    index: Dict[CellKey, Dict[str, Path]] = defaultdict(dict)
    for path in iter_model_files(model, include_prototype):
        rel = path.relative_to(RESULTS_ROOT).parts
        sim, scenario, model_slug, filename = rel
        if model_slug != model:
            continue
        timestamp = filename.replace(".md", "")
        index[CellKey(sim=sim, scenario=scenario)][timestamp] = path
    return index


def choose_cell_path(
    by_timestamp: Dict[str, Path], requested_ts: str, model: str, key: CellKey
) -> Tuple[str, Path]:
    if requested_ts == "latest":
        timestamp = sorted(by_timestamp.keys())[-1]
        return timestamp, by_timestamp[timestamp]
    if requested_ts not in by_timestamp:
        raise RuntimeError(
            f"Missing {model} timestamp {requested_ts} for cell ({key.sim}, {key.scenario})"
        )
    return requested_ts, by_timestamp[requested_ts]


def load_cells(
    model: str, requested_ts: str, keys: Iterable[CellKey], index: Dict[CellKey, Dict[str, Path]]
) -> Dict[CellKey, ResultCell]:
    result: Dict[CellKey, ResultCell] = {}
    for key in keys:
        ts, path = choose_cell_path(index[key], requested_ts, model, key)
        verdict, verdict_line = extract_verdict(path)
        result[key] = ResultCell(
            key=key,
            model=model,
            timestamp=ts,
            path=path,
            verdict=verdict,
            verdict_line=verdict_line,
            score=score_verdict(verdict),
        )
    return result


def simulation_rank(cells: Dict[CellKey, ResultCell]) -> List[Tuple[str, float]]:
    grouped: Dict[str, List[int]] = defaultdict(list)
    for cell in cells.values():
        grouped[cell.key.sim].append(cell.score)
    ranked = [
        (sim, sum(scores) / len(scores))
        for sim, scores in grouped.items()
    ]
    return sorted(ranked, key=lambda item: (-item[1], item[0]))


def scenario_rank(cells: Dict[CellKey, ResultCell]) -> List[Tuple[str, float]]:
    grouped: Dict[str, List[int]] = defaultdict(list)
    for cell in cells.values():
        grouped[cell.key.scenario].append(cell.score)
    ranked = [
        (scenario, sum(scores) / len(scores))
        for scenario, scores in grouped.items()
    ]
    return sorted(ranked, key=lambda item: (item[0]))


def render_report(
    old_model: str,
    new_model: str,
    old_cells: Dict[CellKey, ResultCell],
    new_cells: Dict[CellKey, ResultCell],
) -> str:
    keys = sorted(old_cells.keys(), key=lambda k: (k.sim, k.scenario))
    verdict_text_changes = sum(
        1 for key in keys if old_cells[key].verdict != new_cells[key].verdict
    )
    score_changes = sum(1 for key in keys if old_cells[key].score != new_cells[key].score)

    old_ranked = simulation_rank(old_cells)
    new_ranked = simulation_rank(new_cells)
    old_rank_ix = {sim: idx + 1 for idx, (sim, _) in enumerate(old_ranked)}
    new_rank_ix = {sim: idx + 1 for idx, (sim, _) in enumerate(new_ranked)}
    old_avg_ix = {sim: avg for sim, avg in old_ranked}
    new_avg_ix = {sim: avg for sim, avg in new_ranked}

    old_scenarios = dict(scenario_rank(old_cells))
    new_scenarios = dict(scenario_rank(new_cells))

    lines: List[str] = []
    lines.append(f"# Cross-Model Drift Report: {old_model} vs {new_model}")
    lines.append("")
    lines.append(f"- Baseline model: `{old_model}`")
    lines.append(f"- Comparison model: `{new_model}`")
    lines.append(f"- Cells compared: `{len(keys)}`")
    lines.append(f"- Verdict text changes: `{verdict_text_changes}`")
    lines.append(f"- Score/rank changes: `{score_changes}`")
    lines.append("")
    lines.append("## Simulation Ranking Shift")
    lines.append("")
    lines.append("| Simulation | Old avg score | Old rank | New avg score | New rank | Rank delta |")
    lines.append("|---|---:|---:|---:|---:|---:|")
    for sim, _ in old_ranked:
        lines.append(
            f"| `{sim}` | {old_avg_ix[sim]:.2f} | {old_rank_ix[sim]} | "
            f"{new_avg_ix[sim]:.2f} | {new_rank_ix[sim]} | "
            f"{old_rank_ix[sim] - new_rank_ix[sim]:+d} |"
        )
    lines.append("")
    lines.append("## Per-Cell Drift")
    lines.append("")
    lines.append("| Simulation | Scenario | Old verdict | New verdict | Score delta |")
    lines.append("|---|---|---|---|---:|")
    for key in keys:
        old_cell = old_cells[key]
        new_cell = new_cells[key]
        lines.append(
            f"| `{key.sim}` | `{key.scenario}` | {old_cell.verdict} | "
            f"{new_cell.verdict} | {new_cell.score - old_cell.score:+d} |"
        )
    lines.append("")
    lines.append("## Scenario Aggregates")
    lines.append("")
    lines.append("| Scenario | Old avg score | New avg score | Delta |")
    lines.append("|---|---:|---:|---:|")
    for scenario in sorted(old_scenarios.keys()):
        old_avg = old_scenarios[scenario]
        new_avg = new_scenarios[scenario]
        lines.append(
            f"| `{scenario}` | {old_avg:.2f} | {new_avg:.2f} | {new_avg - old_avg:+.2f} |"
        )
    lines.append("")
    lines.append("## Notes")
    lines.append("")
    lines.append(
        "- This report compares verdict lines from result artifacts; it does not execute protocol harness code."
    )
    lines.append(
        "- Scripted prototype plumbing outputs are excluded unless the comparison tool is run with `--include-prototype`."
    )
    lines.append(
        "- Paths and line numbers for each verdict are in the source result files under `results/<sim>/<scenario>/<model>/<timestamp>.md`."
    )
    return "\n".join(lines) + "\n"


def main() -> int:
    args = parse_args()
    old_index = index_model_cells(args.old_model, args.include_prototype)
    new_index = index_model_cells(args.new_model, args.include_prototype)
    shared_keys = sorted(set(old_index.keys()) & set(new_index.keys()), key=lambda k: (k.sim, k.scenario))
    if not shared_keys:
        raise RuntimeError("No shared (sim, scenario) cells between the two models.")

    old_cells = load_cells(args.old_model, args.old_ts, shared_keys, old_index)
    new_cells = load_cells(args.new_model, args.new_ts, shared_keys, new_index)

    report = render_report(args.old_model, args.new_model, old_cells, new_cells)

    if args.output:
        output_path = Path(args.output).resolve()
    else:
        inferred_ts = sorted({cell.timestamp for cell in new_cells.values()})[-1]
        output_name = f"{args.old_model}_vs_{args.new_model}_{inferred_ts}.md"
        output_path = DEFAULT_OUTPUT_DIR / output_name
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(report)
    print(output_path)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
