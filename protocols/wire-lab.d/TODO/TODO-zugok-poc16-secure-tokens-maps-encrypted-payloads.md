# TODO-zugok: POC16 secure tokens, maps, and encrypted payloads

## Status

Planned. Owns the future `implementations/poc16-secure-tokens-maps-encrypted-payloads/`
proof of concept. POC16 must be a strict superset of POC15 and should add
security and payload-shape pressure before the constrained M4/LoRa work planned
for POC17. Source: `DI-ruvot`.

## Decision Intent Log

ID: DI-ruvot
Date: 2026-06-16 14:07:10
Status: active
Decision: Plan POC16 as a strict POC15 superset that permits pCID-owned CBOR maps, requires cryptographically secure tokens with CWT recommended, and adds encrypted payload coverage.
Intent: POC15 now has multihop forwarding, multiarity envelopes, message DAGs, COSE specimens, per-agent filesystem CAS, sparse peer storage, and bearer-token-shaped incentives. POC16 should preserve that baseline while testing the next protocol questions: when maps are useful for self-documenting payloads, how secure capability tokens should be encoded and verified, and how encrypted payloads affect pCID dispatch, CAS storage, parent links, proof placement, and local trust decisions.
Constraints: Preserve one top-level semantic action `promise`; preserve `grid([42(pCID), ...])`; do not regress POC15 route, DAG, WASM, stdio, shipping, CAS, GC, analyzer, or clean-run behavior; CBOR maps are permissible only when the pCID spec chooses them and remain discouraged for constrained/IoT protocols; all tokens, bearer or non-transferable, must be cryptographically secure; CWT is the recommended token format unless a TE/DI selects a narrower profile; encrypted payloads must not create global authorization, conformance, policy-enforcement, or trust-authority semantics.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md; protocols/wire-lab.d/TODO/TODO-komon-poc17-m4-lora-runtime.md; DEV-GUIDE-RESOURCES.md.

## Scope

- POC16 is executable design evidence, not production software and not a final
  PromiseGrid token, encryption, payload-map, or key-management standard.
- POC16 starts from POC15 behavior and adds new pressure without removing any
  POC15 acceptance gate unless a later explicit DI authorizes a scoped exception.
- POC16 should keep the current PromiseGrid principle that protocol objects help
  agents make, recognize, remember, and evaluate promises; they do not command
  other agents.
- POC16 should preserve local-trust-only semantics: every trust update, token
  acceptance, token refusal, decryption choice, map payload acceptance, and
  storage/GC choice is local to the observing agent.
- POC16 should make failure cases explicit: expired tokens, malformed CWTs,
  wrong audience, wrong issuer, revoked tokens, decryption failure, unknown key,
  unsupported map payloads, constrained-device map refusal, and encrypted
  payloads whose cleartext is unavailable to an intermediate relay.

## Superset Baseline

- Copy or evolve POC15 so the POC16 clean run still demonstrates:
  - multihop forwarding through voluntary route promises,
  - pCID-owned multiarity envelope slots,
  - parent links in envelope slots and payloads,
  - retained raw CBOR message artifacts and diagnostic rendering,
  - per-agent filesystem-backed sparse CAS stores,
  - per-agent partial message DAGs with missing-parent records,
  - peer CAS storage/retrieval promises and local sparse-store non-commitments,
  - local CAS retention/GC promises,
  - route economics and reusable/asymmetric route promises,
  - WASM and stdio runtime-adapter agents doing useful work,
  - shipping, storage, compute, verification, trust, replay, pressure, and
    migration coverage inherited from the POC chain.
- Add a POC16 analyzer rule that reports a regression if any POC15 acceptance
  category disappears without an explicit POC16 exception record.

## CBOR Map Payload Targets

- Add one or more pCID specs whose payload is a CBOR map because the protocol
  benefits from self-documenting fields, sparse optional fields, or human
  diagnostic readability.
- Add one or more pCID specs whose payload remains a compact CBOR array because
  the protocol is intended to interact with limited devices or very small
  parsers.
- Demonstrate that maps are permissible, not universal:
  - a map-accepting agent promises to parse a pCID-owned map payload,
  - a constrained-profile agent promises only an array payload for the same
    high-level work, or explicitly does not promise the map profile,
  - a relay forwards opaque encrypted or cleartext bytes without needing to know
    whether the payload is a map or array unless its own promise requires that.
- Add diagnostics showing map key choices, map canonicalization expectations,
  and how unknown optional keys are treated by the pCID spec.
