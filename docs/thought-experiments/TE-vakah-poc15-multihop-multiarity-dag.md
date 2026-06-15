# TE-vakah: POC15 multihop, multiarity, and raw message DAG

TE ID: TE-vakah
## Status
decided, refined

## Decision under test

Should POC15 expand from a multihop routing plan into
`poc15-multihop-multiarity-dag`, testing promise-based multihop forwarding,
pCID-owned slot-vector variety, raw-message CAS retention, wire-visible parent
links, COSE payload/proof specimens, and transport-session proof alternatives?

## Assumptions

- POC15 is a planned successor to `poc14-wasm`, not executable yet.
- POC15 should be a POC14 superset unless a later DI explicitly authorizes a
  scoped exception.
- PromiseGrid messages use a CBOR `grid(...)` array with slot 0 as `42(pCID)`.
- The protocol named by pCID defines slots 1..N, payload shape, signable view,
  proof encoding, parent-link meaning, validation, and failure behavior.
- Promise Theory remains primary: agents make voluntary promises only for their
  own behavior; route forwarding is not command/control, authorization, policy
  enforcement, or global route truth.
- Run-scoped retention is acceptable for POCs; cross-run persistence is not a
  default because clean runs should remain clean.

## Alternatives

- **Alt A — Multihop only.** Keep POC15 focused on route promises and useful
  routed WASM/stdio work; leave raw CAS/DAG and arity variants for a later POC.
- **Alt B — Archive-only DAG.** Store all raw artifacts in run-scoped CAS and
  build a local review DAG, but keep wire messages shaped like POC14.
- **Alt C — Wire parent DAG only.** Add parent links to valid messages, but do
  not retain every raw artifact.
- **Alt D — Multihop multiarity DAG.** Combine multihop routing, run-scoped raw
  CAS/DAG, pCID-owned arity specimens, parent links in envelope and payload,
  COSE payload/proof specimens, and transport-session proof pressure.

## Scenario analysis

### S1 — Normal multihop route setup

Alice wants to send a message to Dave through Bob and Carol. Alice sends a route
setup promise to Bob promising reciprocal value if Bob can assemble a route under
Alice's constraints. Bob promises only Bob's next hop to Carol. Carol promises
only Carol's next hop to Dave. Dave returns a reachability promise.

- **Alt A** proves the route-promise idea, but the message lineage is visible only
  through local events and ad hoc logs.
- **Alt B** gives excellent review after the run, but the valid messages do not
  carry causal parent links that peers can interpret during the protocol.
- **Alt C** gives valid messages parent links, but rejected or malformed route
  setup bytes can disappear unless separately archived.
- **Alt D** gives both live protocol parent links and complete run review. It
  costs more parser/analyzer work, but the cost matches the design questions POC15
  is supposed to answer.

### S2 — Incentivized forwarding

Bob has limited capacity and asks why Bob should forward Alice's traffic. Alice
can offer reciprocal forwarding, a bearer token Bob values, a non-transferable
future forwarding token, or stake/collateral.

- **Alt A** can model incentives in route payloads, but has no broader way to
  connect payment messages, route setup, and carried messages by exact parent CID.
- **Alt B** can reconstruct the economics afterward, but peers do not see
  protocol-level parent links while evaluating new promises.
- **Alt C** can link payment and route messages, but misses raw monitor/decision
  artifacts that explain why agents chose their offers.
- **Alt D** preserves both the route economics and the raw decision trail.

### S3 — Route durability and asymmetric replies

Alice wants to reuse a route for five messages. Dave replies along a different
path because the reverse route is cheaper or more trusted.

- **Alt A** can implement cached route promises, but route reuse and asymmetric
  causality are harder to review without message DAG links.
- **Alt B** supports post-run review, but the protocol does not exercise how a
  peer handles parent links inside messages.
- **Alt C** supports protocol-visible causal links, but loses invalid or rejected
  bytes that shaped the route state.
- **Alt D** tests one-shot routes, bounded durable routes, and asymmetric reply
  DAGs while retaining every raw artifact needed for later debugging.

### S4 — Malformed and adversarial traffic

Mallory sends malformed CBOR, an unknown pCID, a bad COSE algorithm claim, or an
LLM decision containing command/control wording.

- **Alt A** may reject the messages, but exact raw text/bytes can be lost as in
  current POC14 `decision_rejected` records.
- **Alt B** is strongest for preserving malformed and non-message artifacts.
- **Alt C** is weak because malformed bytes cannot carry valid parent links.
- **Alt D** retains malformed artifacts in the run CAS DAG while keeping valid
  message DAG semantics clean.

### S5 — pCID-owned arity and COSE pressure

Alice, Bob, Carol, and Dave exchange messages under multiple pCIDs: one with no
proof slot, one with a native proof slot, one with envelope parents, one with
payload parents, one with COSE as payload, and one with COSE as detached proof.

- **Alt A** does not test the current design claim that pCID owns arity and slot
  meaning.
- **Alt B** archives the bytes but does not prove parsers can dispatch by
  pCID-owned slot vectors.
- **Alt C** tests parent links but not enough arity/proof variety.
- **Alt D** directly tests the contested design area and should reveal whether the
  envelope rule is too vague, too flexible, or appropriately pCID-owned.

## Conclusions

Alt D is the best POC15 target. POC15 should be renamed to
`poc15-multihop-multiarity-dag` and should explicitly test multihop forwarding,
route economics, raw run CAS/DAG review, wire-visible parent links, pCID-owned
slot-vector variety, COSE payload/proof specimens, and transport-session proof
pressure.

The key separation is:

- Run CAS/DAG preserves every raw artifact for local review.
- Wire parent links are valid-message protocol semantics.
- Transport/session proofs are hop-local.
- Envelope/message proofs are durable object-level authorship or promise proofs
  only when the pCID defines them.

## Implications for open TODOs and pending DIs

- `DI-pamob` remains valid but is narrowed by the POC15 multiarity/DAG expansion.
- `DI-podut` records the rename and expanded POC15 target.
- The POC15 docs should add a message-shape plan alongside route and kernel-role
  docs.
- The first executable implementation should not claim production readiness; it
  should provide analyzer gates for route setup, forwarding, incentives,
  durability, asymmetric routes, raw artifact retention, parent DAG
  reconstruction, and COSE validation.

## Decision status

Decided by `DI-podut`.

## Refinements

### 2026-06-15 — POC15 convergence slice for normal-traffic route multiarity

`DI-kohuj` implements the next executable comparison implied by this TE without
changing the TE's original conclusion. The POC now keeps the universal
`grid([42(pCID), ...])` parse surface while letting `route_v1` use
`grid([42(pCID), parents, payload, proof])` in selected normal app traffic;
records payload-owned parent links through `route_v1` payload fields; traverses
the retained raw-message DAG; models route lifetime, asymmetric response-path
terms, and reciprocal route credits; carries exact `cid_compute_v1` envelopes to
Peggy and Victor over `route_v1`; and records that transport/session signatures
remain hop-local comparison context rather than replacing durable object-level
envelope proofs.
