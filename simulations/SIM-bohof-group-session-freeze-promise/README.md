# SIM-bohof: Group-Session Freeze Promise

This simulation captures the protocol/specimen question in `TODO-bisur` 012.8:
what it means to freeze a group-session specimen and what a
`merge-group-transport-spec` promise must claim. It is a standalone
design-point simulation, not a freeze action, not a signed merge promise, and
not a privileged group-session home. Source: `DI-pukap`.

## Question

What evidence, scope, and promise wording are required before a group-session
specimen can be treated as frozen or merged, especially when outer/feed rules
and group-session semantics are separate concerns? Source: `DI-pukap`;
`TODO-bisur`.

## Candidate Shapes

- **Specimen freeze evidence:** A concrete message DAG, verified CIDs, and
  reproducible parser/checker behavior are sufficient evidence for v0 freeze.
- **Two-surface freeze gate:** Outer/feed rules freeze separately from
  group-session semantics, and the merge promise names both scopes.
- **Human merge promise:** Steve signs a plain-English
  `merge-group-transport-spec` promise that points at exact frozen artifacts.
- **Deferred freeze:** The specimen remains provisional until cryptographic
  promise tooling and frozen pCID references exist.

## Boundaries

This simulation does not freeze any spec, sign any merge promise, or decide
whether the historical `grid <pcid>` envelope becomes canonical PromiseGrid
wire format. It gives the open freeze-gate question a simulation home so the
group-session lineage can evolve independently. Source: `DI-pukap`.
