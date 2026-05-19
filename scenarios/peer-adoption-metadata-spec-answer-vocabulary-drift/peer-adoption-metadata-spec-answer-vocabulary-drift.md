# Spec answer vocabulary drift

## Scenario ID

peer-adoption-metadata-spec-answer-vocabulary-drift

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-dihiz-peer-adoption-metadata/SCENARIOS.md`
- Source simulation: `SIM-dihiz-peer-adoption-metadata/`
- Source row/title: Spec answer vocabulary drift
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-dihiz-peer-adoption-metadata/`.

## Setup

A later frozen spec supersedes pCID X and renames or removes Q9.

## Stimulus

Run the candidate simulation against this source test: Whether answer keys are spec-local, profile-local, globally named, or mapped through migration records.

## Expected Pressure

Adoption metadata must survive spec evolution without making old claims ambiguous.
