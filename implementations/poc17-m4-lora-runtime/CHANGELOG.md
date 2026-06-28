# POC17 Changelog

## Unreleased

- Scaffold the Go behavior simulator for the POC17 M4/LoRa runtime.
- Add bintags-shaped order status traffic under `order_status_v1`.
- Add RFC-like embedded protocol specs and derive slot-0 pCIDs from their exact
  bytes.
- Raise the simulated LoRa application MTU to 200 bytes.
- Convert printable artifact, CAS, parent-link, event, and spec-alias
  identifiers to canonical CID text.
- Add host-local `local_lifecycle_v1` CWT/COSE lifecycle tokens and resource
  withdrawal evidence without sending token bytes over LoRa.
- Add config-driven resource-limit snapshots, fresh-agent restart recovery, and
  radio-visible `peer_storage` grant/put/get promises using Bob-issued compact
  capability tokens.
- Remove synthetic resource-activity usage values from POC17 evidence; only
  configured limits are reported until activity is actually measured.
