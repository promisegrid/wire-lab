# poc8-autonomous-promise-economy

`poc8-autonomous-promise-economy` is executable POC evidence for
`scenarios/promise-economy-capability-token-exchange/` and TODO-rotul. It is a
successor to POC7, not an in-place rewrite of POC7. Source: `DI-sirus`.

The demo uses five containers in a ring:

```text
Alice -- Bob -- Carol -- Dave -- Mallory -- Alice
```

Each container runs one local kernel boundary plus local promise-economy app
roles inside a single bounded process. Peers use length-framed TCP. Each app
message is a signed CBOR `grid(...)` tag (`0x67726964`, decimal `1735551332`)
wrapping `[42(pCID), payload, proof]`.

POC8 uses **one pCID for the whole promise-economy protocol**. The payload field
`kind` selects protocol-defined variants under that one pCID:
`need_advertisement`, `offer_promise`, `counter_promise`,
`acceptance_promise`, `token_issue_promise`, `token_redemption_promise`,
`outcome_observation`, and `exchange_rate_quote`. These are payload kinds, not
separate pCIDs. Source: `DI-sirus`.

## What Takes Place

- **Alice advertises storage and compute needs.** A need advertisement is Alice's
  own promise that the need is current, that Alice will not send data before a
  voluntary acceptance, and that Alice will record keep/break evidence locally.
  It asks for reciprocal promises; it does not command Bob or Carol.
- **Bob offers storage, Alice counters, Bob accepts.** Bob locally scores Alice's
  storage need and offers storage terms. Alice counters the price and stake.
  Bob voluntarily accepts and issues separate non-transferable store/read tokens
  plus a bearer stake token.
- **Alice redeems Bob's storage tokens.** Alice presents the store token with a
  key/value pair, then presents the read token to get the value back. Bob's
  storage app does real deterministic work.
- **Carol offers compute and issues compute access.** Carol locally scores
  Alice's compute need and issues non-transferable compute access after Alice's
  acceptance. Alice redeems the compute token and receives Fibonacci 10 = `55`.
- **Alice trades Bob's bearer stake token to Carol.** Alice offers Bob's bearer
  stake promise to Carol in return for non-transferable compute access. Carol
  locally accepts and issues Alice a scoped compute token. This demonstrates
  bearer-for-non-transferable exchange without a central exchange.
- **Dave asks for a peer-local quote.** Dave asks Bob for an exchange-rate quote;
  Bob answers from Bob's local wallet state, not from a market authority.
- **Mallory circulates stale Alice tokens.** Mallory starts with historical
  stale Alice bearer-token bytes. Dave accepts the first one for local evidence,
  redeems it with Alice, observes a broken outcome, and lowers local trust in
  both Alice and Mallory. Dave refuses Mallory's later stale-token offer because
  Dave's local trust changed.

## Promise Theory Fit

- Agents only promise their own behavior or the behavior of resources they
  control. No agent makes a promise on behalf of another agent.
- Trust is local. Wallet trust scores are private local judgments derived from
  observed keep/break/refuse outcomes.
- Tokens are issuer promises, not global permission objects. A bearer token is a
  transferable issuer promise; a non-transferable token is scoped to the original
  issuee.
- Collateral is modeled as a bearer stake promise, not enforcement. Its value is
  that breaking or keeping it changes later local trust and exchange terms.
- The kernel/relay layer carries signed exact bytes and records local delivery
  evidence. App-level promise-economy logic makes offers, counters, accepts,
  refuses, redeems, and observes.
- POC8 avoids a central exchange, central scheduler after startup, global trust
  score, global price oracle, shared token-status ledger, permission authority,
  authorization authority, and contract-enforcement framing. Source: `DI-sirus`.

## Incentives

- **Alice:** gains storage and compute service only after peers voluntarily make
  promises she can accept; she preserves selective sending by not sending data
  until terms clear her local threshold.
- **Bob:** gains exchange value and trust evidence by keeping storage promises;
  he can refuse or counter if Alice's terms do not clear Bob's local utility.
- **Carol:** gains Bob bearer stake value and relationship evidence by issuing
  scoped compute access that Alice can redeem.
- **Dave:** gains evidence from the first stale-token redemption, then avoids a
  later bad trade after local trust in Alice and Mallory decreases.
- **Mallory:** can try short-term stale-token circulation, but Dave's later
  refusal demonstrates the relationship cost of broken promise history.

## Current Limits

POC8 is still a bounded POC. The local plans are deterministic so tests and
Docker output stay reproducible. It does not define a final PromiseGrid token
standard, economics model, exchange protocol, trust API, kernel API, storage API,
compute API, transport standard, or SDK. Source: `DI-sirus`.

Run:

```sh
go test ./...
POC8_RUN_ID="$(date -u +%Y%m%d%H%M%S)" docker compose up --build --abort-on-container-exit
```
