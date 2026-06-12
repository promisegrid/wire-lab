# poc14-wasm

`poc14-wasm` is executable POC evidence for PromiseGrid storage, compute,
kernel, shipping, live-agent, WASM-boundary, stdio-boundary, and decentralized
monitoring behavior. It is intentionally a superset of POC13: it keeps POC11's
autonomous sparse-mesh relationship/economics pressure, POC12's separate
app/kernel processes and shipping/device workflow, and POC13's CAS, compute,
replica recovery, token lifecycle, verifier disagreement, run-scoped durability,
retention/GC, backpressure, rate-limit, replay protection, bounded trust, and
dynamic topology gates. POC14 adds Peggy as a WASM-boundary process, Victor as a
stdio-worker process behind a local adapter, and analyzer gates for
decentralized monitoring signals that do not assume a production-wide observer.
Source: `DI-sihuz`; `DI-sifot`; `DI-fimoh`; `DI-lulof`; `DI-linof`.

## What This Tests

- Multiple app pCIDs through one local container kernel: `relationship_v1`,
  `postal_scale_v1`, `ups_label_v1`, `printer_port_v1`, `accounting_v1`,
  `cas_storage_v1`, `cid_compute_v1`, and `identity_key_v1`.
- Real app/kernel process boundary: each container runs one `poc14-kernel`
  process plus separate local app processes for relationship, fulfillment,
  postal scale, UPS label printer, printer port, and accounting roles.
- Kernel-style pCID routing: the kernel parses slot 0 `42(pCID)`, checks local
  app receive promises, and delivers exact bytes to the app process that
  promised the target pCID.
- Deterministic device/system agents for scale, printer-port hardware access,
  label-printer business logic, and accounting alongside live LLM
  business/social agents.
- Deterministic heterogeneous-boundary agents: Peggy validates a WASM module
  fixture in her own app process, and Victor's worker process sends and receives
  exact PromiseGrid envelopes only through stdin/stdout with a local adapter.
  Source: `DI-linof`.
- One top-level semantic act, `promise`; workflow steps are payload meanings,
  not RPC verbs.
- Explicit direct TCP relationship transition evidence:
  `direct_peer_added`, `direct_peer_removed`, and `direct_peer_unchanged`.
- Corrected promise evidence semantics: ordinary `not_promised` /
  `non_commitment` evidence is neutral rather than broken peer-promise
  evidence, provider/runtime decision failures stay local, and duplicate
  shipment updates are checkpointed without repeated trust inflation.
  Source: `DI-jinoz`.
- App-local promise journal: POC14 now records an outstanding promise before
  applying kept/broken/malformed trust evidence, separates `send_unavailable`
  from `send_not_promised`, suppresses repeated live-agent promises, and keeps
  local budget/capacity exhaustion out of peer trust. Source: `DI-vujob`.
- Sender-side non-commitment restraint: receiver `not_promised` evidence is
  remembered in an app-local `nonCommitmentJournal`, and later same-run semantic
  retries are recorded as `promise_not_promised_suppressed` without changing
  peer trust. Generic app-local checkpoints now use `checkpointJournal` rather
  than a shipment-only map. Source: `DI-zapab`.
- Local hardware promise tokens: the UPS label printer must ask the local
  `printer_port` resource owner for a scoped future-print promise token, then
  redeem that token with bounded label bytes before it can return print evidence
  to fulfillment. Source: `DI-pohaj`; `DI-vutok`.
- CAS and compute protocol promises: Alice stores exact content bytes with Bob,
  Bob replicates to Frank, Alice retrieves by primary and replica promise-token
  evidence, Alice asks Carol for CID-named function execution, Dave caches and
  verifies the result, Grace supplies disagreement pressure, and Mallory sends
  corrupt bytes, an unknown pCID, an unsupported variant, a bad proof, a key
  rotation promise, and a capacity probe. Source: `DI-sinur`.
- Run-scoped durability promises: apps persist CAS objects, compute cache
  checkpoints, capability tokens, replay windows, and local evidence journals
  under the current run root so an app can recover inside one run, while
  `scripts/run-clean.sh` remains the experiment boundary that resets state.
  Source: `DI-sunuf`.
