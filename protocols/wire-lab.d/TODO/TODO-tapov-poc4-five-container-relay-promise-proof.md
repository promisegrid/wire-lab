# TODO-tapov: poc4 five-container relay promise proof

## Prior aliases

This TODO was minted directly under the post-`TE-mumuv` proquint handle scheme.
No integer or timestamp alias exists.

## Status

Implemented and demo-validated. `poc4` follows `poc3` and tests multi-hop app
promises across a five-container topology where each container talks to only two
neighboring containers, runs one kernel, one relay, and two non-relay apps, and
some apps must rely on relays on third or fourth containers to reach the app
that can fulfill the promise.

## Decision Intent Log

### DI-hikun

ID: DI-hikun
Date: 2026-05-25 20:57:58
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Insert `poc4` at the top of the active queue as a five-container,
multi-hop relay proof with app packages under `implementations/poc4/hello/`,
`implementations/poc4/echo/`, `implementations/poc4/signed/`,
`implementations/poc4/fibonacci/`, `implementations/poc4/relay/`,
`implementations/poc4/storage/`, `implementations/poc4/kernel/`, and optional
shared `implementations/poc4/lib/` only as needed.
Intent: Extend `poc3` from same-node multi-app kernel pressure to multi-node
promise routing pressure. `poc4` should test whether applications can make
promises that depend on reciprocal receive promises and relay promises across a
bounded mesh without reintroducing RPC, service discovery, permission authority,
or a global trust source.
Constraints: Promise Theory first; kernels route/refuse and record local kernel
evidence only; relay is an app that promises to relay, not a kernel authority;
fibonacci, storage, and relay promises are conditional on reciprocal promises to
receive results or confirmations; promisees judge keep/break locally; each
container talks to exactly two neighboring containers in the demo topology;
implementation-specific command names, package/type names, protocol surfaces,
runtime paths, and exact relay route representation require a later DF/DI before
code edits.
Affects: `implementations/poc4/**`; `protocols/wire-lab.d/TODO/TODO.md`; future
kernel/app/relay guide evidence.

### DI-simuk

ID: DI-simuk
Date: 2026-05-25 21:02:46
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Reuse `poc3`'s package/command split, pCID-selected envelope library,
receive-promise registration, local kernel evidence discipline, app-local
promise judgment, signed proof-slot pattern, Go tests, Dockerfile shape, and
Compose-driven demo approach as starting material for `poc4`; do not reuse
`poc3`'s single-peer kernel forwarding, two-container script assumptions,
single-hop app destination field, or direct app-to-target mental model as final
`poc4` architecture.
Intent: `poc4` exists to test multi-hop relay promises, reciprocal receive
promises, and app fulfillment across a bounded five-node mesh. The useful
`poc3` parts keep the wire and package baseline stable, while the rejected parts
are exactly the assumptions that would hide relay behavior or collapse back into
RPC-like direct calls.
Constraints: Keep `grid([42(pCID), payload, ...])` as the envelope baseline; keep
kernels out of application keep/break judgment; make relay an app-level promiser;
run DF before choosing exact command names, package/type names, route
representation, reciprocal-promise message fields, runtime paths, or whether
relay paths are source-routed, next-hop-routed, or locally promised hop-by-hop.
Affects: `implementations/poc4/**`; `implementations/poc3/**`; this TODO.

### DI-bigub

