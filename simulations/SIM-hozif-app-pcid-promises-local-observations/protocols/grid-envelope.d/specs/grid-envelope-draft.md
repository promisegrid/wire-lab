# Grid-envelope draft for app pCID promises

> Status: DRAFT. Simulation-local specimen only. Source: `DI-dikat`.

## Shape

Messages in this simulation use the current envelope direction:

```text
grid([42(pCID), payload, ...])
```

`pCID` is Protocol CID: the content identifier of the protocol specification document. Slot `0` is an IPLD tag-42 CID link for ecosystem compatibility. Slot `1` is interpreted only by the protocol named by `pCID`. Later slots, including proof bytes, are owned by that same protocol.

## Promise carried by the envelope

The current sender's envelope is evidence for the sender's own scoped promise:

> I promise that this payload is shaped according to the protocol specification named by this pCID, and any later slots are used as that protocol specifies.

The envelope does not claim that a receiver accepts the payload, trusts the sender, or will perform any requested action. A receiver that does not recognize the pCID may preserve exact bytes, refuse locally, or ignore the message according to its own promises and local trust.

## App/kernel boundary

The same envelope shape is used between apps and kernels in this simulation. A local API may wrap message construction or delivery, but the PromiseGrid boundary remains the pCID-selected message. The wrapper is an adapter promise by the local implementation, not a separate RPC authority.

Hozif's app/kernel promise-accounting payloads use one stable protocol pCID. A payload `kind` field distinguishes promise records from observation records inside that protocol instead of requiring separate pCIDs for tightly coupled message variants. Source: `DI-bozid`.
