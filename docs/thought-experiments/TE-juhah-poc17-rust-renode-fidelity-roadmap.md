# TE-juhah: POC17 Rust/Renode fidelity roadmap

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-juhah

## Status

needs DF

## Prior aliases

None. This TE was minted directly under the proquint handle scheme.

## Decision under test

This TE tests the next development path for POC17 now that the Go behavior
simulator has passed radio-only and artifact gates. The question is not whether
to keep the Go simulator. `DI-libis` already locked that lane as the first
behavior-evidence path and as the fast regression path. The question is what the
next fidelity lane should prove first, and how each path changes the full POC17
roadmap.

The TE specifically includes a Rust port that still runs on Linux, because that
path may reduce Go-runtime assumptions and prove cross-implementation behavior
before the team pays the Renode platform and radio-modeling cost.

This decision affects `komon.25` in
`protocols/wire-lab.d/TODO/TODO-komon-poc17-m4-lora-runtime.md`.

## Assumptions

- POC17 keeps the current Go behavior simulator as the fast deterministic
  regression lane.
- The current Go lane proved PromiseGrid behavior evidence: radio-only
  simulated transport, exact CBOR artifacts, binary pCID CIDs in slot 0,
  bintags-shaped order-status messages, peer-storage capability tokens, restart
  recovery, expected refusals, and no authority drift.
- The Go lane did not prove Feather M4 Express, SAMD51, SPI, RFM95W/RFM95,
  CircuitPython runtime, radio-driver, packet timing, memory, energy,
  regulatory, or production-device fidelity.
- `bintags` remains prior-art vocabulary and behavior pressure, not a wire
  format or runtime dependency.
- Rust/Renode should add stronger evidence only where it is clearer than the Go
  simulator. It should not replace the Go lane just to reimplement the same
  behavior with weaker tests.
- Production readiness still requires later hardware, SPI, radio-driver,
  memory, energy, and packet-semantics evidence.

## Developer and maintenance model

Alice needs the next step to be small enough to start without losing the POC17
behavior already proven. Bob needs analyzer and fixture parity so a second lane
does not become an uncheckable demo. Carol needs a credible route toward
firmware, Renode, and hardware fidelity. Dave needs the roadmap to say which
claims can eventually support deployment. Ellen needs long-term maintainability
and clear boundaries between behavior evidence and fidelity evidence. Frank
needs peer-storage and gateway behavior to stay PromiseGrid-shaped across
implementations. Mallory tests for hidden bridges, overclaimed fidelity, stale
tokens, malformed frames, and host authority drift.

## Alternatives

### Alt A: Rust-on-Linux port first

Port the current Go behavior simulator shape to Rust while still running on
Linux. Keep the same high-level scenario, exact CBOR fixtures, pCID documents,
message CAS artifacts, order-status flow, peer-storage capability flow, restart
recovery, and analyzer-visible event semantics.

This gives Alice and Bob a serious cross-implementation step without Renode
setup. It proves that the protocol core and simulator structure are not
accidentally Go-only. It also gives Carol a better language base for later
firmware work. It does not prove SPI, RFM95/SX127x, SAMD51, CircuitPython
runtime, or real memory and energy behavior.

Roadmap under this path:

1. Preserve the Go lane as the reference behavior and artifact oracle.
2. Create a Rust-on-Linux lane that reads the same config and emits the same
   classes of events and CBOR artifacts.
3. Port CID, CBOR, pCID dispatch, order-status payloads, peer-storage payloads,
   and sparse CAS before adding new features.
4. Compare Rust artifacts against Go fixtures where exact bytes should match,
   and explain any intentional differences.
5. Add analyzer support only after the Rust lane emits enough evidence to be
   checked by the existing gates or clearly parallel gates.
6. Use Rust ownership and explicit buffers to expose packet and storage limits,
   but do not call them measured M4 resource limits.
7. Keep failure cases aligned with the Go lane: malformed CBOR, unknown pCID,
   MTU refusal, loss, delay, replay, asymmetric reachability, stale token, and
   missing parent recovery.
