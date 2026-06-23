# TODO-zugok: POC16 secure tokens, maps, and encrypted payloads

## Status

Planned. Owns the future `implementations/poc16-secure-tokens-maps-encrypted-payloads/`
proof of concept. POC16 must be a strict superset of POC15, should add
security and payload-shape pressure before the constrained M4/LoRa work planned
for POC17, and should make the pCID-selected parser/builder kernel-role split
explicit. Source: `DI-ruvot`; `DI-mubul`; `DI-nogij`; `DI-rigup`; `DI-vulit`.
The post-implementation parser-role correction is owned by `TODO-sosoj` /
`DI-gazin`; that follow-up makes `TE-ritig`'s parser model 3 executable by
inserting a real parser-role process between apps and the transport kernel.
`DI-magug` supersedes the stale root-level `docs/protocols/` path portions of
the earlier DIs; current POC16 protocol specs live under the implementation-local
`implementations/poc16-secure-tokens-maps-encrypted-payloads/docs/protocols/`
source of truth.

## Decision Intent Log

ID: DI-ruvot
Date: 2026-06-16 14:07:10
Status: active
Decision: Plan POC16 as a strict POC15 superset that permits pCID-owned CBOR maps, requires cryptographically secure tokens with CWT recommended, and adds encrypted payload coverage.
Intent: POC15 now has multihop forwarding, multiarity envelopes, message DAGs, COSE specimens, per-agent filesystem CAS, sparse peer storage, and bearer-token-shaped incentives. POC16 should preserve that baseline while testing the next protocol questions: when maps are useful for self-documenting payloads, how secure capability tokens should be encoded and verified, and how encrypted payloads affect pCID dispatch, CAS storage, parent links, proof placement, and local trust decisions.
Constraints: Preserve one top-level semantic action `promise`; preserve `grid([42(pCID), ...])`; do not regress POC15 route, DAG, WASM, stdio, shipping, CAS, GC, analyzer, or clean-run behavior; CBOR maps are permissible only when the pCID spec chooses them and remain discouraged for constrained/IoT protocols; all tokens, bearer or non-transferable, must be cryptographically secure; CWT is the recommended token format unless a TE/DI selects a narrower profile; encrypted payloads must not create global authorization, conformance, policy-enforcement, or trust-authority semantics.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md; protocols/wire-lab.d/TODO/TODO-komon-poc17-m4-lora-runtime.md; DEV-GUIDE-RESOURCES.md.

ID: DI-mubul
Date: 2026-06-18 14:25:44
Status: superseded
Decision: Plan POC16 so the kernel dispatches normal application messages only by slot-0 pCID, parses payloads only for explicitly kernel-handled pCIDs, stores real protocol spec docs under `docs/protocols/`, freezes each spec with a `docs/protocols/<cidv1-encoded>.md` symlink, and embeds relevant spec prose into LLM-agent prompt context with `go:embed`.
Intent: The pCID already names the protocol spec that defines payload shape, slot meanings, proof semantics, and promise vocabulary. POC16 should stop treating map-vs-array as a kernel negotiation problem, should prove that handler/app code owns pCID-specific payload decoding, and should give LLM agents the same spec prose that deterministic handlers are written against.
Constraints: Do not make runtime kernel routing depend on reading spec prose. Do not let the kernel parse application payloads merely to discover whether they are maps or arrays. Do not create global conformance semantics for malformed or unsupported payloads. LLM agents must receive exact embedded spec text for every pCID they promise to send, receive, redeem, verify, store, compute, or route. Spec docs remain provenance and prompt-context artifacts, not central registry authority.
Affects: docs/protocols/; implementations/poc16-secure-tokens-maps-encrypted-payloads/; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md; DEV-GUIDE-RESOURCES.md.

