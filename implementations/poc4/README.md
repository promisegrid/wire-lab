# poc4 five-container relay promise proof

`poc4` is an executable thought experiment for app-level relay promises across a
five-node ring. Each container runs one kernel, one relay app, and two non-relay
apps.

```text
Alice -- Bob -- Carol -- Dave -- Ellen -- Alice
```

The demo uses these cross-container app promises:

- Alice `fibonacci-client` asks Carol `fibonacci` for `fib(10)=55` through Bob.
- Carol `storage-client` asks Dave `storage` to store and return
  `poc4-key=poc4-value`.
- Alice `hello` asks Dave `signed` for signed evidence through Ellen.
- Bob `echo` asks Ellen `echo` through Alice.

Run it with Docker Compose v2:

```sh
cd implementations/poc4
POC4_RUN_ID="$(date -u +%Y%m%d%H%M%S)" docker compose up --build --abort-on-container-exit
```

The output is intentionally verbose. Kernels record local delivery evidence,
relays record hop evidence, and apps print their own local keep/break judgments.
No app, relay, or kernel claims global trust authority.

Each container writes a demo-only completion marker under the shared
`/run/poc4` volume after its expected local app promises complete, then waits for
all five markers before exiting. That makes `--abort-on-container-exit` mean
"stop after the bounded proof completes" rather than "stop when the first
short-lived app finishes."
