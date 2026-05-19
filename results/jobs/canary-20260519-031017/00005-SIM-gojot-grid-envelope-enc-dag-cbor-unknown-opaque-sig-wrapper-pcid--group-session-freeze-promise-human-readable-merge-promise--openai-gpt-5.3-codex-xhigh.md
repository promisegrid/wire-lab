# LLM Matrix Cell Job

## Cell

- Run group ID: `canary-20260519-031017`
- Simulation ID: `SIM-gojot-grid-envelope-enc-dag-cbor-unknown-opaque-sig-wrapper-pcid`
- Scenario ID: `group-session-freeze-promise-human-readable-merge-promise`
- Model ID: `openai-gpt-5.3-codex-xhigh`
- Intended result path: `results/SIM-gojot-grid-envelope-enc-dag-cbor-unknown-opaque-sig-wrapper-pcid/group-session-freeze-promise-human-readable-merge-promise/openai-gpt-5.3-codex-xhigh/20260519-031715.md`

## Required Source Inputs

Read only source/design inputs before producing the verdict:

- `simulations/SIM-gojot-grid-envelope-enc-dag-cbor-unknown-opaque-sig-wrapper-pcid/README.md`
- `simulations/SIM-gojot-grid-envelope-enc-dag-cbor-unknown-opaque-sig-wrapper-pcid/QUESTION.md` if present
- local draft specs under `simulations/SIM-gojot-grid-envelope-enc-dag-cbor-unknown-opaque-sig-wrapper-pcid/` if present
- `scenarios/group-session-freeze-promise-human-readable-merge-promise/group-session-freeze-promise-human-readable-merge-promise.md`
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

`results/SIM-gojot-grid-envelope-enc-dag-cbor-unknown-opaque-sig-wrapper-pcid/group-session-freeze-promise-human-readable-merge-promise/openai-gpt-5.3-codex-xhigh/20260519-031715.md`

The result must follow the section contract in `results/RUN-PROTOCOL.md` and
must include:

- `- Run mode: llm-doc-eval-blind`
- a line starting with `Evidence verdict:`
- an explicit `Authority Boundary` section.
