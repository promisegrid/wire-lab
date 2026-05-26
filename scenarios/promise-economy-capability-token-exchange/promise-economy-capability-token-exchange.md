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
  `DI-pidag`; `DI-hosuk`

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
