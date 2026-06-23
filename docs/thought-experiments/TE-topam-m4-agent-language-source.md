# TE-topam: M4 agent language source for POC17

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-topam

## Status

decided

## Prior aliases

None. This TE was minted directly under the proquint handle scheme.

## Decision under test

This TE tests what should be the source of language, structure, and runtime
vocabulary for the simulated Cortex-M4 agent in POC17.

POC17 already targets a Feather M4/RFM95-shaped LoRa runtime agent. `DI-solih`
adds `/home/angela/lab/bintags` as prior art because that repo contains a real
Feather M4/RFM9x-style LoRa project: CircuitPython code for the M4 bin tag,
Python code for the Raspberry Pi gateway, and Go code for the host. The open
decision is not whether `bintags` matters. It is how much of its language and
program shape should guide POC17 before the simulator choice is locked.

The decision affects `komon.15`, the later simulator-choice TE in `komon.2`, and
the future POC17 implementation path under
`implementations/poc17-m4-lora-runtime/`.

## Assumptions

- POC17 remains executable design evidence, not a final embedded runtime,
  LoRaWAN implementation, Feather board model, or production API.
- The M4 agent's PromiseGrid traffic must cross only the simulated
  RFM95/SX127x-style LoRa path. UART, stdout, semihosting, debug logs, host
  files, and simulator monitor channels remain diagnostic only.
- The outer message shape remains `grid([42(pCID), ...])`; pCID-owned protocol
  specs define the remaining slots.
- POC16 contributes parser/builder role pressure, exact raw CBOR artifacts,
  sparse CAS behavior, parent links, secure-token pressure, encrypted-payload
  pressure, and the rule that transport code must not parse normal app payloads
  for routing semantics.
- `bintags` contributes real device behavior and vocabulary: `Button`,
  `Display`, `Radio`, `Message`, `Order`, button-driven status updates, LoRa
  send/receive loops, ACK/retry behavior, display refresh pressure, a Pi radio
  gateway, and a Go host.
- `bintags` does not contribute a PromiseGrid wire shape. Its CSV-like text
  packets, MQTT bridge, and host-mediated order control are prior-art pressure,
  not protocol defaults.

## Developer and maintenance model

The named actors are developers and maintainers who need to work with the POC17
code and its descendants over time.

Alice starts the POC17 implementation. She needs a language source that lets her
build the first evidence without hiding the device constraints.

Bob maintains the simulator and analyzer. He needs the implementation to expose
radio path, CBOR, retry, loss, sparse CAS, and local non-commitment evidence in
ways tests can inspect.

Carol ports or compares the work against hardware-shaped runtimes. She needs a
clear trail from the POC17 code back to Feather M4/RFM9x behavior and forward to
native firmware.

Dave works on the production deployment path. He needs to know which POC17
choices were behavior scaffolding, which were simulator shortcuts, and which can
survive into real devices.

Ellen inherits the code after several POC iterations. She needs names,
boundaries, docs, and tests that still explain why the POC17 agent was shaped
the way it was.

Frank integrates constrained devices with stronger peers. He needs POC17 to
avoid hidden host bridges and to keep gateway or storage helpers from promising
on behalf of the device.

Mallory is a security and failure-mode reviewer. She writes malformed-frame,
replay, oversized-packet, stale-token, and hidden-bridge tests to find places
where implementation convenience masks protocol claims.

## Alternatives

### Alt A: CircuitPython-shaped behavior from bintags

Use the `bintags` M4 program as the primary language source for the simulated
agent's behavior and object vocabulary. POC17 would keep names and concepts such
as radio, display, button, message, order/status, ACK, retry, and event loop as
the first human-readable model. The actual simulator implementation could still
be Go, native firmware, or another language after `komon.2` selects the
simulator.

This gives developers a concrete prior system to reason from. It makes the first
POC17 code easier to explain and keeps the work tied to Feather M4/RFM9x
pressure. It also creates a duty to say exactly where POC17 diverges from
`bintags`: CBOR instead of text packets, bounded retries instead of open-ended
retry loops, local PromiseGrid judgment instead of host-mediated workflow, and
radio-only protocol transport.

### Alt B: Native firmware source in C/C++ or Rust

Use native firmware as the source language from the beginning. POC17 would model
the device close to what could eventually run on a Cortex-M4, with tighter
control over memory, packets, interrupts, SPI, and cryptographic library choices.

This gives Carol and Dave the cleanest path toward hardware and production
deployment. It also makes Alice's first milestone harder because she must solve
toolchain, simulator, driver, and board-support problems before the repo gets
evidence about PromiseGrid behavior over a lossy radio.

