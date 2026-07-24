# PromiseGrid Good Practices

This document collects practical guidance for PromiseGrid kernel, app, protocol, and POC developers. It is intentionally written as one practice per paragraph so each item can be discussed, copied, promoted into a spec, or turned into an analyzer check later.

Treat `pCID` as the CID of the protocol specification document, never as a payload CID, message CID, destination address, operation name, or message type.

Use `pCID` only to select the parser, handler, or spec for interpreting the remaining envelope slots. Payload-level routing belongs inside the pCID-defined payload semantics.

Prefer the envelope shape `grid([42(pCID), ...protocol-defined-slots])`. The pCID owns the meaning and arity of every following slot.

Render CIDs as binary on the wire and CIDv1 base32 when printable. Do not render protocol CIDs, message CIDs, or object CIDs as raw hex in filenames, logs, configs, or docs unless explicitly diagnostic.

Use full CID bytes as registry keys. Do not key pCID registries by digest-only debug strings, because different CID prefixes, codecs, or lengths could otherwise be misinterpreted.

Protocol specs must be RFC-like: complete message shapes, slot meanings, expected behavior, state machines when useful, error behavior, examples, and enough detail for an independent implementation.

If a protocol spec changes, its CID changes. Avoid `_v1` or `-v1` naming as the primary versioning mechanism when the pCID already identifies the exact protocol text.

Keep pCID counts low unless there is a real protocol-boundary reason. A pCID is closer to an IP or API version boundary than a per-message operation code.

Do not use universal payload schemas. Each pCID defines its own payload shape, including whether payloads are arrays, maps, nested CBOR, COSE objects, encrypted bytes, or something else.

Use CBOR arrays for compact fixed-shape protocols, especially constrained devices. Use CBOR maps only when the pCID explicitly wants self-documenting payloads.

Do not use key/value pair arrays when a CBOR map would do the same job more safely and clearly.

Never allow body fields to overwrite reserved positional compatibility fields such as promiser, promisee, routing fields, or payload protocol metadata.

Avoid `field_*` prefixes in payload names and code. They are noisy and imply a generic schema where the pCID should define protocol-specific meaning.

The kernel should route exact message bytes by pCID to registered parser or handler roles. It should not grow into a universal parser for every app protocol.

A pCID-specific parser or builder can be a kernel role, kernel module, sidecar process, or local adapter. It does not have to be compiled into one monolithic kernel.

Apps should usually be local processes. Network-facing behavior belongs to transport, listener, sender, and local kernel roles, not directly to arbitrary app code.

Unknown pCIDs, unsupported arities, unsupported operations, and wrong local roles should produce explicit local non-commitments or errors, not kept acknowledgements.

Do not silently acknowledge work a node did not actually promise or perform.

All inter-agent communication in POCs should use real promise-shaped grid messages over TCP unless the artifact is explicitly labeled as a fixture or simulation.

Do not claim a POC uses live transport when it only writes shared files, logs, or CAS entries. That distinction must stay explicit.

Keep raw message bytes intact for review. Diagnostic tools should decode those bytes without replacing the original artifact.

Annotated CBOR diagnostics should include the `grid()` tag, array header, tag-42 wrapper, pCID byte layout, and printable base32 CID.

Use real CBOR I/O for stdio agents. Hex is acceptable for diagnostics, not as the claimed production interface.

Use real WASM runtimes when claiming WASM behavior. Do not substitute scripted host functions or fake adapters without labeling them clearly.

All trust is local, peer-relative, and relationship-specific. There is no global trust authority and no absolute trust score.

No agent can promise on behalf of another agent. Designs that imply delegated authority must be reframed as local promises, capability tokens, or local trust decisions.

Avoid RPC and control vocabulary such as `permission`, `authorization`, `policy enforcement`, `conformance`, and `who may` unless carefully reframed as voluntary promises and local resource decisions.

Use positive PromiseGrid vocabulary: voluntary cooperation, local trust, peer relationships, promises, kept and broken histories, resilient communities, and decentralized operation.

Treat observations as promises from a local vantage point, not global facts. `Alice promises that she observed X` is different from `X is globally true`.

Prefer `event` over `evidence` for production-facing logs unless the word is deliberately about design or test evidence.

Avoid stale vocabulary such as `bindings`, generic `claims`, overloaded `boundary`, and `pressure`. Use standards terms such as CWT claims only when the underlying standard uses them.

Capability tokens are promises by the issuer. Redeeming a token is asking the issuer to keep that promise under the token's terms.

Bearer and non-bearer tokens must be cryptographically secure in serious POCs. CWT and COSE should be preferred over custom token formats unless explicitly experimental.

Do not add a separate `proof` field when the token or COSE object already carries the relevant signature, unless the pCID clearly defines why both are needed.

