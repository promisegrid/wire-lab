# TE-juzif: POC17 simulator choice and timing

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-juzif

## Status

decided

## Prior aliases

None. This TE was minted directly under the proquint handle scheme.

## Decision under test

This TE tests which simulator path should drive POC17 and how early that
simulator should appear in the implementation sequence.

`TODO-komon` already lists the candidate runtime alternatives: Renode custom
platform, QEMU MPS2 plus external radio model, Wokwi/custom chips,
hardware-in-loop, and pure Go radio simulation. `TE-topam` adds the key
interaction: the M4 agent language source can be `bintags` CircuitPython-shaped
behavior even when the first implementation substrate is Go, native firmware, or
another simulator language. The user preference for this TE is to get a
simulator in earlier, as long as POC17 remains honest about what the simulator
does and does not prove.

The decision affects `komon.2`, `komon.3`, and the first scaffold under
`implementations/poc17-m4-lora-runtime/`.

## Assumptions

- POC17 needs executable evidence early. A paper-only simulator choice does not
  satisfy the TODO.
- The first simulator must prove radio-only PromiseGrid traffic for the M4
  agent, exact CBOR artifacts, pCID-owned payloads, bounded retry/loss/MTU
  behavior, sparse local state, and no hidden host bridge.
- The first simulator does not need to prove exact Feather M4 Express plus
  RFM95W fidelity. Exact board, SAMD51, SPI, driver, and packet semantics need a
  later locked fidelity target before the repo can claim them.
- `TE-topam` recommends `bintags` CircuitPython-shaped behavior as the M4 agent
  language source, but it leaves the implementation language for later DF.
- Renode documentation shows relevant features for later fidelity work: platform
  descriptions, peripheral modeling, testing, wireless media, configurable
  packet loss, and deterministic seeds. As of 2026-06-23, its built-in wireless
  docs describe IEEE 802.15.4 and BLE media, not LoRa as a built-in radio.
- QEMU documents Cortex-M4-capable MPS2 board models, but those boards are not
  Feather M4 Express boards and do not provide an RFM95/SX127x model by default.
- Wokwi custom chips can model new simulated hardware with a C/Rust/WASM-style
  chip API, including SPI-facing devices, but the API is marked beta.

## Developer and maintenance model

Alice starts POC17 and needs a simulator path that gets to a clean run quickly.
Bob owns analyzer and test evidence. Carol tracks the later firmware/hardware
fidelity path. Dave works toward production deployment. Ellen inherits the code
after several POCs and needs the simulator story to remain readable. Frank
integrates stronger peers, relays, and peer storage. Mallory writes tests that
look for hidden bridges, overclaimed fidelity, malformed frames, and stale trust
or token behavior.

## Alternatives

### Alt A: Renode custom platform first

Build POC17 around Renode from the beginning, with a custom Feather/RFM95-shaped
platform and a modeled radio peripheral.

This gives Carol the clearest path toward firmware-like evidence. Renode's
platform and peripheral modeling features make it a serious candidate for a
later high-fidelity pass. It also risks delaying Alice and Bob: the first POC17
milestone would need Renode setup, a custom platform, radio/peripheral modeling,
and firmware decisions before the PromiseGrid behavior can run.

### Alt B: QEMU MPS2 plus external radio model first

Run Cortex-M4 firmware on a QEMU MPS2 board and connect it to a separate
RFM95/SX127x-shaped radio model.

This gives a native-firmware-flavored path and uses a well-known emulator. It
also starts from a board that is not the target Feather M4 and still leaves the
radio model outside QEMU. POC17 would risk proving "firmware can run on a
generic MPS2-like Cortex-M4 board" before proving the LoRa/PromiseGrid behavior.

### Alt C: Wokwi/custom chips first

Model the device and radio in Wokwi, using custom chips for missing hardware
pieces such as the RFM95/SX127x-style peripheral.

This may make board-and-peripheral visualization attractive and gives a path to
custom SPI devices. It also adds beta API risk, browser/cloud workflow risk, and
unclear fit with wire-lab's deterministic clean-run and artifact requirements.

### Alt D: Hardware-in-loop first

Use real or near-real Feather/RFM95 hardware early, with the harness driving and
observing the run.

This gives the strongest reality check. It is also the worst fit for a first
repo-level POC17 simulator: it adds hardware availability, flaky physical radio
conditions, board flashing, and reproducibility issues before Bob can build
deterministic analyzer gates.

### Alt E: Pure Go radio simulator first

Build the first POC17 simulator in Go, using `bintags` CircuitPython-shaped
behavior as the M4 agent language source. Model the radio path, RFM95-shaped
packet constraints, loss, retry, MTU, asymmetric reachability, malformed frames,
and artifacts directly in the repo.

