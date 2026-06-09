# poc13-cas-compute-functions

`poc13-cas-compute-functions` is executable POC evidence for PromiseGrid storage,
compute, kernel, shipping, and live-agent behavior. It is intentionally a
superset repair of POC11 and POC12: it keeps POC11's autonomous sparse-mesh
relationship/economics pressure, keeps POC12's separate app/kernel processes and
shipping/device workflow, and adds POC13's decentralized CAS storage,
CID-named compute, replica recovery, token lifecycle, cache, verifier
disagreement, bad-proof, unknown-pCID, and evidence-report pressure. Source:
`DI-timah`; `DI-bikit`; `DI-galin`; `DI-pohaj`; `DI-zapab`; `DI-sinur`.

## What This Tests

- Multiple app pCIDs through one local container kernel: `relationship_v1`,
  `postal_scale_v1`, `ups_label_v1`, `printer_port_v1`, `accounting_v1`,
  `cas_storage_v1`, `cid_compute_v1`, and `evidence_report_v1`.
- Real app/kernel process boundary: each container runs one `poc13-kernel`
  process plus separate local app processes for relationship, fulfillment,
  postal scale, UPS label printer, printer port, and accounting roles.
- Kernel-style pCID routing: the kernel parses slot 0 `42(pCID)`, checks local
  app receive promises, and delivers exact bytes to the app process that
  promised the target pCID.
- Deterministic device/system agents for scale, printer-port hardware access,
  label-printer business logic, and accounting alongside live LLM
  business/social agents.
- One top-level semantic act, `promise`; workflow steps are payload meanings,
  not RPC verbs.
- Explicit direct TCP relationship transition evidence:
  `direct_peer_added`, `direct_peer_removed`, and `direct_peer_unchanged`.
- Corrected promise evidence semantics: ordinary `not_promised` /
  `non_commitment` evidence is neutral rather than broken peer-promise
  evidence, provider/runtime decision failures stay local, and duplicate
  shipment updates are checkpointed without repeated trust inflation.
  Source: `DI-jinoz`.
- App-local promise journal: POC13 now records an outstanding promise before
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
- Superset analyzer gates: `poc13-analyze` now fails if inherited POC11/POC12
  behavior or POC13 storage/compute evidence disappears. Source: `DI-sinur`.
- Observer-only monitor lifecycle: completed nodes wait for `monitor.done` using
  a config-derived provider/turn/grace budget, and a completed run can still
  write the non-authoritative marker if the observer report fails. Source:
  `DI-jupob`.

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

## Shipping Agents

- `poc13-fulfillment`: hybrid workflow coordinator that executes one deterministic
  startup shipment sequence across the production pCIDs, then continues normal
  live LLM relationship turns. It sends shipping pCIDs but receives only the
  relationship pCID in this POC.
- `poc13-postal-scale`: deterministic app for `postal_scale_v1`; promises
  package weight evidence only.
- `poc13-ups-label-printer`: deterministic app for `ups_label_v1`; promises
  label, cost, and tracking evidence only after it receives and redeems a local
  printer-port capability-promise token.
- `poc13-printer-port`: deterministic app for `printer_port_v1`; represents
  the local hardware-access kernel role for a simulated USB printer port. It
  promises bounded future printing by issuing a scoped token to the label
  printer, then promises print evidence when that token is redeemed with label
  bytes.
- `poc13-accounting`: deterministic app for `accounting_v1`; promises address
  lookup and shipment update evidence only.
- `poc13-relationship-agent`: generic live LLM relationship app. Depending on
  the configured agent, it can also promise `cas_storage_v1`,
  `cid_compute_v1`, or `evidence_report_v1` handling while keeping the same
  single top-level `promise` action.
- `poc13-kernel`: container-local transport process that records operational
  routing evidence only; it does not own trust, workflow, device behavior, or
  promise judgment. Source: `DI-galin`.

## Run

Copy the committed template and keep secrets out of it:

```sh
cp config.example.json config.json
printf '%s' "$OPENAI_API_KEY" > openai_api_key.txt
chmod 600 openai_api_key.txt
scripts/run-clean.sh
```

For manual runs, `poc13-analyze` accepts either the parent run directory or its
`run/` JSONL directory and fails loudly if no JSONL evidence exists:

```sh
docker compose up --build --abort-on-container-exit
docker compose run --rm --entrypoint /usr/local/bin/poc13-analyze dave /run/poc13/<run_id>
```

`config.json`, `poc13.env`, `openai_api_key.txt`, provider outputs, and Docker
volume state are ignored and must not be committed.

## Current Limits

POC13 is provisional executable evidence, not a stable shipping, device,
accounting, kernel-routing, monitor, trust, provider, or workflow API. Docker
networking remains static; dynamic TCP relationships are local promises to dial
or accept direct exchanges between app agents, carried by local kernels. The
kernel records only app receive-promise registration, exact-byte delivery,
peer-forwarding, unregistered pCID, and transport outcomes. Apps own trust,
keep/break/non-commitment judgment, relationship ledgers, workflow state, and
deterministic device/system behavior. The fulfillment startup sequence is a POC
guardrail so the run produces concrete pCID-routed shipment evidence instead of
relying on a live LLM to choose that sequence unaided. Source: `DI-timah`;
`DI-bikit`; `DI-parok`; `DI-galin`; `DI-pohaj`; `DI-sinur`.

The repaired POC13 should be treated as the current superset baseline for future
POCs unless a later scoped DI explicitly declares a non-superset exception and
lists the features being dropped. Source: `DI-sinur`.

Post-split evidence semantics are intentionally conservative:
`not_promised` means the receiver did not promise the requested exchange, not
that the receiver broke a promise. Transient provider/runtime decision failures
are local app/runtime evidence, not peer trust evidence. The accounting app
keeps an app-local shipment checkpoint keyed by order, tracking number, and
cost; duplicate confirmations remain visible in logs but do not repeatedly
increase trust. Source: `DI-jinoz`.

The 2026-06-05 `poc13-jinoz-20260605-055916` validation run exited cleanly and
the analyzer reported 454 events: 433 kept, 20 non-commitment, and 1 broken
resource promise. Shipping evidence appeared exactly once for address lookup,
package weight, label printing, and accounting update/confirmation. The monitor
still flagged a real remaining concern: some trust changes are driven by local
budget/capacity exhaustion and should be separated from peer trust in a later
POC. Source: `DI-jinoz`.

`DI-vujob` keeps those follow-up corrections in POC13. The runtime now records
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
2026-06-05 `poc13-vujob-20260605-063213` Docker run exited cleanly and analyzed
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

`DI-pohaj` adds a local printer-port kernel-role app without turning the
message kernel into a USB authority or RPC service. During label printing,
`ups_label_printer` sends `issue_print_capability` to `printer_port`;
`printer_port` returns a deterministic token bound to the issuee, scope, token
ID, and byte limit; `ups_label_printer` sends `redeem_print_capability` with the
token and exact hex label bytes; `printer_port` returns deterministic local
spool evidence. Unit tests cover the token issue/redemption path, wrong-token
rejection, routing pCID registration, and analyzer shipping-event recognition.
The fresh Docker run `poc13-jupob-20260606-030610` exited cleanly after Dave
wrote `monitor_done`; analyzer output reported 625 events, 17 `node_done`, 17
`shutdown_grace_elapsed`, one `monitor_done`, all printer-port shipping counts,
and empty `resource_trust_coupling_counts`. Source: `DI-pohaj`; `DI-vutok`;
`DI-jupob`.
