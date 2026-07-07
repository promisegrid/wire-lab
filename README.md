# PromiseGrid Wire Lab

Wire Lab is the experimental workspace where PromiseGrid protocol, kernel,
runtime, and application ideas are tested before they are written as guide or
production material.

PromiseGrid is a decentralized computing model based on promises between
agents. Agents do not command each other. They make local promises, exchange
messages, keep or break those promises, and update local trust from their own
observations. Wire Lab exists to find the message formats and runtime patterns
that make that practical across machines, organizations, legal entities, and
long periods of time.

## Current Status

The design is still provisional. The strongest current direction is:

- messages use compact CBOR `grid(...)` envelopes;
- slot 0 carries `42(pCID)`, where the protocol CID identifies the protocol spec;
- each pCID-defined spec owns the remaining slots and payload shape;
- pCID is a protocol selector, not a peer address, message type, or operation
  code;
- protocol objects help agents make, recognize, remember, and evaluate promises;
- trust is local and relationship-specific, not global;
- CAS and parent-linked message DAGs are becoming the durable substrate for
  storage, synchronization, collaboration, and version control.

The most recent proof-of-concept work has moved from simple message exchange to
containerized agents communicating over TCP with exact `grid()` CBOR messages,
signed CWT/COSE capability tokens, sparse per-agent CAS stores, CAR payloads for
object transfer, and observer-collected raw artifacts for later review.

## What This Repo Contains

- `simulations/` contains generated and hand-curated protocol experiments.
- `implementations/` contains executable POCs that pressure-test the design.
- `docs/` contains thought experiments, research notes, and design notes.
- `protocols/` contains TODOs, decision records, and protocol-level planning.
- `DEV-GUIDE-RESOURCES.md` is the source map for people writing the PromiseGrid
  Development Guide.

## Lessons From The POCs

The POCs have repeatedly pushed the design away from RPC and command/control
systems. The durable pattern is promise-first:

- an agent promises what it is willing to do;
- another agent decides locally whether to rely on that promise;
- exact message bytes and CIDs make the promise auditable later;
- capability tokens are signed promises, not permissions from a central
  authority;
- sparse CAS stores let each agent retain only the objects it chooses to retain;
- observer and analyzer tools are test machinery, not production monitors or
  trust authorities.

POC18 is currently exploring whether PromiseGrid CAS, parent-linked grid
messages, reference-set promises, and continuous peer sync can become a
Git/GitHub replacement that also works for large files, issue/review flows,
DevOps-style filesystem management, and future LLM-scale collaboration.

## Where To Read Next

Start with `DEV-GUIDE-RESOURCES.md` for the current design state and links to
the relevant DIs, TEs, TODOs, design notes, and POC outputs. Detailed technical
material belongs there and in the linked docs; this README is only the plain
English orientation.

## License

GPL-3.0, matching the rest of PromiseGrid.
