# SIM-guhum-transport-cluster-of-clusters

This simulation turns the cluster-of-clusters transport alternative from
`SIM-narok-transport-family-bakeoff` into a concrete candidate specimen. It
tests whether clusters should be first-class routing and policy units that can
compose into larger PromiseGrid transport structures. Source: `DI-fibuv`.

## Design Under Test

Each cluster promises local membership, routing, and relay behavior. Clusters
then exchange summary promises with other clusters rather than exposing every
member peer to global transport state.

## Boundaries

This simulation does not define final addressing or governance for clusters. It
tests whether hierarchy reduces scale pressure or introduces hidden authority
and migration complexity.