ID: DI-nogij
Date: 2026-06-19 06:33:32
Status: active
Decision: File TE-ritig before POC16 implementation and supersede DI-mubul's over-strong kernel-dispatch wording with a parser/builder kernel-role plan: slot 0 selects a pCID-specific parser or builder role, while pCID-defined payloads or nested payloads carry app, operation, destination, route, and local-addressing semantics.
Intent: pCID should remain the protocol spec selector, not a destination address, app name, RPC method, service registry entry, or universal routing key. POC16 should test whether production-shaped PromiseGrid kernels are better described as cooperating roles--transport listener, pCID router, parser, builder, CAS, key, hardware, resource allocator, and app adapter--rather than as one monolithic dispatcher that parses every application payload.
Constraints: Preserve `grid([42(pCID), ...protocol-defined-slots])`; preserve `docs/protocols/` spec docs and `go:embed` LLM spec context from DI-mubul; do not implement POC16 until DF resolves TE-ritig's surviving alternatives; do not require a universal `to` field or generic routing fact returned to the transport listener; model ACKs, errors, backpressure, malformed input, unsupported pCIDs, and resource allocation as local promises unless a pCID spec explicitly defines a wire-visible promise.
Affects: docs/thought-experiments/TE-ritig-pcid-cardinality-parser-builder-kernel-roles.md; docs/thought-experiments/README.md; protocols/wire-lab.d/specs/harness-spec-draft.md; docs/protocols/README.md; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md; protocols/wire-lab.d/TODO/TODO.md; DEV-GUIDE-RESOURCES.md; future implementations/poc16-secure-tokens-maps-encrypted-payloads/.
Supersedes: DI-mubul

ID: DI-rigup
Date: 2026-06-19 06:50:38
Status: active
Decision: POC16 must be a strict superset of POC15 in named agents and executable functionality, not merely a superset of abstract acceptance categories.
Intent: POC16 is adding secure tokens, maps, encrypted payloads, protocol spec docs, embedded LLM spec context, and parser/builder role pressure, but those additions must not hide regressions by dropping POC15's actual agent population or working behaviors. The POC chain is useful only if each successor preserves the concrete prior-world pressure: relationship agents, adversarial agents, shipping/device/system agents, WASM and stdio runtime agents, multihop route agents, CAS/compute/storage agents, persistent sessions, sparse CAS, raw-message DAGs, route economics, signed capability-token bytes, and analyzer coverage.
Constraints: Preserve all POC15 named agents unless a later scoped non-superset DI explicitly lists and justifies a dropped agent; preserve all POC15 functionality unless a later scoped non-superset DI explicitly lists and justifies a dropped behavior; container grouping, process shape, parser/builder role split, and implementation internals may change only if the clean-run analyzer proves equivalent or stronger agent/function coverage; do not weaken the POC superset discipline in `DI-sinur`.
Affects: protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md; protocols/wire-lab.d/TODO/TODO.md; DEV-GUIDE-RESOURCES.md; future implementations/poc16-secure-tokens-maps-encrypted-payloads/.

ID: DI-vulit
Date: 2026-06-19 06:54:32
Status: active
Decision: Implement POC16 from the POC15 baseline using TE-ritig's default conclusion: moderate protocol-family/major-version pCIDs, slot 0 selecting pCID-specific parser/builder kernel roles, and pCID-defined payloads or nested payloads carrying operation, destination, route, economic, and local-addressing semantics.
Intent: The POC16 implementation needs a concrete locked default before code changes begin. This keeps the transport listener from becoming an RPC dispatcher or universal address parser while still letting POC16 test secure tokens, CBOR maps, encrypted payloads, protocol spec docs, embedded LLM spec context, and POC15 superset preservation in one executable run.
Constraints: Preserve all `DI-rigup` POC15 named-agent and functionality requirements; preserve one top-level semantic action `promise`; preserve `grid([42(pCID), ...protocol-defined-slots])`; do not add pCID-per-operation fragmentation as the default; do not require a universal `to` field; keep parser/builder, ACK/error, and backpressure events as local implementation promises unless a pCID spec explicitly defines a wire-visible promise; keep analyzer and collector passive.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/; docs/protocols/; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md; DEV-GUIDE-RESOURCES.md.

ID: DI-titik
Date: 2026-06-19 23:21:16
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC16 supervisors must use the configured `shutdown_grace_millis` value as the SIGTERM grace window for kernel, parser-role, and app child processes instead of a hard-coded shorter timeout.
Intent: POC16 clean runs emit many parser/kernel/app terminal session records at shutdown. A shorter supervisor kill window can terminate children before those local lifecycle promises are closed and forwarded to the passive collector, causing false analyzer failures even when the protocol traffic completed. The configured shutdown grace is already part of the POC16 runtime contract and monitor wait budget, so the supervisor should honor that same value.
Constraints: Keep the collector passive; do not weaken analyzer terminal-session gates; do not use observer files for agent coordination; preserve the existing process-per-role architecture; treat this as shutdown accounting and validation plumbing, not protocol behavior.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/cmd/poc16-supervisor/main.go; implementations/poc16-secure-tokens-maps-encrypted-payloads/config.json; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md.

