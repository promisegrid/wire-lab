# TE-ritig: pCID cardinality and parser/builder kernel roles

TE ID: TE-ritig

## Status

decided, refined

## Decision under test

POC16 needs to decide how many production-shaped pCIDs PromiseGrid should
normally use, and what a kernel does with a pCID after it reads slot 0 of
`grid([42(pCID), ...protocol-defined-slots])`.

The immediate correction is that `pCID` is not a destination address, an app
name, a message kind, or a service registry entry. A pCID names the protocol
specification that defines the rest of the envelope slots and the payload
contract. It is closer to an IP-version/API-version selector than to an IP
address, TCP port, RPC method, or routing table entry. Payloads or nested
payloads may carry local and network routing information, but that routing
information can vary by protocol and should not force the transport listener to
learn every application addressing scheme.

## Assumptions

- The outer envelope shape remains CBOR `grid([42(pCID),
  ...protocol-defined-slots])`.
- The pCID in slot 0 is the CID of a protocol spec document, not the CID of a
  payload or message instance.
- Protocol specs may define CBOR arrays, CBOR maps, COSE objects, encrypted
  bytes, parent links, proof slots, or other slot meanings after slot 0.
- PromiseGrid kernels are not necessarily monolithic. A given runtime may have
  a transport listener, pCID router, payload parser, payload builder, CAS role,
  key role, hardware role, resource allocator, and app adapter as separate
  processes, objects, WASM modules, firmware functions, or linked library
  routines.
- Promise Theory vocabulary remains load-bearing: each component promises only
  its own behavior; no component commands another agent; all trust and outcome
  judgment are local.
- POC16 should avoid overfitting to IP, TCP, OSPF, RPC, HTTP, or operating-system
  kernel analogies. Those analogies are useful only when they clarify a concrete
  role boundary.

## Alternatives

### Alternative A: Many pCIDs

Every message kind, operation, app handler, or local service gets its own pCID.
For example, `need_advertisement_v1`, `offer_v1`, `storage_put_v1`,
`storage_get_v1`, `compute_v1`, and `route_probe_v1` all become separate
protocol selectors.

This makes small handlers easy to write, but it makes pCID behave like an RPC
method table. It also fragments compatibility: a community that changes a single
operation grows another pCID even when the larger protocol family is still the
same.

### Alternative B: Moderate pCIDs by protocol family and major version

A pCID identifies a coherent protocol family and major-version contract. The
payload contains a pCID-defined promise body, local addressing terms, nested
selectors, or operation details. For example, a storage protocol pCID might
cover store, retrieve, retention, GC, and storage-price promises, while a
shipping protocol pCID might cover scale, label, accounting, and fulfillment
promise bodies.

This treats pCID as a stable protocol selector while allowing payloads to carry
the destination, topic, route, key fingerprint, DID, CID-rooted path, or other
addressing term that the protocol family actually needs.

### Alternative C: Very few broad pCIDs

A small number of broad PromiseGrid pCIDs cover most application semantics. The
payload or nested payloads carry nearly all local routing, operation, promise
kind, economic, and relationship information.

This maximizes long-term stability and minimizes pCID churn, but it risks
creating a large universal payload language. That can become as problematic as a
universal claim-card map: convenient for experiments, expensive for small
devices, and vague about which component promised to understand which subset.

### Alternative D: Layered pCIDs

The outer pCID selects a transport/session/profile family, and an inner
protocol-selected field or nested pCID selects application semantics. This
resembles layered stacks, but each layer is still a PromiseGrid promise surface
rather than a command hierarchy.

This can preserve a small outer parser while allowing independent evolution of
inner semantics. The danger is accidental recursion: if every nested selector
gets called a pCID too early, the design can recreate Alternative A under a
different name.

## Parser and builder role alternatives

### Parser model 1: Monolithic kernel parser

The kernel reads slot 0, parses every known application payload, extracts
destination-like fields, and routes messages to apps.

This is easy in an early POC, but it is the wrong center of gravity for
PromiseGrid. It makes the transport kernel grow with every protocol family and
tempts it to become a service registry, app router, conformance authority, or
RPC dispatcher.

### Parser model 2: Parser returns universal routing facts to the kernel

The transport listener reads pCID, invokes a pCID-specific parser, and asks that
parser to return a generic routing record such as `to`, `from`, `kind`, or
`route`.

This still over-assumes that all protocol families have comparable addressing.
One pCID may address by app nickname, another by key fingerprint, another by
DID, another by CID-rooted path, another by route promise, and another by local
hardware capability token. A generic routing record either becomes too narrow or
turns into a universal schema by stealth.

### Parser model 3: Parser/builder as kernel roles

The transport listener reads only enough of the envelope to identify the pCID and
select a pCID-specific parser role. The parser role is a kernel-adjacent module,
driver, or process that understands the pCID and talks directly to apps using
the app-facing shape that pCID defines. The reverse path sends app intent to a
pCID-specific builder role, which produces exact envelope bytes for the
transport sender.

The receive flow is:

```text
network -> transport listener / pCID router -> pCID-specific payload parser -> app
```

