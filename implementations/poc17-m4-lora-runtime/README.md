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

The clean run writes artifacts under `/tmp/wire-lab-poc17/poc17-demo/`:

- `events.jsonl` records simulator and analyzer events.
- `message-cas/<cid>.cbor` keeps exact emitted CBOR envelope bytes.
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
- Ivan and Bob exchange bintags-shaped order status messages with order
  number, status, source, destination, counter, and MSG/ACK flow under
  `order_status_v1`. The workflow comes from bintags, while the wire shape
  stays PromiseGrid CBOR. Source: `DI-mokit`.
- The radio model covers MTU refusal, loss, duplicate, replay, delay, malformed
  bytes, and asymmetric reachability as transport effects.
- The default simulated application MTU is 200 bytes, leaving margin under the
  238-byte application-buffer ceiling discussed in `DN-zaraz`. Source:
  `DI-dutah`.
- The device keeps a tiny sparse CAS, records missing parents, asks for
  peer-storage help, and performs local GC under pressure.
