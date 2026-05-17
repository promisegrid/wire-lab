# SIM-dihiz: Peer Adoption Metadata

This simulation captures the follow-on question from `TODO-nivus` 011.10: the
wire shape of a peer-level adoption promise such as "I, peer P, promise to
behave as pCID X with open-question answers Q7=yes, Q9=variant-B." It is a
standalone design-point simulation, not a final pCID adoption protocol and not a
shared metadata bundle for other simulations. Source: `DI-pukap`.

## Question

How should a PromiseGrid peer advertise, bind, update, and prove the exact spec
pCID plus open-question answers it promises to follow? Source: `DI-pukap`;
`TODO-nivus`.

## Candidate Shapes

- **Structured adoption object:** A content-addressed object records peer,
  spec-pCID, answer set, time/scope, and signature.
- **Promise message:** Adoption is expressed as a normal PromiseGrid promise
  message that can be stored, relayed, and superseded.
- **Spec-side metadata:** The spec or manifest defines the answer vocabulary,
  while each peer publishes only a compact answer binding.
- **Hybrid claim:** A compact peer promise points at a richer adoption object or
  answer-profile object.

## Boundaries

This simulation does not reopen the settled pCID-storage decisions in
`TODO-nivus` and does not choose the adoption metadata wire shape. It gives the
missing peer-level adoption half of the spec-doc-as-promise machinery a
simulation home before a later TE/DR locks it. Source: `DI-pukap`.
