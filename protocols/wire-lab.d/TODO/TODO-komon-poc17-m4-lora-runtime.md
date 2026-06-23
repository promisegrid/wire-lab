# TODO-komon: POC17 M4 LoRa runtime

## Status

Planned. Owns the future `implementations/poc17-m4-lora-runtime/` proof of
concept, targeted after POC16. POC17 should add a constrained Cortex-M4/LoRa
runtime agent without using a UART or host bridge as the PromiseGrid message
transport. POC16 planning now lives in
`TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md` and should feed the
constrained-device token, encryption, and payload-shape decisions here. Source:
`DI-lazal`; `DI-ruvot`.

## Decision Intent Log

ID: DI-lazal
Date: 2026-06-16 08:38:09
Status: active
Decision: Plan POC17 as a Cortex-M4/LoRa-shaped runtime-agent proof of concept, not POC16, with the simulated device sending and receiving PromiseGrid messages only through a simulated RFM95/SX127x-style radio path.
Intent: The M4/LoRa runtime is useful enough to deserve its own POC, but it should not be rushed into POC16 or hidden behind a reliable UART/host bridge. The point is to test PromiseGrid on a small-device, lossy-radio, constrained-frame runtime where the device firmware owns its own promises and the host harness only observes or resets the simulation. UART or debug logs may exist for diagnostics, but they must not carry protocol messages, influence delivery, or promise on behalf of the M4 agent.
Constraints: Preserve one top-level semantic action `promise`; preserve `grid([42(pCID), ...])`; no global monitor, global CAS, global route authority, global trust authority, service registry, permission authority, authorization authority, or conformance authority; no UART/host bridge as message transport; no claim of exact Adafruit Feather M4 Express plus RFM95W simulation until pin mappings, SAMD51 peripheral behavior, SPI behavior, radio-driver behavior, and packet semantics are validated; all device-to-peer PromiseGrid traffic must cross the simulated LoRa path; any harness observer must be passive and must not affect trust, routing, ACKs, retransmission, or protocol semantics.
Affects: implementations/poc17-m4-lora-runtime/; protocols/wire-lab.d/TODO/TODO-komon-poc17-m4-lora-runtime.md; DEV-GUIDE-RESOURCES.md; future POC16/POC17 planning docs.

ID: DI-dugog
Date: 2026-06-22 16:10:18 PDT
Status: active
Decision: Keep the POC17 Codex handoff bootstrap in this TODO so a separate Codex operator can start POC17 while POC16 cleanup continues in parallel.
Intent: POC17 needs POC16's executable lessons and session-only operating habits without forcing the next operator to reconstruct them from chat logs. The handoff must preserve the radio-only constrained-device goal, avoid POC16 regressions already discovered, and name POC16 cleanup items that should influence POC17 even if they are not finished in POC16 yet.
Constraints: Do not scaffold or implement POC17 as part of the handoff; do not move POC16 cleanup into POC17; keep this handoff as TODO-local coordination, not a second source of protocol truth; preserve no-UART/no-host-bridge transport, Promise Theory vocabulary, pCID-as-protocol-spec, `grid([42(pCID), ...])`, passive harness/analyzer semantics, and honest simulator-fidelity claims.
Affects: protocols/wire-lab.d/TODO/TODO-komon-poc17-m4-lora-runtime.md; protocols/wire-lab.d/TODO/TODO.md; future implementations/poc17-m4-lora-runtime/.

