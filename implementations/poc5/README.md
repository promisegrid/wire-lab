# poc5 local trust selective sending proof

`poc5` is an executable thought experiment for app-local trust, broken promises,
and selective sending across the same five-node ring used by `poc4`.

```text
Alice -- Bob -- Carol -- Dave -- Ellen -- Alice
```

The demo uses this trust sequence:

- Alice sends Bob a low-sensitivity storage probe.
- Bob confirms the store but deliberately returns the wrong value on readback.
- Alice records `promise_broken` as Alice-local evidence, not global truth.
- Alice records `trust_decreased` for Bob in Alice's own local trust table.
- Alice records `selective_send_declined` and does not send sensitive data to Bob.
- Alice records `selective_send_chosen` and sends the sensitive storage promise to Dave.
- Dave stores and returns the value correctly; Alice records `promise_kept`.

Run it with Docker Compose v2:

```sh
cd implementations/poc5
POC5_RUN_ID="$(date -u +%Y%m%d%H%M%S)" docker compose up --build --abort-on-container-exit
```

The output is intentionally verbose. Kernels record local delivery evidence,
relays record hop evidence, and apps print or emit their own local promise
judgments. No app, relay, or kernel claims global trust authority, permission,
authorization, or conformance power.

Each container writes a demo-only completion marker under the shared
`/run/poc5` volume after its expected local app promises complete, then waits for
all five markers before exiting. That makes `--abort-on-container-exit` mean
"stop after the bounded proof completes" rather than "stop when the first
short-lived app finishes."

`poc5` is POC evidence, not a final SDK, stable storage API, stable trust API, or
kernel API. Trust state is in-memory and local to Alice for this proof.
