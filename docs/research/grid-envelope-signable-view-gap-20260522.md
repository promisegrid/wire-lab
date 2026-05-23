# Grid-Envelope Signable-View Gap Note

Date: 2026-05-22
Owner TODO: `TODO-tugoz`
Source: `DI-dunat`; `DI-kafot`

## Purpose

Pin down what the next explicit signable-view simulation must add beyond the
nearest existing envelope specimens.

## Existing Near Neighbors

### `SIM-janov`

- Outer shape is still `[pcid_a, payload_a]`.
- `pcid_a` defines `payload_a` as nested canonical bytes
  `[pcid_b, payload_b, signature_b]`.
- The nested signature covers `[pcid_b, payload_b]`.
- Main strength: keeps the universal outer parse surface minimal while binding
  the inner payload bytes to their inner pCID.
- Main weakness: the outer conformance promise is still effectively attributed
  by transport or surrounding context because the outer layer has no proof of
  who promised that `payload_a` really matches the `pcid_a` contract.

### `SIM-riliz`

- Outer shape is `[pcid, payload, sig_pcid, sig_payload]`.
- The signature/proof is a universal mandatory outer slot pair.
- Main strength: generic peers can always see that a proof-bearing wrapper is
  present, and the default unsigned prefix is explicit.
- Main weakness: it solves the signable-view question by adding universal outer
  slots, which is precisely what the new dedicated signable-view specimen is
  trying not to assume.

### `SIM-gojot`

- Outer shape is `[pcid, payload]`.
- Any proof is carried by a wrapper protocol chosen by `pcid`.
- Main strength: maximally small outer envelope.
- Main weakness: it does not isolate the case where a single payload contract
  itself carries the proof and defines the exact unsigned/signable projection.

### `SIM-maraz`

- Builds on Janov-style nested payloads, but adds a fixed signed outer summary
  header.
- Main strength: bounded outer audit hooks for sparse peers.
- Main weakness: it is already a richer audit-header specimen, not the minimal
  "proof carried inside the pCID-owned payload with an explicit signable view"
  case.

### `SIM-natim`

- Keeps a small outer envelope but adds a mandatory outer attestation proof.
- Main strength: better attribution of origin/relay/holder promises.
- Main weakness: it explores outer attestation and multisig pressure, not the
  simpler inside-payload signable-view rule.

### `SIM-sajar`

- Outer shape is variable arity: `[pcid, field_1, ... field_n]`.
- Main strength: demonstrates that pCID-defined arity can push flexibility into
  schema-local rules.
- Main weakness: it moves too much of the universal parse burden into the first
  pCID to answer the narrow signable-view question cleanly.

## Exact Gap

The corpus still lacks one dedicated specimen with all of these properties at
once:

1. the universal outer envelope stays exactly `grid([pCID, payload])`;
2. the payload contract itself requires a proof field inside the payload;
3. the payload contract defines an explicit unsigned/signable projection so the
   proof is not self-referential;
4. the specimen keeps that rule minimal enough to compare directly against
   `SIM-riliz` (mandatory outer proof), `SIM-gojot` (wrapper-proof), and
   `SIM-janov` (nested signed payload under a layer pCID).

The missing question is not "can a pCID define fields?" The corpus already says
yes. The missing question is narrower: what is the simplest viable way for a
payload-pCID to say "this payload carries its own proof, and these are the exact
bytes that are signed" without adding universal outer signature slots or a
heavier audit header.

## What The New Signable-View Sim Must Add

The new specimen should add these comparison pressures explicitly:

- **Named unsigned/signable view.** The payload contract should define an
  explicit projection such as `payload_without_sig`, not an implicit "remove the
  last field if present" rule.
- **Single-layer comparison.** The specimen should avoid Janov's extra
  `pcid_a` / `pcid_b` layering so the signable-view rule can be evaluated on its
  own merits.
- **Proof-inside-payload rule.** The proof belongs to the payload contract, not
  to a universal outer tuple and not to a separate generic audit header.
- **Outer-envelope minimalism.** Generic peers still only parse
  `[pcid, payload]`; they do not gain a new universal proof slot.
- **Bounded unknown behavior.** Unknown-pCID peers may retain or blind-carry the
  bytes, but they cannot verify the inner proof or claim semantic acceptance.

## Recommended First Dedicated Shape To Test

This note does not lock the winning shape, but the cleanest first dedicated
specimen to compare is:

- outer envelope: `grid([pcid, payload])`
- payload shape defined by that `pcid`:
  - `[payload_core, sig_pcid, sig_payload]`
- explicit signable view:
  - `payload_without_sig = [payload_core]`
  - or, if `sig_pcid` must also be signed as associated data,
    `payload_without_sig = [payload_core, sig_pcid]`

The point is not which of those two projections wins yet. The point is that the
projection must be explicit, named, and reconstructable from the payload spec
without relying on transport identity or outer universal proof slots.

## Minimum Comparison Set

The new signable-view specimen should be scored directly against:

- `SIM-gojot` for wrapper-proof minimalism
- `SIM-riliz` for mandatory outer proof slots
- `SIM-janov` for nested signed payload layering

`SIM-maraz` and `SIM-natim` remain useful follow-on pressure tests, but they are
already more elaborate than the narrow gap now being isolated.
