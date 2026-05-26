# poc3 signed protocol

The signed protocol carries a message plus a proof slot. The signature witnesses
exact signable envelope bytes. It does not create global trust.

Payload fields:

- `kind`: `signed_note_v1`
- `from`: sending app name
- `from_node`: sending node name
- `to`: receiving node name
- `text`: note text

Envelope slot 2 carries a POC Ed25519 public-key/signature byte string.
