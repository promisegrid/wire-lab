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
- Compare RFC-1005's content-addressable test-driven fabric pattern -- test tree
  CID, executable tree CID, arguments, and cache-on-pass semantics -- as prior
  art without making it the base PromiseGrid economics model. Source:
  `DI-nulak`.

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
- [ ] rajig.5 Map RFC-1005's test-tree / executable-tree / cache-on-pass
  vocabulary against PromiseGrid promise-economy terms before adopting or
  rejecting any of it.
- [x] rajig.6 Add a root scenario for capability-token access promises,
  redemption, revocation, bearer/non-transferable token exchange, and peer-local
  floating exchange rates without choosing a base PromiseGrid economics model.
  Done under `DI-hosuk`.

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
- 2026-05-17: Turn 184 routes RFC-1005 Option 2's test-driven fabric as
  promise-economy prior art into this TODO and `SIM-haros`, without treating it
  as settled PromiseGrid economics. Source: `DI-nulak`.

## Decision Intent Log

The creation and routing decision for this TODO is recorded as `DI-pidag` in
`TODO-juhub-turns-149-208-chronological-rewalk.md`. Future locked decisions
local to the promise-economy mechanism spectrum should be appended here.

ID: DI-hosuk
Date: 2026-05-26 13:28:21
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Add
`scenarios/promise-economy-capability-token-exchange/promise-economy-capability-token-exchange.md`
as a root scenario combining capability-token access promises, redemption,
revocation, bearer-token transfer, non-transferable access tokens, peer-local
exchange offers, and floating trust-based exchange rates.
Intent: Existing root scenarios split permissioned capability tokens,
transferable promise tokens, and floating exchange rates into separate pressure
cases. The repo also needs one combined scenario where agents issue security
capability tokens for resources they control, redeem tokens to fulfill promises,
observe broken promises or revocations, trade bearer tokens as personal
currencies, exchange bearer tokens for non-transferable access tokens, and value
each issuer's token through local trust and keep/break history without a central
exchange.
Constraints: Keep this as root scenario evidence pressure only. Do not choose a
base PromiseGrid economics model, token standard, central exchange, global price
oracle, or permission authority. Do not edit simulations, GA config, rubric, or
`DEV-GUIDE-RESOURCES.md` in this step. Preserve existing uncommitted atproto
scenario changes.
Affects:
`scenarios/promise-economy-capability-token-exchange/promise-economy-capability-token-exchange.md`;
`protocols/wire-lab.d/TODO/TODO-rajig-promise-economy-spectrum.md`.
