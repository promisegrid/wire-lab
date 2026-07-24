# PromiseGrid Good Practices

This document collects practical guidance for PromiseGrid kernel, app, protocol, and POC developers. It is intentionally written as one practice per paragraph so each item can be discussed, copied, promoted into a spec, or turned into an analyzer check later.

`pCID` is the CID of the protocol specification document.  Any operation name or message type should be described in a payload field. Use `pCID` to select the parser, handler, or spec for interpreting the remaining envelope slots. Payload-level routing belongs inside the pCID-defined payload semantics.

Prefer the envelope shape `grid([42(pCID), ...protocol-defined-slots])`. The pCID owns the meaning and arity of every following slot.

Render CIDs as binary on the wire and CIDv1 base32 when printable. Do not render protocol CIDs, message CIDs, or object CIDs as raw hex in filenames, logs, configs, or docs unless explicitly diagnostic.

Use full CID bytes as registry keys. Do not key pCID registries by digest-only debug strings, because different CID prefixes, codecs, or lengths could otherwise be misinterpreted.

Protocol specs must be RFC-like: complete message shapes, slot meanings, expected behavior, state machines when useful, error behavior, examples, and enough detail for an independent implementation.

If a protocol spec changes, its CID changes. Avoid `_v1` or `-v1` naming as the primary versioning mechanism when the pCID already identifies the exact protocol text.

Keep pCID counts low unless there is a real protocol-boundary reason.
A pCID is an indicator for routing a message to the right agent,
similar in spirit to a TCP port being an indicator for routing a
packet to the right application.

Do not use universal payload schemas. Each pCID defines its own payload shape, including whether payloads are arrays, maps, nested CBOR, COSE objects, encrypted bytes, or something else.

Use CBOR arrays for compact fixed-shape protocols, especially
constrained devices. Use CBOR maps only when the pCID explicitly wants
self-documenting payloads. Do not use key/value pair arrays when a
CBOR map would do the same job more safely and clearly.

The kernel should route exact message bytes by pCID to registered parser or handler roles. It should not grow into a universal parser for every app protocol.

A pCID-specific parser or builder can be a kernel role, kernel module, sidecar process, or local adapter. It does not have to be compiled into one monolithic kernel.

The PromiseGrid kernel is best understood as a set of roles served by one or more agents, not as one mandatory process. Message routing is one kernel role; resource ownership, parsing, supervision, storage, and transport may be separate roles.

Apps are agents, usually implemented as native processes or WASM/WASI modules. An app may also assume a kernel role for a local resource or transport.

A local routing agent should receive promises from apps about which pCIDs they accept, and make conditional promises about delivering messages. 

For complex payloads, a routing agent may forward pCID-selected bytes to a parser agent that understands that pCID and then forwards to the correct app agents. 

A single protocol spec may define multiple payload message types under one pCID. Put operation or message-type fields in the payload by default; do not add an envelope-level message type without an explicit pCID/spec decision.

Choose pCID granularity by protocol coherence, not by trying to minimize or maximize count. Too many pCIDs fragment specs and startup receive promises; too few pCIDs make for large spec docs and force agents to parse irrelevant messages, creating higher CPU load. 

Spec docs and pCIDs should evolve organically by adding new frozen specs and handlers rather than mutating old specs. 

When possible, load agent executables or modules from CAS.  Keep
agents small and focused on one role.  Avoid monolithic agents that try to do everything.

See `docs/kernel-app-relationship.md` for the longer explanation of kernel roles, app agents, parser roles, pCID routing promises, and pCID granularity tradeoffs. 

Do not silently acknowledge work a node did not actually promise or perform.
Unknown pCIDs, unsupported arities, unsupported operations, and wrong
local roles should produce explicit local non-commitments or errors,
not kept acknowledgements. 

All inter-agent communication in POCs or example code should use real
promise-shaped grid messages over TCP, websocket, or other network or
IPC transport rather than "cheating" by sharing disk files, unless the
transport to be tested or exemplified is explicitly file-based. 

In POC code, keep raw message bytes intact for review. Diagnostic
tools should decode those bytes without replacing the original
artifact.

Annotated CBOR diagnostics should include the `grid()` tag, annotated
hex array header, tag-42 wrapper, pCID byte layout, and printable
base32 CID.

All messages between agents should be CBOR-encoded and wrapped in a
`grid()` envelope with a pCID.  Do not use JSON, XML, or other formats
for inter-agent messages -- they are not PromiseGrid-native and will
require extra translation layers.

All trust is local, peer-relative, and relationship-specific. There is
no global trust authority and no absolute trust score.

No agent can promise on behalf of another agent. Designs that imply
delegated authority must be reframed as local promises, capability
tokens, or local trust decisions.

