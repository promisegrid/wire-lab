# POC13 Run Narrative

POC13 is now the current executable superset baseline for POC11, POC12, and the
original POC13 CAS/compute work. Source: `DI-sinur`.

## Expected Message Sequence

1. Every app process registers receive promises for its configured pCIDs with
   the local container kernel. The kernel records only operational delivery
   evidence; each app owns trust, economics, workflow, and keep/break judgment.
2. Alice performs CAS storage pressure:
   - Alice first runs a dynamic TCP topology probe against fulfillment: local
     broken-promise evidence removes direct reachability, an ordinary promise is
     blocked before kernel send, repair evidence restores direct reachability,
     and a later `relationship_v1` promise crosses the real app/kernel path.
   - Alice offers Bob too little credit for `store_content`.
   - Bob records `economics_price_refused` as local non-commitment evidence.
   - Alice offers enough credit; Bob stores exact bytes, issues a serve token,
     and asks Frank to store a replica.
   - Frank stores the replica and issues Alice a replica token.
   - Alice retrieves once from Bob, then records Bob as unavailable for a later
     modeled outage and retrieves from Frank by redeeming Frank's token.
3. Alice performs CID-named compute pressure:
   - Alice asks Dave for a cache entry and receives a miss.
   - Alice asks Carol to execute a function over explicit function/input/context
     CIDs.
   - Carol executes bounded function bytes, returns result bytes and a
     deliberately bad-result probe, and records the compute promise.
   - Alice recomputes locally, accepts the good result, rejects the bad-result
     probe, and records credits spent or withheld accordingly.
   - Alice asks Dave and Grace to verify. Dave verifies and checkpoints cache
     evidence; Grace produces disagreement pressure; Alice resolves locally.
   - Alice asks Dave again and receives a cache hit.
4. Mallory exercises adversarial pressure:
   - Mallory sends corrupt bytes to Grace.
   - Mallory sends a future repair promise, an unknown pCID, an unsupported
     variant, a deliberately bad proof, a key-rotation promise, and a compute
     capacity probe.
   - Grace and Carol record local non-commitment, malformed, or evidence-report
     outcomes without treating Mallory's text as authority.
5. Fulfillment performs the POC12 shipping workflow:
   - Fulfillment asks accounting for an address and postal_scale for weight.
   - Fulfillment asks ups_label_printer for label/cost/tracking evidence.
   - ups_label_printer asks printer_port for a future-print promise token and
     redeems it with bounded label bytes.
   - Fulfillment updates accounting and repeats the same update once so the
     duplicate checkpoint path is exercised.

## Acceptance Meaning

`poc13-analyze` enforces the repaired superset by requiring:

- POC11/POC12 inherited monitor, relationship, non-commitment, checkpoint, and
  shipping evidence.
- POC13 CAS storage, replica, token lifecycle, corrupt-byte, compute, cache,
  verifier disagreement, bad-proof, unknown-pCID, key-rotation, economics, and
  local-trust evidence.
- `DI-fijov` trust-caution behavior: malformed/broken evidence records recovery
  caution, ordinary kept evidence is delayed while caution remains, and
  future-only repair promises do not immediately restore trust.
- Dynamic topology behavior: removed direct relationships block ordinary
  sends, and restored relationships carry actual kernel-routed messages again.
- Empty resource/trust coupling counts.
- No obvious RPC drift terms in event names or event details.

POC13 remains provisional executable evidence, not a stable PromiseGrid API.

## Latest Clean Run: `poc13-demo`

The 2026-06-11 hardened clean run passed and produced 1,948 events. Analyzer
scoring was 5/5 across all deterministic dimensions, including trust-caution
and dynamic-topology gates. The run recorded 64 non-commitment outcomes,
3 `trust_caution_recorded`, 4 `trust_caution_consumed`, 1
`trust_recovery_delayed`, 1 `trust_repair_future_only`, and one full dynamic
topology sequence: probe started, send blocked, send succeeded after repair.
The observer-only monitor scored autonomy 5 and scored promise-theory fit,
protocol validity, local-trust correctness, and imposition avoidance at 4. The
remaining production-fitness blockers are protocol cleanliness around
adversarial pressure, shutdown coordination, and evidence/log clarity, not
missing POC13 regression coverage.
Source: `DI-sihuz`.