The send flow is:

```text
app -> pCID-specific payload builder -> transport sender -> network
```

This model keeps the transport listener small without pretending that the
listener can understand every application addressing scheme. It also matches
non-monolithic kernels: the parser/builder role can be a Go object, subprocess,
WASM module, MCU firmware function, dynamically loaded module, or app-owned
adapter depending on runtime constraints.

### Parser model 4: App-only parsing

The transport listener reads slot 0, stores exact bytes, and delivers every
matching pCID directly to apps that promised to receive that pCID. Apps parse
their own payloads.

This is attractive for tiny kernels, but it can duplicate parser code across
apps and complicate resource control. It also gives the kernel fewer places to
apply local byte-size limits, proof checks, decryption handoff, malformed-byte
containment, and app-facing compatibility adapters.

## Scenario analysis

### Scenario 1: Normal receive with local app addressing

Alice sends Bob a shipping-workflow message. The outer pCID names a shipping
protocol family. Inside the payload, the protocol may address `ups_label_printer`
by local app name on Bob's node.

Under Alternative A, each shipping step may be a separate pCID, so Bob's
transport listener looks like an RPC method dispatcher. Under Alternative B, one
shipping pCID selects the shipping parser role, and the shipping parser decides
which local app receives the payload. Under Alternative C, a broad PromiseGrid
pCID selects a large general parser, which then needs to understand shipping
subset semantics. Under Alternative D, an outer session pCID selects a session
parser, which then hands a nested shipping selector to a shipping parser.

Model 3 is the cleanest fit: Bob's transport listener does not need to know what
`ups_label_printer` means. It only promises to recognize the pCID and hand exact
bytes to the parser role that promised to handle that pCID.

### Scenario 2: Normal send with protocol-specific builders

Carol's app wants Dave to compute a CID-named function. The app should not need
to manually assemble every CBOR slot, proof, COSE object, encryption layer, and
parent-link convention. It can send a local app-facing promise to the
compute-protocol builder role. The builder promises to produce exact
`grid([42(pCID), ...])` bytes according to the pCID spec.

If the builder cannot sign, encrypt, encode, or allocate memory, that is a local
non-commitment or local broken implementation promise. It is not evidence that
Dave broke anything, and it is not a global protocol-conformance judgment.

### Scenario 3: Variable addressing schemes

Ellen's storage protocol payload addresses an object by CID-rooted path.
Frank's identity protocol payload addresses a peer by DID. Grace's encrypted
mail protocol payload has no visible destination after slot 0 because routing is
implicit in a recipient key hint. Heidi's local hardware protocol addresses a
printer-port token.

A generic `to` field returned to the transport kernel does not survive this
scenario. It either loses meaning or becomes a large universal address union.
The parser role should own these semantics. If a protocol family needs local
delivery to an app, that parser should have its own app-facing promise boundary.

### Scenario 4: ACKs and positive responses

ACKs must not be treated as one universal kernel feature. There are at least
four different meanings:

- a transport/session event that exact bytes were written to or read from a peer
  connection;
- a pCID-router event that a local parser role promised to accept bytes for that
  pCID;
- a parser/builder event that bytes were parsed or built according to that
  parser's local promise;
- an application-level promise response, such as Bob promising that Alice's
  requested storage object was retained.

Only the last kind is usually visible protocol meaning between application
agents. The first three are local implementation events unless the relevant pCID
explicitly defines a wire-visible ACK promise. When a pCID does define an ACK or
response, it should reference the exact parent message CID or pCID-defined
request identifier instead of relying on RPC-style ambient call state.

### Scenario 5: Errors and malformed input

If the outer CBOR grid tag is malformed, the transport listener cannot even
select a pCID parser. It can retain exact bytes locally if promised, but it
cannot claim that the sender broke the application protocol.

If slot 0 contains an unsupported pCID, the listener records local
non-commitment. It may store exact bytes in CAS or drop them under local
retention promises.

If the parser role is unavailable, crashes, or is out of memory, that is local
resource evidence about the receiver's implementation. It should not lower trust
in the sender.

If the parser role accepts the pCID but finds malformed pCID-owned payload bytes,
the parser may record a local malformed-input event, may pass a promise-shaped
failure to the app, or may cause a pCID-defined negative response to be built.
Which of those happens is protocol-specific.

If the destination app has not promised to receive the resulting app-facing
message, the parser records local non-commitment. It does not impose delivery.

### Scenario 6: Resource allocation and backpressure

The transport listener promises socket/session/buffer behavior. The pCID router
promises a bounded parser queue. The parser promises CPU and memory use within a
local budget. The builder promises serialization, proof, encryption, and parent
link construction within its own budget. CAS, key, crypto, clock, hardware, and
process-lifecycle roles promise their own resource surfaces.

Backpressure should be promise-based at each boundary. A parser can decline to
promise more queue capacity. A builder can decline expensive encryption. A CAS
role can decline retention without enough token value or local storage budget.
The app can decline a semantic promise because opportunity cost is too high.
None of these decisions are global authorization failures.

### Scenario 7: Encrypted payloads and hidden routing