This gives Alice and Bob the earliest clean run and the best fit with existing
wire-lab test and analyzer habits. It is the easiest path for proving the
PromiseGrid behavior before choosing hardware fidelity. It also creates the
highest overclaim risk: the docs and analyzer must say "behavior simulator"
until a later Renode/QEMU/Wokwi/hardware pass proves stronger device fidelity.

## Timeline scenario analysis

### T0: First POC17 scaffold

Alice has just finished the required reads. She needs to create the smallest
directory that can run one M4-shaped agent and one non-M4 radio peer. Bob wants
testable artifacts immediately. Carol wants the scaffold not to block a later
native-firmware path.

Alt A makes the scaffold look like the later fidelity target, but Alice must
write simulator plumbing before the first PromiseGrid behavior exists.

Alt B starts with a real Cortex-M4 CPU target, but the external radio model and
board mismatch become the first blockers.

Alt C may be quick for visual hardware experiments, but it makes repo-local
headless clean runs and exact artifact collection uncertain.

Alt D fails the first-scaffold test because it depends on physical hardware and
radio conditions.

Alt E fits the early-simulator preference best. Alice can scaffold a clean-run
simulator immediately, Bob can add analyzer gates immediately, and Carol can keep
the code honest by making the README reserve hardware-fidelity claims.

### T1: First radio-only PromiseGrid clean run

Alice sends exact CBOR envelopes through the simulated radio path. Bob verifies
that no UART, stdout, shared file, host callback, or direct function call carried
the device's protocol messages. Mallory adds malformed CBOR, wrong pCID, replay,
duplicate, and oversized-frame tests.

Alt A can prove this well once the custom platform and radio model exist. The
risk is that POC17 spends too long proving simulator integration before proving
radio-only PromiseGrid behavior.

Alt B can prove firmware-side parse behavior, but the external radio model is
still the part POC17 most needs to trust.

Alt C can model SPI-facing radio behavior, but the beta custom-chip path may
slow Bob's artifact and analyzer work.

Alt D can prove real radio behavior but cannot give deterministic clean runs
early.

Alt E proves the core invariant fastest. It must add hard negative gates against
hidden bridge paths because Go makes those shortcuts easy.

### T2: Language choice pressure

After `TE-topam`, the team wants `bintags` CircuitPython-shaped behavior as the
language source. Carol asks whether this blocks native firmware later. Ellen asks
whether future maintainers will understand why the first code is not
CircuitPython.

Alt A can host native firmware later, but using it first may force the language
choice too early.

Alt B forces native firmware pressure early and makes CircuitPython-shaped
behavior more of a documentation convention than a code shape.

Alt C may be able to model CircuitPython-like device behavior, but it adds a
separate UI/tooling layer.

Alt D says little about language source; real hardware could run several
firmware stacks, but it does not help Alice choose the first executable shape.

Alt E composes best with `TE-topam`: the first Go simulator can use `bintags`
local objects and behavior names while preserving a later native-fidelity path.

### T3: Parser/builder and raw-artifact pressure from POC16

POC17 inherits POC16's corrected parser/builder boundary, exact CBOR artifacts,
parent links, sparse CAS, and local non-commitment behavior. Bob needs tests for
all of these before POC17 grows a hardware simulator. Mallory tries to smuggle a
universal route field or host callback into the transport path.

Alt A can eventually model this in a stronger runtime, but early progress
depends on peripheral and firmware work.

Alt B risks focusing on firmware plumbing before parser/builder evidence.

Alt C risks focusing on chip simulation before parser/builder evidence.

Alt D is too hard to make deterministic at this stage.

Alt E gives the best first analyzer story. It can make the parser/builder split,
raw CBOR, missing parents, retry budgets, and radio-only gates executable before
hardware fidelity work begins.

### T4: Mid-POC fidelity upgrade

After the Go behavior simulator passes, Carol asks for a stronger hardware path.
Dave wants to know whether POC17 can still move toward production. Alice and Bob
do not want the already-proven analyzer gates to disappear.

Alt A becomes attractive here. Renode can be introduced as a second simulator
lane if the Go lane has already fixed the protocol and analyzer contract.

Alt B can be reconsidered if native firmware becomes the main question and if
the board mismatch is acceptable for a firmware-only experiment.

Alt C can be reconsidered if a visual/custom-chip workflow helps validate SPI or
display behavior that the Go simulator approximated.

Alt D can be introduced as a late validation lane, not as the clean-run source
of truth.

Alt E should remain as the fast regression lane even if Alt A, B, C, or D is
added later.

### T5: Production pilot preparation

