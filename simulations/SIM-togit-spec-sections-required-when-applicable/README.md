# SIM-togit-spec-sections-required-when-applicable

This simulation turns the required-when-applicable spec-sections alternative
from `SIM-ranib-spec-requirement-sections` into a concrete candidate specimen.
It tests whether specs should include explanatory sections only when the layer
makes them load-bearing. Source: `DI-fibuv`.

## Design Under Test

Each spec promises to state whether promise vocabulary, 100-year pressure, and
mental-model sections are applicable. Required sections appear only when their
absence would make implementation or review ambiguous.

## Boundaries

This simulation does not define final applicability criteria. It tests whether
conditional requirements reduce boilerplate or let important clarity disappear.