ID: DI-bigub
Date: 2026-05-25 21:15:15
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement `poc4` with relay-owned neighbor TCP, not kernel-owned
multi-hop routing; use packages `kernel`, `relay`, `hello`, `echo`, `signed`,
`fibonacci`, `storage`, and `lib`; expose named commands `poc4-kernel`,
`poc4-relay`, `poc4-hello`, `poc4-echo`, `poc4-signed`, `poc4-fibonacci`, and
`poc4-storage`; use per-package pCIDs; encode relay carriage as relay-pCID
wrappers that carry exact inner envelope bytes; correlate confirmations/results
by exact byte hash; use local route promises for only the demo flows; run five
containers when Docker is available.
Intent: Preserve the promise-first boundary from `poc3` while testing multi-hop
application promises where relays are app-level promisers, not kernel routers,
service registries, permission authorities, or hidden RPC dispatchers.
Constraints: Kernels remain local app/kernel boundaries and local evidence
writers; relays own neighbor transport and local next-hop promises; promisee apps
judge keep/break locally; demo values are `fib(10)=55` and
`poc4-key=poc4-value`; runtime temp paths may use `/tmp/wire-lab-gocache/**` and
`/tmp/wire-lab-poc4-run/**`; committed implementation lives under
`implementations/poc4/**`.
Affects: `implementations/poc4/**`;
`protocols/wire-lab.d/TODO/TODO-tapov-poc4-five-container-relay-promise-proof.md`.

### DI-rinuv

ID: DI-rinuv
Date: 2026-05-25 21:35:33
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Make the `poc4` Compose demo safe to run with
`docker compose up --build --abort-on-container-exit` by adding an explicit
shared completion gate. Each container marks local completion only after its
expected app promises finish, then keeps its kernel and relay alive until all
five containers have marked completion.
Intent: Preserve the useful Compose shutdown behavior without letting
short-lived idle/client apps abort the run before multi-hop relay, computation,
storage, signed-evidence, and echo promises finish.
Constraints: Keep one kernel, one relay, and two non-relay apps per container;
do not add a central trust or routing authority; completion markers are local
demo process coordination, not PromiseGrid protocol evidence; use the existing
`implementations/poc4/scripts/*.sh` shape and a shared Compose volume mounted at
`/run/poc4/$POC4_RUN_ID`.
Affects: `implementations/poc4/compose.yaml`;
`implementations/poc4/README.md`; `implementations/poc4/scripts/*.sh`; this
TODO.

### DI-ponor

ID: DI-ponor
Date: 2026-05-25 23:00:33
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Treat the successful `poc4` Docker Compose run as provisional
executable evidence for multi-hop app promises over app-level relay promises,
and update `DEV-GUIDE-RESOURCES.md` accordingly.
Intent: Preserve the observed result that five containers each ran one kernel,
one relay, and two non-relay apps; signed hello, fibonacci, storage
confirmation/readback, and echo flows completed; relay hops crossed third and
fourth containers; kernels stayed local app/kernel evidence boundaries; app
promisees judged keep/break locally.
Constraints: `poc4` remains executable POC evidence, not a final SDK, final
kernel API, final relay protocol, or stable storage/computation API. The
completion markers under `/run/poc4/$POC4_RUN_ID` are demo process coordination,
not PromiseGrid protocol evidence. The result should inform guide-writer
resources and `TODO-binag`/`SIM-fovip` review, but should not close `DR-davod`
by itself.
Affects: `implementations/poc4/**`; `DEV-GUIDE-RESOURCES.md`;
`protocols/wire-lab.d/TODO/TODO-binag-promisegrid-kernel-design-resolution.md`;
`simulations/SIM-fovip-kernel-promise-boundary-port-contract/`; this TODO.

## Topology

`poc4` uses a five-node ring:

```text
            [Ellen]
          /         \
     [Alice]       [Dave]
        |             |
       [Bob] ---- [Carol]
```

Each container talks only to its two ring neighbors:

| Container | Neighbors |
|---|---|
| Alice | Bob, Ellen |
| Bob | Alice, Carol |
| Carol | Bob, Dave |
| Dave | Carol, Ellen |
| Ellen | Dave, Alice |

Each container runs one kernel, one relay, and two non-relay apps:

| Container | Kernel | Relay | App 1 | App 2 |
|---|---|---|---|---|
| Alice | `poc4-kernel` | `poc4-relay` | `hello` | `fibonacci-client` |
| Bob | `poc4-kernel` | `poc4-relay` | `echo` | `signed` |
| Carol | `poc4-kernel` | `poc4-relay` | `fibonacci` | `storage-client` |
| Dave | `poc4-kernel` | `poc4-relay` | `storage` | `signed` |
| Ellen | `poc4-kernel` | `poc4-relay` | `hello` | `echo` |

