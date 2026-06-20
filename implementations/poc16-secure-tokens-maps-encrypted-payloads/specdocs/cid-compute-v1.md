# cid_compute_v1

## Status

Active POC16 CID-named computation protocol. The embedded Markdown bytes derive
this pCID. Source: `DI-bitug`; `DI-gahuh`.

## Abstract

`cid_compute_v1` lets an agent promise to compute a result from explicit function,
input, and context CIDs; promise cache lookups; provide context bytes; and verify
another agent's result. Computation is trustable only to the extent the requester
can verify locally or through trusted peer promises.

## pCID and envelope

The active shape is:

```text
grid([42(pCID), payload, proof])
```

Function, input, context, and result CIDs are payload fields. The pCID names this
compute protocol, not a function address.

## Promise Theory model

A compute provider promises only its own computation attempt and result. A
verifier promises only its own verification result. A requester may promise
credits after local verification. Bad results lower trust only when a local agent
judges an outstanding promise was broken.

## Payload grammar

The payload is the strict pCID-owned array profile:

```text
[promiser, promisee, promise_about, [outcome, promise_text, reason], body]
```

| promise_about | Body slots |
|---|---|
| `execute_function` | `compute_status`, `exchange_id`, `function_cid`, `function_b64`, `input_cid`, `input_b64`, `context_cid`, `context_b64`, `result_cid`, `result_b64`, `bad_result_cid`, `bad_result_b64`, `credit_offer`, `units`, `capacity_probe` |
| `lookup_compute_cache` | `cache_key`, `exchange_id`, `cache_status`, `function_cid`, `function_b64`, `input_cid`, `input_b64`, `context_cid`, `context_b64`, `result_cid`, `result_b64` |
| `provide_compute_context` | `function_cid`, `exchange_id`, `input_cid`, `context_cid`, `context_b64` |
| `verify_compute_result` | `verdict`, `exchange_id`, `subject_peer`, `subject_result_cid`, `result_promiser`, `function_cid`, `function_b64`, `input_cid`, `input_b64`, `context_cid`, `context_b64`, `result_cid`, `result_b64`, `disagreement_probe` |

## Sender behavior

A compute request SHOULD include exact CIDs and bytes for every object not already
available to the receiver. A compute result MUST include `result_cid` and
`result_b64` when promising a result. A capacity probe is a request for a promise,
not an obligation to compute.

## Receiver and parser behavior

A parser MUST reject unsupported `promise_about`, wrong body length, non-text
slots, or trailing CBOR. A receiver MUST verify byte/CID pairs before using them.
A verifier MUST state `verdict` from its own local computation or validation; it
must not speak for the original result promiser.

## Protocol state machine

```text
[inputs unknown]
    | provide_compute_context / bytes available
    v
[inputs available] --execute_function accepted--> [compute running]
      | capacity refused                          |
      v                                           | result_cid verified
[not promised]                                    v
                                            [result locally accepted]
                                                   |
                                  verifier disagreement / bad result
                                                   v
                                            [local trust update]
```

## State, CAS, DAG, and retention

Implementations SHOULD store function/input/context/result bytes in local CAS and
link compute results to their request parent. Cache keys are local indexes over
pCID, function CID, input CID, context CID, and result CID; cache hits are local
promises, not global facts.

## Security considerations

Never trust `result_cid` without verifying `result_b64` or recomputing. WASM or
other executable functions must run under local resource promises and sandboxing.
Payment should be tied to local acceptance, not mere receipt.

## Interoperability notes

The protocol can describe native, WASM, external, or manual computation as long
as exact function/input/context/result objects are named by CID.

## Examples

```text
grid([42(pCID),
  ["carol", "alice", "execute_function",
    ["kept", "I promise this result is the output of the named function over the named input and context.",
     "compute completed"],
    ["computed", "cmp-9", "cid:function", "...", "cid:input", "...",
     "cid:context", "...", "cid:result", "...", "", "", "alice-credit:5", "5", ""]
  ], proof
])
```
