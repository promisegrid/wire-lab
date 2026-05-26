# poc4 relay protocol

The relay protocol is spoken by relay apps. A relay wrapper carries exact inner
envelope bytes from one relay hop toward a target node/app.

Payload fields:

- `kind`: `relay_forward_v1`
- `origin_node`: node where the wrapped promise started
- `origin_app`: app where the wrapped promise started
- `target_node`: node expected to deliver the inner envelope locally
- `target_app`: app expected to have promised the inner pCID
- `inner_hex`: exact inner envelope bytes encoded as hex
- `request_hash`: hash of the original request envelope, if this is a follow-up

Each relay promises only its local forwarding attempt to the next promised hop.
Relays do not create permission, authority, service discovery, or global routing.
