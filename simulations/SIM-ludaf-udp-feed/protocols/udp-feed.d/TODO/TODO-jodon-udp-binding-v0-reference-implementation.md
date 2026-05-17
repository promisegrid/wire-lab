# TODO-jodon — UDP-feed v0 reference implementation

## Prior aliases

Before the TE-39 proquint migration, this file was known as:

- `TODO-18` (integer alias)
- `TODO-20260501-224805` (timestamp alias and pre-migration filename)

Source: TE-29 (`docs/thought-experiments/TE-20260501-215027-protocols-as-simulated-repos-and-binding-layer.md`).
Spec: `protocols/udp-feed.d/specs/udp-binding-draft.md` (one-page
draft, 10 normative promises, anti-promises, and test-vector
placeholders).

Simulation note: this TODO now lives with the `udp-feed` protocol specimen at
`simulations/SIM-ludaf-udp-feed/protocols/udp-feed.d/` (moved in `rusis.6`,
`DI-loluk`). Historical `udp-binding` / `SIM-piloh` paths remain provenance
only unless a line explicitly names a historical reference.

Flesh out UDP-feed v0 from one-page sketch into a usable v0
artifact: Go reference implementation, conformance test vectors, and
a minimal ns-3 harness scaffold that proves end-to-end round trip.
The protocol/specimen question about what evidence should count as a usable v0
conformance surface is now captured in
`simulations/SIM-kuful-udp-feed-v0-conformance/`, without implementing or
closing this TODO. Source: `DI-pukap`.

This is the first concrete binding implementation under the TE-29
layer decomposition, so it doubles as the proving ground for the
binding-layer abstraction itself: anything painful here is an
abstraction defect, not just a UDP issue.

## Subtasks

1. **Reference implementation** at
   `implementations/go-udp-feed-reference/` in Go (per Steve's
   standing language preference, and per TE-liviv which locks B-side
   implementation artifacts under `implementations/<impl-name>/`
   rather than under any `protocols/<slug>.d/` subtree). The
   implementation tree carries its own `CHANGELOG.md` recording
   conformance claims (e.g. `claim: implements` with the upstream
   `udp-binding-draft.md` doc-CID once frozen).

   - `Send(msg []byte, addr Addr) error` honoring promises 1, 2, 3,
     and 7 of the spec (one datagram per message, 1232 max, no
     in-band reliability, DSCP 0).
   - Recv loop with caller-supplied `Handle(msg []byte, src Addr)`
     callback honoring promise 6 (peer-set filter optional).
   - Local-error path for promise 2 size violations.
   - Promise 8: do not disable UDP checksum (`SO_NO_CHECK = 0`).
   - Stateless per promise 9.

2. **Simulation-artifact writer** honoring promise 10. When invoked
   under the wire-lab simulator, write each datagram payload to:

   ```
   transports/udp/<this-binding-pCID>/<session-pCID>/<message-pCID>/<message-id>.txt
   ```

   The session-pCID, message-pCID, and message-id come from the
   session layer above; this binding does not parse them. Behind a
   simulator-mode flag so production sends do not write artifacts.

3. **Test vectors.** Author the placeholder TVs from the spec
   (TV-1 through TV-5):
   - TV-1: 612-byte session message round-trips byte-for-byte over
     loopback UDP/4646.
   - TV-2: 1232-byte message round-trips byte-for-byte.
   - TV-3: 1233-byte message produces a local sender error before
     any datagram leaves the host.
   - TV-4: malformed datagram is handed up to the session layer
     unmodified by the binding.
   - TV-5: simulation-artifact file written for TV-1 contains exactly
     the 612 bytes of TV-1's session message.

   Test vectors live at `tools/udp-feed/testvectors/`.

4. **`/tmp/spec` walks `protocols/`.** Currently the spec checker
   only walks `specs/`. Update it (or write its successor) to walk
   `protocols/<slug>.d/specs/` so this draft is recognized.

5. **TODO-bihon — ns-3 harness scaffold for UDP-feed v0.** A
   minimal 2-node ns-3 scenario that proves round-trip works through
   the Go reference implementation talking over an ns-3-emulated
   UDP wire. See `protocols/wire-lab.d/TODO/TODO-bihon-ns3-harness-scaffold.md` for full
   subtasks. Tracked as a sibling TODO rather than a subtask of this
   one because the scaffold has its own follow-on lifecycle (loss
   scenarios, multi-binding scenarios) that long outlives the v0
   binding implementation. TODO-jodon is done when the scaffold from
   TODO-bihon successfully runs the v0 reference implementation
   end-to-end at least once.

6. **Update `protocols/udp-feed.d/specs/udp-binding-draft.md`**
   to reference (a) the implementation path under `tools/`, (b) the
   test-vector files, (c) the ns-3 harness scenario name. Replace
   "to be added in TODO-jodon" placeholders with concrete pointers.

## Out of scope for TODO-jodon

- Multiple concurrent UDP bindings on different ports (use case:
  one binding per protocol stack on a host). Possible v1 extension.
- IPv6 specifics beyond the 1232-byte size derivation. Should Just
  Work but is not exercised by v0 test vectors.
- NAT traversal. Explicit anti-promise in the spec.
- Path MTU discovery beyond "error out below 1232." Implementations
  may add it; not required by v0.
- Multicast. Permitted by the spec but not exercised by v0 test
  vectors; receivers can join groups manually for ad-hoc testing.

## Out of scope, but flagged as next likely TODO

- A second feed/binding (TCP feed v0 or WebSocket feed v0) so the
  binding-layer abstraction is exercised across at least two
  qualitatively different real-world transports. This is when the
  per-binding-pCID forking property at C-4 actually starts paying
  off. (Originally hedged as "Likely TODO 020" before the TODO 020 slot was claimed by an unrelated TE-editing-policy item; this is a future TODO not yet filed.)
- A Go reference implementation of group-session v0 to ride above
  UDP-feed v0 and prove the layer composition. (Originally hedged as "Likely TODO 021" before that slot was claimed by session-replay-cleanup; this is a future TODO not yet filed.)

## Dependencies

- TODO-vuhuj (protocols-as-simulated-repos migration) ideally lands
  first so the spec lives at its final path. If TODO-vuhuj has not
  landed, this TODO uses `protocols/udp-feed.d/specs/udp-binding-draft.md`
  (already created in TE-vipir's commit).
- TODO-losoh (group-session rename) is not required; UDP-feed does
  not depend on the session protocol's slug.

## Done when

- `tools/udp-feed/` Go package builds and passes all five test
  vectors.
- `protocols/udp-feed.d/specs/udp-binding-draft.md` is updated
  with concrete pointers replacing "to be added" placeholders.
- TODO-bihon's ns-3 harness scenario runs the Go reference
  implementation end-to-end at least once and produces matching
  PCAP and `.txt` artifacts.
- `/tmp/spec check` (or its successor) reports OK with the spec
  recognized at its new path under `protocols/udp-feed.d/`.
