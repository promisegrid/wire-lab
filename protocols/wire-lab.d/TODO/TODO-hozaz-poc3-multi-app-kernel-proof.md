# TODO-hozaz: poc3 multi-app kernel proof

## Prior aliases

This TODO was minted directly under the post-`TE-mumuv` proquint handle scheme.
No integer or timestamp alias exists.

## Status

Planned. `poc3` follows `poc2` and should test the next level of kernel/app
shape: several local apps sharing one PromiseGrid-style kernel boundary without
turning the kernel into an RPC authority or app-specific service registry.

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

## Subtasks

- [ ] hozaz.1 Review `poc2` evidence and decide which code patterns can be
  reused without carrying over accidental single-app assumptions.
- [ ] hozaz.2 Run DF for exact command names, package/type names, protocol
  surfaces, and runtime topology before creating code.
- [ ] hozaz.3 Define the `hello` app's promises and non-promises.
- [ ] hozaz.4 Define the `echo` app's promises and non-promises, including how
  echoing differs from obeying a remote command.
- [ ] hozaz.5 Define the `signed` app's promises and non-promises, including
  what the signature witnesses and what remains local trust judgment.
- [ ] hozaz.6 Define the kernel's multi-app implementation promises: local
  app boundary, pCID dispatch, evidence records, local refusal, and unsupported
  behavior.
- [ ] hozaz.7 Implement the directory layout under `implementations/poc3/` after
  DF/DI is locked.
- [ ] hozaz.8 Add deterministic tests for app/kernel message flow, unsupported
  pCID refusal, echo semantics, signed-message evidence, and multi-app routing.
- [ ] hozaz.9 Add a deterministic local or container demo command that Steve can
  run directly.
- [ ] hozaz.10 Record final outcome: what worked, what was fake, what changed
  from `poc2`, and whether this evidence should update `DEV-GUIDE-RESOURCES.md`
  or `SIM-fovip`.

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