## Promise flows under test

- Alice's `fibonacci-client` asks Carol's `fibonacci` app to calculate `fib(n)`
  through Alice relay -> Bob relay -> Carol relay.
- Carol's `storage-client` asks Dave's `storage` app to store a key-value pair
  through Carol relay -> Dave relay, then later asks for the value by key.
- Alice's `hello` app reaches Dave's `signed` app through Alice relay -> Ellen
  relay -> Dave relay.
- Bob's `echo` app reaches Ellen's `echo` app through Bob relay -> Alice relay ->
  Ellen relay, proving that a same-protocol app can still require relay help when
  the target is not a direct neighbor.

## Promise semantics

- `fibonacci` promises to calculate the nth Fibonacci number only in return for
  a promise to receive the result when the calculation is complete.
- `storage` promises to store a key-value pair only in return for a promise to
  receive confirmation when storage is complete.
- `storage` promises to return a value later when presented with the key as part
  of a promise to receive the value.
- `relay` promises to relay a message from one app to another only in return for
  a promise to receive confirmation when the relay is complete.
- Kernels promise local app/kernel and kernel/kernel message handling; they do
  not promise relay semantics, route success, or application keep/break status.
- Apps and relays make their own promises and record their own local judgments.

## Initial layout

- `implementations/poc4/hello/*.go`
- `implementations/poc4/echo/*.go`
- `implementations/poc4/signed/*.go`
- `implementations/poc4/fibonacci/*.go`
- `implementations/poc4/relay/*.go`
- `implementations/poc4/storage/*.go`
- `implementations/poc4/kernel/*.go`
- `implementations/poc4/lib/` only if needed

## Subtasks

- [x] tapov.1 Review `poc3` evidence and decide which package, command, kernel,
  envelope, pCID, signature, test, and Compose patterns can be reused without
  preserving accidental two-container assumptions.
- [x] tapov.2 Run DF for exact command names, package/type names, protocol
  surfaces, route representation, reciprocal promise message shapes, runtime
  paths, and demo topology before creating code.
- [x] tapov.3 Define the relay app's promises and non-promises, including relay
  confirmation and refusal behavior.
- [x] tapov.4 Define the fibonacci app and fibonacci-client promises and
  non-promises, including delayed result delivery.
- [x] tapov.5 Define the storage app and storage-client promises and
  non-promises, including store confirmation and later read-by-key.
- [x] tapov.6 Decide how hello, echo, and signed from `poc3` are reused, extended,
  or intentionally simplified for the larger topology.
- [x] tapov.7 Define the kernel's `poc4` implementation promises and confirm
  what remains unchanged from `poc3`.
- [x] tapov.8 Implement the directory layout under `implementations/poc4/` after
  DF/DI is locked.
- [x] tapov.9 Add deterministic tests for relay path selection, reciprocal
  receive promises, fibonacci result delivery, storage confirmation/readback,
  signed relay carriage, local refusal, and app-local promise judgment.
- [x] tapov.10 Add a deterministic five-container Compose demo command that Steve
  can run directly.
- [x] tapov.11 Record final outcome: what worked, what was fake, what changed
  from `poc3`, and whether this evidence should update `DEV-GUIDE-RESOURCES.md`,
  `SIM-fovip`, or a new relay/storage/fibonacci simulation.

## Acceptance criteria

- The demo runs five containers with exactly one kernel, one relay, and two
  non-relay apps per container.
- Each container has only two container neighbors in the configured topology.
- At least two successful app promises require relays on a third or fourth
  container before reaching the fulfilling app.
- Fibonacci, storage, and relay behavior are described as reciprocal promises,
  not RPC calls or permissioned service requests.
- Kernel evidence remains local kernel evidence; application keep/break judgment
  remains local to the relevant promisee app or relay.
- The POC remains small enough to read as an executable thought experiment.

## `poc3` pattern review for `poc4`

Reusable starting points:

- Keep the importable package plus `cmd/` wrapper shape from `poc3`; it made
  application behavior testable without turning command entrypoints into the
  architecture.
- Keep the shared `lib` envelope, pCID, CBOR, frame, evidence, receive-promise,
  and signature helpers as the starting point.
- Keep per-app pCIDs and a kernel-local receive-promise pCID. `poc4` can add
  relay, fibonacci, and storage pCIDs without changing the outer envelope.
- Keep the signed app's proof-slot pressure: exact signable envelope bytes can
  be witnessed without claiming global trust.
- Keep local kernel evidence phrasing: kernels record receive, send, deliver,
  refusal, and broken transport evidence about their own behavior only.

Reuse only after refactor:

- The kernel's `PeerAddress` field must become a bounded neighbor set. `poc4`
  containers have exactly two neighbors, not one peer.
- App payloads need an explicit relay target and route/hop representation chosen
  during DF. The `poc3` `to` field is too direct for third/fourth-container relay
  promises.
- The receive-promise helper should probably grow correlation IDs or promise IDs
  so fibonacci, storage, and relay confirmations can be matched to the promise
  they answer.
- The Docker scripts should keep explicit process status handling, but need a
  five-container startup plan and no hidden all-to-all network assumptions.

Do not carry forward as architecture:

- Do not let the kernel choose multi-hop routes. Relay is an app-level promise;
  the kernel only handles local and neighbor message carriage.
- Do not model relay as a service registry or transparent router. Each relay hop
  is a promise by that relay app to carry a message and later account for the
  relay outcome.
- Do not preserve the direct app-to-target mental model where every app can name
  a reachable destination and the kernel handles the rest.
- Do not treat storage, fibonacci, or relay confirmations as kernel delivery
  receipts. They are application promises and application-local judgments.

## Final outcome

What worked:

- `docker compose up --build --abort-on-container-exit` completed with all five
  containers exiting `0` after the cleanup fixes in `505df60`.
- The demo ran Alice, Bob, Carol, Dave, and Ellen with one kernel, one relay, and
  two non-relay apps per container.
- Alice received signed hello evidence from Dave through Alice -> Ellen -> Dave
  relay promises.
- Alice received `fib(10)=55` from Carol through Alice -> Bob -> Carol relay
  promises.
- Carol stored and read back `poc4-key=poc4-value` through Carol -> Dave relay
  promises.
- Bob completed an echo round trip with Ellen through Bob -> Alice -> Ellen
  relay promises.
- Kernels recorded local app/kernel receive and delivery evidence only; relays
  owned neighbor TCP and hop promises; apps made reciprocal receive promises and
  judged their own app-level outcomes.

What was fake or still bounded:

- The completion markers under `/run/poc4/$POC4_RUN_ID` are Docker-demo process
  coordination, not PromiseGrid protocol evidence.
- Route tables are hand-configured local relay promises for the demo, not
  discovery, routing consensus, or a stable relay protocol.
- Storage is in-memory process-local state; fibonacci is a toy computation; echo
  and signed hello are proof surfaces, not final app APIs.
- The kernel accepts receive promises by pCID and delivers exact bytes; it does
  not implement durable identity, trust decay, resource accounting, sandboxing,
  host-driver promises, or a frozen porting contract.

What changed from `poc3`:

- `poc3` showed two containers with multiple local apps sharing one
  pCID-selected app/kernel boundary.
- `poc4` keeps that local-kernel boundary but moves cross-node reachability into
  relay apps, proving the kernel does not need to become a multi-hop router or
  service registry for app promises to cross non-neighbor paths.

Follow-on implications:

- `DEV-GUIDE-RESOURCES.md` should cite `poc4` as provisional executable evidence
  for multi-hop app/relay/kernel separation.
- `TODO-binag` and `SIM-fovip` should absorb the evidence when deciding whether
  the current kernel implementation-promise packet is sufficient.
- New focused sims may be useful for relay promises, storage promises, and
  computation promises, but `poc4` itself should not be promoted as a stable API
  surface.
