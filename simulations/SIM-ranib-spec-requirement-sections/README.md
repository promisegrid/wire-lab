# SIM-ranib: Spec Requirement Sections

This simulation captures the protocol/specimen question from `TODO-kulih` 010.9
and `DR-robon`: whether PromiseGrid protocol specs should require
layer-specific promise-vocabulary sections, 100-year pressure-test sections, and
layperson-readable mental-model / easy-implementation summaries. It is a
standalone design-point simulation, not an answer to `DR-robon` and not a final
spec template. Source: `DI-pukap`.

## Question

Which explanatory sections should be required in PromiseGrid protocol specs so
the specs remain precise for implementers, understandable to laypeople, and
durable across long horizons? Source: `DI-pukap`; `TODO-kulih`; `DR-robon`.

## Candidate Shapes

- **Required sections:** Every protocol spec carries promise vocabulary,
  100-year pressure test, and layperson/easy-implementation summaries.
- **Required when applicable:** Specs must include the sections only when the
  layer makes them load-bearing.
- **Guide-only prose:** The development guide owns the summaries while specs
  stay minimal and cite guide prose.
- **Split template:** Specs require concise normative hooks, while longer
  mental-model prose lives in companion docs.

## Boundaries

This simulation does not answer `DR-robon`, edit a spec template, or decide the
final PromiseGrid document structure. It supplies scenario pressure so
`TODO-kulih` / TE-nibar can decide whether these sections become required spec
obligations. Source: `DI-pukap`.