Dave prepares a pilot with real constrained devices. Carol evaluates what POC17
proved about firmware, SPI, radio, energy, and memory. Frank integrates relays
and peer storage. Ellen checks whether old POC17 docs overstated the first
simulator's fidelity.

Alt A is the best surviving candidate for pre-production simulator fidelity if
its custom platform and radio model have landed.

Alt B may help if the production firmware is close enough to a generic
Cortex-M4/MPS2 path, but it still needs target-board caveats.

Alt C may help with custom peripherals, displays, or demos, but its beta custom
chip layer should not be the only production-readiness evidence.

Alt D becomes necessary before real deployment claims, but it should supplement
deterministic simulation instead of replacing it.

Alt E remains useful for regression and protocol behavior, but by itself it
should not support a production-device readiness claim.

### T6: Long-term maintenance

Years later, Ellen and Bob still need fast tests. Carol may maintain a Renode or
native-firmware branch. Dave may support real deployments. Mallory still checks
for hidden host bridges and overclaimed simulator fidelity.

Alt A is strong if the team paid the modeling cost and kept it in CI.

Alt B is strong only if the firmware branch stayed active and did not fork away
from the behavior tests.

Alt C is useful if Wokwi remains available and its custom-chip API stabilizes,
but it is a riskier long-term source of repo-local evidence.

Alt D is useful for periodic validation, but too brittle as the daily regression
path.

Alt E is the most durable daily test lane. Its long-term weakness remains the
same: it must be labeled as behavioral evidence until paired with stronger
fidelity evidence.

## Surviving and rejected alternatives

Alt E survives as the recommended first simulator. It best matches the stated
preference to get a simulator in early, and it gives the fastest path to
radio-only transport gates, exact CBOR artifacts, loss/retry/MTU scenarios,
sparse CAS, and no-hidden-bridge analyzer evidence.

Alt A survives as the recommended second simulator lane or mid-POC fidelity
upgrade. Renode looks like the best place to spend later hardware-modeling
effort because it supports custom platforms, peripheral modeling, deterministic
testing, and wireless-medium concepts. POC17 should not claim LoRa fidelity from
Renode until an RFM95/SX127x-style model exists and is gated.

Alt B survives as a narrower firmware experiment, not the primary POC17
simulator. QEMU MPS2 can exercise Cortex-M4 firmware shape, but the board and
radio mismatch make it less useful than Renode for the eventual Feather/RFM95
story and less useful than Go for early protocol evidence.

Alt C survives as an optional custom-peripheral or demo lane. It should not be
the first POC17 simulator because the custom-chip API and browser/cloud workflow
add uncertainty to clean repo-local runs.

Alt D survives as late validation before production claims. It should not be the
first simulator or the daily regression source.

Rejected: Alt A, B, C, or D as the first required simulator. Each one front-loads
hardware, toolchain, or reproducibility work before POC17 proves the
PromiseGrid-over-lossy-radio behavior.

Rejected: Alt E as the only simulator for all future claims. Pure Go can prove
behavior and analyzer gates, but not exact Feather M4/RFM95 firmware or hardware
fidelity.

## Conclusions

The recommended DF outcome is: POC17 should implement an early pure-Go behavior
simulator first, using `bintags` CircuitPython-shaped behavior as the device
language source, and should plan a later Renode/custom-platform fidelity lane
once the first clean-run gates pass.

The first implementation should not wait for native firmware, Renode, Wokwi, or
hardware-in-loop. It should create a deterministic simulator early in POC17 that
models the M4 agent, one non-M4 radio peer, an RFM95/SX127x-shaped packet path,
loss/duplicate/delay/MTU/asymmetry, exact CBOR artifacts, sparse local state,
and analyzer gates against hidden host transport.

`DI-libis` locks the answer: POC17 starts with a pure-Go behavior simulator and
plans Rust/Renode as the later fidelity lane. QEMU, Wokwi, and hardware-in-loop
remain narrower follow-on options, not blockers for the first simulator.

## Implications for the repo's open TODOs and pending DIs

- `komon.2` is complete and decided by `DI-libis`.
- `komon.3` should lock simulator sequence, fidelity target, naming decisions,
  and path decisions before creating `implementations/poc17-m4-lora-runtime/`.
- `komon.4` should scaffold a clean-run simulator early if DF accepts this TE's
  recommendation.
- `komon.6` should treat the first RFM95/SX127x-shaped radio as a behavior
  model unless a later DI locks a stronger Renode, firmware, Wokwi, or hardware
  model.
- `TE-topam` remains compatible with this TE: `bintags` supplies the device
  language source; the first implementation substrate can still be Go.
- `DEV-GUIDE-RESOURCES.md` should not cite POC17 simulator results as hardware
  fidelity unless a later fidelity lane proves and gates that claim.
