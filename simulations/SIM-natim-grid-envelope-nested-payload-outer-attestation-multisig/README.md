# SIM-natim-grid-envelope-nested-payload-outer-attestation-multisig: Nested payload with outer attestation Multisig probe

This promoted simulation breeds two parent ideas into one bounded specimen:

- from `SIM-janov-grid-envelope-layer-pcid-nested-signed-payload`, keep the minimal outer envelope plus a nested payload whose actual payload bytes are signed with explicit pCID binding;
- from `SIM-lotiv-grid-envelope-cryptid-multisig-signature-proofs`, reuse exact-byte proof handling, conservative unsupported/quarantine behavior, and Cryptid-style Multisig as the outer attestation representation.

## Design move

The main design move is to stop leaving outer-layer conformance and relay authorship implicit. This specimen requires a mandatory outer attestation proof, carried as a Cryptid-style Multisig object, that signs the exact nested payload bytes together with a small role-and-context transcript.

That transcript is intentionally small but directly pressures the sampled scenarios:

- sparse chunk advertisement can distinguish `holder` from mere `relay` and bind freshness or availability-scope references;
- onward-restraint forwarding can bind a `forwarder` promise to a restraint/policy reference;
- routing can distinguish `origin` from `relay` and bind local freshness or policy references without hidden global state.

## Why this should score higher

Both parents were penalized because they were strong envelope/proof specimens but weak at durable attribution of the current actor's promise. This specimen repairs that specific weakness while keeping the design standalone, auditable, and still payload-extensible.

## Status

This is a simulation-local draft specimen for comparison. It does not freeze the PromiseGrid envelope, does not require a central registry, and does not settle final routing, chunk-feed, or conditional-release object schemas.

This specimen was promoted from review-stage child proposal
`SIM-natim-child-nested-payload-outer-attestation-multisig` from
`ga-canary-20260521-003110` under `DI-fihub`. The ignored proposal artifacts
remain local raw evidence; this directory is the canonical non-child simulation
home.
