# TE-mazop: Kernel promise boundary and minimum port contract

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-mazop

## Status

needs DF

## Decision under test

What additional evidence is needed before `DR-davod` can define a stable
PromiseGrid kernel-developer porting boundary?

`TE-jimar` already rejected one universal process shape and recommended treating
daemon, microkernel, host-runtime, MCU/library, and split-object deployments as
profiles. This TE tests the part `TE-jimar` did not make concrete enough: the
local promises between apps, local infrastructure, host runtimes, storage,
compute, networking, keys, devices, and evidence records, plus the minimum
credible first-port contract.

This TE does not close `DR-davod`. It creates the next unanswered packet and the
shape of the simulation needed to answer it.

## Assumptions and constraints

- PromiseGrid uses Promise Theory vocabulary. Each agent promises only its own
  behavior; no kernel, host, app, or peer promises on behalf of another agent.
- Trust is local and relationship-specific. A local infrastructure surface may
  help an agent record evidence, but it does not become a trust authority.
- The current outer envelope direction is `grid([42(pCID), payload, ...])`.
  Slot `0` names the protocol spec, slot `1` is the payload anchor, and later
  outer slots are defined by the protocol named by `pCID`.
- A credible port can be partial, but it must say which pCID-selected protocols
  it promises to handle, which host promises it depends on, and which features
  it does not promise.
- Carriage is not semantic acceptance. A port may preserve unknown bytes as
  local evidence without promising to interpret them.
- `DR-tumus` still owns the concrete L6 CAS profile. This TE may mention CAS
  roles, but it must not require IPLD, CAR, selectors, or a specific CAS profile
  for every kernel/runtime port.

## Candidate port contracts

### Alt 1 - Message-only dispatcher

The port promises to parse PromiseGrid envelopes, recover pCID, route known
pCIDs to handlers, preserve exact bytes where local policy allows, and return
clear unsupported-pCID evidence for everything else.

**Easier:** Smallest credible contract for constrained devices, bootstrap tools,
and language ports.

**Harder:** Too thin for apps that need storage, computation, networking,
lifecycle, keys, or evidence retention from the local infrastructure.

**New obligations:** The port must say exactly which pCIDs are supported, what
happens to unsupported pCIDs, and how exact bytes are preserved or refused.

### Alt 2 - App-host promise boundary

The port is the local boundary where apps ask for storage, compute, networking,
key use, device access, and evidence recording. It does not command apps or
hosts. It makes local promises about what it will attempt, what evidence it will
record, and what it does not support.

**Easier:** Gives guide writers the practical app/kernel surface they need
without requiring one daemon shape.

**Harder:** Needs careful vocabulary so "host", "kernel", "runtime", and "app"
do not collapse into one authority.

**New obligations:** Each app-facing operation needs a promise statement, a
failure/refusal shape, and an evidence record shape.

### Alt 3 - Full local node profile

The port promises a broader local profile: pCID dispatch, feed adapter, CAS
store, key custody or key delegation, lifecycle supervision, execution host,
network adapter, and evidence recorder.

**Easier:** Closest to a rich native node and easiest for operator guidance.

**Harder:** Overstates the baseline for browser, WASM, mobile, and MCU ports.

**New obligations:** The profile must separate required roles from optional
roles and make each role's promises independently testable.

### Alt 4 - Declared-profile port contract

The port publishes a small implementation promise record that names its profile,
supported pCIDs, app-facing promises, host/runtime assumptions, unsupported
features, evidence records, and known failure modes. The implementation may be
a daemon, a microkernel service graph, a host boundary, a library, or a split
local mesh.

**Easier:** Makes partial ports honest and comparable. It preserves small-device
and host-runtime portability while giving guide writers a concrete target.

**Harder:** Requires a profile vocabulary and enough examples to avoid becoming
a vague checklist.

**New obligations:** The repo needs a simulation that forces candidate ports to
write their promised profile explicitly and then tests it against app scenarios.

## Scenario analysis

### S1 - Browser/WASM app

Alice ships a browser app. The browser controls filesystem access, network
access, clock behavior, lifecycle, sandboxing, and sometimes key access. The
PromiseGrid code can parse envelopes, dispatch known pCIDs, ask the browser for
storage/network operations, and record what it attempted.

- **Alt 1** works only for pure message handling.
- **Alt 2** fits the real boundary: the app and local PromiseGrid code promise
  what they will request and record, while the browser's behavior is a local
  host assumption.
- **Alt 3** over-promises if it implies feed, CAS, keys, and lifecycle are all
  under PromiseGrid control.
- **Alt 4** is strongest if the port declares "browser/WASM profile" and lists
  storage, network, and key assumptions explicitly.

The evidence gap is not process shape. It is whether the port can make clear
app-facing promises without pretending to control the host.

### S2 - Native multi-app node

Bob runs a native service for several local apps. It has durable storage,
network listeners, local keys, feed adapters, a CAS role, and an evidence
recorder. Apps ask it to send, receive, store, compute, and sign under specific
pCIDs.

- **Alt 1** is too thin.
- **Alt 2** captures the app-facing promise surface.
- **Alt 3** is useful as a rich-node profile.
- **Alt 4** still matters because Bob should declare which roles are actually
  implemented and which pCIDs the node promises to handle.

