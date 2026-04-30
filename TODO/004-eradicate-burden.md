# TODO 004 - Eradicate the word `burden`

Replace every occurrence of the word `burden` (and the assertion-sense uses of `claim`) in the wire-lab prose with `assertion`. The word `burden` was acting as a stand-in for "what a promise asserts" but its English connotations of obligation and heavy-thing-to-carry leak imposition-flavor into a vocabulary where promises are autonomous assertions of state. The word `assertion` is already used in this sense in `specs/harness-spec-draft.md` (e.g., "signature frame asserts authorship") and matches the canonical Promise-Theory speech-act of stating something as so.

## Subtasks

- [x] 004.1 Replace every `burden` occurrence in `specs/harness-spec-draft.md` with `assertion` (15 places, including the `Promise` struct field name, the §2.4 heading, and adjective forms like `transport burden`).
- [x] 004.2 Replace every `burden` occurrence in `docs/thought-experiments/TE-*.md` files (3 places, one per affected TE).
- [x] 004.3 Replace the assertion-sense `claim` in `specs/harness-spec-draft.md:479` ("a causal claim …") with `assertion`. Leave the verb-sense `claiming` in `specs/harness-spec-draft.md:174` ("two handlers claiming the same port") unchanged because it is the demand-sense, not the assertion-sense.
- [x] 004.4 Preserve the superseded `proposals/pending/ppx-dr-001-bootstrap/contest-20260429-033208-steve-traugott.md` artifact unmodified, per `DI-003-20260429-162212`'s preserve-history clause.

## Decision Intent Log

ID: DI-004-20260429-163000
Date: 2026-04-29 16:30:00
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: The wire-lab prose vocabulary uses `assertion` (and the verb form `asserts`) for what a promise asserts about state. The word `burden` is eradicated from the existing prose, including the `Promise` struct field name and trust-ledger phrases like `per-burden-type` (which becomes `per-assertion-type`). The verb-sense use of `claim` (`claiming the same port`) is unchanged; only the assertion-sense `claim` (`a causal claim`) is replaced.
Intent: Stop importing obligation/heavy-thing-to-carry connotations into a vocabulary where promises are autonomous assertions of state. Match the existing assert/asserts usage already present in the spec (e.g., "signature frame asserts authorship") and align with Promise-Theory speech-act terminology. Avoid `claim` because its primary English sense (a demand on resources) drags imposition-flavor and because compound forms like `claim CID` would collide visually with `pCID`.
Constraints: This DI is a vocabulary-level cleanup, not a structural change to the promise-stack model in §1.1, the trust-ledger shape in §2.1, or any other design. The replacement is mechanical except for the §2.4 heading rename (`Trust is per-burden-type` → `Trust is per-assertion-type`) and the `Promise` struct field name (`burden:` → `assertion:`). Future references to existing `burden_type`-style identifiers (e.g., `map[burden_type]float`) become `map[assertion_type]float`.
Affects: `specs/harness-spec-draft.md`; `docs/thought-experiments/TE-20260427-180000-promise-stack-ordering.md`; `docs/thought-experiments/TE-20260428-080000-harness-spec-change-walks-through-unified-flow.md`; `docs/thought-experiments/TE-20260428-094500-should-this-design-become-promisegrid-readme.md`.
Linked DR: none (chat-directed cleanup; no Decision Record needed per Steve's instruction).

## Notes

This correction was approved in chat after pressure-testing six candidate words against the existing `assert/asserts` usage in `specs/harness-spec-draft.md` and against composability with `pCID`. The chosen word `assertion` extends an existing convention rather than introducing a parallel one.
