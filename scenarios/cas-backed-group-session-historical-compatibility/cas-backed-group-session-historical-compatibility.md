# Historical compatibility

## Scenario ID

cas-backed-group-session-historical-compatibility

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-jurar-cas-backed-group-session/SCENARIOS.md`
- Source simulation: `SIM-jurar-cas-backed-group-session/`
- Source row/title: Historical compatibility
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-jurar-cas-backed-group-session/`.

## Setup

A reader encounters both old `.txt` group-session files and new CAS-backed records.

## Stimulus

Run the candidate simulation against this source test: Whether readers can classify historical evidence and successor records without rewriting either.

## Expected Pressure

The migration contract needs explicit compatibility and provenance rules.
