# SIM-lilok-cas-type-path-suffix-negative-control

This simulation turns the path-suffix object-type binding negative-control
alternative from `SIM-kohad-cas-object-type-binding-bakeoff` into a concrete
candidate specimen. It tests the failure mode where local filenames or path
suffixes, rather than content-addressed bytes, carry object type meaning.
Source: `DI-fibuv`.

## Design Under Test

Local storage paths promise object type by suffix, for example treating the same
CID bytes differently when stored as different filenames.

## Boundaries

This is a negative-control specimen. It should remain scoreable pressure for the
GA, but it is expected to expose why path-local names cannot be the primary CAS
type binding when content moves across peers and archives.
