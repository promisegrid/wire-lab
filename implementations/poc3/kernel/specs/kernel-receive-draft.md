# poc3 kernel receive-promise protocol

This protocol is local to the `poc3` proof of concept.

Payload fields:

- `kind`: `receive_promise_v1`
- `app`: the local app name making the receive promise
- `node`: the local node name
- `pcid`: the Protocol CID string the app promises to receive
- `text`: human-readable promise wording

The kernel treats this as a local promise by an app. It is not permission,
authorization, conformance, or a global registry.
