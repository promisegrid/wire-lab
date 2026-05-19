# Refusal or non-response

## Scenario ID

chunk-feed-replication-refusal-or-non-response

## Source / Provenance

- Source type: mined simulation row
- Source path: `simulations/SIM-zazit-chunk-feed-replication/SCENARIOS.md`
- Source simulation: `SIM-zazit-chunk-feed-replication/`
- Source row/title: Refusal or non-response
- Source DI / TODO / TE: `DI-faros`; `DI-vabor`; `DI-botup`; `DI-nanih`

## Purpose

Transform the source row into reusable root comparison pressure so multiple
simulations can be evaluated against the same question rather than only inside
`SIM-zazit-chunk-feed-replication/`.

## Setup

Alice refuses to send C or never answers.

## Stimulus

Run the candidate simulation against this source test: Whether refusal is a normal observed outcome that can feed future local decisions.

## Expected Pressure

The feed protocol needs space for refusal and timeout outcomes without treating every miss as corruption.
