# poc7-capability-token-exchange

`poc7-capability-token-exchange` is executable POC evidence for
`scenarios/promise-economy-capability-token-exchange/`.

The demo uses five containers in a ring:

```text
Alice -- Bob -- Carol -- Dave -- Mallory -- Alice
```

Each container runs one local kernel role, one app-level relay role, and
role-specific token/resource apps inside a single bounded process. Kernels record
local delivery evidence only. Relays carry app messages. Apps make local promise
judgments.

Run:

```sh
go test ./...
POC7_RUN_ID="$(date -u +%Y%m%d%H%M%S)" docker compose up --build --abort-on-container-exit
```

`poc7` is not a final token standard, economics model, exchange protocol, trust
API, kernel API, or SDK. Source: `DI-tugih`.