Use positive PromiseGrid vocabulary: voluntary cooperation, local
trust, peer relationships, promises, kept and broken histories,
resilient communities, and decentralized operation. Avoid RPC and
control vocabulary such as `permission`, `authorization`, `policy
enforcement`, `conformance`, and `evidence` unless carefully reframed
as voluntary promises and local resource decisions. Prefer `event`
over `evidence` for production-facing logs unless the word is
deliberately about design or test evidence. Avoid stale vocabulary
such as `bindings`, generic `claims`, overloaded `boundary`, and
`pressure`. Use standards terms such as CWT claims only when the
underlying standard uses them.

Treat observations as promises from a local vantage point, not global
facts. `Alice promises that she observed X` is different from `X is
globally true`.

Capability tokens are promises by the issuer. Redeeming a token is
asking the issuer to keep that promise under the token's terms.

Bearer tokens are bearer promises and might be redeemable for more
specific tokens by anyone who holds them.  Bearer tokens might be
traded for other bearer tokens issued by other agents.  The relative
exchange rate of bearer tokens might be indicative of relative trust
between agents, but it is not a global trust score, and the perception
of relative value is local to the perceiving agent.

Bearer and non-bearer tokens must be cryptographically secure in
serious POCs or applications. CWT and COSE should be preferred over
custom token formats unless experience shows a real need for a new
format.

Do not add a separate `proof` field when the token or COSE object
already carries the relevant signature, unless the pCID clearly
defines why both are needed.

Kernel-role agents represent local resource owners and may withdraw
local resource promises such as message routing, CPU, RAM, sockets,
files, devices, or subprocess access. That is local resource
protection as described in the kernel's promises to hosted agents, not
command authority over another agent.

In POC and demo scenarios, clean shutdown should be promise-based:
apps promise shutdown behavior, supervisors invoke shutdown capability
tokens, and failures become local events.  Some of these same
principles apply to production code, where a kernel-role agent may
request shutdown within a local resource promise.

Do not assume rollback is safe. Once newer code or data has produced
side effects, corrective roll-forward is more realistic than
pretending the universe can be restored.  For example, an upgrade
agent that produces undesirable results cannot be assumed to be
safely reverted by simply rolling back to an older version.
Particularly in the case of side effects such as network messages,
file writes, or database updates, rollback is not a safe assumption.

The stable `grid` binary should be treated as minimal stage0 bootstrap
code. It should fetch stage1 microkernel modules and app or runtime
code from CAS by CID with direct or indirect approval from the local
resource owner.

CAS is the source of truth. SQLite, JSON state files, indexes, and
analyzers are rebuildable caches or views, not authoritative state.

State should be reconstructible from chronological CAS event streams
where practical, and when side effects can be accounted for.  Event
streams should be append-only, with each message referencing its
parents by CID.  

Each agent may have its own partial CAS. No design should assume one
complete global CAS.  It's likely that the best assumption is that
each agent's CAS contains a full or partial DAG of the agent's own
messages, plus a partial DAG of other agents' messages that the local
agent has chosen to store. An agent's DAG might best be thought of as
a local "view" of the global DAG, not a complete copy of it.  

An agent's DAG might best be organized as one or more parallel
timelines, roughly similar to branches in a git repository.  Unlike
git, an event or message on one timeline might also be referenced from
another timeline, without causing a complete merge.

Tokens, both bearer and non-bearer, might in certain case be
double-spendable, for example of different timeline in the same DAG.
This allows for speculative execution, for example in a grid-based
simulation or forecast where multiple timelines are explored in
parallel.  Double-spend tokens should be treated as local promises,
not global facts.  Merging timelines should account for the presnce of
double-spend, either by rejecting one timeline, or by appending
reconciliation messages to the merged timeline.  This process is
roughly analogous to git's merge conflict resolution, or
reconciliation of transactions in an accounting ledger.

Some CAS content should remain strictly local; some may be shared only
encrypted; sharing must remain selective and promise/trust-based.
Local-only content should be clearly marked as such, and might exist
in a separate local-only CAS store.  

Garbage collection, retention, backpressure, and rate limits should be
local promise-based decisions, not global policy commands.

Use persistent TCP connections and multiplex messages where possible.
Rebuilding a connection per message is a performance smell.

Use message CIDs and parent links for message identity, threading, and
DAG references instead of inventing ad hoc IDs.

Parent links may belong in envelope slots or payloads as described in
the protocol spec that the pCID references.

CAR files are useful for exchange, export, and import of CAS subsets,
and can be transported or stored in a `grid()` envelope.

DAG-CBOR might be used as the structure of a payload slot if the pCID
document specifies that.

A protocol spec document might be markdown, JSON, IPLD DAG-CBOR, or
other formats.  Self-documenting formats should be preferred.
Consideration should be given to avoiding proliferation of format
parsers unless the application requires an unusual format.  The
multicodec in the pCID identifies the content codec of the spec
document bytes, not the hash encoding. For spec docs in markdown
format, the `raw` multicodec (0x55) should be used.

Use well-known libraries for CID, multicodec, COSE, CWT, CAR, and
cryptography where practical. 