ID: DI-jojoj
Date: 2026-06-21 10:40:07
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC16 kernel peer-session setup must not close duplicate or late-started persistent sessions while holding the kernel routing mutex.
Intent: Clean shutdown must be able to snapshot app, outbound-peer, and inbound-peer sessions and emit terminal records for every opened stream. Closing a TCP-backed persistent session while holding `kernel.mu` is unnecessary and can couple transport close behavior to routing-table locks, which makes late duplicate sessions and shutdown races harder to reason about.
Constraints: Do not weaken analyzer terminal-session gates; do not treat transport lifecycle accounting as peer trust evidence; preserve persistent TCP reuse keyed by endpoint; preserve shutdown rejection of sessions opened after the kernel is stopping.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/kernel/kernel.go; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md.

ID: DI-nuhit
Date: 2026-06-21 10:46:34
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC16 kernel shutdown must wait for accepted app/peer handler goroutines to drain using a bounded drain timer that is independent of the already-canceled run context.
Intent: The run context is canceled before shutdown starts, so using that same context as a handler-drain stop condition lets the kernel return immediately, close its log, and let the process exit while handler goroutines are still emitting persistent-session terminal records. Clean-run lifecycle accounting requires the kernel to close sessions and then give handlers a bounded chance to observe those closures and finish.
Constraints: Do not make shutdown unbounded; do not weaken analyzer terminal-session gates; keep drain records as local transport lifecycle accounting, not peer trust evidence; preserve supervisor-level SIGTERM grace as the outer kill bound.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/kernel/kernel.go; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md.

ID: DI-vazoz
Date: 2026-06-22 09:24:56 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC16 container supervisors must shut child roles down in phases: app processes finish first, the parser-role process closes its parser/kernel and parser/app sessions second, and the transport kernel closes listeners, app sessions, peer sessions, and handler drains last.
Intent: Parser roles are the local app-facing kernel role. If the supervisor cancels parser and transport-kernel children at the same time, busy containers can truncate one side of a persistent parser/kernel or kernel peer stream before terminal records are emitted. A phased shutdown preserves the local process model while letting each role close the streams it owns in dependency order.
Constraints: Do not weaken analyzer terminal-session gates; do not add observer-volume coordination; keep supervisor shutdown bounded by `shutdown_grace_millis`; preserve the process-per-role architecture and local Promise Theory semantics.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/cmd/poc16-supervisor/main.go; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md.

ID: DI-kiduj
Date: 2026-06-22 09:34:11 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC16 adversarial unknown-pCID probes must be recorded as expected local non-commitments and must not abort Mallory's remaining adversarial startup probes.
Intent: Unknown protocol CIDs are precisely the case where a receiver or parser role may decline to promise a parse, ACK, or delivery. Treating that expected non-commitment as a fatal startup error suppresses the later corrupt-CAS, bad-proof, key-rotation, capacity-refusal, and repair-promise probes that make the clean-run regression meaningful.
Constraints: Do not weaken unknown-pCID analyzer coverage; do not reinterpret unknown pCIDs as globally invalid; keep the event as local non-commitment, not enforcement or authority; preserve the rest of Mallory's startup workflow so later probes still exercise receiver autonomy.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/runtime/node.go; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md.

ID: DI-kunad
Date: 2026-06-22 09:59:53 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC16 persistent-session shutdown must emit local closed and terminal lifecycle records before invoking the underlying TCP close.
Intent: The analyzer is checking local lifecycle accounting, not global transport authority. A busy or slow socket close must not hide the local fact that the owning session promised to close, notified pending requests, and reached a terminal state during bounded shutdown.
Constraints: Do not weaken analyzer terminal-session gates; do not treat session lifecycle records as peer trust evidence; still attempt the underlying TCP close and record close failures; preserve exactly one terminal record per persistent session.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/transport/persistent_session.go; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md.

ID: DI-pajih
Date: 2026-06-22 10:15:53 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC16 parser-role shutdown must close and drain app-facing parser sessions before closing the parser/kernel control session.
Intent: Parser roles bridge local app promises to the transport kernel. During shutdown, app-facing sessions may still be delivering ACKs or local non-commitments through the parser role. Closing and draining those local app sessions first removes concurrent parser activity before the parser records terminal lifecycle accounting for its kernel control stream.
Constraints: Do not weaken analyzer terminal-session gates; keep shutdown bounded by supervisor grace and existing app-session close behavior; preserve the parser role as the owner of pCID payload parsing, not a transport-kernel RPC dispatcher.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/parserrole/parserrole.go; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md.

