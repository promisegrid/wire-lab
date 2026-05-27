# poc7-capability-token-exchange

`poc7-capability-token-exchange` is executable POC evidence for
`scenarios/promise-economy-capability-token-exchange/`.

The demo uses five containers in a ring:

```text
Alice -- Bob -- Carol -- Dave -- Mallory -- Alice
```

Each container runs one local kernel role, one app-level relay role, and
role-specific token/resource apps inside a single bounded process. Peers use
boring length-framed TCP; each app message is a signed CBOR `grid(...)` tag
(`0x67726964`, decimal `1735551332`) wrapping the
`[42(pCID), payload, proof]` slot vector. Kernels record local delivery
evidence only. Relays carry app messages. Apps make local promise judgments.
The visible app message kinds are promise-shaped: resource promise requests,
promise presentations for fulfillment, promise receipts, reciprocal exchange
promises, holder-initiated transfers, and issuer-local revocation notices.

The scenario now performs real POC work:

- Bob and Carol redeem signed token bytes through issuer-local state.
- Bob stores a key/value for Carol, then returns it through a separate read token.
- Carol computes Fibonacci 10 for Bob and returns `55`.
- Carol offers a Bob bearer token to Alice and receives a non-transferable Alice
  data token promised specifically to Carol.
- Carol redeems that Alice data token and receives Alice's local data payload.
- Mallory voluntarily circulates a revoked Alice bearer token to Dave; Dave
  redeems it with Alice and locally updates trust in both Alice and Mallory from
  the broken redemption outcome.
- Dave then refuses Mallory's later stale-token transfer because Dave's own
  deterministic local policy now scores Mallory and Alice below Dave's threshold.

## Assessment

POC7 works well as executable protocol evidence. A full five-container run now
exercises signed CBOR `grid(...)` app messages, pCID-selected payload parsing,
signed token bytes, real storage / compute / data fulfillment, issuer-local
revocation and redemption state, local accept/refuse decisions, and local trust
updates. It does not work as a finished PromiseGrid economy: peer discovery,
independent goal formation, negotiated terms, collateral, opportunity cost,
durable strategy, and realistic exchange rates remain outside this POC. Source:
`DI-tugih`; `DI-fibok`; `DI-tanat`; `DI-pabot`; `DI-rodog`; `DI-hanih`.

## Transaction Walkthrough

- **Alice to Bob:** Alice sends Bob bearer and non-transferable Alice data
  promise tokens. Bob locally accepts the tokens, then Bob presents the
  non-transferable token back to Alice for fulfillment and receives Alice's data.
- **Alice to Bob for Carol:** Alice asks Bob to issue storage promise tokens for
  Carol. Bob locally decides to issue separate store, read, and trade bearer
  tokens because Bob's local policy scores those opportunities above threshold.
- **Carol to Bob:** Carol accepts Bob's storage tokens, presents one token to ask
  Bob to store `carol-note`, then presents a different token to ask Bob to return
  the stored value.
- **Alice to Carol to Bob:** Alice asks Carol to issue a compute token for Bob.
  Bob accepts that token and presents it to Carol; Carol computes Fibonacci 10
  and returns `55`.
- **Carol to Alice:** Carol offers Alice a Bob bearer storage-trade token. Alice
  locally accepts the reciprocal exchange and issues Carol a non-transferable
  Alice data token, which Carol later redeems with Alice.
- **Mallory to Dave to Alice:** Mallory accepts Alice's revoked bearer token and
  voluntarily circulates it to Dave. Dave accepts the first stale token for local
  evidence, presents it to Alice, observes a broken redemption outcome, and
  lowers local trust in both Alice and Mallory.
- **Mallory to Dave again:** Mallory circulates a second stale Alice token. Dave
  refuses it because Dave's local policy now scores Alice and Mallory below
  Dave's threshold.

## Promise Theory Fit

- **Good fit:** Tokens are issuer promises, not global permission objects. A
  token means "the issuer promises this resource behavior to the bearer or named
  promisee," and redemption produces local evidence that the promise was kept,
  broken, or refused.
- **Good fit:** Trust is local. Bob, Carol, Dave, and Mallory do not consult a
  central trust authority; each wallet records local evidence and updates local
  trust scores from observed outcomes.
- **Good fit:** Agents can refuse. The scenario may introduce an opportunity,
  but issue, accept, redeem, transfer, trade, and quote actions pass through each
  agent's local deterministic policy before action.
- **Good fit:** Kernels and relays carry signed exact bytes and record delivery
  observations; issuer, trader, and resource apps make and judge application
  promises.
- **Current gap:** Alice's harness script still creates most opportunities,
  routes, and timing. That is useful POC scaffolding, but not yet autonomous
  multi-agent economy behavior.

## Impositions And Incentives

POC7 avoids the main protocol-level imposition: no peer is treated as commanded
by another peer or by a central authority. The remaining imposition is harness
scaffolding: Alice's script initiates the sequence and asks other nodes to
consider opportunities. Those nodes can refuse, but they do not yet independently
discover needs, advertise offers, bargain, or choose long-term strategies.

- **Alice:** gains data-sharing evidence, receives Bob's bearer storage token in
  trade, and observes whether her own tokens are being redeemed after revocation.
- **Bob:** gains trust/evidence by keeping storage promises and receives Alice or
  Carol tokens that Bob values enough to accept or redeem.
- **Carol:** gains storage service from Bob, compute exchange evidence with Bob,
  and non-transferable access to Alice's private data after trading a bearer
  token.
- **Dave:** gains evidence value from accepting the first stale token, then acts
  defensively by refusing later stale-token circulation after local trust drops.
- **Mallory:** has a short-term incentive to circulate bearer tokens, including
  stale ones, but Dave's later refusal shows the cost: broken promise history
  degrades Mallory's local relationship with Dave.
- **System pressure:** promise keeping tends to increase local trust; broken,
  refused, or stale-token behavior tends to reduce local trust and make later
  exchanges less likely. That is the key POC7 incentive loop.

Run:

```sh
go test ./...
POC7_RUN_ID="$(date -u +%Y%m%d%H%M%S)" docker compose up --build --abort-on-container-exit
```

`poc7` is not a final token standard, economics model, exchange protocol, trust
API, kernel API, storage API, compute API, TCP transport standard, or SDK.
Source: `DI-tugih`; `DI-fibok`; `DI-tanat`; `DI-rodog`; `DI-hanih`.
The Mallory-to-Dave stale-token flow is corrected under `DI-pabot`.