8. Exit when Rust can run the current POC17 scenario on Linux and pass analyzer
   gates equivalent to the Go clean run.
9. Then decide whether to move the Rust code toward Renode firmware, a shared
   Rust protocol crate, or both.
10. Continue to document this as implementation-portability evidence, not
    hardware fidelity.

### Alt B: Rust protocol core plus Go simulator interop

Build only the Rust protocol core first: CID handling, CBOR envelope parsing and
building, embedded pCID document handling, order-status payloads, peer-storage
payloads, lifecycle-token fixture reading where useful, and CBOR diagnostic
helpers. Keep the Go simulator as the only scenario runner for now.

This is lower risk than a full Rust simulator. It gives Bob exact fixture tests
and gives Carol reusable code for a later firmware lane. It delays a second
running POC17 scenario, but it may create the cleanest path to Rust firmware
without duplicating simulator orchestration too early.

Roadmap under this path:

1. Keep the Go lane as the only clean-run scenario.
2. Add Rust tests that consume Go-generated message artifacts and protocol spec
   documents.
3. Require Rust to parse and rebuild selected exact envelopes with binary pCID
   slot 0 and canonical CID text.
4. Add cross-language fixture tests for order-status and peer-storage payloads.
5. Add Rust error tests for malformed CBOR, unknown pCID, bad CID bytes, and
   overlong frames.
6. Do not add Renode until the Rust protocol core can pass fixture parity.
7. Treat resource work as buffer-size and allocation discipline only; do not
   report hardware resource usage.
8. Exit when Rust protocol code can be reused by either a Linux simulator or a
   Renode firmware lane.
9. Then choose between Alt A and Alt C as the next runtime step.
10. Continue to cite the Go lane for behavior evidence and Rust for protocol
    portability evidence.

### Alt C: Renode radio-driver seam first

Create a narrow Rust/Renode slice that proves the firmware-to-radio seam:
firmware writes and reads an RFM95/SX127x-shaped SPI register/FIFO/IRQ model,
and the model emits analyzer-visible packet evidence. Keep full PromiseGrid
behavior in the Go lane until the radio seam is trustworthy.

This targets the biggest thing the Go lane cannot prove: whether a
firmware-shaped runtime can drive a radio-shaped peripheral instead of calling a
Go medium directly. It is useful even before the full protocol stack is ported.
It requires Renode setup and a custom or adapted peripheral model.

Roadmap under this path:

1. Keep the Go lane as the full behavior oracle.
2. Define a minimal Rust firmware loop that can send and receive opaque packet
   bytes through an SPI/FIFO/IRQ seam.
3. Model only the RFM95/SX127x behaviors needed for packet send/receive,
   packet length checks, FIFO movement, IRQ completion, and configured loss or
   delivery events.
4. Emit artifacts that prove bytes crossed the modeled radio seam, not a host
   bridge.
5. Add analyzer gates for SPI/FIFO/IRQ evidence before claiming radio-driver
   fidelity.
6. Do not port peer-storage or lifecycle behavior until the seam is stable.
7. Keep CircuitPython behavior as vocabulary pressure, not as a runtime claim.
8. Exit when Renode can prove a firmware-shaped endpoint exchanged exact packet
   bytes through the modeled radio path under deterministic tests.
9. Then port enough protocol behavior to build real `grid([42(pCID), payload])`
   packets inside the firmware lane.
10. Continue to say this is radio-driver-seam evidence, not full production
    readiness.

### Alt D: Full Rust/Renode protocol run

Port enough POC17 behavior into Rust/Renode to run the current order-status,
peer-storage, restart, and failure scenario inside the fidelity lane.

This would give the strongest single-lane evidence if it succeeds. It is also
the highest-risk path because it combines Rust porting, firmware shape, Renode
platform work, radio modeling, protocol parity, artifacts, and analyzer gates in
one step.

Roadmap under this path:

1. Freeze the current Go lane as the expected behavior.
2. Build a Rust firmware/runtime that owns pCID dispatch, CBOR envelopes,
   order-status, peer-storage, sparse CAS, and restart state.
