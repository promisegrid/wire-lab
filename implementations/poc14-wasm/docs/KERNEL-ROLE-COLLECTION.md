# Kernel As Role Collection

POC14 event suggests that "kernel" should be treated as a collection of local
roles, not one mandatory monolithic process. Some runtimes can run those roles in
one process; others should split them across processes, objects, firmware
functions, browser APIs, or host adapters. Source: `DI-galin`; `DI-pohaj`;
`DI-vutok`; `DI-pamob`.

## Roles Visible In POC14

1. **Transport role.** The local `poc14-kernel` accepts length-framed TCP from
   local apps, forwards exact CBOR envelopes to peer kernels, and records
   transport outcomes. It does not judge trust for the app.
2. **App-boundary role.** Local apps register receive promises by pCID. The
   kernel can deliver exact bytes only to apps that made local receive promises
   for that pCID.
3. **pCID routing role.** The kernel parses slot 0 `42(pCID)` and chooses the
   local app receive stream that promised that pCID. Payload interpretation stays
   with the app/protocol handler.
4. **Local resource role.** `printer_port` behaves like a kernel-adjacent local
   resource owner: it promises scoped future access tokens and later promises
   print event records after token redemption. It does not grant global permission.
5. **Boundary adapter role.** Victor's stdio adapter and Peggy's WASM process
   show that a runtime boundary can be a local role that translates process I/O
   into exact PromiseGrid envelopes without inventing RPC commands.
6. **Trust/workflow role.** Trust judgment, relationship ledgers, storage,
   compute, shipping workflow, and keep/break interpretation stay in apps, not in
   the transport kernel.

## Why This Matters

The common denominator is not a kernel binary. The common denominator is the
ability to send, receive, sign, parse, route, remember, and evaluate
`grid([42(pCID), payload, proof])` promise envelopes according to local promises.

That means:

- A Docker runtime may run one kernel process plus many app processes.
- A WASM host may expose only an envelope send/receive boundary to the sandbox.
- A stdio subprocess may never see the network and still exchange exact
  PromiseGrid envelopes through an adapter.
- A microcontroller may collapse transport, app-boundary, and local-resource
  roles into firmware functions.
- A production node owned by a legal entity may split roles into separate
  services for operations, safety, or resource isolation.

## POC15 Implications

POC15 should make this role split more explicit:

- Transport role: direct TCP and later peer-to-peer forwarding.
- App-boundary role: local app registration and receive queues by pCID.
- Routing role: multi-hop route-promise selection without global authority.
- Local-resource role: device/storage/compute capability promises.
- Event role: app-owned local journals and optional voluntary summaries.

The analyzer can inspect all roles in a POC run, but production agents cannot
assume a global analyzer. Production-shaped monitoring must be ordinary local
event promises, peer-carried attestations, and economic/trust signals such as
bearer-token exchange-rate offers.
