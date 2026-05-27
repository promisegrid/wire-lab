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

Run:

```sh
go test ./...
POC7_RUN_ID="$(date -u +%Y%m%d%H%M%S)" docker compose up --build --abort-on-container-exit
```

`poc7` is not a final token standard, economics model, exchange protocol, trust
API, kernel API, storage API, compute API, TCP transport standard, or SDK.
Source: `DI-tugih`; `DI-fibok`; `DI-tanat`; `DI-rodog`; `DI-hanih`.
The Mallory-to-Dave stale-token flow is corrected under `DI-pabot`.
