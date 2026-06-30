# TE-juhah: POC17 M4-first Rust/Renode fidelity roadmap

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## Rewrite note

This TE was substantively rewritten in place on 2026-06-29 by explicit user
direction. The normal TE editing policy would treat this as Cat-5 supersedence,
but the user stated that `TE-juhah` should be rewritten in place because the
original draft did not match POC17's intent. The TE ID and filename are
preserved; the current body is the M4-first roadmap analysis.

## TE ID

TE-juhah

## Status

decided

## Prior aliases

None. This TE was minted directly under the proquint handle scheme.

## Decision under test

This TE tests the next development path for POC17 after the Go behavior
simulator passed radio-only and artifact gates.

The main purpose of POC17 is to find out whether a PromiseGrid node can run on a
Cortex-M4-class device and communicate over LoRa-class hardware. The Go
simulator is useful behavior evidence, but it does not answer the hardware-fit
question. The next roadmap must therefore prioritize M4 compatibility first:
memory model, allocation discipline, firmware shape, boot/reset behavior,
bounded packet handling, and the path to an RFM95/SX127x-style LoRa radio.

Radio-driver fidelity and protocol parity remain necessary. They are supporting
gates, not the organizing purpose. Protocol work matters because the M4 node
must still be a PromiseGrid node. Radio work matters because LoRa is the target
transport. But neither should outrank the question "can this node actually fit
and run on M4-class hardware?"

This decision affects `komon.25` in
`protocols/wire-lab.d/TODO/TODO-komon-poc17-m4-lora-runtime.md`.

## Assumptions

- POC17 keeps the current Go behavior simulator as the fast deterministic
  behavior oracle.
- The current Go lane proved radio-only simulated transport, exact CBOR
  artifacts, binary pCID CIDs in slot 0, bintags-shaped order-status messages,
  peer-storage capability tokens, restart recovery, expected refusals, and no
  authority drift.
- The Go lane did not prove M4 fit, firmware viability, SAMD51 behavior,
  RFM95W/RFM95 behavior, SPI/FIFO/IRQ behavior, CircuitPython runtime behavior,
  measured memory, measured energy, packet timing, regulatory behavior, or
  production readiness.
- `bintags` remains prior-art vocabulary and hardware pressure. It does not
  define the PromiseGrid wire format and does not require CircuitPython as the
  next runtime.
- Rust is valuable only if it is shaped toward embedded constraints. A
  comfortable Rust-on-Linux simulator that depends on heap-rich Linux behavior
  would not answer the POC17 hardware question.
- Renode is valuable only if it proves M4/platform/radio facts that the Go lane
  cannot prove.
- Any reported resource usage must be measured. Configured limits may be
  documented, but fake RAM, flash, CPU, radio airtime, or energy use is not
  evidence.

## Developer and maintenance model

Alice needs the next step to attack the real POC17 question without throwing
away the clean Go behavior lane. Bob needs analyzer evidence that distinguishes
M4-fit claims from protocol parity claims. Carol needs a credible path from Go
behavior evidence to Rust firmware and Renode. Dave needs to know which claims
can support eventual deployment. Ellen needs the roadmap to remain readable
after multiple POCs. Frank needs peer-storage and gateway behavior to stay
PromiseGrid-shaped across implementations. Mallory tests for hidden host
bridges, overclaimed hardware fidelity, malformed frames, stale tokens, memory
growth, and authority drift.

## Alternatives

### Alt A: M4 viability wedge first

Build the smallest embedded-shaped Rust core that could plausibly become the
M4 node. It may run under Linux tests at first, but it is designed as firmware
code: bounded buffers, explicit maximum frame sizes, no hidden heap dependency,
clear panic/reset behavior, small sparse state, and no OS services in the node
core.

This path puts the main POC17 question first. It does not try to port the whole
Go simulator. It asks whether the node's core duties can fit an M4-shaped
runtime: receive bytes, parse `grid([42(pCID), payload])`, judge a small set of
pCIDs, build bounded responses, retain sparse local state, and expose enough
evidence for tests.

Roadmap under this path:

1. Keep the Go simulator as the behavior oracle.
2. Define a Rust node core with M4-shaped constraints: fixed buffers, bounded
   state, explicit error paths, and no required heap.
