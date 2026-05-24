# Promise-Accounting Evidence Draft

> **Status: DRAFT.** Not frozen. The pCID for this spec is not yet minted.
> Source: `DI-kafiz`.

## Purpose

This higher-layer payload protocol lets autonomous agents publish compact
evidence about their own promises, refusals, observations, and timeouts. It is
not a base-envelope extension and not a universal assertion artifact.

## Candidate Payload Shape

```cddl
promise-accounting-evidence = [
  observer: agent-ref,
  subject: agent-ref,
  protocol_pcid: cid,
  evidence_kind: "promise" / "refusal" / "observation" / "timeout",
  observed_bytes: bstr,
  observer_note: bstr / nil,
  signer_proof: bstr / nil
]
```

## Promise-Theory Rules

- The signer only signs the signer's own promise, refusal, or observation.
- A timeout is Alice's local observation that Bob did not answer in Alice's
  chosen window; it is not proof of Bob's intent.
- Exact bytes are retained so later peers can reassess evidence without trusting
  Alice's summary.
- Carol may use Alice's evidence, ignore it, or ask Bob for fresh evidence; no
  artifact decides trust for Carol.

## Non-Goals

This draft does not define rejected multi-selector envelope stacks or generic
assertion machinery.