ID: DI-nuriv
Date: 2026-06-22 14:25:38 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: POC16 persistent-session close must bound the underlying TCP close after local terminal lifecycle records have been emitted.
Intent: Clean-run shutdown closes many app, parser, inbound-peer, and outbound-peer streams. One slow or wedged socket close must not prevent the owning process from emitting terminal lifecycle records for later sessions. The local session terminal record is the analyzer-relevant promise that the owner reached a terminal state; TCP close is still attempted and failures or deferred close are recorded, but it must not serialize all other session accounting behind one stream.
Constraints: Do not weaken analyzer terminal-session gates; preserve exactly one terminal record per persistent session; still attempt underlying TCP close; keep deferred-close records as transport lifecycle accounting, not peer trust evidence; do not add observer-volume coordination or global monitoring semantics.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/transport/persistent_session.go; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md.

ID: DI-mapah
Date: 2026-06-23 11:12:52 PDT
Status: active
Author: stevegt@t7a.org (Steve Traugott)
Decision: Replace POC16 flexible array-of-pairs body payloads with nested CBOR map bodies, without backwards parsing for old pair-list bodies.
Intent: If a flexible body already carries text keys, an array of `[key, value]` pairs has little benefit over a CBOR map and creates a design hazard by encouraging flattened Go projections where body keys can overwrite core promise fields such as `from`, `to`, `promise_about`, or derived `payload_protocol`. Flexible pCID-owned protocols should keep core promise slots positional while putting keyed details in a separate nested CBOR map namespace. Constrained protocols that need compactness should use pCID-specific positional body arrays instead.
Constraints: Preserve `grid([42(pCID), ...])`; preserve pCID-owned payload semantics; do not parse historical pair-list bodies; reject reserved/core body-map keys even though the body is nested; keep runtime compatibility projections local to parser/runtime edges; update POC16 protocol specs and the POC17 handoff TODO so future constrained-device work does not copy the old pair-list shape.
Affects: implementations/poc16-secure-tokens-maps-encrypted-payloads/protocol/; implementations/poc16-secure-tokens-maps-encrypted-payloads/docs/protocols/; protocols/wire-lab.d/TODO/TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md; protocols/wire-lab.d/TODO/TODO-komon-poc17-m4-lora-runtime.md.

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
- POC16 should preserve pCID-selected parser/builder boundaries: slot 0 selects
  the protocol-specific parser or builder role, while pCID-defined payloads or
  nested payloads carry app, operation, destination, route, and local-addressing
  semantics.
- POC16 should make failure cases explicit: expired tokens, malformed CWTs,
  wrong audience, wrong issuer, revoked tokens, decryption failure, unknown key,
  unsupported map payloads, constrained-device map refusal, and encrypted
  payloads whose cleartext is unavailable to an intermediate relay.

## Superset Baseline

- Copy or evolve POC15 so the POC16 clean run preserves the POC15 named-agent
  population unless a later scoped non-superset DI explicitly authorizes a
  named exception:
  - Alice, Bob, Carol, Dave, Ellen, Frank, Grace, Heidi, Ivan, Judy, Mallory,
    Oscar,
  - Fulfillment, Postal Scale, UPS Label Printer, Printer Port, Accounting,
  - Peggy and Victor.
- Preserve POC15's agent roles even if container grouping or process internals
  change:
  - privacy-sensitive sender,
  - storage provider,
  - compute provider,
  - skeptical local accountant,
  - relationship broker,
  - relay helper,
  - market maker,
  - repair-oriented operator,
  - constrained edge peer,
  - event-forwarding peer,
  - adversarial peers,
  - shipping workflow coordinator,
  - device/system agents,
  - local hardware-resource kernel-role app,
  - WASM runtime-adapter agent,
  - stdio-only runtime-adapter agent.
- Preserve POC15's executable functionality, not just category names, so the
  POC16 clean run still demonstrates:
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
    migration coverage inherited from the POC chain,
  - persistent multiplexed app/kernel and kernel/kernel TCP sessions,
  - signed CBOR/COSE capability-token bytes for storage incentives,
  - route expiry, renewal, transit-exclusion, trust-driven session
    reconfiguration, and peer CAS retrieval outcomes.