3. Build or adapt a Renode platform and radio peripheral model.
4. Recreate the full POC17 scenario inside Renode.
5. Emit exact CBOR artifacts and event logs that the analyzer can verify.
6. Add gates for no host bridge, SPI/FIFO/IRQ packet path, and parity with the
   Go scenario.
7. Add memory and packet-size gates only for values actually measured by this
   lane.
8. Exit when the Renode lane can pass the same behavior gates plus stronger
   radio-seam gates.
9. Then consider hardware-in-loop validation.
10. Expect slower delivery and higher maintenance cost than the other paths.

### Alt E: Renode board realism first

Focus first on SAMD51/Feather-like platform description, memory map, clock
shape, pin mapping, and boot/runtime constraints before PromiseGrid protocol
behavior.

This helps future hardware claims and may reduce later surprises. It does not
quickly improve POC17's current PromiseGrid evidence, and it risks spending
time on board shape before proving the radio and protocol seams that matter most
to POC17.

Roadmap under this path:

1. Keep the Go lane as the behavior oracle.
2. Define the minimum board features needed for a Feather/RFM95-shaped claim.
3. Build or document the Renode platform shape and missing pieces.
4. Add memory and boot constraints before protocol behavior.
5. Defer full radio behavior until the board model is credible.
6. Defer order-status and peer-storage until firmware can run predictably.
7. Exit when the platform model can host a small firmware loop with documented
   gaps.
8. Then add Alt C's radio-driver seam.
9. Only then consider full protocol behavior.
10. Treat this as platform-shape evidence, not PromiseGrid or radio evidence.

### Alt F: Hardware-in-loop late validation

Use real or near-real Feather/RFM95 hardware later to validate selected claims.
This is not a good next daily development lane, but it belongs in the full POC17
roadmap because production claims eventually require reality checks.

Roadmap under this path:

1. Keep Go as the fast behavior lane.
2. Keep Rust/Renode or Rust-on-Linux as the repeatable fidelity lane.
3. Add hardware only after protocol fixtures, packet limits, and analyzer
   expectations are stable.
4. Use hardware to validate packet size, radio-driver behavior, memory, energy,
   and operational assumptions.
5. Do not require hardware for ordinary clean-run CI.
6. Record physical setup, firmware version, radio region, antenna assumptions,
   and observed limits.
7. Treat failures as evidence that may revise simulator gates.
8. Exit when hardware confirms or rejects a specific claim already stated by
   the simulator roadmap.
9. Keep hardware results separate from deterministic clean-run artifacts.
10. Use this lane before production-device readiness claims.

## Scenario analysis

### T0: Immediately after the Go clean run

Alice wants a next step that does not throw away the working analyzer. Bob wants
exact fixture reuse. Carol wants movement toward firmware. Dave wants the docs
to stop overclaiming. Ellen wants a clear roadmap.

Alt A is attractive because it creates a second implementation while preserving
the current scenario. Alt B is even smaller and may be best if fixture parity is
the most important first step. Alt C is attractive if the team wants the next
proof to attack the hardware-facing gap directly. Alt D is too broad for this
moment. Alt E is useful but does not improve the most important current gap:
radio-driver behavior. Alt F is premature.

### T1: Parser, CBOR, and pCID parity

Mallory tries malformed CBOR, text pCID slot 0, wrong CIDs, and unknown pCIDs.
Frank checks that peer-storage still uses Bob-issued capability tokens and
request CIDs for correlation.

Alt A and Alt B handle this best because they can reuse current artifacts and
tests before Renode complexity enters. Alt C can defer most protocol parity
while proving packet bytes cross the radio seam. Alt D handles everything but
risks hiding parity problems inside a large port. Alt E and Alt F do not target
this stage well.

### T2: Radio-driver and packet-seam fidelity

Carol asks what proves the firmware did not use a host bridge. Bob wants
evidence of SPI register operations, FIFO movement, IRQ completion, packet
length handling, and deterministic delivery effects.

