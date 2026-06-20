# TODO-milov: POC16 RFC-complete specdocs

## Status

Implemented and validated by the current specdoc rewrite task. This TODO exists
because POC16's embedded specdocs are pCID inputs, LLM-agent prompt context, and
developer-facing protocol specifications; short notes are not sufficient for
production-shaped implementation.

## Decision Intent Log

ID: DI-bitug
Date: 2026-06-19 20:27:21
Status: active
Decision: Rewrite every POC16 embedded Markdown specdoc as a complete RFC-style protocol specification, including message shape, expected behavior, payload grammar, Promise Theory semantics, security considerations, examples, and ASCII protocol state-machine diagrams where the protocol has meaningful lifecycle states.
Intent: A pCID names the content of a protocol specification. If the specdoc is too thin, developers and LLM-backed agents cannot implement or evaluate the protocol from the pCID's named document, and the pCID becomes a label rather than a complete protocol promise. POC16 must therefore make each specdoc stand on its own as production-shaped implementation guidance while still marking POC16 itself as executable design evidence rather than a final PromiseGrid standard.
Constraints: Cover all Markdown files embedded by `specdocs/*.md`, including active, specimen/profile, and retired/historical pCIDs; preserve `grid([42(pCID), ...protocol-defined-slots])`; keep pCID as protocol spec selector, not destination address; keep trust local and promise-first; do not introduce global authority, conformance, permission, or command semantics; acknowledge that editing spec bytes changes derived pCIDs; keep runtime behavior unchanged unless validation exposes a mismatch.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/specdocs/*.md; implementations/poc16-secure-tokens-maps-encrypted-payloads/specdocs/specdocs_test.go; protocols/wire-lab.d/TODO/TODO-milov-poc16-rfc-complete-specdocs.md; protocols/wire-lab.d/TODO/TODO.md.

## Scope

- [x] milov.1 Rewrite every active POC16 specdoc as an RFC-style protocol spec.
- [x] milov.2 Rewrite every specimen/profile specdoc as a complete profile spec.
- [x] milov.3 Rewrite every retired/historical specdoc with complete historical behavior and migration notes.
- [x] milov.4 Include ASCII protocol state-machine diagrams where lifecycle states are meaningful.
- [x] milov.5 Add a completeness regression test for required specdoc sections.
- [x] milov.6 Run `go test ./...`, `errcheck ./...`, and the POC16 clean regression.
- [x] milov.7 Commit the specdoc rewrite separately from the parser-role checkpoint.

## Required Spec Structure

- Title, status, version, and pCID derivation note.
- Abstract and intended use.
- Promise Theory model: promisers, promisees, reciprocal promises, local trust, and no promises on behalf of another agent.
- Envelope shape and slot semantics.
- Payload grammar with CBOR diagnostic notation and field tables.
- Sender behavior, receiver/parser behavior, and malformed/non-commitment behavior.
- Protocol state machine, when the protocol has lifecycle states.
- State, CAS/DAG, replay, expiry, retention, and garbage-collection behavior when applicable.
- Security considerations, interoperability notes, and examples.