- Add POC16 analyzer rules that report a regression if:
  - any POC15 named agent disappears without an explicit POC16 exception record,
  - any POC15 agent role disappears without an explicit POC16 exception record,
  - any POC15 executable functionality disappears without an explicit POC16
    exception record,
  - POC16 adds parser/builder, token, map, or encryption pressure by removing
    older route, shipping, CAS, compute, WASM, stdio, persistent-session,
    sparse-CAS, raw-message, or analyzer behavior.

## Protocol Specification Source Targets

- Keep POC16 protocol spec docs that are meant to be hashed into pCIDs under
  `implementations/poc16-secure-tokens-maps-encrypted-payloads/docs/protocols/`.
- Keep editable draft specs at human-readable slug paths such as
  `implementations/poc16-secure-tokens-maps-encrypted-payloads/docs/protocols/route-v1.md`.
- Do not maintain a root-level `docs/protocols/` POC16 mirror or stale CIDv1
  symlink corpus; POC16 uses embedded spec bytes as the pCID source of truth.
- Add validation that every POC16 pCID has a corresponding embedded local
  protocol spec and that runtime pCIDs derive from those exact bytes.
- Treat spec docs as provenance, review, diagnostics, and LLM prompt-context
  sources. Deterministic runtime routing may use compiled/static pCID registry
  data and does not need to read prose specs during message dispatch.

## Kernel Dispatch And Handler Targets

- Treat slot 0 as the protocol-family selector, not as an app address, message
  kind, RPC method, service registry entry, or universal routing key.
- Route slot-0 pCIDs to pCID-specific parser/builder roles. Those roles may be
  processes, objects, WASM modules, firmware functions, or app-owned adapters
  depending on the runtime.
- Keep application pCID payload shape and local routing semantics under the
  parser/builder role selected by the pCID. The role already knows from the pCID
  spec whether to expect a CBOR map, CBOR array, COSE object, encrypted bytes,
  nested selector, app-local address, DID, key fingerprint, CID-rooted path,
  route promise, or another protocol-owned shape.
- Do not negotiate map-vs-array payload shape in the kernel. Instead, use local
  non-commitment and malformed-input behavior: if a receiver has not promised
  the pCID, or the bytes do not match that pCID's spec, record a local event
  without claiming global conformance failure.
- Add analyzer checks that flag any transport-listener branch that chooses app
  behavior by parsing arbitrary application payloads, assuming a universal `to`
  field, or treating pCID as an address rather than selecting the parser/builder
  role.

## Required TE Before POC16 Implementation

- TE-ritig tests pCID cardinality and parser/builder kernel roles before POC16
  implementation.
- Use TE-ritig's surviving alternatives to lock DF choices for:
  - whether POC16 uses moderate protocol-family pCIDs, a very small broad-pCID
    pressure test, or a layered outer/inner pCID model,
  - which pCID-specific parser/builder roles exist in POC16,
  - which ACKs and errors are local implementation events versus wire-visible
    pCID-defined promises,
  - how transport, parser, builder, CAS, key, hardware, resource allocation, and
    app roles exchange backpressure and non-commitment promises.
- Do not start POC16 implementation until those DF choices are locked in this
  TODO's Decision Intent Log.

## LLM Agent Spec Context Targets

- Use `go:embed` to bundle the relevant protocol spec docs into the POC16
  executable context builder for LLM-backed agents.
- Include embedded spec prose in each LLM agent's prompt context for every pCID
  the agent promises to send, receive, redeem, verify, store, compute, or route.
- Record, in run logs, which spec CIDs were supplied to each LLM agent context
  so the analyzer and later human review can verify prompt provenance.
- Validate that embedded spec bytes hash to the expected pCID before an LLM
  agent is allowed to rely on that context.
- Keep deterministic handlers code-driven, but tie their tests and comments to
  the same spec docs so deterministic and LLM behavior share one pCID source.

## CBOR Map Payload Targets

- Add one or more pCID specs whose payload is a CBOR map because the protocol
  benefits from self-documenting fields, sparse optional fields, or human
  diagnostic readability.
- Add one or more pCID specs whose payload remains a compact CBOR array because
  the protocol is intended to interact with limited devices or very small
  parsers.
- Demonstrate that maps are permissible, not universal:
  - a map-accepting agent promises to parse a pCID-owned map payload,
  - a constrained-profile agent promises only a distinct array-shaped pCID for
    the same high-level work, or explicitly does not promise the map-shaped pCID,
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
  - POC15 named-agent and agent-role preservation,
  - POC15 executable-function preservation,
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
  - the POC15 named-agent and executable-functionality baseline that POC16
    preserves,
	  - the implementation-local `docs/protocols/` spec-doc layout,
	  - embedded spec-byte hashing for pCID provenance,
  - `go:embed` spec context for LLM agents,
  - slot-0 pCID selection of parser/builder roles,
  - why pCID is not an app address, message kind, RPC method, or service-registry
    entry,
  - app/routing semantics carried by pCID-defined payloads or nested payloads,
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

