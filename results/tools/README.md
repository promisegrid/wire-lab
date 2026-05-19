# Results Tooling

This directory contains matrix-run preflight and execution helpers.

## Generate Manifest

Generate full matrix for one model:

```bash
python3 results/tools/generate_matrix_manifest.py \
  --models openai-gpt-5.3-codex-xhigh
```

Generate deterministic canary manifest (30 cells):

```bash
python3 results/tools/generate_matrix_manifest.py \
  --models openai-gpt-5.3-codex-xhigh \
  --shuffle-seed 42 \
  --limit-cells 30
```

## Generate LLM Jobs

Generate prompt files from a manifest:

```bash
python3 results/tools/generate_llm_jobs.py \
  --manifest results/manifests/<manifest>.csv
```

The generator does not write final result verdicts. It writes prompt files under
`results/jobs/<run-group-id>/` for Codex or another LLM runner to evaluate.
Source: `DI-moduf`.

Generate only part of a manifest:

```bash
python3 results/tools/generate_llm_jobs.py \
  --manifest results/manifests/<manifest>.csv \
  --start-index 500
```

## Validate Batch

Validate all results for a model/timestamp and enforce matrix links:

```bash
python3 results/tools/validate_results.py \
  --model openai-gpt-5.3-codex-xhigh \
  --timestamp <YYYYMMDD-HHMMSS> \
  --strict-matrix
```

Prototype plumbing outputs are rejected by default. Use `--allow-prototype` only
when deliberately checking preserved scripted-output fixtures. Source:
`DI-moduf`.

## Compare Models

```bash
python3 results/comparisons/compare_model_results.py \
  --old-model openai-gpt-5.5-xhigh \
  --new-model openai-gpt-5.3-codex-xhigh
```

Comparison excludes scripted prototype outputs by default; pass
`--include-prototype` only for tooling audits. Source: `DI-moduf`.
