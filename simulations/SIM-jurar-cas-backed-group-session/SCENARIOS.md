# CAS-Backed Group-Session Scenarios

These scenarios test a turn-177 successor shape where L7 group-session
semantics point at L6 CAS roots and pointer objects. They are inputs for
TODO-pipus migration design, not a rewrite of the existing `.txt`
group-session draft and not a frozen PromiseGrid message format. Source:
`DI-pator`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Additive successor specimen | Existing `.txt` messages remain historical while a new CAS-backed specimen is added beside them. | Whether migration can be additive without rewriting historical message bytes or invalidating existing CIDs. | TODO-pipus must design overlap / successor mechanics rather than mutate old evidence. |
| Group-visible identity | Alice posts a group message whose visible identifier could be a pointer-object CID, message-root CID, or envelope CID. | Which CID humans, tools, parent links, and acknowledgements should cite. | The migration must avoid two competing identities for one logical message. |
| Known-member group | Alice, Bob, and Carol participate with known relationship identities while Mallory attempts anonymous or unrecognized participation. | Whether group-session semantics require member-authenticated provenance, permit pseudonyms, or explicitly represent anonymous participation as a policy exception. | Turn 178's anti-anonymity stance should become group-identity pressure without silently banning every future privacy-preserving design. |
| Parent links through CAS | Bob replies to Alice by referencing a parent that may be a pointer object or Merkle root. | Whether parent semantics survive when the parent is no longer an inline text file. | Group-session successor work must distinguish L7 parent meaning from L6 object resolution. |
| Arbitrary body shape | Alice's message body is a CBOR text string, CBOR map, encrypted blob, signed payload, or large file root. | Whether group semantics can stay stable while body bytes vary. | TODO-pipus and TE-43 must agree on what the group layer sees versus what CAS stores. |
| Missing pointee | Bob sees pointer object CID Y but lacks root CID X or some child chunks. | Whether the group view can show pending / unresolved state without treating the message as invalid. | Sparse-CAS behavior must be a normal group-session state. |
| Envelope independence | The message is wrapped by one candidate grid-envelope variant in one experiment and another in a different experiment. | Whether group-session semantics depend only on resolved payload meaning, not on a chosen grid-envelope winner. | This sim must not backdoor a preferred grid-envelope variant. |
| Historical compatibility | A reader encounters both old `.txt` group-session files and new CAS-backed records. | Whether readers can classify historical evidence and successor records without rewriting either. | The migration contract needs explicit compatibility and provenance rules. |

## Expected Outputs

- A TODO-pipus migration checklist for additive CAS-backed group-session
  specimens.
- A TE-43 interface requirement describing what the group layer receives after
  resolving pointer objects and CAS roots.
- A group-identity pressure case for deciding how known membership,
  pseudonymity, and anonymous participation are represented.
- A guardrail that grid-envelope variants remain separate envelope specimens,
  not hidden dependencies of this group-session successor.
