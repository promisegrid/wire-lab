# local_lifecycle_v1

## Status

POC16 implementation-local protocol specification. This document is the pCID
source text for `local_lifecycle_v1`. Normative intent is recorded in
`DI-jafoj`. Source: `DI-jafoj`.

## Abstract

`local_lifecycle_v1` describes local lifecycle promises between a supervisor and
the local roles whose processes the supervisor starts. The protocol exists so
shutdown can be modeled as promise fulfillment instead of command/control. A
role promises its own shutdown behavior by issuing a signed capability token at
startup. The supervisor may later present that exact token back to the issuer to
invoke the promise under the token terms.

This pCID is local to the POC16 runtime family. It is not a global process
control protocol and does not grant authority over remote agents.

## pCID and envelope

Messages use the PromiseGrid grid CBOR tag. This protocol is an explicit example
of the general `grid([42(pCID), ...])` rule:

```text
grid([42(local_lifecycle_v1_pCID), payload])
```

There is no generic outer proof slot for this protocol. The pCID defines all
slots after slot 0. For this protocol, slot 1 is one CBOR map payload. Any proof
needed by this protocol lives inside the payload as COSE/CWT bytes. In
particular, the startup lifecycle capability token is a CWT payload protected by
COSE_Sign1, so an additional generic envelope proof would be redundant.

## Promise Theory model

The supervisor is not a global authority and does not make promises on behalf of
apps, parser roles, or transport kernels. The supervisor controls only local
process resources it started. If a role breaks or times out on a lifecycle
promise, the supervisor may withdraw local CPU/process resources, including
SIGTERM or SIGKILL, as a local resource-protection act.

Token issuance means the role promises:

- if this exact signed token is later presented by the local supervisor;
- under the signed run, audience, pCID, role, channel, expiry, and invocation
  terms;
- then the role promises to stop starting new work, drain accepted work, close
  local persistent sessions, flush local events, emit a role summary, and exit
  voluntarily.

Token invocation means the supervisor presents the issuer's own promise token
back to the issuer. Invalid, expired, wrong-audience, wrong-run, wrong-pCID, or
replayed tokens are rejected as non-promises for that local issuer.

## Payload grammar

The slot-1 payload is a definite-length CBOR map with UTF-8 text keys. Binary
COSE/CWT token bytes are base64 encoded when placed in the map so diagnostic
tools and line-oriented event logs can display them without rewriting the token.

Required common keys:

- `kind`: one of `token_issued`, `ready`, `token_invoked`, `token_fulfilled`.
- `promiser`: local role that is making the promise in this message.
- `promisee`: local role that is receiving or presenting the promise.
- `role_id`: stable local role identifier such as `agent:alice`,
  `parser:shipping`, or `kernel:shipping`.
- `role_kind`: one of `app`, `parser`, or `kernel`.
- `channel_profile`: one of `tcp`, `parser_path`, or `stdio`.
- `run_id`: local POC run id.

`token_issued` additionally requires:

- `token_cose_b64`: base64 COSE_Sign1 bytes containing a CWT claims map.
- `token_cid`: canonical CIDv1 base32 text for the exact COSE token bytes.
- `public_key_b64`: base64 Ed25519 public key corresponding to the token issuer.

`ready` additionally requires:

- `token_cid`: token being held ready for invocation.
- `detail`: local description of readiness, for example `work_complete`.

`token_invoked` additionally requires:

- `token_cose_b64`: exact token bytes being presented.
- `token_cid`: CID of the presented token.
- `reason`: local invocation reason, for example `run_complete`.
- `deadline_unix`: Unix timestamp by which the presenting supervisor expects a
  local outcome.

`token_fulfilled` additionally requires:

- `token_cid`: token whose promise reached an outcome.
- `outcome`: one of `kept`, `broken`, `rejected`, or `timed_out`.
- `detail`: local outcome summary.

## CWT/COSE Token Profile

