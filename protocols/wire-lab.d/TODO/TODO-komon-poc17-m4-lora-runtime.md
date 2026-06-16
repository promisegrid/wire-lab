# TODO-komon: POC17 M4 LoRa runtime

## Status

Planned. Owns the future `implementations/poc17-m4-lora-runtime/` proof of
concept, targeted after POC16. POC17 should add a constrained Cortex-M4/LoRa
runtime agent without using a UART or host bridge as the PromiseGrid message
transport. Source: `DI-lazal`.

## Decision Intent Log

ID: DI-lazal
Date: 2026-06-16 08:38:09
Status: active
Decision: Plan POC17 as a Cortex-M4/LoRa-shaped runtime-agent proof of concept, not POC16, with the simulated device sending and receiving PromiseGrid messages only through a simulated RFM95/SX127x-style radio path.
Intent: The M4/LoRa runtime is useful enough to deserve its own POC, but it should not be rushed into POC16 or hidden behind a reliable UART/host bridge. The point is to test PromiseGrid on a small-device, lossy-radio, constrained-frame runtime where the device firmware owns its own promises and the host harness only observes or resets the simulation. UART or debug logs may exist for diagnostics, but they must not carry protocol messages, influence delivery, or promise on behalf of the M4 agent.
Constraints: Preserve one top-level semantic action `promise`; preserve `grid([42(pCID), ...])`; no global monitor, global CAS, global route authority, global trust authority, service registry, permission authority, authorization authority, or conformance authority; no UART/host bridge as message transport; no claim of exact Adafruit Feather M4 Express plus RFM95W simulation until pin mappings, SAMD51 peripheral behavior, SPI behavior, radio-driver behavior, and packet semantics are validated; all device-to-peer PromiseGrid traffic must cross the simulated LoRa path; any harness observer must be passive and must not affect trust, routing, ACKs, retransmission, or protocol semantics.
Affects: implementations/poc17-m4-lora-runtime/; protocols/wire-lab.d/TODO/TODO-komon-poc17-m4-lora-runtime.md; DEV-GUIDE-RESOURCES.md; future POC16/POC17 planning docs.

## Scope

- Treat POC17 as executable design evidence for constrained embedded
  PromiseGrid agents, not as a final embedded runtime, radio stack, LoRaWAN
  implementation, Feather board model, or production device API.
- Target a Cortex-M4-class firmware agent shaped by Adafruit Feather M4 Express
  constraints and an RFM95W/SX127x-style SPI LoRa radio, but call it
  "Feather/RFM95-shaped" until simulator fidelity is proven.
- Keep the firmware agent autonomous: it makes its own receive promises, link
  promises, retry promises, storage promises, and local trust judgments.
- Make the radio path the only PromiseGrid message path for the M4 agent.
  UART, stdout, semihosting, debug logs, and simulator monitor channels are
  diagnostics only.
- Preserve the POC superset rule unless an explicit later DI narrows POC17:
  POC17 should inherit the relevant POC15/POC16 protocol lessons rather than
  resetting to a toy transport.

## Architecture Targets

- Simulated M4 firmware owns a small PromiseGrid runtime loop:
  receive radio frame, parse CBOR envelope, match pCID, judge local promise
  constraints, optionally update local event/trust/CAS state, and send an
  envelope response if it chooses to promise a response.
- Simulated RFM95/SX127x-style peripheral owns SPI register/FIFO/IRQ behavior
  sufficient for the firmware driver to send and receive packets.
- The simulated radio medium owns packet delivery effects: constrained MTU,
  loss, duplication, delay, and asymmetric reachability. It must not decide
  trust or promise outcomes.
- Host harness owns build/run/reset/artifact collection only. It may observe
  radio events and retained message artifacts, but it must not relay protocol
  messages for the M4 agent.
- Existing Go/Docker agents may participate as non-M4 peers through an adapter
  that behaves as another radio endpoint, not as a hidden reliable bridge for
  the M4 firmware.

## Protocol Targets

- Define a small-device receive promise profile that the M4 firmware can parse
  with bounded memory and CPU.
- Add `lora_link_v1` or equivalent pCID-owned payload semantics for voluntary
  radio-link promises, including retry budget, frame-size limits, packet-loss
  expectations, and route/relay willingness.
- Add at least one useful device promise, such as `sensor_reading_v1`,
  `device_status_v1`, `store_forward_v1`, or `radio_relay_v1`.
- Keep all message kinds as pCID-owned payload meanings, not top-level wire
  actions and not separate pCIDs unless the whole protocol actually changes.
- Prefer compact pCID-owned CBOR arrays for firmware payloads, but allow maps
  only if a specific pCID spec deliberately chooses them.