Supervisors and local resource owners may withdraw local resource promises such as CPU, RAM, sockets, files, devices, or subprocess access. That is local resource protection, not command authority over another agent.

Clean shutdown should be promise-based: apps promise shutdown behavior, supervisors invoke shutdown capability tokens, and failures become local events.

Do not assume rollback is safe. Once newer code or data has produced side effects, corrective roll-forward is more realistic than pretending the universe can be restored.

The stable `grid` binary should be treated as minimal stage0 bootstrap code. It should fetch stage1 microkernel modules and app or runtime code from CAS by CID with owner approval.

CAS is the source of truth. SQLite, JSON state files, indexes, and analyzers are rebuildable caches or views, not authoritative state.

State should be reconstructible from chronological CAS event streams where practical.

Each agent may have its own partial CAS. No design should assume one complete global CAS.

Some CAS content should remain strictly local; some may be shared only encrypted; sharing must remain selective and promise/trust-based.

Garbage collection, retention, backpressure, and rate limits should be local promise-based decisions, not global policy commands.

Use persistent TCP connections and multiplex messages where possible. Rebuilding a connection per message is a performance smell.

Use message CIDs and parent links for message identity, threading, and DAG references instead of inventing ad hoc IDs.

Parent links may belong in envelope slots or payloads only when the pCID says so. Do not make universal parent-link semantics accidentally.

For CAS storage, distinguish raw block bytes, wrapped CBOR or grid objects, CAR files, indexes, and diagnostics. Content identity must be clear for each layer.

CAR files are useful for exchange, export, and import of CAS subsets, but they should not automatically replace the wire envelope unless a TE or DI chooses that.

DAG-CBOR gives IPLD traversal benefits but restricts custom CBOR tags. If using custom `grid()` tags, be explicit whether the codec is DAG-CBOR, GRID-CBOR, or ordinary CBOR.

Use well-known libraries for CID, multicodec, COSE, CWT, CAR, and cryptography where practical. Hand-rolled code must be POC-local and clearly marked.

Spec docs used by LLM agents should be embedded or otherwise supplied as authoritative prompt context, preferably with provenance showing which pCID/spec informed the decision.

Analyzer and monitor behavior is development apparatus, not production global authority. Production cannot assume one observer sees the whole system.

Observer-collected logs should be treated as a subset from one vantage point unless the POC explicitly proves full capture.

Every new POC should be a superset of the previous POC's implemented lessons unless a DI explicitly authorizes a non-superset.

Keep the apparatus/specimen split clear. Harness code, analyzers, scoring, and monitors are not automatically PromiseGrid APIs.

Run containerized POCs in containers. Analyzer-only runs are not full POC runs.

Use clean-run scripts when evaluating POCs so stale runtime state does not contaminate behavior.

Avoid shared Docker volumes as a hidden communication channel between agents unless the POC explicitly models local shared storage.

Protocol and app logs should say who talked to whom, over what transport, using which pCID, what promise was made, and what local result was observed.

For production-shaped examples, each agent should have incentives to participate and should be free to refuse.

Do not add workflow-specific top-level action kinds by default. The default semantic action is a promise; operation details belong inside pCID-defined payload semantics.

Storage, compute, routing, relay, discovery, token exchange, shutdown, and repair should all be modeled as promises, not commands.

For shipping, device, or hardware POCs, hardware access should be represented as a local resource promise, often mediated by capability tokens.

Route exclusion cannot be guaranteed globally. Alice can only seek promises from peers and their peers; she cannot command the whole network to avoid Mallory.

Continuous sync should be promise-based and peer-selective. Native PromiseGrid sync should not inherit Git's explicit push/pull assumptions unless bridging to Git.

For VCS work, reference sets can unify tags, directories, branches, and fetch roots. A branch is a use of a reference set, not a fundamentally separate mechanism.

Large files should be first-class. Rabin chunking and in-band large blob storage avoid Git-LFS-style external dependency problems.

POSIX support must account for regular files, directories, symlinks, hard links, character devices, block devices, FIFOs, sockets, ownership, modes, timestamps, xattrs, ACLs, sparse files, and portability loss records.

DevOps use must preserve ordered journals of machine changes. Apply one change, run relevant triggers, then continue, matching infrastructures.org-style convergence.

Porcelain and plumbing should use the same core library. Do not repeat Git's mistake of building automation around a separate implementation path.

Consider whether local UI/backend communication should also use grid messages, because that may enable local and remote UI control through the same protocol model.

Documentation should keep README plain-English and current, while putting detailed protocol, POC, and design-guide material in `DEV-GUIDE-RESOURCES.md` and protocol specs.

Do not let old POC machinery leak into current guidance as if it were production architecture.

When a common mistake is found, capture it in the dev guide, relevant TODO or DI, and future POC acceptance gates so it does not recur.
