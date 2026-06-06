# Capability-token exchange

## Scenario ID

promise-economy-capability-token-exchange

## Source / Provenance

- Source type: new harness scenario
- Source path: `simulations/SIM-haros-promise-economy-spectrum/SCENARIOS.md`;
  `protocols/wire-lab.d/TODO/TODO-rajig-promise-economy-spectrum.md`;
  `scenarios/promise-economy-spectrum-permissioned-capability-token/promise-economy-spectrum-permissioned-capability-token.md`;
  `scenarios/promise-economy-spectrum-transferable-promise-token/promise-economy-spectrum-transferable-promise-token.md`;
  `scenarios/promise-economy-spectrum-floating-exchange-rates/promise-economy-spectrum-floating-exchange-rates.md`
- Source simulation: `SIM-haros-promise-economy-spectrum/`
- Source row/title: Combined capability-token access, bearer-token transfer,
  non-transferable redemption, revocation, and peer-local floating exchange rates
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`;
  `DI-pidag`; `DI-hosuk`; `DI-sirus`; `DI-vorus`; `DI-bibom`

## Purpose

Test whether a candidate PromiseGrid design can model security capability tokens
as promises by resource-controlling agents without turning those promises into a
central permission system, central exchange, universal currency, or global price
oracle.

## Setup

Alice controls storage, compute, or data resources and issues tokens that promise
access to those resources when the token is redeemed. Some of Alice's tokens are
non-transferable and can only be used by the original issuee. Other tokens are
bearer tokens: whoever presents the token can ask Alice to keep the access
promise, subject to Alice's own revocation and local trust rules.

Bob, Carol, and Dave each hold different mixes of bearer tokens,
non-transferable access tokens, and local observations about prior promise
keeping or breaking. They may trade bearer tokens with each other, trade bearer
tokens for non-transferable access tokens, or reject offers that do not fit their
needs. Mallory offers tokens from an issuer with poor make/break history or tries
to circulate revoked, stale, or misleading tokens.

There is no central exchange and no global exchange rate. Each agent sees only
the offers made by its peers and values each issuer's bearer tokens according to
local needs, scarcity, perceived resource usefulness, revocation risk, and the
issuer's promise keep/break history.

## Stimulus

Run the candidate simulation against this source test: Alice issues bearer and
non-transferable capability-token promises, Bob redeems one token for access,
Carol trades bearer tokens for a non-transferable access token, Dave quotes a
different exchange rate based on his own trust history, Alice revokes a token
class after a broken promise, and Mallory attempts to trade a token whose
redemption promise is stale or locally distrusted.

## Expected Pressure

The simulation should preserve the distinction between a token as evidence of an
agent's promise and a token as global authority. It should show how redemption
fulfills promises, how revocation or refusal becomes local broken-promise
evidence, how bearer and non-transferable tokens differ, how peer-local exchange
offers create floating rates, and how agents accept or reject offers based on
their own needs and trust judgments.

## POC8 Evidence So Far

`implementations/poc8-autonomous-promise-economy/` is useful evidence for this
scenario because it keeps the whole promise-economy protocol under one pCID and
uses signed CBOR `grid([42(pCID), payload, proof])` messages over framed TCP.
Its payload `kind` values model need advertisements, offer promises,
counter-promises, acceptance/refusal decisions, token issue promises, token
redemption promises, outcome observations, collateral/stake promises, and
peer-local exchange-rate quotes as variants of that one protocol rather than as
separate protocols.

POC8 also improves the agency model over POC7: Alice advertises needs as
Alice-owned promises instead of commands, Bob and Carol locally choose whether
to offer, counter, accept, issue, redeem, or refuse, Alice can trade a bearer
stake promise for non-transferable compute access, and Dave can refuse Mallory's
later stale-token circulation after Dave's own local trust decreases.

## Additional Pressure After POC8

The next pressure is discovery. POC8 still starts from a bounded deterministic
world: the role plans, known peers, topology, stale-token history, routes, and
completion marker are mostly fixed by the harness. A stronger design must show
how agents learn who might be useful peers, what those peers promise, and how
much risk each relationship can carry without introducing a central directory,
global trust score, service registry, exchange authority, or permission system.

Discovery should start from existing relationships and low-risk trials. Alice
may begin with configured seed peers, prior trusted peers, local or physical
pairing, imported references, or introductions from a peer she already knows.
Bob may tell Alice that Carol promises storage, but that referral is Bob's
promise and Alice-local evidence; it is not transitive authority. If Bob carries
Carol's signed offer, Alice still decides locally whether Carol is trusted enough
for this specific storage, compute, relay, or exchange risk.

Transport behavior is evidence, not authority. Alice's local
promise-accounting state for Bob may decide whether Alice opens, keeps, retries,
or uses a TCP connection to Bob, and what data Alice is willing to send over it.
TCP delivery, malformed bytes, disconnects, timeouts, retries, and latency then
become Alice-local observations about Bob's transport promises. An open TCP
connection never means Alice trusts Bob, and Alice may trust Carol while
reaching Carol only through a relay promise from Bob or Dave.

Storage and compute access tokens should be tested against the POC13 model. A
token can be a promise that Alice will store, retain, serve, or compute under
specific terms, but the token is not a global permission. For compute, a
non-transferable token may promise one call to a CAS-stored `function_cid` with
explicit input/context objects; a bearer token may be traded as personal
currency only because each holder locally values Alice's future promise.

## Scenario-Specific Evaluation Questions

- Which agent issues each token, what resource does that agent control, and what
  exactly does the token issuer promise to do on redemption?
- How does the design distinguish bearer tokens from non-transferable tokens
  without relying on a central ownership registry?
- What local evidence is recorded when a token is redeemed, refused, revoked, or
  presented by an unexpected holder?
- How can Bob, Carol, and Dave trade bearer tokens without creating a hidden
  central exchange, universal price oracle, or global currency?
- How does an agent decide whether to exchange bearer tokens for a
  non-transferable access token when the resource value and issuer trust are both
  local and time-varying?
- How do broken promises, revocations, stale offers, and Mallory's misleading
  trades change later exchange rates and access decisions for each observing
  agent?
- How does Alice discover that Bob, Carol, or Dave might make useful storage,
  compute, relay, exchange, or introduction promises without consulting a central
  directory?
- Which discovery statements are the speaker's own promises, which are signed
  offers from another agent, and which are merely local observations?
- What low-risk probe promises let Alice build enough local trust before sending
  sensitive data, relying on returned computation, or accepting bearer-token
  exchange risk?
- How do promise relationships decide whether to open, keep, retry, or avoid TCP
  connections, and how does TCP behavior feed back into local promise accounting?
- When a bearer token buys access to storage or compute, who promises retention,
  serving, execution, result delivery, and result evidence?
- How do local exchange rates change when a compute result was cached from a pure
  `function_cid` call versus produced from impure context objects?