ID: DI-solih
Date: 2026-06-22 17:14:57 PDT
Status: active
Decision: Treat `/home/angela/lab/bintags` as POC17 prior art and require a language-source thought experiment before choosing how to model the simulated M4 agent.
Intent: `bintags` gives POC17 a concrete Feather M4/RFM9x LoRa reference system with CircuitPython device behavior, a Raspberry Pi LoRa gateway, and a Go host, but POC17 still needs to decide whether its simulated M4 agent should follow the CircuitPython shape, native firmware, or a Go-only simulator model. The TE keeps that choice explicit instead of silently inheriting the prior project's language or wire format.
Constraints: Use `bintags` for behavior, vocabulary, hardware pressure, ACK/retry examples, display/button constraints, and role boundaries; do not inherit its CSV-like text packet format, MQTT bridge semantics, or host-mediated control model as PromiseGrid protocol design; do not edit `/home/angela/lab/bintags` as part of POC17 planning; preserve the radio-only POC17 transport rule.
Affects: protocols/wire-lab.d/TODO/TODO-komon-poc17-m4-lora-runtime.md; /home/angela/lab/bintags/README.md; /home/angela/lab/bintags/devices/m4/m4.py; /home/angela/lab/bintags/devices/pi/pi.py; /home/angela/lab/bintags/application/bt/main.go; future POC17 TE and implementation decisions.

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
- Treat POC16's map, secure-token, CWT-shaped, and encrypted-payload work as
  design pressure that constrained M4/LoRa agents may profile down rather than
  silently omit.

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

## Bintags Prior Art

- Treat `/home/angela/lab/bintags` as concrete prior art for a Feather
  M4/RFM9x-style LoRa device, a Raspberry Pi radio gateway, and a Go host.
  Source: `DI-solih`.
- Use the `bintags` M4 code as behavior evidence for button-driven status
  changes, display refresh pressure, radio send/receive loops, ACK/retry
  behavior, and small object vocabulary such as `Button`, `Display`, `Radio`,
  `Message`, and `Order`. Source: `DI-solih`.
- Do not carry over the `bintags` CSV-like text packet format, MQTT bridge
  semantics, or host-mediated control model as PromiseGrid protocol design.
  POC17 still needs pCID-owned CBOR payloads, radio-only M4 traffic, local
  agent judgment, and passive harness evidence. Source: `DI-solih`.
- Run a TE before choosing the simulated M4 agent language source. The TE must
  compare CircuitPython-shaped behavior from `bintags`, native firmware in C/C++
  or Rust, and a Go-only simulation under the same simulator, radio-fidelity,
  memory, parser/builder, testing, and migration scenarios. Source: `DI-solih`.

## Codex Handoff Bootstrap

This section is the starting context for the next Codex operator. Complete
`komon.14` before editing POC17 files.

### Required First Reads

- Read repo rules and workflow first: root `AGENTS.md`, then this TODO from top
  to bottom.
- Read POC16 decision context: `TODO-zugok-poc16-secure-tokens-maps-encrypted-payloads.md`
  and `TODO-sosoj-poc16-parser-role-followup.md`.
- Read POC16 executable design notes:
  `implementations/poc16-secure-tokens-maps-encrypted-payloads/README.md`,
  `implementations/poc16-secure-tokens-maps-encrypted-payloads/docs/KERNEL-ROLES.md`,
  and `implementations/poc16-secure-tokens-maps-encrypted-payloads/docs/MESSAGE-SHAPES.md`.
- Read `bintags` prior art:
  `/home/angela/lab/bintags/README.md`, `/home/angela/lab/bintags/devices/m4/m4.py`,
  `/home/angela/lab/bintags/devices/pi/pi.py`, and
  `/home/angela/lab/bintags/application/bt/main.go`.
- Read `DEV-GUIDE-RESOURCES.md` for current cross-POC design state, but treat
  POC evidence as pressure and lessons, not final PromiseGrid APIs.

### Non-Negotiables

- POC17 PromiseGrid traffic from the M4 agent crosses only the simulated
  RFM95/SX127x-style LoRa path. UART, stdout, semihosting, debug logs, host
  files, and simulator monitor channels are diagnostics only.
- No component may become a global monitor, trust authority, CAS authority,
  route authority, registry authority, permission service, authorization service,
  conformance authority, or RPC controller.
- pCID is a protocol-spec selector, not an address, app name, operation code,
  message kind, route target, service name, or universal dispatch key.
- Preserve the outer message invariant `grid([42(pCID), ...])`; every slot after
  slot 0 is defined by the pCID spec.