3. Implement only the pCID/CBOR pieces needed for the current POC17 scenario.
4. Use Linux only as a test harness around the constrained core.
5. Prove fixture parity for selected Go CBOR artifacts and pCID documents.
6. Add tests for overlong frames, malformed CBOR, unknown pCIDs, missing
   parents, stale tokens, and restart state within fixed limits.
7. Record configured limits and any actually measured allocations or stack-like
   usage separately.
8. Exit when the Rust core can run the minimum node loop under M4-style limits
   and prove it can build/parse the current critical message shapes.
9. Move the same core toward Renode after this wedge passes.
10. Keep this claim narrow: M4 viability evidence, not yet LoRa radio fidelity.

### Alt B: Rust-on-Linux full constrained port

Port the current Go behavior scenario to Rust on Linux, but require the code to
be structured as an embedded-constrained implementation rather than a generic
Linux simulator.

This gives a second implementation faster than Renode and may expose Go-specific
assumptions. It is useful only if it keeps M4 pressure. If it starts by
recreating the whole Go simulator with dynamic maps, host-rich files, large
allocations, and easy process assumptions, it distracts from POC17's purpose.

Roadmap under this path:

1. Preserve the Go lane as the behavior oracle.
2. Port scenario behavior only around a constrained Rust node core.
3. Keep artifact writing, fixture comparison, and diagnostics in the Linux
   harness, not in the M4-shaped node core.
4. Require exact CBOR/pCID parity for selected artifacts.
5. Keep memory and buffer limits visible in tests.
6. Avoid Linux-only node assumptions such as unbounded files, threads, dynamic
   queues, or host callbacks.
7. Exit when Rust-on-Linux can run the POC17 scenario while preserving embedded
   constraints at the node boundary.
8. Then move the node core into Renode.
9. Treat the Linux lane as portability evidence only.
10. Do not claim hardware fidelity from this path alone.

### Alt C: Rust protocol core first

Build only the Rust CID/CBOR/pCID protocol core first. Use Go-generated message
artifacts and RFC-like pCID documents as fixtures.

This is a useful support step, but it is not the main POC17 roadmap by itself.
It proves that Rust can speak selected PromiseGrid wire shapes. It does not
prove that a node can fit on M4 hardware or talk over LoRa.

Roadmap under this path:

1. Keep the Go lane as the only behavior scenario.
2. Parse and rebuild selected exact CBOR envelopes.
3. Check binary CID bytes in slot 0 and canonical CID text in diagnostics.
4. Cover order-status and peer-storage payloads needed by the current scenario.
5. Keep APIs compatible with `no_std` or clearly isolate any host-only pieces.
6. Add malformed-CBOR, bad-CID, unknown-pCID, and overlong-frame tests.
7. Exit when the protocol core can be used by Alt A or Alt D without redesign.
8. Do not stop here. Protocol parity is a prerequisite, not the POC17 goal.
9. Use this path if the M4 viability wedge needs a clean parser foundation.
10. Avoid broad protocol work that is not needed by the M4 node.

### Alt D: Renode M4 platform first

Bring up a minimal M4-shaped Renode target before radio modeling or full
PromiseGrid behavior. Prove firmware boot, reset/panic behavior, memory layout,
artifact extraction, and the ability to run a tiny node loop under deterministic
tests.

This path attacks M4 compatibility directly. It may be the right first Renode
step because the radio model is not useful until the firmware runtime is
credible.

Roadmap under this path:

1. Keep Go as the behavior oracle.
2. Define the minimum M4 platform facts POC17 needs to claim: memory map,
   reset behavior, firmware entry, panic path, and diagnostic extraction.
3. Run a minimal Rust firmware loop in Renode.
4. Prove fixed-buffer packet handling before adding real radio semantics.
5. Measure or bound memory where Renode and tooling can actually support it.
6. Keep diagnostics separate from protocol transport.
7. Exit when a firmware-shaped node loop can run deterministically in Renode
   with documented M4 constraints and no host bridge.
8. Then add the RFM95/SX127x seam.
9. Do not port full peer-storage until the M4 platform is stable.
10. Treat this as M4 platform evidence, not full LoRa evidence.

### Alt E: Renode LoRa radio seam next

Add an RFM95/SX127x-shaped SPI/FIFO/IRQ model after the M4 runtime shape is
credible. Firmware should exchange packet bytes through the modeled radio path,
not through a host callback or shared simulator shortcut.

This is necessary for POC17, but it should not outrank M4 viability. A beautiful
radio seam is not enough if the node core cannot fit the M4 runtime.

