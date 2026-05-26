# TE-lamun: pCID-defined slot vector

*Thought experiment, part of the [PromiseGrid Wire Lab](../../protocols/wire-lab.d/specs/harness-spec-draft.md). This file is content-addressable; its hash is its pCID.*

## TE ID

TE-lamun

## Status

decided, refined

## Decision under test

`TE-fikoj` locked the current envelope direction as
`grid([42(pCID), payload, ...])`: slot `0` is `42(pCID)`, slot `1` is the
payload anchor, and later outer slots are defined by the protocol named by
`pCID`. That decision was correct for the immediate fixed-three-slot versus
variable-arity question, but it left a sharper generalization unresolved:

Should the universal PromiseGrid envelope rule say that slot `1` is always the
payload, or should it say that slot `0` selects the protocol and the protocol
named by `pCID` defines every following slot?

The candidate general rule is:

```text
grid([42(pCID), ...protocol-defined-slots])
```

The candidate profile remains:

```text
grid([42(pCID), payload, ...])
```

This TE tests whether the general rule improves the 100-year envelope story
without collapsing into "anything goes" or weakening developer guidance.

## Assumptions

- `pCID` means Protocol CID: the CID of the protocol specification document,
  never the CID of a payload object.
- Slot `0` remains `42(pCID)` in the current standard envelope instance.
- The whole message is a CBOR array, so the array header carries arity.
- The protocol named by `pCID` already defines payload shape, proof encoding,
  signable view, validation rules, unsupported behavior, and promise semantics.
- PromiseGrid objects help agents make, recognize, remember, and evaluate
  promises; they do not command other agents or create global trust.
- Receivers that do not support a `pCID` may preserve exact bytes under local
  policy, but they do not semantically accept or reinterpret the message.
- Most protocols still benefit from the simple profile where slot `1` is the
  primary payload/body anchor.

## Alternatives

- **Alt-A - Hard slot-1 payload rule.** Keep `grid([42(pCID), payload, ...])`
  as both the formal envelope rule and the recommended profile. Slot `1` is
  always the payload, and slots `2..N` are protocol-defined.
- **Alt-B - pCID-defined slot vector with payload-example profile.** Make the
  formal envelope rule `grid([42(pCID), ...protocol-defined-slots])`. The
  protocol named by `pCID` defines slots `1..N`. The recommended profile remains
  `grid([42(pCID), payload, ...])`, and deviations must be justified in the
  protocol spec.
- **Alt-C - Unconstrained pCID-owned layout.** Make the formal rule
  `grid([42(pCID), ...])` and avoid any recommended slot-1 payload convention.
  Every protocol chooses its own slot vector without repo-level profile
  pressure.

## Scenario analysis

### Scenario 1 - Alice runs a constrained receiver

Alice has a tiny CBOR receiver. She can parse arrays, tag `42`, the CID bytes,
and a few known `pCID`s. She cannot afford a large dynamic object system.

- **Alt-A** is easiest to explain: after slot `0`, Alice knows slot `1` is the
  payload. If the protocol is unsupported, she can still label the exact bytes
  as a PromiseGrid envelope with a payload-shaped second slot.
- **Alt-B** makes Alice consult the supported protocol's spec before naming slot
  `1`, but that cost only applies when she supports the protocol. If she does
  not support it, she still has a stable slot `0` selector and exact bytes.
- **Alt-C** gives Alice no default shape to expect after slot `0`; every
  supported protocol is a one-off parser.

Finding: Alt-A has the simplest universal mental model, but Alt-B is still
small-device compatible because slot `0` stays universal and supported pCIDs
already require protocol-specific code. Alt-C adds avoidable cognitive and code
surface area.

### Scenario 2 - Bob writes a normal application protocol

Bob defines a plain application protocol for sending promise bodies, replies,
and local evidence. He wants other implementers to read his spec quickly.

- **Alt-A** guides Bob toward `[42(pCID), payload, ...]`, which is probably what
  he should use.
- **Alt-B** also guides Bob toward `[42(pCID), payload, ...]`, but states that
  this is the profile, not a universal law. If Bob follows the profile, his
  protocol remains boring and easy to implement.
- **Alt-C** gives Bob too much freedom and too little pressure to preserve the
  common shape.

Finding: Bob benefits from a strong profile. Alt-B keeps that profile without
  pretending every future protocol must have the same slot `1` semantics.

### Scenario 3 - Carol defines a proof-first compact protocol

Carol has a protocol where a compact proof or selector must be inspected before
the body can be safely parsed. The protocol still uses slot `0` as
`42(pCID)`, and the spec clearly says slot `1` is proof metadata while slot `2`
is the body.

- **Alt-A** forces Carol either to call proof metadata a "payload" even though
  that is misleading, or to bury the proof metadata inside a slot-1 payload
  wrapper solely to satisfy the universal rule.
- **Alt-B** lets Carol specify `[42(pCID), proof, body]`, but requires the spec
  to define the slot roles, signable view, and validation order. Generic
  readers still recover slot `0` and know that interpretation begins at the
  pCID spec.
- **Alt-C** also permits Carol's shape, but lacks the "deviation from
  payload-example must be justified" discipline.

Finding: Alt-B is cleaner. It avoids a fake payload wrapper while keeping the
  obligation to specify exactly what the slots mean.

