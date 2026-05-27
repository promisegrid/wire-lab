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
- [x] rajig.7 Add `implementations/poc7-capability-token-exchange/` as
  executable evidence for the combined capability-token exchange scenario. Done
  under `DI-tugih`.
- [x] rajig.8 Deepen POC7 from a shallow JSON sketch into a protocol POC with
  CBOR grid envelopes, protocol-owned payloads, signed token bytes,
  storage/compute redemption work, peer-local exchange-state mutation, and the
  Carol access-token bug fixed. Done under `DI-fibok`.
- [x] rajig.9 Refactor POC7 away from HTTP wrapper transport to framed TCP and
  reframe app message names around promises, reciprocal promises, local
  evidence, and local trust judgments. Done under `DI-tanat`.

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

ID: DI-tugih
Date: 2026-05-26 15:58:15
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement `implementations/poc7-capability-token-exchange/` as a
five-container proof of concept for
`scenarios/promise-economy-capability-token-exchange/`. Each container runs one
local kernel role, one app-level relay role, and role-specific token/resource
apps. The demo uses fixed ticks, indivisible signed bearer and
non-transferable tokens, issuer-local redemption/revocation state, peer-local
exchange offers, and local trust updates.
Intent: Produce executable evidence for the hardest promise-economy pressure
without choosing a base economics model: resource-controlling agents issue
tokens as promises, token redemption keeps or breaks those promises, bearer
tokens can circulate as personal currency, non-transferable tokens remain
relationship-scoped, and exchange rates remain local judgments instead of a
central market fact.
Constraints: Keep POC7 as evidence only, not a final SDK, kernel API, token
standard, exchange protocol, trust API, or economics model. Do not introduce a
central exchange, shared token-status ledger, global price oracle, global trust
authority, or permission authority. Reuse POC5-style app/kernel/relay
discipline but do not mutate POC5.
Affects: `implementations/poc7-capability-token-exchange/**`;
`implementations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-rajig-promise-economy-spectrum.md`.

ID: DI-tanat
Date: 2026-05-26 23:03:49
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Refactor `implementations/poc7-capability-token-exchange/` so POC7
uses framed TCP instead of HTTP wrapper transport, and rename/reframe its app
message kinds so the visible protocol vocabulary is promise-first rather than
RPC-like. The locked message vocabulary is
`resource_promise_request_v1`, `promise_revocation_notice_v1`,
`promise_presented_for_fulfillment_v1`, `promise_received_v1`,
`exchange_offer_promise_v1`, `reciprocal_exchange_promise_v1`,
`promise_evidence_request_v1`, `held_promise_fulfillment_request_v1`, and
`held_reciprocal_exchange_request_v1`.
Intent: POC7 already models signed capability tokens as issuer promises and
keeps trust local, but its HTTP wrapper and command-like names made the
implementation read too much like a conventional RPC service. This follow-on
keeps the executable evidence but makes transport continuity match POC2 through
POC5, and makes each visible message read as a promise request, promise
presentation, promise receipt, reciprocal exchange promise, or local evidence
request instead of as global permission, authorization, or command semantics.
Constraints: Keep POC7 evidence-only and do not define a final PromiseGrid
token, economics, storage, compute, kernel, transport, or app API. Keep trust,
pricing, revocation, and redemption local. Do not add a central exchange,
global trust authority, permission authority, shared token-status ledger, HTTP
server, or non-standard daemon dependency. Preserve the signed CBOR
`grid([42(pCID), payload, proof])` message envelope; TCP framing is transport
plumbing and not the semantic protocol.
Affects: `implementations/poc7-capability-token-exchange/**`;
`implementations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-rajig-promise-economy-spectrum.md`.

ID: DI-fibok
Date: 2026-05-26 22:43:45
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Deepen `implementations/poc7-capability-token-exchange/` after the
initial POC7 commit by replacing JSON-as-protocol with CBOR
`grid([42(pCID), payload, proof])` envelopes, making storage and compute
redemption perform real app payload work, using Ed25519 signatures over
canonical token bytes and envelope signable views, making exchange offers mutate
peer-local wallet state, and fixing the mistaken Alice-self-issued access token
so the trade issues the promised non-transferable access token to Carol.
Intent: POC7 should remain bounded evidence, but it should exercise the actual
PromiseGrid protocol pressure points instead of only naming them: pCID-selected
CBOR envelopes, exact signed bytes, resource work at redemption time, local
economics without a central exchange, and transaction semantics that match the
Alice/Bob/Carol/Dave/Mallory story.
Constraints: Keep POC7 evidence-only and do not declare a final token,
economics, storage, compute, kernel, or exchange API. Keep all trust,
revocation, pricing, and redemption judgment local to the observing agent or
issuer. Do not introduce a central exchange, global trust authority, shared
token-status ledger, permission authority, or non-standard daemon dependency.
Prefer existing POC4/POC5 CBOR-grid patterns and POC6 DAG-CBOR/IPLD evidence
where practical, but keep the POC cheap and runnable in the existing five
containers.
Affects: `implementations/poc7-capability-token-exchange/**`;
`implementations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-rajig-promise-economy-spectrum.md`.