Roadmap under this path:

1. Start after Alt A or Alt D has made the firmware/node shape credible.
2. Model only the RFM95/SX127x behaviors needed for packet send/receive,
   length checks, FIFO movement, IRQ completion, and configured delivery
   effects.
3. Emit evidence that bytes crossed SPI/FIFO/IRQ, not a host bridge.
4. Add MTU, loss, delay, replay, malformed-frame, and asymmetric-link evidence.
5. Keep PromiseGrid payload parsing inside the node, not in the radio model.
6. Exit when the firmware-shaped endpoint exchanges exact packet bytes through
   the modeled radio path under deterministic tests.
7. Then port enough PromiseGrid behavior into the firmware lane.
8. Continue to say this is radio-seam evidence, not full production readiness.
9. Add real timing and regulatory evidence only when actually modeled.
10. Use hardware-in-loop later to validate simulator assumptions.

### Alt F: Full Rust/Renode POC17 scenario

Run order status, peer storage, restart recovery, and failure cases inside the
Rust/Renode lane.

This is the right integration milestone, but it is too large as the first next
step. It should follow M4 viability, protocol fixture parity, and the radio seam
unless a later decision explicitly accepts the combined risk.

Roadmap under this path:

1. Freeze the Go lane as the behavior oracle.
2. Reuse the M4-shaped Rust core.
3. Reuse the Renode M4 platform and radio seam.
4. Add selected pCID behavior: order status, peer storage, device status, and
   link evidence.
5. Emit exact CBOR artifacts and analyzer-visible events.
6. Add no-host-bridge gates, radio-seam gates, and measured-resource gates
   where measurement exists.
7. Exit when the Renode lane can pass the current behavior gates plus stronger
   M4 and radio gates.
8. Then consider hardware-in-loop.
9. Keep the Go lane as the fast regression lane.
10. Treat this as strong simulator fidelity, not final production proof.

### Alt G: Hardware-in-loop late validation

Use real or near-real M4/RFM95 hardware after the deterministic lanes have
stable behavior, M4, and radio-seam evidence.

This is eventually required for deployment-facing claims. It is not a good next
daily development lane because it adds physical setup, flashing, antenna,
region, and flaky-environment variables before the simulator questions are
settled.

Roadmap under this path:

1. Keep Go as behavior oracle.
2. Keep Rust/Renode as deterministic fidelity lane.
3. Add hardware only after packet limits, resource limits, and analyzer
   expectations are stable.
4. Record board, firmware, radio region, antenna, power, and observed limits.
5. Use hardware failures to revise simulator assumptions.
6. Do not require hardware for ordinary clean-run CI.
7. Exit when hardware confirms or rejects a specific claim already stated by a
   simulator gate.
8. Keep hardware results separate from deterministic artifacts.
9. Use this lane before production-device readiness claims.
10. Do not let hardware tests replace deterministic regression tests.

## Scenario analysis

### T0: Immediately after the Go clean run

Alice wants the next step to answer what the Go simulator did not answer. Bob
wants the analyzer to stay useful. Carol wants firmware progress. Dave wants
deployment claims to be honest. Ellen wants the roadmap to say why each lane
exists.

Alt A is the best first answer because it asks whether the node can be shaped
for M4 constraints without waiting for Renode. Alt D is also strong if the team
wants Renode immediately, because it starts with M4 platform viability rather
than radio or full protocol behavior. Alt C is useful only as a support task.
Alt E is necessary but should follow proof that the node and platform shape are
credible. Alt F is too large. Alt G is too late-stage.

### T1: M4 memory and runtime fit

Mallory looks for hidden heap growth, unbounded queues, large CBOR buffers,
host-only storage, and panic paths that cannot be recovered on a small device.
Carol asks whether the node core can plausibly become firmware.

Alt A targets this directly. Alt D targets this through Renode. Alt B can help
only if it enforces embedded boundaries. Alt C proves too little by itself. Alt
E and Alt F depend on this work. Alt G can validate later but should not be the
first place these problems are discovered.

### T2: PromiseGrid identity of the M4 node

Frank checks whether the M4-shaped node is still a real PromiseGrid node: binary
pCID CID in slot 0, exact CBOR artifacts, local pCID judgment, order-status
payloads, peer-storage capability tokens, and no central authority.