Lifecycle tokens are CWT payloads protected by COSE_Sign1 with Ed25519. POC16
uses a well-known Go COSE library for COSE_Sign1 construction and verification.
The token issuer is the role that promises lifecycle behavior. The token
audience is the local supervisor.

The CWT payload contains standard registered CWT fields:

- `iss`: issuer role id.
- `sub`: `local lifecycle shutdown promise`.
- `aud`: local supervisor role id.
- `iat`: issued-at Unix timestamp.
- `nbf`: not-before Unix timestamp.
- `exp`: expiration Unix timestamp.
- `cti`: token identifier bytes.

The CWT payload also contains private PromiseGrid signed promise terms:

- `pg_pcid`: binary CID bytes for this `local_lifecycle_v1` pCID.
- `pg_run_id`: local run id.
- `pg_role_kind`: `app`, `parser`, or `kernel`.
- `pg_role_id`: issuer role id.
- `pg_channel_profile`: `tcp`, `parser_path`, or `stdio`.
- `pg_shutdown_terms`: array of text terms the issuer promises, including
  quiesce, drain accepted work, close persistent sessions, flush local events,
  emit a role summary, and exit voluntarily.
- `pg_grace_ms`: local grace window expected by the issuer.
- `pg_max_invocations`: maximum accepted invocation count, normally `1`.

The token is a promise by its issuer. Possession of the token is not global
permission and is not a command. An issuer may reject a presented token when the
signature, audience, pCID, run id, expiry, not-before, role id, role kind,
channel profile, or invocation count does not match the issuer's local promise
terms.

## Sender behavior

Role senders must issue `token_issued` exactly once during startup for a given
role/run/channel profile. App roles send `ready` after their normal work is
complete and they are waiting for lifecycle invocation. Parser and kernel roles
send `ready` after startup because their normal work is to keep serving until
the lifecycle token is invoked.

The supervisor sender stores the exact token bytes and sends `token_invoked`
only by presenting those same token bytes. The supervisor must not invent a
role's token or alter token terms.

## Receiver and parser behavior

Receivers parse slot 0 as `42(local_lifecycle_v1_pCID)` and then parse slot 1
as the lifecycle map. The generic envelope parser must not require
`payload, proof` for this pCID.

Role receivers verify COSE_Sign1, decode the CWT signed promise terms, and check
issuer, audience, pCID, run id, role id, role kind, channel profile, not-before,
expiry, and invocation count before acting. Supervisor receivers perform the
same token verification when a role issues a token.

Malformed lifecycle messages are local non-promises. They do not prove global
invalidity and do not create trust authority.

## Channel Profiles

The same `local_lifecycle_v1` messages may travel over multiple local channel
profiles:

- `tcp`: the role connects to a supervisor-owned local lifecycle TCP listener.
  Startup, readiness, invocation, and fulfillment frames travel on the dedicated
  lifecycle stream.
- `parser_path`: the supervisor invokes the parser role by sending the
  `token_invoked` grid message through the parser role's existing app-facing
  pCID parser path. The parser role still sends startup, readiness, and
  fulfillment records to the supervisor over its lifecycle stream.
- `stdio`: the supervisor invokes a stdio-adapted role by writing the
  `token_invoked` grid message to the child's stdin using the same
  length-prefixed binary frame convention. The role still sends startup,
  readiness, and fulfillment records to the supervisor over its lifecycle
  stream.

These profiles are local runtime choices. They do not change the pCID, token
meaning, or Promise Theory semantics.

## Protocol state machine

```text
 role process                         local supervisor
     |                                      |
     | token_issued(COSE_Sign1(CWT))       |
     |------------------------------------->|
     |                                      | stores token by token_cid
     | ready(token_cid)                    |
     |------------------------------------->|
     |                                      |
     |              token_invoked(token)    |
     |<-------------------------------------|
     | verify token and invocation terms    |
     | quiesce/drain/flush/summarize/exit   |
     | token_fulfilled(outcome)             |
     |------------------------------------->|
     |                                      | records local outcome
```

