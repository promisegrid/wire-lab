# DR-jalom - Grid-envelope protocol-owned signature slot

DR-ID: DR-jalom
Date: 2026-05-23 23:27:04
Asked by: jessica@t7a.org (Jessica)
State: decided
Question: Should wire-lab add a new standalone simulation for the three-slot
outer envelope `[pcid, payload, signature]` where `pcid` owns the signature
rules, and should that sim be compared against existing outer-shape families as
the next preferred-default outer-shape question?
Why this blocks progress: The repo has close neighbors — mandatory opaque bytes,
outer `sig_pcid`, payload-owned proof, wrapper proof, and outer attestation —
but it does not have a direct specimen for the cleaner “three-slot outer
signature, proof semantics owned by `pcid`” design. Without a tracked specimen,
the repo keeps discussing the shape in chat without generating comparable
simulation evidence.
Affects: `protocols/wire-lab.d/TODO/TODO-mujad-grid-envelope-protocol-owned-signature-slot.md`;
`simulations/SIM-dalor-grid-envelope-protocol-owned-signature-slot/`;
`simulations/README.md`; `DEV-GUIDE-RESOURCES.md`.
Unblocks: TODO-mujad implementation; later scored comparison against existing
outer-shape families.
Waiting on: DI-kukuk
Decision: Add one new standalone sim for `[pcid, payload, signature]` with
proof semantics owned by `pcid`, and compare it against existing
`[pcid,payload]`, explicit `sig_pcid`, payload-owned proof, and wrapper-proof
families. Reuse the existing six-scenario slice for the first scored pass.
Linked DI: DI-kukuk
Related commits:
Last updated: 2026-05-23 23:27:04

## Event log

- 2026-05-23 23:27:04 — Opened and decided from chat after the user approved a
  broad five-family comparison, one new sim only, the existing six-scenario
  slice, and the exact rule that the third slot signs canonical `[pcid,
  payload]` bytes.

## Evidence

- The positional matrix already contains multiple `[pcid, payload, signature]`
  specimens with generic mandatory-opaque-bytes wording.
- `SIM-jufag` and siblings cover the explicit outer `sig_pcid` family.
- `SIM-pamap` covers payload-owned proof with an explicit signable view.
- `SIM-jumav` covers wrapper proof over linked content.
- `SIM-lotiv` keeps multisig placement modes open but does not isolate the
  simpler pCID-owned three-slot envelope as a standalone specimen.

