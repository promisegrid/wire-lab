# TODO-hipak: poc5 local trust selective sending

## Prior aliases

This TODO was minted directly under the post-`TE-mumuv` proquint handle scheme.
No integer or timestamp alias exists.

## Status

Implemented and demo-validated. `poc5` follows `poc4` and tests the local trust loop
that `poc4` intentionally did not decide: Alice should send data for storage or
computation only when her own local promise history gives her enough trust in the
specific peer relationship.

## Decision Intent Log

### DI-rarim

ID: DI-rarim
Date: 2026-05-25 23:24:58
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Insert `poc5` as a new local-trust and selective-sending proof under
`implementations/poc5/**`, reusing the `poc4` five-container ring, package split,
Compose shape, app/kernel boundary, relay-owned neighbor transport, and
pCID-selected envelope baseline while adding app-local trust evidence, deliberate
promise breakage, trust decrease, and later selective refusal to send sensitive
data to an untrusted peer.
Intent: Extend the executable kernel evidence from `poc4` without turning the
kernel, relay, or topology into an authority. `poc5` should prove that
PromiseGrid selective sending is an app-local promise and trust decision: Alice
observes Bob break a storage/computation promise, records local evidence, lowers
her own trust in Bob, refuses a later sensitive send to Bob, and chooses a
better-trusted peer instead.
Constraints: No global trust authority, reputation service, permission server,
policy-enforcement service, or kernel-mediated trust decision. Trust records are
in-memory and demo-local for this proof. Kernels deliver local app messages and
record local observations only. Relays forward neighbor messages only. Runtime
paths are limited to Docker Compose state under `/run/poc5/$POC5_RUN_ID` and Go
cache state under `/tmp/wire-lab-gocache/**`. The first implementation may reuse
`poc4` code structure aggressively to keep the proof readable.
Affects: `implementations/poc5/**`;
`protocols/wire-lab.d/TODO/TODO-hipak-poc5-local-trust-selective-sending.md`;
`protocols/wire-lab.d/TODO/TODO.md`; `DEV-GUIDE-RESOURCES.md` after demo
evidence exists.


### DI-fofik

ID: DI-fofik
Date: 2026-05-25 23:35:03
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Treat the successful `poc5` Docker Compose run as provisional executable evidence for app-local trust decrease and selective sending after observed promise breakage, and update guide-writer resources accordingly.
Intent: Preserve the observed result that Alice sent Bob only low-sensitivity probe data, observed Bob return an incorrect read result, recorded Alice-local `promise_broken` and `trust_decreased` evidence, declined a later sensitive send to Bob, chose Dave instead, and judged Dave's storage promise kept. This evidence demonstrates the PromiseGrid trust loop without making the kernel, relay, or any shared service into a trust authority.
Constraints: `poc5` remains executable POC evidence, not a final SDK, stable storage API, stable trust API, durable trust store, or final kernel API. Trust state is in-memory and local to Alice. Bob's broken promise is deterministic demo behavior. Completion markers under `/run/poc5/$POC5_RUN_ID` are Docker process coordination, not PromiseGrid protocol evidence.
Affects: `implementations/poc5/**`; `DEV-GUIDE-RESOURCES.md`; this TODO.

## Topology

`poc5` uses the same five-node ring as `poc4`:

```text
            [Ellen]
          /         \
     [Alice]       [Dave]
        |             |
       [Bob] ---- [Carol]
```

Each container talks only to its two ring neighbors. Each container runs one
kernel, one relay, and two non-relay apps where possible. The non-relay apps are
demo roles, not final PromiseGrid APIs.

## Promise flows under test

- Alice probes Bob with a low-sensitivity storage/computation promise request.
- Bob deliberately breaks that promise in a deterministic way.
- Alice records her own local evidence of the broken promise and lowers her own
  trust in Bob.
- Alice refuses to send later sensitive data to Bob because Bob no longer meets
  Alice's local trust threshold.
- Alice sends the sensitive request to Carol or Dave instead after their
  reciprocal promises meet Alice's local threshold.
- Carol or Dave keeps the promise, and Alice records that local outcome.

## Subtasks

- [x] hipak.1 Reuse the `poc4` scaffold under `implementations/poc5/` while
  changing module names, commands, Compose network/volume names, run IDs, and
  documentation to `poc5`.
- [x] hipak.2 Add app-local trust/evidence records that distinguish local
  observation from global authority.
- [x] hipak.3 Add deterministic broken-promise behavior for Bob's provider role.
- [x] hipak.4 Add Alice-side selective sending that declines Bob after a broken
  promise and chooses a better-trusted peer.
- [x] hipak.5 Keep kernel and relay behavior non-authoritative and local.
- [x] hipak.6 Add or update tests for trust decrease, selective refusal, and kept
  fallback promise behavior.
- [x] hipak.7 Add a bounded Compose demo that exits cleanly with
  `--abort-on-container-exit`.
- [x] hipak.8 Update `DEV-GUIDE-RESOURCES.md` only after executable evidence is
  validated.

## Acceptance criteria

- The demo runs five containers with the same bounded neighbor topology as
  `poc4`.
- Logs show `promise_broken`, `trust_decreased`, `selective_send_declined`,
  `selective_send_chosen`, and `promise_kept`.
- Alice refuses a later sensitive send to Bob because of Alice's own local trust
  judgment, not because of an external authority.
- Kernels and relays do not decide trust, permission, authorization,
  conformance, or app-level keep/break.
- The proof remains small enough to read as an executable thought experiment.

## Final outcome

What worked:

- `docker compose up --build --abort-on-container-exit` completed with all five
  containers exiting `0`.
- Alice observed Bob return `poc5-probe-value-broken` when Alice expected
  `poc5-probe-value`.
- Alice emitted local `promise_broken`, `trust_decreased`,
  `selective_send_declined`, and `selective_send_chosen` evidence.
- Alice sent `poc5-sensitive-key=poc5-sensitive-value` to Dave instead of Bob.
- Dave stored and returned the sensitive value, and Alice emitted `promise_kept`.
- Kernels stayed local app/kernel delivery boundaries; relays owned neighbor hop
  promises; app-local code judged trust and keep/break.

What remains bounded:

- Trust state is in-memory and seeded by the demo.
- Bob's broken promise is deterministic test behavior, not an adversarial model.
- Storage remains toy process-local key-value state.
- The proof does not define a final trust API, storage API, or kernel API.
