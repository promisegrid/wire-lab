# SIM-zubuf-cas-type-cid-codec-only

This simulation turns the CID-codec-only object-type binding alternative from
`SIM-kohad-cas-object-type-binding-bakeoff` into a concrete candidate specimen.
It tests whether the CIDv1 codec or multicodec value should be the sole object
type discriminator for raw chunks, pointer objects, Merkle nodes, and later
application-shaped CAS objects. Source: `DI-fibuv`.

## Design Under Test

Object bytes promise their broad interpretation through the CID codec only.
There is no second in-object kind field, so type identity is bound entirely by
the content-addressed codec selection.

## Boundaries

This simulation does not allocate final codec values. It tests whether one
codec-level discriminator avoids duplicate sources of truth or becomes too
coarse for PromiseGrid object evolution.
