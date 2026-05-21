# SIM-zisan: Device-bound agent

This simulation is a provisional question home for `FB-nojit`, `FB-tisuf`, and
`FB-tulit`: apps that expose a host-owned physical device, sensor, actuator, or
driver stack to grid peers. Source: `DI-ragaz`.

## Question

What can the guide safely teach about a device-bound PromiseGrid app before the
final app and kernel boundaries settle? Source: `DI-ragaz`.

## Decision Axes

- **Agent role:** a long-running userspace daemon acts as the protocol-facing
  voice of a host-owned device.
- **Owner and operator authority:** owner keys, delegated operator keys, and
  physical custody need explicit evidence.
- **Irreversible effects:** actuator requests need at-most-once posture,
  deduplication, receipts, and break-witnesses for non-idempotent outcomes.
- **Sensor versus actuator:** sensors and actuators share device-bound identity
  and driver-stack concerns; actuators add irreversible-effect pressure.
- **B-side conformance:** code may claim a grid protocol while honestly naming
  host dependencies such as CUPS, libusb, Linux IIO, i2c, IPP, or vendor SDKs.

## Related Root Scenario

- `scenarios/device-bound-agent-physical-effect/device-bound-agent-physical-effect.md`

## Boundaries

This simulation does not freeze a device-agent protocol or driver API. It tests
guide-safe patterns for provisional tutorials and conformance claims. Source:
`DI-ragaz`.
