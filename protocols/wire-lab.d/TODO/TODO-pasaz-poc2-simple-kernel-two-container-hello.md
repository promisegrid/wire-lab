# TODO-pasaz: poc2 simple kernel two-container hello

## Prior aliases

This TODO was minted directly under the post-`TE-mumuv` proquint handle scheme.
No integer or timestamp alias exists.

## Status

Planned. `poc2` is the second PromiseGrid proof-of-concept system and is
separate from the external `grid-poc` repo. The immediate target is a minimal
two-container executable proof that app/kernel and kernel/kernel boundaries can
both use pCID-selected `grid([42(pCID), payload, ...])` messages without
turning the kernel into an RPC authority.

## Decision Intent Log

### DI-ratij

ID: DI-ratij
Date: 2026-05-25 16:04:35
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Create `poc2` as the second PromiseGrid proof-of-concept system,
separate from the external `grid-poc` repo, under `implementations/poc2/`.
Intent: Build a small executable proof that app/kernel and kernel/kernel
boundaries can both use pCID-selected `grid([42(pCID), payload, ...])` messages
without turning the kernel into an RPC authority.
Constraints: Promise Theory first; no permissions, authorization,
conformance-authority, or command-and-control framing; no final API claims;
`pCID` means Protocol CID; exact message bytes are evidence; use two Docker
containers named for Alice and Bob; use Go unless a later DI changes
implementation language; do not run provider-backed GA or canary work as part of
this implementation task.
Affects: `implementations/poc2/`; `protocols/wire-lab.d/TODO/TODO.md`;
`DEV-GUIDE-RESOURCES.md` if the result changes guide evidence.

## Summary

Build `poc2`, a minimal two-container kernel/app proof of concept.

`poc2` demonstrates:

- Alice's app sends a local PromiseGrid hello message to Alice's kernel.
- Alice's kernel forwards the same promise-message semantics to Bob's kernel.
- Bob's kernel delivers the message to Bob's local hello app.
- Both kernels record local evidence of receive, send, deliver, and refusal
  events.
- Unsupported pCIDs are locally refused with evidence, not globally rejected by
  authority.

This is a proof of concept, not the final PromiseGrid kernel, SDK, transport,
identity, storage, or security model.

## Locked implementation decisions

- **Architecture:** two Docker containers, `alice` and `bob`; each container
  runs one local kernel process and one local hello app process.
- **Boundary shape:** all app/kernel and kernel/kernel promise boundaries use
  CBOR `grid([42(pCID), payload, ...])` messages.
- **Transport:** TCP inside the Compose network for kernel/kernel traffic and
  localhost TCP for app/kernel traffic.
- **Language:** Go, matching the repo's current executable tooling shape.
- **Path:** implementation lives under `implementations/poc2/`.
- **Promise Theory constraint:** kernels route, record, deliver, and refuse
  local promises; they do not command apps, command peers, or act as trust
  authorities.

## Subtasks

- [ ] pasaz.1 Create `implementations/poc2/` with `README.md`, `CHANGELOG.md`,
  `go.mod`, Docker/Compose files, and a small Go implementation.
- [ ] pasaz.2 Define the POC's minimal kernel implementation promises:
  supported pCIDs, unsupported-pCID behavior, app-facing send/receive promises,
  kernel-to-kernel transport promise, local evidence log, and explicit
  non-promises.
- [ ] pasaz.3 Add a minimal POC hello protocol spec document whose pCID is the
  Protocol CID for the hello payload shape.
- [ ] pasaz.4 Encode all app/kernel and kernel/kernel boundary messages as CBOR
  `grid([42(pCID), payload, ...])`; use a small CBOR library or tightly scoped
  encoder rather than inventing a different wire shape.
- [ ] pasaz.5 Implement `poc2-kernel`: local app listener, remote peer listener,
  pCID dispatcher, exact-byte logging, unsupported-pCID refusal, and app
  delivery.
- [ ] pasaz.6 Implement `poc2-hello`: app process that connects to the local
  kernel, sends or receives a hello promise-message, and prints received
  messages.
- [ ] pasaz.7 Run two Docker containers, `alice` and `bob`, on one Compose
  network; each container starts one local kernel and one hello app process.
- [ ] pasaz.8 Make the demo command deterministic:
  `cd implementations/poc2 && docker compose up --build --abort-on-container-exit`.
- [ ] pasaz.9 Ensure expected output clearly shows Alice app -> Alice kernel ->
  Bob kernel -> Bob app, plus local evidence records on both sides.
- [ ] pasaz.10 Add tests for envelope round-trip, pCID dispatch,
  unsupported-pCID refusal, local app/kernel delivery, remote kernel/kernel
  delivery, and evidence-log content.
- [ ] pasaz.11 Update `DEV-GUIDE-RESOURCES.md` only if the POC produces evidence
  worth citing as provisional kernel/app-boundary guidance.
- [ ] pasaz.12 Record final outcome in this TODO: what worked, what was fake,
  what remains unresolved, and whether `DR-davod` or `SIM-fovip` should be
  updated.

## Minimum interfaces

- `poc2-kernel`
  - `--node alice|bob`
  - `--app-listen 127.0.0.1:<port>`
  - `--peer-listen 0.0.0.0:<port>`
  - `--peer bob:<port>` for Alice's outbound path
  - writes newline-delimited evidence records to stdout and optionally
    `./run/evidence.jsonl`
- `poc2-hello`
  - `--node alice|bob`
  - `--kernel 127.0.0.1:<port>`
  - `--mode send|receive`
  - `--to bob`
  - `--text "hello from Alice"`
- Message boundary:
  - outer shape: CBOR `grid([42(pCID), payload, ...])`
  - slot 0: tag-42 Protocol CID for the POC hello spec
  - slot 1: protocol-owned payload bytes
  - slots 2..N: absent in v0 unless the hello spec explicitly defines one

## Acceptance criteria

- `go test ./...` passes under `implementations/poc2/`.
- Docker demo completes without manual timing hacks.
- Bob's app prints the hello text sent by Alice's app.
- Alice and Bob evidence logs include exact-byte hashes for local app send,
  kernel forward, remote receive, app delivery, and any refusal.
- The code and docs do not describe the kernel as commanding apps or peers.
- The README explicitly says `poc2` is separate from `grid-poc` and is not a
  final PromiseGrid API.

## Assumptions

- Use Go for the POC because this repo's executable tooling is Go-shaped and
  `implementations/README.md` expects local reference implementations.
- Use two containers, not four: each container runs a local kernel process plus
  a local hello app process.
- Use TCP inside the Compose network for kernel/kernel transport and localhost
  TCP for app/kernel transport.
- No cryptographic identity or production signature scheme in v0; the POC
  records exact bytes and promise evidence but does not claim secure peer
  identity.
- Existing TE/SIM evidence is sufficient for this TODO; no new TE is required
  before building unless implementation exposes a new semantic decision.
