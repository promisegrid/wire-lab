# Chunking Identity Scenarios

These scenarios are evidence for `DR-tumus` DF-tumus.3. They are not a decision
and do not coin final pCID-adjacent terminology. Source: `DI-bukoh`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Same bytes, different chunkers | Alice chunks a file with a small FastCDC-style target; Bob chunks the same file with promisebase / pitbase-style Rabin defaults. | Whether leaf CIDs and Merkle roots diverge, and where the differing chunking parameters are visible. | TE-43 must either lock chunking parameters or bind them into object identity explicitly. |
| pCID-driven object | Alice stores an object whose governing pCID defines the chunker. Bob has the object root but not the spec text cached. | Whether chunk verification and fetch planning require resolving the pCID first. | pCID-driven chunking couples object interpretation to protocol-spec availability. |
| Separate chunking descriptor | Alice stores a chunked object whose root points at a chunking descriptor, provisionally called a chunking CID or cCID candidate. | Whether peers can verify chunks and compare roots without overloading the object pCID. | A separate descriptor may isolate chunker evolution but adds another identity object. |
| Profile negotiation mismatch | Alice advertises profile `fastcdc-small`; Bob supports only `rabin-large`. | Whether peers can refuse, bridge, or request raw bytes without silently accepting mismatched roots. | Negotiated profiles need failure behavior and may still need persistent identity binding. |
| Raw-only migration | Alice migrates one historical message as a single raw chunk behind a pointer object. | Whether the first CAS migration can proceed without chunked Merkle roots. | Raw-only may unblock migration but does not answer large-object replication. |

## Expected Outputs

- Evidence for whether `DR-tumus` DF-tumus.3 should use pCID-driven chunking, a
  chunking-descriptor / cCID-like mechanism, negotiated profiles, or raw-only
  deferral.
- A precise list of chunking facts TE-43 must lock before large-object Merkle
  roots can be considered stable.
