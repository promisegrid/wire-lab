# TODO-nasat: Cross-repo build hazard capture

## Status

Open. This TODO owns the turn-186 process lesson that operational build and
toolchain diagnostics discovered during cross-repo work must be captured in a
discoverable owner artifact instead of living only in chat narration. It is a
harness/process owner, not a PromiseGrid protocol-design question home. Source:
`DI-zagus`.

## Scope

- Capture cross-repo build hazards that affect wire-lab / promisebase
  coordination but do not themselves decide PromiseGrid design.
- Decide whether each hazard belongs in a source repo troubleshooting document,
  a wire-lab coordination note, an owner TODO, or an agent instruction rule.
- Preserve enough detail for future agents to recognize and avoid the failure
  without rereading the original chat log.
- Do not create promisebase commits from this replay pass.

## Current Hazard Notes

- **GOTOOLCHAIN auto-bump.** During turn 186, the assistant observed a
  `GOTOOLCHAIN=auto` / newer-Go lazy `go.mod` write hazard: running tests with a
  newer installed Go toolchain can inject a higher `toolchain` line into
  `go.mod`. The assistant checked out the drift before pushing the promisebase
  randStream fix so promisebase main kept Steve's explicit Go version/toolchain
  setting from the pre-fix commit. Source: `DI-zagus`.
- **Promisebase dependency drift cluster.** Related hazards from nearby turns
  include the Go 1.20 `math/rand.Seed` behavior change fixed in turn 183 and the
  Docker SDK / FUSE / server partial-rot evidence routed to TODO-kituj by
  `DI-nulak`.

## Subtasks

- [ ] nasat.1 Decide the durable home for cross-repo build/toolchain hazard notes
  that are discovered while working from wire-lab.
- [ ] nasat.2 Move or duplicate the GOTOOLCHAIN auto-bump note into the chosen
  durable troubleshooting home.
- [ ] nasat.3 Define whether future hazards should be captured immediately in
  AGENTS-ppx, a wire-lab troubleshooting note, or the affected repo's own docs.
- [ ] nasat.4 Cross-reference TODO-kituj and promisebase adoption work when a
  build hazard affects promisebase / pitbase readiness.

## Question Log

- 2026-05-17: Turn 186 surfaced the GOTOOLCHAIN auto-bump diagnostic during
  cross-repo promisebase work. The diagnostic is useful operational knowledge,
  but the replay should not edit promisebase or invent a permanent
  troubleshooting-file structure in this pass. Source: `DI-zagus`.

## Decision Intent Log

The creation and routing decision for this TODO is recorded as `DI-zagus` in
`TODO-juhub-turns-149-208-chronological-rewalk.md`. Future locked decisions
local to cross-repo build-hazard capture should be appended here.
