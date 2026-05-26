# poc5 signed app draft

`poc5-signed` uses `grid([42(pCID), payload, sig])` for signed results and
`grid([42(pCID), payload])` for sign requests. The pCID names this draft.

The app-level promise is deliberately narrow:

- A `sign_request_v1` asks one signed app to sign a text statement and return a
  `signed_result_v1` to the app named by `from_node` and `from`.
- A `signed_result_v1` carries the request hash, signed text, and one proof
  slot. The promisee judges the result locally from exact bytes and signature
  verification.
- No kernel or relay treats the signature as authority. It is evidence a local
  observer may use when updating trust.

