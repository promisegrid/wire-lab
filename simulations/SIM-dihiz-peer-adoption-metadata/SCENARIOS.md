# Peer Adoption Metadata Scenarios

These scenarios are evidence for `TODO-nivus` 011.10. They are not a decision
and not a frozen adoption-metadata spec. Source: `DI-pukap`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Direct adoption claim | Bob tells Alice he follows pCID X with Q7=yes and Q9=variant-B. | Whether the claim is parseable, signed, scoped, and content-addressable. | The wire shape must be concrete enough for Alice to make protocol-behavior assumptions. |
| Spec answer vocabulary drift | A later frozen spec supersedes pCID X and renames or removes Q9. | Whether answer keys are spec-local, profile-local, globally named, or mapped through migration records. | Adoption metadata must survive spec evolution without making old claims ambiguous. |
| Relationship-scoped adoption | Bob follows one profile when talking to Alice and another when talking to Carol. | Whether adoption promises bind to relationship context instead of only peer identity. | Peer-local variation may be legitimate and should not look like equivocation unless it violates a promise. |
| Revocation and supersedence | Bob stops following pCID X and adopts pCID Y. | Whether the old claim remains auditable while the current claim becomes discoverable. | The design needs explicit current-pointer, supersedence, or freshness semantics. |
| Offline verification | Alice receives Bob's adoption record through Dave while Bob is offline. | Whether third-party carriage works without a central registry. | Content-addressed adoption objects get stronger if offline verification is a first-class scenario. |

## Expected Outputs

- Evidence for the later TE/DR that locks peer-level adoption metadata.
- A list of fields that adoption metadata likely needs: peer/principal, spec
  pCID, answer set, scope, time/freshness, signature/proof, and supersedence.
