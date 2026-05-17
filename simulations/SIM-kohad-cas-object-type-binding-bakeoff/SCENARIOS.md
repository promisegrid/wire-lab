# CAS Object Type Binding Scenarios

These scenarios are evidence for `DR-tumus` DF-tumus.2. They are not a decision
and not a codec registry. Source: `DI-bukoh`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Raw chunk versus pointer bytes | Alice receives bytes whose hash is known, but the local filename is missing. | Whether CID codec identity alone tells Bob whether to parse the bytes as a pointer object or treat them as raw payload. | Object type must survive transport without relying on local paths. |
| Broad codec plus internal kind | Alice stores a CBOR object under a broad PromiseGrid object codec, and the object bytes include `kind = pointer`. | Whether internal kind improves forward compatibility or merely duplicates a CID-level type claim. | TE-43 must avoid two independent type authorities unless the split has a clear rule. |
| Filename suffix negative control | Carol renames a local file from `.ptr` to `.raw` without changing bytes or CID. | Whether path suffixes can safely carry type meaning in sparse replication, export/import, and archival storage. | If suffix changes alter interpretation without changing content identity, suffixes are unsuitable as the primary discriminator. |
| Unknown typed object | Dave receives a CID whose codec he does not implement. | Whether the peer can store, advertise, and forward the object opaquely while avoiding unsafe parsing. | Type binding must define unknown-type behavior for long-lived mixed-version networks. |
| Application object family | Ellen proposes a future application-level CAS object distinct from raw chunks, Merkle nodes, and pointer objects. | Whether the chosen binding model leaves room for new object families without reinterpreting old bytes. | The first type-binding rule should be extensible without changing old CIDs. |

## Expected Outputs

- Evidence for whether `DR-tumus` DF-tumus.2 should prefer codec-only,
  codec-plus-kind, or another type-binding rule.
- A negative-control record explaining why filename and path suffixes should
  not be treated as CAS object identity.
