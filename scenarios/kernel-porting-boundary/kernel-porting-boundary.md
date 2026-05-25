# Kernel Porting Boundary

## Scenario ID

kernel-porting-boundary

## Source / Provenance

- Source type: new harness scenario
- Source path: `/home/stevegt/lab/promisegrid-dev-guide/FEEDBACK.md`
- Source row/title: `FB-vitih`, `FB-mulum`, and `FB-potin`
- Source DI / TODO / TE: `DI-ragaz`; `TODO-rozas`; `DR-davod`

## Purpose

Exercise the developer-facing boundary for a first real PromiseGrid port while
the kernel/runtime terminology remains unsettled. Source: `DI-ragaz`;
`DI-fidot`.

## Setup

Alice wants to port PromiseGrid infrastructure to a new host environment. Bob
offers a native service, Carol offers a browser/WASM adapter, Dave offers an
MCU/header-only profile, Ellen offers split local services, and Mallory claims
that copying the wire-lab harness or one microkernel shape is the porting
target. The available specs are drafts, and only some future frozen pCIDs will
become obligations.

## Stimulus

Alice writes a first porting plan and a port promise record. The plan must say
what it implements now, which pCID-selected messages it exposes, what draft
evidence it follows, what it refuses or cannot promise, what host assumptions it
depends on, what evidence it records, and what it defers until `DR-davod`
decides the guide-facing boundary.

## Expected Pressure

The candidate design must separate harness apparatus from porting target,
identify which binding/session/message/CAS/runtime obligations are provisional,
and preserve a clear path to future frozen-spec implementation promises.

It must also show whether:

- app/kernel operations are pCID-selected `grid([42(pCID), payload, ...])`
  messages, even when local APIs wrap them;
- storage, compute, network send/receive, key use, device access, lifecycle,
  dispatch, refusal, receipt, evidence, namespace, reference, and checkpoint
  operations each name their promiser;
- host/runtime assumptions are separate from PromiseGrid promises;
- unsupported pCIDs and unsupported roles are direct refusals or non-promises;
- exact bytes are preserved where needed for proof, replay, unsupported-pCID
  carriage, or broken-promise evidence;
- voluntary group namespaces are reciprocal promises, not global truth;
- CID-rooted references let Alice share a resource that Bob maps into Bob's own
  local view;
- file/resource current state can be reconstructed as a checkpoint over a
  selected promise-log frontier.

## Scenario Variants

- **Native node:** Bob's service promises storage, dispatch, networking, keys,
  lifecycle, and evidence, but must name every pCID it supports and every role it
  does not promise.
- **Browser/WASM:** Carol's adapter depends on browser storage, network, clocks,
  and lifecycle. Carol can promise adapter behavior, but not that the browser
  will keep promises the browser has not made.
- **Mobile sandbox:** Dave promises work only while the OS offers foreground or
  background execution. Unavailable background work must be recorded as an
  unavailable promise, refusal, or host assumption rather than hidden success.
- **MCU/header-only:** Erin supports one actuator pCID, one bounded evidence
  store, and no general namespace service. The port is credible only if it says
  exactly what it cannot promise.
- **Split local services:** Frank separates dispatch, storage, keys, networking,
  and evidence among local promisers. The record must show which service makes
  each promise and how evidence moves between them.
- **Voluntary namespace:** Alice, Bob, and Carol maintain `/project/report` by
  reciprocal namespace promises. Mallory's lookalike namespace is rejected unless
  a local agent trusts Mallory's promise history.
- **CID-rooted reference:** Alice sends Bob a reference rooted at a CID with
  pCID, selector/path, frontier, promise body, and evidence. Bob chooses whether
  and where to mount it.
- **Checkpointed resource:** Grace reconstructs a file from old pCID specs,
  promise/event logs, and a selected frontier. A different branch may produce a
  different current file because it selects a different promise history.

## Scenario-Specific Evaluation Questions

- Should the guide say kernel, runtime, dispatcher, handler host, or library?
- What is the minimum viable porting target before final freeze?
- Which K1-K5 ingress, feed, CAS, session, and app-layer details should remain
  blocked versus provisional orientation?
- Does the candidate preserve local trust, autonomous promisers, and
  make/break evidence?
- Does the candidate avoid global permission, global namespace authority, and
  universal process-shape assumptions?
- Are local APIs faithful adapters over pCID-selected messages, or do they hide
  the promise boundary?
