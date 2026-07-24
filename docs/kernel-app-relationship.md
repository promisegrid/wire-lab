# PromiseGrid Kernel/App Relationship

This document expands the kernel/app relationship notes that were split out of
`GOOD-PRACTICES.md`. It is practical developer guidance, not a frozen
PromiseGrid protocol specification. Source: `DI-vofih`; `DI-rusup`;
`DI-gumum`; `DI-funaf`.

## Core Model

The PromiseGrid kernel is best understood as a set of roles served by agents,
not as one mandatory process. A rich host might run those roles as a daemon,
microkernel-like collection of processes, sidecars, or libraries. A constrained
device might implement only a small subset in firmware. A browser or WASM host
might expose the same promises through runtime adapters.

Apps are agents too. An app can make and receive PromiseGrid messages, and it
can also assume a kernel role when it promises to manage a local resource,
transport, parser, storage surface, hardware device, or runtime capability. The
important practice is to make that role explicit so a local resource promise is
not mistaken for a global command authority.

Kernel-role agents should describe what they promise, what they do not promise,
which pCIDs they handle, which local resources they own, and what local events
they record. This keeps ports, apps, and tests from treating one implementation
shape as the universal PromiseGrid kernel.

## Routing Promises

One common kernel role is message routing. A routing agent routes exact
`grid()` message bytes by pCID to other local agents that have promised to
accept those pCIDs. The routing agent should not parse every application
payload or judge whether an application-level promise was kept.

Startup registration is itself promise-shaped. An app can promise, "I accept
messages with these pCIDs under these local conditions." A routing agent can
promise, "I will deliver matching message bytes to you under these local
conditions." Those are reciprocal local promises, not service-registry
authority and not permission from a central kernel.

If the routing agent cannot deliver, does not recognize a pCID, or has no local
agent promising to accept that pCID, it should record a local non-commitment or
failure event rather than pretending the work was kept. Unsupported pCIDs and
wrong local roles must not silently become successful acknowledgements.

## Parser And Builder Roles

For simple protocols, the route from network to app can be direct:

```text
network -> transport listener -> pCID router -> app
app -> pCID builder -> transport sender -> network
```

For more complex protocols, a parser or builder agent can sit between the
router and the app:

```text
network -> transport listener -> pCID router -> pCID parser -> app
app -> pCID builder -> transport sender -> network
```

The parser role knows how to interpret the payload for one or more pCIDs and
can forward the resulting protocol-owned information to app agents. This is an
implementation pattern, not a requirement that all apps must communicate only
through one router or parser.

Direct app-to-app messaging is also expected in some scenarios. A
routing agent is just another agent making local delivery promises;
there may be more than one routing agent on a node, and agents with
matching promises may exchange messages directly when the local
deployment supports that.

## pCID And Payload Meaning

A pCID is the address, in content-hash space, of a protocol specification
document. It identifies the protocol spec used to interpret the envelope slots
after slot 0. 

It is normal for one protocol spec to describe multiple message types
or operations under the same pCID. In that case, a message type or
operation name usually belongs in the pCID-defined payload. (Adding a
separate envelope-level message-type slot is a valid protocol-design
choice that has not been thoroughly tested as of July 2026.)

The pCID defines how to interpret the payload, including whether the payload is
a CBOR array, CBOR map, COSE object, encrypted byte string, nested `grid()`
message, DAG-CBOR object, or some other protocol-owned shape.

## pCID Granularity Tradeoffs

Do not choose pCID count directly by rule of thumb. First decide which protocol
spec documents are coherent, then let their CIDs become the pCIDs.

Too many pCIDs fragment related concepts across multiple specs. That increases
developer cognitive load, causes duplication, and makes each app promise a
longer list of accepted pCIDs at startup.

Too few pCIDs overload a spec with unrelated behavior. Because routing is based
on pCID, an agent may receive many messages under an overloaded pCID and then
have to parse payload-level message types only to discard most of them.

The practical balance is cognitive load versus CPU and filtering load. Fewer
larger specs can be easier to understand at the routing layer but more
expensive for receiving apps. More smaller specs can reduce irrelevant
delivery but make the protocol ecosystem harder to understand and maintain.

## Evolution Pattern

Spec docs are expected to grow organically. Developers cannot predict
every future message type in a problem domain, so new specs and pCIDs
should be added as new needs become clear rather than mutating old
specs after their pCIDs are in use.

Existing agents do not need to be edited every time a new pCID appears. A node
can add new agents, parser roles, or builder roles that promise to handle the
new pCID. 

The goal is organized, self-documenting growth. Each new protocol spec states
the promises it supports, the message shapes it defines, and the implementation
roles expected to handle it. Each implementation then promises which pCIDs it
handles and under what local conditions.
