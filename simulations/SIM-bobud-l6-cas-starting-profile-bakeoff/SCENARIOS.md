# L6 CAS Starting Profile Scenarios

These scenarios are evidence for `DR-tumus` DF-tumus.1. They are not a decision
and not a frozen L6 CAS spec. Source: `DI-bukoh`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Bridge-first exchange | Alice stores a pointer object and Bob tries to inspect it with IPFS / IPLD-shaped tooling. | Whether DAG-CBOR-compatible bytes, CIDv1 values, and multicodec conventions make useful interop possible without importing a full external stack. | If bridgeability removes ambiguity cheaply, the IPFS / IPLD-aligned profile becomes a stronger starting point. |
| Prior-art substrate reuse | Alice maps promisebase / pitbase block and tree structures into PromiseGrid object identities. | Whether class headers, stream conventions, and chunker assumptions can be wrapped cleanly without leaking prototype-specific semantics. | If the adapter is natural, promisebase may be a substrate; if not, it remains prior-art evidence only. |
| Minimal migration seed | Alice migrates one historical inline group-session message into a pointer object plus raw CAS bytes. | Whether raw chunks plus a minimal pointer object are enough to test sparse fetch, verification, and preservation of historical bytes. | If the minimal profile proves the migration path, TE-43 can defer chunked Merkle complexity. |
| Mixed-version peer | Alice uses the minimal profile while Bob experiments with DAG-CBOR Merkle nodes. | Whether peers can reject, store opaquely, or bridge objects whose profile they do not yet implement. | The first profile needs an extension story even if it starts small. |
| Long-horizon reprofile | A future implementation wants to replace the first profile with a richer object graph. | Whether old pointer objects and raw chunks remain addressable and explainable after a later profile lands. | The starting profile should avoid identity choices that become migration debt. |

## Expected Outputs

- Evidence for whether `DR-tumus` DF-tumus.1 should choose an IPFS/IPLD-aligned,
  promisebase-adapter, or minimal pointer/raw starting profile.
- A list of assumptions that TE-43 must lock before guide prose can describe a
  concrete L6 CAS profile.