### Alt C: Go-only simulation source

Model the M4 agent directly in Go as a small-device simulation. POC17 would
reuse the repo's existing Go habits and can test pCID parsing, CBOR artifacts,
sparse CAS, radio loss, retries, and analyzer gates quickly.

This helps Alice and Bob produce testable evidence early. It also creates the
largest risk for Carol, Dave, and Ellen: future readers may confuse a Go
behavior simulator with proof that the design fits a Feather M4/RFM95-class
device. This option needs strong simulator-honesty language and hard analyzer
gates against hidden host delivery.

## Timeline scenario analysis

### T0: POC17 intake and first planning handoff

Alice reads `TODO-komon`, the POC16 notes, and `bintags`. She needs a plan that
lets her start without deciding the full simulator stack too early. Bob wants
the future analyzer requirements to remain visible from day one. Carol asks
whether the model will be useful for a later native port. Ellen wants the TODO
and TE trail to be readable after several more POCs land.

Alt A gives the team a shared starting language: radio, display, button, message,
order/status, ACK, retry, and event loop. It also lets the TE say which
`bintags` habits POC17 must not inherit. Alice can describe the first agent in
plain device terms while deferring the implementation language.

Alt B gives Carol a stronger hardware story, but it forces Alice to settle
firmware and simulator questions before the repo has a small-device PromiseGrid
model. That blocks the planning handoff on toolchain details.

Alt C lets Alice start fastest, but it gives Ellen less durable context. If the
first source language is just "Go objects," the historical reason for the
Feather/RFM9x shape becomes easier to lose.

### T1: First executable POC17 slice

Alice builds the smallest runnable agent. Bob needs a deterministic clean run
that records radio send/receive events, exact CBOR envelopes, parse failures,
retry counts, and local non-commitments. Mallory adds malformed CBOR, wrong pCID,
replay, and oversized-frame tests.

Alt A helps Alice choose behavior names and test stories, but it should not force
CircuitPython as the executable language. The first code can be Go if `komon.2`
selects a Go simulator, while the model still uses `bintags`-shaped local
objects.

Alt B makes the executable slice more faithful to firmware constraints, but it
may delay the first clean run while Alice brings up firmware build and simulated
SPI/radio behavior.

Alt C makes Bob's tests easiest. It needs explicit budget constants and
radio-only gates so the Go runtime does not erase the constrained-device
pressure.

### T2: Simulator choice and fidelity lock

After `komon.2`, Alice and Bob must choose among Renode, QEMU plus a radio model,
Wokwi/custom chips, hardware-in-loop, and pure Go simulation. Carol checks
whether the choice preserves a path to hardware. Dave asks which claims can later
survive into a production deployment.

Alt A composes well with this decision because it separates language source from
implementation substrate. The team can choose Go for the first simulator,
Renode for a later firmware pass, or another option without losing the
`bintags` behavior trail.

Alt B should win only if the simulator TE proves that native firmware can run
soon enough to support the POC17 evidence. If not, native-first turns the
simulator decision into a toolchain gate.

Alt C should win only as an implementation substrate, not as the design language.
If pure Go wins, Bob must add analyzer language that calls the result a
small-device behavior simulator until stronger hardware modeling lands.

### T3: Parser/builder, CBOR, and radio-only gates

POC17 starts inheriting POC16 lessons: pCID-selected parser/builder roles, exact
raw CBOR artifacts, parent links, sparse CAS, bounded local storage, and
non-commitment for unsupported input. Frank checks that the gateway does not
carry messages for the device. Mallory checks that stdout, files, shared memory,
and host callbacks cannot act as hidden protocol transport.

Alt A keeps the code understandable if the `Message` idea becomes a local object
that builds and parses pCID-owned CBOR envelopes. It also forces a clear break
from `bintags` CSV-like text packets and MQTT bridge behavior.

Alt B gives the strongest radio and memory discipline if the firmware/runtime
can support the parser and artifact requirements. It also makes every diagnostic
path harder to keep separate from protocol behavior.

Alt C lets Bob add parser/builder and analyzer gates quickly. It requires extra
discipline because one Go process can accidentally bypass the simulated radio
with direct calls or shared data.

### T4: POC17 hardening and follow-on POCs

Ellen inherits the POC17 code after Alice has moved on. Bob has added analyzer
gates. Carol starts a native-fidelity branch. Frank adds store-and-forward with
stronger peers. Mallory expands failure tests to include stale tokens,
encrypted-payload pressure, duplicate packets, and asymmetric reachability.

Alt A gives Ellen the best maintenance trail. She can read the TE, see why
objects like radio, button, display, message, and status exist, and understand
which parts came from prior art rather than PromiseGrid protocol law.