- Retention, GC, backpressure, rate-limit, and replay evidence: apps promise
  local retain-until/delete-after/token-expiry/disk-pressure behavior, record
  retained/removed/ended/broken GC cases, model sender and receiver rate limits
  as reciprocal self-promises, and reject exact envelope or serve-once token
  replays as local non-commitment evidence. Source: `DI-sunuf`.
- POC14 hardening evidence: saved evidence summaries count all local
  non-commitment outcomes, non-mutating ACKs distinguish true duplicates from
  refusals/cache misses/replay refusals/future-only repair, trust scores stay on
  a bounded local scale, recovery caution is analyzer-visible, and dynamic TCP
  topology changes affect real app/kernel send reachability. Source:
  `DI-sihuz`.
- Superset analyzer gates: `poc14-analyze` now fails if inherited POC11/POC12
  behavior, POC14 storage/compute evidence, run-scoped durability, retention/GC,
  pressure, rate-limit, or replay evidence disappears. Source: `DI-sinur`;
  `DI-sunuf`; `DI-sihuz`.
- Observer-only monitor lifecycle: completed nodes wait for `monitor.done` using
  a config-derived provider/turn/grace budget, and a completed run can still
  write the non-authoritative marker if the observer report fails. Source:
  `DI-jupob`.
- Decentralized monitoring evidence: POC14 records local evidence summaries,
  peer-carried attestations, bearer-token exchange-rate signals, relationship
  topology signals, and voluntary gossip as ordinary local promises rather than
  global monitor facts. Source: `DI-lulof`.
- Local hard trust boundaries: Alice records a permanent local distrust decision
  about Mallory and a local promise that Alice's inbound/outbound traffic should
  not transit Mallory. The ledger now blocks Alice's future Mallory sends and
  rejects route candidates with Mallory as a transit hop. This is local evidence
  and local route choice, not a network-wide ban or authorization policy. Source:
  `DI-kinaf`; `DI-dubih`.
- Mixed-version and restart evidence: POC14 records local pCID migration
  promises and same-run restart/recovery promises so future work can test
  protocol evolution and process crashes without relying on cross-run state.
  Source: `DI-linof`.

## Protocol Shape

Every message is a signed CBOR envelope:

```text
grid([42(pCID), payload, proof])
```

The pCID identifies the protocol spec. Message variants such as
`weigh_package`, `address_lookup`, `print_label`, `issue_print_capability`,
`redeem_print_capability`, `shipment_update`, `store_content`,
`serve_replica_content`, `execute_function`, `lookup_compute_cache`, and
`verify_compute_result` are payload meanings inside their protocol, not
separate pCIDs. Source: `DI-bikit`; `DI-pohaj`; `DI-sinur`.

Payload shape is owned by the pCID. Most POC14 payloads still use the older
`field_*` map scaffold so the existing agents can interoperate while the POC
evolves, but that map is not a PromiseGrid-wide payload standard.
`identity_key_v1` is the first narrow cleanup example: its key-rotation request
and ACK payloads are CBOR arrays defined by that pCID, then decoded into
compatibility fields only at the local runtime boundary. Source: `DI-vipih`.

## Shipping Agents

- `poc14-fulfillment`: hybrid workflow coordinator that executes one deterministic
  startup shipment sequence across the production pCIDs, then continues normal
  live LLM relationship turns. It sends shipping pCIDs but receives only the
  relationship pCID in this POC.
- `poc14-postal-scale`: deterministic app for `postal_scale_v1`; promises
  package weight evidence only.
- `poc14-ups-label-printer`: deterministic app for `ups_label_v1`; promises
  label, cost, and tracking evidence only after it receives and redeems a local
  printer-port capability-promise token.
- `poc14-printer-port`: deterministic app for `printer_port_v1`; represents
  the local hardware-access kernel role for a simulated USB printer port. It
  promises bounded future printing by issuing a scoped token to the label
  printer, then promises print evidence when that token is redeemed with label
  bytes.
- `poc14-accounting`: deterministic app for `accounting_v1`; promises address
  lookup and shipment update evidence only.
- `poc14-relationship-agent`: generic live LLM relationship app. Depending on
  the configured agent, it can also promise `cas_storage_v1` or
  `cid_compute_v1` handling while keeping the same single top-level `promise`
  action. `identity_key_v1` is reserved for scripted key-rotation array payloads
  in this cleanup slice. Source: `DI-vipih`.
