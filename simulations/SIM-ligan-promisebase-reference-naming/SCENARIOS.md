# Promisebase Reference Naming Scenarios

These scenarios isolate promisebase's human-readable reference problem from the
broader CAS object-model work. They are design pressure, not a final naming
protocol. Source: `DI-tibis`.

## Scenario Matrix

| Scenario | Setup | What to test | Decision pressure |
|---|---|---|---|
| Immutable local nickname | Alice locally calls root CID X `spring-paper` without publishing that name. | Whether local-only convenience names can stay outside protocol identity. | The first profile may avoid reference protocol scope if local nicknames are enough. |
| Published mutable ref | Alice publishes `project/latest` first pointing at root X and later at root Y. | How signed update history, replay protection, and reader expectations work. | Mutable refs are not CAS roots; they need explicit update semantics if adopted. |
| Reference object as CAS object | Alice writes a CBOR reference object whose CID R points at root X and includes a human-readable label. | Whether the reference object is L6 CAS, L7 metadata, or a separate protocol object. | Reference-object identity must not be confused with the target root identity. |
| Name collision | Alice and Bob both publish `release` pointing at different roots. | Whether names are scoped by peer, group, pCID, site, or another authority. | Human-readable names need scope rules or they become a new central registry problem. |
| Maliciously similar name | Mallory publishes a visually confusing or policy-confusing name near Alice's ref. | Whether readers need normalization, warnings, or refusal policy. | Reference naming imports UI and social-trust risks that raw CIDs avoid. |
| Promisebase custom syntax migration | A promisebase-era reference uses non-CID custom syntax for a root. | Whether migration wraps it, rejects it, or maps it into CID-backed reference objects. | Prior-art adoption must be deliberate and not preserve known-bad syntax by accident. |

## Expected Outputs

- A decision packet for TODO-kituj `kituj.5` on whether reference naming belongs
  in L6, L7, a separate protocol, or a later deferral.
- A distinction table for reference name, reference-object CID, pointer-object
  CID, and target root CID.
- A prior-art disposition list for promisebase reference syntax and behavior.
