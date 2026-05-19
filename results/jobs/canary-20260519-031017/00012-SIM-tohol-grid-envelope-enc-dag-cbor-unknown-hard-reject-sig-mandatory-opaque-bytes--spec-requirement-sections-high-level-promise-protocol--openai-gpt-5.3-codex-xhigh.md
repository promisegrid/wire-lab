# LLM Matrix Cell Job

## Cell

- Run group ID: `canary-20260519-031017`
- Simulation ID: `SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes`
- Scenario ID: `spec-requirement-sections-high-level-promise-protocol`
- Model ID: `openai-gpt-5.3-codex-xhigh`
- Intended result path: `results/SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes/spec-requirement-sections-high-level-promise-protocol/openai-gpt-5.3-codex-xhigh/20260519-031715.md`

## Required Source Inputs

Read only source/design inputs before producing the verdict:

- `simulations/SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes/README.md`
- `simulations/SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes/QUESTION.md` if present
- local draft specs under `simulations/SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes/` if present
- `scenarios/spec-requirement-sections-high-level-promise-protocol/spec-requirement-sections-high-level-promise-protocol.md`
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

`results/SIM-tohol-grid-envelope-enc-dag-cbor-unknown-hard-reject-sig-mandatory-opaque-bytes/spec-requirement-sections-high-level-promise-protocol/openai-gpt-5.3-codex-xhigh/20260519-031715.md`

The result must follow the section contract in `results/RUN-PROTOCOL.md` and
must include:

- `- Run mode: llm-doc-eval-blind`
- a line starting with `Evidence verdict:`
- an explicit `Authority Boundary` section.
