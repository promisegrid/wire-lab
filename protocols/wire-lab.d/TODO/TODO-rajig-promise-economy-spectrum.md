# TODO-rajig: Promise-economy mechanism-spectrum simulation owner

## Status

Open. Owns turn-179 promise-economy mechanism-neutrality questions and
coordinates `simulations/SIM-haros-promise-economy-spectrum/`. This TODO does
not choose an economics mechanism, does not make transferable promise tokens part
of the base protocol, and does not make promisebase authoritative. Source:
`DI-pidag`.

## Scope

- Keep the promise-economy spectrum explicit: peer-local social assessment,
  reciprocal promises, capability tokens, permissioned commitments,
  transferability, floating exchange rates, and cryptocurrency-like failure
  modes.
- Use `SIM-haros-promise-economy-spectrum/` as the simulation question home for
  mechanism-neutrality pressure before any protocol field, wire object, or
  group/feed behavior chooses a point on that spectrum.
- Coordinate with TODO-ralud / `SIM-zarud` for conditional release,
  onward-restraint, geofencing, and recursive promise-graph questions.
- Coordinate with TODO-rusap-style promise-accounting simulations where
  peer-local records inform pull, keep, advertise, or refusal decisions, while
  keeping transferable-token questions separate until a later DF explicitly
  accepts them.

## Subtasks

- [ ] rajig.1 Compare social-assessment, reciprocal-promise, capability-token,
  transferable-token, floating-rate, and no-token variants inside `SIM-haros`
  without declaring a base-protocol winner.
- [ ] rajig.2 Identify protocol fields, object shapes, or feed/group behaviors
  that would accidentally bake in fungibility, token balances, universal pricing,
  or central-exchange assumptions.
- [ ] rajig.3 Decide what must remain out of the base PromiseGrid protocol until
  simulations show a mechanism avoids cryptocurrency-like failure modes.
- [ ] rajig.4 Route any surviving mechanism-specific design into successor sims,
  TEs, DRs, or protocol TODOs before closing this owner.

## Routed Elsewhere

- Conditional-release, onward-restraint, geofencing, and recursive promise-graph
  ownership belong to TODO-ralud and `SIM-zarud-conditional-release-geofencing/`.
- Promisebase prior-art adoption and prototype-vs-canon questions belong to
  TODO-kituj, `DR-tumus`, TODO-dozak, and later promisebase review artifacts.
- Dogfood process lessons from the turn-179 doc-vs-code mistake belong to the
  replay and agent-instruction owners; this TODO only owns the PromiseGrid
  design-space question.

## Question Log

- 2026-05-17: Turn 179's promise-economy spectrum is routed here so
  `SIM-haros-promise-economy-spectrum/` has a TODO owner and the design question
  does not live only in a simulation directory or closed replay notes. Source:
  `DI-pidag`.

## Decision Intent Log

The creation and routing decision for this TODO is recorded as `DI-pidag` in
`TODO-juhub-turns-149-208-chronological-rewalk.md`. Future locked decisions
local to the promise-economy mechanism spectrum should be appended here.
