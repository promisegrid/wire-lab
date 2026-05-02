# TODO 019 — ns-3 harness scaffold for UDP-binding v0

Source: TE-29 OQ-29.9 (`docs/thought-experiments/TE-20260501-215027-protocols-as-simulated-repos-and-binding-layer.md`),
revised 2026-05-01 with empirical sandbox data showing ns-3 3.43 is
the lowest-friction network-realism option in the Perplexity Computer
sandbox (Debian trixie, kernel 6.1.158). `tc netem` is unavailable
(kernel lacks `sch_netem`), so packet-loss / latency / jitter
adversarial scenarios must live in ns-3.

Sibling of TODO 018. TODO 018 is gated on this scaffold proving
end-to-end round trip at least once.

## Goal

A minimal ns-3 harness that:

1. Stands up a 2-node UDP topology in ns-3.
2. Lets the Go UDP-binding v0 reference implementation (TODO 018)
   send and receive datagrams over that emulated wire.
3. Captures PCAP for offline analysis.
4. Produces simulation-artifact files at the wire-lab's canonical
   path:

   ```
   transports/udp/<udp-binding-pCID>/<session-pCID>/<message-pCID>/<message-id>.msg
   ```

5. Becomes the regression baseline for every later binding/session
   change.

The scaffold's first scenario proves the integration works at all.
Loss/latency/jitter scenarios are explicitly out of scope for v0
of the scaffold; they are TODO 020+ once the integration is solid.

## Scaffold layout

Per TE-32, this scaffold is a B-side artifact and lives under the
new top-level `implementations/` directory rather than under
`tools/`. The wire-lab harness has both an A-side (design tree at
`protocols/wire-lab.d/`) and a B-side (reference implementation
including the ns-3 fixture). This TODO targets the B-side. The
implementation tree carries its own `CHANGELOG.md` recording
conformance claims against upstream spec doc-CIDs (the harness
spec, the UDP-binding spec, and any session/message specs the
fixture exercises).

```
implementations/ns3-harness-fixture/
├── README.md                     setup, build, run, integration notes
├── CMakeLists.txt                pkg-config-driven build against Debian ns-3
├── scenarios/
│   └── 001-udp-roundtrip/
│       ├── scenario.cc           ns-3 topology + tap-bridge wiring
│       ├── run.sh                build + launch driver
│       └── expected/             reference PCAPs and .msg files for diff
├── pcap-to-msg/                  Go tool: post-process PCAP into .msg files
│   ├── main.go
│   └── README.md
└── env.sh                        sourceable env: pkg-config paths, etc.
```

## Subtasks

1. **`tools/ns3-harness/CMakeLists.txt`**: pkg-config build against
   Debian ns-3 3.43 (`ns3-core`, `ns3-network`, `ns3-internet`,
   `ns3-point-to-point`, `ns3-applications`, `ns3-tap-bridge`).
   Pattern follows the smoke-test recipe already proven in the
   sandbox session.

2. **Scenario 001 — UDP round trip.** Two ns-3 nodes on a 10.1.1.0/24
   point-to-point link. Each node has a `TapBridge` exposing a
   Linux tap interface (`ns3-tap0`, `ns3-tap1`) that the Go
   UDP-binding implementation binds to. PCAP enabled on both
   interfaces.

3. **Driver `run.sh`** that:
   - Builds `scenario.cc` against Debian ns-3.
   - Brings up tap interfaces (sudo required).
   - Launches two instances of the Go binding (one sender, one
     receiver) bound to the tap interfaces.
   - Sends N test messages from sender to receiver using TODO 018's
     reference implementation.
   - Lets ns-3 simulate for enough wall-clock time to drain the
     queue.
   - Captures PCAP for both sides.
   - Tears tap interfaces down on exit (trap handler).

4. **`pcap-to-msg/`**: a Go post-processor that reads a PCAP, filters
   to UDP datagrams matching the binding's ports, and writes one
   `.msg` file per datagram to:

   ```
   transports/udp/<udp-binding-pCID>/<session-pCID>/<message-pCID>/<message-id>.msg
   ```

   File contents = exact UDP payload bytes. The session-pCID and
   message-pCID are passed in as flags (the scaffold does not parse
   session-layer headers; TODO 021 may change that once group-session
   is implemented). The `<message-id>` derives from a hash of the
   payload bytes (consistent with TE-29 OQ-29.2 lean toward content
   hash).

5. **Reference outputs.** Commit a deterministic scenario seed and
   reference PCAPs / `.msg` files under
   `scenarios/001-udp-roundtrip/expected/`. The scenario passes when
   re-running produces byte-identical output. Determinism is critical
   for C-2 (multi-generational durability) and for catching
   regressions.

6. **README.** Cover (a) sandbox prerequisites (`apt-get install ns3
   libns3.43 libns3-dev cmake ninja-build pkg-config tcpdump`),
   (b) sudo requirement for tap-bridge, (c) what scenario 001
   proves, (d) where output files land, (e) how to add a new
   scenario.

7. **Integration with TODO 018.** TODO 018's done-criteria reference
   this scaffold. Once both TODOs are mergeable, land them in the
   same merge twig (TE-29 was a single coherent locking; the
   implementation that exercises TE-29 should land coherently too).

## Out of scope for TODO 019

- Loss / latency / jitter / reordering scenarios. Tracked as TODO
  020 once the v0 scaffold is proven.
- Multi-node (>2) scenarios. Tracked separately; will exercise
  group-session v0 once that exists.
- Mininet path. Empirically Mininet works in the sandbox but is
  redundant with ns-3 for this first scaffold. Possible later if a
  topology-only experiment needs it.
- Direct Code Execution (DCE). The tap-bridge approach is simpler
  and lets the Go binding stay an ordinary user-space binary. DCE
  is a possible later optimization.

## Dependencies

- **TODO 018** for the Go UDP-binding v0 implementation. The
  scaffold cannot prove round trip without something to bind to the
  tap interface. TODOs 018 and 019 are mutually gating: 018's done
  criterion references 019's scenario; 019's scenario needs 018's
  binary. They land together.
- ns-3 Debian package and tap-bridge support must work in whatever
  environment runs the harness. Confirmed working in the Perplexity
  Computer sandbox (TE-29 OQ-29.9 empirics). Should also work on
  any modern Debian/Ubuntu dev box.

## Done when

- `tools/ns3-harness/scenarios/001-udp-roundtrip/run.sh` exits 0.
- The PCAP captured by ns-3 contains the expected number of UDP
  datagrams in the expected direction.
- The post-processor produces the expected `.msg` files at the
  canonical `transports/udp/.../*.msg` path.
- Re-running with the same seed produces byte-identical output.
- The README is complete enough that a fresh contributor can
  reproduce the scenario in under 30 minutes on a fresh Debian
  install.

## Future scenarios to track separately (preview, not commitments)

- 002: UDP round trip with 5 percent symmetric packet loss.
- 003: UDP round trip with 50 ms one-way latency, 10 ms jitter.
- 004: 1-of-N receivers in a multicast group.
- 005: Network partition (link drops to 0 percent for window W).
- 006: Asymmetric loss (loss only in one direction).
- 007: TCP-binding round trip (gates on TODO 020 second-binding spec).
- 008: group-session v0 over UDP-binding v0 with one Byzantine node
  (gates on group-session v0 reference implementation, TODO 021).
