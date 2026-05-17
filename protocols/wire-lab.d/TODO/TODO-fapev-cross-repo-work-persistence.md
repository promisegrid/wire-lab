# TODO-fapev: Cross-repo work persistence before PAT-gated pushes

## Status

Open. This TODO owns the turn-183 process lesson that local-only cross-repo
commits prepared before PAT-granted push need a session-survivable persistence
path. It is a harness/process owner, not a PromiseGrid protocol-design question
home. Source: `DI-kegar`.

## Scope

- Define when a local-only cross-repo commit, patch, or branch needs persistence
  before handoff, session reset, or PAT grant.
- Choose the persistence mechanism: private patch file under an approved
  private path, private mirror push, bundle file, or another repo-safe method.
- Specify cleanup, redaction, and lifecycle expectations for patch, bundle, or
  mirror artifacts.
- Route the final rule into the appropriate agent instruction file after the
  persistence mechanism is chosen.

## Subtasks

- [ ] fapev.1 Decide the persistence mechanism for local-only cross-repo commits
  pending authentication or push authorization.
- [ ] fapev.2 Define path, privacy, and cleanup rules for patch, bundle, or
  mirror artifacts.
- [ ] fapev.3 Add the resulting process rule to the appropriate agent
  instruction file.
- [ ] fapev.4 Cross-reference the rule from future cross-repo PAT, push, and
  handoff workflows.

## Question Log

- 2026-05-17: Turn 183 left the promisebase randStream fix as a local-only
  commit until the later PAT grant and push. Future cross-repo work needs a
  durable persistence rule so the work survives reset before authentication is
  available. Source: `DI-kegar`.

## Decision Intent Log

The creation and routing decision for this TODO is recorded as `DI-kegar` in
`TODO-juhub-turns-149-208-chronological-rewalk.md`. Future locked decisions
local to cross-repo work persistence should be appended here.
