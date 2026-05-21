# SIM-robot: App semantics and conformance

This simulation is a provisional question home for App Dev feedback items
`FB-dodos`, `FB-hisis`, `FB-kutub`, `FB-gomod`, and `FB-tahof`. It tests what the
guide can say about app semantics and honest conformance before `DR-tuhaz`
settles the stable app-developer contract. Source: `DI-ragaz`.

## Question

Which app-facing semantic patterns are safe as provisional guide prose, and which
must remain blocked until stronger upstream decisions land? Source: `DI-ragaz`.

## Decision Axes

- **Vocabulary status:** promise, assertion, authorship, forwarding,
  conformance, capability, and witness language.
- **Local versus wire identity:** local IDs and storage handles may exist, but
  protocol-boundary identity must be spec-defined.
- **Partial conformance:** useful first slices can be honest if they do not claim
  full draft-spec behavior.
- **Provisional signing:** current signature carriage may be adapter-local until
  grid-envelope/signature decisions freeze.
- **Policy surface:** ingress models and economic patterns may be orientation or
  blocked, depending on what a candidate spec claims.

## Related Root Scenario

- `scenarios/app-semantics-partial-conformance/app-semantics-partial-conformance.md`

## Boundaries

This simulation does not define a universal app SDK, handler ABI, capability
token standard, or final witness format. It exists to keep guide wording honest
while draft specs and DRs remain open. Source: `DI-ragaz`.
