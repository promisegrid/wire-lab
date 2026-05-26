# Scenario walkthroughs

## Kernel porting boundary

Bob's app promises Bob's local kernel that it receives one hello-family pCID and keeps bounded evidence:

```text
{
  kind: "surface_promise_v1",
  promiser: "bob-hello-app",
  promisee_hint: "bob-local-kernel",
  surface: "receive",
  subject_pcid: "hello-pCID",
  promise_text: "I receive hello-pCID messages carried as grid([42(hello-pCID), payload, ...]) while this app is running.",
  message_boundary: "grid([42(pCID), payload, ...])",
  assumptions: ["local runtime keeps this app running"],
  non_promises: ["unknown pCIDs", "delivery after app stop"],
  evidence_promise: "retain exact received bytes until bounded evidence rotation",
  local_adapter: "onMessage envelope adapter",
  validity: "until app stop or evidence store rotation"
}
```

Bob's kernel may later record its own local observation:

```text
{
  kind: "promise_observation_v1",
  observer: "bob-local-kernel",
  subject_promiser: "bob-hello-app",
  promise_ref: "cid:promise-record",
  subject_pcid: "hello-pCID",
  surface: "receive",
  attempt_ref: "cid:exact-envelope-bytes",
  outcome: "kept",
  evidence_refs: ["cid:exact-envelope-bytes", "cid:delivery-note"],
  observed_at: "bob-local-event-1042",
  local_trust_effect: "continue attempting hello-pCID delivery while the app keeps this promise"
}
```

## Kept storage promise

Alice sends data to Bob only after Bob has made a bounded storage promise and Alice's local trust threshold is met. Alice records a kept observation only for Alice's own future judgment:

```text
{
  kind: "promise_observation_v1",
  observer: "alice",
  subject_promiser: "bob-storage-agent",
  promise_ref: "cid:bob-storage-promise",
  subject_pcid: "store-chunk-pCID",
  surface: "store",
  attempt_ref: "cid:alice-chunk-C",
  outcome: "kept",
  evidence_refs: ["cid:stored-byte-witness", "cid:bob-receipt"],
  observed_at: "alice-local-event-205",
  local_trust_effect: "Alice may continue sending similarly scoped data to Bob under this relationship"
}
```

## Refusal and non-promise

Carol asks Alice for data while promising to keep it private inside Carol's local group. Alice may refuse if Carol's promise history is insufficient, and Alice can record that refusal as Alice's own evidence rather than a negative fact about Carol:

```text
{
  kind: "promise_observation_v1",
  observer: "alice",
  subject_promiser: "alice",
  promise_ref: "cid:alice-selective-send-promise",
  subject_pcid: "selective-send-pCID",
  surface: "send",
  attempt_ref: "cid:carol-request-envelope",
  outcome: "refused",
  evidence_refs: ["cid:alice-refusal-envelope", "cid:carol-current-promise-history"],
  observed_at: "alice-local-event-244",
  local_trust_effect: "Alice does not send this data now; Alice may reconsider after different reciprocal promises or more kept history"
}
```