The guide needs examples where Bob's node refuses a request it does not promise
to handle, records that refusal, and does not present the refusal as a global
authorization decision.

### S3 - MCU/header-only device

Carol ports one actuator protocol to an MCU. The firmware can parse the current
envelope, verify a small pCID-specific proof, drive one hardware pin, and store
a bounded evidence counter.

- **Alt 1** is plausible.
- **Alt 2** must shrink to local promises about device access and exact-byte
  evidence.
- **Alt 3** is too large.
- **Alt 4** is useful if "MCU profile" can explicitly say "no CAS, no general
  feed, one pCID, bounded evidence, one device promise."

The port should not fail because it lacks daemon features. It should fail only
if it promises more than it can keep or cannot preserve the evidence required by
its pCID.

### S4 - Mobile sandbox

Dave writes a mobile app. The OS may kill background tasks, restrict network
access, move storage, and broker keys. The app can make promises about what it
does while active, what evidence it stores, and which operations become host
assumptions.

- **Alt 1** does not cover app lifecycle.
- **Alt 2** fits storage/network/key/lifecycle mediation.
- **Alt 3** over-promises unless the app includes a full local node profile.
- **Alt 4** can keep the port honest by requiring explicit unsupported features
  and evidence limits.

The guide needs a way to say "this port keeps these promises only while the OS
lets it run" without making the OS a PromiseGrid agent unless a host adapter
itself publishes promises.

### S5 - Split local services

Ellen splits dispatch, storage, key use, feed sync, and evidence recording into
separate local services. Each service can fail independently and can publish its
own local implementation promises.

- **Alt 1** captures only dispatch.
- **Alt 2** needs to show which service makes which app-facing promise.
- **Alt 3** is the co-located version of the same roles.
- **Alt 4** is strongest if profile declarations can include multiple local
  promisers and their evidence handoff promises.

The kernel cannot be one authority here. It is a set of local promises across
local agents that an app and operator assess.

### S6 - Long-horizon reimplementation

Frank reimplements a minimal port decades later. Some host APIs, CAS tooling,
and programming languages have changed, but the old pCID specs, exact bytes,
and implementation promise records remain.

- **Alt 1** gives the easiest archival bootstrap.
- **Alt 2** helps Frank reconstruct app-facing behavior.
- **Alt 3** may be too tied to historical service assumptions.
- **Alt 4** gives the best audit trail if the old implementation's profile and
  unsupported features were explicit.

The 100-year test favors profile declarations and exact evidence over one
runtime architecture.

## Cross-cutting findings

- The most useful next evidence is a concrete port profile record, not another
  abstract term choice.
- App-facing promises should cover storage, compute, network, key use, device
  access, lifecycle, pCID dispatch, exact-byte evidence, and unsupported
  features.
- Host/runtime assumptions must stay separate from promises made by the
  PromiseGrid port.
- A first credible port can be partial if it names its pCID coverage and
  unsupported features honestly.
- Broken-promise evidence must be local and exact enough for Alice, Bob, Carol,
  and later agents to update their own trust judgments.

## Conclusions

- `DR-davod` should not close from `TE-jimar` alone.
- The next DF packet should ask for evidence from a concrete simulation rather
  than asking Steve to choose a kernel term now.
- The strongest candidate is **Alt 4 - Declared-profile port contract**, with
  Alt 2 as the practical app-facing promise surface and Alt 1/Alt 3 as profile
  endpoints.
- The next simulation should force candidate ports to publish profile promises,
  host assumptions, unsupported features, pCID coverage, and evidence records.

## Recommended next DF packet for DR-davod

Before deciding the stable kernel-developer porting boundary, answer these
simulation-backed questions:

1. What is the smallest profile declaration that makes a first PromiseGrid port
   credible?
2. Which app-facing promises are required for storage, compute, networking,
   key use, device access, lifecycle, pCID dispatch, and evidence recording?
3. Which host/runtime assumptions must be stated separately from PromiseGrid
   promises?
4. What must a port record when it keeps, refuses, cannot perform, or breaks an
   app-facing promise?
5. Which pCID-selected specs are enough for a first port to claim support, and
   how does the port name unsupported pCIDs or unsupported roles?

## Implications for open work

- `DR-davod` remains open and should point to this TE plus `SIM-fovip` as the
  next evidence path.
- `SIM-funas` remains the older question home. It should not be treated as the
  concrete answer to the DR.
- `SIM-fovip` should become the concrete simulation for the minimum port
  contract.
- `DEV-GUIDE-RESOURCES.md` should continue to mark kernel guidance provisional
  until `SIM-fovip` evidence is reviewed and a later DI closes or narrows
  `DR-davod`.

## Decision status

`needs DF` - this TE narrows the next evidence packet but does not decide the
stable kernel-developer porting boundary.

## Refinements

### 2026-05-25 - App/kernel grid-message boundary

`TE-pudiv` follows this TE with a narrower question: whether the app/kernel
boundary should use the same `grid([42(pCID), payload, ...])` message format as
the wire boundary. This is a forward pointer only; the analysis above remains
unchanged historical evidence. Source: `DI-gumum`.
