# TE-ragin: Kernel resource protection by conditional capability promises

TE ID: TE-ragin

## Status

decided

## Decision under test

PromiseGrid needs a Promise-Theory-compatible way to describe a local kernel
role that can throttle or terminate local processes that abuse CPU, RAM, network,
storage, device, or lifecycle resources.

The risk is vocabulary drift. If this is described as authorization,
permission, enforcement, command, compliance, or punishment, the kernel starts to
look like a traditional authority over apps. If it is described as ordinary
conditional promises about local resources, it remains compatible with the
PromiseGrid rule that each agent promises only its own behavior.

## Assumptions

- PromiseGrid kernels are non-monolithic local role sets, not necessarily one
  daemon or one operating-system kernel.
- A lifecycle/resource supervisor is a kernel role when it promises local
  process lifetime, CPU, RAM, device, socket, storage, or shutdown behavior.
- Apps and local kernel roles are autonomous promisers. A local app may keep,
  refuse, be unable to keep, or break a lifecycle/resource-use promise.
- Promise Theory treats impositions and requests as directionals that may
  influence an autonomous agent, not as deterministic control.
- Promise Theory allows agents to promise boundary conditions over resources
  they control. A local kernel role can therefore promise how it will allocate,
  limit, withdraw, or terminate access to resources it controls.
- Capability tokens are issuer promises, not central authorization objects.

## Alternatives

### Alternative A: Kernel enforcement authority

The kernel is described as enforcing compliance with local policy. Apps receive
permissions, and the kernel punishes or kills non-compliant apps.

This matches conventional operating-system vocabulary, but it is a poor fit for
PromiseGrid. It encourages command/control language, hides reciprocal promises,
and makes the kernel appear to promise what the app will do.

### Alternative B: Harness-only supervisor

The Docker supervisor is treated as POC harness machinery. Shutdown, throttling,
and process termination are excluded from the PromiseGrid kernel model.

This avoids overdesigning POC16, but it leaves a real production obligation
undocumented. Rich hosts, browsers, WASM runtimes, mobile systems, and MCUs all
need some role that promises local resource protection.

### Alternative C: Conditional capability promises

The lifecycle/resource kernel role issues conditional capability promises. A
token or local capability promise says, in effect: "I promise this app access to
this local resource under these bounds, in return for the app promising bounded
use and lifecycle cooperation."

If the app breaks its reciprocal promise or the kernel's local resource budget
changes, the kernel role may narrow, revoke, throttle, or terminate access. This
is not a command over the app. It is the kernel role changing what it promises
about resources it controls.

## Scenario analysis

### S1: Normal local resource use

Alice's app receives a conditional capability promise for CPU, RAM, and a local
CAS write budget. The app promises bounded use. The lifecycle/resource role
observes ordinary use, keeps the capability live, and records a local event that
the resource promise remained within terms.

Alternative C makes the reciprocal promises visible. Alternative A would call the
same outcome compliance. Alternative B would hide the resource promise outside
the kernel model.

### S2: CPU or RAM abuse

Alice's app loops or allocates beyond its promised budget. The app has not been
forced to behave; it has simply stopped matching the condition under which the
kernel role promised resource access. The lifecycle/resource role throttles CPU,
reduces memory, pauses the process, or terminates it to keep its own promise to
protect the node and other local apps.

Alternative C preserves autonomy: the app's promise broke, and the kernel's
resource promise changed. Alternative A falsely suggests the kernel made the app
comply. Alternative B cannot explain why resource protection is part of the
kernel model.

### S3: Ignored shutdown request

The lifecycle/resource role sends a local lifecycle promise message asking apps
to quiesce, drain, close sessions, flush final local event records, and finish.
Most apps reciprocate. Bob's app ignores the request and keeps opening work.

Under Alternative C, the lifecycle role records Bob's broken lifecycle promise
and withdraws the process-lifetime capability before escalating to termination.
SIGTERM or host kill is the mechanical means by which the local resource promise
is withdrawn, not a claim that Bob's app was under command authority.

### S4: Wedged parser or transport role

A parser role is a kernel role, but it is still an autonomous local process. If
it wedges during shutdown, the lifecycle/resource role cannot make it keep its
promise. It can only record the local break, withdraw CPU/RAM/socket access, and
protect the rest of the node.

This is important for POC16: analyzer gates should distinguish "role did not
promise / role broke local lifecycle promise" from "supervisor authority forced
success."

### S5: Malicious app

Mallory's app receives a bearer or holder-bound local resource token and then
attempts prompt injection, malformed CBOR floods, replay, or expensive compute
requests. The lifecycle/resource role may rate-limit, revoke, or kill because it
promised resource protection to the local operator and other apps.

Trust consequences remain local. No global authority is created, and no peer is
required to share the same judgment of Mallory.

### S6: Runtime portability

In Docker, withdrawal may be SIGTERM/SIGKILL. In WASM, it may be fuel limits,
memory limits, or host-call denial. In stdio, it may be pipe closure. In a
browser tab, it may be worker termination or quota denial. On an MCU, it may be
not scheduling a loop, resetting a peripheral, or refusing a buffer.

Alternative C is portable because the promise is abstract: conditional resource
access under local resource protection. The host mechanism varies by runtime.

## Conclusions

- Adopt Alternative C.
- A local supervisor that owns process lifetime, CPU, RAM, resource quotas, and
  shutdown sequencing is a lifecycle/resource kernel role.
- "Conditional capability promise" is the preferred vocabulary.
- A local resource token is an issuer promise of access under terms, usually in
  return for a reciprocal promise of bounded use, quiescence, drain, or other
  lifecycle cooperation.
- Throttle, pause, revoke, close, or kill are host mechanisms for withdrawing or
  narrowing the kernel role's own resource promise.
- These mechanisms must not be described as central authorization, permission,
  compliance, enforcement, punishment, or global trust authority.

## Implications for open TODOs and DIs

- `TODO-binag` should record the vocabulary and update the kernel role/profile
  synthesis.
- `DN-lujad` should describe lifecycle/resource protection as a first-class
  kernel role surface.
- POC16 docs may use this as conceptual guidance, but active embedded pCID
  specdocs should not be edited in the same pass because their bytes derive
  pCIDs.
- A future `local_lifecycle_v1` pCID can implement this by carrying quiesce,
  drain, terminal-summary, flush, resource-budget, and capability-withdrawal
  promise messages.

## Decision status

Locked by `DI-vuruz`: use "conditional capability promise" for kernel-local
resource access and protection. The lifecycle/resource kernel role may withdraw,
throttle, or terminate local resource access to keep its own resource-protection
promises, without claiming authority over another autonomous agent.

