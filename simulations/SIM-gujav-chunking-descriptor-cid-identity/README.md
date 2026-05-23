# SIM-gujav-chunking-descriptor-cid-identity

This simulation turns the chunking-descriptor identity alternative from
`SIM-gobaz-chunking-identity-bakeoff` into a concrete candidate specimen. It
tests whether a separate content-addressed or protocol-addressed descriptor
should identify the chunking algorithm and parameters used for a stream.
Source: `DI-fibuv`.

## Design Under Test

Objects carry or reference a chunking descriptor CID that promises the exact
algorithm, parameter set, and boundary rules used to produce CAS leaves and
roots.

## Boundaries

This simulation does not settle the `cCID` name. It tests whether a distinct
descriptor improves reuse and migration or adds another indirection that weakens
simple object identity.
