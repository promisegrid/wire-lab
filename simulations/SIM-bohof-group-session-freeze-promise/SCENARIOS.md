# Group-Session Freeze Promise Scenarios

These scenarios are evidence for `TODO-bisur` 012.8. They are not a freeze and
not a signed merge promise. Source: `DI-pukap`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Verified message DAG | Alice and Bob reproduce a four-message group-session DAG and verify every CID. | Whether specimen evidence alone is enough to freeze the group-session semantics. | If evidence is enough, freeze can be a documented reproducibility claim; if not, conformance tooling is required first. |
| Outer/feed split | The group-session semantics are stable, but a lower outer/feed rule remains unsettled. | Whether the freeze promise can split scopes without pretending both layers are solved. | The merge promise must name exactly what is frozen and what remains provisional. |
| Human-readable merge promise | Steve signs a plain-English promise to merge a group-session specimen. | Whether the promise is precise enough to bind artifacts, scope, and unresolved questions. | Human promises need enough structure to avoid future ambiguity. |
| Cryptographic promise tooling absent | The desired signed-promise tooling is not yet available. | Whether a temporary DI/DR-backed promise is acceptable or freeze must wait. | The project needs a fallback policy for early specimens. |
| Post-freeze mutation request | Carol proposes a breaking group-session change after freeze evidence exists. | Whether the change becomes a new specimen, a superseding pCID, or an amendment to the old lineage. | Freeze semantics must preserve independent evolution without mutating history. |

## Expected Outputs

- Evidence for the later decision on `TODO-bisur` 012.8.
- A minimal freeze-promise checklist that separates evidence, scope, signer,
  artifact identifiers, unresolved questions, and supersedence behavior.
