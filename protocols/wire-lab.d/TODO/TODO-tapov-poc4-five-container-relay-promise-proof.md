# TODO-tapov: poc4 five-container relay promise proof

## Prior aliases

This TODO was minted directly under the post-`TE-mumuv` proquint handle scheme.
No integer or timestamp alias exists.

## Status

Planned. `poc4` follows `poc3` and should test multi-hop app promises across a
five-container topology where each container talks to only two neighboring
containers, runs one kernel, one relay, and two non-relay apps, and some apps
must rely on relays on third or fourth containers to reach the app that can
fulfill the promise.

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

- [ ] tapov.1 Review `poc3` evidence and decide which package, command, kernel,
  envelope, pCID, signature, test, and Compose patterns can be reused without
  preserving accidental two-container assumptions.
- [ ] tapov.2 Run DF for exact command names, package/type names, protocol
  surfaces, route representation, reciprocal promise message shapes, runtime
  paths, and demo topology before creating code.
- [ ] tapov.3 Define the relay app's promises and non-promises, including relay
  confirmation and refusal behavior.
- [ ] tapov.4 Define the fibonacci app and fibonacci-client promises and
  non-promises, including delayed result delivery.
- [ ] tapov.5 Define the storage app and storage-client promises and
  non-promises, including store confirmation and later read-by-key.
- [ ] tapov.6 Decide how hello, echo, and signed from `poc3` are reused, extended,
  or intentionally simplified for the larger topology.
- [ ] tapov.7 Define the kernel's `poc4` implementation promises and confirm
  what remains unchanged from `poc3`.
- [ ] tapov.8 Implement the directory layout under `implementations/poc4/` after
  DF/DI is locked.
- [ ] tapov.9 Add deterministic tests for relay path selection, reciprocal
  receive promises, fibonacci result delivery, storage confirmation/readback,
  signed relay carriage, local refusal, and app-local promise judgment.
- [ ] tapov.10 Add a deterministic five-container Compose demo command that Steve
  can run directly.
- [ ] tapov.11 Record final outcome: what worked, what was fake, what changed
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