### Scenario 4 - Dave maintains a generic relay

Dave relays messages for peers he trusts locally. He does not interpret most
payloads, but he wants useful generic diagnostics.

- **Alt-A** lets Dave say "unknown pCID, slot 1 appears to be payload" even when
  he cannot parse it.
- **Alt-B** lets Dave say "unknown pCID, exact bytes carried; slot 0 names the
  protocol; slots 1..N are uninterpreted until that protocol is supported."
  This is less semantically presumptive and better matches carriage-only
  behavior.
- **Alt-C** is similar to Alt-B for unknown pCIDs, but gives Dave less profile
  guidance for common supported protocols.

Finding: Alt-B gives the most honest generic relay behavior. Carriage is not
semantic acceptance, so Dave should not claim slot `1` is payload unless the
named protocol says so.

### Scenario 5 - Ellen reads a century-old archive

Ellen finds archived PromiseGrid messages decades after the original tools have
faded. She can reconstruct CBOR, tag `42`, CIDs, and some old protocol specs.

- **Alt-A** gives Ellen a default archaeology story: slot `1` is payload unless
  she has reason to think otherwise.
- **Alt-B** gives Ellen a more precise archaeology story: slot `0` identifies
  the spec; that spec is the source of truth for every following slot; if the
  spec used the normal profile, slot `1` is payload.
- **Alt-C** gives Ellen no conventional clue about common protocols.

Finding: Alt-B is the best long-horizon rule because the pCID spec remains the
source of truth. The payload-example profile still helps common cases.

### Scenario 6 - Frank scores or generates future sims

Frank runs the GA scorer and generator. He wants future designs to be simple,
Promise-Theory-clean, and not overfit to one envelope specimen.

- **Alt-A** may over-penalize a protocol that has a legitimate reason to place
  proof, routing wrapper, or negotiation bytes before the primary body.
- **Alt-B** gives the scorer a better rule: reward the default
  `[42(pCID), payload, ...]` profile when it fits, but do not penalize a
  different slot vector if the pCID spec clearly defines the promise semantics
  and validation rules.
- **Alt-C** makes scoring harder because there is no default simplicity profile
  to reward.

Finding: Alt-B gives GA work the right bias. It is pro-simple by default but
not anti-protocol-owned slot vectors.

## Cross-cutting findings

### Slot 0 is the universal bootstrap point

The receiver's first durable obligation is to recover `42(pCID)` from slot `0`.
That is the universal shape PromiseGrid should protect most strongly.

### The pCID spec is already the source of truth

PromiseGrid already relies on the pCID-named spec for payload shape, proof
encoding, signable views, and validation. Extending that authority to all
following slot roles is coherent rather than new machinery.

### An example profile is still valuable

Most protocols should not be clever. `grid([42(pCID), payload, ...])` remains
the example profile because it is easy to teach, easy to debug, and compatible
with the existing `TE-fikoj` lineage.

### "Protocol-defined" is not "anything goes"

A protocol that does not use slot `1` as the primary payload/body anchor must
say why. Its spec must define slot count, slot order, slot meanings, signable
view, validation order, unsupported behavior, and promise semantics.

## Conclusions

- Reject Alt-A as too rigid. Slot `1` as payload is an excellent default, but
  making it a universal law forces misleading wrappers in protocols that have a
  legitimate proof-first or negotiation-first shape.
- Reject Alt-C as too loose. It throws away the practical value of the existing
  payload-example profile.
- Adopt Alt-B:
  - formal rule: `grid([42(pCID), ...protocol-defined-slots])`;
  - slot `0`: current universal protocol selector, `42(pCID)`;
  - slots `1..N`: defined by the protocol spec named by `pCID`;
  - recommended profile: `grid([42(pCID), payload, ...])`;
  - deviation from the payload-example profile must be explicit and justified
    in the protocol spec.

## Implications for open work

- `TE-fikoj` remains the historical source for universal `42(pCID)` and
  variable CBOR arity, but its slot-1 payload wording is superseded by this TE.
- `DN-jotob` should describe both the formal rule and the recommended profile.
- `DEV-GUIDE-RESOURCES.md` should stop presenting slot `1` as a universal law
  and instead call it the default primary body/payload profile.
- Existing sims that use `[42(pCID), payload, varsig]` remain valid specimens of
  the recommended profile.
- Future sims may test non-payload slot `1` layouts, but only when the protocol
  spec clearly justifies the deviation.

## Decision status

Locked by `DI-rojij` in `TODO-mopob`. `DI-rojij` supersedes `DI-sisak` only for
the hard slot-1 payload rule. It preserves universal slot `0` as `42(pCID)`,
preserves CBOR array arity, preserves pCID-owned slot semantics, and reframes
`grid([42(pCID), payload, ...])` as the recommended profile rather than the full
formal envelope law.

## Refinements

### 2026-05-26 — Example profile wording

`DI-punam` narrows the wording for the payload-first profile. Current guidance
should say `grid([42(pCID), payload, ...])` is the recommended example profile,
not an inherited slot-layout rule. This refinement preserves the formal
`grid([42(pCID), ...protocol-defined-slots])` rule and preserves `42(pCID)` in
slot `0`, but removes any implication that ordinary protocols inherit a default
slot layout before their pCID-named spec defines one.
