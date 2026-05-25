# Kernel implementation promise record draft

This draft is a simulation-local specimen for `SIM-fovip`. It is not a frozen
PromiseGrid protocol spec. Its purpose is to make the `DR-davod` question
testable by forcing a port to publish the promises, assumptions, unsupported
features, and evidence records that make a first implementation credible. Source:
`DI-funaf`; `DI-fidot`.

## Role

A kernel implementation promise record says what a PromiseGrid implementation
promises to local apps and operators, what it depends on from the host runtime,
and what it explicitly does not promise.

The record is not a global certificate, permission, or authority. Each receiving
agent still assesses the record and the port's make/break history locally.

## Record shape

```text
kernel_implementation_promise_record = [
  record_pcid,
  port_identity,
  profile,
  supported_pcids,
  app_facing_promises,
  host_assumptions,
  unsupported_features,
  evidence_policy,
  adapter_promises,
  namespace_promises,
  reference_promises
]
```

## Fields

- `record_pcid` is the Protocol CID for this promise-record shape.
- `port_identity` names the local implementation or agent making the promises.
- `profile` names the runtime class: native node, browser/WASM, mobile sandbox,
  MCU/header-only, split local services, or another pCID-defined profile.
- `supported_pcids` lists the pCID-selected protocols the implementation promises to
  parse, dispatch, validate, or preserve.
- `app_facing_promises` states what the implementation promises for storage,
  compute, network send/receive, key use, device access, lifecycle, pCID
  dispatch, refusal, receipt, evidence recording, namespace projection,
  reference resolution, and resource checkpoint behavior.
- `host_assumptions` states what the port depends on from a browser, OS, mobile
  sandbox, language runtime, hardware platform, or local service graph.
- `unsupported_features` states what the port refuses or cannot perform.
- `evidence_policy` states what exact bytes and local records the implementation promises
  to keep for kept, refused, unavailable, and broken promises.
- `adapter_promises` states which local APIs wrap which pCID-selected messages
  and what evidence the adapter records.
- `namespace_promises` states whether the port projects voluntary group
  namespaces and which promisers maintain those namespace frontiers.
- `reference_promises` states how the port handles CID-rooted promise-bound
  references and local path mounting.

## Promise rules

- An implementation promises only its own behavior.
- An implementation may cite host/runtime assumptions, but it does not promise
  that the host will keep them unless the host is also an explicit promiser.
- Unsupported features must be named directly. They must not be hidden behind a
  generic "partial implementation" label.
- Evidence records are local. They help Alice, Bob, Carol, and future agents
  update their own trust judgments; they are not a global trust authority.
- Local APIs are adapters. If an API call crosses a PromiseGrid promise
  boundary, the record must identify the corresponding pCID-selected message or
  state that the operation is outside the PromiseGrid boundary.
- Voluntary group namespaces are promises among agents. The record must not
  describe a namespace as universal truth.
- File-like resources are projections over promises, logs, and checkpoints. The
  record must say what evidence is preserved for the selected frontier.

## Minimum credible first port

A first credible port is allowed to be small. It must still publish:

- at least one supported pCID or a bounded exact-byte carriage profile;
- clear unsupported-pCID behavior;
- app-facing promises for every operation it exposes;
- host/runtime assumptions for every operation it delegates;
- evidence records for kept, refused, unavailable, and broken promises;
- adapter, namespace, reference, and checkpoint promises where the port exposes
  those surfaces;
- an implementation promise record that can be compared with later behavior.

## Scenario pressure

The same record shape must be tested against:

- a native node with broad local services;
- a browser/WASM host with delegated storage, network, key, and lifecycle
  behavior;
- a mobile sandbox with restricted background execution;
- an MCU/header-only port with one pCID and bounded evidence;
- a split local service graph with multiple local promisers;
- a voluntary group namespace maintained by Alice, Bob, and Carol;
- a CID-rooted promise-bound reference from Alice that Bob mounts locally;
- a file/resource checkpoint reconstructed from a selected promise-log frontier.
