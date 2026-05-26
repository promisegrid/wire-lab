# poc3 hello protocol

The hello protocol carries a small message from one app agent to another.

Payload fields:

- `kind`: `hello_v1`
- `from`: sending app name
- `from_node`: sending node name
- `to`: receiving node name
- `text`: hello text

The receiving app judges the message locally. The kernel only routes bytes and
records its own receive/deliver/refusal evidence.