- Avoid old prefixed compatibility maps as the design target. If compatibility
  projections remain in code, label them as runtime adapters, not protocol
  payload recommendations.

## Secure Token Targets

- Replace toy token strings with cryptographically secure capability-token
  objects for all bearer and non-transferable token flows.
- Prefer CWT-shaped tokens for POC16 unless a later TE/DI narrows the profile:
  issuer, subject/holder or bearer semantics, audience, expiration, not-before
  time, token identifier, capability body, confirmation or proof-of-possession
  material for non-transferable tokens, and signature/MAC metadata.
- Model bearer tokens as transferable promises by the issuer that a redeemer may
  present the token for a promised capability, subject to the issuer's local
  redemption promise terms.
- Model non-transferable tokens as holder-bound promises where the issuer can
  locally reject redemption if the presenter does not prove the expected holder
  relationship.
- Add secure-token failure cases:
  - invalid signature or MAC,
  - expired token that does not lower trust in the issuer,
  - revoked token whose status is judged by the issuer's local promise terms,
  - wrong audience or wrong pCID,
  - replay after a one-time redemption promise,
  - bearer token transfer that succeeds,
  - non-transferable token transfer attempt that is locally not promised.
- Keep token validation PromiseGrid-correct: a token is not global permission,
  authorization, or command authority; it is a signed promise artifact whose
  usefulness depends on the redeemer's and issuer's local judgments.

## Encrypted Payload Targets

- Add pCID specs whose envelope slot or payload slot carries encrypted payload
  bytes.
- Distinguish at least three cases:
  - end-to-end encrypted payload where relays can forward and store exact bytes
    but cannot inspect cleartext,
  - storage-encrypted payload where the CAS object CID names ciphertext and the
    cleartext CID is optional local metadata or a cleartext promise body,
  - recipient-specific encrypted payload where only selected promisees can
    decrypt.
- Demonstrate that pCID dispatch still works when only slot 0 is clear and the
  payload slot is encrypted.
- Demonstrate parent-link behavior for encrypted messages:
  - envelope parent links are visible to relays and DAG stores,
  - payload parent links are hidden unless the payload is decrypted,
  - agents record local missing-parent state without treating hidden links as a
    global consistency failure.
- Add encrypted payload failure cases:
  - wrong recipient key,
  - missing decryption key,
  - tampered ciphertext,
  - valid ciphertext under an unsupported pCID,
  - relay attempts to inspect cleartext and records a local non-commitment
    rather than claiming failure of the sender's promise.
- Keep encryption promise-based: Alice may promise that ciphertext is intended
  for Bob and shaped by a pCID; Bob locally decides whether he can decrypt,
  verify, store, respond, or update trust.

## Key And Proof Targets

- Add a minimal key-discovery and key-rotation path sufficient for secure tokens
  and encrypted payloads.
- Keep key rotation under an identity/key pCID rather than a generic
  observation protocol.
- Compare token-level proof, payload-level proof, envelope-level proof, and
  transport/session proof without prematurely declaring one universal answer.
- Record when a relay only verifies the outer envelope proof versus when an
  endpoint verifies a token or decrypts payload cleartext.
- Add negative tests for mismatched algorithm identifiers, unsupported key types,
  stale keys, and token/payload proof confusion.

## CAS, DAG, And Retention Targets

- Store secure tokens, encrypted payloads, and map payload specimens in each
  relevant agent's filesystem-backed CAS using the POC15 sparse-store model.
- Preserve exact-byte CIDs: encrypted objects are addressed by ciphertext CID
  unless a pCID explicitly defines a cleartext-derived reference.
- Add CAS metadata that distinguishes local byte format, encrypted status,
  token status, and visible parent-link status without creating a global schema.
- Add GC behavior for secure tokens and encrypted payloads:
  - expired token objects may be removed under local retention promises,
  - revoked-token records may be retained as local relationship history,
  - encrypted payloads may be retained without decryption keys,
  - agents may decline retention when token value or local storage budget is
    insufficient.

## Economics And Incentives Targets

- Extend route/storage/compute economics so secure bearer tokens can pay for or
  incentivize forwarding, storage, decryption assistance, verification, or
  compute.
- Demonstrate local exchange-rate movement when a token issuer keeps or breaks
  promises.
- Demonstrate opportunity cost: an agent may prefer a compact array protocol or
  decline a large map/encrypted payload because storage, CPU, decryption, or
  verification cost exceeds the promised value.
