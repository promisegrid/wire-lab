# TODO-hozaz: poc3 multi-app kernel proof

## Prior aliases

This TODO was minted directly under the post-`TE-mumuv` proquint handle scheme.
No integer or timestamp alias exists.

## Status

Implemented and Docker-demo verified. `poc3` follows `poc2` and tests the next
level of kernel/app shape: several local apps sharing one PromiseGrid-style
kernel boundary without turning the kernel into an RPC authority or app-specific
service registry.

## Decision Intent Log

### DI-vosof

ID: DI-vosof
Date: 2026-05-25 19:20:33
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Insert `poc3` after `poc2` as a multi-app proof with code organized
under `implementations/poc3/hello/*.go`, `implementations/poc3/echo/*.go`,
`implementations/poc3/signed/*.go`, `implementations/poc3/kernel/*.go`, and an
optional shared `implementations/poc3/lib/` only as needed.
Intent: Keep the next proof of concept simple but broader than `poc2`: multiple
apps should make, receive, echo, and sign promise messages through the same
kernel boundary so we can test whether the kernel shape remains promise-first
under app diversity.
Constraints: Promise Theory first; no external permission authority, service
registry, RPC dispatch model, or global trust authority; all app/kernel
boundaries should continue to use pCID-selected `grid([42(pCID), payload, ...])`
messages unless a later TE/DI changes that; `lib/` is optional and must not
become a dumping ground for premature framework code; implementation-specific
function names, variable names, command names, and exact protocol surfaces still
require a later DF/DI before code edits.
Affects: `implementations/poc3/`; `protocols/wire-lab.d/TODO/TODO.md`;
future kernel/app design evidence.

### DI-fubir

ID: DI-fubir
Date: 2026-05-25 19:22:37
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Reuse `poc2`'s envelope, pCID, framing, evidence, and promise-first
boundary patterns as starting material for `poc3`, but do not reuse its
single-app receiver slot, single-protocol dispatcher, single package layout, or
hello-specific kernel behavior as final `poc3` architecture.
Intent: `poc2` proved that app/kernel and kernel/kernel boundaries can share the
same pCID-selected `grid([42(pCID), payload, ...])` shape. `poc3` should build
on that evidence while intentionally testing app diversity, multiple payload
kinds, local echo promises, and signed-message evidence without letting the
kernel become a service registry or RPC dispatcher.
Constraints: Keep Promise Theory vocabulary first; preserve local observation
and local refusal semantics; keep exact-byte evidence; avoid global trust claims;
run DF before naming commands, packages, functions, variables, or exact runtime
topology for `poc3`.
Affects: `implementations/poc3/`; `implementations/poc2/`; this TODO.

### DI-horak

ID: DI-horak
Date: 2026-05-25 19:44:15
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Implement `poc3` as two containers, each running one kernel and all
three apps; use importable Go packages under `implementations/poc3/kernel/`,
`implementations/poc3/hello/`, `implementations/poc3/echo/`,
`implementations/poc3/signed/`, and shared `implementations/poc3/lib/`, with
command wrappers under each package's `cmd/`; use per-app pCIDs plus a
kernel-local receive-promise pCID; keep kernel evidence internal/stdout-only;
let promisee apps make local keep/break/uncertain judgments; sign the canonical
envelope bytes for the `signed` app.
Intent: Preserve `poc2`'s useful pCID-selected envelope evidence while testing
the next kernel pressure point: multiple apps on each node sharing one local
kernel boundary without converting the kernel into an RPC dispatcher, service
registry, permission authority, or global trust judge.
Constraints: All app/kernel and kernel/kernel payloads use CBOR
`grid([42(pCID), payload, ...])`; pCID remains Protocol CID; receive promises
describe what a local app promises to receive and do not grant permission;
kernel evidence is local evidence about kernel actions, not application
promise-status judgment; demo runtime evidence may use `/tmp/wire-lab-poc3-run/**`;
Go validation may use `/tmp/wire-lab-gocache/**`.
Affects: `implementations/poc3/**`;
`protocols/wire-lab.d/TODO/TODO-hozaz-poc3-multi-app-kernel-proof.md`.

## Summary

Build `poc3`, a multi-app continuation of `poc2`.

`poc3` should demonstrate:

- A `hello` app that sends and receives a minimal promise-message.
- An `echo` app that reflects a received promise-message as a new local promise,
  not as a commanded RPC response.
- A `signed` app that demonstrates the current envelope/proof direction without
  pretending signatures create global trust.
- A shared `kernel` boundary that routes, records evidence, and refuses locally
  unsupported pCIDs without becoming an application authority.
- Optional shared `lib` code only where duplication would otherwise obscure the
  promise/message model.

This is a proof of concept, not the final PromiseGrid SDK, kernel API, app
framework, identity model, signature suite, storage layer, or trust system.

## Initial layout

- `implementations/poc3/hello/*.go`
- `implementations/poc3/echo/*.go`
- `implementations/poc3/signed/*.go`
- `implementations/poc3/kernel/*.go`
- `implementations/poc3/lib/` only if needed

## `poc2` pattern review for `poc3`

Reusable starting points:

- The outer envelope shape is the right default for `poc3`: CBOR
  `grid([42(pCID), payload, ...])` with slot 0 as the Protocol CID link and
  slot 1 as protocol-owned payload bytes.
- The lightweight `ProtocolCID` pattern is useful for a POC: it hashes embedded
  spec bytes and produces deterministic CIDv1 raw/sha2-256 bytes without making
  the POC depend on a full IPLD library.
- The length-framed `FrameConn` transport boundary is useful because it keeps
  the demo focused on message semantics rather than stream parsing.
- The `EvidenceLog` pattern is worth carrying forward: every kept, refused,
  broken, or not-promised outcome is a local observation, not a global authority
  fact.
