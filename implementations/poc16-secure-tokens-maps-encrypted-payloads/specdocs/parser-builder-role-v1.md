# parser_builder_role_v1

Agents use `parser_builder_role_v1` to describe local parser and builder roles
selected by slot-0 pCID. The parser role decodes payload bytes and talks to local
apps; it does not turn pCID into a destination address or force the kernel to
understand every payload's addressing scheme.