- Keep payment voluntary: a forwarding/storage/compute agent may promise service
  in exchange for a token, may decline an offered token, or may accept only from
  trusted issuers.

## Analyzer And Run-Review Targets

- Add analyzer score/report sections for:
  - POC15 superset preservation,
  - secure token validity and failure containment,
  - CWT-shaped token coverage,
  - bearer versus non-transferable semantics,
  - encrypted payload coverage,
  - pCID-owned map payload coverage,
  - constrained-profile array fallback,
  - local trust correctness,
  - imposition/authorization vocabulary regression,
  - CAS/DAG behavior for encrypted and token artifacts.
- Keep analyzer and collector passive. They may summarize retained artifacts for
  developer review but must not affect routing, trust, token redemption,
  decryption, or promise outcomes.
- Include raw-message diagnostic examples for:
  - a map payload,
  - an array payload for a constrained profile,
  - a CWT-shaped bearer token,
  - a holder-bound non-transferable token,
  - an encrypted payload with visible envelope parent links,
  - a failed token or decryption case.

## Documentation Targets

- Add `implementations/poc16-secure-tokens-maps-encrypted-payloads/README.md`
  explaining the POC16 scope, superset baseline, and clean-run commands.
- Add POC16 protocol notes describing:
  - when maps are appropriate,
  - why maps are discouraged for limited devices,
  - secure token shapes and CWT recommendation,
  - encrypted payload cases,
  - pCID dispatch with encrypted payloads,
  - exact-byte CAS semantics for ciphertext,
  - local-trust-only interpretation of tokens and decryption outcomes.
- Update `DEV-GUIDE-RESOURCES.md` after implementation so guide authors can
  distinguish POC16 evidence from final PromiseGrid APIs.
- Cross-reference POC17 so the future M4/LoRa runtime inherits POC16 lessons
  but can reject map-heavy or expensive token/encryption profiles when device
  constraints require compact alternatives.

## Tasks

- [ ] zugok.1 Scaffold `implementations/poc16-secure-tokens-maps-encrypted-payloads/` from the POC15 baseline.
- [ ] zugok.2 Add an analyzer superset gate proving POC15 acceptance categories still appear in POC16 clean runs.
- [ ] zugok.3 Add pCID-owned CBOR map payload specimens and at least one normal map-based protocol flow.
- [ ] zugok.4 Add pCID-owned compact CBOR array alternatives for constrained-profile agents.
- [ ] zugok.5 Add map/array negotiation or local non-commitment behavior without global conformance semantics.
- [ ] zugok.6 Replace toy bearer tokens with cryptographically secure bearer tokens.
- [ ] zugok.7 Replace toy non-transferable tokens with holder-bound cryptographically secure tokens.
- [ ] zugok.8 Prefer a CWT-shaped token format and document any POC-local simplifications.
- [ ] zugok.9 Add token validation failures for invalid proof, expiry, revocation, wrong audience, replay, and transfer mismatch.
- [ ] zugok.10 Add encrypted payload pCIDs for end-to-end, storage-encrypted, and recipient-specific cases.
- [ ] zugok.11 Preserve pCID dispatch when payload slots are encrypted.
- [ ] zugok.12 Add encrypted-parent-link cases covering visible envelope parents and hidden payload parents.
- [ ] zugok.13 Add decryption failure and unsupported encrypted-pCID containment cases.
- [ ] zugok.14 Add minimal key discovery and key rotation for token and payload use.
- [ ] zugok.15 Store token, map, array, and encrypted payload artifacts in per-agent filesystem CAS.
- [ ] zugok.16 Add CAS/DAG metadata and analyzer checks for ciphertext CID semantics and sparse hidden-parent behavior.
- [ ] zugok.17 Add GC and retention behavior for expired tokens, revoked-token history, and encrypted payloads.
- [ ] zugok.18 Extend route/storage/compute economics to use secure bearer-token incentives.
- [ ] zugok.19 Add exchange-rate and opportunity-cost events affected by token issuer promise history.
- [ ] zugok.20 Add raw CBOR diagnostic examples for map, array, CWT-shaped token, holder-bound token, encrypted payload, and failures.
- [ ] zugok.21 Update POC16 README and protocol notes.
- [ ] zugok.22 Update `DEV-GUIDE-RESOURCES.md` after executable POC16 evidence exists.
- [ ] zugok.23 Reconcile POC16 lessons into the POC17 M4/LoRa plan before implementing constrained-device token/encryption behavior.
- [ ] zugok.24 Run Go validation, errcheck, and clean POC16 containers after implementation.
