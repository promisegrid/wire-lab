# SIM-hozif: App pCID promises and local observations

This simulation is a promise-first successor to `SIM-vinag-child-tag42-port-observation-pairing`, with rejected `SIM-lasuv-child-kernel-port-capability-outcome-ledger` used only as cautionary evidence. Source: `DI-dikat`.

Vinag had useful heritage from `SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig` and `SIM-dalor-grid-envelope-protocol-owned-signature-slot`: keep the ecosystem-compatible `grid([42(pCID), payload, ...])` envelope, keep `pCID` as Protocol CID, and let the pCID-named protocol own later-slot proof rules. It also showed that first-port work needs auditable records for what an app or kernel promised and what a local observer later saw.

The correction is vocabulary and structure: this sim does not define a registry, a capability table, a dispatch authority, or a universal outcome ledger. It tests whether the same evidence can be expressed as one stable pCID-selected payload protocol with two record kinds:

1. `app_pcid_promise_v1`: an app promises a local kernel, peer, or operator which pCID-selected messages it will receive or handle, and names any bounded assumptions or non-promises;
2. `local_observation_record_v1`: a local observer records what it saw about one attempted promise, including exact-byte or CID-rooted evidence when needed.

This single-pCID shape is intentional. A pCID names the protocol specification, not an individual message type. Promise and observation records are closely related messages in one local app/kernel promise-accounting protocol, so the payload `kind` field carries message-level variation. Source: `DI-bozid`.

## Question

Can a PromiseGrid kernel/app boundary stay simple if apps make explicit local pCID receive/handle promises, kernels make best-effort delivery and evidence promises, and all kept/refused/broken/timed-out outcomes are local observations rather than authority facts?

## Candidate specimen

The specimen under test is a small three-part pairing:

- the current envelope direction: CBOR `grid([42(pCID), payload, ...])`;
- one `app-kernel-local-promises` payload protocol for app, kernel, peer, storage, compute, or device agents to promise their own pCID behavior and remember what they observed about specific promise attempts;
- two payload record kinds under that one protocol pCID: `app_pcid_promise_v1` and `local_observation_record_v1`.

## Promise Theory principles under test

- No agent promises on behalf of another agent.
- A kernel does not know or certify what an app can do; it can only remember what the app promised this kernel and what this kernel later observed.
- An app/kernel crossing is a pCID-selected grid message, even when a local API gives a convenient wrapper.
- Trust remains local and relationship-scoped; an observation record may affect Alice's future choices but does not create a global score.
- A refusal is not failure by itself. It is useful evidence when the promiser honestly says what it does not promise.
- Exact bytes or CID-rooted references are evidence, not authority.

## Scenario pressure

This sim should be scored against:

- `kernel-porting-boundary`: Bob's app promises Bob's kernel it will receive one hello pCID; Bob's kernel attempts delivery and records what happened.
- `promise-accounting-records-kept-storage-promise`: Alice sends data only after Bob has made bounded storage promises and Alice has enough local trust in Bob.
- `promise-accounting-records-refused-service`: Carol asks Alice for data while making reciprocal promises; Alice or Bob records an honest refusal distinctly from silence.
- `cas-object-model-dag-cbor-interop`: exact-byte or CID-rooted witnesses remain compatible with CBOR/DAG-CBOR without making an IPLD runtime mandatory for every small device.
- `l6-cas-starting-profile-bakeoff-long-horizon-reprofile`: older observation records remain meaningful if later pCID specs or storage profiles evolve.

## Why this should score better than Vinag and Lasuv

Vinag scored well because it kept the outer envelope disciplined and kept local trust boundaries visible, but its `service-*`, `supported_*`, and local API vocabulary still invited RPC-shaped readings. Lasuv made the problem clearer by turning local evidence into a port ledger with pCID capability entries. This successor keeps the useful evidence pressure while dropping the registry-shaped machinery.

## Boundaries

- This is not a final PromiseGrid kernel API.
- This does not define a service registry, capability registry, permission issuer, conformance record, or global namespace truth.
- This does not require every app to publish every promise globally; promises may be local to the relationship between app and kernel or between peers.
- This does not add universal envelope semantics beyond `grid([42(pCID), payload, ...])` and pCID-owned later-slot roles.
