# Device-Bound Agent Physical Effect

## Scenario ID

device-bound-agent-physical-effect

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-nojit`, `FB-tisuf`, and `FB-tulit`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-tuhaz`; `DR-davod`

## Purpose

Exercise apps where Bob's local daemon controls a host-owned physical device and
Alice asks it to perform or report something through a grid protocol.

## Setup

Bob owns a label printer and a temperature sensor. He delegates the front-desk
machine to print shipping labels and the rack sensor to report readings. Alice
sends a print request. Carol audits the receipt. Mallory replays the request
after the label has already been printed.

## Stimulus

The daemon restarts after receiving the request but before all receipts propagate.
It must decide whether replay should print another label, report an already-done
effect, or emit a break-witness.

## Expected Pressure

The candidate design must handle owner/operator identity, host-driver
dependencies, non-idempotent physical effects, at-most-once posture, receipts,
break-witnesses, and 100-year interpretation of evidence after the device and
driver stack are gone.

## Scenario-Specific Evaluation Questions

- What evidence proves the owner delegated the device-bound agent?
- How is a physical effect deduplicated across replay and restart?
- What should a conformance claim say about CUPS, libusb, IIO, i2c, IPP, or
  vendor SDK dependencies?