- `poc14-wasm-agent`: deterministic Peggy app process that validates WASM module
  bytes and sends WASM-boundary evidence as a normal `relationship_v1` promise.
- `poc14-stdio-adapter`: deterministic Victor adapter process that starts
  `poc14-stdio-worker`, receives exact envelope bytes over stdout, forwards
  those bytes through the local kernel, and returns the exact peer ACK over
  stdin.
- `poc14-stdio-worker`: subprocess agent whose application messaging path is
  stdin/stdout only; it signs one PromiseGrid envelope and locally verifies the
  returned ACK envelope.
- `poc14-kernel`: container-local transport process that records operational
  routing evidence only; it does not own trust, workflow, device behavior, or
  promise judgment. Source: `DI-galin`.

## Run

Copy the committed template and keep secrets out of it. `compose.yaml` mounts
the host `OPENAI_API_KEY` value as `/run/secrets/openai_api_key`, so the key does
not need a repo-local secret file and does not need to appear on the command
line. Source: `DI-pohoh`.

```sh
cp config.example.json config.json
scripts/run-clean.sh
```

For manual runs, `poc14-analyze` accepts either the parent run directory or its
`run/` JSONL directory and fails loudly if no JSONL evidence exists:

```sh
docker compose up --build --abort-on-container-exit
docker compose run --rm --entrypoint /usr/local/bin/poc14-analyze dave /run/poc14/<run_id>
```

The expected run narrative is in `docs/RUN-NARRATIVE.md`; current POC-to-
production fitness notes are in `docs/PRODUCTION-FITNESS.md`; implementation
notes are in `docs/IMPLEMENTATION-NOTES.md`.

`config.json`, `poc14.env`, provider outputs, and Docker volume state are
ignored and must not be committed. `openai_api_key.txt` remains ignored for
older local workflows but is not required by the current Compose setup.

## Current Limits

POC14 is provisional executable evidence, not a stable shipping, device,
accounting, WASM, stdio, kernel-routing, monitor, trust, provider, or workflow
API. Docker networking remains static; dynamic TCP relationships are local
promises to dial or accept direct exchanges between app agents, carried by local
kernels. The kernel records only app receive-promise registration, exact-byte
delivery, peer-forwarding, unregistered pCID, and transport outcomes. Apps own
trust, keep/break/non-commitment judgment, relationship ledgers, workflow state,
and deterministic device/system/boundary behavior. The fulfillment startup
sequence is a POC guardrail so the run produces concrete pCID-routed shipment
evidence instead of relying on a live LLM to choose that sequence unaided.
Production monitoring cannot rely on POC14's whole-run analyzer because real
agents are distributed across legal entities; POC14's decentralized-monitoring
events are candidate local evidence signals, not a global dashboard. Source:
`DI-timah`; `DI-bikit`; `DI-parok`; `DI-galin`; `DI-pohaj`; `DI-sinur`;
`DI-lulof`; `DI-linof`.

The repaired POC14 should be treated as the current superset baseline for future
POCs unless a later scoped DI explicitly declares a non-superset exception and
lists the features being dropped. Source: `DI-sinur`.

Post-split evidence semantics are intentionally conservative:
`not_promised` means the receiver did not promise the requested exchange, not
that the receiver broke a promise. Transient provider/runtime decision failures
are local app/runtime evidence, not peer trust evidence. The accounting app
keeps an app-local shipment checkpoint keyed by order, tracking number, and
cost; duplicate confirmations remain visible in logs but do not repeatedly
increase trust. Source: `DI-jinoz`.

The 2026-06-05 `poc14-jinoz-20260605-055916` validation run exited cleanly and
the analyzer reported 454 events: 433 kept, 20 non-commitment, and 1 broken
resource promise. Shipping evidence appeared exactly once for address lookup,
package weight, label printing, and accounting update/confirmation. The monitor
still flagged a real remaining concern: some trust changes are driven by local
budget/capacity exhaustion and should be separated from peer trust in a later
POC. Source: `DI-jinoz`.

