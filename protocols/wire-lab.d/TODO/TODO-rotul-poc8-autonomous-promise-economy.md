# TODO-rotul: POC8 autonomous promise economy

## Status

Implemented. Owns `implementations/poc8-autonomous-promise-economy/`, a successor to
POC7 that keeps the signed CBOR grid envelope and promise-token evidence while
moving the economy from an Alice-scripted transaction path to autonomous
peer-local needs, offers, bargaining, collateral/stake promises, repeated
rounds, local trust, and peer-local exchange rates. Source: `DI-sirus`.

## Scope

- Preserve POC7 as historical executable evidence; build POC8 as a new sibling
  implementation.
- Use one protocol pCID for the POC8 promise-economy protocol. Message kinds are
  protocol-owned payload variants under that single pCID, not separate pCIDs.
- Use signed CBOR `grid([42(pCID), payload, proof])` envelopes over framed TCP.
- Keep all trust, valuation, offer acceptance, counteroffers, collateral/stake
  interpretation, and exchange-rate quotes local to the observing agent.
- Keep tokens as promises by their issuers. A bearer token is a transferable
  issuer promise; a non-transferable token is scoped to the original issuee.
- Do not add a central exchange, central scheduler after startup, global trust
  score, global price oracle, shared token-status ledger, permission authority,
  authorization authority, or contract-enforcement framing.

## Subtasks

- [x] rotul.1 Create the POC8 sibling implementation and preserve POC7 as prior evidence.
- [x] rotul.2 Record one pCID-selected POC8 payload vocabulary with `need_advertisement`, `offer_promise`, `counter_promise`, `acceptance_promise`, `token_issue_promise`, `token_redemption_promise`, and `outcome_observation` variants.
- [x] rotul.3 Implement autonomous peer-local decisions for advertising needs, making offers, counteroffering, accepting, refusing, redeeming tokens, transferring bearer tokens, and quoting exchange rates.
- [x] rotul.4 Add full-economy pressure: collateral/stake promises, scarcity, repeated rounds, local floating exchange rates, bearer-for-nontransferable exchange, and stale-token trust decay.
- [x] rotul.5 Keep all inter-peer and app/kernel messages on signed CBOR grid envelopes over framed TCP.
- [x] rotul.6 Add deterministic tests for payload encoding, proof verification, local valuation, counteroffers, refusals, trust updates, transferability, revocation, and collateral/stake behavior.
- [x] rotul.7 Update implementation and guide resource docs to describe POC8 as evidence, not final API.

## Decision Intent Log

ID: DI-sirus
Date: 2026-05-31 16:00:28
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement `implementations/poc8-autonomous-promise-economy/` as a new
successor to POC7. POC8 uses one protocol pCID for the entire promise-economy
protocol, with payload `kind` values selecting protocol-defined variants. It
keeps signed CBOR `grid([42(pCID), payload, proof])` envelopes over framed TCP,
and it adds autonomous local need advertising, offer promises, counteroffers,
acceptance/refusal decisions, token issuance/redemption, collateral/stake
promises, repeated local rounds, and peer-local exchange-rate behavior.
Intent: POC7 proved useful byte-level and token-level mechanics, but Alice's
harness script still created most opportunities, routes, and timing. POC8 should
exercise the next pressure case: agents make their own promises, advertise their
own needs or offers, decide from local trust and utility, negotiate reciprocal
terms, and suffer local relationship consequences when promises are broken,
without pretending that one peer can command another or that a central authority
can define trust, prices, permission, authorization, or conformance.
Constraints: POC8 is executable design evidence, not a final PromiseGrid API,
token standard, economics model, kernel API, storage API, compute API, transport
standard, or SDK. Preserve POC7 unchanged. Do not add a central exchange, global
price oracle, global trust authority, shared token-status ledger, non-voluntary
cooperation, permission authority, authorization authority, or contract-
enforcement framing. Runtime state remains in-memory plus stdout logs and the
existing Docker completion marker pattern.
Affects: `implementations/poc8-autonomous-promise-economy/**`;
`implementations/README.md`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO.md`;
`protocols/wire-lab.d/TODO/TODO-rotul-poc8-autonomous-promise-economy.md`.
