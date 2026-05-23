# SIM-bihim-cas-type-codec-plus-internal-kind

This simulation turns the codec-plus-internal-kind object-type binding
alternative from `SIM-kohad-cas-object-type-binding-bakeoff` into a concrete
candidate specimen. It tests whether CAS identity should combine a broad CID
codec with a deterministic in-object kind field for finer protocol distinctions.
Source: `DI-fibuv`.

## Design Under Test

The CID codec promises the broad object family, while the object bytes carry a
small canonical kind field that promises the exact payload interpretation within
that family.

## Boundaries

This simulation does not choose final kind names or registry mechanics. It tests
whether an internal kind improves forward compatibility or creates two competing
type authorities.