Alt C is strongest here. Alt D can also prove this, but only after a much larger
port. Alt A and Alt B prepare Rust code for this step but do not prove the seam.
Alt E supports this later by improving platform shape. Alt F validates reality
later but is not deterministic enough for first-line regression.

### T3: Memory, flash, CPU, and energy claims

Dave asks whether POC17 can make constrained-device claims. The current Go lane
reports configured limits only, not measured activity. Ellen needs future docs
to preserve that distinction.

Alt A can expose Rust allocation and buffer choices, but it still runs on Linux.
Alt B can keep buffer sizes explicit in protocol code. Alt C and Alt D can begin
measuring firmware/Renode memory and packet behavior if the model supports it.
Alt E may help with memory-map realism. Alt F is needed before final energy and
deployment claims.

### T4: Full POC17 behavior in a fidelity lane

Frank wants order-status and peer-storage to run outside the Go simulator.
Mallory wants stale-token and missing-parent recovery tests in the stronger
lane.

Alt A gets there sooner than Renode but stays Linux-only. Alt D is the direct
Renode answer, but it should follow either Alt B or Alt C to avoid combining too
many unknowns. Alt C needs a later protocol-port phase. Alt E needs both radio
and protocol work after board modeling.

### T5: Long-term maintenance

Bob keeps the Go clean-run lane fast. Carol maintains the fidelity lane. Ellen
needs docs that explain which lane proves which claim.

A layered path is easiest to maintain: Go remains behavior oracle, Rust protocol
or Rust-on-Linux proves implementation portability, Renode radio seam proves
driver/peripheral fidelity, and hardware-in-loop validates production-facing
claims. A single full Rust/Renode rewrite is harder to maintain unless the team
has already proved the seams independently.

## Surviving and rejected alternatives

Alt A survives as a strong next implementation step because it ports the current
behavior to Rust without taking on Renode immediately. It is especially useful
if the team wants Rust to become the long-term firmware language.

Alt B survives as the lowest-risk first Rust step because it proves exact
protocol parity and produces reusable code for Alt A, Alt C, or Alt D.

Alt C survives as the strongest first Renode step because it proves the
radio-driver seam the Go simulator cannot prove.

Alt D survives as a later integration milestone, not as the recommended next
step. It should follow protocol parity and radio-seam work.

Alt E survives as a supporting platform-realism task, not the first fidelity
milestone. It should not block protocol parity or radio-seam work.

Alt F survives as late validation before production claims, not as the daily
regression or next POC17 lane.

Rejected: replacing the Go simulator with Rust/Renode. The Go lane remains the
fast behavior oracle.

Rejected: claiming production readiness from Rust-on-Linux, Go, or Renode alone.
Production claims require measured hardware-facing evidence and later physical
validation.

## Conclusions

The recommended roadmap is layered:

1. Keep the Go simulator as the behavior oracle and fast regression lane.
2. Build Rust protocol-core fixture parity first if the team wants the smallest
   safe Rust step.
3. Build a Rust-on-Linux scenario port next if the team wants a second full
   implementation before Renode.
4. Build the Renode radio-driver seam before porting the full POC17 scenario
   into Renode.
5. Treat full Rust/Renode POC17 behavior as an integration milestone after the
   protocol core and radio seam are proven.
6. Use hardware-in-loop only for late validation and production-facing claims.

The TE leaves one DF question open: should the next implementation step be Alt
B, Alt A, or Alt C? Alt B is the lowest-risk technical foundation. Alt A gives
the fastest second full implementation. Alt C attacks the most important
hardware-fidelity gap first.

## Implications for the repo's open TODOs and pending DIs

- `komon.25` now has a completed TE but still needs DF before implementation.
- A later DI should lock the chosen first post-Go lane: Rust protocol core,
  Rust-on-Linux simulator, or Renode radio-driver seam.
- `komon.12` should cite this TE when documenting simulator fidelity limits.
- `komon.13` should keep archiving Go clean-run artifacts because the Go lane
  remains the behavior oracle even after Rust/Renode work starts.
- `DEV-GUIDE-RESOURCES.md` must not cite any Rust/Renode path as hardware or
  production fidelity until the specific gates in that path exist and pass.
