# Promise Accounting Record Scenarios

These scenarios test peer-local promise accounting records for turn 177's
promise-economy requirement. They are inputs for TODO-kulih / TE-nibar and for
later cross-layer decisions; they do not define a central accounting authority,
global reputation database, or final economics mechanism. Source: `DI-pator`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Kept storage promise | Alice promises to store chunk C for a period; Bob later asks and Alice serves it. | What Bob records locally and how that affects future pull / keep / advertise choices. | Specs may need promise-vocabulary sections that name observable promise outcomes. |
| Refused service | Alice refuses to send C because of policy, cost, group context, or missing authorization. | Whether refusal is recorded differently from failure, corruption, or timeout. | Promise accounting must support honest refusal instead of treating every refusal as misbehavior. |
| Corrupt data | Mallory sends bytes that fail CID verification. | What observation is recorded and how future local decisions change. | CAS verification and peer-local records should compose without central enforcement. |
| Cross-layer decision | L7 group policy says Bob values a root; L6 knows missing chunks; L5 sees offers from Alice and Carol. | What information flows between layers when Bob decides which chunks to pull. | The turn-178 "decides" issue should be made explicit without collapsing all accounting into one layer. |
| Sparse retention | Bob cannot keep all chunks and must choose what to retain. | Whether local records can include promises, costs, group value, and peer reliability. | Sparse-CAS makes retention a policy decision, not a background storage detail. |
| Key or identity rotation | Alice rotates keys or moves to a new site identity. | Whether Bob can preserve useful relationship history without assuming a global identity registry. | 100-year design needs repairable continuity and migration-friendly records. |
| Layperson explanation | A non-expert asks why Bob chooses Alice's chunks over Mallory's. | Whether the record can be explained as "sites make promises and keep them or do not." | TODO-kulih must decide whether specs require easy mental-model sections. |

## Expected Outputs

- A list of record concepts that future specs may need to name: promise made,
  promise kept, refusal, timeout, corruption, retention, cost, and context.
- A TODO-kulih input for deciding whether protocol specs require
  promise-vocabulary, 100-year pressure-test, and layperson mental-model
  sections.
- A cross-layer cleanup input: L5/L6/L7 decisions can consult peer-local records
  without making those records a harness-owned service.
