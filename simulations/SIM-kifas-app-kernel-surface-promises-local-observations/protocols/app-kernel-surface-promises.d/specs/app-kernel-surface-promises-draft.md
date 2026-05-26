# App/kernel surface promises draft

> Status: DRAFT. Simulation-local payload protocol only. Source: `DI-pavub`.

## Purpose

This single payload protocol lets agents make and remember local app/kernel promises at the level of explicit surfaces without turning the kernel into external authority. One stable pCID names this protocol family. Message-level variation is carried by a `kind` field inside the payload, not by minting a separate pCID for each closely related record type.

The usual case is an app promising its local kernel which pCID-selected messages it will receive or handle. The same protocol can express local storage, compute, send, key-use, lifecycle, namespace-view, reference-view, evidence, and device-effect promises when the promiser is explicit and the surface is local to that promiser.

The records are promises and observations by named agents. Other agents may remember them, compare promises with later observations, and update their own local trust. No record in this protocol speaks for an agent that did not make the promise or observation.

## Payload shape

A deterministic CBOR map has a required `kind` field and then the fields for that record kind.

```text
app_kernel_surface_promises_payload = {
  kind: "surface_promise_v1" | "promise_observation_v1",
  ... kind-specific fields ...
}
```

The protocol pCID is stable for both record kinds. A future incompatible change to this whole protocol family gets a new pCID; compatible additions should be optional fields or new `kind` values defined by a successor spec.

## `surface_promise_v1`

A `surface_promise_v1` record is one agent's promise about its own behavior for one local surface and one pCID-selected message family or resource role.

Fields:

- `kind`: `surface_promise_v1`
- `promiser`: the agent making the promise
- `promisee_hint`: optional intended local promisee or observer, such as the kernel, app, peer, or operator this promise is made to
- `surface`: one local surface name, such as `receive`, `handle`, `store`, `compute`, `send`, `key-use`, `lifecycle`, `namespace-view`, `reference-view`, `evidence`, or `device-effect`
- `subject_pcid`: the Protocol CID this promise concerns, when the surface is pCID-selected
- `promise_text`: short human-readable promise wording
- `message_boundary`: normally `grid([42(pCID), payload, ...])`
- `assumptions[]`: facts or conditions the promiser depends on but does not promise
- `non_promises[]`: pCIDs, surfaces, conditions, or resource roles the promiser explicitly does not promise
- `evidence_promise`: what exact bytes, CIDs, timestamps, bounded local notes, or witness records the promiser promises to retain or emit
- `local_adapter`: optional local API wrapper name, when the wrapper is merely an adapter over pCID-selected messages
- `validity`: local time, event, or frontier bounds for the promise
- `notes`: optional explanatory text

Rules:

- The `promiser` promises only its own behavior.
- `subject_pcid` is Protocol CID, not payload CID.
- `promisee_hint` does not bind anyone else; it only names the relationship where the promise is expected to be useful.
- `assumptions[]` are not promises unless another agent separately promises them.
- `non_promises[]` are first-class evidence. A direct non-promise should not be treated as silence.
- `local_adapter` must not imply a method call that bypasses the pCID-selected message boundary.
- Omitted surfaces mean no promise about that surface.
- Split local implementations add more `surface_promise_v1` records with their own promisers instead of hiding multiple agents behind one kernel-shaped name.

Example reading:

> Bob's hello app promises Bob's kernel: I receive and handle `hello-pCID` messages sent to me through `grid([42(hello-pCID), payload, ...])`; I do not promise to receive unknown pCIDs; I keep exact received bytes for local debugging until my bounded evidence store rotates.

Bob's kernel does not thereby certify a global ability. It only has a local promise to use when deciding whether to attempt delivery and how to evaluate later make/break history.

## `promise_observation_v1`

A `promise_observation_v1` record is one observer's memory of what it observed about one promise attempt. It is local evidence for local trust decisions. It is not a global verdict, global reputation fact, or authority record.

Fields:

- `kind`: `promise_observation_v1`
- `observer`: the agent recording the observation
- `subject_promiser`: the agent whose promise was being assessed
- `promise_ref`: CID or exact-byte reference to the relevant promise record, if available
- `subject_pcid`: the pCID involved in the observed message or attempt
- `surface`: the local surface involved, when known
- `attempt_ref`: CID, exact bytes, or local attempt identifier for the observed message or action
- `outcome`: one of `kept`, `refused`, `not-promised`, `unavailable`, `broken`, `timed-out`, `unreadable`, `corrupt-observation`, or another locally named outcome
- `evidence_refs[]`: exact bytes, CIDs, local logs, or witness notes retained by the observer
- `observed_at`: observer-local timestamp or event reference
- `local_trust_effect`: optional local-only note about how this observation may affect future send, receive, store, compute, keep, pull, or advertise choices
- `notes`: optional explanatory text

Outcome rules:

- `kept`: the observer saw behavior matching the promise closely enough for its local judgment.
- `refused`: the observer has explicit refusal evidence from the subject promiser.
- `not-promised`: the observer has evidence that the subject promiser did not make the requested promise.
- `unavailable`: the observer saw temporary inability reported by the subject promiser without claiming success.
- `broken`: the observer has evidence that the subject promiser made the promise and then did not keep it.
- `timed-out`: the observer did not receive kept, refused, not-promised, or unavailable evidence before its local deadline.
- `unreadable`: bytes arrived but could not be parsed or checked as the relevant pCID protocol required.
- `corrupt-observation`: the observer has evidence that the received bytes or witness data are corrupt or internally inconsistent.

A receiver must not collapse `refused`, `not-promised`, `unavailable`, and `timed-out` into one bucket when evidence distinguishes them. Refusal and non-promise records can preserve trust better than silence because they are explicit evidence about what will not be done.

`local_trust_effect` is not portable authority. Alice may decide not to send future data to Bob after a broken storage promise. Carol may increase trust in Bob after repeated honest refusals. Dave may ignore both observations because Dave's trust relationship and evidence threshold are different.
