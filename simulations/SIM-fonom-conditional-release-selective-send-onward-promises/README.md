# SIM-fonom-conditional-release-selective-send-onward-promises

This simulation is the Promise-Theory-first successor to
`SIM-gibut-conditional-release-group-session-local`. It keeps the same design
pressure—conditional release, onward restraint, and geofencing—but reframes the
mechanism as selective sending plus explicit onward-handling promises between
autonomous agents rather than as session-local policy. Source: `DI-tavaz`.

## Design Under Test

Alice does not rely on a session artifact to impose handling restrictions on
Bob. Instead:

- Alice decides locally whether to send based on Bob's prior promise history and
  Bob's current handling promise.
- Bob may promise onward restraint, storage discipline, computation limits, or
  later deletion, but Bob remains autonomous.
- Alice records whether Bob appears to keep or break those promises and updates
  her local trust accordingly.
- Carol or later peers receive data only when the preceding peer voluntarily
  sends it under a new relationship-specific promise.

## Why this differs from `gibut`

`gibut` centered the release conditions inside group/session semantics, which
made storage, routing, and audit observations ambiguous and drifted toward
policy framing. `fonom` keeps the question at the agent layer: who promises
what, what evidence later shows, and what remains local trust judgment.

## Boundaries

This simulation does not claim that transport, session, or storage layers can
enforce Bob's behavior. It only asks whether conditional release can be modeled
cleanly as selective sending plus explicit promises and local trust updates.
