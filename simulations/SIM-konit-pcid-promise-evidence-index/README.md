# SIM-konit-pcid-promise-evidence-index

This simulation is the Promise-Theory-first successor to
`SIM-tizad-scoped-conformance-citation-ledger`. It keeps the useful indexing and
citation idea, but reframes the index as evidence about pCID-referenced
promises rather than as a conformance or compliance ledger. Source: `DI-tavaz`.

## Design Under Test

Each durable index entry answers a Promise-Theory question:

- which agent promised that a payload or service followed pCID X;
- what exact evidence object supports that promise or a later observation;
- what a peer observed later about keep, break, refusal, or ambiguity;
- which parts remain local trust judgment for Alice, Bob, or Carol.

The index is descriptive evidence only. It does not certify conformance, grant
authorization, or settle disputes by itself.

## Why this differs from `tizad`

`tizad` kept useful audit structure but still leaned on conformance-manifest and
citation-ledger language. `konit` narrows the purpose: keep exact citations and
durable evidence references, but tie them directly to promises and observed
outcomes.

## Boundaries

This simulation does not define a final guide-wide evidence index. It tests
whether pCID-scoped promise-evidence indexing is a better fit than
conformance-centered ledger language.