- Keep the top-level semantic action minimal: ordinary protocol behavior should
  be voluntary `promise` payload semantics, not new action kinds.
- Do not claim exact Feather M4 Express plus RFM95W simulation fidelity until pin
  mappings, SAMD51 peripherals, SPI behavior, radio-driver behavior, and packet
  semantics are explicitly modeled and gated.

### POC16 Lessons To Inherit

- Preserve relevant POC16 coverage as a regression floor: parser/builder role
  separation, pCID-owned payload and arity, raw CBOR artifacts, exact-message
  parent links, per-agent sparse CAS, promise-based GC/backpressure, secure
  token pressure, encrypted payload pressure, and route economics.
- Keep kernel roles as promise boundaries. A runtime may collapse roles into one
  process or firmware loop, but the design must still identify who promises
  transport, app delivery, radio access, resource allocation, storage, route
  selection, and event retention.
- Treat radio delivery as transport only. The radio medium may lose, duplicate,
  delay, or bound packets; it must not judge promises or mutate trust.
- Keep ACK/reply correlation based on exact request message hashes and parent
  links, not payload-level RPC request IDs.
- Keep local scarcity separate from peer promise-breaking. Battery, RAM, flash,
  MTU, retry-budget, and airtime pressure can justify local non-commitment, but
  must not automatically lower trust in a peer.
- Preserve raw-message review. Retain exact CBOR envelopes and malformed radio
  bytes for operator review, but keep harness artifacts separate from production
  protocol behavior.
- Preserve honest analyzer language. Passing gates means "POC evidence complete"
  or "candidate for the current POC scope", not production readiness.

### POC16 Cleanup Still Relevant To POC17

- POC16 still needs implementation-local CID-named spec aliases under its
  `docs/protocols/` directory; POC17 should plan CID-named spec aliases from the
  start once its spec docs exist.
- POC16 still needs a guard preventing stale root `docs/protocols/` POC mirrors;
  POC17 should not reintroduce a root-level protocol mirror.
- POC16 analyzer wording around "production-candidate" needs tightening; POC17
  analyzer output must avoid stronger production claims.
- POC16 builder-role behavior needs more real coverage beyond profile/specimen
  evidence; POC17 should make parser/builder or firmware parser/builder behavior
  executable if it claims coverage.
- POC16 needs diagnostic CBOR examples for each active protocol; POC17 should
  include diagnostic renderings for every small-device pCID it introduces.
- POC16 pCID inventory needs auditing against spec docs, parser, builder, and
  runtime use; POC17 should keep that inventory explicit from the beginning.

### Session-Only Operating Habits

- Always run the clean simulator/container command after behavior changes before
  calling a POC change done. If POC17 has no clean command yet, create one before
  claiming executable success.
- Do not describe shortcuts, harness-only behavior, fake WASM/radio behavior,
  simulated security, or diagnostics as production or fully functioning.
- Keep all raw messages intact for later review. Prefer exact CBOR artifacts and
  diagnostic renderings over summarized logs when validating wire behavior.
- Separate harness/analyzer/observer facts from agent-visible protocol facts.
  Analyzer output is design evidence, not a global system view.
- Record approximations explicitly: exact modeled behavior, approximate modeled
  behavior, diagnostics-only behavior, and not-yet-modeled hardware behavior.
- When unsure whether a new concept should be a top-level action, default to a
  pCID-owned promise payload unless a TE/DI proves otherwise.

## Subtasks

- [ ] komon.14 Read and acknowledge the Codex handoff bootstrap before any POC17
  implementation work.
- [ ] komon.15 Run a TE on simulated M4 agent language sources: CircuitPython-shaped
  behavior from `bintags`, native firmware in C/C++ or Rust, and Go-only
  simulation. Compare each option under simulator support, LoRa/radio fidelity,
  constrained-memory behavior, PromiseGrid parser/builder coverage, testability,
  and long-term migration.
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