- The refusal path is useful: unsupported or unreadable pCIDs produce local
  evidence and an observation instead of an authorization failure or exception
  framed as command rejection.

Reuse only after refactor:

- The kernel should keep separate app/kernel and kernel/kernel boundaries, but
  `poc2`'s single `receiver` field must become an app-relationship or local
  promise table keyed by pCID, app identity, or receive promise semantics chosen
  during DF.
- The `kind` field inside one hello payload was adequate for `poc2`; `poc3`
  needs a clearer rule for whether `hello`, `echo`, and `signed` are separate
  pCIDs or one pCID with protocol-owned payload kinds.
- The observation envelope is useful, but `poc3` should distinguish local app
  observations, peer observations, echo promises, and signed evidence without
  collapsing them into a generic RPC response.
- The Docker two-container demo is useful, but `poc3` may need either multiple
  app processes per container or a local-only mode so tests remain deterministic.

Do not carry forward as architecture:

- Do not keep hello-specific kernel methods as the main kernel design; `poc3`
  should make app behavior live in apps and keep the kernel focused on local
  message acceptance, evidence, routing, delivery, and refusal promises.
- Do not keep a single registered local receiver as the model for multiple apps;
  that assumption hides the app-diversity question `poc3` exists to test.
- Do not treat signatures as trust decisions. A signed app can witness exact
  bytes or a claimed promiser, but Alice, Bob, and Carol still make local trust
  judgments from evidence and relationships.
- Do not let an `echo` app become a disguised remote procedure call. Echoing is
  a local promise by the echo app to make a new message in response to a message
  it chose to interpret and keep.

## Subtasks

- [x] hozaz.1 Review `poc2` evidence and decide which code patterns can be
  reused without carrying over accidental single-app assumptions.
- [x] hozaz.2 Run DF for exact command names, package/type names, protocol
  surfaces, and runtime topology before creating code.
- [x] hozaz.3 Define the `hello` app's promises and non-promises.
- [x] hozaz.4 Define the `echo` app's promises and non-promises, including how
  echoing differs from obeying a remote command.
- [x] hozaz.5 Define the `signed` app's promises and non-promises, including
  what the signature witnesses and what remains local trust judgment.
- [x] hozaz.6 Define the kernel's multi-app implementation promises: local
  app boundary, pCID dispatch, evidence records, local refusal, and unsupported
  behavior.
- [x] hozaz.7 Implement the directory layout under `implementations/poc3/` after
  DF/DI is locked.
- [x] hozaz.8 Add deterministic tests for app/kernel message flow, unsupported
  pCID refusal, echo semantics, signed-message evidence, and multi-app routing.
- [x] hozaz.9 Add a deterministic local or container demo command that Steve can
  run directly.
- [x] hozaz.10 Record final outcome: what worked, what was fake, what changed
  from `poc2`, and whether this evidence should update `DEV-GUIDE-RESOURCES.md`
  or `SIM-fovip`.

## Implementation notes

The first `poc3` implementation adds:

- A standalone Go module under `implementations/poc3/`.
- Importable `kernel`, `hello`, `echo`, `signed`, and `lib` packages with
  command wrappers under each package's `cmd/` directory.
- Per-app pCIDs for hello, echo, and signed messages plus a kernel-local
  receive-promise pCID.
- A two-container Compose demo where Alice and Bob each run one kernel plus all
  three app commands.
- Tests for envelope/proof-slot round-trip, signature-slot behavior,
  signed-envelope verification, and kernel delivery to a promised local
  receiver.

Validation run during implementation:

```sh
cd implementations/poc3
GOCACHE=/tmp/wire-lab-gocache go test ./...
bash -n scripts/run-alice.sh scripts/run-bob.sh
GOCACHE=/tmp/wire-lab-gocache errcheck ./...
docker compose config
```

Docker demo verification:

```sh
cd implementations/poc3
docker compose up --build --abort-on-container-exit
```

The first demo run exposed a bug: Bob's kernel had no peer address, so Bob's
echo response was interpreted as local delivery and produced broken-pipe
evidence. The implementation was corrected so Bob starts with `--peer
alice:9100`.

The corrected demo completed with both containers exiting `0`. Observed outcome:

- Bob's hello app judged Alice's hello message kept locally.
- Bob's echo app judged Alice's echo request kept locally and made a new echo
  response promise-message back to Alice.
- Alice's echo app judged Bob's echo response kept locally.
- Bob's signed app verified Alice's signed envelope bytes and judged the signed
  note kept locally.
- Both kernels recorded local receive, send, deliver, and receive-promise
  evidence without making app-level promise-status judgments.

What worked: `poc3` demonstrates two containers, one kernel per container, all
three apps per container, per-app pCIDs, a kernel-local receive-promise pCID,
and an envelope proof slot for the signed app.

What remains fake or intentionally narrow: deterministic demo keys, no durable
storage, no production peer identity, no real trust history, no app lifecycle
supervision, no queueing for late receive promises, and no final SDK/API claim.

Follow-up: this evidence is strong enough to consider when updating
`DEV-GUIDE-RESOURCES.md` or `SIM-fovip`, but those updates should be separate
tasks so they can distinguish executable POC evidence from final PromiseGrid
guidance.

## Acceptance criteria

- Multiple apps use the same app/kernel message boundary without kernel code
  becoming app-specific RPC dispatch.
- All app/kernel messages remain pCID-selected PromiseGrid envelopes unless a
  later locked decision changes the envelope shape.
- The `echo` behavior is framed as a local promise made by the echo app, not as
  Alice commanding Bob's app to respond.
- The `signed` behavior distinguishes exact-byte/proof evidence from local trust
  judgment.
- The POC remains small enough to read as a thought experiment made executable.
