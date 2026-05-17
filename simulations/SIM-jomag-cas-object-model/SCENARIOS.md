# CAS Object Model Scenarios

These scenarios make the turn-177 L6 CAS object-model pressure concrete. They
are simulation inputs for `TODO-kituj` / TE-43, not decisions and not a frozen
CAS spec. Source: `DI-pator`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Deterministic CBOR agreement | Alice and Bob encode the same pointer object with independent implementations. | Whether map ordering, integer/string choices, tags, and byte-string boundaries produce identical bytes and therefore identical CIDs. | TE-43 must lock a precise CBOR profile rather than saying "use CBOR" generically. |
| DAG-CBOR interop | Alice stores a Merkle node or pointer object using a DAG-CBOR-compatible representation. | Whether CID links, byte strings, and tags stay compatible with IPFS / IPLD-style tooling without requiring those stacks. | TE-43 must decide whether DAG-CBOR is the default object format or only one allowed profile. |
| CIDv1 object typing | The same content hash could be interpreted as a raw chunk, Merkle node, or pointer object unless type is bound into identity. | Whether CIDv1 codec / multicodec values carry object type cleanly enough to avoid filename suffixes. | TE-43 must lock object typing through CID codecs or explicitly choose another type-binding mechanism. |
| Pointer object identity | Alice creates pointer object CID Y that points at root CID X; Bob has Y but not X. | Whether pointer objects are verifiable CAS objects in their own right and how sparse-CAS resolution behaves when the pointee is absent. | TODO-pipus migration must preserve the distinction between pointer-object identity and pointed-at root identity. |
| Chunker parameter mismatch | Alice chunks with turn-177 FastCDC-style small averages; Bob chunks with promisebase / pitbase Rabin defaults. | Whether the same logical file produces different leaf CIDs and Merkle roots under different parameters. | TE-43 must lock chunking algorithm and full parameter set or make parameterized chunking explicit in object identity. |
| Promisebase adapter | Promisebase / pitbase stores blocks and trees with class headers and stream symlinks. | Which prior-art pieces map to CIDv1 codecs, CBOR pointer objects, and PromiseGrid sparse-CAS behavior. | TE-43 must adopt, wrap, or reject each prior-art mechanism deliberately. |
| Small-object degenerate case | A small group message fits in one chunk. | Whether it still uses the same pointer / root / object-typing rules as a large object. | The model should avoid special cases that create a second identity path for small messages. |

## Expected Outputs

- A TE-43 decision list for CBOR profile, CIDv1 codecs, pointer-object shape,
  chunker parameters, and promisebase / pitbase adoption stance.
- A migration constraint for `TODO-pipus`: historical `.txt` message bytes stay
  historical, while any successor specimen uses explicit pointer and CAS-root
  identities.
- A guide-resource warning that these are provisional object-model scenarios,
  not final PromiseGrid APIs.
