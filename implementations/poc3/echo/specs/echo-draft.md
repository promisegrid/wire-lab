# poc3 echo protocol

The echo protocol lets one app ask another app to make a new echo message. The
receiving echo app chooses whether to keep that local promise; this is not an
RPC command.

Payload fields:

- `kind`: `echo_request_v1` or `echo_response_v1`
- `from`: sending app name
- `from_node`: sending node name
- `to`: receiving node name
- `text`: echoed text