Alice sends Bob an encrypted payload whose visible envelope reveals only the
outer pCID and perhaps visible parent links. The parser may need a key role to
decrypt enough inner data to decide local delivery, or it may hand opaque bytes
to a key-holder app. Relays may forward exact bytes without knowing the inner
recipient.

This argues against putting universal routing extraction in the transport
listener. Some protocols will intentionally hide routing details from
intermediate agents. POC16 should preserve pCID dispatch when only slot 0 is
clear and should model the local consequences when a parser cannot decrypt or
cannot access the needed key role.

### Scenario 8: Multi-hop forwarding and route promises

Alice discovers a route to Ivan through Bob and Carol. A route protocol may
contain route promises, reusable route terms, asymmetric return paths,
forwarding-price promises, parent links, and failure responses.

The transport listener should not understand route semantics. A route parser
role can interpret route payloads, maintain route-local state, and ask the
transport sender to forward exact bytes to a promised next hop. If a peer stops
forwarding, local trust and token exchange rates may change, but those updates
belong to the agent observing the promise outcome.

### Scenario 9: Constrained POC17-like runtime

On a Cortex-M4/LoRa runtime, there may be no separate OS process called a
kernel. The parser role may be a small compiled function that recognizes one or
two pCIDs and maps them to firmware callbacks. The same conceptual role split
still applies: the firmware does not need a universal `to` schema, and it should
not grow handlers for unrelated protocol families.

This favors Alternative B or a constrained subset of Alternative D over
Alternative C's broad parser.

### Scenario 10: LLM-backed app agents

An LLM-backed agent needs embedded spec docs for the pCIDs it promises to use.
The LLM can reason about promise terms, economics, and local trust, but exact
CBOR parsing, proof verification, decryption, and parent-link extraction should
be handled by deterministic parser/builder roles whenever possible.

This gives the LLM the right semantic context without asking it to be the
transport kernel or byte-level parser.

## Conclusions

- Reject Alternative A as the default. It overuses pCID as an RPC method table
  or service registry and will likely recreate command-dispatch vocabulary.
- Keep Alternative B as the strongest current default candidate: pCID by
  protocol family and major version, with payload/nested payload semantics for
  operation, destination, route, economic terms, and local addressing.
- Keep Alternative C as a pressure-test, not a default. It may reveal useful
  simplifications, but broad universal payload languages are risky for small
  devices and protocol clarity.
- Keep Alternative D where a real layer boundary exists, especially for
  encrypted/session/transport profiles, but do not let nested selectors recreate
  pCID-per-operation fragmentation.
- Prefer parser model 3 for POC16: slot 0 selects a pCID-specific parser or
  builder role; the transport listener/pCID router does not parse arbitrary
  application payloads and does not require a universal routing record.
- POC16 should explicitly test ACK, error, backpressure, resource allocation,
  encryption, CAS, route, and app-delivery behavior across the role boundaries.

## Implications for open TODOs and pending DIs

- `DI-nogij` supersedes `DI-mubul` for the over-strong phrase "kernel dispatches
  normal application messages only by slot-0 pCID." The corrected rule is that
  slot 0 selects the protocol-specific parser/builder role; pCID-defined
  payloads or nested payloads own app/routing semantics.
- `TODO-zugok` should keep `docs/protocols/` and `go:embed` spec-context work,
  but its POC16 implementation tasks should be framed around parser/builder
  kernel roles rather than transport-listener app dispatch.
- POC16 still needs DF before implementation for:
  - exact pCID inventory and cardinality policy,
  - which pCIDs are parser/builder roles versus pure app-handled protocols,
  - which ACK/error events are local implementation events versus wire-visible
    pCID-defined promises,
  - how resource allocation and backpressure promises cross transport, parser,
    builder, CAS, key, hardware, and app roles.
- `DEV-GUIDE-RESOURCES.md` should describe the current state as "TE/DF pending"
  rather than final API guidance.

## Decision status

Decided. `DI-nogij` locked the parser/builder-role interpretation before POC16
implementation: slot 0 selects the pCID-specific parser or builder role, while
pCID-defined payloads or nested payloads carry app, operation, destination,
route, economic, and local-addressing semantics. `DI-vulit` locked POC16
implementation from that default, and `DI-gazin` later corrected the executable
POC16 runtime so a real parser-role process sits between apps and the transport
kernel.

## Refinements

### 2026-06-20 — POC16 executable parser-role correction

`TODO-sosoj` / `DI-gazin` implements this TE's parser model 3 in POC16 rather
than leaving it as profile pressure. The corrected executable flow is:

```text
app -> parser role -> transport kernel -> peer kernel -> parser role -> app
```

The transport kernel now accepts normal app-side frames only through
`kernel_transport_v1`, routes exact envelopes by pCID to parser roles that have
promised receive capability, and does not decode normal application payload
fields such as `to`. The parser role owns pCID-owned payload projection, local
app delivery, ACK/non-commitment behavior, and raw diagnostic artifacts for the
app/parser/kernel byte flow. This refinement does not change the TE's
alternatives or conclusions; it records that the surviving/default path has now
been made executable.