Alt B gives Carol the best branch target, but only if earlier work kept behavior
clean enough to port. Native firmware should follow a proven behavior model, not
replace it before the model exists.

Alt C gives Bob and Mallory the easiest regression suite. It becomes a liability
if later docs fail to separate the simulator implementation from the M4/LoRa
claim.

### T5: Production pilot with real or near-real devices

Dave prepares a production pilot. Some devices resemble Feather M4/RFM95 boards.
Some stronger peers retain messages or relay over better links. Frank integrates
gateway and storage helpers. Carol checks the port against real hardware
constraints. Ellen reviews whether POC17 claims still match what the code proved.

Alt A remains useful as the behavioral bridge. The production team can trace
radio, button/display, retry, and status behavior back to `bintags` while also
seeing where PromiseGrid changed the wire shape and trust model.

Alt B becomes more important in this phase. Native firmware or a closer
hardware simulator should take over as the implementation source before anyone
claims production device readiness.

Alt C should stay in the test and reference lane. It can keep running as a
regression model, but it should not be the only evidence for deployment on
constrained devices.

### T6: Long-term maintenance and migration

Years later, Ellen and Dave need to decide whether to keep the POC17 line,
supersede it, or port it to a new constrained radio platform. Bob's analyzer has
changed. Carol's hardware branch may have diverged. Frank's gateway code may
support more peer-storage behavior. Mallory still checks for hidden bridges and
overstated fidelity claims.

Alt A gives the most durable explanation of why POC17 exists and how it relates
to real prior work. It does not preserve hardware fidelity by itself, so later
docs must record which POC or production branch proved native constraints.

Alt B gives the best long-term production path if it has landed by this point.
If it never landed, a native-first POC17 would leave little useful evidence
behind.

Alt C gives the most durable automated regression path. It needs a clear label:
behavior simulator, not firmware proof.

## Surviving and rejected alternatives

Alt A survives as the primary language source. It gives developers the most
useful prior-art bridge, keeps the device model understandable, and preserves the
specific Feather/RFM9x pressure that motivated POC17.

Alt B survives as a later fidelity target and possible implementation choice
after `komon.2`, not as the unconditional first source. Native firmware should
become central when the simulator and toolchain can support it without delaying
the core PromiseGrid evidence.

Alt C survives as a possible first implementation substrate and long-term
regression model, not as the source of device language. A Go-only simulator can
build clean evidence quickly, but it should borrow the local behavior vocabulary
from Alt A and carry explicit simulator-honesty gates.

Rejected: Alt B as the unconditional first language source. It creates too much
toolchain and board-support risk before POC17 has proved the radio-only
PromiseGrid behavior.

Rejected: Alt C as the only language source. It would make POC17 too easy to
mistake for a generic Go protocol simulator and too weak as Feather/RFM9x
evidence.

## Conclusions

The recommended DF outcome is: POC17 should use `bintags` CircuitPython-shaped
M4 behavior as the language source for the simulated agent, while allowing the
implementation language to be chosen by `komon.2`.

In practical terms, POC17 should first describe the device with `bintags`-like
local objects and actions: radio receive/send, button status change, display
refresh budget, message parse/build, ACK/retry budget, local state, and
order/status or device-status promises. POC17 must translate the old `Message`
string format into pCID-owned CBOR envelopes and must keep gateway/host behavior
from becoming hidden transport.

`DI-libis` locks the implementation sequence after this TE: POC17 starts with a
Go behavior simulator that uses `bintags` CircuitPython-shaped device vocabulary,
then pursues Rust/Renode as a later fidelity lane after the Go simulator passes
clean radio-only and artifact gates.

## Implications for the repo's open TODOs and pending DIs

- `komon.15` is complete and decided by `DI-libis`.
- `komon.2` should use this TE as an input when comparing Renode, QEMU plus an
  external radio model, Wokwi/custom chips, hardware-in-loop, and pure Go radio
  simulation.
- `komon.3` should still lock remaining naming, fidelity target, path, command,
  package, firmware-language, and runtime-generated path decisions before
  scaffold work begins.
- `komon.4` must not scaffold POC17 until the DF after `komon.2` and this TE
  resolves simulator choice, fidelity target, path approvals, command names,
  package names, firmware language, and runtime-generated paths.
- `DI-solih` remains active. This TE refines how to use `bintags` prior art but
  does not supersede that DI.
- `DEV-GUIDE-RESOURCES.md` should not present this TE as a PromiseGrid API. If
  guide resources later cite POC17, they should cite this TE only as simulator
  language-source pressure for constrained LoRa agents.