- [x] zugok.1 Scaffold `implementations/poc16-secure-tokens-maps-encrypted-payloads/` from the POC15 baseline.
- [x] zugok.2 Add an analyzer superset gate proving POC15 acceptance categories still appear in POC16 clean runs.
- [x] zugok.3 Add implementation-local `docs/protocols/` specs for every POC16 pCID and remove the stale root mirror per `DI-magug`.
- [x] zugok.4 Add pCID-owned CBOR map payload specimens and at least one normal map-based protocol flow.
- [x] zugok.5 Add pCID-owned compact CBOR array alternatives for constrained-profile agents.
- [x] zugok.6 Replace toy bearer tokens with cryptographically secure bearer tokens.
- [x] zugok.7 Replace toy non-transferable tokens with holder-bound cryptographically secure tokens.
- [x] zugok.8 Prefer a CWT-shaped token format and document any POC-local simplifications.
- [x] zugok.9 Add token validation failures for invalid proof, expiry, revocation, wrong audience, replay, and transfer mismatch.
- [x] zugok.10 Add encrypted payload pCIDs for end-to-end, storage-encrypted, and recipient-specific cases.
- [x] zugok.11 Preserve pCID dispatch when payload slots are encrypted.
- [x] zugok.12 Add encrypted-parent-link cases covering visible envelope parents and hidden payload parents.
- [x] zugok.13 Add decryption failure and unsupported encrypted-pCID containment cases.
- [x] zugok.14 Add minimal key discovery and key rotation for token and payload use.
- [x] zugok.15 Store token, map, array, and encrypted payload artifacts in per-agent filesystem CAS.
- [x] zugok.16 Add CAS/DAG metadata and analyzer checks for ciphertext CID semantics and sparse hidden-parent behavior.
- [x] zugok.17 Add GC and retention behavior for expired tokens, revoked-token history, and encrypted payloads.
- [x] zugok.18 Extend route/storage/compute economics to use secure bearer-token incentives.
- [x] zugok.19 Add exchange-rate and opportunity-cost events affected by token issuer promise history.
- [x] zugok.20 Add raw CBOR diagnostic examples for map, array, CWT-shaped token, holder-bound token, encrypted payload, and failures.
- [x] zugok.21 Update POC16 README and protocol notes.
- [x] zugok.22 Update `DEV-GUIDE-RESOURCES.md` after executable POC16 evidence exists.
- [ ] zugok.23 Reconcile POC16 lessons into the POC17 M4/LoRa plan before implementing constrained-device token/encryption behavior.
- [x] zugok.24 Run Go validation, errcheck, and clean POC16 containers after implementation.
- [x] zugok.25 File TE-ritig for pCID cardinality and parser/builder kernel-role analysis.
- [x] zugok.26 Lock DF choices from TE-ritig before POC16 implementation.
- [x] zugok.27 Embed relevant protocol specs with `go:embed` for every LLM agent pCID promise.
- [x] zugok.28 Log and analyze the spec CIDs supplied to each LLM agent prompt context.
- [x] zugok.29 Add malformed-payload and unsupported-pCID local non-commitment tests without global conformance semantics.
- [x] zugok.30 Route slot-0 pCIDs to pCID-specific parser/builder roles rather than parsing arbitrary app payloads in the transport listener.
- [x] zugok.31 Add analyzer checks for pCID-as-address, universal-`to`, RPC-method, and service-registry regressions.
- [x] zugok.32 Add ACK/error/backpressure/resource-allocation tests across transport, parser, builder, CAS, key, hardware, and app roles.
- [x] zugok.33 Preserve all POC15 named agents in the POC16 config or document a scoped non-superset DI exception.
- [x] zugok.34 Preserve all POC15 executable functionality in the POC16 clean run or document a scoped non-superset DI exception.
- [x] zugok.35 Add analyzer gates for POC15 named-agent, agent-role, and executable-function preservation.
- [x] zugok.36 Document the POC15 agent/function baseline in the POC16 README before adding new POC16-specific behavior.
- [x] zugok.37 Replace flexible pair-list payload bodies with nested CBOR map bodies and update the affected POC16/POC17 documentation.
