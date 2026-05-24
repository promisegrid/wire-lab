# SIM-lurip-capability-promise-token-payload

This higher-layer simulation tests capability-like behavior as a promise token
payload. It preserves transfer and capability pressure from the rejected
`hadit`/`jogoh` review while keeping capability semantics inside the
pCID-selected payload protocol, not the base envelope. Source: `DI-kafiz`.

## Design Under Test

Alice issues a token that represents Alice's promise to perform a scoped action
for a holder or named promisee. Bob may present or transfer evidence of that
promise, but the token does not command Alice and does not grant global
permission:

- Alice signs only Alice's own promise and the redemption rules Alice is willing
  to honor.
- Bob's transfer is Bob's own promise or observation about how Bob handled the
  token.
- Carol locally decides whether Alice's promise and Bob's transfer history are
  trustworthy enough to rely on.
- Freeze, expiry, revocation, and transfer rules are payload-protocol semantics
  named by this simulation's future pCID.

## Local Draft Spec

The local draft spec in
`protocols/capability-promise-token.d/specs/capability-promise-token-draft.md`
defines a candidate payload protocol for capability-as-promise-token behavior.
The draft does not freeze a pCID and does not define a bearer authority object.

## Boundaries

This simulation does not settle PromiseGrid's final capability-token design. It
only gives GA scoring and guide-resource prose a PT-clean specimen for testing
capability pressure without generic assertion machinery.
