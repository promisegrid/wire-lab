# App/kernel local promises draft

> Status: DRAFT. Simulation-local payload protocol only. Source: `DI-dikat`; `DI-bozid`.

## Purpose

This single payload protocol lets agents make and remember local app/kernel promises without turning the kernel into an RPC registry or capability authority. One stable pCID names this protocol family. Message-level variation is carried by the `kind` field inside the payload, not by minting a separate pCID for each closely related record type.

The usual case is an app promising its local kernel which pCID-selected messages it will receive or handle. The same protocol can also carry local observations about later kept, refused, not-promised, broken, timed-out, unreadable, or host-unavailable outcomes.

The records are not registry entries and not capability grants. They are promises and observations by named agents. Other agents may remember them, compare promises with later observations, and update their own local trust.

## Payload shape

A deterministic CBOR map has a required `kind` field and then the fields for that record kind.

```text
app_kernel_local_promises_payload = {
  kind: "app_pcid_promise_v1" | "local_observation_record_v1",
  ... kind-specific fields ...
}
```

The protocol pCID is stable for both record kinds. A future incompatible change to this whole protocol family gets a new pCID; compatible additions should be optional fields or new `kind` values defined by a successor spec.

## `app_pcid_promise_v1`

An `app_pcid_promise_v1` record is one agent's promise about its own behavior for one pCID-selected message family.

Fields:

- `kind`: `app_pcid_promise_v1`
- `promiser`: the agent making the promise
- `observer_hint`: optional intended local observer, such as the kernel or peer this promise is made to
- `subject_pcid`: the pCID this promise concerns
- `promise_role`: one of `receive`, `handle`, `store`, `compute`, `send`, `device-effect`, `evidence`, `namespace-view`, `reference-view`, or another explicitly named local role
- `promise_text`: short human-readable promise wording
- `message_boundary`: normally `grid([42(pCID), payload, ...])`
- `assumptions[]`: facts the promiser depends on but does not promise
- `non_promises[]`: pCIDs, roles, or conditions the promiser explicitly does not promise
- `evidence_promise`: what exact bytes, CIDs, timestamps, or local notes the promiser promises to retain or emit
- `local_adapter`: optional local API wrapper name, when the wrapper is merely an adapter over pCID-selected messages
- `validity`: local time or event bounds for the promise
- `notes`: optional explanatory text

Rules:

- The `promiser` promises only its own behavior.
- `subject_pcid` is Protocol CID, not payload CID.
- `observer_hint` does not bind anyone else; it only names the relationship where the promise is expected to be useful.
- `assumptions[]` are not promises unless another agent separately promises them.
- `non_promises[]` are first-class evidence. A direct non-promise should not be treated as silence.
- `local_adapter` must not imply an RPC method that bypasses the pCID-selected message boundary.

Example reading:

> Bob's hello app promises Bob's kernel: I receive and handle `hello-pCID` messages sent to me through `grid([42(hello-pCID), payload, ...])`; I do not promise to receive unknown pCIDs; I keep exact received bytes for local debugging until my bounded evidence store rotates.

Bob's kernel does not thereby certify a global capability. It only has a local promise to use when deciding whether to attempt delivery and how to evaluate later make/break history.

## `local_observation_record_v1`

A `local_observation_record_v1` record is one observer's memory of what it observed about one promise attempt. It is local evidence for local trust decisions. It is not a global verdict, global reputation fact, or authority record.

Fields:

- `kind`: `local_observation_record_v1`
- `observer`: the agent recording the observation
- `subject_promiser`: the agent whose promise was being assessed
- `promise_ref`: CID or exact-byte reference to the relevant promise record, if available
- `subject_pcid`: the pCID involved in the observed message or attempt
- `attempt_ref`: CID, exact bytes, or local attempt identifier for the observed message or action
- `outcome`: one of `kept`, `refused`, `not-promised`, `broken`, `timed-out`, `unreadable`, `host-unavailable`, or another locally named outcome
- `evidence_refs[]`: exact bytes, CIDs, local logs, or witness notes retained by the observer
- `observed_at`: observer-local timestamp or event reference
- `local_trust_effect`: optional local-only note about how this observation may affect future send, receive, store, compute, keep, pull, or advertise choices
- `notes`: optional explanatory text

Outcome rules:

- `kept`: the observer saw behavior matching the promise closely enough for its local judgment.
- `refused`: the observer has explicit refusal evidence from the subject promiser.
- `not-promised`: the observer has evidence that the subject promiser did not make the requested promise.
- `broken`: the observer has evidence that the subject promiser made the promise and then did not keep it.
- `timed-out`: the observer did not receive kept or refused evidence before its local deadline.
- `unreadable`: bytes arrived but could not be parsed or checked as the relevant pCID protocol required.
- `host-unavailable`: the observer saw a host/runtime condition that the promiser had named as an assumption rather than a promise.

A receiver must not collapse `refused`, `not-promised`, and `timed-out` into one bucket when evidence distinguishes them. Refusal and non-promise records can preserve trust better than silence because they are explicit promises about what will not be done.

`local_trust_effect` is not portable authority. Alice may decide not to send future data to Bob after a broken storage promise. Carol may increase trust in Bob after repeated honest refusals. Dave may ignore both observations because Dave's trust relationship and evidence threshold are different.
