# Claim shapes

## Honest partial-conformance wording

Use wording of this form:

`This embodiment conforms only to the contract subset listed here for App Manifest <AM CID>. Live transport, local cache durability, and any storage, retention, authorization, or physical-effect guarantees not listed here are outside this claim. Durable PromiseGrid evidence is published only through Savepoint Audit Envelopes that cite the object IDs listed below.`

## Required interpretation rules

- Same app means same `AM CID` or a declared successor chain.
- Authoritative protocol-boundary identity means the current signing key validated by the `ICR` chain.
- Display names, usernames, colors, and local adapter IDs are presentation hints only.
- A content hash alone identifies content; it does not prove replication, retention, readability, or authorization.
- Live-channel behavior is never implied by an audit publication unless the EC explicitly claims it, which this profile does not.

## Minimal field checklist

### AM
- app label
- contract-family reference
- embodiment classes
- shared semantics
- excluded semantics
- successor reference if any

### EC
- AM reference
- embodiment class
- implemented subset
- runtime limits
- storage limits
- excluded guarantees
- signing key reference

### ICR
- prior key
- next key
- scope
- reason
- signatures or break-witness route

### SAE
- AM reference
- EC refs
- ICR refs
- cited durable object ref
- human-readable promise body
- statement excluding live transport from claim scope

## Audit review shortcut

A third party can review an embodiment by checking, in order:
1. AM for what app is being claimed.
2. EC for what this embodiment actually implements.
3. ICR chain for which key is authoritative now.
4. SAE for which durable object was actually published and what promise was made about it.