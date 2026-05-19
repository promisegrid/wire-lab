# Executable Evidence Check: Top 3 Cells

- Run timestamp UTC: `2026-05-19 03:02:31 UTC`
- Scope: top 3 cells by latest numeric verdict ranking in `openai-gpt-5.3-codex-xhigh` results.
- Method: run executable artifact checks available in-repo (JSON parse + structural spec assertions).

## `SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` × `journalism-source-provenance`

- Result under test: `results/SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/journalism-source-provenance/openai-gpt-5.3-codex-xhigh/20260519-025634.md`
- Executable harness availability: `0` runnable code/build files found in sim tree.
- Harness candidates: none in this simulation tree; packet/runtime protocol replay is not available here.
- Executable artifact checks passed: `5/5`
- `manifest-protocol-grid-envelope`: `PASS`
- `spec-shape-sig_pcid_payload`: `PASS`
- `spec-opaque-policy`: `PASS`
- `scenario-overarching-goal-checks`: `PASS`
- `result-verdict-line`: `PASS`

## `SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` × `promisebase-reference-naming-promisebase-custom-syntax-migration`

- Result under test: `results/SIM-riliz-grid-envelope-enc-dag-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/promisebase-reference-naming-promisebase-custom-syntax-migration/openai-gpt-5.3-codex-xhigh/20260519-030134.md`
- Executable harness availability: `0` runnable code/build files found in sim tree.
- Harness candidates: none in this simulation tree; packet/runtime protocol replay is not available here.
- Executable artifact checks passed: `5/5`
- `manifest-protocol-grid-envelope`: `PASS`
- `spec-shape-sig_pcid_payload`: `PASS`
- `spec-opaque-policy`: `PASS`
- `scenario-overarching-goal-checks`: `PASS`
- `result-verdict-line`: `PASS`

## `SIM-vamaz-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload` × `promisebase-reference-naming-promisebase-custom-syntax-migration`

- Result under test: `results/SIM-vamaz-grid-envelope-enc-cbor-unknown-opaque-sig-mandatory-sig-pcid-payload/promisebase-reference-naming-promisebase-custom-syntax-migration/openai-gpt-5.3-codex-xhigh/20260519-025634.md`
- Executable harness availability: `0` runnable code/build files found in sim tree.
- Harness candidates: none in this simulation tree; packet/runtime protocol replay is not available here.
- Executable artifact checks passed: `5/5`
- `manifest-protocol-grid-envelope`: `PASS`
- `spec-shape-sig_pcid_payload`: `PASS`
- `spec-opaque-policy`: `PASS`
- `scenario-overarching-goal-checks`: `PASS`
- `result-verdict-line`: `PASS`

## Conclusion

- No runnable protocol harness was available for these top 3 cells in their simulation trees.
- Executable artifact-level checks passed for all three cells, confirming manifest/spec/scenario/result structural consistency.
- To get runtime evidence for these cells, a dedicated executable grid-envelope harness would need to be added first.
