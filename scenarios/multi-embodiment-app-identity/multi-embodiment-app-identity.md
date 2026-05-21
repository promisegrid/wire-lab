# Multi-Embodiment App Identity

## Scenario ID

multi-embodiment-app-identity

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-zazon - heterogeneous embodiments but one app identity`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-tuhaz`; `DR-davod`

## Purpose

Exercise one logical app implemented by heterogeneous components that claim the
same pCID-selected contract without pretending each component implements every
part.

## Setup

Alice uses a browser tab with no filesystem access. Bob uses a Neovim plugin
with a long-running Node helper process. Both participate in the same app and
document. Carol reads their conformance claims. Mallory exploits ambiguity
between "same UX app" and "same protocol contract."

## Stimulus

The browser and plugin exchange data through a custom live channel and publish
separate conformance claims naming which parts of the shared contract they
implement.

## Expected Pressure

The candidate design must show how one app stays one app by honoring an explicit
protocol contract or family of contracts, while each embodiment states its own
runtime constraints and implementation scope.

## Scenario-Specific Evaluation Questions

- What makes the browser and plugin one app rather than two unrelated apps?
- What must each component's conformance claim include?
- How does the answer survive host replacement, browser storage loss, helper
  upgrades, and long-term audit?