Spec docs should be pre-hashed and referenced by pCID.  Adding a
symlink to the spec doc is preferred, where the symlink name contains
the CIDv1 base32 rendering of the pCID.  The base32 rendering of the
pCID should be hardcoded in protocol code, rather than re-generating
the pCID at startup.  Yes, this is contrary to the usual "don't
hardcode" advice in software engineering, but prevents the possibility
of a spec doc being changed, either accidentally or maliciously, which
would then change the pCID and break protocol compatibility.  

An agent might validate local copies of a spec doc by re-hashing it
and comparing the result to the hardcoded pCID.  

Spec docs used by LLM-based agents should be embedded or otherwise
supplied as prompt context.

In POCs and tests, analyzer and monitor behavior is development
apparatus, not production global authority. Production cannot assume
one observer sees the whole system. Observer-collected logs should be
treated as a subset from one vantage point unless the POC explicitly
proves full capture. Keep the apparatus/specimen split clear. Harness
code, analyzers, scoring, and monitors are not automatically
PromiseGrid APIs.

POCs, examples, and demos should host all runtime code in containers
for portability and isolation.  Containers and container networks
should simulate realistic production-shaped deployments.  Runtime data
and logs MUST be written to e.g. /tmp in the host filesystem or to a
mounted container volume.  POC and demo runtime data and logs MUST NOT
be written into the source code repository, even when there is a
.gitignore rule.  

Write clearly named scripts or Makefiles for POC or demo setup,
execution teardown, and reset.  All of these functions might be in the
same script or Makefile.  The default behavior should be to reset the
POC or demo to a clean state before execution, while leaving the
runtime state and logs in-place after execution for analysis and
debug.

Avoid shared Docker volumes as a hidden communication channel between
agents unless the POC or example explicitly models local shared storage.

POC analysis tools and logs should be able to describe who talked to
whom, over what transport, using which pCID, what promise was made,
and what local result was observed.

For POCs, demos, and examples, each agent should have incentives to
participate and should be free to refuse.  For production code,
similar incentives MUST be explicitly modeled in the protocol specs
and agent algorithms.  Failure to provide incentives, reciprocal
promises, or other trust mechanisms will lead to either a
dysfunctional design or a design that collapses to conventional RPC.

Storage, compute, routing, relay, discovery, token exchange, shutdown,
and repair should all be modeled as promises, not commands.

For an agent representing or managing access to devices or hardware,
hardware access should be represented as a local resource promise,
often mediated by capability tokens.  This agent is serving in a
kernel role and might best be thought of as being similar to a server
process managing access to hardware in a microkernel OS.  

Incorporate local security policy into local protocol design.  Improve
on traditional network-layer protocols that are generic across
organizations and agnostic to local policy and regulatory compliance.
For example, a financial organization might require that members not
import or export data on a thumb drive, with consequences as described
in regulations or employee agreements. Protocol specs should reflect
these real-world expectations and codify protocol-level measures that
support them.  This is an example of a local security policy that
calls for "route exclusion", where the route is via a thumb drive.  In
any network, route exclusion cannot be guaranteed globally, but in
PromiseGrid, the protocol spec can and should account for failure.
Alice can seek promises from peers and indirectly from their peers and
so on.  But she cannot command the whole network to avoid Mallory, so
in those cases that require geofencing, organizational boundaries, or
other constraints, the protocol spec should describe how to request
and verify such promises. These descriptions should mirror or refer to
the real-world expectations of outcomes, e.g. ITAR, GDPR, or other
regulatory compliance; doing so ensures that protocol designers are
aware of and focused on the security and compliance implications of
their design choices.  Risk mitigation should include encryption of
payloads, secret storage of keys, and other measures that reduce the
risk of accidental or malicious data exfiltration.  

Encrypt payloads by default, using strong algorithms.  Use secret
servers or other secure storage for keys and passphrases.  Do not rely
on the network to provide security against exfiltration.  

Think of the CAS as being a VCS-like DAG store of content-addressed
objects, with references to those objects.  Objects are best organized
in timelines of events, with parent links to previous events.  Doing
so allows the CAS to itself serve as a git-like VCS. Reference sets
serve as tags, directories, branches, and merkle roots. A timeline
(branch) is a use of a reference set, not a fundamentally separate
mechanism.

For grid-based VCS design, large files MUST be first-class. Rabin chunking and
in-band large blob storage avoid Git-LFS-style external dependency
problems. POSIX support must account for regular files, directories,
symlinks, hard links, character devices, block devices, FIFOs,
sockets, ownership, modes, timestamps, xattrs, ACLs, sparse files, and
portability loss records.

For grid-based DevOps design, tools must be able to preserve ordered
journals of machine changes.  Rather than checkout or update the
latest version of all files in a tree or root filesystem, instead
apply one change set, run relevant triggers, then continue, matching
infrastructures.org-style congruence.

Porcelain and plumbing should use the same core library. Do not repeat
Git's mistake of building automation around a separate implementation
path.

GUI, TUI, CLI, and web UIs should be built on top of the same core
library as each other.

UI communication to backend or business-logic agents should also use
grid messages, because that enables local and remote UI control
through the same protocol model.