If the supervisor cannot obtain a token, cannot invoke the token, or does not
receive fulfillment before the local deadline, the supervisor records the
promise as broken or timed out and may withdraw local process resources.

## State, CAS, DAG, and retention

Lifecycle messages are normal PromiseGrid CBOR messages and may be stored in the
same run-local CAS/DAG as other POC16 message artifacts. A lifecycle token is
identified by the CID of its exact COSE_Sign1 bytes. The lifecycle token CID is
diagnostic and local-retention-friendly; it is not a destination address.

POC16 retention is run-local. A supervisor may keep token issue/invocation/
fulfillment messages until the run is analyzed, then remove them under local
retention/GC promises. A production node may retain lifecycle records longer as
local operational history, but no peer can require another legal entity to keep
them.

## Security considerations

Lifecycle tokens must be signed CWT/COSE tokens, not opaque strings. Token
verification must reject invalid signatures, wrong issuer, wrong audience, wrong
run id, wrong pCID, wrong role id, wrong channel profile, expired tokens, early
tokens, and replay beyond `pg_max_invocations`.

The supervisor's fallback to SIGTERM/SIGKILL is local resource withdrawal. It is
not evidence that the supervisor has authority over independent agents outside
the local runtime resources it controls.

## Interoperability notes

This protocol uses CWT and COSE because those are existing CBOR-native security
formats. It uses Ed25519 COSE_Sign1 in POC16. Other lifecycle pCIDs could define
different token algorithms, encrypted tokens, detached proofs, or narrower
payload shapes for constrained runtimes, but those choices would produce
different pCID bytes.

The `tcp`, `parser_path`, and `stdio` channel profiles are runtime-local. They
are not separate network protocols and do not change the signed token meaning.

## Examples

Token issue:

```text
grid([42(local_lifecycle_v1_pCID), {
  kind: "token_issued",
  promiser: "agent:alice",
  promisee: "supervisor:alpha",
  role_id: "agent:alice",
  role_kind: "app",
  channel_profile: "tcp",
  run_id: "poc16-clean",
  token_cose_b64: "...",
  token_cid: "bafkrei...",
  public_key_b64: "..."
}])
```

Token invocation:

```text
grid([42(local_lifecycle_v1_pCID), {
  kind: "token_invoked",
  promiser: "supervisor:alpha",
  promisee: "agent:alice",
  role_id: "agent:alice",
  role_kind: "app",
  channel_profile: "tcp",
  run_id: "poc16-clean",
  token_cose_b64: "...same exact token...",
  token_cid: "bafkrei...",
  reason: "run_complete",
  deadline_unix: "1780000000"
}])
```

## Expected behavior

Issuers must:

- issue at most one active lifecycle token per role/run/channel profile;
- sign the CWT with the issuer's Ed25519 key;
- verify token invocation before acting;
- reject invalid or replayed tokens as non-promises;
- stop starting new work after valid invocation;
- drain accepted work and close local persistent sessions;
- flush local event records to stdout/stderr so the passive collector can
  observe them;
- emit a role summary before exit.

Supervisors must:

- store token bytes exactly as received;
- invoke only tokens issued for the same run and local supervisor audience;
- present the exact token bytes back to the issuer;
- record invocation and outcome locally;
- use SIGTERM/SIGKILL only after token failure or timeout.

## Analyzer Expectations

A clean POC16 run should contain, for every supervised app/parser/kernel role:

- one `local_lifecycle_token_issued` event;
- one `local_lifecycle_ready` event;
- one `local_lifecycle_token_invoked` event;
- one `local_lifecycle_token_verified` event;
- one `local_lifecycle_role_summary` event;
- one `local_lifecycle_token_fulfilled` event with outcome `kept`;
- no `local_lifecycle_sigterm_fallback_used` event;
- zero unterminated persistent sessions.
