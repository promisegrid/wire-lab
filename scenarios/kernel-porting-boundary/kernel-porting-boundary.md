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
the kernel/runtime terminology remains unsettled.

## Setup

Alice wants to port PromiseGrid infrastructure to a new host environment. Bob
offers a library, Carol offers a dispatcher/runtime, and Mallory claims that
copying the wire-lab harness is the porting target. The available specs are
drafts, and only some future frozen pCIDs will become obligations.

## Stimulus

Alice writes a first porting plan and a conformance claim. The plan must say what
it implements now, what draft evidence it follows, and what it defers until
`DR-davod` closes.

## Expected Pressure

The candidate design must separate harness apparatus from porting target,
identify which binding/session/message/CAS/runtime obligations are provisional,
and preserve a clear path to future frozen-spec conformance.

## Scenario-Specific Evaluation Questions

- Should the guide say kernel, runtime, dispatcher, handler host, or library?
- What is the minimum viable porting target before final freeze?
- Which K1-K5 ingress, feed, CAS, session, and app-layer details should remain
  blocked versus provisional orientation?