- Decide whether the M4 slice can afford native Ed25519, a smaller proof
  profile, delegated transport proof, or unsigned local-link specimens; document
  the consequence as a PromiseGrid design pressure rather than hiding it.

## CAS And State Targets

- Give the M4 firmware a tiny sparse CAS or ring buffer for exact message bytes,
  recent parent links, and selected local state.
- Treat missing parents as normal sparse-store state, not protocol failure.
- Add peer-storage promises so the M4 agent can ask a stronger peer to retain
  selected CAS objects in exchange for local credit, reciprocal relay, or bearer
  capability tokens.
- Keep all retention and GC promise-based: the device promises bounded local
  retention under battery, flash, RAM, packet, and energy constraints.
- Avoid cross-run durable state unless a later DI explicitly asks for it;
  POC17 clean runs should reset simulator state to avoid muddying experiments.

## Radio And Failure Scenarios

- Normal direct packet exchange between the M4 agent and one nearby peer.
- Packet loss where the M4 agent retries only within its promised retry budget.
- MTU pressure where a larger PromiseGrid envelope must be refused,
  fragmented under a pCID-defined profile, or delegated to peer storage.
- Asymmetric reachability where Alice can hear the M4 agent but the M4 agent
  cannot hear Alice, forcing local non-commitment or relay promises.
- Store-and-forward through a trusted relay where each hop promises only its own
  bounded forwarding behavior.
- Battery or energy pressure where the M4 agent declines low-value promises
  without treating local scarcity as peer promise-breaking.
- Malformed CBOR, wrong pCID, bad proof, replayed frame, duplicate packet, and
  prompt-injection text probes that are locally rejected or treated as
  non-commitments.
- Permanent distrust or route exclusion where the M4 agent avoids packets or
  relay promises involving a locally distrusted peer.

## Analyzer And Evidence Targets

- Confirm all M4 PromiseGrid messages crossed the simulated radio path, not
  UART, stdout, shared files, or a host bridge.
- Count exact CBOR envelopes emitted by the firmware and received through radio
  peers.
- Count LoRa packet outcomes: sent, received, lost, duplicated, delayed,
  retried, refused for size, and refused for energy/capacity.
- Count pCID-owned payload coverage for the small-device protocols.
- Count local sparse CAS stores, missing parents, retained objects, and GC
  removals.
- Gate against authority drift: no command/control, permission, authorization,
  global trust, global registry, or monitor-as-authority language.
- Gate simulator honesty: distinguish exact modeled behavior, approximate
  modeled behavior, diagnostics-only behavior, and not-yet-modeled hardware
  behavior.

## Subtasks

- [ ] komon.1 Decide whether POC16 must land first and what POC16 contributes
  that POC17 should inherit.
- [ ] komon.2 Run a TE for M4/LoRa runtime alternatives: Renode custom platform,
  QEMU MPS2 plus external radio model, Wokwi/custom chips, hardware-in-loop, and
  pure Go radio simulation.
- [ ] komon.3 Lock simulator choice, fidelity target, and naming decisions via
  DI before creating `implementations/poc17-m4-lora-runtime/`.
- [ ] komon.4 Scaffold the POC17 directory only after DF approval for paths,
  package names, command names, firmware language, and runtime-generated paths.
- [ ] komon.5 Build the smallest firmware agent that can parse one
  `grid([42(pCID), payload])` or `grid([42(pCID), payload, proof])` profile from
  radio bytes.
- [ ] komon.6 Add the RFM95/SX127x-shaped SPI radio model or chosen equivalent
  and make it the only protocol message path for the M4 firmware.
- [ ] komon.7 Add one non-M4 peer that can exchange PromiseGrid envelopes with
  the simulated radio endpoint without acting as a hidden reliable bridge.
- [ ] komon.8 Add small-device pCID specs and compact payload shapes for link,
  status, and one useful device promise.
- [ ] komon.9 Add bounded retry, packet-loss, MTU, duplicate, replay,
  malformed-frame, energy-pressure, and asymmetric-link scenarios.
- [ ] komon.10 Add tiny sparse CAS, parent-link retention, peer-storage promises,
  and local GC behavior for the M4 agent.
- [ ] komon.11 Add analyzer gates proving radio-only transport, exact CBOR
  artifacts, pCID-owned payloads, sparse CAS behavior, failure handling, and
  no authority drift.
- [ ] komon.12 Document simulator fidelity limits honestly in the README and
  DEV guide resources before any clean-run result is cited.
- [ ] komon.13 Run deterministic tests and a clean container/simulator run after
  implementation, then archive exact commands and key message artifacts.