`DI-vujob` keeps those follow-up corrections in POC14. The runtime now records
`promise_outstanding` before resolving local promise evidence, records
`promise_resolved` before applying trust effects, treats local send failures as
`send_unavailable`, treats receiver non-commitment as `send_not_promised`,
records `local_resource_exhausted` without changing peer trust, and suppresses
repeated live-agent promise text for the same target/protocol. The startup
shipping workflow intentionally repeats the accounting update once so Docker
runs exercise duplicate shipment checkpointing through the real app/kernel path.
The analyzer reports `local_resource_counts` and
`resource_trust_coupling_counts` so future runs can show whether local
budget/capacity state is leaking back into peer-trust transitions. The
2026-06-05 `poc14-vujob-20260605-063213` Docker run exited cleanly and analyzed
to 598 events: 569 kept, 29 non-commitment, 86 `promise_outstanding`, 86
`promise_resolved`, 7 `local_resource_exhausted`, 6 `send_not_promised`, 4
`promise_repeated_suppressed`, 1 `accounting_update_duplicate`, 1
`accounting_update_duplicate_confirmed`, and empty
`resource_trust_coupling_counts`. Source: `DI-vujob`.

`DI-zapab` tightens the same evidence model. The runtime now remembers receiver
`not_promised` outcomes by target, pCID name, and `field_promise_about`, then
suppresses a later live-agent retry for the same semantic promise as
`promise_not_promised_suppressed`. That suppression is Alice's local restraint,
not a penalty against Bob. Duplicate evidence now flows through a reusable
`checkpointJournal`; the accounting shipment update remains the first concrete
checkpoint. Source: `DI-zapab`.

`DI-vahan` makes the alternate compute path follow local trust evidence. Alice
still obtains arbitrary payload-defined compute coverage, but after Carol exposes
malformed bad-result evidence, Alice sends the second sum-function promise to
Dave rather than forcing another fresh compute promise to Carol.

`DI-pohaj` adds a local printer-port kernel-role app without turning the
message kernel into a USB authority or RPC service. During label printing,
`ups_label_printer` sends `issue_print_capability` to `printer_port`;
`printer_port` returns a deterministic token bound to the issuee, scope, token
ID, and byte limit; `ups_label_printer` sends `redeem_print_capability` with the
token and exact hex label bytes; `printer_port` returns deterministic local
spool evidence. Unit tests cover the token issue/redemption path, wrong-token
rejection, routing pCID registration, and analyzer shipping-event recognition.
The fresh Docker run `poc14-jupob-20260606-030610` exited cleanly after Dave
wrote `monitor_done`; analyzer output reported 625 events, 17 `node_done`, 17
`shutdown_grace_elapsed`, one `monitor_done`, all printer-port shipping counts,
and empty `resource_trust_coupling_counts`. Source: `DI-pohaj`; `DI-vutok`;
`DI-jupob`.

`DI-sunuf` adds run-scoped durability and operational-pressure evidence without
promoting POC14 state into cross-run infrastructure. Each app writes
`stores/<agent>/durable-state.json` under the current run root with CAS bytes,
compute checkpoints, capability tokens, replay hashes, and local journals.
Retention and GC are promise/evidence records (`retention_until_promised`,
`delete_after_promised`, `gc_object_retained`, `gc_object_removed`,
`retention_promise_broken`) rather than a central cleanup authority. Sender and
receiver rate/capacity behavior is modeled through `send_rate_promised`,
`accept_rate_promised`, and `backpressure_capacity_promised`; replay handling
records exact-envelope and serve-once-token rejections as local non-commitment
evidence. Analyzer scoring now includes durability, retention, pressure, and
replay dimensions, and tests cover within-run recovery, token replay, exact
envelope replay, CBOR fuzzing, delayed partial reads, and short writes. Source:
`DI-sunuf`.

`DI-fijov` tightens local trust recovery after malformed or broken peer
evidence. Corrupt CAS evidence now enters the same relationship-ledger path as
bad proofs and bad compute results. The ledger records per-peer recovery caution
inside the run-scoped relationship state, so ordinary kept promises after recent
malformed/broken evidence first work off caution before they can raise trust.
Future-only repair promises are retained as local evidence but do not
immediately prove that repair has already been kept. Source: `DI-fijov`.
