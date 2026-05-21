# Question

What is the smallest auditable claim package that lets a porter or app author say something useful now, without collapsing runtime boundary, app semantics, and host-local behavior into one misleading conformance claim?

## Open decision points

- Which fields are mandatory in each claim card?
- Which promise-accounting record rows are mandatory for blob publish/retrieve, physical effects, audit snapshots, and key rotation?
- Should the guide standardize on `runtime`, `dispatcher`, or another term inside the Port Boundary Card while `DR-davod` remains open?
- Which current signature-carriage statements are safe as `provisional` versus `local-only`?
- Which host dependencies may be named as implementation facts without being read as PromiseGrid protocol obligations?