Alt C is useful here as a narrow fixture-parity support step. Alt A should
include only enough of Alt C to prove the node speaks the current POC17 shapes.
Alt B can prove more behavior but risks becoming a Linux simulator. Alt D and
Alt E should not parse app payloads in the platform or radio model.

### T3: LoRa radio path

Bob wants evidence that packet bytes cross a radio-shaped boundary. Carol wants
SPI/FIFO/IRQ behavior. Mallory tries to sneak in a host bridge.

Alt E is strongest, but only after the M4 node or platform shape is credible.
Alt D prepares for it. Alt A gives the firmware core something real to send
once the seam exists. Alt C does not prove this. Alt F can prove it later, but
combines too many risks if attempted first.

### T4: Full POC17 behavior under stronger fidelity

Frank wants peer storage and order status outside the Go simulator. Mallory
wants stale-token, restart, missing-parent, malformed-frame, and MTU cases in
the stronger lane.

Alt F is the milestone for this, but it should follow Alt A or Alt D plus Alt E.
Trying Alt F first would mix M4 fit, protocol porting, platform bring-up, radio
modeling, and analyzer work into one hard-to-debug step.

### T5: Production-facing claims

Dave asks whether POC17 can support a claim that a PromiseGrid node can run on
M4/LoRa hardware. Ellen checks whether old documents overclaim.

The answer needs a chain: Go behavior oracle, M4-shaped Rust core or Renode M4
platform, radio seam, full Renode behavior, then hardware-in-loop validation.
Any single lane alone is not enough. The roadmap must say which claim each lane
supports.

## Surviving and rejected alternatives

Alt A survives as the recommended next software step because it prioritizes M4
viability while keeping iteration fast.

Alt D survives as the recommended first Renode step because M4 platform
viability should precede radio seam and full protocol behavior.

Alt E survives as the required next fidelity step after M4 viability, not as the
first roadmap priority.

Alt C survives as a support task that should be scoped to the current POC17
message shapes.

Alt B survives only if it is constrained by Alt A's embedded boundaries. A
generic Rust-on-Linux simulator is rejected.

Alt F survives as a later integration milestone, not as the next step.

Alt G survives as late validation before production-facing claims.

Rejected: protocol parity as the main roadmap driver. It is necessary, but not
the purpose of POC17.

Rejected: radio seam as higher priority than M4 compatibility. It is required,
but a LoRa seam does not answer whether the node can fit and run on M4-class
hardware.

Rejected: claiming hardware readiness from Go, Rust-on-Linux, or Renode alone.
The claim needs measured and validated evidence.

## Conclusions

The recommended roadmap is M4-first:

1. Keep the Go simulator as the behavior oracle and fast regression lane.
2. Build an M4-shaped Rust node core with bounded buffers, explicit limits, and
   no hidden host dependency.
3. Add only enough protocol fixture parity to prove the M4-shaped node speaks
   the current POC17 pCID/CBOR shapes.
4. Use Linux only as a harness around the constrained Rust core.
5. Bring the Rust node core into a minimal Renode M4 platform.
6. Add an RFM95/SX127x-shaped SPI/FIFO/IRQ radio seam.
7. Run the full POC17 order-status, peer-storage, restart, and failure scenario
   in the Rust/Renode lane.
8. Use hardware-in-loop later to validate production-facing assumptions.

The DF result is locked by `DI-pokin`: choose Alt D, Renode M4 platform first,
as the next implementation lane. Alt D best matches POC17's main purpose
because it proves firmware boot, reset/panic behavior, memory layout,
diagnostic extraction, and deterministic M4-shaped execution before radio seam
or full protocol behavior. Protocol parity remains scoped to what the M4 node
needs, and the LoRa seam follows M4 viability rather than outranking it.

## Implications for the repo's open TODOs and pending DIs

- `komon.25` now has an M4-first TE with the first post-Go fidelity lane locked
  by `DI-pokin`: Renode M4 platform first.
- Implementation should start with Renode M4 platform viability before the
  RFM95/SX127x radio seam or full POC17 scenario port.
- `komon.12` should cite this TE when documenting simulator fidelity limits and
  should say clearly that POC17's main unproven claim is M4/LoRa hardware fit.
- `komon.13` should keep archiving Go clean-run artifacts because the Go lane
  remains the behavior oracle while M4/Renode work begins.
- `DEV-GUIDE-RESOURCES.md` must not cite POC17 as hardware-ready until M4,
  radio, resource, packet, and hardware validation gates exist and pass.
