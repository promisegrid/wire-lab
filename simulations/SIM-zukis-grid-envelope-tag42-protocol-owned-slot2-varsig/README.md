# SIM-zukis-grid-envelope-tag42-protocol-owned-slot2-varsig

This simulation is a direct specimen of the locked `TE-fikoj` /
`DI-sisak` outer-envelope direction. It does not reopen the fixed-three-slot
versus variable-arity decision. Instead, it makes that decision concrete with
one protocol-owned example:

```text
grid([42(pCID), payload, varsig])
```

For this specimen:

- slot `0` is always the tagged protocol selector `42(pCID)`;
- slot `1` is the primary payload anchor;
- slot `2` is this protocol's own `varsig` proof slot.

The broader family rule remains larger than this one specimen: PromiseGrid's
current direction is `grid([42(pCID), payload, ...])`, and the protocol named
by `pCID` defines whether later outer slots exist and what they mean. `SIM-zukis`
therefore tests one direct, PT-clean member of that family without turning
slot `2 = varsig` into universal envelope law. Source: `DI-sisak`; `DI-mabit`.

## Promise-Theory framing

The outer envelope helps a receiver interpret another agent's promise. It does
not command behavior, grant global permission, or decide trust centrally. In
this specimen, the current sender's `varsig` is evidence for the sender's own
scoped promise:

> "I promise these payload bytes and this outer-slot arrangement meet the
> protocol specification named by this `pCID`."

Each receiver still decides locally whether it recognizes the protocol, trusts
the sender, verifies the `varsig`, stores the bytes, relays them, or uses the
payload. Carriage is not semantic acceptance. Source: `DI-pagin`; `DI-sisak`;
`DI-mabit`.

## What this sim is testing

This sim tests whether a tagged selector in slot `0`, a stable payload anchor
in slot `1`, and one protocol-owned proof slot in slot `2` give PromiseGrid a
good balance of:

- DAG-CBOR / CID ecosystem interop;
- small deterministic outer parsing;
- protocol-owned evolution of later outer slots;
- clean separation between base-envelope promises and higher-layer promise
  accounting.

## Comparison targets

Primary comparison targets:

- `SIM-dalor-grid-envelope-protocol-owned-signature-slot`
- `SIM-pobod-grid-envelope-outer-promise-nested-signed-payload`
- `SIM-jufag-grid-envelope-quarantine-sig-pcid-outcomes`

`SIM-dalor` is the nearest fixed-three-slot neighbor. `SIM-pobod` pressures
explicit outer promise wording and nested signed payload structure. `SIM-jufag`
is the contrasting explicit-`sig_pcid` selector-shopping branch. Source:
`DI-mabit`.

## Boundaries

This sim does not declare that every PromiseGrid protocol must use slot `2`
as `varsig`. It only tests whether one direct specimen of the locked
`grid([42(pCID), payload, ...])` family performs well when the protocol named
by `pCID` chooses that shape for itself.
