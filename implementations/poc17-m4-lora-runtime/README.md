# poc17-m4-lora-runtime

`poc17-m4-lora-runtime` is the first Go behavior simulator for a
Feather/RFM95-shaped constrained PromiseGrid agent. It starts with the sequence
locked by `DI-libis`: prove the radio-only PromiseGrid behavior in Go first,
then leave Rust/Renode for a later fidelity lane.

This POC is behavior evidence only. It does not claim exact Feather M4 Express,
SAMD51, SPI, RFM95W/RFM95, CircuitPython runtime, radio-driver, packet, memory,
energy, regulatory, or production-device fidelity. Source: `DI-libis`.

## Run

```sh
go test ./...
./scripts/run-clean.sh
```

The first Renode M4 platform smoke slice uses containers for both the Rust
cross-build and Renode execution:

```sh
./scripts/run-renode-m4-smoke.sh
```

This builds a tiny Rust `no_std` firmware ELF for `thumbv7em-none-eabihf` and
loads it on Renode's ATSAMD51 Cortex-M4F platform. It proves only that the
containerized build path, platform load, reset vector, start, pause, and
diagnostic output work. It does not prove RFM95/SX127x, SPI, packet timing,
energy, regulatory, CircuitPython, or production hardware readiness. Source:
`DI-pokin`.

The next Renode risk to retire is the radio seam. The next slice should prove
that firmware-visible packet bytes can move through an RFM95/SX127x-shaped
SPI/FIFO/IRQ path without using UART, host callbacks, shared files, or any other
hidden PromiseGrid message bridge. Full order-status, peer-storage, restart,
and failure-scenario parity should wait until that radio path is credible.
Source: `DI-togag`.

The clean run writes artifacts under `/tmp/wire-lab-poc17/poc17-demo/`:

- `events.jsonl` records simulator and analyzer events.
- `message-cas/<cid>.cbor` keeps exact emitted CBOR envelope bytes.
- `lifecycle-cas/<cid>.cbor` keeps exact host-local lifecycle frames.
- `malformed/<label>-<cid>.bin` keeps malformed radio bytes for review.

## Current Slice

- A Go simulator models one M4-shaped device and one non-M4 peer.
- PromiseGrid messages cross the simulated radio path only.
- The device parses `grid([42(pCID), payload])` and
  `grid([42(pCID), payload, proof])`.
- Slot 0 carries actual CIDv1 raw sha2-256 bytes for the embedded RFC-like
  protocol document, not a readable placeholder name. Source: `DI-dutah`.
- Logs, filenames, parent links, CAS IDs, and diagnostic pCIDs use canonical
  CIDv1 base32 text with the multibase `b` prefix. Source: `DI-nopiv`.
- Host-local simulator lifecycle/resource promises use POC16-shaped
  `local_lifecycle_v1` CWT/COSE tokens. These token bytes do not cross the
  simulated LoRa path in this slice; a later DR/DI is required before making
  lifecycle tokens radio-visible. Source: `DI-zopub`.
- Runtime limits are config-driven and analyzer-visible: RAM bytes, flash bytes,
  energy units, radio airtime bytes, retry count, and local CAS object count.
  The current run reports configured limits only; it does not report activity
  usage unless the simulator actually measures that usage. Source: `DI-gidul`;
  `DI-rujod`.
- Ivan and Bob exchange bintags-shaped order status messages with order
  number, status, source, destination, counter, and MSG/ACK flow under
  `order_status_v1`. The workflow comes from bintags, while the wire shape
  stays PromiseGrid CBOR. Source: `DI-mokit`.
- The radio model covers MTU refusal, loss, duplicate, replay, delay, malformed
  bytes, and asymmetric reachability as transport effects.
- The default simulated application MTU is 200 bytes, leaving margin under the
  238-byte application-buffer ceiling discussed in `DN-zaraz`. Source:
  `DI-dutah`.
- The device keeps a tiny sparse CAS, records missing parents, and uses
  radio-visible `peer_storage` promises. Bob first grants Ivan a compact bearer
  capability; Ivan presents it on `put` and `get`; Bob retains or returns exact
  bytes; Ivan verifies returned bytes against the requested CID before accepting
  them. Source: `DI-gidul`.
- The run includes a fresh-agent restart: Ivan reloads durable identity, order,
  token, and missing-CID state into a new agent process, then recovers the
  missing parent from Bob rather than assuming volatile CAS bytes survived.
  Source: `DI-gidul`.
